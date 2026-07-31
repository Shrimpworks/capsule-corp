# ADR-0006: Signed execution plans and per-device identity

- Status: Accepted
- Date: 2026-07-30

## Context

An authenticated client request is not automatically an authorized execution. Capsule must prevent
proposal substitution, stale approval, replay, cross-device use, and changes between what the user
reviews and what the backend executes.

A digest alone detects changed bytes only when the expected digest arrives through a trusted
channel. It does not authenticate the party that authorized those bytes. Capsule also needs
device-scoped key custody, purpose separation, revocation, and a human-readable approval
experience.

W3C DID Core provides identifiers and verification relationships without requiring a blockchain.
The `did:key` method can represent a P-256 public key and resolve offline. Apple Secure Enclave can
generate and use non-exportable P-256 keys on supported hardware.

## Decision

Capsule uses a prepare, approve, execute protocol:

1. The daemon resolves an untrusted proposal into an immutable execution plan.
2. The trusted host shows a human-readable view generated from that plan.
3. The user explicitly approves the exact plan digest.
4. The daemon executes only after revalidating and atomically consuming the single-use grant.

Each Capsule installation has a per-device root identity. The initial identifier is an offline
P-256 `did:key` DID. On supported Apple hardware, the root private key is generated in the Secure
Enclave. Capsule maintains a local trust registry because a DID identifies a verification method
but does not grant authority.

The root authorizes separate operational keys for user approval, receipt assertion, and transport.
Signed grants bind type, issuer, purpose, audience, installation, subject digest, nonce, issue and
expiry times, and verification method. Unknown DID methods, algorithms, keys, purposes, audiences,
or statuses fail closed. v0 performs no network DID resolution.

The exact canonical signed-envelope representation is selected during contract freeze from a
reviewed standard with maintained Go and TypeScript implementations. Capsule will not invent a
cryptographic primitive.

## Consequences

- A changed source, input, profile, limit, output audience, backend requirement, or policy produces
  a different plan digest and invalidates approval.
- Nonces, expiry, installation identity, audience, and atomic consumption prevent straightforward
  replay and cross-context reuse.
- Human-readable presentation is part of the security boundary and must derive from the signed
  typed model.
- Device identities are not portable in v0. A replacement root creates a new DID and requires
  trusted reenrollment.
- `did:key` has no built-in update or deactivation, so Capsule owns operational-key rotation,
  revocation, and replacement records.
- Receipt signing authority cannot approve jobs.
- Verification of a signature authenticates an enrolled key; policy still decides whether the
  requested action is allowed.
- Shared deterministic digest, signature, replay, wrong-audience, and revoked-key test vectors are
  required before execution is implemented.

## References

- [W3C DID Core](https://www.w3.org/TR/did-core/)
- [`did:key` Method](https://w3c-ccg.github.io/did-key-spec/)
- [Apple Secure Enclave key protection](https://developer.apple.com/documentation/Security/protecting-keys-with-the-secure-enclave)
