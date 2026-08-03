/**
 * Provenance tracking for {@link DecodedJobProposalCandidate}. A WeakSet of
 * object identities is the only way to prove a candidate actually came from
 * the decoder — resolveJobProposal (job-proposal-resolver.ts) refuses any
 * structurally identical object that wasn't retained here, so a caller
 * cannot forge a "decoded" candidate by hand-constructing an object that
 * merely matches the shape.
 */
import type { DecodedJobProposalCandidate } from "./job-proposal.js";

const decodedCandidates = new WeakSet<object>();

/** Marks a decoder-issued candidate as retained. Call exactly once per decode. */
export function retainDecodedJobProposalCandidate(
  proposal: DecodedJobProposalCandidate,
): DecodedJobProposalCandidate {
  decodedCandidates.add(proposal);
  return proposal;
}

/** True only for the exact object identity returned by a prior retain call. */
export function isRetainedDecodedJobProposalCandidate(
  proposal: DecodedJobProposalCandidate,
): boolean {
  return typeof proposal === "object" && proposal !== null && decodedCandidates.has(proposal);
}
