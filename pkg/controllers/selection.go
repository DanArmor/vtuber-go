package controllers

import (
	"context"
	"net/http"

	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/user"
	"github.com/DanArmor/vtuber-go/ent/vtuber"
	"github.com/DanArmor/vtuber-go/pkg/controllers/resp"
	"github.com/gin-gonic/gin"
)

func (s *Service) SelectVtuber(c *gin.Context) {
	type SearchVtubersInput struct {
		VtuberId int `json:"vtuber_id" binding:"required"`
	}
	var input SearchVtubersInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeCantBindJsonBody, "Can't bind json body"))
		return
	}
	payload := getTokenPayload(c)
	userId := payload.UserId

	exists, err := s.Db.User.Query().
		Where(user.And(
			user.IDEQ(userId), user.HasVtubersWith(vtuber.IDEQ(input.VtuberId)),
		)).
		Exist(c.Request.Context())
	if err != nil && !ent.IsNotFound(err) {
		c.JSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeDbError, "Internal error"))
		return
	}
	var selected bool
	if exists {
		err := s.Db.User.UpdateOneID(userId).RemoveVtuberIDs(input.VtuberId).Exec(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeDbError, "Internal error"))
			return
		}
		selected = false
	} else {
		err := s.Db.User.UpdateOneID(userId).AddVtuberIDs(input.VtuberId).Exec(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeDbError, "Internal error"))
			return
		}
		selected = true
	}

	c.JSON(http.StatusOK, resp.HandlerResult(gin.H{"vtuber_id": input.VtuberId, "selected": selected}))
}
