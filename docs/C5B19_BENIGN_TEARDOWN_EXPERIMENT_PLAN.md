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
clears live custody/ticks,
preserves preparation and uncertainty, blocks replacement and reports recovery
required. Supervisor death, suspend/power loss and installed recovery remain
separate gates; this single-child run makes no success claim for them.

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

1. **Mechanism review — `BLOCKED` on primary-source and local source evidence.**
   Confirm Darwin process/clock semantics, single-waiter custody, bounded worker
   and lock separation, exact frozen operation settlement, and the setup/stop
   race. Reject or revise candidate before implementation if any premise fails.
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
