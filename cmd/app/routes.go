package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter().PathPrefix("/api/v1/").Subrouter()
	standardMiddleware := alice.New(defaultHeaders, cors.Default().Handler, app.logRequest, app.rateLimit)
	authMiddleware := alice.New(app.requireAuth)

	router.HandleFunc("/status", app.showStatus).Methods(http.MethodGet)

	// Auth handlers
	router.HandleFunc("/auth/login", app.login).Methods(http.MethodPost)

	// User handlers
	router.Handle("/users", authMiddleware.ThenFunc(app.createUser)).Methods(http.MethodPost)
	router.Handle("/users/me", authMiddleware.ThenFunc(app.getUserByID)).Methods(http.MethodGet)

	// Folder handlers
	router.Handle("/users/me/folders", authMiddleware.ThenFunc(app.createFolder)).Methods(http.MethodPost)
	router.Handle("/users/me/folders", authMiddleware.ThenFunc(app.getFoldersByUser)).Methods(http.MethodGet)
	router.Handle("/folders/{id:[0-9]+}", authMiddleware.ThenFunc(app.getFolderByID)).Methods(http.MethodGet)
	router.Handle("/folders/{id:[0-9]+}", authMiddleware.ThenFunc(app.updateFolder)).Methods(http.MethodPatch)
	router.Handle("/folders/{id:[0-9]+}", authMiddleware.ThenFunc(app.removeFolder)).Methods(http.MethodDelete)

	// Task handlers
	router.Handle("/folders/{id:[0-9]+}/tasks", authMiddleware.ThenFunc(app.createTask)).Methods(http.MethodPost)
	router.Handle("/folders/{id:[0-9]+}/tasks", authMiddleware.ThenFunc(app.getTasksByFolder)).Methods(http.MethodGet)
	router.Handle("/tasks/{id:[0-9]+}", authMiddleware.ThenFunc(app.getTaskByID)).Methods(http.MethodGet)
	router.Handle("/tasks/{id:[0-9]+}", authMiddleware.ThenFunc(app.updateTask)).Methods(http.MethodPatch)
	router.Handle("/tasks/{id:[0-9]+}", authMiddleware.ThenFunc(app.removeTask)).Methods(http.MethodDelete)

	return standardMiddleware.Then(router)
}
