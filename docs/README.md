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

Accepted [ADR-0034](adr/0034-freeze-mjs-first-release-contract.md) puts the first-release
authenticated local IPC S1 fixture slice on the single-member `.mjs` plan-v0 path. M1 source/
manifest fixtures precede S1/M2 registration/fetch fixtures. The retained
[S1 consistency stop](AUTHENTICATED_LOCAL_IPC_S1_CONSISTENCY_STOP.md) still records why a
conditional later TypeScript cutover may not reinterpret the 562-byte v0 binding record or treat
the 626-byte arithmetic as a layout.

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
- [Protocol object model](protocol/OBJECT_MODEL.md)

The top-level JSON Schemas and current mixed `Job` TypeScript types remain canonical for the
buildable scaffold and tests, but they are explicitly **pre-freeze**. Passive Phase 2A proposal and
internal-object candidates are verified separately and are not activated target contracts. The
blocking evidence and contract decisions determine honest semantics before a coordinated
schema/type/example/API replacement. See [Schema status](../schemas/README.md).

Development setup is documented in [Development](DEVELOPMENT.md). GitHub configuration and public
release checks are documented in [Repository setup](REPOSITORY_SETUP.md).
