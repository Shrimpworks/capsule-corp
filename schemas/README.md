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
`pnpm generate:conformance` and verify it with `pnpm verify:schemas`. The passive Source Validator
subcorpus freezes fixed frames and test oracles only; it launches no parser or process.

[`conformance/typescript-approved-byte-v0/`](conformance/typescript-approved-byte-v0/) is a
separate nine-fixture, fourteen-mutation passive corpus for Proposed ADR-0030. It has no consumer
and does not change the active scaffold or the existing v0 conformance manifest.

[`conformance/c1-governed-deno-core/`](conformance/c1-governed-deno-core/) contains the exact
passive governed-runtime C1 controlled-development composition. Its separate schema and Go/Node
known-answer validators pin intended app/runtime composition and non-admitted construction
identities; they create no process, guest, consumer, or admission state.

[`conformance/c2a-governed-deno-core/`](conformance/c2a-governed-deno-core/) contains the exact
passive C2A execution-profile preparation. It consumes the unchanged C1 fixture, closes the
numeric descriptor, candidate machine, transport, teardown, and C2B evidence contracts, and
explicitly refuses while final runnable artifact identities and exactly enforceable resource
fields remain unresolved. Its schema and independent Go/Node validators create no process, guest,
runtime consumer, credential, release, or admission state.

[`conformance/c2b-governed-deno-core/`](conformance/c2b-governed-deno-core/) contains the exact
8,221-byte passive C2B binding plus byte-for-byte mirrors of its two draft-head input records. The
binding preserves C1/C2A, fixes one source/input/completion known answer and six artifact identities,
and records the merge/reverification gate. Independent Go and TypeScript validators consume only
these passive bytes; no product consumer, runtime, profile, backend, VM, guest, or admission exists.

[`conformance/governed-deno-core-release-candidate/`](conformance/governed-deno-core-release-candidate/)
contains the exact unsigned Linux/arm64 candidate consumption manifest and bounded mutation corpus.
Its self-digest, schema, and offline verifier bind public merged evidence without copying experiment
code or large artifacts; publication, execution, selection, and admission remain inactive.

[`authority/`](authority/) contains the closed passive field-authority manifest, its JSON Schema,
and exact coverage notes. Repository verification compares its 713 field entries with 45 selected
current JSON Schema, CDDL, and Go targets and fails on missing, unknown, duplicate, stale, or
nonexistent classifications. This is a pre-freeze development invariant, not schema admission.
