import { sha256 } from "./lib/fixture-bytes.mjs";

const sourceMaximumBytes = 262_144;
const requestHeaderBytes = 216;
const requestMaximumBytes = requestHeaderBytes + sourceMaximumBytes;
const resultFrameBytes = 248;
const policyFrameBytes = 256;
const profileFrameBytes = 256;
const consumerFrameBytes = 192;
const mediaType = "application/capsule.source-validator-frame;v=1";

const implementations = {
  "fixture-integrity": "verified",
  go: "verified",
  typescript: "verified",
  swift: "pending",
  rust: "pending",
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

const roles = [
  { name: "daemon", tag: 1 },
  { name: "approval-broker", tag: 2 },
];

export function addMjsSourceValidatorV1RulesAndCases({
  addCase,
  addRule,
  fixtureBytes,
  retainFixture,
}) {
  addRule(
    "source-validator-v1.role-families",
    "ADR-0036#two-role-specific-private-launcher-services",
    "The daemon and Approval Broker use distinct v1 request, result, process, artifact, and consumer families; cross-role substitution refuses.",
    ordinaryRequirements,
  );
  addRule(
    "source-validator-v1.request-boundary",
    "MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1#role-separated-request-ownership",
    `Each role-specific request is exactly ${requestHeaderBytes} through ${requestMaximumBytes} bytes and owns one copied M1 source.`,
    boundaryRequirements,
  );
  addRule(
    "source-validator-v1.result-binding",
    "MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1#role-separated-result-ownership",
    `Each fixed ${resultFrameBytes}-byte result binds its exact current role request and complete cleanup disposition.`,
    ordinaryRequirements,
  );
  addRule(
    "source-validator-v1.inactive-resource-policy",
    "MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1#resource-policy-boundary",
    "R1 freezes the inactive resource-policy shape and structural caps while refusing every invented active measurement.",
    ordinaryRequirements,
  );
  addRule(
    "source-validator-v1.profile-bindings",
    "MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1#required-passive-r1-targets",
    "Fixture-only process, artifact, and consumer projections bind one exact role family without enrolling an artifact or consumer.",
    ordinaryRequirements,
  );
  addRule(
    "source-validator-v1.cross-version-refusal",
    "ADR-0036#new-v1-protocol-and-profile-families",
    "Historical v0 request/result bytes are never relabeled or accepted by the v1 family.",
    ordinaryRequirements,
  );

  const ordinarySource = fixtureBytes("mjs-source/ordinary.mjs");
  const maximumSource = fixtureBytes("mjs-source/exact-maximum.mjs");
  const v0Request = fixtureBytes("mjs-source-validator/request-ordinary.bin");

  for (const role of roles) {
    const prefix = `mjs-source-validator-v1/${role.name}`;
    const policy = encodeInactivePolicy(role);
    const policyDigest = sha256(policy);
    const process = encodeProcessProfile(role, policyDigest);
    const processDigest = sha256(process);
    const artifact = encodeArtifactProfile(role, processDigest, policyDigest);
    const artifactDigest = sha256(artifact);
    const consumer = encodeConsumerProjection(role, processDigest, artifactDigest, policyDigest);
    const ordinaryRequest = encodeRequest(role, ordinarySource, artifactDigest, policyDigest);
    const maximumRequest = encodeRequest(role, maximumSource, artifactDigest, policyDigest);
    const ordinaryResult = encodeResult(role, ordinarySource, artifactDigest, policyDigest);

    const policyReference = retainFixture(`${prefix}/resource-policy-inactive.bin`, policy);
    const processReference = retainFixture(`${prefix}/process-profile.bin`, process);
    const artifactReference = retainFixture(`${prefix}/artifact-profile.bin`, artifact);
    const consumerReference = retainFixture(`${prefix}/consumer-projection.bin`, consumer);
    const ordinaryRequestReference = retainFixture(
      `${prefix}/request-ordinary.bin`,
      ordinaryRequest,
    );

    for (const [object, suffix, bytes, context] of [
      ["SourceValidatorV1ResourcePolicy", "resource-policy-inactive", policy, noContext()],
      ["SourceValidatorV1ProcessProfile", "process-profile", process, noContext()],
      ["SourceValidatorV1ArtifactProfile", "artifact-profile", artifact, noContext()],
      ["SourceValidatorV1ConsumerProjection", "consumer-projection", consumer, noContext()],
      [
        "SourceValidatorV1Request",
        "request-ordinary",
        ordinaryRequest,
        sourceContext(retainFixture("mjs-source/ordinary.mjs", ordinarySource)),
      ],
      [
        "SourceValidatorV1Result",
        "result-ordinary",
        ordinaryResult,
        resultContext(ordinaryRequestReference, artifactReference),
      ],
    ]) {
      addV1Case({
        addCase,
        id: `source-validator-v1.${role.name}.${suffix}`,
        description: `Accept the fixture-only ${role.name} ${suffix.replaceAll(
          "-",
          " ",
        )} known answer.`,
        ruleIds: [
          suffix === "resource-policy-inactive"
            ? "source-validator-v1.inactive-resource-policy"
            : suffix === "result-ordinary"
              ? "source-validator-v1.result-binding"
              : suffix === "request-ordinary"
                ? "source-validator-v1.request-boundary"
                : "source-validator-v1.profile-bindings",
          "source-validator-v1.role-families",
          ...(suffix === "request-ordinary" ? ["source-validator-v1.cross-version-refusal"] : []),
        ],
        object,
        path: `${prefix}/${suffix}.bin`,
        bytes,
        context,
      });
    }

    addV1Case({
      addCase,
      id: `source-validator-v1.${role.name}.request-exact-maximum`,
      description: `Accept the exact ${role.name} request maximum without clamping.`,
      ruleIds: ["source-validator-v1.request-boundary"],
      object: "SourceValidatorV1Request",
      variant: "exact-maximum",
      path: `${prefix}/request-exact-maximum.bin`,
      bytes: maximumRequest,
      context: sourceContext(retainFixture("mjs-source/exact-maximum.mjs", maximumSource)),
    });

    const capPlusOne = Buffer.concat([maximumRequest, Buffer.from([0])]);
    addV1Reject({
      addCase,
      id: `source-validator-v1.${role.name}.request-cap-plus-one`,
      description: `Refuse the ${role.name} request cap plus one before source allocation growth.`,
      ruleIds: ["source-validator-v1.request-boundary"],
      object: "SourceValidatorV1Request",
      variant: "cap-plus-one",
      path: `${prefix}/reject-request-cap-plus-one.bin`,
      bytes: capPlusOne,
      classification: "MALFORMED",
    });

    for (const [suffix, object, original, offset, classification] of [
      ["request-cross-role", "SourceValidatorV1Request", ordinaryRequest, 16, "DOMAIN"],
      ["result-cross-role", "SourceValidatorV1Result", ordinaryResult, 16, "DOMAIN"],
      ["process-cross-role", "SourceValidatorV1ProcessProfile", process, 16, "DOMAIN"],
      ["artifact-cross-role", "SourceValidatorV1ArtifactProfile", artifact, 16, "DOMAIN"],
      ["consumer-cross-role", "SourceValidatorV1ConsumerProjection", consumer, 16, "DOMAIN"],
    ]) {
      addV1Reject({
        addCase,
        id: `source-validator-v1.${role.name}.${suffix}`,
        description: `Refuse a ${role.name} ${suffix.replaceAll("-", " ")} substitution.`,
        ruleIds: [
          object === "SourceValidatorV1ProcessProfile" ||
          object === "SourceValidatorV1ArtifactProfile" ||
          object === "SourceValidatorV1ConsumerProjection"
            ? "source-validator-v1.profile-bindings"
            : "source-validator-v1.role-families",
        ],
        object,
        path: `${prefix}/reject-${suffix}.bin`,
        bytes: mutateU16(original, offset, role.tag === 1 ? 2 : 1),
        classification,
        context:
          object === "SourceValidatorV1Result"
            ? resultContext(ordinaryRequestReference, artifactReference)
            : noContext(),
      });
    }

    addV1Reject({
      addCase,
      id: `source-validator-v1.${role.name}.request-unknown-version`,
      description: `Refuse an unknown ${role.name} request protocol version.`,
      ruleIds: ["source-validator-v1.cross-version-refusal"],
      object: "SourceValidatorV1Request",
      path: `${prefix}/reject-request-unknown-version.bin`,
      bytes: mutateU16(ordinaryRequest, 12, 2),
      classification: "UNSUPPORTED",
    });
    for (const [suffix, offset, width] of [
      ["correlation-mismatch", 44, 1],
      ["installation-mismatch", 60, 1],
      ["epoch-sequence-mismatch", 76, 8],
      ["epoch-digest-mismatch", 84, 1],
      ["artifact-profile-mismatch", 152, 1],
      ["resource-policy-mismatch", 184, 1],
    ]) {
      addV1Reject({
        addCase,
        id: `source-validator-v1.${role.name}.request-${suffix}`,
        description: `Refuse the ${role.name} request ${suffix.replaceAll(
          "-",
          " ",
        )} against the trusted current context.`,
        ruleIds: ["source-validator-v1.role-families"],
        object: "SourceValidatorV1Request",
        path: `${prefix}/reject-request-${suffix}.bin`,
        bytes: flipAt(ordinaryRequest, offset, width),
        classification: "BINDING",
      });
    }
    addV1Reject({
      addCase,
      id: `source-validator-v1.${role.name}.result-trailing`,
      description: `Refuse trailing bytes after the fixed ${role.name} result.`,
      ruleIds: ["source-validator-v1.result-binding"],
      object: "SourceValidatorV1Result",
      path: `${prefix}/reject-result-trailing.bin`,
      bytes: Buffer.concat([ordinaryResult, Buffer.from([0])]),
      classification: "MALFORMED",
      context: resultContext(ordinaryRequestReference, artifactReference),
    });
    addV1Reject({
      addCase,
      id: `source-validator-v1.${role.name}.policy-invented-threshold`,
      description: `Refuse an invented ${role.name} physical-footprint threshold in inactive R1 policy.`,
      ruleIds: ["source-validator-v1.inactive-resource-policy"],
      object: "SourceValidatorV1ResourcePolicy",
      path: `${prefix}/reject-policy-invented-threshold.bin`,
      bytes: mutateU64(policy, 178, 1n),
      classification: "SCHEMA",
    });

    // Keep retained references live and therefore protected by corpus file-set verification.
    void policyReference;
    void processReference;
    void consumerReference;
  }

  for (const role of roles) {
    addV1Reject({
      addCase,
      id: `source-validator-v1.${role.name}.v0-request-as-v1`,
      description: `Refuse the historical v0 request as a ${role.name} v1 request.`,
      ruleIds: ["source-validator-v1.cross-version-refusal"],
      object: "SourceValidatorV1Request",
      path: `mjs-source-validator-v1/${role.name}/reject-v0-request-as-v1.bin`,
      bytes: v0Request,
      classification: "UNSUPPORTED",
    });
  }
}

function encodeRequest(role, source, artifactDigest, policyDigest) {
  const frame = Buffer.alloc(requestHeaderBytes + source.length);
  writeHeader(frame, "CSV1REQ0", 1, role);
  writeRoleTags(frame, role, 11);
  requestID(role.tag).copy(frame, 44);
  installationID().copy(frame, 60);
  frame.writeBigUInt64BE(7n, 76);
  digestByte(0x22).copy(frame, 84);
  frame.writeUInt32BE(source.length, 116);
  sha256(source).copy(frame, 120);
  artifactDigest.copy(frame, 152);
  policyDigest.copy(frame, 184);
  source.copy(frame, requestHeaderBytes);
  return frame;
}

function encodeResult(role, source, artifactDigest, policyDigest) {
  const frame = Buffer.alloc(resultFrameBytes);
  writeHeader(frame, "CSV1RES0", 2, role);
  writeRoleTags(frame, role, 11);
  requestID(role.tag).copy(frame, 44);
  installationID().copy(frame, 60);
  frame.writeBigUInt64BE(7n, 76);
  digestByte(0x22).copy(frame, 84);
  frame.writeUInt32BE(source.length, 116);
  sha256(source).copy(frame, 120);
  artifactDigest.copy(frame, 152);
  policyDigest.copy(frame, 184);
  frame[216] = 1;
  frame[217] = 1;
  frame[218] = 0;
  frame[219] = 1;
  return frame;
}

function encodeInactivePolicy(role) {
  const frame = Buffer.alloc(policyFrameBytes);
  writeHeader(frame, "CSV1POL0", 3, role);
  frame.writeUInt16BE(0, 18);
  frame.writeUInt16BE(1, 20);
  frame.writeUInt16BE(1, 128);
  frame.writeUInt16BE(0, 130);
  frame.writeUInt16BE(1, 132);
  frame.writeUInt16BE(2, 134);
  frame.writeUInt32BE(requestMaximumBytes, 136);
  frame.writeUInt32BE(resultFrameBytes, 140);
  frame.writeUInt32BE(0, 144);
  frame.writeUInt16BE(1, 148);
  frame.writeUInt16BE(0, 150);
  frame.writeUInt32BE(requestMaximumBytes + resultFrameBytes, 152);
  frame.writeUInt16BE(1, 168);
  return frame;
}

function encodeProcessProfile(role, policyDigest) {
  const frame = Buffer.alloc(profileFrameBytes);
  writeHeader(frame, "CSV1PRC0", 4, role);
  writeRoleTags(frame, role, 6);
  const digests = [
    policyDigest,
    ...Array.from({ length: 6 }, (_, index) => digestByte(0x60 + index + role.tag)),
  ];
  digests.forEach((digest, index) => {
    digest.copy(frame, 32 + index * 32);
  });
  return frame;
}

function encodeArtifactProfile(role, processDigest, policyDigest) {
  const frame = Buffer.alloc(profileFrameBytes);
  writeHeader(frame, "CSV1ART0", 5, role);
  writeRoleTags(frame, role, 6);
  const digests = [
    processDigest,
    policyDigest,
    ...Array.from({ length: 5 }, (_, index) => digestByte(0x70 + index + role.tag)),
  ];
  digests.forEach((digest, index) => {
    digest.copy(frame, 32 + index * 32);
  });
  return frame;
}

function encodeConsumerProjection(role, processDigest, artifactDigest, policyDigest) {
  const frame = Buffer.alloc(consumerFrameBytes);
  writeHeader(frame, "CSV1CON0", 6, role);
  writeRoleTags(frame, role, 6);
  installationID().copy(frame, 32);
  frame.writeBigUInt64BE(7n, 48);
  digestByte(0x22).copy(frame, 56);
  artifactDigest.copy(frame, 88);
  policyDigest.copy(frame, 120);
  const fixtureProcessDigest = Buffer.alloc(32);
  fixtureProcessDigest[0] = 0x80 + role.tag;
  fixtureProcessDigest.copy(frame, 152);
  void processDigest;
  return frame;
}

function writeHeader(frame, magic, kind, role) {
  frame.writeUInt32BE(frame.length - 4, 0);
  frame.write(magic, 4, "ascii");
  frame.writeUInt16BE(1, 12);
  frame.writeUInt16BE(kind, 14);
  frame.writeUInt16BE(role.tag, 16);
}

function writeRoleTags(frame, role, count) {
  for (let index = 0; index < count; index += 1) {
    frame.writeUInt16BE(role.tag * 0x100 + index + 1, 18 + index * 2);
  }
}

function requestID(seed) {
  const value = Buffer.alloc(16);
  value[0] = 0x50 + seed;
  return value;
}

function installationID() {
  const value = Buffer.alloc(16);
  value[0] = 0x11;
  return value;
}

function digestByte(value) {
  return Buffer.alloc(32, value);
}

function mutateU16(value, offset, replacement) {
  const result = Buffer.from(value);
  result.writeUInt16BE(replacement, offset);
  return result;
}

function mutateU64(value, offset, replacement) {
  const result = Buffer.from(value);
  result.writeBigUInt64BE(replacement, offset);
  return result;
}

function flipAt(value, offset, width) {
  const result = Buffer.from(value);
  result[offset + width - 1] ^= 0x01;
  return result;
}

function addV1Case({
  addCase,
  id,
  description,
  ruleIds,
  object,
  variant = "ordinary",
  path,
  bytes,
  context = noContext(),
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
    owner: "source-validator-passive-v1-contract",
    implementations,
  });
}

function addV1Reject(options) {
  const { addCase, classification, ...rest } = options;
  addCase({
    ...rest,
    wireFormat: "fixed-frame",
    mediaType,
    variant: rest.variant ?? "malformed",
    context: rest.context ?? noContext(),
    decision: "reject",
    classification,
    owner: "source-validator-passive-v1-contract",
    implementations,
    effects: zeroEffects,
  });
}

function noContext() {
  return { kind: "none" };
}

function sourceContext(source) {
  return {
    kind: "source-validator",
    source,
    request: null,
    artifactProfile: null,
    engineeringCandidate: null,
    sourceManifest: null,
  };
}

function resultContext(request, artifactProfile) {
  return {
    kind: "source-validator",
    source: null,
    request,
    artifactProfile,
    engineeringCandidate: null,
    sourceManifest: null,
  };
}
