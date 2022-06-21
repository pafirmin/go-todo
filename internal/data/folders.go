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
	Updated time.Time `json:"updated"`
}

type CreateFolderDTO struct {
	Name string `json:"name"`
}

func (d *CreateFolderDTO) Validate(v *validator.Validator) {
	v.ValidLength("name", d.Name, 1, 30)
}

type UpdateFolderDTO struct {
	Name *string `json:"name"`
}

func (d *UpdateFolderDTO) Validate(v *validator.Validator) {
	if d.Name != nil {
		v.ValidLength("name", *d.Name, 1, 30)
	}
}

func (m FolderModel) Insert(userId int, dto *CreateFolderDTO) (*Folder, error) {
	stmt := `INSERT INTO folders (name, user_id, created, updated)
	VALUES($1, $2, DEFAULT, DEFAULT)
	RETURNING *`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f := &Folder{}
	args := []interface{}{dto.Name, userId}

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&f.ID, &f.Name, &f.Created, &f.Updated, &f.UserID)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (m FolderModel) GetByID(id int) (*Folder, error) {
	stmt := `SELECT id, name, created, updated, user_id
	FROM folders
	WHERE folders.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f := &Folder{}

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&f.ID, &f.Name, &f.Created, &f.Updated, &f.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return f, nil
}

func (m FolderModel) GetByUser(userId int, filters Filters) ([]*Folder, MetaData, error) {
	stmt := fmt.Sprintf(`SELECT count(*) OVER(), id, name, created, updated, user_id
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
		err := rows.Scan(&totalRecords, &f.ID, &f.Name, &f.Created, &f.Updated, &f.UserID)
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
	SET name = COALESCE($1, name), updated = now()
	WHERE folders.id = $2
	RETURNING *
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	f := &Folder{}

	err := m.DB.QueryRowContext(ctx, stmt, dto.Name, id).Scan(&f.ID, &f.Name, &f.Created, &f.Updated, &f.UserID)

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
