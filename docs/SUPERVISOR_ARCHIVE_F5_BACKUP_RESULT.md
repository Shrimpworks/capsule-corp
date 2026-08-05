# Supervisor archive F5 coherent-backup and offline-reporting result

Status: `PASSED` for the exact owner-held fixed-store v2 local-conformance scope described here.

Date: 2026-08-05

Parent status: `IN_PROGRESS — TRENDING_GOOD`. F6 production-engine selection, installed protected-
root/owner evidence, real APFS/power-loss durability, a restore/repair activation ceremony,
independently protected checkpoint production, referenced-history retention/deletion policy,
continuous service, product consumers, and every adapter/runtime/backend/guest integration remain
`BLOCKED`.

## Defensive scope

F5 uses only repository-owned authority/lifecycle/archive fixtures, owned temporary files and
directories, injected local store/filesystem faults, and bounded local subprocess-death tests. It
calls no adapter or lifecycle driver, creates no process other than its test subprocesses, and
accesses no runtime, backend, VM, guest, product IPC/consumer, signing identity, credential,
external system, or user data. It adds no dependency and does not begin F6.

## Coherent backup format and transaction

The internal `registrationstate` F5 API accepts an already-created owner-only `BackupRoot`
capability, never a public/IPC path. Under the same live `ArchiveOwner` session used for archive
mutation, `CreateCoherentBackup` fully verifies the source v2 world, copies exact active bytes as
`state.json`, copies every referenced digest-addressed immutable segment under
`state.json.archive`, writes the closed `manifest.json` last, syncs every copied file and containing
directory, then independently reopens the complete copy.

Backup-manifest v0 binds all of the following:

- exact active and segment file names, byte lengths, SHA-256 file digests, segment semantic
  digests, ordinals, source/archive generations, and predecessor checkpoint digests;
- installation, Supervisor, epoch sequence/digest, durable time high water, snapshot/archive
  generations, previous/current checkpoint references, and migration-genesis digest;
- descriptor-set, retained-global combined-index, hot registration/approval/attempt/lifecycle and
  independent effect-tombstone set digests;
- visible-v1 effect seed count/digest and exact hot/archived/total counts; and
- a domain-separated digest over the complete manifest with its digest field zeroed.

The retained one-segment F3-world manifest digest is
`deaec00c9500340be0f56a306deade47d262d0b0f4bce18f9efe26a8512195d2`.

A destination without the final manifest is incomplete, never a coherent backup. Unknown or extra
destination entries, unsafe file shapes, missing/substituted/oversized files, manifest corruption,
and mixed checkpoint/generation content refuse without rewrite. Repeated verification is
deterministic and read-only. Backup creation never changes the live predecessor.

## Restore-admission and rollback semantics

F5 intentionally implements read-only restore admission, not restore activation. It first verifies
both the complete live predecessor and complete candidate backup under the owner session. Exact
candidate equality with an injected `IndependentLatestCheckpoint` fixture is classified only as
`eligible-for-future-restore`; no bytes are copied over the live store, no epoch changes, and no
attempt becomes enabled.

Missing, older, future, or incomparable anchors are `rollback-uncertain` and refuse admission.
Without an independently protected latest checkpoint, even a structurally and cryptographically
valid coherent backup can be an older world that hides later registrations, approvals, attempts,
nonces, effects, instances, replay identities, lifecycle joins, and retained-global counts. F5
therefore makes no rollback-prevention or global non-reuse claim. An actual forward restore/repair
ceremony remains a later ADR and implementation.

At every F5 failure or response-loss boundary the live world remains the fully verified predecessor;
because F5 has no activation operation, it never exposes a partially restored successor or hybrid.

## Offline inventory, reporting, and orphan handling

`VerifyArchiveSet(..., VerificationFull)` returns one deterministic fixed-shape report containing
only versions, installation/Supervisor/epoch identity, generations, time high water, checkpoint,
counts, segment count, and bounded classifications. It emits no filesystem path, plan, approval,
stored label, user content, guest text, or arbitrary decoded error.

`InventoryArchiveArtifacts` classifies archive-directory entries as referenced,
known-unreferenced, backup-referenced, unknown, corrupt, mixed-generation, or cross-installation.
Unknown temporary/name/type/link/mode/owner shapes and corrupt/mixed/cross-installation artifacts
refuse every deletion candidate and preserve bytes as evidence. Inventory never adopts, merges,
renames, rewrites, or deletes an orphan.

The sole deletion API requires an explicitly selected sealed `KnownUnreferencedOrphan` issued by
the inventory. It rechecks the unchanged active checkpoint, owner session, complete live reference
set, every supplied coherent-backup manifest, digest-addressed name, closed segment decode, exact
same-installation next-generation ancestry, and device/inode before unlink. It deletes no
referenced segment, history, tombstone, checkpoint, backup member, unknown file, or temporary file.
An indeterminate unlink retry is idempotent only for that sealed candidate and unchanged checkpoint.
Directory-sync or response loss never changes active authority. This is physical free-space
recovery for one known artifact, not secure deletion.

## Retained fault and boundary evidence

Focused tests retain:

- prepare, active/segment/manifest temporary create, copy/write, sync, close, reopen, rename, and
  directory-sync faults;
- owner loss at entry and immediately before manifest completion, manifest-last response loss, and
  process death after active, segment, and manifest publication;
- missing, extra, substituted, cap-plus-one manifest, rollback, future/incomparable anchor, and
  mixed-backup cases;
- exact one-, two-, and 64-segment coherent backups, inherited segment-65 no-rewrite refusal, and
  twenty-run repeated exact verification;
- known-orphan pre-unlink refusal, post-unlink/directory-sync/reopen response loss, and idempotent
  exact recovery; and
- unknown/corrupt evidence preservation with no deletion candidate.

All F2-F4C known answers, mutation/replay/recovery behavior, immutable segment bytes, visible-v1
limitations, independent effect-tombstone semantics, and fixed capacity rules remain unchanged.

## Verification

The final 2026-08-05 worktree passed the complete contributor gate with Node.js 22.22.1,
pnpm 10.28.2, and Go 1.26.5:

- `pnpm install`, `pnpm check`, `pnpm lint`, `pnpm test`, `pnpm verify:schemas`,
  `pnpm verify:adrs`, and `pnpm format:check`;
- `go test ./...`, `go vet ./...`, `go build ./...`, and `golangci-lint run ./...` with zero
  lint findings;
- `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`, with zero called
  vulnerabilities; and
- `git diff --check`.

Focused F5 evidence additionally passed the full backup/orphan selection, a race-enabled fault and
contention selection, five process-death repetitions, twenty deterministic known-answer
repetitions, exact-cap/cap-plus-one coverage, and a final post-format focused run. The full
`registrationstate` package completed in the repository-wide suite; no test was skipped to obtain
the result.

## Boundary after F5

This is a finite local fixed-store conformance result, not a production database, backup product,
restore mechanism, continuous-service system, APFS durability claim, anti-rollback anchor, or
secure-deletion result. Backup discovery/retention configuration and the injected checkpoint are
fixtures, not installed authority. F6 must separately select and validate a production engine and
exact configuration; it is not started by F5. ADR-0031 remains Proposed.
