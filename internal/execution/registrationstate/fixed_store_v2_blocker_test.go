package registrationstate

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// TestFixedStoreV2F4BHistoricalEffectTombstoneCannotSatisfyF4AReconstruction
// retains the exact F4B stop witness. ADR-0031 requires the visible v1 effect
// seed and every v2-issued effect tombstone to survive when a later lifecycle
// operation replaces the record's current EffectID. F4A's full verifier,
// however, reconstructs exactly one effect from that current field. A closed,
// sorted two-entry ledger is therefore rejected before lookup can follow it.
func TestFixedStoreV2F4BHistoricalEffectTombstoneCannotSatisfyF4AReconstruction(t *testing.T) {
	state, template := stateAndLifecycleRecord(t)
	current := lifecycleRecordForV2State(t, state, template, lifecyclestate.StateDestroyConfirmed)
	currentView := current.View()
	currentView.RecordVersion++
	currentView.OperationSequence = 2
	currentView.EffectID = effectIDForF4BBlocker(t, 0x22)
	current, err := lifecyclestate.NewRecord(currentView, state.TimeHighWaterUnixSeconds)
	if err != nil {
		t.Fatalf("construct current post-v2 lifecycle: %v", err)
	}

	source := loadedV1State{
		SourceFormatVersion: SupervisorStoreFormatV0,
		State:               state,
		LifecycleSetDigest:  lifecycleSetDigest([]lifecyclestate.Record{current}),
		Lifecycles:          []lifecyclestate.Record{current},
	}
	envelope, err := buildEnvelopeV2(source)
	if err != nil {
		t.Fatalf("build closed v2 witness: %v", err)
	}
	currentIndexes, err := archiveIndexesFromDisk(envelope.Indexes)
	if err != nil {
		t.Fatal(err)
	}
	view := currentIndexes.View()
	if len(view.Effects) != 1 {
		t.Fatalf("current lifecycle reconstructed %d effects, want 1", len(view.Effects))
	}

	// This is the required history after a visible-v1 seed was replaced by a
	// newly issued v2 destroy intent. Both entries are valid tombstone values;
	// only the newer one can equal the lifecycle record's current effect field.
	historical := view.Effects[0]
	historical.EffectID = effectIDForF4BBlocker(t, 0x11)
	historical.OperationSequence = 1
	historical.Operation = lifecyclestate.OperationPrepare
	historical.VisibleV1Seed = true
	view.Effects[0].VisibleV1Seed = false
	view.Effects = append(view.Effects, historical)
	sort.Slice(view.Effects, func(left, right int) bool {
		return bytes.Compare(view.Effects[left].EffectID[:], view.Effects[right].EffectID[:]) < 0
	})
	retainedLedger, err := archivestate.NewArchiveIndexes(view)
	if err != nil {
		t.Fatalf("closed append-only effect ledger is not passively representable: %v", err)
	}
	if len(retainedLedger.View().Effects) != 2 {
		t.Fatal("append-only effect ledger lost its historical tombstone")
	}

	envelope.Indexes = archiveIndexesToDisk(retainedLedger)
	envelope.IndexDigests = retainedLedger.Digests()
	envelope.CombinedIndexDigest = retainedLedger.CombinedDigest()
	envelope.HotCounts.Effects = 2
	envelope.TotalCounts.Effects = 2
	seedDigest, err := archivestate.DigestVisibleV1EffectSeed([]lifecyclestate.EffectID{historical.EffectID})
	if err != nil {
		t.Fatal(err)
	}
	envelope.VisibleV1EffectSeedCount = 1
	envelope.VisibleV1EffectSeedDigest = seedDigest

	path := filepath.Join(t.TempDir(), "supervisor-state.json")
	writeV2Envelope(t, path, envelope)
	before := mustReadFile(t, path)
	if _, err := OpenFixedFileStoreV2(path); err == nil || !strings.Contains(err.Error(), "reconstructed index mismatch") {
		t.Fatalf("F4A verifier accepted or misclassified retained historical effect ledger: %v", err)
	}
	if after := mustReadFile(t, path); !bytes.Equal(after, before) {
		t.Fatal("F4A refusal rewrote the blocked v2 witness")
	}

	if current.View().EffectID == historical.EffectID {
		t.Fatal("historical tombstone unexpectedly equals the current lifecycle effect")
	}
}

func effectIDForF4BBlocker(t *testing.T, marker byte) lifecyclestate.EffectID {
	t.Helper()
	domain, err := lifecyclestate.NewDomainIdentifier(
		lifecyclestate.DomainEffectID,
		bytes.Repeat([]byte{marker}, 16),
	)
	if err != nil {
		t.Fatal(err)
	}
	identifier, err := lifecyclestate.NewEffectID(domain)
	if err != nil {
		t.Fatal(err)
	}
	return identifier
}

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
