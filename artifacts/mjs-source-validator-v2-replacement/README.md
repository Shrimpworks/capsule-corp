# Source Validator V2 replacement resource observation

Status: `PASSED` for the exact unprivileged public-SDK footprint-limit question. The product Source
Validator remains `BLOCKED`.

This directory retains one bounded resource-only probe used by the
[supported macOS replacement decision](../../docs/MJS_SOURCE_VALIDATOR_MACOS_PROFILE_REPLACEMENT.md).
It does not parse or execute JavaScript, use a signing identity, enter App Sandbox, open a Capsule
store, use a key, connect to a network, launch a runtime/backend/guest, or build product code.

The public macOS 26.5 SDK declares `task_set_phys_footprint_limit`. The probe asks the current
unprivileged process to set its own physical-footprint limit to 128 MiB before any controlled
allocation. On the observed host the call returned `KERN_NO_ACCESS`; no allocation ran. This makes
the interface unavailable as Capsule's unprivileged hard memory control on this host. It does not
prove behavior on another OS/build or under an entitlement Capsule has not selected.

Reproduce without retaining a binary:

```sh
xcrun clang -Wall -Wextra -Werror -O2 -mmacosx-version-min=13.0 \
  artifacts/mjs-source-validator-v2-replacement/testdata/footprint_limit_probe.c \
  -o /private/tmp/capsule-footprint-limit-probe
/private/tmp/capsule-footprint-limit-probe hold 192
```

Expected observed-host result:

```text
set-128 result=8 ((os/kern) no access) old_limit_mb=0
```

The source also contains a bounded `raise` mode for a future host on which the initial call
succeeds. It attempts to raise 128 MiB to 512 MiB before touching at most 512 MiB, so a successful
self-limit that a compromised child can later raise will also fail closed as a candidate. The
current observation stops before that branch.

No compiled binary is committed. See [`OBSERVATIONS.md`](OBSERVATIONS.md) for the exact environment,
hashes, result, and limitations.
