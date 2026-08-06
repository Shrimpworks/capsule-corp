# Documentation

Capsule is a security boundary under design. Read these documents before changing protocol,
policy, execution, identity, content, runtimes, backends, updates, or evidence:

1. [Project definition](PROJECT.md) — scope, principles, agreed direction, and success criteria
2. [Security overview](SECURITY_OVERVIEW.md) — approachable explanation of the layers, mechanisms,
   evidence, and current limitations
3. [Architecture](ARCHITECTURE.md) — daemon/Broker/Supervisor authority and storage boundaries
4. [Technical design](TECHNICAL_DESIGN.md) — integrated v0 protocol, trust, lifecycle, and evidence
5. [Threat model](security/THREAT_MODEL.md) — adversaries, invariants, attack surface, and severity
6. [Feasibility spikes](FEASIBILITY_SPIKES.md) — disposable evidence gates before schema freeze
7. [Roadmap](ROADMAP.md) — risk-reduction order and phase exit evidence
8. [Architecture decisions](adr/README.md) — accepted decisions and historical supersession
9. [Related systems and design influences](RELATED_SYSTEMS.md) — public precedents, reusable
   lessons, rejected patterns, and the boundary between comparison and Capsule evidence
10. [Ecosystem reuse and adoption map](ECOSYSTEM_REUSE_AND_ADOPTION.md) — roadmap-specific
    build-vs-adopt decisions, dependency authority/TCB classification, bounded spikes, and the
    mandatory dependency-policy checklist
11. [Work status language](STATUS_LANGUAGE.md) — the required `PASSED`, trended `IN_PROGRESS`,
    `BLOCKED`, and true-path-abandonment `NO_GO` vocabulary
12. [Experiment archive](EXPERIMENT_ARCHIVE.md) — immutable migration identity, evidence-linking
    rules, and the boundary between product conformance fixtures and disposable spike code
13. [macOS installation and distribution plan](MACOS_INSTALLATION_AND_DISTRIBUTION_PLAN.md) — one-
    application packaging direction, role boundaries, setup/update/repair/uninstall sequencing,
    staged MVP scope, blockers, and the focused Apple-platform research brief
14. [macOS I2A protected-root bootstrap decision](MACOS_INSTALLATION_I2A_PROTECTED_ROOT_BOOTSTRAP_DECISION.md)
    — selected Coordinator/Supervisor authority split, signed request/record contract,
    descriptor-relative ordering, fault oracles, and exact I2B slices
15. [macOS I2B3 signing preflight and stale-profile blocker](MACOS_INSTALLATION_I2B3_SIGNING_PREFLIGHT_AND_STALE_PROFILE_BLOCKER.md)
    — exact Team-3DDR profile/signing readback, mandatory stale-profile mutation stop, cleanup,
    and the architecture decision required before protected-root execution may resume
16. [Apple certificates, credentials, identifiers, entitlements, and Capsule keys](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md)
    — canonical Team-ID decision, environment/component matrices, safe setup and verification,
    storage/rotation policy, redacted inventory, and Dylan's next actions
17. [Internal-alpha architecture and release audit](ALPHA_ARCHITECTURE_AND_RELEASE_AUDIT.md) — the
    reconciled owner-only alpha scope, ranked blockers, risks, minimum hostile corpus, and external-
    alpha boundary
18. [Alpha vertical-flow readiness](ALPHA_VERTICAL_FLOW_READINESS.md) — exact existing
    contracts/mechanics, missing consumers, corrected authority boundaries, and ordered product
    closure
18. [Broker rendering and approval conformance v0](BROKER_APPROVAL_CONFORMANCE_V0.md) — passive
    Supervisor-owned projection, public-key-only COSE verification, exact claim boundary, and
    installed-product blockers
19. [Public-key approval to FakeBackend integration](APPROVAL_FAKE_LIFECYCLE_INTEGRATION.md) —
    unwired public-fixture verification, durable one-attempt creation, response-loss/restart
    convergence, and `AttemptID`-only no-guest fake lifecycle evidence

For live per-item status, see the
[current workstream dashboard](STATUS_LANGUAGE.md#current-workstream-dashboard). The paragraphs
below are a reading guide to the underlying decision and evidence documents, not a status
restatement.

The completed Gate C tracks are synthesized in
[Gate C implementation readiness](GATE_C_READINESS_CHECKPOINT.md). Their independent review and
the exact pre-user-byte branch point are recorded in
[Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md).

The first backend-independent implementation record is
[Phase 2A contract foundation](PHASE_2A_CONTRACT_FOUNDATION.md). Its three independent contract,
migration, and conformance reviews are consolidated in
[Phase 2A parallel-review synthesis](PHASE_2A_PARALLEL_REVIEW_SYNTHESIS.md).
The proposed exact validator and registration rules for the next slice are in
[Phase 2B boundary decisions](PHASE_2B_BOUNDARY_DECISIONS.md).
The next Supervisor retention boundary is defined in
[Proposed ADR-0031](adr/0031-checkpoint-closed-supervisor-cohorts.md) and its
[archive/compaction conformance plan](SUPERVISOR_ARCHIVE_COMPACTION_PLAN.md). The documentation-only
[F6 SQLite research and execution checkpoint](SUPERVISOR_ARCHIVE_F6_SQLITE_EXPERIMENT_PACKET.md)
freezes one separately authorized comparator packet without selecting a dependency or changing
F6/product-admission status.

The cross-phase provenance, task-to-evidence mapping, merged integration checkpoints, and current
handoff are maintained in the [workstream and evidence ledger](WORKSTREAM_EVIDENCE_LEDGER.md).

The governed runtime's first composed-profile slice is the passive
[C1 controlled-development composition contract](protocol/GOVERNED_DENO_CORE_C1_COMPOSITION.md).
The exact unsigned governed Linux/arm64 candidate and its publication/admission boundary are in the
[governed runtime release-candidate contract](GOVERNED_RUNTIME_RELEASE_CANDIDATE.md).
It fixes the intended `.mjs` JSON-in/JSON-out surface and exact governed construction identities
without creating a guest or admitting a runtime. The follow-on passive
[C2A execution-profile contract](protocol/GOVERNED_DENO_CORE_C2A_EXECUTION_PROFILE.md) freezes the
refusing descriptor, candidate machine, transport, teardown, artifact-blocker, known-answer, and
C2B evidence profile without executing it. C2B owns the first separately authorized composed
execution evidence and remains blocked.
The follow-on
[C2B fixed-fixture passive binding](protocol/GOVERNED_DENO_CORE_C2B_PASSIVE_BINDING.md) freezes the
unchanged C1/C2A relationship to one exact historical governed build candidate. Its immutable
[v2 successor](protocol/GOVERNED_DENO_CORE_C2B_PASSIVE_BINDING_V2.md) binds reviewed merged no-guest
artifact closure under a new identity while keeping final execution/admission inputs null. Neither
wires a consumer or changes C2B execution or admission status. The
[v3 passive successor](protocol/GOVERNED_DENO_CORE_C2B_PASSIVE_BINDING_V3.md) replaces those nulls
with closed roles, omitted unsupported resource fields, exact current governed source identities,
and a composed passive-contract digest while preserving typed blockers for current-source libkrun
and final runner artifacts. The immutable
[v4 materialized successor](protocol/GOVERNED_DENO_CORE_C2B_MATERIALIZED_PROFILE_V4.md) closes only
those build/static blockers with exact accepted header/current-source libkrun, independent ABI
review, unsigned final runner bytes, and a new composed digest. It also wires no consumer,
executes no runner or guest, and authorizes no guest.

Accepted [ADR-0034](adr/0034-freeze-mjs-first-release-contract.md) puts the first-release
authenticated local IPC S1 fixture slice on the single-member `.mjs` plan-v0 path. M1 source/
manifest fixtures precede S1/M2 registration/fetch fixtures. The retained
[S1 consistency stop](AUTHENTICATED_LOCAL_IPC_S1_CONSISTENCY_STOP.md) still records why a
conditional later TypeScript cutover may not reinterpret the 562-byte v0 binding record or treat
the 626-byte arithmetic as a layout.
The follow-on [passive S1 contract](AUTHENTICATED_LOCAL_IPC_S1_PASSIVE_CONTRACT.md) freezes the two
logical method records and envelopes, exact caps and refusal/copy/response-loss oracles, while
leaving native transport fields and peer authentication explicitly blocked.
Accepted [ADR-0044](adr/0044-select-private-xpc-internal-alpha-cli-adapter.md) and the
[passive internal-alpha product-adapter contract](INTERNAL_ALPHA_PRODUCT_ADAPTER_PASSIVE_CONTRACT.md)
select exactly one private-XPC `SubmitMainMJSV0` CLI-to-daemon adapter and freeze bounded aggregate
flow across submission and registration/fetch. The logical contract and in-process refusal oracles
are `PASSED`; no endpoint, signing, peer-authentication, protected state, consumer, or product
admission exists, and diagnostic HTTP remains read-only.

The follow-on grammar/process decision is [Accepted ADR-0035](adr/0035-select-disposable-mjs-source-validator.md)
and its [implementation, conformance, and fault plan](MJS_SOURCE_VALIDATOR_IMPLEMENTATION_PLAN.md).
It selects a one-shot disposable Source Validator and exact Oxc candidate from retained parse-only
evidence. Its first [passive v0 contract slice](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_CONTRACT.md)
fixes bytes and cross-language test oracles. The [bounded V1 artifact](https://github.com/Shrimpworks/capsule-experiments/blob/0944ffd8cfd01ec23e4ae99138b0931d56804077/experiments/completed-compiled-artifact-payloads/payloads/capsule-corp/artifacts/mjs-source-validator-v1/README.md)
retains exact Oxc bytes and supply-chain evidence but has only an identity-free linker ad-hoc
signature; it is not installation-signed, enrolled, wired, or confined. The
[V2 process-profile checkpoint](https://github.com/Shrimpworks/capsule-experiments/blob/0944ffd8cfd01ec23e4ae99138b0931d56804077/experiments/completed-compiled-artifact-payloads/payloads/capsule-corp/artifacts/mjs-source-validator-v2/README.md) retains fixed local
launch/fault mechanics and the exact macOS stop: unsupported `RLIMIT_AS`, ambient authority without
a sandbox, and supported App Sandbox child entitlements changing the fixed V1 bytes. V2 and the
parent product work are `BLOCKED`; no product validator or runtime no-loader boundary is implemented.
The [supported macOS replacement review](MJS_SOURCE_VALIDATOR_MACOS_PROFILE_REPLACEMENT.md) and
[Accepted ADR-0036](adr/0036-select-role-separated-source-validator-launchers.md) pass R0's
architecture scope. Direct inherited helpers are rejected. The selected replacement uses two
role-specific private App-Sandboxed launchers, accepts each private container only as residual
scratch with mandatory cleanup, and uses a later evidence-derived reactive footprint watermark
without a hard-peak or host-availability claim. The next work is the
[passive v1 boundary](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md), followed by unsigned
construction, separately authorized signing/install, confinement/resource/residue evidence,
daemon consumer, Broker consumer, then M2/S1 checkpoint. No service or product consumer exists.

The proposed TypeScript approved-byte cutover is split between the
[atomic cutover plan](TYPESCRIPT_APPROVED_BYTE_CUTOVER_PLAN.md) and the selected, not-yet-
implemented [Source Preparer implementation, conformance, and fault plan](TYPESCRIPT_SOURCE_PREPARER_PLAN.md).
Proposed ADR-0032 assigns transformation and immutable original/emitted/object custody to a
separately enrolled planning component only if TypeScript is later selected; it does not activate
a transformer, endpoint, runtime, or execution path and is not a first-release dependency.

## Detailed authority and trust documents

- [Trust architecture](security/TRUST_ARCHITECTURE.md)
- [Installation trust](security/INSTALLATION_TRUST.md)
- [Runtime integrity](security/RUNTIME_INTEGRITY.md)
- [Component compromise matrix](security/COMPONENT_COMPROMISE_MATRIX.md)
- [Control evidence matrix](security/CONTROL_EVIDENCE_MATRIX.md)
- [Execution Supervisor](EXECUTION_SUPERVISOR.md)
- [Trust repositories](TRUST_REPOSITORIES.md)
- [Update and recovery](UPDATE_AND_RECOVERY.md)
- [macOS installation and distribution plan](MACOS_INSTALLATION_AND_DISTRIBUTION_PLAN.md)
- [macOS I2A protected-root bootstrap decision](MACOS_INSTALLATION_I2A_PROTECTED_ROOT_BOOTSTRAP_DECISION.md)
- [macOS I2B3 signing preflight and stale-profile blocker](MACOS_INSTALLATION_I2B3_SIGNING_PREFLIGHT_AND_STALE_PROFILE_BLOCKER.md)
- [Apple certificates, credentials, identifiers, entitlements, and Capsule keys](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md)
- [Protocol object model](protocol/OBJECT_MODEL.md)

The top-level JSON Schemas and current mixed `Job` TypeScript types remain canonical for the
buildable scaffold and tests, but they are explicitly **pre-freeze**. Passive Phase 2A proposal and
internal-object candidates are verified separately and are not activated target contracts. The
blocking evidence and contract decisions determine honest semantics before a coordinated
schema/type/example/API replacement. See [Schema status](../schemas/README.md).

## Engineering and security-coding references

- [Secure coding standards](security/SECURE_CODING_STANDARDS.md) — cross-language engineering
  practice (closed-world validation, bounded parsing, no ambient authority, provenance over shape,
  dependency supply chain) that every language in this repository follows; a companion to
  `AGENTS.md` and the [threat model](security/THREAT_MODEL.md), not a replacement.
- [Go engineering standards](GO_ENGINEERING_STANDARDS.md) — naming, structure, error handling,
  testing, and lint hygiene for Go code in this repository.
- [JavaScript/TypeScript engineering standards](JAVASCRIPT_TYPESCRIPT_ENGINEERING_STANDARDS.md) —
  the same bar for `packages/*` and `scripts/*`, including the provenance/immutability pattern used
  throughout the protocol package.
- [Rust engineering standards](RUST_ENGINEERING_STANDARDS.md) — the bar for the Rust Source
  Validator artifacts under `artifacts/mjs-source-validator-*`.
- [Capsule domain primer](CAPSULE_DOMAIN_PRIMER.md) — fast-orientation vocabulary and Go package
  map for a contributor about to write or review Go code here.

Development setup is documented in [Development](DEVELOPMENT.md). The sequenced fresh-machine
setup, Apple identity/credential checklist, and machine-loss recovery drill are in
[Development machine setup and recovery](DEV_MACHINE_SETUP_AND_RECOVERY.md). GitHub configuration
and public release checks are documented in [Repository setup](REPOSITORY_SETUP.md).
