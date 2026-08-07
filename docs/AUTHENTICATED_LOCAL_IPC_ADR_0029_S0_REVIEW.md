# ADR-0029 S0 decision review

Date: 2026-08-07

Status: `PASSED`

Parent authenticated product IPC: `BLOCKED`

Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`

## Defensive and authorized scope

Review the already-proposed authenticated local IPC topology using only repository contracts,
fixtures, local mechanics, and retained evidence. This review activates no listener, service,
credential, Keychain item, protected root, store consumer, runtime, backend, process, or guest. It
does not access any system, identity, credential, or data outside the Capsule repository and its
retained controlled evidence.

## Question

Should ADR-0029's core topology be accepted, revised, or rejected before the native XPC research
and controlled harness proceed?

The review tested these invariants from the authenticated-local-IPC S0 packet:

1. one unprivileged per-user Supervisor process and no lifecycle helper;
2. exactly two ordinary role-specific Supervisor services and four ordinary calls;
3. ADR-0044's CLI-to-daemon submission service remains separate from that Supervisor surface;
4. native authentication precedes application-body copy and decode;
5. every native-to-Go bridge is method-specific and copy-only, with no opcode bus;
6. Go remains the sole durable authority, attempt, lifecycle, and recovery owner;
7. request IDs remain live-call correlation only;
8. an `AttemptID` is durably created before any lifecycle effect; and
9. startup and recovery enumerate and accept committed `AttemptID` values only.

## Method

The review compared ADR-0029 with:

- Accepted ADR-0034, ADR-0040, ADR-0043, and ADR-0044;
- the passive `RegisterPlanV0`/`GetRegisteredPlanV0` authority-plane facade;
- the three-method native-XPC S3 dictionary and refusal fixtures;
- the public-key ApprovalGrant-to-FakeBackend integration;
- owner-required `AttemptID`-only startup and recovery; and
- the current dependency-ordered work plan.

It inspected the existing Go method boundaries and retained tests but ran no process or guest. The
review separated topology selection from product activation and platform evidence.

## Decision result

Accept ADR-0029's core topology with the following reconciliations. None changes component
responsibility, adds a helper, activates a listener, or admits a product path.

### Retained invariants

- The Execution Supervisor remains one native-fronted Go process and the sole guest lifecycle
  owner.
- The daemon and Broker use two different Supervisor services; neither service accepts a
  daemon-or-Broker disjunction or generic command envelope.
- `RegisterPlanV0`, `GetRegisteredPlanV0`, `SubmitApprovalV0`, and `RequestAttemptV0` remain the
  complete ordinary Supervisor method set. The separate `SubmitMainMJSV0` CLI-to-daemon adapter
  selected by ADR-0044 is not a fifth Supervisor method.
- The listener requirement, session checks, message-derived code identity, flow slot, outer shape,
  installation/epoch binding, and method binding precede application-body copy and Go decode.
- The native layer owns only live platform objects, copied buffers, and flow/deadline state. It
  cannot create or reinterpret durable authority.
- Registration retains fresh-registration semantics; approval and attempt replay use their
  existing semantic identities; request IDs never become durable deduplication keys.
- Startup opens under one installation owner and recovers only the sorted, duplicate-free set of
  committed `AttemptID` values. No caller can supply recovery bytes, backend selections, or paths.

### Reconciled ordering

The exact existing three-method S3 contract covers `SubmitMainMJSV0`, `RegisterPlanV0`, and
`GetRegisteredPlanV0`. Its controlled native harness may proceed after the native-XPC research
brief and separate exact authorization; it does not wait for the two remaining Supervisor method
envelopes.

`SubmitApprovalV0` and `RequestAttemptV0` remain a separate passive C4 slice. C4 must freeze their
numeric tags, exact dictionaries, case-derived caps and deadlines, refusal precedence, copy rules,
replay behavior, and response-loss fixtures before either method can become an installed
consumer. The existing 528/32-byte arithmetic and five-second deadline are design candidates, not
active or fixture-proven transport authority.

### Reconciled refusal ownership

The S3 native contract is canonical for the three methods it freezes. Its closed numeric status,
reason, and precedence tables distinguish:

- OS peer-requirement rejection, which produces no application delivery or reply;
- delivered outer-envelope refusal before body copy;
- a recognized method-specific core refusal; and
- a local bridge-integrity fault, which terminates the process without an application reply.

ADR-0029 therefore does not maintain a competing flattened refusal table. C4 must extend the same
model for the remaining two methods without renumbering or reinterpreting the existing three.

### Reconciled storage gate

Accepted ADR-0040 permits the owner-only internal alpha to use the bounded fixed-store exception
only after its exact owner, full-verification, refusal-threshold, no-restore, and installation-
retirement rules are enforced. F6 or another reviewed production engine remains mandatory before
external alpha, restore, continuity, multiple users, multiple state-opening processes, or retained
non-disposable user data. Production archive selection is no longer an unconditional prerequisite
for the narrow internal-alpha IPC consumer.

## Findings disposition

| Finding | Severity | Disposition |
| --- | --- | --- |
| ADR acceptance depended on harness and installed evidence while S0 required an accepted-or-revised decision before the harness | Required alignment | Separate design acceptance from implementation/evidence gates; keep activation `BLOCKED` |
| The older IPC dependency graph put both remaining methods before the existing three-method S3 harness | Required alignment | Split S3 harness and C4 into independent prerequisite lanes |
| The consequences retained an unconditional production archive/compaction blocker | Required alignment | Apply ADR-0040's exact internal-alpha exception; retain F6 for external-alpha claims |
| ADR-0029's refusal prose predates the S3 numeric status/reason/precedence contract | Required alignment | Make S3 canonical for its three methods and require C4 to extend it |
| Approval/attempt caps and deadlines were presented as exact before case-derived C4 fixtures exist | Required alignment | Retain them only as candidates until C4 passes |

No code correctness defect, activated security boundary, or product-admission evidence follows
from this document-only review.

## Verification

The retained document alignment is verified by:

```sh
pnpm verify:adrs
pnpm check
```

The future S3 harness still requires the exact defensive authorization and evidence boundaries in
the current work plan. The future C4 slice remains passive and may not activate a listener,
signer, store consumer, process, runtime, backend, or guest.

## Confidence and limitations

Confidence is high that the selected process, service, authority-ownership, copy, replay, and
recovery topology remains consistent with the current repository mechanics and accepted alpha
posture. Confidence in exact platform enforcement remains intentionally unset until the R1 brief
and separately authorized S3 harness complete. Exact `SubmitApprovalV0` and `RequestAttemptV0`
transport values remain unset until C4.
