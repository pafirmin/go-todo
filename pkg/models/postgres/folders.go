package postgres

import (
	"database/sql"
	"errors"

	"github.com/pafirmin/do-daily-go/pkg/models"
)

type FolderModel struct {
	DB *sql.DB
}

type CreateFolderDTO struct {
	Name   string `json:"name"`
	UserID int    `json:"-"`
}

func (m *FolderModel) Get(id int) (*models.Folder, error) {
	stmt := `SELECT id, name, user_id, created
	FROM folders
	WHERE folders.id = $1`

	f := &models.Folder{}

	err := m.DB.QueryRow(stmt, id).Scan(&f.ID, &f.Name, &f.UserID, &f.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return f, nil
}

func (m *FolderModel) Insert(dto *CreateFolderDTO) (*models.Folder, error) {
	stmt := `INSERT INTO folders (name, user_id, created)
	VALUES($1, $2, DEFAULT)
	RETURNING *`

	f := &models.Folder{}

	err := m.DB.QueryRow(stmt, dto.Name, dto.UserID).Scan(&f.ID, &f.Name, &f.UserID, &f.Created)

	if err != nil {
		return nil, err
	}

	return f, nil
}
