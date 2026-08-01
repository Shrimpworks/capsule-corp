# Schema status

The JSON Schemas in this directory are canonical for the repository's current buildable scaffold.
They are **pre-freeze** and are not the intended final Capsule v0 protocol.

The target architecture separates `JobProposal`, `ExecutionPlan`, `PlanRegistration`,
`ApprovalGrant`, `ExecutionAttempt`, `EnforcementTranscript`, `ExecutionReceipt`, and the fixed
`AgentExecutionSummary`. It also removes unsupported future authority and agent-selected guest
paths from the v0 proposal.

Do not extend the current `Job` capability union as a shortcut. The replacement inventory and
schema-freeze gates are documented in
[Protocol Object Model](../docs/protocol/OBJECT_MODEL.md) and
[Feasibility Spikes](../docs/FEASIBILITY_SPIKES.md).

[`candidates/`](candidates/) contains the first passive `JobProposal` replacement candidate. It is
verified with an example and fail-closed cases, but no daemon, SDK, or MCP endpoint accepts it. The
mixed `Job` scaffold therefore remains the current public scaffold until an atomic consumer
cutover.

[`cddl/`](cddl/) contains candidate canonical-CBOR contracts derived from Gate A2 and Phase 2A.
Those files are internal security-object profiles, not replacements for the public JSON API, and
remain pre-freeze while ADR-0019 is Proposed.

[`conformance/v0/`](conformance/v0/) contains the closed manifest and retained byte-exact fixtures
for the proposed Phase 2B boundary rules. The corpus records pending language targets and grants no
runtime authority. Regenerate it with `pnpm generate:conformance` and verify it with
`pnpm verify:schemas`.
