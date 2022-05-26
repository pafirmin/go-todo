package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pafirmin/go-todo/internal/validator"
)

type TaskModel struct {
	DB *sql.DB
}

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Due         time.Time `json:"due"`
	Complete    bool      `json:"complete"`
	Created     time.Time `json:"created"`
	FolderID    int       `json:"folderId"`
}

type CreateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Due         string `json:"due"`
}

func (d *CreateTaskDTO) Validate(v *validator.Validator) {
	v.Check(d.Title != "", "title", "must be provided")
	v.Check(len(d.Title) < 40, "title", "must be shorter than 40 characters")
	v.Check(len(d.Description) < 600, "description", "must be shorter than 600 characters")
	v.Check(validator.ValidDate(d.Due), "due", "must be valid RFC3339 date string")
	v.Check(validator.PermittedValue(d.Priority, "low", "medium", "high"), "priority", "must be one of 'low', 'medium' or 'high'")
}

type UpdateTaskDTO struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Priority    *string `json:"priority,omitempty"`
	Due         *string `json:"due,omitempty"`
	FolderID    *int    `json:"folderId,omitempty"`
}

func (d *UpdateTaskDTO) Validate(v *validator.Validator) {
	if d.Title != nil {
		v.Check(*d.Title != "", "title", "must be provided")
		v.Check(len(*d.Title) < 40, "title", "must be shorter than 40 characters")
	}
	if d.Description != nil {
		v.Check(len(*d.Description) < 600, "description", "must be shorter than 600 characters")
	}
}

func (m TaskModel) Insert(folderID int, dto *CreateTaskDTO) (*Task, error) {
	stmt := `INSERT INTO tasks (title, description, priority, due, complete, created, folder_id)
	VALUES ($1, $2, $3, $4, DEFAULT, DEFAULT, $5)
	RETURNING *`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := &Task{}

	row := m.DB.QueryRowContext(ctx, stmt, dto.Title, dto.Description, dto.Priority, dto.Due, folderID)
	err := row.Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Priority,
		&t.Due,
		&t.Complete,
		&t.Created,
		&t.FolderID,
	)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (m TaskModel) GetByID(id int) (*Task, error) {
	stmt := `SELECT id, title, description, priority, due, complete, created, folder_id
	FROM tasks 
	WHERE tasks.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := &Task{}

	row := m.DB.QueryRowContext(ctx, stmt, id)
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Priority, &t.Due, &t.Complete, &t.Created, &t.FolderID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return t, nil
}

func (m TaskModel) GetByFolder(folderId int, priority string, filters Filters) ([]*Task, MetaData, error) {
	stmt := fmt.Sprintf(`SELECT count(*) OVER(), * FROM tasks
		WHERE tasks.folder_id = $1
		AND tasks.priority LIKE $2 OR $2 = ''
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{folderId, priority, filters.Limit(), filters.Offset()}

	rows, err := m.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, MetaData{}, err
	}

	defer rows.Close()

	totalRecords := 0
	tasks := []*Task{}

	for rows.Next() {
		t := &Task{}
		err := rows.Scan(
			&totalRecords,
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Priority,
			&t.Due,
			&t.Complete,
			&t.Created,
			&t.FolderID,
		)
		if err != nil {
			return nil, MetaData{}, err
		}
		tasks = append(tasks, t)
	}

	metadata := CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return tasks, metadata, nil
}

func (m TaskModel) Update(id int, dto *UpdateTaskDTO) (*Task, error) {
	stmt := `UPDATE tasks
	SET title=COALESCE($1, title), description=COALESCE($2, description), priority=COALESCE($3, priority), due=COALESCE($4, due), folder_id=COALESCE($5, folder_id)
	WHERE tasks.id = $6
	RETURNING *
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := &Task{}

	row := m.DB.QueryRowContext(ctx, stmt, dto.Title, dto.Description, dto.Priority, dto.Due, dto.FolderID, id)
	err := row.Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Priority,
		&t.Due,
		&t.Complete,
		&t.Created,
		&t.FolderID,
	)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (m TaskModel) Delete(id int) (int, error) {
	stmt := `DELETE FROM tasks WHERE tasks.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, stmt, id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
