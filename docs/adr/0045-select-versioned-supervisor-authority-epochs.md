# ADR-0045: Select versioned Supervisor authority epochs

- Status: Proposed
- Date: 2026-08-05
- Refines if accepted: ADR-0012, ADR-0021, ADR-0029, ADR-0033, ADR-0038, and ADR-0040
- Evidence boundary: current Apple public documentation plus the retained Apple Development I2B3
  observation; no Developer ID/notarized or installed successor-epoch evidence

## Context

Proposed ADR-0038 selects who may authorize and create the initial Supervisor protected root. It
does not select what macOS identity owns that authority across replacement. I2B3 proved the gap:
an archived, still-valid Apple Development Supervisor signed with the stable identifier
`com.capsulecorp.capsule.supervisor` rewrote a sentinel created by the new-profile Supervisor in
the same App Sandbox private container. A changed provisioning profile, CDHash, App Group, and
Keychain projection did not version that container.

This is narrower than a general same-user-containment result. It also does not show that a stale
process can cross a distinct App Sandbox container, App Group, or Keychain access group. It shows
that Capsule cannot call the stable Supervisor container a current-epoch authority boundary.

Current Apple documentation, retrieved 2026-08-05, supplies a candidate but not Capsule evidence:

- App Sandbox creates and associates a private container with a sandboxed app. On macOS 14 and
  later the operating system uses code-signing identity in that association, may request user
  authorization for foreign-container access, and tracks launch-agent container association.
- Apple's migration documentation names the ordinary container by bundle identifier, while code
  signing documentation explains that a bundled code-signing identifier is typically its bundle
  identifier and that designated requirements express code identity.
- App Group membership grants shared-container and IPC authority. Apple documents both the
  iOS-style `group.<name>` form and the macOS-only `<TeamIdentifier>.<name>` form; the latter does
  not require Developer-website App Group registration. A group may also add a Keychain access
  group, while the group-prefixed name is the supported Mach/XPC naming form.
- A Keychain item belongs to one access group. A caller outside that group cannot add an item to it,
  and queries must name the exact group or they search every group available to the caller.
- `SMAppService` registers in-bundle LaunchAgents, but does not define Capsule authority,
  continuity, rollback, or state migration.

It is therefore a documented mechanism plus a design inference—not an observed control—that a
fresh explicit App ID/bundle/signing identifier should obtain a distinct private-container
association and that fresh App Groups and Keychain groups should exclude prior members. The
bounded experiment in
[the I2B3 authority-epoch packet](../MACOS_INSTALLATION_I2B3_SUPERVISOR_AUTHORITY_EPOCH_EXPERIMENT.md)
must test that inference before this ADR can be Accepted.

## Proposed decision

### A separate `SupervisorAuthorityEpoch`

Add a sequence-ordered `SupervisorAuthorityEpoch` whose sole purpose is to name the complete OS
and protocol authority that may open one Supervisor state world. It is not:

- an application release, `CFBundleVersion`, marketing version, or CDHash;
- an ADR-0012 installation trust epoch, `InstallationManifest`, TUF version, `TrustSnapshot`,
  runtime/profile identity, or distribution channel; or
- a claim of rollback monotonicity.

One release or installation trust epoch binds exactly one Supervisor-authority-epoch digest. A
release may change without changing Supervisor authority only when the complete authority tuple
below is byte-for-byte unchanged and the ordinary trust-epoch rules admit its exact code. Any
change to a tuple member requires a new Supervisor authority epoch even if the application release
number does not change. The authority sequence starts at one and increments by exactly one; release
versions may skip independently.

`SupervisorAuthorityEpochV0` binds at least:

- object version, installation ID, nonzero authority sequence, predecessor authority-epoch digest
  or the fixed initial-install absence value, and transition reason;
- exact explicit Supervisor App ID, bundle identifier, code-signing identifier, Team/channel,
  allowed CDHashes, effective-entitlement digest, hardened-runtime requirements, and active
  provisioning-profile digest/UUID for Apple Development evidence;
- exact Supervisor `SMAppService` LaunchAgent label and bundle-relative plist/program identities;
- the platform private-container association identified by the versioned application identity and
  obtained only through supported container/standard-directory APIs, never a pathname in a signed
  request or a guessed `~/Library/Containers/...` path;
- exact versioned Trust Coordinator App ID/bundle/signing identity and its current peer
  requirement;
- exact versioned bootstrap App Group and Mach service; the group remains IPC-only residual
  capability and holds no Capsule key, file, default, socket, state, migration bytes, or authority;
- exact role-specific Coordinator installation-root, Supervisor bootstrap-anchor, Supervisor
  evidence, and any later operational Keychain access groups plus their public-key authorization
  digests; every Keychain operation names its group explicitly;
- the exact `SupervisorAuthorityDescriptorV0` digest, bootstrap request/record schema identities,
  state schema and format, selected authority/lifecycle engine identity, owner-lock mechanism,
  and create/open/migration disposition; and
- the separately bound installation-trust-epoch and release-manifest digests used to admit the
  bytes. Those references do not become Supervisor-authority identity.

The descriptor contains identifiers and digests, not private keys, profiles, live host paths, or
caller-selected replacement bytes. ADR-0038 request and record objects must bind its complete
digest. Normal startup must find the same digest in the create-once Supervisor Keychain anchor and
protected-root record before opening the owner or state engine.

### Version every authority-bearing identity

The first candidate uses one generated lowercase decimal authority token `e<N>` consistently, for
example:

| Role | Candidate pattern |
| --- | --- |
| Supervisor App ID, bundle ID, and signing ID | `com.capsulecorp.capsule.supervisor.authority-e<N>` |
| Supervisor LaunchAgent label | `com.capsulecorp.capsule.supervisor.authority-e<N>` |
| Trust Coordinator App ID, bundle ID, and signing ID | `com.capsulecorp.capsule.trust-bootstrap.authority-e<N>` |
| Bootstrap App Group | `3DDR84M4JS.com.capsulecorp.capsule.bootstrap.authority-e<N>` |
| Bootstrap Mach service | `<bootstrap-group>.supervisor` |
| Coordinator installation-root group | `3DDR84M4JS.com.capsulecorp.capsule.trust-bootstrap.installation-root.authority-e<N>` |
| Supervisor bootstrap-anchor group | `3DDR84M4JS.com.capsulecorp.capsule.supervisor.bootstrap-anchor.authority-e<N>` |
| Supervisor evidence group | `3DDR84M4JS.com.capsulecorp.capsule.supervisor.evidence.authority-e<N>` |

These are candidate strings to validate against the exact Team portal/profile and signed installed
composition. A production Developer ID/notarized profile may encode App Group authorization
differently; it must preserve the same logical separation and obtain its own evidence. Apple
Development profile success on one enrolled Mac cannot be cited for shipping distribution.

The frozen bootstrap App Group is deliberately the macOS-style
`<TeamIdentifier>.<group-name>` entitlement, not a Developer-portal App Group resource. The
zero-effect portal preflight retained at `capsule-experiments` merge
`e6390253a274e9ead76366f9869a5e1b272a1595` observed the portal rewriting that input by prepending
`group.` and therefore rejected portal registration before mutation. Merge
`3671a6eb23357ff28de4562dd60e8f68173034ae` corrects the disposition: the portal-registration path
is `NO_GO`; the frozen identity itself remains the intended Proposed candidate and is `BLOCKED` on
signed-profile and E1 execution evidence. Implementations must not substitute the rewritten
iOS-style identifier. Exact explicit App IDs, embedded profiles, signatures, effective
entitlements, and `com.apple.application-identifier` bindings remain mandatory pre-launch
evidence.

The stable Supervisor identity `com.capsulecorp.capsule.supervisor`, its private container, the
I1B and I2B3 profiles, and the unlaunched I2B3 Coordinator identity are **legacy residue**, not
authority epoch zero or a predecessor installation. No successor request, record, state engine,
repair, or migration may adopt bytes from that container. It may be inventoried from outside only
with non-authoritative bounded metadata; its absence is not required to create epoch one.

### OS enforcement, protocol detection, and availability are different claims

The selected boundary composes three distinct controls:

1. **OS-enforced nonmembership.** A stale role lacks the new explicit application identity, App
   Group, and Keychain groups. The required installed corpus must show sandbox denial and
   `errSecMissingEntitlement`/no-match behavior without a consent override. This is the only layer
   that may support a stale role's inability to open or mutate current private state.
2. **Signed/protocol detection.** Exact peer requirements, message-derived code checks, accepted
   CDHashes, entitlement digest, `SupervisorAuthorityEpoch` sequence/digest, installation trust
   epoch, signed request/record, and state/engine bindings reject stale or mixed messages and
   objects. These checks detect substitution; they do not turn a shared writable container into
   separation.
3. **Availability and retirement.** A stale role may still start against, corrupt, delete, or lock
   its own old container/key/service; retain old signatures; request foreign-container consent;
   or consume user and launchd attention. Those are denial/elevated-posture concerns, not current
   state authority. Capsule does not claim that identity rotation removes them.

On macOS 14 and later, a foreign-container access request can present user consent. Capsule never
requests or treats that consent as a repair mechanism. A prompt, a granted foreign-container
exception, Full Disk Access, task-port control, or comparable broad authority moves the observation
to the threat model's elevated user-granted tier. The current authority remains attempts-disabled
and the baseline nonmembership claim is not made until a separately reviewed recovery posture
establishes clean state. Denial is the expected experiment oracle; granting consent is a deliberate
negative mutation, not success.

### Initial transition from legacy residue

Epoch one is a new initial installation transaction, not an update from the stable identity:

1. Verify the candidate descriptor and exact signed bundle while attempts, ordinary Supervisor
   IPC, runtime, backend, and guest remain absent.
2. Prove the legacy LaunchAgent disabled/unregistered and its process absent without entering or
   modifying the legacy private container.
3. Register only the versioned epoch-one Supervisor and require enabled status before invoking the
   versioned Coordinator.
4. Require fresh user presence; create the epoch-one installation-root key only in the exact
   epoch-one Coordinator group; then run ADR-0038's request, protected-root genesis, observation,
   record, and create-once anchor sequence with the authority-descriptor digest bound throughout.
5. Reopen through supported APIs, acquire the enrolled owner, verify the exact state schema/engine
   genesis, and remain attempts-disabled through the complete I2B4/I2B5 matrix.

Any unexpected existing target container contents, foreign-container prompt, legacy path alias,
mixed identity, wrong group, skipped authority sequence, or state/engine mismatch stops before key,
service, root, or store creation where ordering permits; otherwise it enters `repair-required` and
never normalizes or adopts the state.

### Later transition, repair, and rollback rules

A later authority transition is always forward and follows these rules:

- prepare authority sequence `N+1` against the exact active `N` digest; there is no skipped
  authority sequence, wildcard predecessor, or old-or-new active requirement;
- disable new approvals and attempts, close ordinary listeners, and require every active attempt
  durably terminal with authoritative absence before exporting or changing authority;
- create fresh Supervisor/Coordinator application identities, private container, App Group/Mach
  service, role-specific Keychain groups, keys, root, owner, and state-engine world; never reuse or
  relocate the predecessor root/owner;
- require a separately frozen, bounded, signed, content-addressed state-handoff format whose
  schema/engine transition preserves the complete grant, attempt, cleanup, replay/non-reuse,
  archive, quarantine, and trust history. No such product handoff is selected or implemented now;
  until it is accepted with fault evidence, a state-bearing authority transition is `BLOCKED`;
- commit target authority logically once, bind it into a new ADR-0012 trust epoch, accept every
  exact current component, and only then retire the predecessor service/process/private keys;
- never reopen attempts if predecessor retirement, current anchor/root agreement, migration,
  component acceptance, or state-engine verification is indeterminate; and
- treat predecessor-container or predecessor-key loss as availability/evidence loss, never as
  permission to recreate, roll back, or mark cleanup complete.

Same-authority repair may restore byte-identical admitted application bytes and re-register the
same exact service only while the current protected root, anchor, owner identity, descriptor, and
state engine all verify. Changing any authority-tuple member is a forward authority transition.
Missing/corrupt current root, anchor, owner, or authoritative state is `repair-required`; absent an
independently anchored complete handoff or backup admitted by a later ADR, Capsule abandons that
installation identity and invalidates pending authority.

Rollback never reactivates a retired bundle identifier, LaunchAgent label, container, App Group,
Keychain group, key, descriptor, or state snapshot. Running older functional code requires building
and admitting those bytes as a new release under authority sequence `N+1`, with a new descriptor,
fresh identities and keys, and a supported forward state-schema transition. If the older code
cannot read the current schema through that new transition, rollback is unavailable. No coherent
rollback-resistance claim follows without ADR-0012's independent anchor or witness.

## Alternatives and dispositions

- **Stable Supervisor bundle/signing identifier with only a new profile, CDHash, App Group, or
  Keychain group:** `NO_GO` for I2B. The exact I2B3 mutation proved those changes leave stale write
  authority over the stable private container.
- **Authoritative state, migration bytes, epoch pointer, key, or ledger in the bootstrap App
  Group:** `NO_GO`. Membership deliberately spans Coordinator and Supervisor and therefore cannot
  be the single-member Supervisor state boundary.
- **Manually guessed or signed host container path:** `NO_GO`. The platform chooses the container;
  paths are neither authority nor a portable contract.
- **Rollback by re-registering or reopening a retired authority epoch:** `NO_GO`. It restores stale
  authority and can resurrect consumed or unresolved state.
- **Stable container plus cryptographic detection alone:** not selected. The current evidence shows
  stale mutation, while no independently protected, complete detection/repair mechanism exists.
  A future proposal would need to survive arbitrary tear, deletion, replay, and coherent rollback;
  it is not a fallback for this I2B slice.

## Consequences and blockers

- ADR-0038 remains the Proposed owner of bootstrap actor, request/record, and protected-root
  creation. This ADR supplies the missing authority identity and transition boundary rather than
  rewriting ADR-0038's historical decision.
- The no-root topology remains unchanged. No daemon, helper, visible app, updater, or replacer gains
  installation-root, Supervisor-state, or backend authority.
- State remains in one Supervisor-private container. App Groups remain empty residual IPC
  capabilities, and no live path crosses the boundary.
- The initial authority-epoch design is complete enough to freeze the next experiment, but it is
  not implemented or validated. This ADR remains Proposed.
- The bounded research and documentation slice is `PASSED` in its exact retained scope; repository
  adoption remains Proposed pending review. Installed owner-lock G3/I2B remains `BLOCKED` on
  separately authorized profile/signing readback followed by a fresh disposable-container
  nonmembership mutations, Coordinator/Supervisor service/session, Keychain, root, owner, state
  engine, restart, and retirement evidence.
- Apple Development evidence can support only the exact named host, OS, certificate, profiles, and
  development distribution posture. Developer ID, notarization, Gatekeeper, clean-host,
  translocation, minimum-OS, package/update, and shipping retirement claims remain separate I6
  work.

## Evidence and next experiment

- [I2B3 stale-profile blocker](../MACOS_INSTALLATION_I2B3_SIGNING_PREFLIGHT_AND_STALE_PROFILE_BLOCKER.md)
- [I2B3 authority-epoch inert packet and mutation matrix](../MACOS_INSTALLATION_I2B3_SUPERVISOR_AUTHORITY_EPOCH_EXPERIMENT.md)
- [Apple platform-semantics research](../MACOS_INSTALLATION_PLATFORM_RESEARCH.md)
- [Apple: accessing files from the macOS App Sandbox](https://developer.apple.com/documentation/security/accessing-files-from-the-macos-app-sandbox)
- [Apple: App Groups entitlement](https://developer.apple.com/documentation/bundleresources/entitlements/com.apple.security.application-groups)
- [Apple: protecting local app data using containers](https://developer.apple.com/documentation/xcode/protecting-local-app-data-using-containers)
- [Apple: Keychain access groups](https://developer.apple.com/documentation/security/sharing-access-to-keychain-items-among-a-collection-of-apps)
- [Apple: `SMAppService`](https://developer.apple.com/documentation/servicemanagement/smappservice)
- [Apple: code-signing requirements](https://developer.apple.com/documentation/technotes/tn3127-inside-code-signing-requirements)
