package ai

import "context"

// Provider calls a hosted model with a string prompt and returns the raw text
// response (e.g. JSON string from the first text block).
type Provider interface {
	Complete(ctx context.Context, prompt string) (string, error)
	ModelTag() string
}
