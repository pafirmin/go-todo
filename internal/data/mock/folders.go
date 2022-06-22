package mock

import (
	"time"

	"github.com/pafirmin/go-todo/internal/data"
)

var mockFolder = &data.Folder{
	ID:      1,
	Name:    "Test",
	UserID:  1,
	Created: time.Now(),
}

type FolderModel struct{}

func (f FolderModel) Insert(userID int, dto *data.CreateFolderDTO) (*data.Folder, error) {
	return mockFolder, nil
}

func (f FolderModel) GetByID(id int) (*data.Folder, error) {
	switch id {
	case 1:
		return mockFolder, nil
	default:
		return nil, data.ErrNoRecord
	}
}

func (f FolderModel) GetByUser(id int, filters data.Filters) ([]*data.Folder, data.MetaData, error) {
	return []*data.Folder{mockFolder}, data.MetaData{}, nil
}

func (f FolderModel) Update(id int, dto *data.UpdateFolderDTO) (*data.Folder, error) {
	return mockFolder, nil
}

func (f FolderModel) Delete(id int) (int, error) {
	return 1, nil
}
