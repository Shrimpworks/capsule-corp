# Supervisor archive F6 SQLite research and execution checkpoint

Status: `PASSED` for this documentation-only research and experiment-planning scope.
F6 execution, production persistence, SQLite or binding selection, restore activation, an external-
alpha consumer, continuity, and product admission remain `BLOCKED`.

Research date: 2026-08-05. This checkpoint is the canonical input for one separately authorized
experiment in `Shrimpworks/capsule-experiments`. It adds no dependency, binding, store, SQL
implementation, or consumer and did not build or run SQLite.

## Decision boundary and ADR lifecycle

[Proposed ADR-0031](adr/0031-checkpoint-closed-supervisor-cohorts.md) already authorizes a bounded
production-engine comparison after the fixed-store oracle. A second ADR is not needed to authorize
the experiment. ADR-0031 therefore remains **Proposed**. This checkpoint recommends one research
candidate and freezes an executable evidence packet; it does not select SQLite, a Capsule-owned
cgo binding, or any product dependency.

After the experiment, a separate ADR must either select one exact engine/configuration or reject
the tested path. That ADR starts Proposed and cannot become Accepted merely because the experiment
ran. F6 accepts evidence and a decision only: it does not cut over the fixed store, activate
restore, delete history, or wire a consumer in the same slice.

Defensive scope: validate Capsule's durable-authority, archive-publication, replay, corruption,
locking, backup, and refusal controls using only repository fixtures and an explicitly named owned
disposable macOS/APFS experiment environment. Do not access any other system, identity,
credential, or data, and preserve Capsule's existing safeguards. Real power interruption requires
separate authorization naming the disposable host and volume.

## Research recommendation, not selection

The first F6 comparison should use the official SQLite amalgamation behind a Capsule-owned,
experiment-only cgo shim. This is the narrowest path that preserves exact upstream C bytes, exposes
the native Darwin VFS and result codes under test, and permits a fixed function/statement surface.
It is `SPIKE-FIRST`, not `ADOPT-PINNED`.

### Exact upstream material

The current primary release at the research date is
[SQLite 3.53.4](https://sqlite.org/releaselog/3_53_4.html), released 2026-07-24:

| Material | Immutable identity |
| --- | --- |
| source ID | `2026-07-24 19:02:57 bf7c7f30031888f4e796e429ab3978879485813aaca6f641c7b33e4e09459bcc` |
| Fossil check-in | [`bf7c7f30031888f4e796e429ab3978879485813aaca6f641c7b33e4e09459bcc`](https://sqlite.org/src/info/bf7c7f30031888f4e796e429ab3978879485813aaca6f641c7b33e4e09459bcc) |
| official archive | [`sqlite-amalgamation-3530400.zip`](https://sqlite.org/2026/sqlite-amalgamation-3530400.zip), 2,946,650 bytes |
| archive SHA3-256 published by SQLite | `628a44cfe82c66aed1ccbbe85a562d2e33ebe64b3288981ed76285612227934e` |
| archive SHA-256 independently computed | `1e71ddf93849c6a6ecf58b827c0692073d2dd7ee40196158068f7b29f422e87d` |
| `sqlite3.c` | 9,515,341 bytes; SHA-256 `b1dd5d74ec7f29055a6684fa06fb3c2f6821c87dd38f9a458dfd2e8a1db28189` |
| `sqlite3.c` SHA3-256 published by SQLite and independently reproduced | `67f423e9ebbbdc473cbc4772c872ee6b89f31fde4ed0279a5c25d5f65c043a16` |

The [official download page](https://sqlite.org/download.html) publishes the release archive and
SHA3 identities. The experiment must fetch the archive once, verify size and both archive hashes,
extract only `sqlite3.c` and `sqlite3.h`, verify the `sqlite3.c` hashes, and then build offline.
`shell.c`, the SQLite CLI, `sqlite3ext.h`, dynamically loaded extensions, and a system SQLite are
not experiment inputs.

### Maintained binding comparison

All metadata and source archives below were inspected directly on 2026-08-05. Archive SHA-256
values are research-record integrity identifiers, not upstream signatures or attestations.

| Candidate | Exact inspected source and footprint | Disposition for F6 |
| --- | --- | --- |
| Capsule-owned cgo shim over official amalgamation | SQLite 3.53.4 identities above; new shim must be a few fixed wrappers with no Go module dependency | **Research recommendation.** Native SQLite/Darwin behavior, exact official C bytes, fixed API, and smallest auditable authority surface. Not selected or admitted. |
| [`mattn/go-sqlite3` v1.14.49](https://github.com/mattn/go-sqlite3/tree/cc41b8c87686036ea632cede537ffccef69b370a) | tag commit `cc41b8c87686036ea632cede537ffccef69b370a`; module zip SHA-256 `54cc6e644780ed238d1d39456f5bf238c79159a345356e6837d31884cf018ce1`; 103 files, 10,688,175 extracted bytes; no module dependencies; MIT plus embedded SQLite public-domain source; embeds source ID 3.53.4 | Maintained cgo comparator, but not the first F6 binding. Its broad `database/sql`, DSN/URI, build-tag, callback, extension, and optional feature surface is unnecessary. Its `sqlite3-binding.c` carries the same source ID but is not byte-identical to the official `sqlite3.c`, so the experiment would lose the simplest upstream-byte pin. |
| [`modernc.org/sqlite` v1.56.0](https://gitlab.com/cznic/sqlite/-/commit/cc920f9b5d957059ff73ff5600f8c80ef440d3af) | origin commit `cc920f9b5d957059ff73ff5600f8c80ef440d3af`; module zip SHA-256 `6321abf675769662cda07d8cd922eac3eb42e37877a0da3c1b5d7227b2c2417b`; 1,858 files, 152,999,400 extracted bytes; Go 1.25; BSD-3-Clause translation plus SQLite public-domain source; direct and indirect module graph listed below | Maintained pure-Go comparator, but not the native F6 baseline. It carries a project-patched SQLite 3.53.3 rather than primary 3.53.4, a generated translation and large libc/runtime graph, broad `database/sql`/virtual-table surface, and does not exercise the exact upstream C plus Darwin VFS path. |
| [`ncruces/go-sqlite3` v0.35.3](https://github.com/ncruces/go-sqlite3/tree/1389045d3c5c68e1805797b81aa6fc7e8dec04bd) | commit `1389045d3c5c68e1805797b81aa6fc7e8dec04bd`; module zip SHA-256 `e0df1593295a939af7786e444f6b2c91bb5f4dbe9b79428a7b54198482f53c3c`; 276 files, 1,060,565 extracted bytes; Go 1.25; MIT; Wasm engine module v3.2.35304 zip SHA-256 `22908cebfaed085750fa154be1e0596933868ed8954d2a0cf42407cc4f35b97b` | Maintained cgo-free/Wasm comparator with a useful independent implementation path, but not the first native APFS/VFS baseline. It adds a Wasm runtime/VFS boundary and broad C/`database/sql` surface. It may be used later as a test-only semantic comparator after separate graph review. |

The inspected `modernc.org/sqlite` direct graph is `google/pprof` at pseudo-version ending
`ef3492d7dac3`, `golang.org/x/sys v0.47.0`, `modernc.org/fileutil v1.4.0`,
`modernc.org/libc v1.74.4`, and `modernc.org/mathutil v1.7.1`; its indirect graph includes
`go-humanize v1.0.1`, `google/uuid v1.6.0`, `mattn/go-isatty v0.0.24`,
`ncruces/go-strftime v1.0.0`, `remyoudompheng/bigfft` at pseudo-version ending
`24d4a6f8daec`, and `modernc.org/memory v1.11.0`.
The inspected `ncruces/go-sqlite3` runtime graph includes `go-sqlite3-wasm/v3 v3.2.35304`,
`julianday v1.0.0`, `sort v1.0.0`, `wbt v1.0.0`, and `x/sys v0.47.0`; optional packages and tests
add `siphash v1.2.3`, `google/uuid v1.6.0`, `httpreadat v0.1.0`, `x/crypto v0.54.0`,
`x/sync v0.22.0`, `x/text v0.40.0`, and `adiantum v1.1.1` and are not silently admitted. The Wasm
module's tool graph is `wasm2go v0.4.11` plus `x/tools v0.45.0`; it is build-only but must remain in
the source/SBOM inventory.

### Interfaces deliberately excluded

- `database/sql` is ineligible for the first comparator because its generic driver contract,
  connection pool, implicit connection replacement, and driver-specific retry behavior obscure the
  single-owner/first-error evidence. No hidden pool or retry is allowed.
- DSN or URI strings are ineligible because caller-selected query parameters can change open mode,
  cache sharing, mutability, VFS, locking, and journal behavior. The binding accepts typed fixed
  flags and an internally derived descriptor-relative path only.
- An ORM is ineligible because generated SQL, migration/reflection hooks, callbacks, implicit
  transactions, and engine abstractions prevent an exact fixed SQL and fault boundary.
- Dynamic linkage to `/usr/lib/libsqlite3.dylib`, `-lsqlite3`, or another system SQLite is
  ineligible because OS updates can change source ID, compile options, and behavior independently
  of the evidence pin.

## Experiment build and dependency-policy closure

This closes the dependency checklist only far enough to authorize planning. An `unknown` below is
a required experiment output and keeps product admission `BLOCKED`.

| Checklist item | F6 planning disposition |
| --- | --- |
| capability / reuse row | ADR-0031 F6 production-engine comparison; reuse-map DB-1; `SPIKE-FIRST` |
| version, source, hash | exact SQLite 3.53.4 archive, source ID, Fossil check-in, sizes, SHA3-256, and SHA-256 above |
| source/build manifest | retain original archive, extracted-file manifest, Capsule shim sources, exact compile/link commands, environment, all output hashes, and a CycloneDX/SPDX SBOM in the experiment |
| compile flags | freeze the candidate flags below; record `sqlite3_compileoption_get/used` and fail any missing or unexpected option |
| C/native/SDK graph | expected inputs are `sqlite3.c`, `sqlite3.h`, the Capsule C shim, Go cgo, Apple clang/SDK, and `libSystem`; exact Go/Xcode/clang/SDK versions and `otool -L` closure are experiment outputs and currently `unknown` |
| license / notices / SBOM | SQLite source is [public domain](https://sqlite.org/copyright.html); the Capsule shim is Apache-2.0. Retain SQLite provenance/public-domain text and an SBOM. No copyleft source-offer duty is expected, but legal/license review remains an admission gate. |
| provenance | SQLite publishes archive/source hashes and Fossil identity, but no signature or SLSA attestation was verified. Record this limitation; independently reproduce hashes from two clean fetch/extract directories. |
| trust and authority | experiment/test-only during F6; candidate product persistence TCB only after selection. It gains filesystem and all stored Supervisor authority inside the disposable experiment root; no key, network, update, guest/backend, daemon, or Broker authority. |
| reproducibility | build twice from the verified cached archive in clean directories with network disabled; compare source manifests, compile-option readback, dynamic graph, and output hashes. Explain any non-byte-equal binary fields. |
| faults and obligations | run the complete corpus below, upstream SQLite result-code/VFS faults, sanitizers where supported, upstream `make test` or documented equivalent for the exact amalgamation/configuration, and any upstream bug reproducer relevant to the release |
| maintenance | Capsule storage owner monitors SQLite release chronology, security guidance, and correctness notes. Any corruption/durability fix or version/SDK change triggers the entire F6 corpus; response SLA and named human owner remain `unknown` and block admission. |
| upgrade / removal | one pinned release at a time; retain old and new compatibility fixtures, backup/restore and migration mutations, and a rollback plan. If neither F6 comparator passes, remove the experiment-only source/shim and retain the fixed store as oracle; no product rollback is implied. |

Candidate compile definitions:

```text
-DSQLITE_THREADSAFE=1
-DSQLITE_DQS=0
-DSQLITE_DEFAULT_FOREIGN_KEYS=1
-DSQLITE_ENABLE_API_ARMOR
-DSQLITE_OMIT_LOAD_EXTENSION
-DSQLITE_OMIT_SHARED_CACHE
-DSQLITE_OMIT_DEPRECATED
-DSQLITE_OMIT_JSON
-DSQLITE_MAX_MMAP_SIZE=0
```

No FTS, RTree, session, preupdate, DBSTAT, math, ICU, JSON, extension loading, shared-cache, or
memory-mapping feature may be enabled incidentally. The experiment must record SQLite defaults and
all compile options rather than inferring omission. Toolchain, warning, hardening, architecture,
deployment-target, optimization, debug, and sanitizer flags belong in the retained build manifest;
they are not selected by this documentation checkpoint.

## Fixed storage shape and publication protocol

Both candidates preserve two storage classes:

1. one mutable active SQLite database containing hot records, global replay/non-reuse tombstones,
   segment descriptors, generations, checkpoint head, durable-time high water, and exact counts;
2. zero or more closed immutable SQLite segment databases whose finalized bytes are named by
   content digest and contain complete closed cohorts plus segment-derived indexes.

A monolithic database is not an eligible F6 candidate. It would erase the F3 publish-before-
reference boundary and orphan semantics instead of replaying them.

Archive activation remains one Supervisor-owned operation under the existing installation owner
lock. Ordinary work has exactly one live SQLite connection, no pool, and no background connection.
The owner sequentially closes the active database before constructing and fully verifying a new
segment, hashes and publishes the closed segment, syncs its directory, then reopens the active
database and commits the exact digest/length/ordinal/checkpoint reference plus hot-record removal in
one transaction. It finally closes and independently reopens the whole world.

A valid segment published before the active reference is a known unreferenced orphan, not
authority and not evidence of activation. Missing/corrupt referenced segments and hot/archive
overlap are corruption. F6 reports all orphans read-only and never deletes them; F5's explicit
known-orphan deletion stays an oracle behavior and is not silently carried into the comparator.

The backup API is the only bounded connection-count exception: `sqlite3_backup` necessarily owns
one source and one newly created destination handle while the owner lock is held. Both handles are
created by the fixed wrapper, no other connection exists, the destination is fully verified after
both close, and neither handle escapes. This is not a pool.

## Fixed binding and SQL surface

The Capsule shim exports only typed functions for initialize/open/close, schema create/migrate,
full verify, exact retained lookups, fixed state mutations, segment prepare/verify/publish/activate,
backup, restore-admission verification, checkpoint, and diagnostics required by this packet. It
does not export generic exec/query, raw SQL, PRAGMA, table/column name, filename, VFS, open-flag,
or extension interfaces.

Every statement is an internal constant prepared with `sqlite3_prepare_v3` and
`SQLITE_PREPARE_NO_VTAB`; the wrapper rejects a nonempty statement tail, binds exact types and
sizes, accepts only the intended row/column shape, defensively copies result bytes, and always
resets/finalizes on every result. Schema names and migrations are fixed. Caller bytes never become
identifiers or SQL. `sqlite3_busy_timeout(db, 0)` is mandatory; no busy handler,
`unlock_notify`, sleep, retry loop, or automatic replay exists. The first `BUSY`, `LOCKED`, or
extended error is returned to the Supervisor for fail-closed reopen/classification.

### Opens, defensive configuration, and limits

All paths are derived beneath the already opened owned experiment root and validated before the
SQLite call. URI parsing is disabled. Active/segment creation uses
`READWRITE|CREATE|FULLMUTEX|PRIVATECACHE|NOFOLLOW|EXRESCODE`; reopen removes `CREATE`; immutable
segment verification uses `READONLY|FULLMUTEX|PRIVATECACHE|NOFOLLOW|EXRESCODE`. The VFS name is the
fixed built-in `unix` VFS and is read back with `SQLITE_FCNTL_VFSNAME`. A new destination must not
preexist; symlinks and unexpected file types fail before SQLite receives the path.

Immediately after every open and before schema SQL, set and read back with
`sqlite3_db_config`: `DEFENSIVE=1`, `TRUSTED_SCHEMA=0`, `ENABLE_FKEY=1`, `ENABLE_TRIGGER=0`,
`ENABLE_VIEW=0`, `ENABLE_LOAD_EXTENSION=0`, `DQS_DDL=0`, `DQS_DML=0`, `WRITABLE_SCHEMA=0`,
`ENABLE_ATTACH_CREATE=0`, `ENABLE_ATTACH_WRITE=0`, and `ENABLE_COMMENTS=0`. Unsupported or
unconfirmed controls fail the experiment rather than weakening the profile. F6-B also sets and
reads back `NO_CKPT_ON_CLOSE=1`.

Set and read back these `sqlite3_limit` values; a fixed schema/statement that exceeds one is an
experiment failure, not permission to raise it silently:

| limit | value |
| --- | ---: |
| `LENGTH` | 1,048,576 |
| `SQL_LENGTH` | 16,384 |
| `COLUMN` | 64 |
| `EXPR_DEPTH` | 32 |
| `COMPOUND_SELECT` | 4 |
| `VDBE_OP` | 10,000 |
| `FUNCTION_ARG` | 16 |
| `ATTACHED` | 0 |
| `LIKE_PATTERN_LENGTH` | 0 |
| `VARIABLE_NUMBER` | 64 |
| `TRIGGER_DEPTH` | 0 |
| `WORKER_THREADS` | 0 |
| `PARSER_DEPTH` | 100 |

Record `sqlite3_libversion`, `sqlite3_sourceid`, `sqlite3_threadsafe`, compile options, database
read-only/autocommit state, extended result codes, VFS, schema/user/application versions, page
size/count/freelist, foreign-key state/check, journal/locking/sync settings, checkpoint state,
`quick_check`, and full `integrity_check`. Required readback disagreement is a closed failure.

## F6-A and F6-B comparator configurations

Both comparators use the identical amalgamation, shim, schema, fixed SQL, corpus, page size,
ownership, paths, limits, and verification. Only journal-specific configuration differs.

| Setting | F6-A rollback journal | F6-B WAL |
| --- | --- | --- |
| `page_size` before schema | `4096` | `4096` |
| `journal_mode` | `DELETE` | `WAL` |
| `synchronous` | `FULL` | `FULL` |
| `locking_mode` | `NORMAL` | `NORMAL` |
| `fullfsync` | `ON` | `ON` |
| `checkpoint_fullfsync` | record/read back; not used | `ON` |
| `wal_autocheckpoint` | not applicable | `0` |
| close checkpoint | not applicable | disabled with `NO_CKPT_ON_CLOSE=1` |
| checkpoint | not applicable | explicit `sqlite3_wal_checkpoint_v2` only at named corpus boundaries; record mode/result/log/frame counts |
| `auto_vacuum` | `NONE` | `NONE` |
| `journal_size_limit` | `0` | `0` |
| `mmap_size` | `0` | `0` |

No comparator may select `OFF`, `MEMORY`, `NORMAL`, or asynchronous durability; use a hidden WAL
checkpoint; invoke `VACUUM` to conceal growth; or vary a setting without starting a new recorded
candidate identity. Record active database, journal, WAL, SHM, segment, and directory inventory at
every stateful boundary.

## Required corpus and evidence

The fixed store is the semantic oracle. F6 must reproduce its exact known answers, mutations,
fault points, file inventories, digests, generations, counts, lookup results, recovery set,
tombstones, and zero-adapter-call assertions; SQLite success is not defined by SQL tests alone.

| Prior slice | F6 replay obligation |
| --- | --- |
| F2 format/migration | v1 attempt-without-lifecycle union; empty-archive v1-to-v2 migration; exact genesis; version/missing/downgrade refusal; full index/count reconstruction; every corrupt/mismatched field and pre/post-publication fault |
| F3 first activation | complete cohort selection; segment prepare/full verify/publish-before-reference; old-or-new reopen after every fault/process death; known orphan; missing/corrupt/substituted segment; no partial cohort or hot/archive overlap |
| F4A lookup/replay | identical hot/archive registration, approval, attempt, lifecycle, nonce, payload, effect, and instance answers; original replay bytes/state; collisions; corrupt/missing location refusal; defensive copies; cancellation/concurrency; archived attempts excluded from recovery |
| F4B mutation | one atomic authority/lifecycle mutation; independent append-only effect tombstone; current-versus-historical separation; response loss, restart, retry, digest/identifier collision, and every transaction fault with no duplicate external effect |
| F4C growth | second/later segment activation; 64/65 boundary; first-segment bytes remain unchanged; exact caps/refusal; repeated fault/reopen and no referenced eviction |
| F5 backup/restore admission | owner-held manifest-last backup; exact active/segment inventory; backup API copy; complete independent reopen; exact-checkpoint read-only restore admission and rollback-uncertain refusal; partial/extra/corrupt/missing files; offline deterministic report; known/unknown orphan classification without F6 deletion or activation |

Add SQLite-native tests for primary and extended `BUSY`, `LOCKED`, `FULL`, `IOERR`, `CORRUPT`,
`NOTADB`, `READONLY`, and `CANTOPEN` results. A controlled test VFS must inject bounded failures in
`xRead`, `xWrite`, `xSync`, `xTruncate`, `xLock`, `xUnlock`, and relevant `xFileControl` calls at
every recorded ordinal. Mutate database headers/pages/freelists/indexes, rollback journals, WAL
headers/frames/checksums, and SHM state; preserve the exact original and mutation manifest. No
silent repair, fallback, journal deletion, rebuild, or retry counts as success.

Run real same-host multi-process contenders against the owned disposable root and prove
nonblocking first-error behavior plus no bytes changed by the loser. Separately run SIGKILL at each
commit, segment publication, directory sync, backup, WAL frame, checkpoint, and close boundary.
The final gate uses a separately authorized owned disposable Apple-silicon Mac and dedicated APFS
volume for actual power interruption at the same boundaries. After every restart, capture the
filesystem inventory and hashes before opening, then accept exactly the prior or committed world or
fail closed as repair-required. Never run these tests against user, product, or unrelated data.

Quantitative results must report distributions and worst cases for open/full verification,
ordinary mutation, F4B response-loss retry, hot/archive lookup, activation, checkpoint, backup,
restore admission, and offline verification; active/segment/journal/WAL/SHM bytes; memory/CPU; and
lock wait (which must remain zero in the no-retry profile). The experiment may recommend budgets
but does not select a product latency source, window, lifetime, persistence model, or threshold
action.

## ADR-0040 threshold reconciliation

ADR-0040's owner-only disposable exception retains exact stop thresholds: 128 attempts, 8 MiB
active store, 16 referenced segments, plus its timing thresholds. The merged passive checker is a
read-only, point-in-time admission oracle. Merged
[PR #227](https://github.com/Shrimpworks/capsule-corp/pull/227), based on reviewed branch commit
`4fe48f5e407be6253f8a8f0c0a9b2789830391b8`, clarifies that it has no persistent trip latch, no
selected p95 source/window/lifetime/persistence, no `RequestAttempt` wiring, and no response-loss
integration. That merged canonical state is an F6 input, not evidence that enforcement exists.

F6 must test and measure the same numerical boundaries, but passing them does not make their full
operational enforcement exist. Wiring every admission and mutation, persisting a trip state,
selecting the timing observation policy, integrating response loss, and proving retirement remains
`BLOCKED`. The experiment must not claim that a successful passive check enforces the five
thresholds operationally.

## Exit, retained evidence, and prohibitions

The separately authorized experiment must retain an exact source/build/SBOM manifest, both clean
builds, fixed schema/SQL identity, configuration/readback, complete case manifest, per-fault file
inventories and hashes, process/power-loss logs, quantitative results, and limitations at an exact
`capsule-experiments` commit. Its handoff reports F6-A and F6-B separately; a failed candidate is
`NO_GO` only when the exact tested path is abandoned, not merely because evidence is incomplete.

F6 cannot pass until both the required local/VFS corpus and the authorized real APFS/power-loss gate
complete and a decision ADR records the evidence. Even then, product persistence, restore
activation, deletion, anti-rollback, continuous service, consumer wiring, and admission remain
separate gates. The fixed store remains the semantic oracle throughout; there is no simultaneous
cutover.
