package controllers

import (
	"github.com/DanArmor/vtuber-go/ent"
)

type Service struct {
	Db              *ent.Client
	TgBotToken      string
	ExpirationHours int
}

func NewService(db *ent.Client, tgBotToken string, expirationHours int) Service {
	return Service{
		Db:              db,
		TgBotToken:      tgBotToken,
		ExpirationHours: expirationHours,
	}
}
