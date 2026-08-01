// Package execution contains the portable execution-control scaffold.
//
// SupervisorCore is currently an OS-neutral executable contract for plan
// registration, approval consumption, trust-transition fencing, and no-guest
// lifecycle recovery. MemoryStateStore is non-durable, DevelopmentLifecycle is
// test-only, and no implementation here is a macOS authority or hostile-code
// boundary. Product wiring requires the approved protocol, authenticated native
// IPC, a durable transactional store, production cryptographic verification,
// and a separately validated backend adapter.
package execution
