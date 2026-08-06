# Passive authenticated-local-IPC S1 contract result

Date: 2026-08-05

Status: `PASSED` for the exact passive, unwired `RegisterPlanV0` and
`GetRegisteredPlanV0` logical-envelope scope. The full four-method/native authenticated IPC parent
remains `BLOCKED`.

## Defensive scope

This slice defensively validates method separation and copy/no-state behavior using only
repository fixtures, the in-process authority-plane facade, and its fixed transaction oracle. It
creates no XPC listener, Mach service, peer-authenticated connection, key, approval, process,
runtime, backend, VM, guest, installation, or product consumer.

## Frozen passive contract

The two methods remain distinct records and entry points. No opcode or generic command envelope
exists.

| Method | Expected service / role | Purpose | Audience | Method / role-record version | Request / reply application-data caps | Deadline | Response loss |
| --- | --- | --- | --- | --- | ---: | ---: | --- |
| `RegisterPlanV0` | `com.capsulecorp.capsule.supervisor.daemon.v0` / daemon | `capsule.ipc.register-plan.v0` | `capsule.execution-supervisor.local.v0` | `0` / `0` | 328,337 / 4,096 | 5,000 ms | committed retry creates a fresh registration |
| `GetRegisteredPlanV0` | `com.capsulecorp.capsule.supervisor.broker.v0` / Broker | `capsule.ipc.get-registered-plan.v0` | `capsule.execution-supervisor.local.v0` | `0` / `0` | 16 / 332,433 | 2,000 ms | repeatable read by `RegistrationID` |

Every passive request binds protocol version `0`, one nonzero correlation-only request ID, the
current installation ID, and the current epoch sequence/digest. The method record supplies service,
role, audience, and purpose; application data cannot override them. Replies echo the request ID and
return defensive copies only.

`RegisterPlanV0` carries the existing exact plan, 562-byte plan-v0 role projection, canonical
single-member `SourceManifest`, and exact `main.mjs` bytes. `GetRegisteredPlanV0` accepts only a
16-byte nonzero Supervisor-issued `RegistrationID` and returns the exact retained plan, resolved
562-byte role projection, registration, manifest, and source. The method version and 562-byte role
record version are checked together. No 626-byte record, TypeScript role, optional transform,
plan-v1 equivalence, or dual-active authority exists.

The caps are generated from closed field maxima:

- register request: `65,536 + 562 + 95 + 262,144 = 328,337`;
- register reply: `4,096`;
- fetch request: `16`; and
- fetch reply: `65,536 + 562 + 4,096 + 95 + 262,144 = 332,433`.

Exact-boundary and cap-plus-one cases are retained as deterministic logical field-length fixtures.
They make no raw Mach-message or XPC-serialization size claim.

## Oracles and evidence

[`internal/execution/authorityplane/passive_ipc.go`](../internal/execution/authorityplane/passive_ipc.go)
contains the passive method records, logical request/reply projections, fixed refusal projection,
defensive copying, current installation/epoch checks, and an unexported conformance-only adapter for
the two method-specific facade calls.

[`schemas/conformance/authenticated-local-ipc-v0/manifest.json`](../schemas/conformance/authenticated-local-ipc-v0/manifest.json)
retains exact JSON known answers, body-byte identities, eight maximum/cap-plus-one cases, twenty-five
fixed refusal cases, copy-ownership rules, complete pre-core-refusal zero-state/effect assertions, and both
response-loss dispositions. Exact JSON is fixture serialization only, never proposed XPC framing.
Independent Go and Node verifiers agree on the method records, request/reply projections, body
digests, caps, versions, refusals, and response-loss rules.

Wrong method/record version, method, service, role, audience, purpose, protocol version, request
ID, installation, epoch, or registration ID refuses before an authority-core call where that
boundary owns the check. Those refusals create no registration, consume no approval, create no
attempt, call no lifecycle/backend, and create no endpoint/process/runtime/guest. Caller, accessor,
and reply mutation cannot alter copied or retained bytes. A lost registration reply followed by a
retry retains two fresh registrations; a lost fetch reply remains a byte-exact repeatable read.
Oversize reply fixtures are post-core local-integrity faults, not pre-core refusals; they make no
no-state claim because registration may already have committed before an invalid reply is observed.
The passive Go boundary panics for this impossible output violation so it cannot become a soft
refusal; the later native slice owns actual whole-Supervisor termination.

## Deliberately unresolved native boundary

ADR-0029 requires a later native XPC dictionary and numeric message/status mapping, but it does not
yet assign exact ASCII key spellings, numeric message tags, or numeric refusal codes. This slice
does not invent them. It records `transportEncoding`, `numericMessageTags`, and
`peerAuthenticationEvidence` as `null` and therefore supports no endpoint or peer-authentication
claim.

Next native slice remains `BLOCKED` until a reviewed transport contract fixes those values and an
explicitly authorized local harness can test authentication-before-body-copy, message-derived code
identity, session checks, flow control, interruption, and deadline behavior. Installed/product IPC
also remains blocked on protected state, identity/update evidence, signing/key authorization,
Broker consumer work, and the independent runtime/backend/content/evidence gates.

## Successor prerequisite

The later [S3 native-contract prerequisite](AUTHENTICATED_LOCAL_IPC_S3_NATIVE_CONTRACT.md) is now
`PASSED` for exact keys, types, versions, numeric method/status/reason tags, envelope key counts,
and passive refusal fixtures for the three frozen methods. This does not revise S1's historical
scope or supply the still-missing local native harness, listener, or peer-authentication evidence.
