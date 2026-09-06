# C5b14A durable provider-state core checkpoint

Date: 2026-09-05

Status: storage-only Go construction and local verification `PASSED`.
Independent review: `PASSED` / `Ready` (fresh review instance 1 of 3).
Native durable owner, complete composition and product admission: `BLOCKED`.
Parent owner-only internal alpha: `IN_PROGRESS — TRENDING_GOOD`.

## Exact scope and retained evidence

The [C5b13 checkpoint](C5B_NATIVE_LIFECYCLE_PROVIDER_CHECKPOINT.md) retains sixteen
native transport/lifecycle bodies. C5b14A adds a separate one-attempt Go storage
conformance core for intent, fencing, safe recovery cursors and immutable completion
replay. Its IDs, fixed JSON result and observation scope explicitly identify a
storage-only fixture. They are neither C5b13 benign-child nor C5b11 VMM evidence.

Evidence: [C5b14A experiment](https://github.com/Shrimpworks/capsule-experiments/tree/c4d25e2beec0bc2e886a1df29135dea534b53626/experiments/typed-guest-transport-c5b14a-storage-core).
Archive source/evidence commit: `c4d25e2beec0bc2e886a1df29135dea534b53626`; input archive merge:
`44dca5cf73159f6e32764984fc23b7babcc05bfe`. Previous experiment bytes remain unchanged.

This first storage increment does not implement the eight native C providers
(12–15, 21–24), replace either C5b13 gate mock, or link the 24-provider driver. It
provides Go logical operations for their later integration. Native frame-field
conventions, endpoint failure before process intent, recovery-outcome translation,
observation custody and C/Go ownership remain explicit C5b14B obligations.

| Mechanic | Observed storage-only behavior | Claim boundary |
| --- | --- | --- |
| Owner and snapshot | Empty-directory bootstrap; exclusive lock; strict bounded/canonical existing-state reopen; owner and live-state identity checks | Trusted owned fixture root; no installed root, rollback defense or concurrent same-UID pathname protection |
| Intent before effect | Consumed spawn intent is published before return; teardown persists step 16 / resume 17 once | No native effect executes; C5b14B must bind successful durable return to the sole effect owner |
| Fencing and recovery | Post-publication uncertainty fences the handle; lost intent resumes unresolved at 17; cleanup uncertainty stays durable | No PID adoption, restart process custody, automatic orphan deletion or repaired fresh state |
| Completion and replay | Bound trusted fixture observations gate fixed result publication; missing result remains unresolved; existing result reopens/replays exact copied bytes at 22/23 | Test testimony is not terminal/absence/root evidence; a Go return is not authenticated client delivery |

## Verification and limitations

Go 1.25.13 on Darwin/arm64, Node 22.22.1. Retained verification includes 16 top-level
tests and 58 named subtests, the ordinary and race suites, vet/build, two-directory
trimmed test-binary reproduction and nine compiled restoration mutations. Each
mutation must reach its intended test assertion failure; compile errors and timeouts
do not count. Counts include suite parents and are not separate security controls.
The verifier checks source, tests, plan, README and verifier material hashes before
ordinary verification. Exact test names, versions, binary hash and mutation details
are retained in the archive's `evidence/results.json`.

Independent review reproduced the ordinary verifier, confirmed all thirteen frozen
materials unchanged, and returned `Ready` with no actionable findings. The archive
retains the separate bootstrap, author explanation and review report. That review
covers the storage experiment; this canonical documentation is its reconciliation.

Thirteen mutating routes each receive before-publication, after-publication and
after-directory-sync fault injection. Four subprocess cases exit with code 73 at
fixed staging/spawn/teardown/completion edges. Reopen observes the expected complete
snapshot; orphan staging remains present and refuses. Tests also cover consumed
intent, teardown redrive refusal, missing completion, immutable copied replay,
concurrent readers, strict decoding and tested identity substitutions.

These are local process-exit and injected-edge tests, not physical disk failure,
APFS power-loss, hostile filesystem concurrency, installed custody or durability
continuity evidence. The one-attempt 4-KiB snapshot and generation cap of 32 are
fixture bounds, not ADR-0040 attempt/segment/latency admission. No performance or
multi-user result is claimed. No native runner, guest, backend, signing, Keychain,
installed service, user content or unrelated asset participates.

Reuse follows the fixed-store semantic-oracle roadmap, standard Go SHA-256/testing,
platform owner lock and narrowly built archive experiment. Existing FakeBackend
stores retain their original record identities and formats. No package dependency,
new privileged role or second selected product database is introduced. ADR-0029's
native front end / Go authority-core direction and ADR-0041's Supervisor-owned
lifecycle duties remain unchanged. ADR lifecycle, control-evidence levels and
runtime/profile/product admission do not change.

Canonical repository checks: `pnpm install`, `pnpm check`, `pnpm lint`, `pnpm test`,
`pnpm verify:schemas`, `pnpm verify:adrs`, `go test ./...`, `go vet ./...`,
`go build ./...` and `govulncheck@v1.6.0` pass. Full `golangci-lint run ./...` retains
50 pre-existing `revive` documentation findings; the blocking non-revive suite and
new-code revive check pass. This documentation checkpoint adds no product code.

## Milestones from here

| Goalpost | Status | What completing it enables |
| --- | --- | --- |
| C5b12 transport and C5b13 lifecycle mechanics | `PASSED` in their local fixture scopes | Sixteen native bodies ready for a durable-owner bridge |
| C5b14A storage core | `PASSED` for local Go conformance | Tested intent, cursor and replay mechanics for native integration |
| C5b14B native durable-owner integration | `BLOCKED` on eight native providers, two gates, observation/ABI binding and complete benign-driver fault evidence | First complete durable native fixture composition |
| Immutable runner and controlled guest attempt | `BLOCKED` on installed/launch identity, restart custody, timing/provenance, review and exact authorization | End-to-end controlled guest evidence |
| Owner alpha, then external beta | Owner alpha `IN_PROGRESS — TRENDING_GOOD`; external beta `BLOCKED` | Approved work followed by installed signing/distribution, continuity and release evidence |

Next implementation is C5b14B: wire the storage owner to C5b13's retained native
facts, persist every required recovery checkpoint, and review the complete benign
driver before any guest claim. Completion of this core advances the storage
milestone; it does not close that parent milestone. Existing preferred-form
kernel/libkrunfw compliance, missing raw v19/v27 evidence and cross-host reproduction
remain in their existing workstreams.
