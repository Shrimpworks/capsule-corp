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
assert.equal(manifest.methodCount, 5);
assert.deepEqual(manifest.numericMessageTags, {
  invalid: 0,
  SubmitMainMJSV0: 1,
  RegisterPlanV0: 2,
  GetRegisteredPlanV0: 3,
  SubmitApprovalV0: 4,
  RequestAttemptV0: 5,
});
assert.deepEqual(manifest.numericApprovalStateTags, {
  invalid: 0,
  usable: 1,
  consumed: 2,
  invalidated: 3,
});
assert.deepEqual(manifest.numericAttemptStateTags, { invalid: 0, created: 1 });
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
  approvalEnvelope: 512,
  registerPlanV0Request: 328_337,
  registerPlanV0Reply: 4_096,
  getRegisteredPlanV0Request: 16,
  getRegisteredPlanV0Reply: 332_433,
  submitMainMJSV0Request: 2_097_152,
  submitMainMJSV0Reply: 16,
  submitApprovalV0Request: 528,
  submitApprovalV0Reply: 16,
  requestAttemptV0Request: 32,
  requestAttemptV0Reply: 16,
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
assert.equal(manifest.caps.submitApprovalV0Request, 16 + manifest.caps.approvalEnvelope);
assert.equal(manifest.caps.requestAttemptV0Request, 16 + 16);

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
assert.deepEqual(nativeXPC.approvalStateTags, manifest.numericApprovalStateTags);
assert.deepEqual(nativeXPC.attemptStateTags, manifest.numericAttemptStateTags);
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
  "SubmitApprovalV0",
  "RequestAttemptV0",
]);
assert.deepEqual(nativeXPC.methodBindings.SubmitMainMJSV0, {
  entryPoint: "SubmitMainMJSV0",
  service: submitMethod.service,
  expectedRole: submitMethod.expectedRole,
  expectedSigningIdentifier: submitMethod.expectedSigningIdentifier,
  audience: submitMethod.audience,
  purpose: submitMethod.purpose,
  messageTag: nativeXPC.messageTags.SubmitMainMJSV0,
  methodVersion: 0,
  deadlineMilliseconds: 10_000,
});
assert.deepEqual(nativeXPC.methodBindings.RegisterPlanV0, {
  entryPoint: "RegisterPlanV0",
  service: registerMethod.service,
  expectedRole: registerMethod.expectedRole,
  expectedSigningIdentifier: null,
  audience: registerMethod.audience,
  purpose: registerMethod.purpose,
  messageTag: nativeXPC.messageTags.RegisterPlanV0,
  methodVersion: 0,
  deadlineMilliseconds: 5_000,
});
assert.deepEqual(nativeXPC.methodBindings.GetRegisteredPlanV0, {
  entryPoint: "GetRegisteredPlanV0",
  service: getMethod.service,
  expectedRole: getMethod.expectedRole,
  expectedSigningIdentifier: null,
  audience: getMethod.audience,
  purpose: getMethod.purpose,
  messageTag: nativeXPC.messageTags.GetRegisteredPlanV0,
  methodVersion: 0,
  deadlineMilliseconds: 2_000,
});
assert.deepEqual(nativeXPC.methodBindings.SubmitApprovalV0, {
  entryPoint: "SubmitApprovalV0",
  service: "com.capsulecorp.capsule.supervisor.broker.v0",
  expectedRole: "broker",
  expectedSigningIdentifier: null,
  audience: "capsule.execution-supervisor.local.v0",
  purpose: "capsule.ipc.submit-approval.v0",
  messageTag: 4,
  methodVersion: 0,
  deadlineMilliseconds: 5_000,
});
assert.deepEqual(nativeXPC.methodBindings.RequestAttemptV0, {
  entryPoint: "RequestAttemptV0",
  service: "com.capsulecorp.capsule.supervisor.daemon.v0",
  expectedRole: "daemon",
  expectedSigningIdentifier: null,
  audience: "capsule.execution-supervisor.local.v0",
  purpose: "capsule.ipc.request-attempt.v0",
  messageTag: 5,
  methodVersion: 0,
  deadlineMilliseconds: 5_000,
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
for (const method of ["SubmitApprovalV0", "RequestAttemptV0"]) {
  for (const direction of ["request", "successReply"]) {
    const envelope = nativeXPC.envelopes[method][direction];
    const derivedMaximum = envelope.fields
      .filter((field) => field.applicationData === true)
      .reduce((total, field) => total + field.maxDataBytes, 0);
    assert.equal(
      envelope.applicationDataMaxBytes,
      derivedMaximum,
      `${method} ${direction} aggregate maximum`,
    );
  }
}
assert.equal(nativeXPC.envelopes.SubmitMainMJSV0.request.exactKeyCount, 10);
assert.equal(nativeXPC.envelopes.RegisterPlanV0.request.exactKeyCount, 13);
assert.equal(nativeXPC.envelopes.GetRegisteredPlanV0.request.exactKeyCount, 10);
assert.equal(nativeXPC.envelopes.SubmitApprovalV0.request.exactKeyCount, 11);
assert.equal(nativeXPC.envelopes.RequestAttemptV0.request.exactKeyCount, 11);
assert.equal(nativeXPC.envelopes.SubmitMainMJSV0.successReply.exactKeyCount, 7);
assert.equal(nativeXPC.envelopes.RegisterPlanV0.successReply.exactKeyCount, 7);
assert.equal(nativeXPC.envelopes.GetRegisteredPlanV0.successReply.exactKeyCount, 11);
assert.equal(nativeXPC.envelopes.SubmitApprovalV0.successReply.exactKeyCount, 8);
assert.equal(nativeXPC.envelopes.RequestAttemptV0.successReply.exactKeyCount, 8);
assert.deepEqual(
  nativeXPC.envelopes.SubmitApprovalV0.successReply.fields.at(-1).allowedUInt64,
  [1, 2, 3],
);
assert.deepEqual(
  nativeXPC.envelopes.RequestAttemptV0.successReply.fields.at(-1).allowedUInt64,
  [1],
);
assert.equal(manifest.nativeXPCEnvelopeCount, 15);
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
assert.deepEqual(nativeXPC.signedObjectBindings, {
  SubmitApprovalV0: {
    outerField: "capsule.approval-envelope",
    encoding: "tagged-canonical-cose-sign1-with-embedded-canonical-cbor-payload",
    objectType: "capsule.approval-grant",
    objectVersion: 0,
    audience: "capsule.execution-supervisor",
    purpose: "capsule.plan.approve",
    localAudience: "capsule.execution-supervisor.local.v0",
    localPurpose: "capsule.ipc.submit-approval.v0",
    localAndSignedBindingsInterchangeable: false,
  },
});
assert.deepEqual(nativeXPC.identifierDomains, {
  "capsule.request-id": "live-call-correlation-only",
  "capsule.installation-id": "current-installation-identity",
  "capsule.registration-id": "registration-lookup",
  "capsule.approval-id": "approval-reference",
  "capsule.attempt-id": "attempt-reference",
  sameWidthBytesInterchangeable: false,
});
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
assert.deepEqual(nativeXPC.caseTable, expectedNativeXPCCaseTable(nativeXPC));
assert.deepEqual(nativeXPC.cases.map(caseTableEntry), nativeXPC.caseTable);
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
  {
    method: "SubmitApprovalV0",
    disposition:
      "canonical-payload-and-signer-authorization-replay-returns-same-approval-and-current-state",
    semanticIdentity: "canonical-approval-payload+resolved-signer-authorization-identity",
    sameApprovalID: true,
    currentStateReturned: true,
    duplicateAuthorityEffects: 0,
    requestIdIsIdempotencyKey: false,
  },
  {
    method: "RequestAttemptV0",
    disposition:
      "registration-and-approval-reference-replay-returns-same-attempt-and-current-state",
    semanticIdentity: "registration-id+approval-reference",
    sameAttemptID: true,
    currentStateReturned: true,
    duplicateAttempts: 0,
    duplicateLifecycleEffects: 0,
    requestIdIsIdempotencyKey: false,
  },
]);
assert.deepEqual(nativeXPC.deadlineCases, [
  {
    method: "SubmitApprovalV0",
    deadlineMilliseconds: 5_000,
    startsAt: "admission",
    clientExtensionAllowed: false,
    afterDispatchDisposition: "response-unknown-store-semantic-result-or-recovery-fence-controls",
  },
  {
    method: "RequestAttemptV0",
    deadlineMilliseconds: 5_000,
    startsAt: "admission",
    clientExtensionAllowed: false,
    afterDispatchDisposition: "response-unknown-store-semantic-result-or-recovery-fence-controls",
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
assert.deepEqual(
  oracles.flowControl.find(
    (candidate) => candidate.id === "supervisor.mixed-submit-approval-request-attempt-cap-plus-one",
  ),
  {
    id: "supervisor.mixed-submit-approval-request-attempt-cap-plus-one",
    methods: [
      { method: "SubmitApprovalV0", admittedBytes: 528 },
      { method: "RequestAttemptV0", admittedBytes: 32 },
    ],
    isolatedAccountingCapBytes: 560,
    attemptedAdditionalBytes: 1,
    classification: "CAPACITY",
    reason: "ipc-process-in-flight-bytes",
    releasedMethod: "RequestAttemptV0",
    releasedBytes: 32,
    postReleaseReadmissionBytes: 32,
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
assert.deepEqual(oracles.responseLoss.slice(3), nativeXPC.responseLoss.slice(3));

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

function expectedNativeXPCCaseTable(contract) {
  const status = contract.statusTags;
  const reason = contract.reasonTags;
  const row = (
    id,
    method,
    direction,
    mutation,
    decision,
    classification = null,
    reasonTag = null,
    bodyCopied = null,
    coreCalls = null,
  ) => ({
    id,
    method,
    direction,
    mutation,
    decision,
    classification,
    statusTag: classification === null ? null : status[classification],
    reasonTag,
    bodyCopied,
    coreCalls,
  });
  const reject = (id, method, mutation, classification, reasonTag) =>
    row(
      id,
      method,
      "request",
      mutation,
      "reject-before-body-copy",
      classification,
      reasonTag,
      false,
      0,
    );
  const core = (id, method, mutation, classification) =>
    row(
      id,
      method,
      "request",
      mutation,
      "core-refusal-after-bounded-body-copy",
      classification,
      reason.coreRefusal,
      true,
      1,
    );
  const result = [
    row("all.exact-key-and-type-sets", "all-frozen-methods", null, null, "accept-outer-shape-only"),
    row(
      "submit-approval.request.exact-maximum",
      "SubmitApprovalV0",
      "request",
      "registration-id=16-bytes;approval-envelope=512-bytes",
      "accept-outer-shape-and-bounded-body-copy",
      null,
      null,
      true,
      1,
    ),
    row(
      "request-attempt.request.exact-maximum",
      "RequestAttemptV0",
      "request",
      "registration-id=16-bytes;approval-id=16-bytes",
      "accept-outer-shape-and-bounded-body-copy",
      null,
      null,
      true,
      1,
    ),
    reject(
      "submit.body-cap-plus-one",
      "SubmitMainMJSV0",
      "job-proposal=2097153-bytes",
      "MALFORMED",
      reason.dataCap,
    ),
    reject(
      "register.body-cap-plus-one",
      "RegisterPlanV0",
      "source=262145-bytes",
      "MALFORMED",
      reason.dataCap,
    ),
    reject(
      "get.body-cap-plus-one",
      "GetRegisteredPlanV0",
      "registration-id=17-bytes",
      "SCHEMA",
      reason.dataWidth,
    ),
    reject(
      "submit-approval.body-cap-plus-one",
      "SubmitApprovalV0",
      "approval-envelope=513-bytes",
      "MALFORMED",
      reason.dataCap,
    ),
    reject(
      "request-attempt.approval-id-width-plus-one",
      "RequestAttemptV0",
      "approval-id=17-bytes",
      "SCHEMA",
      reason.dataWidth,
    ),
    reject("all.missing-key", "all", "remove-purpose", "MALFORMED", reason.keySet),
    reject("all.extra-key", "all", "add-capsule.extra", "MALFORMED", reason.keySet),
    reject("all.wrong-type", "all", "request-id=XPC_TYPE_STRING", "MALFORMED", reason.valueType),
    reject(
      "all.nested-container",
      "all",
      "body=XPC_TYPE_DICTIONARY",
      "MALFORMED",
      reason.valueType,
    ),
    reject(
      "all.zero-request-id",
      "all",
      "request-id=16-zero-bytes",
      "SCHEMA",
      reason.zeroIdentifier,
    ),
    reject(
      "all.unknown-protocol-version",
      "all",
      "protocol-version=1",
      "UNSUPPORTED",
      reason.protocolVersion,
    ),
    reject(
      "all.unknown-method-version",
      "all",
      "method-version=1",
      "UNSUPPORTED",
      reason.methodVersion,
    ),
    reject(
      "register.joint-protocol-method-record-version-mismatch",
      "RegisterPlanV0",
      "protocol-version=1;method-version=1;role-bindings[0]=1",
      "UNSUPPORTED",
      reason.protocolVersion,
    ),
    reject(
      "register.joint-method-record-version-mismatch",
      "RegisterPlanV0",
      "method-version=1;role-bindings[0]=1",
      "UNSUPPORTED",
      reason.methodVersion,
    ),
    core(
      "register.embedded-record-version-mismatch",
      "RegisterPlanV0",
      "role-bindings[0]=1",
      "UNSUPPORTED",
    ),
    reject(
      "submit-approval.joint-protocol-method-record-version-mismatch",
      "SubmitApprovalV0",
      "protocol-version=1;method-version=1;approval-grant.object-version=1",
      "UNSUPPORTED",
      reason.protocolVersion,
    ),
    reject(
      "submit-approval.joint-method-record-version-mismatch",
      "SubmitApprovalV0",
      "method-version=1;approval-grant.object-version=1",
      "UNSUPPORTED",
      reason.methodVersion,
    ),
    core(
      "submit-approval.embedded-record-version-mismatch",
      "SubmitApprovalV0",
      "approval-grant.object-version=1-with-valid-test-signature",
      "UNSUPPORTED",
    ),
  ];
  const methods = [
    "SubmitMainMJSV0",
    "RegisterPlanV0",
    "GetRegisteredPlanV0",
    "SubmitApprovalV0",
    "RequestAttemptV0",
  ];
  for (const selected of methods) {
    for (const received of methods) {
      if (selected === received) continue;
      result.push(
        reject(
          `cross-service.${selected}.tag-${received}`,
          selected,
          `entry-point=${selected};message-tag=${received}`,
          "UNSUPPORTED",
          reason.messageTag,
        ),
      );
    }
  }
  result.push(
    reject(
      "all.wrong-service",
      "all",
      "selected-method-invoked-with-receiving-service=other-role-service",
      "AUTHENTICATION",
      reason.methodBinding,
    ),
    reject(
      "all.wrong-role",
      "all",
      "message-derived-role=other",
      "AUTHENTICATION",
      reason.methodBinding,
    ),
    reject(
      "all.wrong-session",
      "all",
      "message-derived-euid-or-audit-session=other",
      "AUTHENTICATION",
      reason.methodBinding,
    ),
    reject("all.wrong-audience", "all", "audience=other", "AUTHENTICATION", reason.methodBinding),
    reject("all.wrong-purpose", "all", "purpose=other", "AUTHENTICATION", reason.methodBinding),
    reject(
      "all.local-audience-replaced-by-signed-object-audience",
      "all",
      "audience=capsule.execution-supervisor",
      "AUTHENTICATION",
      reason.methodBinding,
    ),
    reject(
      "all.local-purpose-replaced-by-signed-object-purpose",
      "all",
      "purpose=capsule.plan.approve",
      "AUTHENTICATION",
      reason.methodBinding,
    ),
    core(
      "submit-approval.signed-audience-replaced-by-local-channel-audience",
      "SubmitApprovalV0",
      "approval-grant.audience=capsule.execution-supervisor.local.v0-with-valid-test-signature",
      "BINDING",
    ),
    core(
      "submit-approval.signed-purpose-replaced-by-local-channel-purpose",
      "SubmitApprovalV0",
      "approval-grant.purpose=capsule.ipc.submit-approval.v0-with-valid-test-signature",
      "BINDING",
    ),
    row(
      "all.request-id-replaced-by-installation-id-bytes",
      "all",
      "request",
      "request-id=installation-id-bytes",
      "accept-correlation-field-only",
    ),
    reject(
      "all.installation-id-replaced-by-request-id-bytes",
      "all",
      "installation-id=request-id-bytes",
      "BINDING",
      reason.currentState,
    ),
    row(
      "get.registration-id-replaced-by-request-id-bytes",
      "GetRegisteredPlanV0",
      "request",
      "registration-id=request-id-bytes",
      "continue-to-registration-lookup-without-transport-authority",
    ),
    core(
      "submit-approval.registration-id-replaced-by-request-id-bytes",
      "SubmitApprovalV0",
      "registration-id=request-id-bytes",
      "BINDING",
    ),
    core(
      "request-attempt.registration-id-replaced-by-approval-id-bytes",
      "RequestAttemptV0",
      "registration-id=approval-id-bytes",
      "BINDING",
    ),
    core(
      "request-attempt.approval-id-replaced-by-registration-id-bytes",
      "RequestAttemptV0",
      "approval-id=registration-id-bytes",
      "BINDING",
    ),
    reject(
      "all.wrong-installation",
      "all",
      "installation-id=other",
      "BINDING",
      reason.currentState,
    ),
    reject("all.wrong-epoch-digest", "all", "epoch-digest=other", "BINDING", reason.currentState),
    reject(
      "all.epoch-uint53-cap-plus-one",
      "all",
      "epoch-sequence=9007199254740992",
      "SCHEMA",
      reason.epochSequence,
    ),
    reject(
      "all.extra-file-descriptor",
      "all",
      "add-capsule.smuggled-fd=XPC_TYPE_FD",
      "MALFORMED",
      reason.keySet,
    ),
    reject(
      "all.extra-endpoint",
      "all",
      "add-capsule.smuggled-endpoint=XPC_TYPE_ENDPOINT",
      "MALFORMED",
      reason.keySet,
    ),
    reject(
      "all.extra-mach-send-right",
      "all",
      "add-capsule.smuggled-mach-send-right=XPC_TYPE_MACH_SEND",
      "MALFORMED",
      reason.keySet,
    ),
    row(
      "all.success-reply-extra-key",
      "all-frozen-methods",
      "success-reply",
      "add-any-unknown-key",
      "reject-reply",
      "MALFORMED",
      reason.keySet,
    ),
    row(
      "all.success-reply-request-id-mismatch",
      "all-frozen-methods",
      "success-reply",
      "request-id=other-live-call",
      "reject-reply",
      "BINDING",
      reason.currentState,
    ),
    row(
      "all.refusal-reply-extra-body",
      "all-frozen-methods",
      "refusal-reply",
      "add-any-method-body-key",
      "reject-reply",
      "MALFORMED",
      reason.keySet,
    ),
    row(
      "all.local-integrity-output-fault",
      "all-frozen-methods",
      "success-reply",
      "oversize-short-write-pointer-length-or-bridge-version-mismatch",
      "terminate-without-reply",
      null,
      reason.localIntegrityFault,
    ),
    row(
      "submit-approval.cancel-before-dispatch",
      "SubmitApprovalV0",
      "request",
      "caller-cancel-before-bridge-dispatch",
      "no-core-call-no-application-reply",
      null,
      null,
      false,
      0,
    ),
    row(
      "submit-approval.cancel-after-dispatch",
      "SubmitApprovalV0",
      "request",
      "caller-cancel-after-bridge-dispatch",
      "response-delivery-cancelled-store-semantic-result-controls",
      null,
      null,
      true,
      1,
    ),
    row(
      "request-attempt.cancel-before-dispatch",
      "RequestAttemptV0",
      "request",
      "caller-cancel-before-bridge-dispatch",
      "no-core-call-no-application-reply",
      null,
      null,
      false,
      0,
    ),
    row(
      "request-attempt.cancel-after-dispatch",
      "RequestAttemptV0",
      "request",
      "caller-cancel-after-bridge-dispatch",
      "response-delivery-cancelled-store-semantic-result-controls",
      null,
      null,
      true,
      1,
    ),
  );
  return result;
}

function caseTableEntry(candidate) {
  const expected =
    typeof candidate.expected === "string" ? { decision: candidate.expected } : candidate.expected;
  return {
    id: candidate.id,
    method: candidate.method ?? candidate.methods,
    direction: candidate.direction ?? null,
    mutation: candidate.mutation ?? null,
    decision: expected.decision,
    classification: expected.classification ?? null,
    statusTag: expected.statusTag ?? null,
    reasonTag: expected.reasonTag ?? null,
    bodyCopied: expected.bodyCopied ?? null,
    coreCalls: expected.coreCalls ?? null,
  };
}

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
