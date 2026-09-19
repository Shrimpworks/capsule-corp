# C5b19 gate-2 exact implementation packet candidate

Date: 2026-09-19. Decision owner: Capsule maintainer. Status:
`IN_PROGRESS — TRENDING_GOOD` while independent packet review is open; gate 2
is not `PASSED`. This freezes a proposed **benign, local-only experiment**; it
authorizes neither implementation nor a child run. Gates 3 and 4 require
separate, explicit owner decisions. The parent owner-only hostile-`.mjs` alpha
remains `IN_PROGRESS — TRENDING_GOOD`; runnable C5b19 evidence and product
admission remain `BLOCKED`.

## Decision and exact boundary

Question: can the [gate-1 pipe/control candidate](C5B19_GATE1_MECHANISM_DECISION.md)
be implemented with a closed fixture, wire format, clock, fault corpus and
independent oracle, without restoring storage to the post-create stop path?
The controlled target is **one new fixed inert direct child** in one
owner-created disposable directory on a separately authorized owned macOS
host. No guest, interpreter, daemon, Broker, backend, user content, network,
external PID, descendant or installed service is in scope. Approval of this
packet is a design review, not permission to create that child. The experiment
archive is `Shrimpworks/capsule-experiments` under a new
`experiments/typed-guest-transport-c5b19-benign-teardown/` directory; C5b16
at immutable commit `0efd03def6bc333a92c8b9809bc56b7b3ce9ea80` is a
read-only reference, never an in-place modification or reinterpreted ABI.

The proposed host profile is local **macOS 26.6.2 build 25G83, arm64**, Xcode
26.6 build 17F113, SDK 26.5 and Apple clang 21.0.0. Node 22.22.1 and pnpm
10.28.2 apply only to this repository's checks; Go 1.25.13 is the exact
declared Go toolchain for any new archive Go module. These are read-only
observations from 2026-09-19, **not** a claim that this host is run-authorized.
Gate 4 must capture exact `sw_vers`, `uname -m`, `xcodebuild -version`, SDK,
compiler, `go version`, module sums, source and executable SHA-256 readback;
any mismatch returns to packet review before a first run. Builds must be
offline from pinned local toolchains/modules; no new package or network fetch.

### Source and build closure

Write a **new** experiment source tree, not imports from this repository's
product packages or the C5b16 driver. Freeze this intended file manifest at
gate 3: `PLAN.md`; `inputs/ORIGINS.json`; `inputs/fixture.c`;
`inputs/supervisor.c`; `inputs/protocol.h`; `inputs/protocol.c`;
`inputs/clock.h`; `inputs/clock.c`; `inputs/bridge.h`;
`source/bridge/bridge.go`; `source/store/store.go`;
`source/store/store_test.go`; `tests/model_test.go`;
`tests/oracle.go`; `tests/oracle_test.go`; and `scripts/build.sh`,
`scripts/check.sh`, `scripts/run.sh`. A generated c-archive header and binary,
the fixture and Supervisor executables, and raw run outputs are separately
listed by hash in the later Gate-4 manifest. These names describe ownership,
not pre-existing bytes. Any new linked source, generated unit, runtime module,
signal handler or process waiter must be added to the manifest and reviewed
before running. The C5b18
[`types.go`](../internal/execution/teardownpassive/types.go) and
[`model.go`](../internal/execution/teardownpassive/model.go) are the normative
transition reference; the archive implements a versioned independent adapter
and tests conformance, not a production import.

`inputs/supervisor.c` owns one control lane and one exact child. It creates
the pipes, starts exactly one storage worker before preparation, samples clocks,
owns the C5b18-like state machine, creates child, accepts cancellation, admits
start, latches stop, calls `waitpid`/`kill` and records raw observations. It
must contain no store/bridge call or bridge/store mutex acquisition while child
custody exists. `source/bridge/bridge.go` exposes only one fixed request/settled
result call to that sole worker; `source/store/store.go` owns one serialized
owner-locked record. Neither Go unit may spawn, signal, wait, write a start
token, set a stop/clock anchor or release capacity. The worker may block in the
store; it cannot own a lifecycle decision. `tests/oracle.go` reads raw evidence
after the run and never calls live Supervisor functions. The build closes
linked C, Go, c-archive/runtime and generated-source inventories at Gate 4,
including every `wait*`, `SIGCHLD`, `signal.Notify`, `SA_NOCLDWAIT` and
auto-reap path. If exclusive same-parent wait status cannot be proven, no
signal or absence claim and no run.

### Fixture, process and descriptor caps

`fixture.c` is one fixed, no-argument, no-network C executable. After reading
exactly one `G` byte from descriptor 3, it waits inertly for Supervisor stop;
before that it blocks at the gate. It never forks, execs, opens a path, reads an environment
value or writes output. `alarm(4)` provides *harness-only containment*, never
Supervisor absence evidence. Close all descriptors except 0/1/2 to
`/dev/null` and read-only gate FD 3 in `posix_spawn` file actions; use
`POSIX_SPAWN_CLOEXEC_DEFAULT`, verify the inherited set in a dedicated
negative fixture test, and set an empty explicit environment. Any platform-
synthesized environment entry must be inventoried at Gate 4; an unexpected
entry fails admission. No pipe, lock, trace, directory or store FD enters the
child. The fixed executable path lies under a private 0700 root and is
read-only after hash readback. A pathname-to-inode swap or digest mismatch
fails admission; this is test-only identity, not installed code signing.

The record's 32-byte `ProcessIdentity` is SHA-256 over a domain-separated
encoding of the consumed attempt ID, successful positive spawn PID, native
creation-event sequence/tick and read-back fixture executable digest. This
opaque record binding is **not** signal authority; only the same live parent,
exclusive unreaped child status and retained positive PID support the
conditional direct-child signal argument. A zero, mismatched or substituted
identity refuses start and signal. Gate 4 must prove the digest inputs come
from actual creation/readback observations, not worker or harness text.

At most one harness, one Supervisor, one sole storage worker, and one direct
child are authorized per attempt. Go runtime threads do not grant extra child
or waiter authority. No descendants, second attempt, replacement process,
unbounded queue or PID adoption. One run directory, one owner lock, one record,
one pending write, one request/reply frame per operation, one cancellation
event, one start-token attempt and one signal attempt are maxima. Fixed
iteration count and named fault matrix are set before a campaign; test loops
cannot create an unbounded number of children. Unexpected extra process,
competing reaper, descriptor leak or escaped root stops the campaign.

The harness creates an owner-only 0700 directory under a user-approved
disposable root, records its canonical path/device/inode, and refuses symlinks
or pre-existing occupants. It cleans only its exact owned directory after
authoritative child status (or clearly labeled harness containment), bounded
worker termination and record snapshot. Never recursively delete a broad
parent. If worker remains blocked or child custody is unresolved, stop the
campaign, preserve evidence and require manual owner cleanup; never call a
second PID-targeting fallback. The child self-alarm at 4 s, harness watchdog
at 7 s, verifier timeout at 9 s are independent finite containment caps from
C5b16, not evidence of the 1,000/1,200-ms Supervisor bounds. Harness timeout
may kill its **owned fixture only** after explicit identity/readback and must
record that as harness action, never Supervisor stop. Gate-4 review must prove
the exact watchdog path cannot target an unrelated PID.

### Fixed transport and operation identity

Two anonymous pipes carry worker request and reply; separate one-event cancel
and start-gate pipes are harness/control and control/child only. Control ends
are nonblocking, every writer has `F_SETNOSIGPIPE`, and pipe capacity is
queried. A frame is written once, wholly, at no more than
`min(512, _PC_PIPE_BUF)` bytes; a short/zero write, `EAGAIN`, `EPIPE`, excess,
duplicate, EOF mid-frame or unknown version is fatal/refusal, not a retry.
Readers assemble one bounded exact frame without trusting a length larger than
the schema maximum. Explicit big-endian integers and byte arrays, no native
padding, pointers, paths, PIDs or process-control flags.

Both frames use an 8-byte header: ASCII `C519` (4), protocol version `1`
(1), role (`request=1`, `reply=2`, 1), total frame length (u16). There are
exactly three kinds: preparation=1, runner-identity=2, terminal-join=3.
Outcomes: confirmed=1, failed=2, indeterminate=3. A request body contains
the 240-byte C5b18 `Bindings` in declaration order (five 16-byte identities
and five 32-byte digests); one 32-byte obligation digest; operation ID (16),
generation (u64), kind (u8), candidate (68), and three clock-policy u64s (24).
The 68-byte candidate is prepared (u8), runner-confirmed (u8), runner identity
(32), absence identity (32), terminal-confirmed (u8), terminal disposition
(u8: none=0, absence-recorded=1, timing-violated=2). Thus request maximum is
**397 bytes** (`8+240+32+16+8+1+68+24`). The reply is the same header plus
operation ID (16), generation (8), kind (1), SHA-256 of the exact request
body (32), outcome (1), and durable-record digest (32): **98 bytes**. The
worker's durable digest is never independent authority to start: control also
requires the exact tuple, request hash, confirmed outcome and its frozen
candidate. A failed or indeterminate reply's record digest is all zero and
cannot be treated as confirmed evidence. The obligation digest commits to
the complete versioned cleanup-only record including bindings and fixed
policy; it is derived before preparation and is identical on all three
operations. Every nonzero/different-ID, generation, candidate and digest
validation rule from C5b18 applies before a request is accepted.

These sizes are proposed **exact maxima**, not asserted implemented bytes.
Gate 4 must compare generated C/Go encoders byte-for-byte against this schema
and independently query `_PC_PIPE_BUF`; a changed field or frame length
reopens Gate 2. Logical pending covers a queued request, executing store call
and unconsumed reply; no second request may begin until the exact predecessor
settles. An indeterminate/preparation failure never creates a child. After
creation a stuck/refused identity operation leaves start closed, but control
still stops/reaps under exact custody. Terminal write begins only after
authoritative absence and a settled predecessor. `BridgeApply` is called by
the worker only, including preparation before spawn. `poll` on nonblocking
control descriptors is capped at 1 ms and the next deadline; `EINTR`
resamples immediately. No read, write, lock, worker join or completion drain
can block the post-create control iteration.

### Durable format and frozen fault sites

Use a new `capsule.c5b19.teardown-obligation/v1` record in the owned private
directory, not C5b16 cursor/record/24-provider ABI. It stores the exact
bindings, obligation digest, fixed clock policy, last settled generation,
operation ID/kind/candidate/outcome, and uncertainty/terminal disposition.
One worker holds an exclusive owner lock; each operation writes a fixed
candidate to a temp file in the same directory, fsyncs file, renames to the
record name and fsyncs directory before confirmed reply. Recovery never
converts an unknown write into confirmation. At most one owner/record/generation
is permitted. The record is cleanup-only and cannot issue content, start,
approval, completion or capacity authority.

Fault injection is configured before an attempt and has no execute-time
parameter. Named sites are (A) before temp-file publication, (B) after rename
before directory fsync, (C) after directory fsync before reply, (D) reply
drop, (E) control delay immediately before final start-admission check, and
(F) control delay after final check before start-token write. A–D apply
separately to identity and terminal operations; preparation uses refusal,
pending and late-reply variants without spawning. Fixed delays are 0, 800,
1,200 or 1,600 ms only; each site fires at most once. Delay or lost reply
must not create a new operation/generation, extend a deadline, change a
candidate or release a start. Unbounded worker hang is a *modeled* unit-test
state only; a runnable fault is finite and the harness timeout still bounds
containment. Any failed fsync, rename, lock or digest check yields refusal or
indeterminate state according to what was durably observed, never guessed
success.

### Clock, raw trace and independent oracle

One `mach_continuous_time()` source supplies every current-lifetime anchor and
observation. Before spawn obtain a nonzero `mach_timebase_info`, enforce
nondecreasing ticks, and compare checked u128 products
`(now-anchor)*numer >= bound_ns*denom`. Reject overflow, tick regression or
mixed clock. The original 1,000-ms setup bound starts at endpoint creation;
check immediately before spawn and again immediately before start, then
record any post-sample scheduler stall as timing violation. `t_start` is the
final pre-token sample; `t_wall=t_start+1,000 ms`. A wall callback's *actual*
service tick is distinct from its anchor. Earliest valid action and actual
service ticks are retained. Stop latency requires both
`t_absent-t_force <= 1,000 ms` and `t_absent-t_action <= 1,200 ms` when those
anchors exist. Before start, cancel or fatal setup expiry provides the action;
no synthetic wall anchor exists. Sleep counts; no budget reset on wake,
storage settlement or repeated trigger. Exact `waitpid` terminal status for
the same unreaped child is the only absence event. `kill` return, `ESRCH`,
`ECHILD`, pipe EOF, child alarm and harness cleanup are not.

The raw append-only trace records sequence number, attempt/binding digest,
event kind, raw tick and ratio, operation/generation/request digest, process
identity, syscall/result/errno, actual cancellation ingress and service,
start-token attempt, signal-attempt latch before `kill`, each `waitpid`
observation, store edge/outcome, and harness-only containment separately.
Fixed capacity is 512 entries; trace overflow fails the case, not a truncated
pass. `tests/oracle.go` computes transitions and bounds from this trace plus
durable snapshot and executable digests; it ignores the Supervisor's claimed
pass flag. Missing event is null/failure. It verifies preparation before
spawn, identity confirmation before start, one pending operation, stop over
same-tick reply, no start after stop, no repeated signal, exact absence,
terminal-after-absence and no public completion. It must reject a forged
success flag, changed tuple/candidate, unpaired terminal record or harness
kill masquerading as Supervisor stop.

Test-first implementation is required at Gate 3: add failing parser/model/
oracle tests before encoders, store and control code; add controlled negative
fixtures for frame corruption, state/clock races and descriptor leakage
before enabling their paths. Red/green evidence stays in the experiment
archive. Gate 4 independently reviews both tests and implemented bytes.
No child test runs before separate first-run authorization. Pre-run unit tests
must use simulated process and store effects only; any test that invokes the
fixture is a Gate-5 run and prohibited until Gate 4 closes.

Run matrix after authorization: normal; preparation pending/refused/late;
identity/terminal A–D at each fixed delay; cancel/setup/wall/fatal races;
natural exit before signal (simulated pre-run, or child self-alarm in a bounded
live negative case, never a Supervisor timing pass); lost signal response;
`ECHILD`, `ESRCH`, changed
disposition and competing-waiter refusal; duplicate/mismatched frame;
suspend/late-service simulation; ordinary and supported sanitizer builds.
Mutants restore stop-to-store wait, preparation bypass, overlapping requests,
stale setup check, late start, second signal, forged absence, changed candidate,
restarted/rounded clock or premature terminal release. Independent oracle must
kill each relevant mutant and retain its failing raw trace. Run only fixed
finite cases; a true suspend/wake or Supervisor-restart probe is **excluded**
from first run because it needs a separately reviewed containment packet.

## Gate-2 review questions and stop rules

Independent reviewer should first challenge this packet blind against merged
Gate-1, C5b18, pinned C5b16 source, current platform contracts and the
ecosystem custom-primitive checklist. Then compare author rationale. Specifically
check (1) wire-size arithmetic and frozen-field completeness; (2) Go/C linked
call graph and no store on stop path; (3) exact child/descriptor/reaper premise;
(4) clock and scheduler-failure classification; (5) bounded cleanup and oracle
independence; (6) whether every run case stays within the authorized benign
scope. Any changed authority, helper, process tree, host or ABI reopens Gate 1
and may require an ADR. A blocking finding keeps Gate 2 `BLOCKED`, even if
repository checks pass. If reviewer accepts exact packet, record reviewed
commit, limitations and Gate-2 `PASSED`, then request **separate owner Gate-3
authorization** naming this packet. No implementation or run follows from a
document merge. Gate 4 requires actual-byte review and another owner decision
before first child execution.

Source classes: documented platform facts and conditional PID argument are in
[Gate 1](C5B19_GATE1_MECHANISM_DECISION.md); source observations are from
[immutable C5b16](https://github.com/Shrimpworks/capsule-experiments/tree/0efd03def6bc333a92c8b9809bc56b7b3ce9ea80/experiments/typed-guest-transport-c5b16-timing-fault/source);
proposed frame/fixture/cleanup details above are **design choices**, not
tested facts. Unknown until Gate 4: generated bytes, actual linked waiters,
host signal disposition, executable hashes, timer service and fixture behavior.
