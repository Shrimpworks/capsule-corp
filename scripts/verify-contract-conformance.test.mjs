import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { cp, mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import Ajv2020 from "ajv/dist/2020.js";
import { verifyConformanceCorpus } from "./verify-contract-conformance.mjs";

const schemaPath = new URL("../schemas/conformance/v0/manifest.schema.json", import.meta.url);
const checkedInCorpus = new URL("../schemas/conformance/v0/", import.meta.url);

test("verifies the checked-in foundational conformance corpus", async () => {
  const result = await verifyConformanceCorpus({ rootDirectory: checkedInCorpus });

  assert.deepEqual(result, { caseCount: 147, fixtureCount: 130, ruleCount: 55 });
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

test("retains independently specified source-manifest and canonical-input known answers", async () => {
  const expectedSourceManifestHex =
    "a5017763617073756c652e736f757263652d6d616e69666573740200036b7372632f6d61696e2e747304838364412e74735820f8bfa98819f02db4df7eca4c610e70f3f8fa171d31e5abfaaebfe2dfec05f67c15836b7372632f6d61696e2e747358208d154e7e242fc07090f574f7e79237f4f4f1599eea50d98e053a8a2bdcc3be12182d83647a2e74735820e6b6daa9299b27161cefff62532ecd829cd9db9513179a0466be04a74d345593181905185b";
  const expectedCanonicalInlineInput =
    '{"A":{"a":"é","b":false},"a":"quote:\\" slash:/ backslash:\\\\ control:\\u0000\\u001f rocket:🚀","z":[3,0,-2,true,null]}';
  const sourceManifest = await corpusBytes("source-manifest/ordinary.cbor");
  const canonicalInlineInput = await corpusBytes("canonical-inline-input/ordinary.json");

  assert.equal(sourceManifest.toString("hex"), expectedSourceManifestHex);
  assert.equal(
    sha256Hex(sourceManifest),
    "e5e09b2435baedf897526a89c698c0b0531437a69472372ae426f62d801fc171",
  );
  assert.equal(canonicalInlineInput.toString("utf8"), expectedCanonicalInlineInput);
  assert.equal(
    sha256Hex(canonicalInlineInput),
    "bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72",
  );

  for (const path of ["job-proposal/ordinary.json", "job-proposal/equivalent-public-json.json"]) {
    const proposal = JSON.parse(await corpusBytes(path));
    assert.deepEqual(encodeSourceManifestForTest(proposal.source), sourceManifest, path);
    assert.deepEqual(canonicalInlineJsonForTest(proposal.input.value), canonicalInlineInput, path);
  }
});

test("retains exact proposal semantic boundaries and fixed resolver results", async () => {
  const ordinaryProposal = await corpusJson("job-proposal/ordinary.json");
  assert.deepEqual(Object.keys(ordinaryProposal), [
    "apiVersion",
    "kind",
    "source",
    "runtimeProfile",
    "input",
    "requestedLimits",
    "outputs",
    "labels",
  ]);

  const exactPathProposal = await corpusJson("job-proposal/source-path-exact.json");
  const overPathProposal = await corpusJson("job-proposal/source-path-over.json");
  assert.equal(Buffer.byteLength(exactPathProposal.source.entrypoint), 256);
  assert.equal(Buffer.byteLength(overPathProposal.source.entrypoint), 257);

  const exactSegmentProposal = await corpusJson("job-proposal/source-segment-exact.json");
  const overSegmentProposal = await corpusJson("job-proposal/source-segment-over.json");
  assert.equal(
    Buffer.byteLength(Object.keys(exactSegmentProposal.source.files)[1].split("/")[0]),
    64,
  );
  assert.equal(
    Buffer.byteLength(Object.keys(overSegmentProposal.source.files)[1].split("/")[0]),
    65,
  );

  const exactFileCount = await corpusJson("job-proposal/source-file-count-exact.json");
  const overFileCount = await corpusJson("job-proposal/source-file-count-over.json");
  assert.equal(Object.keys(exactFileCount.source.files).length, 32);
  assert.equal(Object.keys(overFileCount.source.files).length, 33);

  const exactFileProposal = await corpusJson("job-proposal/source-file-bytes-exact.json");
  const overFileProposal = await corpusJson("job-proposal/source-file-bytes-over.json");
  assert.equal(Buffer.byteLength(exactFileProposal.source.files["main.ts"]), 262_144);
  assert.equal(Buffer.byteLength(overFileProposal.source.files["main.ts"]), 262_145);

  const exactAggregate = await corpusJson("job-proposal/source-aggregate-exact.json");
  const overAggregate = await corpusJson("job-proposal/source-aggregate-over.json");
  assert.equal(sourceAggregateBytes(exactAggregate), 1_048_576);
  assert.equal(sourceAggregateBytes(overAggregate), 1_048_577);

  const exactInput = await corpusJson("job-proposal/inline-input-bytes-exact.json");
  const overInput = await corpusJson("job-proposal/inline-input-bytes-over.json");
  assert.equal(canonicalInlineJsonForTest(exactInput.input.value).length, 262_144);
  assert.equal(canonicalInlineJsonForTest(overInput.input.value).length, 262_145);

  assert.equal(
    Object.keys((await corpusJson("job-proposal/labels-count-exact.json")).labels).length,
    8,
  );
  assert.equal(
    Object.keys((await corpusJson("job-proposal/labels-count-over.json")).labels).length,
    9,
  );
  assert.equal(
    Buffer.byteLength(
      Object.keys((await corpusJson("job-proposal/label-key-exact.json")).labels)[0],
    ),
    32,
  );
  assert.equal(
    Buffer.byteLength(
      Object.keys((await corpusJson("job-proposal/label-key-over.json")).labels)[0],
    ),
    33,
  );
  assert.equal(
    Buffer.byteLength(
      Object.values((await corpusJson("job-proposal/label-value-exact.json")).labels)[0],
    ),
    128,
  );
  assert.equal(
    Buffer.byteLength(
      Object.values((await corpusJson("job-proposal/label-value-over.json")).labels)[0],
    ),
    129,
  );

  const profileRegistry = await corpusJson("contexts/profile-registry.json");
  const userPolicy = await corpusJson("contexts/user-policy.json");
  assert.deepEqual(profileRegistry.profiles, [
    {
      alias: "fixture-active@1",
      status: "active",
      exactWallTimeMs: { minimum: 1, maximum: 10_000 },
    },
    {
      alias: "fixture-inactive@1",
      status: "inactive",
      exactWallTimeMs: { minimum: 1, maximum: 10_000 },
    },
  ]);
  assert.deepEqual(userPolicy.wallTimeMs, { trustedDefault: 5_000, ceiling: 10_000 });

  const manifest = await corpusJson("manifest.json");
  const cases = new Map(manifest.cases.map((entry) => [entry.id, entry]));
  assert.deepEqual(cases.get("job-proposal.policy.wall-time.requested").context.oracle.wallTime, {
    milliseconds: 5_000,
    origin: "requested",
  });
  assert.deepEqual(
    cases.get("job-proposal.policy.wall-time.trusted-default").context.oracle.wallTime,
    { milliseconds: 5_000, origin: "trusted-default" },
  );
  assert.deepEqual(
    cases.get("job-proposal.policy.wall-time.exact-ceiling").context.oracle.wallTime,
    { milliseconds: 10_000, origin: "requested" },
  );
  assert.equal(
    cases.get("job-proposal.policy.wall-time.ceiling-plus-one").context.oracle.wallTime,
    null,
  );
});

test("assigns Task 2.3 rejections to the first owner without authority change", async () => {
  const manifest = await corpusJson("manifest.json");
  const cases = new Map(manifest.cases.map((entry) => [entry.id, entry]));
  const expectedRejections = new Map([
    [
      "job-proposal.schema.authority-omission.reject-placeholder",
      ["job-proposal-schema", "UNSUPPORTED"],
    ],
    ["job-proposal.source.path-grammar.reject-dot-segment", ["job-proposal-schema", "SCHEMA"]],
    ["job-proposal.source.path-bytes.cap-plus-one", ["job-proposal-schema", "SCHEMA"]],
    ["job-proposal.source.segment-bytes.cap-plus-one", ["job-proposal-schema", "SCHEMA"]],
    [
      "job-proposal.source.entrypoint-membership.reject-missing",
      ["job-proposal-semantic-validator", "SEMANTIC"],
    ],
    [
      "job-proposal.source.aggregate-bytes.cap-plus-one",
      ["job-proposal-semantic-validator", "SEMANTIC"],
    ],
    [
      "job-proposal.input.canonical-bytes.cap-plus-one",
      ["job-proposal-semantic-validator", "SEMANTIC"],
    ],
    ["job-proposal.slots.fixed-roles.reject-input-substitution", ["job-proposal-schema", "DOMAIN"]],
    ["job-proposal.slots.fixed-roles.reject-unknown", ["job-proposal-schema", "SCHEMA"]],
    [
      "job-proposal.profile.resolution.reject-unknown",
      ["job-proposal-semantic-validator", "BINDING"],
    ],
    [
      "job-proposal.profile.resolution.reject-inactive",
      ["job-proposal-semantic-validator", "POLICY"],
    ],
    [
      "job-proposal.policy.wall-time.ceiling-plus-one",
      ["job-proposal-semantic-validator", "POLICY"],
    ],
  ]);
  for (const [id, [owner, classification]] of expectedRejections) {
    const entry = cases.get(id);
    assert.ok(entry, id);
    assert.deepEqual(
      entry.expected,
      { decision: "reject", classification, owner, authorityStateChanged: false },
      id,
    );
    if (entry.context.kind === "proposal-resolution") {
      assert.deepEqual(
        entry.context.oracle,
        {
          sourceManifest: null,
          sourceManifestDigest: null,
          canonicalInlineInput: null,
          inlineInputDigest: null,
          wallTime: null,
        },
        id,
      );
    }
  }
});

test("keeps schema-owned and semantic proposal cases on their documented side of the schema", async () => {
  const schema = await corpusJson("../../candidates/job-proposal-v0.schema.json");
  const validate = new Ajv2020({ allErrors: true, strict: true }).compile(schema);
  const manifest = await corpusJson("manifest.json");
  const proposalCases = manifest.cases.filter(
    (entry) => entry.fixture.path.startsWith("job-proposal/") && entry.wireFormat === "json",
  );

  for (const entry of proposalCases) {
    const proposal = await corpusJson(entry.fixture.path);
    const schemaAccepted = validate(proposal);
    if (entry.expected.decision === "reject" && entry.expected.owner === "job-proposal-schema") {
      assert.equal(schemaAccepted, false, entry.id);
    } else {
      assert.equal(schemaAccepted, true, `${entry.id}: ${JSON.stringify(validate.errors)}`);
    }
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

test("rejects an unknown proposal-resolution oracle field", async (t) => {
  const corpus = await createCorpus(t, (manifest) => {
    const fixture = manifest.cases[0].fixture;
    manifest.cases[0].context = proposalResolutionContextForTest(fixture);
    manifest.cases[0].context.oracle.unreviewed = true;
  });

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /manifest schema validation failed/u,
  );
});

test("rejects a proposal-resolution digest that differs from retained bytes", async (t) => {
  const corpus = await createCorpus(t, (manifest) => {
    const fixture = manifest.cases[0].fixture;
    manifest.cases[0].context = proposalResolutionContextForTest(fixture, {
      sourceManifest: fixture,
      sourceManifestDigest: "0".repeat(64),
    });
  });

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /source manifest digest must match the retained fixture digest/u,
  );
});

test("rejects a resolution result attached to a rejected proposal", async (t) => {
  const corpus = await createCorpus(t, (manifest) => {
    const fixture = manifest.cases[1].fixture;
    manifest.cases[1].context = proposalResolutionContextForTest(fixture, {
      wallTime: { milliseconds: 5_000, origin: "trusted-default" },
    });
  });

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /rejected proposal case .* cannot retain a resolution result/u,
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

function proposalResolutionContextForTest(fixture, oracle = {}) {
  return {
    kind: "proposal-resolution",
    profileRegistry: fixture,
    userPolicy: fixture,
    oracle: {
      sourceManifest: null,
      sourceManifestDigest: null,
      canonicalInlineInput: null,
      inlineInputDigest: null,
      wallTime: null,
      ...oracle,
    },
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
  return corpusBytes(`shared/${name}`);
}

async function corpusJson(path) {
  return JSON.parse(await corpusBytes(path));
}

async function corpusBytes(path) {
  return readFile(new URL(`../schemas/conformance/v0/${path}`, import.meta.url));
}

function sha256Hex(bytes) {
  return createHash("sha256").update(bytes).digest("hex");
}

function sourceAggregateBytes(proposal) {
  return Object.values(proposal.source.files).reduce(
    (total, content) => total + Buffer.byteLength(content),
    0,
  );
}

function encodeSourceManifestForTest(source) {
  const entries = Object.entries(source.files)
    .sort(([left], [right]) => Buffer.compare(Buffer.from(left), Buffer.from(right)))
    .map(([path, content]) => {
      const bytes = Buffer.from(content);
      return [path, createHash("sha256").update(bytes).digest(), bytes.length];
    });
  return cborForTest(
    new Map([
      [1, "capsule.source-manifest"],
      [2, 0],
      [3, source.entrypoint],
      [4, entries],
      [5, entries.reduce((total, entry) => total + entry[2], 0)],
    ]),
  );
}

function cborForTest(value) {
  if (Buffer.isBuffer(value) || value instanceof Uint8Array) {
    const bytes = Buffer.from(value);
    return Buffer.concat([cborHeadForTest(2, bytes.length), bytes]);
  }
  if (typeof value === "string") {
    const bytes = Buffer.from(value);
    return Buffer.concat([cborHeadForTest(3, bytes.length), bytes]);
  }
  if (typeof value === "number") {
    assert.ok(Number.isSafeInteger(value));
    return value >= 0 ? cborHeadForTest(0, value) : cborHeadForTest(1, -1 - value);
  }
  if (Array.isArray(value)) {
    return Buffer.concat([cborHeadForTest(4, value.length), ...value.map(cborForTest)]);
  }
  if (value instanceof Map) {
    const entries = [...value].map(([key, child]) => [cborForTest(key), cborForTest(child)]);
    entries.sort(([left], [right]) =>
      left.length === right.length ? Buffer.compare(left, right) : left.length - right.length,
    );
    return Buffer.concat([cborHeadForTest(5, entries.length), ...entries.flat()]);
  }
  assert.fail(`unsupported test CBOR value: ${String(value)}`);
}

function cborHeadForTest(majorType, argument) {
  if (argument < 24) {
    return Buffer.of((majorType << 5) | argument);
  }
  if (argument <= 0xff) {
    return Buffer.of((majorType << 5) | 24, argument);
  }
  if (argument <= 0xffff) {
    return Buffer.of((majorType << 5) | 25, argument >> 8, argument & 0xff);
  }
  if (argument <= 0xffffffff) {
    return Buffer.of(
      (majorType << 5) | 26,
      Math.floor(argument / 0x1000000) & 0xff,
      Math.floor(argument / 0x10000) & 0xff,
      Math.floor(argument / 0x100) & 0xff,
      argument & 0xff,
    );
  }
  assert.fail(`test CBOR argument is too large: ${argument}`);
}

function canonicalInlineJsonForTest(value) {
  return Buffer.from(canonicalInlineJsonTextForTest(value));
}

function canonicalInlineJsonTextForTest(value) {
  if (value === null) {
    return "null";
  }
  if (typeof value === "boolean") {
    return value ? "true" : "false";
  }
  if (typeof value === "number") {
    assert.ok(Number.isSafeInteger(value) && !Object.is(value, -0));
    return String(value);
  }
  if (typeof value === "string") {
    return quoteCanonicalStringForTest(value);
  }
  if (Array.isArray(value)) {
    return `[${value.map(canonicalInlineJsonTextForTest).join(",")}]`;
  }
  const entries = Object.entries(value).sort(([left], [right]) =>
    Buffer.compare(Buffer.from(left), Buffer.from(right)),
  );
  return `{${entries
    .map(
      ([key, child]) =>
        `${quoteCanonicalStringForTest(key)}:${canonicalInlineJsonTextForTest(child)}`,
    )
    .join(",")}}`;
}

function quoteCanonicalStringForTest(value) {
  let result = '"';
  for (const scalar of value) {
    const codePoint = scalar.codePointAt(0);
    if (scalar === '"') {
      result += '\\"';
    } else if (scalar === "\\") {
      result += "\\\\";
    } else if (codePoint <= 0x1f) {
      result += `\\u00${codePoint.toString(16).padStart(2, "0")}`;
    } else {
      result += scalar;
    }
  }
  return `${result}"`;
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
