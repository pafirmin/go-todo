package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/pafirmin/go-todo/internal/data"
)

func TestGetTask(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
	}{
		{"Valid ID", "/tasks/1", http.StatusOK, []byte("Test"), "123"},
		{"Forbidden user", "/tasks/1", http.StatusForbidden, nil, "456"},
		{"Unauthorised user", "/tasks/1", http.StatusUnauthorized, nil, "invalid"},
		{"Non-existent ID", "/tasks/2", http.StatusNotFound, nil, "123"},
		{"Negative ID", "/tasks/-1", http.StatusNotFound, nil, "123"},
		{"Decimal ID", "/tasks/1.23", http.StatusNotFound, nil, "123"},
		{"String ID", "/tasks/foo", http.StatusNotFound, nil, "123"},
		{"Empty ID", "/tasks/", http.StatusNotFound, nil, "123"},
		{"Trailing slash", "/tasks/1/", http.StatusNotFound, nil, "123"},
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

func TestCreateTask(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
		dto      *data.CreateTaskDTO
	}{
		{"Valid request", "/folders/1/tasks", http.StatusCreated, []byte("Test"), "123",
			&data.CreateTaskDTO{Title: "Test", Description: "Test", Datetime: time.Now().Format(time.RFC3339)}},
		{"Invalid user", "/folders/1/tasks", http.StatusUnauthorized, nil, "invalid",
			&data.CreateTaskDTO{Title: "Test", Description: "Test", Datetime: time.Now().Format(time.RFC3339)}},
		{"Forbidden user", "/folders/1/tasks", http.StatusForbidden, nil, "456",
			&data.CreateTaskDTO{Title: "Test", Description: "Test", Datetime: time.Now().Format(time.RFC3339)}},
		{"Trailing slash", "/folders/1/tasks/", http.StatusNotFound, nil, "123",
			&data.CreateTaskDTO{Title: "Test", Description: "Test", Datetime: time.Now().Format(time.RFC3339)}},
		{"Invalid body", "/folders/1/tasks", http.StatusUnprocessableEntity, nil, "123",
			&data.CreateTaskDTO{Title: "", Description: "Test", Datetime: time.Now().Format(time.RFC3339)}},
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

func GetTasksByFolder(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
	}{
		{"Valid ID", "/folders/1/tasks", http.StatusOK, []byte("Test"), "123"},
		{"Invalid ID", "/folders/2/tasks", http.StatusNotFound, nil, "123"},
		{"Invalid user", "/folders/1/tasks", http.StatusUnauthorized, nil, "invalid"},
		{"Forbidden user", "/folders/1/tasks", http.StatusForbidden, nil, "456"},
		{"Trailing slash", "/folders/1/tasks/", http.StatusNotFound, nil, "123"},
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
