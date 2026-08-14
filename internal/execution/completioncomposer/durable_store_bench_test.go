package completioncomposer

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"path/filepath"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// attemptIDForBenchmarkIndex derives a distinct, deterministic AttemptID for
// benchmark record index. BenchmarkCommitCompletionFillToCapacity needs many
// independently committable records, one per AttemptID, without paying to
// build a full attempt/approval/registration/lifecycle fixture graph per
// record the way the package's other tests do through retainedTruthSource.
func attemptIDForBenchmarkIndex(index int) approvalattempt.AttemptID {
	var attemptID approvalattempt.AttemptID
	attemptID[0] = 0xb2
	binary.BigEndian.PutUint64(attemptID[8:], uint64(index)) //nolint:gosec // benchmark fixture identity only
	return attemptID
}

// buildBenchmarkDurableCompletion constructs a self-consistent
// DurableCompletion for benchmark record index at the given commit
// sequence, following the same shape Producer.build/Compositor.Compose
// produce, so it exercises FixedFileCompletionStore.CommitCompletion's real
// encode/validate/publish path rather than a shortcut.
func buildBenchmarkDurableCompletion(
	installationID v0candidate.InstallationID,
	index int,
	sequence v0candidate.PositiveUInt53,
) (DurableCompletion, error) {
	attemptID := attemptIDForBenchmarkIndex(index)
	resultJSON := []byte(fmt.Sprintf(`{"benchmarkIndex":%d}`, index))
	projection := TerminalProjectionView{
		AttemptID: attemptID, ApprovalID: repeated16[approvalattempt.ApprovalID](0xa1),
		RegistrationID: repeated16[v0candidate.RegistrationID](0x33),
		PlanDigest:     repeated32[v0candidate.ExecutionPlanDigest](0x44),
		InstallationID: installationID, EpochDigest: repeated32[v0candidate.TrustEpochDigest](0x22),
		SupervisorID:                  repeated16[v0candidate.SupervisorID](0x55),
		ApprovalPayloadDigest:         repeated32[approvalattempt.ApprovalPayloadDigest](0xd4),
		AuthorizationIdentity:         repeated32[approvalattempt.ApprovalKeyAuthorizationIdentity](0x99),
		SourceManifestDigest:          repeated32[v0candidate.SourceManifestDigest](0xc3),
		RuntimeBundleManifestDigest:   repeated32[v0candidate.RuntimeBundleManifestDigest](0x55),
		ProfileRegistryEntryDigest:    repeated32[v0candidate.ProfileRegistryEntryDigest](0x77),
		BackendValidationRecordDigest: repeated32[v0candidate.BackendValidationRecordDigest](0x88),
		BackendConfigurationDigest:    repeated32[v0candidate.BackendConfigurationDigest](0x99),
		BackendImplementationDigest:   repeated32[lifecyclestate.BackendImplementationDigest](0xe1),
		BackendInstanceDigest:         repeated32[lifecyclestate.BackendInstanceDigest](0xf0),
		ImmutableBindingDigest:        repeated32[lifecyclestate.ImmutableBindingDigest](0x12),
		ResultJSON:                    resultJSON, ResultJSONDigest: ResultJSONDigest(sha256.Sum256(resultJSON)),
		LifecycleState: lifecyclestate.StateDestroyed, LifecycleOperationSequence: 6,
		LifecycleLastReconciliation: lifecyclestate.ReconciliationAuthoritativelyAbsent,
	}
	completion, err := NewCompletionRecord(CompletionRecordView{
		RecordVersion: CompletionRecordVersion, CommitSequence: sequence,
		AttemptID: attemptID, RegistrationID: projection.RegistrationID, PlanDigest: projection.PlanDigest,
		ImmutableBindingDigest: projection.ImmutableBindingDigest, Kind: CompletionFakeNoGuest,
		Status: CompletionSucceeded, Commit: CompletionCommittedLast, ResultJSON: resultJSON,
		ResultJSONByteLength: v0candidate.UInt53(len(resultJSON)), ResultJSONDigest: projection.ResultJSONDigest,
	})
	if err != nil {
		return DurableCompletion{}, err
	}
	projection.CompletionRecordDigest = completion.Digest()
	transcript, transcriptDigest, err := encodeTranscript(projection)
	if err != nil {
		return DurableCompletion{}, err
	}
	summary, err := encodePublicSummary(attemptID, transcriptDigest)
	if err != nil {
		return DurableCompletion{}, err
	}
	view := DurableCompletionView{
		ObjectVersion: DurableCompletionObjectVersion, CommitSequence: sequence,
		Commit: CompletionCommittedLast, ResultIntegrity: ResultIntegrityTypedJSONSHA256,
		RunnerIdentity: RunnerIdentityUnresolvedFake, Teardown: TeardownFakeDestroyConfirmed,
		DescendantAbsence: DescendantsFakeAbsent, Terminal: TerminalCompletedFakeLocal,
		Projection: projection, Completion: completion, Transcript: transcript,
		TranscriptDigest: transcriptDigest, PublicSummary: summary,
	}
	view.DurableRecordDigest = digestDurableCompletion(view)
	return restoreDurableCompletion(view)
}

// BenchmarkCommitCompletionFillToCapacity fills a fresh store to `records`
// commits and measures total wall time. Before the item-2 fix,
// validateStoreSnapshot re-decoded/re-marshaled every previously-committed
// record on every commit, so total fill time scaled quadratically with the
// target record count (roughly a 4x runtime increase for each doubling of
// records). After the fix, per-commit validation cost is independent of
// store size, so total fill time scales roughly linearly with record count
// (roughly 2x runtime for each doubling).
func BenchmarkCommitCompletionFillToCapacity(b *testing.B) {
	installationID := repeated16[v0candidate.InstallationID](0x11)
	for _, count := range []int{512, 1024, 2048, MaximumRetainedCompletions} {
		b.Run(fmt.Sprintf("records=%d", count), func(b *testing.B) {
			for iteration := 0; iteration < b.N; iteration++ {
				b.StopTimer()
				path := filepath.Join(b.TempDir(), "bench-store.json")
				store, err := CreateFixedFileCompletionStore(path, installationID, nil)
				if err != nil {
					b.Fatalf("create store: %v", err)
				}
				b.StartTimer()
				for index := 0; index < count; index++ {
					attemptID := attemptIDForBenchmarkIndex(index)
					_, _, err := store.CommitCompletion(context.Background(), attemptID,
						func(sequence v0candidate.PositiveUInt53) (DurableCompletion, error) {
							return buildBenchmarkDurableCompletion(installationID, index, sequence)
						})
					if err != nil {
						b.Fatalf("commit %d: %v", index, err)
					}
				}
			}
		})
	}
}
