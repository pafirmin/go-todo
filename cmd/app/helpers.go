package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/pafirmin/do-daily-go/pkg/jwt"
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

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) ctxClaims(ctx context.Context) (*jwt.UserClaims, error) {
	claims, ok := ctx.Value(ctxKeyUserClaims).(*jwt.UserClaims)
	fmt.Println(claims)

	if ok {
		return claims, nil
	}

	return nil, errNoUser
}
