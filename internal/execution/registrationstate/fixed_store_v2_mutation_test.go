package registrationstate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"path/filepath"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
)

func TestFixedStoreV2AuthorityAndLifecycleMutationsReopenDeterministically(t *testing.T) {
	harness := newApprovalHarness(t)
	vector := fixtureVector(t, harness, 0x72, 1_785_456_000, 1_785_456_300, 0)
	if _, err := MigrateFixedFileStoreV0ToV1(context.Background(), harness.path, V0ToV1MigrationOptions{Lock: &migrationLockStub{held: true}}); err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateFixedFileStoreV1ToV2(context.Background(), harness.path, V1ToV2MigrationOptions{Lock: &migrationLockStub{held: true}}); err != nil {
		t.Fatal(err)
	}
	store, err := OpenFixedFileStoreV2(harness.path)
	if err != nil {
		t.Fatal(err)
	}
	component := mustApprovalAttemptComponent(t, approvalHarness{
		path: harness.path, store: harness.store, clock: harness.clock,
		registrationID: harness.registrationID, registration: harness.registration, planDigest: harness.planDigest,
	}, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 21}, &attemptIDSequence{next: 21}, nil)
	component.store = store
	submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
	if err != nil {
		t.Fatalf("v2 approval mutation: %v", err)
	}
	created, err := component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
	if err != nil {
		t.Fatalf("v2 attempt mutation: %v", err)
	}
	registration := store.state.Registrations[0]
	backend, err := lifecyclestate.NewBackendBinding(lifecyclestate.BackendBindingView{
		Kind: lifecyclestate.BackendFakeNoGuest, ProtocolVersion: lifecyclestate.FakeBackendProtocolVersion,
		ImplementationIdentityDigest:  lifecyclestate.BackendImplementationDigest{1},
		BackendConfigurationDigest:    registration.PlanBindings.BackendConfigurationDigest,
		BackendValidationRecordDigest: registration.PlanBindings.BackendValidationRecordDigest,
		CreatesGuest:                  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, createdLifecycle, err := store.EnsureLifecycle(context.Background(), created.Reference.AttemptID(), backend); err != nil || !createdLifecycle {
		t.Fatalf("v2 lifecycle mutation = %v, %v", createdLifecycle, err)
	}
	firstBytes := mustReadFile(t, harness.path)
	reopened, err := OpenFixedFileStoreV2(harness.path)
	if err != nil {
		t.Fatalf("reopen v2 authority/lifecycle world: %v", err)
	}
	if got, err := reopened.ResolveAttempt(context.Background(), created.Reference.AttemptID()); err != nil || !got.LifecyclePresent {
		t.Fatalf("reopened v2 attempt/lifecycle = %#v, %v", got, err)
	}
	if secondBytes := mustReadFile(t, harness.path); !bytes.Equal(firstBytes, secondBytes) {
		t.Fatal("read-only reopen changed deterministic v2 bytes")
	}
}

func newV2LifecycleMutationHarness(t *testing.T) (lifecycleTransactionHarness, *FixedFileStoreV2) {
	t.Helper()
	harness := newLifecycleTransactionHarness(t, 1)
	record, _, err := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
	if err != nil || record.View().State != lifecyclestate.StatePreparePending {
		t.Fatalf("establish v1 lifecycle: %#v, %v", record.View(), err)
	}
	if _, err := MigrateFixedFileStoreV1ToV2(context.Background(), harness.path, V1ToV2MigrationOptions{
		Lock: &migrationLockStub{held: true},
	}); err != nil {
		t.Fatalf("migrate mutation harness to v2: %v", err)
	}
	effectIDs := &lifecycleEffectIDSequence{next: 101}
	store, err := OpenFixedFileStoreV2WithOptions(harness.path, FixedFileStoreV1Options{
		EffectIDs: effectIDs, OwnerSessionID: harness.owner,
	})
	if err != nil {
		t.Fatalf("open v2 mutation harness: %v", err)
	}
	harness.effectIDs = effectIDs
	return harness, store
}

func TestFixedStoreV2BeginEffectAppendsExactlyOneAtomicHistoricalTombstone(t *testing.T) {
	harness, store := newV2LifecycleMutationHarness(t)
	attemptID := harness.attemptIDs[0]
	record, err := store.ReadLifecycle(context.Background(), attemptID)
	if err != nil {
		t.Fatal(err)
	}
	var first lifecyclestate.EffectID
	for _, operation := range []lifecyclestate.Operation{lifecyclestate.OperationPrepare, lifecyclestate.OperationCreate} {
		permit, beginErr := store.BeginEffect(context.Background(), attemptID, record.View().RecordVersion, operation)
		if beginErr != nil {
			t.Fatalf("begin %s: %v", operation, beginErr)
		}
		replay, replayErr := store.BeginEffect(context.Background(), attemptID, record.View().RecordVersion, operation)
		if replayErr != nil || replay.View().EffectID != permit.View().EffectID {
			t.Fatalf("exact intent replay %s changed effect: %#v, %v", operation, replay.View(), replayErr)
		}
		if operation == lifecyclestate.OperationPrepare {
			first = permit.View().EffectID
		}
		resultInstance := lifecyclestate.BackendInstanceIdentity{}
		if operation == lifecyclestate.OperationCreate {
			instance, instanceErr := lifecyclestate.NewBackendInstanceIdentity(lifecyclestate.BackendInstanceFake, []byte("f4b-instance"))
			if instanceErr != nil {
				t.Fatal(instanceErr)
			}
			resultInstance = instance
		}
		result, resultErr := lifecyclestate.NewEffectResult(operation, lifecyclestate.EffectResultApplied, resultInstance)
		if resultErr != nil {
			t.Fatal(resultErr)
		}
		record, err = store.ConfirmEffect(context.Background(), permit, result)
		if err != nil {
			t.Fatalf("confirm %s: %v", operation, err)
		}
	}
	reopened, err := OpenFixedFileStoreV2(harness.path)
	if err != nil {
		t.Fatalf("reopen two-effect v2 world: %v", err)
	}
	historical, err := reopened.ResolveEffect(context.Background(), first)
	if err != nil {
		t.Fatal(err)
	}
	if historical.Classification != RetainedEffectSupersededByCurrent || historical.LifecyclePresent ||
		historical.Operation != lifecyclestate.OperationPrepare || historical.VisibleV1Seed {
		t.Fatalf("historical resolution invented issuing lifecycle: %#v", historical)
	}
	current, err := reopened.ResolveEffect(context.Background(), record.View().EffectID)
	if err != nil || current.Classification != RetainedEffectCurrent || !current.LifecyclePresent {
		t.Fatalf("current effect resolution = %#v, %v", current, err)
	}
	if got := len(reopened.active.View().Indexes.View().Effects); got != 2 || harness.effectIDs.calls != 2 {
		t.Fatalf("effect ledger/calls = %d/%d, want 2/2", got, harness.effectIDs.calls)
	}
	if reopened.active.View().HotSetDigests.EffectTombstones == (archivestate.EffectTombstoneSetDigest{}) {
		t.Fatal("materialized F4B world lacks effect-tombstone hot-set digest")
	}
}

func TestFixedStoreV2BeginEffectFaultsAreNeitherOrBothAndNoRewriteOnCorruption(t *testing.T) {
	cases := []struct {
		name      string
		fault     LifecycleStoreFault
		committed bool
		wantIndet bool
	}{
		{name: "confirmed-abort", fault: FaultLifecycleIntentAbort},
		{name: "pre-state-indeterminate", fault: FaultLifecycleIntentIndeterminatePreState, wantIndet: true},
		{name: "post-rename-indeterminate-process-death", fault: FaultLifecycleIntentIndeterminate, committed: true, wantIndet: true},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			harness, store := newV2LifecycleMutationHarness(t)
			attemptID := harness.attemptIDs[0]
			record, _ := store.ReadLifecycle(context.Background(), attemptID)
			before := mustReadFile(t, harness.path)
			store.InjectLifecycleFailure(test.fault, errors.New("injected"))
			permit, beginErr := store.BeginEffect(context.Background(), attemptID, record.View().RecordVersion, lifecyclestate.OperationPrepare)
			if beginErr == nil || (test.wantIndet && !errors.Is(beginErr, ErrCommitOutcomeIndeterminate)) {
				t.Fatalf("fault result = %#v, %v", permit.View(), beginErr)
			}
			after := mustReadFile(t, harness.path)
			if bytes.Equal(before, after) == test.committed {
				t.Fatalf("committed-byte classification = %v, want %v", !bytes.Equal(before, after), test.committed)
			}
			reopened, err := OpenFixedFileStoreV2WithOptions(harness.path, FixedFileStoreV1Options{OwnerSessionID: harness.owner})
			if err != nil {
				t.Fatal(err)
			}
			if got := len(reopened.active.View().Indexes.View().Effects); got != boolCount(test.committed) {
				t.Fatalf("reopened effect count = %d, want %d", got, boolCount(test.committed))
			}
			if test.committed {
				// The caller lost the response and the original process is treated as
				// dead. Reopen plus exact retry must recover the one durable issuance.
				replay, replayErr := reopened.BeginEffect(context.Background(), attemptID, record.View().RecordVersion, lifecyclestate.OperationPrepare)
				if replayErr != nil || replay.View().EffectID.IsZero() || len(reopened.effectTombstones) != 1 {
					t.Fatalf("response-loss retry = %#v, %v", replay.View(), replayErr)
				}
			}
		})
	}
}

func boolCount(value bool) int {
	if value {
		return 1
	}
	return 0
}

func TestFixedStoreV2BeginEffectConcurrentReplayOneIssuance(t *testing.T) {
	harness, store := newV2LifecycleMutationHarness(t)
	record, err := store.ReadLifecycle(context.Background(), harness.attemptIDs[0])
	if err != nil {
		t.Fatal(err)
	}
	const callers = 24
	permits := make(chan lifecyclestate.EffectPermit, callers)
	errorsSeen := make(chan error, callers)
	var group sync.WaitGroup
	for range callers {
		group.Add(1)
		go func() {
			defer group.Done()
			permit, beginErr := store.BeginEffect(context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare)
			permits <- permit
			errorsSeen <- beginErr
		}()
	}
	group.Wait()
	close(permits)
	close(errorsSeen)
	var effectID lifecyclestate.EffectID
	for beginErr := range errorsSeen {
		if beginErr != nil {
			t.Fatalf("concurrent exact retry: %v", beginErr)
		}
	}
	for permit := range permits {
		if effectID.IsZero() {
			effectID = permit.View().EffectID
		}
		if permit.View().EffectID != effectID {
			t.Fatal("concurrent exact retry returned a second effect ID")
		}
	}
	reopened, err := OpenFixedFileStoreV2(harness.path)
	if err != nil {
		t.Fatal(err)
	}
	if len(reopened.effectTombstones) != 1 || harness.effectIDs.calls != 1 {
		t.Fatalf("concurrent tombstones/source calls = %d/%d, want 1/1", len(reopened.effectTombstones), harness.effectIDs.calls)
	}
}

func TestFixedStoreV2MutationOwnerFenceAndEffectCollectionCorruptionNoRewrite(t *testing.T) {
	t.Run("owner-session-loss-before-publish", func(t *testing.T) {
		harness, _ := newV2LifecycleMutationHarness(t)
		owner := &archiveOwnerStub{held: true, session: harness.owner}
		store, err := OpenFixedFileStoreV2WithOwner(context.Background(), harness.path, owner, FixedFileStoreV1Options{
			EffectIDs: &lifecycleEffectIDSequence{next: 301}, OwnerSessionID: harness.owner,
		})
		if err != nil {
			t.Fatal(err)
		}
		record, _ := store.ReadLifecycle(context.Background(), harness.attemptIDs[0])
		before := mustReadFile(t, harness.path)
		owner.held = false
		if _, err := store.BeginEffect(context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare); !errors.Is(err, ErrStoreOwnerFenced) {
			t.Fatalf("owner loss result = %v", err)
		}
		if !bytes.Equal(before, mustReadFile(t, harness.path)) {
			t.Fatal("owner loss before publication changed bytes")
		}
	})

	harness, store := newV2LifecycleMutationHarness(t)
	record, _ := store.ReadLifecycle(context.Background(), harness.attemptIDs[0])
	if _, err := store.BeginEffect(context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare); err != nil {
		t.Fatal(err)
	}
	valid := readV2Envelope(t, harness.path)
	mutations := []struct {
		name   string
		mutate func(*diskEnvelopeV2)
	}{
		{name: "digest", mutate: func(envelope *diskEnvelopeV2) { (*envelope.EffectTombstoneSetDigest)[0] ^= 0xff }},
		{name: "cross-link", mutate: func(envelope *diskEnvelopeV2) { (*envelope.EffectTombstones)[0].OperationSequence++ }},
		{name: "missing-after-F4B", mutate: func(envelope *diskEnvelopeV2) {
			envelope.EffectTombstones = nil
			envelope.EffectTombstoneSetDigest = nil
			envelope.MigrationGenesisIndexes = nil
		}},
	}
	for _, test := range mutations {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "state-v2.json")
			writeV2Envelope(t, path, cloneV2Envelope(t, valid))
			envelope := readV2Envelope(t, path)
			test.mutate(&envelope)
			writeV2Envelope(t, path, envelope)
			assertV2OpenDoesNotRewrite(t, path)
		})
	}
}

func TestFixedStoreV2MutationAfterArchivePreservesHistoricalSegmentAndLookup(t *testing.T) {
	path, store, owner := newEligibleFixedStoreV2(t)
	keys := lookupKeysForStore(t, store)
	store = activateLookupStore(t, store, owner)
	segmentPath := archiveSegmentPath(path, store.active.View().Descriptors[0].View().SegmentDigest)
	segmentBefore := mustReadFile(t, segmentPath)
	if err := store.persistTimeHighWater(context.Background(), store.state.TimeHighWaterUnixSeconds+1); err != nil {
		t.Fatal(err)
	}
	reopened, err := OpenFixedFileStoreV2(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(segmentBefore, mustReadFile(t, segmentPath)) || reopened.active.View().ArchivedCounts.Effects != 1 ||
		reopened.active.View().HotCounts.Effects != 0 || reopened.active.View().SnapshotGeneration != 3 {
		t.Fatalf("post-archive mutation changed retained shape: %#v", reopened.active.View())
	}
	resolved, err := reopened.ResolveEffect(context.Background(), keys.effectID)
	if err != nil || resolved.Classification != RetainedEffectCurrent || !resolved.LifecyclePresent {
		t.Fatalf("archived effect after hot mutation = %#v, %v", resolved, err)
	}
}

func TestFixedStoreV2F4BKnownAnswers(t *testing.T) {
	harness, store := newV2LifecycleMutationHarness(t)
	record, _ := store.ReadLifecycle(context.Background(), harness.attemptIDs[0])
	if _, err := store.BeginEffect(context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare); err != nil {
		t.Fatal(err)
	}
	reopened, err := OpenFixedFileStoreV2(harness.path)
	if err != nil {
		t.Fatal(err)
	}
	fileDigest := sha256.Sum256(mustReadFile(t, harness.path))
	tombstoneDigest := reopened.active.View().HotSetDigests.EffectTombstones
	combined := reopened.active.View().CombinedIndexDigest
	answers := map[string]struct {
		got  []byte
		want string
	}{
		"active-file":       {got: fileDigest[:], want: "fe7fc82e0b38a1c1f7a1f99e7898090323d7b1c84565e30557f44ab2e07a070c"},
		"effect-tombstones": {got: tombstoneDigest[:], want: "690afbf0d8fff394110858f0fd20f8a493854e91b2bdaf7160320ed79898a1e1"},
		"combined-index":    {got: combined[:], want: "eba9f827a4cc402aa8a1a209e441f22838f8587e5487742979b631e107b1b9f8"},
	}
	for name, answer := range answers {
		if got := hex.EncodeToString(answer.got); got != answer.want {
			t.Errorf("%s = %s, want %s", name, got, answer.want)
		}
	}
	if reopened.active.View().SnapshotGeneration != 2 || len(reopened.effectTombstones) != 1 {
		t.Fatalf("known-answer shape = generation %d tombstones %d", reopened.active.View().SnapshotGeneration, len(reopened.effectTombstones))
	}
}
