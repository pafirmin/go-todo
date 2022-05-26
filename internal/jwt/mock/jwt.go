package mock

import (
	"errors"
	"time"

	"github.com/pafirmin/go-todo/internal/jwt"
)

type JWTService struct {
	Secret string
}

func (j *JWTService) Parse(tokenStr string) (*jwt.UserClaims, error) {
	if tokenStr == "invalid" {
		return nil, errors.New("invalid token")
	}

	claims := jwt.UserClaims{
		UserID: 1,
	}
	claims.ExpiresAt = time.Now().Add(24 * time.Hour).UnixMicro()
	if tokenStr == "123" {
		claims.UserID = 1
	} else {
		claims.UserID = 2
	}

	return &claims, nil
}

func (j *JWTService) Sign(id int, expires time.Time) (string, error) {
	return "123", nil
}
