package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tldr-ai-be/internal/ratelimit"
	"tldr-ai-be/internal/usage"
)

func testEnvAll(k string) string {
	switch k {
	case "RATE_LIMIT_RPS":
		return "0"
	}
	return ""
}

func TestProcessText_503(t *testing.T) {
	h := NewHandler(&RouterDeps{Processor: nil, Budget: usage.NewFromEnv(testEnvAll)}, false, "", ratelimit.Init(testEnvAll))
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

func TestUsageGet(t *testing.T) {
	b := usage.NewFromEnv(testEnvAll)
	h := NewHandler(&RouterDeps{Processor: nil, Budget: b}, false, "", nil)
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	resp, err := http.Get(srv.URL + "/api/usage")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: %d", resp.StatusCode)
	}
	var out map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out["monthUtc"] == nil {
		t.Fatalf("json: %v", out)
	}
}

func TestUsageAdminReset(t *testing.T) {
	secret := "abc"
	b := usage.NewFromEnv(testEnvAll)
	_, _ = b.TryReserve() // in-flight
	h := NewHandler(&RouterDeps{Processor: nil, Budget: b, UsageResetSecret: secret}, false, "", nil)
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	req, _ := http.NewRequest("POST", srv.URL+"/api/admin/usage-reset", strings.NewReader(""))
	req.Header.Set(adminResetHeader, secret)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("status: %d", resp.StatusCode)
	}
	if s := b.Snapshot(); s.Used != 0 {
		t.Fatalf("used: %d", s.Used)
	}
}
