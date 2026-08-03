# Governed deno_core work plan

Status: active planning checkpoint after accepted ADR-0028, the merged fork-governance branches,
and the first fork-native Linux/arm64 fail-fast result.
This plan selects work order; it does not admit a runtime, backend, profile, or guest.

## Current checkpoint

- Capsule `main` includes durable lifecycle Slice E5, the V8 source/license closure result, and the
  self-contained dynamic-root result.
- [`dills122/deno`](https://github.com/dills122/deno) is a real fork of `denoland/deno`.
- [`dills122/rusty_v8`](https://github.com/dills122/rusty_v8) is a real fork of
  `denoland/rusty_v8`.
- Deno governed PR #1 is merged at head `9adb0b68b55bca81644827f1e7749a3acb091bed`
  and merge `ea18b9dc21ff8ebd19347be7095f47937ee14ec2`.
- `rusty_v8` governed follow-up PR #2 is merged at head
  `a43ee7486c3e05bce5d6e5db586b3e2e688c33cf` and merge
  `a31b8f39dc6933d5635367e8ccb67d70f2cc2385`. The stale original head
  `17698ca...f313a` is not a terminal integration input.
- The first fork-native integration check stopped before prefetch or build because `rusty_v8`
  provides only Linux/amd64 and no Linux/arm64 builder/publication profile. No release or new
  runtime artifact exists.
- Governed `deno_core` is the first runtime engineering candidate. `RUNTIME-001` remains
  unsupported and no product execution path may use these bytes.
- Durable lifecycle Slices E1 through E5 are implemented locally and unwired.

## Priority 1: completed governance bootstrap

The three original parallel tasks completed as local/fork governance checkpoints. None admitted a
runtime or release.

### 1A. Capsule durable lifecycle Slice E4

Repository: `dills122/capsule-corp`.

- Replace the registered lifecycle `MemoryStore` path with the E3 fixed-store transaction port.
- Give every fake operation a stable `EffectID` and exact fake instance identity.
- Serialize `Drive`, `Recover`, and startup per `AttemptID` under one injected owner session.
- Exercise prepare/create/start/observe/stop/destroy intent, result, indeterminate, process-death,
  reopen, and reconciliation paths.
- Preserve `FakeBackend.CreatesGuest() == false`; add no process, runtime, backend, guest, consumer,
  IPC, content, or evidence path.

Exit: completed through E5; the fake-only lifecycle remains unwired and creates no guest.

### 1B. Bootstrap the governed Deno fork

Repository: `dills122/deno`.

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

Repository: `dills122/rusty_v8`.

- Create the governed line from exact upstream commit
  `d305e6afa7736f6e298c30ae6646f7709ee9382b`.
- Replace mutable publisher assumptions with pinned builder inputs and retained effective GN/Ninja
  configuration and linked-component closure.
- Publish the complete corresponding source, generated build metadata, licenses/notices, SBOM,
  provenance, archive digest, and exact release procedure.
- Add advisory ownership, upstream-update/rebase policy, review gates, and byte/restoration tests.
- Do not fork `denoland/v8` unless the governed work actually changes that source.

Exit: partially completed by merged follow-up head `a43ee748...33cf`. The exact 20-gitlink source
lock and an offline build/publication contract exist for Linux/amd64 only. The intended
Linux/arm64 exit remains blocked, and no release candidate was built or published.

## Priority 2: Linux/arm64 construction and evidence closure

Start only after the smallest required `rusty_v8` fork change lands:

1. Add a fully digest-pinned Linux/arm64 sibling builder/publication profile to `rusty_v8` at or
   after `a43ee748...33cf`, covering the builder image, host tools, aarch64 sysroot/target, offline
   build/test, generated GN/Ninja closure, output verification, source/notices, SBOM, and unsigned
   provenance. Do not substitute the existing amd64 profile.
2. Reconstruct the governed Deno, `rusty_v8`, binary, snapshot, and 22-entry root first in clean
   same-host Linux/arm64 builders, then on an independently controlled Linux/arm64 builder; retain
   the distinct provenance claims and compare every declared output byte.
3. Decide the production owner and process topology for the ADR-0026 strip-only transformation,
   then implement the coordinated versioned original-TypeScript/emitted-JavaScript plan migration.
4. Assemble one immutable runtime bundle manifest from governed fork releases and the exact dynamic
   root; rerun physical-op, final-link, file-open, syscall-seal, descriptor, and restoration tests.
5. Run an independent review of the combined libkrun FD-native, direct-block-root, and console
   patches and reproduce the 43 P0-3 frames in the selected host and launcher languages.

## Priority 3: composed development profile

- Bind the governed runtime bundle, exact kernel/init, libkrun patches, launcher, closed descriptor
  manifests, source/input/completion ports, and teardown behavior into one candidate profile.
- Run the fixed JSON-in/JSON-out job only in the owned disposable development guest.
- Complete signed installed P0-1C/P0-4B custody, App Sandbox, Team identity, notarization, staple,
  Gatekeeper, clean-host, session/recovery, and supported-OS-floor evidence.
- Make a separate admission ADR only after every required control row has exact retained evidence.

## Explicit non-goals for this checkpoint

- No general Deno runtime, Node-first reopening, package manager, network, FFI, subprocess, Worker,
  inspector, native addon, ambient filesystem, or post-approval transformation authority.
- No product import from `experiments/` and no shipping dependency on a local source copy.
- No claim that a fork, successful build, or fixed benign fixture makes the runtime secure.
- No weakening of `RUNTIME-001`, `VMM-001`, `SUPPLY-001`, `PATCH-001`, or the external-isolation
  requirement to make sequencing easier.
