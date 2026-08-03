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
runtime authority. Its fixed proposal resolver contexts and known-answer source/canonical-input
bytes are conformance inputs, not an activated planner or accepted ADR. Regenerate it with
`pnpm generate:conformance` and verify it with `pnpm verify:schemas`.

[`conformance/typescript-approved-byte-v0/`](conformance/typescript-approved-byte-v0/) is a
separate nine-fixture, fourteen-mutation passive corpus for Proposed ADR-0029. It has no consumer
and does not change the active scaffold or the existing v0 conformance manifest.
