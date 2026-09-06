# C5b14B native durable-owner integration checkpoint

Date: 2026-09-06

Status: native construction, full local verification and independent review
`PASSED` / **Ready** (review instance 1 of 3).
Installed lifecycle boundary, controlled guest execution and product admission: `BLOCKED`.
Parent owner-only internal alpha: `IN_PROGRESS — TRENDING_GOOD`.

## Exact scope and retained evidence

[C5b14A](C5B_DURABLE_PROVIDER_STATE_CORE_CHECKPOINT.md) tested a storage-only Go core.
C5b14B connects its derived persistence mechanics to the actual private observations
from the [C5b13 benign native lifecycle](C5B_NATIVE_LIFECYCLE_PROVIDER_CHECKPOINT.md).
All 24 native fixture provider bodies, both durable gates, recovery checkpoints
17–20 and a fixture-specific registration-only driver now link together.

Evidence: [C5b14B experiment](https://github.com/Shrimpworks/capsule-experiments/tree/294559e78b410f4b25e0e25fcf373ec0c7d1cfb5/experiments/typed-guest-transport-c5b14b-native-storage).
Archive source/evidence commit: `294559e78b410f4b25e0e25fcf373ec0c7d1cfb5`; input archive merge
`a688aabee4989b5000340bf995188b07219fe7c5`. Source origins pin C5b13
`6dc12dca10c9cb88370bf3f16d0e3d6305b7fd7f` and C5b14A
`c4d25e2beec0bc2e886a1df29135dea534b53626`. Original experiment bytes remain unchanged.

The [independent review](https://github.com/Shrimpworks/capsule-experiments/blob/294559e78b410f4b25e0e25fcf373ec0c7d1cfb5/experiments/typed-guest-transport-c5b14b-native-storage/review/INDEPENDENT_REVIEW.md)
returned **Ready** with no actionable findings after rerunning ordinary verification.
All 34 frozen materials and the evidence hash remained unchanged. The
[review closure](https://github.com/Shrimpworks/capsule-experiments/blob/294559e78b410f4b25e0e25fcf373ec0c7d1cfb5/experiments/typed-guest-transport-c5b14b-native-storage/review/CLOSURE.md)
records the later verdict; frozen README/PLAN pre-review wording remains historical.

The native C entry point links Go 1.25.13 as a `c-archive`. An early Go-main probe
failed the native default-SIGCHLD reaper precondition; that unchanged-entry-point
candidate is `NO_GO`. The selected native-fronted library build preserves the
precondition without resetting or bypassing a signal handler. This follows existing
ADR-0029 ownership, rather than creating another lifecycle process or helper.

| Mechanic | Observed fixture behavior | Claim boundary |
| --- | --- | --- |
| Native/Go ownership | Native code retains child/reaper, descriptors, root and transport observations; Go owns the durable snapshot; scalar fields and bytes are copied synchronously | One serialized native caller, fixed registered fixture and trusted local test controls; no daemon endpoint or retained foreign pointer |
| Durable-before-effect gates | Actual C spawn and teardown paths require Go publication; teardown persists resume 17 before signaling; 17–20 checkpoint before observation | No restart PID adoption or installed process-custody reconstruction |
| Completion-last | Go obtains exact bytes through a private native observer requiring terminal exit zero, absence, removed root and valid completion; persisted bytes precede delivery/replay | Benign-child observations only; copied test output is not authenticated client delivery |
| Recovery | Lost completion responses replay stored bytes; missing completion remains unresolved; pre-process endpoint failure stays permanently refused; fresh processes with consumed intent perform no new endpoint/spawn/signal | Process-exit and injected publication faults; no power-loss or rollback claim |
| Store loss | Missing, corrupt, orphaned or mismatched state refuses; copying previously delivered bytes rechecks the store | Trusted exclusive fixture root; no hostile same-UID pathname-race protection |

Native request fields are explicitly translated to the Go core's closed frame
bounds. Recovery uncertainty normalizes conservatively to indeterminate. A new
failure-1/cursor-21 refusal covers endpoint failure before process intent. Storage
cursor replies preserve the native lifecycle phase. Fixture profile/plan/frame and
snapshot schema identities distinguish native observations from C5b14A's test
booleans and existing FakeBackend records.

## Verification and limits

The retained verifier builds four benign child variants with distinct identities:
ordinary success, pre-ready timeout, completion followed by nonzero exit, and
corrupt completion. It executes an uninterposed C entry point as well as local fault
hooks, native sanitizer variants, Go store tests/race and a separate Go race archive.
Thirteen mutating storage routes each receive three publication-edge faults. Four
abrupt-process/reopen pairs cover spawn intent, completion, teardown intent and
orphan staging; missing/corrupt/cross-profile stores also refuse.

The corpus passes 126 native case records, 17 Go suites with 58 subtests, eight
Go-race integration cases and ten compiled mutations requiring the intended
assertion. Twenty ordinary mode-0 artifacts reproduce across two directories.
Material hashes, tool versions, native import/export inventory, reproduction and
mutation outcomes are retained in
`evidence/results.json`. Native ASan/UBSan instrument the parent C code, not the
ordinary Go archive or benign child. The race archive is a separate variant; its
linker diagnostic is retained. Debug sanitizer binaries are hashed but excluded
from the ordinary two-directory reproduction claim. Darwin archive metadata and
Mach-O object-symbol paths are normalized using physical output paths; source
identity, reaper checks and validation guards remain bound.

One pre-freeze run refused during case 20 nominal setup; original output did not
identify its step or timing, so cause remains unclassified. Added diagnostics, the
next full corpus and 30 focused repetitions passed. The verifier has no retries or
extended deadline; timing and load tolerance remain unproven. This failed run is
retained as a limitation rather than recast as a successful check.

Canonical `pnpm install`, `check`, `lint`, `test`, `verify:schemas`, `verify:adrs`,
Go tests/vet/build and pinned `govulncheck@v1.6.0` pass. Full `golangci-lint` still
reports 50 pre-existing revive documentation findings; non-revive checks and the
new-code revive ratchet pass. This documentation-only checkpoint adds no Go code.

Only fixed benign local processes, anonymous pipes, a deterministic 64-KiB root and
owned temporary directories participate. Child containment alarms after four
seconds; the harness alarms after seven. The harness may reap its captured child
after assertions, providing no lifecycle evidence or signal authority. No guest,
interpreter, libkrun/HVF, signing, Keychain, installed service, user content or
unrelated asset participates. The normal linked driver requires only libSystem.

The one-attempt 4-KiB snapshot and 32-generation cap remain fixture bounds. No
physical disk/power-loss, backup/restore, rollback defense, installed identity,
multi-user service, selected product engine, ADR-0040 performance/admission or
external delivery claim is established. No product code imports the archive.
Existing FakeBackend stores, schemas, ADR lifecycles and control-evidence levels
remain unchanged. Reuse follows the fixed-snapshot semantic oracle, standard Go
SHA-256/testing/cgo, existing Darwin primitives and narrowly built archive policy.

## Next goalposts

| Goalpost | Status | What it enables |
| --- | --- | --- |
| Durable native benign fixture | Construction, full local verification and independent review `PASSED` / **Ready** | One complete native transport/lifecycle/store flow with recovery and replay evidence |
| Immutable runner composition | `BLOCKED` on installed/launch identity, restart custody, timing/provenance and independent candidate review | A concrete candidate for separately authorized controlled guest execution |
| Owner alpha, then external beta | Owner alpha `IN_PROGRESS — TRENDING_GOOD`; external beta `BLOCKED` | Approved end-to-end work, followed by installation/distribution, continuity and release evidence |

The next bounded task should reconcile the completed provider composition with the
runner candidate's remaining identity, custody, timing and provenance requirements
before selecting its next executable slice. Existing preferred-form kernel/libkrunfw
source compliance, missing raw v19/v27 evidence and cross-host reproduction remain
in their existing workstreams. Completing this benign milestone grants no authority
to execute a guest.
