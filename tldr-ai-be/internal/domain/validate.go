package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ValidateProcessRequest checks trimmed input length in runes and bytes.
func ValidateProcessRequest(in *ProcessRequest) error {
	if in == nil {
		return &ValidationError{Message: "Request body is required", LogDetail: "nil ProcessRequest"}
	}
	t := strings.TrimSpace(in.Text)
	if t == "" {
		return &ValidationError{Message: "Text is required", LogDetail: "text empty or whitespace only"}
	}
	if n := utf8.RuneCountInString(t); n < MinInputTextRunes {
		return &ValidationError{
			Message:   fmt.Sprintf("Text must be at least %d characters", MinInputTextRunes),
			LogDetail: fmt.Sprintf("rune count=%d after trim", n),
		}
	}
	if len(t) > MaxInputTextBytes {
		return &ValidationError{
			Message:   "Text is too long",
			LogDetail: fmt.Sprintf("byte length=%d max=%d", len(t), MaxInputTextBytes),
		}
	}
	return nil
}

// ValidateProcessResponse enforces a fixed schema: non-empty summary, exactly
// three non-empty action items, and rune limits per field.
func ValidateProcessResponse(in *ProcessResponse) error {
	if in == nil {
		return &ValidationError{Message: "Response is required", LogDetail: "nil ProcessResponse"}
	}
	summary := strings.TrimSpace(in.Summary)
	if summary == "" {
		return &ValidationError{Message: "Summary is required", LogDetail: "summary empty or whitespace only"}
	}
	if n := utf8.RuneCountInString(summary); n > MaxSummaryRunes {
		return &ValidationError{
			Message:   "Summary is too long",
			LogDetail: fmt.Sprintf("summary rune count=%d max=%d", n, MaxSummaryRunes),
		}
	}
	if len(in.ActionItems) != 3 {
		return &ValidationError{
			Message:   "Expected exactly 3 action items",
			LogDetail: fmt.Sprintf("len(action_items)=%d", len(in.ActionItems)),
		}
	}
	for i, it := range in.ActionItems {
		item := strings.TrimSpace(it)
		if item == "" {
			return &ValidationError{
				Message:   "Each action item must be non-empty",
				LogDetail: fmt.Sprintf("action_items[%d] is empty", i),
			}
		}
		if n := utf8.RuneCountInString(item); n > MaxActionItemRunes {
			return &ValidationError{
				Message:   "An action item is too long",
				LogDetail: fmt.Sprintf("action_items[%d] rune count=%d max=%d", i, n, MaxActionItemRunes),
			}
		}
	}
	return nil
}
