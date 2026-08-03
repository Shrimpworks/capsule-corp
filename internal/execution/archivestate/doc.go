// Package archivestate defines passive, runtime-neutral projections for the
// proposed Supervisor archive boundary. It validates and hashes immutable
// values and selects closed cohorts in memory. It performs no filesystem,
// store, lifecycle, backend, IPC, process, or guest operation.
package archivestate
