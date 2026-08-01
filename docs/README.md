# Documentation

Capsule is a security boundary under design. Read these documents before changing protocol,
policy, execution, identity, content, runtimes, backends, updates, or evidence:

1. [Project definition](PROJECT.md) — scope, principles, agreed direction, and success criteria
2. [Architecture](ARCHITECTURE.md) — daemon/Broker/Supervisor authority and storage boundaries
3. [Technical design](TECHNICAL_DESIGN.md) — integrated v0 protocol, trust, lifecycle, and evidence
4. [Threat model](security/THREAT_MODEL.md) — adversaries, invariants, attack surface, and severity
5. [Feasibility spikes](FEASIBILITY_SPIKES.md) — disposable evidence gates before schema freeze
6. [Roadmap](ROADMAP.md) — risk-reduction order and phase exit evidence
7. [Architecture decisions](adr/README.md) — accepted decisions and historical supersession

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
