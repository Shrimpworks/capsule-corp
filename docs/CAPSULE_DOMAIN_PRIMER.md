# Capsule domain primer for Go contributors

This is a fast-orientation map of Capsule's vocabulary and Go package layout for someone about to
write or review Go code here. It summarizes; it does not supersede. Where this primer and a
canonical document disagree, the canonical document wins:

1. `AGENTS.md` — binding security/process rules.
2. `docs/PROJECT.md` — current status, principles, non-goals.
3. `docs/ARCHITECTURE.md` — authority and storage boundaries.
4. `docs/TECHNICAL_DESIGN.md` — integrated protocol/trust/lifecycle detail.
5. `docs/security/THREAT_MODEL.md` — adversaries and invariants.
6. The relevant ADR(s) in `docs/adr/` for the exact component being touched.

## The one-paragraph model

Capsule is a local-first trusted execution platform for bounded JavaScript/TypeScript jobs
proposed by AI agents. It splits authority three ways so that compromising the agent-facing
service does not also grant approval, content, or execution authority: an untrusted-input-facing
Go daemon that may only propose and plan; a trusted native Broker that owns human approval and
user-only content; and an Execution Supervisor that is the only component allowed to create a
hostile guest. "The task, not the container/VM/shell," is the public abstraction.

## Core actors

| Actor | Language | Authority | Notes |
|---|---|---|---|
| Agent-facing daemon | Go (`cmd/capsuled`, `internal/api`) | Proposal validation and planning only | Never holds approval, content, or execution authority (`AGENTS.md`) |
| Trusted Host Broker | Native (Swift/ObjC preferred) | Approval Broker: user-presence approval. Content Broker: user-only file content | No agent-facing endpoint; logically two interfaces |
| Execution Supervisor | Go core, native IPC front end | Sole authority to create/terminate/destroy a hostile guest | Independently enforces hard-safety rules that nothing else can override |
| Hostile guest | Runtime-dependent | None by default | The sandboxed thing that actually runs untrusted job code |

## Job lifecycle vocabulary

```
JobProposal -> ExecutionPlan -> PlanRegistration -> Approval -> Attempt -> execution -> Receipt
```

- **JobProposal**: the agent's untrusted request for work.
- **ExecutionPlan**: the daemon's validated, narrowed plan built from a proposal.
- **PlanRegistration**: the Supervisor's durable, immutable record of a specific plan, identified
  by a Supervisor-issued registration ID.
- **Approval**: one-use, attempt-bound human authorization. Binds plan digest, registration,
  installation, trust epoch, expected Supervisor, attempt nonce, purpose, audience, and expiry.
  Rendered from Supervisor-registered typed data — the Broker never renders daemon-supplied text.
- **Attempt**: the single consumption of an approval; at most one attempt per approval, ever.
- **Receipt**: composed from Broker approval evidence plus Supervisor enforcement transcript —
  not independent platform attestation, and not to be described as one.

The load-bearing rule that shows up throughout the Go code: **execute by Supervisor-issued
registration ID only.** Never accept replacement plan bytes, backend flags, images, mounts, or
guest paths at execute time — see `internal/execution/registeredlifecycle`, which resolves only a
Supervisor-issued `AttemptID` and revalidates the committed attempt and plan before doing anything.

## Go package map

- `cmd/capsuled` — the daemon binary entry point.
- `internal/api` — the daemon's local HTTP surface. Currently read-only/unauthenticated
  (`/healthz`, `/v1/version`, `/v1/runtimes`). Future job-facing endpoints must proxy to the
  Supervisor rather than implement authority here.
- `internal/protocol/v0candidate`, `v0cbor`, `tsapprovedbytecandidate`, `sourcevalidatorpassive` —
  passive/candidate wire-format and contract code: byte-exact CBOR encode/decode and
  conformance-fixture-verified types. **Candidate** and **passive** here mean exactly that: not an
  activated endpoint, not a frozen protocol. Don't treat presence of code as admission of a
  contract.
- `internal/execution/lifecyclestate` — durable attempt-lifecycle types colocated with Supervisor
  authority state (see ADR-0025).
- `internal/execution/registrationstate` — the fixed local file-snapshot store for registrations,
  approvals, and attempts, including the v1→v2 migration and archive slices (F1–F4 in
  `docs/PROJECT.md`/`docs/SUPERVISOR_ARCHIVE_*`).
- `internal/execution/approvalattempt` — passive typed domains, a closed classification
  vocabulary, and a bounded fixture-only verifier.
- `internal/execution/registeredlifecycle` — the **fake**, no-guest lifecycle: it never launches a
  real guest. It exists to exercise registration/approval/attempt plumbing before a real backend is
  admitted. Don't mistake "lifecycle passes" here for "execution works."
- `internal/execution/archivestate` — immutable retained-segment archive for completed
  registration cohorts (ADR-0031); a local conformance oracle, not a production storage engine.
- `internal/execution/installationowner` — per-installation ownership/locking mechanics (the
  owner-lock G-series work, ADR-0033).

## Glossary

- **Trust epoch**: a signed, sequence-ordered binding of component identity and policy. Trust
  changes are explicit and crash-safe, never implicit.
- **DID (Decentralized Identifier)**: an optional *external* representation of a Capsule
  principal. It is never an authority or trust source on its own — trust comes from local key
  authorization, purpose binding, installation identity, and trust epoch.
- **Fixed store**: this repository's term for its finite, local, conformance-oracle persistence
  format. It is explicitly not a production database engine.
- **Conformance fixture / known-answer test**: a byte-exact input/output vector proving an
  encode/decode/migration/lifecycle transition is correct. This repository has hundreds of these
  already (`docs/PROJECT.md` cites a 94-rule/458-case/561-fixture corpus for one slice alone) —
  the expectation is to grow this corpus alongside new protocol code, not to substitute ad hoc
  assertions for it.
- **Classification**: a closed, internal, content-free validation-outcome vocabulary (see
  `internal/execution/lifecyclestate/errors.go`'s `Classification` type). Not a public protocol
  error, not something to `%v`-print to a user.
- **Status language**: the required vocabulary for describing work state —
  `PASSED`, `IN_PROGRESS — TRENDING_GOOD`, `IN_PROGRESS — TRENDING_BAD`, `BLOCKED`, `NO_GO`
  (`docs/STATUS_LANGUAGE.md`). Use this exactly in comments, commit messages, and PR descriptions
  that describe status — not "done," "working," or "complete."
- **Spike / experiment code**: disposable, non-production. Product packages must never import it;
  retained evidence lives in the separate `Shrimpworks/capsule-experiments` archive, pinned to an
  exact commit.

## Domain-specific anti-patterns

Things that are wrong specifically *because of this project's threat model*, even if they would be
unremarkable in an ordinary Go service:

- Giving anything other than the Execution Supervisor the ability to create, terminate, or
  reconcile a guest.
- Accepting plan bytes, backend flags, images, mounts, or guest paths at execute time instead of a
  Supervisor-issued registration ID.
- Treating a device identifier or DID as authority instead of as an identifier.
- Passing a live host path into a guest instead of an immutable, content-addressed snapshot.
- Silently clamping or silently expanding a user-set resource limit — resolve defaults before
  approval, reject anything above the user's ceiling, then enforce exactly or refuse the job.
- Adding rich document/archive/image/media parsing inside the daemon or Supervisor instead of a
  bounded Broker validator or a future disposable parser sandbox.
- Adding a new Supervisor responsibility or privileged helper without an ADR.
- Describing a backend, profile, control, or security tier as implemented, validated, secure, or
  production-ready without its exact mechanism and retained adversarial evidence to back that claim.

## Where new Go work usually starts

For a change to execution, policy, protocol, or identity: read `AGENTS.md`, then the ADR(s) that
govern the exact component, then the current status paragraph in `docs/PROJECT.md` so you know
whether the area is `PASSED`, `IN_PROGRESS`, or `BLOCKED` before extending it. For an ordinary
internal refactor with no behavior change, `docs/GO_ENGINEERING_STANDARDS.md` and the package's own
`doc.go` are usually enough context.
