import { createHash } from "node:crypto";

const maxUInt53 = Number.MAX_SAFE_INTEGER;
const effectiveNow = 1_785_456_000;
const installationId = repeatedBytes(0x11, 16);
const supervisorId = repeatedBytes(0x55, 16);
const registrationIdA = repeatedBytes(0x77, 16);
const registrationIdB = repeatedBytes(0x78, 16);
const epochDigest = repeatedBytes(0x22, 32);
const sourceManifestDigest = hexBytes(
  "c387c80094027ffbcacb573f44f5f6b4dec4d243bb436b24dd644434feaa1d14",
);
const inlineInputDigest = hexBytes(
  "bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72",
);
const emptyRecordSetDigest = sha256Hex(Buffer.from("[]\n"));

export function addPlanRegistrationRulesAndCases({
  addCase,
  addRule,
  cborEncode,
  jsonDocument,
  retainFixture,
  scalarRoleContext,
}) {
  const passiveImplementations = {
    "fixture-integrity": "verified",
    go: "verified",
    typescript: "pending",
    swift: "pending",
  };
  const stateImplementations = {
    "fixture-integrity": "verified",
    go: "verified",
    typescript: "pending",
    swift: "pending",
  };

  const plan = executionPlan();
  const planBytes = cborEncode(plan);
  const planDigest = sha256Bytes(planBytes);
  const registrationBytes = cborEncode(
    planRegistration({ registrationId: registrationIdA, sequence: 1, planDigest }),
  );

  addRules(addRule);

  const addPassiveCase = (options) =>
    addCase({
      object: "ExecutionPlan",
      wireFormat: "cbor",
      mediaType: "application/capsule.execution-plan+cbor;v=0",
      owner: "role-binding-validator",
      implementations: passiveImplementations,
      ...options,
    });

  const addStateCase = ({
    id,
    description,
    ruleId,
    bytes = planBytes,
    path = "execution-plan/ordinary.cbor",
    object = "ExecutionPlan",
    mediaType = "application/capsule.execution-plan+cbor;v=0",
    owner = "supervisor-registration-state",
    variant = "ordinary",
    before = registrationState(),
    after,
    operation = registrationOperation(),
    decision = "accept",
    classification = null,
    authorityStateChanged = decision === "accept" && operation.method === "register-plan",
  }) => {
    const slug = id.replaceAll(".", "-");
    const operationFixture = retainFixture(
      `contexts/registration/${slug}.operation.json`,
      jsonDocument(operation),
    );
    const beforeFixture = retainFixture(
      `contexts/registration/${slug}.before.json`,
      jsonDocument(before),
    );
    const afterFixture = retainFixture(
      `contexts/registration/${slug}.after.json`,
      jsonDocument(after ?? before),
    );
    addCase({
      id,
      description,
      ruleIds: [ruleId],
      object,
      wireFormat: "cbor",
      mediaType,
      variant,
      path,
      bytes,
      context: { kind: "registration-state", operation: operationFixture, before: beforeFixture },
      decision,
      classification,
      owner,
      implementations: stateImplementations,
      authorityStateChanged,
      stateDelta: { kind: "exact", after: afterFixture },
    });
  };

  const baseRecord = storedRecord({
    wireBytes: registrationBytes,
    planBytes,
    registeredAtUnixSeconds: effectiveNow,
  });
  const baseBefore = registrationState();
  const baseAfter = appendRecord(baseBefore, baseRecord, { sequence: 1 });

  addStateCase({
    id: "execution-plan.registration.exact-bytes",
    description:
      "Register the exact integrated plan bytes and recompute their SHA-256 identity before commit.",
    ruleId: "execution-plan.registration.exact-bytes",
    before: baseBefore,
    after: baseAfter,
  });
  addStateCase({
    id: "execution-plan.registration.mutation-after-submission",
    description: "Caller mutation after submission cannot change retained plan bytes or digest.",
    ruleId: "execution-plan.registration.exact-bytes",
    variant: "nonzero",
    before: baseBefore,
    after: baseAfter,
    operation: registrationOperation({ mutation: "caller-buffer-after-submission" }),
  });
  addStateCase({
    id: "execution-plan.registration.validator-copy-mutation",
    description: "Mutation of the validator's private copy cannot change retained authority.",
    ruleId: "execution-plan.registration.exact-bytes",
    variant: "exact-maximum",
    before: baseBefore,
    after: baseAfter,
    operation: registrationOperation({ mutation: "validator-private-copy" }),
  });

  addPassiveCase({
    id: "execution-plan.domain.roles-match",
    description: "Accept exact plan bytes whose role-specific identifiers and digests match.",
    ruleIds: ["execution-plan.registration.domain-separation"],
    variant: "ordinary",
    path: "execution-plan/ordinary.cbor",
    bytes: planBytes,
    context: scalarRoleContext("source-manifest", "source-manifest"),
  });
  addPassiveCase({
    id: "execution-plan.domain.reject-plan-registration-object",
    description: "Reject a valid PlanRegistration object at the ExecutionPlan entry point.",
    ruleIds: ["execution-plan.registration.domain-separation"],
    variant: "wrong-domain",
    path: "execution-plan/plan-registration-object.cbor",
    bytes: registrationBytes,
    decision: "reject",
    classification: "DOMAIN",
  });
  addCase({
    id: "plan-registration.domain.reject-execution-plan-object",
    description: "Reject a valid ExecutionPlan object at the PlanRegistration entry point.",
    ruleIds: ["execution-plan.registration.domain-separation"],
    object: "PlanRegistration",
    wireFormat: "cbor",
    mediaType: "application/capsule.plan-registration+cbor;v=0",
    variant: "wrong-domain",
    path: "plan-registration/execution-plan-object.cbor",
    bytes: planBytes,
    decision: "reject",
    classification: "DOMAIN",
    owner: "role-binding-validator",
    implementations: passiveImplementations,
  });

  for (const entry of domainCases({ cborEncode, planDigest })) {
    addCase({
      ...entry,
      ruleIds: ["execution-plan.registration.domain-separation"],
      wireFormat: "cbor",
      mediaType:
        entry.object === "ExecutionPlan"
          ? "application/capsule.execution-plan+cbor;v=0"
          : "application/capsule.plan-registration+cbor;v=0",
      variant: "wrong-domain",
      decision: "reject",
      classification: "DOMAIN",
      owner: "role-binding-validator",
      implementations: passiveImplementations,
    });
  }

  addStateCase({
    id: "plan-registration.binding.plan-a-registration-b",
    description: "A correct-role registration B cannot substitute for registration A and plan A.",
    ruleId: "execution-plan.registration.exact-bytes",
    object: "PlanRegistration",
    mediaType: "application/capsule.plan-registration+cbor;v=0",
    bytes: cborEncode(
      planRegistration({ registrationId: registrationIdB, sequence: 2, planDigest }),
    ),
    path: "plan-registration/plan-a-registration-b.cbor",
    variant: "wrong-domain",
    before: baseAfter,
    after: baseAfter,
    operation: registrationOperation({ method: "use-registration", identifier: null }),
    decision: "reject",
    classification: "BINDING",
    authorityStateChanged: false,
  });

  addStateCase({
    id: "plan-registration.channel.authenticated-daemon",
    description: "Only the authenticated daemon registration method reaches stateful admission.",
    ruleId: "plan-registration.channel.caller-role",
    owner: "channel-authentication",
    before: baseBefore,
    after: baseAfter,
  });
  addStateCase({
    id: "plan-registration.creation.first-fresh-registration",
    description: "The first successful call creates a fresh nonzero registration at sequence one.",
    ruleId: "plan-registration.creation.fresh",
    before: baseBefore,
    after: baseAfter,
  });
  for (const [suffix, caller] of [
    ["unauthenticated", { authenticated: false, role: "daemon", purpose: "register-plan" }],
    ["broker", { authenticated: true, role: "broker", purpose: "register-plan" }],
    ["updater", { authenticated: true, role: "updater", purpose: "register-plan" }],
  ]) {
    addStateCase({
      id: `plan-registration.channel.reject-${suffix}`,
      description: `Reject the ${suffix} caller before trusted-time persistence.`,
      ruleId: "plan-registration.channel.caller-role",
      owner: "channel-authentication",
      variant: "wrong-domain",
      before: baseBefore,
      after: baseBefore,
      operation: registrationOperation({
        caller,
        trustedClockObservationUnixSeconds: effectiveNow + 10,
      }),
      decision: "reject",
      classification: "AUTHENTICATION",
      authorityStateChanged: false,
    });
  }

  const duplicateWire = cborEncode(
    planRegistration({ registrationId: registrationIdB, sequence: 2, planDigest }),
  );
  const duplicateRecord = storedRecord({
    wireBytes: duplicateWire,
    planBytes,
    registeredAtUnixSeconds: effectiveNow,
  });
  addStateCase({
    id: "plan-registration.creation.duplicate-plan-bytes",
    description: "Identical plan bytes create a fresh ID and increasing registration sequence.",
    ruleId: "plan-registration.creation.fresh",
    variant: "nonzero",
    before: baseAfter,
    after: appendRecord(baseAfter, duplicateRecord, { sequence: 2 }),
    operation: registrationOperation({ identifier: identifierValue(registrationIdB) }),
  });

  const epochEightDigest = repeatedBytes(0x23, 32);
  const epochEightPlanBytes = cborEncode(
    executionPlan({ epochSequence: 8, epochDigest: epochEightDigest }),
  );
  const epochEightPlanDigest = sha256Bytes(epochEightPlanBytes);
  const epochEightWire = cborEncode(
    planRegistration({
      registrationId: registrationIdB,
      sequence: 2,
      planDigest: epochEightPlanDigest,
      epochSequence: 8,
      epochDigest: epochEightDigest,
    }),
  );
  const epochEightBefore = registrationState({
    epochSequence: 8,
    epochDigestHex: hex(epochEightDigest),
    lastRegistrationSequence: 1,
    storedRegistrationCount: 1,
    unexpiredRegistrationCount: 1,
    registrationSetDigest: populationDigest("old-epoch", 1, 1),
  });
  addStateCase({
    id: "plan-registration.sequence.continues-across-epoch",
    description: "The installation-global registration sequence does not reset at a trust epoch.",
    ruleId: "plan-registration.sequence.installation-global",
    path: "execution-plan/epoch-eight.cbor",
    bytes: epochEightPlanBytes,
    before: epochEightBefore,
    after: appendRecord(
      epochEightBefore,
      storedRecord({
        wireBytes: epochEightWire,
        planBytes: epochEightPlanBytes,
        registeredAtUnixSeconds: effectiveNow,
      }),
      { sequence: 2 },
    ),
    operation: registrationOperation({ identifier: identifierValue(registrationIdB) }),
  });

  const maxWire = cborEncode(
    planRegistration({ registrationId: registrationIdA, sequence: maxUInt53, planDigest }),
  );
  const beforeMaximum = registrationState({ lastRegistrationSequence: maxUInt53 - 1 });
  addStateCase({
    id: "plan-registration.sequence.uint53-maximum",
    description: "Issue the inclusive UInt53 maximum as a successful registration sequence.",
    ruleId: "plan-registration.sequence.installation-global",
    variant: "exact-maximum",
    before: beforeMaximum,
    after: appendRecord(
      beforeMaximum,
      storedRecord({
        wireBytes: maxWire,
        planBytes,
        registeredAtUnixSeconds: effectiveNow,
      }),
      { sequence: maxUInt53 },
    ),
  });
  const exhaustedBefore = registrationState({
    lastRegistrationSequence: maxUInt53,
    storedRegistrationCount: 1,
    unexpiredRegistrationCount: 1,
    registrationSetDigest: populationDigest("maximum-sequence", 1, 1),
  });
  addStateCase({
    id: "plan-registration.sequence.exhausted",
    description: "The request after UInt53 exhaustion fails stale and enters repair-required.",
    ruleId: "plan-registration.sequence.installation-global",
    variant: "cap-plus-one",
    before: exhaustedBefore,
    after: registrationState({
      ...stateValues(exhaustedBefore),
      trustPhase: "repair-required",
      trustReason: "registration-sequence-exhausted",
    }),
    operation: registrationOperation({ identifier: null }),
    decision: "reject",
    classification: "STALE",
    authorityStateChanged: false,
  });

  addExpiryCases({
    addStateCase,
    cborEncode,
    planBytes,
    baseBefore,
    baseAfter,
    registrationBytes,
  });
  addTimeAndTrustCases({
    addStateCase,
    cborEncode,
    planBytes,
    registrationBytes,
    baseBefore,
    baseAfter,
  });
  addCapacityCases({ addStateCase, cborEncode, planBytes, planDigest });
  addFailureCases({
    addStateCase,
    planBytes,
    registrationBytes,
    baseBefore,
    baseAfter,
  });

  const separationContext = retainFixture(
    "contexts/registration/wire-stored-separation.json",
    jsonDocument({
      contextType: "capsule.conformance.wire-stored-separation",
      contextVersion: 0,
      wireFields: [
        "objectType",
        "objectVersion",
        "registrationId",
        "registrationSequence",
        "planDigest",
        "installationId",
        "epochSequence",
        "epochDigest",
        "supervisorId",
        "expiresAt",
      ],
      storedRecordFields: [
        "wireRegistrationHex",
        "exactPlanHex",
        "recomputedPlanDigestHex",
        "registeredAtUnixSeconds",
        "storageFormatVersion",
        "retentionState",
      ],
    }),
  );
  addCase({
    id: "plan-registration.storage.wire-omits-stored-fields",
    description:
      "Wire PlanRegistration contains no plan bytes, registration time, or recovery data.",
    ruleIds: ["plan-registration.storage.separation"],
    object: "PlanRegistration",
    wireFormat: "cbor",
    mediaType: "application/capsule.plan-registration+cbor;v=0",
    variant: "ordinary",
    path: "plan-registration/ordinary.cbor",
    bytes: registrationBytes,
    context: { kind: "fixture", fixture: separationContext },
    owner: "supervisor-registration-predecoder",
    implementations: passiveImplementations,
  });
  addStateCase({
    id: "plan-registration.storage.minimum-stored-record",
    description:
      "Commit the exact closed six-field stored record separately from the wire payload.",
    ruleId: "plan-registration.storage.separation",
    variant: "exact-maximum",
    before: baseBefore,
    after: baseAfter,
  });
}

function addRules(addRule) {
  const ordinary = [
    { decision: "accept", variant: "ordinary" },
    { decision: "reject", variant: "wrong-domain" },
  ];
  addRule(
    "execution-plan.registration.exact-bytes",
    "ADR-0023#registration-creation-replay-sequence-and-expiry",
    "Registration copies, hashes, and stores the exact submitted plan bytes.",
    ordinary,
  );
  addRule(
    "execution-plan.registration.domain-separation",
    "ADR-0023#scalar-zero-rules",
    "Valid objects, identifiers, and digests cannot substitute across nominal roles.",
    ordinary,
  );
  addRule(
    "plan-registration.channel.caller-role",
    "ADR-0023#boundary-ownership-and-media-types",
    "Only the authenticated daemon registration method accepts candidate plan bytes.",
    ordinary,
  );
  addRule(
    "plan-registration.creation.fresh",
    "ADR-0023#registration-creation-replay-sequence-and-expiry",
    "Every successful call creates a fresh nonzero registration, including duplicate plan bytes.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "accept", variant: "nonzero" },
    ],
  );
  addRule(
    "plan-registration.sequence.installation-global",
    "ADR-0023#proposed-task-24-registration-oracle-addendum",
    "The inclusive UInt53 sequence is installation-global and exhaustion fails stale.",
    [
      { decision: "accept", variant: "exact-maximum" },
      { decision: "reject", variant: "cap-plus-one" },
    ],
  );
  addRule(
    "plan-registration.expiry.lifetime",
    "ADR-0023#registration-creation-replay-sequence-and-expiry",
    "Registration requires effectiveNow < expiresAt <= effectiveNow + 300 exactly.",
    [
      { decision: "accept", variant: "exact-maximum" },
      { decision: "reject", variant: "cap-plus-one" },
    ],
  );
  addRule(
    "plan-registration.time.monotonic",
    "ADR-0023#proposed-task-24-registration-oracle-addendum",
    "Effective time uses a durable nondecreasing Unix-seconds high-water mark.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "reject", variant: "malformed" },
    ],
  );
  addRule(
    "plan-registration.trust.gate",
    "ADR-0023#proposed-task-24-registration-oracle-addendum",
    "Epoch fencing, quarantine, and repair fail closed at their distinct state oracle.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "reject", variant: "malformed" },
    ],
  );
  for (const [suffix, description] of [
    ["unexpired", "At most 256 unexpired registrations are admitted."],
    ["stored", "At most 4,096 stored registration records are retained."],
  ]) {
    addRule(
      `plan-registration.capacity.${suffix}`,
      "ADR-0023#registration-creation-replay-sequence-and-expiry",
      description,
      [
        { decision: "accept", variant: "exact-maximum" },
        { decision: "reject", variant: "cap-plus-one" },
      ],
    );
  }
  addRule(
    "plan-registration.failure.atomicity",
    "ADR-0023#proposed-task-24-registration-oracle-addendum",
    "Local prerequisite and confirmed-abort failures create no registration or sequence.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "reject", variant: "malformed" },
    ],
  );
  addRule(
    "plan-registration.storage.separation",
    "ADR-0023#proposed-task-24-registration-oracle-addendum",
    "Wire registration bytes remain separate from the minimum exact stored record.",
    [
      { decision: "accept", variant: "ordinary" },
      { decision: "accept", variant: "exact-maximum" },
    ],
  );
}

function addExpiryCases({
  addStateCase,
  cborEncode,
  planBytes,
  baseBefore,
  baseAfter,
  registrationBytes,
}) {
  const cases = [
    {
      suffix: "equality-expired",
      expiresAt: effectiveNow,
      variant: "malformed",
      decision: "reject",
      classification: "STALE",
      description: "Reject expiry equality because effectiveNow < expiresAt is strict.",
    },
    {
      suffix: "one-second",
      expiresAt: effectiveNow + 1,
      variant: "ordinary",
      description: "Accept the minimum one-second remaining lifetime unchanged.",
    },
    {
      suffix: "exact-300-seconds",
      expiresAt: effectiveNow + 300,
      variant: "exact-maximum",
      description: "Accept the exact 300-second lifetime and copy expiry unchanged.",
    },
    {
      suffix: "301-seconds",
      expiresAt: effectiveNow + 301,
      variant: "cap-plus-one",
      decision: "reject",
      classification: "POLICY",
      description: "Reject 301 seconds rather than clamping or extending the lifetime.",
    },
  ];
  for (const entry of cases) {
    const candidatePlanBytes =
      entry.expiresAt === effectiveNow + 300
        ? planBytes
        : cborEncode(executionPlan({ expiresAt: entry.expiresAt }));
    const candidateDigest = sha256Bytes(candidatePlanBytes);
    const wire = cborEncode(
      planRegistration({
        registrationId: registrationIdA,
        sequence: 1,
        planDigest: candidateDigest,
        expiresAt: entry.expiresAt,
      }),
    );
    const accepted = entry.decision !== "reject";
    addStateCase({
      id: `plan-registration.expiry.${entry.suffix}`,
      description: entry.description,
      ruleId: "plan-registration.expiry.lifetime",
      path: `execution-plan/expires-${entry.suffix}.cbor`,
      bytes: candidatePlanBytes,
      variant: entry.variant,
      before: baseBefore,
      after: accepted
        ? appendRecord(
            baseBefore,
            storedRecord({
              wireBytes: wire,
              planBytes: candidatePlanBytes,
              registeredAtUnixSeconds: effectiveNow,
            }),
            { sequence: 1 },
          )
        : baseBefore,
      decision: entry.decision ?? "accept",
      classification: entry.classification ?? null,
      authorityStateChanged: accepted,
    });
  }

  const equalityBefore = registrationState({
    ...stateValues(baseAfter),
    timeHighWaterUnixSeconds: effectiveNow + 300,
    unexpiredRegistrationCount: 0,
  });
  addStateCase({
    id: "plan-registration.expiry.use-at-equality",
    description: "A stored registration is unusable when effectiveNow equals its expiry.",
    ruleId: "plan-registration.time.monotonic",
    object: "PlanRegistration",
    mediaType: "application/capsule.plan-registration+cbor;v=0",
    path: "plan-registration/ordinary.cbor",
    bytes: registrationBytes,
    variant: "malformed",
    before: equalityBefore,
    after: equalityBefore,
    operation: registrationOperation({
      method: "use-registration",
      identifier: null,
      trustedClockObservationUnixSeconds: effectiveNow + 300,
    }),
    decision: "reject",
    classification: "STALE",
    authorityStateChanged: false,
  });
}

function addTimeAndTrustCases({
  addStateCase,
  cborEncode,
  planBytes,
  registrationBytes,
  baseBefore,
  baseAfter,
}) {
  addStateCase({
    id: "plan-registration.time.initial-high-water",
    description: "A valid registration records the final effective Unix second.",
    ruleId: "plan-registration.time.monotonic",
    before: baseBefore,
    after: baseAfter,
  });

  const concurrentEffectiveState = registrationState({
    ...stateValues(baseBefore),
    timeHighWaterUnixSeconds: effectiveNow + 10,
  });
  const concurrentAfter = appendRecord(
    concurrentEffectiveState,
    storedRecord({
      wireBytes: registrationBytes,
      planBytes,
      registeredAtUnixSeconds: effectiveNow + 10,
    }),
    { sequence: 1 },
  );
  addStateCase({
    id: "plan-registration.time.final-reread-uses-later-high-water",
    description: "The final transaction re-read records a concurrent later high-water second.",
    ruleId: "plan-registration.time.monotonic",
    before: baseBefore,
    after: concurrentAfter,
    operation: registrationOperation({
      concurrentHighWaterUnixSeconds: effectiveNow + 10,
    }),
  });

  const concurrentExpiryAfter = registrationState({
    ...stateValues(baseBefore),
    timeHighWaterUnixSeconds: effectiveNow + 300,
  });
  addStateCase({
    id: "plan-registration.time.final-reread-can-expire-plan",
    description: "A concurrent high-water advance to expiry makes final admission stale.",
    ruleId: "plan-registration.time.monotonic",
    variant: "malformed",
    before: baseBefore,
    after: concurrentExpiryAfter,
    operation: registrationOperation({
      concurrentHighWaterUnixSeconds: effectiveNow + 300,
    }),
    decision: "reject",
    classification: "STALE",
    authorityStateChanged: false,
  });

  const rollbackBefore = registrationState({
    ...stateValues(baseAfter),
    timeHighWaterUnixSeconds: effectiveNow + 300,
    unexpiredRegistrationCount: 0,
  });
  addStateCase({
    id: "plan-registration.time.clock-rollback-cannot-resurrect",
    description: "A backward clock observation cannot reactivate an expired registration.",
    ruleId: "plan-registration.time.monotonic",
    object: "PlanRegistration",
    mediaType: "application/capsule.plan-registration+cbor;v=0",
    path: "plan-registration/ordinary.cbor",
    bytes: registrationBytes,
    variant: "malformed",
    before: rollbackBefore,
    after: rollbackBefore,
    operation: registrationOperation({
      method: "use-registration",
      identifier: null,
      trustedClockObservationUnixSeconds: effectiveNow,
    }),
    decision: "reject",
    classification: "STALE",
    authorityStateChanged: false,
  });

  const wrongEpochPlanBytes = cborEncode(
    executionPlan({ epochSequence: 6, epochDigest: repeatedBytes(0x21, 32) }),
  );
  const highWaterAfterBinding = registrationState({
    ...stateValues(baseBefore),
    timeHighWaterUnixSeconds: effectiveNow + 10,
  });
  addStateCase({
    id: "plan-registration.time.binding-reject-advances-high-water",
    description: "A stateful binding rejection may retain only the trusted time observation.",
    ruleId: "plan-registration.time.monotonic",
    path: "execution-plan/wrong-epoch-binding.cbor",
    bytes: wrongEpochPlanBytes,
    variant: "wrong-domain",
    before: baseBefore,
    after: highWaterAfterBinding,
    operation: registrationOperation({ trustedClockObservationUnixSeconds: effectiveNow + 10 }),
    decision: "reject",
    classification: "BINDING",
    authorityStateChanged: false,
  });

  addStateCase({
    id: "execution-plan.registration.invalid-bytes-no-time-write",
    description: "Malformed plan bytes reject before trusted-time persistence.",
    ruleId: "plan-registration.failure.atomicity",
    path: "execution-plan/invalid-break.cbor",
    bytes: Uint8Array.of(0xff),
    owner: "supervisor-plan-predecoder",
    variant: "malformed",
    before: baseBefore,
    after: baseBefore,
    operation: registrationOperation({ trustedClockObservationUnixSeconds: effectiveNow + 10 }),
    decision: "reject",
    classification: "MALFORMED",
    authorityStateChanged: false,
  });

  addStateCase({
    id: "plan-registration.time.high-water-write-failure",
    description: "High-water persistence failure stops admission without durable state change.",
    ruleId: "plan-registration.failure.atomicity",
    variant: "malformed",
    before: baseBefore,
    after: baseBefore,
    operation: registrationOperation({
      trustedClockObservationUnixSeconds: effectiveNow + 10,
      fault: "time-high-water-write",
      identifier: null,
    }),
    decision: "reject",
    classification: "LOCAL_FAILURE",
    authorityStateChanged: false,
  });

  addStateCase({
    id: "plan-registration.trust.stable",
    description: "Stable matching trust state admits the exact plan.",
    ruleId: "plan-registration.trust.gate",
    before: baseBefore,
    after: baseAfter,
  });
  for (const entry of [
    { suffix: "transition-fenced", trustPhase: "transition-fenced", classification: "STALE" },
    { suffix: "quarantined", quarantined: true, classification: "TRUST_STATE" },
    {
      suffix: "repair-required",
      trustPhase: "repair-required",
      trustReason: "incomplete-transition",
      classification: "TRUST_STATE",
    },
  ]) {
    const before = registrationState({
      trustPhase: entry.trustPhase ?? "stable",
      trustReason: entry.trustReason ?? null,
      quarantined: entry.quarantined ?? false,
    });
    const after = registrationState({
      ...stateValues(before),
      timeHighWaterUnixSeconds: effectiveNow + 10,
    });
    addStateCase({
      id: `plan-registration.trust.reject-${entry.suffix}`,
      description: `Reject registration while trust state is ${entry.suffix}.`,
      ruleId: "plan-registration.trust.gate",
      variant: "malformed",
      before,
      after,
      operation: registrationOperation({ trustedClockObservationUnixSeconds: effectiveNow + 10 }),
      decision: "reject",
      classification: entry.classification,
      authorityStateChanged: false,
    });
  }
}

function addCapacityCases({ addStateCase, cborEncode, planBytes, planDigest }) {
  const scenarios = [
    {
      id: "plan-registration.capacity.unexpired-exact-256",
      ruleId: "plan-registration.capacity.unexpired",
      variant: "exact-maximum",
      beforeCount: 255,
      beforeUnexpired: 255,
      afterCount: 256,
      afterUnexpired: 256,
      sequence: 256,
      identifier: repeatedBytes(0x90, 16),
      accept: true,
      description: "Admit the 256th unexpired registration exactly.",
    },
    {
      id: "plan-registration.capacity.unexpired-cap-plus-one",
      ruleId: "plan-registration.capacity.unexpired",
      variant: "cap-plus-one",
      beforeCount: 256,
      beforeUnexpired: 256,
      sequence: 257,
      identifier: null,
      accept: false,
      description: "Reject the 257th unexpired registration without eviction.",
    },
    {
      id: "plan-registration.capacity.expiry-equality-not-unexpired",
      ruleId: "plan-registration.capacity.unexpired",
      variant: "ordinary",
      beforeCount: 256,
      beforeUnexpired: 0,
      afterCount: 257,
      afterUnexpired: 1,
      sequence: 257,
      identifier: repeatedBytes(0x91, 16),
      accept: true,
      description: "Records expiring at effectiveNow do not consume the unexpired allowance.",
    },
    {
      id: "plan-registration.capacity.stored-exact-4096",
      ruleId: "plan-registration.capacity.stored",
      variant: "exact-maximum",
      beforeCount: 4095,
      beforeUnexpired: 0,
      afterCount: 4096,
      afterUnexpired: 1,
      sequence: 4096,
      identifier: repeatedBytes(0x92, 16),
      accept: true,
      description: "Admit the 4,096th stored record exactly.",
    },
    {
      id: "plan-registration.capacity.stored-cap-plus-one",
      ruleId: "plan-registration.capacity.stored",
      variant: "cap-plus-one",
      beforeCount: 4096,
      beforeUnexpired: 0,
      sequence: 4097,
      identifier: null,
      accept: false,
      description: "Reject record 4,097 even when all stored records are expired.",
    },
  ];

  for (const scenario of scenarios) {
    const beforeDigest = populationDigest(
      scenario.id,
      scenario.beforeCount,
      scenario.beforeUnexpired,
    );
    const before = registrationState({
      lastRegistrationSequence: scenario.sequence - 1,
      storedRegistrationCount: scenario.beforeCount,
      unexpiredRegistrationCount: scenario.beforeUnexpired,
      registrationSetDigest: beforeDigest,
    });
    let after = before;
    if (scenario.accept) {
      const wire = cborEncode(
        planRegistration({
          registrationId: scenario.identifier,
          sequence: scenario.sequence,
          planDigest,
        }),
      );
      after = appendRecord(
        before,
        storedRecord({ wireBytes: wire, planBytes, registeredAtUnixSeconds: effectiveNow }),
        {
          sequence: scenario.sequence,
          storedRegistrationCount: scenario.afterCount,
          unexpiredRegistrationCount: scenario.afterUnexpired,
        },
      );
    }
    addStateCase({
      id: scenario.id,
      description: scenario.description,
      ruleId: scenario.ruleId,
      variant: scenario.variant,
      before,
      after,
      operation: registrationOperation({
        identifier: scenario.identifier ? identifierValue(scenario.identifier) : null,
      }),
      decision: scenario.accept ? "accept" : "reject",
      classification: scenario.accept ? null : "CAPACITY",
      authorityStateChanged: scenario.accept,
    });
  }
}

function addFailureCases({ addStateCase, planBytes, registrationBytes, baseBefore, baseAfter }) {
  const timeOnlyAfter = registrationState({
    ...stateValues(baseBefore),
    timeHighWaterUnixSeconds: effectiveNow + 10,
  });
  const duplicateBefore = baseAfter;
  const duplicateAfter = registrationState({
    ...stateValues(duplicateBefore),
    timeHighWaterUnixSeconds: effectiveNow + 10,
  });
  for (const entry of [
    {
      suffix: "source-error",
      identifier: { kind: "failure", failure: "source-error" },
      before: baseBefore,
      after: timeOnlyAfter,
    },
    {
      suffix: "zero",
      identifier: identifierValue(new Uint8Array(16)),
      before: baseBefore,
      after: timeOnlyAfter,
    },
    {
      suffix: "duplicate",
      identifier: identifierValue(registrationIdA),
      before: duplicateBefore,
      after: duplicateAfter,
    },
  ]) {
    addStateCase({
      id: `plan-registration.identifier.${entry.suffix}`,
      description: `Identifier ${entry.suffix} fails locally without sequence consumption.`,
      ruleId: "plan-registration.failure.atomicity",
      path: "execution-plan/ordinary.cbor",
      bytes: planBytes,
      variant: "malformed",
      before: entry.before,
      after: entry.after,
      operation: registrationOperation({
        trustedClockObservationUnixSeconds: effectiveNow + 10,
        identifier: entry.identifier,
      }),
      decision: "reject",
      classification: "LOCAL_FAILURE",
      authorityStateChanged: false,
    });
  }
  addStateCase({
    id: "plan-registration.transaction.confirmed-abort",
    description: "A confirmed-abort commit exposes neither a sequence nor a partial record.",
    ruleId: "plan-registration.failure.atomicity",
    path: "execution-plan/ordinary.cbor",
    bytes: planBytes,
    variant: "malformed",
    before: baseBefore,
    after: timeOnlyAfter,
    operation: registrationOperation({
      trustedClockObservationUnixSeconds: effectiveNow + 10,
      fault: "registration-commit-confirmed-abort",
    }),
    decision: "reject",
    classification: "LOCAL_FAILURE",
    authorityStateChanged: false,
  });
  addStateCase({
    id: "plan-registration.transaction.atomic-commit",
    description: "A successful commit exposes its complete record and sequence together.",
    ruleId: "plan-registration.failure.atomicity",
    before: baseBefore,
    after: appendRecord(
      baseBefore,
      storedRecord({
        wireBytes: registrationBytes,
        planBytes,
        registeredAtUnixSeconds: effectiveNow,
      }),
      { sequence: 1 },
    ),
  });
}

function domainCases({ cborEncode, planDigest }) {
  const registrationBase = {
    registrationId: registrationIdA,
    sequence: 1,
    planDigest,
  };
  return [
    {
      id: "execution-plan.domain.reject-supervisor-as-installation",
      description: "Reject a Supervisor ID in the installation-ID role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-supervisor-as-installation.cbor",
      bytes: cborEncode(executionPlan({ installationId: supervisorId })),
      context: {
        kind: "scalar-role",
        scalarKind: "id",
        expectedRole: "installation",
        providedRole: "supervisor",
      },
    },
    {
      id: "execution-plan.domain.reject-inline-input-as-epoch",
      description: "Reject an inline-input digest in the trust-epoch role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-inline-input-as-epoch.cbor",
      bytes: cborEncode(executionPlan({ epochDigest: inlineInputDigest })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "trust-epoch",
        providedRole: "inline-input",
      },
    },
    {
      id: "execution-plan.domain.reject-inline-input-as-source-manifest",
      description: "Reject an inline-input digest in the source-manifest role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-inline-input-as-source.cbor",
      bytes: cborEncode(executionPlan({ sourceManifestDigest: inlineInputDigest })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "source-manifest",
        providedRole: "inline-input",
      },
    },
    {
      id: "execution-plan.domain.reject-source-manifest-as-inline-input",
      description: "Reject a source-manifest digest in the inline-input role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-source-manifest-as-inline-input.cbor",
      bytes: cborEncode(executionPlan({ inlineInputDigest: sourceManifestDigest })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "inline-input",
        providedRole: "source-manifest",
      },
    },
    {
      id: "execution-plan.domain.reject-profile-registry-as-runtime-bundle",
      description: "Reject a profile-registry-entry digest in the runtime-bundle role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-profile-registry-as-runtime-bundle.cbor",
      bytes: cborEncode(executionPlan({ runtimeBundleManifestDigest: repeatedBytes(0x77, 32) })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "runtime-bundle-manifest",
        providedRole: "profile-registry-entry",
      },
    },
    {
      id: "execution-plan.domain.reject-runtime-bundle-as-review-attestation",
      description: "Reject a runtime-bundle digest in the profile-review-attestation role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-runtime-bundle-as-review-attestation.cbor",
      bytes: cborEncode(
        executionPlan({ profileReviewAttestationDigests: [repeatedBytes(0x55, 32)] }),
      ),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "profile-review-attestation",
        providedRole: "runtime-bundle-manifest",
      },
    },
    {
      id: "execution-plan.domain.reject-review-attestation-as-profile-registry",
      description: "Reject a profile-review-attestation digest in the profile-registry-entry role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-review-attestation-as-profile-registry.cbor",
      bytes: cborEncode(executionPlan({ profileRegistryEntryDigest: repeatedBytes(0x66, 32) })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "profile-registry-entry",
        providedRole: "profile-review-attestation",
      },
    },
    {
      id: "execution-plan.domain.reject-backend-config-as-validation",
      description: "Reject a backend-configuration digest in the backend-validation-record role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-backend-config-as-validation.cbor",
      bytes: cborEncode(executionPlan({ backendValidationRecordDigest: repeatedBytes(0x99, 32) })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "backend-validation-record",
        providedRole: "backend-configuration",
      },
    },
    {
      id: "execution-plan.domain.reject-backend-validation-as-config",
      description: "Reject a backend-validation-record digest in the backend-configuration role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-backend-validation-as-config.cbor",
      bytes: cborEncode(executionPlan({ backendConfigurationDigest: repeatedBytes(0x88, 32) })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "backend-configuration",
        providedRole: "backend-validation-record",
      },
    },
    {
      id: "execution-plan.domain.reject-policy-as-trust-snapshot",
      description: "Reject a policy-decision digest in the trust-snapshot role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-policy-as-trust-snapshot.cbor",
      bytes: cborEncode(executionPlan({ trustSnapshotDigest: repeatedBytes(0xbb, 32) })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "trust-snapshot",
        providedRole: "policy-decision",
      },
    },
    {
      id: "execution-plan.domain.reject-trust-snapshot-as-policy",
      description: "Reject a trust-snapshot digest in the policy-decision role.",
      object: "ExecutionPlan",
      path: "execution-plan/domain-trust-snapshot-as-policy.cbor",
      bytes: cborEncode(executionPlan({ policyDecisionDigest: repeatedBytes(0xaa, 32) })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "policy-decision",
        providedRole: "trust-snapshot",
      },
    },
    {
      id: "plan-registration.domain.reject-installation-as-registration",
      description: "Reject an installation ID in the registration-ID role.",
      object: "PlanRegistration",
      path: "plan-registration/domain-installation-as-registration.cbor",
      bytes: cborEncode(planRegistration({ ...registrationBase, registrationId: installationId })),
      context: {
        kind: "scalar-role",
        scalarKind: "id",
        expectedRole: "registration",
        providedRole: "installation",
      },
    },
    {
      id: "plan-registration.domain.reject-epoch-as-plan-digest",
      description: "Reject a trust-epoch digest in the execution-plan digest role.",
      object: "PlanRegistration",
      path: "plan-registration/domain-epoch-as-plan-digest.cbor",
      bytes: cborEncode(planRegistration({ ...registrationBase, planDigest: epochDigest })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "execution-plan",
        providedRole: "trust-epoch",
      },
    },
    {
      id: "plan-registration.domain.reject-plan-as-epoch-digest",
      description: "Reject an execution-plan digest in the trust-epoch digest role.",
      object: "PlanRegistration",
      path: "plan-registration/domain-plan-as-epoch-digest.cbor",
      bytes: cborEncode(planRegistration({ ...registrationBase, epochDigest: planDigest })),
      context: {
        kind: "scalar-role",
        scalarKind: "digest",
        expectedRole: "trust-epoch",
        providedRole: "execution-plan",
      },
    },
    {
      id: "plan-registration.domain.reject-registration-as-supervisor",
      description: "Reject a registration ID in the Supervisor-ID role.",
      object: "PlanRegistration",
      path: "plan-registration/domain-registration-as-supervisor.cbor",
      bytes: cborEncode(planRegistration({ ...registrationBase, supervisorId: registrationIdA })),
      context: {
        kind: "scalar-role",
        scalarKind: "id",
        expectedRole: "supervisor",
        providedRole: "registration",
      },
    },
  ];
}

function executionPlan({
  installationId: planInstallationId = installationId,
  epochSequence = 7,
  epochDigest: planEpochDigest = epochDigest,
  sourceManifestDigest: planSourceManifestDigest = sourceManifestDigest,
  inlineInputDigest: planInlineInputDigest = inlineInputDigest,
  runtimeBundleManifestDigest = repeatedBytes(0x55, 32),
  profileReviewAttestationDigests = [repeatedBytes(0x66, 32), repeatedBytes(0x67, 32)],
  profileRegistryEntryDigest = repeatedBytes(0x77, 32),
  backendValidationRecordDigest = repeatedBytes(0x88, 32),
  backendConfigurationDigest = repeatedBytes(0x99, 32),
  trustSnapshotDigest = repeatedBytes(0xaa, 32),
  policyDecisionDigest = repeatedBytes(0xbb, 32),
  expiresAt = effectiveNow + 300,
} = {}) {
  return new Map([
    [1, "capsule.execution-plan"],
    [2, 0],
    [3, planInstallationId],
    [4, epochSequence],
    [5, planEpochDigest],
    [6, planSourceManifestDigest],
    [7, "main.mjs"],
    [8, 50],
    [9, "primary-data"],
    [10, planInlineInputDigest],
    [11, 118],
    [12, "fixture-active@1"],
    [13, runtimeBundleManifestDigest],
    [14, profileReviewAttestationDigests],
    [15, profileRegistryEntryDigest],
    [16, backendValidationRecordDigest],
    [17, backendConfigurationDigest],
    [18, trustSnapshotDigest],
    [19, policyDecisionDigest],
    [20, 5_000],
    [21, "requested"],
    [22, "transformed-json"],
    [23, 65_536],
    [24, expiresAt],
  ]);
}

function planRegistration({
  registrationId,
  sequence,
  planDigest,
  installationId: registrationInstallationId = installationId,
  epochSequence = 7,
  epochDigest: registrationEpochDigest = epochDigest,
  supervisorId: registrationSupervisorId = supervisorId,
  expiresAt = effectiveNow + 300,
}) {
  return new Map([
    [1, "capsule.plan-registration"],
    [2, 0],
    [3, registrationId],
    [4, sequence],
    [5, planDigest],
    [6, registrationInstallationId],
    [7, epochSequence],
    [8, registrationEpochDigest],
    [9, registrationSupervisorId],
    [10, expiresAt],
  ]);
}

function registrationOperation({
  method = "register-plan",
  caller = { authenticated: true, role: "daemon", purpose: "register-plan" },
  trustedClockObservationUnixSeconds = effectiveNow,
  concurrentHighWaterUnixSeconds = null,
  identifier = identifierValue(registrationIdA),
  mutation = "none",
  fault = null,
} = {}) {
  return {
    contextType: "capsule.conformance.registration-operation",
    contextVersion: 0,
    method,
    caller,
    trustedClockObservationUnixSeconds,
    concurrentHighWaterUnixSeconds,
    identifier,
    mutation,
    fault,
  };
}

function registrationState({
  epochSequence = 7,
  epochDigestHex = hex(epochDigest),
  trustPhase = "stable",
  trustReason = null,
  quarantined = false,
  timeHighWaterUnixSeconds = effectiveNow,
  lastRegistrationSequence = 0,
  storedRegistrationCount = 0,
  unexpiredRegistrationCount = 0,
  registrationSetDigest = emptyRecordSetDigest,
  materializedRecords = [],
} = {}) {
  return {
    contextType: "capsule.conformance.registration-state",
    contextVersion: 0,
    installationIdHex: hex(installationId),
    supervisorIdHex: hex(supervisorId),
    epochSequence,
    epochDigestHex,
    trustPhase,
    trustReason,
    quarantined,
    timeHighWaterUnixSeconds,
    lastRegistrationSequence,
    registrationPopulation: {
      storedCount: storedRegistrationCount,
      unexpiredCount: unexpiredRegistrationCount,
      setDigest: registrationSetDigest,
    },
    materializedRecords,
  };
}

function stateValues(state) {
  return {
    epochSequence: state.epochSequence,
    epochDigestHex: state.epochDigestHex,
    trustPhase: state.trustPhase,
    trustReason: state.trustReason,
    quarantined: state.quarantined,
    timeHighWaterUnixSeconds: state.timeHighWaterUnixSeconds,
    lastRegistrationSequence: state.lastRegistrationSequence,
    storedRegistrationCount: state.registrationPopulation.storedCount,
    unexpiredRegistrationCount: state.registrationPopulation.unexpiredCount,
    registrationSetDigest: state.registrationPopulation.setDigest,
    materializedRecords: state.materializedRecords,
  };
}

function storedRecord({ wireBytes, planBytes, registeredAtUnixSeconds }) {
  return {
    wireRegistrationHex: hex(wireBytes),
    exactPlanHex: hex(planBytes),
    recomputedPlanDigestHex: sha256Hex(planBytes),
    registeredAtUnixSeconds,
    storageFormatVersion: 0,
    retentionState: "retained",
  };
}

function appendRecord(
  state,
  record,
  {
    sequence,
    storedRegistrationCount = state.registrationPopulation.storedCount + 1,
    unexpiredRegistrationCount = state.registrationPopulation.unexpiredCount + 1,
  },
) {
  return registrationState({
    ...stateValues(state),
    lastRegistrationSequence: sequence,
    storedRegistrationCount,
    unexpiredRegistrationCount,
    registrationSetDigest: sha256Hex(
      Buffer.from(`${state.registrationPopulation.setDigest}:${JSON.stringify(record)}`),
    ),
    materializedRecords: [...state.materializedRecords, record],
  });
}

function populationDigest(label, count, unexpiredCount) {
  return sha256Hex(Buffer.from(`${label}:${count}:${unexpiredCount}`));
}

function identifierValue(value) {
  return { kind: "value", valueHex: hex(value) };
}

function repeatedBytes(value, length) {
  return Uint8Array.from({ length }, () => value);
}

function hexBytes(value) {
  return Uint8Array.from(Buffer.from(value, "hex"));
}

function hex(value) {
  return Buffer.from(value).toString("hex");
}

function sha256Bytes(value) {
  return createHash("sha256").update(value).digest();
}

function sha256Hex(value) {
  return createHash("sha256").update(value).digest("hex");
}
