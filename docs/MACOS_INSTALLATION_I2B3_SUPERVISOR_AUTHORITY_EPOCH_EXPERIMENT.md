# macOS I2B3 Supervisor-authority-epoch inert packet and mutation matrix

Date: 2026-08-05

```text
Work item: passive I2B3 Supervisor-authority-epoch packet, C3a E0 construction, and C3b
  profile/signature preflight
Status: PASSED
Scope: bounded inert packet, reproducible unsigned fixtures, and separately authorized no-launch
  profile/signature readback for versioned Supervisor/Coordinator identities, container
  nonmembership, App Group/Mach separation, Keychain-group separation, and retirement oracles; no
  bundle launch, container, service, Keychain item, LocalAuthentication prompt, root, store,
  runtime, backend, VM, or guest operation was performed
Evidence or reason: Proposed ADR-0045 and this checked-in packet close the candidate identity tuple,
  evidence classes, stop order, mutation matrix, pass/refusal oracles, and claim boundary; immutable
  capsule-experiments merge dee784d40684100f8315720fab9a5cd3399f492b retains the exact unsigned
  sources, bundles, inputs, manifest, independent verifier, and mutation results, while merge
  ee00ae2abbce64ae6458b82d0b53d904ee39aeb6 retains the exact signed-profile preflight result
Remaining work: separately authorize the disposable-container C3b/E1 mutation run; later
  Keychain/service/root/store work requires another authorization after this gate passes
Next action: obtain a fresh authorization for E1-01..E1-12 and E1-14..E1-15 against the retained
  identities; E1-13 remains excluded, and work stops again before later Keychain/service/root/store
Parent status: installed owner-lock G3/I2B remains BLOCKED
```

## Defensive authorization boundary

This document authorizes no experiment. It freezes the next bounded defensive test of the
ADR-0045 stale-authority control using only repository-owned inert fixtures and, after separate
authorization, disposable Capsule Team `3DDR84M4JS` development probes on the owner's named Mac.
Do not access any other system, identity, credential, container, process, or data, and preserve all
existing Capsule safeguards.

The experiment is split at hard review boundaries:

1. **I2B3-E0 inert packet:** this checked-in closed definition of deterministic source, bundle
   layout, plist, entitlement, expected requirement, profile-request metadata, test vectors, and
   verifier. It performs no portal action, signing, profile download, install, registration,
   launch, Keychain operation, container creation, or filesystem mutation.
2. **I2B3 profile/signature gate:** separately authorized explicit App IDs/profiles and
   signing/readback only. It embeds and independently reads the exact development profiles and
   effective entitlements for the current Supervisor, never-launched Coordinator, and legacy
   negative probe. This gate is now `PASSED` at the immutable merge pinned below. No bundle was
   launched and no container was opened.
3. **I2B3-E1 identity-separation mutation:** separately authorized disposable probe launches
   against only platform-created test containers.
   It must stop and clean the exact disposable sentinels before any `SMAppService` registration,
   bootstrap Mach service, Coordinator launch, LocalAuthentication, Keychain item, installation
   root, Supervisor protected root, owner lock, store, runtime, backend, VM, guest, or attempt.

Passing E0 does not authorize E1. Passing E1 does not authorize the later Coordinator/key/service/
root corpus.

## Retained C3a construction

The exact E0 materialization is retained at
[`capsule-experiments` merge `dee784d40684100f8315720fab9a5cd3399f492b`](https://github.com/Shrimpworks/capsule-experiments/tree/dee784d40684100f8315720fab9a5cd3399f492b/experiments/macos-installation-i2b3-supervisor-authority-epoch-e0).
It contains reproducible unsigned current/legacy probe bundles, a Coordinator bundle that was
never launched, exact plist/entitlement/profile-request/LaunchAgent/descriptor inputs, a closed
manifest, independent verification, and 23 refusal mutations. It used no Apple identity, profile,
signing, container, service, Keychain, protected state, runtime, backend, VM, or guest. This closes
C3a construction only; it does not observe E1 nonmembership or accept ADR-0045.

## Retained C3b preflight reconciliation

The exact legacy I2B3 negative profile was restored through an owner-controlled non-Git path and
reconfirmed as UUID `c45a058b-ffdd-4a6b-bd8c-d746772a2702` with CMS SHA-256
`964f79980edf22a7280fe19e52893a1e40b0a8639d5bbe3d5dc8fdfada9c6c76`. Raw profile bytes remain
outside Git.

The zero-effect [initial E1 preflight](https://github.com/Shrimpworks/capsule-experiments/tree/50c494d4841c5d42e8e2120b82c0481a706a5236/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1)
records the earlier exact-profile absence. The later
[App Group portal preflight](https://github.com/Shrimpworks/capsule-experiments/tree/e6390253a274e9ead76366f9869a5e1b272a1595/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1-app-group-preflight)
observed that the Developer portal prepends `group.` when asked to register the frozen
`3DDR84M4JS...` identifier. The canonical
[correction](https://github.com/Shrimpworks/capsule-experiments/tree/3671a6eb23357ff28de4562dd60e8f68173034ae/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1-app-group-preflight)
narrows the resulting `NO_GO` to that portal-registration path. Apple's current documentation
supports the frozen `<team identifier>.<group name>` form on macOS without Developer-website App
Group registration. The exact identity therefore remains intended; its platform claim remains
`BLOCKED` on E1 container evidence. Do not silently rewrite it to an iOS-style `group.` value.

The exact no-launch
[signed-profile preflight](https://github.com/Shrimpworks/capsule-experiments/tree/ee00ae2abbce64ae6458b82d0b53d904ee39aeb6/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1-signed-profile-preflight)
is `PASSED` at immutable merge `ee00ae2abbce64ae6458b82d0b53d904ee39aeb6`. On owner Mac
`dsteele-shrimp-mbp18-4-01`, it retained the two exact explicit App IDs and Mac development-profile
UUID/CMS/certificate/device bindings plus strict signature, designated-requirement, CDHash,
hardened-runtime, absent-debug, App Sandbox, macOS-style App Group, and role-specific Keychain-group
readback for the current Supervisor, never-launched Coordinator, and legacy negative probe. Raw
profiles remain outside Git. No bundle or Coordinator was launched, and no container, sentinel,
service, Keychain item, LocalAuthentication prompt, root/store, runtime, backend, VM, or guest was
accessed.

## Frozen identities and packet contents

E0 uses authority sequence one and these exact candidate identifiers from Proposed ADR-0045:

- Supervisor: `com.capsulecorp.capsule.supervisor.authority-e1`;
- Coordinator: `com.capsulecorp.capsule.trust-bootstrap.authority-e1`;
- LaunchAgent label: `com.capsulecorp.capsule.supervisor.authority-e1`;
- bootstrap App Group: `3DDR84M4JS.com.capsulecorp.capsule.bootstrap.authority-e1`;
- bootstrap Mach service:
  `3DDR84M4JS.com.capsulecorp.capsule.bootstrap.authority-e1.supervisor`;
- Coordinator installation-root group:
  `3DDR84M4JS.com.capsulecorp.capsule.trust-bootstrap.installation-root.authority-e1`;
- Supervisor bootstrap-anchor group:
  `3DDR84M4JS.com.capsulecorp.capsule.supervisor.bootstrap-anchor.authority-e1`; and
- Supervisor evidence group:
  `3DDR84M4JS.com.capsulecorp.capsule.supervisor.evidence.authority-e1`.

The legacy negative identities are the exact retained I1B/I2B3 stable Supervisor signing ID
`com.capsulecorp.capsule.supervisor` and the unlaunched I2B3 Coordinator ID
`com.capsulecorp.capsule.trust-bootstrap.v1`. They are residue, not epoch zero. E0 must not copy,
open, enumerate, delete, normalize, or migrate their private-container contents.

The inert packet contains:

- one minimal fixed-byte current Supervisor probe and one byte-distinct legacy Supervisor probe;
- one minimal current Coordinator no-launch bundle for entitlement and requirement readback only;
- exact `Info.plist`, entitlements, LaunchAgent plist, inside-out bundle inventory, and inactive
  `SupervisorAuthorityDescriptorV0` known answer;
- exact Team/channel, signing identifier, profile UUID/digest, certificate fingerprint, CDHash,
  entitlement digest, App Group, Keychain-group, service, and bundle-relative program readback
  fields, with credential values left unresolved until an authorized run;
- a verifier that rejects missing, extra, mixed, stable-ID-substituted, sequence-zero, sequence-two,
  wrong-Team, wildcard-ID, wrong-group, wrong-service, path-bearing, or unknown-field packets; and
- fixed sentinel names and contents below. Probe source and bundles remain experiment-only and
  cannot be imported by product packages.

No descriptor contains a private key, provisioning profile bytes, device UDID, live container URL,
home-relative path, or secret. Public profile metadata may be retained later; raw profiles and
private credentials remain excluded from the repository.

## Evidence classes and source dates

| Input | Date | Class | Limitation |
| --- | --- | --- | --- |
| Apple App Sandbox file-access documentation | retrieved 2026-08-05 | documented mechanism | macOS 14+ code-signature/container association and foreign-container consent behavior; it does not prove the Capsule identities are distinct on the named host |
| Apple App Groups entitlement documentation | retrieved 2026-08-05 | documented mechanism | defines shared container, IPC naming, and possible Keychain group; it does not prove Capsule nonuse or residue |
| Apple Keychain access-group documentation | retrieved 2026-08-05 | documented mechanism | defines membership and exact-group queries for data-protection/synchronizable macOS items; it does not prove Secure Enclave/user-presence behavior for this packet |
| Apple `SMAppService` documentation | retrieved 2026-08-05 | documented mechanism | defines registration substrate only; E1 deliberately stops before registration |
| I2B3 stable-profile sentinel result | observed 2026-08-05 on macOS 26.5.2 build `25F84`, arm64 | observation | proves only stable-identity cross-profile mutation under Apple Development |
| Fresh explicit IDs yield separate private-container authority | proposed 2026-08-05 | inference to test | not a control until E1 positive and negative mutations pass without consent |
| Developer ID/notarized distribution behaves equivalently | none | unknown | outside E0/E1 and remains I6 `BLOCKED` |

## E1 ordering and immediate stops

After separate authorization, E1 must execute in this order and retain every command, exit status,
effective entitlement, code requirement, identity, platform-returned container URL digest, and
sentinel digest needed to reproduce the result:

1. Reconfirm the exact named host/OS/SDK, Team `3DDR84M4JS`, selected Apple Development certificate,
   device, already-created explicit App IDs/profiles, and absence of Developer ID use against
   immutable preflight merge `ee00ae2abbce64ae6458b82d0b53d904ee39aeb6`.
2. Do not create or register a macOS-style App Group in the Developer portal. Use only the retained
   exact epoch-one profile identities; stop if a profile projection rewrites them or otherwise
   differs.
3. Before launch, require strict inside-out signature, profile UUID/CMS/certificate/device binding,
   `com.apple.application-identifier`, Team OU/TeamIdentifier, signing ID, effective-entitlement,
   hardened-runtime, absent-debug, App Group, Keychain-group, LaunchAgent-label, and
   peer-requirement readback to match the retained gate exactly.
4. Under a fresh authorization, run only the two Supervisor container probes needed for the
   matrix. Each obtains its own
   platform-selected container/standard-directory URL through supported APIs. The harness records
   those returned URLs and supplies the peer URL only to the fixed cross-container test operation;
   no path is guessed, accepted by product code, or retained as authority.
5. Stop immediately on a foreign-container consent prompt, unexpected existing target contents,
   shared container URL identity, successful cross-open/link/rename/map/write, wrong denial class,
   legacy alias, or cleanup ambiguity. Do not grant consent in the passing run.
6. Remove only exact experiment-created sentinel entries through their owning probe, re-read
   absence, and inventory platform residue without deleting ambiguous metadata.
7. Confirm no service registered, no Coordinator launched, and no Keychain/root/owner/store object
   exists. Stop and publish the result for review.

The current probe may create only
`Library/Application Support/CapsuleAuthorityEpochProbe/current-e1-sentinel` with mode `0600` and
fixed `current-authority-e1` bytes. The legacy probe may create only
`Library/Application Support/CapsuleAuthorityEpochProbe/legacy-stable-sentinel` in its own legacy
container with fixed `legacy-stable` bytes. Neither creates a `supervisor.state`, owner, request,
record, ledger, database, socket, default, or key.

## Mutation matrix

All rows require exact pre/post byte and absence readback. A denial is `PASSED` only when current
state is unchanged and the observed refusal is classified; a crash, prompt, timeout, or missing
observation is `BLOCKED`, not a denial.

| Case | Actor and target | Expected oracle |
| --- | --- | --- |
| E1-01 | current probe creates/reads current sentinel through its platform URL | exact success; establishes only current self-access |
| E1-02 | legacy probe creates/reads legacy sentinel through its platform URL | exact success; confirms the negative probe is runnable without touching current state |
| E1-03 | compare platform-returned current and legacy container identities | distinct; no path is supplied or guessed |
| E1-04 | legacy probe attempts descriptor-relative open/read of current sentinel | OS denial with no consent; current bytes unchanged |
| E1-05 | legacy probe attempts descriptor-relative truncate/write of current sentinel | OS denial with no consent; current bytes unchanged |
| E1-06 | legacy probe attempts create/link/rename into current probe parent | OS denial with no consent; no current entry added or replaced |
| E1-07 | legacy probe attempts retained writable descriptor/mapping acquisition to current sentinel | OS denial; no descriptor/mapping retained; current bytes unchanged |
| E1-08 | current probe attempts symmetric read/write of legacy sentinel | OS denial with no consent; legacy bytes unchanged |
| E1-09 | current Supervisor bundle has only epoch-one App Group and role groups; legacy signed entitlement readback lacks all of them | exact structural nonmembership; no group operation yet |
| E1-10 | no-launch wrong-group Keychain query fixture for each role | inert verifier selects exact `kSecAttrAccessGroup`; live operation deferred |
| E1-11 | stable signing ID substituted into the epoch-one descriptor/profile fixture | verifier refusal before launch |
| E1-12 | epoch-one bundle with sequence zero/two, legacy service, or mixed Coordinator group | verifier refusal before launch |
| E1-13 | user grants a foreign-container prompt in a separately authorized negative run | current state is quarantined; result classified elevated-posture/availability, never nonmembership success |
| E1-14 | current owner removes exact current sentinel; legacy owner removes exact legacy sentinel | exact removal and parent absence; platform metadata inventoried, not normalized |
| E1-15 | final inventory of launchd, Keychain service names, install locations, roots, owners, stores | all absent; any presence fails the stop boundary |

E1 does not test Keychain membership dynamically. The later Keychain slice must separately prove
old-to-new, new-to-old, wrong-App-Group, omitted-`kSecAttrAccessGroup`, deletion, replacement-key,
fresh-user-presence, process-death, and cleanup behavior. It may begin only after E1 passes and a
new authorization names the exact groups, services, and keys.

## Pass, block, and candidate disposition

This E0 packet freeze, C3a materialization, and C3b profile/signature-only preflight are `PASSED`
because the retained definitions and immutable archives close the identities, fields, order,
matrix, oracles, unsigned bytes, profile/signature bindings, and independent verification without
bundle launch or container access. Those deterministic and signed identities are E1 inputs, not E1
platform evidence. E1 is `PASSED` only
when exact Apple Development identities on the
named host show distinct platform container associations, every current/legacy cross-mutation
denies without consent, cleanup is exact, and the stop inventory is empty.

E1 is `BLOCKED` on any missing profile, consent prompt in the normal run, indeterminate result,
cleanup ambiguity, platform/API mismatch, or unavailable authorization. The fresh-identity
candidate is `NO_GO` only if a determinate authorized run shows the legacy identity can mutate the
epoch-one private container without elevated consent, or if the required explicit identities
cannot be provisioned. The exact stable-identifier/profile-rotation candidate is already `NO_GO`
from I2B3; that disposition does not make E0/E1 or the installed parent `NO_GO`.

## Claim boundary after a pass

An E1 pass supports only an Apple Development, exact-host `spike-observed` claim that these exact
versioned identities separate their disposable App Sandbox private containers from the retained
legacy identity without user consent. It does not support:

- Coordinator session, Keychain, LocalAuthentication, Secure Enclave, installation-root,
  bootstrap service, protected root, owner lock, state engine, restart, repair, or retirement;
- baseline denial after the user grants foreign-container access, Full Disk Access, task-port, or
  another elevated capability;
- Developer ID, notarization, Gatekeeper, translocation, clean-host, minimum-OS, update/package,
  production distribution, or cross-version behavior; or
- runtime, backend, VM, guest, attempt, product admission, or an Accepted ADR.

## Primary sources

- [Apple: accessing files from the macOS App Sandbox](https://developer.apple.com/documentation/security/accessing-files-from-the-macos-app-sandbox)
- [Apple: migrating files to an App Sandbox container](https://developer.apple.com/documentation/security/migrating-your-app-s-files-to-its-app-sandbox-container)
- [Apple: App Groups entitlement](https://developer.apple.com/documentation/bundleresources/entitlements/com.apple.security.application-groups)
- [Apple: protecting local app data using containers](https://developer.apple.com/documentation/xcode/protecting-local-app-data-using-containers)
- [Apple: Keychain access groups](https://developer.apple.com/documentation/security/sharing-access-to-keychain-items-among-a-collection-of-apps)
- [Apple: `kSecAttrAccessGroup`](https://developer.apple.com/documentation/security/ksecattraccessgroup)
- [Apple: `SMAppService`](https://developer.apple.com/documentation/servicemanagement/smappservice)
- [Apple: code-signing requirements](https://developer.apple.com/documentation/technotes/tn3127-inside-code-signing-requirements)
- [I2B3 stale-profile blocker](MACOS_INSTALLATION_I2B3_SIGNING_PREFLIGHT_AND_STALE_PROFILE_BLOCKER.md)
- [Proposed ADR-0045](adr/0045-select-versioned-supervisor-authority-epochs.md)
