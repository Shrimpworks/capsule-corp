# Phase 2B and Gate C parallel-task checkpoint

Date: 2026-08-02

Status: integrated local checkpoint. This document records the completed parallel-task group and
the next implementation boundary. It does not accept ADR-0023, activate a public endpoint, admit a
runtime or backend, or authorize user bytes.

Historical-status note: the P0-2 row below records the earlier minimal-removal result at `d3be865`.
PR #30 later selected `GOVERNED-PATCH` after a direct-block-root prototype booted and reran the
bounded corpus without `NullFs`; removal remains credible but unadmitted. Current status lives in
the [Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md) and
[workstream ledger](WORKSTREAM_EVIDENCE_LEDGER.md).

## Task group

| Workstream | Task provenance | Integration | Result |
| --- | --- | --- | --- |
| Phase 2B Task 2.3: proposal/source/input fixtures | `019fc2e6-cf9f-7482-82b5-23e67a399ce6` | `4afbdfa` | Added 18 rules, 42 cases, and 39 fixtures; the corpus reached 55 rules, 147 cases, and 130 fixtures. |
| Gate C P0-2: `NullFs` disposition | replacement `019fc2e8-445e-7cb2-b4c2-54d84282c3fe`; original `019fc2e6-cf9d-7210-b2f3-f3bf2244e83a` ended at the app layer | `d3be865` | The smallest removal built but broke guest bootstrap. The exact block-root profile remains unsupported. |
| Gate C P0-0: stock Bun runtime authority | `019fc2e6-cf9f-7482-82b5-240992c79419` | `004c047` | Stock Bun 1.3.14 failed the prohibited-power contract. `RUNTIME-001` remains unsupported. |
| Phase 2B Task 2.4: plan/registration/state fixtures | `019fc2e6-cf9d-7210-b2f3-f3dfc37458a7` | `fd51ac4` | Added 12 rules, 59 cases, and 148 fixtures; the corpus reached 67 rules, 206 cases, and 278 fixtures. |
| Cross-workstream reconciliation | coordinator `019fc2de-552d-77a0-aa47-35ac39d02edc` | `99192d4` | Reconciled the roadmap, boundary decisions, and durable evidence ledger. |

The app-layer failure in the original P0-2 task is not a product finding. Its replacement retained
the bounded investigation and evidence in the repository.

## Durable contract outcomes

Task 2.3 closed the passive fixture oracles for proposal structure, source paths and ordering,
entrypoint membership, source counts and byte budgets, fixed input/output slots, trusted profile
resolution, labels, canonical inline JSON, and requested/defaulted/over-ceiling wall time.

Task 2.4 closed the passive fixture oracles for exact plan and registration bytes, scalar and
digest domains, authenticated caller roles, fresh registration, installation-global sequence,
expiry equality, monotonic time, trust fencing, fixed capacity, mutation isolation, atomic
failure, and wire/stored-record separation. Its registration-state matrix has 40 cases: 18 accept
and 22 reject. Every stateful case retains exact before, operation, and after fixtures.

Retained known answers:

| Object | Exact length | SHA-256 |
| --- | ---: | --- |
| Ordinary source manifest | retained fixture | `e5e09b2435baedf897526a89c698c0b0531437a69472372ae426f62d801fc171` |
| Ordinary canonical inline input | retained fixture | `bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72` |
| Ordinary `ExecutionPlan` | 530 bytes | `627f9524479000dab6f3cee1d70c0428c63285bcadbc2cb3c6e8018b2dea008c` |
| Ordinary `PlanRegistration` | 165 bytes | `f3569d37ad6d787c2cdd575ef9ec6c369bbe495157c43110fc9e9d610a277614` |

These fixtures are executable specifications for future language implementations. They are not a
claim that strict decoders, semantic planners, registration storage, or public consumers are
implemented. ADR-0023 remains Proposed.

## Durable Gate C outcomes

### Stock Bun

The exact Bun 1.3.14 investigation observed reachable subprocess, executable replacement,
FFI/native-loading, inspector, Worker, and inherited-descriptor authority despite all relevant
stock deny flags. Stock flags cannot satisfy the advertised v0 runtime contract.

The governed construction-level branch also failed its explicit reviewability gate: its minimum
honest closure spans at least 40 hand-authored files plus 10 generated outputs, while a narrow
process/exec self-seal cannot independently close Worker or native loading. Capsule must now choose
an alternate runtime and supersede ADR-0003's Bun-first implementation choice. Execution requiring
`RUNTIME-001` refuses throughout.

### `NullFs`

`NullFs` is not introduced merely by building or loading libkrun, but it is unconditional for the
exact retained `krun_set_root_disk_remount` block-root route. Removing only that device prevented
the guest from mounting its root because the current dummy virtiofs route supplies bootstrap
files and mount points.

The current profile remains unsupported. Next work must either provide an alternate bootstrap
with a closed device list or formally review and exercise the complete residual guest-reachable
virtio queue, FUSE decoder, worker, overlay, and `NullFs` surface under an applicable retained
sanitizer/coverage corpus. No libkrun adapter may handle user bytes before that and the other Gate C
P0 campaigns close.

## Application reality at this checkpoint

The repository contains a runnable daemon shell with health, version, and unavailable-profile
endpoints. It also contains an OS-neutral Supervisor state machine exercised through an in-memory
store and a fake lifecycle that refuses guest-creating backends.

It does not yet contain a job-submission endpoint, wired strict planner, durable Supervisor
process/store, native Broker approval UI, real approval signatures, content custody flow, real
guest execution, or end-to-end SDK/MCP job workflow. The conformance corpus and Supervisor state
machine are durable product foundations; `experiments/` implementations remain non-production
evidence and cannot be imported into product packages.

## Combined verification

The integrated tree passed:

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

The conformance runner passed 22 tests and validated 67 rules, 206 cases, and 278 fixtures. The Go
commands completed successfully; the restricted local run emitted only a non-fatal user-module-
cache write warning.

## Next implementation group

Two backend-independent tasks may begin in parallel from this checkpoint:

1. Implement strict public JSON proposal decoding and semantic planning as an unwired library,
   making the applicable Task 2.3 language targets executable.
2. Implement strict deterministic-CBOR internal wrappers and exact received-byte validation in Go,
   making the applicable plan/registration fixture targets executable without activating a store
   or endpoint.

After both integrate, connect them to the existing Supervisor state machine through a
fault-injectable fake backend. That vertical slice remains sequential because it depends on both
decoding/planning and exact internal wrapper behavior. Public daemon, SDK, and MCP cutover remains
blocked until the fake-backend slice and the separately reviewed Supervisor archival/compaction
design are complete.

Backend P0 research may continue independently, but it must not add backend flags, paths, runtime
claims, or experiment code to the backend-independent contract implementation.
