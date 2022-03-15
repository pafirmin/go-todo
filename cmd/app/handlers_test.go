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
	}{
		{"Valid ID", "/folders/1", http.StatusOK, []byte("Test")},
		{"Non-existent ID", "/folders/2", http.StatusNotFound, nil},
		{"Negative ID", "/folders/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/folders/1.23", http.StatusNotFound, nil},
		{"String ID", "/folders/foo", http.StatusNotFound, nil},
		{"Empty ID", "/folders/", http.StatusNotFound, nil},
		{"Trailing slash", "/folders/1/", http.StatusNotFound, nil},
	}
	rp := requestPerformer(app.routes(), "GET", t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rp(tt.urlPath)

			if code := r.Code; code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if body := r.Body.Bytes(); !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}
