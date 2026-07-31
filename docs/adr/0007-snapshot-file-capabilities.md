# ADR-0007: Regular-file capabilities reference immutable snapshots

- Status: Accepted
- Date: 2026-07-30

## Context

Giving a guest a live host path exposes path traversal, link replacement, time-of-check/time-of-use,
special-file, concurrent mutation, and accidental over-broad mount risks. An agent must not gain
authority by naming a path, and an approval must identify the exact bytes that will enter the
execution environment.

The first product wedge needs user-selected arbitrary file content, but does not require arbitrary
filesystem object semantics.

## Decision

A trusted user action selects a candidate file. The input broker opens it without following links,
verifies it is a bounded regular file, copies from the opened descriptor into private staging while
hashing the copied bytes, rejects material mutation during snapshotting, and records a signed input
snapshot manifest.

The broker issues an opaque, short-lived, job-bound, read-only capability to the immutable snapshot.
The execution plan binds its digest. The guest receives a read-only staged copy and never receives
the original host path or a live host mount.

Arbitrary v0 file content is treated as opaque bytes. Complex parsing occurs inside the guest.
Directories, symlinks, devices, sockets, FIFOs, host-side archive extraction, and other special
filesystem objects require future explicit capability contracts and are rejected by the v0 broker.

## Consequences

- User-visible paths remain trusted-host metadata and do not enter agent-facing contracts or guest
  configuration.
- Approval identifies immutable input bytes rather than mutable filesystem names.
- Snapshotting consumes time and storage proportional to input size; user policy must bound both.
- Changes to the original file after snapshotting do not change an approved job.
- Snapshot staging and retention become sensitive data-handling responsibilities.
- Platform-specific safe-open and mutation-detection implementations require positive, negative,
  race, and adversarial tests.
- Rich host-side parsers and archive extractors remain outside the trusted control plane.
