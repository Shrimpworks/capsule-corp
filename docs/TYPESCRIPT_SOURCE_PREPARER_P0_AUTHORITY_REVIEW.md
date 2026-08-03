# TypeScript Source Preparer P0 authority and TCB review

Date: 2026-08-03

Status: retained adversarial design review; **P1 HOLD / conditional topology decision**. No Source
Preparer service, store, worker contract, endpoint, key, runtime, backend, guest, or deployment is
implemented or admitted by this review.

Decision owner: Capsule architecture and security review.

## Defensive scope and method

This review used only the Capsule repository at commit
`c68dfb1535b6763ad7c89d5f401fa9002f225b26`, its retained passive fixtures, source, plans, and local
evidence. It accessed no unrelated system, identity, credential, user data, runtime, backend,
guest, deployment, signing service, or third-party target. It created no endpoint and executed no
untrusted workload.

The review traced the proposed Source Preparer against the project, architecture, technical design,
threat model, feasibility gates, installation/update model, ADR-0010/0011/0013/0018/0026/0029/
0030/0031/0032, the authenticated-IPC stop, source cutover and archive plans, the field-authority
manifest/verifier, passive approved-byte fixtures, and the current proposal, plan, registration,
approval, and fake-lifecycle code. Three independent issues-only reviews separately challenged
authority/topology, same-user/store/worker confinement, and lifecycle/update/archive behavior.

No live platform probe was needed or authorized. The exact Source Preparer package does not exist,
and the repository already retains the decisive negative fact that ordinary same-user mode bits
permit a known-path writer or retained writable mapping to mutate a custody object. Existing Gate B
evidence establishes that macOS protected-container mechanisms are plausible, but it does not
validate the proposed Source Preparer or child Node composition.

## Decision

The separate, unprivileged Source Preparer remains the least dangerous selected topology **if
Capsule ships TypeScript preparation**. It keeps Node/Amaro and the source store out of the daemon,
Broker, Supervisor, updater, runtime, and backend; adds no operational key or guest route; and
preserves pre-registration transformation and execute-by-registration-ID.

The topology is not ready for P1 contract or byte freeze as written. P1 is a bounded **NO-GO now**
and may become a conditional GO only after the pre-entry criteria in this review are retained in
the ADR, plan, field-authority model, and canonical security documents.

Two propositions must not be conflated:

1. one role-namespaced store is the smallest coherent ownership topology; and
2. one role-namespaced store is a security boundary.

The first is accepted conditionally. The second is rejected. Role namespaces prevent honest
cross-role lookup and deduplication confusion; they do not constrain a compromised store owner and
do not stop an ordinary same-user writer. The security boundary requires an exact OS-protected,
single-member storage mechanism plus installed adversarial evidence.

The least-dangerous viable alternative is JavaScript-only pass-through with TypeScript unsupported.
It removes the parser, transformation claim, worker, source-set saga, and semantic-substitution
trust. User planning direction accepts that bounded first-release fallback only as modern ESM
`.mjs`: no CommonJS, package resolution, legacy Node module surface, or widening of the governed
runtime contract. This is not an already frozen media/profile decision; the exact contract and
applicable ADRs must be updated before implementation. Moving Node into the daemon, Broker,
Supervisor, updater, runtime, or backend is
rejected. A second independent transformer could change the trust claim, but adds an authority,
protocol, store/update seam, and must be decided in a new ADR.

## Accepted trust consequence

The native front, Go core, immutable-store implementation, and Node/Amaro worker are one planning
and approval-understanding TCB for the selected exact-byte claim. Native and Go share an address
space; their separation is not a security boundary. Unless a separately evidenced child profile
proves less authority, Node compromise has the Source Preparer's complete blast radius.

A compromised Source Preparer can:

- read proposal source handled or retained by the component;
- choose arbitrary emitted JavaScript and create internally consistent original, executable,
  profile, options, record, and record-set claims before plan construction;
- retain private copies despite a later logical release;
- replay, lose, corrupt, quarantine, or refuse source sets and exhaust bounded service capacity;
- manipulate preparation, registration-pending, retention, and release behavior if its state logic
  is compromised; and
- lie about having run Node/Amaro while recomputing every digest consistently.

Parent-side hashing prevents trusting worker-supplied digests, counts, media, or records. It does
not detect arbitrary but internally consistent JavaScript selected by a compromised parent or
worker. Supervisor and Broker validators prove canonical bytes and exact digest/media/path/record
relationships; they do not prove that Node performed the operation, that TypeScript and JavaScript
are semantically equivalent, that the user understood either, or that the preparer did not retain
a copy.

The remaining separation is still material. A compromised Source Preparer must have no
Installation-root, Approval, Supervisor-evidence, Content Broker, updater/TUF, backend, or other
operational private key; no user-only Broker content; no plan-registration or approval-signing
authority; no grant/attempt/trust reset; and no backend, runtime, guest, or evidence-signing route.
After registration, different stored bytes must fail the independent plan-digest and role checks
rather than execute transparently. A Source Preparer compromise therefore changes planning and
what the user is asked to approve, but does not by itself create an approved attempt or guest.

Independent deterministic reproduction is not required for the limited claim that the user
approved exact emitted bytes produced by an enrolled Source Preparer. Repeating the same tool does
not prove semantics or correct signer logic. Offline reproduction is useful release, audit, and
retention evidence only. If Capsule instead requires faithful TypeScript erasure despite preparer
compromise, or removes the preparer from approval-understanding TCB, this topology is NO-GO until
an independent pre-approval implementation reproduces or verifies the transformation under a new
ADR.

## Compromise and exposure matrix

| Actor or failure | Can cause, learn, replace, retain, replay, or deny | Must remain unable to do |
| --- | --- | --- |
| Compromised daemon caller | Choose arbitrary original source and plans; replay exact prepare; abandon its prepared set; request registration with a known set; saturate daemon-facing capacity; learn its fixed plan projection | Select profile/options/emitted bytes; read the store by ID/path; call Supervisor-role methods; approve, use keys, release a registered set, or launch |
| Compromised native front | Observe/copy all accepted bodies and replies; bypass its own peer-derived context; invoke Go directly; corrupt in-process buffers; obtain full Source Preparer compromise | Reach any operational key/group, Broker content, Supervisor authority, backend, or guest |
| Compromised Go core/store owner | Create coherent malicious source sets; change lifecycle state; retain/replay/drop bytes; lie to callers; deny all operations | Make altered bytes satisfy an already registered digest; sign approval/evidence; create an attempt or guest |
| Compromised Node/Amaro worker | Read one input set; emit arbitrary bytes/refusal; flood/hang/crash; attack frame parsing and any inherited ambient authority | Access store, owner lock, listener/Mach rights, keys, network, subprocess/native loading, arbitrary files, or survive parent death **only after an exact child profile proves those denials** |
| Compromised immutable-store implementation | Replace, omit, retain, roll back, or replay coherent pre-registration state; corrupt refcounts; deny reads | Hide different post-registration bytes from independent digest/relationship validation |
| Compromised Supervisor | Read copied registered source objects, ignore validation, register unsafe state, and use its existing critical launch authority | This is already outside the Source Preparer containment claim; Broker approval remains a separate claim |
| Baseline same-user process | Probe services, race known names, pressure host resources, attempt link/replace/open/mmap/attach, and submit malformed IPC if it controls an enrolled caller | Obtain a writable store/package alias, trusted IPC delivery, worker/source buffers, keys, or task/debug authority in the validated baseline posture |

Proposal source is agent-facing data, so Source Preparer placement adds no user-content
confidentiality claim. It should still minimize disclosure. Full Disk Access, explicit foreign-
container grants, successful task ports, Accessibility, root, kernel, or administrator control use
the threat model's elevated or trusted-platform tiers and are not silently described as contained.

## Required approval rendering

The Broker independently decodes Supervisor-fetched registered plan and source-object copies. Its
trusted view must show, without daemon or Source Preparer display prose:

- installation, trust epoch, registration, plan, Source Preparer release, and bootstrap identity;
- entrypoint and every logical path;
- every original and executable media type, exact byte length, and digest;
- transformed versus byte-exact JavaScript pass-through status;
- record membership/order and both original/executable aggregate lengths;
- exact Node, Amaro, platform, architecture, source archive, distribution, executable, transformer
  profile, and normalized-options identity;
- source-map/source-URL absence and `reject-any`/zero-diagnostic disposition; and
- an explicit warning that the emitted JavaScript was selected by the enrolled Source Preparer,
  was not independently reproduced, and is not proven semantically equivalent to the authoring
  source.

Every authority-bearing field must be rendered from registered bytes. Exact original and emitted
text must be available through bounded plain-text inspection with control, bidi, normalization,
line-ending, BOM, and trailing-newline distinctions made unambiguous. The approval ceremony must
state whether all bytes were inspected; it must not imply full review from presence alone. Hidden
truncation, rich parsing, diff normalization, or an unrenderable field refuses before the signing
operation.

## Boundary controls required before P1

### Store and path custody

P1 must name one exact code-signature-associated, single-member protected app-data container (or a
separately reviewed mechanism with equivalent enforcement), its minimum macOS floor, signing and
entitlement requirements, SIP/user-override assumptions, fixed internal root, and update behavior.
No shared app group may contain source bytes, state, checkpoints, descriptors, or authority.

The store root is a trusted installed capability, never a caller string. Store operations use
directory-FD-relative no-follow traversal; create-exclusive regular staging files; exact device,
inode, owner, mode, link, length, and generation checks; close and unmap every writable alias;
file sync; atomic same-volume rename; directory sync; read-only reopen; complete graph rehash; and
fresh defensive copies. Logical paths are manifest data, never filesystem paths. Accessors return
no path, URL, bookmark, descriptor, mmap, mutable slice, directory capability, or store name.

The installed corpus must give an arbitrary baseline same-session process the in-progress name and
still deny writable open, link, replace, rename, clone, truncate, mapping, enumeration where
claimed, and attach. It must cover daemon, Supervisor, Broker, unsigned, copied, stale same-team,
current-but-wrong-role, and child worker identities. Mode bits, pathname secrecy, PID, process name,
EUID, CDHash, or same-Team identity alone are insufficient.

An identity-changing epoch must not leave a stale same-component build with new-store authority.
P1 must either select a fresh epoch-scoped single-member protected container/access group and
lock-held migration, or cite evidence that the selected container association rejects stale
enrolled bytes. If neither is implementable, the topology requires ADR revision before P1.

### Keys, access groups, entitlements, and services

The active manifest must positively record that the Source Preparer has no Installation-root,
Approval, Supervisor-evidence, Content Broker, updater/TUF, backend, or runtime key/group; no
broad/shared app group; no network, Hypervisor, Accessibility, Screen Recording, file-provider,
automation, debug, `get-task-allow`, dynamic-library exception, arbitrary file, or user-selected
file entitlement; and no backend/guest Mach service.

The Node worker has no Keychain or app-group entitlement and receives no inherited key/store/
listener/log/trust/owner-lock authority. Absence is verified from installed signatures and live
denial, not inferred from documentation.

The two role-specific services and eight methods in ADR-0032 are the maximum surface. Authentication
follows ADR-0029: listener requirement before activation; message-derived `SecCode`; exact channel,
Team, signing identifier, active CDHash, entitlement/runtime/debug, effective user and audit
session; flow slot; closed outer shape/caps; then installation, epoch, audience, and fixed
method-purpose before body copy. There is no Broker/updater listener, generic command, arbitrary
lookup/enumeration, caller path/descriptor, transform-by-digest, execute-time transform, or
caller-selected cleanup/profile/options method.

### Worker launch and confinement

P1 may freeze worker frames only after an implementable child mechanism is selected. The target is
an exact protected absolute Node/bootstrap identity, `posix_spawn`-style close-on-exec-default
behavior, explicit pipe-only FD mapping, fixed argv/env/cwd/stdio, no source before durable child
identity binding, and an OS-enforced worker profile denying store/container access, network, Mach
services, fork/exec/setsid, inspector/debug/task access, package/config/environment discovery,
native addons/FFI/dlopen outside the admitted closure, arbitrary filesystem/temp output, core
dumps, and persistent children.

The worker process tree is exactly one. CPU, wall time, address space/RSS, file size, descriptors,
threads/processes, stdout, stderr, request, result, and pipe buffering have hard external limits or
the dimension is explicitly unsupported. Node heap flags and measured peaks are not hard RSS
enforcement. Stdout and stderr are drained concurrently; cap-plus-one kills, continues bounded
drain as needed, and reaps rather than deadlocking on a full pipe. No partial/trailing/second frame
or exit status can become success.

A durable pre-spawn nonce and child-start barrier must close the spawn-to-identity gap. Source is
released only after PID/start/code identity and lease binding are durable. Parent death or control-
FD EOF before that release causes OS-owned or child self-termination. Recovery never kills by PID,
path, or code identity alone. Tests cover death after spawn syscall/before return, after return/
before durable bind, after bind/before input, during each pipe phase, after output/before reap, and
during kill/reap. If this isolation requires a new service/helper or cannot coexist with the exact
Node profile, stop and revise ADR-0032.

### Resource, cancellation, and diagnostic closure

One table must define per-service and process-wide connections, per-method in-flight calls,
transform workers, reads, body/reply/native/Go/store/worker copies, byte reservations, spawn rate,
and diagnostic cardinality. Four connections and eight simultaneous reads cannot both be global
claims. The aggregate reservation includes Node native RSS/parser amplification, XPC/native/Go
copies, pipe buffers, staging, rehash buffers, read/reply copies, and continuous drains—not only
the two source aggregates.

Saturation refuses before body copy, queue, worker launch, or store mutation. There is no unbounded
queue. Daemon saturation may deny preparation but cannot starve already retained Supervisor
retain/read/release work; otherwise that accepted denial consequence and priority rule must be
explicit.

Pre-dispatch cancellation has zero effect. After bridge dispatch, disconnect/cancellation is
response cancellation only; service-owned deadline and store truth determine whether one set
committed. `AbandonPreparedSourceSetV0` is not in-flight cancellation. The ten-second envelope must
include ingress, artifact verification, spawn, transform, validation, publication, cleanup, and
reply, with monotonic phase budgets and bounded kill/reap escalation.

Diagnostics are fixed codes and bounded numeric observations. Source, emitted bytes, logical/host
paths, Node stderr, arbitrary parser text, approval data, guest strings, and high-cardinality input
never enter replies or trusted logs.

### Genesis, update, epoch, rollback, and recovery

The current proposal has no closed channel for first-store creation or changing installation,
epoch, trust-snapshot, package, and migration authority. P1 must select a fixed-cap, versioned,
installation-root-authorized sealed genesis/update descriptor at a protected installer-owned
internal location. Ordinary methods are disabled while it is consumed. The Source Preparer remains
the sole store mutator; the installer supplies authorization and exact target facts, not records
or source bytes. Missing, replayed, partial, unknown, or indeterminate genesis/update state enters
repair-required. Startup never creates an empty authority store.

Channel epoch and immutable source-object/registration epoch are distinct. A current Supervisor
may read or release a retained old-epoch set only from its current authenticated channel and exact
durable registration/retention proof. A stale registration cannot gain a new approval or attempt.
Compatible migration preserves exact bytes and bindings without re-transformation; incompatible
migration refuses until forward repair. Old binaries refuse new stores. Coherent rollback remains
detectable only with an independently protected checkpoint; without one, a coherent older source
world is rollback-uncertain/repair-required and cannot re-enable approval or attempt.

Source-intent recovery precedes attempt enablement and is ordered with ADR-0029 startup recovery:
a committed registration retries one retain, an uncommitted intent retries one abort, and an
indeterminate result fences. Approval/attempt remains disabled until the exact retain
acknowledgement is durable. No recovery path registers a second plan, recreates a set, transforms,
or accepts replacement bytes.

### Retention, release, and archive

Physical blob release is not freezeable yet. ADR-0032 and the implementation plan currently
disagree about whether reproduction may be explicitly ended or requires an archive, while P1 asks
for a positive final release vector. P1 may define idempotent release refusal and a retained
tombstone, but it must not freeze a positive blob-decrement path until one rule is selected.

The eventual command must be authenticated Supervisor retained-state, bind the exact installation,
channel/object epochs, registration, source set, source descriptor, cohort, release intent, and
policy, and carry either a committed immutable reproduction-archive/checkpoint identity or an
explicit policy/epoch identity that ends reproduction retention. A free Boolean or cohort expiry
is insufficient. Supervisor archive eligibility includes the internal source reference and
resolved retain/release disposition; absence never means released. Same-role deduplicated blob
refcounts require exact replay, underflow, partial-release, corruption, and response-loss oracles.

Referenced bytes, tombstones, and replay state are never pressure-evicted. Until archive/retention
policy is accepted, capacity refuses. Secure deletion on APFS, SSD, snapshots, or backups is not
claimed.

## Required substitution and refusal corpus

P1 must assign a fixed outcome and before/after state oracle to every case below.

- **Role:** original/executable/profile/options/record/record-set/plan/internal-reference
  substitution, including accidentally equal digest bytes. Lookup and dedup keys include nominal
  role plus digest, length, and media.
- **Version:** worker, frame, method, store, profile, source-object, plan, registration, lifecycle,
  archive, and field-authority version; v0 presented to v1; v1 presented to v0; and the unselected
  626-byte arithmetic record. No fallback or dual active acceptance.
- **Installation/epoch:** wrong installation, channel epoch sequence/digest, object epoch,
  transition state, release epoch, and coherent restored state.
- **Source set:** unknown/stale/released/quarantined set, nonce/body conflict, request digest,
  source descriptor, prepared-set ID, registration intent, plan digest, registration, attempt,
  retain/release intent, refcount, hold, and tombstone swap.
- **Path/source/media:** zero/33 files; duplicate, unsorted, dot, overlong or bad-segment path;
  entrypoint/member/order/count mismatch; TypeScript/JavaScript media swap; transformed file without
  exactly one record; JavaScript pass-through not byte-identical or with a record; per-file and
  aggregate mismatch; invalid UTF-8/BOM; LF/CRLF, Unicode normalization, and trailing-newline
  distinctions.
- **Tool/profile/options:** Node, Amaro, source archive, distribution, executable, bootstrap,
  platform, architecture, mode, input/output media, source map, source URL, diagnostic policy/count,
  worker-reported digest/count/media/record, and installed artifact identity.
- **Record:** missing/extra/duplicate/reordered record, wrong original/emitted length/digest/media,
  options bytes/digest disagreement, cross-record membership, and record-set exact-byte collision.
- **Lifecycle/fault:** expiry equality and clock rollback, abandon/resolve race, response loss,
  cancellation at every phase, worker crash/hang/flood/OOM/signal/partial/trailing result, parent
  death at every lease boundary, publication/sync/rename/reopen ambiguity, retain/abort/release
  replay, store corruption, update interruption, rollback, and capacity exact/cap-plus-one.

Pre-authentication, outer-shape, cap, unsupported, binding, or capacity refusal causes zero body
copy, queue, worker, staging, ID/nonce, or durable mutation. A worker refusal/protocol/cap/timeout
creates no published set; it closes every FD/mapping, releases reservations, and removes only
durably proven unpublished staging. An indeterminate publication fences and never creates a
second set or empty store. Post-commit response loss replays one ID. Validation failure before
registration creates no registration/approval; only the exact pending intent may abort. Approval
read failure causes no signature. Executable read failure after attempt commit leaves the approval
consumed and attempt non-success, invokes zero backend calls, and never rolls authority back.
Release mismatch decrements nothing and deletes nothing.

The only refusal-side durable changes allowed are trusted effective-time high-water advancement,
an exact cleanup record for already-created internal staging/worker state, or tightening to
recovery-fenced, quarantined, or repair-required. None may create or widen plan, registration,
approval, attempt, release, backend, or guest authority.

## Field-authority finding

The current manifest correctly covers 164 top-level fields across 15 passive targets, but it cannot
satisfy the present P0 exit language. Source Preparer methods/store/worker objects do not yet have
canonical definitions, the current profile vocabulary names passive-fixture origin/validation, and
the verifier checks numbered top-level CDDL map entries rather than nested tuple members. One
profile also names one validator/resolver even though Source Preparer-originated claims receive
independent Source Preparer, Supervisor, and Broker validation phases.

P0 therefore retains the complete required field inventory and trust rules in this review. P1 must
land each closed canonical object, recursive field classification, calculated maximum, and verifier
support atomically before calling any bytes frozen. The vocabulary must distinguish origin from
all independent validators/resolvers and cover nested manifest entries, record tuples, worker
frames, request/reply/refusal fields, lifecycle/lease/generation/time state, registration intent,
internal source reference, retain/release/archive state, and cleanup observations. A parallel
manifest or prose-only exception is forbidden.

## Exact P1 entry criteria

P1 may start only when all of these are reviewed and retained:

1. ADR-0032 and the security matrices explicitly accept the native/Go/store/Node preparation TCB
   and the trusted-not-proven transformation claim.
2. Retained bounded platform evidence demonstrates the exact Source Preparer private-container
   mechanism, minimum OS/distribution assumptions, epoch isolation, and the complete baseline
   same-user negative-access matrix.
3. Retained bounded child-profile evidence demonstrates the worker's exact OS confinement, launch,
   artifact custody, FD/stdio, process-tree, resource, deadline, cancellation, parent-death,
   kill/reap, and cleanup mechanisms without a new authority. Failure to deny only the parent's
   source store keeps Node/Amaro in the full store TCB; failure of the network, arbitrary-file,
   native-loading, process-tree, resource, or death bounds is a stop condition.
4. The sealed genesis/update descriptor and sole-writer bootstrap/migration sequence are closed.
5. Exact key, Keychain/app-group, entitlement, Mach-service, content, backend, and guest absences
   are part of the active manifest and evidence plan.
6. The two-service/eight-method surface, every field origin/validator/resolver/authority effect,
   peer authentication, state transition, idempotency identity, and refusal-side effect are closed.
7. One source-set lifecycle table resolves trusted time, expiry/abandon races, nonce/tombstone
   replay, quarantine scope, retain/abort/release response loss, refcounting, and crash recovery.
8. Exact per-method encoded maxima and one complete process-wide connection, copy, memory, worker,
   read, reply, diagnostic, and resource envelope are calculable from the contracts.
9. Broker rendering and inspection satisfies this review's complete registered-byte disclosure.
10. Epoch migration, old-registration read/release, Supervisor source-intent startup ordering,
    rollback limitation, and no-empty-store behavior are fixed.
11. Supervisor archive/source retention is integrated; positive blob release stays unfrozen until
    the archive-or-explicit-end rule and shared-blob refcount semantics are selected.
12. The canonical field-authority vocabulary/schema/verifier can classify every nested P1 field
    and all independent validation phases in the same change as its canonical contract.
13. The full substitution and no-authority-state-change oracle above is represented in the P1
    fixture plan.
14. A later-reviewed modern ESM `.mjs`-only JavaScript contract remains the documented planning
    fallback if any selected TypeScript mechanism fails; CommonJS, package resolution, legacy Node
    module surface, and governed-runtime widening remain excluded.

## Stop conditions

Stop P1 and return to ADR review if any of these occurs:

- the exact installed container cannot deny baseline same-user store mutation;
- the official Node/profile cannot coexist with required sandbox, library-validation, store,
  network, process, native-loading, resource, or death constraints;
- parent death can leave an unidentifiable or unbounded worker/process tree;
- a new key/group, privileged helper, third service, generic method, second store writer,
  updater mutation route, Broker content access, Supervisor/Broker/daemon transformation, or
  backend/guest route is required;
- a caller path, URL, bookmark, descriptor, mmap, store enumeration, profile/options selector, or
  cleanup choice enters an active contract;
- independent faithful-transformation proof becomes a product requirement;
- safe complete Broker rendering cannot be defined;
- any field owner/validator/effect, exact maximum, queue/reservation, commit/cancellation edge,
  release authority, epoch/rollback behavior, or refusal-side state remains ambiguous;
- v0/v1 are simultaneously active, 562 bytes are reinterpreted, or 626 bytes are assumed as a
  layout; or
- restore/update/release/archive behavior can make a registered set abandonable, delete held
  bytes, recreate an empty store, publish a second set, decrement twice, or enable approval/attempt
  before retained-source acknowledgement.

If confinement requires moving transformation into the daemon, Broker, Supervisor, updater,
runtime, or backend, or requires another service/helper, revise ADR-0032 first. Do not implement
around the stop.

## Confidence and limitations

Confidence is high that the authority, lifecycle, field-classification, rendering, and evidence
gaps above exist because they follow directly from the checked-in contracts and retained negative
evidence. Confidence is medium on ultimate macOS feasibility: Gate B supports the protected-
container direction, but no signed Source Preparer/stock-Node child composition exists and this
review did not use an Apple identity or external documentation.

This is a design and evidence decision, not proof that a future implementation is impossible or
secure. No P1 bytes, store, worker, endpoint, consumer, runtime, backend, guest, or production claim
was created.
