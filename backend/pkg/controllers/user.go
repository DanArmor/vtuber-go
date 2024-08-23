package controllers

import (
	"context"
	"net/http"

	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/pkg/controllers/resp"
	"github.com/gin-gonic/gin"
)

func (s *Service) UserGetTimezone(c *gin.Context) {
	payload := getTokenPayload(c)
	userId := payload.UserId
	user, err := s.Db.User.Get(context.Background(), userId)
	if err != nil && !ent.IsNotFound(err) {
		c.JSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeDbError, "Internal error"))
		return
	}

	c.JSON(http.StatusOK, resp.HandlerResult(gin.H{"timezone_shift": user.TimezoneShift}))
}

func (s *Service) UserChangeTimezone(c *gin.Context) {
	type UserChangeTimezoneInput struct {
		TimezoneShift *int `json:"timezone_shift" binding:"required"`
	}
	var input UserChangeTimezoneInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeCantBindJsonBody, "Can't bind json body"))
		return
	}
	if *input.TimezoneShift < -14 || *input.TimezoneShift > 14 {
		c.JSON(http.StatusBadRequest, resp.HandlerError(resp.ErrCodeDbError, "Wrong timezone"))
		return
	}
	payload := getTokenPayload(c)
	userId := payload.UserId
	err := s.Db.User.UpdateOneID(userId).SetTimezoneShift(*input.TimezoneShift).Exec(c.Request.Context())
	if err != nil && !ent.IsNotFound(err) {
		c.JSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeDbError, "Internal error"))
		return
	}

	c.JSON(http.StatusOK, resp.HandlerResult(gin.H{"timezone_shift": *input.TimezoneShift}))
}
