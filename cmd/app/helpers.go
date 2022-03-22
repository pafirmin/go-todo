package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/validator/v10"
	"github.com/pafirmin/go-todo/pkg/jwt"
)

var errNoUser = errors.New("no user in request context")

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	statusText := http.StatusText(http.StatusInternalServerError)
	http.Error(w, statusText, http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) validationError(w http.ResponseWriter, err validator.ValidationErrors) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) unauthorized(w http.ResponseWriter) {
	app.clientError(w, http.StatusUnauthorized)
}

func (app *application) forbidden(w http.ResponseWriter) {
	app.clientError(w, http.StatusForbidden)
}

func (app *application) claimsFromContext(ctx context.Context) (*jwt.UserClaims, bool) {
	claims, ok := ctx.Value(ctxKeyUserClaims).(*jwt.UserClaims)

	return claims, ok
}

func (app *application) writeJSON(w http.ResponseWriter, status int, body interface{}) {
	jsonRsp, err := json.Marshal(body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	w.Write(jsonRsp)
}
