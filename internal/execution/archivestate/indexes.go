package archivestate

import (
	"bytes"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// ArchiveLocation identifies a record inside a referenced immutable segment.
// It is never authority without the corresponding verified record and digest.
type ArchiveLocation struct {
	SegmentOrdinal ArchiveOrdinal
	CohortOrdinal  ArchiveOrdinal
	RecordOrdinal  ArchiveOrdinal
}

func (location ArchiveLocation) valid() bool {
	return validPositive(uint64(location.SegmentOrdinal)) &&
		validPositive(uint64(location.CohortOrdinal)) &&
		validPositive(uint64(location.RecordOrdinal))
}

type RegistrationIndexEntry struct {
	RegistrationID       v0candidate.RegistrationID
	RegistrationSequence v0candidate.PositiveUInt53
	PlanDigest           v0candidate.ExecutionPlanDigest
	ExpiresAt            v0candidate.UInt53
	Location             ArchiveLocation
	FullRecordDigest     ArchiveRecordDigest
}

type ApprovalIndexEntry struct {
	ApprovalID            approvalattempt.ApprovalID
	RegistrationID        v0candidate.RegistrationID
	PayloadDigest         approvalattempt.ApprovalPayloadDigest
	AuthorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity
	AttemptNonce          approvalattempt.AttemptNonce
	State                 approvalattempt.ApprovalState
	ConsumedAttemptID     approvalattempt.AttemptID
	ExpiresAt             v0candidate.UInt53
	Location              ArchiveLocation
	FullRecordDigest      ArchiveRecordDigest
}

type AttemptIndexEntry struct {
	AttemptID        approvalattempt.AttemptID
	ApprovalID       approvalattempt.ApprovalID
	RegistrationID   v0candidate.RegistrationID
	CreatedAt        v0candidate.UInt53
	LifecycleState   lifecyclestate.LifecycleState
	Location         ArchiveLocation
	FullRecordDigest ArchiveRecordDigest
}

type NonceIndexEntry struct {
	AttemptNonce  approvalattempt.AttemptNonce
	PayloadDigest approvalattempt.ApprovalPayloadDigest
	ApprovalID    approvalattempt.ApprovalID
}

type EffectIndexEntry struct {
	EffectID                   lifecyclestate.EffectID
	AttemptID                  approvalattempt.AttemptID
	OperationSequence          v0candidate.PositiveUInt53
	Operation                  lifecyclestate.Operation
	IssuanceSnapshotGeneration v0candidate.PositiveUInt53
	VisibleV1Seed              bool
}

type InstanceIndexEntry struct {
	InstanceDigest lifecyclestate.BackendInstanceDigest
	AttemptID      approvalattempt.AttemptID
	Location       ArchiveLocation
}

type ApprovalReplayIndexEntry struct {
	PayloadDigest         approvalattempt.ApprovalPayloadDigest
	ExactPayloadDigest    ArchiveRecordDigest
	AuthorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity
	ApprovalID            approvalattempt.ApprovalID
	State                 approvalattempt.ApprovalState
	Location              ArchiveLocation
}

type AttemptReplayIndexEntry struct {
	RegistrationID v0candidate.RegistrationID
	ApprovalID     approvalattempt.ApprovalID
	AttemptID      approvalattempt.AttemptID
	State          approvalattempt.AttemptState
}

// ArchiveIndexesView is the complete closed sorted archive index projection.
// Every slice is required to be non-nil, including in the empty projection.
type ArchiveIndexesView struct {
	Registrations  []RegistrationIndexEntry
	Approvals      []ApprovalIndexEntry
	Attempts       []AttemptIndexEntry
	Nonces         []NonceIndexEntry
	Effects        []EffectIndexEntry
	Instances      []InstanceIndexEntry
	ApprovalReplay []ApprovalReplayIndexEntry
	AttemptReplay  []AttemptReplayIndexEntry
}

// ArchiveIndexDigests binds every complete index independently.
type ArchiveIndexDigests struct {
	Registrations  ArchiveIndexDigest
	Approvals      ArchiveIndexDigest
	Attempts       ArchiveIndexDigest
	Nonces         ArchiveIndexDigest
	Effects        ArchiveIndexDigest
	Instances      ArchiveIndexDigest
	ApprovalReplay ArchiveIndexDigest
	AttemptReplay  ArchiveIndexDigest
}

// ArchiveIndexes is an immutable validated collection projection.
type ArchiveIndexes struct {
	view     ArchiveIndexesView
	digests  ArchiveIndexDigests
	combined ArchiveCombinedIndexDigest
}

func NewArchiveIndexes(view ArchiveIndexesView) (ArchiveIndexes, error) {
	if err := validateArchiveIndexesView(view); err != nil {
		return ArchiveIndexes{}, err
	}
	cloned := cloneArchiveIndexesView(view)
	digests := digestArchiveIndexes(cloned)
	return ArchiveIndexes{
		view: cloned, digests: digests, combined: digestCombinedIndexes(digests),
	}, nil
}

func EmptyArchiveIndexes() ArchiveIndexes {
	view := ArchiveIndexesView{
		Registrations:  make([]RegistrationIndexEntry, 0),
		Approvals:      make([]ApprovalIndexEntry, 0),
		Attempts:       make([]AttemptIndexEntry, 0),
		Nonces:         make([]NonceIndexEntry, 0),
		Effects:        make([]EffectIndexEntry, 0),
		Instances:      make([]InstanceIndexEntry, 0),
		ApprovalReplay: make([]ApprovalReplayIndexEntry, 0),
		AttemptReplay:  make([]AttemptReplayIndexEntry, 0),
	}
	digests := digestArchiveIndexes(view)
	return ArchiveIndexes{view: view, digests: digests, combined: digestCombinedIndexes(digests)}
}

func (indexes ArchiveIndexes) View() ArchiveIndexesView {
	return cloneArchiveIndexesView(indexes.view)
}

func (indexes ArchiveIndexes) Digests() ArchiveIndexDigests { return indexes.digests }

func (indexes ArchiveIndexes) CombinedDigest() ArchiveCombinedIndexDigest {
	return indexes.combined
}

func (indexes ArchiveIndexes) counts() ArchiveCounts {
	view := indexes.view
	return ArchiveCounts{
		Registrations: uint64(len(view.Registrations)), Approvals: uint64(len(view.Approvals)),
		Attempts: uint64(len(view.Attempts)), Lifecycles: uint64(len(view.Attempts)),
		Nonces: uint64(len(view.Nonces)), Effects: uint64(len(view.Effects)),
		Instances: uint64(len(view.Instances)), ApprovalReplay: uint64(len(view.ApprovalReplay)),
		AttemptReplay: uint64(len(view.AttemptReplay)),
	}
}

func cloneArchiveIndexesView(view ArchiveIndexesView) ArchiveIndexesView {
	view.Registrations = cloneSlice(view.Registrations)
	view.Approvals = cloneSlice(view.Approvals)
	view.Attempts = cloneSlice(view.Attempts)
	view.Nonces = cloneSlice(view.Nonces)
	view.Effects = cloneSlice(view.Effects)
	view.Instances = cloneSlice(view.Instances)
	view.ApprovalReplay = cloneSlice(view.ApprovalReplay)
	view.AttemptReplay = cloneSlice(view.AttemptReplay)
	return view
}

func validateArchiveIndexesView(view ArchiveIndexesView) error {
	if view.Registrations == nil || view.Approvals == nil || view.Attempts == nil ||
		view.Nonces == nil || view.Effects == nil || view.Instances == nil ||
		view.ApprovalReplay == nil || view.AttemptReplay == nil {
		return classified(ClassificationSchema, "archive-index-nil")
	}
	lengths := [...]int{
		len(view.Registrations), len(view.Approvals), len(view.Attempts), len(view.Nonces),
		len(view.Effects), len(view.Instances), len(view.ApprovalReplay), len(view.AttemptReplay),
	}
	for _, length := range lengths {
		if length > MaxArchiveIndexEntries {
			return classified(ClassificationCapacity, "archive-index-entries")
		}
	}
	if err := validateRegistrationIndex(view.Registrations); err != nil {
		return err
	}
	if err := validateApprovalIndex(view.Approvals); err != nil {
		return err
	}
	if err := validateAttemptIndex(view.Attempts); err != nil {
		return err
	}
	if err := validateNonceIndex(view.Nonces); err != nil {
		return err
	}
	if err := validateEffectIndex(view.Effects); err != nil {
		return err
	}
	if err := validateInstanceIndex(view.Instances); err != nil {
		return err
	}
	if err := validateApprovalReplayIndex(view.ApprovalReplay); err != nil {
		return err
	}
	return validateAttemptReplayIndex(view.AttemptReplay)
}

func validApprovalState(state approvalattempt.ApprovalState) bool {
	return state == approvalattempt.ApprovalUsable || state == approvalattempt.ApprovalConsumed ||
		state == approvalattempt.ApprovalInvalidated
}

func validateRegistrationIndex(entries []RegistrationIndexEntry) error {
	return validateStrictEntries(entries, func(entry RegistrationIndexEntry) ([]byte, error) {
		if entry.RegistrationID == (v0candidate.RegistrationID{}) ||
			!validPositive(uint64(entry.RegistrationSequence)) ||
			uint64(entry.ExpiresAt) > v0candidate.MaxSafeInteger || !entry.Location.valid() ||
			zeroDigest(entry.FullRecordDigest) {
			return nil, classified(ClassificationBinding, "registration-index-entry")
		}
		return append([]byte(nil), entry.RegistrationID[:]...), nil
	}, "registration-index-order")
}

func validateApprovalIndex(entries []ApprovalIndexEntry) error {
	return validateStrictEntries(entries, func(entry ApprovalIndexEntry) ([]byte, error) {
		if entry.ApprovalID == (approvalattempt.ApprovalID{}) ||
			entry.RegistrationID == (v0candidate.RegistrationID{}) ||
			entry.AttemptNonce == (approvalattempt.AttemptNonce{}) ||
			zeroDigest(entry.PayloadDigest) || zeroDigest(entry.AuthorizationIdentity) ||
			!validApprovalState(entry.State) || uint64(entry.ExpiresAt) > v0candidate.MaxSafeInteger ||
			!entry.Location.valid() || zeroDigest(entry.FullRecordDigest) {
			return nil, classified(ClassificationBinding, "approval-index-entry")
		}
		if (entry.State == approvalattempt.ApprovalConsumed) !=
			(entry.ConsumedAttemptID != (approvalattempt.AttemptID{})) {
			return nil, classified(ClassificationBinding, "approval-index-consumed-link")
		}
		return append([]byte(nil), entry.ApprovalID[:]...), nil
	}, "approval-index-order")
}

func validateAttemptIndex(entries []AttemptIndexEntry) error {
	return validateStrictEntries(entries, func(entry AttemptIndexEntry) ([]byte, error) {
		if entry.AttemptID == (approvalattempt.AttemptID{}) ||
			entry.ApprovalID == (approvalattempt.ApprovalID{}) ||
			entry.RegistrationID == (v0candidate.RegistrationID{}) ||
			uint64(entry.CreatedAt) > v0candidate.MaxSafeInteger ||
			!validLifecycleState(entry.LifecycleState) || !entry.Location.valid() ||
			zeroDigest(entry.FullRecordDigest) {
			return nil, classified(ClassificationBinding, "attempt-index-entry")
		}
		return append([]byte(nil), entry.AttemptID[:]...), nil
	}, "attempt-index-order")
}

func validateNonceIndex(entries []NonceIndexEntry) error {
	return validateStrictEntries(entries, func(entry NonceIndexEntry) ([]byte, error) {
		if entry.AttemptNonce == (approvalattempt.AttemptNonce{}) ||
			zeroDigest(entry.PayloadDigest) || entry.ApprovalID == (approvalattempt.ApprovalID{}) {
			return nil, classified(ClassificationBinding, "nonce-index-entry")
		}
		return append([]byte(nil), entry.AttemptNonce[:]...), nil
	}, "nonce-index-order")
}

func validateEffectIndex(entries []EffectIndexEntry) error {
	return validateStrictEntries(entries, func(entry EffectIndexEntry) ([]byte, error) {
		if entry.EffectID.IsZero() || entry.AttemptID == (approvalattempt.AttemptID{}) ||
			!validPositive(uint64(entry.OperationSequence)) || !validLifecycleOperation(entry.Operation) ||
			!validPositive(uint64(entry.IssuanceSnapshotGeneration)) {
			return nil, classified(ClassificationBinding, "effect-index-entry")
		}
		return append([]byte(nil), entry.EffectID[:]...), nil
	}, "effect-index-order")
}

func validateInstanceIndex(entries []InstanceIndexEntry) error {
	return validateStrictEntries(entries, func(entry InstanceIndexEntry) ([]byte, error) {
		if zeroDigest(entry.InstanceDigest) || entry.AttemptID == (approvalattempt.AttemptID{}) ||
			!entry.Location.valid() {
			return nil, classified(ClassificationBinding, "instance-index-entry")
		}
		return append([]byte(nil), entry.InstanceDigest[:]...), nil
	}, "instance-index-order")
}

func validateApprovalReplayIndex(entries []ApprovalReplayIndexEntry) error {
	return validateStrictEntries(entries, func(entry ApprovalReplayIndexEntry) ([]byte, error) {
		if zeroDigest(entry.PayloadDigest) || zeroDigest(entry.ExactPayloadDigest) || zeroDigest(entry.AuthorizationIdentity) ||
			entry.ApprovalID == (approvalattempt.ApprovalID{}) || !validApprovalState(entry.State) ||
			!entry.Location.valid() {
			return nil, classified(ClassificationBinding, "approval-replay-index-entry")
		}
		key := append([]byte(nil), entry.PayloadDigest[:]...)
		key = append(key, entry.AuthorizationIdentity[:]...)
		return key, nil
	}, "approval-replay-index-order")
}

func validLifecycleOperation(operation lifecyclestate.Operation) bool {
	for _, candidate := range lifecyclestate.Operations() {
		if operation == candidate {
			return true
		}
	}
	return false
}

func validateAttemptReplayIndex(entries []AttemptReplayIndexEntry) error {
	return validateStrictEntries(entries, func(entry AttemptReplayIndexEntry) ([]byte, error) {
		if entry.RegistrationID == (v0candidate.RegistrationID{}) ||
			entry.ApprovalID == (approvalattempt.ApprovalID{}) ||
			entry.AttemptID == (approvalattempt.AttemptID{}) || entry.State != approvalattempt.AttemptCreated {
			return nil, classified(ClassificationBinding, "attempt-replay-index-entry")
		}
		key := append([]byte(nil), entry.RegistrationID[:]...)
		key = append(key, entry.ApprovalID[:]...)
		return key, nil
	}, "attempt-replay-index-order")
}

func validateStrictEntries[T any](entries []T, key func(T) ([]byte, error), code string) error {
	var previous []byte
	for index, entry := range entries {
		current, err := key(entry)
		if err != nil {
			return err
		}
		if index > 0 && bytes.Compare(previous, current) >= 0 {
			return classified(ClassificationBinding, code)
		}
		previous = append(previous[:0], current...)
	}
	return nil
}

func digestArchiveIndexes(view ArchiveIndexesView) ArchiveIndexDigests {
	return ArchiveIndexDigests{
		Registrations:  digestIndex("registration", view.Registrations, encodeRegistrationIndex),
		Approvals:      digestIndex("approval", view.Approvals, encodeApprovalIndex),
		Attempts:       digestIndex("attempt", view.Attempts, encodeAttemptIndex),
		Nonces:         digestIndex("nonce", view.Nonces, encodeNonceIndex),
		Effects:        digestIndex("effect", view.Effects, encodeEffectIndex),
		Instances:      digestIndex("instance", view.Instances, encodeInstanceIndex),
		ApprovalReplay: digestIndex("approval-replay", view.ApprovalReplay, encodeApprovalReplayIndex),
		AttemptReplay:  digestIndex("attempt-replay", view.AttemptReplay, encodeAttemptReplayIndex),
	}
}

func digestIndex[T any](name string, entries []T, encode func(*digestEncoder, T)) ArchiveIndexDigest {
	encoder := newDigestEncoder("capsule.supervisor.archive-index.v0")
	encoder.text(name)
	encoder.uint64(uint64(len(entries)))
	for _, entry := range entries {
		encode(encoder, entry)
	}
	return digestSum[ArchiveIndexDigest](encoder)
}

func digestCombinedIndexes(digests ArchiveIndexDigests) ArchiveCombinedIndexDigest {
	encoder := newDigestEncoder("capsule.supervisor.archive-combined-index.v0")
	for _, digest := range [...][32]byte{
		[32]byte(digests.Registrations), [32]byte(digests.Approvals), [32]byte(digests.Attempts),
		[32]byte(digests.Nonces), [32]byte(digests.Effects), [32]byte(digests.Instances),
		[32]byte(digests.ApprovalReplay), [32]byte(digests.AttemptReplay),
	} {
		encoder.bytes(digest[:])
	}
	return digestSum[ArchiveCombinedIndexDigest](encoder)
}

// DigestVisibleV1EffectSeed binds only the nonzero effect IDs visible at the
// v1-to-v2 boundary. It deliberately cannot claim overwritten pre-v2 history.
func DigestVisibleV1EffectSeed(ids []lifecyclestate.EffectID) (ArchiveEffectSeedDigest, error) {
	if ids == nil {
		return ArchiveEffectSeedDigest{}, classified(ClassificationSchema, "visible-v1-effect-seed-nil")
	}
	encoder := newDigestEncoder("capsule.supervisor.visible-v1-effect-seed.v0")
	encoder.uint64(uint64(len(ids)))
	var previous lifecyclestate.EffectID
	for index, identifier := range ids {
		if identifier.IsZero() || (index > 0 && bytes.Compare(previous[:], identifier[:]) >= 0) {
			return ArchiveEffectSeedDigest{}, classified(ClassificationBinding, "visible-v1-effect-seed-order")
		}
		previous = identifier
		encoder.bytes(identifier[:])
	}
	return digestSum[ArchiveEffectSeedDigest](encoder), nil
}

func effectSeedForIndexes(indexes ArchiveIndexes) (uint64, ArchiveEffectSeedDigest, error) {
	ids := make([]lifecyclestate.EffectID, 0)
	for _, entry := range indexes.view.Effects {
		if entry.VisibleV1Seed {
			ids = append(ids, entry.EffectID)
		}
	}
	digest, err := DigestVisibleV1EffectSeed(ids)
	return uint64(len(ids)), digest, err
}

func encodeLocation(encoder *digestEncoder, location ArchiveLocation) {
	encoder.uint64(uint64(location.SegmentOrdinal))
	encoder.uint64(uint64(location.CohortOrdinal))
	encoder.uint64(uint64(location.RecordOrdinal))
}

func encodeRegistrationIndex(encoder *digestEncoder, entry RegistrationIndexEntry) {
	encoder.bytes(entry.RegistrationID[:])
	encoder.uint64(uint64(entry.RegistrationSequence))
	encoder.bytes(entry.PlanDigest[:])
	encoder.uint64(uint64(entry.ExpiresAt))
	encodeLocation(encoder, entry.Location)
	encoder.bytes(entry.FullRecordDigest[:])
}

func encodeApprovalIndex(encoder *digestEncoder, entry ApprovalIndexEntry) {
	encoder.bytes(entry.ApprovalID[:])
	encoder.bytes(entry.RegistrationID[:])
	encoder.bytes(entry.PayloadDigest[:])
	encoder.bytes(entry.AuthorizationIdentity[:])
	encoder.bytes(entry.AttemptNonce[:])
	encoder.text(string(entry.State))
	encoder.bytes(entry.ConsumedAttemptID[:])
	encoder.uint64(uint64(entry.ExpiresAt))
	encodeLocation(encoder, entry.Location)
	encoder.bytes(entry.FullRecordDigest[:])
}

func encodeAttemptIndex(encoder *digestEncoder, entry AttemptIndexEntry) {
	encoder.bytes(entry.AttemptID[:])
	encoder.bytes(entry.ApprovalID[:])
	encoder.bytes(entry.RegistrationID[:])
	encoder.uint64(uint64(entry.CreatedAt))
	encoder.text(string(entry.LifecycleState))
	encodeLocation(encoder, entry.Location)
	encoder.bytes(entry.FullRecordDigest[:])
}

func encodeNonceIndex(encoder *digestEncoder, entry NonceIndexEntry) {
	encoder.bytes(entry.AttemptNonce[:])
	encoder.bytes(entry.PayloadDigest[:])
	encoder.bytes(entry.ApprovalID[:])
}

func encodeEffectIndex(encoder *digestEncoder, entry EffectIndexEntry) {
	encoder.bytes(entry.EffectID[:])
	encoder.bytes(entry.AttemptID[:])
	encoder.uint64(uint64(entry.OperationSequence))
	encoder.text(string(entry.Operation))
	encoder.uint64(uint64(entry.IssuanceSnapshotGeneration))
	encoder.boolean(entry.VisibleV1Seed)
}

func encodeInstanceIndex(encoder *digestEncoder, entry InstanceIndexEntry) {
	encoder.bytes(entry.InstanceDigest[:])
	encoder.bytes(entry.AttemptID[:])
	encodeLocation(encoder, entry.Location)
}

func encodeApprovalReplayIndex(encoder *digestEncoder, entry ApprovalReplayIndexEntry) {
	encoder.bytes(entry.PayloadDigest[:])
	encoder.bytes(entry.ExactPayloadDigest[:])
	encoder.bytes(entry.AuthorizationIdentity[:])
	encoder.bytes(entry.ApprovalID[:])
	encoder.text(string(entry.State))
	encodeLocation(encoder, entry.Location)
}

func encodeAttemptReplayIndex(encoder *digestEncoder, entry AttemptReplayIndexEntry) {
	encoder.bytes(entry.RegistrationID[:])
	encoder.bytes(entry.ApprovalID[:])
	encoder.bytes(entry.AttemptID[:])
	encoder.text(string(entry.State))
}
