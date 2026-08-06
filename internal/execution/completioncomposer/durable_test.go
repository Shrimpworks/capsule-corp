package completioncomposer

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"capsule.local/capsule/internal/execution/lifecyclestate"
)

type durableKnownAnswer struct {
	DurableRecordSHA256 string `json:"durableRecordSha256"`
	TranscriptSHA256    string `json:"transcriptSha256"`
	PublicSummaryJSON   string `json:"publicSummaryJson"`
}

func TestDurableProducerKnownAnswerReopenAndImmutableReplay(t *testing.T) {
	fixture := readKnownAnswer(t)
	known := readDurableKnownAnswer(t)
	source := newRetainedTruthSource(t, []byte(fixture.ResultJSON))
	path := filepath.Join(t.TempDir(), "completion-store.json")
	store, err := CreateFixedFileCompletionStore(path, source.created.Attempt.InstallationID, nil)
	if err != nil {
		t.Fatal(err)
	}
	producer := mustProducer(t, source, store)
	request := commitRequest(source, []byte(fixture.ResultJSON))
	first, err := producer.Complete(context.Background(), request)
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
	durableSHA256 := first.Digest()
	transcriptSHA256 := first.TranscriptSHA256()
	durableDigest := hex.EncodeToString(durableSHA256[:])
	transcriptDigest := hex.EncodeToString(transcriptSHA256[:])
	if known.DurableRecordSHA256 == "" || known.TranscriptSHA256 == "" || known.PublicSummaryJSON == "" {
		t.Fatalf("durable known answer incomplete\ndurable=%s\ntranscript=%s\nsummary=%s", durableDigest, transcriptDigest, first.PublicSummaryBytes())
	}
	if durableDigest != known.DurableRecordSHA256 || transcriptDigest != known.TranscriptSHA256 ||
		string(first.PublicSummaryBytes()) != known.PublicSummaryJSON {
		t.Fatalf("durable known answer mismatch\ndurable=%s\ntranscript=%s\nsummary=%s", durableDigest, transcriptDigest, first.PublicSummaryBytes())
	}
	store.Close()
	reopened, err := OpenFixedFileCompletionStore(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	replayed, err := mustProducer(t, source, reopened).Complete(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if replayed.Digest() != first.Digest() || !bytes.Equal(replayed.TranscriptBytes(), first.TranscriptBytes()) ||
		!bytes.Equal(replayed.PublicSummaryBytes(), first.PublicSummaryBytes()) {
		t.Fatal("same AttemptID replay changed public result")
	}
	stored, err := reopened.ReadCommittedCompletion(context.Background(), source.created.Attempt.AttemptID)
	if err != nil || stored.Digest() != first.CompletionRecord().Digest() {
		t.Fatalf("committed reader: %v", err)
	}
	for _, forbidden := range []string{"path", "artifactName", "diagnostic", "exitCode", "eof", "startedAt", "finishedAt"} {
		if strings.Contains(string(first.TranscriptBytes()), forbidden) || strings.Contains(string(first.PublicSummaryBytes()), forbidden) {
			t.Fatalf("forbidden public field %q", forbidden)
		}
	}
	view := first.View()
	if view.RunnerIdentity != RunnerIdentityUnresolvedFake || view.Teardown != TeardownFakeDestroyConfirmed ||
		view.DescendantAbsence != DescendantsFakeAbsent || view.Terminal != TerminalCompletedFakeLocal {
		t.Fatalf("fixed dispositions = %+v", view)
	}
}

func TestDurableProducerCapsMalformedAndMixedVersions(t *testing.T) {
	source := newRetainedTruthSource(t, []byte("null"))
	producer, _ := newProducerAndStore(t, source, nil)
	over := append(append([]byte{'"'}, bytes.Repeat([]byte{'a'}, MaximumResultJSONBytes-1)...), '"')
	request := commitRequest(source, over)
	if _, err := producer.Complete(context.Background(), request); classification(t, err) != ClassificationMalformed {
		t.Fatalf("cap+1 = %v", err)
	}
	for _, exact := range [][]byte{[]byte(`{"x":1,"x":2}`), []byte(`1e2`), []byte(`null true`)} {
		request = commitRequest(source, exact)
		if _, err := producer.Complete(context.Background(), request); classification(t, err) != ClassificationMalformed {
			t.Fatalf("malformed %q = %v", exact, err)
		}
	}
	request = commitRequest(source, []byte("null"))
	request.ObjectVersion++
	if _, err := producer.Complete(context.Background(), request); classification(t, err) != ClassificationUnsupported {
		t.Fatalf("mixed request version = %v", err)
	}

	path := filepath.Join(t.TempDir(), "mixed-store.json")
	store, err := CreateFixedFileCompletionStore(path, source.created.Attempt.InstallationID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := mustProducer(t, source, store).Complete(context.Background(), commitRequest(source, []byte("null"))); err != nil {
		t.Fatal(err)
	}
	snapshot := readSnapshotForTest(t, path)
	snapshot.Records[0].ObjectVersion++
	writeSnapshotForTest(t, path, snapshot)
	if _, err := OpenFixedFileCompletionStore(path, nil); classification(t, err) != ClassificationRecoveryRequired {
		t.Fatalf("mixed stored version = %v", err)
	}
}

func TestDurableStoreRefusesEarlyMissingDuplicateAndForgedCompletion(t *testing.T) {
	mutations := []struct {
		name  string
		alter func(*durableStoreSnapshot)
	}{
		{name: "early-commit", alter: func(value *durableStoreSnapshot) { value.Records[0].Commit = "prepared" }},
		{name: "missing-commit", alter: func(value *durableStoreSnapshot) { value.Records[0].Commit = "" }},
		{name: "duplicate-commit", alter: func(value *durableStoreSnapshot) {
			duplicate := value.Records[0]
			duplicate.CommitSequence = 2
			value.Records = append(value.Records, duplicate)
			value.Generation = 2
			value.LastCommitSequence = 2
		}},
		{name: "forged-attempt", alter: func(value *durableStoreSnapshot) { value.Records[0].AttemptID = strings.Repeat("fe", 16) }},
	}
	for _, testCase := range mutations {
		t.Run(testCase.name, func(t *testing.T) {
			source := newRetainedTruthSource(t, []byte("null"))
			path := filepath.Join(t.TempDir(), "store.json")
			store, err := CreateFixedFileCompletionStore(path, source.created.Attempt.InstallationID, nil)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := mustProducer(t, source, store).Complete(context.Background(), commitRequest(source, []byte("null"))); err != nil {
				t.Fatal(err)
			}
			snapshot := readSnapshotForTest(t, path)
			testCase.alter(&snapshot)
			writeSnapshotForTest(t, path, snapshot)
			if _, err := OpenFixedFileCompletionStore(path, nil); classification(t, err) != ClassificationRecoveryRequired {
				t.Fatalf("classification = %v", err)
			}
		})
	}
}

func TestEOFRunnerExitZeroAndUnresolvedTeardownNeverCommit(t *testing.T) {
	for _, testCase := range []struct {
		name  string
		alter func(*retainedTruthSource)
	}{
		{name: "eof-and-runner-exit-zero-only", alter: func(source *retainedTruthSource) {
			source.lifecycle = observedLifecycle(t, source.lifecycle)
		}},
		{name: "unresolved-teardown", alter: func(source *retainedTruthSource) {
			source.lifecycle = unresolvedLifecycle(t, source.lifecycle)
		}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			source := newRetainedTruthSource(t, []byte("null"))
			testCase.alter(source)
			producer, store := newProducerAndStore(t, source, nil)
			if _, err := producer.Complete(context.Background(), commitRequest(source, []byte("null"))); classification(t, err) != ClassificationNotTerminal {
				t.Fatalf("classification = %v", err)
			}
			if _, err := store.ReadDurableCompletion(context.Background(), source.created.Attempt.AttemptID); !errors.Is(err, ErrDurableCompletionNotFound) {
				t.Fatalf("nonterminal completion became durable: %v", err)
			}
		})
	}
}

func TestCrashEdgesResponseLossRestartAndStaleReplay(t *testing.T) {
	source := newRetainedTruthSource(t, []byte("null"))
	beforeFault := true
	beforeHook := func(checkpoint StoreCommitCheckpoint) error {
		if checkpoint == CheckpointBeforeCompletionRename && beforeFault {
			beforeFault = false
			return errors.New("injected before rename")
		}
		return nil
	}
	producer, store := newProducerAndStore(t, source, beforeHook)
	request := commitRequest(source, []byte("null"))
	if _, err := producer.Complete(context.Background(), request); classification(t, err) != ClassificationLocalFailure {
		t.Fatalf("before-rename fault = %v", err)
	}
	if _, err := store.ReadDurableCompletion(context.Background(), request.AttemptID); !errors.Is(err, ErrDurableCompletionNotFound) {
		t.Fatalf("pre-rename crash committed: %v", err)
	}

	afterFault := true
	path := filepath.Join(t.TempDir(), "indeterminate.json")
	indeterminate, err := CreateFixedFileCompletionStore(path, source.created.Attempt.InstallationID, func(checkpoint StoreCommitCheckpoint) error {
		if checkpoint == CheckpointAfterCompletionRename && afterFault {
			afterFault = false
			return errors.New("injected response loss")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := mustProducer(t, source, indeterminate).Complete(context.Background(), request); classification(t, err) != ClassificationRecoveryRequired {
		t.Fatalf("after-rename fault = %v", err)
	}
	if _, err := indeterminate.ReadDurableCompletion(context.Background(), request.AttemptID); !errors.Is(err, ErrCompletionCommitIndeterminate) {
		t.Fatalf("indeterminate store not fenced: %v", err)
	}
	reopened, err := OpenFixedFileCompletionStore(path, nil)
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := mustProducer(t, source, reopened).Complete(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	stale := commitRequest(source, []byte(`{"different":true}`))
	if _, err := mustProducer(t, source, reopened).Complete(context.Background(), stale); classification(t, err) != ClassificationReplay {
		t.Fatalf("stale replay = %v", err)
	}
	readback, err := reopened.ReadDurableCompletion(context.Background(), request.AttemptID)
	if err != nil || readback.Digest() != recovered.Digest() {
		t.Fatalf("stale replay changed result: %v", err)
	}
}

func TestDurableCompletionDefensiveCopies(t *testing.T) {
	source := newRetainedTruthSource(t, []byte("null"))
	producer, _ := newProducerAndStore(t, source, nil)
	input := []byte(`{"ok":true}`)
	request := commitRequest(source, input)
	record, err := producer.Complete(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	input[2] = 'X'
	view := record.View()
	view.Projection.ResultJSON[0] = 'X'
	view.Transcript[0] = 'X'
	view.PublicSummary[0] = 'X'
	if record.View().Projection.ResultJSON[0] == 'X' || record.TranscriptBytes()[0] == 'X' || record.PublicSummaryBytes()[0] == 'X' {
		t.Fatal("durable completion exposed aliased bytes")
	}
}

func TestDurableStoreCreateIsAtomicAndExclusive(t *testing.T) {
	source := newRetainedTruthSource(t, []byte("null"))
	directory := t.TempDir()
	path := filepath.Join(directory, "store.json")
	if _, err := CreateFixedFileCompletionStore(path, source.created.Attempt.InstallationID, nil); err != nil {
		t.Fatal(err)
	}
	before := readSnapshotForTest(t, path)

	// A second create at the same path must not silently replace the first
	// store's committed content (ADR-0042's atomic-publish invariant applies
	// to creation too, not only later commits).
	if _, err := CreateFixedFileCompletionStore(path, source.created.Attempt.InstallationID, nil); err == nil {
		t.Fatal("create did not refuse an already-occupied path")
	}
	if after := readSnapshotForTest(t, path); after.InstallationID != before.InstallationID {
		t.Fatal("refused create mutated the existing store")
	}

	// No leftover temporary file survives a successful create.
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.Name() != filepath.Base(path) {
			t.Fatalf("create left a stray file behind: %s", entry.Name())
		}
	}

	// A create whose write fails before the atomic link must not leave a
	// partial file at the final path: a retry at the same path must still
	// succeed as a fresh creation.
	failingPath := filepath.Join(directory, "unwritable", "store.json")
	if err := os.MkdirAll(filepath.Dir(failingPath), 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(filepath.Dir(failingPath), 0o700) })
	if _, err := CreateFixedFileCompletionStore(failingPath, source.created.Attempt.InstallationID, nil); err == nil {
		t.Fatal("create in a read-only directory unexpectedly succeeded")
	}
	if _, err := os.Lstat(failingPath); !os.IsNotExist(err) {
		t.Fatalf("failed create left a partial file: %v", err)
	}
}

func TestDurableStoreRefusesDuplicateObjectName(t *testing.T) {
	source := newRetainedTruthSource(t, []byte("null"))
	path := filepath.Join(t.TempDir(), "store.json")
	store, err := CreateFixedFileCompletionStore(path, source.created.Attempt.InstallationID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := mustProducer(t, source, store).Complete(context.Background(), commitRequest(source, []byte("null"))); err != nil {
		t.Fatal(err)
	}

	exact, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	// Insert an earlier conflicting objectType member. Go's json.Decoder
	// otherwise accepts this with last-value-wins behavior; the store must
	// treat it as malformed recovery state instead of silently reopening.
	tampered := bytes.Replace(exact, []byte(`"objectType"`),
		[]byte(`"objectType":"forged","objectType"`), 1)
	if bytes.Equal(tampered, exact) {
		t.Fatal("test setup did not tamper the snapshot")
	}
	if err := os.WriteFile(path, tampered, 0o600); err != nil { //nolint:gosec // G703: path is the t.TempDir() store file created earlier in this test.
		t.Fatal(err)
	}
	if _, err := OpenFixedFileCompletionStore(path, nil); classification(t, err) != ClassificationRecoveryRequired {
		t.Fatalf("duplicate object name = %v", err)
	}
}

func commitRequest(source *retainedTruthSource, exact []byte) CommitRequest {
	return CommitRequest{ObjectVersion: DurableCompletionObjectVersion, AttemptID: source.created.Attempt.AttemptID,
		ResultJSON: exact, ResultIntegrity: ResultIntegrityTypedJSONSHA256}
}

func newProducerAndStore(t *testing.T, source *retainedTruthSource, hook StoreFaultHook) (*Producer, *FixedFileCompletionStore) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "completion-store.json")
	store, err := CreateFixedFileCompletionStore(path, source.created.Attempt.InstallationID, hook)
	if err != nil {
		t.Fatal(err)
	}
	return mustProducer(t, source, store), store
}

func mustProducer(t *testing.T, source *retainedTruthSource, store DurableCompletionStore) *Producer {
	t.Helper()
	producer, err := NewProducer(ProducerSources{Attempts: source, Lifecycles: source, Store: store})
	if err != nil {
		t.Fatal(err)
	}
	return producer
}

func unresolvedLifecycle(t *testing.T, terminal lifecyclestate.Record) lifecyclestate.Record {
	t.Helper()
	view := terminal.View()
	view.RecordVersion--
	view.OperationSequence = 4
	view.Operation = lifecyclestate.OperationObserve
	view.EffectStatus = lifecyclestate.EffectIndeterminate
	view.State = lifecyclestate.StateUnresolved
	view.CleanupRequired = true
	view.LastConfirmedCheckpoint = lifecyclestate.CheckpointStart
	view.FirstFailure = lifecyclestate.FailureCleanupUnresolved
	view.FailureOperation = lifecyclestate.OperationObserve
	view.LastReconciliation = lifecyclestate.ReconciliationUnknown
	view.AutomaticRecoveryCount = 1
	view.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{Present: true, Value: 125}
	view.RecoveryFence = lifecyclestate.RecoveryFenceReconcileUnknown
	view.TerminalAt = lifecyclestate.OptionalUnixSeconds{}
	record, err := lifecyclestate.NewRecord(view, 130)
	if err != nil {
		t.Fatalf("new unresolved lifecycle: %v", err)
	}
	return record
}

func readSnapshotForTest(t *testing.T, path string) durableStoreSnapshot {
	t.Helper()
	exact, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var snapshot durableStoreSnapshot
	if err := json.Unmarshal(exact, &snapshot); err != nil {
		t.Fatal(err)
	}
	return snapshot
}

func writeSnapshotForTest(t *testing.T, path string, snapshot durableStoreSnapshot) {
	t.Helper()
	exact, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, exact, 0o600); err != nil {
		t.Fatal(err)
	}
}

func readDurableKnownAnswer(t *testing.T) durableKnownAnswer {
	t.Helper()
	exact, err := os.ReadFile(filepath.Join("testdata", "durable-known-answer.json"))
	if err != nil {
		t.Fatal(err)
	}
	var result durableKnownAnswer
	if err := json.Unmarshal(exact, &result); err != nil {
		t.Fatal(err)
	}
	return result
}
