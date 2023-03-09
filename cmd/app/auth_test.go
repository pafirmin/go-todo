package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/pafirmin/go-todo/internal/data"
)

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
