# ADR-0043: Freeze Broker rendering and approval verification

- Status: Accepted
- Date: 2026-08-05
- Refines: ADR-0010, ADR-0011, ADR-0019, ADR-0021, ADR-0024, ADR-0034, and ADR-0040
- Decision scope: owner-only internal-alpha Broker projection, Approval-key policy, and
  Supervisor verification contract

## Context

ADR-0040 requires the native Broker to fetch and render Supervisor-owned typed bytes, obtain fresh
human presence for one Approval-key operation, and produce a one-use attempt-bound grant. The
repository already retains passive atomic plan/source custody and an unwired durable
approval/attempt mechanic. It did not freeze the render projection, exact Approval-key policy, or a
Capsule-owned cryptographic consumer.

`RegistrationID` is caller-copyable and therefore cannot supply display authority. Host text can
contain controls, bidirectional text, or markup-sensitive bytes. UI focus and synthetic input are
not evidence of approval. General COSE/CBOR libraries also do not provide Capsule's closed headers,
canonical-on-wire bytes, role binding, key authorization, time policy, or replay identity.

The current Supervisor readback retains the exact plan, complete role bindings, registration,
canonical source manifest, and exact `main.mjs`. It does not retain inline JSON content bytes. The
first projection must expose that absence rather than manufacture or fetch content.

## Decision

### Read-only projection

The Broker treats `RegistrationID` only as an untrusted locator. It renders nothing until one
Supervisor readback has been strictly decoded and the following agree: locator and retained
registration; exact-plan digest; installation identity; epoch sequence and digest; Supervisor
identity; plan and registration expiry; complete resolved role bindings; source-manifest digest,
entrypoint, member digest, and exact source length.

The bounded v0 projection contains:

- exact registration identity/sequence, plan digest, installation, epoch, and Supervisor identity;
- exact source profile, media types, `main.mjs` entrypoint, one-member count, digest, manifest
  digest, length, content policy, and exact content;
- inline JSON slot, digest, length, and the explicit fact that content bytes are not available or
  shown by this readback;
- exact runtime alias and bundle, review, registry, backend-validation, backend-configuration,
  trust-snapshot, and policy-decision digests;
- supported wall-time and JSON-output limits, their origin and slots, and plan/registration expiry;
  and
- fixed warnings about data-bearing outputs/timing, human-presence limits, focus and synthetic
  input, missing input content, separate runtime syntax/loader enforcement, and lack of product
  admission.

Source display is the reversible ASCII-only `capsule.bytewise-ascii-escape/v0` encoding. Safe
printable ASCII is retained; backslash/newline/carriage-return/tab use named escapes; markup
delimiters `<`, `>`, `&`, double quote, single quote, and backtick plus every other byte use
lowercase `\xhh`. The source remains strict UTF-8 and bounded to 262,144 bytes; escaped display is
bounded to 1,048,576 bytes. Raw control, bidi, non-ASCII, normalization-sensitive, or markup bytes
never reach the display string.

The projection records that focus, activation, and synthetic input are not approval evidence. A
cancel, focus loss, or session replacement invalidates the pending interaction; replacement also
requires a fresh Supervisor fetch. This is a contract only. Until installed UI evidence and an
actual key operation exist, the projection is explicitly approval-ineligible.

### Approval key and human-presence contract

The Approval key is exactly a nonexportable Secure Enclave P-256 key when supported. Its policy is:

- Team ID `3DDR84M4JS`;
- role `capsule.role.approval-broker/v0`;
- access group `3DDR84M4JS.com.capsulecorp.capsule.broker.approval.epoch-` followed by the
  canonical base-10 epoch sequence;
- explicit `userPresence` and `privateKeyUsage` access control;
- purpose `capsule.plan.approve` and audience `capsule.execution-supervisor`;
- installation, epoch sequence/digest, validity, active status, public key, and key ID bound by the
  trusted authorization identity; and
- no software-key fallback.

Each signing operation creates one fresh, nonreused `LAContext`. It is invalidated after success,
cancel, failure, focus loss, timeout, or session replacement. Rendering and key use remain separate:
untrusted display state cannot select a key or signing payload. Unsupported Secure Enclave or
access-control behavior refuses approval rather than weakening policy.

This ADR freezes that contract but does not authorize live Keychain or LocalAuthentication work.
Live signing requires a separate explicit task.

### Capsule-owned COSE_Sign1 verification

The Supervisor uses a narrow Capsule-owned wrapper and standard-library P-256 verification. The
only accepted object is a tagged RFC 9052 `COSE_Sign1` with:

- a four-element array;
- deterministic protected map keys in exact order: algorithm `-7`, content type
  `application/capsule.approval-grant+cbor;v=0`, and a 32-byte key ID;
- an empty unprotected map, embedded exact payload, empty external AAD, and a 64-byte raw `R || S`
  signature; and
- the closed canonical 12-field `capsule.approval-grant` v0 payload.

Raw predecode runs before typed decode. Indefinite, nonpreferred, duplicate, unknown, out-of-order,
detached, DER, tagged-inner, trailing, over-depth, over-collection, over-byte, or noncanonical forms
refuse. Verification accepts every mathematically valid P-256 signature, including complementary
`S`; payload replay identity is the exact canonical payload rather than signature bytes.

The authenticated payload must exactly bind installation, epoch digest, registration, plan,
Supervisor, attempt nonce, purpose, audience, issuance, and expiry. Trusted authorization must
also bind Team, role, access group, key/public key, installation, epoch sequence/digest, purpose,
audience, validity, protection, access control, context policy, active status, and no fallback.
Effective Supervisor time must be within the grant, registration, and key intervals; expiry
equality refuses.

One-use behavior remains Supervisor-owned: canonical-payload replay is idempotent, retained nonces
are unique, and consuming one approval plus creating its immutable `AttemptID` is one durable
transaction. Failure after consumption burns the approval. The Broker never consumes approval or
creates an attempt.

## Consequences

- Spoof-resistant passive rendering and public-key-only cryptographic verification can be reviewed
  without key access or UI activation.
- Inline JSON content is not represented as shown or approved until Supervisor custody/readback is
  extended by a coordinated contract change.
- The wrapper deliberately does not add `go-cose`; deterministic CBOR remains behind Capsule raw
  predecode, closed typed decode, canonical re-encoding, and role binding.
- No live Keychain, `LAContext`, prompt, signing identity, service/IPC endpoint, approval activation,
  runtime, backend, VM, or guest is created by this decision or its passive conformance harness.
- Installed same-byte render/sign/verify evidence, authenticated Broker/Supervisor IPC, protected
  key authorization, UI spoof/focus/cancel testing, and product wiring remain `BLOCKED`.
