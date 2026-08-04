# Installation Trust

Status: intended design; Gate B/F mechanism and process-fault evidence is observed, while product
store, installer, distribution, and power-loss validation remain pending.

The staged product packaging and bootstrap questions are tracked in the
[macOS installation and distribution plan](../MACOS_INSTALLATION_AND_DISTRIBUTION_PLAN.md). One
visible app/DMG is the current direction, but it does not select an installation authority,
minimum OS, updater, replacer, or protected-container bootstrap mechanism.

## Installation identity

Each installation has:

- a random opaque `installationId` used as the normative internal identity;
- one hardware-backed installation root public key where supported;
- purpose-separated operational public keys and authorizations;
- optional DID representations for interoperability/export;
- a signed `InstallationManifest`;
- a sequence-ordered trust-epoch chain;
- a verified external-trust checkpoint and bounded local `TrustSnapshot`.

Loss of the installation root creates a new installation identity. Ordinary backup/restore does not
silently migrate private installation keys.

## InstallationManifest

The signed manifest records at least:

- schema version and manifest sequence;
- installation ID and installation-root public-key identity;
- platform, architecture, and distribution authority;
- expected daemon, Broker, Supervisor, updater, optional Guardian, and—only if the TypeScript
  profile is admitted—Source Preparer roles;
- signing identifiers, team IDs, peer code requirements, exact active code-directory hashes, and
  relevant entitlements;
- Approval Broker and Supervisor evidence-key authorizations;
- policy-bundle and runtime-profile-registry digests;
- pinned TUF root identity and accepted metadata checkpoint;
- component and storage-format versions;
- the Supervisor private state-root identity, closed store/owner-lock names, and Proposed
  ADR-0033 owner-lock UID/device/inode/mode/link-count enrollment;
- previous manifest/epoch digest;
- transition reason: install, update, repair, recovery, or authority change.

The manifest is configuration evidence, not proof that signed code is logically correct.

## Trust epochs

Every successful component- or active-policy-changing transition creates an immutable trust epoch.
Trusted component IPC, plan registration, approval, attempt creation, and receipts bind the same
epoch number and digest.

Epoch checks detect:

- partial component update;
- old daemon with new Supervisor or the reverse;
- stale Broker or updater process;
- restored policy/profile state;
- enrolled component replacement;
- many downgrade or partial-state rollback conditions.

The chain is **sequence-ordered**, not inherently monotonic. If an attacker with privileged local
state access coherently restores an older manifest, keys, stores, and checkpoint, local signature
verification alone may accept that older world. Stronger rollback evidence requires at least one of:

- a platform non-rollbackable counter/anchor with suitable semantics;
- an independently protected latest checkpoint;
- a privacy-reviewed external witness or transparency service.

Receipts must state which mechanism, if any, was active.

## Storage and access

| Data | Owner/location | Required property |
| --- | --- | --- |
| Root and operational key references | Component-specific Keychain/access group | Daemon excluded; user presence on trust-changing/approval keys |
| Current epoch checkpoint | Supervisor protected state plus narrow Keychain checkpoint where feasible | Atomic update and rollback comparison |
| Complete manifest/epoch chain | Supervisor-owned store | Hash-linked, bounded, recoverable |
| Verified public trust metadata | Updater/trust cache | Pinned roots and rollback/freshness checkpoint |
| Proposed TypeScript source set/store | Source Preparer-only protected container | Exact single-member access; sealed installer-owned genesis/update descriptor; installation/epoch/store-format binding; no shared app group or live-path handoff |
| Optional witnessed digest | External or independent protected store | No job content; privacy-reviewed correlation |

A general shared app group is not used merely for convenience. Cross-component data moves through
authenticated typed IPC or narrow handles.

Proposed ADR-0033 currently assigns one-time private state-root and owner-lock creation to the
trusted containing application/installer. The installation review identifies a competing
component-private-container composition: the authenticated setup ceremony authorizes creation,
the Supervisor creates the exact objects inside its own container, and a separately authorized
bootstrap role enrolls the returned closed identity projection. That alternative is not selected
by documentation alone. The signed protected-container spike must determine the supported owner
and amend ADR-0033 before product code.

Whichever bootstrap composition is selected, ordinary updates preserve the same root and lock
object. Normal Supervisor startup, the daemon, and store openers never create or replace it. Loss,
relocation, or restore to a different inode is repair-required and needs an authorized forward
epoch/new-installation decision. The BSD advisory lock serializes cooperating Supervisors;
installed protected-directory enforcement is what must deny baseline same-UID path replacement.
Mode `0600` is not that containment boundary.

Pairwise App Groups, if selected for sandboxed Mach/XPC naming, remain real shared-container and
Keychain namespace capabilities. Capsule places no authority, key, file, defaults, source,
content, or migration state in those groups and must retain installed negative evidence. An empty
group by policy is not structural absence of the capability and is not peer authentication.

A stable data-protection Keychain access group is a Team/profile/entitlement boundary, not an
exact-build or trust-epoch boundary. Gate B demonstrated that a stale same-team Broker rejected by
an enrolled code-directory hash could still read a new group item and use a new Secure Enclave key
in its historical group. Exact XPC requirements do not revoke Keychain membership.

The preferred v0 design therefore creates a fresh group and non-migrated key for every
identity-changing security epoch. Old/new cross-use denial, create-if-absent fingerprint binding,
logical replacement, physical retirement, rollback/forward-repair boundaries, component
acceptance, and modeled/provisioned process death passed. Shipping adoption remains conditional on
the installed Developer ID package, profile capacity, locked-Keychain, restore, migration, and
power-loss matrices. See ADR-0021.

## Development builds

Unsigned, ad-hoc, locally signed, or unrecognized builds default to local `development`
distribution posture. A custom pinned trust repository can enroll custom signed builds for testing,
but receipts identify the custom authority and never call them official Capsule production builds.

Debugged or signature-invalid processes cannot claim validated runtime-integrity posture.

## TrustSnapshot

A network-capable updater/trust verifier produces a compact signed local snapshot containing only
the release/profile/review/validation/revocation state required by local policy. It binds:

- pinned TUF root identity and accepted role versions;
- snapshot/timestamp freshness classification;
- active component release identities;
- active runtime bundle and review/validation identities;
- Capsule-defined revocation/disable records;
- local organization/administrator policy checkpoint;
- creation and expiry/freshness bounds;
- prior local snapshot checkpoint where applicable.

The Supervisor parses this bounded object rather than general TUF metadata during execution.

## Transition states

```text
stable
  → preparing-update
  → prepared
  → swapping
  → pending-verification
  → finalizing-epoch
  → awaiting-component-acceptance
  → stable

failure paths:
  → repair-required
  → quarantined
```

No attempt starts while a component-changing transition is incomplete. The daemon cannot clear
`repair-required` or authorize an epoch transition.

## Key lifecycle

- An identity-changing security epoch creates a fresh component access group, fresh non-migrated
  Secure Enclave key, and exact public-key authorization. The old key is logically rejected at the
  epoch commit, physically retired before attempts re-enable, and never replaced inside a stale
  component's historical group.
- Installation-root use is rare, explicit, and auditable.
- Approval keys require fresh user presence for each v0 approval.
- Supervisor evidence keys are background-usable only by the enrolled Supervisor identity.
- A revoked or replaced key never regains authority because a stale daemon cached it.
- `did:key` identity changes when its key changes; Capsule replacement records provide local
  continuity rather than pretending the old DID was updated.

## Recovery limitations

An interrupted transition must not fall back to ordinary execution because some binaries happen to
start. A trusted installer can restore exact current components or authorize a new epoch, but does
not reset the grant ledger, trust history, or pending backend reconciliation.

See [Update and Recovery](../UPDATE_AND_RECOVERY.md) for the transition protocol.
