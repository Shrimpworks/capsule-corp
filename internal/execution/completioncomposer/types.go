package completioncomposer

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	CompletionRecordVersion   = uint64(0)
	MaximumResultJSONBytes    = 262_144
	MaximumTranscriptBytes    = 4_096
	MaximumPublicSummaryBytes = 256
)

type Classification string

const (
	ClassificationMalformed        Classification = "MALFORMED"
	ClassificationUnsupported      Classification = "UNSUPPORTED"
	ClassificationSchema           Classification = "SCHEMA"
	ClassificationBinding          Classification = "BINDING"
	ClassificationNotTerminal      Classification = "NOT_TERMINAL"
	ClassificationRecoveryRequired Classification = "RECOVERY_REQUIRED"
	ClassificationCapacity         Classification = "CAPACITY"
	ClassificationReplay           Classification = "REPLAY"
	ClassificationLocalFailure     Classification = "LOCAL_FAILURE"
)

type compositorError struct {
	classification Classification
	code           string
}

func (e *compositorError) Error() string {
	return fmt.Sprintf("%s: %s", e.classification, e.code)
}

func classified(classification Classification, code string) error {
	return &compositorError{classification: classification, code: code}
}

func ErrorClassification(err error) (Classification, bool) {
	var target *compositorError
	if errors.As(err, &target) {
		return target.classification, true
	}
	return "", false
}

var (
	// ErrTerminalFactMissing tells the compositor that the Supervisor store has
	// no committed fact yet. Absence never becomes success.
	ErrTerminalFactMissing = errors.New("terminal fact missing")
	// ErrTerminalTruthCorrupt tells the compositor that a retained source could
	// not produce a coherent validated fact. It requires recovery or repair.
	ErrTerminalTruthCorrupt = errors.New("terminal truth corrupt")
)

type CompletionKind string
type CompletionStatus string
type CompletionCommit string

const (
	CompletionFakeNoGuest   CompletionKind   = "fake-no-guest-json"
	CompletionSucceeded     CompletionStatus = "succeeded"
	CompletionCommittedLast CompletionCommit = "committed-last"
)

type ResultJSONDigest [32]byte
type CompletionRecordDigest [32]byte
type TranscriptDigest [32]byte

// CompletionRecordView is the fixed input returned only by the future
// Supervisor-owned committed-completion reader. It has no exit status, EOF,
// prose, diagnostics, paths, artifact names, or timing fields.
type CompletionRecordView struct {
	RecordVersion          uint64
	CommitSequence         v0candidate.PositiveUInt53
	AttemptID              approvalattempt.AttemptID
	RegistrationID         v0candidate.RegistrationID
	PlanDigest             v0candidate.ExecutionPlanDigest
	ImmutableBindingDigest lifecyclestate.ImmutableBindingDigest
	Kind                   CompletionKind
	Status                 CompletionStatus
	Commit                 CompletionCommit
	ResultJSON             []byte
	ResultJSONByteLength   v0candidate.UInt53
	ResultJSONDigest       ResultJSONDigest
}

// CompletionRecord is immutable and owns exact result bytes.
type CompletionRecord struct {
	view   CompletionRecordView
	digest CompletionRecordDigest
}

func NewCompletionRecord(view CompletionRecordView) (CompletionRecord, error) {
	view.ResultJSON = bytes.Clone(view.ResultJSON)
	if err := validateCompletionRecordView(view); err != nil {
		return CompletionRecord{}, err
	}
	record := CompletionRecord{view: view}
	record.digest = digestCompletionRecord(view)
	return record, nil
}

func RestoreCompletionRecord(
	view CompletionRecordView,
	digest CompletionRecordDigest,
) (CompletionRecord, error) {
	record, err := NewCompletionRecord(view)
	if err != nil {
		return CompletionRecord{}, err
	}
	if record.digest != digest {
		return CompletionRecord{}, classified(ClassificationBinding, "completion-record-digest")
	}
	return record, nil
}

func (record CompletionRecord) View() CompletionRecordView {
	view := record.view
	view.ResultJSON = bytes.Clone(record.view.ResultJSON)
	return view
}

func (record CompletionRecord) Digest() CompletionRecordDigest { return record.digest }

func (record CompletionRecord) Validate() error {
	if err := validateCompletionRecordView(record.view); err != nil {
		return err
	}
	if digestCompletionRecord(record.view) != record.digest {
		return classified(ClassificationBinding, "completion-record-digest")
	}
	return nil
}

func validateCompletionRecordView(view CompletionRecordView) error {
	if view.RecordVersion != CompletionRecordVersion {
		return classified(ClassificationUnsupported, "completion-record-version")
	}
	if view.CommitSequence == 0 || uint64(view.CommitSequence) > v0candidate.MaxSafeInteger {
		return classified(ClassificationSchema, "completion-commit-sequence")
	}
	if view.AttemptID == (approvalattempt.AttemptID{}) ||
		view.RegistrationID == (v0candidate.RegistrationID{}) ||
		view.PlanDigest == (v0candidate.ExecutionPlanDigest{}) ||
		view.ImmutableBindingDigest == (lifecyclestate.ImmutableBindingDigest{}) {
		return classified(ClassificationSchema, "completion-identities")
	}
	if view.Kind != CompletionFakeNoGuest || view.Status != CompletionSucceeded ||
		view.Commit != CompletionCommittedLast {
		return classified(ClassificationUnsupported, "completion-disposition")
	}
	if len(view.ResultJSON) > MaximumResultJSONBytes {
		return classified(ClassificationMalformed, "completion-result-cap")
	}
	if uint64(view.ResultJSONByteLength) != uint64(len(view.ResultJSON)) ||
		uint64(view.ResultJSONByteLength) > v0candidate.MaxSafeInteger {
		return classified(ClassificationBinding, "completion-result-length")
	}
	computed := ResultJSONDigest(sha256.Sum256(view.ResultJSON))
	if computed != view.ResultJSONDigest {
		return classified(ClassificationBinding, "completion-result-digest")
	}
	if err := validateTypedJSON(view.ResultJSON); err != nil {
		return err
	}
	return nil
}

func digestCompletionRecord(view CompletionRecordView) CompletionRecordDigest {
	hash := sha256.New()
	writeDigestField(hash, []byte("capsule.unwired.fake-completion-record/v0"))
	writeDigestUint64(hash, view.RecordVersion)
	writeDigestUint64(hash, uint64(view.CommitSequence))
	writeDigestField(hash, view.AttemptID[:])
	writeDigestField(hash, view.RegistrationID[:])
	writeDigestField(hash, view.PlanDigest[:])
	writeDigestField(hash, view.ImmutableBindingDigest[:])
	writeDigestField(hash, []byte(view.Kind))
	writeDigestField(hash, []byte(view.Status))
	writeDigestField(hash, []byte(view.Commit))
	writeDigestField(hash, view.ResultJSON)
	writeDigestUint64(hash, uint64(view.ResultJSONByteLength))
	writeDigestField(hash, view.ResultJSONDigest[:])
	var result CompletionRecordDigest
	copy(result[:], hash.Sum(nil))
	return result
}

type digestWriter interface{ Write([]byte) (int, error) }

func writeDigestField(writer digestWriter, value []byte) {
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(value)))
	_, _ = writer.Write(length[:])
	_, _ = writer.Write(value)
}

func writeDigestUint64(writer digestWriter, value uint64) {
	var encoded [8]byte
	binary.BigEndian.PutUint64(encoded[:], value)
	_, _ = writer.Write(encoded[:])
}

type TerminalProjectionView struct {
	AttemptID                     approvalattempt.AttemptID
	ApprovalID                    approvalattempt.ApprovalID
	RegistrationID                v0candidate.RegistrationID
	PlanDigest                    v0candidate.ExecutionPlanDigest
	InstallationID                v0candidate.InstallationID
	EpochDigest                   v0candidate.TrustEpochDigest
	SupervisorID                  v0candidate.SupervisorID
	ApprovalPayloadDigest         approvalattempt.ApprovalPayloadDigest
	AuthorizationIdentity         approvalattempt.ApprovalKeyAuthorizationIdentity
	SourceManifestDigest          v0candidate.SourceManifestDigest
	RuntimeBundleManifestDigest   v0candidate.RuntimeBundleManifestDigest
	ProfileRegistryEntryDigest    v0candidate.ProfileRegistryEntryDigest
	BackendValidationRecordDigest v0candidate.BackendValidationRecordDigest
	BackendConfigurationDigest    v0candidate.BackendConfigurationDigest
	BackendImplementationDigest   lifecyclestate.BackendImplementationDigest
	BackendInstanceDigest         lifecyclestate.BackendInstanceDigest
	ImmutableBindingDigest        lifecyclestate.ImmutableBindingDigest
	CompletionRecordDigest        CompletionRecordDigest
	ResultJSON                    []byte
	ResultJSONDigest              ResultJSONDigest
	LifecycleState                lifecyclestate.LifecycleState
	LifecycleOperationSequence    v0candidate.PositiveUInt53
	LifecycleLastReconciliation   lifecyclestate.ReconciliationStatus
	CleanupRequired               bool
}

type TerminalProjection struct{ view TerminalProjectionView }

func (projection TerminalProjection) View() TerminalProjectionView {
	view := projection.view
	view.ResultJSON = bytes.Clone(projection.view.ResultJSON)
	return view
}

type Result struct {
	projection       TerminalProjection
	transcript       []byte
	transcriptDigest TranscriptDigest
	publicSummary    []byte
}

func (result Result) Projection() TerminalProjection     { return result.projection }
func (result Result) TranscriptBytes() []byte            { return bytes.Clone(result.transcript) }
func (result Result) TranscriptDigest() TranscriptDigest { return result.transcriptDigest }
func (result Result) PublicSummaryBytes() []byte         { return bytes.Clone(result.publicSummary) }
