# ADR-0046: Freeze the passive typed guest transport

- Status: Accepted
- Date: 2026-08-10
- Refines: ADR-0015, ADR-0017, ADR-0022, ADR-0040, ADR-0041, and ADR-0042

## Context

R2 reconciled C2A's narrowed first-release caps, C2B v4's fixed descriptor/port topology, the v27
diagnostic-console limitations, and the passive completion-last transaction. The historical P0-3
candidate established useful 152/160/64-byte framing shapes but used an obsolete 1 MiB source cap
and was retained outside this repository. C5a must freeze a repository-owned successor before any
separately authorized real transport experiment can run.

Without a closed successor, implementations could disagree on offsets, status values, canonical
JSON, failure precedence, commit placement, cap-plus-one handling, cancellation, descriptor
ownership, or what a guest frame proves. Treating EOF, exit zero, or a self-consistent guest frame
as durable completion would violate ADR-0015 and ADR-0042.

## Decision

Capsule accepts passive typed guest transport v1 as the C5a contract. It retains three fresh,
role-distinct, attempt-bound, one-frame streams created by the Execution Supervisor. Source and
input have 152-byte headers and at most 262,144 payload bytes. Completion has a 160-byte header, at
most 262,144 typed canonical-JSON bytes, and one 64-byte exact-offset trailer written last. All
integers are unsigned big-endian. The exact layouts, magic, roles, bindings, caps, canonical JSON,
and refusal precedence are frozen in the linked protocol contract and generated fixtures.

The sole source payload is the exact registered UTF-8 `main.mjs` bytes with no container or path.
Input and completion share `capsule.canonical-inline-json/v0`. Completion statuses are exactly
`succeeded`, `workload-failed`, `result-invalid`, and `child-terminated`; every non-success payload
is exact JSON `null`, so guest prose, paths, raw errors, and timing do not enter the frame.

The Supervisor owns endpoint creation, directional writers/drains, deadlines, cancellation,
teardown, authoritative absence, and durable publication. The launcher validates both inputs before
child creation, owns the completion writer only through final trailer closure, and never gives the
workload that endpoint or a reopenable node. Cap-plus-one is fully counted and drained without
partial authority effects. Reset/cancel/fault state is monotonic.

A valid frame establishes only `FRAME_OBSERVED`. Durable completion additionally requires
independent authority, input, runtime/profile, lifecycle, teardown, child-tree and runner absence,
cleanup, and publication facts. Only the Supervisor-owned completion transaction publishes
`committed-last`; EOF and process exit never commit.

The passive Go package and Node verifier independently reproduce 48 generated frame dispositions,
13 state/fault cases, and 23 restoration cases. They create no endpoint or product consumer. This
decision adds no dependency and no new Supervisor responsibility: it makes the already-selected
R2/C2A/C2B/ADR-0042 responsibilities byte-exact.

## Consequences

- C5a is `PASSED` in its exact passive contract/conformance scope.
- The obsolete 1 MiB source cap and any attempt to reuse v27 diagnostic console bytes as typed
  completion are rejected.
- Any change to magic, version, offsets, caps, status vocabulary, canonical JSON, endpoint custody,
  refusal precedence, or terminal-proof join requires a versioned successor and ADR review.
- C5b controlled execution remains `BLOCKED` on separate authorization naming the exact governed
  runtime/profile successor, owned disposable guest/host, fixture digest, process names, and
  evidence destination.
- Installed composition/recovery, real transport and process-tree evidence, runtime/profile
  admission, sealed adapter wiring, hostile owner source, and product admission remain `BLOCKED`.
