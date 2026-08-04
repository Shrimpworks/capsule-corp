# Ecosystem reuse and adoption map

Status: canonical planning input; research current through 2026-08-03. This document selects
planning direction, not product dependencies or security posture. No candidate below is admitted
to production unless its consuming roadmap slice closes the stated acceptance evidence.

## Executive summary

Capsule should reuse four kinds of mature work aggressively, while keeping authority-specific
policy narrow:

1. **Adopt macOS enforcement primitives.** `SMAppService`, XPC/Mach services,
   `SecCodeCreateWithXPCMessage`, Code Signing Services, Keychain access controls,
   LocalAuthentication, App Sandbox, Hardened Runtime, notarization, and unified logging already
   own platform lifecycle, identity, key, user-presence, containment, and diagnostic mechanics.
   Capsule still owns the enrolled identities, exact requirements, epochs, typed messages, and
   deny-by-default policy around them.
2. **Adopt standards and small implementations, not generic envelopes.** RFC 8949 deterministic
   CBOR, RFC 9052 COSE, RFC 8610 CDDL, TUF, in-toto/SLSA, SPDX/CycloneDX, and Test262 should be
   inputs or oracles. The bounded production comparison selects pinned `fxamacker/cbor` only for
   typed object encoding/decoding behind Capsule-owned profiles; it rejects `go-cose` as a product
   envelope dependency. Handwritten scanners, policy wrappers, and both candidates remain
   conformance oracles as described below.
3. **Spike transactional storage before expanding the fixed snapshot.** The fixed JSON snapshot is
   the correct F2 conformance oracle, not a production database. A single bounded comparison
   against SQLite must reuse the existing logical, fault, corruption, lock, archive, backup, and
   restoration corpus before a production engine is selected.
4. **Keep governed forks exceptional.** Governed `deno_core`/`rusty_v8` and libkrun are justified
   only where physical omission or a narrow upstream gap is part of the security property. Shared
   CI, SBOM, provenance, source, patch, and release-manifest generation should become a versioned
   build-only tool or generated contract, never a cross-repository runtime dependency.

The highest-risk non-adoptions are equally important: no full Deno/Bun runtime, no generic XPC
`Codable`/JSON command bus, no network TUF/DID resolution on live approval or execution paths, no
rich parser in the daemon/Supervisor, no SQLCipher/ORM by default, no OCI registry/image framework
for the v0 raw-root path, and no hand-written ECMAScript lexer. These choices prevent convenience
frameworks from silently acquiring network, parsing, storage, process, filesystem, update, or
guest authority.

## How to use this map

Every task that proposes a dependency or a new custom primitive must identify the matching row,
record the recommendation and consuming slice, and complete the dependency-policy checklist below.
If no row applies, update this document before implementation. `ADOPT` does not mean “add now”; it
means the named roadmap slice may adopt the candidate after its acceptance criteria pass.

Closed recommendation vocabulary:

- **ADOPT-PLATFORM** — use an OS/framework primitive and keep Capsule policy around it.
- **ADOPT-PINNED** — use an exact version/digest after dependency admission.
- **GOVERN/FORK** — retain a minimal reviewed fork when upstream cannot express a structural
  property.
- **SPIKE-FIRST** — run the exact bounded experiment before selection.
- **TEST-ONLY** — oracle, fixture, fuzz, or evidence dependency; never linked into product paths.
- **DEFER** — wait for its named prerequisite or later roadmap boundary.
- **BUILD-NARROWLY** — Capsule-specific policy or fixed framing is smaller and safer than adoption.

Trust classes are `product TCB`, `planning/approval-understanding TCB`, `build-only`, `test-only`,
and `evidence-only`. “Footprint” counts direct/transitive packages where the official metadata or
current lock is available; otherwise it says `unknown` instead of guessing.

## Current dependency and scaffold baseline

- Product Go has one direct module, `golang.org/x/sys v0.28.0`; the TypeScript workspace has five
  root development dependencies and only internal workspace runtime links.
- Current production-shaped CBOR parsing/encoding remains handwritten and bounded. The retained
  Gate A2 experiment used `fxamacker/cbor` 2.9.1, `cbor2` 2.3.0 plus one transitive package, and
  `thecoolwinter/CBOR` 1.1.2 only as spike inputs. The later production-profile comparison pins
  `fxamacker/cbor` 2.9.2 and `go-cose` 1.3.0 in a standalone experiment; neither is in the root
  product module.
- `FixedFileStore` rewrites a versioned, bounded JSON snapshot with file and directory durability
  barriers. It explicitly is not a production database. Archive F1/F2 and owner-lock G1 are
  passive/unwired mechanics.
- The native XPC/Security, libkrun, Deno, `rusty_v8`, Swift, and C code lives in retained
  experiments or governed forks, not admitted product paths.

## Candidate register

| ID | Candidate and exact primary source | Maturity, maintenance, license, provenance | Platforms / language / footprint | Security and TCB disposition |
| --- | --- | --- | --- | --- |
| APL-1 | Apple [`SMAppService`](https://developer.apple.com/documentation/servicemanagement/smappservice) | Shipped macOS framework, macOS 13+; proprietary OS component, Apple-signed distribution | macOS native API; no package dependency | Product TCB platform primitive; grants service registration/lifecycle, not Capsule approval or guest authority |
| APL-2 | Apple [XPC](https://developer.apple.com/documentation/xpc) and [`SecCodeCreateWithXPCMessage`](https://developer.apple.com/documentation/security/seccodecreatewithxpcmessage%28_%3A_%3A_%3A%29) | Shipped OS IPC and Security APIs; Apple lifecycle/advisories | macOS C/Objective-C boundary; no package dependency | Product TCB; observes live sender and transports messages, while Capsule owns exact peer requirements, caps, roles, and replay semantics |
| APL-3 | Apple [Code Signing Services](https://developer.apple.com/documentation/security/code-signing-services), [Keychain Services](https://developer.apple.com/documentation/security/keychain-services), and [LocalAuthentication](https://developer.apple.com/documentation/localauthentication) | Shipped OS security APIs; proprietary/Apple-signed | macOS native; no package dependency | Product and approval-understanding TCB; gains code-identity, key-use, secret-store, and user-presence authority only in the owning component |
| APL-4 | Apple [App Sandbox](https://developer.apple.com/documentation/security/app-sandbox), [Hardened Runtime](https://developer.apple.com/documentation/security/hardened-runtime), and [notarization](https://developer.apple.com/documentation/security/notarizing-macos-software-before-distribution) | Supported Apple distribution controls | macOS build/install; entitlement and signing identities required | Product/platform and build TCB; constrains process/filesystem/device authority but does not replace runtime/backend controls |
| APL-5 | Apple [Unified Logging](https://developer.apple.com/documentation/os/logging) and [signposts](https://developer.apple.com/documentation/os/os_signpost) | Supported OS observability APIs | macOS native; no package dependency | Product TCB only through bounded adapters; diagnostic authority, with privacy marking and no user content/keys/source |
| APL-6 | Apple [CryptoKit](https://developer.apple.com/documentation/cryptokit) and CommonCrypto as exposed by the installed SDK | Shipped OS cryptography; Apple recommends CryptoKit for high-level Swift use; proprietary/Apple-signed | macOS Swift/C; no package dependency | Product cryptographic TCB; digest/signature operations only, with Capsule retaining serialization, key authorization, purpose, and replay policy |
| STD-1 | IETF [RFC 8949 CBOR](https://www.rfc-editor.org/rfc/rfc8949), [RFC 9052 COSE](https://www.rfc-editor.org/rfc/rfc9052), and [RFC 8610 CDDL](https://www.rfc-editor.org/rfc/rfc8610) | Internet Standards / Standards Track; public errata process; specification licenses | Cross-language standards; zero runtime footprint by themselves | Planning/approval-understanding TCB specifications; no authority until implemented |
| GO-1 | [`fxamacker/cbor` v2.9.2](https://github.com/fxamacker/cbor/releases/tag/v2.9.2) | Active; 2026-05-04 release; MIT; latest-release security owner; exact tag, module sums, zip hash, and license hash retained; mutable GitHub release has no asset/attestation | Go 1.20+; adds `x448/float16` 0.8.4; 1,204 + 76 KiB extracted source and +153,104-byte stripped arm64 comparison delta | Selected future product parsing TCB only for typed encode/decode behind Capsule predecode and object wrappers; no network/filesystem/process authority |
| GO-2 | [`veraison/go-cose` v1.3.0](https://github.com/veraison/go-cose/releases/tag/v1.3.0) | 2024-09-27 release; MPL-2.0; exact tag/sums/zip/license hashes retained; mutable release has no asset/attestation; SECURITY page still lists v1.0.0 as the only supported stable version | Go 1.21+; MVS selects GO-1; 1,916 KiB/68 source files and +150,112 bytes over the fxamacker stripped comparison | **TEST-ONLY / production NO-GO.** Sig_structure and ES256 agree, but the generic envelope/parser surface does not own caps, canonical bytes, key authorization, object binding, or replay |
| GO-3 | Go standard [`crypto/sha256`](https://pkg.go.dev/crypto/sha256), [`crypto/ecdsa`](https://pkg.go.dev/crypto/ecdsa), [`crypto/rand`](https://pkg.go.dev/crypto/rand), and [`crypto/subtle`](https://pkg.go.dev/crypto/subtle) | Go toolchain-maintained, BSD-3-Clause distribution | Cross-platform Go; zero external modules | Product TCB; cryptographic primitives only. Capsule owns purpose, key authorization, serialization, replay, and state |
| TUF-1 | [`go-tuf/v2`](https://github.com/theupdateframework/go-tuf) and the [TUF specification](https://theupdateframework.github.io/specification/latest/) | CNCF project, active maintainers, Apache-2.0, SECURITY/CODEOWNERS; exact future version/digest and graph unknown | Go, portable; network-capable client surface must be isolated from live paths | Updater/build product TCB only; gains local trust-metadata parsing and optional fetch authority, never approval/execution network authority |
| SUP-1 | [in-toto Attestation Framework](https://github.com/in-toto/attestation), [SLSA provenance](https://slsa.dev/spec/v1.2/provenance), [SPDX 3.0](https://spdx.github.io/spdx-spec/v3.0/), and [CycloneDX 1.7](https://github.com/CycloneDX/specification/releases/tag/1.7) | Open standards; in-toto/SLSA under community governance; SPDX ISO/IEC 5962; CycloneDX ECMA-424 and schemas Apache-2.0 | Generated JSON/CBOR documents; generators have separate graphs to pin | Build/evidence-only for v0; records materials and components but never proves runtime admission by itself |
| DB-1 | [SQLite 3.53.3](https://sqlite.org/releaselog/3_53_3.html), [atomic commit](https://sqlite.org/atomiccommit.html), [backup API](https://sqlite.org/backup.html), and [`integrity_check`](https://sqlite.org/pragma.html#pragma_integrity_check) | Mature, actively maintained, public domain core; 2026-06-26 source ID and SHA3 published; recent WAL-reset corruption fix demonstrates active security/correctness maintenance | C library on macOS; a Go C-ABI or pure-Go binding adds a separately measured graph | Candidate product persistence TCB; gains all Supervisor durable authority in one database and therefore requires the strictest spike and restoration evidence |
| JS-1 | [ECMA-262](https://tc39.es/ecma262/) and [`test262`](https://github.com/tc39/test262) | TC39 normative spec and official conformance suite; BSD-3-Clause test suite; rolling corpus rather than stable product library | Test-only JS fixtures/harness; corpus size is large and selected commit must be pinned | Test-only; exercises parser/runtime behavior without becoming an approval or runtime dependency |
| RT-1 | Governed Deno/`deno_core` and `rusty_v8` exact refs recorded in [the work plan](GOVERNED_DENO_CORE_WORK_PLAN.md) | Upstream MIT plus third-party notices; Capsule fork governance and exact commits required; ARM64 handoff in progress | Rust/C++/V8; very large transitive/source/build footprint | Govern/fork product runtime TCB only after admission; would gain hostile-source execution inside external isolation, never host authority |
| BK-1 | Upstream [`libkrun`](https://github.com/libkrun/libkrun) plus Capsule-governed fork identity in [the evidence ledger](WORKSTREAM_EVIDENCE_LEDGER.md) | Apache-2.0 upstream; exact v1.19.4-derived base/patch/merge identities retained; advisories and final release provenance remain open | Rust/C/macOS HVF, large VMM/guest-facing TCB | Govern/fork backend TCB; gains guest/VMM/device/FD authority only inside the Supervisor-owned runner boundary |
| BK-2 | `imago` 0.2.3 as pinned by the retained libkrun supply-chain/root-custody evidence | Exact crate version is retained; maintenance, license, advisory owner, release signing, and standalone transitive graph are **unknown** pending final bundle admission | Rust storage dependency inside the governed libkrun graph | Product backend TCB; positional block I/O over the transferred root descriptor, with no authority to choose a path or format |
| KRN-1 | Linux kernel [seccomp](https://www.kernel.org/doc/html/latest/userspace-api/seccomp_filter.html), [`no_new_privs`](https://www.kernel.org/doc/html/latest/userspace-api/no_new_privs.html), and [Landlock](https://www.kernel.org/doc/html/latest/userspace-api/landlock.html) | Stable kernel ABIs with upstream security/advisory process; GPL-2.0 kernel implementation | Linux guest/worker only; no package dependency but exact kernel/config becomes TCB | Runtime/worker defense-in-depth; process/syscall/filesystem restriction, never a replacement for the external VM/gVisor boundary |
| TEST-1 | Go [fuzzing](https://go.dev/doc/security/fuzz/), [race detector](https://go.dev/doc/articles/race_detector), and [`govulncheck`](https://go.dev/doc/tutorial/govulncheck) | Toolchain/Go security team maintained; Go/BSD licensing; vulnerability DB update process | Build/test only; `govulncheck` currently fetches a database unless an offline mirror is pinned | Test/evidence-only; local code execution and optional advisory-network authority in CI, never product authority |
| RUST-1 | Rust [cargo-fuzz](https://rust-fuzz.github.io/book/cargo-fuzz.html), [Miri](https://github.com/rust-lang/miri), [sanitizers](https://doc.rust-lang.org/beta/unstable-book/compiler-flags/sanitizer.html), and [LLVM coverage](https://doc.rust-lang.org/rustc/instrument-coverage.html) | Rust/LLVM project tooling; exact toolchain pins required; licenses vary by toolchain component | Build/test only, large toolchain footprint | Test/evidence-only; executes repository fixtures on owned CI/hosts, no product authority |

Apple APIs are versioned with the OS rather than package releases. “No package dependency” does
not mean no TCB: the corresponding OS framework and kernel/service implementation remain trusted.

## Roadmap adoption ledger

| Roadmap capability / current mechanism | Candidate | Trust placement and authority consequence | Recommendation and rationale | Consuming task, dependencies, and acceptance evidence |
| --- | --- | --- | --- | --- |
| Phase 1 / installed per-user services; current LaunchAgent/descriptor experiments | APL-1 | Product platform TCB; registers only embedded per-user services | **ADOPT-PLATFORM.** Do not build a launchd installer/controller. Keep registration state distinct from enrollment or approval. | P0-4 installed topology and Phase 8 distribution; exact minimum macOS, user approval, crash/relaunch, update, uninstall, wrong-session, clean-host, signing/notarization readback |
| Phase 1/2 authenticated IPC; current C probes and proposed four-call XPC surface | APL-2 | Product TCB; sender observation and local transport only | **ADOPT-PLATFORM.** Use separate role services and closed XPC dictionaries; do not adopt generic NSXPC/Codable/JSON routing. | ADR-0029 S1 after M1 contract availability; exact audit-token sender, EUID/session, dynamic/static validity, cap/type/key-count, interruption/replay, C-to-Go ownership corpus |
| Installation identity and exact-build checks | APL-2 + APL-3 | Product TCB; code identity and trust-epoch observations | **ADOPT-PLATFORM.** `SecCode` supplies evidence, not policy. Build identifiers, epoch, and allowed requirements remain Supervisor data. | IPC S1 and Gate B installed evidence; wrong-role/team/build, copied/replaced/debugged peer, mixed epoch, process-death cases |
| Approval and operational private keys; experiment-only Keychain/Secure Enclave | APL-3 + GO-3 | Approval/product TCB; key-use authority remains role-separated | **ADOPT-PLATFORM.** Do not add a portable keystore or software-key fallback to product. Crypto APIs do not authorize a key by themselves. | Phase 3 approval and Phase 8 updates; access-group/epoch rotation, user-presence, stale-build group, backup/restore, Secure Enclave availability, no-daemon-use evidence |
| Broker user-presence UI; no product UI | APL-3 plus SwiftUI/AppKit accessibility APIs | Approval-understanding TCB; presents typed registered data and can request user presence | **ADOPT-PLATFORM** for presentation/authentication; **BUILD-NARROWLY** for the fixed approval view and anti-confusion rules. | Phase 3 UI-001; exact bytes-to-view mapping, accessibility, overlay/focus/synthetic-input, stale-session, cancellation, and user-presence-signing evidence |
| Signing, notarization, App Sandbox, Hardened Runtime; ad hoc experimental signatures | APL-4 | Build/platform TCB; constrains installed code and release identity | **ADOPT-PLATFORM.** No custom package signer or entitlement shim. | P0-4/Phase 8; final bundle, entitlements, staple/Gatekeeper, clean/minimum-OS hosts, mixed-update and readback corpus |
| Privacy-safe diagnostics; current ad hoc experiment logs/evidence files | APL-5 | Product diagnostic authority; log store may be read by system/admin tooling | **ADOPT-PLATFORM** through a closed event adapter. Never log source, content, guest strings, keys, raw approvals, paths, or unbounded IDs. | Phase 3 operations and every installed slice; privacy annotations, fixed event taxonomy, cap tests, redaction mutation tests, signpost IDs unrelated to security identity |
| Phase 2 deterministic CBOR predecoder/encoder; handwritten bounded Go code and cross-language fixtures | STD-1 + GO-1 | Product parsing TCB; attacker-controlled proposal/signed-object bytes, no other authority | **ADOPT-PINNED narrowly after schema freeze:** v2.9.2 may replace deterministic typed encoding and field decode only. Keep raw-byte cap/predecode, exact object limits, canonical byte comparison, closed schema/bindings, exact-byte identity, and the handwritten oracle. | [Production comparison](../experiments/production-cbor-cose-profile/RESULTS.md): 17 SourceManifest, 23 role/domain, 40 CBOR, 10,000 property cases, 2.79M fuzz inputs, restoration and resource evidence. Next: frozen-object integration and independent review before root-module admission |
| COSE_Sign1 wrappers; experiment currently handwrites envelope and Sig_structure | STD-1 + GO-2 + GO-3 | Product crypto-envelope TCB; header parsing/signature input, not key selection or authorization | **TEST-ONLY for go-cose v1.3.0; production NO-GO.** Keep a Capsule-owned Sign1-only wrapper with standard crypto. The candidate saves too little after required caps/header/binding/replay controls and adds broad generic COSE surface plus unresolved support/provenance concerns. | [Production comparison](../experiments/production-cbor-cose-profile/RESULTS.md): 82 approval cases, 2.89M fuzz inputs, exact Sig_structure/equivalent-S, restored header families, trusted-key and binding refusals. Next exact test: production-shaped Swift and same-byte Supervisor/Broker integration |
| SHA-256, randomness, P-256 verification, constant-time comparisons | GO-3 + APL-3 + APL-6 | Product crypto TCB; no parsing or policy authority | **ADOPT-PLATFORM/STDLIB.** Prefer CryptoKit in Swift and the Go standard library in Go; use CommonCrypto only at the narrow C boundary already evidenced. Never implement primitives. HKDF is **DEFER** until a protocol requires it. | All digest/signature slices; known-answer and cross-language vectors, algorithm availability, malformed key/signature cases, raw/DER boundary, FIPS/export constraints recorded if applicable |
| CDDL/schema conformance; current custom Node verification scripts | STD-1 plus current generators | Build/test TCB only | **BUILD-NARROWLY** around the official CDDL grammar. Existing scripts encode Capsule caps/field authority; do not replace them with a permissive generic code generator. A CDDL validator may be **TEST-ONLY** after pin/footprint review. | Phase 2 schema freeze; schema-to-fixture drift, unknown field, exact-max, generation idempotence, cross-language decoder agreement |
| Phase 2 fixed snapshot, F1/F2 archive oracle | DB-1 versus current `FixedFileStore` | A production DB would own all Supervisor state, locks, journals, backups, migrations, and corruption responses | **SPIKE-FIRST.** Keep fixed v2 oracle through F2. Do not extend it into an unbounded production engine, and do not select SQLite/driver by reputation alone. | After F2 format/full verifier and G2 composition: one storage-engine spike replays logical/fault/mutation corpus; exact SQLite/driver pin, compile flags, graph, single-writer mode, `FULL` sync, WAL/rollback decision, directory/APFS semantics, integrity/foreign-key checks, online/offline backup, restore/rollback, migration, disk-full, torn/corrupt pages, lock loss, process/power interruption |
| Store encryption | Keychain file-protection design; SQLCipher not selected | Product TCB; encryption layer gains full state plaintext/key authority | **DEFER / DO NOT ADOPT SQLCipher by default.** Encryption does not solve same-role compromise, rollback, corruption, or owner composition and adds a fork/key/migration surface. | Revisit only after threat model names an at-rest attacker not covered by protected storage; then an ADR and measured restore/key-loss/rekey corpus are mandatory |
| Owner lock G1 and G2; exact BSD `flock` object selected | Darwin `openat`/`fstat`/`fcntl`/`flock` through `x/sys` | Product TCB; cooperating-process ownership only, no same-UID containment | **ADOPT-PLATFORM** and keep wrapper narrow. No file-coordination framework or distributed lock library. | Active G2 owns composition; do not duplicate. Acceptance remains owner-before-store, exact inode/enrollment, lock lifetime/loss, contender no-access, migration, crash/restart, installed protected root. G3 later supplies session/update evidence |
| Archive/compaction; F2 format correction in progress | Filesystem/APFS primitives now; DB-1 only after logical oracle | Product persistence TCB; retained replay/tombstone authority | **BUILD-NARROWLY** for archive semantics and immutable segment format; **SPIKE-FIRST** for engine mechanics. Generic log compaction may delete referenced history or blur hot/archive index domains. | Active F2 owns correction and v2 verifier; do not duplicate. Later production-store spike must prove complete-cohort atomic move, immutable digest, global indexes, no referenced deletion, restore and mutation behavior |
| Backup/restore and corruption handling | DB-1 backup/integrity facilities plus Capsule epoch/checkpoint rules | Product persistence TCB; can copy or restore all durable authority | **ADOPT-PINNED** only with selected engine; **BUILD-NARROWLY** for trust-epoch and repair-required policy. | Phase 2/8 storage selection; coherent/partial restore, stale epoch, missing archive, checksum/index mismatch, destination durability, read-only diagnostics, no silent repair/deletion |
| `.mjs` import/module validation; passive V0 frames only | JS-1 plus Proposed ADR-0035's exact Oxc candidate | Planning/approval-understanding TCB when implemented; current Go/Rust codecs are build/test only and gain no process or parser authority | **SPIKE-FIRST** for Oxc artifact/process admission; **BUILD-NARROWLY** for the fixed Capsule frame. The Rust oracle pins `sha2` 0.10.9 and test-only `serde_json` 1.0.151 under lock digest `a45ad0e2b2311d33b16e46e0bf1f66c1563dd240a35f1f9fe431c7bea5894c98`; it is not a product dependency. | V0 passive contract is observed. V1-V6 still require artifact provenance, sandbox/fault evidence, independent consumers, grammar corpus, and runtime no-loader enforcement. Do not handwrite a lexer or substitute the passive `allow` observation for launch authority. |
| `.mjs` passive source/manifest fixtures | Current schema/generator infrastructure | Build/test only; no user-content authority | **BUILD-NARROWLY.** This is Capsule protocol policy, not an ecosystem primitive. | Active M1 owns source/manifest work; preserve byte-exact pass-through, one `main.mjs`, no dependency graph, exact caps, recursive field authority; this audit does not preempt it |
| Runtime ECMAScript semantics and regressions | JS-1 | Test-only; executes selected fixtures only in controlled builds/guests | **TEST-ONLY.** Pin a Test262 commit and a small applicability manifest; do not claim whole-suite conformance or import the suite into product. | Governed runtime validation after admitted build input; selected syntax/evaluation/global/module-loader cases, timeout/memory caps, expected-fail ledger and update procedure |
| Runtime implementation | RT-1 | Product hostile-code TCB inside external isolation; very large source/build closure | **GOVERN/FORK.** Reuse V8/`deno_core` internals only with physical omission and governed construction. Do not adopt full Deno/Bun or generic loaders. | Governed runtime plan; wait for separate ARM64 `rusty_v8` handoff. Exact fork refs, independent builders, source/notices, patch/removal map, SBOM/provenance, no-loader/global/ops evidence, final signed bytes and external-isolation composition |
| Runtime source maps and diagnostics | Governed V8/`deno_core` utilities | Runtime/evidence path; diagnostic output can leak source/content | **DEFER.** Use only bounded user-visible Broker diagnostics after runtime selection; agent-facing output remains fixed. | Later runtime UX slice; no ambient file lookup, bounded frames/text, source-map bytes registered/approved or omitted, adversarial location/string caps |
| Backend/VMM, virtio block/console and raw FD | BK-1 + BK-2 | Backend product TCB; guest/device/FD/process authority. `imago` receives only the finalized duplicate root descriptor | **GOVERN/FORK** for the five exact last-mile patches; reuse upstream libkrun/virtio and pinned imago positional-I/O utilities. Do not invent a second VMM/virtio stack or allow imago pathname/format selection. | Gate C P0-1..P0-4 and later owned-guest composition; exact governed head/patch order, imago license/graph/advisory closure, route/restoration mutations, sanitizers/coverage/fuzz, descriptor manifests, typed transport, teardown, installed signed bytes, clean-host profile evidence |
| Linux runtime/launcher syscall and filesystem confinement | KRN-1 | Product runtime/worker defense-in-depth; filters process/syscall/filesystem authority in one exact kernel | **ADOPT-PLATFORM** only after runtime construction. Keep structural physical omission and external isolation; do not replace them with a syscall/path denylist or a stateful broker by convenience. | Governed runtime and OCI/gVisor profile gates; exact kernel/architecture, TSYNC/thread timing, inherited descriptors, `execve`/clone/socket mutations, Landlock ABI/handled-access set, restoration tests and syscall trace |
| Fixed source/input/completion envelopes; current Go/Node 43-vector model | Standard SHA-256 plus Capsule framing | Product transport TCB; bounded attempt-specific bytes, no filesystem/network authority | **BUILD-NARROWLY.** Exact roles, caps, binding, continuous drain and last-written commit are Capsule-specific and smaller than Protobuf/gRPC/JSON frameworks. | P0-3 real transport; independent native implementation, exact/cap+1/partial/zero/stall/death/role-swap/trailer mutations through governed console, launcher and teardown composition |
| Guest immutable root and artifact snapshots | OS descriptor APIs + SHA-256; no OCI client | Broker/Supervisor content and backend TCB; immutable byte custody and guest attachment | **BUILD-NARROWLY** for v0 raw roots and CAS manifests. **DEFER** OCI tooling until a concrete OCI backend/profile requires it. | P0-1 and Phase 6; descriptor-relative open, no paths/live mounts, digest-before-use/readback, mutation race, special-file/link refusal, quota and cleanup evidence |
| File artifacts / ext4 extraction | Existing experiment tools only | Future disposable parser TCB; gains bounded artifact bytes, never daemon/Supervisor authority | **DEFER / SPIKE-FIRST.** No `debugfs`, archive library, or rich parser in product control-plane processes. | Phase 6/7 artifact gate; select a disposable parser sandbox with fixed image/profile and output slots; fuzz malformed metadata, sparse/overlap/xattr/device/link, crash/timeout and partial-release cleanup |
| TUF trust metadata and updates | TUF-1 | Updater product TCB; local metadata parsing and separately bounded fetch authority | **ADOPT-PINNED** later. No live network/parser in approval or execution; Supervisor consumes only a compact verified local snapshot. | Phase 8; exact version/digest/graph, root bootstrap/rotation, freeze/rollback/mix/delegation/expiry, offline mirror/reproducible inputs, size/depth/role caps, interruption and repair corpus |
| DIDs and JSON-LD resolution | Standards documentation only | Evidence/interoperability; a resolver would add network and rich parsing authority | **DO NOT ADOPT** a resolver for v0. Keep DIDs inert identifiers and never load remote contexts. | Deferred interoperability only after ADR; fixed method allowlist, offline documents, no authorization consequence and bounded parsing would be prerequisites |
| SBOM, source inventory, provenance, admission manifest | SUP-1 + existing evidence generators | Build/evidence-only; reads build graph/materials and emits metadata | **ADOPT-PINNED** formats; **BUILD-NARROWLY** one shared release generator/contract. Prefer CycloneDX for component/vulnerability detail and in-toto/SLSA provenance; emit SPDX when distribution/license consumers require it. | Near-term governed-fork/release governance: versioned schema, exact generator digest, offline input manifest, deterministic output, restore/mutation tests, notices, signatures, two-builder comparison. Metadata never promotes a profile automatically |
| Cross-repo CI/release checks | Generated contract or a small pinned build-only tool | Build-only; repository/source/release mutation authority in CI, never runtime | **BUILD-NARROWLY.** Centralize patch/source/SBOM/provenance schema and verification without a shared runtime library. | Repository governance slice; identical known answers across capsule-corp/Deno/rusty_v8/libkrun, exact tool pin, least-privilege CI token, no release publication in verification, backwards-compatible contract and independent reconstruction |
| Go quality/security tests | TEST-1 | Test/evidence-only; local process and optional advisory DB network authority | **ADOPT-PINNED/TEST-ONLY.** Add fuzz/race lanes and pin `govulncheck` plus its database snapshot for release evidence; ordinary developer checks may use current DB with timestamp recorded. | All Go product slices; seed corpora, max-size/fault targets, race-enabled lifecycle/store tests, crash retention, call-graph limitations, offline release replay |
| Rust/native quality/security tests | RUST-1 plus Clippy/ASan already used in governed work | Test/evidence-only; owned builders and fixtures | **TEST-ONLY.** Reuse upstream toolchain instrumentation; do not call library coverage “guest/VMM validation.” | Governed runtime/libkrun PR and release gates; exact compiler/tool pins, Miri-applicable crates, ASan/UBSan, cargo-fuzz corpora, mutation restoration, coverage mapped to exact source, arm64 and macOS lanes |
| TypeScript schemas/SDK/CLI | JSON Schema, CDDL fixtures, TypeScript compiler, current workspace | Agent-facing planning/SDK TCB; parsing only, no approval/content/backend authority | **BUILD-NARROWLY** for SDK projections and CLI; consider code generation only from frozen authority-classified schemas. Avoid broad CLI frameworks until command complexity proves need. | Phase 2/3 public cutover; generated-file provenance/idempotence, strict unknown-field/numeric behavior, field-authority completeness, help/error text contains no guest/source data, exact protocol fixtures |
| Broker presentation controls and accessibility | SwiftUI/AppKit + LocalAuthentication | Approval-understanding TCB; high consequence UI | **ADOPT-PLATFORM** widgets/accessibility, **BUILD-NARROWLY** fixed review interaction. No web view/HTML renderer or third-party design system in approval UI. | Phase 3 UI-001; typed registered-plan rendering, focus/overlay/synthetic input, accessibility tree, localization/truncation, stale data/session, screenshot redaction and user-presence proof |

## Prioritized implementation dependency graph

```text
ADOPT NOW (standards/platform/test choices; no product activation)
  APL-1..6 + GO-3 + STD-1 + TEST-1/RUST-1 + SUP-1 formats
       |
       +--> M1 passive .mjs source/manifest (active owner)
       +--> F2 fixed-store v2 oracle (active owner)
       +--> G2 owner/store startup composition (active owner)
       +--> S1 native XPC fixtures after M1/M2 contract inputs

BOUNDED DECISIONS
  [CBOR/COSE production profile] (complete: narrow fxamacker selection; go-cose product NO-GO)
       +--> Phase 2 signed-object/profile freeze
  [SQLite production-engine comparison] <-- F2 full verifier + G2 logical ownership
       +--> archive/store engine ADR only if evidence selects it
  [.mjs parser boundary] (already active; not owned by this audit)
       +--> M1 validator implementation; never replaces runtime no-loader evidence

LATER ACTIVATION
  ARM64 rusty_v8 handoff --> independent builders/source+notice closure --> runtime admission
       --> governed runtime + governed libkrun owned-guest composition
  Apple identities/profiles --> signed/notarized complete bundle --> clean/minimum-OS hosts
  TUF updater --> verified local TrustSnapshot --> production update/repair activation

DO NOT ADOPT
  full Bun/Deno | generic XPC/JSON/NSXPC bus | hand-written JS lexer
  live TUF/DID/network resolution | SQLCipher/ORM by default | rich parser in daemon/Supervisor
  OCI registry/image framework for v0 raw root | web approval UI | cross-repo runtime library
```

Only the three named spike decisions above are proposed. Runtime, backend, signing/distribution,
and independent-builder work remains the already-defined evidence program rather than new reuse
spikes.

## Dependency-policy checklist for future tasks

Copy this block into the task plan or retained decision record before a dependency is added:

```text
Capability and roadmap slice:
Reuse-map row and recommendation:
Candidate name, ecosystem, exact version/commit/content digest:
Primary official sources and retrieval date:
Maintenance/release cadence and named advisory/reporting owner:
License(s), notices, source-offer/corresponding-source obligations:
Artifact provenance/signatures and what was actually verified:
Supported OS/architecture/toolchain/language boundary:
Direct dependencies and complete transitive graph/lock digest:
Binary/source/runtime-root size and TCB added/removed:
Trust class: product | approval-understanding | build | test | evidence
Authority gained: key | storage | parser | network | process | filesystem | guest/backend | update
Required features and features physically omitted/disabled:
Input/output byte, depth, count, time, memory, concurrency, and queue caps:
Fault, crash, cancellation, partial-write, corruption, and overload behavior:
Offline/reproducible fetch and build procedure:
Positive, negative, cross-language, restoration, and mutation tests:
Vulnerability monitoring, response SLA/owner, and upgrade trigger:
Upgrade/rebase cadence, rollback/removal plan, and retained compatibility fixtures:
Decision and evidence gaps: ADOPT-PLATFORM | ADOPT-PINNED | GOVERN/FORK |
  SPIKE-FIRST | TEST-ONLY | DEFER | BUILD-NARROWLY
```

Admission fails closed if any applicable item is unknown. An explicitly recorded `unknown` can
justify `SPIKE-FIRST` or `DEFER`; it cannot justify production activation.

## Scope and limitations

This audit used repository source, manifests, retained experiments/evidence, governed-fork records,
and read-only official sources. It did not install or execute a candidate package, inspect private
fork state, use signing identities, create a guest, or validate product IPC/runtime/backend/store
behavior. Package release pages establish current public metadata, not independent source or
artifact provenance. Transitive counts remain unknown where no bounded spike has resolved the
candidate with the pinned Capsule toolchain.

The map therefore changes planning order and prevents redundant invention; it does not establish
`DAEMON-001`, `UI-001`, `RUNTIME-001`, backend/profile admission, production signing, protected
storage, continuous service, independent builds, or production readiness.
