package archivestate

import (
	"bytes"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// ArchiveLocation identifies a record inside a referenced immutable segment.
// It is only the archive arm of a RecordLocation and is never authority without
// the corresponding verified record and digest.
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

func (location ArchiveLocation) isZero() bool { return location == (ArchiveLocation{}) }

// RecordLocationKind is the closed hot/archive discriminant for a retained
// record projection. Hot ordinals index the canonical sorted hot collection;
// archive ordinals index a referenced immutable segment.
type RecordLocationKind string

const (
	RecordLocationHot     RecordLocationKind = "hot"
	RecordLocationArchive RecordLocationKind = "archive"
)

// RecordLocationView keeps the discriminant, record domain, and only the
// coordinates valid for that discriminant. RecordKind prevents a registration
// location from being reused for an approval, attempt, or lifecycle projection.
type RecordLocationView struct {
	Kind             RecordLocationKind
	RecordKind       RecordKind
	HotRecordOrdinal v0candidate.PositiveUInt53
	Archive          ArchiveLocation
}

// RecordLocation is an immutable validated hot/archive location projection.
type RecordLocation struct{ view RecordLocationView }

// NewHotRecordLocation constructs a location in one canonical hot record set.
func NewHotRecordLocation(kind RecordKind, ordinal v0candidate.PositiveUInt53) (RecordLocation, error) {
	return newRecordLocation(RecordLocationView{Kind: RecordLocationHot, RecordKind: kind, HotRecordOrdinal: ordinal})
}

// NewArchiveRecordLocation constructs a location in one referenced segment.
func NewArchiveRecordLocation(kind RecordKind, location ArchiveLocation) (RecordLocation, error) {
	return newRecordLocation(RecordLocationView{Kind: RecordLocationArchive, RecordKind: kind, Archive: location})
}

func newRecordLocation(view RecordLocationView) (RecordLocation, error) {
	if !validRecordKind(view.RecordKind) {
		return RecordLocation{}, classified(ClassificationDomain, "record-location-record-kind")
	}
	switch view.Kind {
	case RecordLocationHot:
		if !validPositive(uint64(view.HotRecordOrdinal)) || !view.Archive.isZero() {
			return RecordLocation{}, classified(ClassificationBinding, "record-location-hot")
		}
	case RecordLocationArchive:
		if view.HotRecordOrdinal != 0 || !view.Archive.valid() {
			return RecordLocation{}, classified(ClassificationBinding, "record-location-archive")
		}
	default:
		return RecordLocation{}, classified(ClassificationDomain, "record-location-kind")
	}
	return RecordLocation{view: view}, nil
}

func (location RecordLocation) View() RecordLocationView { return location.view }

func (location RecordLocation) valid(kind RecordKind) bool {
	validated, err := newRecordLocation(location.view)
	return err == nil && validated.view.RecordKind == kind
}

type RegistrationIndexEntry struct {
	RegistrationID       v0candidate.RegistrationID
	RegistrationSequence v0candidate.PositiveUInt53
	PlanDigest           v0candidate.ExecutionPlanDigest
	ExpiresAt            v0candidate.UInt53
	Location             RecordLocation
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
	Location              RecordLocation
	FullRecordDigest      ArchiveRecordDigest
}

// AttemptLifecyclePresence is the closed discriminant for whether one retained
// immutable attempt has a colocated lifecycle record. Absence is an ordinary
// recoverable hot-state fact: approval consumption and attempt creation commit
// before lifecycle establishment.
type AttemptLifecyclePresence string

const (
	AttemptLifecycleAbsent  AttemptLifecyclePresence = "absent"
	AttemptLifecyclePresent AttemptLifecyclePresence = "present"
)

// AttemptLifecycleView is the lifecycle arm of one retained attempt entry.
// The absent arm carries no invented state, location, or digest. The present
// arm binds the exact lifecycle disposition and its independently typed record
// anchor. Attempt and lifecycle records remain separate full records.
type AttemptLifecycleView struct {
	Presence         AttemptLifecyclePresence
	State            lifecyclestate.LifecycleState
	Location         RecordLocation
	FullRecordDigest ArchiveRecordDigest
}

// AttemptLifecycle is an immutable validated lifecycle-presence union.
type AttemptLifecycle struct{ view AttemptLifecycleView }

// NoAttemptLifecycle represents a committed attempt before lifecycle
// establishment. It is not a lifecycle state and grants no cleanup authority.
func NoAttemptLifecycle() AttemptLifecycle {
	return AttemptLifecycle{view: AttemptLifecycleView{Presence: AttemptLifecycleAbsent}}
}

// NewPresentAttemptLifecycle binds an existing lifecycle record to its exact
// state, typed record location, and full-record digest.
func NewPresentAttemptLifecycle(
	state lifecyclestate.LifecycleState,
	location RecordLocation,
	digest ArchiveRecordDigest,
) (AttemptLifecycle, error) {
	view := AttemptLifecycleView{
		Presence: AttemptLifecyclePresent, State: state, Location: location, FullRecordDigest: digest,
	}
	if !validLifecycleState(state) || !location.valid(RecordLifecycle) || zeroDigest(digest) {
		return AttemptLifecycle{}, classified(ClassificationBinding, "attempt-lifecycle-present")
	}
	return AttemptLifecycle{view: view}, nil
}

func (lifecycle AttemptLifecycle) View() AttemptLifecycleView { return lifecycle.view }

func (lifecycle AttemptLifecycle) valid() bool {
	switch lifecycle.view.Presence {
	case AttemptLifecycleAbsent:
		return lifecycle.view.State == "" && lifecycle.view.Location == (RecordLocation{}) &&
			zeroDigest(lifecycle.view.FullRecordDigest)
	case AttemptLifecyclePresent:
		return validLifecycleState(lifecycle.view.State) &&
			lifecycle.view.Location.valid(RecordLifecycle) && !zeroDigest(lifecycle.view.FullRecordDigest)
	default:
		return false
	}
}

func (lifecycle AttemptLifecycle) present() bool {
	return lifecycle.view.Presence == AttemptLifecyclePresent
}

type AttemptIndexEntry struct {
	AttemptID        approvalattempt.AttemptID
	ApprovalID       approvalattempt.ApprovalID
	RegistrationID   v0candidate.RegistrationID
	CreatedAt        v0candidate.UInt53
	Lifecycle        AttemptLifecycle
	Location         RecordLocation
	FullRecordDigest ArchiveRecordDigest
}

type NonceIndexEntry struct {
	AttemptNonce  approvalattempt.AttemptNonce
	PayloadDigest approvalattempt.ApprovalPayloadDigest
	ApprovalID    approvalattempt.ApprovalID
	Location      RecordLocation
}

type EffectIndexEntry struct {
	EffectID                   lifecyclestate.EffectID
	AttemptID                  approvalattempt.AttemptID
	OperationSequence          v0candidate.PositiveUInt53
	Operation                  lifecyclestate.Operation
	IssuanceSnapshotGeneration v0candidate.PositiveUInt53
	VisibleV1Seed              bool
	AttemptLocation            RecordLocation
}

type InstanceIndexEntry struct {
	InstanceDigest lifecyclestate.BackendInstanceDigest
	AttemptID      approvalattempt.AttemptID
	Location       RecordLocation
}

type ApprovalReplayIndexEntry struct {
	PayloadDigest         approvalattempt.ApprovalPayloadDigest
	ExactPayloadDigest    ArchiveRecordDigest
	AuthorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity
	ApprovalID            approvalattempt.ApprovalID
	State                 approvalattempt.ApprovalState
	Location              RecordLocation
}

type AttemptReplayIndexEntry struct {
	RegistrationID v0candidate.RegistrationID
	ApprovalID     approvalattempt.ApprovalID
	AttemptID      approvalattempt.AttemptID
	State          approvalattempt.AttemptState
	Location       RecordLocation
}

// ArchiveIndexScope separates the complete retained global projection in the
// active snapshot from one segment's derived archive-only projection.
type ArchiveIndexScope string

const (
	ArchiveIndexScopeRetainedGlobal ArchiveIndexScope = "retained-global"
	ArchiveIndexScopeSegmentDerived ArchiveIndexScope = "segment-derived"
)

// ArchiveIndexesView is a closed sorted index projection. Retained-global
// contains every hot and archived identity; segment-derived contains only one
// segment's archive-derived entries. Every slice is non-nil, even when empty.
type ArchiveIndexesView struct {
	Scope          ArchiveIndexScope
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

// ArchiveIndexes is an immutable validated, scope-separated projection.
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

func emptyArchiveIndexes(scope ArchiveIndexScope) ArchiveIndexes {
	view := ArchiveIndexesView{
		Scope:          scope,
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

// EmptyRetainedIndexes returns the empty complete hot-plus-archive projection.
func EmptyRetainedIndexes() ArchiveIndexes {
	return emptyArchiveIndexes(ArchiveIndexScopeRetainedGlobal)
}

// EmptySegmentDerivedIndexes returns the empty archive-only segment projection.
func EmptySegmentDerivedIndexes() ArchiveIndexes {
	return emptyArchiveIndexes(ArchiveIndexScopeSegmentDerived)
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
	counts := ArchiveCounts{
		Registrations: uint64(len(view.Registrations)), Approvals: uint64(len(view.Approvals)),
		Attempts: uint64(len(view.Attempts)),
		Nonces:   uint64(len(view.Nonces)), Effects: uint64(len(view.Effects)),
		Instances: uint64(len(view.Instances)), ApprovalReplay: uint64(len(view.ApprovalReplay)),
		AttemptReplay: uint64(len(view.AttemptReplay)),
	}
	for _, entry := range view.Attempts {
		if entry.Lifecycle.present() {
			counts.Lifecycles++
		}
	}
	return counts
}

func (indexes ArchiveIndexes) scope() ArchiveIndexScope { return indexes.view.Scope }

func (indexes ArchiveIndexes) allLocationsAre(kind RecordLocationKind) bool {
	view := indexes.view
	for _, entry := range view.Registrations {
		if entry.Location.View().Kind != kind {
			return false
		}
	}
	for _, entry := range view.Approvals {
		if entry.Location.View().Kind != kind {
			return false
		}
	}
	for _, entry := range view.Attempts {
		if entry.Location.View().Kind != kind {
			return false
		}
		if entry.Lifecycle.present() && entry.Lifecycle.View().Location.View().Kind != kind {
			return false
		}
	}
	for _, entry := range view.Nonces {
		if entry.Location.View().Kind != kind {
			return false
		}
	}
	for _, entry := range view.Effects {
		if entry.AttemptLocation.View().Kind != kind {
			return false
		}
	}
	for _, entry := range view.Instances {
		if entry.Location.View().Kind != kind {
			return false
		}
	}
	for _, entry := range view.ApprovalReplay {
		if entry.Location.View().Kind != kind {
			return false
		}
	}
	for _, entry := range view.AttemptReplay {
		if entry.Location.View().Kind != kind {
			return false
		}
	}
	return true
}

func (indexes ArchiveIndexes) countsAt(kind RecordLocationKind) ArchiveCounts {
	view := indexes.view
	counts := ArchiveCounts{}
	for _, entry := range view.Registrations {
		if entry.Location.View().Kind == kind {
			counts.Registrations++
		}
	}
	for _, entry := range view.Approvals {
		if entry.Location.View().Kind == kind {
			counts.Approvals++
		}
	}
	for _, entry := range view.Attempts {
		if entry.Location.View().Kind == kind {
			counts.Attempts++
		}
		if entry.Lifecycle.present() && entry.Lifecycle.View().Location.View().Kind == kind {
			counts.Lifecycles++
		}
	}
	for _, entry := range view.Nonces {
		if entry.Location.View().Kind == kind {
			counts.Nonces++
		}
	}
	for _, entry := range view.Effects {
		if entry.AttemptLocation.View().Kind == kind {
			counts.Effects++
		}
	}
	for _, entry := range view.Instances {
		if entry.Location.View().Kind == kind {
			counts.Instances++
		}
	}
	for _, entry := range view.ApprovalReplay {
		if entry.Location.View().Kind == kind {
			counts.ApprovalReplay++
		}
	}
	for _, entry := range view.AttemptReplay {
		if entry.Location.View().Kind == kind {
			counts.AttemptReplay++
		}
	}
	return counts
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
	if view.Scope != ArchiveIndexScopeRetainedGlobal && view.Scope != ArchiveIndexScopeSegmentDerived {
		return classified(ClassificationDomain, "archive-index-scope")
	}
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
			uint64(entry.ExpiresAt) > v0candidate.MaxSafeInteger || !entry.Location.valid(RecordRegistration) ||
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
			!entry.Location.valid(RecordApproval) || zeroDigest(entry.FullRecordDigest) {
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
			!entry.Lifecycle.valid() || !entry.Location.valid(RecordAttempt) ||
			zeroDigest(entry.FullRecordDigest) {
			return nil, classified(ClassificationBinding, "attempt-index-entry")
		}
		attemptLocation := entry.Location.View()
		if !entry.Lifecycle.present() {
			if attemptLocation.Kind != RecordLocationHot {
				return nil, classified(ClassificationBinding, "attempt-index-absent-lifecycle-location")
			}
			return append([]byte(nil), entry.AttemptID[:]...), nil
		}
		lifecycleLocation := entry.Lifecycle.View().Location.View()
		if lifecycleLocation.Kind != attemptLocation.Kind ||
			(lifecycleLocation.Kind == RecordLocationArchive &&
				(lifecycleLocation.Archive.SegmentOrdinal != attemptLocation.Archive.SegmentOrdinal ||
					lifecycleLocation.Archive.CohortOrdinal != attemptLocation.Archive.CohortOrdinal)) {
			return nil, classified(ClassificationBinding, "attempt-index-lifecycle-location")
		}
		return append([]byte(nil), entry.AttemptID[:]...), nil
	}, "attempt-index-order")
}

func validateNonceIndex(entries []NonceIndexEntry) error {
	return validateStrictEntries(entries, func(entry NonceIndexEntry) ([]byte, error) {
		if entry.AttemptNonce == (approvalattempt.AttemptNonce{}) ||
			zeroDigest(entry.PayloadDigest) || entry.ApprovalID == (approvalattempt.ApprovalID{}) ||
			!entry.Location.valid(RecordApproval) {
			return nil, classified(ClassificationBinding, "nonce-index-entry")
		}
		return append([]byte(nil), entry.AttemptNonce[:]...), nil
	}, "nonce-index-order")
}

func validateEffectIndex(entries []EffectIndexEntry) error {
	return validateStrictEntries(entries, func(entry EffectIndexEntry) ([]byte, error) {
		if entry.EffectID.IsZero() || entry.AttemptID == (approvalattempt.AttemptID{}) ||
			!validPositive(uint64(entry.OperationSequence)) || !validLifecycleOperation(entry.Operation) ||
			!validPositive(uint64(entry.IssuanceSnapshotGeneration)) ||
			!entry.AttemptLocation.valid(RecordAttempt) {
			return nil, classified(ClassificationBinding, "effect-index-entry")
		}
		return append([]byte(nil), entry.EffectID[:]...), nil
	}, "effect-index-order")
}

func validateInstanceIndex(entries []InstanceIndexEntry) error {
	return validateStrictEntries(entries, func(entry InstanceIndexEntry) ([]byte, error) {
		if zeroDigest(entry.InstanceDigest) || entry.AttemptID == (approvalattempt.AttemptID{}) ||
			!entry.Location.valid(RecordLifecycle) {
			return nil, classified(ClassificationBinding, "instance-index-entry")
		}
		return append([]byte(nil), entry.InstanceDigest[:]...), nil
	}, "instance-index-order")
}

func validateApprovalReplayIndex(entries []ApprovalReplayIndexEntry) error {
	return validateStrictEntries(entries, func(entry ApprovalReplayIndexEntry) ([]byte, error) {
		if zeroDigest(entry.PayloadDigest) || zeroDigest(entry.ExactPayloadDigest) || zeroDigest(entry.AuthorizationIdentity) ||
			entry.ApprovalID == (approvalattempt.ApprovalID{}) || !validApprovalState(entry.State) ||
			!entry.Location.valid(RecordApproval) {
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
			entry.AttemptID == (approvalattempt.AttemptID{}) || entry.State != approvalattempt.AttemptCreated ||
			!entry.Location.valid(RecordAttempt) {
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
		Registrations:  digestIndex(view.Scope, "registration", view.Registrations, encodeRegistrationIndex),
		Approvals:      digestIndex(view.Scope, "approval", view.Approvals, encodeApprovalIndex),
		Attempts:       digestIndex(view.Scope, "attempt", view.Attempts, encodeAttemptIndex),
		Nonces:         digestIndex(view.Scope, "nonce", view.Nonces, encodeNonceIndex),
		Effects:        digestIndex(view.Scope, "effect", view.Effects, encodeEffectIndex),
		Instances:      digestIndex(view.Scope, "instance", view.Instances, encodeInstanceIndex),
		ApprovalReplay: digestIndex(view.Scope, "approval-replay", view.ApprovalReplay, encodeApprovalReplayIndex),
		AttemptReplay:  digestIndex(view.Scope, "attempt-replay", view.AttemptReplay, encodeAttemptReplayIndex),
	}
}

func digestIndex[T any](scope ArchiveIndexScope, name string, entries []T, encode func(*digestEncoder, T)) ArchiveIndexDigest {
	encoder := newDigestEncoder("capsule.supervisor.archive-index.v0")
	encoder.text(string(scope))
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
	ids := indexes.visibleV1EffectSeed()
	digest, err := DigestVisibleV1EffectSeed(ids)
	return uint64(len(ids)), digest, err
}

func (indexes ArchiveIndexes) visibleV1EffectSeed() []lifecyclestate.EffectID {
	ids := make([]lifecyclestate.EffectID, 0)
	for _, entry := range indexes.view.Effects {
		if entry.VisibleV1Seed {
			ids = append(ids, entry.EffectID)
		}
	}
	return ids
}

func encodeArchiveLocation(encoder *digestEncoder, location ArchiveLocation) {
	encoder.uint64(uint64(location.SegmentOrdinal))
	encoder.uint64(uint64(location.CohortOrdinal))
	encoder.uint64(uint64(location.RecordOrdinal))
}

func encodeRecordLocation(encoder *digestEncoder, location RecordLocation) {
	view := location.View()
	encoder.text(string(view.Kind))
	encoder.text(string(view.RecordKind))
	if view.Kind == RecordLocationHot {
		encoder.uint64(uint64(view.HotRecordOrdinal))
		return
	}
	encodeArchiveLocation(encoder, view.Archive)
}

func encodeRegistrationIndex(encoder *digestEncoder, entry RegistrationIndexEntry) {
	encoder.bytes(entry.RegistrationID[:])
	encoder.uint64(uint64(entry.RegistrationSequence))
	encoder.bytes(entry.PlanDigest[:])
	encoder.uint64(uint64(entry.ExpiresAt))
	encodeRecordLocation(encoder, entry.Location)
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
	encodeRecordLocation(encoder, entry.Location)
	encoder.bytes(entry.FullRecordDigest[:])
}

func encodeAttemptIndex(encoder *digestEncoder, entry AttemptIndexEntry) {
	encoder.bytes(entry.AttemptID[:])
	encoder.bytes(entry.ApprovalID[:])
	encoder.bytes(entry.RegistrationID[:])
	encoder.uint64(uint64(entry.CreatedAt))
	lifecycle := entry.Lifecycle.View()
	encoder.text(string(lifecycle.Presence))
	if lifecycle.Presence == AttemptLifecyclePresent {
		encoder.text(string(lifecycle.State))
		encodeRecordLocation(encoder, lifecycle.Location)
		encoder.bytes(lifecycle.FullRecordDigest[:])
	}
	encodeRecordLocation(encoder, entry.Location)
	encoder.bytes(entry.FullRecordDigest[:])
}

func encodeNonceIndex(encoder *digestEncoder, entry NonceIndexEntry) {
	encoder.bytes(entry.AttemptNonce[:])
	encoder.bytes(entry.PayloadDigest[:])
	encoder.bytes(entry.ApprovalID[:])
	encodeRecordLocation(encoder, entry.Location)
}

func encodeEffectIndex(encoder *digestEncoder, entry EffectIndexEntry) {
	encoder.bytes(entry.EffectID[:])
	encoder.bytes(entry.AttemptID[:])
	encoder.uint64(uint64(entry.OperationSequence))
	encoder.text(string(entry.Operation))
	encoder.uint64(uint64(entry.IssuanceSnapshotGeneration))
	encoder.boolean(entry.VisibleV1Seed)
	encodeRecordLocation(encoder, entry.AttemptLocation)
}

func encodeInstanceIndex(encoder *digestEncoder, entry InstanceIndexEntry) {
	encoder.bytes(entry.InstanceDigest[:])
	encoder.bytes(entry.AttemptID[:])
	encodeRecordLocation(encoder, entry.Location)
}

func encodeApprovalReplayIndex(encoder *digestEncoder, entry ApprovalReplayIndexEntry) {
	encoder.bytes(entry.PayloadDigest[:])
	encoder.bytes(entry.ExactPayloadDigest[:])
	encoder.bytes(entry.AuthorizationIdentity[:])
	encoder.bytes(entry.ApprovalID[:])
	encoder.text(string(entry.State))
	encodeRecordLocation(encoder, entry.Location)
}

func encodeAttemptReplayIndex(encoder *digestEncoder, entry AttemptReplayIndexEntry) {
	encoder.bytes(entry.RegistrationID[:])
	encoder.bytes(entry.ApprovalID[:])
	encoder.bytes(entry.AttemptID[:])
	encoder.text(string(entry.State))
	encodeRecordLocation(encoder, entry.Location)
}
