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
14. [Apple certificates, credentials, identifiers, entitlements, and Capsule keys](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md)
    — canonical Team-ID decision, environment/component matrices, safe setup and verification,
    storage/rotation policy, redacted inventory, and Dylan's next actions

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
[archive/compaction conformance plan](SUPERVISOR_ARCHIVE_COMPACTION_PLAN.md).

The cross-phase provenance, task-to-evidence mapping, merged integration checkpoints, and current
handoff are maintained in the [workstream and evidence ledger](WORKSTREAM_EVIDENCE_LEDGER.md).

The governed runtime's first composed-profile slice is the passive
[C1 controlled-development composition contract](protocol/GOVERNED_DENO_CORE_C1_COMPOSITION.md).
It fixes the intended `.mjs` JSON-in/JSON-out surface and exact governed construction identities
without creating a guest or admitting a runtime. The follow-on passive
[C2A execution-profile contract](protocol/GOVERNED_DENO_CORE_C2A_EXECUTION_PROFILE.md) freezes the
refusing descriptor, candidate machine, transport, teardown, artifact-blocker, known-answer, and
C2B evidence profile without executing it. C2B owns the first separately authorized composed
execution evidence and remains blocked.

Accepted [ADR-0034](adr/0034-freeze-mjs-first-release-contract.md) puts the first-release
authenticated local IPC S1 fixture slice on the single-member `.mjs` plan-v0 path. M1 source/
manifest fixtures precede S1/M2 registration/fetch fixtures. The retained
[S1 consistency stop](AUTHENTICATED_LOCAL_IPC_S1_CONSISTENCY_STOP.md) still records why a
conditional later TypeScript cutover may not reinterpret the 562-byte v0 binding record or treat
the 626-byte arithmetic as a layout.

The follow-on grammar/process decision is [Accepted ADR-0035](adr/0035-select-disposable-mjs-source-validator.md)
and its [implementation, conformance, and fault plan](MJS_SOURCE_VALIDATOR_IMPLEMENTATION_PLAN.md).
It selects a one-shot disposable Source Validator and exact Oxc candidate from retained parse-only
evidence. Its first [passive v0 contract slice](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_CONTRACT.md)
fixes bytes and cross-language test oracles. The [bounded V1 artifact](../artifacts/mjs-source-validator-v1/README.md)
retains exact Oxc bytes and supply-chain evidence but has only an identity-free linker ad-hoc
signature; it is not installation-signed, enrolled, wired, or confined. The
[V2 process-profile checkpoint](../artifacts/mjs-source-validator-v2/README.md) retains fixed local
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
- [Apple certificates, credentials, identifiers, entitlements, and Capsule keys](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md)
- [Protocol object model](protocol/OBJECT_MODEL.md)

The top-level JSON Schemas and current mixed `Job` TypeScript types remain canonical for the
buildable scaffold and tests, but they are explicitly **pre-freeze**. Passive Phase 2A proposal and
internal-object candidates are verified separately and are not activated target contracts. The
blocking evidence and contract decisions determine honest semantics before a coordinated
schema/type/example/API replacement. See [Schema status](../schemas/README.md).

## Go engineering references

- [Go engineering standards](GO_ENGINEERING_STANDARDS.md) — naming, structure, error handling,
  testing, and lint hygiene for Go code in this repository; a companion to `AGENTS.md`, not a
  replacement.
- [Capsule domain primer](CAPSULE_DOMAIN_PRIMER.md) — fast-orientation vocabulary and Go package
  map for a contributor about to write or review Go code here.

Development setup is documented in [Development](DEVELOPMENT.md). GitHub configuration and public
release checks are documented in [Repository setup](REPOSITORY_SETUP.md).
