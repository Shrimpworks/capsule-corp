# ADR-0028: Select governed deno_core as the first runtime candidate

- Status: Accepted
- Date: 2026-08-03
- Fork-integration evidence refinement: 2026-08-03
- Supersedes: ADR-0003's Bun-first implementation ordering only; its runtime-neutral protocol
  decision remains accepted
- Refines: ADR-0013, ADR-0016, and ADR-0026

## Context

ADR-0003 chose Bun as Capsule's first runtime implementation while keeping the protocol
runtime-neutral. The retained Gate C campaign later rejected both stock Bun 1.3.14 and a governed
Bun construction for the unchanged v0 authority contract. The governed construction crossed its
predeclared reviewability threshold at a conservative minimum of 40 hand-authored source files and
10 generated outputs, without a narrow way to remove Worker and native-loading authority while
preserving required runtime threads.

Full Deno v2.9.4 and the first unpatched `deno_core` 0.409.0 construction also failed their exact
tests. Full Deno retained authority outside the proposed profile, while unpatched `deno_core`
registered 99 built-in ops. Later bounded experiments established a narrower governed candidate:

- a reviewable pre-registration patch physically registers and links only the three required
  bootstrap ops;
- fixed restoration mutations and the fixed dependency-free JavaScript fixture fail or pass as
  expected;
- two clean same-host builds reproduce the exact binary and snapshot;
- a 22-entry snapshot-package-derived dynamic root runs the fixed fixture without ambient host
  library, cache, NSS, locale, timezone, or package-database fallback; and
- Proposed ADR-0026 defines a deterministic pre-approval strip-only TypeScript boundary without
  placing a transformer dependency graph in the live runtime.

These results do not admit a runtime. Independent Linux/arm64 reconstruction, complete V8
source/build/notice publication, production TypeScript ownership and protocol migration, external
isolation composition, installed-byte custody, and the complete runtime-profile corpus remain
open. The official prebuilt `rusty_v8` archive cannot satisfy the required publication evidence,
so continuing with experiment-only patches or copied registry sources would not produce a
governed product dependency.

The repository owner created two real GitHub forks on 2026-08-03. The governed Deno fork was later
transferred while idle; the current integration repositories are:

- [`Shrimpworks/deno`](https://github.com/Shrimpworks/deno), forked from `denoland/deno` and
  transferred from the historical `dills122/deno` location with its branches, merged PR #1, and
  Actions history intact; and
- [`dills122/rusty_v8`](https://github.com/dills122/rusty_v8), forked from
  `denoland/rusty_v8`.

Both forks now contain merged Capsule-governed branches. The governed Deno line was merged from
head `9adb0b68b55bca81644827f1e7749a3acb091bed` at merge
`ea18b9dc21ff8ebd19347be7095f47937ee14ec2`. The governed `rusty_v8` follow-up was merged from head
`a43ee7486c3e05bce5d6e5db586b3e2e688c33cf` at merge
`a31b8f39dc6933d5635367e8ccb67d70f2cc2385`. No governed release or admitted artifact exists. The
first fork-native Linux/arm64 construction stopped before building because the merged `rusty_v8`
publication profile supported Linux/amd64 only.

## Decision

Capsule selects a governed `deno_core` construction as the first runtime engineering candidate.
Node remains a later portability and contingency target. Full Deno is not the selected profile,
and Bun is no longer on the first-runtime path.

This is an implementation-order and governance decision, not runtime admission:

- the runtime-neutral plan, profile, adapter, and evidence contracts remain unchanged;
- `RUNTIME-001` remains `unsupported`, and execution requiring it must refuse;
- no current experiment binary, official prebuilt `rusty_v8` archive, copied Cargo source tree, or
  patch file is a product dependency;
- no runtime may be wired to a consumer or real backend until an explicit admission decision binds
  the exact governed fork commits, release artifacts, runtime root, profile, external enforcement,
  and retained evidence; and
- the external isolation boundary remains mandatory and must contain a fully hostile guest even
  when the runtime-level prohibited-power contract also passes.

The governed source line starts from these retained upstream anchors:

- Deno v2.9.4 commit `14eea3160ae5834476aa3b9d317b8d41d991b982`; and
- `rusty_v8` v150.2.0 commit `d305e6afa7736f6e298c30ae6646f7709ee9382b`.

`Shrimpworks/deno` will carry the ordered `deno_core` physical-omission and deterministic-ordering
changes as reviewable commits on a Capsule-governed branch. `dills122/rusty_v8` will own the
reconstructible Linux/arm64 builder, exact source and generated-build closure, licenses/notices,
SBOM, provenance, restoration tests, advisory ownership, and update/rebase policy. Capsule will not
fork `denoland/v8` unless Capsule must change that source line directly.

The Capsule repository will consume only immutable governed fork commits and release digests after
their evidence passes. The current experiment patch stack remains retained evidence and the oracle
for bootstrapping the forks, not the shipping integration mechanism.

## Consequences

- Runtime work now has one primary direction instead of repeatedly reopening Bun, full Deno, or
  Node as equal first-slice candidates.
- Maintaining two security-sensitive forks adds release engineering, review, advisory response,
  corresponding-source publication, and upstream-rebase obligations.
- The governed fork branches are now established at exact merged commits. The first fork-native
  integration check stopped before building because the `rusty_v8` release pipeline supports only
  Linux/amd64 and has no intended Linux/arm64 profile. A pinned arm64 builder/publication profile,
  actual fork-native outputs, and a release still must exist before another packaging or admission
  claim can advance.
- Independent-builder evidence, TypeScript contract migration, libkrun/launcher/transport
  composition, installed signing/notarization, and final profile admission remain separate gates.
- A failure of the governed fork to retain a small reviewable authority surface, reproduce exact
  bytes, publish complete materials, or pass the composed hostile-runtime corpus reopens the
  candidate decision. It does not silently weaken the v0 authority contract.
