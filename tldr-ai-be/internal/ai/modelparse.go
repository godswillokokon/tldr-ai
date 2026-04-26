package ai

import (
	"bytes"
	"encoding/json"
	"strings"

	"tldr-ai-be/internal/domain"
)

type modelMessageJSON struct {
	Summary     string   `json:"summary"`
	ActionItems []string `json:"actionItems"`
}

// ParseModelPayload trims optional ```/```json code fences, unmarshals the JSON payload, and
// returns a domain.ProcessResponse without a Model (caller may set it from Provider.ModelTag()).
// The UTF-8 byte size must not exceed domain.MaxModelJSONBytes.
func ParseModelPayload(raw string) (*domain.ProcessResponse, error) {
	s := strings.TrimSpace(raw)
	s = trimJSONFences(s)
	if len(s) > domain.MaxModelJSONBytes {
		return nil, &domain.InvalidAIOutputError{
			Message:   "The AI service returned an invalid response",
			LogDetail: "model json exceeds size limit",
		}
	}
	dec := json.NewDecoder(bytes.NewReader([]byte(s)))
	dec.DisallowUnknownFields()
	var m modelMessageJSON
	if err := dec.Decode(&m); err != nil {
		return nil, &domain.InvalidAIOutputError{
			Message:   "The AI service returned an invalid response",
			LogDetail: err.Error(),
		}
	}
	// Reject any trailing data after a single value.
	if dec.More() {
		return nil, &domain.InvalidAIOutputError{
			Message:   "The AI service returned an invalid response",
			LogDetail: "trailing data after model JSON",
		}
	}
	rem := s[dec.InputOffset():]
	if strings.TrimSpace(rem) != "" {
		return nil, &domain.InvalidAIOutputError{
			Message:   "The AI service returned an invalid response",
			LogDetail: "trailing data after model JSON",
		}
	}
	return &domain.ProcessResponse{Summary: m.Summary, ActionItems: m.ActionItems}, nil
}

func trimJSONFences(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "json")
	s = strings.TrimSpace(s)
	if i := strings.LastIndex(s, "```"); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}
