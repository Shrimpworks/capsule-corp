package registrationstate

import (
	"context"
	"encoding/json"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// TestFixedStoreV1AttemptWithoutLifecycleBlocksV2Projection retains the exact
// F2 stop witness. A committed attempt without a lifecycle record is a valid
// v1 startup-recovery state, but the corrected passive v2 indexes require a
// lifecycle disposition on every attempt entry and derive the lifecycle count
// from the attempt count. F2 therefore cannot migrate this valid v1 world
// without inventing a lifecycle state/record or violating the count equations.
func TestFixedStoreV1AttemptWithoutLifecycleBlocksV2Projection(t *testing.T) {
	harness := newApprovalHarness(t)
	vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
	component := mustApprovalAttemptComponent(
		t, harness, []approvalattempt.FixtureVector{vector},
		&approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil,
	)
	submission, err := component.SubmitApproval(
		context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes,
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := component.RequestAttempt(
		context.Background(), attemptCall(), harness.registrationID, submission.Reference,
	); err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateFixedFileStoreV0ToV1(
		context.Background(), harness.path,
		V0ToV1MigrationOptions{Lock: &migrationLockStub{held: true}},
	); err != nil {
		t.Fatal(err)
	}

	store, err := OpenFixedFileStoreV1(harness.path)
	if err != nil {
		t.Fatalf("valid v1 reopen: %v", err)
	}
	snapshot, err := store.snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.State.Attempts) != 1 || len(snapshot.Lifecycles) != 0 {
		t.Fatalf("stop witness source counts = attempts %d, lifecycles %d, want 1/0",
			len(snapshot.State.Attempts), len(snapshot.Lifecycles))
	}
	if err := validateV1State(snapshot.State, snapshot.Lifecycles, snapshot.LifecycleSetDigest); err != nil {
		t.Fatalf("stop witness is not a complete valid v1 state: %v", err)
	}

	attempt := snapshot.State.Attempts[0]
	location, err := archivestate.NewHotRecordLocation(archivestate.RecordAttempt, 1)
	if err != nil {
		t.Fatal(err)
	}
	exactAttempt, err := json.Marshal(attempt)
	if err != nil {
		t.Fatal(err)
	}
	recordDigest, err := archivestate.DigestRecord(archivestate.RecordAttempt, exactAttempt)
	if err != nil {
		t.Fatal(err)
	}
	view := archivestate.ArchiveIndexesView{
		Scope:         archivestate.ArchiveIndexScopeRetainedGlobal,
		Registrations: make([]archivestate.RegistrationIndexEntry, 0),
		Approvals:     make([]archivestate.ApprovalIndexEntry, 0),
		Attempts: []archivestate.AttemptIndexEntry{{
			AttemptID: attempt.AttemptID, ApprovalID: attempt.ApprovalID,
			RegistrationID: attempt.RegistrationID, CreatedAt: attempt.CreatedAt,
			Location: location, FullRecordDigest: recordDigest,
		}},
		Nonces: make([]archivestate.NonceIndexEntry, 0), Effects: make([]archivestate.EffectIndexEntry, 0),
		Instances:      make([]archivestate.InstanceIndexEntry, 0),
		ApprovalReplay: make([]archivestate.ApprovalReplayIndexEntry, 0),
		AttemptReplay:  make([]archivestate.AttemptReplayIndexEntry, 0),
	}
	if _, err := archivestate.NewArchiveIndexes(view); err == nil {
		t.Fatal("v2 attempt index accepted the absent lifecycle disposition")
	}

	// Supplying a lifecycle state makes the index constructible only by
	// inventing successor state that is absent from the valid v1 source.
	view.Attempts[0].LifecycleState = lifecyclestate.StatePreparePending
	indexes, err := archivestate.NewArchiveIndexes(view)
	if err != nil {
		t.Fatalf("construct invented-state witness: %v", err)
	}
	descriptorDigest, err := archivestate.DigestArchiveDescriptorSet([]archivestate.ArchiveDescriptor{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = archivestate.NewMigrationGenesisCheckpoint(archivestate.MigrationGenesisCheckpointView{
		StoreFormatVersion:       archivestate.SupervisorStoreFormatV2,
		MigrationSourceVersion:   archivestate.MigrationSourceFormatV1,
		ResultSnapshotGeneration: 1, ArchiveGeneration: 1,
		DescriptorSetDigest: descriptorDigest, Indexes: indexes,
		HotSetDigests: archivestate.HotSetDigests{
			Registrations: [32]byte{1}, Approvals: [32]byte{1},
			Attempts: [32]byte{1}, Lifecycles: [32]byte{1},
		},
		VisibleV1EffectSeed: []lifecyclestate.EffectID{},
		HotCounts:           archivestate.ArchiveCounts{Attempts: 1, Lifecycles: 0},
		InstallationID:      v0candidate.InstallationID{1}, SupervisorID: v0candidate.SupervisorID{1},
		EpochDigest: v0candidate.TrustEpochDigest{1},
	})
	if err == nil {
		t.Fatal("v2 migration genesis represented one hot attempt and zero hot lifecycles")
	}
}
