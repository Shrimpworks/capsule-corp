# Typed source/input/completion transport research

Date: 2026-08-07

Work item: R2 typed transport design

Status: `PASSED`

Parent governed-runtime engineering: `IN_PROGRESS — TRENDING_GOOD`

C5a passive contract: `PASSED` in the later
[v1 contract](protocol/TYPED_GUEST_TRANSPORT_V1.md)

C5b controlled harness: `BLOCKED`

Product runtime/profile admission: `BLOCKED`

## Defensive scope and result

This read-only research reconciles the retained backend-independent P0-3 candidate, C2A, C2B v4,
the exact v27 hostile-denial checkpoint, and the passive durable completion-last transaction. It
creates no process, descriptor, runtime adapter, VMM, guest, artifact, signature, installation,
store mutation, or product authority.

The result selects one passive design for C5a: three role-distinct, attempt-bound, one-frame byte
streams with Supervisor-created endpoints; a distinct trusted launcher that completely validates
source and input before starting the workload; a workload that never receives the completion
endpoint; continuous host completion draining with cap-plus-one retention and independent total
accounting; and a launcher-written commit trailer after bounded result validation and child-tree
absence. A valid guest frame is only one observation. The Supervisor may publish terminal state
only after separately proving lifecycle, teardown, authoritative absence, and cleanup, then
committing the immutable durable completion transaction.

At this research checkpoint, this closed the R2 design question only and the byte contract was not
yet frozen. C5a later froze the passive v1 successor without execution. The governed transport is
still not executed, and v27's diagnostic console proof is not reclassified as final typed
transport.

## Reconciled inputs

| Input | Controlling result for C5a |
| --- | --- |
| P0-3 backend-independent candidate at archive commit `0d8233b55f153b27a901a9ec45a3834208e3aa86` | Retain the 152-byte source/input header, 160-byte completion header, 64-byte commit trailer, role/magic/domain separation, exact-offset commit, cap-plus-one drain, and 43-vector disposition model as the compatibility baseline. Its 1,048,576-byte source candidate is superseded for the first slice by C2A's narrower limit. |
| C2A passive execution profile | Use inclusive 262,144-byte source, canonical-input, and result-JSON payload limits; 262,296-byte source/input physical limits; 262,368-byte completion physical limit; and 262,369-byte completion retention limit. Enforce an approved lower value exactly; never clamp or resize. |
| C2B v4 materialized runner | Preserve runner FDs 5/6/7 and ordered ports `capsule.source`/`capsule.input`/`capsule.completion`; guest launcher FDs 3/4/5; exact directionality; implicit console/init/vsock disablement; and Supervisor-only start, drain, cancellation, teardown, and absence authority. |
| Exact v27 hostile-denial checkpoint | Carry forward the exact denial/restoration controls and lifecycle ordering. V27 used immutable in-root source/input plus diagnostic console transport, so it contributes containment and teardown evidence but no typed-transport implementation evidence. |
| Durable completion-last contract | Treat the guest frame as an input to, never a substitute for, the Supervisor-owned immutable store transaction. Response loss and restart resolve by `AttemptID` and return byte-identical retained state. |

The first source-transfer object is one byte-exact UTF-8 `main.mjs`, represented by the already
registered single-member `SourceManifest` and its exact source bytes. C5a must bind the exact
registered transfer representation without expanding it into a multi-file container or restoring
the earlier 1 MiB aggregate candidate.

## C5a framing baseline

C5a should freeze a versioned successor only after its selected Supervisor and launcher
implementations independently reproduce the retained known answers. The baseline is:

| Channel | Direction | Exact physical maximum | Required shape |
| --- | --- | ---: | --- |
| source | Supervisor to trusted launcher | 262,296 | 152-byte role-1 header plus at most 262,144 exact registered source-transfer bytes |
| canonical input | Supervisor to trusted launcher | 262,296 | 152-byte role-2 header plus at most 262,144 canonical inline-JSON bytes |
| completion | trusted launcher to Supervisor | 262,368 | 160-byte role-3 header, at most 262,144 strict typed-JSON bytes, then one 64-byte commit trailer |

All integers remain unsigned big-endian. Each header must bind its role-specific magic, version,
role, exact header length, nonzero role-distinct attempt and registration identifiers, retained plan
digest, selected runtime-profile digest, declared payload length, and exact payload SHA-256. The
completion header additionally binds one closed terminal status, zero flags, and zero reserved
space. Its trailer binds its own magic/version/length/role, the same attempt, and SHA-256 of the
complete completion header plus payload.

The trailer is valid only at `160 + declared payload length`, is written after all other completion
bytes, and must be the final 64 bytes. Embedded magic is payload; an early, missing, duplicate,
changed, or trailing trailer refuses. EOF is stream termination only. It never commits completion.

C5a must decide whether to retain the candidate status values `succeeded`, `workload failed`,
`result invalid`, and `child terminated` or replace them with a separately versioned closed result
contract. Any non-success payload stays a fixed value such as JSON `null`; guest-controlled error
prose, paths, diagnostics, names, and timing may not enter the frame or public summary.

## Passive state machine

The transport is one attempt-scoped join of three directional machines, not one duplex session.
Every transition is monotonic; reset or cancellation moves the affected machine to refusal and
cannot return it to a writable/acceptable state.

| State | Required action and invariant | Successful next state |
| --- | --- | --- |
| `BOUND` | Resolve the committed `AttemptID`, registration, exact plan/source/input, runtime profile, approved limits, and descriptor roles from Supervisor/store truth. Accept no execute-time replacement bytes or paths. | `ENDPOINTS_READY` |
| `ENDPOINTS_READY` | Create fresh directional endpoints with `CLOEXEC`; validate access modes, object identity, aliasing, flags, and exact runner/launcher manifests. Start every required drain before guest authorization. | `RUNNER_READY` |
| `RUNNER_READY` | Revalidate the runner identity and durable-before-start record; issue exactly `G` plus EOF only for the committed attempt. | `INPUT_TRANSFER` |
| `INPUT_TRANSFER` | Write the fixed source and input frames through distinct writers. Handle partial progress, backpressure, and interruption without changing bytes or caps. Close each writer after its complete frame. | `LAUNCHER_VALIDATED` |
| `LAUNCHER_VALIDATED` | Launcher receives exactly one frame on each role, validates endpoint/magic/role/version/bindings/length/digest/trailing bytes, and creates immutable child inputs. No child exists before both pass. | `CHILD_RUNNING` |
| `CHILD_RUNNING` | Start one fixed child with fixed argv/environment/cwd/FDs. Withhold completion FD/node. Bound and drain result/stdout/stderr, then wait for the exact child tree to be absent. | `RESULT_VALIDATED` |
| `RESULT_VALIDATED` | Validate status, typed JSON, exact result length/digest, fixed and approved caps, and runtime/input integrity. Compose the completion header and payload. | `TRAILER_WRITTEN` |
| `TRAILER_WRITTEN` | Write the complete 64-byte commit trailer last; close the launcher completion writer. No later byte is permitted. | `FRAME_OBSERVED` |
| `FRAME_OBSERVED` | Host drain retains at most physical cap plus one, counts and drains all bytes, and validates the complete exact-offset frame independently of EOF and runner status. | `TERMINAL_PROOF` |
| `TERMINAL_PROOF` | Join frame observation with terminal runner lifecycle, child-tree absence, authoritative runner absence, cleanup false, and every immutable authority/runtime/backend binding. | `DURABLE_COMMIT` |
| `DURABLE_COMMIT` | Atomically publish one immutable completion record and fixed summary. Only now may a caller receive terminal output. | `COMPLETE` |

Source or input acceptance alone authorizes no child start. A completion frame accepted before
runner teardown remains only `FRAME_OBSERVED`. Runner exit, including zero, authorizes no frame and
no durable commit.

## Ownership and drain rules

- The Supervisor owns endpoint creation, writer/drainer tasks, runner creation, start authorization,
  cancellation, signaling, teardown, authoritative absence, and durable commit.
- The runner receives only the sealed FDs 0 through 7 and configuration-free fixed call sequence.
  It cannot receive attempt, plan, path, backend, image, mount, or replacement profile bytes.
- The trusted launcher owns the guest completion writer only until it closes it after the trailer.
  It validates both inbound frames before child creation and retains authority over the workload's
  fixed inputs and bounded result capture.
- The workload receives immutable source/input projections and bounded result/stdout/stderr sinks.
  It receives neither the completion FD nor a reopenable completion node, and descendants inherit
  no path to it.
- Each host writer handles short writes explicitly. Zero progress is bounded as a stall; partial
  progress followed by error is a transport fault, never a smaller valid frame.
- The completion reader begins before authorization, retains no more than 262,369 bytes, maintains
  a separate overflow-safe total count, and continues draining oversize output until producer
  closure or external teardown. Oversize is irreversible even if the retained prefix is valid.
- Cancellation and wall expiry act independently of blocked writers/readers. External teardown
  closes Supervisor-owned endpoints only as part of convergent cancellation; it never converts
  resulting EOF into success.

## Cancellation, reset, and response-loss taxonomy

| Event point | Required disposition |
| --- | --- |
| Before `G` | Refuse start, close attempt endpoints, reap the runner, prove absence, and retain no completion. |
| During source/input transfer | Stop new writes, close owned writers, require launcher refusal without child creation, drain/reap/teardown, and retain a transport-fault fact. |
| During child execution | Prevent any new completion acceptance, externally terminate/reconcile the runner and child tree, and retain no ordinary success. |
| After payload but before complete trailer | `MISSING_COMMIT` or the more specific transport reset/fault; never accept payload digest or EOF alone. |
| After a complete valid trailer but before runner absence | Preserve `FRAME_OBSERVED`; ordinary success remains blocked until lifecycle and absence proofs pass. |
| Guest port close/reset | Terminal fault for that attempt. A source/input reset forbids child creation; a completion reset preserves a frame observation only if all exact bytes, final trailer, and no trailing data were already retained. |
| Host completion reader death | Terminal transport fault and immediate teardown. A producer cannot prove that unread bytes would have been valid. |
| Backpressure or zero progress | Continue bounded draining/writing while the independent deadline is live. On expiry, classify stall, teardown externally, and refuse. |
| Runner/launcher death before trailer | Preserve partial accounting and classify missing commit plus lifecycle failure. |
| Runner/launcher death after trailer | Preserve valid framing evidence but fail ordinary success unless the separately required lifecycle, child-tree, teardown, and absence facts are valid. |
| Caller response lost after durable commit | Reopen/retry by `AttemptID`; return byte-identical stored completion and summary. Do not rerun the guest or accept replacement result bytes. |
| Store outcome indeterminate | Fence the store instance. Reopen must validate exactly one complete predecessor or successor before replay; transport bytes do not resolve storage truth. |
| Cancellation concurrent with commit | The durable transaction may publish only if all terminal inputs were already true. Otherwise cancellation wins and no completion record is created. A published immutable record is never rolled back by a later lost reply. |

## Completion-last proof inputs

The durable transaction consumes independently retained facts. The guest-controlled frame cannot
provide or override Supervisor-owned facts.

| Proof class | Required facts |
| --- | --- |
| authority | Attempt, consumed approval, registration, plan, installation, trust epoch, Supervisor, and approval authorization cross-link exactly |
| input | source manifest and exact source bytes, canonical input, declared/observed lengths, digests, role, version, and endpoint integrity all pass |
| runtime/profile | runtime bundle, profile registry, backend validation/configuration/implementation/instance, root, runner, and immutable lifecycle bindings match retained truth |
| result | exactly one bounded role-3 frame, closed status, strict typed JSON, exact length/digest, fixed cap, approved plan cap, and valid trailer last |
| lifecycle | trusted launcher reports child-tree absence; runner lifecycle is terminal; destroy/teardown is resolved; authoritative runner absence is independently proven; cleanup is false |
| publication | deterministic bounded transcript and fixed non-guest-controlled summary are atomically published once with `committed-last` |

The transport completion digest and payload are observations from the admitted trusted-launcher
model, not attestation that the guest kernel or workload was uncompromised or correct.

## C5a restoration matrix

C5a must make each mutation mechanically detectable without executing libkrun or a guest. C5b must
later carry the applicable cases through the separately authorized controlled composition.

| Mutation or fault | Required refusing control |
| --- | --- |
| source/input/completion endpoint swap | endpoint identity plus role-specific magic and role |
| duplicated, aliased, reused, missing, extra, or wrong-mode FD | pre-start object/access/open-description manifest |
| changed `CLOEXEC`, nonblocking, or status flags | pre/post descriptor-flag canary; dedicated open descriptions |
| completion FD/node inherited or reopened by workload/descendant | launcher child manifest, node access denial, and descendant inventory |
| wrong, duplicate, equal-domain, stale, or zero attempt/registration ID | typed ID/domain and current-attempt binding validation |
| wrong plan, runtime-profile, source, input, or payload digest | exact retained binding and recomputed digest |
| header version/role/length/magic/flags/reserved mutation | closed structural decoder before allocation/copy |
| source/input/result exact cap plus one | declared-length cap before allocation plus physical total accounting |
| truncated frame, partial-error, zero progress, or early EOF | exact-length machine and bounded stall/fault disposition |
| duplicate frame or any trailing byte | one-frame stream and exact physical length |
| malformed/duplicate-key/second-document/unsafe-integer/depth/node JSON | bounded strict typed-JSON validator |
| unknown status or guest-controlled failure payload | closed status set and fixed non-success payload |
| early, missing, duplicate, wrong-attempt, or changed trailer | calculated-offset commit parser and header-plus-payload digest |
| payload/trailer valid prefix followed by output flood | full drained-byte count; irreversible oversize/trailing refusal |
| host reader stall/death, `SIGPIPE`, peer close, or backpressure | stop-aware continuous drain, explicit partial progress, external deadline/teardown |
| launcher crash before/after trailer | separate frame and lifecycle dispositions; no exit-derived success |
| runner zero without completion | `MISSING_COMMIT` plus lifecycle-only exit classification |
| cancellation before start/during input/during child/during commit | monotonic cancellation state and no post-cancel child or commit mutation |
| guest port open/close/reset race or invalid control ID/event | governed bounds-checked console corpus and terminal reset classification |
| queue index/size/descriptor-chain corruption or resource exhaustion | governed sanitizer/fuzz/coverage corpus and bounded teardown |
| re-enable implicit console, init, vsock, network, virtiofs, or extra block/device | exact runner call-order and device-inventory verifier |
| use v27 diagnostic console bytes as typed completion | typed channel/magic/frame requirement; diagnostic evidence is ineligible |
| omit child-tree absence, runner absence, or cleanup resolution | terminal-proof join refuses durable publication |
| retry with changed result after response loss | immutable `AttemptID` replay returns `REPLAY` and preserves the first record |
| missing/early/duplicate/corrupt durable commit or mixed version | store reopen validation returns `RECOVERY_REQUIRED` without rewrite |

## C5a acceptance and stop boundary

C5a may become `PASSED` only when one passive, dependency-free contract and at least two independent
decoders reproduce the narrowed exact/cap-plus-one known answers; every state transition,
disposition precedence, descriptor role, restoration mutation, and completion-to-store projection
is frozen; and no product package activates a runner or listener. C5a must explicitly resolve the
closed status vocabulary and exact single-`main.mjs` transfer encoding.

C5b remains separately authorized work. It must name the exact runtime/profile successor, owned
disposable guest, host, fixture digest, process/service names, defensive intent, and evidence
destination. Stop if progress would require changing C2A/C2B/v27 retained evidence, widening an
approved cap, treating diagnostic console output or EOF as completion, giving the workload the
completion endpoint, running a guest without that authorization, or promoting a passive or fixed
experiment into product admission.

## Sources

- [`GOVERNED_DENO_CORE_C2A_EXECUTION_PROFILE.md`](protocol/GOVERNED_DENO_CORE_C2A_EXECUTION_PROFILE.md)
- [`GOVERNED_DENO_CORE_C2B_MATERIALIZED_PROFILE_V4.md`](protocol/GOVERNED_DENO_CORE_C2B_MATERIALIZED_PROFILE_V4.md)
- [`FIRST_OWNED_GUEST_EXECUTION_CHECKPOINT.md`](FIRST_OWNED_GUEST_EXECUTION_CHECKPOINT.md)
- [`DURABLE_COMPLETION_CONTRACT.md`](DURABLE_COMPLETION_CONTRACT.md)
- [Pinned P0-3 framing result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-p0-3-protocol-conformance/RESULTS.md)
- [Pinned P0-3 candidate and fault-harness source](https://github.com/Shrimpworks/capsule-experiments/tree/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-p0-3-protocol-conformance)
