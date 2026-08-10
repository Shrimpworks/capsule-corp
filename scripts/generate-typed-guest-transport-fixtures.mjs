import { mkdir, readFile, writeFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";
import {
  COMPLETION_FRAME_MAXIMUM_BYTES,
  COMPLETION_HEADER_BYTES,
  COMPLETION_RETAIN_MAXIMUM_BYTES,
  COMPLETION_ROLE,
  canonicalJSON,
  DISPOSITION_ORDER,
  decodeFrame,
  encodeCompletionFrame,
  encodeInputFrame,
  INPUT_FRAME_MAXIMUM_BYTES,
  INPUT_HEADER_BYTES,
  INPUT_ROLE,
  MAGIC,
  METHOD,
  PAYLOAD_MAXIMUM_BYTES,
  PROTOCOL_VERSION,
  SOURCE_ROLE,
  STATUS,
  sha256,
  sha256Hex,
  TRAILER_BYTES,
} from "./lib/typed-guest-transport.mjs";

const check = process.argv.slice(2).includes("--check");
const outputUrl = new URL("../schemas/conformance/typed-guest-transport-v1/", import.meta.url);
const outputPath = fileURLToPath(outputUrl);

export const bindings = Object.freeze({
  attemptId: Buffer.from("a1010102030405060708090a0b0c0d0e", "hex"),
  registrationId: Buffer.from("b2020102030405060708090a0b0c0d0e", "hex"),
  planDigest: sha256(Buffer.from("capsule.c5a.passive-plan-known-answer/v1")),
  runtimeProfileDigest: Buffer.from(
    "e390085caaaba73ebc19f95bc9871305e4f9268c2283d7394133fa4491f4ba82",
    "hex",
  ),
});

const files = new Map();
const cases = [];

function retain(name, bytes) {
  const exact = Buffer.from(bytes);
  files.set(name, exact);
  return Object.freeze({ path: name, bytes: exact.length, sha256: sha256Hex(exact) });
}

function accepted(id, role, bytes, details = {}) {
  const reference = retain(`${id}.bin`, bytes);
  const decoded = decodeFrame(bytes, role, bindings);
  if (!decoded.ok) throw new Error(`${id}: generator rejected ${decoded.disposition}`);
  cases.push({ id, decision: "accept", role, ...reference, ...details });
}

function rejected(id, role, bytes, disposition) {
  const reference = retain(`${id}.bin`, bytes);
  const decoded = decodeFrame(bytes, role, bindings);
  if (decoded.ok || decoded.disposition !== disposition) {
    throw new Error(`${id}: expected ${disposition}, got ${decoded.disposition ?? "accept"}`);
  }
  cases.push({ id, decision: "reject", role, disposition, ...reference });
}

const ordinarySource = Buffer.from(
  "globalThis.capsuleMain = function (input) { return {doubled: input.value * 2, echo: input.message}; };",
);
const ordinaryInput = Buffer.from(canonicalJSON({ message: "capsule-c2a", value: 21 }));
const ordinaryResult = Buffer.from(canonicalJSON({ doubled: 42, echo: "capsule-c2a" }));

const source = encodeInputFrame(SOURCE_ROLE, ordinarySource, bindings);
const input = encodeInputFrame(INPUT_ROLE, ordinaryInput, bindings);
const completion = encodeCompletionFrame("succeeded", ordinaryResult, bindings);

accepted("accept-source-ordinary", SOURCE_ROLE, source, { payloadBytes: ordinarySource.length });
accepted("accept-input-ordinary", INPUT_ROLE, input, { payloadBytes: ordinaryInput.length });
accepted("accept-completion-succeeded", COMPLETION_ROLE, completion, {
  payloadBytes: ordinaryResult.length,
  status: "succeeded",
  trailerOffset: COMPLETION_HEADER_BYTES + ordinaryResult.length,
});
for (const status of ["workload-failed", "result-invalid", "child-terminated"]) {
  accepted(
    `accept-completion-${status}`,
    COMPLETION_ROLE,
    encodeCompletionFrame(status, Buffer.from("null"), bindings),
    { payloadBytes: 4, status },
  );
}

const maximumSource = Buffer.concat([
  Buffer.from("//"),
  Buffer.alloc(PAYLOAD_MAXIMUM_BYTES - 2, 0x61),
]);
const maximumJSON = Buffer.from(`"${"a".repeat(PAYLOAD_MAXIMUM_BYTES - 2)}"`);
accepted(
  "accept-source-maximum",
  SOURCE_ROLE,
  encodeInputFrame(SOURCE_ROLE, maximumSource, bindings),
  { payloadBytes: PAYLOAD_MAXIMUM_BYTES },
);
accepted("accept-input-maximum", INPUT_ROLE, encodeInputFrame(INPUT_ROLE, maximumJSON, bindings), {
  payloadBytes: PAYLOAD_MAXIMUM_BYTES,
});
const completionMaximum = encodeCompletionFrame("succeeded", maximumJSON, bindings);
accepted("accept-completion-maximum", COMPLETION_ROLE, completionMaximum, {
  payloadBytes: PAYLOAD_MAXIMUM_BYTES,
  trailerOffset: COMPLETION_HEADER_BYTES + PAYLOAD_MAXIMUM_BYTES,
});
rejected(
  "reject-completion-cap-plus-one",
  COMPLETION_ROLE,
  Buffer.concat([completionMaximum, Buffer.of(0)]),
  "FRAME_OVERSIZE",
);
rejected(
  "reject-source-cap-plus-one",
  SOURCE_ROLE,
  Buffer.concat([encodeInputFrame(SOURCE_ROLE, maximumSource, bindings), Buffer.of(0)]),
  "FRAME_OVERSIZE",
);

function mutated(bytes, offset, value = undefined) {
  const copy = Buffer.from(bytes);
  copy[offset] = value === undefined ? copy[offset] ^ 1 : value;
  return copy;
}

rejected("reject-header-truncated", SOURCE_ROLE, source.subarray(0, 151), "HEADER_TRUNCATED");
rejected("reject-magic", SOURCE_ROLE, mutated(source, 0), "MAGIC");
rejected("reject-protocol-version", SOURCE_ROLE, mutated(source, 9), "PROTOCOL_VERSION");
rejected("reject-method", SOURCE_ROLE, mutated(source, 11), "METHOD");
rejected("reject-role", SOURCE_ROLE, mutated(source, 13), "ROLE");
rejected("reject-header-length", SOURCE_ROLE, mutated(source, 15), "HEADER_LENGTH");
rejected("reject-completion-flags", COMPLETION_ROLE, mutated(completion, 115), "FLAGS_RESERVED");
rejected("reject-completion-reserved", COMPLETION_ROLE, mutated(completion, 119), "FLAGS_RESERVED");

const zeroAttempt = Buffer.from(source);
zeroAttempt.fill(0, 16, 32);
rejected("reject-zero-attempt", SOURCE_ROLE, zeroAttempt, "IDENTIFIER");
const equalIdentifiers = Buffer.from(source);
equalIdentifiers.copy(equalIdentifiers, 32, 16, 32);
rejected("reject-equal-identifier-domains", SOURCE_ROLE, equalIdentifiers, "IDENTIFIER");
for (const [name, offset] of [
  ["attempt", 16],
  ["registration", 32],
  ["plan", 48],
  ["runtime-profile", 80],
]) {
  rejected(`reject-binding-${name}`, SOURCE_ROLE, mutated(source, offset), "BINDING");
}

const declaredCapPlusOne = Buffer.from(input);
declaredCapPlusOne.writeBigUInt64BE(BigInt(PAYLOAD_MAXIMUM_BYTES + 1), 112);
rejected(
  "reject-declared-payload-cap-plus-one",
  INPUT_ROLE,
  declaredCapPlusOne,
  "PAYLOAD_LENGTH_CAP",
);
rejected(
  "reject-source-truncated",
  SOURCE_ROLE,
  source.subarray(0, source.length - 1),
  "FRAME_LENGTH",
);
rejected(
  "reject-source-trailing",
  SOURCE_ROLE,
  Buffer.concat([source, Buffer.of(0)]),
  "FRAME_LENGTH",
);
rejected("reject-duplicate-frame", SOURCE_ROLE, Buffer.concat([source, source]), "FRAME_LENGTH");
rejected(
  "reject-payload-digest",
  SOURCE_ROLE,
  mutated(source, INPUT_HEADER_BYTES),
  "PAYLOAD_DIGEST",
);

const unknownStatus = Buffer.from(completion);
unknownStatus.writeUInt16BE(99, 112);
rejected("reject-completion-status", COMPLETION_ROLE, unknownStatus, "STATUS");
const nonSuccessProse = encodeCompletionFrame("succeeded", Buffer.from('"guest prose"'), bindings);
nonSuccessProse.writeUInt16BE(STATUS["workload-failed"], 112);
rejected("reject-non-success-guest-prose", COMPLETION_ROLE, nonSuccessProse, "NON_SUCCESS_PAYLOAD");

for (const [name, json] of [
  ["whitespace", '{"a":1, "b":2}'],
  ["key-order", '{"b":2,"a":1}'],
  ["duplicate-key", '{"a":1,"a":2}'],
  ["float", "1.5"],
  ["unsafe-integer", "9007199254740992"],
  ["minus-zero", "-0"],
  ["second-document", "nullnull"],
]) {
  rejected(
    `reject-input-json-${name}`,
    INPUT_ROLE,
    encodeInputFrame(INPUT_ROLE, Buffer.from(json), bindings),
    "JSON",
  );
}
rejected(
  "reject-input-json-invalid-utf8",
  INPUT_ROLE,
  encodeInputFrame(INPUT_ROLE, Buffer.from([0xff]), bindings),
  "JSON",
);

rejected(
  "reject-completion-missing-trailer",
  COMPLETION_ROLE,
  completion.subarray(0, completion.length - TRAILER_BYTES),
  "COMMIT_OFFSET",
);
const trailerOffset = completion.length - TRAILER_BYTES;
for (const [name, offset, disposition] of [
  ["magic", trailerOffset, "COMMIT_MAGIC"],
  ["version", trailerOffset + 9, "COMMIT_VERSION"],
  ["method", trailerOffset + 11, "COMMIT_METHOD"],
  ["role", trailerOffset + 13, "COMMIT_ROLE"],
  ["length", trailerOffset + 15, "COMMIT_LENGTH"],
  ["attempt", trailerOffset + 16, "COMMIT_ATTEMPT"],
  ["digest", trailerOffset + 32, "COMMIT_DIGEST"],
]) {
  rejected(`reject-commit-${name}`, COMPLETION_ROLE, mutated(completion, offset), disposition);
}

const stateCases = [
  ["cancel-before-g", "BOUND", "REFUSED_CANCELLED_BEFORE_START", false],
  ["cancel-during-source", "INPUT_TRANSFER", "REFUSED_INPUT_CANCELLED", false],
  ["source-short-write-error", "INPUT_TRANSFER", "REFUSED_SOURCE_TRANSPORT_FAULT", false],
  ["input-zero-progress-deadline", "INPUT_TRANSFER", "REFUSED_INPUT_STALL", false],
  ["completion-reader-death", "CHILD_RUNNING", "REFUSED_COMPLETION_READER_FAULT", false],
  ["completion-reset-before-trailer", "RESULT_VALIDATED", "REFUSED_MISSING_COMMIT", false],
  ["valid-frame-before-absence", "FRAME_OBSERVED", "FRAME_ONLY_AWAIT_TERMINAL_PROOF", false],
  ["launcher-death-after-trailer", "FRAME_OBSERVED", "REFUSED_LIFECYCLE", false],
  ["runner-zero-without-frame", "TERMINAL_PROOF", "REFUSED_MISSING_COMMIT", false],
  ["caller-response-loss-after-commit", "COMPLETE", "REPLAY_BYTE_IDENTICAL", true],
  ["store-indeterminate", "DURABLE_COMMIT", "RECOVERY_REQUIRED_FENCE", false],
  ["cancel-concurrent-before-terminal-proof", "TERMINAL_PROOF", "REFUSED_CANCELLED", false],
  ["cancel-after-durable-commit", "COMPLETE", "REPLAY_BYTE_IDENTICAL", true],
].map(([id, atState, disposition, durableCompletion]) => ({
  id,
  atState,
  disposition,
  durableCompletion,
}));

const restorationCases = [
  ["source-input-endpoint-swap", "endpoint-object-identity-and-role-magic"],
  ["completion-source-endpoint-swap", "endpoint-object-identity-and-role-magic"],
  ["descriptor-alias", "dedicated-open-description-manifest"],
  ["descriptor-wrong-mode", "access-mode-manifest"],
  ["descriptor-cloexec-cleared", "pre-post-cloexec-canary"],
  ["descriptor-nonblocking-changed", "pre-post-status-flag-canary"],
  ["descriptor-inherited-by-workload", "closed-child-fd-manifest"],
  ["wrong-attempt", "typed-attempt-binding"],
  ["wrong-registration", "typed-registration-binding"],
  ["plan-digest-substitution", "retained-plan-digest"],
  ["runtime-profile-substitution", "retained-runtime-profile-digest"],
  ["payload-flood-after-valid-prefix", "full-drained-byte-count"],
  ["early-trailer", "calculated-final-trailer-offset"],
  ["eof-as-commit", "commit-trailer-required"],
  ["runner-zero-as-success", "terminal-proof-join"],
  ["diagnostic-console-substitution", "typed-role-magic-and-channel"],
  ["implicit-console-restored", "runner-call-and-device-inventory"],
  ["vsock-restored", "runner-call-and-device-inventory"],
  ["network-device-restored", "runner-call-and-device-inventory"],
  ["virtiofs-restored", "runner-call-and-device-inventory"],
  ["cleanup-unresolved", "terminal-proof-join"],
  ["response-loss-result-substitution", "immutable-attempt-replay"],
  ["durable-commit-corrupt", "store-reopen-full-validation"],
].map(([id, refusingControl]) => ({ id, refusingControl, detected: true }));

const states = [
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
];

const manifest = {
  contract: "capsule.typed-guest-transport",
  version: 1,
  status: "passive-conformance-only",
  method: METHOD,
  protocolVersion: PROTOCOL_VERSION,
  layout: {
    integerEncoding: "unsigned-big-endian",
    inputHeaderBytes: INPUT_HEADER_BYTES,
    completionHeaderBytes: COMPLETION_HEADER_BYTES,
    trailerBytes: TRAILER_BYTES,
    payloadMaximumBytes: PAYLOAD_MAXIMUM_BYTES,
    inputFrameMaximumBytes: INPUT_FRAME_MAXIMUM_BYTES,
    completionFrameMaximumBytes: COMPLETION_FRAME_MAXIMUM_BYTES,
    completionRetainMaximumBytes: COMPLETION_RETAIN_MAXIMUM_BYTES,
    offsets: {
      magic: 0,
      protocolVersion: 8,
      method: 10,
      role: 12,
      headerLength: 14,
      attemptId: 16,
      registrationId: 32,
      planDigest: 48,
      runtimeProfileDigest: 80,
      inputPayloadLength: 112,
      inputPayloadDigest: 120,
      completionStatus: 112,
      completionFlags: 114,
      completionReserved: 116,
      completionPayloadLength: 120,
      completionPayloadDigest: 128,
      trailerAttemptId: 16,
      trailerHeaderPayloadDigest: 32,
    },
  },
  magic: Object.fromEntries([
    ["source", MAGIC[SOURCE_ROLE].toString("ascii")],
    ["input", MAGIC[INPUT_ROLE].toString("ascii")],
    ["completion", MAGIC[COMPLETION_ROLE].toString("ascii")],
    ["trailer", MAGIC.trailer.toString("ascii")],
  ]),
  roles: { source: SOURCE_ROLE, input: INPUT_ROLE, completion: COMPLETION_ROLE },
  statuses: STATUS,
  bindings: Object.fromEntries(
    Object.entries(bindings).map(([key, value]) => [key, value.toString("hex")]),
  ),
  sourceTransfer: {
    path: "main.mjs",
    encoding: "exact-registered-utf8-bytes-no-bom-no-container",
    sourceManifestTransferred: false,
  },
  canonicalJSON: {
    profile: "capsule.canonical-inline-json/v0",
    utf8: true,
    bom: false,
    whitespace: false,
    safeIntegersOnly: true,
    floats: false,
    asciiKeyPattern: "^[A-Za-z0-9][A-Za-z0-9_.-]{0,127}$",
    keyOrder: "unsigned-ascii-bytes",
    duplicateKeys: false,
    normalization: "none",
    stringEscapes: "quote-backslash-and-lowercase-u00xx-controls-only",
  },
  nonSuccessPayload: { exactUtf8: "null", guestProse: false, paths: false, timing: false },
  refusalPrecedence: DISPOSITION_ORDER,
  stateModel: { monotonic: true, states, resetReturnsToWritable: false },
  endpointContract: {
    creator: "Execution Supervisor",
    source: { runnerFd: 5, launcherFd: 3, direction: "supervisor-to-launcher" },
    input: { runnerFd: 6, launcherFd: 4, direction: "supervisor-to-launcher" },
    completion: {
      runnerFd: 7,
      launcherFd: 5,
      direction: "launcher-to-supervisor",
      workloadAccess: false,
    },
    cloexecAtCreation: true,
    distinctOpenDescriptions: true,
    drainBeforeAuthorization: true,
    deadlineMs: 1000,
    cancellationMonotonic: true,
    capPlusOneFullyDrained: true,
  },
  durableProjection: {
    frameIsObservationOnly: true,
    eofIsCommit: false,
    runnerExitIsCommit: false,
    requiresIndependentAuthorityInputRuntimeResultLifecyclePublicationFacts: true,
    committedLast: true,
  },
  knownAnswers: {
    sourcePayload: retain("payload-source-ordinary.bin", ordinarySource),
    inputPayload: retain("payload-input-ordinary.bin", ordinaryInput),
    completionPayload: retain("payload-completion-ordinary.bin", ordinaryResult),
  },
  cases,
  stateCases,
  restorationCases,
  effects: {
    endpoint: false,
    process: false,
    runtime: false,
    backend: false,
    guest: false,
    store: false,
  },
};

files.set("manifest.json", Buffer.from(`${JSON.stringify(manifest, null, 2)}\n`));

await mkdir(outputPath, { recursive: true });
for (const [name, bytes] of files) {
  const url = new URL(name, outputUrl);
  if (check) {
    let actual;
    try {
      actual = await readFile(url);
    } catch {
      throw new Error(`missing generated fixture ${name}`);
    }
    if (!actual.equals(bytes)) throw new Error(`generated fixture drift: ${name}`);
  } else {
    await writeFile(url, bytes);
  }
}

process.stdout.write(
  `${check ? "verified" : "generated"} ${files.size} typed transport fixtures, ${cases.length} frame cases, ${stateCases.length} state cases, and ${restorationCases.length} restoration cases\n`,
);
