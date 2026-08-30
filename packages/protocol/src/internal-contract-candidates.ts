/**
 * Decoded views of passive internal deterministic-CBOR candidates.
 *
 * These types are not wire codecs. Exact received canonical bytes remain the
 * authoritative object and no product component consumes these candidates yet.
 */

declare const candidateScalarBrand: unique symbol;

/** A nonzero 16-byte installation identity, produced only by {@link asCandidateInstallationId}. */
export type CandidateInstallationId = Readonly<Uint8Array> & {
  readonly [candidateScalarBrand]: "installation-id:16";
};
/** A nonzero 16-byte registration identity, produced only by {@link asCandidateRegistrationId}. */
export type CandidateRegistrationId = Readonly<Uint8Array> & {
  readonly [candidateScalarBrand]: "registration-id:16";
};
/** A nonzero 16-byte Supervisor identity, produced only by {@link asCandidateSupervisorId}. */
export type CandidateSupervisorId = Readonly<Uint8Array> & {
  readonly [candidateScalarBrand]: "supervisor-id:16";
};

/** The distinct SHA-256 digest roles a candidate object may bind. Each role is a separate nominal type. */
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

/** A role-tagged 32-byte SHA-256 digest, produced only by {@link asCandidateDigest}. Every 32-byte bit pattern is structurally accepted, including all-zero; existence and byte equality are binding checks owned by the role-specific resolver (ADR-0023). */
export type CandidateDigest<Role extends CandidateDigestRole> = Readonly<Uint8Array> & {
  readonly [candidateScalarBrand]: `sha256:${Role}:32`;
};

/** An unsigned safe integer (0 through 2^53-1), produced only by {@link asCandidateUInt53}. */
export type CandidateUInt53 = number & {
  readonly [candidateScalarBrand]: "uint53";
};
/** A positive (nonzero) safe integer, produced only by {@link asCandidatePositiveUInt53}. */
export type CandidatePositiveUInt53 = number & {
  readonly [candidateScalarBrand]: "positive-uint53";
};

/**
 * Decoded view of the Supervisor-issued `PlanRegistration` candidate object.
 * Carries no replacement plan bytes — registration binds a caller's exact
 * submitted `ExecutionPlan` bytes to an ID/sequence/expiry, it never
 * substitutes new plan content (ADR-0011, ADR-0023).
 */
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

/** Validates and defensively copies a 16-byte nonzero installation ID. */
export function asCandidateInstallationId(value: Uint8Array): CandidateInstallationId {
  return copyCandidateId(value, "installation") as CandidateInstallationId;
}

/** Validates and defensively copies a 16-byte nonzero registration ID. */
export function asCandidateRegistrationId(value: Uint8Array): CandidateRegistrationId {
  return copyCandidateId(value, "registration") as CandidateRegistrationId;
}

/** Validates and defensively copies a 16-byte nonzero Supervisor ID. */
export function asCandidateSupervisorId(value: Uint8Array): CandidateSupervisorId {
  return copyCandidateId(value, "Supervisor") as CandidateSupervisorId;
}

/**
 * Validates and defensively copies a 32-byte SHA-256 digest for the given
 * role. Unlike the ID constructors above, the all-zero value is accepted —
 * digests are structurally width-checked only; existence/equality binding is
 * a separate resolver concern (ADR-0023).
 */
export function asCandidateDigest<Role extends CandidateDigestRole>(
  value: Uint8Array,
  role: Role,
): CandidateDigest<Role> {
  if (value.length !== 32) {
    throw new TypeError(`${role} SHA-256 digest must contain exactly 32 bytes`);
  }
  return Uint8Array.from(value) as CandidateDigest<Role>;
}

/** Validates an unsigned safe integer (0 through 2^53-1). */
export function asCandidateUInt53(value: number): CandidateUInt53 {
  if (!Number.isSafeInteger(value) || value < 0) {
    throw new TypeError("value must be an unsigned safe integer");
  }
  return value as CandidateUInt53;
}

/** Validates a positive (nonzero) safe integer. */
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
