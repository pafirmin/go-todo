package mock

import (
	"time"

	"github.com/pafirmin/go-todo/internal/data"
)

var mockTask = &data.Task{
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

func (t TaskModel) Insert(id int, dto *data.CreateTaskDTO) (*data.Task, error) {
	return mockTask, nil
}

func (t TaskModel) GetByFolder(id int, priority string, filters data.Filters) ([]*data.Task, data.MetaData, error) {
	return []*data.Task{mockTask}, data.MetaData{}, nil
}

func (t TaskModel) GetByID(id int) (*data.Task, error) {
	switch id {
	case 1:
		return mockTask, nil
	default:
		return nil, data.ErrNoRecord
	}
}

func (t TaskModel) Update(id int, dto *data.UpdateTaskDTO) (*data.Task, error) {
	return mockTask, nil
}

func (t TaskModel) Delete(id int) (int, error) {
	return 1, nil
}
