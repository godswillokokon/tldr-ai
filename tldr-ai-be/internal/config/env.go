package config

import (
	"os"
	"strings"
)

// GetEnv returns os.Getenv(key) trimmed, or fallback when empty.
func GetEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
