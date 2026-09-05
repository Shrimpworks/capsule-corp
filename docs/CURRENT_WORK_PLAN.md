# Current work plan

Date: 2026-09-05

Work item: retain the first seven native C5b transport providers and their local pipe evidence,
preserving the blocked complete-provider, execution, and product boundaries.

Status: `PASSED` for this canonical reconciliation and the completed child scopes named below.

Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`.

Product admission and the installed security boundary: `BLOCKED`.

This is the current execution index. Detailed security claims remain in the linked ADRs, readiness
map, experiment checkpoints, and evidence ledger. A completed passive contract or controlled
experiment is not an activated product path.

The preceding control reconciliation addressed
[issue #306](https://github.com/Shrimpworks/capsule-corp/issues/306) and
[issue #308](https://github.com/Shrimpworks/capsule-corp/issues/308). It updated `RUNTIME-001`,
`SUPPLY-001`, and `SOURCE-000` in place and recorded the implemented passive
proposal/plan/custody symbols. The preceding follow-up pinned the reviewed C5b8/C5b9 no-run
experiment merges. This reconciliation pins the later build-only compatibility preflight: its
static question `PASSED`, the exact provider-only binding candidate is `NO_GO`, and controlled C5b
execution remains `BLOCKED`. The later C5b11 fixed-runner construction is now reconciled below;
provider implementation/provenance, reviewed runnable composition, and final exact authorization
remain required. Authenticated consumers, protected installed state,
runtime/profile admission, and product activation remain `BLOCKED`. No reconciliation changes an
ADR lifecycle or product-admission result.

## 2026-09-05 C5b11 reconciliation

Status: `PASSED` for the [retained-successor checkpoint](C5B_FIXED_RUNNER_SUCCESSOR_CHECKPOINT.md).

Experiments PR #30 (C5b10) and PR #31 (C5b11) merged on August 18 but were omitted from the
August 21 planning sync. C5b10 is preserved history, not accepted evidence. C5b11 merge
`f206e4ef2cd326ee74e5b7b2739c62efe6da7d6d` retains the reviewed-head descendant
`27b9011fb80edc0f23c11b3f3fa76d00cebc2365`. Fresh local static, mutation, and unlinked-object
reproduction checks pass. PR #31 reports the final exact-head review as `PASSED` / `Ready`;
archived pre-review prose remains unchanged and is distinguished in the checkpoint.

The first implementation slice is now the seven native transport providers in the
[C5b12 checkpoint](C5B_NATIVE_TRANSPORT_PROVIDER_CHECKPOINT.md). The fixed-runner no-run candidate
itself no longer needs construction. Seventeen lifecycle/root/durable-store providers remain
unimplemented; complete composition, controlled execution, installed composition,
runtime/profile admission, and product admission remain `BLOCKED`.

## 2026-09-05 C5b12 native provider milestone

Status: local transport implementation/checks and independent review `PASSED`.
Parent complete-provider composition: `BLOCKED`.

Seven actual Darwin pipe providers now cover endpoints, readiness, fixed source/input writes,
writer closure, start, and exact known-answer completion. The retained source has no caller
replacement bytes, paths, descriptors, or backend callbacks. Local tests link only transport code;
no runner, driver, libkrun/HVF, VM, guest, signing, installed state, or protected store executes.

Milestone order:

1. Native transport mechanics: seven providers and local adversarial pipe evidence `PASSED`.
2. Supervisor runner lifecycle and root custody: nine providers `BLOCKED` on implementation,
   exact identity/FD transfer, terminal/absence evidence, and teardown/root reconciliation.
3. Supervisor durable attempt/completion store: eight providers `BLOCKED` on implementation and
   independently persisted fault/reopen cursors, one-shot teardown, and stored replay evidence.
4. Complete immutable composition: `BLOCKED` on all owners, exact timing/identity/provenance and
   independent review. Only then can a separately authorized controlled typed-transport guest
   attempt be considered. That result would still not establish product admission.

## 2026-08-21 planning sync

Work item: current documentation and executable work order.

Status: `PASSED` for this documentation/planning reconciliation.

Scope: compare the canonical plan, architecture, status dashboard, latest merged history, and open
repository work. The planning sync is accepted on GitHub `main` at exact commit
`d54bd5a351ccfb2abe9956f07acbfdd6e96d07a5`; issue #314 is the first focused hardening slice after
that baseline and changes no product mechanism or retained experiment evidence.

Evidence or reason: the Source Validator's R1-R3 scopes and the `registrationstate`
exported-contract batch are already complete. Issue
[#314](https://github.com/Shrimpworks/capsule-corp/issues/314) now has explicit `lstat` refusal and
repository regression tests for both included archive roots. Issues
[#315](https://github.com/Shrimpworks/capsule-corp/issues/315) through
[#321](https://github.com/Shrimpworks/capsule-corp/issues/321) retain one contradictory test
comment, one narrowed high-segment test path, and five bounded duplication/test-discovery cleanups.
None changes an ADR lifecycle, control-evidence row, runtime/profile status, or product admission.

Remaining work at this historical sync: construct and independently review the C5b no-run
fixed-runner successor (superseded by the September 5 C5b11 checkpoint above); obtain separate exact
authorization before any native XPC, container,
live-signing, libkrun/HVF, VM, or guest execution.

### 2026-08-29 Q1-Q8 closure

Issues #314-#321 (the full Q1-Q8 immediate repository-quality order below) are now `PASSED` and
closed. #320 merged directly to `main` in [PR #329](https://github.com/Shrimpworks/capsule-corp/pull/329).
#315-#319 and #321 merged as a stacked branch sequence ([PR #330](https://github.com/Shrimpworks/capsule-corp/pull/330)-[PR #335](https://github.com/Shrimpworks/capsule-corp/pull/335)) and landed on `main`
via the consolidating [PR #336](https://github.com/Shrimpworks/capsule-corp/pull/336). Q4/#316's test-policy
decision resolved to the `testing.Short()` gate: `main`'s `go` CI job now runs `go test -short`, and
a new `go-full-suite` job runs the complete (non-`-short`) suite nightly plus on `workflow_dispatch`.
None of the eight changes an ADR lifecycle, control-evidence row, runtime/profile status, or product
admission. The C5b next-action wording at this checkpoint is superseded by the September 5
C5b11 reconciliation above.

Parent status: owner-only hostile-`.mjs` internal alpha remains
`IN_PROGRESS — TRENDING_GOOD`; product admission and the installed security boundary remain
`BLOCKED`.

### Immediate repository-quality order

| Order | Work item | Status | Acceptance and verification | Blocker and next action |
| --- | --- | --- | --- | --- |
| Q1 | [#314 archive-root symlink refusal](https://github.com/Shrimpworks/capsule-corp/issues/314) | `PASSED` | Both included roots are `lstat`-checked as real directories before traversal; root-symlink mutations fail closed; the focused verifier test, `pnpm test`, and `pnpm lint` pass without changing archive identities or counts. | Complete in the focused #314 hardening slice; no product admission or concurrent-mutation claim. |
| Q2 | [#320 script test discovery and shared SHA-256 helper](https://github.com/Shrimpworks/capsule-corp/issues/320) | `PASSED` | New `scripts/*.test.mjs` files cannot be silently omitted (`package.json`'s test script now globs); verifier and test reuse the shared `sha256Hex` byte helper; discovered script tests and ordinary pnpm gates pass. | Closed in [PR #329](https://github.com/Shrimpworks/capsule-corp/pull/329); no product admission or concurrent-mutation claim. |
| Q3 | [#315 contradictory lifecycle test comment](https://github.com/Shrimpworks/capsule-corp/issues/315) | `PASSED` | The comment names the existing `duplicate-instance-on-applied-quarantines-without-adoption` subtest and no nonexistent test or open defect; the focused `registrationstate` test passes. | Closed in [PR #330](https://github.com/Shrimpworks/capsule-corp/pull/330); documentation-only. |
| Q4 | [#316 exact-64-segment pipeline coverage](https://github.com/Shrimpworks/capsule-corp/issues/316) | `PASSED` | Restored the disk-backed full-pipeline test gated behind `testing.Short()`; `main`'s `go` CI job runs `go test -short`, a new `go-full-suite` job runs the complete suite nightly plus on `workflow_dispatch`. | Closed in [PR #331](https://github.com/Shrimpworks/capsule-corp/pull/331); resolves the test-policy decision as the gated full-path corpus option. |
| Q5 | [#317 completion result-cap duplication](https://github.com/Shrimpworks/capsule-corp/issues/317) | `PASSED` | One cap calculation/classification retained (`CommitCompletion`'s pre-encode check); all completion-store refusal/replay tests pass unchanged. | Closed in [PR #332](https://github.com/Shrimpworks/capsule-corp/pull/332); behavior-preserving refactor. |
| Q6 | [#318 reconciliation quarantine duplication](https://github.com/Shrimpworks/capsule-corp/issues/318) | `PASSED` | Quarantine field-setting shared via `applyIdentityQuarantine`/`markIdentityQuarantineTrust`; trust reason, reconciliation, recovery, and durable bytes unchanged; focused lifecycle fault/reopen tests pass. | Closed in [PR #333](https://github.com/Shrimpworks/capsule-corp/pull/333). |
| Q7 | [#319 archive marshal/digest duplication](https://github.com/Shrimpworks/capsule-corp/issues/319) | `PASSED` | Centralized via `marshalAndDigestRecord`, preserving exact bytes, digest classifications, and known answers across all 8 call sites. | Closed in [PR #334](https://github.com/Shrimpworks/capsule-corp/pull/334); left issue #219's larger staged refactor unaffected. |
| Q8 | [#321 test-only CBOR helper duplication](https://github.com/Shrimpworks/capsule-corp/issues/321) | `PASSED` | Shared test scan/head encoding helpers (`topLevelFieldRange`, `encodeCBORHead`) with no fixture, refusal, or known-answer change; focused `v0candidate` mutation tests pass. | Closed in [PR #335](https://github.com/Shrimpworks/capsule-corp/pull/335); production decoding untouched. |

## Reconciled baseline

The retained execution baseline begins at fetched `origin/main` commit
`bd7cc9c98c07c91b4d96d3efa2f6261aba350971`, merge PR #256. PR #255 reconciled PR #254 at
`e5401a81b727915ec01afe9012a77e7586a57c13` with the independently completed C3b
profile/signature evidence without changing PR #254's historical input. This checkpoint adds the
later C5b2 governed-input closure at experiments PR #18 and the five no-run C5b construction
results at experiments PRs #19-#23. PR #252 closes C2a's passive S3 deadline-equality fixture
prerequisite. Four separately reviewed `capsule-experiments` merges retain the earlier construction
inputs:

- [PR #6](https://github.com/Shrimpworks/capsule-experiments/tree/b357d0c0fb29100c180494e67cebd7809aabe3c5/experiments/typed-guest-transport-c5b0-v19-successor)
  at merge `b357d0c0fb29100c180494e67cebd7809aabe3c5` retains the C5b0 deterministic no-run packet;
- [PR #7](https://github.com/Shrimpworks/capsule-experiments/tree/dee784d40684100f8315720fab9a5cd3399f492b/experiments/macos-installation-i2b3-supervisor-authority-epoch-e0)
  at merge `dee784d40684100f8315720fab9a5cd3399f492b` retains the C3a reproducible unsigned E0 packet;
- [PR #8](https://github.com/Shrimpworks/capsule-experiments/tree/4a2447d4bd0e03132dc616e608031ca313630cdd/experiments/broker-live-signing-c6b1)
  at merge `4a2447d4bd0e03132dc616e608031ca313630cdd` retains the C6b1a unsigned Broker harness; and
- [PR #9](https://github.com/Shrimpworks/capsule-experiments/tree/067fe2beb40361bb714507cab1331004e0a656fa/experiments/broker-live-signing-c6b1-supervisor-seam)
  at merge `067fe2beb40361bb714507cab1331004e0a656fa` retains the C6b1b test-only Supervisor seam.

Nine later `capsule-experiments` merges retain the next construction and platform-preflight
wave:

- [PR #10](https://github.com/Shrimpworks/capsule-experiments/tree/db08ebf277432e06d6cba3b7f7338e3bd4a61252/experiments/typed-guest-transport-c5b1-executable-successor)
  at merge `db08ebf277432e06d6cba3b7f7338e3bd4a61252` retains the fresh deterministic C5b1
  executable-successor candidates;
- [PR #11](https://github.com/Shrimpworks/capsule-experiments/tree/3d7bd46352506bf6018286749c2c85a3e2f683df/experiments/authenticated-local-ipc-s3-native-xpc-c2b0-ce7220e523bc43ba-c7ae502b0742bab1e)
  at merge `3d7bd46352506bf6018286749c2c85a3e2f683df` retains the inert C2b0 native-XPC harness;
- [PR #12](https://github.com/Shrimpworks/capsule-experiments/tree/50c494d4841c5d42e8e2120b82c0481a706a5236/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1)
  and [PR #13](https://github.com/Shrimpworks/capsule-experiments/tree/cd06bd84690a16bb4d0924a8a4cd64845ebb0159/experiments/broker-live-signing-c6b1c-identity-profile-readback)
  retain the exact zero-effect C3b missing-profile stop and C6b1c portal/download stop;
- [PR #14](https://github.com/Shrimpworks/capsule-experiments/tree/e6390253a274e9ead76366f9869a5e1b272a1595/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1-app-group-preflight)
  retains the C3b portal-form observation, while
  [PR #16](https://github.com/Shrimpworks/capsule-experiments/tree/3671a6eb23357ff28de4562dd60e8f68173034ae/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1-app-group-preflight)
  supplies its canonical interpretation: only Developer-portal registration of the frozen
  macOS-style App Group is `NO_GO`; at that checkpoint the exact identity remained intended and
  `BLOCKED` on signed platform evidence; and
- [PR #15](https://github.com/Shrimpworks/capsule-experiments/tree/82d1a799f70482856aaa6030f612d701b39cec67/experiments/broker-live-signing-c6b1c-signed-artifact-readback)
  retains the `PASSED` no-install C6b1c identity/profile and signed-artifact readback; and
- [PR #17](https://github.com/Shrimpworks/capsule-experiments/tree/ee00ae2abbce64ae6458b82d0b53d904ee39aeb6/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1-signed-profile-preflight)
  at merge `ee00ae2abbce64ae6458b82d0b53d904ee39aeb6` retains the `PASSED` C3b
  profile/signature-only gate over the exact current Supervisor, never-launched Coordinator, and
  legacy negative probe; and
- [PR #18](https://github.com/Shrimpworks/capsule-experiments/tree/5a2f835e8c9df8279237f940f5af757e119593bd/experiments/typed-guest-transport-c5b2-governed-input-closure)
  at merge `5a2f835e8c9df8279237f940f5af757e119593bd` retains the `PASSED` C5b2
  governed-input closure over the available current-source libkrun/header/ABI/final-runner bytes.

Five subsequent `capsule-experiments` merges retain the bounded C5b no-run input and controller
closure:

- [PR #19](https://github.com/Shrimpworks/capsule-experiments/tree/b5db7bcbbf7fe814faec3950ebfbf2d2ac7786e2/experiments/typed-guest-transport-c5b3-runtime-input-recovery)
  at merge `b5db7bcbbf7fe814faec3950ebfbf2d2ac7786e2` retains exact recovered `rusty_v8`
  archive/binding custody plus the bounded runtime/libkrunfw reconstruction plan;
- [PR #20](https://github.com/Shrimpworks/capsule-experiments/tree/60234e22674e46a42e8e5c382d85217a930c2c13/experiments/typed-guest-transport-c5b3-controlled-test-controller)
  at merge `60234e22674e46a42e8e5c382d85217a930c2c13` retains the pure C17 controller core,
  byte-equal non-executable objects, and closed state/fault/replay/cleanup vectors;
- [PR #21](https://github.com/Shrimpworks/capsule-experiments/tree/068e221dafa7cf3e9a945cee7e8bf077eeed1c6b/experiments/typed-guest-transport-c5b4-libkrunfw-recovery)
  at merge `068e221dafa7cf3e9a945cee7e8bf077eeed1c6b` retains two exact network-denied
  reproductions of `libkrunfw.5.dylib` and the official generated-source input;
- [PR #22](https://github.com/Shrimpworks/capsule-experiments/tree/3cfe7db16c55894be444d4c783659043dbd25c95/experiments/typed-guest-transport-c5b5-no-run-effect-adapter)
  at merge `3cfe7db16c55894be444d4c783659043dbd25c95` retains the compile-only descriptive
  action adapter and exact static libkrun symbol closure; and
- [PR #23](https://github.com/Shrimpworks/capsule-experiments/tree/d9967e80a6155a65c6876dc686d8f8498b4a908f/experiments/typed-guest-transport-c5b6-deno-static-reproduction)
  at merge `d9967e80a6155a65c6876dc686d8f8498b4a908f` retains two independent exact Cargo
  acquisitions and byte-identical network-disabled static reproductions of the governed fixed-
  fixture Deno runtime, snapshot, and bundle.

This C5b7 follow-up starts from fetched `origin/main` commit
`b2ab848a4551fc29e845cbc5178bb312de3da1cb`. The later
[experiments PR #24](https://github.com/Shrimpworks/capsule-experiments/tree/78485fb91a31733c568fe43e5fa295474e5956e1/experiments/typed-guest-transport-c5b7-deterministic-runtime-root)
at merge `78485fb91a31733c568fe43e5fa295474e5956e1` retains two byte-identical
100,663,296-byte roots at SHA-256
`5ad18f20cbc97c7a70ead3e795fd3649672513323041e913b0eb55b7acc88775`, a closed
19-node inventory, independent raw-filesystem verification, and 15 mutation refusals.

This C5b8/C5b9 reconciliation starts from fetched `origin/main` commit
`6f246dcf4f5244ebbc319f948d39eac2818f9650`. Three reviewed experiment merges retain the
completed no-run successors:

- [PR #26](https://github.com/Shrimpworks/capsule-experiments/tree/e83614af34d5c39c12a4a3d6e6cda8dcf0304030/experiments/typed-guest-transport-c5b8-controlled-test-effects)
  at merge `e83614af34d5c39c12a4a3d6e6cda8dcf0304030` retains the sealed C5b8
  controlled-test operation layer and its test-double-only fault/replay/cleanup evidence;
- [PR #27](https://github.com/Shrimpworks/capsule-experiments/tree/b0819d76883eb86cbbc03b2b7033fe55bedbf713/experiments/typed-guest-transport-c5b8-c5b7-root-binding-successor)
  at merge `b0819d76883eb86cbbc03b2b7033fe55bedbf713` binds that layer to the exact
  100,663,296-byte C5b7 root without loading a retained runtime artifact; and
- [PR #28](https://github.com/Shrimpworks/capsule-experiments/tree/3965e6b5cc87d476da7f431d7ed8a5758011a1b8/experiments/typed-guest-transport-c5b9-immutable-no-run-composite)
  at merge `3965e6b5cc87d476da7f431d7ed8a5758011a1b8` retains the C5b9 immutable
  no-run composite over six exact component roles, closed static ABI/load and archive inventory,
  typed fixtures, nine unit tests, and 14 mutation refusals.

The C5b9 packet deliberately leaves `_c5b8_controlled_test_operation` without a provider and
records host, guest, authorization, and every effect as absent. It does not recover or recreate the
lost v19/v27 bytes, authorize controlled execution, or advance runtime/profile or product
admission.

This C5b controlled-harness preflight reconciliation starts from fetched `origin/main` commit
`ac38be1d73ef31cd2e84873f003c005ba60f5afc`. The later
[experiments PR #29](https://github.com/Shrimpworks/capsule-experiments/tree/7fc3af9c46895b340c3118a96cb50abb26b1d977/experiments/typed-guest-transport-c5b-controlled-harness-preflight)
at merge `7fc3af9c46895b340c3118a96cb50abb26b1d977` retains the `PASSED` build-only
compatibility preflight and ten mutation refusals. It makes only the exact candidate—turn C5b9
runnable by supplying `_c5b8_controlled_test_operation` without changing the retained runner,
root, or effect ordering—`NO_GO`. The retained runner and root disagree on byte identity; the
effect plan reaches `krun_start_enter` before later source/input writes; the standalone runner does
not implement the per-effect operation ABI; and both candidate owners import libkrun execution.
No artifact was loaded or executed. The retained host/disposable-guest confirmation authorized
preparation only; the preflight records `executionAuthorized: false` and requires later final
manifest authorization.

This architecture-gate reconciliation starts from fetched `origin/main` commit
`686942a34b36a12d353224859fcf835fc916d048`. It preserves every retained experiment and adds
only dependency and lifecycle clarity: C3b evidence must feed a separate C3c adoption decision
before C6a, and Source Validator R4/R5 remains a gate for admitting that later product control
rather than the owner-only internal-alpha runtime/profile candidate reviewed at C7.

This control/documentation reconciliation starts from fetched `origin/main` commit
`7576dfb69d75f6ef11ed9708c2ab407aec06e9be`, merge PR #310. It changes no mechanism or evidence
state. It makes the existing fork-native C1 checkpoint canonical inside `RUNTIME-001` and
`SUPPLY-001`, and aligns Architecture, Technical Design, the readiness dashboard, and `SOURCE-000`
with the already-passed passive proposal, plan, and custody mechanics.

These are immutable evidence pins, not product dependencies. None activates a runtime, governed
profile, installed identity, service, product consumer, backend, VM, or guest.

The following bounded preparation and construction tasks close what can be known before controlled
platform execution or product authorization:

- C1 recovery discovery `PASSED`, but raw v10-v27 publication remains `BLOCKED`: the former local
  archive commit `3fdcf2cebda087ecc99fbc73acfd21a3eae06b5b`, branch, and workspace are absent from
  bounded Capsule paths, Git refs/worktrees/objects, and the remote archive;
- C2 authorization/fixture preparation, C2a's nine strict S3 deadline-boundary cases, and C2b0's
  inert reproducible harness construction are `PASSED`; native execution remains `BLOCKED` on a
  fresh exact host/session authorization;
- C3a deterministic E0 materialization and the bounded C3b preflights are `PASSED` in their exact
  no-launch scopes. The required legacy profile has been restored. The portal-registration path
  for the frozen macOS-style App Group is `NO_GO`, not the identity itself. The exact
  profile/signature-only gate is also `PASSED`; C3b/E1 remains `BLOCKED` on a fresh
  launch/container authorization;
- C5b0-C5b9 no-run packet, executable/input, controller/effect, root, immutable-composite
  construction, and the later compatibility preflight are `PASSED` in their exact bounded scopes.
  Exact `rusty_v8`, governed fixed-
  fixture Deno runtime, libkrunfw, controller-core, descriptive adapter, size-compatible C5b8
  successor, runtime-root, and C5b9 composite identities are retained. C5b8 resolves C5b7's exact
  root-size binding only for the reviewed test-double path. Direct provider-only composition of
  the retained C5b9 inputs is `NO_GO`. The later C5b11 successor construction is `PASSED` in
  its no-run scope. C5b12 adds seven native transport providers with local pipe evidence;
  the remaining 17 providers, complete composition, controlled C5b run, preferred-form
  libkrunfw/kernel source compliance, and admission remain `BLOCKED`; and
- C6b1a unsigned Broker-harness construction, C6b1b test-only Supervisor-seam construction, and
  C6b1c no-install identity/profile/signed-artifact readback are `PASSED`; C6b1d installed live
  signing remains `BLOCKED` on its own exact Keychain/LocalAuthentication authorization. Its
  Capsule commit `16fb810b...` remains the immutable C6b1c construction input rather than a stale
  claim about the current documentation baseline.

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
S3 deadline closure -> inert native XPC harness (`PASSED`) -> controlled run +-> installed authenticated IPC
                                                               │
E0 materialization (`PASSED`) -> signed-profile gate (`PASSED`)
                              -> ADR-0045 E1 identity separation
                              -> ADR-0038/0045 adoption decision +-> key/service/root corpus
                                                               │
C4 passive approval/attempt evidence (PASSED) + R3 passive research (PASSED)
          -> unsigned Broker harness (`PASSED`) -> test Supervisor seam (`PASSED`)
          -> identity/profile readback (`PASSED`) -> installed signing harness -----------------┐
installed authenticated IPC boundary (BLOCKED) ----------------------------------------------+-> product Broker/approval/attempt wiring
                                                                                              -> protected one-attempt path
                                                               │
typed transport design -> passive contract (`PASSED`) -> v19/103-byte no-run packet (`PASSED`)
                                                   -> fresh executable construction (`PASSED`)
                                                   -> root/effects/composite (`PASSED`)
                                                   -> direct provider binding (`NO_GO`)
                                                   -> C5b11 no-run successor (`PASSED`) -> providers + review -> controlled harness ----┐
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
The C7 admission review is limited to ADR-0040's owner-only internal-alpha runtime/profile.
Source Validator R4/R5 is not a C7 prerequisite; it remains the blocked post-alpha delivery path
for admitting the Source Validator itself as a product defense-in-depth control.

## Next work owned by Codex

These are implementation or integration tasks. Each retained task uses its own `codex/<topic>`
branch and pull request unless the orchestrator explicitly groups it before work begins.

| Order | Work item | Current status | Start condition and acceptance boundary |
| --- | --- | --- | --- |
| C1a | Recover the unpublished raw v10-v27 archive | `BLOCKED` | Read-only recovery discovery `PASSED`, but the former `/tmp/capsule-owned-guest.njcPIL/capsule-experiments` workspace and local commit `3fdcf2c...` are unavailable in every authorized location checked. Resume only from an owner-supplied backup, clone, bundle, object database, or filesystem snapshot. Do not reconstruct the 279-entry manifest from chat history. |
| C1b | Publish recovered raw v10-v27 evidence | `BLOCKED` | Requires C1a. Verify the 280-file packet and 279 retained-file manifest, publish `experiments/gate-c-c2b-fixed-owned-guest` in one reviewed immutable `capsule-experiments` commit, read it back, rerun its verifier, then add exact links. If recovery is impossible, a separately authorized rerun must produce new evidence. |
| C2a | Freeze S3 deadline equality boundaries | `PASSED` | Ordered before/at/after cases now cover SubmitMain 9,999/10,000/10,001 ms, RegisterPlan 4,999/5,000/5,001 ms, and GetRegisteredPlan 1,999/2,000/2,001 ms with equality expired, complete zero-state projections, and independent Go/Node checks. The unified passive contract has 15 deadline cases; no listener or process exists. |
| C2b | Run the one-time native XPC S3 harness | `BLOCKED` | Requires explicit authorization naming `Shrimpworks/capsule-experiments`, the owner-confirmed Mac/session, Capsule commit, manifest `c76e1f6c...8b59`, native contract `7ae502b0...962c`, ordered case digest `9ac6845b...f68e`, experimental service alias map, disposable names/root, and defensive no-product scope. Retain OS peer refusal, exact-message identity, EUID/session, copy, cap/flow, deadline, interruption, response-loss, process-fault, and cleanup evidence. |
| C3a | Materialize deterministic E0 fixtures | `PASSED` | Archive merge [`dee784d40684100f8315720fab9a5cd3399f492b`](https://github.com/Shrimpworks/capsule-experiments/tree/dee784d40684100f8315720fab9a5cd3399f492b/experiments/macos-installation-i2b3-supervisor-authority-epoch-e0) retains exact current/legacy probe sources and reproducible unsigned bundles, a never-launched Coordinator, plists, entitlement/profile requests, disabled LaunchAgent and inactive descriptor inputs, a closed manifest, independent verification, and 23 mutation refusals. No portal, identity, profile, signing, container, service, Keychain, runtime, backend, VM, or guest was accessed or activated. |
| C3b | Run ADR-0045 E1 identity separation | `BLOCKED` | C3a, exact legacy-profile restoration, the App Group portal preflight, and the [exact profile/signature-only gate](https://github.com/Shrimpworks/capsule-experiments/tree/ee00ae2abbce64ae6458b82d0b53d904ee39aeb6/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1-signed-profile-preflight) are complete. That gate `PASSED` exact Team/application identifiers, profile UUID/CMS/certificate/device binding, CDHashes/designated requirements, effective App Group/Keychain entitlements, hardened runtime, and absent debug entitlement without launching a bundle or opening a container. The frozen `3DDR84M4JS...` App Group remains the macOS-style entitlement value and is not a Developer-portal App Group resource; do not rewrite it to `group.`. A fresh authorization may now run only E1-01..E1-12 and E1-14..E1-15 against the exact retained identities; E1-13 remains excluded. ADR-0045 remains Proposed. |
| C3c | Decide Supervisor bootstrap and authority-epoch adoption | `BLOCKED` | Requires passed C3b evidence. Reconcile the exact E1 result against ADR-0038 and ADR-0045, then accept, amend, reject, or supersede both Proposed decisions. Any accepted decision must freeze the complete authority descriptor, transition/state-engine binding, bootstrap ceremony, and create/open disposition consumed by C6a. If evidence is insufficient, C3c remains `BLOCKED` and names the additional bounded evidence; C6a cannot supply the missing decision after implementation starts. |
| C4 | Freeze `SubmitApprovalV0` and `RequestAttemptV0` | `PASSED` | PR #248 is the canonical predecessor and PR #249 closes the focused follow-up. Ordered 4,999/5,000/5,001-ms cases for both methods, equality-as-expired behavior, complete closed dictionaries/maps, every ordered field, required `noState`, cancellation/deadline commit truth, all 20 foreign-tag collisions, complete refusal and five-entry response-loss tables, and bounded Go/Node mutation proofs are retained. No listener, signer, store consumer, process, or guest is active. |
| C5a | Freeze the final typed source/input/completion transport | `PASSED` | The passive v1 contract freezes exact 152/160/64-byte layouts, 262,144-byte payload caps, completion cap-plus-one, big-endian bindings, four statuses, canonical JSON, refusal precedence, monotonic state/fault behavior, endpoint custody, completion-last projection, deterministic fixtures, and independent Go/Node verification. No endpoint, process, runtime, backend, guest, or store mutation occurred. |
| C5b0 | Materialize the deterministic no-run typed-transport packet | `PASSED` | Archive merge [`b357d0c0fb29100c180494e67cebd7809aabe3c5`](https://github.com/Shrimpworks/capsule-experiments/tree/b357d0c0fb29100c180494e67cebd7809aabe3c5/experiments/typed-guest-transport-c5b0-v19-successor) binds the v19 lineage digest, governed 103-byte source and SourceManifest, exact role contracts, no-run profile/plan, fresh typed frames, closed inventory, independent verifier, and six mutations. No v19 raw bytes were recreated; executable runner/root/init/launcher/controller identities remain explicitly null. |
| C5b1 | Construct the fresh executable typed-transport successor | `PASSED` | Archive merge [`db08ebf277432e06d6cba3b7f7338e3bd4a61252`](https://github.com/Shrimpworks/capsule-experiments/tree/db08ebf277432e06d6cba3b7f7338e3bd4a61252/experiments/typed-guest-transport-c5b1-executable-successor) retains fresh reproducible runner, raw root, trusted init, launcher, and hard-stop controller candidates, a closed 41-file inventory, provenance/SBOM, independent parsing, and seven mutation refusals. It does not recover v19 or bind/run the governed runtime, libkrun/libkrunfw, kernel, firmware, or a real controller. |
| C5b2 | Bind the available governed successor inputs | `PASSED` | Archive merge [`5a2f835e8c9df8279237f940f5af757e119593bd`](https://github.com/Shrimpworks/capsule-experiments/tree/5a2f835e8c9df8279237f940f5af757e119593bd/experiments/typed-guest-transport-c5b2-governed-input-closure) independently binds the exact current-source libkrun header, ABI audit, unsigned dylib, and final runner source/binary by hash, raw Mach-O metadata, dependencies, exports, imports, and C17 audit. In that immutable slice, libkrunfw/kernel receipts were identity evidence only and the governed `deno_core` executable/libkrunfw bytes and complete controller/composite were absent; C5b3-C5b7 later close only the no-run input/controller/runtime-root portions without rewriting C5b2. Separate firmware remains inapplicable under ADR-0041. Nothing was loaded or executed. |
| C5b3 | Recover runtime inputs and construct the controlled-test core | `PASSED` | Merges [`b5db7bcbbf7fe814faec3950ebfbf2d2ac7786e2`](https://github.com/Shrimpworks/capsule-experiments/tree/b5db7bcbbf7fe814faec3950ebfbf2d2ac7786e2/experiments/typed-guest-transport-c5b3-runtime-input-recovery) and [`60234e22674e46a42e8e5c382d85217a930c2c13`](https://github.com/Shrimpworks/capsule-experiments/tree/60234e22674e46a42e8e5c382d85217a930c2c13/experiments/typed-guest-transport-c5b3-controlled-test-controller) retain exact `rusty_v8` archive/binding custody and the pure no-effect C17 state machine. Its byte-equal `MH_OBJECT` files have no entry point, imports, effect adapter, authorization profile, or runnable composition and were not loaded or executed. |
| C5b4 | Recover exact libkrunfw bytes | `PASSED` | Merge [`068e221dafa7cf3e9a945cee7e8bf077eeed1c6b`](https://github.com/Shrimpworks/capsule-experiments/tree/068e221dafa7cf3e9a945cee7e8bf077eeed1c6b/experiments/typed-guest-transport-c5b4-libkrunfw-recovery) retains the official release input and two byte-identical, network-denied builds of the exact 24,339,104-byte `libkrunfw.5.dylib`. Full preferred-form kernel source/configuration/patch/tool closure remains `BLOCKED`; the extracted kernel stays evidence-only and separate firmware remains inapplicable under Accepted ADR-0041. |
| C5b5 | Construct the compile-only descriptive effect adapter | `PASSED` | Merge [`3cfe7db16c55894be444d4c783659043dbd25c95`](https://github.com/Shrimpworks/capsule-experiments/tree/3cfe7db16c55894be444d4c783659043dbd25c95/experiments/typed-guest-transport-c5b5-no-run-effect-adapter) retains two byte-equal non-executable objects, closed profile/action translation, and exactly 13 reviewed undefined libkrun symbols. It describes requested operations but implements/invokes none; the security-critical real effect implementation remains `BLOCKED`. |
| C5b6 | Reproduce the governed fixed-fixture Deno runtime | `PASSED` | Merge [`d9967e80a6155a65c6876dc686d8f8498b4a908f`](https://github.com/Shrimpworks/capsule-experiments/tree/d9967e80a6155a65c6876dc686d8f8498b4a908f/experiments/typed-guest-transport-c5b6-deno-static-reproduction) retains two independent exact Cargo acquisitions and two byte-identical network-disabled builds of the 68,496,520-byte runtime, 699,988-byte snapshot, and deterministic bundle. The retained static-only builder removed all candidate invocation; no output was loaded or executed. This fixed-fixture identity does not by itself close governed release publication or admission. |
| C5b7 | Rebuild the immutable runtime root | `PASSED` | Archive merge [`78485fb91a31733c568fe43e5fa295474e5956e1`](https://github.com/Shrimpworks/capsule-experiments/tree/78485fb91a31733c568fe43e5fa295474e5956e1/experiments/typed-guest-transport-c5b7-deterministic-runtime-root) retains two independently assembled byte-identical 100,663,296-byte ext4 roots at SHA-256 `5ad18f20cbc97c7a70ead3e795fd3649672513323041e913b0eb55b7acc88775`. They bind the exact C5b6 runtime/snapshot, C5b1 trusted init/launcher, C5b0 source/manifest/input, and C5b3/C5b5 descriptor/transport metadata; a closed 19-node inventory, independent raw-filesystem verifier, and 15 mutation refusals pass without loading or executing an artifact. The root is an explicit successor to C5b1, not byte-equivalent to it. C5b5's historical 134,217,728-byte contract remains unchanged; C5b8's later immutable successor resolves the selected-root binding without rewriting either predecessor. |
| C5b8 | Implement sealed controlled-test effect sequencing | `PASSED` | Archive merges [`e83614af34d5c39c12a4a3d6e6cda8dcf0304030`](https://github.com/Shrimpworks/capsule-experiments/tree/e83614af34d5c39c12a4a3d6e6cda8dcf0304030/experiments/typed-guest-transport-c5b8-controlled-test-effects) and [`b0819d76883eb86cbbc03b2b7033fe55bedbf713`](https://github.com/Shrimpworks/capsule-experiments/tree/b0819d76883eb86cbbc03b2b7033fe55bedbf713/experiments/typed-guest-transport-c5b8-c5b7-root-binding-successor) retain the sealed operation-sequencing layer, byte-equal no-run objects, complete test-double success/fault/replay/cleanup coverage, and the exact C5b7-root binding. No caller-selected path, flag, image, mount, endpoint, backend configuration, retained dylib, runtime, libkrun, HVF, VM, or guest was used. The fixed real operation provider remains outside this passed scope. |
| C5b9 | Bind the complete immutable no-run composite | `PASSED` | Archive merge [`3965e6b5cc87d476da7f431d7ed8a5758011a1b8`](https://github.com/Shrimpworks/capsule-experiments/tree/3965e6b5cc87d476da7f431d7ed8a5758011a1b8/experiments/typed-guest-transport-c5b9-immutable-no-run-composite) binds the exact runner, libkrun, libkrunfw, 100,663,296-byte root, controller, and root-bound effects object. Static verification closes the controller and 13-symbol libkrun surfaces, 14-file archive inventory, typed caps and completion-last fixture, teardown ordering, all predecessor verifiers, nine unit tests, and 14 mutations. `_c5b8_controlled_test_operation` deliberately has no provider; host, guest, authorization, and every effect remain absent. Nothing was loaded or executed, and no v19/v27 identity was reused. |
| C5b compatibility preflight | Test direct provider-only composition | `PASSED`; exact candidate `NO_GO` | Archive merge [`7fc3af9c46895b340c3118a96cb50abb26b1d977`](https://github.com/Shrimpworks/capsule-experiments/tree/7fc3af9c46895b340c3118a96cb50abb26b1d977/experiments/typed-guest-transport-c5b-controlled-harness-preflight) retains exact component identities, four closed contradictions, static source/Mach-O verification, ten mutations, and a closed archive inventory. It abandons only binding the retained C5b9 inputs by supplying the missing operation symbol: runner/root identity, effect order, operation ABI, and single-libkrun-owner requirements do not compose truthfully. No native artifact, libkrun/HVF, runner, VM, or guest executed. |
| C5b11 | Bind the fault-convergent fixed-runner no-run successor | `PASSED` for construction/static evidence | [Immutable merge `f206e4ef2cd326ee74e5b7b2739c62efe6da7d6d`](https://github.com/Shrimpworks/capsule-experiments/tree/f206e4ef2cd326ee74e5b7b2739c62efe6da7d6d/experiments/typed-guest-transport-c5b11-bound-fault-convergent-no-run-successor) retains exact plan/payload/profile binding, one runner importing 13 libkrun symbols, a Supervisor driver importing zero libkrun and 24 closed providers, distinct restart cursors, and fault/replay/teardown models. PR #31 reports exact-head C5b-S5 review `PASSED` / `Ready`; the [checkpoint](C5B_FIXED_RUNNER_SUCCESSOR_CHECKPOINT.md) separates that publication from archived pre-review prose and fresh verification. Providers and effects remain absent. C5b10 is not accepted evidence. |
| C5b providers | Construct the missing fixed process/transport/recovery providers | `BLOCKED` on implementation | Start from the exact C5b11 ABI, attempt bindings, and independent recovery oracle. Freeze the bounded no-run source/provenance/test packet, implement real provider source with closed inputs/imports, reproduce its artifacts, and independently review the new composition. Declarations or test-double results cannot prove provider/platform behavior; no guest execution is authorized by this task row. |
| C5b | Run the controlled typed-transport harness | `BLOCKED` | Requires provider implementation/provenance and independent review of the complete exact composition, then final owner authorization naming its immutable merge and manifest. Retain directional copy, chunk/cap+1, stall/reset/cancel, descriptor substitution, response-loss, completion-last, teardown, and restoration evidence without making an admission decision. |
| C6a | Build the installed authenticated service and protected-state boundary | `BLOCKED` | Requires passed C2b and C3c under Accepted ADR-0029, then separate authorization for the Keychain/service/protected-root corpus. C3c must supply Accepted ADR-0038/0045 decisions or accepted replacements that freeze the authority descriptor and state-engine binding. Implement only method-specific listeners, peer authentication, owner/store startup, and the four passively frozen Supervisor consumers. |
| C6b1a | Build the unsigned Broker evidence harness | `PASSED` | Archive merge [`4a2447d4bd0e03132dc616e608031ca313630cdd`](https://github.com/Shrimpworks/capsule-experiments/tree/4a2447d4bd0e03132dc616e608031ca313630cdd/experiments/broker-live-signing-c6b1) retains the unsigned Swift/Objective-C target, requested entitlement inputs, deterministic closed fixture corpus, public-only signature/binding checks, no-credential interaction double, independent verifier, and stable future seam interface. It used no Apple identity/profile, Keychain, LocalAuthentication, signing, installation, listener, runtime, backend, VM, guest, or product consumer. |
| C6b1b | Build the test-only Supervisor evidence seam | `PASSED` | Archive merge [`067fe2beb40361bb714507cab1331004e0a656fa`](https://github.com/Shrimpworks/capsule-experiments/tree/067fe2beb40361bb714507cab1331004e0a656fa/experiments/broker-live-signing-c6b1-supervisor-seam) retains six ordered approval/attempt commit, replay, response-loss, reopen, and concurrency rows. Canonical payload plus resolved signer authorization is replay identity; the Supervisor experiment store is the only durable authority owner. The model is test-only and is never imported or promoted into product code. |
| C6b1c | Provision and read back the disposable identity | `PASSED` | Archive merge [`82d1a799f70482856aaa6030f612d701b39cec67`](https://github.com/Shrimpworks/capsule-experiments/tree/82d1a799f70482856aaa6030f612d701b39cec67/experiments/broker-live-signing-c6b1c-signed-artifact-readback) retains the exact development profile metadata and no-install signed app. Strict signature, exact bundle/Team/designated requirement, hardened runtime, and the closed App Sandbox plus one Approval Keychain-group entitlement readback pass. The profile wildcard is only an allowlist; no raw profile is embedded, and the app was never installed or launched. |
| C6b1d | Run the installed Broker signing evidence matrix | `BLOCKED` | Requires a fresh authorization naming the exact C6b1c archive/artifact/profile, owner account/container, allowed Keychain and LocalAuthentication operations, prompt handling, evidence destination, and cleanup. Proposed first-run destructive rows D1-D4 and cleanup D14-D16 are not yet authorized; D5-D13 and D17-D18 remain deferred. No product consumer, runtime, backend, VM, or guest. |
| C6b2 | Connect the product Broker and approval/attempt methods | `BLOCKED` | C4 is `PASSED`; C6a and C6b1d remain required. Implement native rendering/UI, installed signing/public-key verification, and method-specific `SubmitApprovalV0`/`RequestAttemptV0` consumers without runtime or guest activation. Research and passive conformance cannot satisfy either installed dependency. |
| C6c | Wire attempt admission and the fixed-store stop policy | `BLOCKED` | Requires C6b2 and an explicit decision for p95 provenance/window/lifetime and any persistent timing-trip semantics. Apply the re-evaluated guard transaction-locally after owner/full verification and before a new consume/create mutation; replay of an existing `AttemptID` converges first. |
| C6d | Run the pre-admission installed owner-only internal-alpha runtime/profile matrix | `BLOCKED` | Requires C1b, C5b, C6a, governed release/artifact review, and separate authorization naming the exact signed-installed candidate and owned test environments. Retain identity compatibility, runtime/root/no-loader and string-codegen restoration, transport, teardown/recovery, and the required broader lifecycle/platform evidence without accepting user source or making an admission decision. Source Validator R4/R5 is not part of this ADR-0040 matrix. |
| C7 | Review one exact owner-only internal-alpha runtime/profile candidate for admission | `BLOCKED` | Requires C6d. Produce an explicit admit-or-refuse result over the exact ADR-0040 candidate and retained evidence; controlled v19/v27 experiments alone cannot admit it. R4/R5 is not a prerequisite for this review and no C7 result admits the post-alpha Product Source Validator or any external-alpha profile. |
| C8 | Connect the sealed real adapter and completion-last path | `BLOCKED` | Requires C6c, C7, and a separately authorized owned guest. Execute only by committed `AttemptID`; consume real result-integrity, runner, teardown, and absence facts. |
| C9 | Run the installed hostile-`.mjs` admission corpus | `BLOCKED` | Requires C8. Response loss, restart, recovery, restoration, pressure, sleep/wake, update, and the minimum hostile source/authority/transport/root/lifecycle cases must converge in the exact signed-installed profile. |

Independent repository-quality work may continue without changing security claims: issue #217 in
one-package documentation batches after the now-passed `lifecyclestate` and `registrationstate`
batches; issue #219 as sequential behavior-preserving archive refactors; and issue #216 only after
its threshold/exemption policy is frozen. Issue #218 package/API reduction follows the high-churn
#219 work. The Q1-Q8 intake above is ordered ahead of those broad backlogs where it closes a
specific false-green or fail-closed gap; each remaining issue remains a separate slice.

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
  R4-v2 is unexecuted and `BLOCKED`. Those slices gate admission of the Product Source Validator
  itself under ADR-0035/0036; they do not gate the ADR-0040 owner-only internal-alpha
  runtime/profile review at C7. Any external-alpha requirement must be frozen separately before
  that later admission path opens.
- TypeScript Source Preparer, automatic TUF updates, Developer ID/notarized distribution,
  clean-host/minimum-OS coverage, independent-builder equality, restore activation, continuity,
  and external-alpha distribution are not on the owner-only internal-alpha critical path.

## Next checkpoint

C2b0, C3b's complete no-launch profile/signature preflight, C5b1-C5b9 no-run construction, the
C5b compatibility preflight, and C6b1c are now `PASSED` in the exact bounded scopes pinned above.
The direct C5b9 provider-only binding candidate is `NO_GO`; C5b controlled execution remains
`BLOCKED` on the missing C5b11 provider implementations/provenance, independent review of the
complete composition, and final authorization naming its exact immutable merge and manifest.
C5b11 construction/static re-verification is `PASSED`; see the September 5 checkpoint above.
The immediate C3b step is a freshly
authorized E1 container matrix
over E1-01..E1-12 and E1-14..E1-15; E1-13 remains excluded. If C3b passes, C3c must reconcile that
evidence and decide ADR-0038/0045 before C6a begins. C2b native execution, the C3b
container matrix, C5b controlled execution, and C6b1d live signing each require their own exact
authorization. C1a resumes only if the owner
supplies a backup or snapshot; otherwise a later exact rerun needs separate authorization and new
identities.

This plan records but does not itself authorize new signing, Keychain use, service registration,
installed listener activation, libkrun/HVF execution, VM, or guest work. C2b, the C3b container
matrix, C5b, and C6b1d must each stop again for their exact owner/environment authorization. No
task may promote internal alpha or
product admission from research, a passive contract, an ad-hoc harness, or a fixed guest
experiment alone.
