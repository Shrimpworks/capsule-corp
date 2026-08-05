# Apple certificates, credentials, identifiers, entitlements, and Capsule keys

Date: 2026-08-04; reconciled 2026-08-05 after I1B/R3

Status: canonical practical inventory and setup guide. The documentation work is `PASSED` when
merged. The exact credentialed I1B/Source Validator R3 development experiment is `PASSED`; product
Source Validator work and installed I2B remain `BLOCKED`, while Developer ID distribution,
notarization, production release, and key ceremonies remain deferred exactly as identified below.
This guide changes no product behavior, signing state, entitlement file, provisioning profile,
credential, key, artifact byte, or admission status.

## Purpose, evidence labels, and hard boundaries

This guide answers five different questions that are easy to conflate:

1. What does Apple require to identify a development team, sign code, authorize entitlements, and
   distribute software?
2. Which identifiers and entitlements does Capsule's proposed macOS composition need?
3. Which secrets belong to Apple workflows and which belong to Capsule's own protocol?
4. Where may each secret live, and who may use it?
5. What exact evidence must exist before Capsule may call an installed or distributed path tested?

The labels used throughout are:

- **Apple fact** — documented by a current Apple primary source linked near the claim.
- **Capsule decision** — selected by a named Capsule document or ADR.
- **Observation** — retained evidence from an exact historical inspection or experiment.
- **Inference** — a conclusion drawn from facts or observations, not an Apple guarantee.
- **Unknown** — not established without a later authorized portal, signing, installation, or
  clean-host task.

This document is metadata-only. It does not authorize inspecting a Keychain, enumerating local
identities or profiles, opening an Apple account, generating or exporting a key, signing code,
notarizing, installing a service, or changing a portal resource. Those operations require a
separately authorized task with the exact selected team and targets.

## Team ID and signed-development inputs reconciled; later credentialed work remains gated

Apple Membership Details confirms `3DDR84M4JS` is Dylan's Individual Apple Developer Program Team
ID. `W4QUR9FUL4` in the Apple Development common name is a member/display suffix, not Capsule's
Team ID. Apple defines the Team ID as the unique 10-character identifier assigned to a program
membership; a certificate's friendly name is not that definition
([Apple Team ID glossary](https://developer.apple.com/help/glossary/team-id/)). Later I2B3 and any
other credentialed slice still stop until their exact role/profile inputs exist and their use is
separately authorized; the completed I1B/R3 experiment is not repeated as an open prerequisite.
Its public certificate/profile metadata, signed readback, installed matrix, and cleanup result are
retained in the
[pinned R3 archive](https://github.com/Shrimpworks/capsule-experiments/tree/0944ffd8cfd01ec23e4ae99138b0931d56804077/experiments/completed-compiled-artifact-payloads/payloads/capsule-corp/artifacts/macos-i1b-r3-signed-development-composition).

### Retained observations

| Material | Public metadata observed | What it proves | What it does not prove |
| --- | --- | --- | --- |
| Apple Development certificate | SHA-1 `1638CFBD9250A00B4DBD81AE8FD1C790B42F61E3`; display text ended in `W4QUR9FUL4`; X.509 subject OU and signed-byte `TeamIdentifier` were `3DDR84M4JS` | The exact certificate belongs to Team `3DDR84M4JS` at the code-signing boundary | Authorization for current use or an exact Capsule role profile |
| Selected Apple Development identity | SHA-1 `80A4969BCD1B3926020888094B9D812A283D3793`; exact I1B/R3 signed-byte readback emitted TeamIdentifier `3DDR84M4JS` | The identity/private-key pairing and selected Team matched for the exact authorized I1B/R3 experiment | Authorization for any later signing task, Developer ID use, or product admission |
| Development profiles | The earlier three Gate B/wildcard profiles were nonmatching; I1B/R3 later created and selected three exact Team-3DDR explicit-App-ID profiles for Broker, daemon, and Supervisor | Exact profile/team/application/certificate/device/entitlement metadata matched the completed I1B/R3 scope | A Coordinator/bootstrap profile, later epoch profile, distribution profile, or reusable product admission |
| Developer ID Application certificate | User-run discovery reported SHA-1 `AD70CEDCA605604676C2853A229AA4664AD3F750`; Team `3DDR84M4JS` | A local Developer ID Application identity for the selected Team was reported | Authorization for use or current notarization/distribution admission |
| Membership | Apple Membership Details showed Team ID `3DDR84M4JS`, Apple Developer Program, Individual | The selected account Team ID and account type | Accepted agreements, role state, App ID/profile eligibility, or credential authorization |

**Conclusion:** the display string `Apple Development: Dylan Steele (W4QUR9FUL4)` was misleading
when treated as a Team ID. Neither that display string nor the default designated requirement is
Team admission evidence. Require exact subject OU/profile Team metadata and emitted signed-byte
TeamIdentifier `3DDR84M4JS`.

### Exact decision and replacement inputs

The selected path is Team `3DDR84M4JS`. `W4QUR9FUL4` is retired only as a Team-ID interpretation;
it remains historical common-name evidence. Every identifier, profile, certificate selection,
installation manifest, peer requirement, and epoch input must now be created or verified under
3DDR. The selected Apple Development identity passed exact subject/profile/signed-byte readback
only for the authorized I1B/R3 experiment; availability never selects it implicitly for later
work.

Before any new credentialed build, the task must receive these inputs as public metadata or exact
local selections, never as secret bytes in Git. I1B/R3 fulfilled them only for its retained exact
experiment:

1. selected Team ID and Membership-details evidence;
2. the ADR-selected containing-app/daemon/Broker/Supervisor bundle topology and explicit bundle
   IDs, plus any newly authorized role such as the I2B3 Trust Coordinator;
3. final pairwise App Group decision for Supervisor IPC, or the supported private alternative;
4. Apple Development certificate fingerprint, serial, validity, subject OU, public-key algorithm,
   and confirmation that its private key is usable on the authorized signing host;
5. one matching Mac App Development profile per independently provisioned executable that needs
   profile-authorized capabilities, including UUID, expiration, application identifier,
   certificate fingerprints, device UDID membership, and entitlements digest;
6. exact effective entitlements, signing identifiers, CodeDirectory hashes, Hardened Runtime
   state, launch/library constraints, and embedded-service placement for the built bytes; and
7. for distribution only, a matching Developer ID Application identity, any required Developer ID
   profiles, and a notarization credential associated with the same selected team.

Mixed-team output is a refusal, not a repair opportunity.

## Plain-English Apple map

### Account, membership, team, and roles

An Apple Account is the login. Apple Developer Program membership supplies the team, Team ID,
Certificates/Identifiers/Profiles access, Developer ID, and notarization eligibility. Enrollment
as an individual makes that person the Account Holder. Users an individual adds in App Store
Connect are not members of the Apple Developer Program team and do not gain its certificate and
profile access ([Apple roles](https://developer.apple.com/help/account/access/roles),
[account overview](https://developer.apple.com/help/account/basics/about-your-developer-account)).

For Capsule's reported Individual membership this means Dylan is the operational bottleneck for
Developer ID creation, agreements, membership renewal, and developer-team resources. App Store
Connect collaborators can help with App Store Connect tasks but cannot become equivalent signing
team maintainers. An organization membership is the supported route if multiple maintainers must
hold Apple Developer team roles; changing membership or transferring identifiers is an account
decision outside this guide.

### Bundle ID, App ID, Services ID, and service names

- A **bundle ID** is the reverse-DNS `CFBundleIdentifier` that identifies one executable bundle.
- An **App ID** is the team-scoped identifier used in a provisioning profile. Capsule needs
  explicit App IDs for provisioned components because explicit IDs identify one app and allow the
  capability allow-list; wildcard IDs are insufficient for the intended role-specific capability
  posture ([register an App ID](https://developer.apple.com/help/account/identifiers/register-an-app-id)).
- A **Services ID** identifies a website using Apple web services such as Sign in with Apple or
  WeatherKit. Capsule's current local macOS architecture needs no Services ID. Do not create one
  merely because Capsule has XPC or Mach *service names*
  ([register a Services ID](https://developer.apple.com/help/account/identifiers/register-a-services-id)).
- An **XPC/Mach service name** is a local IPC name, not a portal Services ID. Capsule's proposed
  Supervisor names and private Source Validator identities are enrolled installation metadata and
  signed-code identities.
- An **App Group** is a signed entitlement with real shared-container and IPC capability. On
  macOS it can enable IPC between sandboxed processes; it must not be treated as peer
  authentication ([App Groups entitlement](https://developer.apple.com/documentation/BundleResources/Entitlements/com.apple.security.application-groups)).

Changing a capability invalidates affected profiles, which must be regenerated
([enable app capabilities](https://developer.apple.com/help/account/identifiers/enable-app-capabilities)).
Therefore identifiers and capability sets must be frozen before each credentialed evidence slice;
the completed I1B/R3 set does not silently authorize a later I2B3 or distribution set.

### Development certificate and development profile

An **Apple Development certificate** binds an Apple-issued signing certificate to a developer and
team. Its private key performs signing. A **Mac App Development provisioning profile** binds one
App ID, selected development certificate(s), selected registered Mac device(s), and allowed
entitlements. Apple states those prerequisites directly and notes that Xcode automatic signing can
manage development profiles
([create a development profile](https://developer.apple.com/help/account/provisioning-profiles/create-a-development-provisioning-profile)).

The certificate and profile are both required for Capsule's Apple Development signed installed
tests when profile-backed capabilities are exercised. The profile cannot repair a wrong-team
certificate, and the certificate cannot authorize capabilities missing from the profile. Ad-hoc
signing is useful for static local inspection but supplies no Apple team or provisioning evidence.

### Developer ID Application and Developer ID Installer

**Developer ID Application** signs apps, command-line tools, frameworks, dylibs, nested helpers,
and other code distributed outside the Mac App Store. Capsule needs it for a Developer ID signed
DMG/test distribution and production direct distribution. Developer ID software should be
notarized and stapled before release
([Developer ID](https://developer.apple.com/support/developer-id/),
[notarization overview](https://developer.apple.com/documentation/security/notarizing-macos-software-before-distribution)).

**Developer ID Installer** signs a flat `.pkg` installer package. It does not sign the app inside
the package. Capsule does **not** need this certificate for its current one-app-in-a-DMG direction.
It becomes necessary only if a later managed/MDM distribution decision selects a signed `.pkg`.
Apple distinguishes the two certificate types and restricts Developer ID creation to the Account
Holder ([create Developer ID certificates](https://developer.apple.com/help/account/certificates/create-developer-id-certificates/)).

Developer ID profiles are required only when the distributed app uses capabilities that require
one. The final export must determine that from the effective entitlement set; do not embed a
development profile in Developer ID output. Apple warns that an expired Developer ID profile can
prevent launch even when the app certificate was valid when signed
([Developer ID expiration behavior](https://developer.apple.com/support/developer-id/)).

### Platform controls are different controls

| Control | What it does | Capsule use and boundary |
| --- | --- | --- |
| App Sandbox | Limits a process to signed entitlements and user-granted access ([entitlement](https://developer.apple.com/documentation/bundleresources/entitlements/com.apple.security.app-sandbox)) | Intended per native role and launcher; does not replace VM isolation or revoke inherited descriptors |
| Hardened Runtime | Applies runtime code-hardening restrictions; notarization expects it ([common notarization issues](https://developer.apple.com/documentation/security/resolving-common-notarization-issues)) | Enable on shipping-posture code; refuse `get-task-allow` and unapproved exceptions |
| Keychain access groups | Allow specifically entitled same-developer targets to access selected Keychain items ([keychain sharing](https://developer.apple.com/documentation/xcode/configuring-keychain-sharing)) | Separate Broker and Supervisor epoch groups; daemon and Source Validators get none |
| App Groups | Grant shared-container and IPC capability, not identity authentication | At most two pairwise groups if ADR-0029's installed evidence admits them; never store Capsule keys or authority state there |
| Private XPC service | Embeds a role-private service in a containing bundle | The exact two Source Validator routes passed inactive-policy R3 reachability; active parser/product admission remains blocked |
| `SMAppService` | Registers helpers located in an app bundle; a LaunchAgent is bootstrapped per user and on later logins ([Service Management](https://developer.apple.com/documentation/servicemanagement/), [`register()`](https://developer.apple.com/documentation/servicemanagement/smappservice/register%28%29)) | Exact I1B per-user daemon/Supervisor registration passed with execution disabled; protected-state, replacement, and repair semantics remain unproved |
| Launch/library constraints | Restrict where executables may launch and which libraries may load ([Apple constraints](https://developer.apple.com/documentation/security/defining-launch-environment-and-library-constraints)) | Bind the Source Validator launchers/parser children and later Runner to reviewed signed closure; not a provisioning profile substitute |
| Secure Enclave | Hardware security processor available on Apple silicon and T2 Macs; can protect supported key operations ([Apple Platform Security](https://support.apple.com/guide/security/hardware-security-overview-secf020d1074/web)) | Preferred for nonexportable Capsule operational keys when the exact API/algorithm/host supports it; not used for exportable Apple code-signing identities |
| LocalAuthentication | Requests system-mediated user authentication and returns success/failure without biometric data ([LocalAuthentication](https://developer.apple.com/documentation/localauthentication)) | Gates every v0 Approval-key operation and rare installation-root ceremonies; it does not prove comprehension |
| Hypervisor entitlement | Authorizes Hypervisor APIs in the entitled process ([entitlement](https://developer.apple.com/documentation/bundleresources/entitlements/com.apple.security.hypervisor)) | Belongs only to the later per-attempt Runner/VMM if libkrun/HVF is admitted, never daemon, Broker, Supervisor, or parser |

Temporary sandbox exceptions, JIT, unsigned executable memory, disabled library validation, broad
file access, Full Disk Access, Accessibility, Automation, Endpoint Security, root helpers, and
shared authority stores are not baseline Capsule requirements. Any proposed exception needs a
separate exact mechanism, owner, compromise analysis, and evidence.

## Apple credentials versus Capsule protocol keys

These systems are independent. An Apple signature says who signed executable code under Apple's
PKI. A Capsule protocol signature authorizes or attests one narrow Capsule object. Neither can be
substituted for the other.

| Material | System | Custodian | Purpose | Private/exportable? | Current state |
| --- | --- | --- | --- | --- | --- |
| Apple Development identity | Apple | Authorized developer/signing host | Development signing for installed tests | Private key is normally Keychain-held and may be exportable | SHA-1 `80A4...D3793` passed exact readback/use for I1B/R3 and the separately authorized I2B3 signing preflight; later use still requires exact authorization |
| Mac App Development profile | Apple | Build/Xcode; public signed metadata | Authorize App ID, device, certificate, entitlements | No private key | Three exact I1B/R3 profiles passed; I2B3 exact Coordinator/Supervisor profiles and signed entitlements also passed, but the old Supervisor profile retained write access to the stable private container, blocking root execution |
| Developer ID Application identity | Apple | Release custodian | Sign direct-distribution code | High-value exportable private key unless separately hardware-backed | Team-3DDR identity present; use/distribution admission deferred and unauthorized |
| Developer ID Installer identity | Apple | Package release custodian | Sign selected `.pkg` outer package | High-value private key | Unnecessary unless `.pkg` selected |
| Notarization API key or app-specific password | Apple/App Store Connect | Release/notarization operator | Authenticate uploads and status/log reads | Secret; API `.p8` is downloadable once | Deferred |
| Installation-root key | Capsule protocol | On-demand Trust Coordinator ceremony selected by Proposed ADR-0038 | Sign the protected-root bootstrap request/record, enroll operational keys, and authorize later trust transitions | Prefer nonexportable Secure Enclave reference in a Coordinator-only group; never visible app, daemon, Supervisor, updater, or replacer | I2A owner/contract selected; passive and installed I2B evidence blocked |
| Approval key | Capsule protocol | Approval Broker | Sign `capsule.plan.approve` after fresh user presence | Prefer nonexportable, user-presence-gated | Not implemented |
| Supervisor evidence key | Capsule protocol | Execution Supervisor | Sign `capsule.execution.attest` objects | Prefer nonexportable, noninteractive, narrow access group | Not implemented |
| Optional content-attestation key | Capsule protocol | Content Broker, only if required | Attest content-custody claims | Separate purpose and access group | Deferred/not selected |
| Trust epoch | Capsule protocol state, not automatically a new key type | Supervisor plus authorized trust-transition ceremony | Bind one coherent component/key/policy/profile state | Signed state; fresh Broker/Supervisor keys/groups on identity-changing epochs per ADR-0021 | Designed, not product-implemented |
| TUF root keys | Capsule release provenance | Offline release custodians | Authorize top-level TUF role keys and root rotation | Offline, split custody; not installed operational keys | Designed, operations not selected |
| TUF targets/delegation keys | Capsule release provenance | Scoped release/profile maintainers | Sign delegated artifacts/profile metadata | Offline or tightly controlled release systems | Designed, operations not selected |
| TUF snapshot/timestamp keys | Capsule release provenance | Release repository automation with bounded role | Bind repository state and freshness | Online only if later admitted; narrow scope | Designed, operations not selected |
| Package/runtime provenance signature | Capsule supply chain | Release/review system | Bind governed runtime/package bytes and review record | Separate from Apple signing and TUF role authorization | Required design; exact operations open |

The installation root, Approval key, Supervisor evidence key, and TUF keys are not Apple
certificates. The daemon must possess none of their private material. TUF roots, not URLs, DIDs, or
Apple certificates, anchor Capsule release/profile trust. See
[trust architecture](security/TRUST_ARCHITECTURE.md),
[installation trust](security/INSTALLATION_TRUST.md), and
[ADR-0014](adr/0014-tuf-trust-and-optional-dids.md).

## Exact environment matrix

| Environment | Required materials | Explicitly unnecessary or prohibited |
| --- | --- | --- |
| Unsigned/local conformance | Source tree; pinned ordinary build dependencies; fixtures and known-answer bytes; no product service activation; `FakeBackend.CreatesGuest() == false` where applicable | Apple account, Team ID decision, certificate, profile, Keychain identity, portal access, entitlements admission, notarization, Developer ID, CI signing secrets, Capsule production private keys, installed-service claims |
| Ad-hoc local artifact inspection | Exact built bytes; ad-hoc signature where Mach-O tooling requires one; static `codesign`/bundle/entitlement/dependency inspection; generated local fingerprints | Apple certificate/profile, team/channel enrollment, Keychain group or protected-container claim, Gatekeeper/notarization claim, Developer ID, production key ceremony, portal/device registration |
| Apple Development signed local installed testing | Selected Team ID; exact explicit App IDs; matching Apple Development identity and private key on authorized host; matching Mac App Development profiles and registered test Mac; exact role entitlements; Hardened Runtime shipping-posture variant; installed bundle/service manifest; fresh test-only Capsule keys only where the exact test requires them | Developer ID Application/Installer, notarization credential, stapling, DMG distribution, production TUF/release keys, CI, clean-host distribution claim; no historical other-team certificate/profile |
| Developer ID signed/notarized test distribution | All exact release bytes and nested-code signatures; selected-team Developer ID Application identity; secure timestamp; Hardened Runtime; Developer ID profiles only where capabilities require them; notarization credential; accepted log; stapled artifact where supported; Gatekeeper test on an authorized test host | Apple Development certificate/profile in release output; Developer ID Installer unless a `.pkg` is tested; production release/TUF root ceremony; automatic updater; production admission from one-host evidence |
| Production release/update | Reviewed release manifest and reproducible exact bytes; all nested Developer ID Application signatures; any required Developer ID profiles; notarization and staple evidence; final distribution container signature; pinned TUF roots and role-separated release metadata; governed runtime/package provenance; prepared-update/repair evidence; support-matrix and clean-host tests; controlled release custody and audit record | Development/ad-hoc identities; debug entitlements; mixed teams/epochs; ambient downloads; daemon-held release, root, approval, or evidence keys; `.pkg`/Developer ID Installer if package distribution remains unselected; unsupported rollback or production-readiness claims |

CI signing is intentionally later and optional. Local unsigned/ad-hoc work and the first authorized
installed tests do not need CI signing. Moving Developer ID material into CI expands the signing
attack surface and requires a separate release-workflow threat model, protected environment,
review gate, ephemeral Keychain, pinned actions, artifact readback, and secret-destruction evidence.

## Component, identity, profile, and nested-code matrix

All identifiers below are proposed inputs from the current installation and Source Validator
plans, not evidence that the App IDs or targets exist in Apple's portal. Final bundle topology and
the Supervisor App Group/private-service choice remain unresolved.

| Component | Proposed identity | Packaging/signing unit | Apple provisioning need | Key/capability posture |
| --- | --- | --- | --- | --- |
| Containing visible Swift app / Approval Broker | `com.capsulecorp.capsule.broker` (final visible-app/container identity still needs one explicit topology decision) | Outermost `Capsule.app`; independently signed executable and provisioned app target | Explicit App ID and dev profile for installed tests; Developer ID profile only if final capabilities require it | App Sandbox, Hardened Runtime, user-selected read-only files, Broker-only epoch Keychain group, Approval key; no backend or Supervisor store |
| Agent-facing daemon | `com.capsulecorp.capsule.daemon` | Separately enrolled nested executable/service inside `Capsule.app`; signed nested code | Explicit App ID/profile if profile-backed sandbox/IPC capability requires it | App Sandbox; no Keychain group, user-file entitlement, hypervisor, approval, evidence, updater, or replacer authority |
| Execution Supervisor | `com.capsulecorp.capsule.supervisor` | Native-fronted Go executable registered as a per-user `SMAppService` agent; signed nested code | Explicit App ID/profile for App Sandbox, Supervisor Keychain group, and any admitted pairwise IPC group | Supervisor-only epoch Keychain group and evidence key; two closed peer services; no user-file, network-update, parser, or hypervisor authority |
| Daemon Source Validator launcher | `com.capsulecorp.capsule.source-validator.daemon.v1` | Private XPC service embedded only in daemon role; separately signed nested bundle | Its containing/profile composition must be proved in R3; do not assume an independent portal profile until Xcode/effective entitlement evidence requires one | Private App Sandbox container as scratch only; no App/Keychain group, global Mach name, network, Capsule key, store, runtime, or backend |
| Daemon parser child | `com.capsulecorp.capsule.source-validator-parser.daemon.v1` | Exact embedded executable launched only by daemon launcher; separately signed nested code, not an independently provisioned service | Normally no independent profile unless an entitlement requires one; no portal App ID merely for a signing identifier | Role-specific launch/library/parent constraints; empty environment and closed descriptors; no persistent authority |
| Broker Source Validator launcher | `com.capsulecorp.capsule.source-validator.approval-broker.v1` | Private XPC service embedded only in Broker app; separately signed nested bundle | Same R3 evidence rule as daemon launcher | Separate private scratch container; no Approval-key use and no cross-role route |
| Broker parser child | `com.capsulecorp.capsule.source-validator-parser.approval-broker.v1` | Exact embedded executable launched only by Broker launcher; separately signed nested code | Normally no independent profile unless required by entitlement | Broker-role-specific constraints; cannot accept daemon requests/results |
| Optional update verifier | Not selected | Future separately enrolled executable or agent; if executable, independently signed nested code | Exact App ID/profile only after ADR and capability selection | May fetch/verify pinned TUF metadata and emit bounded `PreparedUpdate`; no replacement, installation root, Supervisor mutation, or execution |
| Trust/bootstrap coordinator | `com.capsulecorp.capsule.trust-bootstrap.v1` proposed by ADR-0038 | On-demand private XPC service embedded in `Capsule.app`; separately signed nested bundle | Exact App ID/profile, Coordinator-only installation-root Keychain group, and bootstrap-only Coordinator/Supervisor App Group require I2B evidence | Rare user-authorized request/record signing only; no visible-app key access, background registration, network, Supervisor state, updater/replacer, or execution authority |
| Optional bundle replacer | Not selected | Future mechanical helper; independent signed identity if selected | Requires ADR; Developer ID Installer is still unnecessary unless distribution is a `.pkg` | Installs one pre-authorized exact bundle; no target selection, network, trust, repair, or execution |
| Governed runtime bundle, kernel, firmware, runtime root | Closed release artifacts | Data/resource bundles unless a contained file is executable code; sign any Mach-O as nested code and bind all bytes in release/TUF manifests | No independent App ID/profile for inert data | No ambient download or mutation; exact provenance, hash, version, license, review, and runtime-profile binding |
| libkrun and native libraries | Exact dylibs/frameworks in reviewed closure | Nested code signed inside-out before containing executable/app; libraries do not receive app entitlements or profiles | No independent profile for libraries | Hardened library validation and exact closure; not an authority process by themselves |
| Per-attempt Runner/VMM | Proposed, not yet admitted | Independently signed executable nested in app and launched only from sealed Supervisor descriptor | Explicit identity/profile if Hypervisor/App Sandbox entitlements require it | Sole holder of `com.apple.security.hypervisor`; no registration, approval, update, root, or arbitrary path authority |

Apple recommends signing nested code from the deepest item outward and signing the main app last;
frameworks should not be given app entitlements or provisioning profiles
([latest code-signature format](https://developer.apple.com/documentation/xcode/using-the-latest-code-signature-format)).

### Entitlement invariants by role

- All shipping-posture executable roles: Hardened Runtime; no `get-task-allow`; no JIT, unsigned
  executable memory, or disabled library validation unless an exact later admission proves it.
- Broker: App Sandbox, user-selected read-only file access, one Broker-only epoch Keychain group.
- Supervisor: App Sandbox, one Supervisor-only epoch Keychain group, and only the pairwise IPC
  groups later admitted by ADR-0029 evidence.
- Daemon: App Sandbox and only the exact IPC capability required to reach the Supervisor; no
  Keychain, content, backend, or updater/replacer group.
- Source Validator launchers and children: role-private App Sandbox/launch constraints selected by
  ADR-0036; no shared group, global service, or temporary Mach exception.
- Runner/VMM only: Hypervisor entitlement if the libkrun/HVF path is admitted.
- App Groups, if selected: exactly daemon/Supervisor and Broker/Supervisor; never store files,
  defaults, sockets, Keychain items, manifests, source, content, state, or migration material in
  their shared containers.

## Setup without exposing secrets

Every example uses placeholders. Replace them only in an authorized local session. Never paste a
private key, password, token, or `.p12` into a command, terminal transcript, issue, PR, or shell
history.

### 1. Confirm membership and team in the portal

1. Sign in interactively to Apple Developer.
2. Open Membership details and record only the public legal entity/account type and
   `SELECTED_TEAM_ID`.
3. Confirm membership is active and agreements are accepted.
4. If the page does not show the selected ID, stop. Do not use Xcode display labels to override it.
5. For an Individual membership, record that only the Account Holder has Apple Developer team
   resource authority; App Store Connect users are not substitute team members.

This step uses the Apple portal and the paid Apple Developer Program membership. It creates no
certificate and retrieves no secret.

### 2. Freeze identifiers and capabilities

1. Resolve whether the visible Swift Broker is also the containing app identity and freeze the
   complete bundle tree.
2. In Certificates, Identifiers & Profiles, register each required explicit App ID under
   `SELECTED_TEAM_ID` using the final bundle ID.
3. Enable only the capabilities in the reviewed role matrix.
4. Do not create a Services ID.
5. Do not create App Groups until the ADR-0029 private-service/App Group residual-authority
   decision is explicit. If groups are selected, create exactly the two pairwise identifiers.
6. Export or screenshot no portal secrets. Record only the redacted metadata inventory below.

### 3. Obtain the Apple Development identity

Preferred developer workflow: in Xcode, select the exact account/team in Settings > Accounts >
Manage Certificates and create an Apple Development certificate. Xcode creates the private key in
the current user's login Keychain. Manual portal workflow: create a CSR with Keychain Access, upload
the CSR, download the public `.cer`, and install it. Apple's CSR workflow intentionally creates the
key locally ([create a CSR](https://developer.apple.com/help/account/certificates/create-a-certificate-signing-request)).

Immediately verify public metadata and refuse if subject OU is not `SELECTED_TEAM_ID`. Do not
export the identity merely to inspect it. A future authorized task may use read-only commands such
as:

```sh
security find-identity -v -p codesigning
security find-certificate -c "Apple Development" -p > "${TMPDIR}/capsule-apple-development.cer"
openssl x509 -in "${TMPDIR}/capsule-apple-development.cer" -noout -subject -issuer -serial -dates -fingerprint -sha256
```

The temporary file contains only the public certificate, but it still belongs outside the
repository and should be removed when the authorized inspection ends. Select by exact SHA-256 or
SHA-1 after discovery; never select a production identity only by common-name substring.

### 4. Register the test Mac and create development profiles

For manual signing, register the authorized test Mac's Provisioning UDID, then create one Mac App
Development profile per required explicit App ID. Select only the matching Apple Development
certificate and intended Mac. Xcode automatic signing may manage these profiles once real targets
exist ([registered-device distribution](https://developer.apple.com/documentation/xcode/distributing-your-app-to-registered-devices)).

For each downloaded `.provisionprofile`, inspect signed public metadata without installing it:

```sh
security cms -D -i "/path/to/PROFILE.provisionprofile" > "${TMPDIR}/capsule-profile.plist"
plutil -p "${TMPDIR}/capsule-profile.plist"
plutil -extract Entitlements xml1 -o - "${TMPDIR}/capsule-profile.plist"
```

Record the UUID, name, dates, TeamIdentifier, application identifier, developer-certificate
fingerprints, device membership, and an SHA-256 digest of the canonicalized entitlement plist.
Profile metadata can expose registered device identifiers and developer names, so it is not a
repository artifact even though it contains no signing private key.

When an authorized installed-test task needs to compare the host's installed provisioning state,
this read-only command emits profile metadata but no private key:

```sh
profiles show -type provisioning -output stdout-xml
```

It enumerates installed provisioning metadata, so do not run it during a narrower certificate-only
inspection or retain its raw output as a repository artifact.

### 5. Configure Xcode targets

For every signable target:

1. set the exact bundle/signing identifier from the frozen tree;
2. select `SELECTED_TEAM_ID` explicitly;
3. enable Hardened Runtime and only reviewed capabilities;
4. keep per-role entitlement files distinct;
5. set `get-task-allow` only in an explicitly labeled debug build, never a shipping-posture test;
6. embed services, parser children, libraries, runtime, and resources at their fixed bundle paths;
7. archive with reproducible inputs, then sign nested code inside-out; and
8. record exact CodeDirectory hashes and effective entitlements after every signed-byte change.

Read-only artifact checks:

```sh
codesign --verify --deep --strict --verbose=4 "/path/to/Capsule.app"
codesign --display --verbose=4 "/path/to/Capsule.app"
codesign --display --entitlements :- "/path/to/Capsule.app"
codesign --display --requirements :- "/path/to/Capsule.app"
spctl --assess --type execute --verbose=4 "/path/to/Capsule.app"
```

`--deep` is a verification convenience, not a signing strategy. Inspect and verify each enrolled
nested component independently as well.

### 6. Obtain Developer ID only when distribution begins

The Account Holder creates a selected-team Developer ID Application certificate in the portal or
Xcode. Verify its subject OU, TeamIdentifier, validity, fingerprint, and usable private-key binding
on the dedicated release host. Create Developer ID profiles only for targets whose final
capabilities require them. Never reuse a development profile in a Developer ID export.

Create Developer ID Installer only after an ADR selects `.pkg` distribution. A DMG containing
`Capsule.app` does not require it.

### 7. Choose notarization authentication

Apple supports two `notarytool` credential families:

- App Store Connect API key: key ID, issuer ID, and private `.p8`; the same key family can
  authenticate the Notary API
  ([App Store Connect API keys](https://developer.apple.com/help/app-store-connect/get-started/app-store-connect-api),
  [Notary API](https://developer.apple.com/documentation/NotaryAPI/submitting-software-for-notarization-over-the-web)).
- Apple Account plus app-specific password and Team ID. Two-factor authentication is required;
  app-specific passwords can be individually revoked and all are revoked when the primary account
  password changes ([Apple app-specific passwords](https://support.apple.com/en-us/102654)).

**Recommended staged Capsule path:** for the first manually operated Developer ID test, use an
interactive `notarytool store-credentials` prompt on the dedicated release Mac and keep the saved
profile in that user's login Keychain. Prefer a dedicated least-privilege **team** App Store
Connect API key for repeatable release automation after API access and notary compatibility are
verified; keep its `.p8` in offline release custody and inject it only for the notarization job.
Use an app-specific password as the simpler Individual-account fallback, not as a CI default. Do
not use the Apple Account's primary password. Do not assume an individual API key works with
`notarytool`; verify tool support before selecting it.

Run the interactive command without secret-valued command-line options so the tool prompts rather
than writing a password or key path into shell history:

```sh
xcrun notarytool store-credentials "capsule-notary-SELECTED_TEAM_ID"
```

Later submission and readback use only the Keychain profile name:

```sh
xcrun notarytool submit "/path/to/Capsule.zip" --keychain-profile "capsule-notary-SELECTED_TEAM_ID" --wait
xcrun notarytool history --keychain-profile "capsule-notary-SELECTED_TEAM_ID"
xcrun notarytool log "SUBMISSION_ID" --keychain-profile "capsule-notary-SELECTED_TEAM_ID" "${TMPDIR}/capsule-notary-log.json"
xcrun stapler staple "/path/to/Capsule.app"
xcrun stapler validate "/path/to/Capsule.app"
```

Always retain the submission ID and reviewed log, even on success. Do not commit a log until it has
been reviewed for paths, account names, identifiers, or other unintended metadata. Notarization is
an automated malware and signing check, not App Review or proof of Capsule correctness
([custom notarization workflow](https://developer.apple.com/documentation/security/customizing-the-notarization-workflow)).

### 8. Keep CI signing deferred until admitted

If a later ADR/workflow selects GitHub Actions signing:

1. use a protected `release-signing` environment whose secrets are unavailable until required
   review passes; GitHub documents that environment secrets are released only to jobs referencing
   the environment and, when configured, after approval
   ([deployments and environments](https://docs.github.com/en/actions/reference/workflows-and-actions/deployments-and-environments));
2. run only on reviewed release refs, never pull requests from forks or arbitrary workflow input;
3. pin third-party actions to full commit SHAs and allow only reviewed actions;
4. create an ephemeral Keychain, import the minimum signing identity for one job, restrict its ACL,
   sign, verify, export only public evidence, then destroy the runner/Keychain;
5. use a dedicated team notarization API key rather than an Apple Account password;
6. never echo secrets, enable shell tracing, place secrets in command arguments, upload Keychains,
   or persist `.p12`, `.p8`, passwords, or tokens as artifacts; and
7. treat a self-hosted runner as non-isolated and require a clean, single-use release-host design.

GitHub warns that command-line processes may expose arguments and that masking must happen before
values reach logs
([using Actions secrets](https://docs.github.com/en/actions/how-tos/write-workflows/choose-what-workflows-do/use-secrets),
[workflow masking](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-commands)).

## Storage and access policy

| Store | Allowed material | Prohibited material/use | Recovery and multi-maintainer note |
| --- | --- | --- | --- |
| User login Keychain | Local Apple Development identity; tightly controlled Developer ID identity on a dedicated release account; `notarytool` saved credential; references to Capsule keys | Daemon access; broad ACLs; routine installation-root availability; copied Keychain database in Git/artifacts | Back up only through an explicit encrypted custody procedure. Individual membership makes Dylan the Apple-team continuity point. |
| System Keychain | Public trust anchors or system-wide identity only if a later installer decision proves it necessary | Default storage for Capsule user-scoped signing, approval, evidence, or notarization secrets; granting root services access for convenience | Current per-user architecture needs none. Adding it expands authority and requires an ADR. |
| Secure Enclave plus Keychain reference | Nonexportable installation-root, Approval, and Supervisor evidence keys where algorithm/host/API support is proved; role- and epoch-specific access control | Apple code-signing identities, TUF offline keys, backup-by-export, daemon or shared-group access | Nonexportability means loss may require a new installation identity. Multiple maintainers cannot copy the key; recovery is an authorized replacement ceremony. |
| CI environment secret store | Later dedicated notarization team API key, encrypted `.p12` and separate password only if CI signing is admitted, public IDs/fingerprints | Apple Account primary password; long-lived repository/org-wide secrets; raw Keychain; secrets available to ordinary CI or PR jobs | Use reviewed environment gates, shortest lifetime, audit, revocation drill, and ephemeral runner. CI is deferred. |
| Offline release custody | TUF root and sensitive delegated role keys; master copy of notarization `.p8`; encrypted Developer ID recovery package only if policy permits export | Routine build-host availability, daemon/Supervisor access, one-person unrecorded copies | Prefer split custody and two-person release approval. Individual Apple membership still prevents a second person from independently managing Apple Developer team certificates. |
| Repository and artifacts | Public certificates, fingerprints, serials, validity, public keys, profile/entitlement digests, key IDs, signed TUF metadata, redacted logs, revocation status | Private keys; `.p12` or password; `.p8`; app-specific password; session/JWT/S3 tokens; raw Keychain export; unredacted profile dumps; recovery codes | Public metadata must still be reviewed for privacy and path leakage. Git history is not a secret store. |

Keychain access groups are capabilities, not folder names. Broker and Supervisor groups must be
disjoint and fresh for every identity-changing security epoch under
[ADR-0021](adr/0021-security-epoch-keychain-groups.md). No private key is migrated automatically
between epochs. App Groups must never be reused as authority-bearing Keychain groups merely because
the platform permits related sharing.

## Rotation, revocation, expiry, recovery, update, and rollback

### Apple Development certificate

1. Before expiry or planned replacement, create a new selected-team certificate on the authorized
   developer host.
2. Regenerate every affected development profile to include the new certificate.
3. Rebuild and re-sign every target; verify OU, TeamIdentifier, profile, entitlements, and CDHashes.
4. Run the complete installed wrong-role/stale-build/mixed-epoch corpus.
5. Revoke the old certificate only after no active development profile or authorized host needs it.

Wrong Team ID, missing private key, expired certificate, profile exclusion, or unverified subject
is `BLOCKED`. It is never repaired by relaxing a code requirement.

### Development or Developer ID provisioning profile

Regenerate after certificate, device, App ID, or capability change, and before expiration. Re-sign
the app with the replacement profile and test launch. An expired Developer ID profile is a release
incident because Apple may evaluate it at every launch. Keep a public expiration inventory and
alerts; do not rely on a cached Xcode filename.

### Developer ID certificate

- Planned rotation: create replacement, sign a new update, notarize, staple, test on a clean host,
  and update release metadata. Existing correctly timestamped apps may continue after ordinary
  certificate expiry, but new releases need a valid certificate.
- Suspected compromise: stop releases, revoke with Apple, quarantine all builds signed after the
  earliest possible compromise, rotate notarization credentials independently, publish Capsule
  release/TUF disable or replacement records as policy permits, and require forward repair.
- Loss without suspected compromise: recover only from approved encrypted custody; otherwise
  create a new certificate and release transition. Do not restore an unknown Keychain dump.

Apple states that revocation can prevent Developer ID apps from installing or launching
([Developer ID certificates](https://developer.apple.com/help/account/certificates/create-developer-id-certificates/)).

### Notarization credentials

- API key: revoke in App Store Connect, generate a replacement, update only the protected release
  environment, validate `notarytool history`, and destroy stale copies. A revoked key cannot be
  reinstated and its private half is downloadable only once.
- App-specific password: revoke it at the Apple Account page and store a replacement interactively.
  A primary Apple Account password reset revokes all app-specific passwords.
- Notarization rejection: preserve submission ID/log, classify signing, Hardened Runtime,
  entitlement, nested-code, or service-side failure, fix exact bytes, and submit a new artifact.
  Never staple or release a rejected/unknown submission.

### Capsule operational keys

- Installation-root loss creates a new installation identity. It does not authorize silent restore
  from backup.
- Approval-key loss or access-control failure disables new approval. Pending approvals do not move
  to a replacement key; authorize the replacement under the installation root and a new epoch.
- Supervisor evidence-key loss disables trustworthy new enforcement evidence and execution until
  an authorized forward epoch enrolls a fresh key.
- Suspected operational-key compromise sets the key `suspended` or `revoked`, disables attempts,
  preserves history, and performs a user-authorized forward trust transition. The daemon cannot
  clear this state.
- Identity-changing epochs use fresh, non-migrated Broker/Supervisor keys and Keychain groups.

### TUF and release provenance keys

Rotate online timestamp/snapshot roles through signed metadata with overlap only as the TUF profile
permits. Rotate targets/delegations within their path scope. Root rotation requires the threshold
and old/new root rules of the selected TUF profile, offline custody, independent review, and a
client update corpus. Losing the threshold of root keys is a release-trust recovery incident, not
an invitation to replace pinned roots through an unsigned channel.

### Update, repair, and rollback classifications

| Failure | Classification | Required response |
| --- | --- | --- |
| Team OU, TeamIdentifier, profile, or selected team disagree | `BLOCKED` | Stop signing/install; replace inputs under one team |
| Certificate/profile expired, absent, or missing private-key binding | `BLOCKED` | Renew/regenerate/recover through authorized workflow |
| Exact candidate path is proved unsupported, such as direct inherited Source Validator helper | `NO_GO` for that path only | Return to architecture; do not generalize to the whole capability |
| Notarization rejected or log unavailable | `BLOCKED` | Preserve ID/log and correct exact artifact |
| Mixed component version, CDHash, entitlement digest, key group, or epoch | `repair-required`; parent work `BLOCKED` | Keep execution disabled; finish or perform authorized forward repair |
| Revoked/compromised Developer ID or Capsule key | Security incident; `BLOCKED` | Revoke/suspend, quarantine affected output, rotate independently, forward repair |
| Coherent restoration of an older valid local world | Known rollback limitation | Do not claim detection without an independent non-rollbackable anchor or witness |
| Clean-host/minimum-OS evidence absent | Distribution admission `BLOCKED` | Label one-host result narrowly; fund/run required matrix later |

Repair restores exact enrolled components or creates an explicitly authorized forward epoch. It
never recreates missing authoritative state, resets consumed approvals, clears quarantine, or
rewinds trust history. See [Update and Recovery](UPDATE_AND_RECOVERY.md).

## What Dylan needs to do next

The completed I1B/R3 work is an input, not a task to repeat. Current actions and triggers are:

1. **Maintain active Individual Team `3DDR84M4JS`.** `[COST] [APPLE PORTAL]` Treat
   `W4QUR9FUL4` only as a historical certificate display suffix and keep membership, agreements,
   and eligibility current.
2. **Preserve the selected containing topology.** Accepted ADR-0037 fixes the visible Broker app
   with separately enrolled embedded daemon and Supervisor identities. Do not register the earlier
   Option A/Option B identifiers or other speculative duplicates.
3. **Treat I1B/R3 as `PASSED` only in its exact scope.** SHA-1 `80A4...D3793` and three exact
   Broker/daemon/Supervisor profiles passed the development-signed, installed,
   execution-disabled matrix retained in the pinned archive. They do not authorize later signing
   or product admission.
4. **Defer a separate central operational inventory while there is one maintainer.** Current
   redacted I1B/R3 metadata remains in the pinned archive and evidence ledger. Create a shared
   inventory when a second maintainer, production signing, or CI notarization requires it; do not
   combine personal Apple-account recovery, restricted release secrets, nonexportable Capsule
   keys, and offline TUF custody into one shared credential set.
5. **Resolve the I2B3 stale-profile blocker before another mutation run.** The exact Team-3DDR
   Coordinator/Supervisor profiles and signed entitlements passed, but the archived I1B Supervisor
   profile rewrote current-profile state in the stable private container. Select a versioned
   signing/container epoch in an ADR, provision it exactly, and separately authorize the remaining
   caller/key, App Group/service, protected-root, and descriptor-relative corpus. Do not reuse I1B
   profiles as generic Team authority. See the
   [I2B3 blocker result](MACOS_INSTALLATION_I2B3_SIGNING_PREFLIGHT_AND_STALE_PROFILE_BLOCKER.md).
6. **Rehearse machine-loss recovery deliberately.** Use the
   [recovery runbook](DEV_MACHINE_SETUP_AND_RECOVERY.md#part-c--machine-loss--reconstruction-drill)
   on a spare/replacement Mac or perform a supervised non-mutating portal walkthrough; retain the
   actual corrections and timing.
7. **Create or use Developer ID and notarization credentials only for separately authorized
   distribution work.** `[DEFERRED]` Use dedicated release custody; do not create Developer ID
   Installer unless `.pkg` is selected, and never store secret values in Git or shell history.
8. **Before production, fund clean-host/minimum-OS evidence and design release/TUF custody.**
   `[DEFERRED]` Select thresholds, operators, offline storage, update verifier/replacer authority,
   rollback posture, and incident drills. Do not give release or installation-root keys to the
   daemon or routine CI.

## Redacted public-metadata inventory template

Keep this table in a restricted operational inventory or a reviewed redacted repository record.
Rows contain public metadata only. Do not attach certificates/profiles unless the repository has a
specific reviewed reason; never attach private material.

| Item ID | Certificate/profile type | SHA-256 fingerprint | SHA-1 fingerprint | Subject OU / Team ID | Serial | Valid from / to | Public key algorithm | Profile UUID / name / expiration | Application identifier | Entitlements digest | Source | Owner | Status | Last verification |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `apple-dev-<team>-<date>` | Apple Development | `<HEX>` | `<HEX>` | `<TEAM_ID>` | `<HEX>` | `<UTC> / <UTC>` | `<ALG/BITS>` | N/A | N/A | N/A | `<Xcode or portal; host class>` | `<role/person>` | `active/suspended/revoked/replaced/expired` | `<UTC; command/tool version>` |
| `mac-dev-<role>-<epoch>` | Mac App Development profile | `<profile-file SHA-256>` | N/A | `<TEAM_ID>` | N/A | `<created / expires>` | `<embedded cert algorithms>` | `<UUID> / <redacted name> / <UTC>` | `<PREFIX>.<BUNDLE_ID>` | `<SHA-256 of normalized entitlements>` | `<portal or Xcode-managed>` | `<role/person>` | `<status>` | `<UTC>` |
| `developer-id-app-<team>-<date>` | Developer ID Application | `<HEX>` | `<HEX>` | `<TEAM_ID>` | `<HEX>` | `<UTC> / <UTC>` | `<ALG/BITS>` | `<Developer ID profile metadata if used>` | `<PREFIX>.<BUNDLE_ID if profiled>` | `<DIGEST if profiled>` | `<portal; release host class>` | `<release custodian>` | `<status>` | `<UTC>` |
| `notary-team-key-<date>` | App Store Connect team API key metadata | N/A | N/A | `<TEAM_ID>` | `<KEY_ID only>` | `<created / revoked>` | `<documented algorithm>` | N/A | N/A | N/A | `<App Store Connect>` | `<release custodian>` | `<status>` | `<UTC>` |

Known retained rows, with unknown fields deliberately left unknown:

| Item ID | Type | SHA-1 | Subject OU / Team ID | Status | Observation boundary |
| --- | --- | --- | --- | --- | --- |
| `apple-dev-legacy-3ddr` | Apple Development | `1638CFBD9250A00B4DBD81AE8FD1C790B42F61E3` | `3DDR84M4JS` | legacy; not selected | Exact prior X.509/signed-byte readback; current use is not authorized |
| `apple-dev-i1b-r3-3ddr` | Apple Development | `80A4969BCD1B3926020888094B9D812A283D3793` | `3DDR84M4JS` | selected for completed I1B/R3 only; later use not authorized | Exact subject/profile/signed-byte and installed readback retained in the pinned R3 archive |
| `developer-id-app-3ddr` | Developer ID Application | `AD70CEDCA605604676C2853A229AA4664AD3F750` | `3DDR84M4JS` | present; distribution use deferred and unauthorized | User-run identity output; no current distribution/notarization admission follows |

Never inventory private-key bytes, `.p12`/`.p8` contents, passwords, recovery codes, JWTs, session
tokens, temporary notarization upload credentials, raw Keychain exports, or secret-store values.

A separate central operational inventory is intentionally deferred while Capsule has one
maintainer and production signing/CI notarization remain deferred. Current redacted I1B/R3
metadata and hashes remain in the pinned archive and evidence ledger. Create the shared location
when a second maintainer, production signing, or CI notarization requires it; record only its
concrete pointer in the recovery runbook.

## Verification checklist for a credentialed task

The task handoff must retain only redacted evidence and answer each item:

- Membership Team ID equals certificate OU equals emitted `TeamIdentifier` equals profile team and
  application-identifier prefix.
- Every bundle/signing identifier is role-specific and matches the frozen installation profile.
- Every effective entitlement is allowed by the profile and the role matrix; no debug or broad
  exception leaks into shipping posture.
- Every nested Mach-O has a valid exact signature, expected team/signing identifier, expected
  CDHash, Hardened Runtime, and declared launch/library closure.
- Source Validator services are private to their containing role and parser results cannot cross
  roles.
- Daemon cannot access Broker/Supervisor Keychain groups, user-only content, updater/replacer,
  installation-root, Approval, evidence, or backend authority.
- App Groups, if admitted, are pairwise, empty of authority material, and not treated as peer
  authentication.
- Development output contains only development profiles; Developer ID output contains no
  development profile.
- Notarization status is accepted, the log is reviewed, the ticket is stapled where supported, and
  Gatekeeper is tested on the declared host.
- Installation, update, repair, and rollback results are labeled with the exact host, OS, SDK,
  Xcode, certificate/profile, team, build, and evidence limitations.

## Canonical related decisions and limits

- [Development machine setup and recovery](DEV_MACHINE_SETUP_AND_RECOVERY.md) sequences this
  document's setup sections into an execution checklist and adds an unrehearsed machine-loss
  recovery drill; this document remains canonical for the policy and rationale behind each step.
- [Apple Development provisioning plan](APPLE_DEVELOPMENT_PROVISIONING_PLAN.md) retains the exact
  W4/3DDR discovery, historical installed-test proposal, and current I1B/R3 reconciliation.
- [macOS installation and distribution plan](MACOS_INSTALLATION_AND_DISTRIBUTION_PLAN.md) defines
  the one-app packaging direction and unresolved bootstrap/updater/replacer authorities.
- [ADR-0029](adr/0029-select-authenticated-local-ipc-topology.md) defines the proposed
  daemon/Broker/Supervisor services and exact peer checks.
- [ADR-0036](adr/0036-select-role-separated-source-validator-launchers.md) selects two private
  role-specific launchers; exact inactive-policy R3 passed while active product work remains
  blocked.
- [ADR-0021](adr/0021-security-epoch-keychain-groups.md) requires fresh epoch-scoped Keychain
  groups and non-migrated keys.
- [Ecosystem reuse and adoption](ECOSYSTEM_REUSE_AND_ADOPTION.md) remains the dependency and
  provenance admission checklist for runtime, libkrun, update, TUF, and signing tooling.

Proposed ADR-0038 now resolves the Supervisor protected-root bootstrap owner in design: an
on-demand Trust Coordinator signs and the Supervisor creates. This guide does not make that path
installed evidence. ADR-0037 selects the containing-app identity. I2B3 exact Coordinator/
Supervisor profiles and signed entitlements passed, but the stable Supervisor container remained
writable by the archived I1B profile; a signing/container epoch decision now blocks the remaining
bootstrap App Group/service/key/root work. Ordinary Supervisor IPC App Group/private-service
topology, update verifier, Bundle Replacer, TUF operations, distribution location, minimum macOS
version, and `.pkg` path remain explicit decision points.
