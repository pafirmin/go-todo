package mock

import (
	"time"

	"github.com/pafirmin/do-daily-go/pkg/models"
	"github.com/pafirmin/do-daily-go/pkg/models/postgres"
)

var mockUser = &models.User{
	ID:             1,
	Email:          "mock@example.com",
	HashedPassword: "1234",
	Created:        time.Now(),
}

type UserModel struct{}

func (m *UserModel) Insert(dto *postgres.CreateUserDTO) (*models.User, error) {
	return mockUser, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) Authenticate(cred *postgres.Credentials) (int, error) {
	switch cred.Email {
	case "mock@example.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}
