# Technical Design

Status: agreed target architecture; pre-implementation and subject to the blocking feasibility
gates in [Feasibility Spikes](FEASIBILITY_SPIKES.md).

This document is the integration-level design for Capsule v0. Narrow companion documents own the
detailed protocol, trust, integrity, Supervisor, update, compromise, and evidence rules linked
below.

## Current versus intended state

The repository currently contains buildable Go and TypeScript scaffolding. It does not yet launch a
guest or implement the authorities in this design.

The JSON Schemas in `schemas/` and types in `packages/protocol` are canonical for the current
scaffold only. They are explicitly pre-freeze and inconsistent with the intended object model in
several ways: the current `Job` mixes proposal and effective authority, exposes guest paths, and
models capabilities that v0 will not support. The current receipt also lacks composed Broker and
Supervisor evidence. Phase 2 replaces these contracts after the platform and interoperability
spikes determine honest enforceable semantics.

The Phase 2A contract-foundation slice now provides passive candidates for a narrow first
`JobProposal` and minimum `ExecutionPlan`/`PlanRegistration` payloads. Their fixtures and decoded
Go/TypeScript views are verification artifacts only: no product endpoint consumes them, and the
minimum plan omits unresolved controls and cannot authorize execution.

The current unwired implementation also retains an exact plan-registration conformance handoff.
TypeScript prepares only a defensive copy of provenance-bearing constructed plan bytes plus a
complete separately issued role-binding set. A local test command passes those inert values to the
Go `registrationstate` component, which independently predecodes, strictly decodes, role-binds,
hashes, registers, and re-reads the exact bytes. The command's JSON wrapper is test serialization,
not a local IPC contract or product transport. No authenticated TypeScript/Go process seam, daemon
consumer, Broker, approval, real backend, or guest is connected.

## Reference workflow

The first complete workflow is intentionally narrow:

1. An agent proposes a dependency-free, one-shot TypeScript transformation with inline JSON input
   and bounded JSON output.
2. The Go daemon strictly decodes the proposal, resolves trusted policy and manifests, and creates
   canonical `ExecutionPlan` bytes.
3. The daemon sends those exact bytes to the Execution Supervisor. The Supervisor independently
   validates them, applies non-overridable hard-safety rules, stores them durably, and returns a
   `PlanRegistration`.
4. The Trusted Host Broker fetches the registered bytes directly from the Supervisor, independently
   validates them, renders a bounded human-readable view, requires fresh user presence, and signs
   one attempt-bound `ApprovalGrant`.
5. The Supervisor performs runtime-integrity preflight, atomically consumes the grant, and creates
   one `ExecutionAttempt` before any hostile side effect.
6. The selected runtime executes in a disposable development backend with no network, ambient
   environment, subprocesses, native addons, FFI, macros, inspector, package installation, or live
   host path. The exact runtime profile must prove those restrictions. Stock Bun 1.3.14, its
   governed-construction branch, hardened full Deno v2.9.4, and the tested minimal `deno_core`
   0.409.0 construction failed P0-0. No real runtime may enter this step until another exact
   construction closes the unchanged authority contract.
7. The Supervisor applies filesystem-safety collection and the Broker performs bounded content
   validation and user delivery.
8. The Supervisor destroys or explicitly classifies unresolved backend state and signs an
   `EnforcementTranscript`.
9. A user receipt composes the registered plan, Broker approval evidence, Supervisor transcript,
   and artifact manifest. The agent receives only a fixed `AgentExecutionSummary` by default.

The task remains the public abstraction. Runtimes, process trees, containers, VMs, and management
channels remain backend details.

## Formal invariants

### Execution invariant

A hostile guest may be created only by an enrolled Execution Supervisor after it independently
verifies a previously registered immutable plan, an active one-use attempt-bound approval, matching
installation and trust epoch, trusted profile state, current runtime-integrity evidence, and a
backend able to enforce every required control.

### Approval invariant

The Broker renders and approves the same typed plan bytes registered by the Supervisor. It never
signs daemon-supplied display text or an opaque digest alone.

### Hard-safety invariant

The Supervisor independently rejects unsupported protocol power and enforces versioned
non-overridable safety rules. A compromised or mistaken daemon cannot enable network, subprocess,
environment, native addon, FFI, package, untrusted-image, arbitrary-path, or unsupported backend
authority in v0.

### Data invariant

An agent cannot turn a path, URL, environment variable, image reference, secret name, DID, or
arbitrary identifier into authority. A trusted action creates a capability for exact immutable
bytes or a narrowly defined operation.

### Attempt invariant

One approval authorizes at most one attempt. Once the Supervisor consumes the grant, a crash,
timeout, or indeterminate launch never restores it to unused state.

### Integrity-failure invariant

A required runtime-integrity failure after approval consumption terminates or quarantines the
attempt, blocks automatic artifact release, and prevents ordinary success classification.

### Evidence invariant

A receipt contains cryptographically attributable Broker approval and Supervisor enforcement
claims bound to the same registration, plan digest, installation, trust epoch, and attempt. It does
not claim independent hardware/platform attestation or prove human understanding, guest
correctness, kernel integrity, or absence of encoded exfiltration.

## Component and authority model

### Agent-facing daemon

The Go daemon is the public protocol and planning service. It may authenticate agent clients,
strictly validate proposals, store agent source, resolve policy and trusted references, construct
plans, register exact plan bytes, request registered attempts, expose fixed status, and request
authorized cancellation.

Before the candidate public endpoint activates, its transport and work scheduler enforce a bounded
aggregate envelope for connections, concurrent requests, in-flight bytes, queues, deadlines,
cancellation, downstream stalls, and diagnostics. Strict object budgets are per-request controls,
not a substitute for service-wide backpressure and deterministic overload behavior.

It may not:

- access Approval, installation-root, or Supervisor evidence private keys;
- create or modify key authorizations or trust epochs;
- launch a backend or reach a privileged launcher directly;
- provide new plan bytes at execute time;
- clear quarantine, repair-required, epoch, or grant-consumption state;
- receive user-only input or artifact content by default;
- forge Supervisor terminal evidence or declare success without it.

### Trusted Host Broker

The Broker is a native signed macOS process with no agent-facing interface. One v0 process may host
two logically separate surfaces:

- **Approval Broker:** fetches registered plans, independently validates and hashes them, renders
  typed fields, protects against spoofing and hostile labels, requires LocalAuthentication/Keychain
  user presence, and signs one-use approval grants.
- **Content Broker:** performs trusted file selection, creates immutable snapshots, owns user labels
  and original paths, stores user-only content, transfers attempt-scoped handles to the Supervisor,
  validates bounded output content, and releases a fixed summary to the daemon.

UI activation, focus, and synthetic input are not approval evidence. The Approval-key signature
requires the configured LocalAuthentication/Keychain-gated private-key operation over the exact
registered binding. Accessibility, Screen Recording, overlay/window automation, and comparable
broad user-granted capabilities use an explicit elevated-adversary posture rather than an
unsupported claim that ordinary UI hardening defeats them.

The Broker never launches the backend. Approval and content interfaces, keys, and records remain
separable even while deployed in one process.

### Execution Supervisor

The Supervisor is the sole hostile-guest launch authority. It owns registered plan bytes,
registration sequences, independent hard-safety validation, approval verification and consumption,
attempts, integrity/quarantine state, backend capability matching, staging verification, backend
handles, cleanup leases, safe filesystem collection, transcript events, and its evidence key.

It does not own public agent parsing, policy authoring, file-picking UX, arbitrary DID/TUF/network
resolution, rich content parsing, profile building, dependency installation, or general updater UX.
New responsibilities require an ADR.

See [Execution Supervisor](EXECUTION_SUPERVISOR.md).

### Operating-system enforcement

Process separation is a design claim only after the macOS spike proves:

- XPC peer code requirements for each trusted message/channel;
- component signing identifier, team identity, effective user/session, entitlements, and exact active
  build identity;
- component-specific Keychain access groups and user-presence rules;
- separate protected storage containers and narrow filesystem access;
- no daemon route to Broker/Supervisor keys, stores, or backend-control interfaces.

Same-user file modes alone are not a sufficient authority boundary. A shared app group is forbidden
unless a documented narrow capability requires it.

### Least privilege and languages

Go remains the initial daemon language. Swift is the preferred Broker language because the Broker
is primarily a macOS UI, Keychain, LocalAuthentication, XPC, and file-selection component.

The Supervisor language and privilege model are deferred until the feasibility spike. The project
will compare a native implementation, Go with narrow platform bindings, and an unprivileged
Supervisor plus a tiny sealed-descriptor launcher if required. It does not assume root. Adding Rust
or another language requires a narrow interface and a demonstrated reduction in privileged risk,
not an assumption that language choice alone creates the security boundary.

## End-to-end protocol

### 1. Content capability creation

The v0 file path, when added after inline JSON, is:

1. The user chooses a file through the Broker.
2. The Broker opens it without following links, verifies a bounded regular file, and copies the
   opened data-fork bytes into private storage while hashing.
3. The copied bytes are authoritative. Mutation checks are a race/usability signal, not a promise
   that the original path remained stable.
4. The Broker records an `InputSnapshotManifest` and returns only an opaque capability reference
   and safe manifest metadata to the daemon.
5. Original path, resource forks, extended attributes, and other filesystem semantics never enter
   the agent or guest contract.

Directories, symlinks, devices, sockets, FIFOs, archives, and host-side rich parsing are outside v0.

### 2. Proposal and planning

Before normal object decoding, the daemon enforces raw bytes, media type, strict UTF-8, duplicate
keys, nesting, string length, collection counts, numeric range, one-document, and trailing-data
rules.

The narrow `JobProposal` contains source, an inline JSON or opaque input slot, requested user-owned
limits, an output slot, and a trusted profile alias. It contains no network grants, subprocess
names, environment inheritance, packages, native powers, arbitrary image, backend flags, host
paths, guest paths, or agent-supplied general JSON Schema.

The daemon resolves content manifests, source identity, policy, runtime bundle and registry entry,
exact requested/effective limits, fixed audiences, required backend controls, installation, and
trust epoch into canonical `ExecutionPlan` bytes.

### 3. Plan registration

The daemon sends exact canonical bytes through authenticated local IPC. The Supervisor repeats raw
limits and schema parsing, recomputes canonical bytes and digest, verifies referenced objects, and
enforces the v0 hard-safety rules.

If acceptable, it durably stores the bytes and returns a `PlanRegistration` binding:

- opaque registration ID and sequence;
- plan digest;
- installation ID and trust epoch identifier;
- Supervisor identity;
- expiry.

Attempt requests contain only the registration ID. Replacement bytes are never accepted.

### 4. Approval

The Broker connects directly to the Supervisor and fetches the registered bytes. Its approval view
shows source identity, complete input-read authority, runtime/profile posture, network/process/
environment state, exact limits, declared output, user and agent observation channels, and the
warning that generated code can encode inputs through permitted outputs or metadata.

Untrusted labels are length-bounded, escaped, Unicode-bidi-safe, and visually separated from
trusted UI. The UI does not claim that the user understands the source. A plan that cannot be
rendered completely and safely cannot be approved.

After fresh user presence, the Broker signs an `ApprovalGrant` binding the plan digest,
registration, installation, epoch, expected Supervisor, attempt nonce, purpose, audience, issue
time, expiry, and boot/session context where reliable.

### 5. Runtime-integrity preflight

Before grant consumption, the Supervisor creates a `RuntimeIntegrityAssessment` that verifies the
installation manifest and active epoch, its own expected dynamic code state, authenticated peers,
key authorization, local trust snapshot freshness, runtime profile/review/registry state, exact
backend implementation and controls, durable ledger health, and absence of quarantine or
repair-required state.

This creates a short internal permission to proceed. It is not a public attestation. Without an
independent Guardian, evidence is point-in-time and is labeled accordingly.

### 6. Grant consumption and attempt creation

In one durable Supervisor transaction:

1. verify the active, unused, unexpired, correctly scoped grant;
2. verify the registered bytes remain exact;
3. create the unique attempt;
4. bind the grant to the attempt;
5. mark the grant consumed;
6. commit before hostile side effects.

A failure after this transaction burns the approval. Retrying requires a new approval.

### 7. Stage and execute

The Supervisor assigns all fixed guest paths, obtains attempt-scoped content handles, verifies exact
bytes and digests, verifies the runtime/backend identities, creates a fresh guest, and starts the
declared entrypoint without shell interpolation. In the first inline slice it delivers registered
source and input through dedicated bounded attempt-bound virtio-console ports. Later file slices may
stage dedicated volumes only after their custody and parser gates pass.

The trusted runtime root is immutable under a proven host-custody mechanism, not only guest-read-
only. For the libkrun candidate, P0 separately proves stable attachment identity, frozen-object
construction, and adversarial end-to-end custody. The finalized root digest is computed through the
exact retained genuine read-only descriptor only after every writable alias/mapping is closed and
the sole pathname is unlinked. Input transport, scratch, completion/result transport, and any later
artifact output are separate and bounded. The Supervisor observes time, resources, cancellation,
backend state, and transcript events.

### 8. Collect and release

For the first slice, a trusted launcher writes exactly one fixed-cap typed attempt-bound completion
frame containing bounded inline JSON to a dedicated port that the unprivileged workload cannot own.
Before implementation, the protocol freezes separate exact source, canonical-input, completion-
frame, and JSON-payload caps plus per-channel role, version, attempt, registration/plan, runtime-
profile, length, digest, terminal-status, and commit-trailer semantics. The launcher writes the
commit trailer only after the complete payload; every channel fails rather than resizing. The host
continuously drains cap-plus-one and never treats stream EOF as completion. The Supervisor validates
framing, binding, limits, JSON, and the separate runtime, input, runner-lifecycle, result, and
teardown dispositions. Runner exit status is never workload success, and the accepted record is
guest-reported completion—not attestation of an uncompromised guest kernel or correct execution.

When file artifacts are added, the Supervisor filesystem gate accepts only declared fixed output
slots, bounded counts/bytes, regular files, safe fixed paths, and expected storage semantics. A
disposable bounded filesystem-image parser rejects symlinks, hard-link tricks, devices, sockets,
FIFOs, sparse-file abuse, malformed metadata, and undeclared output. The Supervisor never delivers
content directly to the agent.

The Broker content gate performs bounded UTF-8/JSON and later JSONL/text/CSV validation, protects UI
rendering from formula, terminal, bidi, and HTML hazards, stores user content, and returns either a
fixed agent summary or content separately authorized for that audience.

Rich document, spreadsheet, PDF, archive, image, audio, or media interpretation belongs in another
disposable parser sandbox—not the daemon, Broker authority core, or Supervisor.

### 9. Destroy and compose evidence

Every path after backend creation passes through terminate/destroy/reconcile. Missing backend state
never proves destruction.

The Supervisor signs a terminal `EnforcementTranscript`. The user receipt composes the registered
plan, approval and signer authorization, integrity assessment summary, exact backend/profile/control
identities, transcript, artifact manifest, result classification, teardown state, and optional
privacy-reviewed trust/witness proof.

The daemon may package and index this material but cannot manufacture either embedded authority
claim.

## Protocol object model

The object model separates:

- public untrusted proposal and fixed summary objects;
- content-addressed source, input, policy, backend, and artifact objects;
- installation, key, profile, validation, and trust-snapshot objects;
- plan registration and authorization objects;
- attempt, integrity, event, transcript, and receipt evidence.

`JobProposal`, `ExecutionPlan`, `PlanRegistration`, and `ExecutionAttempt` are distinct concepts and
must never be used interchangeably. Logical slot identifiers replace agent-provided guest paths;
Capsule assigns paths such as `/capsule/inputs/0/data` internally.

See [Protocol Object Model](protocol/OBJECT_MODEL.md).

## Cryptographic and serialization profile

Gate A rejected the tested RFC 8785/JWS direction after a real Swift/Foundation number
representation disagreement. Gate A2 conditionally passed a bounded deterministic-CBOR/COSE_Sign1
profile across Go, Swift, and TypeScript. ADR-0019 therefore proposes this internal security-object
baseline, subject to its production hardening conditions:

- SHA-256 content identities;
- hardware-backed P-256 ECDSA keys where supported;
- RFC 8949 deterministic CBOR with closed, object-specific CDDL maps;
- tagged RFC 9052 COSE_Sign1 envelopes with embedded payloads;
- ES256 with exact 64-byte `R || S` signatures on the wire;
- exact protected headers for algorithm, object/version media type, and bounded byte-string key ID;
- an empty unprotected-header map and no dynamic key-discovery headers;
- explicit object type, schema version, purpose, audience, installation, epoch, and state binding.

Mandatory rules include:

- reject duplicate keys, indefinite lengths, invalid UTF-8, unknown tags/fields, floats, and
  non-preferred encodings before ordinary object processing;
- verify canonical-on-wire payload bytes rather than silently normalizing alternate encodings;
- constrain integers to an exactly interoperable range or encode large counters as canonical
  decimal strings;
- reject algorithm mismatch, unprotected headers, detached payloads, embedded keys/certificates,
  URLs, DIDs, and unrelated key-discovery mechanisms;
- reject DER signatures, accept either mathematically valid ECDSA S form, and never use signature
  bytes as object or replay identity;
- keep mutually exclusive signed object schemas and object-specific signing purposes;
- use shared byte-exact, numeric, signature, malformed, resource-bound, and cross-object confusion
  fixtures;
- retain exact registered payload bytes as authoritative instead of replacing them with
  decode-and-re-encode output.

Capsule will not implement ECDSA primitives. The narrow CBOR/COSE wrappers and their dependencies
must be reviewed and fuzzed before ADR-0019 can become Accepted. The public agent API remains strict
JSON with a separate decoder and schema boundary.

Not every content-addressed object is signed. Source, plan, policy, and artifact manifests can use
hash identity because signed/registered parent objects bind them. Distributed runtime bundles,
review claims, key authorization, installation epochs, approvals, Supervisor transcripts, TUF
metadata, and witness checkpoints require purpose-appropriate signatures.

## Identity, DIDs, and key authority

The normative internal installation identity is a random `installationId` plus locally enrolled
public keys. The installation root performs rare trust ceremonies and is not available to the
daemon or used for routine receipts.

Operational keys are purpose-separated:

| Key | Custodian | Purpose | User presence |
| --- | --- | --- | --- |
| Installation root | Broker/install ceremony | Operational-key and epoch transitions | Required for v0 trust-changing ceremonies |
| Approval key | Approval Broker | `capsule.plan.approve` | Required for each v0 plan |
| Supervisor evidence key | Supervisor | `capsule.execution.attest` | No; narrow noninteractive use |
| Optional content key | Broker | Snapshot/delivery claim across a process boundary | Normally no |

DIDs are first-class external identifiers for organizations, reviewers, workers, portable signer
references, and exported receipts. They do not confer authorization. An offline `did:key` may
render an operational public key, but Capsule policy still binds that key's installation, purpose,
validity, sequence, status, and replacement. Because `did:key` cannot rotate or deactivate, it is
not the sole long-lived installation root or recovery mechanism.

Local v0 authorization performs no network DID resolution, arbitrary method loading, resolver
plugin execution, or remote JSON-LD context retrieval.

See [Trust Architecture](security/TRUST_ARCHITECTURE.md).

## Installation Trust Domain

An `InstallationManifest` binds the installation ID/root, expected component roles and exact code
requirements, relevant entitlements, operational-key authorizations, runtime-profile registry and
policy digests, pinned trust root/checkpoint, storage formats, prior epoch digest, and transition
reason.

Every component-changing install, update, repair, or authority transition produces a
sequence-ordered trust epoch. Plans, registrations, approvals, attempts, trusted IPC, and receipts
bind its number and digest.

Epochs detect partial update, stale peer, enrolled-component mismatch, and many restored-state
conditions. They do not prove monotonicity against coherent privileged rollback. The design uses
the term `sequence-ordered` unless a non-rollbackable anchor or external witness supports the
stronger claim.

See [Installation Trust](security/INSTALLATION_TRUST.md) and
[Update and Recovery](UPDATE_AND_RECOVERY.md).

## External trust repositories

Production release/profile distribution uses TUF-style root, targets, snapshot, timestamp, and
delegated roles. Pinned root metadata—not a URL, DID, TLS certificate, or `valid: true` response—is
the trust anchor.

An updater/trust verifier performs network I/O and full TUF processing outside the live execution
TCB. It emits a compact bounded signed `TrustSnapshot` containing the exact locally accepted
release/profile/review/validation/revocation state and freshness classification. The Supervisor
consumes that snapshot without fetching or resolving network data.

TUF distributes Capsule-defined revocation and emergency-disable objects; TUF itself is not a
general authorization or revocation policy. Local execution can use cached verified state under an
explicit freshness policy. Unavailability never causes acceptance of unsigned or rollback data.

See [Trust Repositories](TRUST_REPOSITORIES.md).

## Runtime integrity

The initial macOS posture combines OS-enforced XPC peer code requirements, dynamic code validity,
exact enrolled build identities, expected entitlements, effective user/session, and shared trust
epoch. PID, path, and process name alone are never identity.

Runtime-integrity states are:

```text
uninitialized → verifying → verified → ready → executing
                         ↘ degraded | quarantined | repair-required | compromised
```

- `degraded` means optional evidence or freshness is insufficient for the requested posture.
- `quarantined` means material identity, debug, backend, or state mismatch.
- `repair-required` means an interrupted or inconsistent trust transition.
- `compromised` means strong evidence of unauthorized modification.

The daemon cannot clear these states. Evidence records whether observations were startup-only,
point-in-time, periodically revalidated, continuously monitored, or platform-attested. Only
implemented mechanisms may use the corresponding label.

See [Runtime Integrity](security/RUNTIME_INTEGRITY.md).

## Policy and exact limits

Policy precedence is:

```text
Capsule non-overridable hard-safety invariants
  → required organization policy and pinned trust roots
  → local administrator policy
  → user defaults and ceilings
  → one-plan user approval
  → agent request
```

Each layer may narrow lower authority and may not widen a higher safety invariant. Missing values
resolve to trusted defaults during planning. Requests above a ceiling are rejected rather than
clamped. The exact resolved values appear in the plan and approval view. A backend enforces each
value exactly or refuses execution.

The Gate C readiness corpus supports a narrow limit vocabulary: exact wall time, per-stream retained
console prefix, integer vCPU topology, closed workload-specific guest-RAM profiles, bounded port
frames, physical scratch-image bytes, and—only for later file artifacts—artifact count,
per-artifact logical bytes, total artifact logical bytes, and separate parser limits. CPU
percentage/time, arbitrary RAM, and exact total host/VMM memory are unsupported. Backend-independent
vocabulary can freeze before the remaining P0 gates; exact libkrun values freeze only after their
mechanisms pass.

## Runtime bundles and profile evidence

Profile trust uses separate objects:

- `RuntimeBundleManifest`: immutable runtime, root/image, dependencies, lock, SBOM, provenance,
  hardening configuration, exact guest-kernel image/configuration/boot/module/debug policy,
  launcher restrictions, required controls, and build identity;
- `ProfileReviewAttestation`: signed reviewer claim about one exact bundle and attack-corpus
  evidence, with verdict, date, expiry, and limitations;
- `ProfileRegistryEntry`: mutable local alias, activation state, accepted review authorities,
  limits, and backend compatibility;
- `BackendValidationRecord`: explicit validation verdict and posture ceiling for exact bundle/
  backend/host/configuration claims, known limitations, expiry/invalidation triggers, and retained
  evidence. A P0 `development-admitted` verdict cannot authorize `validated-local`.

Verifying a bundle signature attributes the exact manifest to an accepted publisher key; it does
not activate or validate the bundle. The first bundle contains no third-party guest packages.

## Backend contract and posture

Each backend exposes a bounded `BackendCapabilityReport` identifying exact implementation and
binary digest, platform requirements, network and IPC mechanisms, filesystem controls, root
immutability, resource controls, process-tree kill, storage quotas, orphan discovery, teardown
evidence, management channels, unsupported controls, and validation-record digests.

The Supervisor matches the plan to supported mechanisms; the backend cannot self-assert an
authoritative tier. The native libkrun/HVF candidate binds one signed App-Sandboxed VMM process to
one attempt, records and verifies its PID/start/code identity before authorizing VM start, compiles
out network, and disables implicit vsock. Its current spike is conditional evidence only: the
pathname disk API has an unresolved same-user mutation race, the block-root path creates a
`NullFs` virtiofs device, and neither the Bun nor tested Deno-family constructions close
`RUNTIME-001`; current runtime bytes are not admissible. P0 will test a genuine inherited read-only
root descriptor through `/dev/fd/N` and dedicated virtio-console ports for source/input and
fixed-cap typed inline results. Root custody must separately pass attachment identity, construction
of a frozen object, and
end-to-end same-user attacks; `/dev/fd/N` alone is not an immutability mechanism. Runner exit is
never guest success; guest-reported completion, input integrity, result validation, and teardown
remain distinct evidence. The pinned multiport implementation's unchecked guest port IDs, non-stop-
aware output wait, undocumented directional-FD convention, shared-status mutation, and partial-
then-error handling are explicit P0 hazards, not trusted stream semantics. Filesystem-image parsing
is a later artifact gate. gVisor resource limits bind the outer Linux worker, engine, host cgroup/OCI
configuration, and exact `runsc`/shim identity. Direct Apple Containerization remains
development-only because it has no supported durable VM/helper identity or restart reconciliation
surface.

Guest-kernel minimization is profile-bound defense-in-depth and never replaces a VMM corpus that
assumes a hostile kernel. Any validated platform profile also binds its exact hardware model,
OS/hypervisor build, documented vendor mitigation state, and concurrency/co-residency policy. Those
facts support configuration matching, not a claim of speculative-execution or shared-cache
noninterference; residual microarchitectural risk remains explicit unless a separate exact campaign
supports a narrower statement.

Posture is recorded across independent dimensions:

- isolation: `development`, `validated-local`, `hosted-hardened`, `high-assurance`;
- runtime integrity: `startup-only`, `preflight-point-in-time`, `periodically-revalidated`,
  `continuously-monitored`, future `platform-attested`;
- trust freshness: `current-online-verified`, `current-cached`, `offline-grace`, `stale-degraded`;
- distribution identity: official Capsule, enterprise/custom authority, or local development.

## Egress and observation channels

Guest output is an attacker-controlled channel. Capsule cannot prove noninterference: code that can
read an input may encode it into permitted content, byte counts, timing, status, or repeated
behavior.

The default agent response has a fixed shape such as:

```json
{
  "state": "completed",
  "attemptId": "attempt_...",
  "receiptId": "receipt_..."
}
```

Exact artifact names, sizes, user paths, guest messages, previews, rich violation strings, timings,
and resource metrics remain user-only unless a separate policy explicitly grants them. Trusted
logs contain opaque IDs, bounded codes, build identities, and approved numeric metrics—not paths,
content, source excerpts, environment values, guest strings, approval text, secrets, or keys.

The channel budget reduces leakage; it does not eliminate state/timing leakage or prove
confidentiality.

## Persistence, transactions, and recovery

Storage ownership follows the authority split. Content is bounded access-controlled storage rather
than unbounded SQLite blobs. The Supervisor keeps the authoritative grant ledger and backend state;
the daemon never reconstructs them from receipt indexes.

Operations spanning stores use explicit saga records and idempotent messages. Required recovery
rules include:

- approval consumption and attempt creation commit before backend side effects;
- content staging and output release bind attempt and content digests;
- user artifact release waits for terminal integrity and collection state;
- a crash after consumption never restores a grant;
- a missing backend handle is not proof of destruction;
- a partial trust transition enters `repair-required`;
- repeated cleanup records unresolved or teardown failure instead of success;
- garbage collection never deletes the only record needed to revoke, reconcile, or explain an
  active attempt.

[Proposed ADR-0025](adr/0025-colocate-durable-attempt-lifecycle-state.md) narrows the next fake-only
implementation boundary to one colocated Supervisor snapshot/transaction domain with durable
effect intents and `AttemptID`-only startup recovery. The proposal and its conformance plan do not
make that lifecycle durable today or select production archive, locking, rollback, backup, or
backend reconciliation mechanisms.

## Error and violation taxonomy

Stable bounded codes distinguish malformed/unsupported protocol, authentication, untrusted or
wrong-purpose key, stale/replayed/expired/wrong-audience approval, policy denial, profile/trust
failure, integrity/quarantine failure, backend capability mismatch, guest/resource failure,
content/egress rejection, cancellation, recovery, and teardown failure.

Human-readable messages do not echo secrets, host paths, approval text, or arbitrary guest output.
Agent-visible classifications come from a smaller fixed allowlist than user/admin diagnostics.

## Verification strategy

The immediate work has two lanes: backend-independent product implementation against a hard-fenced
fake backend, and the bounded fail-fast P0 program in the
[Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md). Each remaining spike has a security
hypothesis, exact environment, adversarial cases, pass/fail criteria, retained evidence, contract
consequence, ADR decision, and disposal rule.

After contract freeze, every implementation layer uses shared positive and negative fixtures.
Every backend candidate then runs the same corpus across filesystem, process, network/IPC,
resources, runtime features, outputs, identity, approval, handoff, component substitution,
recovery, update, and cross-job state.

No claim becomes implemented or validated merely because a mechanism is designed or a prototype
passes a happy path. See [Control Evidence Matrix](security/CONTROL_EVIDENCE_MATRIX.md).

## Ordered implementation plan

1. Retain the completed architecture, claim baseline, feasibility results, and pivot decisions.
2. Freeze backend-independent contracts using the measured results.
3. Implement registered-plan, approval-ledger, fake-backend, crash-recovery, and composed-evidence
   lifecycle using a locally seeded development trust snapshot.
4. Implement inline JSON ownership, bounded JSON output, and fixed agent summary.
5. In parallel, close runtime authority, immutable root custody, `NullFs`, typed port transport,
   and complete installed-bundle admission; do not connect user bytes to libkrun before all pass.
6. After P0-0 selects an admitted runtime by ADR, add one dependency-free inline-JSON vertical
   slice through the admitted libkrun/HVF development profile, preserving Apple Containerization
   only as a regression fixture.
7. Add immutable regular-file snapshots, a disposable bounded filesystem-image parser, and broader
   bounded outputs.
8. Compare the exact libkrun/HVF and OCI/gVisor profiles before stronger posture; keep Apple
   Containerization explicitly development-only unless a future supported lifecycle API reopens
   its gate.
9. Operationalize production TUF/update infrastructure.
10. Evaluate optional Guardian and external witness mechanisms.

See [Roadmap](ROADMAP.md) for exit evidence.

## Blocking decisions

- Production review and freeze of the bounded deterministic-CBOR/object-specific COSE profiles
- Distribution-package/session validation of the per-user macOS process, storage, XPC, Keychain,
  and privilege topology
- Exact libkrun/HVF runtime-authority restrictions, immutable root custody, `NullFs` disposition,
  typed port transport/completion, resource, App Sandbox, installed-bundle, identity, and recovery
  mechanisms before user bytes; disposable filesystem-image parsing before file artifacts
- Exact OCI/gVisor worker, engine, runtime, resource, management, identity, and recovery mechanisms
- Broker-to-Supervisor content-handle mechanism
- Non-rollbackable or externally witnessed trust checkpoint, if required
- Retention, encryption, garbage collection, and secure-deletion limitations
- Quantitative performance and resource budgets
- Daemon aggregate request/backpressure budgets before candidate endpoint activation
- Exact platform/co-residency assumptions and microarchitectural residual-risk treatment before a
  `validated-local` claim

## External design references

- [W3C Decentralized Identifiers v1.0](https://www.w3.org/TR/did-core/)
- [W3C Decentralized Identifiers v1.1](https://www.w3.org/TR/did-1.1/)
- [`did:key` Method](https://w3c-ccg.github.io/did-key-spec/)
- [The Update Framework metadata](https://theupdateframework.io/docs/metadata/)
- [RFC 8785: JSON Canonicalization Scheme](https://www.rfc-editor.org/rfc/rfc8785)
- [RFC 7515: JSON Web Signature](https://www.rfc-editor.org/rfc/rfc7515)
- [RFC 8725: JSON Web Token Best Current Practices](https://www.rfc-editor.org/rfc/rfc8725)
- [Apple XPC peer requirements](https://developer.apple.com/documentation/xpc/xpc_connection_set_peer_requirement)
- [Apple dynamic code validity](https://developer.apple.com/documentation/security/seccodecheckvalidity%28_%3A_%3A_%3A%29)
- [Apple Keychain user presence](https://developer.apple.com/documentation/localauthentication/accessing-keychain-items-with-face-id-or-touch-id)
- [Apple Container](https://github.com/apple/container)
- [Apple Containerization](https://apple.github.io/containerization/documentation/containerization/)
- [gVisor security model](https://gvisor.dev/docs/architecture_guide/security/)
