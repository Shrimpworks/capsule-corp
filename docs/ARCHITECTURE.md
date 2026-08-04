# Architecture

Status: intended architecture. The repository implements only unwired local contract and
FakeBackend lifecycle mechanics; no connected authority boundary, real runtime/backend, or guest
is implemented or validated.

## System context

```text
                         Public or configured trust
                  ┌─────────────────────────────────┐
                  │ TUF trust repository            │
                  │ releases, profiles, revocations │
                  └───────────────┬─────────────────┘
                                  │ install, update, explicit refresh
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Local installation                          │
│                                                                     │
│  Untrusted agent / MCP / SDK                                        │
│             │ proposal, status, cancellation                        │
│             ▼                                                       │
│  Agent-facing Capsule daemon (Go)                                   │
│  ├── strict public protocol boundary                                │
│  ├── proposal handling and plan construction                        │
│  ├── policy resolution and immutable planning                       │
│  └── no approval key, backend launch, or user-only content          │
│             ├── exact first-release `main.mjs` pass-through source   │
│             │   └── atomic plan/source registration and readback     │
│             ├── future conditional TypeScript Source Preparer       │
│             │                                                       │
│             └── register exact plan / request registered attempt    │
│             ▼                                                       │
│  Execution Supervisor                                               │
│  ├── independently validates and stores exact plan bytes            │
│  ├── enforces hard-safety invariants                                │
│  ├── owns approval ledger and attempt lifecycle                     │
│  ├── owns backend handles, cleanup, and enforcement transcript      │
│  └── sole authority allowed to create a hostile guest               │
│       ▲                              │                               │
│       │ plan fetch, approval,        ▼                               │
│       │ job-scoped content handles  Isolation backend               │
│       │                         ├── Apple Container dev backend      │
│  Trusted Host Broker          ├── native libkrun/HVF candidate      │
│  ├── native approval UI       └── OCI + gVisor comparison/fallback  │
│  ├── user-presence key                 │                             │
│  ├── file selection                    ▼                             │
│  ├── user-content store          Disposable governed-runtime guest  │
│  └── fixed agent summary                                             │
│                                                                     │
│  Optional Runtime Guardian                                          │
│  └── observes relevant events; never authorizes execution           │
└─────────────────────────────────────────────────────────────────────┘
```

TUF processing, software downloads, DID resolution, and other network trust work do not occur in
the live execution path. A dedicated updater or trust verifier reduces verified external metadata
to a compact signed local `TrustSnapshot` consumed by the Supervisor.

## Authority model

The architecture deliberately splits three primary authorities. Proposed ADR-0032 further
subdivides the planning side with a method-specific Source Preparer that has no approval, content,
or execution authority:

| Component | May do | Must not do |
| --- | --- | --- |
| Agent-facing daemon | Authenticate agents, validate proposals, resolve policy, construct plans, register plans, request attempts, expose fixed status | Use approval/evidence keys, launch a backend, replace registered plan bytes, retrieve user-only content, clear quarantine or grant state |
| Proposed TypeScript Source Preparer | Transform copied untrusted source before planning, own exact original/emitted/profile/options/record bytes, return fixed plan projections and registered executable copies | Approve or register plans, access user-only Broker content, expose a generic parser/path API, launch a backend, transform after registration |
| Trusted Host Broker | Render registered plans, require user presence, sign one-use approvals, select files, own user content, release fixed summaries | Expose agent endpoints, launch a guest, accept daemon display prose as authoritative, make enforcement claims |
| Execution Supervisor | Independently validate registered plans, enforce hard safety, consume approvals, create attempts, manage backend lifecycle, sign enforcement transcripts | Parse the public agent protocol, author general policy, select files, perform rich parsing, fetch network trust data |

Combining Approval and Content Broker responsibilities in one native process is an accepted v0
deployment simplification, not a conceptual collapse. Their interfaces, keys, and persisted records
remain separable so a later high-assurance deployment can split them.

## Components

### Agent-facing daemon

The daemon remains the product-facing orchestration and planning service, implemented initially in
Go. It treats every proposal as untrusted after authentication. It resolves a narrow proposal,
trusted policy, content manifests, a runtime bundle, and required backend controls into canonical
`ExecutionPlan` bytes.

The daemon cannot authorize its own plan. It sends exact bytes to the Supervisor for independent
validation and durable registration. After registration, execute APIs accept only the returned
registration identifier.

Accepted ADR-0034 assigns the first release one byte-exact pass-through `main.mjs` source under the
existing plan-v0 source role. Registration atomically validates and retains exact plan, bindings,
canonical source manifest, and source bytes; Broker fetch reads defensive Supervisor-retained
copies. Only the passive source-byte/SourceManifest foundation exists. The module-request validator
is on the retained grammar-counterexample hold, so JobProposal narrowing, plan construction, source
custody, authenticated transport, consumer, runtime, backend, and guest remain unimplemented.
Proposed ADR-0032's separately enrolled
Source Preparer is now only a conditional later TypeScript topology and remains on P1 HOLD.

### Proposed TypeScript Source Preparer

If TypeScript is later selected, the proposed unprivileged Source Preparer owns exact
pre-registration Node 22.22.1/Amaro 1.1.5
strip-only emission and a single role-namespaced immutable store for original, executable, profile,
options, manifest, record, and record-set bytes. It accepts no caller path or transform option,
returns defensive copies only, and has no Approval/content/backend/evidence key or guest route.

The Supervisor resolves and independently validates a prepared source set before plan-v1
registration; the Broker independently validates copied registered objects before approval; and
attempt staging can retrieve only executable-manifest bytes through retained Supervisor state. See
[Proposed ADR-0032](adr/0032-select-enrolled-typescript-source-preparer.md).

### Trusted Host Broker

The macOS Broker is a signed native process, preferably Swift because it directly integrates with
AppKit/SwiftUI, LocalAuthentication, Keychain, XPC, and platform file-selection APIs.

The Approval interface fetches a plan from the Supervisor, independently validates and hashes it,
renders a bounded and spoof-resistant typed view, requires fresh user presence, and signs one grant
for one registration and attempt nonce. It never signs daemon-supplied prose or an opaque digest
alone.

The Content interface safely snapshots user-selected regular-file data-fork bytes, stores input and
user-only output content, and issues opaque job-scoped handles. The Supervisor necessarily receives
transient read/write handles required for staging and collection; it does not receive ambient
access to the Broker's content store.

### Execution Supervisor

The Supervisor is the primary local execution authority and the smallest, most security-sensitive
component. It independently enforces both schema validity and versioned non-overridable hard-safety
rules. For v0 those rules reject network, subprocess, environment, native-addon, FFI, package,
untrusted-image, arbitrary-path, and unsupported authority regardless of daemon policy.

It owns exact registered plan bytes, approval consumption, attempt state, runtime-integrity state,
backend capability matching, backend handles, cleanup leases, staged-digest verification, safe
filesystem collection, and a hash-linked enforcement transcript.

Proposed ADR-0029 selects one unprivileged per-user Supervisor process: a small native
C/Objective-C XPC/Security front end linked in-process with the existing Go authority/lifecycle
core. It adds no Swift Supervisor service, root LaunchDaemon, or privileged helper. The exact
installed identity/session evidence remains open. Proposed ADR-0033 selects a pre-created enrolled
sibling object held by nonblocking BSD `flock` for the process lifetime. Passive G1 implements
the internal opaque Go/Darwin acquisition capability. G2 now requires it before the current v1
store opener and sorted no-guest recovery, binds the store and coordinator to its one session,
fences on post-open ownership failure, and closes the store before releasing the descriptor. The
signed bootstrap record and installed protected-root/session/update matrix remain open, so this is
still an unwired local mechanic rather than a product authority boundary.

See [Execution Supervisor](EXECUTION_SUPERVISOR.md).

### Isolation backends

Backends implement a durable lifecycle equivalent to:

```text
probe → prepare → create → stage → start → wait/inspect
      → terminate → collect → destroy → reconcile
```

Direct Apple Containerization is a macOS development backend, not a production boundary: Gate C
found no supported durable host-side VM/helper identity or restart enumeration. The preferred
native candidate is now one libkrun/Hypervisor.framework VMM process per attempt, gated by a
durable-record-before-start handshake and verified with a PID/start/code-identity tuple. Its spike
passed lifecycle mechanics conditionally, but the integrated readiness result found unresolved
same-user block-image mutation and a guest-visible `NullFs` virtiofs device. The reconciled first
slice proposes genuine read-only descriptor custody for the trusted root plus bounded
virtio-console ports for source, inline input, completion, and inline JSON output. Those mechanisms,
the selected governed `deno_core` profile, complete installed distribution, and admissible release
bytes remain P0 hypotheses, not validated controls. Filesystem-image output parsing is deferred until
file artifacts. OCI plus gVisor remains an independent candidate and contingency. Each exact
backend reports mechanisms, unsupported controls, management channels, recovery behavior, and
retained validation evidence.

The fake backend creates no guest and exists to test plan registration, approval consumption, state
transitions, fault recovery, and evidence composition.

The current unwired E5/G2 checkpoint exercises this fake through a fixed colocated Supervisor
snapshot. It retains durable intents, stable effect IDs, cleanup/reconciliation state, exact
256-active and 4,096-retained ceilings, and idempotent repeated startup under the local Darwin
owner-required composition. This is repository-local mechanic evidence only: it is not installed
protected-root evidence, a production database, a real adapter, guest lifecycle evidence, or
evidence composition.

The development-only ADR-0033 experiment observed duplicate-process refusal, last-descriptor
release, `CLOEXEC`, enrolled file checks, and replacement limitations on one host. It selected the
mechanism. G2 now wires only the current fixed snapshot/no-guest recovery mechanic to it. Installed
protected-directory custody is a separate required boundary because advisory locks and mode bits
do not contain a same-UID process.

Proposed ADR-0031 defines the retention boundary. Only a complete expired
registration cohort whose attempts are all durably destroyed with cleanup false after
authoritative absence may move into an immutable Supervisor-owned archive segment. Exact records,
replay/non-reuse tombstones, cross-record indexes, and checkpoint digests remain retained; active
or unresolved work never archives. Passive F1 projections, exact limits/known answers, defensive
copies, and complete-cohort eligibility now exist, but they write no file, migrate no store, move no
cohort, and activate no archive. The passive
[F2 format blocker resolution](SUPERVISOR_ARCHIVE_F2_FORMAT_BLOCKER.md) now defines scope-separated
global/segment indexes, typed hot/archive record locations and exact counts, and a distinct
migration-genesis checkpoint with generated known answers. The subsequent stateful review stopped
because a valid v1 committed attempt may have no lifecycle record, a world the corrected v2
attempt index/count model cannot represent without invention. The exact
[F2 v1 mapping blocker](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md) is retained. No v2 store,
migration, or verifier was implemented; another passive format decision is required. The
selected fixed checkpoint is a finite conformance oracle, not a production engine or continuous-
service mechanism, and referenced archive history is not deletable under that proposal.

### Trust verifier and repository

TUF root, targets, snapshot, and timestamp metadata anchor release/profile distribution,
delegations, and Capsule-defined disable or revocation records. URLs and DIDs do not anchor trust.
The live Supervisor consumes only a bounded locally verified snapshot, avoiding a general TUF or
network parser in the execution TCB.

See [Trust Repositories](TRUST_REPOSITORIES.md).

### Optional Runtime Guardian

A future Endpoint Security system extension may provide notify-only observations about component
execution, signature invalidation, debugging, unexpected children, or protected-state changes. It
never grants authority and is not required for the initial point-in-time integrity posture.

## Installation Trust Domain

Each installation has:

- a random normative `installationId`;
- a hardware-backed installation root public key where supported;
- purpose-separated Approval and Supervisor evidence keys;
- optional DID representations for interoperability and exported evidence;
- a signed `InstallationManifest`;
- a sequence-ordered chain of trust epochs;
- component-specific code requirements, entitlements, storage, IPC, and Keychain access.

A trust epoch detects partial updates, stale peers, restored state, and enrolled-component mismatch.
It does not by itself prove global monotonicity or defeat coherent rollback by a privileged local
administrator. Stronger rollback evidence needs a non-rollbackable anchor or external witness.

See [Installation Trust](security/INSTALLATION_TRUST.md).

## Local process boundaries

Separate processes become separate authorities only when macOS controls support the distinction.
The feasibility work must prove:

- OS-enforced XPC peer code-signing requirements;
- expected effective user and session;
- component-specific Keychain access groups and user-presence rules;
- separate app/app-group containers and file permissions;
- no shared app group unless a narrow capability requires it;
- exact build/epoch validation and dynamic code checks;
- no daemon access to Broker or Supervisor keys and stores;
- no daemon access to backend launch interfaces.

Same-user mode bits alone are not sufficient containment against a compromised same-user process.

## Storage ownership

| Store | Owner | Contents |
| --- | --- | --- |
| Daemon store | Daemon | Proposals, untrusted agent-source copies, planning state, conditional prepared-source references, plan references, fixed agent summaries, non-authoritative receipt indexes |
| Proposed Source Preparer store | TypeScript Source Preparer | Immutable role-namespaced original/emitted/profile/options/manifest/record bytes plus bounded prepare/retain/release/quarantine state |
| Broker store | Trusted Host Broker | Input snapshots, original-path/user labels, user-only artifacts, retention policy, approval audit records, Approval key reference |
| Supervisor store | Execution Supervisor | Registered plan bytes and first-release exact `main.mjs` manifest/bytes, grant ledger, attempts, integrity/quarantine state, backend handles, cleanup leases, transcript chain, evidence-key reference, trust-epoch checkpoint |
| Trust cache | Updater/trust verifier | Pinned roots, verified TUF metadata, profile/review/revocation material, signed local trust snapshots |

Cross-store operations use explicit idempotent saga states. No component treats another store's
absence as proof that approval was unused, a backend was destroyed, or content was released.

## End-to-end control flow

```text
select/snapshot input (Broker)
        ↓ opaque content identity
propose and plan (daemon)
        ↓ exact canonical bytes
independently validate/register (Supervisor)
        ↓ registration ID
fetch/render/user-presence approve (Broker ↔ Supervisor)
        ↓ attempt-bound ApprovalGrant
preflight integrity + atomically consume (Supervisor)
        ↓ one ExecutionAttempt
stage/start/wait/collect/destroy (Supervisor ↔ backend)
        ↓ EnforcementTranscript + ArtifactManifest
validate/release content (Broker)
        ↓
composed user receipt + fixed agent summary
```

## Trust boundaries

- **Agent to daemon:** authenticated content remains untrusted; strict raw and schema limits apply.
- **Daemon to Supervisor:** authenticated peer plus exact registered bytes; Supervisor repeats
  structural and hard-safety validation.
- **Broker to Supervisor:** code-identity-authenticated IPC; approval and content handles are
  purpose-, epoch-, registration-, and attempt-bound.
- **Supervisor to backend:** generated typed configuration; no guest-controlled shell, paths,
  images, flags, or management channels.
- **Guest to host:** hostile syscalls, process, filesystem, IPC, network, and resource behavior.
- **Guest output to Broker/user/agent:** filesystem safety precedes bounded content parsing; the
  agent receives a minimized fixed-shape summary by default.
- **External trust to local install:** pinned TUF roots and verified snapshots; never a live
  execution authorization call.

## Identity and DIDs

DIDs remain a first-class interoperability feature. They can identify public organizations,
reviewers, workers, receipt signers, and verification methods. They do not grant Capsule authority.
Local key authorization binds installation, purpose, status, sequence, validity, and replacement.

`did:key` may render an operational public key offline, but its lack of update and deactivation
makes it unsuitable as the only long-lived installation recovery mechanism. Network DID resolution,
arbitrary methods, resolver plugins, and remote JSON-LD contexts do not enter local v0 approval or
execution.

## Evidence and posture

The user-facing receipt composes a cryptographically attributable Approval Broker claim, a
cryptographically attributable Supervisor enforcement transcript, the registered plan, and the
artifact manifest. These signatures do not independently attest that the host, UI, or signer logic
was uncompromised.

Posture is multidimensional:

- isolation assurance;
- runtime-integrity evidence mode;
- trust freshness;
- distribution authority.

No single `secure` or `authoritative` label hides unsupported dimensions.

## Portability

The independent adapter boundaries are:

- client adapter: MCP, CLI, SDK, or future HTTP;
- runtime adapter: governed `deno_core` as the first engineering candidate, with Node or other
  reviewed profiles remaining future portability/contingency options; Bun is retained only as
  rejected historical evidence for the first profile;
- isolation backend: fake lifecycle, Apple development-only Containerization, native
  libkrun/Hypervisor.framework candidate, OCI/gVisor comparison/contingency, or future microVM;
- platform key/UI/IPC provider: native macOS first, later other operating systems;
- external trust provider: default Capsule TUF repository or pinned self-hosted authorities.

Public job semantics should remain portable. Source compatibility, platform controls, and security
posture are not assumed identical.

## Detailed design

See [Technical Design](TECHNICAL_DESIGN.md), [Threat Model](security/THREAT_MODEL.md),
[Trust Architecture](security/TRUST_ARCHITECTURE.md), and [Protocol Object Model](protocol/OBJECT_MODEL.md).
