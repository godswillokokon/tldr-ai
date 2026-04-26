package config

import (
	"log"
	"os"
	"strings"
)

// IsAnthropicPlaceholderKey reports true if the value is empty or looks like
// documentation / sample / placeholder text rather than a real key.
func IsAnthropicPlaceholderKey(key string) bool {
	k := strings.TrimSpace(key)
	if k == "" {
		return true
	}
	lower := strings.ToLower(k)
	needles := []string{
		"placeholder", "example", "changeme", "replace_me", "replace-me",
		"your-anthropic", "your-anthropic-key", "your_anthropic_key", "your-anthropic_api_key",
		"your_api_key", "insert-key", "put-key-here", "not-a-real",
		"test-key-only", "fake-key", "api_key_here", "00000000-0000",
	}
	for _, n := range needles {
		if strings.Contains(lower, n) {
			return true
		}
	}
	// Tutorials often use an all-literal mask; block obvious sk-ant- masks that are only 'x' or '-'.
	if strings.HasPrefix(k, "sk-ant-") {
		rest := strings.TrimPrefix(k, "sk-ant-")
		rest = strings.TrimLeft(rest, "0123456789") // drop "api03" style prefix
		if rest == "" {
			return true
		}
		var nonMask int
		for i := 0; i < len(rest); i++ {
			c := rest[i]
			if c != 'x' && c != 'X' && c != '-' {
				nonMask++
			}
		}
		if nonMask == 0 && len(rest) >= 4 {
			return true
		}
	}
	return false
}

// LogStartupEnvHint logs a single non-secret line when ANTHROPIC_API_KEY is missing
// or looks like a sample value, so operators know which env vars to set.
func LogStartupEnvHint() {
	k := strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY"))
	if k != "" && !IsAnthropicPlaceholderKey(k) {
		return
	}
	log.Print("tldr-ai-be: set ANTHROPIC_API_KEY to a real key (not your_anthropic_key or other placeholders); ANTHROPIC_MODEL is optional (default claude-sonnet-4-6).")
}
