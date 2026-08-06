import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { sha256Hex } from "./lib/fixture-bytes.mjs";

const repository = new URL("../", import.meta.url);
const root = new URL("schemas/conformance/authenticated-local-ipc-v0/", repository);
const manifest = await json("manifest.json");

assert.equal(manifest.objectType, "capsule.authenticated-local-ipc-passive-conformance");
assert.equal(manifest.objectVersion, 0);
assert.equal(manifest.status, "passive-unwired-native-contract-no-listener");
assert.equal(manifest.protocolVersion, 0);
assert.equal(manifest.audience, "capsule.execution-supervisor.local.v0");
assert.equal(manifest.submissionAudience, "capsule.daemon.local.v0");
assert.equal(manifest.roleBindingRecordVersion, 0);
assert.equal(manifest.transportEncoding, "xpc-dictionary-v0");
assert.equal(manifest.transportEncodingVersion, 0);
assert.equal(manifest.methodCount, 3);
assert.deepEqual(manifest.numericMessageTags, {
  invalid: 0,
  SubmitMainMJSV0: 1,
  RegisterPlanV0: 2,
  GetRegisteredPlanV0: 3,
});
assert.deepEqual(manifest.numericStatusTags, {
  OK: 0,
  MALFORMED: 1,
  UNSUPPORTED: 2,
  SCHEMA: 3,
  BINDING: 4,
  AUTHENTICATION: 5,
  STALE: 6,
  REPLAY: 7,
  CAPACITY: 8,
  TRUST_STATE: 9,
  LOCAL_FAILURE: 10,
  RECOVERY_REQUIRED: 11,
  SEMANTIC: 12,
  DOMAIN: 13,
});
assert.equal(manifest.numericReasonTags.none, 0);
assert.equal(manifest.numericReasonTags.localIntegrityFault, 14);
assert.equal(manifest.peerAuthenticationEvidence, null);
assert.deepEqual(manifest.caps, {
  jobProposal: 2_097_152,
  executionPlan: 65_536,
  roleBindings: 562,
  sourceManifest: 95,
  source: 262_144,
  planRegistration: 4_096,
  registerPlanV0Request: 328_337,
  registerPlanV0Reply: 4_096,
  getRegisteredPlanV0Request: 16,
  getRegisteredPlanV0Reply: 332_433,
  submitMainMJSV0Request: 2_097_152,
  submitMainMJSV0Reply: 16,
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

const submitMethod = await json("submit-main-mjs-v0.method.json");
const registerMethod = await json("register-plan-v0.method.json");
const getMethod = await json("get-registered-plan-v0.method.json");
assert.deepEqual([submitMethod.method, submitMethod.methodVersion], ["SubmitMainMJSV0", 0]);
assert.equal(submitMethod.expectedRole, "internal-alpha-cli");
assert.equal(submitMethod.expectedSigningIdentifier, "com.capsulecorp.capsule.cli");
assert.equal(submitMethod.requestMediaType, "application/capsule.job-proposal+json;v=0");
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
for (const method of [submitMethod, registerMethod, getMethod]) {
  assert.equal(method.peerRequirementBeforeDeliveryRequired, true);
  assert.equal(method.methodAuthorityDisposition, "service-role-purpose-audience-derived");
  assert.equal(method.requestIdAuthorityDisposition, "correlation-only");
  assert.equal(method.connectionMaxInFlight, 1);
  assert.equal(method.applicationQueueCapacity, 0);
}

const nativeXPC = await json("native-xpc-v0.contract.json");
assert.equal(nativeXPC.objectType, "capsule.authenticated-local-ipc-native-xpc-contract");
assert.equal(nativeXPC.status, "passive-unwired-no-listener");
assert.equal(nativeXPC.transportEncoding, "xpc-dictionary-v0");
assert.equal(nativeXPC.transportEncodingVersion, 0);
assert.equal(nativeXPC.topLevelObjectType, "XPC_TYPE_DICTIONARY");
assert.deepEqual(nativeXPC.messageTags, manifest.numericMessageTags);
assert.deepEqual(nativeXPC.statusTags, manifest.numericStatusTags);
assert.deepEqual(nativeXPC.reasonTags, manifest.numericReasonTags);
assert.equal(nativeXPC.extraObjectsAllowed, 0);
assert.equal(nativeXPC.fileDescriptorsAllowed, 0);
assert.equal(nativeXPC.endpointsAllowed, 0);
assert.equal(nativeXPC.machRightsAllowed, 0);
assert.equal(nativeXPC.nestedContainersAllowed, 0);
assert.equal(nativeXPC.listenerActivated, false);
assert.equal(nativeXPC.serviceRegistered, false);
assert.equal(nativeXPC.peerAuthenticationEvidence, null);
assert.equal(nativeXPC.messageTagDisposition, "method-specific-cross-check-not-dispatch-opcode");
assert.equal(nativeXPC.requestIdDisposition, "correlation-only");
assert.deepEqual(nativeXPC.validationPrecedence, [
  "protocol-version",
  "method-version",
  "service-entry-point-role-message-tag-audience-purpose",
  "installation-epoch-current-state",
  "application-data-copy",
  "embedded-record-version-and-core-validation",
]);
assert.deepEqual(Object.keys(nativeXPC.envelopes), [
  "SubmitMainMJSV0",
  "RegisterPlanV0",
  "GetRegisteredPlanV0",
]);
assert.deepEqual(nativeXPC.methodBindings.SubmitMainMJSV0, {
  service: submitMethod.service,
  expectedRole: submitMethod.expectedRole,
  expectedSigningIdentifier: submitMethod.expectedSigningIdentifier,
  audience: submitMethod.audience,
  purpose: submitMethod.purpose,
  messageTag: nativeXPC.messageTags.SubmitMainMJSV0,
  methodVersion: 0,
});
assert.deepEqual(nativeXPC.methodBindings.RegisterPlanV0, {
  service: registerMethod.service,
  expectedRole: registerMethod.expectedRole,
  expectedSigningIdentifier: null,
  audience: registerMethod.audience,
  purpose: registerMethod.purpose,
  messageTag: nativeXPC.messageTags.RegisterPlanV0,
  methodVersion: 0,
});
assert.deepEqual(nativeXPC.methodBindings.GetRegisteredPlanV0, {
  service: getMethod.service,
  expectedRole: getMethod.expectedRole,
  expectedSigningIdentifier: null,
  audience: getMethod.audience,
  purpose: getMethod.purpose,
  messageTag: nativeXPC.messageTags.GetRegisteredPlanV0,
  methodVersion: 0,
});
for (const [method, envelopes] of Object.entries(nativeXPC.envelopes)) {
  for (const envelope of Object.values(envelopes)) {
    assert.equal(envelope.method, method);
    assert.equal(envelope.fields.length, envelope.exactKeyCount);
    assert.equal(new Set(envelope.fields.map((field) => field.key)).size, envelope.exactKeyCount);
    for (const field of envelope.fields) {
      assert.ok(Object.values(nativeXPC.keys).includes(field.key), `${method} ${field.key}`);
      assert.ok(
        Object.values(nativeXPC.valueTypes).includes(field.valueType),
        `${method} ${field.valueType}`,
      );
    }
  }
  assert.equal(envelopes.request.fields[0].key, nativeXPC.keys.protocolVersion);
  assert.equal(envelopes.request.fields[8].key, nativeXPC.keys.purpose);
  assert.equal(envelopes.successReply.fields[4].fixedUInt64, nativeXPC.statusTags.OK);
  assert.equal(envelopes.successReply.fields[5].fixedUInt64, nativeXPC.reasonTags.none);
  assert.equal(envelopes.refusalReply.exactKeyCount, 6);
  assert.equal(envelopes.refusalReply.applicationDataMaxBytes, 0);
}
assert.equal(nativeXPC.envelopes.SubmitMainMJSV0.request.exactKeyCount, 10);
assert.equal(nativeXPC.envelopes.RegisterPlanV0.request.exactKeyCount, 13);
assert.equal(nativeXPC.envelopes.GetRegisteredPlanV0.request.exactKeyCount, 10);
assert.equal(nativeXPC.envelopes.SubmitMainMJSV0.successReply.exactKeyCount, 7);
assert.equal(nativeXPC.envelopes.RegisterPlanV0.successReply.exactKeyCount, 7);
assert.equal(nativeXPC.envelopes.GetRegisteredPlanV0.successReply.exactKeyCount, 11);
assert.equal(manifest.nativeXPCEnvelopeCount, 9);
assert.equal(manifest.nativeXPCCaseCount, nativeXPC.cases.length);
assert.equal(manifest.nativeXPCRefusalReplyCount, nativeXPC.refusalReplies.length);
for (const [classification, statusTag] of Object.entries(nativeXPC.classificationToStatus)) {
  assert.equal(statusTag, nativeXPC.statusTags[classification]);
  assert.notEqual(statusTag, nativeXPC.statusTags.OK);
}
assert.deepEqual(nativeXPC.structuralReasonToStatus, {
  keySet: nativeXPC.statusTags.MALFORMED,
  valueType: nativeXPC.statusTags.MALFORMED,
  dataWidth: nativeXPC.statusTags.SCHEMA,
  dataCap: nativeXPC.statusTags.MALFORMED,
  zeroIdentifier: nativeXPC.statusTags.SCHEMA,
  epochSequence: nativeXPC.statusTags.SCHEMA,
  protocolVersion: nativeXPC.statusTags.UNSUPPORTED,
  methodVersion: nativeXPC.statusTags.UNSUPPORTED,
  messageTag: nativeXPC.statusTags.UNSUPPORTED,
  methodBinding: nativeXPC.statusTags.AUTHENTICATION,
  currentState: nativeXPC.statusTags.BINDING,
  capacity: nativeXPC.statusTags.CAPACITY,
});
assert.equal(nativeXPC.coreRefusalMapping, "classification-selects-status;reason=coreRefusal");
assert.equal(
  nativeXPC.localIntegrityMapping,
  "terminate-without-reply;reason-tag-is-fixture-diagnostic-only",
);
assert.equal(nativeXPC.refusalReplies.length, 13);
for (const refusalReply of nativeXPC.refusalReplies) {
  assert.equal(refusalReply.statusTag, nativeXPC.statusTags[refusalReply.classification]);
  assert.equal(refusalReply.reasonTag, nativeXPC.reasonTags.coreRefusal);
  assert.equal(refusalReply.bodyKeysAllowed, 0);
  assert.equal(refusalReply.exactKeyCount, 6);
}
for (const candidate of nativeXPC.cases.filter((entry) => entry.noState)) {
  assertZeroState(candidate.noState, candidate.id);
  assert.equal(candidate.expected.bodyCopied, false, candidate.id);
  assert.equal(candidate.expected.coreCalls, 0, candidate.id);
}
assert.deepEqual(
  nativeXPC.cases.map((candidate) => candidate.id),
  [
    "all.exact-key-and-type-sets",
    "submit.body-cap-plus-one",
    "register.body-cap-plus-one",
    "get.body-cap-plus-one",
    "all.missing-key",
    "all.extra-key",
    "all.wrong-type",
    "all.nested-container",
    "all.zero-request-id",
    "all.unknown-protocol-version",
    "all.unknown-method-version",
    "register.joint-protocol-method-record-version-mismatch",
    "register.joint-method-record-version-mismatch",
    "register.embedded-record-version-mismatch",
    "cross-service.SubmitMainMJSV0.tag-RegisterPlanV0",
    "cross-service.SubmitMainMJSV0.tag-GetRegisteredPlanV0",
    "cross-service.RegisterPlanV0.tag-SubmitMainMJSV0",
    "cross-service.RegisterPlanV0.tag-GetRegisteredPlanV0",
    "cross-service.GetRegisteredPlanV0.tag-SubmitMainMJSV0",
    "cross-service.GetRegisteredPlanV0.tag-RegisterPlanV0",
    "all.wrong-service",
    "all.wrong-role",
    "all.wrong-audience",
    "all.wrong-purpose",
    "all.local-audience-replaced-by-signed-object-audience",
    "all.local-purpose-replaced-by-signed-object-purpose",
    "all.request-id-replaced-by-installation-id-bytes",
    "all.installation-id-replaced-by-request-id-bytes",
    "get.registration-id-replaced-by-request-id-bytes",
    "all.wrong-installation",
    "all.wrong-epoch-digest",
    "all.epoch-uint53-cap-plus-one",
    "all.extra-file-descriptor",
    "all.extra-endpoint",
    "all.extra-mach-send-right",
    "all.success-reply-extra-key",
    "all.success-reply-request-id-mismatch",
    "all.refusal-reply-extra-body",
    "all.local-integrity-output-fault",
  ],
);
assert.deepEqual(nativeXPC.responseLoss, [
  {
    method: "SubmitMainMJSV0",
    disposition: "committed-retry-may-create-fresh-registration",
  },
  {
    method: "RegisterPlanV0",
    disposition: "committed-retry-creates-fresh-registration",
  },
  {
    method: "GetRegisteredPlanV0",
    disposition: "repeatable-read-by-registration-id",
  },
]);
assert.deepEqual(nativeXPC.futureNativeHarnessOracles, [
  {
    id: "os-peer-requirement-mismatch",
    scope: "future-external-native-harness-only",
    expected: "no-message-delivery-and-no-application-reply",
    currentEvidence: false,
    inBandRefusal: false,
  },
]);

const submitRequest = await json("submit-main-mjs-v0.request.json");
const submitReply = await json("submit-main-mjs-v0.reply.json");
const registerRequest = await json("register-plan-v0.request.json");
const registerReply = await json("register-plan-v0.reply.json");
const getRequest = await json("get-registered-plan-v0.request.json");
const getReply = await json("get-registered-plan-v0.reply.json");
for (const envelope of [
  submitRequest,
  submitReply,
  registerRequest,
  registerReply,
  getRequest,
  getReply,
]) {
  assert.equal(envelope.objectVersion, 0);
  assert.equal(envelope.fixtureSerialization, "exact-json-not-xpc-framing");
}
assert.equal(submitRequest.requestId, submitReply.requestId);
assert.equal(submitReply.applicationDataBytes, 16);
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
for (const candidate of oracles.flowControl) assertZeroState(candidate.noState, candidate.id);
assert.deepEqual(
  oracles.flowControl.find(
    (candidate) => candidate.id === "supervisor.mixed-register-get-aggregate-byte-cap-plus-one",
  ),
  {
    id: "supervisor.mixed-register-get-aggregate-byte-cap-plus-one",
    methods: [
      { method: "RegisterPlanV0", admittedBytes: 328_337 },
      { method: "GetRegisteredPlanV0", admittedBytes: 16 },
    ],
    isolatedAccountingCapBytes: 328_353,
    attemptedAdditionalBytes: 1,
    classification: "CAPACITY",
    reason: "ipc-process-in-flight-bytes",
    releasedMethod: "GetRegisteredPlanV0",
    postReleaseReadmissionBytes: 1,
    postReleaseReadmission: "admitted",
    productionAggregateCapBytes: 2_626_696,
    queueDepth: 0,
    noState: {
      authorityStateChanged: false,
      coreCalls: 0,
      registrationsCreated: 0,
      approvalsConsumed: 0,
      attemptsCreated: 0,
      lifecycleCalls: 0,
      backendCalls: 0,
    },
  },
);
assertZeroState(oracles.cancellationAndDeadline[0].noState, "cancel.before-dispatch");
assert.equal(oracles.cancellationAndDeadline[1].admittedSlotHeldUntilCoreReturns, true);
assert.equal(oracles.cancellationAndDeadline[2].admittedSlotHeldUntilCoreReturns, true);
assert.equal(oracles.sourceCustody.exactProposalSourceEqualsRegisteredSource, true);
assert.equal(oracles.sourceCustody.executeTimeReplacementBytesAccepted, false);
assert.equal(
  oracles.sourceCustody.extractedMainMJS.sha256,
  oracles.sourceCustody.registeredSourceBytes.sha256,
);
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
assert.equal(oracles.responseLoss[0].retryMayCreateFreshRegistration, true);
assert.equal(oracles.responseLoss[1].requestIdIsIdempotencyKey, false);
assert.equal(oracles.responseLoss[1].retryCommittedRegistrations, 2);
assert.equal(oracles.responseLoss[1].bothRegistrationsSeparatelyReadable, true);
assert.equal(oracles.responseLoss[2].requestIdIsIdempotencyKey, false);
assert.equal(oracles.responseLoss[2].repeatedReplyBodyByteEqual, true);

const combined = JSON.stringify({
  manifest,
  submitMethod,
  registerMethod,
  getMethod,
  registerRequest,
  registerReply,
  getRequest,
  getReply,
  submitRequest,
  submitReply,
  oracles,
});
assert.equal(combined.includes("RegisterPlanV1"), false);
assert.equal(combined.includes("GetRegisteredPlanV1"), false);
assert.equal(combined.includes("typescript"), false);
assert.equal(combined.includes("626-byte"), false);

process.stdout.write(
  "verified independent TypeScript passive authenticated-local-IPC known answers\n",
);

async function json(path) {
  return JSON.parse(await readFile(new URL(path, root), "utf8"));
}

function verifyReference(bytes, known, label) {
  assert.equal(bytes.length, known.byteLength, `${label} byte length`);
  assert.equal(sha256Hex(bytes), known.sha256, `${label} digest`);
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
