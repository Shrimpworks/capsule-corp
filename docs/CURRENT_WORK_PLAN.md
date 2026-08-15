# Current work plan

Date: 2026-08-14

Work item: reconcile the merged C5b7 deterministic no-run runtime-root evidence and advance the
exact next composition gates.

Status: `PASSED` for this canonical reconciliation and the completed child scopes named below.

Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`.

Product admission and the installed security boundary: `BLOCKED`.

This is the current execution index. Detailed security claims remain in the linked ADRs, readiness
map, experiment checkpoints, and evidence ledger. A completed passive contract or controlled
experiment is not an activated product path.

## Reconciled baseline

This reconciliation starts from fetched `origin/main` commit
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
- C5b0-C5b2 packet/executable/input closure, the five C5b3-C5b6 no-run slices, and C5b7's
  deterministic runtime-root construction are `PASSED`. Exact `rusty_v8`, governed fixed-fixture
  Deno runtime, libkrunfw, controller-core, descriptive adapter, and runtime-root bytes are now
  retained. The C5b7 root is intentionally versioned and cannot be combined as-is with C5b5's
  frozen 134,217,728-byte root contract. A reviewed size-compatible real effect implementation,
  immutable complete composite, controlled C5b run, preferred-form libkrunfw/kernel source
  compliance, and admission remain `BLOCKED`; and
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
                              -> ADR-0045 E1 identity separation +-> key/service/root corpus
                                                               │
C4 passive approval/attempt evidence (PASSED) + R3 passive research (PASSED)
          -> unsigned Broker harness (`PASSED`) -> test Supervisor seam (`PASSED`)
          -> identity/profile readback (`PASSED`) -> installed signing harness -----------------┐
installed authenticated IPC boundary (BLOCKED) ----------------------------------------------+-> product Broker/approval/attempt wiring
                                                                                              -> protected one-attempt path
                                                               │
typed transport design -> passive contract (`PASSED`) -> v19/103-byte no-run packet (`PASSED`)
                                                   -> fresh executable construction (`PASSED`)
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
| C3a | Materialize deterministic E0 fixtures | `PASSED` | Archive merge [`dee784d40684100f8315720fab9a5cd3399f492b`](https://github.com/Shrimpworks/capsule-experiments/tree/dee784d40684100f8315720fab9a5cd3399f492b/experiments/macos-installation-i2b3-supervisor-authority-epoch-e0) retains exact current/legacy probe sources and reproducible unsigned bundles, a never-launched Coordinator, plists, entitlement/profile requests, disabled LaunchAgent and inactive descriptor inputs, a closed manifest, independent verification, and 23 mutation refusals. No portal, identity, profile, signing, container, service, Keychain, runtime, backend, VM, or guest was accessed or activated. |
| C3b | Run ADR-0045 E1 identity separation | `BLOCKED` | C3a, exact legacy-profile restoration, the App Group portal preflight, and the [exact profile/signature-only gate](https://github.com/Shrimpworks/capsule-experiments/tree/ee00ae2abbce64ae6458b82d0b53d904ee39aeb6/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1-signed-profile-preflight) are complete. That gate `PASSED` exact Team/application identifiers, profile UUID/CMS/certificate/device binding, CDHashes/designated requirements, effective App Group/Keychain entitlements, hardened runtime, and absent debug entitlement without launching a bundle or opening a container. The frozen `3DDR84M4JS...` App Group remains the macOS-style entitlement value and is not a Developer-portal App Group resource; do not rewrite it to `group.`. A fresh authorization may now run only E1-01..E1-12 and E1-14..E1-15 against the exact retained identities; E1-13 remains excluded. ADR-0045 remains Proposed. |
| C4 | Freeze `SubmitApprovalV0` and `RequestAttemptV0` | `PASSED` | PR #248 is the canonical predecessor and PR #249 closes the focused follow-up. Ordered 4,999/5,000/5,001-ms cases for both methods, equality-as-expired behavior, complete closed dictionaries/maps, every ordered field, required `noState`, cancellation/deadline commit truth, all 20 foreign-tag collisions, complete refusal and five-entry response-loss tables, and bounded Go/Node mutation proofs are retained. No listener, signer, store consumer, process, or guest is active. |
| C5a | Freeze the final typed source/input/completion transport | `PASSED` | The passive v1 contract freezes exact 152/160/64-byte layouts, 262,144-byte payload caps, completion cap-plus-one, big-endian bindings, four statuses, canonical JSON, refusal precedence, monotonic state/fault behavior, endpoint custody, completion-last projection, deterministic fixtures, and independent Go/Node verification. No endpoint, process, runtime, backend, guest, or store mutation occurred. |
| C5b0 | Materialize the deterministic no-run typed-transport packet | `PASSED` | Archive merge [`b357d0c0fb29100c180494e67cebd7809aabe3c5`](https://github.com/Shrimpworks/capsule-experiments/tree/b357d0c0fb29100c180494e67cebd7809aabe3c5/experiments/typed-guest-transport-c5b0-v19-successor) binds the v19 lineage digest, governed 103-byte source and SourceManifest, exact role contracts, no-run profile/plan, fresh typed frames, closed inventory, independent verifier, and six mutations. No v19 raw bytes were recreated; executable runner/root/init/launcher/controller identities remain explicitly null. |
| C5b1 | Construct the fresh executable typed-transport successor | `PASSED` | Archive merge [`db08ebf277432e06d6cba3b7f7338e3bd4a61252`](https://github.com/Shrimpworks/capsule-experiments/tree/db08ebf277432e06d6cba3b7f7338e3bd4a61252/experiments/typed-guest-transport-c5b1-executable-successor) retains fresh reproducible runner, raw root, trusted init, launcher, and hard-stop controller candidates, a closed 41-file inventory, provenance/SBOM, independent parsing, and seven mutation refusals. It does not recover v19 or bind/run the governed runtime, libkrun/libkrunfw, kernel, firmware, or a real controller. |
| C5b2 | Bind the available governed successor inputs | `PASSED` | Archive merge [`5a2f835e8c9df8279237f940f5af757e119593bd`](https://github.com/Shrimpworks/capsule-experiments/tree/5a2f835e8c9df8279237f940f5af757e119593bd/experiments/typed-guest-transport-c5b2-governed-input-closure) independently binds the exact current-source libkrun header, ABI audit, unsigned dylib, and final runner source/binary by hash, raw Mach-O metadata, dependencies, exports, imports, and C17 audit. In that immutable slice, libkrunfw/kernel receipts were identity evidence only and the governed `deno_core` executable/libkrunfw bytes and complete controller/composite were absent; C5b3-C5b7 later close only the no-run input/controller/runtime-root portions without rewriting C5b2. Separate firmware remains inapplicable under ADR-0041. Nothing was loaded or executed. |
| C5b3 | Recover runtime inputs and construct the controlled-test core | `PASSED` | Merges [`b5db7bcbbf7fe814faec3950ebfbf2d2ac7786e2`](https://github.com/Shrimpworks/capsule-experiments/tree/b5db7bcbbf7fe814faec3950ebfbf2d2ac7786e2/experiments/typed-guest-transport-c5b3-runtime-input-recovery) and [`60234e22674e46a42e8e5c382d85217a930c2c13`](https://github.com/Shrimpworks/capsule-experiments/tree/60234e22674e46a42e8e5c382d85217a930c2c13/experiments/typed-guest-transport-c5b3-controlled-test-controller) retain exact `rusty_v8` archive/binding custody and the pure no-effect C17 state machine. Its byte-equal `MH_OBJECT` files have no entry point, imports, effect adapter, authorization profile, or runnable composition and were not loaded or executed. |
| C5b4 | Recover exact libkrunfw bytes | `PASSED` | Merge [`068e221dafa7cf3e9a945cee7e8bf077eeed1c6b`](https://github.com/Shrimpworks/capsule-experiments/tree/068e221dafa7cf3e9a945cee7e8bf077eeed1c6b/experiments/typed-guest-transport-c5b4-libkrunfw-recovery) retains the official release input and two byte-identical, network-denied builds of the exact 24,339,104-byte `libkrunfw.5.dylib`. Full preferred-form kernel source/configuration/patch/tool closure remains `BLOCKED`; the extracted kernel stays evidence-only and separate firmware remains inapplicable under Accepted ADR-0041. |
| C5b5 | Construct the compile-only descriptive effect adapter | `PASSED` | Merge [`3cfe7db16c55894be444d4c783659043dbd25c95`](https://github.com/Shrimpworks/capsule-experiments/tree/3cfe7db16c55894be444d4c783659043dbd25c95/experiments/typed-guest-transport-c5b5-no-run-effect-adapter) retains two byte-equal non-executable objects, closed profile/action translation, and exactly 13 reviewed undefined libkrun symbols. It describes requested operations but implements/invokes none; the security-critical real effect implementation remains `BLOCKED`. |
| C5b6 | Reproduce the governed fixed-fixture Deno runtime | `PASSED` | Merge [`d9967e80a6155a65c6876dc686d8f8498b4a908f`](https://github.com/Shrimpworks/capsule-experiments/tree/d9967e80a6155a65c6876dc686d8f8498b4a908f/experiments/typed-guest-transport-c5b6-deno-static-reproduction) retains two independent exact Cargo acquisitions and two byte-identical network-disabled builds of the 68,496,520-byte runtime, 699,988-byte snapshot, and deterministic bundle. The retained static-only builder removed all candidate invocation; no output was loaded or executed. This fixed-fixture identity does not by itself close governed release publication or admission. |
| C5b7 | Rebuild the immutable runtime root | `PASSED` | Archive merge [`78485fb91a31733c568fe43e5fa295474e5956e1`](https://github.com/Shrimpworks/capsule-experiments/tree/78485fb91a31733c568fe43e5fa295474e5956e1/experiments/typed-guest-transport-c5b7-deterministic-runtime-root) retains two independently assembled byte-identical 100,663,296-byte ext4 roots at SHA-256 `5ad18f20cbc97c7a70ead3e795fd3649672513323041e913b0eb55b7acc88775`. They bind the exact C5b6 runtime/snapshot, C5b1 trusted init/launcher, C5b0 source/manifest/input, and C5b3/C5b5 descriptor/transport metadata; a closed 19-node inventory, independent raw-filesystem verifier, and 15 mutation refusals pass without loading or executing an artifact. The root is an explicit successor to C5b1, not byte-equivalent to it. C5b5 freezes 134,217,728 root bytes and is metadata-only/incompatible as-is with this 100,663,296-byte root, so C5b8 must bind the exact selected root size before C5b9. |
| C5b8 | Implement the real controlled-test effects | `BLOCKED` | Implement and independently review the narrow operation layer behind the C5b3 core and C5b5 descriptive adapter. It must add no caller-selected paths, flags, images, mounts, endpoints, or backend configuration and must be exercised only through test doubles before composition. |
| C5b9 | Bind the complete immutable no-run composite | `BLOCKED` | Requires passed C5b7 and a C5b8 implementation that explicitly resolves C5b7's 100,663,296-byte root against C5b5's historical 134,217,728-byte contract. Bind exact runner, libkrun, libkrunfw, runtime root, controller, effect implementation, fixtures, and authorization-profile placeholders; independently verify closed inventory, ABI/load surfaces, caps, completion-last, teardown/cleanup, and restoration mutations. Stop before loading libkrun. |
| C5b | Run the controlled typed-transport harness | `BLOCKED` | Requires C5b9 plus separate authorization naming that exact successor, owner-confirmed host, and owned disposable guest. Retain directional copy, chunk/cap+1, stall/reset/cancel, descriptor substitution, response-loss, completion-last, teardown, and restoration evidence without making an admission decision. |
| C6a | Build the installed authenticated service and protected-state boundary | `BLOCKED` | Requires passed C2b and C3b evidence under Accepted ADR-0029, then separate authorization for the Keychain/service/protected-root corpus. Implement only method-specific listeners, peer authentication, owner/store startup, and the four passively frozen Supervisor consumers. |
| C6b1a | Build the unsigned Broker evidence harness | `PASSED` | Archive merge [`4a2447d4bd0e03132dc616e608031ca313630cdd`](https://github.com/Shrimpworks/capsule-experiments/tree/4a2447d4bd0e03132dc616e608031ca313630cdd/experiments/broker-live-signing-c6b1) retains the unsigned Swift/Objective-C target, requested entitlement inputs, deterministic closed fixture corpus, public-only signature/binding checks, no-credential interaction double, independent verifier, and stable future seam interface. It used no Apple identity/profile, Keychain, LocalAuthentication, signing, installation, listener, runtime, backend, VM, guest, or product consumer. |
| C6b1b | Build the test-only Supervisor evidence seam | `PASSED` | Archive merge [`067fe2beb40361bb714507cab1331004e0a656fa`](https://github.com/Shrimpworks/capsule-experiments/tree/067fe2beb40361bb714507cab1331004e0a656fa/experiments/broker-live-signing-c6b1-supervisor-seam) retains six ordered approval/attempt commit, replay, response-loss, reopen, and concurrency rows. Canonical payload plus resolved signer authorization is replay identity; the Supervisor experiment store is the only durable authority owner. The model is test-only and is never imported or promoted into product code. |
| C6b1c | Provision and read back the disposable identity | `PASSED` | Archive merge [`82d1a799f70482856aaa6030f612d701b39cec67`](https://github.com/Shrimpworks/capsule-experiments/tree/82d1a799f70482856aaa6030f612d701b39cec67/experiments/broker-live-signing-c6b1c-signed-artifact-readback) retains the exact development profile metadata and no-install signed app. Strict signature, exact bundle/Team/designated requirement, hardened runtime, and the closed App Sandbox plus one Approval Keychain-group entitlement readback pass. The profile wildcard is only an allowlist; no raw profile is embedded, and the app was never installed or launched. |
| C6b1d | Run the installed Broker signing evidence matrix | `BLOCKED` | Requires a fresh authorization naming the exact C6b1c archive/artifact/profile, owner account/container, allowed Keychain and LocalAuthentication operations, prompt handling, evidence destination, and cleanup. Proposed first-run destructive rows D1-D4 and cleanup D14-D16 are not yet authorized; D5-D13 and D17-D18 remain deferred. No product consumer, runtime, backend, VM, or guest. |
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

C2b0, C3b's complete no-launch profile/signature preflight, C5b1-C5b7 no-run construction, and
C6b1c are now `PASSED` in the exact bounded scopes pinned above. C5b8 real-effect implementation
may proceed without runtime execution and must resolve C5b7's exact root-size binding before the
C5b9 immutable no-run composite. The immediate C3b step is a freshly authorized E1 container matrix
over E1-01..E1-12 and E1-14..E1-15; E1-13 remains excluded. C2b native execution, the C3b
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
