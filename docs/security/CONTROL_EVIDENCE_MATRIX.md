# Control Evidence Matrix

Status: living claim registry; all initial rows are proposed.

No security claim advances because it appears in architecture text. A claim becomes implemented
only when its mechanism exists, and validated only when exact attack tests and retained evidence
cover the pinned configuration.

Status values:

- `proposed`: design only;
- `spike-observed`: observed in a named non-production prototype/environment;
- `implemented`: mechanism exists but authoritative validation is incomplete;
- `validated`: exact pinned configuration passed the required retained tests;
- `degraded`: formerly supported claim cannot currently meet its evidence/freshness requirement;
- `unsupported`: mechanism is unavailable and execution requiring it must refuse or the contract
  must change.

| Claim ID | Claim | Threat | Owner | Mechanism | Live verification | Required attack tests | Receipt/transcript evidence | Known limitation | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| AUTH-001 | One approval authorizes at most one attempt | Replay, concurrent double execution | Broker + Supervisor | Attempt nonce plus durable atomic grant consumption | Grant-ledger transaction and restart reconciliation | Replay, concurrent execute, crash after consumption | Grant digest, consumption event, attempt ID | Denial of service can burn approval | proposed |
| AUTH-002 | Daemon cannot approve or launch | Compromised agent-facing service | Broker + Supervisor + macOS | Keychain isolation, XPC peer requirements, Supervisor-only backend endpoint | Peer/key access checks | Daemon key use, IPC impersonation, direct launch attempts | Integrity assessment peer results | Trusted admin/kernel out of scope | proposed |
| PLAN-001 | Executed bytes equal registered/approved plan | Substitution after approval | Supervisor + Broker | Exact durable registration; execute by ID only; independent digest | Registration lookup and staged digest checks | Plan A/B substitution, modified bytes, stale registration | Plan digest, registration ID/sequence, grant binding | Correct signed logic still trusted | proposed |
| SAFE-001 | V0 unsupported powers cannot be enabled by daemon policy | Policy bug or daemon compromise | Supervisor | Versioned non-overridable hard-safety validator | Independent plan check before registration/attempt | Network/process/env/package/image/path fields and type confusion | Validation version and result | Validator correctness remains TCB | proposed |
| RI-001 | Trusted peers match enrolled component identity | Same-user impersonation, partial update | macOS + Supervisor/Broker | XPC peer requirement, exact build/epoch checks | OS-enforced connection/message check and preflight | Unsigned, wrong ID/team, stale/debugged peer | RuntimeIntegrityAssessment | Compromised kernel out of scope | proposed |
| KEY-001 | Daemon cannot use Approval or Supervisor evidence key | Authority collapse | Keychain + component identity | Component-specific access group/access control | Failed key operation from daemon | Direct and injected daemon access attempts | Key authorization IDs; no secret material | Admin/keychain compromise out of scope | proposed |
| DATA-001 | Daemon cannot retrieve user-only content | Agent content bypass | Content Broker | Separate store/API plus attempt-scoped handles | Credential/ACL/handle checks | MCP/daemon read, forged/stale/cross-job handle | Delivery audit record and content digest | Supervisor/Broker necessarily access scoped bytes | proposed |
| DATA-002 | Guest input equals approved immutable bytes | Live-file race/path substitution | Broker + Supervisor | Snapshot content ID and staged digest verification | Digest before guest start | Mutation, link/special-file, partial transfer, cross-job reuse | Input manifest and staged digest event | Original path semantics intentionally excluded | proposed |
| NET-001 | Guest has no unauthorized network/IPC | Exfiltration, metadata/host-service access | Supervisor + backend | Proven no-interface/network-none and no host socket mechanism | Exact config inspection plus guest probes | TCP/UDP/DNS/IPv4/6/loopback/Unix/vsock/metadata | Backend config/capability/validation digest | Backend/kernel compromise out of scope | proposed |
| RES-001 | Required resource controls are externally enforced exactly | Host availability loss | Supervisor + backend/host | Mechanism per limit dimension | Capability report and runtime observation | CPU, memory, PID, disk, output, cancellation abuse | Exact approved values, mechanism IDs, violation events | Vocabulary pending Apple spike | proposed |
| CLEAN-001 | Every created guest is destroyed or explicitly unresolved | Persistent hostile execution | Supervisor + backend | Durable handle/cleanup lease and reconcile lifecycle | Backend enumeration and teardown observation | Crash/cancel/timeout/orphan/partial destroy | Destroy/reconcile events and terminal class | Missing handle never proves absence | proposed |
| EGR-001 | Agent receives only fixed default summary | Content/metadata exfiltration | Broker + daemon | Fixed typed response with no guest strings | Response schema enforcement | Names/sizes/timing/errors/log/output injection | Summary schema/version and release decision | State/timing channel remains | proposed |
| TRUST-001 | Partial component update fails closed | Stale/mixed components | Installation Trust Domain | Signed manifest and common epoch binding | Connection and preflight epoch comparison | Old/new component combinations and restored stores | Epoch number/digest and assessment result | Coherent rollback needs stronger anchor | proposed |
| TUF-001 | External trust data is locally verified without live execution dependency | Malicious/stale service, rollback | Updater/trust verifier | Pinned TUF root plus signed local TrustSnapshot | Version/expiry/hash/checkpoint verification | Freeze, rollback, mix-and-match, delegated-scope abuse | Trust snapshot/checkpoint/freshness | Capsule defines revocation semantics | proposed |
| EVID-001 | User receipt contains attributable approval and enforcement claims | Daemon-forged success | Broker + Supervisor | Purpose-separated signatures over same registration/attempt | Independent signature and binding verification | Remove/swap/replay/cross-object evidence | Embedded grant, transcript, authorizations | Not platform attestation or proof of truth | proposed |

## Update rule

When a mechanism or spike changes, update the corresponding row with an evidence link, exact
version/configuration, and limitations. Never promote an entire backend because one row passes; all
claims required by the requested posture must be validated together.
