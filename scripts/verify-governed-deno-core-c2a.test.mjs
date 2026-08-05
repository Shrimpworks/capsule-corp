import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";
import { sha256, sha256Hex } from "./lib/fixture-bytes.mjs";

const fixtureUrl = new URL(
  "../schemas/conformance/c2a-governed-deno-core/passive-execution-profile.json",
  import.meta.url,
);
const c1Url = new URL(
  "../schemas/conformance/c1-governed-deno-core/controlled-development-profile.json",
  import.meta.url,
);
const exact = await readFile(fixtureUrl);
const c1Exact = await readFile(c1Url);
const contract = JSON.parse(exact);

test("independently verifies exact passive C2A and unchanged C1 bytes", () => {
  assert.equal(exact.length, 26_850);
  assert.equal(
    sha256Hex(exact),
    "d4ce88888186266f5d251e6246c889b1fd46d7746bb0ba56bcc4b3ce4675992f",
  );
  assert.equal(c1Exact.length, 9_289);
  assert.equal(
    sha256Hex(c1Exact),
    "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e",
  );
  assert.equal(contract.c1Binding.bytes, c1Exact.length);
  assert.equal(contract.c1Binding.sha256, sha256Hex(c1Exact));
  assert.equal(contract.c1Binding.consumption, "exact-byte-read-only");
  assert.equal(contract.c1Binding.mutation, "forbidden");
  assert.deepEqual(contract.c1PlanAndProfileBindings.executionPlan.requiredReferenceRoles, [
    "runtimeBundleDigest",
    "profileRegistryDigest",
    "backendValidationDigest",
    "backendConfigurationDigest",
  ]);
  assert.equal(contract.c1PlanAndProfileBindings.runtimeProfile.selectedIdentity, null);
  assert.equal(contract.c1PlanAndProfileBindings.runtimeProfile.selectedDigest, null);
});

test("recomputes the C2A main.mjs, SourceManifest, input, and completion known answer", () => {
  const source = Buffer.from(contract.knownAnswer.source.utf8);
  const input = Buffer.from(contract.knownAnswer.canonicalInput.utf8);
  const completion = Buffer.from(contract.knownAnswer.expectedCompletion.jsonUtf8);
  const manifest = encodeSourceManifest("main.mjs", source);

  assert.equal(source.length, 103);
  assert.equal(sha256Hex(source), contract.knownAnswer.source.sha256);
  assert.equal(manifest.length, 89);
  assert.equal(sha256Hex(manifest), contract.knownAnswer.source.sourceManifestSha256);
  assert.equal(input.length, 36);
  assert.equal(sha256Hex(input), contract.knownAnswer.canonicalInput.sha256);
  assert.deepEqual(JSON.parse(input), { message: "capsule-c2a", value: 21 });
  assert.equal(completion.length, 35);
  assert.equal(sha256Hex(completion), contract.knownAnswer.expectedCompletion.jsonSha256);
  assert.deepEqual(JSON.parse(completion), { doubled: 42, echo: "capsule-c2a" });
  assert.equal(contract.knownAnswer.expectedCompletion.commit, "required-last");
  assert.equal(contract.knownAnswer.expectedCompletion.runnerExitAlone, "never-success");
});

test("closes numeric descriptor ownership and exact transport limits", () => {
  assert.deepEqual(
    contract.runnerDescriptorProfile.entries.map(({ fd, role, accessMode }) => ({
      fd,
      role,
      accessMode,
    })),
    [
      { fd: 0, role: "runner-stdin-null", accessMode: "O_RDONLY" },
      { fd: 1, role: "runner-stdout", accessMode: "O_WRONLY" },
      { fd: 2, role: "runner-stderr", accessMode: "O_WRONLY" },
      { fd: 3, role: "record-before-start-control", accessMode: "O_RDONLY" },
      { fd: 4, role: "runtime-root", accessMode: "O_RDONLY" },
      { fd: 5, role: "registered-source-port-input", accessMode: "O_RDONLY" },
      { fd: 6, role: "approved-inline-input-port-input", accessMode: "O_RDONLY" },
      { fd: 7, role: "completion-result-port-output", accessMode: "O_WRONLY" },
    ],
  );
  assert.deepEqual(contract.runnerDescriptorProfile.numericRange, {
    minimum: 0,
    maximum: 7,
    closeFromInclusive: 8,
  });
  assert.deepEqual(
    contract.runnerDescriptorProfile.portCalls.map(({ portId, inputFd, outputFd }) => ({
      portId,
      inputFd,
      outputFd,
    })),
    [
      { portId: 0, inputFd: 5, outputFd: -1 },
      { portId: 1, inputFd: 6, outputFd: -1 },
      { portId: 2, inputFd: -1, outputFd: 7 },
    ],
  );
  assert.deepEqual(contract.transportProfile.completion, {
    jsonMaximumBytes: 262_144,
    headerBytes: 160,
    trailerBytes: 64,
    physicalMaximumBytes: 262_368,
    retainMaximumBytes: 262_369,
    hostRunnerFd: 7,
    guestLauncherFd: 5,
    eofIsCommit: false,
  });
});

test("keeps unsupported resources, final artifacts, execution, and admission refusing", () => {
  assert.deepEqual(
    {
      vcpus: contract.machineProfile.vcpus,
      guestRamMiB: contract.machineProfile.guestRamMiB,
      wallTimeMs: contract.machineProfile.wallTimeMs,
      concurrency: contract.machineProfile.concurrency,
    },
    { vcpus: 1, guestRamMiB: 256, wallTimeMs: 1_000, concurrency: 1 },
  );
  assert.equal(contract.machineProfile.cpuTimeLimitMs, null);
  assert.equal(contract.machineProfile.hostVmmMemoryLimitBytes, null);
  assert.equal(contract.machineProfile.scratch.maximumBytes, null);
  assert.equal(contract.artifactClosure.required.length, 9);
  assert.ok(
    contract.artifactClosure.required.every(
      ({ sha256: digest, admitted }) => digest === null && admitted === false,
    ),
  );
  assert.equal(contract.guestDescriptorProfile.childManifest, null);
  assert.equal(contract.c2bMatrix.length, 11);
  assert.equal(contract.c2bMatrix.flatMap(({ cases }) => cases).length, 91);
  assert.equal(contract.restorationMutations.length, 18);
  assert.deepEqual(contract.workStatus, {
    c2a: "PASSED-passive-preparation-only",
    parentGovernedRuntime: "IN_PROGRESS-TRENDING_GOOD",
    c2b: "BLOCKED",
    runtimeProfileAdmission: "BLOCKED",
    runtime001: "unsupported",
    vmm001: "unsupported",
  });
  assert.ok(Object.values(contract.effects).every((effect) => effect === false));
});

function encodeSourceManifest(path, content) {
  return cbor(
    new Map([
      [1, "capsule.source-manifest"],
      [2, 0],
      [3, path],
      [4, [[path, sha256(content), content.length]]],
      [5, content.length],
    ]),
  );
}

function cbor(value) {
  if (Buffer.isBuffer(value) || value instanceof Uint8Array) {
    const bytes = Buffer.from(value);
    return Buffer.concat([cborHead(2, bytes.length), bytes]);
  }
  if (typeof value === "string") {
    const bytes = Buffer.from(value);
    return Buffer.concat([cborHead(3, bytes.length), bytes]);
  }
  if (typeof value === "number") {
    assert.ok(Number.isSafeInteger(value));
    return value >= 0 ? cborHead(0, value) : cborHead(1, -1 - value);
  }
  if (Array.isArray(value)) {
    return Buffer.concat([cborHead(4, value.length), ...value.map(cbor)]);
  }
  if (value instanceof Map) {
    const entries = [...value].map(([key, child]) => [cbor(key), cbor(child)]);
    entries.sort(([left], [right]) =>
      left.length === right.length ? Buffer.compare(left, right) : left.length - right.length,
    );
    return Buffer.concat([cborHead(5, entries.length), ...entries.flat()]);
  }
  assert.fail(`unsupported CBOR value: ${String(value)}`);
}

function cborHead(majorType, argument) {
  if (argument < 24) return Buffer.of((majorType << 5) | argument);
  if (argument <= 0xff) return Buffer.of((majorType << 5) | 24, argument);
  if (argument <= 0xffff) {
    return Buffer.of((majorType << 5) | 25, argument >> 8, argument & 0xff);
  }
  assert.fail(`CBOR argument exceeds test encoder range: ${argument}`);
}
