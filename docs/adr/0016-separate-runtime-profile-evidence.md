# ADR-0016: Runtime bundle, review, activation, and backend validation are separate

- Status: Accepted
- Date: 2026-07-31

## Context

Publisher signing establishes artifact origin/integrity but does not prove security review,
activate a profile, or show that one backend/host enforces the required controls. Combining mutable
review/activation state with immutable bundle identity creates ambiguity and invalidates caching.

## Decision

Capsule separates:

- immutable publisher-signed `RuntimeBundleManifest`;
- immutable `ProfileReviewAttestation` for one exact bundle/evidence set;
- mutable local `ProfileRegistryEntry` for alias and activation policy;
- `BackendValidationRecord` binding exact bundle/backend/host/configuration and retained evidence;
- machine-readable `BackendCapabilityReport` describing current mechanisms/unsupported controls.

Activation requires accepted publisher, reviewer, registry, trust snapshot, backend capability, and
validation state. The first executable bundle has no third-party guest packages.

## Consequences

- A signed malicious bundle is not automatically executable.
- Review can expire/revoke without changing bundle bytes.
- Validation claims cannot float across backend or host configuration.
- More objects and signatures require clear object/type/purpose separation.
- Friendly aliases always resolve to exact immutable identities before planning/approval.
