import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { readFile } from "node:fs/promises";

const repository = new URL("../", import.meta.url);
const root = new URL("schemas/conformance/authenticated-local-ipc-v0/", repository);
const manifest = await json("manifest.json");

assert.equal(manifest.objectType, "capsule.authenticated-local-ipc-passive-conformance");
assert.equal(manifest.objectVersion, 0);
assert.equal(manifest.status, "passive-unwired-no-transport");
assert.equal(manifest.protocolVersion, 0);
assert.equal(manifest.audience, "capsule.execution-supervisor.local.v0");
assert.equal(manifest.roleBindingRecordVersion, 0);
assert.equal(manifest.transportEncoding, null);
assert.equal(manifest.numericMessageTags, null);
assert.equal(manifest.peerAuthenticationEvidence, null);
assert.deepEqual(manifest.caps, {
  executionPlan: 65_536,
  roleBindings: 562,
  sourceManifest: 95,
  source: 262_144,
  planRegistration: 4_096,
  registerPlanV0Request: 328_337,
  registerPlanV0Reply: 4_096,
  getRegisteredPlanV0Request: 16,
  getRegisteredPlanV0Reply: 332_433,
});
assert.equal(
  manifest.caps.registerPlanV0Request,
  manifest.caps.executionPlan +
    manifest.caps.roleBindings +
    manifest.caps.sourceManifest +
    manifest.caps.source,
);
assert.equal(
  manifest.caps.getRegisteredPlanV0Reply,
  manifest.caps.registerPlanV0Request + manifest.caps.planRegistration,
);

for (const [path, known] of Object.entries(manifest.knownAnswers)) {
  const bytes = await readFile(new URL(path, root));
  verifyReference(bytes, known, path);
  JSON.parse(bytes);
}
for (const [name, known] of Object.entries(manifest.bodyFixtures)) {
  const bytes = await readFile(new URL(known.path, repository));
  verifyReference(bytes, known, name);
}

const registerMethod = await json("register-plan-v0.method.json");
const getMethod = await json("get-registered-plan-v0.method.json");
assert.deepEqual(
  [registerMethod.method, registerMethod.methodVersion, registerMethod.roleBindingRecordVersion],
  ["RegisterPlanV0", 0, 0],
);
assert.deepEqual(
  [getMethod.method, getMethod.methodVersion, getMethod.roleBindingRecordVersion],
  ["GetRegisteredPlanV0", 0, 0],
);
assert.notEqual(registerMethod.service, getMethod.service);
assert.notEqual(registerMethod.expectedRole, getMethod.expectedRole);
assert.notEqual(registerMethod.purpose, getMethod.purpose);
assert.equal(registerMethod.audience, getMethod.audience);

const registerRequest = await json("register-plan-v0.request.json");
const registerReply = await json("register-plan-v0.reply.json");
const getRequest = await json("get-registered-plan-v0.request.json");
const getReply = await json("get-registered-plan-v0.reply.json");
for (const envelope of [registerRequest, registerReply, getRequest, getReply]) {
  assert.equal(envelope.objectVersion, 0);
  assert.equal(envelope.fixtureSerialization, "exact-json-not-xpc-framing");
}
assert.equal(registerRequest.requestId, registerReply.requestId);
assert.equal(getRequest.requestId, getReply.requestId);
assert.equal(registerRequest.installationId, getRequest.installationId);
assert.equal(registerRequest.epochSequence, getRequest.epochSequence);
assert.equal(registerRequest.epochDigest, getRequest.epochDigest);
assert.equal(sumReferences(registerRequest.body), registerRequest.applicationDataBytes);
assert.equal(sumReferences(registerReply.body), registerReply.applicationDataBytes);
assert.equal(getRequest.applicationDataBytes, 16);
assert.equal(sumReferences(getReply.body), getReply.applicationDataBytes);

const oracles = await json("oracles.json");
assert.equal(oracles.maxima.length, manifest.maximumCaseCount);
assert.equal(oracles.refusals.length, manifest.refusalCaseCount);
for (const candidate of oracles.maxima) {
  assert.equal(
    Object.values(candidate.fieldLengths).reduce((sum, length) => sum + length, 0),
    candidate.applicationDataBytes,
    candidate.id,
  );
  if (candidate.expected.decision === "reject") assertZeroState(candidate.noState, candidate.id);
  if (candidate.expected.decision === "local-integrity-fault") {
    assert.equal(candidate.noState, null, candidate.id);
    assert.equal(candidate.postCoreState, "not-claimed-by-passive-output-oracle", candidate.id);
    assert.equal(candidate.processTerminationRequired, true, candidate.id);
  }
}
for (const candidate of oracles.refusals) assertZeroState(candidate.noState, candidate.id);
assert.deepEqual(oracles.zeroEffects, {
  ipcEndpoint: false,
  peerAuthenticated: false,
  keyUsed: false,
  processCreated: false,
  runtimeCreated: false,
  backendCreated: false,
  guestCreated: false,
});
assert.equal(oracles.responseLoss[0].requestIdIsIdempotencyKey, false);
assert.equal(oracles.responseLoss[0].retryCommittedRegistrations, 2);
assert.equal(oracles.responseLoss[0].bothRegistrationsSeparatelyReadable, true);
assert.equal(oracles.responseLoss[1].requestIdIsIdempotencyKey, false);
assert.equal(oracles.responseLoss[1].repeatedReplyBodyByteEqual, true);

const combined = JSON.stringify({
  manifest,
  registerMethod,
  getMethod,
  registerRequest,
  registerReply,
  getRequest,
  getReply,
  oracles,
});
assert.equal(combined.includes("RegisterPlanV1"), false);
assert.equal(combined.includes("GetRegisteredPlanV1"), false);
assert.equal(combined.includes("typescript"), false);
assert.equal(combined.includes("626"), false);

process.stdout.write(
  "verified independent TypeScript passive authenticated-local-IPC known answers\n",
);

async function json(path) {
  return JSON.parse(await readFile(new URL(path, root), "utf8"));
}

function verifyReference(bytes, known, label) {
  assert.equal(bytes.length, known.byteLength, `${label} byte length`);
  assert.equal(createHash("sha256").update(bytes).digest("hex"), known.sha256, `${label} digest`);
}

function sumReferences(value) {
  return Object.values(value).reduce((sum, field) => sum + field.byteLength, 0);
}

function assertZeroState(value, label) {
  assert.deepEqual(
    value,
    {
      authorityStateChanged: false,
      coreCalls: 0,
      registrationsCreated: 0,
      approvalsConsumed: 0,
      attemptsCreated: 0,
      lifecycleCalls: 0,
      backendCalls: 0,
    },
    label,
  );
}
