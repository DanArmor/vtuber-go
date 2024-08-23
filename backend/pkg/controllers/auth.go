package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/user"
	"github.com/DanArmor/vtuber-go/pkg/controllers/resp"
	"github.com/DanArmor/vtuber-go/pkg/controllers/types"
	"github.com/DanArmor/vtuber-go/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (s *Service) AuthUser(c *gin.Context) {
	var input types.InitData
	if err := c.BindQuery(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.HandlerError(resp.ErrCodeCantBindJsonBody, "Can't bind query string"))
		return
	}
	if input.User.Id == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.HandlerError(resp.ErrCodeNoTgId, "No Tg Id"))
		return
	}
	values := c.Request.URL.Query()
	if err := utils.CheckIntegrityInitData(values, s.TgBotToken, s.ExpirationHours); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.HandlerError(resp.ErrCodeCantValidateInitData, "Can't validate init data"))
		return
	}

	id, err := s.Db.User.Query().Where(user.TgID(input.User.Id)).FirstID(context.Background())
	if err != nil && !ent.IsNotFound(err) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeCantValidateInitData, "Internal error"))
		return
	}
	if id == 0 {
		createdUser, err := s.Db.User.Create().
			SetTgID(input.User.Id).
			SetFirstName(input.User.FirstName).
			SetLastName(input.User.LastName).
			SetUsername(input.User.Username).
			SetLanguageCode(input.User.LanguageCode).
			SetTimezoneShift(0).
			SetPhotoURL(input.User.PhotoUrl).Save(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeCantValidateInitData, "Internal error"))
			return
		}
		id = createdUser.ID
	}
	token, _, err := s.TokenMaker.CreateToken(
		input,
		id,
		time.Duration(s.ExpirationHours),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.HandlerError(resp.ErrCodeDbError, "Internal error"))
		return
	}
	c.JSON(http.StatusOK, resp.HandlerResult(gin.H{"token": token}))
}
