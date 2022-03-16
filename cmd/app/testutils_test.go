package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	mockJwt "github.com/pafirmin/go-todo/pkg/jwt/mock"
	"github.com/pafirmin/go-todo/pkg/models/mock"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog:   log.New(io.Discard, "", 0),
		folders:    &mock.FolderModel{},
		infoLog:    log.New(io.Discard, "", 0),
		jwtService: &mockJwt.JWTService{Secret: "123"},
		tasks:      &mock.TaskModel{},
		users:      &mock.UserModel{},
		validator:  validator.New(),
	}
}

func getRequestMaker(r http.Handler, method string, t *testing.T) func(string, string, string) *httptest.ResponseRecorder {
	return func(path, body, token string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		r.ServeHTTP(w, req)

		return w
	}
}
