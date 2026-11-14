# C5b19 stop-the-line passive adapter race correction

Status: PASSED for the bounded correction and targeted verification. Parent Gate 3 remains IN_PROGRESS — TRENDING_GOOD. Child execution remains BLOCKED. Full canonical checks and independent review remain Brain-owned.

## Assignment and scope

Internal CI-RACE assignment in `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-decision-reconciliation`, base `8e4fd22cc00cd5e961544f6a803974600dc0e751`. Owned only `internal/execution/authorityplane/passive_product_adapter.go` and `passive_product_adapter_test.go`; no branch, commit, push or PR created. Other agents' documentation/planning changes and the preexisting `.serena/project.yml` edit were preserved.

Defensively validate result/deadline/cancellation arbitration with local synthetic passive harnesses. No child, guest, fixture executable, external identity or remote system was exercised. Deadline constants, admission limits, caller cancellation policy and branch protections remain unchanged.

## Finding and correction

Post-merge CI run 43543205245 failed `TestPassiveCancellationDeadlineStallAndResponseLossSemantics` at the response-loss retry. The hosted failed-job log retained by Brain is `/tmp/c5b19-g3a8-postmerge-ci-failure.log`, SHA-256 `143b82f60fe539b692e89e5a1b91ae7564ef7ee7d4be210f2f0e1dffc1b4177e`.

`runPassiveAdmitted` publishes a buffered result, releases the admission slot, then cancels its method context during worker cleanup. Its select could choose that already-ready cleanup notification and misclassify it as an expired method deadline. The cleanup owns the only explicit cancellation function; its context is rooted in Background, independent of caller cancellation.

The seven-line correction distinguishes `context.Canceled` from actual `context.DeadlineExceeded` in the method-context branch. Cleanup cancellation consumes the already-published buffered result, preserving its exact value and error. Real deadline expiry retains the existing response-unknown disposition. Worker ownership of cleanup/release remains intact, so caller cancellation cannot cancel dispatched work or release its slot early. Simultaneous actual deadline/result and caller/result arbitration is unchanged.

## Tests-first evidence

Added `TestPassiveCompletedResultSurvivesMethodContextCleanup` before the fix. A test-only context delays select evaluation until worker cleanup, making result and cleanup notification ready together; 128 bounded trials each cover nil error and a sentinel method error, without timing sleeps. The original source failed on attempt 0: `(0, method-deadline-after-dispatch-response-unknown)` instead of `(73, nil)`.

This test synchronizes readiness, not the randomized select choice. It is not a deterministic branch-selection proof. Under independent uniform selection of the two ready cases, the chance that 128 trials all avoid the faulty branch is 2^-128; Go's selection specification does not establish independent choices across runs. The observed red failure supplies direct regression evidence. A later test-only `sync.Once` refinement makes repeated/concurrent Done calls safe; final green/race checks include it.

Added `TestPassiveCallerCancellationRetainsMethodContextAndSlot`: after caller cancellation returns response-unknown, the method context remains live and its connection slot remains occupied; a controlled release lets the worker finish and cleanup is awaited. The existing real-deadline/stall/response-loss test was retained unchanged.

Commands used exact Go 1.26.11 at `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64/bin/go`, with `GOTOOLCHAIN=local`:

- RED: `go test ./internal/execution/authorityplane -run 'TestPassive(CompletedResultSurvivesMethodContextCleanup|CallerCancellationRetainsMethodContextAndSlot)$' -count=1` — FAIL, expected behavioral assertion. The shell then printed the log, so its final command status was 0; the test output explicitly records FAIL.
- GREEN final: `go test ./internal/execution/authorityplane -run 'TestPassive(CompletedResultSurvivesMethodContextCleanup|CallerCancellationRetainsMethodContextAndSlot|CancellationDeadlineStallAndResponseLossSemantics)$' -count=100` — PASSED, 2.211s.
- Final package race check: `go test -race ./internal/execution/authorityplane -count=1` — PASSED, 1.984s.
- `gofmt` on both owned files; `git diff --check` — PASSED.

## Evidence hashes

| Artifact | SHA-256 |
| --- | --- |
| `/tmp/c5b19-ci-race-red.log` | `26159331e6f7b351b8a3c6dd92b1c17dfad15116b3f107446a2ef039760a0149` |
| `/tmp/c5b19-ci-race-green.log` | `3e33b02b4ee67c20c361c1bbebb7eb454c3add2a7d7fc5c0a09d045747ed21a3` |
| `/tmp/c5b19-ci-race-package.log` | `6703715216c0870a78957d70e970b92a31016f13d38181827f589bf11ffd33f5` |
| `internal/execution/authorityplane/passive_product_adapter.go` | `785ac54438cea50838b39d8de8a4cf1b8f9b0dd6b887bfedd51beea5543cc9f4` |
| `internal/execution/authorityplane/passive_product_adapter_test.go` | `4621e65f44dd4fcc1ded6f553d8ba29e2c38245fdd6550a51ef89a09274f9484` |

## Method and limitations

Used codebase-memory Tier 2: index ready at 2026-10-04T22:18:10Z, not watched, primary checkout root. Search found runPassiveAdmitted, complete depth-one trace reported three callers (submission, registration, registered-plan retrieval); snippet and coverage checks covered both owned files with metadata_match. Direct exact worktree reads established current source. The graph's inferred `controlledSubmissionCore.release` callee is an imprecise closure-name association; the source shows the local admission-release closure. No broader graph completeness claim.

Confidence high for this narrow race. Local checks run on macOS, not the hosted Linux runner; hosted CI remains required. This passive harness correction does not admit runtime/child execution or establish the full Gate 3 oracle. Brain owns retained evidence copying, full checks, fresh review and publication.
