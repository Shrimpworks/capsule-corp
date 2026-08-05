import { createECDH, createHmac } from "node:crypto";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { fromHex, sequence, sha256 } from "./lib/fixture-bytes.mjs";

const root = new URL("../schemas/conformance/i2b-bootstrap-v0/", import.meta.url);
const check = process.argv.includes("--check");
const privateScalar = fromHex("1f1e1d1c1b1a191817161514131211100f0e0d0c0b0a09080706050403020100");
const ecdh = createECDH("prime256v1");
ecdh.setPrivateKey(privateScalar);
const point = ecdh.getPublicKey(undefined, "uncompressed");
const publicKey = encodeMap([
  [1, 2],
  [3, -7],
  [-1, 1],
  [-2, point.subarray(1, 33)],
  [-3, point.subarray(33, 65)],
]);
const keyId = sha256(publicKey);

const installationId = sequence(16, 0x10);
const supervisorId = sequence(16, 0x30);
const nonce = sequence(32, 0x50);
const epochDigest = repeated(32, 0xa6);
const authorizationIdentity = sha256(
  concatenate([
    Buffer.from("capsule.installation-root-key-authorization/v0\0", "utf8"),
    keyId,
    installationId,
    uint64be(1),
    epochDigest,
  ]),
);

const request = new Map([
  [1, "capsule.supervisor-bootstrap-request"],
  [2, 0],
  [3, "capsule.installation.bootstrap.request"],
  [4, "capsule.execution-supervisor.bootstrap"],
  [5, installationId],
  [6, publicKey],
  [7, keyId],
  [8, authorizationIdentity],
  [9, supervisorId],
  [10, 501],
  [11, repeated(32, 0xa1)],
  [12, repeated(32, 0xa2)],
  [13, repeated(32, 0xa3)],
  [14, repeated(32, 0xa4)],
  [15, 1],
  [16, epochDigest],
  [17, "supervisor-private-app-sandbox"],
  [18, "supervisor.state"],
  [19, "supervisor.bootstrap-request"],
  [20, "supervisor.bootstrap-record"],
  [21, "supervisor.owner"],
  [22, 384],
  [23, 1],
  [24, "darwin-openat-flock-v0"],
  [25, "supervisor.store"],
  [26, "capsule.supervisor-store/fixed-v1"],
  [27, "conformance-non-product-no-guest"],
  [28, false],
  [29, false],
  [30, nonce],
  [31, 2_000_000_000],
  [32, 2_000_000_000],
  [33, 2_000_000_300],
  [34, "one-shot-exact-payload-replay-only"],
]);

const requestPayload = encode(request);
const requestEnvelope = envelope(
  "application/capsule.supervisor-bootstrap-request+cbor;v=0",
  requestPayload,
);

const record = new Map([
  [1, "capsule.supervisor-bootstrap-record"],
  [2, 0],
  [3, "capsule.installation.bootstrap.record"],
  [4, "capsule.execution-supervisor"],
  [5, sha256(requestPayload)],
  [6, sha256(requestEnvelope)],
  [7, requestEnvelope.length],
  [8, nonce],
  [9, installationId],
  [10, publicKey],
  [11, keyId],
  [12, authorizationIdentity],
  [13, supervisorId],
  [14, 501],
  [15, repeated(32, 0xa1)],
  [16, repeated(32, 0xa2)],
  [17, repeated(32, 0xa3)],
  [18, repeated(32, 0xa4)],
  [19, 1],
  [20, epochDigest],
  [21, "supervisor-private-app-sandbox"],
  [22, "supervisor.state"],
  [23, 1001],
  [24, 1002],
  [25, "directory"],
  [26, 501],
  [27, 448],
  [28, "supervisor.owner"],
  [29, 1001],
  [30, 1003],
  [31, "regular-file"],
  [32, 501],
  [33, 384],
  [34, 1],
  [35, "darwin-openat-flock-v0"],
  [36, "supervisor.store"],
  [37, 1001],
  [38, "regular-file"],
  [39, 501],
  [40, 384],
  [41, 1],
  [42, "capsule.supervisor-store/fixed-v1"],
  [43, "conformance-non-product-no-guest"],
  [44, 4096],
  [45, repeated(32, 0xa7)],
  [46, "supervisor.bootstrap-request"],
  [47, "supervisor.bootstrap-record"],
  [48, "retain-exact-envelope-mode-0400-single-link"],
  [49, "retain-exact-envelope-file-and-epoch-anchor"],
  [50, "installation-candidate-to-protected-root-disabled"],
  [51, "protected-root-validated-disabled"],
  [52, 2_000_000_100],
  [53, false],
  [54, false],
  [55, false],
  [56, false],
  [57, "commit-once"],
  [58, "exact-retained-envelope-response-replay-only"],
  [59, 2_000_000_000],
  [60, 2_000_000_000],
  [61, 2_000_000_300],
  [62, "one-shot-exact-payload-replay-only"],
  [63, "capsule.supervisor-bootstrap-request"],
  [64, 0],
  [65, "capsule.installation.bootstrap.request"],
  [66, "capsule.execution-supervisor.bootstrap"],
  [67, 2_000_000_100],
  [68, 2_000_000_100],
  [69, 2_000_000_300],
]);
const recordPayload = encode(record);
const recordEnvelope = envelope(
  "application/capsule.supervisor-bootstrap-record+cbor;v=0",
  recordPayload,
);

const maximumRequest = cloneMap(request);
maximumRequest.set(10, 4_294_967_295);
maximumRequest.set(31, 9_007_199_254_740_691);
maximumRequest.set(32, 9_007_199_254_740_691);
maximumRequest.set(33, 9_007_199_254_740_991);
const maximumRequestPayload = encode(maximumRequest);
const maximumRequestEnvelope = envelope(
  "application/capsule.supervisor-bootstrap-request+cbor;v=0",
  maximumRequestPayload,
);

const maximumRecord = cloneMap(record);
for (const label of [23, 24, 29, 30, 37, 44]) maximumRecord.set(label, 9_007_199_254_740_991);
for (const label of [14, 26, 32, 39]) maximumRecord.set(label, 4_294_967_295);
maximumRecord.set(5, sha256(maximumRequestPayload));
maximumRecord.set(6, sha256(maximumRequestEnvelope));
maximumRecord.set(7, maximumRequestEnvelope.length);
maximumRecord.set(59, 9_007_199_254_740_691);
maximumRecord.set(60, 9_007_199_254_740_691);
maximumRecord.set(61, 9_007_199_254_740_991);
maximumRecord.set(52, 9_007_199_254_740_691);
maximumRecord.set(67, 9_007_199_254_740_691);
maximumRecord.set(68, 9_007_199_254_740_691);
maximumRecord.set(69, 9_007_199_254_740_991);
const maximumRecordPayload = encode(maximumRecord);
const maximumRecordEnvelope = envelope(
  "application/capsule.supervisor-bootstrap-record+cbor;v=0",
  maximumRecordPayload,
);

const files = new Map();
add("request/ordinary.payload.cbor", requestPayload);
add("request/ordinary.cose", requestEnvelope);
add("record/ordinary.payload.cbor", recordPayload);
add("record/ordinary.cose", recordEnvelope);
add("request/calculated-maximum.payload.cbor", maximumRequestPayload);
add("request/calculated-maximum.cose", maximumRequestEnvelope);
add("record/calculated-maximum.payload.cbor", maximumRecordPayload);
add("record/calculated-maximum.cose", maximumRecordEnvelope);

const cases = [];
accept(
  "request-ordinary",
  "request",
  "request/ordinary.cose",
  "admit-once",
  2_000_000_001,
  "fresh",
);
accept(
  "request-exact-pending-replay",
  "request",
  "request/ordinary.cose",
  "resume-exact",
  2_000_000_001,
  "pending-exact",
);
const complementaryRequestEnvelope = envelope(
  "application/capsule.supervisor-bootstrap-request+cbor;v=0",
  requestPayload,
  { complementaryS: true },
);
add("request/complementary-s.cose", complementaryRequestEnvelope);
accept(
  "request-complementary-s-fresh",
  "request",
  "request/complementary-s.cose",
  "admit-once",
  2_000_000_001,
  "fresh",
);
accept(
  "request-complementary-s-pending-replay",
  "request",
  "request/complementary-s.cose",
  "resume-exact",
  2_000_000_001,
  "pending-exact",
);
accept("record-ordinary", "record", "record/ordinary.cose", "commit-once", 2_000_000_101, "fresh");
accept(
  "record-response-loss",
  "record",
  "record/ordinary.cose",
  "return-retained-envelope",
  2_000_000_101,
  "completed-exact",
);
const complementaryRecordEnvelope = envelope(
  "application/capsule.supervisor-bootstrap-record+cbor;v=0",
  recordPayload,
  { complementaryS: true },
);
add("record/complementary-s.cose", complementaryRecordEnvelope);
accept(
  "record-complementary-s-response-loss",
  "record",
  "record/complementary-s.cose",
  "return-retained-envelope",
  2_000_000_101,
  "completed-exact",
);
accept(
  "request-calculated-maximum",
  "request",
  "request/calculated-maximum.cose",
  "admit-once",
  9_007_199_254_740_691,
  "fresh",
  true,
);
accept(
  "record-calculated-maximum",
  "record",
  "record/calculated-maximum.cose",
  "commit-once",
  9_007_199_254_740_691,
  "fresh",
  true,
  "maximum",
);

mutatedRequest("request-wrong-domain", 1, "capsule.supervisor-bootstrap-record");
mutatedRequest("request-wrong-version", 2, 1);
mutatedRequest("request-wrong-purpose", 3, "capsule.installation.bootstrap.record");
mutatedRequest("request-wrong-audience", 4, "capsule.execution-supervisor");
mutatedRequest("request-wrong-installation", 5, repeated(16, 0x91));
mutatedRequest("request-wrong-epoch-sequence", 15, 2);
mutatedRequest("request-wrong-epoch-digest", 16, repeated(32, 0x92));
mutatedRequest("request-wrong-supervisor", 9, repeated(16, 0x93));
mutatedRequest("request-wrong-release", 11, repeated(32, 0x94));
mutatedRequest("request-wrong-root", 18, "other.state");
mutatedRequest("request-wrong-lock", 24, "other-lock-profile");
mutatedRequest("request-wrong-store", 26, "other-store-format");
mutatedRequest("request-zero-nonce", 30, repeated(32, 0));
mutatedRequest("request-unknown-signer-key", 7, repeated(32, 0x95));
const requestUnknownAuthorization = cloneMap(request);
requestUnknownAuthorization.set(8, repeated(32, 0x96));
signedRequestCase(
  "request-unknown-signer-authorization",
  requestUnknownAuthorization,
  2_000_000_001,
  true,
);

const future = cloneMap(request);
future.set(32, 2_000_000_100);
future.set(33, 2_000_000_200);
signedRequestCase("request-future-time", future, 2_000_000_050, true);
signedRequestCase("request-stale-time", request, 2_000_000_300, true);
signedRequestCase("request-expired-time", request, 2_000_000_301, true);
const invalidWindow = cloneMap(request);
invalidWindow.set(33, 2_000_000_301);
signedRequestCase("request-window-over-300", invalidWindow, 2_000_000_001, true);

mutatedRecord("record-wrong-domain", 1, "capsule.supervisor-bootstrap-request");
mutatedRecord("record-wrong-version", 2, 1);
mutatedRecord("record-wrong-purpose", 3, "capsule.installation.bootstrap.request");
mutatedRecord("record-wrong-audience", 4, "capsule.execution-supervisor.bootstrap");
mutatedRecord("record-wrong-installation", 9, repeated(16, 0x81));
mutatedRecord("record-wrong-epoch", 20, repeated(32, 0x82));
mutatedRecord("record-wrong-supervisor", 13, repeated(16, 0x83));
mutatedRecord("record-wrong-release", 15, repeated(32, 0x84));
mutatedRecord("record-wrong-root", 22, "other.state");
mutatedRecord("record-wrong-lock", 35, "other-lock-profile");
mutatedRecord("record-wrong-store", 42, "other-store-format");
mutatedRecord("record-zero-nonce", 8, repeated(32, 0));
mutatedRecord("record-unknown-signer-key", 11, repeated(32, 0x87));
const recordUnknownAuthorization = cloneMap(record);
recordUnknownAuthorization.set(12, repeated(32, 0x88));
signedRecordCase(
  "record-unknown-signer-authorization",
  recordUnknownAuthorization,
  2_000_000_101,
  true,
);
mutatedRecord("record-request-payload-substitution", 5, repeated(32, 0x85));
mutatedRecord("record-request-envelope-substitution", 6, repeated(32, 0x86));
mutatedRecord("record-wrong-request-time-projection", 59, 2_000_000_001);
mutatedRecord(
  "record-wrong-request-purpose-projection",
  65,
  "capsule.installation.bootstrap.record",
);

const futureRecord = cloneMap(record);
futureRecord.set(68, 2_000_000_200);
signedRecordCase("record-future-time", futureRecord, 2_000_000_150, true);
rejectContext("record-stale-time", "record", "record/ordinary.cose", 2_000_000_300, "fresh");
rejectContext("record-expired-time", "record", "record/ordinary.cose", 2_000_000_301, "fresh");

rejectContext(
  "request-reused-nonce",
  "request",
  "request/ordinary.cose",
  2_000_000_001,
  "pending-other",
);
rejectContext(
  "request-completed-replay",
  "request",
  "request/ordinary.cose",
  2_000_000_001,
  "completed-exact",
);
rejectContext(
  "record-reused-nonce",
  "record",
  "record/ordinary.cose",
  2_000_000_101,
  "completed-other",
);

const noncanonicalPayload = encodeMap([...request.entries()].reverse(), false);
addSignedReject("request-noncanonical-cbor", "request", noncanonicalPayload);
const duplicatePayload = concatenate([
  replaceMapCount(requestPayload, 35),
  encode(1),
  encode(request.get(1)),
]);
addSignedReject("request-duplicate-field", "request", duplicatePayload);
const unknownPayload = concatenate([
  replaceMapCount(requestPayload, 35),
  encode(35),
  encode(false),
]);
addSignedReject("request-unknown-field", "request", unknownPayload);
addSignedReject(
  "request-payload-trailing",
  "request",
  concatenate([requestPayload, Uint8Array.of(0)]),
);
addSignedReject("request-cross-object-payload", "request", recordPayload);

const noncanonicalRecordPayload = encodeMap([...record.entries()].reverse(), false);
addSignedReject("record-noncanonical-cbor", "record", noncanonicalRecordPayload);
const duplicateRecordPayload = concatenate([
  replaceMapCount(recordPayload, 70),
  encode(1),
  encode(record.get(1)),
]);
addSignedReject("record-duplicate-field", "record", duplicateRecordPayload);
const unknownRecordPayload = concatenate([
  replaceMapCount(recordPayload, 70),
  encode(70),
  encode(false),
]);
addSignedReject("record-unknown-field", "record", unknownRecordPayload);
addSignedReject(
  "record-payload-trailing",
  "record",
  concatenate([recordPayload, Uint8Array.of(0)]),
);
addSignedReject("record-cross-object-payload", "record", requestPayload);

// Restore the exact COSE profile independently for both signed objects. The
// envelopes below are validly signed wherever signature verification is the
// control under test; structural profiles still fail before signature use.
for (const profile of [
  {
    object: "request",
    media: "application/capsule.supervisor-bootstrap-request+cbor;v=0",
    payload: requestPayload,
    map: request,
  },
  {
    object: "record",
    media: "application/capsule.supervisor-bootstrap-record+cbor;v=0",
    payload: recordPayload,
    map: record,
  },
]) {
  const wrongType = cloneMap(profile.map);
  wrongType.set(2, "0");
  addSignedReject(`${profile.object}-wrong-field-type`, profile.object, encode(wrongType));

  const wrongAlgorithm = protectedHeaderMap(profile.media);
  wrongAlgorithm.set(1, -8);
  addEnvelopeReject(`${profile.object}-protected-wrong-algorithm`, profile, {
    protectedBytes: encode(wrongAlgorithm),
  });

  const wrongMedia = protectedHeaderMap(profile.media);
  wrongMedia.set(3, "application/capsule.wrong+cbor;v=0");
  addEnvelopeReject(`${profile.object}-protected-wrong-content-type`, profile, {
    protectedBytes: encode(wrongMedia),
  });

  const unknownProtected = protectedHeaderMap(profile.media);
  unknownProtected.set(1000, 1);
  addEnvelopeReject(`${profile.object}-protected-unknown-label`, profile, {
    protectedBytes: encode(unknownProtected),
  });

  addEnvelopeReject(`${profile.object}-protected-noncanonical-order`, profile, {
    protectedBytes: encodeMap([...protectedHeaderMap(profile.media).entries()].reverse(), false),
  });
  addEnvelopeReject(`${profile.object}-nonempty-unprotected`, profile, {
    unprotectedBytes: encode(new Map([[4, keyId]])),
  });
  addEnvelopeReject(`${profile.object}-nonempty-external-aad`, profile, {
    externalAAD: Buffer.from("capsule-test-nonempty-external-aad", "utf8"),
  });
  addEnvelopeReject(`${profile.object}-detached-payload`, profile, { detachedPayload: true });
  addEnvelopeReject(`${profile.object}-der-signature`, profile, { derSignature: true });
}

// Record verification must prove that retained request envelope and payload
// are one exact signed pair, then bind every repeated request field. These
// records are validly signed so signature verification alone cannot save an
// omitted Capsule-owned cross-object check.
const splitPairRecord = cloneMap(maximumRecord);
splitPairRecord.set(6, sha256(requestEnvelope));
splitPairRecord.set(7, requestEnvelope.length);
addRecordBindingReject(
  "record-request-envelope-payload-split",
  splitPairRecord,
  "request/ordinary.cose",
  "request/calculated-maximum.payload.cbor",
);

const componentProfileRequest = cloneMap(request);
componentProfileRequest.set(13, repeated(32, 0xb3));
const componentProfileRequestPayload = encode(componentProfileRequest);
const componentProfileRequestEnvelope = envelope(
  "application/capsule.supervisor-bootstrap-request+cbor;v=0",
  componentProfileRequestPayload,
);
add("request/record-binding-component-profile.payload.cbor", componentProfileRequestPayload);
add("request/record-binding-component-profile.cose", componentProfileRequestEnvelope);
const componentProfileRecord = cloneMap(record);
componentProfileRecord.set(5, sha256(componentProfileRequestPayload));
componentProfileRecord.set(6, sha256(componentProfileRequestEnvelope));
componentProfileRecord.set(7, componentProfileRequestEnvelope.length);
addRecordBindingReject(
  "record-request-component-profile-substitution",
  componentProfileRecord,
  "request/record-binding-component-profile.cose",
  "request/record-binding-component-profile.payload.cbor",
);

const tamperedBoundRequestEnvelope = Buffer.from(requestEnvelope);
tamperedBoundRequestEnvelope[tamperedBoundRequestEnvelope.length - 1] ^= 1;
add("request/record-binding-tampered-signature.cose", tamperedBoundRequestEnvelope);
const tamperedRequestRecord = cloneMap(record);
tamperedRequestRecord.set(6, sha256(tamperedBoundRequestEnvelope));
addRecordBindingReject(
  "record-request-signature-substitution",
  tamperedRequestRecord,
  "request/record-binding-tampered-signature.cose",
  "request/ordinary.payload.cbor",
);

const signatureTamper = Buffer.from(requestEnvelope);
signatureTamper[signatureTamper.length - 1] ^= 1;
addRejectBytes("request-envelope-substitution", "request", signatureTamper);
addRejectBytes(
  "request-envelope-trailing",
  "request",
  concatenate([requestEnvelope, Uint8Array.of(0)]),
);
const wrongSignature = envelope(
  "application/capsule.supervisor-bootstrap-request+cbor;v=0",
  requestPayload,
  { signatureLength: 63 },
);
addRejectBytes("request-wrong-signature-shape", "request", wrongSignature);
const recordSignatureTamper = Buffer.from(recordEnvelope);
recordSignatureTamper[recordSignatureTamper.length - 1] ^= 1;
addRejectBytes("record-envelope-substitution", "record", recordSignatureTamper);
addRejectBytes(
  "record-envelope-trailing",
  "record",
  concatenate([recordEnvelope, Uint8Array.of(0)]),
);
const wrongRecordSignature = envelope(
  "application/capsule.supervisor-bootstrap-record+cbor;v=0",
  recordPayload,
  { signatureLength: 63 },
);
addRejectBytes("record-wrong-signature-shape", "record", wrongRecordSignature);
addRejectBytes("request-envelope-cap-plus-one", "request", repeated(4097, 0));
addRejectBytes("record-envelope-cap-plus-one", "record", repeated(6145, 0));
addRejectBytes(
  "request-payload-cap-plus-one",
  "request",
  framedEnvelope(
    repeated(2049, 0),
    repeated(64, 0),
    "application/capsule.supervisor-bootstrap-request+cbor;v=0",
  ),
);
addRejectBytes(
  "record-payload-cap-plus-one",
  "record",
  framedEnvelope(
    repeated(4097, 0),
    repeated(64, 0),
    "application/capsule.supervisor-bootstrap-record+cbor;v=0",
  ),
);
addRejectBytes(
  "request-calculated-maximum-plus-one",
  "request",
  framedEnvelope(
    repeated(maximumRequestPayload.length + 1, 0),
    repeated(64, 0),
    "application/capsule.supervisor-bootstrap-request+cbor;v=0",
  ),
);
addRejectBytes(
  "record-calculated-maximum-plus-one",
  "record",
  framedEnvelope(
    repeated(maximumRecordPayload.length + 1, 0),
    repeated(64, 0),
    "application/capsule.supervisor-bootstrap-record+cbor;v=0",
  ),
);

const manifest = {
  manifestVersion: "capsule.i2b-bootstrap-conformance/v0",
  generator: "independent-node-rfc6979-p256",
  publicKeyCoseKeyHex: Buffer.from(publicKey).toString("hex"),
  keyIdHex: Buffer.from(keyId).toString("hex"),
  authorizationIdentityHex: Buffer.from(authorizationIdentity).toString("hex"),
  maxima: {
    request: {
      payload: maximumRequestPayload.length,
      protected: protectedHeader("application/capsule.supervisor-bootstrap-request+cbor;v=0")
        .length,
      envelope: maximumRequestEnvelope.length,
    },
    record: {
      payload: maximumRecordPayload.length,
      protected: protectedHeader("application/capsule.supervisor-bootstrap-record+cbor;v=0").length,
      envelope: maximumRecordEnvelope.length,
    },
  },
  effects: {
    stateMutations: 0,
    keyOperations: 0,
    keychainOperations: 0,
    ipcEndpoints: 0,
    processes: 0,
    services: 0,
    runtimeOperations: 0,
    backendOperations: 0,
    guestOperations: 0,
  },
  cases,
};
add("manifest.json", Buffer.from(`${JSON.stringify(manifest, null, 2)}\n`));

for (const [path, contents] of files) {
  const url = new URL(path, root);
  if (check) {
    const current = await readFile(url);
    if (!current.equals(Buffer.from(contents)))
      throw new Error(`stale I2B bootstrap fixture: ${path}`);
  } else {
    await mkdir(new URL("./", url), { recursive: true });
    await writeFile(url, contents);
  }
}
process.stdout.write(
  `${check ? "verified" : "generated"} I2B bootstrap corpus: ${cases.length} cases; request maxima ${maximumRequestPayload.length}/${maximumRequestEnvelope.length}; record maxima ${maximumRecordPayload.length}/${maximumRecordEnvelope.length}\n`,
);

function mutatedRequest(id, label, value) {
  const m = cloneMap(request);
  m.set(label, value);
  signedRequestCase(id, m, 2_000_000_001, false);
}
function signedRequestCase(id, map, trustedNow, selfExpected) {
  const path = `request/${id}.cose`;
  add(path, envelope("application/capsule.supervisor-bootstrap-request+cbor;v=0", encode(map)));
  cases.push({
    id,
    object: "request",
    fixture: path,
    expected: "REFUSE",
    trustedNow,
    replay: "fresh",
    selfExpected,
  });
}
function mutatedRecord(id, label, value) {
  const m = cloneMap(record);
  m.set(label, value);
  const path = `record/${id}.cose`;
  add(path, envelope("application/capsule.supervisor-bootstrap-record+cbor;v=0", encode(m)));
  cases.push({
    id,
    object: "record",
    fixture: path,
    expected: "REFUSE",
    trustedNow: 2_000_000_101,
    replay: "fresh",
    selfExpected: false,
  });
}
function signedRecordCase(id, map, trustedNow, selfExpected) {
  const path = `record/${id}.cose`;
  add(path, envelope("application/capsule.supervisor-bootstrap-record+cbor;v=0", encode(map)));
  cases.push({
    id,
    object: "record",
    fixture: path,
    expected: "REFUSE",
    trustedNow,
    replay: "fresh",
    selfExpected,
  });
}
function accept(
  id,
  object,
  fixture,
  decision,
  trustedNow,
  replay,
  selfExpected = false,
  requestVariant = "ordinary",
) {
  cases.push({
    id,
    object,
    fixture,
    expected: "ACCEPT",
    decision,
    trustedNow,
    replay,
    selfExpected,
    requestVariant,
  });
}
function rejectContext(id, object, fixture, trustedNow, replay) {
  cases.push({ id, object, fixture, expected: "REFUSE", trustedNow, replay, selfExpected: false });
}
function addSignedReject(id, object, payload) {
  const media =
    object === "request"
      ? "application/capsule.supervisor-bootstrap-request+cbor;v=0"
      : "application/capsule.supervisor-bootstrap-record+cbor;v=0";
  const path = `${object}/${id}.cose`;
  add(path, envelope(media, payload));
  cases.push({
    id,
    object,
    fixture: path,
    expected: "REFUSE",
    trustedNow: object === "request" ? 2_000_000_001 : 2_000_000_101,
    replay: "fresh",
    selfExpected: false,
  });
}
function addRejectBytes(id, object, value) {
  const path = `${object}/${id}.cose`;
  add(path, value);
  cases.push({
    id,
    object,
    fixture: path,
    expected: "REFUSE",
    trustedNow: object === "request" ? 2_000_000_001 : 2_000_000_101,
    replay: "fresh",
    selfExpected: false,
  });
}
function addEnvelopeReject(id, profile, options) {
  addRejectBytes(id, profile.object, envelope(profile.media, profile.payload, options));
}
function addRecordBindingReject(id, recordMap, boundRequestEnvelope, boundRequestPayload) {
  const path = `record/${id}.cose`;
  add(
    path,
    envelope("application/capsule.supervisor-bootstrap-record+cbor;v=0", encode(recordMap)),
  );
  cases.push({
    id,
    object: "record",
    fixture: path,
    expected: "REFUSE",
    trustedNow: 2_000_000_101,
    replay: "fresh",
    selfExpected: true,
    requestEnvelope: boundRequestEnvelope,
    requestPayload: boundRequestPayload,
  });
}
function add(path, value) {
  files.set(path, Buffer.from(value));
}

function protectedHeader(media) {
  return encode(protectedHeaderMap(media));
}
function protectedHeaderMap(media) {
  return new Map([
    [1, -7],
    [3, media],
    [4, keyId],
  ]);
}
function envelope(media, payload, options = {}) {
  const protectedBytes = options.protectedBytes ?? protectedHeader(media);
  const externalAAD = options.externalAAD ?? new Uint8Array();
  let signature =
    options.signatureLength === 63
      ? repeated(63, 1)
      : sign(sigStructure(protectedBytes, payload, externalAAD));
  if (options.complementaryS) signature = complementarySignature(signature);
  if (options.derSignature) signature = derSignature(signature);
  return concatenate([
    Uint8Array.of(0xd2, 0x84),
    encodeBytes(protectedBytes),
    options.unprotectedBytes ?? Uint8Array.of(0xa0),
    options.detachedPayload ? Uint8Array.of(0xf6) : encodeBytes(payload),
    encodeBytes(signature),
  ]);
}
function complementarySignature(signature) {
  const n = BigInt("0xffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc632551");
  const result = Buffer.from(signature);
  const s = toBigInt(result.subarray(32));
  if (s === 0n || s >= n) throw new Error("invalid fixture S scalar");
  result.set(fromBigInt(n - s, 32), 32);
  return result;
}
function derSignature(signature) {
  const integer = (scalar) => {
    let value = Buffer.from(scalar);
    while (value.length > 1 && value[0] === 0) value = value.subarray(1);
    if (value[0] & 0x80) value = concatenate([Uint8Array.of(0), value]);
    return concatenate([Uint8Array.of(0x02, value.length), value]);
  };
  const r = integer(signature.subarray(0, 32));
  const s = integer(signature.subarray(32, 64));
  return concatenate([Uint8Array.of(0x30, r.length + s.length), r, s]);
}
function framedEnvelope(payload, signature, media) {
  return concatenate([
    Uint8Array.of(0xd2, 0x84),
    encodeBytes(protectedHeader(media)),
    Uint8Array.of(0xa0),
    encodeBytes(payload),
    encodeBytes(signature),
  ]);
}
function sigStructure(protectedBytes, payload, externalAAD = new Uint8Array()) {
  return encode(["Signature1", protectedBytes, externalAAD, payload]);
}

function sign(message) {
  const hash = sha256(message);
  const n = BigInt("0xffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc632551");
  const d = toBigInt(privateScalar);
  const z = toBigInt(hash) % n;
  let v = repeated(32, 1),
    k = repeated(32, 0);
  k = hmac(k, concatenate([v, Uint8Array.of(0), privateScalar, fromBigInt(z, 32)]));
  v = hmac(k, v);
  k = hmac(k, concatenate([v, Uint8Array.of(1), privateScalar, fromBigInt(z, 32)]));
  v = hmac(k, v);
  for (;;) {
    v = hmac(k, v);
    const candidate = toBigInt(v);
    if (candidate > 0n && candidate < n) {
      const pointK = createECDH("prime256v1");
      pointK.setPrivateKey(fromBigInt(candidate, 32));
      const r = toBigInt(pointK.getPublicKey(undefined, "uncompressed").subarray(1, 33)) % n;
      const s = (inverse(candidate, n) * (z + r * d)) % n;
      if (r !== 0n && s !== 0n) return concatenate([fromBigInt(r, 32), fromBigInt(s, 32)]);
    }
    k = hmac(k, concatenate([v, Uint8Array.of(0)]));
    v = hmac(k, v);
  }
}

function encode(value) {
  if (value instanceof Uint8Array || Buffer.isBuffer(value)) return encodeBytes(value);
  if (typeof value === "string") {
    const b = Buffer.from(value, "utf8");
    return concatenate([head(3, b.length), b]);
  }
  if (typeof value === "boolean") return Uint8Array.of(value ? 0xf5 : 0xf4);
  if (typeof value === "number") {
    if (!Number.isSafeInteger(value)) throw new Error("unsafe integer");
    return value >= 0 ? head(0, value) : head(1, -1 - value);
  }
  if (Array.isArray(value)) return concatenate([head(4, value.length), ...value.map(encode)]);
  if (value instanceof Map) return encodeMap([...value.entries()]);
  throw new Error(`unsupported CBOR value ${String(value)}`);
}
function encodeMap(entries, sort = true) {
  const encoded = entries.map(([k, v]) => [encode(k), encode(v)]);
  if (sort) encoded.sort(([a], [b]) => a.length - b.length || Buffer.compare(a, b));
  return concatenate([head(5, encoded.length), ...encoded.flat()]);
}
function encodeBytes(value) {
  return concatenate([head(2, value.length), value]);
}
function head(major, value) {
  if (value < 24) return Uint8Array.of((major << 5) | value);
  if (value <= 0xff) return Uint8Array.of((major << 5) | 24, value);
  if (value <= 0xffff) return Uint8Array.of((major << 5) | 25, value >> 8, value);
  if (value <= 0xffffffff)
    return Uint8Array.of(
      (major << 5) | 26,
      value / 2 ** 24,
      value / 2 ** 16,
      value / 2 ** 8,
      value,
    );
  const out = new Uint8Array(9);
  out[0] = (major << 5) | 27;
  new DataView(out.buffer).setBigUint64(1, BigInt(value));
  return out;
}
function replaceMapCount(value, count) {
  const prefix = head(5, count);
  const oldPrefix = value[0] === 0xb8 ? 2 : 1;
  return concatenate([prefix, value.subarray(oldPrefix)]);
}
function cloneMap(value) {
  return new Map(
    [...value].map(([key, child]) => [
      key,
      child instanceof Uint8Array || Buffer.isBuffer(child) ? Buffer.from(child) : child,
    ]),
  );
}
function hmac(key, value) {
  return createHmac("sha256", key).update(value).digest();
}
function inverse(a, n) {
  let [t, nt, r, nr] = [0n, 1n, n, a];
  while (nr) {
    const q = r / nr;
    [t, nt] = [nt, t - q * nt];
    [r, nr] = [nr, r - q * nr];
  }
  return t < 0n ? t + n : t;
}
function toBigInt(value) {
  return BigInt(`0x${Buffer.from(value).toString("hex")}`);
}
function fromBigInt(value, width) {
  return fromHex(value.toString(16).padStart(width * 2, "0"));
}
function uint64be(value) {
  const out = new Uint8Array(8);
  new DataView(out.buffer).setBigUint64(0, BigInt(value));
  return out;
}
function concatenate(parts) {
  return Buffer.concat(parts.map((part) => Buffer.from(part)));
}
function repeated(length, value) {
  return Buffer.alloc(length, value);
}
