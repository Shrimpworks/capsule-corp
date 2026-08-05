import assert from "node:assert/strict";
import { readdir, readFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";
import {
  C2B_PASSIVE_BINDING_V2_KNOWN_ANSWER,
  decodeC2BPassiveBindingV2,
  validateC2BPassiveBindingV2,
} from "./governed-deno-core-c2b-passive-binding-v2.js";
import type { StrictJsonValue } from "./strict-job-proposal-json.js";

type MutableJsonValue =
  | null
  | boolean
  | number
  | string
  | MutableJsonValue[]
  | { [key: string]: MutableJsonValue };

const root = new URL("../../../", import.meta.url);
const exact = await read("schemas/conformance/c2b-governed-deno-core/passive-binding-v2.json");
const v1 = await read("schemas/conformance/c2b-governed-deno-core/passive-binding.json");

test("decodes the immutable no-guest build-closure successor", () => {
  const binding = decodeC2BPassiveBindingV2(exact);
  assert.equal(binding.schemaVersion, 2);
  assert.equal(binding.archiveEvidence.mergeCommit, "50108417ebf1aa45788a4e9a6b4ca6b4448e9972");
  assert.equal(binding.archiveEvidence.reviewedHead, "518eea04e1f81d27e61178b7f4ff164b955dea76");
  assert.equal(binding.constructedArtifacts.length, 6);
  assert.equal(binding.runtimeManifestCandidate.canonicalAuthority, null);
  assert.equal(binding.hostPreflightHarness.finalHostRunnerIdentity, null);
  assert.equal(binding.hostPreflightHarness.final, false);
  assert.equal(binding.hostPreflightHarness.callsLibkrun, false);
  assert.deepEqual(Object.values(binding.effects), Array(13).fill(false));
  assert.deepEqual(
    Object.entries(binding.unresolved)
      .filter(([key]) => key !== "refusal")
      .map(([, value]) => value),
    Array(11).fill(null),
  );
  assert.equal(binding.nextConsumerGate.implemented, false);
  assert.equal(binding.nextConsumerGate.admissionEffect, false);
});

test("rejects exact cap plus one and cross-version bytes", () => {
  assert.equal(exact.length, C2B_PASSIVE_BINDING_V2_KNOWN_ANSWER.bytes);
  assert.throws(() => decodeC2BPassiveBindingV2(Buffer.concat([exact, Buffer.from(" ")])), /CAP/u);
  assert.throws(() => decodeC2BPassiveBindingV2(v1), /C2B_V2_BINDING_/u);
});

test("strict decoding rejects unknown, missing, duplicate, and trailing bytes", () => {
  const source = new TextDecoder().decode(exact);
  const original = JSON.parse(source) as Record<string, unknown>;
  const missing = structuredClone(original);
  delete missing.domain;
  const cases = [
    source.replace("{\n", '{\n  "unknown": true,\n'),
    source.replace("{\n", '{\n  "objectType": "capsule.governed-deno-core-c2b-passive-binding",\n'),
    `${JSON.stringify(missing, null, 2)}\n`,
    `${source}{}`,
  ];
  for (const candidate of cases) {
    assert.throws(() => decodeC2BPassiveBindingV2(new TextEncoder().encode(candidate)));
  }
});

test("rejects substitution, cross-version, null-closure, order, and admission mutations", () => {
  const original = JSON.parse(new TextDecoder().decode(exact)) as MutableJsonValue;
  const mutations: readonly [string, (value: MutableJsonValue) => void][] = [
    ["version", (value) => setAt(value, ["schemaVersion"], 1)],
    ["domain", (value) => setAt(value, ["domain"], "capsule.other/v2")],
    [
      "archive-substitution",
      (value) => setAt(value, ["archiveEvidence", "mergeCommit"], "0".repeat(40)),
    ],
    [
      "artifact-substitution",
      (value) => setAt(value, ["constructedArtifacts", 0, "sha256"], "0".repeat(64)),
    ],
    ["artifact-order", (value) => swapAt(value, ["constructedArtifacts"], 0, 1)],
    ["runner-final", (value) => setAt(value, ["hostPreflightHarness", "final"], true)],
    [
      "runner-identity",
      (value) => setAt(value, ["hostPreflightHarness", "finalHostRunnerIdentity"], "invented"),
    ],
    [
      "firmware-identity",
      (value) => setAt(value, ["unresolved", "separateFirmwareIdentity"], "invented"),
    ],
    ["cpu-limit", (value) => setAt(value, ["unresolved", "cpuTimeLimitMs"], 1)],
    ["scratch-limit", (value) => setAt(value, ["unresolved", "scratchMaximumBytes"], 1)],
    ["guest-evidence", (value) => setAt(value, ["unresolved", "guestEvidence"], "invented")],
    [
      "admission-null",
      (value) => setAt(value, ["unresolved", "runtimeProfileAdmission"], "admitted"),
    ],
    ["admission-effect", (value) => setAt(value, ["effects", "admission"], true)],
    ["consumer", (value) => setAt(value, ["nextConsumerGate", "implemented"], true)],
  ];
  for (const [name, mutate] of mutations) {
    const candidate = structuredClone(original);
    mutate(candidate);
    assert.throws(
      () => validateC2BPassiveBindingV2(candidate as StrictJsonValue),
      /C2B_V2_SEMANTIC_KNOWN_ANSWER_/u,
      name,
    );
  }
});

test("returns independently owned deeply frozen values", () => {
  const first = decodeC2BPassiveBindingV2(exact);
  const second = decodeC2BPassiveBindingV2(exact);
  assert.notEqual(first, second);
  assert.notEqual(first.constructedArtifacts, second.constructedArtifacts);
  assert.ok(Object.isFrozen(first));
  assert.ok(Object.isFrozen(first.constructedArtifacts));
  assert.ok(Object.isFrozen(first.constructedArtifacts[0]));
  assert.ok(Object.isFrozen(first.unresolved));
});

test("has no TypeScript product consumer call site", async () => {
  const repositoryRoot = fileURLToPath(root);
  const calls: string[] = [];
  for (const base of ["packages", "internal"]) {
    for (const path of await walk(join(repositoryRoot, base))) {
      if (
        !path.endsWith(".ts") ||
        path.endsWith(".test.ts") ||
        path.endsWith("passive-binding-v2.ts")
      ) {
        continue;
      }
      if ((await readFile(path, "utf8")).includes("decodeC2BPassiveBindingV2(")) calls.push(path);
    }
  }
  assert.deepEqual(calls, []);
});

async function read(path: string): Promise<Buffer> {
  return readFile(new URL(path, root));
}

async function walk(directory: string): Promise<string[]> {
  const paths: string[] = [];
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    if (entry.name === "node_modules" || entry.name === "dist") continue;
    const path = join(directory, entry.name);
    if (entry.isDirectory()) paths.push(...(await walk(path)));
    else paths.push(path);
  }
  return paths;
}

function setAt(
  value: MutableJsonValue,
  path: readonly (string | number)[],
  replacement: MutableJsonValue,
): void {
  const last = path.at(-1);
  if (last === undefined) throw new Error("empty mutation path");
  const parent = mutableAt(value, path.slice(0, -1));
  if (typeof last === "number") {
    if (!Array.isArray(parent)) throw new Error("mutation parent is not an array");
    parent[last] = replacement;
  } else {
    if (parent === null || typeof parent !== "object" || Array.isArray(parent)) {
      throw new Error("mutation parent is not an object");
    }
    parent[last] = replacement;
  }
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
    } else {
      if (current === null || typeof current !== "object" || Array.isArray(current)) {
        throw new Error("mutation path is not an object");
      }
      current = current[part] as MutableJsonValue;
    }
  }
  return current;
}
