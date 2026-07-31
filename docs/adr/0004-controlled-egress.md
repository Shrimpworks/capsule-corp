# ADR-0004: Guest output is controlled egress

- Status: Accepted
- Date: 2026-07-30

## Context

A guest with no network can still disclose all input through stdout, errors, structured results, or
artifacts if those channels are automatically returned to an AI client.

## Decision

All guest-controlled output passes through an artifact/egress broker. Metadata and content are
separate permissions. Every artifact declares validation and audience policy.

The default for user-granted data is full delivery to the user and metadata-only delivery to the
agent. Agent content access requires a separate grant.

## Consequences

- MCP tools expose descriptors before content.
- stdout and stderr are bounded and not automatically forwarded in full.
- Structured guest results are not inherently trusted.
- User interfaces must distinguish task completion from agent-readable output.
