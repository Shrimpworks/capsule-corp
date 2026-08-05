package registrationstate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

type migrationLockStub struct {
	held   bool
	checks int
}

func (lock *migrationLockStub) CheckOfflineMigrationLock(context.Context) error {
	lock.checks++
	if !lock.held {
		return errors.New("lock not held")
	}
	return nil
}

type migrationFaultStub struct {
	point MigrationFault
	hits  map[MigrationFault]int
}

func (fault *migrationFaultStub) FailMigrationAt(point MigrationFault) error {
	if fault.hits == nil {
		fault.hits = make(map[MigrationFault]int)
	}
	fault.hits[point]++
	if point == fault.point {
		return errors.New("injected migration fault")
	}
	return nil
}

func TestFixedStoreV0ToV1MigrationAndDowngradeRefusal(t *testing.T) {
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
	v0State, _ := harness.store.snapshot(context.Background())
	v0Bytes := mustReadFile(t, harness.path)

	if _, err := OpenFixedFileStoreV1(harness.path); err == nil {
		t.Fatal("v1 open automatically accepted v0")
	}
	if got := mustReadFile(t, harness.path); !bytes.Equal(got, v0Bytes) {
		t.Fatal("refused automatic migration rewrote v0")
	}

	lock := &migrationLockStub{held: true}
	migrated, err := MigrateFixedFileStoreV0ToV1(
		context.Background(), harness.path, V0ToV1MigrationOptions{Lock: lock},
	)
	if err != nil {
		t.Fatalf("migrate v0 to v1: %v", err)
	}
	if lock.checks != 3 {
		t.Fatalf("migration lock checks = %d, want 3", lock.checks)
	}
	snapshot, err := migrated.snapshotV1(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.SourceFormatVersion != SupervisorStoreFormatV0 || len(snapshot.Lifecycles) != 0 ||
		snapshot.LifecycleSetDigest != lifecycleSetDigest(nil) {
		t.Fatalf("migrated v1 metadata = %#v", snapshot)
	}
	if !authorityStatesEqual(v0State, snapshot.State) {
		t.Fatal("migration changed v0 authority state")
	}

	reopened, err := OpenFixedFileStoreV1(harness.path)
	if err != nil {
		t.Fatalf("reopen v1: %v", err)
	}
	reopenedSnapshot, _ := reopened.snapshotV1(context.Background())
	if !authorityStatesEqual(snapshot.State, reopenedSnapshot.State) ||
		reopenedSnapshot.LifecycleSetDigest != snapshot.LifecycleSetDigest {
		t.Fatal("reopen changed migrated state")
	}

	v1Bytes := mustReadFile(t, harness.path)
	if _, err := NewFixedFileStore(harness.path, InitialState{}); err == nil ||
		!strings.Contains(err.Error(), "unsupported fixed registration store version") {
		t.Fatalf("v0 opener downgrade result = %v", err)
	}
	if got := mustReadFile(t, harness.path); !bytes.Equal(got, v1Bytes) {
		t.Fatal("v0 downgrade refusal rewrote v1")
	}
	if _, err := MigrateFixedFileStoreV0ToV1(
		context.Background(), harness.path,
		V0ToV1MigrationOptions{Lock: &migrationLockStub{held: true}},
	); err == nil {
		t.Fatal("repeat migration accepted v1 as v0")
	}
	if got := mustReadFile(t, harness.path); !bytes.Equal(got, v1Bytes) {
		t.Fatal("repeat migration rewrote v1")
	}
}

func TestFixedStoreV1MigrationFaultOracles(t *testing.T) {
	for iteration := 0; iteration < 12; iteration++ {
		t.Run("pre-rename-"+string(rune('a'+iteration)), func(t *testing.T) {
			path := newV0StorePath(t)
			before := mustReadFile(t, path)
			fault := &migrationFaultStub{point: FaultMigrationBeforeRename}
			_, err := MigrateFixedFileStoreV0ToV1(
				context.Background(), path,
				V0ToV1MigrationOptions{Lock: &migrationLockStub{held: true}, Faults: fault},
			)
			if err == nil || errors.Is(err, ErrMigrationOutcomeIndeterminate) {
				t.Fatalf("pre-rename error = %v", err)
			}
			if after := mustReadFile(t, path); !bytes.Equal(after, before) {
				t.Fatal("confirmed pre-rename failure changed v0 bytes")
			}
			if _, err := NewFixedFileStore(path, InitialState{}); err != nil {
				t.Fatalf("v0 did not reopen after confirmed failure: %v", err)
			}
			if _, err := OpenFixedFileStoreV1(path); err == nil {
				t.Fatal("confirmed failure produced v1")
			}
		})

		t.Run("post-rename-"+string(rune('a'+iteration)), func(t *testing.T) {
			path := newV0StorePath(t)
			fault := &migrationFaultStub{point: FaultMigrationAfterRename}
			_, err := MigrateFixedFileStoreV0ToV1(
				context.Background(), path,
				V0ToV1MigrationOptions{Lock: &migrationLockStub{held: true}, Faults: fault},
			)
			if !errors.Is(err, ErrMigrationOutcomeIndeterminate) {
				t.Fatalf("post-rename error = %v", err)
			}
			if _, err := OpenFixedFileStoreV1(path); err != nil {
				t.Fatalf("indeterminate post-rename v1 reopen: %v", err)
			}
			if _, err := NewFixedFileStore(path, InitialState{}); err == nil {
				t.Fatal("old opener accepted post-rename v1")
			}
		})
	}
}

func TestFixedStoreV1RefusesMissingAndUnlockedMigration(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.json")
	if _, err := OpenFixedFileStoreV1(missing); err == nil {
		t.Fatal("v1 open created a missing store")
	}
	if _, err := MigrateFixedFileStoreV0ToV1(
		context.Background(), missing,
		V0ToV1MigrationOptions{Lock: &migrationLockStub{held: true}},
	); err == nil {
		t.Fatal("migration created a missing store")
	}
	if _, err := os.Lstat(missing); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing fallback created state: %v", err)
	}

	path := newV0StorePath(t)
	before := mustReadFile(t, path)
	for _, options := range []V0ToV1MigrationOptions{
		{},
		{Lock: &migrationLockStub{held: false}},
	} {
		if _, err := MigrateFixedFileStoreV0ToV1(context.Background(), path, options); err == nil {
			t.Fatal("migration proceeded without a held lock")
		}
		if after := mustReadFile(t, path); !bytes.Equal(after, before) {
			t.Fatal("unlocked migration rewrote v0")
		}
	}
}

func TestFixedStoreV1MigrationFullyValidatesV0WithoutRewrite(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func([]byte) []byte
	}{
		{name: "trailing", mutate: func(data []byte) []byte {
			return append(data, []byte("{}\n")...)
		}},
		{name: "missing version", mutate: func(data []byte) []byte {
			return bytes.Replace(data, []byte(`"storeFormatVersion":0,`), nil, 1)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := newV0StorePath(t)
			data := test.mutate(mustReadFile(t, path))
			if err := os.WriteFile(path, data, 0o600); err != nil {
				t.Fatal(err)
			}
			before := mustReadFile(t, path)
			if _, err := MigrateFixedFileStoreV0ToV1(
				context.Background(), path,
				V0ToV1MigrationOptions{Lock: &migrationLockStub{held: true}},
			); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
				t.Fatalf("corrupt v0 migration error = %v", err)
			}
			if after := mustReadFile(t, path); !bytes.Equal(after, before) {
				t.Fatal("failed v0 validation rewrote source evidence")
			}
		})
	}
}

func TestFixedStoreV1ReopenValidatesLifecycleCrossLinksAndCopies(t *testing.T) {
	state, record := stateAndLifecycleRecord(t)
	path := filepath.Join(t.TempDir(), "state-v1.json")
	writeV1Envelope(t, path, encodedEnvelopeV1(state, []lifecyclestate.Record{record}))

	store, err := OpenFixedFileStoreV1(path)
	if err != nil {
		t.Fatalf("open valid lifecycle v1: %v", err)
	}
	first, _ := store.snapshotV1(context.Background())
	if len(first.Lifecycles) != 1 || first.LifecycleSetDigest != lifecycleSetDigest(first.Lifecycles) {
		t.Fatal("valid lifecycle collection did not round trip")
	}
	view := first.Lifecycles[0].View()
	bindingView := view.Bindings.View()
	bindingView.ProfileReviewAttestationDigests[0][0] ^= 0xff
	second, _ := store.snapshotV1(context.Background())
	if second.Lifecycles[0].Bindings().View().ProfileReviewAttestationDigests[0] ==
		bindingView.ProfileReviewAttestationDigests[0] {
		t.Fatal("snapshot lifecycle projection aliased caller-owned bytes")
	}
}

func TestFixedStoreV1CorruptionRequiresRepairWithoutRewrite(t *testing.T) {
	state, record := stateAndLifecycleRecord(t)
	valid := encodedEnvelopeV1(state, []lifecyclestate.Record{record})
	tests := []struct {
		name      string
		wantError string
		mutate    func(*diskEnvelopeV1)
	}{
		{name: "unsupported store version", mutate: func(envelope *diskEnvelopeV1) {
			value := uint64(2)
			envelope.StoreFormatVersion = &value
		}},
		{name: "missing store version", mutate: func(envelope *diskEnvelopeV1) {
			envelope.StoreFormatVersion = nil
		}},
		{name: "unsupported source version", mutate: func(envelope *diskEnvelopeV1) {
			value := uint64(1)
			envelope.MigrationSourceVersion = &value
		}},
		{name: "missing lifecycle collection", mutate: func(envelope *diskEnvelopeV1) {
			envelope.Lifecycles = nil
		}},
		{name: "lifecycle set digest", mutate: func(envelope *diskEnvelopeV1) {
			envelope.LifecycleSetDigest[0] ^= 0xff
		}},
		{name: "approval set digest", wantError: "fixed approval/attempt set digest mismatch", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.ApprovalSetDigest[0] ^= 0xff
		}},
		{name: "attempt set digest", wantError: "fixed approval/attempt set digest mismatch", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.AttemptSetDigest[0] ^= 0xff
		}},
		{name: "missing registration set digest", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.RegistrationSetDigest = [32]byte{}
		}},
		{name: "record format version", mutate: func(envelope *diskEnvelopeV1) {
			envelope.Lifecycles[0].FormatVersion++
		}},
		{name: "partial backend binding", mutate: func(envelope *diskEnvelopeV1) {
			envelope.Lifecycles[0].Bindings.Backend = lifecyclestate.BackendBindingView{}
		}},
		{name: "immutable binding digest", mutate: func(envelope *diskEnvelopeV1) {
			envelope.Lifecycles[0].ImmutableBindingDigest[0] ^= 0xff
		}},
		{name: "lifecycle attempt cross-link", mutate: func(envelope *diskEnvelopeV1) {
			envelope.State.Attempts[0].AttemptID[0] ^= 0xff
			attemptDigest, err := attemptSetDigest(envelope.State.Attempts)
			if err != nil {
				t.Fatalf("attempt set digest: %v", err)
			}
			envelope.State.AttemptSetDigest = attemptDigest
		}},
		{name: "timestamp above high water", mutate: func(envelope *diskEnvelopeV1) {
			envelope.Lifecycles[0].LastTransitionAt = envelope.State.TimeHighWaterUnixSeconds + 1
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			envelope := cloneDiskEnvelopeV1(t, valid)
			test.mutate(&envelope)
			path := filepath.Join(t.TempDir(), "corrupt-v1.json")
			writeV1Envelope(t, path, envelope)
			before := mustReadFile(t, path)
			_, err := OpenFixedFileStoreV1(path)
			if err == nil || !errors.Is(err, ErrStoreRepairRequired) {
				t.Fatalf("corrupt open error = %v, want repair required", err)
			}
			if test.wantError != "" && !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("corrupt open error = %v, want %q", err, test.wantError)
			}
			if after := mustReadFile(t, path); !bytes.Equal(after, before) {
				t.Fatal("corrupt v1 open rewrote evidence")
			}
		})
	}

	t.Run("truncated", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "truncated-v1.json")
		writeV1Envelope(t, path, valid)
		data := mustReadFile(t, path)
		if err := os.WriteFile(path, data[:len(data)/2], 0o600); err != nil {
			t.Fatal(err)
		}
		assertV1OpenDoesNotRewrite(t, path)
	})

	t.Run("trailing data", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "trailing-v1.json")
		writeV1Envelope(t, path, valid)
		data := append(mustReadFile(t, path), []byte("{}\n")...)
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
		assertV1OpenDoesNotRewrite(t, path)
	})

	t.Run("unknown field", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "unknown-v1.json")
		writeV1Envelope(t, path, valid)
		data := mustReadFile(t, path)
		data = bytes.Replace(data, []byte(`"storeFormatVersion":1`), []byte(`"storeFormatVersion":1,"unexpected":true`), 1)
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
		assertV1OpenDoesNotRewrite(t, path)
	})

	t.Run("duplicate field", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "duplicate-v1.json")
		writeV1Envelope(t, path, valid)
		data := mustReadFile(t, path)
		data = bytes.Replace(data, []byte(`"storeFormatVersion":1`), []byte(`"storeFormatVersion":1,"storeFormatVersion":1`), 1)
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
		assertV1OpenDoesNotRewrite(t, path)
	})

	t.Run("lifecycle retained capacity", func(t *testing.T) {
		envelope := cloneDiskEnvelopeV1(t, valid)
		envelope.Lifecycles = make([]lifecycleRecordDisk, MaxRetainedLifecycleRecords+1)
		path := filepath.Join(t.TempDir(), "capacity-v1.json")
		writeV1Envelope(t, path, envelope)
		assertV1OpenDoesNotRewrite(t, path)
	})
}

func TestFixedStoreV1RejectsValidButWrongLifecycleCrossLinks(t *testing.T) {
	state, record := stateAndLifecycleRecord(t)
	tests := []struct {
		name   string
		mutate func(*lifecyclestate.ImmutableBindingsView)
	}{
		{name: "attempt", mutate: func(view *lifecyclestate.ImmutableBindingsView) {
			view.AttemptID[0] ^= 0xff
		}},
		{name: "approval", mutate: func(view *lifecyclestate.ImmutableBindingsView) {
			view.ApprovalID[0] ^= 0xff
		}},
		{name: "registration", mutate: func(view *lifecyclestate.ImmutableBindingsView) {
			view.RegistrationID[0] ^= 0xff
		}},
		{name: "plan role", mutate: func(view *lifecyclestate.ImmutableBindingsView) {
			view.SourceManifestDigest[0] ^= 0xff
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wrong := lifecycleRecordWithMutatedBindings(t, state, record, test.mutate)
			path := filepath.Join(t.TempDir(), "wrong-cross-link-v1.json")
			writeV1Envelope(t, path, encodedEnvelopeV1(state, []lifecyclestate.Record{wrong}))
			assertV1OpenDoesNotRewrite(t, path)
		})
	}
}

func TestFixedStoreV1RejectsUnsafeFileShapeWithoutRewrite(t *testing.T) {
	state, record := stateAndLifecycleRecord(t)
	path := filepath.Join(t.TempDir(), "permissions-v1.json")
	writeV1Envelope(t, path, encodedEnvelopeV1(state, []lifecyclestate.Record{record}))
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	before := mustReadFile(t, path)
	if _, err := OpenFixedFileStoreV1(path); err == nil {
		t.Fatal("v1 open accepted unsafe file permissions")
	}
	if after := mustReadFile(t, path); !bytes.Equal(after, before) {
		t.Fatal("unsafe file refusal rewrote evidence")
	}
}

func TestFixedStoreV1CapacityAndFormatIdentities(t *testing.T) {
	state := ordinaryInitialStateForV1()
	tooMany := make([]lifecyclestate.Record, MaxRetainedLifecycleRecords+1)
	if err := validateV1State(state, tooMany, lifecycleSetDigest(tooMany)); err == nil ||
		!strings.Contains(err.Error(), "retained capacity") {
		t.Fatalf("retained cap-plus-one result = %v", err)
	}
	emptyDigest := lifecycleSetDigest(nil)
	emptyHex := hex.EncodeToString(emptyDigest[:])
	const wantEmpty = "5b65d09d0b92bc6c4629fc2283aefb51a1ac3533f2ee81fed968db438e790ad4"
	if emptyHex != wantEmpty {
		t.Fatalf("empty lifecycle-set digest = %s, want %s", emptyHex, wantEmpty)
	}
}

func TestFixedStoreV1ExactActiveCapacityReleasesOnlyAfterDurableDestroy(t *testing.T) {
	exactState, exactRecords := generatedLifecyclePopulation(
		t, MaxActiveLifecycleRecords, lifecyclestate.StatePreparePending,
	)
	if got := activeLifecycleWork(exactState, nil); got != MaxActiveLifecycleRecords {
		t.Fatalf("exact active population = %d, want %d", got, MaxActiveLifecycleRecords)
	}
	if err := validateV1State(exactState, nil, lifecycleSetDigest(nil)); err != nil {
		t.Fatalf("exact active population rejected: %v", err)
	}
	if len(exactRecords) != MaxActiveLifecycleRecords {
		t.Fatal("generated exact active population lost records")
	}

	states := []lifecyclestate.LifecycleState{
		lifecyclestate.StateObserved,
		lifecyclestate.StateStopped,
		lifecyclestate.StateDestroyConfirmed,
		lifecyclestate.StateUnresolved,
		lifecyclestate.StateQuarantined,
	}
	for _, state := range states {
		t.Run(string(state)+"-does-not-release", func(t *testing.T) {
			candidate, records := generatedLifecyclePopulation(t, MaxActiveLifecycleRecords+1, state)
			one := records[:1]
			if got := activeLifecycleWork(candidate, one); got != MaxActiveLifecycleRecords+1 {
				t.Fatalf("%s active population = %d, want %d", state, got, MaxActiveLifecycleRecords+1)
			}
			if err := validateV1State(candidate, one, lifecycleSetDigest(one)); err == nil ||
				!strings.Contains(err.Error(), "active capacity") {
				t.Fatalf("%s cap-plus-one result = %v", state, err)
			}
		})
	}

	releasedState, releasedRecords := generatedLifecyclePopulation(
		t, MaxActiveLifecycleRecords+1, lifecyclestate.StateDestroyed,
	)
	oneDestroyed := releasedRecords[:1]
	view := oneDestroyed[0].View()
	if view.CleanupRequired || view.LastReconciliation != lifecyclestate.ReconciliationAuthoritativelyAbsent {
		t.Fatalf("capacity-releasing record = %#v", view)
	}
	for _, test := range []struct {
		name   string
		mutate func(*lifecyclestate.RecordView)
	}{
		{name: "cleanup-still-required", mutate: func(candidate *lifecyclestate.RecordView) {
			candidate.CleanupRequired = true
		}},
		{name: "absence-not-authoritative", mutate: func(candidate *lifecyclestate.RecordView) {
			candidate.LastReconciliation = lifecyclestate.ReconciliationUnknown
		}},
	} {
		t.Run(test.name+"-cannot-be-destroyed", func(t *testing.T) {
			candidate := view
			test.mutate(&candidate)
			if _, err := lifecyclestate.NewRecord(candidate, releasedState.TimeHighWaterUnixSeconds); err == nil {
				t.Fatal("invalid destroyed disposition passed durable record validation")
			}
		})
	}
	if got := activeLifecycleWork(releasedState, oneDestroyed); got != MaxActiveLifecycleRecords {
		t.Fatalf("destroyed capacity release = %d, want %d", got, MaxActiveLifecycleRecords)
	}
	if err := validateV1State(releasedState, oneDestroyed, lifecycleSetDigest(oneDestroyed)); err != nil {
		t.Fatalf("destroyed capacity release rejected: %v", err)
	}

	path := filepath.Join(t.TempDir(), "active-capacity-v1.json")
	writeV1Envelope(t, path, encodedEnvelopeV1(releasedState, oneDestroyed))
	reopened, err := OpenFixedFileStoreV1(path)
	if err != nil {
		t.Fatalf("reopen destroyed capacity release: %v", err)
	}
	recovery, err := reopened.RecoveryAttemptIDs(context.Background())
	if err != nil || len(recovery) != MaxActiveLifecycleRecords {
		t.Fatalf("recovery after one durable destroy = %d, %v", len(recovery), err)
	}

	capPlusOnePath := filepath.Join(t.TempDir(), "active-capacity-plus-one-v1.json")
	writeV1Envelope(t, capPlusOnePath, encodedEnvelopeV1(releasedState, nil))
	before := mustReadFile(t, capPlusOnePath)
	if _, err := OpenFixedFileStoreV1(capPlusOnePath); err == nil ||
		!strings.Contains(err.Error(), "active capacity") {
		t.Fatalf("active cap-plus-one open = %v", err)
	}
	if after := mustReadFile(t, capPlusOnePath); !bytes.Equal(after, before) {
		t.Fatal("active cap-plus-one refusal evicted or rewrote retained state")
	}
}

func TestFixedStoreV1ExactRetainedLifecycleCapacityNeverEvicts(t *testing.T) {
	state, records := generatedLifecyclePopulation(
		t, MaxRetainedLifecycleRecords, lifecyclestate.StateDestroyed,
	)
	if len(records) != MaxRetainedLifecycleRecords {
		t.Fatalf("retained lifecycle count = %d", len(records))
	}
	if err := validateV1State(state, records, lifecycleSetDigest(records)); err != nil {
		t.Fatalf("exact retained lifecycle population rejected: %v", err)
	}

	path := filepath.Join(t.TempDir(), "retained-lifecycle-capacity-v1.json")
	envelope := encodedEnvelopeV1(state, records)
	encoded, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	if int64(len(encoded)) > MaxSupervisorStateV1Bytes {
		t.Fatalf("exact retained lifecycle population is %d bytes, exceeds fixed bound %d", len(encoded), MaxSupervisorStateV1Bytes)
	}
	t.Logf("exact retained lifecycle population: %d records, %d encoded bytes", len(records), len(encoded))
	writeV1Envelope(t, path, envelope)
	reopened, err := OpenFixedFileStoreV1(path)
	if err != nil {
		t.Fatalf("reopen exact retained lifecycle population: %v", err)
	}
	recovery, err := reopened.RecoveryAttemptIDs(context.Background())
	if err != nil || len(recovery) != 0 {
		t.Fatalf("terminal retained recovery set = %d, %v", len(recovery), err)
	}

	tooMany := append(append([]lifecyclestate.Record(nil), records...), lifecyclestate.Record{})
	if err := validateV1State(state, tooMany, lifecycleSetDigest(tooMany)); err == nil ||
		!strings.Contains(err.Error(), "retained capacity") {
		t.Fatalf("retained cap-plus-one result = %v", err)
	}
	capPlusOnePath := filepath.Join(t.TempDir(), "retained-lifecycle-capacity-plus-one-v1.json")
	writeV1Envelope(t, capPlusOnePath, encodedEnvelopeV1(state, tooMany))
	before := mustReadFile(t, capPlusOnePath)
	if _, err := OpenFixedFileStoreV1(capPlusOnePath); err == nil ||
		!strings.Contains(err.Error(), "retained capacity") {
		t.Fatalf("retained cap-plus-one open = %v", err)
	}
	if after := mustReadFile(t, capPlusOnePath); !bytes.Equal(after, before) {
		t.Fatal("retained cap-plus-one refusal evicted or rewrote retained state")
	}
}

func generatedLifecyclePopulation(
	t *testing.T,
	count int,
	state lifecyclestate.LifecycleState,
) (installationState, []lifecyclestate.Record) {
	t.Helper()
	templateState, templateRecord := stateAndLifecycleRecord(t)
	templateApproval := templateState.Approvals[0]
	templateAttempt := templateState.Attempts[0]
	templateBindings := templateRecord.Bindings().View()

	result := cloneState(templateState)
	result.Approvals = make([]approvalattempt.ApprovalRecord, count)
	result.Attempts = make([]approvalattempt.ExecutionAttempt, count)
	records := make([]lifecyclestate.Record, count)
	for index := range count {
		ordinal := uint16(index + 1)
		approval := approvalattempt.CloneApprovalRecord(templateApproval)
		approval.ApprovalID = templateApproval.ApprovalID
		approval.AttemptNonce = templateApproval.AttemptNonce
		binary.BigEndian.PutUint16(approval.ApprovalID[14:], ordinal)
		binary.BigEndian.PutUint16(approval.AttemptNonce[14:], ordinal)
		binary.BigEndian.PutUint16(approval.ExactEnvelopeBytes[len(approval.ExactEnvelopeBytes)-2:], ordinal)
		binary.BigEndian.PutUint16(approval.ExactPayloadBytes[len(approval.ExactPayloadBytes)-2:], ordinal)
		approval.EnvelopeDigest = approvalattempt.ApprovalEnvelopeDigest(sha256.Sum256(approval.ExactEnvelopeBytes))
		approval.PayloadDigest = approvalattempt.ApprovalPayloadDigest(sha256.Sum256(approval.ExactPayloadBytes))
		approval.ConsumedAttemptID = templateAttempt.AttemptID
		binary.BigEndian.PutUint16(approval.ConsumedAttemptID[14:], ordinal)

		attempt := templateAttempt
		attempt.AttemptID = approval.ConsumedAttemptID
		attempt.ApprovalID = approval.ApprovalID
		attempt.AttemptNonce = approval.AttemptNonce
		attempt.ApprovalPayloadDigest = approval.PayloadDigest
		result.Approvals[index] = approval
		result.Attempts[index] = attempt

		bindingsView := templateBindings
		bindingsView.ProfileReviewAttestationDigests = append(
			[]v0candidate.ProfileReviewAttestationDigest(nil),
			templateBindings.ProfileReviewAttestationDigests...,
		)
		bindingsView.AttemptID = attempt.AttemptID
		bindingsView.ApprovalID = attempt.ApprovalID
		bindingsView.AttemptNonce = attempt.AttemptNonce
		bindingsView.ApprovalPayloadDigest = attempt.ApprovalPayloadDigest
		bindings, err := lifecyclestate.NewImmutableBindings(bindingsView)
		if err != nil {
			t.Fatalf("generated lifecycle bindings %d: %v", index, err)
		}
		records[index] = generatedLifecycleRecord(t, result.TimeHighWaterUnixSeconds, bindings, state)
	}
	approvalDigest, err := approvalSetDigest(result.Approvals)
	if err != nil {
		t.Fatal(err)
	}
	attemptDigest, err := attemptSetDigest(result.Attempts)
	if err != nil {
		t.Fatal(err)
	}
	result.ApprovalSetDigest = approvalDigest
	result.AttemptSetDigest = attemptDigest
	return result, records
}

func generatedLifecycleRecord(
	t *testing.T,
	highWater v0candidate.UInt53,
	bindings lifecyclestate.ImmutableBindings,
	state lifecyclestate.LifecycleState,
) lifecyclestate.Record {
	t.Helper()
	view := lifecyclestate.RecordView{
		FormatVersion: lifecyclestate.LifecycleRecordFormatVersion, RecordVersion: 1,
		SnapshotGeneration: 1, Bindings: bindings, ImmutableBindingDigest: bindings.Digest(),
		State: state, CleanupRequired: true, LastConfirmedCheckpoint: lifecyclestate.CheckpointNone,
		FirstFailure: lifecyclestate.FailureNone, FailureOperation: lifecyclestate.OperationNone,
		LastReconciliation: lifecyclestate.ReconciliationNone,
		RecoveryFence:      lifecyclestate.RecoveryFenceNone,
		OpenedAt:           bindings.View().AttemptCreatedAt,
		LastTransitionAt:   bindings.View().AttemptCreatedAt,
	}
	setEffect := func(operation lifecyclestate.Operation, status lifecyclestate.EffectStatus, checkpoint lifecyclestate.Checkpoint) {
		view.Operation = operation
		view.EffectStatus = status
		view.LastConfirmedCheckpoint = checkpoint
		attemptID := bindings.View().AttemptID
		identifier, err := lifecyclestate.NewDomainIdentifier(lifecyclestate.DomainEffectID, attemptID[:])
		if err != nil {
			t.Fatal(err)
		}
		view.EffectID, err = lifecyclestate.NewEffectID(identifier)
		if err != nil {
			t.Fatal(err)
		}
		view.OperationSequence = 1
		view.Instance, err = lifecyclestate.NewBackendInstanceIdentity(
			lifecyclestate.BackendInstanceFake, attemptID[:],
		)
		if err != nil {
			t.Fatal(err)
		}
	}
	switch state {
	case lifecyclestate.StatePreparePending:
		view.Operation = lifecyclestate.OperationNone
		view.EffectStatus = lifecyclestate.EffectNone
	case lifecyclestate.StateObserved:
		setEffect(lifecyclestate.OperationObserve, lifecyclestate.EffectConfirmed, lifecyclestate.CheckpointObserve)
	case lifecyclestate.StateStopped:
		setEffect(lifecyclestate.OperationStop, lifecyclestate.EffectConfirmed, lifecyclestate.CheckpointStop)
	case lifecyclestate.StateDestroyConfirmed:
		setEffect(lifecyclestate.OperationDestroy, lifecyclestate.EffectConfirmed, lifecyclestate.CheckpointDestroy)
	case lifecyclestate.StateUnresolved:
		setEffect(lifecyclestate.OperationObserve, lifecyclestate.EffectIndeterminate, lifecyclestate.CheckpointStart)
		view.FirstFailure = lifecyclestate.FailureCleanupUnresolved
		view.FailureOperation = lifecyclestate.OperationObserve
		view.LastReconciliation = lifecyclestate.ReconciliationUnknown
		view.AutomaticRecoveryCount = 1
		view.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{Present: true, Value: highWater + 1}
		view.RecoveryFence = lifecyclestate.RecoveryFenceReconcileUnknown
	case lifecyclestate.StateQuarantined:
		setEffect(lifecyclestate.OperationStart, lifecyclestate.EffectIndeterminate, lifecyclestate.CheckpointCreate)
		view.FirstFailure = lifecyclestate.FailureBinding
		view.FailureOperation = lifecyclestate.OperationStart
		view.LastReconciliation = lifecyclestate.ReconciliationIdentityMismatch
		view.AutomaticRecoveryCount = 1
		view.RecoveryFence = lifecyclestate.RecoveryFenceIdentityMismatch
	case lifecyclestate.StateDestroyed:
		view.Operation = lifecyclestate.OperationNone
		view.EffectStatus = lifecyclestate.EffectNone
		view.CleanupRequired = false
		view.LastReconciliation = lifecyclestate.ReconciliationAuthoritativelyAbsent
		view.AutomaticRecoveryCount = 1
		view.TerminalAt = lifecyclestate.OptionalUnixSeconds{Present: true, Value: view.LastTransitionAt}
	default:
		t.Fatalf("unsupported generated lifecycle state %q", state)
	}
	record, err := lifecyclestate.NewRecord(view, highWater)
	if err != nil {
		t.Fatalf("generated lifecycle record %s: %v", state, err)
	}
	return record
}

func stateAndLifecycleRecord(t *testing.T) (installationState, lifecyclestate.Record) {
	t.Helper()
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
	state, _ := harness.store.snapshot(context.Background())
	attempt := state.Attempts[0]
	approval := state.Approvals[0]
	registration := state.Registrations[0]
	implementation := lifecyclestate.BackendImplementationDigest(sha256.Sum256([]byte("fake-no-guest-e2")))
	backend, err := lifecyclestate.NewBackendBinding(lifecyclestate.BackendBindingView{
		Kind:                          lifecyclestate.BackendFakeNoGuest,
		ProtocolVersion:               lifecyclestate.FakeBackendProtocolVersion,
		ImplementationIdentityDigest:  implementation,
		BackendConfigurationDigest:    registration.PlanBindings.BackendConfigurationDigest,
		BackendValidationRecordDigest: registration.PlanBindings.BackendValidationRecordDigest,
		CreatesGuest:                  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	bindings, err := lifecyclestate.NewImmutableBindings(lifecyclestate.ImmutableBindingsView{
		AttemptID:                       attempt.AttemptID,
		ApprovalID:                      attempt.ApprovalID,
		AttemptNonce:                    attempt.AttemptNonce,
		RegistrationID:                  attempt.RegistrationID,
		RegistrationSequence:            attempt.RegistrationSequence,
		PlanDigest:                      attempt.PlanDigest,
		InstallationID:                  attempt.InstallationID,
		EpochSequence:                   attempt.EpochSequence,
		EpochDigest:                     attempt.EpochDigest,
		SupervisorID:                    attempt.SupervisorID,
		ApprovalPurpose:                 attempt.Purpose,
		ApprovalAudience:                attempt.Audience,
		ApprovalPayloadDigest:           attempt.ApprovalPayloadDigest,
		AuthorizationIdentity:           attempt.AuthorizationIdentity,
		AttemptCreatedAt:                attempt.CreatedAt,
		AttemptStorageVersion:           attempt.StorageFormatVersion,
		ApprovalStorageVersion:          approval.StorageFormatVersion,
		RegistrationStorageVersion:      registration.Record.StorageFormatVersion,
		SourceManifestDigest:            registration.PlanBindings.SourceManifestDigest,
		InlineInputDigest:               registration.PlanBindings.InlineInputDigest,
		RuntimeBundleManifestDigest:     registration.PlanBindings.RuntimeBundleManifestDigest,
		ProfileReviewAttestationDigests: registration.PlanBindings.ProfileReviewAttestationDigests,
		ProfileRegistryEntryDigest:      registration.PlanBindings.ProfileRegistryEntryDigest,
		BackendValidationRecordDigest:   registration.PlanBindings.BackendValidationRecordDigest,
		BackendConfigurationDigest:      registration.PlanBindings.BackendConfigurationDigest,
		TrustSnapshotDigest:             registration.PlanBindings.TrustSnapshotDigest,
		PolicyDecisionDigest:            registration.PlanBindings.PolicyDecisionDigest,
		Backend:                         backend,
	})
	if err != nil {
		t.Fatal(err)
	}
	record, err := lifecyclestate.NewRecord(lifecyclestate.RecordView{
		FormatVersion:           lifecyclestate.LifecycleRecordFormatVersion,
		RecordVersion:           1,
		SnapshotGeneration:      1,
		Bindings:                bindings,
		ImmutableBindingDigest:  bindings.Digest(),
		Operation:               lifecyclestate.OperationNone,
		EffectStatus:            lifecyclestate.EffectNone,
		State:                   lifecyclestate.StatePreparePending,
		CleanupRequired:         true,
		LastConfirmedCheckpoint: lifecyclestate.CheckpointNone,
		FirstFailure:            lifecyclestate.FailureNone,
		FailureOperation:        lifecyclestate.OperationNone,
		LastReconciliation:      lifecyclestate.ReconciliationNone,
		RecoveryFence:           lifecyclestate.RecoveryFenceNone,
		OpenedAt:                attempt.CreatedAt,
		LastTransitionAt:        attempt.CreatedAt,
	}, state.TimeHighWaterUnixSeconds)
	if err != nil {
		t.Fatal(err)
	}
	return state, record
}

func lifecycleRecordWithMutatedBindings(
	t *testing.T,
	state installationState,
	record lifecyclestate.Record,
	mutate func(*lifecyclestate.ImmutableBindingsView),
) lifecyclestate.Record {
	t.Helper()
	view := record.View()
	bindingsView := view.Bindings.View()
	mutate(&bindingsView)
	bindings, err := lifecyclestate.NewImmutableBindings(bindingsView)
	if err != nil {
		t.Fatalf("mutated bindings should remain structurally valid: %v", err)
	}
	view.Bindings = bindings
	view.ImmutableBindingDigest = bindings.Digest()
	result, err := lifecyclestate.NewRecord(view, state.TimeHighWaterUnixSeconds)
	if err != nil {
		t.Fatalf("mutated record should remain structurally valid: %v", err)
	}
	return result
}

func ordinaryInitialStateForV1() installationState {
	state := installationState{
		InitialState:          ordinaryInitialState(),
		RegistrationSetDigest: emptyRegistrationSetDigest(),
		Registrations:         make([]registrationEntry, 0),
		ApprovalSetDigest:     emptyApprovalSetDigest(),
		Approvals:             make([]approvalattempt.ApprovalRecord, 0),
		AttemptSetDigest:      emptyAttemptSetDigest(),
		Attempts:              make([]approvalattempt.ExecutionAttempt, 0),
	}
	return state
}

func newV0StorePath(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "state-v0.json")
	if _, err := NewFixedFileStore(path, ordinaryInitialState()); err != nil {
		t.Fatalf("new v0 store: %v", err)
	}
	return path
}

func writeV1Envelope(t *testing.T, path string, envelope diskEnvelopeV1) {
	t.Helper()
	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func cloneDiskEnvelopeV1(t *testing.T, envelope diskEnvelopeV1) diskEnvelopeV1 {
	t.Helper()
	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	var clone diskEnvelopeV1
	if err := json.Unmarshal(data, &clone); err != nil {
		t.Fatal(err)
	}
	return clone
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func assertV1OpenDoesNotRewrite(t *testing.T, path string) {
	t.Helper()
	before := mustReadFile(t, path)
	if _, err := OpenFixedFileStoreV1(path); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
		t.Fatalf("corrupt v1 open error = %v", err)
	}
	if after := mustReadFile(t, path); !bytes.Equal(after, before) {
		t.Fatal("corrupt v1 open rewrote evidence")
	}
}

func authorityStatesEqual(left, right installationState) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	return leftErr == nil && rightErr == nil && bytes.Equal(leftJSON, rightJSON)
}
