# Project Definition

## Summary

Capsule is a local-first, capability-controlled JavaScript and TypeScript task runtime for AI
agents. It accepts a bounded job declaration, applies user and host policy, executes generated code
inside a disposable external sandbox, controls every result channel, and destroys the environment.

The task—not the container, VM, shell, or development environment—is the unit of abstraction.

## Problem

Agents frequently need to perform work that is more reliable as code: transform data, validate
configuration, analyze source, generate reports, call an approved API, or run bounded checks.
Executing generated code directly on a user's machine gives it excessive ambient authority. Generic
cloud sandboxes can be opaque, heavyweight, old, or disconnected from local workflows.

Capsule aims to make constrained execution fast and understandable enough to become the default
place for agent-generated JS/TS tasks.

## Goal

Provide a fast, testable execution boundary in which:

- The agent can propose source code and requested capabilities.
- Trusted host or organizational policy can deny or narrow requested authority; the user explicitly
  approves the resulting immutable, digest-bound execution plan in v0.
- Every installation has a device-scoped identity and purpose-separated signing keys.
- Modern runtimes are replaceable, pinned artifacts rather than control-plane dependencies.
- Isolation backends can evolve without changing the job protocol.
- User-defined limits are exact, visible, and never silently expanded or rewritten.
- Results and artifacts have explicit validation and audience policies.
- Every execution produces evidence of the effective boundary.

## Product scope and initial wedge

The platform scope is broader than file processing: it is intended for bounded agent-generated
JS/TS tasks. The first vertical slice is deliberately smaller:

- Local desktop control experience
- Bun and TypeScript
- One-shot execution
- Explicit prepare, human-readable approval, and execute phases
- Per-installation identity with offline verification
- Explicitly granted regular-file or inline JSON inputs
- Immutable, digest-bound snapshots instead of live host-file mounts
- Curated dependencies
- No network, subprocesses, native addons, FFI, or dynamic installation
- Validated CSV, JSON, JSONL, or text artifacts
- User-full and agent-metadata-only artifact delivery by default

This wedge exercises the essential security architecture without turning Capsule into a complete
coding-agent computer.

## Primary users

- Developers building agents that need a safe computation surface
- Desktop AI users who want bounded local execution
- Teams introducing policy around generated code
- Tool authors who need a vendor-neutral execution contract

## Non-goals for the first release

- General remote shell access
- Long-lived development environments
- Browser automation
- Docker-in-Docker
- Arbitrary languages
- Arbitrary package installation during a task
- Background services and public preview URLs
- Multi-tenant hosted scheduling
- Proof that guest code is correct or aligned with user intent

## Principles

1. **No ambient authority.** Every external effect requires an explicit grant.
2. **Fail closed.** Unknown capabilities, profiles, paths, and output types are rejected.
3. **External isolation is mandatory.** Language permissions are supplemental hardening.
4. **Egress is a capability.** Logs, structured results, and artifacts are untrusted output.
5. **Jobs are disposable.** Guest state is never reused across trust boundaries.
6. **Approval is exact.** The user approves one immutable execution-plan digest, not a mutable job.
7. **Identity is scoped.** Device, approval, receipt, and transport authority are explicit and
   purpose limited.
8. **Versions are evidence.** Source, profiles, runtimes, policy, inputs, and outputs are hashed;
   trusted handoffs are signed.
9. **Limits belong to the user.** Defaults and ceilings come from trusted user policy and are not
   silently clamped.
10. **Portability is contractual.** Job semantics are portable; source compatibility is not
    promised.
11. **Security claims are testable.** The attack corpus is part of the product.

## Agreed v0 decisions

- An untrusted proposal is resolved into an immutable execution plan before approval.
- Explicit user approval is required for every v0 plan.
- Approval binds the plan digest, device installation, audience, nonce, and expiry.
- Each installation uses a per-device P-256 DID identity with offline `did:key` resolution.
- Device root, approval, receipt, and transport signing purposes are separated.
- File capabilities identify private immutable snapshots, never agent-supplied or live host paths.
- Resource requests above user-owned ceilings are rejected rather than silently reduced.
- Apple Container is the macOS integration candidate; OCI plus gVisor is the Linux reference
  backend.
- Backends remain development tier until their exact pinned implementations pass the required
  adversarial corpus.
- A blockchain is not required. Content addressing, signed operations, replay protection, and
  receipts borrow useful distributed-system patterns without introducing consensus.

## Success criteria

The first milestone succeeds when a client can submit the same job contract through CLI and MCP,
resolve it into a stable plan, obtain explicit device-bound approval, run it with Bun in a
disposable environment, enforce deny-by-default capabilities and exact user-owned limits, deliver
audience-controlled artifacts, prove teardown, and produce a signed receipt.

The first reference demonstration transforms inline JSON or one user-selected regular-file snapshot
into declared JSON, JSONL, CSV, or text artifacts with no network or ambient authority.

An authoritative milestone additionally requires the exact backend and runtime profile to survive
the documented adversarial corpus before Capsule presents them as a security boundary.
