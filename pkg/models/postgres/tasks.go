package postgres

import (
	"database/sql"
	"errors"

	"github.com/pafirmin/go-todo/pkg/models"
)

type TaskModel struct {
	DB *sql.DB
}

type CreateTaskDTO struct {
	Title       string `json:"title" validate:"required,min=1,max=30"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Due         string `json:"due" validate:"required"`
}

type UpdateTaskDTO struct {
	Title       *string `json:"title" validate:"required,min=1,max=30,omitempty"`
	Description *string `json:"description"`
	Priority    *string `json:"priority,omitempty"`
	Due         *string `json:"due" validate:"required,omitempty"`
	FolderID    *int    `json:"folderId" validate:"required,omitempty"`
}

func (m *TaskModel) Insert(folderID int, dto *CreateTaskDTO) (*models.Task, error) {
	stmt := `INSERT INTO tasks (title, description, priority, due, complete, created, folder_id)
	VALUES ($1, $2, $3, $4, DEFAULT, DEFAULT, $5)
	RETURNING *`

	t := &models.Task{}

	row := m.DB.QueryRow(stmt, dto.Title, dto.Description, dto.Priority, dto.Due, folderID)
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

func (m *TaskModel) GetByID(id int) (*models.Task, error) {
	stmt := `SELECT id, title, description, priority, due, complete, created, folder_id
	FROM tasks 
	WHERE tasks.id = $1`

	t := &models.Task{}

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Priority, &t.Due, &t.Complete, &t.Created, &t.FolderID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return t, nil
}

func (m *TaskModel) GetByFolder(folderId int) ([]*models.Task, error) {
	stmt := `SELECT * FROM tasks WHERE tasks.folder_id = $1`

	rows, err := m.DB.Query(stmt, folderId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []*models.Task{}

	for rows.Next() {
		t := &models.Task{}
		err := rows.Scan(
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
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (m *TaskModel) Update(id int, dto *UpdateTaskDTO) (*models.Task, error) {
	stmt := `UPDATE tasks
	SET title=COALESCE($1, title), description=COALESCE($2, description), priority=COALESCE($3, priority), due=COALESCE($4, due), folderId=COALESCE($5, folder_id)
	WHERE tasks.id = $6
	RETURNING *
	`
	t := &models.Task{}

	row := m.DB.QueryRow(stmt, dto.Title, dto.Description, dto.Priority, dto.FolderID, dto.Due)
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

func (m *TaskModel) Delete(id int) (int, error) {
	stmt := `DELETE FROM tasks WHERE tasks.id = $1`
	_, err := m.DB.Exec(stmt, id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
