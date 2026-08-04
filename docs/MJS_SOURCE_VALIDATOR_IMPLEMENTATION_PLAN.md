# `.mjs` Source Validator implementation, conformance, and fault plan

## Purpose and non-goals

This plan converts Accepted ADR-0035 and ADR-0036's parser/process choice into independently
reviewable gates. It does not implement the validator or replace M1's passive
source/manifest contract, a product endpoint, runtime, loader, backend, or guest.
M1 owns authoritative proposal bytes, manifest fixtures, and scanner-stop evidence;
this work consumes those exact bytes additively after they land.

The end state is a parser-based pre-approval policy/usability check plus separate,
unconditional runtime no-loader enforcement. Neither layer substitutes for the
other.

## Historical V0 typed operation

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

These V0 bytes remain passive historical evidence. The supported macOS implementation consumes
the new role-specific v1 families in the
[passive v1 boundary](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md); it never relabels or
accepts V0/V1/V2 bytes as v1.

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
generator consumes M1's canonical bytes in place and adds 5 rules, 128 cases, and 128 fixtures.
Before R1, the full authority manifest covered 228 fields across 20 passive targets. This closes
only V0's byte/ownership question. It does not enroll Oxc or an executable, create a process or
endpoint, or begin the later product slices.

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

### V2 — historical direct-child disposable process profile

- Define the platform-specific fixed launch descriptor and sandbox in its own
  reviewable implementation slice.
- Permit only the copied request and result descriptors plus required loader
  mechanics for the enrolled executable; close all other inherited descriptors.
- Deny network, durable writes, caller-selected paths, environment authority,
  package caches, key access, Supervisor state, runtime/backend control sockets,
  and guest lifecycle operations.
- The historical acceptance rule required exact memory, CPU, output, process-count, and wall
  ceilings plus kill/reap. ADR-0036 replaces only that path's unavailable memory rule for the new
  role-specific v1 profile; it does not reinterpret V2.

Exit evidence: positive descriptor inventory and adversarial attempts for each
denial, fork/child escape, inherited-FD audit, memory/CPU/output exhaustion,
deadline/reap behavior, and post-crash clean restart. Until this evidence exists,
process separation is crash isolation only, not an admitted security boundary.

Observed V2 checkpoint: [`artifacts/mjs-source-validator-v2/`](../artifacts/mjs-source-validator-v2/)
retains the exact macOS arm64 stop plus a Darwin-only diagnostic harness. The strict test bootstrap
fixes the V1 target and profile argument, empty environment, fresh cwd, three post-exec role FDs,
one process group, and CPU/file/FD/child/output/wall limits, but lowering `RLIMIT_AS` to 256 MiB
returns `EINVAL` and refuses before V1 `exec`. An explicitly named unbounded-memory diagnostic
mutation then verifies ordinary/exact-maximum V0 results, partial/duplicate/trailing/oversize
output refusal, crash/signal/timeout/cancellation, group kill/reap, and a clean later invocation.
That mutation also permits a 512 MiB mapping, an owned out-of-cwd file read, IPv4 and Unix socket
creation without connect, cwd metadata change, and empty-file creation. Apple's supported embedded
App Sandbox child entitlement shape changes the fixed V1 Mach-O bytes; deprecated custom sandbox
interfaces are not used. No Keychain or Supervisor state was probed. Therefore V2 is `BLOCKED`,
not passed or `NO_GO`. The diagnostic mutation is test-only counterevidence and may not be wired;
the supported replacement starts at R1 with new role-specific identities.

Supported replacement review: the exact design slice in
[`MJS_SOURCE_VALIDATOR_MACOS_PROFILE_REPLACEMENT.md`](MJS_SOURCE_VALIDATOR_MACOS_PROFILE_REPLACEMENT.md)
is `PASSED`, while the historical V2 path and product parent remain `BLOCKED`. Apple's supported direct-helper App Sandbox
inheritance gives the parser its daemon or Broker parent's static rights and is `NO_GO` for this
exact boundary. ADR-0036 completes R0 by selecting two separately sandboxed role-specific private
XPC launchers, each owning one fresh matching parser child and no shared service/result/cache. It
accepts the role-private writable container only as residual scratch authority with mandatory
cleanup, and replaces the unavailable hard ceiling with an evidence-derived reactive footprint
watermark. The watermark is not a hard peak/exact cap or host-availability guarantee.

The replacement proceeds only through R1-R5B below. Any build uses new role-specific protocol,
process-profile, artifact-profile, signed-bundle, entitlement, Hardened Runtime, launch/library-
constraint, resource-policy, and supported-host identities. V0/V1/V2 bytes are never relabeled. No
signing operation is authorized by this plan, and no threshold, sample interval, baseline,
overshoot, or kill latency is chosen before the signed R4 corpus.

## Accepted replacement sequence

### R1 — passive v1 contracts and fixtures

Status: `PASSED` for the unwired passive-contract scope. The generated corpus retains 46 v1 cases
with independent Go/Node decoding, separate daemon/Broker known answers, exact/cap-plus-one
boundaries, strict cross-role/cross-version refusal, and an inactive resource-policy record that
rejects invented measurements. No launcher, parser child, signing, service, consumer, or active
policy exists.

Freeze the two role-specific request/result/process/artifact-profile families, the reactive-policy
record shape, copied ownership, exact caps, cleanup/refusal dispositions, and cross-role/cross-
version/mixed-update refusals. Update the canonical field-authority manifest in the same change.
Preserve every V0/V1/V2 byte and fixture unchanged. R1 may define field widths and missing/inactive
policy refusals but must not invent active reactive-resource measurements.

### R2 — unsigned role-specific construction

Construct the two smallest private XPC launcher bundles and their matching parser children from the
reviewed lock, offline and without a product consumer. Retain exact source/SBOM/notice/provenance,
two-builder, static linked-closure, descriptor, and no-generic-bus evidence. Do not sign, install,
use an Apple credential, or run arbitrary/user source.

### R3 — separately authorized signing and installation

Only a separately authorized task may use exact matching Apple identities/profiles to sign and
install the new role-specific bytes. Prove each containing role privately reaches only its own
service, the two sandboxes/containers and code requirements are distinct, launch/library
constraints and entitlements match, and old/mixed packages refuse. Unsupported private-XPC
reachability stops; it never widens to a shared/global service or app group.

### R4 — confinement, reactive-resource, and residue corpus

Use only fixed benign parse-only evidence to run the complete authority-denial, descriptor,
network/filesystem/native-loading, child/death/orphan, output/deadline, restart/update/startup, and
container-residue matrix. Measure per-role and simultaneous two-role baseline, threshold,
sampling cadence, maximum observed overshoot, kill latency, and host-pressure behavior. A profile
review may select values only from this signed corpus and must state host-availability limitations.
Every request, crash, launcher restart, update, and startup requires cleanup/empty-inventory
evidence; cleanup is not confidentiality or secure-erasure proof.

### R5D — daemon consumer

After R1-R4 pass, implement only the daemon-facing v1 client over copied decoded bytes before plan
construction. It uses one fresh daemon-role child, accepts no Broker service/result, has no fallback
parser, and maps every failure to a fixed refusal.

### R5B — Approval Broker consumer

After R5D passes, independently implement only the Broker-facing v1 client over exact bytes fetched
from Supervisor registration state. It uses one fresh Broker-role child, accepts no daemon
service/result/cache, renders only fixed facts, and performs zero Approval-key operations on every
refusal.

### M2/S1 checkpoint

Only after R5B passes may the project reconcile JobProposal narrowing and the plan-v0
registration/fetch field-authority/fixture slice. This checkpoint does not waive V6 runtime
no-loader evidence or any independent Supervisor IPC/product-admission gate.

### V3 / R5D — daemon pre-plan integration

- Decode M1's exact copied source bytes, validate cap/UTF-8/BOM according to its
  canonical contract, then invoke a fresh validator before plan construction.
- Bind the returned digest/length to exactly the supplied bytes and permit plan
  construction only for `valid + allow + all five zero counts`.
- Map every child/protocol/local failure to a fixed refusal without retrying a
  different parser or falling back to lexical scanning.

V3 consumes the accepted role-specific replacement protocol/profile version, not the passive V0 or
unchanged V1 test artifact. It may begin only as R5D after R1-R4 pass.

Exit evidence: byte mutation between decode/validation/plan, wrong-child artifact,
timeout/crash/malformed result, duplicate reply, process exhaustion, and daemon
restart cases. No path may accept caller-supplied validation facts.

### V4 / R5B — independent Broker integration

- Fetch only Supervisor-retained plan, registration, source manifest, and exact
  source bytes.
- Revalidate all ADR-0034 custody bindings, invoke a new validator instance, and
  bind its result to the fetched bytes before rendering or key use.
- Render fixed typed node facts and the explicit runtime-layer warning; never
  render child-supplied arbitrary text.
- Refuse if the daemon's historical result is present, substituted, or treated as
  trusted.

V4 invokes its own fresh parser child through the Broker-private launcher only after R5D passes. It
may not share a cached child, result, service state, app group, Keychain group, or daemon-owned
container.

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
| hang / CPU or reactive-footprint watermark | kill/drain/reap entire child group, fixed refusal | cleanup; new clean child for a later explicit request only |
| crash / signal / unexpected exit | fixed refusal | reap; never reuse partial output |
| partial, oversized, duplicate, trailing, or unknown result | fixed refusal | close channel and child |
| digest, length, version, or artifact mismatch | fixed refusal and bounded security event | quarantine artifact/installation when the enrolled policy requires it |
| parent cancellation or caller disconnect | kill/reap child; no authority mutation | later caller starts a fresh operation |
| daemon skips or forges validation | Broker independently refuses | no Approval-key operation |
| validator falsely allows | runtime loader still refuses requests | investigate; parser is not runtime boundary |
| validator falsely refuses | no plan/approval | bounded denial of service; no fallback parser |

## Acceptance rule

ADR-0035 and ADR-0036 are Accepted architecture decisions. Product Source Validator remains
`BLOCKED` until R2-R5B retain reviewed evidence on every supported host and the canonical docs,
field-authority material, and M2/S1 checkpoint agree. Runtime/profile admission additionally
requires V6. A parser binary, endpoint, signed package, cleanup pass, or happy-path result alone is
insufficient and does not make the product control passed or admitted.

For macOS, the active profile must bind evidence-derived reactive watermark values and measured
overshoot/kill/pressure limitations from R4. Parent-side caps and sampling alone do not establish a
hard peak or host-availability guarantee. Any ADR-0036 stop condition halts the exact candidate
without abandoning the Source Validator capability.
