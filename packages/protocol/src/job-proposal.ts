/**
 * Passive Phase 2 candidate types. No daemon, SDK, or MCP endpoint accepts this
 * object until the coordinated protocol cutover is complete.
 */

export const JOB_PROPOSAL_API_VERSION = "capsule.dev/v0" as const;
export const JOB_PROPOSAL_KIND = "JobProposal" as const;
export const PRIMARY_DATA_INPUT_SLOT = "primary-data" as const;
export const TRANSFORMED_JSON_OUTPUT_SLOT = "transformed-json" as const;

declare const sourcePathBrand: unique symbol;
declare const sourceEntrypointBrand: unique symbol;
declare const runtimeProfileAliasBrand: unique symbol;
declare const positiveSafeIntegerBrand: unique symbol;
declare const safeJsonIntegerBrand: unique symbol;
declare const decodedJobProposalCandidateBrand: unique symbol;

/** A relative source-file path accepted by {@link asSourcePath}'s candidate grammar. */
export type SourcePath = string & { readonly [sourcePathBrand]: "SourcePath" };
/** A relative source entrypoint filename accepted by {@link asSourceEntrypoint}'s candidate grammar. */
export type SourceEntrypoint = string & {
  readonly [sourceEntrypointBrand]: "SourceEntrypoint";
};
/** A `name@version` runtime profile alias accepted by {@link asRuntimeProfileAlias}. */
export type RuntimeProfileAlias = string & {
  readonly [runtimeProfileAliasBrand]: "RuntimeProfileAlias";
};
/** A positive (nonzero) safe integer, produced only by {@link asPositiveSafeInteger}. */
export type PositiveSafeInteger = number & {
  readonly [positiveSafeIntegerBrand]: "PositiveSafeInteger";
};
/** A safe integer usable as an inline JSON number, produced only by {@link asSafeJsonInteger}. */
export type SafeJsonInteger = number & {
  readonly [safeJsonIntegerBrand]: "SafeJsonInteger";
};

/**
 * A JSON value restricted to the candidate's supported scalar/container
 * shapes. Numbers are safe integers only; floats, `NaN`, and `Infinity` are
 * not representable.
 */
export type InlineJsonValue =
  | null
  | boolean
  | SafeJsonInteger
  | string
  | readonly InlineJsonValue[]
  | { readonly [key: string]: InlineJsonValue };

/** The agent-proposed source bundle: one entrypoint plus its file contents, keyed by path. */
export interface SourceBundleProposal {
  readonly entrypoint: SourceEntrypoint;
  readonly files: Readonly<Record<SourcePath, string>>;
}

/** The sole supported v0 input: inline JSON bound to the fixed `primary-data` slot. */
export interface InlineJsonInputProposal {
  readonly slot: typeof PRIMARY_DATA_INPUT_SLOT;
  readonly kind: "inline-json";
  readonly value: InlineJsonValue;
}

/**
 * v0 exposes only the independently supported wall-time request (ADR-0017).
 * An omitted value resolves to the trusted profile/policy default; see
 * resolveJobProposal's WALL_TIME_CEILING/WALL_TIME_PROFILE_RANGE checks.
 */
export interface RequestedLimits {
  readonly wallTimeMs?: PositiveSafeInteger;
}

/**
 * `maxBytes` is a required structural bound only. Unlike `wallTimeMs`
 * (see resolveJobProposal's WALL_TIME_CEILING check), it is not yet
 * checked against a trusted policy ceiling anywhere in the resolver.
 * This is a deliberately deferred gap, not an oversight: ADR-0023 and
 * PHASE_2A_CONTRACT_FOUNDATION.md defer output-payload limits pending
 * P0 transport evidence. Activating enforcement requires adding an
 * `outputMaxJsonBytes` ceiling dimension to TrustedUserPolicyContext and
 * an ADR-0023 addendum recording the mechanism and evidence — do not
 * wire a fixed schema maximum here as a shortcut.
 */
export interface InlineJsonOutputProposal {
  readonly slot: typeof TRANSFORMED_JSON_OUTPUT_SLOT;
  readonly kind: "inline-json";
  readonly maxBytes: PositiveSafeInteger;
}

/**
 * The full passive v0 candidate proposal shape. Omits network, process,
 * environment, native/FFI, package, image/backend, and arbitrary path
 * authority entirely (ADR-0017) — there is no false-valued placeholder for
 * unsupported power; those fields simply do not exist here.
 */
export interface JobProposal {
  readonly apiVersion: typeof JOB_PROPOSAL_API_VERSION;
  readonly kind: typeof JOB_PROPOSAL_KIND;
  readonly source: SourceBundleProposal;
  readonly runtimeProfile: RuntimeProfileAlias;
  readonly input: InlineJsonInputProposal;
  readonly requestedLimits: RequestedLimits;
  readonly outputs: readonly [InlineJsonOutputProposal];
  readonly labels?: Readonly<Record<string, string>>;
}

/**
 * A passive candidate produced only after the Task 3A raw and closed-schema
 * boundaries accept caller-owned bytes. It is not semantically resolved
 * authority and cannot be used as an ExecutionPlan.
 */
export type DecodedJobProposalCandidate = JobProposal & {
  readonly [decodedJobProposalCandidateBrand]: "DecodedJobProposalCandidate";
};

const SOURCE_PATH_PATTERN =
  /^[A-Za-z0-9][A-Za-z0-9._-]{0,63}(?:\/[A-Za-z0-9][A-Za-z0-9._-]{0,63})*$/u;
const SOURCE_ENTRYPOINT_PATTERN =
  /^(?:[A-Za-z0-9][A-Za-z0-9._-]{0,63}\/)*[A-Za-z0-9][A-Za-z0-9._-]{0,55}\.(?:js|mjs|cjs|ts|mts|cts)$/u;
const RUNTIME_PROFILE_ALIAS_PATTERN = /^[a-z][a-z0-9-]*@[1-9][0-9]*$/u;

/**
 * Validates an ASCII-only relative source path. Every path segment must
 * start with an alphanumeric character, which structurally blocks `.`/`..`
 * traversal segments without needing platform filesystem semantics
 * (ADR-0023: "Canonicalize Unicode source paths" was rejected for this
 * reason). Capsule never materializes these as live host paths.
 */
export function asSourcePath(value: string): SourcePath {
  if (value.length > 256 || !SOURCE_PATH_PATTERN.test(value)) {
    throw new TypeError("source path is outside the candidate grammar");
  }
  return value as SourcePath;
}

/** Validates a relative entrypoint path ending in a supported JS/TS extension. */
export function asSourceEntrypoint(value: string): SourceEntrypoint {
  if (value.length > 256 || !SOURCE_ENTRYPOINT_PATTERN.test(value)) {
    throw new TypeError("source entrypoint is outside the candidate grammar");
  }
  return value as SourceEntrypoint;
}

/** Validates a `name@version` runtime profile alias (e.g. `bun-data@1`). */
export function asRuntimeProfileAlias(value: string): RuntimeProfileAlias {
  if (value.length > 128 || !RUNTIME_PROFILE_ALIAS_PATTERN.test(value)) {
    throw new TypeError("runtime profile alias is outside the candidate grammar");
  }
  return value as RuntimeProfileAlias;
}

/** Validates a positive (nonzero) safe integer, e.g. a requested limit or byte count. */
export function asPositiveSafeInteger(value: number): PositiveSafeInteger {
  if (!Number.isSafeInteger(value) || value <= 0) {
    throw new TypeError("value must be a positive safe integer");
  }
  return value as PositiveSafeInteger;
}

/** Validates a safe integer for use as an {@link InlineJsonValue} number. */
export function asSafeJsonInteger(value: number): SafeJsonInteger {
  if (!Number.isSafeInteger(value)) {
    throw new TypeError("value must be a safe JSON integer");
  }
  return value as SafeJsonInteger;
}
