package ai

import (
	"strings"
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	in := "hello there this is at least twenty runes of text for sure"
	p := BuildPrompt(in)
	if !strings.HasSuffix(p, in) {
		t.Fatalf("expected user text appended, got: %q", p)
	}
	if !strings.Contains(p, "Text:\n") {
		t.Fatalf("expected 'Text:' section")
	}
	for _, needle := range []string{`"summary"`, "actionItems", "exactly 3", "no markdown", "JSON"} {
		if !strings.Contains(p, needle) {
			t.Fatalf("expected prompt to contain %q", needle)
		}
	}
}
