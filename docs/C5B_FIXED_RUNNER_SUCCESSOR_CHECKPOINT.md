# C5b fixed-runner successor checkpoint

Date: 2026-09-05

Status: `PASSED` for reconciliation and fresh local verification of retained C5b11
construction/static evidence. Parent controlled C5b execution: `BLOCKED`.
Owner-only internal alpha: `IN_PROGRESS — TRENDING_GOOD`; runtime/profile and product admission:
`BLOCKED`.

## Immutable target and historical disposition

The August 21 plan omitted two experiment merges from August 18. Do not rebuild the already-retained
no-run successor or promote the superseded C5b10 candidate:

| Scope | Exact identity | Disposition |
| --- | --- | --- |
| C5b10, experiments PR #30 | Merge `ecc3e5efb835931d2d2113d1bc20831a35aba8b4`; candidate `6eb030130734882de4529e647a5a0ac29af362f6` | Historical construction only; not accepted evidence after review found stale attempt/profile binding and incomplete fault reconciliation. |
| Earlier C5b11 review targets | `d4a805ab6fc6fb700d06f57896a2775680755d0f`, then `5a671198a61280ce343e2ba03787430da27fc1b7` | Preserved ancestors; their findings were addressed in normal descendant commits. |
| C5b11, experiments PR #31 | Reviewed head `27b9011fb80edc0f23c11b3f3fa76d00cebc2365`; merge `f206e4ef2cd326ee74e5b7b2739c62efe6da7d6d` | Construction/static evidence `PASSED`; the PR publication records C5b-S5 review as `PASSED` / `Ready` for that exact head. |

The [immutable C5b11 packet](https://github.com/Shrimpworks/capsule-experiments/tree/f206e4ef2cd326ee74e5b7b2739c62efe6da7d6d/experiments/typed-guest-transport-c5b11-bound-fault-convergent-no-run-successor)
is the current no-run input. Its README, results, handoff, and review packet intentionally retain
their pre-review checkpoint wording. The later review disposition is reported by the
[PR #31 publication](https://github.com/Shrimpworks/capsule-experiments/pull/31), not by rewriting
those archived bytes. That PR has no separate GitHub review/comment report; the publication's
review summary and this local re-verification are distinct evidence sources. Neither grants
execution authority or product admission.

## What the candidate closes

The fixed runner binds the exact 100,663,296-byte C5b7 root and alone imports the 13 closed libkrun
symbols. The Supervisor driver imports zero libkrun symbols and 24 closed provider symbols. Its
entry accepts only the fixed registration ID; source/input writes and writer closure precede the
start byte. Attempt-plan payload bindings, distinct recovery and durable-resume cursors,
ambiguous-spawn reconciliation, one-shot teardown, and stored completion replay are checked in the
literal no-run model.

| Retained object | SHA-256 |
| --- | --- |
| Attempt runtime profile | `829bdd048210c14d67f4cfcb659c39db69fe5ed2ff4edb74f3f2d9f3c869f82d` |
| Attempt plan | `891359ad03c420b658f0ce66769fd9996eae0022bdd0ea92a3884a8c7723bf29` |
| Fixed runner object | `e3d249c29885e9f1b3d961296835830d1b735d246d5b58d5656b9a7dbb46a65b` |
| Supervisor driver object | `468b76662f87bd9e9599b99e9c417e2575b8805fc1a44073952238f6d1388342` |

## Fresh verification on 2026-09-05

Defensive scope: exact archived C5b11 sources, fixtures, and predecessor bytes in an isolated local
archive checkout and disposable build/mutation directories. No candidate object or dylib was
linked, loaded, or invoked; no provider, runner, libkrun/HVF, VM, guest, signing, Keychain, installed
service, or product-state operation was performed.

Commands run from the C5b11 experiment directory:

```sh
node scripts/generate-bindings.mjs --check
node scripts/generate.mjs --check
node --test scripts/verify-profile.test.mjs
node scripts/verify.mjs
node scripts/test-mutations.mjs
```

Environment: Node.js 22.22.1 and Apple clang 21.0.0 (`clang-2100.1.1.101`), arm64 Darwin 25.6.0.
The experiment directory is byte-identical between reviewed head and merge; the archive checkout
remained clean after verification.

All passed: 34-file inventory, 22 tests, 95 restored-invalid mutations, 65 primary reconciliation
cases, 50 recovery/failure crossings, 11 reopen/resume cases, and performed effects `NONE`.
`sh scripts/build.sh` was also run in a disposable copy: its two independent unlinked-object builds
matched each other and both exact retained object hashes above. This is same-host reproduction;
it does not establish cross-host equality or behavior of absent providers.

## Next bounded implementation

At this C5b11 checkpoint, provider implementations and their provenance were absent.
The later [C5b12 native transport checkpoint](C5B_NATIVE_TRANSPORT_PROVIDER_CHECKPOINT.md)
adds seven actual pipe providers and local fixture evidence; 17 lifecycle/root/store providers
and complete composition remain blocked. The original next-slice requirement was to map
the 24 declarations in `source/supervisor_effect_abi.h` to fixed implementation owners and the
retained recovery oracle. Each implementation must freeze its exact construction/test boundary
before adding native effects. Preserve one libkrun owner, registration-only entry, independently durable restart
cursors, completion-before-delivery, and authoritative absence before root removal.

Completion of that slice requires actual provider source, closed inputs/imports and provenance,
reproducible construction, failure-sensitive verification, and independent review of its exact
new composition. Do not count interface declarations, a test double, or unconditional refusal
stubs as implemented platform behavior. Any local effect tests need their own explicitly bounded
fixture/process scope; guest execution remains separately authorized only after a complete
reviewed immutable composition exists.

Preferred-form libkrunfw/kernel source compliance, raw v19/v27 evidence recovery or separately
versioned replacement, cross-host evidence, installed composition, runtime/profile admission,
and product admission remain `BLOCKED` in their existing scopes. No ADR lifecycle or control
evidence state changes in this reconciliation.
