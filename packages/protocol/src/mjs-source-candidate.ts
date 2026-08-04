export const JOB_PROPOSAL_MEDIA_TYPE = "application/capsule.job-proposal+json;v=0" as const;
export const MJS_SOURCE_PROFILE = "capsule.mjs-source/v0" as const;
export const MJS_SOURCE_MEDIA_TYPE =
  "application/capsule.javascript-source;v=0;module=esm" as const;
export const SOURCE_MANIFEST_MEDIA_TYPE = "application/capsule.source-manifest+cbor;v=0" as const;
export const MJS_MAIN_PATH = "main.mjs" as const;
export const MJS_SOURCE_MAX_BYTES = 262_144 as const;

declare const retainedMjsSourceBytesBrand: unique symbol;
declare const mjsSourceScalarBrand: unique symbol;

export type JobProposalMediaType = typeof JOB_PROPOSAL_MEDIA_TYPE & {
  readonly [mjsSourceScalarBrand]: "media:job-proposal:v0";
};
export type MjsSourceProfileIdentifier = typeof MJS_SOURCE_PROFILE & {
  readonly [mjsSourceScalarBrand]: "profile:mjs-source:v0";
};
export type MjsSourceMemberMediaType = typeof MJS_SOURCE_MEDIA_TYPE & {
  readonly [mjsSourceScalarBrand]: "media:mjs-source:v0";
};
export type SourceManifestMediaType = typeof SOURCE_MANIFEST_MEDIA_TYPE & {
  readonly [mjsSourceScalarBrand]: "media:source-manifest:v0";
};
export type MjsMainLogicalPath = typeof MJS_MAIN_PATH & {
  readonly [mjsSourceScalarBrand]: "path:mjs-source-member:v0";
};
export type MjsMainEntrypoint = typeof MJS_MAIN_PATH & {
  readonly [mjsSourceScalarBrand]: "entrypoint:mjs-source:v0";
};

export type MjsSourceByteRefusalCode = "SOURCE_BYTES" | "SOURCE_UTF8" | "SOURCE_BOM";

export interface MjsSourceByteRefusal {
  readonly owner: "mjs-source-byte-validator";
  readonly classification: "MALFORMED" | "SEMANTIC";
  readonly code: MjsSourceByteRefusalCode;
}

/**
 * Exact passive byte identity for the one `main.mjs` member.
 *
 * This type proves only profile/media/path, size, strict UTF-8, and leading-BOM
 * rules. It deliberately does not claim ECMAScript parsing or module-request
 * validation; ADR-0034 plan construction remains blocked until that separate
 * boundary is selected and implemented.
 */
export interface RetainedMjsSourceBytes {
  readonly profile: MjsSourceProfileIdentifier;
  readonly memberMediaType: MjsSourceMemberMediaType;
  readonly logicalPath: MjsMainLogicalPath;
  readonly entrypoint: MjsMainEntrypoint;
  readonly byteLength: number;
  readonly [retainedMjsSourceBytesBrand]: "RetainedMjsSourceBytes";
  copyExactBytes(): Uint8Array;
}

export type MjsSourceByteValidationResult =
  | { readonly ok: true; readonly source: RetainedMjsSourceBytes }
  | { readonly ok: false; readonly refusal: MjsSourceByteRefusal };

/** Copies and validates exact source bytes without parsing or executing them. */
export function validateMjsSourceBytes(received: Uint8Array): MjsSourceByteValidationResult {
  if (!(received instanceof Uint8Array) || received.byteLength > MJS_SOURCE_MAX_BYTES) {
    return rejected("SEMANTIC", "SOURCE_BYTES");
  }
  const exactBytes = Uint8Array.from(received);
  if (exactBytes[0] === 0xef && exactBytes[1] === 0xbb && exactBytes[2] === 0xbf) {
    return rejected("MALFORMED", "SOURCE_BOM");
  }
  try {
    new TextDecoder("utf-8", { fatal: true, ignoreBOM: true }).decode(exactBytes);
  } catch {
    return rejected("MALFORMED", "SOURCE_UTF8");
  }

  const retained = Object.freeze({
    profile: MJS_SOURCE_PROFILE as MjsSourceProfileIdentifier,
    memberMediaType: MJS_SOURCE_MEDIA_TYPE as MjsSourceMemberMediaType,
    logicalPath: MJS_MAIN_PATH as MjsMainLogicalPath,
    entrypoint: MJS_MAIN_PATH as MjsMainEntrypoint,
    byteLength: exactBytes.byteLength,
    copyExactBytes: () => Uint8Array.from(exactBytes),
  }) as RetainedMjsSourceBytes;
  return Object.freeze({ ok: true, source: retained });
}

export function validateMjsSourceProfile(value: string): value is MjsSourceProfileIdentifier {
  return value === MJS_SOURCE_PROFILE;
}

export function validateMjsSourceMediaType(value: string): value is MjsSourceMemberMediaType {
  return value === MJS_SOURCE_MEDIA_TYPE;
}

export function validateSourceManifestMediaType(value: string): value is SourceManifestMediaType {
  return value === SOURCE_MANIFEST_MEDIA_TYPE;
}

function rejected(
  classification: MjsSourceByteRefusal["classification"],
  code: MjsSourceByteRefusalCode,
): { readonly ok: false; readonly refusal: MjsSourceByteRefusal } {
  return Object.freeze({
    ok: false,
    refusal: Object.freeze({ owner: "mjs-source-byte-validator", classification, code }),
  });
}
