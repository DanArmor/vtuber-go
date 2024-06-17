package controllers

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/DanArmor/vtuber-go/ent/reportedstream"
	"github.com/DanArmor/vtuber-go/ent/vtuber"
	"gopkg.in/telebot.v3"
)

func (s *Service) NotifyUsers() {
	vtubers, err := s.Db.Vtuber.Query().
		Where(vtuber.HasUsers()).
		All(context.Background())
	if err != nil {
		log.Printf("Notify error: %v", err)
		return
	}

	ids := make([]string, 0, len(vtubers))
	for i := range vtubers {
		if vtubers[i].YoutubeChannelID != "" {
			ids = append(ids, vtubers[i].YoutubeChannelID)
		}
	}
	queryParam := strings.Join(ids, ",")

	request := s.HolodexClient.DefaultApi.GetCachedLive(context.Background()).Channels(queryParam)

	videos, response, err := s.HolodexClient.DefaultApi.GetCachedLiveExecute(request)
	defer response.Body.Close()
	if err != nil || response.StatusCode != 200 {
		log.Printf("Notify error: %v / %d", err, response.StatusCode)
		return
	}
	now := time.Now()
	for i := range videos {
		if videos[i].AvailableAt == nil {
			log.Printf("Notify nil available time")
			continue
		}
		if now.Add(time.Duration(s.TimeNotifyAfter) * time.Minute).After(*videos[i].AvailableAt) {
			if videos[i].ChannelId == nil {
				log.Printf("Notify nil Channel.ID")
				continue
			}
			exist, err := s.Db.ReportedStream.Query().Where(reportedstream.VideoIDEQ(videos[i].GetId())).Exist(context.Background())
			if err != nil {
				log.Printf("Notify error: %v", err)
				return
			}
			if exist {
				log.Printf("Stream %s exists", videos[i].GetId())
				continue
			}
			users, err := s.Db.Vtuber.Query().Where(vtuber.YoutubeChannelID(*videos[i].ChannelId)).QueryUsers().All(context.Background())
			if err != nil {
				log.Printf("Notify error: %v", err)
				return
			}
			for j := range users {
				c := telebot.ChatID(users[j].TgID)
				s.TgBot.Send(c, "Test message for channel "+*videos[i].ChannelId+" and video "+videos[i].GetId()+" and title: "+videos[i].GetTitle())
			}
		}
	}
}

func (s *Service) SchedulerGo() {
	interval := time.Duration(time.Minute * time.Duration(s.TimeStep))
	quit := make(chan int, 1)
	for {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Recovered: %v", err)
					quit <- 1
				}
			}()
			ticker := time.NewTicker(interval)

			for {
				select {
				case <-ticker.C:
					s.NotifyUsers()
				}
			}
		}()
		<-quit
	}
}
