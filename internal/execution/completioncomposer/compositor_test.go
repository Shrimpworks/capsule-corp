package completioncomposer

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

type knownAnswerFixture struct {
	ResultJSON        string `json:"resultJson"`
	TranscriptJSON    string `json:"transcriptJson"`
	TranscriptSHA256  string `json:"transcriptSha256"`
	PublicSummaryJSON string `json:"publicSummaryJson"`
}

type retainedTruthSource struct {
	created       registrationstate.CreatedAttempt
	lifecycle     lifecyclestate.Record
	completion    CompletionRecord
	attemptErr    error
	lifecycleErr  error
	completionErr error
}

func (source *retainedTruthSource) ResolveCreated(
	context.Context,
	approvalattempt.AttemptID,
) (registrationstate.CreatedAttempt, error) {
	if source.attemptErr != nil {
		return registrationstate.CreatedAttempt{}, source.attemptErr
	}
	return cloneCreatedAttempt(source.created), nil
}

func (source *retainedTruthSource) ReadLifecycle(
	context.Context,
	approvalattempt.AttemptID,
) (lifecyclestate.Record, error) {
	if source.lifecycleErr != nil {
		return lifecyclestate.Record{}, source.lifecycleErr
	}
	return source.lifecycle, nil
}

func (source *retainedTruthSource) ReadCommittedCompletion(
	context.Context,
	approvalattempt.AttemptID,
) (CompletionRecord, error) {
	if source.completionErr != nil {
		return CompletionRecord{}, source.completionErr
	}
	return source.completion, nil
}

func TestKnownAnswerTerminalProjection(t *testing.T) {
	fixture := readKnownAnswer(t)
	source := newRetainedTruthSource(t, []byte(fixture.ResultJSON))
	result := compose(t, source)
	transcript := string(result.TranscriptBytes())
	summary := string(result.PublicSummaryBytes())
	transcriptDigest := result.TranscriptDigest()
	digest := hex.EncodeToString(transcriptDigest[:])
	if fixture.TranscriptJSON == "" || fixture.TranscriptSHA256 == "" || fixture.PublicSummaryJSON == "" {
		t.Fatalf("known answer incomplete\ntranscript=%s\ndigest=%s\nsummary=%s", transcript, digest, summary)
	}
	if transcript != fixture.TranscriptJSON {
		t.Fatalf("transcript mismatch\n got: %s\nwant: %s", transcript, fixture.TranscriptJSON)
	}
	if digest != fixture.TranscriptSHA256 {
		t.Fatalf("transcript digest = %s, want %s", digest, fixture.TranscriptSHA256)
	}
	if summary != fixture.PublicSummaryJSON {
		t.Fatalf("public summary mismatch\n got: %s\nwant: %s", summary, fixture.PublicSummaryJSON)
	}
	view := result.Projection().View()
	if !bytes.Equal(view.ResultJSON, []byte(fixture.ResultJSON)) ||
		view.LifecycleState != lifecyclestate.StateDestroyed || view.CleanupRequired ||
		view.LifecycleLastReconciliation != lifecyclestate.ReconciliationAuthoritativelyAbsent {
		t.Fatalf("terminal projection = %+v", view)
	}
	for _, forbidden := range []string{"exitCode", "eof", "message", "diagnostic", "path", "artifact", "startedAt", "finishedAt"} {
		if strings.Contains(transcript, forbidden) || strings.Contains(summary, forbidden) {
			t.Fatalf("forbidden field %q entered public material", forbidden)
		}
	}
}

func TestResponseLossReopenAndReplayConverge(t *testing.T) {
	fixture := readKnownAnswer(t)
	firstSource := newRetainedTruthSource(t, []byte(fixture.ResultJSON))
	first := compose(t, firstSource)
	// Model response loss by discarding first result. Reopen restores the
	// committed completion record and rereads the other retained facts; the
	// compositor owns no state.
	reopenedSource := newRetainedTruthSource(t, []byte(fixture.ResultJSON))
	restored, err := RestoreCompletionRecord(
		firstSource.completion.View(),
		firstSource.completion.Digest(),
	)
	if err != nil {
		t.Fatalf("restore committed completion: %v", err)
	}
	reopenedSource.completion = restored
	reopened := compose(t, reopenedSource)
	if !bytes.Equal(first.TranscriptBytes(), reopened.TranscriptBytes()) ||
		first.TranscriptDigest() != reopened.TranscriptDigest() ||
		!bytes.Equal(first.PublicSummaryBytes(), reopened.PublicSummaryBytes()) {
		t.Fatal("response-loss reopen did not converge on exact terminal projection")
	}
	for index := 0; index < 8; index++ {
		replayed := compose(t, reopenedSource)
		if replayed.TranscriptDigest() != first.TranscriptDigest() ||
			!bytes.Equal(replayed.TranscriptBytes(), first.TranscriptBytes()) {
			t.Fatalf("replay %d changed transcript", index)
		}
	}
}

func TestCompletionJSONExactCapAndMaxPlusOne(t *testing.T) {
	source := newRetainedTruthSource(t, []byte("null"))
	binding := source.lifecycle.ImmutableBindingDigest()
	exact := append(append([]byte{'"'}, bytes.Repeat([]byte{'a'}, MaximumResultJSONBytes-2)...), '"')
	record := completionForJSON(t, source.created, binding, exact)
	if len(record.View().ResultJSON) != MaximumResultJSONBytes {
		t.Fatal("exact completion cap not retained")
	}
	over := append(append([]byte{'"'}, bytes.Repeat([]byte{'a'}, MaximumResultJSONBytes-1)...), '"')
	view := record.View()
	view.ResultJSON = over
	view.ResultJSONByteLength = v0candidate.UInt53(len(over))
	view.ResultJSONDigest = ResultJSONDigest(sha256.Sum256(over))
	if _, err := NewCompletionRecord(view); classification(t, err) != ClassificationMalformed {
		t.Fatalf("cap-plus-one classification = %v", err)
	}
}

func TestCompletionJSONRefusesMalformedAndUntypedValues(t *testing.T) {
	source := newRetainedTruthSource(t, []byte("null"))
	tests := []string{
		`{"a":1,"a":2}`,
		`1.5`,
		`1e2`,
		`9007199254740992`,
		`"\ud800"`,
		`null true`,
	}
	for _, exact := range tests {
		view := source.completion.View()
		view.ResultJSON = []byte(exact)
		view.ResultJSONByteLength = v0candidate.UInt53(len(exact))
		view.ResultJSONDigest = ResultJSONDigest(sha256.Sum256([]byte(exact)))
		if _, err := NewCompletionRecord(view); classification(t, err) != ClassificationMalformed {
			t.Errorf("%q classification = %v", exact, err)
		}
	}
}

func TestMissingCorruptAndContradictoryTruthRefuses(t *testing.T) {
	fixture := readKnownAnswer(t)
	tests := []struct {
		name  string
		alter func(*retainedTruthSource)
		want  Classification
	}{
		{name: "missing-attempt", alter: func(source *retainedTruthSource) { source.attemptErr = ErrTerminalFactMissing }, want: ClassificationNotTerminal},
		{name: "missing-completion", alter: func(source *retainedTruthSource) { source.completionErr = ErrTerminalFactMissing }, want: ClassificationNotTerminal},
		{name: "missing-lifecycle", alter: func(source *retainedTruthSource) { source.lifecycleErr = ErrTerminalFactMissing }, want: ClassificationNotTerminal},
		{name: "corrupt-source", alter: func(source *retainedTruthSource) { source.completionErr = ErrTerminalTruthCorrupt }, want: ClassificationRecoveryRequired},
		{name: "corrupt-completion", alter: func(source *retainedTruthSource) { source.completion.digest[0] ^= 0xff }, want: ClassificationRecoveryRequired},
		{name: "wrong-completion-attempt", alter: func(source *retainedTruthSource) {
			view := source.completion.View()
			view.AttemptID[0] ^= 0xff
			source.completion = mustCompletion(t, view)
		}, want: ClassificationBinding},
		{name: "nonterminal-lifecycle", alter: func(source *retainedTruthSource) { source.lifecycle = observedLifecycle(t, source.lifecycle) }, want: ClassificationNotTerminal},
		{name: "failed-lifecycle", alter: func(source *retainedTruthSource) { source.lifecycle = failedDestroyedLifecycle(t, source.lifecycle) }, want: ClassificationBinding},
		{name: "plan-result-cap", alter: func(source *retainedTruthSource) {
			exact := append(append([]byte{'"'}, bytes.Repeat([]byte{'a'}, 65_536)...), '"')
			source.completion = completionForJSON(t, source.created, source.lifecycle.ImmutableBindingDigest(), exact)
		}, want: ClassificationBinding},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			source := newRetainedTruthSource(t, []byte(fixture.ResultJSON))
			testCase.alter(source)
			compositor, err := New(Sources{Attempts: source, Lifecycles: source, Completions: source})
			if err != nil {
				t.Fatal(err)
			}
			_, err = compositor.Compose(context.Background(), source.created.Attempt.AttemptID)
			if got := classification(t, err); got != testCase.want {
				t.Fatalf("classification = %s (%v), want %s", got, err, testCase.want)
			}
		})
	}
}

func TestDefensiveCopies(t *testing.T) {
	input := []byte(`{"ok":true}`)
	source := newRetainedTruthSource(t, input)
	input[2] = 'X'
	if !bytes.Equal(source.completion.View().ResultJSON, []byte(`{"ok":true}`)) {
		t.Fatal("completion aliased constructor input")
	}
	completionView := source.completion.View()
	completionView.ResultJSON[0] = 'X'
	if source.completion.View().ResultJSON[0] == 'X' {
		t.Fatal("completion view aliases retained bytes")
	}
	result := compose(t, source)
	projection := result.Projection().View()
	projection.ResultJSON[0] = 'X'
	if result.Projection().View().ResultJSON[0] == 'X' {
		t.Fatal("projection aliases retained result bytes")
	}
	transcript := result.TranscriptBytes()
	summary := result.PublicSummaryBytes()
	transcript[0] = 'X'
	summary[0] = 'X'
	if result.TranscriptBytes()[0] == 'X' || result.PublicSummaryBytes()[0] == 'X' {
		t.Fatal("encoded result accessors alias retained bytes")
	}
}

func newRetainedTruthSource(t *testing.T, resultJSON []byte) *retainedTruthSource {
	t.Helper()
	plan := readFile(t, filepath.Join("..", "..", "..", "schemas", "conformance", "v0", "execution-plan", "ordinary.cbor"))
	planDigest := v0candidate.ExecutionPlanDigest(sha256.Sum256(plan))
	attempt := approvalattempt.ExecutionAttempt{
		AttemptID:            repeated16[approvalattempt.AttemptID](0xb2),
		ApprovalID:           repeated16[approvalattempt.ApprovalID](0xa1),
		AttemptNonce:         repeated16[approvalattempt.AttemptNonce](0xc3),
		RegistrationID:       repeated16[v0candidate.RegistrationID](0x33),
		RegistrationSequence: 1, PlanDigest: planDigest,
		InstallationID: repeated16[v0candidate.InstallationID](0x11), EpochSequence: 7,
		EpochDigest:  repeated32[v0candidate.TrustEpochDigest](0x22),
		SupervisorID: repeated16[v0candidate.SupervisorID](0x55),
		Purpose:      approvalattempt.ApprovalGrantPurpose, Audience: approvalattempt.ApprovalGrantAudience,
		ApprovalPayloadDigest: repeated32[approvalattempt.ApprovalPayloadDigest](0xd4),
		AuthorizationIdentity: repeated32[approvalattempt.ApprovalKeyAuthorizationIdentity](0x99),
		CreatedAt:             100, State: approvalattempt.AttemptCreated,
		StorageFormatVersion: lifecyclestate.SliceBAttemptStorageVersion,
	}
	created := registrationstate.CreatedAttempt{
		Attempt: attempt,
		Approval: registrationstate.ConsumedApprovalBinding{
			ApprovalID: attempt.ApprovalID, ConsumedAttemptID: attempt.AttemptID,
			AttemptNonce: attempt.AttemptNonce, RegistrationID: attempt.RegistrationID,
			RegistrationSequence: attempt.RegistrationSequence, PlanDigest: attempt.PlanDigest,
			InstallationID: attempt.InstallationID, EpochSequence: attempt.EpochSequence,
			EpochDigest: attempt.EpochDigest, SupervisorID: attempt.SupervisorID,
			Purpose: attempt.Purpose, Audience: attempt.Audience,
			PayloadDigest:         attempt.ApprovalPayloadDigest,
			AuthorizationIdentity: attempt.AuthorizationIdentity,
			State:                 approvalattempt.ApprovalConsumed,
			StorageFormatVersion:  lifecyclestate.SliceBApprovalStorageVersion,
		},
		Registration: registrationstate.StoredRegistration{
			WireRegistrationBytes: []byte{0xa0}, ExactPlanBytes: plan,
			RecomputedPlanDigest: planDigest, RegisteredAtUnixSeconds: 90,
			StorageFormatVersion: lifecyclestate.SliceBRegistrationStorageVersion,
			RetentionState:       registrationstate.RetentionStateRetained,
		},
		PlanRoleBindings: ordinaryPlanBindings(t),
	}
	backend := mustBackendBinding(t)
	bindings, err := lifecyclestate.NewImmutableBindings(lifecyclestate.ImmutableBindingsView{
		AttemptID: attempt.AttemptID, ApprovalID: attempt.ApprovalID, AttemptNonce: attempt.AttemptNonce,
		RegistrationID: attempt.RegistrationID, RegistrationSequence: attempt.RegistrationSequence,
		PlanDigest: attempt.PlanDigest, InstallationID: attempt.InstallationID,
		EpochSequence: attempt.EpochSequence, EpochDigest: attempt.EpochDigest,
		SupervisorID: attempt.SupervisorID, ApprovalPurpose: attempt.Purpose,
		ApprovalAudience: attempt.Audience, ApprovalPayloadDigest: attempt.ApprovalPayloadDigest,
		AuthorizationIdentity: attempt.AuthorizationIdentity, AttemptCreatedAt: attempt.CreatedAt,
		AttemptStorageVersion:           attempt.StorageFormatVersion,
		ApprovalStorageVersion:          created.Approval.StorageFormatVersion,
		RegistrationStorageVersion:      created.Registration.StorageFormatVersion,
		SourceManifestDigest:            created.PlanRoleBindings.SourceManifestDigest,
		InlineInputDigest:               created.PlanRoleBindings.InlineInputDigest,
		RuntimeBundleManifestDigest:     created.PlanRoleBindings.RuntimeBundleManifestDigest,
		ProfileReviewAttestationDigests: created.PlanRoleBindings.ProfileReviewAttestationDigests,
		ProfileRegistryEntryDigest:      created.PlanRoleBindings.ProfileRegistryEntryDigest,
		BackendValidationRecordDigest:   created.PlanRoleBindings.BackendValidationRecordDigest,
		BackendConfigurationDigest:      created.PlanRoleBindings.BackendConfigurationDigest,
		TrustSnapshotDigest:             created.PlanRoleBindings.TrustSnapshotDigest,
		PolicyDecisionDigest:            created.PlanRoleBindings.PolicyDecisionDigest,
		Backend:                         backend,
	})
	if err != nil {
		t.Fatalf("new immutable bindings: %v", err)
	}
	instance, err := lifecyclestate.NewBackendInstanceIdentity(
		lifecyclestate.BackendInstanceFake,
		[]byte("fake-instance-known-answer"),
	)
	if err != nil {
		t.Fatalf("new fake instance: %v", err)
	}
	effectDomain, err := lifecyclestate.NewDomainIdentifier(
		lifecyclestate.DomainEffectID,
		bytes.Repeat([]byte{0xe6}, 16),
	)
	if err != nil {
		t.Fatal(err)
	}
	effectID, err := lifecyclestate.NewEffectID(effectDomain)
	if err != nil {
		t.Fatal(err)
	}
	lifecycle, err := lifecyclestate.NewRecord(lifecyclestate.RecordView{
		FormatVersion: lifecyclestate.LifecycleRecordFormatVersion,
		RecordVersion: 14, SnapshotGeneration: 9,
		Bindings: bindings, ImmutableBindingDigest: bindings.Digest(),
		OperationSequence: 6, Operation: lifecyclestate.OperationDestroy,
		EffectID: effectID, EffectStatus: lifecyclestate.EffectConfirmed,
		Instance: instance, State: lifecyclestate.StateDestroyed, CleanupRequired: false,
		LastConfirmedCheckpoint: lifecyclestate.CheckpointDestroy,
		FirstFailure:            lifecyclestate.FailureNone, FailureOperation: lifecyclestate.OperationNone,
		LastReconciliation:     lifecyclestate.ReconciliationAuthoritativelyAbsent,
		AutomaticRecoveryCount: 1, RecoveryFence: lifecyclestate.RecoveryFenceNone,
		OpenedAt: 100, LastTransitionAt: 120,
		TerminalAt: lifecyclestate.OptionalUnixSeconds{Present: true, Value: 120},
	}, 120)
	if err != nil {
		t.Fatalf("new terminal lifecycle: %v", err)
	}
	completion := completionForJSON(t, created, bindings.Digest(), resultJSON)
	return &retainedTruthSource{created: created, lifecycle: lifecycle, completion: completion}
}

func ordinaryPlanBindings(t *testing.T) v0candidate.ExecutionPlanRoleBindings {
	t.Helper()
	return v0candidate.ExecutionPlanRoleBindings{
		InstallationID:              repeated16[v0candidate.InstallationID](0x11),
		EpochDigest:                 repeated32[v0candidate.TrustEpochDigest](0x22),
		SourceManifestDigest:        mustHex32[v0candidate.SourceManifestDigest](t, "c387c80094027ffbcacb573f44f5f6b4dec4d243bb436b24dd644434feaa1d14"),
		InlineInputDigest:           mustHex32[v0candidate.InlineInputDigest](t, "bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72"),
		RuntimeBundleManifestDigest: repeated32[v0candidate.RuntimeBundleManifestDigest](0x55),
		ProfileReviewAttestationDigests: []v0candidate.ProfileReviewAttestationDigest{
			repeated32[v0candidate.ProfileReviewAttestationDigest](0x66),
			repeated32[v0candidate.ProfileReviewAttestationDigest](0x67),
		},
		ProfileRegistryEntryDigest:    repeated32[v0candidate.ProfileRegistryEntryDigest](0x77),
		BackendValidationRecordDigest: repeated32[v0candidate.BackendValidationRecordDigest](0x88),
		BackendConfigurationDigest:    repeated32[v0candidate.BackendConfigurationDigest](0x99),
		TrustSnapshotDigest:           repeated32[v0candidate.TrustSnapshotDigest](0xaa),
		PolicyDecisionDigest:          repeated32[v0candidate.PolicyDecisionDigest](0xbb),
	}
}

func mustBackendBinding(t *testing.T) lifecyclestate.BackendBinding {
	t.Helper()
	result, err := lifecyclestate.NewBackendBinding(lifecyclestate.BackendBindingView{
		Kind:                          lifecyclestate.BackendFakeNoGuest,
		ProtocolVersion:               lifecyclestate.FakeBackendProtocolVersion,
		ImplementationIdentityDigest:  repeated32[lifecyclestate.BackendImplementationDigest](0xe1),
		BackendConfigurationDigest:    repeated32[v0candidate.BackendConfigurationDigest](0x99),
		BackendValidationRecordDigest: repeated32[v0candidate.BackendValidationRecordDigest](0x88),
		CreatesGuest:                  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func completionForJSON(
	t *testing.T,
	created registrationstate.CreatedAttempt,
	binding lifecyclestate.ImmutableBindingDigest,
	exact []byte,
) CompletionRecord {
	t.Helper()
	return mustCompletion(t, CompletionRecordView{
		RecordVersion: CompletionRecordVersion, CommitSequence: 1,
		AttemptID: created.Attempt.AttemptID, RegistrationID: created.Attempt.RegistrationID,
		PlanDigest: created.Attempt.PlanDigest, ImmutableBindingDigest: binding,
		Kind: CompletionFakeNoGuest, Status: CompletionSucceeded, Commit: CompletionCommittedLast,
		ResultJSON: exact, ResultJSONByteLength: v0candidate.UInt53(len(exact)),
		ResultJSONDigest: ResultJSONDigest(sha256.Sum256(exact)),
	})
}

func mustCompletion(t *testing.T, view CompletionRecordView) CompletionRecord {
	t.Helper()
	record, err := NewCompletionRecord(view)
	if err != nil {
		t.Fatalf("new completion record: %v", err)
	}
	return record
}

func observedLifecycle(t *testing.T, terminal lifecyclestate.Record) lifecyclestate.Record {
	t.Helper()
	view := terminal.View()
	view.RecordVersion--
	view.OperationSequence = 4
	view.Operation = lifecyclestate.OperationObserve
	view.EffectStatus = lifecyclestate.EffectConfirmed
	view.State = lifecyclestate.StateObserved
	view.CleanupRequired = true
	view.LastConfirmedCheckpoint = lifecyclestate.CheckpointObserve
	view.LastReconciliation = lifecyclestate.ReconciliationNone
	view.AutomaticRecoveryCount = 0
	view.TerminalAt = lifecyclestate.OptionalUnixSeconds{}
	view.LastTransitionAt = 110
	record, err := lifecyclestate.NewRecord(view, 120)
	if err != nil {
		t.Fatalf("new observed lifecycle: %v", err)
	}
	return record
}

func failedDestroyedLifecycle(t *testing.T, terminal lifecyclestate.Record) lifecyclestate.Record {
	t.Helper()
	view := terminal.View()
	view.FirstFailure = lifecyclestate.FailureLifecycle
	view.FailureOperation = lifecyclestate.OperationObserve
	record, err := lifecyclestate.NewRecord(view, 120)
	if err != nil {
		t.Fatalf("new failed terminal lifecycle: %v", err)
	}
	return record
}

func compose(t *testing.T, source *retainedTruthSource) Result {
	t.Helper()
	compositor, err := New(Sources{Attempts: source, Lifecycles: source, Completions: source})
	if err != nil {
		t.Fatal(err)
	}
	result, err := compositor.Compose(context.Background(), source.created.Attempt.AttemptID)
	if err != nil {
		t.Fatalf("compose: %v", err)
	}
	return result
}

func readKnownAnswer(t *testing.T) knownAnswerFixture {
	t.Helper()
	exact := readFile(t, filepath.Join("testdata", "known-answer.json"))
	var fixture knownAnswerFixture
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&fixture); err != nil {
		t.Fatal(err)
	}
	return fixture
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	exact, err := os.ReadFile(path) //nolint:gosec // test reads fixed repository fixtures only
	if err != nil {
		t.Fatal(err)
	}
	return exact
}

func classification(t *testing.T, err error) Classification {
	t.Helper()
	if err == nil {
		return ""
	}
	value, ok := ErrorClassification(err)
	if !ok {
		t.Fatalf("unclassified error: %v", err)
	}
	return value
}

func repeated16[T ~[16]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func repeated32[T ~[32]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func mustHex32[T ~[32]byte](t *testing.T, value string) T {
	t.Helper()
	exact, err := hex.DecodeString(value)
	if err != nil || len(exact) != 32 {
		t.Fatal(errors.New("invalid fixed digest"))
	}
	var result T
	copy(result[:], exact)
	return result
}
