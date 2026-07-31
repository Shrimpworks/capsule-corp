# ADR-0011: Execute only Supervisor-registered immutable plans

- Status: Accepted
- Date: 2026-07-31
- Refines: ADR-0006, ADR-0008, and ADR-0009

## Context

Approving a digest is insufficient if the execution component later receives replacement bytes or
interprets policy differently. A daemon-owned grant ledger also lets a daemon compromise replay or
reset authorization.

## Decision

The daemon sends exact canonical `ExecutionPlan` bytes to the Supervisor before approval. The
Supervisor independently enforces raw/schema limits and versioned hard-safety rules, stores exact
bytes durably, and returns a `PlanRegistration`.

The Broker fetches those registered bytes directly from the Supervisor and renders their typed
fields. Its `ApprovalGrant` binds plan digest, registration, installation, trust epoch, expected
Supervisor, purpose, audience, attempt nonce, and expiry.

Attempt APIs accept only the registration ID. In one durable Supervisor transaction, the grant is
verified and consumed and one `ExecutionAttempt` is created before any backend side effect. A crash
after consumption does not restore the grant.

## Consequences

- The daemon cannot substitute plan bytes at execute time.
- Plan A approval cannot authorize plan B or a second attempt.
- Supervisor validation must be independently implemented and versioned.
- Grant-ledger and attempt recovery become Supervisor responsibilities.
- Approval can be burned by a post-consumption failure; safe retry requires new user approval.
- Exact backend capability matching happens before attempt launch; unsupported controls refuse the
  attempt rather than being clamped or approximated.
