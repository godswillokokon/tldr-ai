package httpapi

import (
	"context"
	"net/http"
	"strings"
	"time"

	"tldr-ai-be/internal/domain"
	"tldr-ai-be/internal/errs"
	"tldr-ai-be/internal/service"
	"tldr-ai-be/internal/usage"
	"tldr-ai-be/internal/web"
)

const processTextTimeout = 42 * time.Second
const adminResetHeader = "X-Usage-Reset-Secret"

// RouterDeps holds dependencies for HTTP handlers.
type RouterDeps struct {
	Processor        *service.TextProcessor
	Budget           *usage.Budget
	UsageResetSecret string
}

// processText handles POST /api/processText: usage reserve, decode, 42s Process, commit/release, JSON 200.
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
	b := d.Budget
	var res *usage.Reservation
	if b != nil {
		var uerr error
		res, uerr = b.TryReserve()
		if uerr != nil {
			web.HandleError(w, r, uerr)
			return
		}
		defer b.Release(res)
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
	if b != nil && res != nil {
		b.Commit(res)
	}
	_ = web.WriteJSON(w, http.StatusOK, out)
}

// usageGet returns GET /api/usage.
func (d *RouterDeps) usageGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	_ = web.WriteJSON(w, http.StatusOK, d.Budget.Snapshot())
}

// usageAdminReset is POST /api/admin/usage-reset.
func (d *RouterDeps) usageAdminReset(w http.ResponseWriter, r *http.Request) {
	if d.UsageResetSecret == "" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if strings.TrimSpace(r.Header.Get(adminResetHeader)) != strings.TrimSpace(d.UsageResetSecret) {
		web.HandleError(w, r, errs.BadRequest("Invalid reset secret"))
		return
	}
	if d.Budget == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	d.Budget.AdminReset()
	w.WriteHeader(http.StatusNoContent)
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
