package teardownpassive

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreparationOutcomesFailClosedBeforeCustody(t *testing.T) {
	for _, outcome := range []WriteOutcome{WriteFailed, WriteIndeterminate} {
		t.Run(string(outcome), func(t *testing.T) {
			model := newModel(t)
			pending, err := model.BeginPreparation(operationID(0x30))
			if err != nil {
				t.Fatal(err)
			}
			if decision := model.Decision(); decision.MayCreate || decision.MayStart {
				t.Fatal("pending preparation granted create or start")
			}
			if err := model.ObserveCreatedCustody(processIdentity(0x41)); !errors.Is(err, ErrState) {
				t.Fatalf("custody accepted before preparation: %v", err)
			}
			if err := model.SettleWrite(pending, outcome); err != nil {
				t.Fatal(err)
			}
			snapshot := model.Snapshot()
			if snapshot.Prepared || snapshot.Custody != CustodyNone || model.Decision().MayCreate {
				t.Fatal("failed or indeterminate preparation widened authority")
			}
			if _, err := model.BeginPreparation(operationID(0x32)); !errors.Is(err, ErrState) {
				t.Fatalf("preparation retry invented: %v", err)
			}
		})
	}
}

func TestConfirmedPreparationAllowsOneExactCurrentLifetimeCustody(t *testing.T) {
	model := preparedModel(t)
	if !model.Decision().MayCreate {
		t.Fatal("confirmed preparation did not permit one custody observation")
	}
	identity := processIdentity(0x41)
	if err := model.ObserveCreatedCustody(identity); err != nil {
		t.Fatal(err)
	}
	if model.Snapshot().ProcessIdentity != identity || model.Decision().MayCreate {
		t.Fatal("custody identity was not retained exactly once")
	}
	if err := model.ObserveCreatedCustody(processIdentity(0x42)); !errors.Is(err, ErrState) {
		t.Fatalf("second custody identity accepted: %v", err)
	}
}

func TestMutationNaturalAbsenceCannotReopenCreation(t *testing.T) {
	model := custodyModel(t)
	if err := model.ObserveAbsence(80); err != nil {
		t.Fatal(err)
	}
	if model.Decision().MayCreate {
		t.Fatal("natural absence reopened creation authority for the consumed attempt")
	}
	if err := model.ObserveCreatedCustody(processIdentity(0x42)); !errors.Is(err, ErrState) {
		t.Fatalf("replacement custody accepted after natural absence: %v", err)
	}
}

func TestCancellationWhilePreparationPendingPreventsCreateAfterLateCommit(t *testing.T) {
	model := newModel(t)
	pending, err := model.BeginPreparation(operationID(0x30))
	if err != nil {
		t.Fatal(err)
	}
	if err := model.LatchStop(TriggerCancellation, 90); err != nil {
		t.Fatal(err)
	}
	if err := model.SettleWrite(pending, WriteConfirmed); err != nil {
		t.Fatal(err)
	}
	if !model.Snapshot().Prepared || model.Decision().MayCreate {
		t.Fatal("late preparation success ignored pre-create cancellation")
	}
	if err := model.ObserveCreatedCustody(processIdentity(0x41)); !errors.Is(err, ErrState) {
		t.Fatalf("create proceeded after pre-create cancellation: %v", err)
	}
}

func TestMutationStorageWaitCannotBlockStopOrFirstSignal(t *testing.T) {
	model := custodyModel(t)
	pending, err := model.BeginRunnerIdentityWrite(operationID(0x32))
	if err != nil {
		t.Fatal(err)
	}
	if err := model.LatchStop(TriggerCancellation, 100); err != nil {
		t.Fatal(err)
	}
	decision := model.Decision()
	if !decision.MaySignal || !decision.StopIndependentOfStorage || !model.Snapshot().Pending.Active {
		t.Fatal("pending storage restored a dependency to the stop path")
	}
	if err := model.RequestSignal(101); err != nil {
		t.Fatalf("first signal request waited for storage: %v", err)
	}
	if err := model.SettleWrite(pending, WriteConfirmed); err != nil {
		t.Fatal(err)
	}
	if !model.Snapshot().RunnerIdentityConfirmed || !model.Snapshot().Stop.Latched {
		t.Fatal("late runner settlement erased identity or stop")
	}
}

func TestMutationLateRunnerSuccessCannotReleaseStart(t *testing.T) {
	model := custodyModel(t)
	pending, err := model.BeginRunnerIdentityWrite(operationID(0x32))
	if err != nil {
		t.Fatal(err)
	}
	if err := model.LatchStop(TriggerFatalFault, 200); err != nil {
		t.Fatal(err)
	}
	if err := model.SettleWrite(pending, WriteConfirmed); err != nil {
		t.Fatal(err)
	}
	if err := model.AttemptStart(201); !errors.Is(err, ErrState) {
		t.Fatalf("late storage success released start after stop: %v", err)
	}
	if model.Snapshot().Start.Attempted || model.Decision().MayStart {
		t.Fatal("late start became visible")
	}
}

func TestStopLatchPreventsNewRunnerPublication(t *testing.T) {
	model := custodyModel(t)
	if err := model.LatchStop(TriggerCancellation, 210); err != nil {
		t.Fatal(err)
	}
	if _, err := model.BeginRunnerIdentityWrite(operationID(0x32)); !errors.Is(err, ErrState) {
		t.Fatalf("runner publication began after stop: %v", err)
	}
	if !model.Decision().MaySignal {
		t.Fatal("runner-publication refusal removed cleanup authority")
	}
}

func TestRunnerPublicationFailureRetainsCleanupButWithholdsStart(t *testing.T) {
	for _, outcome := range []WriteOutcome{WriteFailed, WriteIndeterminate} {
		t.Run(string(outcome), func(t *testing.T) {
			model := custodyModel(t)
			pending, err := model.BeginRunnerIdentityWrite(operationID(0x32))
			if err != nil {
				t.Fatal(err)
			}
			if err := model.SettleWrite(pending, outcome); err != nil {
				t.Fatal(err)
			}
			if err := model.LatchStop(TriggerFatalFault, 220); err != nil {
				t.Fatal(err)
			}
			decision := model.Decision()
			if decision.MayStart || !decision.MaySignal {
				t.Fatal("runner publication failure granted start or removed cleanup")
			}
		})
	}
}

func TestMutationUncertainSignalCannotBeRepeated(t *testing.T) {
	model := runnerModel(t)
	if err := model.LatchStop(TriggerCancellation, 300); err != nil {
		t.Fatal(err)
	}
	if err := model.RequestSignal(301); err != nil {
		t.Fatal(err)
	}
	if err := model.MarkSignalUncertain(); err != nil {
		t.Fatal(err)
	}
	if err := model.RequestSignal(302); !errors.Is(err, ErrState) {
		t.Fatalf("uncertain signal was redriven: %v", err)
	}
	if got := model.Snapshot().Signal.Attempts; got != 1 {
		t.Fatalf("signal attempts = %d, want 1", got)
	}
	if err := model.ObserveAbsence(350); err != nil {
		t.Fatalf("uncertain response blocked exact absence reconciliation: %v", err)
	}
}

func TestMutationRestartRetainsPreparationButCannotAdoptPID(t *testing.T) {
	model := custodyModel(t)
	model.Restart()
	snapshot := model.Snapshot()
	if !snapshot.Prepared || !snapshot.CreationConsumed {
		t.Fatal("restart forgot confirmed preparation or consumed creation")
	}
	if !snapshot.RecoveryRequired || snapshot.Custody != CustodyNone || snapshot.ProcessIdentity != (ProcessIdentity{}) {
		t.Fatal("restart restored live custody or omitted recovery requirement")
	}
	if err := model.ObserveCreatedCustody(processIdentity(0x41)); !errors.Is(err, ErrRecovery) {
		t.Fatalf("restart adopted a PID/process identity: %v", err)
	}
	if err := model.RequestSignal(400); !errors.Is(err, ErrRecovery) {
		t.Fatalf("restart manufactured signal authority: %v", err)
	}
}

func TestRestartMakesOutstandingWriteIndeterminate(t *testing.T) {
	model := custodyModel(t)
	if _, err := model.BeginRunnerIdentityWrite(operationID(0x32)); err != nil {
		t.Fatal(err)
	}
	model.Restart()
	snapshot := model.Snapshot()
	if snapshot.Pending.Active || snapshot.LastWriteOutcome != WriteIndeterminate || !snapshot.RecoveryRequired {
		t.Fatal("restart failed to retain pending-write uncertainty")
	}
	if snapshot.Custody != CustodyNone || snapshot.RunnerIdentityConfirmed {
		t.Fatal("restart invented custody or runner publication")
	}
}

func TestMutationRepeatedTriggersRetainEarliestDeadline(t *testing.T) {
	model := runnerModel(t)
	if err := model.AttemptStart(100); err != nil {
		t.Fatal(err)
	}
	if err := model.LatchStop(TriggerWallDeadline, 1100); err != nil {
		t.Fatal(err)
	}
	if err := model.LatchStop(TriggerCancellation, 300); err != nil {
		t.Fatal(err)
	}
	if err := model.LatchStop(TriggerFatalFault, 450); err != nil {
		t.Fatal(err)
	}
	if err := model.LatchStop(TriggerFatalFault, 250); err != nil {
		t.Fatal(err)
	}
	stop := model.Snapshot().Stop
	if stop.Trigger != TriggerFatalFault || stop.ActionTick != 250 {
		t.Fatalf("earliest action reset: %+v", stop)
	}
	if err := model.RequestSignal(260); err != nil {
		t.Fatal(err)
	}
	if err := model.ObserveAbsence(1201); err != nil {
		t.Fatal(err)
	}
	if got := model.Snapshot().Timing; got != TimingSatisfied {
		t.Fatalf("timing = %s, want satisfied", got)
	}
}

func TestWallDeadlineMustUseOriginalStartAnchor(t *testing.T) {
	model := runnerModel(t)
	if err := model.AttemptStart(100); err != nil {
		t.Fatal(err)
	}
	if err := model.LatchStop(TriggerWallDeadline, 1101); !errors.Is(err, ErrClock) {
		t.Fatalf("shifted wall deadline accepted: %v", err)
	}
	if err := model.LatchStop(TriggerWallDeadline, 1100); err != nil {
		t.Fatal(err)
	}
}

func TestTimingViolationCannotBecomeSuccess(t *testing.T) {
	model := runnerModel(t)
	if err := model.AttemptStart(100); err != nil {
		t.Fatal(err)
	}
	if err := model.LatchStop(TriggerCancellation, 200); err != nil {
		t.Fatal(err)
	}
	if err := model.RequestSignal(300); err != nil {
		t.Fatal(err)
	}
	if err := model.ObserveAbsence(1501); err != nil {
		t.Fatal(err)
	}
	if got := model.Snapshot().Timing; got != TimingViolated {
		t.Fatalf("timing = %s, want violated", got)
	}
	if err := model.LatchStop(TriggerCancellation, 1400); err != nil {
		t.Fatal(err)
	}
	if got := model.Snapshot().Timing; got != TimingViolated {
		t.Fatalf("late trigger rewrote failed timing to %s", got)
	}
}

func TestInvalidAndOverflowingClockInputsLeaveStateUnchanged(t *testing.T) {
	model := runnerModel(t)
	before := model.Snapshot()
	if err := model.AttemptStart(maxSafeTick); !errors.Is(err, ErrClock) {
		t.Fatalf("overflowing start anchor accepted: %v", err)
	}
	if got := model.Snapshot(); got != before {
		t.Fatal("failed start mutated state")
	}
	if err := model.LatchStop("unknown", 1); !errors.Is(err, ErrClock) {
		t.Fatalf("unknown trigger accepted: %v", err)
	}
	if got := model.Snapshot(); got != before {
		t.Fatal("failed trigger mutated state")
	}
}

func TestMutationTerminalCannotBeginOrCommitBeforeAbsence(t *testing.T) {
	model := runnerModel(t)
	if _, err := model.BeginTerminalWrite(operationID(0x33)); !errors.Is(err, ErrState) {
		t.Fatalf("terminal write began before absence: %v", err)
	}
	if err := model.LatchStop(TriggerCancellation, 500); err != nil {
		t.Fatal(err)
	}
	if err := model.RequestSignal(501); err != nil {
		t.Fatal(err)
	}
	if _, err := model.BeginTerminalWrite(operationID(0x33)); !errors.Is(err, ErrState) {
		t.Fatalf("signal return substituted for absence: %v", err)
	}
	if err := model.ObserveAbsence(550); err != nil {
		t.Fatal(err)
	}
	pending, err := model.BeginTerminalWrite(operationID(0x33))
	if err != nil {
		t.Fatal(err)
	}
	if err := model.SettleWrite(pending, WriteFailed); err != nil {
		t.Fatal(err)
	}
	snapshot := model.Snapshot()
	if snapshot.TerminalConfirmed || snapshot.OutputReleased || snapshot.CapacityReleased {
		t.Fatal("failed terminal publication released completion")
	}
}

func TestConfirmedTerminalJoinIsTheOnlyReleasePoint(t *testing.T) {
	model := runnerModel(t)
	if err := model.ObserveAbsence(600); err != nil {
		t.Fatal(err)
	}
	pending, err := model.BeginTerminalWrite(operationID(0x33))
	if err != nil {
		t.Fatal(err)
	}
	if model.Decision().MayRelease {
		t.Fatal("pending terminal publication released state")
	}
	if err := model.SettleWrite(pending, WriteConfirmed); err != nil {
		t.Fatal(err)
	}
	snapshot := model.Snapshot()
	if !snapshot.TerminalConfirmed || !snapshot.OutputReleased || !snapshot.CapacityReleased || !model.Decision().MayRelease {
		t.Fatal("confirmed terminal join did not release all three projections")
	}
}

func TestRestartAfterConfirmedTerminalKeepsOnlyDurableRelease(t *testing.T) {
	model := runnerModel(t)
	if err := model.ObserveAbsence(600); err != nil {
		t.Fatal(err)
	}
	pending, err := model.BeginTerminalWrite(operationID(0x33))
	if err != nil {
		t.Fatal(err)
	}
	if err := model.SettleWrite(pending, WriteConfirmed); err != nil {
		t.Fatal(err)
	}
	if err := model.LatchStop(TriggerCancellation, 700); !errors.Is(err, ErrState) {
		t.Fatalf("terminal record accepted a new stop trigger: %v", err)
	}
	model.Restart()
	snapshot := model.Snapshot()
	if snapshot.RecoveryRequired || !snapshot.TerminalConfirmed ||
		!snapshot.OutputReleased || !snapshot.CapacityReleased {
		t.Fatal("restart lost confirmed durable terminal release")
	}
	if snapshot.Custody != CustodyNone || snapshot.ProcessIdentity != (ProcessIdentity{}) ||
		snapshot.Stop.Latched || snapshot.Signal.Requested || snapshot.Absence.Observed ||
		snapshot.Start.Attempted {
		t.Fatal("restart retained current-lifetime control state")
	}
}

func TestOnlyOneExactStorageOperationMayBeOutstanding(t *testing.T) {
	model := custodyModel(t)
	pending, err := model.BeginRunnerIdentityWrite(operationID(0x32))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := model.BeginTerminalWrite(operationID(0x33)); !errors.Is(err, ErrPendingWrite) {
		t.Fatalf("second storage operation bypassed pending worker: %v", err)
	}
	mutated := pending
	mutated.Generation++
	if err := model.SettleWrite(mutated, WriteConfirmed); !errors.Is(err, ErrPendingWrite) {
		t.Fatalf("wrong generation settled pending write: %v", err)
	}
	mutated = pending
	mutated.OperationID = operationID(0x34)
	if err := model.SettleWrite(mutated, WriteConfirmed); !errors.Is(err, ErrPendingWrite) {
		t.Fatalf("wrong operation settled pending write: %v", err)
	}
	if !model.Snapshot().Pending.Active {
		t.Fatal("mismatched settlement consumed pending write")
	}
}

func TestDecisionMatchesIndependentSnapshotOracle(t *testing.T) {
	model := newModel(t)
	assertDecision(t, model)
	pending, err := model.BeginPreparation(operationID(0x30))
	if err != nil {
		t.Fatal(err)
	}
	assertDecision(t, model)
	if err := model.SettleWrite(pending, WriteConfirmed); err != nil {
		t.Fatal(err)
	}
	assertDecision(t, model)
	if err := model.ObserveCreatedCustody(processIdentity(0x41)); err != nil {
		t.Fatal(err)
	}
	_, err = model.BeginRunnerIdentityWrite(operationID(0x32))
	if err != nil {
		t.Fatal(err)
	}
	assertDecision(t, model)
	if err := model.LatchStop(TriggerCancellation, 700); err != nil {
		t.Fatal(err)
	}
	assertDecision(t, model)
	if err := model.RequestSignal(701); err != nil {
		t.Fatal(err)
	}
	assertDecision(t, model)
	if err := model.ObserveAbsence(750); err != nil {
		t.Fatal(err)
	}
	assertDecision(t, model)
}

func TestBindingsAndClockPolicyAreClosed(t *testing.T) {
	bindings := validBindings()
	model, err := NewModel(bindings)
	if err != nil {
		t.Fatal(err)
	}
	if model.Snapshot().Contract != ContractIdentity || model.Snapshot().RecordVersion != RecordVersionV1 {
		t.Fatal("model identity mismatch")
	}
	if got := model.ClockPolicy(); got != (ClockPolicy{WallAfterStartMS: 1000, ForceAbsenceMS: 1000, ActionToAbsenceMS: 1200}) {
		t.Fatalf("clock policy = %+v", got)
	}
	bindings.PlanDigest[0] ^= 0xff
	if model.Bindings().PlanDigest[0] != 0x16 {
		t.Fatal("model retained caller-owned binding bytes")
	}
	zero := validBindings()
	zero.AttemptID = AttemptID{}
	if _, err := NewModel(zero); !errors.Is(err, ErrBinding) {
		t.Fatalf("zero binding accepted: %v", err)
	}
	substituted := validBindings()
	substituted.ApprovalID = ApprovalID(substituted.AttemptID)
	if _, err := NewModel(substituted); !errors.Is(err, ErrBinding) {
		t.Fatalf("same-width identifier substitution accepted: %v", err)
	}
	digestSubstitution := validBindings()
	digestSubstitution.ProfileDigest = ProfileDigest(digestSubstitution.PlanDigest)
	if _, err := NewModel(digestSubstitution); !errors.Is(err, ErrBinding) {
		t.Fatalf("same-width digest substitution accepted: %v", err)
	}
}

func TestNoProductConsumerImportsPassiveTeardownModel(t *testing.T) {
	root := os.DirFS(filepath.Join("..", "..", ".."))
	needle := "internal/execution/teardownpassive"
	err := fs.WalkDir(root, ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" || entry.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" || strings.Contains(path, "teardownpassive") {
			return nil
		}
		contents, readErr := fs.ReadFile(root, path)
		if readErr != nil {
			return readErr
		}
		if bytes.Contains(contents, []byte(needle)) {
			t.Errorf("product consumer imports passive teardown model: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func assertDecision(t *testing.T, model *Model) {
	t.Helper()
	snapshot := model.Snapshot()
	if got, want := model.Decision(), independentDecision(snapshot); got != want {
		t.Fatalf("decision = %+v, independent = %+v", got, want)
	}
}

func independentDecision(snapshot Snapshot) Decision {
	blocked := snapshot.RecoveryRequired || snapshot.TerminalConfirmed
	return Decision{
		MayCreate: snapshot.Prepared && !snapshot.CreationConsumed && !blocked &&
			snapshot.Custody == CustodyNone &&
			!snapshot.Stop.Latched && !snapshot.Pending.Active,
		MayStart: snapshot.Prepared && !blocked && snapshot.Custody == CustodyExact &&
			snapshot.RunnerIdentityConfirmed && !snapshot.Stop.Latched &&
			!snapshot.Start.Attempted && !snapshot.Pending.Active,
		MaySignal: !blocked && snapshot.Custody == CustodyExact && snapshot.Stop.Latched &&
			!snapshot.Signal.Requested && !snapshot.Absence.Observed,
		MayPublishTerminal: !blocked && snapshot.Absence.Observed &&
			!snapshot.TerminalWriteStarted && !snapshot.Pending.Active,
		MayRelease:               snapshot.TerminalConfirmed,
		StopIndependentOfStorage: !blocked && snapshot.Custody == CustodyExact && snapshot.Stop.Latched,
		RecoveryRequired:         snapshot.RecoveryRequired,
	}
}

func newModel(t *testing.T) *Model {
	t.Helper()
	model, err := NewModel(validBindings())
	if err != nil {
		t.Fatal(err)
	}
	return model
}

func preparedModel(t *testing.T) *Model {
	t.Helper()
	model := newModel(t)
	pending, err := model.BeginPreparation(operationID(0x30))
	if err != nil {
		t.Fatal(err)
	}
	if err := model.SettleWrite(pending, WriteConfirmed); err != nil {
		t.Fatal(err)
	}
	return model
}

func custodyModel(t *testing.T) *Model {
	t.Helper()
	model := preparedModel(t)
	if err := model.ObserveCreatedCustody(processIdentity(0x41)); err != nil {
		t.Fatal(err)
	}
	return model
}

func runnerModel(t *testing.T) *Model {
	t.Helper()
	model := custodyModel(t)
	pending, err := model.BeginRunnerIdentityWrite(operationID(0x32))
	if err != nil {
		t.Fatal(err)
	}
	if err := model.SettleWrite(pending, WriteConfirmed); err != nil {
		t.Fatal(err)
	}
	return model
}

func validBindings() Bindings {
	return Bindings{
		InstallationID:           installationID(0x11),
		TrustEpochDigest:         trustEpochDigest(0x12),
		SupervisorIdentityDigest: supervisorIdentityDigest(0x13),
		AttemptID:                attemptID(0x14),
		ApprovalID:               approvalID(0x15),
		RegistrationID:           registrationID(0x16),
		PlanDigest:               planDigest(0x16),
		ProfileDigest:            profileDigest(0x17),
		RunnerDigest:             runnerDigest(0x18),
		PreparationOperationID:   operationID(0x30),
	}
}

func installationID(value byte) (result InstallationID)     { fill(result[:], value); return result }
func trustEpochDigest(value byte) (result TrustEpochDigest) { fill(result[:], value); return result }
func supervisorIdentityDigest(value byte) (result SupervisorIdentityDigest) {
	fill(result[:], value)
	return result
}
func attemptID(value byte) (result AttemptID)             { fill(result[:], value); return result }
func approvalID(value byte) (result ApprovalID)           { fill(result[:], value); return result }
func registrationID(value byte) (result RegistrationID)   { fill(result[:], value); return result }
func planDigest(value byte) (result PlanDigest)           { fill(result[:], value); return result }
func profileDigest(value byte) (result ProfileDigest)     { fill(result[:], value); return result }
func runnerDigest(value byte) (result RunnerDigest)       { fill(result[:], value); return result }
func operationID(value byte) (result OperationID)         { fill(result[:], value); return result }
func processIdentity(value byte) (result ProcessIdentity) { fill(result[:], value); return result }

func fill(target []byte, value byte) {
	for index := range target {
		target[index] = value
	}
}
