package main

import (
	"context"
	"time"

	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/queuetask"
	"github.com/DanArmor/vtuber-go/ent/reportedstream"
	"github.com/DanArmor/vtuber-go/ent/vtuber"
	"github.com/DanArmor/vtuber-go/internal/config"
	"github.com/DanArmor/vtuber-go/internal/infrastructure/task"
	"github.com/DanArmor/vtuber-go/internal/setup"
	holodexcontroller "github.com/DanArmor/vtuber-go/pkg/api/holodex_controller"
	telegramcontroller "github.com/DanArmor/vtuber-go/pkg/api/telegram_controller"
	"github.com/DanArmor/vtuber-go/pkg/utils"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
)

type Options struct {
	ConfigPath string `long:"config" description:"Config path (with extension)" required:"true"`
}

type VtuberGoWorker struct {
	holodexController  *holodexcontroller.HolodexController
	telegramController *telegramcontroller.TelegramController
	db                 *ent.Client
	logger             *zap.Logger
}

type CleanTaskInput struct {
	TaskName string `json:"task_name"`
	After    int    `json:"after"`
}

func (vgw *VtuberGoWorker) ExtractCleanTaskInput(content task.TaskInput) (CleanTaskInput, error) {
	var step CleanTaskInput
	err := utils.JSONDecode(content, &step)
	return step, err
}

func (vgw *VtuberGoWorker) CleanTasks(ctx context.Context, input task.TaskInput) error {
	data, err := vgw.ExtractCleanTaskInput(input)
	if err != nil {
		panic(err)
	}
	_, err = vgw.db.QueueTask.Delete().
		Where(
			queuetask.And(
				queuetask.TaskName(data.TaskName),
				queuetask.CreatedAtLT(time.Now().Add(time.Duration(-data.After)*time.Millisecond)),
				queuetask.StatusNotIn(string(task.TaskStatusRunning)),
			),
		).
		Exec(ctx)
	if err != nil {
		panic(err)
	}
	return nil
}

func (vgw *VtuberGoWorker) GetVtubersWithActiveUsers(ctx context.Context) []*ent.Vtuber {
	vtubers, err := vgw.db.Vtuber.Query().
		Where(vtuber.HasUsers()).
		All(ctx)
	if err != nil {
		vgw.logger.Error("Get active vtubers error", zap.Error(err))
		return nil
	}
	return vtubers
}

func (vgw *VtuberGoWorker) NotifyUsers(ctx context.Context, _ task.TaskInput) error {
	timeNotifyAfter := 30
	vtubers := vgw.GetVtubersWithActiveUsers(ctx)
	vgw.logger.Debug("Got vtubers with active users", zap.Int("VtubersCount", len(vtubers)))
	videos := vgw.holodexController.GetVtubersUpcomingVideos(ctx, vtubers)
	vgw.logger.Debug("Got vtubers upcoming videos", zap.Int("VideosCount", len(videos)))
	now := time.Now()
	for videoIndex := range videos {
		vgw.logger.Debug("Process video", zap.String("VideoID", *videos[videoIndex].Id))

		if videos[videoIndex].AvailableAt == nil {
			vgw.logger.Warn("Available time nil")
			continue
		}
		if !now.Add(time.Duration(timeNotifyAfter) * time.Minute).After(*videos[videoIndex].AvailableAt) {
			vgw.logger.Debug("Video is not soon to start", zap.Time("AvailableAt", *videos[videoIndex].AvailableAt))
			continue
		}
		if videos[videoIndex].Channel.Id == nil {
			vgw.logger.Warn("Nil channel id")
			continue
		}
		exist, err := vgw.db.ReportedStream.Query().
			Where(reportedstream.VideoIDEQ(videos[videoIndex].GetId())).
			Exist(ctx)
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
			zap.String("VideoID", *videos[videoIndex].Id),
			zap.String("ChannelId", *videos[videoIndex].Channel.Id),
		)
		vgw.logger.Debug("Query users for vtuber")
		// We don't check the len of users, because before that we've queried only vtubers with active users
		users, err := vgw.db.Vtuber.Query().
			Where(vtuber.YoutubeChannelID(*videos[videoIndex].Channel.Id)).
			QueryUsers().
			All(ctx)
		if err != nil && !ent.IsNotFound(err) {
			vgw.logger.Error("Can't query users for vtuber", zap.Error(err))
			return err
		}
		if len(users) == 0 {
			vgw.logger.Warn("No users for vtuber video. Probably bug of holodex api. See TRE-54")
			continue
		}
		vgw.logger.Debug("Notify users on video in telegram", zap.Int("UsersCount", len(users)))
		vgw.telegramController.NotifyUsersOnVideo(users, videos[videoIndex])
		vgw.logger.Debug("Query authorId")
		authorID, err := vgw.db.Vtuber.Query().
			Where(vtuber.YoutubeChannelID(*videos[videoIndex].Channel.Id)).
			FirstID(ctx)
		if err != nil {
			vgw.logger.Error("Author search error", zap.Error(err))
			return err
		}
		vgw.logger.Debug("Set ReportedStream")
		err = vgw.db.ReportedStream.Create().
			SetAuthorID(authorID).
			SetAvailableAt(*videos[videoIndex].AvailableAt).
			SetVideoID(*videos[videoIndex].Id).
			SetVtuberID(authorID).
			Exec(ctx)
		if err != nil {
			vgw.logger.Error("Reported stream set error", zap.Error(err))
		}
		vgw.logger.Debug("Process video. Done", zap.String("VideoID", *videos[videoIndex].Id))
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

	entClient := setup.MustDatabaseSetupNoMigrations(config.DriverName, config.SQLURL)

	worker := task.NewWorker("MainWorker", entClient, zap.L())

	holodexController := holodexcontroller.NewHolodexController(config.HolodexAPIKey, zap.L())
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
	task.RegisterTask(task.TaskDescriptor{
		Name: "clean_tasks",
		Func: vgw.CleanTasks,
	})

	worker.Run(context.Background())
}
