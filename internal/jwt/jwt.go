package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	Secret []byte
}

type UserClaims struct {
	Email  string `json:"email"`
	UserID int    `json:"userId"`
	jwt.StandardClaims
}

var (
	ErrInvalidSigningMethod = errors.New("jwt: invalid signing method")
	ErrNoUser               = errors.New("jwt: could not parse user data from token")
)

func NewJWTService(secret []byte) *JWTService {
	return &JWTService{
		Secret: secret,
	}
}

func (j *JWTService) Sign(id int, email string, expires time.Time) (string, error) {
	claims := &UserClaims{
		Email:  email,
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

func (j *JWTService) Parse(tokenStr string) (*UserClaims, error) {
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
