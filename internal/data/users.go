package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/pafirmin/go-todo/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Created        time.Time `json:"created"`
}

type CreateUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (d *CreateUserDTO) Validate(v *validator.Validator) {
	v.Check(d.Email != "", "email", "must be provided")
	v.Check(validator.IsEmail(d.Email), "email", "must be a valid email address")
	v.Check(len(d.Password) >= 8, "password", "must be at least 8 characters")
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m UserModel) Get(id int) (*User, error) {
	stmt := `SELECT * FROM users WHERE users.id = $1`

	u := &User{}

	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return u, nil
}

func (m UserModel) Insert(dto *CreateUserDTO) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 12)
	if err != nil {
		return nil, err
	}

	stmt := `INSERT INTO users (email, hashed_password, created)
	VALUES($1, $2, now())
	RETURNING *`

	u := User{}

	err = m.DB.
		QueryRow(stmt, dto.Email, string(hashedPassword)).
		Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Created)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m UserModel) Authenticate(creds *Credentials) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := `SELECT id, hashed_password FROM users WHERE email = $1`

	row := m.DB.QueryRow(stmt, creds.Email)
	if err := row.Scan(&id, &hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(creds.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}
