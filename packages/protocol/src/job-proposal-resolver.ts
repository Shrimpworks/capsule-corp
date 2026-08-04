import { createHash } from "node:crypto";

import { isRetainedDecodedJobProposalCandidate } from "./decoded-job-proposal-candidate.js";
import {
  concatenateBytes,
  encodeCborArrayHeader,
  encodeCborByteString,
  encodeCborMapHeader,
  encodeCborText,
  encodeCborUnsigned,
} from "./internal-cbor-primitives.js";
import {
  asPositiveSafeInteger,
  asRuntimeProfileAlias,
  type DecodedJobProposalCandidate,
  type InlineJsonValue,
  type PositiveSafeInteger,
  PRIMARY_DATA_INPUT_SLOT,
  type RuntimeProfileAlias,
  type SourceEntrypoint,
  type SourcePath,
  TRANSFORMED_JSON_OUTPUT_SLOT,
} from "./job-proposal.js";
import { retainResolvedJobProposalPlanInputs } from "./resolved-job-proposal-provenance.js";

const MAX_SOURCE_FILE_BYTES = 262_144;
const MAX_SOURCE_AGGREGATE_BYTES = 1_048_576;
const MAX_CANONICAL_INLINE_INPUT_BYTES = 262_144;

declare const resolvedScalarBrand: unique symbol;
declare const trustedContextBrand: unique symbol;

/** The digest roles this resolver computes. Distinct from {@link CandidateDigestRole}; these are resolver-local, hex-encoded, and computed from resolved bytes rather than decoded from the wire. */
export type ResolvedDigestRole = "source-content" | "source-manifest" | "inline-input";
/** A role-tagged lowercase-hex SHA-256 digest string computed by this resolver. */
export type ResolvedDigestHex<Role extends ResolvedDigestRole> = string & {
  readonly [resolvedScalarBrand]: `sha256:${Role}:hex`;
};

/** An immutable defensive copy of resolved bytes, safe to hand to callers without aliasing internal buffers. */
export interface RetainedExactBytes {
  readonly byteLength: number;
  readonly bytes: readonly number[];
}

/** Caller-supplied input to {@link createTrustedProfileRegistryContext}, before validation/branding. */
export interface TrustedRuntimeProfileEntryInput {
  readonly alias: string;
  readonly status: "active" | "inactive";
  readonly exactWallTimeMs: {
    readonly minimum: number;
    readonly maximum: number;
  };
}

/** A validated, branded runtime profile entry as it appears inside a {@link TrustedProfileRegistryContext}. */
export interface TrustedRuntimeProfileEntry {
  readonly alias: RuntimeProfileAlias;
  readonly status: "active" | "inactive";
  readonly exactWallTimeMs: {
    readonly minimum: PositiveSafeInteger;
    readonly maximum: PositiveSafeInteger;
  };
}

/**
 * An immutable, provenance-checked registry of trusted runtime profiles,
 * produced only by {@link createTrustedProfileRegistryContext}.
 * resolveJobProposal refuses any object not created by that function, even
 * if it has an identical shape.
 */
export interface TrustedProfileRegistryContext {
  readonly profiles: readonly TrustedRuntimeProfileEntry[];
  readonly [trustedContextBrand]: "profile-registry";
}

/**
 * An immutable, provenance-checked trusted user resource policy, produced
 * only by {@link createTrustedUserPolicyContext}. Encodes ADR-0009's
 * trusted-default-plus-hard-ceiling model: a request above the ceiling is
 * rejected, never silently clamped.
 */
export interface TrustedUserPolicyContext {
  readonly wallTimeMs: {
    readonly trustedDefault: PositiveSafeInteger;
    readonly ceiling: PositiveSafeInteger;
  };
  readonly [trustedContextBrand]: "user-policy";
}

/** One resolved source file: its validated path, retained UTF-8 bytes, and content digest. */
export interface ResolvedSourceFile {
  readonly path: SourcePath;
  readonly utf8Bytes: RetainedExactBytes;
  readonly contentDigest: ResolvedDigestHex<"source-content">;
}

/** The deterministic-CBOR-encoded source manifest and its digest, built from the resolved file set. */
export interface ResolvedSourceManifest {
  readonly exactBytes: RetainedExactBytes;
  readonly digest: ResolvedDigestHex<"source-manifest">;
}

/** The fully resolved source bundle: entrypoint, sorted files, aggregate size, and manifest. */
export interface ResolvedSourceBundle {
  readonly entrypoint: SourceEntrypoint;
  readonly files: readonly ResolvedSourceFile[];
  readonly aggregateByteLength: number;
  readonly manifest: ResolvedSourceManifest;
}

/** The resolved inline JSON input: its canonical encoded bytes and digest. */
export interface ResolvedInlineInput {
  readonly slot: typeof PRIMARY_DATA_INPUT_SLOT;
  readonly exactBytes: RetainedExactBytes;
  readonly digest: ResolvedDigestHex<"inline-input">;
}

/** The resolved wall-time limit and whether it came from the request or the trusted default. */
export interface ResolvedWallTime {
  readonly milliseconds: PositiveSafeInteger;
  readonly origin: "requested" | "trusted-default";
}

/**
 * The exact proposal-owned fields that a later Task 3C plan builder may consume.
 * This value is not an ExecutionPlan and contains no installation, trust,
 * runtime-bundle, backend, policy-decision, expiry, approval, or launch authority.
 */
export interface ResolvedJobProposalPlanInputs {
  readonly sourceManifestDigest: ResolvedDigestHex<"source-manifest">;
  readonly sourceEntrypoint: SourceEntrypoint;
  readonly sourceByteLength: number;
  readonly inputSlot: typeof PRIMARY_DATA_INPUT_SLOT;
  readonly inlineInputDigest: ResolvedDigestHex<"inline-input">;
  readonly inlineInputByteLength: number;
  readonly runtimeProfileAlias: RuntimeProfileAlias;
  readonly wallTimeMs: PositiveSafeInteger;
  readonly wallTimeOrigin: "requested" | "trusted-default";
  readonly outputSlot: typeof TRANSFORMED_JSON_OUTPUT_SLOT;
  readonly outputMaxJsonBytes: PositiveSafeInteger;
}

/**
 * The full output of a successful {@link resolveJobProposal} call: every
 * field the caller requested resolved against trusted context, plus the
 * narrower {@link ResolvedJobProposalPlanInputs} subset that
 * constructExecutionPlan actually consumes.
 */
export interface ResolvedJobProposal {
  readonly source: ResolvedSourceBundle;
  readonly inlineInput: ResolvedInlineInput;
  readonly runtimeProfile: TrustedRuntimeProfileEntry & { readonly status: "active" };
  readonly wallTime: ResolvedWallTime;
  readonly output: {
    readonly slot: typeof TRANSFORMED_JSON_OUTPUT_SLOT;
    readonly maxJsonBytes: PositiveSafeInteger;
  };
  readonly labels?: Readonly<Record<string, string>>;
  readonly planInputs: ResolvedJobProposalPlanInputs;
}

/**
 * Why a resolution was refused. `SEMANTIC` means the proposal's own content
 * is invalid (e.g. an oversized file); `POLICY` means it's well-formed but
 * violates trusted policy (e.g. a wall-time ceiling); `BINDING` means a
 * provenance/identity check failed (e.g. an unretained candidate or trusted
 * context).
 */
export type JobProposalResolutionClassification = "SEMANTIC" | "POLICY" | "BINDING";

/** The exhaustive set of refusal codes {@link resolveJobProposal} can return. */
export type JobProposalResolutionRefusalCode =
  | "CANDIDATE_PROVENANCE"
  | "TRUSTED_CONTEXT_PROVENANCE"
  | "SOURCE_FILE_BYTES"
  | "SOURCE_AGGREGATE_BYTES"
  | "SOURCE_ENTRYPOINT_MEMBERSHIP"
  | "INLINE_INPUT_BYTES"
  | "PROFILE_UNKNOWN"
  | "PROFILE_INACTIVE"
  | "WALL_TIME_CEILING"
  | "WALL_TIME_PROFILE_RANGE";

export interface JobProposalResolutionRefusal {
  readonly owner: "job-proposal-semantic-validator";
  readonly classification: JobProposalResolutionClassification;
  readonly code: JobProposalResolutionRefusalCode;
}

/** The outcome of {@link resolveJobProposal}: either a resolved proposal or a classified refusal. */
export type JobProposalResolutionResult =
  | { readonly ok: true; readonly resolved: ResolvedJobProposal }
  | { readonly ok: false; readonly refusal: JobProposalResolutionRefusal };

const trustedProfileRegistries = new WeakSet<object>();
const trustedUserPolicies = new WeakSet<object>();

/**
 * Validates and freezes a set of runtime profile entries into a
 * provenance-branded {@link TrustedProfileRegistryContext}. Rejects
 * duplicate aliases, unsupported statuses, and an inverted wall-time range.
 * Callers must treat the returned object as opaque — resolveJobProposal
 * only accepts contexts it can prove came from this function.
 */
export function createTrustedProfileRegistryContext(
  entries: readonly TrustedRuntimeProfileEntryInput[],
): TrustedProfileRegistryContext {
  const aliases = new Set<string>();
  const profiles = entries.map((entry) => {
    const alias = asRuntimeProfileAlias(entry.alias);
    if (aliases.has(alias)) {
      throw new TypeError("trusted profile registry contains a duplicate alias");
    }
    aliases.add(alias);
    if (entry.status !== "active" && entry.status !== "inactive") {
      throw new TypeError("trusted profile registry contains an unsupported status");
    }
    const minimum = asPositiveSafeInteger(entry.exactWallTimeMs.minimum);
    const maximum = asPositiveSafeInteger(entry.exactWallTimeMs.maximum);
    if (minimum > maximum) {
      throw new TypeError("trusted profile registry contains an invalid wall-time range");
    }
    return Object.freeze({
      alias,
      status: entry.status,
      exactWallTimeMs: Object.freeze({ minimum, maximum }),
    });
  });
  const context = Object.freeze({
    profiles: Object.freeze(profiles),
  }) as TrustedProfileRegistryContext;
  trustedProfileRegistries.add(context);
  return context;
}

/**
 * Validates and freezes a trusted-default-plus-ceiling policy into a
 * provenance-branded {@link TrustedUserPolicyContext} (ADR-0009). Rejects a
 * default that exceeds its own ceiling.
 */
export function createTrustedUserPolicyContext(input: {
  readonly wallTimeMs: { readonly trustedDefault: number; readonly ceiling: number };
}): TrustedUserPolicyContext {
  const trustedDefault = asPositiveSafeInteger(input.wallTimeMs.trustedDefault);
  const ceiling = asPositiveSafeInteger(input.wallTimeMs.ceiling);
  if (trustedDefault > ceiling) {
    throw new TypeError("trusted wall-time default exceeds its ceiling");
  }
  const context = Object.freeze({
    wallTimeMs: Object.freeze({ trustedDefault, ceiling }),
  }) as TrustedUserPolicyContext;
  trustedUserPolicies.add(context);
  return context;
}

/**
 * Resolves only a decoder-issued passive candidate against separately issued,
 * immutable trusted contexts. It performs no I/O and creates no authority.
 */
export function resolveJobProposal(
  proposal: DecodedJobProposalCandidate,
  profileRegistry: TrustedProfileRegistryContext,
  userPolicy: TrustedUserPolicyContext,
): JobProposalResolutionResult {
  if (!isRetainedDecodedJobProposalCandidate(proposal)) {
    return rejected("BINDING", "CANDIDATE_PROVENANCE");
  }
  if (!trustedProfileRegistries.has(profileRegistry) || !trustedUserPolicies.has(userPolicy)) {
    return rejected("BINDING", "TRUSTED_CONTEXT_PROVENANCE");
  }

  const sourceResult = resolveSource(proposal);
  if (!sourceResult.ok) {
    return sourceResult;
  }
  const inlineInputResult = resolveInlineInput(proposal.input.value);
  if (!inlineInputResult.ok) {
    return inlineInputResult;
  }
  const profile = profileRegistry.profiles.find((entry) => entry.alias === proposal.runtimeProfile);
  if (!profile) {
    return rejected("BINDING", "PROFILE_UNKNOWN");
  }
  if (profile.status !== "active") {
    return rejected("POLICY", "PROFILE_INACTIVE");
  }

  const requestedWallTime = proposal.requestedLimits.wallTimeMs;
  const wallTime = Object.freeze({
    milliseconds: requestedWallTime ?? userPolicy.wallTimeMs.trustedDefault,
    origin: requestedWallTime === undefined ? "trusted-default" : "requested",
  }) as ResolvedWallTime;
  if (wallTime.milliseconds > userPolicy.wallTimeMs.ceiling) {
    return rejected("POLICY", "WALL_TIME_CEILING");
  }
  if (
    wallTime.milliseconds < profile.exactWallTimeMs.minimum ||
    wallTime.milliseconds > profile.exactWallTimeMs.maximum
  ) {
    return rejected("POLICY", "WALL_TIME_PROFILE_RANGE");
  }

  const source = sourceResult.source;
  const inlineInput = inlineInputResult.inlineInput;
  const output = Object.freeze({
    slot: TRANSFORMED_JSON_OUTPUT_SLOT,
    maxJsonBytes: proposal.outputs[0].maxBytes,
  });
  const planInputs = retainResolvedJobProposalPlanInputs(
    Object.freeze({
      sourceManifestDigest: source.manifest.digest,
      sourceEntrypoint: source.entrypoint,
      sourceByteLength: source.aggregateByteLength,
      inputSlot: PRIMARY_DATA_INPUT_SLOT,
      inlineInputDigest: inlineInput.digest,
      inlineInputByteLength: inlineInput.exactBytes.byteLength,
      runtimeProfileAlias: profile.alias,
      wallTimeMs: wallTime.milliseconds,
      wallTimeOrigin: wallTime.origin,
      outputSlot: TRANSFORMED_JSON_OUTPUT_SLOT,
      outputMaxJsonBytes: output.maxJsonBytes,
    }),
  );
  const labels = proposal.labels
    ? Object.freeze(Object.fromEntries(Object.entries(proposal.labels).sort(compareEntries)))
    : undefined;
  const resolved = Object.freeze({
    source,
    inlineInput,
    runtimeProfile: profile as TrustedRuntimeProfileEntry & { readonly status: "active" },
    wallTime,
    output,
    ...(labels === undefined ? {} : { labels }),
    planInputs,
  });
  return Object.freeze({ ok: true, resolved });
}

function resolveSource(
  proposal: DecodedJobProposalCandidate,
):
  | { readonly ok: true; readonly source: ResolvedSourceBundle }
  | { readonly ok: false; readonly refusal: JobProposalResolutionRefusal } {
  const sourceEntries = Object.entries(proposal.source.files).sort(compareEntries);
  if (!Object.hasOwn(proposal.source.files, proposal.source.entrypoint)) {
    return rejected("SEMANTIC", "SOURCE_ENTRYPOINT_MEMBERSHIP");
  }

  let aggregateByteLength = 0;
  const files: ResolvedSourceFile[] = [];
  for (const [path, content] of sourceEntries) {
    const contentBytes = new TextEncoder().encode(content);
    if (contentBytes.byteLength > MAX_SOURCE_FILE_BYTES) {
      return rejected("SEMANTIC", "SOURCE_FILE_BYTES");
    }
    aggregateByteLength += contentBytes.byteLength;
    if (aggregateByteLength > MAX_SOURCE_AGGREGATE_BYTES) {
      return rejected("SEMANTIC", "SOURCE_AGGREGATE_BYTES");
    }
    files.push(
      Object.freeze({
        path: path as SourcePath,
        utf8Bytes: retainExactBytes(contentBytes),
        contentDigest: digestHex(contentBytes, "source-content"),
      }),
    );
  }

  const manifestBytes = encodeSourceManifest(
    proposal.source.entrypoint,
    files,
    aggregateByteLength,
  );
  const manifest = Object.freeze({
    exactBytes: retainExactBytes(manifestBytes),
    digest: digestHex(manifestBytes, "source-manifest"),
  });
  return Object.freeze({
    ok: true,
    source: Object.freeze({
      entrypoint: proposal.source.entrypoint,
      files: Object.freeze(files),
      aggregateByteLength,
      manifest,
    }),
  });
}

function resolveInlineInput(
  value: InlineJsonValue,
):
  | { readonly ok: true; readonly inlineInput: ResolvedInlineInput }
  | { readonly ok: false; readonly refusal: JobProposalResolutionRefusal } {
  const writer = new BoundedByteWriter(MAX_CANONICAL_INLINE_INPUT_BYTES);
  if (!writeCanonicalInlineJson(writer, value)) {
    return rejected("SEMANTIC", "INLINE_INPUT_BYTES");
  }
  const exactBytes = writer.finish();
  return Object.freeze({
    ok: true,
    inlineInput: Object.freeze({
      slot: PRIMARY_DATA_INPUT_SLOT,
      exactBytes: retainExactBytes(exactBytes),
      digest: digestHex(exactBytes, "inline-input"),
    }),
  });
}

function writeCanonicalInlineJson(writer: BoundedByteWriter, value: InlineJsonValue): boolean {
  if (value === null) {
    return writer.writeAscii("null");
  }
  if (value === true) {
    return writer.writeAscii("true");
  }
  if (value === false) {
    return writer.writeAscii("false");
  }
  if (typeof value === "number") {
    return writer.writeAscii(String(value));
  }
  if (typeof value === "string") {
    return writeCanonicalString(writer, value);
  }
  if (Array.isArray(value)) {
    if (!writer.writeAscii("[")) {
      return false;
    }
    for (let index = 0; index < value.length; index += 1) {
      if (
        (index > 0 && !writer.writeAscii(",")) ||
        !writeCanonicalInlineJson(writer, value[index])
      ) {
        return false;
      }
    }
    return writer.writeAscii("]");
  }

  if (!writer.writeAscii("{")) {
    return false;
  }
  const entries = Object.entries(value).sort(compareEntries);
  for (let index = 0; index < entries.length; index += 1) {
    const [key, entry] = entries[index] as [string, InlineJsonValue];
    if (
      (index > 0 && !writer.writeAscii(",")) ||
      !writeCanonicalString(writer, key) ||
      !writer.writeAscii(":") ||
      !writeCanonicalInlineJson(writer, entry)
    ) {
      return false;
    }
  }
  return writer.writeAscii("}");
}

function writeCanonicalString(writer: BoundedByteWriter, value: string): boolean {
  if (!writer.writeAscii('"')) {
    return false;
  }
  let runStart = 0;
  for (let index = 0; index < value.length; ) {
    const codePoint = value.codePointAt(index) as number;
    const scalarWidth = codePoint > 0xffff ? 2 : 1;
    const scalar = value.slice(index, index + scalarWidth);
    let escapedScalar: string | undefined;
    if (scalar === '"') {
      escapedScalar = '\\"';
    } else if (scalar === "\\") {
      escapedScalar = "\\\\";
    } else if (codePoint <= 0x1f) {
      escapedScalar = `\\u00${codePoint.toString(16).padStart(2, "0")}`;
    }
    if (escapedScalar !== undefined) {
      if (
        (runStart < index && !writer.writeUtf8(value.slice(runStart, index))) ||
        !writer.writeAscii(escapedScalar)
      ) {
        return false;
      }
      runStart = index + scalarWidth;
    }
    index += scalarWidth;
  }
  if (runStart < value.length && !writer.writeUtf8(value.slice(runStart))) {
    return false;
  }
  return writer.writeAscii('"');
}

class BoundedByteWriter {
  readonly #chunks: Uint8Array[] = [];
  #byteLength = 0;

  constructor(readonly maximumByteLength: number) {}

  writeAscii(value: string): boolean {
    return this.writeBytes(new TextEncoder().encode(value));
  }

  writeUtf8(value: string): boolean {
    return this.writeBytes(new TextEncoder().encode(value));
  }

  finish(): Uint8Array {
    const bytes = new Uint8Array(this.#byteLength);
    let offset = 0;
    for (const chunk of this.#chunks) {
      bytes.set(chunk, offset);
      offset += chunk.byteLength;
    }
    return bytes;
  }

  private writeBytes(bytes: Uint8Array): boolean {
    if (this.#byteLength + bytes.byteLength > this.maximumByteLength) {
      return false;
    }
    this.#chunks.push(bytes);
    this.#byteLength += bytes.byteLength;
    return true;
  }
}

function encodeSourceManifest(
  entrypoint: SourceEntrypoint,
  files: readonly ResolvedSourceFile[],
  aggregateByteLength: number,
): Uint8Array {
  return concatenateBytes([
    encodeCborMapHeader(5),
    encodeCborUnsigned(1),
    encodeCborText("capsule.source-manifest"),
    encodeCborUnsigned(2),
    encodeCborUnsigned(0),
    encodeCborUnsigned(3),
    encodeCborText(entrypoint),
    encodeCborUnsigned(4),
    encodeCborArrayHeader(files.length),
    ...files.flatMap((file) => [
      encodeCborArrayHeader(3),
      encodeCborText(file.path),
      encodeCborByteString(hexToBytes(file.contentDigest)),
      encodeCborUnsigned(file.utf8Bytes.byteLength),
    ]),
    encodeCborUnsigned(5),
    encodeCborUnsigned(aggregateByteLength),
  ]);
}

function retainExactBytes(bytes: Uint8Array): RetainedExactBytes {
  return Object.freeze({
    byteLength: bytes.byteLength,
    bytes: Object.freeze(Array.from(bytes)),
  });
}

function digestHex<Role extends ResolvedDigestRole>(
  bytes: Uint8Array,
  _role: Role,
): ResolvedDigestHex<Role> {
  return createHash("sha256").update(bytes).digest("hex") as ResolvedDigestHex<Role>;
}

function hexToBytes(value: string): Uint8Array {
  return Uint8Array.from(Buffer.from(value, "hex"));
}

function compareEntries(
  left: readonly [string, unknown],
  right: readonly [string, unknown],
): number {
  return compareAscii(left[0], right[0]);
}

/**
 * Compares UTF-16 code units. Callers must guarantee ASCII-only input: for
 * ASCII strings, UTF-16-code-unit order matches canonical UTF-8-byte order,
 * but the equivalence does not hold for non-ASCII input. Do not widen a
 * caller's key/path pattern to allow non-ASCII without revisiting this sort.
 */
function compareAscii(left: string, right: string): number {
  const length = Math.min(left.length, right.length);
  for (let index = 0; index < length; index += 1) {
    const difference = left.charCodeAt(index) - right.charCodeAt(index);
    if (difference !== 0) {
      return difference;
    }
  }
  return left.length - right.length;
}

function rejected(
  classification: JobProposalResolutionClassification,
  code: JobProposalResolutionRefusalCode,
): { readonly ok: false; readonly refusal: JobProposalResolutionRefusal } {
  return Object.freeze({
    ok: false,
    refusal: Object.freeze({
      owner: "job-proposal-semantic-validator",
      classification,
      code,
    }),
  });
}
