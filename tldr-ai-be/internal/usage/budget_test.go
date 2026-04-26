package usage

import (
	"errors"
	"testing"

	"tldr-ai-be/internal/errs"
)

func testGet(m map[string]string) func(key string) string {
	return func(k string) string {
		if m == nil {
			return ""
		}
		return m[k]
	}
}

func TestNewFromEnv_defaults(t *testing.T) {
	b := NewFromEnv(testGet(nil))
	if b.perCall != defaultPerCallUSD {
		t.Fatalf("perCall: %v", b.perCall)
	}
}

func TestTryReserve_UsdCap(t *testing.T) {
	b := NewFromEnv(testGet(map[string]string{
		"USAGE_BUDGET_USD":   "0.01",
		"USAGE_PER_CALL_USD": "0.01",
	}))
	_, err := b.TryReserve()
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.TryReserve()
	if err == nil {
		t.Fatal("expected cap")
	}
	var he *errs.HTTPError
	if !errors.As(err, &he) || he.Status != 403 {
		t.Fatalf("got %v", err)
	}
}

func TestTryReserve_CallCap(t *testing.T) {
	b := NewFromEnv(testGet(map[string]string{
		"USAGE_MAX_CALLS": "1",
	}))
	_, err := b.TryReserve()
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.TryReserve()
	if err == nil {
		t.Fatal("expected cap")
	}
}

func TestCommitRelease(t *testing.T) {
	b := NewFromEnv(testGet(map[string]string{
		"USAGE_BUDGET_USD":   "1",
		"USAGE_PER_CALL_USD": "0.1",
	}))
	r, err := b.TryReserve()
	if err != nil {
		t.Fatal(err)
	}
	b.Release(r)
	if len(b.pending) != 0 {
		t.Fatalf("pending: %d", len(b.pending))
	}
	r2, err := b.TryReserve()
	if err != nil {
		t.Fatal(err)
	}
	b.Commit(r2)
	if b.usedTotal != 1 {
		t.Fatalf("used: %d", b.usedTotal)
	}
}

func TestAdminReset(t *testing.T) {
	b := NewFromEnv(testGet(map[string]string{"USAGE_MAX_CALLS": "5"}))
	r, _ := b.TryReserve()
	b.Commit(r)
	if b.Snapshot().Used != 1 {
		t.Fatal()
	}
	b.AdminReset()
	if b.Snapshot().Used != 0 {
		t.Fatal()
	}
}
