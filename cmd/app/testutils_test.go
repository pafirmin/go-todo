package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pafirmin/go-todo/internal/data"
	"github.com/pafirmin/go-todo/internal/data/mock"
	mockJwt "github.com/pafirmin/go-todo/internal/jwt/mock"
)

func newTestApplication(t *testing.T) *application {
	models := data.Models{
		Folders: mock.FolderModel{},
		Users:   mock.UserModel{},
		Tasks:   mock.TaskModel{},
		Tokens:  mock.TokenModel{},
	}
	return &application{
		errorLog:   log.New(io.Discard, "", 0),
		infoLog:    log.New(io.Discard, "", 0),
		jwtService: &mockJwt.JWTService{Secret: "123"},
		models:     models,
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
