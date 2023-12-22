package controllers

import (
	"net/http"

	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/org"
	"github.com/DanArmor/vtuber-go/ent/user"
	"github.com/DanArmor/vtuber-go/ent/vtuber"
	"github.com/DanArmor/vtuber-go/ent/wave"
	"github.com/DanArmor/vtuber-go/pkg/controllers/resp"
	"github.com/DanArmor/vtuber-go/pkg/utils/selected"
	"github.com/gin-gonic/gin"
)

func (s *Service) GetOrgs(c *gin.Context) {
	orgs, err := s.Db.Org.Query().WithWaves().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeDbError, "Internal error"))
		return
	}
	c.JSON(http.StatusOK, resp.HandlerResult(gin.H{"orgs": orgs}))
}

func (s *Service) SearchVtubers(c *gin.Context) {
	type SearchVtubersInput struct {
		Name     string            `json:"name"`
		Org      []int             `json:"orgs"`
		Wave     []int             `json:"waves"`
		Selected selected.Selected `json:"selected"`
	}
	var input SearchVtubersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeCantBindJsonBody, "Can't bind json body"))
		return
	}

	initData := getInitData(c)

	query := s.Db.Vtuber.Query()
	if input.Name != "" {
		query.Where(vtuber.EnglishNameContains(input.Name))
	}
	if len(input.Org) != 0 {
		query.Where(
			vtuber.HasWaveWith(
				wave.HasOrgWith(
					org.IDIn(input.Org...),
				),
			),
		)
	}
	if len(input.Wave) != 0 {
		query.Where(
			vtuber.HasWaveWith(
				wave.IDIn(input.Wave...),
			),
		)
	}
	if input.Selected == selected.Yes {
		query.Where(
			vtuber.HasUsersWith(
				user.TgIDEQ(initData.User.Id),
			),
		)
	} else if input.Selected == selected.No {
		query.Where(
			vtuber.Not(
				vtuber.HasUsersWith(
					user.TgIDEQ(initData.User.Id),
				),
			),
		)
	}

	vtubers, err := query.All(c.Request.Context())
	if err != nil {
		if !ent.IsNotFound(err) {
			c.JSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeDbError, "Internal error"))
			return
		}
	}
	c.JSON(http.StatusOK, resp.HandlerResult(gin.H{"vtubers": vtubers}))
}
