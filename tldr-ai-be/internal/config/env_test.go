package config

import "testing"

func TestGetEnv(t *testing.T) {
	t.Setenv("TEST_GETENV_KEY", "")
	if g := GetEnv("TEST_GETENV_KEY", "fb"); g != "fb" {
		t.Fatalf("empty: got %q", g)
	}
	t.Setenv("TEST_GETENV_KEY", "  val  ")
	if g := GetEnv("TEST_GETENV_KEY", "fb"); g != "val" {
		t.Fatalf("trim: got %q", g)
	}
}
