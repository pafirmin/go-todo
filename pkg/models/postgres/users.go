package postgres

import (
	"database/sql"

	"github.com/pafirmin/do-daily-go/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

type CreateUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *UserModel) Insert(dto *CreateUserDTO) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 12)
	if err != nil {
		return nil, err
	}

	stmt := `INSERT INTO users (email, hashed_password, created)
	VALUES($1, $2, now())
	RETURNING *`

	u := &models.User{}

	err = m.DB.
		QueryRow(stmt, dto.Email, hashedPassword).
		Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Created)

	if err != nil {
		return nil, err
	}

	return u, nil
}
