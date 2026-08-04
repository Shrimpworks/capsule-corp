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

// TestFixedStoreV1AttemptWithoutLifecycleHasExactV2Projection retains the
// original F2 stop witness and proves the passive resolution. A committed
// attempt without a lifecycle record remains a valid v1 startup-recovery state;
// the v2 projection retains the attempt on the explicit absent arm and counts
// zero lifecycles without inventing state or bytes.
func TestFixedStoreV1AttemptWithoutLifecycleHasExactV2Projection(t *testing.T) {
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
			Lifecycle: archivestate.NoAttemptLifecycle(), Location: location, FullRecordDigest: recordDigest,
		}},
		Nonces: make([]archivestate.NonceIndexEntry, 0), Effects: make([]archivestate.EffectIndexEntry, 0),
		Instances:      make([]archivestate.InstanceIndexEntry, 0),
		ApprovalReplay: make([]archivestate.ApprovalReplayIndexEntry, 0),
		AttemptReplay:  make([]archivestate.AttemptReplayIndexEntry, 0),
	}
	indexes, err := archivestate.NewArchiveIndexes(view)
	if err != nil {
		t.Fatalf("v2 attempt index rejected explicit lifecycle absence: %v", err)
	}

	// Presence cannot be asserted with an attempt-record anchor. A real present
	// arm requires a separately typed lifecycle location and full-record digest.
	if _, err := archivestate.NewPresentAttemptLifecycle(
		lifecyclestate.StatePreparePending, location, recordDigest,
	); err == nil {
		t.Fatal("present lifecycle accepted an attempt-record anchor")
	}
	descriptorDigest, err := archivestate.DigestArchiveDescriptorSet([]archivestate.ArchiveDescriptor{})
	if err != nil {
		t.Fatal(err)
	}
	genesis, err := archivestate.NewMigrationGenesisCheckpoint(archivestate.MigrationGenesisCheckpointView{
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
	if err != nil {
		t.Fatalf("v2 migration genesis rejected one hot attempt and zero hot lifecycles: %v", err)
	}
	if genesis.View().HotCounts.Attempts != 1 || genesis.View().HotCounts.Lifecycles != 0 {
		t.Fatalf("resolved genesis counts = attempts %d, lifecycles %d, want 1/0",
			genesis.View().HotCounts.Attempts, genesis.View().HotCounts.Lifecycles)
	}
}
