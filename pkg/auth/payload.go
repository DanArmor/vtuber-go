package auth

import (
	"time"

	"github.com/DanArmor/vtuber-go/pkg/controllers/types"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Payload contains the payload data of the token
type Payload struct {
	jwt.StandardClaims
	Id       uuid.UUID      `json:"id"`
	InitData types.InitData `json:"init_data"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(initData types.InitData, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		Id:       tokenId,
		InitData: initData,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Local().Add(time.Hour * duration).Unix(),
		},
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (p *Payload) Valid() error {
	if time.Now().Unix() > p.ExpiresAt {
		return ErrExpiredToken
	}
	return nil
}
