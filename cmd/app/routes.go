package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter().PathPrefix("/api/v1/").Subrouter()
	standardMiddleware := alice.New(defaultHeaders, cors.Default().Handler, app.logRequest)
	authMiddleware := alice.New(app.requireAuth)

	router.HandleFunc("/auth/login", app.login).Methods("POST")

	router.Handle("/users", authMiddleware.ThenFunc(app.createUser)).Methods("POST")
	router.Handle("/users/me", authMiddleware.ThenFunc(app.getUserByID)).Methods("GET")
	router.Handle("/users/me/folders", authMiddleware.ThenFunc(app.getFoldersByUser)).Methods("GET")
	router.Handle("/users/me/folders", authMiddleware.ThenFunc(app.createFolder)).Methods("POST")

	router.Handle("/folders/{id:[0-9]+}", authMiddleware.ThenFunc(app.getFolderByID)).Methods("GET")
	router.Handle("/folders/{id:[0-9]+}", authMiddleware.ThenFunc(app.updateFolder)).Methods("PATCH")
	router.Handle("/folders/{id:[0-9]+}", authMiddleware.ThenFunc(app.deleteFolder)).Methods("DELETE")
	router.Handle("/folders/{id:[0-9]+}/tasks", authMiddleware.ThenFunc(app.createTask)).Methods("POST")
	router.Handle("/folders/{id:[0-9]+}/tasks", authMiddleware.ThenFunc(app.getTasksByFolder)).Methods("GET")

	return standardMiddleware.Then(router)
}
