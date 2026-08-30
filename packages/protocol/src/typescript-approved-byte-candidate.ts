import { createHash } from "node:crypto";

import {
  concatenateBytes,
  encodeCborArrayHeader,
  encodeCborByteString,
  encodeCborMapHeader,
  encodeCborText,
  encodeCborUnsigned,
} from "./internal-cbor-primitives.js";

/**
 * Passive, unwired Slice A views for ADR-0026. Nothing in the product calls
 * these helpers and they do not transform, register, approve, or execute code.
 */

export const APPROVED_BYTE_CAPS = Object.freeze({
  sourceFiles: 32,
  originalFileBytes: 262_144,
  originalAggregateBytes: 1_048_576,
  emittedFileBytes: 262_144,
  emittedAggregateBytes: 1_048_576,
});

export const TYPESCRIPT_SOURCE_MEDIA_TYPE =
  "application/capsule.typescript-source;v=0;module=esm" as const;
export const JAVASCRIPT_SOURCE_MEDIA_TYPE =
  "application/capsule.javascript-source;v=0;module=esm" as const;

export const APPROVED_BYTE_TOOLCHAIN = Object.freeze({
  transformer: "node:module.stripTypeScriptTypes",
  nodeVersion: "22.22.1",
  amaroVersion: "1.1.5",
  platform: "darwin",
  architecture: "arm64",
  sourceArchiveSha256: "87104b07e7acee748bcc5391e1bc69cf3571caa0fdfb8b1d6b5fd3f9599b7849",
  distributionArchiveSha256: "261da057fb25ff2912dd6abb7842fc915ddf7947a2cb3c8cce90875d2b9bb667",
  executableSha256: "245e0321af97d3c21dd4e7104457334dfe3c3ba7982d0db75363e354565f8cbb",
});

export const APPROVED_BYTE_OPTIONS = Object.freeze({
  mode: "strip",
  diagnosticPolicy: "reject-any",
  sourceMap: "absent",
  sourceUrl: "absent",
  inputMediaType: TYPESCRIPT_SOURCE_MEDIA_TYPE,
  outputMediaType: JAVASCRIPT_SOURCE_MEDIA_TYPE,
});

/** The fixed set of byte artifacts {@link ApprovedByteDigest} can be bound to. */
export type ApprovedByteDigestRole =
  | "emitted-javascript"
  | "executable-javascript-source-manifest"
  | "normalized-options"
  | "original-authoring-source-manifest"
  | "original-source"
  | "transformation-record-set"
  | "transformer-profile";

declare const approvedByteDigestBrand: unique symbol;
/** A SHA-256 digest bound to exactly one {@link ApprovedByteDigestRole}, producible only via {@link asApprovedByteDigest}. */
export type ApprovedByteDigest<Role extends ApprovedByteDigestRole> = readonly number[] & {
  readonly [approvedByteDigestBrand]: `sha256:${Role}:32`;
};

/** One source file's already-emitted bytes, supplied by the caller — never produced by a transformer call. */
export type ApprovedByteFixtureSource = Readonly<{
  logicalPath: string;
  originalBytes: Uint8Array;
  emittedBytes: Uint8Array;
}>;

/** Input to {@link buildApprovedByteFixtureCandidate}: an entrypoint and its already-emitted source set. */
export type ApprovedByteFixtureInput = Readonly<{
  entrypoint: string;
  sources: readonly [ApprovedByteFixtureSource, ...ApprovedByteFixtureSource[]];
}>;

/** The raw encoded byte arrays a built {@link ApprovedByteFixtureCandidate}'s digests are computed over. */
export type ApprovedByteCandidateBytes = Readonly<{
  transformerProfile: readonly number[];
  normalizedOptions: readonly number[];
  originalManifest: readonly number[];
  executableManifest: readonly number[];
  transformationRecords: readonly (readonly number[])[];
  transformationRecordSet: readonly number[];
  planSourceBindings: readonly number[];
}>;

/** The SHA-256 digests computed over a built {@link ApprovedByteFixtureCandidate}'s {@link ApprovedByteCandidateBytes}. */
export type ApprovedByteCandidateDigests = Readonly<{
  transformerProfile: ApprovedByteDigest<"transformer-profile">;
  normalizedOptions: ApprovedByteDigest<"normalized-options">;
  originalManifest: ApprovedByteDigest<"original-authoring-source-manifest">;
  executableManifest: ApprovedByteDigest<"executable-javascript-source-manifest">;
  transformationRecordSet: ApprovedByteDigest<"transformation-record-set">;
}>;

/** The result of {@link buildApprovedByteFixtureCandidate}: encoded bytes paired with their digests. */
export type ApprovedByteFixtureCandidate = Readonly<{
  bytes: ApprovedByteCandidateBytes;
  digests: ApprovedByteCandidateDigests;
}>;

/** The retained expected values {@link verifyApprovedByteFixtureKnownAnswers} checks a candidate against. */
export type ApprovedByteKnownAnswers = Readonly<{
  transformerProfile: string;
  normalizedOptions: string;
  originalManifest: string;
  executableManifest: string;
  transformationRecords: readonly string[];
  transformationRecordSet: string;
  planSourceBindings: string;
}>;

/** The exhaustive set of refusal codes {@link buildApprovedByteFixtureCandidate} can return. */
export type ApprovedByteFixtureBuildRefusalCode =
  | "SOURCE_FILE_COUNT"
  | "SOURCE_PATH"
  | "SOURCE_PATH_DUPLICATE"
  | "SOURCE_ENTRYPOINT"
  | "ORIGINAL_SOURCE_BYTES"
  | "ORIGINAL_SOURCE_BOM"
  | "ORIGINAL_SOURCE_UTF8"
  | "EMITTED_SOURCE_BYTES"
  | "EMITTED_SOURCE_BOM"
  | "EMITTED_SOURCE_UTF8"
  | "ORIGINAL_AGGREGATE_BYTES"
  | "EMITTED_AGGREGATE_BYTES";

/** A classified refusal from {@link buildApprovedByteFixtureCandidate}. */
export interface ApprovedByteFixtureBuildRefusal {
  readonly owner: "approved-byte-fixture-builder";
  readonly classification: "MALFORMED" | "SCHEMA" | "DOMAIN";
  readonly code: ApprovedByteFixtureBuildRefusalCode;
}

/** The outcome of {@link buildApprovedByteFixtureCandidate}: either the built candidate or a classified refusal. */
export type ApprovedByteFixtureBuildResult =
  | { readonly ok: true; readonly candidate: ApprovedByteFixtureCandidate }
  | { readonly ok: false; readonly refusal: ApprovedByteFixtureBuildRefusal };

/** The exhaustive set of refusal codes {@link verifyApprovedByteFixtureKnownAnswers} can return. */
export type ApprovedByteFixtureVerifyRefusalCode =
  | "KNOWN_ANSWER_MISMATCH"
  | "TRANSFORMATION_RECORD_MISMATCH"
  | "DIGEST_BINDING_MISMATCH";

/** A classified refusal from {@link verifyApprovedByteFixtureKnownAnswers}. */
export interface ApprovedByteFixtureVerifyRefusal {
  readonly owner: "approved-byte-fixture-verifier";
  readonly classification: "DOMAIN";
  readonly code: ApprovedByteFixtureVerifyRefusalCode;
}

/** The outcome of {@link verifyApprovedByteFixtureKnownAnswers}: either a confirmed match or a classified refusal. */
export type ApprovedByteFixtureVerifyResult =
  | { readonly ok: true }
  | { readonly ok: false; readonly refusal: ApprovedByteFixtureVerifyRefusal };

class ApprovedByteFixtureBuildFailure extends Error {
  constructor(
    readonly classification: ApprovedByteFixtureBuildRefusal["classification"],
    readonly code: ApprovedByteFixtureBuildRefusalCode,
  ) {
    super(code);
  }
}

const sourcePathPattern =
  /^[A-Za-z0-9][A-Za-z0-9._-]{0,63}(?:\/[A-Za-z0-9][A-Za-z0-9._-]{0,63})*$/u;
const strictUtf8 = new TextDecoder("utf-8", { fatal: true, ignoreBOM: true });

/**
 * Builds only deterministic fixture candidates from already-emitted bytes.
 * It intentionally has no transformer call and grants no execute eligibility.
 */
export function buildApprovedByteFixtureCandidate(
  input: ApprovedByteFixtureInput,
): ApprovedByteFixtureBuildResult {
  try {
    return Object.freeze({ ok: true, candidate: buildCandidate(input) });
  } catch (error) {
    if (error instanceof ApprovedByteFixtureBuildFailure) {
      return buildRejected(error.classification, error.code);
    }
    throw error;
  }
}

function buildCandidate(input: ApprovedByteFixtureInput): ApprovedByteFixtureCandidate {
  if (input.sources.length > APPROVED_BYTE_CAPS.sourceFiles) {
    buildFail("DOMAIN", "SOURCE_FILE_COUNT");
  }
  const copied = input.sources
    .map((source) => ({
      logicalPath: validateSourcePath(source.logicalPath),
      originalBytes: validateSourceBytes(source.originalBytes),
      emittedBytes: validateEmittedBytes(source.emittedBytes),
    }))
    .sort((left, right) => compareAscii(left.logicalPath, right.logicalPath));
  for (let index = 1; index < copied.length; index += 1) {
    if (copied[index - 1]?.logicalPath === copied[index]?.logicalPath) {
      buildFail("DOMAIN", "SOURCE_PATH_DUPLICATE");
    }
  }
  if (!copied.some((source) => source.logicalPath === input.entrypoint)) {
    buildFail("DOMAIN", "SOURCE_ENTRYPOINT");
  }
  const originalAggregate = copied.reduce(
    (total, source) => total + source.originalBytes.length,
    0,
  );
  const emittedAggregate = copied.reduce((total, source) => total + source.emittedBytes.length, 0);
  if (originalAggregate > APPROVED_BYTE_CAPS.originalAggregateBytes) {
    buildFail("DOMAIN", "ORIGINAL_AGGREGATE_BYTES");
  }
  if (emittedAggregate > APPROVED_BYTE_CAPS.emittedAggregateBytes) {
    buildFail("DOMAIN", "EMITTED_AGGREGATE_BYTES");
  }

  const transformerProfile = encodeTransformerProfile();
  const normalizedOptions = encodeNormalizedOptions();
  const transformerDigest = digest(transformerProfile, "transformer-profile");
  const optionsDigest = digest(normalizedOptions, "normalized-options");
  const transformationRecords = copied.map((source) =>
    encodeTransformationRecord(source, transformerDigest, normalizedOptions, optionsDigest),
  );
  const originalManifest = encodeSourceManifest(
    "capsule.original-authoring-source-manifest",
    input.entrypoint,
    copied.map((source) => [
      source.logicalPath,
      TYPESCRIPT_SOURCE_MEDIA_TYPE,
      sha256(source.originalBytes),
      source.originalBytes.length,
    ]),
    originalAggregate,
  );
  const executableManifest = encodeSourceManifest(
    "capsule.executable-javascript-source-manifest",
    input.entrypoint,
    copied.map((source) => [
      source.logicalPath,
      JAVASCRIPT_SOURCE_MEDIA_TYPE,
      sha256(source.emittedBytes),
      source.emittedBytes.length,
    ]),
    emittedAggregate,
  );
  const transformationRecordSet = cborMap([
    [1, "capsule.typescript-transformation-record-set"],
    [2, 0],
    [3, transformationRecords],
  ]);
  const originalManifestDigest = digest(originalManifest, "original-authoring-source-manifest");
  const executableManifestDigest = digest(
    executableManifest,
    "executable-javascript-source-manifest",
  );
  const transformationRecordSetDigest = digest(
    transformationRecordSet,
    "transformation-record-set",
  );
  const planSourceBindings = cborMap([
    [1, "capsule.typescript-plan-source-bindings"],
    [2, 0],
    [3, Uint8Array.from(originalManifestDigest)],
    [4, Uint8Array.from(executableManifestDigest)],
    [5, Uint8Array.from(transformationRecordSetDigest)],
  ]);

  return Object.freeze({
    bytes: Object.freeze({
      transformerProfile: retain(transformerProfile),
      normalizedOptions: retain(normalizedOptions),
      originalManifest: retain(originalManifest),
      executableManifest: retain(executableManifest),
      transformationRecords: Object.freeze(transformationRecords.map(retain)),
      transformationRecordSet: retain(transformationRecordSet),
      planSourceBindings: retain(planSourceBindings),
    }),
    digests: Object.freeze({
      transformerProfile: transformerDigest,
      normalizedOptions: optionsDigest,
      originalManifest: originalManifestDigest,
      executableManifest: executableManifestDigest,
      transformationRecordSet: transformationRecordSetDigest,
    }),
  });
}

/** Wraps a raw 32-byte SHA-256 digest as an {@link ApprovedByteDigest} bound to `_role`; throws if `value` is not exactly 32 bytes. */
export function asApprovedByteDigest<Role extends ApprovedByteDigestRole>(
  value: Uint8Array,
  _role: Role,
): ApprovedByteDigest<Role> {
  if (value.length !== 32) {
    throw new TypeError("approved-byte SHA-256 digest must contain exactly 32 bytes");
  }
  return Object.freeze(Array.from(value)) as ApprovedByteDigest<Role>;
}

/**
 * Bounded fixture-only verification. The caller supplies retained known
 * answers; this function is not a general CBOR decoder or admission port.
 */
export function verifyApprovedByteFixtureKnownAnswers(
  candidate: ApprovedByteFixtureCandidate,
  expected: ApprovedByteKnownAnswers,
): ApprovedByteFixtureVerifyResult {
  const actual = {
    transformerProfile: hashHex(candidate.bytes.transformerProfile),
    normalizedOptions: hashHex(candidate.bytes.normalizedOptions),
    originalManifest: hashHex(candidate.bytes.originalManifest),
    executableManifest: hashHex(candidate.bytes.executableManifest),
    transformationRecords: candidate.bytes.transformationRecords.map(hashHex),
    transformationRecordSet: hashHex(candidate.bytes.transformationRecordSet),
    planSourceBindings: hashHex(candidate.bytes.planSourceBindings),
  };
  for (const role of [
    "transformerProfile",
    "normalizedOptions",
    "originalManifest",
    "executableManifest",
    "transformationRecordSet",
    "planSourceBindings",
  ] as const) {
    if (actual[role] !== expected[role]) {
      return verifyRejected("KNOWN_ANSWER_MISMATCH");
    }
  }
  if (
    actual.transformationRecords.length !== expected.transformationRecords.length ||
    actual.transformationRecords.some(
      (value, index) => value !== expected.transformationRecords[index],
    )
  ) {
    return verifyRejected("TRANSFORMATION_RECORD_MISMATCH");
  }
  if (
    !Buffer.from(candidate.digests.transformerProfile).equals(
      Buffer.from(expected.transformerProfile, "hex"),
    ) ||
    !Buffer.from(candidate.digests.normalizedOptions).equals(
      Buffer.from(expected.normalizedOptions, "hex"),
    ) ||
    !Buffer.from(candidate.digests.originalManifest).equals(
      Buffer.from(expected.originalManifest, "hex"),
    ) ||
    !Buffer.from(candidate.digests.executableManifest).equals(
      Buffer.from(expected.executableManifest, "hex"),
    ) ||
    !Buffer.from(candidate.digests.transformationRecordSet).equals(
      Buffer.from(expected.transformationRecordSet, "hex"),
    )
  ) {
    return verifyRejected("DIGEST_BINDING_MISMATCH");
  }
  return Object.freeze({ ok: true });
}

function validateSourcePath(value: string): string {
  if (value.length > 256 || !sourcePathPattern.test(value)) {
    buildFail("SCHEMA", "SOURCE_PATH");
  }
  return value;
}

function validateSourceBytes(value: Uint8Array): Uint8Array {
  if (value.length > APPROVED_BYTE_CAPS.originalFileBytes) {
    buildFail("DOMAIN", "ORIGINAL_SOURCE_BYTES");
  }
  validateUtf8WithoutBom(value, "original");
  return Uint8Array.from(value);
}

function validateEmittedBytes(value: Uint8Array): Uint8Array {
  if (value.length > APPROVED_BYTE_CAPS.emittedFileBytes) {
    buildFail("DOMAIN", "EMITTED_SOURCE_BYTES");
  }
  validateUtf8WithoutBom(value, "emitted");
  return Uint8Array.from(value);
}

function validateUtf8WithoutBom(value: Uint8Array, role: "original" | "emitted"): void {
  if (value[0] === 0xef && value[1] === 0xbb && value[2] === 0xbf) {
    buildFail("MALFORMED", role === "original" ? "ORIGINAL_SOURCE_BOM" : "EMITTED_SOURCE_BOM");
  }
  try {
    strictUtf8.decode(value);
  } catch {
    buildFail("MALFORMED", role === "original" ? "ORIGINAL_SOURCE_UTF8" : "EMITTED_SOURCE_UTF8");
  }
}

function buildFail(
  classification: ApprovedByteFixtureBuildRefusal["classification"],
  code: ApprovedByteFixtureBuildRefusalCode,
): never {
  throw new ApprovedByteFixtureBuildFailure(classification, code);
}

function buildRejected(
  classification: ApprovedByteFixtureBuildRefusal["classification"],
  code: ApprovedByteFixtureBuildRefusalCode,
): { readonly ok: false; readonly refusal: ApprovedByteFixtureBuildRefusal } {
  return Object.freeze({
    ok: false,
    refusal: Object.freeze({ owner: "approved-byte-fixture-builder", classification, code }),
  });
}

function verifyRejected(code: ApprovedByteFixtureVerifyRefusalCode): {
  readonly ok: false;
  readonly refusal: ApprovedByteFixtureVerifyRefusal;
} {
  return Object.freeze({
    ok: false,
    refusal: Object.freeze({
      owner: "approved-byte-fixture-verifier",
      classification: "DOMAIN",
      code,
    }),
  });
}

function encodeTransformerProfile(): Uint8Array {
  return cborMap([
    [1, "capsule.typescript-transformer-profile"],
    [2, 0],
    [3, APPROVED_BYTE_TOOLCHAIN.transformer],
    [4, APPROVED_BYTE_TOOLCHAIN.nodeVersion],
    [5, APPROVED_BYTE_TOOLCHAIN.amaroVersion],
    [6, toolchainDigestBytes(APPROVED_BYTE_TOOLCHAIN.sourceArchiveSha256)],
    [7, APPROVED_BYTE_TOOLCHAIN.platform],
    [8, APPROVED_BYTE_TOOLCHAIN.architecture],
    [9, toolchainDigestBytes(APPROVED_BYTE_TOOLCHAIN.distributionArchiveSha256)],
    [10, toolchainDigestBytes(APPROVED_BYTE_TOOLCHAIN.executableSha256)],
  ]);
}

function encodeNormalizedOptions(): Uint8Array {
  return cborMap([
    [1, "capsule.typescript-normalized-options"],
    [2, 0],
    [3, APPROVED_BYTE_OPTIONS.mode],
    [4, APPROVED_BYTE_OPTIONS.diagnosticPolicy],
    [5, APPROVED_BYTE_OPTIONS.sourceMap],
    [6, APPROVED_BYTE_OPTIONS.sourceUrl],
    [7, APPROVED_BYTE_OPTIONS.inputMediaType],
    [8, APPROVED_BYTE_OPTIONS.outputMediaType],
  ]);
}

function encodeTransformationRecord(
  source: { logicalPath: string; originalBytes: Uint8Array; emittedBytes: Uint8Array },
  transformerDigest: ApprovedByteDigest<"transformer-profile">,
  normalizedOptions: Uint8Array,
  optionsDigest: ApprovedByteDigest<"normalized-options">,
): Uint8Array {
  return cborMap([
    [1, "capsule.typescript-transformation-record"],
    [2, 0],
    [3, source.logicalPath],
    [4, TYPESCRIPT_SOURCE_MEDIA_TYPE],
    [5, source.originalBytes.length],
    [6, sha256(source.originalBytes)],
    [7, JAVASCRIPT_SOURCE_MEDIA_TYPE],
    [8, source.emittedBytes.length],
    [9, sha256(source.emittedBytes)],
    [10, Uint8Array.from(transformerDigest)],
    [11, APPROVED_BYTE_TOOLCHAIN.nodeVersion],
    [12, APPROVED_BYTE_TOOLCHAIN.amaroVersion],
    [13, toolchainDigestBytes(APPROVED_BYTE_TOOLCHAIN.sourceArchiveSha256)],
    [14, toolchainDigestBytes(APPROVED_BYTE_TOOLCHAIN.distributionArchiveSha256)],
    [15, toolchainDigestBytes(APPROVED_BYTE_TOOLCHAIN.executableSha256)],
    [16, normalizedOptions],
    [17, Uint8Array.from(optionsDigest)],
    [18, "absent"],
    [19, "absent"],
    [20, "reject-any"],
    [21, 0],
  ]);
}

function encodeSourceManifest(
  objectType: string,
  entrypoint: string,
  entries: readonly (readonly [string, string, Uint8Array, number])[],
  aggregate: number,
): Uint8Array {
  return cborMap([
    [1, objectType],
    [2, 0],
    [3, entrypoint],
    [4, entries],
    [5, aggregate],
  ]);
}

type CborValue = number | string | Uint8Array | readonly CborValue[];

function cborMap(entries: readonly (readonly [number, CborValue])[]): Uint8Array {
  return concatenateBytes([
    encodeCborMapHeader(entries.length),
    ...entries.flatMap(([key, value]) => [encodeCborUnsigned(key), cbor(value)]),
  ]);
}

function cbor(value: CborValue): Uint8Array {
  if (value instanceof Uint8Array) {
    return encodeCborByteString(value);
  }
  if (typeof value === "number") {
    return encodeCborUnsigned(value);
  }
  if (typeof value === "string") {
    return encodeCborText(value);
  }
  return concatenateBytes([encodeCborArrayHeader(value.length), ...value.map(cbor)]);
}

function digest<Role extends ApprovedByteDigestRole>(
  value: Uint8Array,
  role: Role,
): ApprovedByteDigest<Role> {
  return asApprovedByteDigest(sha256(value), role);
}

function sha256(value: Uint8Array): Uint8Array {
  return Uint8Array.from(createHash("sha256").update(value).digest());
}

function hashHex(value: ArrayLike<number>): string {
  return createHash("sha256").update(Uint8Array.from(value)).digest("hex");
}

function fromHex(value: string): Uint8Array {
  return Uint8Array.from(Buffer.from(value, "hex"));
}

function toolchainDigestBytes(value: string): Uint8Array {
  const bytes = fromHex(value);
  if (bytes.length !== 32) {
    throw new TypeError("approved-byte toolchain SHA-256 digest must contain exactly 32 bytes");
  }
  return bytes;
}

function retain(value: Uint8Array): readonly number[] {
  return Object.freeze(Array.from(value));
}

function compareAscii(left: string, right: string): number {
  return Buffer.compare(Buffer.from(left, "ascii"), Buffer.from(right, "ascii"));
}
