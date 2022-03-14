package postgres

import (
	"database/sql"
	"errors"
	"fmt"

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

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *UserModel) Get(id int) (*models.User, error) {
	stmt := `SELECT * FROM users WHERE users.id = $1`

	u := &models.User{}

	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

func (m *UserModel) Insert(dto *CreateUserDTO) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 12)
	if err != nil {
		return nil, err
	}

	stmt := `INSERT INTO users (email, hashed_password, created)
	VALUES($1, $2, now())
	RETURNING *`

	fmt.Println(dto.Email, string(hashedPassword))

	u := &models.User{}

	err = m.DB.
		QueryRow(stmt, dto.Email, string(hashedPassword)).
		Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Created)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (m *UserModel) Authenticate(creds *Credentials) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM users WHERE email = $1`

	row := m.DB.QueryRow(stmt, creds.Email)
	if err := row.Scan(&id, &hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(creds.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}
