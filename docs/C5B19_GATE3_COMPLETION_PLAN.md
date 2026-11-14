# C5b19 Gate 3 implementation completion plan

Date: 2026-11-13. Owner: Brain orchestrator.

This reset replaces the open-ended sequence of public validation micro-slices with
four deliverables for the existing [frozen Gate-2 packet](C5B19_GATE2_IMPLEMENTATION_PACKET.md).
It changes scheduling and acceptance reporting, not architecture, ABI, dependencies
or execution authority. The earlier parser and correspondence work remains useful;
it has not yet produced a complete experiment or an integrated evidence verdict.

## Status and evidence boundary

| Scope | Work status | Meaning / next action |
| --- | --- | --- |
| Gate-3 owner implementation authorization | `PASSED` | Existing implementation-only authority for the integrated exact packet; no new authorization is claimed here. |
| Implementation under that authority | `IN_PROGRESS — TRENDING_BAD` at reset | Useful partial checks, but repeated fragmentation left lifecycle integration and all native components unfinished. Execute D1–D4 against their exits; revise trend when integrated evidence demonstrates progress. |
| Gate 4 implemented-byte review and first-run authorization | `BLOCKED` | Depends on the complete candidate and independent actual-byte review, then a separate owner first-run decision. |
| Gates 5–6, child execution and product admission | `BLOCKED` | No fixture has been run by this reset; this plan grants no run or product authority. |

The reset's inventory and requirements audits are `PASSED` for read-only source
inspection. They are not a new test pass or recertification of historical reviews.
The owner-only product parent remains `IN_PROGRESS — TRENDING_GOOD` on its broader
roadmap; its installed security boundary and admission remain `BLOCKED`.

Audited canonical baseline: `a9e1d32afbb1fa1ef9c94175321c9c3bbb3d3361`, the head of
[PR #410](https://github.com/Shrimpworks/capsule-corp/pull/410), open at reset.
The reset PR is stacked on that branch and does not merge it. Audited archive:
[`a9c457cb3b89a4ecd8c38684fd3d804d2089214c`](https://github.com/Shrimpworks/capsule-experiments/tree/a9c457cb3b89a4ecd8c38684fd3d804d2089214c/experiments/typed-guest-transport-c5b19-benign-teardown),
in existing draft archive [PR #38](https://github.com/Shrimpworks/capsule-experiments/pull/38).
No archive bytes change in this reset. The [inventory](../.planning/c5b19-gate3-reset/inventory.md)
retains exact source/test anchors; the [requirements audit](../.planning/c5b19-gate3-reset/requirements.md)
maps packet obligations and review-reopening conditions. The graph did not cover
the archive; immutable Git tree and source reads establish this map.

## Implemented, partial and missing

There are **six present and fourteen absent paths** in the packet's twenty-path
manifest. This is file inventory, not a completion percentage. Only `go.mod` is
complete for its narrow declaration role; none of D1–D4 is complete.
All paths below are relative to the pinned experiment directory.

| Required paths | Assessment | Delivery owner |
| --- | --- | --- |
| `go.mod` | Implemented declaration: Go 1.25.13, standard library only; full build closure missing | D4 |
| `PLAN.md`, `inputs/ORIGINS.json` | Partial scope/provenance history; complete implemented/build handoff missing | D4 |
| `tests/model_test.go`, `tests/oracle.go`, `tests/oracle_test.go` | Partial pure validators and synthetic histories; integrated lifecycle/evidence verdict missing | D1; D4 executable packaging |
| `inputs/protocol.h`, `inputs/protocol.c`, `inputs/bridge.h` | Missing C codec and fixed bridge ABI | D2 |
| `source/bridge/bridge.go`, `source/store/store.go`, `source/store/store_test.go` | Missing sole-worker bridge, durable implementation and fault tests | D2 |
| `inputs/supervisor.c`, `inputs/clock.h`, `inputs/clock.c` | Missing native control, clocks and trace producer | D3 |
| `inputs/fixture.c`, `inputs/harness.c` | Missing inert child and containment harness | D4 |
| `scripts/build.sh`, `scripts/check.sh`, `scripts/run.sh` | Missing reproducible build, integrated checks and gated case entry point | D4 |

The following narrow components are implemented with corresponding tests in the
archive; their integrated roles remain partial:

| Implemented component | What remains |
| --- | --- |
| Cancellation/process/obligation bytes; request397/reply98/record395 parsing and hashes | Serialized evidence ingestion, C/Go correspondence and observed-effect integration |
| One-pending/three-operation synthetic write histories | Full control lifecycle and independent C5b18 model conformance |
| Bounded event envelopes and field/lane shapes for **all 29 kinds** | Lifecycle meaning, diagnostic variants, complete capture and timing interpretation |
| Captured-frame, received/applied, pending-decision and durable-snapshot correspondence | Compose a complete attempt; detect required omissions and reconcile stop/admission/terminal state |
| Worker-received/sent, write/result and transport-lane correspondence | Whole-stream causality, syscall outcome semantics and required-event completeness |

For example, the lane helper intentionally accepts an empty selection or a lone
kind-8 event. That is valid for a local predicate, but cannot establish a complete
attempt. Current `tests/oracle.go` is `package oracle`, without an executable entry
point. No native custody, actual timing, real durable publication or runnable
experiment is established by these helpers. Historical G3a review-provenance gaps
remain recorded; this reset does not recertify them.

## Four remaining deliverables

Brain owns integration and acceptance for all four. Assigned workers own named
files within that delivery; independent reviewers do not implement their fixes.
Each deliverable requires tests first, retained red/green and relevant mutation
results, fresh independent review, an immutable archive pin and a PR handoff.
Statuses below describe dependencies, not abandoned candidates. Internal coding
packages may be small; a new helper is not automatically a new public milestone.

### D1 — Integrated offline oracle and finite scenario corpus

Status: `IN_PROGRESS — TRENDING_GOOD` for reusable foundations and the now-bounded
integration target; acceptance remains open. Dependency: existing frozen packet
and audited helpers. This is the next highest-priority deliverable.

Owns `tests/model_test.go`, `tests/oracle.go`, `tests/oracle_test.go`: one bounded,
read-only input contract and independent verdict for a complete supplied attempt,
with a versioned C5b18 adapter. No calls into live Supervisor functions.

Exit criteria:

- Compose existing checks over scenario-appropriate complete lanes, raw frames,
  snapshots, identity and executable digests. Missing/contradictory evidence fails;
  a valid preparation refusal legitimately lacks child observations. Ignore claimed
  Supervisor success; worker/harness observations cannot establish control application.
- Derive preparation/spawn/identity/start/stop/signal/wait/absence/terminal transitions,
  exact current-lifetime identity and wait-to-absence causality; one pending operation,
  stop precedence, no late start/repeated signal, terminal only after absence and
  predecessor settlement, and permanently closed public completion/capacity.
- Derive setup/service/wall/action/force bounds using exact ticks and checked
  arithmetic. Use lane-local sequence and cross-lane identity joins, never a global
  timestamp sort. Include the formerly proposed same-operation ordering slice here.
- Enumerate synthetic success/refusal/failure scenarios and mutations: frozen-field
  corruption, omitted/reordered causal events, replacement candidates, forged success,
  regressed/mixed clocks, delayed service, custody failures and modeled restoration.
  Retain integrated failing traces and actual checker/mutant outcomes, not only guards.
- Specify the independent expected verdicts for D2/D3 simulations, including 6,999
  paired waits plus 512 other control events, full budgets/overflow, readiness/EOF
  consumption and EINTR storms. D3 must exercise its actual loop against those cases.

Output: integrated oracle contract, finite scenario/mutation matrix and synthetic
verdict evidence. Standalone executable/build packaging is D4; native evidence is
not credited here. All process and store effects in these tests are simulated.

### D2 — Frozen wire, bridge and sole-worker durable implementation

Status: `BLOCKED` for acceptance on D1's byte/state contract. Owns the six D2 paths
in the inventory; depends on D1 conformance vectors. Internal implementation can
proceed once those contracts are fixed; acceptance includes their integration.

Exit criteria:

- C/Go correspondence for exact 397/98/57/395-byte request/reply/cancellation/record
  formats, preimages, field/zero rules and malformed refusal. Single whole-frame
  writes; short/zero/error outcomes refuse rather than retry or infer success.
- Implement the real serialized owner-lock/temp-write/file-fsync/rename/directory-fsync
  store path for later review, but test it pre-run with **simulated store effects**.
  Visible record bytes alone never confirm a candidate; unknown results stay uncertain.
- Test refused/pending/late preparation, identity/terminal fault sites A–D and the
  fixed 0/800/1200/1600-ms delays as simulated time, with no candidate replacement,
  pipelining, generation substitution or inferred settlement.
- Fixed bridge request/result call serves one worker. No spawn/signal/wait/start,
  stop-anchor/completion/capacity authority in Go. Retain dependency/import/interface
  and lock inventories and D1 correspondence results.

Output: complete wire/store/bridge source and simulated fault/restoration evidence.
Worker transport integration belongs to D3; full linked-byte admission belongs to Gate 4.

### D3 — Native control, clocks and trace production

Status: `BLOCKED` for acceptance on D1 oracle/corpus and D2 worker contract. Owns
`inputs/supervisor.c`, `inputs/clock.h`, `inputs/clock.c` and native worker plumbing.

Exit criteria:

- Implement preconfigured bounded pipes, nonblocking control, F_SETNOSIGPIPE and
  queried PIPE_BUF, one precreated worker and preparation-before-create. Control
  alone owns custody/application/start/stop; no creation after refused/late/pending prep.
- Exercise actual control code with simulated syscalls: fresh pre-spawn/final start
  checks, cancellation-before-reply precedence, separate sampled/actual service ticks,
  one signal, exact wait/terminal-status/absence and post-absence terminal eligibility.
  Include fixed A–F faults, syscall failures, stuck worker and both same-tick orderings.
- No control-to-Go/store/lock/join/allocation/file-I/O path during custody. Verify
  this in source/call-graph review; worker stalls cannot stall the stop decision.
- One nonzero continuous-time ratio, private control watermark and checked u128
  products preserve original 1,000-ms setup/wall, 1,000-ms force-to-absence and
  1,200-ms action-to-absence bounds. Simulate regression, overflow, stalls and late
  service without restarting bounds or claiming a scheduler guarantee.
- Emit exact 256-byte records in independent 32,768/1,024/64-entry arrays, bounded
  captures and fail-closed overflow. Run D1 capacity/readiness/EOF/EINTR cases through
  the control loop and oracle. Serialize control only after absence, without waiting
  for worker; final reply/EOF handling has a two-second bound and no unbounded join.

Output: complete native source, simulated trace verdicts and explicit custody,
waiter/signal, lock, clock and authority inventories. This proves simulated behavior
and inspectable implementation, not actual native custody or timing.

### D4 — Containment, build closure and Gate-4 handoff

Status: `BLOCKED` on D1–D3 completion. Owns fixture/harness/scripts/provenance paths
and integrates all twenty manifest entries, generated headers and executables.

Exit criteria:

- Implement inert fixture with pre-read four-second alarm, no arguments/network/
  fork/exec/path/environment/output; explicit empty spawn environment, /dev/null
  descriptors 0/1/2 plus gate FD3, default/unblocked SIGALRM and required CLOEXEC.
  Descriptor/alarm/identity-negative tests simulate effects; no fixture launch.
- Harness alone parents/waits for Supervisor; Supervisor alone parents/waits for
  fixture. Finite-case `run.sh` execs harness, which contains/reaps Supervisor before
  execing the read-only oracle. Verify owned 0700 roots/path/device/inode/hashes,
  one seven-second watchdog action against its own live unreaped Supervisor, no
  fixture-PID fallback and preserve-on-uncertainty/manual cleanup. Oracle's separate
  nine-second bound grants no lifecycle authority.
- Build/hash offline with exact approved pins and no implicit dependency fetch.
  Close all source/module/generated/linked-runtime and executable inventories;
  compile ordinary/supported-sanitizer configurations without running the fixture.
  Reconcile oracle executable packaging and test injection with the fixed manifest.
- Prepare complete source/trace/corpus evidence and source-to-packet coverage with
  no missing implementation item. Include waiter/SIGCHLD/SA_NOCLDWAIT/auto-reap/runtime,
  descriptor, lock/clock, oracle independence and cleanup premises for Gate-4 review.
- Verify missing/changed authorization/digest/environment refusal before entering
  any child path. Retain finite post-authorization case/fault/mutation matrix and
  two-directory reproduction procedure; actual run/reproduction remains Gate 5.

Output: exact buildable candidate and complete pre-run review packet. Ordinary
independent deliverable reviews do not replace Gate 4's full actual-byte review or
its separate owner first-run authorization. D4 completion cannot authorize a run.

## Coverage, finish line and change control

| Packet obligation | Accountable deliverable / integration |
| --- | --- |
| Source manifest, toolchains, dependencies, generated/runtime closure | D4 integrates all D1–D3 paths |
| Fixture/process/descriptor/environment and containment premises | D3 custody + D4 fixture/harness; D1 judgments |
| Fixed wire/preimages/pending candidate and durable A–D faults | D2; D1 verdicts and D3 transport |
| Start races E/F, stop/custody/terminal behavior and clock bounds | D3; D1 independent scenario verdicts |
| All 29 kinds, independent lanes, raw sidecars, capacity/pacing | D1 evidence contract + D3 producers + D4 harness lane |
| Independent full lifecycle/timing oracle and restoration corpus | D1, exercised against D2/D3 and packaged by D4 |
| Complete pre-run source/build/digest/cleanup handoff | D4; independent admission remains Gate 4 |

Implementation is complete only when **D1–D4 all meet their exits**, every frozen
packet requirement has implementation and permitted pre-run evidence, and blocking
implementation findings are closed. Gate 4 reviews that complete candidate; missing
Gate-3 implementation cannot be deferred into review. Gate 5 subsequently runs only
an explicitly authorized finite corpus; Gate 6 records the supported conclusion.
Even a later benign single-child pass does not admit descendants, restart, guests,
installed product paths or a stronger security claim.

Preserve immutable historical evidence and accepted regression behavior. Reviewed
integration/refactoring is allowed; an append-only source constraint must not prevent
executable composition. Any manifest/ABI/dependency/authority/host mismatch reopens
the applicable packet or mechanism review and authorization before proceeding.
New lifecycle owners/process trees require Gate-1 review and potentially an ADR.
Real suspend/wake or Supervisor restart remains outside the first-run packet.

Brain reports progress by accepted exits, remaining failed/missing scenarios and
next dependency, not PR/helper/test counts or file-based percentages. Correction
commits and internal test-first packages may be necessary; this is four deliverables,
not a promise of exactly four PRs. Adding a public deliverable requires an explicit,
justified update to this plan. Historical G3a next-slice recommendations are superseded
by D1–D4; their code, pins, review results and limitations remain historical evidence.
