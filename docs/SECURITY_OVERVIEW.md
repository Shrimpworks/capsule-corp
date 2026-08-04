# Capsule security overview

Capsule is trying to make one narrow kind of local work safer to authorize: a bounded JavaScript
or TypeScript job proposed by an AI agent. The goal is not to give the agent a computer, shell, or
general development environment. The goal is to describe one task exactly, ask a person to approve
that exact task, and run it without giving either the agent or its code broad access to the host.

> [!WARNING]
> Capsule is not yet a usable security boundary. The repository contains local lifecycle mechanics,
> design requirements, and retained experiments, but it does not yet provide a complete path for
> running hostile code. Do not use it for that purpose unless another trusted isolation system
> supplies the boundary.

This document is for readers who are comfortable with software concepts but are not expected to
know security engineering, cryptography, virtualization, or macOS internals.

## How to read status statements

Capsule keeps four kinds of claim separate:

- **Implemented mechanics** are code and tests in the repository. They may still be internal,
  unwired, and unsuitable for production.
- **Experiment evidence** is a result observed in a named, retained test environment. It supports
  only the mechanism and conditions that were tested.
- **Selected design** is an accepted or proposed requirement for the intended system. It is not a
  claim that the requirement has been implemented.
- **Open production work** is required before Capsule can offer a usable hostile-code boundary or a
  stronger security posture.

These labels appear throughout the document because “designed,” “observed in a prototype,” and
“implemented in the product” are materially different claims.

## The problem: ambient authority

**Selected design.** A normal program launched under a user's account often inherits *ambient
authority*: access that comes from its surroundings rather than from an explicit grant for the
current task. Examples include the user's files, environment variables, network, credentials,
other processes, and installed developer tools. Agent-generated code should not receive all of
that merely because it needs to transform one JSON value.

Capsule instead treats the **job** as the public unit. A job describes source, input, exact limits,
and a declared result. A runtime, process, virtual machine, or container is an internal way to
perform the job, not something the agent gets to configure directly.

**Selected design.** The first intended executable slice is deliberately small: dependency-free
JavaScript derived from a narrow TypeScript subset, inline JSON input, and bounded JSON output. It
has no network, subprocess, inherited environment, native add-on, foreign-function interface,
package installation, arbitrary image, or agent-selected host or guest path.

**Open production work.** No complete executable slice exists today. The current daemon scaffold
does not start a guest, and no runtime or isolation backend has been admitted to handle user bytes.

## Layers, not one magic barrier

**Selected design.** Capsule assumes any single defense can have a bug or be used in the wrong
context. It therefore combines separate authorities, exact human authorization, runtime reduction,
external virtual-machine isolation, operating-system process controls, bounded data channels,
durable recovery state, and attributable evidence. No layer is described as sufficient by itself.

```text
untrusted agent
      │ proposes a narrow job
      ▼
daemon ── constructs exact typed plan bytes
      │ registers those bytes
      ▼
Execution Supervisor ◄── Approval Broker fetches and displays the registered plan
      │                    then signs one attempt after user presence
      │ consumes approval and records the attempt before effects
      ▼
runtime adapter ── gives the runtime only its fixed source/input/result channels
      ▼
isolation backend ── contains a disposable hostile guest outside the runtime
      │
      ▼
bounded result + cleanup state + composed evidence
```

The runtime's own permissions are defense in depth. An in-process JavaScript sandbox shares the
host process and its native runtime implementation, so Capsule does not treat it as the boundary
between hostile code and the host operating system. External isolation remains mandatory even if
runtime-level powers have also been physically removed.

## Exact plans and human approval

Except for the paragraph explicitly labeled **Implemented mechanics**, this section describes the
selected design rather than a connected product path.

### From a request to registered bytes

**Selected design.** An agent submits an untrusted proposal. The daemon resolves trusted defaults,
policy, source and input identities, a runtime profile, backend requirements, audiences, and exact
resource limits into a typed `ExecutionPlan`. “Typed” means each field has a defined role and
shape; a source digest cannot be silently treated as an input digest, for example.

The daemon sends the exact plan bytes to the Execution Supervisor before approval. The Supervisor
independently decodes them, applies non-overridable safety rules, stores the exact bytes, and returns
a `PlanRegistration` containing a Supervisor-issued registration ID and the plan digest.

An execution request later supplies only that registration ID. It has no place to provide new plan
bytes, backend flags, images, mounts, or paths. This removes an entire substitution opportunity
between approval and execution.

### Digests make byte changes visible

A **SHA-256 digest** is a 32-byte fingerprint calculated from exact input bytes. SHA-256 is a hash
function, not encryption: it does not hide the input and there is no key that decrypts it. Its job
here is to make byte identity compact and checkable.

For example, if an approved source file contains `return 1;` and one byte changes it to
`return 2;`, its source digest changes. That changes the plan digest, so the old registration and
approval do not apply. The changed bytes require a new plan, registration, and approval.

A digest alone does not say who created or approved the bytes. It becomes meaningful only when a
trusted record or signature binds that digest to a role such as “the plan registered by this
Supervisor.”

### The approval screen cannot trust daemon prose

**Selected design.** The Approval Broker fetches registered plan bytes directly from the
Supervisor, validates and hashes them independently, and renders their typed fields. It must show
the complete relevant authority: source and input identity, runtime/profile posture, network and
process state, exact limits, declared output, and who can observe the result.

The Broker does not display a daemon-supplied sentence such as “convert this harmless file.” A
compromised daemon could make that sentence reassuring while registering a different plan. The
Broker also refuses a plan it cannot render completely and safely; an opaque digest alone is not a
useful human approval screen.

### One approval, one attempt

**Selected design.** After a fresh macOS user-presence check, the Broker signs one `ApprovalGrant`.
The grant binds the plan digest, registration, installation, trust epoch, expected Supervisor,
purpose, audience, expiry, and a fresh attempt nonce. The Supervisor verifies those bindings and,
in one durable transaction, marks the approval consumed and creates one unique attempt before any
backend effect.

A **nonce** is a fresh value used once to make replay detectable. It is separate from the approval
ID and attempt ID. IDs name records; they do not prove that a caller is trusted. Durable state must
remember used nonces and consumed approvals, or a crash could make an old grant look unused.

**Implemented mechanics.** The repository has unwired Go mechanics for exact plan registration,
candidate approval submission, atomic approval consumption and attempt creation, replay and
concurrency handling, durable effective-time high-water, and attempt-keyed recovery. The approval
verifier accepts retained fixtures only. It is not production cryptography, Broker signing, user
presence, or authenticated interprocess communication.

## Who owns which authority

**Selected design.** Splitting code into processes matters only if the operating system and
protocol enforce the split. The intended roles are:

| Role | What it may do | What it must not do |
| --- | --- | --- |
| Agent-facing daemon | Authenticate clients, validate proposals, resolve policy, build and register plans, request attempts, return fixed status | Approve, use approval or evidence keys, read user-only content, launch a backend, replace registered bytes, clear consumed grants or quarantine |
| Approval and Content Broker | Fetch and render registered plans, require user presence, sign approvals, select files, hold user content, release a fixed agent summary | Expose an agent endpoint, launch a guest, trust daemon display prose, make Supervisor enforcement claims |
| Execution Supervisor | Independently validate and store plans, verify and consume approvals, create attempts, enforce hard safety, own backend lifecycle, cleanup, quarantine, and enforcement evidence | Parse the public agent protocol, invent general policy, select files, run rich parsers, fetch network trust data |
| Runtime adapter | Translate one admitted plan into a fixed runtime invocation and fixed data channels | Add packages, ambient environment, shell interpolation, paths, descriptors, or powers not present in the plan |
| Isolation backend | Create, observe, stop, destroy, and reconcile one externally isolated guest using typed configuration | Treat its own capability report as proof of validation, accept guest-controlled configuration, infer cleanup from a missing response |

The Approval and Content Broker may initially share one native macOS process, but their interfaces,
keys, and records remain logically separate. A backend adapter is separate from the isolation
backend so the same task contract can be carried by different runtimes and isolation mechanisms.

**Open production work.** The signed Broker, deployed Supervisor, authenticated role-specific
interfaces, protected production stores, content handles, and real runtime/backend adapters do not
exist as a connected product path.

## What the macOS mechanisms contribute

These mechanisms solve different problems. None of them proves that the signed program is correct,
that a virtual machine cannot escape, or that a user understood generated code.

### Code signing and code requirements

macOS **code signing** attaches identity and integrity metadata to executable code. A
**designated requirement** is the system's expression of the identity associated with signed code;
Capsule expects stricter explicit code requirements that also constrain the accepted signer/team,
component identifier, role, distribution channel, entitlements, and enrolled build identity.

**Selected design.** Trusted local peers will check those requirements on actual messages and
running processes. Exact code-directory hashes identify the enrolled active builds; the shared
installation trust epoch supplies the protocol state that code signing does not represent.

**Experiment evidence.** Apple Development and Developer ID fixtures rejected unsigned,
same-team/wrong-role, stale-build, and debug-entitled peers under the tested exact requirements.
That establishes useful platform primitives, not a shipping topology. A valid signature still does
not prove correct logic, current authorization, or exclusive key control.

### Hardened Runtime

The macOS **Hardened Runtime** is a collection of launch and runtime restrictions around signed
code. Among other things, it can make debugging, code injection, and loading untrusted dynamic
libraries harder unless an entitlement explicitly permits them.

**Selected design.** Shipping Capsule components are intended to use Hardened Runtime without the
development `get-task-allow` entitlement, and preflight will reject checked debugged or dynamically
invalid components.

Hardened Runtime is not a filesystem, network, process, or virtual-machine sandbox. It also cannot
make correctly signed but flawed code trustworthy.

### App Sandbox

The macOS **App Sandbox** limits which operating-system resources an app process can access based
on its signed entitlements and user-granted access.

**Selected design.** Capsule expects the virtual-machine monitor runner and trusted native
components to use narrow, role-specific sandboxes and protected containers. The daemon must not
share the Broker's or Supervisor's key groups or authority-bearing storage.

**Experiment evidence.** Signed test components used distinct protected containers, and a
Developer ID libkrun runner ran under App Sandbox in one earlier experiment. The later complete
topology harness lacked a valid signing identity, so its sandboxed runner stopped before `main`.
Final entitlements, inherited descriptor checks, protected construction, and clean-host behavior
remain open. App Sandbox does not revoke file descriptors inherited at launch and does not replace
hostile-guest isolation.

### Keychain, LocalAuthentication, and user presence

The macOS **Keychain** stores keys and secrets under access-control and signed-entitlement rules.
**LocalAuthentication** is the system interface for checks such as Touch ID or a device credential.
Capsule uses **user presence** to mean that the operating system allowed a protected key operation
after the configured local check; it does not mean the user understood the program.

**Selected design.** Approval, installation-root, and Supervisor evidence keys have different
purposes and custodians. The daemon has none of their private keys. Each approval requires a
user-presence-gated Approval-key signature over the exact registered binding; the Supervisor's
evidence key is noninteractive but can sign only enforcement objects under policy.

**Experiment evidence.** Provisioned fixtures denied daemon and sibling access to the Broker and
Supervisor Keychain groups. A persistent Secure Enclave P-256 elliptic-curve approval key refused
noninteractive use and signed after interactive presence. The experiments also found that a stable
Keychain group does not distinguish an old build from a newly enrolled build. The proposed
mitigation uses a new group and non-migrated key for an identity-changing trust epoch; shipping
install, locked-Keychain, restore, migration, and power-loss behavior remain open.

### XPC and audit tokens

**XPC** is a macOS interprocess communication system. An **audit token** is operating-system
metadata attached to a process or message that includes facts such as effective user and login
session. Capsule expects XPC peer code requirements to reject the wrong signed role before it can
invoke a trusted method, then uses typed protocol checks for role, purpose, installation, epoch,
registration, and attempt.

**Experiment evidence.** Test services applied exact code requirements in both directions,
revalidated the sender of the actual message, and checked effective user and audit session. Stale
clients and services were refused. An audit token is an identity observation, not authorization:
the expected user or process ID still cannot approve a plan or select a method outside its role.

**Open production work.** The product has no authenticated XPC interface today. Cross-user,
fast-user-switching, logout/login, in-flight retry, update, and durable epoch cases remain untested
in the final topology.

### SMAppService and launchd

`SMAppService` is Apple's API for registering an app's helper service. **launchd** is the macOS
service manager that starts and restarts processes in a system or per-user domain.

**Selected design.** Capsule's current direction is an unprivileged per-user Supervisor registered
inside the app, with no root daemon or privileged helper. launchd supplies lifecycle and activation;
it does not approve jobs or decide which backend operation is authorized.

**Experiment evidence.** A per-user three-role service experiment showed on-demand activation,
exact bidirectional peer checks, and a new Supervisor instance after a crash. A later 18-role
topology used embedded `SMAppService` registration and explicit activation without root, but pure
on-demand activation and the final signed package remain open.

### Hypervisor.framework and HVF

Apple's **Hypervisor.framework**, often shortened to **HVF**, exposes hardware virtualization to a
macOS process. Capsule's lead native candidate uses libkrun, a virtual-machine monitor library, so
one signed libkrun process owns one disposable Linux guest for one attempt.

**Selected design.** The Supervisor will record and verify the exact runner process identity before
a private control descriptor authorizes virtual-machine start. The runner is intended to expose
only a closed device and descriptor set, with network and implicit host communication paths absent.

**Experiment evidence.** The libkrun/HVF work supports continued engineering: exact-process
lifecycle, bounded resource mechanics, forced teardown, some App Sandbox behavior, and controlled
guest probes were observed. It did not validate the final device surface, runtime, installed
package, or complete hostile-guest boundary. Hardware virtualization does not prove the virtual-
machine monitor has no exploitable bug, that the guest kernel is trustworthy, or that
microarchitectural side channels are absent.

### Notarization and Gatekeeper

**Notarization** is Apple's automated submission and ticketing process for distributed software.
**Gatekeeper** assesses whether downloaded software is signed and accepted under the system's
distribution policy; a stapled notarization ticket lets that result travel with the app.

**Selected design.** The final complete app, nested services, runtime, libraries, entitlements, and
minimum operating-system build will be signed, notarized, stapled, assessed, and read back on clean
hosts before a development profile can be admitted.

**Experiment evidence.** One isolated Supervisor fixture reached an accepted, stapled, Gatekeeper-
accepted result. The complete Capsule topology has not. Notarization and Gatekeeper do not review
the correctness of Capsule's approval logic, runtime restrictions, isolation, or cleanup.

## Cryptography without mystery

### Hashing is byte identity, not secrecy

Capsule uses SHA-256 to content-address source, input, plans, policy, runtime bundles, artifacts,
and evidence records. **Content addressing** means referring to data by a digest of its exact bytes
rather than by a mutable path or friendly name. Hashes detect substitution when the expected digest
arrives through a trusted binding. They do not encrypt data, prove its origin, or make a mutable
file immutable.

### Signatures attribute claims to authorized keys

A **digital signature** lets a verifier check that a party holding a private key signed exact
bytes. Capsule's selected direction uses purpose-specific signatures for approvals, Supervisor
transcripts, installation changes, profile review, and external trust metadata.

A mathematically valid signature is not enough. Capsule must also verify that the public key is
enrolled for this installation, object type, purpose, audience, epoch, validity period, and status.
A Supervisor evidence key cannot approve a plan, and an Approval key cannot authorize a software
update. Signatures support authenticity and attribution; they do not encrypt the signed content or
prove the signed statement is true.

**Experiment evidence.** Go, Swift, and TypeScript agreed on a deliberately bounded deterministic
Concise Binary Object Representation (CBOR) and CBOR Object Signing and Encryption (COSE)
signed-object candidate. The governing ADR remains Proposed: production wrappers, complete object
profiles, Swift integration, review, and fuzzing are open.

### IDs, nonces, purpose, and audience prevent reuse in the wrong place

An installation ID, registration ID, approval ID, attempt ID, and effect ID name different records.
An attempt nonce is fresh replay material. Capsule treats these as separate types even though some
have the same byte length, so a value from one role is rejected in another.

Signed objects also state their **purpose** and **audience**. Purpose says what kind of authority is
being claimed, such as plan approval. Audience says which component is meant to accept it. Together
with installation, epoch, registration, attempt, and expiry bindings, these fields stop a valid
object from being reused across roles or contexts. The durable ledger supplies the other half of
replay prevention by remembering what has already been consumed.

### Installation identity, trust epochs, and key separation

**Selected design.** A random installation ID plus locally enrolled public keys is the internal
trust root. A decentralized identifier (DID) may identify a key or public principal for exported
evidence, but an identifier is not trust. Local authorization records decide what the key may do;
Capsule performs no live DID lookup in approval or execution.

A **trust epoch** is a sequence-ordered record of the installed components, keys, policy, and trust
checkpoints that are expected to work together. Plans, approvals, attempts, and trusted IPC bind
the same epoch so a partial update or stale peer fails closed. Epochs are not inherently rollback-
proof: a privileged party that coherently restores all older local state may need an independent
checkpoint or witness to be detected.

### Evidence and receipt chains

**Selected design.** The Supervisor will record bounded, hash-linked enforcement events and sign a
terminal transcript. A user receipt will compose that transcript with the Broker's separate
approval, the registered plan, runtime/backend identities, artifact manifest, result classification,
and teardown state. The agent gets a much smaller fixed summary by default.

Hash links reveal removed, reordered, or substituted records when the trusted head is known.
Purpose-separated signatures make the Broker and Supervisor claims attributable. They do not make
the receipt independent platform attestation: a compromised authorized Supervisor could still lie
in its own transcript, and a signed approval cannot prove comprehension.

## Why Capsule is governing its runtime and backend

General-purpose runtimes and virtual-machine libraries are built for many users. Capsule needs a
much smaller, reviewable surface whose prohibited powers are structurally absent or externally
denied, plus exact source, build, update, and evidence ownership. “Governed” means Capsule carries
reviewed forks or patches with pinned sources, reproducible builds, tests, publication materials,
advisory responsibility, and a plan to rebase or remove the changes. It does not mean the fork is
automatically safe.

### Governed `deno_core` and `rusty_v8`

`deno_core` is a low-level Rust wrapper around the V8 JavaScript engine. `rusty_v8` supplies the
Rust bindings and V8 build artifacts. Full Deno and stock Bun exposed powers outside Capsule's v0
contract. An initial minimal `deno_core` construction still registered 99 built-in operations.

**Experiment evidence.** A one-file governed patch physically reduced that registry to the three
bootstrap operations required by the fixed experiment and caught deliberate restoration changes.
The resulting binary and snapshot were reproduced in clean same-host containers. A separate
22-entry Linux runtime root ran the fixed fixture without ambient host libraries or configuration.
These are narrow construction results, not a complete runtime.

**Selected design.** Accepted ADR-0028 chooses governed `deno_core` as the first runtime engineering
candidate. The real Deno and `rusty_v8` forks have governed commits. Governed `rusty_v8` PR #4 is
unmerged external work in progress at head
`80e863ddb942a4aa2b384e794fc23e35b9d2bb15`; its clean Linux/arm64 build and fixed test
pass, and the corrected GN evidence-query diagnostic passes. One complete clean bundle run,
evidence review, and merge remain. It has no accepted artifact or admission effect. Exact V8
build/source/notice closure, independent builders, closed `.mjs` module loading, external
isolation composition, and runtime admission remain open.

**Experiment evidence.** A fixed Node/Amaro strip-only TypeScript experiment produced deterministic
JavaScript for a narrow syntax subset. Proposed ADR-0026 requires both the original TypeScript and
emitted JavaScript, plus the exact transformer identity, to be bound before plan registration and
approval. No production component owns or performs that transformation today. PR #72 kept the
proposed Source Preparer's P1 contracts `BLOCKED` until protected-store, worker-confinement,
genesis/update, retention/release, recursive field-authority, and lifecycle evidence closes.
Accepted ADR-0034 now freezes the first release as one byte-exact pass-through `main.mjs`, without
CommonJS, static/dynamic dependencies, package resolution, legacy Node module authority, or a
module-loader fallback. TypeScript and its Source Preparer remain conditional later work.

### libkrun/HVF, immutable roots, and direct block attachment

libkrun is the lead native Apple isolation candidate because one virtual-machine monitor process
can be tied to one attempt and reconciled after the Supervisor crashes. Its current direction uses
HVF for hardware virtualization and a standalone Linux runtime root containing the governed
runtime and its required libraries.

**Experiment evidence.** A narrow raw-only libkrun patch accepted an already finalized read-only
file descriptor and booted four owned unsandboxed guests with matching host and guest root digests.
It is a patch candidate, not admitted immutable custody. The final design still has to prove that
all writable aliases are gone, the sole path is unlinked, the digest is calculated through that
exact retained descriptor, and the installed sandboxed runner cannot substitute or reacquire the
root.

An earlier libkrun block-root path created an unexpected guest-visible `NullFs` virtiofs device.
**Experiment evidence.** A governed direct-block-root prototype booted without that device and
reran the bounded device corpus. This makes removal credible; it does not admit the patch or prove
the full virtual-machine device surface safe.

### Directional console transport and commit-last completion

**Selected design.** The first slice plans separate, bounded, one-directional virtio-console ports:
host-to-guest source, host-to-guest canonical input, and launcher-to-host completion containing
bounded JSON. Each message binds its role, version, attempt, plan/profile, length, and digest. The
host drains up to the cap plus one byte and rejects oversize data; end-of-file is never treated as
proof of a complete message.

The trusted guest launcher keeps the completion channel away from the workload. It verifies all
input before starting the workload, waits for the exact child process tree, and writes a fixed
commit trailer last. This **commit-last** rule distinguishes a complete result from a partial write
caused by a crash.

**Experiment and governed-source evidence.** Go and Node agreed on 43 backend-independent framing
vectors. The later public governed libkrun merge added bounded console/property and raw-FD library
tests, fixed two locally observed shutdown/lifecycle defects, and moved the four changed console
files from 13/88 to 37/88 covered functions and 90/728 to 298/733 covered lines. It still leaves
measured uncovered code, independent review, and the real transport, launcher, guest, App Sandbox,
teardown, and installed topology untested. The transport gate remains open.

## Deny by default and fail closed

**Selected design.** Unknown object versions, fields, powers, identities, algorithms, purposes,
audiences, states, and backend controls are rejected. Unsupported power is omitted from the v0
contract rather than represented as a switch that happens to be false. No “compatibility” fallback
may add a shell, writable path, network route, generic image, privileged helper, or alternate
runtime when the selected mechanism fails.

### Exact caps, or refusal

Trusted policy resolves missing resource values before approval. A request above the user's
ceiling is rejected, not silently reduced. The exact resolved values are shown in the plan. A
backend must enforce each value exactly or refuse the attempt; approximate CPU or memory controls
cannot be relabeled as exact ones.

**Experiment evidence.** The retained backend work supports only a narrow vocabulary, including
wall time, bounded console prefixes, integer virtual-CPU topology, closed guest-memory profiles,
bounded port frames, and fixed scratch-image bytes. CPU percentage/time, arbitrary memory, and
exact total host/virtual-machine memory remain unsupported.

### Record intent before effects

**Selected design.** For every external action, Capsule first commits a durable intent with a
stable effect ID, then performs the action, then commits the observed result. If the process dies
between the action and its result record, the intent remains as reconciliation work. A missing
response, process, path, or handle is not proof that the effect never happened.

**Implemented mechanics.** The fake lifecycle uses this intent/effect/reconciliation pattern for
prepare, create, start, observe, stop, and destroy. It retains cleanup obligations and stable fake
instance identities across controlled restart tests. It has exact ceilings of 256 active and 4,096
retained lifecycle records and never evicts explanatory state to make room.

### Cleanup, quarantine, and recovery

**Selected design.** Every path after backend creation must reach terminate, destroy, and
reconcile. Cleanup is complete only after authoritative absence for the exact backend identity.
Unknown state stays unresolved. Identity or cross-link mismatch enters quarantine and repair-
required state; the daemon cannot clear it. Output is not released as ordinary success while
integrity or teardown is indeterminate.

After a crash, startup takes one installation owner lock, validates the durable store, enumerates
attempts with missing or unfinished lifecycle work, and performs bounded reconciliation. Repeated
unknown results stop automatic retry rather than creating a restart storm. A consumed approval
remains consumed throughout; safe retry of the job requires a new human approval.

**Implemented local mechanics.** Archive Slice F1 adds passive archive projections, exact limits
and known answers, defensive copies, and a pure complete-cohort eligibility selector. F2 adds only
an owner-asserted v1-to-v2 all-hot migration and read-only empty-archive full verifier. Neither
moves a cohort, writes a segment, activates an archive, routes retained lookup, or calls an adapter.

**Owner G2 local mechanics.** Proposed ADR-0033 selects a pre-created enrolled sibling object
plus lifetime nonblocking BSD `flock`. G1 implements the internal Go/Darwin acquisition, and G2
composes it with the current owner-required v1/no-guest startup, same-session coordinator, sorted
recovery, post-open fence, and ordered close under owned temporary roots.

**Open production work.** The authenticated bootstrap and Apple-signed protected-state-root/
session/update/reboot matrix remain unimplemented. Archive F3+ segment/activation/retained-lookup
work, production archive/compaction, backup and rollback handling, real power-loss tests,
real-backend reconciliation, signed evidence, and installed recovery are unresolved.

## What exists today

### Capsule can claim today

- **Implemented mechanics:** exact-byte internal registration; candidate approval and immutable
  attempt records; atomic one-use consumption/create behavior; attempt-keyed durable fake lifecycle;
  explicit intents, effect IDs, reconciliation, capacity refusal, and controlled restart tests.
- **Implemented mechanics:** `FakeBackend.CreatesGuest()` is fixed to `false`. These tests create no
  virtual machine or hostile process.
- **Implemented mechanics:** archive F1 provides passive types, digests, defensive copies, and
  eligibility selection only; it writes and activates no archive.
- **Experiment evidence:** bounded macOS code identity, XPC, Keychain, user-presence, protected
  storage, per-user service, libkrun/HVF, runtime-construction, root-custody, device, framing, and
  packaging questions have retained results with stated environments and limitations.
- **Selected design:** the daemon/Broker/Supervisor authority split, registered exact plans,
  one-use attempt-bound approval, external isolation, exact-or-refused limits, controlled egress,
  and composed evidence are the intended requirements.

### Capsule cannot claim today

- It cannot run hostile agent-generated code safely as a product.
- It has no production Broker signing or approval interface, authenticated product IPC, protected
  production authority store, or complete content custody and release path.
- It has no admitted runtime, runtime profile, isolation backend, guest kernel, launcher, or final
  signed and notarized installed topology.
- It has no real directional guest transport, runtime-root custody composition, completed virtual-
  machine device corpus, or usable libkrun execution adapter.
- It has no production enforcement transcript or composed receipt chain.
- It cannot claim `validated-local`, continuous monitoring, platform attestation, production
  readiness, rollback-proof epochs, absence of data leakage through allowed output/timing, or proof
  that generated code is correct.

## Glossary

- **Ambient authority:** access inherited from the surrounding user or process environment rather
  than granted specifically for one task.
- **Approval Broker:** the trusted native role that displays registered typed plan data and signs a
  one-use grant after user presence.
- **Content address:** an identity derived from a hash of exact bytes rather than a mutable name or
  path.
- **Daemon:** the agent-facing planning service; it is not allowed to approve or launch work.
- **Digest:** a fixed-size hash result used here as an exact-byte fingerprint.
- **Effect ID:** the durable name for one intended backend operation, reused during reconciliation.
- **Execution Supervisor:** the sole intended authority for approval consumption, attempt creation,
  backend lifecycle, cleanup, and enforcement evidence.
- **Guest:** the disposable Linux environment that will contain the hostile workload.
- **HVF:** Hypervisor.framework, Apple's hardware-virtualization interface.
- **IPC:** interprocess communication, meaning structured messages between separate local
  processes.
- **Nonce:** a fresh one-use value retained to detect replay.
- **Runtime adapter:** the narrow layer that maps an admitted plan to one fixed runtime invocation.
- **Trust epoch:** a sequence-ordered identity for one coherent installed set of components, keys,
  policy, and trust state.
- **VMM:** virtual-machine monitor, the host process that creates and manages a virtual machine.
- **XPC:** Apple's local interprocess communication system.

## Deeper technical reading

- [Project definition](PROJECT.md) and [roadmap](ROADMAP.md)
- [Architecture](ARCHITECTURE.md), [technical design](TECHNICAL_DESIGN.md), and
  [Execution Supervisor](EXECUTION_SUPERVISOR.md)
- [Threat model](security/THREAT_MODEL.md) and
  [control/evidence matrix](security/CONTROL_EVIDENCE_MATRIX.md)
- [Workstream and evidence ledger](WORKSTREAM_EVIDENCE_LEDGER.md) and
  [Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md)
- [Related systems and design influences](RELATED_SYSTEMS.md), including which external patterns
  Capsule reuses, defers, or deliberately rejects
- [ADR-0010: authority separation](adr/0010-separate-planning-approval-content-and-execution.md),
  [ADR-0011: registered plans and one-use attempts](adr/0011-supervisor-registered-plans-and-one-use-attempts.md),
  and [ADR-0015: composed receipts](adr/0015-supervisor-transcripts-and-composed-receipts.md)
- [ADR-0018: platform-specific trusted components](adr/0018-platform-specific-trusted-components.md),
  [ADR-0022: libkrun/HVF evaluation](adr/0022-evaluate-libkrun-hvf-native-backend.md), and
  [ADR-0028: governed deno_core selection](adr/0028-select-governed-deno-core-first.md)
- Proposed details remain explicitly proposed in
  [ADR-0019: bounded CBOR/COSE](adr/0019-bounded-deterministic-cbor-and-cose.md),
  [ADR-0024: approval/attempt boundary](adr/0024-approval-consumption-and-attempt-creation.md),
  [ADR-0025: durable lifecycle state](adr/0025-colocate-durable-attempt-lifecycle-state.md), and
  [ADR-0026: pre-approval TypeScript emission](adr/0026-bind-pre-approval-typescript-erasure.md).
