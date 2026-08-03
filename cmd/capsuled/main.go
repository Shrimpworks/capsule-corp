// Command capsuled is the capsule daemon entrypoint. It currently only
// serves internal/api's read-only diagnostic endpoints on a loopback
// address by default; it holds no Approval, Supervisor, or execution
// authority (see internal/api's package doc and AGENTS.md).
package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"capsule.local/capsule/internal/api"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	listenAddress := flag.String("listen", "127.0.0.1:7777", "daemon listen address")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	serverHandler := api.NewServer(api.Options{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	})

	server := &http.Server{
		Addr:              *listenAddress,
		Handler:           serverHandler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("daemon shutdown failed", "error", err)
		}
	}()

	logger.Info("starting capsule daemon", "address", server.Addr, "version", version)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("daemon stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
