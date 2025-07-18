package telegramcontroller

import (
	"fmt"
	"time"

	_ "time/tzdata"

	"github.com/DanArmor/go-holodex"
	"github.com/DanArmor/vtuber-go/ent"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

type TelegramController struct {
	tgBotToken string
	tgBot      *telebot.Bot
	logger     *zap.Logger
}

func NewTelegramController(tgBotToken string, logger *zap.Logger) *TelegramController {
	newLogger := logger.With(zap.String("Service", "TelegramController"))
	newLogger.Debug("Constructing TelegramController")
	pref := telebot.Settings{
		Token:  tgBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		newLogger.Fatal("Fatal error on NewBot call", zap.Error(err))
	}
	return &TelegramController{
		tgBotToken: tgBotToken,
		tgBot:      b,
		logger:     newLogger,
	}
}

func (tc *TelegramController) NotifyUsersOnVideo(users []*ent.User, video holodex.Video) {
	now := time.Now()

	for user_index := range users {
		loc := time.FixedZone("temp-zone", users[user_index].TimezoneShift*60*60)
		userChat := telebot.ChatID(users[user_index].TgID)
		channelName := video.Channel.EnglishName.Get()
		if !video.Channel.EnglishName.IsSet() {
			channelName = video.Channel.Name
		}
		msg := telebot.Photo{}
		msg.File = telebot.FromURL(fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", *video.Id))
		msg.Caption = fmt.Sprintf("Stream of %s wiil begin in <b>~%d</b> minutes\n\nTitle: %s\n\n<a href='https://www.youtube.com/watch?v=%s'>▶️ Stream link</a>\nStart time: <b>%02d:%02d</b> (GMT %+d)",
			*channelName,
			int(video.AvailableAt.Sub(now).Minutes()),
			*video.Title,
			*video.Id,
			video.AvailableAt.In(loc).Hour(),
			video.AvailableAt.In(loc).Minute(),
			users[user_index].TimezoneShift,
		)
		_, err := tc.tgBot.Send(userChat, &msg, telebot.ModeHTML)
		if err != nil {
			tc.logger.Error("Can't send notify: %v")
		}
	}
}
