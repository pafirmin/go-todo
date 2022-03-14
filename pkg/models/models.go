package models

import "time"

type User struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Created        time.Time `json:"created"`
}

type Folder struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	UserID  int       `json:"userId"`
	Created time.Time `json:"created"`
}

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"Description"`
	Priority    string    `json:"Priority"`
	Due         time.Time `json:"due"`
	Complete    bool      `json:"complete"`
	Created     time.Time `json:"created"`
	FolderID    int       `json:"folderId"`
}
