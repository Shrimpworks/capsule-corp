# Gate B Results: macOS Authority and Storage Separation

Date: 2026-07-31  
Authoritative repository baseline: `9bfd2ac` (`Document hardened architecture and spike plan (#7)`)  
Decision: **conditional-pass for feasibility; Gate B is not yet validated for a shipping configuration**

## Hypothesis and gate

Hypothesis: independently signed daemon, Broker, and Supervisor components can use macOS-enforced
peer identity, Keychain authority, hardware-backed keys, user-presence access control, and protected
containers to form real authorities. Trust epoch, plan, attempt, and migration semantics remain
Capsule protocol responsibilities.

Gate B requires the operating system and protocol jointly to deny unauthorized peers, key use, and
storage access. The local spike establishes that the required primitive families exist and exercises
several critical negative cases. It cannot establish the complete gate because the host has no Apple
Development, Developer ID, or Mac App Store signing identity and therefore cannot create validated
production entitlements or a representative installed three-component XPC topology.

## Environment and tools

Observed locally:

| Item | Value |
| --- | --- |
| Repository | exact detached commit `9bfd2ac`; worktree clean before spike |
| macOS | 26.5.2 (build 25F84), Darwin 25.5.0 |
| Hardware | arm64; Secure Enclave P-256 creation and signing observed |
| System Integrity Protection | enabled |
| Developer directory | `/Library/Developer/CommandLineTools` |
| macOS SDK | 26.5 |
| Apple clang | 21.0.0 (`clang-2100.1.1.101`) |
| Swift | 6.3.3; unusable for this spike because the SDK Swift interfaces were built by 6.3.2 |
| LLDB | 2100.0.17.203 |
| Signing identities | `0 valid identities found` for code signing |
| Go | 1.23.4 darwin/arm64 |
| Ambient Node | 16.15.0; unrelated to the native experiment and below repository requirements |

The Swift compiler/SDK mismatch was observed, not inferred. The probe was moved to Objective-C/C
and the public Security, LocalAuthentication, CoreFoundation, and XPC interfaces. No Swift result is
claimed.

## Retained prototype

Run:

```sh
./experiments/macos-authority-separation/run.sh --with-debugger
```

Derived binaries remain ignored under `build/`. The Keychain probe uses unique process-scoped tags,
deletes any persistent test items, suppresses interactive authentication, and retains no key
material. The final run exited successfully.

## Observed local evidence

### Code identity and stale/copy behavior

The experiment compiled distinct daemon/Broker fixtures and ad-hoc signed them with explicit
signing identifiers.

| Case | Observation |
| --- | --- |
| Correct Broker identifier | identifier requirement matched |
| Wrong daemon identifier | denied |
| Unsigned fixture | denied |
| Impostor ad-hoc signed with the expected identifier | **matched an identifier-only requirement** |
| Second/stale build with the same identifier | **matched an identifier-only requirement** |
| Apple-chain plus identifier requirement against ad-hoc code | denied |
| Exact Broker v1 code-directory hash | v1 matched; v2 denied |
| Exact copy of Broker v1 | matched both identifier and exact-hash requirements |
| `get-task-allow` present | denied by an explicit entitlement-absent requirement |

Interpretation:

- Signing identifier alone is not authority. A production requirement must include the accepted
  Apple-issued signing chain/team and component identifier.
- Team plus identifier is a release identity, not an exact-build identity. Exact build enforcement
  needs the enrolled code-directory hash set (including all accepted architecture/hash variants).
- A byte-for-byte copied signed binary retains the same code identity. This is expected; path and
  filename are not identity. If launch context matters, use launch/responsible/parent constraints
  and protected installation/storage rules rather than a path check.
- Trust epoch is not a code-signing fact. The XPC/channel check must be combined with typed epoch
  binding and fail-closed protocol state.

### Dynamic validity and debugging

The public Security framework reported both ordinary fixtures dynamically valid and not debugged.
LLDB launched the debug-entitled fixture and the process then reported:

```text
seccode.dynamic-valid=true
seccode.debugged=true
```

The debugged flag is therefore independent from dynamic validity. Validated posture must require
both a valid signature/dynamic state and absence of `get-task-allow`/debug attachment. Apple
documents the debugged status as sticky for the process lifetime, but this remains point-in-time
evidence at the check boundary.

### Keychain and Secure Enclave

| Case | Observation |
| --- | --- |
| Add item to an unentitled data-protection Keychain access group | `errSecMissingEntitlement` (`-34018`) |
| Persist Secure Enclave key from the unprovisioned fixture | denied with `-34018` |
| Create nonpersistent Secure Enclave P-256 evidence key | succeeded |
| Sign SHA-256 digest with that key without interaction | succeeded |
| Create nonpersistent Secure Enclave P-256 approval key with `userPresence` | succeeded |
| Sign with authentication interaction prohibited | denied with LocalAuthentication `-1004`, “User interaction required” |

This directly observes hardware support, noninteractive evidence signing, user-presence gating, and
the missing-entitlement denial. It does **not** establish persistent per-component access groups;
that requires a validated application identifier/keychain-access-groups entitlement and matching
provisioning profile.

The private-key export case is covered by Apple’s platform contract rather than this persistent
fixture: the unprovisioned process was denied before it could persist the key. Apple states that
Secure Enclave private key material cannot be transferred into or out of the enclave and supports
only generated P-256 keys.

### SDK/API availability observed in the macOS 26.5 SDK

- `xpc_connection_set_peer_code_signing_requirement`: macOS 12.0+; XPC checks all received
  messages and drops/cancels mismatching peers.
- team identity, entitlement, and lightweight peer requirements: macOS 14.4+.
- `xpc_peer_requirement_t`, `xpc_connection_set_peer_requirement`, and conversion from a
  `ProcessCodeRequirement` to an XPC peer requirement: macOS 26.0+.
- `SecCodeCreateWithXPCMessage`: creates a dynamic code object from the actual received XPC message,
  avoiding PID/path/name substitution.
- `SecCodeCheckValidityWithProcessRequirement`: macOS 15.0+; process constraints expose dynamic
  validity, debug capability, Hardened Runtime, library validation, ad-hoc status, and debugged
  status.
- XPC exposes peer effective UID and audit-session ID at connection establishment. The correct
  Capsule GUI/login-session policy and fast-user-switching cases were not exercised.

## Current primary-source platform evidence

### XPC and runtime identity

Apple’s XPC peer requirement is applied to messages received on the connection; mismatching
listener requests are dropped and reply paths receive a peer-code-signing error. The current SDK
also permits a process code requirement to become an XPC peer requirement on macOS 26. See
[xpc_connection_set_peer_requirement](https://developer.apple.com/documentation/xpc/xpc_connection_set_peer_requirement),
[TN3127: Inside Code Signing Requirements](https://developer.apple.com/documentation/Technotes/tn3127-inside-code-signing-requirements),
and [SecCodeCreateWithXPCMessage](https://developer.apple.com/documentation/security/seccodecreatewithxpcmessage%28_%3A_%3A_%3A%29).

Apple distinguishes signer/team identity from signing identifier and notes that ad-hoc signatures
cannot reliably preserve an identity across versions. Dynamic validity checks combine the running
code state with static signature validation. Debug attachment is separately observable through
[`kSecCodeStatusDebugged`](https://developer.apple.com/documentation/security/seccodestatus/debugged),
and [`SecCodeCheckValidity`](https://developer.apple.com/documentation/security/seccodecheckvalidity%28_%3A_%3A_%3A%29)
performs dynamic validation.

For launch authority, launch, parent, responsible-process, library, and launchd spawn constraints
are kernel/launchd mechanisms, not XPC authorization replacements. See
[Applying launch environment and library constraints](https://developer.apple.com/documentation/security/applying-launch-environment-and-library-constraints)
and [Defining launch environment and library constraints](https://developer.apple.com/documentation/security/defining-launch-environment-and-library-constraints).

### Keychain, Secure Enclave, and user presence

For macOS data-protection Keychain items, access groups are derived from validated entitlements.
An item belongs to one group; using an unauthorized group fails with `errSecMissingEntitlement`.
This applies on macOS when `kSecUseDataProtectionKeychain` is used, and Secure Enclave keys use that
model. See [Sharing access to keychain items among a collection of apps](https://developer.apple.com/documentation/security/sharing-access-to-keychain-items-among-a-collection-of-apps),
[`kSecAttrAccessGroup`](https://developer.apple.com/documentation/Security/kSecAttrAccessGroup), and
[`kSecUseDataProtectionKeychain`](https://developer.apple.com/documentation/security/ksecusedataprotectionkeychain).

Apple documents Secure Enclave support on M1-or-later Macs (and certain earlier Touch ID Macs), P-256
only, generated in the enclave, and not importable/exportable. User-presence access control may use
biometry or device credentials; it proves a system-gated key operation, not user comprehension or
correct Broker rendering. See [Protecting keys with the Secure Enclave](https://developer.apple.com/documentation/Security/protecting-keys-with-the-secure-enclave),
[`kSecAccessControlUserPresence`](https://developer.apple.com/documentation/security/secaccesscontrolcreateflags/userpresence),
and [Accessing Keychain Items with Face ID or Touch ID](https://developer.apple.com/documentation/localauthentication/accessing-keychain-items-with-face-id-or-touch-id).

Legacy macOS Keychain ACLs are not a substitute for this model: Apple documents `kSecAttrAccess` as
mutually exclusive with `kSecAttrAccessControl` and inapplicable to data-protection Keychain items.
See [`kSecAttrAccess`](https://developer.apple.com/documentation/security/ksecattraccess).

### Protected storage

Apple documents code-signature-associated app data containers on macOS 14+ and SIP-protected app
group containers for non-sandboxed bundled components on macOS 15+. Access by a nonmember triggers
user authorization rather than silently succeeding. See
[Accessing files from the macOS App Sandbox](https://developer.apple.com/documentation/security/accessing-files-from-the-macos-app-sandbox),
[Protecting local app data using containers on macOS](https://developer.apple.com/documentation/xcode/protecting-local-app-data-using-containers),
and [Accessing app group containers in your existing macOS app](https://developer.apple.com/documentation/xcode/accessing-app-group-containers).

Therefore protected containers are real OS enforcement, but the protection is conditional on
correct bundle signing/validated entitlements, SIP, and the trusted-user/administrator assumption.
A user can authorize another app, and a privileged administrator/kernel remains out of scope.

## Enforcement versus packaging/protocol assumptions

| Property | Classification | Boundary/condition |
| --- | --- | --- |
| Apple signer/team + component signing identifier on XPC messages | OS enforcement | Apple-issued dev/distribution signature; exact requirement installed before activation |
| Exact active build | OS enforcement when code-directory hashes are required | Manifest must contain all accepted hashes; update replaces accepted set |
| Dynamic validity and debugged status | OS observation, point-in-time | Check actual XPC-message sender; debugged is distinct from validity |
| Effective UID and audit session | OS observation | Capsule must define correct user/login-session policy and test switching |
| Trust epoch, registration, nonce, purpose, audience | Protocol enforcement | Not represented by code signing; typed messages and durable state required |
| Keychain access group separation by component | OS enforcement | Data-protection Keychain, validated entitlement/profile, no shared operational group |
| Exact build allowed to use a Keychain key | **Not provided by a stable access group** | Access-group membership is entitlement/team based, not code-directory-hash based |
| Secure Enclave nonexportable P-256 operations | Hardware/platform enforcement where available | Generated on supported hardware; no import; unavailable hardware must fail or use an explicit lower posture |
| Approval user presence | OS-gated key operation | Does not attest correct UI, comprehension, or a particular person remotely |
| Broker/Supervisor protected store | OS enforcement with user override/admin limitation | App data container or single-component app group; signed bundled process; SIP enabled |
| No broad shared app group | Packaging rule | A narrow IPC-only group may be unavoidable for sandboxed/nonsandboxed IPC but must contain no authority/state |
| Supervisor-only backend creation | Architecture plus IPC/launch enforcement | Exact backend/launch-helper topology still depends on Gates C/E |

## Required adversarial cases: disposition

| Required case | Result |
| --- | --- |
| Unsigned peer | static requirement denial observed; live XPC denial pending signed harness |
| Same-team wrong signing identifier | API/platform support substantiated; no local team identity, so live test pending |
| Stale build | identifier-only acceptance and exact-hash denial observed; epoch protocol pending |
| Copied binary / PID/path/name substitution | exact copy accepted by design; message-derived `SecCode` is the required identity mechanism |
| Debug/development signing | `get-task-allow` denial and live `debugged=true` observation passed; Apple dev vs distribution channel matrix pending |
| Wrong entitlement / Keychain group | unentitled access-group operation denied with `-34018`; positive provisioned group pending |
| Daemon use of Approval/Supervisor keys | negative entitlement path observed; actual three-bundle cross-group test pending |
| Broker/Supervisor store access | current Apple protection substantiated; actual signed containers and deny/cancel prompt cases pending |
| Migration/team change | documentation evidence only; no signing identities or alternate team available |
| Unavailable Secure Enclave | documentation evidence and test plan only; current hardware has a working enclave |
| User-presence unavailable/noninteractive | noninteractive use denied; not-enrolled, lockout, password fallback, cancel, and successful interactive cases pending |
| Wrong trust epoch | protocol design only; no Gate F lifecycle implementation exists |

## Decision and conditions

**Conditional-pass.** Current macOS provides credible enforcement primitives for the intended
authority split, and the local negative evidence found no fundamental platform blocker. The gate
must remain unvalidated/proposed until all conditions below pass on the minimum supported OS and
distribution channel:

1. Three distinct Apple Development and then distribution-signed installed components use unique
   signing identifiers and validated entitlements/profiles.
2. Every trusted XPC direction enforces signer/team, role identifier, accepted code-directory hash
   set, required/forbidden entitlements, effective UID/session, and actual message-sender dynamic
   state before parsing or state transition.
3. The active epoch is independently bound in protocol state; stale components fail even if their
   historical signatures remain valid.
4. Broker Approval, installation-root, and Supervisor evidence keys use disjoint data-protection
   Keychain groups. The daemon belongs to none of them.
5. Broker and Supervisor stores use distinct protected containers; a shared group, if required only
   for IPC naming, contains no authoritative state, content, keys, or trust checkpoints.
6. The project explicitly selects a minimum macOS path: macOS 26 XPC process requirements, or an
   earlier compatibility design based on legacy XPC code requirements plus
   `SecCodeCreateWithXPCMessage`/dynamic checks. Both paths need their own retained matrix.
7. Backend launch ownership is proven with the Gate C/E topology; Gate B alone cannot establish it.

## Concrete architecture/document changes proposed

No broad architecture files were edited in this spike. Apply these focused changes after review:

1. **`TRUST_ARCHITECTURE.md` / `RUNTIME_INTEGRITY.md`:** name
   `SecCodeCreateWithXPCMessage` as the compatibility-path source of peer identity and prohibit PID,
   path, or process-name lookup. Add a macOS-26 path using `XPCPeerRequirement` constructed from a
   `ProcessCodeRequirement`.
2. **`RUNTIME_INTEGRITY.md`:** make the preflight predicate explicit: signature valid, exact team,
   role identifier, active code-directory hash, Hardened Runtime/library-validation policy,
   `get-task-allow` absent, and debugged flag absent. Record checks as point-in-time unless a real
   monitor exists.
3. **`INSTALLATION_TRUST.md`:** distinguish release identity (team + identifier + channel) from
   exact build identity (code-directory hash set + entitlement digest). State that an exact copy is
   the same code identity and that epoch, not path, rejects stale state.
4. **`INSTALLATION_TRUST.md` / ADR-0012:** record that a stable Keychain access group is
   component-scoped, not build-scoped. Require operational key rotation on epoch change and resolve
   whether per-epoch access groups, a separate installation-authority component, or another
   reviewed mechanism prevents a stale same-team build from discovering/using a newly enrolled key.
5. **`TRUST_ARCHITECTURE.md` / ADR-0018:** move the installation root out of the routine Broker
   process. Prefer a rarely launched installation-authority/repair ceremony with its own identifier,
   Keychain group, and protected store.
6. **`ARCHITECTURE.md`:** define storage packaging: sandboxed app data container where compatible;
   otherwise a single-component SIP-protected app group container on macOS 15+. An IPC-only common
   group is non-authoritative and documented as a narrow exception.
7. **`COMPONENT_COMPROMISE_MATRIX.md`:** add the same-team stale-build/key-group residual. Exact
   XPC hash denial does not itself revoke Keychain group membership.
8. **`FEASIBILITY_SPIKES.md`:** require development, Developer ID, and any Mac App Store variants
   separately; a development-signing pass cannot promote distribution posture.

## Distribution and entitlement prerequisites

- Stable Apple Developer Team ID and distinct bundle/signing identifiers for daemon, Broker,
  Supervisor, updater, and any launcher.
- Apple Development identities for the full local negative harness; Developer ID Application plus
  Hardened Runtime and notarization for direct distribution, or App Sandbox for Mac App Store.
- A validated provisioning profile for each bundled executable that claims restricted Keychain
  access groups. Apple’s distribution guidance says each program using those groups needs the
  corresponding profile. See [Creating distribution-signed code for macOS](https://developer.apple.com/documentation/xcode/creating-distribution-signed-code-for-the-mac).
- App-like bundle packaging for daemon-style executables that require provisioning-profile-backed
  entitlements. See [Signing a daemon with a restricted entitlement](https://developer.apple.com/documentation/xcode/signing-a-daemon-with-a-restricted-entitlement).
- Distinct single-member storage groups and distinct keychain groups. Do not assume a macOS-style
  team-prefixed app group is also a Keychain access group; Apple documents that it is not.
- Explicit update/migration handling for Team ID/App ID prefix, distribution-channel, designated
  requirement, profile, entitlement, and container association changes.

## Open risks

1. **Stale same-team key use:** access groups do not encode an exact build. A retained old component
   with the same validated group entitlement may remain able to query group keys. Protocol key
   revocation alone is insufficient if the stale process can find the replacement private key.
2. **Container user override:** protected container access can be user-authorized; the UX and
   behavior under Full Disk Access, MDM, fast user switching, and unattended launch agents need
   adversarial testing.
3. **Interactive approval:** only noninteractive denial was exercised. Successful user presence,
   cancel, lockout, no enrollment, password fallback, session switching, and screen-lock behavior
   remain unknown for Capsule UX.
4. **Hardware absence/failure:** Intel without suitable hardware, virtual machines, disabled or
   failed token services, OS restore, and hardware replacement were not available locally.
5. **Actual XPC lifecycle:** launchd registration, listener/client symmetric requirements, stale
   service activation, reconnect/replay, malformed first message, and process death were not run.
6. **Backend authority:** the Supervisor/launcher privilege and Apple Container control endpoint
   remain gated by Gates C and E.
7. **Migration:** Secure Enclave keys do not migrate; access-group/team changes and app-container
   designated-requirement changes can strand state or prompt. Recovery must create a new
   installation identity when continuity cannot be proven.

## Next smallest test

Build one minimal Xcode workspace on a Mac with an Apple Development team:

- a sandboxed Go daemon embedded in a signed bundle or XPC/login-item host;
- a native Broker app with an Approval Keychain group and private app container;
- an unprivileged Supervisor app-like launch agent with its own evidence group and one-member
  protected app-group container;
- one echo-only XPC method in each required direction, with exact message-derived peer validation;
- v1/v2 and wrong-role targets signed by the same team, plus unsigned/ad-hoc/debug variants;
- cross-group key/store probes and an epoch mismatch field.

Run the full wrong-peer, stale-build, debug attach, key-use, container prompt-deny, reconnect, and
session matrix first under Apple Development signing. Then export the same topology under the chosen
distribution channel and repeat it. This is the smallest test that can turn the present
conditional-pass into a real Gate B decision.

## Verification

Passed on the environment above:

- `./experiments/macos-authority-separation/run.sh --with-debugger`
- `sh -n experiments/macos-authority-separation/run.sh`
- `pnpm install --frozen-lockfile` under Node 22.22.1 / pnpm 10.28.2
- `pnpm check`
- `pnpm lint`
- `pnpm test`
- `pnpm verify:schemas`
- `go test ./...`
- `go vet ./...`
- `go build ./...`
