package registrationstate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
)

func TestFixedStoreV1ToV2MigrationAndDowngradeRefusal(t *testing.T) {
	state, _ := stateAndLifecycleRecord(t)
	path := newV1PathFromState(t, state, nil)
	v1Bytes := mustReadFile(t, path)
	v1, err := OpenFixedFileStoreV1(path)
	if err != nil {
		t.Fatal(err)
	}
	recovery, err := v1.RecoveryAttemptIDs(context.Background())
	if err != nil || len(recovery) != 1 || recovery[0] != state.Attempts[0].AttemptID {
		t.Fatalf("v1 missing-lifecycle recovery set = %#v, %v", recovery, err)
	}
	if _, err := OpenFixedFileStoreV2(path); err == nil {
		t.Fatal("v2 opener accepted v1")
	}
	if got := mustReadFile(t, path); !bytes.Equal(got, v1Bytes) {
		t.Fatal("v2 downgrade refusal rewrote v1")
	}

	lock := &migrationLockStub{held: true}
	store, err := MigrateFixedFileStoreV1ToV2(
		context.Background(), path, V1ToV2MigrationOptions{Lock: lock},
	)
	if err != nil {
		t.Fatalf("migrate v1 to v2: %v", err)
	}
	if lock.checks != 3 {
		t.Fatalf("owner checks = %d, want 3", lock.checks)
	}
	snapshot, err := store.snapshotV2(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !authorityStatesEqual(state, snapshot.State) || len(snapshot.Lifecycles) != 0 {
		t.Fatal("v2 migration changed exact v1 authority/lifecycle state")
	}
	active := snapshot.Active.View()
	if active.StoreFormatVersion != 2 || active.MigrationSourceVersion != 1 ||
		active.SnapshotGeneration != 1 || active.ArchiveGeneration != 1 ||
		len(active.Descriptors) != 0 || active.ArchivedCounts != (archivestate.ArchiveCounts{}) {
		t.Fatalf("v2 migration metadata = %#v", active)
	}
	if active.HotCounts.Attempts != 1 || active.HotCounts.Lifecycles != 0 ||
		active.TotalCounts != active.HotCounts {
		t.Fatalf("missing-lifecycle v2 counts = %#v", active.HotCounts)
	}
	indexes := active.Indexes.View()
	if len(indexes.Registrations) != 1 || len(indexes.Approvals) != 1 ||
		len(indexes.Attempts) != 1 || len(indexes.Nonces) != 1 ||
		len(indexes.ApprovalReplay) != 1 || len(indexes.AttemptReplay) != 1 ||
		len(indexes.Effects) != 0 || len(indexes.Instances) != 0 ||
		indexes.Attempts[0].Lifecycle.View().Presence != archivestate.AttemptLifecycleAbsent {
		t.Fatalf("missing-lifecycle retained-global indexes = %#v", indexes)
	}
	if snapshot.Genesis.Reference() != active.CurrentCheckpoint ||
		snapshot.Genesis.View().HotCounts != active.HotCounts {
		t.Fatal("migration genesis is not the exact active checkpoint head")
	}

	v2Bytes := mustReadFile(t, path)
	if _, err := OpenFixedFileStoreV1(path); err == nil {
		t.Fatal("v1 opener accepted v2")
	}
	if _, err := NewFixedFileStore(path, InitialState{}); err == nil {
		t.Fatal("v0 opener accepted v2")
	}
	if _, err := MigrateFixedFileStoreV1ToV2(
		context.Background(), path,
		V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}},
	); err == nil {
		t.Fatal("repeat v1-to-v2 migration accepted v2")
	}
	if got := mustReadFile(t, path); !bytes.Equal(got, v2Bytes) {
		t.Fatal("downgrade or repeat migration rewrote v2")
	}

	firstReport, err := VerifyFixedFileStoreV2(path)
	if err != nil {
		t.Fatal(err)
	}
	secondReport, err := VerifyFixedFileStoreV2(path)
	if err != nil || firstReport != secondReport {
		t.Fatalf("read-only verification is nondeterministic: %#v %#v %v", firstReport, secondReport, err)
	}

	copySnapshot := snapshot
	copySnapshot.State.Approvals[0].ExactPayloadBytes[0] ^= 0xff
	copyIndexes := copySnapshot.Active.View().Indexes.View()
	copyIndexes.Attempts[0].AttemptID[0] ^= 0xff
	reopened, err := OpenFixedFileStoreV2(path)
	if err != nil {
		t.Fatal(err)
	}
	reopenedSnapshot, _ := reopened.snapshotV2(context.Background())
	if !authorityStatesEqual(state, reopenedSnapshot.State) ||
		reopenedSnapshot.Active.View().Indexes.View().Attempts[0].AttemptID != state.Attempts[0].AttemptID {
		t.Fatal("v2 read projection aliased caller-owned memory")
	}
}

func TestFixedStoreV2MigrationDeterministicSerializationAndDigests(t *testing.T) {
	state, record := stateAndLifecycleRecord(t)
	paths := []string{
		newV1PathFromState(t, state, []lifecyclestate.Record{record}),
		newV1PathFromState(t, state, []lifecyclestate.Record{record}),
	}
	for _, path := range paths {
		if _, err := MigrateFixedFileStoreV1ToV2(
			context.Background(), path,
			V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}},
		); err != nil {
			t.Fatal(err)
		}
	}
	left := mustReadFile(t, paths[0])
	right := mustReadFile(t, paths[1])
	if !bytes.Equal(left, right) {
		t.Fatal("identical v1 worlds produced different v2 bytes")
	}
	store, err := OpenFixedFileStoreV2(paths[0])
	if err != nil {
		t.Fatal(err)
	}
	snapshot, _ := store.snapshotV2(context.Background())
	view := snapshot.Active.View()
	if view.VisibleV1EffectSeedCount != 0 ||
		view.VisibleV1EffectSeedDigest != snapshot.Genesis.SeedDigest() ||
		view.CombinedIndexDigest != snapshot.Genesis.View().Indexes.CombinedDigest() {
		t.Fatal("deterministic v2 digest cross-links disagree")
	}
	bytesDigest := sha256.Sum256(left)
	genesisDigest := snapshot.Genesis.Digest()
	t.Logf("v2 bytes sha256=%s combined-index=%s genesis=%s counts=%+v",
		hex.EncodeToString(bytesDigest[:]),
		hex.EncodeToString(view.CombinedIndexDigest[:]),
		hex.EncodeToString(genesisDigest[:]), view.HotCounts)
}

type fixedStoreV2KnownAnswer struct {
	BytesSHA256       string
	CombinedIndex     string
	GenesisCheckpoint string
	VisibleSeed       string
	Counts            archivestate.ArchiveCounts
}

func TestFixedStoreV2KnownAnswers(t *testing.T) {
	emptyState := ordinaryInitialStateForV1()
	missingState, template := stateAndLifecycleRecord(t)
	observed := lifecycleRecordForV2State(t, missingState, template, lifecyclestate.StateObserved)
	got := []fixedStoreV2KnownAnswer{
		fixedStoreV2KnownAnswerFor(t, emptyState, nil),
		fixedStoreV2KnownAnswerFor(t, missingState, nil),
		fixedStoreV2KnownAnswerFor(t, missingState, []lifecyclestate.Record{observed}),
	}
	want := []fixedStoreV2KnownAnswer{
		{
			BytesSHA256:       "c845344aa1b464d7cd40ba86d33ecf3e5797cfed8406767dff0db12be9d9bb04",
			CombinedIndex:     "78e817b6a07989095010743601a017e43e3b660ea78ad0231e01d900227e207c",
			GenesisCheckpoint: "37336388a2f775b79d4f56e1af8ff1afd45de6cea96a64bfc27cec564290c88c",
			VisibleSeed:       "17de5f44f523dab94ca4b215ce7779358146fb094fa6d208e0190cb0ba69e0a1",
			Counts:            archivestate.ArchiveCounts{},
		},
		{
			BytesSHA256:       "569ba7c1aa25432a1001b1ca7122a7772ccfe954a18a484d26ed71b4255b8dca",
			CombinedIndex:     "924c78b9508123feb1b78fd62b71df6cace9c97b4887f2d67bbdc6ef2a9a7de5",
			GenesisCheckpoint: "983c2474dbef1fa6908de0fa02f96e2aaf7245bce48399e224a0f8a2c349a23e",
			VisibleSeed:       "17de5f44f523dab94ca4b215ce7779358146fb094fa6d208e0190cb0ba69e0a1",
			Counts: archivestate.ArchiveCounts{
				Registrations: 1, Approvals: 1, Attempts: 1, Nonces: 1,
				ApprovalReplay: 1, AttemptReplay: 1,
			},
		},
		{
			BytesSHA256:       "29c706be5d8a55958acacae7ad01a001576a307d847fd610a6f2c0e57f291235",
			CombinedIndex:     "2dc21df5c66bdb46f7ba80ed6566b12912a6b9b79a1f5b7a3d325f507bea65c2",
			GenesisCheckpoint: "3b5de809e6cbcda85b94439a9142cdc69243b1d18f912dbe8fc81ba6e4101a99",
			VisibleSeed:       "acee3fa25e62c185eb1e9b26313b6f04b9c88e57857e0fb400371c2c5a67295f",
			Counts: archivestate.ArchiveCounts{
				Registrations: 1, Approvals: 1, Attempts: 1, Lifecycles: 1,
				Nonces: 1, Effects: 1, Instances: 1, ApprovalReplay: 1, AttemptReplay: 1,
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("fixed-store v2 known answers changed\n got: %#v\nwant: %#v", got, want)
	}
}

func fixedStoreV2KnownAnswerFor(
	t *testing.T,
	state installationState,
	records []lifecyclestate.Record,
) fixedStoreV2KnownAnswer {
	t.Helper()
	path := newV1PathFromState(t, state, records)
	store := mustMigrateV2(t, path)
	snapshot := mustV2Snapshot(t, store)
	active := snapshot.Active.View()
	bytesDigest := sha256.Sum256(mustReadFile(t, path))
	genesisDigest := snapshot.Genesis.Digest()
	return fixedStoreV2KnownAnswer{
		BytesSHA256:       hex.EncodeToString(bytesDigest[:]),
		CombinedIndex:     hex.EncodeToString(active.CombinedIndexDigest[:]),
		GenesisCheckpoint: hex.EncodeToString(genesisDigest[:]),
		VisibleSeed:       hex.EncodeToString(active.VisibleV1EffectSeedDigest[:]),
		Counts:            active.HotCounts,
	}
}

func TestFixedStoreV2MapsZeroMixedAndEveryLifecycleState(t *testing.T) {
	t.Run("zero-attempts", func(t *testing.T) {
		path := newV1PathFromState(t, ordinaryInitialStateForV1(), nil)
		store := mustMigrateV2(t, path)
		view := mustV2Snapshot(t, store).Active.View()
		if view.HotCounts != (archivestate.ArchiveCounts{}) {
			t.Fatalf("zero-state counts = %#v", view.HotCounts)
		}
	})

	t.Run("mixed-absent-present", func(t *testing.T) {
		state, records := generatedLifecyclePopulation(t, 2, lifecyclestate.StatePreparePending)
		path := newV1PathFromState(t, state, records[:1])
		view := mustV2Snapshot(t, mustMigrateV2(t, path)).Active.View()
		if view.HotCounts.Attempts != 2 || view.HotCounts.Lifecycles != 1 {
			t.Fatalf("mixed counts = %#v", view.HotCounts)
		}
		attempts := view.Indexes.View().Attempts
		present := 0
		absent := 0
		for _, attempt := range attempts {
			switch attempt.Lifecycle.View().Presence {
			case archivestate.AttemptLifecyclePresent:
				present++
			case archivestate.AttemptLifecycleAbsent:
				absent++
			}
		}
		if present != 1 || absent != 1 {
			t.Fatalf("mixed lifecycle arms = present %d absent %d", present, absent)
		}
	})

	for _, state := range lifecyclestate.LifecycleStates() {
		state := state
		t.Run(string(state), func(t *testing.T) {
			authority, template := stateAndLifecycleRecord(t)
			record := lifecycleRecordForV2State(t, authority, template, state)
			path := newV1PathFromState(t, authority, []lifecyclestate.Record{record})
			view := mustV2Snapshot(t, mustMigrateV2(t, path)).Active.View()
			attempt := view.Indexes.View().Attempts[0]
			if attempt.Lifecycle.View().Presence != archivestate.AttemptLifecyclePresent ||
				attempt.Lifecycle.View().State != state || view.HotCounts.Lifecycles != 1 {
				t.Fatalf("state %s mapping = %#v", state, attempt.Lifecycle.View())
			}
		})
	}
}

func TestFixedStoreV2MigrationFaultMatrixPreservesOldOrNewWorld(t *testing.T) {
	preRename := []MigrationFault{
		FaultV2MigrationAfterTempCreate, FaultV2MigrationAfterTempMode,
		FaultV2MigrationAfterTempWrite, FaultV2MigrationAfterTempSync,
		FaultV2MigrationAfterTempClose, FaultV2MigrationBeforeRename,
	}
	for _, point := range preRename {
		t.Run(string(point), func(t *testing.T) {
			state, _ := stateAndLifecycleRecord(t)
			path := newV1PathFromState(t, state, nil)
			before := mustReadFile(t, path)
			_, err := MigrateFixedFileStoreV1ToV2(
				context.Background(), path,
				V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}, Faults: &migrationFaultStub{point: point}},
			)
			if err == nil || errors.Is(err, ErrMigrationOutcomeIndeterminate) {
				t.Fatalf("confirmed fault result = %v", err)
			}
			if after := mustReadFile(t, path); !bytes.Equal(after, before) {
				t.Fatal("confirmed pre-rename fault changed v1 bytes")
			}
			if _, err := OpenFixedFileStoreV1(path); err != nil {
				t.Fatalf("old v1 world did not reopen: %v", err)
			}
		})
	}

	for _, point := range []MigrationFault{FaultV2MigrationAfterRename, FaultV2MigrationAfterDirSync} {
		t.Run(string(point), func(t *testing.T) {
			state, _ := stateAndLifecycleRecord(t)
			path := newV1PathFromState(t, state, nil)
			_, err := MigrateFixedFileStoreV1ToV2(
				context.Background(), path,
				V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}, Faults: &migrationFaultStub{point: point}},
			)
			if !errors.Is(err, ErrMigrationOutcomeIndeterminate) {
				t.Fatalf("post-rename result = %v", err)
			}
			if _, err := OpenFixedFileStoreV2(path); err != nil {
				t.Fatalf("new v2 world did not reopen: %v", err)
			}
			if _, err := OpenFixedFileStoreV1(path); err == nil {
				t.Fatal("post-rename world merged v1 and v2")
			}
		})
	}
}

func TestFixedStoreV2MigrationFullyValidatesV1BeforeAnyWrite(t *testing.T) {
	state, template := stateAndLifecycleRecord(t)
	record := lifecycleRecordForV2State(t, state, template, lifecyclestate.StateObserved)
	valid := encodedEnvelopeV1(state, []lifecyclestate.Record{record})
	tests := []struct {
		name   string
		mutate func(*diskEnvelopeV1)
	}{
		{name: "registration-set-digest", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.RegistrationSetDigest[0] ^= 0xff
		}},
		{name: "approval-set-digest", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.ApprovalSetDigest[0] ^= 0xff
		}},
		{name: "attempt-set-digest", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.AttemptSetDigest[0] ^= 0xff
		}},
		{name: "lifecycle-set-digest", mutate: func(envelope *diskEnvelopeV1) {
			envelope.LifecycleSetDigest[0] ^= 0xff
		}},
		{name: "duplicate-lifecycle", mutate: func(envelope *diskEnvelopeV1) {
			envelope.Lifecycles = append(envelope.Lifecycles, envelope.Lifecycles[0])
			envelope.LifecycleSetDigest = lifecycleSetDigest([]lifecyclestate.Record{record, record})
		}},
		{name: "time-high-water", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.TimeHighWaterUnixSeconds = 0
		}},
		{name: "cross-link", mutate: func(envelope *diskEnvelopeV1) {
			envelope.Lifecycles[0].Bindings.ApprovalID[0] ^= 0xff
		}},
		{name: "trust-transition", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.TrustPhase = TrustTransitionFenced
		}},
		{name: "quarantined", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.Quarantined = true
		}},
		{name: "attempts-disabled", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.AttemptsDisabled = true
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			envelope := cloneDiskEnvelopeV1(t, valid)
			test.mutate(&envelope)
			path := filepath.Join(t.TempDir(), "invalid-source-v1.json")
			writeV1Envelope(t, path, envelope)
			before := mustReadFile(t, path)
			if _, err := MigrateFixedFileStoreV1ToV2(
				context.Background(), path,
				V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}},
			); err == nil {
				t.Fatal("invalid v1 source migrated")
			}
			if got := mustReadFile(t, path); !bytes.Equal(got, before) {
				t.Fatal("invalid v1 source refusal wrote state")
			}
		})
	}

	for _, raw := range []struct {
		name   string
		mutate func([]byte) []byte
	}{
		{name: "truncated", mutate: func(data []byte) []byte { return data[:len(data)/2] }},
		{name: "unknown-field", mutate: func(data []byte) []byte {
			return bytes.Replace(data, []byte("{\"storeFormatVersion\""), []byte("{\"unknown\":1,\"storeFormatVersion\""), 1)
		}},
		{name: "duplicate-field", mutate: func(data []byte) []byte {
			return bytes.Replace(data, []byte("{\"storeFormatVersion\":1"), []byte("{\"storeFormatVersion\":1,\"storeFormatVersion\":1"), 1)
		}},
	} {
		t.Run(raw.name, func(t *testing.T) {
			path := newV1PathFromState(t, state, []lifecyclestate.Record{record})
			data := raw.mutate(mustReadFile(t, path))
			if err := os.WriteFile(path, data, 0o600); err != nil {
				t.Fatal(err)
			}
			before := mustReadFile(t, path)
			if _, err := MigrateFixedFileStoreV1ToV2(
				context.Background(), path,
				V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}},
			); err == nil {
				t.Fatal("raw-invalid v1 source migrated")
			}
			if got := mustReadFile(t, path); !bytes.Equal(got, before) {
				t.Fatal("raw-invalid v1 source refusal wrote state")
			}
		})
	}
}

func TestFixedStoreV2RefusesUnsafeSourceAndSuccessorFileShapes(t *testing.T) {
	state, _ := stateAndLifecycleRecord(t)
	path := newV1PathFromState(t, state, nil)
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	before := mustReadFile(t, path)
	if _, err := MigrateFixedFileStoreV1ToV2(
		context.Background(), path,
		V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}},
	); err == nil {
		t.Fatal("migration accepted wrong source mode")
	}
	if got := mustReadFile(t, path); !bytes.Equal(got, before) {
		t.Fatal("wrong-mode refusal rewrote source")
	}

	target := newV1PathFromState(t, state, nil)
	link := filepath.Join(t.TempDir(), "state-link.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateFixedFileStoreV1ToV2(
		context.Background(), link,
		V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}},
	); err == nil {
		t.Fatal("migration followed a source symlink")
	}
	if _, err := OpenFixedFileStoreV2(link); err == nil {
		t.Fatal("v2 opener followed a store symlink")
	}
	if _, err := OpenFixedFileStoreV1(target); err != nil {
		t.Fatalf("symlink refusal changed target: %v", err)
	}
}

func TestFixedStoreV2MigrationOwnerChecksAndConcurrency(t *testing.T) {
	state, _ := stateAndLifecycleRecord(t)
	path := newV1PathFromState(t, state, nil)
	before := mustReadFile(t, path)
	if _, err := MigrateFixedFileStoreV1ToV2(context.Background(), path, V1ToV2MigrationOptions{}); err == nil {
		t.Fatal("migration accepted a missing owner")
	}
	if _, err := MigrateFixedFileStoreV1ToV2(
		context.Background(), path,
		V1ToV2MigrationOptions{Lock: &migrationLockStub{held: false}},
	); err == nil {
		t.Fatal("migration accepted a failed entry owner check")
	}
	if got := mustReadFile(t, path); !bytes.Equal(got, before) {
		t.Fatal("owner refusal rewrote v1")
	}
	for _, failAt := range []int{2, 3} {
		t.Run("owner-check-"+string(rune('0'+failAt)), func(t *testing.T) {
			candidatePath := newV1PathFromState(t, state, nil)
			candidateBefore := mustReadFile(t, candidatePath)
			lock := &stagedMigrationLock{failAt: failAt}
			if failAt == 2 {
				lock.onCheck = func(check int) {
					if check != 2 {
						return
					}
					matches, globErr := filepath.Glob(filepath.Join(filepath.Dir(candidatePath), ".capsule-supervisor-v2-migration-*.tmp"))
					if globErr != nil || len(matches) != 1 {
						t.Fatalf("owner check immediately before rename saw %d temporary files: %v", len(matches), globErr)
					}
				}
			}
			_, err := MigrateFixedFileStoreV1ToV2(
				context.Background(), candidatePath, V1ToV2MigrationOptions{Lock: lock},
			)
			if err == nil || !errors.Is(err, ErrMigrationLockRequired) || lock.checks != failAt {
				t.Fatalf("owner check %d result = %v after %d checks", failAt, err, lock.checks)
			}
			if failAt == 2 {
				if got := mustReadFile(t, candidatePath); !bytes.Equal(got, candidateBefore) {
					t.Fatal("pre-commit owner loss changed v1 bytes")
				}
				if _, err := OpenFixedFileStoreV1(candidatePath); err != nil {
					t.Fatalf("pre-commit owner loss did not preserve v1: %v", err)
				}
			} else if _, err := OpenFixedFileStoreV2(candidatePath); err != nil {
				t.Fatalf("pre-reopen owner loss did not leave complete committed v2: %v", err)
			}
		})
	}

	const contenders = 8
	var wait sync.WaitGroup
	results := make(chan error, contenders)
	for range contenders {
		wait.Add(1)
		go func() {
			defer wait.Done()
			_, err := MigrateFixedFileStoreV1ToV2(
				context.Background(), path,
				V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}},
			)
			results <- err
		}()
	}
	wait.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("concurrent migration successes = %d, want 1", successes)
	}
	if _, err := OpenFixedFileStoreV2(path); err != nil {
		t.Fatalf("concurrent result is not complete v2: %v", err)
	}
}

func TestFixedStoreV2MigrationCapacityBoundaries(t *testing.T) {
	t.Run("exact-active-missing-lifecycle", func(t *testing.T) {
		state, _ := generatedLifecyclePopulation(t, MaxActiveLifecycleRecords, lifecyclestate.StatePreparePending)
		path := newV1PathFromState(t, state, nil)
		view := mustV2Snapshot(t, mustMigrateV2(t, path)).Active.View()
		if view.HotCounts.Attempts != MaxActiveLifecycleRecords || view.HotCounts.Lifecycles != 0 {
			t.Fatalf("exact active counts = %#v", view.HotCounts)
		}
	})

	t.Run("active-cap-plus-one-refuses", func(t *testing.T) {
		state, _ := generatedLifecyclePopulation(t, MaxActiveLifecycleRecords+1, lifecyclestate.StatePreparePending)
		path := newV1PathFromStateUnchecked(t, state, nil)
		before := mustReadFile(t, path)
		if _, err := MigrateFixedFileStoreV1ToV2(
			context.Background(), path,
			V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}},
		); err == nil {
			t.Fatal("active cap-plus-one migrated")
		}
		if got := mustReadFile(t, path); !bytes.Equal(got, before) {
			t.Fatal("active cap-plus-one refusal rewrote v1")
		}
	})

	t.Run("exact-retained", func(t *testing.T) {
		state, records := generatedLifecyclePopulation(t, MaxRetainedLifecycleRecords, lifecyclestate.StateDestroyed)
		path := newV1PathFromState(t, state, records)
		view := mustV2Snapshot(t, mustMigrateV2(t, path)).Active.View()
		if view.HotCounts.Attempts != MaxRetainedLifecycleRecords ||
			view.HotCounts.Lifecycles != MaxRetainedLifecycleRecords {
			t.Fatalf("exact retained counts = %#v", view.HotCounts)
		}
	})

	t.Run("retained-cap-plus-one-refuses", func(t *testing.T) {
		state, records := generatedLifecyclePopulation(t, MaxRetainedLifecycleRecords, lifecyclestate.StateDestroyed)
		records = append(records, records[len(records)-1])
		path := newV1PathFromStateUnchecked(t, state, records)
		before := mustReadFile(t, path)
		if _, err := MigrateFixedFileStoreV1ToV2(
			context.Background(), path,
			V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}},
		); err == nil {
			t.Fatal("retained cap-plus-one migrated")
		}
		if got := mustReadFile(t, path); !bytes.Equal(got, before) {
			t.Fatal("retained cap-plus-one refusal rewrote v1")
		}
	})

	t.Run("v2-byte-cap-plus-one", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "oversized-v2.json")
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		if err := file.Truncate(archivestate.MaxSupervisorStateV2Bytes + 1); err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
		before, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := OpenFixedFileStoreV2(path); err == nil {
			t.Fatal("v2 byte cap-plus-one opened")
		}
		after, err := os.Stat(path)
		if err != nil || before.Size() != after.Size() {
			t.Fatal("v2 byte cap refusal rewrote evidence")
		}
	})
}

func TestFixedStoreV2CorruptionRefusesWithoutRewrite(t *testing.T) {
	state, record := stateAndLifecycleRecord(t)
	record = lifecycleRecordForV2State(t, state, record, lifecyclestate.StateObserved)
	validPath := newV1PathFromState(t, state, []lifecyclestate.Record{record})
	mustMigrateV2(t, validPath)
	valid := readV2Envelope(t, validPath)
	tests := []struct {
		name   string
		mutate func(*diskEnvelopeV2)
	}{
		{name: "unsupported-version", mutate: func(envelope *diskEnvelopeV2) { value := uint64(3); envelope.StoreFormatVersion = &value }},
		{name: "missing-version", mutate: func(envelope *diskEnvelopeV2) { envelope.StoreFormatVersion = nil }},
		{name: "missing-index-field", mutate: func(envelope *diskEnvelopeV2) { envelope.Indexes.Attempts = nil }},
		{name: "descriptor-digest", mutate: func(envelope *diskEnvelopeV2) { envelope.DescriptorSetDigest[0] ^= 0xff }},
		{name: "index-digest", mutate: func(envelope *diskEnvelopeV2) { envelope.IndexDigests.Attempts[0] ^= 0xff }},
		{name: "combined-index-digest", mutate: func(envelope *diskEnvelopeV2) { envelope.CombinedIndexDigest[0] ^= 0xff }},
		{name: "count", mutate: func(envelope *diskEnvelopeV2) { envelope.HotCounts.Attempts++ }},
		{name: "location", mutate: func(envelope *diskEnvelopeV2) { envelope.Indexes.Attempts[0].Location.HotRecordOrdinal++ }},
		{name: "record-digest", mutate: func(envelope *diskEnvelopeV2) { envelope.Indexes.Attempts[0].FullRecordDigest[0] ^= 0xff }},
		{name: "cross-link", mutate: func(envelope *diskEnvelopeV2) { envelope.Indexes.Attempts[0].ApprovalID[0] ^= 0xff }},
		{name: "genesis-digest", mutate: func(envelope *diskEnvelopeV2) { envelope.MigrationGenesis.Digest[0] ^= 0xff }},
		{name: "checkpoint-kind", mutate: func(envelope *diskEnvelopeV2) {
			envelope.CurrentCheckpoint.Kind = archivestate.ArchiveCheckpointActivation
		}},
		{name: "time-high-water", mutate: func(envelope *diskEnvelopeV2) { envelope.State.TimeHighWaterUnixSeconds++ }},
		{name: "snapshot-generation", mutate: func(envelope *diskEnvelopeV2) { envelope.SnapshotGeneration++ }},
		{name: "archive-generation", mutate: func(envelope *diskEnvelopeV2) { envelope.ArchiveGeneration++ }},
		{name: "effect-seed-count", mutate: func(envelope *diskEnvelopeV2) { envelope.VisibleV1EffectSeedCount++ }},
		{name: "effect-index-on-different-attempt", mutate: func(envelope *diskEnvelopeV2) { envelope.Indexes.Effects[0].AttemptID[0] ^= 0xff }},
		{name: "instance-index-on-different-attempt", mutate: func(envelope *diskEnvelopeV2) { envelope.Indexes.Instances[0].AttemptID[0] ^= 0xff }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			envelope := cloneV2Envelope(t, valid)
			test.mutate(&envelope)
			path := filepath.Join(t.TempDir(), "corrupt-v2.json")
			writeV2Envelope(t, path, envelope)
			assertV2OpenDoesNotRewrite(t, path)
		})
	}

	rawCases := []struct {
		name   string
		mutate func([]byte) []byte
	}{
		{name: "truncated", mutate: func(data []byte) []byte { return data[:len(data)/2] }},
		{name: "trailing", mutate: func(data []byte) []byte { return append(data, []byte("{}\n")...) }},
		{name: "unknown-field", mutate: func(data []byte) []byte {
			return bytes.Replace(data, []byte("{\"storeFormatVersion\""), []byte("{\"unknown\":1,\"storeFormatVersion\""), 1)
		}},
		{name: "duplicate-field", mutate: func(data []byte) []byte {
			return bytes.Replace(data, []byte("{\"storeFormatVersion\":2"), []byte("{\"storeFormatVersion\":2,\"storeFormatVersion\":2"), 1)
		}},
	}
	for _, test := range rawCases {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "raw-corrupt-v2.json")
			data := test.mutate(mustReadFile(t, validPath))
			if err := os.WriteFile(path, data, 0o600); err != nil {
				t.Fatal(err)
			}
			assertV2OpenDoesNotRewrite(t, path)
		})
	}
}

func TestFixedStoreV2AbsentLifecycleRejectsInventedEffectAndInstance(t *testing.T) {
	state, record := stateAndLifecycleRecord(t)
	record = lifecycleRecordForV2State(t, state, record, lifecyclestate.StateObserved)
	absentPath := newV1PathFromState(t, state, nil)
	presentPath := newV1PathFromState(t, state, []lifecyclestate.Record{record})
	mustMigrateV2(t, absentPath)
	mustMigrateV2(t, presentPath)
	absent := readV2Envelope(t, absentPath)
	present := readV2Envelope(t, presentPath)
	absent.Indexes.Effects = append([]effectIndexEntryDisk(nil), present.Indexes.Effects...)
	absent.Indexes.Instances = append([]instanceIndexEntryDisk(nil), present.Indexes.Instances...)
	locallyValid, err := archiveIndexesFromDisk(absent.Indexes)
	if err != nil {
		t.Fatalf("invented local index should be structurally valid before full reconstruction: %v", err)
	}
	absent.IndexDigests = locallyValid.Digests()
	absent.CombinedIndexDigest = locallyValid.CombinedDigest()
	path := filepath.Join(t.TempDir(), "invented-lifecycle-state-v2.json")
	writeV2Envelope(t, path, absent)
	assertV2OpenDoesNotRewrite(t, path)
}

func TestFixedStoreV2MigrationProcessDeathOldOrNewWorld(t *testing.T) {
	for _, point := range []MigrationFault{FaultV2MigrationBeforeRename, FaultV2MigrationAfterRename} {
		t.Run(string(point), func(t *testing.T) {
			state, _ := stateAndLifecycleRecord(t)
			path := newV1PathFromState(t, state, nil)
			before := mustReadFile(t, path)
			command := exec.Command(os.Args[0], "-test.run=TestFixedStoreV2MigrationProcessDeathHelper") //nolint:gosec // G204: exact current test binary and fixed test selector.
			command.Env = append(os.Environ(), "CAPSULE_F2_CHILD_PATH="+path, "CAPSULE_F2_CHILD_POINT="+string(point))
			if err := command.Run(); err == nil {
				t.Fatal("process-death child unexpectedly returned success")
			}
			switch point {
			case FaultV2MigrationBeforeRename:
				if got := mustReadFile(t, path); !bytes.Equal(got, before) {
					t.Fatal("pre-rename process death changed v1 bytes")
				}
				if _, err := OpenFixedFileStoreV1(path); err != nil {
					t.Fatalf("pre-rename process death did not leave complete v1: %v", err)
				}
			case FaultV2MigrationAfterRename:
				if _, err := OpenFixedFileStoreV2(path); err != nil {
					t.Fatalf("post-rename process death did not leave complete v2: %v", err)
				}
			}
		})
	}
}

func TestFixedStoreV2MigrationProcessDeathHelper(t *testing.T) {
	path := os.Getenv("CAPSULE_F2_CHILD_PATH")
	point := MigrationFault(os.Getenv("CAPSULE_F2_CHILD_POINT"))
	if path == "" || point == "" {
		return
	}
	_, _ = MigrateFixedFileStoreV1ToV2(
		context.Background(), path,
		V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}, Faults: exitMigrationFault{point: point}},
	)
	t.Fatal("child migration returned without injected process death")
}

type exitMigrationFault struct{ point MigrationFault }

func (fault exitMigrationFault) FailMigrationAt(point MigrationFault) error {
	if point == fault.point {
		os.Exit(97)
	}
	return nil
}

type stagedMigrationLock struct {
	checks  int
	failAt  int
	onCheck func(int)
}

func (lock *stagedMigrationLock) CheckOfflineMigrationLock(context.Context) error {
	lock.checks++
	if lock.onCheck != nil {
		lock.onCheck(lock.checks)
	}
	if lock.checks == lock.failAt {
		return errors.New("injected owner loss")
	}
	return nil
}

func newV1PathFromState(t *testing.T, state installationState, records []lifecyclestate.Record) string {
	t.Helper()
	if err := validateV1State(state, records, lifecycleSetDigest(records)); err != nil {
		t.Fatalf("invalid v1 test source: %v", err)
	}
	return newV1PathFromStateUnchecked(t, state, records)
}

func newV1PathFromStateUnchecked(t *testing.T, state installationState, records []lifecyclestate.Record) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "state-v1.json")
	writeV1Envelope(t, path, encodedEnvelopeV1(state, records))
	return path
}

func mustMigrateV2(t *testing.T, path string) *FixedFileStoreV2 {
	t.Helper()
	store, err := MigrateFixedFileStoreV1ToV2(
		context.Background(), path,
		V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}},
	)
	if err != nil {
		t.Fatalf("migrate v2: %v", err)
	}
	return store
}

func mustV2Snapshot(t *testing.T, store *FixedFileStoreV2) fixedStoreV2Snapshot {
	t.Helper()
	snapshot, err := store.snapshotV2(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func readV2Envelope(t *testing.T, path string) diskEnvelopeV2 {
	t.Helper()
	var envelope diskEnvelopeV2
	if err := decodeOneClosedJSON(mustReadFile(t, path), &envelope); err != nil {
		t.Fatal(err)
	}
	return envelope
}

func cloneV2Envelope(t *testing.T, envelope diskEnvelopeV2) diskEnvelopeV2 {
	t.Helper()
	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	var cloned diskEnvelopeV2
	if err := json.Unmarshal(data, &cloned); err != nil {
		t.Fatal(err)
	}
	return cloned
}

func writeV2Envelope(t *testing.T, path string, envelope diskEnvelopeV2) {
	t.Helper()
	data, err := encodeEnvelopeV2(envelope)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func assertV2OpenDoesNotRewrite(t *testing.T, path string) {
	t.Helper()
	before := mustReadFile(t, path)
	if _, err := OpenFixedFileStoreV2(path); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
		t.Fatalf("corrupt v2 open result = %v", err)
	}
	if after := mustReadFile(t, path); !bytes.Equal(after, before) {
		t.Fatal("corrupt v2 refusal rewrote evidence")
	}
	if report, err := VerifyFixedFileStoreV2(path); err == nil || report != (FixedFileStoreV2VerificationReport{}) {
		t.Fatalf("corrupt v2 verifier returned report %#v, %v", report, err)
	}
}

func lifecycleRecordForV2State(
	t *testing.T,
	authority installationState,
	template lifecyclestate.Record,
	state lifecyclestate.LifecycleState,
) lifecyclestate.Record {
	t.Helper()
	view := template.View()
	view.State = state
	view.OperationSequence = 0
	view.Operation = lifecyclestate.OperationNone
	view.EffectID = lifecyclestate.EffectID{}
	view.EffectStatus = lifecyclestate.EffectNone
	view.Instance = lifecyclestate.BackendInstanceIdentity{}
	view.CleanupRequired = true
	view.LastConfirmedCheckpoint = lifecyclestate.CheckpointNone
	view.FirstFailure = lifecyclestate.FailureNone
	view.FailureOperation = lifecyclestate.OperationNone
	view.LastReconciliation = lifecyclestate.ReconciliationNone
	view.AutomaticRecoveryCount = 0
	view.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{}
	view.RecoveryFence = lifecyclestate.RecoveryFenceNone
	view.AutomaticRecoveryExhausted = false
	view.TerminalAt = lifecyclestate.OptionalUnixSeconds{}
	setEffect := func(operation lifecyclestate.Operation, status lifecyclestate.EffectStatus, checkpoint lifecyclestate.Checkpoint, instance bool) {
		view.OperationSequence = 1
		view.Operation = operation
		view.EffectStatus = status
		view.LastConfirmedCheckpoint = checkpoint
		attemptID := view.Bindings.View().AttemptID
		identifier, err := lifecyclestate.NewDomainIdentifier(lifecyclestate.DomainEffectID, attemptID[:])
		if err != nil {
			t.Fatal(err)
		}
		view.EffectID, err = lifecyclestate.NewEffectID(identifier)
		if err != nil {
			t.Fatal(err)
		}
		if instance {
			view.Instance, err = lifecyclestate.NewBackendInstanceIdentity(lifecyclestate.BackendInstanceFake, attemptID[:])
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	switch state {
	case lifecyclestate.StatePreparePending:
	case lifecyclestate.StatePrepareIntent:
		setEffect(lifecyclestate.OperationPrepare, lifecyclestate.EffectIntent, lifecyclestate.CheckpointNone, false)
	case lifecyclestate.StatePrepared:
		setEffect(lifecyclestate.OperationPrepare, lifecyclestate.EffectConfirmed, lifecyclestate.CheckpointPrepare, false)
	case lifecyclestate.StateCreateIntent:
		setEffect(lifecyclestate.OperationCreate, lifecyclestate.EffectIntent, lifecyclestate.CheckpointPrepare, false)
	case lifecyclestate.StateCreated:
		setEffect(lifecyclestate.OperationCreate, lifecyclestate.EffectConfirmed, lifecyclestate.CheckpointCreate, true)
	case lifecyclestate.StateStartIntent:
		setEffect(lifecyclestate.OperationStart, lifecyclestate.EffectIntent, lifecyclestate.CheckpointCreate, true)
	case lifecyclestate.StateStarted:
		setEffect(lifecyclestate.OperationStart, lifecyclestate.EffectConfirmed, lifecyclestate.CheckpointStart, true)
	case lifecyclestate.StateObserveIntent:
		setEffect(lifecyclestate.OperationObserve, lifecyclestate.EffectIntent, lifecyclestate.CheckpointStart, true)
	case lifecyclestate.StateObserved:
		setEffect(lifecyclestate.OperationObserve, lifecyclestate.EffectConfirmed, lifecyclestate.CheckpointObserve, true)
	case lifecyclestate.StateStopIntent:
		setEffect(lifecyclestate.OperationStop, lifecyclestate.EffectIntent, lifecyclestate.CheckpointObserve, true)
	case lifecyclestate.StateStopped:
		setEffect(lifecyclestate.OperationStop, lifecyclestate.EffectConfirmed, lifecyclestate.CheckpointStop, true)
	case lifecyclestate.StateDestroyIntent:
		setEffect(lifecyclestate.OperationDestroy, lifecyclestate.EffectIntent, lifecyclestate.CheckpointStop, true)
	case lifecyclestate.StateDestroyConfirmed:
		setEffect(lifecyclestate.OperationDestroy, lifecyclestate.EffectConfirmed, lifecyclestate.CheckpointDestroy, true)
	case lifecyclestate.StateDestroyed:
		view.CleanupRequired = false
		view.LastReconciliation = lifecyclestate.ReconciliationAuthoritativelyAbsent
		view.AutomaticRecoveryCount = 1
		view.TerminalAt = lifecyclestate.OptionalUnixSeconds{Present: true, Value: view.LastTransitionAt}
	case lifecyclestate.StateUnresolved:
		setEffect(lifecyclestate.OperationObserve, lifecyclestate.EffectIndeterminate, lifecyclestate.CheckpointStart, true)
		view.FirstFailure = lifecyclestate.FailureCleanupUnresolved
		view.FailureOperation = lifecyclestate.OperationObserve
		view.LastReconciliation = lifecyclestate.ReconciliationUnknown
		view.AutomaticRecoveryCount = 1
		view.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{Present: true, Value: authority.TimeHighWaterUnixSeconds + 1}
		view.RecoveryFence = lifecyclestate.RecoveryFenceReconcileUnknown
	case lifecyclestate.StateQuarantined:
		setEffect(lifecyclestate.OperationStart, lifecyclestate.EffectIndeterminate, lifecyclestate.CheckpointCreate, true)
		view.FirstFailure = lifecyclestate.FailureBinding
		view.FailureOperation = lifecyclestate.OperationStart
		view.LastReconciliation = lifecyclestate.ReconciliationIdentityMismatch
		view.AutomaticRecoveryCount = 1
		view.RecoveryFence = lifecyclestate.RecoveryFenceIdentityMismatch
	default:
		t.Fatalf("unknown lifecycle state %q", state)
	}
	record, err := lifecyclestate.NewRecord(view, authority.TimeHighWaterUnixSeconds)
	if err != nil {
		t.Fatalf("construct lifecycle state %s: %v", state, err)
	}
	return record
}

func TestFixedStoreV2TestHelpersRemainDeterministic(t *testing.T) {
	state, record := stateAndLifecycleRecord(t)
	left := lifecycleRecordForV2State(t, state, record, lifecyclestate.StateObserved)
	right := lifecycleRecordForV2State(t, state, record, lifecyclestate.StateObserved)
	if !reflect.DeepEqual(lifecycleRecordToDisk(left), lifecycleRecordToDisk(right)) {
		t.Fatal("v2 lifecycle fixture construction is nondeterministic")
	}
	if sha256.Sum256(mustJSON(t, lifecycleRecordToDisk(left))) != sha256.Sum256(mustJSON(t, lifecycleRecordToDisk(right))) {
		t.Fatal("v2 lifecycle fixture bytes are nondeterministic")
	}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
