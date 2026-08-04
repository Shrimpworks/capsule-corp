# ADR-0033: Select an enrolled BSD `flock` object for Supervisor ownership

- Status: Proposed
- Date: 2026-08-03
- Refines if accepted: ADR-0012, ADR-0018, ADR-0025, ADR-0029, and ADR-0031

## Context

The current E5 lifecycle path opens a fixed v1 snapshot and creates an in-memory owner-session ID.
Its injected coordinator serializes attempts only inside one process. It cannot refuse a second
Supervisor process before store access, migration, recovery, archive preparation, or adapter work.

ADR-0025 proposed a pre-created mode-0600 sibling lock file held for the process lifetime, but did
not select the macOS lock operation or define who creates, enrolls, preserves, and replaces the
object. Atomic store replacement also means locking the mutable snapshot inode is incorrect.

The retained local experiment compared POSIX record locks, macOS open-file-description locks,
BSD `flock`, and `O_EXLOCK` using only owned temporary files and child processes. It also exercised
path/object substitution and startup ordering. The experiment did not test an installed protected
container, another user/session, an update package, reboot, or product service.

## Proposed decision

### Exact primitive

Use one pre-created regular sibling object named by a closed installation bootstrap descriptor.
The Darwin implementation opens it relative to an already validated Supervisor-private state-root
directory descriptor with `openat(O_RDONLY | O_NOFOLLOW | O_CLOEXEC)`, validates the opened object
with `fstat`, and acquires `flock(fd, LOCK_EX | LOCK_NB)`. It retains that exact descriptor without
duplication for the complete Supervisor lifetime and releases it last during orderly shutdown.
Process death releases the lock when the last inherited description closes.

The pre-lock and post-lock object checks require:

- a regular file;
- effective UID equal to the enrolled per-user Supervisor UID;
- mode exactly `0600`;
- link count exactly one;
- device and inode equal to the enrolled bootstrap values; and
- a second `openat` of the closed lock name resolving to the same device/inode before store access.

The implementation verifies `FD_CLOEXEC` after open. It exposes no raw descriptor, duplicate,
unlock method, caller path, or create-if-missing fallback. A fork-only child would retain the open
file description, so the Supervisor must not create one; ordinary `fork`+`exec` children rely on
`CLOEXEC`, and the installed descriptor corpus must confirm the lock is absent after exec.

The advisory lock serializes cooperating enrolled Supervisor processes. It is not peer
authentication, storage protection, a key, a lease, a rollback anchor, or permission to mutate
state.

It is also distinct from the Source Preparer P0 boundary merged in PR #72 from head `a12041c` as
`2e268b0`.
Single-Supervisor exclusion does not provide source-store confidentiality or integrity, prove a
single-member protected container, confine a preparation worker, or close source-store
genesis/update authority. A future Source Preparer may compose an independently enrolled owner
lock only after those separate properties are selected and proven, using its own protected root,
installation/component identity, lock object, and owner session. It must never share the
Supervisor state root, lock descriptor, or owner-session domain.

### Trusted bootstrap and installation responsibility

The trusted containing application/installer creates the lock object exactly once while creating
a new Supervisor-private state root. It uses exclusive no-follow creation, sets mode `0600`, syncs
the file and directory, reopens it, and records this closed bootstrap projection in the signed
installation/epoch material:

- bootstrap format/version;
- installation ID and Supervisor ID;
- expected effective UID;
- state-root device/inode and closed store/lock entry names;
- owner-lock device/inode/mode/link-count policy; and
- expected store format and active trust-epoch binding.

The daemon, Broker, backend, updater network path, normal Supervisor startup, and store opener may
not create or replace the lock. Ordinary binary updates preserve the private state root and the
same lock inode. An updater that cannot preserve it must leave attempts disabled and enter the
prepared repair/new-installation path; it may not silently create another object.

This ADR selects no general same-installation replacement procedure. Store relocation, restore to
a new inode, or loss of the lock object requires a separately authorized offline repair/restore
ceremony that verifies retained state and an independent checkpoint, creates and enrolls a new
object, and commits a forward trust epoch before attempts can resume. If continuity cannot be
proven, Capsule creates a new installation identity and invalidates pending authority. Normal
startup treats a missing or mismatched object as repair-required.

The mutable store is not inode-bound because its commit protocol replaces that inode. After the
lock is held, the Supervisor opens the store by its closed sibling name relative to the same
retained state-root descriptor and validates its own file policy plus the installation,
Supervisor, epoch, format, digest, and cross-link contents. The bootstrap state-root and lock
identity, not the current snapshot inode, bind ownership to the installation.

### Startup, recovery, and shutdown contract

The process order is exact:

1. resolve and validate only the installation-configured private state-root capability and signed
   bootstrap projection;
2. open and validate the enrolled lock object without creation;
3. acquire the nonblocking exclusive `flock`;
4. revalidate the held object, directory entry, and `CLOEXEC` state, then issue one fresh nonzero
   owner-session ID from that live opaque capability;
5. open the authoritative store without creation, fully validate it, and bind it to that exact
   owner-session ID;
6. create the per-attempt coordinator with the same owner-session ID;
7. advance trusted time where allowed, enumerate the sorted recovery `AttemptID` set, and call only
   `Recover(AttemptID)`;
8. keep attempts and all four proposed IPC calls disabled until recovery is clean; and
9. retain the descriptor through all mutation, archive, backup, and repair checks, then close it
   only after listeners, work, store handles, and native/Go queues have stopped.

`EWOULDBLOCK`/`EAGAIN` is the fixed duplicate-owner result. It occurs before a store read, write,
recovery enumeration, archive operation, owner-session creation, or adapter call. Missing,
symlinked, non-regular, wrongly owned, wrongly permissioned, multiply linked, or identity-mismatched
objects are repair-required and are not rewritten. Unexpected descriptor loss or an invariant
violation terminates the process; it is not converted to ordinary readiness.

The owner-session ID remains a random in-process anti-confusion value. It is created only after the
OS lock plus its complete pre-store revalidation, changes on every reacquisition/reopen, and binds
the subsequently opened store, coordinator, and sealed effect/archive values to one live owner. It
neither acquires nor proves the OS lock and never substitutes for the retained descriptor.

### Same-UID and pathname limitation

`flock` is advisory. A process that can open the object may ignore the lock, hold it for denial of
service, or—if it can modify the parent directory—rename the held inode and install a separately
lockable replacement. The local experiment observed that exact split-object behavior. Device/inode
enrollment detects a replacement at startup but cannot prevent a later rename by an actor with
directory mutation authority.

Therefore product use is blocked on the installed Gate B property that the enrolled Supervisor's
private state root denies baseline same-UID pathname access to other processes. Mode bits and the
lock alone are insufficient. Full Disk Access, user-authorized foreign-container access, task-port
control, root/administrator, kernel compromise, and comparable elevated capabilities remain
outside this claim. No cross-user protection is claimed by this ADR or its local evidence.

## Alternatives considered

### POSIX `F_SETLK` record lock

Rejected. It is process-associated: fork did not inherit the lock, a second same-process open
joined it, and closing another descriptor for the same file released the process's locks. That
behavior is unnecessarily fragile for a long-lived Go/native process with migrations and offline
verification.

### macOS `F_OFD_SETLK` open-file-description lock

Not selected. On macOS 26.5.2 it had the desired duplicate/fork/exec/last-close behavior and
contended with both POSIX and BSD lock attempts. It is newly exposed in the macOS 26 SDK, while the
product support floor and compatibility path remain open. It supplies no path-replacement defense
that BSD `flock` lacks and provides no current benefit that justifies a newer platform dependency.

### `O_EXLOCK | O_NONBLOCK`

Not selected. It interoperated with `flock` and can acquire the same lock during `open`, but that
orders acquisition before descriptor identity validation and does not stop rename/unlink or a new
path inode. Explicit open, validate, `flock`, and revalidate makes the refusal boundary reviewable.

### Lock the store inode or rely on SQLite/launchd

Rejected. Atomic store commits replace the snapshot inode. SQLite is not yet the selected
production engine and the owner must also cover migrations, archives, backups, and offline repair.
launchd service identity and restart behavior do not lock the installation state for explicit
offline operations and do not replace store truth.

### Create-on-start sentinel

Rejected. An `O_CREAT|O_EXCL` marker becomes stale after process death and either blocks recovery
or requires unsafe absence/PID inference. Normal startup never creates the owner object.

## Consequences and blockers

- The selected mechanism can replace the installation-global injected owner assertion without
  adding a process, helper, daemon route, backend authority, or new public method.
- The existing per-attempt coordinator remains useful inside the sole owner; it is not the
  installation lock.
- Bounded G2 now composes the opaque held-owner guard with the current v1/no-guest startup: the
  owner-required opener precedes sorted recovery, store and coordinator share one session,
  post-open held-owner failure permanently fences, and store shutdown precedes descriptor release.
- Future migration, archive, backup, repair, and product mutation paths must require the same opaque
  guard; G2 intentionally exposes none of those ports.
- The local harness supports the primitive and ordering only. Protected-container, Apple-signed
  installed identity, wrong-user/session, update/restore, logout/login, reboot, abrupt shutdown,
  minimum-OS, and product Go/native integration remain acceptance blockers.
- The first bounded G3 checkpoint returned NO-GO before installed build. The authorized Apple
  Development certificate's common-name suffix says `W4QUR9FUL4`, but its subject OU and an exact
  signed-byte TeamIdentifier are `3DDR84M4JS`; no W4 profile is cached. Its noncredential fixture
  retains exact experimental role/state/bootstrap/update fields only. Protected-root creation by a
  trusted installer, the signed per-installation bootstrap envelope/parser, and descriptor-relative
  closed-store opening remain separate blockers even after matching credentials exist.
- This decision does not select a production store, authenticate IPC, prevent coherent rollback,
  provide source-store protection or worker confinement, provide continuous service, or admit a
  runtime/backend/guest.

## Evidence and implementation plan

- [Local owner-lock results](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/supervisor-owner-lock-boundary/RESULTS.md)
- [Passive G1 Go/Darwin owner package](../../internal/execution/installationowner/)
- [Bounded G2 owned v1/no-guest startup](../../internal/execution/registeredlifecycle/owned_startup.go)
- [G3 installed identity/session/update NO-GO](https://github.com/Shrimpworks/capsule-experiments/blob/3e9c9cbc3e0314439771151f1fd99c2b3a5a50b9/experiments/supervisor-owner-lock-installed-g3/RESULTS.md)
- [Owner-lock implementation and fault plan](../SUPERVISOR_OWNER_LOCK_PLAN.md)
- [ADR-0025 durable lifecycle](0025-colocate-durable-attempt-lifecycle-state.md)
- [ADR-0029 Supervisor topology](0029-select-authenticated-local-ipc-topology.md)
- [ADR-0031 archive boundary](0031-checkpoint-closed-supervisor-cohorts.md)
