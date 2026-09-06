# C5b15 immutable runner reconciliation checkpoint

Date: 2026-09-06

Status: exact-input static reconciliation, refusal tests and independent review
`PASSED` / **Ready** (instance 2 of 3).
Parent owner-only internal alpha: `IN_PROGRESS — TRENDING_GOOD`.
Immutable runner composition, installed lifecycle, controlled guest execution and
product admission: `BLOCKED`.

## Merged inputs and scope

C5b14B archive [PR #35](https://github.com/Shrimpworks/capsule-experiments/pull/35)
merged at `bb33895d82853c35b99e53b234c464dd55e6685a`; canonical
[PR #360](https://github.com/Shrimpworks/capsule-corp/pull/360) merged at
`2c30faa8c0e88dd43f289bbea18f58232d02cf23`. The reviewed C5b14B source/evidence
pin remains `294559e78b410f4b25e0e25fcf373ec0c7d1cfb5`; merge metadata does not
replace it. The [native durable-owner checkpoint](C5B_NATIVE_DURABLE_OWNER_CHECKPOINT.md)
continues to report a completed benign-fixture slice.

The [C5b15 archive](https://github.com/Shrimpworks/capsule-experiments/tree/7c22ad119c21d44e74d0a525f96ddcee915330f8/experiments/typed-guest-transport-c5b15-runner-reconciliation)
at `7c22ad119c21d44e74d0a525f96ddcee915330f8` retains 13 exact input copies with repository, commit, full
path, byte count and SHA-256, an offline Node audit, deterministic report and refusal
tests. Inputs include C5b11 merge `f206e4ef2cd326ee74e5b7b2739c62efe6da7d6d`,
reviewed C5b14B and canonical C2A/ADR-0041 requirements at the merge above.

Defensive scope: read those exact repository sources/JSON using local Node tools
and owned fixture buffers/directories. Copied C/Go/generator source is data. No
candidate is compiled, linked, loaded or executed; no benign child, VMM, guest,
endpoint, signal, installed store/service, signing, Keychain or credential is used.
No product consumer, schema, ADR lifecycle or control-evidence maturity changes.

## Concrete composition differences

| Dimension | C5b14B benign fixture | Retained C5b11 candidate | Consequence |
| --- | --- | --- | --- |
| Root | 65,536 bytes | 100,663,296 bytes (96 MiB) | Exact root/runner/source/plan identities require a versioned composition |
| argv zero | `fixture-runner` | `capsule-c5b11-fixed-runner` | Existing native launch does not satisfy candidate preflight unchanged |
| Providers | 24 `_c5b14b_` exports | 24 `_c5b11_` imports | Same logical ABI lineage, different compiled symbols and generated bindings |
| Attempt/runtime profile | `f1f6cb0eb8d5f663a5b90ce77fdb1e5830eac0c47e952bb587d4dfaa7d8a137c` | `829bdd048210c14d67f4cfcb659c39db69fe5ed2ff4edb74f3f2d9f3c869f82d` | Old frames and store identity cannot stand in for a newly composed candidate |

The exact path **substitute unchanged C5b14B providers into retained C5b11** is
`NO_GO`. This abandons that direct substitution only. It neither invalidates the
completed benign experiment nor declares the broader runner workstream abandoned.
Missing evidence and unimplemented successor mechanisms below remain `BLOCKED`.

## Timing finding and remaining gates

C5b14B starts its 1,000-ms setup clock at endpoint creation. That budget includes
executable/root work and durable spawn publication. Its teardown function calls the
durable before-teardown gate, then reads monotonic time and starts a new 1,000-ms
cleanup clock. Therefore that cleanup clock excludes time waiting in the gate.
It does not alone establish C2A's bound from initial action to authoritative absence.

This is an exact-source ordering observation, not a measured storage hang or a
claimed cause of the earlier unexplained setup refusal. The latter remains
unclassified. Passing the benign corpus does not establish timing/load tolerance.
C5b15 reads its retained results as testimony and does not report a fresh native run.

Restart lookup preserves four cases: already-fenced state returns the existing
cursor; unfenced state without spawn intent is fresh; unfenced consumed intent
with completion publishes fenced resume 22 for replay; unfenced consumed intent
without completion publishes unresolved resume 17. Both consumed-intent branches
require successful publication before returning the new cursor. These store
outcomes do not themselves prove that a live process has been recovered.

| Gate | Status | Required next evidence |
| --- | --- | --- |
| Versioned composition | `BLOCKED` | Exact Supervisor/native/Go/runner/root/library/frame/store identities and closed construction/import/load provenance |
| Launch identity and root custody | `BLOCKED` | Supervisor-owned enrolled executable/code identity, exact descriptor/root handoff, replacement/debug/wrong-role refusals; reuse existing APL-2/APL-3 and installation work |
| Restart custody | `BLOCKED` | Installed AttemptID-bound reconciliation before authority reopens; authoritative custody/absence or unresolved refusal. Unfenced consumed intent without completion publishes unresolved resume 17; that refusal is not reconstructed custody |
| Timing | `BLOCKED` | Clock anchors and total elapsed bounds across setup, storage, execution, cancellation and absence; controlled stall evidence with durable gates intact |
| Provenance and source obligations | `BLOCKED` | Exact new composition and candidate review; existing preferred-form kernel source, raw v19/v27 and cross-host work remain separately open |
| Controlled guest execution | `BLOCKED` | Complete reviewed candidate followed by separate authorization naming immutable manifest, host and owned disposable guest |

No new lifecycle helper, PID/path-based authority, storage engine or resource-limit
change is selected. Reuse follows [the adoption map](ECOSYSTEM_REUSE_AND_ADOPTION.md):
platform identity and root primitives, existing Supervisor/G1/G2 ownership and narrow
offline standard-library verification. Identity observed for an XPC peer does not
by itself prove custody of a launched child.

## Next bounded implementation: C5b16 timing/fault baseline

Instrument a versioned derivative of the C5b14B benign C/Go fixture before promoting
an immutable runnable composition. This is the first uncertainty to test because
a durable publication stall may consume time outside the current cleanup clock.
Its result will guide the later launch/custody work without silently changing policy.

Acceptance criteria:

1. Record bounded monotonic phase timings for setup, executable/root checks, durable
   spawn intent, readiness/start, completion/terminal observation, teardown request,
   gate return, signal and absence. Preserve exact refusal phase with fixed labels;
   distinguish observed cause from an unreproduced failure.
2. Inject finite publication delays, including beyond the selected initial-action
   bound, and retain clock-local plus total elapsed intervals. Keep current one-second
   clocks, four/seven-second child/harness containment, no retries, no extra effects,
   no completion/custody fabrication and no widened limits.
3. Compare with exact clock anchors/limits and C2A/ADR-0041 obligations; rerun affected
   native/refusal/mutation checks and independent review. A remedy requiring different
   authority, durability order or architecture needs an explicit decision gate; the
   probe must not bypass a durable gate or signal without native custody.

Scope remains fixed benign local child/C/Go fixtures and owned temporary directories.
No guest, backend, interpreter or installed service is authorized by this selection.
C5b16 closes a timing/fault evidence question; it is not a promise that the current
synchronous store satisfies every bound or that the eventual remedy is implemented.

## Verification

`node scripts/verify.mjs` verifies frozen materials, exact inputs, selected internal
provenance edges and the retained report, then runs 21 positive/substitution/false-
report tests. Unknown, missing, changed or redirected inputs and false execution,
composition, clock, root or restart claims refuse. These tests exercise the audit;
they do not replace native-control mutations or an independent candidate review.
The exact-source extracts are not a general C/Go parser or exhaustive security proof.

Fresh [independent review](https://github.com/Shrimpworks/capsule-experiments/blob/7c22ad119c21d44e74d0a525f96ddcee915330f8/experiments/typed-guest-transport-c5b15-runner-reconciliation/review/REVIEW_2.md)
returned **Ready** after confirming all 19 frozen materials and 13 immutable inputs.
The first review correctly required conditional restart wording; that finding is
closed in the report, source assertions, tests and prose. The
[closure](https://github.com/Shrimpworks/capsule-experiments/blob/7c22ad119c21d44e74d0a525f96ddcee915330f8/experiments/typed-guest-transport-c5b15-runner-reconciliation/review/CLOSURE.md)
records both reviews; frozen README/PLAN status remains historical.

Canonical `pnpm install`, `check`, `lint`, `test`, `verify:schemas`, `verify:adrs`,
Go 1.25.13 tests/vet/build and `govulncheck@v1.6.0` pass. Full golangci retains
50 pre-existing revive documentation findings; blocking non-revive checks and
new-code revive ratchet pass. These are orchestrator-run checks; the reviewer
independently reran the archive verifier, not the complete product suite.
