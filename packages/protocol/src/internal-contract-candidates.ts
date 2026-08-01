/**
 * Decoded views of passive internal deterministic-CBOR candidates.
 *
 * These types are not wire codecs. Exact received canonical bytes remain the
 * authoritative object and no product component consumes these candidates yet.
 */

declare const candidateScalarBrand: unique symbol;

export type CandidateInstallationId = Readonly<Uint8Array> & {
  readonly [candidateScalarBrand]: "installation-id:16";
};
export type CandidateRegistrationId = Readonly<Uint8Array> & {
  readonly [candidateScalarBrand]: "registration-id:16";
};
export type CandidateSupervisorId = Readonly<Uint8Array> & {
  readonly [candidateScalarBrand]: "supervisor-id:16";
};

export type CandidateDigestRole =
  | "backend-configuration"
  | "backend-validation-record"
  | "execution-plan"
  | "inline-input"
  | "policy-decision"
  | "profile-registry-entry"
  | "profile-review-attestation"
  | "runtime-bundle-manifest"
  | "source-manifest"
  | "trust-epoch"
  | "trust-snapshot";

export type CandidateDigest<Role extends CandidateDigestRole> = Readonly<Uint8Array> & {
  readonly [candidateScalarBrand]: `sha256:${Role}:32`;
};

export type CandidateUInt53 = number & {
  readonly [candidateScalarBrand]: "uint53";
};
export type CandidatePositiveUInt53 = number & {
  readonly [candidateScalarBrand]: "positive-uint53";
};

export interface ExecutionPlanCandidate {
  objectType: "capsule.execution-plan";
  objectVersion: 0;
  installationId: CandidateInstallationId;
  epochSequence: CandidateUInt53;
  epochDigest: CandidateDigest<"trust-epoch">;
  sourceManifestDigest: CandidateDigest<"source-manifest">;
  sourceEntrypoint: string;
  sourceByteLength: CandidateUInt53;
  inputSlot: "primary-data";
  inlineInputDigest: CandidateDigest<"inline-input">;
  inlineInputByteLength: CandidateUInt53;
  runtimeProfileAlias: string;
  runtimeBundleManifestDigest: CandidateDigest<"runtime-bundle-manifest">;
  profileReviewAttestationDigests: readonly [
    CandidateDigest<"profile-review-attestation">,
    ...CandidateDigest<"profile-review-attestation">[],
  ];
  profileRegistryEntryDigest: CandidateDigest<"profile-registry-entry">;
  backendValidationRecordDigest: CandidateDigest<"backend-validation-record">;
  backendConfigurationDigest: CandidateDigest<"backend-configuration">;
  trustSnapshotDigest: CandidateDigest<"trust-snapshot">;
  policyDecisionDigest: CandidateDigest<"policy-decision">;
  wallTimeMs: CandidatePositiveUInt53;
  wallTimeOrigin: "requested" | "trusted-default";
  outputSlot: "transformed-json";
  outputMaxJsonBytes: CandidatePositiveUInt53;
  expiresAt: CandidateUInt53;
}

export interface PlanRegistrationCandidate {
  objectType: "capsule.plan-registration";
  objectVersion: 0;
  registrationId: CandidateRegistrationId;
  registrationSequence: CandidatePositiveUInt53;
  planDigest: CandidateDigest<"execution-plan">;
  installationId: CandidateInstallationId;
  epochSequence: CandidateUInt53;
  epochDigest: CandidateDigest<"trust-epoch">;
  supervisorId: CandidateSupervisorId;
  expiresAt: CandidateUInt53;
}

export function asCandidateInstallationId(value: Uint8Array): CandidateInstallationId {
  return copyCandidateId(value, "installation") as CandidateInstallationId;
}

export function asCandidateRegistrationId(value: Uint8Array): CandidateRegistrationId {
  return copyCandidateId(value, "registration") as CandidateRegistrationId;
}

export function asCandidateSupervisorId(value: Uint8Array): CandidateSupervisorId {
  return copyCandidateId(value, "Supervisor") as CandidateSupervisorId;
}

export function asCandidateDigest<Role extends CandidateDigestRole>(
  value: Uint8Array,
  role: Role,
): CandidateDigest<Role> {
  if (value.length !== 32) {
    throw new TypeError(`${role} SHA-256 digest must contain exactly 32 bytes`);
  }
  return Uint8Array.from(value) as CandidateDigest<Role>;
}

export function asCandidateUInt53(value: number): CandidateUInt53 {
  if (!Number.isSafeInteger(value) || value < 0) {
    throw new TypeError("value must be an unsigned safe integer");
  }
  return value as CandidateUInt53;
}

export function asCandidatePositiveUInt53(value: number): CandidatePositiveUInt53 {
  if (!Number.isSafeInteger(value) || value <= 0) {
    throw new TypeError("value must be a positive safe integer");
  }
  return value as CandidatePositiveUInt53;
}

function copyCandidateId(value: Uint8Array, role: string): Readonly<Uint8Array> {
  if (value.length !== 16) {
    throw new TypeError(`${role} ID must contain exactly 16 bytes`);
  }
  if (value.every((byte) => byte === 0)) {
    throw new TypeError(`${role} ID must be nonzero`);
  }
  return Uint8Array.from(value);
}
