# Phase 2 unwired decoder checkpoint

Date: 2026-08-02

Status: integrated implementation checkpoint. The code remains backend-independent and unwired.
It does not accept ADR-0019 or ADR-0023, activate an endpoint or consumer, authorize execution,
handle user content, or create a guest. The fixed Task 4B file store is limited to test/development
conformance and cannot be activated before archival/compaction review.

## Integrated tasks

| Slice | Task provenance | Integration commit | Executable coverage |
| --- | --- | --- | --- |
| Task 3A strict `JobProposal` byte decoder | `019fc327-80ed-7fe2-8fc0-7a9aea12afb6` | `49ae76b` | 62 TypeScript manifest cases: 38 raw and 24 closed-schema cases across 60 fixtures |
| Task 3B semantic proposal resolution | `019fc343-699c-7e70-9175-37dc44fca6f3` | `1e1be4e` | 18 TypeScript semantic-resolution cases: 11 accept and 7 reject across 14 fixtures |
| Task 4A strict internal CBOR wrappers | `019fc327-c094-7592-9cb0-b0a270af30ad` | `b55abbd` | 81 Go manifest cases covering media type, shared scalars, both predecoders/wire objects, and role/cross-object binding |
| PlanRegistration depth-fixture correction | coordinator `019fc2de-552d-77a0-aa47-35ac39d02edc` | `b7940ad` | Replaced array-based depth fixtures with contract-valid nested-map fixtures without changing any limit or decision |
| Task 4B exact registration state | `019fc343-699c-7e70-9175-37bdce364c2a` | `43cc268` | 40 Go registration-state cases: 18 accept and 22 reject, with exact post-state fixtures |

The manifest remains 67 rules, 206 cases, and 278 fixtures. It now records 80 verified TypeScript
targets: 62 Task 3A raw/schema targets plus all 18 Task 3B semantic-resolution targets. It records
121 verified Go targets: 81 Task 4A wrapper targets plus all 40 Task 4B registration-state targets.
Swift targets and exact `ExecutionPlan` construction remain pending.

## Task 3A outcome

`decodeJobProposal(bytes: Uint8Array)` now performs the strict public raw and closed-schema
boundary without accepting a pre-parsed object. It enforces the retained UTF-8, BOM, single-
document, duplicate decoded-key, integer grammar/range, depth, node, member, element, decoded-text,
closed-field, object type/version, path grammar, slot, label, and scalar-shape fixtures.

Success returns only a deeply frozen passive `DecodedJobProposalCandidate`. Failure returns fixed
internal owner/classification/code data without caller-controlled prose. Caller-owned input is
copied before parsing. No parsing dependency, media-type endpoint, policy resolution, plan
construction, SDK/MCP consumer, content path, runtime, or backend behavior was added.

## Task 3B outcome

`resolveJobProposal` accepts only a decoder-issued candidate plus separately constructed immutable
profile-registry and user-policy contexts. It enforces per-file and aggregate source-byte caps,
entrypoint membership, deterministic SourceManifest bytes/digest, bounded canonical inline-JSON
bytes/digest, exact active-profile resolution, and requested/defaulted wall time without clamping.

Success returns only deeply frozen passive source, input, profile, wall-time, output, label, and
plan-input values. It adds no installation, trust, runtime-bundle, backend, policy-decision, expiry,
approval, or launch authority. Task 3C must combine these passive inputs with separately supplied
trusted bindings to construct and encode the exact minimum `ExecutionPlan`.

## Task 4A outcome

The Go `v0candidate` package now provides bounded deterministic-CBOR scanning and object-specific
strict wrappers for `ExecutionPlan` and `PlanRegistration`. The wrapper retains a private immutable
copy of the exact received authoritative bytes, computes the plan digest over only that copy, and
returns defensive copies of bytes and decoded views.

Decoding requires complete trusted role bindings. Same-width IDs and digests do not establish their
own nominal role, and there is deliberately no unbound decode path that implies authority. The
implementation rejects malformed, nonpreferred, tagged, floating-point, duplicate/noncanonical
map, trailing, wrong-object, wrong-role, out-of-budget, out-of-range, and mutated-buffer cases. It
adds no CBOR dependency and does not call `SupervisorCore` or persist a registration.

The later registration slice must call `DecodeExecutionPlan(received, trustedRoleBindings)` and
atomically persist the wrapper's `AuthoritativeBytes()` and `Digest()` with the separately issued
wire registration and state oracle.

## Task 4B outcome

`internal/execution/registrationstate` now calls the Task 4A wrapper before trusted-time or
registration state can change. It persists a monotonic trusted-time high water separately, then
performs a final transactional reread for installation/epoch binding, trust fencing, strict expiry,
fixed capacity, installation-global sequence, fresh identifier, deterministic wire encoding, and
the complete retained record. The sequence-exhaustion transition commits `repair-required` without
creating authority; confirmed-abort injection exposes neither a sequence nor a partial record.

The fixed store retains every record without eviction, uses whole-state fsync/rename transactions,
and returns defensive copies. It remains a bounded single-process test/development store with no
archive, compaction, encryption, secure deletion, or production multi-process coordination. The
`RegistrationResolver` handoff for the later fake backend accepts only a Supervisor-issued
registration ID and returns the stored exact plan after time/trust/binding checks. It accepts no
replacement plan bytes or backend options. The old `SupervisorCore.RegisterPlan` path remains
unchanged and unwired.

## Fixture correction

The initial Task 4A pass identified that the retained PlanRegistration depth-exact fixture used
one-element arrays even though ADR-0023 sets that object's array maximum to zero. The product
predecoder correctly rejected it. The generator now expresses the same exact depth and cap-plus-one
cases with deterministic one-entry nested maps. This preserves the depth maximum, zero-array rule,
case IDs, decisions, and corpus size while making the exact-maximum fixture independently valid.

## Verification

The combined post-merge branch passed:

```sh
fnm exec --using=22.22.1 -- pnpm install --frozen-lockfile
fnm exec --using=22.22.1 -- pnpm check
fnm exec --using=22.22.1 -- pnpm lint
fnm exec --using=22.22.1 -- pnpm test
fnm exec --using=22.22.1 -- pnpm verify:schemas
fnm exec --using=22.22.1 -- pnpm format:check
fnm exec --using=22.22.1 -- pnpm site:build
go test -race ./internal/protocol/v0candidate ./internal/execution/registrationstate
go test ./...
go vet ./...
go build ./...
git diff --check
```

The restricted Go run used a task-specific cache under `/tmp`; all listed commands completed.

## Next dependency boundary

Task 3C exact plan construction can now proceed against the frozen Task 3B resolved inputs. In
parallel, an unwired registered-plan/fault-injectable fake-backend slice can consume only the Task
4B `RegistrationResolver` ID-only interface. Consumer activation remains blocked on the separately
reviewed Supervisor archival/compaction design. Public daemon, SDK, and MCP cutover remains later
and atomic.
