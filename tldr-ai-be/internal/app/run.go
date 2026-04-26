package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"tldr-ai-be/internal/ai"
	"tldr-ai-be/internal/config"
	"tldr-ai-be/internal/httpapi"
	"tldr-ai-be/internal/ratelimit"
	"tldr-ai-be/internal/service"
	"tldr-ai-be/internal/usage"
)

const (
	defaultPort       = "8080"
	shutdownTimeout   = 15 * time.Second
	readHeaderTimeout = 5 * time.Second
	readWriteTimeout  = 2 * time.Minute
	idleTimeout       = 120 * time.Second
	maxHeaderBytes    = 1 << 20 // 1 MiB
)

func newHandler() http.Handler {
	return newHandlerWithEnv(os.Getenv)
}

func newHandlerWithEnv(get func(string) string) http.Handler {
	config.LogStartupEnvHint()
	p, err := ai.NewProviderFromEnv(get)
	if err != nil {
		log.Printf("startup: AI provider: %v (set a valid ANTHROPIC_API_KEY; POST /api/processText returns 503 until then)", err)
	}
	var proc *service.TextProcessor
	if p != nil {
		proc = service.NewTextProcessor(p)
	}
	trust := strings.TrimSpace(get("TRUST_PROXY")) == "1"
	cors := strings.TrimSpace(get("CORS_ALLOW_ORIGIN"))
	lim := ratelimit.Init(get)
	bud := usage.NewFromEnv(get)
	usageReset := strings.TrimSpace(get("USAGE_RESET_SECRET"))
	return httpapi.NewHandler(&httpapi.RouterDeps{
		Processor:        proc,
		Budget:           bud,
		UsageResetSecret: usageReset,
	}, trust, cors, lim)
}

// Run starts the HTTP server and blocks until SIGINT/SIGTERM or ListenAndServe fails.
func Run() error {
	port := config.GetEnv("PORT", defaultPort)
	addr := ":" + port

	srv := &http.Server{
		Addr:              addr,
		Handler:           newHandler(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readWriteTimeout,
		WriteTimeout:      readWriteTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("tldr-ai-be: started successfully — http://localhost:%s (API: /api/processText, /api/usage, /health) — press Ctrl+C to stop", port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(ln)
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
