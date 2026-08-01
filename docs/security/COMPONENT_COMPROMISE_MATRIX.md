# Component Compromise Matrix

Status: normative intended design assumptions; controls are not yet implemented.

This matrix describes intended blast-radius reduction. It does not claim that merely separating
processes implements these protections; the macOS authority spike must validate the OS controls.

| Compromised component | Attacker can | Intended remaining protection | Required response |
| --- | --- | --- | --- |
| Agent client / MCP | Submit malicious proposals/source, request cancellation, exploit public parsing | Cannot select files, approve, launch, alter trust state, or retrieve user-only content | Reject/limit requests; no trust downgrade |
| Agent-facing daemon | Corrupt proposed plans/metadata, cause denial of service, lie in unsigned status, request malicious registered plans | Cannot sign user approval, directly launch backend, consume/reset grants, forge Supervisor transcript, or access Broker content by default | Broker/Supervisor reject mismatches; quarantine component on identity failure; repair/reinstall |
| Approval Broker | Misrender a plan, misuse user-presence UX, sign arbitrary approvals while its key is usable | Cannot directly launch guest or forge Supervisor enforcement claim; Supervisor still enforces registration, epoch, hard safety, expiry, and backend posture | Revoke Broker/key, invalidate pending grants, repair |
| Content Broker | Disclose or corrupt content within its accessible user-local domain | No agent/general network endpoint by design; cannot approve via evidence key or launch a guest | Revoke component, quarantine content, repair, assess disclosure |
| Approval key | Sign arbitrary objects allowed by its key authorization | Cannot launch; Supervisor still checks object type/purpose, registration, epoch, expected Supervisor, attempt, audience, expiry, and hard safety | Revoke/replace key; invalidate unused grants |
| Execution Supervisor | Launch unauthorized guests, misuse staged content, forge enforcement transcripts using its evidence key, hide cleanup state | Cannot forge Broker user-presence approval; optional external witness may expose later history rewriting | Critical local compromise; disable execution, quarantine artifacts, repair/reinstall and rotate evidence key |
| Supervisor evidence key | Forge Supervisor-attributable enforcement claims within its purpose | Cannot create a Broker approval or independently reach backend unless component access is also compromised | Revoke/replace key; distrust affected evidence interval |
| Updater/trust verifier | Withhold refresh, emit malicious local snapshot if its local signing authority is compromised | Pinned TUF roots and delegated metadata constrain acceptable external artifacts; Supervisor hard safety remains | Quarantine trust snapshot/component, refresh through trusted installer, rotate local authority |
| External trust service | Withhold or replay metadata, observe refresh traffic, serve malicious unsigned artifacts | Pinned TUF roots, threshold/delegated signatures, versions, expiration, hashes, and checkpoints | Fail refresh; use allowed verified cache or refuse/downgrade |
| TUF online role key | Sign malicious metadata within delegated role | Offline/threshold root and delegated path/role scope constrain damage | Root-led revoke/rotate; publish corrected higher versions |
| Runtime publisher | Publish malicious bundle | Independent review, registry activation, backend validation, and hard safety remain required | Revoke bundle/publisher delegation; deactivate exact digest |
| Profile reviewer | Approve a malicious exact bundle | Cannot change signed bundle bytes or install component binaries | Remove review authority; deactivate affected entries; rerun validation |
| Isolation backend | Escape or falsify control/teardown claims | Supervisor-side staging/identity/evidence checks and external host controls may limit some effects; backend boundary itself is lost | Critical/high depending on reach; disable exact configuration and revoke validation record |
| Optional Guardian | Suppress or fabricate optional observations | Never grants authority; core point-in-time checks remain | Mark evidence unavailable/degraded; repair Guardian separately |
| Same-user malicious process | Probe local services/files, impersonate by name/path/PID, exploit shared state | XPC code requirements, protected storage, Keychain groups, entitlements, exact epoch/build checks if proven | Deny connection/access; quarantine on enrolled-state modification |
| Stale same-team component | Retain its historical signing identifier/profile/access-group entitlement and attempt to use replacement group keys | Exact-build/epoch XPC denies trusted channels; a fresh security-epoch access group/key denies old/new private-key cross-use when the transition completes | Fence execution, finish or repair the authorized security-epoch transition, retire the old key, and never place replacement authority in the stale group |
| Local administrator / kernel | Can likely control components, keys in use, storage, backend, and local observations | Optional independent checkpoint may reveal historical inconsistency; no complete local containment | Explicitly out of validated-local scope; reinstall/re-enroll and distrust affected interval |

## Severity anchors

- Supervisor compromise is Critical because the sole execution authority and its evidence claims are
  lost.
- Agent or daemon compromise becomes Critical if it can reach approval keys, user-only content, or
  backend launch despite the intended split.
- Broker compromise is serious but should not independently create a guest or forge Supervisor
  evidence.
- External trust compromise is contained only to the degree pinned roles, delegated scope, local
  checkpoints, and independent review remain intact.

## Validation use

Each “remaining protection” must map to a mechanism and attack test in
[Control Evidence Matrix](CONTROL_EVIDENCE_MATRIX.md). Until then it remains `proposed`, not an
implemented property.
