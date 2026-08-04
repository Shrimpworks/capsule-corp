import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { readFile } from "node:fs/promises";
import test from "node:test";

const contractUrl = new URL(
  "../schemas/conformance/c1-governed-deno-core/controlled-development-profile.json",
  import.meta.url,
);
const exact = await readFile(contractUrl);
const contract = JSON.parse(exact);

const artifacts = Object.freeze({
  rustyV8Raw: "e964d6b1b3689e91f8cf488d8a9f05764a03434b2e2e8347be5067300d39a7de",
  rustyV8Gzip: "1ae209c9e4ba5803d010d2c79ee4cc0af0126c5a7ebcca211c7e41deaede4cd2",
  denoCoreBinary: "56d3acefd2cc2f5136a0b8143c47131e49a58fbf66382dfd3e84f715ce8e2898",
  snapshot: "4e8965217d5a6675a880326eee6f690bbeec7e7cb243decf2f3e9f453a871a2c",
  twoFileBundle: "0cc08f93e82fcfe68b033e8807975a3bd67068a817da811a87a73aedc3f23937",
  rootManifest: "100832dbb37737f29341bc5404df6d4405b8d6b706f274028892801fa88e7de8",
  rootTar: "9c46b45c4d220aedcc47c9ee53e875bc71d31d0b881b51740aaa9b882b5741e6",
  rootGzip: "e847651b35cd425dd8f6fe3bd45d693aff0af244e3a7bd30c629fa125cac62e8",
});

test("independently verifies the passive C1 known answer", () => {
  assert.equal(exact.length, 9_289);
  assert.equal(
    createHash("sha256").update(exact).digest("hex"),
    "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e",
  );
  assert.equal(
    contract.contract,
    "capsule.governed-deno-core.controlled-development-composition/c1",
  );
  assert.equal(contract.status, "passive-not-admitted");
  assert.equal(contract.evidence.mergeCommit, "fa03d7043b4f0653081d6c5733d597f49f6efd1c");
  assert.equal(contract.forks.deno.head, "9adb0b68b55bca81644827f1e7749a3acb091bed");
  assert.equal(contract.forks.rustyV8.head, "80e863ddb942a4aa2b384e794fc23e35b9d2bb15");
  assert.deepEqual(
    Object.fromEntries(
      Object.entries(contract.artifacts).map(([name, value]) => [name, value.sha256]),
    ),
    artifacts,
  );
  assert.ok(Object.values(contract.artifacts).every((artifact) => artifact.admitted === false));
});

test("closes the C1 app and runtime surfaces", () => {
  assert.deepEqual(contract.application.source, {
    profile: "capsule.mjs-source/v0",
    mediaType: "application/capsule.javascript-source;v=0;module=esm",
    fileCount: 1,
    path: "main.mjs",
    entrypoint: "main.mjs",
    minimumBytes: 0,
    maximumBytes: 262_144,
    binding: "ExecutionPlan.sourceManifestDigest+sourceByteLength+sourceEntrypoint",
    transformation: "none-byte-exact",
  });
  assert.deepEqual(contract.application.main, {
    binding: "globalThis.capsuleMain",
    receiver: "globalThis",
    requiredType: "function",
    invocations: 1,
    awaitResult: true,
    arguments: ["parsed-canonical-inline-json"],
  });
  assert.equal(contract.application.input.maximumBytes, 262_144);
  assert.equal(
    contract.application.input.binding,
    "ExecutionPlan.inlineInputDigest+inlineInputByteLength",
  );
  assert.equal(contract.application.completion.maximumJsonBytes, 262_144);
  assert.equal(contract.application.completion.maximumPhysicalFrameBytes, 262_368);
  assert.equal(contract.application.completion.limitBinding, "ExecutionPlan.outputMaxJsonBytes");
  assert.equal(
    contract.application.completion.attemptBinding,
    "P0-3-attempt-bound-completion-frame",
  );
  assert.equal(contract.application.completion.workloadOwnsEndpoint, false);

  assert.deepEqual(contract.runtimeSurface.builtinOps, [
    "op_get_ext_import_meta_proto",
    "op_get_extras_binding_object",
    "op_set_captured_bootstrap",
  ]);
  assert.equal(contract.runtimeSurface.moduleLoader, "none");
  assert.deepEqual(contract.runtimeSurface.extensions, []);
  assert.deepEqual(contract.runtimeSurface.modules, [
    { identity: "capsule:main.mjs", role: "main", source: "registered-main.mjs" },
  ]);
  assert.deepEqual(contract.runtimeSurface.workloadFiles.readable, []);
  assert.deepEqual(contract.runtimeSurface.workloadFiles.writable, []);

  const overlap = contract.runtimeSurface.permittedGlobals.filter((name) =>
    contract.runtimeSurface.forbiddenGlobals.includes(name),
  );
  assert.deepEqual(overlap, []);
  for (const forbidden of [
    "Deno",
    "Function",
    "WebAssembly",
    "Worker",
    "eval",
    "fetch",
    "process",
  ]) {
    assert.ok(contract.runtimeSurface.forbiddenGlobals.includes(forbidden));
  }
});

test("keeps descriptors, resources, and all authority effects inactive", () => {
  assert.deepEqual(contract.descriptors.constructionProbeObserved, [0, 1, 2]);
  assert.deepEqual(
    contract.descriptors.logicalRoles.map((role) => role.role),
    ["runtime-root", "registered-source", "approved-inline-input", "completion-result"],
  );
  assert.equal(contract.descriptors.numericAssignment, "c2-required-unselected");
  assert.equal(contract.descriptors.workloadCompletionAccess, "none");
  assert.equal(contract.resources.transportProfileRef, "capsule.gate-c.p0-3.measured-limits/v0");
  assert.equal(contract.resources.machineProfileRef, null);
  assert.equal(contract.resources.activation, "refuse-until-c2-and-admission");
  assert.deepEqual(contract.effects, {
    process: false,
    runtime: false,
    backend: false,
    guest: false,
    admission: false,
    signing: false,
    release: false,
  });
  assert.ok(contract.refusalCodes.includes("C1_RUNTIME_NOT_ADMITTED"));
  assert.equal(contract.stopConditions.length, 10);
});
