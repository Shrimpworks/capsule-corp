# M2/S1 conformance slice plan: `RegisterPlanV0` / `GetRegisteredPlanV0`

Date: 2026-08-03

Status: retained dependency-aware implementation plan for the next
[S1/M2 passive contracts and fixed fixtures](AUTHENTICATED_LOCAL_IPC_PLAN.md#m1-source-validator-r1-r5b-and-s1m2-sequential-contract-boundary)
slice. No byte maximum in this document is calculated or invented; every number is either quoted
verbatim from ADR-0029/ADR-0034 prose (marked as unverified until fixture-generated) or explicitly
marked "blocked on M1."

Current dependency reconciliation (2026-08-04): M1 and Source Validator R1 have passed, but this
slice remains `BLOCKED`. Accepted ADR-0036 requires the remaining sequential Source Validator path
first: R2 unsigned launcher/parser construction, separately authorized R3
signing/install, R4 confinement/reactive-resource/residue corpus, R5D daemon consumer, and R5B
Approval Broker consumer. Only then does the M2/S1 checkpoint update the historical M1-blocked
items below and decide whether to run this plan. V0/V1/V2 evidence remains unchanged; no Source
Validator result becomes a Supervisor registration/fetch field.

Reviewer: Claude, independent read-only planning at the request of the Capsule orchestrator
(codex).

## Scope and method

Read-only planning. It read `AGENTS.md`, [ADR-0029](adr/0029-select-authenticated-local-ipc-topology.md),
[ADR-0034](adr/0034-freeze-mjs-first-release-contract.md),
[`AUTHENTICATED_LOCAL_IPC_PLAN.md`](AUTHENTICATED_LOCAL_IPC_PLAN.md),
[`AUTHENTICATED_LOCAL_IPC_S1_CONSISTENCY_STOP.md`](AUTHENTICATED_LOCAL_IPC_S1_CONSISTENCY_STOP.md),
`schemas/authority/field-authority-manifest.json` and its JSON Schema, existing plan-v0 CDDL
(`schemas/cddl/candidates/execution-plan-v0.cddl`, `plan-registration-v0.cddl`), and existing Go
code in `internal/protocol/v0candidate/`, `internal/execution/registrationstate/`, and
`internal/execution/approvalattempt/`.

**Grounding fact:** no Go, TypeScript, or `.mjs` fixture file implements the `RegisterPlanV0`/
`GetRegisteredPlanV0` wire shapes anywhere in the repository today; the only hits are prose in
ADRs/plan docs. What already exists and is reusable as-is: `internal/protocol/v0candidate/
contracts.go` (decoded Go views `ExecutionPlan`/`PlanRegistration` plus nominal ID/digest types),
`internal/execution/registrationstate/component.go`'s `Component.RegisterPlan` (the exact Go
operation ADR-0029 says `RegisterPlanV0` must call), and the closed field-authority manifest (30
profiles, 15 targets, 164 fields as of this review). No Go struct or CDDL yet exists for the
562-byte role-binding record, the single-member `SourceManifest` v0, or any
`RegisterPlanV0`/`GetRegisteredPlanV0` request/reply envelope — this slice must define them. There
is also no read-only Broker fetch facade in Go yet; ADR-0029/the IPC plan both scope its full
implementation to S2, but its typed projection must be designed now so S1/M2 fixtures and S2 share
one contract.

## Field ownership

### Common header (both calls)

| Field | Owner / origin | Authority effect | Existing profile fit |
| --- | --- | --- | --- |
| protocol version (`0`) | structurally fixed by native front end | none | new: `ipc-message-discriminator` |
| request ID (16 bytes, nonzero) | native front end, per call | none — correlation only, never a store key | new: `ipc-correlation-only` |
| installation ID (16 bytes) | caller-submitted; authority from matching already-open Supervisor state | narrows-authority | `plan-trust-state` (closest existing analogue) |
| epoch sequence/digest | caller-submitted, resolved against active epoch | narrows-authority | `plan-trust-state` |
| message tag | derived from service + tag by native layer; caller values cannot override | none | new: `ipc-message-discriminator` |
| audience / purpose strings | derived from service + tag, not caller-suppliable | none (selects predeclared alternative) | new profile needed — no existing profile uses `selects-predeclared-alternative` for IPC audience/purpose |

### `RegisterPlanV0` request body

| Field | Origin role | Independent checks / authority effect | Profile fit |
| --- | --- | --- | --- |
| exact plan bytes (1..65,536) | `daemon-planner` | Supervisor predecodes/decodes/role-binds/hashes/retains; authority only on commit | `plan-content` (exact existing match) |
| 562-byte plan-role-binding record | `daemon-planner`-supplied nominal identities | Supervisor resolves every nominal role from trusted local state — daemon labels alone grant nothing | new: `role-binding-nominal` (as-submitted) + `role-binding-resolved` (Supervisor-owned, `changes-authority`) |
| exact canonical source-manifest bytes (87..95) | `daemon-planner`-derived from agent source | Supervisor canonically decodes, hashes, checks plan source role, retains | new: `source-manifest-content`, modeled on `plan-content` |
| exact `main.mjs` bytes B (0..262,144) | `agent`-originated, daemon-copied | Supervisor checks length/UTF-8/BOM/profile/digest/manifest-membership/pass-through identity, retains defensive copy | new: `source-bytes-content`, `originRole: agent` (not `daemon-planner` — the daemon only copies), `guestControl: derived-from-untrusted` |

### `GetRegisteredPlanV0`

| Field | Origin | Effect | Profile fit |
| --- | --- | --- | --- |
| request: nonzero `RegistrationID` (16 bytes) | Approval Broker (caller) | lookup key only; must bind to an existing registration | new: `registration-lookup-key` |
| reply: exact retained plan (≤65,536) | Execution Supervisor (defensive copy) | evidence/read-only; must equal committed `plan-content` bytes exactly | `plan-content` (reused for the retained copy) |
| reply: 562-byte binding record | Execution Supervisor | Broker fetch returns the Supervisor-retained resolved projection, "not caller labels reinterpreted on read" | `role-binding-resolved` |
| reply: retained registration (≤4,096) | Execution Supervisor | unchanged | existing `registration-state`/`registration-plan` |
| reply: source manifest / `main.mjs` bytes | Execution Supervisor | retained copies | `source-manifest-content` / `source-bytes-content` |

**Key ownership rule to carry into implementation:** the daemon's role labels in the 562-byte record
are never trusted merely because the daemon supplied them — the Supervisor's own resolvers construct
the trusted `ExecutionPlanRoleBindings`. This is why the binding record needs two profiles (nominal
vs. resolved), matching the manifest's existing pattern for fields with different authority strength
at different pipeline stages (e.g. `plan-content` vs. `plan-trust-state`).

## Request/reply structure derivation methodology

Process, not final layout:

1. **Enumerate fields from canonical source objects only** — `execution-plan-v0.cddl`,
   `plan-registration-v0.cddl`, ADR-0029's 562-byte role-projection table, ADR-0034's single-member
   `SourceManifest` v0 CBOR map, and the raw `main.mjs` byte sequence itself (no wrapper object). A
   field with no traceable canonical source is out of scope; adding one requires a new ADR per
   `AGENTS.md`.
2. **Distinguish wire-CBOR objects from bridge-only structures.** `ExecutionPlan`,
   `PlanRegistration`, and `SourceManifest` are deterministic-CBOR objects (closed CDDL map, integer
   labels, exact bounds, known-answer fixture). The 562-byte role-binding record is explicitly not
   one of these — ADR-0029 calls it "a bridge-only, fixed-layout v0 projection," never signed, never
   stored as a plan. The IPC envelope is a third category: an application-level byte-budget
   accounting, not a CBOR or Mach-message framing claim.
3. **Derive per-field independent-check ownership before layout** — which validator structurally
   checks each field, which resolver (if any) turns a submitted label into a trusted value, and
   which refusal classification applies if wrong (`MALFORMED`/`SCHEMA`/`BINDING`/`STALE`/`REPLAY`/
   `UNSUPPORTED`/`AUTHENTICATION`). This ordering happens before any layout/cap decision.
4. **Derive the reply shape as defensive copies, never re-derivation.** `GetRegisteredPlanV0` must
   return the exact retained bytes committed at registration time — the reply's structure is the
   request's committed structure plus the registration wrapper, never a re-encoded view.
5. **Never derive a cap from the envelope; derive the envelope's cap from its fields.** Aggregate
   byte budgets are computed after every individual field's cap is independently fixed, never
   assumed top-down — exactly what ADR-0034 requires for the candidate aggregate numbers below.
6. **Cross-language identity requirement.** Once a canonical definition is written, Go/TypeScript/
   Swift-readable fixtures must agree byte-for-byte before any cap is trusted.

## Cap calculation method (methodology only — no new numbers)

ADR-0034 states candidate figures — 95-byte source-manifest cap, 262,144-byte source cap,
328,337-byte `RegisterPlanV0` request cap, 332,433-byte `GetRegisteredPlanV0` reply cap — but is
explicit that "the passive fixture slice must generate them from the closed message contract and
field-authority manifest and must retain exact-boundary/cap-plus-one vectors before either value is
frozen in code... implementations may not copy the arithmetic alone... If generated canonical bytes
disagree, the fixture evidence controls and this ADR must be revised."

This slice's cap-calculation step is therefore: write a fixture generator, not a spreadsheet. It
must read the same closed canonical definitions named above, compute each aggregate cap as a sum of
independently-fixed per-field maxima, emit both the boundary and boundary-plus-one values as
fixtures, and assert the generated aggregate equals the currently-stated ADR figure — a mismatch
blocks acceptance and forces an ADR revision. The same predecoder-budget discipline ADR-0023 already
uses for other objects (raw bytes / nesting depth / total data items / max map entries / max array
elements per object) needs an equivalent row for the new `SourceManifest` v0 target and the
bridge-only 562-byte record and IPC envelope. A cap may narrow but must never silently widen approved
authority (`AGENTS.md`, ADR-0023); this is a hard constraint on the generator's output, which becomes
the frozen cap only after ADR review, not automatically. Transport framing stays separate from
application-data budget — ADR-0029's caps are "aggregate application data budgets, not raw
Mach-message or XPC-serialization claims," and only S3's native harness is authorized to make
transport-level observations.

## Native/Go copy-ownership tests

Directly from ADR-0029's byte-ownership section and the IPC plan's S4 fault matrix. These are
Go-side tests runnable now against a test-double bridge (no real XPC needed — that's S3), because
the contract is fully specified independent of the native implementation:

- native → Go input copy discipline: mutating the native-side buffer after copy must not change what
  the Go decoder sees;
- Go → native output copy discipline: mutating a returned byte slice after a call returns must never
  mutate Supervisor-retained state;
- no native pointer retained by Go, no Go pointer retained by native — a synthetic violation (short
  write, oversize output, pointer/length mismatch) must terminate the process as a local integrity
  fault, not produce a soft refusal;
- fixed-cap output buffer never grows past its declared cap;
- caller mutation after dispatch and accessor/reply mutation after return cannot change retained
  bytes;
- one function per method, no opcode/generic-command entry point; and
- one dedicated core queue serializes store-affecting bridge calls (concurrency test, no interleaved
  partial writes).

None require M1 — they test the bridge contract using synthetic byte payloads of the correct
already-fixed widths.

## Wrong-role/purpose/audience/install/epoch refusal cases

All buildable now with synthetic fixture bytes; none require M1-specific values.

**Role/service:** message delivered on the wrong role's service (Broker-role service receiving
`RegisterPlanV0`/`RequestAttemptV0`, or daemon-role service receiving `GetRegisteredPlanV0`/
`SubmitApprovalV0`) → `AUTHENTICATION`, zero core calls; a valid-service message naming a method not
in that service's closed call set → `UNSUPPORTED`.

**Purpose/audience:** a request body attempting to smuggle a different purpose value than the one
the native layer derives from service+tag → `MALFORMED`/`SCHEMA` (caller values cannot override
purpose/audience per ADR-0029); mismatched audience field → `AUTHENTICATION`/`UNSUPPORTED`; local
channel purpose substituted for signed-object purpose/audience anywhere approval-adjacent → refused
(explicit ADR-0029 rule, worth its own regression test here since it's a common implementer
mistake).

**Installation/epoch:** request installation ID not equal to already-open Supervisor state →
`BINDING`, zero registration created; stale epoch sequence/digest → `BINDING` or transition
`TRUST_STATE`; `GetRegisteredPlanV0` fetch of a `RegistrationID` under a different
installation/epoch than currently active → `BINDING`, no mutation.

**562-byte role-binding record:** review count byte > 8 → `SCHEMA`/`MALFORMED`; nonzero value in an
unused review-digest slot beyond the declared count → rejected; nominal identities that don't
resolve against any locally trusted resolver → refused before state change, `BINDING`/`SCHEMA`
depending on which resolver fails; record length not exactly 562 bytes → `MALFORMED`/`SCHEMA`.

**Structural:** unknown/missing XPC dictionary key, wrong value type, extra file
descriptor/endpoint/Mach right, or nested dictionary/array in the closed message → refused before a
Go call, zero core calls.

Every case must additionally assert the ADR-0029 no-state rule: zero widened authority, zero
consumed approval, zero created attempt, zero lifecycle/backend effect, for every rejection prior to
a committed Go transaction.

## Response-loss and replay behavior

- `RegisterPlanV0` is deliberately non-idempotent: resubmitting the exact same plan/bindings/
  manifest/source after a simulated lost reply must create a second, independent registration — no
  deduplication key is invented by the transport, and both registrations must remain separately
  readable via `GetRegisteredPlanV0` until each expires on its own. This is the trap case an
  implementer is most likely to get wrong.
- `GetRegisteredPlanV0` is a repeatable read: two fetches of the same `RegistrationID` must return
  byte-identical copies with zero store mutation on either call.
- Response loss before dispatch has no effect (zero Go calls, zero state change).
- Response loss after dispatch is response-cancellation only, never an inferred abort — durable
  state must reflect the committed registration regardless of whether the client received the reply.
- An indeterminate store-commit outcome must surface or recover as `RECOVERY_REQUIRED`; no XPC cache
  or request-ID table overrides store truth.
- Request-ID reuse across separate connections/calls after the first completed is a wholly fresh
  transport call, never an implicit "same operation" signal.
- Concurrent identical requests converge (don't duplicate) for the two idempotent calls
  (`SubmitApprovalV0`/`RequestAttemptV0`) — worth cross-referencing here so a reader doesn't assume
  `RegisterPlanV0` behaves the same way; it explicitly does not.

All seven are pure Go/fixture-store tests against the already-implemented
`registrationstate.Component`, only newly wrapped with the four-call IPC-shaped entry points — none
require M1.

## Field-authority manifest additions

These additions are deferred to the post-R5B M2/S1 checkpoint. R1 owns a separate additive set of
role-specific Source Validator v1 request/result/process/profile/resource-policy classifications;
it must not pre-create these Supervisor IPC targets or treat its unmeasured resource fields as an
active policy.

Per ADR-0034's explicit instruction that the passive field-authority manifest must add canonical
targets for the source manifest and the method-specific registration/fetch projections "in the same
change as their candidate definitions... no parallel manifest or prose exception is allowed":

**New targets:** `capsule.source-manifest` v0 (blocked on M1's finalized CDDL and 87/95-byte
fixtures); `capsule.register-plan-v0-request` and `capsule.get-registered-plan-v0-reply` (partially
startable now for the plan-bytes/binding-record fields; blocked on M1 for the source-manifest/
source-bytes fields).

**New profiles** (all startable now, using only the existing closed enum vocabulary in
`field-authority-manifest.schema.json` — `authorityClass`, `originRole`, `validator`, `resolver`,
`authorityEffect`, `approvalVisibility`, `contentStatus`, `binding`, `unknownBehavior`):
`ipc-message-discriminator`, `ipc-correlation-only`, `role-binding-nominal`, `role-binding-resolved`,
`source-manifest-content`, `source-bytes-content`, `registration-lookup-key`.

`scripts/verify-field-authority-manifest.mjs` already enforces no duplicate target identity, no
duplicate field path, no duplicate `sourceField`, and every field's `profile` existing in
`profiles`. New target/profile pairs from this slice pass that verifier unchanged.

## Fixture generation and independent-verifier responsibilities

Modeled on the existing `docs/PHASE_2B_APPROVAL_ATTEMPT_BOUNDARY.md` pattern (manifest of
rules/cases/fixtures, exact case counts, `authorityStateChanged` assertions) and
`scripts/verify-cddl-fixtures.mjs`/`scripts/verify-field-authority-manifest.mjs` (schema-driven,
script-enforced, not hand-authored trust).

**Generator responsibilities:** read the canonical CDDL/Go-struct definitions from the manifest
(never re-derive a field width by hand); emit byte-exact known-answer fixtures for every call's
success reply, every refusal classification's minimal triggering case, every boundary and
cap-plus-one vector, and every copy-ownership/mutation-after-return case above; emit derived
aggregate caps with a self-check that fails the build if the generated aggregate disagrees with
current ADR prose, forcing an ADR revision rather than silent drift; cross-generate Go, TypeScript,
and (where applicable) Swift-readable fixtures from one shared source of truth, matching the
existing `execution-plan-v0.json`/`plan-registration-v0.json` pattern.

**Independent-verifier responsibilities:** decode every fixture with an implementation sharing no
code with the generator and confirm byte-for-byte/field-for-field agreement; independently re-run the
refusal and response-loss matrices above without importing the generator's classification logic;
confirm every field in the two new IPC targets has a semantically-correct manifest entry (the
verifier script already checks presence, not semantic correctness); confirm no fixture imports
experiment code into product packages or claims ADR-0019 acceptance; record exact tool/library
versions and distinguish "generated and independently verified" from "generated only" in retained
evidence.

This process can be designed and scaffolded now but cannot run to completion until M1 supplies the
`SourceManifest` encoder and its own fixture corpus.

## Dependency-aware implementation plan

**Can start now, no M1 dependency:**

1. Add the seven new field-authority profiles to `schemas/authority/field-authority-manifest.json`.
2. Define the 562-byte role-binding-record Go struct (bridge-only) with the exact field layout from
   ADR-0029's table — already closed independent of M1.
3. Add `capsule.register-plan-v0-request`/`capsule.get-registered-plan-v0-reply` manifest targets
   for the plan-bytes and binding-record fields only.
4. Write the native/Go copy-ownership tests against a synthetic test-double bridge.
5. Write the wrong-role/purpose/audience/install/epoch refusal fixtures (all cases above) using
   synthetic installation/epoch/registration IDs against the existing `registrationstate.Component`.
6. Write the response-loss/replay fixtures against the existing `Component.RegisterPlan`/
   `ResolveUsable`.
7. Design the read-only Broker fetch facade's Go signature (defensive-copy discipline only, not the
   full S2 implementation) using the plan-bytes/binding-record portion of the reply; revisit once
   source-manifest/source-bytes reply fields are added.
8. Scaffold the fixture generator script structure using the plan-bytes/binding-record/registration
   portions of the cap arithmetic (65,536-byte plan cap and 562-byte binding record are both already
   closed, independent of M1).

**Blocked on M1:**

| Item | Blocked on |
| --- | --- |
| `capsule.source-manifest` v0 target's canonical CDDL and 87/95-byte boundary fixtures | M1's finalized `SourceManifest` CBOR encoder/decoder and its fixture files |
| `source-manifest-content`/`source-bytes-content` field entries in the two new IPC targets | Same — cannot classify a field whose canonical source object has no `definition` yet |
| `0..262,144`-byte `main.mjs` source-bytes cap enforcement and its cap-plus-one fixtures | M1's finalized strict-UTF-8/BOM/newline/Unicode source validator and byte-budget vectors |
| The 328,337-byte request cap and 332,433-byte reply cap as fixture-generated (not prose-copied) values | M1's completed manifest/source fixtures, which the generator must consume as source of truth |
| Full completion of the two new IPC manifest targets (all fields present) | Same M1 outputs |
| Independent cross-language byte-equality confirmation for the two new IPC targets | Cannot run until the generator (blocked above) produces something to verify |
| Non-executing ECMAScript module-request validator refusal fixtures | M1-owned; this slice only needs to know they exist as a precondition for source-bytes acceptance — see the retained [`.mjs` module-request validation boundary review](MJS_VALIDATION_BOUNDARY_REVIEW.md) |

**Explicit non-goals for this slice** (already out of scope per the IPC plan, restated for
completeness): S2's full Go facade implementation (only the fetch-facade signature is in scope
here); S3's native XPC harness; any TypeScript plan-v1/approved-byte migration work (removed from
the first-release critical path by ADR-0034); any claim of production authenticated IPC, activated
endpoint, or admitted backend/runtime.

## Summary of every value/artifact blocked on M1

1. `SourceManifest` v0 canonical CDDL definition and file location.
2. 87-byte (zero-file) / 95-byte (max-file) source-manifest fixture bytes, fixture-generated and
   cross-verified — currently only prose in ADR-0034, not yet code-verified.
3. `main.mjs` byte-cap (0..262,144) enforcement code and its boundary/cap-plus-one fixtures.
4. The two new field-authority manifest targets' source-manifest/source-bytes field entries.
5. The 328,337-byte request cap and 332,433-byte reply cap, as fixture-generated values.
6. Independent cross-language verification of items 1–5.
7. The non-executing module-request validator's refusal fixtures (consumed here, owned by M1).

Everything else in this plan — field ownership assignments, the 562-byte binding record, refusal-
classification test cases, response-loss/replay test cases, copy-ownership tests, and the
fixture-generator/verifier process design — can proceed now, in parallel with M1.
