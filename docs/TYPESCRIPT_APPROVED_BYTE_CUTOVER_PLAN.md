# TypeScript approved-byte atomic cutover plan

Status: proposed implementation plan; passive Slice A conformance artifacts and the Slice B design
decision exist, but no owner/store implementation or consumer is active.

## Decision checkpoint

Proposed ADR-0030 defines the exact object family and originally records `OWNER-UNRESOLVED`.
[Proposed ADR-0032](adr/0032-select-enrolled-typescript-source-preparer.md) now selects a separate
unprivileged Source Preparer with an exact one-shot Node worker and one role-namespaced immutable
source store. The decision is not accepted or implemented. No consumer may activate until its
evidence gates pass and the coordinated cutover completes. SDK-produced transformations remain
untrusted hints; post-registration or post-approval transformation is always forbidden.

The retained
[P0 authority/TCB checkpoint](TYPESCRIPT_SOURCE_PREPARER_P0_AUTHORITY_REVIEW.md) places passive
Source Preparer P1 work on hold. The topology is conditionally retained as the least-dangerous
TypeScript option, but its role-namespaced store is not a same-user security boundary by itself and
Node/Amaro remains full planning/store TCB until proven confined. P1 cannot freeze APIs or bytes
until every checkpoint entry criterion is satisfied. User planning direction accepts a bounded
modern ESM `.mjs`-only JavaScript first-release fallback if a stop condition fires, with no
CommonJS, package resolution, legacy Node module surface, or governed-runtime widening. That exact
media/profile contract is not frozen here and requires an applicable ADR/contract update before
implementation.

## Dependency graph

```text
A. passive object family + known answers (complete)
  -> B. Proposed ADR-0032 owner/store design (selected; acceptance and implementation pending)
       -> B0. P0 authority/TCB platform-boundary closure (HOLD)
            -> B1. passive Source Preparer/store/field-authority fixtures
                 -> B2. fault-injected store + governed transformer/package/installed evidence
  -> C. final CDDL/media/profile review and Swift agreement
  -> D. immutable original/emitted/options/profile/record source-store ownership
  -> E. ExecutionPlan v1 builder + complete role resolver
  -> F. PR #65 copy-only IPC record/method versioning + registration/Broker migration
  -> G. approval/attempt/lifecycle projection migration
  -> H. transcript/receipt migration and full downstream fixture replacement
  -> I. one offline store migration + one active-version cutover
```

Runtime admission, module loading, and real-backend work depend on `I`; they are not part of this
graph and do not become valid because the graph completes.

## Atomic implementation slices

### Slice A — passive contract and conformance

- Retain the v0 approved-byte CDDL family and illustrative `ExecutionPlan` v1 CDDL separately from
  current v0 wrappers.
- Retain exact profile/options, authoring/executable manifests, record, ordered set, and plan-source
  binding projection known answers.
- Verify exact path ordering, defensive ownership, strict UTF-8/BOM/newline/Unicode distinctions,
  all caps and cap-plus-one, object mutations, and nominal digest-role confusion in TypeScript and
  independent Go fixture verification.
- Do not call the transformer, create a consumer, or change current product behavior.

### Slice B — owner and source-store decision

- Review and accept Proposed ADR-0032's separately enrolled, unprivileged Source Preparer; reject
  daemon, SDK, Broker, Supervisor, updater, generic parser/helper, runtime, and split-store routes.
- Bind exact Node/Amaro packaging, executable custody, updater responsibilities, process identity,
  authentication, timeout, cancellation, diagnostics, aggregate capacity, crash recovery, and
  cleanup.
- Define immutable storage ownership for exact original/emitted file bytes and exact profile,
  options, record, record-set, and manifest bytes. Every accessor returns defensive copies or
  read-only snapshots; no live host path becomes authority.
- Follow the exact dependency, conformance, and fault gates in
  [`TYPESCRIPT_SOURCE_PREPARER_PLAN.md`](TYPESCRIPT_SOURCE_PREPARER_PLAN.md). Slice B is not complete
  merely because the topology is written down.
- Before B1, prove one OS-enforced single-member store container, close worker process-tree and
  package/store authority, select the sealed store-genesis/update channel, settle archive/release
  and cancellation/death semantics, and design recursive nested-member field-authority coverage.

### Slice C — complete object-model implementation

- Finalize plan v1 CDDL and object-specific bounded decoders in Go, TypeScript, and Swift.
- Replace the single source role with nominal original/executable/record-set roles in builders and
  complete-role enforcement. Validate path/member/order/media/length/digest relationships and
  require every transformed file record and every pass-through relationship.
- Retain exact received canonical bytes; never decode/re-encode authoritative objects.
- Produce independent cross-language bytes and reject cross-object, cross-version, cross-domain,
  unknown-field, malformed-CBOR, and raw-budget cases.

### Slice D — authority-state migration

- Introduce a new top-level store/protocol version containing only plan v1 registrations.
- Update registration records, approval rendering snapshots, immutable attempt projections,
  lifecycle copied plan bindings and binding-set digest, transcripts, and receipts.
- Replace the 530-byte plan known answer and every downstream registration, approval, attempt,
  lifecycle, snapshot, and receipt fixture together.
- Implement an explicit offline, lock-held, validate-transform-write-sync-rename-sync-reopen
  migration. Old binaries refuse the new version. Failure or indeterminate rename enters the
  existing recovery/repair model; no fallback creates an empty or mixed store.

### PR #65 authenticated-IPC reconciliation

- Preserve PR #65's selected unprivileged per-user SMAppService Supervisor topology and synchronous
  method-specific copy-only native C/Objective-C-to-Go ABI. Native continues to own only live
  XPC/Security peer objects, flow/deadline accounting, and bounded copied buffers; Go continues to
  own authority and lifecycle state. Neither layer transforms TypeScript or becomes a generic
  helper.
- Treat the proposed 562-byte `RegisterPlanV0` submitted complete-role binding record as v0-only.
  Do not freeze or reinterpret it for the three new plan-v1 digest roles. Define a versioned
  registration method/binding record with a newly calculated inclusive cap, cap-plus-one refusal,
  defensive copies on both sides of the C ABI, and independent native/Go known answers.
- Reconcile `RequestAttemptV0`, `GetRegisteredPlanV0`, and `SubmitApprovalV0` only where the atomic
  plan-v1 projections require a version change. No optional transform fields, generic envelope,
  fallback method, or dual v0/v1 authority acceptance is allowed.

### Slice E — activation

- Activate one plan version only after full fixture agreement, Broker complete rendering, exact
  storage readback, owner/topology evidence, and archive/rollback implications are reviewed.
- Remove old active plan acceptance in the same release. No feature flag, retry path, or endpoint
  accepts both shapes as equivalent authority.
- Keep attempts registration-ID-only. Runtime staging reads only the executable manifest role.

## Migration invariants

- One active installation/store version, one plan version, and one complete copied-binding shape.
- No mixed `Job` extension, optional transform fields, dual decoding fallback, or field inference.
- No transformation after plan construction, registration, Broker rendering, approval, or attempt
  creation.
- Any original, emitted, media, profile, toolchain, options, record, order, map/URL disposition, or
  diagnostic change requires a new plan, registration, and approval.
- `ApprovalGrant` binds the exact plan/registration; it never duplicates transform fields.
- Only executable-manifest bytes are execute-eligible.

## Required cutover verification

- Exact/cap-plus-one for 32 files, 262,144 bytes per original/emitted file, and independent
  1,048,576-byte aggregates.
- ASCII path order/member/entrypoint checks; composed/decomposed Unicode, LF/CRLF, BOM, invalid
  UTF-8, and trailing-newline distinctions.
- Profile, Node, Amaro, archive, executable, options bytes/digest, input/output media, original,
  emitted, record order, source-map, source-URL, and diagnostic mutations.
- Cross-domain digest, cross-object, cross-version, missing/unknown field, canonical-CBOR, raw
  budget, and defensive-copy refusal.
- Exact copy-only native/Go IPC agreement for the new complete-role binding record, including the
  replacement for the v0-only 562-byte cap/known answer and cap-plus-one refusal.
- Go/TypeScript/Swift known-answer agreement; independent Go/TypeScript hashing; fixture generator
  idempotence; Broker and Supervisor exact-byte readback.
- Every old/new store and process-death boundary, including no mixed snapshot and old-binary refusal.

## Explicit non-claims

Completion of this plan would define the authority contract. It would not prove semantic
equivalence, type correctness, Node API stability, cross-platform stability, an independent
transformer, governed runtime admission, external isolation, module loading correctness, or
production readiness. Each requires its own retained mechanism and evidence.
