package controllers

import (
	"github.com/DanArmor/vtuber-go/ent"
	holodex "github.com/watsonindustries/go-holodex"
)

type Service struct {
	Db              *ent.Client
	TgBotToken      string
	ExpirationHours int
	HolodexClient   *holodex.APIClient
}

func NewService(db *ent.Client, tgBotToken string, expirationHours int, holodexClient *holodex.APIClient) Service {
	return Service{
		Db:              db,
		TgBotToken:      tgBotToken,
		ExpirationHours: expirationHours,
		HolodexClient:   holodexClient,
	}
}
