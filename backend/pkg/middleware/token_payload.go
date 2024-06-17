package middleware

import (
	"errors"

	"github.com/DanArmor/vtuber-go/pkg/auth"
	"github.com/gin-gonic/gin"
)

func GetTokenPayload(c *gin.Context) (auth.Payload, error) {
	val, exists := c.Get("token-payload")
	if !exists {
		return auth.Payload{}, errors.New("no payload in context")
	}
	return val.(auth.Payload), nil
}
