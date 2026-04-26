package domain

// ValidationError is returned for client-side fixable / bad input. Handlers map
// it to HTTP 400 via web.HandleError.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return "validation failed"
}

// NewValidationError returns a *ValidationError with a public message.
func NewValidationError(message string) *ValidationError {
	if message == "" {
		message = "Invalid request"
	}
	return &ValidationError{Message: message}
}
