package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	defaultHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	contentType := rs.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("want %q; got %q", "1; mode=block", contentType)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestRecoverPanic(t *testing.T) {
	app := newTestApplication(t)
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("")
	})

	app.recoverPanic(next).ServeHTTP(rr, r)

	rs := rr.Result()

	if h := rs.Header.Get("Connection"); h != "close" {
		t.Errorf("want %q, got %q", "close", h)
	}
}

func TestRequireAuth(t *testing.T) {
	tests := []struct {
		token    string
		wantCode int
	}{
		{"Bearer 123", http.StatusOK},
		{"123", http.StatusUnauthorized},
		{"Bearer invalid", http.StatusUnauthorized},
	}

	for _, test := range tests {
		app := newTestApplication(t)
		rr := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.Header.Set("Authorization", test.token)

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		app.requireAuth(next).ServeHTTP(rr, r)

		rs := rr.Result()

		if sc := rs.StatusCode; sc != test.wantCode {
			t.Errorf("want %q, got %q", test.wantCode, rs.StatusCode)
		}
	}

}
