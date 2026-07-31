# License-free feasibility spike results

Date: 2026-07-31

Status: research synthesis. This records development evidence and current decisions; it does not
promote any control, backend, profile, or component to production-ready.

## Outcome

An Apple Developer Program account is **not blocking the next implementation slice**. Local ad-hoc
signing and public macOS APIs were enough to test the serialization profile, exact live-process and
XPC code requirements, message-derived peer identity, cross-process descriptor custody, direct
VM-backed Containerization, no-network and resource controls, bounded output, controller crashes,
and real `SIGKILL` recovery of the trust-transition stores.

The account is still required before Capsule can validate its shipping identity and storage
boundary: Team-ID requirements, distinct distribution-signed components, provisioned Keychain
groups, protected app/app-group containers, production entitlements, Hardened Runtime,
notarization, and update/migration behavior.

## Current gate position

| Gate | License-free result | Current decision | Smallest remaining work |
| --- | --- | --- | --- |
| A: JCS/JWS | RFC 8785 JSON failed because the Swift/Foundation number representation differed and no maintained strict Swift JCS path passed. ES256 itself interoperated. | **Fail the original format.** | Do not revive generic JCS/JWS without new evidence. |
| A2: deterministic CBOR/COSE | Go, Swift, and TypeScript emitted identical bounded `ApprovalGrant` payload/protected bytes, verified all producers, and rejected all 12 negative vectors. | **Conditional pass for the bounded format.** | Freeze CDDL/profile rules, broaden confusion/resource/fuzz cases, review wrappers and dependencies. |
| B: macOS authority separation | Exact static and running-process hashes accepted v1 and denied stale v2; live launchd XPC accepted the exact client/copy, denied stale and unsigned fixtures, revalidated the message sender, and transferred an exact read-only FD. Debug state, missing Keychain entitlement, Secure Enclave signing, and noninteractive user-presence denial were also observed. | **Conditional pass for mechanism feasibility.** | Repeat with distinct Apple Development and distribution identities, Team-ID/identifier/entitlement/session requirements, provisioned key groups, protected stores, and product service packaging. |
| C: execution backend | Stock Apple Container 1.0.0 failed crash reconciliation and is disqualified as the security boundary. Direct Containerization 0.33.3 ran locally without a developer account. Empty interfaces/sockets, uid 1000, no-new-privileges, empty capabilities, read-only root, bounded tmpfs/output, exact memory, and discrete storage behavior were observed. A retained four-hunk patch enforced `pids.max=16` against root and non-root fork attacks. A simultaneous two-controller test removed one VM helper per controller crash. | **Stock API fail; direct backend remains a focused candidate, not a pass.** | Obtain a supported durable VM/helper identity or documented controller-death guarantee, then stress lifecycle races and management-vsock exposure. Govern or upstream the PID patch. |
| D: content custody | Snapshotting rejected symlink, directory, FIFO, device, Unix socket, oversized sparse input, mutation/substitution, cross-attempt reuse, duplicate redemption, and early/oversized output release. Exact read-only descriptors crossed both a real child process and live exact-peer XPC. | **Conditional pass for the custody contract.** | Put the contract behind production-signed Broker/Supervisor XPC, protected stores, the selected backend staging path, and a reviewed multi-process SQLite ledger. |
| E: Supervisor topology | Per-user direct Swift avoids the same-user global API service and needs no root helper. The exact package compiled with Command Line Tools; direct VM runs and controller-crash tests succeeded as development evidence. | **Provisional Swift/direct choice.** | Finalize only after Gate C identity/recovery and production-signed Gate B tests. Keep Go for portable policy/orchestration code; do not introduce a root helper. |
| F: trust transition/recovery | The 29-case state model passed. A new parent/child harness sent exact-PID `SIGKILL` after 23 durable checkpoints and repeated both ambiguous external-effect boundaries ten times each: 43 real process deaths, all reconciled to the expected stable or fail-closed state. | **Conditional pass for ordering/state semantics.** | Move the kill harness to the real Supervisor store/processes; add disk-full, WAL/checkpoint damage, APFS restore, installer faults, and VM power-cut testing. A coherent local rollback still needs an independent anchor if the threat model requires detection. |

## Important observations

### The Apple path is narrower, not blocked

The original stock CLI/API remains unsuitable: a restarted API server reported a still-running
helper as stopped. The lower-level direct library behaved materially better. Controller death
removed its Virtualization helper in the original single case and in a simultaneous two-controller
control, while an unrelated baseline helper remained alive.

That causal result is strong enough to continue one focused Apple experiment. It is not sufficient
for a terminal cleanup receipt because the public library exposes neither the helper PID nor a
durable identity that a restarted Supervisor can enumerate. A process-name scan must never become
the product control.

The missing PID cgroup was a small API-surface gap, not a missing guest mechanism. The retained
patch maps one optional configuration field to Containerization's existing `LinuxPids` OCI value.
With Capsule explicitly setting 16 and raising `RLIMIT_NPROC` to 256 for the attack, cgroup v2
denied both uid 1000 and uid 0 workloads after 13 children. The patch still needs validation for
invalid values, concurrency, compatibility, dependency updates, and preferably upstream review.

### The authority and custody design composes locally

The live XPC result closes an important feasibility uncertainty without overstating ad-hoc signing.
The OS applied an exact code requirement on the listener before message delivery. Accepted
messages were independently tied to their audit-token sender through `SecCodeCreateWithXPCMessage`,
then carried an already-open read-only descriptor. The stale client did not reach the protocol;
the exact authenticated client with a malformed operation reached it and was rejected there.

This supports the planned split:

```text
code requirement + message-derived sender  -> component authority
typed installation/epoch/attempt binding   -> protocol authority
one-use Broker ledger                      -> content authority
read-only FD copied into attempt storage   -> byte custody
```

Ad-hoc exact hashes prove that the mechanics compose. They do not replace the production signer,
Team ID, component identifiers, entitlements, session checks, or protected stores.

### The update model now survived real process death

The first Gate F model raised an exception after closing SQLite, which tested ordering but not an
unclean process. The retained child harness now pauses after the committed checkpoint and is killed
with uncatchable `SIGKILL`; neither SQLite connection is closed. A fresh process opens the WAL-backed
stores and reconciles the result. Every incomplete update remained execution-disabled, the final
re-enable checkpoint recovered stable, consumed grants stayed consumed, backend-create ambiguity
kept cleanup responsibility, and completed external release was observed rather than rolled back.

This is still process-crash evidence, not power-loss evidence. It says the state machine is worth
implementing; it does not certify the final macOS file layout or installer protocol.

## Work that genuinely waits for Apple credentials

1. Distinct daemon, Broker, Supervisor, updater, and installation-authority identifiers signed by
   one real Team ID, first with Apple Development and then the intended distribution channel.
2. Live XPC negative matrix for same-team/wrong-role, stale accepted hash, wrong entitlement,
   `get-task-allow`, debug attachment, effective UID/audit session, reconnect, activation, update,
   and service replacement.
3. Provisioned, disjoint data-protection Keychain access groups; persistent Secure Enclave evidence
   and approval keys; proof that daemon and sibling components cannot use them.
4. Distinct protected Broker/Supervisor stores and any narrowly justified IPC-only app group,
   including user override, Full Disk Access, fast-user-switching, MDM, and migration cases.
5. Hardened Runtime, library validation, Developer ID or App Store packaging, notarization,
   designated-requirement stability, and team/profile/container migration.

Successful interactive user-presence behavior can be explored without a paid account using an
ephemeral key, but it requires a person at the Mac and does not validate persistent provisioned key
custody. It is intentionally separate from the autonomous run recorded here.

## License-free next slice

Proceed without waiting for the Apple account:

1. Draft the narrow CBOR/COSE ADR and CDDL from Gate A2, preserving exact registered payload bytes.
2. Implement the Go authoritative state machines and ports around fake Broker/backend adapters,
   including the Gate F transition fence, component-acceptance barrier, cleanup obligations, and
   idempotent release intent.
3. Build the per-user Swift Supervisor skeleton around a fake backend, with the live XPC exact-peer
   harness retained as a development conformance test—not copied as production ad-hoc policy.
4. Run the final direct-Containerization identity/recovery spike and management-vsock attack
   characterization. Stop the Apple path if it requires private APIs or cannot produce independently
   attributable teardown evidence.
5. Translate the content-custody contract into schemas/state tables and a real multi-process
   ledger harness, while keeping all user paths and redeemable handles out of daemon-visible data.

## Evidence index

- Gate A: [`../experiments/gate-a-signing-canonicalization/README.md`](../experiments/gate-a-signing-canonicalization/README.md)
- Gate A2: [`../experiments/gate-a2-cbor-cose/README.md`](../experiments/gate-a2-cbor-cose/README.md)
- Gate B: [`../experiments/macos-authority-separation/RESULTS.md`](../experiments/macos-authority-separation/RESULTS.md)
- Gate C stock API: [`../experiments/apple-container-gate-c/RESULTS.md`](../experiments/apple-container-gate-c/RESULTS.md)
- Gate C direct backend: [`../experiments/apple-containerization-direct/RESULTS.md`](../experiments/apple-containerization-direct/RESULTS.md)
- Gate D: [`../experiments/gate-d-content-custody/README.md`](../experiments/gate-d-content-custody/README.md)
- Gate E: [`../experiments/gate-e-supervisor-topology/RESULTS.md`](../experiments/gate-e-supervisor-topology/RESULTS.md)
- Gate F: [`../experiments/gate-f-trust-transition/RESULTS.md`](../experiments/gate-f-trust-transition/RESULTS.md)
