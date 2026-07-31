import assert from "node:assert/strict";
import test from "node:test";

import { CAPSULE_TOOL_CATALOG } from "./index.js";

test("no MCP tool exposes artifact content without a separate approval path", () => {
  const unsafeTools = CAPSULE_TOOL_CATALOG.filter(
    (tool) => tool.contentAccess !== "approval-required" && tool.name.includes("read_artifact"),
  );

  assert.deepEqual(unsafeTools, []);
  assert.ok(CAPSULE_TOOL_CATALOG.some((tool) => tool.name === "capsule.request_artifact_access"));
});
