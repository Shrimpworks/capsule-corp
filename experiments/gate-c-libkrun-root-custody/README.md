# Gate C P0-1: immutable runtime-root custody

Status: **development-only retained experiment; PATCH-REQUIRED decision recorded 2026-08-02**.
Nothing here is product code, an admitted backend, or permission to execute user-supplied bytes.

Owner: Capsule core. Retain until an exact installed, signed/notarized App Sandbox bundle using a
narrow FD-native libkrun API has rerun this corpus and the result has been reconciled into the
canonical Gate C decision. Product packages must not import this directory.

## Defensive scope and question

This experiment defensively validates P0-1 using only owned repository fixtures, local processes,
the already pinned libkrun/libkrunfw build, locally cached OCI fixture images, and one owned local
Hypervisor.framework guest. It does not access another system, identity, credential, or data.

It separates three claims:

1. P0-1A: whether both stock libkrun 1.19.4 block consumers open `/dev/fd/N` as the exact inherited
   read-only object without returning to the original image pathname;
2. P0-1B: whether exclusive construction, closure of every writable alias/mapping, unlink, and
   post-finalization digest produce a frozen object; and
3. P0-1C: whether the whole topology resists the locally testable baseline same-user and crash
   cases.

## Exact inputs

- libkrun `v1.19.4`, commit `728df8125077d0db44265f6e997c72b81b65c015`, built with `BLK=1`
  and without `NET`;
- libkrunfw `v5.5.0`, embedded Linux `6.12.91`;
- the two existing Gate C source patches for firmware `@rpath` resolution and exact read-only root
  mount flags;
- Alpine fixture
  `alpine@sha256:14358309a308569c32bdc37e2e0e9694be33a9d99e68afb0f5ff33cc1f695dce`;
- Ubuntu builder
  `ubuntu@sha256:0e0a0fc6d18feda9db1590da249ac93e8d5abfea8f4c3c0c849ce512b5ef8982`.

The selected comparable root is ext4 without a journal. The journaled control is retained because
the guest's mounted block-device view replayed metadata in memory and therefore did not equal the
unchanged host backing bytes.

## Reproduce

Build from the already prepared pinned libkrun tree and locally available tools:

```sh
./experiments/gate-c-libkrun-root-custody/build.sh
./experiments/gate-c-libkrun-root-custody/prepare-root.sh
```

Run the local descriptor/race corpus and source audit:

```sh
./experiments/gate-c-libkrun-root-custody/source-audit.sh
python3 ./experiments/gate-c-libkrun-root-custody/local_custody.py
```

Run the owned unsandboxed guest path:

```sh
python3 ./experiments/gate-c-libkrun-root-custody/run_guest.py --timeout 60
```

Attempt the App Sandbox path. Exit 78 means this host reproduced the retained pre-main signing/
sandbox initialization limitation, not that libkrun or custody passed or failed:

```sh
python3 ./experiments/gate-c-libkrun-root-custody/run_guest.py --sandboxed --timeout 60
```

Focused verification, excluding the Docker/guest rerun by default:

```sh
./experiments/gate-c-libkrun-root-custody/verify.sh
CAPSULE_RUN_GUEST=true ./experiments/gate-c-libkrun-root-custody/verify.sh
```

Generated builds and raw reruns remain ignored under `.build/` and `.runs/`. Selected raw evidence
is retained under `evidence/2026-08-02/`.

## Decision boundary

The observed stock pathname route passed P0-1A and the unsandboxed full-guest construction path,
but P0-1C did not close because same-UID construction protection and the exact installed App
Sandbox runner could not be exercised without a valid signing identity. The retained decision is
therefore **PATCH-REQUIRED**, not PASS: introduce one raw-only FD-native libkrun API with explicit
descriptor ownership and no pathname fallback, then rerun this exact corpus on the final signed
installed bytes. The evidence does not justify rejecting libkrun for v0, and the patch does not
replace the still-required platform-protected construction-store test.
