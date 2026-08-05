# Supervisor archive F4C bounded immutable-segment growth result

Status: `PASSED` for the exact fixed-store v2 local-conformance scope described here.

Date: 2026-08-05

Parent status: `IN_PROGRESS — TRENDING_GOOD`. F5 backup/restore and orphan-cleanup policy, F6
production-engine selection, real power-loss evidence, product admission, and every adapter,
runtime, backend, guest, IPC, service, deployment, and consumer integration remain `BLOCKED`.

## Defensive scope

This slice used only repository-owned fixed-store code, deterministic generated fixtures, owned
temporary directories and files, local fault/process-death injection, and the existing fake
no-guest lifecycle binding. It accessed no other system, identity, credential, user data, runtime,
backend, VM, guest, service, IPC consumer, or deployment. Archive code made zero adapter calls and
the retained fake binding has `CreatesGuest() == false`.

## Implemented growth contract

The owner-required v2 archive transaction now deterministically selects, prepares, fully verifies,
publishes, activates, reopens, and resolves a second and later immutable segment. Ordinal is the
prior descriptor count plus one, archive generation is the prior generation plus one, source
snapshot generation is the freshly verified active generation, and each segment binds the exact
current checkpoint digest as its predecessor. Activation preserves every earlier descriptor and
segment byte, publishes the new digest-addressed segment before installing its reference, and
atomically advances the active checkpoint.

Every reopen reconstructs the retained-global indexes from the current hot records, the
independent append-only F4B effect-tombstone collection, and all referenced segments. Typed hot and
archive locations, lifecycle `absent | present` joins, visible-v1 seed facts, complete registration
cohorts, approval/attempt replay identities, instance identities, and hot-only `AttemptID` recovery
remain exact. Historical effects resolve as `superseded-by-current` without an invented issuing
lifecycle; only the exact current effect returns the retained lifecycle.

The first F4C checkpoint domain-binds the fifth `HotSetDigests.EffectTombstones` digest and
materializes that digest in the activation checkpoint. The encoder retains the legacy four-digest
checkpoint encoding only when the fifth digest is zero, so the exact F2 migration, F3 first-
activation, and F4B mutation bytes are unchanged. No old byte sequence is reinterpreted.

## Exact second-segment answers

The deterministic two-cohort fixture archives one complete cohort per segment. Its second
activation has snapshot generation `3`, archive generation `3`, segment ordinal `2`, two
referenced segments, zero hot cohorts, and archived counts of `2` registrations, approvals,
attempts, lifecycles, nonces, effects, instances, approval-replay entries, and attempt-replay
entries.

| Answer | SHA-256/domain digest |
| --- | --- |
| Active snapshot file | `22d0ab0d3a4d06ab507ed5248b0645591734f2162709d791fd83a2160da366b3` |
| Second segment file | `1aec9165d1026ee312e32043c1aa85dbd06c8804e8b439d3e05e59fa8210316b` |
| Second segment semantic identity | `4adcd3d80596d6e17e21f29f197f39658df0e1debc9ee4851f7d4f72bcdce7d7` |
| Second activation checkpoint | `af26eaee606097864da7bdc2e311a193808d63d5fdce42e1b9b32f566c313374` |
| Retained-global combined index | `7d139e5ce710c54dfead1e3a6cd3fbb1adb50e56a2af428cccc04b9f9956b6c3` |

The first segment file is byte-identical before and after second activation. Both archived
registrations, exact approval replays, attempt replays, lifecycle joins, and recovery exclusion
resolve after reopen. A mutation-between-activations fixture retains seven total archived effects:
one visible-v1 seed from the first cohort plus six independently issued v2 effects from the second
cohort. Its first v2 effect resolves historical, its destroy effect resolves current with the exact
destroyed lifecycle, and the hot effect count is zero after the second activation.

## Bounds and refusal evidence

The existing exact caps remain authoritative: `64 MiB` active v2 bytes, `64 MiB` per segment,
`256` cohorts per segment, `4,096` approvals, attempts, and lifecycles per segment, `64` referenced
segments, and `262,144` entries per retained-global index. Candidate encoding is completed during
prepare, before any segment publication, so active cap-plus-one refuses before bytes. Existing
inclusive cap/cap-plus-one oracles remain unchanged.

The focused maximum-growth oracle activates exactly `64` one-cohort segments, reaches archive
generation `65`, retains one final eligible cohort hot, and fully verifies. Planning segment `65`
then returns `CAPACITY` and leaves the active file and complete archive-directory inventory
byte-identical. No path silently evicts, deletes, compacts, merges, scans for fallback, or accepts a
caller-selected location.

Fresh full verification refuses stale generation/head, wrong owner session, a losing competing
candidate, duplicate identity, duplicate cross-segment location, partial index, partial counts,
missing/corrupt/substituted second segment, wrong predecessor, stale descriptor ordinal or
generation, and partial or corrupt tombstone worlds without rewriting retained evidence. Exact
concurrent retries converge to one activation and stale refusals for every loser.

## Publication, fault, and death evidence

All thirteen inherited activation fault points run against the second activation. Pre-publication
confirmed aborts preserve the complete old world; post-publication/pre-activation faults preserve
the old world and report exactly one valid orphan; post-activation faults reopen the complete new
world and report no orphan. Four subprocess-death points prove the same old-or-complete-new split
before segment publication, after segment publication, before activation, and after activation.
The first referenced segment and active predecessor bytes never change on a refused outcome.

The retained F2, F3, F4A, and F4B known-answer and fault suites continue to pass. F4C changes no
adapter or fake-backend behavior and creates no guest.

## Verification

Focused verification:

```sh
env GOCACHE=/tmp/capsule-f4c-go-cache go test ./internal/execution/archivestate ./internal/execution/registrationstate -count=1
env GOCACHE=/tmp/capsule-f4c-go-cache go test -race ./internal/execution/registrationstate -run 'TestFixedStoreV2Second' -count=1
env GOCACHE=/tmp/capsule-f4c-go-cache go test ./internal/execution/registrationstate -run '^TestFixedStoreV2SecondSegmentKnownAnswerLookupReplayAndRecovery$' -count=20
```

Full repository verification passed:

```sh
pnpm install
pnpm check
pnpm lint
pnpm test
pnpm verify:schemas
pnpm verify:adrs
pnpm format:check
go test ./...
go vet ./...
go build ./...
golangci-lint run ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
git diff --check
```

`golangci-lint` reported `0 issues`. `govulncheck` found zero called vulnerabilities; it reported one
vulnerability in the required module graph that this code does not call.

## Boundary after F4C

This is a finite fixed-store conformance oracle. It is not a production database, continuous-
service or multi-process claim, real power-loss result, coherent backup/restore mechanism,
anti-rollback design, secure-deletion policy, or product-admitted archive. It does not delete even
known orphans. F5 owns backup/restore and the narrow orphan-cleanup rule; F6 owns production-engine
selection and quantitative durability/resource evidence. ADR-0031 remains Proposed.
