/**
 * Minimal deterministic-CBOR head/scalar encoders shared by the source-manifest
 * encoder (job-proposal-resolver.ts) and the ExecutionPlan encoder
 * (execution-plan-builder.ts). Not part of the public package API — neither
 * consumer's callers see CBOR bytes directly, only digests and opaque byte
 * buffers.
 *
 * Every length/count/scalar header supports the full safe-integer range
 * (0 through 2^53-1), matching CandidateUInt53/PositiveSafeInteger elsewhere
 * in this package. There is no reason for one caller's manifest encoding to
 * silently accept a narrower range than the other's plan encoding.
 */

export function encodeCborHead(majorType: number, value: number): Uint8Array {
  if (!Number.isSafeInteger(value) || value < 0) {
    throw new TypeError("CBOR head argument must be an unsigned safe integer");
  }
  const prefix = majorType << 5;
  if (value < 24) {
    return Uint8Array.of(prefix | value);
  }
  if (value <= 0xff) {
    return Uint8Array.of(prefix | 24, value);
  }
  if (value <= 0xffff) {
    return Uint8Array.of(prefix | 25, value >>> 8, value & 0xff);
  }
  if (value <= 0xffff_ffff) {
    const bytes = new Uint8Array(5);
    bytes[0] = prefix | 26;
    new DataView(bytes.buffer).setUint32(1, value, false);
    return bytes;
  }
  const bytes = new Uint8Array(9);
  bytes[0] = prefix | 27;
  new DataView(bytes.buffer).setBigUint64(1, BigInt(value), false);
  return bytes;
}

export function encodeCborUnsigned(value: number): Uint8Array {
  return encodeCborHead(0, value);
}

export function encodeCborByteString(value: Uint8Array | readonly number[]): Uint8Array {
  const bytes = Uint8Array.from(value);
  return concatenateBytes([encodeCborHead(2, bytes.byteLength), bytes]);
}

export function encodeCborText(value: string): Uint8Array {
  const bytes = new TextEncoder().encode(value);
  return concatenateBytes([encodeCborHead(3, bytes.byteLength), bytes]);
}

export function encodeCborArrayHeader(length: number): Uint8Array {
  return encodeCborHead(4, length);
}

export function encodeCborMapHeader(length: number): Uint8Array {
  return encodeCborHead(5, length);
}

export function concatenateBytes(chunks: readonly Uint8Array[]): Uint8Array {
  const byteLength = chunks.reduce((total, chunk) => total + chunk.byteLength, 0);
  const bytes = new Uint8Array(byteLength);
  let offset = 0;
  for (const chunk of chunks) {
    bytes.set(chunk, offset);
    offset += chunk.byteLength;
  }
  return bytes;
}
