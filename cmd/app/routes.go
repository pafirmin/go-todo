package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	dynamicMiddleware := alice.New(app.requireAuth)
	router.HandleFunc("/auth/login", app.login).Methods("POST")

	router.Handle("/users", dynamicMiddleware.ThenFunc(app.createUser)).Methods("POST")
	router.Handle("/users/{id:[0-9]+}", dynamicMiddleware.ThenFunc(app.getUser)).Methods("GET")

	router.Handle("/folders", dynamicMiddleware.ThenFunc(app.getFolders)).Methods("GET")
	router.Handle("/folders/{id:[0-9]+}", dynamicMiddleware.ThenFunc(app.getOneFolder)).Methods("GET")
	router.Handle("/folders", dynamicMiddleware.ThenFunc(app.createFolder)).Methods("POST")
	router.Handle("/folders/{id:[0-9]+}", dynamicMiddleware.ThenFunc(app.updateFolder)).Methods("PATCH")
	router.Handle("/folders/{id:[0-9]+}", dynamicMiddleware.ThenFunc(app.deleteFolder)).Methods("DELETE")

	return router
}
