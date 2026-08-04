# Changelog

All notable changes to this project will be documented in this file.

The project currently follows the structure of [Keep a Changelog](https://keepachangelog.com/),
but it has not selected a public versioning or release policy yet.

## [Unreleased]

- Compose the passive Darwin installation owner with an owner-required current-v1/no-guest
  Supervisor startup: one session binds store and coordinator, sorted `AttemptID` recovery occurs
  only after store validation, post-open ownership loss fences until reopen, and shutdown closes the
  store before the owner descriptor. This remains an unwired local mechanic under owned temporary
  roots and does not add v2/archive, IPC, runtime, real-backend, guest, or installed-storage claims.

### Added

- Initial repository scaffold.
- Canonical project, architecture, security, and roadmap documentation.
- Draft job, runtime profile, and execution receipt schemas.
- Minimal Go daemon and TypeScript protocol, SDK, and MCP package skeletons.
- Reproducible Node toolchain declaration and unified local CI command.
- GitHub publication checklist, support policy, ADR template, and Go control-plane decision.
- Repo-owned AI Central link setup for the full reviewed local skill catalog, including Caveman and
  Hallmark, while keeping symlinks out of Git.
- Passive Phase 2A `JobProposal`, minimum `ExecutionPlan`, and `PlanRegistration` candidates with
  byte-exact fixtures, Go/TypeScript decoded views, and explicit non-activation boundaries.
- Retained Gate C libkrun/HVF storage, console/lifecycle, installed-recovery, adversarial, and
  supply-chain experiments with a reconciled P0 plan and explicit prohibited claims.
- Proposed Phase 2B decoder/registration bounds plus a manifest-driven conformance foundation with
  37 rules, 105 cases, and 91 unique raw/media/scalar/CBOR fixtures.
- A durable workstream/evidence ledger mapping parallel task conclusions and external reviews to
  retained evidence, synthesis documents, merged integration commits, and remaining work.
- An unwired exact plan-registration conformance handoff that passes copied TypeScript plan bytes
  and complete role bindings to the Go `registrationstate` component without defining product IPC
  or activating any daemon, approval, backend, or guest surface.
- A passive internal Go/Darwin Supervisor owner capability that validates one enrolled pre-created
  sibling object, acquires lifetime nonblocking BSD `flock`, retains an opaque `CLOEXEC`
  descriptor, and refuses duplicate or mismatched local owners before downstream test markers.
  Product startup, owner-required store opening, protected storage, and installed evidence remain
  unimplemented.

### Removed

- Moved 38 completed non-production experiment trees (744 tracked files) to the public,
  commit-pinned `Shrimpworks/capsule-experiments` archive. Capsule retains canonical decisions,
  production conformance fixtures, and immutable evidence links.
- The legacy `internal/execution.SupervisorCore` in-memory scaffold (`SupervisorCore`,
  `MemoryStateStore`, `DevelopmentLifecycle`, and the deprecated `Backend`/`RuntimeAdapter`
  interfaces). It had no product wiring and was already superseded as the authoritative unwired
  path by the `approvalattempt`/`registrationstate`/`registeredlifecycle` split. See ADR-0027.
