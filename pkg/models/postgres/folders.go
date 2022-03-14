package postgres

import (
	"database/sql"

	"github.com/pafirmin/do-daily-go/pkg/models"
)

type FolderModel struct {
	DB *sql.DB
}

func (m *FolderModel) Insert(name string, userID int) (*models.Folder, error) {
	stmt := `INERT INTO folders (name, user_id, created)
	VALUES($1, $2, UTC_TIMESTAMP())
	RETURNING *`

	f := &models.Folder{}

	err := m.DB.QueryRow(stmt, name, userID).Scan(&f.ID, &f.Name, &f.UserID, &f.Created)

	if err != nil {
		return nil, err
	}

	return f, nil
}
