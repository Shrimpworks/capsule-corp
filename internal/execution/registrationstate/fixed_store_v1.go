package registrationstate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/installationowner"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	// MaxRetainedLifecycleRecords and MaxActiveLifecycleRecords are inclusive,
	// no-eviction fixed-oracle limits. Only a durably destroyed record with
	// cleanup false after authoritative absence releases active capacity.
	MaxRetainedLifecycleRecords = 4_096
	// MaxActiveLifecycleRecords caps attempts whose cleanup work is not durably closed.
	MaxActiveLifecycleRecords = 256
)

// LifecycleSetDigest identifies the complete sorted v1 lifecycle collection.
// It is a deterministic integrity value, not an origin or authorization proof.
type LifecycleSetDigest [32]byte

// OfflineMigrationLock is an injected assertion that the installation's
// exclusive owner lock remains held. E2 deliberately does not claim a real
// platform lock implementation.
type OfflineMigrationLock interface {
	CheckOfflineMigrationLock(context.Context) error
}

// MigrationFault identifies the two E2 commit-outcome oracles.
type MigrationFault string

const (
	// FaultMigrationBeforeRename is a confirmed-abort test edge;
	// FaultMigrationAfterRename is indeterminate until reopen establishes the
	// complete v0 predecessor or v1 successor. Neither is a product control.
	FaultMigrationBeforeRename MigrationFault = "migration-before-rename"
	// FaultMigrationAfterRename requires reopen to distinguish complete v0/v1 worlds.
	FaultMigrationAfterRename MigrationFault = "migration-after-rename"
)

// MigrationFaultInjector is a local conformance port. Returning an error at
// the pre-rename point is a confirmed failure; returning one after rename is
// indeterminate until the file is reopened.
type MigrationFaultInjector interface {
	FailMigrationAt(MigrationFault) error
}

// V0ToV1MigrationOptions carries only the asserted offline owner-lock check
// and optional local fault injector. It contains no replacement records or
// migration data and cannot relax full source validation.
type V0ToV1MigrationOptions struct {
	Lock   OfflineMigrationLock
	Faults MigrationFaultInjector
}

var (
	// ErrMigrationLockRequired reports missing/lost exclusive ownership;
	// ErrMigrationOutcomeIndeterminate requires reopen before deciding which
	// complete format committed. Callers must not retry mutation in place.
	ErrMigrationLockRequired = errors.New("fixed supervisor-state migration requires the offline owner lock")
	// ErrMigrationOutcomeIndeterminate requires full reopen before later mutation.
	ErrMigrationOutcomeIndeterminate = errors.New("fixed supervisor-state migration outcome is indeterminate")
)

// FixedFileStoreV1 is the E3 transactional view of a validated v1 snapshot.
// G2 can bind it to the opaque Darwin installation owner, but it remains an
// unwired fixed-file checkpoint with no consumer, adapter invocation, installed
// protected-storage, or production platform claim.
type FixedFileStoreV1 struct {
	mu                  sync.Mutex
	path                string
	sourceFormatVersion uint64
	state               installationState
	lifecycleSetDigest  LifecycleSetDigest
	lifecycles          []lifecyclestate.Record
	effectIDs           LifecycleEffectIDSource
	ownerSessionID      lifecyclestate.OwnerSessionID
	faults              map[LifecycleStoreFault]error
	issuedPermits       map[approvalattempt.AttemptID]lifecyclestate.EffectPermit
	effectTombstones    []archivestate.EffectIndexEntry
	retainedEffects     map[lifecyclestate.EffectID]struct{}
	v2Commit            func(context.Context, LifecycleStoreFault, LifecycleStoreFault, LifecycleStoreFault, bool, *installationState, *[]lifecyclestate.Record, *[]archivestate.EffectIndexEntry) error
	recoveryRequired    bool
	repairRequired      bool
	owner               installationowner.InstallationOwner
	v2Owner             ArchiveOwner
	ownerRequired       bool
	ownerFenced         bool
	closed              bool
}

type fixedStoreV1Snapshot struct {
	SourceFormatVersion uint64
	State               installationState
	LifecycleSetDigest  LifecycleSetDigest
	Lifecycles          []lifecyclestate.Record
}

type diskEnvelopeV1 struct {
	StoreFormatVersion     *uint64               `json:"storeFormatVersion"`
	MigrationSourceVersion *uint64               `json:"migrationSourceVersion"`
	State                  installationState     `json:"state"`
	LifecycleSetDigest     LifecycleSetDigest    `json:"lifecycleSetDigest"`
	Lifecycles             []lifecycleRecordDisk `json:"lifecycles"`
}

type lifecycleRecordDisk struct {
	FormatVersion      uint64                       `json:"formatVersion"`
	RecordVersion      lifecyclestate.RecordVersion `json:"recordVersion"`
	SnapshotGeneration v0candidate.PositiveUInt53   `json:"snapshotGeneration"`

	Bindings               lifecycleBindingsDisk                 `json:"bindings"`
	ImmutableBindingDigest lifecyclestate.ImmutableBindingDigest `json:"immutableBindingDigest"`

	OperationSequence v0candidate.PositiveUInt53  `json:"operationSequence"`
	Operation         lifecyclestate.Operation    `json:"operation"`
	EffectID          lifecyclestate.EffectID     `json:"effectId"`
	EffectStatus      lifecyclestate.EffectStatus `json:"effectStatus"`

	Instance lifecyclestate.BackendInstanceIdentityView `json:"instance"`

	State                   lifecyclestate.LifecycleState        `json:"state"`
	CleanupRequired         bool                                 `json:"cleanupRequired"`
	LastConfirmedCheckpoint lifecyclestate.Checkpoint            `json:"lastConfirmedCheckpoint"`
	FirstFailure            lifecyclestate.FailureClassification `json:"firstFailure"`
	FailureOperation        lifecyclestate.Operation             `json:"failureOperation"`

	LastReconciliation         lifecyclestate.ReconciliationStatus `json:"lastReconciliation"`
	AutomaticRecoveryCount     v0candidate.UInt53                  `json:"automaticRecoveryCount"`
	NextRecoveryAt             lifecyclestate.OptionalUnixSeconds  `json:"nextRecoveryAt"`
	RecoveryFence              lifecyclestate.RecoveryFenceReason  `json:"recoveryFence"`
	AutomaticRecoveryExhausted bool                                `json:"automaticRecoveryExhausted"`

	OpenedAt         v0candidate.UInt53                 `json:"openedAt"`
	LastTransitionAt v0candidate.UInt53                 `json:"lastTransitionAt"`
	TerminalAt       lifecyclestate.OptionalUnixSeconds `json:"terminalAt"`
}

type lifecycleBindingsDisk struct {
	AttemptID             approvalattempt.AttemptID                        `json:"attemptId"`
	ApprovalID            approvalattempt.ApprovalID                       `json:"approvalId"`
	AttemptNonce          approvalattempt.AttemptNonce                     `json:"attemptNonce"`
	RegistrationID        v0candidate.RegistrationID                       `json:"registrationId"`
	RegistrationSequence  v0candidate.PositiveUInt53                       `json:"registrationSequence"`
	PlanDigest            v0candidate.ExecutionPlanDigest                  `json:"planDigest"`
	InstallationID        v0candidate.InstallationID                       `json:"installationId"`
	EpochSequence         v0candidate.UInt53                               `json:"epochSequence"`
	EpochDigest           v0candidate.TrustEpochDigest                     `json:"epochDigest"`
	SupervisorID          v0candidate.SupervisorID                         `json:"supervisorId"`
	ApprovalPurpose       string                                           `json:"approvalPurpose"`
	ApprovalAudience      string                                           `json:"approvalAudience"`
	ApprovalPayloadDigest approvalattempt.ApprovalPayloadDigest            `json:"approvalPayloadDigest"`
	AuthorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity `json:"authorizationIdentity"`
	AttemptCreatedAt      v0candidate.UInt53                               `json:"attemptCreatedAt"`

	AttemptStorageVersion      uint64 `json:"attemptStorageVersion"`
	ApprovalStorageVersion     uint64 `json:"approvalStorageVersion"`
	RegistrationStorageVersion uint64 `json:"registrationStorageVersion"`

	SourceManifestDigest            v0candidate.SourceManifestDigest             `json:"sourceManifestDigest"`
	InlineInputDigest               v0candidate.InlineInputDigest                `json:"inlineInputDigest"`
	RuntimeBundleManifestDigest     v0candidate.RuntimeBundleManifestDigest      `json:"runtimeBundleManifestDigest"`
	ProfileReviewAttestationDigests []v0candidate.ProfileReviewAttestationDigest `json:"profileReviewAttestationDigests"`
	ProfileRegistryEntryDigest      v0candidate.ProfileRegistryEntryDigest       `json:"profileRegistryEntryDigest"`
	BackendValidationRecordDigest   v0candidate.BackendValidationRecordDigest    `json:"backendValidationRecordDigest"`
	BackendConfigurationDigest      v0candidate.BackendConfigurationDigest       `json:"backendConfigurationDigest"`
	TrustSnapshotDigest             v0candidate.TrustSnapshotDigest              `json:"trustSnapshotDigest"`
	PolicyDecisionDigest            v0candidate.PolicyDecisionDigest             `json:"policyDecisionDigest"`

	Backend lifecyclestate.BackendBindingView `json:"backend"`
}

// OpenFixedFileStoreV1 opens only an existing v1 file. It never creates a
// missing store and never interprets v0 as an empty-lifecycle v1 snapshot.
func OpenFixedFileStoreV1(path string) (*FixedFileStoreV1, error) {
	return OpenFixedFileStoreV1WithOptions(path, FixedFileStoreV1Options{})
}

// OpenFixedFileStoreV1WithOptions opens an existing validated v1 snapshot and
// installs trusted local identifier sources. The options exist for retained,
// deterministic fault and collision tests; they authorize no external effect.
func OpenFixedFileStoreV1WithOptions(path string, options FixedFileStoreV1Options) (*FixedFileStoreV1, error) {
	return openFixedFileStoreV1WithOptions(path, options)
}

func openFixedFileStoreV1WithOptions(path string, options FixedFileStoreV1Options) (*FixedFileStoreV1, error) {
	if path == "" {
		return nil, errors.New("fixed supervisor v1 store path is required")
	}
	if err := validateExistingStoreFile(path); err != nil {
		return nil, err
	}
	loaded, err := loadV1State(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStoreRepairRequired, err)
	}
	resolved, err := resolveFixedFileStoreV1Options(options)
	if err != nil {
		return nil, err
	}
	return &FixedFileStoreV1{
		path:                path,
		sourceFormatVersion: loaded.SourceFormatVersion,
		state:               cloneState(loaded.State),
		lifecycleSetDigest:  loaded.LifecycleSetDigest,
		lifecycles:          cloneLifecycleRecords(loaded.Lifecycles, loaded.State.TimeHighWaterUnixSeconds),
		effectIDs:           resolved.EffectIDs,
		ownerSessionID:      resolved.OwnerSessionID,
		faults:              make(map[LifecycleStoreFault]error),
		issuedPermits:       make(map[approvalattempt.AttemptID]lifecyclestate.EffectPermit),
	}, nil
}

// OpenFixedFileStoreV1WithOwner is the G2-only v1 opener. It requires the
// opaque installation owner to be live before any store path inspection or
// read, binds the store to that exact owner session, and rechecks ownership
// before returning the validated store. It never creates or migrates state.
func OpenFixedFileStoreV1WithOwner(
	ctx context.Context,
	path string,
	owner installationowner.InstallationOwner,
	options FixedFileStoreV1Options,
) (*FixedFileStoreV1, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if owner == nil {
		return nil, ErrStoreOwnerRequired
	}
	if err := owner.CheckHeld(ctx); err != nil {
		return nil, fmt.Errorf("%w: initial owner check: %v", ErrStoreOwnerFenced, err)
	}
	ownerSessionID := owner.OwnerSessionID()
	if ownerSessionID.IsZero() {
		return nil, fmt.Errorf("%w: zero owner session", ErrStoreOwnerFenced)
	}
	if !options.OwnerSessionID.IsZero() && options.OwnerSessionID != ownerSessionID {
		return nil, ErrStoreOwnerSessionMismatch
	}
	options.OwnerSessionID = ownerSessionID
	store, err := openFixedFileStoreV1WithOptions(path, options)
	if err != nil {
		return nil, err
	}
	store.owner = owner
	store.ownerRequired = true
	if err := owner.CheckHeld(ctx); err != nil {
		store.ownerFenced = true
		store.recoveryRequired = true
		return nil, fmt.Errorf("%w: post-open owner check: %v", ErrStoreOwnerFenced, err)
	}
	if owner.OwnerSessionID() != store.ownerSessionID {
		store.ownerFenced = true
		store.recoveryRequired = true
		return nil, ErrStoreOwnerSessionMismatch
	}
	return store, nil
}

func (store *FixedFileStoreV1) snapshotV1(ctx context.Context) (fixedStoreV1Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return fixedStoreV1Snapshot{}, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.requireReadableLocked(); err != nil {
		return fixedStoreV1Snapshot{}, err
	}
	return fixedStoreV1Snapshot{
		SourceFormatVersion: store.sourceFormatVersion,
		State:               cloneState(store.state),
		LifecycleSetDigest:  store.lifecycleSetDigest,
		Lifecycles:          cloneLifecycleRecords(store.lifecycles, store.state.TimeHighWaterUnixSeconds),
	}, nil
}

// fixedStoreV1MigrationMu serializes every MigrateFixedFileStoreV0ToV1 call
// against the same store path, mirroring fixedStoreV2MigrationMu's guard
// against two concurrent callers independently building and renaming their
// own envelope onto the same destination.
var fixedStoreV1MigrationMu sync.Mutex

// MigrateFixedFileStoreV0ToV1 performs the sole E2 write path. It validates
// v0 under an asserted offline lock, commits exactly one empty lifecycle
// collection through temp/sync/rename/directory-sync, then reopens v1.
func MigrateFixedFileStoreV0ToV1(
	ctx context.Context,
	path string,
	options V0ToV1MigrationOptions,
) (*FixedFileStoreV1, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if options.Lock == nil {
		return nil, ErrMigrationLockRequired
	}
	fixedStoreV1MigrationMu.Lock()
	defer fixedStoreV1MigrationMu.Unlock()
	if err := options.Lock.CheckOfflineMigrationLock(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMigrationLockRequired, err)
	}
	if err := validateExistingStoreFile(path); err != nil {
		return nil, err
	}
	state, err := loadV0StateStrict(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStoreRepairRequired, err)
	}
	if err := options.Lock.CheckOfflineMigrationLock(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMigrationLockRequired, err)
	}
	envelope := encodedEnvelopeV1(state, nil)
	if err := persistEnvelopeV1(ctx, path, envelope, options.Lock, options.Faults); err != nil {
		return nil, err
	}
	if err := options.Lock.CheckOfflineMigrationLock(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMigrationLockRequired, err)
	}
	return OpenFixedFileStoreV1(path)
}

type loadedV1State struct {
	SourceFormatVersion uint64
	State               installationState
	LifecycleSetDigest  LifecycleSetDigest
	Lifecycles          []lifecyclestate.Record
}

func loadV1State(path string) (loadedV1State, error) {
	data, err := readBoundedStoreFile(path, MaxSupervisorStateV1Bytes)
	if err != nil {
		return loadedV1State{}, err
	}
	var envelope diskEnvelopeV1
	if err := decodeOneClosedJSON(data, &envelope); err != nil {
		return loadedV1State{}, errors.New("decode fixed supervisor v1 store")
	}
	if envelope.StoreFormatVersion == nil || *envelope.StoreFormatVersion != SupervisorStoreFormatV1 {
		return loadedV1State{}, errors.New("unsupported fixed supervisor store version")
	}
	if envelope.MigrationSourceVersion == nil || *envelope.MigrationSourceVersion != SupervisorStoreFormatV0 {
		return loadedV1State{}, errors.New("unsupported fixed supervisor migration source version")
	}
	if envelope.Lifecycles == nil {
		return loadedV1State{}, errors.New("fixed supervisor v1 lifecycle collection is missing")
	}
	if len(envelope.Lifecycles) > MaxRetainedLifecycleRecords {
		return loadedV1State{}, errors.New("fixed lifecycle state exceeds retained capacity")
	}
	if err := validateStateWithAttemptCapacity(envelope.State, false); err != nil {
		return loadedV1State{}, err
	}
	if err := validateV1AuthorityDigests(envelope.State); err != nil {
		return loadedV1State{}, err
	}
	records := make([]lifecyclestate.Record, len(envelope.Lifecycles))
	for index, disk := range envelope.Lifecycles {
		record, restoreErr := restoreLifecycleRecord(disk, envelope.State.TimeHighWaterUnixSeconds)
		if restoreErr != nil {
			return loadedV1State{}, fmt.Errorf("invalid lifecycle record %d", index)
		}
		records[index] = record
	}
	if err := validateV1State(envelope.State, records, envelope.LifecycleSetDigest); err != nil {
		return loadedV1State{}, err
	}
	return loadedV1State{
		SourceFormatVersion: *envelope.MigrationSourceVersion,
		State:               cloneState(envelope.State),
		LifecycleSetDigest:  envelope.LifecycleSetDigest,
		Lifecycles:          cloneLifecycleRecords(records, envelope.State.TimeHighWaterUnixSeconds),
	}, nil
}

func loadV0StateStrict(path string) (installationState, error) {
	data, err := readBoundedStoreFile(path, MaxSupervisorStateBytes)
	if err != nil {
		return installationState{}, err
	}
	var envelope struct {
		StoreFormatVersion *uint64           `json:"storeFormatVersion"`
		State              installationState `json:"state"`
	}
	if err := decodeOneClosedJSON(data, &envelope); err != nil {
		return installationState{}, errors.New("decode fixed supervisor v0 store")
	}
	if envelope.StoreFormatVersion == nil || *envelope.StoreFormatVersion != SupervisorStoreFormatV0 {
		return installationState{}, errors.New("unsupported fixed supervisor v0 store version")
	}
	if err := validateState(envelope.State); err != nil {
		return installationState{}, err
	}
	if err := validateV1AuthorityDigests(envelope.State); err != nil {
		return installationState{}, err
	}
	return cloneState(envelope.State), nil
}

func readBoundedStoreFile(path string, maxBytes int64) ([]byte, error) {
	file, err := os.Open(path) //nolint:gosec // G304: path is the store's own installation-configured path, not caller/network-supplied
	if err != nil {
		return nil, errors.New("open fixed supervisor store")
	}
	defer func() { _ = file.Close() }()
	data, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
	if err != nil {
		return nil, errors.New("read fixed supervisor store")
	}
	if int64(len(data)) > maxBytes {
		return nil, errors.New("fixed supervisor store exceeds byte capacity")
	}
	return data, nil
}

func decodeOneClosedJSON(data []byte, destination any) error {
	if err := rejectDuplicateJSONNames(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("fixed supervisor store has trailing data")
	}
	return nil
}

func rejectDuplicateJSONNames(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	var readValue func() error
	readValue = func() error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		delimiter, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delimiter {
		case '{':
			seen := make(map[string]struct{})
			for decoder.More() {
				nameToken, err := decoder.Token()
				if err != nil {
					return err
				}
				name, ok := nameToken.(string)
				if !ok {
					return errors.New("fixed supervisor store has a non-string object name")
				}
				if _, exists := seen[name]; exists {
					return errors.New("fixed supervisor store has a duplicate object name")
				}
				seen[name] = struct{}{}
				if err := readValue(); err != nil {
					return err
				}
			}
			closing, err := decoder.Token()
			if err != nil || closing != json.Delim('}') {
				return errors.New("fixed supervisor store has an incomplete object")
			}
			return nil
		case '[':
			for decoder.More() {
				if err := readValue(); err != nil {
					return err
				}
			}
			closing, err := decoder.Token()
			if err != nil || closing != json.Delim(']') {
				return errors.New("fixed supervisor store has an incomplete array")
			}
			return nil
		default:
			return errors.New("fixed supervisor store has an unexpected delimiter")
		}
	}
	if err := readValue(); err != nil {
		return err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return errors.New("fixed supervisor store has trailing data")
	}
	return nil
}

func validateExistingStoreFile(path string) error {
	info, err := os.Lstat(path) //nolint:gosec // G703: path is the installation-configured store file validated here before any open or mutation.
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New("fixed supervisor store is missing")
		}
		return errors.New("inspect fixed supervisor store")
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return errors.New("fixed supervisor store must be a regular file")
	}
	if info.Mode().Perm() != 0o600 {
		return errors.New("fixed supervisor store permissions must be 0600")
	}
	return nil
}

func encodedEnvelopeV1(state installationState, records []lifecyclestate.Record) diskEnvelopeV1 {
	version := SupervisorStoreFormatV1
	source := SupervisorStoreFormatV0
	disks := make([]lifecycleRecordDisk, len(records))
	for index, record := range records {
		disks[index] = lifecycleRecordToDisk(record)
	}
	return diskEnvelopeV1{
		StoreFormatVersion:     &version,
		MigrationSourceVersion: &source,
		State:                  cloneState(state),
		LifecycleSetDigest:     lifecycleSetDigest(records),
		Lifecycles:             disks,
	}
}

func persistEnvelopeV1(
	ctx context.Context,
	path string,
	envelope diskEnvelopeV1,
	lock OfflineMigrationLock,
	faults MigrationFaultInjector,
) error {
	parent := filepath.Dir(path)
	temporary, err := os.CreateTemp(parent, ".capsule-supervisor-v1-migration-*.tmp")
	if err != nil {
		return errors.New("create fixed supervisor v1 migration")
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
		return errors.New("protect fixed supervisor v1 migration")
	}
	encoder := json.NewEncoder(temporary)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(envelope); err != nil {
		return errors.New("encode fixed supervisor v1 migration")
	}
	if err := temporary.Sync(); err != nil {
		return errors.New("sync fixed supervisor v1 migration")
	}
	if err := temporary.Close(); err != nil {
		return errors.New("close fixed supervisor v1 migration")
	}
	if faults != nil {
		if err := faults.FailMigrationAt(FaultMigrationBeforeRename); err != nil {
			return errors.New("fixed supervisor v1 migration confirmed before rename failure")
		}
	}
	if err := lock.CheckOfflineMigrationLock(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrMigrationLockRequired, err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return errors.New("commit fixed supervisor v1 migration")
	}
	committed = true
	if faults != nil {
		if err := faults.FailMigrationAt(FaultMigrationAfterRename); err != nil {
			return fmt.Errorf("%w: injected after rename", ErrMigrationOutcomeIndeterminate)
		}
	}
	directory, err := os.Open(parent) //nolint:gosec // G304: parent is the store's own installation-configured directory, not caller/network-supplied
	if err != nil {
		return fmt.Errorf("%w: open migration directory", ErrMigrationOutcomeIndeterminate)
	}
	syncErr := directory.Sync()
	closeErr := directory.Close()
	if syncErr != nil || closeErr != nil {
		return fmt.Errorf("%w: sync migration directory", ErrMigrationOutcomeIndeterminate)
	}
	return nil
}

func validateV1State(
	state installationState,
	records []lifecyclestate.Record,
	digest LifecycleSetDigest,
) error {
	if err := validateStateWithAttemptCapacity(state, false); err != nil {
		return err
	}
	if err := validateV1AuthorityDigests(state); err != nil {
		return err
	}
	if len(records) > MaxRetainedLifecycleRecords {
		return errors.New("fixed lifecycle state exceeds retained capacity")
	}
	if lifecycleSetDigest(records) != digest {
		return errors.New("fixed lifecycle set digest mismatch")
	}
	attemptIDs := make(map[approvalattempt.AttemptID]struct{}, len(records))
	effectIDs := make(map[lifecyclestate.EffectID]struct{}, len(records))
	instanceDigests := make(map[lifecyclestate.BackendInstanceDigest]struct{}, len(records))
	var previous approvalattempt.AttemptID
	for index, record := range records {
		if err := record.Validate(state.TimeHighWaterUnixSeconds); err != nil {
			return errors.New("fixed lifecycle state has an invalid record")
		}
		view := record.View()
		bindings := view.Bindings.View()
		if index > 0 && bytes.Compare(previous[:], bindings.AttemptID[:]) >= 0 {
			return errors.New("fixed lifecycle records are not strictly sorted")
		}
		previous = bindings.AttemptID
		if _, exists := attemptIDs[bindings.AttemptID]; exists {
			return errors.New("fixed lifecycle state has a duplicate attempt ID")
		}
		attemptIDs[bindings.AttemptID] = struct{}{}
		if !view.EffectID.IsZero() {
			if _, exists := effectIDs[view.EffectID]; exists {
				return errors.New("fixed lifecycle state has a duplicate effect ID")
			}
			effectIDs[view.EffectID] = struct{}{}
		}
		if view.Instance.Present() {
			instanceDigest := view.Instance.Digest()
			if _, exists := instanceDigests[instanceDigest]; exists {
				return errors.New("fixed lifecycle state has a duplicate instance identity")
			}
			instanceDigests[instanceDigest] = struct{}{}
		}
		if err := validateLifecycleAuthorityCrossLinks(state, record); err != nil {
			return err
		}
	}
	if activeLifecycleWork(state, records) > MaxActiveLifecycleRecords {
		return errors.New("fixed lifecycle state exceeds active capacity")
	}
	return nil
}

func validateV1AuthorityDigests(state installationState) error {
	approvalDigest, err := approvalSetDigest(state.Approvals)
	if err != nil {
		return fmt.Errorf("fixed v1 compute approval set digest: %w", err)
	}
	attemptDigest, err := attemptSetDigest(state.Attempts)
	if err != nil {
		return fmt.Errorf("fixed v1 compute attempt set digest: %w", err)
	}
	switch {
	case len(state.Registrations) == 0 && state.RegistrationSetDigest != emptyRegistrationSetDigest():
		return errors.New("fixed v1 empty registration set digest mismatch")
	case len(state.Registrations) != 0 && state.RegistrationSetDigest == ([32]byte{}):
		return errors.New("fixed v1 registration set digest is missing")
	case approvalDigest != state.ApprovalSetDigest:
		return errors.New("fixed v1 approval set digest mismatch")
	case attemptDigest != state.AttemptSetDigest:
		return errors.New("fixed v1 attempt set digest mismatch")
	default:
		return nil
	}
}

func validateLifecycleAuthorityCrossLinks(state installationState, record lifecyclestate.Record) error {
	binding := record.Bindings().View()
	attemptIndex := -1
	for index, attempt := range state.Attempts {
		if attempt.AttemptID == binding.AttemptID {
			attemptIndex = index
			break
		}
	}
	if attemptIndex < 0 {
		return errors.New("fixed lifecycle state lacks its committed attempt")
	}
	attempt := state.Attempts[attemptIndex]
	if attempt.ApprovalID != binding.ApprovalID ||
		attempt.AttemptNonce != binding.AttemptNonce ||
		attempt.RegistrationID != binding.RegistrationID ||
		attempt.RegistrationSequence != binding.RegistrationSequence ||
		attempt.PlanDigest != binding.PlanDigest ||
		attempt.InstallationID != binding.InstallationID ||
		attempt.EpochSequence != binding.EpochSequence ||
		attempt.EpochDigest != binding.EpochDigest ||
		attempt.SupervisorID != binding.SupervisorID ||
		attempt.Purpose != binding.ApprovalPurpose ||
		attempt.Audience != binding.ApprovalAudience ||
		attempt.ApprovalPayloadDigest != binding.ApprovalPayloadDigest ||
		attempt.AuthorizationIdentity != binding.AuthorizationIdentity ||
		attempt.CreatedAt != binding.AttemptCreatedAt ||
		attempt.StorageFormatVersion != binding.AttemptStorageVersion {
		return errors.New("fixed lifecycle/attempt cross-link mismatch")
	}
	approvalIndex := -1
	for index, approval := range state.Approvals {
		if approval.ApprovalID == binding.ApprovalID {
			approvalIndex = index
			break
		}
	}
	if approvalIndex < 0 || state.Approvals[approvalIndex].State != approvalattempt.ApprovalConsumed ||
		state.Approvals[approvalIndex].ConsumedAttemptID != binding.AttemptID ||
		state.Approvals[approvalIndex].StorageFormatVersion != binding.ApprovalStorageVersion {
		return errors.New("fixed lifecycle/approval cross-link mismatch")
	}
	registrationIndex := findRegistration(state.Registrations, binding.RegistrationID)
	if registrationIndex < 0 {
		return errors.New("fixed lifecycle state lacks its registration")
	}
	registration := state.Registrations[registrationIndex]
	if registration.Index.RegistrationSequence != binding.RegistrationSequence ||
		registration.Record.RecomputedPlanDigest != binding.PlanDigest ||
		registration.Index.InstallationID != binding.InstallationID ||
		registration.Index.EpochSequence != binding.EpochSequence ||
		registration.Index.EpochDigest != binding.EpochDigest ||
		registration.Index.SupervisorID != binding.SupervisorID ||
		registration.Record.StorageFormatVersion != binding.RegistrationStorageVersion {
		return errors.New("fixed lifecycle/registration cross-link mismatch")
	}
	roles := registration.PlanBindings
	if roles.SourceManifestDigest != binding.SourceManifestDigest ||
		roles.InlineInputDigest != binding.InlineInputDigest ||
		roles.RuntimeBundleManifestDigest != binding.RuntimeBundleManifestDigest ||
		!reflect.DeepEqual(roles.ProfileReviewAttestationDigests, binding.ProfileReviewAttestationDigests) ||
		roles.ProfileRegistryEntryDigest != binding.ProfileRegistryEntryDigest ||
		roles.BackendValidationRecordDigest != binding.BackendValidationRecordDigest ||
		roles.BackendConfigurationDigest != binding.BackendConfigurationDigest ||
		roles.TrustSnapshotDigest != binding.TrustSnapshotDigest ||
		roles.PolicyDecisionDigest != binding.PolicyDecisionDigest {
		return errors.New("fixed lifecycle/plan-role cross-link mismatch")
	}
	return nil
}

func lifecycleSetDigest(records []lifecyclestate.Record) LifecycleSetDigest {
	disks := make([]lifecycleRecordDisk, len(records))
	for index, record := range records {
		disks[index] = lifecycleRecordToDisk(record)
	}
	sort.Slice(disks, func(left, right int) bool {
		return bytes.Compare(disks[left].Bindings.AttemptID[:], disks[right].Bindings.AttemptID[:]) < 0
	})
	hash := sha256.New()
	_, _ = hash.Write([]byte("capsule.supervisor.lifecycle-set.v1"))
	var count [8]byte
	binary.BigEndian.PutUint64(count[:], uint64(len(disks)))
	_, _ = hash.Write(count[:])
	for _, disk := range disks {
		encoded, err := json.Marshal(disk)
		if err != nil {
			panic(fmt.Sprintf("encode closed lifecycle record: %v", err))
		}
		var length [8]byte
		binary.BigEndian.PutUint64(length[:], uint64(len(encoded)))
		_, _ = hash.Write(length[:])
		_, _ = hash.Write(encoded)
	}
	var digest LifecycleSetDigest
	copy(digest[:], hash.Sum(nil))
	return digest
}

func lifecycleRecordToDisk(record lifecyclestate.Record) lifecycleRecordDisk {
	view := record.View()
	binding := view.Bindings.View()
	return lifecycleRecordDisk{
		FormatVersion:              view.FormatVersion,
		RecordVersion:              view.RecordVersion,
		SnapshotGeneration:         view.SnapshotGeneration,
		Bindings:                   lifecycleBindingsToDisk(binding),
		ImmutableBindingDigest:     view.ImmutableBindingDigest,
		OperationSequence:          view.OperationSequence,
		Operation:                  view.Operation,
		EffectID:                   view.EffectID,
		EffectStatus:               view.EffectStatus,
		Instance:                   view.Instance.View(),
		State:                      view.State,
		CleanupRequired:            view.CleanupRequired,
		LastConfirmedCheckpoint:    view.LastConfirmedCheckpoint,
		FirstFailure:               view.FirstFailure,
		FailureOperation:           view.FailureOperation,
		LastReconciliation:         view.LastReconciliation,
		AutomaticRecoveryCount:     view.AutomaticRecoveryCount,
		NextRecoveryAt:             view.NextRecoveryAt,
		RecoveryFence:              view.RecoveryFence,
		AutomaticRecoveryExhausted: view.AutomaticRecoveryExhausted,
		OpenedAt:                   view.OpenedAt,
		LastTransitionAt:           view.LastTransitionAt,
		TerminalAt:                 view.TerminalAt,
	}
}

func lifecycleBindingsToDisk(view lifecyclestate.ImmutableBindingsView) lifecycleBindingsDisk {
	return lifecycleBindingsDisk{
		AttemptID:                       view.AttemptID,
		ApprovalID:                      view.ApprovalID,
		AttemptNonce:                    view.AttemptNonce,
		RegistrationID:                  view.RegistrationID,
		RegistrationSequence:            view.RegistrationSequence,
		PlanDigest:                      view.PlanDigest,
		InstallationID:                  view.InstallationID,
		EpochSequence:                   view.EpochSequence,
		EpochDigest:                     view.EpochDigest,
		SupervisorID:                    view.SupervisorID,
		ApprovalPurpose:                 view.ApprovalPurpose,
		ApprovalAudience:                view.ApprovalAudience,
		ApprovalPayloadDigest:           view.ApprovalPayloadDigest,
		AuthorizationIdentity:           view.AuthorizationIdentity,
		AttemptCreatedAt:                view.AttemptCreatedAt,
		AttemptStorageVersion:           view.AttemptStorageVersion,
		ApprovalStorageVersion:          view.ApprovalStorageVersion,
		RegistrationStorageVersion:      view.RegistrationStorageVersion,
		SourceManifestDigest:            view.SourceManifestDigest,
		InlineInputDigest:               view.InlineInputDigest,
		RuntimeBundleManifestDigest:     view.RuntimeBundleManifestDigest,
		ProfileReviewAttestationDigests: append([]v0candidate.ProfileReviewAttestationDigest(nil), view.ProfileReviewAttestationDigests...),
		ProfileRegistryEntryDigest:      view.ProfileRegistryEntryDigest,
		BackendValidationRecordDigest:   view.BackendValidationRecordDigest,
		BackendConfigurationDigest:      view.BackendConfigurationDigest,
		TrustSnapshotDigest:             view.TrustSnapshotDigest,
		PolicyDecisionDigest:            view.PolicyDecisionDigest,
		Backend:                         view.Backend.View(),
	}
}

func restoreLifecycleRecord(disk lifecycleRecordDisk, highWater v0candidate.UInt53) (lifecyclestate.Record, error) {
	backend, err := lifecyclestate.NewBackendBinding(disk.Bindings.Backend)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	bindingView := lifecyclestate.ImmutableBindingsView{
		AttemptID:                       disk.Bindings.AttemptID,
		ApprovalID:                      disk.Bindings.ApprovalID,
		AttemptNonce:                    disk.Bindings.AttemptNonce,
		RegistrationID:                  disk.Bindings.RegistrationID,
		RegistrationSequence:            disk.Bindings.RegistrationSequence,
		PlanDigest:                      disk.Bindings.PlanDigest,
		InstallationID:                  disk.Bindings.InstallationID,
		EpochSequence:                   disk.Bindings.EpochSequence,
		EpochDigest:                     disk.Bindings.EpochDigest,
		SupervisorID:                    disk.Bindings.SupervisorID,
		ApprovalPurpose:                 disk.Bindings.ApprovalPurpose,
		ApprovalAudience:                disk.Bindings.ApprovalAudience,
		ApprovalPayloadDigest:           disk.Bindings.ApprovalPayloadDigest,
		AuthorizationIdentity:           disk.Bindings.AuthorizationIdentity,
		AttemptCreatedAt:                disk.Bindings.AttemptCreatedAt,
		AttemptStorageVersion:           disk.Bindings.AttemptStorageVersion,
		ApprovalStorageVersion:          disk.Bindings.ApprovalStorageVersion,
		RegistrationStorageVersion:      disk.Bindings.RegistrationStorageVersion,
		SourceManifestDigest:            disk.Bindings.SourceManifestDigest,
		InlineInputDigest:               disk.Bindings.InlineInputDigest,
		RuntimeBundleManifestDigest:     disk.Bindings.RuntimeBundleManifestDigest,
		ProfileReviewAttestationDigests: append([]v0candidate.ProfileReviewAttestationDigest(nil), disk.Bindings.ProfileReviewAttestationDigests...),
		ProfileRegistryEntryDigest:      disk.Bindings.ProfileRegistryEntryDigest,
		BackendValidationRecordDigest:   disk.Bindings.BackendValidationRecordDigest,
		BackendConfigurationDigest:      disk.Bindings.BackendConfigurationDigest,
		TrustSnapshotDigest:             disk.Bindings.TrustSnapshotDigest,
		PolicyDecisionDigest:            disk.Bindings.PolicyDecisionDigest,
		Backend:                         backend,
	}
	bindings, err := lifecyclestate.RestoreImmutableBindings(bindingView, disk.ImmutableBindingDigest)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	var instance lifecyclestate.BackendInstanceIdentity
	if disk.Instance.Kind != "" || len(disk.Instance.Value) != 0 || disk.Instance.Digest != (lifecyclestate.BackendInstanceDigest{}) {
		instance, err = lifecyclestate.RestoreBackendInstanceIdentity(disk.Instance)
		if err != nil {
			return lifecyclestate.Record{}, err
		}
	}
	return lifecyclestate.NewRecord(lifecyclestate.RecordView{
		FormatVersion:              disk.FormatVersion,
		RecordVersion:              disk.RecordVersion,
		SnapshotGeneration:         disk.SnapshotGeneration,
		Bindings:                   bindings,
		ImmutableBindingDigest:     disk.ImmutableBindingDigest,
		OperationSequence:          disk.OperationSequence,
		Operation:                  disk.Operation,
		EffectID:                   disk.EffectID,
		EffectStatus:               disk.EffectStatus,
		Instance:                   instance,
		State:                      disk.State,
		CleanupRequired:            disk.CleanupRequired,
		LastConfirmedCheckpoint:    disk.LastConfirmedCheckpoint,
		FirstFailure:               disk.FirstFailure,
		FailureOperation:           disk.FailureOperation,
		LastReconciliation:         disk.LastReconciliation,
		AutomaticRecoveryCount:     disk.AutomaticRecoveryCount,
		NextRecoveryAt:             disk.NextRecoveryAt,
		RecoveryFence:              disk.RecoveryFence,
		AutomaticRecoveryExhausted: disk.AutomaticRecoveryExhausted,
		OpenedAt:                   disk.OpenedAt,
		LastTransitionAt:           disk.LastTransitionAt,
		TerminalAt:                 disk.TerminalAt,
	}, highWater)
}

func cloneLifecycleRecords(records []lifecyclestate.Record, highWater v0candidate.UInt53) []lifecyclestate.Record {
	cloned := make([]lifecyclestate.Record, len(records))
	for index, record := range records {
		restored, err := restoreLifecycleRecord(lifecycleRecordToDisk(record), highWater)
		if err != nil {
			panic(fmt.Sprintf("clone validated lifecycle record: %v", err))
		}
		cloned[index] = restored
	}
	return cloned
}
