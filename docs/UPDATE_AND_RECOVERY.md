# Update and Recovery

Status: intended crash-safe model; trust-transition behavior requires Gate F.

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
9. Start the new Supervisor in `pending-verification` with execution disabled.
10. Validate all installed components, key access, stores, and migrations against the prepared
    record.
11. Create and sign epoch N+1, linking the prior epoch and transition reason.
12. Atomically commit the current epoch checkpoint and local trust snapshot reference.
13. Optionally witness the new digest.
14. Re-enable attempts only after every required component accepts the same epoch.

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
