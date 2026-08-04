import { createHash } from "node:crypto";

const knownEnvelopeBase64Url =
  "0oRYRKMBJgN4K2FwcGxpY2F0aW9uL2NhcHN1bGUuYXBwcm92YWwtZ3JhbnQrY2Jvcjt2PTAEUWFwcHJvdmFsLXRlc3Qta2V5oFjqrAF2Y2Fwc3VsZS5hcHByb3ZhbC1ncmFudAIAA1ARERERERERERERERERERERBFggIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIiIFUDMzMzMzMzMzMzMzMzMzMzMGWCBERERERERERERERERERERERERERERERERERERERERERAdQVVVVVVVVVVVVVVVVVVVVVQhQZmZmZmZmZmZmZmZmZmZmZgl0Y2Fwc3VsZS5wbGFuLmFwcHJvdmUKeBxjYXBzdWxlLmV4ZWN1dGlvbi1zdXBlcnZpc29yCxpqa-WADBpqa-asWEARuN_9scveZ7aMR7epQ4uMVlYSvxAIKm0T1Ep_sB4hhm58OtFfuXKw_qK98t5qWEoVuf4-WDSHbhNQNsDICdMk";
const knownPayloadHex =
  "ac017663617073756c652e617070726f76616c2d6772616e74020003501111111111111111111111111111111104582022222222222222222222222222222222222222222222222222222222222222220550333333333333333333333333333333330658204444444444444444444444444444444444444444444444444444444444444444075055555555555555555555555555555555085066666666666666666666666666666666097463617073756c652e706c616e2e617070726f76650a781c63617073756c652e657865637574696f6e2d73757065727669736f720b1a6a6be5800c1a6a6be6ac";
const knownProtectedHex =
  "a3012603782b6170706c69636174696f6e2f63617073756c652e617070726f76616c2d6772616e742b63626f723b763d300451617070726f76616c2d746573742d6b6579";

const installationId = repeatedBytes(0x11, 16);
const epochDigest = repeatedBytes(0x22, 32);
const registrationId = repeatedBytes(0x33, 16);
const planDigest = repeatedBytes(0x44, 32);
const supervisorId = repeatedBytes(0x55, 16);
const attemptNonce = repeatedBytes(0x66, 16);
const emptySetDigest = sha256Hex(Buffer.from("[]\n"));
const sharedNonzeroId = Uint8Array.from({ length: 16 }, (_, index) => index + 1);

export function addApprovalAttemptRulesAndCases({
  addCase,
  addRule,
  cborEncode,
  jsonDocument,
  retainFixture,
}) {
  const implementations = {
    "fixture-integrity": "verified",
    go: "verified",
    typescript: "not-applicable",
    swift: "pending",
  };
  const emptyState = retainFixture(
    "contexts/approval-attempt/empty.json",
    jsonDocument(approvalAttemptState()),
  );

  const addStateCase = ({
    id,
    description,
    ruleIds,
    object,
    variant,
    path,
    bytes,
    operation,
    payload = null,
    protectedHeader = null,
    decision = "accept",
    classification = null,
    owner = "approval-attempt-contract",
  }) => {
    const operationFixture = retainFixture(
      `contexts/approval-attempt/${id.replaceAll(".", "-")}.operation.json`,
      jsonDocument({
        contextType: "capsule.conformance.approval-attempt-operation",
        contextVersion: 0,
        ...operation,
      }),
    );
    addCase({
      id,
      description,
      ruleIds,
      object,
      wireFormat: "raw-bytes",
      mediaType: object === "ApprovalGrant" ? "application/capsule.approval-grant+cbor;v=0" : null,
      variant,
      path,
      bytes,
      context: {
        kind: "approval-attempt-state",
        operation: operationFixture,
        before: emptyState,
        payload,
        protectedHeader,
      },
      decision,
      classification,
      owner,
      implementations,
      authorityStateChanged: false,
      timeHighWaterChanged: false,
      trustStateTightened: false,
      fakeBackendEffectPermitted: false,
      stateDelta: { kind: "exact", after: emptyState },
    });
  };

  const identifierRequirements = [
    { decision: "accept", variant: "nonzero" },
    { decision: "reject", variant: "all-zero" },
    { decision: "reject", variant: "wrong-domain" },
  ];
  for (const [domain, object] of [
    ["approval", "ApprovalID"],
    ["attempt", "AttemptID"],
    ["attempt-nonce", "AttemptNonce"],
  ]) {
    const ruleId = `approval-attempt.identifier.${domain}`;
    addRule(
      ruleId,
      "ADR-0024#identifier-and-nonce-domains",
      `${object} accepts only a nonzero exact 16-byte value in its own semantic domain.`,
      identifierRequirements,
    );
    addStateCase({
      id: `${ruleId}.nonzero`,
      description: `Accept one exact nonzero ${object} fixture in its own role.`,
      ruleIds: [ruleId],
      object,
      variant: "nonzero",
      path: "shared/id-nonzero-16.bin",
      bytes: sharedNonzeroId,
      operation: identifierOperation(domain, domain),
    });
    for (const [suffix, path, bytes] of [
      ["zero", "shared/id-zero-16.bin", new Uint8Array(16)],
      ["short", "shared/id-short.bin", new Uint8Array(15)],
      ["long", "shared/id-long.bin", new Uint8Array(17)],
    ]) {
      addStateCase({
        id: `${ruleId}.reject-${suffix}`,
        description: `Reject the ${suffix} ${object} representation.`,
        ruleIds: [ruleId],
        object,
        variant: suffix === "zero" ? "all-zero" : "malformed",
        path,
        bytes,
        operation: identifierOperation(domain, domain),
        decision: "reject",
        classification: "SCHEMA",
      });
    }
    for (const providedDomain of ["approval", "attempt", "attempt-nonce"].filter(
      (candidate) => candidate !== domain,
    )) {
      addStateCase({
        id: `${ruleId}.reject-${providedDomain}`,
        description: `Reject correct-width ${providedDomain} bytes in the ${domain} role.`,
        ruleIds: [ruleId],
        object,
        variant: "wrong-domain",
        path: "shared/id-nonzero-16.bin",
        bytes: sharedNonzeroId,
        operation: identifierOperation(domain, providedDomain),
        decision: "reject",
        classification: "DOMAIN",
      });
    }
  }

  for (const [reference, domain] of [
    ["approval-reference", "approval"],
    ["attempt-reference", "attempt"],
  ]) {
    const object = reference === "approval-reference" ? "ApprovalReference" : "AttemptReference";
    const ruleId = `approval-attempt.reference.${reference}`;
    addRule(
      ruleId,
      "ADR-0024#boundary-and-authority",
      `${object} is only the typed local projection of its matching nonzero identifier.`,
      [
        { decision: "accept", variant: "ordinary" },
        { decision: "reject", variant: "wrong-domain" },
      ],
    );
    addStateCase({
      id: `${ruleId}.ordinary`,
      description: `Accept a ${object} constructed from its matching identifier role.`,
      ruleIds: [ruleId],
      object,
      variant: "ordinary",
      path: "shared/id-nonzero-16.bin",
      bytes: sharedNonzeroId,
      operation: referenceOperation(reference, domain),
    });
    addStateCase({
      id: `${ruleId}.reject-wrong-domain`,
      description: `Reject a cross-domain identifier before constructing ${object}.`,
      ruleIds: [ruleId],
      object,
      variant: "wrong-domain",
      path: "shared/id-nonzero-16.bin",
      bytes: sharedNonzeroId,
      operation: referenceOperation(reference, domain === "approval" ? "attempt" : "approval"),
      decision: "reject",
      classification: "DOMAIN",
    });
  }

  const vocabulary = [
    "AUTHENTICATION",
    "MALFORMED",
    "UNSUPPORTED",
    "SCHEMA",
    "DOMAIN",
    "BINDING",
    "STALE",
    "REPLAY",
    "CAPACITY",
    "TRUST_STATE",
    "LOCAL_FAILURE",
    "RECOVERY_REQUIRED",
  ];
  const vocabularyRule = "approval-attempt.classification.fixed-vocabulary";
  addRule(
    vocabularyRule,
    "ADR-0024#fixed-internal-classifications",
    "The approval/attempt slice uses only the fixed content-free internal classification vocabulary.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "reject", variant: "malformed" },
    ],
  );
  addStateCase({
    id: `${vocabularyRule}.ordinary`,
    description: "Accept the complete ordered fixed classification vocabulary.",
    ruleIds: [vocabularyRule],
    object: "approval-state",
    variant: "ordinary",
    path: "approval-attempt/classifications.json",
    bytes: jsonDocument(vocabulary),
    operation: { mode: "classification-vocabulary" },
  });
  addStateCase({
    id: `${vocabularyRule}.reject-incomplete`,
    description: "Reject an incomplete classification vocabulary.",
    ruleIds: [vocabularyRule],
    object: "attempt-state",
    variant: "malformed",
    path: "approval-attempt/classifications-incomplete.json",
    bytes: jsonDocument(vocabulary.slice(0, -1)),
    operation: { mode: "classification-vocabulary" },
    decision: "reject",
    classification: "SCHEMA",
  });

  const ordinary = vectorFixtures({
    name: "ordinary",
    envelope: Buffer.from(knownEnvelopeBase64Url, "base64url"),
    payload: Buffer.from(knownPayloadHex, "hex"),
    protectedHeader: Buffer.from(knownProtectedHex, "hex"),
    keyId: Buffer.from("approval-test-key"),
    view: grantView(),
  });
  const maximumKeyId = repeatedBytes(0x6b, 64);
  const maximumPayload = cborEncode(
    grantPayload({
      issuedAt: Number.MAX_SAFE_INTEGER - 1,
      expiresAt: Number.MAX_SAFE_INTEGER,
    }),
  );
  const maximumProtected = cborEncode(protectedMap(maximumKeyId));
  const maximum = vectorFixtures({
    name: "calculated-maximum",
    envelope: sign1Envelope(maximumProtected, maximumPayload),
    payload: maximumPayload,
    protectedHeader: maximumProtected,
    keyId: maximumKeyId,
    view: grantView({ issuedAt: Number.MAX_SAFE_INTEGER - 1, expiresAt: Number.MAX_SAFE_INTEGER }),
  });
  if (
    maximum.envelopeBytes.length !== 431 ||
    maximum.payloadBytes.length !== 242 ||
    maximum.protectedBytes.length !== 116
  ) {
    throw new Error("approval calculated maxima no longer match ADR-0024");
  }
  const payloadOverBytes = repeatedBytes(0, 257);
  const payloadOver = vectorFixtures({
    name: "payload-cap-plus-one",
    envelope: sign1Envelope(Buffer.from(knownProtectedHex, "hex"), payloadOverBytes),
    payload: payloadOverBytes,
    protectedHeader: Buffer.from(knownProtectedHex, "hex"),
    keyId: Buffer.from("approval-test-key"),
    view: grantView(),
  });
  const protectedOverBytes = repeatedBytes(0, 129);
  const protectedOver = vectorFixtures({
    name: "protected-cap-plus-one",
    envelope: sign1Envelope(protectedOverBytes, Uint8Array.of(0xa0)),
    payload: Uint8Array.of(0xa0),
    protectedHeader: protectedOverBytes,
    keyId: Buffer.from("approval-test-key"),
    view: grantView(),
  });
  const profiles = new Map([
    [ordinary.name, ordinary],
    [maximum.name, maximum],
    [payloadOver.name, payloadOver],
    [protectedOver.name, protectedOver],
  ]);
  for (const mutation of ["object-type", "object-version", "purpose", "audience"]) {
    const view = grantView();
    if (mutation === "object-type") view.objectType = "capsule.execution-attempt";
    if (mutation === "object-version") view.objectVersion = 1;
    if (mutation === "purpose") view.purpose = "capsule.execution.attest";
    if (mutation === "audience") view.audience = "capsule.agent-daemon";
    const payload = cborEncode(grantPayload(view));
    const profile = vectorFixtures({
      name: mutation,
      envelope: sign1Envelope(Buffer.from(knownProtectedHex, "hex"), payload),
      payload,
      protectedHeader: Buffer.from(knownProtectedHex, "hex"),
      keyId: Buffer.from("approval-test-key"),
      view,
    });
    profiles.set(profile.name, profile);
  }
  for (const profile of profiles.values()) {
    retainVectorFixtures(profile, retainFixture);
  }

  const bindingRule = "approval-grant.fixture-verifier.bindings";
  addRule(
    bindingRule,
    "ADR-0024#exact-signed-and-resolved-bindings",
    "The fixture-only verifier accepts only the exact candidate payload and complete signed/resolved role bindings.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "reject", variant: "wrong-domain" },
    ],
  );
  addVerifierCase(addStateCase, ordinary, {
    id: `${bindingRule}.ordinary-known-answer`,
    description:
      "Accept the retained exact Gate A2 candidate known answer with every resolved binding.",
    ruleIds: [bindingRule],
    variant: "ordinary",
  });
  for (const mutation of [
    "installation",
    "epoch-sequence",
    "epoch-digest",
    "registration",
    "plan-digest",
    "supervisor",
    "attempt-nonce",
    "protected-key-id",
    "authorization-identity",
  ]) {
    addVerifierCase(addStateCase, ordinary, {
      id: `${bindingRule}.reject-${mutation}`,
      description: `Reject the exact known answer under a wrong ${mutation} resolved binding.`,
      ruleIds: [bindingRule],
      variant: "wrong-domain",
      operation: { bindingMutation: mutation },
      decision: "reject",
      classification: "BINDING",
    });
  }
  for (const [mutation, classification] of [
    ["object-type", "UNSUPPORTED"],
    ["object-version", "UNSUPPORTED"],
    ["purpose", "BINDING"],
    ["audience", "BINDING"],
  ]) {
    addVerifierCase(addStateCase, profiles.get(mutation), {
      id: `${bindingRule}.reject-${mutation}`,
      description: `Reject a fixture with the wrong signed ${mutation}.`,
      ruleIds: [bindingRule],
      variant: "malformed",
      decision: "reject",
      classification,
    });
  }

  const copyRule = "approval-grant.fixture-verifier.copy-ownership";
  addRule(
    copyRule,
    "ADR-0024#exact-signed-and-resolved-bindings",
    "The verifier copies the received envelope and every returned byte slice.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "accept", variant: "nonzero" },
    ],
  );
  addVerifierCase(addStateCase, ordinary, {
    id: `${copyRule}.ordinary`,
    description: "Return copied exact bytes for the ordinary known answer.",
    ruleIds: [copyRule],
    variant: "ordinary",
  });
  addVerifierCase(addStateCase, ordinary, {
    id: `${copyRule}.caller-and-accessor-mutation`,
    description: "Caller and accessor mutations cannot alter verifier-owned exact bytes.",
    ruleIds: [copyRule],
    variant: "nonzero",
    operation: { callerMutation: true },
  });

  const budgetRules = [
    ["approval-grant.fixture-verifier.envelope-bytes", "complete tagged envelope", 512],
    ["approval-grant.fixture-verifier.payload-bytes", "embedded canonical payload", 256],
    ["approval-grant.fixture-verifier.protected-bytes", "encoded protected-header map", 128],
  ];
  for (const [id, label, maximumBytes] of budgetRules) {
    addRule(
      id,
      "ADR-0024#exact-signed-and-resolved-bindings",
      `The ${label} is bounded at ${maximumBytes} raw bytes before ordinary fixture verification, ` +
        "and separately at the calculated closed-candidate maximum below that raw cap.",
      [
        { decision: "accept", variant: "exact-maximum" },
        { decision: "reject", variant: "cap-plus-one" },
        { decision: "reject", variant: "calculated-maximum-plus-one" },
        { decision: "reject", variant: "raw-maximum" },
      ],
    );
  }
  addVerifierCase(addStateCase, maximum, {
    id: "approval-grant.fixture-verifier.calculated-maxima",
    description: "Accept the calculated closed-candidate maxima of 431/242/116 bytes.",
    ruleIds: budgetRules.map(([id]) => id),
    variant: "exact-maximum",
  });
  addStateCase({
    id: "approval-grant.fixture-verifier.envelope-cap-plus-one",
    description: "Reject a 513-byte envelope before fixture lookup or authority change.",
    ruleIds: [budgetRules[0][0]],
    object: "ApprovalGrant",
    variant: "cap-plus-one",
    path: "approval-grant/envelope-cap-plus-one.cose",
    bytes: repeatedBytes(0, 513),
    operation: verifierOperation("envelope-cap-plus-one"),
    decision: "reject",
    classification: "MALFORMED",
  });
  addVerifierCase(addStateCase, payloadOver, {
    id: "approval-grant.fixture-verifier.payload-cap-plus-one",
    description: "Reject a 257-byte embedded payload before candidate validation.",
    ruleIds: [budgetRules[1][0]],
    variant: "cap-plus-one",
    decision: "reject",
    classification: "MALFORMED",
  });
  addVerifierCase(addStateCase, protectedOver, {
    id: "approval-grant.fixture-verifier.protected-cap-plus-one",
    description: "Reject a 129-byte protected map before candidate validation.",
    ruleIds: [budgetRules[2][0]],
    variant: "cap-plus-one",
    decision: "reject",
    classification: "MALFORMED",
  });

  // The calculated closed-candidate maxima (431/242/116, see `maximum` above) are
  // strictly smaller than the raw predecoder gates (512/256/128): no fixture above
  // exercised the raw gate's own true boundary. These six vectors are synthetic
  // filler bytes, not realistic COSE/CBOR structures -- FixtureVerifier is a
  // byte-equality fixture lookup, not a parser, so envelope/payload/protected
  // content is irrelevant to it and only byte length matters here. Each vector
  // isolates one dimension at its true boundary while holding the other two well
  // under their own calculated maxima, so the SCHEMA rejection it produces is
  // attributable to exactly that dimension.
  const envelopeCalculatedMaximumPlusOne = vectorFixtures({
    name: "envelope-calculated-maximum-plus-one",
    envelope: repeatedBytes(0, 432),
    payload: repeatedBytes(0, 1),
    protectedHeader: repeatedBytes(0, 1),
    keyId: Buffer.from("approval-test-key"),
    view: grantView(),
  });
  const envelopeRawMaximum = vectorFixtures({
    name: "envelope-raw-maximum",
    envelope: repeatedBytes(0, 512),
    payload: repeatedBytes(0, 1),
    protectedHeader: repeatedBytes(0, 1),
    keyId: Buffer.from("approval-test-key"),
    view: grantView(),
  });
  const payloadCalculatedMaximumPlusOne = vectorFixtures({
    name: "payload-calculated-maximum-plus-one",
    envelope: repeatedBytes(0, 8),
    payload: repeatedBytes(0, 243),
    protectedHeader: repeatedBytes(0, 1),
    keyId: Buffer.from("approval-test-key"),
    view: grantView(),
  });
  const payloadRawMaximum = vectorFixtures({
    name: "payload-raw-maximum",
    envelope: repeatedBytes(0, 8),
    payload: repeatedBytes(0, 256),
    protectedHeader: repeatedBytes(0, 1),
    keyId: Buffer.from("approval-test-key"),
    view: grantView(),
  });
  const protectedCalculatedMaximumPlusOne = vectorFixtures({
    name: "protected-calculated-maximum-plus-one",
    envelope: repeatedBytes(0, 8),
    payload: repeatedBytes(0, 1),
    protectedHeader: repeatedBytes(0, 117),
    keyId: Buffer.from("approval-test-key"),
    view: grantView(),
  });
  const protectedRawMaximum = vectorFixtures({
    name: "protected-raw-maximum",
    envelope: repeatedBytes(0, 8),
    payload: repeatedBytes(0, 1),
    protectedHeader: repeatedBytes(0, 128),
    keyId: Buffer.from("approval-test-key"),
    view: grantView(),
  });
  for (const profile of [
    envelopeCalculatedMaximumPlusOne,
    envelopeRawMaximum,
    payloadCalculatedMaximumPlusOne,
    payloadRawMaximum,
    protectedCalculatedMaximumPlusOne,
    protectedRawMaximum,
  ]) {
    retainVectorFixtures(profile, retainFixture);
  }
  addVerifierCase(addStateCase, envelopeCalculatedMaximumPlusOne, {
    id: "approval-grant.fixture-verifier.envelope-calculated-maximum-plus-one",
    description:
      "Reject a 432-byte envelope: within the raw 512-byte gate but one over the calculated 431-byte candidate maximum.",
    ruleIds: [budgetRules[0][0]],
    variant: "calculated-maximum-plus-one",
    decision: "reject",
    classification: "SCHEMA",
  });
  addVerifierCase(addStateCase, envelopeRawMaximum, {
    id: "approval-grant.fixture-verifier.envelope-raw-maximum",
    description:
      "Reject an exactly-512-byte envelope: passes the raw byte gate and reaches the calculated-candidate-maximum rejection.",
    ruleIds: [budgetRules[0][0]],
    variant: "raw-maximum",
    decision: "reject",
    classification: "SCHEMA",
  });
  addVerifierCase(addStateCase, payloadCalculatedMaximumPlusOne, {
    id: "approval-grant.fixture-verifier.payload-calculated-maximum-plus-one",
    description:
      "Reject a 243-byte embedded payload: within the raw 256-byte gate but one over the calculated 242-byte candidate maximum.",
    ruleIds: [budgetRules[1][0]],
    variant: "calculated-maximum-plus-one",
    decision: "reject",
    classification: "SCHEMA",
  });
  addVerifierCase(addStateCase, payloadRawMaximum, {
    id: "approval-grant.fixture-verifier.payload-raw-maximum",
    description:
      "Reject an exactly-256-byte embedded payload: passes the raw byte gate and reaches the calculated-candidate-maximum rejection.",
    ruleIds: [budgetRules[1][0]],
    variant: "raw-maximum",
    decision: "reject",
    classification: "SCHEMA",
  });
  addVerifierCase(addStateCase, protectedCalculatedMaximumPlusOne, {
    id: "approval-grant.fixture-verifier.protected-calculated-maximum-plus-one",
    description:
      "Reject a 117-byte protected map: within the raw 128-byte gate but one over the calculated 116-byte candidate maximum.",
    ruleIds: [budgetRules[2][0]],
    variant: "calculated-maximum-plus-one",
    decision: "reject",
    classification: "SCHEMA",
  });
  addVerifierCase(addStateCase, protectedRawMaximum, {
    id: "approval-grant.fixture-verifier.protected-raw-maximum",
    description:
      "Reject an exactly-128-byte protected map: passes the raw byte gate and reaches the calculated-candidate-maximum rejection.",
    ruleIds: [budgetRules[2][0]],
    variant: "raw-maximum",
    decision: "reject",
    classification: "SCHEMA",
  });

  addDurableStoreCases({
    addCase,
    addRule,
    emptyState,
    implementations,
    jsonDocument,
    ordinary,
    retainFixture,
  });
}

function addDurableStoreCases({
  addCase,
  addRule,
  emptyState,
  implementations,
  jsonDocument,
  ordinary,
  retainFixture,
}) {
  const state = (name, overrides = {}) =>
    retainFixture(
      `contexts/approval-attempt/store-${name}.json`,
      jsonDocument(approvalAttemptState(overrides)),
    );
  const usable = state("usable", {
    approvalPopulation: population("usableCount", 1, 1, "usable"),
    materializedApprovals: [{ approvalIdHex: "a1".padEnd(32, "0"), state: "usable" }],
  });
  const consumed = state("consumed-created", {
    approvalPopulation: population("usableCount", 0, 1, "consumed"),
    attemptPopulation: population("nonterminalCount", 1, 1, "created"),
    materializedApprovals: [{ approvalIdHex: "a1".padEnd(32, "0"), state: "consumed" }],
    materializedAttempts: [{ attemptIdHex: "b2".padEnd(32, "0"), state: "created" }],
  });
  const recoveryFenced = state("recovery-fenced", { recoveryFence: true });
  const usableRecoveryFenced = state("usable-recovery-fenced", {
    recoveryFence: true,
    approvalPopulation: population("usableCount", 1, 1, "usable"),
    materializedApprovals: [{ approvalIdHex: "a1".padEnd(32, "0"), state: "usable" }],
  });
  const approvalCapacity = state("approval-capacity", {
    approvalPopulation: population("usableCount", 256, 256, "approval-capacity"),
  });
  const attemptCapacity = state("attempt-capacity", {
    approvalPopulation: population("usableCount", 1, 257, "attempt-approval-capacity"),
    attemptPopulation: population("nonterminalCount", 256, 256, "attempt-capacity"),
  });

  const addStoreCase = ({
    id,
    description,
    ruleIds,
    method,
    scenario,
    before,
    after,
    decision = "accept",
    classification = null,
    variant = "ordinary",
    authorityStateChanged = false,
    timeHighWaterChanged = false,
    recovery = false,
  }) => {
    const operation = retainFixture(
      `contexts/approval-attempt/${id.replaceAll(".", "-")}.operation.json`,
      jsonDocument({
        contextType: "capsule.conformance.approval-attempt-operation",
        contextVersion: 0,
        mode: "store-transition",
        method,
        scenario,
        recovery,
      }),
    );
    addCase({
      id,
      description,
      ruleIds,
      object: method === "submit-approval" ? "approval-state" : "attempt-state",
      wireFormat: "raw-bytes",
      mediaType: null,
      variant,
      path: `approval-grant/${ordinary.name}.cose`,
      bytes: ordinary.envelopeBytes,
      context: {
        kind: "approval-attempt-state",
        operation,
        before,
        payload: ordinary.payloadFixture,
        protectedHeader: ordinary.protectedFixture,
      },
      decision,
      classification,
      owner: "approval-attempt-fixed-store",
      implementations,
      authorityStateChanged,
      timeHighWaterChanged,
      trustStateTightened: false,
      fakeBackendEffectPermitted: false,
      stateDelta: { kind: "exact", after },
    });
  };

  addRule(
    "approval-store.submission-idempotency",
    "ADR-0024#approval-submission-durable-states-and-replay",
    "One canonical payload creates one durable usable record; exact replay returns it without mutation.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "reject", variant: "malformed" },
    ],
  );
  addStoreCase({
    id: "approval-store.submission.first-usable",
    description: "Commit one verified candidate as one usable approval.",
    ruleIds: ["approval-store.submission-idempotency"],
    method: "submit-approval",
    scenario: "first-usable",
    before: emptyState,
    after: usable,
    authorityStateChanged: true,
  });
  addStoreCase({
    id: "approval-store.submission.exact-replay",
    description: "Return the same usable reference for an exact canonical-payload replay.",
    ruleIds: ["approval-store.submission-idempotency"],
    method: "submit-approval",
    scenario: "exact-replay",
    before: usable,
    after: usable,
  });
  addStoreCase({
    id: "approval-store.submission.nonce-replay-reject",
    description: "Reject a different payload that reuses a retained nonce.",
    ruleIds: ["approval-store.submission-idempotency"],
    method: "submit-approval",
    scenario: "nonce-replay",
    before: usable,
    after: usable,
    decision: "reject",
    classification: "REPLAY",
    variant: "malformed",
  });

  addRule(
    "attempt-store.atomic-consume-create",
    "ADR-0024#atomic-consume-and-attempt-creation",
    "One transaction consumes one usable approval and creates exactly one immutable attempt.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "reject", variant: "malformed" },
    ],
  );
  addStoreCase({
    id: "attempt-store.consume-create.atomic",
    description: "Atomically consume one approval and insert its one created attempt.",
    ruleIds: ["attempt-store.atomic-consume-create"],
    method: "request-attempt",
    scenario: "atomic-commit",
    before: usable,
    after: consumed,
    authorityStateChanged: true,
  });
  addStoreCase({
    id: "attempt-store.consume-create.exact-replay",
    description: "Return the same attempt for an exact request replay without redriving effects.",
    ruleIds: ["attempt-store.atomic-consume-create"],
    method: "request-attempt",
    scenario: "exact-replay",
    before: consumed,
    after: consumed,
  });
  addStoreCase({
    id: "attempt-store.consume-create.confirmed-abort",
    description: "A confirmed transaction abort leaves the approval usable and creates no attempt.",
    ruleIds: ["attempt-store.atomic-consume-create"],
    method: "request-attempt",
    scenario: "confirmed-abort",
    before: usable,
    after: usable,
    decision: "reject",
    classification: "LOCAL_FAILURE",
    variant: "malformed",
  });

  addRule(
    "approval-attempt-store.recovery-fence",
    "ADR-0024#request-replay-concurrency-and-process-death",
    "An indeterminate commit fences mutation until reopen validates one coherent pre-state or post-state.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "reject", variant: "malformed" },
    ],
  );
  addStoreCase({
    id: "approval-store.recovery.indeterminate-before-record",
    description: "Fence an indeterminate approval commit without granting authority.",
    ruleIds: ["approval-attempt-store.recovery-fence"],
    method: "submit-approval",
    scenario: "indeterminate-pre-state",
    before: emptyState,
    after: recoveryFenced,
    decision: "reject",
    classification: "RECOVERY_REQUIRED",
    variant: "malformed",
    recovery: true,
  });
  addStoreCase({
    id: "approval-store.recovery.indeterminate-after-record",
    description: "Fence an indeterminate approval commit whose complete record reached the store.",
    ruleIds: ["approval-attempt-store.recovery-fence"],
    method: "submit-approval",
    scenario: "indeterminate-post-state",
    before: emptyState,
    after: usableRecoveryFenced,
    decision: "reject",
    classification: "RECOVERY_REQUIRED",
    variant: "malformed",
    recovery: true,
  });
  addStoreCase({
    id: "approval-store.recovery.reopen-usable",
    description:
      "Reopen validates the committed usable record and clears only the transient fence.",
    ruleIds: ["approval-attempt-store.recovery-fence"],
    method: "submit-approval",
    scenario: "reopen-post-state",
    before: usableRecoveryFenced,
    after: usable,
    recovery: true,
  });

  addRule(
    "approval-attempt-store.capacity",
    "ADR-0024#bounded-first-store-and-retention",
    "Exact approval and attempt ceilings reject cap-plus-one without eviction or consumption.",
    [
      { decision: "accept", variant: "exact-maximum" },
      { decision: "reject", variant: "cap-plus-one" },
    ],
  );
  addStoreCase({
    id: "approval-store.capacity.usable-exact",
    description: "Retain exactly 256 usable approvals without eviction.",
    ruleIds: ["approval-attempt-store.capacity"],
    method: "submit-approval",
    scenario: "usable-exact-256",
    before: approvalCapacity,
    after: approvalCapacity,
    variant: "exact-maximum",
  });
  addStoreCase({
    id: "approval-store.capacity.usable-cap-plus-one",
    description: "Reject the 257th usable approval without eviction.",
    ruleIds: ["approval-attempt-store.capacity"],
    method: "submit-approval",
    scenario: "usable-cap-plus-one",
    before: approvalCapacity,
    after: approvalCapacity,
    decision: "reject",
    classification: "CAPACITY",
    variant: "cap-plus-one",
  });
  addStoreCase({
    id: "attempt-store.capacity.nonterminal-cap-plus-one",
    description: "Reject the 257th nonterminal attempt and leave its approval usable.",
    ruleIds: ["approval-attempt-store.capacity"],
    method: "request-attempt",
    scenario: "nonterminal-cap-plus-one",
    before: attemptCapacity,
    after: attemptCapacity,
    decision: "reject",
    classification: "CAPACITY",
    variant: "cap-plus-one",
  });
}

function population(liveField, liveCount, retainedCount, label) {
  return {
    [liveField]: liveCount,
    retainedCount,
    setDigest: sha256Hex(Buffer.from(label)),
  };
}

function addVerifierCase(addStateCase, profile, options) {
  const { operation: operationOverrides, ...caseOptions } = options;
  addStateCase({
    object: "ApprovalGrant",
    path: `approval-grant/${profile.name}.cose`,
    bytes: profile.envelopeBytes,
    payload: profile.payloadFixture,
    protectedHeader: profile.protectedFixture,
    operation: verifierOperation(profile.name, operationOverrides),
    ...caseOptions,
  });
}

function retainVectorFixtures(profile, retainFixture) {
  profile.payloadFixture = retainFixture(
    `approval-grant/${profile.name}.payload.cbor`,
    profile.payloadBytes,
  );
  profile.protectedFixture = retainFixture(
    `approval-grant/${profile.name}.protected.cbor`,
    profile.protectedBytes,
  );
}

function vectorFixtures({ name, envelope, payload, protectedHeader, keyId, view }) {
  return {
    name,
    envelopeBytes: Buffer.from(envelope),
    payloadBytes: Buffer.from(payload),
    protectedBytes: Buffer.from(protectedHeader),
    keyIdHex: hex(keyId),
    view,
  };
}

function identifierOperation(expectedDomain, providedDomain) {
  return { mode: "identifier", expectedDomain, providedDomain };
}

function referenceOperation(referenceKind, providedDomain) {
  return { mode: "reference", referenceKind, providedDomain };
}

function verifierOperation(vector, overrides = {}) {
  return {
    mode: "fixture-verifier",
    vector,
    bindingMutation: "none",
    callerMutation: false,
    ...overrides,
  };
}

function grantView({
  objectType = "capsule.approval-grant",
  objectVersion = 0,
  purpose = "capsule.plan.approve",
  audience = "capsule.execution-supervisor",
  issuedAt = 1_785_456_000,
  expiresAt = 1_785_456_300,
} = {}) {
  return { objectType, objectVersion, purpose, audience, issuedAt, expiresAt };
}

function grantPayload(view = grantView()) {
  return new Map([
    [1, view.objectType ?? "capsule.approval-grant"],
    [2, view.objectVersion ?? 0],
    [3, installationId],
    [4, epochDigest],
    [5, registrationId],
    [6, planDigest],
    [7, supervisorId],
    [8, attemptNonce],
    [9, view.purpose ?? "capsule.plan.approve"],
    [10, view.audience ?? "capsule.execution-supervisor"],
    [11, view.issuedAt ?? 1_785_456_000],
    [12, view.expiresAt ?? 1_785_456_300],
  ]);
}

function protectedMap(keyId) {
  return new Map([
    [1, -7],
    [3, "application/capsule.approval-grant+cbor;v=0"],
    [4, keyId],
  ]);
}

function sign1Envelope(protectedHeader, payload) {
  return Buffer.concat([
    Buffer.from([0xd2, 0x84]),
    encodeByteString(protectedHeader),
    Buffer.from([0xa0]),
    encodeByteString(payload),
    encodeByteString(new Uint8Array(64)),
  ]);
}

function encodeByteString(value) {
  const bytes = Buffer.from(value);
  if (bytes.length < 24) return Buffer.concat([Buffer.from([0x40 | bytes.length]), bytes]);
  if (bytes.length <= 0xff) return Buffer.concat([Buffer.from([0x58, bytes.length]), bytes]);
  if (bytes.length <= 0xffff) {
    return Buffer.concat([Buffer.from([0x59, bytes.length >> 8, bytes.length & 0xff]), bytes]);
  }
  throw new Error("fixture byte string is too large");
}

function approvalAttemptState(overrides = {}) {
  return {
    contextType: "capsule.conformance.approval-attempt-state",
    contextVersion: 0,
    installationIdHex: hex(installationId),
    supervisorIdHex: hex(supervisorId),
    epochSequence: 7,
    epochDigestHex: hex(epochDigest),
    trustPhase: "stable",
    trustReason: null,
    attemptsEnabled: true,
    recoveryFence: false,
    timeHighWaterUnixSeconds: 1_785_456_000,
    approvalPopulation: {
      usableCount: 0,
      retainedCount: 0,
      setDigest: emptySetDigest,
    },
    attemptPopulation: {
      nonterminalCount: 0,
      retainedCount: 0,
      setDigest: emptySetDigest,
    },
    materializedApprovals: [],
    materializedAttempts: [],
    ...overrides,
  };
}

function repeatedBytes(value, length) {
  return Uint8Array.from({ length }, () => value);
}

function hex(value) {
  return Buffer.from(value).toString("hex");
}

function sha256Hex(value) {
  return createHash("sha256").update(value).digest("hex");
}
