package controllers

import (
	"log"
	"time"

	holodex "github.com/DanArmor/go-holodex"
	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/pkg/auth"
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
}

func NewService(db *ent.Client, tgBotToken string, expirationHours int, timeNotifyAfter int, timeStep int,
	holodexClient *holodex.APIClient, authMaker auth.Maker) Service {
	pref := telebot.Settings{
		Token:  tgBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	return Service{
		Db:              db,
		TgBotToken:      tgBotToken,
		ExpirationHours: expirationHours,
		HolodexClient:   holodexClient,
		TokenMaker:      authMaker,
		TgBot:           b,
		TimeNotifyAfter: timeNotifyAfter,
		TimeStep:        timeStep,
	}
}
