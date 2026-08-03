package registrationstate

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	MaxUsableApprovals     = 256
	MaxRetainedApprovals   = 4_096
	MaxNonterminalAttempts = 256
	MaxRetainedAttempts    = 4_096
	ApprovalStoreVersion   = uint64(0)
	AttemptStoreVersion    = uint64(0)
)

func validateApprovalAttemptState(state installationState) error {
	if len(state.Approvals) > MaxRetainedApprovals || len(state.Attempts) > MaxRetainedAttempts {
		return errors.New("fixed approval/attempt state exceeds retained capacity")
	}
	if countUsableApprovals(state.Approvals, state.TimeHighWaterUnixSeconds) > MaxUsableApprovals ||
		countNonterminalAttempts(state.Attempts) > MaxNonterminalAttempts {
		return errors.New("fixed approval/attempt state exceeds active capacity")
	}
	approvalIDs := make(map[approvalattempt.ApprovalID]int, len(state.Approvals))
	payloads := make(map[approvalattempt.ApprovalPayloadDigest]int, len(state.Approvals))
	nonces := make(map[approvalattempt.AttemptNonce]int, len(state.Approvals))
	for index, record := range state.Approvals {
		if record.ApprovalID == (approvalattempt.ApprovalID{}) ||
			record.AttemptNonce == (approvalattempt.AttemptNonce{}) ||
			record.PayloadDigest == (approvalattempt.ApprovalPayloadDigest{}) ||
			record.AuthorizationIdentity == (approvalattempt.ApprovalKeyAuthorizationIdentity{}) ||
			record.StorageFormatVersion != ApprovalStoreVersion {
			return errors.New("fixed approval state has an invalid record")
		}
		if _, exists := approvalIDs[record.ApprovalID]; exists {
			return errors.New("fixed approval state has a duplicate approval ID")
		}
		if _, exists := payloads[record.PayloadDigest]; exists {
			return errors.New("fixed approval state has a duplicate payload digest")
		}
		if _, exists := nonces[record.AttemptNonce]; exists {
			return errors.New("fixed approval state has a duplicate nonce")
		}
		approvalIDs[record.ApprovalID] = index
		payloads[record.PayloadDigest] = index
		nonces[record.AttemptNonce] = index
		if len(record.ExactEnvelopeBytes) == 0 ||
			len(record.ExactEnvelopeBytes) > approvalattempt.ApprovalEnvelopeRawMaxBytes ||
			len(record.ExactPayloadBytes) == 0 ||
			len(record.ExactPayloadBytes) > approvalattempt.ApprovalPayloadRawMaxBytes ||
			len(record.ExactProtectedBytes) > approvalattempt.ApprovalProtectedRawMaxBytes ||
			len(record.ProtectedKeyID) == 0 || len(record.ProtectedKeyID) > 64 {
			return errors.New("fixed approval state has invalid exact bytes")
		}
		payloadDigest := approvalattempt.ApprovalPayloadDigest(sha256.Sum256(record.ExactPayloadBytes))
		envelopeDigest := approvalattempt.ApprovalEnvelopeDigest(sha256.Sum256(record.ExactEnvelopeBytes))
		if payloadDigest != record.PayloadDigest || envelopeDigest != record.EnvelopeDigest {
			return errors.New("fixed approval state has an exact-byte digest mismatch")
		}
		registrationIndex := findRegistration(state.Registrations, record.RegistrationID)
		if registrationIndex < 0 || !approvalMatchesRegistration(record, state.Registrations[registrationIndex]) {
			return errors.New("fixed approval state has a registration cross-link mismatch")
		}
		if record.Purpose != approvalattempt.ApprovalGrantPurpose ||
			record.Audience != approvalattempt.ApprovalGrantAudience ||
			record.IssuedAt >= record.ExpiresAt ||
			record.ExpiresAt > state.Registrations[registrationIndex].Index.ExpiresAt {
			return errors.New("fixed approval state has invalid signed bindings")
		}
		switch record.State {
		case approvalattempt.ApprovalUsable, approvalattempt.ApprovalInvalidated:
			if record.ConsumedAttemptID != (approvalattempt.AttemptID{}) {
				return errors.New("fixed approval state has a terminal cross-link")
			}
		case approvalattempt.ApprovalConsumed:
			if record.ConsumedAttemptID == (approvalattempt.AttemptID{}) {
				return errors.New("fixed approval state consumed without an attempt")
			}
		default:
			return errors.New("fixed approval state has an unknown state")
		}
	}

	attemptIDs := make(map[approvalattempt.AttemptID]int, len(state.Attempts))
	ownedApprovals := make(map[approvalattempt.ApprovalID]struct{}, len(state.Attempts))
	for index, attempt := range state.Attempts {
		if attempt.AttemptID == (approvalattempt.AttemptID{}) ||
			attempt.ApprovalID == (approvalattempt.ApprovalID{}) ||
			attempt.AttemptNonce == (approvalattempt.AttemptNonce{}) ||
			attempt.State != approvalattempt.AttemptCreated ||
			attempt.StorageFormatVersion != AttemptStoreVersion ||
			attempt.CreatedAt > state.TimeHighWaterUnixSeconds {
			return errors.New("fixed attempt state has an invalid record")
		}
		if _, exists := attemptIDs[attempt.AttemptID]; exists {
			return errors.New("fixed attempt state has a duplicate attempt ID")
		}
		if _, exists := ownedApprovals[attempt.ApprovalID]; exists {
			return errors.New("fixed attempt state has duplicate approval ownership")
		}
		attemptIDs[attempt.AttemptID] = index
		ownedApprovals[attempt.ApprovalID] = struct{}{}
		approvalIndex, exists := approvalIDs[attempt.ApprovalID]
		if !exists || !attemptMatchesApproval(attempt, state.Approvals[approvalIndex]) {
			return errors.New("fixed attempt state has an approval cross-link mismatch")
		}
	}
	for _, approval := range state.Approvals {
		if approval.State == approvalattempt.ApprovalConsumed {
			attemptIndex, exists := attemptIDs[approval.ConsumedAttemptID]
			if !exists || state.Attempts[attemptIndex].ApprovalID != approval.ApprovalID {
				return errors.New("fixed approval state consumed without its attempt")
			}
		}
	}
	approvalDigest, err := approvalSetDigest(state.Approvals)
	if err != nil {
		return fmt.Errorf("compute approval set digest: %w", err)
	}
	attemptDigest, err := attemptSetDigest(state.Attempts)
	if err != nil {
		return fmt.Errorf("compute attempt set digest: %w", err)
	}
	approvalDigestValid := approvalDigest == state.ApprovalSetDigest ||
		(len(state.Approvals) == 0 && state.ApprovalSetDigest == ([32]byte{}))
	attemptDigestValid := attemptDigest == state.AttemptSetDigest ||
		(len(state.Attempts) == 0 && state.AttemptSetDigest == ([32]byte{}))
	if !approvalDigestValid || !attemptDigestValid {
		return errors.New("fixed approval/attempt set digest mismatch")
	}
	return nil
}

func approvalMatchesRegistration(record approvalattempt.ApprovalRecord, registration registrationEntry) bool {
	return record.RegistrationSequence == registration.Index.RegistrationSequence &&
		record.PlanDigest == registration.Record.RecomputedPlanDigest &&
		record.InstallationID == registration.Index.InstallationID &&
		record.EpochSequence == registration.Index.EpochSequence &&
		record.EpochDigest == registration.Index.EpochDigest &&
		record.SupervisorID == registration.Index.SupervisorID
}

func attemptMatchesApproval(attempt approvalattempt.ExecutionAttempt, approval approvalattempt.ApprovalRecord) bool {
	return approval.State == approvalattempt.ApprovalConsumed &&
		approval.ConsumedAttemptID == attempt.AttemptID &&
		attempt.ApprovalID == approval.ApprovalID &&
		attempt.AttemptNonce == approval.AttemptNonce &&
		attempt.RegistrationID == approval.RegistrationID &&
		attempt.RegistrationSequence == approval.RegistrationSequence &&
		attempt.PlanDigest == approval.PlanDigest &&
		attempt.InstallationID == approval.InstallationID &&
		attempt.EpochSequence == approval.EpochSequence &&
		attempt.EpochDigest == approval.EpochDigest &&
		attempt.SupervisorID == approval.SupervisorID &&
		attempt.Purpose == approval.Purpose && attempt.Audience == approval.Audience &&
		attempt.ApprovalPayloadDigest == approval.PayloadDigest &&
		attempt.AuthorizationIdentity == approval.AuthorizationIdentity
}

func approvalSetDigest(records []approvalattempt.ApprovalRecord) ([32]byte, error) {
	if len(records) == 0 {
		return emptyApprovalSetDigest(), nil
	}
	encoded, err := json.Marshal(records)
	if err != nil {
		return [32]byte{}, fmt.Errorf("encode approval set for digest: %w", err)
	}
	return sha256.Sum256(append([]byte("approvals:"), encoded...)), nil
}

func attemptSetDigest(records []approvalattempt.ExecutionAttempt) ([32]byte, error) {
	if len(records) == 0 {
		return emptyAttemptSetDigest(), nil
	}
	encoded, err := json.Marshal(records)
	if err != nil {
		return [32]byte{}, fmt.Errorf("encode attempt set for digest: %w", err)
	}
	return sha256.Sum256(append([]byte("attempts:"), encoded...)), nil
}

func countUsableApprovals(records []approvalattempt.ApprovalRecord, now v0candidate.UInt53) int {
	count := 0
	for _, record := range records {
		if record.State == approvalattempt.ApprovalUsable && uint64(record.ExpiresAt) > uint64(now) {
			count++
		}
	}
	return count
}

func countNonterminalAttempts(records []approvalattempt.ExecutionAttempt) int {
	count := 0
	for _, record := range records {
		if record.State == approvalattempt.AttemptCreated {
			count++
		}
	}
	return count
}
