package ai

import "strings"

// BuildPrompt builds the system-style instruction plus the user text for a single user message.
// The model is instructed to emit JSON only with "summary" and "actionItems" (exactly 3 items).
func BuildPrompt(input string) string {
	var b strings.Builder
	b.WriteString("You are a strict JSON generator. Output a single JSON object and nothing else (no markdown, no code fences, no commentary).\n")
	b.WriteString("The object must have exactly these keys: \"summary\" (string) and \"actionItems\" (array of exactly 3 strings).\n")
	b.WriteString("Requirements:\n")
	b.WriteString("- \"summary\" must be 1-2 short plain sentences in English.\n")
	b.WriteString("- \"actionItems\" must contain exactly 3 non-empty strings.\n")
	b.WriteString("- Do not use markdown in field values. Do not use leading numbers, bullets, or prefixes in action item strings; each item must be a single plain sentence.\n")
	b.WriteString("Text:\n")
	b.WriteString(input)
	return b.String()
}
