# Phase 2A parallel-review synthesis

Status: reconciled against the passive contract-foundation implementation merged by PR #9 as
`3a75098`.

## Sources and verdict

Three independent read-only tasks reviewed the contract shape, deprecated-`Job` migration, and
fail-closed corpus. They agree on the authority split and do not identify a reason to abandon the
current direction.

| Review | Task ID |
| --- | --- |
| Contract shape | `019fbd6b-0de7-7201-befc-be983552e5ed` |
| Deprecated mixed-`Job` migration | `019fbd6b-0de5-7003-82b2-40049e07fcb7` |
| Fail-closed conformance corpus | `019fbd6b-0de5-7003-82b2-3febecc775aa` |

The cross-phase integration history is indexed in the
[workstream and evidence ledger](WORKSTREAM_EVIDENCE_LEDGER.md).

The passive foundation is the correct stopping point for this PR. It deliberately chooses a
narrower contract than some review sketches and does not guess unresolved resource, transport, or
backend semantics. It activates no public or execution path.

No additional isolation-backend spike is required before the registered-plan/fake-backend slice.
The remaining immediate work is contract decision, strict decoding, conformance, and lifecycle
implementation. The existing libkrun/runtime P0 campaigns still block handling user bytes.

## Consensus reconciled with the implementation

| Topic | Reconciled outcome | Current state |
| --- | --- | --- |
| Authority split | `JobProposal`, `ExecutionPlan`, and `PlanRegistration` are closed, distinct objects; this is not a rename of `Job`. | Candidate shapes and role-specific language views implemented. |
| Public authority | Omit network, process, environment, native/FFI, package, image, backend, mount, socket, host-path, guest-path, and arbitrary-schema fields entirely. | Enforced structurally by the passive proposal schema. |
| First input/output | One inline JSON input and one bounded inline JSON output with fixed `primary-data` and `transformed-json` roles. | Implemented; opaque/file inputs and file artifacts remain deferred. |
| Source paths | Use an ASCII-only relative candidate grammar; never accept live host paths or agent-selected guest paths. | Structural grammar implemented; entrypoint existence, aggregate bounds, and manifest canonicalization remain semantic work. |
| JSON numbers | Admit safe integers only in the first candidate, avoiding floating-point/negative-zero/exponent ambiguity. | Schema and branded TypeScript scalar implemented; canonical input bytes/digest remain unresolved. |
| Requested limits | Expose only independently justified wall time; omission means a trusted default and ceiling-plus-one must reject without clamping. | Structural candidate implemented; policy/default resolution is not implemented. |
| Internal identity | Hash exact received deterministic-CBOR plan bytes; never trust a caller-supplied digest or replace bytes with decode/re-encode output. | Candidate fixture rule and decoded role-specific digest types implemented; production decoder/storage integration remains open. |
| Registration | Supervisor-issued response binds registration ID/sequence, plan digest, installation, epoch, Supervisor, and expiry; it carries no replacement plan bytes. | Candidate CDDL and decoded views implemented for authenticated typed local IPC; replay/expiry persistence semantics remain open. |
| Internal IDs | Installation, registration, and Supervisor IDs are distinct 16-byte nonzero domains; digests are distinct 32-byte domains. | Go/TypeScript constructors enforce the rule. CDDL expresses width, with nonzero ID rejection documented as semantic validation. |
| Complete plan | A launch-capable plan must bind every applicable resolved control, profile, posture, content, and transport value. | Intentionally incomplete candidate; it cannot authorize execution. |

Suggestions for vCPU, RAM profile, scratch, console-prefix, safety-display, runtime, posture, and
audience fields are retained as later inventory, not silently added to the minimum candidate. Their
exact vocabulary and values require accepted contract decisions and applicable backend evidence.

## What remains outside this PR

### Contract decisions

- strict JSON and CBOR media types, raw-byte, nesting, string, collection, and allocation caps;
- canonical inline-JSON byte identity and digest construction;
- complete source-manifest canonicalization and aggregate source bounds;
- final internal media types, field sets, labels, identifier textual forms, and digest-zero policy;
- registration duplicate/deduplication behavior, sequence scope, TTL source, expiry boundary,
  replay, retention, and trust-epoch interaction;
- final error/violation vocabulary and its mapping to decoder, schema, semantic, policy, binding,
  control, domain, and stale failures;
- the complete launch-capable `ExecutionPlan` control/reference set.

### Implementation and evidence

- strict pre-schema JSON decoding that retains raw malformed fixtures and rejects duplicate keys,
  invalid UTF-8, trailing documents/data, and cap-plus-one inputs;
- bounded deterministic-CBOR decoders/wrappers and negative canonical-on-wire corpora in every
  participating language, including Swift where the Broker participates;
- semantic proposal planning for source, slots, profile resolution, trusted defaults/ceilings,
  content identity, policy, and exact-or-refused controls;
- Supervisor registration with exact-byte custody, independent validation, durable sequence/expiry,
  wire/stored-record separation, and no authority-bearing state change on rejection;
- a fault-injectable fake backend and registered-plan lifecycle that creates no guest;
- the later coordinated daemon/SDK/MCP cutover and removal of every mixed-`Job` bypass-shaped
  surface.

## Dependency-ordered next work

### Task 1: Freeze boundary decisions needed by executable validators

Record the raw caps/media types, path and inline-JSON identity rules, ID/digest zero semantics,
registration replay/expiry semantics, and internal error classifications in the relevant ADRs or an
explicit addendum. Do not complete the launch-capable plan vocabulary with unevidenced fields.

Proposed resolution: [Phase 2B boundary decisions](PHASE_2B_BOUNDARY_DECISIONS.md) and
[ADR-0023](adr/0023-bound-protocol-decoding-and-registration.md).

Acceptance:

- every value needed by a raw decoder or registration-state test is exact;
- no decision expands public authority or claims a backend control;
- ADR-0019 remains Proposed unless all of its acceptance conditions actually pass.

Dependencies: none.

### Task 2: Build one manifest-driven conformance corpus

Add byte-retaining accept/reject fixtures for shared scalars, raw JSON, deterministic CBOR,
`JobProposal`, `ExecutionPlan`, `PlanRegistration`, and cross-object/domain substitutions. Each case
records its first owning layer, provisional/frozen classification, participating languages, and
`authorityStateChanged: false` where rejected.

Acceptance:

- malformed bytes are stored as bytes rather than parsed JSON values;
- exact maxima and cap-plus-one cases exist;
- the runner fails on skipped/unlisted fixtures and verifies exact accepted bytes/digests.

Dependencies: Task 1 for cases whose exact oracle is currently undecided.

Implementation progress: Phase 2B Tasks 2.1 through 2.4 now provide the closed manifest, repository
integrity runner, foundational raw/media/scalar/CBOR fixtures, fixed proposal resolver contexts,
source/canonical-inline-input known answers, and exact plan/registration/domain/state fixtures.
Task 3A now implements 62 TypeScript raw/schema cases, Task 4A implements 81 Go internal
media/scalar/CBOR/wrapper cases, and Task 4B implements all 40 Go registration-state cases through
a separate unwired component. Semantic planning, Swift, and consumer wiring remain pending; no
endpoint or backend is activated.

### Task 3: Implement strict proposal decoding and semantic planning as an unwired library

Implementation progress: Task 3A is complete. `decodeJobProposal(Uint8Array)` implements the strict
raw and closed-schema boundary and returns only a frozen passive candidate or fixed internal
refusal. Task 3B semantic trusted resolution and minimum-plan construction remain pending.

Implement the public raw-decoder boundary and semantic `JobProposal` validation without adding a
daemon, SDK, or MCP endpoint. Resolve trusted defaults and ceilings into a minimum exact plan or
return a typed refusal.

Acceptance:

- duplicates, trailing data, malformed UTF-8, cap-plus-one inputs, invalid entrypoints, unknown or
  inactive profiles, and over-ceiling limits fail at the documented owner;
- missing wall time becomes an exact trusted-default value and origin; no value is clamped;
- no user content, backend, approval, or launch authority is introduced.

Dependencies: Tasks 1 and 2.

### Task 4: Implement bounded internal wrappers and exact-byte registration

Implementation progress: Task 4A is complete. Go now strictly predecodes, object-decodes, copies,
hashes, and role-binds exact candidate bytes without persistence or authority creation. Task 4B
implements exact Supervisor registration state using a bounded, fault-injectable development
store. Applicable Swift/TypeScript wrappers remain pending, and the older `SupervisorCore`
registration path is not activated.

Implement object-specific deterministic-CBOR validation and separate the Supervisor-issued wire
registration from its stored record containing exact plan bytes. Recompute the digest from copied
received bytes and validate every installation/epoch/Supervisor/reference binding independently.

Acceptance:

- decode/re-encode bytes never replace authoritative received bytes;
- wrong-domain objects, IDs, digests, noncanonical encodings, stale registrations, and mutated
  buffers fail without creating authority-bearing state;
- applicable Go, Swift, and TypeScript implementations agree on retained fixtures.

Dependencies: Tasks 1 and 2.

### Task 5: Build the registered-plan/fake-backend vertical slice

Connect the planner and Supervisor only to a fault-injectable backend that cannot create a guest.
Exercise durable registration, one-use approval/attempt binding, crash recovery, cleanup outcomes,
and evidence composition without user content.

Acceptance:

- Plan A approval cannot realize Plan B or a second attempt;
- every injected post-side-effect failure becomes destroyed, unresolved, or quarantined state;
- the daemon cannot manufacture ordinary terminal success.

Dependencies: Tasks 3 and 4.

### Task 6: Perform the atomic public cutover

Only after the new vertical slice is reviewable, migrate daemon, SDK, MCP, examples, schemas, and
documentation together. Remove `prepare_job`/`run_job`, `PreparedJob`, `Backend.Execute`,
`RuntimeAdapter`, the mixed `Job` type/helper/schema/example, and any execute-time plan bytes or
backend flags. Old clients receive an explicit unsupported-version response; there is no permissive
adapter.

Acceptance:

- repository search and tests show no public mixed-`Job`, guest/mount path, capability-union, or
  replacement-plan execution surface;
- MCP attempt operations accept only frozen identifiers and return the fixed agent summary;
- complete verification and adversarial conformance pass before activation.

Dependencies: Task 5.

## Parallelization

After Task 1 freezes the shared vocabulary, corpus assembly, language-specific strict-wrapper
implementations, and migration search assertions can proceed in parallel. Planner/registration
integration and public cutover remain sequential because they share authority-bearing contracts and
durable state.

Backend/runtime P0 work may continue independently, but no result from it is imported into these
contracts until an accepted ADR and retained evidence justify the exact field or value.
