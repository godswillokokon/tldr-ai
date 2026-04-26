package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"tldr-ai-be/internal/errs"
	"tldr-ai-be/internal/web"
)

// Recover recovers from panics, logs stack (via HandleError on an Internal HTTPError)
// and responds with 500 JSON.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			stack := string(debug.Stack())
			detail := fmt.Sprintf("panic: %v\n%s", rec, stack)
			web.HandleError(w, r, errs.Internal("Internal server error", detail))
		}()
		next.ServeHTTP(w, r)
	})
}
