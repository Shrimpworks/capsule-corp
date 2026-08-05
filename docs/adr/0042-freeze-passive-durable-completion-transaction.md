# ADR-0042: Freeze the passive durable completion-last transaction

- Status: Proposed
- Date: 2026-08-05
- Refines if accepted: ADR-0015, ADR-0024, ADR-0025, and ADR-0040

## Context

The read-only FakeBackend compositor proves that a terminal projection can be derived from a
committed attempt, exact registered plan, immutable lifecycle bindings, one completion-last fact,
and authoritative absence. It deliberately supplies no producer or durable completion store.
Without that boundary, a caller could treat JSON bytes, stream EOF, runner exit zero, or a daemon
response as terminal success, and response loss could create a different public result on retry.

The first implementation remains unwired and no-guest. No admitted runtime, runner, completion
transport, process tree, guest, signing key, authenticated endpoint, or production storage engine
exists. FakeBackend therefore cannot supply a real runner identity or real descendant-absence
evidence. Those fields must remain explicitly unresolved or fake-scoped rather than being inferred.

## Decision

Capsule selects one Supervisor-owned passive transaction for the local FakeBackend oracle. It reads
the committed created attempt and lifecycle record by `AttemptID`, independently revalidates their
registration, plan, source-manifest, runtime/profile, backend, and immutable-binding links, validates
bounded typed JSON and its SHA-256 digest, composes the bounded transcript and fixed summary, then
atomically publishes one immutable completion bundle with `committed-last`. No public result is
returned before that durable publication.

The exact stored object is `capsule.unwired-fake-durable-completion` version 0. The exact store is
`capsule.unwired-fake-completion-store` version 0 and binds one nonzero installation ID at its
top level. Every retained record must bind that same installation. Each record binds:

- attempt, approval, registration, plan, installation, epoch, Supervisor, approval-key
  authorization, and source-manifest identities;
- runtime-bundle, profile-registry, backend-validation/configuration/implementation/instance, and
  immutable lifecycle-binding digests;
- exact typed-JSON bytes, length, SHA-256 digest, and
  `supervisor-validated-typed-json-sha256` disposition;
- lifecycle `destroyed`, operation sequence 6, fake destroy confirmation, cleanup false, and fake
  authoritative absence;
- `unresolved-fake-backend-has-no-runner` instead of a fabricated runner identity;
- `completed-fake-no-guest-local-mechanic`, which is not product success or attestation; and
- exact transcript, transcript SHA-256, fixed summary, completion-record digest, durable-record
  digest, positive commit sequence, and `committed-last` marker.

Inclusive caps are 262,144 result JSON bytes, 4,096 transcript bytes, 256 summary bytes, 368,640
encoded durable-record bytes, 4,096 retained records, 67,108,864 aggregate retained result bytes,
and 100,663,296 store bytes. The typed-JSON decoder also retains the compositor's exact depth,
node, member, element, key, integer, UTF-8, duplicate-key, and trailing-data bounds.

The store writes and synchronizes a complete temporary snapshot, atomically renames it, and
synchronizes the parent directory. A failure before rename commits nothing. An error after rename
is indeterminate, permanently fences that store instance, and requires reopen. Reopen accepts only
one complete old or new version. Missing, early, duplicate, mixed-version, malformed, forged,
cross-linked, or digest-invalid completion state is recovery-required and is never rewritten.

The first committed record for an `AttemptID` is immutable. Exact replay returns the same stored
result, transcript, digest, and summary. A changed result for that `AttemptID` is `REPLAY` and cannot
replace the record. EOF, runner exit, guest prose, diagnostics, timing, paths, and artifact names are
not commit inputs. Nonterminal or unresolved lifecycle state commits no completion, releases no
output, and cannot be used as capacity-release evidence.

The fixed completion file is a bounded local conformance oracle, not a second selected product
database. It requires external single-owner composition and supplies no installed protected root,
owner lock, rollback defense, archive, backup, restore, multi-process, APFS power-loss, continuous-
service, or production-engine claim. A product store must integrate this logical transaction with
the selected Supervisor authority/lifecycle engine rather than treating this file as admission.

## Consequences

- Only a durable Supervisor transaction can establish the fake-local terminal object and public
  projection; EOF and exit zero remain non-authoritative observations.
- Lost responses and restart converge on one immutable public result for the same `AttemptID`.
- Exact source, runtime/profile, result, lifecycle, teardown, and absence bindings are inspectable
  without adding guest-controlled names, paths, diagnostics, or timing to the summary.
- FakeBackend's missing runner/process-tree facts remain explicit limitations instead of invented
  evidence.
- Production completion transport, trusted launcher/result-integrity evidence, real runner identity,
  real authoritative descendant absence, signing, receipt composition, installed ownership, and
  product-store integration remain blocked follow-ups.
