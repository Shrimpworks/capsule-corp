# Current work plan

Date: 2026-08-11

Work item: close C2a's strict S3 deadline-equality fixtures without activating native IPC.

Status: `PASSED` for the exact passive C2a fixture and documentation scope.

Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`.

Product admission and the installed security boundary: `BLOCKED`.

This is the current execution index. Detailed security claims remain in the linked ADRs, readiness
map, experiment checkpoints, and evidence ledger. A completed passive contract or controlled
experiment is not an activated product path.

## Reconciled baseline

This slice starts from fetched `origin/main` commit
`88f3a2c1f968b1aa604ce14a2db4389822e5b193`, merge PR #251. PR #249 closed the exact passive C4
boundary and independent-verifier findings. PR #250 then made durable evidence retention a tracked
steering and archive invariant: important raw evidence cannot live only in a temporary workspace,
and cleanup cannot precede verified immutable remote publication. PR #251 reconciled the five
preparation callbacks and construction order. None activates a runtime, governed profile,
installed identity, service, product consumer, backend, VM, or guest.

Five later read-only preparation tasks now close what can be known before construction or
authorization:

- C1 recovery discovery `PASSED`, but raw v10-v27 publication remains `BLOCKED`: the former local
  archive commit `3fdcf2cebda087ecc99fbc73acfd21a3eae06b5b`, branch, and workspace are absent from
  bounded Capsule paths, Git refs/worktrees/objects, and the remote archive;
- C2 authorization/fixture preparation and C2a's nine strict S3 deadline-boundary cases are
  `PASSED`; the native harness remains `BLOCKED` on exact host/session authorization;
- C3 packet audit completed, but E1 remains `BLOCKED`: E0 is a closed inert specification, not a
  materialized independently hashed probe/bundle/verifier set;
- C5b authorization/fixture preparation `PASSED`; the owner selected the narrower v19 lineage plus
  the governed 103-byte source/SourceManifest pair for a no-run successor, which does not admit it;
  and
- C6b1 preparation `PASSED`; execution remains `BLOCKED` until an unsigned Broker harness,
  deterministic fixture manifest, test-only Supervisor seam, matching profile, disposable account,
  and separately enumerated mutation authorization exist.

The current retained baseline includes these exact slices:

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
raw v10-v27 recovery -> verified publication ------------------┐
                                                               │
S3 deadline closure -> native XPC harness ----------------------+-> installed authenticated IPC
                                                               │
E0 materialization -> ADR-0045 E1 identity separation ----------+-> key/service/root corpus
                                                               │
C4 passive approval/attempt evidence (PASSED) + R3 passive research (PASSED)
          -> unsigned Broker harness -> test Supervisor seam -> installed signing harness -----┐
installed authenticated IPC boundary (BLOCKED) ----------------------------------------------+-> product Broker/approval/attempt wiring
                                                                                              -> protected one-attempt path
                                                               │
typed transport design -> passive contract (`PASSED`) -> v19/103-byte no-run successor
                                                   -> controlled harness ----┐
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
| C1a | Recover the unpublished raw v10-v27 archive | `BLOCKED` | Read-only recovery discovery `PASSED`, but the former `/tmp/capsule-owned-guest.njcPIL/capsule-experiments` workspace and local commit `3fdcf2c...` are unavailable in every authorized location checked. Resume only from an owner-supplied backup, clone, bundle, object database, or filesystem snapshot. Do not reconstruct the 279-entry manifest from chat history. |
| C1b | Publish recovered raw v10-v27 evidence | `BLOCKED` | Requires C1a. Verify the 280-file packet and 279 retained-file manifest, publish `experiments/gate-c-c2b-fixed-owned-guest` in one reviewed immutable `capsule-experiments` commit, read it back, rerun its verifier, then add exact links. If recovery is impossible, a separately authorized rerun must produce new evidence. |
| C2a | Freeze S3 deadline equality boundaries | `PASSED` | Ordered before/at/after cases now cover SubmitMain 9,999/10,000/10,001 ms, RegisterPlan 4,999/5,000/5,001 ms, and GetRegisteredPlan 1,999/2,000/2,001 ms with equality expired, complete zero-state projections, and independent Go/Node checks. The unified passive contract has 15 deadline cases; no listener or process exists. |
| C2b | Run the one-time native XPC S3 harness | `BLOCKED` | Requires explicit authorization naming `Shrimpworks/capsule-experiments`, the owner-confirmed Mac/session, Capsule commit, manifest `c76e1f6c...8b59`, native contract `7ae502b0...962c`, ordered case digest `9ac6845b...f68e`, experimental service alias map, disposable names/root, and defensive no-product scope. Retain OS peer refusal, exact-message identity, EUID/session, copy, cap/flow, deadline, interruption, response-loss, process-fault, and cleanup evidence. |
| C3a | Materialize deterministic E0 fixtures | `BLOCKED` | Turn the `PASSED` inert E0 specification into exact current/legacy probes, a no-launch Coordinator bundle, plists, entitlements, LaunchAgent/descriptor inputs, manifest, and independent verifier with immutable digests. No portal, signing, profile, container, Keychain, service, runtime, or guest mutation. |
| C3b | Run ADR-0045 E1 identity separation | `BLOCKED` | Requires C3a, owner-confirmed host/device binding, one frozen legacy negative profile, evidence workspace, and exact Apple Development authorization. Prove current/legacy cross-container denial and cleanup; stop before Keychain, service registration, protected root, store, runtime, or guest. ADR-0045 remains Proposed. |
| C4 | Freeze `SubmitApprovalV0` and `RequestAttemptV0` | `PASSED` | PR #248 is the canonical predecessor and PR #249 closes the focused follow-up. Ordered 4,999/5,000/5,001-ms cases for both methods, equality-as-expired behavior, complete closed dictionaries/maps, every ordered field, required `noState`, cancellation/deadline commit truth, all 20 foreign-tag collisions, complete refusal and five-entry response-loss tables, and bounded Go/Node mutation proofs are retained. No listener, signer, store consumer, process, or guest is active. |
| C5a | Freeze the final typed source/input/completion transport | `PASSED` | The passive v1 contract freezes exact 152/160/64-byte layouts, 262,144-byte payload caps, completion cap-plus-one, big-endian bindings, four statuses, canonical JSON, refusal precedence, monotonic state/fault behavior, endpoint custody, completion-last projection, deterministic fixtures, and independent Go/Node verification. No endpoint, process, runtime, backend, guest, or store mutation occurred. |
| C5b0 | Materialize the no-run typed-transport successor | `BLOCKED` | The owner selected the v19 benign lineage at composed digest `ac272171...f48fa` plus the governed 103-byte source `c8e940...b475` and SourceManifest `712b1b...61b0`. Produce exact successor profile, runner, root, init, launcher, controller, plan, and frame identities without loading libkrun, calling HVF, starting a guest, or making an admission decision. |
| C5b | Run the controlled typed-transport harness | `BLOCKED` | Requires C5b0 and separate authorization naming the exact successor, owner-confirmed host, and owned disposable guest. Retain directional copy, chunk/cap+1, stall/reset/cancel, descriptor substitution, response-loss, completion-last, teardown, and restoration evidence without making an admission decision. |
| C6a | Build the installed authenticated service and protected-state boundary | `BLOCKED` | Requires passed C2b and C3b evidence under Accepted ADR-0029, then separate authorization for the Keychain/service/protected-root corpus. Implement only method-specific listeners, peer authentication, owner/store startup, and the four passively frozen Supervisor consumers. |
| C6b1a | Build the unsigned Broker evidence harness | `BLOCKED` | Add and review the disposable target/source, deterministic projection/payload/display fixture manifest, effective-entitlement request, public-key-only checks, and no-credential test doubles. No signing, Keychain, LocalAuthentication, provisioning, install, or product consumer. |
| C6b1b | Build the test-only Supervisor evidence seam | `BLOCKED` | Add an explicitly non-product harness that exercises SubmitApproval/RequestAttempt commit, replay, response-loss, and crash convergence while preserving the Supervisor as the only durable authority owner. It must not activate an installed listener or product consumer. |
| C6b1c | Provision and read back the disposable identity | `BLOCKED` | Requires C6b1a/b plus exact owner authorization naming the owned Mac, development target/App ID/profile/access group, disposable account/container, and immutable composite fixture manifest. Proposed host label `dsteele-shrimp-mbp18-4-01` remains subject to owner confirmation. Stop before key creation or prompt/signing. |
| C6b1d | Run the installed Broker signing evidence matrix | `BLOCKED` | Requires C6b1c and separately enumerated mutation authorization. Proposed first-run destructive rows D1-D4 and cleanup D14-D16 are not yet authorized; D5-D13 and D17-D18 remain deferred. No product consumer, runtime, backend, VM, or guest. |
| C6b2 | Connect the product Broker and approval/attempt methods | `BLOCKED` | C4 is `PASSED`; C6a and C6b1d remain required. Implement native rendering/UI, installed signing/public-key verification, and method-specific `SubmitApprovalV0`/`RequestAttemptV0` consumers without runtime or guest activation. Research and passive conformance cannot satisfy either installed dependency. |
| C6c | Wire attempt admission and the fixed-store stop policy | `BLOCKED` | Requires C6b2 and an explicit decision for p95 provenance/window/lifetime and any persistent timing-trip semantics. Apply the re-evaluated guard transaction-locally after owner/full verification and before a new consume/create mutation; replay of an existing `AttemptID` converges first. |
| C6d | Run the pre-admission installed runtime/profile matrix | `BLOCKED` | Requires C1b, C5b, C6a, governed release/artifact review, and separate authorization naming the exact signed-installed candidate and owned test environments. Retain identity compatibility, runtime/root/loader restoration, transport, teardown/recovery, and the required broader lifecycle/platform evidence without accepting user source or making an admission decision. |
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

C2a is now `PASSED`; C2b must stop for exact owner/environment authorization. C3a, C5b0, C6b1a,
and the independent test-only C6b1b seam may continue as separate retained construction tasks. C1a
resumes only if the owner supplies a backup or snapshot; otherwise a later exact rerun needs
separate authorization and new identities.

No signing, Keychain, service registration, installed listener, portal/profile mutation, libkrun/HVF
execution, VM, or guest is authorized by this plan. C2b, C3b, C5b, C6b1c, and C6b1d must each stop
again for their exact owner/environment authorization. No task may promote internal alpha or
product admission from research, a passive contract, an ad-hoc harness, or a fixed guest
experiment alone.
