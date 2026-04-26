package domain

// ValidationError is client-fixable input / contract failures.
// Message is safe to return in API JSON; LogDetail is for server logs only.
type ValidationError struct {
	Message   string
	LogDetail string
}

func (e *ValidationError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return "Invalid request"
}

// InvalidAIOutputError means upstream produced unusable or invalid structured output.
// Message is safe to return; LogDetail is for server logs (parse/schema details).
type InvalidAIOutputError struct {
	Message   string
	LogDetail string
}

func (e *InvalidAIOutputError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return "The AI service returned an invalid response"
}
