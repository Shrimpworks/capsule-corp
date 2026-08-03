# ADR-0026: Bind strip-only TypeScript emission before plan registration

- Status: Proposed
- Date: 2026-08-03
- Refines if accepted: ADR-0003, ADR-0011, ADR-0017, ADR-0023, ADR-0024, and ADR-0025

## Context

The governed `deno_core` physical-omission experiment removed the native built-in registration
blocker for one exact construction, but `deno_core` accepts JavaScript rather than TypeScript. If
Capsule approved only an original TypeScript digest and transformed it later, the executable bytes
could change with transformer version, options, diagnostics, source-map behavior, or mutation after
approval. That would violate ADR-0011's exact registered-byte authority.

The retained Deno-family evidence identified `deno_ast` 0.53.3 with `transpiling` as the compatible
Deno path. Its exact locked marker resolves 180 packages including itself. The bounded
[TypeScript approved-byte experiment](../../experiments/typescript-approved-byte-boundary/RESULTS.md)
compared that path with official TypeScript 6.0.3 and exact Node 22.22.1's built-in
`node:module.stripTypeScriptTypes` backed by Amaro 1.1.5.

For fixed benign fixtures, Node strip-only emission was byte-identical across 20 same-process and
three separate-process runs, preserved Unicode and line endings without normalization, enforced
the selected exact caps, refused malformed and transform-requiring TypeScript, and detected
deliberate source/output/options/transformer/disposition mutations. An independent Go verifier
agreed on the exact byte identities and closed record. Node marks the API Active Development and
warns that output is not stable across Node versions, so a floating Node identity is not safe.

This ADR proposes a contract and ordering decision only. It does not select governed `deno_core`,
admit `RUNTIME-001`, choose a production transformation-process topology, accept ADR-0019, or
change current product schemas/types.

## Proposed decision

### Narrow transformation profile

The first candidate accepts only erasable ESM TypeScript under this exact closed profile:

| Field | Exact value |
| --- | --- |
| transformer API | `node:module.stripTypeScriptTypes` |
| Node | `22.22.1` |
| bundled parser/stripper | Amaro `1.1.5` |
| mode | `strip` |
| source map | explicitly `absent` |
| source URL | explicitly `absent` |
| diagnostics | `reject-any`; successful count exactly zero |
| input media type | `application/capsule.typescript-source;v=0;module=esm` |
| output media type | `application/capsule.javascript-source;v=0;module=esm` |

The transformation profile additionally binds the exact official Node source-archive digest,
platform distribution-archive digest, and installed executable digest. A platform with different
distribution or executable bytes is a different profile until cross-platform output evidence
supports a separately reviewed identity rule.

`.ts` and `.mts` may use this candidate only when every syntax form is erasable. `.tsx`, JSX,
CommonJS-specific `.cts`, decorators, enums, namespaces, parameter properties, and any other form
requiring JavaScript generation refuse. Exact JavaScript source remains pass-through and is still
bound into the executable manifest. Unknown media types, options, fields, diagnostics, source-map
states, versions, or toolchain identities refuse.

The inclusive limits are 32 source files, 262,144 bytes per original TypeScript file, 1,048,576
aggregate original bytes, 262,144 bytes per emitted JavaScript file, and 1,048,576 aggregate
emitted bytes. Original and emitted caps are separately enforced. No path, stream, or buffer grows,
truncates, or clamps at cap-plus-one.

### Ordering and byte ownership

Transformation occurs after strict raw/schema decoding and semantic source resolution but before
executable-source manifest construction, `ExecutionPlan` construction, plan registration, Broker
rendering, or approval. Transformation after plan registration or approval is forbidden.

The same validated logical source path identifies the authoring and executable roles; the emitted
bytes are JavaScript regardless of a retained `.ts`/`.mts` logical suffix. This avoids a path
rewrite and collision surface. A later runtime module loader must use the registered output media
type and exact executable manifest, not infer executable syntax from a path suffix.

The future plan binds:

1. an original-authoring source-manifest digest over exact original TypeScript and pass-through
   JavaScript bytes;
2. an executable source-manifest digest over exact emitted and pass-through JavaScript bytes; and
3. an ordered transformation-record-set digest.

Each closed per-file transformation record binds logical source path, original digest/length/media
type, emitted digest/length/media type, transformation-profile digest and exact toolchain
identities, normalized-options bytes/digest, source-map disposition `absent`, and diagnostic policy
`reject-any` with count zero. Original and emitted digests are distinct nominal roles even if the
32-byte values happen to match.

The exact emitted JavaScript object is the only source role eligible for later runtime delivery.
The original TypeScript object is retained for approval rendering, audit, and reproduction but
never becomes execute-time authority. The runtime receives no transformer, options, original-only
digest, or ability to regenerate JavaScript.

### Registration and approval binding

`ExecutionPlan` owns the detailed transformation bindings. `PlanRegistration` continues to bind
the exact plan digest and store the exact plan bytes. `ApprovalGrant` continues to bind the plan
digest and registration rather than duplicating every transformation field. The Broker fetches the
registered plan from the Supervisor, validates the role-specific manifests/records, and renders
both authoring and executable identities plus the exact transformation profile before signing.

This preserves one authoritative typed plan and avoids a second signed projection that could
disagree. A change to original bytes, emitted bytes, media type, limits, toolchain, options,
diagnostic policy, or source-map disposition produces a different plan and requires a new
registration and approval.

The Supervisor does not accept original TypeScript plus a request to transform. Execute and
attempt APIs remain registration-ID-only and never accept replacement source, emitted bytes,
transformer identity, or options.

### Transformation owner and trust claim

This experiment establishes ordering and exact-byte ownership, not a production process owner.
The existing Go daemon architecture, AGENTS.md prohibition on daemon-to-helper shortcuts, and
deferred Supervisor language/privilege topology make an implicit helper or new Supervisor parser
responsibility unacceptable. The coordinated implementation proposal must name an owner and trust
boundary in a separate accepted decision before consumer activation.

The approval claim is exact byte authorization, not semantic equivalence. The evidence does not
prove that emitted JavaScript faithfully expresses the user's intended TypeScript, that a
compromised planner behaved correctly, or that the user understood either representation.

## Coordinated migration plan

The current pre-freeze mixed `Job`, candidate schemas, TypeScript types, 530-byte plan known answer,
and public contracts remain unchanged by this experiment. A later versioned migration must update
the complete set atomically:

1. define candidate CDDL and role-specific scalar types for original/executable source manifests,
   transformation profiles, records, and record sets;
2. extend the candidate `ExecutionPlan` and its Broker view with the three role-specific digests;
3. retain immutable exact original/emitted bytes and record bytes in the owning source store;
4. update Go, TypeScript, and Swift decoded views, plan construction, complete role bindings,
   registration validation/storage, approval rendering, attempt projection, durable lifecycle
   copied bindings, enforcement transcripts, and receipts;
5. add byte-exact cross-language known answers plus max/cap-plus-one, invalid UTF-8/BOM,
   Unicode/line-ending, diagnostics, unknown-option, wrong-media, wrong-domain, mutation, and
   explicit-source-map-absence fixtures; and
6. replace all downstream plan/registration/approval/attempt fixtures in one reviewed protocol
   version rather than extending or adapting the deprecated mixed `Job` authority model.

No consumer cutover may accept both old and new plan shapes as equivalent authority.

## Alternatives considered

### `deno_ast` 0.53.3 with `transpiling`

Rejected as the narrowest first candidate. It is the exact Deno-family path and was deterministic
for the fixed probe, but its marker resolves 179 dependencies and it transforms a broader language,
including the fixed enum. It remains a valid comparison point if the strip-only syntax is too
narrow and a later decision accepts the larger pre-approval review surface.

### TypeScript 6.0.3 `transpileModule`

Rejected as the narrowest first candidate. The official npm package declares no package
dependencies and is far smaller as a package graph than `deno_ast`, but it exposes a full compiler
and broader transformation/options surface. Its isolated emit also does not perform complete type
checking. No evidence required its broader syntax after strip-only passed.

### Transform inside the runtime after approval

Rejected. It puts the transformer graph into the live runtime TCB and makes executable bytes depend
on post-approval toolchain/options behavior. Binding only the TypeScript digest does not authorize
the emitted JavaScript bytes.

### Bind only emitted JavaScript

Rejected for TypeScript product semantics and approval rendering. It preserves execution-byte
identity but loses the authoring bytes and transformation identity the user was told they supplied.

### Duplicate every field in `ApprovalGrant`

Rejected. The grant already binds the exact registered plan digest. Duplication introduces a
cross-object consistency surface without creating stronger byte authority.

## Consequences and acceptance blockers

- Exact executable-byte authorization can be preserved without a transformer in the live runtime.
- The first TypeScript subset is deliberately narrow and rejects common transform-requiring syntax.
- Every transformer/toolchain/options change invalidates prior emitted bytes and requires a new
  plan, registration, and approval.
- Node's Active Development API, platform-specific executable identity, production owner/topology,
  independently reconstructible packaging/provenance, broader adversarial corpus, coordinated
  protocol migration, and Broker rendering all remain blockers.
- Accepting this ADR would not admit governed `deno_core`; module loading, packaging, restoration
  review, external isolation composition, and complete runtime-profile admission remain open.

This ADR may become Accepted only after the owner/topology decision and coordinated object-model
fixtures prove that the Broker and Supervisor retain the same exact original, emitted, options,
profile, and plan bytes without any post-approval transformation path.
