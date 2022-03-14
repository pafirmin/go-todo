package postgres

import (
	"database/sql"

	"github.com/pafirmin/do-daily-go/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(email string, password string) (*models.User, error) {
	stmt := `INSERT INTO users (email, hashed_password, created)
	VALUES($1, $2, now())
	RETURNING *`

	u := &models.User{}

	err := m.DB.QueryRow(stmt, email, password).Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Created)

	if err != nil {
		return nil, err
	}

	return u, nil
}
