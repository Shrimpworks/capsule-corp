# ADR-0023: Bound protocol decoding and registration semantics

- Status: Proposed
- Date: 2026-08-01
- Refines if accepted: ADR-0009, ADR-0011, ADR-0017, and the object profiles proposed by ADR-0019

## Context

Phase 2A introduced passive candidates for the first `JobProposal`, minimum `ExecutionPlan`, and
`PlanRegistration`. The independent reviews agreed on their authority separation but identified
several values that must become exact before strict decoders, semantic planners, canonical wrappers,
or durable registration tests can share conformance fixtures.

The unresolved values fall into two different evidence domains:

- protocol parsing, canonical identity, and registration state need conservative local resource and
  replay rules before product implementation; and
- guest transport, output payload, runtime, and backend limits still require the existing P0
  measurement and admission campaigns.

This ADR proposes only the first domain. It neither completes a launch-capable `ExecutionPlan` nor
admits a backend, runtime bundle, public endpoint, approval path, or hostile guest.

## Proposed decision

### Boundary ownership and media types

Each object has one object-specific entry point. A caller cannot select a generic JSON/CBOR mode and
then supply an object discriminator that changes the validator.

| Object | Wire profile | Conformance media type | First owner |
| --- | --- | --- | --- |
| `JobProposal` | strict UTF-8 JSON, not canonical JSON | `application/capsule.job-proposal+json;v=0` | public raw decoder |
| `ExecutionPlan` | deterministic CBOR payload | `application/capsule.execution-plan+cbor;v=0` | Supervisor plan decoder |
| `PlanRegistration` | deterministic CBOR payload returned by typed authenticated IPC | `application/capsule.plan-registration+cbor;v=0` | Supervisor IPC method/wrapper |

Type/subtype comparison is case-insensitive after syntactic parsing. The version parameter must be
exactly `v=0`; aliases and additional parameters are rejected. Typed local IPC binds the method and
peer role independently of the conformance media-type label. `PlanRegistration` bytes alone never
grant registration authority and are not an independently signed portable object.

### Strict `JobProposal` raw profile

The raw decoder operates on caller-owned bytes before ordinary JSON parsing or schema allocation.
The first candidate uses these exact denial-of-service budgets:

| Dimension | Maximum |
| --- | ---: |
| Raw JSON document | 2,097,152 bytes |
| Nesting depth, counting the root | 32 |
| Total JSON value nodes | 16,384 |
| Total object members | 4,096 |
| Members in one object | 256 |
| Total array elements | 4,096 |
| Elements in one array | 256 |
| Total decoded UTF-8 text across keys and values | 1,572,864 bytes |
| One decoded object key | 128 bytes |
| Source file entries | 32 |
| One source file after strict UTF-8 encoding | 262,144 bytes |
| Aggregate source file bytes | 1,048,576 bytes |
| Canonical inline-input JSON bytes | 262,144 bytes |
| Non-source JSON string | 65,536 bytes |
| Labels | 8 entries; key 32 ASCII bytes; value 128 printable-ASCII bytes |

Every maximum is inclusive. A cap-plus-one input is drained or rejected without resizing the
budget. These are protocol/parser budgets for the narrow first slice, not guest-port, completion-
frame, output-payload, RAM, CPU, scratch, or backend limits.

The decoder additionally:

- accepts exactly one JSON document and rejects an empty body, UTF-8 BOM, invalid/overlong UTF-8,
  unpaired surrogate escape, trailing non-whitespace data, and a second document;
- rejects duplicate keys after JSON escape decoding, including `kind` plus `k\u0069nd`;
- rejects a value as soon as any raw, depth, node, member, element, text, or object-specific budget
  is exceeded;
- accepts insignificant JSON whitespace and arbitrary object-key order because the public request
  is not canonical JSON;
- accepts numeric tokens only in the grammar `0`, `[1-9][0-9]*`, or `-[1-9][0-9]*`, with the value
  between `-9007199254740991` and `9007199254740991`; fractions, exponents, leading plus, leading
  zeroes, and negative zero are rejected; and
- reports unknown fields/versions after raw decoding to the closed object schema rather than
  ignoring them.

Raw limits count UTF-8 bytes, not JavaScript UTF-16 code units, Swift grapheme clusters, Unicode
scalar count, or JSON escape spelling.

### Source path and source-manifest identity

The first slice keeps the Phase 2A ASCII path grammar:

- `/` is the only separator;
- paths are relative, 1 through 256 ASCII bytes, with 1 through 64 bytes per segment;
- a segment starts with an ASCII alphanumeric and thereafter contains only ASCII alphanumeric,
  `.`, `_`, or `-`;
- empty, `.`, `..`, leading/trailing separator, repeated separator, backslash, NUL/control, Unicode,
  and case-folded aliases are rejected;
- path comparison and ordering are case-sensitive unsigned ASCII-byte comparison; and
- the entrypoint must exactly equal one source-file key and end in `.js`, `.mjs`, `.cjs`, `.ts`,
  `.mts`, or `.cts`.

Decoded source content is encoded as strict UTF-8 without Unicode normalization, newline rewriting,
or trailing-newline insertion. JSON escape spelling is not source identity; the decoded Unicode
scalar sequence is. The source aggregate byte count is the sum of those UTF-8 file bytes.

The future `SourceManifest` candidate uses a closed deterministic-CBOR map:

```text
1: "capsule.source-manifest"
2: 0
3: validated entrypoint bytes
4: array sorted by source-path bytes of [path, SHA-256(content bytes), byte length]
5: aggregate source byte length
```

The `SourceManifestDigest` is SHA-256 over the exact deterministic-CBOR manifest bytes. Source bytes
remain separate immutable content; the manifest neither embeds host paths nor authorizes content
retrieval.

### Canonical inline-JSON identity

Public JSON remains noncanonical, but its decoded inline-input value is serialized to one bounded
canonical JSON byte representation for hashing, storage, and later transport:

- UTF-8 without BOM or insignificant whitespace;
- literals are exactly `null`, `true`, and `false`;
- safe integers use shortest base-10 form, with zero encoded as `0`;
- arrays preserve order and use only `,`, `[` and `]` delimiters;
- object keys use the existing ASCII candidate grammar and sort by unsigned ASCII bytes;
- object delimiters are `{`, `}`, `:`, and `,` with no whitespace;
- strings preserve the decoded Unicode scalar sequence without normalization; `"` and `\` use
  two-character escapes, U+0000 through U+001F use lowercase `\u00xx`, `/` is not escaped, and all
  other scalars are emitted as UTF-8; and
- the complete canonical bytes must not exceed 262,144 bytes.

`InlineInputDigest` is SHA-256 over those exact canonical JSON bytes. Its nominal role plus the
fixed `primary-data` plan field supplies domain separation; a generic digest is never accepted in
its place. `inlineInputByteLength` is the canonical JSON byte length. Implementations must produce
the same bytes directly and may not use platform `JSON.stringify`/`JSONSerialization` behavior as
the contract.

### Internal CBOR predecoder budgets

Before object decoding, internal wrappers apply deterministic-CBOR rules from proposed ADR-0019 and
these object-specific raw budgets:

| Object | Raw bytes | Nesting depth | Total data items | Maximum map entries | Maximum array elements |
| --- | ---: | ---: | ---: | ---: | ---: |
| `ExecutionPlan` | 65,536 | 8 | 256 | 64 | 8 |
| `PlanRegistration` | 4,096 | 4 | 32 | 16 | 0 |

Definite lengths, preferred integer/length encodings, deterministic map order, strict UTF-8, closed
integer labels, and the CDDL field bounds are mandatory. Indefinite containers, floats, bignums,
arbitrary tags, duplicate keys, trailing data, and decoded string/byte totals above the raw cap are
rejected before ordinary object allocation. Exact received plan bytes remain authoritative.

### Scalar zero rules

- Installation, registration, and Supervisor IDs are distinct 16-byte roles and reject the all-zero
  value. IDs are never reused within an installation.
- Attempt/approval IDs and nonces follow the same nonzero rule when their contracts freeze.
- SHA-256 digest fields accept every 32-byte bit pattern structurally, including all zeroes. A
  digest becomes authoritative only after the role-specific resolver finds matching bytes/record;
  an unresolved all-zero digest therefore fails binding like any other unknown digest.
- Internal CBOR uses byte strings. Public textual ID/digest projections remain outside this slice
  because `JobProposal` carries neither; later public summary/receipt contracts must define them
  explicitly.

### Registration creation, replay, sequence, and expiry

The caller submits only exact candidate `ExecutionPlan` bytes over authenticated daemon-to-
Supervisor IPC. The Supervisor independently validates and copies those bytes before any durable
authority change.

- Every successful call creates a new registration, even for identical plan bytes. There is no
  digest-based deduplication or caller-controlled idempotency key.
- The Supervisor generates a fresh nonzero registration ID and never accepts a caller-supplied
  `PlanRegistration` object.
- `registrationSequence` starts at 1 and increases monotonically across the installation, including
  trust-epoch transitions. It never resets or decreases. Exhausting `UInt53` enters
  `repair-required` and refuses registration.
- Counter advancement and the complete stored registration commit atomically. Validation,
  capacity, identifier-generation, or transaction failures create no registration and consume no
  sequence. Consumers do not treat contiguity as authority, so recovery-visible gaps remain safe.
- The trusted first-slice plan lifetime is 300 seconds. Registration accepts a plan only when
  `effectiveNow < plan.expiresAt <= effectiveNow + 300`; it copies `plan.expiresAt` unchanged to the
  registration. The user cannot request or extend this lifetime.
- A registration is usable only while `effectiveNow < expiresAt`; equality is expired. Approval
  expiry cannot exceed registration expiry. Trust fencing, epoch mismatch, quarantine, or repair
  invalidates use even before the timestamp.
- `effectiveNow` never moves backward: it is the maximum of the current trusted clock observation
  and the Supervisor's durable time high-water mark. A backward wall-clock adjustment cannot
  reactivate an expired registration; a forward jump may expire it early and fails safe.
- The first implementation admits at most 256 unexpired registrations and 4,096 stored registration
  records per installation. Capacity exhaustion returns `CAPACITY`; it never evicts an active,
  approved, attempted, unresolved, quarantined, or sole explanatory record.
- Automatic deletion is disabled in the first implementation. Later retention policy may prune
  exact bytes only after all related recovery/evidence obligations are terminal, while retaining a
  tombstone sufficient to prevent ID/sequence reuse and explain replay denial.
- The daemon cannot clear, reset, resurrect, or garbage-collect Supervisor registration state.

The 4,096-record ceiling is therefore an intentional fail-closed limit for the unwired first
implementation, not a viable continuous-service retention policy. Consumer activation requires a
separately reviewed Supervisor-owned archival/compaction design, retained power-loss evidence, and
an exact rule for which terminal records may leave the active store. Increasing the ceiling alone
does not satisfy that requirement.

The stored Supervisor record contains the wire registration, exact plan bytes, their recomputed
digest, registration time, and recovery metadata. Exact plan bytes do not appear inside the wire
`PlanRegistration` payload.

### Conformance classifications

The first corpus uses a fixed internal oracle vocabulary. These classifications identify the first
owning boundary and are not yet the smaller public `ErrorCode` or agent-summary vocabulary.

| Classification | First meaning |
| --- | --- |
| `MALFORMED` | Invalid raw syntax/encoding/media type/bounds, duplicate data, trailing bytes, or noncanonical CBOR |
| `UNSUPPORTED` | Unknown version/field/power/tag or intentionally unavailable protocol feature |
| `SCHEMA` | Known object with a missing field, wrong type/width, or invalid scalar shape |
| `SEMANTIC` | Structurally valid but contradictory source, path, slot, or value relationship |
| `POLICY` | Valid request denied by trusted policy/default/ceiling rules |
| `BINDING` | Referenced bytes/state do not match installation, epoch, registration, plan, profile, or policy |
| `CONTROL` | An exact required control cannot be enforced by the selected admitted profile/backend |
| `DOMAIN` | A valid object, ID, digest, slot, media type, or purpose is used in the wrong semantic role |
| `AUTHENTICATION` | Local channel peer/session/purpose is not authorized for the operation |
| `STALE` | Expired, replayed, superseded, fenced, or sequence-invalid state |
| `CAPACITY` | A fixed local object/store budget is exhausted without evicting authority/evidence |

A rejected proposal or registration case asserts `stateChanged: false` for every authority-bearing
store. Logging or a monotonic clock high-water observation is not authority creation. Human-readable
messages never carry user content, host paths, approval text, or arbitrary guest strings.

## Alternatives considered

### Reuse RFC 8785/JCS for inline JSON

Rejected for this slice. Gate A found a real cross-language number-representation disagreement.
The narrower safe-integer value set permits a smaller explicitly specified canonical JSON encoder
whose bytes are also directly consumable by the guest.

### Canonicalize Unicode source paths

Rejected for the first slice. ASCII paths remove normalization, separator, and case-fold ambiguity
without requiring platform filesystem semantics. Capsule never materializes these names as live
host paths.

### Deduplicate registration by plan digest

Rejected. Digest equality does not imply retry intent, approval equivalence, or lifecycle identity.
Fresh registration records preserve simple one-use reasoning at the cost of bounded unused records.

### Reject all-zero SHA-256 values structurally

Rejected. Width and nominal role are structural; existence and byte equality are binding checks.
Special-casing one digest value adds no meaningful authority defense.

### Freeze backend transport/output values here

Rejected. Parser budgets do not prove virtio-console framing, launcher behavior, output validation,
or backend enforcement. Those exact values remain owned by the documented P0 campaign.

## Consequences

- Task 2 can build byte-exact raw, schema, semantic, registration, and cross-domain fixtures with an
  exact owner and oracle.
- Candidate schema limits that count characters rather than UTF-8 bytes remain only structural
  backstops; the raw/semantic validator owns these exact byte budgets.
- Implementations need streaming or tokenizing predecoders rather than ordinary `JSON.parse` or a
  general-purpose CBOR decode followed by size checks.
- The restricted canonical inline-JSON encoder becomes small enough for independent Go, Swift, and
  TypeScript implementations and fuzzing.
- Registration retries may create more than one record; fixed capacity and no implicit eviction
  keep the failure explicit.
- The first store eventually refuses all new registrations if its terminal records are never
  archived. That is acceptable for unwired conformance/fake-backend work and blocks activation until
  bounded Supervisor-owned retention is implemented and validated.
- A future increase to a public or canonical-on-wire maximum requires a reviewed contract/version
  change and cap-plus-one fixtures. Internal operational limits may narrow but never silently widen
  approved authority.
- ADR-0019 remains Proposed. This ADR does not satisfy its production dependency, wrapper review,
  fuzzing, Swift-boundary, or same-byte integration conditions.

## Acceptance conditions

This ADR may become Accepted only after review confirms the proposed exact values and Task 2 retains
positive, exact-boundary, cap-plus-one, duplicate, malformed, wrong-domain, replay, expiry, clock-
rollback, capacity, and no-state-change fixtures. Acceptance does not activate a public endpoint or
authorize guest execution.
