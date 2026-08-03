# ADR-0030: Define the passive TypeScript approved-byte migration boundary

- Status: Proposed
- Date: 2026-08-03
- Refines if accepted: ADR-0011, ADR-0019, ADR-0023, ADR-0024, ADR-0025, and ADR-0026

## Context

ADR-0026 and the retained TypeScript approved-byte experiment establish the ordering rule: exact
Node 22.22.1/Amaro 1.1.5 strip-only emission can finish before plan construction, and the emitted
JavaScript—not an original-only TypeScript digest—can become the bytes later authorized for
execution. They intentionally leave the production owner and coordinated protocol migration open.

The current passive `ExecutionPlan` v0 has one `sourceManifestDigest`, the current registration and
lifecycle fixtures copy that role, and the current mixed `Job` remains the scaffold API. Extending
any one of those surfaces would create dual or partial authority. This ADR defines the complete
nominal object family and the atomic cutover boundary while retaining only unwired Slice A
fixtures and verification.

## Proposed decision

### Object family and media types

The approved-byte object family is version 0 and uses exact deterministic CBOR under ADR-0019's
still-Proposed profile. Each object has a closed CDDL map, `objectVersion == 0`, and an exact media
type:

| Object | Exact media type |
| --- | --- |
| Original authoring manifest | `application/capsule.original-authoring-source-manifest+cbor;v=0` |
| Executable JavaScript manifest | `application/capsule.executable-javascript-source-manifest+cbor;v=0` |
| Transformer profile | `application/capsule.typescript-transformer-profile+cbor;v=0` |
| Normalized options | `application/capsule.typescript-normalized-options+cbor;v=0` |
| Per-file transformation record | `application/capsule.typescript-transformation-record+cbor;v=0` |
| Ordered record set | `application/capsule.typescript-transformation-record-set+cbor;v=0` |

The Slice A `typescript-plan-source-bindings` object is a fixture-only projection of the three
future plan fields. It is not a registration or executable plan. A coordinated plan cutover uses
`application/capsule.execution-plan+cbor;v=1`; current v0 remains unchanged until that cutover.

All object digests are SHA-256 over exact received canonical CBOR bytes. File-content digests are
SHA-256 over exact file bytes. Original file, emitted file, original manifest, executable
manifest, profile, options, record-set, and execution-plan digests are separate nominal roles.
Equal-width or accidentally equal values are never interchangeable without role-resolved context.

### Original and executable manifests

Both manifests bind:

- the exact canonical ASCII entrypoint;
- one through 32 unique entries sorted by case-sensitive unsigned ASCII path bytes;
- for each entry, logical path, exact media type, byte length, and content SHA-256; and
- the exact aggregate byte length.

Paths retain ADR-0023's relative ASCII grammar and 256-byte path/64-byte segment ceilings. The
original manifest accepts exact TypeScript or JavaScript authoring bytes. TypeScript is restricted
to erasable ESM `.ts`/`.mts`; `.tsx`, JSX, `.cts`, decorators, enums, namespaces, parameter
properties, and other generation-requiring syntax refuse. JavaScript is exact pass-through. The
executable manifest accepts only JavaScript media type. Original and executable manifests have the
same logical paths and entrypoint. A suffix never determines executable syntax.

Each file is at most 262,144 bytes. Original and executable aggregates are independently at most
1,048,576 bytes. Exact caps accept; cap-plus-one refuses without truncation, clamping, resizing, or
partial record creation. Strict UTF-8 without BOM is required. Unicode scalars, normalization form,
LF/CRLF choice, and trailing newline state are byte identity and are never rewritten.

### Profile and normalized options

The only v0 transformer profile is:

| Field | Exact value |
| --- | --- |
| API | `node:module.stripTypeScriptTypes` |
| Node | `22.22.1` |
| Amaro | `1.1.5` |
| platform/architecture | `darwin` / `arm64` |
| official Node source archive SHA-256 | `87104b07e7acee748bcc5391e1bc69cf3571caa0fdfb8b1d6b5fd3f9599b7849` |
| distribution archive SHA-256 | `261da057fb25ff2912dd6abb7842fc915ddf7947a2cb3c8cce90875d2b9bb667` |
| installed executable SHA-256 | `245e0321af97d3c21dd4e7104457334dfe3c3ba7982d0db75363e354565f8cbb` |

Normalized options are the exact closed CBOR values: mode `strip`, diagnostic policy
`reject-any`, source map `absent`, source URL `absent`, TypeScript input media type, and JavaScript
output media type. Unknown fields or values refuse. Platform or executable changes create a new
profile; Node's Active Development API is not claimed stable across versions or platforms.

### Transformation record and ordered set

Every transformed TypeScript file has one record binding:

- its logical path;
- original media type, length, and role-specific content digest;
- emitted media type, length, and role-specific content digest;
- transformer-profile digest plus exact Node, Amaro, source archive, distribution archive, and
  executable identities;
- the complete exact normalized-options CBOR bytes and their role-specific digest;
- source-map and source-URL dispositions, both exactly `absent`; and
- diagnostic policy `reject-any` with successful count exactly zero.

A rejected transformation creates no successful record. JavaScript pass-through entries create no
transformation record. The record-set object contains zero through 32 exact canonical record byte
strings ordered by their decoded unsigned-ASCII logical paths. Its digest is SHA-256 over the exact
record-set CBOR, so record order, membership, or any record mutation changes the plan binding.

### Plan, registration, approval, and execution eligibility

The future `ExecutionPlan` v1 owns three distinct fields:

1. original-authoring source-manifest digest;
2. executable-JavaScript source-manifest digest; and
3. ordered transformation-record-set digest.

The plan also binds both manifest aggregate lengths. Role resolution requires exact immutable
manifest, record, options, and content bytes to exist in the owning source store and requires all
paths, entrypoint, lengths, digests, media types, records, and pass-through relationships to agree.

Only bytes named by the executable JavaScript manifest are execute-eligible. Original authoring
bytes, transformer profile, options, records, or an original-only digest can never be sent to a
runtime or used to regenerate executable bytes after plan construction. Registration, attempt, and
execute operations remain registration-ID-only and accept no replacement bytes or transform
inputs.

`PlanRegistration` continues to retain and bind exact plan bytes and their digest. `ApprovalGrant`
continues to bind that exact registered plan rather than duplicating transformation fields. The
Broker must validate and completely render the registered plan's original/executable/profile
roles before signing. Any mutation requires a new plan, registration, and approval.

### Transformation owner remains an activation blocker

No existing role can be assigned production transformation without another accepted decision:

- An SDK or proposal builder may precompute bytes for ergonomics, but it is agent-controlled and
  cannot establish trusted transformation or plan authority. Its output must be treated as an
  untrusted proposal.
- The Go daemon is the architectural planning and source-handling owner, but invoking a Node
  helper would violate the daemon-to-helper prohibition, and embedding a JavaScript parser/runtime
  changes its review surface and packaging. Neither topology is presently authorized.
- The Broker owns approval and content custody, not general source transformation.
- The Supervisor must not gain a parser/transformer responsibility without a separate ADR, and a
  daemon-to-Supervisor transform request would widen the security-sensitive boundary.
- The updater/trust verifier and runtime adapter do not own proposal construction, and transformation
  inside the runtime is post-approval and forbidden.

Therefore Slice A records an explicit `OWNER-UNRESOLVED` activation blocker. It does not silently
assign transformation to the daemon, Supervisor, Broker, SDK, helper, or runtime. A later ADR must
select a process, packaging, authentication, custody, fault, and update topology consistent with
Capsule's authority split before any consumer activates.

Follow-up Proposed ADR-0032 selects a separately enrolled Source Preparer and immutable source
store. That later design does not remove this Slice A blocker until it is accepted, implemented,
and validated through the coordinated cutover.

The parallel authenticated-IPC decision in PR #65 does not resolve this blocker. Its selected
unprivileged per-user SMAppService Supervisor has a small native C/Objective-C XPC/Security front
end linked in-process to the existing Go authority/lifecycle core through a synchronous,
method-specific, copy-only C ABI. Native owns live peer objects, flow/deadline accounting, and
bounded copied buffers; Go owns registration, approval, attempt, lifecycle, recovery, and the
store. Neither side of that seam may transform TypeScript or turn the IPC layer into a generic
helper.

### Atomic cutover

The migration is one versioned authority cutover, not compatible dual acceptance. It must update
CDDL/schema candidates, TypeScript/Go/Swift views, source ownership/builders, complete-role
enforcement, plan construction, registration validation/storage, Broker rendering, approval and
attempt projections, lifecycle copied bindings/set digests, transcripts, receipts, and every
downstream known answer in one reviewed change. Old plan v0 and new plan v1 may be decoded only by
explicit version-specific offline migration tooling; no active consumer treats them as equivalent
authority and no active store contains a mixed version.

PR #65's current `RegisterPlanV0` submitted complete-role binding record is a 562-byte proposed v0
record. It remains the oracle only for the current plan v0 seam. The approved-byte migration adds
original-manifest, executable-manifest, and transformation-record-set roles, so it must version the
registration method and complete binding record together (for example, `RegisterPlanV1`) and
replace its exact cap and known answer. It must not freeze, extend, or reinterpret the 562-byte v0
record as the v1 shape. The other selected role-specific methods—`RequestAttemptV0`,
`GetRegisteredPlanV0`, and `SubmitApprovalV0`—likewise require an explicit atomic version/copy-shape
review where their typed projections change; no generic IPC fallback is permitted.

The exact dependency and verification gates are retained in
[`TYPESCRIPT_APPROVED_BYTE_CUTOVER_PLAN.md`](../TYPESCRIPT_APPROVED_BYTE_CUTOVER_PLAN.md).

## Passive Slice A evidence

Slice A adds only candidate CDDL, nominal TypeScript/Go fixture views, deterministic generation,
nine exact known-answer files, fourteen refusal mutations, and bounded local fixture verification.
It wires no product consumer. The retained ordinary source and emitted output remain the exact
391-byte Node 22.22.1/Amaro 1.1.5 experiment result.

## Consequences and limitations

- Exact executable-byte authority is representable without post-approval transformation.
- The current plan/schema/type/registration/approval/lifecycle behavior remains unchanged.
- No evidence here proves semantic equivalence, TypeScript correctness, cross-version or
  cross-platform output stability, an independent transformer implementation, runtime admission,
  module loading, product process isolation, production readiness, or user comprehension.
- ADR-0019, acceptance and evidence for Proposed ADR-0032's owner/topology, governed packaging,
  Broker/Supervisor implementations, runtime admission, and the complete atomic cutover all remain
  blockers.
