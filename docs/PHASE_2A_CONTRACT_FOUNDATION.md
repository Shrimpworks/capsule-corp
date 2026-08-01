# Phase 2A contract foundation

Status: implementation specification for the first Phase 2 contract slice.

## Objective

Prepare the coordinated replacement of the deprecated mixed `Job` scaffold with a narrow public
`JobProposal` boundary and define the first candidate internal `ExecutionPlan` and
`PlanRegistration` contracts without granting new authority or selecting a real isolation backend.

This slice establishes four rules in executable fixtures:

1. a proposal is untrusted desired work and cannot express effective backend authority;
2. a plan is exact resolved authority and is a different object from a proposal;
3. a registration acknowledges exact canonical plan bytes and cannot carry replacement bytes; and
4. identifiers, digests, versions, limits, and logical slots are role-specific and bounded.

The user-approved implementation sequence is contract first, registered-plan/fake-backend
lifecycle second, and inline-content ownership third. This document does not freeze approval,
attempt, receipt, runtime-bundle, backend-specific transport, or cryptographic-wrapper contracts.

## Existing decisions

- `JobProposal`, `ExecutionPlan`, and `PlanRegistration` are distinct object types.
- V0 proposals omit network, subprocess, environment, native-addon, FFI, macro, inspector,
  package-install, arbitrary image/backend, mount, socket, host-path, and guest-path authority.
- Source and input/output references use semantic source paths and logical slots, never live host
  paths or agent-selected guest paths.
- Requested values above trusted user ceilings are rejected rather than clamped. Missing values are
  resolved before plan registration and approval.
- The Supervisor receives and retains exact canonical plan bytes. Execute-time operations accept
  only a Supervisor-issued registration ID.
- Public agent objects use strict JSON. Canonically registered internal objects are candidates for
  the bounded deterministic-CBOR profile in ADR-0019; that ADR remains Proposed.

## Scope

### Included

- shared role-specific identifier, SHA-256 digest, object-version, bounded-integer, source-path,
  logical-slot, and JSON-value rules;
- a passive, closed candidate JSON Schema and TypeScript projection for `JobProposal`;
- candidate examples and verification that do not activate a public proposal endpoint;
- closed candidate CDDL and byte-exact fixtures for `ExecutionPlan` and `PlanRegistration` where the
  merged architecture already determines the field;
- Go and TypeScript projections for fields admitted by those candidate contracts;
- positive fixtures and fail-closed cases for unknown fields, unsupported powers, wrong object
  types/versions, unsafe numbers, invalid source paths, invalid slots, digest/identifier confusion,
  and replacement-plan attempts;
- documentation that distinguishes frozen public JSON, candidate internal CBOR, semantic policy
  validation, and deferred objects.

### Excluded

- `ApprovalGrant`, `TrustSnapshot`, `ExecutionAttempt`, integrity, transcript, receipt, artifact,
  or content-access contract freeze;
- production CBOR/COSE encoders, decoders, signature verification, keys, or dependency selection;
- daemon, Broker, XPC, HTTP, database, or user-interface endpoints;
- policy resolution, user-default selection, or backend capability matching behavior;
- a libkrun, Apple Containerization, OCI, gVisor, runtime-adapter, or guest execution path;
- final backend-specific numeric values or port byte layouts before the remaining P0 evidence;
- activation of both `Job` and `JobProposal` as competing public contracts;
- deletion of the old scaffold before the internal plan/registration contracts and public
  consumers can cut over atomically.

## Contract boundaries

### Public JSON

`JobProposal` is the intended object accepted from an agent. During this foundation slice it is a
passive candidate with no daemon, SDK, or MCP activation. Schema acceptance proves only structural
validity. A later planner still must apply raw-decoder limits, semantic source-path
canonicalization, trusted defaults and ceilings, profile resolution, content resolution, and policy.

The first candidate is deliberately narrower than the final inventory: inline JSON input only,
exactly one inline JSON output, fixed `primary-data` and `transformed-json` slots, an ASCII-only
relative source-path grammar, and only the independently supported wall-time request. Closed RAM,
vCPU, scratch, console, source/input transport, completion-frame, and payload limits remain trusted
profile/policy values until their exact public request semantics and caps are evidenced.

The proposal must not contain false-valued placeholders for unsupported authority. Unknown fields
fail rather than being ignored.

### Internal canonical objects

`ExecutionPlan` and `PlanRegistration` use mutually exclusive object types and candidate closed
CDDL maps. Fixture encodings are deterministic and byte exact, but product serialization is not
declared implemented until reviewed Go, Swift, and TypeScript wrappers satisfy ADR-0019.

Received canonical plan bytes remain authoritative. A decoder may inspect them but cannot replace
them with decode-and-re-encode output before digesting, registration, approval, or execution.

### Semantic validation

JSON Schema and CDDL do not decide trusted policy. The following remain semantic checks:

- an entrypoint identifies exactly one canonical source file;
- source paths are unique after semantic canonicalization;
- requested limits are permitted by user policy and exactly supported by the selected profile;
- logical slots have unique roles and resolve to exact content identities;
- plan references resolve to active, purpose-correct installation, epoch, policy, runtime, and
  backend records;
- registration expiry and sequence are valid for the Supervisor installation state.

## Project structure

- `schemas/*.schema.json`: current public JSON scaffold contracts.
- `schemas/candidates/`: passive Phase 2 candidate contracts that are not agent endpoints.
- `schemas/cddl/*.cddl`: candidate canonical internal-object contracts.
- `schemas/fixtures/`: human-readable fixture descriptions and byte-exact known-answer values.
- `examples/jobs/`: current scaffold and passive proposal examples until atomic cutover.
- `packages/protocol/src/`: TypeScript projections and small boundary helpers.
- `internal/execution/`: existing OS-neutral fake lifecycle; no real backend is connected here.
- `scripts/verify-schemas.mjs`: JSON Schema positive and negative conformance checks.
- `scripts/verify-cddl-fixtures.mjs`: deterministic fixture-byte checks, not a production codec.

## Code style

Use closed discriminated objects and role-specific names. Do not use a generic identifier where a
registration, installation, Supervisor, plan digest, source digest, or epoch digest is intended.

```ts
export interface PlanRegistration {
  objectType: "capsule.plan-registration";
  objectVersion: 0;
  registrationId: RegistrationId;
  planDigest: Sha256Digest;
  installationId: InstallationId;
  epochDigest: Sha256Digest;
  supervisorId: SupervisorId;
}
```

This example illustrates naming and separation only. The accepted schema/fixture owns the final
field set and wire representation.

## Testing strategy

Follow red-green-refactor for behavior changes:

1. add a focused test or invalid fixture and confirm it fails against the old scaffold;
2. add the smallest contract/schema/type implementation that passes it;
3. run the focused package or schema verification;
4. refactor only after the focused tests pass;
5. run the repository verification matrix before handoff.

Tests assert observable acceptance or rejection, not implementation call sequences. Invalid
fixtures identify the rejecting boundary: raw JSON decoder, JSON Schema, semantic planner,
canonical internal decoder, or Supervisor registration.

## Commands

```sh
fnm exec --using=22.22.1 -- pnpm install --frozen-lockfile
fnm exec --using=22.22.1 -- pnpm verify:schemas
fnm exec --using=22.22.1 -- pnpm check
fnm exec --using=22.22.1 -- pnpm lint
fnm exec --using=22.22.1 -- pnpm test
go test ./...
go vet ./...
go build ./...
git diff --check
```

## Boundaries

Always:

- preserve deny-by-default authority and execute-by-registration-ID;
- keep schemas, fixtures, TypeScript, Go, examples, verifier behavior, and documentation aligned;
- reject unknown fields, versions, powers, unsafe integers, and cross-object substitutions;
- retain exact received canonical plan bytes separately from decoded values;
- label internal CDDL/fixture work candidate while ADR-0019 remains Proposed.

Ask first:

- adding a production serialization or cryptographic dependency;
- changing the public v0 product wedge or adding an authority not present in accepted ADRs;
- selecting a backend-specific limit or transport shape before its evidence closes.

Never:

- make a daemon-controlled value an approval, content, key, backend, or trust authority;
- accept replacement plan bytes at attempt/execute time;
- import experiment code into product packages;
- claim contract freeze, backend validation, or security posture beyond exact retained evidence.

## Success criteria

- A valid narrow passive `JobProposal` candidate passes its schema and TypeScript tests without
  activating a daemon, SDK, or MCP endpoint.
- The old mixed `Job` remains unchanged and explicitly pre-freeze until a later atomic cutover; it
  is not extended, adapted, or accepted as a `JobProposal`.
- Unsupported authority and agent-selected paths are structurally unrepresentable.
- Unknown fields/versions, unsafe numbers, invalid paths/slots, and malformed object roles fail the
  documented boundary.
- `ExecutionPlan` and `PlanRegistration` cannot be confused with each other or with a proposal.
- Registration binds exact plan digest, installation, epoch, Supervisor, sequence, and expiry if
  those fields are confirmed by the merged contract review.
- Candidate internal fixture bytes are deterministic and verified without importing experiment
  implementations.
- Existing fake-lifecycle tests still prove that guest-creating backends are rejected.
- The complete repository verification matrix passes.

Before public activation, the atomic cutover must also remove the MCP catalog's direct
`prepare_job`/`run_job` shape and the deprecated Go `Execute(PreparedJob)`/`RuntimeAdapter` scaffold.
Those dormant surfaces are not used to implement or exercise the passive candidate.

## Open questions

These are decided before a corresponding field is marked frozen:

- later backend-independent requested-limit additions after removing unsupported CPU-time,
  arbitrary RAM, PID, scratch, log, and artifact semantics from the old scaffold;
- source-bundle aggregate raw-byte/depth limits beyond the candidate's per-field structural bounds;
- the complete `ExecutionPlan` reference set and whether any counters require decimal strings;
- registration sequence and expiry representation and their relationship to trust epoch changes;
- which shared objects require Swift participation in this slice versus the later production
  canonical-wrapper slice.

Unresolved questions narrow the implemented candidate. They never justify carrying forward fields
from the deprecated `Job` authority model.
