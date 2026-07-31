# Threat Model

Status: draft

This document describes the intended security boundary. The current scaffold does not yet implement
or satisfy these properties.

## Security objective

Contain an AI-generated JS/TS task so that it can exercise only authority explicitly granted by a
trusted user, host application, or policy, while bounding resource consumption and controlling all
data returned from the guest.

## Assets

- Host filesystem and user data not granted to the job
- Host and AI-client credentials
- Control-plane process and configuration
- Other jobs and their inputs, outputs, and state
- Runtime profile registry and signing trust roots
- Integrity of the effective policy and execution receipt
- Host availability and bounded resource use
- Confidentiality of guest output not approved for the agent

## Adversaries and untrusted inputs

Capsule assumes the following may be malicious:

- AI-generated source code
- Prompts and model output
- Submitted source bundles
- User-selected input content
- Third-party dependencies
- Runtime parsers and libraries
- Guest stdout, stderr, filenames, structured results, and artifacts
- A client attempting to request more authority than the user intended

The local operating-system administrator and authoritative backend operator are trusted in the
initial model. Defending against a compromised host kernel is out of scope for local execution.

## Trust boundaries

### Client to daemon

The client is authenticated and authorized to request jobs, but generated request content remains
untrusted. A valid request is not automatically an authorized request.

### Capability issuance

Only trusted user or host actions may turn a host resource into a capability. An agent may reference
an issued capability but may not mint one from a path, URL, environment variable, or identifier it
invented.

### Daemon to isolation backend

The daemon passes a fully resolved effective policy to the launcher. Backend configuration is
generated from trusted code and must not contain guest-controlled shell interpolation, mount flags,
image references, or seccomp rules.

### Guest to host

The guest is hostile. Its syscalls, filesystem access, process creation, IPC, and network access are
controlled externally. Runtime permission systems are supplemental only.

### Guest output to user or agent

All guest-controlled output crosses an egress broker. Content delivery is separate from metadata
delivery and follows an explicit audience policy.

## Mandatory security properties

### Isolation

- The guest cannot read or write arbitrary host paths.
- Inputs are staged or mounted with the exact granted access.
- The runtime root is immutable during a job.
- Scratch and output storage are isolated and size-limited.
- Host environment variables, open descriptors, sockets, and credentials are not inherited.
- Guest processes cannot signal, inspect, or attach to host or other-job processes.
- Guest state is destroyed after completion or cancellation.

### Network and IPC

- Network is externally denied by default.
- DNS, loopback, link-local, metadata services, and Unix sockets are included in the denial.
- Future network access must use either an isolated network policy or a broker that revalidates
  destinations and redirects.
- Container-engine, agent, SSH, credential, display, and host-service sockets are never inherited.

### Resource control

- Wall time and CPU time are bounded independently.
- Memory, PIDs, temporary storage, output count, and output bytes are bounded externally.
- Output flooding cannot exhaust daemon memory.
- Cancellation terminates the complete guest process tree.

### Runtime hardening

- Bun automatic installation and `.env` loading are disabled.
- Native addons, FFI, dynamic package installation, macros, subprocesses, and inspector access are
  disabled unless a future profile explicitly grants them.
- Profiles are selected from a trusted registry and resolved to immutable digests.

### Egress

- stdout and stderr are capped and not automatically exposed in full to the agent.
- Only declared artifacts from the dedicated output volume are considered.
- v0 artifacts must be regular files with exact paths, allowed types, and bounded sizes.
- Symlinks, hard-link tricks, device files, sockets, FIFOs, sparse-file abuse, and archives are
  rejected in v0.
- Agent content access requires an audience grant separate from user delivery.

## Non-guarantees

Capsule does not prove:

- That guest code performs the requested task correctly
- That an approved output does not contain copied or encoded input data
- That secret-pattern redaction can identify every sensitive value
- That a supported runtime or host kernel contains no unknown vulnerability
- That source written for one runtime behaves identically on another
- That inputs supplied through an AI client's normal attachment flow remain outside model context

## Abuse cases and required tests

The authoritative backend must cover at least:

| Category | Cases |
| --- | --- |
| Filesystem | traversal, absolute paths, symlinks, hard links, `/proc`, home and credential reads |
| Process | fork bomb, worker creation, signals, inspector activation, orphan processes |
| Network | TCP, UDP, DNS, loopback, IPv6, Unix sockets, metadata endpoints |
| Runtime | native addons, FFI, Wasm abuse, dynamic imports, Bun auto-install, `.env` loading |
| Resources | busy loop, heap OOM, native allocation, PID exhaustion, disk fill, output flood |
| Artifacts | undeclared paths, oversized output, sparse files, devices, FIFOs, archive bombs |
| Isolation | cross-job state, cached writable data, inherited descriptors, environment leakage |
| Protocol | unknown fields, duplicate identifiers, invalid limits, unsupported capabilities |

## Security tiers

- **Development:** convenient local backend; useful for development, not an authoritative claim.
- **Authoritative local:** isolated Linux worker with the documented mandatory controls.
- **Hosted hardened:** gVisor or stronger boundary, signed profiles, multi-user isolation, monitoring.
- **High assurance:** future microVM backend with dedicated tenant boundaries.

Every receipt and user-facing status should identify the backend tier used.

## Open questions

- Exact local Linux virtualization strategy on macOS and Windows
- Profile signing and trust-root distribution
- Capability persistence and revocation semantics
- Safe handling of user-approved network access
- Whether structured guest results need a distinct policy from file artifacts
- Retention and encryption policy for receipts and staged artifacts
