import type { DecodedJobProposalCandidate } from "./job-proposal.js";

const decodedCandidates = new WeakSet<object>();

export function retainDecodedJobProposalCandidate(
  proposal: DecodedJobProposalCandidate,
): DecodedJobProposalCandidate {
  decodedCandidates.add(proposal);
  return proposal;
}

export function isRetainedDecodedJobProposalCandidate(
  proposal: DecodedJobProposalCandidate,
): boolean {
  return typeof proposal === "object" && proposal !== null && decodedCandidates.has(proposal);
}
