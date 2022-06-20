package data

import (
	"context"
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
	stmt := `SELECT id, email, hashed_password, created FROM users WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := &User{}

	rows := m.DB.QueryRowContext(ctx, stmt, id)

	err := rows.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Created)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := User{}

	rows := m.DB.QueryRowContext(ctx, stmt, dto.Email, string(hashedPassword))

	err = rows.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Created)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m UserModel) Authenticate(creds *Credentials) (*User, error) {
	stmt := `SELECT id, email, hashed_password, created FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := User{}

	row := m.DB.QueryRowContext(ctx, stmt, creds.Email)
	if err := row.Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		} else {
			return nil, err
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(creds.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		} else {
			return nil, err
		}
	}
	return &u, nil
}
