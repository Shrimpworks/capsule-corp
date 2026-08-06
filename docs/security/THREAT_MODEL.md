# Threat Model

Status: intended repository-wide model. The current scaffold does not implement or satisfy these
properties.

## Overview

Capsule is a local-first trusted execution platform for bounded JS/TS jobs proposed by AI agents.
Its security objective is to contain hostile generated code, prevent an agent-facing compromise
from becoming user authorization or backend launch authority, limit resources and data flows, and
produce honest attributable evidence about the exact approved attempt.

The target architecture has three principal local authorities:

- the agent-facing daemon proposes and plans;
- the Trusted Host Broker owns user presence and user content;
- the Execution Supervisor alone creates hostile guests and signs enforcement transcripts.

The initial executable workload is one byte-exact dependency-free `main.mjs`, inline JSON in, and
bounded JSON out. Static/dynamic dependency requests, CommonJS, packages, and module-loader
fallbacks are outside the accepted first-release source contract. File
snapshots, broader output formats, stronger backend posture, and production updates follow only
after their evidence gates.

### Security objectives

1. A guest receives only exact authority approved for one immutable registered plan.
2. One approval creates at most one attempt and cannot be replayed after crash/restart.
3. A compromised agent or daemon cannot approve, directly launch, retrieve user-only content, reset
   trust state, or forge Supervisor terminal success.
4. The Supervisor independently rejects unsupported power and refuses backends unable to enforce
   exact required controls.
5. Inputs and outputs cross bounded content-addressed interfaces rather than ambient paths.
6. Runtime/component identity and trust transitions fail closed on mismatch or indeterminate state.
7. Receipts distinguish attributable claims from independent attestation and record limitations.

### Assets

- Host filesystem, credentials, environment, sockets, processes, and services not granted to a job
- User-selected input and produced output content
- Agent-client and host-application credentials
- Installation root, Approval, Supervisor evidence, and future transport/delegation private keys
- Installation manifest, component identities, trust epochs, policy, and quarantine state
- Pinned TUF roots, verified metadata checkpoints, local trust snapshots, and revocation state
- Registered plan bytes/digests, approvals, nonces, and grant-consumption ledger
- Original authoring bytes, emitted executable bytes, transformation objects, prepared source-set
  lifecycle, and the proposed Source Preparer's immutable-store integrity and availability
- Backend handles, cleanup leases, guest configuration, and management channels
- The exact registered `main.mjs` manifest/bytes and their pass-through custody from proposal
  through Broker rendering and runtime staging
- The integrity, availability, artifact identity, and fixed results of the proposed disposable
  `.mjs` Source Validator before planning and approval
- Other jobs' source, inputs, outputs, state, and cached data
- Runtime bundles, review attestations, registry activation, and backend validation records
- Artifact manifests, enforcement transcripts, receipt composition, and witness checkpoints
- Host availability and bounded CPU, memory, processes, storage, time, log, and output consumption
- Confidentiality of data not approved for the agent audience

## Threat Model, Trust Boundaries, and Assumptions

### Adversaries and attacker-controlled inputs

Capsule assumes the following may be malicious:

- AI-generated source, prompts, model output, labels, and source bundles;
- the agent/MCP/SDK client and authenticated agent requests;
- a compromised agent-facing daemon;
- a same-user local process attempting IPC, file, key, or component impersonation;
- user-selected file bytes and third-party data;
- the selected workload runtime, parsers, libraries, dependencies, and guest-controlled runtime
  behavior;
- all guest stdout, stderr, filenames, metadata, structured results, files, timings, and exit state;
- signed envelopes from unknown, revoked, expired, replaced, wrong-purpose, wrong-audience,
  wrong-installation, wrong-epoch, or compromised keys;
- DID strings, verification methods, headers, nonces, timestamps, signatures, and claims;
- OCI images, manifests, archives, registries, runtime/profile metadata, and backend reports before
  trust verification;
- external mirrors/services serving stale, equivocated, rollback, malformed, or malicious data;
- partial/corrupt local update, restored stores, stale components, and crash-replayed messages.

Authentication identifies a caller; it does not make its content trusted or authorize every
operation. A valid signature supports attribution to an enrolled key for a checked object; it does
not establish truth, human understanding, exclusive key control, or correct signer logic.

### Trusted assumptions

For validated-local posture, Capsule initially trusts:

- the intended user for approval and trust-changing ceremonies;
- the local operating-system administrator;
- the macOS kernel, code-signing/Keychain/LocalAuthentication services, Secure Enclave where used,
  virtualization stack, and hardware;
- the enrolled Broker and Supervisor logic and their protected keys;
- accepted trust roots, release processes, reviewers, and exact validation records within their
  delegated scopes;
- the exact isolation backend implementation/configuration after it passes required evidence.

If the proposed TypeScript profile is admitted, Capsule additionally trusts the enrolled Source
Preparer native front, Go core, immutable-store implementation, and—until a child boundary is
proven—Node/Amaro worker to choose emitted JavaScript and preparation claims faithfully before plan
construction. Supervisor and Broker validation proves internal byte/digest/media/path/record
consistency, not correct erasure or semantic equivalence. That planning and approval-understanding
trust consequence must be rendered to the approver; it is not independent attestation.

Under Accepted ADR-0035/0036, approval correctness additionally depends on the enrolled Oxc parser
artifact and the separately proven role-specific launcher/process profile classifying the exact copied
bytes faithfully. The validator is still attacker-exposed: parser input, diagnostics, crashes,
resource behavior, and returned facts are untrusted until the parent checks the fixed protocol,
artifact identity, digest, length, status, and deadline. Its compromise may cause false source
facts or denial of service, but must not expose an Approval key, mutate Supervisor state, or create
a guest. Independent runtime no-loader enforcement limits a false allow from becoming module-
loading authority.

The daemon is trusted for plan correctness in ordinary operation, but the design explicitly limits
what its compromise can authorize or disclose.

Defending against a compromised administrator, kernel, hypervisor, Secure Enclave, hardware, or
authorized Supervisor is outside local containment guarantees. An optional independent witness may
detect some historical inconsistency; it does not restore complete local protection.

Capsule does not currently claim resistance to cross-domain microarchitectural leakage such as
speculative-execution or shared-cache attacks between a hostile guest, the host, or another job.
Any stronger posture must bind an exact hardware model, OS/hypervisor build, applicable vendor
mitigation state, and concurrency/co-residency policy to retained evidence. A hardware model name,
Hypervisor.framework use, or absence of simultaneous multithreading alone is not evidence of
noninterference.

### Same-user attacker tiers

`Same-user` is not a complete claim without naming the authority already held by the attacker.
Capsule separates three tiers:

| Tier | Included capability | Claim treatment |
| --- | --- | --- |
| Baseline same-user attacker | Arbitrary unprivileged process under the same login; ordinary same-UID pathname/directory access; races, links, replacement, and observation; retained writable descriptors/mappings to original user files; malformed IPC; component impersonation and attach attempts | The v0 boundary must resist it. A new custody object is not custodied while this attacker can acquire or retain writable authority to it. |
| Elevated user-granted attacker | Successful task-port/debug attachment, Full Disk Access, Accessibility control, Screen Recording/window automation, explicit foreign-container consent, private VFS privilege, or another broad user-authorized capability | Separate elevated-adversary posture. Shipping denial/fail-closed cases remain required, but Capsule makes no general resistance claim without exact evidence for that capability. |
| Trusted-platform compromise | Malicious root/administrator, SIP bypass, kernel/hypervisor compromise, or an authorized compromised Supervisor | Outside local containment guarantees. |

The baseline attacker may retain writable authority to an original selected file. Broker snapshotting
must create a new object for which that authority cannot be acquired. Shipping enrolled components
use hardened runtime without `com.apple.security.get-task-allow`; attempted attachment belongs in
the exact installed corpus, while a successful task port moves the case to the elevated tier.

### Agent to daemon

The agent-facing transport exposes proposal, fixed status, authorized cancellation, and fixed
summary operations. Agent credentials cannot call approval, file selection, content redemption,
trust administration, quarantine reset, or backend control.

Raw input is bounded before ordinary decoding. Unknown versions, fields, powers, and malformed
content fail closed.

Under Accepted ADR-0035/0036, the daemon sends one exact copied `main.mjs` frame through its
daemon-private launcher to a fresh one-shot Source Validator child before plan construction. The
child has no path/package/loader API, persistent Capsule product store, keys,
network, backend route, or inherited ambient descriptors. A parser/semantic diagnostic, crash,
timeout, resource ceiling, malformed/extra output, wrong artifact/version, or digest/length
mismatch refuses. There is no lexical-scanner or second-parser fallback.

Per-request parser limits do not establish daemon availability. Before activating the candidate
public endpoint, the daemon must also enforce a configured aggregate envelope for connections,
concurrent requests, in-flight bytes, queues, request time, cancellation, and stalled downstream
work. Overload sheds work with bounded errors and no authority change; it never creates an
unbounded queue or silently widens an approved job budget.

### Daemon to Supervisor

The daemon sends exact canonical plan bytes through code-identity-authenticated local IPC. The
Supervisor independently repeats raw/schema validation and applies non-overridable v0 hard-safety
rules before durable registration.

Attempt APIs accept only a registration ID. No API accepts replacement plan bytes, backend flags,
images, mounts, guest paths, or policy overrides at execute time.

### Broker to Supervisor

The Broker fetches registered bytes directly, independently validates/hashes/renders them, and signs
an attempt-bound approval after user presence. It transfers only attempt-scoped content handles.

Trusted UI must safely render untrusted labels, sizes, and source/content metadata. It never claims
that the user understood generated source merely because presence was proven.

For the owner-only internal alpha, Accepted ADR-0040 does not require the blocked Source Validator.
The Broker trusts neither daemon prose nor child-supplied prose: it renders only strictly bound
Supervisor-retained plan, role-binding, registration, manifest, and exact-source bytes. Accepted
ADR-0043 freezes the bounded ASCII-only source projection and explicitly discloses that the current
readback contains inline-input digest and length but not content. A later fresh Source Validator may
add typed grammar facts as defense in depth; it is not approval authority.

UI activation, focus, or a synthetic click is not user presence. Approval requires the configured
LocalAuthentication/Keychain-gated private-key operation over the exact registered-plan binding.
Baseline synthetic-event attempts must fail to produce a grant without that operation. Accessibility,
overlay/window automation, and similar broad user-granted capabilities remain an explicit elevated
posture with separately recorded limitations and tests; Capsule does not imply they are defeated by
ordinary anti-spoofing UI.

The retained ADR-0043 implementation is passive evidence only: no installed UI, `LAContext`,
Keychain item, private key, prompt, signing operation, IPC endpoint, or activation path exists. Its
public-key-only verifier proves strict Sign1/key/time/role binding for checked-in vectors, not human
presence, comprehension, correct UI logic, or synthetic-input containment.

### Local process and storage boundary

The accepted Source Validator architecture adds hostile-input local launcher/child boundaries.
Until fixed role-specific artifact enrollment, private-XPC reachability, OS sandboxing, descriptor
closure, reactive-resource evidence, kill/drain/reap recovery, cleanup, and compromised-child cases
are retained on supported hosts, the product control is only a design plus historical crash-
isolation evidence. The child has zero authority effect; every abnormal outcome is a refusal. A
same-UID process name, PID, or executable path is not identity evidence.

The retained V2 macOS checkpoint now supplies negative evidence for this exact artifact/profile
pair. The strict child refuses before `exec` because `RLIMIT_AS` cannot be lowered. When that limit
is explicitly omitted only for diagnostic testing, the child process can open an owned out-of-cwd
file, create IPv4 and Unix socket descriptors, mutate cwd metadata, create an empty file, and map
512 MiB. The supported App Sandbox child entitlement shape changes the fixed V1 bytes; deprecated
custom sandbox interfaces were not used. CPU/file/FD/child/wall ceilings, output faults, group
kill/reap, and a clean later invocation are local-mechanic evidence only. V2 remains `BLOCKED` and
no Keychain, Supervisor-state, installed-identity, sandbox, or product-containment claim follows.

The later [supported-profile review](../MJS_SOURCE_VALIDATOR_MACOS_PROFILE_REPLACEMENT.md) rejects
direct App Sandbox inheritance because the child receives daemon/Broker static rights. Accepted
ADR-0036 selects two distinct role-private App-Sandboxed XPC launchers, each owning one fresh
matching parser child. No service, result, cache, container, group, key, profile, or accepted request
crosses daemon/Broker roles; neither launcher reaches Supervisor state or backend authority.

App Sandbox gives each launcher/child read/write access to its private container. That is accepted
only as residual scratch authority: persistent Capsule product state, cache, source/diagnostic log,
and reusable result are forbidden, and cleanup/residue testing is required after every request,
crash, restart, update, and startup. Cleanup reduces persistence and supplies evidence but cannot
prove confidentiality or secure erasure against a compromised parser, swap, backup, or forensic
recovery.

The observed host exposes no usable public unprivileged hard memory cap; the public-SDK footprint
setter returned `KERN_NO_ACCESS`. ADR-0036 accepts an evidence-derived reactive physical-footprint
watermark with one child per launcher request, bounded per-role and combined two-role concurrency,
fixed sampling cadence, and process-group kill/drain/reap. Sampling can overshoot before kill, so
the claim is neither a hard peak/exact cap nor a host-availability guarantee. Threshold, cadence,
baseline, maximum observed overshoot, and kill latency remain unset until the separately authorized
signed corpus. Product use remains blocked until that evidence and both independent consumers pass.

The replacement threat set includes cross-role request/result replay; a generic or shared service;
old/new consumer-service-parser-profile mixes; parent/responsible/self or library-constraint
substitution; native/JIT/DYLD/debug loading; inherited file descriptors, Mach/bootstrap rights,
network/IPC, out-of-container filesystem and authority-store access; container residue across calls
or updates; descendants and parser survival after launcher death; delayed/missed footprint samples;
two-role concurrent overshoot; and system-pressure collateral damage. Any unsupported private-XPC
reachability, authority/native-loading/filesystem/network escape, unbounded orphan/cleanup result,
accepted mixed version, or measured host risk beyond the reviewed profile stops the exact candidate
rather than enabling a shared bus, private API, custom sandbox, or privileged helper.

Separate authorities require proven macOS enforcement: XPC peer code requirements, effective user/
session and exact build/epoch checks, Keychain access-group separation, protected storage
containers, relevant entitlements, and Supervisor-only backend control.

PID, path, process name, same-user mode bits, or a diagram alone do not establish identity or
containment. Broad shared app groups are disallowed.

Proposed ADR-0033 uses one enrolled pre-created regular sibling object and lifetime nonblocking
BSD `flock` to serialize cooperating Supervisor processes before store access. The advisory lock
is not same-UID containment: the retained local harness observed that an actor able to mutate the
parent directory can rename the locked inode and create a separately lockable replacement.
Product use therefore requires the installed Supervisor-private protected state root, exact
device/inode enrollment, and fail-closed repair on mismatch. No cross-user or elevated local-
capability claim follows from the lock.

Proposed ADR-0038 selects the bootstrap authority without advancing installed evidence. A
separately enrolled, on-demand Trust Coordinator uses a Coordinator-only user-presence-gated
installation-root key to sign one closed request and the final observed record. The Supervisor
alone creates the fixed root/owner/disabled-store genesis in its private container. The visible
app, daemon, updater, and replacer receive neither key nor state authority. A dedicated
Coordinator/Supervisor App Group supplies only the setup Mach name; Capsule creates no shared-
container entry or shared Keychain item and treats both namespaces as residual capability. Apple
peer/code identity authenticates the handoff but never substitutes for the Capsule installation-
root signature.

The I2B3 exact Team-3DDR profile/signing preflight later passed, but the mandatory stale-profile
test found an unresolved boundary: a Supervisor signed with the still-valid archived I1B profile
could rewrite a current-profile sentinel because both stable signing identifiers resolve the same
App Sandbox private container. New CDHash, profile, App Group, and role-private Keychain-group
projections did not version that filesystem authority. The run stopped and cleaned the exact
sentinel before any key, service, installation, or root existed. Proposed ADR-0045 now selects a
versioned signing/container authority epoch, but installed I2 remains `BLOCKED` on its exact
Apple Development evidence; peer requirements alone do not close separately launched stale-code
write authority. See the
[I2B3 blocker result](../MACOS_INSTALLATION_I2B3_SIGNING_PREFLIGHT_AND_STALE_PROFILE_BLOCKER.md).

Proposed ADR-0045 selects a fresh explicit Supervisor and Coordinator application identity,
private-container association, LaunchAgent label, bootstrap App Group/Mach service, role-specific
Keychain groups, descriptor digest, and state schema/engine binding for each independent
`SupervisorAuthorityEpoch`. The stable Supervisor identity is legacy residue and is never epoch
zero. OS nonmembership must deny stale access to the new private container and groups without a
consent override; exact peer requirements and signed epoch objects only detect stale protocol
traffic and cannot substitute for that denial. A foreign-container prompt or granted consent moves
the case to the elevated user-granted tier and keeps attempts disabled. The passive decision and
inert mutation matrix supply no installed control evidence; I2B remains `BLOCKED`.

The proposed Source Preparer's sole role-namespaced store is an ownership topology, not a security
boundary by itself. TypeScript P1 requires an exact OS-enforced protected container with the Source
Preparer as its single enrolled member and negative open/link/replace/map/rename/handle-retention
evidence for daemon, Broker, Supervisor, updater, stale component, and unrelated same-user process.
Its one-shot Node worker also requires an exact child profile and process-tree cleanup proof; fixed
argv/environment/cwd/descriptors alone do not prevent inherited same-user or parent authority.

### Supervisor to backend

Backend configuration is generated from trusted typed data and contains no guest-controlled shell
interpolation, path, image, flag, mount, socket, seccomp, or management channel.

The Supervisor matches every required control against an exact capability report and accepted
validation record. A backend that cannot enforce a required value refuses the attempt.

### Guest to host

The guest is hostile. External isolation controls syscalls, filesystem, processes, IPC, network,
resources, and lifecycle. Runtime permission systems are supplemental only. Capsule never runs an
untrusted workload runtime directly on the host.

Every admitted runtime bundle binds the exact guest-kernel image, configuration, boot arguments,
module policy, debug/inspection facilities, and launcher restrictions. Removing unused kernel and
device authority is defense-in-depth that reduces reachable surface and must be reviewed and
tested, but the VMM boundary must still contain a fully hostile guest kernel.

### Guest completion evidence scope

For host containment, the entire guest—including its kernel—is hostile and must remain inside the
VMM boundary. For ordinary completion semantics, the exact admitted guest kernel and trusted
launcher are part of the runtime TCB. They may report one syntactically valid attempt-bound
completion envelope; Capsule does not claim that report survives guest-kernel compromise.

Attempt/profile binding, length/digest checks, and commit-trailer framing reject stale, torn, or
ordinary user-process-forged records. They are not attestation: a malicious guest kernel can observe
or misuse guest-held authority. The Supervisor transcript therefore records host-observed profile,
device/port topology, limits, runner lifecycle, envelope validity, bounded-result disposition, and
teardown, with the explicit limitation that workload completion is guest-reported. It does not
prove correct execution, user intent, or an uncompromised guest kernel.

The launcher is a distinct admitted guest process, not the current experiment's `exec`-and-replace
shim. It fully verifies source/input before starting the admitted runtime, withholds the completion
endpoint and node,
uses a fixed child FD/argv/environment/cwd manifest, bounds the child result, waits for the exact
child tree, and writes the commit trailer last. The host runner separately starts with an exact
role-specific FD allowlist because a VMM compromise acquires any ambient descriptor it inherits.

The passive C2B v3 contract makes that boundary exact without claiming execution: host FDs 0–7,
launcher FDs 0–5, and runtime-child FDs 0–2 are closed; only balloon, RNG, one three-port console,
and one read-only raw-root block device are allowed; implicit console/init/vsock are explicitly
disabled; TSI, network, virtiofs/`NullFs`, extra block devices, sockets, and live host paths are
absent. Libkrunfw supplies the sole non-EFI boot kernel, so external kernel/firmware path calls are
forbidden. The contract does not authorize a guest and blocks stale libkrun or missing runner bytes.

The immutable C2B v4 build/static successor closes those two artifact blockers without weakening
the boundary. Its retained unsigned runner imports only the reviewed fixed calls, validates exact
FDs 0–7, closes from 8, disables implicit console/init/vsock, creates the three ports in fixed
order, and accepts no path, plan, profile, image, mount, backend flag, or replacement configuration.
Exact source confirms libkrunfw is the sole non-EFI boot-kernel carrier; no separate firmware
identity is invented. The v4 runner and libkrun were never loaded or executed in that immutable
scope. A later separately authorized
[v19 fixed-owned-guest experiment](../FIRST_OWNED_GUEST_EXECUTION_CHECKPOINT.md) loaded its exact
experimental successor bytes and passed one fixed benign boot/runtime/console-proof/reap/teardown
attempt. It did not use the final typed transport or hostile source. The next exact v20-v23 runners
refused before readiness; the Supervisor issued no start authorization and no guest launched. V23
retained authoritative staged-path/open-FD/runner hash evidence and proved the actual root was
correct while the embedded expected digest bytes were malformed. V24 replaces them with an exact
literal plus static length assertion, then passes one exact known answer and the early non-root/
capability/descriptor/root/host-path/mount/root-regain denial controls in a real guest. It stops in
the vsock-check family without connecting or sending. Semantic review rejects v25's socket-open
property as non-authoritative. V26 instead tests local-CID ioctl capability, still without any
connect/send path, and retains fixed failure detail plus bounded console progress. Its exact
invocation passed active local-CID/raw-block denials and then reported expected down/unbacked
`dummy0`. V27 accepts only loopback plus that bounded passive inventory and rejects virtio backing
and non-loopback/default routes without connecting or sending. Its exact authorized guest completed
that control plus the entire fixed 30-marker denial corpus, exact completion/console proofs, normal
reap, unlinked-root teardown, unchanged canary, and zero network/credential authority or traffic.
This one local reproduction does not cover every kernel/hypervisor path. Broader hostile-guest
transport, pressure/lifecycle/platform validation, installed confinement/recovery, and runtime/
profile admission therefore remain `BLOCKED`.

The virtio-console implementation is part of the hostile-guest-to-VMM attack surface. Application
framing does not validate guest-controlled control IDs/events, queues, descriptor chains, reset/
open/close ordering, or cancellation/backpressure behavior. The exact pinned implementation and any
fixes require their own sanitizer/coverage corpus and fail-closed teardown evidence.

### Content and egress

The Broker owns persisted user content. The Supervisor receives only transient attempt-scoped
handles needed for staging and collection. The daemon receives no user-only content by default.

The Supervisor first applies filesystem safety; the Broker then applies bounded content validation
and user delivery. The agent receives a smaller fixed summary, not a redacted full receipt.

### External trust and updates

Pinned TUF root metadata anchors distribution. A network-capable updater verifies full external
metadata and produces a compact local trust snapshot. Live execution performs no trust-service
call, arbitrary DID resolution, or general TUF/network parsing.

Component-changing updates use a prepared, authorized, crash-safe trust transition. Partial state
enters `repair-required` rather than accepting whichever components start.

### Mandatory security properties

#### Identity and key authority

- Normative local identity is installation ID plus locally authorized public keys.
- Installation-root, Approval, and Supervisor evidence keys are purpose-separated. The
  installation-root key is available only to the on-demand Trust Coordinator ceremony; the
  visible app, daemon, Supervisor, updater, and replacer cannot use it.
- Approval requires fresh user presence for every v0 plan.
- Unknown, revoked, suspended, replaced, expired, wrong-purpose, wrong-type, wrong-audience,
  wrong-installation, and wrong-epoch keys/objects fail closed.
- DIDs identify principals or verification methods but never grant authority.
- Local v0 authorization uses no live DID resolution, arbitrary method, plugin, or remote context.

#### Plan and approval integrity

- Hashes are byte identities, not origin proof without a trusted binding.
- Broker approval and Supervisor registration bind the same canonical plan digest/bytes.
- For the first release, registration atomically retains one canonical single-member source
  manifest and the exact pass-through `main.mjs` bytes; Broker and execution rehash those retained
  bytes. There is no transform, second source role, dependency request, or execute-time byte input.
- When TypeScript is supported, the registered plan separately binds exact original authoring
  bytes, exact emitted executable JavaScript bytes, and the closed transformation identity before
  approval; transformation from only an original digest after approval is forbidden.
- Proposed ADR-0032 assigns those pre-registration bytes to one enrolled Source Preparer store.
  The daemon cannot mint a stored source set; the Supervisor and Broker independently validate all
  digest/media/path/record relationships; neither validation proves the selected toolchain ran,
  that erasure was faithful, or semantic equivalence. Approval must distinguish every original and
  executable role, show the fixed toolchain/options/record identity and dispositions, permit
  bounded inspection of both byte roles, and disclose that executable selection is trusted rather
  than independently reproduced.
- Approval binds registration, installation, epoch, expected Supervisor, attempt nonce, purpose,
  audience, and expiry.
- Approval is atomically consumed with attempt creation before backend side effects.
- A changed plan, content, profile, policy, limit, audience, backend requirement, or epoch requires a
  new plan/registration/approval.
- The Supervisor independently enforces v0 hard-safety rules.

#### Isolation and resources

- No arbitrary host path, live user-file mount, credential, environment, descriptor, or host socket
  reaches the guest.
- Root is immutable under host custody; first-slice source/input and completion/result ports,
  scratch, and later artifact output are separate and bounded.
- Network denial includes TCP, UDP, DNS, IPv4/IPv6, loopback expectations, metadata services, Unix
  sockets, vsock/management channels, and inherited host IPC.
- Process isolation prevents host/other-job signal, inspect, attach, and state reuse.
- Required wall/CPU semantics, memory, process/PID, storage, log, output, and cancellation controls
  are externally enforced exactly or rejected.
- Every post-create path reaches terminate/destroy/reconcile; teardown failure is distinct.

#### Content and observations

- Agents cannot mint authority from paths, URLs, environment names, DIDs, or arbitrary identifiers.
- File capability is exact immutable regular-file data-fork bytes, never the original path.
- Only declared fixed output slots, regular files, and bounded byte counts reach content parsing.
- Rich parsing occurs in a future disposable parser sandbox, not the daemon or Supervisor.
- Agent summary contains no guest-controlled strings, artifact names/sizes, paths, timings, rich
  violations, or full metrics by default.
- Metadata minimization reduces but cannot eliminate state/timing/covert-channel leakage.

#### Runtime and installation integrity

- Trusted IPC checks code identity, effective user/session, exact enrolled build, relevant
  entitlements, and common epoch using supported OS mechanisms.
- Supervisor state additionally binds a distinct versioned `SupervisorAuthorityEpoch`; release and
  installation-trust epochs reference its digest but do not replace its application/container/
  group/key/state-engine identity.
- Debugged, dynamically invalid, mismatched, stale, or partially updated components cannot claim
  validated posture. Stable App Sandbox container identity is not treated as a stale-binary write
  fence; I2B3 proved an archived Supervisor profile could reach current-profile container state.
- The daemon cannot clear degraded/quarantined/repair-required/compromised state.
- Trust epochs are sequence-ordered and never described as rollback-proof without an additional
  anchor/witness.
- An integrity failure after grant consumption burns the approval, quarantines output as needed,
  and cannot become success-with-warning.

#### External trust

- Pinned TUF roots, not URLs or DIDs, anchor release/profile trust.
- Versions, delegated scope, expiration, hashes, snapshot consistency, and rollback checkpoints are
  verified before producing a local trust snapshot.
- TUF carries Capsule-defined revocation/disable records; Capsule policy defines their semantics.
- Offline/unavailable service never causes acceptance of unsigned or rollback state.

#### Persistence and recovery

- One enrolled installation owner lock is acquired before store read/mutation, recovery, archive,
  backup, repair, or adapter work; duplicate ownership refuses without creating a lock/store or
  issuing an owner-session permit.
- Initial root/owner/store creation occurs only under an exact installation-root-signed request
  from the enrolled Coordinator. Ordinary startup reads the fixed signed anchor and opens without
  creation; missing or mismatched protected state is repair-required.
- Grant consumption, attempt identity, trust/quarantine state, backend handles/cleanup leases, and
  content release state are durable before dependent side effects.
- Cross-store messages are authenticated, bounded, idempotent, and attempt/epoch/content-bound.
- Missing backend state is not assumed to prove destruction.
- Collected content is not released while integrity/teardown state is indeterminate.
- Repair preserves or explicitly replaces trust, grant, attempt, cleanup, and evidence history.
- Archive/compaction never treats expiry or terminal state as permission to forget replay,
  identifier, nonce, effect, instance, cleanup, or explanatory history. A cohort may leave hot
  state only when all of its attempts are durably destroyed after authoritative absence.
- An indeterminate archive publication/activation outcome fences until reopen establishes the old
  complete hot world or the new complete referenced-archive world. Missing referenced archive data
  is repair-required, not permission to resurrect hot authority.

#### Evidence

- A user receipt composes Broker approval and Supervisor enforcement evidence bound to the same
  plan/registration/attempt/installation/epoch.
- The daemon cannot forge either embedded authority claim.
- Posture retains isolation, runtime-integrity, trust-freshness, and distribution dimensions.
- Receipts state observation/claim limitations and do not imply independent attestation.

## Attack Surface, Mitigations, and Attacker Stories

### Public protocol and canonicalization

Risks include parser differentials, duplicate keys, Unicode/numeric disagreement, oversized input,
type confusion, algorithm confusion, and cross-object signature reuse. Mitigations are strict
pre-schema decoding, canonical-on-wire rules, protected-header allowlists, purpose/type separation,
safe numeric constraints, and shared Go/Swift/TypeScript negative fixtures.

Attacker story: the daemon and Broker interpret two spellings of a plan differently. Capsule must
reject non-canonical or duplicate-key input before any approval/state transition.

### Daemon service exhaustion

Risks include concurrent maximum-size requests, connection churn, slow or partial bodies, parser
CPU amplification, cancellation races, stalled Supervisor calls, queue growth, and high-cardinality
diagnostics. Per-object byte/depth/count limits are necessary but insufficient. The daemon requires
bounded per-peer and global in-flight accounting, concurrency and queue ceilings, deadlines,
cancellation propagation, bounded diagnostics, and deterministic overload behavior.

### Daemon compromise

The daemon can propose a malicious plan or lie in unsigned metadata. It must still be unable to use
Approval/Supervisor keys, access Broker content, bypass Supervisor hard safety, replace registered
bytes, launch a backend, reset grant/quarantine state, or forge terminal enforcement evidence.

Attacker story: a malicious MCP payload obtains code execution in the daemon and tries to run a
networked guest. Independent Supervisor schema/hard-safety validation and Supervisor-only backend
control must deny it without relying on the compromised daemon.

### Same-user impersonation and component substitution

Risks include connecting under an expected process name, copied/stale binaries, same-team wrong
components, debugger attachment, shared Keychain groups, broad app groups, restored stores, and
partial update. Mitigations depend on proven XPC code requirements, Keychain/access-container
separation, exact build/epoch checks, dynamic validation, and repair-required transitions.

### Approval/UI abuse

Risks include daemon-supplied prose, hidden fields, truncated/bidi labels, deceptive provenance,
approval fatigue, stale sessions, wrong Supervisor, replay, synthetic input, Accessibility-driven
automation, overlay/window spoofing, and focus substitution. The Broker renders Supervisor-
registered typed data with bounded safe UI and a LocalAuthentication/Keychain-gated signing
operation; the Supervisor consumes the grant once. A click or focused window alone is never
approval evidence.

Passive ADR-0043 conformance neutralizes raw bidi/control bytes, rejects mixed or stale Supervisor
readback, freezes context invalidation requirements, and verifies public-key ApprovalGrant vectors.
It deliberately remains approval-ineligible and does not advance the installed UI control.

The Broker can produce a signature whose local key operation was gated by configured user-presence
controls. A remote verifier can check key attribution, but cannot independently verify user
presence, comprehension, or correct UI logic.

### Guest/backend attack surface

Risks include filesystem escape, process breakout, network/metadata access, management-channel use,
runtime bugs, fork/output/disk/memory exhaustion, cancellation escape, writable cache reuse, and
orphaned guests. External backend mechanisms, fresh state, exact capability matching, durable
handles, and a retained attack corpus are mandatory.

Apple Containerization remains development-only after Gate C found no supported durable host-side
VM/helper identity or restart enumeration; ambiguous controller loss cannot produce ordinary
success or capability release. The native libkrun/HVF candidate makes one signed VMM process the
VM lifecycle object and gates start on a durable PID/start/code-identity record. That reduces the
hidden-helper risk but does not trust PID/path alone or prove safety against a malicious guest,
VMM exploit, output flood, or untested disk/recovery path. Its readiness corpus additionally found
that a same-user process could mutate a live raw backing image and that the block-root path exposed
a `NullFs` virtiofs device without a host-backed share. Those are unresolved custody and VMM-surface
risks, not evidence of an observed escape. The reconciled first slice rejected stock Bun's
subprocess and FFI surfaces as an authority mismatch, selected governed `deno_core` as the first
engineering candidate without admitting it, and proposes bounded dedicated
virtio-console ports for source/input and typed inline completion. The pinned guest kernel and
trusted launcher are part of that profile's TCB; an in-guest completion record is not attestation
against a compromised guest kernel. gVisor validation must bind its
host/outer-VM, engine, cgroup/OCI configuration, Sentry/gofer identity, management endpoints, and
exact `runsc` binary. The existing runc control run validates only the surrounding harness, not
gVisor's isolation boundary.

The admitted guest kernel is an exact runtime-profile component with its own configuration and
provenance. A minimal configuration and disabled unused facilities reduce the workload-to-kernel
and kernel-to-device attack surface, but cannot justify relaxing malicious-kernel VMM tests or
treating guest completion as attestation.

### Input and output content

Risks include path traversal, symlink/hard-link/special-file races, mutation, archive/parser bugs,
sparse files, terminal escapes, CSV formula injection, bidi/HTML spoofing, output flood, and content
release to the agent. v0 snapshots exact regular-file bytes, assigns paths, separates filesystem and
content gates, bounds parsers, and defers rich formats to another sandbox. The first inline JSON
slice avoids a guest filesystem-image parser by returning one bounded typed result frame; that does
not waive parser isolation before file artifacts are supported.

Generated code may intentionally encode complete granted input into allowed output, size, timing,
state, or repeated calls. Capsule warns users and budgets channels but does not claim
noninterference.

### Trust service and update attacks

Risks include malicious mirror/service, freeze, rollback, mix-and-match, delegated-role escape,
expired metadata, compromised publisher/reviewer, partial install, coherent local rollback, and
repair history reset. Pinned TUF roles, independent review/activation/validation, local snapshots,
prepared epochs, and fault-injected recovery reduce these risks.

Archive introduces a related local rollback and omission surface. An attacker or fault may remove
a segment, substitute an older otherwise valid segment/index/checkpoint, split a registration
cohort, or restore a coherent pre-archive world so that a consumed approval or identifier appears
absent. Proposed ADR-0031 requires immutable full-record segments, reconstructed exact tombstone
indexes, a hash-linked visible checkpoint, publish-before-activate ordering, and fail-closed reopen.
Those mechanisms detect corruption and mix-and-match only inside the visible set. They do not
prevent coherent rollback without an independently protected latest checkpoint.

An external witness is optional and privacy-sensitive. It may improve historical rollback evidence
but does not authorize local execution.

### Component compromise summary

The detailed blast-radius and response matrix is
[Component Compromise Matrix](COMPONENT_COMPROMISE_MATRIX.md). The highest-impact local compromise
is the Supervisor because it owns launch authority and enforcement claims. Supervisor compromise is
not concealed by the presence of a valid user approval.

A compromised proposed Source Preparer can choose mutually consistent executable JavaScript and
preparation metadata before approval, retain or disclose agent-supplied source, exploit any omitted
cross-role/version/install/epoch/source/path/media/digest/options/profile/record binding, deny
preparation, and corrupt its own unregistered store state. It cannot thereby create Approval or
Supervisor-evidence signatures, register or launch a plan, read Broker-only content, clear
quarantine, or replace already registered bytes without a digest mismatch. Exact blast radius and
P1 stop conditions are retained in the
[P0 authority/TCB checkpoint](../TYPESCRIPT_SOURCE_PREPARER_P0_AUTHORITY_REVIEW.md).

### Abuse cases and required tests

| Category | Cases |
| --- | --- |
| Protocol | raw byte/depth/count limits, invalid UTF-8, duplicate keys, unsafe numbers, trailing data, unknown fields/versions/powers |
| Daemon service | concurrent maximum-size requests, connection churn, slow/partial bodies, parser amplification, queue saturation, stalled downstream calls, cancellation races, bounded diagnostics, deterministic overload |
| Cryptography | non-canonical bytes, wrong type/purpose/audience/epoch, `none`/algorithm confusion, key substitution, malformed DER/raw/high-S handling |
| Identity/IPC | unsigned/same-team-wrong-ID/stale/debugged peer, wrong user/session, Keychain/store access, PID/path/name substitution |
| Plan/approval | plan mutation, registration swap, replay, wrong Supervisor, stale nonce, expiry, concurrent/double attempt, crash-after-consume, synthetic input, Accessibility automation, overlay/focus substitution, signing without configured user presence |
| Source validation | cross-role/cross-version result, shared launcher/cache, wrong private service/profile/epoch, direct inherited helper, container residue, out-of-container/store/key/network access, native/JIT/DYLD/debug loading, inherited FD/Mach/bootstrap rights, descendant/orphan, launcher/child/consumer death, missed footprint sample, combined-role overshoot, cleanup/update/startup refusal |
| Source preparation | compromised front/Go/Node/store, fabricated coherent emission, original/executable/profile/options/record substitution, cross-role/version/install/epoch/path/media replay, protected-store same-user access, stale-process handle retention, worker store/package/network/process/native-loading access, cap-plus-one, cancellation/death/orphan/stdio flood, genesis/update/rollback/release without authority |
| Filesystem/content | traversal, live path, mutation, symlink/hard link, devices/FIFO/socket, sparse file, archive, parser/formula/terminal/bidi/HTML hazards |
| Network/IPC | TCP, UDP, DNS, IPv4/6, loopback, metadata, Unix/vsock/management sockets, inherited descriptors |
| Process/runtime | workers/fork bomb, signals, inspector, orphan processes, subprocess APIs, native addons, FFI, macros, auto-install, `.env`, dynamic import abuse, completion-descriptor forgery |
| Guest kernel | image/config/boot substitution, module/debug/inspection authority, unnecessary drivers/protocols, launcher escape, malformed device traffic after kernel compromise |
| Resources | busy loop, CPU semantics, heap/native OOM, PID exhaustion, disk/log/output flood, cancellation tree kill |
| Isolation | host credentials/env/files, cross-job state, writable caches, backend escape, malicious runtime/image/config, mismatched host mitigation state, forbidden guest co-residency, speculative/shared-cache probe limitations |
| Egress | undeclared/oversized/malformed output, agent content access, metadata/timing channels, guest strings in errors/logs |
| Trust/update | TUF freeze/rollback/mix-and-match/delegation, revoked profile, partial component update, restored state, coherent rollback |
| Recovery | crash at every saga/side-effect boundary, orphan enumeration, missing handle, teardown failure, repair history preservation |
| Evidence | daemon-forged success, swapped/removed grant or transcript, cross-attempt receipt, overstated posture/attestation |

### Existing and planned mitigations

The current repository implements schema/tooling plus named unwired local mechanics for exact
registration, atomic approval consumption/attempt creation, and a durable no-guest fake lifecycle.
The [public-key approval integration](../APPROVAL_FAKE_LIFECYCLE_INTEGRATION.md) now connects those
mechanics for one public-only signed exact-plan fixture: plan A cannot authorize plan B, equivalent
signatures share canonical-payload replay identity, the attempt commits before fake effects, and
response-loss/restart recovery does not redrive an effect. This adds no live signer, authenticated
IPC, installed authorization, real adapter/backend, runtime, guest, or product control.
The E5 fixed-store corpus proves exact lifecycle ceilings, no eviction, destroyed-only active
capacity release, and bounded repeated startup. G2 now composes that no-guest path with the local
Darwin owner-required v1 opener and same-session coordinator. It implements no runtime security
control, installed protected-root boundary, consumer, real backend, or guest. The authoritative
claim registry is [Control Evidence Matrix](CONTROL_EVIDENCE_MATRIX.md); rows advance only for the
exact retained mechanism and evidence they name.

The [passive durable completion-last contract](../DURABLE_COMPLETION_CONTRACT.md) now retains one
additional local oracle: only a transaction over committed attempt/lifecycle truth can publish the
FakeBackend completion, bounded transcript, and fixed summary. EOF, runner exit zero, forged or
mixed records, changed replay, response loss, and unresolved teardown cannot manufacture or replace
that object. FakeBackend has no runner or process tree, so the retained record says runner identity
is unresolved and limits absence to the simulator; it supplies no real teardown or product-success
evidence.

Passive owner-lock G1 implements the internal Darwin enrolled-object acquisition and opaque
lifetime capability. G2 composes it with only the current v1/no-guest startup under owned temporary
roots, including duplicate-before-store refusal, sorted recovery, response-loss/process-death
reopen, post-open ownership fencing, and store-before-owner shutdown. It remains an unwired local
mechanic and supplies no protected-storage, same-UID-containment, signed-bootstrap,
installed-identity, rollback, archive, runtime, real-backend, or guest claim.

Proposed ADR-0031 and its conformance plan define the full-record archive and exact replay/non-reuse
tombstone boundary. Passive F1 types, exact limits/known answers, defensive copies, and complete-
cohort eligibility now exist. The passive
[F2 format blocker resolution](../SUPERVISOR_ARCHIVE_F2_FORMAT_BLOCKER.md) adds scope-separated
global/segment indexes, record-kind-bound hot/archive locations and count equations, and a distinct
generation-one migration-genesis checkpoint with generated answers. The follow-on valid-v1
mapping contradiction is also passively resolved with an explicit absent/present lifecycle union,
a typed lifecycle-record anchor on the present arm, and independent attempt/lifecycle counts. The
executable [F2 v1 mapping resolution](../SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md) retains the
real crash witness and no-invention rule. The
[stateful F2 result](../SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md) now implements only the owner-
asserted v1-to-v2 file migration and empty-archive full verifier, with strict reconstruction,
downgrade refusal, and local fault/corruption/capacity/process-death evidence. The follow-on
[stateful F3 result](../SUPERVISOR_ARCHIVE_F3_ACTIVATION_RESULT.md) implements one owner-asserted
sealed immutable-segment prepare/verify/publish/activate transaction. It retains every selected
record and visible tombstone, publishes before reference, atomically installs the generation-two
checkpoint, reports valid orphans without deletion, and refuses missing/corrupt/substituted
referenced bytes without fallback. It performs no retained lookup, v2 authority mutation, second
segment, backup, orphan cleanup, or adapter call. The follow-on
[read-only F4A result](../SUPERVISOR_ARCHIVE_F4A_LOOKUP_RESULT.md) routes exact retained objects,
replays, passive collision checks, and hot-only recovery through freshly full-verified retained-
global typed locations. It refuses corrupt/substituted/stale/wrong-domain storage without fallback
or rewrite and performs no v2 mutation, new tombstone commit, second segment, or adapter call. The
[F4B implementation review](../SUPERVISOR_ARCHIVE_F4B_MUTATION_BLOCKER.md) retains that former
contradiction. The [F4B result](../SUPERVISOR_ARCHIVE_F4B_MUTATION_RESULT.md) now uses the selected
independent append-only source for same-transaction issuance, full hot/archive reconstruction, and
historical `superseded-by-current` lookup without a false lifecycle record. Corrupt, omitted,
partial, substituted, or cross-linked tombstone state refuses without rewrite or fallback. The
[F4C result](../SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md) now carries those independent tombstones
and complete cohorts through deterministic second/later immutable-segment activation, exact
segment 64/65 capacity behavior, and old-or-complete-new reopen without deletion or adapter calls.
The [F5 result](../SUPERVISOR_ARCHIVE_F5_BACKUP_RESULT.md) now retains an exact manifest-last copy
of active v2 plus every referenced segment, full copied-set reopen, read-only exact-anchor restore
admission, fixed-shape offline inventory, and explicit removal of only a sealed known-unreferenced
segment after live and supplied-backup reference scans. Unknown, corrupt, mixed-generation, and
cross-installation artifacts refuse and remain evidence. The finite fixed-store checkpoint still
provides no production engine, installed multi-process/protected-root evidence, power-loss result,
restore activation/anti-rollback mechanism, referenced-history deletion, continuous service,
consumer, or guest evidence.

### Non-guarantees

Capsule does not prove:

- guest code performs the intended task correctly or is aligned;
- a signed claim or DID/credential describes a true real-world identity;
- user presence means the user understood source or every consequence;
- permitted outputs/metadata/timing contain no copied or encoded input;
- pattern redaction finds every sensitive value;
- supported runtime, parser, kernel, backend, hypervisor, Secure Enclave, or cryptographic library
  has no unknown vulnerability;
- hostile guest/host/job domains cannot leak through speculative execution, shared caches, branch
  predictors, memory deduplication, or another microarchitectural channel;
- elevated Accessibility, Screen Recording, task-port, or comparable user-granted authority cannot
  influence or spoof approval UX unless an exact posture and retained evidence say otherwise;
- a signature remains trustworthy after private-key or signer-logic compromise;
- trust epochs defeat coherent rollback without a stronger checkpoint;
- receipt claims are independently true merely because signatures verify;
- source behaves identically across runtimes/platforms;
- inputs already supplied through an AI client's attachment flow were never in model context.

## Severity Calibration

Severity assumes the vulnerable path is reachable in the documented deployment. Development-only
spike failures without product authority may be lower, but a false production security claim can
raise impact.

### Critical

- Guest escape to arbitrary host code execution, credentials, or filesystem access
- Supervisor compromise or a non-Supervisor route to create hostile guests
- Agent/daemon ability to approve plans, use Approval keys, or reset consumed grants
- Execution of substituted/unregistered bytes through approval, canonicalization, or replay bypass
- Agent/daemon retrieval of arbitrary user-only input/output content
- Forged validated receipt/terminal success that defeats both Broker and Supervisor evidence
- Cross-user/cross-tenant compromise in a future hosted authoritative deployment

### High

- Reliable unauthorized network, metadata-service, host-IPC, or management-channel access
- Acceptance of revoked/wrong-purpose/wrong-installation/wrong-epoch approval or trust state
- Live host-file mounts or content-handle confusion exposing more than granted bytes
- Resource-control bypass materially affecting host availability outside approved budgets
- Runtime component substitution, partial update, or debug state accepted as validated posture
- Teardown failure hidden as success while hostile execution may remain active
- User-only artifact release from an integrity-failed or indeterminate attempt

### Medium

- Bounded output/metadata disclosure to the wrong audience without broader authority
- Receipt, metric, profile-state, or transcript-integrity defect that materially weakens
  auditability but does not authorize execution
- Denial of service confined to one local job and approved resource budget
- TUF/profile freshness misclassification that downgrades evidence without activating untrusted code
- Unsafe approval rendering that misleads but still requires separate exploitation of authority

### Low

- Non-sensitive diagnostic or documentation inaccuracy with no effect on authorization, isolation,
  content, trust, posture, or cleanup
- Development-only availability failure that cannot execute guest code or alter production evidence
- Harmless metadata mismatch rejected closed

## Security posture

Posture is multidimensional rather than one “secure” tier:

- isolation assurance;
- runtime-integrity evidence mode;
- trust freshness;
- distribution identity.

Only exact pinned configurations with complete required rows in the control-evidence matrix may use
`validated-local` or stronger isolation language. Every receipt and UI preserves the underlying
dimensions and known limitations.
