package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/pafirmin/go-todo/internal/jwt"
	"github.com/pafirmin/go-todo/internal/validator"
)

type responsePayload map[string]interface{}

func (app *application) errorResponse(w http.ResponseWriter, status int, message interface{}) {
	app.writeJSON(w, status, responsePayload{"message": message})
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	app.errorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *application) validationError(w http.ResponseWriter, v *validator.Validator) {
	app.writeJSON(w, http.StatusUnprocessableEntity, responsePayload{"message": "validation failed", "errors": v.Errors})
}

func (app *application) badRequest(w http.ResponseWriter) {
	app.errorResponse(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
}

func (app *application) notFound(w http.ResponseWriter) {
	msg := "the requested resource could not be found"
	app.errorResponse(w, http.StatusNotFound, msg)
}

func (app *application) unauthorized(w http.ResponseWriter) {
	msg := "you are not authorised to access this resource"
	app.errorResponse(w, http.StatusUnauthorized, msg)
}

func (app *application) forbidden(w http.ResponseWriter) {
	msg := "you do not have permission to access this resource"
	app.errorResponse(w, http.StatusForbidden, msg)
}

func (app *application) rateLimitExceeded(w http.ResponseWriter) {
	app.errorResponse(w, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
}

func (app *application) claimsFromContext(ctx context.Context) (*jwt.UserClaims, bool) {
	claims, ok := ctx.Value(ctxKeyUserClaims).(*jwt.UserClaims)

	return claims, ok
}

func (app *application) writeJSON(w http.ResponseWriter, status int, body responsePayload) {
	jsonRsp, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	w.Write(jsonRsp)
}

func (app *application) stringFromQuery(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) intFromQuery(qs url.Values, key string, defaultValue int) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return i
}

func (app *application) dateFromQuery(qs url.Values, key string, defaultValue time.Time) time.Time {
	t, err := time.Parse("2006-01-02", qs.Get(key))
	if err != nil {
		return defaultValue
	}

	return t
}

func Version() string {
	var revision string
	var modified bool

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}
	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}

	return revision
}
