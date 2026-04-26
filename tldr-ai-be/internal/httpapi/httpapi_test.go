package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tldr-ai-be/internal/ratelimit"
)

func testEnvAll(k string) string {
	switch k {
	case "RATE_LIMIT_RPS":
		return "0"
	}
	return ""
}

func TestProcessText_503(t *testing.T) {
	h := NewHandler(&RouterDeps{Processor: nil}, false, "", ratelimit.Init(testEnvAll))
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	resp, err := http.Post(srv.URL+"/api/processText", "application/json", strings.NewReader(`{"text":""}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status: %d", resp.StatusCode)
	}
}
