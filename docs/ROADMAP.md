# Roadmap

The roadmap is ordered by risk reduction rather than feature count.

## Phase 0A: contract freeze

- Freeze job, execution-plan, approval, profile, capability, artifact, receipt, error, and lifecycle
  schemas.
- Define canonical source, input, policy, plan, artifact, and receipt digests.
- Select the canonical signed-envelope representation and maintained Go and TypeScript libraries.
- Add deterministic digest, signature, nonce, replay, wrong-audience, and revocation test vectors.
- Generate or validate Go and TypeScript contract views against shared positive and negative
  fixtures.
- Tighten total request, source, inline input, snapshot, log, and artifact limits.

Exit evidence:

- All wire objects reject unknown fields and unsupported versions.
- Go, TypeScript, and JSON Schema agree on every shared fixture.
- The same logical plan produces the same digest in Go and TypeScript.
- The security properties and known non-guarantees are current.

## Phase 0B: trust and planning

- Implement the per-installation P-256 device root and platform key-provider interface.
- Implement offline `did:key` resolution and the local trust registry.
- Implement purpose-separated approval, receipt, and transport key authorization.
- Implement trusted user defaults and ceilings with no silent clamping.
- Resolve untrusted proposals into immutable execution plans.
- Generate human-readable approval summaries from the same typed model as the signed plan.
- Implement single-use, audience-bound, expiring approval grants.

Exit evidence:

- Changed source, input, profile, limit, audience, backend, or policy invalidates approval.
- Unknown, wrong-purpose, revoked, replayed, expired, or wrong-installation signatures fail closed.
- The daemon cannot execute a plan that differs from the plan displayed to the user.

## Phase 1: durable lifecycle without guest execution

- Implement the job state machine behind durable repository interfaces.
- Use a fake backend to exercise prepare, approval, staging, execution, collection, destruction,
  cancellation, timeout, restart recovery, and teardown failure.
- Persist approval consumption, identity state, capability metadata, cleanup leases, artifact
  descriptors, and receipt indexes.
- Produce signed receipts without guest-controlled content.

Exit evidence:

- Every post-create path reaches destruction.
- Restart recovery cannot replay an approval or report unknown cleanup as success.
- State transitions are atomic, monotonic, and idempotent.

## Phase 2: minimal Bun execution

- Implement the Bun runtime adapter and verify the draft `bun-data@1` profile.
- Implement the Apple Container development backend with one fresh lightweight VM per job.
- Pin the Apple Container implementation, Linux kernel, init environment, and runtime image.
- Execute inline source and inline JSON with network, environment, subprocesses, native addons,
  FFI, macros, inspector, and installation denied.
- Enforce exact user-approved resource and output limits.
- Return bounded logs and one validated structured artifact through controlled egress.

Exit evidence:

- The guest receives no ambient host filesystem, credential, environment, socket, or network
  authority.
- The plan and receipt identify the exact backend, runtime profile, inputs, limits, outputs, and
  teardown result.
- The backend remains labeled development.

## Phase 3: file-to-artifact reference slice

- Implement trusted regular-file selection and immutable snapshot capabilities.
- Stage snapshots read-only with dedicated bounded scratch and output filesystems.
- Complete declared JSON, JSONL, CSV, and text artifact validation.
- Enforce user-full and agent-metadata-only delivery by default.
- Implement separate signed content-access grants.
- Complete CLI and MCP adapters over the same daemon contract.
- Demonstrate one user-selected input transformed into validated artifacts and a signed receipt.

Exit evidence:

- The guest never receives the original host path or a live host-file mount.
- Link, mutation, special-file, oversized-input, malformed-output, and audience-bypass tests pass.
- MCP cannot issue capabilities, approve plans, or read user-only content.

## Phase 4: authoritative backend validation

- Implement the OCI plus gVisor Linux backend.
- Run the shared adversarial corpus against Apple Container and gVisor.
- Exercise network, filesystem, process, resource, runtime, egress, management-channel, recovery, and
  cross-job-state attacks.
- Build and verify immutable runtime images, SBOMs, signatures, and independent review attestations.
- Measure cold and warm startup, memory, CPU, I/O, and teardown behavior.
- Assign an authoritative tier only to exact pinned configurations supported by retained evidence.

Exit evidence:

- Every mandatory threat-model test passes for the exact backend/profile identity.
- Security claims and limitations are published without overstating unsupported platforms or
  controls.

## Phase 5: runtime portability proof

- Add the Node adapter.
- Build a runtime conformance suite and runtime-neutral fixtures.
- Define runtime discovery, compatibility, deprecation, and version policy.
- Prove that backend and client contracts do not depend on Bun-specific control-plane behavior.

## Phase 6: broader task capabilities

- Custom immutable dependency preparation
- Deno adapter
- Directory and repository snapshots
- Approved network and API brokers
- Approved subprocesses and patch artifacts
- Temporary services
- Hosted scheduling and Firecracker backend

## Deferred decisions

- Public license
- Permanent Go module path
- Public repository and release model
- Windows local-isolation backend
- Portable multi-device identity and recovery
- Externally resolved DID methods and general Verifiable Credentials
- Public transparency or blockchain anchoring
- Hosted offering and multi-tenancy
- Quantitative performance service levels
