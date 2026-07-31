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
