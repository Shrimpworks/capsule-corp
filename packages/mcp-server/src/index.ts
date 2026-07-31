export interface CapsuleToolDescriptor {
  name: string;
  description: string;
  contentAccess: "none" | "metadata" | "approval-required";
}

export const CAPSULE_TOOL_CATALOG: readonly CapsuleToolDescriptor[] = [
  {
    name: "capsule.list_profiles",
    description: "List runtime profiles available to this client.",
    contentAccess: "none",
  },
  {
    name: "capsule.prepare_job",
    description: "Validate a proposed job and return its effective capability request.",
    contentAccess: "none",
  },
  {
    name: "capsule.run_job",
    description: "Run an approved job in a disposable runtime.",
    contentAccess: "none",
  },
  {
    name: "capsule.get_job",
    description: "Read status, metrics, violations, and artifact metadata for a job.",
    contentAccess: "metadata",
  },
  {
    name: "capsule.cancel_job",
    description: "Cancel a running job and destroy its guest environment.",
    contentAccess: "none",
  },
  {
    name: "capsule.list_artifacts",
    description: "List policy-approved artifact descriptors without returning content.",
    contentAccess: "metadata",
  },
  {
    name: "capsule.request_artifact_access",
    description: "Request separately authorized access to guest-controlled artifact content.",
    contentAccess: "approval-required",
  },
  {
    name: "capsule.get_receipt",
    description: "Read the execution receipt without guest-controlled artifact content.",
    contentAccess: "metadata",
  },
] as const;
