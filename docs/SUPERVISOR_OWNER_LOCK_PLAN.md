# Supervisor owner-lock implementation and fault plan

Status: proposed and unimplemented product boundary. The retained development-only experiment
selects BSD `flock` semantics for later implementation; no product startup, service, protected
store, backend, runtime, or guest is wired by this plan.

Decision: [Proposed ADR-0033](adr/0033-select-enrolled-flock-supervisor-owner.md).

## Defensive scope and invariant

Defensively validate and later implement only the installation-global single-Supervisor owner
guard using repository tests, owned temporary state roots/processes, and separately authorized
installed fixtures. Do not access another user, session, process, store, identity, credential,
service, backend, runtime, or guest.

No store read/mutation, migration, recovery enumeration, `Recover`, archive, backup, repair, IPC
dispatch, or adapter call may occur before one exact enrolled object is exclusively held. The lock
does not grant any of those operations; existing typed state and authority checks remain required.

This plan does not close the Source Preparer P0 HOLD merged in PR #72 (head `a12041c`, merge
`2e268b0`). Process
exclusion is not source-store confidentiality or integrity, protected-container membership,
worker confinement, or store-genesis/update authority. Any later composition uses a separate
Source Preparer root, component identity, enrolled lock, and session only after those boundaries
are independently selected and evidenced; it does not reuse Supervisor ownership material.

## Narrow Go/Darwin boundary

The implementation target is an internal opaque capability, not a public or IPC contract:

```go
type OwnerLockEnrollment struct {
    InstallationID  v0candidate.InstallationID
    SupervisorID    v0candidate.SupervisorID
    ExpectedUID     uint32
    StateRootDevice uint64
    StateRootInode  uint64
    LockDevice      uint64
    LockInode       uint64
}

type InstallationOwner interface {
    OwnerSessionID() lifecyclestate.OwnerSessionID
    CheckHeld(context.Context) error
    CloseAfterShutdown(context.Context) error
    sealedInstallationOwner()
}

func AcquireDarwinInstallationOwner(
    context.Context,
    TrustedStateRoot,
    OwnerLockEnrollment,
) (InstallationOwner, error)
```

`TrustedStateRoot` is constructed only from installed bootstrap/platform state. It retains a
validated directory descriptor and closed entry names; no caller path crosses this API. The Darwin
implementation uses the smallest reviewed syscall surface: `openat`, `fstat`, `fcntl(F_GETFD)`,
`flock`, and `close`. A native shim is unnecessary if the pinned Go syscall dependency exposes
those exact operations; any shim remains in-process, method-specific, and descriptor-opaque.

The capability owns exactly one descriptor. It has no `File`, `FD`, `Dup`, `Unlock`, or transfer
accessor. `CheckHeld` verifies the wrapper is live, the descriptor identity remains enrolled, and
`FD_CLOEXEC` remains set; it cannot prove that an unrelated process obeys advisory locks. Store
openers, migrations, lifecycle transactions, archive sealed values, backups, and repairs accept
the opaque capability and reject a mismatched owner session.

The first integration must add an owner-required v1 opener rather than silently changing the
existing conformance constructor. Existing path-based constructors remain explicitly test-only
until all callers migrate. No product build may expose an ownerless store opener.

## Bootstrap record and store binding

Define one closed bootstrap projection and byte-exact fixture before implementation. It binds the
installation/Supervisor IDs, expected UID, protected state-root device/inode, closed lock/store
names, enrolled lock device/inode/policy, expected store format, and trust epoch. Unknown fields,
versions, missing identities, zero values, path separators, `.`/`..`, or non-ASCII entry names
refuse.

The trusted installer creates the state root and lock, syncs them, reopens them, and signs/enrolls
the observed identities. Startup reads no daemon-supplied path. The store may change inode on each
atomic transaction; its closed name and decoded installation/Supervisor/epoch state, not inode,
bind it to the held root and owner.

## Fault and recovery oracles

Every case records store bytes/digest, marker inventory, owner-session presence, recovery calls,
archive calls, and adapter calls before and after:

| Boundary/fault | Required result |
| --- | --- |
| bootstrap/state-root missing, symlinked, replaced, wrong identity | repair-required; zero lock/store/core/adapter work |
| lock missing, symlinked, non-regular, wrong UID/mode/dev/inode/link count | repair-required; no creation or rewrite |
| open/fstat/CLOEXEC failure | fixed local/repair refusal; zero store/core/adapter work |
| `flock` busy | fixed duplicate-owner result; zero store read/mutation/recovery/archive/adapter work |
| death before lock | no owner; successor follows ordinary startup |
| death after lock and before store open | descriptor closes; successor validates unchanged store |
| entry replacement between open, lock, and recheck | repair-required; held descriptor closes; zero store/core/adapter work |
| store missing/corrupt/unsupported after lock | lock remains held until process exit; original bytes unchanged; repair-required |
| death during store validation or recovery enumeration | successor reacquires and restarts full validation/enumeration |
| death during each `Recover(AttemptID)` | existing intent/checkpoint oracle decides; no second owner overlaps |
| fork-only child before exec | forbidden integration path; test proves inherited lock delays successor |
| exec child | lock descriptor absent under `CLOEXEC`; exact installed FD corpus verifies it |
| accidental close/unlock/descriptor mismatch | invariant termination before another mutation |
| rename/unlink/replacement by authorized test fixture | observe split inode; enrollment mismatch refuses new object; record same-UID limitation |
| process death / abrupt kill | last-close releases; successor reopens and fully recovers |

Concurrent duplicate-start tests must use independent process descriptions, not goroutines or two
objects in one process. One contender holds the lock while the other proves no store read marker,
temporary transaction, owner session, recovery call, archive file, or fake call exists.

## Implementation slices

### G1: passive bootstrap and opaque owner types

- Add closed enrollment/state-root/owner-session types and validation.
- Add Darwin build-tagged acquisition using owned temporary roots only.
- Retain process, descriptor, inheritance, replacement, and no-state fault tests.
- Adapt the existing `OfflineMigrationLock` assertion to an actual held owner in tests.

Acceptance: internal package only; no product command/service wiring, consumer, backend, or guest.

### G2: owner-required fixed-store and startup composition

- Add owner-required v1 open/migration and the future v2/archive ports.
- Generate the owner session only after lock plus store validation.
- Compose the no-guest fake startup order and all E5 recovery tests across child processes.
- Keep the per-attempt coordinator under the same owner session.

Acceptance: disposable process harness only, `FakeBackend.CreatesGuest() == false`, duplicate owner
refuses before store mutation or fake work, and ownerless product constructors are impossible.

### G3: installed identity/session/update evidence

Resource status: current Individual membership is Team `W4QUR9FUL4`, and local read-only discovery
reports a valid Apple Development identity for that Team. G3 can proceed after exact W4 role
identifiers, entitlements, and profiles are deliberately created. The three Xcode 26.6-cached
profiles are historical Team `3DDR84M4JS` fixtures and are not reusable for W4 tests. A separate
Developer ID Application identity for historical Team `3DDR84M4JS` is later distribution authority
requiring explicit authorization and matching-Team package design; it is not W4 development
evidence and does not make Developer ID/notarization work current. Paid owned clean-host/minimum-OS
coverage remains deferred activation evidence.

- Package the exact Supervisor fixture under the selected protected private state root.
- Exercise Apple Development then Developer ID/notarized builds, wrong/stale identity, wrong
  user/session, fast-user switching, logout/login, reboot, launchd restart/backoff, update
  preservation, prepared repair, lock loss, state-root relocation, and protected-path attacks.
- Read back lock/root/store identities and inherited descriptors after every transition.

Acceptance: exact tested distribution matrix only. This does not accept the production store,
authenticate the proposed IPC calls, prove cross-user protection generally, or claim production
readiness.

## Required focused tests

- `TestDarwinOwnerAcceptsOnlyEnrolledRegularObject`;
- `TestDarwinOwnerRefusesSymlinkModeOwnerLinkAndIdentityMismatch`;
- `TestDarwinOwnerDoesNotCreateMissingObject`;
- `TestDarwinOwnerDuplicateProcessRefusesBeforeStoreRead`;
- `TestDarwinOwnerDescriptorDupForkExecAndCLOEXEC`;
- `TestDarwinOwnerProcessDeathAllowsFullReopenRecovery`;
- `TestDarwinOwnerRenameUnlinkReplacementNeverCreatesSecondEnrolledOwner`;
- `TestOwnerRequiredStoreOpenRejectsOwnerSessionMismatch`;
- `TestOwnerRequiredMigrationChecksHeldOwnerAtEveryCommitBoundary`;
- `TestOwnerRequiredStartupRecoversAttemptIDsBeforeReadiness`;
- `TestOwnerRequiredArchiveAndBackupRejectAfterOwnerLoss`; and
- `TestDuplicateOwnerMakesZeroFakeCallsAndFakeStillCreatesNoGuest`.

## Remaining blockers

The retained local semantics do not close protected installed storage, signing/entitlements,
session/reboot/update behavior, support floor, real power loss, production engine, coherent
restore/rollback, authenticated IPC, production approvals, content/evidence, runtime/backend, or
guest gates. `flock` remains advisory and supplies no same-UID containment without the installed
protected-directory property. It also supplies no Source Preparer store protection or worker
boundary.
