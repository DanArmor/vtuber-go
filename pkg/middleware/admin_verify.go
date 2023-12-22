package middleware

import (
	"net/http"

	"github.com/DanArmor/vtuber-go/pkg/controllers/resp"
	"github.com/gin-gonic/gin"
)

func AdminVerify(adminToken string) func(c *gin.Context) {
	return func(c *gin.Context) {

		val := c.Request.Header.Get("ADMIN-TOKEN")
		if val != adminToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.HandlerError(resp.ErrCodeWrongAdminToken, "Wrong admin token"))
			return
		}

		c.Next()
	}
}
