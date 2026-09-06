# C5b pre-C5b16 focused security audit

Date: 2026-09-06

Status: focused audit `PASSED`; no reportable vulnerability found in the exact
reviewed scope. C5b16 remains `BLOCKED` on implementation and review. Installed
lifecycle, controlled guest execution, runtime/profile admission and product
admission remain `BLOCKED`.

## Decision

Proceed to a bounded C5b16 timing/fault evidence slice on the existing benign
C/Go fixture. No security fix blocks that experiment. Do not promote the fixture,
change durability order, signal before durable teardown intent, or interpret a
fixture self-alarm as Supervisor-enforced absence.

C5b16 must record both durable-gate latency and post-gate cleanup latency, plus
total elapsed time from the applicable initial action. Its first retained fault
case must cover teardown-gate pre-publication failure and assert the exact phase,
`spawn=1`, `kill=0`, unresolved cursor 17, no completion and fixture-alarm-only
containment. Finite delays below, at and above the 1,200-ms C2A total bound must
remain distinguishable from the unchanged 1,000-ms local cleanup clock and the
four/seven-second fixture/harness stops. Any remedy that changes authority,
durability order or architecture requires a separate decision gate.

## Exact scope and baselines

Defensive scope was confined to this repository, the owned `capsule-experiments`
archive and controlled local disposable fixtures. No real guest, third-party
system, external credential, installed service or unrelated data was accessed.

- Canonical baseline: `a26f83d5a99d46c7a4c7a1eb2c72bad290b7decc`.
- Owned archive merge: `cbfef30751f0a2dd8d8a71356c6b975b52d3e752`.
- C5b15 retained source/evidence: `7c22ad119c21d44e74d0a525f96ddcee915330f8`.
- C5b14B reviewed source/evidence: `294559e78b410f4b25e0e25fcf373ec0c7d1cfb5`.
- Reviewed contracts: C5b11-C5b15 checkpoints, protocol object model, durable
  completion contract, C2A execution profile, threat model and relevant ADRs.
- Reviewed implementation: native lifecycle/providers, Go c-archive bridge and
  store, fixed transport ABI, fixtures, generators and verification scripts.

Method combined independent architecture/threat, baseline and focused
source-to-sink passes with direct source inspection and fresh deterministic
execution. All passes returned no reportable finding.

## Security result

| Surface | Result | Evidence boundary |
| --- | --- | --- |
| Supervisor-only authority and registration/attempt binding | `PASSED` | Fixed registered effects remain Supervisor-owned; execute accepts no replacement plan, backend flag, image, mount, guest path or daemon authority. |
| Native PID custody, teardown and completion | `PASSED` for fixture semantics | Signal/reap/absence observations use the spawned child; completion follows durable terminal state and observed absence. This does not establish installed process-tree or restart custody. |
| Go/C ABI and fixed transport | `PASSED` for private fixture ABI | Inputs are bounded and synchronously copied; no retained caller pointers or unsafe outcome translation found. Fixed-width trusted framing is not a general hostile-input parser. |
| Durable intent, publication, replay and restart | `PASSED` for declared model | Effects follow durable gates; indeterminate publication fences; consumed intent without completion becomes unresolved; restart does not adopt a PID or redrive authority. Same-UID rollback resistance is not implemented. |
| Claims versus evidence | `PASSED` | C5b14B execution and C5b15 static reconciliation remain fixture/static evidence only; neither supports installed, guest, runtime/profile or product admission. |

No candidate met the reportable-vulnerability threshold. Four important limits
remain explicit rather than findings in the current declared scope:

1. Cleanup clock starts after durable teardown publication returns. It therefore
   does not bound publication wait from initial action and cannot by itself prove
   C2A forced-absence timing.
2. Teardown-gate pre-publication failure sends no signal. Fresh case 7 reproduced
   `result=-2 spawn=1 kill=0 phase=2 cursor=17/17 unresolved=1 complete=0` in
   4.47 seconds. This is conservative state handling with fixed-child self-alarm
   containment, not product cleanup compliance.
3. Fixture executable hashing followed by pathname launch has no installed
   same-UID replacement defense. Exclusive controlled-directory scope makes this
   an explicit unimplemented installed control, not a supported security claim.
4. Fixture root/store paths have no protected installation identity or coherent
   rollback defense. The exclusive local root assumption must not cross into an
   installed-candidate claim.

## Fresh verification

- C5b14B `node scripts/verify.mjs`: `PASSED` for 126 native cases, 17 Go suites
  with 58 subtests, 24 providers, 10 compiled assertion mutations, race, vet and
  build checks.
- C5b15 `node scripts/verify.mjs`: `PASSED` for 13 exact inputs, retained
  provenance/report checks and 21 refusal tests.
- Focused C5b14B case 7: `PASSED` as a deterministic conservative-refusal
  reproduction; it is not evidence of bounded Supervisor teardown.
- Canonical `pnpm install`, `check`, `lint`, `test`, `verify:schemas` and
  `verify:adrs` passed. Go 1.25.13 tests, vet, build and
  `govulncheck@v1.6.0` passed.
- Full `golangci-lint` retains the known 50 pre-existing `revive` documentation
  findings. Blocking non-`revive` analysis and the new-code `revive` ratchet
  passed with zero issues.

The host-default Go 1.26.5 is not the repository verification toolchain. A
preliminary vulnerability run under that unpinned version reported standard
library advisories; rerunning all reported Go gates under the required Go
1.25.13 pin produced the results above.

No new experiment harness or archive evidence was needed. Existing immutable
fixtures reproduced the relevant edge, so this audit retains only canonical
conclusions and the next decision in `capsule-corp`.

## Deferred installed-candidate audit

Before any installed or beta claim, audit exact enrolled executable identity and
protected launch, installation/store ownership and rollback resistance,
AttemptID-bound restart/process-tree custody, real transport parser bounds,
resource enforcement and total deadlines, signing and authenticated IPC,
Approval Broker transaction binding, real isolation backend/guest behavior, and
power-loss/sleep/wake/update/cross-host evidence. Current result supplies none of
those claims.
