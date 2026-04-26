package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestID_generates(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/x", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h := RequestID(mux)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	h.ServeHTTP(rr, req)
	if id := rr.Header().Get("X-Request-ID"); id == "" || id == "unknown" {
		t.Fatalf("id: %q", id)
	}
}

func TestRequestID_echo(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/x", func(w http.ResponseWriter, r *http.Request) {})
	h := RequestID(mux)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Request-ID", "  client-id  ")
	h.ServeHTTP(rr, req)
	if got := strings.TrimSpace(rr.Header().Get("X-Request-ID")); got != "client-id" {
		t.Fatalf("echo: %q", got)
	}
}

func TestSecurityHeaders_nosniff(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	h := SecurityHeaders(mux)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("header: %q", rr.Header().Get("X-Content-Type-Options"))
	}
}

func TestRecover_panicTo500(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})
	h := Recover(SecurityHeaders(RequestID(CORS("")(mux))))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/panic", nil))
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Internal server error") {
		t.Fatalf("body: %s", rr.Body.String())
	}
	if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("ct: %q", ct)
	}
}
