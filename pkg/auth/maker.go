package auth

import (
	"time"

	"github.com/DanArmor/vtuber-go/pkg/controllers/types"
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(initData types.InitData, userId int, duration time.Duration) (string, *Payload, error)
	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
