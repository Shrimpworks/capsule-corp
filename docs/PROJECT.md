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
- Only the user, trusted host, or organizational policy can grant authority.
- Modern runtimes are replaceable, pinned artifacts rather than control-plane dependencies.
- Isolation backends can evolve without changing the job protocol.
- Results and artifacts have explicit validation and audience policies.
- Every execution produces evidence of the effective boundary.

## Product scope and initial wedge

The platform scope is broader than file processing: it is intended for bounded agent-generated
JS/TS tasks. The first vertical slice is deliberately smaller:

- Local desktop control experience
- Bun and TypeScript
- One-shot execution
- Explicitly granted regular-file or inline JSON inputs
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
6. **Versions are evidence.** Source, profiles, runtimes, policy, inputs, and outputs are hashed.
7. **Portability is contractual.** Job semantics are portable; source compatibility is not promised.
8. **Security claims are testable.** The attack corpus is part of the product.

## Success criteria

The first milestone succeeds when a client can submit the same job contract through CLI and MCP,
run it with Bun in a disposable environment, enforce deny-by-default capabilities and hard resource
limits, deliver audience-controlled artifacts, and produce a verifiable receipt. The authoritative
backend must survive the documented adversarial corpus before being presented as a security
boundary.
