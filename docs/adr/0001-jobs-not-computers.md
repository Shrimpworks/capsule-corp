# ADR-0001: Jobs, not computers

- Status: Accepted
- Date: 2026-07-30

## Context

Many agent sandbox products expose a remote machine, shell, container, or long-lived development
environment. That abstraction encourages broad filesystem, process, package, network, and lifecycle
authority.

Capsule needs a smaller unit that can be validated and authorized before execution.

## Decision

The public abstraction is a bounded execution job containing source, inputs, requested capabilities,
limits, and result contracts. Containers, VMs, processes, and pools are private backend details.

The first API will not expose arbitrary commands, container flags, image references, or general
interactive shell sessions.

## Consequences

- Policies can be described in task terms.
- Backends can change without changing client integrations.
- Interactive development and arbitrary shell workflows are deferred.
- Some useful workloads require new explicit capabilities rather than escape hatches.
