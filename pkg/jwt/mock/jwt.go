package mock

import (
	"errors"
	"time"

	"github.com/pafirmin/do-daily-go/pkg/jwt"
)

type JWTService struct {
	Secret string
}

func (j *JWTService) Parse(tokenStr string) (*jwt.UserClaims, error) {
	if tokenStr == "invalid" {
		return nil, errors.New("Invalid token")
	}

	claims := jwt.UserClaims{
		UserID: 1,
		Email:  "test@example.com",
	}
	claims.ExpiresAt = time.Now().Add(24 * time.Hour).UnixMicro()
	claims.Email = "test@example.com"
	if tokenStr == "123" {
		claims.UserID = 1
	} else {
		claims.UserID = 2
	}

	return &claims, nil
}

func (j *JWTService) Sign(id int, email string, expires time.Time) (string, error) {
	return "123", nil
}
