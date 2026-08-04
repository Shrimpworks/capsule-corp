# macOS installation I0 passive contract

```text
Work item: macOS installation Slice I0 passive contracts
Status: PASSED
Scope: one visible no-guest application profile, seven exact roles, entitlement/service/bootstrap
  projections, deterministic inventory/update/repair/uninstall validators, generated fixtures,
  and field-authority mappings only
Evidence or reason: canonical Go values and generated JSON known answers reject missing, mixed,
  extra, inactive-signing, interrupted-bootstrap, mixed-update, unsafe-repair, and unsafe-uninstall
  cases without any external effect
Remaining work: I1-I6 and every signed installed, protected-root, IPC, replacement, distribution,
  runtime, backend, and guest gate remain outside I0
Next action: I1 may construct only the exact execution-disabled developer app tree after matching
  signing material and supported Apple bundle/service placement are available
Parent status: macOS installation and distribution are IN_PROGRESS — TRENDING_GOOD
```

Decision: [Accepted ADR-0037](../adr/0037-freeze-passive-macos-installation-i0-contract.md).

The canonical machine-readable Go test profile is
`internal/installation/macosplan/testdata/profile.json`; its manifest records exact sizes and
SHA-256 digests. The Go contract and generator are in `internal/installation/macosplan`.

## Boundary

I0 performs no signing, provisioning, Keychain access, application construction, installation,
Service Management registration, XPC activation, process launch, state creation, update, repair,
uninstall, runtime, backend, or guest operation. Every profile capable of causing one of those
effects is inactive and activation-refusing.

The profile composes:

- ADR-0029's native-fronted Go Supervisor identity, two services, and four-call role split;
- ADR-0033's exact owner-lock object policy and no-create startup order;
- ADR-0036/R1's two containing bundles, role-private services, methods, and parser identities;
- R2's unsigned launcher and parser resource paths without treating them as enrolled bytes; and
- G2 as current-v1/no-guest local evidence while retaining G3's Team/profile/bootstrap blockers.

## Closed refusal order

Bundle evaluation validates the canonical profile and exact ordered role tree first. Fewer roles
refuse as `component-missing`, more roles as `component-extra`, and any role/path/bundle/service/
signing substitution as `component-mixed`. Only an exact inventory reaches profile activation,
where the checked-in profile refuses as `signing-profile-inactive`.

Bootstrap encodes a sole clean recovery edge but I0 always refuses a caller-asserted active profile,
so it cannot produce attempt enablement. Early enablement, mixed inventory, missing protected
state, inactive bootstrap authority, wrong epoch, or unclean recovery enters `repair-required` or
a fixed trust-state refusal.

Update comparison binds exact role/service/entitlement-set digests, bootstrap and store identities,
Supervisor IPC v0, and validator v1. There is no old-or-new compatibility window. Repair and
uninstall outputs are classifications only and never authorize mutation.

## Generated corpus

The corpus contains 19 cases across inventory, bootstrap, update, repair, and uninstall plus the
canonical profile and digest manifest. Go regenerates it with:

```sh
go generate ./internal/installation/macosplan
```

Focused tests re-render every file and fail on byte drift. The fixtures are language-neutral, but
I0 adds no TypeScript decoder because no current TypeScript consumer needs installation authority.

## Dependencies after I0

| Slice | Required dependency before work can pass |
| --- | --- |
| I1 | matching Apple Development Team/profile set; supported exact `Capsule.app` helper and private-XPC placement; Swift app shell and daemon/Supervisor targets; execution always disabled |
| I2 | selected protected-root bootstrap owner; signed bootstrap envelope; Supervisor-private container denial corpus; descriptor-relative store open; matching identities; product store decision |
| I3 | I2; ADR-0029 pairwise App Group/private-service decision and signed peer matrix; ADR-0036 R3 signing/install plus R4 confinement/resource/residue evidence |
| I4 | I2/I3; reviewed whole-bundle replacement authority, service stop/re-registration/recovery, exact mixed-version refusal, and forward repair |
| I5 | I4; pinned TUF/update-verifier profile and separately reviewed mechanical replacer; not an MVP dependency |
| I6 | Developer ID/notarization/stapling/Gatekeeper; support-floor and clean-host matrix; backup/restore; final uninstall/local-erasure and distribution admission |

Runtime/backend admission remains separate and cannot be inferred from any I0-I6 result.
