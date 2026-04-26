package ai

import (
	"fmt"
	"strings"

	"tldr-ai-be/internal/config"
	"tldr-ai-be/internal/httpclient"
)

const defaultAnthropicModel = "claude-sonnet-4-6"

// NewProviderFromEnv returns an Anthropic Provider using env (via get) for credentials.
// It requires a non-empty ANTHROPIC_API_KEY that is not a placeholder (see config.IsAnthropicPlaceholderKey).
// ANTHROPIC_MODEL defaults to "claude-sonnet-4-6" when empty.
func NewProviderFromEnv(get func(key string) string) (Provider, error) {
	key := strings.TrimSpace(get("ANTHROPIC_API_KEY"))
	if key == "" {
		return nil, fmt.Errorf("ai: ANTHROPIC_API_KEY is required")
	}
	if config.IsAnthropicPlaceholderKey(key) {
		return nil, fmt.Errorf("ai: ANTHROPIC_API_KEY must not be a placeholder or example value")
	}
	model := strings.TrimSpace(get("ANTHROPIC_MODEL"))
	if model == "" {
		model = defaultAnthropicModel
	}
	return NewAnthropic(httpclient.Default, key, model), nil
}
