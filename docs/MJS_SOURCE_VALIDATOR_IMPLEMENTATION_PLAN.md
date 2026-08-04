# `.mjs` Source Validator implementation, conformance, and fault plan

## Purpose and non-goals

This plan converts Proposed ADR-0035's parser/process choice into independently
reviewable gates. It does not implement the validator or replace M1's passive
source/manifest contract, a product endpoint, runtime, loader, backend, or guest.
M1 owns authoritative proposal bytes, manifest fixtures, and scanner-stop evidence;
this work consumes those exact bytes additively after they land.

The end state is a parser-based pre-approval policy/usability check plus separate,
unconditional runtime no-loader enforcement. Neither layer substitutes for the
other.

## Proposed typed operation

`ValidateMjsSourceV0` accepts only:

- fixed operation and protocol version;
- one copied byte string named by the protocol as `main.mjs`;
- declared byte length, which must equal the copied length and be at most 262,144;
  and
- a parent-computed SHA-256 used only for copy-binding confirmation.

The child recomputes SHA-256 and returns a fixed-size result containing:

- operation and implementation/artifact version;
- recomputed digest and byte length;
- one closed parse status and one closed policy status; and
- bounded counts for static import, export-from, import expression, and
  `import.meta` nodes; and
- a bounded count of unresolved references to the unavailable CommonJS bindings
  `require`, `module`, `exports`, `__dirname`, and `__filename`.

No unbounded strings or diagnostics cross the production channel. Internally
useful diagnostics stay in bounded test-only evidence and must not echo source or
specifier text into authority logs. Unknown fields, extra frames, partial frames,
wrong versions, counter overflow, trailing bytes, or multiple requests refuse.

## Dependency-ordered gates

### V0 — contract reconciliation

- Preserve the merged authoritative M1 source/manifest fixtures and consume them
  directly without duplicating fixtures.
- Import their exact strict-UTF-8/BOM/cap semantics.
- Define canonical request/result bytes, maximum frame sizes, status enums,
  counter bounds, digest binding, and one-request EOF behavior.
- Classify every new field in the canonical field-authority material before an
  endpoint exists.

Exit evidence: byte-level known answers, exact/cap-plus-one frames, decoder
predecode tables, and no disagreement with ADR-0034 or M1.

Observed passive slice: the exact
[`capsule.source-validator.protocol/v0`](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_CONTRACT.md)
fixed-frame projection is retained with independent unwired Go and test-only Rust codecs. The
generator consumes M1's canonical bytes in place and adds 5 rules, 128 cases, and 128 fixtures;
228 fields across 20 passive targets are now classified. This closes only V0's byte/ownership
question. It does not enroll Oxc or an executable, create a process or endpoint, or begin V1-V6.

### V1 — hermetic parser artifact

- Package only the exact Oxc 0.140.0 parser/AST/visitor/semantic mode selected by
  the experiment; no formatter, transformer, resolver, runtime, or loader.
- Complete source/license/notice review for the full locked graph and record a
  reproducible source and binary manifest.
- Enroll the immutable executable digest, toolchain, target, build flags, and
  artifact review under the installation trust domain.
- Define an update rule that never lets an unreviewed parser replace an active
  artifact and makes rollback/trust-epoch consequences explicit.

Exit evidence: reproducible-build comparison, SBOM/notice review, artifact
signature and assessment, tamper/rollback/missing-artifact refusals.

Observed bounded V1 checkpoint: [`artifacts/mjs-source-validator-v1/`](../artifacts/mjs-source-validator-v1/)
packages only the selected Oxc mode plus pinned `sha2` for the V0 digest, emits exact V0 results for
all 28 M1 HOLD cases, retains a 74-dependency lock, complete checksum/source/license/notice and
CycloneDX inventories, and reproduces the 1,146,656-byte macOS arm64 executable byte-for-byte in
two clean source/target directories on one host. The executable SHA-256 is
`ba2a6b38be6b4eea8c067887cf80988756e2f4a551d128bf2dabdaf7f2ecb600`; its not-enrolled V0
artifact-profile identity is `cfadcedc3e983377b964e0465c1f7127a307acbfda15ad8a02d7a302e82b4ce7`.

This checkpoint does not satisfy V1 enrollment: the Mach-O has only an identity-free linker ad-hoc
signature, provenance and assessment are unsigned, both
builds share one host/toolchain/cache/administrator, no installation authority has enrolled the
profile, and no vulnerability-monitoring owner or response SLA exists. Missing, tampered, unknown,
unsigned, or rolled-back profiles therefore remain refusals rather than accepted artifacts. V2
still owns the fixed launch descriptor, sandbox, descriptors, resources, deadline kill/reap, and
restart evidence; the retained partial-input kill observation is not V2 confinement proof.

### V2 — disposable process profile

- Define the platform-specific fixed launch descriptor and sandbox in its own
  reviewable implementation slice.
- Permit only the copied request and result descriptors plus required loader
  mechanics for the enrolled executable; close all other inherited descriptors.
- Deny network, durable writes, caller-selected paths, environment authority,
  package caches, key access, Supervisor state, runtime/backend control sockets,
  and guest lifecycle operations.
- Enforce exact memory, CPU, output, process-count, and wall ceilings; kill and
  reap the entire child on deadline or parent cancellation.

Exit evidence: positive descriptor inventory and adversarial attempts for each
denial, fork/child escape, inherited-FD audit, memory/CPU/output exhaustion,
deadline/reap behavior, and post-crash clean restart. Until this evidence exists,
process separation is crash isolation only, not an admitted security boundary.

### V3 — daemon pre-plan integration

- Decode M1's exact copied source bytes, validate cap/UTF-8/BOM according to its
  canonical contract, then invoke a fresh validator before plan construction.
- Bind the returned digest/length to exactly the supplied bytes and permit plan
  construction only for `valid + allow + all five zero counts`.
- Map every child/protocol/local failure to a fixed refusal without retrying a
  different parser or falling back to lexical scanning.

Exit evidence: byte mutation between decode/validation/plan, wrong-child artifact,
timeout/crash/malformed result, duplicate reply, process exhaustion, and daemon
restart cases. No path may accept caller-supplied validation facts.

### V4 — independent Broker integration

- Fetch only Supervisor-retained plan, registration, source manifest, and exact
  source bytes.
- Revalidate all ADR-0034 custody bindings, invoke a new validator instance, and
  bind its result to the fetched bytes before rendering or key use.
- Render fixed typed node facts and the explicit runtime-layer warning; never
  render child-supplied arbitrary text.
- Refuse if the daemon's historical result is present, substituted, or treated as
  trusted.

Exit evidence: compromised-daemon simulations, source/result replay, digest and
length mismatch, child failure, view mutation, and proof that zero Approval-key
operations occur on every refusal.

### V5 — grammar conformance and mutation suite

Run the retained corpus plus the authoritative M1 corpus and a reviewed,
version-pinned bounded Test262 selection. Required categories include:

- static, side-effect, attributed, and multiline imports;
- named/all/namespace export-from declarations;
- literal/computed import expressions in nesting, regex/division ambiguity,
  template interpolation, and generated-code strings;
- `import.meta` versus ordinary property names/comments;
- object/class method names, strings, comments, regexps, escapes, and Unicode;
- free CommonJS bindings versus locally shadowed same-name bindings;
- every parser diagnostic and early error, including recovery-before-forbidden
  syntax; and
- one mutation/restoration test for each forbidden AST category and every status
  branch.

All cases must be deterministic across supported architectures and repeated fresh
processes. Differential parser results are investigation triggers, not a voting
oracle. Corpus passage does not claim semantic equivalence or arbitrary hostile
code safety.

### V6 — runtime no-loader evidence

Independently prove, for the exact admitted runtime bundle and profile, that:

- no filesystem, network, package, URL, import-map, Capsule-internal, or fallback
  module resolver exists;
- every static and dynamically constructed request refuses, including through
  eval/generated code;
- `import.meta` cannot expose a URL, resolver, host path, or ambient authority;
- source reaches the runtime only from Supervisor-retained attempt-bound bytes;
  and
- bypassing or lying in both validators does not create loader authority.

This gate belongs to runtime admission. Parser passage alone never satisfies it.

## Fault and recovery matrix

| Fault | Required result | Recovery rule |
| --- | --- | --- |
| invalid UTF-8 / cap+1 | refuse before parse | no child retry |
| parser or semantic diagnostic | typed refuse | no recovered AST use |
| hang / CPU or memory exhaustion | kill entire child, fixed refusal | new clean child for a later explicit request |
| crash / signal / unexpected exit | fixed refusal | reap; never reuse partial output |
| partial, oversized, duplicate, trailing, or unknown result | fixed refusal | close channel and child |
| digest, length, version, or artifact mismatch | fixed refusal and bounded security event | quarantine artifact/installation when the enrolled policy requires it |
| parent cancellation or caller disconnect | kill/reap child; no authority mutation | later caller starts a fresh operation |
| daemon skips or forges validation | Broker independently refuses | no Approval-key operation |
| validator falsely allows | runtime loader still refuses requests | investigate; parser is not runtime boundary |
| validator falsely refuses | no plan/approval | bounded denial of service; no fallback parser |

## Acceptance rule

ADR-0035 may move from Proposed only when V0–V5 have retained, reviewed evidence
on every supported host and the canonical docs reflect the final protocol and
authority ownership. Runtime/profile admission additionally requires V6. A parser
binary, IPC endpoint, or passing happy-path test alone is a NO-GO.
