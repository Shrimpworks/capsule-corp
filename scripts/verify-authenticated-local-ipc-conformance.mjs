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
  if (candidate.id === "all.exact-key-and-type-sets") {
    assert.equal(candidate.expected.bodyCopied, null, candidate.id);
    assert.equal(candidate.expected.coreCalls, null, candidate.id);
  } else {
    assert.equal(candidate.expected.bodyCopied, false, candidate.id);
    assert.equal(candidate.expected.coreCalls, 0, candidate.id);
  }
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
assert.equal(
  nativeXPC.deadlineBoundaryRule,
  "dispatch-only-when-elapsed-milliseconds-is-strictly-less-than-deadline-milliseconds",
);
assert.deepEqual(nativeXPC.deadlineCases, expectedDeadlineCases());
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
verifyHardenedNativeXPC(nativeXPC, oracles);
runHardenedMutationProofs(nativeXPC, oracles);

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

process.stdout.write("verified independent Node passive authenticated-local-IPC known answers\n");

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
  return result.map((entry) => ({ ...entry, ...expectedNativeXPCCaseControls(entry) }));
}

function expectedNativeXPCCaseControls(entry) {
  let replyDisposition = "core-semantic-reply";
  if (
    entry.decision === "reject-before-body-copy" ||
    entry.decision === "core-refusal-after-bounded-body-copy"
  )
    replyDisposition = "application-refusal-reply";
  else if (entry.decision === "reject-reply") replyDisposition = "received-reply-rejected";
  else if (
    entry.decision === "terminate-without-reply" ||
    entry.decision === "no-core-call-no-application-reply" ||
    entry.decision === "response-delivery-cancelled-store-semantic-result-controls"
  )
    replyDisposition = "no-application-reply";
  else if (entry.id === "all.exact-key-and-type-sets")
    replyDisposition = "not-applicable-passive-shape-check";
  if (["all.wrong-service", "all.wrong-role"].includes(entry.id))
    replyDisposition = "delivered-in-band-refusal";
  if (entry.id === "all.wrong-session")
    replyDisposition = "delivered-in-band-refusal-only-for-this-passive-oracle";

  let postCoreState = null;
  if (entry.decision === "core-refusal-after-bounded-body-copy")
    postCoreState = "core-refusal-no-authority-effect";
  else if (
    entry.decision === "response-delivery-cancelled-store-semantic-result-controls" ||
    entry.decision === "terminate-without-reply"
  )
    postCoreState = "store-semantic-result-or-recovery-fence-controls";
  else if (
    entry.coreCalls === 1 ||
    entry.id === "all.request-id-replaced-by-installation-id-bytes" ||
    entry.id === "get.registration-id-replaced-by-request-id-bytes"
  )
    postCoreState = "core-semantic-result-controls";

  const requiresNoState =
    entry.decision === "reject-before-body-copy" ||
    entry.decision === "no-core-call-no-application-reply" ||
    entry.decision === "accept-outer-shape-only";
  return {
    replyDisposition,
    terminationDisposition:
      entry.decision === "terminate-without-reply" ? "terminate-process" : "continue-process",
    noState: requiresNoState ? zeroStateValue() : null,
    postCoreState,
  };
}

function caseTableEntry(candidate) {
  const expected = candidate.expected;
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
    replyDisposition: candidate.replyDisposition,
    terminationDisposition: candidate.terminationDisposition,
    noState: candidate.noState,
    postCoreState: candidate.postCoreState,
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

function verifyHardenedNativeXPC(contract, oracleSet) {
  assertExactKeys(
    contract,
    [
      "objectType",
      "objectVersion",
      "status",
      "transportEncoding",
      "transportEncodingVersion",
      "fixtureSerialization",
      "topLevelObjectType",
      "valueTypes",
      "keys",
      "messageTags",
      "approvalStateTags",
      "attemptStateTags",
      "statusTags",
      "reasonTags",
      "classificationToStatus",
      "structuralReasonToStatus",
      "coreRefusalMapping",
      "localIntegrityMapping",
      "methodBindings",
      "signedObjectBindings",
      "identifierDomains",
      "requestCommonKeyCount",
      "replyCommonKeyCount",
      "extraObjectsAllowed",
      "fileDescriptorsAllowed",
      "endpointsAllowed",
      "machRightsAllowed",
      "nestedContainersAllowed",
      "messageTagDisposition",
      "validationPrecedence",
      "requestIdDisposition",
      "copyDisposition",
      "localIntegrityDisposition",
      "envelopes",
      "cases",
      "caseTable",
      "refusalReplies",
      "responseLoss",
      "deadlineBoundaryRule",
      "deadlineCases",
      "futureNativeHarnessOracles",
      "peerAuthenticationEvidence",
      "listenerActivated",
      "serviceRegistered",
    ],
    "native XPC contract",
  );
  assert.deepEqual(contract.valueTypes, {
    uint64: "XPC_TYPE_UINT64",
    data: "XPC_TYPE_DATA",
    string: "XPC_TYPE_STRING",
  });
  assert.deepEqual(contract.keys, {
    protocolVersion: "capsule.protocol-version",
    methodVersion: "capsule.method-version",
    messageTag: "capsule.message-tag",
    requestId: "capsule.request-id",
    installationId: "capsule.installation-id",
    epochSequence: "capsule.epoch-sequence",
    epochDigest: "capsule.epoch-digest",
    audience: "capsule.audience",
    purpose: "capsule.purpose",
    status: "capsule.status",
    reason: "capsule.reason",
    jobProposal: "capsule.job-proposal",
    executionPlan: "capsule.execution-plan",
    roleBindings: "capsule.role-bindings",
    sourceManifest: "capsule.source-manifest",
    source: "capsule.source",
    registrationId: "capsule.registration-id",
    planRegistration: "capsule.plan-registration",
    approvalEnvelope: "capsule.approval-envelope",
    approvalId: "capsule.approval-id",
    approvalState: "capsule.approval-state",
    attemptId: "capsule.attempt-id",
    attemptState: "capsule.attempt-state",
  });
  assert.deepEqual(contract.messageTags, {
    invalid: 0,
    SubmitMainMJSV0: 1,
    RegisterPlanV0: 2,
    GetRegisteredPlanV0: 3,
    SubmitApprovalV0: 4,
    RequestAttemptV0: 5,
  });
  assert.deepEqual(contract.approvalStateTags, {
    invalid: 0,
    usable: 1,
    consumed: 2,
    invalidated: 3,
  });
  assert.deepEqual(contract.attemptStateTags, { invalid: 0, created: 1 });
  assert.deepEqual(contract.statusTags, {
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
  assert.deepEqual(contract.reasonTags, {
    none: 0,
    keySet: 1,
    valueType: 2,
    dataWidth: 3,
    dataCap: 4,
    zeroIdentifier: 5,
    epochSequence: 6,
    protocolVersion: 7,
    methodVersion: 8,
    messageTag: 9,
    methodBinding: 10,
    currentState: 11,
    capacity: 12,
    coreRefusal: 13,
    localIntegrityFault: 14,
  });
  assert.deepEqual(contract.classificationToStatus, {
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
  assert.deepEqual(contract.structuralReasonToStatus, {
    keySet: 1,
    valueType: 1,
    dataWidth: 3,
    dataCap: 1,
    zeroIdentifier: 3,
    epochSequence: 3,
    protocolVersion: 2,
    methodVersion: 2,
    messageTag: 2,
    methodBinding: 5,
    currentState: 4,
    capacity: 8,
  });
  assert.deepEqual(contract.methodBindings, expectedNativeXPCMethodBindings());
  assert.deepEqual(contract.signedObjectBindings, {
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
  assert.deepEqual(contract.identifierDomains, {
    "capsule.request-id": "live-call-correlation-only",
    "capsule.installation-id": "current-installation-identity",
    "capsule.registration-id": "registration-lookup",
    "capsule.approval-id": "approval-reference",
    "capsule.attempt-id": "attempt-reference",
    sameWidthBytesInterchangeable: false,
  });
  assert.equal(contract.fixtureSerialization, "exact-json-description-of-xpc-dictionaries");
  assert.equal(contract.topLevelObjectType, "XPC_TYPE_DICTIONARY");
  assert.equal(contract.coreRefusalMapping, "classification-selects-status;reason=coreRefusal");
  assert.equal(
    contract.localIntegrityMapping,
    "terminate-without-reply;reason-tag-is-fixture-diagnostic-only",
  );
  assert.equal(contract.requestCommonKeyCount, 9);
  assert.equal(contract.replyCommonKeyCount, 6);
  assert.deepEqual(
    [
      contract.extraObjectsAllowed,
      contract.fileDescriptorsAllowed,
      contract.endpointsAllowed,
      contract.machRightsAllowed,
      contract.nestedContainersAllowed,
    ],
    [0, 0, 0, 0, 0],
  );
  assert.equal(contract.messageTagDisposition, "method-specific-cross-check-not-dispatch-opcode");
  assert.deepEqual(contract.validationPrecedence, [
    "protocol-version",
    "method-version",
    "service-entry-point-role-message-tag-audience-purpose",
    "installation-epoch-current-state",
    "application-data-copy",
    "embedded-record-version-and-core-validation",
  ]);
  assert.equal(contract.requestIdDisposition, "correlation-only");
  assert.equal(
    contract.copyDisposition,
    "body-copy-only-after-peer-flow-shape-current-state-binding",
  );
  assert.equal(
    contract.localIntegrityDisposition,
    "oversize-output-short-write-pointer-length-or-bridge-version-fault-terminates-process-without-reply",
  );
  assert.deepEqual(contract.envelopes, expectedNativeXPCEnvelopes());
  assert.deepEqual(contract.refusalReplies, expectedRefusalReplies(contract));
  verifyCompleteNativeXPCCases(contract);
  assert.deepEqual(contract.responseLoss, expectedResponseLoss());
  assert.deepEqual(contract.deadlineCases, expectedDeadlineCases());
  assert.deepEqual(contract.futureNativeHarnessOracles, [
    {
      id: "os-peer-requirement-mismatch",
      scope: "future-external-native-harness-only",
      expected: "no-message-delivery-and-no-application-reply",
      currentEvidence: false,
      inBandRefusal: false,
    },
  ]);
  assert.equal(contract.peerAuthenticationEvidence, null);
  assert.equal(contract.listenerActivated, false);
  assert.equal(contract.serviceRegistered, false);
  assert.deepEqual(oracleSet.responseLoss, expectedOracleResponseLoss());
  assert.deepEqual(oracleSet.refusals, expectedOracleRefusals());
  assert.deepEqual(oracleSet.cancellationAndDeadline, expectedCancellationAndDeadline());
  assert.deepEqual(oracleSet.flowControl, expectedFlowControl());
}

function expectedNativeXPCMethodBindings() {
  return {
    SubmitMainMJSV0: {
      entryPoint: "SubmitMainMJSV0",
      service: "com.capsulecorp.capsule.daemon.cli.v0",
      expectedRole: "internal-alpha-cli",
      expectedSigningIdentifier: "com.capsulecorp.capsule.cli",
      audience: "capsule.daemon.local.v0",
      purpose: "capsule.ipc.submit-main-mjs.v0",
      messageTag: 1,
      methodVersion: 0,
      deadlineMilliseconds: 10_000,
    },
    RegisterPlanV0: {
      entryPoint: "RegisterPlanV0",
      service: "com.capsulecorp.capsule.supervisor.daemon.v0",
      expectedRole: "daemon",
      expectedSigningIdentifier: null,
      audience: "capsule.execution-supervisor.local.v0",
      purpose: "capsule.ipc.register-plan.v0",
      messageTag: 2,
      methodVersion: 0,
      deadlineMilliseconds: 5_000,
    },
    GetRegisteredPlanV0: {
      entryPoint: "GetRegisteredPlanV0",
      service: "com.capsulecorp.capsule.supervisor.broker.v0",
      expectedRole: "broker",
      expectedSigningIdentifier: null,
      audience: "capsule.execution-supervisor.local.v0",
      purpose: "capsule.ipc.get-registered-plan.v0",
      messageTag: 3,
      methodVersion: 0,
      deadlineMilliseconds: 2_000,
    },
    SubmitApprovalV0: {
      entryPoint: "SubmitApprovalV0",
      service: "com.capsulecorp.capsule.supervisor.broker.v0",
      expectedRole: "broker",
      expectedSigningIdentifier: null,
      audience: "capsule.execution-supervisor.local.v0",
      purpose: "capsule.ipc.submit-approval.v0",
      messageTag: 4,
      methodVersion: 0,
      deadlineMilliseconds: 5_000,
    },
    RequestAttemptV0: {
      entryPoint: "RequestAttemptV0",
      service: "com.capsulecorp.capsule.supervisor.daemon.v0",
      expectedRole: "daemon",
      expectedSigningIdentifier: null,
      audience: "capsule.execution-supervisor.local.v0",
      purpose: "capsule.ipc.request-attempt.v0",
      messageTag: 5,
      methodVersion: 0,
      deadlineMilliseconds: 5_000,
    },
  };
}

function expectedNativeXPCEnvelopes() {
  const type = {
    uint64: "XPC_TYPE_UINT64",
    data: "XPC_TYPE_DATA",
    string: "XPC_TYPE_STRING",
  };
  const key = {
    protocol: "capsule.protocol-version",
    method: "capsule.method-version",
    tag: "capsule.message-tag",
    request: "capsule.request-id",
    installation: "capsule.installation-id",
    epochSequence: "capsule.epoch-sequence",
    epochDigest: "capsule.epoch-digest",
    audience: "capsule.audience",
    purpose: "capsule.purpose",
    status: "capsule.status",
    reason: "capsule.reason",
  };
  const data = (name, min, max, applicationData, nonZeroData) => ({
    key: name,
    valueType: type.data,
    required: true,
    minDataBytes: min,
    maxDataBytes: max,
    applicationData,
    nonZeroData,
  });
  const uint = (name, fixedUInt64) => ({
    key: name,
    valueType: type.uint64,
    required: true,
    fixedUInt64,
  });
  const string = (name, fixedString) => ({
    key: name,
    valueType: type.string,
    required: true,
    fixedString,
  });
  const allowed = (name, allowedUInt64) => ({
    key: name,
    valueType: type.uint64,
    required: true,
    allowedUInt64,
  });
  const methods = [
    [
      "SubmitMainMJSV0",
      1,
      "capsule.daemon.local.v0",
      "capsule.ipc.submit-main-mjs.v0",
      2_097_152,
      16,
      [data("capsule.job-proposal", 1, 2_097_152, true, false)],
      [data("capsule.registration-id", 16, 16, true, true)],
    ],
    [
      "RegisterPlanV0",
      2,
      "capsule.execution-supervisor.local.v0",
      "capsule.ipc.register-plan.v0",
      328_337,
      4_096,
      [
        data("capsule.execution-plan", 1, 65_536, true, false),
        data("capsule.role-bindings", 562, 562, true, false),
        data("capsule.source-manifest", 87, 95, true, false),
        data("capsule.source", 0, 262_144, true, false),
      ],
      [data("capsule.plan-registration", 1, 4_096, true, false)],
    ],
    [
      "GetRegisteredPlanV0",
      3,
      "capsule.execution-supervisor.local.v0",
      "capsule.ipc.get-registered-plan.v0",
      16,
      332_433,
      [data("capsule.registration-id", 16, 16, true, true)],
      [
        data("capsule.execution-plan", 1, 65_536, true, false),
        data("capsule.role-bindings", 562, 562, true, false),
        data("capsule.plan-registration", 1, 4_096, true, false),
        data("capsule.source-manifest", 87, 95, true, false),
        data("capsule.source", 0, 262_144, true, false),
      ],
    ],
    [
      "SubmitApprovalV0",
      4,
      "capsule.execution-supervisor.local.v0",
      "capsule.ipc.submit-approval.v0",
      528,
      16,
      [
        data("capsule.registration-id", 16, 16, true, true),
        data("capsule.approval-envelope", 1, 512, true, false),
      ],
      [
        data("capsule.approval-id", 16, 16, true, true),
        allowed("capsule.approval-state", [1, 2, 3]),
      ],
    ],
    [
      "RequestAttemptV0",
      5,
      "capsule.execution-supervisor.local.v0",
      "capsule.ipc.request-attempt.v0",
      32,
      16,
      [
        data("capsule.registration-id", 16, 16, true, true),
        data("capsule.approval-id", 16, 16, true, true),
      ],
      [data("capsule.attempt-id", 16, 16, true, true), allowed("capsule.attempt-state", [1])],
    ],
  ];
  const result = {};
  for (const [
    method,
    tag,
    audience,
    purpose,
    requestMax,
    replyMax,
    requestBody,
    replyBody,
  ] of methods) {
    const requestHeader = [
      uint(key.protocol, 0),
      uint(key.method, 0),
      uint(key.tag, tag),
      data(key.request, 16, 16, false, true),
      data(key.installation, 16, 16, false, false),
      { key: key.epochSequence, valueType: type.uint64, required: true },
      data(key.epochDigest, 32, 32, false, false),
      string(key.audience, audience),
      string(key.purpose, purpose),
    ];
    const replyHeader = (success) => [
      uint(key.protocol, 0),
      uint(key.method, 0),
      uint(key.tag, tag),
      data(key.request, 16, 16, false, true),
      success
        ? uint(key.status, 0)
        : { key: key.status, valueType: type.uint64, required: true, fixedUInt64: null },
      success
        ? uint(key.reason, 0)
        : { key: key.reason, valueType: type.uint64, required: true, fixedUInt64: null },
    ];
    result[method] = {
      request: {
        method,
        direction: "request",
        protocolVersion: 0,
        methodVersion: 0,
        messageTag: tag,
        exactKeyCount: requestHeader.length + requestBody.length,
        requiredKeyCount: requestHeader.length + requestBody.length,
        optionalKeyCount: 0,
        closedMap: true,
        applicationDataMaxBytes: requestMax,
        fields: [...requestHeader, ...requestBody],
      },
      successReply: {
        method,
        direction: "success-reply",
        protocolVersion: 0,
        methodVersion: 0,
        messageTag: tag,
        exactKeyCount: 6 + replyBody.length,
        requiredKeyCount: 6 + replyBody.length,
        optionalKeyCount: 0,
        closedMap: true,
        applicationDataMaxBytes: replyMax,
        fields: [...replyHeader(true), ...replyBody],
      },
      refusalReply: {
        method,
        direction: "refusal-reply",
        protocolVersion: 0,
        methodVersion: 0,
        messageTag: tag,
        exactKeyCount: 6,
        requiredKeyCount: 6,
        optionalKeyCount: 0,
        closedMap: true,
        applicationDataMaxBytes: 0,
        fields: replyHeader(false),
      },
    };
  }
  return result;
}

function expectedRefusalReplies(contract) {
  return Object.entries({
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
  }).map(([classification, statusTag]) => ({
    classification,
    statusTag,
    reasonTag: contract.reasonTags.coreRefusal,
    bodyKeysAllowed: 0,
    exactKeyCount: 6,
  }));
}

function expectedResponseLoss() {
  return [
    { method: "SubmitMainMJSV0", disposition: "committed-retry-may-create-fresh-registration" },
    { method: "RegisterPlanV0", disposition: "committed-retry-creates-fresh-registration" },
    { method: "GetRegisteredPlanV0", disposition: "repeatable-read-by-registration-id" },
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
  ];
}

function expectedDeadlineCases() {
  const rule =
    "dispatch-only-when-elapsed-milliseconds-is-strictly-less-than-deadline-milliseconds";
  const after = "response-unknown-store-semantic-result-or-recovery-fence-controls";
  return ["SubmitApprovalV0", "RequestAttemptV0"].flatMap((method) => [
    {
      id: `${method}.deadline.immediately-before`,
      method,
      deadlineMilliseconds: 5_000,
      startsAt: "admission",
      clientExtensionAllowed: false,
      boundaryRule: rule,
      afterDispatchDisposition: after,
      boundary: "immediately-before",
      elapsedMilliseconds: 4_999,
      expected: {
        decision: "dispatch-core-before-deadline",
        statusTag: null,
        reasonTag: null,
        replyDisposition: "store-semantic-reply-unless-deadline-cancels-delivery",
        bodyCopied: true,
        coreCalls: 1,
      },
      noState: null,
      postCoreState: "store-semantic-result-or-recovery-fence-controls",
    },
    {
      id: `${method}.deadline.exactly-at`,
      method,
      deadlineMilliseconds: 5_000,
      startsAt: "admission",
      clientExtensionAllowed: false,
      boundaryRule: rule,
      afterDispatchDisposition: after,
      boundary: "exactly-at",
      elapsedMilliseconds: 5_000,
      expected: {
        decision: "deadline-expired-before-dispatch",
        statusTag: null,
        reasonTag: null,
        replyDisposition: "no-application-reply",
        bodyCopied: false,
        coreCalls: 0,
      },
      noState: zeroStateValue(),
      postCoreState: null,
    },
    {
      id: `${method}.deadline.immediately-after`,
      method,
      deadlineMilliseconds: 5_000,
      startsAt: "admission",
      clientExtensionAllowed: false,
      boundaryRule: rule,
      afterDispatchDisposition: after,
      boundary: "immediately-after",
      elapsedMilliseconds: 5_001,
      expected: {
        decision: "deadline-expired-before-dispatch",
        statusTag: null,
        reasonTag: null,
        replyDisposition: "no-application-reply",
        bodyCopied: false,
        coreCalls: 0,
      },
      noState: zeroStateValue(),
      postCoreState: null,
    },
  ]);
}

function verifyCompleteNativeXPCCases(contract) {
  const expectedTable = expectedNativeXPCCaseTable(contract);
  assert.deepEqual(contract.caseTable, expectedTable);
  assert.equal(contract.cases.length, expectedTable.length);
  for (let index = 0; index < expectedTable.length; index += 1) {
    const candidate = contract.cases[index];
    const expectedRow = expectedTable[index];
    assert.deepEqual(caseTableEntry(candidate), expectedRow, candidate.id);
    const topKeys = ["id", candidate.method === undefined ? "methods" : "method"];
    if (expectedRow.direction !== null) topKeys.push("direction");
    if (expectedRow.mutation !== null) topKeys.push("mutation");
    if (
      candidate.id.includes("joint-protocol-method-record-version-mismatch") ||
      candidate.id.includes("joint-method-record-version-mismatch")
    ) {
      topKeys.push("mismatchSet", "selectedFailure", "validationPrecedence");
      assert.deepEqual(candidate.validationPrecedence, contract.validationPrecedence, candidate.id);
      assert.equal(
        candidate.selectedFailure,
        candidate.id.includes("joint-protocol") ? "protocol-version" : "method-version",
        candidate.id,
      );
    }
    if (candidate.id.startsWith("cross-service.")) {
      topKeys.push("dispatchIdentity");
      assert.equal(
        candidate.dispatchIdentity,
        "service-entry-point-and-role-derived-before-tag-cross-check",
        candidate.id,
      );
    }
    if (["all.wrong-service", "all.wrong-role", "all.wrong-session"].includes(candidate.id)) {
      topKeys.push("deliveryPrecondition");
      assert.equal(typeof candidate.deliveryPrecondition, "string", candidate.id);
      const expectedReply =
        candidate.id === "all.wrong-session"
          ? "delivered-in-band-refusal-only-for-this-passive-oracle"
          : "delivered-in-band-refusal";
      assert.equal(candidate.replyDisposition, expectedReply, candidate.id);
    }
    topKeys.push(
      "expected",
      "replyDisposition",
      "terminationDisposition",
      "noState",
      "postCoreState",
    );
    assertExactKeys(candidate, topKeys, candidate.id);
    if (expectedRow.noState === null) assert.equal(candidate.noState, null, candidate.id);
    else assertZeroState(candidate.noState, candidate.id);

    const expectedKeys = [
      "decision",
      "classification",
      "statusTag",
      "reasonTag",
      "bodyCopied",
      "coreCalls",
    ];
    if (
      ["submit-approval.request.exact-maximum", "request-attempt.request.exact-maximum"].includes(
        candidate.id,
      )
    ) {
      expectedKeys.push("applicationDataBytes");
      assert.equal(
        candidate.expected.applicationDataBytes,
        candidate.method === "SubmitApprovalV0" ? 528 : 32,
        candidate.id,
      );
    }
    if (candidate.id === "all.request-id-replaced-by-installation-id-bytes") {
      expectedKeys.push("authorityGrantedByRequestId", "requestIdIsIdempotencyKey");
      assert.equal(candidate.expected.authorityGrantedByRequestId, false);
      assert.equal(candidate.expected.requestIdIsIdempotencyKey, false);
    }
    if (candidate.id === "get.registration-id-replaced-by-request-id-bytes") {
      expectedKeys.push("registrationAuthorityGrantedByWidth");
      assert.equal(candidate.expected.registrationAuthorityGrantedByWidth, false);
    }
    if (
      [
        "submit-approval.registration-id-replaced-by-request-id-bytes",
        "request-attempt.registration-id-replaced-by-approval-id-bytes",
        "request-attempt.approval-id-replaced-by-registration-id-bytes",
      ].includes(candidate.id)
    ) {
      expectedKeys.push("authorityGrantedByWidth");
      assert.equal(candidate.expected.authorityGrantedByWidth, false, candidate.id);
    }
    if (candidate.id.endsWith("cancel-before-dispatch")) {
      expectedKeys.push("authorityStateChanged");
      assert.equal(candidate.expected.authorityStateChanged, false, candidate.id);
    }
    if (candidate.id.endsWith("cancel-after-dispatch")) {
      expectedKeys.push("retryIdentity");
      assert.equal(
        candidate.expected.retryIdentity,
        candidate.method === "SubmitApprovalV0"
          ? "canonical-approval-payload+resolved-signer-authorization-identity"
          : "registration-id+approval-reference",
        candidate.id,
      );
    }
    assertExactKeys(candidate.expected, expectedKeys, `${candidate.id}.expected`);
  }
}

function expectedOracleRefusals() {
  const rows = [
    [
      "submit.method-version",
      "SubmitMainMJSV0",
      "methodVersion",
      "UNSUPPORTED",
      "submit-main-mjs-method-version",
    ],
    ["submit.method", "SubmitMainMJSV0", "method", "UNSUPPORTED", "submit-main-mjs-method"],
    [
      "submit.service",
      "SubmitMainMJSV0",
      "service",
      "AUTHENTICATION",
      "submit-main-mjs-method-binding",
    ],
    [
      "submit.role",
      "SubmitMainMJSV0",
      "expectedRole",
      "AUTHENTICATION",
      "submit-main-mjs-method-binding",
    ],
    [
      "submit.signing-identifier",
      "SubmitMainMJSV0",
      "expectedSigningIdentifier",
      "AUTHENTICATION",
      "submit-main-mjs-method-binding",
    ],
    [
      "submit.audience",
      "SubmitMainMJSV0",
      "audience",
      "AUTHENTICATION",
      "submit-main-mjs-method-binding",
    ],
    [
      "submit.purpose",
      "SubmitMainMJSV0",
      "purpose",
      "AUTHENTICATION",
      "submit-main-mjs-method-binding",
    ],
    [
      "submit.protocol-version",
      "SubmitMainMJSV0",
      "protocolVersion",
      "UNSUPPORTED",
      "ipc-protocol-version",
    ],
    ["submit.zero-request-id", "SubmitMainMJSV0", "requestId", "SCHEMA", "ipc-request-id"],
    [
      "submit.installation",
      "SubmitMainMJSV0",
      "installationId",
      "BINDING",
      "ipc-current-supervisor-state",
    ],
    [
      "submit.epoch-sequence",
      "SubmitMainMJSV0",
      "epochSequence",
      "BINDING",
      "ipc-current-supervisor-state",
    ],
    [
      "submit.epoch-digest",
      "SubmitMainMJSV0",
      "epochDigest",
      "BINDING",
      "ipc-current-supervisor-state",
    ],
    [
      "register.method-version",
      "RegisterPlanV0",
      "methodVersion",
      "UNSUPPORTED",
      "register-plan-method-record-version",
    ],
    [
      "register.role-binding-record-version",
      "RegisterPlanV0",
      "roleBindingRecordVersion",
      "UNSUPPORTED",
      "register-plan-method-record-version",
    ],
    ["register.method", "RegisterPlanV0", "method", "UNSUPPORTED", "register-plan-method"],
    [
      "register.service",
      "RegisterPlanV0",
      "service",
      "AUTHENTICATION",
      "register-plan-method-binding",
    ],
    [
      "register.role",
      "RegisterPlanV0",
      "expectedRole",
      "AUTHENTICATION",
      "register-plan-method-binding",
    ],
    [
      "register.audience",
      "RegisterPlanV0",
      "audience",
      "AUTHENTICATION",
      "register-plan-method-binding",
    ],
    [
      "register.purpose",
      "RegisterPlanV0",
      "purpose",
      "AUTHENTICATION",
      "register-plan-method-binding",
    ],
    [
      "register.protocol-version",
      "RegisterPlanV0",
      "protocolVersion",
      "UNSUPPORTED",
      "ipc-protocol-version",
    ],
    ["register.zero-request-id", "RegisterPlanV0", "requestId", "SCHEMA", "ipc-request-id"],
    [
      "register.installation",
      "RegisterPlanV0",
      "installationId",
      "BINDING",
      "ipc-current-supervisor-state",
    ],
    [
      "register.epoch-sequence",
      "RegisterPlanV0",
      "epochSequence",
      "BINDING",
      "ipc-current-supervisor-state",
    ],
    [
      "register.epoch-digest",
      "RegisterPlanV0",
      "epochDigest",
      "BINDING",
      "ipc-current-supervisor-state",
    ],
    [
      "fetch.method-version",
      "GetRegisteredPlanV0",
      "methodVersion",
      "UNSUPPORTED",
      "get-registered-plan-method-record-version",
    ],
    [
      "fetch.role-binding-record-version",
      "GetRegisteredPlanV0",
      "roleBindingRecordVersion",
      "UNSUPPORTED",
      "get-registered-plan-method-record-version",
    ],
    ["fetch.method", "GetRegisteredPlanV0", "method", "UNSUPPORTED", "get-registered-plan-method"],
    [
      "fetch.service",
      "GetRegisteredPlanV0",
      "service",
      "AUTHENTICATION",
      "get-registered-plan-method-binding",
    ],
    [
      "fetch.role",
      "GetRegisteredPlanV0",
      "expectedRole",
      "AUTHENTICATION",
      "get-registered-plan-method-binding",
    ],
    [
      "fetch.audience",
      "GetRegisteredPlanV0",
      "audience",
      "AUTHENTICATION",
      "get-registered-plan-method-binding",
    ],
    [
      "fetch.purpose",
      "GetRegisteredPlanV0",
      "purpose",
      "AUTHENTICATION",
      "get-registered-plan-method-binding",
    ],
    [
      "fetch.protocol-version",
      "GetRegisteredPlanV0",
      "protocolVersion",
      "UNSUPPORTED",
      "ipc-protocol-version",
    ],
    ["fetch.zero-request-id", "GetRegisteredPlanV0", "requestId", "SCHEMA", "ipc-request-id"],
    [
      "fetch.installation",
      "GetRegisteredPlanV0",
      "installationId",
      "BINDING",
      "ipc-current-supervisor-state",
    ],
    [
      "fetch.epoch-sequence",
      "GetRegisteredPlanV0",
      "epochSequence",
      "BINDING",
      "ipc-current-supervisor-state",
    ],
    [
      "fetch.epoch-digest",
      "GetRegisteredPlanV0",
      "epochDigest",
      "BINDING",
      "ipc-current-supervisor-state",
    ],
    [
      "fetch.zero-registration-id",
      "GetRegisteredPlanV0",
      "registrationId",
      "SCHEMA",
      "registration-id",
    ],
  ];
  return rows.map(([id, method, mutation, classification, reason]) => ({
    id,
    method,
    mutation,
    expected: { decision: "reject", classification, reason },
    noState: zeroStateValue(),
  }));
}

function expectedOracleResponseLoss() {
  return [
    {
      method: "SubmitMainMJSV0",
      disposition: "committed-retry-may-create-fresh-registration",
      downstreamDisposition: "committed-retry-creates-fresh-registration",
      retryMayCreateFreshRegistration: true,
      requestIdIsIdempotencyKey: false,
    },
    {
      method: "RegisterPlanV0",
      disposition: "committed-retry-creates-fresh-registration",
      firstCommittedRegistrations: 1,
      retryCommittedRegistrations: 2,
      bothRegistrationsSeparatelyReadable: true,
      requestIdIsIdempotencyKey: false,
    },
    {
      method: "GetRegisteredPlanV0",
      disposition: "repeatable-read-by-registration-id",
      storeMutation: false,
      repeatedReplyBodyByteEqual: true,
      requestIdIsIdempotencyKey: false,
    },
    ...expectedResponseLoss().slice(3),
  ];
}

function expectedCancellationAndDeadline() {
  return [
    {
      id: "cancel.before-dispatch",
      disposition: "caller-cancelled-before-dispatch",
      coreCalls: 0,
      noState: zeroStateValue(),
    },
    {
      id: "cancel.after-dispatch",
      disposition: "caller-cancelled-after-dispatch-response-unknown",
      coreCalls: 1,
      state: "method-semantic-result-controls;transport-does-not-infer-abort",
      admittedSlotHeldUntilCoreReturns: true,
    },
    {
      id: "deadline.after-dispatch",
      disposition: "method-deadline-after-dispatch-response-unknown",
      coreCalls: 1,
      state: "method-semantic-result-or-recovery-fence-controls;transport-does-not-infer-abort",
      admittedSlotHeldUntilCoreReturns: true,
    },
    {
      id: "submit-approval.deadline-after-dispatch-commit",
      method: "SubmitApprovalV0",
      disposition: "deadline-cancels-response-delivery-only",
      coreCalls: 1,
      replyDisposition: "no-application-reply-on-expired-call",
      statusTag: null,
      reasonTag: null,
      postCoreState: {
        durableAuthority: "approval-commit-remains-authoritative",
        replayIdentity: "canonical-approval-payload+resolved-signer-authorization-identity",
        replayResult: "same-approval-id-and-current-state",
        newRefusalAllowed: false,
      },
      noState: null,
    },
    {
      id: "request-attempt.deadline-after-dispatch-commit",
      method: "RequestAttemptV0",
      disposition: "deadline-cancels-response-delivery-only",
      coreCalls: 1,
      replyDisposition: "no-application-reply-on-expired-call",
      statusTag: null,
      reasonTag: null,
      postCoreState: {
        durableAuthority: "attempt-commit-remains-authoritative",
        replayIdentity: "registration-id+approval-reference",
        replayResult: "same-attempt-id-and-current-state",
        newRefusalAllowed: false,
      },
      noState: null,
    },
    {
      id: "downstream.stall",
      disposition: "method-deadline-after-dispatch-response-unknown",
      newWorkOnSameConnection: "CAPACITY",
      queueDepth: 0,
      processTerminationEvidence: false,
    },
  ];
}

function expectedFlowControl() {
  const zero = zeroStateValue;
  return [
    {
      id: "submit.connection.concurrent-cap-plus-one",
      method: "SubmitMainMJSV0",
      admitted: 1,
      attempted: 2,
      sameConnection: true,
      classification: "CAPACITY",
      reason: "ipc-connection-in-flight",
      queueDepth: 0,
      noState: zero(),
    },
    {
      id: "submit.service-connections-cap-plus-one",
      method: "SubmitMainMJSV0",
      admitted: 4,
      attempted: 5,
      classification: "CAPACITY",
      reason: "ipc-service-connections",
      queueDepth: 0,
      noState: zero(),
    },
    {
      id: "supervisor.combined-process-cap-plus-one",
      methods: [{ method: "RegisterPlanV0" }, { method: "GetRegisteredPlanV0" }],
      admitted: 8,
      attempted: 9,
      classification: "CAPACITY",
      reason: "ipc-service-connections",
      alsoAtProcessCap: true,
      queueDepth: 0,
      noState: zero(),
    },
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
      noState: zero(),
    },
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
      noState: zero(),
    },
  ];
}

function runHardenedMutationProofs(contract, oracleSet) {
  const mutations = [
    ["missing dictionary field", (c) => c.envelopes.SubmitApprovalV0.request.fields.pop()],
    [
      "extra dictionary field",
      (c) =>
        c.envelopes.RequestAttemptV0.request.fields.push({
          key: "capsule.extra",
          valueType: "XPC_TYPE_DATA",
        }),
    ],
    [
      "changed type or width",
      (c) => {
        c.envelopes.SubmitApprovalV0.request.fields.at(-1).maxDataBytes = 511;
      },
    ],
    [
      "absent required noState",
      (c) => {
        delete c.cases.find((entry) => entry.id === "all.missing-key").noState;
      },
    ],
    [
      "modified cancellation result",
      (_c, o) => {
        o.cancellationAndDeadline[1].coreCalls = 0;
      },
    ],
    [
      "missing refusal case",
      (c) => {
        c.cases = c.cases.filter((entry) => entry.id !== "all.wrong-role");
      },
    ],
    [
      "response-loss drift",
      (c) => {
        c.responseLoss[3].sameApprovalID = false;
      },
    ],
    [
      "exact-at-boundary inversion",
      (c) => {
        c.deadlineCases.find(
          (entry) => entry.id === "SubmitApprovalV0.deadline.exactly-at",
        ).expected.decision = "dispatch-core-before-deadline";
      },
    ],
  ];
  for (const [label, mutate] of mutations) {
    const changedContract = JSON.parse(JSON.stringify(contract));
    const changedOracles = JSON.parse(JSON.stringify(oracleSet));
    mutate(changedContract, changedOracles);
    assert.throws(() => verifyHardenedNativeXPC(changedContract, changedOracles), undefined, label);
  }
}

function assertExactKeys(value, expected, label) {
  assert.deepEqual(Object.keys(value), expected, `${label} exact keys/order`);
}

function zeroStateValue() {
  return {
    authorityStateChanged: false,
    coreCalls: 0,
    registrationsCreated: 0,
    approvalsConsumed: 0,
    attemptsCreated: 0,
    lifecycleCalls: 0,
    backendCalls: 0,
  };
}
