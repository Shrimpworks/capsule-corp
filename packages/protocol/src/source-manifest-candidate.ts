import { createHash } from "node:crypto";

import {
  concatenateBytes,
  encodeCborArrayHeader,
  encodeCborByteString,
  encodeCborMapHeader,
  encodeCborText,
  encodeCborUnsigned,
} from "./internal-cbor-primitives.js";
import {
  MJS_MAIN_PATH,
  type MjsMainEntrypoint,
  type MjsMainLogicalPath,
  type RetainedMjsSourceBytes,
  SOURCE_MANIFEST_MEDIA_TYPE,
  type SourceManifestMediaType,
  validateMjsSourceBytes,
  validateSourceManifestMediaType,
} from "./mjs-source-candidate.js";

export const SOURCE_MANIFEST_OBJECT_TYPE = "capsule.source-manifest" as const;
export const SOURCE_MANIFEST_OBJECT_VERSION = 0 as const;
export const SOURCE_MANIFEST_MIN_CBOR_BYTES = 87 as const;
export const SOURCE_MANIFEST_MAX_CBOR_BYTES = 95 as const;

declare const sourceManifestScalarBrand: unique symbol;

export type MjsSourceContentDigest = string & {
  readonly [sourceManifestScalarBrand]: "sha256:mjs-source-content:32";
};
export type SourceManifestDigest = string & {
  readonly [sourceManifestScalarBrand]: "sha256:source-manifest:32";
};

export interface SourceManifestMemberV0 {
  readonly logicalPath: MjsMainLogicalPath;
  readonly contentDigest: MjsSourceContentDigest;
  readonly byteLength: number;
}

export interface SourceManifestV0 {
  readonly objectType: typeof SOURCE_MANIFEST_OBJECT_TYPE;
  readonly objectVersion: typeof SOURCE_MANIFEST_OBJECT_VERSION;
  readonly entrypoint: MjsMainEntrypoint;
  readonly members: readonly [SourceManifestMemberV0];
  readonly aggregateByteLength: number;
}

export interface DecodedSourceManifestV0 {
  readonly mediaType: SourceManifestMediaType;
  readonly byteLength: number;
  readonly digest: SourceManifestDigest;
  readonly view: SourceManifestV0;
  copyAuthoritativeBytes(): Uint8Array;
  copyDigestBytes(): Uint8Array;
}

export type SourceManifestRefusalCode =
  | "MANIFEST_MEDIA_TYPE"
  | "MANIFEST_CBOR"
  | "MANIFEST_OBJECT_TYPE"
  | "MANIFEST_OBJECT_VERSION"
  | "MANIFEST_SHAPE"
  | "MANIFEST_PATH"
  | "MANIFEST_SOURCE_DIGEST"
  | "MANIFEST_SOURCE_LENGTH";

export interface SourceManifestRefusal {
  readonly owner: "source-manifest-validator";
  readonly classification: "MALFORMED" | "UNSUPPORTED" | "SCHEMA" | "DOMAIN";
  readonly code: SourceManifestRefusalCode;
}

export type SourceManifestDecodeResult =
  | { readonly ok: true; readonly manifest: DecodedSourceManifestV0 }
  | { readonly ok: false; readonly refusal: SourceManifestRefusal };

class ManifestDecodeFailure extends Error {
  constructor(
    readonly classification: SourceManifestRefusal["classification"],
    readonly code: SourceManifestRefusalCode,
  ) {
    super(code);
  }
}

/** Encodes the one canonical SourceManifest from already validated retained source bytes. */
export function encodeSourceManifestV0(source: RetainedMjsSourceBytes): Uint8Array {
  const sourceBytes = source.copyExactBytes();
  const contentDigest = createHash("sha256").update(sourceBytes).digest();
  return concatenateBytes([
    encodeCborMapHeader(5),
    encodeCborUnsigned(1),
    encodeCborText(SOURCE_MANIFEST_OBJECT_TYPE),
    encodeCborUnsigned(2),
    encodeCborUnsigned(SOURCE_MANIFEST_OBJECT_VERSION),
    encodeCborUnsigned(3),
    encodeCborText(MJS_MAIN_PATH),
    encodeCborUnsigned(4),
    encodeCborArrayHeader(1),
    encodeCborArrayHeader(3),
    encodeCborText(MJS_MAIN_PATH),
    encodeCborByteString(contentDigest),
    encodeCborUnsigned(sourceBytes.byteLength),
    encodeCborUnsigned(5),
    encodeCborUnsigned(sourceBytes.byteLength),
  ]);
}

/**
 * Strictly decodes, hashes, role-binds, and defensively retains one passive
 * SourceManifest and its separately supplied exact source bytes.
 */
export function decodeSourceManifestV0(
  received: Uint8Array,
  sourceBytes: Uint8Array,
  mediaType: string,
): SourceManifestDecodeResult {
  if (!validateSourceManifestMediaType(mediaType)) {
    return rejected("UNSUPPORTED", "MANIFEST_MEDIA_TYPE");
  }
  const source = validateMjsSourceBytes(sourceBytes);
  if (!source.ok) {
    return rejected("DOMAIN", "MANIFEST_SOURCE_DIGEST");
  }
  const authoritative = Uint8Array.from(received);
  if (
    authoritative.byteLength < SOURCE_MANIFEST_MIN_CBOR_BYTES ||
    authoritative.byteLength > SOURCE_MANIFEST_MAX_CBOR_BYTES
  ) {
    return rejected("MALFORMED", "MANIFEST_CBOR");
  }

  try {
    const reader = new CborReader(authoritative);
    if (reader.containerLength(5) !== 5) schemaFail("MANIFEST_SHAPE");
    reader.requiredUnsigned(1);
    const objectType = reader.text();
    if (objectType !== SOURCE_MANIFEST_OBJECT_TYPE) {
      schemaFail("MANIFEST_OBJECT_TYPE", "UNSUPPORTED");
    }
    reader.requiredUnsigned(2);
    const objectVersion = reader.unsigned();
    if (objectVersion !== SOURCE_MANIFEST_OBJECT_VERSION) {
      schemaFail("MANIFEST_OBJECT_VERSION", "UNSUPPORTED");
    }
    reader.requiredUnsigned(3);
    if (reader.text() !== MJS_MAIN_PATH) schemaFail("MANIFEST_PATH", "DOMAIN");
    reader.requiredUnsigned(4);
    if (reader.containerLength(4) !== 1) schemaFail("MANIFEST_SHAPE");
    if (reader.containerLength(4) !== 3) schemaFail("MANIFEST_SHAPE");
    if (reader.text() !== MJS_MAIN_PATH) schemaFail("MANIFEST_PATH", "DOMAIN");
    const contentDigestBytes = reader.bytes();
    if (contentDigestBytes.byteLength !== 32) schemaFail("MANIFEST_SHAPE");
    const memberByteLength = reader.unsigned();
    reader.requiredUnsigned(5);
    const aggregateByteLength = reader.unsigned();
    reader.finish();

    const exactSource = source.source.copyExactBytes();
    const expectedContentDigest = createHash("sha256").update(exactSource).digest();
    if (!Buffer.from(contentDigestBytes).equals(expectedContentDigest)) {
      schemaFail("MANIFEST_SOURCE_DIGEST", "DOMAIN");
    }
    if (
      memberByteLength !== exactSource.byteLength ||
      aggregateByteLength !== exactSource.byteLength
    ) {
      schemaFail("MANIFEST_SOURCE_LENGTH", "DOMAIN");
    }

    const contentDigest = Buffer.from(contentDigestBytes).toString("hex") as MjsSourceContentDigest;
    const digestBytes = createHash("sha256").update(authoritative).digest();
    const digest = digestBytes.toString("hex") as SourceManifestDigest;
    const view = Object.freeze({
      objectType: SOURCE_MANIFEST_OBJECT_TYPE,
      objectVersion: SOURCE_MANIFEST_OBJECT_VERSION,
      entrypoint: MJS_MAIN_PATH as MjsMainEntrypoint,
      members: Object.freeze([
        Object.freeze({
          logicalPath: MJS_MAIN_PATH as MjsMainLogicalPath,
          contentDigest,
          byteLength: memberByteLength,
        }),
      ]) as readonly [SourceManifestMemberV0],
      aggregateByteLength,
    });
    const manifest = Object.freeze({
      mediaType: SOURCE_MANIFEST_MEDIA_TYPE as SourceManifestMediaType,
      byteLength: authoritative.byteLength,
      digest,
      view,
      copyAuthoritativeBytes: () => Uint8Array.from(authoritative),
      copyDigestBytes: () => Uint8Array.from(digestBytes),
    }) as DecodedSourceManifestV0;
    return Object.freeze({ ok: true, manifest });
  } catch (error) {
    if (error instanceof ManifestDecodeFailure) {
      return rejected(error.classification, error.code);
    }
    throw error;
  }
}

class CborReader {
  #offset = 0;

  constructor(private readonly value: Uint8Array) {}

  finish(): void {
    if (this.#offset !== this.value.byteLength) schemaFail("MANIFEST_CBOR", "MALFORMED");
  }

  requiredUnsigned(expected: number): void {
    if (this.unsigned() !== expected) schemaFail("MANIFEST_CBOR", "MALFORMED");
  }

  unsigned(): number {
    const [major, value] = this.head();
    if (major !== 0) schemaFail("MANIFEST_CBOR", "MALFORMED");
    return value;
  }

  containerLength(expectedMajor: 4 | 5): number {
    const [major, value] = this.head();
    if (major !== expectedMajor) schemaFail("MANIFEST_CBOR", "MALFORMED");
    return value;
  }

  bytes(): Uint8Array {
    return this.string(2);
  }

  text(): string {
    const bytes = this.string(3);
    try {
      return new TextDecoder("utf-8", { fatal: true }).decode(bytes);
    } catch {
      schemaFail("MANIFEST_CBOR", "MALFORMED");
    }
  }

  private string(expectedMajor: 2 | 3): Uint8Array {
    const [major, length] = this.head();
    if (major !== expectedMajor || length > this.value.byteLength - this.#offset) {
      schemaFail("MANIFEST_CBOR", "MALFORMED");
    }
    const result = this.value.slice(this.#offset, this.#offset + length);
    this.#offset += length;
    return result;
  }

  private head(): readonly [number, number] {
    if (this.#offset >= this.value.byteLength) schemaFail("MANIFEST_CBOR", "MALFORMED");
    const initial = this.value[this.#offset++] as number;
    const major = initial >>> 5;
    const additional = initial & 0x1f;
    if (additional < 24) return [major, additional];
    const width = additional === 24 ? 1 : additional === 25 ? 2 : additional === 26 ? 4 : -1;
    if (width < 0 || width > this.value.byteLength - this.#offset) {
      schemaFail("MANIFEST_CBOR", "MALFORMED");
    }
    let result = 0;
    for (let index = 0; index < width; index += 1) {
      result = result * 256 + (this.value[this.#offset++] as number);
    }
    const minimum = width === 1 ? 24 : width === 2 ? 0x100 : 0x1_0000;
    if (result < minimum) schemaFail("MANIFEST_CBOR", "MALFORMED");
    return [major, result];
  }
}

function schemaFail(
  code: SourceManifestRefusalCode,
  classification: SourceManifestRefusal["classification"] = "SCHEMA",
): never {
  throw new ManifestDecodeFailure(classification, code);
}

function rejected(
  classification: SourceManifestRefusal["classification"],
  code: SourceManifestRefusalCode,
): { readonly ok: false; readonly refusal: SourceManifestRefusal } {
  return Object.freeze({
    ok: false,
    refusal: Object.freeze({ owner: "source-manifest-validator", classification, code }),
  });
}
