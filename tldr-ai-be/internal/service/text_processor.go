package service

import (
	"context"
	"strings"

	"tldr-ai-be/internal/ai"
	"tldr-ai-be/internal/domain"
)

// TextProcessor runs the end-to-end summarization pipeline: validate request,
// build the model prompt, call the provider, parse JSON, set model id, and
// validate the structured result. Errors are domain-typed or *errs.HTTPError
// from the model client; they are already safe to expose or log via web.HandleError.
type TextProcessor struct {
	provider ai.Provider
}

// NewTextProcessor returns a TextProcessor with injected model provider.
func NewTextProcessor(p ai.Provider) *TextProcessor {
	return &TextProcessor{provider: p}
}

// Process runs ValidateProcessRequest → BuildPrompt → provider.Complete →
// ParseModelPayload → set Model from ModelTag → ValidateProcessResponse.
func (t *TextProcessor) Process(ctx context.Context, in *domain.ProcessRequest) (*domain.ProcessResponse, error) {
	if err := domain.ValidateProcessRequest(in); err != nil {
		return nil, err
	}
	text := strings.TrimSpace(in.Text)
	prompt := ai.BuildPrompt(text)
	raw, err := t.provider.Complete(ctx, prompt)
	if err != nil {
		return nil, err
	}
	out, err := ai.ParseModelPayload(raw)
	if err != nil {
		// *domain.InvalidAIOutputError: safe public Message, LogDetail for logs
		return nil, err
	}
	out.Model = t.provider.ModelTag()
	if err := domain.ValidateProcessResponse(out); err != nil {
		return nil, err
	}
	return out, nil
}
