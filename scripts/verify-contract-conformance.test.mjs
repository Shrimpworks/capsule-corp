import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { cp, mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { verifyConformanceCorpus } from "./verify-contract-conformance.mjs";

const schemaPath = new URL("../schemas/conformance/v0/manifest.schema.json", import.meta.url);
const checkedInCorpus = new URL("../schemas/conformance/v0/", import.meta.url);

test("verifies the checked-in raw JSON and scalar corpus", async () => {
  const result = await verifyConformanceCorpus({ rootDirectory: checkedInCorpus });

  assert.deepEqual(result, { caseCount: 65, fixtureCount: 61, ruleCount: 25 });
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
