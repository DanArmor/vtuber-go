package controllers

import (
	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/pkg/auth"
	holodex "github.com/watsonindustries/go-holodex"
)

type Service struct {
	Db              *ent.Client
	TgBotToken      string
	ExpirationHours int
	HolodexClient   *holodex.APIClient
	TokenMaker      auth.Maker
}

func NewService(db *ent.Client, tgBotToken string, expirationHours int, holodexClient *holodex.APIClient, authMaker auth.Maker) Service {
	return Service{
		Db:              db,
		TgBotToken:      tgBotToken,
		ExpirationHours: expirationHours,
		HolodexClient:   holodexClient,
		TokenMaker:      authMaker,
	}
}
