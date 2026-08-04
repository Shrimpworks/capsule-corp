# Governed deno_core work plan

Status: active planning checkpoint after accepted ADR-0028, the governed Linux/arm64 construction
handoff, the passed passive C1 controlled-development composition contract, and the passed passive
C2A execution-profile preparation contract.
This plan selects work order; it does not admit a runtime, backend, profile, or guest.

## Current checkpoint

- Capsule `main` includes durable lifecycle Slice E5, the V8 source/license closure result, and the
  self-contained dynamic-root result.
- [`Shrimpworks/deno`](https://github.com/Shrimpworks/deno) is the current governed Deno
  integration repository and remains a fork of `denoland/deno`. It transferred from
  `dills122/deno` with its default branch, `capsule/upstream-v2.9.4`, merged PR #1, and Actions
  history intact while idle. The transfer is repository governance only, not runtime, artifact,
  review, release, or security evidence.
- [`Shrimpworks/rusty_v8`](https://github.com/Shrimpworks/rusty_v8) is a real fork of
  `denoland/rusty_v8`. It transferred from `dills122/rusty_v8` with its governed branches, four
  merged PRs, and Actions history intact.
- Deno governed PR #1 is merged at head `9adb0b68b55bca81644827f1e7749a3acb091bed`
  and merge `ea18b9dc21ff8ebd19347be7095f47937ee14ec2`.
- `rusty_v8` governed follow-up PR #2 is merged at head
  `a43ee7486c3e05bce5d6e5db586b3e2e688c33cf` and merge
  `a31b8f39dc6933d5635367e8ccb67d70f2cc2385`. The stale original head
  `17698ca...f313a` is not a terminal integration input.
- The first fork-native integration check stopped before prefetch or build because `rusty_v8`
  provides only Linux/amd64 and no Linux/arm64 builder/publication profile. No release or new
  runtime artifact exists.
- Governed `rusty_v8` PR #4 is merged at exact head
  `80e863ddb942a4aa2b384e794fc23e35b9d2bb15` and merge
  `cbf56de2e1156b1cf1561fdbaea7172a0aa056f4`. Its clean Linux/arm64 build, fixed test, corrected
  GN evidence query, network-disabled full build, and unsigned bundle upload pass. Independent-
  builder equality, exact evidence review, governed release publication, and admission remain.
- Governed `deno_core` is the first runtime engineering candidate. `RUNTIME-001` remains
  unsupported and no product execution path may use these bytes.
- `Shrimpworks/capsule-experiments` PR #1 is merged at
  `fa03d7043b4f0653081d6c5733d597f49f6efd1c`. Its exact Linux/arm64 construction handoff is
  `PASSED`; all eight raw/compressed runtime and root identities are pinned in the passive
  [C1 composition contract](protocol/GOVERNED_DENO_CORE_C1_COMPOSITION.md). C1 itself is `PASSED`
  only as a no-effect contract. Independent-builder equality, governed release publication,
  composed execution, and admission remain open.
- [C2A execution-profile preparation](protocol/GOVERNED_DENO_CORE_C2A_EXECUTION_PROFILE.md) is
  `PASSED` only for its passive, refusing contract. It consumes the unchanged C1 bytes and freezes
  the exact host/guest descriptor maps, bounded candidate machine and transport values, teardown,
  known answer, restoration corpus, and complete C2B matrix. It creates no guest or runtime. C2B
  is `BLOCKED` on final runnable artifact identities, the launcher-child contract, exact scratch
  enforcement, remaining governed review/P0 closures, and separate guest authorization.
- Durable lifecycle Slices E1 through E5 are implemented locally and unwired.

## Priority 1: completed governance bootstrap

The three original parallel tasks completed as local/fork governance checkpoints. None admitted a
runtime or release.

### 1A. Capsule durable lifecycle Slice E4

Repository: `Shrimpworks/capsule-corp`.

- Replace the registered lifecycle `MemoryStore` path with the E3 fixed-store transaction port.
- Give every fake operation a stable `EffectID` and exact fake instance identity.
- Serialize `Drive`, `Recover`, and startup per `AttemptID` under one injected owner session.
- Exercise prepare/create/start/observe/stop/destroy intent, result, indeterminate, process-death,
  reopen, and reconciliation paths.
- Preserve `FakeBackend.CreatesGuest() == false`; add no process, runtime, backend, guest, consumer,
  IPC, content, or evidence path.

Exit: completed through E5; the fake-only lifecycle remains unwired and creates no guest.

### 1B. Bootstrap the governed Deno fork

Repository: `Shrimpworks/deno` (transferred from the historical `dills122/deno` location without
changing the retained source or merge identities below).

- Create the governed line from exact upstream commit
  `14eea3160ae5834476aa3b9d317b8d41d991b982`.
- Carry the retained physical-omission and deterministic-ordering changes as separate reviewable
  commits, preserving their exact source identities and restoration mutations.
- Add fork-local scope, ownership, upstream-sync, patch-removal, review, CI, and release rules.
- Reproduce the three-op registry/final-link result, snapshot and binary known answers, fixed
  fixture, and prohibited restoration cases from clean builders.
- Do not publish an admitted Capsule runtime or wire the fork into Capsule in this slice.

Exit: completed by merged Deno head `9adb0b68...91bed`; exact ancestry, patch/fixture identities,
and fork-local three-op/restoration verification pass. The fork did not produce fresh binary or
snapshot bytes.

### 1C. Bootstrap the governed rusty_v8 fork

Repository: `Shrimpworks/rusty_v8`.

- Create the governed line from exact upstream commit
  `d305e6afa7736f6e298c30ae6646f7709ee9382b`.
- Replace mutable publisher assumptions with pinned builder inputs and retained effective GN/Ninja
  configuration and linked-component closure.
- Publish the complete corresponding source, generated build metadata, licenses/notices, SBOM,
  provenance, archive digest, and exact release procedure.
- Add advisory ownership, upstream-update/rebase policy, review gates, and byte/restoration tests.
- Do not fork `denoland/v8` unless the governed work actually changes that source.

Exit: completed through the later merged Linux/arm64 head `80e863dd...2bb15`. The exact
20-gitlink source lock and offline build/publication contract now have retained clean Linux/arm64
construction evidence. No governed release candidate was published or admitted.

## Priority 2: Linux/arm64 construction and evidence closure

The fork-native Linux/arm64 construction subset is `PASSED`; remaining Priority 2 evidence is
`IN_PROGRESS — TRENDING_GOOD`. Neither status implies runtime admission:

1. Add a fully digest-pinned Linux/arm64 sibling builder/publication profile to `rusty_v8` at or
   after `a43ee748...33cf`, covering the builder image, host tools, aarch64 sysroot/target, offline
   build/test, generated GN/Ninja closure, output verification, source/notices, SBOM, and unsigned
   provenance. Do not substitute the existing amd64 profile.
2. Reconstruct the governed Deno, `rusty_v8`, binary, snapshot, and 22-entry root first in clean
   same-host Linux/arm64 builders. A genuinely independent Linux/arm64 builder is viable but not
   currently planned; same-host and GitHub-CI equality must remain explicitly limited, and
   independent-builder equality is deferred.
3. Follow accepted ADR-0034 for the first release: one byte-exact pass-through `main.mjs`, one
   plan-v0 source role, no static/dynamic dependency request, and no filesystem/network/package/
   fallback module loader. The passive M1 source-byte/SourceManifest foundation is retained, but
   the module-request validator and downstream M2/S1 work are blocked by the exact
   division-versus-regexp counterexample until a separate parser-boundary decision. Complete M1-M4
   before runtime integration, then prove M5 only under a separate runtime/profile admission plan.
   Proposed ADR-0032 remains on P1 HOLD for a conditional later TypeScript feature and is not a
   first-release dependency. Do not add CommonJS, package resolution, legacy Node module surface,
   or wider runtime authority.
4. Assemble one immutable runtime bundle manifest from governed fork releases and the exact dynamic
   root; rerun physical-op, final-link, file-open, syscall-seal, descriptor, and restoration tests.
5. Run an independent review of the combined libkrun FD-native, direct-block-root, and console
   patches and reproduce the 43 P0-3 frames in the selected host and launcher languages.

Items 1 through the same-host construction part of item 4 passed in the merged experiments handoff.
Independent-builder equality, a governed release, the combined libkrun review, and real
transport/launcher reproduction remain later requirements.

## Priority 3: composed development profile

- C1 is `PASSED`: the passive fixture freezes exact artifact identities, one byte-exact
  `main.mjs`, JSON input/completion bindings, desired globals/ops/modules/files, loader absence,
  logical descriptor roles, resource references, and refusal/stop behavior without creating a
  guest or activating any consumer.
- C2A is `PASSED` only as passive preparation: its refusing fixture binds C1, governed source
  identities, numeric FDs 0 through 7, one vCPU, 256 MiB guest RAM, 1,000 ms wall time,
  concurrency one, source/input/completion caps, `SIGKILL` teardown deadlines, the fixed benign
  known answer, and the 91-case C2B plus 18-mutation corpus. Unsupported CPU, host/VMM-memory, and
  scratch ceilings remain null and refusing; all nine final runnable artifact identities and the
  composed runtime-profile identity/digest remain unresolved.
- C2B is `BLOCKED`. Once its named artifact, launcher, resource, review, and P0 dependencies close,
  request separate authorization and run the fixed JSON-in/JSON-out job only in the owned
  disposable development guest. C2B must not broaden the C2A fixture or imply admission.
- Complete signed installed P0-1C/P0-4B custody, App Sandbox, Team identity, notarization, staple,
  Gatekeeper, clean-host, session/recovery, and supported-OS-floor evidence.
- Make a separate admission ADR only after every required control row has exact retained evidence.

The composed guest step requires explicit authorization for an owned disposable development guest.
Apple Development identities and provisioning profiles must likewise be deliberately authorized
before the existing installed matrices run. Exact G3 readback and later Apple Membership Details
confirm Team `3DDR84M4JS`; the W4 common-name suffix is not a Team ID. Local signed/provisioned
experiments require an explicitly authorized exact identity plus role identifiers, entitlements,
and matching 3DDR profiles. The three cached Gate B/wildcard profiles have the wrong App IDs and
are not reusable. A separate Developer ID Application identity for Team `3DDR84M4JS` is later
distribution authority requiring explicit authorization and
matching-Team package design; it is not current development authority or evidence that Developer ID/
notarization work is currently planned. Paid owned clean-host/
minimum-OS coverage is not currently planned; it remains deferred activation/distribution evidence,
not a blocker for current local mechanics.

## Explicit non-goals for this checkpoint

- No general Deno runtime, Node-first reopening, package manager, network, FFI, subprocess, Worker,
  inspector, native addon, ambient filesystem, or post-approval transformation authority.
- No product import from `capsule-experiments` and no shipping dependency on an archived or local
  source copy.
- No claim that a fork, successful build, or fixed benign fixture makes the runtime secure.
- No weakening of `RUNTIME-001`, `VMM-001`, `SUPPLY-001`, `PATCH-001`, or the external-isolation
  requirement to make sequencing easier.
- No claim that an eventual transfer of related forks to the Shrimpworks organization changes
  source review, build provenance, release, or admission evidence; that transfer is later repository
  governance work.
