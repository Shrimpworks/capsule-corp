import { retainDecodedJobProposalCandidate } from "./decoded-job-proposal-candidate.js";
import {
  asPositiveSafeInteger,
  asRuntimeProfileAlias,
  asSafeJsonInteger,
  asSourceEntrypoint,
  asSourcePath,
  type DecodedJobProposalCandidate,
  type InlineJsonValue,
  JOB_PROPOSAL_API_VERSION,
  JOB_PROPOSAL_KIND,
  type JobProposal,
  PRIMARY_DATA_INPUT_SLOT,
  type RuntimeProfileAlias,
  type SourceEntrypoint,
  type SourcePath,
  TRANSFORMED_JSON_OUTPUT_SLOT,
} from "./job-proposal.js";
import {
  type JobProposalRawRefusal,
  predecodeStrictJobProposalJson,
  type StrictJsonValue,
} from "./strict-job-proposal-json.js";

export type {
  JobProposalRawRefusal,
  JobProposalRawRefusalCode,
} from "./strict-job-proposal-json.js";

export type JobProposalSchemaClassification = "UNSUPPORTED" | "SCHEMA" | "DOMAIN";

export type JobProposalSchemaRefusalCode =
  | "ROOT_TYPE"
  | "MISSING_FIELD"
  | "UNKNOWN_FIELD"
  | "OBJECT_VERSION"
  | "OBJECT_TYPE"
  | "FIELD_TYPE"
  | "FIELD_VALUE"
  | "SOURCE_PATH"
  | "SOURCE_ENTRYPOINT"
  | "RUNTIME_PROFILE_ALIAS"
  | "POSITIVE_INTEGER"
  | "INLINE_JSON_KEY"
  | "INPUT_SLOT"
  | "INPUT_SLOT_DOMAIN"
  | "OUTPUT_SLOT"
  | "OUTPUT_SLOT_DOMAIN"
  | "LABELS";

export interface JobProposalSchemaRefusal {
  readonly owner: "job-proposal-schema";
  readonly classification: JobProposalSchemaClassification;
  readonly code: JobProposalSchemaRefusalCode;
}

export type JobProposalDecodeRefusal = JobProposalRawRefusal | JobProposalSchemaRefusal;

export type JobProposalDecodeResult =
  | { readonly ok: true; readonly proposal: DecodedJobProposalCandidate }
  | { readonly ok: false; readonly refusal: JobProposalDecodeRefusal };

class SchemaFailure extends Error {
  constructor(
    readonly classification: JobProposalSchemaClassification,
    readonly code: JobProposalSchemaRefusalCode,
  ) {
    super(code);
  }
}

const TOP_LEVEL_REQUIRED = [
  "apiVersion",
  "kind",
  "source",
  "runtimeProfile",
  "input",
  "requestedLimits",
  "outputs",
] as const;
const TOP_LEVEL_OPTIONAL = ["labels"] as const;
const SOURCE_REQUIRED = ["entrypoint", "files"] as const;
const INPUT_REQUIRED = ["slot", "kind", "value"] as const;
const REQUESTED_LIMITS_OPTIONAL = ["wallTimeMs"] as const;
const OUTPUT_REQUIRED = ["slot", "kind", "maxBytes"] as const;
const INLINE_JSON_KEY_PATTERN = /^[A-Za-z0-9][A-Za-z0-9_.-]{0,127}$/u;
const LABEL_KEY_PATTERN = /^[a-z][a-z0-9-]{0,31}$/u;
const LABEL_VALUE_PATTERN = /^[ -~]{0,128}$/u;

/**
 * Strictly decodes caller-owned bytes into a passive raw/schema-validated
 * candidate. The returned value contains no resolved policy, profile, content,
 * backend, approval, registration, or execution authority.
 */
export function decodeJobProposal(bytes: Uint8Array): JobProposalDecodeResult {
  const raw = predecodeStrictJobProposalJson(bytes);
  if (!raw.ok) {
    return raw;
  }

  try {
    const proposal = retainDecodedJobProposalCandidate(decodeClosedProposal(raw.value));
    return Object.freeze({ ok: true, proposal });
  } catch (error) {
    if (error instanceof SchemaFailure) {
      return schemaRejected(error.classification, error.code);
    }
    throw error;
  }
}

function decodeClosedProposal(value: StrictJsonValue): DecodedJobProposalCandidate {
  const root = requireObject(value, "ROOT_TYPE");
  requireClosedFields(root, TOP_LEVEL_REQUIRED, TOP_LEVEL_OPTIONAL);

  const apiVersion = requireString(root.apiVersion);
  if (apiVersion !== JOB_PROPOSAL_API_VERSION) {
    schemaFail("UNSUPPORTED", "OBJECT_VERSION");
  }
  const kind = requireString(root.kind);
  if (kind !== JOB_PROPOSAL_KIND) {
    schemaFail("DOMAIN", "OBJECT_TYPE");
  }

  const source = decodeSource(root.source);
  const runtimeProfile = decodeRuntimeProfile(root.runtimeProfile);
  const input = decodeInput(root.input);
  const requestedLimits = decodeRequestedLimits(root.requestedLimits);
  const outputs = decodeOutputs(root.outputs);
  const labels = root.labels === undefined ? undefined : decodeLabels(root.labels);

  const proposal: JobProposal = Object.freeze({
    apiVersion: JOB_PROPOSAL_API_VERSION,
    kind: JOB_PROPOSAL_KIND,
    source,
    runtimeProfile,
    input,
    requestedLimits,
    outputs,
    ...(labels === undefined ? {} : { labels }),
  });
  return proposal as DecodedJobProposalCandidate;
}

function decodeSource(value: StrictJsonValue | undefined): JobProposal["source"] {
  const source = requireObject(value);
  requireClosedFields(source, SOURCE_REQUIRED);

  const entrypointValue = requireString(source.entrypoint);
  let entrypoint: SourceEntrypoint;
  try {
    entrypoint = asSourceEntrypoint(entrypointValue);
  } catch {
    schemaFail("SCHEMA", "SOURCE_ENTRYPOINT");
  }

  const filesValue = requireObject(source.files);
  const fileEntries = Object.entries(filesValue);
  if (fileEntries.length < 1 || fileEntries.length > 32) {
    schemaFail("SCHEMA", "FIELD_VALUE");
  }
  const files: Array<readonly [SourcePath, string]> = [];
  for (const [pathValue, contentValue] of fileEntries) {
    let path: SourcePath;
    try {
      path = asSourcePath(pathValue);
    } catch {
      schemaFail("SCHEMA", "SOURCE_PATH");
    }
    files.push([path, requireString(contentValue)]);
  }

  return Object.freeze({
    entrypoint,
    files: Object.freeze(Object.fromEntries(files)) as Readonly<Record<SourcePath, string>>,
  });
}

function decodeRuntimeProfile(value: StrictJsonValue | undefined): RuntimeProfileAlias {
  const alias = requireString(value);
  try {
    return asRuntimeProfileAlias(alias);
  } catch {
    schemaFail("SCHEMA", "RUNTIME_PROFILE_ALIAS");
  }
}

function decodeInput(value: StrictJsonValue | undefined): JobProposal["input"] {
  const input = requireObject(value);
  requireClosedFields(input, INPUT_REQUIRED);
  const slot = requireString(input.slot);
  if (slot !== PRIMARY_DATA_INPUT_SLOT) {
    schemaFail(
      slot === TRANSFORMED_JSON_OUTPUT_SLOT ? "DOMAIN" : "SCHEMA",
      slot === TRANSFORMED_JSON_OUTPUT_SLOT ? "INPUT_SLOT_DOMAIN" : "INPUT_SLOT",
    );
  }
  const kind = requireString(input.kind);
  if (kind !== "inline-json") {
    schemaFail("UNSUPPORTED", "FIELD_VALUE");
  }
  return Object.freeze({
    slot: PRIMARY_DATA_INPUT_SLOT,
    kind: "inline-json",
    value: decodeInlineJson(input.value),
  });
}

function decodeRequestedLimits(value: StrictJsonValue | undefined): JobProposal["requestedLimits"] {
  const requestedLimits = requireObject(value);
  requireClosedFields(requestedLimits, [], REQUESTED_LIMITS_OPTIONAL);
  if (requestedLimits.wallTimeMs === undefined) {
    return Object.freeze({});
  }
  return Object.freeze({ wallTimeMs: decodePositiveInteger(requestedLimits.wallTimeMs) });
}

function decodeOutputs(value: StrictJsonValue | undefined): JobProposal["outputs"] {
  if (!Array.isArray(value) || value.length !== 1) {
    schemaFail("SCHEMA", "FIELD_VALUE");
  }
  const output = requireObject(value[0]);
  requireClosedFields(output, OUTPUT_REQUIRED);
  const slot = requireString(output.slot);
  if (slot !== TRANSFORMED_JSON_OUTPUT_SLOT) {
    schemaFail(
      slot === PRIMARY_DATA_INPUT_SLOT ? "DOMAIN" : "SCHEMA",
      slot === PRIMARY_DATA_INPUT_SLOT ? "OUTPUT_SLOT_DOMAIN" : "OUTPUT_SLOT",
    );
  }
  const kind = requireString(output.kind);
  if (kind !== "inline-json") {
    schemaFail("UNSUPPORTED", "FIELD_VALUE");
  }
  return Object.freeze([
    Object.freeze({
      slot: TRANSFORMED_JSON_OUTPUT_SLOT,
      kind: "inline-json",
      maxBytes: decodePositiveInteger(output.maxBytes),
    }),
  ]) as JobProposal["outputs"];
}

function decodeLabels(value: StrictJsonValue): Readonly<Record<string, string>> {
  const labelsValue = requireObject(value);
  const entries = Object.entries(labelsValue);
  if (entries.length > 8) {
    schemaFail("SCHEMA", "LABELS");
  }
  const labels: Array<readonly [string, string]> = [];
  for (const [key, rawValue] of entries) {
    if (!LABEL_KEY_PATTERN.test(key)) {
      schemaFail("SCHEMA", "LABELS");
    }
    const labelValue = requireString(rawValue);
    if (!LABEL_VALUE_PATTERN.test(labelValue)) {
      schemaFail("SCHEMA", "LABELS");
    }
    labels.push([key, labelValue]);
  }
  return Object.freeze(Object.fromEntries(labels));
}

function decodeInlineJson(value: StrictJsonValue | undefined): InlineJsonValue {
  if (value === null || typeof value === "boolean" || typeof value === "string") {
    return value;
  }
  if (typeof value === "number") {
    return asSafeJsonInteger(value);
  }
  if (Array.isArray(value)) {
    return Object.freeze(value.map((entry) => decodeInlineJson(entry)));
  }
  const object = requireObject(value);
  const entries: Array<readonly [string, InlineJsonValue]> = [];
  for (const [key, entry] of Object.entries(object)) {
    if (!INLINE_JSON_KEY_PATTERN.test(key)) {
      schemaFail("SCHEMA", "INLINE_JSON_KEY");
    }
    entries.push([key, decodeInlineJson(entry)]);
  }
  return Object.freeze(Object.fromEntries(entries)) as Readonly<Record<string, InlineJsonValue>>;
}

function decodePositiveInteger(
  value: StrictJsonValue | undefined,
): JobProposal["outputs"][0]["maxBytes"] {
  if (typeof value !== "number") {
    schemaFail("SCHEMA", "FIELD_TYPE");
  }
  try {
    return asPositiveSafeInteger(value);
  } catch {
    schemaFail("SCHEMA", "POSITIVE_INTEGER");
  }
}

function requireObject(
  value: StrictJsonValue | undefined,
  code: JobProposalSchemaRefusalCode = "FIELD_TYPE",
): { readonly [key: string]: StrictJsonValue } {
  if (value === null || typeof value !== "object" || Array.isArray(value)) {
    schemaFail("SCHEMA", code);
  }
  return value as { readonly [key: string]: StrictJsonValue };
}

function requireString(value: StrictJsonValue | undefined): string {
  if (typeof value !== "string") {
    schemaFail("SCHEMA", "FIELD_TYPE");
  }
  return value;
}

function requireClosedFields(
  object: { readonly [key: string]: StrictJsonValue },
  required: readonly string[],
  optional: readonly string[] = [],
): void {
  for (const field of required) {
    if (!Object.hasOwn(object, field)) {
      schemaFail("SCHEMA", "MISSING_FIELD");
    }
  }
  const admitted = new Set([...required, ...optional]);
  for (const field of Object.keys(object)) {
    if (!admitted.has(field)) {
      schemaFail("UNSUPPORTED", "UNKNOWN_FIELD");
    }
  }
}

function schemaRejected(
  classification: JobProposalSchemaClassification,
  code: JobProposalSchemaRefusalCode,
): JobProposalDecodeResult {
  return Object.freeze({
    ok: false,
    refusal: Object.freeze({ owner: "job-proposal-schema", classification, code }),
  });
}

function schemaFail(
  classification: JobProposalSchemaClassification,
  code: JobProposalSchemaRefusalCode,
): never {
  throw new SchemaFailure(classification, code);
}
