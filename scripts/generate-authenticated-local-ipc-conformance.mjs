import { mkdir, readFile, writeFile } from "node:fs/promises";
import { sha256Hex } from "./lib/fixture-bytes.mjs";

const repository = new URL("../", import.meta.url);
const outputRoot = new URL("schemas/conformance/authenticated-local-ipc-v0/", repository);

const protocolVersion = 0;
const audience = "capsule.execution-supervisor.local.v0";
const submissionAudience = "capsule.daemon.local.v0";
const roleBindingRecordVersion = 0;
const authorityDisposition = "service-role-purpose-audience-derived";
const requestIdDisposition = "correlation-only";
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
  status: "passive-unwired-no-transport",
  protocolVersion,
  audience,
  submissionAudience,
  roleBindingRecordVersion,
  transportEncoding: null,
  numericMessageTags: null,
  peerAuthenticationEvidence: null,
  caps,
  methodCount: 3,
  refusalCaseCount: oracles.refusals.length,
  maximumCaseCount: oracles.maxima.length,
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
