package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1/").Subrouter()
	standardMiddleware := alice.New(app.recoverPanic, defaultHeaders, cors.Default().Handler, app.logRequest, app.rateLimit)
	authMiddleware := alice.New(app.requireAuth)

	s.HandleFunc("/status", app.showStatus).Methods(http.MethodGet)

	// Auth handlers
	s.HandleFunc("/auth/login", app.login).Methods(http.MethodPost)

	// User handlers
	s.Handle("/users", authMiddleware.ThenFunc(app.createUser)).Methods(http.MethodPost)
	s.Handle("/users/me", authMiddleware.ThenFunc(app.getUserByID)).Methods(http.MethodGet)

	// Folder handlers
	s.Handle("/users/me/folders", authMiddleware.ThenFunc(app.createFolder)).Methods(http.MethodPost)
	s.Handle("/users/me/folders", authMiddleware.ThenFunc(app.getFoldersByUser)).Methods(http.MethodGet)
	s.Handle("/folders/{id:[0-9]+}", authMiddleware.ThenFunc(app.getFolderByID)).Methods(http.MethodGet)
	s.Handle("/folders/{id:[0-9]+}", authMiddleware.ThenFunc(app.updateFolder)).Methods(http.MethodPatch)
	s.Handle("/folders/{id:[0-9]+}", authMiddleware.ThenFunc(app.removeFolder)).Methods(http.MethodDelete)

	// Task handlers
	s.Handle("/folders/{id:[0-9]+}/tasks", authMiddleware.ThenFunc(app.createTask)).Methods(http.MethodPost)
	s.Handle("/folders/{id:[0-9]+}/tasks", authMiddleware.ThenFunc(app.getTasksByFolder)).Methods(http.MethodGet)
	s.Handle("/tasks/{id:[0-9]+}", authMiddleware.ThenFunc(app.getTaskByID)).Methods(http.MethodGet)
	s.Handle("/tasks/{id:[0-9]+}", authMiddleware.ThenFunc(app.updateTask)).Methods(http.MethodPatch)
	s.Handle("/tasks/{id:[0-9]+}", authMiddleware.ThenFunc(app.removeTask)).Methods(http.MethodDelete)

	return standardMiddleware.Then(r)
}
