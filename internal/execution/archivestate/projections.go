package archivestate

import (
	"bytes"

	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

type ArchiveRecordReference struct {
	Kind   RecordKind
	Digest ArchiveRecordDigest
}

func (reference ArchiveRecordReference) valid(kind RecordKind) bool {
	return reference.Kind == kind && !zeroDigest(reference.Digest)
}

// CohortProjectionView binds one complete registration cohort without
// carrying a path or mutable store handle.
type CohortProjectionView struct {
	RegistrationID       v0candidate.RegistrationID
	RegistrationSequence v0candidate.PositiveUInt53
	ExpiresAt            v0candidate.UInt53
	Registration         ArchiveRecordReference
	Approvals            []ArchiveRecordReference
	Attempts             []ArchiveRecordReference
	Lifecycles           []ArchiveRecordReference
	EncodedByteLength    uint64
}

// CohortProjection is an immutable, validated, indivisible segment unit.
type CohortProjection struct {
	view   CohortProjectionView
	digest ArchiveCohortDigest
}

func NewCohortProjection(view CohortProjectionView) (CohortProjection, error) {
	if err := validateCohortProjectionView(view); err != nil {
		return CohortProjection{}, err
	}
	cloned := cloneCohortProjectionView(view)
	return CohortProjection{view: cloned, digest: digestCohort(cloned)}, nil
}

func (cohort CohortProjection) View() CohortProjectionView {
	return cloneCohortProjectionView(cohort.view)
}

func (cohort CohortProjection) Digest() ArchiveCohortDigest { return cohort.digest }

func cloneCohortProjectionView(view CohortProjectionView) CohortProjectionView {
	view.Approvals = cloneSlice(view.Approvals)
	view.Attempts = cloneSlice(view.Attempts)
	view.Lifecycles = cloneSlice(view.Lifecycles)
	return view
}

func validateCohortProjectionView(view CohortProjectionView) error {
	if view.RegistrationID == (v0candidate.RegistrationID{}) ||
		!validPositive(uint64(view.RegistrationSequence)) ||
		uint64(view.ExpiresAt) > v0candidate.MaxSafeInteger ||
		!view.Registration.valid(RecordRegistration) || view.Approvals == nil ||
		view.Attempts == nil || view.Lifecycles == nil || view.EncodedByteLength == 0 ||
		view.EncodedByteLength > uint64(MaxSupervisorArchiveBytes) {
		return classified(ClassificationBinding, "cohort-projection")
	}
	if len(view.Approvals) > MaxArchiveApprovalsPerSegment ||
		len(view.Attempts) > MaxArchiveAttemptsPerSegment ||
		len(view.Lifecycles) > MaxArchiveLifecyclesPerSegment {
		return classified(ClassificationCapacity, "cohort-record-count")
	}
	if len(view.Attempts) != len(view.Lifecycles) {
		return classified(ClassificationBinding, "cohort-attempt-lifecycle-count")
	}
	for _, group := range []struct {
		kind    RecordKind
		records []ArchiveRecordReference
	}{
		{RecordApproval, view.Approvals}, {RecordAttempt, view.Attempts}, {RecordLifecycle, view.Lifecycles},
	} {
		seen := make(map[ArchiveRecordDigest]struct{}, len(group.records))
		for _, reference := range group.records {
			if !reference.valid(group.kind) {
				return classified(ClassificationDomain, "cohort-record-domain")
			}
			if _, exists := seen[reference.Digest]; exists {
				return classified(ClassificationBinding, "cohort-record-collision")
			}
			seen[reference.Digest] = struct{}{}
		}
	}
	return nil
}

func digestCohort(view CohortProjectionView) ArchiveCohortDigest {
	encoder := newDigestEncoder("capsule.supervisor.archive-cohort.v0")
	encoder.bytes(view.RegistrationID[:])
	encoder.uint64(uint64(view.RegistrationSequence))
	encoder.uint64(uint64(view.ExpiresAt))
	encodeRecordReference(encoder, view.Registration)
	for _, records := range [][]ArchiveRecordReference{view.Approvals, view.Attempts, view.Lifecycles} {
		encoder.uint64(uint64(len(records)))
		for _, record := range records {
			encodeRecordReference(encoder, record)
		}
	}
	encoder.uint64(view.EncodedByteLength)
	return digestSum[ArchiveCohortDigest](encoder)
}

func encodeRecordReference(encoder *digestEncoder, reference ArchiveRecordReference) {
	encoder.text(string(reference.Kind))
	encoder.bytes(reference.Digest[:])
}

// ArchiveCounts records exact collection and index counts. In active v2,
// hot plus archived equals total, retained-global indexes equal total, and
// record-location arms independently equal hot and archived. It does not imply
// deletion or a continuous-service budget.
type ArchiveCounts struct {
	Registrations  uint64
	Approvals      uint64
	Attempts       uint64
	Lifecycles     uint64
	Nonces         uint64
	Effects        uint64
	Instances      uint64
	ApprovalReplay uint64
	AttemptReplay  uint64
}

func (counts ArchiveCounts) add(other ArchiveCounts) (ArchiveCounts, error) {
	values := []*uint64{
		&counts.Registrations, &counts.Approvals, &counts.Attempts, &counts.Lifecycles,
		&counts.Nonces, &counts.Effects, &counts.Instances, &counts.ApprovalReplay, &counts.AttemptReplay,
	}
	otherValues := [...]uint64{
		other.Registrations, other.Approvals, other.Attempts, other.Lifecycles,
		other.Nonces, other.Effects, other.Instances, other.ApprovalReplay, other.AttemptReplay,
	}
	for index, value := range values {
		if ^uint64(0)-*value < otherValues[index] {
			return ArchiveCounts{}, classified(ClassificationCapacity, "archive-count-overflow")
		}
		*value += otherValues[index]
	}
	return counts, nil
}

// ArchiveSetDigests binds all full-record and derived projection sets in a
// segment. These values are integrity checkpoints, not signatures.
type ArchiveSetDigests struct {
	Registrations  ArchiveIndexDigest
	Approvals      ArchiveIndexDigest
	Attempts       ArchiveIndexDigest
	Lifecycles     ArchiveIndexDigest
	Nonces         ArchiveIndexDigest
	Effects        ArchiveIndexDigest
	Instances      ArchiveIndexDigest
	ApprovalReplay ArchiveIndexDigest
	AttemptReplay  ArchiveIndexDigest
}

// ArchiveSegmentView is the complete passive archive-segment v0 projection.
type ArchiveSegmentView struct {
	FormatVersion            uint64
	Ordinal                  ArchiveOrdinal
	InstallationID           v0candidate.InstallationID
	SupervisorID             v0candidate.SupervisorID
	EpochSequence            v0candidate.UInt53
	EpochDigest              v0candidate.TrustEpochDigest
	SourceSnapshotGeneration v0candidate.PositiveUInt53
	ArchiveGeneration        ArchiveGeneration
	DurableTimeHighWater     v0candidate.UInt53
	PriorCheckpointDigest    ArchiveCheckpointDigest
	Cohorts                  []CohortProjection
	DerivedIndexes           ArchiveIndexes
	SetDigests               ArchiveSetDigests
	Counts                   ArchiveCounts
	EncodedByteLength        uint64
}

// ArchiveSegment is immutable and contains no file name or path.
type ArchiveSegment struct {
	view   ArchiveSegmentView
	digest ArchiveSegmentDigest
}

func NewArchiveSegment(view ArchiveSegmentView) (ArchiveSegment, error) {
	if err := validateArchiveSegmentView(view); err != nil {
		return ArchiveSegment{}, err
	}
	cloned := cloneArchiveSegmentView(view)
	derivedIndexes, err := NewArchiveIndexes(cloned.DerivedIndexes.View())
	if err != nil {
		return ArchiveSegment{}, err
	}
	cloned.DerivedIndexes = derivedIndexes
	return ArchiveSegment{view: cloned, digest: digestArchiveSegment(cloned)}, nil
}

func (segment ArchiveSegment) View() ArchiveSegmentView { return cloneArchiveSegmentView(segment.view) }

func (segment ArchiveSegment) Digest() ArchiveSegmentDigest { return segment.digest }

func cloneArchiveSegmentView(view ArchiveSegmentView) ArchiveSegmentView {
	view.Cohorts = cloneSlice(view.Cohorts)
	return view
}

func validateArchiveSegmentView(view ArchiveSegmentView) error {
	if view.FormatVersion != SupervisorArchiveFormatV0 {
		return classified(ClassificationUnsupported, "archive-segment-version")
	}
	if !validPositive(uint64(view.Ordinal)) || view.InstallationID == (v0candidate.InstallationID{}) ||
		view.SupervisorID == (v0candidate.SupervisorID{}) ||
		uint64(view.EpochSequence) > v0candidate.MaxSafeInteger ||
		!validPositive(uint64(view.SourceSnapshotGeneration)) ||
		!validPositive(uint64(view.ArchiveGeneration)) ||
		uint64(view.DurableTimeHighWater) > v0candidate.MaxSafeInteger || view.Cohorts == nil ||
		view.EncodedByteLength == 0 || view.EncodedByteLength > uint64(MaxSupervisorArchiveBytes) {
		return classified(ClassificationBinding, "archive-segment-metadata")
	}
	if len(view.Cohorts) > MaxArchiveCohortsPerSegment {
		return classified(ClassificationCapacity, "archive-segment-cohorts")
	}
	derivedIndexes, err := NewArchiveIndexes(view.DerivedIndexes.View())
	if err != nil {
		return err
	}
	if derivedIndexes.scope() != ArchiveIndexScopeSegmentDerived ||
		!derivedIndexes.allLocationsAre(RecordLocationArchive) {
		return classified(ClassificationDomain, "archive-segment-derived-index-scope")
	}
	derivedView := derivedIndexes.View()
	if len(derivedView.Registrations) != 0 || len(derivedView.Approvals) != 0 || len(derivedView.Attempts) != 0 {
		return classified(ClassificationDomain, "archive-segment-derived-index-domain")
	}
	var previous CohortProjectionView
	for index, cohort := range view.Cohorts {
		cohortView := cohort.View()
		if err := validateCohortProjectionView(cohortView); err != nil || cohort.Digest() != digestCohort(cohortView) {
			return classified(ClassificationBinding, "archive-segment-cohort")
		}
		if index > 0 && compareCohortIdentity(previous, cohortView) >= 0 {
			return classified(ClassificationBinding, "archive-segment-cohort-order")
		}
		previous = cohortView
	}
	wantCounts := countsForCohorts(view.Cohorts)
	derivedCounts := derivedIndexes.counts()
	wantCounts.Nonces = derivedCounts.Nonces
	wantCounts.Effects = derivedCounts.Effects
	wantCounts.Instances = derivedCounts.Instances
	wantCounts.ApprovalReplay = derivedCounts.ApprovalReplay
	wantCounts.AttemptReplay = derivedCounts.AttemptReplay
	if view.Counts.Registrations != wantCounts.Registrations ||
		view.Counts.Approvals != wantCounts.Approvals || view.Counts.Attempts != wantCounts.Attempts ||
		view.Counts.Lifecycles != wantCounts.Lifecycles {
		return classified(ClassificationBinding, "archive-segment-counts")
	}
	if view.Counts.Registrations > MaxArchiveCohortsPerSegment ||
		view.Counts.Approvals > MaxArchiveApprovalsPerSegment ||
		view.Counts.Attempts > MaxArchiveAttemptsPerSegment ||
		view.Counts.Lifecycles > MaxArchiveLifecyclesPerSegment {
		return classified(ClassificationCapacity, "archive-segment-record-counts")
	}
	if wantDigests := archiveSetDigests(view.Cohorts, derivedIndexes); wantDigests != view.SetDigests {
		return classified(ClassificationBinding, "archive-segment-set-digests")
	}
	return nil
}

func compareCohortIdentity(left, right CohortProjectionView) int {
	if left.RegistrationSequence < right.RegistrationSequence {
		return -1
	}
	if left.RegistrationSequence > right.RegistrationSequence {
		return 1
	}
	return bytes.Compare(left.RegistrationID[:], right.RegistrationID[:])
}

func countsForCohorts(cohorts []CohortProjection) ArchiveCounts {
	counts := ArchiveCounts{Registrations: uint64(len(cohorts))}
	for _, cohort := range cohorts {
		view := cohort.View()
		counts.Approvals += uint64(len(view.Approvals))
		counts.Attempts += uint64(len(view.Attempts))
		counts.Lifecycles += uint64(len(view.Lifecycles))
	}
	return counts
}

func digestArchiveSegment(view ArchiveSegmentView) ArchiveSegmentDigest {
	encoder := newDigestEncoder("capsule.supervisor.archive-segment.v0")
	encoder.uint64(view.FormatVersion)
	encoder.uint64(uint64(view.Ordinal))
	encoder.bytes(view.InstallationID[:])
	encoder.bytes(view.SupervisorID[:])
	encoder.uint64(uint64(view.EpochSequence))
	encoder.bytes(view.EpochDigest[:])
	encoder.uint64(uint64(view.SourceSnapshotGeneration))
	encoder.uint64(uint64(view.ArchiveGeneration))
	encoder.uint64(uint64(view.DurableTimeHighWater))
	encoder.bytes(view.PriorCheckpointDigest[:])
	encoder.uint64(uint64(len(view.Cohorts)))
	for _, cohort := range view.Cohorts {
		digest := cohort.Digest()
		encoder.bytes(digest[:])
	}
	encodeSetDigests(encoder, view.SetDigests)
	encodeCounts(encoder, view.Counts)
	encoder.uint64(view.EncodedByteLength)
	return digestSum[ArchiveSegmentDigest](encoder)
}

func archiveSetDigests(cohorts []CohortProjection, derived ArchiveIndexes) ArchiveSetDigests {
	digests := derived.Digests()
	return ArchiveSetDigests{
		Registrations:  digestCohortRecordSet("registration", cohorts, RecordRegistration),
		Approvals:      digestCohortRecordSet("approval", cohorts, RecordApproval),
		Attempts:       digestCohortRecordSet("attempt", cohorts, RecordAttempt),
		Lifecycles:     digestCohortRecordSet("lifecycle", cohorts, RecordLifecycle),
		Nonces:         digests.Nonces,
		Effects:        digests.Effects,
		Instances:      digests.Instances,
		ApprovalReplay: digests.ApprovalReplay,
		AttemptReplay:  digests.AttemptReplay,
	}
}

// DeriveArchiveSegmentSets returns the exact record/projection digests and
// counts that NewArchiveSegment requires. It preserves the package's domain-
// separated digest implementation while allowing the fixed-store adapter to
// serialize full records without duplicating digest primitives.
func DeriveArchiveSegmentSets(
	cohorts []CohortProjection,
	derived ArchiveIndexes,
) (ArchiveSetDigests, ArchiveCounts, error) {
	if cohorts == nil {
		return ArchiveSetDigests{}, ArchiveCounts{}, classified(ClassificationSchema, "archive-segment-cohorts-nil")
	}
	validated, err := NewArchiveIndexes(derived.View())
	if err != nil {
		return ArchiveSetDigests{}, ArchiveCounts{}, err
	}
	if validated.scope() != ArchiveIndexScopeSegmentDerived ||
		!validated.allLocationsAre(RecordLocationArchive) {
		return ArchiveSetDigests{}, ArchiveCounts{}, classified(ClassificationDomain, "archive-segment-derived-index-scope")
	}
	for _, cohort := range cohorts {
		if err := validateCohortProjectionView(cohort.View()); err != nil {
			return ArchiveSetDigests{}, ArchiveCounts{}, err
		}
	}
	counts := countsForCohorts(cohorts)
	derivedCounts := validated.counts()
	counts.Nonces = derivedCounts.Nonces
	counts.Effects = derivedCounts.Effects
	counts.Instances = derivedCounts.Instances
	counts.ApprovalReplay = derivedCounts.ApprovalReplay
	counts.AttemptReplay = derivedCounts.AttemptReplay
	return archiveSetDigests(cohorts, validated), counts, nil
}

func digestCohortRecordSet(name string, cohorts []CohortProjection, kind RecordKind) ArchiveIndexDigest {
	encoder := newDigestEncoder("capsule.supervisor.archive-segment-set.v0")
	encoder.text(name)
	count := uint64(0)
	for _, cohort := range cohorts {
		view := cohort.View()
		switch kind {
		case RecordRegistration:
			count++
		case RecordApproval:
			count += uint64(len(view.Approvals))
		case RecordAttempt:
			count += uint64(len(view.Attempts))
		case RecordLifecycle:
			count += uint64(len(view.Lifecycles))
		}
	}
	encoder.uint64(count)
	for _, cohort := range cohorts {
		view := cohort.View()
		var records []ArchiveRecordReference
		switch kind {
		case RecordRegistration:
			records = []ArchiveRecordReference{view.Registration}
		case RecordApproval:
			records = view.Approvals
		case RecordAttempt:
			records = view.Attempts
		case RecordLifecycle:
			records = view.Lifecycles
		}
		for _, record := range records {
			encoder.bytes(record.Digest[:])
		}
	}
	return digestSum[ArchiveIndexDigest](encoder)
}

func encodeSetDigests(encoder *digestEncoder, digests ArchiveSetDigests) {
	for _, digest := range [...][32]byte{
		[32]byte(digests.Registrations), [32]byte(digests.Approvals), [32]byte(digests.Attempts),
		[32]byte(digests.Lifecycles), [32]byte(digests.Nonces), [32]byte(digests.Effects),
		[32]byte(digests.Instances), [32]byte(digests.ApprovalReplay), [32]byte(digests.AttemptReplay),
	} {
		encoder.bytes(digest[:])
	}
}

func encodeCounts(encoder *digestEncoder, counts ArchiveCounts) {
	for _, count := range [...]uint64{
		counts.Registrations, counts.Approvals, counts.Attempts, counts.Lifecycles, counts.Nonces,
		counts.Effects, counts.Instances, counts.ApprovalReplay, counts.AttemptReplay,
	} {
		encoder.uint64(count)
	}
}

// ArchiveDescriptorView is the active-state projection of one referenced
// immutable segment.
type ArchiveDescriptorView struct {
	Ordinal                  ArchiveOrdinal
	SourceSnapshotGeneration v0candidate.PositiveUInt53
	ArchiveGeneration        ArchiveGeneration
	SegmentDigest            ArchiveSegmentDigest
	EncodedByteLength        uint64
	Counts                   ArchiveCounts
}

type ArchiveDescriptor struct {
	view   ArchiveDescriptorView
	digest ArchiveIndexDigest
}

func NewArchiveDescriptor(view ArchiveDescriptorView) (ArchiveDescriptor, error) {
	if !validPositive(uint64(view.Ordinal)) || !validPositive(uint64(view.SourceSnapshotGeneration)) ||
		!validPositive(uint64(view.ArchiveGeneration)) || zeroDigest(view.SegmentDigest) ||
		view.EncodedByteLength == 0 || view.EncodedByteLength > uint64(MaxSupervisorArchiveBytes) {
		return ArchiveDescriptor{}, classified(ClassificationBinding, "archive-descriptor")
	}
	if view.Counts.Registrations > MaxArchiveCohortsPerSegment ||
		view.Counts.Approvals > MaxArchiveApprovalsPerSegment ||
		view.Counts.Attempts > MaxArchiveAttemptsPerSegment ||
		view.Counts.Lifecycles > MaxArchiveLifecyclesPerSegment ||
		view.Counts.Nonces > MaxArchiveIndexEntries || view.Counts.Effects > MaxArchiveIndexEntries ||
		view.Counts.Instances > MaxArchiveIndexEntries || view.Counts.ApprovalReplay > MaxArchiveIndexEntries ||
		view.Counts.AttemptReplay > MaxArchiveIndexEntries {
		return ArchiveDescriptor{}, classified(ClassificationCapacity, "archive-descriptor-counts")
	}
	encoder := newDigestEncoder("capsule.supervisor.archive-descriptor.v0")
	encodeDescriptorView(encoder, view)
	return ArchiveDescriptor{view: view, digest: digestSum[ArchiveIndexDigest](encoder)}, nil
}

func (descriptor ArchiveDescriptor) View() ArchiveDescriptorView { return descriptor.view }

func (descriptor ArchiveDescriptor) Digest() ArchiveIndexDigest { return descriptor.digest }

func encodeDescriptorView(encoder *digestEncoder, view ArchiveDescriptorView) {
	encoder.uint64(uint64(view.Ordinal))
	encoder.uint64(uint64(view.SourceSnapshotGeneration))
	encoder.uint64(uint64(view.ArchiveGeneration))
	encoder.bytes(view.SegmentDigest[:])
	encoder.uint64(view.EncodedByteLength)
	encodeCounts(encoder, view.Counts)
}

func DigestArchiveDescriptorSet(descriptors []ArchiveDescriptor) (ArchiveDescriptorSetDigest, error) {
	if descriptors == nil {
		return ArchiveDescriptorSetDigest{}, classified(ClassificationSchema, "archive-descriptor-set-nil")
	}
	if len(descriptors) > MaxReferencedArchiveSegments {
		return ArchiveDescriptorSetDigest{}, classified(ClassificationCapacity, "archive-descriptor-count")
	}
	encoder := newDigestEncoder("capsule.supervisor.archive-descriptor-set.v0")
	encoder.uint64(uint64(len(descriptors)))
	var previous ArchiveOrdinal
	seenSegments := make(map[ArchiveSegmentDigest]struct{}, len(descriptors))
	for index, descriptor := range descriptors {
		view := descriptor.View()
		if index > 0 && previous >= view.Ordinal {
			return ArchiveDescriptorSetDigest{}, classified(ClassificationBinding, "archive-descriptor-order")
		}
		if _, exists := seenSegments[view.SegmentDigest]; exists {
			return ArchiveDescriptorSetDigest{}, classified(ClassificationBinding, "archive-descriptor-collision")
		}
		seenSegments[view.SegmentDigest] = struct{}{}
		previous = view.Ordinal
		digest := descriptor.Digest()
		encoder.bytes(digest[:])
	}
	return digestSum[ArchiveDescriptorSetDigest](encoder), nil
}

// HotSetDigests binds the current mutable record collections into a checkpoint.
type HotSetDigests struct {
	Registrations [32]byte
	Approvals     [32]byte
	Attempts      [32]byte
	Lifecycles    [32]byte
	// EffectTombstones is deliberately excluded from the legacy F2/F3 JSON
	// representation. Exact predecessor checkpoint bytes retain zero here;
	// F4B active worlds bind the separately persisted nonzero collection digest.
	EffectTombstones EffectTombstoneSetDigest `json:"-"`
}

func validHotSetDigests(digests HotSetDigests) bool {
	return digests.Registrations != ([32]byte{}) && digests.Approvals != ([32]byte{}) &&
		digests.Attempts != ([32]byte{}) && digests.Lifecycles != ([32]byte{})
}

func validF4BHotSetDigests(digests HotSetDigests) bool {
	return validHotSetDigests(digests) && digests.EffectTombstones != (EffectTombstoneSetDigest{})
}

// ArchiveCheckpointKind prevents a migration genesis from being interpreted
// as an activation checkpoint (or the reverse).
type ArchiveCheckpointKind string

const (
	ArchiveCheckpointNone             ArchiveCheckpointKind = "none"
	ArchiveCheckpointMigrationGenesis ArchiveCheckpointKind = "migration-genesis"
	ArchiveCheckpointActivation       ArchiveCheckpointKind = "activation"
)

type ArchiveCheckpointReference struct {
	Kind   ArchiveCheckpointKind
	Digest ArchiveCheckpointDigest
}

func NoArchiveCheckpointReference() ArchiveCheckpointReference {
	return ArchiveCheckpointReference{Kind: ArchiveCheckpointNone}
}

func (reference ArchiveCheckpointReference) valid(allowNone bool) bool {
	switch reference.Kind {
	case ArchiveCheckpointNone:
		return allowNone && zeroDigest(reference.Digest)
	case ArchiveCheckpointMigrationGenesis, ArchiveCheckpointActivation:
		return !zeroDigest(reference.Digest)
	default:
		return false
	}
}

// MigrationGenesisCheckpointView is the generation-one v1-to-v2 checkpoint.
// Unlike an activation checkpoint it binds no previous checkpoint or segment;
// it binds the complete all-hot retained indexes and the exact visible-v1 seed.
type MigrationGenesisCheckpointView struct {
	StoreFormatVersion       uint64
	MigrationSourceVersion   uint64
	ResultSnapshotGeneration v0candidate.PositiveUInt53
	ArchiveGeneration        ArchiveGeneration
	DescriptorSetDigest      ArchiveDescriptorSetDigest
	Indexes                  ArchiveIndexes
	HotSetDigests            HotSetDigests
	VisibleV1EffectSeed      []lifecyclestate.EffectID
	HotCounts                ArchiveCounts
	InstallationID           v0candidate.InstallationID
	SupervisorID             v0candidate.SupervisorID
	EpochSequence            v0candidate.UInt53
	EpochDigest              v0candidate.TrustEpochDigest
	DurableTimeHighWater     v0candidate.UInt53
}

type MigrationGenesisCheckpoint struct {
	view       MigrationGenesisCheckpointView
	seedDigest ArchiveEffectSeedDigest
	digest     ArchiveCheckpointDigest
}

func NewMigrationGenesisCheckpoint(view MigrationGenesisCheckpointView) (MigrationGenesisCheckpoint, error) {
	if view.StoreFormatVersion != SupervisorStoreFormatV2 ||
		view.MigrationSourceVersion != MigrationSourceFormatV1 {
		return MigrationGenesisCheckpoint{}, classified(ClassificationUnsupported, "migration-genesis-version")
	}
	if view.ResultSnapshotGeneration != 1 || view.ArchiveGeneration != 1 ||
		view.InstallationID == (v0candidate.InstallationID{}) ||
		view.SupervisorID == (v0candidate.SupervisorID{}) ||
		uint64(view.EpochSequence) > v0candidate.MaxSafeInteger ||
		uint64(view.DurableTimeHighWater) > v0candidate.MaxSafeInteger ||
		!validHotSetDigests(view.HotSetDigests) || view.VisibleV1EffectSeed == nil {
		return MigrationGenesisCheckpoint{}, classified(ClassificationBinding, "migration-genesis-metadata")
	}
	emptyDescriptorDigest, err := DigestArchiveDescriptorSet([]ArchiveDescriptor{})
	if err != nil || view.DescriptorSetDigest != emptyDescriptorDigest {
		return MigrationGenesisCheckpoint{}, classified(ClassificationBinding, "migration-genesis-descriptors")
	}
	indexes, err := NewArchiveIndexes(view.Indexes.View())
	if err != nil {
		return MigrationGenesisCheckpoint{}, err
	}
	if indexes.scope() != ArchiveIndexScopeRetainedGlobal ||
		!indexes.allLocationsAre(RecordLocationHot) {
		return MigrationGenesisCheckpoint{}, classified(ClassificationDomain, "migration-genesis-index-scope")
	}
	if indexes.counts() != view.HotCounts {
		return MigrationGenesisCheckpoint{}, classified(ClassificationBinding, "migration-genesis-counts")
	}
	seed := cloneSlice(view.VisibleV1EffectSeed)
	seedDigest, err := DigestVisibleV1EffectSeed(seed)
	if err != nil {
		return MigrationGenesisCheckpoint{}, err
	}
	indexedSeed := indexes.visibleV1EffectSeed()
	if len(seed) != len(indexedSeed) {
		return MigrationGenesisCheckpoint{}, classified(ClassificationBinding, "migration-genesis-effect-seed")
	}
	for index := range seed {
		if seed[index] != indexedSeed[index] {
			return MigrationGenesisCheckpoint{}, classified(ClassificationBinding, "migration-genesis-effect-seed")
		}
	}
	cloned := view
	cloned.Indexes = indexes
	cloned.VisibleV1EffectSeed = seed
	encoder := newDigestEncoder("capsule.supervisor.archive-migration-genesis.v0")
	encodeMigrationGenesisView(encoder, cloned, seedDigest)
	return MigrationGenesisCheckpoint{view: cloned, seedDigest: seedDigest, digest: digestSum[ArchiveCheckpointDigest](encoder)}, nil
}

func (checkpoint MigrationGenesisCheckpoint) View() MigrationGenesisCheckpointView {
	view := checkpoint.view
	view.Indexes = checkpoint.view.Indexes
	view.VisibleV1EffectSeed = cloneSlice(checkpoint.view.VisibleV1EffectSeed)
	return view
}

func (checkpoint MigrationGenesisCheckpoint) Digest() ArchiveCheckpointDigest {
	return checkpoint.digest
}

func (checkpoint MigrationGenesisCheckpoint) SeedDigest() ArchiveEffectSeedDigest {
	return checkpoint.seedDigest
}

func (checkpoint MigrationGenesisCheckpoint) Reference() ArchiveCheckpointReference {
	return ArchiveCheckpointReference{Kind: ArchiveCheckpointMigrationGenesis, Digest: checkpoint.digest}
}

func encodeMigrationGenesisView(encoder *digestEncoder, view MigrationGenesisCheckpointView, seedDigest ArchiveEffectSeedDigest) {
	encoder.uint64(view.StoreFormatVersion)
	encoder.uint64(view.MigrationSourceVersion)
	encoder.uint64(uint64(view.ResultSnapshotGeneration))
	encoder.uint64(uint64(view.ArchiveGeneration))
	encoder.bytes(view.DescriptorSetDigest[:])
	indexDigests := view.Indexes.Digests()
	for _, digest := range [...][32]byte{
		[32]byte(indexDigests.Registrations), [32]byte(indexDigests.Approvals), [32]byte(indexDigests.Attempts),
		[32]byte(indexDigests.Nonces), [32]byte(indexDigests.Effects), [32]byte(indexDigests.Instances),
		[32]byte(indexDigests.ApprovalReplay), [32]byte(indexDigests.AttemptReplay),
	} {
		encoder.bytes(digest[:])
	}
	combined := view.Indexes.CombinedDigest()
	encoder.bytes(combined[:])
	encoder.bytes(view.HotSetDigests.Registrations[:])
	encoder.bytes(view.HotSetDigests.Approvals[:])
	encoder.bytes(view.HotSetDigests.Attempts[:])
	encoder.bytes(view.HotSetDigests.Lifecycles[:])
	encoder.uint64(uint64(len(view.VisibleV1EffectSeed)))
	encoder.bytes(seedDigest[:])
	encodeCounts(encoder, view.HotCounts)
	encoder.bytes(view.InstallationID[:])
	encoder.bytes(view.SupervisorID[:])
	encoder.uint64(uint64(view.EpochSequence))
	encoder.bytes(view.EpochDigest[:])
	encoder.uint64(uint64(view.DurableTimeHighWater))
}

type ArchiveCheckpointView struct {
	PreviousCheckpoint       ArchiveCheckpointReference
	NewSegmentDigest         ArchiveSegmentDigest
	DescriptorSetDigest      ArchiveDescriptorSetDigest
	ArchiveIndexDigest       ArchiveCombinedIndexDigest
	HotSetDigests            HotSetDigests
	SourceSnapshotGeneration v0candidate.PositiveUInt53
	ResultSnapshotGeneration v0candidate.PositiveUInt53
	ArchiveGeneration        ArchiveGeneration
	InstallationID           v0candidate.InstallationID
	SupervisorID             v0candidate.SupervisorID
	EpochSequence            v0candidate.UInt53
	EpochDigest              v0candidate.TrustEpochDigest
	DurableTimeHighWater     v0candidate.UInt53
}

type ArchiveCheckpoint struct {
	view   ArchiveCheckpointView
	digest ArchiveCheckpointDigest
}

func NewArchiveCheckpoint(view ArchiveCheckpointView) (ArchiveCheckpoint, error) {
	if !validPositive(uint64(view.SourceSnapshotGeneration)) ||
		!validPositive(uint64(view.ResultSnapshotGeneration)) ||
		view.ResultSnapshotGeneration <= view.SourceSnapshotGeneration ||
		!validPositive(uint64(view.ArchiveGeneration)) ||
		view.InstallationID == (v0candidate.InstallationID{}) ||
		view.SupervisorID == (v0candidate.SupervisorID{}) ||
		uint64(view.EpochSequence) > v0candidate.MaxSafeInteger ||
		uint64(view.DurableTimeHighWater) > v0candidate.MaxSafeInteger ||
		!view.PreviousCheckpoint.valid(false) || zeroDigest(view.NewSegmentDigest) ||
		zeroDigest(view.DescriptorSetDigest) || zeroDigest(view.ArchiveIndexDigest) ||
		!validHotSetDigests(view.HotSetDigests) {
		return ArchiveCheckpoint{}, classified(ClassificationBinding, "archive-checkpoint")
	}
	encoder := newDigestEncoder("capsule.supervisor.archive-checkpoint.v0")
	encodeCheckpointView(encoder, view)
	return ArchiveCheckpoint{view: view, digest: digestSum[ArchiveCheckpointDigest](encoder)}, nil
}

func (checkpoint ArchiveCheckpoint) View() ArchiveCheckpointView { return checkpoint.view }

func (checkpoint ArchiveCheckpoint) Digest() ArchiveCheckpointDigest { return checkpoint.digest }

func (checkpoint ArchiveCheckpoint) Reference() ArchiveCheckpointReference {
	return ArchiveCheckpointReference{Kind: ArchiveCheckpointActivation, Digest: checkpoint.digest}
}

func encodeCheckpointView(encoder *digestEncoder, view ArchiveCheckpointView) {
	encoder.text(string(view.PreviousCheckpoint.Kind))
	encoder.bytes(view.PreviousCheckpoint.Digest[:])
	encoder.bytes(view.NewSegmentDigest[:])
	encoder.bytes(view.DescriptorSetDigest[:])
	encoder.bytes(view.ArchiveIndexDigest[:])
	encoder.bytes(view.HotSetDigests.Registrations[:])
	encoder.bytes(view.HotSetDigests.Approvals[:])
	encoder.bytes(view.HotSetDigests.Attempts[:])
	encoder.bytes(view.HotSetDigests.Lifecycles[:])
	// Retain the exact legacy F3 checkpoint answer when the independent F4B
	// tombstone source was not materialized. Every F4C checkpoint carries and
	// domain-binds the fifth hot-set digest.
	if view.HotSetDigests.EffectTombstones != (EffectTombstoneSetDigest{}) {
		encoder.bytes(view.HotSetDigests.EffectTombstones[:])
	}
	encoder.uint64(uint64(view.SourceSnapshotGeneration))
	encoder.uint64(uint64(view.ResultSnapshotGeneration))
	encoder.uint64(uint64(view.ArchiveGeneration))
	encoder.bytes(view.InstallationID[:])
	encoder.bytes(view.SupervisorID[:])
	encoder.uint64(uint64(view.EpochSequence))
	encoder.bytes(view.EpochDigest[:])
	encoder.uint64(uint64(view.DurableTimeHighWater))
}

// ActiveStateV2View is the passive archive-bearing portion of a v2 snapshot.
// Hot authority collections remain owned by registrationstate in F2.
type ActiveStateV2View struct {
	StoreFormatVersion        uint64
	MigrationSourceVersion    uint64
	SnapshotGeneration        v0candidate.PositiveUInt53
	ArchiveGeneration         ArchiveGeneration
	InstallationID            v0candidate.InstallationID
	SupervisorID              v0candidate.SupervisorID
	EpochSequence             v0candidate.UInt53
	EpochDigest               v0candidate.TrustEpochDigest
	DurableTimeHighWater      v0candidate.UInt53
	Descriptors               []ArchiveDescriptor
	DescriptorSetDigest       ArchiveDescriptorSetDigest
	Indexes                   ArchiveIndexes
	IndexDigests              ArchiveIndexDigests
	CombinedIndexDigest       ArchiveCombinedIndexDigest
	HotSetDigests             HotSetDigests
	EffectTombstoneCoverage   string
	VisibleV1EffectSeedCount  uint64
	VisibleV1EffectSeedDigest ArchiveEffectSeedDigest
	PreviousCheckpoint        ArchiveCheckpointReference
	CurrentCheckpoint         ArchiveCheckpointReference
	HotCounts                 ArchiveCounts
	ArchivedCounts            ArchiveCounts
	TotalCounts               ArchiveCounts
}

type ActiveStateV2 struct{ view ActiveStateV2View }

func NewActiveStateV2(view ActiveStateV2View) (ActiveStateV2, error) {
	if view.StoreFormatVersion != SupervisorStoreFormatV2 ||
		view.MigrationSourceVersion != MigrationSourceFormatV1 {
		return ActiveStateV2{}, classified(ClassificationUnsupported, "active-v2-version")
	}
	if !validPositive(uint64(view.SnapshotGeneration)) || !validPositive(uint64(view.ArchiveGeneration)) ||
		view.InstallationID == (v0candidate.InstallationID{}) || view.SupervisorID == (v0candidate.SupervisorID{}) ||
		uint64(view.EpochSequence) > v0candidate.MaxSafeInteger ||
		uint64(view.DurableTimeHighWater) > v0candidate.MaxSafeInteger || view.Descriptors == nil ||
		view.EffectTombstoneCoverage != EffectTombstoneCoverageVisibleV1AndV2 {
		return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-metadata")
	}
	descriptorDigest, err := DigestArchiveDescriptorSet(view.Descriptors)
	if err != nil || descriptorDigest != view.DescriptorSetDigest {
		return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-descriptor-digest")
	}
	indexes, indexErr := NewArchiveIndexes(view.Indexes.View())
	if indexErr != nil {
		return ActiveStateV2{}, indexErr
	}
	if indexes.Digests() != view.IndexDigests ||
		indexes.CombinedDigest() != view.CombinedIndexDigest {
		return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-index-digest")
	}
	if indexes.scope() != ArchiveIndexScopeRetainedGlobal {
		return ActiveStateV2{}, classified(ClassificationDomain, "active-v2-index-scope")
	}
	if view.HotSetDigests != (HotSetDigests{}) && !validF4BHotSetDigests(view.HotSetDigests) {
		return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-hot-set-digests")
	}
	seedCount, seedDigest, seedErr := effectSeedForIndexes(indexes)
	if seedErr != nil || seedCount != view.VisibleV1EffectSeedCount ||
		seedDigest != view.VisibleV1EffectSeedDigest {
		return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-effect-seed")
	}
	descriptorCounts := ArchiveCounts{}
	for _, descriptor := range view.Descriptors {
		descriptorCounts, err = descriptorCounts.add(descriptor.View().Counts)
		if err != nil {
			return ActiveStateV2{}, err
		}
	}
	if descriptorCounts != view.ArchivedCounts {
		return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-descriptor-counts")
	}
	total, err := view.HotCounts.add(view.ArchivedCounts)
	if err != nil || total != view.TotalCounts {
		return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-counts")
	}
	if indexes.counts() != view.TotalCounts {
		return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-global-index-counts")
	}
	if indexes.countsAt(RecordLocationHot) != view.HotCounts ||
		indexes.countsAt(RecordLocationArchive) != view.ArchivedCounts {
		return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-location-counts")
	}
	if !view.CurrentCheckpoint.valid(false) || !view.PreviousCheckpoint.valid(true) {
		return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-checkpoint-reference")
	}
	switch view.CurrentCheckpoint.Kind {
	case ArchiveCheckpointMigrationGenesis:
		if view.PreviousCheckpoint.Kind != ArchiveCheckpointNone || len(view.Descriptors) != 0 ||
			view.ArchivedCounts != (ArchiveCounts{}) || view.ArchiveGeneration != 1 {
			return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-genesis-world")
		}
	case ArchiveCheckpointActivation:
		if view.PreviousCheckpoint.Kind == ArchiveCheckpointNone || len(view.Descriptors) == 0 {
			return ActiveStateV2{}, classified(ClassificationBinding, "active-v2-activation-world")
		}
	default:
		return ActiveStateV2{}, classified(ClassificationDomain, "active-v2-checkpoint-kind")
	}
	cloned := view
	cloned.Descriptors = cloneSlice(view.Descriptors)
	cloned.Indexes = indexes
	return ActiveStateV2{view: cloned}, nil
}

func (state ActiveStateV2) View() ActiveStateV2View {
	view := state.view
	view.Descriptors = cloneSlice(state.view.Descriptors)
	return view
}

func cloneSlice[T any](values []T) []T {
	if values == nil {
		return nil
	}
	result := make([]T, len(values))
	copy(result, values)
	return result
}
