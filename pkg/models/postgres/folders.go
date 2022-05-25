package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pafirmin/go-todo/pkg/models"
)

type FolderModel struct {
	DB *sql.DB
}

type CreateFolderDTO struct {
	Name string `json:"name" validate:"required,min=1,max=30"`
}

type UpdateFolderDTO struct {
	Name *string `json:"name" validate:"required,min=1,max=30,omitempty"`
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

func (m *FolderModel) GetByUser(userId int, filters models.Filters) ([]*models.Folder, models.MetaData, error) {
	stmt := fmt.Sprintf(`SELECT count(*) OVER(), id, name, user_id, created
	FROM folders
	WHERE folders.user_id = $1
	ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3
	`, filters.SortColumn(), filters.SortDirection())

	args := []interface{}{userId, filters.Limit(), filters.Offset()}

	rows, err := m.DB.Query(stmt, args...)
	if err != nil {
		return nil, models.MetaData{}, err
	}

	defer rows.Close()

	folders := []*models.Folder{}
	totalRecords := 0

	for rows.Next() {
		f := &models.Folder{}
		err := rows.Scan(&totalRecords, &f.ID, &f.Name, &f.UserID, &f.Created)
		if err != nil {
			return nil, models.MetaData{}, err
		}
		folders = append(folders, f)
	}

	metadata := models.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return folders, metadata, nil
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

func (m *FolderModel) Update(id int, dto *UpdateFolderDTO) (*models.Folder, error) {
	stmt := `UPDATE folders
	SET name = COALESCE($1, name)
	WHERE folders.id = $2
	RETURNING *
	`

	f := &models.Folder{}

	err := m.DB.QueryRow(stmt, dto.Name, id).Scan(&f.Name)

	if err != nil {
		return nil, err
	}

	return f, err
}

func (m *FolderModel) Delete(id int) (int, error) {
	stmt := `DELETE FROM folders WHERE folders.id = $1`
	_, err := m.DB.Exec(stmt, id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
