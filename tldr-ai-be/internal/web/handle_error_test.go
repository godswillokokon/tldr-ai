package web

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"tldr-ai-be/internal/domain"
	"tldr-ai-be/internal/errs"
)

func TestHandleError_HTTPError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Request-ID", "req-1")
	rr := httptest.NewRecorder()
	HandleError(rr, req, errs.BadRequest("nope"))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rr.Code)
	}
	var body map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&body)
	if body["error"] != "nope" {
		t.Fatalf("body: %v", body)
	}
}

func TestHandleError_validation(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Request-ID", "r2")
	rr := httptest.NewRecorder()
	HandleError(rr, req, domain.NewValidationError("bad field"))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status: %d", rr.Code)
	}
	b, _ := io.ReadAll(rr.Body)
	if string(b) != `{"error":"bad field"}` {
		t.Fatalf("body: %s", b)
	}
}

func TestHandleError_unexpected(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Request-ID", "r3")
	rr := httptest.NewRecorder()
	HandleError(rr, req, io.EOF)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status: %d", rr.Code)
	}
	b, _ := io.ReadAll(rr.Body)
	if string(b) != `{"error":"Internal server error"}` {
		t.Fatalf("body: %s", b)
	}
}

func TestHandleError_nil(t *testing.T) {
	rr := httptest.NewRecorder()
	HandleError(rr, httptest.NewRequest(http.MethodGet, "/", nil), nil)
	if rr.Body.Len() != 0 || rr.Header().Get("Content-Type") != "" {
		t.Fatalf("expected no response body or json header, code=%d body=%q", rr.Code, rr.Body.String())
	}
}
