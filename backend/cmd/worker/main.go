package main

import (
	"context"
	"time"

	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/reportedstream"
	"github.com/DanArmor/vtuber-go/ent/vtuber"
	"github.com/DanArmor/vtuber-go/internal/config"
	"github.com/DanArmor/vtuber-go/internal/infrastructure/task"
	"github.com/DanArmor/vtuber-go/internal/setup"
	holodexcontroller "github.com/DanArmor/vtuber-go/pkg/api/holodex_controller"
	telegramcontroller "github.com/DanArmor/vtuber-go/pkg/api/telegram_controller"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
)

type Options struct {
	ConfigPath string `long:"config" description:"Config path (with extension)" required:"true"`
}

const defaultTrustedProxy = "127.0.0.1"

func getBaseRouter(router *gin.Engine, basePath string) *gin.RouterGroup {
	if basePath != "" {
		return router.Group(basePath)
	}
	return router.Group("")
}

type VtuberGoWorker struct {
	holodexController  *holodexcontroller.HolodexController
	telegramController *telegramcontroller.TelegramController
	db                 *ent.Client
	logger             *zap.Logger
}

func (vgw *VtuberGoWorker) GetVtubersWithActiveUsers() []*ent.Vtuber {
	vtubers, err := vgw.db.Vtuber.Query().
		Where(vtuber.HasUsers()).
		All(context.Background())
	if err != nil {
		vgw.logger.Error("Get active vtubers error", zap.Error(err))
		return nil
	}
	return vtubers
}

func (vgw *VtuberGoWorker) NotifyUsers(m task.TaskInput) error {
	timeNotifyAfter := 30
	vtubers := vgw.GetVtubersWithActiveUsers()
	vgw.logger.Debug("Got vtubers with active users", zap.Int("VtubersCount", len(vtubers)))
	videos := vgw.holodexController.GetVtubersUpcomingVideos(vtubers)
	vgw.logger.Debug("Got vtubers upcoming videos", zap.Int("VideosCount", len(videos)))
	now := time.Now()
	for video_index := range videos {
		vgw.logger.Debug("Process video", zap.String("VideoID", *videos[video_index].Id))

		if videos[video_index].AvailableAt == nil {
			vgw.logger.Warn("Available time nil")
			continue
		}
		if !now.Add(time.Duration(timeNotifyAfter) * time.Minute).After(*videos[video_index].AvailableAt) {
			vgw.logger.Debug("Video is not soon to start", zap.Time("AvailableAt", *videos[video_index].AvailableAt))
			continue
		}
		if videos[video_index].Channel.Id == nil {
			vgw.logger.Warn("Nil channel id")
			continue
		}
		exist, err := vgw.db.ReportedStream.Query().
			Where(reportedstream.VideoIDEQ(videos[video_index].GetId())).
			Exist(context.Background())
		if err != nil {
			vgw.logger.Error("ReportedStreams query error", zap.Error(err))
			return err
		}
		if exist {
			vgw.logger.Debug("Stream exists")
			continue
		}
		vgw.logger.Debug(
			"ChannelID for video",
			zap.String("VideoID", *videos[video_index].Id),
			zap.String("ChannelId", *videos[video_index].Channel.Id),
		)
		vgw.logger.Debug("Query users for vtuber")
		// We don't check the len of users, because before that we've queried only vtubers with active users
		users, err := vgw.db.Vtuber.Query().
			Where(vtuber.YoutubeChannelID(*videos[video_index].Channel.Id)).
			QueryUsers().
			All(context.Background())
		if err != nil && !ent.IsNotFound(err) {
			vgw.logger.Error("Can't query users for vtuber", zap.Error(err))
			return err
		}
		if len(users) == 0 {
			vgw.logger.Warn("No users for vtuber video. Probably bug of holodex api. See TRE-54")
			continue
		}
		vgw.logger.Debug("Notify users on video in telegram", zap.Int("UsersCount", len(users)))
		vgw.telegramController.NotifyUsersOnVideo(users, videos[video_index])
		vgw.logger.Debug("Query authorId")
		authorId, err := vgw.db.Vtuber.Query().
			Where(vtuber.YoutubeChannelID(*videos[video_index].Channel.Id)).
			FirstID(context.Background())
		if err != nil {
			vgw.logger.Error("Author search error", zap.Error(err))
			return err
		}
		vgw.logger.Debug("Set ReportedStream")
		err = vgw.db.ReportedStream.Create().
			SetAuthorID(authorId).
			SetAvailableAt(*videos[video_index].AvailableAt).
			SetVideoID(*videos[video_index].Id).
			SetVtuberID(authorId).
			Exec(context.Background())
		if err != nil {
			vgw.logger.Error("Reported stream set error", zap.Error(err))
		}
		vgw.logger.Debug("Process video. Done", zap.String("VideoID", *videos[video_index].Id))
	}
	return nil
}

func main() {
	var options Options
	parser := flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		panic("Can't parse cmd args")
	}

	config, err := config.LoadConfig(options.ConfigPath)
	if err != nil {
		panic("Can't load config")
	}

	cfg := zap.NewDevelopmentConfig()

	logger, err := cfg.Build()
	if err != nil {
		panic("Can't create new logger: " + err.Error())
	}

	zap.ReplaceGlobals(logger)

	zap.L().Info("Service started")

	entClient := setup.MustDatabaseSetupNoMigrations(config.DriverName, config.SqlUrl)

	worker := task.NewTaskWorker("MainWorker", entClient, zap.L())

	holodexController := holodexcontroller.NewHolodexController(config.HolodexApiKey, zap.L())
	telegramController := telegramcontroller.NewTelegramController(config.TgBotToken, zap.L())

	vgw := VtuberGoWorker{
		holodexController:  holodexController,
		telegramController: telegramController,
		db:                 entClient,
		logger:             zap.L(),
	}

	task.RegisterTask(task.TaskDescriptor{
		Name: "notify_telegram_users",
		Func: vgw.NotifyUsers,
	})

	worker.Run(context.Background())
}
