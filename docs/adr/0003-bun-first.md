# ADR-0003: Bun-first, runtime-neutral protocol

- Status: Accepted
- Date: 2026-07-30

## Context

Capsule aims to support modern JS/TS runtimes, but implementing Bun, Node, and Deno simultaneously
would expand the initial scope. Bun offers fast startup and native TypeScript, while lacking a
complete hostile-code security model on its own.

## Decision

The job protocol and adapter interfaces are runtime-neutral. Bun is the first runtime implementation,
Node is the second portability proof, and Deno follows later.

Source compatibility across runtimes is not required. Runtime selection is explicit and resolved to
an immutable profile digest.

## Consequences

- The initial product has a clear technical wedge.
- Runtime-specific APIs are allowed in guest source.
- The Node adapter will test whether the abstraction is real.
- The external isolation boundary must carry the primary security guarantee.
