package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pafirmin/go-todo/internal/validator"
)

type TaskModel struct {
	DB *sql.DB
}

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Datetime    time.Time `json:"datetime"`
	Status      string    `json:"status"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	FolderID    int       `json:"folder_id"`
}

type CreateTaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Datetime    string `json:"datetime"`
}

func (d *CreateTaskDTO) Validate(v *validator.Validator) {
	v.ValidLength("title", d.Title, 1, 50)
	v.ValidLength("description", d.Description, 0, 500)
	v.ValidDatetime("datetime", d.Datetime)
}

type UpdateTaskDTO struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Datetime    *string `json:"datetime,omitempty"`
	Status      *string `json:"status,omitempty"`
	FolderID    *int    `json:"folder_id,omitempty"`
}

func (d *UpdateTaskDTO) Validate(v *validator.Validator) {
	if d.Title != nil {
		v.ValidLength("title", *d.Title, 1, 50)
	}
	if d.Description != nil {
		v.ValidLength("description", *d.Description, 0, 500)
	}
	if d.Datetime != nil {
		v.ValidDatetime("datetime", *d.Datetime)
	}
	if d.Status != nil {
		v.PermittedValue("status", *d.Status, "default", "important", "cancelled")
	}
}

func (m TaskModel) Insert(folderID int, dto *CreateTaskDTO) (*Task, error) {
	stmt := `INSERT INTO tasks (title, description, status, datetime, created, updated, folder_id)
	VALUES ($1, $2, DEFAULT, $3, DEFAULT, DEFAULT, $4)
	RETURNING *`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := &Task{}
	args := []interface{}{dto.Title, dto.Description, dto.Datetime, folderID}

	row := m.DB.QueryRowContext(ctx, stmt, args...)
	err := row.Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Datetime,
		&t.Created,
		&t.Updated,
		&t.FolderID,
	)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (m TaskModel) GetByID(id int) (*Task, error) {
	stmt := `SELECT id, title, description, datetime, status, created, updated, folder_id
	FROM tasks 
	WHERE tasks.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := &Task{}

	row := m.DB.QueryRowContext(ctx, stmt, id)
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Datetime, &t.Status, &t.Created, &t.Updated, &t.FolderID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return t, nil
}

func (m TaskModel) GetByUser(
	userID int,
	folderIDs []string,
	status string,
	minDate,
	maxDate time.Time,
	filters Filters,
) ([]*Task, MetaData, error) {
	var minDateStmt string
	var maxDateStmt string

	if !minDate.IsZero() {
		minDateStmt = fmt.Sprintf("AND DATE_TRUNC('day', tasks.datetime) >= '%s'", minDate.Format("2006-01-02"))
	}
	if !maxDate.IsZero() {
		maxDateStmt = fmt.Sprintf("AND DATE_TRUNC('day', tasks.datetime) <= '%s'", maxDate.Format("2006-01-02"))
	}

	stmt := fmt.Sprintf(`SELECT count(*) OVER(),
	tasks.id, tasks.title, tasks.description, tasks.status, tasks.datetime, tasks.created, tasks.updated, tasks.folder_id
	FROM tasks
	INNER JOIN folders ON folders.id = tasks.folder_id 
	WHERE folders.user_id = $1
	AND (tasks.folder_id = ANY ($2::int[]) OR $2 = '{}')
	AND (tasks.status LIKE $3 OR $3 = '')
	%s
	%s
	ORDER BY tasks.%s %s, tasks.id ASC
	LIMIT $4 OFFSET $5`, minDateStmt, maxDateStmt, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{userID, pq.Array(folderIDs), status, filters.Limit(), filters.Offset()}

	rows, err := m.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, MetaData{}, err
	}

	defer rows.Close()

	totalRecords := 0
	tasks := []*Task{}

	for rows.Next() {
		t := &Task{}
		err = rows.Scan(
			&totalRecords,
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Datetime,
			&t.Created,
			&t.Updated,
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

func (m TaskModel) GetByFolder(
	folderID int,
	status string,
	minDate,
	maxDate time.Time,
	filters Filters,
) ([]*Task, MetaData, error) {
	var minDateStmt string
	var maxDateStmt string

	if !minDate.IsZero() {
		minDateStmt = fmt.Sprintf("AND DATE_TRUNC('day', tasks.datetime) >= '%s'", minDate.Format("2006-01-02"))
	}
	if !maxDate.IsZero() {
		maxDateStmt = fmt.Sprintf("AND DATE_TRUNC('day', tasks.datetime) <= '%s'", maxDate.Format("2006-01-02"))
	}

	stmt := fmt.Sprintf(`SELECT count(*) OVER(), * FROM tasks
		WHERE tasks.folder_id = $1
		AND (tasks.status LIKE $2 OR $2 = '')
		%s
		%s
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, minDateStmt, maxDateStmt, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{folderID, status, filters.Limit(), filters.Offset()}

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
			&t.Status,
			&t.Datetime,
			&t.Created,
			&t.Updated,
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
	SET title = COALESCE($1, title),
		description = COALESCE($2, description),
		status = COALESCE($3, status),
		datetime = COALESCE($4, datetime),
		folder_id = COALESCE($5, folder_id),
		updated = now()
	WHERE tasks.id = $6
	RETURNING *
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := &Task{}
	args := []interface{}{dto.Title, dto.Description, dto.Status, dto.Datetime, dto.FolderID, id}

	row := m.DB.QueryRowContext(ctx, stmt, args...)
	err := row.Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Datetime,
		&t.Created,
		&t.Updated,
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
