// verify-schemas.mjs reads schemas/, examples/, and profiles/ relative to its
// own file location and hardcodes the fixture paths it validates, so it cannot
// be pointed at a synthetic tree the way verify-adr-index.mjs can. These tests
// therefore run it against the real repository and assert the property that
// matters: what the script reports having validated must equal what is
// actually on disk, so a narrowing of its discovery fails the run.
import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { readdir } from "node:fs/promises";
import { sep } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);
const scriptPath = fileURLToPath(new URL("./verify-schemas.mjs", import.meta.url));
const schemaDirectory = new URL("../schemas/", import.meta.url);

test("validates the checked-in schemas cleanly", async () => {
  const { stdout } = await execFileAsync(process.execPath, [scriptPath]);
  assert.match(stdout, /^validated schemas\/job\.schema\.json$/mu);
});

test("validates every schema in the tree, not only the top level", async () => {
  // Regression: a non-recursive readdir once found only schemas/*.schema.json
  // plus one hardcoded candidate path, so the schemas under schemas/authority/
  // and schemas/conformance/ silently escaped the draft and $id checks. The
  // script reported success either way, so only comparing against an
  // independent walk catches it.
  const { stdout } = await execFileAsync(process.execPath, [scriptPath]);
  const reported = [...stdout.matchAll(/^validated schemas\/(.+\.schema\.json)$/gmu)]
    .map((match) => match[1])
    .sort();

  const onDisk = (await readdir(schemaDirectory, { recursive: true }))
    .map((entry) => entry.split(sep).join("/"))
    .filter((entry) => entry.endsWith(".schema.json"))
    .sort();

  assert.deepEqual(reported, onDisk);
  assert.ok(
    onDisk.some((entry) => entry.includes("/")),
    "expected at least one schema in a subdirectory, otherwise this test proves nothing",
  );
});

test("every validated schema declares a distinct $id", async () => {
  // The schemas all share one strict Ajv instance, so a duplicate $id would
  // make addSchema throw rather than silently shadow. Assert the property here
  // too, so the failure is legible rather than an Ajv internal error.
  const files = (await readdir(schemaDirectory, { recursive: true }))
    .map((entry) => entry.split(sep).join("/"))
    .filter((entry) => entry.endsWith(".schema.json"));

  const identifiers = new Map();
  for (const file of files) {
    const { default: schema } = await import(new URL(file, schemaDirectory).href, {
      with: { type: "json" },
    });
    assert.equal(
      schema.$schema,
      "https://json-schema.org/draft/2020-12/schema",
      `${file} must declare JSON Schema draft 2020-12`,
    );
    assert.ok(
      typeof schema.$id === "string" && schema.$id.length > 0,
      `${file} must declare a stable $id`,
    );
    assert.ok(
      !identifiers.has(schema.$id),
      `${file} reuses $id from ${identifiers.get(schema.$id)}`,
    );
    identifiers.set(schema.$id, file);
  }
});
