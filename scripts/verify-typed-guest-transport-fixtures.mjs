import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import {
  COMPLETION_FRAME_MAXIMUM_BYTES,
  COMPLETION_ROLE,
  DISPOSITION_ORDER,
  decodeFrame,
  INPUT_FRAME_MAXIMUM_BYTES,
  INPUT_ROLE,
  SOURCE_ROLE,
  sha256,
  sha256Hex,
} from "./lib/typed-guest-transport.mjs";

const root = new URL("../schemas/conformance/typed-guest-transport-v1/", import.meta.url);
const manifestBytes = await readFile(new URL("manifest.json", root));
const manifest = JSON.parse(manifestBytes);
const bindings = Object.fromEntries(
  Object.entries(manifest.bindings).map(([key, value]) => [key, Buffer.from(value, "hex")]),
);

assert.equal(manifest.contract, "capsule.typed-guest-transport");
assert.equal(manifest.version, 1);
assert.equal(manifest.status, "passive-conformance-only");
assert.deepEqual(manifest.refusalPrecedence, DISPOSITION_ORDER);
assert.deepEqual(manifest.stateModel.states, [
  "BOUND",
  "ENDPOINTS_READY",
  "RUNNER_READY",
  "INPUT_TRANSFER",
  "LAUNCHER_VALIDATED",
  "CHILD_RUNNING",
  "RESULT_VALIDATED",
  "TRAILER_WRITTEN",
  "FRAME_OBSERVED",
  "TERMINAL_PROOF",
  "DURABLE_COMMIT",
  "COMPLETE",
]);
assert.deepEqual(manifest.statuses, {
  succeeded: 1,
  "workload-failed": 2,
  "result-invalid": 3,
  "child-terminated": 4,
});
assert.equal(manifest.endpointContract.completion.workloadAccess, false);
assert.equal(manifest.durableProjection.eofIsCommit, false);
assert.equal(manifest.durableProjection.frameIsObservationOnly, true);
assert.ok(Object.values(manifest.effects).every((effect) => effect === false));

const roles = new Set([SOURCE_ROLE, INPUT_ROLE, COMPLETION_ROLE]);
for (const fixture of manifest.cases) {
  assert.ok(roles.has(fixture.role), fixture.id);
  const bytes = await readFile(new URL(fixture.path, root));
  assert.equal(bytes.length, fixture.bytes, fixture.id);
  assert.equal(sha256Hex(bytes), fixture.sha256, fixture.id);
  const decoded = decodeFrame(bytes, fixture.role, bindings);
  if (fixture.decision === "accept") {
    assert.equal(decoded.ok, true, `${fixture.id}: ${decoded.disposition}`);
  } else {
    assert.equal(decoded.ok, false, fixture.id);
    assert.equal(decoded.disposition, fixture.disposition, fixture.id);
  }
}

const ordinarySource = manifest.cases.find(({ id }) => id === "accept-source-ordinary");
const ordinarySourceBytes = await readFile(new URL(ordinarySource.path, root));
const sourcePayload = ordinarySourceBytes.subarray(manifest.layout.inputHeaderBytes);
const independentlyBound = {
  ...bindings,
  expectedPayloadBytes: sourcePayload.length,
  expectedPayloadDigest: sha256(sourcePayload),
};
assert.equal(decodeFrame(ordinarySourceBytes, SOURCE_ROLE, independentlyBound).ok, true);
const wrongExpectedDigest = Buffer.from(independentlyBound.expectedPayloadDigest);
wrongExpectedDigest[0] ^= 1;
assert.equal(
  decodeFrame(ordinarySourceBytes, SOURCE_ROLE, {
    ...independentlyBound,
    expectedPayloadDigest: wrongExpectedDigest,
  }).disposition,
  "BINDING",
);

const sourceMaximum = manifest.cases.find(({ id }) => id === "accept-source-maximum");
const inputMaximum = manifest.cases.find(({ id }) => id === "accept-input-maximum");
const completionMaximum = manifest.cases.find(({ id }) => id === "accept-completion-maximum");
const completionCapPlusOne = manifest.cases.find(
  ({ id }) => id === "reject-completion-cap-plus-one",
);
assert.equal(sourceMaximum.bytes, INPUT_FRAME_MAXIMUM_BYTES);
assert.equal(inputMaximum.bytes, INPUT_FRAME_MAXIMUM_BYTES);
assert.equal(completionMaximum.bytes, COMPLETION_FRAME_MAXIMUM_BYTES);
assert.equal(completionCapPlusOne.bytes, COMPLETION_FRAME_MAXIMUM_BYTES + 1);

assert.deepEqual(
  manifest.stateCases.map(({ id, disposition, durableCompletion }) => [
    id,
    disposition,
    durableCompletion,
  ]),
  [
    ["cancel-before-g", "REFUSED_CANCELLED_BEFORE_START", false],
    ["cancel-during-source", "REFUSED_INPUT_CANCELLED", false],
    ["source-short-write-error", "REFUSED_SOURCE_TRANSPORT_FAULT", false],
    ["input-zero-progress-deadline", "REFUSED_INPUT_STALL", false],
    ["completion-reader-death", "REFUSED_COMPLETION_READER_FAULT", false],
    ["completion-reset-before-trailer", "REFUSED_MISSING_COMMIT", false],
    ["valid-frame-before-absence", "FRAME_ONLY_AWAIT_TERMINAL_PROOF", false],
    ["launcher-death-after-trailer", "REFUSED_LIFECYCLE", false],
    ["runner-zero-without-frame", "REFUSED_MISSING_COMMIT", false],
    ["caller-response-loss-after-commit", "REPLAY_BYTE_IDENTICAL", true],
    ["store-indeterminate", "RECOVERY_REQUIRED_FENCE", false],
    ["cancel-concurrent-before-terminal-proof", "REFUSED_CANCELLED", false],
    ["cancel-after-durable-commit", "REPLAY_BYTE_IDENTICAL", true],
  ],
);
assert.deepEqual(
  manifest.stateCases.filter(({ durableCompletion }) => durableCompletion).map(({ id }) => id),
  ["caller-response-loss-after-commit", "cancel-after-durable-commit"],
);
assert.deepEqual(
  manifest.restorationCases.map(({ id, detected }) => [id, detected]),
  [
    "source-input-endpoint-swap",
    "completion-source-endpoint-swap",
    "descriptor-alias",
    "descriptor-wrong-mode",
    "descriptor-cloexec-cleared",
    "descriptor-nonblocking-changed",
    "descriptor-inherited-by-workload",
    "wrong-attempt",
    "wrong-registration",
    "plan-digest-substitution",
    "runtime-profile-substitution",
    "payload-flood-after-valid-prefix",
    "early-trailer",
    "eof-as-commit",
    "runner-zero-as-success",
    "diagnostic-console-substitution",
    "implicit-console-restored",
    "vsock-restored",
    "network-device-restored",
    "virtiofs-restored",
    "cleanup-unresolved",
    "response-loss-result-substitution",
    "durable-commit-corrupt",
  ].map((id) => [id, true]),
);

process.stdout.write(
  `independently verified ${manifest.cases.length} Node frame cases, ${manifest.stateCases.length} state cases, and ${manifest.restorationCases.length} restoration cases\n`,
);
