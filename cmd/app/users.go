package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/pafirmin/go-todo/internal/data"
	"github.com/pafirmin/go-todo/internal/validator"
)

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	creds := &data.Credentials{}

	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		app.unauthorized(w)
		return
	}

	id, err := app.models.Users.Authenticate(creds)
	if err != nil {
		app.unauthorized(w)
		return
	}

	exp := time.Now().Add(24 * time.Hour)
	token, err := app.jwtService.Sign(id, exp)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, responsePayload{"token": token})
}

func (app *application) getUserByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := app.claimsFromContext(r.Context())
	if !ok || claims.UserID < 1 {
		app.unauthorized(w)
		return
	}

	u, err := app.models.Users.Get(claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFound(w)
		default:
			app.serverError(w, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, responsePayload{"user": u})
}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	dto := &data.CreateUserDTO{}

	err := json.NewDecoder(r.Body).Decode(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	v := validator.New()
	if v.Exec(dto); !v.Valid() {
		app.validationError(w, v)
		return
	}

	u, err := app.models.Users.Insert(dto)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, responsePayload{"user": u})
}
