# C5b16 benign timing/fault checkpoint

Date: 2026-09-06

Status: construction, verification and independent review `PASSED` / **Ready**
(review instance 2 of 3).
Parent owner-only internal alpha: `IN_PROGRESS — TRENDING_GOOD`.
Installed lifecycle, controlled guest execution and product admission: `BLOCKED`.

## Question and scope

The merged [focused audit](C5B_PRE_C5B16_SECURITY_AUDIT.md) at `cf835768c3cb04bf6ad9ad3cb2ed8e90ab313db5`
selected the [C5b15 timing probe](C5B_RUNNER_RECONCILIATION_CHECKPOINT.md).
C5b16 asks how durable publication delay affects the existing fixture's setup and
teardown clock accounting. It uses a versioned derivative of C5b14B source
`294559e78b410f4b25e0e25fcf373ec0c7d1cfb5`, on archive main
`cbfef30751f0a2dd8d8a71356c6b975b52d3e752`.

Defensive execution is confined to owned disposable local directories and the
fixed benign C child with a Go 1.25.13 c-archive. No guest, interpreter, backend,
installed service, signing identity, private key or user content is involved.
No product source, authority, resource limit, durability order or ADR lifecycle changes.

## Mechanism and evidence interpretation

A bounded 512-entry native `CLOCK_MONOTONIC` trace observes serialized phase and
publication-edge events. The drain thread does not write the trace. One trusted
fixture-only delay can be armed before the first drive, for spawn intent or
teardown intent at one of the three existing publication edges. Its maximum is
1,600 ms; it cannot be rearmed after use. It is not an execute-time plan option.

The unchanged setup clock begins at endpoint creation and lasts 1,000 ms. The
unchanged cleanup clock begins after the durable teardown gate returns and lasts
1,000 ms. The comparison initial action is fixture entry to the checked teardown
request, before the durable gate; it is not a measured guest wall/cancel dispatch.
The [C2A contract](protocol/GOVERNED_DENO_CORE_C2A_EXECUTION_PROFILE.md) separately
requires forced absence within 1,000 ms and a 1,200-ms maximum from initial action.
The four-second child self-alarm and seven-second harness stop are containment
for this experiment only.

Raw observations distinguish gate latency, post-gate cleanup, total time to native
absence, and harness-only reap after refusal. Missing native absence stays null;
it is never replaced with a self-alarm or harness observation. Finite injected
sleep measures an artificial publication stall, not a kernel disk hang or a
latency distribution. Observed elapsed time includes scheduler overhead.

## Observed timing

The retained matrix contains 19 cases in ordinary and parent-C ASan/UBSan builds
(38 runs). These ranges combine both variants and all three publication edges;
they are observations from this corpus, not latency guarantees or percentiles.

| Requested delay (ms) | Gate elapsed (ms) | Post-gate to native absence (ms) | Total teardown to native absence (ms) |
| --- | --- | --- | --- |
| 0 | 10.053–14.461 | 29.915–39.909 | 39.968–54.370 |
| 800 | 812.009–819.302 | 30.256–50.652 | 845.228–869.954 |
| 1200 | 1213.131–1219.914 | 30.485–42.154 | 1247.739–1258.763 |
| 1600 | 1614.221–1637.714 | 30.195–72.075 | 1644.416–1709.789 |

Every retained 1,200-ms and 1,600-ms successful gate delay exceeded the 1,200-ms
initial-action bound, while its post-gate cleanup remained below 1,000 ms. The
existing synchronous gate plus a newly started cleanup clock therefore does not
satisfy the total timing requirement for these controlled delays.

All four refusal cases per variant (the initial no-delay pre-publication refusal
and delayed refusals at all three publication edges) retained `spawn=1`, `kill=0`,
phase 2, unresolved cursor 17 and no completion. Native absence was unobserved;
the harness subsequently reaped a child terminated by its own `SIGALRM`.

A 1,200-ms delay in spawn-intent publication at each edge prevented any spawn.
Both retained 800-ms spawn-delay cases spawned once. The ordinary case refused
at readiness (effect 3) after observed setup expiry and completed conservative
cleanup; the sanitizer case completed normally. This is not a latency guarantee.
The initial exploratory test incorrectly demanded success from every sub-deadline
spawn delay; its failed assertion is retained without an inferred cause. The final
checker permits success or conservative refusal only with observed setup expiry.
The separate zero-delay normal case still requires exact completion and delivery.

## Verification and retained evidence

Timing runs and three compiled timing mutations plus seven false-evidence mutations
are `PASSED`. Removing teardown-request or absence observations, omitting the
injected sleep, inventing absence or a total-bound pass, reversing time, removing
a case or refusal observation, substituting a sanitizer artifact identity, or
claiming product admission refuses. Each run binds its ordinary or sanitizer
executable digest. Path-sensitive debug sanitizer hashes are retained but excluded
from deterministic reproduction. Exact source and artifact hashes bind the runs.

The [immutable archive](https://github.com/Shrimpworks/capsule-experiments/tree/0efd03def6bc333a92c8b9809bc56b7b3ce9ea80/experiments/typed-guest-transport-c5b16-timing-fault) retains the plan, predecessor origin ledger,
source, raw timing and regression evidence, exploratory failures and both review
packets. Full inherited regression passed 126 native cases, 17 Go suites with
58 subtests, eight Go-race integration cases, ten compiled control mutations,
24 provider exports and 22 reproduced artifacts. Independent review instance 2
returned **Ready** after the first review's sanitizer-identity and missing-refusal
findings were corrected and rechecked.

Required canonical Node/pnpm checks, Go test/vet/build and pinned govulncheck passed.
Full golangci reports 50 pre-existing revive documentation findings; the blocking
non-revive gate and new-code revive ratchet passed. Parent-run check results and
reviewer-run verification are identified separately in the archive review packet.

The probe uses no new dependency. Reuse follows the native test/platform and narrow
Supervisor transport rows in the ecosystem adoption map; the archive plan retains
the test-only dependency/authority checklist and exact predecessor origin ledger.

## Next goalpost: teardown-intent and deadline decision

Before further runnable composition, prepare a passive design decision that names
exact wall/cancel and forced-absence anchors, identifies every storage operation
on the critical teardown path, and covers pending, failed, and indeterminate
publication while a child is live. Compare feasible Supervisor-owned ordering
options, including whether sufficient durable teardown intent can exist before
launch. Require exact authority/attempt binding and crash/restart/unresolved-state
consequences for each option. This is a question to resolve, not an adopted design.

Any selected change to durable-before-effect ordering, responsibility, or recovery
semantics requires the applicable ADR and review before implementation. Timing
observability alone does not fix the issue. Broader launch identity, protected
state, and restart custody remain separate gates after that decision.

## Decision boundary

The bounded evidence slice and independent review are `PASSED`.
A post-gate timer alone cannot establish the total initial-action bound. Later
work must reconcile durable-before-effect ordering with timely teardown and
publication-failure handling before runnable composition is promoted. A remedy
that changes authority, durability ordering or architecture needs a separate
explicit decision gate. No such remedy is selected or implemented here.

Installed executable/store protection, rollback resistance, restart/process-tree
custody, real guest controls, sleep/wake/power-loss and cross-host timing evidence
remain `BLOCKED`. Earlier unclassified C5b14B setup refusal remains unclassified.
