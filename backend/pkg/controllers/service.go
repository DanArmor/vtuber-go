package controllers

import (
	"bytes"
	"io/fs"
	"time"

	holodex "github.com/DanArmor/go-holodex"
	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/internal/infrastructure/auth"
	"github.com/DanArmor/vtuber-go/internal/infrastructure/task"
	"github.com/DanArmor/vtuber-go/internal/resources/photos"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

type Service struct {
	Db              *ent.Client
	TgBotToken      string
	ExpirationHours int
	HolodexClient   *holodex.APIClient
	TokenMaker      auth.Maker
	TgBot           *telebot.Bot
	TimeNotifyAfter int
	TimeStep        int
	Scheduler       *task.TaskScheduler
}

func NewService(db *ent.Client, tgBotToken string, expirationHours int, timeNotifyAfter int, timeStep int,
	holodexClient *holodex.APIClient, authMaker auth.Maker) *Service {
	pref := telebot.Settings{
		Token:  tgBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		zap.L().Fatal("Can't create telegram bot", zap.Error(err))
	}

	var photoSlice []*telebot.Photo
	photoFiles, err := fs.ReadDir(photos.PhotosFS, photos.PhotosDir)
	if err != nil {
		panic(err)
	}
	for _, file := range photoFiles {
		if file.IsDir() {
			continue
		}
		photoSlice = append(photoSlice, &telebot.Photo{})
		data, _ := fs.ReadFile(photos.PhotosFS, photos.PhotosDir+"/"+file.Name())
		reader := bytes.NewReader(data)
		photoSlice[len(photoSlice)-1].File = telebot.FromReader(reader)
	}
	photoSlice[0].Caption = "Telegram bot that reminds users about streams of selected vtubers \\(30 minutes before the stream\\)\\.\nIt has user interface implemented via Telegram Mini App mechanism, you can access it by pressing 'Menu' button in the lower left corner\\.\n\nSelect the streamers whose broadcasts you want to be notified about\\. You can use filters by:\n\\- *name*;\n\\- *company*;\n\\- *wave*;\n\\- *all/selected/not selected property*\\.\nAlso in 'Show' tab you can hide tags of vtubers and avatars\\.\n\nIn 'Settings' tab you can select the shift from GMT of your local time to get proper 'start time' value\\. If you didn't select any \\- it will format start time as GMT \\+0\\.\n\nAlso you can click on streamer in list and got streamer card with links to their accounts on youtube/twitch/twitter and etc\\."
	album := telebot.Album{}
	for i := range photoSlice {
		album = append(album, photoSlice[i])
	}

	b.Handle("/start", func(c telebot.Context) error {
		err := c.SendAlbum(album, telebot.ModeMarkdownV2)
		return err
	})
	go func() {
		counter := 1
		quit := make(chan int, 1)
		for {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						zap.L().Warn("Recovered", zap.Any("error", err))
						quit <- 1
					}
				}()
				b.Start()
			}()
			<-quit
			counter += 1
			if counter > 10 {
				panic("telegram bot restart limit exceeded")
			}
		}
	}()

	taskScheduler := task.NewTaskScheduler(db, zap.L())

	interval := 120000
	taskScheduler.AddScheduleTask(task.AddScheduleTaskInput{
		ScheduleTaskName: "notify_telegram_users",
		TaskName:         "notify_telegram_users",
		Interval:         &interval,
	})

	s := &Service{
		Db:              db,
		TgBotToken:      tgBotToken,
		ExpirationHours: expirationHours,
		HolodexClient:   holodexClient,
		TokenMaker:      authMaker,
		TgBot:           b,
		TimeNotifyAfter: timeNotifyAfter,
		TimeStep:        timeStep,
		Scheduler:       taskScheduler,
	}

	return s
}
