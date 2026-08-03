package registrationstate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// StoreFault identifies the fixed store's confirmed-abort and indeterminate
// durability boundaries.
type StoreFault string

const (
	SupervisorStoreFormatV0 = uint64(0)
	SupervisorStoreFormatV1 = uint64(1)
	MaxSupervisorStateBytes = int64(16 << 20)

	FaultTimeHighWaterWrite                  StoreFault = "time-high-water-write"
	FaultTimeHighWaterIndeterminatePreState  StoreFault = "time-high-water-indeterminate-pre-state"
	FaultTimeHighWaterCommitIndeterminate    StoreFault = "time-high-water-commit-indeterminate"
	FaultRegistrationCommitAbort             StoreFault = "registration-commit-confirmed-abort"
	FaultApprovalCommitAbort                 StoreFault = "approval-commit-confirmed-abort"
	FaultApprovalCommitIndeterminatePreState StoreFault = "approval-commit-indeterminate-pre-state"
	FaultApprovalCommitIndeterminate         StoreFault = "approval-commit-indeterminate"
	FaultAttemptCommitAbort                  StoreFault = "attempt-commit-confirmed-abort"
	FaultAttemptCommitIndeterminatePreState  StoreFault = "attempt-commit-indeterminate-pre-state"
	FaultAttemptCommitIndeterminate          StoreFault = "attempt-commit-indeterminate"
)

// ErrCommitOutcomeIndeterminate means the atomic rename completed but the
// directory durability barrier failed. It is intentionally not classified as
// a confirmed-abort LOCAL_FAILURE; callers must recover the store before
// deciding whether authority committed.
var ErrCommitOutcomeIndeterminate = errors.New("fixed supervisor-state commit outcome is indeterminate")

// ErrStoreRepairRequired reports that startup validation refused a corrupt or
// cross-linked snapshot without rewriting or deleting the retained evidence.
var ErrStoreRepairRequired = errors.New("fixed supervisor state requires repair")

var errNoStoreMutation = errors.New("fixed store transaction has no mutation")

// FixedFileStore is the deliberately bounded, no-eviction first Supervisor store. It
// writes a complete versioned snapshot with fsync/rename/fsync ordering and is
// fault-injectable for conformance and fake-backend development. It is not a
// reviewed production database: it has no archive, compaction, encryption, or
// multi-process coordination, and no consumer may activate it.
type FixedFileStore struct {
	mu               sync.Mutex
	path             string
	state            installationState
	faults           map[StoreFault]error
	recoveryRequired bool
}

type diskEnvelope struct {
	StoreFormatVersion uint64            `json:"storeFormatVersion"`
	State              installationState `json:"state"`
}

func NewFixedFileStore(path string, initial InitialState) (*FixedFileStore, error) {
	if path == "" {
		return nil, errors.New("fixed registration store path is required")
	}
	store := &FixedFileStore{path: path, faults: make(map[StoreFault]error)}
	info, err := os.Lstat(path)
	switch {
	case err == nil:
		if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			return nil, errors.New("fixed registration store must be a regular file")
		}
		state, loadErr := loadState(path)
		if loadErr != nil {
			return nil, fmt.Errorf("%w: %v", ErrStoreRepairRequired, loadErr)
		}
		store.state = state
		return store, nil
	case !errors.Is(err, os.ErrNotExist):
		return nil, errors.New("inspect fixed registration store")
	}

	state := installationState{
		InitialState:          initial,
		RegistrationSetDigest: emptyRegistrationSetDigest(),
		Registrations:         make([]registrationEntry, 0),
		ApprovalSetDigest:     emptyApprovalSetDigest(),
		Approvals:             make([]approvalattempt.ApprovalRecord, 0),
		AttemptSetDigest:      emptyAttemptSetDigest(),
		Attempts:              make([]approvalattempt.ExecutionAttempt, 0),
	}
	if state.TrustPhase == "" {
		state.TrustPhase = TrustStable
	}
	if state.LastRegistrationSequence != 0 {
		return nil, errors.New("new fixed registration store must begin before sequence one")
	}
	if err := validateState(state); err != nil {
		return nil, err
	}
	if err := persistState(path, state); err != nil {
		return nil, err
	}
	store.state = state
	return store, nil
}

// InjectFailure installs a one-shot failure at an exact store boundary.
func (store *FixedFileStore) InjectFailure(point StoreFault, err error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if err == nil {
		err = errors.New("injected fixed-store failure")
	}
	store.faults[point] = err
}

func (store *FixedFileStore) snapshot(ctx context.Context) (installationState, error) {
	if err := ctx.Err(); err != nil {
		return installationState{}, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	return cloneState(store.state), nil
}

func (store *FixedFileStore) recoveryFenced() bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.recoveryRequired
}

func (store *FixedFileStore) persistTimeHighWater(
	ctx context.Context,
	observation v0candidate.UInt53,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.recoveryRequired {
		return ErrCommitOutcomeIndeterminate
	}
	if err := store.takeFault(FaultTimeHighWaterWrite); err != nil {
		return err
	}
	candidate := cloneState(store.state)
	if observation > candidate.TimeHighWaterUnixSeconds {
		candidate.TimeHighWaterUnixSeconds = observation
	}
	if err := validateState(candidate); err != nil {
		return err
	}
	if err := store.takeFault(FaultTimeHighWaterIndeterminatePreState); err != nil {
		store.recoveryRequired = true
		return fmt.Errorf("%w: injected before time commit", ErrCommitOutcomeIndeterminate)
	}
	indeterminate := store.takeFault(FaultTimeHighWaterCommitIndeterminate)
	if err := persistStateWithBoundary(store.path, candidate, indeterminate); err != nil {
		if errors.Is(err, ErrCommitOutcomeIndeterminate) {
			store.state = candidate
			store.recoveryRequired = true
		}
		return err
	}
	store.state = candidate
	return nil
}

func (store *FixedFileStore) update(
	ctx context.Context,
	update func(*installationState) error,
) error {
	if update == nil {
		return errors.New("registration state update is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return store.commit(ctx, FaultRegistrationCommitAbort, "", "", update)
}

func (store *FixedFileStore) commitApproval(
	ctx context.Context,
	update func(*installationState) error,
) error {
	return store.commit(ctx, FaultApprovalCommitAbort, FaultApprovalCommitIndeterminatePreState, FaultApprovalCommitIndeterminate, update)
}

func (store *FixedFileStore) commitAttempt(
	ctx context.Context,
	update func(*installationState) error,
) error {
	return store.commit(ctx, FaultAttemptCommitAbort, FaultAttemptCommitIndeterminatePreState, FaultAttemptCommitIndeterminate, update)
}

func (store *FixedFileStore) commit(
	ctx context.Context,
	abortPoint StoreFault,
	preIndeterminatePoint StoreFault,
	indeterminatePoint StoreFault,
	update func(*installationState) error,
) error {
	if update == nil {
		return errors.New("registration state update is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.recoveryRequired {
		return ErrCommitOutcomeIndeterminate
	}
	candidate := cloneState(store.state)
	if err := update(&candidate); err != nil {
		if errors.Is(err, errNoStoreMutation) {
			return nil
		}
		return err
	}
	if err := validateState(candidate); err != nil {
		return err
	}
	if err := store.takeFault(abortPoint); err != nil {
		return err
	}
	if preIndeterminatePoint != "" {
		if err := store.takeFault(preIndeterminatePoint); err != nil {
			store.recoveryRequired = true
			return fmt.Errorf("%w: injected before state commit", ErrCommitOutcomeIndeterminate)
		}
	}
	var indeterminate error
	if indeterminatePoint != "" {
		indeterminate = store.takeFault(indeterminatePoint)
	}
	if err := persistStateWithBoundary(store.path, candidate, indeterminate); err != nil {
		if errors.Is(err, ErrCommitOutcomeIndeterminate) {
			store.state = candidate
			store.recoveryRequired = true
		}
		return err
	}
	store.state = candidate
	return nil
}

func (store *FixedFileStore) takeFault(point StoreFault) error {
	err := store.faults[point]
	delete(store.faults, point)
	return err
}

func loadState(path string) (installationState, error) {
	data, err := readBoundedStoreFile(path)
	if err != nil {
		return installationState{}, err
	}
	var header struct {
		StoreFormatVersion *uint64 `json:"storeFormatVersion"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return installationState{}, errors.New("decode fixed registration store")
	}
	if header.StoreFormatVersion == nil || *header.StoreFormatVersion != SupervisorStoreFormatV0 {
		return installationState{}, errors.New("unsupported fixed registration store version")
	}
	var envelope diskEnvelope
	if err := decodeOneClosedJSON(data, &envelope); err != nil {
		return installationState{}, errors.New("decode fixed registration store")
	}
	if err := validateState(envelope.State); err != nil {
		return installationState{}, err
	}
	return cloneState(envelope.State), nil
}

func persistState(path string, state installationState) error {
	return persistStateWithBoundary(path, state, nil)
}

func persistStateWithBoundary(path string, state installationState, afterRename error) error {
	parent := filepath.Dir(path)
	temporary, err := os.CreateTemp(parent, ".capsule-supervisor-state-*.tmp")
	if err != nil {
		return errors.New("create fixed registration transaction")
	}
	temporaryPath := temporary.Name()
	committed := false
	defer func() {
		_ = temporary.Close()
		if !committed {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(0o600); err != nil {
		return errors.New("protect fixed registration transaction")
	}
	encoder := json.NewEncoder(temporary)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(diskEnvelope{StoreFormatVersion: SupervisorStoreFormatV0, State: state}); err != nil {
		return errors.New("encode fixed registration transaction")
	}
	if err := temporary.Sync(); err != nil {
		return errors.New("sync fixed registration transaction")
	}
	if err := temporary.Close(); err != nil {
		return errors.New("close fixed registration transaction")
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return errors.New("commit fixed registration transaction")
	}
	committed = true
	if afterRename != nil {
		return fmt.Errorf("%w: injected after rename", ErrCommitOutcomeIndeterminate)
	}
	directory, err := os.Open(parent)
	if err != nil {
		return fmt.Errorf("%w: open store directory", ErrCommitOutcomeIndeterminate)
	}
	syncErr := directory.Sync()
	closeErr := directory.Close()
	if syncErr != nil || closeErr != nil {
		return fmt.Errorf("%w: sync store directory", ErrCommitOutcomeIndeterminate)
	}
	return nil
}

func validateState(state installationState) error {
	if state.InstallationID == (v0candidate.InstallationID{}) ||
		state.SupervisorID == (v0candidate.SupervisorID{}) {
		return errors.New("fixed registration state lost installation identity")
	}
	if uint64(state.EpochSequence) > v0candidate.MaxSafeInteger ||
		uint64(state.TimeHighWaterUnixSeconds) > v0candidate.MaxSafeInteger ||
		uint64(state.LastRegistrationSequence) > v0candidate.MaxSafeInteger {
		return errors.New("fixed registration state contains a value outside UInt53")
	}
	if state.TrustPhase != TrustStable && state.TrustPhase != TrustTransitionFenced &&
		state.TrustPhase != TrustRepairRequired {
		return errors.New("fixed registration state has an unknown trust phase")
	}
	if state.TrustPhase == TrustRepairRequired && state.TrustReason == "" {
		return errors.New("repair-required registration state lacks a reason")
	}
	if len(state.Registrations) > MaxRetainedRegistrations {
		return errors.New("fixed registration state exceeds retained capacity")
	}
	seen := make(map[v0candidate.RegistrationID]struct{}, len(state.Registrations))
	var previousSequence uint64
	for _, entry := range state.Registrations {
		if entry.Index.RegistrationID == (v0candidate.RegistrationID{}) ||
			entry.Index.InstallationID == (v0candidate.InstallationID{}) ||
			entry.Index.SupervisorID == (v0candidate.SupervisorID{}) ||
			entry.Index.RegistrationSequence == 0 ||
			uint64(entry.Index.RegistrationSequence) > v0candidate.MaxSafeInteger ||
			uint64(entry.Index.EpochSequence) > v0candidate.MaxSafeInteger ||
			uint64(entry.Index.ExpiresAt) > v0candidate.MaxSafeInteger ||
			uint64(entry.Index.RegistrationSequence) > uint64(state.LastRegistrationSequence) {
			return errors.New("fixed registration state has an invalid registration index")
		}
		if uint64(entry.Index.RegistrationSequence) <= previousSequence {
			return errors.New("fixed registration sequences are not strictly increasing")
		}
		previousSequence = uint64(entry.Index.RegistrationSequence)
		if _, exists := seen[entry.Index.RegistrationID]; exists {
			return errors.New("fixed registration state has a duplicate registration ID")
		}
		seen[entry.Index.RegistrationID] = struct{}{}
		if entry.Record.RegisteredAtUnixSeconds > state.TimeHighWaterUnixSeconds {
			return errors.New("fixed registration record is newer than the time high water")
		}
		if err := validateStoredRecord(entry); err != nil {
			return err
		}
	}
	return validateApprovalAttemptState(state)
}

func validateStoredRecord(entry registrationEntry) error {
	record := entry.Record
	if len(record.WireRegistrationBytes) == 0 ||
		len(record.WireRegistrationBytes) > v0candidate.PlanRegistrationMaxCBORBytes ||
		len(record.ExactPlanBytes) == 0 ||
		len(record.ExactPlanBytes) > v0candidate.ExecutionPlanMaxCBORBytes ||
		record.StorageFormatVersion != StorageFormatVersion ||
		record.RetentionState != RetentionStateRetained ||
		uint64(record.RegisteredAtUnixSeconds) > v0candidate.MaxSafeInteger ||
		record.RegisteredAtUnixSeconds >= entry.Index.ExpiresAt {
		return errors.New("fixed registration state has an invalid stored record")
	}
	digest := sha256.Sum256(record.ExactPlanBytes)
	if v0candidate.ExecutionPlanDigest(digest) != record.RecomputedPlanDigest {
		return errors.New("fixed registration stored plan digest mismatch")
	}
	decodedPlan, err := v0candidate.DecodeExecutionPlan(record.ExactPlanBytes, entry.PlanBindings)
	if err != nil || decodedPlan.Digest() != record.RecomputedPlanDigest {
		return errors.New("fixed registration plan validation failed")
	}
	planView := decodedPlan.View()
	if planView.InstallationID != entry.Index.InstallationID ||
		planView.EpochSequence != entry.Index.EpochSequence ||
		planView.EpochDigest != entry.Index.EpochDigest ||
		planView.ExpiresAt != entry.Index.ExpiresAt {
		return errors.New("fixed registration plan/index mismatch")
	}
	decoded, err := v0candidate.DecodePlanRegistration(
		record.WireRegistrationBytes,
		v0candidate.PlanRegistrationRoleBindings{
			RegistrationID: entry.Index.RegistrationID,
			PlanDigest:     record.RecomputedPlanDigest,
			InstallationID: entry.Index.InstallationID,
			EpochDigest:    entry.Index.EpochDigest,
			SupervisorID:   entry.Index.SupervisorID,
		},
	)
	if err != nil {
		return errors.New("fixed registration wire validation failed")
	}
	view := decoded.View()
	if view.RegistrationSequence != entry.Index.RegistrationSequence ||
		view.EpochSequence != entry.Index.EpochSequence ||
		view.ExpiresAt != entry.Index.ExpiresAt {
		return errors.New("fixed registration wire/index mismatch")
	}
	return nil
}

func cloneState(state installationState) installationState {
	clone := state
	clone.Registrations = make([]registrationEntry, len(state.Registrations))
	for index, entry := range state.Registrations {
		clone.Registrations[index] = entry
		clone.Registrations[index].PlanBindings = clonePlanBindings(entry.PlanBindings)
		clone.Registrations[index].Record = cloneStoredRegistration(entry.Record)
	}
	clone.Approvals = make([]approvalattempt.ApprovalRecord, len(state.Approvals))
	for index, record := range state.Approvals {
		clone.Approvals[index] = approvalattempt.CloneApprovalRecord(record)
	}
	clone.Attempts = make([]approvalattempt.ExecutionAttempt, len(state.Attempts))
	copy(clone.Attempts, state.Attempts)
	return clone
}

func snapshotProjection(state installationState) stateSnapshot {
	records := make([]StoredRegistration, len(state.Registrations))
	for index, entry := range state.Registrations {
		records[index] = cloneStoredRegistration(entry.Record)
	}
	return stateSnapshot{
		InitialState:               state.InitialState,
		StoredRegistrationCount:    len(state.Registrations),
		UnexpiredRegistrationCount: countUnexpired(state.Registrations, state.TimeHighWaterUnixSeconds),
		RegistrationSetDigest:      state.RegistrationSetDigest,
		Records:                    records,
	}
}

type storedRecordJSON struct {
	WireRegistrationHex        string `json:"wireRegistrationHex"`
	ExactPlanHex               string `json:"exactPlanHex"`
	RecomputedPlanDigestHex    string `json:"recomputedPlanDigestHex"`
	RegisteredAtUnixSeconds    uint64 `json:"registeredAtUnixSeconds"`
	StoredRecordFormatVersion  uint64 `json:"storageFormatVersion"`
	StoredRecordRetentionState string `json:"retentionState"`
}

func appendRegistrationSetDigest(previous [32]byte, record StoredRegistration) ([32]byte, error) {
	projection := storedRecordJSON{
		WireRegistrationHex:        hex.EncodeToString(record.WireRegistrationBytes),
		ExactPlanHex:               hex.EncodeToString(record.ExactPlanBytes),
		RecomputedPlanDigestHex:    hex.EncodeToString(record.RecomputedPlanDigest[:]),
		RegisteredAtUnixSeconds:    uint64(record.RegisteredAtUnixSeconds),
		StoredRecordFormatVersion:  record.StorageFormatVersion,
		StoredRecordRetentionState: record.RetentionState,
	}
	encoded, err := json.Marshal(projection)
	if err != nil {
		return [32]byte{}, fmt.Errorf("encode closed stored registration record: %w", err)
	}
	material := make([]byte, 0, 64+1+len(encoded))
	material = append(material, hex.EncodeToString(previous[:])...)
	material = append(material, ':')
	material = append(material, encoded...)
	return sha256.Sum256(material), nil
}

func recordsEqual(left, right StoredRegistration) bool {
	return bytes.Equal(left.WireRegistrationBytes, right.WireRegistrationBytes) &&
		bytes.Equal(left.ExactPlanBytes, right.ExactPlanBytes) &&
		left.RecomputedPlanDigest == right.RecomputedPlanDigest &&
		left.RegisteredAtUnixSeconds == right.RegisteredAtUnixSeconds &&
		left.StorageFormatVersion == right.StorageFormatVersion &&
		left.RetentionState == right.RetentionState
}
