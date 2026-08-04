# Supervisor archive F2 stateful migration result

Status: `PASSED` for the defensive repository-local fixed-store v1-to-v2 migration and empty-
archive full-verifier scope described here.

Date: 2026-08-04

This result uses only repository-owned state fixtures, owned temporary files/directories, injected
file faults, and local subprocesses. It invokes no lifecycle adapter, creates no backend process or
guest, moves no cohort, publishes no segment, routes no retained lookup, and adds no consumer or
product IPC behavior.

The applicable persistence/archive rows in
[Ecosystem reuse and adoption](ECOSYSTEM_REUSE_AND_ADOPTION.md) remain `BUILD-NARROWLY` for this
fixed oracle and `SPIKE-FIRST` for a later engine. F2 adds no dependency and makes no SQLite or
other production-engine selection.

## Exact predecessor and successor

The only predecessor is the closed v1 JSON envelope:

- `storeFormatVersion = 1`, `migrationSourceVersion = 0`;
- the complete existing `installationState` with registration, approval, and attempt set digests;
- `lifecycleSetDigest`; and
- a non-null lifecycle-record array.

The only successor is the closed v2 JSON envelope implemented by
`registrationstate.FixedFileStoreV2`:

- `storeFormatVersion = 2`, `migrationSourceVersion = 1`;
- `snapshotGeneration = 1`, `archiveGeneration = 1`;
- the byte-equivalent v1 authority state, lifecycle set digest, and sorted lifecycle records;
- a non-null empty descriptor array and exact empty descriptor-set digest;
- complete sorted `retained-global` registration, approval, attempt, nonce, effect, instance,
  approval-replay, and attempt-replay indexes with per-family and combined digests;
- record-kind-bound `hot` locations for every full record and no archive locations;
- an explicit lifecycle `absent` or `present` arm on every attempt entry, with lifecycle count
  derived only from present arms;
- `visible-v1-seed-plus-all-v2-issued` effect coverage, retaining every nonzero effect ID still
  visible in v1 and inventing none;
- no previous checkpoint and one exact `migration-genesis` current checkpoint; and
- exact hot, archived, and total counts, with every archived count equal to zero.

The v1 opener refuses v2, the v2 opener refuses v0/v1, and neither creates, rewrites, or falls back
to another version. The v2 full verifier closed-decodes the entire file, validates the unchanged v1
authority/lifecycle world, reconstructs every index and digest from full records, and compares the
complete active projection and migration genesis. An absent lifecycle remains absent and has no
effect or instance entry.

## Retained known answers

The exact literals are asserted by `TestFixedStoreV2KnownAnswers`:

| World | File SHA-256 | Combined retained-global index | Migration genesis | Visible-v1 seed | Counts `(registration, approval, attempt, lifecycle, nonce, effect, instance, approval-replay, attempt-replay)` |
| --- | --- | --- | --- | --- | --- |
| Empty | `c845344aa1b464d7cd40ba86d33ecf3e5797cfed8406767dff0db12be9d9bb04` | `78e817b6a07989095010743601a017e43e3b660ea78ad0231e01d900227e207c` | `37336388a2f775b79d4f56e1af8ff1afd45de6cea96a64bfc27cec564290c88c` | `17de5f44f523dab94ca4b215ce7779358146fb094fa6d208e0190cb0ba69e0a1` | `(0,0,0,0,0,0,0,0,0)` |
| One committed attempt, lifecycle absent | `569ba7c1aa25432a1001b1ca7122a7772ccfe954a18a484d26ed71b4255b8dca` | `924c78b9508123feb1b78fd62b71df6cace9c97b4887f2d67bbdc6ef2a9a7de5` | `983c2474dbef1fa6908de0fa02f96e2aaf7245bce48399e224a0f8a2c349a23e` | `17de5f44f523dab94ca4b215ce7779358146fb094fa6d208e0190cb0ba69e0a1` | `(1,1,1,0,1,0,0,1,1)` |
| One observed lifecycle with effect and instance | `29c706be5d8a55958acacae7ad01a001576a307d847fd610a6f2c0e57f291235` | `2dc21df5c66bdb46f7ba80ed6566b12912a6b9b79a1f5b7a3d325f507bea65c2` | `3b5de809e6cbcda85b94439a9142cdc69243b1d18f912dbe8fc81ba6e4101a99` | `acee3fa25e62c185eb1e9b26313b6f04b9c88e57857e0fb400371c2c5a67295f` | `(1,1,1,1,1,1,1,1,1)` |

The second row is the required valid v1 crash witness. Migration preserves its one committed
attempt and zero lifecycle records exactly; it runs no recovery ceremony and invents no lifecycle
state.

## Commit and fault result

Migration requires the asserted installation owner at entry, immediately before commit, and before
reopen. It fully validates v1 before creating a temporary file. The sole commit sequence is:
temporary file in the same directory, mode `0600`, complete write, file sync, close, atomic rename,
directory sync, ownership recheck, and version-specific full reopen verification.

Injected failures through the pre-rename boundary preserve the source v1 file byte for byte.
Injected post-rename and post-directory-sync outcomes return indeterminate; a fresh opener then sees
one complete v1 world or one complete v2 world, never a merge. The retained subprocess death oracle
observes v1 before rename and complete v2 after rename. Eight concurrent migration contenders
serialize: exactly one succeeds and the rest refuse the now-v2 source.

The focused corpus covers zero/missing/mixed/all lifecycle states; exact 256-active and
4,096-retained limits plus cap+1; unsafe file shapes; source corruption before write; successor
truncation, unsupported/missing/duplicate/unknown/trailing fields; digest, sorted-index, count,
location, cross-link, generation, checkpoint, capacity, time-high-water, and invented absent-arm
effect/instance faults; deterministic serialization; no-rewrite refusal; and defensive-copy read
projections.

## Boundary after F2

This is a finite fixed-store conformance oracle, not a production store or product-admitted
archive. Slice F3 segment prepare/publish/activation, Slice F4 retained lookup and v2 mutations,
Slice F5 backup/orphan/offline-report work, and Slice F6 production-engine comparison remain open.
No production durability, real power-loss, backup/restore, rollback resistance, continuous service,
runtime/backend/guest security, referenced-history deletion, or product-admission claim advances.
