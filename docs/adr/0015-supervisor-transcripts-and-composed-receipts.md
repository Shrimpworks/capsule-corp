# ADR-0015: Supervisor transcripts and composed execution receipts

- Status: Accepted
- Date: 2026-07-31

## Context

A receipt signed only by the agent-facing daemon lets that daemon author both the plan narrative and
the claim that it executed correctly. It does not preserve the new authority separation.

## Decision

The Supervisor records bounded hash-linked enforcement events and signs a terminal
`EnforcementTranscript` with a purpose-limited evidence key. The Broker supplies the separate
`ApprovalGrant`. A user `ExecutionReceipt` composes the registered plan, signer authorizations,
approval, integrity assessment summary, exact backend/profile/control identities, Supervisor
transcript, artifact manifest, result classification, teardown state, and optional witness proof.

The daemon may package/index the receipt but cannot forge either embedded authority claim. The
agent receives a separate fixed `AgentExecutionSummary`, not the full receipt by default.

Evidence language is calibrated as “cryptographically attributable approval claim,”
“cryptographically attributable enforcement transcript,” and “signature-verifiable composed
receipt.” It is not called independent platform attestation.

## Consequences

- Daemon-forged ordinary success is independently detectable.
- Verifiers must validate embedded object types, purposes, key authorizations, bindings, and status,
  not only an outer package signature.
- Supervisor compromise can forge its own enforcement claims and remains Critical.
- Receipts retain multidimensional posture and known limitations.
- Agent metadata channels are deliberately smaller than user evidence.
