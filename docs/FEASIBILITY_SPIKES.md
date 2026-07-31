# Feasibility Spikes

Status: planned decision work. Spike code is non-production and may be discarded.

## Purpose

Capsule must not freeze contracts around controls that macOS, the cryptographic libraries, or an
isolation backend cannot actually enforce. These spikes answer narrow blocking questions before
broad implementation.

A spike is successful when it produces a trustworthy decision, including a negative one. It is not
judged by product polish or how much code can be reused.

## Rules for every spike

Each spike records:

- the security hypothesis and threat it addresses;
- the exact operating-system, hardware, tool, library, and backend versions;
- a minimal architecture and every privilege/entitlement used;
- positive, negative, misuse, and failure-injection cases;
- commands or tests that reproduce the observation;
- observed results separated from inference;
- unsupported behavior and residual risk;
- a pass, conditional-pass, or fail decision;
- the contract/ADR consequence and safe fallback;
- whether prototype code, fixtures, or both are retained.

Spike code must be isolated from production packages, labeled `development-only`, and excluded from
authoritative receipts or posture. Product code may not import it. If retained under `experiments/`,
the experiment README names an owner and removal/replacement condition.

## Gate A: signing and canonicalization interoperability

Question: can maintained Go, Swift, and TypeScript implementations agree on strict JSON parsing,
RFC 8785 canonical bytes, and JWS ES256 signatures without unsafe normalization?

Prototype:

- one tiny producer/verifier in each language;
- shared valid, boundary, and malformed fixtures;
- canonical-on-wire payload verification;
- protected-header allowlist with explicit object/media type and version.

Required adversarial cases:

- duplicate keys before ordinary object decoding;
- invalid UTF-8, lone surrogates, Unicode normalization differences, and UTF-16 sort-order edges;
- unsafe integers, negative zero, exponent forms, and trailing data;
- wrong object type, purpose, audience, installation, epoch, registration, and attempt;
- `none`, unknown algorithms, unknown `crit`, embedded `jwk`, `jku`, and `x5u`;
- raw versus DER P-256 signatures, malformed lengths, and the selected high-S policy.

Pass condition: all three implementations produce or verify one normative vector set with identical
payload bytes and fail the complete negative set. Otherwise document COSE/deterministic-CBOR or
another reviewed fallback in a new ADR. Do not hand-roll ECDSA.

## Gate B: macOS authority and storage separation

Question: can the daemon, Broker, and Supervisor be real security authorities rather than merely
same-user processes?

Prototype:

- separately signed daemon, Broker, and Supervisor test targets;
- OS-enforced XPC peer code requirements;
- component-specific Keychain access groups and Secure Enclave P-256 keys;
- Approval-key user-presence flow;
- separate protected data containers without a broad shared app group;
- exact build, entitlement, effective-user/session, and trust-epoch checks.

Required attacks:

- unsigned peer, same-team wrong signing identifier, stale build, copied binary, debugged peer, and
  wrong epoch;
- daemon attempts to sign with Approval and Supervisor keys;
- daemon attempts to open Broker/Supervisor stores;
- replayed connection identity and PID/path/name substitution.

Pass condition: the operating system and protocol jointly deny unauthorized peers, key use, and
storage access. A conditional pass must identify which separation relies on the trusted local
administrator assumption.

## Gate C: Apple Container capability coverage

Question: can the exact Apple Container API/configuration enforce the minimum v0 isolation
contract, or must Capsule use a lower-level API or another backend?

Prototype one dependency-free Bun JSON job and probe:

- no default or usable network interface; TCP, UDP, DNS, IPv4/IPv6, loopback, Unix socket, metadata,
  and management-channel behavior;
- read-only root and no host sockets, credentials, environment, or live user-file mounts;
- memory enforcement, CPU control semantics, PID/process-tree ceilings, wall timeout, scratch/output
  storage quotas, log/output bounds, and cancellation;
- no inherited descriptors and bounded source/input transfer;
- orphan discovery after Supervisor crash and independently observable teardown evidence.

Pass condition: every required control is enforceable and distinguishable from best-effort
accounting. Unsupported controls are removed from the v0 contract or cause a backend pivot before
schema freeze. Network inconvenience is not evidence of network denial.

## Gate D: content handles and custody

Question: can the Broker give the Supervisor only job-scoped transient access to exact content
without giving the daemon user-content access or using live user-file mounts?

Prototype:

- Broker-owned immutable content objects;
- opaque content identity exposed to the daemon;
- authenticated Broker-to-Supervisor handle transfer;
- staged digest verification and bounded returned-output handle;
- revocation/expiry, crash, duplicate redemption, and garbage-collection behavior.

Required attacks include forged identities, stale handles, cross-job reuse, daemon redemption,
path substitution, symlink/special-file inputs, partial transfer, Broker/Supervisor crash, and output
release before terminal integrity classification.

Pass condition: the daemon never receives content or redeemable authority, the Supervisor receives
only attempt-scoped access, and crash recovery cannot widen or resurrect authority.

## Gate E: Supervisor language and privilege

Question: what is the smallest maintainable component that can use macOS security APIs and control
the selected backend without unnecessary privilege?

Compare:

- native Swift Supervisor;
- Go Supervisor with narrow audited native bindings;
- a hybrid unprivileged Supervisor plus a sealed-descriptor launcher helper.

Measure platform API coverage, IPC surface, parsing/serialization TCB, deployment/update behavior,
privilege, memory/runtime footprint, testability, crash recovery, and developer cost. Do not select
Rust or another language solely as a security label; a new language must reduce a demonstrated
privileged risk through a narrow interface.

Pass condition: select the least-privileged design that can enforce the proven backend contract.
If no privileged helper is required, prohibit one in v0.

## Gate F: trust transition and recovery

Question: can an installation change components and trust state without accepting partial update,
grant replay, or silent rollback?

Prototype a signed `PreparedUpdate`, pending-verification startup, epoch finalization, and repair
workflow with fault injection before and after every durable write or component swap.

Required cases:

- old daemon/new Supervisor and the reverse;
- stale Broker, restored store, missing manifest, changed entitlement, or mismatched policy/profile
  checkpoint;
- crash before grant consumption, after consumption, during backend creation, and during epoch
  finalization;
- coherent rollback to an older locally valid epoch.

Pass condition: partial transitions enter `repair-required`, consumed grants never become unused,
and coherent rollback limitations are documented. Sequence-ordered epochs must not be called
monotonic unless backed by a non-rollbackable anchor or external witness.

## Optional Gate G: Runtime Guardian

Endpoint Security feasibility is non-blocking. Determine entitlement, root/admin, user-approval,
Full Disk Access, deployment, event coverage, performance, and false-positive characteristics.
Begin notify-only. The Guardian may contribute observations but never approval or launch authority.

## Decision record template

```text
Spike:
Revision and environment:
Hypothesis:
Threat/control:
Observed evidence:
Counterevidence and limitations:
Decision: pass | conditional-pass | fail
Contract consequence:
ADR consequence:
Retained artifacts:
Prototype disposal/replacement:
```
