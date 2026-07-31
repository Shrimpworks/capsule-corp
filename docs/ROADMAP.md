# Roadmap

The roadmap is ordered by uncertainty and risk reduction rather than feature count. Disposable
spikes may be built outside the final product shape. Their retained evidence—not prototype code
quality—is the deliverable.

## Phase 0: architecture and claim baseline

- Align project, architecture, technical design, threat model, security policy, and ADRs on the
  daemon/Broker/Supervisor authority split.
- Define the component-compromise and control-to-evidence matrices.
- Inventory the target protocol objects and mark current schemas/types as pre-freeze scaffolding.
- State evidence and runtime-integrity claims without implying platform attestation.

Exit evidence:

- No document says the daemon can approve, launch, own user content, and author authoritative
  execution evidence.
- DIDs are consistently described as identifiers, not authorization roots.
- Apple Container and all runtime profiles are development posture.
- Current versus intended behavior is obvious.

## Phase 1: blocking feasibility spikes

Run bounded, disposable prototypes in parallel where practical:

- Go/Swift/TypeScript RFC 8785 plus JWS ES256 interoperability and strict-decoder behavior.
- macOS XPC peer requirements, dynamic validation, Keychain/access-group isolation, and protected
  storage separation.
- Apple Container network, filesystem, resource, management-channel, orphan-recovery, and teardown
  controls.
- Broker-to-Supervisor immutable content-handle transfer without daemon content access or live user
  mounts.
- Supervisor language, per-user versus privileged process model, and need for any tiny launcher.
- Prepared-update, epoch-finalization, interruption, and repair state transitions.

Exit evidence:

- Every spike has recorded setup, exact platform/tool versions, adversarial cases, results,
  limitations, and a decision.
- Daemon key/content/backend access attempts fail in the macOS prototype.
- Exact enforceable backend controls and unsupported limits are known.
- Cryptographic implementations agree or a documented alternative is selected.
- Resulting ADRs and contract vocabulary are updated before schema freeze.

See [Feasibility Spikes](FEASIBILITY_SPIKES.md).

## Phase 2: contract and cryptographic freeze

- Replace the mixed `Job` schema with narrow `JobProposal` semantics.
- Add schemas for plan, registration, approval, attempt, trust snapshot, integrity assessment,
  transcript, artifact manifest, agent summary, and composed receipt.
- Define semantic source-path canonicalization and logical input/output slots.
- Freeze strict raw decoding, canonical bytes, digest, signature, type/domain separation, and safe
  numeric rules using retained cross-language fixtures.
- Define stable error, violation, posture, lifecycle, and recovery records.

Exit evidence:

- JSON Schema, Go, Swift, and TypeScript agree on all applicable fixtures.
- Duplicate keys, unknown fields/versions, unsafe numbers, unsupported powers, wrong types, and
  cross-protocol substitutions fail closed.
- The backend contract requests only controls supported or explicitly rejected by candidates.

## Phase 3: registered-plan and fake-backend lifecycle

- Implement daemon plan generation and Supervisor plan registration.
- Implement direct Broker fetch/render/user-presence approval.
- Implement a locally seeded, signed development `TrustSnapshot`; production TUF service remains
  later work.
- Implement a durable atomic grant ledger and one-attempt semantics.
- Build a fault-injectable fake backend that never runs guest code.
- Implement multi-store saga/reconciliation states and crash injection at every side-effect edge.
- Produce a Supervisor enforcement transcript and composed receipt without guest content.

Exit evidence:

- The daemon cannot execute unregistered or replacement bytes.
- Plan A approval cannot execute plan B, and one grant cannot create two attempts.
- The daemon cannot forge ordinary terminal success.
- Every post-create path reaches explicit destroy, unresolved, or quarantine state.

## Phase 4: inline JSON and content separation

- Implement bounded inline JSON content ownership.
- Implement fixed logical slots and Supervisor staging verification.
- Add a bounded JSON output gate and user-only storage.
- Return a fixed-shape agent summary with no guest-controlled strings, names, sizes, timings, or
  violation detail by default.

Exit evidence:

- Daemon/MCP credentials cannot retrieve user-only content.
- Agent-observable fields stay within the documented channel budget.
- Rich parsing is absent from the daemon and Supervisor.

## Phase 5: dependency-free Bun development execution

- Build one exact package-free runtime bundle.
- Execute one JSON-in/JSON-out job through the selected macOS backend path.
- Deny network and all ambient host resources using the proven mechanisms.
- Bind source, input, runtime, backend, controls, output, integrity evidence, and teardown into the
  attempt transcript.

Exit evidence:

- The exact configuration passes the minimum development attack suite.
- Unsupported controls refuse execution rather than being silently approximated.
- The backend remains clearly labeled development.

## Phase 6: regular-file snapshot vertical slice

- Implement native file selection and immutable regular-file data-fork snapshots.
- Transfer job-scoped handles to the Supervisor without daemon content access.
- Add bounded scratch/output storage and JSON/JSONL/text artifacts, then CSV.
- Add audience-controlled release and separate content grants.
- Complete CLI and MCP adapters over the same daemon protocol.

Exit evidence:

- No agent or guest contract contains an original host path.
- Link, mutation, special-file, oversized-input, formula/terminal/bidi, malformed-output, and
  audience-bypass tests pass.

## Phase 7: authoritative validation

- Implement the OCI plus gVisor Linux backend.
- Run the complete shared attack corpus against exact Apple and gVisor configurations.
- Build runtime SBOM, provenance, review attestation, registry, and backend validation records.
- Test runtime-integrity failure, update, repair, cancellation, restart, orphan, and teardown paths.
- Measure startup, CPU, memory, storage, I/O, and cleanup behavior.

Exit evidence:

- Only exact pinned configurations supported by retained evidence may use `validated-local` or a
  stronger isolation posture.
- Published claims include limitations and never collapse posture dimensions into one unsupported
  label.

## Phase 8: production trust repository and updates

- Operate root, targets, snapshot, timestamp, and delegated TUF roles.
- Publish releases, runtime bundles, review records, validation records, and Capsule-defined
  revocation/disable objects.
- Produce compact signed local trust snapshots outside the live Supervisor path.
- Support offline bundles and pinned self-hosted repositories.
- Implement explicit crash-safe install, update, repair, and key-replacement ceremonies.

## Phase 9: optional Guardian and external witness

- Evaluate Endpoint Security entitlement and deployment requirements.
- Begin with notify-only observations; never make Guardian approval authority.
- Evaluate privacy-reviewed installation/receipt checkpoint witnessing and transparency monitoring.
- Send no job content, filenames, or per-execution identifiers by default.

## Deferred

- Node and Deno portability
- network and API brokers
- secrets brokerage
- directory and repository snapshots
- rich document/media parsing sandboxes
- automatic background update delegation
- portable multi-device identity and recovery
- externally resolved DID methods and general Verifiable Credentials
- Windows
- hosted multi-tenancy
- Firecracker
- platform attestation
- public transparency or blockchain anchoring
