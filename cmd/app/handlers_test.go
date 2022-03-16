package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pafirmin/do-daily-go/pkg/models/postgres"
)

func TestGetFolder(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
	}{
		{"Valid ID", "/folders/1", http.StatusOK, []byte("Test"), "123"},
		{"Forbidden user", "/folders/1", http.StatusForbidden, nil, "456"},
		{"Unauthorised user", "/folders/1", http.StatusUnauthorized, nil, "invalid"},
		{"Non-existent ID", "/folders/2", http.StatusNotFound, nil, "123"},
		{"Negative ID", "/folders/-1", http.StatusNotFound, nil, "123"},
		{"Decimal ID", "/folders/1.23", http.StatusNotFound, nil, "123"},
		{"String ID", "/folders/foo", http.StatusNotFound, nil, "123"},
		{"Empty ID", "/folders/", http.StatusNotFound, nil, "123"},
		{"Trailing slash", "/folders/1/", http.StatusNotFound, nil, "123"},
	}
	rp := getRequest(app.routes(), t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rp("/api/v1"+tt.urlPath, tt.token)

			if code := r.Code; code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if body := r.Body.Bytes(); !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestGetFoldersByUser(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
	}{
		{"Valid ID", "/users/me/folders", http.StatusOK, []byte("Test"), "123"},
		{"Invalid user", "/users/me/folders", http.StatusUnauthorized, nil, "invalid"},
		{"Trailing slash", "/users/me/folders/", http.StatusNotFound, nil, "123"},
	}
	rp := getRequest(app.routes(), t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rp("/api/v1"+tt.urlPath, tt.token)

			if code := r.Code; code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if body := r.Body.Bytes(); !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestCreateFolder(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
		dto      *postgres.CreateFolderDTO
	}{
		{"Valid user", "/users/me/folders", http.StatusCreated, []byte("Test"), "123",
			&postgres.CreateFolderDTO{Name: "Test"}},
		{"Invalid user", "/users/me/folders", http.StatusUnauthorized, nil, "invalid",
			&postgres.CreateFolderDTO{Name: "Test"}},
		{"Trailing slash", "/users/me/folders/", http.StatusNotFound, nil, "123",
			&postgres.CreateFolderDTO{Name: "Test"}},
		{"Invalid body", "/users/me/folders", http.StatusBadRequest, nil, "123",
			&postgres.CreateFolderDTO{Name: ""}},
	}
	rp := postRequest(app.routes(), t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.dto)
			r := rp("/api/v1"+tt.urlPath, string(body), tt.token)

			if code := r.Code; code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if body := r.Body.Bytes(); !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
		dto      *postgres.CreateUserDTO
	}{
		{"Valid body", "/users", http.StatusCreated, []byte("mock@example.com"), "123",
			&postgres.CreateUserDTO{Email: "mock@example.com", Password: "Test1234"}},
		{"Invalid email", "/users", http.StatusBadRequest, nil, "123",
			&postgres.CreateUserDTO{Email: "invalid", Password: "Test1234"}},
		{"Invalid password", "/users", http.StatusBadRequest, nil, "123",
			&postgres.CreateUserDTO{Email: "mock@example.com", Password: "Test123"}},
		{"Trailing slash", "/users/", http.StatusNotFound, nil, "123",
			&postgres.CreateUserDTO{Email: "mock@example.com", Password: "Test1234"}},
	}
	rp := postRequest(app.routes(), t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.dto)
			r := rp("/api/v1"+tt.urlPath, string(body), tt.token)

			if code := r.Code; code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if body := r.Body.Bytes(); !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q; got %q", tt.wantBody, r.Body)
			}
		})
	}
}
