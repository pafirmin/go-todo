package main

import (
	"database/sql"
	"flag"
	"log"
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

type config struct {
	port    int
	dbAddr  string
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config     config
	errorLog   *log.Logger
	folders    foldersService
	infoLog    *log.Logger
	jwtService jwtService
	tasks      tasksService
	users      usersService
	validator  *validator.Validate
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "Server port")
	flag.StringVar(&cfg.dbAddr, "db-address", "", "Postgres DB Address")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	secret := os.Getenv("SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Print(secret)

	db, err := openDB(cfg.dbAddr)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	infoLog.Print("database connection pool established")

	app := &application{
		config:     cfg,
		errorLog:   errorLog,
		infoLog:    infoLog,
		folders:    &postgres.FolderModel{DB: db},
		tasks:      &postgres.TaskModel{DB: db},
		users:      &postgres.UserModel{DB: db},
		jwtService: jwt.NewJWTService(secret),
		validator:  validator.New(),
	}

	err = app.serve()
	if err != nil {
		errorLog.Fatal(err, nil)
	}
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
