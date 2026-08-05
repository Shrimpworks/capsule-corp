import assert from "node:assert/strict";
import { mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { verifyFieldAuthorityManifest } from "./verify-field-authority-manifest.mjs";

const repositoryRoot = new URL("../", import.meta.url);
const manifestUrl = new URL("../schemas/authority/field-authority-manifest.json", import.meta.url);
const checkedInManifest = JSON.parse(await readFile(manifestUrl, "utf8"));
const schemaUrl = new URL(
  "../schemas/authority/field-authority-manifest.schema.json",
  import.meta.url,
);
const checkedInSchema = JSON.parse(await readFile(schemaUrl, "utf8"));

test("verifies the checked-in passive field-authority manifest", async () => {
  assert.deepEqual(await verifyFieldAuthorityManifest({ rootDirectory: repositoryRoot }), {
    fieldCount: 1203,
    profileCount: 95,
    targetCount: 60,
  });
});

test("rejects a missing recursive C2B passive field classification", async () => {
  const manifest = cloneManifest();
  const fields = target(manifest, "capsule.governed-deno-core-c2b-passive-binding").fields;
  fields.splice(
    fields.findIndex((field) => field.path === "/fixedFixture/source/sha256"),
    1,
  );
  await assert.rejects(
    verifyFieldAuthorityManifest({ manifest, rootDirectory: repositoryRoot }),
    /missing field classifications.*fixedFixture\/source\/sha256/u,
  );
});

test("rejects an incomplete C2B passive consumer classification", async () => {
  const manifest = cloneManifest();
  manifest.profiles["governed-runtime-c2b-passive-evidence"].allowedConsumers.pop();
  await assert.rejects(
    verifyFieldAuthorityManifest({ manifest, rootDirectory: repositoryRoot }),
    /incomplete retention or consumer classification/u,
  );
});

test("rejects an incomplete C2B v2 passive consumer classification", async () => {
  const manifest = cloneManifest();
  manifest.profiles["governed-runtime-c2b-v2-passive-evidence"].allowedConsumers.pop();
  await assert.rejects(
    verifyFieldAuthorityManifest({ manifest, rootDirectory: repositoryRoot }),
    /incomplete retention or consumer classification/u,
  );
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

test("rejects a missing nested SourceManifest member classification", async () => {
  const manifest = cloneManifest();
  target(manifest, "capsule.source-manifest").fields.splice(5, 1);
  await assert.rejects(
    verifyFieldAuthorityManifest({ manifest, rootDirectory: repositoryRoot }),
    /missing field classifications for capsule\.source-manifest@0: 4\/\*\/2/u,
  );
});

test("classifies an embedded Go struct field with no separate field name", async () => {
  const root = await mkdtemp(join(tmpdir(), "field-authority-go-"));
  try {
    await mkdir(join(root, "internal", "fixture"), { recursive: true });
    await writeFile(
      join(root, "internal", "fixture", "types.go"),
      [
        "package fixture",
        "",
        "// field-authority-object: capsule.embedded-fixture v0",
        "type Fixture struct {",
        "\tEmbeddedRole",
        "\tName string",
        "}",
        "",
      ].join("\n"),
    );
    const manifest = minimalManifest({
      object: "capsule.embedded-fixture",
      version: 0,
      definition: { kind: "go-struct", path: "internal/fixture/types.go", type: "Fixture" },
      fields: [
        { path: "/embeddedRole", sourceField: "EmbeddedRole", profile: "supervisor-reference" },
        { path: "/name", sourceField: "Name", profile: "supervisor-reference" },
      ],
    });
    const result = await verifyFieldAuthorityManifest({
      manifest,
      schema: checkedInSchema,
      rootDirectory: root,
    });
    assert.equal(result.fieldCount, 2);

    // Proves the embedded field is actually required, not merely tolerated:
    // dropping its classification must now be caught as missing coverage.
    const withoutEmbedded = structuredClone(manifest);
    withoutEmbedded.targets[0].fields.shift();
    await assert.rejects(
      verifyFieldAuthorityManifest({
        manifest: withoutEmbedded,
        schema: checkedInSchema,
        rootDirectory: root,
      }),
      /missing field classifications for capsule\.embedded-fixture@0: EmbeddedRole/u,
    );
  } finally {
    await rm(root, { recursive: true, force: true });
  }
});

test("finds the closing brace of a CDDL rule with a brace inside a string and a comment", async () => {
  const root = await mkdtemp(join(tmpdir(), "field-authority-cddl-"));
  try {
    await mkdir(join(root, "schemas", "fixture"), { recursive: true });
    await writeFile(
      join(root, "schemas", "fixture", "rule.cddl"),
      [
        "fixture-rule = {",
        '  1: "capsule.embedded-fixture" ; a trailing comment with a stray } brace',
        "  2: 0,",
        '  3: "a text value containing a { brace",',
        "}",
        "",
      ].join("\n"),
    );
    const manifest = minimalManifest({
      object: "capsule.embedded-fixture",
      version: 0,
      definition: { kind: "cddl-map", path: "schemas/fixture/rule.cddl", rule: "fixture-rule" },
      fields: [
        { path: "/1", sourceField: "1", profile: "supervisor-reference" },
        { path: "/2", sourceField: "2", profile: "supervisor-reference" },
        { path: "/3", sourceField: "3", profile: "supervisor-reference" },
      ],
    });
    const result = await verifyFieldAuthorityManifest({
      manifest,
      schema: checkedInSchema,
      rootDirectory: root,
    });
    assert.equal(result.fieldCount, 3);
  } finally {
    await rm(root, { recursive: true, force: true });
  }
});

function cloneManifest() {
  return structuredClone(checkedInManifest);
}

function minimalManifest(soleTarget) {
  return {
    manifestVersion: checkedInManifest.manifestVersion,
    status: checkedInManifest.status,
    vocabularyVersion: checkedInManifest.vocabularyVersion,
    profiles: { "supervisor-reference": checkedInManifest.profiles["supervisor-reference"] },
    targets: [soleTarget],
  };
}

function target(manifest, object) {
  const value = manifest.targets.find((candidate) => candidate.object === object);
  assert.ok(value, `missing fixture target ${object}`);
  return value;
}
