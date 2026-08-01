# Phase 2B boundary-decision specification

Status: proposed specification for human review. No validator or authority-bearing behavior is
implemented by this document.

## Objective

Resolve only the exact protocol and registration questions needed to build the manifest-driven
conformance corpus and then strict unwired validators. The normative proposal is
[ADR-0023](adr/0023-bound-protocol-decoding-and-registration.md).

Success means Task 2 can assign an exact byte fixture, first owning layer, expected classification,
and authority-state oracle without guessing. This slice does not complete the launch-capable plan,
select production dependencies, activate consumers, or handle user bytes.

## Assumptions for review

- Phase 2A object shapes remain passive candidates.
- Parser/registration budgets can be conservative product decisions; real guest transport,
  runtime, and backend limits still require P0 evidence.
- The first public wedge remains dependency-free source, inline JSON input, one inline JSON output,
  fixed logical slots, and optional requested wall time.
- Public JSON is strict but noncanonical. Canonical inline-input bytes are a separate bounded
  representation derived only after raw/schema validation.
- Proposed ADR-0019 remains Proposed throughout this slice.

If any assumption changes, update this specification and ADR before implementing a corpus oracle.

## Decisions under review

ADR-0023 proposes exact rules for:

1. object-specific media types and first owning boundary;
2. raw JSON bytes, depth, nodes, members, elements, decoded text, source, inline-input, and label
   budgets;
3. duplicate keys, UTF-8, BOM, one-document, trailing-data, and integer-token rejection;
4. ASCII source-path semantics and deterministic `SourceManifest` identity;
5. a restricted canonical JSON representation for exact inline-input bytes;
6. plan/registration deterministic-CBOR predecoder budgets;
7. nonzero role-specific IDs and role-resolved digest behavior;
8. fresh registration per call, installation-global monotonic sequence, 300-second lifetime,
   monotonic effective time, fixed capacity, and non-evicting retention; and
9. fixed internal conformance classifications separated from future public error codes.

## Scope

### Included

- proposed ADR and documentation reconciliation;
- exact Task 2 fixture/runner requirements;
- a migration note for Phase 2A character-count schema limits versus semantic UTF-8 byte limits;
- explicit non-goals for backend, approval, attempt, receipt, and consumer activation.

### Excluded

- JSON tokenizer, canonical encoder, CBOR/COSE wrapper, database, daemon, SDK, MCP, Broker, or
  Supervisor behavior;
- changes to candidate wire bytes or schemas before the proposed decisions are reviewed;
- new dependencies or reuse of experiment packages in product code;
- a production retention/archive policy for terminal registration records; the proposed bounded
  first store is intentionally suitable only for unwired conformance and fake-backend work;
- completion-frame, console-port, inline-result, RAM, vCPU, CPU, scratch, artifact, or backend caps;
- public error objects, full receipt/profile migration, or legacy `Job` removal;
- real guest creation or execution.

## Task 2 implementation plan

### Task 2.1: Define the conformance manifest schema

Each case records a stable case ID, object, wire format/media type, exact fixture path/hash, context,
expected decision/classification/owner, authority-state oracle, and applicable implementations.

Acceptance:

- the manifest is closed and rejects an unknown field/version;
- every listed fixture exists and every fixture is listed;
- raw malformed fixtures remain bytes rather than embedded parsed JSON.

Likely files: `schemas/conformance/v0/manifest.schema.json`, `schemas/conformance/v0/manifest.json`,
and one verifier test. Estimated scope: medium.

### Task 2.2: Add shared raw/scalar fixtures

Seed exact/max/cap-plus-one cases for media type, empty/trailing/duplicate/UTF-8 JSON, deterministic
CBOR, depths/counts/text budgets, integer grammar/range, ID widths/zero, and digest widths/domain.

Acceptance:

- every independently reachable ADR-0023 raw/scalar rule has one accept and at least one focused
  reject case; the derived JSON-node cap has an exact-maximum accept and a cap-plus-one reject that
  explicitly records its unavoidable collection-cap overlap;
- escaped-equivalent duplicate keys and noncanonical CBOR are retained byte-exact;
- each rejection owns no authority-bearing state.

Likely files: `schemas/conformance/v0/shared/` and manifest entries. Estimated scope: medium.

### Task 2.3: Add proposal/source/input fixtures

Cover structural authority omission, exact source-path/manifest ordering, entrypoint membership,
source aggregate bytes, canonical inline JSON, fixed slot roles, labels, and requested/defaulted/
over-ceiling wall time.

Acceptance:

- independently written encoders agree on accepted source/input bytes and digests;
- ceiling+1 rejects while exact ceiling and trusted default are unchanged;
- invalid paths/slots/profile resolution fail at the documented first owner.

Likely files: `schemas/conformance/v0/job-proposal/` and fixed resolver contexts. Estimated scope:
medium.

### Task 2.4: Add plan/registration and state fixtures

Cover exact plan bytes/digest, cross-object/domain substitutions, authenticated caller role,
duplicate plan registrations, sequence, expiry equality, epoch fencing, clock rollback, capacity,
mutation after submission, and durable no-state-change failures.

Acceptance:

- Plan A/registration B and correct-width wrong-role values fail closed;
- identical plan bytes create distinct increasing registrations;
- expiry, rollback, overflow, and capacity cannot resurrect or evict authority.

Likely files: `schemas/conformance/v0/{execution-plan,plan-registration,contexts}/`. Estimated scope:
medium.

### Task 2.5: Add one repository conformance runner

The initial runner verifies manifest completeness, raw fixture hashes, structural candidate
expectations, exact known-answer bytes/digests, and explicitly marks semantic/registration cases
pending until their unwired implementations land. It must fail for unlisted or silently skipped
cases.

Acceptance:

- `pnpm verify:schemas` includes the runner without importing experiment code;
- the runner distinguishes implemented assertions from declared pending language targets;
- a required case cannot be skipped by deleting its file or manifest entry independently.

Likely files: `scripts/verify-contract-conformance.mjs`, `package.json`, and focused tests/fixtures.
Estimated scope: medium.

## Checkpoints

After Tasks 2.1 and 2.2:

- manifest and raw fixture corpus are reviewable independently of semantic implementation;
- repository check, lint, tests, schema verification, Go test/vet/build, and diff checks pass.

After Tasks 2.3 through 2.5:

- every accepted decision has an exact fixture oracle;
- no case implies production wrapper/backend validation;
- Task 3 strict unwired proposal decoding can start without inventing a contract value.

Consumer activation additionally remains blocked on the Supervisor-owned registration
archive/compaction design called out by ADR-0023; increasing the test-store ceiling is not a
substitute.

## Commands

```sh
fnm exec --using=22.22.1 -- pnpm install --frozen-lockfile
fnm exec --using=22.22.1 -- pnpm check
fnm exec --using=22.22.1 -- pnpm lint
fnm exec --using=22.22.1 -- pnpm test
fnm exec --using=22.22.1 -- pnpm verify:schemas
fnm exec --using=22.22.1 -- pnpm format:check
fnm exec --using=22.22.1 -- pnpm site:build
go test ./...
go vet ./...
go build ./...
git diff --check
```

## Boundaries

Always:

- retain malformed inputs and accepted canonical bytes exactly;
- separate raw, schema, semantic, policy, registration, and channel-authentication ownership;
- assert no authority-bearing state change on every rejected registration case;
- preserve exact-or-refused limits and execute-by-registration-ID.

Ask first:

- changing any proposed maximum or lifetime after fixture review;
- adding a production parsing/canonicalization dependency;
- changing the public wedge, canonical object fields, or registration authority model.

Never:

- activate a public proposal/attempt endpoint in this slice;
- treat ordinary JSON parsing or decode/re-encode output as authoritative bytes;
- import experiment implementations into product code;
- add backend flags, paths, transport framing, guest execution, or stronger posture claims.

## Success criteria

- ADR-0023 is reviewed as an exact proposed decision rather than hidden implementation policy.
- Every Task 2 case can be assigned an exact expected owner/classification or is explicitly deferred
  because it belongs to backend P0 or a later object contract.
- The task breakdown contains no cross-subsystem task larger than a focused reviewable session.
- Documentation consistently distinguishes parser budgets from unproven backend/transport limits.
- No product behavior, schema candidate, or canonical fixture bytes change before human review.
