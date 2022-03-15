package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pafirmin/do-daily-go/pkg/jwt"
)

type contextKey string

const ctxKeyUserClaims = contextKey("user")

func defaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			app.clientError(w, http.StatusUnauthorized)
			return
		}

		token := authHeader[1]
		claims, err := jwt.Parse(token)
		if err != nil {
			app.clientError(w, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyUserClaims, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
