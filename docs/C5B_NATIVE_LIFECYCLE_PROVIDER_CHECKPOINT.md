# C5b13 native lifecycle fixture checkpoint

Date: 2026-09-05

Status: same-session benign lifecycle construction and local checks `PASSED`.
Independent review: `PASSED` / `Ready` (fresh review instance 2 of 3).
Complete lifecycle boundary, durable store, guest composition and product admission: `BLOCKED`.
Parent owner-only internal alpha: `IN_PROGRESS — TRENDING_GOOD`.

## Exact scope and retained evidence

The [C5b12 checkpoint](C5B_NATIVE_TRANSPORT_PROVIDER_CHECKPOINT.md) supplies seven native
transport bodies. C5b13 adds bodies for effects 2, 9–11 and 16–20 in a separate benign
fixture composition: fixed child launch, terminal observation, authoritative direct-child
absence, one-shot teardown, and unlinked-root cleanup. The new executable, 64-KiB root,
profile, plan, registration, attempt and frames have their own identities. They do not
reuse the approved-byte identities of C5b11's VMM candidate.

Evidence: [C5b13 experiment](https://github.com/Shrimpworks/capsule-experiments/tree/6dc12dca10c9cb88370bf3f16d0e3d6305b7fd7f/experiments/typed-guest-transport-c5b13-native-lifecycle).
Archive source/evidence commit: `6dc12dca10c9cb88370bf3f16d0e3d6305b7fd7f`; input archive merge:
`ad5217fc17b97b3c3cac44d383e6b296c111d310`. Existing C5b11/C5b12 bytes remain unchanged.

| Mechanic | Local evidence | Claim boundary |
| --- | --- | --- |
| Fixed descriptor handoff | Hash-bound benign executable/root; private six-pipe peers; child FDs 0–7; extra inherited FD refused by fixture checks | Point-in-time executable hash inside an exclusive test directory; no installed launch identity or same-UID pathname protection |
| Terminal and absence | Sole direct-child PID returned by successful spawn; exclusive waitpid receipt; valid completion followed by exit 23 refuses nominal success | No PID adoption after failed spawn, no restart reconstruction, process-tree or guest-resource absence claim |
| Teardown reconciliation | Fixed store gate before one-shot signal; safe resume cursor 17 before effect 16; repeated request cannot signal again; signal-response loss reconciles | Gate acknowledgments are test-only and non-durable; no store implementation |
| Root custody/removal | Exact copied root reopened read-only and unlinked; retained identity checked; absence required before release; replacement/ambiguous close withholds removal | Same-session disposable fixture only; no installed or crash/power-loss evidence |

The unlinked object exports exactly 16 providers. Two fixed internal store gates remain
undefined, and the eight store providers (12–15, 21–24) remain absent. Tests call the
fixture providers directly; the original 24-provider registration driver is not linked.
This is not a completed nine-provider product lifecycle boundary or a runnable guest
composition. No product code imports the archive.

## Verification and limitations

Four separately built benign variants cover ordinary success, pre-ready timeout, valid
completion followed by nonzero exit, and corrupted completion. Each has a distinct
executable/profile/plan/frame identity. Two output directories reproduce all four
executables, provider objects, generated sources, bindings and frames byte-for-byte.
27 native and 27 ASan/UBSan parent cases plus nine compiled assertion mutations pass.
The full material/import/export inventory and toolchain are retained in archive evidence.
Ordinary `node scripts/verify.mjs` refuses stale source/test/input/generated evidence.
Children use the same ordinary artifacts in both suites and are not instrumented.

First independent review found that transport writer-close ambiguity could be forgotten
before root cleanup. The fix retains uncertainty in shared private state. SOURCE/START
close faults, including a deliberately retained-open fixture FD, now withhold cleanup
without retry. A compiled mutation proves that guard is failure-sensitive. Second fresh
review reran ordinary verification and returned `Ready`; the archive retains both results.

Tests use only the owned Mac, compiled benign child, deterministic root, anonymous pipes
and disposable local directories. No actual runner, interpreter, libkrun/HVF, VM, guest,
signing, Keychain, installed service, protected store or unrelated data participates.
The child has its own four-second containment alarm; the parent has a five-second
watchdog. Setup and cleanup each use a fixed 1,000-ms budget. These are fixture bounds;
complete composition still must reconcile setup with the approved guest runtime policy.

One serialized caller and exclusive child reaper are preconditions. Non-default or
auto-reaping SIGCHLD disposition, forked owner, lost spawn response, lost reaper custody,
wrong identities/cursors, replaced root and ambiguous close refuse. After any attempted
spawn, terminal/absence require retained child custody. A definite pre-spawn refusal can
settle a no-child/no-root state; that terminal fact does not claim an observed exit.
No signal is sent after reap, and no failed spawn result can supply replacement PID
authority. The fixture's executable pathname race is explicitly outside this evidence.

Reuse: existing Supervisor responsibilities from ADR-0041; public Darwin descriptor,
spawn/wait and CommonCrypto APIs; narrowly built root-copy custody; test-only native
instrumentation. No package, persistence engine, new privileged helper or daemon route.
The archive plan records the dependency-policy checklist. ADR lifecycle, control-evidence
levels, installed security boundary and runtime/profile/product admission do not change.

## Milestones from here

| Goalpost | Status | What completing it enables |
| --- | --- | --- |
| Native transport mechanics | `PASSED` in C5b12's local pipe scope | Fixed source/input/start/completion flow |
| Native same-session lifecycle mechanics | `PASSED` in C5b13's benign fixture scope | Concrete launch, teardown and root-custody code for storage integration |
| Durable attempt/completion owner | `BLOCKED` on eight provider bodies and two fixed gate implementations | Durable-before-effect intent, fencing, unresolved cleanup, completion-before-delivery and restart replay |
| Complete immutable runner composition | `BLOCKED` on durable owner, installed/launch identity, restart custody, timing, provenance and independent review | Candidate for a separately authorized controlled typed-transport guest attempt |
| Owner alpha, then external beta | `IN_PROGRESS — TRENDING_GOOD` for owner alpha; external beta `BLOCKED` | End-to-end approved work first; installed signing/distribution, storage continuity and release evidence later |

Next bounded implementation is the durable attempt/completion owner using the existing
fixed-store policy and C5b11 recovery oracle. A test acknowledgment must never substitute
for persisted intent. Teardown must persist cursor 17 before effect 16; lost completion
responses must reopen/replay stored bytes without re-execution or recommit. Installed
identity/custody, preferred-form kernel/libkrunfw source compliance, missing raw v19/v27
evidence and cross-host reproduction remain in their existing workstreams.
