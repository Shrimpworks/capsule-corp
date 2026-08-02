# Gate C P0-3 local console correctness experiment

Status: **non-production, local dependency-reliability experiment**.

This experiment reviews the pinned libkrun 1.19.4 console implementation used by Capsule's
development backend candidate. It uses only the retained local source checkout and synthetic unit
inputs. It does not run user code, connect to an external system, or establish that P0-3 is closed.

## Question

Can the current console route handle fixed control messages, port state transitions, partial output,
and shutdown without panicking or waiting indefinitely, or does Capsule require a governed local
correctness patch before further integration?

## Reproduce

The default source path is the retained checkout used by the earlier Gate C work. Override it only
with another owned local checkout of the same exact commit.

```sh
CAPSULE_LIBKRUN_SOURCE=/private/tmp/capsule-libkrun-v1.19.4 ./verify.sh
```

`verify.sh` archives the exact pinned commit into a fresh temporary directory, applies the retained
candidate patch, checks the four changed Rust files, and runs the `krun-devices` library tests with
Cargo networking disabled. Temporary products are removed when the command exits.

## Retained scope

- [`patches/0001-console-correctness.patch`](patches/0001-console-correctness.patch) is an
  experimental downstream patch, not product code or an admitted dependency update.
- [`RESULTS.md`](RESULTS.md) records the observed test result, decision, and unresolved work.
- [`evidence/2026-08-02/verification.txt`](evidence/2026-08-02/verification.txt) retains the exact
  local verification summary.

Product packages do not import this directory. Any later use of the patch requires independent
review, complete profile composition, exact final-byte admission, and a resulting ADR/evidence
update.
