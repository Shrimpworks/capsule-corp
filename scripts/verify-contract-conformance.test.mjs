import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { cp, mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { verifyConformanceCorpus } from "./verify-contract-conformance.mjs";

const schemaPath = new URL("../schemas/conformance/v0/manifest.schema.json", import.meta.url);
const checkedInCorpus = new URL("../schemas/conformance/v0/", import.meta.url);

test("verifies the checked-in foundational conformance corpus", async () => {
  const result = await verifyConformanceCorpus({ rootDirectory: checkedInCorpus });

  assert.deepEqual(result, { caseCount: 105, fixtureCount: 91, ruleCount: 37 });
});

test("retains exact JSON boundary values and their cap-plus-one pairs", async () => {
  assert.equal((await fixtureBytes("json-raw-bytes-exact.bin")).length, 2_097_152);
  assert.equal((await fixtureBytes("json-raw-bytes-over.bin")).length, 2_097_153);

  const depthExact = measureJson(await fixtureJson("json-depth-exact.bin"));
  const depthOver = measureJson(await fixtureJson("json-depth-over.bin"));
  assert.equal(depthExact.maxDepth, 32);
  assert.equal(depthOver.maxDepth, 33);

  const nodesExact = measureJson(await fixtureJson("json-nodes-exact.bin"));
  const nodesOver = measureJson(await fixtureJson("json-nodes-over.bin"));
  assert.deepEqual(
    [nodesExact.nodes, nodesExact.totalMembers, nodesExact.totalElements],
    [8_193, 4_096, 4_096],
  );
  assert.deepEqual(
    [nodesOver.nodes, nodesOver.totalMembers, nodesOver.totalElements],
    [8_194, 4_096, 4_097],
  );

  await assertJsonMetric("json-total-members", "totalMembers", 4_096, 4_097);
  await assertJsonMetric("json-object-members", "maxObjectMembers", 256, 257);
  await assertJsonMetric("json-total-elements", "totalElements", 4_096, 4_097);
  await assertJsonMetric("json-array-elements", "maxArrayElements", 256, 257);
  await assertJsonMetric("json-decoded-text", "decodedTextBytes", 1_572_864, 1_572_865);
  await assertJsonMetric("json-key-bytes", "maxKeyBytes", 128, 129);
  await assertJsonMetric("json-string-bytes", "maxStringBytes", 65_536, 65_537);
});

test("retains exact CBOR predecoder boundaries for both candidate objects", async () => {
  for (const [slug, limits] of [
    [
      "execution-plan",
      { rawBytes: 65_536, depth: 8, items: 256, mapEntries: 64, arrayElements: 8 },
    ],
    [
      "plan-registration",
      { rawBytes: 4_096, depth: 4, items: 33, mapEntries: 16, arrayElements: 0 },
    ],
  ]) {
    assert.equal((await fixtureBytes(`cbor-${slug}-raw-bytes-exact.bin`)).length, limits.rawBytes);
    assert.equal(
      (await fixtureBytes(`cbor-${slug}-raw-bytes-over.bin`)).length,
      limits.rawBytes + 1,
    );
    await assertCborMetric(slug, "depth", "maxDepth", limits.depth, limits.depth + 1);
    await assertCborMetric(slug, "items", "items", limits.items, limits.items + 1);
    await assertCborMetric(
      slug,
      "map-entries",
      "maxMapEntries",
      limits.mapEntries,
      limits.mapEntries + 1,
    );
    await assertCborMetric(
      slug,
      "array-elements",
      "maxArrayElements",
      limits.arrayElements,
      limits.arrayElements + 1,
    );
  }
});

test("retains representative noncanonical CBOR bytes without decoding and re-encoding", async () => {
  const expectedHex = new Map([
    ["cbor-profile-canonical-map.bin", "a10100"],
    ["cbor-profile-indefinite-container.bin", "9f00ff"],
    ["cbor-profile-float.bin", "f93c00"],
    ["cbor-profile-bignum.bin", "c24101"],
    ["cbor-profile-tag.bin", "c000"],
    ["cbor-profile-duplicate-key.bin", "a201000101"],
    ["cbor-profile-trailing-data.bin", "0000"],
    ["cbor-profile-nonpreferred-integer.bin", "1800"],
    ["cbor-profile-invalid-utf8.bin", "61ff"],
    ["cbor-profile-map-order.bin", "a202000100"],
  ]);
  for (const [name, hex] of expectedHex) {
    assert.equal((await fixtureBytes(name)).toString("hex"), hex, name);
  }
});

test("verifies a closed, complete, byte-exact conformance corpus", async (t) => {
  const corpus = await createCorpus(t);

  const result = await verifyConformanceCorpus({ rootDirectory: corpus.root });

  assert.deepEqual(result, { caseCount: 2, fixtureCount: 2, ruleCount: 1 });
});

test("rejects an unknown manifest field", async (t) => {
  const corpus = await createCorpus(t, (manifest) => {
    manifest.unreviewed = true;
  });

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /manifest schema validation failed/u,
  );
});

test("rejects a listed fixture that is missing", async (t) => {
  const corpus = await createCorpus(t);
  await rm(join(corpus.root, "shared/reject.bin"));

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /listed fixture does not exist/u,
  );
});

test("rejects a fixture that is present but unlisted", async (t) => {
  const corpus = await createCorpus(t);
  await writeFile(join(corpus.root, "shared/unlisted.bin"), "unlisted");

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /fixture is not listed/u,
  );
});

test("rejects fixture bytes that do not match the retained hash", async (t) => {
  const corpus = await createCorpus(t);
  await writeFile(join(corpus.root, "shared/accept.bin"), "mutated");

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /fixture byte length mismatch|fixture SHA-256 mismatch/u,
  );
});

test("rejects deletion of a required rule outcome and its fixture", async (t) => {
  const corpus = await createCorpus(t, (manifest) => {
    manifest.cases = manifest.cases.filter((entry) => entry.expected.decision !== "reject");
  });
  await rm(join(corpus.root, "shared/reject.bin"));

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /required case is not covered/u,
  );
});

async function createCorpus(t, mutateManifest = () => {}) {
  const root = await mkdtemp(join(tmpdir(), "capsule-conformance-"));
  t.after(async () => rm(root, { recursive: true, force: true }));
  await mkdir(join(root, "shared"));
  await cp(schemaPath, join(root, "manifest.schema.json"));

  const acceptBytes = Buffer.from("accepted", "utf8");
  const rejectBytes = Buffer.from("rejected", "utf8");
  await writeFile(join(root, "shared/accept.bin"), acceptBytes);
  await writeFile(join(root, "shared/reject.bin"), rejectBytes);

  const implementations = {
    "fixture-integrity": "verified",
    go: "pending",
    typescript: "pending",
    swift: "not-applicable",
  };
  const manifest = {
    manifestVersion: "capsule.conformance/v0",
    rules: [
      {
        id: "job-proposal.raw.example",
        source: "ADR-0023#strict-jobproposal-raw-profile",
        description: "Test-only representative rule.",
        requiredCases: [
          { decision: "accept", variant: "ordinary" },
          { decision: "reject", variant: "malformed" },
        ],
      },
    ],
    cases: [
      {
        id: "job-proposal.raw.example.accept",
        description: "Accept representative bytes.",
        ruleIds: ["job-proposal.raw.example"],
        object: "JobProposal",
        wireFormat: "raw-bytes",
        mediaType: "application/capsule.job-proposal+json;v=0",
        variant: "ordinary",
        fixture: reference("shared/accept.bin", acceptBytes),
        context: { kind: "none" },
        expected: {
          decision: "accept",
          classification: null,
          owner: "public-raw-decoder",
          authorityStateChanged: false,
        },
        implementations,
      },
      {
        id: "job-proposal.raw.example.reject",
        description: "Reject representative bytes.",
        ruleIds: ["job-proposal.raw.example"],
        object: "JobProposal",
        wireFormat: "raw-bytes",
        mediaType: "application/capsule.job-proposal+json;v=0",
        variant: "malformed",
        fixture: reference("shared/reject.bin", rejectBytes),
        context: { kind: "none" },
        expected: {
          decision: "reject",
          classification: "MALFORMED",
          owner: "public-raw-decoder",
          authorityStateChanged: false,
        },
        implementations,
      },
    ],
  };
  mutateManifest(manifest);
  await writeFile(join(root, "manifest.json"), `${JSON.stringify(manifest, null, 2)}\n`);

  return { root };
}

function reference(path, bytes) {
  return {
    path,
    sha256: createHash("sha256").update(bytes).digest("hex"),
    byteLength: bytes.length,
  };
}

async function assertJsonMetric(prefix, property, exact, over) {
  assert.equal(measureJson(await fixtureJson(`${prefix}-exact.bin`))[property], exact);
  assert.equal(measureJson(await fixtureJson(`${prefix}-over.bin`))[property], over);
}

async function assertCborMetric(slug, dimension, property, exact, over) {
  const exactMetrics = measureCbor(await fixtureBytes(`cbor-${slug}-${dimension}-exact.bin`));
  const overMetrics = measureCbor(await fixtureBytes(`cbor-${slug}-${dimension}-over.bin`));
  assert.equal(exactMetrics[property], exact);
  assert.equal(overMetrics[property], over);
}

async function fixtureJson(name) {
  return JSON.parse(await fixtureBytes(name));
}

async function fixtureBytes(name) {
  return readFile(new URL(`../schemas/conformance/v0/shared/${name}`, import.meta.url));
}

function measureJson(value) {
  const metrics = {
    decodedTextBytes: 0,
    maxArrayElements: 0,
    maxDepth: 0,
    maxKeyBytes: 0,
    maxObjectMembers: 0,
    maxStringBytes: 0,
    nodes: 0,
    totalElements: 0,
    totalMembers: 0,
  };
  visit(value, 1);
  return metrics;

  function visit(node, depth) {
    metrics.nodes += 1;
    metrics.maxDepth = Math.max(metrics.maxDepth, depth);
    if (typeof node === "string") {
      const bytes = Buffer.byteLength(node);
      metrics.decodedTextBytes += bytes;
      metrics.maxStringBytes = Math.max(metrics.maxStringBytes, bytes);
    } else if (Array.isArray(node)) {
      metrics.totalElements += node.length;
      metrics.maxArrayElements = Math.max(metrics.maxArrayElements, node.length);
      for (const child of node) {
        visit(child, depth + 1);
      }
    } else if (node !== null && typeof node === "object") {
      const entries = Object.entries(node);
      metrics.totalMembers += entries.length;
      metrics.maxObjectMembers = Math.max(metrics.maxObjectMembers, entries.length);
      for (const [key, child] of entries) {
        const bytes = Buffer.byteLength(key);
        metrics.decodedTextBytes += bytes;
        metrics.maxKeyBytes = Math.max(metrics.maxKeyBytes, bytes);
        visit(child, depth + 1);
      }
    }
  }
}

function measureCbor(bytes) {
  const metrics = { items: 0, maxArrayElements: 0, maxDepth: 0, maxMapEntries: 0 };
  let offset = 0;
  visit(1);
  assert.equal(offset, bytes.length, "accepted CBOR boundary fixture must contain one data item");
  return metrics;

  function visit(depth) {
    metrics.items += 1;
    metrics.maxDepth = Math.max(metrics.maxDepth, depth);
    const initial = bytes[offset];
    offset += 1;
    const majorType = initial >> 5;
    const argument = readArgument(initial & 0x1f);
    if (majorType === 0 || majorType === 1) {
      return;
    }
    if (majorType === 2 || majorType === 3) {
      offset += argument;
      assert.ok(offset <= bytes.length, "CBOR string length must stay inside retained bytes");
      return;
    }
    if (majorType === 4) {
      metrics.maxArrayElements = Math.max(metrics.maxArrayElements, argument);
      for (let index = 0; index < argument; index += 1) {
        visit(depth + 1);
      }
      return;
    }
    if (majorType === 5) {
      metrics.maxMapEntries = Math.max(metrics.maxMapEntries, argument);
      for (let index = 0; index < argument * 2; index += 1) {
        visit(depth + 1);
      }
      return;
    }
    assert.fail(`unsupported accepted CBOR fixture major type: ${majorType}`);
  }

  function readArgument(additionalInformation) {
    if (additionalInformation < 24) {
      return additionalInformation;
    }
    const byteLength = { 24: 1, 25: 2, 26: 4, 27: 8 }[additionalInformation];
    assert.ok(byteLength, "accepted CBOR fixture must use a definite integer argument");
    assert.ok(offset + byteLength <= bytes.length, "CBOR argument must stay inside retained bytes");
    let value = 0;
    for (let index = 0; index < byteLength; index += 1) {
      value = value * 256 + bytes[offset + index];
      assert.ok(Number.isSafeInteger(value), "CBOR fixture argument must remain a safe integer");
    }
    offset += byteLength;
    return value;
  }
}
