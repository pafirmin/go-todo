package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pafirmin/go-todo/pkg/models"
	"github.com/pafirmin/go-todo/pkg/models/postgres"
)

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dto := &postgres.CreateTaskDTO{}
	err = json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := app.validator.Struct(dto); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	f, err := app.folders.GetByID(id)
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

	t, err := app.tasks.Insert(f.ID, dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, t)
}

func (app *application) getTasksByFolder(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	var input struct {
		Priority string
		models.Filters
	}

	qs := r.URL.Query()

	input.Priority = app.stringFromQuery(qs, "priority", "")
	input.Filters.Sort = app.stringFromQuery(qs, "sort", "id")
	input.Filters.Page = app.intFromQuery(qs, "page", 1)
	input.Filters.PageSize = app.intFromQuery(qs, "page_size", 20)
	input.Filters.SortSafeList = []string{"id", "due", "-id", "-due"}

	if !input.Filters.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	f, err := app.folders.GetByID(id)
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

	tasks, err := app.tasks.GetByFolder(f.ID, input.Priority, input.Filters)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, tasks)
}

func (app *application) getTaskByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	t, err := app.tasks.GetByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	f, err := app.folders.GetByID(t.FolderID)
	if err != nil {
		app.serverError(w, err)
		return
	} else if f.UserID != claims.UserID {
		app.forbidden(w)
		return
	}

	app.writeJSON(w, http.StatusOK, f)
}

func (app *application) updateTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok {
		app.unauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dto := &postgres.UpdateTaskDTO{}
	err = json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := app.validator.Struct(dto); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	t, err := app.tasks.GetByID(id)
	if err != nil {
		app.notFound(w)
		return
	}

	f, err := app.folders.GetByID(t.FolderID)
	if f.ID != claims.UserID {
		app.forbidden(w)
		return
	}

	if dto.FolderID != nil {
		f, err := app.folders.GetByID(*dto.FolderID)
		if err != nil {
			app.notFound(w)
			return
		}

		if f.UserID != claims.UserID {
			app.forbidden(w)
			return
		}
	}

	t, err = app.tasks.Update(id, dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, t)
}

func (app *application) removeTask(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok {
		app.unauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	t, err := app.tasks.GetByID(id)
	if err != nil {
		app.notFound(w)
		return
	}

	f, err := app.folders.GetByID(t.FolderID)
	if f.ID != claims.UserID {
		app.forbidden(w)
		return
	}

	_, err = app.tasks.Delete(t.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
