import assert from "node:assert/strict";
import test from "node:test";

import { CAPSULE_TOOL_CATALOG } from "./index.js";

// The only tool in the catalog whose description promises guest-controlled
// artifact content: every other entry either grants no content access or
// explicitly documents that it returns metadata/receipts "without returning
// content" / "without guest-controlled artifact content". Matching by name
// here (rather than the no-op `tool.name.includes("read_artifact")`, which
// matches no real catalog entry) means this test actually fails if
// `capsule.request_artifact_access` is ever misconfigured, renamed away from
// "approval-required", or removed from the catalog.
const CONTENT_BEARING_TOOL_NAMES = ["capsule.request_artifact_access"];

test("no MCP tool exposes artifact content without a separate approval path", () => {
  const unsafeTools = CAPSULE_TOOL_CATALOG.filter(
    (tool) =>
      CONTENT_BEARING_TOOL_NAMES.includes(tool.name) && tool.contentAccess !== "approval-required",
  );

  assert.deepEqual(unsafeTools, []);
  assert.ok(CAPSULE_TOOL_CATALOG.some((tool) => tool.name === "capsule.request_artifact_access"));

  for (const name of CONTENT_BEARING_TOOL_NAMES) {
    const tool = CAPSULE_TOOL_CATALOG.find((entry) => entry.name === name);
    assert.ok(tool, `expected content-bearing tool "${name}" to exist in the catalog`);
    assert.equal(tool?.contentAccess, "approval-required", name);
  }

  // Every non-content-bearing tool must NOT claim approval-required content
  // access either, so a future entry can't silently widen this guard's
  // blind spot by being both content-bearing and mis-tagged as safe, nor
  // can an unrelated tool be over-tagged in a way that masks a real gap.
  for (const tool of CAPSULE_TOOL_CATALOG) {
    if (CONTENT_BEARING_TOOL_NAMES.includes(tool.name)) {
      continue;
    }
    assert.notEqual(tool.contentAccess, "approval-required", tool.name);
  }
});
