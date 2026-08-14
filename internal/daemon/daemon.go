// Package daemon holds the capsule daemon's process lifecycle: starting the
// HTTP listener, waiting for a shutdown signal, draining in-flight requests
// within a bounded timeout, and choosing the process exit code.
//
// This logic used to live directly in cmd/capsuled's run(), where it called
// flag.String against the global flag.CommandLine and so could not be
// constructed more than once inside a single test binary. Run in this
// package takes an already-built Config plus an injectable Deps (listener,
// clock, and signal source), so it can be exercised repeatedly and
// deterministically in tests. cmd/capsuled remains the only place flags are
// parsed; it builds a Config from them and calls Run.
package daemon

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DefaultShutdownTimeout is the bound on how long Run waits for in-flight
// requests to drain after a shutdown signal before forcing the listener
// closed. It matches the daemon's behavior prior to this package's
// extraction; Config.ShutdownTimeout falls back to it when unset.
const DefaultShutdownTimeout = 5 * time.Second

// Config is the daemon's already-built configuration. cmd/capsuled is
// responsible for parsing flags and constructing this from them; Run itself
// parses nothing.
type Config struct {
	// ListenAddress is the "host:port" the daemon listens on when Deps does
	// not supply a pre-bound Listener. It is also recorded in the startup
	// log line.
	ListenAddress string
	// Handler serves every accepted connection. cmd/capsuled builds this
	// from internal/api.NewServer.
	Handler http.Handler
	// ShutdownTimeout bounds how long Run waits for in-flight requests to
	// drain after a shutdown signal. Zero means DefaultShutdownTimeout.
	ShutdownTimeout time.Duration
	// Version is reported in the startup log line.
	Version string
}

// shutdownTimeout returns the effective shutdown timeout, defaulting an
// unset Config.ShutdownTimeout to DefaultShutdownTimeout.
func (cfg Config) shutdownTimeout() time.Duration {
	if cfg.ShutdownTimeout <= 0 {
		return DefaultShutdownTimeout
	}
	return cfg.ShutdownTimeout
}

// Clock abstracts the passage of time so tests can force the shutdown-
// timeout path without waiting out a real timeout.
type Clock interface {
	// After returns a channel that receives the current time once duration
	// d has elapsed, matching time.After's contract.
	After(d time.Duration) <-chan time.Time
}

// systemClock is the production Clock, backed by time.After.
type systemClock struct{}

// After implements Clock using the real wall clock.
func (systemClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// Deps carries the daemon's injectable environment: the listening socket,
// logger, signal source, and clock. Every field is optional; an unset field
// falls back to the same production behavior cmd/capsuled used before this
// package existed.
type Deps struct {
	// Listener, if set, is served directly instead of Run deriving one from
	// Config.ListenAddress. Tests use this to inject a pre-bound loopback
	// listener, or to observe the actual ephemeral port chosen by the OS.
	Listener net.Listener
	// Logger receives the daemon's structured lifecycle log lines. Defaults
	// to a JSON handler writing to os.Stderr, matching prior behavior.
	Logger *slog.Logger
	// NotifyContext returns a context canceled when one of the given
	// signals is received, and a matching stop function. Defaults to
	// signal.NotifyContext. Tests inject a fake to simulate a received
	// signal without sending a real one.
	NotifyContext func(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc)
	// Clock is used to bound the shutdown-drain wait. Defaults to the real
	// wall clock.
	Clock Clock
}

// withDefaults fills every unset Deps field with its production default.
func (deps Deps) withDefaults() Deps {
	if deps.Logger == nil {
		deps.Logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}
	if deps.NotifyContext == nil {
		deps.NotifyContext = signal.NotifyContext
	}
	if deps.Clock == nil {
		deps.Clock = systemClock{}
	}
	return deps
}

// Run executes the daemon's lifecycle to completion and returns the process
// exit code: it starts serving cfg.Handler, waits for a shutdown signal (via
// deps.NotifyContext, which is signal.NotifyContext for os.Interrupt and
// syscall.SIGTERM by default), drains in-flight requests within
// cfg.shutdownTimeout(), and reports whether serving stopped cleanly.
//
// This is a direct extraction of cmd/capsuled's former run(): the same
// listen handling, shutdown signal set, shutdown timeout, and exit codes,
// parameterized so it can be constructed and torn down repeatedly within one
// test binary.
func Run(cfg Config, deps Deps) int {
	deps = deps.withDefaults()

	server := &http.Server{
		Addr:              cfg.ListenAddress,
		Handler:           cfg.Handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := deps.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)

		<-ctx.Done()

		// shutdownCtx is canceled either when server.Shutdown finishes (the
		// deferred cancel below) or when deps.Clock signals that
		// cfg.shutdownTimeout() has elapsed, whichever comes first. Using
		// the injected clock here — rather than context.WithTimeout, which
		// is pinned to the real wall clock — is what lets a test force this
		// path without waiting out a real timeout.
		shutdownCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			select {
			case <-deps.Clock.After(cfg.shutdownTimeout()):
				cancel()
			case <-shutdownCtx.Done():
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			deps.Logger.Error("daemon shutdown failed", "error", err)
		}
	}()

	deps.Logger.Info("starting capsule daemon", "address", cfg.ListenAddress, "version", cfg.Version)

	var serveErr error
	if deps.Listener != nil {
		serveErr = server.Serve(deps.Listener)
	} else {
		serveErr = server.ListenAndServe()
	}

	// server.Serve/ListenAndServe can also return for a reason other than a
	// caught signal (for example a listen failure). Cancel ctx
	// unconditionally so the shutdown goroutine is never left blocked on
	// <-ctx.Done() waiting for a signal that will not arrive; stop is
	// idempotent, so the deferred call above stays safe.
	stop()

	// server.Shutdown causes Serve/ListenAndServe to return promptly, but
	// active handlers keep draining in the background goroutine above until
	// the shutdown timeout elapses or every connection finishes. Wait for
	// that goroutine so Run does not return while requests are still in
	// flight.
	<-shutdownDone

	if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		deps.Logger.Error("daemon stopped unexpectedly", "error", serveErr)
		return 1
	}
	return 0
}
