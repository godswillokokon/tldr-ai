package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"tldr-ai-be/internal/domain"
	"tldr-ai-be/internal/errs"
	"tldr-ai-be/internal/httpclient"
)

const (
	anthropicMessagesURL = "https://api.anthropic.com/v1/messages"
	anthropicAPIVersion  = "2023-06-01"
	anthropicMaxTokens   = 400
)

// Anthropic implements Provider against the Anthropic Messages API.
type Anthropic struct {
	client *http.Client
	apiKey string
	model  string
}

// NewAnthropic builds a provider. A nil client uses httpclient.Default.
func NewAnthropic(client *http.Client, apiKey, model string) *Anthropic {
	if client == nil {
		client = httpclient.Default
	}
	return &Anthropic{client: client, apiKey: apiKey, model: model}
}

// ModelTag returns the configured model id.
func (a *Anthropic) ModelTag() string { return a.model }

type anthropicRequest struct {
	Model     string         `json:"model"`
	MaxTokens int            `json:"max_tokens"`
	Messages  []anthropicMsg `json:"messages"`
}

type anthropicMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

// Complete calls POST /v1/messages and returns the first "text" block, if any.
func (a *Anthropic) Complete(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(anthropicRequest{
		Model:     a.model,
		MaxTokens: anthropicMaxTokens,
		Messages: []anthropicMsg{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", errs.Internal("Could not build model request", err.Error())
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicMessagesURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", anthropicAPIVersion)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", errs.Internal("Model request failed", err.Error())
	}
	defer resp.Body.Close()

	lim := int64(domain.MaxModelResponseBytes) + 1
	limited := io.LimitReader(resp.Body, lim)
	data, err := io.ReadAll(limited)
	if err != nil {
		return "", errs.Internal("Could not read model response", err.Error())
	}
	if int64(len(data)) > int64(domain.MaxModelResponseBytes) {
		return "", errs.Internal("Model response is too large", fmt.Sprintf("bytes=%d", len(data)))
	}

	if resp.StatusCode >= 300 {
		return "", &errs.HTTPError{
			PublicMessage: "The model service returned an error",
			Status:        http.StatusBadGateway,
			LogDetail:     fmt.Sprintf("anthropic: http_status=%d", resp.StatusCode),
		}
	}

	var ar anthropicResponse
	if err := json.Unmarshal(data, &ar); err != nil {
		return "", errs.InvalidAIOutput("The AI service returned an invalid response", err.Error())
	}
	for i := range ar.Content {
		if ar.Content[i].Type == "text" && ar.Content[i].Text != "" {
			return ar.Content[i].Text, nil
		}
	}
	return "", errs.InvalidAIOutput("The AI service returned an invalid response", "no text content in response")
}
