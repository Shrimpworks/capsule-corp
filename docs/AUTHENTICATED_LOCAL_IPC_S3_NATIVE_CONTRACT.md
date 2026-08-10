# Passive authenticated-local-IPC S3 native-contract prerequisite

Date: 2026-08-05

Conformance amendment: 2026-08-06

C4 successor: 2026-08-10. This document preserves the original three-method S3 scope. The current
generated `native-xpc-v0` contract is extended, without renumbering these methods, by the
[C4 passive contract](AUTHENTICATED_LOCAL_IPC_C4_NATIVE_CONTRACT.md).

Work item: exact native XPC dictionary contract for `SubmitMainMJSV0`, `RegisterPlanV0`, and
`GetRegisteredPlanV0`.

Status: `PASSED` for the exact passive, unwired encoding-contract and cross-language conformance
scope. The S3 native authentication/cap harness remains `BLOCKED` on a separately authorized local
harness; installed authenticated product IPC remains `BLOCKED`.

## Defensive scope

This prerequisite defensively freezes the dictionary field-and-type boundary that a later local
S3 harness must consume. It uses only repository Go constants/specifications, generated JSON
fixtures, and independent Go/Node verification. It creates no XPC listener, Mach service, service
registration, connection, Apple identity/profile/credential, key, installation, protected store,
daemon or Broker consumer, process, runtime, backend, VM, or guest. It supplies no peer-
authentication or native-delivery evidence.

No new dependency or custom platform primitive is introduced. The applicable ecosystem row is
`APL-2`: adopt Apple XPC as the later platform transport while Capsule builds the narrow exact
roles, dictionaries, caps, and replay policy. This slice freezes only that Capsule-owned policy and
does not call the platform.

An external experiment or new ADR was not required for this prerequisite. Accepted ADR-0029
already selects the two role-specific Supervisor services, authentication-before-body-copy
ordering, fixed common fields, copy ownership, refusal classes, and method-specific bridge
topology; Accepted ADR-0044 selects the single CLI-to-daemon method and its flow contract. The S3
plan expressly calls for the exact native dictionary and numeric mapping prerequisite before the
external harness. Assigning exact keys and numeric encodings therefore implements an already-scoped
prerequisite; it changes no component responsibility, listener topology, authority path, or
privileged helper and needs no ADR addendum. The later harness remains responsible for proving that
the selected XPC APIs enforce the contract.

## Frozen dictionary profile

The transport encoding is `xpc-dictionary-v0`, encoding version `0`. Every message is one top-level
`XPC_TYPE_DICTIONARY`; nested dictionaries/arrays, file descriptors, endpoints, Mach rights, and
extra objects are forbidden. Keys are fixed ASCII strings and values use only
`XPC_TYPE_UINT64`, `XPC_TYPE_DATA`, and `XPC_TYPE_STRING`.

Every request has these exact nine common keys:

| Key | XPC type | Constraint |
| --- | --- | --- |
| `capsule.protocol-version` | `XPC_TYPE_UINT64` | exactly `0` |
| `capsule.method-version` | `XPC_TYPE_UINT64` | exactly the selected method's version `0` |
| `capsule.message-tag` | `XPC_TYPE_UINT64` | exactly the method-derived numeric tag |
| `capsule.request-id` | `XPC_TYPE_DATA` | exactly 16 nonzero bytes; correlation only |
| `capsule.installation-id` | `XPC_TYPE_DATA` | exactly 16 bytes; current-state binding |
| `capsule.epoch-sequence` | `XPC_TYPE_UINT64` | `UInt53`; current-state binding |
| `capsule.epoch-digest` | `XPC_TYPE_DATA` | exactly 32 bytes; current-state binding |
| `capsule.audience` | `XPC_TYPE_STRING` | exact method-derived value |
| `capsule.purpose` | `XPC_TYPE_STRING` | exact method-derived value |

Every success or delivered refusal reply has exactly six common keys:
`capsule.protocol-version`, `capsule.method-version`, `capsule.message-tag`,
`capsule.request-id`, `capsule.status`, and `capsule.reason`. A success has status/reason `0/0` and
only its method-specific success body. A refusal has a nonzero closed status/reason pair and no
body key. Output cap/shape, short-write, pointer/length, or bridge-version violations remain local
integrity faults that terminate the future process without an application reply.

## Method-specific envelopes

The numeric message tags are `1` for `SubmitMainMJSV0`, `2` for `RegisterPlanV0`, and `3` for
`GetRegisteredPlanV0`; `0` is invalid. A tag only cross-checks the method already selected by the
role-specific service and entry point. No generic dispatcher or opcode bus exists.

The generated case set covers all six ordered collisions between one selected service/entry point
and either other frozen method's tag. Each refuses as `UNSUPPORTED/messageTag`; no foreign tag can
retarget dispatch. The binding identity remains the selected service + entry point + admitted peer
role, then the numeric tag is cross-checked. It is never tag-only, including where two services
reuse the same numeric namespace.

The method binding also retains the existing exact service, expected role, audience, and purpose:
`com.capsulecorp.capsule.daemon.cli.v0` / `internal-alpha-cli` /
`capsule.daemon.local.v0` / `capsule.ipc.submit-main-mjs.v0`;
`com.capsulecorp.capsule.supervisor.daemon.v0` / daemon /
`capsule.execution-supervisor.local.v0` / `capsule.ipc.register-plan.v0`; and
`com.capsulecorp.capsule.supervisor.broker.v0` / Broker /
`capsule.execution-supervisor.local.v0` / `capsule.ipc.get-registered-plan.v0`. The CLI binding
additionally retains signing-identifier candidate `com.capsulecorp.capsule.cli`; exact installed
daemon/Broker code identities remain S5 profile inputs rather than invented values here.

| Method | Request body keys | Request keys / cap | Success body keys | Reply keys / cap |
| --- | --- | ---: | --- | ---: |
| `SubmitMainMJSV0` | `capsule.job-proposal` (`DATA`, 1..2,097,152) | 10 / 2,097,152 | `capsule.registration-id` (`DATA`, exactly 16 nonzero) | 7 / 16 |
| `RegisterPlanV0` | `capsule.execution-plan` (`DATA`, 1..65,536); `capsule.role-bindings` (`DATA`, exactly 562); `capsule.source-manifest` (`DATA`, 87..95); `capsule.source` (`DATA`, 0..262,144) | 13 / 328,337 | `capsule.plan-registration` (`DATA`, 1..4,096) | 7 / 4,096 |
| `GetRegisteredPlanV0` | `capsule.registration-id` (`DATA`, exactly 16 nonzero) | 10 / 16 | plan, role bindings, plan registration, source manifest, and source under the same keys above | 11 / 332,433 |

All request application-data maxima, the Supervisor aggregate 2,626,696-byte limit, the daemon
aggregate 8,388,608-byte limit, zero application queue, and 10,000/5,000/2,000-millisecond method
deadlines remain unchanged. XPC framing bytes are not counted or claimed.

The flow corpus additionally isolates mixed-method byte bookkeeping at a smaller test limit:
one maximum `RegisterPlanV0` request plus one maximum `GetRegisteredPlanV0` request fills 328,353
bytes, the next byte refuses as `CAPACITY/ipc-process-in-flight-bytes`, and releasing the Get slot
admits one byte again. This exercises aggregate accounting and release without changing or claiming
that the production 2,626,696-byte cap is reachable independently of the existing connection and
request-count caps.

## Fixed status and refusal mapping

Status tags are closed: `OK=0`, `MALFORMED=1`, `UNSUPPORTED=2`, `SCHEMA=3`, `BINDING=4`,
`AUTHENTICATION=5`, `STALE=6`, `REPLAY=7`, `CAPACITY=8`, `TRUST_STATE=9`,
`LOCAL_FAILURE=10`, `RECOVERY_REQUIRED=11`, `SEMANTIC=12`, and `DOMAIN=13`. The last two preserve
the strict plan decoder's existing closed classifications rather than collapsing a recognized core
refusal into a different transport class.

Reason tags are closed: `none=0`, `keySet=1`, `valueType=2`, `dataWidth=3`, `dataCap=4`,
`zeroIdentifier=5`, `epochSequence=6`, `protocolVersion=7`, `methodVersion=8`, `messageTag=9`,
`methodBinding=10`, `currentState=11`, `capacity=12`, `coreRefusal=13`, and
`localIntegrityFault=14`.

Structural mappings are fixed: missing/extra keys and wrong/nested types are `MALFORMED`; exact-
width and zero-identifier failures are `SCHEMA`; cap-plus-one variable data is `MALFORMED`;
unknown protocol/method/tag is `UNSUPPORTED`; purpose/audience method-binding mismatch is
`AUTHENTICATION`; installation/epoch mismatch is `BINDING`; and a missed flow slot is `CAPACITY`.
For a recognized core refusal, the core classification selects the status and reason is always
`coreRefusal`. A local integrity fault has no reply; reason tag 14 is retained only as a fixture
diagnostic.

Validation precedence is fixed as protocol version, method version, method binding, current-state
binding, bounded application-data copy, then embedded-record/core validation. Joint
protocol/method/562-byte role-binding-record version mutations therefore select protocol first;
with protocol corrected, method version wins; with both outer versions correct, record version `1`
reaches the existing bounded Go decoder and returns `UNSUPPORTED/coreRefusal`.

Wrong-service and wrong-role fixture cases are delivered in-band refusals only under the explicit
precondition that the future receiving service's peer requirement and session gate admitted the
connection and the message reached the method-binding recheck. They are not OS peer-drop evidence.
A distinct future-harness oracle says an OS peer-requirement mismatch must produce no message
delivery and no application reply, while recording `currentEvidence=false`. A validating
dictionary cannot be used to infer peer authentication.

## Ownership, replay, and response loss

The future native owner must validate peer/session/flow slot, exact outer key/type/width shape,
installation, epoch, service, tag, audience, and purpose before copying application data. Go then
copies every input before return; native output buffers are fixed-cap and receive no Go pointer.

The request ID remains correlation only. Reuse after completion or reconnect is a fresh transport
call and creates no deduplication authority. Response-loss semantics remain method-specific:

- `SubmitMainMJSV0`: retry may create a fresh registration if the earlier downstream registration
  committed;
- `RegisterPlanV0`: a successful retry creates a fresh registration and retains the earlier one;
  and
- `GetRegisteredPlanV0`: retry is a repeatable exact-byte read by `RegistrationID`.

The same-width identifier substitutions make the domains explicit without claiming bytes carry a
hidden type tag: request-ID bytes copied from an installation ID remain correlation-only and grant
no authority; installation-ID substitution fails current-state binding; registration-ID bytes
copied from a request ID proceed only to the method-specific registration lookup and gain no
authority from their width. Local channel fields substituted with the signed ApprovalGrant values
`capsule.execution-supervisor` / `capsule.plan.approve` fail method binding. The inverse signed-
object substitutions cannot occur in these three envelopes because none carries an ApprovalGrant;
`SubmitApprovalV0` and `RequestAttemptV0` remain blocked and receive no invented encoding.

## Retained evidence and limitations

[`passive_native_xpc_contract.go`](../internal/execution/authorityplane/passive_native_xpc_contract.go)
contains the production-side passive constants and method-specific request/success/refusal specs.
It imports no XPC API and exposes no dispatcher or endpoint.

[`native-xpc-v0.contract.json`](../schemas/conformance/authenticated-local-ipc-v0/native-xpc-v0.contract.json)
is generated from the Node oracle. It retains exact keys, types, versions, tags, statuses, reasons,
key counts, caps, all three method envelopes, classification mappings, refusal-only replies,
cap-plus-one, missing/extra/wrong-type/nested, six ordered cross-service/tag collisions, joint
version precedence, same-width identifier substitutions, local/signed purpose-audience swaps,
wrong installation/epoch, explicitly keyed FD/endpoint/Mach-send-right smuggling,
reply-correlation/extra-key, response-loss, and local-integrity cases. The original three
response-loss entries remain byte-for-byte semantic entries in the C4 successor's exact five-entry
table. The Go verifier independently compares those fixtures with the production-side
specifications and rejects generic command/service/role/backend/path/image/mount fields.

Focused verification:

```sh
node scripts/generate-authenticated-local-ipc-conformance.mjs --check
node scripts/verify-authenticated-local-ipc-conformance.mjs
go test ./internal/execution/authorityplane
```

Remaining S3 work is the separately authorized one-time C/Objective-C harness in
`Shrimpworks/capsule-experiments`, consuming this exact fixture digest and testing actual peer-
requirement-before-delivery, message-derived `SecCode`, EUID/session, flow, copy, interruption,
deadline, response-loss, and process-fault behavior. Capsule must retain only its canonical result,
generated product conformance fixtures, and non-activating production-side contract. S5 installed
identity/update evidence, daemon/Broker consumers, protected state, approval, runtime/backend, and
product admission remain separate blockers.

Exactly three methods are frozen here. Both `SubmitApprovalV0` and `RequestAttemptV0` remain
outside this historical S3 scope. The C4 successor now freezes their passive tags, envelopes,
case-derived deadlines, replay, and response-loss evidence; activated callers remain `BLOCKED`.

Parent status: owner-only hostile-`.mjs` internal alpha remains
`IN_PROGRESS — TRENDING_GOOD`; product admission remains `BLOCKED`.
