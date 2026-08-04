# ADR-0021: Scope operational Keychain groups to security epochs

- Status: Proposed
- Date: 2026-07-31
- Refines: ADR-0012 and ADR-0013

## Context

Exact XPC enrollment and Keychain access groups enforce different identities. Gate B demonstrated
that a stale same-Team Broker denied by an exact code-directory requirement could still use a new
Secure Enclave key placed in its historical stable access group.

A follow-up provisioned two disjoint Broker groups and keys. Old-to-new and new-to-old private-key
access failed with `errSecMissingEntitlement`. A two-store transition model, 14 modeled process
deaths, and eight provisioned platform process deaths then exercised fencing, create-if-absent key
creation, public-key authorization, epoch commit, old-key retirement, component acceptance,
rollback, and forward repair.

## Proposed decision

Whenever an enrolled Broker or Supervisor identity changes in a security epoch:

- provision a fresh component access group and generate a fresh non-migrated Secure Enclave key;
- durably fence execution and invalidate old unused grants before creating target authority;
- create the target key idempotently and bind its exact public-key fingerprint into an
  installation-root-authorized target epoch;
- commit the new authorization once, logically reject the old key immediately, and keep execution
  disabled;
- permit restoration of the prior world only before the commit and only when its exact authority
  remains intact;
- require forward repair after commit, including physical old-key retirement and acceptance by
  every exact current component process; and
- re-enable execution only after the new authorization, key, component identities, old-key
  absence, and epoch state agree.

An always-on signing mediator is not selected for v0 because it creates a high-value signing oracle
and another IPC, availability, recovery, and update boundary. Reconsider it only if installed
Developer ID packaging cannot sustain security-epoch group transitions safely.

## Consequences

- “Per release” means per identity-changing security epoch, not every source commit.
- Provisioning, App ID/group capacity, installer cleanup authority, locked-Keychain behavior,
  migration, and distribution-package transitions remain acceptance blockers.
- Old-key cleanup needs a narrowly authenticated retained actor; it does not authorize the old
  component for ordinary work.
- A coherent restore of the complete older local world remains outside this mechanism and still
  needs an independent anti-rollback anchor if that threat is in scope.

## Evidence

- [`experiments/macos-authority-separation/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/macos-authority-separation/RESULTS.md)
- [`experiments/gate-b-key-rotation/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-b-key-rotation/RESULTS.md)
