package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Service struct {
	Secret []byte
}

type UserClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var (
	ErrInvalidSigningMethod = errors.New("jwt: invalid signing method")
)

func NewService(secret []byte) *Service {
	return &Service{
		Secret: secret,
	}
}

func (j *Service) Sign(id int, expires time.Time) (string, error) {
	claims := &UserClaims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ret, err := token.SignedString(j.Secret)
	if err != nil {
		return "", err
	}

	return ret, nil
}

func (j *Service) Parse(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(j.Secret), nil
	})

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
