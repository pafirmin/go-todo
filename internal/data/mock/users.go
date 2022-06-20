package mock

import (
	"time"

	"github.com/pafirmin/go-todo/internal/data"
)

var mockUser = &data.User{
	ID:             1,
	Email:          "mock@example.com",
	HashedPassword: "1234",
	Created:        time.Now(),
}

type UserModel struct{}

func (m UserModel) Insert(dto *data.CreateUserDTO) (*data.User, error) {
	return mockUser, nil
}

func (m UserModel) Get(id int) (*data.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, data.ErrNoRecord
	}
}

func (m UserModel) Authenticate(cred *data.Credentials) (*data.User, error) {
	switch cred.Email {
	case "mock@example.com":
		return mockUser, nil
	default:
		return nil, data.ErrInvalidCredentials
	}
}
