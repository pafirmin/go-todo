package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pafirmin/do-daily-go/pkg/models"
	"github.com/pafirmin/do-daily-go/pkg/models/postgres"
)

func (app *application) getFolders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get all tasks"))
}

func (app *application) getFolderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	claims, err := app.ctxClaims(r.Context())
	if err != nil || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	f, err := app.folders.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	} else if f.UserID != claims.UserID {
		app.forbidden(w)
		return
	}

	jsonRsp, err := json.Marshal(f)
	if err != nil {
		app.serverError(w, err)
	}

	w.Write(jsonRsp)
}

func (app *application) createFolder(w http.ResponseWriter, r *http.Request) {
	claims, err := app.ctxClaims(r.Context())
	if err != nil || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	dto := &postgres.CreateFolderDTO{UserID: claims.UserID}
	err = json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	f, err := app.folders.Insert(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	jsonRsp, err := json.Marshal(f)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRsp)
}

func (app *application) updateFolder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update a task"))
}

func (app *application) deleteFolder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete a task"))
}

func (app *application) getTasksByFolder(w http.ResponseWriter, r *http.Request) {
	claims, err := app.ctxClaims(r.Context())
	if err != nil || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	f, err := app.folders.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	} else if f.UserID != claims.UserID {
		app.forbidden(w)
		return
	}

	tasks, err := app.tasks.GetByFolder(f.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	jsonRsp, err := json.Marshal(tasks)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Write(jsonRsp)
}

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	claims, err := app.ctxClaims(r.Context())
	if err != nil || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	f, err := app.folders.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	} else if f.UserID != claims.UserID {
		app.forbidden(w)
		return
	}

	dto := &postgres.CreateTaskDTO{FolderID: id}
	err = json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	t, err := app.tasks.Insert(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}
	jsonRsp, err := json.Marshal(t)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRsp)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	creds := &postgres.Credentials{}

	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		app.unauthorized(w)
		return
	}

	id, err := app.users.Authenticate(creds)
	if err != nil {
		app.unauthorized(w)
		return
	}

	exp := time.Now().Add(24 * time.Hour)
	token, err := app.jwtService.Sign(id, creds.Email, exp)
	if err != nil {
		app.serverError(w, err)
		return
	}

	rsp := map[string]string{token: token}
	jsonRsp, err := json.Marshal(rsp)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Write([]byte(jsonRsp))
}

func (app *application) getUserByID(w http.ResponseWriter, r *http.Request) {
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
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	jsonRsp, err := json.Marshal(u)
	if err != nil {
		app.serverError(w, err)
	}

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

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonRsp)
}
