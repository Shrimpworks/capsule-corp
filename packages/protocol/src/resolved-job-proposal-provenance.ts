import type { ResolvedJobProposalPlanInputs } from "./job-proposal-resolver.js";

const retainedPlanInputs = new WeakSet<object>();

export function retainResolvedJobProposalPlanInputs(
  value: ResolvedJobProposalPlanInputs,
): ResolvedJobProposalPlanInputs {
  retainedPlanInputs.add(value);
  return value;
}

export function isRetainedResolvedJobProposalPlanInputs(
  value: unknown,
): value is ResolvedJobProposalPlanInputs {
  return typeof value === "object" && value !== null && retainedPlanInputs.has(value);
}
