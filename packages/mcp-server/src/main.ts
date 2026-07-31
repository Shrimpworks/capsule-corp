#!/usr/bin/env node

import { CAPSULE_TOOL_CATALOG } from "./index.js";

if (process.argv.includes("--describe")) {
  process.stdout.write(`${JSON.stringify({ tools: CAPSULE_TOOL_CATALOG }, null, 2)}\n`);
} else {
  process.stderr.write(
    "capsule-mcp is scaffolded but not connected to an MCP transport yet; use --describe\n",
  );
  process.exitCode = 64;
}
