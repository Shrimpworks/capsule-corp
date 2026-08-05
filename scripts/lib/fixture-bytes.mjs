// Shared byte/digest helpers for the conformance fixture generate/verify
// scripts under scripts/. These were previously copy-pasted verbatim (or
// near-verbatim) into a dozen-plus scripts; centralized here so there is one
// place to get sha256/hex/repeated-byte helpers right.
import { createHash } from "node:crypto";

/** SHA-256 digest of `value`, as raw bytes. */
export function sha256(value) {
  return createHash("sha256").update(value).digest();
}

/** SHA-256 digest of `value`, as a lowercase hex string. */
export function sha256Hex(value) {
  return createHash("sha256").update(value).digest("hex");
}

/** Encode bytes as a lowercase hex string. */
export function toHex(bytes) {
  return Buffer.from(bytes).toString("hex");
}

/** Decode a hex string to a Buffer. */
export function fromHex(value) {
  return Buffer.from(value, "hex");
}

/** `length` bytes, each set to `value`. */
export function repeatedBytes(value, length) {
  return Uint8Array.from({ length }, () => value);
}

/** `length` bytes counting up from `start`, wrapping at 0xff. */
export function sequence(length, start) {
  return Buffer.from(Array.from({ length }, (_, i) => (start + i) & 0xff));
}
