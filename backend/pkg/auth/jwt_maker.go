package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/DanArmor/vtuber-go/pkg/controllers/types"
	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

// JwtMaker is a JSON Web Token maker
type JwtMaker struct {
	secretKey string
}

// NewJwtMaker creates a new JwtMaker
func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JwtMaker{
		secretKey: secretKey,
	}, nil
}

// CreateToken creates a new token for a specific username and duration
func (m *JwtMaker) CreateToken(initData types.InitData, userId int, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(initData, userId, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(m.secretKey))
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (m *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
