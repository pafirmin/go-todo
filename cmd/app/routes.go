package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/users", app.createUser).Methods("POST")

	router.HandleFunc("/tasks", app.getTasks).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", app.getOneTask).Methods("GET")
	router.HandleFunc("/tasks", app.createTask).Methods("POST")
	router.HandleFunc("/tasks/{id:[0-9]+}", app.updateTask).Methods("PATCH")
	router.HandleFunc("/tasks/{id:[0-9]+}", app.deleteTask).Methods("DELETE")

	return router
}
