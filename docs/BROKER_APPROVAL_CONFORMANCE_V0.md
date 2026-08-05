# Broker rendering and approval conformance v0

## Result

Work item: passive Broker projection and strict ApprovalGrant verification

Status: `PASSED`

Scope: Supervisor-owned passive readback fixtures, public P-256 fixture key material, one local
signed known-answer vector, and unwired in-process Go harnesses. No live private key is retained.

Evidence or reason: the implementation freezes a byte-exact, ASCII-safe render projection and a
Capsule-owned tagged `COSE_Sign1` verifier. Focused tests pass exact source/input/runtime/limit/time
projection, bidi/control neutralization, source exact/cap-plus-one bounds, stale/mixed/source
substitution refusal, fixed warnings, approval-ineligible state, exact key policy, canonical
headers and payload, empty unprotected map, embedded payload, raw `R || S`, public-key signature
verification, complementary `S`, every grant/key/time binding, defensive copies, and malformed or
restored-feature refusals. Existing Supervisor tests retain one-use canonical-payload replay,
nonce uniqueness, and atomic approval-consume/attempt-create behavior.

Remaining work: authenticated installed Broker-to-Supervisor fetch, Supervisor custody/readback of
inline JSON content if the product must display it, native UI rendering, fresh `LAContext` and
Secure Enclave key creation/use, installed key authorization, same-byte sign/verify composition,
focus/cancel/synthetic-input adversarial testing, durable product wiring, and consumer activation.

Next action: separately authorize and implement the installed Broker/Approval-key slice. Any live
signing task must name the identity, host, fixtures, and exact credential-bearing operations.

Parent status: owner-only hostile-`.mjs` internal alpha is `IN_PROGRESS — TRENDING_GOOD`; Broker
product activation and internal-alpha admission are `BLOCKED`.

## Frozen implementation

`internal/execution/brokerapproval` accepts one read-only defensive Supervisor result and an
independently trusted Supervisor context. `RegistrationID` is only a locator. No rendered field is
taken from caller-supplied display text. The component decodes and cross-binds plan, role bindings,
registration, manifest, exact source, installation, epoch, Supervisor, digest, length, and expiry
before returning a projection.

The source is represented with `capsule.bytewise-ascii-escape/v0`. This is a reversible byte
encoding, not source rewriting: it prevents raw control, bidi, non-ASCII, and markup-sensitive
bytes from entering the display string while preserving the exact approved byte sequence. The
projection reports the exact inline JSON digest and length but states that the current Supervisor
readback does not contain its content. It cannot be treated as content shown to the user.

`internal/execution/approvalattempt.ProductionShapedVerifier` owns Sign1 framing, raw bounded predecode,
closed typed decode, deterministic re-encoding, exact headers, exact payload replay bytes, public
P-256 key selection, signature verification, key authorization, role binding, and time binding.
The general CBOR library is not allowed to restore open COSE behavior. The only checked-in key
material is public X/Y coordinates and a signed envelope.

The passive key authorization projection freezes Team `3DDR84M4JS`, Broker role, epoch-scoped
access group, Secure Enclave P-256 protection, `userPresence|privateKeyUsage`, fresh nonreused
context policy, active status, purpose/audience, installation/epoch/validity, and no software
fallback. Constructing it does not enroll or authorize a key.

## Claim boundary

This work contains no Swift/AppKit UI, Keychain access, LocalAuthentication prompt, signing API,
private key, Apple signing identity, service or IPC endpoint, installed protected state, approval
submission, runtime adapter, backend, VM, or guest. It does not show that a human saw or understood
the bytes, that focus was trustworthy, that synthetic input was contained, or that any product
path can accept or execute source.

The existing fixed approval/attempt store remains an unwired local mechanic. It demonstrates
one-use attempt-bound state transitions but is not newly connected to this verifier. Rendering and
key use are intentionally separate.

## Reproduction

Focused verification:

```sh
GOCACHE=/tmp/capsule-go-cache go test ./internal/execution/brokerapproval ./internal/execution/approvalattempt
```

Repository handoff verification remains the complete command set required by `AGENTS.md`.
