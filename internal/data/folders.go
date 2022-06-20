package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pafirmin/go-todo/internal/validator"
)

type FolderModel struct {
	DB *sql.DB
}

type Folder struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	UserID  int       `json:"user_id"`
	Created time.Time `json:"created"`
}

type CreateFolderDTO struct {
	Name string `json:"name" validate:"required,min=1,max=30"`
}

func (d *CreateFolderDTO) Validate(v *validator.Validator) {
	v.Check(d.Name != "", "name", "folder name must be provided")
	v.Check(len(d.Name) < 40, "name", "must be shorter than 40 characters")
}

type UpdateFolderDTO struct {
	Name *string `json:"name" validate:"required,min=1,max=30,omitempty"`
}

func (d *UpdateFolderDTO) Validate(v *validator.Validator) {
	if d.Name != nil {
		v.Check(*d.Name == "", "name", "folder name must be provided")
		v.Check(len(*d.Name) < 40, "name", "must be shorter than 40 characters")
	}
}

func (m FolderModel) Insert(userId int, dto *CreateFolderDTO) (*Folder, error) {
	stmt := `INSERT INTO folders (name, user_id, created)
	VALUES($1, $2, DEFAULT)
	RETURNING *`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f := &Folder{}
	args := []interface{}{dto.Name, userId}

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&f.ID, &f.Name, &f.UserID, &f.Created)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (m FolderModel) GetByID(id int) (*Folder, error) {
	stmt := `SELECT id, name, user_id, created
	FROM folders
	WHERE folders.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f := &Folder{}

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&f.ID, &f.Name, &f.UserID, &f.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return f, nil
}

func (m FolderModel) GetByUser(userId int, filters Filters) ([]*Folder, MetaData, error) {
	stmt := fmt.Sprintf(`SELECT count(*) OVER(), id, name, user_id, created
	FROM folders
	WHERE folders.user_id = $1
	ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3
	`, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{userId, filters.Limit(), filters.Offset()}

	rows, err := m.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, MetaData{}, err
	}

	defer rows.Close()

	folders := []*Folder{}
	totalRecords := 0

	for rows.Next() {
		f := Folder{}
		err := rows.Scan(&totalRecords, &f.ID, &f.Name, &f.UserID, &f.Created)
		if err != nil {
			return nil, MetaData{}, err
		}
		folders = append(folders, &f)
	}

	metadata := CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return folders, metadata, nil
}

func (m FolderModel) Update(id int, dto *UpdateFolderDTO) (*Folder, error) {
	stmt := `UPDATE folders
	SET name = COALESCE($1, name)
	WHERE folders.id = $2
	RETURNING *
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f := &Folder{}

	err := m.DB.QueryRowContext(ctx, stmt, dto.Name, id).Scan(&f.Name)

	if err != nil {
		return nil, err
	}

	return f, err
}

func (m FolderModel) Delete(id int) (int, error) {
	stmt := `DELETE FROM folders WHERE folders.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
