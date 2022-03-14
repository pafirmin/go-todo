package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID int    `json:"userId"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

var secret = []byte(os.Getenv("SECRET"))

func Sign(id int, email string, expires time.Time) (string, error) {
	claims := &Claims{
		UserID: id,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}
	fmt.Println(secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ret, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return ret, nil
}
