import { createHash } from "node:crypto";

const protocolVersion = 0;
const requestKind = 1;
const resultKind = 2;
const method = 1;
const validatorProfile = 1;
const sourceProfile = 1;
const sourceMedia = 1;
const correlationDomain = 1;
const sourceDigestDomain = 0x0101;
const candidateDigestDomain = 0x0102;
const executableDigestDomain = 0x0103;
const buildDigestDomain = 0x0104;
const assessmentDigestDomain = 0x0105;
const artifactProfileDomain = 0x0106;
const requestHeaderBytes = 80;
const sourceMaximumBytes = 262_144;
const resultFrameBytes = 138;
const candidateFrameBytes = 292;
const artifactProfileFrameBytes = 160;
const mediaType = "application/capsule.source-validator-frame;v=0";

const implementations = {
  "fixture-integrity": "verified",
  go: "verified",
  typescript: "not-applicable",
  swift: "pending",
  rust: "verified",
};

const zeroEffects = {
  state: false,
  approval: false,
  key: false,
  ipcEndpoint: false,
  process: false,
  runtime: false,
  backend: false,
  guest: false,
};

const ordinaryRequirements = [
  { decision: "accept", variant: "ordinary" },
  { decision: "reject", variant: "malformed" },
];
const boundaryRequirements = [
  { decision: "accept", variant: "exact-maximum" },
  { decision: "reject", variant: "cap-plus-one" },
];

const holdObservations = [
  ["property-import-meta", [0, 0, 0, 0, 0]],
  ["method-import", [0, 0, 0, 0, 0]],
  ["template-interpolation-import", [0, 0, 1, 0, 0]],
  ["eval-string-data", [0, 0, 0, 0, 0]],
  ["division-regexp-counterexample", [0, 0, 1, 0, 0]],
  ["static-import", [1, 0, 0, 0, 0]],
  ["side-effect-import", [1, 0, 0, 0, 0]],
  ["export-from", [0, 1, 0, 0, 0]],
  ["export-star", [0, 1, 0, 0, 0]],
  ["import-meta", [0, 0, 0, 1, 0]],
  ["specifier-absolute", [0, 0, 1, 0, 0]],
  ["specifier-bare", [0, 0, 1, 0, 0]],
  ["specifier-node", [0, 0, 1, 0, 0]],
  ["specifier-npm", [0, 0, 1, 0, 0]],
  ["specifier-http", [0, 0, 1, 0, 0]],
  ["specifier-https", [0, 0, 1, 0, 0]],
  ["specifier-data", [0, 0, 1, 0, 0]],
  ["specifier-blob", [0, 0, 1, 0, 0]],
  ["specifier-file", [0, 0, 1, 0, 0]],
  ["specifier-capsule", [0, 0, 1, 0, 0]],
  ["commonjs-require", [0, 0, 0, 0, 1]],
  ["commonjs-require-resolve", [0, 0, 0, 0, 1]],
  ["commonjs-module-exports", [0, 0, 0, 0, 1]],
  ["commonjs-exports", [0, 0, 0, 0, 1]],
  ["commonjs-dirname", [0, 0, 0, 0, 1]],
  ["commonjs-filename", [0, 0, 0, 0, 1]],
  ["local-export", [0, 0, 0, 0, 0]],
  ["noncode-spellings", [0, 0, 0, 0, 0]],
];

export function addMjsSourceValidatorPassiveRulesAndCases({
  addCase,
  addRule,
  fixtureBytes,
  retainFixture,
}) {
  const candidate = encodeCandidate();
  const candidateDigest = domainHash(
    "capsule.source-validator.engineering-candidate/v0",
    candidate,
  );
  const artifactProfile = encodeArtifactProfile(candidateDigest);
  const artifactProfileDigest = domainHash(
    "capsule.source-validator.artifact-profile/v0",
    artifactProfile,
  );
  const candidateReference = retainFixture(
    "mjs-source-validator/engineering-candidate.bin",
    candidate,
  );
  const profileReference = retainFixture(
    "mjs-source-validator/artifact-profile.bin",
    artifactProfile,
  );

  addRule(
    "source-validator.identity-records",
    "ADR-0035#exact-parser-candidate",
    "The passive Oxc engineering-candidate and future enrolled-artifact profile have exact domain-separated fixed records.",
    ordinaryRequirements,
  );
  addPassiveCase({
    addCase,
    id: "source-validator.identity-records.engineering-candidate",
    description: "Accept the exact Oxc 0.140.0 engineering-candidate identity record.",
    ruleIds: ["source-validator.identity-records"],
    object: "SourceValidatorEngineeringCandidate",
    path: candidateReference.path,
    bytes: candidate,
  });
  addPassiveCase({
    addCase,
    id: "source-validator.identity-records.artifact-profile",
    description:
      "Accept the fixture-only enrolled-artifact profile shape without admitting its placeholder artifact.",
    ruleIds: ["source-validator.identity-records"],
    object: "SourceValidatorArtifactProfile",
    path: profileReference.path,
    bytes: artifactProfile,
    context: sourceValidatorContext({ engineeringCandidate: candidateReference }),
  });

  for (const [id, offset, classification] of [
    ["candidate-version", 12, "UNSUPPORTED"],
    ["candidate-lock-digest", 44, "UNSUPPORTED"],
    ["candidate-crate-checksum", 80, "UNSUPPORTED"],
  ]) {
    const mutated = Buffer.from(candidate);
    mutated[offset] ^= 0x01;
    addPassiveReject({
      addCase,
      id: `source-validator.identity-records.${id}`,
      description: `Reject the ${id.replaceAll("-", " ")} mutation.`,
      ruleIds: ["source-validator.identity-records"],
      object: "SourceValidatorEngineeringCandidate",
      path: `mjs-source-validator/reject-${id}.bin`,
      bytes: mutated,
      classification,
    });
  }
  addPassiveReject({
    addCase,
    id: "source-validator.identity-records.candidate-trailing",
    description: "Reject trailing bytes after the engineering-candidate record.",
    ruleIds: ["source-validator.identity-records"],
    object: "SourceValidatorEngineeringCandidate",
    path: "mjs-source-validator/reject-candidate-trailing.bin",
    bytes: Buffer.concat([candidate, Buffer.from([0])]),
    classification: "MALFORMED",
  });

  const profileMutations = [
    ["profile-version", mutateU16(artifactProfile, 12, 1), "UNSUPPORTED"],
    ["profile-media", mutateU16(artifactProfile, 22, 2), "UNSUPPORTED"],
    ["candidate-domain", mutateU16(artifactProfile, 24, sourceDigestDomain), "DOMAIN"],
    ["candidate-digest", flip(artifactProfile, 26), "BINDING"],
    ["executable-domain", mutateU16(artifactProfile, 58, sourceDigestDomain), "DOMAIN"],
    ["executable-zero", fill(artifactProfile, 60, 32, 0), "SCHEMA"],
    ["build-domain", mutateU16(artifactProfile, 92, sourceDigestDomain), "DOMAIN"],
    ["assessment-domain", mutateU16(artifactProfile, 126, sourceDigestDomain), "DOMAIN"],
  ];
  for (const [id, bytes, classification] of profileMutations) {
    addPassiveReject({
      addCase,
      id: `source-validator.identity-records.${id}`,
      description: `Reject the ${id.replaceAll("-", " ")} mutation.`,
      ruleIds: ["source-validator.identity-records"],
      object: "SourceValidatorArtifactProfile",
      path: `mjs-source-validator/reject-${id}.bin`,
      bytes,
      classification,
      context: sourceValidatorContext({ engineeringCandidate: candidateReference }),
    });
  }

  addRule(
    "source-validator.request-frame-boundary",
    "MJS_SOURCE_VALIDATOR_IMPLEMENTATION_PLAN#v0-contract-reconciliation",
    "A request is exactly 80 through 262224 bytes and carries one copied M1 source without clamping.",
    boundaryRequirements,
  );
  addRule(
    "source-validator.request-binding",
    "MJS_SOURCE_VALIDATOR_IMPLEMENTATION_PLAN#proposed-typed-operation",
    "Request protocol, method, profiles, media, correlation, length, digest, source profile, and one-frame EOF are closed and bound.",
    ordinaryRequirements,
  );
  addRule(
    "source-validator.m1-additive-mapping",
    "ADR-0035#status-and-evidence-limits",
    "Requests and passive result oracles consume the exact merged M1 byte, identity, HOLD, and SourceManifest fixtures additively.",
    ordinaryRequirements,
  );

  const ordinaryID = requestID(1);
  const ordinarySource = fixtureBytes("mjs-source/ordinary.mjs");
  const ordinaryRequest = encodeRequest(ordinaryID, ordinarySource);
  const ordinaryRequestReference = retainFixture(
    "mjs-source-validator/request-ordinary.bin",
    ordinaryRequest,
  );
  const ordinarySourceReference = retainFixture("mjs-source/ordinary.mjs", ordinarySource);
  const ordinaryManifestReference = retainFixture(
    "mjs-source-manifest/ordinary.cbor",
    fixtureBytes("mjs-source-manifest/ordinary.cbor"),
  );
  addPassiveCase({
    addCase,
    id: "source-validator.request-binding.ordinary",
    description: "Accept the exact ordinary one-shot request known answer.",
    ruleIds: ["source-validator.request-binding", "source-validator.m1-additive-mapping"],
    object: "SourceValidatorRequest",
    path: ordinaryRequestReference.path,
    bytes: ordinaryRequest,
    context: sourceValidatorContext({
      source: ordinarySourceReference,
      sourceManifest: ordinaryManifestReference,
    }),
  });

  const requestSources = [
    ["minimum", "mjs-source/minimum.mjs", "mjs-source-manifest/minimum.cbor", "minimum"],
    [
      "exact-maximum",
      "mjs-source/exact-maximum.mjs",
      "mjs-source-manifest/exact-maximum.cbor",
      "exact-maximum",
    ],
    ...[
      "lf",
      "crlf",
      "lone-cr",
      "line-separator",
      "paragraph-separator",
      "composed",
      "decomposed",
      "embedded-bom",
      "no-trailing-newline",
    ].map((id) => [id, `mjs-source/identity-${id}.mjs`, null, "ordinary"]),
  ];
  for (const [index, [id, sourcePath, manifestPath, variant]] of requestSources.entries()) {
    const source = fixtureBytes(sourcePath);
    const sourceReference = retainFixture(sourcePath, source);
    const manifestReference = manifestPath
      ? retainFixture(manifestPath, fixtureBytes(manifestPath))
      : null;
    const request = encodeRequest(requestID(10 + index), source);
    addPassiveCase({
      addCase,
      id: `source-validator.m1-additive-mapping.${id}`,
      description: `Frame the exact merged M1 ${id.replaceAll("-", " ")} source without rewriting.`,
      ruleIds: ["source-validator.m1-additive-mapping", "source-validator.request-frame-boundary"],
      object: "SourceValidatorRequest",
      variant,
      path: `mjs-source-validator/request-${id}.bin`,
      bytes: request,
      context: sourceValidatorContext({
        source: sourceReference,
        sourceManifest: manifestReference,
      }),
    });
  }

  for (const [index, [id, counts]] of holdObservations.entries()) {
    const sourcePath = `mjs-source/language-hold-${id}.mjs`;
    const source = fixtureBytes(sourcePath);
    const sourceReference = retainFixture(sourcePath, source);
    const idBytes = requestID(40 + index);
    const request = encodeRequest(idBytes, source);
    const requestReference = retainFixture(`mjs-source-validator/request-hold-${id}.bin`, request);
    addPassiveCase({
      addCase,
      id: `source-validator.m1-additive-mapping.request-${id}`,
      description: `Frame the exact merged M1 HOLD ${id.replaceAll("-", " ")} source bytes.`,
      ruleIds: ["source-validator.m1-additive-mapping"],
      object: "SourceValidatorRequest",
      path: requestReference.path,
      bytes: request,
      context: sourceValidatorContext({ source: sourceReference }),
    });
    const result = encodeResult({
      id: idBytes,
      source,
      artifactProfileDigest,
      parse: 1,
      policy: counts.some((count) => count > 0) ? 2 : 1,
      classification: resultClassification(counts),
      counts,
    });
    addPassiveCase({
      addCase,
      id: `source-validator.m1-additive-mapping.result-${id}`,
      description: `Retain the fixed passive Oxc observation for the exact M1 HOLD ${id.replaceAll("-", " ")} bytes.`,
      ruleIds: ["source-validator.m1-additive-mapping"],
      object: "SourceValidatorResult",
      path: `mjs-source-validator/result-hold-${id}.bin`,
      bytes: result,
      context: sourceValidatorContext({
        source: sourceReference,
        request: requestReference,
        artifactProfile: profileReference,
        engineeringCandidate: candidateReference,
      }),
    });
  }

  for (const [id, sourcePath] of [
    ["cap-plus-one", "mjs-source/cap-plus-one.mjs"],
    ["invalid-utf8", "mjs-source/invalid-utf8.mjs"],
    ["unpaired-surrogate-utf8", "mjs-source/unpaired-surrogate-utf8.mjs"],
    ["leading-bom", "mjs-source/leading-bom.mjs"],
  ]) {
    const source = fixtureBytes(sourcePath);
    const request = encodeRequest(requestID(90), source);
    addPassiveReject({
      addCase,
      id: `source-validator.request-binding.${id}`,
      description: `Reject the exact merged M1 ${id.replaceAll("-", " ")} source before parsing.`,
      ruleIds: [
        id === "cap-plus-one"
          ? "source-validator.request-frame-boundary"
          : id === "invalid-utf8"
            ? "source-validator.m1-additive-mapping"
            : "source-validator.request-binding",
      ],
      object: "SourceValidatorRequest",
      variant: id === "cap-plus-one" ? "cap-plus-one" : "malformed",
      path: `mjs-source-validator/request-${id}.bin`,
      bytes: request,
      classification: "DOMAIN",
      context: sourceValidatorContext({ source: retainFixture(sourcePath, source) }),
    });
  }

  const requestMutations = [
    ["protocol-version", mutateU16(ordinaryRequest, 12, 1), "UNSUPPORTED"],
    ["frame-kind", mutateU16(ordinaryRequest, 14, resultKind), "UNSUPPORTED"],
    ["method", mutateU16(ordinaryRequest, 16, 2), "UNSUPPORTED"],
    ["validator-profile", mutateU16(ordinaryRequest, 18, 2), "UNSUPPORTED"],
    ["source-profile", mutateU16(ordinaryRequest, 20, 2), "UNSUPPORTED"],
    ["media", mutateU16(ordinaryRequest, 22, 2), "UNSUPPORTED"],
    ["correlation-domain", mutateU16(ordinaryRequest, 24, 2), "DOMAIN"],
    ["correlation-zero", fill(ordinaryRequest, 28, 16, 0), "SCHEMA"],
    ["digest-domain", mutateU16(ordinaryRequest, 26, candidateDigestDomain), "DOMAIN"],
    ["length", mutateU32(ordinaryRequest, 44, ordinarySource.length + 1), "BINDING"],
    ["digest", flip(ordinaryRequest, 48), "BINDING"],
    ["source-mutation", flip(ordinaryRequest, requestHeaderBytes), "BINDING"],
    ["truncated", ordinaryRequest.subarray(0, ordinaryRequest.length - 1), "MALFORMED"],
    ["trailing", Buffer.concat([ordinaryRequest, Buffer.from([0])]), "MALFORMED"],
    ["duplicate-frame", Buffer.concat([ordinaryRequest, ordinaryRequest]), "MALFORMED"],
  ];
  for (const [id, bytes, classification] of requestMutations) {
    addPassiveReject({
      addCase,
      id: `source-validator.request-binding.${id}`,
      description: `Reject the request ${id.replaceAll("-", " ")} mutation.`,
      ruleIds: ["source-validator.request-binding"],
      object: "SourceValidatorRequest",
      path: `mjs-source-validator/reject-request-${id}.bin`,
      bytes,
      classification,
      context: sourceValidatorContext({ source: ordinarySourceReference }),
    });
  }

  addRule(
    "source-validator.result-binding",
    "MJS_SOURCE_VALIDATOR_IMPLEMENTATION_PLAN#proposed-typed-operation",
    "The fixed 138-byte result binds request, source, artifact profile, closed status/classification, four syntax counts, and one free-CommonJS count.",
    ordinaryRequirements,
  );
  const ordinaryResult = encodeResult({
    id: ordinaryID,
    source: ordinarySource,
    artifactProfileDigest,
    parse: 1,
    policy: 1,
    classification: 0,
    counts: [0, 0, 0, 0, 0],
  });
  addPassiveCase({
    addCase,
    id: "source-validator.result-binding.ordinary",
    description: "Accept the fixed ordinary request/result/artifact-profile known answer.",
    ruleIds: ["source-validator.result-binding"],
    object: "SourceValidatorResult",
    path: "mjs-source-validator/result-ordinary.bin",
    bytes: ordinaryResult,
    context: sourceValidatorContext({
      source: ordinarySourceReference,
      request: ordinaryRequestReference,
      artifactProfile: profileReference,
      engineeringCandidate: candidateReference,
    }),
  });

  for (const [id, parse, classification] of [
    ["parser-diagnostic", 2, 4],
    ["semantic-diagnostic", 3, 5],
  ]) {
    addPassiveCase({
      addCase,
      id: `source-validator.result-binding.${id}`,
      description: `Accept the closed synthetic ${id.replaceAll("-", " ")} status branch as a passive protocol oracle.`,
      ruleIds: ["source-validator.result-binding"],
      object: "SourceValidatorResult",
      path: `mjs-source-validator/result-${id}.bin`,
      bytes: encodeResult({
        id: ordinaryID,
        source: ordinarySource,
        artifactProfileDigest,
        parse,
        policy: 3,
        classification,
        counts: [0, 0, 0, 0, 0],
      }),
      context: sourceValidatorContext({
        source: ordinarySourceReference,
        request: ordinaryRequestReference,
        artifactProfile: profileReference,
        engineeringCandidate: candidateReference,
      }),
    });
  }

  const resultMutations = [
    ["protocol-version", mutateU16(ordinaryResult, 12, 1), "UNSUPPORTED"],
    ["method", mutateU16(ordinaryResult, 16, 2), "UNSUPPORTED"],
    ["media", mutateU16(ordinaryResult, 22, 2), "UNSUPPORTED"],
    ["correlation-domain", mutateU16(ordinaryResult, 24, 2), "DOMAIN"],
    ["correlation", flip(ordinaryResult, 28), "BINDING"],
    ["digest-domain", mutateU16(ordinaryResult, 26, candidateDigestDomain), "DOMAIN"],
    ["source-length", mutateU32(ordinaryResult, 44, ordinarySource.length + 1), "BINDING"],
    ["source-digest", flip(ordinaryResult, 48), "BINDING"],
    ["artifact-domain", mutateU16(ordinaryResult, 80, sourceDigestDomain), "DOMAIN"],
    ["artifact-digest", flip(ordinaryResult, 82), "BINDING"],
    ["parse-status", fill(ordinaryResult, 114, 1, 0xff), "DOMAIN"],
    ["policy-status", fill(ordinaryResult, 115, 1, 2), "DOMAIN"],
    ["classification", fill(ordinaryResult, 116, 1, 1), "DOMAIN"],
    ["reserved", fill(ordinaryResult, 117, 1, 1), "UNSUPPORTED"],
    ["count-on-allow", mutateU32(ordinaryResult, 118, 1), "DOMAIN"],
    ["static-import-count-cap", mutateU32(ordinaryResult, 118, sourceMaximumBytes + 1), "DOMAIN"],
    ["export-from-count-cap", mutateU32(ordinaryResult, 122, sourceMaximumBytes + 1), "DOMAIN"],
    [
      "import-expression-count-cap",
      mutateU32(ordinaryResult, 126, sourceMaximumBytes + 1),
      "DOMAIN",
    ],
    ["import-meta-count-cap", mutateU32(ordinaryResult, 130, sourceMaximumBytes + 1), "DOMAIN"],
    ["commonjs-count-cap", mutateU32(ordinaryResult, 134, sourceMaximumBytes + 1), "DOMAIN"],
    ["truncated", ordinaryResult.subarray(0, ordinaryResult.length - 1), "MALFORMED"],
    ["trailing", Buffer.concat([ordinaryResult, Buffer.from([0])]), "MALFORMED"],
    ["duplicate-frame", Buffer.concat([ordinaryResult, ordinaryResult]), "MALFORMED"],
  ];
  for (const [id, bytes, classification] of resultMutations) {
    addPassiveReject({
      addCase,
      id: `source-validator.result-binding.${id}`,
      description: `Reject the result ${id.replaceAll("-", " ")} mutation.`,
      ruleIds: ["source-validator.result-binding"],
      object: "SourceValidatorResult",
      path: `mjs-source-validator/reject-result-${id}.bin`,
      bytes,
      classification,
      context: sourceValidatorContext({
        source: ordinarySourceReference,
        request: ordinaryRequestReference,
        artifactProfile: profileReference,
        engineeringCandidate: candidateReference,
      }),
    });
  }

  const replayRequest = encodeRequest(requestID(2), ordinarySource);
  const replayReference = retainFixture(
    "mjs-source-validator/request-cross-request.bin",
    replayRequest,
  );
  addPassiveReject({
    addCase,
    id: "source-validator.result-binding.cross-request-substitution",
    description: "Reject an otherwise valid result replayed onto a distinct current request ID.",
    ruleIds: ["source-validator.result-binding"],
    object: "SourceValidatorResult",
    path: "mjs-source-validator/result-ordinary.bin",
    bytes: ordinaryResult,
    classification: "BINDING",
    context: sourceValidatorContext({
      source: ordinarySourceReference,
      request: replayReference,
      artifactProfile: profileReference,
      engineeringCandidate: candidateReference,
    }),
  });
}

function addPassiveCase({
  addCase,
  id,
  description,
  ruleIds,
  object,
  variant = "ordinary",
  path,
  bytes,
  context = { kind: "none" },
}) {
  addCase({
    id,
    description,
    ruleIds,
    object,
    wireFormat: "fixed-frame",
    mediaType,
    variant,
    path,
    bytes,
    context,
    owner: "source-validator-passive-contract",
    implementations,
  });
}

function addPassiveReject({
  addCase,
  id,
  description,
  ruleIds,
  object,
  variant = "malformed",
  path,
  bytes,
  classification,
  context = { kind: "none" },
}) {
  addCase({
    id,
    description,
    ruleIds,
    object,
    wireFormat: "fixed-frame",
    mediaType,
    variant,
    path,
    bytes,
    context,
    decision: "reject",
    classification,
    owner: "source-validator-passive-contract",
    implementations,
    effects: zeroEffects,
  });
}

function sourceValidatorContext({
  source = null,
  request = null,
  artifactProfile = null,
  engineeringCandidate = null,
  sourceManifest = null,
}) {
  return {
    kind: "source-validator",
    source,
    request,
    artifactProfile,
    engineeringCandidate,
    sourceManifest,
  };
}

function encodeRequest(id, source) {
  const frame = Buffer.alloc(requestHeaderBytes + source.length);
  frame.writeUInt32BE(frame.length - 4, 0);
  frame.write("CAPMJSRQ", 4, "ascii");
  writeCommon(frame, requestKind, id, source);
  source.copy(frame, requestHeaderBytes);
  return frame;
}

function encodeResult({
  id,
  source,
  artifactProfileDigest,
  parse,
  policy,
  classification,
  counts,
}) {
  const frame = Buffer.alloc(resultFrameBytes);
  frame.writeUInt32BE(frame.length - 4, 0);
  frame.write("CAPMJSRS", 4, "ascii");
  writeCommon(frame, resultKind, id, source);
  frame.writeUInt16BE(artifactProfileDomain, 80);
  artifactProfileDigest.copy(frame, 82);
  frame[114] = parse;
  frame[115] = policy;
  frame[116] = classification;
  counts.forEach((count, index) => {
    frame.writeUInt32BE(count, 118 + index * 4);
  });
  return frame;
}

function writeCommon(frame, kind, id, source) {
  [
    protocolVersion,
    kind,
    method,
    validatorProfile,
    sourceProfile,
    sourceMedia,
    correlationDomain,
    sourceDigestDomain,
  ].forEach((value, index) => {
    frame.writeUInt16BE(value, 12 + index * 2);
  });
  id.copy(frame, 28);
  frame.writeUInt32BE(source.length, 44);
  sha256(source).copy(frame, 48);
}

function encodeCandidate() {
  const frame = Buffer.alloc(candidateFrameBytes);
  frame.writeUInt32BE(frame.length - 4, 0);
  frame.write("CAPMJSCI", 4, "ascii");
  [0, 1, 0, 140, 0, 1, 95, 0, 0x000f, 6, 65, 0].forEach((value, index) => {
    frame.writeUInt16BE(value, 12 + index * 2);
  });
  frame.writeUInt32BE(24_449_903, 36);
  frame.writeUInt32BE(1_854_528, 40);
  hex("505669a07338603876bc96c242f8d5af386d3a13139e70110a8b52f39bae69ac").copy(frame, 44);
  const crates = [
    "0f8245ba555b465d3577732d5f9d9306babb0aaa7b80e97a2ce21f74fae442a3",
    "3305400b90fff2a30b272b58fe6080d25369407b2ac37c4ac652996a9677efe0",
    "4640e6d0de2e0f6c820d1444a468d070c710111df76ce90a1694ac386641e133",
    "8abd68f81349d37ea79f1d99d2370e15f282cc9fbe66e8544d072595744ab38e",
    "5967f96881e1694d10b453311fa681b4df0f38760628e1de613b046566cd8c8e",
    "e83fa0a0fe6e5e2f5abb173a64afe8db711bb612acbc002c663fd13b08a8cbf3",
  ];
  crates.forEach((checksum, index) => {
    const offset = 76 + index * 36;
    frame.writeUInt16BE(index + 1, offset);
    hex(checksum).copy(frame, offset + 4);
  });
  return frame;
}

function encodeArtifactProfile(candidateDigest) {
  const frame = Buffer.alloc(artifactProfileFrameBytes);
  frame.writeUInt32BE(frame.length - 4, 0);
  frame.write("CAPMJSAP", 4, "ascii");
  [0, protocolVersion, method, validatorProfile, sourceProfile, sourceMedia].forEach(
    (value, index) => {
      frame.writeUInt16BE(value, 12 + index * 2);
    },
  );
  frame.writeUInt16BE(candidateDigestDomain, 24);
  candidateDigest.copy(frame, 26);
  frame.writeUInt16BE(executableDigestDomain, 58);
  Buffer.alloc(32, 0x11).copy(frame, 60);
  frame.writeUInt16BE(buildDigestDomain, 92);
  Buffer.alloc(32, 0x22).copy(frame, 94);
  frame.writeUInt16BE(assessmentDigestDomain, 126);
  Buffer.alloc(32, 0x33).copy(frame, 128);
  return frame;
}

function resultClassification(counts) {
  const syntax = counts.slice(0, 4).some((count) => count > 0);
  const commonjs = counts[4] > 0;
  if (syntax && commonjs) return 3;
  if (syntax) return 1;
  if (commonjs) return 2;
  return 0;
}

function requestID(seed) {
  return Buffer.from(Array.from({ length: 16 }, (_, index) => (seed + index) & 0xff));
}

function mutateU16(value, offset, replacement) {
  const result = Buffer.from(value);
  result.writeUInt16BE(replacement, offset);
  return result;
}

function mutateU32(value, offset, replacement) {
  const result = Buffer.from(value);
  result.writeUInt32BE(replacement, offset);
  return result;
}

function flip(value, offset) {
  const result = Buffer.from(value);
  result[offset] ^= 0x01;
  return result;
}

function fill(value, offset, length, replacement) {
  const result = Buffer.from(value);
  result.fill(replacement, offset, offset + length);
  return result;
}

function domainHash(domain, value) {
  return createHash("sha256")
    .update(domain, "ascii")
    .update(Buffer.from([0]))
    .update(value)
    .digest();
}

function sha256(value) {
  return createHash("sha256").update(value).digest();
}

function hex(value) {
  return Buffer.from(value, "hex");
}
