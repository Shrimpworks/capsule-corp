# macOS installation I2B3 signing preflight and stale-profile blocker

Date: 2026-08-05

```text
Work item: I2B3 exact Apple Development protected-root/bootstrap evidence
Status: BLOCKED
Scope: exact Team-3DDR Coordinator/Supervisor profile creation, signed-entitlement preflight,
  and the required stale-Supervisor current-container mutation stop test on the owner's named Mac
Evidence or reason: both new profiles and frozen entitlements signed exactly, but a Supervisor
  signed with the still-valid archived I1B profile rewrote a sentinel created by the new I2B3
  Supervisor profile in the same App Sandbox container
Remaining work: select and review a version/epoch boundary that prevents or safely composes stale
  Supervisor private-container mutation before any installation-root key, service, or root exists
Blocker and owner: installation architecture; Proposed ADR-0038 and the I2B plan do not yet close
  stale-profile write authority over the stable Supervisor container identity
Next action: decide the container/signing-identity epoch strategy in an ADR, provision its exact
  profiles, and separately authorize a rerun from a residue-free preflight
Parent status: installed owner-lock G3/I2B remains BLOCKED; macOS installation remains
  IN_PROGRESS — TRENDING_BAD until the stale-profile boundary has a selected replacement
```

## Defensive and authorized scope

The owner explicitly authorized the reconciled I2B3 development-signing, service, App Group,
Keychain, container, protected-root, fault, and cleanup mutations. The run remained local to the
owner's named Mac, the exact Capsule Team `3DDR84M4JS` identities, disposable test bundles, and
Capsule-owned containers. It did not use Developer ID, notarization, a product store, ordinary
daemon/Broker IPC, a runtime, backend, VM, guest, unrelated identity, credential, or data.

The stale-profile fault was deliberately ordered before service registration, LocalAuthentication,
installation-root key creation, Coordinator launch, installation, or protected-root creation. The
task required an immediate stop if a stale signer could mutate current-root state. The sentinel is
the smallest safe pre-root oracle for that property.

## Exact profile and signing readback

The selected certificate remained the exact I1B/R3 Apple Development identity:

- Team identifier `3DDR84M4JS`;
- SHA-1 `80A4969BCD1B3926020888094B9D812A283D3793`;
- SHA-256 `D3E9FBDDBC342F747C3649B5A6FFB307A575827404E02D638C11B6B795A09629`;
- serial `2680E3A814E45A8A4AC3C2B2EF09023E`;
- validity 2026-08-04 through 2027-08-04.

The Apple Developer account created the explicit Coordinator App ID
`com.capsulecorp.capsule.trust-bootstrap.v1`. No portal `group.*` App Group resource was created:
the frozen macOS entitlement uses the documented Team-ID form
`3DDR84M4JS.com.capsulecorp.capsule.bootstrap.v0`.

| Component | Profile name | UUID | CMS SHA-256 | Exact profile binding |
| --- | --- | --- | --- | --- |
| Trust Coordinator | `Capsule I2B3 Trust Coordinator Development 3DDR` | `c0446281-ba4b-451b-8c73-9ee9d8ef97a2` | `6f3bcf42fcf208352bce23b201941222aa1269d698bdf1ddcefb5e836f0c0ad8` | explicit Coordinator App ID, selected certificate, owner-named Mac |
| Supervisor | `Capsule I2B3 Supervisor Bootstrap Development 3DDR` | `c45a058b-ffdd-4a6b-bd8c-d746772a2702` | `964f79980edf22a7280fe19e52893a1e40b0a8639d5bbe3d5dc8fdfada9c6c76` | explicit stable Supervisor App ID, selected certificate, owner-named Mac |

OpenSSL CMS verification and certificate extraction agreed with the certificate fingerprints,
Team, explicit application identifiers, device selection, and one-year development expiry above.
The profile payloads authorize the Team/application identifier and wildcard Team Keychain prefix;
effective `codesign` readback, not the profile projection alone, proved the exact narrower
entitlements.

The no-launch signed preflight `PASSED` for both components. Hardened-runtime signatures had Team
`3DDR84M4JS`, their exact signing identifiers, and only these role projections:

- Coordinator: App Sandbox, bootstrap App Group, and
  `3DDR84M4JS.com.capsulecorp.capsule.trust-bootstrap.installation-root.epoch-1`;
- Supervisor: App Sandbox, bootstrap App Group, and
  `3DDR84M4JS.com.capsulecorp.capsule.supervisor.bootstrap-anchor.epoch-1`.

`codesign --verify --strict` passed both bundles. The preflight Coordinator and Supervisor
CDHashes were `c9ff866e3b2c3ae7e24eb1ca94cc8d485e249478` and
`8649ce18817e09617914e37d918f5769ea6eaccf`, respectively. No preflight binary was launched.

## Required stale-profile stop test

The fault probe used two independently signed App Sandbox app bundles with the exact stable
Supervisor signing identifier `com.capsulecorp.capsule.supervisor`:

| Probe | Provisioning/profile state | Effective capability relevant to the test | CDHash at the mutation step |
| --- | --- | --- | --- |
| current | new I2B3 profile `c45a058b-ffdd-4a6b-bd8c-d746772a2702` | App Sandbox plus new bootstrap groups | `45af8b8a724f673f15ce0ed0fcc91e002ba0f0d6` |
| stale | archived I1B profile `7e27acf1-48fc-4c57-ad43-4b1caa9c12a8` | App Sandbox only; no I2B3 App Group or role-private Keychain group | `86e4fde749208964829a8c13fa0fef2464024138` |

Both processes reported the same platform-selected private home:
`~/Library/Containers/com.capsulecorp.capsule.supervisor/Data`. The current-profile probe created
`Library/Application Support/CapsuleSupervisor/i2b3-stale-profile-sentinel` with mode `0600` and
fixed current-profile bytes. The stale-profile probe then opened that same descriptor-relative
entry with `O_NOFOLLOW|O_CLOEXEC`, truncated it, and wrote fixed stale-profile bytes. It exited
successfully. Current-profile readback returned exit `72` and the exact changed value
`stale-i1b-profile-mutated-current-state`.

This result is not a general same-user containment claim and does not show Keychain cross-group
access. It establishes the narrower blocker: changing the provisioning profile, CDHash, App Group,
and Keychain access group does not by itself version the stable App Sandbox private container.
XPC peer requirements can reject a stale caller at a listener, but they do not prevent a separately
launched stale Supervisor identity from reaching its own stable container.

## Stop boundary and cleanup

The run stopped immediately after current-profile detection of the stale mutation. Therefore:

- no Supervisor `SMAppService` registration or approval occurred;
- no Supervisor Mach service or Coordinator private XPC service launched;
- no LocalAuthentication prompt occurred;
- no Coordinator or Supervisor component performed a Keychain access-group operation, and no
  Keychain add/update/delete occurred; a post-stop shell lookup only confirmed the planned test
  service names absent;
- no installation-root key, request, observation, record, anchor, protected root, owner, store,
  pending journal, publish intent, runtime, backend, VM, or guest existed;
- no app was installed in `~/Applications` or `/Applications`.

The current-profile probe removed only the exact sentinel and then the empty
`CapsuleSupervisor` test parent. Descriptor-relative readback confirmed that parent absent.
`launchctl` confirmed the Supervisor service absent, and both installation paths and the
Coordinator container were absent.

Launching the current App-Group-entitled Supervisor caused the platform to materialize the App
Group container. Entitled enumeration found only the platform-created `Library` hierarchy,
`Library/Application Scripts/3DDR84M4JS.com.capsulecorp.capsule.bootstrap.v0`, and
`.com.apple.containermanagerd.metadata.plist`: seven directories/metadata entries and zero
Capsule-created authority or payload entries. The design already classifies this platform metadata
as residual capability, not Capsule state. The pre-existing Supervisor container remains, with no
fixed I2B parent/root/sentinel retained.

The two downloaded public provisioning-profile files and the portal App ID/profiles are retained
as audit inputs. They contain no private signing key. The disposable probe bundles and source are
not product artifacts and must not be admitted or archived as production code.

## Architecture consequence

Installed I2B cannot resume merely by repeating the same stable Supervisor bundle identifier with
a new profile. Before another protected-root run, an ADR must select one closed approach and its
replacement/fault tests. Candidate questions include a versioned Supervisor signing/container
identity, a separately versioned Supervisor-only protected container, or a cryptographically
authenticated corruption-detection design whose stale-write and repair semantics meet the exact
threat model. Revoking or deleting a portal profile is not assumed to invalidate already-signed
bytes and is not an accepted fence without retained platform evidence.

Any successor must still prove Supervisor enabled before Coordinator, `JoinExistingSession=true`,
fresh user presence, explicit access group on every Keychain operation, exact request/record
agreement, create-once descriptor-relative no-follow behavior, repair-required fencing, complete
cleanup, and the remaining I2B4-I2B5 session/update matrix. This result advances no installed
control, product-store, runtime, backend, or guest claim.
