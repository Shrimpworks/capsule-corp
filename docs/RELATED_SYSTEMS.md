# Related systems and design influences

Status: public-source design comparison and planning input, not implementation evidence, runtime
admission, a competitor-completeness claim, or a legal novelty/patentability opinion.

Capsule is not trying to provide a general-purpose computer for an agent. Its intended boundary
turns one untrusted proposal into exact registered bytes, separate human authorization, one
attempt, bounded content and execution authority, durable cleanup, and role-attributable evidence.
No single public project reviewed on 2026-08-03 supplied that complete composition. Individual
mechanisms do have strong precedents, and the useful lessons below refine Capsule without changing
its authority split or current no-admission posture.

Public projects change independently of Capsule. Links identify the reviewed design source, not a
pinned dependency or retained Capsule evidence. Any mechanism considered for implementation still
requires an exact version, threat review, bounded experiment, and the applicable ADR and
control-evidence update.

## Comparisons that directly reinforce Capsule

| System | Useful lesson | Capsule boundary and limitation |
| --- | --- | --- |
| [Deno permission broker](https://docs.deno.com/runtime/fundamentals/security/#permission-broker) | Version requests and decisions, correlate monotonic request IDs, and terminate on malformed, mismatched, reordered, or interrupted broker exchanges. | This is a full-Deno feature, not part of Capsule's governed `deno_core` candidate. Capsule may reuse its protocol-failure ideas, but v0 adds no dynamic permission grant and continues to require an external guest boundary. |
| [HARP](https://harp-protocol.github.io/) | Render the exact artifact being authorized, bind the decision to its byte identity, purpose, scope, expiry, and replay state, and keep the approver outside the proposing interface. | HARP v0.2 is a draft approval protocol. It supplies no hostile-guest, content-custody, backend-lifecycle, cleanup, or enforcement-evidence boundary and is not a Capsule dependency. |
| [Redan](https://github.com/getredan/redan/blob/main/docs/security-model.md) | Exhaustively classify configuration fields that can exercise host authority; require explicit trust for the exact configuration bytes; make new fields fail closed until classified. Its host-mediated network and credential model is also a useful future reference. | Redan's alpha design intentionally accepts a broader development-machine and same-user posture, including a writable workspace. Capsule does not adopt that trust model, raw workspace exposure, response scrubbing as a confidentiality boundary, or a combined launch/credential authority. |
| [Qubes qrexec](https://doc.qubes-os.org/en/r4.3/developer/services/qrexec.html) | Expose typed cross-domain services under default-deny policy. Human confirmation may select only among policy-permitted targets; it cannot turn a denied target into allowed authority. | Capsule already requires the same rule: approval may accept, reject, or select a predeclared Supervisor-validated alternative, but any broadened plan requires new registered bytes and approval. |
| [XDG document portal](https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.Documents.html) | Let trusted UI select content and expose a scoped document identity and revocable grant instead of general filesystem authority. | Capsule needs a stronger same-user model: the Broker creates an immutable snapshot, withholds original host paths, transfers only attempt-scoped access, and durably records expiry, release, quarantine, and cleanup. The portal's FUSE implementation is not selected. |
| [Confidential Containers policies](https://confidentialcontainers.org/docs/attestation/policies/) and [Kata init-data advisory](https://github.com/kata-containers/kata-containers/security/advisories/GHSA-989w-4xr2-ww9m) | Keep evidence production, policy evaluation, and resource release separate. Never treat an I/O or verification error as confirmed absence. | Capsule's durable lifecycle already distinguishes confirmed, unsupported, failed, indeterminate, and repair-required outcomes. Real adapters, content custody, and artifact release must preserve that distinction. |
| [Bazel Remote Execution](https://github.com/bazelbuild/remote-apis), [BuildKit LLB](https://docs.docker.com/build/buildkit/), and [Nix](https://nixos.org/guides/how-nix-works/) | Resolve a flexible request into an immutable content-addressed action and dependency closure before execution. | Future package support should use a separately governed preparation job that emits a `PreparedClosure`. The approved execution attempt remains network-free, package-manager-free, and bound to the complete closure digest. |
| [DSSE](https://github.com/secure-systems-lab/dsse) and [in-toto](https://in-toto.io/) | Give every signed statement an explicit type and role, bind materials and products by digest, and support small offline verification. | Capsule retains bounded CBOR/COSE candidates and its own purpose-separated keys. The lesson is statement semantics, not a format migration or an attestation claim. |
| [Gondolin](https://github.com/earendil-works/gondolin) | Synthetic host-provided resources can avoid a direct writable host mount and can keep guest-visible names separate from host paths. | Its experimental TypeScript control plane and broader development environment are transport references only, not Capsule's authority model. |

## Capability-oriented references

[SES](https://docs.endojs.org/modules/ses.html),
[LavaMoat](https://github.com/LavaMoat/LavaMoat), the
[WASI Component Model](https://component-model.bytecodealliance.org/design/worlds.html), and
[Fuchsia capability routing](https://fuchsia.dev/fuchsia-src/concepts/principles/secure) reinforce
explicit endowments, typed imports/exports, and handle-based capability routing. They may inform
later package analysis, defense-in-depth, and alternative runtime profiles. None replaces the
microVM boundary, and none should enter the v0 trusted computing base without a separate need and
review.

Agent sandbox products such as
[Anthropic Sandbox Runtime](https://github.com/anthropic-experimental/sandbox-runtime),
[Docker Sandboxes](https://docs.docker.com/ai/sandboxes/security/isolation/), and hosted execution
services are useful product benchmarks for denial explanations, cancellation, streaming,
artifacts, and repository-change adoption. Their usual abstraction is a capable computer or
workspace. Capsule's first abstraction remains one bounded task.

## Refinements retained in Capsule planning

### Field-authority classification

Before a target protocol object freezes, every field must have a machine-readable classification
covering:

- the role permitted to originate it;
- the component that validates or resolves it;
- whether it changes, narrows, selects, or merely describes authority;
- whether trusted approval UI must render it;
- whether it contains user content or guest-controlled material;
- the exact digest, signature, registration, or record that binds it; and
- whether an unknown value or field rejects the object.

The manifest complements strict unknown-field decoding. Strict decoding protects the running
system from an unknown wire field; classification protects development from adding a known field
without deciding its authority consequences. Once the manifest exists, CI must reject a target
field that has no closed classification. The first implementation should cover current candidate
objects and expand only with coordinated object-model migrations.

### Approval cannot widen

The Approval Broker may approve or reject the exact registered plan. A future selectable option
may choose only from a closed set already represented in the registered bytes and independently
validated by the Supervisor. Approval cannot add a destination, path, capability, output audience,
backend option, resource increase, or other authority. A change requires a new plan, registration,
and approval. This is an existing Capsule invariant; future approval fixtures must make it
directly falsifiable.

### Error is not absence

Security-sensitive observations use explicit closed outcomes. An I/O failure, short or corrupt
read, invalid signature, digest mismatch, stale epoch, unsupported version, unavailable report,
or post-restart uncertainty never becomes confirmed absence or success-with-warning. The existing
durable lifecycle applies this rule to fake effects; the same oracle must apply at every future
backend, content, and artifact boundary.

## Deferred profiles, not v0 expansion

- **Content lifecycle:** selected, snapshotted, digested, registered, leased, staged, consumed or
  expired, collected, validated, released or quarantined, and retained or destroyed. Gate D and
  the later content phases own the exact storage and crash-recovery contract.
- **Prepared dependencies:** a separately isolated preparation profile may resolve dependencies
  and emit an immutable `PreparedClosure`; live approved execution never installs packages.
- **Receipts:** plan registration, approval, Supervisor enforcement, artifact custody, and final
  composition remain separate attributable statements with a future small offline verifier.
- **Network:** a future profile must specify destinations, protocols, DNS, redirects, private and
  metadata ranges, methods, paths, credential purpose, byte/concurrency/time limits, and audit
  fields. V0 remains no-network.
- **Repository work:** a future repository-agent profile operates on an immutable snapshot,
  private clone, or private worktree and returns a patch, commit bundle, or tree identity for
  deliberate adoption. It never receives a writable mount of the user's actual repository and
  Git control state.

## Patterns deliberately not adopted

- a generic privileged `run(command)` authority surface;
- a writable mount of the user's real workspace;
- guest environment variables or files containing host credentials;
- destination allowlists, DNS policy, or response scrubbing presented as complete egress or
  confidentiality controls;
- persistent or warm hostile guests by default;
- dependency installation during the approved execution attempt;
- a runtime permission system presented as the host isolation boundary;
- one component that presents approval, holds user content or credentials, launches the guest,
  and declares enforcement success; or
- a same-user writable trust file presented as resistance to a malicious same-user process.

## Planning consequence

This comparison creates no new runtime or backend admission blocker and does not reorder the
current governed-fork, TypeScript approved-byte, authenticated IPC, archive/owner-lock, or libkrun
work. The field-authority manifest is the only new near-term protocol control. The other retained
lessons either make existing invariants more directly testable or inform later content,
preparation, evidence, network, and repository profiles.
