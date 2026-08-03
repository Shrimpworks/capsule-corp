package registeredlifecycle

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const ordinaryPlanPath = "../../../schemas/conformance/v0/execution-plan/ordinary.cbor"

var durableOperations = []Operation{
	OperationPrepare, OperationCreate, OperationStart,
	OperationObserve, OperationStop, OperationDestroy,
}

func TestDriveResolvesCommittedSliceBAttemptAndExposesNoSuccessResult(t *testing.T) {
	harness := newHarness(t, nil)
	if harness.backend.CreatesGuest() || harness.backend.Binding().View().CreatesGuest {
		t.Fatal("fake backend reported guest creation")
	}
	snapshot, err := harness.component.Drive(context.Background(), harness.attemptID)
	if err != nil {
		t.Fatalf("drive fake lifecycle: %v", err)
	}
	if snapshot.AttemptID != harness.attemptID || snapshot.RegistrationID != harness.registrationID ||
		snapshot.State != StateDestroyed || snapshot.CleanupRequired ||
		snapshot.LastReconciliation != lifecyclestate.ReconciliationAuthoritativelyAbsent {
		t.Fatalf("lifecycle disposition = %+v, want bound destroyed attempt", snapshot)
	}
	if snapshot.Failure != "" || snapshot.FailureAt != "" {
		t.Fatalf("ordinary fake lifecycle recorded a job result/failure: %+v", snapshot)
	}
	backend := harness.backend.Snapshot(harness.attemptID)
	if !backend.Prepared || !backend.Created || !backend.Started || !backend.Observed ||
		!backend.Stopped || !backend.Destroyed || backend.InstanceDigest != snapshot.InstanceDigest {
		t.Fatalf("incomplete fake lifecycle: %+v", backend)
	}
	for _, operation := range durableOperations {
		if backend.ApplicationCounts[operation] != 1 || backend.EffectIDs[operation].IsZero() {
			t.Fatalf("%s application/effect = %d/%x", operation, backend.ApplicationCounts[operation], backend.EffectIDs[operation])
		}
	}
	componentType := reflect.TypeOf((*Component)(nil))
	method, ok := componentType.MethodByName("Drive")
	if !ok || method.Type.NumIn() != 3 || method.Type.In(2) != reflect.TypeOf(approvalattempt.AttemptID{}) {
		t.Fatalf("Drive inputs = %v, want receiver, context, AttemptID only", method.Type)
	}
	if _, exists := componentType.MethodByName("Execute"); exists {
		t.Fatal("obsolete RegistrationID-keyed Execute method remains")
	}
}

func TestDriveRejectsMissingMutatedAndCrossLinkedAttemptsBeforePrepare(t *testing.T) {
	tests := []struct {
		name  string
		alter func(*registrationstate.CreatedAttempt)
		want  Classification
	}{
		{name: "wrong plan role binding", alter: func(created *registrationstate.CreatedAttempt) {
			created.PlanRoleBindings.RuntimeBundleManifestDigest[0] ^= 0xff
		}, want: ClassificationDomain},
		{name: "wrong attempt identity", alter: func(created *registrationstate.CreatedAttempt) {
			created.Attempt.AttemptID[0] ^= 0xff
		}, want: ClassificationBinding},
		{name: "non created attempt state", alter: func(created *registrationstate.CreatedAttempt) {
			created.Attempt.State = ""
		}, want: ClassificationBinding},
		{name: "cross linked approval", alter: func(created *registrationstate.CreatedAttempt) {
			created.Approval.ApprovalID[0] ^= 0xff
		}, want: ClassificationBinding},
		{name: "mutated exact plan", alter: func(created *registrationstate.CreatedAttempt) {
			created.Registration.ExactPlanBytes[len(created.Registration.ExactPlanBytes)-1] ^= 0xff
		}, want: ClassificationBinding},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			harness := newHarness(t, nil)
			component := harness.newComponent(t, &alteringResolver{base: harness.attempts, alter: test.alter}, nil)
			_, err := component.Drive(context.Background(), harness.attemptID)
			assertLifecycleClassification(t, err, test.want)
			if _, err := harness.store.ReadLifecycle(context.Background(), harness.attemptID); !errors.Is(err, registrationstate.ErrLifecycleNotFound) {
				t.Fatalf("invalid resolved attempt lifecycle read = %v", err)
			}
			if calls := harness.backend.Snapshot(harness.attemptID).CallCounts; len(calls) != 0 {
				t.Fatalf("invalid resolved attempt reached fake backend: %v", calls)
			}
		})
	}
	harness := newHarness(t, nil)
	unknown := harness.attemptID
	unknown[15] ^= 0xff
	_, err := harness.component.Drive(context.Background(), unknown)
	assertLifecycleClassification(t, err, ClassificationBinding)
	_, err = harness.component.Drive(context.Background(), approvalattempt.AttemptID{})
	assertLifecycleClassification(t, err, ClassificationBinding)
}

func TestRecoveryFencedAttemptStoreRefusesBeforePrepare(t *testing.T) {
	tests := []struct {
		name  string
		fault registrationstate.LifecycleStoreFault
		want  Classification
	}{
		{name: "confirmed-abort", fault: registrationstate.FaultLifecycleEnsureAbort, want: ClassificationLocalFailure},
		{name: "indeterminate-pre-state", fault: registrationstate.FaultLifecycleEnsureIndeterminatePreState, want: ClassificationRecoveryRequired},
		{name: "indeterminate-post-rename", fault: registrationstate.FaultLifecycleEnsureIndeterminate, want: ClassificationRecoveryRequired},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			harness := newHarness(t, nil)
			harness.store.InjectLifecycleFailure(test.fault, errors.New("simulated lifecycle creation fault"))
			_, err := harness.component.Drive(context.Background(), harness.attemptID)
			assertLifecycleClassification(t, err, test.want)
			if calls := harness.backend.Snapshot(harness.attemptID).CallCounts; len(calls) != 0 {
				t.Fatalf("faulted lifecycle creation reached fake backend: %v", calls)
			}
			harness.reopen(t, nil)
			recovered, recoverErr := harness.component.Recover(context.Background(), harness.attemptID)
			if recoverErr != nil || recovered.State != StateDestroyed || recovered.CleanupRequired {
				t.Fatalf("creation reopen recovery = %+v, %v", recovered, recoverErr)
			}
		})
	}
}

func TestConcurrentAndSequentialDriveNeverRedriveEffects(t *testing.T) {
	harness := newHarness(t, nil)
	const workers = 32
	var group sync.WaitGroup
	errorsSeen := make(chan error, workers)
	for range workers {
		group.Add(1)
		go func() {
			defer group.Done()
			_, err := harness.component.Drive(context.Background(), harness.attemptID)
			errorsSeen <- err
		}()
	}
	group.Wait()
	close(errorsSeen)
	for err := range errorsSeen {
		if err != nil {
			t.Fatalf("concurrent drive: %v", err)
		}
	}
	if _, err := harness.component.Drive(context.Background(), harness.attemptID); err != nil {
		t.Fatalf("sequential replay: %v", err)
	}
	backend := harness.backend.Snapshot(harness.attemptID)
	for _, operation := range durableOperations {
		if backend.ApplicationCounts[operation] != 1 || backend.CallCounts[operation] != 1 {
			t.Fatalf("%s calls/applications = %d/%d, want 1/1", operation, backend.CallCounts[operation], backend.ApplicationCounts[operation])
		}
	}
}

func TestTwoApprovalsForOneRegistrationHaveIndependentAttemptLifecycles(t *testing.T) {
	harness := newHarness(t, []fixtureSpec{{nonce: 0x66}, {nonce: 0x67, variant: 1}})
	second := harness.attemptIDs[1]
	firstSnapshot, err := harness.component.Drive(context.Background(), harness.attemptID)
	if err != nil {
		t.Fatal(err)
	}
	secondSnapshot, err := harness.component.Drive(context.Background(), second)
	if err != nil {
		t.Fatal(err)
	}
	if firstSnapshot.RegistrationID != secondSnapshot.RegistrationID ||
		firstSnapshot.AttemptID == secondSnapshot.AttemptID || firstSnapshot.ApprovalID == secondSnapshot.ApprovalID {
		t.Fatalf("attempt lifecycles conflated: first=%+v second=%+v", firstSnapshot, secondSnapshot)
	}
	if !harness.backend.Snapshot(harness.attemptID).Destroyed || !harness.backend.Snapshot(second).Destroyed {
		t.Fatal("independent fake instances did not both destroy")
	}
}

func TestOneApprovalCannotDriveTwoAttemptIDs(t *testing.T) {
	harness := newHarness(t, nil)
	if _, err := harness.component.Drive(context.Background(), harness.attemptID); err != nil {
		t.Fatal(err)
	}
	aliasID := attemptIDFor(99)
	component := harness.newComponent(t, aliasingResolver{base: harness.attempts, sourceID: harness.attemptID}, nil)
	_, err := component.Drive(context.Background(), aliasID)
	assertLifecycleClassification(t, err, ClassificationBinding)
	if calls := harness.backend.Snapshot(aliasID).CallCounts; len(calls) != 0 {
		t.Fatalf("approval alias reached fake backend: %v", calls)
	}
}

func TestFakeFaultMatrixEndsDestroyedOrUnresolved(t *testing.T) {
	for _, operation := range durableOperations {
		for _, moment := range []FaultMoment{FaultBeforeEffect, FaultAfterEffect} {
			t.Run(string(operation)+"/"+string(moment), func(t *testing.T) {
				harness := newHarness(t, nil)
				if err := harness.backend.InjectFault(harness.attemptID, operation, moment); err != nil {
					t.Fatal(err)
				}
				initial, err := harness.component.Drive(context.Background(), harness.attemptID)
				if err == nil {
					t.Fatal("injected fault returned nil error")
				}
				if moment == FaultAfterEffect {
					if initial.State != StateUnresolved || !initial.CleanupRequired {
						t.Fatalf("post-effect fault disposition = %+v", initial)
					}
					harness.reopen(t, nil)
					recovered, recoverErr := harness.component.Recover(context.Background(), harness.attemptID)
					if recovered.State != StateDestroyed || recovered.CleanupRequired || recoverErr == nil {
						t.Fatalf("recovered fault = %+v, %v", recovered, recoverErr)
					}
				} else if initial.State != StateDestroyed || initial.CleanupRequired {
					t.Fatalf("confirmed not-applied retry = %+v", initial)
				}
				backend := harness.backend.Snapshot(harness.attemptID)
				if backend.ApplicationCounts[operation] != 1 {
					t.Fatalf("%s applications = %d, want exactly 1", operation, backend.ApplicationCounts[operation])
				}
			})
		}
	}
}

func TestPostSideEffectInterruptionRecoversByAttemptAcrossComponentRestart(t *testing.T) {
	checkpoints := []Checkpoint{
		CheckpointAfterPrepareEffect, CheckpointAfterCreateEffect, CheckpointAfterStartEffect,
		CheckpointAfterObserveEffect, CheckpointAfterStopEffect, CheckpointAfterDestroyEffect,
	}
	for _, checkpoint := range checkpoints {
		t.Run(string(checkpoint), func(t *testing.T) {
			harness := newHarness(t, nil)
			triggered := false
			harness.component = harness.newComponent(t, harness.attempts, func(
				_ context.Context, actual Checkpoint, _ approvalattempt.AttemptID,
			) error {
				if actual == checkpoint && !triggered {
					triggered = true
					return errors.New("simulated process death")
				}
				return nil
			})
			interrupted, err := harness.component.Drive(context.Background(), harness.attemptID)
			assertLifecycleClassification(t, err, ClassificationLocalFailure)
			if interrupted.State == StateDestroyed {
				t.Fatal("checkpoint did not leave durable recovery work")
			}
			harness.reopen(t, nil)
			recovered, err := harness.component.Recover(context.Background(), harness.attemptID)
			if err != nil || recovered.State != StateDestroyed || recovered.CleanupRequired {
				t.Fatalf("recovered disposition = %+v, %v", recovered, err)
			}
			backend := harness.backend.Snapshot(harness.attemptID)
			for _, operation := range durableOperations {
				if backend.ApplicationCounts[operation] != 1 {
					t.Fatalf("%s applications = %d", operation, backend.ApplicationCounts[operation])
				}
			}
		})
	}
}

func TestUnknownRecoveryRemainsUnresolvedAndHonorsDurableBackoff(t *testing.T) {
	harness := newHarness(t, nil)
	harness.component = harness.newComponent(t, harness.attempts, oneShotCheckpoint(CheckpointAfterCreateEffect))
	if _, err := harness.component.Drive(context.Background(), harness.attemptID); err == nil {
		t.Fatal("expected simulated interruption")
	}
	if err := harness.backend.InjectFault(harness.attemptID, OperationReconcile, FaultBeforeEffect); err != nil {
		t.Fatal(err)
	}
	harness.reopen(t, nil)
	unresolved, err := harness.component.Recover(context.Background(), harness.attemptID)
	assertLifecycleClassification(t, err, ClassificationCleanupUnresolved)
	if unresolved.State != StateUnresolved || unresolved.AutomaticRecoveryCount != 1 || !unresolved.NextRecoveryAt.Present {
		t.Fatalf("unknown recovery = %+v", unresolved)
	}
	before := harness.backend.Snapshot(harness.attemptID).CallCounts[OperationReconcile]
	harness.reopen(t, nil)
	_, err = harness.component.Recover(context.Background(), harness.attemptID)
	assertLifecycleClassification(t, err, ClassificationCleanupUnresolved)
	after := harness.backend.Snapshot(harness.attemptID).CallCounts[OperationReconcile]
	if after != before {
		t.Fatalf("backoff redrove reconciliation: before=%d after=%d", before, after)
	}
	harness.clock.set(1_785_456_001)
	harness.reopen(t, nil)
	recovered, err := harness.component.Recover(context.Background(), harness.attemptID)
	assertLifecycleClassification(t, err, ClassificationCleanupUnresolved)
	if recovered.State != StateDestroyed || recovered.CleanupRequired {
		t.Fatalf("eligible recovery = %+v, %v", recovered, err)
	}
}

func TestStartupEnumerationDrivesCreatedAttemptOnceAndIgnoresExpiry(t *testing.T) {
	harness := newHarness(t, nil)
	harness.clock.set(1_785_456_301)
	results, err := harness.component.RecoverCreatedAttempts(context.Background())
	if err != nil || len(results) != 1 || results[0].State != StateDestroyed {
		t.Fatalf("startup recovery = %+v, %v", results, err)
	}
	before := harness.backend.Snapshot(harness.attemptID).CallCounts
	harness.reopen(t, nil)
	results, err = harness.component.RecoverCreatedAttempts(context.Background())
	if err != nil || len(results) != 0 {
		t.Fatalf("terminal startup replay = %+v, %v", results, err)
	}
	if after := harness.backend.Snapshot(harness.attemptID).CallCounts; !reflect.DeepEqual(before, after) {
		t.Fatalf("startup replay redrove effects: before=%v after=%v", before, after)
	}
}

func TestDefensiveCopiesSnapshotsAndFixedErrors(t *testing.T) {
	harness := newHarness(t, nil)
	resolved, err := harness.attempts.ResolveCreated(context.Background(), harness.attemptID)
	if err != nil {
		t.Fatal(err)
	}
	resolved.Registration.ExactPlanBytes[0] ^= 0xff
	resolved.PlanRoleBindings.ProfileReviewAttestationDigests[0][0] ^= 0xff
	again, err := harness.attempts.ResolveCreated(context.Background(), harness.attemptID)
	if err != nil || !bytes.Equal(again.Registration.ExactPlanBytes, harness.plan) ||
		again.PlanRoleBindings.ProfileReviewAttestationDigests[0][0] == resolved.PlanRoleBindings.ProfileReviewAttestationDigests[0][0] {
		t.Fatal("created-attempt resolution did not defensively copy bindings")
	}
	if err := harness.backend.InjectFault(harness.attemptID, OperationObserve, FaultAfterEffect); err != nil {
		t.Fatal(err)
	}
	_, err = harness.component.Drive(context.Background(), harness.attemptID)
	if err == nil || bytes.Contains([]byte(err.Error()), harness.plan) || bytes.Contains([]byte(err.Error()), []byte("injected")) {
		t.Fatalf("fixed error leaked content: %q", err)
	}
	first := harness.backend.Snapshot(harness.attemptID)
	first.CallCounts[OperationPrepare] = 999
	first.EffectIDs[OperationPrepare] = lifecyclestate.EffectID{}
	second := harness.backend.Snapshot(harness.attemptID)
	if second.CallCounts[OperationPrepare] != 1 || second.EffectIDs[OperationPrepare].IsZero() {
		t.Fatal("backend snapshot did not defensively copy maps")
	}
}

func TestLifecycleStoreRefusesBeyondAttemptKeyedCapacity(t *testing.T) {
	harness := newHarness(t, nil)
	if _, ok := any(harness.store).(registrationstate.DurableLifecycleStore); !ok {
		t.Fatal("active lifecycle path is not the durable transaction port")
	}
	if _, err := New(Options{
		Attempts: harness.attempts, Store: harness.store, Backend: harness.backend,
		Coordinator: nil, Clock: harness.clock,
	}); err == nil {
		t.Fatal("constructor accepted a missing injected owner coordinator")
	}
	wrongOwner := ownerID(t, 0xee)
	wrongCoordinator, _ := NewCoordinator(wrongOwner)
	if _, err := New(Options{
		Attempts: harness.attempts, Store: harness.store, Backend: harness.backend,
		Coordinator: wrongCoordinator, Clock: harness.clock,
	}); err == nil {
		t.Fatal("constructor accepted a mismatched owner session")
	}
}

func TestLifecycleIntentCommitsBeforeEveryFakeEffect(t *testing.T) {
	harness := newHarness(t, nil)
	seen := make(map[Operation]lifecyclestate.EffectID)
	harness.component = harness.newComponent(t, harness.attempts, func(
		_ context.Context, checkpoint Checkpoint, attemptID approvalattempt.AttemptID,
	) error {
		operation := operationForCheckpoint(checkpoint)
		reopened, err := registrationstate.OpenFixedFileStoreV1(harness.path)
		if err != nil {
			t.Fatalf("reopen at %s: %v", checkpoint, err)
		}
		record, err := reopened.ReadLifecycle(context.Background(), attemptID)
		if err != nil {
			t.Fatal(err)
		}
		view := record.View()
		if view.Operation != durableOperation(operation) || view.EffectStatus != lifecyclestate.EffectIntent || view.EffectID.IsZero() {
			t.Fatalf("%s durable pre-effect state = %#v", operation, view)
		}
		seen[operation] = view.EffectID
		return nil
	})
	if _, err := harness.component.Drive(context.Background(), harness.attemptID); err != nil {
		t.Fatal(err)
	}
	if len(seen) != len(durableOperations) {
		t.Fatalf("observed %d durable intents, want %d", len(seen), len(durableOperations))
	}
}

func TestIntentConfirmedAbortCallsNoEffect(t *testing.T) {
	for index, target := range durableOperations {
		t.Run(string(target), func(t *testing.T) {
			harness := newHarness(t, nil)
			if index == 0 {
				harness.store.InjectLifecycleFailure(registrationstate.FaultLifecycleIntentAbort, errors.New("abort"))
			} else {
				predecessor := durableOperations[index-1]
				harness.component = harness.newComponent(t, harness.attempts, func(
					_ context.Context, checkpoint Checkpoint, _ approvalattempt.AttemptID,
				) error {
					if operationForCheckpoint(checkpoint) == predecessor {
						harness.store.InjectLifecycleFailure(registrationstate.FaultLifecycleIntentAbort, errors.New("abort"))
					}
					return nil
				})
			}
			_, err := harness.component.Drive(context.Background(), harness.attemptID)
			if err == nil {
				t.Fatal("intent abort returned nil")
			}
			backend := harness.backend.Snapshot(harness.attemptID)
			if backend.CallCounts[target] != 0 || backend.ApplicationCounts[target] != 0 {
				t.Fatalf("%s reached adapter: %+v", target, backend)
			}
		})
	}
}

func TestConfirmedEffectCommitFailureFencesUntilReopen(t *testing.T) {
	faults := []struct {
		name  string
		fault registrationstate.LifecycleStoreFault
		want  Classification
	}{
		{name: "confirmed-abort", fault: registrationstate.FaultLifecycleAfterEffectAbort, want: ClassificationLocalFailure},
		{name: "commit-indeterminate", fault: registrationstate.FaultLifecycleAfterEffectIndeterminate, want: ClassificationRecoveryRequired},
	}
	for _, fault := range faults {
		for _, target := range durableOperations {
			t.Run(fault.name+"/"+string(target), func(t *testing.T) {
				harness := newHarness(t, nil)
				harness.component = harness.newComponent(t, harness.attempts, func(
					_ context.Context, checkpoint Checkpoint, _ approvalattempt.AttemptID,
				) error {
					if operationForCheckpoint(checkpoint) == target {
						harness.store.InjectLifecycleFailure(fault.fault, errors.New("commit fault"))
					}
					return nil
				})
				_, err := harness.component.Drive(context.Background(), harness.attemptID)
				assertLifecycleClassification(t, err, fault.want)
				if !harness.store.LifecycleRecoveryFenced() || harness.backend.Snapshot(harness.attemptID).ApplicationCounts[target] != 1 {
					t.Fatal("confirmed post-effect commit failure did not fence after one application")
				}
				harness.reopen(t, nil)
				recovered, recoverErr := harness.component.Recover(context.Background(), harness.attemptID)
				if recovered.State != StateDestroyed || recovered.CleanupRequired || recoverErr != nil {
					t.Fatalf("reopen recovery = %+v, %v", recovered, recoverErr)
				}
				if harness.backend.Snapshot(harness.attemptID).ApplicationCounts[target] != 1 {
					t.Fatalf("%s effect replayed", target)
				}
			})
		}
	}
}

func TestIntentIndeterminateReopensAndRetriesSameEffectID(t *testing.T) {
	for index, target := range durableOperations {
		t.Run(string(target), func(t *testing.T) {
			harness := newHarness(t, nil)
			if index == 0 {
				harness.store.InjectLifecycleFailure(registrationstate.FaultLifecycleIntentIndeterminate, errors.New("indeterminate"))
			} else {
				predecessor := durableOperations[index-1]
				harness.component = harness.newComponent(t, harness.attempts, func(
					_ context.Context, checkpoint Checkpoint, _ approvalattempt.AttemptID,
				) error {
					if operationForCheckpoint(checkpoint) == predecessor {
						harness.store.InjectLifecycleFailure(registrationstate.FaultLifecycleIntentIndeterminate, errors.New("indeterminate"))
					}
					return nil
				})
			}
			_, err := harness.component.Drive(context.Background(), harness.attemptID)
			assertLifecycleClassification(t, err, ClassificationRecoveryRequired)
			if harness.backend.Snapshot(harness.attemptID).CallCounts[target] != 0 {
				t.Fatalf("indeterminate intent called %s", target)
			}
			durable, readErr := harness.store.ReadLifecycle(context.Background(), harness.attemptID)
			if readErr != nil || durable.View().Operation != durableOperation(target) || durable.View().EffectID.IsZero() {
				t.Fatalf("durable intent = %#v, %v", durable.View(), readErr)
			}
			effectID := durable.View().EffectID
			harness.reopen(t, nil)
			recovered, recoverErr := harness.component.Recover(context.Background(), harness.attemptID)
			assertLifecycleClassification(t, recoverErr, ClassificationLifecycleFailure)
			backend := harness.backend.Snapshot(harness.attemptID)
			if recovered.State != StateDestroyed || backend.ApplicationCounts[target] != 1 || backend.EffectIDs[target] != effectID {
				t.Fatalf("same-ID intent recovery = %+v backend=%+v", recovered, backend)
			}
		})
	}
}

func TestReconciliationResultCommitFailureRepeatsObservationNotEffect(t *testing.T) {
	harness := newHarness(t, nil)
	harness.component = harness.newComponent(t, harness.attempts, oneShotCheckpoint(CheckpointAfterCreateEffect))
	if _, err := harness.component.Drive(context.Background(), harness.attemptID); err == nil {
		t.Fatal("expected process-death checkpoint")
	}
	harness.reopen(t, nil)
	harness.store.InjectLifecycleFailure(registrationstate.FaultLifecycleReconciliationResultAbort, errors.New("result abort"))
	_, err := harness.component.Recover(context.Background(), harness.attemptID)
	assertLifecycleClassification(t, err, ClassificationLocalFailure)
	before := harness.backend.Snapshot(harness.attemptID)
	if before.CallCounts[OperationReconcile] != 1 || before.ApplicationCounts[OperationCreate] != 1 {
		t.Fatalf("first reconciliation = %+v", before)
	}
	harness.reopen(t, nil)
	recovered, recoverErr := harness.component.Recover(context.Background(), harness.attemptID)
	if recoverErr != nil || recovered.State != StateDestroyed {
		t.Fatalf("result-loss reopen = %+v, %v", recovered, recoverErr)
	}
	after := harness.backend.Snapshot(harness.attemptID)
	if after.ApplicationCounts[OperationCreate] != 1 || after.CallCounts[OperationReconcile] != 3 {
		t.Fatalf("result loss redrove effect or skipped observation: %+v", after)
	}
}

func TestLostEffectResponseReconcilesSameEffectIDAndExactInstance(t *testing.T) {
	harness := newHarness(t, nil)
	if err := harness.backend.InjectFault(harness.attemptID, OperationCreate, FaultAfterEffect); err != nil {
		t.Fatal(err)
	}
	unresolved, err := harness.component.Drive(context.Background(), harness.attemptID)
	assertLifecycleClassification(t, err, ClassificationLifecycleFailure)
	if unresolved.EffectID.IsZero() || unresolved.EffectID != harness.backend.Snapshot(harness.attemptID).EffectIDs[OperationCreate] {
		t.Fatalf("lost response effect identity = %x", unresolved.EffectID)
	}
	harness.reopen(t, nil)
	recovered, _ := harness.component.Recover(context.Background(), harness.attemptID)
	backend := harness.backend.Snapshot(harness.attemptID)
	if recovered.State != StateDestroyed || backend.ApplicationCounts[OperationCreate] != 1 ||
		recovered.InstanceDigest != backend.InstanceDigest {
		t.Fatalf("same-effect exact-instance recovery = %+v backend=%+v", recovered, backend)
	}
}

func TestMissingAndMismatchedBackendStateFailClosed(t *testing.T) {
	for _, mismatch := range []bool{false, true} {
		name := "missing"
		if mismatch {
			name = "identity-mismatch"
		}
		t.Run(name, func(t *testing.T) {
			harness := newHarness(t, nil)
			harness.component = harness.newComponent(t, harness.attempts, oneShotCheckpoint(CheckpointAfterStartEffect))
			if _, err := harness.component.Drive(context.Background(), harness.attemptID); err == nil {
				t.Fatal("expected interruption")
			}
			if mismatch {
				identity, err := lifecyclestate.NewBackendInstanceIdentity(lifecyclestate.BackendInstanceFake, []byte("wrong-instance"))
				if err != nil {
					t.Fatal(err)
				}
				harness.backend.replaceAttemptIdentity(harness.attemptID, identity)
			} else {
				harness.backend.forgetAttemptState(harness.attemptID)
			}
			harness.reopen(t, nil)
			snapshot, err := harness.component.Recover(context.Background(), harness.attemptID)
			if mismatch {
				assertLifecycleClassification(t, err, ClassificationTrustState)
				if snapshot.State != StateQuarantined {
					t.Fatalf("identity mismatch = %+v", snapshot)
				}
			} else {
				assertLifecycleClassification(t, err, ClassificationCleanupUnresolved)
				if snapshot.State != StateUnresolved || !snapshot.CleanupRequired {
					t.Fatalf("missing state = %+v", snapshot)
				}
			}
			backend := harness.backend.Snapshot(harness.attemptID)
			if backend.CallCounts[OperationStop] != 0 || backend.CallCounts[OperationDestroy] != 0 {
				t.Fatalf("mismatched object received cleanup effect: %+v", backend)
			}
		})
	}
}

func TestTrustedClockFailureBlocksForwardEffectsButAllowsCleanup(t *testing.T) {
	harness := newHarness(t, nil)
	harness.component = harness.newComponent(t, harness.attempts, oneShotCheckpoint(CheckpointAfterStartEffect))
	if _, err := harness.component.Drive(context.Background(), harness.attemptID); err == nil {
		t.Fatal("expected process-death checkpoint")
	}
	harness.clock.fail(errors.New("trusted clock unavailable"))
	harness.reopen(t, nil)
	recovered, err := harness.component.Recover(context.Background(), harness.attemptID)
	if err != nil || recovered.State != StateDestroyed || recovered.CleanupRequired {
		t.Fatalf("clock-failure cleanup = %+v, %v", recovered, err)
	}
	backend := harness.backend.Snapshot(harness.attemptID)
	if backend.ApplicationCounts[OperationObserve] != 0 ||
		backend.ApplicationCounts[OperationStop] != 1 || backend.ApplicationCounts[OperationDestroy] != 1 {
		t.Fatalf("clock-failure effects = %+v", backend.ApplicationCounts)
	}

	fresh := newHarness(t, nil)
	fresh.clock.fail(errors.New("trusted clock unavailable"))
	_, err = fresh.component.Drive(context.Background(), fresh.attemptID)
	assertLifecycleClassification(t, err, ClassificationLocalFailure)
	if calls := fresh.backend.Snapshot(fresh.attemptID).CallCounts; len(calls) != 0 {
		t.Fatalf("clock failure allowed a forward effect: %v", calls)
	}
}

type testHarness struct {
	component      *Component
	attempts       *registrationstate.ApprovalAttemptComponent
	clock          *testClock
	store          *registrationstate.FixedFileStoreV1
	backend        *FakeBackend
	path           string
	effectIDs      *effectIDSequence
	ownerNext      byte
	registrationID v0candidate.RegistrationID
	attemptID      approvalattempt.AttemptID
	attemptIDs     []approvalattempt.AttemptID
	plan           []byte
}

type fixtureSpec struct{ nonce, variant byte }

func newHarness(t *testing.T, specs []fixtureSpec) *testHarness {
	t.Helper()
	if len(specs) == 0 {
		specs = []fixtureSpec{{nonce: 0x66}}
	}
	plan := mustRead(t, ordinaryPlanPath)
	path := filepath.Join(t.TempDir(), "supervisor-state.json")
	stateStore, err := registrationstate.NewFixedFileStore(path, registrationstate.InitialState{
		InstallationID: repeated16[v0candidate.InstallationID](0x11),
		SupervisorID:   repeated16[v0candidate.SupervisorID](0x55), EpochSequence: 7,
		EpochDigest: repeated32[v0candidate.TrustEpochDigest](0x22), TrustPhase: registrationstate.TrustStable,
		TimeHighWaterUnixSeconds: 1_785_456_000,
	})
	if err != nil {
		t.Fatal(err)
	}
	clock := &testClock{value: 1_785_456_000}
	registrations, err := registrationstate.New(registrationstate.Options{
		Store: stateStore, Clock: clock,
		Identifiers: fixedRegistrationIDSource{id: repeated16[v0candidate.RegistrationID](0x33)},
	})
	if err != nil {
		t.Fatal(err)
	}
	issued, err := registrations.RegisterPlan(context.Background(), registrationstate.AuthenticatedCallContext{
		Authenticated: true, Role: registrationstate.CallerDaemon, Purpose: registrationstate.RegisterPlanPurpose,
	}, plan, ordinaryBindings())
	if err != nil {
		t.Fatal(err)
	}
	vectors := make([]approvalattempt.FixtureVector, len(specs))
	for index, spec := range specs {
		vectors[index] = lifecycleFixtureVector(t, issued.View(), spec.nonce, spec.variant)
	}
	verifier, err := approvalattempt.NewFixtureVerifier(vectors)
	if err != nil {
		t.Fatal(err)
	}
	attempts, err := registrationstate.NewApprovalAttempt(registrationstate.ApprovalAttemptOptions{
		Store: stateStore, Clock: clock, Verifier: verifier,
		ApprovalIdentifiers: &approvalIDSequence{next: 1}, AttemptIdentifiers: &attemptIDSequence{next: 1},
		Integrity: fixedIntegrity{assessedAt: 1_785_456_000},
	})
	if err != nil {
		t.Fatal(err)
	}
	attemptIDs := make([]approvalattempt.AttemptID, 0, len(vectors))
	for index, vector := range vectors {
		submission, submitErr := attempts.SubmitApproval(context.Background(), registrationstate.AuthenticatedCallContext{
			Authenticated: true, Role: registrationstate.CallerBroker, Purpose: registrationstate.SubmitApprovalPurpose,
		}, issued.View().RegistrationID, vector.EnvelopeBytes)
		if submitErr != nil {
			t.Fatalf("submit approval %d: %v", index, submitErr)
		}
		created, createErr := attempts.RequestAttempt(context.Background(), registrationstate.AuthenticatedCallContext{
			Authenticated: true, Role: registrationstate.CallerDaemon, Purpose: registrationstate.RequestAttemptPurpose,
		}, issued.View().RegistrationID, submission.Reference)
		if createErr != nil {
			t.Fatalf("create attempt %d: %v", index, createErr)
		}
		attemptIDs = append(attemptIDs, created.Reference.AttemptID())
	}
	if _, err := registrationstate.MigrateFixedFileStoreV0ToV1(
		context.Background(), path,
		registrationstate.V0ToV1MigrationOptions{Lock: offlineLock{}},
	); err != nil {
		t.Fatal(err)
	}
	effectIDs := &effectIDSequence{next: 1}
	store, err := registrationstate.OpenFixedFileStoreV1WithOptions(path, registrationstate.FixedFileStoreV1Options{
		EffectIDs: effectIDs, OwnerSessionID: ownerID(t, 0x40),
	})
	if err != nil {
		t.Fatal(err)
	}
	bindings := ordinaryBindings()
	backendBinding, err := lifecyclestate.NewBackendBinding(lifecyclestate.BackendBindingView{
		Kind: lifecyclestate.BackendFakeNoGuest, ProtocolVersion: lifecyclestate.FakeBackendProtocolVersion,
		ImplementationIdentityDigest:  lifecyclestate.BackendImplementationDigest(sha256.Sum256([]byte("fake-no-guest-e4"))),
		BackendConfigurationDigest:    bindings.BackendConfigurationDigest,
		BackendValidationRecordDigest: bindings.BackendValidationRecordDigest, CreatesGuest: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	backend, err := NewFakeBackend(backendBinding)
	if err != nil || backend.CreatesGuest() {
		t.Fatalf("new fake backend: %v", err)
	}
	harness := &testHarness{
		attempts: attempts, clock: clock, store: store, backend: backend, path: path,
		effectIDs: effectIDs, ownerNext: 0x41, registrationID: issued.View().RegistrationID,
		attemptID: attemptIDs[0], attemptIDs: attemptIDs, plan: plan,
	}
	harness.component = harness.newComponent(t, attempts, nil)
	return harness
}

func (harness *testHarness) newComponent(t *testing.T, resolver AttemptResolver, checkpoint CheckpointHook) *Component {
	t.Helper()
	coordinator, err := NewCoordinator(harness.store.OwnerSessionID())
	if err != nil {
		t.Fatal(err)
	}
	component, err := New(Options{
		Attempts: resolver, Store: harness.store, Backend: harness.backend,
		Coordinator: coordinator, Clock: harness.clock, Checkpoint: checkpoint,
	})
	if err != nil || component.backend.CreatesGuest() {
		t.Fatalf("new registered lifecycle: %v", err)
	}
	return component
}

func (harness *testHarness) reopen(t *testing.T, checkpoint CheckpointHook) {
	t.Helper()
	store, err := registrationstate.OpenFixedFileStoreV1WithOptions(harness.path, registrationstate.FixedFileStoreV1Options{
		EffectIDs: harness.effectIDs, OwnerSessionID: ownerID(t, harness.ownerNext),
	})
	harness.ownerNext++
	if err != nil {
		t.Fatal(err)
	}
	harness.store = store
	harness.component = harness.newComponent(t, harness.attempts, checkpoint)
}

type alteringResolver struct {
	base  AttemptResolver
	alter func(*registrationstate.CreatedAttempt)
}

func (resolver *alteringResolver) ResolveCreated(ctx context.Context, attemptID approvalattempt.AttemptID) (registrationstate.CreatedAttempt, error) {
	created, err := resolver.base.ResolveCreated(ctx, attemptID)
	if err == nil && resolver.alter != nil {
		resolver.alter(&created)
	}
	return created, err
}

type aliasingResolver struct {
	base     AttemptResolver
	sourceID approvalattempt.AttemptID
}

func (resolver aliasingResolver) ResolveCreated(ctx context.Context, attemptID approvalattempt.AttemptID) (registrationstate.CreatedAttempt, error) {
	created, err := resolver.base.ResolveCreated(ctx, resolver.sourceID)
	if err != nil {
		return registrationstate.CreatedAttempt{}, err
	}
	created.Attempt.AttemptID = attemptID
	created.Approval.ConsumedAttemptID = attemptID
	return created, nil
}

func oneShotCheckpoint(target Checkpoint) CheckpointHook {
	triggered := false
	return func(_ context.Context, checkpoint Checkpoint, _ approvalattempt.AttemptID) error {
		if checkpoint == target && !triggered {
			triggered = true
			return errors.New("simulated process death")
		}
		return nil
	}
}

func operationForCheckpoint(checkpoint Checkpoint) Operation {
	switch checkpoint {
	case CheckpointAfterPrepareEffect:
		return OperationPrepare
	case CheckpointAfterCreateEffect:
		return OperationCreate
	case CheckpointAfterStartEffect:
		return OperationStart
	case CheckpointAfterObserveEffect:
		return OperationObserve
	case CheckpointAfterStopEffect:
		return OperationStop
	case CheckpointAfterDestroyEffect:
		return OperationDestroy
	default:
		return ""
	}
}

type offlineLock struct{}

func (offlineLock) CheckOfflineMigrationLock(context.Context) error { return nil }

type effectIDSequence struct {
	mu   sync.Mutex
	next uint64
}

func (source *effectIDSequence) NewEffectID(context.Context) (lifecyclestate.EffectID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	value := make([]byte, 16)
	value[0] = 0xe4
	binary.BigEndian.PutUint64(value[8:], source.next)
	source.next++
	domain, err := lifecyclestate.NewDomainIdentifier(lifecyclestate.DomainEffectID, value)
	if err != nil {
		return lifecyclestate.EffectID{}, err
	}
	return lifecyclestate.NewEffectID(domain)
}

func ownerID(t *testing.T, value byte) lifecyclestate.OwnerSessionID {
	t.Helper()
	domain, err := lifecyclestate.NewDomainIdentifier(lifecyclestate.DomainOwnerSessionID, bytes.Repeat([]byte{value}, 16))
	if err != nil {
		t.Fatal(err)
	}
	owner, err := lifecyclestate.NewOwnerSessionID(domain)
	if err != nil {
		t.Fatal(err)
	}
	return owner
}

type testClock struct {
	mu    sync.Mutex
	value uint64
	err   error
}

func (clock *testClock) ObserveUnixSeconds(context.Context) (uint64, error) {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	return clock.value, clock.err
}
func (clock *testClock) set(value uint64) {
	clock.mu.Lock()
	clock.value = value
	clock.err = nil
	clock.mu.Unlock()
}
func (clock *testClock) fail(err error) {
	clock.mu.Lock()
	clock.err = err
	clock.mu.Unlock()
}

type fixedIntegrity struct{ assessedAt v0candidate.UInt53 }

func (integrity fixedIntegrity) Assess(_ context.Context, preflight registrationstate.IntegrityPreflight) (registrationstate.RuntimeIntegrityAssessment, error) {
	return registrationstate.RuntimeIntegrityAssessment{Preflight: preflight, AssessedAt: integrity.assessedAt, Permitted: true}, nil
}

type fixedRegistrationIDSource struct{ id v0candidate.RegistrationID }

func (source fixedRegistrationIDSource) NewRegistrationID(context.Context) (v0candidate.RegistrationID, error) {
	return source.id, nil
}

type approvalIDSequence struct {
	mu   sync.Mutex
	next uint64
}

func (source *approvalIDSequence) NewApprovalID(context.Context) (approvalattempt.ApprovalID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	var result approvalattempt.ApprovalID
	result[0] = 0xa1
	binary.BigEndian.PutUint64(result[8:], source.next)
	source.next++
	return result, nil
}

type attemptIDSequence struct {
	mu   sync.Mutex
	next uint64
}

func (source *attemptIDSequence) NewAttemptID(context.Context) (approvalattempt.AttemptID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	result := attemptIDFor(source.next)
	source.next++
	return result, nil
}

func attemptIDFor(value uint64) approvalattempt.AttemptID {
	var result approvalattempt.AttemptID
	result[0] = 0xb2
	binary.BigEndian.PutUint64(result[8:], value)
	return result
}

func lifecycleFixtureVector(t *testing.T, registration v0candidate.PlanRegistration, nonceByte, variant byte) approvalattempt.FixtureVector {
	t.Helper()
	envelope := mustRead(t, "../../../schemas/conformance/v0/approval-grant/ordinary.cose")
	payload := mustRead(t, "../../../schemas/conformance/v0/approval-grant/ordinary.payload.cbor")
	protected := mustRead(t, "../../../schemas/conformance/v0/approval-grant/ordinary.protected.cbor")
	if variant != 0 {
		envelope[len(envelope)-1] ^= variant
		payload[len(payload)-1] ^= variant
	}
	return approvalattempt.FixtureVector{
		EnvelopeBytes: envelope, PayloadBytes: payload, ProtectedHeaderBytes: protected,
		ProtectedKeyID: []byte("approval-test-key"),
		View: approvalattempt.ApprovalGrant{
			ObjectType: approvalattempt.ApprovalGrantObjectType, ObjectVersion: 0,
			InstallationID: registration.InstallationID, EpochDigest: registration.EpochDigest,
			RegistrationID: registration.RegistrationID, PlanDigest: registration.PlanDigest,
			SupervisorID: registration.SupervisorID, AttemptNonce: repeated16[approvalattempt.AttemptNonce](nonceByte),
			Purpose: approvalattempt.ApprovalGrantPurpose, Audience: approvalattempt.ApprovalGrantAudience,
			IssuedAt: 1_785_456_000, ExpiresAt: 1_785_456_300,
		},
		ResolvedEpochSequence: registration.EpochSequence,
		AuthorizationIdentity: repeated32[approvalattempt.ApprovalKeyAuthorizationIdentity](0x99), SignatureAccepted: true,
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	value, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return value
}

func ordinaryBindings() v0candidate.ExecutionPlanRoleBindings {
	return v0candidate.ExecutionPlanRoleBindings{
		InstallationID: repeated16[v0candidate.InstallationID](0x11), EpochDigest: repeated32[v0candidate.TrustEpochDigest](0x22),
		SourceManifestDigest:            hex32[v0candidate.SourceManifestDigest]("e5e09b2435baedf897526a89c698c0b0531437a69472372ae426f62d801fc171"),
		InlineInputDigest:               hex32[v0candidate.InlineInputDigest]("bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72"),
		RuntimeBundleManifestDigest:     repeated32[v0candidate.RuntimeBundleManifestDigest](0x55),
		ProfileReviewAttestationDigests: []v0candidate.ProfileReviewAttestationDigest{repeated32[v0candidate.ProfileReviewAttestationDigest](0x66), repeated32[v0candidate.ProfileReviewAttestationDigest](0x67)},
		ProfileRegistryEntryDigest:      repeated32[v0candidate.ProfileRegistryEntryDigest](0x77),
		BackendValidationRecordDigest:   repeated32[v0candidate.BackendValidationRecordDigest](0x88),
		BackendConfigurationDigest:      repeated32[v0candidate.BackendConfigurationDigest](0x99),
		TrustSnapshotDigest:             repeated32[v0candidate.TrustSnapshotDigest](0xaa), PolicyDecisionDigest: repeated32[v0candidate.PolicyDecisionDigest](0xbb),
	}
}

func assertLifecycleClassification(t *testing.T, err error, want Classification) {
	t.Helper()
	classification, ok := ErrorClassification(err)
	if !ok || classification != want {
		t.Fatalf("classification = %q (%t), want %q: %v", classification, ok, want, err)
	}
}

func repeated16[T ~[16]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}
func repeated32[T ~[32]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}
func hex32[T ~[32]byte](value string) T {
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != 32 {
		panic("invalid test digest")
	}
	var result T
	copy(result[:], decoded)
	return result
}
