package registrationstate

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/archivestate"
)

type backupFaultStub struct{ point BackupFault }

func (fault backupFaultStub) FailBackupAt(point BackupFault) error {
	if point == fault.point {
		return errors.New("injected backup fault")
	}
	return nil
}

type backupFaultFunc func(BackupFault) error

func (fault backupFaultFunc) FailBackupAt(point BackupFault) error { return fault(point) }

type orphanFaultStub struct{ point OrphanFault }

func (fault orphanFaultStub) FailOrphanAt(point OrphanFault) error {
	if point == fault.point {
		return errors.New("injected orphan fault")
	}
	return nil
}

func newActivatedF5Store(t *testing.T) (string, *FixedFileStoreV2, *archiveOwnerStub) {
	t.Helper()
	path, store, owner := newEligibleFixedStoreV2(t)
	activated, err := store.ActivateArchive(context.Background(), owner, mustPreparedArchive(t, store, owner), nil)
	if err != nil {
		t.Fatal(err)
	}
	return path, activated, owner
}

func newEmptyBackupRoot(t *testing.T) BackupRoot {
	t.Helper()
	path := filepath.Join(t.TempDir(), "backup")
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatal(err)
	}
	root, err := OpenBackupRoot(path)
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func TestCoherentBackupIncludesAndVerifiesExactArchiveSet(t *testing.T) {
	path, store, owner := newActivatedF5Store(t)
	activeBefore := mustReadFile(t, path)
	segmentPath := archiveSegmentPath(path, store.segments[0].Segment.Digest())
	segmentBefore := mustReadFile(t, segmentPath)
	destination := newEmptyBackupRoot(t)

	manifest, err := store.CreateCoherentBackup(context.Background(), owner, destination, nil)
	if err != nil {
		t.Fatal(err)
	}
	verified, report, err := VerifyCoherentBackup(context.Background(), destination)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.View() != verified.View() || report.Classification != ArchiveVerificationVerified ||
		report.SegmentCount != 1 || report.SnapshotGeneration != store.active.View().SnapshotGeneration ||
		report.CheckpointDigest != store.active.View().CurrentCheckpoint.Digest {
		t.Fatalf("backup manifest/report mismatch: %#v %#v", manifest.View(), report)
	}
	knownDigest := manifest.View().ManifestDigest
	if got, want := hex.EncodeToString(knownDigest[:]), "e9cee67bba449f9ebd2397f5e6c1c2e786201f7deee65d4e81643e4ef827ec1a"; got != want {
		t.Fatalf("coherent backup manifest known answer = %s", got)
	}
	if !bytes.Equal(activeBefore, mustReadFile(t, filepath.Join(destination.path, backupActiveName))) ||
		!bytes.Equal(segmentBefore, mustReadFile(t, filepath.Join(destination.path, backupArchiveDirName, filepath.Base(segmentPath)))) {
		t.Fatal("backup did not retain exact active/segment bytes")
	}
	if !bytes.Equal(activeBefore, mustReadFile(t, path)) || !bytes.Equal(segmentBefore, mustReadFile(t, segmentPath)) {
		t.Fatal("backup changed live predecessor")
	}
	storeRoot, err := OpenStoreRoot(path)
	if err != nil {
		t.Fatal(err)
	}
	offline, err := VerifyArchiveSet(context.Background(), storeRoot, VerificationFull)
	if err != nil || offline.Classification != ArchiveVerificationVerified || offline.TotalCounts != report.TotalCounts {
		t.Fatalf("offline report = %#v, %v", offline, err)
	}

	rollback, err := store.VerifyRestoreAdmission(context.Background(), owner, destination, nil)
	if !errors.Is(err, ErrRestoreRollbackUncertain) || rollback.Eligible ||
		rollback.Classification != ArchiveVerificationRollbackUncertain {
		t.Fatalf("no-anchor restore = %#v, %v", rollback, err)
	}
	anchor := LatestCheckpointFromManifest(manifest)
	eligible, err := store.VerifyRestoreAdmission(context.Background(), owner, destination, &anchor)
	if err != nil || !eligible.Eligible || eligible.Classification != ArchiveVerificationEligibleFutureRestore {
		t.Fatalf("anchored restore = %#v, %v", eligible, err)
	}
	if !bytes.Equal(activeBefore, mustReadFile(t, path)) {
		t.Fatal("restore admission replaced live predecessor")
	}
}

func TestCoherentBackupFaultMatrixNeverCreatesHybridCompleteBackup(t *testing.T) {
	points := []BackupFault{
		FaultBackupAfterPrepare,
		FaultBackupAfterActiveTempCreate, FaultBackupAfterActiveCopy, FaultBackupAfterActiveSync,
		FaultBackupAfterActiveClose, FaultBackupAfterActiveReopen,
		FaultBackupAfterActiveRename, FaultBackupAfterActiveDirectorySync,
		FaultBackupAfterSegmentTempCreate, FaultBackupAfterSegmentCopy, FaultBackupAfterSegmentSync,
		FaultBackupAfterSegmentClose, FaultBackupAfterSegmentReopen,
		FaultBackupAfterSegmentRename, FaultBackupAfterSegmentDirectorySync,
		FaultBackupAfterCheckpointValidation,
		FaultBackupAfterManifestTempCreate, FaultBackupAfterManifestWrite, FaultBackupAfterManifestSync,
		FaultBackupAfterManifestClose, FaultBackupAfterManifestReopen,
		FaultBackupAfterManifestRename, FaultBackupAfterRootDirectorySync, FaultBackupAfterReopen,
	}
	for _, point := range points {
		t.Run(string(point), func(t *testing.T) {
			path, store, owner := newActivatedF5Store(t)
			liveBefore := mustReadFile(t, path)
			destination := newEmptyBackupRoot(t)
			_, createErr := store.CreateCoherentBackup(context.Background(), owner, destination, backupFaultStub{point: point})
			if createErr == nil {
				t.Fatal("injected backup fault succeeded")
			}
			_, report, verifyErr := VerifyCoherentBackup(context.Background(), destination)
			completePoint := point == FaultBackupAfterManifestRename || point == FaultBackupAfterRootDirectorySync || point == FaultBackupAfterReopen
			if completePoint {
				if verifyErr != nil || report.Classification != ArchiveVerificationVerified {
					t.Fatalf("post-manifest fault did not leave complete successor: %#v, %v", report, verifyErr)
				}
			} else if verifyErr == nil || report.Classification == ArchiveVerificationVerified {
				t.Fatalf("pre-manifest fault created complete backup: %#v, %v", report, verifyErr)
			}
			if !bytes.Equal(liveBefore, mustReadFile(t, path)) {
				t.Fatal("backup fault changed live predecessor")
			}
		})
	}
}

func TestCoherentBackupMissingExtraSubstitutedAndCapRefusePreserve(t *testing.T) {
	for _, mutate := range []struct {
		name string
		fn   func(*testing.T, BackupRoot)
	}{
		{"missing-segment", func(t *testing.T, root BackupRoot) {
			entries, _ := os.ReadDir(filepath.Join(root.path, backupArchiveDirName))
			if err := os.Remove(filepath.Join(root.path, backupArchiveDirName, entries[0].Name())); err != nil {
				t.Fatal(err)
			}
		}},
		{"extra", func(t *testing.T, root BackupRoot) {
			if err := os.WriteFile(filepath.Join(root.path, "unknown"), []byte("evidence"), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
		{"substituted-active", func(t *testing.T, root BackupRoot) {
			path := filepath.Join(root.path, backupActiveName)
			data := mustReadFile(t, path)
			data[len(data)/2] ^= 1
			if err := os.WriteFile(path, data, 0o600); err != nil {
				t.Fatal(err)
			}
		}},
		{"manifest-cap-plus-one", func(t *testing.T, root BackupRoot) {
			path := filepath.Join(root.path, backupManifestName)
			if err := os.Chmod(path, 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.Truncate(path, maxBackupManifest+1); err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(path, 0o400); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(mutate.name, func(t *testing.T) {
			_, store, owner := newActivatedF5Store(t)
			root := newEmptyBackupRoot(t)
			if _, err := store.CreateCoherentBackup(context.Background(), owner, root, nil); err != nil {
				t.Fatal(err)
			}
			mutate.fn(t, root)
			before := backupInventoryBytes(t, root.path)
			if _, _, err := VerifyCoherentBackup(context.Background(), root); err == nil {
				t.Fatal("invalid backup verified")
			}
			if after := backupInventoryBytes(t, root.path); !bytes.Equal(before, after) {
				t.Fatal("backup refusal rewrote evidence")
			}
		})
	}
}

func TestCoherentBackupStaleOwnerAndMixedCheckpointRefuse(t *testing.T) {
	t.Run("owner-at-entry", func(t *testing.T) {
		_, store, owner := newActivatedF5Store(t)
		owner.held = false
		root := newEmptyBackupRoot(t)
		if _, err := store.CreateCoherentBackup(context.Background(), owner, root, nil); err == nil {
			t.Fatal("stale owner created backup")
		}
		entries, _ := os.ReadDir(root.path)
		if len(entries) != 0 {
			t.Fatal("stale owner changed destination")
		}
	})
	t.Run("owner-before-manifest", func(t *testing.T) {
		path, store, owner := newActivatedF5Store(t)
		before := mustReadFile(t, path)
		root := newEmptyBackupRoot(t)
		faults := backupFaultFunc(func(point BackupFault) error {
			if point == FaultBackupAfterSegmentDirectorySync {
				owner.mu.Lock()
				owner.held = false
				owner.mu.Unlock()
			}
			return nil
		})
		if _, err := store.CreateCoherentBackup(context.Background(), owner, root, faults); err == nil {
			t.Fatal("lost owner completed backup")
		}
		if _, _, err := VerifyCoherentBackup(context.Background(), root); err == nil {
			t.Fatal("lost owner left complete backup")
		}
		if !bytes.Equal(before, mustReadFile(t, path)) {
			t.Fatal("lost owner changed predecessor")
		}
	})
	t.Run("mixed-manifest", func(t *testing.T) {
		_, store, owner := newActivatedF5Store(t)
		left := newEmptyBackupRoot(t)
		if _, err := store.CreateCoherentBackup(context.Background(), owner, left, nil); err != nil {
			t.Fatal(err)
		}
		_, newer, newerOwner, _ := newEligibleGrowthStoreV2(t, 2)
		newer = activateNextArchiveSegment(t, newer, newerOwner)
		newer = activateNextArchiveSegment(t, newer, newerOwner)
		right := newEmptyBackupRoot(t)
		if _, err := newer.CreateCoherentBackup(context.Background(), newerOwner, right, nil); err != nil {
			t.Fatal(err)
		}
		leftManifest := filepath.Join(left.path, backupManifestName)
		if err := os.Chmod(leftManifest, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(leftManifest, mustReadFile(t, filepath.Join(right.path, backupManifestName)), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(leftManifest, 0o400); err != nil {
			t.Fatal(err)
		}
		if _, report, err := VerifyCoherentBackup(context.Background(), left); err == nil ||
			report.Classification == ArchiveVerificationVerified {
			t.Fatalf("mixed backup verified: %#v, %v", report, err)
		}
	})
}

func TestCoherentBackupContendersPublishExactlyOneCompleteSet(t *testing.T) {
	_, store, owner := newActivatedF5Store(t)
	destination := newEmptyBackupRoot(t)
	const contenders = 8
	results := make(chan error, contenders)
	var group sync.WaitGroup
	for range contenders {
		group.Add(1)
		go func() {
			defer group.Done()
			_, err := store.CreateCoherentBackup(context.Background(), owner, destination, nil)
			results <- err
		}()
	}
	group.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
			continue
		}
		if !errors.Is(err, ErrBackupDestinationNotEmpty) {
			t.Fatalf("backup contender = %v", err)
		}
	}
	if successes != 1 {
		t.Fatalf("successful contenders = %d, want 1", successes)
	}
	if _, report, err := VerifyCoherentBackup(context.Background(), destination); err != nil ||
		report.Classification != ArchiveVerificationVerified {
		t.Fatalf("contended backup = %#v, %v", report, err)
	}
}

func backupInventoryBytes(t *testing.T, root string) []byte {
	t.Helper()
	var result []byte
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		result = append(result, []byte(path[len(root):])...)
		result = append(result, mustReadFile(t, path)...)
		return nil
	})
	return result
}

func TestKnownUnreferencedOrphanDeletionRequiresInventoryAndReferenceScan(t *testing.T) {
	path, store, owner := newEligibleFixedStoreV2(t)
	verified := mustPreparedArchive(t, store, owner)
	if _, err := store.ActivateArchive(context.Background(), owner, verified, archiveFaultStub{point: FaultArchiveAfterSegmentPublish}); !errors.Is(err, ErrArchiveOutcomeIndeterminate) {
		t.Fatalf("publish fault = %v", err)
	}
	reopened, err := OpenFixedFileStoreV2(path)
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := reopened.InventoryArchiveArtifacts(context.Background(), owner, nil)
	if err != nil || inventory.Report().KnownUnreferenced != 1 || len(inventory.KnownUnreferenced()) != 1 {
		t.Fatalf("orphan inventory = %#v, %v", inventory.Report(), err)
	}
	candidate := inventory.KnownUnreferenced()[0]
	orphanPath := archiveSegmentPath(path, candidate.Digest())
	orphanBefore := mustReadFile(t, orphanPath)
	if err := reopened.DeleteKnownUnreferencedOrphan(context.Background(), owner, inventory, KnownUnreferencedOrphan{}, nil); !errors.Is(err, ErrOrphanEvidencePreserved) {
		t.Fatalf("zero candidate deletion = %v", err)
	}
	if !bytes.Equal(orphanBefore, mustReadFile(t, orphanPath)) {
		t.Fatal("invalid deletion request changed orphan")
	}
	if err := reopened.DeleteKnownUnreferencedOrphan(context.Background(), owner, inventory, candidate, orphanFaultStub{point: FaultOrphanBeforeRemove}); err == nil {
		t.Fatal("pre-remove fault succeeded")
	}
	if !bytes.Equal(orphanBefore, mustReadFile(t, orphanPath)) {
		t.Fatal("pre-remove fault changed orphan")
	}
	if err := reopened.DeleteKnownUnreferencedOrphan(context.Background(), owner, inventory, candidate, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphanPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("known orphan still exists: %v", err)
	}
	final, err := OpenFixedFileStoreV2(path)
	if err != nil || final.orphans != 0 || final.active.View().CurrentCheckpoint != store.active.View().CurrentCheckpoint {
		t.Fatalf("predecessor after orphan removal = %#v, %v", final, err)
	}
}

func TestBackupReferencedOrphanIsReportedAndNeverOfferedForDeletion(t *testing.T) {
	path, predecessor, owner := newEligibleFixedStoreV2(t)
	predecessorBytes := mustReadFile(t, path)
	successor, err := predecessor.ActivateArchive(context.Background(), owner, mustPreparedArchive(t, predecessor, owner), nil)
	if err != nil {
		t.Fatal(err)
	}
	backup := newEmptyBackupRoot(t)
	manifest, err := successor.CreateCoherentBackup(context.Background(), owner, backup, nil)
	if err != nil {
		t.Fatal(err)
	}
	segmentPath := archiveSegmentPath(path, manifest.disk.Segments[0].SegmentDigest)
	segmentBefore := mustReadFile(t, segmentPath)
	if err := os.WriteFile(path, predecessorBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	reopened, err := OpenFixedFileStoreV2(path)
	if err != nil {
		t.Fatal(err)
	}
	inventory, err := reopened.InventoryArchiveArtifacts(context.Background(), owner, []BackupRoot{backup})
	if err != nil || inventory.Report().BackupReferenced != 1 ||
		inventory.Report().KnownUnreferenced != 0 || len(inventory.KnownUnreferenced()) != 0 {
		t.Fatalf("backup-protected inventory = %#v, %v", inventory.Report(), err)
	}
	if !bytes.Equal(segmentBefore, mustReadFile(t, segmentPath)) {
		t.Fatal("backup-referenced orphan changed")
	}
}

func TestKnownOrphanPostRemoveFaultsReopenOnlyCompletePredecessor(t *testing.T) {
	for _, point := range []OrphanFault{FaultOrphanAfterRemove, FaultOrphanAfterDirectorySync, FaultOrphanAfterReopen} {
		t.Run(string(point), func(t *testing.T) {
			path, store, owner := newEligibleFixedStoreV2(t)
			activeBefore := mustReadFile(t, path)
			verified := mustPreparedArchive(t, store, owner)
			if _, err := store.ActivateArchive(context.Background(), owner, verified, archiveFaultStub{point: FaultArchiveAfterSegmentPublish}); !errors.Is(err, ErrArchiveOutcomeIndeterminate) {
				t.Fatal(err)
			}
			reopened, err := OpenFixedFileStoreV2(path)
			if err != nil {
				t.Fatal(err)
			}
			inventory, err := reopened.InventoryArchiveArtifacts(context.Background(), owner, nil)
			if err != nil {
				t.Fatal(err)
			}
			candidate := inventory.KnownUnreferenced()[0]
			if err := reopened.DeleteKnownUnreferencedOrphan(context.Background(), owner, inventory, candidate, orphanFaultStub{point: point}); !errors.Is(err, ErrOrphanOutcomeIndeterminate) {
				t.Fatalf("post-remove fault = %v", err)
			}
			if !bytes.Equal(activeBefore, mustReadFile(t, path)) {
				t.Fatal("orphan response loss changed active predecessor")
			}
			if err := reopened.DeleteKnownUnreferencedOrphan(context.Background(), owner, inventory, candidate, nil); err != nil {
				t.Fatalf("exact recovery retry = %v", err)
			}
			fresh, err := OpenFixedFileStoreV2(path)
			if err != nil || fresh.orphans != 0 {
				t.Fatalf("post-remove reopen = %#v, %v", fresh, err)
			}
		})
	}
}

func TestUnknownOrCorruptOrphanRefusesAndPreservesEvidence(t *testing.T) {
	for _, name := range []string{"unknown.bin", "segment-not-a-digest.json"} {
		t.Run(name, func(t *testing.T) {
			path, store, owner := newEligibleFixedStoreV2(t)
			if _, err := ensureArchiveRoot(path); err != nil {
				t.Fatal(err)
			}
			evidencePath := filepath.Join(archiveRoot(path), name)
			if err := os.WriteFile(evidencePath, []byte("preserve-me"), 0o400); err != nil {
				t.Fatal(err)
			}
			before := mustReadFile(t, evidencePath)
			inventory, err := store.InventoryArchiveArtifacts(context.Background(), owner, nil)
			if !errors.Is(err, ErrOrphanEvidencePreserved) ||
				(inventory.Report().UnknownArtifacts == 0 && inventory.Report().CorruptArtifacts == 0) {
				t.Fatalf("unsafe inventory = %#v, %v", inventory.Report(), err)
			}
			if !bytes.Equal(before, mustReadFile(t, evidencePath)) || len(inventory.KnownUnreferenced()) != 0 {
				t.Fatal("unsafe artifact was changed or offered for deletion")
			}
		})
	}
}

func TestFailedFullVerificationSuppressesOtherwiseKnownOrphanDeletion(t *testing.T) {
	path, store, owner, _ := newEligibleGrowthStoreV2(t, 2)
	store = activateNextArchiveSegment(t, store, owner)
	verified := mustPreparedArchive(t, store, owner)
	if _, err := store.ActivateArchive(context.Background(), owner, verified, archiveFaultStub{point: FaultArchiveAfterSegmentPublish}); !errors.Is(err, ErrArchiveOutcomeIndeterminate) {
		t.Fatal(err)
	}
	orphanPath := archiveSegmentPath(path, verified.SegmentDigest())
	orphanBefore := mustReadFile(t, orphanPath)
	referencedPath := archiveSegmentPath(path, store.segments[0].Segment.Digest())
	referencedBytes := mustReadFile(t, referencedPath)
	referencedBytes[len(referencedBytes)/2] ^= 1
	if err := os.Chmod(referencedPath, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(referencedPath, referencedBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(referencedPath, 0o400); err != nil {
		t.Fatal(err)
	}
	inventory, err := store.InventoryArchiveArtifacts(context.Background(), owner, nil)
	if !errors.Is(err, ErrOrphanEvidencePreserved) ||
		inventory.Report().Classification != ArchiveVerificationCorrupt ||
		len(inventory.KnownUnreferenced()) != 0 {
		t.Fatalf("failed-verifier inventory = %#v, %v", inventory.Report(), err)
	}
	if !bytes.Equal(orphanBefore, mustReadFile(t, orphanPath)) {
		t.Fatal("failed full verification changed orphan evidence")
	}
}

func TestBackupProcessDeathLeavesOnlyIncompleteOrCompleteSet(t *testing.T) {
	if os.Getenv("CAPSULE_F5_BACKUP_DEATH") != "" {
		path := os.Getenv("CAPSULE_F5_STORE")
		destinationPath := os.Getenv("CAPSULE_F5_BACKUP")
		store, err := OpenFixedFileStoreV2(path)
		if err != nil {
			os.Exit(91)
		}
		root, err := OpenBackupRoot(destinationPath)
		if err != nil {
			os.Exit(92)
		}
		owner := archiveOwnerForTest(t)
		_, _ = store.CreateCoherentBackup(context.Background(), owner, root, exitBackupFault{point: BackupFault(os.Getenv("CAPSULE_F5_BACKUP_DEATH"))})
		os.Exit(93)
	}
	for _, point := range []BackupFault{FaultBackupAfterActiveRename, FaultBackupAfterSegmentRename, FaultBackupAfterManifestRename} {
		t.Run(string(point), func(t *testing.T) {
			path, _, _ := newActivatedF5Store(t)
			root := newEmptyBackupRoot(t)
			command := exec.Command(os.Args[0], "-test.run=^TestBackupProcessDeathLeavesOnlyIncompleteOrCompleteSet$") //nolint:gosec // current test binary and fixed literal argument.
			command.Env = append(os.Environ(), "CAPSULE_F5_BACKUP_DEATH="+string(point), "CAPSULE_F5_STORE="+path, "CAPSULE_F5_BACKUP="+root.path)
			if err := command.Run(); err == nil {
				t.Fatal("death child unexpectedly succeeded")
			}
			_, report, err := VerifyCoherentBackup(context.Background(), root)
			complete := point == FaultBackupAfterManifestRename
			if complete && (err != nil || report.Classification != ArchiveVerificationVerified) {
				t.Fatalf("post-manifest death backup = %#v, %v", report, err)
			}
			if !complete && err == nil {
				t.Fatal("pre-manifest death produced complete backup")
			}
		})
	}
}

type exitBackupFault struct{ point BackupFault }

func (fault exitBackupFault) FailBackupAt(point BackupFault) error {
	if point == fault.point {
		os.Exit(94)
	}
	return nil
}

func TestBackupExactSegmentCapacityAndRepeatedVerification(t *testing.T) {
	path, store, owner, _ := newEligibleGrowthStoreV2(t, 2)
	store = activateNextArchiveSegment(t, store, owner)
	store = activateNextArchiveSegment(t, store, owner)
	root := newEmptyBackupRoot(t)
	manifest, err := store.CreateCoherentBackup(context.Background(), owner, root, nil)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.View().SegmentCount != 2 {
		t.Fatalf("segment count = %d", manifest.View().SegmentCount)
	}
	for range 20 {
		repeated, report, verifyErr := VerifyCoherentBackup(context.Background(), root)
		if verifyErr != nil || repeated.View() != manifest.View() || report.SegmentCount != 2 {
			t.Fatalf("repeated verification = %#v %#v, %v", repeated.View(), report, verifyErr)
		}
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestCoherentBackupAcceptsExact64SegmentWorldAndPreservesCapPlusOneRefusal(t *testing.T) {
	path, store, owner, _ := newEligibleGrowthStoreV2(t, archivestate.MaxReferencedArchiveSegments+1)
	for range archivestate.MaxReferencedArchiveSegments {
		store = activateNextArchiveSegment(t, store, owner)
	}
	before := inventoryDigest(t, path)
	root := newEmptyBackupRoot(t)
	manifest, err := store.CreateCoherentBackup(context.Background(), owner, root, nil)
	if err != nil || manifest.View().SegmentCount != archivestate.MaxReferencedArchiveSegments {
		t.Fatalf("exact-cap backup = %#v, %v", manifest.View(), err)
	}
	verified, report, err := VerifyCoherentBackup(context.Background(), root)
	if err != nil || verified.View() != manifest.View() || report.SegmentCount != archivestate.MaxReferencedArchiveSegments {
		t.Fatalf("exact-cap verification = %#v %#v, %v", verified.View(), report, err)
	}
	if _, err := store.PlanArchive(context.Background(), owner, archiveOneCohortLimits(t)); err == nil ||
		!bytes.Contains([]byte(err.Error()), []byte("CAPACITY")) {
		t.Fatalf("segment cap-plus-one plan = %v", err)
	}
	if after := inventoryDigest(t, path); before != after {
		t.Fatal("backup or cap-plus-one refusal changed live history")
	}
}

func TestRestoreRejectsOlderFutureAndIncomparableAnchor(t *testing.T) {
	_, store, owner := newActivatedF5Store(t)
	root := newEmptyBackupRoot(t)
	manifest, err := store.CreateCoherentBackup(context.Background(), owner, root, nil)
	if err != nil {
		t.Fatal(err)
	}
	base := LatestCheckpointFromManifest(manifest)
	mutations := []func(*IndependentLatestCheckpoint){
		func(anchor *IndependentLatestCheckpoint) { anchor.SnapshotGeneration-- },
		func(anchor *IndependentLatestCheckpoint) { anchor.SnapshotGeneration++ },
		func(anchor *IndependentLatestCheckpoint) {
			anchor.CurrentCheckpoint = archivestate.NoArchiveCheckpointReference()
		},
	}
	for _, mutate := range mutations {
		anchor := base
		mutate(&anchor)
		report, admissionErr := store.VerifyRestoreAdmission(context.Background(), owner, root, &anchor)
		if !errors.Is(admissionErr, ErrRestoreRollbackUncertain) || report.Eligible {
			t.Fatalf("mismatched anchor admitted: %#v, %v", report, admissionErr)
		}
	}
}
