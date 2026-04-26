package ai

import (
	"strings"
	"testing"

	"tldr-ai-be/internal/domain"
)

func validJSON() string {
	return `{
  "summary": "A short test summary.",
  "actionItems": ["do one thing", "do two thing", "do three thing"]
}`
}

func TestParseModelPayload(t *testing.T) {
	t.Run("plain", func(t *testing.T) {
		out, err := ParseModelPayload(validJSON())
		if err != nil {
			t.Fatal(err)
		}
		if out.Summary != "A short test summary." {
			t.Fatalf("summary: %q", out.Summary)
		}
		if len(out.ActionItems) != 3 {
			t.Fatalf("actionItems: %v", out.ActionItems)
		}
	})

	t.Run("fenced", func(t *testing.T) {
		raw := "```json\n" + validJSON() + "\n```"
		out, err := ParseModelPayload(raw)
		if err != nil {
			t.Fatal(err)
		}
		if out.Summary == "" {
			t.Fatal("empty summary")
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		_, err := ParseModelPayload(`{"summary":"x","actionItems":["a","b","c"],"extra":1}`)
		if err == nil {
			t.Fatal("expected error")
		}
		if _, ok := err.(*domain.InvalidAIOutputError); !ok {
			t.Fatalf("got %T", err)
		}
	})

	t.Run("trailing", func(t *testing.T) {
		_, err := ParseModelPayload(validJSON() + "  null")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestParseModelPayload_tooBig(t *testing.T) {
	// Exceeds byte cap before meaningful parse.
	huge := `{"summary":"` + strings.Repeat("a", domain.MaxModelJSONBytes) + `"}`
	_, err := ParseModelPayload(huge)
	if err == nil {
		t.Fatal("expected error")
	}
}
