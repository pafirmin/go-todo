package main

import (
	"errors"
	"net/http"

	"github.com/pafirmin/go-todo/internal/data"
	"github.com/pafirmin/go-todo/internal/validator"
)

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

	err := app.readJSON(w, r, dto)
	if err != nil {
		app.badRequest(w, err.Error())
		return
	}

	v := validator.New()
	if v.Exec(dto); !v.Valid() {
		app.validationFailed(w, v)
		return
	}

	u, err := app.models.Users.Insert(dto)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "email already in use")
			app.validationFailed(w, v)
		default:
			app.serverError(w, err)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, responsePayload{"user": u})
}
