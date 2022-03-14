package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/pafirmin/do-daily-go/pkg/models/postgres"
)

type application struct {
	errorLog *log.Logger
	folders  *postgres.FolderModel
	infoLog  *log.Logger
	users    *postgres.UserModel
}

func main() {
	port := os.Getenv("PORT")
	dsn := "postgresql://postgres@localhost:5432/dodaily?sslmode=disable"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		folders:  &postgres.FolderModel{DB: db},
		infoLog:  infoLog,
		users:    &postgres.UserModel{DB: db},
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
