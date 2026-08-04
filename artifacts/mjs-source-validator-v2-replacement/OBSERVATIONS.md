# V2 replacement public memory-interface observation

Date: 2026-08-04

## Work item

```text
Work item: unprivileged public-SDK physical-footprint limit for a future macOS Source Validator
Status: PASSED
Scope: one fixed C probe on macOS 26.5.2 arm64; no JavaScript, signing, sandbox, store, key,
  network, runtime, backend, or guest
Evidence or reason: the public SDK interface returned KERN_NO_ACCESS for the process's own task
  before any allocation.
Remaining work: this negative answer does not provide a memory control; architecture must select a
  supported hard limit or explicitly accept a quantified reactive monitoring policy.
Next action: keep product V2 replacement BLOCKED pending that decision.
Parent status: Product Source Validator V1-V5 and downstream M2/S1 are BLOCKED.
```

## Environment

- macOS 26.5.2 build `25F84`
- Darwin 25.5.0 `RELEASE_ARM64_T6000`
- Apple arm64 host
- macOS 26.5 SDK at
  `/Applications/Xcode.app/Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs/MacOSX26.5.sdk`
- Apple clang 21.0.0 (`clang-2100.1.1.101`)
- probe source: 1,730 bytes, SHA-256
  `593ba843c40b00c9d9297d35a7388e4cf8a71722a62caae8529f733db5b4d92d`
- temporary Mach-O from the recorded command: SHA-256
  `8fa95aa3349bc25e1a1915bd0e22ff64f89e08b2f211b0b6307b1a387d845ffe`

The temporary binary hash is an observation, not an enrolled artifact or reproducible-build claim.
The binary was written only to `/private/tmp` and is not retained.

## Public interface and command

The installed public SDK declares `task_set_phys_footprint_limit` in `mach/task.defs` and
`mach/task.h` with the description “Change the task's physical footprint limit (in MB).” The
probe calls it only for `mach_task_self()`.

```sh
xcrun clang -Wall -Wextra -Werror -O2 -mmacosx-version-min=13.0 \
  artifacts/mjs-source-validator-v2-replacement/testdata/footprint_limit_probe.c \
  -o /private/tmp/capsule-footprint-limit-probe
/private/tmp/capsule-footprint-limit-probe hold 192
```

Observed result and exit:

```text
set-128 result=8 ((os/kern) no access) old_limit_mb=0
exit 65
```

No allocation branch ran. No signing identity or entitlement was requested or used.

## Related SDK observations

- `RLIMIT_RSS` is a source-compatibility alias of `RLIMIT_AS`, not a separate control.
- `RLIMIT_DATA` is the data-segment/`sbrk` limit described by Apple's `setrlimit(2)` page; it is
  not a total mapping, allocator, resident, or address-space ceiling.
- `RLIMIT_MEMLOCK` limits locked memory only.
- `proc_pid_rusage` is public and reports current/lifetime physical footprint for live or zombie
  processes. It is an observation interface, not a setter.
- The SDK exports several spawn/jetsam-related symbols that have no public declaration in the SDK
  headers. They are private implementation surface and were not called.

## Inference and limitation

The observed unprivileged child cannot use this interface to create its own hard footprint limit
on this host. That result joins, but does not replace, the retained V2 `RLIMIT_AS=EINVAL`
observation. It does not establish the absence of every kernel mechanism, behavior on another
supported OS, or behavior available to Apple/system processes. Capsule rejects private APIs and
does not infer product support from an exported symbol.

Parent sampling through `proc_pid_rusage` remains possible, but any kill threshold is reactive:
the child may exceed it between samples. The supported macOS replacement decision therefore keeps
the current exact-memory gate `BLOCKED`.
