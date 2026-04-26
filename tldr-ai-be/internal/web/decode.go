package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"

	"tldr-ai-be/internal/domain"
	"tldr-ai-be/internal/errs"
)

// DecodeProcessRequest reads and decodes the request body as JSON into a ProcessRequest.
// Bodies may not exceed maxBytes. Unknown JSON fields, trailing JSON, or multiple values are rejected.
func DecodeProcessRequest(r *http.Request, maxBytes int64) (*domain.ProcessRequest, error) {
	if maxBytes < 0 {
		maxBytes = 0
	}
	limit := maxBytes
	if maxBytes < math.MaxInt64 {
		limit = maxBytes + 1
	}
	limited := io.LimitReader(r.Body, limit)
	b, err := io.ReadAll(limited)
	if err != nil {
		return nil, errs.Internal("Could not read request", err.Error())
	}
	if int64(len(b)) > maxBytes {
		return nil, errs.PayloadTooLarge("")
	}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var out domain.ProcessRequest
	if err := dec.Decode(&out); err != nil {
		return nil, mapJSONDecodeError(err)
	}
	if dec.More() {
		return nil, errs.BadRequest("Request body contains additional JSON data")
	}
	if rem := b[dec.InputOffset():]; len(bytes.TrimSpace(rem)) > 0 {
		return nil, errs.BadRequest("Request body contains additional JSON data")
	}
	return &out, nil
}

func mapJSONDecodeError(err error) *errs.HTTPError {
	if err == nil {
		return nil
	}
	if errors.Is(err, io.EOF) {
		return errs.BadRequest("Request body is empty")
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return &errs.HTTPError{
			PublicMessage: "Invalid JSON in request body",
			Status:        http.StatusBadRequest,
			LogDetail:     err.Error(),
		}
	}
	var se *json.SyntaxError
	if errors.As(err, &se) {
		return &errs.HTTPError{
			PublicMessage: "Invalid JSON in request body",
			Status:        http.StatusBadRequest,
			LogDetail:     fmt.Sprintf("json syntax: offset %d: %s", se.Offset, err.Error()),
		}
	}
	var te *json.UnmarshalTypeError
	if errors.As(err, &te) {
		return &errs.HTTPError{
			PublicMessage: "Invalid JSON in request body",
			Status:        http.StatusBadRequest,
			LogDetail:     fmt.Sprintf("json type: field %q: %s", te.Field, err.Error()),
		}
	}
	return &errs.HTTPError{
		PublicMessage: "Invalid JSON in request body",
		Status:        http.StatusBadRequest,
		LogDetail:     err.Error(),
	}
}
