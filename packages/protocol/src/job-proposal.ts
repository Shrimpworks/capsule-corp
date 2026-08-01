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

export type SourcePath = string & { readonly [sourcePathBrand]: "SourcePath" };
export type SourceEntrypoint = string & {
  readonly [sourceEntrypointBrand]: "SourceEntrypoint";
};
export type RuntimeProfileAlias = string & {
  readonly [runtimeProfileAliasBrand]: "RuntimeProfileAlias";
};
export type PositiveSafeInteger = number & {
  readonly [positiveSafeIntegerBrand]: "PositiveSafeInteger";
};

export type InlineJsonValue =
  | null
  | boolean
  | number
  | string
  | readonly InlineJsonValue[]
  | { readonly [key: string]: InlineJsonValue };

export interface SourceBundleProposal {
  entrypoint: SourceEntrypoint;
  files: Readonly<Record<SourcePath, string>>;
}

export interface InlineJsonInputProposal {
  slot: typeof PRIMARY_DATA_INPUT_SLOT;
  kind: "inline-json";
  value: InlineJsonValue;
}

export interface RequestedLimits {
  wallTimeMs?: PositiveSafeInteger;
}

export interface InlineJsonOutputProposal {
  slot: typeof TRANSFORMED_JSON_OUTPUT_SLOT;
  kind: "inline-json";
  maxBytes: PositiveSafeInteger;
}

export interface JobProposal {
  apiVersion: typeof JOB_PROPOSAL_API_VERSION;
  kind: typeof JOB_PROPOSAL_KIND;
  source: SourceBundleProposal;
  runtimeProfile: RuntimeProfileAlias;
  input: InlineJsonInputProposal;
  requestedLimits: RequestedLimits;
  outputs: readonly [InlineJsonOutputProposal];
  labels?: Readonly<Record<string, string>>;
}

const SOURCE_PATH_PATTERN =
  /^[A-Za-z0-9][A-Za-z0-9._-]{0,63}(?:\/[A-Za-z0-9][A-Za-z0-9._-]{0,63})*$/u;
const SOURCE_ENTRYPOINT_PATTERN =
  /^(?:[A-Za-z0-9][A-Za-z0-9._-]{0,63}\/)*[A-Za-z0-9][A-Za-z0-9._-]{0,55}\.(?:js|mjs|cjs|ts|mts|cts)$/u;
const RUNTIME_PROFILE_ALIAS_PATTERN = /^[a-z][a-z0-9-]*@[1-9][0-9]*$/u;

export function asSourcePath(value: string): SourcePath {
  if (value.length > 256 || !SOURCE_PATH_PATTERN.test(value)) {
    throw new TypeError("source path is outside the candidate grammar");
  }
  return value as SourcePath;
}

export function asSourceEntrypoint(value: string): SourceEntrypoint {
  if (value.length > 256 || !SOURCE_ENTRYPOINT_PATTERN.test(value)) {
    throw new TypeError("source entrypoint is outside the candidate grammar");
  }
  return value as SourceEntrypoint;
}

export function asRuntimeProfileAlias(value: string): RuntimeProfileAlias {
  if (value.length > 128 || !RUNTIME_PROFILE_ALIAS_PATTERN.test(value)) {
    throw new TypeError("runtime profile alias is outside the candidate grammar");
  }
  return value as RuntimeProfileAlias;
}

export function asPositiveSafeInteger(value: number): PositiveSafeInteger {
  if (!Number.isSafeInteger(value) || value <= 0) {
    throw new TypeError("value must be a positive safe integer");
  }
  return value as PositiveSafeInteger;
}
