package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	mockJwt "github.com/pafirmin/do-daily-go/pkg/jwt/mock"
	"github.com/pafirmin/do-daily-go/pkg/models/mock"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog:   log.New(io.Discard, "", 0),
		folders:    &mock.FolderModel{},
		infoLog:    log.New(io.Discard, "", 0),
		jwtService: &mockJwt.JWTService{Secret: "123"},
		tasks:      &mock.TaskModel{},
		users:      &mock.UserModel{},
	}
}

func requestPerformer(r http.Handler, method string, t *testing.T) func(string, string) *httptest.ResponseRecorder {
	return func(path string, token string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(method, path, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		r.ServeHTTP(w, req)

		return w
	}
}
