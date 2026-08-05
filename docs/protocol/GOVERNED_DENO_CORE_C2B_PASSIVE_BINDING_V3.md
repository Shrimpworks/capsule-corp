# Governed `deno_core` C2B passive successor contract v3

Status: `PASSED` for the passive successor contract only. Parent governed-runtime work remains
`IN_PROGRESS — TRENDING_GOOD`. Fixed-owned-guest eligibility and runtime/profile admission remain
`BLOCKED`; `RUNTIME-001` and `VMM-001` remain unsupported.

This record closes the smallest passive contract needed before a separately authorized run of one
fixed benign owned guest. It creates no process, runtime consumer, adapter, backend, VM, guest,
credential, signature, publication, release, or admission effect.

## Immutable object

The canonical fixture is
[`passive-binding-v3.json`](../../schemas/conformance/c2b-governed-deno-core/passive-binding-v3.json).
Its exact identity is:

- bytes: `18,357`;
- fixture SHA-256: `d72327bba369484a56db7d543a32e8bbd4eac403230ac65d63709ac3ba3bbdfb`;
- composed passive-contract SHA-256:
  `8b1ec936a7b56370716d28557125e46866dea8f21a149704a01f251a0dddbcc1`.

The composed digest is SHA-256 over the exact fixture after replacing only
`composedProfile.contractDigestSha256` with 64 ASCII zeroes. It binds every role, source identity,
artifact disposition, descriptor, device, runtime control, supported resource value, transport
cap, teardown rule, blocker, and zero effect. It identifies this passive contract; it is not a
runnable materialization, signature, review record, backend-validation record, or admission.

V3 consumes C1, C2A, C2B v1, and C2B v2 by their exact byte lengths and digests. All predecessor
bytes remain immutable. V3 contains no JSON `null`: missing runnable identities are typed blockers,
and unsupported resource fields are absent.

## Governed source and artifact reconciliation

No mutable branch name is authority. V3 binds these immutable accepted source records:

| role | upstream or anchor | accepted commit | accepted tree |
| --- | --- | --- | --- |
| Deno | `14eea3160ae5834476aa3b9d317b8d41d991b982` | `3fa21d1ae7705ab4bcb4bc98955f25301b20122a` | `6060cb0eb4cd3395a4c141f054634968744617d2` |
| `rusty_v8` | `d305e6afa7736f6e298c30ae6646f7709ee9382b` | `d09221062280ae1675fe26c53c3f43871aae2055` | `2632901e6e7e9ac88662756ceb658d4e3e49fceb` |
| libkrun | `728df8125077d0db44265f6e997c72b81b65c015` | `7432eda5a49220976b0167005aa43ee622f9d632` | `7671440cfbafa58fe20aebf8d4deb2a843ebe346` |

The accepted Deno change is a governance-only descendant of the runtime-source merge used by the
retained build. The accepted `rusty_v8` run retained unsigned artifact `8941046558` with archive
digest `6073159d...3473`; its static archive remains byte-equal to C1. Neither fact admits those
bytes.

The retained 4,426,736-byte libkrun dylib at SHA-256 `f8e05177...ae52` was built from superseded
source commit `cf0333cd...3ecc`, tree `ffa4131d...3fee`. Accepted commit `7432eda5...d632` changes
runtime console validation. Its CI compiled a dylib but retained no artifact or digest. V3 therefore
preserves the old dylib as historical evidence only and blocks selection until exact bytes built
from the accepted commit/tree are retained and verified.

## Boot and trusted guest roles

The non-EFI build uses only libkrun feature `blk`; default features and `init-blob`, `net`, `efi`,
GPU, sound, input, and TEE features are absent. The retained `libkrunfw.5.dylib` is the sole runtime
boot-kernel carrier. Libkrun resolves `@rpath/libkrunfw.5.dylib` and obtains its kernel through
`krunfw_get_kernel`. The extracted Linux 6.12.91 image at SHA-256 `b50a4165...22d` is derived
evidence, not a second runtime input. A separate firmware role is inapplicable. The runner must not
call `krun_set_kernel` or `krun_set_firmware`, and no caller path may select either role.

The retained raw ext4 root is 134,217,728 bytes at SHA-256 `390a4786...6798`. It is the only block
device and is attached from exact FD 4 through the governed raw read-only root API as `/dev/vda`,
then mounted ext4 read-only, `nosuid`, and `nodev`. The trusted PID 1 and launcher are fixed at:

| role | guest path | bytes | SHA-256 |
| --- | --- | ---: | --- |
| init | `/usr/local/libexec/capsule-init.krun` | 930,144 | `4f4f2c8b...85cd` |
| launcher | `/usr/local/libexec/capsule-launcher` | 995,920 | `fd255394...32d` |

The init remounts root, mounts only bounded `devtmpfs` and `proc`, closes inherited descriptors,
opens the exact guest descriptor set, and execs the launcher with an empty environment. The
launcher verifies complete fixed source/input frames, withholds completion authority from the
runtime child, waits and drains it, validates the fixed result and seal marker, proves the child
terminated, and writes the completion frame last.

## Final host-runner role

V3 resolves the final role as one per-attempt App-Sandboxed VMM process, created and owned only by
the Execution Supervisor. It is not a daemon, service, or privileged helper. The Supervisor owns
durable identity recording, the one-byte `G` start authorization, continuous drains, completion
interpretation, cancellation, signaling, and authoritative absence. The runner owns only exact
preflight, fixed libkrun configuration, and `krun_start_enter`; its exit status is lifecycle
evidence, never workload success.

The fixture freezes the ordered libkrun call contract. Critical points include explicit
`krun_disable_implicit_console`, `krun_disable_implicit_init`, and
`krun_disable_implicit_vsock`; the raw-root FD API; one multiport console and exactly three fixed
ports; fixed kernel console `hvc0`; fixed init path, one-element argv, three-entry init environment,
and `/` workdir; then the exact `G` plus EOF handshake before `krun_start_enter`. All path-selecting
kernel, firmware, legacy disk, root, virtiofs, vsock, network, GPU, and sound calls are forbidden.

The role is resolved, but no final runner artifact exists. The 34,864-byte build-only preflight at
SHA-256 `4d56480b...0922` calls no libkrun API and cannot fill that role. V3 blocks guest eligibility
until final runner bytes, digest, call audit, and mutation evidence are retained.

## Exact descriptors, ports, and devices

The runner receives only FDs 0 through 7: null stdin; Supervisor-drained stdout/stderr; start
control; read-only root; source input; approved input; and completion output. It closes FD 8 and
above before any libkrun call. Guest init realizes only FDs 0 through 5. The runtime child receives
only null stdin and bounded stdout/stderr pipes; completion FD 5 is `CLOEXEC` before child start.

One virtio-console device exposes ports 0 `capsule.source` at `/dev/hvc0`, 1 `capsule.input` at
`/dev/vport0p1`, and 2 `capsule.completion` at `/dev/vport0p2`. The exact virtio inventory is
balloon, RNG, that multiport console, and the read-only root block device. Additional console,
network, vsock, virtiofs/`NullFs`, writable or additional block, GPU, sound, and input devices are
forbidden. Implicit console, init, and vsock are explicitly disabled; without a vsock device, TSI
is absent. Host sockets and live host paths are absent.

## Runtime, resources, and teardown

Only the fixed 103-byte benign source, 36-byte canonical input, and 35-byte expected JSON result
are in scope. The 68,496,520-byte fixed-fixture runtime is SHA-256
`e781a90236cdf1272a9a16189c6be033164fa25a5aa9e52376ef998982ec0a77`. It uses `--jitless` and
`--random-seed=42`, no module loader, no
extensions or inspector, the three bootstrap-required ops only, and an internal fixed main-module
URL. It calls V8 context `set_allow_generation_from_strings(false)` and verifies the result. Static,
re-export, dynamic, eval/`Function`, WebAssembly, filesystem, network, and package loading are
forbidden. The retained TSYNC seccomp seal denies clone, exec, socket, and executable-memory paths.

The only supported resource fields are exactly one vCPU, 256 MiB guest RAM, 1,000 ms wall time,
and concurrency one. The fixed guest has no scratch device. CPU-time, exact total host/VMM memory,
and scratch-maximum fields do not exist in v3; adding any unknown resource field refuses. Transport
and diagnostic caps remain exactly those in C2A.

At wall expiry or cancellation, the Supervisor acts independently of guest cooperation and drains.
The 200 ms grace cannot support success. It revalidates PID, start time, UID/GID, executable path,
and code identity before `SIGKILL`, then requires authoritative absence within 1,000 ms and no more
than 1,200 ms after the action. Identity mismatch is unresolved: no signal and no success. Ordinary
success also requires source/input/runtime integrity, valid commit-last completion, bounded result,
terminal runner lifecycle, child-tree absence, authoritative runner absence, and no required
cleanup.

## Remaining gate

Fixed-owned-guest eligibility remains `BLOCKED` on four exact steps:

1. retain a dylib built from accepted libkrun commit `7432eda5...d632` and verify its tree;
2. construct, retain, and review the final runner bytes, digest, call sequence, and mutations;
3. create a new versioned materialization binding those identities to this contract; and
4. obtain separate owned-disposable-guest authorization naming that new profile digest.

V3 authorizes none of those actions.

The later [host-runner source contract v1](GOVERNED_DENO_CORE_C2B_HOST_RUNNER_SOURCE_V1.md) is
`PASSED` only for dependency-free passive C17 call-plan bytes and local mutation verification. It
does not fill step 2: the accepted header/current-source dylib are not retained locally, no exact
ABI runner was compiled, and no final runner bytes or digest exist. V3 remains immutable and fixed-
owned-guest eligibility remains `BLOCKED`.

The later [v4 materialized successor](GOVERNED_DENO_CORE_C2B_MATERIALIZED_PROFILE_V4.md) is
`PASSED` for steps 1 through 3 only: it retains exact current-source libkrun and final-runner bytes,
an independent ABI audit, and a new composed digest without loading or executing either artifact.
Step 4 remains `BLOCKED`; separate guest authorization must name that exact v4 digest. V3 bytes and
this historical gate remain unchanged.
