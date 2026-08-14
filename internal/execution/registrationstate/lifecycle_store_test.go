package registrationstate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
)

type lifecycleEffectIDSequence struct {
	mu        sync.Mutex
	next      uint64
	calls     int
	collision lifecyclestate.EffectID
	err       error
}

func (source *lifecycleEffectIDSequence) NewEffectID(context.Context) (lifecyclestate.EffectID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	source.calls++
	if source.err != nil {
		return lifecyclestate.EffectID{}, source.err
	}
	if !source.collision.IsZero() {
		return source.collision, nil
	}
	value := make([]byte, 16)
	value[0] = 0xe3
	binary.BigEndian.PutUint64(value[8:], source.next)
	source.next++
	domain, err := lifecyclestate.NewDomainIdentifier(lifecyclestate.DomainEffectID, value)
	if err != nil {
		return lifecyclestate.EffectID{}, err
	}
	return lifecyclestate.NewEffectID(domain)
}

type lifecycleTransactionHarness struct {
	path       string
	store      *FixedFileStoreV1
	attemptIDs []approvalattempt.AttemptID
	backend    lifecyclestate.BackendBinding
	effectIDs  *lifecycleEffectIDSequence
	owner      lifecyclestate.OwnerSessionID
}

func newLifecycleTransactionHarness(t *testing.T, attemptCount int) lifecycleTransactionHarness {
	t.Helper()
	harness := newApprovalHarness(t)
	vectors := make([]approvalattempt.FixtureVector, attemptCount)
	for index := range vectors {
		vectors[index] = fixtureVector(
			t,
			harness,
			byte(0x60+index),
			1_785_456_000,
			1_785_456_300,
			byte(index),
		)
	}
	component := mustApprovalAttemptComponent(
		t,
		harness,
		vectors,
		&approvalIDSequence{next: 1},
		&attemptIDSequence{next: 1},
		nil,
	)
	attemptIDs := make([]approvalattempt.AttemptID, 0, attemptCount)
	for _, vector := range vectors {
		submission, err := component.SubmitApproval(
			context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes,
		)
		if err != nil {
			t.Fatalf("submit approval: %v", err)
		}
		created, err := component.RequestAttempt(
			context.Background(), attemptCall(), harness.registrationID, submission.Reference,
		)
		if err != nil {
			t.Fatalf("request attempt: %v", err)
		}
		attemptIDs = append(attemptIDs, created.Reference.AttemptID())
	}
	state, err := harness.store.snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	registration := state.Registrations[0]
	backend, err := lifecyclestate.NewBackendBinding(lifecyclestate.BackendBindingView{
		Kind:                          lifecyclestate.BackendFakeNoGuest,
		ProtocolVersion:               lifecyclestate.FakeBackendProtocolVersion,
		ImplementationIdentityDigest:  lifecyclestate.BackendImplementationDigest(sha256.Sum256([]byte("fake-no-guest-e3"))),
		BackendConfigurationDigest:    registration.PlanBindings.BackendConfigurationDigest,
		BackendValidationRecordDigest: registration.PlanBindings.BackendValidationRecordDigest,
		CreatesGuest:                  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateFixedFileStoreV0ToV1(
		context.Background(),
		harness.path,
		V0ToV1MigrationOptions{Lock: &migrationLockStub{held: true}},
	); err != nil {
		t.Fatalf("migrate transaction harness: %v", err)
	}
	effectIDs := &lifecycleEffectIDSequence{next: 1}
	owner := lifecycleOwnerID(t, 0x41)
	store, err := OpenFixedFileStoreV1WithOptions(harness.path, FixedFileStoreV1Options{
		EffectIDs: effectIDs, OwnerSessionID: owner,
	})
	if err != nil {
		t.Fatalf("open transaction store: %v", err)
	}
	return lifecycleTransactionHarness{
		path: harness.path, store: store, attemptIDs: attemptIDs, backend: backend,
		effectIDs: effectIDs, owner: owner,
	}
}

func lifecycleOwnerID(t *testing.T, value byte) lifecyclestate.OwnerSessionID {
	t.Helper()
	domain, err := lifecyclestate.NewDomainIdentifier(
		lifecyclestate.DomainOwnerSessionID, bytes.Repeat([]byte{value}, 16),
	)
	if err != nil {
		t.Fatal(err)
	}
	owner, err := lifecyclestate.NewOwnerSessionID(domain)
	if err != nil {
		t.Fatal(err)
	}
	return owner
}

func TestDurableLifecycleEnsureReadReplayConcurrencyAndReopen(t *testing.T) {
	harness := newLifecycleTransactionHarness(t, 2)
	firstID := harness.attemptIDs[0]

	const goroutines = 24
	var wait sync.WaitGroup
	created := make(chan bool, goroutines)
	errorsSeen := make(chan error, goroutines)
	for range goroutines {
		wait.Add(1)
		go func() {
			defer wait.Done()
			record, wasCreated, err := harness.store.EnsureLifecycle(context.Background(), firstID, harness.backend)
			if err == nil && record.Bindings().View().AttemptID != firstID {
				err = errors.New("wrong attempt binding")
			}
			created <- wasCreated
			errorsSeen <- err
		}()
	}
	wait.Wait()
	close(created)
	close(errorsSeen)
	createdCount := 0
	for wasCreated := range created {
		if wasCreated {
			createdCount++
		}
	}
	for err := range errorsSeen {
		if err != nil {
			t.Fatalf("concurrent ensure: %v", err)
		}
	}
	if createdCount != 1 {
		t.Fatalf("created count = %d, want 1", createdCount)
	}
	read, err := harness.store.ReadLifecycle(context.Background(), firstID)
	if err != nil {
		t.Fatal(err)
	}
	if read.View().State != lifecyclestate.StatePreparePending || !read.View().CleanupRequired || read.View().RecordVersion != 1 {
		t.Fatalf("initial lifecycle = %#v", read.View())
	}

	second, made, err := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[1], harness.backend)
	if err != nil || !made || second.Bindings().View().RegistrationID != read.Bindings().View().RegistrationID {
		t.Fatalf("second attempt ensure = made %v, err %v", made, err)
	}
	if second.Bindings().View().AttemptID == read.Bindings().View().AttemptID ||
		second.Bindings().View().ApprovalID == read.Bindings().View().ApprovalID {
		t.Fatal("two approvals for one registration shared lifecycle authority")
	}

	reopened, err := OpenFixedFileStoreV1WithOptions(harness.path, FixedFileStoreV1Options{
		EffectIDs: &lifecycleEffectIDSequence{next: 10}, OwnerSessionID: lifecycleOwnerID(t, 0x42),
	})
	if err != nil {
		t.Fatal(err)
	}
	readAgain, err := reopened.ReadLifecycle(context.Background(), firstID)
	if err != nil || readAgain.View().ImmutableBindingDigest != read.View().ImmutableBindingDigest {
		t.Fatalf("reopen lifecycle = %#v, %v", readAgain.View(), err)
	}
	returned := readAgain.Bindings().View()
	returned.ProfileReviewAttestationDigests[0][0] ^= 0xff
	unchanged, _ := reopened.ReadLifecycle(context.Background(), firstID)
	if unchanged.Bindings().View().ProfileReviewAttestationDigests[0] == returned.ProfileReviewAttestationDigests[0] {
		t.Fatal("read lifecycle returned aliased plan bindings")
	}
}

func TestDurableLifecycleBeginReplayAndIntentFaultOracles(t *testing.T) {
	for _, test := range []struct {
		name       string
		fault      LifecycleStoreFault
		wantFence  bool
		wantIntent bool
	}{
		{name: "confirmed-abort", fault: FaultLifecycleIntentAbort},
		{name: "indeterminate-pre-state", fault: FaultLifecycleIntentIndeterminatePreState, wantFence: true},
		{name: "indeterminate-post-rename", fault: FaultLifecycleIntentIndeterminate, wantFence: true, wantIntent: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			harness := newLifecycleTransactionHarness(t, 1)
			record, _, err := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
			if err != nil {
				t.Fatal(err)
			}
			before := mustReadFile(t, harness.path)
			harness.store.InjectLifecycleFailure(test.fault, errors.New("fault"))
			if _, err := harness.store.BeginEffect(
				context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare,
			); err == nil {
				t.Fatal("faulted begin succeeded")
			}
			if harness.store.LifecycleRecoveryFenced() != test.wantFence {
				t.Fatalf("fence = %v, want %v", harness.store.LifecycleRecoveryFenced(), test.wantFence)
			}
			if !test.wantIntent && !bytes.Equal(before, mustReadFile(t, harness.path)) {
				t.Fatal("pre-rename intent fault changed durable state")
			}
			reopened, err := OpenFixedFileStoreV1WithOptions(harness.path, FixedFileStoreV1Options{
				EffectIDs: &lifecycleEffectIDSequence{next: 100}, OwnerSessionID: lifecycleOwnerID(t, 0x52),
			})
			if err != nil {
				t.Fatal(err)
			}
			after, err := reopened.ReadLifecycle(context.Background(), harness.attemptIDs[0])
			if err != nil {
				t.Fatal(err)
			}
			if gotIntent := after.View().State == lifecyclestate.StatePrepareIntent; gotIntent != test.wantIntent {
				t.Fatalf("reopen intent = %v, want %v", gotIntent, test.wantIntent)
			}
		})
	}

	harness := newLifecycleTransactionHarness(t, 1)
	record, _, _ := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
	permit, err := harness.store.BeginEffect(
		context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare,
	)
	if err != nil {
		t.Fatal(err)
	}
	replay, err := harness.store.BeginEffect(
		context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare,
	)
	if err != nil || replay.View().EffectID != permit.View().EffectID || harness.effectIDs.calls != 1 {
		t.Fatalf("intent replay = %#v, calls %d, err %v", replay.View(), harness.effectIDs.calls, err)
	}
}

func TestDurableLifecycleEffectCheckpointsAndAuthoritativeDestroy(t *testing.T) {
	harness := newLifecycleTransactionHarness(t, 1)
	attemptID := harness.attemptIDs[0]
	record, _, err := harness.store.EnsureLifecycle(context.Background(), attemptID, harness.backend)
	if err != nil {
		t.Fatal(err)
	}
	instance, err := lifecyclestate.NewBackendInstanceIdentity(lifecyclestate.BackendInstanceFake, []byte("instance-one"))
	if err != nil {
		t.Fatal(err)
	}
	operations := []lifecyclestate.Operation{
		lifecyclestate.OperationPrepare,
		lifecyclestate.OperationCreate,
		lifecyclestate.OperationStart,
		lifecyclestate.OperationObserve,
		lifecyclestate.OperationStop,
		lifecyclestate.OperationDestroy,
	}
	for _, operation := range operations {
		permit, err := harness.store.BeginEffect(context.Background(), attemptID, record.View().RecordVersion, operation)
		if err != nil {
			t.Fatalf("begin %s: %v", operation, err)
		}
		resultInstance := lifecyclestate.BackendInstanceIdentity{}
		if operation == lifecyclestate.OperationCreate {
			resultInstance = instance
		}
		result, err := lifecyclestate.NewEffectResult(operation, lifecyclestate.EffectResultApplied, resultInstance)
		if err != nil {
			t.Fatal(err)
		}
		record, err = harness.store.ConfirmEffect(context.Background(), permit, result)
		if err != nil {
			t.Fatalf("confirm %s: %v", operation, err)
		}
		replayed, err := harness.store.ConfirmEffect(context.Background(), permit, result)
		if err != nil || replayed.View().RecordVersion != record.View().RecordVersion {
			t.Fatalf("confirm replay %s: %#v, %v", operation, replayed.View(), err)
		}
	}
	if record.View().State != lifecyclestate.StateDestroyConfirmed || !record.View().CleanupRequired {
		t.Fatalf("destroy confirmation cleared cleanup: %#v", record.View())
	}
	absent, err := lifecyclestate.NewReconcileResult(
		lifecyclestate.OperationDestroy, lifecyclestate.ReconciliationAuthoritativelyAbsent, lifecyclestate.BackendInstanceIdentity{},
	)
	if err != nil {
		t.Fatal(err)
	}
	record, err = harness.store.RecordReconciliation(context.Background(), attemptID, record.View().RecordVersion, absent)
	if err != nil {
		t.Fatal(err)
	}
	if record.View().State != lifecyclestate.StateDestroyed || record.View().CleanupRequired || !record.View().TerminalAt.Present {
		t.Fatalf("authoritative absence = %#v", record.View())
	}
	if ids, err := harness.store.RecoveryAttemptIDs(context.Background()); err != nil || len(ids) != 0 {
		t.Fatalf("terminal recovery set = %x, %v", ids, err)
	}
}

func TestDurableLifecycleAfterEffectAndIndeterminateReopenOracles(t *testing.T) {
	for _, test := range []struct {
		name      string
		fault     LifecycleStoreFault
		wantState lifecyclestate.LifecycleState
	}{
		{name: "confirmed-abort", fault: FaultLifecycleAfterEffectAbort, wantState: lifecyclestate.StatePrepareIntent},
		{name: "pre-state-indeterminate", fault: FaultLifecycleAfterEffectIndeterminatePreState, wantState: lifecyclestate.StatePrepareIntent},
		{name: "post-rename-indeterminate", fault: FaultLifecycleAfterEffectIndeterminate, wantState: lifecyclestate.StatePrepared},
	} {
		t.Run(test.name, func(t *testing.T) {
			harness := newLifecycleTransactionHarness(t, 1)
			record, _, _ := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
			permit, err := harness.store.BeginEffect(
				context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare,
			)
			if err != nil {
				t.Fatal(err)
			}
			result, _ := lifecyclestate.NewEffectResult(
				lifecyclestate.OperationPrepare, lifecyclestate.EffectResultApplied, lifecyclestate.BackendInstanceIdentity{},
			)
			harness.store.InjectLifecycleFailure(test.fault, errors.New("fault"))
			if _, err := harness.store.ConfirmEffect(context.Background(), permit, result); err == nil {
				t.Fatal("faulted confirmation succeeded")
			}
			if !harness.store.LifecycleRecoveryFenced() {
				t.Fatal("after-effect fault did not fence handle")
			}
			reopened, err := OpenFixedFileStoreV1(harness.path)
			if err != nil {
				t.Fatal(err)
			}
			after, _ := reopened.ReadLifecycle(context.Background(), harness.attemptIDs[0])
			if after.View().State != test.wantState {
				t.Fatalf("reopen state = %s, want %s", after.View().State, test.wantState)
			}
		})
	}

	harness := newLifecycleTransactionHarness(t, 1)
	record, _, _ := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
	permit, _ := harness.store.BeginEffect(
		context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare,
	)
	record, err := harness.store.RecordIndeterminate(context.Background(), permit, ClassificationRecoveryRequired)
	if err != nil || record.View().EffectStatus != lifecyclestate.EffectIndeterminate || record.View().State != lifecyclestate.StateUnresolved {
		t.Fatalf("record indeterminate = %#v, %v", record.View(), err)
	}
	reopened, err := OpenFixedFileStoreV1(harness.path)
	if err != nil {
		t.Fatal(err)
	}
	applied, _ := lifecyclestate.NewReconcileResult(
		lifecyclestate.OperationPrepare, lifecyclestate.ReconciliationApplied, lifecyclestate.BackendInstanceIdentity{},
	)
	reconciled, err := reopened.RecordReconciliation(context.Background(), harness.attemptIDs[0], record.View().RecordVersion, applied)
	if err != nil || reconciled.View().State != lifecyclestate.StatePrepared || reconciled.View().EffectID != permit.View().EffectID {
		t.Fatalf("same-ID reconciliation = %#v, %v", reconciled.View(), err)
	}
}

func TestDurableLifecycleReconciliationFaultsWrongOwnerAndRecoverySort(t *testing.T) {
	harness := newLifecycleTransactionHarness(t, 2)
	first, _, _ := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
	permit, _ := harness.store.BeginEffect(
		context.Background(), harness.attemptIDs[0], first.View().RecordVersion, lifecyclestate.OperationPrepare,
	)
	if ids, err := harness.store.RecoveryAttemptIDs(context.Background()); err != nil || len(ids) != 2 || bytes.Compare(ids[0][:], ids[1][:]) >= 0 {
		t.Fatalf("sorted recovery set = %x, %v", ids, err)
	}

	otherOwner, err := OpenFixedFileStoreV1WithOptions(harness.path, FixedFileStoreV1Options{
		EffectIDs: &lifecycleEffectIDSequence{next: 100}, OwnerSessionID: lifecycleOwnerID(t, 0x77),
	})
	if err != nil {
		t.Fatal(err)
	}
	result, _ := lifecyclestate.NewEffectResult(
		lifecyclestate.OperationPrepare, lifecyclestate.EffectResultApplied, lifecyclestate.BackendInstanceIdentity{},
	)
	if _, err := otherOwner.ConfirmEffect(context.Background(), permit, result); !errors.Is(err, ErrLifecycleBindingMismatch) {
		t.Fatalf("wrong-owner permit error = %v", err)
	}

	indeterminate, err := harness.store.RecordIndeterminate(context.Background(), permit, ClassificationRecoveryRequired)
	if err != nil {
		t.Fatal(err)
	}
	unknown, _ := lifecyclestate.NewReconcileResult(
		lifecyclestate.OperationPrepare, lifecyclestate.ReconciliationUnknown, lifecyclestate.BackendInstanceIdentity{},
	)
	harness.store.InjectLifecycleFailure(FaultLifecycleReconciliationEligibilityAbort, errors.New("fault"))
	if _, err := harness.store.RecordReconciliation(
		context.Background(), harness.attemptIDs[0], indeterminate.View().RecordVersion, unknown,
	); err == nil {
		t.Fatal("eligibility abort succeeded")
	}
	unchanged, _ := harness.store.ReadLifecycle(context.Background(), harness.attemptIDs[0])
	if unchanged.View().RecordVersion != indeterminate.View().RecordVersion || unchanged.View().AutomaticRecoveryCount != 0 {
		t.Fatal("eligibility abort changed reconciliation state")
	}

	harness.store.InjectLifecycleFailure(FaultLifecycleReconciliationResultAbort, errors.New("fault"))
	if _, err := harness.store.RecordReconciliation(
		context.Background(), harness.attemptIDs[0], unchanged.View().RecordVersion, unknown,
	); err == nil {
		t.Fatal("result abort succeeded")
	}
	if !harness.store.LifecycleRecoveryFenced() {
		t.Fatal("lost reconciliation result did not fence handle")
	}
	reopened, err := OpenFixedFileStoreV1(harness.path)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := reopened.ReadLifecycle(context.Background(), harness.attemptIDs[0])
	if after.View().AutomaticRecoveryCount != 1 || after.View().LastReconciliation != lifecyclestate.ReconciliationNone {
		t.Fatalf("eligibility checkpoint after result loss = %#v", after.View())
	}
}

func TestDurableLifecycleEffectIDCollisionAndBindingMismatchFailClosed(t *testing.T) {
	harness := newLifecycleTransactionHarness(t, 2)
	first, _, _ := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
	firstPermit, err := harness.store.BeginEffect(
		context.Background(), harness.attemptIDs[0], first.View().RecordVersion, lifecyclestate.OperationPrepare,
	)
	if err != nil {
		t.Fatal(err)
	}
	second, _, _ := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[1], harness.backend)
	harness.effectIDs.collision = firstPermit.View().EffectID
	before := mustReadFile(t, harness.path)
	if _, err := harness.store.BeginEffect(
		context.Background(), harness.attemptIDs[1], second.View().RecordVersion, lifecyclestate.OperationPrepare,
	); err == nil {
		t.Fatal("effect-ID collision accepted")
	}
	if !bytes.Equal(before, mustReadFile(t, harness.path)) {
		t.Fatal("effect-ID collision changed state")
	}

	wrongView := harness.backend.View()
	wrongView.ImplementationIdentityDigest[0] ^= 0xff
	wrongBackend, err := lifecyclestate.NewBackendBinding(wrongView)
	if err != nil {
		t.Fatal(err)
	}
	before = mustReadFile(t, harness.path)
	if _, _, err := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[1], wrongBackend); err == nil {
		t.Fatal("mismatched established backend binding replay accepted")
	}
	if !harness.store.LifecycleRecoveryFenced() || !bytes.Equal(before, mustReadFile(t, harness.path)) {
		t.Fatal("binding mismatch did not enter no-rewrite repair fence")
	}
}

func TestDurableLifecycleEstablishmentFaultMatrix(t *testing.T) {
	for _, test := range []struct {
		name       string
		fault      LifecycleStoreFault
		wantRecord bool
		wantFence  bool
	}{
		{name: "confirmed-abort", fault: FaultLifecycleEnsureAbort},
		{name: "indeterminate-pre-state", fault: FaultLifecycleEnsureIndeterminatePreState, wantFence: true},
		{name: "indeterminate-post-rename", fault: FaultLifecycleEnsureIndeterminate, wantRecord: true, wantFence: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			harness := newLifecycleTransactionHarness(t, 1)
			before := mustReadFile(t, harness.path)
			harness.store.InjectLifecycleFailure(test.fault, errors.New("fault"))
			if _, _, err := harness.store.EnsureLifecycle(
				context.Background(), harness.attemptIDs[0], harness.backend,
			); err == nil {
				t.Fatal("faulted lifecycle establishment succeeded")
			}
			if harness.store.LifecycleRecoveryFenced() != test.wantFence {
				t.Fatalf("fence = %v, want %v", harness.store.LifecycleRecoveryFenced(), test.wantFence)
			}
			if !test.wantRecord && !bytes.Equal(before, mustReadFile(t, harness.path)) {
				t.Fatal("pre-rename establishment fault changed state")
			}
			reopened, err := OpenFixedFileStoreV1(harness.path)
			if err != nil {
				t.Fatal(err)
			}
			_, err = reopened.ReadLifecycle(context.Background(), harness.attemptIDs[0])
			if (err == nil) != test.wantRecord {
				t.Fatalf("reopen record error = %v, want record %v", err, test.wantRecord)
			}
		})
	}
}

func TestDurableLifecycleIndeterminateFaultMatrix(t *testing.T) {
	for _, test := range []struct {
		name      string
		fault     LifecycleStoreFault
		wantState lifecyclestate.LifecycleState
	}{
		{name: "confirmed-abort", fault: FaultLifecycleRecordIndeterminateAbort, wantState: lifecyclestate.StatePrepareIntent},
		{name: "indeterminate-pre-state", fault: FaultLifecycleRecordIndeterminatePreState, wantState: lifecyclestate.StatePrepareIntent},
		{name: "indeterminate-post-rename", fault: FaultLifecycleRecordIndeterminateCommitIndeterminate, wantState: lifecyclestate.StateUnresolved},
	} {
		t.Run(test.name, func(t *testing.T) {
			harness := newLifecycleTransactionHarness(t, 1)
			record, _, _ := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
			permit, _ := harness.store.BeginEffect(
				context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare,
			)
			harness.store.InjectLifecycleFailure(test.fault, errors.New("fault"))
			if _, err := harness.store.RecordIndeterminate(
				context.Background(), permit, ClassificationRecoveryRequired,
			); err == nil {
				t.Fatal("faulted indeterminate checkpoint succeeded")
			}
			if !harness.store.LifecycleRecoveryFenced() {
				t.Fatal("failed indeterminate checkpoint did not fence handle")
			}
			reopened, err := OpenFixedFileStoreV1(harness.path)
			if err != nil {
				t.Fatal(err)
			}
			after, _ := reopened.ReadLifecycle(context.Background(), harness.attemptIDs[0])
			if after.View().State != test.wantState {
				t.Fatalf("reopen state = %s, want %s", after.View().State, test.wantState)
			}
		})
	}
}

func TestDurableLifecycleReconciliationPrePostRenameMatrix(t *testing.T) {
	for _, test := range []struct {
		name       string
		fault      LifecycleStoreFault
		wantCount  uint64
		wantResult lifecyclestate.ReconciliationStatus
	}{
		{name: "eligibility-pre", fault: FaultLifecycleReconciliationEligibilityIndeterminatePreState, wantResult: lifecyclestate.ReconciliationNone},
		{name: "eligibility-post", fault: FaultLifecycleReconciliationEligibilityIndeterminate, wantCount: 1, wantResult: lifecyclestate.ReconciliationNone},
		{name: "result-pre", fault: FaultLifecycleReconciliationResultIndeterminatePreState, wantCount: 1, wantResult: lifecyclestate.ReconciliationNone},
		{name: "result-post", fault: FaultLifecycleReconciliationResultIndeterminate, wantCount: 1, wantResult: lifecyclestate.ReconciliationUnknown},
	} {
		t.Run(test.name, func(t *testing.T) {
			harness := newLifecycleTransactionHarness(t, 1)
			record, _, _ := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
			permit, _ := harness.store.BeginEffect(
				context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare,
			)
			record, _ = harness.store.RecordIndeterminate(context.Background(), permit, ClassificationRecoveryRequired)
			unknown, _ := lifecyclestate.NewReconcileResult(
				lifecyclestate.OperationPrepare, lifecyclestate.ReconciliationUnknown, lifecyclestate.BackendInstanceIdentity{},
			)
			harness.store.InjectLifecycleFailure(test.fault, errors.New("fault"))
			if _, err := harness.store.RecordReconciliation(
				context.Background(), harness.attemptIDs[0], record.View().RecordVersion, unknown,
			); err == nil {
				t.Fatal("faulted reconciliation succeeded")
			}
			if !harness.store.LifecycleRecoveryFenced() {
				t.Fatal("indeterminate reconciliation did not fence handle")
			}
			reopened, err := OpenFixedFileStoreV1(harness.path)
			if err != nil {
				t.Fatal(err)
			}
			after, _ := reopened.ReadLifecycle(context.Background(), harness.attemptIDs[0])
			if uint64(after.View().AutomaticRecoveryCount) != test.wantCount || after.View().LastReconciliation != test.wantResult {
				t.Fatalf("reopen reconciliation = count %d result %s", after.View().AutomaticRecoveryCount, after.View().LastReconciliation)
			}
		})
	}
}

func TestDurableLifecycleCreateIdentityCollisionQuarantinesWithoutAdoption(t *testing.T) {
	harness := newLifecycleTransactionHarness(t, 2)
	instance, err := lifecyclestate.NewBackendInstanceIdentity(lifecyclestate.BackendInstanceFake, []byte("same-instance"))
	if err != nil {
		t.Fatal(err)
	}
	createdPermits := make([]lifecyclestate.EffectPermit, 2)
	for index, attemptID := range harness.attemptIDs {
		record, _, _ := harness.store.EnsureLifecycle(context.Background(), attemptID, harness.backend)
		prepare, _ := harness.store.BeginEffect(context.Background(), attemptID, record.View().RecordVersion, lifecyclestate.OperationPrepare)
		preparedResult, _ := lifecyclestate.NewEffectResult(
			lifecyclestate.OperationPrepare, lifecyclestate.EffectResultApplied, lifecyclestate.BackendInstanceIdentity{},
		)
		record, err = harness.store.ConfirmEffect(context.Background(), prepare, preparedResult)
		if err != nil {
			t.Fatal(err)
		}
		createdPermits[index], err = harness.store.BeginEffect(
			context.Background(), attemptID, record.View().RecordVersion, lifecyclestate.OperationCreate,
		)
		if err != nil {
			t.Fatal(err)
		}
	}
	createResult, _ := lifecyclestate.NewEffectResult(lifecyclestate.OperationCreate, lifecyclestate.EffectResultApplied, instance)
	if _, err := harness.store.ConfirmEffect(context.Background(), createdPermits[0], createResult); err != nil {
		t.Fatal(err)
	}
	if _, err := harness.store.ConfirmEffect(context.Background(), createdPermits[1], createResult); !errors.Is(err, ErrLifecycleBindingMismatch) {
		t.Fatalf("identity collision result = %v", err)
	}
	quarantined, err := harness.store.ReadLifecycle(context.Background(), harness.attemptIDs[1])
	if err != nil {
		t.Fatal(err)
	}
	if quarantined.View().State != lifecyclestate.StateQuarantined || quarantined.View().Instance.Present() {
		t.Fatalf("colliding instance was adopted: %#v", quarantined.View())
	}
	reopened, err := OpenFixedFileStoreV1(harness.path)
	if err != nil {
		t.Fatal(err)
	}
	recovery, err := reopened.RecoveryAttemptIDs(context.Background())
	if err != nil || len(recovery) != 2 {
		t.Fatalf("quarantine recovery set = %x, %v", recovery, err)
	}
}

// TestDurableLifecycleCompleteReconciliationQuarantinesOnInstanceValidationFailure
// exercises validateReconciledInstanceLocked's quarantine branch in
// CompleteReconciliation directly (lifecycle_store.go), covering both
// dispositions that can trip it while landing the record in
// StateQuarantined: a ReconciliationPresent report with no recorded
// instance, and a ReconciliationIdentityMismatch report that conflicts with
// the recorded instance. Each case must set the store's
// TrustRepairRequired/AttemptsDisabled trust-fencing fields, mirroring the
// already-tested ConfirmEffect-side identity collision.
//
// The third disposition named in issue #271 - a duplicate instance observed
// via a ReconciliationApplied create report - is deliberately not covered
// here as a passing "quarantines" case. It uncovers a real defect; see
// TestDurableLifecycleCompleteReconciliationDuplicateInstanceOnAppliedFailsClosedInstead
// below and the accompanying report.
func TestDurableLifecycleCompleteReconciliationQuarantinesOnInstanceValidationFailure(t *testing.T) {
	assertQuarantinedTrustState := func(t *testing.T, harness lifecycleTransactionHarness) {
		t.Helper()
		if harness.store.state.TrustPhase != TrustRepairRequired ||
			harness.store.state.TrustReason != "lifecycle-instance-identity-mismatch" ||
			!harness.store.state.AttemptsDisabled {
			t.Fatalf(
				"quarantine trust state = phase %s, reason %q, disabled %v",
				harness.store.state.TrustPhase, harness.store.state.TrustReason, harness.store.state.AttemptsDisabled,
			)
		}
	}

	t.Run("present-without-recorded-instance-quarantines", func(t *testing.T) {
		harness := newLifecycleTransactionHarness(t, 1)
		attemptID := harness.attemptIDs[0]
		record, _, err := harness.store.EnsureLifecycle(context.Background(), attemptID, harness.backend)
		if err != nil {
			t.Fatal(err)
		}
		permit, err := harness.store.BeginEffect(
			context.Background(), attemptID, record.View().RecordVersion, lifecyclestate.OperationPrepare,
		)
		if err != nil {
			t.Fatal(err)
		}
		record, err = harness.store.RecordIndeterminate(context.Background(), permit, ClassificationRecoveryRequired)
		if err != nil {
			t.Fatal(err)
		}
		reportedInstance, err := lifecyclestate.NewBackendInstanceIdentity(
			lifecyclestate.BackendInstanceFake, []byte("present-but-unrecorded"),
		)
		if err != nil {
			t.Fatal(err)
		}
		present, err := lifecyclestate.NewReconcileResult(
			lifecyclestate.OperationPrepare, lifecyclestate.ReconciliationPresent, reportedInstance,
		)
		if err != nil {
			t.Fatal(err)
		}
		quarantined, err := harness.store.RecordReconciliation(
			context.Background(), attemptID, record.View().RecordVersion, present,
		)
		if err != nil {
			t.Fatalf("present-status reconciliation returned an error instead of quarantining: %v", err)
		}
		view := quarantined.View()
		if view.State != lifecyclestate.StateQuarantined ||
			view.EffectStatus != lifecyclestate.EffectIndeterminate ||
			view.LastReconciliation != lifecyclestate.ReconciliationIdentityMismatch ||
			view.RecoveryFence != lifecyclestate.RecoveryFenceIdentityMismatch {
			t.Fatalf("present-status quarantine transition = %#v", view)
		}
		assertQuarantinedTrustState(t, harness)
	})

	t.Run("identity-mismatch-against-recorded-instance-quarantines", func(t *testing.T) {
		harness := newLifecycleTransactionHarness(t, 1)
		attemptID := harness.attemptIDs[0]
		record, _, err := harness.store.EnsureLifecycle(context.Background(), attemptID, harness.backend)
		if err != nil {
			t.Fatal(err)
		}
		preparePermit, err := harness.store.BeginEffect(
			context.Background(), attemptID, record.View().RecordVersion, lifecyclestate.OperationPrepare,
		)
		if err != nil {
			t.Fatal(err)
		}
		preparedResult, err := lifecyclestate.NewEffectResult(
			lifecyclestate.OperationPrepare, lifecyclestate.EffectResultApplied, lifecyclestate.BackendInstanceIdentity{},
		)
		if err != nil {
			t.Fatal(err)
		}
		record, err = harness.store.ConfirmEffect(context.Background(), preparePermit, preparedResult)
		if err != nil {
			t.Fatal(err)
		}
		createPermit, err := harness.store.BeginEffect(
			context.Background(), attemptID, record.View().RecordVersion, lifecyclestate.OperationCreate,
		)
		if err != nil {
			t.Fatal(err)
		}
		recordedInstance, err := lifecyclestate.NewBackendInstanceIdentity(
			lifecyclestate.BackendInstanceFake, []byte("originally-recorded"),
		)
		if err != nil {
			t.Fatal(err)
		}
		createdResult, err := lifecyclestate.NewEffectResult(
			lifecyclestate.OperationCreate, lifecyclestate.EffectResultApplied, recordedInstance,
		)
		if err != nil {
			t.Fatal(err)
		}
		record, err = harness.store.ConfirmEffect(context.Background(), createPermit, createdResult)
		if err != nil {
			t.Fatal(err)
		}
		startPermit, err := harness.store.BeginEffect(
			context.Background(), attemptID, record.View().RecordVersion, lifecyclestate.OperationStart,
		)
		if err != nil {
			t.Fatal(err)
		}
		record, err = harness.store.RecordIndeterminate(context.Background(), startPermit, ClassificationRecoveryRequired)
		if err != nil {
			t.Fatal(err)
		}
		conflictingInstance, err := lifecyclestate.NewBackendInstanceIdentity(
			lifecyclestate.BackendInstanceFake, []byte("actually-observed"),
		)
		if err != nil {
			t.Fatal(err)
		}
		mismatch, err := lifecyclestate.NewReconcileResult(
			lifecyclestate.OperationStart, lifecyclestate.ReconciliationIdentityMismatch, conflictingInstance,
		)
		if err != nil {
			t.Fatal(err)
		}
		quarantined, err := harness.store.RecordReconciliation(
			context.Background(), attemptID, record.View().RecordVersion, mismatch,
		)
		if err != nil {
			t.Fatalf("identity-mismatch reconciliation returned an error instead of quarantining: %v", err)
		}
		view := quarantined.View()
		if view.State != lifecyclestate.StateQuarantined ||
			view.EffectStatus != lifecyclestate.EffectIndeterminate ||
			view.LastReconciliation != lifecyclestate.ReconciliationIdentityMismatch ||
			view.RecoveryFence != lifecyclestate.RecoveryFenceIdentityMismatch {
			t.Fatalf("identity-mismatch quarantine transition = %#v", view)
		}
		assertQuarantinedTrustState(t, harness)
	})

}

// TestDurableLifecycleCompleteReconciliationDuplicateInstanceOnAppliedFailsClosedInstead
// documents a defect found while adding issue #271's coverage, rather than
// silently patching it: when CompleteReconciliation's validation catches a
// duplicate instance identity reported via a ReconciliationApplied create
// result (as opposed to the already-tested ConfirmEffect-side collision in
// TestDurableLifecycleCreateIdentityCollisionQuarantinesWithoutAdoption), its
// error-handling branch (lifecycle_store.go around L645-658) resets State,
// EffectStatus, LastReconciliation, RecoveryFence, and NextRecoveryAt to
// quarantine values, but never clears the Instance field that
// applyReconciliationResult/applyConfirmedEffect had already adopted onto
// the record a few lines earlier. The quarantined record is then committed
// still carrying the colliding instance identity, so the store's own
// whole-state invariant check (validateV1State, fixed_store_v1.go ~L629-633)
// rejects the candidate state with "fixed lifecycle state has a duplicate
// instance identity" before it can ever be persisted. The intended
// quarantine-without-adoption transition for this specific disposition can
// never durably complete: RecordReconciliation/CompleteReconciliation always
// return this internal invariant error instead of quarantining, and the
// second attempt is left stuck rather than fenced.
//
// This is a security-relevant gap (the deny-by-default quarantine defense
// silently fails to fire, and callers just get an opaque store error), but
// fixing it is out of scope for issue #271, which asks to stop and report
// rather than fix this defect under a test-only issue. This test pins today's
// actual (defective) behavior so any future change to it is deliberate; do
// not "fix" it by loosening this assertion without also fixing
// CompleteReconciliation to clear the colliding instance before quarantining.
func TestDurableLifecycleCompleteReconciliationDuplicateInstanceOnAppliedFailsClosedInstead(t *testing.T) {
	harness := newLifecycleTransactionHarness(t, 2)
	sharedInstance, err := lifecyclestate.NewBackendInstanceIdentity(
		lifecyclestate.BackendInstanceFake, []byte("shared-instance"),
	)
	if err != nil {
		t.Fatal(err)
	}
	preparedResult, err := lifecyclestate.NewEffectResult(
		lifecyclestate.OperationPrepare, lifecyclestate.EffectResultApplied, lifecyclestate.BackendInstanceIdentity{},
	)
	if err != nil {
		t.Fatal(err)
	}

	firstID := harness.attemptIDs[0]
	firstRecord, _, err := harness.store.EnsureLifecycle(context.Background(), firstID, harness.backend)
	if err != nil {
		t.Fatal(err)
	}
	firstPrepare, err := harness.store.BeginEffect(
		context.Background(), firstID, firstRecord.View().RecordVersion, lifecyclestate.OperationPrepare,
	)
	if err != nil {
		t.Fatal(err)
	}
	firstRecord, err = harness.store.ConfirmEffect(context.Background(), firstPrepare, preparedResult)
	if err != nil {
		t.Fatal(err)
	}
	firstCreate, err := harness.store.BeginEffect(
		context.Background(), firstID, firstRecord.View().RecordVersion, lifecyclestate.OperationCreate,
	)
	if err != nil {
		t.Fatal(err)
	}
	firstCreatedResult, err := lifecyclestate.NewEffectResult(
		lifecyclestate.OperationCreate, lifecyclestate.EffectResultApplied, sharedInstance,
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := harness.store.ConfirmEffect(context.Background(), firstCreate, firstCreatedResult); err != nil {
		t.Fatal(err)
	}

	secondID := harness.attemptIDs[1]
	secondRecord, _, err := harness.store.EnsureLifecycle(context.Background(), secondID, harness.backend)
	if err != nil {
		t.Fatal(err)
	}
	secondPrepare, err := harness.store.BeginEffect(
		context.Background(), secondID, secondRecord.View().RecordVersion, lifecyclestate.OperationPrepare,
	)
	if err != nil {
		t.Fatal(err)
	}
	secondRecord, err = harness.store.ConfirmEffect(context.Background(), secondPrepare, preparedResult)
	if err != nil {
		t.Fatal(err)
	}
	secondCreate, err := harness.store.BeginEffect(
		context.Background(), secondID, secondRecord.View().RecordVersion, lifecyclestate.OperationCreate,
	)
	if err != nil {
		t.Fatal(err)
	}
	secondRecord, err = harness.store.RecordIndeterminate(context.Background(), secondCreate, ClassificationRecoveryRequired)
	if err != nil {
		t.Fatal(err)
	}
	duplicateApplied, err := lifecyclestate.NewReconcileResult(
		lifecyclestate.OperationCreate, lifecyclestate.ReconciliationApplied, sharedInstance,
	)
	if err != nil {
		t.Fatal(err)
	}

	// BeginReconciliation's eligibility checkpoint (AutomaticRecoveryCount++,
	// LastReconciliation reset to none) is a separate durable commit from the
	// result commit that fails below; only the result side is broken here, so
	// this intentionally does not assert the file is byte-for-byte unchanged.
	_, err = harness.store.RecordReconciliation(
		context.Background(), secondID, secondRecord.View().RecordVersion, duplicateApplied,
	)
	if err == nil {
		t.Fatal("expected the known duplicate-instance-on-applied defect to still fail closed; " +
			"if this now succeeds, CompleteReconciliation was fixed to clear the colliding instance " +
			"before quarantining - update this test to assert the correct quarantine transition and " +
			"fold it back into TestDurableLifecycleCompleteReconciliationQuarantinesOnInstanceValidationFailure")
	}
	if !strings.Contains(err.Error(), "duplicate instance identity") {
		t.Fatalf("unexpected failure mode for the known defect, got: %v", err)
	}

	stuck, err := harness.store.ReadLifecycle(context.Background(), secondID)
	if err != nil {
		t.Fatal(err)
	}
	if stuck.View().State == lifecyclestate.StateQuarantined {
		t.Fatal("second attempt unexpectedly reached quarantine; the defect appears fixed, see comment above")
	}
	firstStillOwns, err := harness.store.ReadLifecycle(context.Background(), firstID)
	if err != nil {
		t.Fatal(err)
	}
	if !firstStillOwns.View().Instance.Present() || firstStillOwns.View().Instance.Digest() != sharedInstance.Digest() {
		t.Fatalf("original instance owner lost its binding: %#v", firstStillOwns.View())
	}
}

func TestDurableLifecycleCapacityReservationAndExpiredAuthorityRecovery(t *testing.T) {
	state := ordinaryInitialStateForV1()
	state.Attempts = make([]approvalattempt.ExecutionAttempt, MaxActiveLifecycleRecords)
	if active := activeLifecycleWork(state, nil); active != MaxActiveLifecycleRecords {
		t.Fatalf("matched missing-lifecycle reservation = %d", active)
	}
	state.Attempts = append(state.Attempts, approvalattempt.ExecutionAttempt{})
	if active := activeLifecycleWork(state, nil); active != MaxActiveLifecycleRecords+1 {
		t.Fatalf("cap-plus-one reservation = %d", active)
	}

	harness := newLifecycleTransactionHarness(t, 1)
	envelope := loadEnvelopeV1ForTest(t, harness.path)
	envelope.State.TimeHighWaterUnixSeconds = 1_785_457_000
	writeV1Envelope(t, harness.path, envelope)
	reopened, err := OpenFixedFileStoreV1(harness.path)
	if err != nil {
		t.Fatal(err)
	}
	ids, err := reopened.RecoveryAttemptIDs(context.Background())
	if err != nil || len(ids) != 1 || ids[0] != harness.attemptIDs[0] {
		t.Fatalf("expired registration/approval suppressed recovery: %x, %v", ids, err)
	}
}

func TestDurableLifecycleNotAppliedDispositionAndSameIDReconciliationRetry(t *testing.T) {
	t.Run("confirmed-not-applied", func(t *testing.T) {
		harness := newLifecycleTransactionHarness(t, 1)
		record, _, _ := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
		permit, _ := harness.store.BeginEffect(
			context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare,
		)
		notApplied, _ := lifecyclestate.NewEffectResult(
			lifecyclestate.OperationPrepare, lifecyclestate.EffectResultNotApplied, lifecyclestate.BackendInstanceIdentity{},
		)
		record, err := harness.store.ConfirmEffect(context.Background(), permit, notApplied)
		if err != nil {
			t.Fatal(err)
		}
		if record.View().State != lifecyclestate.StateUnresolved ||
			record.View().FirstFailure != lifecyclestate.FailureLifecycle ||
			!record.View().CleanupRequired {
			t.Fatalf("not-applied disposition = %#v", record.View())
		}
		if _, err := harness.store.BeginEffect(
			context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationCreate,
		); err == nil {
			t.Fatal("successor operation ran through not-applied disposition")
		}
	})

	t.Run("reconciled-not-applied-reuses-effect", func(t *testing.T) {
		harness := newLifecycleTransactionHarness(t, 1)
		record, _, _ := harness.store.EnsureLifecycle(context.Background(), harness.attemptIDs[0], harness.backend)
		permit, _ := harness.store.BeginEffect(
			context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare,
		)
		record, _ = harness.store.RecordIndeterminate(context.Background(), permit, ClassificationRecoveryRequired)
		notApplied, _ := lifecyclestate.NewReconcileResult(
			lifecyclestate.OperationPrepare, lifecyclestate.ReconciliationNotApplied, lifecyclestate.BackendInstanceIdentity{},
		)
		record, err := harness.store.RecordReconciliation(
			context.Background(), harness.attemptIDs[0], record.View().RecordVersion, notApplied,
		)
		if err != nil {
			t.Fatal(err)
		}
		retry, err := harness.store.BeginEffect(
			context.Background(), harness.attemptIDs[0], record.View().RecordVersion, lifecyclestate.OperationPrepare,
		)
		if err != nil {
			t.Fatal(err)
		}
		if retry.View().EffectID != permit.View().EffectID || harness.effectIDs.calls != 1 {
			t.Fatalf("not-applied retry allocated another effect: %x, calls %d", retry.View().EffectID, harness.effectIDs.calls)
		}
		applied, _ := lifecyclestate.NewEffectResult(
			lifecyclestate.OperationPrepare, lifecyclestate.EffectResultApplied, lifecyclestate.BackendInstanceIdentity{},
		)
		record, err = harness.store.ConfirmEffect(context.Background(), retry, applied)
		if err != nil || record.View().State != lifecyclestate.StatePrepared {
			t.Fatalf("same-ID retry confirmation = %#v, %v", record.View(), err)
		}
	})
}

func loadEnvelopeV1ForTest(t *testing.T, path string) diskEnvelopeV1 {
	t.Helper()
	data := mustReadFile(t, path)
	var envelope diskEnvelopeV1
	if err := json.Unmarshal(data, &envelope); err != nil {
		t.Fatal(err)
	}
	return envelope
}
