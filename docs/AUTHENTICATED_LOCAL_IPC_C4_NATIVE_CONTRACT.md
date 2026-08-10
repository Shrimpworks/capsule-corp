# Passive authenticated-local-IPC C4 candidate and CL4 audit disposition

Date: 2026-08-10

Work item: freeze `SubmitApprovalV0` and `RequestAttemptV0` in the existing passive
`xpc-dictionary-v0` contract.

Status: `BLOCKED` for the current passive evidence claim pending the focused evidence-hardening
implementation required by the completed CL4 audit.

CL4 audit: `PASSED`; disposition: `AMEND`.

Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`.

Installed authenticated IPC, live approval signing, protected-state consumers, runtime/profile
admission, and product admission: `BLOCKED`.

## Defensive scope

This slice uses only Capsule repository constants, generators, generated JSON fixtures, the passive
Go contract model, Go tests, and an independent Node verifier. It adds no listener, Mach service,
service registration, peer authentication, entitlement, identity, key, `LAContext`, signer,
Keychain item, store consumer, threshold-policy hook, process, endpoint, runtime, backend, VM, or
guest. APL-2 remains the applicable reuse-map row: a later product slice may adopt Apple XPC while
Capsule owns the closed roles, dictionaries, bounds, and replay policy. No dependency is added.

## Frozen method and state tags

The existing tags are unchanged. `0` remains invalid.

| Method | Tag | Service | Derived role | Local purpose | Deadline |
| --- | ---: | --- | --- | --- | ---: |
| `SubmitMainMJSV0` | 1 | daemon CLI service | internal-alpha CLI | `capsule.ipc.submit-main-mjs.v0` | 10,000 ms |
| `RegisterPlanV0` | 2 | Supervisor daemon service | daemon | `capsule.ipc.register-plan.v0` | 5,000 ms |
| `GetRegisteredPlanV0` | 3 | Supervisor Broker service | Broker | `capsule.ipc.get-registered-plan.v0` | 2,000 ms |
| `SubmitApprovalV0` | 4 | Supervisor Broker service | Broker | `capsule.ipc.submit-approval.v0` | 5,000 ms |
| `RequestAttemptV0` | 5 | Supervisor daemon service | daemon | `capsule.ipc.request-attempt.v0` | 5,000 ms |

The new reply state tags are closed. Approval state is `invalid=0`, `usable=1`, `consumed=2`, and
`invalidated=3`. Attempt state is `invalid=0` and `created=1`. Success replies never use either
invalid tag. The five-second values are no longer prose arithmetic: the merged generator emits
admission-start deadline cases. CL4 found that these do not yet establish explicit before, exactly
at, and after boundary behavior and that the Go and Node comparison must be strengthened. These are
passive protocol deadlines, not observed platform latency or installed-service evidence.

## Exact dictionaries and mechanically derived caps

All requests retain the exact nine common keys and all replies retain the exact six common keys
from the S3 contract. Values remain top-level `XPC_TYPE_UINT64`, `XPC_TYPE_DATA`, or
`XPC_TYPE_STRING`; nested containers, file descriptors, endpoints, Mach rights, and extra objects
remain forbidden.

`SubmitApprovalV0` has exactly 11 request keys: the nine common keys, a nonzero 16-byte
`capsule.registration-id`, and `capsule.approval-envelope` containing 1..512 exact bytes. Its
success reply has exactly eight keys: the six common reply keys, a nonzero 16-byte
`capsule.approval-id`, and `capsule.approval-state` as one allowed numeric tag. Its application-data
caps are therefore exactly `16 + 512 = 528` request bytes and `16` reply bytes.

`RequestAttemptV0` has exactly 11 request keys: the nine common keys, a nonzero 16-byte
`capsule.registration-id`, and a nonzero 16-byte `capsule.approval-id`. Its success reply has
exactly eight keys: the six common reply keys, a nonzero 16-byte `capsule.attempt-id`, and
`capsule.attempt-state == created`. Its application-data caps are exactly `16 + 16 = 32` request
bytes and `16` reply bytes. XPC framing and numeric-value storage are not counted as application
data and no raw Mach-message cap is claimed.

The approval envelope remains the exact bounded tagged canonical COSE_Sign1 object with an
embedded canonical `capsule.approval-grant` v0 payload. Its signed audience is
`capsule.execution-supervisor` and its signed purpose is `capsule.plan.approve`. The local channel
audience remains `capsule.execution-supervisor.local.v0`; local and signed purpose/audience values
are never interchangeable. `RequestAttemptV0` transports only the two typed lookup identities and
accepts no plan, approval bytes, backend setting, image, mount, path, or threshold-policy input.

## Dispatch, copy, and refusal precedence

Dispatch identity is derived from exact receiving service plus the explicitly recorded closed
entry point plus admitted peer role. Each entry point has the same spelling as its method name;
the numeric tag is checked afterward and cannot select another entry point. The generated ordered
matrix now contains all 20 foreign-tag collisions across the five methods. No generic dispatcher,
opcode bus, command field, service field, or role field exists.

The frozen native-owner contract requires peer/session/flow, exact outer key/type/width, current
installation and epoch, and method binding to complete before copying either new body. It then
requires a separate copy-only bridge per method: native-owned bytes are copied before entry into
Go, Go-owned inputs cannot alias native memory, and native fixed-cap output buffers receive no Go
pointer. Output cap, short-write, pointer/length, or bridge-version failure terminates without an
application reply. This is a passive bridge requirement, not an activated bridge implementation.

Validation precedence remains exact:

1. protocol version;
2. method version;
3. service + entry point + role + message tag + local audience + local purpose;
4. installation/epoch current state;
5. bounded application-data copy; and
6. embedded record version and core validation.

Joint protocol/method/ApprovalGrant-version mutations therefore select protocol first, then method,
then the copied core refusal for ApprovalGrant object version. Missing/extra keys, wrong or nested
types, cap-plus-one data, zero identifiers, wrong current state, wrong role/session, signed/local
purpose swaps, same-width ID substitutions, and keyed FD/endpoint/Mach-send-right smuggling all
retain exact status/reason/body-copy/core-call outcomes in the ordered case table. The S3 status
and reason numbers are unchanged. Recognized core failures keep their exact status and use
`coreRefusal`; local integrity faults have no reply.

Same-width bytes do not carry hidden authority. A request ID remains correlation only. A
registration lookup, approval reference, and attempt reference remain separate domains because the
method and field key select their typed Go input; substituting bytes from another 16-byte domain
does not grant that domain and reaches only the exact current-state/core binding refusal.

## Replay, cancellation, response loss, and flow

- `SubmitApprovalV0` replay identity is the verified canonical approval payload plus resolved
  signer-authorization identity. An exact or mathematically equivalent envelope returns the same
  `ApprovalID` and current state with zero duplicate authority effect.
- `RequestAttemptV0` replay identity is `(RegistrationID, ApprovalReference)`. Exact replay or
  concurrency returns the same committed `AttemptID` and current state with zero duplicate attempt
  or lifecycle effect.
- Request IDs are not idempotency keys. Cancellation before dispatch makes no core call and creates
  no reply or state effect. Cancellation or deadline after dispatch cancels only response delivery;
  store truth or the recovery fence controls the result.
- The generated response-loss table retains the original three entries and appends complete exact
  entries for both new methods. Go and Node deep-compare it.
- A generated isolated aggregate case admits one maximum 528-byte approval request and one
  32-byte attempt request, refuses byte 561 as `CAPACITY`, releases the 32-byte attempt slot, and
  re-admits exactly 32 bytes. This is accounting evidence; it does not claim the production
  2,626,696-byte Supervisor cap is independently reachable.

## Retained evidence

- `internal/execution/authorityplane/passive_native_xpc_contract.go` is the passive production-side
  specification. It imports no XPC API and exposes no dispatcher or endpoint.
- `schemas/conformance/authenticated-local-ipc-v0/native-xpc-v0.contract.json` contains all five
  method bindings/envelopes, numeric tables, the complete ordered native cases, deadline cases,
  response-loss table, and no-effect/platform-evidence nulls.
- `schemas/conformance/authenticated-local-ipc-v0/oracles.json` retains the exact/cap-plus-one,
  copy, aggregate release/re-admission, cancellation, and response-loss oracles.
- `scripts/generate-authenticated-local-ipc-conformance.mjs` derives the aggregates and generated
  references. `scripts/verify-authenticated-local-ipc-conformance.mjs` reconstructs the ordered
  case table and new method tables, but CL4 requires a stronger independent complete-map and
  all-field comparison before the evidence claim can pass.
- `passive_native_xpc_contract_test.go` compares the Go method/envelope/state/deadline/response-loss
  model and ordered native cases, but CL4 requires complete dictionary, closed-map, required
  `noState`, cancellation/deadline, and refusal-table completeness checks.

Focused verification:

```sh
node scripts/generate-authenticated-local-ipc-conformance.mjs --check
node scripts/verify-authenticated-local-ipc-conformance.mjs
go test ./internal/execution/authorityplane
```

## Limitations and ADR impact

This implements Accepted ADR-0029 and the S0 review without changing component responsibility,
service topology, signed-object policy, or ADR lifecycle. No ADR addendum is required. CL4 found no
runtime authority bypass, but its `AMEND` disposition blocks the current C4 passive evidence claim.
A separate focused evidence-hardening PR must:

1. generate explicit before, exactly at, and after 5,000-ms deadline cases for
   `SubmitApprovalV0` and `RequestAttemptV0`, with exact-at-boundary behavior defined;
2. strengthen the Go and Node verification to independently compare complete dictionaries, closed
   maps, every case field, required `noState` entries, cancellation/deadline oracles, and
   refusal-table completeness; and
3. rerun generator `--check`, the independent Node verifier, focused Go contract/replay tests, and
   `git diff --check`.

This documentation integration does not implement that follow-up and does not mark it `PASSED`.
Native OS pre-delivery enforcement, installed identities/session/update, Broker live signing,
protected-state consumers, threshold-policy wiring, runtime/profile admission, and product
admission remain separately `BLOCKED`.
