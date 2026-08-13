# Apple Development provisioning and installed-test plan (Team `3DDR84M4JS`)

Date: 2026-08-04; reconciled 2026-08-05 after I1B/R3

Status: retained historical provisioning proposal plus current reconciliation. The original
documentation slice did not mutate the Apple Developer portal, Keychain, profiles, signed bytes,
or installed services. Accepted ADR-0037 later selected the containing topology, and the exact
Apple Development I1B/Source Validator R3 installed experiment is `PASSED` using three explicit
Team-`3DDR84M4JS` profiles. Sections that describe the earlier Option A/Option B choice or the S5
matrix are retained planning history, not current instructions. Proposed ADR-0045 now selects a
versioned Supervisor-authority-epoch candidate and inert experiment packet. Installed I2B remains
`BLOCKED` on separately authorized epoch-one App IDs/profiles, signing, disposable-container
nonmembership, and the later Coordinator/bootstrap evidence gates.
The completed experiment's public evidence is retained in the
[pinned R3 archive](https://github.com/Shrimpworks/capsule-experiments/tree/0944ffd8cfd01ec23e4ae99138b0931d56804077/experiments/completed-compiled-artifact-payloads/payloads/capsule-corp/artifacts/macos-i1b-r3-signed-development-composition).

The canonical practical inventory, Team-ID reconciliation, replacement-input checklist, environment and
component matrices, credential/key custody rules, and safe setup/verification commands now live in
[Apple certificates, credentials, identifiers, entitlements, and Capsule keys](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md).
This retained plan remains the detailed historical installed-test proposal, W4/3DDR discovery,
and current provisioning reconciliation.

Identity reconciliation: the earlier inference that the parenthesized value in an Apple
Development certificate display name was the signing Team ID is false. Certificate SHA-1
`1638CFBD9250A00B4DBD81AE8FD1C790B42F61E3` is displayed as
`Apple Development: Dylan Steele (W4QUR9FUL4)`, but its X.509 subject OU and an exact signed-byte
TeamIdentifier are `3DDR84M4JS`. On 2026-08-04 the account's Apple Membership Details independently
identified `3DDR84M4JS` as the Apple Developer Program Team ID for the Individual membership.
`W4QUR9FUL4` is therefore a certificate common-name/member display suffix, not Capsule's Team ID.
Installed G3/I2B remains `BLOCKED` on its protected-root/bootstrap composition and exact
Coordinator/bootstrap profiles;
see the retained
[historical G3 result](https://github.com/Shrimpworks/capsule-experiments/blob/3e9c9cbc3e0314439771151f1fd99c2b3a5a50b9/experiments/supervisor-owner-lock-installed-g3/RESULTS.md).

The 2026-08-04 exact-selector follow-up also proved that the `codesign` default designated
requirement repeats the misleading W4 common name without binding the certificate Team OU. That
common name is not Team admission evidence. Every future run must require Team OU and emitted
TeamIdentifier `3DDR84M4JS`, the role signing identifier, enrolled exact CDHashes, and the effective
entitlement digest. The cached 3DDR Gate B profiles still do not match the Capsule product bundle
IDs and must not be reused as generic Team profiles.

A user-run `security find-identity -v -p codesigning` on 2026-08-04 reported a new valid Apple
Development identity, SHA-1 `80A4969BCD1B3926020888094B9D812A283D3793`, and Keychain Access showed
it under My Certificates. The later separately authorized I1B/R3 experiment selected that identity
and retained exact subject/profile/signed-byte readback for TeamIdentifier `3DDR84M4JS`. This does
not authorize later use. The older Apple Development identity `1638...61E3` and Developer ID
Application identity `AD70...F750` remain distinct and were not selected for I1B/R3.

Reviewer: Claude, independent read-only planning at the request of the Capsule orchestrator
(codex).

## Scope and method

This plan originally relied on the false display-name inference that `W4QUR9FUL4` was the signing
Team. Exact G3 readback and the later Apple Membership Details screen establish `3DDR84M4JS` as the
Team ID. It read `AGENTS.md`,
[`ARCHITECTURE.md`](ARCHITECTURE.md), [`AUTHENTICATED_LOCAL_IPC_PLAN.md`](AUTHENTICATED_LOCAL_IPC_PLAN.md)
and its [S1 consistency stop](AUTHENTICATED_LOCAL_IPC_S1_CONSISTENCY_STOP.md),
[`SUPERVISOR_OWNER_LOCK_PLAN.md`](SUPERVISOR_OWNER_LOCK_PLAN.md),
[`security/THREAT_MODEL.md`](security/THREAT_MODEL.md), the Gate C P0-4A checkpoint documents
(`GATE_C_P0_RECONCILIATION.md`, `GATE_C_READINESS_CHECKPOINT.md`, `WORKSTREAM_EVIDENCE_LEDGER.md`,
`FEASIBILITY_SPIKES.md`), ADR-0012, ADR-0013, ADR-0018, ADR-0021, and the currently **Proposed**
ADR-0029, plus the retained real spikes in the commit-pinned
[`gate-c-installed-development-topology`](https://github.com/Shrimpworks/capsule-experiments/tree/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-installed-development-topology)
(P0-4A) and
[`macos-authority-separation`](https://github.com/Shrimpworks/capsule-experiments/tree/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/macos-authority-separation)
(Gate B, the same Team but different experimental bundle identifiers and profiles).

The original planning task did not access the Apple Developer portal, create any identifier,
download any provisioning profile, or sign any component. The later separately authorized I1B/R3
experiment did those exact scoped operations and retained its public evidence independently.

The only real product binary in the repository today is `cmd/capsuled/main.go`. There is no Swift
Broker source, no Supervisor native front-end source, and no Xcode project for any product
component — only spike/experiment Xcode projects and Info.plist templates. Everything below is a
proposal against an architecture that exists on paper plus disposable spikes, not against real
target-of-signing binaries.

## Resource state (already confirmed, not re-derived)

Per S5 and `WORKSTREAM_EVIDENCE_LEDGER.md`:

- Apple Membership Details confirms Individual Apple Developer Program membership under Team
  `3DDR84M4JS`.
- `security find-identity -v -p codesigning` reported two valid Apple Development identities whose
  common names end in `W4QUR9FUL4`: SHA-1 `1638...61E3` and `80A4...D3793`. That suffix is not the
  Team ID. Exact I1B/R3 readback later selected `80A4...D3793` and emitted TeamIdentifier
  `3DDR84M4JS` for that experiment only.
- The earlier Xcode cache contained three nonmatching Gate B/wildcard profiles. I1B/R3 later
  created and selected three exact explicit-App-ID profiles for Broker, daemon, and Supervisor;
  their public metadata and hashes are retained in the pinned R3 archive.
- A separate Developer ID Application identity SHA-1 `AD70...F750` also belongs to Team
  `3DDR84M4JS`. It is later distribution authority and does not authorize Developer ID signing or
  make notarization current.
- The separately authorized I2B3 preflight created and cryptographically reconciled exact
  Coordinator profile `c0446281-ba4b-451b-8c73-9ee9d8ef97a2` and Supervisor bootstrap profile
  `c45a058b-ffdd-4a6b-bd8c-d746772a2702` under the selected Team-3DDR certificate. Signed effective
  entitlement readback passed, but the required stale-profile test then blocked I2B3 before keys,
  service registration, installation, or root creation.
- No paid owned clean-host / minimum-OS validation hardware is currently planned.

Any later App ID/profile capacity and portal availability must be observed under Team
`3DDR84M4JS` during the separately authorized task; I1B/R3 evidence does not establish future
portal state.

## Component inventory and selected bundle topology

The original plan, following `ARCHITECTURE.md` and ADR-0029, described three ordinary local
authorities plus one already-decided non-authority:

| Authority | Language/packaging (as decided) | Source |
| --- | --- | --- |
| Agent-facing daemon | Go (`cmd/capsuled`) | ADR-0018 |
| Trusted Host Broker (Approval Broker) | preferably Swift/native | ADR-0018 |
| Execution Supervisor | one unprivileged per-user `SMAppService.agent` executable containing a small in-process native C/Objective-C XPC/Security front end linked with the existing Go core | ADR-0029 |

ADR-0029 is explicit that the Supervisor's native front end is not a separate process, binary, or
bundle — it shares the Supervisor's code identity, and the ADR explicitly rejects a "separate Swift
XPC front end plus Go Supervisor service" topology. The TypeScript Source Preparer's native front
(ADR-0032) remains `BLOCKED`, not selected for the first release. Proposed ADR-0038 later adds one
installation-only Trust Coordinator as the eighth I2 role. Any other signable authority still
requires its own ADR.

The original plan considered Option A (daemon as main app) and Option B (a separate installer
app). Accepted ADR-0037 supersedes that choice: one visible Broker `Capsule.app` contains the
separately enrolled daemon and Supervisor roles. Neither historical option authorizes registering
`com.capsulecorp.capsule` or `com.capsulecorp.capsule.installer`.

| Component | Selected bundle/signing ID | Current status |
| --- | --- | --- |
| Visible Approval Broker app | `com.capsulecorp.capsule.broker` | Selected by ADR-0037; exact development profile passed I1B/R3 |
| Agent-facing daemon | `com.capsulecorp.capsule.daemon` | Selected by ADR-0037; exact development profile passed I1B/R3 |
| Execution Supervisor | `com.capsulecorp.capsule.supervisor` | Historical stable identity; I1B/R3 and I2B3 preflight evidence retained, but Proposed ADR-0045 classifies it as legacy residue rather than an admitted authority epoch |
| Daemon-private Source Validator XPC service | `com.capsulecorp.capsule.source-validator.daemon.v1` | Nested signed identity passed inactive-policy R3 without an independent portal profile |
| Broker-private Source Validator XPC service | `com.capsulecorp.capsule.source-validator.approval-broker.v1` | Nested signed identity passed inactive-policy R3 without an independent portal profile |
| Daemon-role parser child | `com.capsulecorp.capsule.source-validator-parser.daemon.v1` | Nested signed identity only; parser spawn remained disabled in R3 |
| Broker-role parser child | `com.capsulecorp.capsule.source-validator-parser.approval-broker.v1` | Nested signed identity only; parser spawn remained disabled in R3 |
| Trust Coordinator | `com.capsulecorp.capsule.trust-bootstrap.v1` | Historical unlaunched I2B3 preflight identity and legacy residue; epoch one uses a fresh versioned identity only after separate authorization |
| Authority-epoch-one Supervisor candidate | `com.capsulecorp.capsule.supervisor.authority-e1` | Selected only by Proposed ADR-0045; inert construction and exact no-launch Apple Development profile/signature readback are `PASSED`; container evidence remains `BLOCKED` |
| Authority-epoch-one Coordinator candidate | `com.capsulecorp.capsule.trust-bootstrap.authority-e1` | Selected only by Proposed ADR-0045; inert construction and exact no-launch Apple Development profile/signature readback are `PASSED`; Coordinator launch and later service/root evidence remain `BLOCKED` |
| Native front-end (Supervisor's C/Obj-C shim) | N/A — not a separate bundle | Shares `com.capsulecorp.capsule.supervisor`'s code identity |
| Guest launcher/runner (flagged for completeness, outside this provisioning plan) | e.g. `com.capsulecorp.capsule.runner` | Would eventually hold `com.apple.security.hypervisor`; blocked behind P0-0..P0-3 per `GATE_C_P0_RECONCILIATION.md` |

The `com.capsulecorp.capsule.*` namespace is not invented — it is lifted directly from ADR-0029's
own Mach service names (below). The spike namespace `com.capsulecorp.spike.*` and the Gate B
historical namespace are disposable/experiment-only and tied to the wrong team; neither is reused
for product App IDs. I1B/R3 created the three containing-role App IDs listed above; this repository
retains public metadata, not portal authority or raw profiles.

ADR-0036 narrows daemon packaging for the Source Validator path: the daemon execution role has the
enrolled containing-bundle identity `com.capsulecorp.capsule.daemon` even if a user-visible main app
also exists. Its private validator service is not embedded in or reachable from the installer/main
app as an alternate peer. R1-R2 used no signing identity; the later exact R3 signed installed
inactive-policy composition is `PASSED`. Active parser/product work remains blocked and the passed
R3 evidence authorizes no later portal or signing action.

## Required App IDs and provisioning profiles (development-only)

The development lane requires an explicitly authorized Apple Development identity whose
certificate subject OU and emitted TeamIdentifier both equal `3DDR84M4JS`; every profile below is
then a macOS Development profile,
not Developer ID/Distribution:

| Explicit App ID | Capabilities in exact evidence/candidate | Profile type | Current status |
| --- | --- | --- | --- |
| `com.capsulecorp.capsule.broker` | App Sandbox | macOS App Development | Exact profile passed I1B/R3 |
| `com.capsulecorp.capsule.daemon` | App Sandbox | macOS App Development | Exact profile passed I1B/R3 |
| `com.capsulecorp.capsule.supervisor` | App Sandbox | macOS App Development | Exact I1B/R3 profile passed; new I2B3 bootstrap profile/signing preflight also passed, but the old profile retained write access to the stable private container |
| `com.capsulecorp.capsule.trust-bootstrap.v1` | App Sandbox plus candidate bootstrap App Group and Coordinator-only Keychain group | macOS App Development | Exact App ID/profile and signed-entitlement preflight passed; Coordinator was not launched and no key/container state was created |

For each newly authorized App ID: use one macOS App Development provisioning profile scoped to the
exact App ID, matching authorized Apple Development certificate after exact verification, and the
owned test Mac's device UDID. Do not reuse I1B profiles for a changed I2B3 entitlement set. The
[I2B3 stale-profile result](MACOS_INSTALLATION_I2B3_SIGNING_PREFLIGHT_AND_STALE_PROFILE_BLOCKER.md)
additionally proves that merely issuing a new profile for the same Supervisor signing identifier
does not version its App Sandbox private container.
Development profiles are inherently single-Team, non-notarizable, and Gatekeeper-rejecting by
design — P0-4A's `spctl --assess` already returned status 3 against ad-hoc-signed bytes, and this
plan should not chase a Gatekeeper-pass result with a Development profile.

No product Xcode project exists yet. I1B/R3 used a deliberately inert provisioning probe and exact
manual gates; future product targets must bind the selected identifiers and exact profiles rather
than rely on a wildcard or reserve speculative IDs. The archived
[`GateBProvisioned.xcodeproj`](https://github.com/Shrimpworks/capsule-experiments/tree/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/macos-authority-separation/Provisioned/GateBProvisioned.xcodeproj)
shows the historical `CODE_SIGN_STYLE = Automatic` and explicit `PRODUCT_BUNDLE_IDENTIFIER`
pattern only; it is not a current product target.

## Historical minimum-entitlement proposal

Grounded in the real per-role split already in the archived
[`Provisioned/Entitlements`](https://github.com/Shrimpworks/capsule-experiments/tree/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/macos-authority-separation/Provisioned/Entitlements)
and in
`security/THREAT_MODEL.md`'s local process/storage boundary and mandatory security properties.
This table predates I1B/R3 and is not its effective entitlement inventory. I1B/R3 retained the
exact development-only effective entitlements in the pinned R3 archive. Later I2B3 must consume
the narrower candidate groups and identities frozen by I2B2 rather than activate this historical
table wholesale.

| Role | Required entitlement keys | Why |
| --- | --- | --- |
| Main app / daemon | `com.apple.security.app-sandbox = true` | Matches Gate B's `Daemon.entitlements`; `security/INSTALLATION_TRUST.md` excludes the daemon from all Keychain/root/operational key access |
| | Local IPC to reach `...supervisor.daemon.v0` | Open item — under App Sandbox this likely needs `com.apple.security.temporary-exception.mach-lookup-global-name` naming the exact service; ADR-0029 hasn't reached S3 yet, so S5 must determine this empirically |
| Execution Supervisor | `com.apple.security.app-sandbox = true` | "Unprivileged per-user," no host-root helper (`EXECUTION_SUPERVISOR.md`) |
| | `keychain-access-groups` = one Supervisor-owned, epoch-scoped group, e.g. `$(AppIdentifierPrefix)com.capsulecorp.capsule.supervisor.evidence.epoch-<N>` | Matches Gate B's `Supervisor.entitlements`; ADR-0021 requires a fresh, non-migrated group per identity-changing security epoch |
| | Listener side of the two Mach services (§6) via the embedded `LaunchAgent.plist` | Confirmed real pattern from the P0-4A spike |
| Approval Broker | `com.apple.security.app-sandbox = true` | Same least-privilege baseline |
| | `keychain-access-groups` = one Broker-owned, epoch-scoped group, disjoint from the Supervisor's | Matches Gate B's `Broker.entitlements`; Gate B proved cross-group access fails with `errSecMissingEntitlement` (`-34018`), the desired outcome |
| | `com.apple.security.files.user-selected.read-only` | Broker "safely snapshots user-selected regular-file data-fork bytes" (`ARCHITECTURE.md`) |
| | LocalAuthentication use (implicit via the Keychain-group entitlement + `SecAccessControl`) | Approval requires "fresh user presence for every v0 plan" (Threat Model) |

All roles in this historical proposal should also set, as build settings:
`ENABLE_HARDENED_RUNTIME = YES` and no `get-task-allow` in shipping-posture builds.

## Prohibited entitlements per role

| Role | Prohibited entitlement | Cited reason |
| --- | --- | --- |
| Daemon | `keychain-access-groups` (any) | `security/INSTALLATION_TRUST.md`: daemon excluded from key references |
| Daemon | `com.apple.security.hypervisor` | `ARCHITECTURE.md`: daemon "must not... launch a backend" |
| Daemon | `com.apple.security.files.user-selected.*` | File selection belongs solely to the Broker |
| Supervisor | `com.apple.security.hypervisor` | Belongs to the separate runner/launcher process per the libkrun spikes, not the Supervisor itself — giving it directly collapses a deliberately separated privilege boundary |
| Supervisor | `com.apple.security.get-task-allow` | Shipping components ship without it per the Threat Model |
| Supervisor | `com.apple.security.files.user-selected.*` / general filesystem entitlement | Supervisor "never gets an original host path" (Threat Model) |
| Supervisor | Any shared/broad `com.apple.security.application-groups` | `ARCHITECTURE.md`: "no shared app group unless a narrow capability requires it" |
| Supervisor | Broker's or another epoch's Keychain group | ADR-0021 requires disjoint, epoch-fresh groups |
| Approval Broker | `com.apple.security.hypervisor` | Broker never launches a backend or guest |
| Approval Broker | Any daemon- or agent-facing network listener entitlement | `ARCHITECTURE.md`: Broker "must not... expose agent endpoints" |
| Approval Broker | `com.apple.security.cs.disable-library-validation`, `.allow-unsigned-executable-memory`, `.allow-jit` | Not evidenced as required anywhere; Hardened Runtime posture requires the minimal exception set |
| All four | `com.apple.security.temporary-exception.*` beyond the single named-Mach-service exception in §4 (if even required) | Broad sandbox temporary exceptions are exactly the unreviewed privilege expansion the least-privilege rules exist to prevent |

## XPC/Mach service identities and peer-validation requirements

From Accepted ADR-0029, which is the governing document that names services and the validation
order. Two role-specific Mach services, recorded in the `InstallationManifest`:

| Service name | Enrolled peer | Closed calls |
| --- | --- | --- |
| `com.capsulecorp.capsule.supervisor.daemon.v0` | daemon | `RegisterPlanV0`, `RequestAttemptV0` |
| `com.capsulecorp.capsule.supervisor.broker.v0` | Approval Broker | `GetRegisteredPlanV0`, `SubmitApprovalV0` |

Peer-validation order any provisioning/entitlement plan must support (ADR-0029):

1. install the exact listener peer requirement before listener activation;
2. require XPC-observed EUID and audit-session ID to match the enrolled current login session;
3. derive `SecCode` via `SecCodeCreateWithXPCMessage` from the exact XPC message and re-validate
   Team ID, role-specific signing identifier, accepted CDHash set, entitlement digest, Hardened
   Runtime constraints, and debug-state rejection;
4. acquire a fixed connection/global flow-control slot;
5. validate the closed XPC dictionary shape before decoding the method body;
6. require the request's installation ID and epoch sequence/digest to match already-open Supervisor
   state; and
7. dispatch only the method implied by service + closed message tag.

ADR-0029's own S3 slice (native authentication/cap harness) has not run yet; the most recent passive
evidence exercised only requirement parsing on ad-hoc processes and "does not authenticate a peer,
inspect or use an Apple-issued identity, establish Team or distribution enrollment, or validate
production IPC." Any Team-3DDR-enrolled test of the table above is new work, not a rerun.

## Team ID / update-channel / CDHash / update expectations

- `InstallationManifest` must record, per component: signing identifier, Team ID `3DDR84M4JS`, peer
  code requirements, and the exact
  active CodeDirectory hashes — the designated requirement/Team match never replaces the active
  CDHash set (ADR-0029).
- Every trusted IPC call, plan registration, approval, attempt, and receipt binds the active epoch
  (installation ID + sequence-ordered epoch digest); a component from a different epoch must fail
  closed as partial-update/stale-peer (`security/INSTALLATION_TRUST.md`).
- Epochs are sequence-ordered, not rollback-proof — never describe a development test result as defeating
  rollback; only a non-rollbackable anchor or external witness would, and none exists here.
- Keychain access groups are epoch-scoped, not identity-scoped (ADR-0021): every identity-changing
  epoch should provision a fresh access group and a fresh, non-migrated Secure Enclave key. Gate B's
  evidence proved a stable Keychain group is a Team/profile boundary, not a build/epoch boundary.
- Update-channel expectations follow `UPDATE_AND_RECOVERY.md`'s 14-step prepared-update ceremony;
  none of it exists as running code yet, so a development install validates the mechanics categories it
  implies (old/new component pairing, stale-peer refusal, partial-update fail-closed), not the
  ceremony itself.
- Every test rebuild changes the CDHash and must rebuild-and-repin at every step that changes any
  signed byte, per the explicit warning in the archived
  [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-installed-development-topology/RESULTS.md).

## Historical S5 installation and rollback test matrix

Narrowed to what's reachable with only a free/Individual Apple Development identity and a single
owned Mac unless noted.
I1B/R3 passed its own exact signed-install/refusal/reachability/cleanup matrix; it did not claim
every row below. Later installed slices must select and authorize their exact rows from current
I2B/I3 plans rather than rerun this pre-I1B matrix by implication.

| # | Scenario | Oracle / expected result | Source |
| --- | --- | --- | --- |
| 1 | Fresh install: register Supervisor as embedded `SMAppService.agent`, explicit `launchctl kickstart` activation | Registration returns `enabled`; explicit kickstart starts the process — P0-4A found pure on-demand Mach activation not established | P0-4A `RESULTS.md` |
| 2 | Correct-peer daemon → Supervisor `...daemon.v0` call | Delivered, passes all checks | ADR-0029 |
| 3 | Correct-peer Broker → Supervisor `...broker.v0` call | Delivered, passes all checks | ADR-0029 |
| 4 | Wrong-role peer | Dropped by XPC listener requirement before delivery | ADR-0029 refusal table |
| 5 | Same-identifier, different-CDHash peer (stale rebuild) | Refused / connection interrupted | P0-4A `RESULTS.md` |
| 6 | Missing/mixed/unexpected installed component | Bounded refusal, no service use | P0-4A `RESULTS.md` |
| 7 | Wrong EUID / wrong Aqua session | `AUTHENTICATION` refusal, zero body copy | ADR-0029 refusal table |
| 8 | Debug-entitled or actively-debugged peer | Rejected | ADR-0029 |
| 9 | Same-session crash/reconnect | Fresh PID/instance, same-session recovery succeeds | P0-4A `RESULTS.md` |
| 10 | Duplicate Supervisor start / owner-lock contention | Second owner refuses before any store mutation | `SUPERVISOR_OWNER_LOCK_PLAN.md` |
| 11 | Owner-lock object relocated/replaced by another same-UID process | Repair-required; refuses rather than silently accepting a new inode | `SUPERVISOR_OWNER_LOCK_PLAN.md`; `security/INSTALLATION_TRUST.md` |
| 12 | Upgrade: old daemon + new Supervisor, and reverse | Refused / fenced until one active epoch | Threat Model conformance table |
| 13 | Interrupted/partial update | `repair-required`; daemon cannot clear it | `UPDATE_AND_RECOVERY.md` |
| 14 | Downgrade / coherent local rollback | Not achievable as a positive containment test — test only that a mismatched epoch is detected and fenced, not prevented | `security/INSTALLATION_TRUST.md` |
| 15 | Revoked/expired provisioning profile | Standard macOS launch/re-sign failure, not a Capsule mechanism | N/A |
| 16 | Wrong Team entirely | Refused before delivery | Analog of Gate B's stale/wrong-Team cases |
| 17 | Tampered binary vs. byte-identical copy | Tamper refused; byte-identical copy accepted as same code identity | ADR-0029 |
| 18 | Reboot / logout-login / fast user switch | Should reconnect correctly — not yet exercised anywhere in this repo | P0-4A `RESULTS.md` open item |
| 19 | Sleep/wake | Not evidenced anywhere in the repo yet | Same |
| 20 | Locked Keychain during approval-signing | Should fail closed, not silently proceed | Threat Model |

## What's testable on the current single Mac vs. what needs another machine

**Historically assessed as testable on this one Mac:** rows 1–13, 16, 17, and 20 above once their
exact current-slice profiles and entitlements exist. I1B/R3 later proved its narrower
development-signed install, private-XPC reachability, refusal, `SMAppService` registration, and
cleanup scope; it did not complete this whole S5 matrix. Ordinary Supervisor IPC, operational
Keychain groups, epoch rotation, protected state, and later fault rows require their current
I2B/I3 authorization and evidence.

**Needs a second machine** (not necessarily clean): multi-machine trust-epoch behavior; genuine
fast-user-switch/session concurrency under real load.

**Needs a genuinely clean host, cannot be done on this Mac at all:** Gatekeeper/quarantine/
translocation on first launch (`spctl --assess` on the build machine is not equivalent to a fresh
machine's first-launch path); supported macOS floor validation (P0-4A's effective floor was
measured at macOS 26.0, driven by a dependency, not `LSMinimumSystemVersion`); install/update races
across genuinely independent boot/session state.

## Claims that stay impossible without clean-host testing or authorized Developer ID evidence

Per this plan's explicit scope (no authorized Developer ID/notarization use):

- Notarization and stapling.
- Gatekeeper assessment as an end user would experience it — a Development-signed, unnotarized app
  is supposed to fail Gatekeeper; P0-4A already confirmed that (`spctl` status 3).
- "Validated"/"production-ready" language of any kind for any of this — `AGENTS.md` and the Threat
  Model both forbid that without exact mechanism plus retained adversarial evidence. I1B/R3 passed
  only its exact development experiment; the parent product and later installed slices remain
  blocked.
- Coherent-rollback resistance (needs a non-rollbackable anchor or external witness; neither exists
  in this repo's design yet).
- Minimum supported macOS floor, and cross-machine/clean-host install-race/update evidence (§9).
- Anything about the actual guest/backend (libkrun/HVF, Apple Containerization, gVisor) —
  orthogonal to this plan; P0-4A "cannot admit the backend" per `GATE_C_P0_RECONCILIATION.md`.

## Current action list for the user/orchestrator

1. **Preserve completed I1B/R3 evidence.** Team `3DDR84M4JS`, SHA-1 `80A4...D3793`, the ADR-0037
   topology, and three exact Broker/daemon/Supervisor profiles passed only the retained
   development-signed, installed, execution-disabled scope. Do not repeat or widen it implicitly.
2. **Do not register the historical Option A/Option B identifiers.** ADR-0037 selects the visible
   Broker app plus separately enrolled embedded daemon and Supervisor identities.
3. **Keep the central operational inventory deferred for the current one-maintainer stage.** The
   pinned R3 archive and evidence ledger retain current redacted metadata. Create a shared location
   when a second maintainer, production signing, or CI notarization requires it.
4. **Execute the ADR-0045 gates in order.** First retain and review the inert epoch-one packet.
   Then separately authorize only its exact versioned App IDs/profiles, signing/readback, and
   disposable-container nonmembership matrix. If that passes, separately authorize caller/key,
   App Group/service/root, and descriptor-relative evidence; keep runtime, backend, guest, and
   attempts absent.
5. **Retain any later credentialed evidence outside chat.** Archive the exact public metadata and
   results in `capsule-experiments`, pin its immutable commit here, and update the evidence ledger.
   Raw profiles, private keys, and credentials remain excluded.
6. **Keep Developer ID/notarization separate and deferred.** Distribution work requires its own
   authorization, custody design, clean-host matrix, and evidence; I1B/R3 does not advance it.

## Explicit non-goals

This plan does not authorize a new Apple ID/App ID/profile, portal action, profile download,
signature, Keychain/Secure Enclave key, installed mutation, or claim that any Capsule component is
"validated," "secure," or "production-ready." The original planning slice performed none of
those actions; the later separately authorized I1B/R3 experiment performed only its exact retained
development scope. ADR-0037 supersedes the Option A/Option B choice. Developer ID/notarization
remains explicitly excluded and deferred.
