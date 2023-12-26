package controllers

import (
	"context"
	"net/http"

	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/org"
	"github.com/DanArmor/vtuber-go/ent/vtuber"
	"github.com/DanArmor/vtuber-go/ent/wave"
	"github.com/DanArmor/vtuber-go/pkg/controllers/resp"
	"github.com/DanArmor/vtuber-go/pkg/controllers/types"
	"github.com/gin-gonic/gin"
)

func (s *Service) PostVtubers(c *gin.Context) {
	type PostVtubersInput struct {
		Vtubers []types.Vtuber `json:"vtubers" binding:"required"`
	}
	var input PostVtubersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeCantBindJsonBody, "Can't bind json body"))
		return
	}
	vtubers := make([]ent.Vtuber, 0, len(input.Vtubers))
	for i := range input.Vtubers {
		vtuberExist, err := s.Db.Vtuber.Query().
			Where(vtuber.EnglishNameEQ(input.Vtubers[i].EnglishName)).
			Where(vtuber.HasWaveWith(
				wave.And(
					wave.NameEQ(input.Vtubers[i].WaveName), wave.HasOrgWith(org.NameEQ(input.Vtubers[i].CompanyName)),
				),
			)).
			Exist(c.Request.Context())
		if err != nil && !ent.IsNotFound(err) {
			c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeDbError, err.Error()))
			return
		}
		if vtuberExist {
			continue
		}
		var company *ent.Org
		company, err = s.Db.Org.Query().Where(org.NameEQ(input.Vtubers[i].CompanyName)).First(c.Request.Context())
		if err != nil {
			if !ent.IsNotFound(err) {
				c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeDbError, err.Error()))
				return
			} else {
				company, err = s.Db.Org.Create().SetName(input.Vtubers[i].CompanyName).Save(context.Background())
				if err != nil {
					c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeDbError, err.Error()))
					return
				}
			}
		}
		var vtuberWave *ent.Wave
		vtuberWave, err = s.Db.Org.QueryWaves(company).Where(wave.NameEQ(input.Vtubers[i].WaveName)).First(c.Request.Context())
		if err != nil {
			if !ent.IsNotFound(err) {
				c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeDbError, err.Error()))
				return
			} else {
				vtuberWave, err = s.Db.Wave.Create().SetName(input.Vtubers[i].WaveName).SetOrg(company).Save(context.Background())
				if err != nil {
					c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeDbError, err.Error()))
					return
				}
			}
		}

		new_vtuber, err := s.Db.Vtuber.Create().
			SetYoutubeChannelID(input.Vtubers[i].YoutubeChannelId).
			SetChannelName(input.Vtubers[i].ChannelName).
			SetEnglishName(input.Vtubers[i].EnglishName).
			SetPhotoURL(input.Vtubers[i].PhotoURL).
			SetTwitter(input.Vtubers[i].Twitter).
			SetVideoCount(input.Vtubers[i].VideoCount).
			SetSubscriberCount(input.Vtubers[i].SubscriberCount).
			SetClipCount(input.Vtubers[i].ClipCount).
			SetTopTopics(input.Vtubers[i].TopTopics).
			SetInactive(input.Vtubers[i].Inactive).
			SetTwitch(input.Vtubers[i].Twitch).
			SetBannerURL(input.Vtubers[i].BannerURL).
			SetWave(vtuberWave).
			Save(context.Background())
		if err != nil {
			c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeDbError, err.Error()))
			return
		}
		vtubers = append(vtubers, *new_vtuber)
	}

	c.JSON(http.StatusOK, resp.HandlerResult(gin.H{"vtubers": vtubers}))
}
