import assert from "node:assert/strict";
import { readdir, readFile } from "node:fs/promises";
import { join } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";
import {
  C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER,
  decodeC2BPassiveBindingV3,
  validateC2BPassiveBindingV3,
} from "./governed-deno-core-c2b-passive-binding-v3.js";
import type { StrictJsonValue } from "./strict-job-proposal-json.js";

type MutableJsonValue =
  | null
  | boolean
  | number
  | string
  | MutableJsonValue[]
  | { [key: string]: MutableJsonValue };

const root = new URL("../../../", import.meta.url);
const exact = await read("schemas/conformance/c2b-governed-deno-core/passive-binding-v3.json");
const v2 = await read("schemas/conformance/c2b-governed-deno-core/passive-binding-v2.json");

test("decodes the passive fixed-owned-guest successor contract", () => {
  const binding = decodeC2BPassiveBindingV3(exact);
  assert.equal(binding.schemaVersion, 3);
  assert.equal(
    objectAt(binding.governedInputs, "deno").acceptedCommit,
    "3fa21d1ae7705ab4bcb4bc98955f25301b20122a",
  );
  assert.equal(
    objectAt(binding.governedInputs, "rustyV8").acceptedCommit,
    "d09221062280ae1675fe26c53c3f43871aae2055",
  );
  assert.equal(
    objectAt(binding.governedInputs, "libkrun").acceptedCommit,
    "7432eda5a49220976b0167005aa43ee622f9d632",
  );
  assert.equal(
    binding.composedProfile.contractDigestSha256,
    C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER.contractSha256,
  );
  assert.equal(binding.nextGuestGate.state, "BLOCKED");
  assert.equal(binding.nextGuestGate.guestAuthorization, false);
  assert.deepEqual(Object.values(binding.effects), Array(13).fill(false));
  assert.equal(containsNull(binding as unknown as StrictJsonValue), false);
  assert.equal("cpuTimeLimitMs" in binding.resourceContract, false);
  assert.equal("hostVmmMemoryLimitBytes" in binding.resourceContract, false);
  assert.equal("scratchMaximumBytes" in binding.resourceContract, false);
});

test("rejects exact cap plus one and cross-version bytes", () => {
  assert.equal(exact.length, C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER.bytes);
  assert.throws(() => decodeC2BPassiveBindingV3(Buffer.concat([exact, Buffer.from(" ")])), /CAP/u);
  assert.throws(() => decodeC2BPassiveBindingV3(v2), /C2B_V3_BINDING_/u);
});

test("strict decoding rejects unknown, missing, duplicate, trailing, and null bytes", () => {
  const source = new TextDecoder().decode(exact);
  const original = JSON.parse(source) as Record<string, unknown>;
  const missing = structuredClone(original);
  delete missing.domain;
  const withNull = structuredClone(original) as Record<string, unknown>;
  withNull.unknown = null;
  const cases = [
    source.replace("{\n", '{\n  "unknown": true,\n'),
    source.replace("{\n", '{\n  "objectType": "capsule.governed-deno-core-c2b-passive-binding",\n'),
    `${JSON.stringify(missing, null, 2)}\n`,
    `${JSON.stringify(withNull, null, 2)}\n`,
    `${source}{}`,
  ];
  for (const candidate of cases) {
    assert.throws(
      () => decodeC2BPassiveBindingV3(new TextEncoder().encode(candidate)),
      /C2B_V3_BINDING_/u,
    );
  }
});

test("rejects source, runner, device, runtime, resource, teardown, digest, and authority mutations", () => {
  const original = JSON.parse(new TextDecoder().decode(exact)) as MutableJsonValue;
  const mutations: readonly [string, (value: MutableJsonValue) => void][] = [
    [
      "libkrun-source",
      (value) => setAt(value, ["governedInputs", "libkrun", "acceptedCommit"], "0".repeat(40)),
    ],
    [
      "stale-libkrun-selection",
      (value) => setAt(value, ["artifactRoles", "libkrun", "selectedArtifactState"], "selected"),
    ],
    ["runner-final", (value) => setAt(value, ["hostRunnerContract", "artifactState"], "PASSED")],
    [
      "implicit-vsock",
      (value) => setAt(value, ["descriptorAndDeviceContract", "implicitVsock"], "enabled"),
    ],
    ["module-loader", (value) => setAt(value, ["runtimeContract", "moduleLoader"], "filesystem")],
    [
      "string-code-generation",
      (value) => setAt(value, ["runtimeContract", "stringCodeGeneration"], "enabled"),
    ],
    ["cpu-limit", (value) => setAt(value, ["resourceContract", "cpuTimeLimitMs"], 1)],
    ["teardown", (value) => setAt(value, ["teardownContract", "forcedSignal"], "SIGTERM")],
    [
      "profile-digest",
      (value) => setAt(value, ["composedProfile", "contractDigestSha256"], "0".repeat(64)),
    ],
    ["guest-authority", (value) => setAt(value, ["nextGuestGate", "guestAuthorization"], true)],
    ["admission", (value) => setAt(value, ["effects", "admission"], true)],
  ];
  for (const [name, mutate] of mutations) {
    const candidate = structuredClone(original);
    mutate(candidate);
    assert.throws(
      () => validateC2BPassiveBindingV3(candidate as StrictJsonValue),
      /C2B_V3_/u,
      name,
    );
  }
});

test("returns independently owned deeply frozen values", () => {
  const first = decodeC2BPassiveBindingV3(exact);
  const second = decodeC2BPassiveBindingV3(exact);
  assert.notEqual(first, second);
  assert.notEqual(first.artifactRoles, second.artifactRoles);
  assert.ok(Object.isFrozen(first));
  assert.ok(Object.isFrozen(first.artifactRoles));
  assert.ok(Object.isFrozen(objectAt(first.artifactRoles, "libkrun")));
  assert.ok(Object.isFrozen(first.nextGuestGate));
});

test("has no TypeScript product consumer call site", async () => {
  const repositoryRoot = fileURLToPath(root);
  const calls: string[] = [];
  for (const base of ["packages", "internal"]) {
    for (const path of await walk(join(repositoryRoot, base))) {
      if (
        !path.endsWith(".ts") ||
        path.endsWith(".test.ts") ||
        path.endsWith("passive-binding-v3.ts")
      ) {
        continue;
      }
      if ((await readFile(path, "utf8")).includes("decodeC2BPassiveBindingV3(")) calls.push(path);
    }
  }
  assert.deepEqual(calls, []);
});

function objectAt(value: Readonly<Record<string, StrictJsonValue>>, key: string) {
  const child = value[key];
  if (child === null || typeof child !== "object" || Array.isArray(child)) {
    throw new Error(`expected object at ${key}`);
  }
  return child as Readonly<Record<string, StrictJsonValue>>;
}

function containsNull(value: StrictJsonValue): boolean {
  if (value === null) return true;
  if (Array.isArray(value)) return value.some(containsNull);
  if (typeof value === "object") return Object.values(value).some(containsNull);
  return false;
}

function read(path: string): Promise<Buffer> {
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
