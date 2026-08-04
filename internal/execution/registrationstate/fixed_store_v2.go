package registrationstate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	FaultV2MigrationAfterTempCreate MigrationFault = "v2-migration-after-temp-create"
	FaultV2MigrationAfterTempMode   MigrationFault = "v2-migration-after-temp-mode"
	FaultV2MigrationAfterTempWrite  MigrationFault = "v2-migration-after-temp-write" //nolint:gosec // G101: fixed fault label, not a credential.
	FaultV2MigrationAfterTempSync   MigrationFault = "v2-migration-after-temp-sync"
	FaultV2MigrationAfterTempClose  MigrationFault = "v2-migration-after-temp-close"
	FaultV2MigrationBeforeRename    MigrationFault = "v2-migration-before-rename"
	FaultV2MigrationAfterRename     MigrationFault = "v2-migration-after-rename"
	FaultV2MigrationAfterDirSync    MigrationFault = "v2-migration-after-directory-sync"
)

const migrationGenesisKind = "migration-genesis"

// V1ToV2MigrationOptions contains only the asserted installation owner and a
// local fault-injection port. It carries no records, paths, adapter values, or
// replacement authority bytes.
type V1ToV2MigrationOptions struct {
	Lock   OfflineMigrationLock
	Faults MigrationFaultInjector
}

// FixedFileStoreV2 is the completely verified fixed-oracle v2 view. F2 adds
// migration genesis and F3 adds exactly one sealed immutable-segment
// activation. Neither slice authorizes retained lookup, v2 authority mutation,
// adapter invocation, a consumer, or a guest.
type FixedFileStoreV2 struct {
	path       string
	state      installationState
	lifecycles []lifecyclestate.Record
	active     archivestate.ActiveStateV2
	genesis    archivestate.MigrationGenesisCheckpoint
	segments   []loadedArchiveSegmentV0
	orphans    uint16
}

// FixedFileStoreV2VerificationReport is the bounded, content-free F2/F3 full
// verification result. It deliberately includes no path or retained bytes.
type FixedFileStoreV2VerificationReport struct {
	StoreFormatVersion     uint64
	MigrationSourceVersion uint64
	SnapshotGeneration     v0candidate.PositiveUInt53
	ArchiveGeneration      archivestate.ArchiveGeneration
	DurableTimeHighWater   v0candidate.UInt53
	CurrentCheckpoint      archivestate.ArchiveCheckpointReference
	HotCounts              archivestate.ArchiveCounts
	ArchivedCounts         archivestate.ArchiveCounts
	TotalCounts            archivestate.ArchiveCounts
	SegmentCount           uint16
	OrphanSegmentCount     uint16
}

type fixedStoreV2Snapshot struct {
	State      installationState
	Lifecycles []lifecyclestate.Record
	Active     archivestate.ActiveStateV2
	Genesis    archivestate.MigrationGenesisCheckpoint
}

type diskEnvelopeV2 struct {
	StoreFormatVersion     *uint64 `json:"storeFormatVersion"`
	MigrationSourceVersion *uint64 `json:"migrationSourceVersion"`

	SnapshotGeneration v0candidate.PositiveUInt53     `json:"snapshotGeneration"`
	ArchiveGeneration  archivestate.ArchiveGeneration `json:"archiveGeneration"`

	State              installationState     `json:"state"`
	LifecycleSetDigest LifecycleSetDigest    `json:"lifecycleSetDigest"`
	Lifecycles         []lifecycleRecordDisk `json:"lifecycles"`

	Descriptors         []archivestate.ArchiveDescriptorView    `json:"descriptors"`
	DescriptorSetDigest archivestate.ArchiveDescriptorSetDigest `json:"descriptorSetDigest"`
	Indexes             archiveIndexesDisk                      `json:"indexes"`
	IndexDigests        archivestate.ArchiveIndexDigests        `json:"indexDigests"`
	CombinedIndexDigest archivestate.ArchiveCombinedIndexDigest `json:"combinedIndexDigest"`

	EffectTombstoneCoverage   string                               `json:"effectTombstoneCoverage"`
	VisibleV1EffectSeedCount  uint64                               `json:"visibleV1EffectSeedCount"`
	VisibleV1EffectSeedDigest archivestate.ArchiveEffectSeedDigest `json:"visibleV1EffectSeedDigest"`

	PreviousCheckpoint   archivestate.ArchiveCheckpointReference `json:"previousCheckpoint"`
	CurrentCheckpoint    archivestate.ArchiveCheckpointReference `json:"currentCheckpoint"`
	MigrationGenesis     migrationGenesisDisk                    `json:"migrationGenesis"`
	ActivationCheckpoint *archiveCheckpointDisk                  `json:"activationCheckpoint,omitempty"`

	HotCounts      archivestate.ArchiveCounts `json:"hotCounts"`
	ArchivedCounts archivestate.ArchiveCounts `json:"archivedCounts"`
	TotalCounts    archivestate.ArchiveCounts `json:"totalCounts"`
}

type migrationGenesisDisk struct {
	Kind                      string                                  `json:"kind"`
	StoreFormatVersion        uint64                                  `json:"storeFormatVersion"`
	MigrationSourceVersion    uint64                                  `json:"migrationSourceVersion"`
	ResultSnapshotGeneration  v0candidate.PositiveUInt53              `json:"resultSnapshotGeneration"`
	ArchiveGeneration         archivestate.ArchiveGeneration          `json:"archiveGeneration"`
	DescriptorSetDigest       archivestate.ArchiveDescriptorSetDigest `json:"descriptorSetDigest"`
	IndexDigests              archivestate.ArchiveIndexDigests        `json:"indexDigests"`
	CombinedIndexDigest       archivestate.ArchiveCombinedIndexDigest `json:"combinedIndexDigest"`
	HotSetDigests             archivestate.HotSetDigests              `json:"hotSetDigests"`
	VisibleV1EffectSeed       []lifecyclestate.EffectID               `json:"visibleV1EffectSeed"`
	VisibleV1EffectSeedDigest archivestate.ArchiveEffectSeedDigest    `json:"visibleV1EffectSeedDigest"`
	HotCounts                 archivestate.ArchiveCounts              `json:"hotCounts"`
	InstallationID            v0candidate.InstallationID              `json:"installationId"`
	SupervisorID              v0candidate.SupervisorID                `json:"supervisorId"`
	EpochSequence             v0candidate.UInt53                      `json:"epochSequence"`
	EpochDigest               v0candidate.TrustEpochDigest            `json:"epochDigest"`
	DurableTimeHighWater      v0candidate.UInt53                      `json:"durableTimeHighWater"`
	Digest                    archivestate.ArchiveCheckpointDigest    `json:"digest"`
}

type recordLocationDisk struct {
	Kind             archivestate.RecordLocationKind `json:"kind"`
	RecordKind       archivestate.RecordKind         `json:"recordKind"`
	HotRecordOrdinal v0candidate.PositiveUInt53      `json:"hotRecordOrdinal"`
	Archive          archivestate.ArchiveLocation    `json:"archive"`
}

type attemptLifecycleDisk struct {
	Presence         archivestate.AttemptLifecyclePresence `json:"presence"`
	State            lifecyclestate.LifecycleState         `json:"state"`
	Location         recordLocationDisk                    `json:"location"`
	FullRecordDigest archivestate.ArchiveRecordDigest      `json:"fullRecordDigest"`
}

type registrationIndexEntryDisk struct {
	RegistrationID       v0candidate.RegistrationID       `json:"registrationId"`
	RegistrationSequence v0candidate.PositiveUInt53       `json:"registrationSequence"`
	PlanDigest           v0candidate.ExecutionPlanDigest  `json:"planDigest"`
	ExpiresAt            v0candidate.UInt53               `json:"expiresAt"`
	Location             recordLocationDisk               `json:"location"`
	FullRecordDigest     archivestate.ArchiveRecordDigest `json:"fullRecordDigest"`
}

type approvalIndexEntryDisk struct {
	ApprovalID            approvalattempt.ApprovalID                       `json:"approvalId"`
	RegistrationID        v0candidate.RegistrationID                       `json:"registrationId"`
	PayloadDigest         approvalattempt.ApprovalPayloadDigest            `json:"payloadDigest"`
	AuthorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity `json:"authorizationIdentity"`
	AttemptNonce          approvalattempt.AttemptNonce                     `json:"attemptNonce"`
	State                 approvalattempt.ApprovalState                    `json:"state"`
	ConsumedAttemptID     approvalattempt.AttemptID                        `json:"consumedAttemptId"`
	ExpiresAt             v0candidate.UInt53                               `json:"expiresAt"`
	Location              recordLocationDisk                               `json:"location"`
	FullRecordDigest      archivestate.ArchiveRecordDigest                 `json:"fullRecordDigest"`
}

type attemptIndexEntryDisk struct {
	AttemptID        approvalattempt.AttemptID        `json:"attemptId"`
	ApprovalID       approvalattempt.ApprovalID       `json:"approvalId"`
	RegistrationID   v0candidate.RegistrationID       `json:"registrationId"`
	CreatedAt        v0candidate.UInt53               `json:"createdAt"`
	Lifecycle        attemptLifecycleDisk             `json:"lifecycle"`
	Location         recordLocationDisk               `json:"location"`
	FullRecordDigest archivestate.ArchiveRecordDigest `json:"fullRecordDigest"`
}

type nonceIndexEntryDisk struct {
	AttemptNonce  approvalattempt.AttemptNonce          `json:"attemptNonce"`
	PayloadDigest approvalattempt.ApprovalPayloadDigest `json:"payloadDigest"`
	ApprovalID    approvalattempt.ApprovalID            `json:"approvalId"`
	Location      recordLocationDisk                    `json:"location"`
}

type effectIndexEntryDisk struct {
	EffectID                   lifecyclestate.EffectID    `json:"effectId"`
	AttemptID                  approvalattempt.AttemptID  `json:"attemptId"`
	OperationSequence          v0candidate.PositiveUInt53 `json:"operationSequence"`
	Operation                  lifecyclestate.Operation   `json:"operation"`
	IssuanceSnapshotGeneration v0candidate.PositiveUInt53 `json:"issuanceSnapshotGeneration"`
	VisibleV1Seed              bool                       `json:"visibleV1Seed"`
	AttemptLocation            recordLocationDisk         `json:"attemptLocation"`
}

type instanceIndexEntryDisk struct {
	InstanceDigest lifecyclestate.BackendInstanceDigest `json:"instanceDigest"`
	AttemptID      approvalattempt.AttemptID            `json:"attemptId"`
	Location       recordLocationDisk                   `json:"location"`
}

type approvalReplayIndexEntryDisk struct {
	PayloadDigest         approvalattempt.ApprovalPayloadDigest            `json:"payloadDigest"`
	ExactPayloadDigest    archivestate.ArchiveRecordDigest                 `json:"exactPayloadDigest"`
	AuthorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity `json:"authorizationIdentity"`
	ApprovalID            approvalattempt.ApprovalID                       `json:"approvalId"`
	State                 approvalattempt.ApprovalState                    `json:"state"`
	Location              recordLocationDisk                               `json:"location"`
}

type attemptReplayIndexEntryDisk struct {
	RegistrationID v0candidate.RegistrationID   `json:"registrationId"`
	ApprovalID     approvalattempt.ApprovalID   `json:"approvalId"`
	AttemptID      approvalattempt.AttemptID    `json:"attemptId"`
	State          approvalattempt.AttemptState `json:"state"`
	Location       recordLocationDisk           `json:"location"`
}

type archiveIndexesDisk struct {
	Scope          archivestate.ArchiveIndexScope `json:"scope"`
	Registrations  []registrationIndexEntryDisk   `json:"registrations"`
	Approvals      []approvalIndexEntryDisk       `json:"approvals"`
	Attempts       []attemptIndexEntryDisk        `json:"attempts"`
	Nonces         []nonceIndexEntryDisk          `json:"nonces"`
	Effects        []effectIndexEntryDisk         `json:"effects"`
	Instances      []instanceIndexEntryDisk       `json:"instances"`
	ApprovalReplay []approvalReplayIndexEntryDisk `json:"approvalReplay"`
	AttemptReplay  []attemptReplayIndexEntryDisk  `json:"attemptReplay"`
}

type loadedV2State struct {
	Envelope   diskEnvelopeV2
	State      installationState
	Lifecycles []lifecyclestate.Record
	Active     archivestate.ActiveStateV2
	Genesis    archivestate.MigrationGenesisCheckpoint
	Segments   []loadedArchiveSegmentV0
	Orphans    uint16
}

var fixedStoreV2MigrationMu sync.Mutex

// OpenFixedFileStoreV2 opens only an existing, complete, version-specific v2
// world and runs the full fixed-oracle verifier. It never creates, migrates,
// downgrades, normalizes, repairs, or rewrites a store or segment.
func OpenFixedFileStoreV2(path string) (*FixedFileStoreV2, error) {
	if path == "" {
		return nil, errors.New("fixed supervisor v2 store path is required")
	}
	if err := validateExistingStoreFile(path); err != nil {
		return nil, err
	}
	loaded, err := loadV2State(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStoreRepairRequired, err)
	}
	return &FixedFileStoreV2{
		path:       path,
		state:      cloneState(loaded.State),
		lifecycles: cloneLifecycleRecords(loaded.Lifecycles, loaded.State.TimeHighWaterUnixSeconds),
		active:     loaded.Active, genesis: loaded.Genesis,
		segments: cloneLoadedArchiveSegments(loaded.Segments), orphans: loaded.Orphans,
	}, nil
}

// VerifyFixedFileStoreV2 runs the same complete read-only validation as the v2
// opener and returns only fixed metadata and counts.
func VerifyFixedFileStoreV2(path string) (FixedFileStoreV2VerificationReport, error) {
	store, err := OpenFixedFileStoreV2(path)
	if err != nil {
		return FixedFileStoreV2VerificationReport{}, err
	}
	view := store.active.View()
	return FixedFileStoreV2VerificationReport{
		StoreFormatVersion: view.StoreFormatVersion, MigrationSourceVersion: view.MigrationSourceVersion,
		SnapshotGeneration: view.SnapshotGeneration, ArchiveGeneration: view.ArchiveGeneration,
		DurableTimeHighWater: view.DurableTimeHighWater, CurrentCheckpoint: view.CurrentCheckpoint,
		HotCounts: view.HotCounts, ArchivedCounts: view.ArchivedCounts, TotalCounts: view.TotalCounts,
		SegmentCount:       uint16(len(store.segments)), //nolint:gosec // G115: full verifier caps descriptors at 64 before construction.
		OrphanSegmentCount: store.orphans,
	}, nil
}

func (store *FixedFileStoreV2) snapshot(ctx context.Context) (fixedStoreV2Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return fixedStoreV2Snapshot{}, err
	}
	return fixedStoreV2Snapshot{
		State:      cloneState(store.state),
		Lifecycles: cloneLifecycleRecords(store.lifecycles, store.state.TimeHighWaterUnixSeconds),
		Active:     store.active, Genesis: store.genesis,
	}, nil
}

// MigrateFixedFileStoreV1ToV2 performs the sole stateful F2 write. It fully
// validates v1 before any write, maps every record into an all-hot retained
// projection without lifecycle invention, commits once, then owner-checks and
// fully reopens the v2 result.
func MigrateFixedFileStoreV1ToV2(
	ctx context.Context,
	path string,
	options V1ToV2MigrationOptions,
) (*FixedFileStoreV2, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if options.Lock == nil {
		return nil, ErrMigrationLockRequired
	}
	fixedStoreV2MigrationMu.Lock()
	defer fixedStoreV2MigrationMu.Unlock()
	if err := options.Lock.CheckOfflineMigrationLock(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMigrationLockRequired, err)
	}
	if err := validateExistingStoreFile(path); err != nil {
		return nil, err
	}
	source, err := loadV1State(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrStoreRepairRequired, err)
	}
	if source.State.TrustPhase != TrustStable || source.State.Quarantined || source.State.AttemptsDisabled {
		return nil, fmt.Errorf("%w: fixed supervisor v1 state is fenced", ErrStoreRepairRequired)
	}
	envelope, err := buildEnvelopeV2(source)
	if err != nil {
		return nil, fmt.Errorf("%w: build fixed supervisor v2 state: %v", ErrStoreRepairRequired, err)
	}
	encoded, err := encodeEnvelopeV2(envelope)
	if err != nil {
		return nil, err
	}
	if err := persistEnvelopeV2(ctx, path, encoded, options.Lock, options.Faults); err != nil {
		return nil, err
	}
	if err := options.Lock.CheckOfflineMigrationLock(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMigrationLockRequired, err)
	}
	return OpenFixedFileStoreV2(path)
}

func buildEnvelopeV2(source loadedV1State) (diskEnvelopeV2, error) {
	if err := validateV1State(source.State, source.Lifecycles, source.LifecycleSetDigest); err != nil {
		return diskEnvelopeV2{}, err
	}
	if err := validateV2RegistrationSetDigest(source.State); err != nil {
		return diskEnvelopeV2{}, err
	}
	indexes, err := reconstructV2Indexes(source.State, source.Lifecycles)
	if err != nil {
		return diskEnvelopeV2{}, err
	}
	descriptorDigest, err := archivestate.DigestArchiveDescriptorSet([]archivestate.ArchiveDescriptor{})
	if err != nil {
		return diskEnvelopeV2{}, err
	}
	indexView := indexes.View()
	hotCounts := countsForIndexView(indexView)
	visibleSeed := visibleV1Seed(indexView)
	hotSetDigests := archivestate.HotSetDigests{
		Registrations: source.State.RegistrationSetDigest, Approvals: source.State.ApprovalSetDigest,
		Attempts: source.State.AttemptSetDigest, Lifecycles: [32]byte(source.LifecycleSetDigest),
	}
	genesis, err := archivestate.NewMigrationGenesisCheckpoint(archivestate.MigrationGenesisCheckpointView{
		StoreFormatVersion:       archivestate.SupervisorStoreFormatV2,
		MigrationSourceVersion:   archivestate.MigrationSourceFormatV1,
		ResultSnapshotGeneration: 1, ArchiveGeneration: 1,
		DescriptorSetDigest: descriptorDigest, Indexes: indexes, HotSetDigests: hotSetDigests,
		VisibleV1EffectSeed: visibleSeed, HotCounts: hotCounts,
		InstallationID: source.State.InstallationID, SupervisorID: source.State.SupervisorID,
		EpochSequence: source.State.EpochSequence, EpochDigest: source.State.EpochDigest,
		DurableTimeHighWater: source.State.TimeHighWaterUnixSeconds,
	})
	if err != nil {
		return diskEnvelopeV2{}, err
	}
	active, err := archivestate.NewActiveStateV2(archivestate.ActiveStateV2View{
		StoreFormatVersion:     archivestate.SupervisorStoreFormatV2,
		MigrationSourceVersion: archivestate.MigrationSourceFormatV1,
		SnapshotGeneration:     1, ArchiveGeneration: 1,
		InstallationID: source.State.InstallationID, SupervisorID: source.State.SupervisorID,
		EpochSequence: source.State.EpochSequence, EpochDigest: source.State.EpochDigest,
		DurableTimeHighWater: source.State.TimeHighWaterUnixSeconds,
		Descriptors:          []archivestate.ArchiveDescriptor{}, DescriptorSetDigest: descriptorDigest,
		Indexes: indexes, IndexDigests: indexes.Digests(), CombinedIndexDigest: indexes.CombinedDigest(),
		EffectTombstoneCoverage:  archivestate.EffectTombstoneCoverageVisibleV1AndV2,
		VisibleV1EffectSeedCount: uint64(len(visibleSeed)), VisibleV1EffectSeedDigest: genesis.SeedDigest(),
		PreviousCheckpoint: archivestate.NoArchiveCheckpointReference(), CurrentCheckpoint: genesis.Reference(),
		HotCounts: hotCounts, ArchivedCounts: archivestate.ArchiveCounts{}, TotalCounts: hotCounts,
	})
	if err != nil {
		return diskEnvelopeV2{}, err
	}
	version := archivestate.SupervisorStoreFormatV2
	sourceVersion := archivestate.MigrationSourceFormatV1
	lifecycleDisks := make([]lifecycleRecordDisk, len(source.Lifecycles))
	for index, record := range source.Lifecycles {
		lifecycleDisks[index] = lifecycleRecordToDisk(record)
	}
	activeView := active.View()
	genesisView := genesis.View()
	return diskEnvelopeV2{
		StoreFormatVersion: &version, MigrationSourceVersion: &sourceVersion,
		SnapshotGeneration: activeView.SnapshotGeneration, ArchiveGeneration: activeView.ArchiveGeneration,
		State: cloneState(source.State), LifecycleSetDigest: source.LifecycleSetDigest, Lifecycles: lifecycleDisks,
		Descriptors: []archivestate.ArchiveDescriptorView{}, DescriptorSetDigest: activeView.DescriptorSetDigest,
		Indexes: archiveIndexesToDisk(indexes), IndexDigests: activeView.IndexDigests,
		CombinedIndexDigest:       activeView.CombinedIndexDigest,
		EffectTombstoneCoverage:   activeView.EffectTombstoneCoverage,
		VisibleV1EffectSeedCount:  activeView.VisibleV1EffectSeedCount,
		VisibleV1EffectSeedDigest: activeView.VisibleV1EffectSeedDigest,
		PreviousCheckpoint:        activeView.PreviousCheckpoint, CurrentCheckpoint: activeView.CurrentCheckpoint,
		MigrationGenesis: migrationGenesisDisk{
			Kind: migrationGenesisKind, StoreFormatVersion: genesisView.StoreFormatVersion,
			MigrationSourceVersion:   genesisView.MigrationSourceVersion,
			ResultSnapshotGeneration: genesisView.ResultSnapshotGeneration,
			ArchiveGeneration:        genesisView.ArchiveGeneration, DescriptorSetDigest: genesisView.DescriptorSetDigest,
			IndexDigests: indexes.Digests(), CombinedIndexDigest: indexes.CombinedDigest(),
			HotSetDigests:             genesisView.HotSetDigests,
			VisibleV1EffectSeed:       append([]lifecyclestate.EffectID(nil), genesisView.VisibleV1EffectSeed...),
			VisibleV1EffectSeedDigest: genesis.SeedDigest(), HotCounts: genesisView.HotCounts,
			InstallationID: genesisView.InstallationID, SupervisorID: genesisView.SupervisorID,
			EpochSequence: genesisView.EpochSequence, EpochDigest: genesisView.EpochDigest,
			DurableTimeHighWater: genesisView.DurableTimeHighWater, Digest: genesis.Digest(),
		},
		HotCounts: activeView.HotCounts, ArchivedCounts: activeView.ArchivedCounts,
		TotalCounts: activeView.TotalCounts,
	}, nil
}

func encodeEnvelopeV2(envelope diskEnvelopeV2) ([]byte, error) {
	encoded, err := json.Marshal(envelope)
	if err != nil {
		return nil, errors.New("encode fixed supervisor v2 migration")
	}
	encoded = append(encoded, '\n')
	if int64(len(encoded)) > archivestate.MaxSupervisorStateV2Bytes {
		return nil, errors.New("fixed supervisor v2 store exceeds byte capacity")
	}
	return encoded, nil
}

func persistEnvelopeV2(
	ctx context.Context,
	path string,
	encoded []byte,
	lock OfflineMigrationLock,
	faults MigrationFaultInjector,
) error {
	parent := filepath.Dir(path)
	temporary, err := os.CreateTemp(parent, ".capsule-supervisor-v2-migration-*.tmp")
	if err != nil {
		return errors.New("create fixed supervisor v2 migration")
	}
	temporaryPath := temporary.Name()
	committed := false
	defer func() {
		_ = temporary.Close()
		if !committed {
			_ = os.Remove(temporaryPath) //nolint:gosec // G703: exact path returned by CreateTemp in the validated store directory.
		}
	}()
	if err := failV2Migration(faults, FaultV2MigrationAfterTempCreate, false); err != nil {
		return err
	}
	if err := temporary.Chmod(0o600); err != nil {
		return errors.New("protect fixed supervisor v2 migration")
	}
	if err := failV2Migration(faults, FaultV2MigrationAfterTempMode, false); err != nil {
		return err
	}
	if _, err := temporary.Write(encoded); err != nil {
		return errors.New("write fixed supervisor v2 migration")
	}
	if err := failV2Migration(faults, FaultV2MigrationAfterTempWrite, false); err != nil {
		return err
	}
	if err := temporary.Sync(); err != nil {
		return errors.New("sync fixed supervisor v2 migration")
	}
	if err := failV2Migration(faults, FaultV2MigrationAfterTempSync, false); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return errors.New("close fixed supervisor v2 migration")
	}
	if err := failV2Migration(faults, FaultV2MigrationAfterTempClose, false); err != nil {
		return err
	}
	if err := failV2Migration(faults, FaultV2MigrationBeforeRename, false); err != nil {
		return err
	}
	if err := lock.CheckOfflineMigrationLock(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrMigrationLockRequired, err)
	}
	if err := os.Rename(temporaryPath, path); err != nil { //nolint:gosec // G703: exact CreateTemp path and installation-configured store path.
		return errors.New("commit fixed supervisor v2 migration")
	}
	committed = true
	if err := failV2Migration(faults, FaultV2MigrationAfterRename, true); err != nil {
		return err
	}
	directory, err := os.Open(parent) //nolint:gosec // G304: trusted installation-configured store parent
	if err != nil {
		return fmt.Errorf("%w: open v2 migration directory", ErrMigrationOutcomeIndeterminate)
	}
	syncErr := directory.Sync()
	closeErr := directory.Close()
	if syncErr != nil || closeErr != nil {
		return fmt.Errorf("%w: sync v2 migration directory", ErrMigrationOutcomeIndeterminate)
	}
	if err := failV2Migration(faults, FaultV2MigrationAfterDirSync, true); err != nil {
		return err
	}
	return nil
}

func failV2Migration(faults MigrationFaultInjector, point MigrationFault, indeterminate bool) error {
	if faults == nil {
		return nil
	}
	if err := faults.FailMigrationAt(point); err != nil {
		if indeterminate {
			return fmt.Errorf("%w: injected at %s", ErrMigrationOutcomeIndeterminate, point)
		}
		return fmt.Errorf("fixed supervisor v2 migration confirmed failure at %s", point)
	}
	return nil
}

func loadV2State(path string) (loadedV2State, error) {
	data, err := readBoundedStoreFile(path, archivestate.MaxSupervisorStateV2Bytes)
	if err != nil {
		return loadedV2State{}, err
	}
	var envelope diskEnvelopeV2
	if err := decodeOneClosedJSON(data, &envelope); err != nil {
		return loadedV2State{}, errors.New("decode fixed supervisor v2 store")
	}
	if envelope.StoreFormatVersion == nil || *envelope.StoreFormatVersion != archivestate.SupervisorStoreFormatV2 {
		return loadedV2State{}, errors.New("unsupported fixed supervisor v2 store version")
	}
	if envelope.MigrationSourceVersion == nil || *envelope.MigrationSourceVersion != archivestate.MigrationSourceFormatV1 {
		return loadedV2State{}, errors.New("unsupported fixed supervisor v2 migration source version")
	}
	if envelope.Lifecycles == nil || envelope.Descriptors == nil {
		return loadedV2State{}, errors.New("fixed supervisor v2 collection is missing")
	}
	if len(envelope.Descriptors) > archivestate.MaxReferencedArchiveSegments {
		return loadedV2State{}, errors.New("fixed supervisor v2 store exceeds archive descriptor capacity")
	}
	if len(envelope.Lifecycles) > MaxRetainedLifecycleRecords {
		return loadedV2State{}, errors.New("fixed lifecycle state exceeds retained capacity")
	}
	if err := validateStateWithAttemptCapacity(envelope.State, false); err != nil {
		return loadedV2State{}, err
	}
	if err := validateV1AuthorityDigests(envelope.State); err != nil {
		return loadedV2State{}, err
	}
	records := make([]lifecyclestate.Record, len(envelope.Lifecycles))
	for index, disk := range envelope.Lifecycles {
		record, restoreErr := restoreLifecycleRecord(disk, envelope.State.TimeHighWaterUnixSeconds)
		if restoreErr != nil {
			return loadedV2State{}, fmt.Errorf("invalid v2 lifecycle record %d", index)
		}
		records[index] = record
	}
	if err := validateV1State(envelope.State, records, envelope.LifecycleSetDigest); err != nil {
		return loadedV2State{}, err
	}
	if err := validateV2RegistrationSetDigest(envelope.State); err != nil {
		return loadedV2State{}, err
	}
	storedIndexes, err := archiveIndexesFromDisk(envelope.Indexes)
	if err != nil {
		return loadedV2State{}, err
	}
	segments, orphans, err := loadArchiveSegmentsForEnvelope(path, envelope)
	if err != nil {
		return loadedV2State{}, err
	}
	reconstructed, err := reconstructV2IndexesForWorld(envelope.State, records, segments)
	if err != nil {
		return loadedV2State{}, err
	}
	if !reflect.DeepEqual(storedIndexes.View(), reconstructed.View()) {
		return loadedV2State{}, errors.New("fixed supervisor v2 reconstructed index mismatch")
	}
	if storedIndexes.Digests() != envelope.IndexDigests ||
		storedIndexes.CombinedDigest() != envelope.CombinedIndexDigest {
		return loadedV2State{}, errors.New("fixed supervisor v2 index digest mismatch")
	}
	genesisIndexes, err := reconstructMigrationGenesisIndexes(envelope.State, records, segments)
	if err != nil {
		return loadedV2State{}, err
	}
	genesisHotDigests, err := reconstructMigrationGenesisHotSetDigests(envelope.State, records, segments)
	if err != nil {
		return loadedV2State{}, err
	}
	genesis, err := validateMigrationGenesisDisk(envelope, genesisIndexes, genesisHotDigests)
	if err != nil {
		return loadedV2State{}, err
	}
	descriptors := make([]archivestate.ArchiveDescriptor, len(envelope.Descriptors))
	for index, view := range envelope.Descriptors {
		descriptor, descriptorErr := archivestate.NewArchiveDescriptor(view)
		if descriptorErr != nil {
			return loadedV2State{}, descriptorErr
		}
		descriptors[index] = descriptor
	}
	active, err := archivestate.NewActiveStateV2(archivestate.ActiveStateV2View{
		StoreFormatVersion:     *envelope.StoreFormatVersion,
		MigrationSourceVersion: *envelope.MigrationSourceVersion,
		SnapshotGeneration:     envelope.SnapshotGeneration, ArchiveGeneration: envelope.ArchiveGeneration,
		InstallationID: envelope.State.InstallationID, SupervisorID: envelope.State.SupervisorID,
		EpochSequence: envelope.State.EpochSequence, EpochDigest: envelope.State.EpochDigest,
		DurableTimeHighWater: envelope.State.TimeHighWaterUnixSeconds,
		Descriptors:          descriptors, DescriptorSetDigest: envelope.DescriptorSetDigest,
		Indexes: storedIndexes, IndexDigests: envelope.IndexDigests,
		CombinedIndexDigest:       envelope.CombinedIndexDigest,
		EffectTombstoneCoverage:   envelope.EffectTombstoneCoverage,
		VisibleV1EffectSeedCount:  envelope.VisibleV1EffectSeedCount,
		VisibleV1EffectSeedDigest: envelope.VisibleV1EffectSeedDigest,
		PreviousCheckpoint:        envelope.PreviousCheckpoint, CurrentCheckpoint: envelope.CurrentCheckpoint,
		HotCounts: envelope.HotCounts, ArchivedCounts: envelope.ArchivedCounts, TotalCounts: envelope.TotalCounts,
	})
	if err != nil {
		return loadedV2State{}, err
	}
	switch envelope.CurrentCheckpoint.Kind {
	case archivestate.ArchiveCheckpointMigrationGenesis:
		if envelope.SnapshotGeneration != 1 || envelope.ArchiveGeneration != 1 ||
			envelope.ArchivedCounts != (archivestate.ArchiveCounts{}) || len(segments) != 0 ||
			envelope.CurrentCheckpoint != genesis.Reference() ||
			envelope.PreviousCheckpoint != archivestate.NoArchiveCheckpointReference() ||
			envelope.ActivationCheckpoint != nil {
			return loadedV2State{}, errors.New("fixed supervisor v2 genesis world mismatch")
		}
	case archivestate.ArchiveCheckpointActivation:
		if len(segments) != 1 || envelope.SnapshotGeneration != 2 || envelope.ArchiveGeneration != 2 ||
			envelope.PreviousCheckpoint != genesis.Reference() || envelope.ActivationCheckpoint == nil {
			return loadedV2State{}, errors.New("fixed supervisor F3 v2 activation world mismatch")
		}
		if _, checkpointErr := validateActivationCheckpointDisk(envelope, storedIndexes, segments[0]); checkpointErr != nil {
			return loadedV2State{}, checkpointErr
		}
	default:
		return loadedV2State{}, errors.New("fixed supervisor v2 checkpoint kind is unsupported")
	}
	return loadedV2State{
		Envelope: envelope, State: cloneState(envelope.State),
		Lifecycles: cloneLifecycleRecords(records, envelope.State.TimeHighWaterUnixSeconds),
		Active:     active, Genesis: genesis, Segments: cloneLoadedArchiveSegments(segments), Orphans: orphans,
	}, nil
}

func validateMigrationGenesisDisk(
	envelope diskEnvelopeV2,
	indexes archivestate.ArchiveIndexes,
	hotSetDigests archivestate.HotSetDigests,
) (archivestate.MigrationGenesisCheckpoint, error) {
	disk := envelope.MigrationGenesis
	emptyDescriptors, emptyDescriptorErr := archivestate.DigestArchiveDescriptorSet([]archivestate.ArchiveDescriptor{})
	if disk.Kind != migrationGenesisKind || disk.IndexDigests != indexes.Digests() ||
		disk.CombinedIndexDigest != indexes.CombinedDigest() ||
		disk.StoreFormatVersion != *envelope.StoreFormatVersion ||
		disk.MigrationSourceVersion != *envelope.MigrationSourceVersion ||
		disk.ResultSnapshotGeneration != 1 || disk.ArchiveGeneration != 1 || emptyDescriptorErr != nil ||
		disk.DescriptorSetDigest != emptyDescriptors || disk.HotCounts != countsForIndexView(indexes.View()) ||
		disk.InstallationID != envelope.State.InstallationID ||
		disk.SupervisorID != envelope.State.SupervisorID || disk.EpochSequence != envelope.State.EpochSequence ||
		disk.EpochDigest != envelope.State.EpochDigest || disk.DurableTimeHighWater != envelope.State.TimeHighWaterUnixSeconds {
		return archivestate.MigrationGenesisCheckpoint{}, errors.New("fixed supervisor v2 migration genesis cross-link mismatch")
	}
	if disk.HotSetDigests != hotSetDigests {
		return archivestate.MigrationGenesisCheckpoint{}, errors.New("fixed supervisor v2 migration genesis hot-set mismatch")
	}
	seed := make([]lifecyclestate.EffectID, len(disk.VisibleV1EffectSeed))
	copy(seed, disk.VisibleV1EffectSeed)
	genesis, err := archivestate.NewMigrationGenesisCheckpoint(archivestate.MigrationGenesisCheckpointView{
		StoreFormatVersion: disk.StoreFormatVersion, MigrationSourceVersion: disk.MigrationSourceVersion,
		ResultSnapshotGeneration: disk.ResultSnapshotGeneration, ArchiveGeneration: disk.ArchiveGeneration,
		DescriptorSetDigest: disk.DescriptorSetDigest, Indexes: indexes, HotSetDigests: disk.HotSetDigests,
		VisibleV1EffectSeed: seed,
		HotCounts:           disk.HotCounts, InstallationID: disk.InstallationID, SupervisorID: disk.SupervisorID,
		EpochSequence: disk.EpochSequence, EpochDigest: disk.EpochDigest,
		DurableTimeHighWater: disk.DurableTimeHighWater,
	})
	if err != nil {
		return archivestate.MigrationGenesisCheckpoint{}, err
	}
	if genesis.SeedDigest() != disk.VisibleV1EffectSeedDigest || genesis.Digest() != disk.Digest ||
		uint64(len(disk.VisibleV1EffectSeed)) != envelope.VisibleV1EffectSeedCount ||
		disk.VisibleV1EffectSeedDigest != envelope.VisibleV1EffectSeedDigest {
		return archivestate.MigrationGenesisCheckpoint{}, errors.New("fixed supervisor v2 migration genesis digest mismatch")
	}
	return genesis, nil
}

func reconstructV2Indexes(
	state installationState,
	records []lifecyclestate.Record,
) (archivestate.ArchiveIndexes, error) {
	view := archivestate.ArchiveIndexesView{
		Scope:          archivestate.ArchiveIndexScopeRetainedGlobal,
		Registrations:  make([]archivestate.RegistrationIndexEntry, 0, len(state.Registrations)),
		Approvals:      make([]archivestate.ApprovalIndexEntry, 0, len(state.Approvals)),
		Attempts:       make([]archivestate.AttemptIndexEntry, 0, len(state.Attempts)),
		Nonces:         make([]archivestate.NonceIndexEntry, 0, len(state.Approvals)),
		Effects:        make([]archivestate.EffectIndexEntry, 0, len(records)),
		Instances:      make([]archivestate.InstanceIndexEntry, 0, len(records)),
		ApprovalReplay: make([]archivestate.ApprovalReplayIndexEntry, 0, len(state.Approvals)),
		AttemptReplay:  make([]archivestate.AttemptReplayIndexEntry, 0, len(state.Attempts)),
	}

	registrations := append([]registrationEntry(nil), state.Registrations...)
	sort.Slice(registrations, func(left, right int) bool {
		return bytes.Compare(registrations[left].Index.RegistrationID[:], registrations[right].Index.RegistrationID[:]) < 0
	})
	for index, entry := range registrations {
		location, err := archivestate.NewHotRecordLocation(archivestate.RecordRegistration, v0candidate.PositiveUInt53(index+1))
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		digest, err := digestV2Record(archivestate.RecordRegistration, entry)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		view.Registrations = append(view.Registrations, archivestate.RegistrationIndexEntry{
			RegistrationID: entry.Index.RegistrationID, RegistrationSequence: entry.Index.RegistrationSequence,
			PlanDigest: entry.Record.RecomputedPlanDigest, ExpiresAt: entry.Index.ExpiresAt,
			Location: location, FullRecordDigest: digest,
		})
	}

	approvals := make([]approvalattempt.ApprovalRecord, len(state.Approvals))
	for index, record := range state.Approvals {
		approvals[index] = approvalattempt.CloneApprovalRecord(record)
	}
	sort.Slice(approvals, func(left, right int) bool {
		return bytes.Compare(approvals[left].ApprovalID[:], approvals[right].ApprovalID[:]) < 0
	})
	for index, record := range approvals {
		location, err := archivestate.NewHotRecordLocation(archivestate.RecordApproval, v0candidate.PositiveUInt53(index+1))
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		digest, err := digestV2Record(archivestate.RecordApproval, record)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		exactPayloadDigest, err := archivestate.DigestRecord(archivestate.RecordApproval, record.ExactPayloadBytes)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		view.Approvals = append(view.Approvals, archivestate.ApprovalIndexEntry{
			ApprovalID: record.ApprovalID, RegistrationID: record.RegistrationID,
			PayloadDigest: record.PayloadDigest, AuthorizationIdentity: record.AuthorizationIdentity,
			AttemptNonce: record.AttemptNonce, State: record.State, ConsumedAttemptID: record.ConsumedAttemptID,
			ExpiresAt: record.ExpiresAt, Location: location, FullRecordDigest: digest,
		})
		view.Nonces = append(view.Nonces, archivestate.NonceIndexEntry{
			AttemptNonce: record.AttemptNonce, PayloadDigest: record.PayloadDigest,
			ApprovalID: record.ApprovalID, Location: location,
		})
		view.ApprovalReplay = append(view.ApprovalReplay, archivestate.ApprovalReplayIndexEntry{
			PayloadDigest: record.PayloadDigest, ExactPayloadDigest: exactPayloadDigest,
			AuthorizationIdentity: record.AuthorizationIdentity, ApprovalID: record.ApprovalID,
			State: record.State, Location: location,
		})
	}

	lifecycleByAttempt := make(map[approvalattempt.AttemptID]struct {
		record   lifecyclestate.Record
		location archivestate.RecordLocation
		digest   archivestate.ArchiveRecordDigest
	}, len(records))
	for index, record := range records {
		attemptID := record.Bindings().View().AttemptID
		location, err := archivestate.NewHotRecordLocation(archivestate.RecordLifecycle, v0candidate.PositiveUInt53(index+1))
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		digest, err := digestV2Record(archivestate.RecordLifecycle, lifecycleRecordToDisk(record))
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		lifecycleByAttempt[attemptID] = struct {
			record   lifecyclestate.Record
			location archivestate.RecordLocation
			digest   archivestate.ArchiveRecordDigest
		}{record: record, location: location, digest: digest}
	}

	attempts := append([]approvalattempt.ExecutionAttempt(nil), state.Attempts...)
	sort.Slice(attempts, func(left, right int) bool {
		return bytes.Compare(attempts[left].AttemptID[:], attempts[right].AttemptID[:]) < 0
	})
	for index, attempt := range attempts {
		location, err := archivestate.NewHotRecordLocation(archivestate.RecordAttempt, v0candidate.PositiveUInt53(index+1))
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		digest, err := digestV2Record(archivestate.RecordAttempt, attempt)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		lifecycle := archivestate.NoAttemptLifecycle()
		if joined, present := lifecycleByAttempt[attempt.AttemptID]; present {
			lifecycle, err = archivestate.NewPresentAttemptLifecycle(joined.record.View().State, joined.location, joined.digest)
			if err != nil {
				return archivestate.ArchiveIndexes{}, err
			}
			recordView := joined.record.View()
			if !recordView.EffectID.IsZero() {
				view.Effects = append(view.Effects, archivestate.EffectIndexEntry{
					EffectID: recordView.EffectID, AttemptID: attempt.AttemptID,
					OperationSequence: recordView.OperationSequence, Operation: recordView.Operation,
					IssuanceSnapshotGeneration: recordView.SnapshotGeneration,
					VisibleV1Seed:              true, AttemptLocation: location,
				})
			}
			if recordView.Instance.Present() {
				view.Instances = append(view.Instances, archivestate.InstanceIndexEntry{
					InstanceDigest: recordView.Instance.Digest(), AttemptID: attempt.AttemptID,
					Location: joined.location,
				})
			}
		}
		view.Attempts = append(view.Attempts, archivestate.AttemptIndexEntry{
			AttemptID: attempt.AttemptID, ApprovalID: attempt.ApprovalID,
			RegistrationID: attempt.RegistrationID, CreatedAt: attempt.CreatedAt,
			Lifecycle: lifecycle, Location: location, FullRecordDigest: digest,
		})
		view.AttemptReplay = append(view.AttemptReplay, archivestate.AttemptReplayIndexEntry{
			RegistrationID: attempt.RegistrationID, ApprovalID: attempt.ApprovalID,
			AttemptID: attempt.AttemptID, State: attempt.State, Location: location,
		})
	}

	sort.Slice(view.Nonces, func(left, right int) bool {
		return bytes.Compare(view.Nonces[left].AttemptNonce[:], view.Nonces[right].AttemptNonce[:]) < 0
	})
	sort.Slice(view.Effects, func(left, right int) bool {
		return bytes.Compare(view.Effects[left].EffectID[:], view.Effects[right].EffectID[:]) < 0
	})
	sort.Slice(view.Instances, func(left, right int) bool {
		return bytes.Compare(view.Instances[left].InstanceDigest[:], view.Instances[right].InstanceDigest[:]) < 0
	})
	sort.Slice(view.ApprovalReplay, func(left, right int) bool {
		comparison := bytes.Compare(view.ApprovalReplay[left].PayloadDigest[:], view.ApprovalReplay[right].PayloadDigest[:])
		if comparison != 0 {
			return comparison < 0
		}
		return bytes.Compare(view.ApprovalReplay[left].AuthorizationIdentity[:], view.ApprovalReplay[right].AuthorizationIdentity[:]) < 0
	})
	sort.Slice(view.AttemptReplay, func(left, right int) bool {
		comparison := bytes.Compare(view.AttemptReplay[left].RegistrationID[:], view.AttemptReplay[right].RegistrationID[:])
		if comparison != 0 {
			return comparison < 0
		}
		return bytes.Compare(view.AttemptReplay[left].ApprovalID[:], view.AttemptReplay[right].ApprovalID[:]) < 0
	})
	return archivestate.NewArchiveIndexes(view)
}

func digestV2Record(kind archivestate.RecordKind, record any) (archivestate.ArchiveRecordDigest, error) {
	exact, err := json.Marshal(record)
	if err != nil {
		return archivestate.ArchiveRecordDigest{}, errors.New("encode fixed supervisor v2 record")
	}
	return archivestate.DigestRecord(kind, exact)
}

func validateV2RegistrationSetDigest(state installationState) error {
	digest := emptyRegistrationSetDigest()
	for _, entry := range state.Registrations {
		var err error
		digest, err = appendRegistrationSetDigest(digest, entry.Record)
		if err != nil {
			return errors.New("fixed supervisor v2 registration set digest computation failed")
		}
	}
	if digest != state.RegistrationSetDigest {
		return errors.New("fixed supervisor v2 registration set digest mismatch")
	}
	return nil
}

func countsForIndexView(view archivestate.ArchiveIndexesView) archivestate.ArchiveCounts {
	counts := archivestate.ArchiveCounts{
		Registrations: uint64(len(view.Registrations)), Approvals: uint64(len(view.Approvals)),
		Attempts: uint64(len(view.Attempts)), Nonces: uint64(len(view.Nonces)),
		Effects: uint64(len(view.Effects)), Instances: uint64(len(view.Instances)),
		ApprovalReplay: uint64(len(view.ApprovalReplay)), AttemptReplay: uint64(len(view.AttemptReplay)),
	}
	for _, attempt := range view.Attempts {
		if attempt.Lifecycle.View().Presence == archivestate.AttemptLifecyclePresent {
			counts.Lifecycles++
		}
	}
	return counts
}

func visibleV1Seed(view archivestate.ArchiveIndexesView) []lifecyclestate.EffectID {
	result := make([]lifecyclestate.EffectID, 0, len(view.Effects))
	for _, effect := range view.Effects {
		if effect.VisibleV1Seed {
			result = append(result, effect.EffectID)
		}
	}
	return result
}

func archiveIndexesToDisk(indexes archivestate.ArchiveIndexes) archiveIndexesDisk {
	view := indexes.View()
	disk := archiveIndexesDisk{
		Scope:          view.Scope,
		Registrations:  make([]registrationIndexEntryDisk, len(view.Registrations)),
		Approvals:      make([]approvalIndexEntryDisk, len(view.Approvals)),
		Attempts:       make([]attemptIndexEntryDisk, len(view.Attempts)),
		Nonces:         make([]nonceIndexEntryDisk, len(view.Nonces)),
		Effects:        make([]effectIndexEntryDisk, len(view.Effects)),
		Instances:      make([]instanceIndexEntryDisk, len(view.Instances)),
		ApprovalReplay: make([]approvalReplayIndexEntryDisk, len(view.ApprovalReplay)),
		AttemptReplay:  make([]attemptReplayIndexEntryDisk, len(view.AttemptReplay)),
	}
	for index, entry := range view.Registrations {
		disk.Registrations[index] = registrationIndexEntryDisk{
			RegistrationID: entry.RegistrationID, RegistrationSequence: entry.RegistrationSequence,
			PlanDigest: entry.PlanDigest, ExpiresAt: entry.ExpiresAt,
			Location: recordLocationToDisk(entry.Location), FullRecordDigest: entry.FullRecordDigest,
		}
	}
	for index, entry := range view.Approvals {
		disk.Approvals[index] = approvalIndexEntryDisk{
			ApprovalID: entry.ApprovalID, RegistrationID: entry.RegistrationID,
			PayloadDigest: entry.PayloadDigest, AuthorizationIdentity: entry.AuthorizationIdentity,
			AttemptNonce: entry.AttemptNonce, State: entry.State, ConsumedAttemptID: entry.ConsumedAttemptID,
			ExpiresAt: entry.ExpiresAt, Location: recordLocationToDisk(entry.Location),
			FullRecordDigest: entry.FullRecordDigest,
		}
	}
	for index, entry := range view.Attempts {
		disk.Attempts[index] = attemptIndexEntryDisk{
			AttemptID: entry.AttemptID, ApprovalID: entry.ApprovalID,
			RegistrationID: entry.RegistrationID, CreatedAt: entry.CreatedAt,
			Lifecycle: attemptLifecycleToDisk(entry.Lifecycle), Location: recordLocationToDisk(entry.Location),
			FullRecordDigest: entry.FullRecordDigest,
		}
	}
	for index, entry := range view.Nonces {
		disk.Nonces[index] = nonceIndexEntryDisk{
			AttemptNonce: entry.AttemptNonce, PayloadDigest: entry.PayloadDigest,
			ApprovalID: entry.ApprovalID, Location: recordLocationToDisk(entry.Location),
		}
	}
	for index, entry := range view.Effects {
		disk.Effects[index] = effectIndexEntryDisk{
			EffectID: entry.EffectID, AttemptID: entry.AttemptID, OperationSequence: entry.OperationSequence,
			Operation: entry.Operation, IssuanceSnapshotGeneration: entry.IssuanceSnapshotGeneration,
			VisibleV1Seed: entry.VisibleV1Seed, AttemptLocation: recordLocationToDisk(entry.AttemptLocation),
		}
	}
	for index, entry := range view.Instances {
		disk.Instances[index] = instanceIndexEntryDisk{
			InstanceDigest: entry.InstanceDigest, AttemptID: entry.AttemptID,
			Location: recordLocationToDisk(entry.Location),
		}
	}
	for index, entry := range view.ApprovalReplay {
		disk.ApprovalReplay[index] = approvalReplayIndexEntryDisk{
			PayloadDigest: entry.PayloadDigest, ExactPayloadDigest: entry.ExactPayloadDigest,
			AuthorizationIdentity: entry.AuthorizationIdentity, ApprovalID: entry.ApprovalID,
			State: entry.State, Location: recordLocationToDisk(entry.Location),
		}
	}
	for index, entry := range view.AttemptReplay {
		disk.AttemptReplay[index] = attemptReplayIndexEntryDisk{
			RegistrationID: entry.RegistrationID, ApprovalID: entry.ApprovalID,
			AttemptID: entry.AttemptID, State: entry.State, Location: recordLocationToDisk(entry.Location),
		}
	}
	return disk
}

func archiveIndexesFromDisk(disk archiveIndexesDisk) (archivestate.ArchiveIndexes, error) {
	if disk.Registrations == nil || disk.Approvals == nil || disk.Attempts == nil || disk.Nonces == nil ||
		disk.Effects == nil || disk.Instances == nil || disk.ApprovalReplay == nil || disk.AttemptReplay == nil {
		return archivestate.ArchiveIndexes{}, errors.New("fixed supervisor v2 index collection is missing")
	}
	view := archivestate.ArchiveIndexesView{
		Scope:          disk.Scope,
		Registrations:  make([]archivestate.RegistrationIndexEntry, len(disk.Registrations)),
		Approvals:      make([]archivestate.ApprovalIndexEntry, len(disk.Approvals)),
		Attempts:       make([]archivestate.AttemptIndexEntry, len(disk.Attempts)),
		Nonces:         make([]archivestate.NonceIndexEntry, len(disk.Nonces)),
		Effects:        make([]archivestate.EffectIndexEntry, len(disk.Effects)),
		Instances:      make([]archivestate.InstanceIndexEntry, len(disk.Instances)),
		ApprovalReplay: make([]archivestate.ApprovalReplayIndexEntry, len(disk.ApprovalReplay)),
		AttemptReplay:  make([]archivestate.AttemptReplayIndexEntry, len(disk.AttemptReplay)),
	}
	for index, entry := range disk.Registrations {
		location, err := recordLocationFromDisk(entry.Location)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		view.Registrations[index] = archivestate.RegistrationIndexEntry{
			RegistrationID: entry.RegistrationID, RegistrationSequence: entry.RegistrationSequence,
			PlanDigest: entry.PlanDigest, ExpiresAt: entry.ExpiresAt,
			Location: location, FullRecordDigest: entry.FullRecordDigest,
		}
	}
	for index, entry := range disk.Approvals {
		location, err := recordLocationFromDisk(entry.Location)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		view.Approvals[index] = archivestate.ApprovalIndexEntry{
			ApprovalID: entry.ApprovalID, RegistrationID: entry.RegistrationID,
			PayloadDigest: entry.PayloadDigest, AuthorizationIdentity: entry.AuthorizationIdentity,
			AttemptNonce: entry.AttemptNonce, State: entry.State, ConsumedAttemptID: entry.ConsumedAttemptID,
			ExpiresAt: entry.ExpiresAt, Location: location, FullRecordDigest: entry.FullRecordDigest,
		}
	}
	for index, entry := range disk.Attempts {
		location, err := recordLocationFromDisk(entry.Location)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		lifecycle, err := attemptLifecycleFromDisk(entry.Lifecycle)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		view.Attempts[index] = archivestate.AttemptIndexEntry{
			AttemptID: entry.AttemptID, ApprovalID: entry.ApprovalID,
			RegistrationID: entry.RegistrationID, CreatedAt: entry.CreatedAt,
			Lifecycle: lifecycle, Location: location, FullRecordDigest: entry.FullRecordDigest,
		}
	}
	for index, entry := range disk.Nonces {
		location, err := recordLocationFromDisk(entry.Location)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		view.Nonces[index] = archivestate.NonceIndexEntry{
			AttemptNonce: entry.AttemptNonce, PayloadDigest: entry.PayloadDigest,
			ApprovalID: entry.ApprovalID, Location: location,
		}
	}
	for index, entry := range disk.Effects {
		location, err := recordLocationFromDisk(entry.AttemptLocation)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		view.Effects[index] = archivestate.EffectIndexEntry{
			EffectID: entry.EffectID, AttemptID: entry.AttemptID, OperationSequence: entry.OperationSequence,
			Operation: entry.Operation, IssuanceSnapshotGeneration: entry.IssuanceSnapshotGeneration,
			VisibleV1Seed: entry.VisibleV1Seed, AttemptLocation: location,
		}
	}
	for index, entry := range disk.Instances {
		location, err := recordLocationFromDisk(entry.Location)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		view.Instances[index] = archivestate.InstanceIndexEntry{
			InstanceDigest: entry.InstanceDigest, AttemptID: entry.AttemptID, Location: location,
		}
	}
	for index, entry := range disk.ApprovalReplay {
		location, err := recordLocationFromDisk(entry.Location)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		view.ApprovalReplay[index] = archivestate.ApprovalReplayIndexEntry{
			PayloadDigest: entry.PayloadDigest, ExactPayloadDigest: entry.ExactPayloadDigest,
			AuthorizationIdentity: entry.AuthorizationIdentity, ApprovalID: entry.ApprovalID,
			State: entry.State, Location: location,
		}
	}
	for index, entry := range disk.AttemptReplay {
		location, err := recordLocationFromDisk(entry.Location)
		if err != nil {
			return archivestate.ArchiveIndexes{}, err
		}
		view.AttemptReplay[index] = archivestate.AttemptReplayIndexEntry{
			RegistrationID: entry.RegistrationID, ApprovalID: entry.ApprovalID,
			AttemptID: entry.AttemptID, State: entry.State, Location: location,
		}
	}
	return archivestate.NewArchiveIndexes(view)
}

func recordLocationToDisk(location archivestate.RecordLocation) recordLocationDisk {
	view := location.View()
	return recordLocationDisk{
		Kind: view.Kind, RecordKind: view.RecordKind,
		HotRecordOrdinal: view.HotRecordOrdinal, Archive: view.Archive,
	}
}

func recordLocationFromDisk(disk recordLocationDisk) (archivestate.RecordLocation, error) {
	switch disk.Kind {
	case archivestate.RecordLocationHot:
		if disk.Archive != (archivestate.ArchiveLocation{}) {
			return archivestate.RecordLocation{}, errors.New("fixed supervisor v2 hot location carries archive coordinates")
		}
		return archivestate.NewHotRecordLocation(disk.RecordKind, disk.HotRecordOrdinal)
	case archivestate.RecordLocationArchive:
		if disk.HotRecordOrdinal != 0 {
			return archivestate.RecordLocation{}, errors.New("fixed supervisor v2 archive location carries hot ordinal")
		}
		return archivestate.NewArchiveRecordLocation(disk.RecordKind, disk.Archive)
	default:
		return archivestate.RecordLocation{}, errors.New("fixed supervisor v2 location has unknown kind")
	}
}

func attemptLifecycleToDisk(lifecycle archivestate.AttemptLifecycle) attemptLifecycleDisk {
	view := lifecycle.View()
	disk := attemptLifecycleDisk{
		Presence: view.Presence, State: view.State, FullRecordDigest: view.FullRecordDigest,
	}
	if view.Presence == archivestate.AttemptLifecyclePresent {
		disk.Location = recordLocationToDisk(view.Location)
	}
	return disk
}

func attemptLifecycleFromDisk(disk attemptLifecycleDisk) (archivestate.AttemptLifecycle, error) {
	switch disk.Presence {
	case archivestate.AttemptLifecycleAbsent:
		if disk.State != "" || disk.Location != (recordLocationDisk{}) ||
			disk.FullRecordDigest != (archivestate.ArchiveRecordDigest{}) {
			return archivestate.AttemptLifecycle{}, errors.New("fixed supervisor v2 absent lifecycle carries state")
		}
		return archivestate.NoAttemptLifecycle(), nil
	case archivestate.AttemptLifecyclePresent:
		location, err := recordLocationFromDisk(disk.Location)
		if err != nil {
			return archivestate.AttemptLifecycle{}, err
		}
		return archivestate.NewPresentAttemptLifecycle(disk.State, location, disk.FullRecordDigest)
	default:
		return archivestate.AttemptLifecycle{}, errors.New("fixed supervisor v2 lifecycle presence is unknown")
	}
}
