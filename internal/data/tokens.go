package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

const (
	ScopeRefresh = "refresh"
)

type TokenModel struct {
	DB *sql.DB
}

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int, exp time.Time, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: exp,
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func (m TokenModel) New(userID int, exp time.Time, scope string) (*Token, error) {
	token, err := generateToken(userID, exp, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)

	return token, err
}

func (m TokenModel) Insert(token *Token) error {
	stmt := `INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, args...)
	return err
}

func (m TokenModel) DeleteForUser(scope string, userID int) error {
	stmt := `DELETE FROM tokens WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, scope, userID)
	return err
}

func (m TokenModel) Delete(tokenText string) error {
	hash := sha256.Sum256([]byte(tokenText))
	stmt := `DELETE FROM tokens WHERE hash = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, hash[:])
	return err
}
