package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Models struct {
	Users interface {
		Insert(*CreateUserDTO) (*User, error)
		Get(int) (*User, error)
		GetByEmail(string) (*User, error)
		Authenticate(*Credentials) (*User, error)
		GetByToken(string, string) (*User, error)
	}
	Folders interface {
		Insert(int, *CreateFolderDTO) (*Folder, error)
		GetByID(int) (*Folder, error)
		GetByUser(int, Filters) ([]*Folder, MetaData, error)
		Update(int, *UpdateFolderDTO) (*Folder, error)
		Delete(int) (int, error)
	}
	Tasks interface {
		Insert(int, *CreateTaskDTO) (*Task, error)
		GetByUser(int, []string, string, time.Time, time.Time, Filters) ([]*Task, MetaData, error)
		GetByFolder(int, string, time.Time, time.Time, Filters) ([]*Task, MetaData, error)
		GetByID(int) (*Task, error)
		Update(int, *UpdateTaskDTO) (*Task, error)
		Delete(int) (int, error)
	}
	Tokens interface {
		New(int, time.Time, string) (*Token, error)
		Insert(*Token) error
		DeleteForUser(string, int) error
		Delete(string) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:   UserModel{DB: db},
		Folders: FolderModel{DB: db},
		Tasks:   TaskModel{DB: db},
		Tokens:  TokenModel{DB: db},
	}
}
