package web

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"tldr-ai-be/internal/domain"
	"tldr-ai-be/internal/errs"
)

const maxLogDetail = 1024

type errorResponse struct {
	Error string `json:"error"`
}

// HandleError maps err to a JSON error response. It never writes for a nil err.
// *errs.HTTPError uses status + PublicMessage; LogDetail is written to logs
// (truncated) with the request id when set.
// *domain.ValidationError returns 400.
// *domain.InvalidAIOutputError returns 502 with a safe message; LogDetail is logged.
// Any other error returns 500 with a generic message; the full error is logged
// (truncated) with the request id.
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	reqID := strings.TrimSpace(r.Header.Get("X-Request-ID"))
	logPrefix := "[request_id=" + reqID + "] "

	var he *errs.HTTPError
	if errors.As(err, &he) {
		if he.LogDetail != "" {
			log.Printf("%shttp_error: %s (detail: %s)", logPrefix, he.PublicMessage, truncate(he.LogDetail, maxLogDetail))
		} else {
			log.Printf("%shttp_error: %s", logPrefix, he.PublicMessage)
		}
		_ = WriteJSON(w, he.Status, errorResponse{Error: he.PublicMessage})
		return
	}

	var ve *domain.ValidationError
	if errors.As(err, &ve) {
		if ve.LogDetail != "" {
			log.Printf("%svalidation: %s (detail: %s)", logPrefix, ve.Message, truncate(ve.LogDetail, maxLogDetail))
		} else {
			log.Printf("%svalidation: %s", logPrefix, ve.Message)
		}
		_ = WriteJSON(w, http.StatusBadRequest, errorResponse{Error: ve.Error()})
		return
	}

	var aio *domain.InvalidAIOutputError
	if errors.As(err, &aio) {
		if aio.LogDetail != "" {
			log.Printf("%sinvalid_ai_output: %s (detail: %s)", logPrefix, aio.Message, truncate(aio.LogDetail, maxLogDetail))
		} else {
			log.Printf("%sinvalid_ai_output: %s", logPrefix, aio.Message)
		}
		_ = WriteJSON(w, http.StatusBadGateway, errorResponse{Error: aio.Error()})
		return
	}

	log.Printf("%sunexpected error: %s", logPrefix, truncate(err.Error(), maxLogDetail))
	_ = WriteJSON(w, http.StatusInternalServerError, errorResponse{Error: "Internal server error"})
}

func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
