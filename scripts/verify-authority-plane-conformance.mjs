import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { sha256Hex } from "./lib/fixture-bytes.mjs";

const root = new URL("../schemas/conformance/authority-plane-v0/", import.meta.url);
const manifest = JSON.parse(await readFile(new URL("manifest.json", root), "utf8"));
assert.equal(manifest.objectType, "capsule.authority-plane-conformance");
assert.equal(manifest.objectVersion, 0);
assert.deepEqual(manifest.caps, {
  executionPlan: 65_536,
  roleBindings: 562,
  sourceManifest: 95,
  source: 262_144,
  planRegistration: 4_096,
  registerPlanV0: 328_337,
  getRegisteredPlanV0: 332_433,
});
for (const [path, known] of Object.entries(manifest.fixtures)) {
  const bytes = await readFile(new URL(path, root));
  assert.equal(bytes.length, known.byteLength, path);
  assert.equal(sha256Hex(bytes), known.sha256, path);
}
const bindings = await readFile(new URL("role-bindings.bin", root));
assert.equal(bindings.length, 562);
assert.equal(bindings[0], 0);
assert.equal(bindings[145], 2);
assert.equal(
  bindings.subarray(146 + 2 * 32, 146 + 8 * 32).every((byte) => byte === 0),
  true,
);
const source = await readFile(new URL("main.mjs", root));
assert.equal(source[0] === 0xef && source[1] === 0xbb && source[2] === 0xbf, false);
new TextDecoder("utf-8", { fatal: true, ignoreBOM: true }).decode(source);
process.stdout.write("verified independent TypeScript authority-plane known answers\n");
