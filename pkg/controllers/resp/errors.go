package resp

import "github.com/gin-gonic/gin"

const (
	// 1xx - Specific auth errors
	ErrCodeUserAlreadyExists = 100
	ErrCodeNoSuchUser        = 101
	ErrCodeWrongPass         = 102
	// 2xx - Input errors
	ErrCodeCantBindJsonBody     = 200
	ErrCodeCantBindQueryBody    = 201
	ErrCodeCantValidateInitData = 202
	ErrCodeNoTgId               = 203
	// 5xx - Internal Errors
	ErrCodeDbError = 500
)

// HandlerError is a wrapper for errors in handlers
func HandlerError(code int, msg string) gin.H {
	return gin.H{"error": map[string]any{"code": code, "msg": msg}}
}
