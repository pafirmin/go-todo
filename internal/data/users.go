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
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
}

type CreateUserDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (d *CreateUserDTO) Validate(v *validator.Validator) {
	v.ValidEmail("email", d.Email)
	v.ValidLength("password", d.Password, 8, 40)
	v.ValidLength("first_name", d.FirstName, 1, 40)
	v.ValidLength("last_name", d.LastName, 1, 40)
	v.ValidLength("password", d.Password, 8, 40)
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m UserModel) Get(id int) (*User, error) {
	stmt := `SELECT id, email, first_name, last_name, hashed_password, created, updated FROM users WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := &User{}

	rows := m.DB.QueryRowContext(ctx, stmt, id)

	err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.HashedPassword, &u.Created, &u.Updated)
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

	stmt := `INSERT INTO users (email, first_name, last_name, hashed_password, created, updated)
	VALUES($1, $2, $3, $4, DEFAULT, DEFAULT)
	RETURNING *`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := User{}
	args := []interface{}{dto.Email, dto.FirstName, dto.LastName, string(hashedPassword)}

	rows := m.DB.QueryRowContext(ctx, stmt, args...)

	err = rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.HashedPassword, &u.Created, &u.Updated)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (m UserModel) Authenticate(creds *Credentials) (*User, error) {
	stmt := `SELECT id, email, first_name, last_name, hashed_password, created, updated FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := User{}

	row := m.DB.QueryRowContext(ctx, stmt, creds.Email)
	if err := row.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.HashedPassword, &u.Created, &u.Updated); err != nil {
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
