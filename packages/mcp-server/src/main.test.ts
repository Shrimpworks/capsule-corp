import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import test from "node:test";
import { fileURLToPath } from "node:url";

import { CAPSULE_TOOL_CATALOG } from "./index.js";

const mainPath = fileURLToPath(new URL("./main.js", import.meta.url));

test("--describe prints the tool catalog as JSON and exits 0", () => {
  const result = spawnSync(process.execPath, [mainPath, "--describe"], { encoding: "utf8" });

  assert.equal(result.status, 0);
  assert.equal(result.stderr, "");
  const parsed = JSON.parse(result.stdout) as { tools: unknown };
  assert.deepEqual(parsed, { tools: CAPSULE_TOOL_CATALOG });
});

test("without --describe, reports the unconnected-transport fallback and exits 64", () => {
  const result = spawnSync(process.execPath, [mainPath], { encoding: "utf8" });

  assert.equal(result.status, 64);
  assert.equal(result.stdout, "");
  assert.match(
    result.stderr,
    /capsule-mcp is scaffolded but not connected to an MCP transport yet; use --describe/,
  );
});
