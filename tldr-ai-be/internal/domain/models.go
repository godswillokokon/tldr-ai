package domain

// ProcessRequest is a client-issued summarization request.
type ProcessRequest struct {
	Text string `json:"text"`
}

// ProcessResponse is the expected structured output from the model.
type ProcessResponse struct {
	Summary     string   `json:"summary"`
	ActionItems []string `json:"actionItems"`
	Model       string   `json:"model,omitempty"`
}
