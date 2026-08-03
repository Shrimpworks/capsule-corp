/**
 * Provenance tracking for {@link ResolvedJobProposalPlanInputs}, the same
 * WeakSet-identity pattern as decoded-job-proposal-candidate.ts.
 * constructExecutionPlan (execution-plan-builder.ts) refuses any plan-inputs
 * value not retained here, so a caller cannot fabricate resolved plan inputs
 * by hand-constructing a matching object and skipping resolveJobProposal's
 * semantic/policy/binding checks.
 */
import type { ResolvedJobProposalPlanInputs } from "./job-proposal-resolver.js";

const retainedPlanInputs = new WeakSet<object>();

/** Marks a resolver-issued plan-inputs value as retained. Call exactly once per resolution. */
export function retainResolvedJobProposalPlanInputs(
  value: ResolvedJobProposalPlanInputs,
): ResolvedJobProposalPlanInputs {
  retainedPlanInputs.add(value);
  return value;
}

/** True only for the exact object identity returned by a prior retain call. */
export function isRetainedResolvedJobProposalPlanInputs(
  value: unknown,
): value is ResolvedJobProposalPlanInputs {
  return typeof value === "object" && value !== null && retainedPlanInputs.has(value);
}
