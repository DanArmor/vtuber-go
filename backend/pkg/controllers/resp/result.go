package resp

import "github.com/gin-gonic/gin"

// HandlerResult is a wrapper for results of handlers
func HandlerResult(content gin.H) gin.H {
	return gin.H{"result": content}
}
