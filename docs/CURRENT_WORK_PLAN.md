# Current work plan

Date: 2026-08-10

Work item: close the CL4 `AMEND` findings for the exact passive/no-listener C4 evidence claim.

Status: `PASSED` for this planning and documentation scope.

Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`.

Product admission and the installed security boundary: `BLOCKED`.

This is the current execution index. Detailed security claims remain in the linked ADRs, readiness
map, experiment checkpoints, and evidence ledger. A completed passive contract or controlled
experiment is not an activated product path.

## Reconciled baseline

This closure started from fetched `origin/main` commit
`133669ba4b61b4e09baddfd6f7bb0a922ca49c8a`, merge PR #248. PR #248 is the canonical predecessor
that retained R3, CL4's `PASSED` audit and historical `AMEND` disposition, and the focused C4
implementation dependency. The follow-up now closes the exact boundary and independent-verifier
findings, so the passive/no-listener C4 evidence claim is `PASSED`. C4, C5a, R3, and CL4 change no
runtime, governed-profile, installed, or product bytes.

Today closed or retained these exact slices:

- the passive fixed-store threshold checker and its narrowed re-evaluated scope;
- the passive five-method native-XPC contract and adversarial fixture expansion, including C4's
  approval/attempt caps, deadlines, replay, and response-loss tables;
- the proposed Supervisor authority-epoch decision and inert experiment packet;
- the documentation-only F6 SQLite research and executable experiment packet;
- issue labels, issue forms, and pull-request categorization;
- four correctness defects in daemon shutdown, completion-store creation/decoding, and passive
  host-runner defensive copying;
- a new-code exported-contract documentation gate;
- a completed `lifecyclestate` exported-contract documentation batch plus current-checkout lint
  isolation from nested local agent worktrees;
- focused coverage increases and one behavior-preserving archive-state refactor;
- the exact v19 benign guest and v27 30-marker fixed hostile-denial experiment checkpoints; and
- C5a's passive typed source/input/completion byte contract, deterministic fixtures, ordered
  refusal/state/restoration cases, and independent Go/Node verification; and
- R3's canonical read-only Apple-platform evidence brief, with Supervisor-owned durable approval
  authority, experiment-only accessibility/signature candidates, and an explicit installed-harness
  authorization boundary.

The generated repository totals are 95 rules, 502 cases, 624 fixtures, and 1,203 classified fields
across 95 profiles and 60 targets.

## Product-critical dependency graph

```text
raw v10-v27 evidence publication ------------------------------┐
                                                               │
ADR-0029 decision review + native XPC research -> S3 harness --+-> installed authenticated IPC
                                                               │
ADR-0045 E1 identity separation -> key/service/root corpus -----+
                                                               │
C4 passive approval/attempt evidence (PASSED) + R3 passive research (PASSED)
                                          -> disposable installed signing harness (BLOCKED) ---┐
installed authenticated IPC boundary (BLOCKED) ----------------------------------------------+-> product Broker/approval/attempt wiring
                                                                                              -> protected one-attempt path
                                                               │
typed transport design -> passive contract (`PASSED`) -> controlled harness ┐
installed service/identity evidence ------------------------------+-> pre-admission profile matrix
                                                                  -> admission review
                                                                  -> sealed real adapter
                                                               │
existing durable completion-last oracle ------------------------+-> real completion/absence proof
                                                                  -> installed hostile corpus
                                                                  -> internal alpha PASSED

F6 experiment -> engine decision -> restore/continuity -----------> external alpha only
```

The three top evidence lanes may run in parallel. Installed composition begins only after the
relevant lane passes. Product completion may reuse the existing completion-last semantics only
after real runner identity, result integrity, teardown, and authoritative absence exist.

## Next work owned by Codex

These are implementation or integration tasks. Each retained task uses its own `codex/<topic>`
branch and pull request unless the orchestrator explicitly groups it before work begins.

| Order | Work item | Current status | Start condition and acceptance boundary |
| --- | --- | --- | --- |
| C1 | Publish the raw v10-v27 harness, reports, manifests, and receipts to `Shrimpworks/capsule-experiments` | `BLOCKED` | Requires access to the owner-controlled disposable experiment workspace. Publish one verified immutable archive commit, then replace local-only evidence pointers with exact links. No new execution or broader claim. |
| C2 | Run the one-time native XPC S3 harness | `BLOCKED` | CL1 and R1 are `PASSED`; execution still requires explicit authorization naming `Shrimpworks/capsule-experiments`, the owned Mac, exact fixture digest, disposable process/service names, and defensive no-product scope. Retain OS pre-delivery peer refusal, message-derived identity, EUID/session, copy, cap, deadline, interruption, response-loss, and process-fault evidence. |
| C3 | Run ADR-0045 E1 Supervisor-authority-epoch identity separation | `BLOCKED` | Requires explicit Apple Development authorization for the named Team-`3DDR84M4JS` profiles and owned disposable containers. Prove current/legacy cross-container denial and cleanup; stop before Keychain, service registration, protected root, store, runtime, or guest. |
| C4 | Freeze `SubmitApprovalV0` and `RequestAttemptV0` | `PASSED` | PR #248 is the canonical predecessor. The focused follow-up freezes ordered 4,999/5,000/5,001-ms cases for both methods, equality-as-expired behavior, complete closed dictionaries/maps, every ordered field, required `noState`, cancellation/deadline commit truth, all 20 foreign-tag collisions, complete refusal and five-entry response-loss tables, and bounded Go/Node mutation proofs. No listener, signer, store consumer, process, or guest is active. |
| C5a | Freeze the final typed source/input/completion transport | `PASSED` | The passive v1 contract freezes exact 152/160/64-byte layouts, 262,144-byte payload caps, completion cap-plus-one, big-endian bindings, four statuses, canonical JSON, refusal precedence, monotonic state/fault behavior, endpoint custody, completion-last projection, deterministic fixtures, and independent Go/Node verification. No endpoint, process, runtime, backend, guest, or store mutation occurred. |
| C5b | Run the controlled typed-transport harness | `BLOCKED` | Requires C5a and separate authorization naming the exact runtime/profile successor and owned disposable guest. Retain directional copy, chunk/cap+1, stall/reset/cancel, descriptor substitution, response-loss, completion-last, and restoration evidence without making an admission decision. |
| C6a | Build the installed authenticated service and protected-state boundary | `BLOCKED` | Requires passed C2 and C3 evidence plus the accepted-or-revised ADR-0029 result, then separate authorization for the Keychain/service/protected-root corpus. Implement only method-specific listeners, peer authentication, owner/store startup, and the four passively frozen Supervisor consumers. |
| C6b1 | Run the disposable installed Broker signing evidence harness | `BLOCKED` | R3 is `PASSED`; execution still requires separate authorization naming the owned Mac, exact development-signed Broker test target/profile/access group, disposable user/container, immutable fixture digest, permitted Keychain mutations, and every destructive row. No product consumer, runtime, backend, VM, or guest. |
| C6b2 | Connect the product Broker and approval/attempt methods | `BLOCKED` | C4 is now `PASSED`; C6a and C6b1 remain required. Implement native rendering/UI, installed signing/public-key verification, and method-specific `SubmitApprovalV0`/`RequestAttemptV0` consumers without runtime or guest activation. Research and passive conformance cannot satisfy either installed dependency. |
| C6c | Wire attempt admission and the fixed-store stop policy | `BLOCKED` | Requires C6b2 and an explicit decision for p95 provenance/window/lifetime and any persistent timing-trip semantics. Apply the re-evaluated guard transaction-locally after owner/full verification and before a new consume/create mutation; replay of an existing `AttemptID` converges first. |
| C6d | Run the pre-admission installed runtime/profile matrix | `BLOCKED` | Requires C1, C5b, C6a, governed release/artifact review, and separate authorization naming the exact signed-installed candidate and owned test environments. Retain identity compatibility, runtime/root/loader restoration, transport, teardown/recovery, and the required broader lifecycle/platform evidence without accepting user source or making an admission decision. |
| C7 | Review one exact runtime/profile candidate for admission | `BLOCKED` | Requires C6d. Produce an explicit admit-or-refuse result over the exact candidate and retained evidence; controlled v19/v27 experiments alone cannot admit it. |
| C8 | Connect the sealed real adapter and completion-last path | `BLOCKED` | Requires C6c, C7, and a separately authorized owned guest. Execute only by committed `AttemptID`; consume real result-integrity, runner, teardown, and absence facts. |
| C9 | Run the installed hostile-`.mjs` admission corpus | `BLOCKED` | Requires C8. Response loss, restart, recovery, restoration, pressure, sleep/wake, update, and the minimum hostile source/authority/transport/root/lifecycle cases must converge in the exact signed-installed profile. |

Independent repository-quality work may continue without changing security claims: issue #217 in
one-package documentation batches, next `registrationstate`; issue #219 as sequential
behavior-preserving archive refactors;
and issue #216 only after its threshold/exemption policy is frozen. Issue #218 package/API
reduction follows the high-churn #219 work.

## Focused tasks for Claude

These are deliberately small reviews, tests, or documentation batches. They must not activate a
service, use credentials, run a guest, or turn a review into an architecture decision.

| Task | Packet and execution status | Deliverable |
| --- | --- | --- |
| CL1 — ADR-0029 S0 decision review | `PASSED` | The retained review accepts the two Supervisor services/four calls plus ADR-0044's separate CLI call, authentication-before-copy, no opcode bus, Go authority ownership, and `AttemptID`-only recovery; it reconciles S3/C4 ordering, refusal ownership, transport candidates, and the ADR-0040 fixed-store exception. |
| CL2 — Issue #216 ratchet packet | Packet `PASSED`; execution `BLOCKED` on assignment | Fixed complexity thresholds, narrow function/path-specific exemptions, owner/removal conditions, and preservation rules for intentionally linear protocol validators. |
| CL3 — Issue #217 package documentation batch | Packet `PASSED`; execution `BLOCKED` on assignment | One package at a time, starting with authority-bearing `registrationstate` exports. Document purpose, provenance/authority, caller obligations, and passive/product limitations without API changes. |
| CL4 — IPC refusal-matrix review | `PASSED`; historical disposition `AMEND`; follow-up closed | The independent read-only audit found no runtime authority bypass. This focused implementation closes its before/at/after 5,000-ms, complete-dictionary/map, every-field, required-`noState`, cancellation/deadline, refusal, replay, response-loss, and mutation-proof findings. CL4 remains historically `AMEND`; the exact passive C4 claim is now `PASSED`. |

## Deep research tasks for ChatGPT

Research must prefer primary sources, separate documented behavior from inference and observation,
and return a canonical repository artifact before its conclusion drives implementation.

| Task | Packet and execution status | Research question and stop boundary |
| --- | --- | --- |
| R1 — Native XPC enforcement brief | `PASSED` | Primary Apple sources and macOS 26.5 SDK declarations select the low-level `xpc_connection_t` controlled-harness baseline, requirement-before-activation, exact-message `SecCode` validation, connection-time EUID/ASID checks, non-preemptive cancellation, protocol-owned deadlines, and store-owned response-loss convergence. No portal, credential, service, or process mutation occurred. |
| R2 — Typed transport design | `PASSED` | The retained research reconciles the narrowed single-`main.mjs` caps and three role-distinct streams, freezes the passive state-machine input for C5a, separates frame observation from durable terminal truth, and supplies cancellation/reset/response-loss and restoration matrices. No guest or artifact mutation occurred. |
| R3 — Broker live-signing evidence brief | `PASSED` | The [canonical passive brief](BROKER_LIVE_SIGNING_EVIDENCE_BRIEF.md) maps fresh `LAContext`, Secure Enclave/Keychain candidates, AppKit focus/spoof/cancel, process death, replay, and update behavior from Apple sources while preserving Supervisor-owned durable approval authority. No live key, credential, prompt, service, or installed evidence. |
| R4 — F6 schema and VFS fault map | Packet `PASSED`; execution `BLOCKED` as external-alpha work | External-alpha only. Map the split active/immutable schema, fixed statements, transactions, SQLite result codes, and VFS injections to every F2-F5 known answer without generic SQL, DSN, pool, retry, or extension authority. |

## Explicitly deferred work

- F6 execution needs a separately authorized `capsule-experiments` task naming an owned disposable
  Apple-silicon host, APFS root/volume, and interruption owner. The experiment selects nothing.
  Any later engine/configuration selection or rejection requires a new Proposed ADR.
- Source Validator R4/R5 remains post-alpha defense-in-depth. Exact R4-v1 candidates are `NO_GO`;
  R4-v2 is unexecuted and `BLOCKED`.
- TypeScript Source Preparer, automatic TUF updates, Developer ID/notarized distribution,
  clean-host/minimum-OS coverage, independent-builder equality, restore activation, continuity,
  and external-alpha distribution are not on the owner-only internal-alpha critical path.

## Next checkpoint

The next orchestration checkpoint should collect this focused C4 implementation closure. C4, C5a,
CL1, CL4, R1, R2, and R3 are `PASSED` in their exact scopes. CL4 retains historical disposition
`AMEND`, now closed by this follow-up.
The owner may authorize C1 and C3 independently when their exact workspace, identity, and host
boundaries are available. C2's CL1/R1 prerequisites are closed, but its exact fixture/host
authorization remains mandatory. The disposable C6b1 Broker signing harness remains independently
`BLOCKED` on its exact host/identity/profile/group/fixture/mutation authorization. After those
handoffs, update this file and the evidence ledger before starting installed composition. No task
may promote internal alpha or product admission from research, a passive contract, an ad-hoc
harness, or a fixed guest experiment alone.
