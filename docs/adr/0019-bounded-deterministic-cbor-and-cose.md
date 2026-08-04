# ADR-0019: Use bounded deterministic CBOR and object-specific COSE profiles

- Status: Proposed
- Date: 2026-07-31
- Replaces if accepted: the RFC 8785/JWS candidate in the pre-implementation technical design

## Context

Capsule authorities must agree on the exact bytes of registered and signed security objects across
Go, Swift, and TypeScript. Gate A found that the tested RFC 8785 JSON Canonicalization Scheme path
did not produce a maintainable cross-language contract: Foundation number representation caused a
real disagreement and no maintained strict Swift JCS path passed the gate.

Gate A2 tested a deliberately narrow alternative. Independent Go, Swift, and TypeScript
implementations produced identical deterministic-CBOR bytes for one `ApprovalGrant`, exchanged
COSE_Sign1 ES256 envelopes, accepted both mathematically valid ECDSA S forms, and rejected twelve
alternate-representation and profile-confusion cases. The evidence is a conditional pass for that
bounded object, not a review of general CBOR or COSE processing.

## Proposed decision

Capsule will use RFC 8949 deterministic CBOR for canonical internal security objects and
object-specific RFC 9052 COSE_Sign1 profiles for objects that require signatures. This does not
change the strict JSON public agent protocol and does not imply that every content-addressed object
needs a signature.

Each object profile must define:

- a closed CDDL map with small unsigned-integer labels, exact required fields, exact types, and
  explicit byte/string/collection bounds;
- a unique object type, object version, signing purpose, audience, and media type;
- definite lengths and preferred deterministic encodings on the wire;
- an exact protected-header allowlist and an empty unprotected-header map;
- a bounded byte-string key ID that only selects a locally authorized key and never grants
  authority by itself;
- a byte-exact known-answer fixture and shared positive, negative, cross-language, and fuzz
  corpora.

The initial signed-envelope profile is:

- tagged COSE_Sign1 (`tag 18`) with an embedded payload;
- protected `alg` (`1`) set to ES256 (`-7`), protected `content type` (`3`) set to the exact
  object/version media type, and protected `kid` (`4`) encoded as a bounded byte string;
- no unprotected headers and no external additional authenticated data;
- an exact 64-byte IEEE P1363 `R || S` signature; DER signatures are rejected;
- verification of either valid low-S or high-S form, while signature bytes are never used as an
  object identity, replay identity, or database key.

Decoders reject duplicate map keys, indefinite lengths, non-preferred encodings, invalid UTF-8,
floats, bignums outside the profile, unknown fields or headers, arbitrary semantic tags, detached
payloads, trailing data, and values outside object-specific resource bounds. Decoding is bounded
before ordinary allocation or semantic processing. No envelope header may supply an embedded key,
URL, DID resolver, certificate chain, or other dynamic trust-discovery mechanism.

Exact registered payload bytes remain authoritative. A component may decode those bytes for
validation or rendering, but must not replace them with decode-and-re-encode output before hashing,
approval, storage, execution, or evidence composition.

The first candidate contract is `ApprovalGrant` v0 in
[`schemas/cddl/approval-grant-v0.cddl`](../../schemas/cddl/approval-grant-v0.cddl). Its field set is
the Gate A2 shape and remains pre-freeze until the final approval/session binding and identifier
formats are resolved.

## Acceptance conditions

This ADR stays Proposed until all of the following hold:

1. Every signed or canonically registered v0 object has a closed CDDL contract and byte-exact
   fixtures; object schemas are mutually exclusive.
2. Cross-language tests cover wrong object, installation, epoch, registration, audience, purpose,
   attempt, time bounds, malformed UTF-8, unknown fields, arbitrary tags, resource ceilings, and
   cross-object substitution.
3. The narrow Go, Swift, and TypeScript wrappers receive independent security review and corpus
   fuzzing, including the Swift byte-string boundary identified by Gate A2.
4. Production dependencies are pinned with provenance and update policy; experiment packages are
   not imported into product components.
5. The Supervisor and Broker demonstrate that they retain, hash, fetch, render, and verify the same
   exact registered payload bytes.

## Consequences

- Capsule avoids the observed cross-language JSON-number canonicalization failure.
- Integer labels, closed maps, and object-specific wrappers keep the signed parsing surface small.
- Public JSON and internal CBOR become separate protocol boundaries with separate strict decoders.
- Generic CBOR/COSE features remain unsupported even if a selected library implements them.
- Signature malleability is harmless only because payload bytes and domain bindings, never
  signature bytes, define semantic identity and replay state.
- The current JSON Schemas and TypeScript scaffold remain unchanged until Phase 2 replaces them.
- A failure to meet the acceptance conditions requires revisiting this ADR rather than silently
  weakening canonical-on-wire validation.

## Evidence

- [Gate A failure](../../experiments/gate-a-signing-canonicalization/README.md)
- [Gate A2 deterministic CBOR/COSE result](../../experiments/gate-a2-cbor-cose/README.md)
- [Gate A2 retained vectors](../../experiments/gate-a2-cbor-cose/fixtures/go-vectors.json)
- [Pinned production-library comparison](../../experiments/production-cbor-cose-profile/RESULTS.md),
  which conditionally passes `fxamacker/cbor` v2.9.2 for typed object encode/decode behind retained
  Capsule predecode and policy controls, and rejects `go-cose` v1.3.0 as a production envelope
  dependency. This is dependency-selection evidence only and does not close the acceptance
  conditions above.
- [Phase 2A contract foundation](../PHASE_2A_CONTRACT_FOUNDATION.md), including passive minimum
  `ExecutionPlan` and `PlanRegistration` CDDL candidates and byte-exact fixtures. These add object
  coverage but do not satisfy the production-wrapper, Swift, fuzzing, or integration acceptance
  conditions above.
