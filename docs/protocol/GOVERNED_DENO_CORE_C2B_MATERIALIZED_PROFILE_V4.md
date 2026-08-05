# Governed `deno_core` C2B materialized host-runner profile v4

```text
Work item: C2B current-source libkrun and final host-runner materialization
Status: PASSED
Scope: exact-source controlled compilation, exact ABI review, unsigned Mach-O retention, and
  passive composed-profile verification only; no runner or libkrun execution, HVF call, VM,
  guest, arbitrary workload, signing, installation, publication, product wiring, or admission
Evidence or reason: two clean package-scoped libkrun builds and two clean runner builds produced
  byte-equal unsigned artifacts. Exact known answers, imports, exports, loader paths, FDs 0–7,
  close-from 8, three console ports, implicit-device disable calls, fixed call order, no replacement
  authority, boot-role resolution, and Supervisor-owned teardown pass retained Go verification and
  mutations.
Remaining work: a separately authorized owned-disposable-guest task must name the exact composed
  profile digest before loading libkrun or creating a guest. Runtime/profile admission, signing,
  installation, product wiring, release review, and complete P0 evidence remain BLOCKED.
Next action: stop before guest execution and call back for separate authorization.
Parent status: governed-runtime engineering is IN_PROGRESS — TRENDING_GOOD; RUNTIME-001 and
  VMM-001 remain unsupported.
```

This immutable successor consumes, but does not alter or reinterpret, C1, C2A, C2B v1/v2/v3, or
the 3,996-byte host-runner source contract v1. It closes v3's build/static blockers only. It is not
a `BackendValidationRecord`, installed profile, runtime release, signature, or permission to run a
guest.

## Immutable object and artifacts

The canonical directory is
`schemas/conformance/c2b-host-runner-materialized-v4`. The materialized profile is 10,301 bytes at
SHA-256 `198688bacd50aaee4f57b4cd7c56cea6b939c10aa220fbbeba7d315de820d1fd`.
Its composed-profile digest is
`e390085caaaba73ebc19f95bc9871305e4f9268c2283d7394133fa4491f4ba82`, calculated over the exact
JSON after replacing only `composedProfile.digestSha256` with 64 ASCII zeroes.

| artifact | bytes | SHA-256 |
| --- | ---: | --- |
| accepted `libkrun.h` | 54,658 | `dce44d1d70ab770b1089e57646e025281a4137fe5052b9dd8eaefb80c01a1bd8` |
| ABI audit source | 2,512 | `419256ea91de9b5e5323e1f1d6d42afb0a5fa85a8835d0d0404734af0ee92356` |
| unsigned current-source `libkrun.1.dylib` | 4,393,448 | `055d9d18dc964fec4aba21948c4a344cb7a51cb48a2c70017484b718eae12f9f` |
| final C17 runner source | 7,917 | `5a5560fa667390253bf504d7c045fcbcc304fa5829b22a8acf1fff00a8e37eb9` |
| unsigned final runner | 100,488 | `a30e3f7cba5f480b6e164536854749b5e1ba3349f20af6c9c8e5d2590bffe1ad` |

The dependency-free verifier at `internal/execution/hostrunnermaterialized` bounds and checks the
profile and every retained artifact, recomputes the composed digest, parses both Mach-O files,
checks runner imports, checks required libkrun exports, rejects code signatures, audits source
ordering and forbidden calls, and proves no other Go package imports this unwired verifier.
One-byte artifact mutations plus authority, FD, port, firmware, and handshake-order mutations
refuse.

## Exact governed source and controlled build

No mutable default branch supplied authority. The libkrun source was fetched directly by accepted
commit `7432eda5a49220976b0167005aa43ee622f9d632`; its tree was independently read back as
`7671440cfbafa58fe20aebf8d4deb2a843ebe346`. Its exact upstream anchor remains
`728df8125077d0db44265f6e997c72b81b65c015`. The accepted Deno and `rusty_v8` identities remain
the exact v3 values.

The accepted libkrun `Cargo.lock` is 45,876 bytes at SHA-256
`9d5dc785636a264794a396ab478821c4ed33acae91650db8d72e8a35733f288c`. Two clean source archives
were built with the package-scoped equivalent of:

```sh
cargo build --release --locked --offline -p libkrun --no-default-features --features blk
```

Build paths were remapped to `/usr/src/libkrun`. The two raw outputs were byte-equal at 4,427,920
bytes and SHA-256 `86d94dbd4f01e80c6c9058966a2d7fe6e823738291447d4e72e648a4f7f4c9e3`.
The linker supplied identity-free ad-hoc signatures; removing those signatures, without applying
any signature or identity, produced the retained byte-equal dylibs above. The dylib is arm64,
minimum macOS 11.0, SDK 26.5, install name `libkrun.1.dylib`, with only Hypervisor.framework,
`libiconv.2.dylib`, and `libSystem.B.dylib` as dynamic dependencies.

An attempted `-Wl,-no_uuid` recipe is `NO_GO` for this exact build environment because Rust build
scripts could not execute as Mach-O files without `LC_UUID`. It produced no retained candidate.
The accepted recipe instead retains deterministic UUID-bearing bytes and removes only the
identity-free linker signature.

The build environment was macOS 26.5.2 build 25F84, Xcode 26.6 build 17F113, SDK 26.5, Rust
1.93.1, Cargo 1.93.1, and LLVM 21.1.8. Same-host clean-directory equality is not independent-builder
reproduction, supported-minimum-OS evidence, or release provenance.

## Independent ABI review

The exact accepted header and implementation were reviewed before the runner was compiled. The
retained C17 audit uses type assertions for every runner call and passes `clang -fsyntax-only`
against that header. The final runner imports only these libkrun functions:

- context creation and exact vCPU/RAM configuration;
- explicit disablement of implicit console, init, and vsock;
- read-only raw-root FD attachment and exact remount selection;
- one multiport console and three ordered ports;
- kernel console, workdir, and fixed exec selection; and
- `krun_start_enter`.

Three implementation details are closed explicitly:

1. The header exposes `krun_create_ctx` with an unspecified parameter list while the implementation
   takes zero parameters. The caller passes zero arguments only.
2. `krun_set_exec` forms 4,096-pointer Rust slices before scanning for null terminators. The runner
   therefore supplies two static, fully allocated 4,096-slot arrays with early null terminators;
   short automatic C arrays would not satisfy this implementation contract.
3. `krun_add_console_port_inout` returns success zero, not the allocated port ID. Port IDs 0, 1,
   and 2 are bound by fixed append order on console ID zero.

The retained runner's Mach-O import table agrees with the static review. It contains none of the
forbidden kernel, firmware, legacy disk/root, virtiofs, vsock, network, GPU, sound, `open`, socket,
spawn, exec, or dynamic-loader imports. The library continues to export broader upstream APIs;
the unwired runner does not import them. That is not product-level library surface reduction.

## Boot-kernel and firmware role

Exact accepted source resolves the disputed roles without inventing an identity. This non-EFI,
`blk`-only libkrun build loads `@rpath/libkrunfw.5.dylib` and obtains its kernel from
`krunfw_get_kernel` when no external kernel, kernel bundle, or firmware was configured.
The exact retained libkrunfw identity remains 24,339,104 bytes at SHA-256
`0b14f4b8005dafd33c38df5935b9efdb6381c724224b3967ba1cecbecf10b6e9`.
The extracted Linux 6.12.91 image remains derived audit evidence only: 24,117,248 bytes at SHA-256
`b50a4165215d5d897ab3614606a2105756cf8f2b2510cbceda9dc06057a5622d`.

Libkrunfw is the sole runtime boot-kernel carrier. There is no separate firmware input, path, or
identity in this profile. The runner calls neither `krun_set_kernel` nor `krun_set_firmware`.

## Final runner and authority boundary

Two clean runner builds were byte-equal after replacing libkrun's install reference with
`@rpath/libkrun.1.dylib` and removing only the linker ad-hoc signature. The retained runner is an
arm64 Mach-O with minimum macOS 14.0, SDK 26.5, and `@executable_path` rpath. It was never launched.

It accepts no argument, AttemptID bytes, plan, profile, image, mount, guest path, backend flag, or
replacement configuration. After removing only macOS's `__CF_USER_TEXT_ENCODING` injection, it
requires an empty environment. It validates exactly FDs 0 through 7 and closes every FD from 8
through the process limit before its first libkrun call:

| FD | exact role |
| ---: | --- |
| 0 | character-device, read-only stdin |
| 1–2 | Supervisor-drained, write-only FIFOs |
| 3 | Supervisor start-control, read-only FIFO |
| 4 | unlinked regular read-only mode-0400 raw root |
| 5 | source, read-only FIFO |
| 6 | approved input, read-only FIFO |
| 7 | completion, write-only FIFO |

FD 4 must be unique, unlinked, nonzero device/inode, exactly 134,217,728 bytes, and match SHA-256
`390a4786a20d45f1c691ec8c203f84f5e9d372a30e98f867cc8309a144ca6798` under bounded `pread`.
It becomes read-only raw `/dev/vda`, remounted ext4 `ro,nosuid,nodev`.

The fixed call order creates one three-port console: `capsule.source` uses `(5,-1)` as port 0,
`capsule.input` uses `(6,-1)` as port 1, and `capsule.completion` uses `(-1,7)` as port 2. It then
sets `hvc0`, `/`, and the fixed init argv/environment, emits `R`, requires exactly `G` followed by
EOF, and only then calls `krun_start_enter`. Implicit console, init, and vsock are disabled before
device attachment.

Only the Execution Supervisor may create one runner for one committed `AttemptID`, write the start
authorization, continuously drain, interpret completion, cancel, signal, and prove authoritative
absence. Teardown stays external: revalidate PID, start, UID/GID, executable path, and code identity
before `SIGKILL`; require authoritative absence within 1,000 ms and no more than 1,200 ms after the
action. Identity mismatch is unresolved: no signal and no success. Runner exit is lifecycle
evidence, never workload success.

## Dependency and admission boundary

This slice reuses the accepted governed libkrun line; it adds no new product dependency or
primitive. Exact source, tree, header, lockfile, features, dynamic closure, and unsigned artifact
are retained. The governed-fork adoption recommendation still does not confer product admission.
Independent source and license review of the full transitive closure, advisory review, two-builder
reproduction, signed/notarized installed composition, upgrade/removal ownership, libkrunfw source
and distribution obligations, minimum-OS evidence, and the remaining P0 restoration corpus stay
open.

## Stop boundary

The build/static materialization slice is `PASSED`. The first fixed benign owned guest is
`BLOCKED` until a separate task explicitly authorizes an owned disposable guest and names composed
profile digest `e390085caaaba73ebc19f95bc9871305e4f9268c2283d7394133fa4491f4ba82`.
That later task must reverify every exact byte before loading libkrun. Product runtime/profile
admission remains `BLOCKED`; parent governed-runtime engineering remains
`IN_PROGRESS — TRENDING_GOOD`; `RUNTIME-001` and `VMM-001` remain unsupported.
