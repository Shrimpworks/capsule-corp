package archivestate

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

func TestArchivePassiveKnownAnswersAndDefensiveCopies(t *testing.T) {
	empty := EmptyArchiveIndexes()
	assertHex(t, "empty registration index", empty.Digests().Registrations, "8144fc094dd2a921cf357c39aba296fba1fc1a35eec8ee9f378c221781f96201")
	assertHex(t, "empty combined index", empty.CombinedDigest(), "2dd78bdddb4e186229d709bdda5b666e4e2d668e5c1216c751be2f4abb46648e")
	descriptorDigest, err := DigestArchiveDescriptorSet([]ArchiveDescriptor{})
	if err != nil {
		t.Fatal(err)
	}
	assertHex(t, "empty descriptor set", descriptorDigest, "a84af9da7e16fadb5aa76f4385558d4bc622ed1ea32ef435899ff02c20e863b3")

	checkpoint, err := NewArchiveCheckpoint(ArchiveCheckpointView{
		SourceSnapshotGeneration: 1, ResultSnapshotGeneration: 2, ArchiveGeneration: 1,
		InstallationID: installationID(1), SupervisorID: supervisorID(2), EpochSequence: 3,
		EpochDigest: epochDigest(4), DurableTimeHighWater: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertHex(t, "empty checkpoint", checkpoint.Digest(), "69dced13926ca3bfdf7324f7862035480af3b6d26c8b942a36a1ec0db5ee7d54")

	one := mustCohortProjection(t, closedCohort(t, 1, 90, 1))
	oneSegment := mustSegment(t, []CohortProjection{one}, 1_024)
	assertHex(t, "one cohort", one.Digest(), "67f6a0da8db580bb9f93b38b3d2d600bc667a49184aae1b7cdbf315ff7603198")
	assertHex(t, "one cohort segment", oneSegment.Digest(), "9a9d4e162edb7e493a7b0255ce87c3eb085d5fb865ca44ff4feb5b077bb66cce")

	multi := mustCohortProjection(t, closedCohort(t, 2, 90, 2))
	multiSegment := mustSegment(t, []CohortProjection{one, multi}, 2_048)
	assertHex(t, "multi attempt cohort", multi.Digest(), "41464e4a282e71fbfc2905cde5395bd55bf8bae8050a80ad56ae39fb0a4fbfb0")
	assertHex(t, "multi cohort segment", multiSegment.Digest(), "d84aa296b8a3b8a5f90fe2824b8371a165428979bb1b39447ec6fa2354362cf8")

	registrationDigest, err := DigestRecord(RecordRegistration, []byte("record"))
	if err != nil {
		t.Fatal(err)
	}
	approvalDigest, err := DigestRecord(RecordApproval, []byte("record"))
	if err != nil {
		t.Fatal(err)
	}
	if registrationDigest == approvalDigest {
		t.Fatal("wrong record domains collided")
	}
	if _, err := DigestRecord(RecordKind("unknown"), []byte("record")); classification(t, err) != ClassificationDomain {
		t.Fatalf("unknown record domain = %v", err)
	}

	view := validStateView(t, []RegistrationCohortCandidate{closedCohort(t, 3, 90, 1)})
	state, err := NewValidatedV1OrV2State(view)
	if err != nil {
		t.Fatal(err)
	}
	want := append([]byte(nil), state.view.Cohorts[0].Registration.ExactRecordBytes...)
	view.Cohorts[0].Registration.ExactRecordBytes[0] ^= 0xff
	projected := state.View()
	projected.Cohorts[0].Registration.ExactRecordBytes[0] ^= 0xff
	if got := state.View().Cohorts[0].Registration.ExactRecordBytes; string(got) != string(want) {
		t.Fatal("validated state aliases caller or accessor bytes")
	}

	indexView := EmptyArchiveIndexes().View()
	indexView.Nonces = []NonceIndexEntry{{AttemptNonce: nonce(1), PayloadDigest: payloadDigest(1), ApprovalID: approvalID(1)}}
	indexes, err := NewArchiveIndexes(indexView)
	if err != nil {
		t.Fatal(err)
	}
	returned := indexes.View()
	returned.Nonces[0].AttemptNonce[15] ^= 0xff
	if indexes.View().Nonces[0].AttemptNonce != nonce(1) {
		t.Fatal("archive indexes alias accessor slice")
	}
}

func TestArchiveSelectorAcceptsOnlyCompleteClosedRegistrationCohorts(t *testing.T) {
	tests := []struct {
		name         string
		cohort       func(*testing.T) RegistrationCohortCandidate
		highWater    v0candidate.UInt53
		wantSelected int
	}{
		{name: "expired-no-approval", cohort: func(t *testing.T) RegistrationCohortCandidate { return closedCohort(t, 1, 90, 0) }, highWater: 100, wantSelected: 1},
		{name: "exact-expiry", cohort: func(t *testing.T) RegistrationCohortCandidate { return closedCohort(t, 1, 100, 0) }, highWater: 100, wantSelected: 1},
		{name: "one-second-before-expiry", cohort: func(t *testing.T) RegistrationCohortCandidate { return closedCohort(t, 1, 101, 0) }, highWater: 100},
		{name: "expired-usable-not-rewritten", cohort: func(t *testing.T) RegistrationCohortCandidate {
			return approvalOnlyCohort(t, 1, 90, approvalattempt.ApprovalUsable, 90)
		}, highWater: 100, wantSelected: 1},
		{name: "usable-unexpired", cohort: func(t *testing.T) RegistrationCohortCandidate {
			return approvalOnlyCohort(t, 1, 200, approvalattempt.ApprovalUsable, 101)
		}, highWater: 100},
		{name: "invalidated", cohort: func(t *testing.T) RegistrationCohortCandidate {
			return approvalOnlyCohort(t, 1, 90, approvalattempt.ApprovalInvalidated, 90)
		}, highWater: 100, wantSelected: 1},
		{name: "consumed-destroyed-absent", cohort: func(t *testing.T) RegistrationCohortCandidate { return closedCohort(t, 1, 90, 1) }, highWater: 100, wantSelected: 1},
		{name: "observed", cohort: func(t *testing.T) RegistrationCohortCandidate {
			return activeCohort(t, 1, 90, lifecyclestate.StateObserved)
		}, highWater: 100},
		{name: "stopped", cohort: func(t *testing.T) RegistrationCohortCandidate {
			return activeCohort(t, 1, 90, lifecyclestate.StateStopped)
		}, highWater: 100},
		{name: "destroy-confirmed", cohort: func(t *testing.T) RegistrationCohortCandidate {
			return activeCohort(t, 1, 90, lifecyclestate.StateDestroyConfirmed)
		}, highWater: 100},
		{name: "unresolved", cohort: func(t *testing.T) RegistrationCohortCandidate {
			return activeCohort(t, 1, 90, lifecyclestate.StateUnresolved)
		}, highWater: 100},
		{name: "quarantined", cohort: func(t *testing.T) RegistrationCohortCandidate {
			return activeCohort(t, 1, 90, lifecyclestate.StateQuarantined)
		}, highWater: 100},
		{name: "recovery-fenced-destroyed", cohort: func(t *testing.T) RegistrationCohortCandidate {
			cohort := closedCohort(t, 1, 90, 1)
			cohort.Lifecycles[0].RecoveryFence = lifecyclestate.RecoveryFenceCommitFailure
			return cohort
		}, highWater: 100},
		{name: "automatic-recovery-exhausted", cohort: func(t *testing.T) RegistrationCohortCandidate {
			cohort := activeCohort(t, 1, 90, lifecyclestate.StateUnresolved)
			cohort.Lifecycles[0].AutomaticRecoveryExhausted = true
			cohort.Lifecycles[0].RecoveryFence = lifecyclestate.RecoveryFenceAutomaticExhausted
			return cohort
		}, highWater: 100},
	}
	limits, _ := NewArchiveLimits(256, uint64(MaxSupervisorArchiveBytes))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cohort := test.cohort(t)
			view := validStateView(t, []RegistrationCohortCandidate{cohort})
			view.DurableTimeHighWater = test.highWater
			view.RecoveryAttemptIDs = recoveryIDs(cohort)
			state, err := NewValidatedV1OrV2State(view)
			if err != nil {
				t.Fatal(err)
			}
			plan, err := SelectClosedRegistrationCohorts(state, limits)
			if err != nil {
				t.Fatal(err)
			}
			if got := len(plan.View().Selected); got != test.wantSelected {
				t.Fatalf("selected = %d, want %d", got, test.wantSelected)
			}
			if test.name == "expired-usable-not-rewritten" && state.View().Cohorts[0].Approvals[0].State != approvalattempt.ApprovalUsable {
				t.Fatal("selection rewrote expired durable usable approval")
			}
		})
	}

	held := closedCohort(t, 1, 90, 0)
	held.RetentionHold = true
	state, _ := NewValidatedV1OrV2State(validStateView(t, []RegistrationCohortCandidate{held}))
	plan, _ := SelectClosedRegistrationCohorts(state, limits)
	if len(plan.View().Selected) != 0 {
		t.Fatal("retention-held cohort selected")
	}

	for _, fence := range []ArchiveFence{ArchiveFenceTrustTransition, ArchiveFenceQuarantined, ArchiveFenceRepairRequired, ArchiveFenceRecovery} {
		view := validStateView(t, []RegistrationCohortCandidate{closedCohort(t, 1, 90, 0)})
		view.Fence = fence
		state, err := NewValidatedV1OrV2State(view)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := SelectClosedRegistrationCohorts(state, limits); classification(t, err) != ClassificationTrustState {
			t.Fatalf("fence %s = %v", fence, err)
		}
	}

	broken := closedCohort(t, 1, 90, 1)
	broken.Lifecycles = nil
	view := validStateView(t, []RegistrationCohortCandidate{})
	view.Cohorts = []RegistrationCohortCandidate{broken}
	view.RecoveryAttemptIDs = []approvalattempt.AttemptID{}
	if _, err := NewValidatedV1OrV2State(view); classification(t, err) != ClassificationRepairNeeded {
		t.Fatalf("missing lifecycle half = %v", err)
	}
	broken = closedCohort(t, 1, 90, 1)
	broken.Attempts = []AttemptCandidate{}
	broken.Lifecycles = []LifecycleCandidate{}
	broken.EncodedByteLength = minimumCohortBytes(broken)
	view.Cohorts = []RegistrationCohortCandidate{broken}
	if _, err := NewValidatedV1OrV2State(view); classification(t, err) != ClassificationRepairNeeded {
		t.Fatalf("missing attempt half = %v", err)
	}
	broken = closedCohort(t, 1, 90, 1)
	broken.Approvals = []ApprovalCandidate{}
	broken.EncodedByteLength = minimumCohortBytes(broken)
	view.Cohorts = []RegistrationCohortCandidate{broken}
	if _, err := NewValidatedV1OrV2State(view); classification(t, err) != ClassificationRepairNeeded {
		t.Fatalf("missing approval half = %v", err)
	}
	for _, mutation := range []struct {
		name  string
		apply func(*LifecycleCandidate)
	}{
		{name: "destroyed-cleanup-true", apply: func(record *LifecycleCandidate) { record.CleanupRequired = true }},
		{name: "destroyed-absence-not-authoritative", apply: func(record *LifecycleCandidate) { record.LastReconciliation = lifecyclestate.ReconciliationUnknown }},
	} {
		broken = closedCohort(t, 1, 90, 1)
		mutation.apply(&broken.Lifecycles[0])
		view.Cohorts = []RegistrationCohortCandidate{broken}
		view.RecoveryAttemptIDs = recoveryIDs(broken)
		if _, err := NewValidatedV1OrV2State(view); classification(t, err) != ClassificationRepairNeeded {
			t.Fatalf("%s = %v", mutation.name, err)
		}
	}
}

func TestArchiveClosedIndexesActiveProjectionAndMutationRefusals(t *testing.T) {
	indexes := completeArchiveIndexes(t)
	assertHex(t, "one-entry combined index", indexes.CombinedDigest(), "eb3ac7659f3504441fd507397762b221b8596c9d4f7d642fbd4ca3e1679abb83")
	emptySeed, err := DigestVisibleV1EffectSeed([]lifecyclestate.EffectID{})
	if err != nil {
		t.Fatal(err)
	}
	assertHex(t, "empty visible-v1 effect seed", emptySeed, "17de5f44f523dab94ca4b215ce7779358146fb094fa6d208e0190cb0ba69e0a1")

	descriptors := []ArchiveDescriptor{}
	descriptorDigest, err := DigestArchiveDescriptorSet(descriptors)
	if err != nil {
		t.Fatal(err)
	}
	zeroCounts := ArchiveCounts{}
	activeIndexes := EmptyArchiveIndexes()
	activeView := ActiveStateV2View{
		StoreFormatVersion:        SupervisorStoreFormatV2,
		MigrationSourceVersion:    MigrationSourceFormatV1,
		SnapshotGeneration:        1,
		ArchiveGeneration:         1,
		InstallationID:            installationID(1),
		SupervisorID:              supervisorID(2),
		EpochSequence:             3,
		EpochDigest:               epochDigest(4),
		DurableTimeHighWater:      100,
		Descriptors:               descriptors,
		DescriptorSetDigest:       descriptorDigest,
		Indexes:                   activeIndexes,
		IndexDigests:              activeIndexes.Digests(),
		CombinedIndexDigest:       activeIndexes.CombinedDigest(),
		EffectTombstoneCoverage:   EffectTombstoneCoverageVisibleV1AndV2,
		VisibleV1EffectSeedDigest: emptySeed,
		PreviousCheckpointDigest:  checkpointDigest(8),
		CurrentCheckpointDigest:   checkpointDigest(9),
		HotCounts:                 zeroCounts,
		ArchivedCounts:            zeroCounts,
		TotalCounts:               zeroCounts,
	}
	active, err := NewActiveStateV2(activeView)
	if err != nil {
		t.Fatal(err)
	}
	returned := active.View()
	returned.Descriptors = append(returned.Descriptors, ArchiveDescriptor{})
	if len(active.View().Descriptors) != 0 {
		t.Fatal("active v2 descriptors alias accessor")
	}

	mutations := []struct {
		name  string
		apply func(*ActiveStateV2View)
		want  Classification
	}{
		{name: "wrong-version", apply: func(view *ActiveStateV2View) { view.StoreFormatVersion = 1 }, want: ClassificationUnsupported},
		{name: "descriptor-digest", apply: func(view *ActiveStateV2View) { view.DescriptorSetDigest[0] ^= 1 }, want: ClassificationBinding},
		{name: "index-digest", apply: func(view *ActiveStateV2View) { view.CombinedIndexDigest[0] ^= 1 }, want: ClassificationBinding},
		{name: "coverage", apply: func(view *ActiveStateV2View) { view.EffectTombstoneCoverage = "all-history" }, want: ClassificationBinding},
		{name: "counts", apply: func(view *ActiveStateV2View) { view.TotalCounts.Attempts = 1 }, want: ClassificationBinding},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			candidate := activeView
			mutation.apply(&candidate)
			if _, err := NewActiveStateV2(candidate); classification(t, err) != mutation.want {
				t.Fatalf("mutation = %v", err)
			}
		})
	}

	cohort := mustCohortProjection(t, closedCohort(t, 1, 90, 1))
	wrongDomain := cohort.View()
	wrongDomain.Approvals[0].Kind = RecordAttempt
	if _, err := NewCohortProjection(wrongDomain); classification(t, err) != ClassificationDomain {
		t.Fatalf("cohort wrong domain = %v", err)
	}
	duplicate := cohort.View()
	duplicate.Approvals = append(duplicate.Approvals, duplicate.Approvals[0])
	if _, err := NewCohortProjection(duplicate); classification(t, err) != ClassificationBinding {
		t.Fatalf("cohort digest collision = %v", err)
	}

	segmentView := mustSegment(t, []CohortProjection{cohort}, 1_024).View()
	segmentView.FormatVersion = 1
	if _, err := NewArchiveSegment(segmentView); classification(t, err) != ClassificationUnsupported {
		t.Fatalf("segment version mutation = %v", err)
	}
	segmentView = mustSegment(t, []CohortProjection{cohort}, 1_024).View()
	segmentView.Counts.Attempts++
	if _, err := NewArchiveSegment(segmentView); classification(t, err) != ClassificationBinding {
		t.Fatalf("segment count mutation = %v", err)
	}
	segmentView = mustSegment(t, []CohortProjection{cohort}, 1_024).View()
	segmentView.SetDigests.Approvals[0] ^= 1
	if _, err := NewArchiveSegment(segmentView); classification(t, err) != ClassificationBinding {
		t.Fatalf("segment set digest mutation = %v", err)
	}
	segmentView = mustSegment(t, []CohortProjection{cohort}, 1_024).View()
	segmentView.DerivedIndexes = completeArchiveIndexes(t)
	if _, err := NewArchiveSegment(segmentView); classification(t, err) != ClassificationDomain {
		t.Fatalf("segment derived-index domain mutation = %v", err)
	}

	if _, err := NewArchiveCheckpoint(ArchiveCheckpointView{
		SourceSnapshotGeneration: 2, ResultSnapshotGeneration: 2, ArchiveGeneration: 1,
		InstallationID: installationID(1), SupervisorID: supervisorID(2),
	}); classification(t, err) != ClassificationBinding {
		t.Fatalf("checkpoint generation mutation = %v", err)
	}
}

func TestArchiveSelectorNeverSplitsMultiAttemptRegistration(t *testing.T) {
	cohort := closedCohort(t, 1, 90, 2)
	cohort.Lifecycles[1] = activeLifecycle(t, cohort.Attempts[1].AttemptID, lifecyclestate.StateObserved)
	cohort.EncodedByteLength = minimumCohortBytes(cohort)
	view := validStateView(t, []RegistrationCohortCandidate{cohort})
	view.RecoveryAttemptIDs = recoveryIDs(cohort)
	state, err := NewValidatedV1OrV2State(view)
	if err != nil {
		t.Fatal(err)
	}
	limits, _ := NewArchiveLimits(256, uint64(MaxSupervisorArchiveBytes))
	plan, err := SelectClosedRegistrationCohorts(state, limits)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.View().Selected) != 0 {
		t.Fatal("selector split active and destroyed attempts in one registration")
	}

	closed := closedCohort(t, 1, 90, 2)
	state, err = NewValidatedV1OrV2State(validStateView(t, []RegistrationCohortCandidate{closed}))
	if err != nil {
		t.Fatal(err)
	}
	plan, err = SelectClosedRegistrationCohorts(state, limits)
	if err != nil {
		t.Fatal(err)
	}
	if got := plan.View().Selected; len(got) != 1 || len(got[0].AttemptDigests) != 2 || len(got[0].LifecycleDigests) != 2 {
		t.Fatalf("multi-attempt selection = %#v", got)
	}
}

func TestArchiveSelectorAgreesWithRecoveryAttemptIDs(t *testing.T) {
	active := activeCohort(t, 1, 90, lifecyclestate.StateObserved)
	view := validStateView(t, []RegistrationCohortCandidate{active})
	view.RecoveryAttemptIDs = []approvalattempt.AttemptID{}
	if _, err := NewValidatedV1OrV2State(view); classification(t, err) != ClassificationBinding {
		t.Fatalf("missing recovery ID = %v", err)
	}

	closed := closedCohort(t, 1, 90, 1)
	view = validStateView(t, []RegistrationCohortCandidate{closed})
	view.RecoveryAttemptIDs = []approvalattempt.AttemptID{closed.Attempts[0].AttemptID}
	if _, err := NewValidatedV1OrV2State(view); classification(t, err) != ClassificationBinding {
		t.Fatalf("terminal ID in recovery = %v", err)
	}

	view = validStateView(t, []RegistrationCohortCandidate{active})
	view.RecoveryAttemptIDs = recoveryIDs(active)
	state, err := NewValidatedV1OrV2State(view)
	if err != nil {
		t.Fatal(err)
	}
	limits, _ := NewArchiveLimits(256, uint64(MaxSupervisorArchiveBytes))
	plan, err := SelectClosedRegistrationCohorts(state, limits)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.View().Selected) != 0 {
		t.Fatal("recovery-eligible attempt selected")
	}
}

func TestArchiveExactBoundsCapPlusOneOrderingAndCollisions(t *testing.T) {
	if _, err := NewArchiveLimits(MaxArchiveCohortsPerSegment, uint64(MaxSupervisorArchiveBytes)); err != nil {
		t.Fatal(err)
	}
	if _, err := NewArchiveLimits(MaxArchiveCohortsPerSegment+1, uint64(MaxSupervisorArchiveBytes)); classification(t, err) != ClassificationCapacity {
		t.Fatalf("cohort cap plus one = %v", err)
	}
	if _, err := NewArchiveLimits(MaxArchiveCohortsPerSegment, uint64(MaxSupervisorArchiveBytes)+1); classification(t, err) != ClassificationCapacity {
		t.Fatalf("byte cap plus one = %v", err)
	}

	descriptors := make([]ArchiveDescriptor, MaxReferencedArchiveSegments)
	for index := range descriptors {
		view := ArchiveDescriptorView{Ordinal: ArchiveOrdinal(index + 1), SourceSnapshotGeneration: 1, ArchiveGeneration: 1,
			SegmentDigest: segmentDigest(index + 1), EncodedByteLength: 1}
		descriptor, err := NewArchiveDescriptor(view)
		if err != nil {
			t.Fatal(err)
		}
		descriptors[index] = descriptor
	}
	if _, err := DigestArchiveDescriptorSet(descriptors); err != nil {
		t.Fatal(err)
	}
	extra, _ := NewArchiveDescriptor(ArchiveDescriptorView{Ordinal: MaxReferencedArchiveSegments + 1, SourceSnapshotGeneration: 1, ArchiveGeneration: 1, SegmentDigest: segmentDigest(65), EncodedByteLength: 1})
	if _, err := DigestArchiveDescriptorSet(append(descriptors, extra)); classification(t, err) != ClassificationCapacity {
		t.Fatalf("descriptor cap plus one = %v", err)
	}

	cohorts := make([]RegistrationCohortCandidate, MaxArchiveCohortsPerSegment+1)
	for index := range cohorts {
		cohorts[index] = closedCohort(t, index+1, 90, 0)
	}
	view := validStateView(t, cohorts)
	state, err := NewValidatedV1OrV2State(view)
	if err != nil {
		t.Fatal(err)
	}
	limits, _ := NewArchiveLimits(MaxArchiveCohortsPerSegment, uint64(MaxSupervisorArchiveBytes))
	plan, err := SelectClosedRegistrationCohorts(state, limits)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.View().Selected) != MaxArchiveCohortsPerSegment {
		t.Fatalf("exact cohort cap selected %d", len(plan.View().Selected))
	}

	exactTwoBytes := cohorts[0].EncodedByteLength + cohorts[1].EncodedByteLength
	limits, _ = NewArchiveLimits(2, exactTwoBytes)
	plan, err = SelectClosedRegistrationCohorts(state, limits)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.View().Selected) != 2 || plan.View().SelectedBytes != exactTwoBytes {
		t.Fatalf("exact byte cap = %#v", plan.View())
	}
	limits, _ = NewArchiveLimits(2, exactTwoBytes-1)
	plan, err = SelectClosedRegistrationCohorts(state, limits)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.View().Selected) != 1 || plan.View().Selected[0].RegistrationSequence != 1 {
		t.Fatalf("cap prefix = %#v", plan.View())
	}

	unsorted := []RegistrationCohortCandidate{closedCohort(t, 2, 90, 0), closedCohort(t, 1, 90, 0)}
	bad := validStateView(t, []RegistrationCohortCandidate{})
	bad.Cohorts = unsorted
	if _, err := NewValidatedV1OrV2State(bad); classification(t, err) != ClassificationBinding {
		t.Fatalf("unsorted cohorts = %v", err)
	}

	duplicate := closedCohort(t, 1, 90, 2)
	duplicate.Approvals[1].ApprovalID = duplicate.Approvals[0].ApprovalID
	bad = validStateView(t, []RegistrationCohortCandidate{})
	bad.Cohorts = []RegistrationCohortCandidate{duplicate}
	if _, err := NewValidatedV1OrV2State(bad); classification(t, err) != ClassificationBinding {
		t.Fatalf("approval collision = %v", err)
	}

	indexView := EmptyArchiveIndexes().View()
	indexView.Nonces = []NonceIndexEntry{
		{AttemptNonce: nonce(1), PayloadDigest: payloadDigest(1), ApprovalID: approvalID(1)},
		{AttemptNonce: nonce(1), PayloadDigest: payloadDigest(2), ApprovalID: approvalID(2)},
	}
	if _, err := NewArchiveIndexes(indexView); classification(t, err) != ClassificationBinding {
		t.Fatalf("nonce collision = %v", err)
	}
	complete := completeArchiveIndexes(t).View()
	complete.Effects = append(complete.Effects, complete.Effects[0])
	if _, err := NewArchiveIndexes(complete); classification(t, err) != ClassificationBinding {
		t.Fatalf("effect collision = %v", err)
	}
	complete = completeArchiveIndexes(t).View()
	complete.Instances = append(complete.Instances, complete.Instances[0])
	if _, err := NewArchiveIndexes(complete); classification(t, err) != ClassificationBinding {
		t.Fatalf("instance collision = %v", err)
	}
}

func TestArchiveSelectorDeterministicKnownAnswerAndNoSideEffects(t *testing.T) {
	cohorts := []RegistrationCohortCandidate{closedCohort(t, 1, 90, 0), closedCohort(t, 2, 90, 2)}
	state, err := NewValidatedV1OrV2State(validStateView(t, cohorts))
	if err != nil {
		t.Fatal(err)
	}
	limits, _ := NewArchiveLimits(2, uint64(MaxSupervisorArchiveBytes))
	first, err := SelectClosedRegistrationCohorts(state, limits)
	if err != nil {
		t.Fatal(err)
	}
	second, err := SelectClosedRegistrationCohorts(state, limits)
	if err != nil {
		t.Fatal(err)
	}
	if first.View().Digest != second.View().Digest {
		t.Fatal("identical selection changed digest")
	}
	assertHex(t, "selector", first.View().Digest, "e335a3823e433ec53917b45f3962e80041ac2ef2daaa6be1c303fed141440a12")
	if len(state.View().Cohorts) != 2 {
		t.Fatal("selector mutated source state")
	}
}

func validStateView(t *testing.T, cohorts []RegistrationCohortCandidate) ValidatedV1OrV2StateView {
	t.Helper()
	return ValidatedV1OrV2StateView{
		SnapshotGeneration: 7, ArchiveGeneration: 3, CheckpointHead: checkpointDigest(9),
		OwnerSessionID: ownerID(8), DurableTimeHighWater: 100, Fence: ArchiveFenceNone,
		Cohorts: cohorts, RecoveryAttemptIDs: recoveryIDs(cohorts...),
	}
}

func closedCohort(t *testing.T, sequence int, expiry uint64, attemptCount int) RegistrationCohortCandidate {
	t.Helper()
	registration := registrationCandidate(t, sequence, expiry)
	cohort := RegistrationCohortCandidate{Registration: registration, Approvals: make([]ApprovalCandidate, attemptCount), Attempts: make([]AttemptCandidate, attemptCount), Lifecycles: make([]LifecycleCandidate, attemptCount)}
	for index := range attemptCount {
		ordinal := sequence*32 + index + 1
		approval := approvalCandidate(t, registration, ordinal, approvalattempt.ApprovalConsumed, expiry)
		attempt := attemptCandidate(t, registration, approval, ordinal)
		approval.ConsumedAttemptID = attempt.AttemptID
		approval.ExactRecordBytes = recordBytes("approval", ordinal, byte(approval.State[0]))
		approval.RecordDigest = mustRecordDigest(t, RecordApproval, approval.ExactRecordBytes)
		cohort.Approvals[index] = approval
		cohort.Attempts[index] = attempt
		cohort.Lifecycles[index] = destroyedLifecycle(t, attempt.AttemptID)
	}
	cohort.EncodedByteLength = minimumCohortBytes(cohort)
	return cohort
}

func approvalOnlyCohort(t *testing.T, sequence int, registrationExpiry uint64, state approvalattempt.ApprovalState, approvalExpiry uint64) RegistrationCohortCandidate {
	t.Helper()
	registration := registrationCandidate(t, sequence, registrationExpiry)
	approval := approvalCandidate(t, registration, sequence*32+1, state, approvalExpiry)
	cohort := RegistrationCohortCandidate{Registration: registration, Approvals: []ApprovalCandidate{approval}, Attempts: []AttemptCandidate{}, Lifecycles: []LifecycleCandidate{}}
	cohort.EncodedByteLength = minimumCohortBytes(cohort)
	return cohort
}

func activeCohort(t *testing.T, sequence int, expiry uint64, state lifecyclestate.LifecycleState) RegistrationCohortCandidate {
	t.Helper()
	cohort := closedCohort(t, sequence, expiry, 1)
	cohort.Lifecycles[0] = activeLifecycle(t, cohort.Attempts[0].AttemptID, state)
	cohort.EncodedByteLength = minimumCohortBytes(cohort)
	return cohort
}

func registrationCandidate(t *testing.T, sequence int, expiry uint64) RegistrationCandidate {
	t.Helper()
	exact := recordBytes("registration", sequence, byte(expiry))
	return RegistrationCandidate{RegistrationID: registrationID(sequence), RegistrationSequence: v0candidate.PositiveUInt53(sequence), ExpiresAt: v0candidate.UInt53(expiry), ExactRecordBytes: exact, RecordDigest: mustRecordDigest(t, RecordRegistration, exact)}
}

func approvalCandidate(t *testing.T, registration RegistrationCandidate, ordinal int, state approvalattempt.ApprovalState, expiry uint64) ApprovalCandidate {
	t.Helper()
	exact := recordBytes("approval", ordinal, byte(state[0]))
	return ApprovalCandidate{ApprovalID: approvalID(ordinal), AttemptNonce: nonce(ordinal), RegistrationID: registration.RegistrationID,
		RegistrationSequence: registration.RegistrationSequence, ExpiresAt: v0candidate.UInt53(expiry), State: state,
		ExactRecordBytes: exact, RecordDigest: mustRecordDigest(t, RecordApproval, exact)}
}

func attemptCandidate(t *testing.T, registration RegistrationCandidate, approval ApprovalCandidate, ordinal int) AttemptCandidate {
	t.Helper()
	exact := recordBytes("attempt", ordinal, 0)
	return AttemptCandidate{AttemptID: attemptID(ordinal), ApprovalID: approval.ApprovalID, RegistrationID: registration.RegistrationID,
		RegistrationSequence: registration.RegistrationSequence, ExactRecordBytes: exact, RecordDigest: mustRecordDigest(t, RecordAttempt, exact)}
}

func destroyedLifecycle(t *testing.T, attempt approvalattempt.AttemptID) LifecycleCandidate {
	t.Helper()
	exact := append([]byte("lifecycle-destroyed-"), attempt[:]...)
	return LifecycleCandidate{AttemptID: attempt, State: lifecyclestate.StateDestroyed,
		TerminalAt:         lifecyclestate.OptionalUnixSeconds{Present: true, Value: 99},
		LastReconciliation: lifecyclestate.ReconciliationAuthoritativelyAbsent, RecoveryFence: lifecyclestate.RecoveryFenceNone,
		ExactRecordBytes: exact, RecordDigest: mustRecordDigest(t, RecordLifecycle, exact)}
}

func activeLifecycle(t *testing.T, attempt approvalattempt.AttemptID, state lifecyclestate.LifecycleState) LifecycleCandidate {
	t.Helper()
	exact := append(append([]byte("lifecycle-active-"), []byte(state)...), attempt[:]...)
	return LifecycleCandidate{AttemptID: attempt, State: state, CleanupRequired: true,
		LastReconciliation: lifecyclestate.ReconciliationNone, RecoveryFence: lifecyclestate.RecoveryFenceNone,
		ExactRecordBytes: exact, RecordDigest: mustRecordDigest(t, RecordLifecycle, exact)}
}

func minimumCohortBytes(cohort RegistrationCohortCandidate) uint64 {
	result := len(cohort.Registration.ExactRecordBytes)
	for _, approval := range cohort.Approvals {
		result += len(approval.ExactRecordBytes)
	}
	for _, attempt := range cohort.Attempts {
		result += len(attempt.ExactRecordBytes)
	}
	for _, lifecycle := range cohort.Lifecycles {
		result += len(lifecycle.ExactRecordBytes)
	}
	return uint64(result)
}

func recoveryIDs(cohorts ...RegistrationCohortCandidate) []approvalattempt.AttemptID {
	result := make([]approvalattempt.AttemptID, 0)
	for _, cohort := range cohorts {
		for _, lifecycle := range cohort.Lifecycles {
			if lifecycle.State != lifecyclestate.StateDestroyed || lifecycle.CleanupRequired || lifecycle.RecoveryFence != lifecyclestate.RecoveryFenceNone || lifecycle.NextRecoveryAt.Present {
				result = append(result, lifecycle.AttemptID)
			}
		}
	}
	for left := range result {
		for right := left + 1; right < len(result); right++ {
			if string(result[right][:]) < string(result[left][:]) {
				result[left], result[right] = result[right], result[left]
			}
		}
	}
	return result
}

func mustCohortProjection(t *testing.T, cohort RegistrationCohortCandidate) CohortProjection {
	t.Helper()
	planned, err := plannedCohort(cohort)
	if err != nil {
		t.Fatal(err)
	}
	projection, err := NewCohortProjection(CohortProjectionView{RegistrationID: planned.RegistrationID, RegistrationSequence: planned.RegistrationSequence,
		ExpiresAt: cohort.Registration.ExpiresAt, Registration: ArchiveRecordReference{Kind: RecordRegistration, Digest: planned.RegistrationDigest},
		Approvals: references(RecordApproval, planned.ApprovalDigests), Attempts: references(RecordAttempt, planned.AttemptDigests), Lifecycles: references(RecordLifecycle, planned.LifecycleDigests), EncodedByteLength: cohort.EncodedByteLength})
	if err != nil {
		t.Fatal(err)
	}
	return projection
}

func mustSegment(t *testing.T, cohorts []CohortProjection, encoded uint64) ArchiveSegment {
	t.Helper()
	derived := EmptyArchiveIndexes()
	segment, err := NewArchiveSegment(ArchiveSegmentView{FormatVersion: SupervisorArchiveFormatV0, Ordinal: 1,
		InstallationID: installationID(1), SupervisorID: supervisorID(2), EpochSequence: 3, EpochDigest: epochDigest(4),
		SourceSnapshotGeneration: 7, ArchiveGeneration: 3, DurableTimeHighWater: 100, PriorCheckpointDigest: checkpointDigest(9),
		Cohorts: cohorts, DerivedIndexes: derived, SetDigests: archiveSetDigests(cohorts, derived),
		Counts: countsForCohorts(cohorts), EncodedByteLength: encoded})
	if err != nil {
		t.Fatal(err)
	}
	return segment
}

func completeArchiveIndexes(t *testing.T) ArchiveIndexes {
	t.Helper()
	location := ArchiveLocation{SegmentOrdinal: 1, CohortOrdinal: 1, RecordOrdinal: 1}
	exactPayloadDigest := mustRecordDigest(t, RecordApproval, []byte("exact-payload"))
	view := ArchiveIndexesView{
		Registrations: []RegistrationIndexEntry{{
			RegistrationID: registrationID(1), RegistrationSequence: 1, PlanDigest: planDigest(1),
			ExpiresAt: 90, Location: location, FullRecordDigest: mustRecordDigest(t, RecordRegistration, []byte("registration")),
		}},
		Approvals: []ApprovalIndexEntry{{
			ApprovalID: approvalID(1), RegistrationID: registrationID(1), PayloadDigest: payloadDigest(1),
			AuthorizationIdentity: authorizationIdentity(1), AttemptNonce: nonce(1), State: approvalattempt.ApprovalConsumed,
			ConsumedAttemptID: attemptID(1), ExpiresAt: 90, Location: location,
			FullRecordDigest: mustRecordDigest(t, RecordApproval, []byte("approval")),
		}},
		Attempts: []AttemptIndexEntry{{
			AttemptID: attemptID(1), ApprovalID: approvalID(1), RegistrationID: registrationID(1), CreatedAt: 80,
			LifecycleState: lifecyclestate.StateDestroyed, Location: location,
			FullRecordDigest: mustRecordDigest(t, RecordAttempt, []byte("attempt")),
		}},
		Nonces: []NonceIndexEntry{{AttemptNonce: nonce(1), PayloadDigest: payloadDigest(1), ApprovalID: approvalID(1)}},
		Effects: []EffectIndexEntry{{EffectID: effectID(1), AttemptID: attemptID(1), OperationSequence: 1,
			Operation: lifecyclestate.OperationDestroy, IssuanceSnapshotGeneration: 1}},
		Instances: []InstanceIndexEntry{{InstanceDigest: instanceDigest(1), AttemptID: attemptID(1), Location: location}},
		ApprovalReplay: []ApprovalReplayIndexEntry{{PayloadDigest: payloadDigest(1), ExactPayloadDigest: exactPayloadDigest,
			AuthorizationIdentity: authorizationIdentity(1), ApprovalID: approvalID(1), State: approvalattempt.ApprovalConsumed, Location: location}},
		AttemptReplay: []AttemptReplayIndexEntry{{RegistrationID: registrationID(1), ApprovalID: approvalID(1), AttemptID: attemptID(1), State: approvalattempt.AttemptCreated}},
	}
	indexes, err := NewArchiveIndexes(view)
	if err != nil {
		t.Fatal(err)
	}
	return indexes
}

func mustRecordDigest(t *testing.T, kind RecordKind, exact []byte) ArchiveRecordDigest {
	t.Helper()
	digest, err := DigestRecord(kind, exact)
	if err != nil {
		t.Fatal(err)
	}
	return digest
}

func recordBytes(kind string, ordinal int, suffix byte) []byte {
	result := []byte(kind + "-record-")
	var encoded [4]byte
	binary.BigEndian.PutUint32(encoded[:], uint32(ordinal))
	return append(append(result, encoded[:]...), suffix)
}

func classification(t *testing.T, err error) Classification {
	t.Helper()
	value, _ := ErrorClassification(err)
	return value
}

func assertHex[T ~[32]byte](t *testing.T, name string, value T, want string) {
	t.Helper()
	got := hex.EncodeToString(value[:])
	if got != want {
		t.Fatalf("%s = %s, want %s", name, got, want)
	}
}

func registrationID(value int) (result v0candidate.RegistrationID) {
	binary.BigEndian.PutUint32(result[12:], uint32(value))
	return
}
func installationID(value int) (result v0candidate.InstallationID) {
	binary.BigEndian.PutUint32(result[12:], uint32(value))
	return
}
func supervisorID(value int) (result v0candidate.SupervisorID) {
	binary.BigEndian.PutUint32(result[12:], uint32(value))
	return
}
func approvalID(value int) (result approvalattempt.ApprovalID) {
	binary.BigEndian.PutUint32(result[12:], uint32(value))
	return
}
func attemptID(value int) (result approvalattempt.AttemptID) {
	binary.BigEndian.PutUint32(result[12:], uint32(value))
	return
}
func nonce(value int) (result approvalattempt.AttemptNonce) {
	binary.BigEndian.PutUint32(result[12:], uint32(value))
	return
}
func ownerID(value int) (result lifecyclestate.OwnerSessionID) {
	binary.BigEndian.PutUint32(result[12:], uint32(value))
	return
}
func payloadDigest(value int) (result approvalattempt.ApprovalPayloadDigest) {
	binary.BigEndian.PutUint32(result[28:], uint32(value))
	return
}
func epochDigest(value int) (result v0candidate.TrustEpochDigest) {
	binary.BigEndian.PutUint32(result[28:], uint32(value))
	return
}
func checkpointDigest(value int) (result ArchiveCheckpointDigest) {
	binary.BigEndian.PutUint32(result[28:], uint32(value))
	return
}
func segmentDigest(value int) (result ArchiveSegmentDigest) {
	binary.BigEndian.PutUint32(result[28:], uint32(value))
	return
}
func planDigest(value int) (result v0candidate.ExecutionPlanDigest) {
	binary.BigEndian.PutUint32(result[28:], uint32(value))
	return
}
func authorizationIdentity(value int) (result approvalattempt.ApprovalKeyAuthorizationIdentity) {
	binary.BigEndian.PutUint32(result[28:], uint32(value))
	return
}
func effectID(value int) (result lifecyclestate.EffectID) {
	binary.BigEndian.PutUint32(result[12:], uint32(value))
	return
}
func instanceDigest(value int) (result lifecyclestate.BackendInstanceDigest) {
	binary.BigEndian.PutUint32(result[28:], uint32(value))
	return
}
