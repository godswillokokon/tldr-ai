package httpapi

import (
	"net/http/httptest"
	"testing"
)

func TestClientIP(t *testing.T) {
	t.Run("remote", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		r.RemoteAddr = "192.0.2.1:12345"
		if g := ClientIP(r, false); g != "192.0.2.1" {
			t.Fatalf("got %q", g)
		}
	})
	t.Run("x-forwarded", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		r.RemoteAddr = "10.0.0.1:1"
		r.Header.Set("X-Forwarded-For", "  192.0.2.2 , 10.0.0.1 ")
		if g := ClientIP(r, true); g != "192.0.2.2" {
			t.Fatalf("got %q", g)
		}
	})
	t.Run("no trust", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		r.RemoteAddr = "10.0.0.1:1"
		r.Header.Set("X-Forwarded-For", "192.0.2.2")
		if g := ClientIP(r, false); g != "10.0.0.1" {
			t.Fatalf("got %q", g)
		}
	})
}
