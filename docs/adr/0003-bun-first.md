# ADR-0003: Bun-first, runtime-neutral protocol

- Status: Accepted
- Date: 2026-07-30
- Gate C refinement: 2026-07-31
- P0-0 stock-runtime refinement: 2026-08-02

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

## Consequences

- The initial product has a clear technical wedge.
- Runtime-specific APIs inside an admitted profile may be used in guest source.
- The Node adapter will test whether the abstraction is real.
- The external isolation boundary must carry the primary security guarantee.
- External isolation does not make a promised runtime-level authority restriction true; unsupported
  profile powers cause admission refusal.
