# ADR-0027: Retire the SupervisorCore in-memory scaffold

- Status: Accepted
- Date: 2026-08-02
- Refines if accepted: none; records completion of a transition already directed by ADR-0024 and
  ADR-0025

## Context

`internal/execution.SupervisorCore` was the original build scaffold for plan registration, approval
consumption, trust-transition fencing, and no-guest lifecycle recovery, backed by the non-durable
`MemoryStateStore` and test-only `DevelopmentLifecycle`. `internal/execution/contracts.go` also
carried the `Backend`/`RuntimeAdapter`/`PreparedJob` interfaces, already marked `Deprecated` in
their own doc comments as unable to satisfy Capsule's durable prepare/create/stage/start/collect/
destroy/reconcile contract.

ADR-0024 and ADR-0025 established a second, colocated implementation
(`internal/execution/approvalattempt`, `registrationstate`, `registeredlifecycle`,
`lifecyclestate`) as the authoritative unwired path going forward. `docs/EXECUTION_SUPERVISOR.md`
and `docs/PHASE_2A_PARALLEL_REVIEW_SYNTHESIS.md` already documented `SupervisorCore` as "not the
oracle for ADR-0024" and behaviorally frozen. It remained in the tree only as a second, unmaintained
implementation of overlapping responsibility, which a codebase audit flagged as a source of
confusing drift (for example, its `RequestAttempt` never re-validated runtime-integrity staleness
inside its commit closure the way `ApprovalAttemptComponent.RequestAttempt` does).

A repo-wide search confirmed zero external consumers: no file outside `internal/execution` itself
imports the top-level `capsule.local/capsule/internal/execution` package, and `cmd/capsuled` /
`internal/api` reference neither `SupervisorCore` nor any of its supporting types
(`OpaqueID`, `InstallationState`, `AttemptRecord`, `PreparedJob`, `Backend`, `RuntimeAdapter`, etc.).
All of it was self-contained within `contracts.go`, `development_lifecycle.go`, `memory_store.go`,
`supervisor_state.go`, and their test file.

## Decision

Delete `internal/execution/{contracts,development_lifecycle,memory_store,supervisor_state,
supervisor_state_test}.go` outright rather than deprecating them in place. `internal/execution`
becomes a documentation-only package root (`doc.go`) pointing to the `approvalattempt`/
`registrationstate`/`registeredlifecycle`/`lifecyclestate` split as the sole current implementation.

This is a removal of dead, superseded scaffold, not a new Supervisor responsibility, authority
change, or security posture change: nothing here was ever wired to a real backend, endpoint, or
guest, and the replacement path was already accepted and further along.

## Consequences

- One Supervisor-shaped implementation exists in the tree instead of two, removing the maintenance
  hazard of two independently drifting copies of the same invariants (the staleness-recheck gap
  above cannot recur because there is only one `RequestAttempt` now).
- `docs/EXECUTION_SUPERVISOR.md`'s living status section and this repository's `CHANGELOG.md` record
  the removal; historical checkpoint documents (`PHASE_2A_PARALLEL_REVIEW_SYNTHESIS.md`,
  `PHASE_2B_APPROVAL_ATTEMPT_BOUNDARY.md`, etc.) are left unchanged as dated snapshots of what was
  true when written.
- No test, build, or product surface depended on the removed code; `go build ./...`, `go vet ./...`,
  and `go test ./...` are unaffected.
- Nothing about the `approvalattempt`/`registrationstate`/`registeredlifecycle` boundary, its
  unwired status, or ADR-0024/ADR-0025's open questions changes.
