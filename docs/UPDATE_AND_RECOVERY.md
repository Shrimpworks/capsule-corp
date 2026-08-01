# Update and Recovery

Status: intended crash-safe model; Gate F process/storage-ordering evidence is positive, while the
real Supervisor store, installer, backend, Keychain anchor, and power-loss behavior remain pending.

## Objectives

An update must not create a window in which old and new components, keys, policy, profile state, or
storage formats are silently mixed. Recovery must not reset approval consumption, trust history,
cleanup obligations, or evidence simply because components were reinstalled.

## Trust refresh versus component transition

A repository refresh verifies new external metadata and may create a new local `TrustSnapshot`.
It does not automatically change installed code.

A new trust epoch is required when active component identities, their relevant entitlements/code
requirements, required policy/profile registry checkpoint, local authority set, or storage format
changes. Routine timestamp refresh that does not change active enrolled state need not create an
epoch.

## Prepared update ceremony

The intended v0 flow is:

1. Refresh and verify pinned TUF metadata outside the Supervisor.
2. Download exact artifacts and verify trusted length/hash.
3. Statically validate signatures, code requirements, entitlements, platform range, and package
   layout.
4. Verify runtime/profile/review/validation records required by the target release.
5. Stop accepting new attempts and wait for or explicitly cancel/reconcile active attempts.
6. Create a signed `PreparedUpdate` binding the current epoch, exact target identities, policy/
   profile checkpoint, storage migration, and recovery plan.
7. Require user/admin authorization for the v0 trust-changing transition.
8. Install/swap through the supported platform mechanism.
9. Enter `swapping` with execution disabled and record every installer, migration, and Keychain
   effect through durable intent plus observation/reconciliation state.
10. Start the new Supervisor in `pending-verification` and validate all installed components, key
    access, stores, and migrations against the prepared record.
11. Enter `finalizing-epoch`; stage and sign epoch N+1, linking the prior epoch and transition
    reason.
12. Atomically commit the authoritative epoch pointer, execution-disabled state, policy/profile/
    trust references, and storage version.
13. Enter `awaiting-component-acceptance`; optionally witness the new digest and require every
    exact current component process to accept the same transition and epoch.
14. Re-enable attempts only in the final transaction after all required acceptances and a fresh
    self-check.

No continuously running daemon or Supervisor receives installation-root authority merely to make
background updates convenient.

## Interrupted update

Any incomplete transition enters `repair-required`. Starting some components successfully does not
restore ordinary execution.

Recovery reads the prepared record and determines whether it can:

- safely finish installation and finalize the planned epoch;
- restore the exact prior enrolled epoch without losing authoritative state;
- require an explicit repair/reinstall ceremony.

It never invents a new epoch, resets grants, or marks unknown backends destroyed.

Before the N+1 commit, repair may restore the exact prior world only when the prepared transition
declared that path reversible and no later external effect makes it unsafe. After the commit,
repair is forward-only: older bytes or authority become a newly authorized N+2 transition rather
than rewinding the current pointer.

## Repair

A trusted signed installer may restore exact release components, but repair must preserve or
explicitly replace:

- grant consumption and attempt history;
- backend handles/cleanup leases;
- manifest and epoch history;
- key revocation/replacement state;
- content retention/quarantine state;
- enforcement transcripts and receipt indexes.

Repair either restores the current enrolled epoch exactly or creates a clearly identified new
epoch after authorization. It cannot silently rewrite history.

## Job and attempt recovery

On Supervisor restart:

1. Load and verify installation/epoch state before accepting peers.
2. Enumerate nonterminal attempts and cleanup leases.
3. Preserve every consumed approval.
4. Reconcile backend handles with independently enumerated backend state where supported.
5. Terminate/destroy orphaned or indeterminate guests.
6. Quarantine collected but unreleased content until terminal integrity and teardown state is
   known.
7. Append recovery events and sign only the terminal classification supported by observations.

A missing process, VM, container, or handle is not alone proof of successful destruction.

Direct Apple Containerization cannot currently satisfy this recovery step through supported public
APIs because it exposes no durable host VM/helper identity or restart enumeration. It therefore
remains development-only and unresolved after ambiguous controller loss. The libkrun/HVF candidate
uses a durable-record-before-start handshake and a live PID/start/code-identity tuple; its spike
passed controller death before record, after record, and after start. A same-host installed
follow-up exactly recovered six reparented runners, but that depended on the enrolled
`AbandonProcessGroup=true` launch profile and did not cover complete distribution, logout/login,
sleep/wake, reboot, pressure, failed-signal races, or power interruption. Those cases and product
store composition remain required before claiming ordinary terminal cleanup. The OCI/gVisor
candidate must independently demonstrate exact engine/runtime identity and recovery.

## Multi-store saga

Cross-component operations cannot use one database transaction. Each message is idempotent and
binds installation, epoch, registration, attempt, content identity, operation, sequence, and expiry
as appropriate.

Important saga checkpoints include:

- plan registered;
- approval accepted/consumed;
- attempt created;
- content handle issued/redeemed;
- backend create intent/handle persisted;
- staged digests verified;
- collection manifest committed;
- terminal integrity result;
- user content released;
- backend destroyed/reconciled;
- receipt composed/indexed.

Replay, duplicate, out-of-order, and missing messages produce stable no-op or recovery states—not a
wider authorization.

## Durable store requirements

Gate F process and fault-injection evidence makes these v0 requirements explicit:

- one authoritative Supervisor writer, with `BUSY`, `LOCKED`, `FULL`, `IOERR`, `CORRUPT`, and
  failed commit treated as refusal outcomes;
- expected state sequence, epoch digest, transition ID, and execution-enable compare-and-swap on
  every security transition;
- a cleanup/reconciliation intent committed before every external backend, installer, Keychain,
  migration, or release effect;
- no automatic recreation of a missing, corrupt, truncated, schema-incompatible, or checkpoint-
  mismatched authoritative store;
- an explicitly tested SQLite journal/WAL/checkpoint, sync, backup, restore, and migration policy;
- same-volume verified temporary files, file sync, atomic rename, and directory sync for critical
  standalone state;
- explicit fail-closed clock rollback/unavailability handling; and
- bounded store growth and enough already-durable recovery identity to survive capacity exhaustion
  after an external effect.

A database plus ordinary sidecar on the same rollback domain detects partial restore but not a
coherent older world. If coherent rollback detection is required, the latest checkpoint must be
bound to a suitably independent Keychain/platform anchor or external witness.

## Lost root, OS reinstall, and machine replacement

Loss of the installation root creates a new installation identity. The old identity is marked
replaced or abandoned through whatever local/external record remains available. Private keys are
not silently restored from ordinary backup.

Pending approvals are invalid across the new identity. User content migration, if offered later,
is a separate explicit operation and does not migrate signing authority.

Portable multi-device identity and automated recovery are deferred. Optional DIDs can express
replacement relationships externally, but local trust still requires enrollment.

## Rollback limitation

Signed epoch chains are sequence-ordered and detect partial or mismatched rollback. They cannot
alone detect an attacker who coherently restores an entire older valid local world. Strong
anti-rollback requires a non-rollbackable anchor, independently protected checkpoint, or external
witness. Documentation and receipts must not claim monotonicity without one.

## Fault-injection plan

Gate F interrupts before and after every durable write, file swap, process stop/start, migration,
epoch checkpoint, grant consumption, backend creation, collection, release, and teardown step.

Required outcomes:

- one approval never becomes unused after consumption;
- no mixed epoch becomes execution-ready;
- content is not released from an indeterminate attempt;
- an orphan retains a cleanup obligation;
- repair cannot erase history;
- coherent rollback limitations remain visible.
