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
		if !ent.IsNotFound(err) {
			c.JSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeDbError, "Internal error"))
			return
		}
	}
	c.JSON(http.StatusOK, resp.HandlerResult(gin.H{"orgs": orgs}))
}

func (s *Service) SearchVtubers(c *gin.Context) {
	const MinPageSize = 10
	const MaxPageSize = 30
	type SearchVtubersInput struct {
		Name     string            `json:"name"`
		Org      []int             `json:"orgs"`
		Wave     []int             `json:"waves"`
		Selected selected.Selected `json:"selected"`
		Offset   *int              `json:"offset" binding:"required"`
		Limit    int               `json:"page_size" binding:"required"`
	}
	var input SearchVtubersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeCantBindJsonBody, "Can't bind json body"))
		return
	}

	if input.Limit < MinPageSize || input.Limit > MaxPageSize {
		c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeCantBindJsonBody, "Page size is incorrect"))
		return
	}

	payload := getTokenPayload(c)
	userId := payload.UserId

	query := s.Db.Vtuber.Query()
	if input.Name != "" {
		query.Where(vtuber.EnglishNameContainsFold(input.Name))
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
				user.IDEQ(userId),
			),
		)
	} else if input.Selected == selected.No {
		query.Where(
			vtuber.Not(
				vtuber.HasUsersWith(
					user.IDEQ(userId),
				),
			),
		)
	}
	vtubers, err := query.
		WithWave(func(wq *ent.WaveQuery) {
			wq.WithOrg()
		}).
		WithUsers(func(uq *ent.UserQuery) {
			uq.Where(user.IDEQ(userId)).IDs(c.Request.Context())
		}).
		Limit(input.Limit).
		Offset(*input.Offset).
		All(c.Request.Context())
	if err != nil {
		if !ent.IsNotFound(err) {
			c.JSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeDbError, "Internal error"))
			return
		}
	}
	c.JSON(http.StatusOK, resp.HandlerResult(gin.H{
		"vtubers": vtubers,
		"page_meta": gin.H{
			"offset":         input.Offset,
			"page_size_req":  input.Limit,
			"page_size_resp": len(vtubers),
		},
	}))
}
