package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	Email  string `json:"email"`
	UserID int    `json:"userId"`
	jwt.StandardClaims
}

var (
	secret                  = []byte(os.Getenv("SECRET"))
	ErrInvalidSigningMethod = errors.New("jwt: invalid signing method")
	ErrNoUser               = errors.New("jwt: could not parse user data from token")
)

func Sign(id int, email string, expires time.Time) (string, error) {
	claims := &UserClaims{
		Email:  email,
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ret, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return ret, nil
}

func Parse(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
