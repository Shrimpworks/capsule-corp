// Package registeredlifecycle implements the unwired Phase 2B created-attempt
// fake lifecycle. Its drive and recovery entries accept only a
// Supervisor-issued attempt ID resolved from the durable Slice B store, and
// its concrete backend creates no guest. The package is intentionally separate
// from execution.SupervisorCore and from all experiment code.
package registeredlifecycle
