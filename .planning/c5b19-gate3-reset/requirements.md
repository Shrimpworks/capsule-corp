# RESET-REQUIREMENTS — Gate 3 completion contract

Status: PASSED for this read-only requirements handoff. Parent documentation reset remains IN_PROGRESS — TRENDING_GOOD. Implementation under Gate 3 is incomplete; Gate 4 and child execution remain BLOCKED.

## Method and limits

Read the complete frozen `docs/C5B19_GATE2_IMPLEMENTATION_PACKET.md`, complete Gate-1 mechanism decision and ADR-0047, reset plan, and experiment-plan context plus ordered gates at canonical branch `codex/c5b19-gate3-reset` (base `a9e1d32afbb1fa1ef9c94175321c9c3bbb3d3361`). This assignment derives obligations; RESET-INVENTORY owns claims about actual implementation. No repository edits, tests, builds, native code, child, live store or network actions were performed. Only this report is retained under `/tmp`.

Codebase-memory coverage checked all five cited requirement/plan paths. Index is best-effort generation 2026-10-04, from the primary checkout, not this audit branch. Four requirement documents had no recorded coverage issue; reset plan was missing. Current full source reads, not that old index, ground these requirements. This is a bounded requirements review, not independently verified actual call-graph closure or a new platform-contract review.

## Correct the gate terminology

The ordered gates call Gate 3 **owner implementation authorization**, already PASSED for implementation-only scope. The unfinished work is implementation performed under that authorization. Report these separately: `Gate 3 authorization: PASSED; implementation: IN_PROGRESS — TRENDING_GOOD; Gate 4 implemented-byte review/first-run authorization: BLOCKED`.

A useful delivery milestone is **Gate-3 implementation complete / ready for Gate-4 review**. It requires the entire exact packet implemented, compiled and supported by the permitted pre-run evidence, with no unimplemented packet component or unresolved blocking implementation finding. It is not a new authority gate and does not mean the actual bytes have passed independent Gate-4 admission. The ordinary independent review for each implementation deliverable also does not replace that full Gate-4 review.

Gate 4 independently checks exact source/generated/binary/build identities and all linked custody, signal, lock, clock, descriptor, cleanup and oracle premises, then closes blocking findings. Only after that does the owner separately authorize a first run on the named owned host. Gate 5 executes the authorized finite corpus and retains raw evidence; Gate 6 records the canonical conclusion. A synthetic pass never claims actual timing, native custody, durable publication, descriptor inheritance or product admission.

## Proposed fixed delivery set: four outcomes

Use four public deliverables, with internal test-first packages as needed. Do not give every additional helper/check a separate public roadmap slice. Each deliverable retains exact source pin, red/green evidence, appropriate permitted tests, independent review and its known limitations. Scope expansion needs a reasoned change to this finite plan.

### D1 — Complete synthetic lifecycle oracle and scenario corpus

**Outcome:** the existing partial validators compose into one independent offline judgment of a complete supplied attempt, rather than passing individually selected fragments.

**Owns:** `tests/model_test.go`, `tests/oracle.go`, `tests/oracle_test.go`; integrated raw-frame, lane, snapshot, identity and executable-digest input contract. Maintain a versioned independent adapter to C5b18; no product imports or calls into live Supervisor functions.

**Acceptance:** complete synthetic success and each applicable failure scenario derive preparation/spawn/identity/start/stop/signal/wait/absence/terminal ordering from raw observations. Check current-lifetime identity, exact wait-result-to-absence causality, one pending operation, stop precedence, no late start or repeated signal, terminal after absence/predecessor settlement, and permanently closed public completion/capacity. Derive setup, service, wall, action and force timing with exact tick arithmetic; missing or contradictory evidence cannot be a pass. Require scenario-appropriate lane, frame, snapshot and digest completeness; preparation refusal legitimately has no child observations. Worker/harness observations cannot substitute for control application or exact fixture absence. Same-tick decisions follow control sequence, not global cross-lane timestamp sorting.

Exercise corruption of every frozen field, known-answer cancellation/obligation/process/record bytes, omitted/reordered causal events, changed candidates, forged success, mixed/regressed clocks, late final/post-start service, custody failures, and the packet's relevant modeled restoration mutations. Retain failing synthetic traces and checker results. Include no-child trace-capacity/pacing cases: 6,999 paired waits plus 512 other control events, full budget/overflow, readiness consumption, EOF removal and EINTR storms. The integrated corpus is finite and enumerated; a passing helper or guard test alone is not evidence that an integrated scenario or executed mutant was caught.

**Dependencies/output:** builds on audited existing parser/model/sidecar work. Outputs complete oracle contract, scenario/mutation matrix with expected verdicts, and synthetic evidence. This provides the independent test target for D2/D3. Standalone executable packaging remains D4.

### D2 — Frozen C/Go wire and sole-worker durable implementation

**Outcome:** the actual cross-language request/result and storage path implements the frozen data contract, with no lifecycle authority.

**Owns:** `inputs/protocol.h`, `inputs/protocol.c`, `inputs/bridge.h`, `source/bridge/bridge.go`, `source/store/store.go`, `source/store/store_test.go`, and their applicable D1 conformance vectors. Worker transport integration joins D3; no new transport, generic IPC or dependency.

**Acceptance:** C/Go encode/decode correspondence for 397-byte request, 98-byte reply, 57-byte cancellation, 395-byte record and all immutable digest preimages; exact fields, byte order, zero rules and bounded malformed-input refusal. One serialized owner lock/record and one fixed pending operation; simulated temp-file write, file fsync, rename, directory fsync and reply observations establish confirmation only for the exact candidate. Failure/unknown outcome remains uncertain even when record bytes are visible. Cover preparation refused/pending/late and identity/terminal fault sites A–D, once per configured site with frozen delays 0/800/1200/1600 ms modeled rather than slept. No replacement candidate, operation, generation or inferred settlement.

Implement real store operations for later review/run, but pre-run tests use **simulated store effects only**. Go bridge exposes one fixed request/result call to the sole worker; no spawn, signal, wait, start token, clock/stop anchoring, completion or capacity authority. Review dependencies/imports/generated interface and store/bridge locks as implementation evidence, while full linked-byte closure remains Gate 4.

**Dependencies/output:** D1 byte/state contract; output complete source and simulated failure/restoration evidence. Internal packages may proceed alongside late D1 scenario integration once contracts are fixed; acceptance requires D1 correspondence tests.

### D3 — Native Supervisor, clock and trace implementation

**Outcome:** the full proposed control mechanism is present and exercised through simulated process/pipe/store observations, including the stop path when worker progress is unavailable.

**Owns:** `inputs/supervisor.c`, `inputs/clock.h`, `inputs/clock.c`; native worker plumbing and lane-local trace/sidecar production. Integrates D2 without moving authority.

**Acceptance:** configure bounded pipes/nonblocking control/F_SETNOSIGPIPE/queried PIPE_BUF before worker/child; one precreated worker performs preparation. No create on refused/late/pending preparation. Actual code owns exact same-parent custody, fresh pre-spawn and final start checks, separate sampled/actual service ticks, cancellation-before-reply precedence, one latched signal, exact wait-status absence and post-absence terminal eligibility. No control-to-Go/store/lock/join/allocation/file-I/O path during custody. All syscall failures and fixed faults A–F leave explicit refusal/uncertainty; only control decides application/start/stop.

Clock uses one nonzero continuous-time ratio, private control watermark and checked u128 products; original 1,000-ms setup/wall, 1,000-ms force-to-absence and 1,200-ms action-to-absence bounds never restart or round away. Tests simulate boundary ticks, scheduler stalls, suspend/late service, regression/overflow and both same-tick orderings.

Produce exact 256-byte events in independent arrays (32,768/1,024/64 entries), separate writer sequences/clock state, capped frame captures, and fail-closed overflow. Worker pre-fault flush cannot block control; control serialization waits until exact absence and never waits for worker. Simulate one-shot readiness, EOF removal, EINTR, every wait-call/result pair, full budgets and stuck worker. Post-absence request close, final reply consumption and at-most-2-second EOF wait must not turn EOF into settlement or use an unbounded join.

**Dependencies/output:** D1 oracle/corpus plus D2 worker contract. Output complete native source, simulated syscall/state/trace evidence checked by the independent oracle, and an explicit source call-graph/lock/authority inventory for later Gate-4 review. A modeled syscall result is not native custody/timing evidence.

### D4 — Containment harness, build closure and Gate-4 handoff

**Outcome:** the entire frozen experiment can be independently reviewed as exact buildable bytes with a finite, gated run plan; nothing runs the fixture yet.

**Owns:** `inputs/fixture.c`, `inputs/harness.c`, `scripts/build.sh`, `scripts/check.sh`, `scripts/run.sh`, and complete `PLAN.md`, `go.mod`, `inputs/ORIGINS.json` provenance. Integrates all D1–D3 source, generated header/c-archive and fixture/Supervisor/harness/oracle executable artifacts.

**Acceptance:** fixture has only inert gate behavior, pre-read 4-second alarm, no arguments/network/fork/exec/path/environment/output; spawn code enforces explicit empty environment and only /dev/null 0/1/2 plus gate FD3, default/unblocked SIGALRM and required CLOEXEC flags. Simulated descriptor/alarm/identity-negative tests exist before paths are enabled; any test actually launching the fixture remains Gate 5.

Harness is sole Supervisor parent/waiter; Supervisor sole fixture parent/waiter. `run.sh` validates a finite case ID and execs harness, which after reaping/containing Supervisor execs read-only oracle. No shell/oracle extra waiter or target probe. Implement exact 0700 owned-root/path/device/inode and hash-readback checks, one permitted watchdog action against its own live unreaped Supervisor at 7 seconds, no fixture PID fallback, safe bounded shutdown and preserve-on-uncertainty/manual-cleanup outcome. Oracle's separate 9-second in-process bound grants no lifecycle authority.

Compile offline using exact approved local pins without running fixture. Produce complete source/build/module/generated/linked-runtime inventory and executable hashes for the Gate-4 pre-run manifest; no extra dependency or implicit fetch. Provide every waiter/SIGCHLD/SA_NOCLDWAIT/auto-reap/init/runtime path, descriptor and executable identity premise, lock/clock path, oracle independence and cleanup checklist for independent review. Include ordinary and supported sanitizer build configuration, finite post-authorization case/fault/mutation matrix, and two-directory reproducibility procedure; actual run/reproduction evidence remains Gate 5. Safely verify run refusal on missing/changed authorization/digests/environment without entering a child path.

Packaging may require reviewed refactoring: the current oracle package is not itself a standalone executable entrypoint. D4 must reconcile executable composition and any test-injection boundary with the frozen manifest before claiming build closure. Do not require perpetual append-only source prefixes; preserve immutable historical commits, accepted regression semantics and provenance while allowing reviewed internal refactoring. Do not invent an extra linked entrypoint, wrapper, test helper or authority path outside the packet; an out-of-packet change reopens the applicable review. This is a completion obligation, not a selection of implementation design.

**Dependencies/output:** D1–D3 complete. Output the exact implemented candidate, permitted pre-run red/green evidence, inventory-to-packet coverage with no missing item, build/digest evidence, deliberate limitations and a Gate-4 review packet. Handoff records Gate 4 and first child run still BLOCKED; it neither self-admits binaries nor supplies first-run authority.

## Packet coverage / non-negotiable dependencies

| Frozen obligation | Deliverable | Primary packet sections |
| --- | --- | --- |
| New archive tree, no product/C5b16 import, all 20 named manifest entries, pinned no-dependency Go module | D4 integrates D1–D3 | Source/build closure, lines 49–98 |
| Closed fixture, exact process identity, descriptor/environment/reaper premises | D3 implementation + D4 containment/closure; D1 oracle | Fixture/process/descriptor caps, lines 100–152 |
| Root/lock/record/process/attempt caps, 2-s worker exit, 4/7/9-s independent containment/verifier roles | D2/D3 + D4 | Lines 131–185 |
| Exact frames, obligation, candidate, cancellation, tuple/pending semantics and native Go/C vectors | D2 with D1/D3 integration | Lines 187–250 |
| Durable publication, unknown outcome and A–D finite fault sites; E/F start races | D2 + D3; D1 evidence | Lines 252–278 |
| Continuous ticks, exact ratios/bounds, no scheduler guarantee | D3 + D1 | Lines 280–298 |
| Closed 29-kind/3-lane raw trace, captures and correspondence, causal rather than cross-lane order | D1 + D3; D4 harness lane | Lines 300–435 |
| Complete independent state/timing oracle and finite negative/restoration corpus | D1 integrated with D2/D3 | Lines 437–465 |
| Full generated/linked source/binary/readback and review/admission boundary | D4 prepares; Gate 4 independently admits | Lines 49–98, 467–491; experiment ordered gates 3–6 |

## Permitted testing and scope-change rules

Pre-run tests are pure parsers/encoders/models/oracles and native/Go unit harnesses using simulated process **and store effects**. Simulate custody, waits, signals, descriptors, clock advancement, filesystem failures and worker delays; unbounded worker hang exists only as modeled state. Compiling/hashing is distinct from fixture execution. Real fixture launches, even descriptor-negative or alarm tests, belong to Gate 5 after Gate 4 and explicit first-run authorization. Do not silently reinterpret simulated tests as permission to perform real store fault campaigns, native lifecycle probes or sleeping/live suspend experiments.

Do not move Gate-4 actual linked-byte/host inspection into a supposedly PASSED Gate-3 admission claim. Gate-3 author prepares evidence; Gate-4 independent review judges it. Host/toolchain/environment disagreement is a packet-review stop, not permission to substitute the current machine.

A new helper/lifecycle owner, process tree, authority, host or ABI reopens Gate 1 and may require ADR work. Changed wire fields/lengths or dependency/module requirements reopen Gate 2. Added linked/generated/runtime source, signal handler or waiter must enter the manifest and independent review before running; material mismatch also reopens packet review and owner implementation authorization. Shared-lock transport, pipelining/multiple attempts, product imports, restart/PID adoption or descendants are not implementation shortcuts within this packet. Real suspend/wake or Supervisor restart is excluded from first run and requires a separately reviewed containment packet. Changing fixture or widening clocks after an adverse result requires a new reviewed packet.

## Recommendation to Brain

Next highest-priority deliverable is D1, because integrated scenario judgments replace the present succession of partial-validator milestones and give native work a finite conformance target. D2 and D3 may have internal test-first packages and parallel source work where contracts allow; publish outcomes only against the four explicit exits. Do not promise a fixed PR count, elapsed time or percentage complete. A deliverable can require correction/review commits; it must not silently create another public micro-slice. Reconcile RESET-INVENTORY evidence into implemented/partial/missing rows before assigning any PASSED implementation status.
