// Command capsuled is the capsule daemon entrypoint. It currently only
// serves internal/api's read-only diagnostic endpoints on a loopback
// address by default; it holds no Approval, Supervisor, or execution
// authority (see internal/api's package doc and AGENTS.md).
package main

import (
	"flag"
	"os"

	"capsule.local/capsule/internal/api"
	"capsule.local/capsule/internal/daemon"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	os.Exit(run())
}

// run parses flags, builds the daemon configuration, and hands off to
// internal/daemon.Run for the actual lifecycle (listen, serve, shut down on
// signal, pick an exit code). It is deliberately thin: internal/daemon.Run
// is what carries the logic under test, since flag.String against the
// global flag.CommandLine can only run once per process and so cannot be
// exercised repeatedly in-process here.
func run() int {
	listenAddress := flag.String("listen", "127.0.0.1:7777", "daemon listen address")
	flag.Parse()

	handler := api.NewServer(api.Options{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	})

	return daemon.Run(daemon.Config{
		ListenAddress:   *listenAddress,
		Handler:         handler,
		ShutdownTimeout: daemon.DefaultShutdownTimeout,
		Version:         version,
	}, daemon.Deps{})
}
