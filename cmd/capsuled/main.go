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
	os.Exit(run())
}

// run holds the daemon lifecycle so every deferred cleanup (signal-context
// stop, in particular) unwinds before the process exits.
func run() int {
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

	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)

		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("daemon shutdown failed", "error", err)
		}
	}()

	logger.Info("starting capsule daemon", "address", server.Addr, "version", version)

	serveErr := server.ListenAndServe()

	// ListenAndServe can also return for a reason other than a caught signal
	// (for example a listen failure). Cancel ctx unconditionally so the
	// shutdown goroutine is never left blocked on <-ctx.Done() waiting for a
	// signal that will not arrive; stop is idempotent, so the deferred call
	// above stays safe.
	stop()

	// server.Shutdown causes ListenAndServe to return promptly, but active
	// handlers keep draining in the background goroutine above until its
	// five-second window elapses or every connection finishes. Wait for that
	// goroutine so main does not os.Exit while requests are still in flight.
	<-shutdownDone

	if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		logger.Error("daemon stopped unexpectedly", "error", serveErr)
		return 1
	}
	return 0
}
