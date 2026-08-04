# ADR-0003: Bun-first, runtime-neutral protocol

- Status: Accepted; first-runtime implementation ordering superseded by ADR-0028
- Date: 2026-07-30
- Gate C refinement: 2026-07-31
- P0-0 stock-runtime refinement: 2026-08-02
- P0-0 governed-construction refinement: 2026-08-02
- Deno-family disposition refinement: 2026-08-02
- Governed `deno_core` physical-omission refinement: 2026-08-02
- Governed `deno_core` reproducible-package refinement: 2026-08-02
- Governed `deno_core` direction and V8 source/license refinement: 2026-08-02
- TypeScript approved-byte refinement: 2026-08-03
- Governed `deno_core` runtime-root and product-direction refinement: 2026-08-03
- Governed `deno_core` selection: superseded in part by ADR-0028 on 2026-08-03

## Context

Capsule aims to support modern JS/TS runtimes, but implementing Bun, Node, and Deno simultaneously
would expand the initial scope. Bun offers fast startup and native TypeScript, while lacking a
complete hostile-code security model on its own.

## Decision

The job protocol and adapter interfaces are runtime-neutral. Bun is the first runtime implementation,
Node is the second portability proof, and Deno follows later.

Source compatibility across runtimes is not required. Runtime selection is explicit and resolved to
an immutable profile digest. Runtime-specific APIs are permitted only when they are inside that
profile's approved authority contract. The Bun-first choice does not admit subprocess, FFI, native
addon, inspector, macro, environment-file, package-install, or other powers prohibited by v0.

Gate C proved that pinned Bun 1.3.14 executes the fixture but did not prove a non-bypassable
no-subprocess/no-FFI profile. Bun therefore remains the intended first runtime only if the
[Gate C P0 runtime-authority campaign](../GATE_C_P0_RECONCILIATION.md) passes. Otherwise Capsule
must use a governed patch/alternate runtime or explicitly revise this decision and the product
contract.

The subsequent [exact stock-runtime investigation](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-bun-runtime-authority/RESULTS.md)
rejected stock Bun 1.3.14 for this profile: relevant flags did not remove direct or aliased process,
`execve`, FFI/native-loader, inspector, Worker, or inherited-descriptor authority. Bun-first now
means only a governed construction-level patch plus exact external enforcement that later passes
P0-0. Until then `RUNTIME-001` is unsupported and execution requiring it refuses. If that governed
branch fails, Capsule selects an alternate runtime and updates this ADR; it does not weaken the
prohibited-power contract implicitly.

The governed-construction branch has now failed its explicit reviewability gate. Exact source
review expanded the minimum closure from the prior 21-file lower bound to at least 40 hand-authored
files and 10 generated outputs across build identity, registries, loaders, native sinks, globals,
configuration, resolution, and restoration backstops. A narrow process/exec self-seal cannot close
Worker or native loading while preserving required lazy runtime/JIT threads. The experiment stopped
before authoring or building a partial patch, as required by its fail-fast rule. See the
[governed-construction review](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-bun-runtime-authority/governed-closure/CONSTRUCTION_REVIEW.md).

Therefore the runtime-neutral protocol portion of this ADR remains accepted, but **Bun is no longer
the selected first implementation candidate**. `RUNTIME-001` remains unsupported. The next runtime
decision must compare alternate runtimes under the unchanged prohibited-power contract and must
supersede this ADR's Bun-first implementation choice. Until that decision is accepted, no real
runtime profile is selected for the first executable slice. ADR-0028 later fulfills the
implementation-order decision without admitting a runtime profile.

The subsequent [Deno-family experiment](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-runtime-authority/RESULTS.md)
did not produce that superseding decision. Full Deno v2.9.4 retained initial-static-graph, Worker,
SIGUSR1-inspector, Node-compatibility, and runtime-managed storage routes under the exact deny
profile. The minimal `deno_core` 0.409.0 construction removed ambient extensions and its module
loader, but still physically registered 99 built-in core ops before middleware disabled 96, and it
did not include a TypeScript pipeline. The result is `DENO-FAMILY-NO-GO`: no candidate was
selected, this ADR is not superseded, its runtime-neutral portion remains accepted, and
`RUNTIME-001` still refuses.

The bounded [governed `deno_core` follow-up](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-physical-omission/RESULTS.md)
passed one narrower question left by that disposition. A one-file pre-registration allowlist
reduced the exact built-in registry from 99 ops to three bootstrap-required ops; two
ASLR-controlled clean builds produced identical snapshots and binaries, and the final binary
exposed only those three built-in op symbols. Fixed prohibited-power and restoration probes still
refused. This removes the physical-built-in-registration objection for a governed patched
construction, but it does not supply the excluded TypeScript transformation, independently
reconstructible runtime
package, complete runtime-profile admission, or external-isolation composition. The earlier
full-Deno and unpatched-`deno_core` NO-GO remains historical fact. At that checkpoint no runtime was
selected and this ADR was not superseded; `RUNTIME-001` remained unsupported.

The subsequent [reproducible-package follow-up](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-reproducible-package/RESULTS.md)
closed the prior local builder-image ambiguity for that bounded candidate. A digest-pinned official
Rust base with no added apt state and a complete 191-crate source bundle reproduced the same
snapshot and binary in two clean same-host containers, and a normalized two-file candidate archive
matched completely. This is not independent-builder provenance. The exact prebuilt V8 archive still
lacks a complete corresponding source/third-party-notice record, the candidate archive depends on
an unbundled dynamic Bookworm root, and production TypeScript ownership/wiring plus
external-isolation/profile admission remain open. The experiment therefore returned NO-GO for
runtime-selection evidence. At that checkpoint no runtime was selected and this ADR was not
superseded; `RUNTIME-001` remained unsupported.

The exact [V8 source/license follow-up](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-v8-source-license-closure/RESULTS.md)
then bound the official Linux/arm64 asset to its release job, exact `rusty_v8`/V8 source revisions,
Chromium V8 base, and four-patch Deno V8 stack. It returned `SOURCE-LICENSE-CLOSURE-NO-GO` because
mutable publisher inputs, missing exact GN/Ninja link metadata, and absent generated notices block
an independent rebuild and complete source/notice publication. Product direction is governed
`deno_core` as the intended first engineering candidate, with Node reserved for a later
portability/contingency proof. ADR-0028 now accepts that implementation ordering and supersedes
this ADR's Bun-first ordering only. It does not select or admit a runtime profile or change
`RUNTIME-001`.

The bounded [self-contained runtime-root follow-up](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-runtime-root/RESULTS.md)
closed the candidate archive's ambient Bookworm dynamic-root dependency. The exact governed
binary executes its fixed fixture through a packaged, snapshot-pinned loader and library closure
inside a 22-entry immutable root; read-only syscall evidence found no host library, cache, NSS,
locale, timezone, or package-database load, and deliberate loader, library, manifest, ownership,
mode, relocation, environment, and entry-cap mutations failed closed. Two clean same-host
containers produced the same normalized root bytes. This is bounded construction evidence, not
independent-builder provenance, external-isolation validation, or profile admission.

After the hard Bun pivot, governed `deno_core` is the intended first runtime engineering direction.
ADR-0028 accepts that direction and records the real governed Deno and `dills122/rusty_v8` forks;
the Deno fork has since transferred from `dills122/deno` to `Shrimpworks/deno` as repository-
governance work. That transfer does not select or admit a runtime profile and does not promote
`RUNTIME-001`. Continuation requires governed branches and a release/provenance workflow in those
forks rather than a maintained copied registry tree or experiment-only patch stack, plus closure
of the remaining V8 source/notices, production transformation and protocol ownership,
installed-root custody, external isolation, and full-profile evidence.

The bounded [TypeScript approved-byte follow-up](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/typescript-approved-byte-boundary/RESULTS.md)
passed another narrow question. Exact Node 22.22.1/Amaro 1.1.5 strip-only emission was
deterministic for fixed fixtures and supports [Proposed ADR-0026](0026-bind-pre-approval-typescript-erasure.md),
which binds original and executable bytes before registration and forbids post-approval
transformation. The experiment did not choose a production transformer owner, change product
contracts, add a `deno_core` module loader, or close packaging, external isolation, or full runtime
admission. ADR-0028 supersedes only this ADR's Bun-first ordering; `RUNTIME-001` continues to
refuse.

## Consequences

- The initial product has a clear technical wedge.
- Runtime-specific APIs inside an admitted profile may be used in guest source.
- A later Node adapter remains a portability/contingency proof rather than first-runtime selection.
- The external isolation boundary must carry the primary security guarantee.
- External isolation does not make a promised runtime-level authority restriction true; unsupported
  profile powers cause admission refusal.
- Governed `deno_core` is the intended first runtime engineering direction after the Bun pivot.
  Its physical-op, reproducible-package, V8 source/license, TypeScript approved-byte, and
  self-contained-root results narrow the investigation but do not admit a runtime. ADR-0028 now
  selects the engineering order and fork-governance direction.
  Capsule does not relax the advertised authority contract to preserve implementation sequencing.
