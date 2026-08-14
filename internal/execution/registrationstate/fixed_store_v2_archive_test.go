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
	"capsule.local/capsule/internal/protocol/v0candidate"
)

type archiveOwnerStub struct {
	mu      sync.Mutex
	held    bool
	session lifecyclestate.OwnerSessionID
	checks  int
}

func (owner *archiveOwnerStub) OwnerSessionID() lifecyclestate.OwnerSessionID {
	owner.mu.Lock()
	defer owner.mu.Unlock()
	return owner.session
}

func (owner *archiveOwnerStub) CheckHeld(context.Context) error {
	owner.mu.Lock()
	defer owner.mu.Unlock()
	owner.checks++
	if !owner.held {
		return errors.New("owner not held")
	}
	return nil
}

type archiveFaultStub struct{ point ArchiveFault }

func (fault archiveFaultStub) FailArchiveAt(point ArchiveFault) error {
	if point == fault.point {
		return errors.New("injected archive fault")
	}
	return nil
}

type archiveFaultFunc func(ArchiveFault) error

func (fault archiveFaultFunc) FailArchiveAt(point ArchiveFault) error { return fault(point) }

func TestFixedStoreV2ArchivePrepareVerifyActivateKnownAnswer(t *testing.T) {
	path, store, owner := newEligibleFixedStoreV2(t)
	before := mustReadFile(t, path)
	limits, err := archivestate.NewArchiveLimits(1, uint64(archivestate.MaxSupervisorArchiveBytes))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := store.PlanArchive(context.Background(), owner, limits)
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := store.PrepareArchive(context.Background(), owner, plan)
	if err != nil {
		t.Fatal(err)
	}
	firstBytes := prepared.SegmentBytes()
	firstBytes[0] ^= 0xff
	if bytes.Equal(firstBytes, prepared.SegmentBytes()) {
		t.Fatal("prepared archive exposed mutable segment bytes")
	}
	if got := mustReadFile(t, path); !bytes.Equal(got, before) {
		t.Fatal("prepare changed active snapshot")
	}
	if _, err := os.Stat(archiveRoot(path)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("prepare created archive root: %v", err)
	}
	verified, err := store.VerifyPreparedArchive(context.Background(), owner, prepared)
	if err != nil {
		t.Fatal(err)
	}
	if got := mustReadFile(t, path); !bytes.Equal(got, before) {
		t.Fatal("verify changed active snapshot")
	}
	reopened, err := store.ActivateArchive(context.Background(), owner, verified, nil)
	if err != nil {
		t.Fatal(err)
	}
	report, err := VerifyFixedFileStoreV2(path)
	if err != nil {
		t.Fatal(err)
	}
	if report.SnapshotGeneration != 2 || report.ArchiveGeneration != 2 || report.SegmentCount != 1 ||
		report.OrphanSegmentCount != 0 || report.HotCounts != (archivestate.ArchiveCounts{}) ||
		report.ArchivedCounts.Registrations != 1 || report.ArchivedCounts.Approvals != 1 ||
		report.ArchivedCounts.Attempts != 1 || report.ArchivedCounts.Lifecycles != 1 ||
		report.ArchivedCounts.Nonces != 1 || report.ArchivedCounts.Effects != 1 ||
		report.ArchivedCounts.Instances != 1 || report.ArchivedCounts.ApprovalReplay != 1 ||
		report.ArchivedCounts.AttemptReplay != 1 {
		t.Fatalf("unexpected activated report: %#v", report)
	}
	if reopened.active.View().CurrentCheckpoint.Kind != archivestate.ArchiveCheckpointActivation {
		t.Fatal("activation did not install activation checkpoint")
	}
	segmentPath := archiveSegmentPath(path, prepared.SegmentDigest())
	segmentInfo, err := os.Lstat(segmentPath)
	if err != nil || segmentInfo.Mode().Perm() != 0o400 {
		t.Fatalf("immutable segment mode: %v, %v", segmentInfo, err)
	}
	activeDigest := sha256.Sum256(mustReadFile(t, path))
	segmentDigest := sha256.Sum256(mustReadFile(t, segmentPath))
	semanticDigest := prepared.SegmentDigest()
	combinedIndexDigest := reopened.active.View().CombinedIndexDigest
	knownAnswers := map[string]struct {
		got  []byte
		want string
	}{
		"active-file":      {activeDigest[:], "91e1978fc55420000b9cfce526d9a840554cb203119488c2539007f81daea1ad"},
		"segment-file":     {segmentDigest[:], "f78aa7b19bd1ef2f3a09fc641adf787e741ec72e7f6dfb6183be3d877e6cddc9"},
		"segment-semantic": {semanticDigest[:], "a6c126be25f24c1f631280d34b64e55d8389b53632391ff250f4cf7e17fb41ea"},
		"checkpoint":       {report.CurrentCheckpoint.Digest[:], "8f71152e057d9c75bd20b50b5a73e03ed3c59e38be353a26b4afa29f933afccc"},
		"combined-index":   {combinedIndexDigest[:], "7b9e25c9e6cb877188e0b1e5373e5950fd1d5bacdbebf71bef18eb25854f635e"},
	}
	for name, answer := range knownAnswers {
		if got := hex.EncodeToString(answer.got); got != answer.want {
			t.Errorf("%s known answer = %s, want %s", name, got, answer.want)
		}
	}
	if owner.checks < 7 {
		t.Fatalf("owner checks = %d, want entry/commit/reopen checks", owner.checks)
	}
}

func newEligibleFixedStoreV2(t *testing.T) (string, *FixedFileStoreV2, *archiveOwnerStub) {
	return newEligibleFixedStoreV2At(t, 1_785_456_300)
}

func newEligibleFixedStoreV2At(t *testing.T, highWater uint64) (string, *FixedFileStoreV2, *archiveOwnerStub) {
	t.Helper()
	state, template := stateAndLifecycleRecord(t)
	state.TimeHighWaterUnixSeconds = v0candidate.UInt53(highWater)
	destroyConfirmed := lifecycleRecordForV2State(t, state, template, lifecyclestate.StateDestroyConfirmed)
	destroyedView := destroyConfirmed.View()
	destroyedView.State = lifecyclestate.StateDestroyed
	destroyedView.CleanupRequired = false
	destroyedView.LastReconciliation = lifecyclestate.ReconciliationAuthoritativelyAbsent
	destroyedView.AutomaticRecoveryCount = 1
	destroyedView.TerminalAt = lifecyclestate.OptionalUnixSeconds{Present: true, Value: destroyedView.LastTransitionAt}
	destroyed, err := lifecyclestate.NewRecord(destroyedView, state.TimeHighWaterUnixSeconds)
	if err != nil {
		t.Fatalf("construct retained-identity destroyed lifecycle: %v", err)
	}
	path := filepath.Join(t.TempDir(), "supervisor-state.json")
	writeV1Envelope(t, path, encodedEnvelopeV1(state, []lifecyclestate.Record{destroyed}))
	store := mustMigrateV2(t, path)
	return path, store, archiveOwnerForTest(t)
}

func archiveOwnerForTest(t *testing.T) *archiveOwnerStub {
	t.Helper()
	domain, err := lifecyclestate.NewDomainIdentifier(lifecyclestate.DomainOwnerSessionID, bytes.Repeat([]byte{0xf3}, 16))
	if err != nil {
		t.Fatal(err)
	}
	session, err := lifecyclestate.NewOwnerSessionID(domain)
	if err != nil {
		t.Fatal(err)
	}
	return &archiveOwnerStub{held: true, session: session}
}

func mustPreparedArchive(t *testing.T, store *FixedFileStoreV2, owner ArchiveOwner) VerifiedFixedStoreV2Archive {
	t.Helper()
	limits, err := archivestate.NewArchiveLimits(1, uint64(archivestate.MaxSupervisorArchiveBytes))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := store.PlanArchive(context.Background(), owner, limits)
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
	return verified
}

func mutateSegmentJSON(t *testing.T, data []byte, mutate func(map[string]any)) []byte {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	mutate(value)
	result, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return append(result, '\n')
}

func TestFixedStoreV2ArchiveTestHelpersDeterministic(t *testing.T) {
	leftPath, leftStore, leftOwner := newEligibleFixedStoreV2(t)
	rightPath, rightStore, rightOwner := newEligibleFixedStoreV2(t)
	left := mustPreparedArchive(t, leftStore, leftOwner)
	right := mustPreparedArchive(t, rightStore, rightOwner)
	if !bytes.Equal(left.SegmentBytes(), right.SegmentBytes()) || left.SegmentDigest() != right.SegmentDigest() {
		t.Fatal("fixed F3 segment preparation is nondeterministic")
	}
	if reflect.DeepEqual(leftPath, rightPath) {
		t.Fatal("test roots unexpectedly alias")
	}
}

func TestArchiveActivationFaultMatrixPreservesExactlyOldOrNewWorld(t *testing.T) {
	tests := []struct {
		point         ArchiveFault
		indeterminate bool
		newWorld      bool
		orphanCount   uint16
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
	}
	for _, test := range tests {
		t.Run(string(test.point), func(t *testing.T) {
			path, store, owner := newEligibleFixedStoreV2(t)
			oldBytes := mustReadFile(t, path)
			verified := mustPreparedArchive(t, store, owner)
			_, err := store.ActivateArchive(context.Background(), owner, verified, archiveFaultStub{point: test.point})
			if err == nil {
				t.Fatal("injected archive fault succeeded")
			}
			if errors.Is(err, ErrArchiveOutcomeIndeterminate) != test.indeterminate {
				t.Fatalf("fault result = %v, indeterminate want %v", err, test.indeterminate)
			}
			reopened, openErr := OpenFixedFileStoreV2(path)
			if openErr != nil {
				t.Fatalf("reopen complete world: %v", openErr)
			}
			report, reportErr := VerifyFixedFileStoreV2(path)
			if reportErr != nil {
				t.Fatal(reportErr)
			}
			if test.newWorld {
				if report.SnapshotGeneration != 2 || report.SegmentCount != 1 || report.HotCounts.Registrations != 0 {
					t.Fatalf("did not reopen complete successor: %#v", report)
				}
			} else {
				if report.SnapshotGeneration != 1 || report.SegmentCount != 0 || report.HotCounts.Registrations != 1 ||
					!bytes.Equal(mustReadFile(t, path), oldBytes) {
					t.Fatalf("did not preserve complete predecessor: %#v", report)
				}
			}
			if report.OrphanSegmentCount != test.orphanCount || reopened.orphans != test.orphanCount {
				t.Fatalf("orphan count = %d/%d, want %d", report.OrphanSegmentCount, reopened.orphans, test.orphanCount)
			}
		})
	}
}

// TestArchiveActivationRetriesReuseExistingSegmentAfterPublishFault exercises
// the idempotent retry path documented in
// docs/SUPERVISOR_ARCHIVE_F3_ACTIVATION_RESULT.md: after an indeterminate
// fault at FaultArchiveAfterSegmentPublish/FaultArchiveAfterSegmentDirSync,
// the segment file is already durably renamed into place even though
// activation never completed. The documented recovery is a fresh
// Plan/Prepare/Verify/Activate cycle on a reopened store; publishArchiveSegment
// must recognize the byte-identical existing segment and reuse it rather than
// fail or rewrite it. Previously this recovery path was only proven by
// code-reading, not by an executed retry (see issue #129, and the identical
// pattern already fixed for owned-startup shutdown in issue #95).
func TestArchiveActivationRetriesReuseExistingSegmentAfterPublishFault(t *testing.T) {
	for _, point := range []ArchiveFault{FaultArchiveAfterSegmentPublish, FaultArchiveAfterSegmentDirSync} {
		t.Run(string(point), func(t *testing.T) {
			path, store, owner := newEligibleFixedStoreV2(t)
			verified := mustPreparedArchive(t, store, owner)
			segmentDigest := verified.SegmentDigest()

			if _, err := store.ActivateArchive(context.Background(), owner, verified, archiveFaultStub{point: point}); !errors.Is(err, ErrArchiveOutcomeIndeterminate) {
				t.Fatalf("injected fault at %s = %v, want ErrArchiveOutcomeIndeterminate", point, err)
			}
			faultedReport, err := VerifyFixedFileStoreV2(path)
			if err != nil {
				t.Fatal(err)
			}
			if faultedReport.SnapshotGeneration != 1 || faultedReport.SegmentCount != 0 ||
				faultedReport.OrphanSegmentCount != 1 {
				t.Fatalf("post-fault report = %#v, want the old world intact with one orphaned segment", faultedReport)
			}

			segmentPath := archiveSegmentPath(path, segmentDigest)
			orphanBytes := mustReadFile(t, segmentPath)
			orphanInfo, err := os.Lstat(segmentPath)
			if err != nil {
				t.Fatal(err)
			}

			reopened, err := OpenFixedFileStoreV2(path)
			if err != nil {
				t.Fatalf("reopen after fault: %v", err)
			}
			limits, err := archivestate.NewArchiveLimits(1, uint64(archivestate.MaxSupervisorArchiveBytes))
			if err != nil {
				t.Fatal(err)
			}
			retryPlan, err := reopened.PlanArchive(context.Background(), owner, limits)
			if err != nil {
				t.Fatalf("retry plan: %v", err)
			}
			retryPrepared, err := reopened.PrepareArchive(context.Background(), owner, retryPlan)
			if err != nil {
				t.Fatalf("retry prepare: %v", err)
			}
			if retryPrepared.SegmentDigest() != segmentDigest {
				t.Fatalf("retry segment digest = %x, want the orphaned segment's %x", retryPrepared.SegmentDigest(), segmentDigest)
			}
			retryVerified, err := reopened.VerifyPreparedArchive(context.Background(), owner, retryPrepared)
			if err != nil {
				t.Fatalf("retry verify: %v", err)
			}
			activated, err := reopened.ActivateArchive(context.Background(), owner, retryVerified, nil)
			if err != nil {
				t.Fatalf("retry activate did not reuse the existing segment: %v", err)
			}

			finalReport, err := VerifyFixedFileStoreV2(path)
			if err != nil {
				t.Fatal(err)
			}
			if finalReport.SnapshotGeneration != 2 || finalReport.ArchiveGeneration != 2 ||
				finalReport.SegmentCount != 1 || finalReport.OrphanSegmentCount != 0 {
				t.Fatalf("retry did not reach a clean single-segment activated world: %#v", finalReport)
			}
			if activated.active.View().CurrentCheckpoint.Kind != archivestate.ArchiveCheckpointActivation {
				t.Fatal("retry activation did not install activation checkpoint")
			}

			reusedInfo, err := os.Lstat(segmentPath)
			if err != nil {
				t.Fatal(err)
			}
			if !reusedInfo.ModTime().Equal(orphanInfo.ModTime()) || reusedInfo.Size() != orphanInfo.Size() {
				t.Fatalf("segment was rewritten instead of reused: mtime %v -> %v, size %d -> %d",
					orphanInfo.ModTime(), reusedInfo.ModTime(), orphanInfo.Size(), reusedInfo.Size())
			}
			if !bytes.Equal(mustReadFile(t, segmentPath), orphanBytes) {
				t.Fatal("reused segment bytes changed")
			}
		})
	}
}

func TestArchiveActivationDuplicateOwnerAndStaleConcurrencyRefuse(t *testing.T) {
	path, store, owner := newEligibleFixedStoreV2(t)
	verified := mustPreparedArchive(t, store, owner)
	owner.held = false
	if _, err := store.ActivateArchive(context.Background(), owner, verified, nil); !errors.Is(err, ErrArchiveOwnerRequired) {
		t.Fatalf("unheld owner activation = %v", err)
	}
	if _, err := os.Stat(archiveRoot(path)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("owner refusal wrote archive bytes: %v", err)
	}
	owner.held = true

	const contenders = 8
	results := make(chan error, contenders)
	for range contenders {
		go func() {
			_, err := store.ActivateArchive(context.Background(), owner, verified, nil)
			results <- err
		}()
	}
	successes := 0
	stale := 0
	for range contenders {
		err := <-results
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrArchiveStaleTransaction):
			stale++
		default:
			t.Fatalf("concurrent activation result = %v", err)
		}
	}
	if successes != 1 || stale != contenders-1 {
		t.Fatalf("concurrent activation successes/stale = %d/%d", successes, stale)
	}
	if report, err := VerifyFixedFileStoreV2(path); err != nil || report.SegmentCount != 1 || report.SnapshotGeneration != 2 {
		t.Fatalf("concurrent successor = %#v, %v", report, err)
	}
}

// TestFixedStoreV2ArchiveEntryAndSnapshotRaceOrdinaryMutationCommit is a
// regression test for #269: archiveEntry (used by PlanArchive and every other
// archive/backup entry point) and snapshotV2 used to read store.active,
// store.state, store.lifecycles, and store.genesis without holding store.mu,
// while commitV2WorldLocked writes those same fields under that mutex. Racing
// an archive read against an ordinary mutation commit on the same store
// reliably tripped `go test -race` before the fix; it must pass cleanly after
// it, with both reads now holding store.mu for their access to shared state.
func TestFixedStoreV2ArchiveEntryAndSnapshotRaceOrdinaryMutationCommit(t *testing.T) {
	_, store, owner := newEligibleFixedStoreV2(t)
	limits, err := archivestate.NewArchiveLimits(1, uint64(archivestate.MaxSupervisorArchiveBytes))
	if err != nil {
		t.Fatal(err)
	}

	// The mutation commit's write to store.active/state/lifecycles/genesis
	// (commitV2WorldLocked) sits behind a comparatively slow fsync/rename
	// sequence, so the actual unsynchronized write instant is a narrow window.
	// Spin the archive/snapshot reads in a tight loop for the whole duration of
	// a bounded sequence of ordinary mutation commits, rather than firing a
	// fixed number of one-shot calls, so the reads keep sampling continuously
	// and reliably land inside that window.
	const commits = 40
	stopReaders := make(chan struct{})
	var group sync.WaitGroup
	group.Add(3)

	go func() {
		defer group.Done()
		for {
			select {
			case <-stopReaders:
				return
			default:
			}
			// archiveEntry's unlocked read of store.active raced here.
			if _, err := store.PlanArchive(context.Background(), owner, limits); err != nil &&
				!errors.Is(err, ErrArchiveStaleTransaction) {
				t.Errorf("concurrent PlanArchive: %v", err)
				return
			}
		}
	}()
	go func() {
		defer group.Done()
		for {
			select {
			case <-stopReaders:
				return
			default:
			}
			// snapshotV2's unlocked read of store.state/lifecycles/active/genesis
			// raced here.
			if _, err := store.snapshotV2(context.Background()); err != nil {
				t.Errorf("concurrent snapshotV2: %v", err)
				return
			}
		}
	}()
	go func() {
		defer group.Done()
		defer close(stopReaders)
		for i := range commits {
			// An ordinary mutation commit, matching commitV2WorldLocked's write
			// path that both reads above raced against.
			observation := v0candidate.UInt53(1_785_456_300 + 1000 + uint64(i))
			if err := store.persistTimeHighWater(context.Background(), observation); err != nil {
				t.Errorf("concurrent persistTimeHighWater: %v", err)
				return
			}
		}
	}()
	group.Wait()
}

func TestArchiveOwnerRecheckedBeforePublicationActivationAndReopen(t *testing.T) {
	tests := []struct {
		name        string
		point       ArchiveFault
		orphanCount uint16
	}{
		{"before-publication", FaultArchiveBeforeSegmentPublish, 0},
		{"before-activation", FaultArchiveBeforeActivation, 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, store, owner := newEligibleFixedStoreV2(t)
			verified := mustPreparedArchive(t, store, owner)
			fault := archiveFaultFunc(func(point ArchiveFault) error {
				if point == test.point {
					owner.mu.Lock()
					owner.held = false
					owner.mu.Unlock()
				}
				return nil
			})
			if _, err := store.ActivateArchive(context.Background(), owner, verified, fault); err == nil ||
				!errors.Is(err, ErrArchiveOwnerRequired) {
				t.Fatalf("lost owner result = %v", err)
			}
			owner.mu.Lock()
			owner.held = true
			owner.mu.Unlock()
			report, err := VerifyFixedFileStoreV2(path)
			if err != nil || report.SnapshotGeneration != 1 || report.HotCounts.Registrations != 1 ||
				report.OrphanSegmentCount != test.orphanCount {
				t.Fatalf("owner-loss predecessor = %#v, %v", report, err)
			}
		})
	}

	t.Run("before-reopen", func(t *testing.T) {
		path, store, owner := newEligibleFixedStoreV2(t)
		verified := mustPreparedArchive(t, store, owner)
		fault := archiveFaultFunc(func(point ArchiveFault) error {
			if point == FaultArchiveAfterActiveDirSync {
				owner.mu.Lock()
				owner.held = false
				owner.mu.Unlock()
			}
			return nil
		})
		if _, err := store.ActivateArchive(context.Background(), owner, verified, fault); err == nil ||
			!errors.Is(err, ErrArchiveOwnerRequired) {
			t.Fatalf("pre-reopen owner loss = %v", err)
		}
		owner.mu.Lock()
		owner.held = true
		owner.mu.Unlock()
		if report, err := VerifyFixedFileStoreV2(path); err != nil || report.SnapshotGeneration != 2 || report.SegmentCount != 1 {
			t.Fatalf("pre-reopen owner-loss successor = %#v, %v", report, err)
		}
	})
}

func TestArchiveRefusesMissingOrActiveLifecycleWithoutInventingHistory(t *testing.T) {
	tests := []struct {
		name       string
		lifecycles func(t *testing.T, state installationState, template lifecyclestate.Record) []lifecyclestate.Record
	}{
		{name: "missing", lifecycles: func(_ *testing.T, _ installationState, _ lifecyclestate.Record) []lifecyclestate.Record { return nil }},
		{name: "active", lifecycles: func(t *testing.T, state installationState, template lifecyclestate.Record) []lifecyclestate.Record {
			return []lifecyclestate.Record{lifecycleRecordForV2State(t, state, template, lifecyclestate.StateObserved)}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state, template := stateAndLifecycleRecord(t)
			state.TimeHighWaterUnixSeconds = 1_785_456_300
			path := filepath.Join(t.TempDir(), "supervisor-state.json")
			writeV1Envelope(t, path, encodedEnvelopeV1(state, test.lifecycles(t, state, template)))
			store := mustMigrateV2(t, path)
			before := mustReadFile(t, path)
			limits, err := archivestate.NewArchiveLimits(1, uint64(archivestate.MaxSupervisorArchiveBytes))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := store.PlanArchive(context.Background(), archiveOwnerForTest(t), limits); err == nil {
				t.Fatal("incomplete lifecycle world produced an archive plan")
			}
			if !bytes.Equal(before, mustReadFile(t, path)) {
				t.Fatal("incomplete lifecycle refusal rewrote authority history")
			}
			if _, err := os.Stat(archiveRoot(path)); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("incomplete lifecycle refusal wrote archive bytes: %v", err)
			}
		})
	}
}

func TestArchiveSegmentCorruptionSubstitutionAndCrossLinksRefuseWithoutRewrite(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(t *testing.T, data []byte) []byte
	}{
		{name: "truncated", mutate: func(_ *testing.T, data []byte) []byte { return bytes.Clone(data[:len(data)/2]) }},
		{name: "trailing", mutate: func(_ *testing.T, data []byte) []byte { return append(bytes.Clone(data), []byte("{}\n")...) }},
		{name: "duplicate-name", mutate: func(_ *testing.T, data []byte) []byte {
			return append([]byte(`{"formatVersion":0,`), data[1:]...)
		}},
		{name: "unknown-field", mutate: func(t *testing.T, data []byte) []byte {
			return mutateSegmentJSON(t, data, func(value map[string]any) { value["unknownF3Field"] = true })
		}},
		{name: "unsupported-version", mutate: mutateArchiveSegmentDisk(func(disk *archiveSegmentDiskV0) {
			version := uint64(1)
			disk.FormatVersion = &version
		})},
		{name: "count-mismatch", mutate: mutateArchiveSegmentDisk(func(disk *archiveSegmentDiskV0) { disk.Counts.Attempts++ })},
		{name: "count-cap-plus-one", mutate: mutateArchiveSegmentDisk(func(disk *archiveSegmentDiskV0) {
			disk.Counts.Approvals = archivestate.MaxArchiveApprovalsPerSegment + 1
		})},
		{name: "duplicate-and-order", mutate: mutateArchiveSegmentDisk(func(disk *archiveSegmentDiskV0) {
			disk.Cohorts = append(disk.Cohorts, disk.Cohorts[0])
		})},
		{name: "attempt-cross-link", mutate: mutateArchiveSegmentDisk(func(disk *archiveSegmentDiskV0) {
			disk.Cohorts[0].Attempts[0].RegistrationID[0] ^= 0xff
		})},
		{name: "cohort-digest", mutate: mutateArchiveSegmentDisk(func(disk *archiveSegmentDiskV0) {
			disk.Cohorts[0].CohortDigest[0] ^= 0xff
		})},
		{name: "wrong-index-domain", mutate: mutateArchiveSegmentDisk(func(disk *archiveSegmentDiskV0) {
			disk.DerivedIndexes.Scope = archivestate.ArchiveIndexScopeRetainedGlobal
		})},
		{name: "effect-tombstone-omission", mutate: mutateArchiveSegmentDisk(func(disk *archiveSegmentDiskV0) {
			disk.DerivedIndexes.Effects = []effectIndexEntryDisk{}
		})},
		{name: "semantic-digest", mutate: mutateArchiveSegmentDisk(func(disk *archiveSegmentDiskV0) {
			disk.SegmentDigest = "00532c64181191d67d6d8b1f5dd4423277f2dcfa228bc2176e8580880c39fbdc"
		})},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, store, owner := newEligibleFixedStoreV2(t)
			verified := mustPreparedArchive(t, store, owner)
			if _, err := store.ActivateArchive(context.Background(), owner, verified, nil); err != nil {
				t.Fatal(err)
			}
			activeBefore := mustReadFile(t, path)
			segmentPath := archiveSegmentPath(path, verified.SegmentDigest())
			corrupt := test.mutate(t, mustReadFile(t, segmentPath))
			replaceArchiveFixtureFile(t, segmentPath, corrupt)
			corruptBefore := mustReadFile(t, segmentPath)
			if _, err := OpenFixedFileStoreV2(path); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
				t.Fatalf("corrupt archive open = %v", err)
			}
			if !bytes.Equal(activeBefore, mustReadFile(t, path)) || !bytes.Equal(corruptBefore, mustReadFile(t, segmentPath)) {
				t.Fatal("corrupt archive refusal rewrote retained bytes")
			}
		})
	}
}

func TestReferencedArchiveMissingAndValidSegmentSubstitutionNeverFallsBack(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		path, store, owner := newEligibleFixedStoreV2(t)
		verified := mustPreparedArchive(t, store, owner)
		if _, err := store.ActivateArchive(context.Background(), owner, verified, nil); err != nil {
			t.Fatal(err)
		}
		activeBefore := mustReadFile(t, path)
		if err := os.Remove(archiveSegmentPath(path, verified.SegmentDigest())); err != nil {
			t.Fatal(err)
		}
		if _, err := OpenFixedFileStoreV2(path); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
			t.Fatalf("missing segment open = %v", err)
		}
		if !bytes.Equal(activeBefore, mustReadFile(t, path)) {
			t.Fatal("missing segment refusal fell back or rewrote active state")
		}
	})

	t.Run("valid-substitution", func(t *testing.T) {
		path, store, owner := newEligibleFixedStoreV2At(t, 1_785_456_300)
		verified := mustPreparedArchive(t, store, owner)
		if _, err := store.ActivateArchive(context.Background(), owner, verified, nil); err != nil {
			t.Fatal(err)
		}
		_, otherStore, otherOwner := newEligibleFixedStoreV2At(t, 1_785_456_301)
		other := mustPreparedArchive(t, otherStore, otherOwner)
		if other.SegmentDigest() == verified.SegmentDigest() {
			t.Fatal("substitution fixture did not change segment identity")
		}
		segmentPath := archiveSegmentPath(path, verified.SegmentDigest())
		replaceArchiveFixtureFile(t, segmentPath, other.SegmentBytes())
		if _, err := OpenFixedFileStoreV2(path); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
			t.Fatalf("valid segment substitution open = %v", err)
		}
	})
}

func TestArchiveActivationProcessDeathOldOrNewWorld(t *testing.T) {
	tests := []struct {
		point       ArchiveFault
		newWorld    bool
		orphanCount uint16
	}{
		{FaultArchiveBeforeSegmentPublish, false, 1},
		{FaultArchiveAfterSegmentPublish, false, 1},
		{FaultArchiveBeforeActivation, false, 1},
		{FaultArchiveAfterActivation, true, 0},
	}
	for _, test := range tests {
		t.Run(string(test.point), func(t *testing.T) {
			path, _, _ := newEligibleFixedStoreV2(t)
			command := exec.Command(os.Args[0], "-test.run=^TestArchiveActivationProcessDeathHelper$") //nolint:gosec // G204: current bounded test binary only.
			command.Env = append(os.Environ(), "CAPSULE_F3_DEATH_PATH="+path, "CAPSULE_F3_DEATH_POINT="+string(test.point))
			err := command.Run()
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) || exitErr.ExitCode() != 77 {
				t.Fatalf("archive death helper = %v", err)
			}
			report, verifyErr := VerifyFixedFileStoreV2(path)
			if verifyErr != nil {
				t.Fatalf("reopen after process death: %v", verifyErr)
			}
			if test.newWorld {
				if report.SnapshotGeneration != 2 || report.SegmentCount != 1 || report.HotCounts.Registrations != 0 {
					t.Fatalf("death successor = %#v", report)
				}
			} else if report.SnapshotGeneration != 1 || report.SegmentCount != 0 || report.HotCounts.Registrations != 1 {
				t.Fatalf("death predecessor = %#v", report)
			}
			if report.OrphanSegmentCount != test.orphanCount {
				t.Fatalf("death orphan count = %d, want %d", report.OrphanSegmentCount, test.orphanCount)
			}
		})
	}
}

func TestArchiveActivationProcessDeathHelper(t *testing.T) {
	path := os.Getenv("CAPSULE_F3_DEATH_PATH")
	point := ArchiveFault(os.Getenv("CAPSULE_F3_DEATH_POINT"))
	if path == "" || point == "" {
		t.Skip("archive process-death helper")
	}
	store, err := OpenFixedFileStoreV2(path)
	if err != nil {
		os.Exit(79)
	}
	owner := archiveOwnerForTest(t)
	verified := mustPreparedArchive(t, store, owner)
	_, _ = store.ActivateArchive(context.Background(), owner, verified, exitArchiveFault{point: point})
	os.Exit(78)
}

type exitArchiveFault struct{ point ArchiveFault }

func (fault exitArchiveFault) FailArchiveAt(point ArchiveFault) error {
	if point == fault.point {
		os.Exit(77)
	}
	return nil
}

func mutateArchiveSegmentDisk(mutate func(*archiveSegmentDiskV0)) func(*testing.T, []byte) []byte {
	return func(t *testing.T, data []byte) []byte {
		t.Helper()
		var disk archiveSegmentDiskV0
		if err := decodeOneClosedJSON(data, &disk); err != nil {
			t.Fatal(err)
		}
		mutate(&disk)
		encoded, err := json.Marshal(disk)
		if err != nil {
			t.Fatal(err)
		}
		return append(encoded, '\n')
	}
}

func replaceArchiveFixtureFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o400); err != nil {
		t.Fatal(err)
	}
}
