# V0 CBOR object-set freeze and unwired Go wrapper result

Work item: applicable v0 CBOR object-set freeze and narrow Go wrapper

Status: `PASSED`

Scope: defensively freeze the exact first-release object subset that may use the pinned Go CBOR
dependency, retain the pre-freeze signed and conditional objects outside that subset, implement
one unwired object-specific `SourceManifest` codec, and admit the passive I2B1 Supervisor bootstrap
request/record payloads to the same narrow typed-decode dependency boundary. The work uses only
repository fixtures, public fixture material, cached exact module sources, and local tests. It
creates no production signing, IPC, store, runtime, backend, or guest authority.

Evidence or reason: accepted ADR-0034 makes the five-field, single-member `SourceManifest` v0
shape exact. I2B1 now separately freezes the 34-field request and 69-field record with closed CDDL,
field authority, actual-shape maxima, 95 independently generated cases, and agreeing Go/Swift
verification. Every other current CBOR or signed candidate still has an explicit schema,
authority, consumer, or ADR gate. The `internal/protocol/v0cbor` package imports exact
`fxamacker/cbor/v2` v2.9.2 only for deterministic encoding and typed field decoding after the
retained Capsule predecoder. Its tests replay all 19 applicable Go/TypeScript manifest cases,
three known answers, 10,000 deterministic round trips, the handwritten oracle, cap-plus-one and
restoration cases, defensive ownership, and a 30-second fuzz run with 673,967 executions and no
failed invariant. The focused race run also passes.

Remaining work: the installed Coordinator/Supervisor same-byte path, production key authorization
and signing, the rest of the signed-object family, and any
`ExecutionPlan`/`PlanRegistration` dependency cutover remain separate coordinated work. ADR-0019
remains Proposed.

Next action: keep the wrapper unwired. Revisit its object set only in the same reviewed change that
freezes another object's CDDL, authority manifest, caps, known answers, consumers, and failure
classifications.

Parent status: the larger production canonicalization/signing workstream remains `BLOCKED` on the
ADR-0019 signed-object, Swift, key-authorization, and same-byte integration conditions.

## Exact object-set decision

“Current” means a repository-retained v0 CBOR shape or signed-payload candidate, not every object
that the target object inventory says may eventually be signed. “Eligible” means the shape and
authority semantics are exact enough for this unwired Go wrapper; it does not mean a product
consumer or authority path is admitted.

| Object/profile | Current owner and intended consumers | Exact retained bounds and known answers | Decision |
| --- | --- | --- | --- |
| `SourceManifest` v0 | Daemon derives it from exact `main.mjs`; Supervisor and Broker will independently validate retained bytes; current consumers are passive Go/TypeScript tests only | CBOR 87..95 bytes; depth 4; 15 items; one map up to 5 pairs; arrays up to 3 elements; separately bound source 0..262,144 bytes. Known answers: minimum 87 bytes / `b29a880488898dbac54c5890ec56b97a8aca461cbf6585af60fa3a56bc9e9044`, ordinary 89 / `052dce0c353e1efeb70f93405a8757ef6fa4d29f91a4d6bcaa67d00c45abc0d6`, maximum 95 / `21b9b521ed0ab6973be563b14f1eae4591097423e17d4d00ba0025c3fe4298d8` | **FROZEN FOR THIS WRAPPER / ELIGIBLE.** ADR-0034 fixes all five fields, one member, media, source binding, and limits. |
| `SupervisorBootstrapRequest` v0 | Future Trust Coordinator producer and Supervisor consumer; current consumers are passive Go and independent Swift verification only | Payload raw/calculated 2,048/860 bytes; protected 256/98; envelope 4,096/1,032; payload depth 2, 69 items, 34 map pairs, no arrays. | **PASSIVE I2B1 FROZEN / ELIGIBLE FOR TYPED DECODE ONLY.** The closed wrapper retains predecode, canonical-byte, signature, authorization, time, payload replay, and copy protections. No production signer or installed consumer exists. |
| `SupervisorBootstrapRecord` v0 | Future Trust Coordinator producer and Supervisor consumer; current consumers are passive Go and independent Swift verification only | Payload raw/calculated 4,096/1,526 bytes; protected 256/97; envelope 6,144/1,697; payload depth 2, 139 items, 69 map pairs, no arrays. | **PASSIVE I2B1 FROZEN / ELIGIBLE FOR TYPED DECODE ONLY.** Exact signed-request pairing, all repeated request fields, observed root/owner/store, record time, transition, retention, one-use, and response-loss bindings are closed. Installed enrollment remains blocked. |
| `ExecutionPlan` v0 | Future daemon producer; future Supervisor/Broker consumers; current unwired registration tests decode it | CBOR up to 65,536 bytes; depth 8; 256 items; map 64; array 8. Ordinary fixture 530 bytes / `627f9524479000dab6f3cee1d70c0428c63285bcadbc2cb3c6e8018b2dea008c` | **PRE-FREEZE / EXCLUDED.** ADR-0023 remains Proposed; the candidate omits unresolved authority/transport values and cannot authorize execution. |
| `PlanRegistration` v0 | Future Supervisor producer; future daemon/Broker consumers; current unwired registration state emits and decodes it | CBOR up to 4,096 bytes; depth 4; 33 items; map 16; arrays forbidden. Ordinary fixture 165 bytes / `f3569d37ad6d787c2cdd575ef9ec6c369bbe495157c43110fc9e9d610a277614` | **PRE-FREEZE / EXCLUDED.** The typed IPC method, final consumer contract, and ADR-0023 acceptance remain open. It is not independently signed portable authority. |
| `ApprovalGrant` v0 payload, protected headers, and tagged COSE_Sign1 envelope | Future Approval Broker signer and Supervisor verifier; current `approvalattempt.FixtureVerifier` is retained-vector-only | Provisional raw copy budgets 512/256/128 bytes; calculated candidate maxima 431/242/116. Ordinary envelope/payload/protected are 375/234/68 bytes with SHA-256 `fb0a9e7c983f6f3986260dce857edf6b18cba99ee386f9532300dbdc31a5a3bd`, `8ed203acb49409cf2c787bcb04e5e40aaed7139e8bc5b599bd53a49fb3c0e6ea`, and `b79d430399eb9d3f3690735f03a021a80a24f1ea76821303cf90fd010033ecbf` | **SIGNED PRE-FREEZE / EXCLUDED.** ADR-0019 and ADR-0024 remain Proposed; final session/key authorization and production COSE/Swift paths are open. The active fixture work is untouched. |
| Spike-only `EnforcementTranscript` signed profile | Experiment-only cross-object oracle; no repository product producer or consumer | The Gate A2 hardening experiment retains a mutually exclusive test shape and cross-language cases, not a final transcript CDDL | **NOT IN THE CURRENT V0 SET.** It is a test oracle only; the final evidence object remains unfrozen. |
| TypeScript approved-byte v0 family: original/executable manifests, transformer profile, normalized options, transformation record, record set, and plan-source bindings | Conditional future Source Preparer, Supervisor, Broker, and plan-v1 consumers; current Go code verifies fixture identity only | Fixture-only object cap 4,096 bytes. Known-answer lengths/digests are 171/`1010ae00c786a6266348173c7760e0190be4cc280be1f71c8549f09727e4b183`, 174/`295138062d0785785373b8c468fee75f77a28131d0974f30f69c4050425e9814`, 215/`ab82cd553d52490fb0e2cc2e6cfc8cc106440a601c99c70cebad43537efc7f97`, 188/`8a43f134b568e983a7e4e24a763209a61405dc535a710e4799611033e4983b2e`, 659/`3ddf8b242afbddaa1856f79a24a46f7e8f9c674699ef6b1bd63215d49e512c39`, 714/`5738283a5accdbd8b736af81982dc46068172ec502f5c43e4113fe7de10c76eb`, and 150/`deb9342fc8c2c0ea18ff280aa3c409246d7c96fab9dc46dddb54862640e4cc28` | **FUTURE-CONDITIONAL PRE-FREEZE / EXCLUDED.** ADR-0034 removes the family from the first-release path; per-object production maxima and consumers do not exist. |
| Illustrative `ExecutionPlan` v1 TypeScript shape | No current owner or consumer | No active wrapper, fixture corpus, or admitted cap | **NOT IN THE CURRENT SET.** It is an illustrative atomic migration candidate only. |

The target inventory also anticipates signed runtime bundles, review/validation claims, key
authorizations, installation/epoch records, trust snapshots, Supervisor transcripts, TUF metadata,
and witness checkpoints. None has a current production-frozen v0 signed payload profile, so this
task does not silently treat the object-model list as a signed wire set.

## Responsibility and byte-ownership matrix

| Behavior | Classification | Failure ownership / byte timing |
| --- | --- | --- |
| Select the `SourceManifest` entry point and exact media type | `RETAIN-CAPSULE` | Reject before CBOR processing; no caller byte copy |
| Enforce 87..95 raw bytes | `RETAIN-CAPSULE` | `MALFORMED`; before copy or library decode |
| Scan definite/preferred CBOR, map order/duplicates, UTF-8, tags, floats/simple values, safe integers, depth, total items, map pairs, array elements, and trailing data | `RETAIN-CAPSULE` | `MALFORMED`; allocation-independent scan of caller bytes before copy |
| Copy exact accepted manifest bytes and separately capped source bytes | `RETAIN-CAPSULE` | Only after the corresponding raw cap; every returned byte slice is another defensive copy |
| Deterministically encode the private fixed wire struct | `DELEGATE-FXAMACKER` | RFC 8949 Core Deterministic mode; tags and marshaler hooks disabled; no caller map/options/value surface |
| Decode fields into the private fixed wire struct | `DELEGATE-FXAMACKER` | Runs only after Capsule predecode and private copy; duplicate, indefinite, tag, UTF-8, bignum, custom-unmarshaler, and unknown-field options are also narrowed as defense in depth |
| Re-encode the typed value and compare it byte-for-byte with the exact received bytes | `RETAIN-CAPSULE` | `MALFORMED`; received bytes, never normalized output, remain identity |
| Enforce five fields, object/version, `main.mjs`, one three-field member, digest width, and source caps | `RETAIN-CAPSULE` | `UNSUPPORTED`, `SCHEMA`, or `DOMAIN` according to the existing first-owning-boundary vocabulary |
| Recompute SHA-256 and bind member digest/length and aggregate length to separately supplied exact source | `RETAIN-CAPSULE` | `DOMAIN`; standard-library hash over Capsule-owned bytes |
| Preserve exact-byte manifest identity and defensive typed projection | `RETAIN-CAPSULE` | Happens only after every structural, canonical, and binding check succeeds |
| Handwritten Go encoder/decoder and TypeScript known answers | `TEST-ORACLE` | Remain independent; no production path is switched in this task |
| Generic `cbor.Marshal`, `cbor.Unmarshal`, `UnmarshalFirst`, `any`, maps, tags, streams, diagnostics, and custom marshalers | `TEST-ORACLE` when used to prove a restoration hazard; otherwise unavailable | No exported wrapper route |
| I2B1 COSE envelope framing, Sig_structure, ES256 verification, key authorization, exact-byte identity, and passive replay decision | `RETAIN-CAPSULE` in the object-specific passive bootstrap verifier; `go-cose` remains product `NO_GO` and absent | Implemented only for checked-in passive verification. Payload/envelope digests and typed replay state, never signature bytes, own identity; no durable ledger or production signer is present. |

The wrapper deliberately maps fxamacker's private typed-decode failures into Capsule's bounded
classification vocabulary and never returns library diagnostics containing decoded input.

## Exact API boundary

The package exposes one constructor and one object-specific codec:

```go
codec, err := v0cbor.NewSourceManifestCodec()
encoded, err := codec.Encode(exactMainMJSBytes)
decoded, err := codec.Decode(receivedManifest, exactMediaType, exactMainMJSBytes)
```

The returned object exposes only `ExactBytes`, `Digest`, `SourceBytes`, and the closed immutable
`SourceManifestView`. There is no generic marshal/unmarshal method, caller-provided struct, map,
tag, decode mode, option, object discriminator, or `any` value. Existing production-shaped
registration and approval packages do not import `internal/protocol/v0cbor`.

## Independent dependency and restoration review

Capability and roadmap slice: Phase 2 deterministic CBOR object-specific wrapper.

Reuse-map row and recommendation: `STD-1 + GO-1`, `ADOPT-PINNED` narrowly after the applicable
object freezes.

Candidate: `github.com/fxamacker/cbor/v2 v2.9.2`, tag commit
`45589abe5c63bea2db4d311e0d0fcc551cd772ae`, module sum
`h1:X4Ksno9+x3cz0TZv69ec1hxP/+tymuR8PXQJyDwfh78=`, module zip SHA-256
`778e7b1e56acefbf9f9b6fd435e55a1dcc6d59bbd5513b9a60eac08886617d06`.

Primary sources and retrieval: exact cached Go module source and origin metadata independently
checked on 2026-08-04 against the commit-pinned experiment. The mutable GitHub release supplies no
asset signature or attestation. Module checksum, zip hash, tag commit, retained source, and
cross-implementation fixtures are sufficient for this unwired source-reviewable pin; they are not
an upstream reproducible-build or distribution attestation claim.

Maintenance and security: upstream states that the latest release receives fixes and accepts
private reports at its documented security address. The Capsule protocol maintainers own monthly
release/advisory review, immediate review on an upstream security notice, and the decision to
upgrade, retain, or remove the pin. No automatic version range or unattended upgrade is allowed.
The 2026-08-04 verbose `govulncheck` review reports no symbol or imported-package vulnerability in
fxamacker or float16. It separately reports GO-2026-5024 in the pre-existing `x/sys v0.28.0`
Windows package, fixed in v0.44.0, but finds no imported package or called symbol; that baseline
module update remains outside this wrapper's dependency change.

License and notices: fxamacker is MIT; license SHA-256
`78cad457d5ea7318230f3d969d4cdf29cef45524a1fc8ca3a97646da1ad7a841`.
The sole transitive module `github.com/x448/float16 v0.8.4` is MIT; zip SHA-256
`73b24a41037ea999ab66851e3798a0973dbb1f214925915b01f0820f7b2f1500`, license SHA-256
`a555f1194fdac34da70fb416968f7e2217b02352c26c1eac2fa45fcb4290ae8d`. Its module source has no
separate security policy; fxamacker and Capsule's upgrade owner absorb that risk. Required notices
are retained in `THIRD_PARTY_NOTICES.md`. `go-cose` and its MPL-2.0 obligations are not added.

Toolchain and graph: Go 1.20+ upstream, Capsule Go 1.23+. The complete added graph is fxamacker
v2.9.2 to float16 v0.8.4. The comparison measured 1,204 KiB plus 76 KiB extracted source and a
153,104-byte stripped arm64 delta. The dependency gains typed parsing/encoding authority only in
the unwired wrapper; it gains no key, store, filesystem, network, process, IPC, backend, or guest
authority.

Restoration risks independently retained:

- generic/default decode accepts features outside Capsule's profile and has much wider default
  depth/collection limits;
- `UnmarshalFirst` permits a sequence tail by design;
- library modes can enable tags, indefinite lengths, simple/floating values, bignums, maps/`any`,
  streams, custom CBOR/binary/text marshalers, and broader field matching;
- duplicate-key enforcement may partially fill a destination, so the wrapper always discards it on
  any error;
- typed decode does not establish canonical-on-wire identity, object/version policy, nominal role
  binding, trusted source bytes, key authorization, or replay semantics; and
- a future generic helper or shared mode would make restoration easier, so neither is exported or
  centralized as a decode-anything service.

The retained mutation corpus proves that duplicate, out-of-order, nonpreferred, invalid-UTF-8,
depth, collection, tag, indefinite, float, trailing, unknown-field, and wrong-type restorations
refuse. A direct mode probe proves fxamacker accepts a representative nonpreferred representation
while Capsule's predecoder and exact-byte comparison reject it. The retained 30-second fuzz run
completed 673,967 executions with 91 total corpus entries and no crash, acceptance disagreement,
or exact-byte invariant failure.

Offline reproduction uses the exact cached graph with `GOPROXY=off`. Upgrade or removal requires
rerunning the full corpus, 10,000-case property comparison, bounded fuzz target, dependency hashes,
licenses/notices, vulnerability scan, and restoration probes. The handwritten oracle remains until
a later coordinated schema freeze and cutover review explicitly removes it.

## Limitations and explicit non-claims

- This freezes only the Go-wrapper-eligible `SourceManifest` plus the passive I2B1 bootstrap
  request/record subset, not the complete future signed-object family and not ADR-0019.
- The package is unwired and has no product consumer, authenticated IPC, store, Broker, Supervisor
  same-byte path, key, signature, replay ledger, runtime, backend, or guest.
- Swift has an independent executable I2B1 fixture verifier, not a production wrapper or signer.
- The dependency's release is identified by exact public source/module evidence but is not signed
  or attested upstream.
- The corpus and bounded fuzzing support this wrapper result only; they do not establish a
  continuous-service or production-readiness claim.
