# C5b19 gate-2 exact implementation packet candidate

Date: 2026-09-19. Decision owner: Capsule maintainer. Gate 2: `BLOCKED`
until independent review of this corrected packet closes; parent workstream:
`IN_PROGRESS — TRENDING_GOOD`. Review instance 1 of 3 on
`6d3ad7be7319c0f3ac32691c1c56fff5df4d8702` returned **Not ready** with
four accepted P2 packet gaps. Instance 2 on
`ae2b42ba68b9efe81933654ffa855498a3fcfa81` returned **Not ready** with
two further accepted P2 trace-budget and harness/verifier-source gaps.
Instance 3 on `f7c9016982a319ad56fc4853968501d1af46f3d1` returned
**Not ready** with accepted P2 causal-event encoding and P3 wait-slot
arithmetic findings. Corrections below are **unreviewed**: configured
three-instance limit is exhausted, and another independent review requires
an explicit owner decision. This
freezes a proposed **benign, local-only experiment**; it
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
product packages or the C5b16 driver. Gate 2 freezes this intended file
manifest for Gate-3 implementation: `PLAN.md`; `go.mod` (Go 1.25.13, no
external requirements; no `go.sum` unless a new reviewed dependency reopens
Gate 2); `inputs/ORIGINS.json`; `inputs/fixture.c`;
`inputs/supervisor.c`; `inputs/harness.c`; `inputs/protocol.h`;
`inputs/protocol.c`;
`inputs/clock.h`; `inputs/clock.c`; `inputs/bridge.h`;
`source/bridge/bridge.go`; `source/store/store.go`;
`source/store/store_test.go`; `tests/model_test.go`;
`tests/oracle.go`; `tests/oracle_test.go`; and `scripts/build.sh`,
`scripts/check.sh`, `scripts/run.sh`. Generated c-archive header and binary,
fixture, Supervisor, harness and oracle executables are listed by hash in the
later Gate-4 pre-run manifest; raw run outputs appear only in the Gate-5
evidence manifest after authorization and execution. These names describe
ownership,
not pre-existing bytes. `inputs/harness.c` is the sole direct parent/waiter of
Supervisor and sole writer of the fixed cancellation event. `scripts/run.sh`
validates a preapproved case ID and uses `exec` to **replace itself** with the
harness: it must not remain as a shell waiter or spawn a wrapper. Once its
direct Supervisor child is reaped or contained, the harness writes its
separate containment log and uses `exec` to replace itself with the compiled
read-only oracle. That oracle runs post-Supervisor even if fixture residue is
unknown, which is then a failure/manual-cleanup result. Neither script nor
oracle may signal, wait for or probe any lifecycle target. Gate-4 review
checks the actual `exec`/wait call graph. Any new linked source, generated
unit, runtime module,
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

`fixture.c` is one fixed, no-argument, no-network C executable. It arms
`alarm(4)` **before** its first blocking gate read. After reading exactly one
`G` byte from descriptor 3, it waits inertly for Supervisor stop; before that
it blocks at the gate. It never forks, execs, opens a path, reads an
environment value or writes output. The alarm provides *harness-only
containment*, never Supervisor teardown evidence. `posix_spawn` uses
`POSIX_SPAWN_CLOEXEC_DEFAULT`, `POSIX_SPAWN_SETSIGMASK` with an empty mask,
and `POSIX_SPAWN_SETSIGDEF` explicitly including `SIGALRM`; verify those flags,
default disposition and unblocked alarm before admission. This mirrors the
pinned C5b16 containment premise, not its synchronous teardown call graph.
Close all descriptors except 0/1/2 to `/dev/null` and read-only gate FD 3
in file actions; verify the inherited set in a dedicated negative fixture
test, and set an empty explicit environment. Any platform-
synthesized environment entry must be inventoried at Gate 4; an unexpected
entry fails admission. No pipe, lock, trace, directory or store FD enters the
child. The fixed executable path lies under a private 0700 root and is
read-only after hash readback. A pathname-to-inode swap or digest mismatch
fails admission; this is test-only identity, not installed code signing.

The record's 32-byte `ProcessIdentity` is SHA-256 over exact bytes
`"capsule.c5b19.process-identity/v1\0" || AttemptID[16] || pid_u32_be ||
creation_sequence_u64_be || creation_tick_u64_be || fixture_sha256[32]`.
PID must be positive and fit signed Darwin `pid_t`; sequence/tick come from
control's actual successful creation observation. This
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

The process tree is exactly harness → Supervisor → fixture. Harness is the
exclusive same-parent waiter for its **Supervisor** child; Supervisor is the
exclusive same-parent waiter for its **fixture** child. Harness never signals,
waits for or probes the fixture PID, even after Supervisor death. Harness may
bind its direct child in its own trace with SHA-256 over
`"capsule.c5b19.harness-child/v1\0" || AttemptID[16] || pid_u32_be ||
spawn_sequence_u64_be || spawn_tick_u64_be || supervisor_sha256[32]`.
That 32-byte digest is evidence, not PID-targeting authority. Harness may
make at most one `SIGKILL` request for its own live, unreaped Supervisor after
the 7-s watchdog, with its own `waitpid(WNOHANG)==0` and unchanged default
`SIGCHLD`/exclusive-waiter premise checked first. That action is harness-only
failure containment, not Supervisor stop or fixture absence. The fixture's
pre-gate 4-s alarm separately contains an orphaned fixture. Negative tests
before first run cover inherited ignored/blocked `SIGALRM`, alarm armed too
late, wrong FD set, watchdog identity mismatch and already-reaped Supervisor.
Any failed premise refuses the first run. A real alarm failure or unknown
fixture residue ends the campaign and requires manual owner cleanup; no
unsafe fallback PID action.

The harness creates an owner-only 0700 directory under a user-approved
disposable root, records its canonical path/device/inode, and refuses symlinks
or pre-existing occupants. It cleans only its exact owned directory after
authoritative child status, bounded worker shutdown and record snapshot.
After exact fixture absence and its last permitted operation, control closes
its request writer without waiting; worker finishes at most its current
operation, observes request EOF, then closes its reply writer. Control waits
at most 2 s for EOF on the nonblocking worker-reply pipe **after** consuming
any exact final reply; that
EOF is a shutdown observation, never a storage settlement. If no EOF, it
records failure, serializes its own trace and
terminates its **own** process without treating the pending write as settled.
No unbounded join is allowed. If it fails to exit, the harness's 7-s direct-
Supervisor watchdog applies as above. Never recursively delete a broad
parent. If worker remains blocked, child custody is unresolved, or any
harness-only containment fires, preserve the exact run directory and require
manual owner review/cleanup; never call a second PID-targeting fallback.
The child self-alarm at 4 s and harness watchdog at 7 s are independent
finite containment caps from C5b16, not evidence of the 1,000/1,200-ms
Supervisor bounds. The 9-s **verifier** timeout is a separate in-process
limit on compiled `tests/oracle.go` after harness has reaped its Supervisor
child and replaced itself. It records an incomplete verification on expiry
and has no process-targeting, signal, wait or cleanup
authority. It is not a third live-run watchdog or another parent of fixture.

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
cannot be treated as confirmed evidence. The obligation digest is SHA-256
over the **immutable** preimage
`"capsule.c5b19.obligation/v1\0" || record_version_u16_be(1) ||
Bindings[240] || ClockPolicy[24] || permitted_trigger_mask_u8(0x07)`.
Bindings use declaration order; the three policy u64s are big-endian
`1000,1000,1200`; mask bits 0/1/2 mean cancel/wall/fatal and all other bits
must be zero. The preimage excludes itself, process identity, operation,
generation, evolving candidate, outcome and record digest; thus its 32-byte
hash is identical across all three operations. Every nonzero/different-ID,
generation, candidate and digest
validation rule from C5b18 applies before a request is accepted.

The one harness cancellation event has a fixed **57-byte** frame: ASCII
`C5CN` (4), version 1 (u8), role 1 (u8), total length 57 (u16_be), consumed
AttemptID[16], ApprovalID[16], RegistrationID[16], and trigger kind 1 (u8).
Control accepts only exact equality with its frozen bindings from its
dedicated precreated harness pipe, at most once; duplicate, truncation, extra
byte, unknown role/kind or substitution is a fixture fault, not a second
action. The event contains no timestamp, target PID, path or content
authority. Its accepted ingress tick is sampled by control. Gate-3 TDD must
include fixed known-answer bytes for this frame, immutable obligation
preimage, process identity and durable record, plus mutation/refusal vectors
for every field and trailing/short bytes before live child tests.

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
directory, not C5b16 cursor/record/24-provider ABI. The **395-byte**
evolving stored candidate is ASCII `C5R1` (4), record version u16_be(1),
Bindings[240], obligation digest[32], ClockPolicy[24], generation u64_be,
operation ID[16], kind u8 and DurableProjection[68], in that exact order
(`4+2+240+32+24+8+16+1+68`). Its SHA-256 is the reply's record digest;
the digest itself is not stored inside those 395 bytes. No outcome is stored
in the candidate: confirmed means the worker observed successful file fsync,
rename and directory fsync for the exact candidate and returned a matching
reply. A present record after failed/unknown publication alone never proves
confirmation; recovery retains uncertainty. Terminal disposition is inside
the frozen projection, and no stored field grants public completion.
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

Raw evidence has **separate writers**. Control owns a preallocated
32,768-entry fixed-size in-memory array and sequence counter; worker owns
1,024 entries and harness owns 64. Each slot is exactly 256 bytes, so reserved
memory is 8 MiB control, 256 KiB worker and 16 KiB harness. Slot bytes are
explicitly encoded, not a native struct: magic `C5EV` (4), version u8(1),
lane u8 (control=1, worker=2, harness=3), kind u16_be, sequence u64_be,
tick u64_be, auxiliary tick u64_be, timebase numerator/denominator u32_be
each, binding digest[32] (SHA-256 of exact Bindings[240]), operation ID[16],
generation u64_be, request digest[32], process identity[32], PID u32_be,
syscall result i32_be, errno i32_be, flags u32_be, wait-status u32_be,
outcome u8, application result u8, durable-record digest[32], candidate
digest[32] (SHA-256 of exact DurableProjection[68]), and 10 zero reserved
bytes (`180+1+1+32+32+10=256`). Outcome codes are 0 none, 1 confirmed,
2 failed, 3 indeterminate; application codes are 0 not applicable,
1 accepted, 2 rejected. Unknown code, nonzero reserved bytes or nonsequential
lane index fails the oracle.

Kind u16 values are closed: 1 endpoint-created, 2 setup-check,
3 write-request, 4 request-write-result, 5 worker-reply,
6 worker-reply-write-result, 7 control-reply-received,
8 control-reply-applied, 9 spawn-call, 10 spawn-result, 11 identity-frozen,
12 cancel-ingress, 13 start-admission, 14 start-token-call,
15 start-token-result, 16 start-post-service, 17 wall-service,
18 fatal-service, 19 stop-latch, 20 signal-latch, 21 signal-call,
22 signal-result, 23 wait-call, 24 wait-result, 25 absence,
26 store-edge, 27 store-outcome, 28 trace-fault and
29 harness-containment. Kind 29 is harness-only. `flags` is zero except
store-edge's site code A=1/B=2/C=3/D=4; wait-status is populated only for
an exact terminal wait result. Fields not named for an event kind are zero.

`write-request` records the frozen operation ID, generation, request and
candidate digests before the one request write; `request-write-result` records
that syscall's returned byte count/errno. `worker-reply` records the same
tuple plus outcome and record digest before reply publication; its write-
result event records returned bytes/errno but has **no** control authority.
`control-reply-received` records the complete exact frame read by control;
`control-reply-applied` records the control validation decision:
code 1 only after exact tuple/request/candidate validation and the C5b18-
equivalent settle transition; code 2 for a rejected frame, which never
changes durable projection and latches fatal fault. Each event carries its
operation ID,
generation, request and candidate digests; confirmed settlement additionally
requires exact nonzero record digest matching the 395-byte stored candidate.
Control also copies each of its at-most-three exact 397-byte request frames
and at-most-three exact 98-byte complete reply frames to preallocated
read-only-after-capture slots, then serializes them beside its trace **after**
absence; worker separately retains its own received/sent frame slots before
any fault delay. These sidecars are evidence only, not transport authority.
The independent oracle decodes raw frames, recomputes request/candidate/
record hashes and expected C5b18 projections from frozen bindings and raw
process observations, then requires a received/applied pair for any claimed
confirmation. It compares those bytes to worker reply and durable snapshot.
A worker reply without control application never licenses start or terminal.

Field-use rules for the 256-byte slot (all unmentioned fields are zero):

| Kinds | Required non-common fields and refusal rule |
| --- | --- |
| 1–2 endpoint/setup | `tick`; setup check carries endpoint-creation tick in `auxiliary tick`. Missing or regressed observation fails. |
| 3–8 write transport | Exact operation ID, generation, request and candidate digests. Kinds 4/6 carry syscall result/errno; 5–8 carry outcome and record digest (zero digest unless confirmed); only kind 8 carries application result 1 or 2. Kinds 3/5/7 must pair with their captured raw frame slot. |
| 9–11 spawn/identity | Spawn result carries returned PID/result/errno; identity-frozen carries positive PID and exact `ProcessIdentity`. Failed spawn has no identity. |
| 12 cancellation | `tick` is control's accepted ingress, frozen binding digest is mandatory; no supplied timestamp or PID field. |
| 13–16 start | Kind 13 is `t_start`; 14 is actual write-call tick; 15 carries write result/errno and after-call tick; 16 carries final service tick and `t_start` as auxiliary tick. Missing 15 or 16 is not timely start evidence. |
| 17–19 wall/fatal/stop | 17 carries scheduled wall anchor in auxiliary tick and actual service in tick; 18 carries actual fatal-service tick; 19 carries earliest action anchor in auxiliary tick and actual stop service in tick. |
| 20–22 signal | Exact identity and positive PID; 20 is latch before call, 21 call, 22 return/errno. No second latch/call. |
| 23–25 wait/absence | Exact identity and positive PID; 23 call, 24 result/errno and terminal wait status when returned PID matches, 25 authoritative absence tick linked to that exact result. `ECHILD`, `ESRCH` or zero return cannot produce kind 25. |
| 26–27 store | Worker lane only; 26 carries one frozen A–D site in flags, 27 carries outcome and operation tuple. Neither may act as control reply application. |
| 28–29 fault/containment | 28 records lane-local trace fault; 29 is harness lane only and never a Supervisor stop or absence event. |

Every event includes its lane-local sequence, current-lane tick/timebase and
binding digest. Kinds 3–8 and 26–27 additionally bind the exact operation;
only control-lane kind 8 may mark a settled write in the oracle. Parser tests
must reject any nonzero forbidden field, missing paired event, changed raw
frame, outcome/digest substitution, and reversed same-tick stop/apply order.

`start-admission` records final `t_start` sample. `start-token-call` records
actual call entry; `start-token-result` records return/errno and tick
immediately after the nonblocking write; `start-post-service` records the
next continuous tick with `t_start` in auxiliary tick; the oracle derives
setup-expiry/timing classification from those raw ticks. A delayed or missing
post-service event fails the timing oracle even if prior admission was
timely. `stop-latch` stores
earliest action in auxiliary tick; wall/fatal service stores actual service
tick. `wait-call` and `wait-result` are paired with the same positive PID,
the latter carrying return/errno/status. The oracle rejects unpaired,
out-of-order, substituted or omitted causal events.

No lane reads or locks another's array. Control does no allocation, file I/O,
shared lock or worker call while it has child custody; each control event
append is a bounded copy to its own
next slot. Overflow latches case failure without overwriting earlier events.
The control lane's nondecreasing tick watermark is private; worker timestamps
are observations only and never advance/control that watermark. After exact
child absence, control serializes its array to its own evidence file without
waiting for worker. Worker separately flushes and fsyncs its own bounded
evidence file **before** entering each configured fault delay; its logging
failure fails the case but cannot block control. It may serialize its final
array after store completion; a stuck worker leaves only the pre-fault events
it already committed. Harness serializes its own bounded containment log
before replacing itself with the read-only oracle. No shared file append,
mutex, global sequence or cross-lane clock state.

Control records every actual `waitpid` call/result and each stop/start event,
not just a final status. After a zero-return wait observation, it may not
repeat `waitpid` until a 1-ms clock-paced `poll` iteration or an actual
cancel/reply event. A ready descriptor is consumed once; EOF/error closes
that control endpoint and removes it from the poll set, so readiness cannot
spin as an unbounded event source. `EINTR` immediately resamples as Gate 1
requires; an interrupt storm or trace overflow is explicit failed/unmeasured
evidence, never a pass. Under ordinary 1-ms pacing, the 7-s harness cap
allows fewer than 7,000 wait iterations. **Each iteration costs two slots**,
`wait-call` and `wait-result`: 6,999 iterations cost 13,998 slots and leave
18,770 of 32,768 control slots for all other events. The one-attempt
protocol permits at most three operations, one cancellation, one start and
one signal; its non-wait control event count is capped at 512 by pre-run
simulation including reply-application and post-start events. An interrupt
storm or repeated readiness beyond that cap fails as overflow/unmeasured
rather than creating passing evidence. This is a capacity design, not a
scheduler or latency guarantee. Pre-run simulated tests must exercise 6,999
paired waits plus 512 other events, full-budget and overflow refusal,
immediate-ready/EOF removal and `EINTR` storms without a child. The live
oracle distinguishes real stop/timing failure from trace exhaustion.
Aggregate cap is 33,856 entries; absent, truncated or overflowed lane
evidence is failure, not a pass.
The oracle merges by exact operation/request identity and causal control
events, **not** by assuming a total timestamp order between lanes. A same-
tick stop/reply race is decided by the control lane's own stop-latch and
reply-application sequence; worker timestamps cannot reorder it. Harness
actions never enter the Supervisor trace or become Supervisor absence.
`tests/oracle.go` computes transitions and bounds from these streams plus
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
