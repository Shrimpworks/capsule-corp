import { createHash } from "node:crypto";

export const PROTOCOL_VERSION = 1;
export const METHOD = 1;
export const SOURCE_ROLE = 1;
export const INPUT_ROLE = 2;
export const COMPLETION_ROLE = 3;
export const INPUT_HEADER_BYTES = 152;
export const COMPLETION_HEADER_BYTES = 160;
export const TRAILER_BYTES = 64;
export const PAYLOAD_MAXIMUM_BYTES = 262_144;
export const INPUT_FRAME_MAXIMUM_BYTES = 262_296;
export const COMPLETION_FRAME_MAXIMUM_BYTES = 262_368;
export const COMPLETION_RETAIN_MAXIMUM_BYTES = 262_369;

export const MAGIC = Object.freeze({
  [SOURCE_ROLE]: Buffer.from("CPSRC001", "ascii"),
  [INPUT_ROLE]: Buffer.from("CPINP001", "ascii"),
  [COMPLETION_ROLE]: Buffer.from("CPCMP001", "ascii"),
  trailer: Buffer.from("CPEND001", "ascii"),
});

export const STATUS = Object.freeze({
  succeeded: 1,
  "workload-failed": 2,
  "result-invalid": 3,
  "child-terminated": 4,
});

export const DISPOSITION_ORDER = Object.freeze([
  "FRAME_OVERSIZE",
  "HEADER_TRUNCATED",
  "MAGIC",
  "PROTOCOL_VERSION",
  "METHOD",
  "ROLE",
  "HEADER_LENGTH",
  "FLAGS_RESERVED",
  "IDENTIFIER",
  "BINDING",
  "PAYLOAD_LENGTH_CAP",
  "FRAME_LENGTH",
  "COMMIT_OFFSET",
  "PAYLOAD_DIGEST",
  "STATUS",
  "NON_SUCCESS_PAYLOAD",
  "JSON",
  "COMMIT_MAGIC",
  "COMMIT_VERSION",
  "COMMIT_METHOD",
  "COMMIT_ROLE",
  "COMMIT_LENGTH",
  "COMMIT_ATTEMPT",
  "COMMIT_DIGEST",
]);

const textDecoder = new TextDecoder("utf-8", { fatal: true });
const inlineKey = /^[A-Za-z0-9][A-Za-z0-9_.-]{0,127}$/u;

export function sha256(bytes) {
  return createHash("sha256").update(bytes).digest();
}

export function sha256Hex(bytes) {
  return sha256(bytes).toString("hex");
}

export function encodeInputFrame(role, payload, bindings) {
  if (role !== SOURCE_ROLE && role !== INPUT_ROLE) throw new Error("input role");
  const exact = Buffer.from(payload);
  if (exact.length > PAYLOAD_MAXIMUM_BYTES) throw new Error("payload cap");
  const header = Buffer.alloc(INPUT_HEADER_BYTES);
  MAGIC[role].copy(header, 0);
  header.writeUInt16BE(PROTOCOL_VERSION, 8);
  header.writeUInt16BE(METHOD, 10);
  header.writeUInt16BE(role, 12);
  header.writeUInt16BE(INPUT_HEADER_BYTES, 14);
  writeBindings(header, bindings);
  header.writeBigUInt64BE(BigInt(exact.length), 112);
  sha256(exact).copy(header, 120);
  return Buffer.concat([header, exact]);
}

export function encodeCompletionFrame(statusName, payload, bindings) {
  const status = STATUS[statusName];
  if (status === undefined) throw new Error("status");
  const exact = Buffer.from(payload);
  if (exact.length > PAYLOAD_MAXIMUM_BYTES) throw new Error("payload cap");
  if (status !== STATUS.succeeded && !exact.equals(Buffer.from("null"))) {
    throw new Error("non-success payload");
  }
  const header = Buffer.alloc(COMPLETION_HEADER_BYTES);
  MAGIC[COMPLETION_ROLE].copy(header, 0);
  header.writeUInt16BE(PROTOCOL_VERSION, 8);
  header.writeUInt16BE(METHOD, 10);
  header.writeUInt16BE(COMPLETION_ROLE, 12);
  header.writeUInt16BE(COMPLETION_HEADER_BYTES, 14);
  writeBindings(header, bindings);
  header.writeUInt16BE(status, 112);
  header.writeUInt16BE(0, 114);
  header.writeUInt32BE(0, 116);
  header.writeBigUInt64BE(BigInt(exact.length), 120);
  sha256(exact).copy(header, 128);
  const committed = Buffer.concat([header, exact]);
  const trailer = Buffer.alloc(TRAILER_BYTES);
  MAGIC.trailer.copy(trailer, 0);
  trailer.writeUInt16BE(PROTOCOL_VERSION, 8);
  trailer.writeUInt16BE(METHOD, 10);
  trailer.writeUInt16BE(COMPLETION_ROLE, 12);
  trailer.writeUInt16BE(TRAILER_BYTES, 14);
  Buffer.from(bindings.attemptId).copy(trailer, 16);
  sha256(committed).copy(trailer, 32);
  return Buffer.concat([committed, trailer]);
}

function writeBindings(header, bindings) {
  const fields = [bindings.attemptId, bindings.registrationId];
  for (const value of fields) {
    if (Buffer.from(value).length !== 16) throw new Error("identifier width");
  }
  for (const value of [bindings.planDigest, bindings.runtimeProfileDigest]) {
    if (Buffer.from(value).length !== 32) throw new Error("digest width");
  }
  Buffer.from(bindings.attemptId).copy(header, 16);
  Buffer.from(bindings.registrationId).copy(header, 32);
  Buffer.from(bindings.planDigest).copy(header, 48);
  Buffer.from(bindings.runtimeProfileDigest).copy(header, 80);
}

export function decodeFrame(bytes, expectedRole, bindings) {
  const exact = Buffer.from(bytes);
  const maximum =
    expectedRole === COMPLETION_ROLE ? COMPLETION_FRAME_MAXIMUM_BYTES : INPUT_FRAME_MAXIMUM_BYTES;
  if (exact.length > maximum) return refused("FRAME_OVERSIZE");
  const headerBytes =
    expectedRole === COMPLETION_ROLE ? COMPLETION_HEADER_BYTES : INPUT_HEADER_BYTES;
  if (exact.length < headerBytes) return refused("HEADER_TRUNCATED");
  if (!exact.subarray(0, 8).equals(MAGIC[expectedRole])) return refused("MAGIC");
  if (exact.readUInt16BE(8) !== PROTOCOL_VERSION) return refused("PROTOCOL_VERSION");
  if (exact.readUInt16BE(10) !== METHOD) return refused("METHOD");
  if (exact.readUInt16BE(12) !== expectedRole) return refused("ROLE");
  if (exact.readUInt16BE(14) !== headerBytes) return refused("HEADER_LENGTH");
  if (
    expectedRole === COMPLETION_ROLE &&
    (exact.readUInt16BE(114) !== 0 || exact.readUInt32BE(116) !== 0)
  ) {
    return refused("FLAGS_RESERVED");
  }
  const attempt = exact.subarray(16, 32);
  const registration = exact.subarray(32, 48);
  if (isZero(attempt) || isZero(registration) || attempt.equals(registration)) {
    return refused("IDENTIFIER");
  }
  if (
    !attempt.equals(bindings.attemptId) ||
    !registration.equals(bindings.registrationId) ||
    !exact.subarray(48, 80).equals(bindings.planDigest) ||
    !exact.subarray(80, 112).equals(bindings.runtimeProfileDigest)
  ) {
    return refused("BINDING");
  }
  const lengthOffset = expectedRole === COMPLETION_ROLE ? 120 : 112;
  const digestOffset = expectedRole === COMPLETION_ROLE ? 128 : 120;
  const declared = exact.readBigUInt64BE(lengthOffset);
  if (declared > BigInt(PAYLOAD_MAXIMUM_BYTES)) return refused("PAYLOAD_LENGTH_CAP");
  if (
    bindings.expectedPayloadBytes !== undefined &&
    declared !== BigInt(bindings.expectedPayloadBytes)
  ) {
    return refused("BINDING");
  }
  const payloadLength = Number(declared);
  const payloadEnd = headerBytes + payloadLength;
  if (exact.length < payloadEnd) return refused("FRAME_LENGTH");
  const expectedLength = payloadEnd + (expectedRole === COMPLETION_ROLE ? TRAILER_BYTES : 0);
  if (expectedRole === COMPLETION_ROLE && exact.length < expectedLength) {
    return refused("COMMIT_OFFSET");
  }
  if (exact.length !== expectedLength) return refused("FRAME_LENGTH");
  const payload = exact.subarray(headerBytes, headerBytes + payloadLength);
  if (!sha256(payload).equals(exact.subarray(digestOffset, digestOffset + 32))) {
    return refused("PAYLOAD_DIGEST");
  }
  if (
    bindings.expectedPayloadDigest !== undefined &&
    !sha256(payload).equals(Buffer.from(bindings.expectedPayloadDigest))
  ) {
    return refused("BINDING");
  }
  let statusName = null;
  if (expectedRole === COMPLETION_ROLE) {
    const status = exact.readUInt16BE(112);
    statusName = Object.keys(STATUS).find((name) => STATUS[name] === status);
    if (statusName === undefined) return refused("STATUS");
    if (statusName !== "succeeded" && !payload.equals(Buffer.from("null"))) {
      return refused("NON_SUCCESS_PAYLOAD");
    }
  }
  if (expectedRole !== SOURCE_ROLE) {
    try {
      parseCanonicalJSON(payload);
    } catch {
      return refused("JSON");
    }
  } else if (
    payload.length === 0 ||
    payload.subarray(0, 3).equals(Buffer.from([0xef, 0xbb, 0xbf]))
  ) {
    return refused("JSON");
  } else {
    try {
      textDecoder.decode(payload);
    } catch {
      return refused("JSON");
    }
  }
  if (expectedRole === COMPLETION_ROLE) {
    const trailerOffset = payloadEnd;
    if (trailerOffset + TRAILER_BYTES !== exact.length) return refused("COMMIT_OFFSET");
    const trailer = exact.subarray(trailerOffset);
    if (!trailer.subarray(0, 8).equals(MAGIC.trailer)) return refused("COMMIT_MAGIC");
    if (trailer.readUInt16BE(8) !== PROTOCOL_VERSION) return refused("COMMIT_VERSION");
    if (trailer.readUInt16BE(10) !== METHOD) return refused("COMMIT_METHOD");
    if (trailer.readUInt16BE(12) !== COMPLETION_ROLE) return refused("COMMIT_ROLE");
    if (trailer.readUInt16BE(14) !== TRAILER_BYTES) return refused("COMMIT_LENGTH");
    if (!trailer.subarray(16, 32).equals(attempt)) return refused("COMMIT_ATTEMPT");
    if (!trailer.subarray(32, 64).equals(sha256(exact.subarray(0, trailerOffset)))) {
      return refused("COMMIT_DIGEST");
    }
  }
  return Object.freeze({ ok: true, payload: Buffer.from(payload), status: statusName });
}

function refused(disposition) {
  return Object.freeze({ ok: false, disposition });
}

function isZero(bytes) {
  return bytes.every((value) => value === 0);
}

export function canonicalJSON(value) {
  if (value === null) return "null";
  if (value === true) return "true";
  if (value === false) return "false";
  if (typeof value === "number") {
    if (!Number.isSafeInteger(value) || Object.is(value, -0)) throw new Error("integer");
    return String(value);
  }
  if (typeof value === "string") return canonicalString(value);
  if (Array.isArray(value)) return `[${value.map(canonicalJSON).join(",")}]`;
  if (typeof value === "object") {
    const entries = Object.entries(value).sort(([left], [right]) =>
      Buffer.compare(Buffer.from(left), Buffer.from(right)),
    );
    for (const [key] of entries) if (!inlineKey.test(key)) throw new Error("key");
    return `{${entries.map(([key, child]) => `${canonicalString(key)}:${canonicalJSON(child)}`).join(",")}}`;
  }
  throw new Error("type");
}

function canonicalString(value) {
  let output = '"';
  for (const scalar of value) {
    const point = scalar.codePointAt(0);
    if (point >= 0xd800 && point <= 0xdfff) throw new Error("surrogate");
    if (scalar === '"') output += '\\"';
    else if (scalar === "\\") output += "\\\\";
    else if (point <= 0x1f) output += `\\u00${point.toString(16).padStart(2, "0")}`;
    else output += scalar;
  }
  return `${output}"`;
}

export function parseCanonicalJSON(bytes) {
  const exact = Buffer.from(bytes);
  if (exact.length === 0 || exact.subarray(0, 3).equals(Buffer.from([0xef, 0xbb, 0xbf]))) {
    throw new Error("empty-or-bom");
  }
  const text = textDecoder.decode(exact);
  rejectDuplicateKeys(text);
  const value = JSON.parse(text);
  if (canonicalJSON(value) !== text) throw new Error("noncanonical");
  return value;
}

function rejectDuplicateKeys(text) {
  let index = 0;
  const whitespace = () => {
    while (/\s/u.test(text[index] ?? "")) index += 1;
  };
  const string = () => {
    const start = index;
    index += 1;
    while (index < text.length) {
      if (text[index] === "\\") index += 2;
      else if (text[index++] === '"') return JSON.parse(text.slice(start, index));
    }
    throw new Error("string");
  };
  const value = () => {
    whitespace();
    if (text[index] === "{") return object();
    if (text[index] === "[") return array();
    if (text[index] === '"') return string();
    const match = /^(?:null|true|false|-?(?:0|[1-9][0-9]*)(?:\.[0-9]+)?(?:[eE][+-]?[0-9]+)?)/u.exec(
      text.slice(index),
    );
    if (!match) throw new Error("value");
    index += match[0].length;
  };
  const object = () => {
    index += 1;
    const keys = new Set();
    whitespace();
    if (text[index] === "}") {
      index += 1;
      return;
    }
    while (true) {
      whitespace();
      if (text[index] !== '"') throw new Error("key");
      const key = string();
      if (keys.has(key)) throw new Error("duplicate");
      keys.add(key);
      whitespace();
      if (text[index++] !== ":") throw new Error("colon");
      value();
      whitespace();
      const separator = text[index++];
      if (separator === "}") return;
      if (separator !== ",") throw new Error("separator");
    }
  };
  const array = () => {
    index += 1;
    whitespace();
    if (text[index] === "]") {
      index += 1;
      return;
    }
    while (true) {
      value();
      whitespace();
      const separator = text[index++];
      if (separator === "]") return;
      if (separator !== ",") throw new Error("separator");
    }
  };
  value();
  whitespace();
  if (index !== text.length) throw new Error("trailing");
}
