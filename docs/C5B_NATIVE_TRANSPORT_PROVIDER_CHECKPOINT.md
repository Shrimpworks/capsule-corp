# C5b12 native transport provider checkpoint

Date: 2026-09-05

Status: local construction/transport checks and independent review `PASSED`.
Parent complete-provider composition and controlled
C5b execution: `BLOCKED`. Owner-only internal alpha: `IN_PROGRESS — TRENDING_GOOD`.
Runtime/profile, installed security boundary and product admission: `BLOCKED`.

## Exact scope and evidence

The [C5b11 checkpoint](C5B_FIXED_RUNNER_SUCCESSOR_CHECKPOINT.md) retained fixed runner
and driver source but no implementations of its 24 Supervisor provider declarations.
The first native implementation now adds seven transport providers in
[the C5b12 experiment](https://github.com/Shrimpworks/capsule-experiments/tree/0e83da32657a0190a7efd17ce203acf3503c905b/experiments/typed-guest-transport-c5b12-native-transport).
Archive source/evidence commit: `0e83da32657a0190a7efd17ce203acf3503c905b`; input archive remains immutable
`f206e4ef2cd326ee74e5b7b2739c62efe6da7d6d`. No existing C5b11 bytes changed.

| Provider effects | Actual responsibility | Evidence boundary |
| --- | --- | --- |
| 1 | Create six distinct pipe pairs, mark CLOEXEC, configure only Supervisor ends nonblocking, start output drain | Local Darwin descriptors; no runner spawn |
| 3 | Accept exact readiness byte within fixed deadline | Extra readiness is rejected before completion acceptance |
| 4–6 | Write exact retained source/input frames and close both writers | Fragmented/EINTR/zero writes, closed reader, backpressure, order and binding refusals |
| 7 | Send the one start byte and close its writer after inputs | No libkrun, guest or lifecycle fact |
| 8 | Continuously drain completion/stderr, retain cap+1, validate exact known answer and final trailer | Fixed 259-byte frame only; no general JSON parser or durable success |

The unlinked provider object reproduces at SHA-256
`a9e2c80ce273f19869a5ba496248b66ac9244fc504cee9c341973644673376de`,
with seven exports and 15 closed system/compiler imports. No process creation,
kill, filesystem-path, database, loader or libkrun import exists. Two builds,
39 native cases, 39 ASan/UBSan cases and six compiled source-mutation refusals
passed. Mutations must fail assertions; build failure or timeout is not a pass.
Retained verification records macOS 26.6.2 (25G83), SDK 26.5, Apple clang 21.0.0
and Node.js 22.22.1 on the owned arm64 host.
Independent review returned `Ready with non-blocking follow-up`. The requested
source/input/test/script hash inventory was added, and the full suite reran;
provider/test source and the reviewed object remained unchanged. Ordinary
verification now rejects stale material/object evidence. The archive retains
the review and author response.

Defensive tests use only exact repository frames, anonymous local pipes, bounded
threads and disposable test processes. They link the new transport code alone.
The peer is a test thread, not a runner or guest. No complete driver/runner,
dylib, libkrun/HVF, VM, guest, signing, Keychain, installed service, protected state
or unrelated system/data was touched. No product package imports the archive.

## Claim limits and next owners

This establishes native pipe behavior for a fixed fixture. One serialized caller
and one attempt per process remain assumptions; no setup/reset/replacement
descriptor or callback API exists. The future Supervisor lifecycle owner must
transfer private peers into the exact runner and close parent copies. Failure
cleanup after endpoint creation in the local tests is harness-owned. Only partial
endpoint-creation cleanup belongs to the implemented providers.

The construction deadline is 1,000 ms from endpoint creation, including setup.
Complete composition must reconcile it with approved guest runtime timing;
these tests do not prove an enforceable guest deadline. Valid bytes, EOF and
readiness never substitute for terminal state, authoritative absence, root
removal or durable completion. The remaining 17 symbols stay unresolved.

| Next milestone | Missing providers | What it enables |
| --- | --- | --- |
| Fixed runner lifecycle and root custody | 2, 9–11, 16–20 | Exact Supervisor-owned identity/FD transfer, terminal/absence checks and teardown/root reconciliation |
| Durable attempt and completion store | 12–15, 21–24 | Durable-before-effect cursors, fencing, unresolved cleanup, completion-before-delivery and restart replay |
| Complete immutable composition review | All owners plus signed identity, provenance, timing and retained fault corpus | A concrete candidate for separate controlled guest authorization |

The C5b11 recovery oracle remains the prerequisite: ambiguous spawn reconciles,
teardown persists safe resume cursor 17 before effect 16, and lost completion
responses reopen/replay stored bytes without re-execution or recommit. No recovery
claim changes here. Preferred-form kernel/libkrunfw source compliance, raw v19/v27
evidence recovery/replacement, cross-host reproduction, installed composition and
runtime/profile/product admission remain `BLOCKED` in their existing scopes.

Reuse decision: P0-3 fixed transport framing `BUILD-NARROWLY`, OS descriptor APIs
`ADOPT-PLATFORM`, native instrumented tests `TEST-ONLY`; no package dependency or
new Supervisor responsibility. The archive plan retains the dependency checklist.
No ADR lifecycle or control-evidence state changes.
