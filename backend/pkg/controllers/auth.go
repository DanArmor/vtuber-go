package controllers

import (
	"net/http"
	"time"

	"github.com/DanArmor/vtuber-go/ent"
	"github.com/DanArmor/vtuber-go/ent/user"
	telegramcontroller "github.com/DanArmor/vtuber-go/pkg/api/telegram_controller"
	"github.com/DanArmor/vtuber-go/pkg/controllers/resp"
	"github.com/DanArmor/vtuber-go/pkg/controllers/types"
	"github.com/gin-gonic/gin"
)

func (s *Service) AuthUser(c *gin.Context) {
	var input types.InitData
	if err := c.BindQuery(&input); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			resp.HandlerError(resp.ErrCodeCantBindJsonBody, "Can't bind query string"),
		)
		return
	}
	if input.User.ID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.HandlerError(resp.ErrCodeNoTgID, "No Tg Id"))
		return
	}
	values := c.Request.URL.Query()
	if err := telegramcontroller.CheckIntegrityInitData(values, s.TgBotToken, s.ExpirationHours); err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			resp.HandlerError(resp.ErrCodeCantValidateInitData, "Can't validate init data"),
		)
		return
	}

	id, err := s.DB.User.Query().Where(user.TgID(input.User.ID)).FirstID(c.Request.Context())
	if err != nil && !ent.IsNotFound(err) {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			resp.HandlerError(resp.ErrCodeCantValidateInitData, "Internal error"),
		)
		return
	}
	if id == 0 {
		createdUser, err := s.DB.User.Create().
			SetTgID(input.User.ID).
			SetFirstName(input.User.FirstName).
			SetLastName(input.User.LastName).
			SetUsername(input.User.Username).
			SetLanguageCode(input.User.LanguageCode).
			SetTimezoneShift(0).
			SetPhotoURL(input.User.PhotoURL).Save(c.Request.Context())
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				resp.HandlerError(resp.ErrCodeCantValidateInitData, "Internal error"),
			)
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
