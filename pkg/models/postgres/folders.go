package postgres

import (
	"database/sql"
	"errors"

	"github.com/pafirmin/go-todo/pkg/models"
)

type FolderModel struct {
	DB *sql.DB
}

type CreateFolderDTO struct {
	Name string `json:"name" validate:"required,min=1,max=30"`
}

func (m *FolderModel) GetByID(id int) (*models.Folder, error) {
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

func (m *FolderModel) GetByUser(userId int) ([]*models.Folder, error) {
	stmt := `SELECT id, name, user_id, created
	FROM folders
	WHERE folders.user_id = $1`

	rows, err := m.DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	folders := []*models.Folder{}

	for rows.Next() {
		f := &models.Folder{}
		err := rows.Scan(&f.ID, &f.Name, &f.UserID, &f.Created)
		if err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}

	return folders, nil
}

func (m *FolderModel) Insert(userId int, dto *CreateFolderDTO) (*models.Folder, error) {
	stmt := `INSERT INTO folders (name, user_id, created)
	VALUES($1, $2, DEFAULT)
	RETURNING *`

	f := &models.Folder{}

	err := m.DB.QueryRow(stmt, dto.Name, userId).Scan(&f.ID, &f.Name, &f.UserID, &f.Created)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (m *FolderModel) Delete(id int) (int, error) {
	stmt := `DELETE FROM folders WHERE folders.id = $1`
	_, err := m.DB.Exec(stmt, id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
