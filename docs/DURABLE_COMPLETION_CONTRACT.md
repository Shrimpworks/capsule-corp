# Passive durable completion-last contract

## Status and authorized scope

Work item: Supervisor-owned durable completion producer/store transaction, bounded transcript, and
fixed agent summary for repository lifecycle/FakeBackend fixtures.

Status: `PASSED`

Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`

Product completion and admission: `BLOCKED`

This slice defensively validates terminal-proof and replay controls using only
`internal/execution/completioncomposer`, retained lifecycle fixtures, `FakeBackend` identity, owned
temporary files, and controlled local commit faults. It creates no process, runtime, adapter, VM,
guest, signing key, IPC endpoint, product store, or product authority.

## Commit-last rule

`Producer.Complete` copies one version-0 request containing only `AttemptID`, exact result JSON, and
the fixed `supervisor-validated-typed-json-sha256` integrity disposition. There is no EOF, exit
status, diagnostic, path, artifact name, timing, backend flag, or caller-selected lifecycle field.

Inside the store transaction, the producer:

1. resolves the immutable created attempt and retained exact plan from Supervisor/store truth;
2. resolves the lifecycle record by the same `AttemptID`;
3. revalidates approval/attempt/registration/plan/source/runtime/profile/backend cross-links;
4. requires the exact completed six-operation FakeBackend lifecycle, destroy confirmation, cleanup
   false, and fake authoritative absence;
5. validates typed JSON, exact length, SHA-256, the fixed result cap, and the approved plan cap;
6. composes one deterministic bounded transcript and fixed summary; and
7. atomically publishes the full record with `committed-last` before returning either object.

An incomplete lifecycle, unresolved cleanup, missing fact, EOF, or exit-zero observation produces no
completion record. An indeterminate rename or directory sync fences the current store instance;
reopen must establish one complete predecessor or successor before retry.

## Exact objects and caps

| Object/budget | Exact value |
| --- | ---: |
| commit request version | 0 |
| `CompletionRecord` version | 0 |
| `capsule.unwired-fake-durable-completion` version | 0 |
| `capsule.unwired-fake-completion-store` version | 0 |
| result JSON bytes | 262,144 |
| transcript bytes | 4,096 |
| public summary bytes | 256 |
| encoded durable completion bytes | 368,640 |
| retained completion records | 4,096 |
| aggregate retained result JSON bytes | 67,108,864 |
| fixed completion-store bytes | 100,663,296 |

The store header binds one nonzero installation ID and every record must match it. The transcript
binds attempt, approval, registration, plan, installation, epoch, Supervisor,
approval authorization, source manifest, runtime bundle, profile registry entry, backend validation,
backend configuration/implementation/instance, immutable lifecycle bindings, completion and result
digests, result integrity, lifecycle sequence, teardown, absence, cleanup, fake-local terminal class,
and the explicit unresolved fake runner identity.

The summary remains only:

```json
{"state":"completed","attemptId":"attempt_<fixed hex>","transcriptId":"transcript_sha256_<fixed hex>"}
```

It contains no result content, guest-controlled string, path, artifact name/size, diagnostic,
detailed violation, timing, or metric. `completed` means only that the no-guest local transaction
completed; the transcript states the limitation and no product consumer is wired.

## Replay, recovery, and refusal

The first completion for an `AttemptID` is immutable. Exact retries, response loss, and reopen
return byte-identical transcript and summary. Changed result bytes are `REPLAY`; they do not mutate
the record. Store validation rejects missing/early/duplicate commit markers, duplicate attempts,
malformed or mixed versions, cap overflow, invalid base64/hex, forged cross-links, and changed
digests as `RECOVERY_REQUIRED` without rewriting evidence.

Focused tests retain:

- completion and durable-record known answers;
- exact result cap plus cap+1, malformed JSON, and mixed request/store versions;
- forged, missing, early, and duplicate completion records;
- EOF/runner-exit-zero-only, nonterminal, unresolved teardown, and missing completion cases;
- confirmed pre-rename failure, post-rename indeterminate response loss, reopen, restart, exact
  replay, and stale replay;
- defensive copies for request, result, transcript, summary, and stored views.

## Explicit limits and next action

`RunnerIdentityUnresolvedFake` is intentional: FakeBackend creates no runner or guest. Fake
destroy/absence proves only the repository simulator's state, not process-tree teardown. Typed JSON
integrity proves bounded syntax/bytes/digest, not a correct workload, uncompromised guest kernel, or
trusted launcher report. The transcript is unsigned and not a receipt or attestation.

The fixed completion file is separate only to keep this conformance slice narrow. It requires an
external sole owner and is not the selected product engine. Product work remains `BLOCKED` on the
real completion transport and trusted launcher, result-integrity owner, materialized runner/profile,
authoritative process-tree absence, installed protected owner/store, completion transaction in the
selected Supervisor engine, evidence signing, Broker release, and composed receipt.
