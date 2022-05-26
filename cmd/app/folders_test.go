package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pafirmin/go-todo/internal/data"
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
	rm := getRequestMaker(app.routes(), "GET", t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rm("/api/v1"+tt.urlPath, "", tt.token)

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
	rm := getRequestMaker(app.routes(), "GET", t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rm("/api/v1"+tt.urlPath, "", tt.token)

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
		dto      *data.CreateFolderDTO
	}{
		{"Valid user", "/users/me/folders", http.StatusCreated, []byte("Test"), "123",
			&data.CreateFolderDTO{Name: "Test"}},
		{"Invalid user", "/users/me/folders", http.StatusUnauthorized, nil, "invalid",
			&data.CreateFolderDTO{Name: "Test"}},
		{"Trailing slash", "/users/me/folders/", http.StatusNotFound, nil, "123",
			&data.CreateFolderDTO{Name: "Test"}},
		{"Invalid body", "/users/me/folders", http.StatusBadRequest, nil, "123",
			&data.CreateFolderDTO{Name: ""}},
	}
	rm := getRequestMaker(app.routes(), "POST", t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.dto)
			r := rm("/api/v1"+tt.urlPath, string(body), tt.token)

			if code := r.Code; code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if body := r.Body.Bytes(); !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestDeleteFolder(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		token    string
	}{
		{"Valid request", "/folders/1", http.StatusNoContent, "123"},
		{"Invalid user", "/folders/1", http.StatusUnauthorized, "invalid"},
		{"Invalid ID", "/folders/test", http.StatusNotFound, "123"},
		{"Foribidden user", "/folders/1", http.StatusForbidden, "456"},
	}

	rm := getRequestMaker(app.routes(), "DELETE", t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rm("/api/v1"+tt.urlPath, "", tt.token)

			if code := r.Code; code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
		})
	}
}
