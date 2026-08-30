package daemon

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

// newTestLogger returns a logger that also writes to a buffer, so tests can
// assert on specific log lines (e.g. that a shutdown error was recorded)
// without depending on stderr.
func newTestLogger(t *testing.T) (*slog.Logger, *bytes.Buffer) {
	t.Helper()

	var buf bytes.Buffer
	return slog.New(slog.NewJSONHandler(&buf, nil)), &buf
}

// fakeSignalSource returns a NotifyContext replacement that ignores the
// requested signals and instead hands back a context the test controls
// directly, standing in for "a signal was received" without sending a real
// OS signal.
func fakeSignalSource(ctx context.Context, cancel context.CancelFunc) func(context.Context, ...os.Signal) (context.Context, context.CancelFunc) {
	return func(context.Context, ...os.Signal) (context.Context, context.CancelFunc) {
		return ctx, cancel
	}
}

// instantClock is a Clock whose After fires immediately regardless of the
// requested duration, used to force the shutdown-timeout path without
// waiting out a real timeout.
type instantClock struct{}

func (instantClock) After(time.Duration) <-chan time.Time {
	fired := make(chan time.Time, 1)
	fired <- time.Now()
	return fired
}

// neverClock is a Clock whose After never fires, used when a test wants the
// shutdown timeout goroutine to never win the race against server.Shutdown
// completing on its own.
type neverClock struct{}

func (neverClock) After(time.Duration) <-chan time.Time {
	return make(chan time.Time)
}

// blockingHandler serves GET / by blocking until unblock is closed, letting
// a test hold a connection open across a shutdown attempt.
type blockingHandler struct {
	received chan struct{}
	unblock  chan struct{}
}

func newBlockingHandler() *blockingHandler {
	return &blockingHandler{
		received: make(chan struct{}),
		unblock:  make(chan struct{}),
	}
}

func (h *blockingHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	close(h.received)
	<-h.unblock
	w.WriteHeader(http.StatusOK)
}

// listenLoopback opens a loopback listener on an OS-assigned ephemeral port.
func listenLoopback(t *testing.T) net.Listener {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on loopback: %v", err)
	}
	return listener
}

func TestRun_StartsAndShutsDownCleanlyOnSignal(t *testing.T) {
	listener := listenLoopback(t)
	logger, logs := newTestLogger(t)

	signalCtx, cancelSignal := context.WithCancel(context.Background())

	runDone := make(chan int, 1)
	go func() {
		runDone <- Run(Config{
			ListenAddress: listener.Addr().String(),
			Handler:       http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }),
		}, Deps{
			Listener:      listener,
			Logger:        logger,
			NotifyContext: fakeSignalSource(signalCtx, cancelSignal),
			Clock:         instantClock{},
		})
	}()

	// Confirm the server is actually serving before triggering shutdown.
	resp, err := http.Get("http://" + listener.Addr().String() + "/healthz")
	if err != nil {
		t.Fatalf("request before shutdown: %v", err)
	}
	_ = resp.Body.Close()

	cancelSignal()

	select {
	case code := <-runDone:
		if code != 0 {
			t.Fatalf("expected exit code 0 on clean signal shutdown, got %d; logs: %s", code, logs.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after signal shutdown")
	}
}

// realSIGTERMHelperEnv gates TestHelperProcess_RealSIGTERM: set only in the
// child process TestRun_RealSIGTERMShutsDownCleanly re-executes itself as, so
// an ordinary `go test` run treats that function as a normal (no-op) test.
const realSIGTERMHelperEnv = "CAPSULE_DAEMON_REAL_SIGTERM_HELPER"

// TestHelperProcess_RealSIGTERM is not a real test on its own; it is a
// well-known name TestRun_RealSIGTERMShutsDownCleanly re-executes the test
// binary as (via `-test.run` and realSIGTERMHelperEnv) so the real default
// signal.NotifyContext path can be exercised with a genuine SIGTERM sent to
// an owned child process, instead of to os.Getpid() of the shared `go test`
// process — which the signal could otherwise terminate outright, racing the
// daemon's own handler (see #326). Run under plain `go test`, with the env
// var unset, it does nothing and passes immediately.
func TestHelperProcess_RealSIGTERM(t *testing.T) {
	if os.Getenv(realSIGTERMHelperEnv) != "1" {
		return
	}

	listener := listenLoopback(t)
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	runDone := make(chan int, 1)
	go func() {
		runDone <- Run(Config{
			ListenAddress: listener.Addr().String(),
			Handler:       http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }),
		}, Deps{
			Listener: listener,
			Logger:   logger,
			// NotifyContext left unset: exercises the real default,
			// signal.NotifyContext(os.Interrupt, syscall.SIGTERM).
			Clock: instantClock{},
		})
	}()

	// Wait for a completed HTTP round-trip, not just an accepted TCP
	// connection: the OS backlog can accept a connection before Run's
	// goroutine ever runs, but a response can only come back once
	// server.Serve is actively looping — which happens strictly after
	// Run registers the real signal.NotifyContext handler earlier in the
	// same goroutine. That ordering is what makes this readiness signal
	// safe to gate the SIGTERM on.
	waitForHealthyResponse(t, listener.Addr().String())

	// A single stdout line the parent scans for; anything else this process
	// writes to stdout/stderr is incidental test-framework output.
	fmt.Println("ready")

	os.Exit(<-runDone)
}

func TestRun_RealSIGTERMShutsDownCleanly(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestHelperProcess_RealSIGTERM$") // #nosec G204 G702 -- os.Args[0] is this same compiled test binary, a fixed self-reexec with a literal flag, not external/tainted input.
	cmd.Env = append(os.Environ(), realSIGTERMHelperEnv+"=1")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("create helper process stdout pipe: %v", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("start helper process: %v", err)
	}
	// Guarantees the helper process is never leaked if this test exits early
	// via a Fatalf below; a no-op once the process has already been waited.
	defer func() { _ = cmd.Process.Kill() }()

	ready := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if scanner.Text() == "ready" {
				ready <- nil
				return
			}
		}
		ready <- fmt.Errorf("helper process ended before signaling ready: %w", scanner.Err())
	}()

	select {
	case err := <-ready:
		if err != nil {
			t.Fatalf("%v; helper stderr: %s", err, stderr.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("helper process never signaled ready; helper stderr: %s", stderr.String())
	}

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("send SIGTERM to helper process: %v", err)
	}

	waitErr := make(chan error, 1)
	go func() { waitErr <- cmd.Wait() }()

	select {
	case err := <-waitErr:
		if err != nil {
			t.Fatalf("expected helper process to exit 0 on SIGTERM shutdown, got: %v; helper stderr: %s", err, stderr.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("helper process did not exit after SIGTERM; helper stderr: %s", stderr.String())
	}
}

func TestRun_ListenFailureReturnsErrorExitCode(t *testing.T) {
	// Occupy the address first, mirroring a real "address already in use"
	// failure, then let Run derive its own listener from cfg.ListenAddress
	// (Deps.Listener unset) exactly like production does.
	occupied := listenLoopback(t)
	defer func() {
		// Best-effort cleanup of the test's own occupying listener; nothing
		// further to do if it fails to close.
		_ = occupied.Close()
	}()

	logger, logs := newTestLogger(t)
	signalCtx, cancelSignal := context.WithCancel(context.Background())
	defer cancelSignal()

	code := Run(Config{
		ListenAddress: occupied.Addr().String(),
		Handler:       http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }),
	}, Deps{
		Logger:        logger,
		NotifyContext: fakeSignalSource(signalCtx, cancelSignal),
		Clock:         instantClock{},
	})

	if code != 1 {
		t.Fatalf("expected exit code 1 on listen failure, got %d; logs: %s", code, logs.String())
	}

	if !strings.Contains(logs.String(), "daemon stopped unexpectedly") {
		t.Fatalf("expected a logged listen failure, got logs: %s", logs.String())
	}
}

func TestRun_ShutdownTimeoutStillExitsCleanly(t *testing.T) {
	// Preserves current behavior: a shutdown that times out logs an error
	// but does not change the process exit code (see daemon.go's Run,
	// mirroring cmd/capsuled's former run() exactly). This test exists to
	// pin that behavior, not to endorse it — see the accompanying report
	// for discussion of whether it should change.
	listener := listenLoopback(t)
	logger, logs := newTestLogger(t)
	handler := newBlockingHandler()

	signalCtx, cancelSignal := context.WithCancel(context.Background())

	runDone := make(chan int, 1)
	go func() {
		runDone <- Run(Config{
			ListenAddress: listener.Addr().String(),
			Handler:       handler,
		}, Deps{
			Listener:      listener,
			Logger:        logger,
			NotifyContext: fakeSignalSource(signalCtx, cancelSignal),
			// instantClock forces the shutdown-timeout branch to fire as
			// soon as shutdown begins, without waiting out a real timeout.
			Clock: instantClock{},
		})
	}()

	// Start a request and hold it open so server.Shutdown cannot complete
	// on its own; only the (instant) timeout can end the shutdown wait.
	reqDone := make(chan error, 1)
	go func() {
		resp, err := http.Get("http://" + listener.Addr().String() + "/")
		if err == nil {
			_ = resp.Body.Close()
		}
		reqDone <- err
	}()

	select {
	case <-handler.received:
	case <-time.After(5 * time.Second):
		t.Fatal("handler never received the in-flight request")
	}

	cancelSignal()

	select {
	case code := <-runDone:
		if code != 0 {
			t.Fatalf("expected exit code 0 even when shutdown times out, got %d; logs: %s", code, logs.String())
		}
	case <-time.After(5 * time.Second):
		close(handler.unblock)
		t.Fatal("Run did not return after shutdown timeout")
	}

	if !strings.Contains(logs.String(), "daemon shutdown failed") {
		t.Fatalf("expected a logged shutdown-timeout error, got logs: %s", logs.String())
	}

	close(handler.unblock)
	<-reqDone
}

func TestRun_ShutdownCompletesBeforeTimeoutWhenNotForced(t *testing.T) {
	// Sanity check for instantClock's counterpart: with a Clock that never
	// fires, a shutdown with no in-flight requests still completes cleanly
	// and quickly, confirming the timeout goroutine only forces things when
	// the clock says to.
	listener := listenLoopback(t)
	logger, logs := newTestLogger(t)

	signalCtx, cancelSignal := context.WithCancel(context.Background())

	runDone := make(chan int, 1)
	go func() {
		runDone <- Run(Config{
			ListenAddress: listener.Addr().String(),
			Handler:       http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }),
		}, Deps{
			Listener:      listener,
			Logger:        logger,
			NotifyContext: fakeSignalSource(signalCtx, cancelSignal),
			Clock:         neverClock{},
		})
	}()

	waitUntilServing(t, listener.Addr().String())
	cancelSignal()

	select {
	case code := <-runDone:
		if code != 0 {
			t.Fatalf("expected exit code 0, got %d; logs: %s", code, logs.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return; shutdown should have completed without needing the clock")
	}
}

// waitUntilServing polls addr until it accepts a connection or the test
// times out, avoiding a fixed sleep before a request or signal that must
// land after Serve has actually started.
func waitUntilServing(t *testing.T, addr string) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("server on %s never started accepting connections", addr)
}

// waitForHealthyResponse polls addr with real HTTP GETs until one completes,
// or the test times out. Unlike waitUntilServing's bare TCP dial, a completed
// HTTP round-trip proves server.Serve's accept loop is actively running, not
// merely that the OS accepted a connection into its listen backlog.
func waitForHealthyResponse(t *testing.T, addr string) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get("http://" + addr + "/")
		if err == nil {
			_ = resp.Body.Close()
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("server on %s never returned a response", addr)
}
