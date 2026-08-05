# Technical Design

Status: agreed target architecture; pre-implementation and subject to the blocking feasibility
gates in [Feasibility Spikes](FEASIBILITY_SPIKES.md).

This document is the integration-level design for Capsule v0. Narrow companion documents own the
detailed protocol, trust, integrity, Supervisor, update, compromise, and evidence rules linked
below.

The intended macOS product packaging and staged setup/update scope is in the
[macOS installation and distribution plan](MACOS_INSTALLATION_AND_DISTRIBUTION_PLAN.md). It keeps
the user experience to one Swift application while preserving the internal daemon, Broker,
Supervisor, Source Validator launcher, Runner, and later update-role boundaries. Proposed ADR-0038
now selects a narrow on-demand Trust Coordinator as the protected-root request/record signer; it is
not installed product authority until I2B's signed corpus supports it. No Bundle Replacer is
product authority until a separate ADR and signed installed corpus support it.

Accepted ADR-0037 and the generated passive I0 profile now freeze the no-guest application's exact
seven-role tree, service/entitlement projections, lifecycle classifiers, and inactive gates. That
slice is `PASSED` only as deterministic local contract evidence; it does not build, sign, install,
register, launch, update, repair, uninstall, or admit any component.

The bounded I1A follow-up now constructs that exact tree as unsigned bytes. Its visible AppKit app
renders only typed inactive status; daemon and Supervisor entries are non-executable test-only
placeholders; the four R2 launcher/parser identities are unchanged; and independent readback
refuses missing, mixed, extra, mode-changed, manifest-changed, or substituted files. I1A invokes no
Apple signing identity and activates no entitlement, service, bootstrap, IPC, runtime, backend, or
guest. See the [I1A construction result](MACOS_INSTALLATION_I1A_UNSIGNED_CONSTRUCTION.md).

I2B2 now extends the unchanged I1A bytes with the eighth, installation-only Trust Coordinator role
and inactive Supervisor bootstrap descriptor. Its unsigned profile freezes exact identifiers,
cross-links, entitlement/constraint candidates, protected-state names, and cleanup/repair
projections while requiring `unsigned-profile-inactive` refusal. This is `PASSED` construction
evidence only; production wrapper review and separately authorized I2B3 signing/key/App Group/
service/container handoff remain `BLOCKED`. See the
[I2B2 construction result](MACOS_INSTALLATION_I2B2_UNSIGNED_CONSTRUCTION.md).

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

The same Go authority path now includes the unwired E1-E5 durable lifecycle checkpoint: one
colocated fixed snapshot retains attempts, lifecycle intents, stable fake effect/instance
identities, reconciliation, exact active/retained ceilings, and repeated-startup behavior while
`FakeBackend.CreatesGuest() == false`. Proposed ADR-0029 selects a native-fronted, in-process Go
Supervisor topology for the future authenticated boundary. No native bridge fixture, product XPC
service, production owner lock, consumer, runtime, backend effect, or guest is connected.

Proposed ADR-0033 now selects the owner-lock mechanism at design level: open and validate one
installation-root-authorized, Supervisor-created, pre-created sibling object, acquire nonblocking
BSD `flock`, and retain the opaque `CLOEXEC` descriptor for the Supervisor lifetime. Proposed
ADR-0038 selects the one-shot Coordinator/Supervisor authorization and creation ceremony. The
passive I2B1 contract now freezes its exact request/record CDDL, raw and calculated maxima,
independently generated signed fixtures, strict Go verifier, independent Swift verifier, and field
authority without activating a signer or installed handoff. The
bounded owner harness observed
process/descriptor semantics and refusal-before-store ordering only. Passive G1 adds the internal
Go/Darwin owner package using the selected descriptor-relative syscalls. G2 now composes it before
the existing v1 store and sorted no-guest recovery, uses its one owner-session ID for both store
and coordinator, permanently fences lifecycle reads/mutations after a failed held-owner check, and
closes lifecycle/store state before the descriptor. The owned-temporary-root fault/process corpus
does not wire product startup or provide the I2B installed signed handoff and protected-root
evidence.

Proposed ADR-0031 defines the next archive boundary. A complete expired
registration cohort may leave the hot snapshot only after every bound attempt is durably destroyed
with cleanup false and authoritative absence. Immutable segments retain full records and exact
registration/approval/attempt/nonce/effect/instance/replay tombstones. The finite fixed-store v2
checkpoint is selected only as a fault-injectable conformance oracle. Slice F1 now implements
passive archive types, exact limits/known-answer digests, defensive copies, and a pure complete-
cohort selector. It performs no file I/O, v2 migration, archive activation, lookup, or authority
mutation. The passive
[F2 format blocker resolution](SUPERVISOR_ARCHIVE_F2_FORMAT_BLOCKER.md) now freezes separate global/
segment index domains, typed hot/archive record locations and counts, a distinct generation-one
migration-genesis checkpoint, and generated answers. The follow-on valid-v1 mapping contradiction
is also passively resolved with a closed absent/present lifecycle union on attempt entries, a typed
lifecycle-record anchor on the present arm, and independently derived lifecycle counts. The
executable [F2 v1 mapping resolution](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md) retains the real
committed-attempt-before-lifecycle witness and exact `attempts = 1, lifecycles = 0` genesis answer.
The [stateful F2 result](SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md) now implements the owner-
asserted closed v1-to-v2 migration and read-only full verifier with deterministic bytes, complete
all-hot retained-global index reconstruction, exact migration genesis, downgrade refusal, and
pre/post-rename fault oracles. It creates no archive segment, moves no cohort, mutates no v2
authority, and calls no lifecycle adapter. Production-engine selection,
referenced-history deletion, continuous service, coherent restore activation, and rollback-
resistant non-reuse remain blocked. The
[stateful F3 result](SUPERVISOR_ARCHIVE_F3_ACTIVATION_RESULT.md) now adds exactly one sealed
immutable-segment prepare/verify/publish/activate transaction. It preserves every selected full
record and visible tombstone, publishes and directory-syncs the verified segment before the atomic
generation-two active reference, and fully reopens either the complete predecessor or successor
under retained fault, response-loss, concurrency, corruption/substitution, owner-loss, and process-
death oracles. It adds no retained lookup, v2 mutation, second segment, backup/orphan cleanup,
production engine, consumer, adapter, runtime, backend, or guest.
The [read-only F4A result](SUPERVISOR_ARCHIVE_F4A_LOOKUP_RESULT.md) now adds fresh full-verification
retained-global lookup/replay/passive-collision routing across typed hot/archive locations and
hot-only `AttemptID` recovery. It adds no v2 mutation, new effect tombstone, second segment,
consumer, adapter, runtime, backend, or guest; F4B mutation and F4C growth remain outside F4A's
passed scope.
The retained [F4B blocker](SUPERVISOR_ARCHIVE_F4B_MUTATION_BLOCKER.md) records the former
current-effect contradiction and ADR-0031's independent-source correction. The
[F4B result](SUPERVISOR_ARCHIVE_F4B_MUTATION_RESULT.md) now passes atomic fixed-store v2 authority/
lifecycle mutation, same-transaction append-only effect issuance, durable exact replay, and direct
hot/archive historical resolution. The
[F4C result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md) now passes deterministic second/later
segment activation and exact segment 64/65 bounded growth in the same fixed-store scope; F5-F6
remain deferred.

## Reference workflow

The first complete workflow is intentionally narrow:

1. An agent proposes one dependency-free, byte-exact `main.mjs` module with inline JSON input and
   bounded JSON output. Static/dynamic dependency requests, `import.meta`, CommonJS, packages, and
   loader fallbacks must refuse under accepted ADR-0034. The M1 passive foundation and bounded
   parser/process selection, passive R1 contracts, and unsigned R2 construction are `PASSED`;
   product Source Validator R3-R5B is `BLOCKED`, and plan construction remains `BLOCKED` until its
   required gates pass.
2. The Go daemon strictly decodes the proposal and, under Accepted ADR-0035/0036, sends an exact
   copied source through its private role-specific Source Validator launcher before planning. The
   launcher owns a fresh parser child; the validator parses but never
   executes source and returns only digest/length and fixed grammar-node facts. The daemon then
   resolves trusted policy and constructs canonical `ExecutionPlan` v0 bytes plus
   the canonical single-member source manifest. Registration atomically validates and retains the
   exact plan, complete bindings, manifest, and pass-through source bytes. Proposed ADR-0032's
   Source Preparer and ADR-0030's plan-v1 cutover remain conditional later TypeScript work; they no
   longer block the first release. These proposal/plan steps are target design, not current
   implementation; only the passive source-byte/SourceManifest foundation exists.
3. The daemon sends those exact bytes to the Execution Supervisor. The Supervisor independently
   validates them, applies non-overridable hard-safety rules, stores them durably, and returns a
   `PlanRegistration`.
4. The Trusted Host Broker fetches the registered bytes directly from the Supervisor, independently
   validates them, invokes its separate Broker-private launcher/fresh child over the exact copied
   source, binds
   the fixed result to digest and length, renders a bounded human-readable view, requires fresh
   user presence, and signs one attempt-bound `ApprovalGrant`. Any validator failure refuses before
   an Approval-key operation.
5. The Supervisor performs runtime-integrity preflight, atomically consumes the grant, and creates
   one `ExecutionAttempt` before any hostile side effect.
6. The selected runtime executes in a disposable development backend with no network, ambient
   environment, subprocesses, native addons, FFI, macros, inspector, package installation, or live
   host path. The exact runtime profile must prove those restrictions. Stock Bun 1.3.14, its
   governed-construction branch, hardened full Deno v2.9.4, and the tested minimal `deno_core`
   0.409.0 construction failed P0-0. Accepted ADR-0028 selects the later governed three-op
   `deno_core` construction as the first engineering candidate, and its real Deno and `rusty_v8`
   governed branches are merged. Governed `rusty_v8` PR #4 merged at head
   `80e863ddb942a4aa2b384e794fc23e35b9d2bb15` and merge
   `cbf56de2e1156b1cf1561fdbaea7172a0aa056f4`. Its exact-head workflow-dispatch run passed the
   clean Linux/arm64 network-disabled build, fixed test, corrected GN evidence query, evidence
   collection, and unsigned bundle upload. This closes the fork's bounded ARM64 construction
   blocker, but it does not admit or publish a runtime artifact. Fork-native Capsule bundle
   reconstruction, evidence review, `.mjs` source-custody/no-loader evidence, complete profile
   composition, and runtime-admission evidence remain open.
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

### Accepted disposable Source Validator architecture; blocked implementation

Accepted ADR-0035 adds one method-specific, one-shot parse-only child process for each invocation of the first-release
`.mjs` grammar policy. Exact Oxc 0.140.0 is the engineering candidate. The child receives only a
copied source frame and returns a closed typed result with recomputed digest, length, parse/policy
status, and four bounded node counts. It has no store, keys, network, paths, package/loader API,
runtime, backend, guest, or authority effect. Crash, hang, parser/semantic diagnostic, malformed
result, artifact mismatch, or sandbox failure refuses.

The daemon and Approval Broker are intended to invoke independent instances over the bytes each
owns at that stage. The parser is not linked into either parent and is never placed in the
Supervisor. V0 fixed the protocol; V1 now retains an unwired exact Oxc artifact, complete locked
inventory, V0/M1 result agreement, and same-host deterministic reproduction. The artifact remains
identity-free linker-ad-hoc-signed rather than installation-signed, not enrolled, and not invoked
by any parent. ADR-0035's design is Accepted and R1's passive role contracts are retained, but the
product boundary remains `BLOCKED` until independent provenance, installation enrollment,
sandbox/resource, consumer, broader conformance, and fault evidence pass the
[`implementation plan`](MJS_SOURCE_VALIDATOR_IMPLEMENTATION_PLAN.md). Runtime no-loader admission
remains a separate mandatory control.

R2 now retains two unsigned role-specific XPC bundle layouts and matching role-specific parser
children. Exact Oxc/Rust sources build offline in two clean same-host directories with identical
bundle bytes; source, dependency, license/notice, SBOM, static dynamic-library closure, and
unsigned provenance records are retained. The bundled R1 resource policies are inactive, so the
launchers perform exact fixed-key/frame/policy predecode and refuse without spawning. This is
construction evidence, not an installed service or active process-boundary result.

The [R3 execution packet](SOURCE_VALIDATOR_R3_EXECUTION_PACKET.md) fixes Team `3DDR84M4JS`, exact
R2 byte identities and placements, entitlement/profile requirements, the reachability/mixed-update
matrix, cleanup, and credentialed mutation boundaries. R3 remains `BLOCKED` until its containing
fixtures, exact role profiles, and finalized signed constraint bytes exist and those mutations are
separately authorized.

The exact V2 local macOS checkpoint does not advance that boundary. Its strict bootstrap refuses
before `exec` because `RLIMIT_AS` cannot be lowered; its explicitly unbounded diagnostic mutation
proves deterministic process/fault mechanics while also proving ambient file, socket, metadata-
write, and memory authority remains. Supported App Sandbox child entitlements change the fixed V1
bytes, and deprecated custom profiles are not substituted. V2 and the product Source Validator
parent are `BLOCKED` pending a new reviewed/enrolled artifact and supported resource/confinement
design.

The [supported replacement review](MJS_SOURCE_VALIDATOR_MACOS_PROFILE_REPLACEMENT.md) and Accepted
ADR-0036 make the architecture correction exact. Direct App Sandbox inheritance is `NO_GO` because
it preserves daemon/Broker static rights. The selected lower-authority design uses two distinct
private App-Sandboxed XPC launchers, one per consumer. Each owns only its matching fresh parser
child, fixed copied pipes, resource observation, process-group kill/drain/reap, and cleanup. The
role-specific services, parser identities, methods, artifact profiles, containers, and accepted
results never cross roles; neither launcher has a Supervisor/backend/key/store route.

App Sandbox grants each role-private container writable scratch authority. The product persists no
Capsule state, cache, source/diagnostic log, or reusable result there; every request, crash,
launcher restart, update, and startup requires cleanup/residue evidence. Cleanup does not prove
confidentiality or secure erasure. No public unprivileged hard memory control is usable on the
observed host, so the accepted policy is a later evidence-derived reactive footprint watermark with
one child per launcher request, bounded per-role/combined concurrency, fixed sampling cadence, and
kill/drain/reap. It is not a hard peak/exact cap or host-availability guarantee. Numeric threshold,
cadence, baseline, overshoot, and kill latency remain unset until the separately authorized signed
corpus.

The passive v1 contract and unsigned R2 construction are now `PASSED`: fixed role-distinct frames
and an explicitly inactive resource-policy shape have independent Go/Node decoders, while two
offline role-specific bundle/parser builds retain same-host byte equality and inactive
predecode/refusal with no child or authority effect. Product implementation remains `BLOCKED` on
signed installed reachability/confinement, measured resource/residue corpus, daemon consumer, and
Broker consumer. Unsupported private-XPC reachability, authority/native-loading/
filesystem/network escape, orphan cleanup, mixed update, or unacceptable measured host risk stops
the candidate rather than widening the bus.

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

The Supervisor is the sole hostile-guest launch authority. It owns registered plan bytes, the exact
first-release registered `main.mjs` manifest/bytes, registration sequences, independent hard-safety
validation, approval verification and consumption, attempts, integrity/quarantine state, backend
capability matching, staging verification, backend handles, cleanup leases, safe filesystem
collection, transcript events, and its evidence key.

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

Go remains the initial daemon language. Swift/native is preferred for the Broker because the Broker
is primarily a macOS UI, Keychain, LocalAuthentication, XPC, and file-selection component.

Proposed ADR-0029 selects one unprivileged per-user Supervisor process with a small native
C/Objective-C XPC/Security front end and the existing Go authority/lifecycle core linked in-process
through a synchronous method-specific copy-only C ABI. No Swift Supervisor service, host-root
process, or privileged helper is selected. Installed signing/session/owner-lock evidence remains
open at the product-evidence level: ADR-0033 selects the mechanism and G2 composes the local current
v1/no-guest port. Bounded G3 discovery stopped before installed build because the certificate's
common-name suffix was mistaken for a Team ID. Apple Membership Details now confirms the emitted
Team `3DDR84M4JS` is the account Team; no exact Capsule role profile is cached. I2A/Proposed
ADR-0038 now resolve the protected-root bootstrap authority and define the signed request/record;
passive fixtures, the installed handoff/container, and descriptor-relative closed-store opening
remain blockers. The installed
protected-root/session/update matrix remains unrun. Any later separate or privileged
component still requires a new ADR. Adding Rust or
another language requires a narrow interface and a demonstrated reduction in privileged risk, not
an assumption that language choice alone creates the security boundary.

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

ADR-0036's Source Validator footprint watermark is not a user-owned guest/job memory limit and does
not weaken ADR-0009. It is an evidence-bound trusted-component availability policy: the launcher
reacts after observing a threshold and therefore cannot claim an exact peak. The signed profile
must enforce its selected threshold and sampling cadence exactly as configured, record measured
overshoot and kill latency, and refuse unsupported hosts, while explicitly leaving host
availability unguaranteed.

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
out network, and disables implicit vsock. Its evidence remains conditional: a narrow raw-only
FD-native API is only a `PATCH-CANDIDATE`; the direct-block-root bootstrap selected
`GOVERNED-PATCH` and removed the prototype's `NullFs` device, but neither change passed the final
signed installed corpus. The public governed fork preserves its exact five-patch aggregate and now
also contains merged bounded console/raw-FD follow-up source at
`cf0333cdba478cc34a8570a65b38412da7fd3ecc`. That follow-up bounds the previously unchecked port-ID
and shutdown paths, fixes a queued-backpressure shutdown loop and inactive-state lifecycle defect,
and raises measured `port.rs`/`process_tx.rs` coverage from zero to 15/17 and 4/4 functions and
111/137 and 82/96 lines. It does not close the remaining 2 functions/26 lines in `port.rs`, 14
lines in `process_tx.rs`, a post-merge branch/verifier pin mismatch, the caller shared-status
hazard, independent review, or real control/queue/descriptor fuzzing and composition. The
backend-independent P0-3 framing candidate
and governed library source have not passed the real transport/launcher/guest corpus. P0-4A
conditionally supported the no-host-root topology but did not establish signed/notarized
distribution or a supported macOS floor. Neither the Bun nor governed `deno_core` evidence closes
`RUNTIME-001`; current runtime bytes are not admissible. Root custody must separately pass
attachment identity, construction of a frozen object, and end-to-end same-user attacks;
`/dev/fd/N` alone is not an immutability mechanism. Runner exit is never guest success;
guest-reported completion, input integrity, result validation, and teardown remain distinct
evidence. Undocumented directional-FD behavior and shared-open-description flag mutation remain
explicit P0 hazards rather than trusted stream semantics. Filesystem-
image parsing is a later artifact gate. gVisor resource limits bind the outer Linux worker, engine,
host cgroup/OCI configuration, and exact `runsc`/shim identity. Direct Apple Containerization remains
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

[Proposed ADR-0025](adr/0025-colocate-durable-attempt-lifecycle-state.md) selects one colocated
Supervisor snapshot/transaction domain with durable effect intents and `AttemptID`-only startup
recovery. Slices E1 through E5 now implement the passive contract, explicit fixed-store v1
migration/open validation, durable lifecycle transactions, the no-guest FakeBackend driver, exact
256-active/4,096-retained capacity behavior, and repeated-startup/exhaustion checks. Active
capacity is released only by a durable `destroyed` record with cleanup false after authoritative
absence. G2 now composes the local Darwin owner with the current v1/no-guest startup and retains the
same-session coordinator; production archive, protected installed storage, rollback, backup, and
real-backend reconciliation mechanisms remain unselected.

Proposed ADR-0033 selects BSD `flock` over POSIX process locks, macOS 26 OFD locks, and `O_EXLOCK`
after one owned local corpus. The selected opener validates the pre-created object by UID, mode,
type, link count, device, and inode relative to a retained protected state-root descriptor before
store access. Passive G1 implements that internal Go/Darwin acquisition and its local refusal
oracles. G2 adds the owner-required v1 opener/startup composition, exact sorted recovery and close
ordering, plus post-open entry-replacement fencing. It still supplies no installed same-UID
pathname protection, signed bootstrap provenance, archive composition, or product service.

[Proposed ADR-0031](adr/0031-checkpoint-closed-supervisor-cohorts.md) selects the local
conformance shape for archive and replay retention. Under the sole owner lock, the Supervisor
publishes a fully verified immutable closed-cohort segment before atomically activating a v2 hot
snapshot that references it, installs exact tombstone indexes, and removes the same records from
hot sets. Indeterminate publication or activation fences until reopen. Full cohort records remain
retained and referenced-history deletion is forbidden; fixed total caps eventually refuse. Passive
F1 types, known answers, defensive copies, and eligibility selection now exist. The passive
[F2 format blocker resolution](SUPERVISOR_ARCHIVE_F2_FORMAT_BLOCKER.md) adds scope-separated global/
segment indexes, typed locations/counts, and a distinct generated migration genesis. F2 now adds
the owner-asserted v1-to-v2 file migration and read-only empty-archive full verifier. F3 now adds
one closed cohort segment and atomic generation-two activation with full reopen verification. It
still adds no retained lookup, v2 authority mutation, second segment, consumer, or deletion. F4A
now adds only read-only retained-global lookup/replay/passive-collision routing and hot-only
`AttemptID` recovery. F4B now adds atomic authority/lifecycle mutation plus the independent
append-only tombstone source and is `PASSED` in the exact fixed-store scope documented in the
[F4B result](SUPERVISOR_ARCHIVE_F4B_MUTATION_RESULT.md). F4C now passes second/later-segment
activation and exact 64-segment bounded growth in the
[F4C result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md). A production engine plus real
power-loss proof remain deferred.

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

Implementation planning also begins with the
[ecosystem reuse and adoption map](ECOSYSTEM_REUSE_AND_ADOPTION.md). Its closed recommendation,
trust classification, authority consequences, dependency checklist, and exact consuming slice must
be recorded before adding a third-party/platform dependency or a new custom primitive. Candidate
metadata and public conformance suites are planning inputs, not evidence that a Capsule control is
implemented.

## Ordered implementation plan

1. Retain the completed architecture, claim baseline, feasibility results, and pivot decisions.
2. Continue from the completed ADR-0031/F1, stateful
   [F2 migration/full verifier](SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md), and
   [F3 first-segment activation](SUPERVISOR_ARCHIVE_F3_ACTIVATION_RESULT.md),
   [F4A retained lookup](SUPERVISOR_ARCHIVE_F4A_LOOKUP_RESULT.md), and the completed
   [F4B atomic mutation result](SUPERVISOR_ARCHIVE_F4B_MUTATION_RESULT.md), and
   [F4C bounded-growth result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md). Continue with F5
   backup/orphan/offline reporting before the bounded SQLite
   comparison. Keep the fixed snapshot as the logical oracle; do not promote it into the product
   store.
3. Continue the Source Validator replacement after accepted R0/ADR-0036 in this strict order:
   completed passive role-specific v1 contracts/fixtures, field authority, and unsigned
   two-launcher/parser construction; separately authorized signing/install and private
   reachability; confinement,
   reactive-resource, and residue corpus; daemon consumer; Broker consumer; then the M2/S1
   checkpoint. V0/V1/V2 bytes remain historical and unchanged. No threshold, sample cadence, or
   overshoot is chosen before the signed corpus, and no direct parent inheritance, shared service,
   deprecated custom profile, private API, or unbounded diagnostic mutation is a fallback.
4. Freeze the applicable signed-object set and independently review the narrow production CBOR
   wrapper. Use pinned `fxamacker/cbor` only for object-specific deterministic encoding and typed
   field decoding; retain Capsule predecode, caps, canonical-on-wire comparison, bindings, and
   replay checks. Keep `go-cose` test-only.
5. After both role-specific Source Validator consumers pass, hold the M2/S1 checkpoint: narrow
   JobProposal and generate plan-v0 registration/fetch fixtures from complete field authority.
   Only then decide whether to implement the authenticated native-to-Go Supervisor IPC bridge;
   ADR-0036 adds no Supervisor call and does not widen ADR-0029's four-call surface.
6. Connect that product path to the already-passed no-guest registration, approval/attempt, durable
   lifecycle, owner-lock G2, and archive oracle. Add composed evidence only from retained Supervisor
   state; keep `FakeBackend.CreatesGuest() == false` until runtime/backend admission.
7. Implement inline JSON ownership, bounded JSON output, and fixed agent summary.
8. Retain the exact governed Linux/arm64 handoff and passed passive C1 composition contract, then
   retain immutable C2B v1 plus its no-guest build-closure v2 successor without wiring a consumer,
   then close C2
   runtime authority, immutable root custody, `NullFs`, closed numeric descriptors,
   enforceable machine resources, typed port transport, and complete installed-bundle admission;
   do not connect user bytes to libkrun before all pass. The former upstream PR dependencies are
   merged; v2 pins the reviewed successor evidence. Every v2 archive, manifest, and artifact
   identity must be reverified before a separately authorized composed-profile/owned-guest task
   may consume it.
9. After the ADR-0028 governed `deno_core` candidate passes a separate runtime/profile admission
   ADR, add one dependency-free inline-JSON vertical slice through the admitted libkrun/HVF
   development profile, preserving Apple Containerization only as a regression fixture.
10. Add immutable regular-file snapshots, a disposable bounded filesystem-image parser, and broader
   bounded outputs.
11. Compare the exact libkrun/HVF and OCI/gVisor profiles before stronger posture; keep Apple
   Containerization explicitly development-only unless a future supported lifecycle API reopens
   its gate.
12. Operationalize production TUF/update infrastructure.
13. Evaluate optional Guardian and external witness mechanisms.

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
- Production Supervisor storage-engine selection and exact archive/compaction, implemented
  ADR-0033 owner lock plus installed protected state root, coherent-backup, power-loss,
  retention/deletion, and replay/non-reuse evidence
- Non-rollbackable or externally witnessed trust checkpoint, if required
- Retention, encryption, garbage collection, and secure-deletion limitations
- Quantitative performance and resource budgets
- Daemon aggregate request/backpressure budgets before candidate endpoint activation
- Exact platform/co-residency assumptions and microarchitectural residual-risk treatment before a
  `validated-local` claim

## External design references

The scoped comparison, reusable lessons, and deliberately rejected patterns are summarized in
[Related systems and design influences](RELATED_SYSTEMS.md). Those public references are planning
inputs, not Capsule implementation evidence or dependency selections. Exact candidate metadata,
trust/authority classification, recommendations, and consuming acceptance criteria are maintained
in the [ecosystem reuse and adoption map](ECOSYSTEM_REUSE_AND_ADOPTION.md).

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
