package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
		app.badRequest(w)
		return
	}

	dto := &data.CreateTaskDTO{}
	err = json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if v := validator.New(); !v.Validate(dto) {
		app.validationError(w, v)
		return
	}

	f, err := app.models.Folders.GetByID(id)
	if errors.Is(err, data.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	} else if f.UserID != claims.UserID {
		app.forbidden(w)
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

	var input struct {
		Priority string
		data.Filters
	}

	qs := r.URL.Query()

	input.Priority = app.stringFromQuery(qs, "priority", "")
	input.Filters.Sort = app.stringFromQuery(qs, "sort", "id")
	input.Filters.Page = app.intFromQuery(qs, "page", 1)
	input.Filters.PageSize = app.intFromQuery(qs, "page_size", 20)
	input.Filters.SortSafeList = []string{"id", "due", "created", "-id", "-due", "-created"}

	if !input.Filters.Valid() {
		app.badRequest(w)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.badRequest(w)
		return
	}

	f, err := app.models.Folders.GetByID(id)
	if errors.Is(err, data.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	} else if f.UserID != claims.UserID {
		app.forbidden(w)
		return
	}

	tasks, metadata, err := app.models.Tasks.GetByFolder(f.ID, input.Priority, input.Filters)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, responsePayload{"metadata": metadata, "tasks": tasks})
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
		app.badRequest(w)
		return
	}

	t, err := app.models.Tasks.GetByID(id)
	if err != nil {
		if errors.Is(err, data.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	f, err := app.models.Folders.GetByID(t.FolderID)
	if err != nil {
		app.serverError(w, err)
		return
	} else if f.UserID != claims.UserID {
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
		app.badRequest(w)
		return
	}

	dto := &data.UpdateTaskDTO{}
	err = json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if v := validator.New(); !v.Validate(dto) {
		app.validationError(w, v)
		return
	}

	t, err := app.models.Tasks.GetByID(id)
	if err != nil {
		app.notFound(w)
		return
	}

	f, err := app.models.Folders.GetByID(t.FolderID)
	if f.ID != claims.UserID {
		app.forbidden(w)
		return
	}

	if dto.FolderID != nil {
		f, err := app.models.Folders.GetByID(*dto.FolderID)
		if err != nil {
			app.notFound(w)
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
		app.badRequest(w)
		return
	}

	t, err := app.models.Tasks.GetByID(id)
	if err != nil {
		app.notFound(w)
		return
	}

	f, err := app.models.Folders.GetByID(t.FolderID)
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
