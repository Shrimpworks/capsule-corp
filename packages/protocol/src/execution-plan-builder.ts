import { createHash } from "node:crypto";

import {
  concatenateBytes,
  encodeCborArrayHeader,
  encodeCborByteString,
  encodeCborMapHeader,
  encodeCborText,
  encodeCborUnsigned,
} from "./internal-cbor-primitives.js";
import type { ResolvedJobProposalPlanInputs, RetainedExactBytes } from "./job-proposal-resolver.js";
import { isRetainedResolvedJobProposalPlanInputs } from "./resolved-job-proposal-provenance.js";

const EXECUTION_PLAN_OBJECT_TYPE = "capsule.execution-plan" as const;
const EXECUTION_PLAN_OBJECT_VERSION = 0 as const;
const EXECUTION_PLAN_FIELD_COUNT = 24;
const MAX_REVIEW_ATTESTATIONS = 8;

declare const executionPlanScalarBrand: unique symbol;
declare const trustedExecutionPlanBindingBrand: unique symbol;
declare const trustedExecutionPlanRegistrationRoleBindingsBrand: unique symbol;

/** Digest roles a trusted (non-proposal-owned) ExecutionPlan binding may carry. */
export type TrustedExecutionPlanDigestRole =
  | "trust-epoch"
  | "source-manifest"
  | "inline-input"
  | "runtime-bundle-manifest"
  | "profile-review-attestation"
  | "profile-registry-entry"
  | "backend-validation-record"
  | "backend-configuration"
  | "trust-snapshot"
  | "policy-decision";

/** {@link TrustedExecutionPlanDigestRole} plus the two non-digest byte roles retained on a built plan. */
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

/** A role-tagged defensive byte copy retained on a plan or binding value. */
export interface RetainedExecutionPlanBytes<Role extends ExecutionPlanRetainedByteRole> {
  readonly role: Role;
  readonly byteLength: number;
  readonly bytes: readonly number[];
  readonly [executionPlanScalarBrand]: Role;
}

/** A nonzero 16-byte installation ID, produced only by {@link createTrustedExecutionPlanInstallationId}. */
export type TrustedExecutionPlanInstallationId = RetainedExecutionPlanBytes<"installation-id"> & {
  readonly [trustedExecutionPlanBindingBrand]: "installation-id";
};

/** A role-tagged 32-byte digest, produced only by {@link createTrustedExecutionPlanDigest} for that exact role. */
export type TrustedExecutionPlanDigest<Role extends TrustedExecutionPlanDigestRole> =
  RetainedExecutionPlanBytes<Role> & {
    readonly [trustedExecutionPlanBindingBrand]: Role;
  };

/** Caller-supplied input to {@link createTrustedExecutionPlanBindings}, before provenance/shape validation. */
export interface TrustedExecutionPlanBindingsInput {
  readonly installationId: TrustedExecutionPlanInstallationId;
  readonly epochSequence: number;
  readonly epochDigest: TrustedExecutionPlanDigest<"trust-epoch">;
  readonly runtimeBundleManifestDigest: TrustedExecutionPlanDigest<"runtime-bundle-manifest">;
  /** An unordered set of independent per-reviewer verdicts; supply in any order (see profileReviewAttestationDigestsMatch). */
  readonly profileReviewAttestationDigests: readonly TrustedExecutionPlanDigest<"profile-review-attestation">[];
  readonly profileRegistryEntryDigest: TrustedExecutionPlanDigest<"profile-registry-entry">;
  readonly backendValidationRecordDigest: TrustedExecutionPlanDigest<"backend-validation-record">;
  readonly backendConfigurationDigest: TrustedExecutionPlanDigest<"backend-configuration">;
  readonly trustSnapshotDigest: TrustedExecutionPlanDigest<"trust-snapshot">;
  readonly policyDecisionDigest: TrustedExecutionPlanDigest<"policy-decision">;
  readonly expiresAt: number;
}

/**
 * The trust-owned (non-proposal) half of an ExecutionPlan's bindings,
 * produced only by {@link createTrustedExecutionPlanBindings}: installation
 * identity, trust epoch, runtime bundle/profile/backend/policy digests, and
 * expiry. Combined with {@link ResolvedJobProposalPlanInputs} by
 * {@link constructExecutionPlan}.
 */
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

/** Caller-supplied input to {@link createTrustedExecutionPlanRegistrationRoleBindings}, before provenance/shape validation. */
export interface TrustedExecutionPlanRegistrationRoleBindingsInput {
  readonly installationId: TrustedExecutionPlanInstallationId;
  readonly epochDigest: TrustedExecutionPlanDigest<"trust-epoch">;
  readonly sourceManifestDigest: TrustedExecutionPlanDigest<"source-manifest">;
  readonly inlineInputDigest: TrustedExecutionPlanDigest<"inline-input">;
  readonly runtimeBundleManifestDigest: TrustedExecutionPlanDigest<"runtime-bundle-manifest">;
  /** An unordered set of independent per-reviewer verdicts; supply in any order (see profileReviewAttestationDigestsMatch). */
  readonly profileReviewAttestationDigests: readonly TrustedExecutionPlanDigest<"profile-review-attestation">[];
  readonly profileRegistryEntryDigest: TrustedExecutionPlanDigest<"profile-registry-entry">;
  readonly backendValidationRecordDigest: TrustedExecutionPlanDigest<"backend-validation-record">;
  readonly backendConfigurationDigest: TrustedExecutionPlanDigest<"backend-configuration">;
  readonly trustSnapshotDigest: TrustedExecutionPlanDigest<"trust-snapshot">;
  readonly policyDecisionDigest: TrustedExecutionPlanDigest<"policy-decision">;
}

/**
 * The role-binding subset a Supervisor registration checks a submitted plan
 * against, produced only by
 * {@link createTrustedExecutionPlanRegistrationRoleBindings}. Deliberately
 * narrower than {@link TrustedExecutionPlanBindings}: no epoch sequence or
 * expiry, since registration binds identity/content roles, not time.
 */
export interface TrustedExecutionPlanRegistrationRoleBindings {
  readonly installationId: TrustedExecutionPlanInstallationId;
  readonly epochDigest: TrustedExecutionPlanDigest<"trust-epoch">;
  readonly sourceManifestDigest: TrustedExecutionPlanDigest<"source-manifest">;
  readonly inlineInputDigest: TrustedExecutionPlanDigest<"inline-input">;
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
  readonly [trustedExecutionPlanRegistrationRoleBindingsBrand]: "registration-role-bindings";
}

/** The decoded field-by-field view of a constructed ExecutionPlan, mirroring the CBOR map's fields (ADR-0019). */
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

/** The output of a successful {@link constructExecutionPlan} call: the decoded view plus its exact canonical bytes and digest. */
export interface ConstructedExecutionPlan {
  readonly candidate: ExecutionPlanCandidateView;
  readonly exactBytes: RetainedExactBytes;
  readonly digest: RetainedExecutionPlanBytes<"execution-plan">;
  copyExactBytes(): Uint8Array;
  copyDigestBytes(): Uint8Array;
}

/** The inert TypeScript-side handoff produced by {@link prepareExecutionPlanRegistrationHandoff}: exact plan bytes plus the role bindings a Supervisor registration would check them against. */
export interface ExecutionPlanRegistrationHandoff {
  readonly exactPlanByteLength: number;
  readonly roleBindings: TrustedExecutionPlanRegistrationRoleBindings;
  copyExactPlanBytes(): Uint8Array;
}

/** The exhaustive set of refusal codes {@link prepareExecutionPlanRegistrationHandoff} can return. */
export type ExecutionPlanRegistrationHandoffRefusalCode =
  | "PLAN_PROVENANCE"
  | "ROLE_BINDING_PROVENANCE"
  | "ROLE_BINDING_MISMATCH"
  | "PLAN_DIGEST_MISMATCH";

export interface ExecutionPlanRegistrationHandoffRefusal {
  readonly owner: "execution-plan-registration-handoff";
  readonly classification: "BINDING";
  readonly code: ExecutionPlanRegistrationHandoffRefusalCode;
}

/** The outcome of {@link prepareExecutionPlanRegistrationHandoff}: either the handoff or a classified refusal. */
export type ExecutionPlanRegistrationHandoffResult =
  | { readonly ok: true; readonly handoff: ExecutionPlanRegistrationHandoff }
  | { readonly ok: false; readonly refusal: ExecutionPlanRegistrationHandoffRefusal };

/** The exhaustive set of refusal codes {@link constructExecutionPlan} can return. */
export type ExecutionPlanConstructionRefusalCode =
  | "PLAN_INPUT_PROVENANCE"
  | "TRUSTED_BINDING_PROVENANCE";

export interface ExecutionPlanConstructionRefusal {
  readonly owner: "execution-plan-builder";
  readonly classification: "BINDING";
  readonly code: ExecutionPlanConstructionRefusalCode;
}

/** The outcome of {@link constructExecutionPlan}: either the constructed plan or a classified refusal. */
export type ExecutionPlanConstructionResult =
  | { readonly ok: true; readonly plan: ConstructedExecutionPlan }
  | { readonly ok: false; readonly refusal: ExecutionPlanConstructionRefusal };

const trustedInstallationIds = new WeakSet<object>();
const trustedDigests = new WeakMap<object, TrustedExecutionPlanDigestRole>();
const trustedBindings = new WeakSet<object>();
const trustedRegistrationRoleBindings = new WeakSet<object>();
const constructedPlans = new WeakSet<object>();

/** Validates and defensively copies a 16-byte nonzero installation ID for use in trusted ExecutionPlan bindings. */
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

/** Validates and defensively copies a 32-byte digest for the given trusted role. */
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

/**
 * Validates a closed-shape {@link TrustedExecutionPlanBindingsInput} (exact
 * key set, correct provenance on every nested ID/digest, one-to-eight
 * review attestations) and freezes it into {@link TrustedExecutionPlanBindings}.
 * Every field is defensively copied — the caller's original input objects
 * are never aliased into the returned value.
 */
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
 * Validates a closed-shape {@link TrustedExecutionPlanRegistrationRoleBindingsInput}
 * and freezes it into {@link TrustedExecutionPlanRegistrationRoleBindings},
 * the same way {@link createTrustedExecutionPlanBindings} does for the
 * fuller binding set.
 */
export function createTrustedExecutionPlanRegistrationRoleBindings(
  input: TrustedExecutionPlanRegistrationRoleBindingsInput,
): TrustedExecutionPlanRegistrationRoleBindings {
  requireExactRegistrationRoleBindingInput(input);
  if (!trustedInstallationIds.has(input.installationId)) {
    throw new TypeError("trusted installation ID provenance is required");
  }
  requireDigestRole(input.epochDigest, "trust-epoch");
  requireDigestRole(input.sourceManifestDigest, "source-manifest");
  requireDigestRole(input.inlineInputDigest, "inline-input");
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

  const bindings = Object.freeze({
    installationId: copyInstallationId(input.installationId),
    epochDigest: copyDigest(input.epochDigest, "trust-epoch"),
    sourceManifestDigest: copyDigest(input.sourceManifestDigest, "source-manifest"),
    inlineInputDigest: copyDigest(input.inlineInputDigest, "inline-input"),
    runtimeBundleManifestDigest: copyDigest(
      input.runtimeBundleManifestDigest,
      "runtime-bundle-manifest",
    ),
    profileReviewAttestationDigests: Object.freeze(
      input.profileReviewAttestationDigests.map((digest) =>
        copyDigest(digest, "profile-review-attestation"),
      ),
    ) as TrustedExecutionPlanRegistrationRoleBindings["profileReviewAttestationDigests"],
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
  }) as TrustedExecutionPlanRegistrationRoleBindings;
  trustedRegistrationRoleBindings.add(bindings);
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
  constructedPlans.add(plan);
  return Object.freeze({ ok: true, plan });
}

/**
 * Prepares the inert TypeScript side of the unwired registration handoff. The
 * returned object contains only a private defensive copy of the exact plan
 * bytes and a copied complete role-binding set. It performs no I/O and grants
 * no registration, approval, attempt, backend, or guest authority.
 */
export function prepareExecutionPlanRegistrationHandoff(
  plan: ConstructedExecutionPlan,
  roleBindings: TrustedExecutionPlanRegistrationRoleBindings,
): ExecutionPlanRegistrationHandoffResult {
  if (!constructedPlans.has(plan)) {
    return rejectedHandoff("PLAN_PROVENANCE");
  }
  if (!trustedRegistrationRoleBindings.has(roleBindings)) {
    return rejectedHandoff("ROLE_BINDING_PROVENANCE");
  }
  const exactPlanBytes = plan.copyExactBytes();
  const recomputedDigest = Uint8Array.from(createHash("sha256").update(exactPlanBytes).digest());
  if (!equalBytes(recomputedDigest, plan.copyDigestBytes())) {
    return rejectedHandoff("PLAN_DIGEST_MISMATCH");
  }
  if (!registrationRoleBindingsMatchPlan(roleBindings, plan.candidate)) {
    return rejectedHandoff("ROLE_BINDING_MISMATCH");
  }
  const copiedRoleBindings = copyRegistrationRoleBindings(roleBindings);
  const handoff = Object.freeze({
    exactPlanByteLength: exactPlanBytes.byteLength,
    roleBindings: copiedRoleBindings,
    copyExactPlanBytes: () => Uint8Array.from(exactPlanBytes),
  });
  return Object.freeze({ ok: true, handoff });
}

const TRUSTED_DIGEST_ROLES: ReadonlySet<TrustedExecutionPlanDigestRole> = new Set([
  "trust-epoch",
  "source-manifest",
  "inline-input",
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

const REGISTRATION_ROLE_BINDING_INPUT_KEYS = Object.freeze([
  "backendConfigurationDigest",
  "backendValidationRecordDigest",
  "epochDigest",
  "inlineInputDigest",
  "installationId",
  "policyDecisionDigest",
  "profileRegistryEntryDigest",
  "profileReviewAttestationDigests",
  "runtimeBundleManifestDigest",
  "sourceManifestDigest",
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

function requireExactRegistrationRoleBindingInput(
  input: TrustedExecutionPlanRegistrationRoleBindingsInput,
): void {
  if (typeof input !== "object" || input === null || Array.isArray(input)) {
    throw new TypeError("trusted registration role bindings must use the closed binding shape");
  }
  const keys = Reflect.ownKeys(input);
  if (
    keys.some((key) => typeof key !== "string") ||
    keys.length !== REGISTRATION_ROLE_BINDING_INPUT_KEYS.length ||
    [...(keys as string[])]
      .sort()
      .some((key, index) => key !== REGISTRATION_ROLE_BINDING_INPUT_KEYS[index])
  ) {
    throw new TypeError("trusted registration role bindings must use the closed binding shape");
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

function copyRegistrationRoleBindings(
  bindings: TrustedExecutionPlanRegistrationRoleBindings,
): TrustedExecutionPlanRegistrationRoleBindings {
  const copy = Object.freeze({
    installationId: copyInstallationId(bindings.installationId),
    epochDigest: copyDigest(bindings.epochDigest, "trust-epoch"),
    sourceManifestDigest: copyDigest(bindings.sourceManifestDigest, "source-manifest"),
    inlineInputDigest: copyDigest(bindings.inlineInputDigest, "inline-input"),
    runtimeBundleManifestDigest: copyDigest(
      bindings.runtimeBundleManifestDigest,
      "runtime-bundle-manifest",
    ),
    profileReviewAttestationDigests: Object.freeze(
      bindings.profileReviewAttestationDigests.map((digest) =>
        copyDigest(digest, "profile-review-attestation"),
      ),
    ) as TrustedExecutionPlanRegistrationRoleBindings["profileReviewAttestationDigests"],
    profileRegistryEntryDigest: copyDigest(
      bindings.profileRegistryEntryDigest,
      "profile-registry-entry",
    ),
    backendValidationRecordDigest: copyDigest(
      bindings.backendValidationRecordDigest,
      "backend-validation-record",
    ),
    backendConfigurationDigest: copyDigest(
      bindings.backendConfigurationDigest,
      "backend-configuration",
    ),
    trustSnapshotDigest: copyDigest(bindings.trustSnapshotDigest, "trust-snapshot"),
    policyDecisionDigest: copyDigest(bindings.policyDecisionDigest, "policy-decision"),
  }) as TrustedExecutionPlanRegistrationRoleBindings;
  trustedRegistrationRoleBindings.add(copy);
  return copy;
}

function registrationRoleBindingsMatchPlan(
  bindings: TrustedExecutionPlanRegistrationRoleBindings,
  candidate: ExecutionPlanCandidateView,
): boolean {
  return (
    equalRetainedBytes(bindings.installationId, candidate.installationId) &&
    equalRetainedBytes(bindings.epochDigest, candidate.epochDigest) &&
    equalRetainedBytes(bindings.sourceManifestDigest, candidate.sourceManifestDigest) &&
    equalRetainedBytes(bindings.inlineInputDigest, candidate.inlineInputDigest) &&
    equalRetainedBytes(
      bindings.runtimeBundleManifestDigest,
      candidate.runtimeBundleManifestDigest,
    ) &&
    bindings.profileReviewAttestationDigests.length ===
      candidate.profileReviewAttestationDigests.length &&
    profileReviewAttestationDigestsMatch(
      bindings.profileReviewAttestationDigests,
      candidate.profileReviewAttestationDigests,
    ) &&
    equalRetainedBytes(bindings.profileRegistryEntryDigest, candidate.profileRegistryEntryDigest) &&
    equalRetainedBytes(
      bindings.backendValidationRecordDigest,
      candidate.backendValidationRecordDigest,
    ) &&
    equalRetainedBytes(bindings.backendConfigurationDigest, candidate.backendConfigurationDigest) &&
    equalRetainedBytes(bindings.trustSnapshotDigest, candidate.trustSnapshotDigest) &&
    equalRetainedBytes(bindings.policyDecisionDigest, candidate.policyDecisionDigest)
  );
}

/**
 * profileReviewAttestationDigests holds independent per-reviewer verdicts
 * about the same bundle (docs/TECHNICAL_DESIGN.md: "ProfileReviewAttestation:
 * signed reviewer claim..."; docs/protocol/OBJECT_MODEL.md: "Independent
 * verdict for exact bundle"), not a priority- or sequence-ordered list --
 * nothing in ADR-0019, ADR-0023, or the object model gives review order any
 * meaning. `bindings` and `candidate` are also built independently (trusted
 * Supervisor-side role bindings vs. the submitted plan's own candidate
 * view), so their array order has no guaranteed relationship to begin with.
 * Compare as sets, sorted by digest bytes, instead of by position, so two
 * callers supplying the same attestation set in a different order do not
 * spuriously refuse registration.
 */
function profileReviewAttestationDigestsMatch(
  left: readonly RetainedExecutionPlanBytes<"profile-review-attestation">[],
  right: readonly RetainedExecutionPlanBytes<"profile-review-attestation">[],
): boolean {
  const sortedLeft = [...left].sort(compareRetainedBytes);
  const sortedRight = [...right].sort(compareRetainedBytes);
  return sortedLeft.every((digest, index) => {
    const other = sortedRight[index];
    return other !== undefined && equalRetainedBytes(digest, other);
  });
}

function compareRetainedBytes(
  left: RetainedExecutionPlanBytes<ExecutionPlanRetainedByteRole>,
  right: RetainedExecutionPlanBytes<ExecutionPlanRetainedByteRole>,
): number {
  const length = Math.min(left.bytes.length, right.bytes.length);
  for (let index = 0; index < length; index += 1) {
    const difference = (left.bytes[index] ?? 0) - (right.bytes[index] ?? 0);
    if (difference !== 0) {
      return difference;
    }
  }
  return left.bytes.length - right.bytes.length;
}

function equalRetainedBytes(
  left: RetainedExecutionPlanBytes<ExecutionPlanRetainedByteRole>,
  right: RetainedExecutionPlanBytes<ExecutionPlanRetainedByteRole>,
): boolean {
  return equalBytes(Uint8Array.from(left.bytes), Uint8Array.from(right.bytes));
}

function equalBytes(left: Uint8Array, right: Uint8Array): boolean {
  return (
    left.byteLength === right.byteLength && left.every((value, index) => value === right[index])
  );
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

function rejectedHandoff(code: ExecutionPlanRegistrationHandoffRefusalCode): {
  readonly ok: false;
  readonly refusal: ExecutionPlanRegistrationHandoffRefusal;
} {
  return Object.freeze({
    ok: false,
    refusal: Object.freeze({
      owner: "execution-plan-registration-handoff",
      classification: "BINDING",
      code,
    }),
  });
}
