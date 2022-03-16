package mock

import (
	"time"

	"github.com/pafirmin/go-todo/pkg/models"
	"github.com/pafirmin/go-todo/pkg/models/postgres"
)

var mockTask = &models.Task{
	ID:          1,
	Title:       "Test",
	Description: "Test",
	Priority:    "low",
	FolderID:    1,
	Complete:    false,
	Created:     time.Now(),
	Due:         time.Now().Add(7 * 24 * time.Hour),
}

type TaskModel struct{}

func (t *TaskModel) Insert(dto *postgres.CreateTaskDTO) (*models.Task, error) {
	return mockTask, nil
}

func (t *TaskModel) GetByFolder(id int) ([]*models.Task, error) {
	return []*models.Task{mockTask}, nil
}

func (t *TaskModel) GetByID(id int) (*models.Task, error) {
	return mockTask, nil
}

func (t *TaskModel) Delete(id int) (int, error) {
	return 1, nil
}
