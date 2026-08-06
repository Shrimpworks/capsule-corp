package completioncomposer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

type StoreCommitCheckpoint string

const (
	CheckpointBeforeCompletionRename    StoreCommitCheckpoint = "before-completion-rename"
	CheckpointAfterCompletionRename     StoreCommitCheckpoint = "after-completion-rename"
	CheckpointBeforeCompletionStoreLink StoreCommitCheckpoint = "before-completion-store-link"
	CheckpointAfterCompletionStoreLink  StoreCommitCheckpoint = "after-completion-store-link"
)

type StoreFaultHook func(StoreCommitCheckpoint) error

type durableStoreSnapshot struct {
	ObjectType         string                    `json:"objectType"`
	ObjectVersion      uint64                    `json:"objectVersion"`
	InstallationID     string                    `json:"installationId"`
	Generation         uint64                    `json:"generation"`
	LastCommitSequence uint64                    `json:"lastCommitSequence"`
	Records            []storedDurableCompletion `json:"records"`
}

type storedDurableCompletion struct {
	ObjectType            string `json:"objectType"`
	ObjectVersion         uint64 `json:"objectVersion"`
	CommitSequence        uint64 `json:"commitSequence"`
	Commit                string `json:"commit"`
	ResultIntegrity       string `json:"resultIntegrity"`
	RunnerIdentity        string `json:"runnerIdentity"`
	Teardown              string `json:"teardown"`
	DescendantAbsence     string `json:"descendantAbsence"`
	Terminal              string `json:"terminalClassification"`
	AttemptID             string `json:"attemptId"`
	ApprovalID            string `json:"approvalId"`
	RegistrationID        string `json:"registrationId"`
	PlanDigest            string `json:"planDigest"`
	InstallationID        string `json:"installationId"`
	EpochDigest           string `json:"epochDigest"`
	SupervisorID          string `json:"supervisorId"`
	ApprovalDigest        string `json:"approvalPayloadDigest"`
	Authorization         string `json:"authorizationIdentity"`
	SourceManifest        string `json:"sourceManifestDigest"`
	RuntimeDigest         string `json:"runtimeBundleManifestDigest"`
	ProfileDigest         string `json:"profileRegistryEntryDigest"`
	BackendValidation     string `json:"backendValidationRecordDigest"`
	BackendConfiguration  string `json:"backendConfigurationDigest"`
	BackendImplementation string `json:"backendImplementationDigest"`
	BackendInstance       string `json:"backendInstanceDigest"`
	ImmutableBinding      string `json:"immutableBindingDigest"`
	CompletionDigest      string `json:"completionRecordDigest"`
	ResultJSON            string `json:"resultJsonBase64"`
	ResultJSONDigest      string `json:"resultJsonDigest"`
	LifecycleState        string `json:"lifecycleState"`
	LifecycleSequence     uint64 `json:"lifecycleOperationSequence"`
	Reconciliation        string `json:"lastReconciliation"`
	CleanupRequired       bool   `json:"cleanupRequired"`
	Transcript            string `json:"transcriptBase64"`
	TranscriptDigest      string `json:"transcriptDigest"`
	PublicSummary         string `json:"publicSummaryBase64"`
	DurableRecordDigest   string `json:"durableRecordDigest"`
}

// FixedFileCompletionStore is a bounded single-owner local conformance
// oracle. Callers must provide external ownership; this type supplies no
// installed lock, protected root, rollback defense, archive, or production DB.
type FixedFileCompletionStore struct {
	mu       sync.Mutex
	path     string
	snapshot durableStoreSnapshot
	fault    StoreFaultHook
	fenced   bool
}

func CreateFixedFileCompletionStore(
	path string,
	installationID v0candidate.InstallationID,
	fault StoreFaultHook,
) (*FixedFileCompletionStore, error) {
	if installationID == (v0candidate.InstallationID{}) {
		return nil, classified(ClassificationSchema, "completion-store-installation")
	}
	snapshot := durableStoreSnapshot{
		ObjectType:     durableCompletionStoreType,
		ObjectVersion:  DurableCompletionStoreVersion,
		InstallationID: hex.EncodeToString(installationID[:]),
		Records:        []storedDurableCompletion{},
	}
	exact, err := encodeStoreSnapshot(snapshot)
	if err != nil {
		return nil, err
	}
	if err := createSnapshotFile(path, exact, fault); err != nil {
		return nil, err
	}
	return &FixedFileCompletionStore{path: path, snapshot: snapshot, fault: fault}, nil
}

// createSnapshotFile writes exact to a new file at path, following the same
// write-temp/sync/atomically-publish/sync-directory invariant as publish
// (ADR-0042). Unlike publish's rename, it must not silently replace an
// existing store, so it links the synced temporary file into path instead:
// os.Link fails closed with an existing-file error if path is already
// occupied, and — like O_EXCL — never leaves a partial file at path, because
// a failed link touches only the temporary name that the deferred cleanup
// then removes.
func createSnapshotFile(path string, exact []byte, fault StoreFaultHook) error {
	directory := filepath.Dir(path)
	temporary, err := os.CreateTemp(directory, ".completion-init-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()

	if err := temporary.Chmod(0o600); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err = temporary.Write(exact); err == nil {
		err = temporary.Sync()
	}
	closeErr := temporary.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}

	if fault != nil {
		if err := fault(CheckpointBeforeCompletionStoreLink); err != nil {
			return err
		}
	}
	if err := os.Link(temporaryPath, path); err != nil {
		return err
	}
	if fault != nil {
		if err := fault(CheckpointAfterCompletionStoreLink); err != nil {
			// The link already landed a complete, synced snapshot at path;
			// only its directory-entry durability is now uncertain. Leave
			// the file in place — a caller can inspect or reopen it — but
			// do not report success for this call.
			return classified(ClassificationRecoveryRequired, "completion-store-create-indeterminate")
		}
	}
	if err := syncDirectory(directory); err != nil {
		return classified(ClassificationRecoveryRequired, "completion-store-create-indeterminate")
	}
	return nil
}

func OpenFixedFileCompletionStore(path string, fault StoreFaultHook) (*FixedFileCompletionStore, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0o077 != 0 {
		return nil, classified(ClassificationRecoveryRequired, "completion-store-file-policy")
	}
	file, err := os.Open(path) //nolint:gosec // unwired store receives one Supervisor-owned fixed local path
	if err != nil {
		return nil, err
	}
	openedInfo, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	if !os.SameFile(info, openedInfo) {
		_ = file.Close()
		return nil, classified(ClassificationRecoveryRequired, "completion-store-file-replaced")
	}
	exact, readErr := io.ReadAll(io.LimitReader(file, MaximumDurableCompletionStoreBytes+1))
	closeErr := file.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if len(exact) > MaximumDurableCompletionStoreBytes {
		return nil, classified(ClassificationRecoveryRequired, "completion-store-cap")
	}
	snapshot, err := decodeStoreSnapshot(exact)
	if err != nil {
		return nil, err
	}
	return &FixedFileCompletionStore{path: path, snapshot: snapshot, fault: fault}, nil
}

func (store *FixedFileCompletionStore) Close() {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.fenced = true
}

func (store *FixedFileCompletionStore) CommitCompletion(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	build func(v0candidate.PositiveUInt53) (DurableCompletion, error),
) (DurableCompletion, bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return DurableCompletion{}, false, err
	}
	if store.fenced {
		return DurableCompletion{}, false, ErrCompletionCommitIndeterminate
	}
	if attemptID == (approvalattempt.AttemptID{}) || build == nil {
		return DurableCompletion{}, false, classified(ClassificationSchema, "completion-store-request")
	}
	for _, stored := range store.snapshot.Records {
		if stored.AttemptID == hex.EncodeToString(attemptID[:]) {
			record, err := decodeStoredCompletion(stored)
			return record, false, err
		}
	}
	if len(store.snapshot.Records) >= MaximumRetainedCompletions {
		return DurableCompletion{}, false, classified(ClassificationCapacity, "completion-record-capacity")
	}
	sequence := store.snapshot.LastCommitSequence + 1
	if sequence == 0 || sequence > v0candidate.MaxSafeInteger {
		return DurableCompletion{}, false, classified(ClassificationCapacity, "completion-sequence-capacity")
	}
	record, err := build(v0candidate.PositiveUInt53(sequence))
	if err != nil {
		return DurableCompletion{}, false, err
	}
	if err := record.Validate(); err != nil {
		return DurableCompletion{}, false, err
	}
	view := record.View()
	if view.Projection.AttemptID != attemptID || uint64(view.CommitSequence) != sequence ||
		hex.EncodeToString(view.Projection.InstallationID[:]) != store.snapshot.InstallationID {
		return DurableCompletion{}, false, classified(ClassificationBinding, "completion-store-built-record")
	}
	aggregate := len(view.Completion.View().ResultJSON)
	for _, retained := range store.snapshot.Records {
		decoded, decodeErr := base64.StdEncoding.Strict().DecodeString(retained.ResultJSON)
		if decodeErr != nil {
			return DurableCompletion{}, false, classified(ClassificationRecoveryRequired, "completion-store-retained-result")
		}
		aggregate += len(decoded)
	}
	if aggregate > MaximumRetainedResultJSONBytes {
		return DurableCompletion{}, false, classified(ClassificationCapacity, "completion-store-result-capacity")
	}
	stored, err := encodeStoredCompletion(record)
	if err != nil {
		return DurableCompletion{}, false, err
	}
	working := store.snapshot
	working.Records = append([]storedDurableCompletion(nil), store.snapshot.Records...)
	working.Records = append(working.Records, stored)
	working.Generation++
	working.LastCommitSequence = sequence
	if err := validateStoreSnapshot(working); err != nil {
		return DurableCompletion{}, false, err
	}
	exact, err := encodeStoreSnapshot(working)
	if err != nil {
		return DurableCompletion{}, false, err
	}
	if err := store.publish(exact); err != nil {
		return DurableCompletion{}, false, err
	}
	store.snapshot = working
	return record, true, nil
}

func (store *FixedFileCompletionStore) ReadDurableCompletion(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (DurableCompletion, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return DurableCompletion{}, err
	}
	if store.fenced {
		return DurableCompletion{}, ErrCompletionCommitIndeterminate
	}
	encodedID := hex.EncodeToString(attemptID[:])
	for _, stored := range store.snapshot.Records {
		if stored.AttemptID == encodedID {
			return decodeStoredCompletion(stored)
		}
	}
	return DurableCompletion{}, ErrDurableCompletionNotFound
}

func (store *FixedFileCompletionStore) ReadCommittedCompletion(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (CompletionRecord, error) {
	record, err := store.ReadDurableCompletion(ctx, attemptID)
	if errors.Is(err, ErrDurableCompletionNotFound) {
		return CompletionRecord{}, ErrTerminalFactMissing
	}
	if err != nil {
		return CompletionRecord{}, err
	}
	return record.CompletionRecord(), nil
}

func (store *FixedFileCompletionStore) publish(exact []byte) error {
	directory := filepath.Dir(store.path)
	temporary, err := os.CreateTemp(directory, ".completion-next-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	remove := true
	defer func() {
		if remove {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(0o600); err != nil {
		_ = temporary.Close()
		return err
	}
	if _, err = temporary.Write(exact); err == nil {
		err = temporary.Sync()
	}
	closeErr := temporary.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	if store.fault != nil {
		if err := store.fault(CheckpointBeforeCompletionRename); err != nil {
			return err
		}
	}
	if err := os.Rename(temporaryPath, store.path); err != nil {
		return err
	}
	remove = false
	if store.fault != nil {
		if err := store.fault(CheckpointAfterCompletionRename); err != nil {
			store.fenced = true
			return ErrCompletionCommitIndeterminate
		}
	}
	if err := syncDirectory(directory); err != nil {
		store.fenced = true
		return ErrCompletionCommitIndeterminate
	}
	return nil
}

func syncDirectory(path string) error {
	directory, err := os.Open(path) //nolint:gosec // path is the fixed parent of the Supervisor-owned local store
	if err != nil {
		return err
	}
	err = directory.Sync()
	closeErr := directory.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func encodeStoreSnapshot(snapshot durableStoreSnapshot) ([]byte, error) {
	if err := validateStoreSnapshot(snapshot); err != nil {
		return nil, err
	}
	exact, err := json.Marshal(snapshot)
	if err != nil || len(exact) > MaximumDurableCompletionStoreBytes {
		return nil, classified(ClassificationCapacity, "completion-store-encoded-cap")
	}
	return exact, nil
}

func decodeStoreSnapshot(exact []byte) (durableStoreSnapshot, error) {
	// encoding/json's Decoder accepts syntactically ambiguous duplicate
	// object member names with last-value-wins behavior; reject them before
	// struct-decoding so a corrupted/tampered snapshot with a conflicting
	// earlier member cannot silently reopen as if it were unambiguous.
	if err := rejectDuplicateJSONNames(exact); err != nil {
		return durableStoreSnapshot{}, classified(ClassificationRecoveryRequired, "completion-store-duplicate-name")
	}
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.DisallowUnknownFields()
	var snapshot durableStoreSnapshot
	if err := decoder.Decode(&snapshot); err != nil {
		return durableStoreSnapshot{}, classified(ClassificationRecoveryRequired, "completion-store-malformed")
	}
	if token, err := decoder.Token(); err != io.EOF || token != nil {
		return durableStoreSnapshot{}, classified(ClassificationRecoveryRequired, "completion-store-trailing")
	}
	if err := validateStoreSnapshot(snapshot); err != nil {
		return durableStoreSnapshot{}, err
	}
	return snapshot, nil
}

// rejectDuplicateJSONNames walks exact and fails if any JSON object in it
// (at any depth) repeats a member name. It mirrors
// registrationstate.rejectDuplicateJSONNames for this package's own strict
// decode boundary.
func rejectDuplicateJSONNames(exact []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(exact))
	var readValue func() error
	readValue = func() error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		delimiter, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delimiter {
		case '{':
			seen := make(map[string]struct{})
			for decoder.More() {
				nameToken, err := decoder.Token()
				if err != nil {
					return err
				}
				name, ok := nameToken.(string)
				if !ok {
					return errors.New("completion store has a non-string object name")
				}
				if _, exists := seen[name]; exists {
					return errors.New("completion store has a duplicate object name")
				}
				seen[name] = struct{}{}
				if err := readValue(); err != nil {
					return err
				}
			}
			closing, err := decoder.Token()
			if err != nil || closing != json.Delim('}') {
				return errors.New("completion store has an incomplete object")
			}
			return nil
		case '[':
			for decoder.More() {
				if err := readValue(); err != nil {
					return err
				}
			}
			closing, err := decoder.Token()
			if err != nil || closing != json.Delim(']') {
				return errors.New("completion store has an incomplete array")
			}
			return nil
		default:
			return errors.New("completion store has an unexpected delimiter")
		}
	}
	if err := readValue(); err != nil {
		return err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return errors.New("completion store has trailing data")
	}
	return nil
}

func validateStoreSnapshot(snapshot durableStoreSnapshot) error {
	if snapshot.ObjectType != durableCompletionStoreType || snapshot.ObjectVersion != DurableCompletionStoreVersion {
		return classified(ClassificationRecoveryRequired, "completion-store-version")
	}
	installationID, err := decode16[v0candidate.InstallationID](snapshot.InstallationID)
	if err != nil || installationID == (v0candidate.InstallationID{}) {
		return classified(ClassificationRecoveryRequired, "completion-store-installation")
	}
	if len(snapshot.Records) > MaximumRetainedCompletions || snapshot.Generation != uint64(len(snapshot.Records)) ||
		snapshot.LastCommitSequence != uint64(len(snapshot.Records)) {
		return classified(ClassificationRecoveryRequired, "completion-store-counts")
	}
	seen := make(map[string]struct{}, len(snapshot.Records))
	var aggregate int
	for index, stored := range snapshot.Records {
		if stored.CommitSequence != uint64(index+1) {
			return classified(ClassificationRecoveryRequired, "completion-store-sequence")
		}
		if _, exists := seen[stored.AttemptID]; exists {
			return classified(ClassificationRecoveryRequired, "completion-store-duplicate")
		}
		seen[stored.AttemptID] = struct{}{}
		record, err := decodeStoredCompletion(stored)
		if err != nil {
			return classified(ClassificationRecoveryRequired, "completion-store-record")
		}
		aggregate += len(record.CompletionRecord().View().ResultJSON)
		if record.View().Projection.InstallationID != installationID {
			return classified(ClassificationRecoveryRequired, "completion-store-record-installation")
		}
		if aggregate > MaximumRetainedResultJSONBytes {
			return classified(ClassificationRecoveryRequired, "completion-store-result-cap")
		}
	}
	return nil
}

func encodeStoredCompletion(record DurableCompletion) (storedDurableCompletion, error) {
	view := record.View()
	projection := view.Projection
	stored := storedDurableCompletion{
		ObjectType: durableCompletionObjectType, ObjectVersion: view.ObjectVersion,
		CommitSequence: uint64(view.CommitSequence), Commit: string(view.Commit),
		ResultIntegrity: string(view.ResultIntegrity), RunnerIdentity: string(view.RunnerIdentity),
		Teardown: string(view.Teardown), DescendantAbsence: string(view.DescendantAbsence), Terminal: string(view.Terminal),
		AttemptID: encodeHex(projection.AttemptID[:]), ApprovalID: encodeHex(projection.ApprovalID[:]),
		RegistrationID: encodeHex(projection.RegistrationID[:]), PlanDigest: encodeHex(projection.PlanDigest[:]),
		InstallationID: encodeHex(projection.InstallationID[:]), EpochDigest: encodeHex(projection.EpochDigest[:]),
		SupervisorID: encodeHex(projection.SupervisorID[:]), ApprovalDigest: encodeHex(projection.ApprovalPayloadDigest[:]),
		Authorization: encodeHex(projection.AuthorizationIdentity[:]), SourceManifest: encodeHex(projection.SourceManifestDigest[:]),
		RuntimeDigest: encodeHex(projection.RuntimeBundleManifestDigest[:]),
		ProfileDigest: encodeHex(projection.ProfileRegistryEntryDigest[:]), BackendValidation: encodeHex(projection.BackendValidationRecordDigest[:]),
		BackendConfiguration: encodeHex(projection.BackendConfigurationDigest[:]), BackendImplementation: encodeHex(projection.BackendImplementationDigest[:]),
		BackendInstance: encodeHex(projection.BackendInstanceDigest[:]), ImmutableBinding: encodeHex(projection.ImmutableBindingDigest[:]),
		CompletionDigest: encodeHex(projection.CompletionRecordDigest[:]), ResultJSON: base64.StdEncoding.EncodeToString(projection.ResultJSON),
		ResultJSONDigest: encodeHex(projection.ResultJSONDigest[:]), LifecycleState: string(projection.LifecycleState),
		LifecycleSequence: uint64(projection.LifecycleOperationSequence), Reconciliation: string(projection.LifecycleLastReconciliation),
		CleanupRequired: projection.CleanupRequired, Transcript: base64.StdEncoding.EncodeToString(view.Transcript),
		TranscriptDigest: encodeHex(view.TranscriptDigest[:]), PublicSummary: base64.StdEncoding.EncodeToString(view.PublicSummary),
		DurableRecordDigest: encodeHex(view.DurableRecordDigest[:]),
	}
	exact, err := json.Marshal(stored)
	if err != nil || len(exact) > MaximumDurableCompletionBytes {
		return storedDurableCompletion{}, classified(ClassificationCapacity, "durable-completion-encoded-cap")
	}
	return stored, nil
}

func decodeStoredCompletion(stored storedDurableCompletion) (DurableCompletion, error) {
	exact, marshalErr := json.Marshal(stored)
	if marshalErr != nil || len(exact) > MaximumDurableCompletionBytes {
		return DurableCompletion{}, classified(ClassificationRecoveryRequired, "durable-completion-encoded-cap")
	}
	if stored.ObjectType != durableCompletionObjectType || stored.ObjectVersion != DurableCompletionObjectVersion {
		return DurableCompletion{}, classified(ClassificationRecoveryRequired, "durable-completion-object-version")
	}
	resultJSON, err := base64.StdEncoding.Strict().DecodeString(stored.ResultJSON)
	if err != nil {
		return DurableCompletion{}, classified(ClassificationRecoveryRequired, "durable-completion-result-encoding")
	}
	transcript, err := base64.StdEncoding.Strict().DecodeString(stored.Transcript)
	if err != nil {
		return DurableCompletion{}, classified(ClassificationRecoveryRequired, "durable-completion-transcript-encoding")
	}
	summary, err := base64.StdEncoding.Strict().DecodeString(stored.PublicSummary)
	if err != nil {
		return DurableCompletion{}, classified(ClassificationRecoveryRequired, "durable-completion-summary-encoding")
	}
	var projection TerminalProjectionView
	if err := decodeProjection(stored, &projection); err != nil {
		return DurableCompletion{}, err
	}
	completion, err := NewCompletionRecord(CompletionRecordView{
		RecordVersion: CompletionRecordVersion, CommitSequence: v0candidate.PositiveUInt53(stored.CommitSequence),
		AttemptID: projection.AttemptID, RegistrationID: projection.RegistrationID, PlanDigest: projection.PlanDigest,
		ImmutableBindingDigest: projection.ImmutableBindingDigest, Kind: CompletionFakeNoGuest,
		Status: CompletionSucceeded, Commit: CompletionCommit(stored.Commit), ResultJSON: resultJSON,
		ResultJSONByteLength: v0candidate.UInt53(len(resultJSON)), ResultJSONDigest: projection.ResultJSONDigest,
	})
	if err != nil {
		return DurableCompletion{}, classified(ClassificationRecoveryRequired, "durable-completion-record-decode")
	}
	storedCompletionDigest, err := decode32[CompletionRecordDigest](stored.CompletionDigest)
	if err != nil || storedCompletionDigest != completion.Digest() {
		return DurableCompletion{}, classified(ClassificationRecoveryRequired, "durable-completion-record-digest")
	}
	projection.ResultJSON = bytes.Clone(resultJSON)
	projection.CompletionRecordDigest = storedCompletionDigest
	transcriptDigest, err := decode32[TranscriptDigest](stored.TranscriptDigest)
	if err != nil {
		return DurableCompletion{}, err
	}
	durableDigest, err := decode32[DurableCompletionDigest](stored.DurableRecordDigest)
	if err != nil {
		return DurableCompletion{}, err
	}
	return restoreDurableCompletion(DurableCompletionView{
		ObjectVersion: stored.ObjectVersion, CommitSequence: v0candidate.PositiveUInt53(stored.CommitSequence),
		Commit: CompletionCommit(stored.Commit), ResultIntegrity: ResultIntegrityDisposition(stored.ResultIntegrity),
		RunnerIdentity: RunnerIdentityDisposition(stored.RunnerIdentity), Teardown: TeardownDisposition(stored.Teardown),
		DescendantAbsence: DescendantAbsenceDisposition(stored.DescendantAbsence), Terminal: TerminalClassification(stored.Terminal),
		Projection: projection, Completion: completion, Transcript: transcript, TranscriptDigest: transcriptDigest,
		PublicSummary: summary, DurableRecordDigest: durableDigest,
	})
}

func decodeProjection(stored storedDurableCompletion, target *TerminalProjectionView) error {
	var err error
	if target.AttemptID, err = decode16[approvalattempt.AttemptID](stored.AttemptID); err != nil {
		return err
	}
	if target.ApprovalID, err = decode16[approvalattempt.ApprovalID](stored.ApprovalID); err != nil {
		return err
	}
	if target.RegistrationID, err = decode16[v0candidate.RegistrationID](stored.RegistrationID); err != nil {
		return err
	}
	if target.PlanDigest, err = decode32[v0candidate.ExecutionPlanDigest](stored.PlanDigest); err != nil {
		return err
	}
	if target.InstallationID, err = decode16[v0candidate.InstallationID](stored.InstallationID); err != nil {
		return err
	}
	if target.EpochDigest, err = decode32[v0candidate.TrustEpochDigest](stored.EpochDigest); err != nil {
		return err
	}
	if target.SupervisorID, err = decode16[v0candidate.SupervisorID](stored.SupervisorID); err != nil {
		return err
	}
	if target.ApprovalPayloadDigest, err = decode32[approvalattempt.ApprovalPayloadDigest](stored.ApprovalDigest); err != nil {
		return err
	}
	if target.AuthorizationIdentity, err = decode32[approvalattempt.ApprovalKeyAuthorizationIdentity](stored.Authorization); err != nil {
		return err
	}
	if target.SourceManifestDigest, err = decode32[v0candidate.SourceManifestDigest](stored.SourceManifest); err != nil {
		return err
	}
	if target.RuntimeBundleManifestDigest, err = decode32[v0candidate.RuntimeBundleManifestDigest](stored.RuntimeDigest); err != nil {
		return err
	}
	if target.ProfileRegistryEntryDigest, err = decode32[v0candidate.ProfileRegistryEntryDigest](stored.ProfileDigest); err != nil {
		return err
	}
	if target.BackendValidationRecordDigest, err = decode32[v0candidate.BackendValidationRecordDigest](stored.BackendValidation); err != nil {
		return err
	}
	if target.BackendConfigurationDigest, err = decode32[v0candidate.BackendConfigurationDigest](stored.BackendConfiguration); err != nil {
		return err
	}
	if target.BackendImplementationDigest, err = decode32[lifecyclestate.BackendImplementationDigest](stored.BackendImplementation); err != nil {
		return err
	}
	if target.BackendInstanceDigest, err = decode32[lifecyclestate.BackendInstanceDigest](stored.BackendInstance); err != nil {
		return err
	}
	if target.ImmutableBindingDigest, err = decode32[lifecyclestate.ImmutableBindingDigest](stored.ImmutableBinding); err != nil {
		return err
	}
	if target.ResultJSONDigest, err = decode32[ResultJSONDigest](stored.ResultJSONDigest); err != nil {
		return err
	}
	target.LifecycleState = lifecyclestate.LifecycleState(stored.LifecycleState)
	target.LifecycleOperationSequence = v0candidate.PositiveUInt53(stored.LifecycleSequence)
	target.LifecycleLastReconciliation = lifecyclestate.ReconciliationStatus(stored.Reconciliation)
	target.CleanupRequired = stored.CleanupRequired
	return nil
}

func encodeHex(value []byte) string { return hex.EncodeToString(value) }

func decode16[T ~[16]byte](value string) (T, error) {
	var result T
	exact, err := hex.DecodeString(value)
	if err != nil || len(exact) != len(result) {
		return result, classified(ClassificationRecoveryRequired, "durable-completion-identity")
	}
	copy(result[:], exact)
	if result == (T{}) {
		return result, classified(ClassificationRecoveryRequired, "durable-completion-zero-identity")
	}
	return result, nil
}

func decode32[T ~[32]byte](value string) (T, error) {
	var result T
	exact, err := hex.DecodeString(value)
	if err != nil || len(exact) != len(result) {
		return result, classified(ClassificationRecoveryRequired, "durable-completion-digest-field")
	}
	copy(result[:], exact)
	if result == (T{}) {
		return result, classified(ClassificationRecoveryRequired, "durable-completion-zero-digest")
	}
	return result, nil
}

var _ DurableCompletionStore = (*FixedFileCompletionStore)(nil)
var _ CommittedCompletionReader = (*FixedFileCompletionStore)(nil)
