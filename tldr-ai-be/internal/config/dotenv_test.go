package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnv_missingFile(t *testing.T) {
	if err := LoadDotEnv(filepath.Join(t.TempDir(), "missing.file")); err != nil {
		t.Fatal(err)
	}
}

func TestLoadDotenvOverride_parses(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "f")
	if err := os.WriteFile(p, []byte("A=b\nB=\"c d\"\nexport E=f\n#x\nG=\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := LoadDotEnvOverride(p); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("A") != "b" {
		t.Fatalf("A=%q", os.Getenv("A"))
	}
	if os.Getenv("B") != "c d" {
		t.Fatalf("B=%q", os.Getenv("B"))
	}
	if os.Getenv("E") != "f" {
		t.Fatalf("E=%q", os.Getenv("E"))
	}
	t.Cleanup(func() {
		for _, k := range []string{"A", "B", "E", "G"} {
			_ = os.Unsetenv(k)
		}
	})
}

func TestLoadDotEnv_doesNotOverride(t *testing.T) {
	t.Setenv("K", "first")
	d := t.TempDir()
	p := filepath.Join(d, "f")
	if err := os.WriteFile(p, []byte("K=second\nM=1\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := LoadDotEnv(p); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("K") != "first" {
		t.Fatalf("K=%q", os.Getenv("K"))
	}
	if os.Getenv("M") != "1" {
		t.Fatalf("M=%q", os.Getenv("M"))
	}
	t.Cleanup(func() { _ = os.Unsetenv("M") })
}
