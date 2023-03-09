package mock

import (
	"time"

	"github.com/pafirmin/go-todo/internal/data"
)

type TokenModel struct{}

func (m TokenModel) New(userID int, exp time.Time, scope string) (*data.Token, error) {
	return &data.Token{}, nil
}

func (m TokenModel) Insert(token *data.Token) error {
	return nil
}

func (m TokenModel) DeleteForUser(scope string, userID int) error {
	return nil
}

func (m TokenModel) Delete(tokenText string) error {
	return nil
}
