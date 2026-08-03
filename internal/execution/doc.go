// Package execution is a documentation-only root for the portable
// execution-control scaffold. It holds no executable code of its own.
//
// The former SupervisorCore in-memory scaffold (plan registration, approval
// consumption, trust-transition fencing, and no-guest lifecycle recovery) was
// retired: it had no product wiring (unreferenced by cmd/capsuled and
// internal/api) and was already superseded as the authoritative unwired path
// by the subpackages below. See docs/EXECUTION_SUPERVISOR.md for current
// implementation status.
//
// The authoritative unwired path is:
//
//   - internal/execution/approvalattempt: passive typed domains, fixed
//     classifications/states, and a retained-vector-only bounded verifier;
//   - internal/execution/registrationstate: a fixed file snapshot colocating
//     exact registrations, effective-time high-water, approvals, and
//     immutable created attempts, including atomic consume/create and reopen
//     validation; and
//   - internal/execution/registeredlifecycle: an AttemptID-keyed, no-guest
//     fake lifecycle that revalidates the committed attempt and exact
//     registered plan before fake prepare.
//   - internal/execution/lifecyclestate: durable attempt lifecycle types
//     colocated with Supervisor authority state (ADR-0025).
//
// These are implemented local mechanics with retained Go tests, not a
// deployed Supervisor. Product wiring requires the approved protocol,
// authenticated native IPC, a durable transactional store, production
// cryptographic verification, and a separately validated backend adapter.
package execution
