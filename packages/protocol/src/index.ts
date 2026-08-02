/**
 * Pre-freeze protocol scaffold used by current examples and tests.
 *
 * The target v0 protocol will replace this mixed Job authority model after the
 * blocking feasibility gates. See docs/protocol/OBJECT_MODEL.md.
 */

export * from "./internal-contract-candidates.js";
export * from "./job-proposal.js";
export * from "./job-proposal-decoder.js";
export * from "./job-proposal-resolver.js";

export const API_VERSION = "capsule.dev/v0" as const;

export type RuntimeName = "bun" | "node" | "deno";
export type UserExposure = "none" | "metadata" | "full";
export type AgentExposure = "none" | "metadata" | "preview" | "full";

export interface RuntimeSelector {
  name: RuntimeName;
  profile: string;
  digest?: `sha256:${string}`;
}

export interface SourceBundle {
  entrypoint: string;
  files: Record<string, string>;
}

export interface InputCapability {
  id: string;
  kind: "file";
  mountPath: string;
  access: "read";
}

export interface RequestedCapabilities {
  network: Array<{ host: string; ports: number[] }>;
  subprocesses: string[];
  environment: string[];
  nativeAddons: boolean;
  ffi: boolean;
  packageInstallation: boolean;
}

export interface ResourceLimits {
  wallTimeMs: number;
  cpuTimeMs: number;
  memoryBytes: number;
  temporaryStorageBytes: number;
  pids: number;
  logBytes: number;
  artifactBytes: number;
}

export interface ExposurePolicy {
  user: UserExposure;
  agent: AgentExposure;
}

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

export interface RuntimeProfileDescriptor {
  name: string;
  runtime: RuntimeName;
  status: "draft" | "active" | "retired";
  available: boolean;
}

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
