# Feasibility Spikes

Status: first-wave, Apple-credentialed, second-wave, and five Gate C implementation-readiness
tracks completed on 2026-07-31. Spike code remains non-production and may be discarded. Gate C
failed Apple Containerization as a production backend; libkrun/Hypervisor.framework remains the
lead native candidate under evaluation, but its exact profile remains unadmitted. The FD-native
root API is a `PATCH-CANDIDATE`; the direct-block-root `NullFs` disposition is `GOVERNED-PATCH`,
making removal credible but not final; and the current runtime bytes are not admissible. A
post-track independent review and
source-research reconciliation narrowed the first-slice P0 work, deferred filesystem parsing until
file artifacts, and added runtime authority as a blocker. Subsequent exact experiments rejected
stock/governed Bun and both the hardened full-Deno and tested minimal-`deno_core` constructions;
the later governed-`deno_core` follow-up physically omitted 96 nonessential built-in ops through a
one-file patch. Its package follow-up reproduced the exact snapshot and binary in two clean
same-host containers from a digest-pinned no-apt builder and complete Cargo source bundle, but
failed the stronger selection-evidence question on independent-builder, V8 notice/source, and
dynamic-runtime-root closure. The exact V8 follow-up proved the official archive-to-source and
four-patch V8 relationships but returned `SOURCE-LICENSE-CLOSURE-NO-GO` because immutable publisher
inputs, complete linked-component metadata, and generated notices are unavailable. Governed
`deno_core` is the intended engineering direction, not an admitted profile. The later TypeScript
approved-byte experiment separately passed a
narrower pre-approval transformation/binding question with exact Node 22.22.1/Amaro 1.1.5
strip-only emission; it did not select a transformer owner or runtime. `RUNTIME-001` remains
unsupported. OCI/gVisor remains an independent comparison and contingency. See the
[Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md) and
[License-free feasibility spike results](LICENSE_FREE_SPIKE_RESULTS.md) for the consolidated gate
decisions, credential-gated work, and next slice.

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
- stale same-team code denied by exact IPC but retaining a historical Keychain group, including use
  of a replacement Secure Enclave key.

Pass condition: the operating system and protocol jointly deny unauthorized peers, key use, and
storage access. Stable access-group membership must not be described as exact-build/epoch key
isolation. A conditional pass must identify which separation relies on the trusted local
administrator assumption and which update/key-rotation mitigation remains unproven.

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

### Gate C native follow-up: libkrun/HVF

Question: can one signed libkrun VMM process per attempt remove Containerization's hidden-helper
lifecycle ambiguity while preserving a native Apple hardware-VM boundary?

The retained follow-up compiled out network, disabled implicit vsock, booted a raw immutable block
root, used a trusted non-root/no-capability guest launcher, exercised resource limits, and packaged
the same-Team runner under App Sandbox. It also tested live Security.framework process identity,
concurrent exact cancellation, a durable-record-before-start pipe across three controller
`SIGKILL` points, and digest-pinned Bun 1.3.14.

Decision: conditional pass as the lead native candidate under evaluation. Five subsequent
readiness tracks proved bounded scratch/output mechanics, console capture, exact forced teardown,
and same-machine installed recovery, while also finding a mutable-path input race, an unexpected
`NullFs` virtiofs device, ambiguous runner-zero completion, incomplete output parsing/distribution,
and a no-go for the current runtime bytes. Subsequent review found that libkrun's multiport console
API may remove source/input and inline results from the raw-block path, while stock Bun's
subprocess/FFI surface adds an earlier runtime-authority admission question. Later retained P0
experiments narrowed, but did not close, those gates: the FD-native raw-root API became a
`PATCH-CANDIDATE`; a direct-block-root prototype selected `GOVERNED-PATCH` and reran the bounded
device corpus without `NullFs`; a 43-vector backend-independent framing model conditionally passed;
a 51-test local console patch established that stock console handling cannot proceed as-is; and
P0-4A conditionally passed an 18-role no-host-root topology while exposing App Sandbox, signing,
notarization, Gatekeeper, clean-host, and macOS-floor gaps. None of those results admits a backend
or permits user bytes. The later
[Deno-family disposition](../experiments/gate-c-deno-runtime-authority/RESULTS.md) also reached
NO-GO: full Deno retained initial-graph, Worker, inspector, compatibility, and persistence routes,
while `deno_core` physically registered 99 built-in ops and did not supply the TypeScript pipeline.
The bounded [physical-omission follow-up](../experiments/gate-c-deno-core-physical-omission/RESULTS.md)
then passed that one construction question: the exact governed patched construction registers and links only three
bootstrap ops with a reviewable one-file patch and reproducible ASLR-controlled snapshot. It did
not address TypeScript, independent builder provenance, runtime-profile admission, or external
isolation composition, so it does not reverse the family NO-GO or support `RUNTIME-001`.
The subsequent [reproducible-package follow-up](../experiments/gate-c-deno-core-reproducible-package/RESULTS.md)
closed the local builder-image ambiguity for a bounded two-file candidate: two clean same-host
containers using a digest-pinned no-apt builder and complete offline Cargo source bundle reproduced
the prior binary and snapshot identities. It still returned NO-GO for runtime-selection evidence
because no independent builder/host, complete archive-corresponding V8 source/notices, standalone
dynamic runtime root, or production TypeScript ownership/wiring was available.
The bounded [approved-byte follow-up](../experiments/typescript-approved-byte-boundary/RESULTS.md)
subsequently passed the exact pre-approval byte-binding question for a strip-only ESM TypeScript
subset. It selected no runtime, wired no component, and left production ownership, protocol
migration, packaging, module loading, and runtime admission open.
See
[`../experiments/gate-c-libkrun-hvf/RESULTS.md`](../experiments/gate-c-libkrun-hvf/RESULTS.md).

### Gate C implementation-readiness follow-ups

The completed comparison, coordination failures, retained evidence, contract consequences, and
remaining risk campaigns are recorded in the
[Gate C implementation-readiness synthesis](GATE_C_READINESS_CHECKPOINT.md).

The native work was split into independently reproducible tracks so evidence could be collected
without turning spike code into product code:

| Track | Decision | Evidence |
| --- | --- | --- |
| Block storage and egress | Conditional pass for raw-block mechanics; live same-user mutation means immutable input custody failed | [`RESULTS.md`](../experiments/gate-c-libkrun-storage-egress/RESULTS.md) |
| Console, timeout, and resources | Conditional pass for 4 KiB prefixes, wall/cancel scheduling, exact forced kill, and closed vCPU/RAM profiles; graceful shutdown and host CPU/memory quotas unsupported | [`RESULTS.md`](../experiments/gate-c-libkrun-console-lifecycle/RESULTS.md) |
| Installed lifecycle and recovery | Conditional same-host mechanics pass; full distribution, authority separation, session/reboot cases, and corrupt-root terminal integrity remain open | [`RESULTS.md`](../experiments/gate-c-libkrun-installed-recovery/RESULTS.md) |
| Adversarial VMM validation | Conditional fail for the exact profile because block-root adds a guest-visible `NullFs` virtiofs device | [`RESULTS.md`](../experiments/gate-c-libkrun-adversarial/RESULTS.md) |
| Runtime supply chain | Conditional feasibility pass; current bytes fail the development admission checklist | [`RESULTS.md`](../experiments/gate-c-libkrun-supply-chain/RESULTS.md) |

Later P0 checkpoints are retained separately: [FD-native root custody](../experiments/gate-c-libkrun-root-custody/RESULTS.md),
[direct-block-root `NullFs` disposition](../experiments/gate-c-libkrun-adversarial/NULLFS_P0_2_DISPOSITION.md),
[backend-independent P0-3 framing](../experiments/gate-c-p0-3-protocol-conformance/RESULTS.md),
[console correctness](../experiments/gate-c-libkrun-console-correctness/RESULTS.md), and
[P0-4A installed topology](../experiments/gate-c-installed-development-topology/RESULTS.md).

Each track owns a separate directory under `experiments/` and records observed evidence separately
from inference. Their integration permits backend-independent contract and fake-backend work, not
implementation of the real adapter against user bytes. The reconciled pre-user-byte gates are
runtime-authority closure, immutable runtime-root custody, independent `NullFs` disposition, typed
port transport/completion with bounded inline JSON, and an admissible complete installed bundle.
The bounded filesystem-image parser moves to the file-artifact gate rather than disappearing.
Passing these campaigns still would not by itself justify a `validated-local` or production claim.
See [Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md) for exit criteria and explicit pivots.

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
