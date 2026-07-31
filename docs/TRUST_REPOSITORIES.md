# Trust Repositories

Status: intended production distribution design; not a live v0 execution dependency.

## Purpose

Capsule needs verifiable release, component, runtime-profile, review, validation, delegation, and
emergency-disable distribution. A service URL, TLS session, DID, or JSON response is not a trust
anchor.

The default distribution uses The Update Framework (TUF) model. Users and organizations may pin
one or more self-hosted repositories with equivalent verified metadata.

## Role model

The initial trusted root metadata is pinned during install or an explicit trust ceremony. Mature
root operation should support threshold signatures and offline custody.

Standard roles provide:

- `root`: trust-root and role-key evolution;
- `targets`: authorized artifact metadata;
- `snapshot`: one consistent view of targets metadata;
- `timestamp`: bounded freshness of the current snapshot.

Delegated targets roles may cover:

- desktop/daemon releases;
- Broker and Supervisor releases;
- runtime bundle manifests;
- profile review attestations;
- backend validation records;
- Capsule-defined revocation and emergency-disable objects;
- enterprise/custom namespaces.

Delegation path/scope and threshold are security boundaries. An online role compromise is contained
only within its authorized scope.

## Separation from live execution

The network-capable updater/trust verifier owns:

- repository transport;
- TUF parsing and signature/delegation/version/expiration verification;
- rollback/freeze/mix-and-match checkpoints;
- artifact length/hash verification;
- local policy selection among pinned repositories;
- output of a compact signed local `TrustSnapshot`.

The Supervisor owns only bounded `TrustSnapshot` verification and consumption. It performs no
repository fetch, general TUF parsing, DID resolution, URL redirect handling, or remote JSON-LD
processing during plan registration or attempts.

This keeps network and complex metadata parsing outside the execution TCB and prevents service
availability from becoming per-job authorization.

## TrustSnapshot

The exact contract is frozen after the feasibility spikes. It is expected to bind:

- object type/version, installation, issuer, prior checkpoint, and validity/freshness;
- pinned root digest and accepted root version;
- accepted timestamp, snapshot, and delegated metadata versions/digests;
- active component release identities;
- active runtime bundle, review, registry, and backend-validation identities;
- Capsule-defined key/profile/component revocation or disable records;
- organization/administrator policy checkpoint;
- distribution authority classification;
- explicit offline/freshness posture and limitations.

This is a derived local authorization input, not a replacement TUF specification or a claim that
TUF defines Capsule revocation semantics.

## Offline behavior

- Local execution may use cached verified state according to explicit policy.
- Offline bundles can initialize or refresh a pinned repository.
- Expired or stale revocation state can refuse execution or downgrade requested posture.
- Network/service failure never causes acceptance of unsigned, unverified, or rollback metadata.
- A receipt records the snapshot identity and freshness classification used for its attempt.

## Multiple authorities

Organizations may require additional pinned roots or delegations. Policy defines intersection and
precedence; an additional repository cannot widen Capsule hard-safety invariants.

Repository identity is pinned root metadata plus local configuration—not merely a URL or DID.
Changing a required root or repository set is an installation trust ceremony and may require a new
trust epoch.

## DIDs

DIDs can identify a public Capsule organization, profile reviewer, enterprise, or verification
method and can provide discovery hints outside execution. They do not replace pinned root metadata
or delegated authorization.

If a DID-linked endpoint is used for discovery, the locally pinned policy still selects and
verifies the repository. Network DID resolution is never required for an already installed local
job.

## Privacy

Routine refresh does not transmit:

- job source or input/output content;
- user filenames or paths;
- content/artifact digests;
- stable per-execution identifiers;
- receipt details.

Optional checkpoint witnessing uses a separate privacy-reviewed protocol and submits only the
minimum installation/receipt checkpoint. Correlation and retention are explicit user/policy
choices.

## Failure and recovery

- Invalid signature, delegated scope, length/hash, version, or snapshot consistency fails refresh.
- Expiration obeys explicit offline-grace policy; it never resets versions.
- Interrupted refresh retains the last complete verified checkpoint until a new one commits.
- Root rotation follows TUF root-update rules and explicit local policy.
- Emergency disable is an authenticated Capsule policy object and must state scope, sequence,
  effective time, recovery, and offline behavior.

## Evidence requirements

Retain test vectors for root rotation, threshold failure, freeze, rollback, mix-and-match,
delegation escape, expired metadata, malicious mirrors, interrupted refresh, conflicting pinned
repositories, and revoked runtime/profile state.

See [Installation Trust](security/INSTALLATION_TRUST.md) and
[Update and Recovery](UPDATE_AND_RECOVERY.md).
