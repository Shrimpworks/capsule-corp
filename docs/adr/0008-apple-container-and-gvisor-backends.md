# ADR-0008: Apple Container and gVisor are independent backend targets

- Status: Accepted
- Date: 2026-07-30

## Context

Capsule is local-first and developed primarily on macOS, while the first planned authoritative
isolation target, gVisor, runs on Linux. Requiring all ordinary development on Linux would harm the
local workflow. Running generated Bun code directly on the macOS host would violate the external
isolation requirement.

Apple Container provides OCI-compatible Linux containers in per-container lightweight virtual
machines on supported Apple silicon and macOS versions. Its implementation and interfaces are
still evolving. gVisor provides an independent OCI-compatible Linux application-kernel boundary.

## Decision

Capsule defines one backend lifecycle and initially implements:

- a fake development backend that never executes guest code
- an Apple Container backend for on-demand macOS integration and candidate authoritative-local use
- an OCI plus gVisor backend as the Linux reference and candidate authoritative or hosted use

Apple Container uses one fresh lightweight VM per job. It must pin the container implementation,
kernel, init environment, OCI image, and runtime profile; provision no network interface; expose no
host bind mounts or host sockets; use a read-only root and dedicated input, scratch, and output
storage; enforce exact limits; and destroy the VM after collection.

The gVisor backend uses one separate sandbox per job with explicitly generated OCI configuration,
filesystem exposure, network denial, resource control, and teardown.

The backends are alternatives under the same protocol. Capsule does not require gVisor inside an
Apple lightweight VM. Neither backend receives an authoritative tier until its exact pinned
implementation passes the mandatory adversarial corpus.

Capsule does not implement a backend that executes untrusted Bun directly on the host.

## Consequences

- Go and TypeScript development, unit tests, protocol tests, and fake-backend tests remain native on
  macOS.
- macOS guest integration does not require a user-managed long-lived Linux VM on supported systems.
- Linux CI or a dedicated worker is required for direct gVisor testing.
- Apple platform and version support is narrower than the portable Capsule protocol.
- Broken or inconvenient networking behavior is never treated as a security control; network
  absence is deliberate and tested.
- Each backend must report its exact identity, controls, security tier, and teardown evidence.
- A backend-specific failure does not change the job contract or justify an authority escape hatch.

## References

- [Apple Container](https://github.com/apple/container)
- [Apple Containerization](https://github.com/apple/containerization)
- [gVisor security model](https://gvisor.dev/docs/architecture_guide/security/)
