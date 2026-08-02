// Package registeredlifecycle implements the unwired Phase 3 registered-plan
// fake lifecycle. Its public execution and recovery entries accept only a
// Supervisor-issued registration ID, and its concrete backend creates no
// guest. The package is intentionally separate from execution.SupervisorCore
// and from all experiment code.
package registeredlifecycle
