package registrationstate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

var _ StateStore = (*FixedFileStoreV2)(nil)
var _ DurableLifecycleStore = (*FixedFileStoreV2)(nil)

// snapshot exposes only mutable hot authority state to the existing internal
// authority components. Archived terminal cohorts remain reachable solely
// through the retained lookup/replay ports and never re-enter active capacity.
func (store *FixedFileStoreV2) snapshot(ctx context.Context) (installationState, error) {
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return installationState{}, err
	}
	return cloneState(loaded.State), nil
}

func (store *FixedFileStoreV2) recoveryFenced() bool {
	return store.LifecycleRecoveryFenced()
}

// InjectFailure installs a one-shot authority-transaction fault for F4B.
func (store *FixedFileStoreV2) InjectFailure(point StoreFault, err error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if err == nil {
		err = errors.New("injected fixed-store v2 failure")
	}
	store.authorityFaults[point] = err
}

func (store *FixedFileStoreV2) persistTimeHighWater(ctx context.Context, observation v0candidate.UInt53) error {
	return store.commitV2Authority(ctx, FaultTimeHighWaterWrite, FaultTimeHighWaterIndeterminatePreState, FaultTimeHighWaterCommitIndeterminate,
		func(state *installationState) error {
			if observation <= state.TimeHighWaterUnixSeconds {
				return errNoStoreMutation
			}
			state.TimeHighWaterUnixSeconds = observation
			return nil
		})
}

func (store *FixedFileStoreV2) update(ctx context.Context, update func(*installationState) error) error {
	return store.commitV2Authority(ctx, FaultRegistrationCommitAbort, "", "", update)
}

func (store *FixedFileStoreV2) commitApproval(ctx context.Context, update func(*installationState) error) error {
	return store.commitV2Authority(ctx, FaultApprovalCommitAbort, FaultApprovalCommitIndeterminatePreState, FaultApprovalCommitIndeterminate, update)
}

func (store *FixedFileStoreV2) commitAttempt(ctx context.Context, update func(*installationState) error) error {
	return store.commitV2Authority(ctx, FaultAttemptCommitAbort, FaultAttemptCommitIndeterminatePreState, FaultAttemptCommitIndeterminate, update)
}

func (store *FixedFileStoreV2) commitV2Authority(
	ctx context.Context,
	abortPoint StoreFault,
	preIndeterminatePoint StoreFault,
	postRenamePoint StoreFault,
	mutation func(*installationState) error,
) error {
	if mutation == nil {
		return errors.New("fixed supervisor v2 authority mutation is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.requireMutableLocked(); err != nil {
		return err
	}
	candidateState := cloneState(store.state)
	if err := mutation(&candidateState); err != nil {
		if errors.Is(err, errNoStoreMutation) {
			return nil
		}
		return err
	}
	if err := validateStateWithAttemptCapacity(candidateState, false); err != nil {
		return err
	}
	if err := store.takeAuthorityFaultLocked(abortPoint); err != nil {
		return err
	}
	if err := store.takeAuthorityFaultLocked(preIndeterminatePoint); err != nil {
		store.recoveryRequired = true
		return fmt.Errorf("%w: injected before v2 authority commit", ErrCommitOutcomeIndeterminate)
	}
	postRename := store.takeAuthorityFaultLocked(postRenamePoint)
	candidateRecords := cloneLifecycleRecords(store.lifecycles, candidateState.TimeHighWaterUnixSeconds)
	candidateTombstones := cloneEffectTombstones(store.effectTombstones)
	if err := store.commitV2WorldLocked(ctx, &candidateState, &candidateRecords, &candidateTombstones, nil, postRename); err != nil {
		return err
	}
	store.state = candidateState
	store.lifecycles = candidateRecords
	store.lifecycleSetDigest = lifecycleSetDigest(candidateRecords)
	store.effectTombstones = candidateTombstones
	return nil
}

func (store *FixedFileStoreV2) commitV2LifecycleCandidateLocked(
	ctx context.Context,
	abortPoint LifecycleStoreFault,
	preIndeterminatePoint LifecycleStoreFault,
	postRenamePoint LifecycleStoreFault,
	fenceOnAbort bool,
	candidateState *installationState,
	candidateRecords *[]lifecyclestate.Record,
	candidateTombstones *[]archivestate.EffectIndexEntry,
) error {
	abort := store.takeLifecycleFaultLocked(abortPoint)
	if abort != nil {
		if fenceOnAbort {
			store.recoveryRequired = true
		}
		return errors.New("fixed supervisor v2 lifecycle transaction confirmed before rename failure")
	}
	if err := store.takeLifecycleFaultLocked(preIndeterminatePoint); err != nil {
		store.recoveryRequired = true
		return fmt.Errorf("%w: injected before v2 lifecycle commit", ErrCommitOutcomeIndeterminate)
	}
	postRename := store.takeLifecycleFaultLocked(postRenamePoint)
	return store.commitV2WorldLocked(ctx, candidateState, candidateRecords, candidateTombstones, store.lifecycles, postRename)
}

func (store *FixedFileStoreV2) commitV2WorldLocked(
	ctx context.Context,
	candidateState *installationState,
	candidateRecords *[]lifecyclestate.Record,
	candidateTombstones *[]archivestate.EffectIndexEntry,
	previousRecords []lifecyclestate.Record,
	postRename error,
) error {
	loaded, err := loadV2State(store.path)
	if err != nil {
		store.repairRequired = true
		return fmt.Errorf("%w: v2 mutation source: %v", ErrStoreRepairRequired, err)
	}
	if loaded.Active.View().SnapshotGeneration != store.active.View().SnapshotGeneration ||
		loaded.Active.View().CurrentCheckpoint != store.active.View().CurrentCheckpoint {
		store.recoveryRequired = true
		return ErrArchiveStaleTransaction
	}
	if uint64(loaded.Active.View().SnapshotGeneration) >= v0candidate.MaxSafeInteger {
		return classified(ClassificationCapacity, "v2-snapshot-generation")
	}
	nextGeneration := v0candidate.PositiveUInt53(uint64(loaded.Active.View().SnapshotGeneration) + 1)
	if previousRecords != nil {
		for index := range *candidateRecords {
			attemptID := (*candidateRecords)[index].Bindings().View().AttemptID
			previous := findLifecycle(previousRecords, attemptID)
			if previous >= 0 && reflect.DeepEqual(
				lifecycleRecordToDisk(previousRecords[previous]), lifecycleRecordToDisk((*candidateRecords)[index]),
			) {
				continue
			}
			view := (*candidateRecords)[index].View()
			// Lifecycle record generations predate the v2 envelope generation and
			// may already be higher at migration. Never move that independent
			// monotonic counter backwards; only raise it when the envelope commit
			// generation has overtaken it.
			if view.SnapshotGeneration < nextGeneration {
				view.SnapshotGeneration = nextGeneration
			}
			record, recordErr := lifecyclestate.NewRecord(view, candidateState.TimeHighWaterUnixSeconds)
			if recordErr != nil {
				return recordErr
			}
			(*candidateRecords)[index] = record
		}
	}
	normalized, err := effectTombstonesAtHotAttemptLocations(*candidateState, *candidateTombstones)
	if err != nil {
		return err
	}
	for index := range normalized {
		if _, alreadyRetained := store.retainedEffects[normalized[index].EffectID]; !alreadyRetained {
			normalized[index].IssuanceSnapshotGeneration = nextGeneration
		}
	}
	*candidateTombstones = normalized
	envelope, active, genesis, err := buildMutatedEnvelopeV2(loaded, *candidateState, *candidateRecords, normalized, nextGeneration)
	if err != nil {
		return err
	}
	encoded, err := encodeEnvelopeV2(envelope)
	if err != nil {
		return err
	}
	if store.v2Owner != nil {
		if err := store.v2Owner.CheckHeld(ctx); err != nil || store.v2Owner.OwnerSessionID() != store.ownerSessionID {
			store.ownerFenced = true
			store.recoveryRequired = true
			return ErrStoreOwnerFenced
		}
	}
	if err := persistV2Mutation(store.path, encoded, postRename); err != nil {
		if errors.Is(err, ErrCommitOutcomeIndeterminate) {
			store.recoveryRequired = true
			store.state = cloneState(*candidateState)
			store.lifecycles = cloneLifecycleRecords(*candidateRecords, candidateState.TimeHighWaterUnixSeconds)
			store.effectTombstones = cloneEffectTombstones(normalized)
		}
		return err
	}
	if store.v2Owner != nil {
		if err := store.v2Owner.CheckHeld(ctx); err != nil || store.v2Owner.OwnerSessionID() != store.ownerSessionID {
			store.ownerFenced = true
			store.recoveryRequired = true
			return ErrStoreOwnerFenced
		}
	}
	store.active = active
	store.genesis = genesis
	for _, tombstone := range normalized {
		store.retainedEffects[tombstone.EffectID] = struct{}{}
	}
	return nil
}

func buildMutatedEnvelopeV2(
	loaded loadedV2State,
	state installationState,
	records []lifecyclestate.Record,
	tombstones []archivestate.EffectIndexEntry,
	nextGeneration v0candidate.PositiveUInt53,
) (diskEnvelopeV2, archivestate.ActiveStateV2, archivestate.MigrationGenesisCheckpoint, error) {
	indexes, err := reconstructV2IndexesForWorld(state, records, tombstones, loaded.Segments)
	if err != nil {
		return diskEnvelopeV2{}, archivestate.ActiveStateV2{}, archivestate.MigrationGenesisCheckpoint{}, err
	}
	hotIndexes, err := reconstructV2Indexes(state, records, tombstones)
	if err != nil {
		return diskEnvelopeV2{}, archivestate.ActiveStateV2{}, archivestate.MigrationGenesisCheckpoint{}, err
	}
	hotCounts := countsForIndexView(hotIndexes.View())
	archivedCounts := loaded.Active.View().ArchivedCounts
	totalCounts, err := addArchiveCounts(hotCounts, archivedCounts)
	if err != nil {
		return diskEnvelopeV2{}, archivestate.ActiveStateV2{}, archivestate.MigrationGenesisCheckpoint{}, err
	}
	tombstoneDigest, err := archivestate.DigestEffectTombstoneSet(tombstones)
	if err != nil {
		return diskEnvelopeV2{}, archivestate.ActiveStateV2{}, archivestate.MigrationGenesisCheckpoint{}, err
	}
	hotDigests := archivestate.HotSetDigests{
		Registrations: state.RegistrationSetDigest, Approvals: state.ApprovalSetDigest,
		Attempts: state.AttemptSetDigest, Lifecycles: [32]byte(lifecycleSetDigest(records)),
		EffectTombstones: tombstoneDigest,
	}
	descriptors := make([]archivestate.ArchiveDescriptor, len(loaded.Envelope.Descriptors))
	for index, view := range loaded.Envelope.Descriptors {
		descriptor, descriptorErr := archivestate.NewArchiveDescriptor(view)
		if descriptorErr != nil {
			return diskEnvelopeV2{}, archivestate.ActiveStateV2{}, archivestate.MigrationGenesisCheckpoint{}, descriptorErr
		}
		descriptors[index] = descriptor
	}
	active, err := archivestate.NewActiveStateV2(archivestate.ActiveStateV2View{
		StoreFormatVersion: archivestate.SupervisorStoreFormatV2, MigrationSourceVersion: archivestate.MigrationSourceFormatV1,
		SnapshotGeneration: nextGeneration, ArchiveGeneration: loaded.Active.View().ArchiveGeneration,
		InstallationID: state.InstallationID, SupervisorID: state.SupervisorID,
		EpochSequence: state.EpochSequence, EpochDigest: state.EpochDigest, DurableTimeHighWater: state.TimeHighWaterUnixSeconds,
		Descriptors: descriptors, DescriptorSetDigest: loaded.Active.View().DescriptorSetDigest,
		Indexes: indexes, IndexDigests: indexes.Digests(), CombinedIndexDigest: indexes.CombinedDigest(), HotSetDigests: hotDigests,
		EffectTombstoneCoverage:   archivestate.EffectTombstoneCoverageVisibleV1AndV2,
		VisibleV1EffectSeedCount:  loaded.Active.View().VisibleV1EffectSeedCount,
		VisibleV1EffectSeedDigest: loaded.Active.View().VisibleV1EffectSeedDigest,
		PreviousCheckpoint:        loaded.Active.View().PreviousCheckpoint, CurrentCheckpoint: loaded.Active.View().CurrentCheckpoint,
		HotCounts: hotCounts, ArchivedCounts: archivedCounts, TotalCounts: totalCounts,
	})
	if err != nil {
		return diskEnvelopeV2{}, archivestate.ActiveStateV2{}, archivestate.MigrationGenesisCheckpoint{}, err
	}
	var genesisIndexes archivestate.ArchiveIndexes
	var genesisHotDigests archivestate.HotSetDigests
	if loaded.Envelope.MigrationGenesisIndexes != nil {
		genesisIndexes, err = archiveIndexesFromDisk(*loaded.Envelope.MigrationGenesisIndexes)
		genesisHotDigests = loaded.Envelope.MigrationGenesis.HotSetDigests
	} else {
		genesisIndexes, err = reconstructMigrationGenesisIndexes(loaded.State, loaded.Lifecycles, loaded.Segments)
		if err == nil {
			genesisHotDigests, err = reconstructMigrationGenesisHotSetDigests(loaded.State, loaded.Lifecycles, loaded.Segments)
		}
	}
	if err != nil {
		return diskEnvelopeV2{}, archivestate.ActiveStateV2{}, archivestate.MigrationGenesisCheckpoint{}, err
	}
	genesis, err := validateMigrationGenesisDisk(loaded.Envelope, genesisIndexes, genesisHotDigests)
	if err != nil {
		return diskEnvelopeV2{}, archivestate.ActiveStateV2{}, archivestate.MigrationGenesisCheckpoint{}, err
	}
	candidate := loaded.Envelope
	candidate.SnapshotGeneration = nextGeneration
	candidate.State = cloneState(state)
	candidate.LifecycleSetDigest = lifecycleSetDigest(records)
	candidate.Lifecycles = make([]lifecycleRecordDisk, len(records))
	for index, record := range records {
		candidate.Lifecycles[index] = lifecycleRecordToDisk(record)
	}
	tombstoneDisk := effectTombstonesToDisk(tombstones)
	candidate.EffectTombstones = &tombstoneDisk
	candidate.EffectTombstoneSetDigest = &tombstoneDigest
	genesisDisk := archiveIndexesToDisk(genesisIndexes)
	candidate.MigrationGenesisIndexes = &genesisDisk
	candidate.Indexes = archiveIndexesToDisk(indexes)
	candidate.IndexDigests = indexes.Digests()
	candidate.CombinedIndexDigest = indexes.CombinedDigest()
	candidate.HotCounts, candidate.ArchivedCounts, candidate.TotalCounts = hotCounts, archivedCounts, totalCounts
	return candidate, active, genesis, nil
}

func persistV2Mutation(path string, encoded []byte, afterRename error) error {
	parent := filepath.Dir(path)
	temporary, err := os.CreateTemp(parent, ".capsule-supervisor-v2-transaction-*.tmp")
	if err != nil {
		return errors.New("create fixed supervisor v2 transaction")
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
		return errors.New("protect fixed supervisor v2 transaction")
	}
	if _, err := temporary.Write(encoded); err != nil {
		return errors.New("write fixed supervisor v2 transaction")
	}
	if err := temporary.Sync(); err != nil {
		return errors.New("sync fixed supervisor v2 transaction")
	}
	if err := temporary.Close(); err != nil {
		return errors.New("close fixed supervisor v2 transaction")
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return errors.New("commit fixed supervisor v2 transaction")
	}
	committed = true
	if afterRename != nil {
		return fmt.Errorf("%w: injected after v2 rename", ErrCommitOutcomeIndeterminate)
	}
	directory, err := os.Open(parent) //nolint:gosec // trusted installation-configured store parent
	if err != nil {
		return fmt.Errorf("%w: open v2 store directory", ErrCommitOutcomeIndeterminate)
	}
	syncErr, closeErr := directory.Sync(), directory.Close()
	if syncErr != nil || closeErr != nil {
		return fmt.Errorf("%w: sync v2 store directory", ErrCommitOutcomeIndeterminate)
	}
	return nil
}

func (store *FixedFileStoreV2) takeAuthorityFaultLocked(point StoreFault) error {
	if point == "" {
		return nil
	}
	err := store.authorityFaults[point]
	delete(store.authorityFaults, point)
	return err
}
