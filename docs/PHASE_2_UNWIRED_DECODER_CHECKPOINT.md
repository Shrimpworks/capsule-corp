# Phase 2 unwired decoder checkpoint

Date: 2026-08-02

Status: integrated implementation checkpoint. The code remains backend-independent and unwired.
It does not accept ADR-0019 or ADR-0023, activate an endpoint or store, authorize a plan, handle
user content, or create a guest.

## Integrated tasks

| Slice | Task provenance | Integration commit | Executable coverage |
| --- | --- | --- | --- |
| Task 3A strict `JobProposal` byte decoder | `019fc327-80ed-7fe2-8fc0-7a9aea12afb6` | `49ae76b` | 62 TypeScript manifest cases: 38 raw and 24 closed-schema cases across 60 fixtures |
| Task 4A strict internal CBOR wrappers | `019fc327-c094-7592-9cb0-b0a270af30ad` | `b55abbd` | 81 Go manifest cases covering media type, shared scalars, both predecoders/wire objects, and role/cross-object binding |
| PlanRegistration depth-fixture correction | coordinator `019fc2de-552d-77a0-aa47-35ac39d02edc` | `b7940ad` | Replaced array-based depth fixtures with contract-valid nested-map fixtures without changing any limit or decision |

The manifest remains 67 rules, 206 cases, and 278 fixtures. It now records 62 verified TypeScript
targets and 81 verified Go targets. Swift targets, proposal semantic planning, and registration
state remain pending.

## Task 3A outcome

`decodeJobProposal(bytes: Uint8Array)` now performs the strict public raw and closed-schema
boundary without accepting a pre-parsed object. It enforces the retained UTF-8, BOM, single-
document, duplicate decoded-key, integer grammar/range, depth, node, member, element, decoded-text,
closed-field, object type/version, path grammar, slot, label, and scalar-shape fixtures.

Success returns only a deeply frozen passive `DecodedJobProposalCandidate`. Failure returns fixed
internal owner/classification/code data without caller-controlled prose. Caller-owned input is
copied before parsing. No parsing dependency, media-type endpoint, policy resolution, plan
construction, SDK/MCP consumer, content path, runtime, or backend behavior was added.

Task 3B must accept only the successful candidate plus separately supplied trusted resolution
context. It still owns source file and aggregate byte semantics, entrypoint membership,
SourceManifest ordering/identity, canonical inline-input bytes/digest, profile activation and
binding, wall-time default/ceiling policy, and minimum `ExecutionPlan` construction.

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
go test -race ./internal/protocol/v0candidate
go test ./...
go vet ./...
go build ./...
git diff --check
```

The restricted Go run emitted only the known non-fatal user-module-cache write warning.

## Next dependency boundary

Task 3B semantic planning can now proceed against the frozen `DecodedJobProposalCandidate` input.
In parallel, the exact Supervisor registration-state implementation can be prepared against the
Go wrapper, but activation remains blocked on the separately reviewed archival/compaction design.
The registered-plan/fault-injectable fake-backend vertical slice follows both. Public daemon, SDK,
and MCP cutover remains later and atomic.
