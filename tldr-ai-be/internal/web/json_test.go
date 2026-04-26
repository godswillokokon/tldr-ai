package web

import (
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	_ = WriteJSON(rr, 201, struct {
		Foo int `json:"foo"`
	}{Foo: 1})
	if rr.Code != 201 {
		t.Fatalf("code: %d", rr.Code)
	}
	if g := rr.Header().Get("Content-Type"); g != "application/json; charset=utf-8" {
		t.Fatalf("ct: %q", g)
	}
	if s := rr.Body.String(); s != `{"foo":1}` {
		t.Fatalf("body: %q", s)
	}
}
