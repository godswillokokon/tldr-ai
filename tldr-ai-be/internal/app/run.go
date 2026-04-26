package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tldr-ai-be/internal/config"
	"tldr-ai-be/internal/middleware"
	"tldr-ai-be/internal/web"
)

const (
	defaultPort       = "8080"
	shutdownTimeout   = 15 * time.Second
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 120 * time.Second
	maxHeaderBytes    = 1 << 20 // 1 MiB
)

func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	h := http.Handler(mux)
	h = middleware.RequestID(h)
	h = middleware.SecurityHeaders(h)
	h = middleware.Recover(h)
	return h
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
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

// Run starts the HTTP server and blocks until SIGINT/SIGTERM or ListenAndServe fails.
func Run() error {
	port := config.GetEnv("PORT", defaultPort)
	addr := ":" + port

	srv := &http.Server{
		Addr:              addr,
		Handler:           newHandler(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}
		err := <-errCh
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}
