package registrationstate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"

	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	backupFormatV0       = uint64(0)
	backupManifestName   = "manifest.json"
	backupActiveName     = "state.json"
	backupArchiveDirName = "state.json.archive"
	maxBackupManifest    = int64(1 << 20)
)

// VerificationMode is deliberately closed to full verification. F5 selects
// no cached, incremental, or startup-fast verification mode.
type VerificationMode string

const VerificationFull VerificationMode = "full"

// ArchiveVerificationClassification is a bounded offline reporting
// vocabulary. It contains no stored text or filesystem path.
type ArchiveVerificationClassification string

const (
	ArchiveVerificationVerified              ArchiveVerificationClassification = "verified"
	ArchiveVerificationVerifiedWithOrphans   ArchiveVerificationClassification = "verified-with-known-unreferenced"
	ArchiveVerificationIncomplete            ArchiveVerificationClassification = "incomplete"
	ArchiveVerificationExtraArtifacts        ArchiveVerificationClassification = "extra-artifacts"
	ArchiveVerificationCorrupt               ArchiveVerificationClassification = "corrupt"
	ArchiveVerificationMixedGeneration       ArchiveVerificationClassification = "mixed-generation"
	ArchiveVerificationCrossInstallation     ArchiveVerificationClassification = "cross-installation"
	ArchiveVerificationRollbackUncertain     ArchiveVerificationClassification = "rollback-uncertain"
	ArchiveVerificationEligibleFutureRestore ArchiveVerificationClassification = "eligible-for-future-restore"
)

// StoreRoot and BackupRoot are internal trusted local capabilities. They are
// not public/IPC path contracts and expose no path accessor.
type StoreRoot struct{ path string }
type BackupRoot struct{ path string }

// OpenStoreRoot accepts only an existing installation-configured v2 store
// file. Full semantic validation remains the verifier's responsibility.
func OpenStoreRoot(path string) (StoreRoot, error) {
	if path == "" {
		return StoreRoot{}, errors.New("fixed supervisor archive store root is required")
	}
	if err := validateExistingStoreFile(path); err != nil {
		return StoreRoot{}, err
	}
	return StoreRoot{path: path}, nil
}

// OpenBackupRoot accepts an existing owner-only directory. Creation is left
// to installation code or the owned local test harness.
func OpenBackupRoot(path string) (BackupRoot, error) {
	if err := validateOwnedDirectory(path, 0o700); err != nil {
		return BackupRoot{}, err
	}
	return BackupRoot{path: path}, nil
}

type BackupFileDigest [32]byte
type BackupManifestDigest [32]byte

// backupHotSetDigestsDisk explicitly serializes the fifth F4B collection.
// archivestate.HotSetDigests omits it from the legacy F2/F3 JSON shape.
type backupHotSetDigestsDisk struct {
	Registrations    [32]byte                              `json:"registrations"`
	Approvals        [32]byte                              `json:"approvals"`
	Attempts         [32]byte                              `json:"attempts"`
	Lifecycles       [32]byte                              `json:"lifecycles"`
	EffectTombstones archivestate.EffectTombstoneSetDigest `json:"effectTombstones"`
}

type backupSegmentDisk struct {
	Ordinal                  archivestate.ArchiveOrdinal          `json:"ordinal"`
	FileName                 string                               `json:"fileName"`
	EncodedByteLength        uint64                               `json:"encodedByteLength"`
	FileDigest               BackupFileDigest                     `json:"fileDigest"`
	SegmentDigest            archivestate.ArchiveSegmentDigest    `json:"segmentDigest"`
	SourceSnapshotGeneration v0candidate.PositiveUInt53           `json:"sourceSnapshotGeneration"`
	ArchiveGeneration        archivestate.ArchiveGeneration       `json:"archiveGeneration"`
	PriorCheckpointDigest    archivestate.ArchiveCheckpointDigest `json:"priorCheckpointDigest"`
}

type backupManifestDisk struct {
	FormatVersion             *uint64                                 `json:"formatVersion"`
	StoreFormatVersion        uint64                                  `json:"storeFormatVersion"`
	ArchiveFormatVersion      uint64                                  `json:"archiveFormatVersion"`
	InstallationID            v0candidate.InstallationID              `json:"installationId"`
	SupervisorID              v0candidate.SupervisorID                `json:"supervisorId"`
	EpochSequence             v0candidate.UInt53                      `json:"epochSequence"`
	EpochDigest               v0candidate.TrustEpochDigest            `json:"epochDigest"`
	SnapshotGeneration        v0candidate.PositiveUInt53              `json:"snapshotGeneration"`
	ArchiveGeneration         archivestate.ArchiveGeneration          `json:"archiveGeneration"`
	DurableTimeHighWater      v0candidate.UInt53                      `json:"durableTimeHighWater"`
	PreviousCheckpoint        archivestate.ArchiveCheckpointReference `json:"previousCheckpoint"`
	CurrentCheckpoint         archivestate.ArchiveCheckpointReference `json:"currentCheckpoint"`
	MigrationGenesisDigest    archivestate.ArchiveCheckpointDigest    `json:"migrationGenesisDigest"`
	DescriptorSetDigest       archivestate.ArchiveDescriptorSetDigest `json:"descriptorSetDigest"`
	CombinedIndexDigest       archivestate.ArchiveCombinedIndexDigest `json:"combinedIndexDigest"`
	HotSetDigests             backupHotSetDigestsDisk                 `json:"hotSetDigests"`
	VisibleV1EffectSeedCount  uint64                                  `json:"visibleV1EffectSeedCount"`
	VisibleV1EffectSeedDigest archivestate.ArchiveEffectSeedDigest    `json:"visibleV1EffectSeedDigest"`
	HotCounts                 archivestate.ArchiveCounts              `json:"hotCounts"`
	ArchivedCounts            archivestate.ArchiveCounts              `json:"archivedCounts"`
	TotalCounts               archivestate.ArchiveCounts              `json:"totalCounts"`
	ActiveFileName            string                                  `json:"activeFileName"`
	ActiveEncodedByteLength   uint64                                  `json:"activeEncodedByteLength"`
	ActiveFileDigest          BackupFileDigest                        `json:"activeFileDigest"`
	Segments                  []backupSegmentDisk                     `json:"segments"`
	ManifestDigest            BackupManifestDigest                    `json:"manifestDigest"`
}

// BackupManifestView is a defensive, path-free projection of one verified
// coherent set. Segment file names are fixed digest-addressed basenames, not
// caller paths.
type BackupManifestView struct {
	FormatVersion             uint64
	InstallationID            v0candidate.InstallationID
	SupervisorID              v0candidate.SupervisorID
	EpochSequence             v0candidate.UInt53
	EpochDigest               v0candidate.TrustEpochDigest
	SnapshotGeneration        v0candidate.PositiveUInt53
	ArchiveGeneration         archivestate.ArchiveGeneration
	DurableTimeHighWater      v0candidate.UInt53
	PreviousCheckpoint        archivestate.ArchiveCheckpointReference
	CurrentCheckpoint         archivestate.ArchiveCheckpointReference
	VisibleV1EffectSeedCount  uint64
	VisibleV1EffectSeedDigest archivestate.ArchiveEffectSeedDigest
	HotCounts                 archivestate.ArchiveCounts
	ArchivedCounts            archivestate.ArchiveCounts
	TotalCounts               archivestate.ArchiveCounts
	SegmentCount              uint16
	ManifestDigest            BackupManifestDigest
}

type BackupManifest struct {
	disk backupManifestDisk
	root BackupRoot
}

func (manifest BackupManifest) View() BackupManifestView {
	disk := manifest.disk
	return BackupManifestView{
		FormatVersion: *disk.FormatVersion, InstallationID: disk.InstallationID,
		SupervisorID: disk.SupervisorID, EpochSequence: disk.EpochSequence, EpochDigest: disk.EpochDigest,
		SnapshotGeneration: disk.SnapshotGeneration, ArchiveGeneration: disk.ArchiveGeneration,
		DurableTimeHighWater: disk.DurableTimeHighWater, PreviousCheckpoint: disk.PreviousCheckpoint,
		CurrentCheckpoint: disk.CurrentCheckpoint, VisibleV1EffectSeedCount: disk.VisibleV1EffectSeedCount,
		VisibleV1EffectSeedDigest: disk.VisibleV1EffectSeedDigest, HotCounts: disk.HotCounts,
		ArchivedCounts: disk.ArchivedCounts, TotalCounts: disk.TotalCounts,
		SegmentCount: uint16(len(disk.Segments)), ManifestDigest: disk.ManifestDigest, //nolint:gosec // bounded by 64.
	}
}

type ArchiveVerificationReport struct {
	Version              uint64
	Classification       ArchiveVerificationClassification
	StoreFormatVersion   uint64
	ArchiveFormatVersion uint64
	InstallationID       v0candidate.InstallationID
	SupervisorID         v0candidate.SupervisorID
	EpochSequence        v0candidate.UInt53
	EpochDigest          v0candidate.TrustEpochDigest
	SnapshotGeneration   v0candidate.PositiveUInt53
	ArchiveGeneration    archivestate.ArchiveGeneration
	TimeHighWater        v0candidate.UInt53
	PreviousCheckpoint   archivestate.ArchiveCheckpointReference
	CheckpointDigest     archivestate.ArchiveCheckpointDigest
	ActiveCounts         archivestate.ArchiveCounts
	ArchivedCounts       archivestate.ArchiveCounts
	TotalCounts          archivestate.ArchiveCounts
	SegmentCount         uint16
	KnownUnreferenced    uint16
	BackupReferenced     uint16
	UnknownArtifacts     uint16
	CorruptArtifacts     uint16
	MixedGeneration      uint16
	CrossInstallation    uint16
}

type BackupFault string

const (
	FaultBackupAfterPrepare              BackupFault = "backup-after-prepare"
	FaultBackupAfterActiveTempCreate     BackupFault = "backup-after-active-temp-create"
	FaultBackupAfterActiveCopy           BackupFault = "backup-after-active-copy"
	FaultBackupAfterActiveSync           BackupFault = "backup-after-active-sync"
	FaultBackupAfterActiveClose          BackupFault = "backup-after-active-close"
	FaultBackupAfterActiveReopen         BackupFault = "backup-after-active-reopen"
	FaultBackupAfterActiveRename         BackupFault = "backup-after-active-rename"
	FaultBackupAfterActiveDirectorySync  BackupFault = "backup-after-active-directory-sync"
	FaultBackupAfterSegmentTempCreate    BackupFault = "backup-after-segment-temp-create"
	FaultBackupAfterSegmentCopy          BackupFault = "backup-after-segment-copy"
	FaultBackupAfterSegmentSync          BackupFault = "backup-after-segment-sync"
	FaultBackupAfterSegmentClose         BackupFault = "backup-after-segment-close"
	FaultBackupAfterSegmentReopen        BackupFault = "backup-after-segment-reopen"
	FaultBackupAfterSegmentRename        BackupFault = "backup-after-segment-rename"
	FaultBackupAfterSegmentDirectorySync BackupFault = "backup-after-segment-directory-sync"
	FaultBackupAfterCheckpointValidation BackupFault = "backup-after-checkpoint-validation"
	FaultBackupAfterManifestTempCreate   BackupFault = "backup-after-manifest-temp-create"
	FaultBackupAfterManifestWrite        BackupFault = "backup-after-manifest-write"
	FaultBackupAfterManifestSync         BackupFault = "backup-after-manifest-sync"
	FaultBackupAfterManifestClose        BackupFault = "backup-after-manifest-close"
	FaultBackupAfterManifestReopen       BackupFault = "backup-after-manifest-reopen"
	FaultBackupAfterManifestRename       BackupFault = "backup-after-manifest-rename"
	FaultBackupAfterRootDirectorySync    BackupFault = "backup-after-root-directory-sync"
	FaultBackupAfterReopen               BackupFault = "backup-after-reopen"
)

type BackupFaultInjector interface{ FailBackupAt(BackupFault) error }

var (
	ErrBackupDestinationNotEmpty  = errors.New("fixed supervisor coherent backup destination is not empty")
	ErrBackupIncomplete           = errors.New("fixed supervisor coherent backup is incomplete")
	ErrBackupOutcomeIndeterminate = errors.New("fixed supervisor coherent backup outcome is indeterminate")
	ErrBackupEvidencePreserved    = errors.New("fixed supervisor backup evidence is preserved")
	ErrRestoreRollbackUncertain   = errors.New("fixed supervisor restore is rollback-uncertain")
	ErrOrphanEvidencePreserved    = errors.New("fixed supervisor orphan evidence is preserved")
	ErrOrphanOutcomeIndeterminate = errors.New("fixed supervisor orphan cleanup outcome is indeterminate")
)

var fixedStoreV2BackupMu sync.Mutex

// CreateCoherentBackup copies one owner-held, fully verified v2 world into an
// already-existing empty owner-only destination. It writes the manifest last
// and never changes the live store.
func (store *FixedFileStoreV2) CreateCoherentBackup(
	ctx context.Context,
	owner ArchiveOwner,
	destination BackupRoot,
	faults BackupFaultInjector,
) (BackupManifest, error) {
	if err := ctx.Err(); err != nil {
		return BackupManifest{}, err
	}
	fixedStoreV2BackupMu.Lock()
	defer fixedStoreV2BackupMu.Unlock()
	loaded, session, err := store.archiveEntry(ctx, owner)
	if err != nil {
		return BackupManifest{}, err
	}
	if err := requireEmptyBackupRoot(destination); err != nil {
		return BackupManifest{}, err
	}
	activeBytes, err := readBoundedStoreFile(store.path, archivestate.MaxSupervisorStateV2Bytes)
	if err != nil {
		return BackupManifest{}, err
	}
	manifest, err := buildBackupManifest(loaded, activeBytes, destination)
	if err != nil {
		return BackupManifest{}, err
	}
	if err := failBackup(faults, FaultBackupAfterPrepare, false); err != nil {
		return BackupManifest{}, err
	}
	archiveDirectory := filepath.Join(destination.path, backupArchiveDirName)
	if err := os.Mkdir(archiveDirectory, 0o700); err != nil {
		return BackupManifest{}, errors.New("create coherent backup archive directory")
	}
	if err := syncDirectory(destination.path); err != nil {
		return BackupManifest{}, errors.New("sync coherent backup root after archive directory creation")
	}
	if err := copyBackupFile(destination.path, backupActiveName, activeBytes, 0o600,
		FaultBackupAfterActiveTempCreate, FaultBackupAfterActiveCopy, FaultBackupAfterActiveSync,
		FaultBackupAfterActiveClose, FaultBackupAfterActiveReopen,
		FaultBackupAfterActiveRename, FaultBackupAfterActiveDirectorySync, faults); err != nil {
		return BackupManifest{}, err
	}
	for index, segment := range loaded.Segments {
		entry := manifest.disk.Segments[index]
		if err := copyBackupFile(archiveDirectory, entry.FileName, segment.Bytes, 0o400,
			FaultBackupAfterSegmentTempCreate, FaultBackupAfterSegmentCopy, FaultBackupAfterSegmentSync,
			FaultBackupAfterSegmentClose, FaultBackupAfterSegmentReopen,
			FaultBackupAfterSegmentRename, FaultBackupAfterSegmentDirectorySync, faults); err != nil {
			return BackupManifest{}, err
		}
	}
	// Check the owner and complete live checkpoint again immediately before
	// the manifest becomes the backup completion marker.
	current, err := loadV2State(store.path)
	if err != nil || current.Active.View().CurrentCheckpoint != loaded.Active.View().CurrentCheckpoint ||
		current.Active.View().SnapshotGeneration != loaded.Active.View().SnapshotGeneration {
		return BackupManifest{}, fmt.Errorf("%w: live checkpoint changed during backup", ErrArchiveStaleTransaction)
	}
	if err := owner.CheckHeld(ctx); err != nil || owner.OwnerSessionID() != session {
		return BackupManifest{}, ErrArchiveOwnerSessionMismatch
	}
	if err := failBackup(faults, FaultBackupAfterCheckpointValidation, false); err != nil {
		return BackupManifest{}, err
	}
	encoded, err := encodeBackupManifest(manifest.disk)
	if err != nil {
		return BackupManifest{}, err
	}
	if err := copyBackupFile(destination.path, backupManifestName, encoded, 0o400,
		FaultBackupAfterManifestTempCreate, FaultBackupAfterManifestWrite, FaultBackupAfterManifestSync,
		FaultBackupAfterManifestClose, FaultBackupAfterManifestReopen,
		FaultBackupAfterManifestRename, FaultBackupAfterRootDirectorySync, faults); err != nil {
		return BackupManifest{}, err
	}
	verified, _, err := VerifyCoherentBackup(ctx, destination)
	if err != nil || verified.View().ManifestDigest != manifest.View().ManifestDigest {
		return BackupManifest{}, fmt.Errorf("reopen coherent backup: %w", err)
	}
	if err := failBackup(faults, FaultBackupAfterReopen, true); err != nil {
		return BackupManifest{}, err
	}
	return verified, nil
}

func buildBackupManifest(loaded loadedV2State, activeBytes []byte, root BackupRoot) (BackupManifest, error) {
	view := loaded.Active.View()
	version := backupFormatV0
	segments := make([]backupSegmentDisk, len(loaded.Segments))
	for index, segment := range loaded.Segments {
		segmentView := segment.Segment.View()
		segments[index] = backupSegmentDisk{
			Ordinal: segmentView.Ordinal, FileName: archiveSegmentName(segment.Segment.Digest()),
			EncodedByteLength: uint64(len(segment.Bytes)), FileDigest: sha256.Sum256(segment.Bytes),
			SegmentDigest: segment.Segment.Digest(), SourceSnapshotGeneration: segmentView.SourceSnapshotGeneration,
			ArchiveGeneration: segmentView.ArchiveGeneration, PriorCheckpointDigest: segmentView.PriorCheckpointDigest,
		}
	}
	disk := backupManifestDisk{
		FormatVersion: &version, StoreFormatVersion: view.StoreFormatVersion,
		ArchiveFormatVersion: archivestate.SupervisorArchiveFormatV0,
		InstallationID:       view.InstallationID, SupervisorID: view.SupervisorID,
		EpochSequence: view.EpochSequence, EpochDigest: view.EpochDigest,
		SnapshotGeneration: view.SnapshotGeneration, ArchiveGeneration: view.ArchiveGeneration,
		DurableTimeHighWater: view.DurableTimeHighWater, PreviousCheckpoint: view.PreviousCheckpoint,
		CurrentCheckpoint: view.CurrentCheckpoint, MigrationGenesisDigest: loaded.Genesis.Digest(),
		DescriptorSetDigest: view.DescriptorSetDigest, CombinedIndexDigest: view.CombinedIndexDigest,
		HotSetDigests: backupHotSetDigestsToDisk(backupHotSetDigests(loaded)), VisibleV1EffectSeedCount: view.VisibleV1EffectSeedCount,
		VisibleV1EffectSeedDigest: view.VisibleV1EffectSeedDigest, HotCounts: view.HotCounts,
		ArchivedCounts: view.ArchivedCounts, TotalCounts: view.TotalCounts,
		ActiveFileName: backupActiveName, ActiveEncodedByteLength: uint64(len(activeBytes)),
		ActiveFileDigest: sha256.Sum256(activeBytes), Segments: segments,
	}
	disk.ManifestDigest = digestBackupManifest(disk)
	return BackupManifest{disk: disk, root: root}, nil
}

func digestBackupManifest(disk backupManifestDisk) BackupManifestDigest {
	disk.ManifestDigest = BackupManifestDigest{}
	payload, _ := json.Marshal(disk)
	hash := sha256.New()
	_, _ = hash.Write([]byte("capsule.supervisor.coherent-backup-manifest.v0"))
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(payload)))
	_, _ = hash.Write(length[:])
	_, _ = hash.Write(payload)
	var result BackupManifestDigest
	copy(result[:], hash.Sum(nil))
	return result
}

func encodeBackupManifest(disk backupManifestDisk) ([]byte, error) {
	encoded, err := json.Marshal(disk)
	if err != nil {
		return nil, errors.New("encode coherent backup manifest")
	}
	encoded = append(encoded, '\n')
	if int64(len(encoded)) > maxBackupManifest {
		return nil, errors.New("coherent backup manifest exceeds capacity")
	}
	return encoded, nil
}

func copyBackupFile(
	directory, finalName string,
	data []byte,
	mode os.FileMode,
	afterCreate, afterWrite, afterSync, afterClose, afterReopen, afterRename, afterDirectorySync BackupFault,
	faults BackupFaultInjector,
) error {
	temporary, err := os.CreateTemp(directory, ".capsule-supervisor-backup-*.tmp")
	if err != nil {
		return errors.New("create coherent backup temporary file")
	}
	temporaryPath := temporary.Name()
	committed := false
	defer func() {
		_ = temporary.Close()
		if !committed {
			_ = os.Remove(temporaryPath) //nolint:gosec // exact CreateTemp path in validated destination.
		}
	}()
	if err := temporary.Chmod(mode); err != nil {
		return errors.New("protect coherent backup temporary file")
	}
	if err := failBackup(faults, afterCreate, false); err != nil {
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		return errors.New("copy coherent backup file")
	}
	if err := failBackup(faults, afterWrite, false); err != nil {
		return err
	}
	if err := temporary.Sync(); err != nil {
		return errors.New("sync coherent backup file")
	}
	if err := failBackup(faults, afterSync, false); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return errors.New("close coherent backup file")
	}
	if err := failBackup(faults, afterClose, false); err != nil {
		return err
	}
	if err := validateArchiveFileShape(temporaryPath, mode); err != nil {
		return err
	}
	reopened, err := readBoundedStoreFile(temporaryPath, int64(len(data)))
	if err != nil || !bytes.Equal(reopened, data) {
		return errors.New("reopen coherent backup temporary file")
	}
	if err := failBackup(faults, afterReopen, false); err != nil {
		return err
	}
	finalPath := filepath.Join(directory, finalName)
	if err := os.Rename(temporaryPath, finalPath); err != nil { //nolint:gosec // fixed or digest-derived child name.
		return errors.New("publish coherent backup file")
	}
	committed = true
	if err := failBackup(faults, afterRename, true); err != nil {
		return err
	}
	if err := syncDirectory(directory); err != nil {
		return fmt.Errorf("%w: sync coherent backup directory", ErrBackupOutcomeIndeterminate)
	}
	return failBackup(faults, afterDirectorySync, true)
}

func failBackup(faults BackupFaultInjector, point BackupFault, indeterminate bool) error {
	if faults == nil {
		return nil
	}
	if err := faults.FailBackupAt(point); err != nil {
		if indeterminate {
			return fmt.Errorf("%w: injected at %s", ErrBackupOutcomeIndeterminate, point)
		}
		return fmt.Errorf("fixed supervisor coherent backup confirmed failure at %s", point)
	}
	return nil
}

func requireEmptyBackupRoot(root BackupRoot) error {
	if err := validateOwnedDirectory(root.path, 0o700); err != nil {
		return err
	}
	entries, err := os.ReadDir(root.path)
	if err != nil {
		return errors.New("read coherent backup destination")
	}
	if len(entries) != 0 {
		return ErrBackupDestinationNotEmpty
	}
	return nil
}

// VerifyCoherentBackup independently verifies the manifest inventory and then
// runs the complete v2 verifier over the copied active/segment set.
func VerifyCoherentBackup(ctx context.Context, root BackupRoot) (BackupManifest, ArchiveVerificationReport, error) {
	if err := ctx.Err(); err != nil {
		return BackupManifest{}, ArchiveVerificationReport{}, err
	}
	report := ArchiveVerificationReport{Version: 0, Classification: ArchiveVerificationIncomplete}
	if err := validateOwnedDirectory(root.path, 0o700); err != nil {
		return BackupManifest{}, report, err
	}
	entries, err := os.ReadDir(root.path)
	if err != nil {
		return BackupManifest{}, report, errors.New("read coherent backup root")
	}
	wantRoot := map[string]bool{backupManifestName: false, backupActiveName: false, backupArchiveDirName: false}
	for _, entry := range entries {
		if _, exists := wantRoot[entry.Name()]; !exists {
			report.Classification = ArchiveVerificationExtraArtifacts
			return BackupManifest{}, report, ErrBackupEvidencePreserved
		}
		wantRoot[entry.Name()] = true
	}
	if !wantRoot[backupManifestName] || !wantRoot[backupActiveName] || !wantRoot[backupArchiveDirName] {
		return BackupManifest{}, report, ErrBackupIncomplete
	}
	manifestPath := filepath.Join(root.path, backupManifestName)
	if err := validateArchiveFileShape(manifestPath, 0o400); err != nil {
		report.Classification = ArchiveVerificationCorrupt
		return BackupManifest{}, report, ErrBackupEvidencePreserved
	}
	manifestBytes, err := readBoundedStoreFile(manifestPath, maxBackupManifest)
	if err != nil {
		report.Classification = ArchiveVerificationCorrupt
		return BackupManifest{}, report, ErrBackupEvidencePreserved
	}
	var disk backupManifestDisk
	if err := decodeOneClosedJSON(manifestBytes, &disk); err != nil || disk.FormatVersion == nil || *disk.FormatVersion != backupFormatV0 ||
		disk.StoreFormatVersion != archivestate.SupervisorStoreFormatV2 || disk.ArchiveFormatVersion != archivestate.SupervisorArchiveFormatV0 ||
		disk.ActiveFileName != backupActiveName || disk.Segments == nil || disk.ManifestDigest != digestBackupManifest(disk) {
		report.Classification = ArchiveVerificationCorrupt
		return BackupManifest{}, report, ErrBackupEvidencePreserved
	}
	activePath := filepath.Join(root.path, backupActiveName)
	if err := validateArchiveFileShape(activePath, 0o600); err != nil {
		report.Classification = ArchiveVerificationCorrupt
		return BackupManifest{}, report, ErrBackupEvidencePreserved
	}
	activeBytes, err := readBoundedStoreFile(activePath, archivestate.MaxSupervisorStateV2Bytes)
	if err != nil || uint64(len(activeBytes)) != disk.ActiveEncodedByteLength || sha256.Sum256(activeBytes) != disk.ActiveFileDigest {
		report.Classification = ArchiveVerificationCorrupt
		return BackupManifest{}, report, ErrBackupEvidencePreserved
	}
	archiveDirectory := filepath.Join(root.path, backupArchiveDirName)
	if err := validateOwnedDirectory(archiveDirectory, 0o700); err != nil {
		report.Classification = ArchiveVerificationCorrupt
		return BackupManifest{}, report, ErrBackupEvidencePreserved
	}
	archiveEntries, err := os.ReadDir(archiveDirectory)
	if err != nil {
		return BackupManifest{}, report, ErrBackupEvidencePreserved
	}
	wantSegments := make(map[string]backupSegmentDisk, len(disk.Segments))
	for index, segment := range disk.Segments {
		if segment.Ordinal != archivestate.ArchiveOrdinal(index+1) || segment.FileName != archiveSegmentName(segment.SegmentDigest) {
			report.Classification = ArchiveVerificationMixedGeneration
			return BackupManifest{}, report, ErrBackupEvidencePreserved
		}
		if _, exists := wantSegments[segment.FileName]; exists {
			report.Classification = ArchiveVerificationCorrupt
			return BackupManifest{}, report, ErrBackupEvidencePreserved
		}
		wantSegments[segment.FileName] = segment
	}
	if len(archiveEntries) != len(wantSegments) {
		report.Classification = ArchiveVerificationExtraArtifacts
		return BackupManifest{}, report, ErrBackupEvidencePreserved
	}
	for _, entry := range archiveEntries {
		segment, exists := wantSegments[entry.Name()]
		if !exists {
			report.Classification = ArchiveVerificationExtraArtifacts
			return BackupManifest{}, report, ErrBackupEvidencePreserved
		}
		path := filepath.Join(archiveDirectory, entry.Name())
		if err := validateArchiveFileShape(path, 0o400); err != nil {
			report.Classification = ArchiveVerificationCorrupt
			return BackupManifest{}, report, ErrBackupEvidencePreserved
		}
		data, readErr := readBoundedStoreFile(path, archivestate.MaxSupervisorArchiveBytes)
		if readErr != nil || uint64(len(data)) != segment.EncodedByteLength || sha256.Sum256(data) != segment.FileDigest {
			report.Classification = ArchiveVerificationCorrupt
			return BackupManifest{}, report, ErrBackupEvidencePreserved
		}
	}
	loaded, err := loadV2State(activePath)
	if err != nil {
		report.Classification = ArchiveVerificationCorrupt
		return BackupManifest{}, report, ErrBackupEvidencePreserved
	}
	report = reportFromLoaded(loaded, ArchiveVerificationVerified)
	if mismatch := backupManifestMismatch(loaded, disk); mismatch != "" {
		report.Classification = ArchiveVerificationMixedGeneration
		return BackupManifest{}, report, fmt.Errorf("%w: %s", ErrBackupEvidencePreserved, mismatch)
	}
	for index, segment := range loaded.Segments {
		manifestSegment := disk.Segments[index]
		segmentView := segment.Segment.View()
		if segment.Segment.Digest() != manifestSegment.SegmentDigest || segmentView.Ordinal != manifestSegment.Ordinal ||
			segmentView.SourceSnapshotGeneration != manifestSegment.SourceSnapshotGeneration ||
			segmentView.ArchiveGeneration != manifestSegment.ArchiveGeneration ||
			segmentView.PriorCheckpointDigest != manifestSegment.PriorCheckpointDigest {
			report.Classification = ArchiveVerificationMixedGeneration
			return BackupManifest{}, report, ErrBackupEvidencePreserved
		}
	}
	return BackupManifest{disk: disk, root: root}, report, nil
}

func backupManifestMismatch(loaded loadedV2State, disk backupManifestDisk) string {
	view := loaded.Active.View()
	checks := []struct {
		code string
		bad  bool
	}{
		{"orphan-count", loaded.Orphans != 0},
		{"installation", view.InstallationID != disk.InstallationID},
		{"supervisor", view.SupervisorID != disk.SupervisorID},
		{"epoch-sequence", view.EpochSequence != disk.EpochSequence},
		{"epoch-digest", view.EpochDigest != disk.EpochDigest},
		{"snapshot-generation", view.SnapshotGeneration != disk.SnapshotGeneration},
		{"archive-generation", view.ArchiveGeneration != disk.ArchiveGeneration},
		{"time-high-water", view.DurableTimeHighWater != disk.DurableTimeHighWater},
		{"previous-checkpoint", view.PreviousCheckpoint != disk.PreviousCheckpoint},
		{"current-checkpoint", view.CurrentCheckpoint != disk.CurrentCheckpoint},
		{"migration-genesis", loaded.Genesis.Digest() != disk.MigrationGenesisDigest},
		{"descriptor-set", view.DescriptorSetDigest != disk.DescriptorSetDigest},
		{"combined-index", view.CombinedIndexDigest != disk.CombinedIndexDigest},
		{"hot-registration-set", backupHotSetDigests(loaded).Registrations != disk.HotSetDigests.Registrations},
		{"hot-approval-set", backupHotSetDigests(loaded).Approvals != disk.HotSetDigests.Approvals},
		{"hot-attempt-set", backupHotSetDigests(loaded).Attempts != disk.HotSetDigests.Attempts},
		{"hot-lifecycle-set", backupHotSetDigests(loaded).Lifecycles != disk.HotSetDigests.Lifecycles},
		{"hot-effect-set", backupHotSetDigests(loaded).EffectTombstones != disk.HotSetDigests.EffectTombstones},
		{"visible-v1-count", view.VisibleV1EffectSeedCount != disk.VisibleV1EffectSeedCount},
		{"visible-v1-digest", view.VisibleV1EffectSeedDigest != disk.VisibleV1EffectSeedDigest},
		{"hot-counts", view.HotCounts != disk.HotCounts},
		{"archived-counts", view.ArchivedCounts != disk.ArchivedCounts},
		{"total-counts", view.TotalCounts != disk.TotalCounts},
		{"segment-count", len(loaded.Segments) != len(disk.Segments)},
	}
	for _, check := range checks {
		if check.bad {
			return check.code
		}
	}
	return ""
}

func backupHotSetDigests(loaded loadedV2State) archivestate.HotSetDigests {
	digests := archivestate.HotSetDigests{
		Registrations: loaded.State.RegistrationSetDigest,
		Approvals:     loaded.State.ApprovalSetDigest,
		Attempts:      loaded.State.AttemptSetDigest,
		Lifecycles:    [32]byte(loaded.Envelope.LifecycleSetDigest),
	}
	if loaded.Envelope.EffectTombstoneSetDigest != nil {
		digests.EffectTombstones = *loaded.Envelope.EffectTombstoneSetDigest
	}
	return digests
}

func backupHotSetDigestsToDisk(digests archivestate.HotSetDigests) backupHotSetDigestsDisk {
	return backupHotSetDigestsDisk{
		Registrations: digests.Registrations, Approvals: digests.Approvals, Attempts: digests.Attempts,
		Lifecycles: digests.Lifecycles, EffectTombstones: digests.EffectTombstones,
	}
}

// VerifyArchiveSet is the deterministic read-only full verifier over a trusted
// live-store capability. It emits no path or retained bytes.
func VerifyArchiveSet(ctx context.Context, root StoreRoot, mode VerificationMode) (ArchiveVerificationReport, error) {
	if err := ctx.Err(); err != nil {
		return ArchiveVerificationReport{}, err
	}
	if mode != VerificationFull {
		return ArchiveVerificationReport{}, errors.New("unsupported fixed supervisor archive verification mode")
	}
	loaded, err := loadV2State(root.path)
	if err != nil {
		return ArchiveVerificationReport{Version: 0, Classification: ArchiveVerificationCorrupt}, fmt.Errorf("%w: %v", ErrStoreRepairRequired, err)
	}
	classification := ArchiveVerificationVerified
	if loaded.Orphans != 0 {
		classification = ArchiveVerificationVerifiedWithOrphans
	}
	return reportFromLoaded(loaded, classification), nil
}

func reportFromLoaded(loaded loadedV2State, classification ArchiveVerificationClassification) ArchiveVerificationReport {
	view := loaded.Active.View()
	segmentCount := len(loaded.Segments)
	if segmentCount == 0 && len(loaded.Envelope.Descriptors) != 0 {
		segmentCount = len(loaded.Envelope.Descriptors)
	}
	return ArchiveVerificationReport{
		Version: 0, Classification: classification, StoreFormatVersion: view.StoreFormatVersion,
		ArchiveFormatVersion: archivestate.SupervisorArchiveFormatV0, InstallationID: view.InstallationID,
		SupervisorID: view.SupervisorID, EpochSequence: view.EpochSequence, EpochDigest: view.EpochDigest,
		SnapshotGeneration: view.SnapshotGeneration, ArchiveGeneration: view.ArchiveGeneration,
		TimeHighWater: view.DurableTimeHighWater, PreviousCheckpoint: view.PreviousCheckpoint,
		CheckpointDigest: view.CurrentCheckpoint.Digest, ActiveCounts: view.HotCounts,
		ArchivedCounts: view.ArchivedCounts, TotalCounts: view.TotalCounts,
		SegmentCount: uint16(segmentCount), KnownUnreferenced: loaded.Orphans, //nolint:gosec // both bounded before conversion.
	}
}

// IndependentLatestCheckpoint is an injected protected-anchor fixture. F5
// does not implement, persist, or select the production anchor.
type IndependentLatestCheckpoint struct {
	InstallationID       v0candidate.InstallationID
	SupervisorID         v0candidate.SupervisorID
	EpochSequence        v0candidate.UInt53
	EpochDigest          v0candidate.TrustEpochDigest
	SnapshotGeneration   v0candidate.PositiveUInt53
	ArchiveGeneration    archivestate.ArchiveGeneration
	DurableTimeHighWater v0candidate.UInt53
	PreviousCheckpoint   archivestate.ArchiveCheckpointReference
	CurrentCheckpoint    archivestate.ArchiveCheckpointReference
}

func LatestCheckpointFromManifest(manifest BackupManifest) IndependentLatestCheckpoint {
	view := manifest.View()
	return IndependentLatestCheckpoint{
		InstallationID: view.InstallationID, SupervisorID: view.SupervisorID,
		EpochSequence: view.EpochSequence, EpochDigest: view.EpochDigest,
		SnapshotGeneration: view.SnapshotGeneration, ArchiveGeneration: view.ArchiveGeneration,
		DurableTimeHighWater: view.DurableTimeHighWater, PreviousCheckpoint: view.PreviousCheckpoint,
		CurrentCheckpoint: view.CurrentCheckpoint,
	}
}

type RestoreAdmissionReport struct {
	Classification ArchiveVerificationClassification
	Eligible       bool
	Manifest       BackupManifestView
	Predecessor    ArchiveVerificationReport
}

// VerifyRestoreAdmission is read-only. Exact anchor equality is merely
// eligible for a future separately authorized ceremony; F5 never replaces the
// live predecessor, enables attempts, or changes an epoch.
func (store *FixedFileStoreV2) VerifyRestoreAdmission(
	ctx context.Context,
	owner ArchiveOwner,
	backup BackupRoot,
	anchor *IndependentLatestCheckpoint,
) (RestoreAdmissionReport, error) {
	predecessor, session, err := store.archiveEntry(ctx, owner)
	if err != nil {
		return RestoreAdmissionReport{}, err
	}
	manifest, _, err := VerifyCoherentBackup(ctx, backup)
	if err != nil {
		return RestoreAdmissionReport{Classification: ArchiveVerificationCorrupt, Predecessor: reportFromLoaded(predecessor, ArchiveVerificationVerified)}, err
	}
	report := RestoreAdmissionReport{
		Classification: ArchiveVerificationRollbackUncertain, Manifest: manifest.View(),
		Predecessor: reportFromLoaded(predecessor, ArchiveVerificationVerified),
	}
	if anchor != nil && *anchor == LatestCheckpointFromManifest(manifest) {
		report.Classification = ArchiveVerificationEligibleFutureRestore
		report.Eligible = true
	}
	current, reopenErr := loadV2State(store.path)
	if reopenErr != nil || current.Active.View().CurrentCheckpoint != predecessor.Active.View().CurrentCheckpoint ||
		current.Active.View().SnapshotGeneration != predecessor.Active.View().SnapshotGeneration ||
		owner.OwnerSessionID() != session || owner.CheckHeld(ctx) != nil {
		return RestoreAdmissionReport{}, fmt.Errorf("%w: live predecessor changed during restore admission", ErrStoreRepairRequired)
	}
	if !report.Eligible {
		return report, ErrRestoreRollbackUncertain
	}
	return report, nil
}

// ArchiveArtifactClassification is the closed orphan inventory vocabulary.
type ArchiveArtifactClassification string

const (
	ArchiveArtifactReferenced        ArchiveArtifactClassification = "referenced"
	ArchiveArtifactKnownUnreferenced ArchiveArtifactClassification = "known-unreferenced"
	ArchiveArtifactBackupReferenced  ArchiveArtifactClassification = "backup-referenced"
	ArchiveArtifactUnknown           ArchiveArtifactClassification = "unknown"
	ArchiveArtifactCorrupt           ArchiveArtifactClassification = "corrupt"
	ArchiveArtifactMixedGeneration   ArchiveArtifactClassification = "mixed-generation"
	ArchiveArtifactCrossInstallation ArchiveArtifactClassification = "cross-installation"
)

// KnownUnreferencedOrphan is sealed to one fully decoded inode, digest,
// checkpoint, owner session, and backup-manifest scan. It contains no caller-
// supplied deletion path.
type KnownUnreferencedOrphan struct {
	digest     archivestate.ArchiveSegmentDigest
	device     string
	inode      uint64
	checkpoint archivestate.ArchiveCheckpointReference
	session    lifecyclestate.OwnerSessionID
}

func (candidate KnownUnreferencedOrphan) Digest() archivestate.ArchiveSegmentDigest {
	return candidate.digest
}

type ArchiveInventory struct {
	report     ArchiveVerificationReport
	candidates []KnownUnreferencedOrphan
	backups    []BackupRoot
}

func (inventory ArchiveInventory) Report() ArchiveVerificationReport { return inventory.report }

func (inventory ArchiveInventory) KnownUnreferenced() []KnownUnreferencedOrphan {
	return append([]KnownUnreferencedOrphan{}, inventory.candidates...)
}

// InventoryArchiveArtifacts performs offline-only classification. Unknown,
// corrupt, mixed-generation, and cross-installation bytes are reported and
// preserved; they suppress every deletion candidate.
func (store *FixedFileStoreV2) InventoryArchiveArtifacts(
	ctx context.Context,
	owner ArchiveOwner,
	backups []BackupRoot,
) (ArchiveInventory, error) {
	loaded, session, err := store.archiveEntry(ctx, owner)
	fullVerificationFailed := false
	if err != nil {
		// archiveEntry is strict about unsafe unreferenced entries. Continue only
		// far enough to produce the bounded evidence-preservation inventory.
		if owner == nil || owner.CheckHeld(ctx) != nil {
			return ArchiveInventory{}, err
		}
		loaded, err = loadActiveV2EnvelopeForInventory(store.path)
		if err != nil {
			return ArchiveInventory{}, fmt.Errorf("%w: %v", ErrStoreRepairRequired, err)
		}
		fullVerificationFailed = true
		session = owner.OwnerSessionID()
	}
	report := reportFromLoaded(loaded, ArchiveVerificationVerified)
	report.KnownUnreferenced = 0
	backupDigests := make(map[archivestate.ArchiveSegmentDigest]struct{})
	for _, backup := range backups {
		manifest, _, verifyErr := VerifyCoherentBackup(ctx, backup)
		if verifyErr != nil {
			report.Classification = ArchiveVerificationCorrupt
			return ArchiveInventory{report: report}, ErrOrphanEvidencePreserved
		}
		for _, segment := range manifest.disk.Segments {
			backupDigests[segment.SegmentDigest] = struct{}{}
		}
	}
	referenced := make(map[string]struct{}, len(loaded.Envelope.Descriptors))
	for _, descriptor := range loaded.Envelope.Descriptors {
		referenced[archiveSegmentName(descriptor.SegmentDigest)] = struct{}{}
	}
	root := archiveRoot(store.path)
	entries, readErr := os.ReadDir(root) //nolint:gosec // derived trusted archive root.
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) && len(referenced) == 0 {
			return ArchiveInventory{report: report, backups: append([]BackupRoot{}, backups...)}, nil
		}
		return ArchiveInventory{}, ErrOrphanEvidencePreserved
	}
	candidates := make([]KnownUnreferencedOrphan, 0)
	// A fallback inventory is reporting-only. If the ordinary full verifier
	// failed for any reason, preserve every unreferenced byte and issue no
	// deletion capability even when an individual orphan decodes cleanly.
	unsafeEvidence := fullVerificationFailed
	for _, entry := range entries {
		name := entry.Name()
		if _, ok := referenced[name]; ok {
			continue
		}
		path := filepath.Join(root, name)
		classification, segment, info := classifyUnreferencedArchiveArtifact(path, name, loaded)
		switch classification {
		case ArchiveArtifactKnownUnreferenced:
			if _, protected := backupDigests[segment.Segment.Digest()]; protected {
				report.BackupReferenced++
				continue
			}
			stat := info.Sys().(*syscall.Stat_t)
			report.KnownUnreferenced++
			candidates = append(candidates, KnownUnreferencedOrphan{
				digest: segment.Segment.Digest(), device: fmt.Sprint(stat.Dev), inode: uint64(stat.Ino),
				checkpoint: loaded.Active.View().CurrentCheckpoint, session: session,
			})
		case ArchiveArtifactCrossInstallation:
			report.CrossInstallation++
			unsafeEvidence = true
		case ArchiveArtifactMixedGeneration:
			report.MixedGeneration++
			unsafeEvidence = true
		case ArchiveArtifactCorrupt:
			report.CorruptArtifacts++
			unsafeEvidence = true
		default:
			report.UnknownArtifacts++
			unsafeEvidence = true
		}
	}
	if unsafeEvidence {
		report.Classification = ArchiveVerificationCorrupt
		return ArchiveInventory{report: report, backups: append([]BackupRoot{}, backups...)}, ErrOrphanEvidencePreserved
	}
	if len(candidates) > 0 {
		report.Classification = ArchiveVerificationVerifiedWithOrphans
	}
	sort.Slice(candidates, func(i, j int) bool { return bytes.Compare(candidates[i].digest[:], candidates[j].digest[:]) < 0 })
	return ArchiveInventory{report: report, candidates: candidates, backups: append([]BackupRoot{}, backups...)}, nil
}

func classifyUnreferencedArchiveArtifact(path, name string, loaded loadedV2State) (ArchiveArtifactClassification, loadedArchiveSegmentV0, os.FileInfo) {
	if strings.HasPrefix(name, ".capsule-supervisor-") || !strings.HasPrefix(name, "segment-") || !strings.HasSuffix(name, ".json") {
		return ArchiveArtifactUnknown, loadedArchiveSegmentV0{}, nil
	}
	info, err := os.Lstat(path) //nolint:gosec // exact child returned by ReadDir.
	if err != nil || validateArchiveFileShape(path, 0o400) != nil {
		return ArchiveArtifactUnknown, loadedArchiveSegmentV0{}, info
	}
	data, err := readBoundedStoreFile(path, archivestate.MaxSupervisorArchiveBytes)
	if err != nil {
		return ArchiveArtifactCorrupt, loadedArchiveSegmentV0{}, info
	}
	segment, err := decodeArchiveSegmentV0(data)
	if err != nil || name != archiveSegmentName(segment.Segment.Digest()) {
		return ArchiveArtifactCorrupt, loadedArchiveSegmentV0{}, info
	}
	view := segment.Segment.View()
	active := loaded.Active.View()
	if view.InstallationID != active.InstallationID || view.SupervisorID != active.SupervisorID ||
		view.EpochSequence != active.EpochSequence || view.EpochDigest != active.EpochDigest {
		return ArchiveArtifactCrossInstallation, segment, info
	}
	if view.Ordinal != archivestate.ArchiveOrdinal(len(loaded.Envelope.Descriptors)+1) ||
		view.SourceSnapshotGeneration != active.SnapshotGeneration ||
		view.ArchiveGeneration != archivestate.ArchiveGeneration(uint64(active.ArchiveGeneration)+1) ||
		view.PriorCheckpointDigest != active.CurrentCheckpoint.Digest ||
		view.DurableTimeHighWater > active.DurableTimeHighWater {
		return ArchiveArtifactMixedGeneration, segment, info
	}
	return ArchiveArtifactKnownUnreferenced, segment, info
}

// DeleteKnownUnreferencedOrphan is the sole explicit F5 deletion operation.
// It rechecks the complete live/backup reference set and inode before unlink.
// It never chooses an orphan automatically.
func (store *FixedFileStoreV2) DeleteKnownUnreferencedOrphan(
	ctx context.Context,
	owner ArchiveOwner,
	inventory ArchiveInventory,
	candidate KnownUnreferencedOrphan,
	faults OrphanFaultInjector,
) error {
	if candidate.digest == (archivestate.ArchiveSegmentDigest{}) || candidate.session.IsZero() {
		return ErrOrphanEvidencePreserved
	}
	fresh, err := store.InventoryArchiveArtifacts(ctx, owner, inventory.backups)
	if err != nil {
		return err
	}
	if owner.OwnerSessionID() != candidate.session || fresh.report.CheckpointDigest != candidate.checkpoint.Digest {
		return ErrOrphanEvidencePreserved
	}
	path := archiveSegmentPath(store.path, candidate.digest)
	var matched *KnownUnreferencedOrphan
	for index := range fresh.candidates {
		if fresh.candidates[index] == candidate {
			matched = &fresh.candidates[index]
			break
		}
	}
	if matched == nil {
		// Exact retry after an indeterminate unlink may observe absence. The
		// sealed candidate, unchanged active checkpoint, same owner session, and
		// repeated complete backup-reference scan make that outcome idempotent.
		if _, statErr := os.Lstat(path); errors.Is(statErr, os.ErrNotExist) { //nolint:gosec // digest-derived trusted child.
			return nil
		}
		return ErrOrphanEvidencePreserved
	}
	info, err := os.Lstat(path) //nolint:gosec // digest-derived child of trusted archive root.
	if err != nil {
		return ErrOrphanEvidencePreserved
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || fmt.Sprint(stat.Dev) != candidate.device || uint64(stat.Ino) != candidate.inode {
		return ErrOrphanEvidencePreserved
	}
	if err := failOrphan(faults, FaultOrphanBeforeRemove, false); err != nil {
		return err
	}
	if err := os.Remove(path); err != nil { //nolint:gosec // sealed digest/inode-verified candidate.
		return errors.New("remove known unreferenced archive segment")
	}
	if err := failOrphan(faults, FaultOrphanAfterRemove, true); err != nil {
		return err
	}
	if err := syncDirectory(archiveRoot(store.path)); err != nil {
		return fmt.Errorf("%w: sync archive directory after orphan removal", ErrOrphanOutcomeIndeterminate)
	}
	if err := failOrphan(faults, FaultOrphanAfterDirectorySync, true); err != nil {
		return err
	}
	if _, err := loadV2State(store.path); err != nil {
		return fmt.Errorf("%w: reopen after orphan removal", ErrStoreRepairRequired)
	}
	return failOrphan(faults, FaultOrphanAfterReopen, true)
}

type OrphanFault string

const (
	FaultOrphanBeforeRemove       OrphanFault = "orphan-before-remove"
	FaultOrphanAfterRemove        OrphanFault = "orphan-after-remove"
	FaultOrphanAfterDirectorySync OrphanFault = "orphan-after-directory-sync"
	FaultOrphanAfterReopen        OrphanFault = "orphan-after-reopen"
)

type OrphanFaultInjector interface{ FailOrphanAt(OrphanFault) error }

func failOrphan(faults OrphanFaultInjector, point OrphanFault, indeterminate bool) error {
	if faults == nil {
		return nil
	}
	if err := faults.FailOrphanAt(point); err != nil {
		if indeterminate {
			return fmt.Errorf("%w: injected at %s", ErrOrphanOutcomeIndeterminate, point)
		}
		return fmt.Errorf("known orphan removal confirmed failure at %s", point)
	}
	return nil
}

// loadActiveV2EnvelopeForInventory decodes only enough active metadata to
// classify unsafe unreferenced evidence. It deliberately supplies no deletion
// candidate until the ordinary full verifier succeeds.
func loadActiveV2EnvelopeForInventory(path string) (loadedV2State, error) {
	data, err := readBoundedStoreFile(path, archivestate.MaxSupervisorStateV2Bytes)
	if err != nil {
		return loadedV2State{}, err
	}
	var envelope diskEnvelopeV2
	if err := decodeOneClosedJSON(data, &envelope); err != nil || envelope.StoreFormatVersion == nil ||
		*envelope.StoreFormatVersion != archivestate.SupervisorStoreFormatV2 {
		return loadedV2State{}, errors.New("decode fixed supervisor v2 inventory state")
	}
	indexes, err := archiveIndexesFromDisk(envelope.Indexes)
	if err != nil {
		return loadedV2State{}, err
	}
	descriptors := make([]archivestate.ArchiveDescriptor, len(envelope.Descriptors))
	for index, descriptor := range envelope.Descriptors {
		descriptors[index], err = archivestate.NewArchiveDescriptor(descriptor)
		if err != nil {
			return loadedV2State{}, err
		}
	}
	active, err := archivestate.NewActiveStateV2(archivestate.ActiveStateV2View{
		StoreFormatVersion: *envelope.StoreFormatVersion, MigrationSourceVersion: *envelope.MigrationSourceVersion,
		SnapshotGeneration: envelope.SnapshotGeneration, ArchiveGeneration: envelope.ArchiveGeneration,
		InstallationID: envelope.State.InstallationID, SupervisorID: envelope.State.SupervisorID,
		EpochSequence: envelope.State.EpochSequence, EpochDigest: envelope.State.EpochDigest,
		DurableTimeHighWater: envelope.State.TimeHighWaterUnixSeconds, Descriptors: descriptors,
		DescriptorSetDigest: envelope.DescriptorSetDigest, Indexes: indexes, IndexDigests: envelope.IndexDigests,
		CombinedIndexDigest: envelope.CombinedIndexDigest, EffectTombstoneCoverage: envelope.EffectTombstoneCoverage,
		VisibleV1EffectSeedCount: envelope.VisibleV1EffectSeedCount, VisibleV1EffectSeedDigest: envelope.VisibleV1EffectSeedDigest,
		PreviousCheckpoint: envelope.PreviousCheckpoint, CurrentCheckpoint: envelope.CurrentCheckpoint,
		HotCounts: envelope.HotCounts, ArchivedCounts: envelope.ArchivedCounts, TotalCounts: envelope.TotalCounts,
	})
	if err != nil {
		return loadedV2State{}, err
	}
	return loadedV2State{Envelope: envelope, State: envelope.State, Active: active}, nil
}

func validateOwnedDirectory(path string, mode os.FileMode) error {
	info, err := os.Lstat(path) //nolint:gosec // trusted installation/test capability.
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.IsDir() || info.Mode().Perm() != mode {
		return errors.New("owner-only directory shape is unsafe")
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || int(stat.Uid) != os.Geteuid() || stat.Nlink < 1 {
		return errors.New("owner-only directory identity is unsafe")
	}
	return nil
}

func syncDirectory(path string) error {
	directory, err := os.Open(path) //nolint:gosec // validated trusted directory capability.
	if err != nil {
		return err
	}
	syncErr, closeErr := directory.Sync(), directory.Close()
	if syncErr != nil {
		return syncErr
	}
	return closeErr
}
