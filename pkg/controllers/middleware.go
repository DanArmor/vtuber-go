package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) CheckToken(c *gin.Context) {
	token := c.Request.Header.Get("vtubergo-token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "no token in header of the request")
	}
	payload, err := s.TokenMaker.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "can't verify token")
	}
	c.Set("token-payload", *payload)
	c.Next()
}
