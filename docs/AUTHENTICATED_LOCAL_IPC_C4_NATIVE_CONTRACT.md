# Passive authenticated-local-IPC C4 contract and CL4 amendment closure

Date: 2026-08-10

Work item: freeze `SubmitApprovalV0` and `RequestAttemptV0` in the existing passive
`xpc-dictionary-v0` contract.

Status: `PASSED` for the exact passive, unwired, no-listener C4 evidence scope.

CL4 audit: `PASSED`; historical disposition: `AMEND`; focused follow-up: closed.

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

PR #248 is the canonical predecessor that retained the CL4 `AMEND` disposition and dispatched this
focused implementation closure. This follow-up changes evidence strength only; it does not alter
topology, authority ownership, replay identity, privileged responsibility, or ADR lifecycle.

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
invalid tag. The five-second deadline check is closed and mechanical: dispatch is allowed only
when elapsed admission time is strictly less than 5,000 milliseconds. Equality is expired. For
each new method, the generated ordered cases are:

| Boundary | Case identity | Decision | Reply/status/reason | Core/state |
| --- | --- | --- | --- | --- |
| 4,999 ms | `<Method>.deadline.immediately-before` | dispatch the core | the store-semantic reply may be delivered unless the deadline later cancels delivery; transport status/reason are `null` | body copied, one core call, store semantic result or recovery fence controls |
| 5,000 ms | `<Method>.deadline.exactly-at` | expired before dispatch | no application reply; status/reason are `null` | no body copy, zero core calls, exact zero `noState` |
| 5,001 ms | `<Method>.deadline.immediately-after` | expired before dispatch | no application reply; status/reason are `null` | no body copy, zero core calls, exact zero `noState` |

`<Method>` is separately `SubmitApprovalV0` and `RequestAttemptV0`, producing six cases in that
order. These are passive protocol deadlines, not observed platform latency, OS cancellation, or
installed-service evidence.

## Exact dictionaries and mechanically derived caps

All requests retain the exact nine common keys and all replies retain the exact six common keys
from the S3 contract. Values remain top-level `XPC_TYPE_UINT64`, `XPC_TYPE_DATA`, or
`XPC_TYPE_STRING`; nested containers, file descriptors, endpoints, Mach rights, and extra objects
remain forbidden. Every envelope now records `requiredKeyCount == exactKeyCount`,
`optionalKeyCount == 0`, and `closedMap == true`; every listed field independently records
`required == true`.

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
- The exact after-dispatch commit oracles prohibit a deadline from turning an already committed
  approval or attempt replay into a new refusal. Approval replay returns the same `ApprovalID` and
  current state; attempt replay returns the same `AttemptID` and current state.
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
  references. `scripts/verify-authenticated-local-ipc-conformance.mjs` independently compares the
  complete closed method dictionaries and maps, every ordered case field, all 20 foreign-tag
  collisions, the complete refusal and response-loss tables, required `noState` entries, deadline,
  cancellation, release/re-admission, cap-plus-one, and no-reply oracles.
- Every ordered native case now carries exact status/reason, body-copy/core-call, reply disposition,
  process-termination disposition, `noState`, and post-core-state fields; `null` is an explicit
  not-applicable value rather than an omitted or silently unchecked field.
- `passive_native_xpc_contract_test.go` independently compares all five complete request/success/
  refusal dictionaries, the full ordered cases, status/reason/state tables, six deadline cases,
  exact five-entry response-loss table, and complete oracle sections.
- Bounded Go and Node mutations prove refusal of a missing or extra dictionary field, changed
  type/width/cap, absent required `noState`, changed cancellation/deadline result, missing refusal
  case, response-loss drift, and exact-at-boundary inversion.

Focused verification:

```sh
node scripts/generate-authenticated-local-ipc-conformance.mjs --check
node scripts/verify-authenticated-local-ipc-conformance.mjs
go test ./internal/execution/authorityplane
```

## Limitations and ADR impact

This implements Accepted ADR-0029 and the S0 review without changing component responsibility,
service topology, signed-object policy, or ADR lifecycle. No ADR addendum is required. CL4 found no
runtime authority bypass; its historical `AMEND` disposition is now closed by the focused evidence
hardening above, so the exact passive/no-listener C4 claim is `PASSED`.

Native OS pre-delivery enforcement, installed identities/session/update, Broker live signing,
protected-state consumers, threshold-policy wiring, runtime/profile admission, and product
admission remain separately `BLOCKED`.
