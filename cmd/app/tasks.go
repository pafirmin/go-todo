package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pafirmin/go-todo/internal/data"
	"github.com/pafirmin/go-todo/internal/validator"
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
		app.notFound(w)
		return
	}

	f, err := app.models.Folders.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}

	if f.UserID != claims.UserID {
		app.forbidden(w)
		return
	}

	dto := &data.CreateTaskDTO{}
	err = app.readJSON(w, r, dto)
	if err != nil {
		app.badRequest(w, err.Error())
		return
	}

	v := validator.New()
	if v.Exec(dto); !v.Valid() {
		app.validationFailed(w, v)
		return
	}

	t, err := app.models.Tasks.Insert(f.ID, dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, responsePayload{"task": t})
}

func (app *application) getTasksByFolder(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.notFound(w)
		return
	}

	f, err := app.models.Folders.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}

	if f.UserID != claims.UserID {
		app.forbidden(w)
		return
	}

	var input struct {
		Status  string
		MinDate time.Time
		MaxDate time.Time
		data.Filters
	}

	qs := r.URL.Query()

	input.Status = app.stringFromQuery(qs, "status", "")
	input.MinDate = app.dateFromQuery(qs, "min_date", time.Time{})
	input.MaxDate = app.dateFromQuery(qs, "max_date", time.Time{})
	input.Filters.Sort = app.stringFromQuery(qs, "sort", "datetime")
	input.Filters.Page = app.intFromQuery(qs, "page", 1)
	input.Filters.PageSize = app.intFromQuery(qs, "page_size", 20)
	input.Filters.SortSafeList = []string{"id", "due", "created", "datetime", "-id", "-due", "-created", "-datetime"}

	v := validator.New()

	if v.Exec(&input.Filters); !v.Valid() {
		app.validationFailed(w, v)
		return
	}

	tasks, metadata, err := app.models.Tasks.GetByFolder(f.ID, input.Status, input.MinDate, input.MaxDate, input.Filters)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, responsePayload{"metadata": metadata, "tasks": tasks})
}

func (app *application) getTaskByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok {
		app.unauthorized(w)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.notFound(w)
		return
	}

	t, err := app.models.Tasks.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}

	f, err := app.models.Folders.GetByID(t.FolderID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if f.UserID != claims.UserID {
		app.forbidden(w)
		return
	}

	app.writeJSON(w, http.StatusOK, responsePayload{"task": t})
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
		app.notFound(w)
		return
	}

	t, err := app.models.Tasks.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}

	f, err := app.models.Folders.GetByID(t.FolderID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if f.UserID != claims.UserID {
		app.forbidden(w)
		return
	}

	dto := &data.UpdateTaskDTO{}
	err = app.readJSON(w, r, dto)
	if err != nil {
		app.badRequest(w, err.Error())
		return
	}

	v := validator.New()
	if v.Exec(dto); !v.Valid() {
		app.validationFailed(w, v)
		return
	}

	if dto.FolderID != nil && *dto.FolderID != f.ID {
		f, err := app.models.Folders.GetByID(*dto.FolderID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNoRecord):
				app.notFound(w)
			default:
				app.serverError(w, err)
			}
			return
		}

		if f.UserID != claims.UserID {
			app.forbidden(w)
			return
		}
	}

	t, err = app.models.Tasks.Update(id, dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, responsePayload{"task": t})
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
		app.notFound(w)
		return
	}

	t, err := app.models.Tasks.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}

	f, err := app.models.Folders.GetByID(t.FolderID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if f.ID != claims.UserID {
		app.forbidden(w)
		return
	}

	_, err = app.models.Tasks.Delete(t.ID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
