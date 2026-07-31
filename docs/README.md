# Documentation

Capsule is a security boundary under design. Read the documents in this order before changing
protocol, policy, execution, identity, inputs, outputs, runtimes, or backends:

1. [Project definition](PROJECT.md) — product scope, principles, and success criteria
2. [Technical design](TECHNICAL_DESIGN.md) — agreed v0 trust, execution, and implementation design
3. [Architecture](ARCHITECTURE.md) — component and trust-boundary map
4. [Threat model](security/THREAT_MODEL.md) — adversaries, mandatory properties, and security claims
5. [Roadmap](ROADMAP.md) — ordered delivery and evidence gates
6. [Architecture decisions](adr/README.md) — accepted decisions and their consequences

The JSON Schemas in `../schemas/` remain the canonical wire contracts. Design text that is not yet
represented in a schema is a Phase 0 requirement, not implemented behavior.

Development setup is documented in [Development](DEVELOPMENT.md). GitHub configuration and public
release checks are documented in [Repository setup](REPOSITORY_SETUP.md).
