package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInit_disabled(t *testing.T) {
	l := Init(func(string) string { return "0" })
	if !l.Allow("9.9.9.9") {
		t.Fatal("disabled should allow")
	}
	rr := httptest.NewRecorder()
	var saw bool
	m := l.HTTPMiddleware(func(*http.Request) string { return "1.1.1.1" })(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { saw = true }))
	m.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))
	if !saw {
		t.Fatal("expected next")
	}
}

func TestLimiter_burstExhausted(t *testing.T) {
	l := newLimiter(1, 2)
	if !l.Allow("a") {
		t.Fatal(1)
	}
	if !l.Allow("a") {
		t.Fatal(2)
	}
	if l.Allow("a") {
		t.Fatal("expected drop")
	}
}

func TestLimiter_HTTP429(t *testing.T) {
	l := newLimiter(1, 8) // burst 2
	chain := l.HTTPMiddleware(func(*http.Request) string { return "9.1.1.1" })(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }),
	)
	codes := make([]int, 5)
	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		chain.ServeHTTP(rr, httptest.NewRequest("GET", "/", nil))
		codes[i] = rr.Code
	}
	if codes[0] != 200 || codes[1] != 200 {
		t.Fatalf("first two: %v", codes)
	}
	found := false
	for _, c := range codes[2:] {
		if c == http.StatusTooManyRequests {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected 429: %v", codes)
	}
}
