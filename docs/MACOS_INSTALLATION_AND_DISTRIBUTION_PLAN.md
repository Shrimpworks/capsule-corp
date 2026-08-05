# macOS installation, update, repair, and uninstall plan

Date: 2026-08-04

## Status

```text
Work item: external installation architecture review and repository reconciliation
Status: PASSED
Scope: architecture and planning review only; no installer, service, key, store, update, repair,
  uninstall, runtime, backend, or guest behavior was implemented or exercised
Evidence or reason: the review supplied a coherent one-application distribution shape, explicit
  setup/update/repair/uninstall state machines, and a bounded spike order. The post-I1B platform
  research now classifies the seven Apple-platform questions and their remaining installed stops.
Remaining work: every installed mechanism and authority-bearing helper below remains subject to
  its named ADR, implementation, signed installed corpus, and product-admission gate.
Next action: after passed I2B2 unsigned construction, review the production signed-object wrapper
  and separately authorize I2B3 signing/key/App Group/service/container handoff evidence.
Parent status: macOS product installation and distribution are IN_PROGRESS — TRENDING_GOOD.
```

The intended product experience is one application, not a collection of separately installed
tools. Internally, Capsule still uses several narrowly enrolled processes because those process
boundaries prevent one compromised role from inheriting every key, state store, parser, updater,
and execution capability.

This plan is not an installer ADR and does not select the Bundle Replacer, updater implementation,
application location, or minimum supported macOS release. Proposed ADR-0038 now selects the narrow
Trust Coordinator/bootstrap responsibility; its installed evidence remains I2B work. This plan
defines what may be implemented now, what requires another decision, and which work can be deferred
without blocking the first developer MVP.

The canonical Apple-material and Capsule-key guide is
[Apple certificates, credentials, identifiers, entitlements, and Capsule keys](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md).
It records the exact W4/3DDR Team-ID stop, required replacement inputs, signable-component matrix,
credential custody, rotation/recovery rules, and the user action order without selecting any of the
unresolved authorities in this plan.

## Product shape

The current engineering direction is:

- distribute one self-contained `Capsule.app` in a signed DMG for ordinary installation;
- expose the Swift Approval Broker as the visible application and owner of setup, status, repair,
  update, and uninstall user interfaces;
- embed the per-user daemon and native-fronted Go Execution Supervisor as separately enrolled
  background components;
- embed governed runtime, libkrun, firmware, kernel, and runtime-root bytes in the application
  release rather than downloading ambient executable dependencies;
- keep the Source Validator behind the two role-specific, private, App-Sandboxed launcher
  boundaries selected by [Accepted ADR-0036](adr/0036-select-role-separated-source-validator-launchers.md),
  each of which owns one fresh Oxc child per validation;
- create one Runner process per execution attempt only after the Supervisor has durably consumed
  approval and created the attempt;
- use no permanent root daemon, `LaunchDaemon`, system extension, kernel extension, or privileged
  execution helper for the first release; and
- reserve a package containing the same application for later managed/MDM deployment. A package
  must not become Capsule's operational trust authority.

Users install, open, update, repair, and remove one product. The internal services are security
roles, not separate products the user configures manually.

Accepted [ADR-0037](adr/0037-freeze-passive-macos-installation-i0-contract.md) and the
[I0 passive contract](protocol/MACOS_INSTALLATION_I0_PASSIVE_CONTRACT.md) now freeze the exact
one-visible-app, seven-role, service, entitlement, bootstrap, update, repair, and uninstall
candidate without activating it. Signing, bootstrap-owner, product-store, and IPC-transport values
that current evidence cannot support are explicit inactive refusal gates.

## Authority boundaries

| Role | May own | Must not own | Current status |
| --- | --- | --- | --- |
| Swift Broker application | setup/status UI, registered-plan rendering, user-presence approval, user-facing repair and uninstall choices | Supervisor state, backend lifecycle, daemon credentials, automatic trust-state repair | `IN_PROGRESS — TRENDING_GOOD`; product target absent |
| Agent-facing daemon | agent protocol, proposal intake, untrusted source submission, registered identifiers returned by the Supervisor | approval keys, installation root, Supervisor evidence key, content custody, updater/replacer or backend authority | architecture selected; product IPC absent |
| Execution Supervisor | authoritative installation/epoch view, registration/approval/attempt/lifecycle state, owner lock, recovery, backend lifecycle | network update resolution, rich UI, arbitrary parsers, daemon-controlled replacement bytes | local mechanics `PASSED`; installed composition `BLOCKED` |
| Daemon Source Validator launcher | one copied validation request, one fresh child, bounded drain/kill/reap, private disposable scratch residual | Broker requests, Capsule keys/stores, network, package loading, runtime/backend paths | ADR-0036 architecture and R1 passive contracts `PASSED`; product work `BLOCKED` on R2-R5D |
| Broker Source Validator launcher | same narrow parser operation for Broker-retained source | daemon requests, approval key use, Capsule stores, runtime/backend paths | ADR-0036 architecture and R1 passive contracts `PASSED`; product work `BLOCKED` on R2-R5B |
| Runner | one sealed Supervisor-issued attempt descriptor and the closed runtime/backend descriptor set | registration, approval, update, repair, installation-root, or arbitrary host-path authority | `IN_PROGRESS — TRENDING_GOOD`; not admitted |
| Update verifier/agent | pinned release metadata retrieval and verification; bounded prepared-update output | installation-root key, Supervisor mutation, replacement authority, execution authority | later beta/production design; not selected |
| Trust/bootstrap coordinator | one user-authorized installation/epoch ceremony over closed inputs; installation-root signing of the I2 request/record | routine background operation, network fetching, Supervisor store mutation, execution | selected by Proposed ADR-0038 as an on-demand XPC role; I2B evidence blocked |
| Bundle replacer | if selected, mechanical replacement of one pre-authorized exact bundle | network, target selection, trust decisions, store repair, execution | new authority role; ADR and replacement evidence required |

The last three rows are intentionally not treated as implementation details. Proposed ADR-0038
supplies the Trust Coordinator's closed method/identity/fault contract. A Bundle Replacer and any
later Coordinator update/repair method still require separate decisions and evidence.

## Installation and protected state

The intended authoritative state location is the Execution Supervisor's private App Sandbox
container. Mode bits and a BSD lock alone do not protect state from another process running as the
same user. The installed composition must prove that the admitted Supervisor identity receives the
container and that baseline same-user processes cannot replace the state-root or owner-lock
entries.

Proposed [ADR-0038](adr/0038-select-one-shot-coordinator-supervisor-bootstrap.md) and the
[I2A decision/fault plan](MACOS_INSTALLATION_I2A_PROTECTED_ROOT_BOOTSTRAP_DECISION.md) now select
the narrower composition. A separately signed, on-demand Trust Coordinator owns the
user-presence-gated installation-root key and signs the closed request/final record. The
authenticated Supervisor alone creates and observes the fixed root/lock/disabled-store genesis in
its own container. The visible app only invokes setup UI and `SMAppService` registration; it gains
no signing key or Supervisor-state authority. I1B is complete; installed protected-state support
remains blocked on I2B.

The selected composition requires:

1. creation is exclusive, no-follow, one-time, synchronized, and fully reopened;
2. the signed bootstrap record binds installation, Supervisor, state-root, owner-lock, store
   format, expected UID, exact names, and enrolled object identities;
3. ordinary startup opens without creation and refuses missing, replaced, relocated, linked,
   incorrectly owned, or incorrectly permissioned objects;
4. binary update preserves the protected root and owner-lock identity; and
5. inability to prove continuity disables attempts and enters explicit repair or new-installation
   handling. It never creates a plausible empty replacement.

## IPC and App Group posture

The product still follows Proposed ADR-0029: one native-fronted Go Supervisor process and two
role-specific Supervisor services, one for the daemon and one for the Broker. Service identity,
dynamic code validity, exact enrolled component identity, session, installation, epoch, closed
method, and closed message shape remain independently required.

An App Group used for Mach/XPC naming is a real residual capability. It also provides a shared
container namespace and can be used as a Keychain access group. Therefore the signed IPC spike
must prove and retain all of the following:

- exactly two pairwise groups if groups are selected at all: daemon/Supervisor and
  Broker/Supervisor;
- no files, defaults, sockets, manifests, source, content, state, or migration material in either
  group container;
- no Keychain item assigned to either IPC group;
- no daemon membership in any installation-root, approval, Supervisor-evidence, content, or
  updater/replacer access group;
- wrong-role, wrong-group, stale-build, mixed-epoch, wrong-session, and unexpected-container
  attempts refuse; and
- service names, bundle identifiers, signing identifiers, entitlements, Team identity, active
  CodeDirectory hashes, and peer requirements are generated from one installation profile.

App Group membership is never described as peer authentication or structural absence of shared
authority. If the residual shared namespaces cannot be accepted and tested, a different supported
private service topology needs a new decision.

## Setup state machine

The first-run UI may be inside the Swift Broker application, but authority remains split. The
planned sequence is:

```text
preflight application location and exact signed bundle
  -> create bounded setup journal with execution disabled
  -> create installation identity through an explicitly authorized ceremony
  -> prepare purpose-separated operational-key enrollments
  -> register and authenticate the Supervisor service
  -> bootstrap and enroll the protected Supervisor root and owner lock
  -> register and authenticate the daemon
  -> validate every embedded role, entitlement, artifact, runtime, and backend profile offline
  -> construct and authorize InstallationManifest plus trust epoch 1
  -> run owner/store/recovery/IPC self-checks
  -> enable new attempts in the final durable transition
```

Before the epoch commit, interruption leaves no usable installation and cleanup may remove only
objects proven to belong to the incomplete setup transaction. At or after the epoch commit,
ordinary startup never starts over. It resumes the exact transition or enters `repair-required`.

Setup must explain in plain language that Capsule installs per-user background components, that
the Supervisor owns security state, and that execution remains disabled until all checks pass. It
must not ask for Full Disk Access, Accessibility, Automation, Endpoint Security, root, or broad
file access as an installation shortcut.

## Update progression

Automatic updating is not an MVP dependency.

### Development MVP

- Developer-signed application assembled and launched from one declared location.
- Manual replacement of the complete application bundle only.
- Explicit stop, replacement, `SMAppService` re-registration where required, restart, exact
  component verification, and recovery.
- No background update schedule, network TUF client, custom Bundle Replacer, rollback claim, or
  production distribution claim.

### Signed beta

- User-triggered update check through an isolated metadata verifier.
- One bounded `PreparedUpdate` binding exact predecessor, target bundle, component identities,
  migrations, trust material, and recovery disposition.
- New processes remain execution-disabled until store, owner-lock, component identity, IPC, key,
  migration, and recovery verification passes.
- Interruption is predecessor-or-successor or `repair-required`; never mixed-ready.

### Production distribution

- Reviewed TUF profile and release-role operations outside the live Supervisor path.
- Separately reviewed replacement mechanism with no target-selection or trust authority.
- Developer ID signing, notarization, stapling, Gatekeeper, nested-code, update/replacement,
  repair, and uninstall evidence on the declared support matrix.
- Explicit rollback and coherent-backup posture; no monotonicity claim without an independent
  non-rollbackable checkpoint or external witness.

Normal updates preserve the Supervisor private root and owner lock. A release that changes their
location or cannot preserve their identities is an explicit migration/repair ceremony, not a
routine bundle swap.

## Repair and uninstall

Repair restores exact enrolled components or performs an authorized forward trust transition. It
does not recreate state, mark unknown cleanup complete, unconsume approvals, or discard retained
archive and replay-prevention records.

The UI must distinguish these operations:

1. **Repair application files** — restore the current authorized release while preserving
   installation identity and security history.
2. **Remove the application, preserve local Capsule state** — stop/unregister services and remove
   application bytes while leaving a clearly documented reinstall/recovery path.
3. **Remove local content and keys where safe** — refuse or defer objects still required for
   cleanup, quarantine, replay prevention, archive, or repair.
4. **Abandon this installation identity** — make the local installation unusable and record the
   limits of local erasure. It must not claim to erase external release logs, witnesses, receipts,
   backups, or already exported evidence.

Deleting `Capsule.app` in Finder cannot be assumed to perform the security-state uninstall
ceremony. On later launch/reinstall, Capsule must detect the retained installation, incomplete
transition, or missing authoritative state and present the corresponding recovery choice.

## Platform and distribution support

No minimum macOS version is frozen. macOS 14 may remain a candidate floor, but final selection is
`BLOCKED` on the exact libkrun/runtime dependencies, Service Management behavior, container and
update-policy mechanisms, source-validator launchers, signed bundle, and a declared support test
matrix. Existing firmware or framework bytes may raise the effective floor.

Paid clean-host/minimum-version coverage is not currently planned. That does not block local
development composition, but it does block a supported distribution claim. Results on the current
owned host must be labeled with the exact host, OS, SDK, signature, and provisioning posture.

The ordinary installation location also remains open. `/Applications` provides conventional UX
but may require a replacement authorization path depending on ownership. `~/Applications` is
user-writable but has different discovery, update, and support consequences. The selected update
mechanism must determine the supported location; the UI cannot promise both before replacement and
service-registration evidence exists.

## Ordered work

| Slice | Scope | Status and exit |
| --- | --- | --- |
| I0 | Passive application bundle, role, entitlement, service, bootstrap, update, repair, and uninstall contract | `PASSED`; exact inactive profile, generated cases/digests, field authority, and pure missing/mixed/extra/transition/retention validators retained; no installed behavior |
| I1A | Unsigned app shell and exact seven-role byte/layout construction, execution permanently disabled | `PASSED`; deterministic Swift/AppKit status shell, inert daemon/Supervisor placeholders, exact R2 identities, closed readback and refusal evidence; no signed or installed behavior |
| I1B | Developer-signed app shell, exact effective entitlements and installed daemon/Supervisor/private-XPC placement, execution always disabled | `PASSED`; exact Team-`3DDR84M4JS` profiles, inside-out signing/readback, private-service reachability, eight refusal cases, per-user registration, and cleanup are retained in the [pinned R3 archive](https://github.com/Shrimpworks/capsule-experiments/tree/0944ffd8cfd01ec23e4ae99138b0931d56804077/experiments/completed-compiled-artifact-payloads/payloads/capsule-corp/artifacts/macos-i1b-r3-signed-development-composition); Apple Development only, no product admission or Release asset claim |
| I2A | Protected-root/bootstrap-owner architecture and passive contract | `PASSED`; Proposed ADR-0038 selects one-shot Coordinator authorization plus Supervisor creation and retains the exact I2B object/order/fault plan without installed behavior |
| I2B | Protected Supervisor container, one-time signed bootstrap, ADR-0033 owner-lock/descriptor-relative fixed-v1 open, same-user mutation and restart corpus | I2B1 passive request/record objects and I2B2 unsigned installation-only construction are `PASSED`; installed I2B is `BLOCKED` on production wrapper review and separately authorized I2B3 exact Team-3DDR signing/key/App Group/service/container handoff evidence, followed by I2B4-I2B5 installed faults; product-store selection remains separate |
| I3 | Pairwise authenticated daemon/Broker IPC plus two role-specific Source Validator launchers | `BLOCKED` on I2, ADR-0029's Supervisor App Group/private-service decision, and ADR-0036 R4 confinement/resource/residue evidence plus later consumers; R1-R3 are `PASSED` in their exact scopes |
| I4 | Manual whole-bundle replacement, service re-registration, mixed-version refusal, forward repair, retained-state recovery | `BLOCKED` on I2/I3 and replacement-authority decision |
| I5 | User-triggered TUF verification and reviewed mechanical replacement | later `IN_PROGRESS` only after I4; not an MVP dependency |
| I6 | Developer ID/notarized beta, support-floor matrix, backup/restore/uninstall and distribution admission | deferred; paid clean-host coverage is not currently planned |

Runtime/backend admission remains separate. I0-I4 can use `FakeBackend.CreatesGuest() == false`
and must not connect user bytes to libkrun merely because installation mechanics pass.

## ADR boundaries

Before implementation crosses the named authority boundaries, retain separate decisions for:

1. the one-app/DMG and embedded per-user service topology, including supported install location;
2. the exact protected-container bootstrap owner, now selected by Proposed ADR-0038, followed by
   its I2B installed evidence;
3. App Group/private-service residual authority and final Mach/XPC naming;
4. the two Source Validator launchers and quantified reactive-resource contract, now selected by
   Accepted ADR-0036 while passive contracts and installed evidence remain blocked;
5. Trust Coordinator authority, now narrowed by Proposed ADR-0038 to the on-demand bootstrap
   request/record signer; later update/repair methods still require their own decision;
6. update verifier/TUF profile and `TrustSnapshot` contract;
7. Bundle Replacer authority and complete-bundle replacement transaction; and
8. uninstall/local-erasure semantics and archive/external-evidence retention.

Do not combine all eight into one approval. Each has a different compromise consequence and
different installed evidence.

## Completed platform research and remaining installed questions

The seven-question brief below is answered by the
[post-I1B Apple-platform semantics research](MACOS_INSTALLATION_PLATFORM_RESEARCH.md). That
research slice is `PASSED`; it does not advance I2B, I4, or I6 installed evidence. In particular,
the exact signed Coordinator session/profile, App Group residue, stale-signer Keychain behavior,
replacement interruption recovery, and clean-host/multi-OS distribution claims remain `BLOCKED`.

The original prompt remains here as the closed scope and challenge criteria for the retained
result.

### Completed Apple-platform research prompt

> Review the current Capsule macOS installation plan as a defensive architecture research task.
> Capsule is one visible Swift application with an unprivileged native-fronted Go Execution
> Supervisor, an agent-facing daemon, two role-specific private Source Validator launchers, and
> later optional update-verifier and bundle-replacer roles. The Supervisor is intended to own its
> authoritative state in its own App Sandbox container. Capsule does not permit a permanent root
> daemon, privileged execution helper, daemon-to-backend path, shared authority store, or automatic
> recreation of missing security state.
>
> Use current official Apple documentation, current public SDK declarations, and WWDC sessions as
> primary sources. Clearly separate documented guarantees, reasoned inferences, observed behavior,
> and unresolved questions. Do not recommend private APIs, deprecated custom sandbox profiles, or
> broad temporary exceptions. Do not assume a paid clean-host/minimum-OS lab is available.
>
> Research these exact questions:
>
> 1. What supported bundle and Service Management layout can embed a per-user App-Sandboxed
>    Supervisor and daemon in one containing app, and what are the exact `SMAppService`
>    registration, approval, replacement, unregister/re-register, login, logout, and restart
>    semantics when the containing app is replaced?
> 2. Validate ADR-0038's selected handoff: the containing app invokes its on-demand private Trust
>    Coordinator, registers the Supervisor LaunchAgent, and only the exact Coordinator/Supervisor
>    pair use one bootstrap App-Group-named XPC service. Determine whether the exact nested bundle,
>    container association, service launch, code-signing/provisioning, and two-message sequence are
>    supported; identify any bootstrap circularity. If Apple does not document a guarantee, say so
>    and keep I2B blocked rather than choosing a wider fallback.
> 3. For sandboxed peer XPC, when are App Groups required for Mach service naming or lookup? What
>    shared-container, Keychain, preferences, or other capabilities follow from membership? Compare
>    two pairwise groups with supported private XPC-service alternatives, and identify what can be
>    structurally denied versus only left unused and tested empty.
> 4. What does `NSUpdateSecurityPolicy` protect, on which macOS releases is it supported, and how
>    does it interact with App Sandbox containers, same-Team stale binaries, nested helper apps,
>    launch agents, packages, app replacement, and Developer ID distribution? Do not assume it
>    solves Keychain stale-membership or coherent rollback.
> 5. For applications in `/Applications` and `~/Applications`, what supported unprivileged or
>    user-authorized mechanisms can atomically replace one signed/notarized app bundle while
>    preserving its per-user services and private containers? Identify ownership, authorization,
>    running-code, quarantine, translocation, Gatekeeper, nested-code, crash, rollback, and
>    same-volume atomicity limitations. Compare a small custom replacer with stock Sparkle only as
>    mechanisms; neither may become Capsule's trust authority.
> 6. Validate the exact authority ADR-0038 assigns to the on-demand Trust Coordinator: one
>    Coordinator-only installation-root Keychain group, fresh user presence, the two closed
>    request/record signing purposes, and the bootstrap-only pairwise App Group. Confirm denial to
>    the visible app, daemon, updater, replacer, and continuously running Supervisor. Include Secure
>    Enclave/LocalAuthentication availability, stale same-Team code, process-death, replay, and
>    update consequences.
> 7. Which claims can be tested on one owned current macOS host with Apple Development and Developer
>    ID identities, and which claims fundamentally require separate clean hosts or multiple OS
>    versions? Produce a minimal local evidence matrix and mark distribution-only evidence as
>    deferred rather than pretending it is available.
>
> Return: (a) a concise decision table for each question; (b) exact primary-source links and quoted
>    API/availability facts within fair-use limits; (c) a proposed supported process/bundle tree;
>    (d) an authority and entitlement matrix; (e) first-run and replacement sequence diagrams;
>    (f) failure and interruption tables; (g) the smallest developer-signed spikes that can run on
>    one current owned Mac; and (h) a list of claims that remain blocked without clean-host or
>    minimum-OS testing. Use Capsule's status terms precisely: `PASSED`, `IN_PROGRESS —
>    TRENDING_GOOD/BAD`, `BLOCKED`, and `NO_GO` only for an abandoned exact path.
