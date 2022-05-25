package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

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
			&data.CreateTaskDTO{Title: "Test", Description: "Test", Priority: "low", Due: time.Now().String()}},
		{"Invalid user", "/folders/1/tasks", http.StatusUnauthorized, nil, "invalid",
			&data.CreateTaskDTO{Title: "Test", Description: "Test", Priority: "low", Due: time.Now().String()}},
		{"Forbidden user", "/folders/1/tasks", http.StatusForbidden, nil, "456",
			&data.CreateTaskDTO{Title: "Test", Description: "Test", Priority: "low", Due: time.Now().String()}},
		{"Trailing slash", "/folders/1/tasks/", http.StatusNotFound, nil, "123",
			&data.CreateTaskDTO{Title: "Test", Description: "Test", Priority: "low", Due: time.Now().String()}},
		{"Invalid body", "/folders/1/tasks", http.StatusBadRequest, nil, "123",
			&data.CreateTaskDTO{Title: "", Description: "Test", Priority: "low", Due: time.Now().String()}},
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

func TestCreateUser(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
		dto      *data.CreateUserDTO
	}{
		{"Valid body", "/users", http.StatusCreated, []byte("mock@example.com"), "123",
			&data.CreateUserDTO{Email: "mock@example.com", Password: "Test1234"}},
		{"Invalid email", "/users", http.StatusBadRequest, nil, "123",
			&data.CreateUserDTO{Email: "invalid", Password: "Test1234"}},
		{"Invalid password", "/users", http.StatusBadRequest, nil, "123",
			&data.CreateUserDTO{Email: "mock@example.com", Password: "Test123"}},
		{"Trailing slash", "/users/", http.StatusNotFound, nil, "123",
			&data.CreateUserDTO{Email: "mock@example.com", Password: "Test1234"}},
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
				t.Errorf("want body to contain %q; got %q", tt.wantBody, r.Body)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	app := newTestApplication(t)

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
		token    string
		dto      *data.Credentials
	}{
		{"Valid credentials", "/auth/login", http.StatusOK, []byte("123"), "",
			&data.Credentials{Email: "mock@example.com", Password: "Test1234"}},
		{"Invalid credentials", "/auth/login", http.StatusUnauthorized, nil, "",
			&data.Credentials{Email: "invalid", Password: "Test1234"}},
		{"Trailing slash", "/auth/login/", http.StatusNotFound, nil, "",
			&data.Credentials{Email: "mock@example.com", Password: "Test1234"}},
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
				t.Errorf("want body to contain %q; got %q", tt.wantBody, r.Body)
			}
		})
	}
}
