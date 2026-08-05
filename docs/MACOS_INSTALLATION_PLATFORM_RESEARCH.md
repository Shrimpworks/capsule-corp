# macOS installation platform-semantics research

Date: 2026-08-04

```text
Work item: post-I1B Apple-platform installation research
Status: PASSED
Scope: current Apple primary documentation, public macOS 26.5 SDK declarations, current-host
  read-only filesystem/platform observations, and a bounded Sparkle mechanism comparison only;
  no signing, profile/App-ID mutation, Keychain item, service registration, installation,
  replacement, runtime, backend, guest, or network-facing Capsule behavior was exercised
Evidence or reason: the seven questions in the installation-plan research brief now have an
  explicit documented/inferred/observed/unknown classification, supported bundle and authority
  candidates, first-run/replacement sequences, failure tables, and smallest follow-on spikes
Remaining work: exact Apple-signed Coordinator/Supervisor profiles and session behavior,
  user-presence Keychain use, stale-signer fencing, installed bootstrap, replacement interruption,
  Developer ID/notarized distribution, and cross-version/clean-host evidence remain blocked
Next action: implement passive I2B1 request/record/observation fixtures; authorize the exact
  developer-signed I2B2 platform mutations separately before creating a key, App Group-enabled
  profile, service, protected root, or installed replacement fixture
Parent status: macOS product installation and distribution are IN_PROGRESS — TRENDING_GOOD;
  installed I2B, replacement I4, and distribution I6 remain BLOCKED
```

## Defensive scope and method

This is a research result, not installed security-control evidence. It defensively reviews only
Capsule's one-app, per-user installation candidate and the already selected I2A bootstrap shape.
It does not authorize use of Apple Development or Developer ID identities, portal changes,
Keychain mutation, service registration, app replacement, or access to another container, user,
host, credential, process, or data.

Claims use these evidence classes:

- **Documented** — an Apple public document or public SDK declaration states the behavior.
- **Observed** — a read-only fact from the owned macOS 26.5.2 arm64 host.
- **Inference** — a bounded design consequence of documented facts; it still needs installed
  evidence.
- **Unknown** — Apple does not publish the guarantee needed by Capsule. The dependent path stays
  `BLOCKED`.

No private API, deprecated custom sandbox profile, temporary Mach exception, privileged helper,
root daemon, or automatic state recreation is proposed.

## Decision table

| Question | Research result | Product consequence |
| --- | --- | --- |
| 1. One-bundle `SMAppService` layout and lifecycle | `PASSED` for the documented substrate. macOS 13+ supports signed apps embedding LaunchAgent plists under `Contents/Library/LaunchAgents` and executable bytes addressed by bundle-relative `BundleProgram`. Registration bootstraps the agent immediately and at later logins. Changed agent plist or executable bytes require re-registration; Apple recommends unregistering first. The asynchronous unregister completion is the only documented point after which re-registration is safe. | Replacement must stop callers, asynchronously unregister and await termination, replace the whole bundle, verify the new bundle, then re-register and re-check status. Approval persistence, coherent replacement, and crash recovery are not guaranteed by `SMAppService`; I4 stays `BLOCKED`. |
| 2. ADR-0038 Coordinator-to-Supervisor bootstrap | `IN_PROGRESS — TRENDING_GOOD`. The direct embedded Coordinator XPC service, per-user Supervisor LaunchAgent, App-Group-named Mach service, and OS-enforced peer code requirement all have public mechanisms. There is no container-bootstrap circularity if Supervisor registration and approval complete before the Coordinator is invoked. The exact signed profile/session composition remains untested. | Refine the order to register/approve Supervisor first, then invoke the Coordinator. Because the Coordinator needs the user's Keychain and interactive LocalAuthentication, its XPC service must test `JoinExistingSession=true`; the documented `false` default is not an admissible assumption. I2B stays `BLOCKED`. |
| 3. App Group requirement and residual authority | `PASSED` for the architecture choice and limitation. Apple directs sandboxed apps to use an App-Group-prefixed name for Mach/XPC IPC. A private embedded XPC service is only directly reachable by its containing main app; an anonymous listener endpoint itself needs an already authenticated transport. | The bootstrap-only App Group is the narrow supported direct route found. It necessarily grants both members a shared container, IPC namespaces, shared preferences capability, and a potential Keychain access-group namespace. Nonmembers are structurally excluded by entitlements; member nonuse can only be tested, not structurally removed. |
| 4. `NSUpdateSecurityPolicy` | `PASSED` for a bounded negative claim. The macOS Ventura behavior narrows which Team/signing identifiers or package Teams may modify an app bundle; the default otherwise permits same-Team writers. | It is an additional app-bundle writer policy, not filesystem authority, App Sandbox authority, peer authentication, container/Keychain protection, version or CDHash binding, coherent update, rollback, or stale-binary revocation. It cannot close I2B or I4 by itself. |
| 5. `/Applications` and `~/Applications` replacement | `IN_PROGRESS — TRENDING_GOOD`. `NSWorkspace.AuthorizationType.replaceFile` plus its authorized `FileManager.replaceItem` is a public user-authorized replacement mechanism on macOS 10.14+. `FileManager.replaceItem` requires source and destination on one volume. Current-host observations show both candidate directories on one APFS data volume, but that is not a support-matrix guarantee. | A small on-demand replacer remains the narrowest candidate. It must run from staged verified bytes, obtain user authorization before shutdown, replace only after all old-bundle processes and registrations are gone, and relaunch the new app for verification/re-registration. Crash durability, rollback, translocation, quarantine, Gatekeeper, and notarization remain `BLOCKED`. |
| 6. Coordinator key and user-presence authority | `BLOCKED` for the complete denial claim. Secure Enclave P-256 signing with `userPresence` plus `privateKeyUsage` and a fresh authentication context is publicly supported; role-specific Keychain/App Group entitlements can exclude the visible app, daemon, updater, replacer, and Supervisor. Static access-group membership does not distinguish an older still-valid Coordinator build with the same entitlement. | The current Supervisor can reject a stale Coordinator at XPC with a current-profile/CDHash requirement, and missing or changed Keychain/ledger state can fail to repair-required. The Keychain itself does not provide release-version revocation. I2B must retain stale-signer use/delete/refusal evidence and accept only safely fenced denial-of-service, or a later ADR must select a stronger rotation mechanism. |
| 7. One-host versus support-matrix evidence | `PASSED`. A minimal current-host matrix is now separated from clean-host, older-OS, and Developer ID/notarized distribution claims. | One owned host can advance I2B and a development-only I4 mechanic. It cannot establish first-install defaults, minimum OS, Gatekeeper/notarization, translocation, approval migration, or supported distribution. |

The research assignment is `PASSED` because it answered and bounded the questions. That does not
make any incomplete installed path `PASSED`.

## Primary-source facts and claim boundary

| Source | Fact retained | Evidence class and limit |
| --- | --- | --- |
| [`SMAppService`](https://developer.apple.com/documentation/servicemanagement/smappservice) and [`register()`](https://developer.apple.com/documentation/servicemanagement/smappservice/register%28%29) | macOS 13+ controls signed, in-bundle login items, agents, and daemons. A registered LaunchAgent is bootstrapped immediately and on subsequent logins. | Documented. It does not promise cross-replacement continuity. |
| [Current `SMAppService.h`](https://developer.apple.com/documentation/servicemanagement/smappservice) | If agent plist or executable bytes change, the service “must be re-registered or it may not launch.” Async unregister completes after the running process has been killed and is then safe to re-register. | Public SDK. The synchronous unregister explicitly does not wait for reap. |
| [Updating helper executables](https://developer.apple.com/documentation/servicemanagement/updating-helper-executables-from-earlier-versions-of-macos) | Agent plists live under `Contents/Library/LaunchAgents`; `BundleProgram` is relative to the containing bundle. Status may require user action in Login Items settings. | Documented. User approval after a changed release must be read back, not inferred. |
| [Creating XPC Services](https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/CreatingXPCServices.html) | An embedded XPC service is on-demand, independently sandboxed, may be killed when idle, and is “available only to the main application that contains it.” `JoinExistingSession=true` is specified for a service that needs the user's Keychain. | Archived Apple guide, still matched by current public interfaces. Exact current-host behavior needs an installed spike. |
| [App Groups entitlement](https://developer.apple.com/documentation/bundleresources/entitlements/com.apple.security.application-groups) | On macOS, App Groups support group-prefixed Mach/XPC names and grant shared-container and IPC access. App Groups can also participate in Keychain sharing. | Documented residual capability. Empty-by-construction is not structural absence. |
| [Keychain access-group sharing](https://developer.apple.com/documentation/security/sharing-access-to-keychain-items-among-a-collection-of-apps) | Access groups are formed from explicit Keychain groups, application identifier, and App Groups. An item belongs to one group; queries without a group can search all groups available to the caller. | Documented. Capsule must always specify the Coordinator-only group and negative-test the bootstrap group. |
| [Protecting local app data using containers](https://developer.apple.com/documentation/xcode/protecting-local-app-data-using-containers) | macOS 15+ protects App Group containers from nonmembers, while group members share them. | Documented. Protection from outsiders does not reduce authority shared by the two members. |
| [`NSXPCConnection` current SDK](https://developer.apple.com/documentation/foundation/nsxpcconnection) and [`xpc_connection_set_peer_code_signing_requirement`](https://developer.apple.com/documentation/xpc/xpc_connection_set_peer_code_signing_requirement) | Objective-C code-signing requirements are public on macOS 13+; the C API is public on macOS 12+. Mismatching listener peers are rejected by the OS before application dispatch. | Public SDK. Use these APIs for the candidate floor; do not depend on the different macOS 26-only `xpc_connection_set_peer_requirement` API. |
| [`NSUpdateSecurityPolicy`](https://developer.apple.com/documentation/bundleresources/information-property-list/nsupdatesecuritypolicy) and [WWDC22 privacy](https://developer.apple.com/videos/play/wwdc2022/10096/?time=272) | The policy maps allowed process Team/signing identifiers and package Teams. Ventura added the app-bundle management protection; same-Team updates are the default. | Documented. No version, CDHash, container, Keychain, transaction, or rollback field exists. |
| [`NSWorkspace.AuthorizationType.replaceFile`](https://developer.apple.com/documentation/appkit/nsworkspace/authorizationtype/replacefile) and [`FileManager.replaceItem`](https://developer.apple.com/documentation/foundation/filemanager/replaceitemat%28_%3Awithitemat%3Abackupitemname%3Aoptions%3A%29) | A user can authorize the replace-file operation; the authorized file manager can perform only the matching replacement primitive. Replacement requires both URLs on one volume. | Public SDK/documented. Apple does not document Capsule's required crash-durable whole-app transaction. |
| [Updating Mac Software](https://developer.apple.com/documentation/security/updating-mac-software) and [TN3126](https://developer.apple.com/documentation/technotes/tn3126-inside-code-signing-hashes) | Code-signing validation is lazy per page; changing code used by a running process can trigger a signing failure. Nested code and resources remain sealed. | Documented. Every old-bundle process must be gone before replacement. |
| [App Translocation notes](https://developer.apple.com/forums/thread/724969) | The exact translocation trigger is not documented and there is no supported detection or original-path recovery API. | Apple DTS guidance. Capsule must install to a declared location before registering or replacing services. |
| [Notarizing macOS software](https://developer.apple.com/documentation/security/notarizing-macos-software-before-distribution) | Developer ID distribution requires notarization for current ordinary distribution, and Gatekeeper consumes the resulting ticket. | Documented distribution fact. No Developer ID action was authorized or performed here. |
| [Secure Enclave key protection](https://developer.apple.com/documentation/security/protecting-keys-with-the-secure-enclave) | Supported hardware keeps P-256 private-key operations in the Secure Enclave; the app receives operation results, not plaintext private-key bytes. | Documented. Hardware availability and recovery limitations require explicit handling. |
| [`kSecUseAuthenticationContext`](https://developer.apple.com/documentation/security/ksecuseauthenticationcontext) and [Touch ID reuse duration](https://developer.apple.com/documentation/localauthentication/lacontext/touchidauthenticationallowablereuseduration) | Omitting a context creates a one-use context; a supplied authenticated context may be reused. Touch ID reuse defaults to zero. | Documented. Capsule should use a new nonreused context for each request and record signature. |
| [Sparkle sandboxing](https://sparkle-project.github.io/documentation/sandboxing/) | Stock sandboxed integration normally uses an Installer XPC service and temporary global Mach lookup exceptions; a source-built alternate service composition is possible. | Identified production reference, not Apple authority. The default temporary-exception path is `NO_GO` for Capsule. |

Quoted text is deliberately short. The linked documents and SDK headers, not the quotations, are
the authority.

## Current-host public SDK observation

Read-only observation date: 2026-08-04.

| Item | Observation |
| --- | --- |
| Host | macOS 26.5.2 build `25F84`, arm64 |
| Toolchain | Xcode 26.6 build `17F113`; macOS 26.5 SDK |
| `SMAppService` | class and registration lifecycle available from macOS 13.0 |
| `NSXPCConnection.setCodeSigningRequirement` | available from macOS 13.0 |
| `xpc_connection_set_peer_code_signing_requirement` | available from macOS 12.0 |
| `xpc_connection_set_peer_requirement` | different typed API, available only from macOS 26.0; not selected |
| `NSWorkspaceAuthorizationTypeReplaceFile` | available from macOS 10.14; authorized `FileManager` ignores replacement backup/options |
| Secure Enclave declaration | P-256 keys, generated in the enclave; import of pregenerated keys is unsupported |
| Current directory facts | `/Applications` is `root:admin` mode `0775`; `~/Applications` is current-user-owned mode `0700`; both report APFS data-volume device `16777232` |

Public header SHA-256 values make the observation reproducible without retaining SDK bytes:

| Header | SHA-256 |
| --- | --- |
| `ServiceManagement/SMAppService.h` | `859ec2a7c4471dc1ac1e36b4de67e76d0a62ce4e4d145c942fd80cc4c2d8f8ff` |
| `AppKit/NSWorkspace.h` | `be46224f91d11b7d0e25a0ed80dda5c1a9e52c404038a5a7ebfc25794befff30` |
| `Foundation/NSXPCConnection.h` | `d34e7696e475abf3accda0a82127a61f5a95f5cb6609cdf5f57bfdea61c5fd1b` |
| `usr/include/xpc/connection.h` | `fee99c719a44ff57f5f3bd9e3e2d2b5518d093a476c0df9d4fbcc6a358750d51` |
| `Security/SecAccessControl.h` | `a3c5c75c87992d21a7bb1ce3e4638440a2fd0246936d8bb3a102b08fc3dc21ed` |
| `Security/SecItem.h` | `1a5a731dc52e85cf7fb6ebe26d75aa467113c07b6ec3389c98e183d75c651eaf` |

These observations do not establish behavior on another host or OS release.

## Proposed supported bundle and process tree

```text
Capsule.app                                      visible signed containing app
└─ Contents
   ├─ MacOS/Capsule                              Broker/setup/status UI
   ├─ Library
   │  ├─ LaunchAgents
   │  │  ├─ com.capsulecorp.capsule.supervisor.plist
   │  │  └─ com.capsulecorp.capsule.daemon.plist
   │  └─ Helpers
   │     ├─ CapsuleSupervisor.app                App-Sandboxed per-user agent
   │     │  └─ Contents/MacOS/CapsuleSupervisor
   │     └─ CapsuleDaemon.app                    App-Sandboxed per-user agent
   │        └─ Contents/XPCServices
   │           └─ CapsuleSourceValidatorDaemon.xpc
   └─ XPCServices
      ├─ CapsuleTrustBootstrap.xpc               on-demand Coordinator
      └─ CapsuleSourceValidatorBroker.xpc
```

`BundleProgram` values are relative to `Capsule.app`. The Coordinator is a direct private XPC
service of the visible app. The Supervisor is not a private service of the Coordinator; it is an
`SMAppService` LaunchAgent exposing only the bootstrap App-Group-prefixed Mach service during I2.
That distinction avoids relying on cross-bundle private-XPC lookup.

The tree is a candidate, not an admission. Exact nested profiles, entitlements, `Info.plist`
values, launch constraints, library constraints, peer requirements, and signed placement remain
I2B evidence.

## Authority and entitlement matrix

| Role | Required candidate authority | Structurally absent | Residual or unresolved authority |
| --- | --- | --- | --- |
| Visible app | App Sandbox; private lookup of its embedded Coordinator; `SMAppService` registration UI | bootstrap App Group; installation-root Keychain group; Supervisor container; daemon credentials; network update/replacement authority | Can invoke Coordinator and request Supervisor registration; cannot choose bootstrap fields or sign bytes |
| Trust Coordinator | App Sandbox; `JoinExistingSession=true` candidate; exact Coordinator-only Keychain group; bootstrap App Group; LocalAuthentication/Security use; two fixed signing purposes | Approval/Supervisor keys; Supervisor container; network; runtime/backend/guest; update/replacement; arbitrary file access | Bootstrap App Group shared container/IPC/potential Keychain namespace; private Keychain group remains reachable by stale builds with the same admitted entitlement until a stronger rotation rule exists |
| Supervisor | App Sandbox; private application container; bootstrap App Group; exact two-message Mach service; current Coordinator peer requirement | installation-root private key; visible-app container; updater/replacer; network; user-content parsing; runtime/backend/guest in I2 | Bootstrap App Group residual; same-user process availability attacks; current-profile verification and repair-required fencing remain mandatory |
| Daemon | App Sandbox; daemon-private Source Validator service only | bootstrap App Group; Coordinator Keychain group; Supervisor container/store; Approval/install/update/backend authority | None added by I2 |
| Broker Source Validator | its already proven private XPC route only | bootstrap App Group; Coordinator Keychain group; Supervisor state | Platform-managed private container residue already recorded by R3 |
| Daemon Source Validator | its already proven daemon-private route only | bootstrap App Group; Coordinator Keychain group; Supervisor state | Platform-managed private container residue already recorded by R3 |
| Future updater/verifier | none in I2 | Coordinator Keychain group; bootstrap App Group; Supervisor state; replacement authority | I5-only and not selected |
| Future Bundle Replacer | one sealed predecessor/target replacement descriptor and user-authorized replace operation, if selected | Keychain groups; bootstrap App Group; network/target selection/trust decisions; Supervisor mutation; execution | On-disk bundle-write authority for one transaction; I4 ADR/evidence required |

Every Keychain query must name its exact access group. Relying on the caller's default or searching
all available groups would silently widen the Coordinator through its App Group membership.

## First-run sequence

```mermaid
sequenceDiagram
    actor User
    participant App as Capsule.app
    participant SM as SMAppService
    participant Coord as Trust Coordinator
    participant Sup as Supervisor
    participant KC as Keychain / Secure Enclave

    App->>App: Verify exact installed bundle; attempts disabled
    App->>SM: Read Supervisor status
    alt Not registered
        App->>SM: register Supervisor LaunchAgent
    end
    alt Approval required or denied
        App-->>User: Explain and open Login Items settings
        App->>App: Stop with no key, request, root, or store
    else Enabled
        App->>Coord: Invoke private bootstrap method
        Coord->>KC: Create/open exact Coordinator-group key with fresh user presence
        Coord->>Coord: Construct, sign, and retain exact request
        Coord->>Sup: prepare(request) over bootstrap App-Group service
        Sup->>Sup: Authenticate current Coordinator; journal then create/observe fixed disabled state
        Sup-->>Coord: Typed observation bound to live connection
        Coord->>Coord: Validate observation; construct, sign, and retain exact record
        Coord->>Sup: finalize(record)
        Sup->>Sup: Verify/reobserve/retain record and disabled anchor
        Sup-->>Coord: Fixed finalized receipt
        Coord-->>App: Closed setup status only
    end
```

Registration and approval precede Coordinator key creation. This removes an avoidable circularity
and prevents a denied LaunchAgent from leaving a plausible installation-root key or request.

There is no exactly-once guarantee in XPC or `SMAppService`. The request and record remain
create-once, exact-byte, durably replayable Capsule objects. XPC interruption is a transport event,
not authority to create new bytes.

## Whole-bundle replacement sequence candidate

```mermaid
sequenceDiagram
    actor User
    participant App as Current Capsule.app
    participant Repl as Staged Replacer
    participant SM as SMAppService
    participant FS as Authorized FileManager
    participant New as New Capsule.app

    App->>App: Verify staged complete bundle and sealed replacement descriptor
    App->>Repl: Launch exact staged replacer with descriptor identifier only
    Repl-->>User: Request replace-file authorization if required
    alt Authorization denied
        Repl-->>App: Fixed denial; current install unchanged
    else Authorization granted
        App->>App: Persist attempts-disabled replacement journal
        App->>SM: Async unregister daemon and Supervisor
        SM-->>App: Completion after processes killed
        App->>App: Verify zero old-bundle process and exit
        Repl->>FS: Same-volume whole-bundle replace
        Repl->>New: Launch new visible app
        New->>New: Verify exact bundle, protected state, owner, epoch, and journal
        New->>SM: Re-register exact agents and read status
        New->>New: Recover forward or remain repair-required; never auto-enable
    end
```

The replacer obtains authorization before service teardown, but does not replace until it observes
all old-bundle processes gone. It receives no paths or target bytes from the daemon. The actual
descriptor, staging ownership, authorization handoff, `replaceItem` result, directory durability,
backup, relaunch, and forward-recovery protocol require a separate I4 ADR and spike.

## Failure and interruption tables

### Registration and bootstrap

| Failure point | Required outcome |
| --- | --- |
| Supervisor status is `requiresApproval`, `notFound`, or indeterminate | no Coordinator invocation, key, request, root, or store; attempts disabled; show exact setup state |
| `register()` reports already registered | read back status and exact installed identity; never treat the error as proof of correct registration |
| Agent bytes/plist changed | async unregister, await completion, re-register from verified containing bundle |
| Coordinator cannot join the interactive user session or present LocalAuthentication | `BLOCKED`; create no fallback software key and do not move signing into the visible app |
| Coordinator dies before retaining request | no Supervisor mutation; a later ceremony starts from absence |
| Coordinator dies after retaining request | only the exact retained request may resume; fresh bytes refuse |
| Supervisor dies before durable request admission | no root; exact retry may begin only if all absence checks hold |
| Supervisor dies after request admission | exact replay resumes the journal; different request is repair-required |
| reply lost after observation or finalization | Coordinator/Supervisor return retained exact bytes; no second key, root, record, or epoch |
| wrong role, Team, identifier, CDHash/profile, entitlement, EUID, audit session, or debug state | OS peer requirement and application checks refuse before body copy |
| stale Coordinator can access or delete its old Keychain-group item | current Supervisor channel rejects stale profile; missing/changed key or ledger is repair-required/new-installation identity; installed proof remains `BLOCKED` |
| bootstrap App Group contains a Capsule-created file/default/socket/key | stop I2B; do not normalize or silently delete ambiguous material |

### Replacement

| Failure point | Required outcome |
| --- | --- |
| staged bundle invalid, mixed, quarantined unexpectedly, translocated, wrong predecessor, or on another volume | refuse before authorization or service stop |
| user denies replace authorization | current bundle and registrations remain unchanged |
| unregister fails or completion is not received | current bundle remains; attempts disabled; no replacement |
| any old-bundle process remains | no replacement |
| death before filesystem replacement | current bundle remains authoritative; replacement journal records pre-publish state |
| `replaceItem` fails | never combine trees; current or staged exact bundle must verify, otherwise repair-required |
| death after replacement before relaunch/re-register | next user launch verifies successor and resumes exact forward recovery; attempts remain disabled |
| new agent requires Login Items approval | no fallback service, root daemon, or stale agent; attempts remain disabled |
| stale same-Team replacer matches `NSUpdateSecurityPolicy` | policy alone is insufficient; descriptor/current-profile checks and a replacement ADR must refuse it |
| rollback requested after successor state migration | refuse unless a separately authorized coherent rollback contract exists |

## Mechanism comparison for replacement

| Candidate | Status | Fit and stop conditions |
| --- | --- | --- |
| Small Capsule on-demand replacer using `NSWorkspace` replace authorization | `IN_PROGRESS — TRENDING_GOOD` | Smallest authority surface found. It can be user-authorized without a permanent root helper and can use the public same-volume replacement primitive. It still needs an ADR, exact staged execution placement, closed descriptor, old-process quiescence, interruption/durability evidence, and forward recovery. |
| Stock Sparkle 2 sandboxed installer path | `NO_GO` in its default exact configuration | Sparkle documents temporary global Mach lookup exceptions for the ordinary sandboxed installer connection. Capsule prohibits that exception and cannot make a generic updater its trust authority. |
| Source-built/narrowed Sparkle alternative XPC composition | `BLOCKED` | Sparkle documents alternate connection/status services, but this changes its build/TCB and still adds download, archive, signature, installer, relaunch, and update-state behavior. A later dependency-policy spike would have to prove it smaller or safer than the custom candidate. |
| Installer package for ordinary user-driven install/replacement | `BLOCKED` for product update; available as a distribution mechanism | A signed package can perform user-authorized installation, and `NSUpdateSecurityPolicy.AllowPackages` can name package Teams only. Package authority is Team-wide and does not become Capsule installation-root or coherent-update authority. |
| Deprecated `AuthorizationExecuteWithPrivileges`, setuid, or permanent privileged helper | `NO_GO` | The API is deprecated, App Sandbox does not support Authorization Services privilege escalation, and Capsule's per-user design does not justify permanent root authority. |

## Smallest current-host spikes

Each spike requires its own explicit mutation authorization. Passing one does not authorize the
next.

| Spike | Exact scope | Retained evidence | Excluded |
| --- | --- | --- | --- |
| P1 passive I2B bundle/profile fixture | Unsigned eighth-role bundle tree, exact plist/entitlement/peer-requirement candidates, request/observation/record fixtures, field authority | deterministic bytes, closed identifiers, availability gates, wrong-role/version/extra-field refusal | signing, App IDs/profiles, Keychain, service, filesystem mutation |
| P2 signed session/reachability | Exact Apple Development Coordinator and updated Supervisor profiles; execution-disabled app in `~/Applications`; register/approve Supervisor; private app-to-Coordinator and App-Group Coordinator-to-Supervisor fixed probes | effective entitlements, `JoinExistingSession` behavior, peer requirements, containers, wrong-role/stale/mixed/refusal, exact cleanup | Keychain key, protected root/store, runtime/backend/guest |
| P3 disposable user-presence key | One experiment-owned Secure Enclave P-256 key in the exact Coordinator-only group; two fixed digests, a fresh context for each sign; process death/relaunch; wrong-role queries; exact deletion | public key/key ID, access-control description, prompt/cancel/death results, no private bytes, cleanup | Approval key, Supervisor store, arbitrary signing, update/replacement |
| P4 protected-root I2B | Exact request/record plus descriptor-relative fixed root/owner/no-guest store on the owned current host | full I2B fault, replay, stale, same-user, restart, cleanup, and no-guest matrices | product store, attempts, runtime/backend/guest |
| P5 development replacement mechanic | Two exact inert Apple Development bundles on the current APFS volume; `~/Applications` first, `/Applications` only with a separate user-authorized operation; stop/unregister/replace/relaunch/re-register | signed-byte identities, authorization outcome, old-process inventory, predecessor/successor/interruption matrix, cleanup | Developer ID, notarization, network update, rollback, runtime/backend/guest |
| P6 Developer ID distribution matrix | Separately authorized Developer ID/notarized DMG/app/package and clean-host/multi-OS runs | Gatekeeper, quarantine, translocation, `NSUpdateSecurityPolicy`, first approval, replacement, logout/login/restart | not authorized by this research task; deferred to I6 |

## One-host evidence matrix

| Claim | One current owned host | Clean host or multiple OS releases required |
| --- | --- | --- |
| Exact bundle layout, effective entitlements, profiles, nested seals, and current SDK availability | yes, after explicit signing/profile authorization | minimum-OS and distribution portability |
| `SMAppService` register/status/unregister/re-register and current-login restart | yes | pristine first-approval defaults, OS migration, managed settings, every login/logout/reboot variant |
| Private app-to-Coordinator and App-Group Coordinator-to-Supervisor reachability | yes | support-floor behavior and clean account/session state |
| `JoinExistingSession`, LocalAuthentication prompt/cancel, Secure Enclave P-256, process death | yes, after explicit Keychain/user-interaction authorization | non-Secure-Enclave hardware, no-biometry, managed policy, older releases |
| Current-profile peer refusal and stale local fixture refusal | yes | revocation, expired credentials, old OS code-signing behavior, independently retained stale installs |
| App Group container remains empty of Capsule files/defaults/sockets/keys | yes for the run | structural absence is impossible; long-duration and multi-version residue |
| `~/Applications` same-volume replacement mechanics | yes | other filesystems, home layouts, migration and support claims |
| `/Applications` user-authorized replacement mechanics | yes, only with separate user action | nonadmin accounts, managed Macs, unusual ownership, external volumes |
| Developer ID, notarization, Gatekeeper ticket, quarantine, translocation, DMG/package | no | yes; distribution-only and explicitly deferred |
| `NSUpdateSecurityPolicy` across stale/current signers and packages | development-only mechanics cannot support the product claim | yes, with authorized Developer ID/notarized releases and clean policy state |
| Minimum macOS version | no | yes, on every declared minimum/current combination |
| Crash-durable replacement and coherent rollback | bounded fault injection can inform it | clean recovery machines and the final selected filesystem/update mechanism; rollback also needs an authority decision |

## Claims that remain blocked

- exact `JoinExistingSession=true` Coordinator sandbox, user-session, Keychain, and UI behavior;
- exact Coordinator and App-Group-enabled Supervisor App IDs/profiles and effective entitlements;
- denial or safe fencing of every stale same-Team/same-access-group Coordinator operation;
- exact App Group container, defaults, IPC, Keychain-query, and residue behavior over updates;
- first-install and post-replacement Login Items approval behavior on a clean account and host;
- automatic preservation of `SMAppService` registration or approval across replacement;
- crash-durable, same-volume whole-app replacement, backup, forward recovery, or rollback;
- `NSUpdateSecurityPolicy` enforcement for the final Developer ID/notarized app, package, and
  replacer identities;
- Gatekeeper, quarantine, notarization, stapling, translocation, and package behavior;
- `/Applications` behavior for nonadmin, managed, external-volume, or unusual-ownership hosts;
- a minimum supported macOS version or any cross-version support claim; and
- production installation, update, repair, restore, uninstall, runtime, backend, guest, or attempt
  admission.

Missing evidence is `BLOCKED`, not `NO_GO`. Only the exact default Sparkle temporary-exception path
and the already rejected privileged/deprecated replacement paths are `NO_GO` here.
