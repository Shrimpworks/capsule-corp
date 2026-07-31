# ADR-0009: Resource policy is user-owned and exact

- Status: Accepted
- Date: 2026-07-30
- Refined by: ADR-0011

## Context

Resource limits are part of the authority a user approves. Silent clamping or backend substitution
can make the displayed request differ from the execution and can hide policy errors. Allowing an
agent to choose unbounded values would expose host availability.

Capsule needs user control without giving the agent or backend permission to expand or reinterpret
that control.

## Decision

Trusted user policy defines explicit default values and hard ceilings. A proposal may request values
at or below the ceilings. Missing values resolve to user defaults while preparing the execution
plan. Requests above a ceiling are rejected rather than silently clamped.

The execution plan records the exact resolved limits shown to and approved by the user. After
approval, the daemon does not modify them. An isolation backend must either enforce every value or
refuse the execution.

v0 includes limits for request, source, inline input, snapshot input, wall time, CPU time, memory,
PIDs, scratch storage, logs, artifact count, per-artifact bytes, and total artifact bytes. Profile
metadata may declare supported ranges, but it does not override stricter user policy.

## Consequences

- User intent remains visible and auditable.
- Policy errors fail closed instead of producing surprising execution.
- Plans and receipts can compare requested, defaulted, and enforced values exactly.
- Backends need explicit capability discovery and enforcement verification.
- A job that could run under a smaller value is still rejected if it requested an unauthorized
  larger value; the user or proposer must submit a new plan.
- Tests must prove that no value changes after approval and that unsupported backend controls reject
  execution.
