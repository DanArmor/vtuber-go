package controllers

import (
	"github.com/DanArmor/vtuber-go/pkg/auth"
	"github.com/gin-gonic/gin"
)

func getTokenPayload(c *gin.Context) auth.Payload {
	v, _ := c.Get("token-payload")
	return v.(auth.Payload)
}
