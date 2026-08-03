# Governed deno_core work plan

Status: active planning checkpoint after merged PRs #54, #56, and #57 and accepted ADR-0028.
This plan selects work order; it does not admit a runtime, backend, profile, or guest.

## Current checkpoint

- Capsule `main` includes durable lifecycle Slice E3, the V8 source/license closure result, and the
  self-contained dynamic-root result.
- [`dills122/deno`](https://github.com/dills122/deno) is a real fork of `denoland/deno`.
- [`dills122/rusty_v8`](https://github.com/dills122/rusty_v8) is a real fork of
  `denoland/rusty_v8`.
- Both forks currently contain only inherited `main`; no governed branch or release exists.
- Governed `deno_core` is the first runtime engineering candidate. `RUNTIME-001` remains
  unsupported and no product execution path may use these bytes.
- Durable lifecycle Slices E1 through E3 are implemented locally and unwired. E4 and E5 remain.

## Priority 1: next parallel batch

These three tasks are independent enough to run in parallel and should use separate branches and
pull requests in their owning repositories.

### 1A. Capsule durable lifecycle Slice E4

Repository: `dills122/capsule-corp`.

- Replace the registered lifecycle `MemoryStore` path with the E3 fixed-store transaction port.
- Give every fake operation a stable `EffectID` and exact fake instance identity.
- Serialize `Drive`, `Recover`, and startup per `AttemptID` under one injected owner session.
- Exercise prepare/create/start/observe/stop/destroy intent, result, indeterminate, process-death,
  reopen, and reconciliation paths.
- Preserve `FakeBackend.CreatesGuest() == false`; add no process, runtime, backend, guest, consumer,
  IPC, content, or evidence path.

Exit: the migrated Slice C behaviors and E4 death/fault matrix pass against a newly reopened store
and newly constructed component. E5 remains a separate checkpoint.

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

Exit: a real governed branch and green draft PR exist in the fork, with exact upstream ancestry and
no experiment-directory or copied-registry dependency.

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

Exit: a real governed branch and green draft PR exist in the fork, and its release candidate can be
reconstructed from declared inputs without relying on the incomplete official prebuilt archive.

## Priority 2: evidence and contract closure

Start after the relevant Priority 1 branches stabilize:

1. Reconstruct the governed Deno, `rusty_v8`, binary, snapshot, and 22-entry root on an independently
   controlled Linux/arm64 builder; compare every declared output byte and retain provenance.
2. Decide the production owner and process topology for the ADR-0026 strip-only transformation,
   then implement the coordinated versioned original-TypeScript/emitted-JavaScript plan migration.
3. Assemble one immutable runtime bundle manifest from governed fork releases and the exact dynamic
   root; rerun physical-op, final-link, file-open, syscall-seal, descriptor, and restoration tests.
4. Run an independent review of the combined libkrun FD-native, direct-block-root, and console
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
