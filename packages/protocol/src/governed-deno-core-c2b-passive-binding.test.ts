import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { readFile } from "node:fs/promises";
import test from "node:test";
import {
  C2B_PASSIVE_BINDING_KNOWN_ANSWER,
  type C2BPassiveBindingInputs,
  decodeC2BPassiveBinding,
  validateC2BPassiveBinding,
} from "./governed-deno-core-c2b-passive-binding.js";
import type { StrictJsonValue } from "./strict-job-proposal-json.js";

type MutableJsonValue =
  | null
  | boolean
  | number
  | string
  | MutableJsonValue[]
  | { [key: string]: MutableJsonValue };

const root = new URL("../../../", import.meta.url);
const bindingExact = await read("schemas/conformance/c2b-governed-deno-core/passive-binding.json");
const inputs = await readInputs();

test("independently decodes the passive binding and all four exact inputs", () => {
  const binding = decodeC2BPassiveBinding(bindingExact, inputs);
  assert.equal(binding.objectType, "capsule.governed-deno-core-c2b-passive-binding");
  assert.equal(binding.schemaVersion, 1);
  assert.equal(binding.workStatus.passiveReconciliation, "PASSED");
  assert.equal(binding.workStatus.parentGovernedRuntime, "IN_PROGRESS — TRENDING_GOOD");
  assert.equal(binding.workStatus.c2bComposedProfileGuestExecution, "BLOCKED");
  assert.equal(binding.workStatus.runtimeProfileAdmission, "BLOCKED");
  assert.equal(binding.workStatus.runtime001, "unsupported");
  assert.equal(binding.workStatus.vmm001, "unsupported");
  assert.equal(binding.limitations.sameHostBuildEvidence, "not-independent-builder-equality");
  assert.deepEqual(Object.values(binding.effects), Array(12).fill(false));
  assert.equal(binding.artifacts.length, 6);
  assert.equal(binding.dependencyMergePolicy.dependencies.length, 2);
  assert.equal(binding.nextConsumerGate.implemented, false);
  assert.equal(binding.nextConsumerGate.admissionEffect, false);
});

test("recomputes exact C2A and fork fixture known answers", () => {
  const binding = decodeC2BPassiveBinding(bindingExact, inputs);
  const source = Buffer.from(binding.fixedFixture.source.utf8 as string);
  const input = Buffer.from(binding.fixedFixture.input.utf8 as string);
  const completion = Buffer.from(binding.fixedFixture.completion.utf8 as string);
  assert.equal(source.length, 103);
  assert.equal(sha256(source), binding.fixedFixture.source.sha256);
  assert.equal(input.length, 36);
  assert.equal(sha256(input), binding.fixedFixture.input.sha256);
  assert.equal(completion.length, 35);
  assert.equal(sha256(completion), binding.fixedFixture.completion.sha256);
  assert.deepEqual(JSON.parse(input.toString()), { message: "capsule-c2a", value: 21 });
  assert.deepEqual(JSON.parse(completion.toString()), { doubled: 42, echo: "capsule-c2a" });
});

test("rejects semantic version, domain, length, digest, cross-link, order, cap, and admission mutations", () => {
  const original = JSON.parse(new TextDecoder().decode(bindingExact)) as MutableJsonValue;
  const mutations: readonly [string, (value: MutableJsonValue) => void][] = [
    ["wrong-version", (value) => setAt(value, ["schemaVersion"], 2)],
    ["wrong-domain", (value) => setAt(value, ["domain"], "capsule.other/v1")],
    ["wrong-length", (value) => setAt(value, ["fixedFixture", "source", "bytes"], 104)],
    ["wrong-digest", (value) => setAt(value, ["fixedFixture", "source", "sha256"], "0".repeat(64))],
    [
      "cross-link",
      (value) => setAt(value, ["buildEvidence", "forkSupplementIdentity"], "capsule.other"),
    ],
    ["substitution", (value) => setAt(value, ["forkSupplement", "head"], "0".repeat(40))],
    ["order", (value) => swapAt(value, ["artifacts"], 0, 1)],
    ["cap", (value) => setAt(value, ["fixedFixture", "source", "maximumBytes"], 262_145)],
    ["admission", (value) => setAt(value, ["effects", "admission"], true)],
  ];
  for (const [name, mutate] of mutations) {
    const candidate = structuredClone(original);
    mutate(candidate);
    assert.throws(
      () => validateC2BPassiveBinding(candidate as StrictJsonValue),
      /C2B_SEMANTIC-KNOWN-ANSWER_/u,
      name,
    );
  }
});

test("strict decoding rejects unknown, duplicate, missing, trailing, and input mutations", () => {
  const original = new TextDecoder().decode(bindingExact);
  const unknown = new TextEncoder().encode(original.replace("{\n", '{\n  "unknown": true,\n'));
  const duplicate = new TextEncoder().encode(
    original.replace(
      "{\n",
      '{\n  "objectType": "capsule.governed-deno-core-c2b-passive-binding",\n',
    ),
  );
  const missingObject = JSON.parse(original) as Record<string, unknown>;
  delete missingObject.domain;
  const missing = new TextEncoder().encode(`${JSON.stringify(missingObject, null, 2)}\n`);
  const trailing = Buffer.concat([bindingExact, Buffer.from("{}")]);
  for (const [name, value] of [
    ["unknown", unknown],
    ["duplicate", duplicate],
    ["missing", missing],
    ["trailing", trailing],
  ] as const) {
    assert.throws(() => decodeC2BPassiveBinding(value, inputs), /C2B_BINDING_/u, name);
  }

  const mutatedInputs = cloneInputs(inputs);
  const mutationIndex = Math.floor(mutatedInputs.buildEvidence.length / 2);
  mutatedInputs.buildEvidence[mutationIndex] =
    (mutatedInputs.buildEvidence[mutationIndex] as number) ^ 1;
  assert.throws(() => decodeC2BPassiveBinding(bindingExact, mutatedInputs), /C2B_EVIDENCE_DIGEST/u);
});

test("returns an independently owned deeply frozen decoded value", () => {
  const first = decodeC2BPassiveBinding(bindingExact, inputs);
  const second = decodeC2BPassiveBinding(bindingExact, inputs);
  assert.notEqual(first, second);
  assert.notEqual(first.artifacts, second.artifacts);
  assert.ok(Object.isFrozen(first));
  assert.ok(Object.isFrozen(first.artifacts));
  assert.ok(Object.isFrozen(first.artifacts[0]));
});

async function readInputs(): Promise<C2BPassiveBindingInputs> {
  return {
    c1: await read("schemas/conformance/c1-governed-deno-core/controlled-development-profile.json"),
    c2a: await read("schemas/conformance/c2a-governed-deno-core/passive-execution-profile.json"),
    forkSupplement: await read(
      "schemas/conformance/c2b-governed-deno-core/inputs/fork-supplement.json",
    ),
    buildEvidence: await read(
      "schemas/conformance/c2b-governed-deno-core/inputs/build-evidence.json",
    ),
  };
}

async function read(path: string): Promise<Buffer> {
  return readFile(new URL(path, root));
}

function sha256(value: Uint8Array): string {
  return createHash("sha256").update(value).digest("hex");
}

function cloneInputs(value: C2BPassiveBindingInputs): {
  c1: Uint8Array;
  c2a: Uint8Array;
  forkSupplement: Uint8Array;
  buildEvidence: Uint8Array;
} {
  return {
    c1: Uint8Array.from(value.c1),
    c2a: Uint8Array.from(value.c2a),
    forkSupplement: Uint8Array.from(value.forkSupplement),
    buildEvidence: Uint8Array.from(value.buildEvidence),
  };
}

function setAt(
  value: MutableJsonValue,
  path: readonly (string | number)[],
  replacement: MutableJsonValue,
): void {
  const [last] = path.slice(-1);
  if (last === undefined) throw new Error("empty mutation path");
  const parent = mutableAt(value, path.slice(0, -1));
  if (typeof last === "number") {
    if (!Array.isArray(parent)) throw new Error("mutation parent is not an array");
    parent[last] = replacement;
    return;
  }
  if (parent === null || typeof parent !== "object" || Array.isArray(parent)) {
    throw new Error("mutation parent is not an object");
  }
  parent[last] = replacement;
}

function swapAt(
  value: MutableJsonValue,
  path: readonly (string | number)[],
  a: number,
  b: number,
): void {
  const array = mutableAt(value, path);
  if (!Array.isArray(array)) throw new Error("mutation target is not an array");
  [array[a], array[b]] = [array[b] as MutableJsonValue, array[a] as MutableJsonValue];
}

function mutableAt(value: MutableJsonValue, path: readonly (string | number)[]): MutableJsonValue {
  let current = value;
  for (const part of path) {
    if (typeof part === "number") {
      if (!Array.isArray(current)) throw new Error("mutation path is not an array");
      current = current[part] as MutableJsonValue;
      continue;
    }
    if (current === null || typeof current !== "object" || Array.isArray(current)) {
      throw new Error("mutation path is not an object");
    }
    current = current[part] as MutableJsonValue;
  }
  return current;
}

assert.equal(bindingExact.length, C2B_PASSIVE_BINDING_KNOWN_ANSWER.bytes);
