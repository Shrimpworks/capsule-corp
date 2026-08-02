# Gate C P0-0 governed Bun runtime-authority closure result

Date: 2026-08-02

Decision: **NO-GO. The exact Bun 1.3.14 source is present, but the required local source-build
inputs are unavailable. Stop the Bun governed-patch branch, investigate an alternate runtime, and
reconsider ADR-0003. Do not weaken the v0 authority contract, and do not admit runtime or backend
bytes.**

## Question and defensive boundary

Question: can the narrowest construction-level governed Bun 1.3.14 profile plus exact external
launcher/kernel enforcement preserve ADR-0003's dependency-free Bun first slice while making
subprocess/execve, FFI/native loading, inspector, Worker, macro/preload/config/environment-file,
and package-install/dynamic-resolution authority structurally unavailable?

This was authorized defensive, local-only research against Capsule's repository fixtures and the
retained exact Bun checkout. No user or third-party workload was run, no unrelated system was
accessed, and the prior stock Docker evidence was not treated as a security boundary.

## Exact retained inputs

| Input | Exact observation |
| --- | --- |
| Capsule task base | `22aa6d7f48019fb537c751b5495e6ee9d3a4e955` |
| Bun source | tag `bun-v1.3.14`, commit `0d9b296af33f2b851fcbf4df3e9ec89751734ba4`; clean retained checkout at `/private/tmp/capsule-gate-c-p0-0-bun-src-network` |
| Stock Linux runtime evidence | Bun SHA-256 `37141662ebed915a2ab89313156e455e2a1374395f5f6760d06407f49406f086` in the retained Linux/arm64 image |
| Local host | Apple M1 Max class arm64 host, macOS 26.5.2 build 25F84, Darwin 25.5.0 |
| Local bootstrap Bun | `1.3.14+0d9b296af`; macOS binary SHA-256 differs from the retained Linux runtime, as expected |

The exact read-only input check and output are retained in `check-inputs.sh` and
`evidence/2026-08-02/input-check.txt`.

## Fail-fast trigger

Bun's pinned `CONTRIBUTING.md` requires an existing Bun release, CMake, Ninja, LLVM 21.1.8, and the
documented platform dependency set; the build scripts also initialize Zig and dependencies. The
retained checkout had no `build`, `node_modules`, or `vendor/zig` tree. The host had a bootstrap Bun,
Go, Rust/Cargo, Ruby, GNU libtool, and `pkg-config`, but lacked `cmake`, `ninja`, `clang-21`,
`automake`, `ccache`, and GNU sed executables in `PATH`. The available Apple clang reported 21.0.0,
not the pinned LLVM 21.1.8 toolchain named by Bun's build instructions. Network-dependent toolchain
bootstrap was outside the available local inputs for this fail-fast campaign.

The task required an immediate NO-GO when the local source/build inputs were unavailable. No
governed source patch was authored after that trigger, no Bun build was attempted, and no mutated
binary was represented as evidence.

## Construction map and candidate mechanism

[`SOURCE_MAP.md`](SOURCE_MAP.md) records the pinned registry, process/native-loader,
module-resolution, inspector, Worker, configuration, and syscall entry routes. The smallest honest
candidate was not an API blacklist:

1. make a mandatory Capsule build profile omit every prohibited global, native module, compatibility
   module, loader, Worker, macro, discovery, and package-resolution route;
2. initialize indispensable JSC/Bun runtime state, then install a mandatory fail-closed Linux
   process/exec seal immediately before registered source evaluation; and
3. combine that with fixed empty environment/cwd/argv, an empty/no-exec workload staging design,
   and exact launcher/child descriptor manifests.

A purely external launcher filter is insufficient at the required timing: it must permit the first
`execve` that starts Bun but deny later replacement-image execution. A stateful seccomp broker would
add a new complex boundary. A runtime self-seal is narrower, but its `clone`/`clone3` and lazy-thread
compatibility can be determined only from the exact built runtime. Native loading has no generic
syscall filter compatible by assumption with Bun/JSC JIT behavior, so it additionally depends on
complete construction-level loader removal.

## Required measurements and mutations

| Required item | Result |
| --- | --- |
| Authored patch size | `0` lines; no patch was authored after the fail-fast trigger |
| Lower-bound maintenance surface | 21 pinned source units before generated outputs, build-profile wiring, tests, or a Linux self-seal; see `SOURCE_MAP.md` |
| Build and final binary digest | Not produced; local pinned build inputs unavailable |
| JIT/runtime breakage | Not measured; no patched binary exists |
| Descriptor inheritance | Not rerun; retained stock evidence still shows main VM and Worker reading inherited FD 3, so a closed exact manifest remains mandatory |
| Primitive restoration mutations | Not run; subprocess, execve, FFI, addon/plugin, SQLite loader, inspector, Worker, macro, preload/config/env, and package-resolution mutations require a built governed runtime |
| Structural closure argument | Not established. The source map identifies the required construction boundaries, but absence cannot be claimed without a governed diff, exact build, syscall trace, and restoration mutations |

## Decision and next action

`RUNTIME-001` remains unsupported and execution requiring it must refuse. The governed Bun branch
did not reach a testable construction because its exact local build inputs were unavailable, and
its lower-bound maintenance surface spans native registry, C++, Zig, generated JS modules,
resolver/configuration, JSC lifecycle, and Linux syscall enforcement.

The exact next decision is to investigate an alternate runtime against the same prohibited-power
contract and then update or supersede ADR-0003 based on retained evidence. The parent orchestrator
must not reinterpret this environmental NO-GO as authority to weaken the profile or admit current
Bun/libkrun bytes.

## Limitations

- This campaign did not install or download missing build dependencies.
- It did not compile Bun, author a speculative unbuildable patch, or run the local Docker oracle.
- The 21-file surface is a reviewed lower bound, not a line-count estimate or proof that a safe
  patch is impossible.
- The syscall-seal analysis is a construction constraint, not validated seccomp/LSM evidence.
- The result decides this governed-Bun campaign under the explicit fail-fast rule; it does not
  select or validate the alternate runtime.

## Verification

Focused evidence checks passed:

```sh
sh -n experiments/gate-c-bun-runtime-authority/governed-closure/check-inputs.sh
shellcheck experiments/gate-c-bun-runtime-authority/governed-closure/check-inputs.sh
./experiments/gate-c-bun-runtime-authority/source-inventory.sh \
  /private/tmp/capsule-gate-c-p0-0-bun-src-network
./experiments/gate-c-bun-runtime-authority/governed-closure/check-inputs.sh \
  /private/tmp/capsule-gate-c-p0-0-bun-src-network
```

Repository-required verification passed with Node 22.22.1, pnpm 10.28.2, and Go 1.26.5:

```sh
pnpm install
pnpm check
pnpm lint
pnpm test
pnpm verify:schemas
go test ./...
go vet ./...
go build ./...
```

The first ambient-shell `pnpm` invocation selected Node 16 and failed before installation; rerunning
through the repository-compatible Node 22.22.1 completed from the existing pnpm store without
downloads. This tooling correction did not change the Bun source-build NO-GO: Bun's separate pinned
build prerequisites and initialized build inputs remained absent.
