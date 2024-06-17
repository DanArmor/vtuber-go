package controllers

import (
	"net/http"

	"github.com/DanArmor/vtuber-go/pkg/controllers/resp"
	"github.com/gin-gonic/gin"
)

func (s *Service) CheckToken(c *gin.Context) {
	token := c.Request.Header.Get("vtubergo-token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.HandlerError(resp.ErrCodeCantValidateInitData, "no token in header of the request"))
		return
	}
	payload, err := s.TokenMaker.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, resp.HandlerError(resp.ErrCodeCantValidateInitData, "can't verify token"))
		return
	}
	c.Set("token-payload", *payload)
	c.Next()
}
