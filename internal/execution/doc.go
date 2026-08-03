// Package execution contains the portable execution-control scaffold.
//
// SupervisorCore is currently an OS-neutral executable contract for plan
// registration, approval consumption, trust-transition fencing, and no-guest
// lifecycle recovery. MemoryStateStore is non-durable, DevelopmentLifecycle is
// test-only, and no implementation here is a macOS authority or hostile-code
// boundary. Product wiring requires the approved protocol, authenticated native
// IPC, a durable transactional store, production cryptographic verification,
// and a separately validated backend adapter.
//
// SupervisorCore is the older in-memory scaffold and is not the oracle for
// ADR-0024: see docs/EXECUTION_SUPERVISOR.md. The current unwired path is the
// registrationstate/approvalattempt/registeredlifecycle split in this tree's
// subpackages; SupervisorCore remains behaviorally frozen alongside it until
// retired.
package execution
