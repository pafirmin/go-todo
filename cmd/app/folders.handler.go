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

func (app *application) createFolder(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	dto := &postgres.CreateFolderDTO{}
	err := json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := app.validator.Struct(dto); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	f, err := app.folders.Insert(claims.UserID, dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, responseWrapper{"folder": f})
}

func (app *application) getFoldersByUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok || claims.UserID < 0 {
		app.unauthorized(w)
		return
	}

	var input struct {
		models.Filters
	}

	qs := r.URL.Query()

	input.Filters.Page = app.intFromQuery(qs, "page", 1)
	input.Filters.PageSize = app.intFromQuery(qs, "page_size", 20)
	input.Filters.Sort = app.stringFromQuery(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "name", "-id", "-name"}

	if !input.Filters.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	folders, metadata, err := app.folders.GetByUser(claims.UserID, input.Filters)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, responseWrapper{"metadata": metadata, "folders": folders})
}

func (app *application) getFolderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	claims, ok := app.claimsFromContext(r.Context())
	if !ok || claims.UserID < 1 {
		app.unauthorized(w)
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

	app.writeJSON(w, http.StatusOK, responseWrapper{"folder": f})
}

func (app *application) updateFolder(w http.ResponseWriter, r *http.Request) {
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

	dto := &postgres.UpdateFolderDTO{}
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
	if f.ID != claims.UserID {
		app.forbidden(w)
		return
	}

	f, err = app.folders.Update(id, dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, responseWrapper{"folder": f})
}

func (app *application) removeFolder(w http.ResponseWriter, r *http.Request) {
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

	_, err = app.folders.Delete(id)

	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
