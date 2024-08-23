package controllers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	_ "time/tzdata"

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
			if videos[i].Channel.Id == nil {
				log.Println("Notify nil Channel.Id")
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
			users, err := s.Db.Vtuber.Query().Where(vtuber.YoutubeChannelID(*videos[i].Channel.Id)).QueryUsers().All(context.Background())
			if err != nil {
				log.Printf("Notify error: %v", err)
				return
			}
			for j := range users {
				loc := time.FixedZone("temp-zone", users[j].TimezoneShift*60*60)
				userChat := telebot.ChatID(users[j].TgID)
				channelName := videos[i].Channel.EnglishName.Get()
				if !videos[i].Channel.EnglishName.IsSet() {
					channelName = videos[i].Channel.Name
				}
				msg := telebot.Photo{}
				msg.File = telebot.FromURL(fmt.Sprintf("https://img.youtube.com/vi/%s/0.jpg", *videos[i].Id))
				msg.Caption = fmt.Sprintf("Stream of %s wiil begin in <b>~%d</b> minutes\n\nTitle: %s\n\n<a href='https://www.youtube.com/watch?v=%s'>▶️ Stream link</a>\nStart time: <b>%02d:%02d</b> (GMT %+d)",
					*channelName,
					int(videos[i].AvailableAt.Sub(now).Minutes()),
					*videos[i].Title,
					*videos[i].Id,
					videos[i].AvailableAt.In(loc).Hour(),
					videos[i].AvailableAt.In(loc).Minute(),
					users[j].TimezoneShift,
				)
				_, err := s.TgBot.Send(userChat, &msg, telebot.ModeHTML)
				if err != nil {
					log.Printf("Can't send notify: %v")
				}
			}
			authorId, err := s.Db.Vtuber.Query().Where(vtuber.YoutubeChannelID(*videos[i].Channel.Id)).FirstID(context.Background())
			if err != nil {
				log.Println(err)
			}
			err = s.Db.ReportedStream.Create().
				SetAuthorID(authorId).
				SetAvailableAt(*videos[i].AvailableAt).
				SetVideoID(*videos[i].Id).
				SetVtuberID(authorId).
				Exec(context.Background())
			if err != nil {
				log.Println(err)
			}

		}
	}
}

func (s *Service) SchedulerGo() {
	interval := time.Duration(time.Minute * time.Duration(s.TimeStep))
	counter := 1
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
		if counter > 10 {
			panic("scheduler restart limit exceeded")
		}
	}
}
