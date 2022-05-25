package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pafirmin/go-todo/pkg/models"
	"github.com/pafirmin/go-todo/pkg/models/postgres"
)

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

	app.writeJSON(w, http.StatusOK, responseWrapper{"token": token})
}

func (app *application) getUserByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	u, err := app.users.Get(claims.UserID)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, responseWrapper{"user": u})
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	dto := &postgres.CreateUserDTO{}

	err := json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if err := app.validator.Struct(dto); err != nil {
		app.validationError(w, err.(validator.ValidationErrors))
		return
	}

	u, err := app.users.Insert(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, responseWrapper{"user": u})
}
