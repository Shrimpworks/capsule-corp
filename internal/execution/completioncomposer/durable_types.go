package completioncomposer

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	DurableCompletionObjectVersion     = uint64(0)
	DurableCompletionStoreVersion      = uint64(0)
	MaximumDurableCompletionBytes      = 368_640
	MaximumRetainedCompletions         = 4_096
	MaximumRetainedResultJSONBytes     = 67_108_864
	MaximumDurableCompletionStoreBytes = 100_663_296
)

const (
	durableCompletionObjectType = "capsule.unwired-fake-durable-completion"
	durableCompletionStoreType  = "capsule.unwired-fake-completion-store"
)

type ResultIntegrityDisposition string
type RunnerIdentityDisposition string
type TeardownDisposition string
type DescendantAbsenceDisposition string
type TerminalClassification string
type DurableCompletionDigest [32]byte

const (
	ResultIntegrityTypedJSONSHA256 ResultIntegrityDisposition   = "supervisor-validated-typed-json-sha256"
	RunnerIdentityUnresolvedFake   RunnerIdentityDisposition    = "unresolved-fake-backend-has-no-runner"
	TeardownFakeDestroyConfirmed   TeardownDisposition          = "fake-destroy-confirmed"
	DescendantsFakeAbsent          DescendantAbsenceDisposition = "fake-authoritatively-absent-no-process-tree"
	TerminalCompletedFakeLocal     TerminalClassification       = "completed-fake-no-guest-local-mechanic"
)

// CommitRequest is copied before validation. It deliberately contains no EOF,
// exit status, diagnostics, paths, artifact names, timing, backend flags, or
// caller-selected lifecycle facts.
type CommitRequest struct {
	ObjectVersion   uint64
	AttemptID       approvalattempt.AttemptID
	ResultJSON      []byte
	ResultIntegrity ResultIntegrityDisposition
}

func (request CommitRequest) clone() CommitRequest {
	request.ResultJSON = bytes.Clone(request.ResultJSON)
	return request
}

type DurableCompletionView struct {
	ObjectVersion       uint64
	CommitSequence      v0candidate.PositiveUInt53
	Commit              CompletionCommit
	ResultIntegrity     ResultIntegrityDisposition
	RunnerIdentity      RunnerIdentityDisposition
	Teardown            TeardownDisposition
	DescendantAbsence   DescendantAbsenceDisposition
	Terminal            TerminalClassification
	Projection          TerminalProjectionView
	Completion          CompletionRecord
	Transcript          []byte
	TranscriptDigest    TranscriptDigest
	PublicSummary       []byte
	DurableRecordDigest DurableCompletionDigest
}

// DurableCompletion is the immutable committed-last object returned by the
// fixed local store. It owns every byte slice.
type DurableCompletion struct{ view DurableCompletionView }

func newDurableCompletion(
	sequence v0candidate.PositiveUInt53,
	completion CompletionRecord,
	result Result,
) (DurableCompletion, error) {
	view := DurableCompletionView{
		ObjectVersion:     DurableCompletionObjectVersion,
		CommitSequence:    sequence,
		Commit:            CompletionCommittedLast,
		ResultIntegrity:   ResultIntegrityTypedJSONSHA256,
		RunnerIdentity:    RunnerIdentityUnresolvedFake,
		Teardown:          TeardownFakeDestroyConfirmed,
		DescendantAbsence: DescendantsFakeAbsent,
		Terminal:          TerminalCompletedFakeLocal,
		Projection:        result.Projection().View(),
		Completion:        completion,
		Transcript:        result.TranscriptBytes(),
		TranscriptDigest:  result.TranscriptDigest(),
		PublicSummary:     result.PublicSummaryBytes(),
	}
	view.DurableRecordDigest = digestDurableCompletion(view)
	record := DurableCompletion{view: cloneDurableView(view)}
	if err := record.Validate(); err != nil {
		return DurableCompletion{}, err
	}
	return record, nil
}

func restoreDurableCompletion(view DurableCompletionView) (DurableCompletion, error) {
	record := DurableCompletion{view: cloneDurableView(view)}
	if err := record.Validate(); err != nil {
		return DurableCompletion{}, err
	}
	return record, nil
}

func (record DurableCompletion) View() DurableCompletionView {
	return cloneDurableView(record.view)
}

func (record DurableCompletion) CompletionRecord() CompletionRecord { return record.view.Completion }
func (record DurableCompletion) TranscriptBytes() []byte            { return bytes.Clone(record.view.Transcript) }
func (record DurableCompletion) TranscriptSHA256() TranscriptDigest {
	return record.view.TranscriptDigest
}
func (record DurableCompletion) PublicSummaryBytes() []byte {
	return bytes.Clone(record.view.PublicSummary)
}
func (record DurableCompletion) Digest() DurableCompletionDigest {
	return record.view.DurableRecordDigest
}

func (record DurableCompletion) Validate() error {
	view := cloneDurableView(record.view)
	if view.ObjectVersion != DurableCompletionObjectVersion {
		return classified(ClassificationUnsupported, "durable-completion-version")
	}
	if view.CommitSequence == 0 || uint64(view.CommitSequence) > v0candidate.MaxSafeInteger {
		return classified(ClassificationSchema, "durable-completion-sequence")
	}
	if view.Commit != CompletionCommittedLast ||
		view.ResultIntegrity != ResultIntegrityTypedJSONSHA256 ||
		view.RunnerIdentity != RunnerIdentityUnresolvedFake ||
		view.Teardown != TeardownFakeDestroyConfirmed ||
		view.DescendantAbsence != DescendantsFakeAbsent ||
		view.Terminal != TerminalCompletedFakeLocal {
		return classified(ClassificationUnsupported, "durable-completion-disposition")
	}
	if err := view.Completion.Validate(); err != nil {
		return classified(ClassificationRecoveryRequired, "durable-completion-record")
	}
	completion := view.Completion.View()
	projection := view.Projection
	if completion.CommitSequence != view.CommitSequence ||
		completion.AttemptID != projection.AttemptID ||
		completion.RegistrationID != projection.RegistrationID ||
		completion.PlanDigest != projection.PlanDigest ||
		completion.ImmutableBindingDigest != projection.ImmutableBindingDigest ||
		completion.ResultJSONDigest != projection.ResultJSONDigest ||
		view.Completion.Digest() != projection.CompletionRecordDigest ||
		!bytes.Equal(completion.ResultJSON, projection.ResultJSON) {
		return classified(ClassificationBinding, "durable-completion-projection")
	}
	if projection.LifecycleState != lifecyclestate.StateDestroyed ||
		projection.LifecycleOperationSequence != 6 || projection.CleanupRequired ||
		projection.LifecycleLastReconciliation != lifecyclestate.ReconciliationAuthoritativelyAbsent {
		return classified(ClassificationBinding, "durable-completion-lifecycle")
	}
	expectedTranscript, expectedDigest, err := encodeTranscript(projection)
	if err != nil || expectedDigest != view.TranscriptDigest || !bytes.Equal(expectedTranscript, view.Transcript) {
		return classified(ClassificationBinding, "durable-completion-transcript")
	}
	expectedSummary, err := encodePublicSummary(projection.AttemptID, expectedDigest)
	if err != nil || !bytes.Equal(expectedSummary, view.PublicSummary) {
		return classified(ClassificationBinding, "durable-completion-summary")
	}
	if len(view.Transcript) > MaximumTranscriptBytes || len(view.PublicSummary) > MaximumPublicSummaryBytes ||
		len(completion.ResultJSON) > MaximumResultJSONBytes {
		return classified(ClassificationMalformed, "durable-completion-cap")
	}
	if digestDurableCompletion(view) != view.DurableRecordDigest {
		return classified(ClassificationBinding, "durable-completion-digest")
	}
	return nil
}

func cloneDurableView(view DurableCompletionView) DurableCompletionView {
	view.Projection.ResultJSON = bytes.Clone(view.Projection.ResultJSON)
	completionView := view.Completion.View()
	if restored, err := RestoreCompletionRecord(completionView, view.Completion.Digest()); err == nil {
		view.Completion = restored
	}
	view.Transcript = bytes.Clone(view.Transcript)
	view.PublicSummary = bytes.Clone(view.PublicSummary)
	return view
}

func digestDurableCompletion(view DurableCompletionView) DurableCompletionDigest {
	hash := sha256.New()
	writeDigestField(hash, []byte("capsule.unwired.fake-durable-completion/v0"))
	writeDigestUint64(hash, view.ObjectVersion)
	writeDigestUint64(hash, uint64(view.CommitSequence))
	writeDigestField(hash, []byte(view.Commit))
	writeDigestField(hash, []byte(view.ResultIntegrity))
	writeDigestField(hash, []byte(view.RunnerIdentity))
	writeDigestField(hash, []byte(view.Teardown))
	writeDigestField(hash, []byte(view.DescendantAbsence))
	writeDigestField(hash, []byte(view.Terminal))
	writeProjectionDigestFields(hash, view.Projection)
	completionDigest := view.Completion.Digest()
	writeDigestField(hash, completionDigest[:])
	writeDigestField(hash, view.Transcript)
	writeDigestField(hash, view.TranscriptDigest[:])
	writeDigestField(hash, view.PublicSummary)
	var result DurableCompletionDigest
	copy(result[:], hash.Sum(nil))
	return result
}

func writeProjectionDigestFields(writer digestWriter, value TerminalProjectionView) {
	for _, field := range [][]byte{
		value.AttemptID[:], value.ApprovalID[:], value.RegistrationID[:], value.PlanDigest[:],
		value.InstallationID[:], value.EpochDigest[:], value.SupervisorID[:],
		value.ApprovalPayloadDigest[:], value.AuthorizationIdentity[:],
		value.SourceManifestDigest[:],
		value.RuntimeBundleManifestDigest[:], value.ProfileRegistryEntryDigest[:],
		value.BackendValidationRecordDigest[:], value.BackendConfigurationDigest[:],
		value.BackendImplementationDigest[:], value.BackendInstanceDigest[:],
		value.ImmutableBindingDigest[:], value.CompletionRecordDigest[:], value.ResultJSON,
		value.ResultJSONDigest[:], []byte(value.LifecycleState),
		[]byte(value.LifecycleLastReconciliation),
	} {
		writeDigestField(writer, field)
	}
	var number [8]byte
	binary.BigEndian.PutUint64(number[:], uint64(value.LifecycleOperationSequence))
	writeDigestField(writer, number[:])
	if value.CleanupRequired {
		writeDigestField(writer, []byte{1})
	} else {
		writeDigestField(writer, []byte{0})
	}
}
