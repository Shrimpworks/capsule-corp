import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";
import { verifyFieldAuthorityManifest } from "./verify-field-authority-manifest.mjs";

const repositoryRoot = new URL("../", import.meta.url);
const manifestUrl = new URL("../schemas/authority/field-authority-manifest.json", import.meta.url);
const checkedInManifest = JSON.parse(await readFile(manifestUrl, "utf8"));

test("verifies the checked-in passive field-authority manifest", async () => {
  assert.deepEqual(await verifyFieldAuthorityManifest({ rootDirectory: repositoryRoot }), {
    fieldCount: 164,
    profileCount: 30,
    targetCount: 15,
  });
});

test("rejects a canonical field missing from the manifest", async () => {
  const manifest = cloneManifest();
  target(manifest, "capsule.execution-plan").fields.pop();
  await assert.rejects(
    verifyFieldAuthorityManifest({ manifest, rootDirectory: repositoryRoot }),
    /missing field classifications for capsule\.execution-plan@0: 24/u,
  );
});

test("rejects an unknown closed authority class", async () => {
  const manifest = cloneManifest();
  manifest.profiles["plan-runtime"].authorityClass = "runtime-ish-authority";
  await assert.rejects(
    verifyFieldAuthorityManifest({ manifest, rootDirectory: repositoryRoot }),
    /field-authority manifest schema failure:.*authorityClass/u,
  );
});

test("rejects duplicate field paths", async () => {
  const manifest = cloneManifest();
  const fields = target(manifest, "capsule.plan-registration").fields;
  fields[1].path = fields[0].path;
  await assert.rejects(
    verifyFieldAuthorityManifest({ manifest, rootDirectory: repositoryRoot }),
    /duplicate classified path for capsule\.plan-registration@0/u,
  );
});

test("rejects a stale manifest object version", async () => {
  const manifest = cloneManifest();
  target(manifest, "capsule.approval-grant").version = 1;
  await assert.rejects(
    verifyFieldAuthorityManifest({ manifest, rootDirectory: repositoryRoot }),
    /stale object version or identity for capsule\.approval-grant@1/u,
  );
});

test("rejects a classified field absent from the canonical target", async () => {
  const manifest = cloneManifest();
  target(manifest, "capsule.typescript-plan-source-bindings").fields.push({
    path: "/futureAuthority",
    sourceField: 99,
    profile: "typescript-runtime",
  });
  await assert.rejects(
    verifyFieldAuthorityManifest({ manifest, rootDirectory: repositoryRoot }),
    /classified fields absent from canonical target capsule\.typescript-plan-source-bindings@0: 99/u,
  );
});

function cloneManifest() {
  return structuredClone(checkedInManifest);
}

function target(manifest, object) {
  const value = manifest.targets.find((candidate) => candidate.object === object);
  assert.ok(value, `missing fixture target ${object}`);
  return value;
}
