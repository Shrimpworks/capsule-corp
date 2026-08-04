# Passive field-authority classification

[`field-authority-manifest.json`](field-authority-manifest.json) is the closed, machine-readable
development invariant for selected current passive targets. Its 172 entries cover 16 objects:

- the `JobProposal` v0 candidate;
- `ExecutionPlan`, `PlanRegistration`, and `ApprovalGrant` v0 candidates;
- passive `ApprovalReference`, `AttemptReference`, durable `ApprovalRecord`, and immutable
  `ExecutionAttempt` v0 structures; and
- the passive TypeScript approved-byte original/executable manifests, transformer profile,
  normalized options, transformation record and set, and future-plan source-binding projection; and
- the passive single-member MJS `SourceManifest`, including nested member path/digest/length.

Each field references one complete classification profile that fixes its origin role,
validator/resolver, authority effect, approval visibility and source, content/guest-control
status, binding, and fail-closed unknown behavior. The vocabulary is closed by
[`field-authority-manifest.schema.json`](field-authority-manifest.schema.json). Evidence-only
observations remain separate from authority-bearing fields; for example, the retained approval
envelope digest is evidence-only and never selects a ledger record.

`pnpm verify:schemas` compares the manifest directly with the canonical JSON Schema, numbered CDDL
maps, and explicitly version-marked Go passive structs. It rejects missing fields, unknown
classification values, duplicate paths/source fields, stale object identities or versions, and
classified fields absent from their target definition. Focused mutation tests retain each failure
mode.

This directory does not freeze a schema, activate a consumer, assign the unresolved production
TypeScript transformation owner, admit a runtime/backend, or authorize execution. It deliberately
does not classify the deprecated mixed `Job` scaffold or the illustrative future
`ExecutionPlan` v1 shape; those require their coordinated versioned cutovers.
