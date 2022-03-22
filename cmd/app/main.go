package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/pafirmin/go-todo/pkg/jwt"
	"github.com/pafirmin/go-todo/pkg/models"
	"github.com/pafirmin/go-todo/pkg/models/postgres"
)

type usersService interface {
	Insert(*postgres.CreateUserDTO) (*models.User, error)
	Get(int) (*models.User, error)
	Authenticate(*postgres.Credentials) (int, error)
}

type foldersService interface {
	Insert(int, *postgres.CreateFolderDTO) (*models.Folder, error)
	GetByID(int) (*models.Folder, error)
	GetByUser(int) ([]*models.Folder, error)
	Update(int, *postgres.UpdateFolderDTO) (*models.Folder, error)
	Delete(int) (int, error)
}

type tasksService interface {
	Insert(int, *postgres.CreateTaskDTO) (*models.Task, error)
	GetByFolder(int) ([]*models.Task, error)
	GetByID(int) (*models.Task, error)
	Update(int, *postgres.UpdateTaskDTO) (*models.Task, error)
	Delete(int) (int, error)
}

type jwtService interface {
	Sign(int, string, time.Time) (string, error)
	Parse(string) (*jwt.UserClaims, error)
}

type application struct {
	errorLog   *log.Logger
	folders    foldersService
	infoLog    *log.Logger
	jwtService jwtService
	tasks      tasksService
	users      usersService
	validator  *validator.Validate
}

func main() {
	port := os.Getenv("PORT")
	secret := os.Getenv("SECRET")
	dsn := "postgresql://postgres@localhost:5432/dodaily?sslmode=disable"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog:   errorLog,
		folders:    &postgres.FolderModel{DB: db},
		infoLog:    infoLog,
		jwtService: jwt.NewJWTService(secret),
		tasks:      &postgres.TaskModel{DB: db},
		users:      &postgres.UserModel{DB: db},
		validator:  validator.New(),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Staring server on %s", port)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
