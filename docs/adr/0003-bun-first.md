# ADR-0003: Bun-first, runtime-neutral protocol

- Status: Accepted
- Date: 2026-07-30
- Gate C refinement: 2026-07-31
- P0-0 stock-runtime refinement: 2026-08-02
- P0-0 governed-construction refinement: 2026-08-02
- Deno-family disposition refinement: 2026-08-02

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

The subsequent [exact stock-runtime investigation](../../experiments/gate-c-bun-runtime-authority/RESULTS.md)
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
[governed-construction review](../../experiments/gate-c-bun-runtime-authority/governed-closure/CONSTRUCTION_REVIEW.md).

Therefore the runtime-neutral protocol portion of this ADR remains accepted, but **Bun is no longer
the selected first implementation candidate**. `RUNTIME-001` remains unsupported. The next runtime
decision must compare alternate runtimes under the unchanged prohibited-power contract and must
supersede this ADR's Bun-first implementation choice. Until that decision is accepted, no real
runtime profile is selected for the first executable slice.

The subsequent [Deno-family experiment](../../experiments/gate-c-deno-runtime-authority/RESULTS.md)
did not produce that superseding decision. Full Deno v2.9.4 retained initial-static-graph, Worker,
SIGUSR1-inspector, Node-compatibility, and runtime-managed storage routes under the exact deny
profile. The minimal `deno_core` 0.409.0 construction removed ambient extensions and its module
loader, but still physically registered 99 built-in core ops before middleware disabled 96, and it
did not include a TypeScript pipeline. The result is `DENO-FAMILY-NO-GO`: no candidate was
selected, this ADR is not superseded, its runtime-neutral portion remains accepted, and
`RUNTIME-001` still refuses.

## Consequences

- The initial product has a clear technical wedge.
- Runtime-specific APIs inside an admitted profile may be used in guest source.
- The Node adapter will test whether the abstraction is real.
- The external isolation boundary must carry the primary security guarantee.
- External isolation does not make a promised runtime-level authority restriction true; unsupported
  profile powers cause admission refusal.
- Runtime construction investigation remains blocking after the Bun and Deno-family NO-GOs;
  Capsule does not relax the advertised authority contract to preserve implementation sequencing.
