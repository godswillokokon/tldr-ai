package errs

import "net/http"

// HTTPError is an error with a public message (JSON), HTTP status, and
// an optional private detail for server logs.
type HTTPError struct {
	PublicMessage string
	Status        int
	LogDetail     string
}

// Error implements error.
func (e *HTTPError) Error() string {
	if e == nil {
		return ""
	}
	return e.PublicMessage
}

// --- helpers ---

// BadRequest returns 400.
func BadRequest(public string) *HTTPError {
	if public == "" {
		public = "Bad request"
	}
	return &HTTPError{PublicMessage: public, Status: http.StatusBadRequest}
}

// TooManyRequests returns 429.
func TooManyRequests(public string) *HTTPError {
	if public == "" {
		public = "Too many requests"
	}
	return &HTTPError{PublicMessage: public, Status: http.StatusTooManyRequests}
}

// ServiceUnavailable returns 503.
func ServiceUnavailable(public string) *HTTPError {
	if public == "" {
		public = "Service temporarily unavailable"
	}
	return &HTTPError{PublicMessage: public, Status: http.StatusServiceUnavailable}
}

// BadGateway returns 502.
func BadGateway(public string) *HTTPError {
	if public == "" {
		public = "Bad gateway"
	}
	return &HTTPError{PublicMessage: public, Status: http.StatusBadGateway}
}

// Internal returns 500. Prefer setting LogDetail for logs.
func Internal(public, logDetail string) *HTTPError {
	if public == "" {
		public = "Internal server error"
	}
	return &HTTPError{PublicMessage: public, Status: http.StatusInternalServerError, LogDetail: logDetail}
}

// PayloadTooLarge returns 413.
func PayloadTooLarge(public string) *HTTPError {
	if public == "" {
		public = "Request body is too large"
	}
	return &HTTPError{PublicMessage: public, Status: http.StatusRequestEntityTooLarge}
}

// UsageCapExceeded returns 403.
func UsageCapExceeded(public string) *HTTPError {
	if public == "" {
		public = "Usage cap exceeded"
	}
	return &HTTPError{PublicMessage: public, Status: http.StatusForbidden}
}

// InvalidAIOutput is returned when upstream AI output cannot be used (e.g. empty or invalid schema).
// Maps to 502 so clients can retry; LogDetail can hold parse/safety details.
func InvalidAIOutput(public, logDetail string) *HTTPError {
	if public == "" {
		public = "The AI service returned an invalid response"
	}
	return &HTTPError{PublicMessage: public, Status: http.StatusBadGateway, LogDetail: logDetail}
}
