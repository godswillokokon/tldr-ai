package config

import "strings"

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
		"your-anthropic", "your_api_key", "insert-key", "put-key-here", "not-a-real",
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
