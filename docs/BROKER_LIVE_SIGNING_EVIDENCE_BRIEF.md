# Approval Broker live-signing evidence brief

Date: 2026-08-11

Work item: R3 read-only Apple-platform research for a future installed Approval Broker signing
harness.

Status: `PASSED` for the exact passive research, authority reconciliation, candidate definition,
and future-evidence-plan scope.

Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`.

Disposable installed Broker signing evidence, authenticated product consumers, the installed
security boundary, and product admission: `BLOCKED`.

This document canonicalizes a supplied read-only research packet against the repository state at
`origin/main` commit `ed4220fe16d1752a75c67da957a25681d79e34f3`. That state includes the merged
passive C5a transport contract in PR #246 and the passive C4 approval/attempt candidate in PR #247.
The later CL4 audit is `PASSED` with historical disposition `AMEND`. PR #248 is the canonical
predecessor to the focused follow-up that closed the audit's exact passive evidence findings, so
the passive/no-listener C4 claim is `PASSED`. Neither merge nor the follow-up activates a platform
or product path.

## Defensive scope

This result uses repository documents and public Apple documentation only. It created, inspected,
used, or changed no key, Keychain item, identity, provisioning profile, credential,
`LAContext`, authentication prompt, installed service, IPC endpoint, runtime, backend, VM, guest,
or external data.

The future installed harness described here is not authorized by this document. It remains
`BLOCKED` until a separate task names the owned Mac, exact development-signed Broker test target,
profile and access group, disposable user/container, immutable fixture digest, permitted Keychain
mutations, and every destructive row. The harness may have no product consumer, runtime, backend,
VM, or guest.

Subsequent construction work now retains an immutable
[unsigned Broker harness](https://github.com/Shrimpworks/capsule-experiments/tree/4a2447d4bd0e03132dc616e608031ca313630cdd/experiments/broker-live-signing-c6b1)
and [test-only Supervisor seam](https://github.com/Shrimpworks/capsule-experiments/tree/067fe2beb40361bb714507cab1331004e0a656fa/experiments/broker-live-signing-c6b1-supervisor-seam).
Those C6b1a/b construction slices are `PASSED` for deterministic public fixtures, no-credential
interaction checks, and Supervisor-only commit/replay/response-loss/reopen/concurrency oracles.
They do not change this document's installed-evidence boundary: no Apple identity/profile,
Keychain, LocalAuthentication, private key, signing, installation, authenticated listener, or
product consumer was used.

The later no-install
[C6b1c signed-artifact readback](https://github.com/Shrimpworks/capsule-experiments/tree/82d1a799f70482856aaa6030f612d701b39cec67/experiments/broker-live-signing-c6b1c-signed-artifact-readback)
is also `PASSED` in its exact bounded scope. It retains the selected development-profile metadata
and exact signed Broker artifact with Team `3DDR84M4JS`, bundle ID
`com.capsulecorp.capsule.broker.c6b1`, CDHash
`029b8d5cabd38e1fde9e23564e4e5b1590cf569d`, hardened runtime, App Sandbox, and exactly one
Approval Keychain access group
`3DDR84M4JS.com.capsulecorp.capsule.broker.approval.epoch-7`. The profile wildcard is an allowlist,
not proof of key access. The raw profile was not embedded, the app was not installed or launched,
and no Keychain, LocalAuthentication, Secure Enclave, service, IPC, or product operation occurred.

## Governing decisions and evidence boundary

[Accepted ADR-0043](adr/0043-freeze-broker-rendering-and-approval-verification.md) remains the
governing decision. It requires all of the following without this research changing their status:

- a nonexportable Secure Enclave P-256 Approval key with no software fallback;
- exact `userPresence | privateKeyUsage` access control;
- one fresh, nonreused `LAContext` for each signing operation;
- public-key authorization bound to Team, Broker role, exact access group, installation, epoch,
  purpose, audience, validity, active status, and key identity;
- the closed canonical tagged COSE_Sign1 profile with a 64-byte raw `R || S` signature;
- purpose `capsule.plan.approve` and audience `capsule.execution-supervisor`;
- registration, plan, Supervisor, attempt nonce, installation, epoch, issuance, and expiry
  bindings; and
- canonical-payload replay, one-use consumption, and attempt creation owned by the Supervisor.

[Accepted ADR-0029](adr/0029-select-authenticated-local-ipc-topology.md) keeps the native front end
non-authoritative and the Go Supervisor core as the only durable owner of registration,
approval, attempt, and lifecycle state. The merged passive
[C4 candidate](AUTHENTICATED_LOCAL_IPC_C4_NATIVE_CONTRACT.md) records `SubmitApprovalV0` and
`RequestAttemptV0` encodings, replay, and response-loss behavior without implementing their
listeners or consumers. Its exact passive/no-listener evidence claim is `PASSED`; CL4 remains
`PASSED` with historical disposition `AMEND`, now closed by the focused follow-up. Installed IPC,
signing, consumers, and product admission remain `BLOCKED`.

[Proposed ADR-0021](adr/0021-security-epoch-keychain-groups.md) remains Proposed. This document
does not accept its identity-changing key/group transition policy.

## Durable approval authority stays Supervisor-owned

The durable `SubmitApprovalV0` commit in the Supervisor is the approval-authority linearization
point. The Broker does not become a second durable authority owner.

The Broker may hold one bounded, non-authoritative in-memory interaction while its process lives:
the immutable fetched/rendered projection, canonical payload and protected header, exact
`Sig_structure`, one fresh context, one sign-call budget, returned signature/envelope bytes, and
submission status. Process death discards all of it. This brief selects no Broker envelope journal,
cache, recovery store, durable replay index, or state authority.

Any future durable Broker cache or journal, recovery owner, or state-authority assignment requires
an explicit reviewed decision and a complete crash-consistency design. It must not silently move
approval admission away from the Supervisor.

The required crash and retry behavior is:

| Boundary | Authority result | Retry behavior |
| --- | --- | --- |
| Before signing | No authority exists. | Start a new fetch, render, context, authentication, and signing interaction. |
| Signature returned but no Supervisor commit | The signature or envelope in Broker memory is not approval authority. Broker crash loses it. | A later request starts a new interaction; it must not recover authority from logs, caches, or partial bytes. |
| Envelope submitted and Supervisor commit is indeterminate | The Broker cannot decide authority from reply absence. Supervisor recovery/store truth controls. | While alive, the Broker may resubmit the exact in-memory envelope. After Broker death, a fresh interaction may create an equivalent envelope over the same canonical payload. |
| Supervisor committed, but reply was lost or the Broker crashed | The Supervisor already owns one durable approval record. | An exact or mathematically equivalent later submission converges at the Supervisor to the same `ApprovalID` and current state. Identical envelope bytes are not promised after Broker death. |
| Supervisor consumed approval and created an attempt, but the reply was lost | One attempt exists in the Supervisor transaction. | `RequestAttemptV0` replay converges to the same `AttemptID` and state. |

Canonical payload plus resolved signer-authorization identity remains the `SubmitApprovalV0`
replay identity. Signature bytes are not authority identity. Complementary mathematically valid
ECDSA signatures over one canonical payload cannot create two approvals.

## Apple-documented mechanism

Apple documents the following relevant building blocks:

- `LAContext.evaluateAccessControl(..., .useKeySign, ...)` can authorize use of a key's access
  control, while a later private-key operation remains separately fallible. See
  [evaluateAccessControl](https://developer.apple.com/documentation/localauthentication/lacontext/evaluateaccesscontrol(_:operation:localizedreason:reply:)).
- An authorized context bound to a key can permit additional signing operations until the context
  is invalidated. Apple's [WWDC22 local-authorization session](https://developer.apple.com/videos/play/wwdc2022/10108/)
  makes this explicit. Capsule's one-sign rule is therefore an application invariant, not a
  Keychain-provided consumption property.
- [`LAContext.invalidate()`](https://developer.apple.com/documentation/localauthentication/lacontext/invalidate())
  invalidates the context and cancels pending evaluation, but does not establish that an already
  computed signature was retracted.
- [`touchIDAuthenticationAllowableReuseDuration`](https://developer.apple.com/documentation/localauthentication/lacontext/touchidauthenticationallowablereuseduration)
  governs reuse of recent Touch ID device unlock. Setting it to zero does not make an already
  authorized context a one-operation capability.
- [`interactionNotAllowed`](https://developer.apple.com/documentation/localauthentication/lacontext/interactionnotallowed)
  controls interactive authentication. Explicit approval sets it false; noninteractive readback
  and preflight must fail rather than unexpectedly prompt or fall back.
- Local Authentication returns cancellation, invalid-context, lockout, unavailable, and
  noninteractive failures. See [`LAError.Code`](https://developer.apple.com/documentation/localauthentication/laerror-swift.struct/code).
- [`userPresence`](https://developer.apple.com/documentation/security/secaccesscontrolcreateflags/userpresence)
  is documented as biometry-any or device-passcode presence, and
  [`privateKeyUsage`](https://developer.apple.com/documentation/security/secaccesscontrolcreateflags/privatekeyusage)
  is required for Secure Enclave private-key operations.
- Secure Enclave P-256 keys are generated with
  [`SecKeyCreateRandomKey`](https://developer.apple.com/documentation/security/generating-new-cryptographic-keys),
  their public key is obtained with
  [`SecKeyCopyPublicKey`](https://developer.apple.com/documentation/security/seckeycopypublickey(_:)),
  and their private representation is not exportable. See
  [`SecKeyCopyExternalRepresentation`](https://developer.apple.com/documentation/security/seckeycopyexternalrepresentation(_:_:)).
- Algorithm compatibility is checked with
  [`SecKeyIsAlgorithmSupported`](https://developer.apple.com/documentation/security/seckeyisalgorithmsupported(_:_:_:)),
  and signing uses
  [`SecKeyCreateSignature`](https://developer.apple.com/documentation/security/seckeycreatesignature(_:_:_:_:)).
- A Keychain item belongs to one access group. The installed profile and signed effective
  entitlements must authorize that group; an App Group would add broader container and IPC
  authority and is not a key-only substitute. See Apple's
  [Keychain sharing guidance](https://developer.apple.com/documentation/security/sharing-access-to-keychain-items-among-a-collection-of-apps),
  [`keychain-access-groups`](https://developer.apple.com/documentation/bundleresources/entitlements/keychain-access-groups),
  and [App Groups](https://developer.apple.com/documentation/bundleresources/entitlements/com.apple.security.application-groups).

These mechanisms do not prove Capsule's exact installed behavior, key authorization, one-sign
state machine, focus races, recovery, or product wiring.

## Experiment candidates, not product policy

Two values are intentionally only experiment candidates:

| Candidate | Current disposition | Required result |
| --- | --- | --- |
| `kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly` | Sole accessibility candidate for the disposable evidence harness. It is not accepted/frozen product policy or installed evidence. | The exact owned Mac, user, OS, hardware, Keychain mode, lock/password/passcode transitions, deletion behavior, and recovery outcome must match the fail-closed contract. |
| `kSecKeyAlgorithmECDSASignatureMessageRFC4754SHA256` | Preferred signing-algorithm candidate because Apple SDK headers describe the RFC 4754 family as fixed-width concatenated `r || s`. It is not proven target support, output, or runtime policy. | The exact target key reports support; output is exactly 64 bytes; Broker same-byte verification and the Supervisor verifier accept it; every alternate encoding refuses. |

There is no silent fallback. If either candidate fails, live activation remains `BLOCKED` and a
reviewed refinement must select the next path. In particular:

- `kSecAttrAccessibleWhenUnlockedThisDeviceOnly` or another weaker accessibility class is not an
  automatic fallback;
- X9.62/DER output plus a DER-to-raw conversion is a separately governed contingency, not runtime
  negotiation;
- a software key, imported key, migrated key, App Group, broader access group, different ACL, or
  alternate signer is not permitted; and
- direct-candidate failure does not weaken ADR-0043.

## Candidate installed key shape

The disposable harness should test exactly this candidate, without treating it as an activated
product policy:

```text
key class:             private EC key
key type:              kSecAttrKeyTypeECSECPrimeRandom
key size:              256 bits
token:                 kSecAttrTokenIDSecureEnclave
permanent:             true
synchronizable:        false
keychain mode:         data-protection Keychain on macOS
access group:          3DDR84M4JS.com.capsulecorp.capsule.broker.approval.epoch-<base10 epoch>
application tag:       fixed versioned Broker tag bound to Approval purpose and key generation
accessibility:         kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly (experiment candidate)
access-control flags:  userPresence | privateKeyUsage (ADR-0043 requirement)
software fallback:     forbidden
```

Key creation is a separately authorized enrollment transaction, never a side effect of an ordinary
approval request. Missing, duplicate, inaccessible, unexpected, or attribute-mismatched key state
refuses approval and enters the separately selected initialization/repair path. The Broker cannot
self-authorize a new public key, silently replace a missing key, rotate authority, or reset trust.

Every query must name the exact class, application tag, access group, token, key class, and
data-protection disposition. Zero or multiple matches refuse. Readback must cover token, type,
size, class, permanence, synchronizability, sign capability, application label, public key, and
the fields the API makes observable. A private-key export attempt must fail, but that negative
result alone does not prove access-group, ACL, code-identity, or public-key authorization.

## One-interaction state machine

The Broker's candidate live state is bounded and process-local:

```text
IDLE
  -> FETCHED
  -> RENDERED
  -> AUTHORIZING
  -> AUTHORIZED
  -> SIGNING             # exactly one sign-call budget
  -> SIGNED_LOCAL_OK     # context invalidated; envelope remains in memory only
  -> SUBMITTING
  -> SUPERVISOR_COMMITTED
  -> REPLIED

Any nonterminal state
  -> CANCELED | REFUSED | FAILED
```

Every terminal transition invalidates the context if present, drops the private-key reference,
zeros the sign-call budget, and ignores callbacks from older interaction generations. The Broker
serializes one interaction process-wide. An identical active request may receive bounded
`IN_FLIGHT` status; a different request receives `BUSY` and is not queued behind an authenticated
context.

The candidate sequence is:

1. Acquire the single-flight interaction slot.
2. Fetch one defensive registered-plan view from the authenticated Supervisor Broker service.
3. Cross-bind registration, plan, source, installation, epoch, Supervisor, nonce, purpose,
   audience, validity, authorized public key, and rendered projection under ADR-0043.
4. Freeze the protected header, canonical payload, exact `Sig_structure`, their digests, the
   displayed projection, deadline, and interaction generation before authentication.
5. Preflight the exact key noninteractively; require one match and the authorized public key/kid.
6. Create a fresh `LAContext`, set interactive use explicitly, keep recent-unlock reuse at zero,
   and call `evaluateAccessControl(..., .useKeySign, ...)` with a bounded concrete reason.
7. After success, recheck generation, active key/window/session/lifecycle state, deadline, frozen
   identities, authorization status, and one-sign budget.
8. Retrieve the same key under the authorized context and recheck its public key/kid.
9. Atomically consume the one-sign budget and call `SecKeyCreateSignature` once over the frozen
   `Sig_structure` using the preferred experiment algorithm.
10. Immediately invalidate the context and drop the private-key reference on every result.
11. Require the exact raw 64-byte shape, locally verify the same bytes with the expected public
    key, construct the closed COSE_Sign1 envelope, and run the strict passive profile checks.
12. Submit the exact envelope to `SubmitApprovalV0`. Only the Supervisor's durable commit creates
    recoverable approval authority.
13. Return the Supervisor-issued `ApprovalID` and current state. Keep any reply-retry bytes only in
    bounded memory for the life of this interaction.

The signed bytes are the exact canonical COSE `Sig_structure`, never UI prose, a caller-provided
digest, mutable state, or a daemon-supplied display string.

## Cancellation and UI claim boundary

App/window/session/sleep notifications are conservative cancellation inputs, not proof that the
OS revoked an authorization before Capsule observed the event. The harness must cancel and
invalidate on app deactivation, key-window loss, session switch-out, screen/device lock, sleep,
logout/power-off, deadline, process termination, and selected Space/background changes. Every late
callback must be rejected by the immutable interaction generation and monotonic terminal state.

Relevant lifecycle hooks include
[`NSApplication.willResignActiveNotification`](https://developer.apple.com/documentation/appkit/nsapplication/willresignactivenotification),
[`NSWorkspace.sessionDidResignActiveNotification`](https://developer.apple.com/documentation/appkit/nsworkspace/sessiondidresignactivenotification),
and [`windowDidResignKey`](https://developer.apple.com/documentation/appkit/nswindowdelegate/windowdidresignkey(_:)).

Ordinary AppKit UI is not secure attention. The Local Authentication prompt, app name/reason,
focus state, and system-owned authentication data improve context but do not prove comprehension,
exact-build identity, overlay resistance, capture resistance, or immunity to Accessibility,
automation, screen sharing, recording, task-port access, or injected code. Apple's
[secure-intent guidance](https://support.apple.com/guide/security/secure-intent-and-connections-to-the-secure-enclave-sec7a94f7d1e/web)
does not generalize hardware secure-intent flows to an arbitrary macOS application window.

The narrow installed claim is only that synthetic UI interaction without the configured
LocalAuthentication/Keychain-gated private-key operation cannot produce an accepted approval.
Broader user-granted capabilities remain in the threat model's elevated tier.

## Future installed evidence

The separately authorized disposable harness must retain privacy-minimized, non-secret evidence.
Harness observations do not become a product Broker journal or authority store.

| Evidence record | Required content | Explicit exclusion |
| --- | --- | --- |
| Key-enrollment observation | Exact test build/source, host/OS/hardware, requested/profile/signed entitlements, access group/tag, requested and observable key attributes, public key/kid, algorithm-support result, private-export negative result, authorization-record digest, and normalized API results. | Private-key bytes, wrapped private representation, credentials, biometric data, passwords/passcodes, broad Keychain dumps. |
| Interaction observation | Fixture/payload/protected-header/`Sig_structure`/display digests, process and interaction generations, context creation/invalidation, reuse setting, lifecycle events, normalized LA/Security results, sign-call count, signature length/format, local verification, submission outcome, and final state. | Reusable context, private key, biometric material, unbounded localized OS prose. |
| Supervisor verification/commit receipt | Envelope/payload/`Sig_structure` digests, resolved public-key authorization, binding/time/profile results, canonical payload replay identity, `ApprovalID`, state, transaction outcome, and any later `AttemptID`. | Broker private-key reference, device identifier as authority, UI prose. |
| Installed experiment receipt | Exact matrix row, immutable target/profile/fixture identities, authorized mutation, environment, timestamps, normalized API/process results, cleanup, and oracle result. | Credentials, raw profiles containing secrets, Keychain exports, biometric or password content. |

At minimum, the matrix must cover:

| Group | Required rows and oracle |
| --- | --- |
| Key construction/readback | Correct enrollment; missing, duplicate, wrong tag/group/token/type/size/accessibility/ACL/permanence/synchronizability; private-export failure; unauthorized component; stale same-Team build with a retained group. Any mismatch refuses before prompt/sign. |
| Algorithm and exact bytes | Target support for the preferred RFC 4754 candidate; exactly 64-byte raw output; valid `r`/`s`; Broker same-byte verify; Supervisor verify; mutation of every protected/payload/`Sig_structure` binding; DER/short/oversize/alternate encodings refuse. |
| Fresh context and one sign | One approval succeeds; an instrumented negative control demonstrates why reused authorization is dangerous; production issues exactly one sign call; a second approval uses a different context generation; no silent authentication or signer fallback. |
| Concurrency | Identical requests join or receive `IN_FLIGHT`; different requests receive `BUSY`; one context, one sign budget, and one canonical payload remain active. |
| Cancellation/races | User cancel, app cancel, system cancel, timeout, lockout/unavailable, noninteractive refusal, invalidate/evaluate/sign callback races, and late callback after every terminal transition all produce no new authority. |
| UI/lifecycle | Focus loss, app deactivation, hide/background/Space change, screen lock, fast-user switch, sleep/wake, logout/login, and normal process termination require a fresh fetch/render/context/auth/sign. Synthetic UI alone cannot create an accepted signature; no secure-attention or capture-confidentiality claim is made. |
| Broker process death | Kill before authentication, after authentication, during sign, and after signature/local verify but before Supervisor commit. All have no durable approval authority and require a fresh interaction. |
| Supervisor commit and response loss | Lose the `SubmitApprovalV0` reply before and after commit, restart either process, and retry exact/equivalent envelopes. Supervisor truth converges to either no record or the same `ApprovalID`/state; Broker byte-identical recovery is not claimed. |
| Attempt transaction | Crash before and after atomic consume/create and lose the `RequestAttemptV0` reply. Replay returns no effect or the same `AttemptID`, never a second attempt. |
| Biometry/credential/accessibility | Enrollment changes under exact `userPresence`; lock/unlock; credential change; no-biometry variants; candidate passcode/accessibility behavior. Unexpected or ambiguous behavior leaves activation `BLOCKED`. |
| Update/rotation/restore | Same-authority update, identity-changing update, old-key refusal, deletion, repair, clean reinstall, and other-device restore. These rows are destructive or architecture-sensitive and require exact separate authorization. ADR-0021 remains Proposed. |

The authorization must enumerate destructive Keychain, credential, deletion, update, rotation,
restore, or rollback rows individually. Unnamed destructive rows do not run.

## Unsupported assumptions and stop rules

- A successful biometric/password event does not automatically authorize exactly one key use.
- Setting reuse duration to zero does not consume an authorized context.
- Focus loss, sleep, switch, or backgrounding does not guarantee pre-signature invalidation.
- Domain state is not a key ID, user identity, installation ID, epoch, or anti-rollback anchor.
- A Team ID, bundle ID, group label, Secure Enclave presence, or device identifier does not
  authorize a public key.
- Key nonexportability does not prove only the intended Broker may invoke it.
- A stable Keychain group does not exclude stale entitled same-Team code.
- Missing or corrupt key state cannot trigger lazy creation or silent replacement.
- DER and raw ECDSA encodings are not interchangeable.
- A second signature byte string over the same payload is not a second approval.
- Local Authentication does not prove anti-overlay, anti-recording, or Accessibility resistance.

An ambiguous, unsupported, unavailable, or mismatched result leaves the harness or activation
`BLOCKED`. It does not trigger a weaker accessibility class, alternate algorithm, software key,
broader group, new helper, or different authority owner.

## Planning and dependency result

The dependency split is now exact:

1. **R3 canonical passive evidence brief — `PASSED`.** This document retains the read-only Apple
   mechanism, candidate values, authority correction, limitations, and future evidence matrix.
2. **C6b1a unsigned Broker harness — `PASSED`.** Its immutable no-credential construction is
   pinned above; it is not a signed or installed target.
3. **C6b1b test-only Supervisor seam — `PASSED`.** Its immutable local model preserves the
   Supervisor as the sole durable authority owner; it is not a product store or consumer.
4. **C6b1c identity/profile and no-install signed-artifact readback — `PASSED`.** The immutable
   archive pin above retains exact profile metadata, signature/requirement/CDHash, hardened-runtime,
   and closed effective-entitlement evidence without installation or launch.
5. **C6b1d installed signing evidence — `BLOCKED`.** It requires a fresh exact owner authorization
   naming the retained C6b1c artifact/profile, owner account/container, Keychain and
   LocalAuthentication operations, prompt rules, destructive rows, evidence destination, and
   cleanup. C6b1c does not authorize any of those effects.
6. **Product Broker/`SubmitApprovalV0`/`RequestAttemptV0` wiring — `BLOCKED`.** It cannot begin
   until C6b1d and the installed authenticated service boundary pass their exact scopes.

A research document cannot satisfy installed evidence, product consumer evidence, or product
admission.

## C6b1d authorization packet inputs

This packet is preparation, not authorization. A later owner authorization must name all of these
immutable inputs together:

- Capsule commit `16fb810b97e7ff2a157a251ae4dc8023dcfc01b4` and `capsule-experiments`
  merge `82d1a799f70482856aaa6030f612d701b39cec67`;
- signed executable SHA-256
  `0a31663736678b0fccefb3f524209167aaed085b3c214cf8af2024a82ea38833`, CDHash
  `029b8d5cabd38e1fde9e23564e4e5b1590cf569d`, and CodeDirectory SHA-256
  `029b8d5cabd38e1fde9e23564e4e5b1590cf569dabc8bf1d307d7f80340c0431`;
- development profile `XT8MS38HWV`, UUID `2e8d338c-5668-4d41-9eb3-eb29634ebecf`, CMS SHA-256
  `a00dca2e4cfb8d4d432ffbeeaa0cc616e74aa8294364286f28dfe998ae0e32ee`, and certificate SHA-256
  `D3E9FBDDBC342F747C3649B5A6FFB307A575827404E02D638C11B6B795A09629`;
- exact bundle ID `com.capsulecorp.capsule.broker.c6b1`, Approval group
  `3DDR84M4JS.com.capsulecorp.capsule.broker.approval.epoch-7`, application tag
  `com.capsulecorp.capsule.broker.c6b1.approval-key.v0.generation-1`, and the owner-confirmed
  disposable account/container and evidence destination; and
- the permitted Keychain/LocalAuthentication calls, prompt behavior, stop conditions, and cleanup
  receipt. No broad Keychain inventory, private-key export, credential/biometric retention,
  product listener, runtime, backend, VM, or guest is allowed.

The proposed first bounded mutation set remains individually opt-in: D1 create exactly one primary
permanent Secure Enclave P-256 key; D2 create one same-tag duplicate and prove zero/multiple-match
refusal; D3 delete only that duplicate and prove exact absence; D4 delete the primary mid-run and
prove missing-key refusal with no lazy recreation; D14 delete exact experiment keys and prove zero
matches; D15 stop/remove only experiment app/process/service/container state; and D16 delete the
disposable account/home only after retained-evidence and path readback. D5-D13 and D17-D18 remain
excluded. Every included destructive row must be repeated verbatim in the eventual owner
authorization; omission means it does not run.

## ADR impact

No new ADR is required for this documentation/evidence integration. ADR-0043 remains Accepted and
governing. ADR-0021 remains Proposed.

Selecting any of the following may require an ADR-0043 refinement or a new ADR:

- a durable Broker journal, cache, recovery owner, or other Broker state authority;
- a weaker accessibility class or alternate credential/biometry ACL;
- an alternate signature encoding/conversion path;
- a new helper, mediator, privileged actor, or IPC route;
- changed Approval-key ownership or private-key access by the Supervisor/daemon;
- an App Group or broader shared group for the key;
- automatic key creation, self-authorization, reset, rotation, deletion, restoration, or rollback;
- identity-changing key/group epochs or changed old-key verification policy; or
- changed purpose, audience, installation, epoch, replay, consume, or attempt semantics.

Experimental evidence alone cannot accept ADR-0021 or change these boundaries.

## Limitations and next action

This research did not observe the exact target OS, hardware, signed entitlements/profile,
access-group behavior, key attributes, Secure Enclave operation, algorithm support/output,
authentication UI, focus/lifecycle races, process death, update, credential transition, or
authenticated consumer.

Confidence is high for the API and accepted authority composition, medium for the two experiment
candidates, and intentionally absent for installed behavior that was not run.

Next action: freeze the exact C6b1d owner authorization over the retained C6b1c bytes. The first
run may include only individually approved D1-D4 and cleanup D14-D16; D5-D13 and D17-D18 remain
deferred unless separately named. Until C6b1d and
the installed authenticated service boundary pass, live signing, product wiring, the installed
security boundary, and product admission remain `BLOCKED`.
