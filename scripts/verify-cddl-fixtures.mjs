import { readFile } from "node:fs/promises";

const fixtureUrl = new URL("../schemas/fixtures/approval-grant-v0.json", import.meta.url);
const fixture = JSON.parse(await readFile(fixtureUrl, "utf8"));

assertExactKeys(fixture, [
  "attemptNonceHex",
  "audience",
  "epochDigestHex",
  "expiresAt",
  "installationIdHex",
  "issuedAt",
  "objectType",
  "objectVersion",
  "payloadHex",
  "planDigestHex",
  "protected",
  "protectedHex",
  "purpose",
  "registrationIdHex",
  "supervisorIdHex",
]);
assertExactKeys(fixture.protected, ["algorithm", "contentType", "keyIdHex"]);

assertEqual(fixture.objectType, "capsule.approval-grant", "object type");
assertEqual(fixture.objectVersion, 0, "object version");
assertEqual(fixture.purpose, "capsule.plan.approve", "purpose");
assertEqual(fixture.audience, "capsule.execution-supervisor", "audience");
assertSafeUnsigned(fixture.issuedAt, "issuedAt");
assertSafeUnsigned(fixture.expiresAt, "expiresAt");
if (fixture.expiresAt <= fixture.issuedAt) {
  throw new Error("ApprovalGrant fixture expiry must follow issue time");
}

const payload = new Map([
  [1, fixture.objectType],
  [2, fixture.objectVersion],
  [3, decodeFixedHex(fixture.installationIdHex, 16, "installationIdHex")],
  [4, decodeFixedHex(fixture.epochDigestHex, 32, "epochDigestHex")],
  [5, decodeFixedHex(fixture.registrationIdHex, 16, "registrationIdHex")],
  [6, decodeFixedHex(fixture.planDigestHex, 32, "planDigestHex")],
  [7, decodeFixedHex(fixture.supervisorIdHex, 16, "supervisorIdHex")],
  [8, decodeFixedHex(fixture.attemptNonceHex, 16, "attemptNonceHex")],
  [9, fixture.purpose],
  [10, fixture.audience],
  [11, fixture.issuedAt],
  [12, fixture.expiresAt],
]);

assertEqual(fixture.protected.algorithm, -7, "protected algorithm");
assertEqual(
  fixture.protected.contentType,
  "application/capsule.approval-grant+cbor;v=0",
  "protected content type",
);
const keyId = decodeHex(fixture.protected.keyIdHex, "protected.keyIdHex");
if (keyId.length < 1 || keyId.length > 64) {
  throw new Error("protected.keyIdHex must encode 1 through 64 bytes");
}
const protectedHeader = new Map([
  [1, fixture.protected.algorithm],
  [3, fixture.protected.contentType],
  [4, keyId],
]);

assertHexEqual(encodeDeterministicCbor(payload), fixture.payloadHex, "ApprovalGrant payload");
assertHexEqual(
  encodeDeterministicCbor(protectedHeader),
  fixture.protectedHex,
  "ApprovalGrant protected header",
);

process.stdout.write(`validated ${fixtureUrl.pathname}\n`);

const registrationFixtureUrl = new URL(
  "../schemas/fixtures/plan-registration-v0.json",
  import.meta.url,
);
const registrationFixture = JSON.parse(await readFile(registrationFixtureUrl, "utf8"));

assertExactKeys(registrationFixture, [
  "epochDigestHex",
  "epochSequence",
  "expiresAt",
  "installationIdHex",
  "objectType",
  "objectVersion",
  "payloadHex",
  "planDigestHex",
  "registrationIdHex",
  "registrationSequence",
  "supervisorIdHex",
]);
assertEqual(registrationFixture.objectType, "capsule.plan-registration", "object type");
assertEqual(registrationFixture.objectVersion, 0, "object version");
assertSafeUnsigned(registrationFixture.registrationSequence, "registrationSequence");
assertSafeUnsigned(registrationFixture.epochSequence, "epochSequence");
assertSafeUnsigned(registrationFixture.expiresAt, "expiresAt");

const registrationPayload = new Map([
  [1, registrationFixture.objectType],
  [2, registrationFixture.objectVersion],
  [3, decodeFixedHex(registrationFixture.registrationIdHex, 16, "registrationIdHex")],
  [4, registrationFixture.registrationSequence],
  [5, decodeFixedHex(registrationFixture.planDigestHex, 32, "planDigestHex")],
  [6, decodeFixedHex(registrationFixture.installationIdHex, 16, "installationIdHex")],
  [7, registrationFixture.epochSequence],
  [8, decodeFixedHex(registrationFixture.epochDigestHex, 32, "epochDigestHex")],
  [9, decodeFixedHex(registrationFixture.supervisorIdHex, 16, "supervisorIdHex")],
  [10, registrationFixture.expiresAt],
]);

assertHexEqual(
  encodeDeterministicCbor(registrationPayload),
  registrationFixture.payloadHex,
  "PlanRegistration payload",
);
process.stdout.write(`validated ${registrationFixtureUrl.pathname}\n`);

const planFixtureUrl = new URL("../schemas/fixtures/execution-plan-v0.json", import.meta.url);
const planFixture = JSON.parse(await readFile(planFixtureUrl, "utf8"));

assertExactKeys(planFixture, [
  "backendConfigurationDigestHex",
  "backendValidationRecordDigestHex",
  "epochDigestHex",
  "epochSequence",
  "expiresAt",
  "inlineInputByteLength",
  "inlineInputDigestHex",
  "inputSlot",
  "installationIdHex",
  "objectType",
  "objectVersion",
  "outputMaxJsonBytes",
  "outputSlot",
  "payloadHex",
  "policyDecisionDigestHex",
  "profileRegistryEntryDigestHex",
  "profileReviewAttestationDigestsHex",
  "runtimeBundleManifestDigestHex",
  "runtimeProfileAlias",
  "sourceByteLength",
  "sourceEntrypoint",
  "sourceManifestDigestHex",
  "trustSnapshotDigestHex",
  "wallTimeMs",
  "wallTimeOrigin",
]);
assertEqual(planFixture.objectType, "capsule.execution-plan", "object type");
assertEqual(planFixture.objectVersion, 0, "object version");
assertEqual(planFixture.inputSlot, "primary-data", "input slot");
assertEqual(planFixture.outputSlot, "transformed-json", "output slot");
if (!["requested", "trusted-default"].includes(planFixture.wallTimeOrigin)) {
  throw new Error("wallTimeOrigin is unsupported");
}
if (
  !Array.isArray(planFixture.profileReviewAttestationDigestsHex) ||
  planFixture.profileReviewAttestationDigestsHex.length < 1 ||
  planFixture.profileReviewAttestationDigestsHex.length > 8
) {
  throw new Error("profileReviewAttestationDigestsHex must contain 1 through 8 digests");
}
for (const field of [
  "epochSequence",
  "sourceByteLength",
  "inlineInputByteLength",
  "wallTimeMs",
  "outputMaxJsonBytes",
  "expiresAt",
]) {
  assertSafeUnsigned(planFixture[field], field);
}

const planPayload = new Map([
  [1, planFixture.objectType],
  [2, planFixture.objectVersion],
  [3, decodeFixedHex(planFixture.installationIdHex, 16, "installationIdHex")],
  [4, planFixture.epochSequence],
  [5, decodeFixedHex(planFixture.epochDigestHex, 32, "epochDigestHex")],
  [6, decodeFixedHex(planFixture.sourceManifestDigestHex, 32, "sourceManifestDigestHex")],
  [7, planFixture.sourceEntrypoint],
  [8, planFixture.sourceByteLength],
  [9, planFixture.inputSlot],
  [10, decodeFixedHex(planFixture.inlineInputDigestHex, 32, "inlineInputDigestHex")],
  [11, planFixture.inlineInputByteLength],
  [12, planFixture.runtimeProfileAlias],
  [
    13,
    decodeFixedHex(
      planFixture.runtimeBundleManifestDigestHex,
      32,
      "runtimeBundleManifestDigestHex",
    ),
  ],
  [
    14,
    planFixture.profileReviewAttestationDigestsHex.map((value, index) =>
      decodeFixedHex(value, 32, `profileReviewAttestationDigestsHex[${index}]`),
    ),
  ],
  [
    15,
    decodeFixedHex(planFixture.profileRegistryEntryDigestHex, 32, "profileRegistryEntryDigestHex"),
  ],
  [
    16,
    decodeFixedHex(
      planFixture.backendValidationRecordDigestHex,
      32,
      "backendValidationRecordDigestHex",
    ),
  ],
  [
    17,
    decodeFixedHex(planFixture.backendConfigurationDigestHex, 32, "backendConfigurationDigestHex"),
  ],
  [18, decodeFixedHex(planFixture.trustSnapshotDigestHex, 32, "trustSnapshotDigestHex")],
  [19, decodeFixedHex(planFixture.policyDecisionDigestHex, 32, "policyDecisionDigestHex")],
  [20, planFixture.wallTimeMs],
  [21, planFixture.wallTimeOrigin],
  [22, planFixture.outputSlot],
  [23, planFixture.outputMaxJsonBytes],
  [24, planFixture.expiresAt],
]);

assertHexEqual(
  encodeDeterministicCbor(planPayload),
  planFixture.payloadHex,
  "ExecutionPlan payload",
);
process.stdout.write(`validated ${planFixtureUrl.pathname}\n`);

function encodeDeterministicCbor(value) {
  if (value instanceof Uint8Array) {
    return concatenate([encodeArgument(2, value.length), value]);
  }
  if (typeof value === "string") {
    const encoded = new TextEncoder().encode(value);
    return concatenate([encodeArgument(3, encoded.length), encoded]);
  }
  if (typeof value === "number") {
    if (!Number.isSafeInteger(value)) {
      throw new Error(`CBOR integer is outside the safe range: ${value}`);
    }
    return value >= 0 ? encodeArgument(0, value) : encodeArgument(1, -1 - value);
  }
  if (Array.isArray(value)) {
    return concatenate([
      encodeArgument(4, value.length),
      ...value.map((child) => encodeDeterministicCbor(child)),
    ]);
  }
  if (value instanceof Map) {
    const encodedEntries = [...value.entries()].map(([key, child]) => [
      encodeDeterministicCbor(key),
      encodeDeterministicCbor(child),
    ]);
    encodedEntries.sort(([left], [right]) => {
      if (left.length !== right.length) {
        return left.length - right.length;
      }
      return Buffer.compare(left, right);
    });
    return concatenate([encodeArgument(5, encodedEntries.length), ...encodedEntries.flat()]);
  }
  throw new Error(`unsupported fixture value: ${String(value)}`);
}

function encodeArgument(majorType, argument) {
  assertSafeUnsigned(argument, "CBOR argument");
  if (argument < 24) {
    return Uint8Array.of((majorType << 5) | argument);
  }
  if (argument <= 0xff) {
    return Uint8Array.of((majorType << 5) | 24, argument);
  }
  if (argument <= 0xffff) {
    return Uint8Array.of((majorType << 5) | 25, argument >> 8, argument & 0xff);
  }
  if (argument <= 0xffffffff) {
    return Uint8Array.of(
      (majorType << 5) | 26,
      Math.floor(argument / 0x1000000) & 0xff,
      Math.floor(argument / 0x10000) & 0xff,
      Math.floor(argument / 0x100) & 0xff,
      argument & 0xff,
    );
  }
  const encoded = new Uint8Array(9);
  encoded[0] = (majorType << 5) | 27;
  new DataView(encoded.buffer).setBigUint64(1, BigInt(argument));
  return encoded;
}

function concatenate(parts) {
  const output = new Uint8Array(parts.reduce((total, part) => total + part.length, 0));
  let offset = 0;
  for (const part of parts) {
    output.set(part, offset);
    offset += part.length;
  }
  return output;
}

function decodeFixedHex(value, byteLength, label) {
  const decoded = decodeHex(value, label);
  if (decoded.length !== byteLength) {
    throw new Error(`${label} must encode exactly ${byteLength} bytes`);
  }
  return decoded;
}

function decodeHex(value, label) {
  if (typeof value !== "string" || !/^(?:[0-9a-f]{2})+$/u.test(value)) {
    throw new Error(`${label} must be nonempty canonical lowercase hexadecimal`);
  }
  return Uint8Array.from(Buffer.from(value, "hex"));
}

function assertSafeUnsigned(value, label) {
  if (!Number.isSafeInteger(value) || value < 0) {
    throw new Error(`${label} must be a nonnegative safe integer`);
  }
}

function assertExactKeys(value, expected) {
  if (value === null || typeof value !== "object" || Array.isArray(value)) {
    throw new Error("fixture object is malformed");
  }
  const actual = Object.keys(value).sort();
  const wanted = [...expected].sort();
  if (actual.length !== wanted.length || actual.some((key, index) => key !== wanted[index])) {
    throw new Error(`fixture keys differ: got ${actual.join(", ")}, want ${wanted.join(", ")}`);
  }
}

function assertHexEqual(actual, expected, label) {
  const actualHex = Buffer.from(actual).toString("hex");
  assertEqual(actualHex, expected, `${label} bytes`);
}

function assertEqual(actual, expected, label) {
  if (actual !== expected) {
    throw new Error(`${label}: got ${String(actual)}, want ${String(expected)}`);
  }
}
