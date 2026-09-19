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
Remaining work: mechanism/source review, exact runnable-plan review, explicit
owner authorization, implementation, bounded run, retained evidence and review.
Next action: review this plan and close the pre-implementation decisions below.
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
implementation. No external dependency or new primitive is selected here.

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
| Documented fact | Apple's [`sigaction(2)`](https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/sigaction.2.html) says `SA_NOCLDWAIT` suppresses zombies. [POSIX process-ID reuse](https://pubs.opengroup.org/onlinepubs/009696699/basedefs/xbd_chap04.html#tag_04_12) prohibits reuse until process lifetime ends; [POSIX exit](https://pubs.opengroup.org/onlinepubs/009695299/functions/exit.html) retains an unwaited child as a zombie absent `SIGCHLD` ignore/`SA_NOCLDWAIT`. These support only a conditional same-parent, exclusive-waiter direct-child argument. |
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

For the first clock candidate, use [`mach_continuous_time()`](https://developer.apple.com/documentation/driverkit/mach_continuous_time)
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

Apple's [uptime-raw clock note](https://developer.apple.com/documentation/driverkit/kiotimerclockuptimeraw)
distinguishes a monotonic clock that excludes sleep from a continuous clock
that includes it. This is a clock-policy decision, not evidence of timely stop
through suspend/wake. The runnable plan must specify the chosen Darwin clock's
behavior across suspend/wake and identify scheduler stalls as timing failures
or unmeasured conditions; it may not silently pause or reset a budget. Restart
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

1. **Mechanism review — `BLOCKED` on exact call graph and clock choice.**
   Specify current-host process/clock semantics, single-waiter custody, bounded
   worker and lock separation, exact frozen operation settlement, and the
   setup/stop race. Reject or revise candidate before implementation if any
   premise fails.
   Check the native test/platform and narrow Supervisor-transport rows of
   [ecosystem reuse map](ECOSYSTEM_REUSE_AND_ADOPTION.md); complete its policy
   checklist for any proposed dependency or custom primitive. No new package is
   assumed.
2. **Exact runnable-plan review — `BLOCKED` on gate 1.** Freeze source/toolchain/
   OS/architecture, artifact hashes, directory and process caps, fault sites,
   timeouts, cleanup procedure, descriptor allowlist, executable identity,
   independent oracle and archive destination.
   Obtain independent security/engineering review of that exact packet. Close
   every blocking finding before requesting execution authority.
3. **Owner authorization — `BLOCKED` on gate 2.** Owner explicitly authorizes
   implementing and running the frozen benign packet on a named owned test host.
   Approval of ADR-0047 or merge of this document is not that authorization.
4. **Experiment and evidence — `BLOCKED` on gate 3.** Implement in the experiment
   archive, run the bounded ordinary/sanitizer/fault/mutation corpus, reproduce
   from a clean second directory, review raw evidence and publish an immutable
   archive commit. A failed case produces an adverse result; do not widen clocks
   or change the fixture after seeing it without a new reviewed packet.
5. **Canonical conclusion — `BLOCKED` on gate 4.** Record exact pass/fail and
   limitations here with commit-pinned evidence. The experiment passes only if
   ordering, storage-independent stop, exact same-lifetime custody, both timing
   bounds for every applicable case, fail-closed uncertainty and mutation
   sensitivity all hold. A passing benign single-child result still does not
   admit installed lifecycle, descendants, restart custody, a guest or product.

Review checkpoints after gates 1, 2 and 4 prevent a failed premise from
becoming runnable authority. Existing repo checks may validate this document;
they do not execute or authorize the fixture.
