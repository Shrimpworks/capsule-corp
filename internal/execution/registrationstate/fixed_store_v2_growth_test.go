package registrationstate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

func TestFixedStoreV2SecondSegmentKnownAnswerLookupReplayAndRecovery(t *testing.T) {
	path, store, owner, keys := newEligibleGrowthStoreV2(t, 2)
	first := activateNextArchiveSegment(t, store, owner)
	firstPath := archiveSegmentPath(path, first.segments[0].Segment.Digest())
	firstBytes := mustReadFile(t, firstPath)
	firstCheckpoint := first.active.View().CurrentCheckpoint

	plan, err := first.PlanArchive(context.Background(), owner, archiveOneCohortLimits(t))
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := first.PrepareArchive(context.Background(), owner, plan)
	if err != nil {
		t.Fatal(err)
	}
	segmentView := prepared.candidate.Descriptors[1]
	if segmentView.Ordinal != 2 || segmentView.ArchiveGeneration != 3 ||
		prepared.plan.View().CheckpointHead != firstCheckpoint.Digest {
		t.Fatalf("second segment coordinates = descriptor %#v plan %#v", segmentView, prepared.plan.View())
	}
	verified, err := first.VerifyPreparedArchive(context.Background(), owner, prepared)
	if err != nil {
		t.Fatal(err)
	}
	second, err := first.ActivateArchive(context.Background(), owner, verified, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(firstBytes, mustReadFile(t, firstPath)) {
		t.Fatal("second activation rewrote first immutable segment")
	}
	report, err := VerifyFixedFileStoreV2(path)
	if err != nil {
		t.Fatal(err)
	}
	if report.SnapshotGeneration != 3 || report.ArchiveGeneration != 3 || report.SegmentCount != 2 ||
		report.HotCounts != (archivestate.ArchiveCounts{}) || report.ArchivedCounts.Registrations != 2 ||
		report.ArchivedCounts.Approvals != 2 || report.ArchivedCounts.Attempts != 2 ||
		report.ArchivedCounts.Lifecycles != 2 || report.ArchivedCounts.Effects != 2 ||
		report.ArchivedCounts.Instances != 2 {
		t.Fatalf("second activation report = %#v", report)
	}
	if second.active.View().PreviousCheckpoint != firstCheckpoint ||
		second.segments[1].Segment.View().PriorCheckpointDigest != firstCheckpoint.Digest {
		t.Fatal("second activation lost exact checkpoint ancestry")
	}
	for _, key := range keys {
		registration, resolveErr := second.ResolveRegistration(context.Background(), key.registrationID)
		if resolveErr != nil || registration.RegistrationID != key.registrationID {
			t.Fatalf("archived registration %x = %#v, %v", key.registrationID, registration, resolveErr)
		}
		approval, state, replayErr := second.ResolveApprovalReplay(
			context.Background(), key.payloadDigest, key.exactPayload, key.authorizationIdentity,
		)
		if replayErr != nil || approval.ApprovalID() != key.approvalID || state != approvalattempt.ApprovalConsumed {
			t.Fatalf("archived approval replay %x = %#v/%s, %v", key.approvalID, approval, state, replayErr)
		}
		attempt, attemptState, replayErr := second.ResolveAttemptReplay(context.Background(), key.registrationID, key.approvalID)
		if replayErr != nil || attempt.AttemptID() != key.attemptID || attemptState != approvalattempt.AttemptCreated {
			t.Fatalf("archived attempt replay %x = %#v/%s, %v", key.attemptID, attempt, attemptState, replayErr)
		}
	}
	recovery, err := second.RecoveryAttemptIDs(context.Background())
	if err != nil || len(recovery) != 0 {
		t.Fatalf("second-segment archived recovery set = %x, %v", recovery, err)
	}

	activeDigest := sha256.Sum256(mustReadFile(t, path))
	secondFileDigest := sha256.Sum256(mustReadFile(t, archiveSegmentPath(path, second.segments[1].Segment.Digest())))
	secondSemanticDigest := second.segments[1].Segment.Digest()
	combinedIndexDigest := second.active.View().CombinedIndexDigest
	known := []struct {
		name string
		got  []byte
		want string
	}{
		{"active-file", activeDigest[:], "0c03403504686204007b26519f71e2bbfe7365acf295cfd0c93e055050e592cf"},
		{"second-segment", secondFileDigest[:], "5638f9a3e10374cf9b6c3cece471dc8cdedd3d074b8a9d38c7a9d1304dbc3d82"},
		{"second-semantic", secondSemanticDigest[:], "db57a840b1e5ea75b0ffe4700664300b0108e818620518be6f8811c4eb3476e8"},
		{"second-checkpoint", report.CurrentCheckpoint.Digest[:], "5e310a41ea2fc17fe43d211cdab6d44c1377700a1992a37c30c94dc4d143b9e8"},
		{"combined-index", combinedIndexDigest[:], "b7d4f6a92313c2c8b629d06c04fa9be03ec891e4ca5157cf1409dc3360396d3c"},
	}
	for _, answer := range known {
		if got := hex.EncodeToString(answer.got); got != answer.want {
			t.Errorf("F4C known answer %s=%s, want %s", answer.name, got, answer.want)
		}
	}
}

func TestFixedStoreV2SecondActivationFaultRaceAndProcessDeathOldOrCompleteNew(t *testing.T) {
	for _, test := range []struct {
		point         ArchiveFault
		indeterminate bool
		newWorld      bool
		orphans       uint16
	}{
		{FaultArchiveAfterSegmentTempCreate, false, false, 0},
		{FaultArchiveAfterSegmentTempWrite, false, false, 0},
		{FaultArchiveAfterSegmentTempSync, false, false, 0},
		{FaultArchiveBeforeSegmentPublish, false, false, 0},
		{FaultArchiveAfterSegmentPublish, true, false, 1},
		{FaultArchiveAfterSegmentDirSync, true, false, 1},
		{FaultArchiveAfterActiveTempCreate, false, false, 1},
		{FaultArchiveAfterActiveTempWrite, false, false, 1},
		{FaultArchiveAfterActiveTempSync, false, false, 1},
		{FaultArchiveBeforeActivation, false, false, 1},
		{FaultArchiveAfterActivation, true, true, 0},
		{FaultArchiveAfterActiveDirSync, true, true, 0},
		{FaultArchiveAfterReopen, true, true, 0},
	} {
		t.Run(string(test.point), func(t *testing.T) {
			path, store, owner, _ := newEligibleGrowthStoreV2(t, 2)
			store = activateNextArchiveSegment(t, store, owner)
			before := inventoryDigest(t, path)
			verified := mustPreparedArchive(t, store, owner)
			_, activateErr := store.ActivateArchive(context.Background(), owner, verified, archiveFaultStub{point: test.point})
			if activateErr == nil || errors.Is(activateErr, ErrArchiveOutcomeIndeterminate) != test.indeterminate {
				t.Fatalf("second activation fault %s = %v", test.point, activateErr)
			}
			reopened, openErr := OpenFixedFileStoreV2(path)
			if openErr != nil {
				t.Fatal(openErr)
			}
			report, verifyErr := VerifyFixedFileStoreV2(path)
			if verifyErr != nil {
				t.Fatal(verifyErr)
			}
			if test.newWorld {
				if report.SnapshotGeneration != 3 || report.SegmentCount != 2 || report.HotCounts.Registrations != 0 {
					t.Fatalf("second fault successor = %#v", report)
				}
			} else if report.SnapshotGeneration != 2 || report.SegmentCount != 1 || report.HotCounts.Registrations != 1 ||
				inventoryDigest(t, path).active != before.active {
				t.Fatalf("second fault predecessor = %#v", report)
			}
			if report.OrphanSegmentCount != test.orphans || reopened.orphans != test.orphans {
				t.Fatalf("second fault orphan count = %d/%d, want %d", report.OrphanSegmentCount, reopened.orphans, test.orphans)
			}
		})
	}

	t.Run("process-death", func(t *testing.T) {
		for _, test := range []struct {
			point       ArchiveFault
			newWorld    bool
			orphanCount uint16
		}{
			{FaultArchiveBeforeSegmentPublish, false, 1},
			{FaultArchiveAfterSegmentPublish, false, 1},
			{FaultArchiveBeforeActivation, false, 1},
			{FaultArchiveAfterActivation, true, 0},
		} {
			t.Run(string(test.point), func(t *testing.T) {
				path, store, owner, _ := newEligibleGrowthStoreV2(t, 2)
				_ = activateNextArchiveSegment(t, store, owner)
				command := exec.Command(os.Args[0], "-test.run=^TestArchiveActivationProcessDeathHelper$") //nolint:gosec // G204: current bounded test binary only.
				command.Env = append(os.Environ(), "CAPSULE_F3_DEATH_PATH="+path, "CAPSULE_F3_DEATH_POINT="+string(test.point))
				err := command.Run()
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) || exitErr.ExitCode() != 77 {
					t.Fatalf("second archive death helper = %v", err)
				}
				report, verifyErr := VerifyFixedFileStoreV2(path)
				if verifyErr != nil {
					t.Fatalf("second reopen after process death: %v", verifyErr)
				}
				if test.newWorld {
					if report.SnapshotGeneration != 3 || report.SegmentCount != 2 || report.HotCounts.Registrations != 0 {
						t.Fatalf("second death successor = %#v", report)
					}
				} else if report.SnapshotGeneration != 2 || report.SegmentCount != 1 || report.HotCounts.Registrations != 1 {
					t.Fatalf("second death predecessor = %#v", report)
				}
				if report.OrphanSegmentCount != test.orphanCount {
					t.Fatalf("second death orphan count = %d, want %d", report.OrphanSegmentCount, test.orphanCount)
				}
			})
		}
	})

	t.Run("concurrent-exact-retries-converge", func(t *testing.T) {
		path, store, owner, _ := newEligibleGrowthStoreV2(t, 2)
		store = activateNextArchiveSegment(t, store, owner)
		verified := mustPreparedArchive(t, store, owner)
		const callers = 12
		results := make(chan error, callers)
		var group sync.WaitGroup
		for range callers {
			group.Add(1)
			go func() {
				defer group.Done()
				_, err := store.ActivateArchive(context.Background(), owner, verified, nil)
				results <- err
			}()
		}
		group.Wait()
		close(results)
		success, stale := 0, 0
		for err := range results {
			switch {
			case err == nil:
				success++
			case errors.Is(err, ErrArchiveStaleTransaction):
				stale++
			default:
				t.Fatalf("competing second activation = %v", err)
			}
		}
		if success != 1 || stale != callers-1 {
			t.Fatalf("second activation success/stale = %d/%d", success, stale)
		}
		if report, err := VerifyFixedFileStoreV2(path); err != nil || report.SegmentCount != 2 {
			t.Fatalf("converged second activation = %#v, %v", report, err)
		}
	})
}

func TestFixedStoreV2SecondSegmentMutationRefusalsNeverRewriteOrScanFallback(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*diskEnvelopeV2)
	}{
		{name: "partial-count", mutate: func(envelope *diskEnvelopeV2) { envelope.ArchivedCounts.Registrations-- }},
		{name: "partial-index", mutate: func(envelope *diskEnvelopeV2) {
			envelope.Indexes.Registrations = envelope.Indexes.Registrations[:len(envelope.Indexes.Registrations)-1]
		}},
		{name: "duplicate-identity", mutate: func(envelope *diskEnvelopeV2) {
			envelope.Indexes.Registrations[1].RegistrationID = envelope.Indexes.Registrations[0].RegistrationID
		}},
		{name: "duplicate-location-cross-segment", mutate: func(envelope *diskEnvelopeV2) {
			envelope.Indexes.Registrations[1].Location = envelope.Indexes.Registrations[0].Location
		}},
		{name: "stale-head", mutate: func(envelope *diskEnvelopeV2) { envelope.CurrentCheckpoint.Digest[0] ^= 0xff }},
	} {
		t.Run(test.name, func(t *testing.T) {
			path, store, owner, _ := newEligibleGrowthStoreV2(t, 2)
			store = activateNextArchiveSegment(t, store, owner)
			_ = activateNextArchiveSegment(t, store, owner)
			envelope := readV2Envelope(t, path)
			test.mutate(&envelope)
			writeV2Envelope(t, path, envelope)
			assertV2OpenDoesNotRewrite(t, path)
		})
	}

	for _, test := range []struct {
		name   string
		mutate func(t *testing.T, firstPath, secondPath string)
	}{
		{name: "missing-second", mutate: func(t *testing.T, _, secondPath string) {
			if err := os.Remove(secondPath); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "corrupt-second", mutate: func(t *testing.T, _, secondPath string) {
			data := mustReadFile(t, secondPath)
			data[len(data)/2] ^= 0xff
			replaceArchiveFixtureFile(t, secondPath, data)
		}},
		{name: "substituted-second", mutate: func(t *testing.T, firstPath, secondPath string) {
			replaceArchiveFixtureFile(t, secondPath, mustReadFile(t, firstPath))
		}},
		{name: "wrong-predecessor", mutate: func(t *testing.T, _, secondPath string) {
			data := mutateArchiveSegmentDisk(func(disk *archiveSegmentDiskV0) {
				disk.PriorCheckpointDigest[0] ^= 0xff
			})(t, mustReadFile(t, secondPath))
			replaceArchiveFixtureFile(t, secondPath, data)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			path, store, owner, _ := newEligibleGrowthStoreV2(t, 2)
			store = activateNextArchiveSegment(t, store, owner)
			firstPath := archiveSegmentPath(path, store.segments[0].Segment.Digest())
			store = activateNextArchiveSegment(t, store, owner)
			secondPath := archiveSegmentPath(path, store.segments[1].Segment.Digest())
			activeBefore := mustReadFile(t, path)
			firstBefore := mustReadFile(t, firstPath)
			test.mutate(t, firstPath, secondPath)
			if _, err := OpenFixedFileStoreV2(path); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
				t.Fatalf("mutated second segment open = %v", err)
			}
			if !bytes.Equal(activeBefore, mustReadFile(t, path)) || !bytes.Equal(firstBefore, mustReadFile(t, firstPath)) {
				t.Fatal("second-segment refusal rewrote active state or predecessor segment")
			}
		})
	}
}

func TestFixedStoreV2SecondSegmentPlanGenerationHeadOwnerAndCandidateFences(t *testing.T) {
	t.Run("stale-generation-and-head", func(t *testing.T) {
		path, store, owner, _ := newEligibleGrowthStoreV2(t, 2)
		store = activateNextArchiveSegment(t, store, owner)
		plan, err := store.PlanArchive(context.Background(), owner, archiveOneCohortLimits(t))
		if err != nil {
			t.Fatal(err)
		}
		if err := store.persistTimeHighWater(context.Background(), store.state.TimeHighWaterUnixSeconds+1); err != nil {
			t.Fatal(err)
		}
		before := inventoryDigest(t, path)
		if _, err := store.PrepareArchive(context.Background(), owner, plan); !errors.Is(err, ErrArchiveStaleTransaction) {
			t.Fatalf("stale second plan prepare = %v", err)
		}
		if after := inventoryDigest(t, path); after != before {
			t.Fatal("stale second plan rewrote archive inventory")
		}
	})

	t.Run("owner-session", func(t *testing.T) {
		path, store, owner, _ := newEligibleGrowthStoreV2(t, 2)
		store = activateNextArchiveSegment(t, store, owner)
		plan, err := store.PlanArchive(context.Background(), owner, archiveOneCohortLimits(t))
		if err != nil {
			t.Fatal(err)
		}
		before := inventoryDigest(t, path)
		wrongOwner := &archiveOwnerStub{held: true, session: lifecycleOwnerID(t, 0xf4)}
		if _, err := store.PrepareArchive(context.Background(), wrongOwner, plan); !errors.Is(err, ErrArchiveOwnerSessionMismatch) {
			t.Fatalf("second plan wrong owner session = %v", err)
		}
		if after := inventoryDigest(t, path); after != before {
			t.Fatal("wrong owner session rewrote archive inventory")
		}
	})

	t.Run("competing-candidates", func(t *testing.T) {
		path, store, owner, _ := newEligibleGrowthStoreV2(t, 3)
		store = activateNextArchiveSegment(t, store, owner)
		onePlan, err := store.PlanArchive(context.Background(), owner, archiveOneCohortLimits(t))
		if err != nil {
			t.Fatal(err)
		}
		twoLimits, err := archivestate.NewArchiveLimits(2, uint64(archivestate.MaxSupervisorArchiveBytes))
		if err != nil {
			t.Fatal(err)
		}
		twoPlan, err := store.PlanArchive(context.Background(), owner, twoLimits)
		if err != nil {
			t.Fatal(err)
		}
		onePrepared, err := store.PrepareArchive(context.Background(), owner, onePlan)
		if err != nil {
			t.Fatal(err)
		}
		twoPrepared, err := store.PrepareArchive(context.Background(), owner, twoPlan)
		if err != nil {
			t.Fatal(err)
		}
		if onePrepared.SegmentDigest() == twoPrepared.SegmentDigest() {
			t.Fatal("competing second candidates unexpectedly share an identity")
		}
		oneVerified, err := store.VerifyPreparedArchive(context.Background(), owner, onePrepared)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := store.ActivateArchive(context.Background(), owner, oneVerified, nil); err != nil {
			t.Fatal(err)
		}
		before := inventoryDigest(t, path)
		if _, err := store.VerifyPreparedArchive(context.Background(), owner, twoPrepared); !errors.Is(err, ErrArchiveStaleTransaction) {
			t.Fatalf("losing second candidate verify = %v", err)
		}
		if after := inventoryDigest(t, path); after != before {
			t.Fatal("losing second candidate rewrote archive inventory")
		}
	})
}

func TestFixedStoreV2SecondSegmentPreservesExactCurrentAndHistoricalEffects(t *testing.T) {
	path, store, owner, keys := newGrowthStoreV2(t, 2, true)
	store = activateNextArchiveSegment(t, store, owner)
	attemptID := keys[1].attemptID
	record, err := store.ReadLifecycle(context.Background(), attemptID)
	if err != nil {
		t.Fatal(err)
	}
	instance, err := lifecyclestate.NewBackendInstanceIdentity(lifecyclestate.BackendInstanceFake, []byte("f4c-second-instance"))
	if err != nil {
		t.Fatal(err)
	}
	var historical lifecyclestate.EffectID
	for _, operation := range []lifecyclestate.Operation{
		lifecyclestate.OperationPrepare,
		lifecyclestate.OperationCreate,
		lifecyclestate.OperationStart,
		lifecyclestate.OperationObserve,
		lifecyclestate.OperationStop,
		lifecyclestate.OperationDestroy,
	} {
		permit, beginErr := store.BeginEffect(context.Background(), attemptID, record.View().RecordVersion, operation)
		if beginErr != nil {
			t.Fatalf("begin second-segment %s: %v", operation, beginErr)
		}
		if operation == lifecyclestate.OperationPrepare {
			historical = permit.View().EffectID
		}
		resultInstance := lifecyclestate.BackendInstanceIdentity{}
		if operation == lifecyclestate.OperationCreate {
			resultInstance = instance
		}
		result, resultErr := lifecyclestate.NewEffectResult(operation, lifecyclestate.EffectResultApplied, resultInstance)
		if resultErr != nil {
			t.Fatal(resultErr)
		}
		record, err = store.ConfirmEffect(context.Background(), permit, result)
		if err != nil {
			t.Fatalf("confirm second-segment %s: %v", operation, err)
		}
	}
	absent, err := lifecyclestate.NewReconcileResult(
		lifecyclestate.OperationDestroy, lifecyclestate.ReconciliationAuthoritativelyAbsent,
		lifecyclestate.BackendInstanceIdentity{},
	)
	if err != nil {
		t.Fatal(err)
	}
	record, err = store.RecordReconciliation(context.Background(), attemptID, record.View().RecordVersion, absent)
	if err != nil {
		t.Fatal(err)
	}
	currentID := record.View().EffectID
	store = activateNextArchiveSegment(t, store, owner)

	historicalResult, err := store.ResolveEffect(context.Background(), historical)
	if err != nil || historicalResult.Classification != RetainedEffectSupersededByCurrent || historicalResult.LifecyclePresent ||
		historicalResult.Operation != lifecyclestate.OperationPrepare || historicalResult.VisibleV1Seed {
		t.Fatalf("second-segment historical effect = %#v, %v", historicalResult, err)
	}
	currentResult, err := store.ResolveEffect(context.Background(), currentID)
	if err != nil || currentResult.Classification != RetainedEffectCurrent || !currentResult.LifecyclePresent ||
		currentResult.Lifecycle.View().State != lifecyclestate.StateDestroyed {
		t.Fatalf("second-segment current effect = %#v, %v", currentResult, err)
	}
	report, err := VerifyFixedFileStoreV2(path)
	if err != nil || report.SegmentCount != 2 || report.ArchivedCounts.Effects != 7 || report.HotCounts.Effects != 0 {
		t.Fatalf("second-segment effect reconstruction = %#v, %v", report, err)
	}
}

type growthLookupKey struct {
	registrationID        v0candidate.RegistrationID
	approvalID            approvalattempt.ApprovalID
	attemptID             approvalattempt.AttemptID
	payloadDigest         approvalattempt.ApprovalPayloadDigest
	exactPayload          []byte
	authorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity
}

func newEligibleGrowthStoreV2(t *testing.T, count int) (string, *FixedFileStoreV2, *archiveOwnerStub, []growthLookupKey) {
	t.Helper()
	return newGrowthStoreV2(t, count, false)
}

func newGrowthStoreV2(t *testing.T, count int, activeLast bool) (string, *FixedFileStoreV2, *archiveOwnerStub, []growthLookupKey) {
	t.Helper()
	state, template := stateAndLifecycleRecord(t)
	templateRegistration := state.Registrations[0]
	templateApproval := state.Approvals[0]
	templateAttempt := state.Attempts[0]
	templateBindings := template.Bindings().View()
	state.Registrations = make([]registrationEntry, count)
	state.Approvals = make([]approvalattempt.ApprovalRecord, count)
	state.Attempts = make([]approvalattempt.ExecutionAttempt, count)
	state.LastRegistrationSequence = v0candidate.UInt53(count)
	state.TimeHighWaterUnixSeconds = templateRegistration.Index.ExpiresAt
	records := make([]lifecyclestate.Record, count)
	keys := make([]growthLookupKey, count)
	for index := range count {
		ordinal := uint16(index + 1) //nolint:gosec // test count is bounded by retained capacity.
		registration := templateRegistration
		registration.Index.RegistrationID = templateRegistration.Index.RegistrationID
		binary.BigEndian.PutUint16(registration.Index.RegistrationID[14:], ordinal)
		registration.Index.RegistrationSequence = v0candidate.PositiveUInt53(index + 1)
		registrationView := v0candidate.PlanRegistration{
			ObjectType: v0candidate.PlanRegistrationObjectType, ObjectVersion: v0candidate.CandidateObjectVersion,
			RegistrationID: registration.Index.RegistrationID, RegistrationSequence: registration.Index.RegistrationSequence,
			PlanDigest: registration.Record.RecomputedPlanDigest, InstallationID: registration.Index.InstallationID,
			EpochSequence: registration.Index.EpochSequence, EpochDigest: registration.Index.EpochDigest,
			SupervisorID: registration.Index.SupervisorID, ExpiresAt: registration.Index.ExpiresAt,
		}
		registration.Record.WireRegistrationBytes = encodePlanRegistration(registrationView)
		state.Registrations[index] = registration

		approval := approvalattempt.CloneApprovalRecord(templateApproval)
		binary.BigEndian.PutUint16(approval.ApprovalID[14:], ordinal)
		binary.BigEndian.PutUint16(approval.AttemptNonce[14:], ordinal)
		binary.BigEndian.PutUint16(approval.ExactEnvelopeBytes[len(approval.ExactEnvelopeBytes)-2:], ordinal)
		binary.BigEndian.PutUint16(approval.ExactPayloadBytes[len(approval.ExactPayloadBytes)-2:], ordinal)
		approval.EnvelopeDigest = approvalattempt.ApprovalEnvelopeDigest(sha256.Sum256(approval.ExactEnvelopeBytes))
		approval.PayloadDigest = approvalattempt.ApprovalPayloadDigest(sha256.Sum256(approval.ExactPayloadBytes))
		approval.RegistrationID = registration.Index.RegistrationID
		approval.RegistrationSequence = registration.Index.RegistrationSequence
		approval.ConsumedAttemptID = templateAttempt.AttemptID
		binary.BigEndian.PutUint16(approval.ConsumedAttemptID[14:], ordinal)
		state.Approvals[index] = approval

		attempt := templateAttempt
		attempt.AttemptID = approval.ConsumedAttemptID
		attempt.ApprovalID = approval.ApprovalID
		attempt.AttemptNonce = approval.AttemptNonce
		attempt.RegistrationID = registration.Index.RegistrationID
		attempt.RegistrationSequence = registration.Index.RegistrationSequence
		attempt.ApprovalPayloadDigest = approval.PayloadDigest
		state.Attempts[index] = attempt

		bindingsView := templateBindings
		bindingsView.ProfileReviewAttestationDigests = append([]v0candidate.ProfileReviewAttestationDigest(nil), templateBindings.ProfileReviewAttestationDigests...)
		bindingsView.AttemptID, bindingsView.ApprovalID, bindingsView.AttemptNonce = attempt.AttemptID, attempt.ApprovalID, attempt.AttemptNonce
		bindingsView.RegistrationID, bindingsView.RegistrationSequence = attempt.RegistrationID, attempt.RegistrationSequence
		bindingsView.ApprovalPayloadDigest = attempt.ApprovalPayloadDigest
		bindings, err := lifecyclestate.NewImmutableBindings(bindingsView)
		if err != nil {
			t.Fatal(err)
		}
		recordView := template.View()
		recordView.Bindings, recordView.ImmutableBindingDigest = bindings, bindings.Digest()
		base, err := lifecyclestate.NewRecord(recordView, state.TimeHighWaterUnixSeconds)
		if err != nil {
			t.Fatal(err)
		}
		if activeLast && index == count-1 {
			records[index] = lifecycleRecordForV2State(t, state, base, lifecyclestate.StatePreparePending)
			keys[index] = growthLookupKey{
				registrationID: registration.Index.RegistrationID, approvalID: approval.ApprovalID,
				attemptID: attempt.AttemptID, payloadDigest: approval.PayloadDigest,
				exactPayload: bytes.Clone(approval.ExactPayloadBytes), authorizationIdentity: approval.AuthorizationIdentity,
			}
			continue
		}
		destroyConfirmed := lifecycleRecordForV2State(t, state, base, lifecyclestate.StateDestroyConfirmed)
		destroyedView := destroyConfirmed.View()
		destroyedView.State = lifecyclestate.StateDestroyed
		destroyedView.CleanupRequired = false
		destroyedView.LastReconciliation = lifecyclestate.ReconciliationAuthoritativelyAbsent
		destroyedView.AutomaticRecoveryCount = 1
		destroyedView.TerminalAt = lifecyclestate.OptionalUnixSeconds{Present: true, Value: destroyedView.LastTransitionAt}
		records[index], err = lifecyclestate.NewRecord(destroyedView, state.TimeHighWaterUnixSeconds)
		if err != nil {
			t.Fatal(err)
		}
		keys[index] = growthLookupKey{
			registrationID: registration.Index.RegistrationID, approvalID: approval.ApprovalID,
			attemptID: attempt.AttemptID, payloadDigest: approval.PayloadDigest,
			exactPayload: bytes.Clone(approval.ExactPayloadBytes), authorizationIdentity: approval.AuthorizationIdentity,
		}
	}
	if err := recomputeAuthoritySetDigests(&state); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "supervisor-state.json")
	writeV1Envelope(t, path, encodedEnvelopeV1(state, records))
	_ = mustMigrateV2(t, path)
	owner := archiveOwnerForTest(t)
	store, err := OpenFixedFileStoreV2WithOptions(path, FixedFileStoreV1Options{
		EffectIDs: &lifecycleEffectIDSequence{next: 401}, OwnerSessionID: owner.session,
	})
	if err != nil {
		t.Fatal(err)
	}
	return path, store, owner, keys
}

func archiveOneCohortLimits(t *testing.T) archivestate.ArchiveLimits {
	t.Helper()
	limits, err := archivestate.NewArchiveLimits(1, uint64(archivestate.MaxSupervisorArchiveBytes))
	if err != nil {
		t.Fatal(err)
	}
	return limits
}

func activateNextArchiveSegment(t *testing.T, store *FixedFileStoreV2, owner ArchiveOwner) *FixedFileStoreV2 {
	t.Helper()
	plan, err := store.PlanArchive(context.Background(), owner, archiveOneCohortLimits(t))
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := store.PrepareArchive(context.Background(), owner, plan)
	if err != nil {
		t.Fatal(err)
	}
	verified, err := store.VerifyPreparedArchive(context.Background(), owner, prepared)
	if err != nil {
		t.Fatal(err)
	}
	reopened, err := store.ActivateArchive(context.Background(), owner, verified, nil)
	if err != nil {
		t.Fatal(err)
	}
	return reopened
}

// TestFixedStoreV2Exact64SegmentsProductionPipelineAcceptAndSegment65Refuses
// drives all 64 archive-segment activations through the real disk-backed
// Plan/Prepare/Verify/Activate transaction pipeline, unlike
// newExactArchiveCapacityStoreV2 below, which hand-assembles ordinals 1-63 in
// memory to stay inside the package's race-test timeout. Replaying all five
// full-verification phases 64 times is comparatively slow, so this is the
// gated full-path corpus: it only runs outside `go test -short`, keeping one
// authoritative construction path exercised (just not on every invocation)
// per issue #316.
func TestFixedStoreV2Exact64SegmentsProductionPipelineAcceptAndSegment65Refuses(t *testing.T) {
	if testing.Short() {
		t.Skip("full disk-backed 64-segment activation pipeline skipped in -short mode; see newExactArchiveCapacityStoreV2 for the fast in-memory equivalent")
	}
	path, store, owner, _ := newEligibleGrowthStoreV2(t, archivestate.MaxReferencedArchiveSegments+1)
	for ordinal := 1; ordinal <= archivestate.MaxReferencedArchiveSegments; ordinal++ {
		store = activateNextArchiveSegment(t, store, owner)
		if got := len(store.segments); got != ordinal {
			t.Fatalf("segment count after ordinal %d = %d", ordinal, got)
		}
	}
	report, err := VerifyFixedFileStoreV2(path)
	if err != nil || report.SegmentCount != archivestate.MaxReferencedArchiveSegments ||
		report.ArchiveGeneration != archivestate.ArchiveGeneration(archivestate.MaxReferencedArchiveSegments+1) ||
		report.HotCounts.Registrations != 1 {
		t.Fatalf("exact 64-segment world = %#v, %v", report, err)
	}
	before := inventoryDigest(t, path)
	if _, err := store.PlanArchive(context.Background(), owner, archiveOneCohortLimits(t)); err == nil ||
		!bytes.Contains([]byte(err.Error()), []byte("CAPACITY")) {
		t.Fatalf("segment 65 plan = %v", err)
	}
	after := inventoryDigest(t, path)
	if before != after {
		t.Fatal("segment 65 refusal rewrote, deleted, or evicted retained history")
	}
}

// newExactArchiveCapacityStoreV2 constructs the closed 64-segment input shared
// by the F4C capacity and F5 backup assertions. The ordinary activation, fault,
// response-loss, and reopen paths are covered by the focused F3/F4C tests.
// Replaying those five full-verification phases for every ordinal made this one
// boundary fixture consume the package's entire race-test timeout, so it stays
// gated as TestFixedStoreV2Exact64SegmentsProductionPipelineAcceptAndSegment65Refuses
// above instead of running here on every invocation. This helper instead
// builds deterministic successors 1 through 63 in memory, publishes their
// sealed segments before the active bytes, and performs one complete
// closed-world load. It then uses the production Plan/Prepare/Verify/Activate
// path for ordinal 64 so the load-bearing 63-to-64 boundary remains covered.
func newExactArchiveCapacityStoreV2(t *testing.T) (string, *FixedFileStoreV2, *archiveOwnerStub) {
	t.Helper()
	path, _, owner, _ := newEligibleGrowthStoreV2(t, archivestate.MaxReferencedArchiveSegments+1)
	root, err := ensureArchiveRoot(path)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := loadV2State(path)
	if err != nil {
		t.Fatal(err)
	}
	for ordinal := 1; ordinal < archivestate.MaxReferencedArchiveSegments; ordinal++ {
		plan, planErr := planArchiveFromLoaded(loaded, owner.OwnerSessionID(), archiveOneCohortLimits(t))
		if planErr != nil {
			t.Fatal(planErr)
		}
		segmentBytes, segment, candidate, buildErr := buildArchiveCandidate(loaded, plan)
		if buildErr != nil {
			t.Fatal(buildErr)
		}
		segmentPath := filepath.Join(root, archiveSegmentName(segment.Segment.Digest()))
		if writeErr := os.WriteFile(segmentPath, segmentBytes, 0o400); writeErr != nil {
			t.Fatal(writeErr)
		}
		loaded = exactArchiveCapacitySuccessor(t, loaded, segment, candidate)
		if got := len(loaded.Segments); got != ordinal || len(candidate.Descriptors) != ordinal {
			t.Fatalf("segment count after ordinal %d = %d", ordinal, got)
		}
	}
	writeV2Envelope(t, path, loaded.Envelope)
	store, err := OpenFixedFileStoreV2WithOptions(path, FixedFileStoreV1Options{
		EffectIDs: &lifecycleEffectIDSequence{next: 401}, OwnerSessionID: owner.session,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(store.segments); got != archivestate.MaxReferencedArchiveSegments-1 {
		t.Fatalf("segment count before production boundary activation = %d", got)
	}
	store = activateNextArchiveSegment(t, store, owner)
	if got := len(store.segments); got != archivestate.MaxReferencedArchiveSegments {
		t.Fatalf("segment count after production boundary activation = %d", got)
	}
	return path, store, owner
}

func exactArchiveCapacitySuccessor(
	t *testing.T,
	prior loadedV2State,
	segment loadedArchiveSegmentV0,
	candidate diskEnvelopeV2,
) loadedV2State {
	t.Helper()
	indexes, err := archiveIndexesFromDisk(candidate.Indexes)
	if err != nil {
		t.Fatal(err)
	}
	tombstones, tombstoneDigest, _, err := effectTombstonesFromEnvelope(candidate, indexes)
	if err != nil {
		t.Fatal(err)
	}
	descriptors := make([]archivestate.ArchiveDescriptor, len(candidate.Descriptors))
	for index, view := range candidate.Descriptors {
		descriptors[index], err = archivestate.NewArchiveDescriptor(view)
		if err != nil {
			t.Fatal(err)
		}
	}
	active, err := buildActiveStateV2(candidate, descriptors, indexes, tombstoneDigest)
	if err != nil {
		t.Fatal(err)
	}
	lifecycles := make([]lifecyclestate.Record, len(candidate.Lifecycles))
	for index, disk := range candidate.Lifecycles {
		lifecycles[index], err = restoreLifecycleRecord(disk, candidate.State.TimeHighWaterUnixSeconds)
		if err != nil {
			t.Fatal(err)
		}
	}
	segments := append(cloneLoadedArchiveSegments(prior.Segments), segment)
	if err := validateActivationCheckpointDisk(candidate, indexes, segments); err != nil {
		t.Fatal(err)
	}
	encoded, err := encodeEnvelopeV2(candidate)
	if err != nil {
		t.Fatal(err)
	}
	return loadedV2State{
		Envelope: candidate, ActiveEncodedBytes: uint64(len(encoded)), State: cloneState(candidate.State),
		Lifecycles:       cloneLifecycleRecords(lifecycles, candidate.State.TimeHighWaterUnixSeconds),
		EffectTombstones: cloneEffectTombstones(tombstones), Active: active, Genesis: prior.Genesis,
		Segments: segments,
	}
}

type growthInventory struct {
	active [32]byte
	files  [32]byte
}

func inventoryDigest(t *testing.T, path string) growthInventory {
	t.Helper()
	active := sha256.Sum256(mustReadFile(t, path))
	hash := sha256.New()
	entries, err := os.ReadDir(archiveRoot(path))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return growthInventory{active: active, files: sha256.Sum256(nil)}
		}
		t.Fatal(err)
	}
	for _, entry := range entries {
		_, _ = hash.Write([]byte(entry.Name()))
		_, _ = hash.Write(mustReadFile(t, filepath.Join(archiveRoot(path), entry.Name())))
	}
	var files [32]byte
	copy(files[:], hash.Sum(nil))
	return growthInventory{active: active, files: files}
}
