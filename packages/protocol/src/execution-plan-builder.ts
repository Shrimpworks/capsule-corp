import { createHash } from "node:crypto";

import type { ResolvedJobProposalPlanInputs, RetainedExactBytes } from "./job-proposal-resolver.js";
import { isRetainedResolvedJobProposalPlanInputs } from "./resolved-job-proposal-provenance.js";

const EXECUTION_PLAN_OBJECT_TYPE = "capsule.execution-plan" as const;
const EXECUTION_PLAN_OBJECT_VERSION = 0 as const;
const EXECUTION_PLAN_FIELD_COUNT = 24;
const MAX_REVIEW_ATTESTATIONS = 8;

declare const executionPlanScalarBrand: unique symbol;
declare const trustedExecutionPlanBindingBrand: unique symbol;

export type TrustedExecutionPlanDigestRole =
  | "trust-epoch"
  | "runtime-bundle-manifest"
  | "profile-review-attestation"
  | "profile-registry-entry"
  | "backend-validation-record"
  | "backend-configuration"
  | "trust-snapshot"
  | "policy-decision";

export type ExecutionPlanRetainedByteRole =
  | "installation-id"
  | "trust-epoch"
  | "source-manifest"
  | "inline-input"
  | "runtime-bundle-manifest"
  | "profile-review-attestation"
  | "profile-registry-entry"
  | "backend-validation-record"
  | "backend-configuration"
  | "trust-snapshot"
  | "policy-decision"
  | "execution-plan";

export interface RetainedExecutionPlanBytes<Role extends ExecutionPlanRetainedByteRole> {
  readonly role: Role;
  readonly byteLength: number;
  readonly bytes: readonly number[];
  readonly [executionPlanScalarBrand]: Role;
}

export type TrustedExecutionPlanInstallationId = RetainedExecutionPlanBytes<"installation-id"> & {
  readonly [trustedExecutionPlanBindingBrand]: "installation-id";
};

export type TrustedExecutionPlanDigest<Role extends TrustedExecutionPlanDigestRole> =
  RetainedExecutionPlanBytes<Role> & {
    readonly [trustedExecutionPlanBindingBrand]: Role;
  };

export interface TrustedExecutionPlanBindingsInput {
  readonly installationId: TrustedExecutionPlanInstallationId;
  readonly epochSequence: number;
  readonly epochDigest: TrustedExecutionPlanDigest<"trust-epoch">;
  readonly runtimeBundleManifestDigest: TrustedExecutionPlanDigest<"runtime-bundle-manifest">;
  readonly profileReviewAttestationDigests: readonly TrustedExecutionPlanDigest<"profile-review-attestation">[];
  readonly profileRegistryEntryDigest: TrustedExecutionPlanDigest<"profile-registry-entry">;
  readonly backendValidationRecordDigest: TrustedExecutionPlanDigest<"backend-validation-record">;
  readonly backendConfigurationDigest: TrustedExecutionPlanDigest<"backend-configuration">;
  readonly trustSnapshotDigest: TrustedExecutionPlanDigest<"trust-snapshot">;
  readonly policyDecisionDigest: TrustedExecutionPlanDigest<"policy-decision">;
  readonly expiresAt: number;
}

export interface TrustedExecutionPlanBindings {
  readonly installationId: TrustedExecutionPlanInstallationId;
  readonly epochSequence: number;
  readonly epochDigest: TrustedExecutionPlanDigest<"trust-epoch">;
  readonly runtimeBundleManifestDigest: TrustedExecutionPlanDigest<"runtime-bundle-manifest">;
  readonly profileReviewAttestationDigests: readonly [
    TrustedExecutionPlanDigest<"profile-review-attestation">,
    ...TrustedExecutionPlanDigest<"profile-review-attestation">[],
  ];
  readonly profileRegistryEntryDigest: TrustedExecutionPlanDigest<"profile-registry-entry">;
  readonly backendValidationRecordDigest: TrustedExecutionPlanDigest<"backend-validation-record">;
  readonly backendConfigurationDigest: TrustedExecutionPlanDigest<"backend-configuration">;
  readonly trustSnapshotDigest: TrustedExecutionPlanDigest<"trust-snapshot">;
  readonly policyDecisionDigest: TrustedExecutionPlanDigest<"policy-decision">;
  readonly expiresAt: number;
  readonly [trustedExecutionPlanBindingBrand]: "execution-plan-bindings";
}

export interface ExecutionPlanCandidateView {
  readonly objectType: typeof EXECUTION_PLAN_OBJECT_TYPE;
  readonly objectVersion: typeof EXECUTION_PLAN_OBJECT_VERSION;
  readonly installationId: RetainedExecutionPlanBytes<"installation-id">;
  readonly epochSequence: number;
  readonly epochDigest: RetainedExecutionPlanBytes<"trust-epoch">;
  readonly sourceManifestDigest: RetainedExecutionPlanBytes<"source-manifest">;
  readonly sourceEntrypoint: string;
  readonly sourceByteLength: number;
  readonly inputSlot: "primary-data";
  readonly inlineInputDigest: RetainedExecutionPlanBytes<"inline-input">;
  readonly inlineInputByteLength: number;
  readonly runtimeProfileAlias: string;
  readonly runtimeBundleManifestDigest: RetainedExecutionPlanBytes<"runtime-bundle-manifest">;
  readonly profileReviewAttestationDigests: readonly [
    RetainedExecutionPlanBytes<"profile-review-attestation">,
    ...RetainedExecutionPlanBytes<"profile-review-attestation">[],
  ];
  readonly profileRegistryEntryDigest: RetainedExecutionPlanBytes<"profile-registry-entry">;
  readonly backendValidationRecordDigest: RetainedExecutionPlanBytes<"backend-validation-record">;
  readonly backendConfigurationDigest: RetainedExecutionPlanBytes<"backend-configuration">;
  readonly trustSnapshotDigest: RetainedExecutionPlanBytes<"trust-snapshot">;
  readonly policyDecisionDigest: RetainedExecutionPlanBytes<"policy-decision">;
  readonly wallTimeMs: number;
  readonly wallTimeOrigin: "requested" | "trusted-default";
  readonly outputSlot: "transformed-json";
  readonly outputMaxJsonBytes: number;
  readonly expiresAt: number;
}

export interface ConstructedExecutionPlan {
  readonly candidate: ExecutionPlanCandidateView;
  readonly exactBytes: RetainedExactBytes;
  readonly digest: RetainedExecutionPlanBytes<"execution-plan">;
  copyExactBytes(): Uint8Array;
  copyDigestBytes(): Uint8Array;
}

export type ExecutionPlanConstructionRefusalCode =
  | "PLAN_INPUT_PROVENANCE"
  | "TRUSTED_BINDING_PROVENANCE";

export interface ExecutionPlanConstructionRefusal {
  readonly owner: "execution-plan-builder";
  readonly classification: "BINDING";
  readonly code: ExecutionPlanConstructionRefusalCode;
}

export type ExecutionPlanConstructionResult =
  | { readonly ok: true; readonly plan: ConstructedExecutionPlan }
  | { readonly ok: false; readonly refusal: ExecutionPlanConstructionRefusal };

const trustedInstallationIds = new WeakSet<object>();
const trustedDigests = new WeakMap<object, TrustedExecutionPlanDigestRole>();
const trustedBindings = new WeakSet<object>();

export function createTrustedExecutionPlanInstallationId(
  value: Uint8Array,
): TrustedExecutionPlanInstallationId {
  requireUint8Array(value, "installation ID");
  const copiedValue = Uint8Array.from(value);
  if (copiedValue.byteLength !== 16) {
    throw new TypeError("trusted installation ID must contain exactly 16 bytes");
  }
  if (copiedValue.every((byte) => byte === 0)) {
    throw new TypeError("trusted installation ID must be nonzero");
  }
  const retained = retainRoleBytes(
    copiedValue,
    "installation-id",
  ) as TrustedExecutionPlanInstallationId;
  trustedInstallationIds.add(retained);
  return retained;
}

export function createTrustedExecutionPlanDigest<Role extends TrustedExecutionPlanDigestRole>(
  value: Uint8Array,
  role: Role,
): TrustedExecutionPlanDigest<Role> {
  requireUint8Array(value, `${role} digest`);
  const copiedValue = Uint8Array.from(value);
  if (!TRUSTED_DIGEST_ROLES.has(role)) {
    throw new TypeError("trusted ExecutionPlan digest role is unsupported");
  }
  if (copiedValue.byteLength !== 32) {
    throw new TypeError(`trusted ${role} digest must contain exactly 32 bytes`);
  }
  const retained = retainRoleBytes(copiedValue, role) as TrustedExecutionPlanDigest<Role>;
  trustedDigests.set(retained, role);
  return retained;
}

export function createTrustedExecutionPlanBindings(
  input: TrustedExecutionPlanBindingsInput,
): TrustedExecutionPlanBindings {
  requireExactBindingInput(input);
  if (!trustedInstallationIds.has(input.installationId)) {
    throw new TypeError("trusted installation ID provenance is required");
  }
  requireUInt53(input.epochSequence, "trusted epoch sequence");
  requireDigestRole(input.epochDigest, "trust-epoch");
  requireDigestRole(input.runtimeBundleManifestDigest, "runtime-bundle-manifest");
  if (
    !Array.isArray(input.profileReviewAttestationDigests) ||
    input.profileReviewAttestationDigests.length < 1 ||
    input.profileReviewAttestationDigests.length > MAX_REVIEW_ATTESTATIONS
  ) {
    throw new TypeError(
      "trusted profile review attestations must contain one through eight digests",
    );
  }
  for (const digest of input.profileReviewAttestationDigests) {
    requireDigestRole(digest, "profile-review-attestation");
  }
  requireDigestRole(input.profileRegistryEntryDigest, "profile-registry-entry");
  requireDigestRole(input.backendValidationRecordDigest, "backend-validation-record");
  requireDigestRole(input.backendConfigurationDigest, "backend-configuration");
  requireDigestRole(input.trustSnapshotDigest, "trust-snapshot");
  requireDigestRole(input.policyDecisionDigest, "policy-decision");
  requireUInt53(input.expiresAt, "trusted plan expiry");

  const bindings = Object.freeze({
    installationId: copyInstallationId(input.installationId),
    epochSequence: input.epochSequence,
    epochDigest: copyDigest(input.epochDigest, "trust-epoch"),
    runtimeBundleManifestDigest: copyDigest(
      input.runtimeBundleManifestDigest,
      "runtime-bundle-manifest",
    ),
    profileReviewAttestationDigests: Object.freeze(
      input.profileReviewAttestationDigests.map((digest) =>
        copyDigest(digest, "profile-review-attestation"),
      ),
    ) as TrustedExecutionPlanBindings["profileReviewAttestationDigests"],
    profileRegistryEntryDigest: copyDigest(
      input.profileRegistryEntryDigest,
      "profile-registry-entry",
    ),
    backendValidationRecordDigest: copyDigest(
      input.backendValidationRecordDigest,
      "backend-validation-record",
    ),
    backendConfigurationDigest: copyDigest(
      input.backendConfigurationDigest,
      "backend-configuration",
    ),
    trustSnapshotDigest: copyDigest(input.trustSnapshotDigest, "trust-snapshot"),
    policyDecisionDigest: copyDigest(input.policyDecisionDigest, "policy-decision"),
    expiresAt: input.expiresAt,
  }) as TrustedExecutionPlanBindings;
  trustedBindings.add(bindings);
  return bindings;
}

/**
 * Constructs only the minimum passive candidate from Task 3B plan inputs and
 * separately issued trusted bindings. It performs no registration or I/O.
 */
export function constructExecutionPlan(
  planInputs: ResolvedJobProposalPlanInputs,
  bindings: TrustedExecutionPlanBindings,
): ExecutionPlanConstructionResult {
  if (!isRetainedResolvedJobProposalPlanInputs(planInputs)) {
    return rejected("PLAN_INPUT_PROVENANCE");
  }
  if (!trustedBindings.has(bindings)) {
    return rejected("TRUSTED_BINDING_PROVENANCE");
  }

  const reviewDigests = Object.freeze(
    bindings.profileReviewAttestationDigests.map(copyRetainedBytes),
  ) as ExecutionPlanCandidateView["profileReviewAttestationDigests"];
  const candidate = Object.freeze({
    objectType: EXECUTION_PLAN_OBJECT_TYPE,
    objectVersion: EXECUTION_PLAN_OBJECT_VERSION,
    installationId: copyRetainedBytes(bindings.installationId),
    epochSequence: bindings.epochSequence,
    epochDigest: copyRetainedBytes(bindings.epochDigest),
    sourceManifestDigest: retainedDigestFromHex(planInputs.sourceManifestDigest, "source-manifest"),
    sourceEntrypoint: planInputs.sourceEntrypoint,
    sourceByteLength: planInputs.sourceByteLength,
    inputSlot: planInputs.inputSlot,
    inlineInputDigest: retainedDigestFromHex(planInputs.inlineInputDigest, "inline-input"),
    inlineInputByteLength: planInputs.inlineInputByteLength,
    runtimeProfileAlias: planInputs.runtimeProfileAlias,
    runtimeBundleManifestDigest: copyRetainedBytes(bindings.runtimeBundleManifestDigest),
    profileReviewAttestationDigests: reviewDigests,
    profileRegistryEntryDigest: copyRetainedBytes(bindings.profileRegistryEntryDigest),
    backendValidationRecordDigest: copyRetainedBytes(bindings.backendValidationRecordDigest),
    backendConfigurationDigest: copyRetainedBytes(bindings.backendConfigurationDigest),
    trustSnapshotDigest: copyRetainedBytes(bindings.trustSnapshotDigest),
    policyDecisionDigest: copyRetainedBytes(bindings.policyDecisionDigest),
    wallTimeMs: planInputs.wallTimeMs,
    wallTimeOrigin: planInputs.wallTimeOrigin,
    outputSlot: planInputs.outputSlot,
    outputMaxJsonBytes: planInputs.outputMaxJsonBytes,
    expiresAt: bindings.expiresAt,
  }) as ExecutionPlanCandidateView;
  const encoded = encodeExecutionPlan(candidate);
  const digestBytes = Uint8Array.from(createHash("sha256").update(encoded).digest());
  const exactBytes = retainExactBytes(encoded);
  const digest = retainRoleBytes(digestBytes, "execution-plan");
  const plan = Object.freeze({
    candidate,
    exactBytes,
    digest,
    copyExactBytes: () => Uint8Array.from(encoded),
    copyDigestBytes: () => Uint8Array.from(digestBytes),
  });
  return Object.freeze({ ok: true, plan });
}

const TRUSTED_DIGEST_ROLES: ReadonlySet<TrustedExecutionPlanDigestRole> = new Set([
  "trust-epoch",
  "runtime-bundle-manifest",
  "profile-review-attestation",
  "profile-registry-entry",
  "backend-validation-record",
  "backend-configuration",
  "trust-snapshot",
  "policy-decision",
]);

const BINDING_INPUT_KEYS = Object.freeze([
  "backendConfigurationDigest",
  "backendValidationRecordDigest",
  "epochDigest",
  "epochSequence",
  "expiresAt",
  "installationId",
  "policyDecisionDigest",
  "profileRegistryEntryDigest",
  "profileReviewAttestationDigests",
  "runtimeBundleManifestDigest",
  "trustSnapshotDigest",
]);

function requireExactBindingInput(input: TrustedExecutionPlanBindingsInput): void {
  if (typeof input !== "object" || input === null || Array.isArray(input)) {
    throw new TypeError("trusted ExecutionPlan bindings must use the closed binding shape");
  }
  const keys = Reflect.ownKeys(input);
  if (
    keys.some((key) => typeof key !== "string") ||
    keys.length !== BINDING_INPUT_KEYS.length ||
    [...(keys as string[])].sort().some((key, index) => key !== BINDING_INPUT_KEYS[index])
  ) {
    throw new TypeError("trusted ExecutionPlan bindings must use the closed binding shape");
  }
}

function requireUint8Array(value: unknown, field: string): asserts value is Uint8Array {
  if (!(value instanceof Uint8Array)) {
    throw new TypeError(`${field} must be supplied as bytes`);
  }
}

function requireUInt53(value: number, field: string): void {
  if (!Number.isSafeInteger(value) || value < 0) {
    throw new TypeError(`${field} must be an unsigned safe integer`);
  }
}

function requireDigestRole<Role extends TrustedExecutionPlanDigestRole>(
  value: TrustedExecutionPlanDigest<Role>,
  expectedRole: Role,
): void {
  if (trustedDigests.get(value) !== expectedRole) {
    throw new TypeError(`trusted ${expectedRole} digest provenance is required`);
  }
}

function copyInstallationId(
  value: TrustedExecutionPlanInstallationId,
): TrustedExecutionPlanInstallationId {
  const copy = retainRoleBytes(
    Uint8Array.from(value.bytes),
    "installation-id",
  ) as TrustedExecutionPlanInstallationId;
  trustedInstallationIds.add(copy);
  return copy;
}

function copyDigest<Role extends TrustedExecutionPlanDigestRole>(
  value: TrustedExecutionPlanDigest<Role>,
  role: Role,
): TrustedExecutionPlanDigest<Role> {
  const copy = retainRoleBytes(
    Uint8Array.from(value.bytes),
    role,
  ) as TrustedExecutionPlanDigest<Role>;
  trustedDigests.set(copy, role);
  return copy;
}

function copyRetainedBytes<Role extends ExecutionPlanRetainedByteRole>(
  value: RetainedExecutionPlanBytes<Role>,
): RetainedExecutionPlanBytes<Role> {
  return retainRoleBytes(Uint8Array.from(value.bytes), value.role);
}

function retainedDigestFromHex<Role extends "source-manifest" | "inline-input">(
  value: string,
  role: Role,
): RetainedExecutionPlanBytes<Role> {
  if (!/^[0-9a-f]{64}$/u.test(value)) {
    throw new TypeError("retained Task 3B digest is outside the exact hexadecimal profile");
  }
  const bytes = new Uint8Array(32);
  for (let index = 0; index < bytes.byteLength; index += 1) {
    bytes[index] = Number.parseInt(value.slice(index * 2, index * 2 + 2), 16);
  }
  return retainRoleBytes(bytes, role);
}

function retainRoleBytes<Role extends ExecutionPlanRetainedByteRole>(
  value: Uint8Array,
  role: Role,
): RetainedExecutionPlanBytes<Role> {
  return Object.freeze({
    role,
    byteLength: value.byteLength,
    bytes: Object.freeze(Array.from(value)),
  }) as RetainedExecutionPlanBytes<Role>;
}

function retainExactBytes(value: Uint8Array): RetainedExactBytes {
  return Object.freeze({
    byteLength: value.byteLength,
    bytes: Object.freeze(Array.from(value)),
  });
}

function encodeExecutionPlan(candidate: ExecutionPlanCandidateView): Uint8Array {
  return concatenateBytes([
    encodeCborMapHeader(EXECUTION_PLAN_FIELD_COUNT),
    encodeCborUnsigned(1),
    encodeCborText(candidate.objectType),
    encodeCborUnsigned(2),
    encodeCborUnsigned(candidate.objectVersion),
    encodeCborUnsigned(3),
    encodeCborByteString(candidate.installationId.bytes),
    encodeCborUnsigned(4),
    encodeCborUnsigned(candidate.epochSequence),
    encodeCborUnsigned(5),
    encodeCborByteString(candidate.epochDigest.bytes),
    encodeCborUnsigned(6),
    encodeCborByteString(candidate.sourceManifestDigest.bytes),
    encodeCborUnsigned(7),
    encodeCborText(candidate.sourceEntrypoint),
    encodeCborUnsigned(8),
    encodeCborUnsigned(candidate.sourceByteLength),
    encodeCborUnsigned(9),
    encodeCborText(candidate.inputSlot),
    encodeCborUnsigned(10),
    encodeCborByteString(candidate.inlineInputDigest.bytes),
    encodeCborUnsigned(11),
    encodeCborUnsigned(candidate.inlineInputByteLength),
    encodeCborUnsigned(12),
    encodeCborText(candidate.runtimeProfileAlias),
    encodeCborUnsigned(13),
    encodeCborByteString(candidate.runtimeBundleManifestDigest.bytes),
    encodeCborUnsigned(14),
    encodeCborArrayHeader(candidate.profileReviewAttestationDigests.length),
    ...candidate.profileReviewAttestationDigests.map((digest) =>
      encodeCborByteString(digest.bytes),
    ),
    encodeCborUnsigned(15),
    encodeCborByteString(candidate.profileRegistryEntryDigest.bytes),
    encodeCborUnsigned(16),
    encodeCborByteString(candidate.backendValidationRecordDigest.bytes),
    encodeCborUnsigned(17),
    encodeCborByteString(candidate.backendConfigurationDigest.bytes),
    encodeCborUnsigned(18),
    encodeCborByteString(candidate.trustSnapshotDigest.bytes),
    encodeCborUnsigned(19),
    encodeCborByteString(candidate.policyDecisionDigest.bytes),
    encodeCborUnsigned(20),
    encodeCborUnsigned(candidate.wallTimeMs),
    encodeCborUnsigned(21),
    encodeCborText(candidate.wallTimeOrigin),
    encodeCborUnsigned(22),
    encodeCborText(candidate.outputSlot),
    encodeCborUnsigned(23),
    encodeCborUnsigned(candidate.outputMaxJsonBytes),
    encodeCborUnsigned(24),
    encodeCborUnsigned(candidate.expiresAt),
  ]);
}

function encodeCborUnsigned(value: number): Uint8Array {
  return encodeCborHead(0, value);
}

function encodeCborMapHeader(length: number): Uint8Array {
  return encodeCborHead(5, length);
}

function encodeCborArrayHeader(length: number): Uint8Array {
  return encodeCborHead(4, length);
}

function encodeCborText(value: string): Uint8Array {
  const bytes = new TextEncoder().encode(value);
  return concatenateBytes([encodeCborHead(3, bytes.byteLength), bytes]);
}

function encodeCborByteString(value: readonly number[]): Uint8Array {
  const bytes = Uint8Array.from(value);
  return concatenateBytes([encodeCborHead(2, bytes.byteLength), bytes]);
}

function encodeCborHead(majorType: number, value: number): Uint8Array {
  requireUInt53(value, "CBOR argument");
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

function concatenateBytes(chunks: readonly Uint8Array[]): Uint8Array {
  const byteLength = chunks.reduce((total, chunk) => total + chunk.byteLength, 0);
  const bytes = new Uint8Array(byteLength);
  let offset = 0;
  for (const chunk of chunks) {
    bytes.set(chunk, offset);
    offset += chunk.byteLength;
  }
  return bytes;
}

function rejected(code: ExecutionPlanConstructionRefusalCode): {
  readonly ok: false;
  readonly refusal: ExecutionPlanConstructionRefusal;
} {
  return Object.freeze({
    ok: false,
    refusal: Object.freeze({
      owner: "execution-plan-builder",
      classification: "BINDING",
      code,
    }),
  });
}
