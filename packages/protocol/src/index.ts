/**
 * Pre-freeze protocol scaffold used by current examples and tests.
 *
 * The target v0 protocol will replace this mixed Job authority model after the
 * blocking feasibility gates. See docs/protocol/OBJECT_MODEL.md.
 */

export * from "./execution-plan-builder.js";
export * from "./governed-deno-core-c2b-passive-binding.js";
export * from "./governed-deno-core-c2b-passive-binding-v2.js";
export * from "./governed-deno-core-c2b-passive-binding-v3.js";
export * from "./job-proposal.js";
export * from "./job-proposal-decoder.js";
export * from "./job-proposal-resolver.js";
export * from "./mjs-source-candidate.js";
export * from "./source-manifest-candidate.js";
export * from "./typescript-approved-byte-candidate.js";

/**
 * The apiVersion string for the deprecated pre-freeze {@link Job} model
 * below, not for {@link JobProposal}. It is declared separately from
 * {@link JOB_PROPOSAL_API_VERSION} in job-proposal.ts — even though both
 * currently equal "capsule.dev/v0" — because the two constants version two
 * unrelated, independently evolving object models (this file's deprecated
 * `Job`/`RuntimeSelector` scaffold vs. the ADR-0026/ADR-0030 `JobProposal`
 * replacement). Deriving one from the other would wrongly imply the models
 * are coupled and must move together; they are not, and the coordinated
 * protocol cutover may retire this constant (with the `Job` model it
 * belongs to) on its own schedule.
 */
export const API_VERSION = "capsule.dev/v0" as const;

/** Runtime a deprecated pre-freeze {@link Job} may select. */
export type RuntimeName = "bun" | "node" | "deno";
/** How much of a deprecated pre-freeze {@link Job}'s result the end user sees. */
export type UserExposure = "none" | "metadata" | "full";
/** How much of a deprecated pre-freeze {@link Job}'s result the invoking agent sees. */
export type AgentExposure = "none" | "metadata" | "preview" | "full";

/** Runtime selection for a deprecated pre-freeze {@link Job}. */
export interface RuntimeSelector {
  name: RuntimeName;
  profile: string;
  digest?: `sha256:${string}`;
}

/** Source files for a deprecated pre-freeze {@link Job}. */
export interface SourceBundle {
  entrypoint: string;
  files: Record<string, string>;
}

/** A single mounted input capability for a deprecated pre-freeze {@link Job}. */
export interface InputCapability {
  id: string;
  kind: "file";
  mountPath: string;
  access: "read";
}

/** Capabilities a deprecated pre-freeze {@link Job} may request; every field starts deny-by-default. */
export interface RequestedCapabilities {
  network: Array<{ host: string; ports: number[] }>;
  subprocesses: string[];
  environment: string[];
  nativeAddons: boolean;
  ffi: boolean;
  packageInstallation: boolean;
}

/** Resource ceilings for a deprecated pre-freeze {@link Job}. */
export interface ResourceLimits {
  wallTimeMs: number;
  cpuTimeMs: number;
  memoryBytes: number;
  temporaryStorageBytes: number;
  pids: number;
  logBytes: number;
  artifactBytes: number;
}

/** User- and agent-facing exposure levels for a deprecated pre-freeze {@link Job} result or artifact. */
export interface ExposurePolicy {
  user: UserExposure;
  agent: AgentExposure;
}

/** A single declared output artifact contract for a deprecated pre-freeze {@link Job}. */
export interface ArtifactContract {
  guestPath: string;
  name: string;
  mediaType: "application/json" | "application/x-ndjson" | "text/csv" | "text/plain";
  maxBytes: number;
  validation: {
    regularFileOnly: true;
    symlinks: "deny";
    mustParse: boolean;
    jsonSchema?: Record<string, unknown>;
  };
  exposure: ExposurePolicy;
}

/** @deprecated Pre-freeze scaffold; do not extend with additional authority. */
export interface Job {
  apiVersion: typeof API_VERSION;
  kind: "Job";
  runtime: RuntimeSelector;
  source: SourceBundle;
  inputs: {
    inline: Record<string, unknown>;
    capabilities: InputCapability[];
  };
  capabilities: RequestedCapabilities;
  limits: ResourceLimits;
  results: {
    structured: {
      enabled: boolean;
      maxBytes: number;
      exposure: ExposurePolicy;
    };
    artifacts: ArtifactContract[];
  };
  labels?: Record<string, string>;
}

/** Describes one available runtime profile a deprecated pre-freeze {@link Job} can select. */
export interface RuntimeProfileDescriptor {
  name: string;
  runtime: RuntimeName;
  status: "draft" | "active" | "retired";
  available: boolean;
}

/** Returns a {@link RequestedCapabilities} with every capability denied — the deny-by-default starting point for a deprecated pre-freeze {@link Job}. */
export function createDenyByDefaultCapabilities(): RequestedCapabilities {
  return {
    network: [],
    subprocesses: [],
    environment: [],
    nativeAddons: false,
    ffi: false,
    packageInstallation: false,
  };
}
