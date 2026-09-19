# C5b19 benign teardown experiment plan

Date: 2026-09-19

Work item: documentation-only plan for a later bounded runnable experiment.
Status: `IN_PROGRESS — TRENDING_GOOD` while exact mechanism and review are open.
Scope: defensive validation of ADR-0047's prepared-obligation ordering and
storage-independent Supervisor stop path, using one pinned fixed benign direct
child in owned disposable local directories on an explicitly authorized macOS
test host. This document authorizes no implementation or run.
Evidence or reason: [C5b16](C5B_TIMING_FAULT_CHECKPOINT.md) found a controlled
publication delay that exceeded the total teardown bound; [C5b18](C5B18_PASSIVE_SUCCESSOR_SPECIFICATION.md)
passed a no-effect ordering model; [ADR-0047](adr/0047-prepare-teardown-obligation-before-launch.md)
is Accepted as architecture only.
Remaining work: mechanism/source review, exact implementation-packet review,
explicit implementation authorization, independent implemented-byte review,
separate first-run authorization, bounded run, retained evidence and review.
Next action: independently review the gate-1 prospective mechanism packet
below; this plan does not approve implementation or a first child run.
Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`.
Runnable experiment, installed lifecycle, guest execution and product admission:
`BLOCKED` pending the named review, authorization and separate later evidence.

## Decision question and boundary

Can one current-lifetime Supervisor control path request one exact child's stop
and observe its physical absence within the existing C2A bounds while a single
post-create durable publication is blocked, refused or completed late, without
releasing start or completion improperly? The experiment tests *mechanism*, not
an installed service, hostile child, backend, guest, process tree or product.
Neither a fixture self-alarm nor an independent watchdog counts as Supervisor
teardown. No daemon, Broker, private key, user content, network, VM, guest,
interpreter or live host path enters this experiment.

Completed one-time source, harness and raw evidence belong in
`Shrimpworks/capsule-experiments`; the canonical outcome and exact archive commit
pin return here. Existing C5b11/C5b14B/C5b16 driver, cursor 17, record format,
24-provider ABI and product code remain unchanged. A successor gets new format,
profile and driver identities, never a reinterpretation of old evidence.

## Proposed mechanism to review, not yet selected for implementation

1. **Preparation:** one closed immutable obligation, bound to the exact consumed
   attempt, approval, registration, installation/epoch, Supervisor, plan/profile/
   runner and Supervisor-generated operation ID, must be durably confirmed before
   child creation. Failed, pending or indeterminate preparation permits no spawn.
   The record contains no fabricated PID or signal/absence claim.
2. **Single owner:** one native Supervisor control lane owns creation, exact
   current-lifetime child custody, cancellation/deadline ingress, stop latch,
   identity check, one signal attempt, reap and absence observation. A separate
   bounded storage lane owns all store calls and their locks. The control lane
   never acquires the storage lock, waits for its worker, joins drain/completion,
   or invokes a synchronous bridge operation before stop/absence. One storage
   operation and one attempt may be outstanding; the owner lock stays held while
   its result is pending. The worker posts one exact result to a bounded mailbox;
   only the control lane applies that result or authorizes start. A full mailbox
   cannot block stop or grow into a queue. No helper or second lifecycle owner
   is introduced.
3. **Direct-child custody candidate:** the fixed child cannot fork or exec and
   has no descendants. Candidate native custody uses the Supervisor's own
   successful child-creation observation and exclusive same-lifetime wait/reap
   ownership; no PID lookup, recovery adoption or second waiter is allowed. A
   reviewer must verify the exact Darwin spawn, signal, child-status and PID
   lifetime semantics from primary platform sources, including the exit-before-
   signal race, before this candidate may be implemented. An unproved or
   ambiguous identity blocks signaling and the run; numeric PID alone is not
   authority. This does not establish process-tree or restart custody.
4. **Start gate:** publish the observed child identity durably before the first
   start-token write. Storage refusal leaves the child gated and start closed.
   The control lane must still service authenticated attempt-bound cancellation
   or Supervisor-detected setup expiry classified as a fatal fault, with exact
   custody. A late success cannot release start after stop. The same control
   lane serializes stop-latch checks and start authorization. Capture `t_start`
   immediately before the first start-token attempt, only after confirmed
   identity publication and no stop latch.
5. **Stop and publication:** a monotonic stop latch holds the earliest action
   anchor and distinct actual service tick. Exact custody permits one forced
   signal request and authoritative same-child reap/absence observation without
   waiting for any pending write. A lost signal response is never redriven.
   A late storage response settles only its exact frozen operation/generation/
   candidate; it cannot clear stop, start a child, overwrite newer evidence,
   release output/capacity or erase uncertainty. Terminal publication follows
   observed absence and remains short of public completion until the full typed
   result/lifecycle/cleanup/durable join exists. Storage failure with exact live
   custody does not itself revoke prepared destructive cleanup; lost custody does.

This is a proposed division of work, not an approved new Supervisor
responsibility. ADR-0047 already selects the Supervisor-owned scheduling split;
if mechanism review finds a helper, additional lifecycle owner, different
authority or durable ordering necessary, stop and write a separate ADR before
implementation. No external package is selected. The platform-clock and
single-slot mailbox candidates still require gate-1 mechanism and custom-
primitive policy review; neither is product-admitted by this plan.

Primary-source starting points (read 2026-09-19): Apple's archived
[`posix_spawn(2)`](https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/posix_spawn.2.html)
defines success/failure child-PID behavior and warns that unclosed descriptors
are inherited; [`kill(2)`](https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/kill.2.html)
targets a positive PID, while [`wait(2)`](https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/wait.2.html)
specifies exact-child `waitpid` and `WNOHANG` observations. These pages do not
by themselves prove no PID reuse across an exit/signal race or a current-macOS
single-waiter implementation. Review current SDK/man-page behavior and the
actual native source before approving custody. Freeze the child's descriptor
allowlist and close-on-exec behavior; a start pipe must not carry other host
authority into the fixture.

## Mechanism source review: evidence, inference and open gate

Read-only review on 2026-09-19 used macOS 26.6.2 (build 25G83), Xcode 26.6,
macOS SDK 26.5, and the immutable C5b16 archive commit
`0efd03def6bc333a92c8b9809bc56b7b3ce9ea80`. No child or fault fixture was
run. The possible later run host has not been named or authorized.

| Class | Evidence and consequence |
| --- | --- |
| Documented fact | Apple's [`posix_spawn(2)`](https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/posix_spawn.2.html) returns a child PID on success, leaves its output undefined on failure, and can inherit descriptors. Apple's [`wait(2)`](https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/wait.2.html) supports positive-PID child-specific `waitpid(..., WNOHANG)` and makes `ECHILD` an error, not absence proof. [`kill(2)`](https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/kill.2.html) addresses a positive PID, but its success or `ESRCH` is not a reap observation. |
| Documented fact | Apple's [`sigaction(2)`](https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/sigaction.2.html) says `SA_NOCLDWAIT` suppresses zombies. [POSIX.1-2024 process-ID reuse](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html) prohibits reuse until process lifetime ends; [POSIX.1-2024 process termination](https://pubs.opengroup.org/onlinepubs/9799919799/functions/_exit.html) retains unwaited child status absent `SIGCHLD` ignore/`SA_NOCLDWAIT`. These support only a conditional same-parent, exclusive-waiter direct-child argument. |
| Documented fact | Go's [`os/signal` documentation](https://go.dev/src/os/signal/doc.go) says `-buildmode=c-archive` does not install asynchronous signal handlers by default, but `signal.Notify` can install one. A future archive's imports, initialization and linked native code still need source closure; the build mode alone does not prove default `SIGCHLD` or exclusive reaping. |
| Source observation | Pinned [C5b16 native lifecycle](https://github.com/Shrimpworks/capsule-experiments/blob/0efd03def6bc333a92c8b9809bc56b7b3ce9ea80/experiments/typed-guest-transport-c5b16-timing-fault/source/native/lifecycle.c) checks default `SIGCHLD` without `SA_NOCLDWAIT`, uses `POSIX_SPAWN_CLOEXEC_DEFAULT`, retains a successful returned PID, and uses child-specific `waitpid(WNOHANG)`. It assumes one serialized caller and exclusive reaper. Its effect 16 still waits for `BeforeTeardown` before signaling, and effects 17–19 checkpoint before reap/absence; it is not the successor mechanism. Pinned [Go bridge](https://github.com/Shrimpworks/capsule-experiments/blob/0efd03def6bc333a92c8b9809bc56b7b3ce9ea80/experiments/typed-guest-transport-c5b16-timing-fault/source/bridge/owner.go) holds one mutex across store calls. |
| Inference | If one live Supervisor lane is the exclusive waiter, `SIGCHLD` remains at default without `SA_NOCLDWAIT`, and no other code reaps the child, an exit between `waitpid(WNOHANG)==0` and positive-PID `kill` should not redirect the signal to a reused PID: an exited child remains unreaped. Still latch the signal attempt before calling `kill`; regardless of return, use exact `waitpid` status for absence. This is not a proved property of a future concurrent implementation. |
| Source conflict / unknown | Apple's [continuous-time page](https://developer.apple.com/documentation/driverkit/mach_continuous_time) says `mach_continuous_time()` advances during sleep and suggests `CLOCK_MONOTONIC_RAW` as equivalent, while its [uptime-raw page](https://developer.apple.com/documentation/driverkit/kiotimerclockuptimeraw) says `CLOCK_MONOTONIC_RAW` excludes sleep. SDK 26.5 `mach/mach_time.h` explicitly describes `mach_continuous_time()` as advancing during sleep; C5b16's [trace](https://github.com/Shrimpworks/capsule-experiments/blob/0efd03def6bc333a92c8b9809bc56b7b3ce9ea80/experiments/typed-guest-transport-c5b16-timing-fault/source/native/timing.h) uses `CLOCK_MONOTONIC` for observation only. Do not select a clock by treating the contradictory aliases as settled. |

Mechanism disposition: `IN_PROGRESS — TRENDING_GOOD` for source narrowing,
not `PASSED` for runnable readiness. Before gate 1 closes, freeze a reviewable
native/Go lane and mailbox call graph, check the existing source closure for
competing waiters or signal-disposition changes, choose one clock and coherent
awake/suspend policy, and specify exit-before-signal, `ECHILD`/`ESRCH` and
blocked-store test oracles. After explicit implementation authorization, inspect
the resulting code and prove those assumptions before its first child run; a
mismatch stops the experiment. The fallback is no signal or start on lost
custody, and no timing/completion claim on missing observations; it is not a
daemon/helper watchdog. If direct-child exclusivity fails, revise the candidate
and seek an ADR for any added lifecycle authority.

The 2026-09-19 read-only source check confirmed the pinned C5b16
`reaper_owned()` predicate and child-specific `waitpid` in `lifecycle.c`, and
`bridgeMu` spanning `BridgeApply`'s store call in `owner.go`. It also confirmed
that the current macOS 26.5 SDK declares `mach_continuous_time()` available since
macOS 10.12 and advancing during sleep. These are observations of old source
and local headers, not a closed successor call graph or a timing measurement.
The design chooses continuous elapsed time for the first successor packet;
the named host, checked conversion, actual call graph and implementation must
still prove that policy. Do not mix an uptime clock into its anchors.

### Proposed first call graph for gate-1 review

This is a design candidate, not executable source or a timing guarantee:

1. Before creation, the storage lane confirms one exact prepared obligation.
   Failure or unknown commit closes creation. Only a confirmed result reaches
   the native control lane.
2. The control lane creates one child, retains its successful PID and sole
   wait/reap right, then submits one frozen runner-identity publication to a
   worker. The worker owns the bridge/store lock and writes one immutable reply
   into a single-slot release/acquire mailbox; it never changes custody, stop,
   clocks, start or capacity. No second write or attempt bypasses that slot.
3. On each control iteration, sample the chosen clock and accept the fixed
   attempt-bound test cancellation slot; evaluate fatal/setup expiry and the
   wall bound **before** applying a ready store reply. Stop wins a same-tick
   race. Only this lane may apply an exact settled reply or write a start token.
   A pending mailbox cannot block the iteration; worker completion cannot
   directly wake or authorize the child.
4. Once stop latches, check current-lifetime reaper ownership and observe the
   exact child via nonblocking `waitpid`. If still live, set the private signal-
   attempt latch **before** one positive-PID `kill`; then continue exact
   `waitpid` observation. `kill` return and `ECHILD`/`ESRCH` never substitute
   for absence. Keep polling/dispatch bounded independently of the worker;
   record actual service/signal/reap times, not a scheduled callback time.
5. After exact absence, defer terminal publication until the prior operation
   has settled and the exact frozen generation can advance. A stuck worker
   retains owner lock and blocks durable completion/capacity indefinitely, not
   stop or absence. No result/output is released from this experiment.

### Gate-1 read-only ownership and proof packet

Status: `IN_PROGRESS — TRENDING_GOOD` for a reviewable candidate, not a gate-1
pass or runnable approval. The maintainer owns the mechanism decision. This
packet compares the selected ADR-0047 scheduling split with the old synchronous
chain; it inspects source and platform contracts only. No child or fault case
was executed.

| Step | Sole owner and data crossing | Forbidden dependency or result |
| --- | --- | --- |
| Prepare before spawn | Storage lane confirms one frozen obligation; control receives only its exact confirmed result. | Pending/refused/indeterminate preparation never creates a child. Waiting here precedes custody, not a post-create stop. |
| Create and publish | Control alone records successful current-lifetime custody and consumes creation; it submits one frozen identity-publication request carrying obligation, operation ID, generation, kind and candidate. | Worker receives no PID-targeting, start, stop or capacity authority. Failed/ambiguous spawn never supplies a signal target. |
| Store settlement | Worker alone holds bridge/store locks and returns one immutable reply with the same tuple and outcome through a bounded slot. Control validates every field before applying it. | A full slot may stall the worker, not control. No second operation, attempt or generation bypasses it; a late reply cannot release start after stop. |
| Start or stop | Control samples one continuous clock, records actual ingress/service, evaluates fatal/setup/wall stop first, then applies an exact ready reply and permits at most one start-token attempt. | Worker cannot wake or authorize the child. Control never calls `BridgeApply`, a store method, a blocking mailbox receive, a drain/completion join or a store lock while custody exists. |
| Signal and absence | Control alone checks the same-parent exclusive-reaper premise and uses positive-PID `waitpid(..., WNOHANG)`. A matching terminal status proves absence. If still live, it latches the one signal attempt before `kill`, then continues exact child-status observation. | `kill` success, `ESRCH`, EOF, child alarm and root removal prove no absence. `ECHILD`, changed disposition, competing reaper or mismatched identity is custody failure, never permission to signal again. |

The worker request/reply slot is a *candidate custom primitive*, not an adopted
implementation. Gate 1 must choose and review its concrete C/Go ownership,
release/acquire or kernel-transport semantics, fixed capacity, cancellation
ingress and wake/poll behavior against the
[ecosystem checklist](ECOSYSTEM_REUSE_AND_ADOPTION.md). A second thread is not
a second lifecycle owner. An unbounded queue, a worker-controlled start, or
blocking on worker completion in the control path rejects this candidate.

The direct-child exit-before-signal argument is conditional: a successful
`posix_spawn` supplies a positive PID; the same parent retains it; `SIGCHLD`
stays at default without `SA_NOCLDWAIT`; no other code waits or reaps; and
`waitpid` has not reaped the child when it returns zero. Under
[POSIX.1-2024's process lifetime and PID-reuse rules](https://pubs.opengroup.org/onlinepubs/9799919799/basedefs/V1_chap04.html),
an exit before `kill` then leaves status for that parent and should not redirect
the signal to a reused PID. This is an inference, not proof that future linked
Go/native code or its host satisfies those premises. A matching `waitpid`
return may instead observe natural exit and must suppress the signal. After a
zero return, the signal attempt stays latched even if `kill` returns an error;
only the same-child status can later establish absence. If the premise or
observation fails, stop the run with unresolved custody and harness-only
containment, not a PID probe or second signal.

Before gate 1's design decision can pass, an independent reviewer must close
these prospective proofs; gate 4 separately checks the implemented bytes:

1. Trace every prospective native/Go entry point from cancellation, deadline,
   setup expiry and identity-publication reply to start, signal and reap. Show
   that no post-create path acquires a bridge/store lock or waits for the worker.
2. Inventory the pinned native/Go source for `wait`, `waitpid`, `waitid`,
   process-wide `SIGCHLD` changes, `signal.Notify`, auto-reaping and competing
   child owners. Specify the same audit over the later complete linked source
   as a gate-4 condition. The pinned C5b16 `reaper_owned()` check is a useful
   guard, not a substitute for either closure.
3. Specify the single-slot primitive and its producer/consumer ordering,
   frozen tuple equality, no lost stop ingress, bounded control polling and
   exact reply behavior when the worker is blocked, refused or returns late.
   Require code-level proof of these properties at gate 4.
4. Freeze the `mach_continuous_time()` conversion rule and nondecreasing tick
   policy for a proposed supported macOS floor; gate 2 names the exact host.
   Late wall service, suspend/wake and scheduler stalls must fail timing rather
   than replenish a budget.
5. Define independent raw-event oracles for blocked publication, child exit
   before signal, `ECHILD`/`ESRCH`, lost signal response, stop/start same-tick
   race and late reply. Each missing authoritative event remains null/failure;
   C5b16's self-alarm never counts as Supervisor absence.

The pinned [C5b16 driver](https://github.com/Shrimpworks/capsule-experiments/blob/0efd03def6bc333a92c8b9809bc56b7b3ce9ea80/experiments/typed-guest-transport-c5b16-timing-fault/inputs/supervisor_effect_driver.c)
still fences and looks up the attempt before effect 16; pinned
[`lifecycle.c`](https://github.com/Shrimpworks/capsule-experiments/blob/0efd03def6bc333a92c8b9809bc56b7b3ce9ea80/experiments/typed-guest-transport-c5b16-timing-fault/source/native/lifecycle.c)
still calls `BeforeTeardown` before signal and checkpoints before reap; pinned
[`owner.go`](https://github.com/Shrimpworks/capsule-experiments/blob/0efd03def6bc333a92c8b9809bc56b7b3ce9ea80/experiments/typed-guest-transport-c5b16-timing-fault/source/bridge/owner.go)
holds `bridgeMu` across the store call. These observations reject reusing that
call chain unchanged. They do not demonstrate the successor's independence.

For the first clock policy, use [`mach_continuous_time()`](https://developer.apple.com/documentation/driverkit/mach_continuous_time)
plus checked [`mach_timebase_info()`](https://developer.apple.com/documentation/driverkit/mach_timebase_info-c.func)
conversion for every same-lifetime anchor and observation: the SDK header and
Apple's direct API description say it advances through sleep.
On wake, a serviced wall callback after its original anchor is a timing failure;
sleep never grants a new budget. The later exact packet must confirm API support
on its named macOS floor and test the conversion/overflow and callback-service
oracles. The contradictory `CLOCK_MONOTONIC_RAW` documentation above is not used
to infer an alias or mix two clocks. The test harness supplies only fixed,
attempt-bound cancellation fixtures; it does not prove product authentication.

## Clock and observation contract

Use the [C5b17 anchors](C5B_TEARDOWN_DEADLINE_DESIGN.md#proposed-clock-anchors)
and C5b18 same-lifetime nondecreasing observed-tick watermark. Monotonic
`t_wall = t_start + 1,000 ms`; `t_action` is the earliest valid wall,
Supervisor-ingress cancellation or detected fatal-fault anchor. Wall callback
service has a separate observed tick; late service fails timing even if the
signal is prompt. First successor has no grace. Require both authoritative
`t_absent <= t_force + 1,000 ms` and
`t_absent <= t_action + 1,200 ms`; neither budget restarts after storage
settlement. Record dispatch, identity-check, signal-call and reap/absence ticks
separately. Signal return, EOF, child alarm, root removal and harness reap are
not `t_absent`.

For a child held before the first start-token attempt, `t_start` and `t_wall`
do not exist. A blocked/refused identity publication must therefore remain a
start refusal, not a synthetic wall-time pass. Accepted cancellation or detected
fatal setup expiry supplies `t_action`; if exact custody permits a forced stop,
measure both force-to-absence and action-to-absence against the same bounds.
Missing action, force or authoritative absence is null/failure, never a passing
latency sample. The independent oracle must distinguish this pre-start case
from a started child whose wall callback was serviced late.

Apple's [uptime-raw clock note](https://developer.apple.com/documentation/driverkit/kiotimerclockuptimeraw)
distinguishes a monotonic clock that excludes sleep from a continuous clock
that includes it. An uptime clock could silently discount sleep and classify a
late wake as on time; continuous elapsed time avoids that accounting error but
does not force timely scheduling. This is a clock-policy decision, not evidence
of timely stop through suspend/wake. The runnable plan must verify the chosen
Darwin clock's behavior on its named host and identify scheduler stalls as
timing failures or unmeasured conditions; it may not silently pause or reset a
budget. Restart
clears live custody/ticks, preserves preparation and uncertainty, blocks
replacement and reports recovery required. Supervisor death, suspend/power loss
and installed recovery remain separate gates; this single-child run makes no
success claim for them.

## Bounded cases and independent oracles

Use fixed, hash-pinned child and harness bytes; a disposable owner-created root;
no externally supplied executable, path, delay, backend flag or process target.
Fault injection is armed once before the run at a named store edge, with fixed
finite delays no greater than the C5b16 1,600-ms comparison, or fixed refusal/
lost-response behavior. It is not an execute-time option. The child is inert,
has no descendants, holds a start gate, and has a separate finite harness-only
containment alarm. Record if that alarm fires; never count it as Supervisor
absence. The owner must approve exact child source/hash, maximum runtime,
directory, expected child count and cleanup before any run.

| Case | Required observation and refusal |
| --- | --- |
| Confirmed preparation; zero-delay normal | One child, one start, one stop, exact absence, durable facts in order; no public-completion claim from this slice. |
| Preparation pending/refused/indeterminate | Zero child creation/start; no guessed confirmation or retry identity. |
| Created child; identity publication blocked at each bounded write edge | Start stays shut. Accepted cancel or setup expiry still reaches exact stop/absence without store settlement. Record actual action and absence ticks. |
| Post-create write refused or response lost | Exact live custody permits prepared stop; failure/uncertainty remains durable and no result/capacity is released. |
| Late write success after stop or absence | Exact operation may settle but cannot reopen start, stop, authority or capacity. |
| Cancel/wall/fatal races and duplicate requests | Earliest anchor and watermark retained; no second signal or deadline extension. |
| Child exits before signal; signal response uncertain; identity mismatch | No signal to a possibly reused PID, no redrive, no invented absence. Exact native child-status observation decides what can be claimed. |
| Terminal write blocked/refused after native absence | Physical absence fact remains distinct; no public completion/output/capacity release. |
| Supervisor process termination/restart probe | Only if separately reviewed as safe: obligation/uncertainty persist, live custody is not restored, no PID adoption/signal/replacement or timely-absence claim. Harness cleanup is separately labeled. |

Run ordinary and sanitizer builds where supported. Test the known C5b16
0/800/1,200/1,600-ms comparison delays; do not infer a latency guarantee from
finite sleep. For every case retain raw monotonic trace, native process-status
observations, store operation/generation/candidate and settled result, exact
artifact digests, exit/refusal classification and any harness-only cleanup.
Independent checker derives allowed transitions and timing from raw observations,
not the fixture's success flag. Missing event is null, never a pass. A timing
miss or custody ambiguity is an explicit failed/unresolved case, not silently
reclassified as normal teardown. Stop the campaign on unexpected additional
processes, a child outside the owned fixture, lost containment, unbounded
storage-worker growth or inability to clean up safely.

Mutation campaign must execute compiled/source variants that restore a stop-to-
storage wait, release start after stop, repeat a signal, replace earliest action
with a later trigger, forge absence from signal/EOF/alarm, accept a substituted
identity, settle a changed operation/candidate, or release completion before
the terminal join. Each relevant mutation must be detected by an independent
oracle; guard tests alone are not mutation proof. Preserve the failing mutant
trace and checker output, then restore the exact clean source.

## Ordered gates and acceptance

1. **Mechanism review — `BLOCKED` on independent prospective call-graph,
   mailbox and custody review.**
   Specify current-host process/clock semantics, single-waiter custody, bounded
   worker and lock separation, exact frozen operation settlement, and the
   setup/stop race. The continuous-time policy is chosen for the candidate,
   not yet validated on a named host or implemented source. Reject or revise
   the candidate before implementation if any premise fails.
   Check the native test/platform and narrow Supervisor-transport rows of
   [ecosystem reuse map](ECOSYSTEM_REUSE_AND_ADOPTION.md); complete its policy
   checklist for any proposed dependency or custom primitive. No new package is
   assumed.
2. **Exact implementation-packet review — `BLOCKED` on gate 1.** Freeze the
   intended source inputs and call graph, toolchain/OS/architecture, fixture and
   harness specifications, directory and process caps, fault sites, timeouts,
   cleanup, descriptor allowlist, executable-identity policy, independent
   oracle and archive destination. Obtain independent security/engineering
   review of that exact packet. Close every blocking finding before requesting
   implementation authority.
   Future implementation and binary digests cannot be asserted at this gate.
3. **Owner implementation authorization — `BLOCKED` on gate 2.** Owner explicitly
   authorizes only implementing the frozen benign packet in the owned experiment
   archive. Approval of ADR-0047, merge of this document, or this authorization
   does not authorize a child run.
4. **Implemented-byte review and first-run authorization — `BLOCKED` on gate 3.**
   Implement without running the child. Freeze the resulting native/Go source,
   fixture and harness source, generated bytes, build inputs and executable
   digests; independently review the actual control/storage/mailbox call graph,
   waiter and signal-disposition closure, lock and clock paths, descriptor
   allowlist, test oracle and cleanup
   against gates 1–2. Read back the exact source and binary digests selected for
   the run. A material mismatch reopens packet review and owner implementation
   authorization. Only after blocking findings close may the owner separately
   authorize the first run on the named owned host; the run must refuse if its
   digests or environment differ from the approved readback.
5. **Experiment and evidence — `BLOCKED` on gate 4.** Run the bounded ordinary/
   sanitizer/fault/mutation corpus, reproduce from a clean second directory,
   review raw evidence and publish an immutable archive commit. A failed case
   produces an adverse result; do not widen clocks or change the fixture after
   seeing it without a new reviewed packet.
6. **Canonical conclusion — `BLOCKED` on gate 5.** Record exact pass/fail and
   limitations here with commit-pinned evidence. The experiment passes only if
   ordering, storage-independent stop, exact same-lifetime custody, both timing
   bounds for every applicable case, fail-closed uncertainty and mutation
   sensitivity all hold. A passing benign single-child result still does not
   admit installed lifecycle, descendants, restart custody, a guest or product.

Review checkpoints after gates 1, 2, 4 and 5 prevent a failed premise from
becoming runnable authority. Existing repo checks may validate this document;
they do not execute or authorize the fixture.
