# Trust Architecture

Status: intended v0 trust model; implementation evidence is pending.

## Objective

Capsule separates public planning, trusted user authorization/content custody, and hostile-guest
execution so that no ordinary component compromise automatically yields all three powers.

The local authorization root is the Installation Trust Domain: a random installation identity,
locally authorized public keys, component identities, policy, and a sequence-ordered trust epoch.
DIDs can name principals and verification methods, but do not grant authority.

## Authorities

| Authority | Owner | Positive authority | Explicit denial |
| --- | --- | --- | --- |
| Proposal/planning | Agent-facing daemon | Validate proposals, resolve trusted inputs/policy/profiles, construct and register a plan | Approval signing, user-content custody, backend launch, grant reset |
| Proposed TypeScript source preparation | Enrolled Source Preparer | Emit and retain exact original/executable/profile/options/record bytes before plan construction | Approval, plan registration, user-content custody, backend launch, post-registration transformation |
| Human authorization | Approval Broker | Fetch/render registered plan, require user presence, sign one attempt-bound approval | Plan creation, backend launch, enforcement claims |
| Content custody | Content Broker | Select/snapshot user files, own user content, transfer attempt-scoped handles, release output | Agent endpoint, general network, backend launch |
| Execution enforcement | Execution Supervisor | Register exact bytes, enforce hard safety, consume grants, create/destroy guests, sign transcript | Public agent parsing, file picker, rich parsing, network trust resolution |
| Installation bootstrap | On-demand Trust Coordinator | With fresh user presence, construct and installation-root-sign one closed Supervisor bootstrap request and record | Supervisor-state creation, Approval/evidence signing, update/replacement, backend launch |
| Release/profile distribution | Updater/trust verifier | Verify TUF metadata and produce a bounded local trust snapshot | Per-job approval, backend launch, user-content access |
| Optional observation | Runtime Guardian | Report relevant Endpoint Security events | Approval or launch authority |

The v0 Trusted Host Broker may combine Human authorization and Content custody in one signed native
process. Their keys, APIs, and persisted records remain purpose-separated.

Proposed ADR-0032 splits the narrow TypeScript parser/store responsibility from the public daemon.
The Source Preparer is additional planning and approval-understanding TCB, but receives no
operational key or execution authority. The proposal remains unimplemented.

## Why process separation is conditional

On macOS, processes running as the same user can often interact with the same user-owned files and
services unless the platform imposes stronger controls. A design diagram does not create a security
boundary.

The v0 separation depends on feasibility evidence for:

- XPC peer code requirements that bind the expected signing identifier/team and component purpose;
- effective user/session checks and protocol-level installation/epoch binding;
- protected app/data containers and absence of broad shared app groups;
- Keychain access groups and access-control policies that exclude the daemon;
- exact active build, entitlement, and dynamic code validation;
- a backend control endpoint that only the Supervisor can reach.

The local administrator and kernel remain trusted. These controls primarily contain an untrusted
agent, hostile guest, and compromised ordinary same-user process within the documented assumption.

Gate B found a deliberate boundary between IPC and key custody: exact-build XPC requirements can
reject a stale component, while a stable data-protection Keychain group still recognizes that
component's historical Team/profile/entitlement. Therefore operational keys are not considered
epoch-isolated merely because their component channel is exact-build authenticated.

The proposed v0 mitigation is a fresh access group and fresh non-migrated key for every identity-
changing security epoch. Modeled and provisioned process-death tests support its fence, exact
fingerprint authorization, old-key retirement, rollback/forward-repair, component-acceptance, and
re-enable ordering. Shipping claims still require installed distribution/profile/Keychain/restore
validation. See ADR-0021.

## Local key hierarchy

| Key | Custodian | Authorized purpose | Availability |
| --- | --- | --- | --- |
| Installation root | On-demand Trust Coordinator ceremony selected by Proposed ADR-0038 | Enroll operational keys and authorize trust transitions | Rare; fresh user presence in v0 |
| Approval key | Approval Broker | `capsule.plan.approve` | Fresh user presence per v0 plan |
| Supervisor evidence key | Execution Supervisor | `capsule.execution.attest` | Noninteractive, narrow process use |
| Optional Broker content key | Content Broker | `capsule.content.attest` | Noninteractive only if cross-process claims require it |
| Transport/delegation keys | Deferred | Remote identity or explicit delegation | Not part of local v0 |

A key authorization binds installation, public key/key ID, purpose, issuer, validity, status,
sequence, replacement/revocation relationship, and allowed object types. A signature is rejected if
any binding fails even when its mathematics is valid.

The installation root is not available to the visible app, daemon, Approval/Content Broker,
Supervisor, updater, or bundle replacer and is not used for routine receipts. Apple code identity
authenticates installed participants but never substitutes for this Capsule signing authority.
Receipt or evidence signing authority never substitutes for user approval.

## DIDs

Normative rule:

> A DID identifies a principal or verification method. Capsule trust metadata and policy decide
> what that principal may do.

First-class uses include:

- a public Capsule organization or enterprise identifier;
- independent profile reviewer or hosted worker identifiers;
- portable signer references in exported receipts;
- offline `did:key` rendering of an operational public key;
- discovery hints for public trust/transparency endpoints outside execution.

Prohibited assumptions:

- possession or resolution of a DID is not authorization;
- a DID document is not inherently trusted;
- `did:key` does not supply rotation, deactivation, or recovery;
- `did:web` or another method does not replace a pinned TUF root;
- no network resolution, arbitrary method, resolver plugin, or remote JSON-LD context enters local
  v0 approval or execution.

The normative internal installation identity is therefore an opaque `installationId` plus enrolled
keys. Optional DID fields are preserved in the target object model so interoperability does not
require retrofitting identity later.

## External trust

External release and profile trust begins with locally pinned TUF root metadata. An updater/trust
verifier checks signatures, delegated scope, versions, expiration, lengths, hashes, snapshot/
timestamp consistency, and local rollback checkpoints.

The updater reduces accepted state to a signed bounded `TrustSnapshot`. The Supervisor consumes
that local object and does not implement a general network client, DID resolver, or full live TUF
refresh path.

TUF distributes Capsule-defined key/profile revocation and emergency-disable records. Their
meaning, precedence, and execution consequence remain Capsule policy.

## Policy precedence

```text
Capsule non-overridable hard-safety invariants
  → required organization policy and pinned trust roots
  → local administrator policy
  → user defaults and ceilings
  → exact one-plan approval
  → agent request
```

Higher layers may narrow authority and may not widen a lower hard-safety invariant. The Supervisor
rechecks v0 hard-safety rules independently of the daemon's policy result.

## Trust state

Key, component, and trust-snapshot status is explicit. At minimum:

- `active`
- `suspended`
- `revoked`
- `replaced`
- `expired`
- `stale-degraded`
- `quarantined`
- `repair-required`

Unknown state fails closed. The daemon cannot change authoritative key, epoch, quarantine, or grant
state.

## Binding requirements

Security-relevant objects bind the minimum relevant subset of:

- unique object type, version, and signing purpose;
- issuer/key and local key authorization;
- installation and trust epoch;
- audience and expected peer/Supervisor;
- plan digest and registration;
- attempt identifier/nonce;
- issue/expiry or other freshness data;
- referenced content/profile/policy/backend identities.

Purpose and object-type separation prevent an Approval signature from becoming a key authorization,
profile review, content grant, or enforcement transcript.

## Trust limitations

- Valid code signing does not prove correct program logic.
- Key possession does not prove the intended person still exclusively controls the key.
- Sequence-ordered trust epochs detect mismatches but not every coherent rollback.
- Secure Enclave and Keychain reduce key-extraction risk but do not make local state immutable.
- A compromised Supervisor can launch unauthorized guests and forge its enforcement claims.
- A compromised kernel or privileged administrator can likely falsify local observations and is out
  of scope for validated-local guarantees.

See [Installation Trust](INSTALLATION_TRUST.md), [Runtime Integrity](RUNTIME_INTEGRITY.md), and
[Component Compromise Matrix](COMPONENT_COMPROMISE_MATRIX.md).
