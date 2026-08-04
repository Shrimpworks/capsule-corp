# Apple Development provisioning and installed-test plan (Team `W4QUR9FUL4`)

Date: 2026-08-03

Status: retained read-only provisioning and installed-test plan. Nothing in this repository, the
Apple Developer portal, or any Keychain was modified to produce it. Every bundle ID, entitlement,
service name, and test step below is a proposal for the user/orchestrator to execute later against
[S5: Apple Development installed matrix](AUTHENTICATED_LOCAL_IPC_PLAN.md#s5-apple-development-installed-matrix),
not a decision this repository has already made unless the exact doc/ADR making it is cited.

Correction after exact G3 discovery: the earlier inference that the displayed certificate was a
W4 Team identity is false. Certificate SHA-1
`1638CFBD9250A00B4DBD81AE8FD1C790B42F61E3` is displayed as
`Apple Development: Dylan Steele (W4QUR9FUL4)`, but its X.509 subject OU and an exact signed-byte
TeamIdentifier are `3DDR84M4JS`. No cached profile belongs to W4. Treat the bundle/profile tables
below as future W4 requirements only; do not use that certificate or the historical profiles as
W4 evidence. Installed G3 is currently `BLOCKED`; see the retained
[historical G3 result](https://github.com/Shrimpworks/capsule-experiments/blob/3e9c9cbc3e0314439771151f1fd99c2b3a5a50b9/experiments/supervisor-owner-lock-installed-g3/RESULTS.md).

The 2026-08-04 exact-selector follow-up also proved that the `codesign` default designated
requirement repeats the misleading W4 common name without binding the certificate Team OU. That
requirement is not W4 admission evidence. The future W4 run must require the explicit Team OU,
emitted TeamIdentifier, role signing identifier, enrolled exact CDHashes, and effective entitlement
digest. Both standard local profile caches were checked; only the same three 3DDR profiles exist.

Reviewer: Claude, independent read-only planning at the request of the Capsule orchestrator
(codex).

## Scope and method

This plan originally relied on the display-name inference for an Apple Development identity
expected to belong to Team `W4QUR9FUL4` (not a paid Developer ID/notarization identity). Exact G3
readback later disproved that Team inference. It read `AGENTS.md`,
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
(Gate B, historical Team `3DDR84M4JS`).

It did not access the Apple Developer portal, create any identifier, download any provisioning
profile, sign any component, or modify the repository beyond adding this document.

The only real product binary in the repository today is `cmd/capsuled/main.go`. There is no Swift
Broker source, no Supervisor native front-end source, and no Xcode project for any product
component — only spike/experiment Xcode projects and Info.plist templates. Everything below is a
proposal against an architecture that exists on paper plus disposable spikes, not against real
target-of-signing binaries.

## Resource state (already confirmed, not re-derived)

Per S5 and `WORKSTREAM_EVIDENCE_LEDGER.md`:

- The user confirms Individual membership under development Team `W4QUR9FUL4`; the locally
  available certificate does not belong to that Team at the code-signing boundary.
- `security find-identity -v -p codesigning` reports a certificate whose display name ends in W4,
  but its subject OU and signed-byte TeamIdentifier are `3DDR84M4JS`; it is not W4 evidence.
- Xcode 26.6 has three cached provisioning profiles, all belonging to the historical Team
  `3DDR84M4JS` (Gate B Broker, Gate B Supervisor, wildcard). None are reusable for W4.
- A separate Developer ID Application identity exists only for historical Team `3DDR84M4JS`. It is
  not W4 evidence and does not make Developer ID/notarization current for W4.
- No matching W4 signing certificate or W4 role profile is locally available.
- No paid owned clean-host / minimum-OS validation hardware is currently planned.

The task confirms Individual membership. Actual App ID/profile capacity and portal availability
must still be observed under the selected W4 team during certificate/profile creation; repository
text cannot establish Apple account state.

## Component inventory and bundle-identifier proposal

Per `ARCHITECTURE.md` and ADR-0029, there are exactly three local authorities plus one already-
decided non-authority:

| Authority | Language/packaging (as decided) | Source |
| --- | --- | --- |
| Agent-facing daemon | Go (`cmd/capsuled`) | ADR-0018 |
| Trusted Host Broker (Approval Broker) | preferably Swift/native | ADR-0018 |
| Execution Supervisor | one unprivileged per-user `SMAppService.agent` executable containing a small in-process native C/Objective-C XPC/Security front end linked with the existing Go core | ADR-0029 |

ADR-0029 is explicit that the Supervisor's native front end is not a separate process, binary, or
bundle — it shares the Supervisor's code identity, and the ADR explicitly rejects a "separate Swift
XPC front end plus Go Supervisor service" topology. The TypeScript Source Preparer's native front
(ADR-0032) remains `BLOCKED`, not selected for the first release. **There is no fourth signable
component today**; adding one requires its own ADR.

No ADR or architecture doc names a distinct product-facing container app separate from the daemon,
Broker, and Supervisor. This plan proposes bundle IDs for both plausible options rather than picking
one:

- **Option A** — the daemon is the main app (`cmd/capsuled` ships as the primary user-visible
  product and registers the Supervisor's `SMAppService.agent`).
- **Option B** — a separate installer/container app exists only to hold the Supervisor's embedded
  `LaunchAgent.plist` and perform first-install/update ceremonies (`UPDATE_AND_RECOVERY.md`'s
  "trusted installer" language is never equated with the daemon).

| Component | Proposed bundle ID | Notes |
| --- | --- | --- |
| Main app — Option A | `com.capsulecorp.capsule` | `cmd/capsuled`'s eventual signed form |
| Main app — Option B | `com.capsulecorp.capsule.installer` | Only if Option B is chosen |
| Agent-facing daemon (if packaged separately) | `com.capsulecorp.capsule.daemon` | Needs its own code identity regardless of A/B — ADR-0029 authenticates it as a distinct Supervisor peer |
| Execution Supervisor | `com.capsulecorp.capsule.supervisor` | Embeds the in-process native front end; no separate ID for that front end |
| Approval Broker | `com.capsulecorp.capsule.broker` | Swift/native per ADR-0018 |
| Native front-end (Supervisor's C/Obj-C shim) | N/A — not a separate bundle | Shares `com.capsulecorp.capsule.supervisor`'s code identity |
| Guest launcher/runner (flagged for completeness, out of this request's four components) | e.g. `com.capsulecorp.capsule.runner` | Would eventually hold `com.apple.security.hypervisor`; blocked behind P0-0..P0-3 per `GATE_C_P0_RECONCILIATION.md` |

The `com.capsulecorp.capsule.*` namespace is not invented — it is lifted directly from ADR-0029's
own Mach service names (below). The spike namespace `com.capsulecorp.spike.*` and the Gate B
historical namespace are disposable/experiment-only and tied to the wrong team; neither is reused
for product App IDs. None of the IDs above exist as registered App IDs yet.

## Required App IDs and provisioning profiles (development-only)

The future W4 lane requires an Apple Development identity whose certificate subject OU and signed
TeamIdentifier both equal `W4QUR9FUL4`; every profile below is then a macOS Development profile,
not Developer ID/Distribution:

| Explicit App ID | Capabilities | Profile type | Notes |
| --- | --- | --- | --- |
| Main app (per Option A/B choice) | App Sandbox | macOS App Development | Must be an explicit App ID, not wildcard — wildcard App IDs cannot carry Keychain Sharing or most entitlement-backed capabilities |
| `com.capsulecorp.capsule.daemon` | App Sandbox | macOS App Development | |
| `com.capsulecorp.capsule.supervisor` | App Sandbox, Keychain Sharing | macOS App Development | `SMAppService`/embedded-login-item capability is Xcode-managed, not an App ID toggle |
| `com.capsulecorp.capsule.broker` | App Sandbox, Keychain Sharing | macOS App Development | |

For each App ID: one macOS App Development provisioning profile scoped to the exact App ID, the one
matching W4 Apple Development certificate after reissue/verification, and this Mac's device UDID.
Development profiles are inherently single-Team, non-notarizable, and Gatekeeper-rejecting by
design — P0-4A's `spctl --assess` already returned status 3 against ad-hoc-signed bytes, and this
plan should not chase a Gatekeeper-pass result with a Development profile.

Given no Xcode project exists yet for these components, default to manual/explicit App ID + profile
creation in the portal first (reserving the bundle IDs above), then let Xcode's automatic signing
pick them up once real targets exist, mirroring the archived
[`GateBProvisioned.xcodeproj`](https://github.com/Shrimpworks/capsule-experiments/tree/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/macos-authority-separation/Provisioned/GateBProvisioned.xcodeproj)
`CODE_SIGN_STYLE = Automatic` with an explicit
`PRODUCT_BUNDLE_IDENTIFIER` per target.

## Minimum entitlements per role

Grounded in the real per-role split already in the archived
[`Provisioned/Entitlements`](https://github.com/Shrimpworks/capsule-experiments/tree/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/macos-authority-separation/Provisioned/Entitlements)
and in
`security/THREAT_MODEL.md`'s local process/storage boundary and mandatory security properties.

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

All four roles should also set, as build settings: `ENABLE_HARDENED_RUNTIME = YES` and no
`get-task-allow` in shipping-posture builds.

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

From the currently Proposed ADR-0029, which is the only doc that names services and the validation
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
production IPC." Any Team-W4-enrolled test of the table above is new work, not a rerun.

## Team ID / update-channel / CDHash / update expectations

- `InstallationManifest` must record, per component: signing identifier, Team ID (`W4QUR9FUL4` for
  every W4-scoped test — never mixed with `3DDR84M4JS`), peer code requirements, and the exact
  active CodeDirectory hashes — the designated requirement/Team match never replaces the active
  CDHash set (ADR-0029).
- Every trusted IPC call, plan registration, approval, attempt, and receipt binds the active epoch
  (installation ID + sequence-ordered epoch digest); a component from a different epoch must fail
  closed as partial-update/stale-peer (`security/INSTALLATION_TRUST.md`).
- Epochs are sequence-ordered, not rollback-proof — never describe a W4 test result as defeating
  rollback; only a non-rollbackable anchor or external witness would, and none exists here.
- Keychain access groups are epoch-scoped, not identity-scoped (ADR-0021): every identity-changing
  epoch should provision a fresh access group and a fresh, non-migrated Secure Enclave key. Gate B's
  evidence proved a stable Keychain group is a Team/profile boundary, not a build/epoch boundary.
- Update-channel expectations follow `UPDATE_AND_RECOVERY.md`'s 14-step prepared-update ceremony;
  none of it exists as running code yet, so a W4 install validates the mechanics categories it
  implies (old/new component pairing, stale-peer refusal, partial-update fail-closed), not the
  ceremony itself.
- Every W4 test rebuild changes the CDHash and must rebuild-and-repin at every step that changes any
  signed byte, per the explicit warning in the archived
  [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-installed-development-topology/RESULTS.md).

## Installation and rollback test matrix

Narrowed to what's reachable with only a free/Individual Apple Development identity and a single
owned Mac unless noted.

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

**Testable on this one Mac** once App IDs/profiles/entitlements exist: rows 1–13, 16, 17, 20 above
(all already demonstrated feasible on one host by P0-4A or Gate B, against a different team but a
host-local mechanism); the full XPC peer-validation matrix (§6) with real W4-signed binaries instead
of ad-hoc signing; real Keychain access-group separation and epoch-key rotation for W4-owned App
IDs; real `SMAppService.agent` lifecycle with a genuine Development-signed, sandboxed bundle (P0-4A's
sandboxed runner failed before `main` under ad-hoc signing with AMFI error `-423`; a real Development
profile should get past that specific failure).

**Needs a second machine** (not necessarily clean): multi-machine trust-epoch behavior; genuine
fast-user-switch/session concurrency under real load.

**Needs a genuinely clean host, cannot be done on this Mac at all:** Gatekeeper/quarantine/
translocation on first launch (`spctl --assess` on the build machine is not equivalent to a fresh
machine's first-launch path); supported macOS floor validation (P0-4A's effective floor was
measured at macOS 26.0, driven by a dependency, not `LSMinimumSystemVersion`); install/update races
across genuinely independent boot/session state.

## Claims that stay impossible without clean-host testing or a paid Developer ID

Per this plan's explicit scope (no paid Developer ID/notarization identity):

- Notarization and stapling.
- Gatekeeper assessment as an end user would experience it — a Development-signed, unnotarized app
  is supposed to fail Gatekeeper; P0-4A already confirmed that (`spctl` status 3).
- "Validated"/"production-ready" language of any kind for any of this — `AGENTS.md` and the Threat
  Model both forbid that without exact mechanism plus retained adversarial evidence, which this plan
  does not produce; it only sets up the ability to start collecting S5 evidence.
- Coherent-rollback resistance (needs a non-rollbackable anchor or external witness; neither exists
  in this repo's design yet).
- Minimum supported macOS floor, and cross-machine/clean-host install-race/update evidence (§9).
- Anything about the actual guest/backend (libkrun/HVF, Apple Containerization, gVisor) —
  orthogonal to this plan; P0-4A "cannot admit the backend" per `GATE_C_P0_RECONCILIATION.md`.

## Explicit action list for the user/orchestrator

Nothing below has been executed:

1. Select the user-confirmed Individual Team `W4QUR9FUL4` explicitly in the Apple portal/Xcode and
   verify the portal's current App ID/profile capacity.
2. Reissue/create an Apple Development certificate under that selected Team whose public subject OU and harmless
   signed-byte TeamIdentifier both read back exactly `W4QUR9FUL4`; revoke/retire the misleading
   certificate according to Apple account policy rather than relabeling it locally.
3. Decide Option A vs. Option B (§3) before reserving App IDs — it changes whether 3 or 4 distinct
   App IDs are needed.
4. In the Apple Developer portal, register the exact explicit App IDs from §4, enabling App Sandbox
   and Keychain Sharing where listed.
5. Register this Mac's device UDID as a development test device under the team.
6. Create one macOS App Development provisioning profile per App ID.
7. Download/install the profiles via Xcode — not by copying Gate B's cached `3DDR84M4JS` profiles,
   which are unusable here.
8. Author the four `.entitlements` files (§5/§6) using the archived
   [`Provisioned/Entitlements`](https://github.com/Shrimpworks/capsule-experiments/tree/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/macos-authority-separation/Provisioned/Entitlements)
   files as a structural template but with W4-scoped,
   epoch-suffixed group names.
9. Once real Xcode targets exist (they don't yet — only `cmd/capsuled/main.go` is real product
   source), wire `CODE_SIGN_ENTITLEMENTS` and `PRODUCT_BUNDLE_IDENTIFIER` per target, matching the
   `GateBProvisioned.xcodeproj` build-setting pattern already in the repo.
10. Build and locally run the §8 test matrix, rows 1–13/16/17/20 first, recording exact
   CDHash/entitlement-digest/Team-ID readbacks at every step per §7.
11. Only after S5 evidence is retained and reviewed should anyone consider requesting Developer
    ID/notarization for this team — a distinct, separately-authorized future task.
12. Whatever S5 evidence is produced must be retained in `capsule-experiments`, while this
    repository receives the resulting decision, an exact commit-pinned evidence link, this plan's
    S5 status update, and the `WORKSTREAM_EVIDENCE_LEDGER.md` update. Chat history alone is not
    retained evidence.

## Explicit non-goals

This plan does not authorize, and did not perform: creating any Apple ID/App ID/profile, touching
the Apple Developer portal, generating or downloading any provisioning profile, signing any binary,
creating any Keychain item or Secure Enclave key, or claiming any Capsule component is "validated,"
"secure," or "production-ready." It does not select between §3's Option A/B, does not invent a
fourth native-front-end component the repository hasn't chosen, and does not extend scope into
Developer ID/notarization, which is explicitly excluded.
