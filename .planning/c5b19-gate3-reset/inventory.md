# C5b19 Gate 3 reset — immutable inventory

Inventory status: **PASSED** (read-only inventory completed). Gate 3 implementation remains **IN_PROGRESS — TRENDING_GOOD**; child execution and product admission remain **BLOCKED**. This report makes no timing, custody, security, Gate-3-completion or new test-pass claim.

## Scope and method

Defensive repository inspection only: canonical packet at Capsule commit `a9e1d32afbb1fa1ef9c94175321c9c3bbb3d3361` (P), and experiment Git tree at archive commit `a9c457cb3b89a4ecd8c38684fd3d804d2089214c` (A). All archive paths below are relative to `experiments/typed-guest-transport-c5b19-benign-teardown/`. Exact source was read with `git show A:<path>` and presence with `git ls-tree -r --name-only A -- <experiment>`, independent of working-tree additions. No implementation, process fixture, tests, live store, native code or child execution occurred.

Graph-first discovery: `capsule-corp` index is ready but points at `/Users/dsteele/repos/capsule-corp`, indexed 2026-10-04. Structural search for the current oracle symbols returned no rows. `check_index_coverage` explicitly returned `outside_project / unavailable` for the archive oracle path. Consequently this report uses exact-commit source fallback, not negative graph results, to establish coverage and missing files.

P's frozen manifest (`docs/C5B19_GATE2_IMPLEMENTATION_PACKET.md:49–99`) has **20 required paths: six present, fourteen absent**. These counts describe file presence, not a completion fraction. Generated binaries/header inventory belongs to Gate 4 and run outputs to Gate 5 (P:62–65); neither is credited as Gate 3 implementation.

## Every manifest path

“Complete” below means the narrow artifact obligation only, not its parent gate. “Partial” means useful implemented bytes exist but the required integrated role is incomplete. Missing paths are established by the immutable tree enumeration; there is no file line to cite for absence.

| Required path | Tree presence | Functional assessment / evidence |
| --- | --- | --- |
| `PLAN.md` | Present | Partial: current scope/review/history ledger; not full implemented experiment. Explicit no-child scope A:3–4, authority gate A:71–83. The old original-slice paragraph A:27–30 is historical, not a current per-feature inventory. |
| `go.mod` | Present | Complete narrow module/toolchain declaration: module `capsule.experiments/c5b19`, Go `1.25.13`, no requirements (A:1–3). Does not establish source/build closure. |
| `inputs/ORIGINS.json` | Present | Partial: immutable-source and increment provenance ledger; not completed native/build/run manifest. `run_authorized: false`, `product_admitted: false` (A:11–12). Historical first G3a third-review verdict remains uncorroborated (A:28; PLAN:66–69). |
| `inputs/fixture.c` | Absent | Missing fixed inert child implementation. |
| `inputs/supervisor.c` | Absent | Missing native control state machine, sole child custody, transport, clocks and observations. |
| `inputs/harness.c` | Absent | Missing sole parent/waiter/cancellation writer and post-containment oracle exec. |
| `inputs/protocol.h` | Absent | Missing frozen C protocol declarations. |
| `inputs/protocol.c` | Absent | Missing C codec/transport implementation and C/Go correspondence. |
| `inputs/clock.h` | Absent | Missing native clock interface. |
| `inputs/clock.c` | Absent | Missing clock sampling and checked arithmetic implementation. |
| `inputs/bridge.h` | Absent | Missing frozen sole-worker bridge ABI. |
| `source/bridge/bridge.go` | Absent | Missing fixed request/settled-result c-archive bridge. |
| `source/store/store.go` | Absent | Missing serialized owner-locked record and durable fault-site implementation. |
| `source/store/store_test.go` | Absent | Missing store fault/restoration/concurrency tests. |
| `tests/model_test.go` | Present | Partial: process-identity vectors and synthetic bounded write-history tests. Explicitly byte histories, not actual control events (A:78); first matrix A:164. No complete independent C5b18 transition adapter conformance. |
| `tests/oracle.go` | Present | Partial: pure in-memory validators and supplied-history correspondence. `package oracle`, imports only bytes/sha256/binary/errors (A:1–10); no executable entry point or evidence-file ingestion/final case verdict. Families below. |
| `tests/oracle_test.go` | Present | Partial: substantial parser/correspondence/unit coverage, with explicit refusal and ownership tests. No full lifecycle trace/timing/case or executable integration test suite. |
| `scripts/build.sh` | Absent | Missing reproducible frozen toolchain/native/c-archive/oracle build. |
| `scripts/check.sh` | Absent | Missing controlled pre-run check entry point and full synthetic/negative/mutation suite integration. |
| `scripts/run.sh` | Absent | Missing fixed-case allowlist and exec-only harness handoff. Child execution still prohibited. |

## Reusable implemented families (A)

These are useful foundations. Every row has production code and corresponding tests in the immutable tree. This audit inspected source; it did not rerun those tests or independently reproduce archived historical mutant results.

| Implemented family | Source / representative test anchors | Boundary still outstanding |
| --- | --- | --- |
| Cancellation frame, process identity, obligation preimage | `oracle.go:23,43,78`; `oracle_test.go:32,154`; `model_test.go:15,52` | Binding/byte predicates do not establish an observed process or actual cancellation ingress. |
| Request397/reply98/record395 decoding and hash correspondence | `oracle.go:137,144,200,217`; `oracle_test.go:256,369,426,525` | No C encoder/parser, serialized file loading, observed transport or durable write. |
| One-pending/three-operation synthetic write history | `oracle.go:275–384`; `model_test.go:164,199,235,370` | Explicitly not the control model, store, timeout or record-loading path. Only a simulated accepted application advances confirmed projection. |
| Event envelope, bounded lane framing | `oracle.go:417,457`; `oracle_test.go:594,762,796,845` | Common fields/framing; no capture completeness, causality or timing proof. Lane capacity checks alone are not required full control-loop capacity simulation. |
| Kind-specific fields and lane ownership for **all 29 kinds** | `oracle.go:517–617`; `oracle_test.go:964,988,1013,1128` | All lifecycle kinds have syntactic shapes. Their lifecycle semantics are not implemented. Kind29 handles payload-free subset only; kind8 excludes malformed-frame diagnostic encodings. No Darwin terminal-status interpretation. |
| Captured raw-frame/event correspondence, kinds3/5/7 | `oracle.go:622–707`; `oracle_test.go:1246,1272,1393,1443` | In-memory selections, not serialized sidecar format. Kind5/7 request context alone does not prove independent worker receipt. |
| Received/applied correspondence, kinds7→8 | `oracle.go:709–773`; `oracle_test.go:1520,1561,1631,1746` | Exact paired data and local order, not valid decision/control-state transition or omitted-pair detection. |
| Supplied pending-write/decision reconciliation | `oracle.go:778–871`; `oracle_test.go:1808,1860,1931,2037,2128` | Serial supplied kind3 and whole7/8 groups; rejected decision monotonically latches fatal. Explicitly lacks stop/fatal admission/full control state (815–817). |
| Supplied record snapshot correspondence | `oracle.go:876–1009`; `oracle_test.go:2202,2245,2307,2419,2597` | Absent/older/pending/failed observation does not prove a consistent final store or settle control. |
| Supplied independent worker-received/sent sidecars | `oracle.go:1014–1111`; `oracle_test.go:2680,2733,2809,2942,3052` | Bounded optional evidence, no serialized format, capture completeness/authenticity, chronology or syscall success. |
| Supplied write/result pairs, kinds3→4 and5→6 | `oracle.go:1136–1196`; `oracle_test.go:3107,3140,3254,3280,3333` | Raw signed result/errno and exact numeric byte count; no full syscall-success semantics, inter-pair chronology, missing-event detection or retry behavior. |
| Supplied transport lane consistency, kinds3–8 | `oracle.go:1201–1250`; `oracle_test.go:3453,3478,3507,3576,3674,3697` | At most18 events/three perkind/three operations; gaps/empty/lone kind8 allowed. Does not compose pair/frame/history obligations. Tests at3674 expressly accept lane-consistent but pair-inconsistent data; tests at3453 expressly accept lone kind8. |

## Major incomplete behavioral obligations

1. **One integrated, bounded, read-only evidence verifier.** Add exact serialized ingestion/selection provenance, required-event and sidecar completeness, executable digest/frozen bindings checks, composition of the existing validators, cross-lane identity joins without total timestamp ordering, durable snapshot reconciliation, and one independent case verdict ignoring Supervisor pass flags. P:350–360,377–393,430–444. An optional in-memory helper returning success for an empty selection is valid locally but cannot be the complete-verifier acceptance contract.
2. **Independent full lifecycle/model conformance and timing oracle.** Preparation before spawn; returned-PID-derived identity before start; no start after stop; cancel/setup/wall/fatal races; same-tick stop/apply ordering; one signal; exact paired wait/terminal-status/absence; terminal after absence/no public completion; checked clock products and actual-service/post-start bounds. P:79–99,289–300,366–393,432–444. Existing 29-kind shapes and write-history predicates are reusable pieces, not this implementation.
3. **Native Supervisor/harness/fixture and storage separation.** All required source paths are absent. The control lane's no-store dependency, single-worker restriction, fixed descriptor inventory, sole wait/custody and cancellation ownership remain design obligations. P:67–99 and101 onward.
4. **Build and complete no-child conformance/mutation harness.** C/Go ABI/bytes, store faults, descriptor-negative fixtures, simulated process/store effects, 6,999 paired waits plus512 other events, full/overflow budgets, immediate-ready/EOF removal and EINTR storms remain unfinished. P:410–429,446–465. Mere presence of envelope cap tests does not close the control pacing/overflow simulation requirement.

## Reset recommendation

Count progress by three integrated outputs rather than each helper: (1) complete offline model/evidence verdict with full synthetic lifecycle/timing failures; (2) frozen native and store/bridge implementation with no-child conformance and fault tests; (3) reproducible source/build closure and whole-experiment pre-run integration/mutation evidence. Root should reconcile their exact dependencies and exit criteria against the requirements audit. Require every existing helper to be consumed by, or explicitly justified outside, the integrated verifier; avoid another open-ended sequence of supplied-selection APIs. Do not discard the tested byte/correspondence foundation.

Gate 3 completion means implemented bytes and no-child tests, followed by Gate 4's independent review of exact source, generated/runtime closure and binaries plus separate owner first-run authorization. Actual child/live evidence is Gate 5, not an excuse to run during Gate 3 (P:446–453). Child execution remains **BLOCKED** throughout this reset.

Confidence: high for immutable manifest presence, private API boundaries and the explicitly documented missing semantics; moderate for completeness of every conceivable edge-case test gap because this is a scoped source inventory, not a fresh implementation review or test run. Historical review/provenance claims were not revalidated against external review systems.
