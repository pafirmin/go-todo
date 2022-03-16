package mock

import (
	"time"

	"github.com/pafirmin/do-daily-go/pkg/models"
	"github.com/pafirmin/do-daily-go/pkg/models/postgres"
)

var mockFolder = &models.Folder{
	ID:      1,
	Name:    "Test",
	UserID:  1,
	Created: time.Now(),
}

type FolderModel struct{}

func (f *FolderModel) Insert(userId int, dto *postgres.CreateFolderDTO) (*models.Folder, error) {
	return mockFolder, nil
}

func (f *FolderModel) Get(id int) (*models.Folder, error) {
	switch id {
	case 1:
		return mockFolder, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (f *FolderModel) GetByUser(id int) ([]*models.Folder, error) {
	return []*models.Folder{mockFolder}, nil
}
