# Candidate schema status

Schemas in this directory are passive Phase 2 contract candidates. They are verified in CI but are
not activated by the daemon, SDK, MCP server, Broker, or Supervisor.

The top-level schemas remain canonical for the current pre-freeze scaffold until the complete
proposal, plan, registration, and consumer migration can cut over atomically. Candidate schemas
must not be used as permissive adapters for the mixed `Job` shape.

[`job-proposal-v0.schema.json`](job-proposal-v0.schema.json) deliberately covers only the first
dependency-free inline-JSON proposal wedge. It structurally omits unsupported authority and fixes
the first logical input/output roles. Raw JSON limits, duplicate-key rejection, semantic source
validation, trusted default/ceiling resolution, and internal canonical plan construction remain
separate boundaries.
