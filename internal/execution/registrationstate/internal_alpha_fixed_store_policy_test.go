package registrationstate

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

type internalAlphaClockStub struct {
	times []time.Time
	next  int
}

func (clock *internalAlphaClockStub) Now() time.Time {
	value := clock.times[clock.next]
	clock.next++
	return value
}

func internalAlphaTestClock(duration time.Duration) *internalAlphaClockStub {
	start := time.Unix(1_800_000_000, 0)
	return &internalAlphaClockStub{times: []time.Time{start, start.Add(duration)}}
}

func allowedInternalAlphaFacts() internalAlphaStoreFacts {
	return internalAlphaStoreFacts{
		known: true, cumulativeAttempts: InternalAlphaMaxCumulativeAttempts - 1,
		activeStoreBytes:            InternalAlphaMaxActiveStoreBytes - 1,
		archiveSegments:             InternalAlphaMaxArchiveSegments - 1,
		startupVerification:         InternalAlphaMaxStartupVerification,
		startupVerificationObserved: true,
		durableCommitP95:            InternalAlphaMaxDurableCommitP95,
		durableCommitP95Observed:    true,
	}
}

func TestInternalAlphaFixedStorePolicyExactBoundariesAndCapPlusOne(t *testing.T) {
	if report := evaluateInternalAlphaFixedStoreFacts(allowedInternalAlphaFacts()); !report.Allowed || report.StopReason != "" {
		t.Fatalf("last admitted boundary = %#v", report)
	}

	for _, test := range []struct {
		name   string
		reason InternalAlphaStoreStopReason
		mutate func(*internalAlphaStoreFacts)
	}{
		{"attempt-exact", InternalAlphaStoreStopCumulativeAttempts, func(facts *internalAlphaStoreFacts) {
			facts.cumulativeAttempts = InternalAlphaMaxCumulativeAttempts
		}},
		{"attempt-cap-plus-one", InternalAlphaStoreStopCumulativeAttempts, func(facts *internalAlphaStoreFacts) {
			facts.cumulativeAttempts = InternalAlphaMaxCumulativeAttempts + 1
		}},
		{"active-bytes-exact", InternalAlphaStoreStopActiveStoreBytes, func(facts *internalAlphaStoreFacts) {
			facts.activeStoreBytes = InternalAlphaMaxActiveStoreBytes
		}},
		{"active-bytes-cap-plus-one", InternalAlphaStoreStopActiveStoreBytes, func(facts *internalAlphaStoreFacts) {
			facts.activeStoreBytes = InternalAlphaMaxActiveStoreBytes + 1
		}},
		{"segments-exact", InternalAlphaStoreStopArchiveSegments, func(facts *internalAlphaStoreFacts) {
			facts.archiveSegments = InternalAlphaMaxArchiveSegments
		}},
		{"segments-cap-plus-one", InternalAlphaStoreStopArchiveSegments, func(facts *internalAlphaStoreFacts) {
			facts.archiveSegments = InternalAlphaMaxArchiveSegments + 1
		}},
		{"startup-over", InternalAlphaStoreStopStartupVerification, func(facts *internalAlphaStoreFacts) {
			facts.startupVerification = InternalAlphaMaxStartupVerification + time.Nanosecond
		}},
		{"commit-p95-over", InternalAlphaStoreStopDurableCommitP95, func(facts *internalAlphaStoreFacts) {
			facts.durableCommitP95 = InternalAlphaMaxDurableCommitP95 + time.Nanosecond
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			facts := allowedInternalAlphaFacts()
			test.mutate(&facts)
			report := evaluateInternalAlphaFixedStoreFacts(facts)
			if report.Allowed || report.StopReason != test.reason {
				t.Fatalf("boundary result = %#v, want %s", report, test.reason)
			}
		})
	}
}

func TestInternalAlphaFixedStorePolicyUnknownStatePrecedesNumericStops(t *testing.T) {
	facts := allowedInternalAlphaFacts()
	facts.known = false
	facts.cumulativeAttempts = InternalAlphaMaxCumulativeAttempts + 1
	facts.activeStoreBytes = InternalAlphaMaxActiveStoreBytes + 1
	facts.archiveSegments = InternalAlphaMaxArchiveSegments + 1
	facts.startupVerification = InternalAlphaMaxStartupVerification + time.Nanosecond
	facts.durableCommitP95 = InternalAlphaMaxDurableCommitP95 + time.Nanosecond

	report := evaluateInternalAlphaFixedStoreFacts(facts)
	if report.Allowed || report.StopReason != InternalAlphaStoreStopUnknownState {
		t.Fatalf("combined unknown/numeric result = %#v", report)
	}
}

func TestInternalAlphaFixedStorePolicyUnknownStateAndCorruptionRefuseWithoutRewrite(t *testing.T) {
	path := newV1PathFromState(t, ordinaryInitialStateForV1(), nil)
	mustMigrateV2(t, path)
	root, err := OpenStoreRoot(path)
	if err != nil {
		t.Fatal(err)
	}
	owner := archiveOwnerForTest(t)
	startup, _, err := verifyInternalAlphaFixedStoreStartupWithClock(
		context.Background(), root, owner, internalAlphaTestClock(time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("missing-p95", func(t *testing.T) {
		before := inventoryDigest(t, path)
		report, checkErr := CheckInternalAlphaFixedStoreAttemptAdmission(
			context.Background(), root, owner, startup, InternalAlphaDurableCommitP95{},
		)
		if !errors.Is(checkErr, ErrInternalAlphaAttemptRefused) || report.StopReason != InternalAlphaStoreStopUnknownState {
			t.Fatalf("missing p95 result = %#v, %v", report, checkErr)
		}
		if after := inventoryDigest(t, path); after != before {
			t.Fatal("unknown timing refusal rewrote store or archive bytes")
		}
	})

	t.Run("corrupt-active", func(t *testing.T) {
		corruptPath := filepath.Join(t.TempDir(), "corrupt-state.json")
		data := mustReadFile(t, path)
		data[len(data)/2] ^= 0xff
		if err := os.WriteFile(corruptPath, data, 0o600); err != nil {
			t.Fatal(err)
		}
		corruptRoot, err := OpenStoreRoot(corruptPath)
		if err != nil {
			t.Fatal(err)
		}
		before := mustReadFile(t, corruptPath)
		p95, err := ObserveInternalAlphaDurableCommitP95(InternalAlphaMaxDurableCommitP95 + time.Nanosecond)
		if err != nil {
			t.Fatal(err)
		}
		report, checkErr := CheckInternalAlphaFixedStoreAttemptAdmission(
			context.Background(), corruptRoot, owner, startup, p95,
		)
		if !errors.Is(checkErr, ErrInternalAlphaAttemptRefused) || !errors.Is(checkErr, ErrStoreRepairRequired) ||
			report.StopReason != InternalAlphaStoreStopUnknownState {
			t.Fatalf("corrupt result = %#v, %v", report, checkErr)
		}
		if after := mustReadFile(t, corruptPath); !bytes.Equal(after, before) {
			t.Fatal("corruption refusal rewrote evidence")
		}
	})
}

func TestInternalAlphaFixedStorePolicyCountsOnlyPrimaryActiveEncodedBytes(t *testing.T) {
	path, store, owner := newEligibleFixedStoreV2(t)
	keys := lookupKeysForStore(t, store)
	verified := mustPreparedArchive(t, store, owner)
	reopened, err := store.ActivateArchive(context.Background(), owner, verified, nil)
	if err != nil {
		t.Fatal(err)
	}
	segmentBytes := mustReadFile(t, archiveSegmentPath(path, verified.SegmentDigest()))
	if len(segmentBytes) == 0 {
		t.Fatal("segment-bearing fixture persisted empty referenced segment bytes")
	}
	approvalBefore, err := reopened.ResolveApproval(context.Background(), keys.approvalID)
	if err != nil {
		t.Fatal(err)
	}
	before := inventoryDigest(t, path)
	activeBytes := uint64(len(mustReadFile(t, path)))

	root, err := OpenStoreRoot(path)
	if err != nil {
		t.Fatal(err)
	}
	startup, _, err := verifyInternalAlphaFixedStoreStartupWithClock(
		context.Background(), root, owner, internalAlphaTestClock(time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	p95, err := ObserveInternalAlphaDurableCommitP95(200 * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	report, err := CheckInternalAlphaFixedStoreAttemptAdmission(
		context.Background(), root, owner, startup, p95,
	)
	if err != nil || !report.Allowed || report.ArchiveSegments != 1 {
		t.Fatalf("segment-bearing admission = %#v, %v", report, err)
	}
	if report.ActiveStoreBytes != activeBytes {
		t.Fatalf("active encoded bytes = %d, want primary file length %d", report.ActiveStoreBytes, activeBytes)
	}
	if report.ActiveStoreBytes >= activeBytes+uint64(len(segmentBytes)) {
		t.Fatal("active encoded bytes did not exclude referenced archive segment bytes")
	}
	if after := inventoryDigest(t, path); after != before {
		t.Fatal("segment-bearing policy check rewrote active or archive bytes")
	}
	approvalAfter, err := reopened.ResolveApproval(context.Background(), keys.approvalID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(approvalAfter, approvalBefore) {
		t.Fatalf("passive policy changed retained approval:\nbefore=%#v\nafter=%#v", approvalBefore, approvalAfter)
	}
}

func TestInternalAlphaFixedStorePolicyRestartAndAdmissionAreReadOnly(t *testing.T) {
	path := newV1PathFromState(t, ordinaryInitialStateForV1(), nil)
	mustMigrateV2(t, path)
	p95, err := ObserveInternalAlphaDurableCommitP95(200 * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	before := inventoryDigest(t, path)
	var previousStartup InternalAlphaStartupVerification

	for restart := 0; restart < 2; restart++ {
		owner := &archiveOwnerStub{held: true, session: lifecycleOwnerID(t, byte(0xa0+restart))}
		root, openErr := OpenStoreRoot(path)
		if openErr != nil {
			t.Fatal(openErr)
		}
		if restart != 0 {
			staleReport, staleErr := CheckInternalAlphaFixedStoreAttemptAdmission(
				context.Background(), root, owner, previousStartup, p95,
			)
			if !errors.Is(staleErr, ErrInternalAlphaAttemptRefused) || staleReport.StopReason != InternalAlphaStoreStopUnknownState {
				t.Fatalf("restart %d stale owner-session result = %#v, %v", restart, staleReport, staleErr)
			}
		}
		startup, startupReport, verifyErr := verifyInternalAlphaFixedStoreStartupWithClock(
			context.Background(), root, owner, internalAlphaTestClock(2*time.Second),
		)
		if verifyErr != nil || startupReport.Posture != InternalAlphaFixedStorePosture {
			t.Fatalf("restart %d startup = %#v, %v", restart, startupReport, verifyErr)
		}
		report, checkErr := CheckInternalAlphaFixedStoreAttemptAdmission(
			context.Background(), root, owner, startup, p95,
		)
		if checkErr != nil || !report.Allowed || report.CumulativeAttempts != 0 ||
			report.ActiveStoreBytes != uint64(len(mustReadFile(t, path))) || report.ArchiveSegments != 0 {
			t.Fatalf("restart %d admission = %#v, %v", restart, report, checkErr)
		}
		previousStartup = startup
	}

	if after := inventoryDigest(t, path); after != before {
		t.Fatal("startup/attempt policy checks rewrote authority or archive bytes")
	}
	reopened, err := OpenFixedFileStoreV2(path)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := reopened.snapshotV2(context.Background())
	if err != nil || len(snapshot.State.Attempts) != 0 || snapshot.State.AttemptsDisabled {
		t.Fatalf("read-only policy changed attempt authority = %#v, %v", snapshot.State, err)
	}
}

func TestInternalAlphaFixedStoreStartupOverTwoSecondsAndKnownOrphanRefuse(t *testing.T) {
	t.Run("startup-over-two-seconds", func(t *testing.T) {
		path := newV1PathFromState(t, ordinaryInitialStateForV1(), nil)
		mustMigrateV2(t, path)
		root, err := OpenStoreRoot(path)
		if err != nil {
			t.Fatal(err)
		}
		before := inventoryDigest(t, path)
		startup, report, verifyErr := verifyInternalAlphaFixedStoreStartupWithClock(
			context.Background(), root, archiveOwnerForTest(t), internalAlphaTestClock(2*time.Second+time.Nanosecond),
		)
		if !errors.Is(verifyErr, ErrInternalAlphaAttemptRefused) || startup.observed ||
			report.StopReason != InternalAlphaStoreStopStartupVerification {
			t.Fatalf("slow startup = %#v %#v, %v", startup, report, verifyErr)
		}
		if after := inventoryDigest(t, path); after != before {
			t.Fatal("slow startup refusal rewrote store")
		}
	})

	t.Run("known-unreferenced-orphan", func(t *testing.T) {
		path, store, owner := newEligibleFixedStoreV2(t)
		verified := mustPreparedArchive(t, store, owner)
		if _, err := store.ActivateArchive(
			context.Background(), owner, verified, archiveFaultStub{point: FaultArchiveAfterSegmentPublish},
		); !errors.Is(err, ErrArchiveOutcomeIndeterminate) {
			t.Fatalf("orphan fixture activation = %v", err)
		}
		root, err := OpenStoreRoot(path)
		if err != nil {
			t.Fatal(err)
		}
		before := inventoryDigest(t, path)
		startup, report, verifyErr := verifyInternalAlphaFixedStoreStartupWithClock(
			context.Background(), root, owner, internalAlphaTestClock(time.Second),
		)
		if !errors.Is(verifyErr, ErrInternalAlphaAttemptRefused) || startup.observed ||
			report.StopReason != InternalAlphaStoreStopUnknownState {
			t.Fatalf("known orphan startup = %#v %#v, %v", startup, report, verifyErr)
		}
		if after := inventoryDigest(t, path); after != before {
			t.Fatal("known orphan refusal rewrote or deleted evidence")
		}
	})
}
