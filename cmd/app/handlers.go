package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pafirmin/do-daily-go/pkg/models"
	"github.com/pafirmin/do-daily-go/pkg/models/postgres"
)

func (app *application) getTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get all tasks"))
}

func (app *application) getOneTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get one task"))
}

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a task"))
}

func (app *application) updateTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update a task"))
}

func (app *application) deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete a task"))
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if id < 1 {
		app.notFound(w)
		return
	}

	u, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	jsonRsp, err := json.Marshal(u)
	if err != nil {
		app.serverError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRsp)
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	dto := &postgres.CreateUserDTO{}

	err := json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	u, err := app.users.Insert(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	jsonRsp, err := json.Marshal(u)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRsp)
}
