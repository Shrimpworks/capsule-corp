# `.mjs` Source Validator passive contract v0

Status: observed passive conformance slice for Proposed ADR-0035. This document fixes bytes and
test oracles only. It does not define a product endpoint, launch a process, invoke Oxc, enroll the
fixture artifact, authorize planning or Approval, or implement runtime no-loader enforcement.

## Closed nominal identities

| Role | Identity |
| --- | --- |
| Protocol | `capsule.source-validator.protocol/v0` |
| Method | `capsule.source-validator.validate-mjs-source/v0` |
| Request | `capsule.source-validator-request/v0` |
| Result | `capsule.source-validator-result/v0` |
| Validator profile | `capsule.source-validator.one-shot/v0` |
| Source profile | `capsule.mjs-source/v0` |
| Source media type | `application/capsule.javascript-source;v=0;module=esm` |
| Correlation value | `capsule.source-validator.validation-request-id/v0` |
| Engineering candidate | `capsule.source-validator.oxc-engineering-candidate/v0` |
| Artifact profile | `capsule.source-validator.artifact-profile/v0` |

The frame tags below are closed projections of those nominal identities. A tag is never accepted
without the expected frame magic, frame kind, and complete profile tuple. No string, optional
field, map, extension area, or generic RPC envelope exists on this wire.

## Framing and scalar rules

Every integer is an unsigned big-endian integer. The first `u32` is the exact byte count after the
length word, so total frame length is `bodyLength + 4`. Decoders widen it before addition and
compare it with the object-specific inclusive maximum before converting to a host index. All
encoded values fit `u32`; JavaScript generators use only safe integers and reject values outside
`0..2^32-1` before encoding. Counts and source lengths have the tighter caps below.

One supplied byte slice must contain exactly one frame. A short length word, truncation, a declared
length different from the supplied slice, a second concatenated frame, or any trailing byte
refuses. Fixed offsets make unknown and duplicate fields unrepresentable; a changed reserved byte,
magic, tag, or unused value refuses. Decoders never clamp.

### Common request/result prefix

| Offset | Width | Value |
| ---: | ---: | --- |
| 0 | 4 | body length |
| 4 | 8 | object magic |
| 12 | 2 | protocol version `0` |
| 14 | 2 | frame kind: request `1`, result `2` |
| 16 | 2 | method tag `1` |
| 18 | 2 | validator-profile tag `1` |
| 20 | 2 | source-profile tag `1` |
| 22 | 2 | source-media tag `1` |
| 24 | 2 | correlation domain `0x0001` |
| 26 | 2 | source-digest domain `0x0101` |
| 28 | 16 | nonzero validation-request correlation ID |
| 44 | 4 | exact source byte length |
| 48 | 32 | SHA-256 of the exact source bytes |

The correlation ID exists only to bind a reply to the current in-memory request. It is not an
artifact, registration, approval, attempt, installation, runtime, backend, or guest identity and
grants no authority.

## Validation request

The request magic is `CAPMJSRQ`. Bytes `80..EOF` are the defensive copy of the exact canonical M1
`main.mjs` bytes. The declared length must equal that suffix and the SHA-256 must recompute exactly.
The source must independently satisfy M1: `0..262,144` bytes inclusive, strict UTF-8, no leading
UTF-8 BOM, and no rewriting or normalization. Therefore:

| Bound | Bytes |
| --- | ---: |
| Fixed header / minimum frame | 80 |
| Maximum body | 262,220 |
| Maximum frame | 262,224 |
| Cap-plus-one refusal frame | 262,225 |

The request carries no path, filename supplied by a caller, package, specifier, network, loader,
environment, runtime, backend, key, Approval, Supervisor, or guest data.

## Validation result

The result magic is `CAPMJSRS` and its only size is 138 bytes.

| Offset | Width | Value |
| ---: | ---: | --- |
| 0..79 | 80 | common prefix, recomputed child source facts |
| 80 | 2 | artifact-profile-digest domain `0x0106` |
| 82 | 32 | domain-separated enrolled-artifact-profile digest |
| 114 | 1 | parse disposition |
| 115 | 1 | child-reported policy disposition, accepted only after parent derivation |
| 116 | 1 | closed classification |
| 117 | 1 | reserved, exactly zero |
| 118 | 4 | static-import node count |
| 122 | 4 | export-from node count |
| 126 | 4 | import-expression node count |
| 130 | 4 | `import.meta` node count |
| 134 | 4 | unresolved free-CommonJS-reference count |

Parse dispositions are `1=valid`, `2=parser-diagnostic-refusal`, and
`3=semantic-diagnostic-refusal`. Reported policy dispositions are `1=allow`, `2=deny`, and
`3=not-evaluated`. Classifications are `0=no-finding`, `1=forbidden-syntax`,
`2=free-commonjs`, `3=forbidden-syntax-and-commonjs`, `4=parser-diagnostic`, and
`5=semantic-diagnostic`.

Each count is inclusively bounded by both `262,144` and the recomputed source byte length. A valid
parse is `allow/no-finding` exactly when all five counts are zero; otherwise it is `deny` with the
classification derived from the four syntax counts and the separate CommonJS count. Either
diagnostic status is `not-evaluated`, has its matching diagnostic classification, and has five zero
counts. No recovery AST is usable. No parser prose, source excerpt, specifier, path, unbounded
collection, key, or authority decision is returned.

The parent independently decodes the artifact profile, recomputes its identity digest, and matches
the result's correlation ID, source length, and source digest to the current request. Replay onto a
different correlation ID, cross-request substitution, source mutation, response mutation, or
artifact substitution refuses before any state or authority change.

The raw child report is a parser observation, not a parent policy decision. The parent derives the
only acceptable policy/classification pair again from the closed parse disposition and five
bounded counts, then requires byte-for-byte agreement. Transport framing and the correlation ID,
artifact enrollment identity, this parser observation, the parent's derived decision, and runtime
no-loader enforcement are five distinct layers; none can stand in for another.

`allow` is only a passive parser/policy observation. It is not plan construction authority, an
Approval decision, launch authority, or evidence of runtime no-loader enforcement.

## Engineering-candidate record

The fixed 292-byte `CAPMJSCI` record is the exact ADR-0035 engineering candidate. Its body-length
word is `288`. The layout is:

| Offset | Width | Value |
| ---: | ---: | --- |
| 12 | 2 | record version `0` |
| 14 | 2 | Oxc implementation tag `1` |
| 16 | 6 | Oxc semantic version `0.140.0` as three `u16` values |
| 22 | 6 | Rust candidate version `1.95.0` as three `u16` values |
| 28 | 2 | semantic-mode mask `0x000f`: parser, module AST, syntax-error checking, unresolved references |
| 30 | 2 | direct Oxc crate count `6` |
| 32 | 2 | locked-package count `65` |
| 34 | 2 | reserved zero |
| 36 | 4 | cached locked-source bytes `24,449,903` |
| 40 | 4 | macOS arm64 release-probe bytes `1,854,528` |
| 44 | 32 | retained parser-experiment `Cargo.lock` SHA-256 |
| 76 | 216 | six 36-byte direct-crate records |

Each direct-crate record is `u16 crate ID`, `u16 reserved zero`, and a 32-byte crates.io checksum.
IDs `1..6` are, in order, `oxc_allocator`, `oxc_ast`, `oxc_ast_visit`, `oxc_parser`,
`oxc_semantic`, and `oxc_span`, all version `0.140.0`. Their exact checksums and the complete lock
identity are encoded by the deterministic generator and checked independently in Go and Rust.
Any byte other than the exact retained candidate refuses; this record does not admit the graph.

## Artifact-profile record

The fixed 160-byte `CAPMJSAP` record defines the future enrollment projection. Its body length is
`156`; offsets `12..23` are record version `0`, protocol version `0`, and method, validator-profile,
source-profile, and source-media tags, all `1`. It then contains exactly:

| Offset | Domain | Bytes |
| ---: | ---: | --- |
| 24 | `0x0102` engineering candidate | digest at 26..57 |
| 58 | `0x0103` executable | digest at 60..91 |
| 92 | `0x0104` build manifest | digest at 94..125 |
| 126 | `0x0105` assessment | digest at 128..159 |

Every digest is nonzero and nominally typed. The retained conformance profile uses `11`, `22`, and
`33` repeated for the executable, build, and assessment digest bytes solely to freeze the layout;
it is not an enrolled or executable product artifact.

Record identities are domain-separated as
`SHA-256(ASCII(domain) || 0x00 || exactRecordBytes)`. The candidate domain is
`capsule.source-validator.engineering-candidate/v0`; the artifact-profile domain is
`capsule.source-validator.artifact-profile/v0`. Raw SHA-256 and either domain-separated digest are
not interchangeable.

## Retained known answers

These values are generated, never hand-edited:

| Fixture | Length | SHA-256 |
| --- | ---: | --- |
| engineering candidate | 292 | `2c39757c40198074f1b1dd6e0ed37fb6c75c1c699c0090e2aa4b8ae88cecc9af` |
| artifact profile | 160 | `2075ee498ce4b3d81843d57c5289f0056092aa5e3b575d885018c7348419fc8b` |
| minimum request | 80 | `4e40f3057a7d4fe7814b74806d794c91e30b5f79f8b55d12c0d76c0c177fc4f7` |
| ordinary request | 137 | `5ea0960c1b8200f483ee29fa2756b15c69a62aadcb1dc675c244569b295795fc` |
| maximum request | 262,224 | `24db6d0599ba05dead0e7eac1f1454262506c8c187b5542694bb6bb12e6a0571` |
| cap-plus-one refusal | 262,225 | `4add49e80567b4c097bbad0f989863fe028b868304b18d4ef40c9349a888428b` |
| ordinary result | 138 | `5a428178447f70607367f0c052f674ce74315a60376cf15826de06562ea1ca25` |

The candidate identity digest is
`11a08cace3ddc4dde925bd93ea3de8ea313f7f5766839deca802078df038d0c6`; the artifact-profile
identity digest is `0fe1523741ac3cedc32e062778836d73f4eb21f422b3587d78ad81df8c180908`.

## Conformance, ownership, and limitations

The deterministic corpus adds 5 rules, 128 cases, and 128 fixture files while referencing M1's
canonical bytes in place. It covers the minimum/ordinary/exact-cap/cap-plus-one, UTF-8/BOM,
newline/Unicode identity, all 28 language HOLD, module-request, scanner-stop, and SourceManifest
oracles. Rejections carry explicit zero state, Approval, key, IPC-endpoint, process, runtime,
backend, and guest effects. Mutations restore enforcement of every version, domain, media,
artifact, digest, length, status, count, reserved/trailing, and cap check.

[`sourcevalidatorpassive`](../../internal/protocol/sourcevalidatorpassive/contract.go) is the
unwired Go parent-language model. The standalone
[`mjs-source-validator-passive-contract`](../../experiments/mjs-source-validator-passive-contract/README.md)
crate is an independent Rust child-language **test-only** oracle; product packages do not import it.
It uses exact `sha2` 0.10.9 and test-only `serde_json` 1.0.151 under a 22-package lock with SHA-256
`a45ad0e2b2311d33b16e46e0bf1f66c1563dd240a35f1f9fe431c7bea5894c98`. This applies the reuse
map's standard-cryptography and Rust-test recommendations: reuse a pinned primitive, keep the
Capsule-specific fixed framing narrow, run offline, and remove the oracle once an independently
reviewed enrolled child owns the same codec. The crate gains no product authority.

Artifact construction/admission (V1), process confinement (V2), daemon and Broker consumers
(V3/V4), broader grammar evidence (V5), and runtime no-loader enforcement (V6) remain
unimplemented. Proposed ADR-0035 therefore remains Proposed.
