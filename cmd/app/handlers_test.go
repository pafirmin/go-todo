package main

import (
	"bytes"
	"net/http"
	"testing"
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
		{"Unauthorised user", "/folders/1", http.StatusForbidden, nil, "456"},
		{"Non-existent ID", "/folders/2", http.StatusNotFound, nil, "123"},
		{"Negative ID", "/folders/-1", http.StatusNotFound, nil, "123"},
		{"Decimal ID", "/folders/1.23", http.StatusNotFound, nil, "123"},
		{"String ID", "/folders/foo", http.StatusNotFound, nil, "123"},
		{"Empty ID", "/folders/", http.StatusNotFound, nil, "123"},
		{"Trailing slash", "/folders/1/", http.StatusNotFound, nil, "123"},
	}
	rp := requestPerformer(app.routes(), "GET", t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rp(tt.urlPath, tt.token)

			if code := r.Code; code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if body := r.Body.Bytes(); !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}
