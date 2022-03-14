package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
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

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	type createUserDTO struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	dto := &createUserDTO{}

	err := json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.errorLog.Output(2, trace)

		http.Error(w, http.StatusText(500), 500)
		return
	}

	u, err := app.users.Insert(dto.Email, dto.Password)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.errorLog.Output(2, trace)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	u.HashedPassword = ""
	jsonRsp, err := json.Marshal(u)
	if err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		app.errorLog.Output(2, trace)

		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRsp)
}
