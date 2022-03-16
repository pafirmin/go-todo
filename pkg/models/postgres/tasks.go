package postgres

import (
	"database/sql"
	"time"

	"github.com/pafirmin/do-daily-go/pkg/models"
)

type TaskModel struct {
	DB *sql.DB
}

type CreateTaskDTO struct {
	Title       string    `json:"title" validate:"required,min=1,max=30"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Due         time.Time `json:"due" validate:"required,datetime"`
	FolderID    int       `json:"folderId"`
}

func (m *TaskModel) Insert(dto *CreateTaskDTO) (*models.Task, error) {
	stmt := `INSERT INTO tasks (title, description, priority, due, complete, created, folder_id)
	VALUES ($1, $2, $3, $4, DEFAULT, DEFAULT, $5)
	RETURNING *`

	t := &models.Task{}

	row := m.DB.QueryRow(stmt, dto.Title, dto.Description, dto.Priority, dto.Due, dto.FolderID)
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
