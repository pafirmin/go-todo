package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/pafirmin/go-todo/internal/data"
	"github.com/pafirmin/go-todo/internal/jwt"
)

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
	infoLog    *log.Logger
	jwtService jwtService
	validator  *validator.Validate
	models     data.Models
}

func main() {
	var cfg config
	var secret string

	flag.IntVar(&cfg.port, "port", 4000, "Server port")
	flag.StringVar(&cfg.dbAddr, "db-address", "", "Postgres DB Address")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&secret, "jwt-secret", "", "JWT Secret key")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

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
		models:     data.NewModels(db),
		jwtService: jwt.NewJWTService([]byte(secret)),
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
