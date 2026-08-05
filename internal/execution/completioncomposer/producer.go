package completioncomposer

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

var (
	ErrDurableCompletionNotFound     = errors.New("durable completion not found")
	ErrCompletionCommitIndeterminate = errors.New("completion commit outcome indeterminate")
)

type DurableCompletionStore interface {
	CommitCompletion(
		context.Context,
		approvalattempt.AttemptID,
		func(v0candidate.PositiveUInt53) (DurableCompletion, error),
	) (DurableCompletion, bool, error)
	ReadDurableCompletion(context.Context, approvalattempt.AttemptID) (DurableCompletion, error)
}

type ProducerSources struct {
	Attempts   CreatedAttemptReader
	Lifecycles LifecycleReader
	Store      DurableCompletionStore
}

// Producer owns the unwired commit-last transaction. It accepts result bytes,
// but it derives every authority, lifecycle, runner, teardown, absence, and
// terminal-classification field from retained readers and fixed FakeBackend
// semantics.
type Producer struct{ sources ProducerSources }

func NewProducer(sources ProducerSources) (*Producer, error) {
	if sources.Attempts == nil || sources.Lifecycles == nil || sources.Store == nil {
		return nil, errors.New("attempt, lifecycle, and durable-completion sources are required")
	}
	return &Producer{sources: sources}, nil
}

func (producer *Producer) Complete(ctx context.Context, request CommitRequest) (DurableCompletion, error) {
	request = request.clone()
	if err := validateCommitRequest(request); err != nil {
		return DurableCompletion{}, err
	}
	record, created, err := producer.sources.Store.CommitCompletion(
		ctx,
		request.AttemptID,
		func(sequence v0candidate.PositiveUInt53) (DurableCompletion, error) {
			return producer.build(ctx, request, sequence)
		},
	)
	if err != nil {
		return DurableCompletion{}, mapCompletionStoreFailure(err)
	}
	if !created {
		view := record.CompletionRecord().View()
		digest := ResultJSONDigest(sha256.Sum256(request.ResultJSON))
		if view.AttemptID != request.AttemptID || view.ResultJSONDigest != digest ||
			!bytes.Equal(view.ResultJSON, request.ResultJSON) {
			return DurableCompletion{}, classified(ClassificationReplay, "completion-stale-replay")
		}
	}
	return record, nil
}

func (producer *Producer) build(
	ctx context.Context,
	request CommitRequest,
	sequence v0candidate.PositiveUInt53,
) (DurableCompletion, error) {
	created, err := producer.sources.Attempts.ResolveCreated(ctx, request.AttemptID)
	if err != nil {
		return DurableCompletion{}, mapSourceFailure(err, "completion-attempt")
	}
	lifecycle, err := producer.sources.Lifecycles.ReadLifecycle(ctx, request.AttemptID)
	if err != nil {
		return DurableCompletion{}, mapSourceFailure(err, "completion-lifecycle")
	}
	lifecycleView, err := validateTerminalLifecycle(lifecycle, cloneCreatedAttempt(created))
	if err != nil {
		return DurableCompletion{}, err
	}
	resultDigest := ResultJSONDigest(sha256.Sum256(request.ResultJSON))
	completion, err := NewCompletionRecord(CompletionRecordView{
		RecordVersion:          CompletionRecordVersion,
		CommitSequence:         sequence,
		AttemptID:              request.AttemptID,
		RegistrationID:         created.Attempt.RegistrationID,
		PlanDigest:             created.Attempt.PlanDigest,
		ImmutableBindingDigest: lifecycleView.ImmutableBindingDigest,
		Kind:                   CompletionFakeNoGuest,
		Status:                 CompletionSucceeded,
		Commit:                 CompletionCommittedLast,
		ResultJSON:             request.ResultJSON,
		ResultJSONByteLength:   v0candidate.UInt53(len(request.ResultJSON)),
		ResultJSONDigest:       resultDigest,
	})
	if err != nil {
		return DurableCompletion{}, err
	}
	completionSource := fixedCompletionReader{record: completion}
	compositor, err := New(Sources{
		Attempts:    producer.sources.Attempts,
		Lifecycles:  producer.sources.Lifecycles,
		Completions: completionSource,
	})
	if err != nil {
		return DurableCompletion{}, classified(ClassificationLocalFailure, "completion-compositor")
	}
	result, err := compositor.Compose(ctx, request.AttemptID)
	if err != nil {
		return DurableCompletion{}, err
	}
	return newDurableCompletion(sequence, completion, result)
}

func validateCommitRequest(request CommitRequest) error {
	if request.ObjectVersion != DurableCompletionObjectVersion {
		return classified(ClassificationUnsupported, "completion-request-version")
	}
	if request.AttemptID == (approvalattempt.AttemptID{}) {
		return classified(ClassificationSchema, "completion-request-attempt")
	}
	if request.ResultIntegrity != ResultIntegrityTypedJSONSHA256 {
		return classified(ClassificationUnsupported, "completion-result-integrity")
	}
	if len(request.ResultJSON) > MaximumResultJSONBytes {
		return classified(ClassificationMalformed, "completion-result-cap")
	}
	return validateTypedJSON(request.ResultJSON)
}

type fixedCompletionReader struct{ record CompletionRecord }

func (reader fixedCompletionReader) ReadCommittedCompletion(
	context.Context,
	approvalattempt.AttemptID,
) (CompletionRecord, error) {
	return reader.record, nil
}

func mapCompletionStoreFailure(err error) error {
	if errors.Is(err, ErrCompletionCommitIndeterminate) {
		return classified(ClassificationRecoveryRequired, "completion-commit-indeterminate")
	}
	if _, ok := ErrorClassification(err); ok {
		return err
	}
	return classified(ClassificationLocalFailure, "completion-store")
}
