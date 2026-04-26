package httpapi

import (
	"context"
	"net/http"
	"strings"
	"time"

	"tldr-ai-be/internal/domain"
	"tldr-ai-be/internal/errs"
	"tldr-ai-be/internal/service"
	"tldr-ai-be/internal/web"
)

const processTextTimeout = 42 * time.Second

// RouterDeps holds dependencies for HTTP handlers.
type RouterDeps struct {
	Processor *service.TextProcessor
}

// processText handles POST /api/processText (Content-Type, decode, 42s timeout, JSON 200 or HandleError).
func (d *RouterDeps) processText(w http.ResponseWriter, r *http.Request) {
	if d.Processor == nil {
		web.HandleError(w, r, errs.ServiceUnavailable("Summarization is not configured for this instance"))
		return
	}
	if h := r.Header.Get("Content-Type"); h != "" {
		ct := strings.ToLower(strings.TrimSpace(strings.Split(h, ";")[0]))
		if ct != "application/json" {
			web.HandleError(w, r, errs.BadRequest("Content-Type must be application/json when provided"))
			return
		}
	}
	in, err := web.DecodeProcessRequest(r, int64(domain.MaxRequestBodyBytes))
	if err != nil {
		web.HandleError(w, r, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), processTextTimeout)
	defer cancel()
	out, err := d.Processor.Process(ctx, in)
	if err != nil {
		web.HandleError(w, r, err)
		return
	}
	_ = web.WriteJSON(w, http.StatusOK, out)
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodHead {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		return
	}
	_ = web.WriteJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{Status: "ok"})
}
