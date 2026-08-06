import { mkdir, readFile, writeFile } from "node:fs/promises";
import { sha256Hex } from "./lib/fixture-bytes.mjs";

const repository = new URL("../", import.meta.url);
const outputRoot = new URL("schemas/conformance/authenticated-local-ipc-v0/", repository);

const protocolVersion = 0;
const audience = "capsule.execution-supervisor.local.v0";
const submissionAudience = "capsule.daemon.local.v0";
const signedApprovalAudience = "capsule.execution-supervisor";
const signedApprovalPurpose = "capsule.plan.approve";
const roleBindingRecordVersion = 0;
const authorityDisposition = "service-role-purpose-audience-derived";
const requestIdDisposition = "correlation-only";
const nativeXPCEncoding = "xpc-dictionary-v0";
const nativeXPCEncodingVersion = 0;
const nativeXPCValidationPrecedence = Object.freeze([
  "protocol-version",
  "method-version",
  "service-entry-point-role-message-tag-audience-purpose",
  "installation-epoch-current-state",
  "application-data-copy",
  "embedded-record-version-and-core-validation",
]);
const nativeXPCTypes = Object.freeze({
  uint64: "XPC_TYPE_UINT64",
  data: "XPC_TYPE_DATA",
  string: "XPC_TYPE_STRING",
});
const nativeXPCKeys = Object.freeze({
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
});
const nativeXPCMessageTags = Object.freeze({
  invalid: 0,
  SubmitMainMJSV0: 1,
  RegisterPlanV0: 2,
  GetRegisteredPlanV0: 3,
});
const nativeXPCStatusTags = Object.freeze({
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
const nativeXPCReasonTags = Object.freeze({
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
const supervisorFlow = Object.freeze({
  peerRequirementBeforeDeliveryRequired: true,
  methodAuthorityDisposition: authorityDisposition,
  requestIdAuthorityDisposition: requestIdDisposition,
  serviceMaxConnections: 4,
  connectionMaxInFlight: 1,
  processMaxAdmittedRequests: 8,
  applicationQueueCapacity: 0,
  inFlightRequestDataMaxBytes: 2_626_696,
});
const daemonFlow = Object.freeze({
  peerRequirementBeforeDeliveryRequired: true,
  methodAuthorityDisposition: authorityDisposition,
  requestIdAuthorityDisposition: requestIdDisposition,
  serviceMaxConnections: 4,
  connectionMaxInFlight: 1,
  processMaxAdmittedRequests: 4,
  applicationQueueCapacity: 0,
  inFlightRequestDataMaxBytes: 8_388_608,
});
const caps = Object.freeze({
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
if (
  caps.registerPlanV0Request !==
    caps.executionPlan + caps.roleBindings + caps.sourceManifest + caps.source ||
  caps.getRegisteredPlanV0Reply !== caps.registerPlanV0Request + caps.planRegistration
) {
  throw new Error("authenticated-local-IPC aggregate cap drift");
}

const methods = {
  submitMainMJSV0: {
    objectType: "capsule.authenticated-local-ipc-submit-main-mjs-v0-method-record",
    objectVersion: 0,
    method: "SubmitMainMJSV0",
    methodVersion: 0,
    service: "com.capsulecorp.capsule.daemon.cli.v0",
    expectedRole: "internal-alpha-cli",
    expectedSigningIdentifier: "com.capsulecorp.capsule.cli",
    audience: submissionAudience,
    purpose: "capsule.ipc.submit-main-mjs.v0",
    requestMediaType: "application/capsule.job-proposal+json;v=0",
    requestDataMaxBytes: caps.submitMainMJSV0Request,
    replyDataMaxBytes: caps.submitMainMJSV0Reply,
    deadlineMilliseconds: 10_000,
    responseLossDisposition: "committed-retry-may-create-fresh-registration",
    ...daemonFlow,
  },
  registerPlanV0: {
    objectType: "capsule.authenticated-local-ipc-register-plan-v0-method-record",
    objectVersion: 0,
    method: "RegisterPlanV0",
    methodVersion: 0,
    roleBindingRecordVersion,
    service: "com.capsulecorp.capsule.supervisor.daemon.v0",
    expectedRole: "daemon",
    audience,
    purpose: "capsule.ipc.register-plan.v0",
    requestDataMaxBytes: caps.registerPlanV0Request,
    replyDataMaxBytes: caps.registerPlanV0Reply,
    deadlineMilliseconds: 5_000,
    responseLossDisposition: "committed-retry-creates-fresh-registration",
    ...supervisorFlow,
  },
  getRegisteredPlanV0: {
    objectType: "capsule.authenticated-local-ipc-get-registered-plan-v0-method-record",
    objectVersion: 0,
    method: "GetRegisteredPlanV0",
    methodVersion: 0,
    roleBindingRecordVersion,
    service: "com.capsulecorp.capsule.supervisor.broker.v0",
    expectedRole: "broker",
    audience,
    purpose: "capsule.ipc.get-registered-plan.v0",
    requestDataMaxBytes: caps.getRegisteredPlanV0Request,
    replyDataMaxBytes: caps.getRegisteredPlanV0Reply,
    deadlineMilliseconds: 2_000,
    responseLossDisposition: "repeatable-read-by-registration-id",
    ...supervisorFlow,
  },
};

const nativeXPCEnvelopes = Object.fromEntries(
  Object.values(methods).map((method) => {
    const tag = nativeXPCMessageTags[method.method];
    const body = nativeXPCBodyFields(method.method);
    return [
      method.method,
      {
        request: nativeXPCEnvelope(
          method,
          "request",
          tag,
          9 + body.request.length,
          method.requestDataMaxBytes,
          [...nativeXPCRequestHeader(method, tag), ...body.request],
        ),
        successReply: nativeXPCEnvelope(
          method,
          "success-reply",
          tag,
          6 + body.reply.length,
          method.replyDataMaxBytes,
          [...nativeXPCReplyHeader(method, tag, true), ...body.reply],
        ),
        refusalReply: nativeXPCEnvelope(
          method,
          "refusal-reply",
          tag,
          6,
          0,
          nativeXPCReplyHeader(method, tag, false),
        ),
      },
    ];
  }),
);

const proposal = await repositoryFixture("schemas/conformance/v0/job-proposal/ordinary.json");
const plan = await repositoryFixture("schemas/conformance/authority-plane-v0/execution-plan.cbor");
const bindings = await repositoryFixture(
  "schemas/conformance/authority-plane-v0/role-bindings.bin",
);
const sourceManifest = await repositoryFixture(
  "schemas/conformance/authority-plane-v0/source-manifest.cbor",
);
const source = await repositoryFixture("schemas/conformance/authority-plane-v0/main.mjs");
const registration = await repositoryFixture(
  "schemas/conformance/v0/plan-registration/ordinary.cbor",
);

if (bindings.byteLength !== caps.roleBindings || bindings.bytes[0] !== roleBindingRecordVersion) {
  throw new Error("role-binding record/version drift");
}
const decodedProposal = JSON.parse(proposal.bytes);
const proposalSource = Buffer.from(decodedProposal.source.files["main.mjs"], "utf8");
const proposalSourceSha256 = sha256Hex(proposalSource);
if (!proposalSource.equals(source.bytes)) {
  throw new Error(
    "submission proposal source differs from authority-plane main.mjs custody fixture",
  );
}

const submitRequest = {
  objectType: "capsule.authenticated-local-ipc-submit-main-mjs-v0-request",
  objectVersion: 0,
  fixtureSerialization: "exact-json-not-xpc-framing",
  protocolVersion,
  requestId: repeatedHex(0x31, 16),
  installationId: repeatedHex(0x11, 16),
  epochSequence: 7,
  epochDigest: repeatedHex(0x22, 32),
  methodRecord: reference("submit-main-mjs-v0.method.json", jsonBytes(methods.submitMainMJSV0)),
  body: { exactJobProposalBytes: proposal.reference },
  applicationDataBytes: proposal.byteLength,
};

const submitReply = {
  objectType: "capsule.authenticated-local-ipc-submit-main-mjs-v0-reply",
  objectVersion: 0,
  fixtureSerialization: "exact-json-not-xpc-framing",
  requestId: submitRequest.requestId,
  body: { registrationId: repeatedHex(0x77, 16) },
  applicationDataBytes: 16,
};

const registerRequest = {
  objectType: "capsule.authenticated-local-ipc-register-plan-v0-request",
  objectVersion: 0,
  fixtureSerialization: "exact-json-not-xpc-framing",
  protocolVersion,
  requestId: repeatedHex(0x41, 16),
  installationId: repeatedHex(0x11, 16),
  epochSequence: 7,
  epochDigest: repeatedHex(0x22, 32),
  methodRecord: reference("register-plan-v0.method.json", jsonBytes(methods.registerPlanV0)),
  body: {
    exactPlanBytes: plan.reference,
    roleBindingBytes: bindings.reference,
    sourceManifestBytes: sourceManifest.reference,
    sourceBytes: source.reference,
  },
  applicationDataBytes:
    plan.byteLength + bindings.byteLength + sourceManifest.byteLength + source.byteLength,
};

const registerReply = {
  objectType: "capsule.authenticated-local-ipc-register-plan-v0-reply",
  objectVersion: 0,
  fixtureSerialization: "exact-json-not-xpc-framing",
  requestId: registerRequest.requestId,
  body: { planRegistrationBytes: registration.reference },
  applicationDataBytes: registration.byteLength,
};

const getRequest = {
  objectType: "capsule.authenticated-local-ipc-get-registered-plan-v0-request",
  objectVersion: 0,
  fixtureSerialization: "exact-json-not-xpc-framing",
  protocolVersion,
  requestId: repeatedHex(0x51, 16),
  installationId: repeatedHex(0x11, 16),
  epochSequence: 7,
  epochDigest: repeatedHex(0x22, 32),
  methodRecord: reference(
    "get-registered-plan-v0.method.json",
    jsonBytes(methods.getRegisteredPlanV0),
  ),
  body: { registrationId: repeatedHex(0x77, 16) },
  applicationDataBytes: 16,
};

const getReply = {
  objectType: "capsule.authenticated-local-ipc-get-registered-plan-v0-reply",
  objectVersion: 0,
  fixtureSerialization: "exact-json-not-xpc-framing",
  requestId: getRequest.requestId,
  body: {
    exactPlanBytes: plan.reference,
    roleBindingBytes: bindings.reference,
    planRegistrationBytes: registration.reference,
    sourceManifestBytes: sourceManifest.reference,
    sourceBytes: source.reference,
  },
  applicationDataBytes:
    plan.byteLength +
    bindings.byteLength +
    registration.byteLength +
    sourceManifest.byteLength +
    source.byteLength,
};

const zeroState = Object.freeze({
  authorityStateChanged: false,
  coreCalls: 0,
  registrationsCreated: 0,
  approvalsConsumed: 0,
  attemptsCreated: 0,
  lifecycleCalls: 0,
  backendCalls: 0,
});
const zeroEffects = Object.freeze({
  ipcEndpoint: false,
  peerAuthenticated: false,
  keyUsed: false,
  processCreated: false,
  runtimeCreated: false,
  backendCreated: false,
  guestCreated: false,
});

const nativeXPCContract = {
  objectType: "capsule.authenticated-local-ipc-native-xpc-contract",
  objectVersion: 0,
  status: "passive-unwired-no-listener",
  transportEncoding: nativeXPCEncoding,
  transportEncodingVersion: nativeXPCEncodingVersion,
  fixtureSerialization: "exact-json-description-of-xpc-dictionaries",
  topLevelObjectType: "XPC_TYPE_DICTIONARY",
  valueTypes: nativeXPCTypes,
  keys: nativeXPCKeys,
  messageTags: nativeXPCMessageTags,
  statusTags: nativeXPCStatusTags,
  reasonTags: nativeXPCReasonTags,
  classificationToStatus: Object.fromEntries(
    Object.entries(nativeXPCStatusTags).filter(([classification]) => classification !== "OK"),
  ),
  structuralReasonToStatus: {
    keySet: nativeXPCStatusTags.MALFORMED,
    valueType: nativeXPCStatusTags.MALFORMED,
    dataWidth: nativeXPCStatusTags.SCHEMA,
    dataCap: nativeXPCStatusTags.MALFORMED,
    zeroIdentifier: nativeXPCStatusTags.SCHEMA,
    epochSequence: nativeXPCStatusTags.SCHEMA,
    protocolVersion: nativeXPCStatusTags.UNSUPPORTED,
    methodVersion: nativeXPCStatusTags.UNSUPPORTED,
    messageTag: nativeXPCStatusTags.UNSUPPORTED,
    methodBinding: nativeXPCStatusTags.AUTHENTICATION,
    currentState: nativeXPCStatusTags.BINDING,
    capacity: nativeXPCStatusTags.CAPACITY,
  },
  coreRefusalMapping: "classification-selects-status;reason=coreRefusal",
  localIntegrityMapping: "terminate-without-reply;reason-tag-is-fixture-diagnostic-only",
  methodBindings: Object.fromEntries(
    Object.values(methods).map((method) => [
      method.method,
      {
        service: method.service,
        expectedRole: method.expectedRole,
        expectedSigningIdentifier: method.expectedSigningIdentifier ?? null,
        audience: method.audience,
        purpose: method.purpose,
        messageTag: nativeXPCMessageTags[method.method],
        methodVersion: method.methodVersion,
      },
    ]),
  ),
  requestCommonKeyCount: 9,
  replyCommonKeyCount: 6,
  extraObjectsAllowed: 0,
  fileDescriptorsAllowed: 0,
  endpointsAllowed: 0,
  machRightsAllowed: 0,
  nestedContainersAllowed: 0,
  messageTagDisposition: "method-specific-cross-check-not-dispatch-opcode",
  validationPrecedence: nativeXPCValidationPrecedence,
  requestIdDisposition,
  copyDisposition: "body-copy-only-after-peer-flow-shape-current-state-binding",
  localIntegrityDisposition:
    "oversize-output-short-write-pointer-length-or-bridge-version-fault-terminates-process-without-reply",
  envelopes: nativeXPCEnvelopes,
  cases: nativeXPCCases(),
  refusalReplies: Object.entries(nativeXPCStatusTags)
    .filter(([classification]) => classification !== "OK")
    .map(([classification, statusTag]) => ({
      classification,
      statusTag,
      reasonTag: nativeXPCReasonTags.coreRefusal,
      bodyKeysAllowed: 0,
      exactKeyCount: 6,
    })),
  responseLoss: [
    { method: "SubmitMainMJSV0", disposition: methods.submitMainMJSV0.responseLossDisposition },
    { method: "RegisterPlanV0", disposition: methods.registerPlanV0.responseLossDisposition },
    {
      method: "GetRegisteredPlanV0",
      disposition: methods.getRegisteredPlanV0.responseLossDisposition,
    },
  ],
  futureNativeHarnessOracles: [
    {
      id: "os-peer-requirement-mismatch",
      scope: "future-external-native-harness-only",
      expected: "no-message-delivery-and-no-application-reply",
      currentEvidence: false,
      inBandRefusal: false,
    },
  ],
  peerAuthenticationEvidence: null,
  listenerActivated: false,
  serviceRegistered: false,
};

const oracles = {
  objectType: "capsule.authenticated-local-ipc-passive-oracles",
  objectVersion: 0,
  maxima: [
    maximumCase(
      "submit-main-mjs-v0.request.exact-maximum",
      "SubmitMainMJSV0",
      "request",
      { exactJobProposalBytes: caps.jobProposal },
      caps.submitMainMJSV0Request,
      "accept-outer-shape",
    ),
    maximumCase(
      "submit-main-mjs-v0.request.cap-plus-one",
      "SubmitMainMJSV0",
      "request",
      { exactJobProposalBytes: caps.jobProposal + 1 },
      caps.submitMainMJSV0Request + 1,
      "reject",
      "MALFORMED",
      "submit-main-mjs-proposal-bytes",
    ),
    maximumCase(
      "submit-main-mjs-v0.reply.exact",
      "SubmitMainMJSV0",
      "reply",
      { registrationId: 16 },
      caps.submitMainMJSV0Reply,
      "accept-outer-shape",
    ),
    maximumCase(
      "submit-main-mjs-v0.reply.cap-plus-one",
      "SubmitMainMJSV0",
      "reply",
      { registrationId: 17 },
      caps.submitMainMJSV0Reply + 1,
      "local-integrity-fault",
      "LOCAL_FAILURE",
      "submit-main-mjs-reply-shape",
    ),
    maximumCase(
      "register-plan-v0.request.exact-maximum",
      "RegisterPlanV0",
      "request",
      {
        exactPlanBytes: caps.executionPlan,
        roleBindingBytes: caps.roleBindings,
        sourceManifestBytes: caps.sourceManifest,
        sourceBytes: caps.source,
      },
      caps.registerPlanV0Request,
      "accept-outer-shape",
    ),
    maximumCase(
      "register-plan-v0.request.source-cap-plus-one",
      "RegisterPlanV0",
      "request",
      {
        exactPlanBytes: caps.executionPlan,
        roleBindingBytes: caps.roleBindings,
        sourceManifestBytes: caps.sourceManifest,
        sourceBytes: caps.source + 1,
      },
      caps.registerPlanV0Request + 1,
      "reject",
      "MALFORMED",
      "register-source-bytes",
    ),
    maximumCase(
      "register-plan-v0.reply.exact-maximum",
      "RegisterPlanV0",
      "reply",
      {
        planRegistrationBytes: caps.planRegistration,
      },
      caps.registerPlanV0Reply,
      "accept-outer-shape",
    ),
    maximumCase(
      "register-plan-v0.reply.cap-plus-one",
      "RegisterPlanV0",
      "reply",
      {
        planRegistrationBytes: caps.planRegistration + 1,
      },
      caps.registerPlanV0Reply + 1,
      "local-integrity-fault",
      "LOCAL_FAILURE",
      "register-plan-reply-cap",
    ),
    maximumCase(
      "get-registered-plan-v0.request.exact-maximum",
      "GetRegisteredPlanV0",
      "request",
      {
        registrationId: 16,
      },
      caps.getRegisteredPlanV0Request,
      "accept-outer-shape",
    ),
    maximumCase(
      "get-registered-plan-v0.request.cap-plus-one",
      "GetRegisteredPlanV0",
      "request",
      {
        registrationId: 17,
      },
      caps.getRegisteredPlanV0Request + 1,
      "reject",
      "SCHEMA",
      "registration-id",
    ),
    maximumCase(
      "get-registered-plan-v0.reply.exact-maximum",
      "GetRegisteredPlanV0",
      "reply",
      {
        exactPlanBytes: caps.executionPlan,
        roleBindingBytes: caps.roleBindings,
        planRegistrationBytes: caps.planRegistration,
        sourceManifestBytes: caps.sourceManifest,
        sourceBytes: caps.source,
      },
      caps.getRegisteredPlanV0Reply,
      "accept-outer-shape",
    ),
    maximumCase(
      "get-registered-plan-v0.reply.source-cap-plus-one",
      "GetRegisteredPlanV0",
      "reply",
      {
        exactPlanBytes: caps.executionPlan,
        roleBindingBytes: caps.roleBindings,
        planRegistrationBytes: caps.planRegistration,
        sourceManifestBytes: caps.sourceManifest,
        sourceBytes: caps.source + 1,
      },
      caps.getRegisteredPlanV0Reply + 1,
      "local-integrity-fault",
      "LOCAL_FAILURE",
      "get-registered-plan-reply-shape",
    ),
  ],
  refusals: [
    refusal(
      "submit.method-version",
      "SubmitMainMJSV0",
      "methodVersion",
      "UNSUPPORTED",
      "submit-main-mjs-method-version",
    ),
    refusal("submit.method", "SubmitMainMJSV0", "method", "UNSUPPORTED", "submit-main-mjs-method"),
    refusal(
      "submit.service",
      "SubmitMainMJSV0",
      "service",
      "AUTHENTICATION",
      "submit-main-mjs-method-binding",
    ),
    refusal(
      "submit.role",
      "SubmitMainMJSV0",
      "expectedRole",
      "AUTHENTICATION",
      "submit-main-mjs-method-binding",
    ),
    refusal(
      "submit.signing-identifier",
      "SubmitMainMJSV0",
      "expectedSigningIdentifier",
      "AUTHENTICATION",
      "submit-main-mjs-method-binding",
    ),
    refusal(
      "submit.audience",
      "SubmitMainMJSV0",
      "audience",
      "AUTHENTICATION",
      "submit-main-mjs-method-binding",
    ),
    refusal(
      "submit.purpose",
      "SubmitMainMJSV0",
      "purpose",
      "AUTHENTICATION",
      "submit-main-mjs-method-binding",
    ),
    refusal(
      "submit.protocol-version",
      "SubmitMainMJSV0",
      "protocolVersion",
      "UNSUPPORTED",
      "ipc-protocol-version",
    ),
    refusal("submit.zero-request-id", "SubmitMainMJSV0", "requestId", "SCHEMA", "ipc-request-id"),
    refusal(
      "submit.installation",
      "SubmitMainMJSV0",
      "installationId",
      "BINDING",
      "ipc-current-supervisor-state",
    ),
    refusal(
      "submit.epoch-sequence",
      "SubmitMainMJSV0",
      "epochSequence",
      "BINDING",
      "ipc-current-supervisor-state",
    ),
    refusal(
      "submit.epoch-digest",
      "SubmitMainMJSV0",
      "epochDigest",
      "BINDING",
      "ipc-current-supervisor-state",
    ),
    refusal(
      "register.method-version",
      "RegisterPlanV0",
      "methodVersion",
      "UNSUPPORTED",
      "register-plan-method-record-version",
    ),
    refusal(
      "register.role-binding-record-version",
      "RegisterPlanV0",
      "roleBindingRecordVersion",
      "UNSUPPORTED",
      "register-plan-method-record-version",
    ),
    refusal("register.method", "RegisterPlanV0", "method", "UNSUPPORTED", "register-plan-method"),
    refusal(
      "register.service",
      "RegisterPlanV0",
      "service",
      "AUTHENTICATION",
      "register-plan-method-binding",
    ),
    refusal(
      "register.role",
      "RegisterPlanV0",
      "expectedRole",
      "AUTHENTICATION",
      "register-plan-method-binding",
    ),
    refusal(
      "register.audience",
      "RegisterPlanV0",
      "audience",
      "AUTHENTICATION",
      "register-plan-method-binding",
    ),
    refusal(
      "register.purpose",
      "RegisterPlanV0",
      "purpose",
      "AUTHENTICATION",
      "register-plan-method-binding",
    ),
    refusal(
      "register.protocol-version",
      "RegisterPlanV0",
      "protocolVersion",
      "UNSUPPORTED",
      "ipc-protocol-version",
    ),
    refusal("register.zero-request-id", "RegisterPlanV0", "requestId", "SCHEMA", "ipc-request-id"),
    refusal(
      "register.installation",
      "RegisterPlanV0",
      "installationId",
      "BINDING",
      "ipc-current-supervisor-state",
    ),
    refusal(
      "register.epoch-sequence",
      "RegisterPlanV0",
      "epochSequence",
      "BINDING",
      "ipc-current-supervisor-state",
    ),
    refusal(
      "register.epoch-digest",
      "RegisterPlanV0",
      "epochDigest",
      "BINDING",
      "ipc-current-supervisor-state",
    ),
    refusal(
      "fetch.method-version",
      "GetRegisteredPlanV0",
      "methodVersion",
      "UNSUPPORTED",
      "get-registered-plan-method-record-version",
    ),
    refusal(
      "fetch.role-binding-record-version",
      "GetRegisteredPlanV0",
      "roleBindingRecordVersion",
      "UNSUPPORTED",
      "get-registered-plan-method-record-version",
    ),
    refusal(
      "fetch.method",
      "GetRegisteredPlanV0",
      "method",
      "UNSUPPORTED",
      "get-registered-plan-method",
    ),
    refusal(
      "fetch.service",
      "GetRegisteredPlanV0",
      "service",
      "AUTHENTICATION",
      "get-registered-plan-method-binding",
    ),
    refusal(
      "fetch.role",
      "GetRegisteredPlanV0",
      "expectedRole",
      "AUTHENTICATION",
      "get-registered-plan-method-binding",
    ),
    refusal(
      "fetch.audience",
      "GetRegisteredPlanV0",
      "audience",
      "AUTHENTICATION",
      "get-registered-plan-method-binding",
    ),
    refusal(
      "fetch.purpose",
      "GetRegisteredPlanV0",
      "purpose",
      "AUTHENTICATION",
      "get-registered-plan-method-binding",
    ),
    refusal(
      "fetch.protocol-version",
      "GetRegisteredPlanV0",
      "protocolVersion",
      "UNSUPPORTED",
      "ipc-protocol-version",
    ),
    refusal(
      "fetch.zero-request-id",
      "GetRegisteredPlanV0",
      "requestId",
      "SCHEMA",
      "ipc-request-id",
    ),
    refusal(
      "fetch.installation",
      "GetRegisteredPlanV0",
      "installationId",
      "BINDING",
      "ipc-current-supervisor-state",
    ),
    refusal(
      "fetch.epoch-sequence",
      "GetRegisteredPlanV0",
      "epochSequence",
      "BINDING",
      "ipc-current-supervisor-state",
    ),
    refusal(
      "fetch.epoch-digest",
      "GetRegisteredPlanV0",
      "epochDigest",
      "BINDING",
      "ipc-current-supervisor-state",
    ),
    refusal(
      "fetch.zero-registration-id",
      "GetRegisteredPlanV0",
      "registrationId",
      "SCHEMA",
      "registration-id",
    ),
  ],
  copyOwnership: [
    "submission-proposal-caller-and-accessor-mutation-do-not-change-copied-bytes",
    "caller-mutation-after-request-construction-does-not-change-copied-body",
    "request-accessor-mutation-does-not-change-copied-body",
    "facade-input-copy-does-not-alias-passive-request",
    "success-reply-accessor-mutation-does-not-change-retained-state-or-repeated-read",
  ],
  responseLoss: [
    {
      method: "SubmitMainMJSV0",
      disposition: methods.submitMainMJSV0.responseLossDisposition,
      downstreamDisposition: methods.registerPlanV0.responseLossDisposition,
      retryMayCreateFreshRegistration: true,
      requestIdIsIdempotencyKey: false,
    },
    {
      method: "RegisterPlanV0",
      disposition: methods.registerPlanV0.responseLossDisposition,
      firstCommittedRegistrations: 1,
      retryCommittedRegistrations: 2,
      bothRegistrationsSeparatelyReadable: true,
      requestIdIsIdempotencyKey: false,
    },
    {
      method: "GetRegisteredPlanV0",
      disposition: methods.getRegisteredPlanV0.responseLossDisposition,
      storeMutation: false,
      repeatedReplyBodyByteEqual: true,
      requestIdIsIdempotencyKey: false,
    },
  ],
  flowControl: [
    {
      id: "submit.connection.concurrent-cap-plus-one",
      method: "SubmitMainMJSV0",
      admitted: 1,
      attempted: 2,
      sameConnection: true,
      classification: "CAPACITY",
      reason: "ipc-connection-in-flight",
      queueDepth: 0,
      noState: zeroState,
    },
    {
      id: "submit.service-connections-cap-plus-one",
      method: "SubmitMainMJSV0",
      admitted: 4,
      attempted: 5,
      classification: "CAPACITY",
      reason: "ipc-service-connections",
      queueDepth: 0,
      noState: zeroState,
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
      noState: zeroState,
    },
    {
      id: "supervisor.mixed-register-get-aggregate-byte-cap-plus-one",
      methods: [
        { method: "RegisterPlanV0", admittedBytes: caps.registerPlanV0Request },
        { method: "GetRegisteredPlanV0", admittedBytes: caps.getRegisteredPlanV0Request },
      ],
      isolatedAccountingCapBytes: caps.registerPlanV0Request + caps.getRegisteredPlanV0Request,
      attemptedAdditionalBytes: 1,
      classification: "CAPACITY",
      reason: "ipc-process-in-flight-bytes",
      releasedMethod: "GetRegisteredPlanV0",
      postReleaseReadmissionBytes: 1,
      postReleaseReadmission: "admitted",
      productionAggregateCapBytes: 2_626_696,
      queueDepth: 0,
      noState: zeroState,
    },
  ],
  cancellationAndDeadline: [
    {
      id: "cancel.before-dispatch",
      disposition: "caller-cancelled-before-dispatch",
      coreCalls: 0,
      noState: zeroState,
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
      id: "downstream.stall",
      disposition: "method-deadline-after-dispatch-response-unknown",
      newWorkOnSameConnection: "CAPACITY",
      queueDepth: 0,
      processTerminationEvidence: false,
    },
  ],
  sourceCustody: {
    proposalBytes: proposal.reference,
    extractedMainMJS: {
      logicalPath: "main.mjs",
      byteLength: proposalSource.length,
      sha256: proposalSourceSha256,
    },
    registeredSourceBytes: source.reference,
    sourceManifestBytes: sourceManifest.reference,
    exactProposalSourceEqualsRegisteredSource: true,
    registrationCommitsPlanBindingsManifestSourceAtomically: true,
    executeTimeReplacementBytesAccepted: false,
  },
  zeroEffects,
};

const expected = new Map([
  ["native-xpc-v0.contract.json", jsonBytes(nativeXPCContract)],
  ["submit-main-mjs-v0.method.json", jsonBytes(methods.submitMainMJSV0)],
  ["submit-main-mjs-v0.request.json", jsonBytes(submitRequest)],
  ["submit-main-mjs-v0.reply.json", jsonBytes(submitReply)],
  ["register-plan-v0.method.json", jsonBytes(methods.registerPlanV0)],
  ["register-plan-v0.request.json", jsonBytes(registerRequest)],
  ["register-plan-v0.reply.json", jsonBytes(registerReply)],
  ["get-registered-plan-v0.method.json", jsonBytes(methods.getRegisteredPlanV0)],
  ["get-registered-plan-v0.request.json", jsonBytes(getRequest)],
  ["get-registered-plan-v0.reply.json", jsonBytes(getReply)],
  ["oracles.json", jsonBytes(oracles)],
]);

const manifest = {
  objectType: "capsule.authenticated-local-ipc-passive-conformance",
  objectVersion: 0,
  status: "passive-unwired-native-contract-no-listener",
  protocolVersion,
  audience,
  submissionAudience,
  roleBindingRecordVersion,
  transportEncoding: nativeXPCEncoding,
  transportEncodingVersion: nativeXPCEncodingVersion,
  numericMessageTags: nativeXPCMessageTags,
  numericStatusTags: nativeXPCStatusTags,
  numericReasonTags: nativeXPCReasonTags,
  peerAuthenticationEvidence: null,
  caps,
  methodCount: 3,
  refusalCaseCount: oracles.refusals.length,
  maximumCaseCount: oracles.maxima.length,
  nativeXPCEnvelopeCount: Object.values(nativeXPCEnvelopes).length * 3,
  nativeXPCCaseCount: nativeXPCContract.cases.length,
  nativeXPCRefusalReplyCount: nativeXPCContract.refusalReplies.length,
  knownAnswers: Object.fromEntries(
    [...expected].map(([path, bytes]) => [path, reference(path, bytes)]),
  ),
  bodyFixtures: {
    exactJobProposalBytes: proposal.reference,
    exactPlanBytes: plan.reference,
    roleBindingBytes: bindings.reference,
    planRegistrationBytes: registration.reference,
    sourceManifestBytes: sourceManifest.reference,
    sourceBytes: source.reference,
  },
};
expected.set("manifest.json", jsonBytes(manifest));

if (process.argv.includes("--check")) {
  for (const [path, bytes] of expected) {
    const actual = await readFile(new URL(path, outputRoot));
    if (!actual.equals(bytes)) throw new Error(`stale authenticated-local-IPC fixture: ${path}`);
  }
  process.stdout.write("verified generated passive authenticated-local-IPC known answers\n");
} else {
  await mkdir(outputRoot, { recursive: true });
  for (const [path, bytes] of expected) await writeFile(new URL(path, outputRoot), bytes);
  process.stdout.write("generated passive authenticated-local-IPC known answers\n");
}

function maximumCase(
  id,
  method,
  direction,
  fieldLengths,
  applicationDataBytes,
  decision,
  classification,
  reason,
) {
  return {
    id,
    method,
    direction,
    byteSource: "deterministic-repeated-byte-per-field",
    fieldLengths,
    applicationDataBytes,
    expected: { decision, classification: classification ?? null, reason: reason ?? null },
    noState: decision === "reject" ? zeroState : null,
    postCoreState:
      decision === "local-integrity-fault" ? "not-claimed-by-passive-output-oracle" : null,
    processTerminationRequired: decision === "local-integrity-fault",
  };
}

function refusal(id, method, mutation, classification, reason) {
  return {
    id,
    method,
    mutation,
    expected: { decision: "reject", classification, reason },
    noState: zeroState,
  };
}

function nativeXPCEnvelope(
  method,
  direction,
  messageTag,
  exactKeyCount,
  applicationDataMaxBytes,
  fields,
) {
  if (fields.length !== exactKeyCount) throw new Error(`${method.method} ${direction} key drift`);
  return {
    method: method.method,
    direction,
    protocolVersion,
    methodVersion: method.methodVersion,
    messageTag,
    exactKeyCount,
    applicationDataMaxBytes,
    fields,
  };
}

function nativeXPCRequestHeader(method, tag) {
  return [
    fixedUInt64Field(nativeXPCKeys.protocolVersion, protocolVersion),
    fixedUInt64Field(nativeXPCKeys.methodVersion, method.methodVersion),
    fixedUInt64Field(nativeXPCKeys.messageTag, tag),
    dataField(nativeXPCKeys.requestId, 16, 16, false, true),
    dataField(nativeXPCKeys.installationId, 16, 16, false, false),
    { key: nativeXPCKeys.epochSequence, valueType: nativeXPCTypes.uint64 },
    dataField(nativeXPCKeys.epochDigest, 32, 32, false, false),
    fixedStringField(nativeXPCKeys.audience, method.audience),
    fixedStringField(nativeXPCKeys.purpose, method.purpose),
  ];
}

function nativeXPCReplyHeader(method, tag, success) {
  return [
    fixedUInt64Field(nativeXPCKeys.protocolVersion, protocolVersion),
    fixedUInt64Field(nativeXPCKeys.methodVersion, method.methodVersion),
    fixedUInt64Field(nativeXPCKeys.messageTag, tag),
    dataField(nativeXPCKeys.requestId, 16, 16, false, true),
    {
      key: nativeXPCKeys.status,
      valueType: nativeXPCTypes.uint64,
      fixedUInt64: success ? nativeXPCStatusTags.OK : null,
    },
    {
      key: nativeXPCKeys.reason,
      valueType: nativeXPCTypes.uint64,
      fixedUInt64: success ? nativeXPCReasonTags.none : null,
    },
  ];
}

function nativeXPCBodyFields(method) {
  switch (method) {
    case "SubmitMainMJSV0":
      return {
        request: [dataField(nativeXPCKeys.jobProposal, 1, caps.jobProposal, true, false)],
        reply: [dataField(nativeXPCKeys.registrationId, 16, 16, true, true)],
      };
    case "RegisterPlanV0":
      return {
        request: [
          dataField(nativeXPCKeys.executionPlan, 1, caps.executionPlan, true, false),
          dataField(nativeXPCKeys.roleBindings, caps.roleBindings, caps.roleBindings, true, false),
          dataField(nativeXPCKeys.sourceManifest, 87, caps.sourceManifest, true, false),
          dataField(nativeXPCKeys.source, 0, caps.source, true, false),
        ],
        reply: [dataField(nativeXPCKeys.planRegistration, 1, caps.planRegistration, true, false)],
      };
    case "GetRegisteredPlanV0":
      return {
        request: [dataField(nativeXPCKeys.registrationId, 16, 16, true, true)],
        reply: [
          dataField(nativeXPCKeys.executionPlan, 1, caps.executionPlan, true, false),
          dataField(nativeXPCKeys.roleBindings, caps.roleBindings, caps.roleBindings, true, false),
          dataField(nativeXPCKeys.planRegistration, 1, caps.planRegistration, true, false),
          dataField(nativeXPCKeys.sourceManifest, 87, caps.sourceManifest, true, false),
          dataField(nativeXPCKeys.source, 0, caps.source, true, false),
        ],
      };
    default:
      throw new Error(`unknown native XPC method ${method}`);
  }
}

function dataField(key, minimum, maximum, applicationData, nonZeroData) {
  return {
    key,
    valueType: nativeXPCTypes.data,
    minDataBytes: minimum,
    maxDataBytes: maximum,
    applicationData,
    nonZeroData,
  };
}

function fixedUInt64Field(key, value) {
  return { key, valueType: nativeXPCTypes.uint64, fixedUInt64: value };
}

function fixedStringField(key, value) {
  return { key, valueType: nativeXPCTypes.string, fixedString: value };
}

function nativeXPCCases() {
  const reject = (id, method, mutation, classification, reasonTag, context = {}) => ({
    id,
    method,
    direction: "request",
    mutation,
    ...context,
    expected: {
      decision: "reject-before-body-copy",
      classification,
      statusTag: nativeXPCStatusTags[classification],
      reasonTag,
      bodyCopied: false,
      coreCalls: 0,
    },
    noState: zeroState,
  });
  const crossServiceTag = (selectedMethod, receivedTag) =>
    reject(
      `cross-service.${selectedMethod}.tag-${receivedTag}`,
      selectedMethod,
      `entry-point=${selectedMethod};message-tag=${receivedTag}`,
      "UNSUPPORTED",
      nativeXPCReasonTags.messageTag,
      {
        dispatchIdentity: "service-entry-point-and-role-derived-before-tag-cross-check",
      },
    );
  return [
    {
      id: "all.exact-key-and-type-sets",
      methods: "all-frozen-methods",
      expected: "accept-outer-shape-only",
    },
    reject(
      "submit.body-cap-plus-one",
      "SubmitMainMJSV0",
      "job-proposal=2097153-bytes",
      "MALFORMED",
      nativeXPCReasonTags.dataCap,
    ),
    reject(
      "register.body-cap-plus-one",
      "RegisterPlanV0",
      "source=262145-bytes",
      "MALFORMED",
      nativeXPCReasonTags.dataCap,
    ),
    reject(
      "get.body-cap-plus-one",
      "GetRegisteredPlanV0",
      "registration-id=17-bytes",
      "SCHEMA",
      nativeXPCReasonTags.dataWidth,
    ),
    reject("all.missing-key", "all", "remove-purpose", "MALFORMED", nativeXPCReasonTags.keySet),
    reject("all.extra-key", "all", "add-capsule.extra", "MALFORMED", nativeXPCReasonTags.keySet),
    reject(
      "all.wrong-type",
      "all",
      "request-id=XPC_TYPE_STRING",
      "MALFORMED",
      nativeXPCReasonTags.valueType,
    ),
    reject(
      "all.nested-container",
      "all",
      "body=XPC_TYPE_DICTIONARY",
      "MALFORMED",
      nativeXPCReasonTags.valueType,
    ),
    reject(
      "all.zero-request-id",
      "all",
      "request-id=16-zero-bytes",
      "SCHEMA",
      nativeXPCReasonTags.zeroIdentifier,
    ),
    reject(
      "all.unknown-protocol-version",
      "all",
      "protocol-version=1",
      "UNSUPPORTED",
      nativeXPCReasonTags.protocolVersion,
    ),
    reject(
      "all.unknown-method-version",
      "all",
      "method-version=1",
      "UNSUPPORTED",
      nativeXPCReasonTags.methodVersion,
    ),
    reject(
      "register.joint-protocol-method-record-version-mismatch",
      "RegisterPlanV0",
      "protocol-version=1;method-version=1;role-bindings[0]=1",
      "UNSUPPORTED",
      nativeXPCReasonTags.protocolVersion,
      {
        mismatchSet: "protocol-version,method-version,role-binding-record-version",
        selectedFailure: "protocol-version",
        validationPrecedence: nativeXPCValidationPrecedence,
      },
    ),
    reject(
      "register.joint-method-record-version-mismatch",
      "RegisterPlanV0",
      "method-version=1;role-bindings[0]=1",
      "UNSUPPORTED",
      nativeXPCReasonTags.methodVersion,
      {
        mismatchSet: "method-version,role-binding-record-version",
        selectedFailure: "method-version",
        validationPrecedence: nativeXPCValidationPrecedence,
      },
    ),
    {
      id: "register.embedded-record-version-mismatch",
      method: "RegisterPlanV0",
      direction: "request",
      mutation: "role-bindings[0]=1",
      expected: {
        decision: "core-refusal-after-bounded-body-copy",
        classification: "UNSUPPORTED",
        statusTag: nativeXPCStatusTags.UNSUPPORTED,
        reasonTag: nativeXPCReasonTags.coreRefusal,
        bodyCopied: true,
        coreCalls: 1,
      },
    },
    crossServiceTag("SubmitMainMJSV0", "RegisterPlanV0"),
    crossServiceTag("SubmitMainMJSV0", "GetRegisteredPlanV0"),
    crossServiceTag("RegisterPlanV0", "SubmitMainMJSV0"),
    crossServiceTag("RegisterPlanV0", "GetRegisteredPlanV0"),
    crossServiceTag("GetRegisteredPlanV0", "SubmitMainMJSV0"),
    crossServiceTag("GetRegisteredPlanV0", "RegisterPlanV0"),
    reject(
      "all.wrong-service",
      "all",
      "selected-method-invoked-with-receiving-service=other-role-service",
      "AUTHENTICATION",
      nativeXPCReasonTags.methodBinding,
      {
        deliveryPrecondition:
          "peer-requirement-and-session-admitted-on-receiving-service;message-reached-method-binding-validator",
        replyDisposition: "delivered-in-band-refusal",
      },
    ),
    reject(
      "all.wrong-role",
      "all",
      "message-derived-role=other",
      "AUTHENTICATION",
      nativeXPCReasonTags.methodBinding,
      {
        deliveryPrecondition:
          "peer-requirement-and-session-admitted-on-receiving-service;message-reached-method-binding-validator",
        replyDisposition: "delivered-in-band-refusal",
      },
    ),
    reject(
      "all.wrong-audience",
      "all",
      "audience=other",
      "AUTHENTICATION",
      nativeXPCReasonTags.methodBinding,
    ),
    reject(
      "all.wrong-purpose",
      "all",
      "purpose=other",
      "AUTHENTICATION",
      nativeXPCReasonTags.methodBinding,
    ),
    reject(
      "all.local-audience-replaced-by-signed-object-audience",
      "all",
      `audience=${signedApprovalAudience}`,
      "AUTHENTICATION",
      nativeXPCReasonTags.methodBinding,
    ),
    reject(
      "all.local-purpose-replaced-by-signed-object-purpose",
      "all",
      `purpose=${signedApprovalPurpose}`,
      "AUTHENTICATION",
      nativeXPCReasonTags.methodBinding,
    ),
    {
      id: "all.request-id-replaced-by-installation-id-bytes",
      method: "all",
      direction: "request",
      mutation: "request-id=installation-id-bytes",
      expected: {
        decision: "accept-correlation-field-only",
        authorityGrantedByRequestId: false,
        requestIdIsIdempotencyKey: false,
      },
    },
    reject(
      "all.installation-id-replaced-by-request-id-bytes",
      "all",
      "installation-id=request-id-bytes",
      "BINDING",
      nativeXPCReasonTags.currentState,
    ),
    {
      id: "get.registration-id-replaced-by-request-id-bytes",
      method: "GetRegisteredPlanV0",
      direction: "request",
      mutation: "registration-id=request-id-bytes",
      expected: {
        decision: "continue-to-registration-lookup-without-transport-authority",
        registrationAuthorityGrantedByWidth: false,
      },
    },
    reject(
      "all.wrong-installation",
      "all",
      "installation-id=other",
      "BINDING",
      nativeXPCReasonTags.currentState,
    ),
    reject(
      "all.wrong-epoch-digest",
      "all",
      "epoch-digest=other",
      "BINDING",
      nativeXPCReasonTags.currentState,
    ),
    reject(
      "all.epoch-uint53-cap-plus-one",
      "all",
      "epoch-sequence=9007199254740992",
      "SCHEMA",
      nativeXPCReasonTags.epochSequence,
    ),
    reject(
      "all.extra-file-descriptor",
      "all",
      "add-capsule.smuggled-fd=XPC_TYPE_FD",
      "MALFORMED",
      nativeXPCReasonTags.keySet,
    ),
    reject(
      "all.extra-endpoint",
      "all",
      "add-capsule.smuggled-endpoint=XPC_TYPE_ENDPOINT",
      "MALFORMED",
      nativeXPCReasonTags.keySet,
    ),
    reject(
      "all.extra-mach-send-right",
      "all",
      "add-capsule.smuggled-mach-send-right=XPC_TYPE_MACH_SEND",
      "MALFORMED",
      nativeXPCReasonTags.keySet,
    ),
    {
      id: "all.success-reply-extra-key",
      methods: "all-frozen-methods",
      direction: "success-reply",
      mutation: "add-any-unknown-key",
      expected: {
        decision: "reject-reply",
        classification: "MALFORMED",
        statusTag: nativeXPCStatusTags.MALFORMED,
        reasonTag: nativeXPCReasonTags.keySet,
      },
    },
    {
      id: "all.success-reply-request-id-mismatch",
      methods: "all-frozen-methods",
      direction: "success-reply",
      mutation: "request-id=other-live-call",
      expected: {
        decision: "reject-reply",
        classification: "BINDING",
        statusTag: nativeXPCStatusTags.BINDING,
        reasonTag: nativeXPCReasonTags.currentState,
      },
    },
    {
      id: "all.refusal-reply-extra-body",
      methods: "all-frozen-methods",
      direction: "refusal-reply",
      mutation: "add-any-method-body-key",
      expected: {
        decision: "reject-reply",
        classification: "MALFORMED",
        statusTag: nativeXPCStatusTags.MALFORMED,
        reasonTag: nativeXPCReasonTags.keySet,
      },
    },
    {
      id: "all.local-integrity-output-fault",
      methods: "all-frozen-methods",
      direction: "success-reply",
      mutation: "oversize-short-write-pointer-length-or-bridge-version-mismatch",
      expected: {
        decision: "terminate-without-reply",
        statusTag: null,
        reasonTag: nativeXPCReasonTags.localIntegrityFault,
      },
    },
  ];
}

async function repositoryFixture(path) {
  const bytes = await readFile(new URL(path, repository));
  return { bytes, byteLength: bytes.length, reference: reference(path, bytes) };
}

function reference(path, bytes) {
  return { path, byteLength: bytes.length, sha256: sha256Hex(bytes) };
}

function jsonBytes(value) {
  return Buffer.from(`${JSON.stringify(value, null, 2)}\n`);
}

function repeatedHex(value, length) {
  return Buffer.alloc(length, value).toString("hex");
}
