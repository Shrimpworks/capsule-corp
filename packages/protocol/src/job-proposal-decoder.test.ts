import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

import {
  type DecodedJobProposalCandidate,
  decodeJobProposal,
  type JobProposalDecodeResult,
} from "./index.js";
import { predecodeStrictJobProposalJson } from "./strict-job-proposal-json.js";

interface ConformanceCase {
  readonly id: string;
  readonly object: string;
  readonly fixture: { readonly path: string };
  readonly expected: {
    readonly decision: "accept" | "reject";
    readonly classification: string | null;
    readonly owner: string;
  };
  readonly implementations: { readonly typescript: string };
}

interface ConformanceManifest {
  readonly cases: readonly ConformanceCase[];
}

const corpusRoot = new URL("../../../schemas/conformance/v0/", import.meta.url);
const manifest = JSON.parse(
  await readFile(new URL("manifest.json", corpusRoot), "utf8"),
) as ConformanceManifest;
const task3aCases = manifest.cases.filter(
  (entry) =>
    entry.object === "JobProposal" &&
    entry.implementations.typescript === "verified" &&
    ["public-raw-decoder", "job-proposal-schema"].includes(entry.expected.owner),
);

test("Task 3A marks exactly its raw and schema manifest slice verified", () => {
  assert.equal(task3aCases.length, 57);
  assert.equal(
    task3aCases.filter((entry) => entry.expected.owner === "public-raw-decoder").length,
    38,
  );
  assert.equal(
    task3aCases.filter((entry) => entry.expected.owner === "job-proposal-schema").length,
    19,
  );
  assert.equal(new Set(task3aCases.map((entry) => entry.fixture.path)).size, 55);
});

for (const entry of task3aCases) {
  test(`Task 3A conformance: ${entry.id}`, async () => {
    const bytes = new Uint8Array(await readFile(new URL(entry.fixture.path, corpusRoot)));
    if (entry.expected.owner === "public-raw-decoder") {
      const result = predecodeStrictJobProposalJson(bytes);
      assert.equal(result.ok, entry.expected.decision === "accept", entry.id);
      if (!result.ok) {
        assert.deepEqual(
          [result.refusal.owner, result.refusal.classification],
          [entry.expected.owner, entry.expected.classification],
        );
      }
      return;
    }

    const result = decodeJobProposal(bytes);
    assert.equal(result.ok, entry.expected.decision === "accept", entry.id);
    if (!result.ok) {
      assert.deepEqual(
        [result.refusal.owner, result.refusal.classification],
        [entry.expected.owner, entry.expected.classification],
      );
    }
  });
}

test("decoder accepts bytes only and returns fixed internal refusal data", () => {
  const preParsed = decodeJobProposal({ kind: "JobProposal" } as unknown as Uint8Array);
  assert.deepEqual(preParsed, {
    ok: false,
    refusal: {
      owner: "public-raw-decoder",
      classification: "MALFORMED",
      code: "INPUT_TYPE",
    },
  });

  const request = ordinaryProposal();
  const marker = "guest-controlled-marker";
  const result = decodeObject({ ...request, [marker]: false });
  assert.deepEqual(result, {
    ok: false,
    refusal: {
      owner: "job-proposal-schema",
      classification: "UNSUPPORTED",
      code: "UNKNOWN_FIELD",
    },
  });
  assert.equal(JSON.stringify(result).includes(marker), false);
});

test("decoder closes object version, type, and nested fields", () => {
  assertRefusal(decodeObject({ ...ordinaryProposal(), apiVersion: "capsule.dev/v1" }), [
    "job-proposal-schema",
    "UNSUPPORTED",
    "OBJECT_VERSION",
  ]);
  assertRefusal(decodeObject({ ...ordinaryProposal(), kind: "ExecutionPlan" }), [
    "job-proposal-schema",
    "DOMAIN",
    "OBJECT_TYPE",
  ]);
  const proposal = ordinaryProposal();
  assertRefusal(
    decodeObject({
      ...proposal,
      requestedLimits: { ...proposal.requestedLimits, cpuTimeMs: 1 },
    }),
    ["job-proposal-schema", "UNSUPPORTED", "UNKNOWN_FIELD"],
  );
});

test("decoder copies caller bytes and freezes the passive candidate graph", () => {
  const bytes = new TextEncoder().encode(JSON.stringify(ordinaryProposal()));
  const result = decodeJobProposal(bytes);
  assert.equal(result.ok, true);
  if (!result.ok) {
    return;
  }
  const candidate: DecodedJobProposalCandidate = result.proposal;
  bytes.fill(0);

  assert.equal(candidate.source.files["main.mjs"], "export {};\n");
  assert.deepEqual(candidate.input.value, { message: "hello" });
  assert.equal(Object.isFrozen(candidate), true);
  assert.equal(Object.isFrozen(candidate.source), true);
  assert.equal(Object.isFrozen(candidate.source.files), true);
  assert.equal(Object.isFrozen(candidate.outputs), true);
  assert.equal(Object.isFrozen(candidate.input.value), true);
  assert.equal(Reflect.set(candidate.source.files, "main.mjs", "changed"), false);
  assert.equal(candidate.source.files["main.mjs"], "export {};\n");
});

test("raw string budget distinguishes source content from non-source strings", () => {
  const sourceAtSemanticLayer = ordinaryProposal();
  sourceAtSemanticLayer.source.files["main.mjs"] = "x".repeat(70_000);
  const sourceResult = decodeObject(sourceAtSemanticLayer);
  assert.equal(sourceResult.ok, true);

  const inlineValueOverRawCap = ordinaryProposal();
  inlineValueOverRawCap.input.value = { value: "x".repeat(65_537) };
  assertRefusal(decodeObject(inlineValueOverRawCap), [
    "public-raw-decoder",
    "MALFORMED",
    "STRING_BYTES",
  ]);
});

function decodeObject(value: object): JobProposalDecodeResult {
  return decodeJobProposal(new TextEncoder().encode(JSON.stringify(value)));
}

function assertRefusal(
  result: JobProposalDecodeResult,
  expected: readonly [string, string, string],
): void {
  assert.equal(result.ok, false);
  if (result.ok) {
    return;
  }
  assert.deepEqual(
    [result.refusal.owner, result.refusal.classification, result.refusal.code],
    expected,
  );
}

function ordinaryProposal(): {
  apiVersion: string;
  kind: string;
  source: { entrypoint: string; files: Record<string, string> };
  runtimeProfile: string;
  input: { slot: string; kind: string; value: unknown };
  requestedLimits: Record<string, number>;
  outputs: Array<{ slot: string; kind: string; maxBytes: number }>;
  labels: Record<string, string>;
} {
  return {
    apiVersion: "capsule.dev/v0",
    kind: "JobProposal",
    source: { entrypoint: "main.mjs", files: { "main.mjs": "export {};\n" } },
    runtimeProfile: "fixture-active@1",
    input: { slot: "primary-data", kind: "inline-json", value: { message: "hello" } },
    requestedLimits: { wallTimeMs: 5_000 },
    outputs: [{ slot: "transformed-json", kind: "inline-json", maxBytes: 65_536 }],
    labels: { example: "decoder-test" },
  };
}
