package completioncomposer

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

type CreatedAttemptReader interface {
	ResolveCreated(context.Context, approvalattempt.AttemptID) (registrationstate.CreatedAttempt, error)
}

type LifecycleReader interface {
	ReadLifecycle(context.Context, approvalattempt.AttemptID) (lifecyclestate.Record, error)
}

type CommittedCompletionReader interface {
	ReadCommittedCompletion(context.Context, approvalattempt.AttemptID) (CompletionRecord, error)
}

type Sources struct {
	Attempts    CreatedAttemptReader
	Lifecycles  LifecycleReader
	Completions CommittedCompletionReader
}

type Compositor struct{ sources Sources }

func New(sources Sources) (*Compositor, error) {
	if sources.Attempts == nil || sources.Lifecycles == nil || sources.Completions == nil {
		return nil, errors.New("attempt, lifecycle, and committed-completion readers are required")
	}
	return &Compositor{sources: sources}, nil
}

func (compositor *Compositor) Compose(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, classified(ClassificationLocalFailure, "terminal-context")
	}
	if attemptID == (approvalattempt.AttemptID{}) {
		return Result{}, classified(ClassificationSchema, "terminal-attempt-id")
	}
	created, err := compositor.sources.Attempts.ResolveCreated(ctx, attemptID)
	if err != nil {
		return Result{}, mapSourceFailure(err, "terminal-attempt")
	}
	created = cloneCreatedAttempt(created)
	plan, err := validateCreatedAttempt(created, attemptID)
	if err != nil {
		return Result{}, err
	}
	completion, err := compositor.sources.Completions.ReadCommittedCompletion(ctx, attemptID)
	if err != nil {
		return Result{}, mapSourceFailure(err, "terminal-completion")
	}
	if err := completion.Validate(); err != nil {
		return Result{}, classified(ClassificationRecoveryRequired, "terminal-completion-corrupt")
	}
	lifecycle, err := compositor.sources.Lifecycles.ReadLifecycle(ctx, attemptID)
	if err != nil {
		return Result{}, mapSourceFailure(err, "terminal-lifecycle")
	}
	lifecycleView, err := validateTerminalLifecycle(lifecycle, created)
	if err != nil {
		return Result{}, err
	}
	completionView := completion.View()
	if completionView.AttemptID != attemptID ||
		completionView.RegistrationID != created.Attempt.RegistrationID ||
		completionView.PlanDigest != created.Attempt.PlanDigest ||
		completionView.ImmutableBindingDigest != lifecycle.ImmutableBindingDigest() {
		return Result{}, classified(ClassificationBinding, "terminal-completion-cross-link")
	}
	if uint64(completionView.ResultJSONByteLength) > uint64(plan.View().OutputMaxJSONBytes) {
		return Result{}, classified(ClassificationBinding, "terminal-result-plan-cap")
	}
	projection := buildProjection(created, lifecycleView, completion)
	transcript, transcriptDigest, err := encodeTranscript(projection)
	if err != nil {
		return Result{}, err
	}
	publicSummary, err := encodePublicSummary(attemptID, transcriptDigest)
	if err != nil {
		return Result{}, err
	}
	return Result{
		projection: TerminalProjection{view: projection}, transcript: transcript,
		transcriptDigest: transcriptDigest, publicSummary: publicSummary,
	}, nil
}

func mapSourceFailure(err error, code string) error {
	if errors.Is(err, ErrTerminalFactMissing) || errors.Is(err, registrationstate.ErrLifecycleNotFound) {
		return classified(ClassificationNotTerminal, code+"-missing")
	}
	if errors.Is(err, ErrTerminalTruthCorrupt) ||
		errors.Is(err, registrationstate.ErrStoreRepairRequired) ||
		errors.Is(err, registrationstate.ErrCommitOutcomeIndeterminate) ||
		errors.Is(err, registrationstate.ErrStoreClosed) ||
		errors.Is(err, registrationstate.ErrStoreOwnerFenced) ||
		errors.Is(err, registrationstate.ErrStoreOwnerSessionMismatch) {
		return classified(ClassificationRecoveryRequired, code+"-unavailable")
	}
	return classified(ClassificationLocalFailure, code+"-read")
}

func cloneCreatedAttempt(created registrationstate.CreatedAttempt) registrationstate.CreatedAttempt {
	created.Registration.WireRegistrationBytes = bytes.Clone(created.Registration.WireRegistrationBytes)
	created.Registration.ExactPlanBytes = bytes.Clone(created.Registration.ExactPlanBytes)
	created.PlanRoleBindings.ProfileReviewAttestationDigests = append(
		[]v0candidate.ProfileReviewAttestationDigest(nil),
		created.PlanRoleBindings.ProfileReviewAttestationDigests...,
	)
	return created
}

func validateCreatedAttempt(
	created registrationstate.CreatedAttempt,
	attemptID approvalattempt.AttemptID,
) (*v0candidate.DecodedExecutionPlan, error) {
	attempt := created.Attempt
	approval := created.Approval
	if attempt.AttemptID != attemptID || attempt.State != approvalattempt.AttemptCreated ||
		approval.State != approvalattempt.ApprovalConsumed || approval.ConsumedAttemptID != attemptID ||
		attempt.ApprovalID != approval.ApprovalID || attempt.AttemptNonce != approval.AttemptNonce ||
		attempt.RegistrationID != approval.RegistrationID ||
		attempt.RegistrationSequence != approval.RegistrationSequence ||
		attempt.PlanDigest != approval.PlanDigest || attempt.InstallationID != approval.InstallationID ||
		attempt.EpochSequence != approval.EpochSequence || attempt.EpochDigest != approval.EpochDigest ||
		attempt.SupervisorID != approval.SupervisorID || attempt.Purpose != approval.Purpose ||
		attempt.Audience != approval.Audience ||
		attempt.ApprovalPayloadDigest != approval.PayloadDigest ||
		attempt.AuthorizationIdentity != approval.AuthorizationIdentity ||
		attempt.StorageFormatVersion != lifecyclestate.SliceBAttemptStorageVersion ||
		approval.StorageFormatVersion != lifecyclestate.SliceBApprovalStorageVersion {
		return nil, classified(ClassificationBinding, "terminal-attempt-approval-cross-link")
	}
	if len(created.Registration.ExactPlanBytes) == 0 ||
		created.Registration.StorageFormatVersion != lifecyclestate.SliceBRegistrationStorageVersion ||
		created.Registration.RetentionState != registrationstate.RetentionStateRetained {
		return nil, classified(ClassificationBinding, "terminal-registration-record")
	}
	plan, err := v0candidate.DecodeExecutionPlan(
		created.Registration.ExactPlanBytes,
		cloneRoleBindings(created.PlanRoleBindings),
	)
	if err != nil {
		return nil, classified(ClassificationRecoveryRequired, "terminal-plan-decode")
	}
	computed := v0candidate.ExecutionPlanDigest(sha256.Sum256(created.Registration.ExactPlanBytes))
	if computed != created.Registration.RecomputedPlanDigest || plan.Digest() != computed ||
		attempt.PlanDigest != computed {
		return nil, classified(ClassificationRecoveryRequired, "terminal-plan-digest")
	}
	return plan, nil
}

func cloneRoleBindings(value v0candidate.ExecutionPlanRoleBindings) v0candidate.ExecutionPlanRoleBindings {
	value.ProfileReviewAttestationDigests = append(
		[]v0candidate.ProfileReviewAttestationDigest(nil),
		value.ProfileReviewAttestationDigests...,
	)
	return value
}

func validateTerminalLifecycle(
	record lifecyclestate.Record,
	created registrationstate.CreatedAttempt,
) (lifecyclestate.RecordView, error) {
	view := record.View()
	if !view.TerminalAt.Present {
		return lifecyclestate.RecordView{}, classified(ClassificationNotTerminal, "terminal-lifecycle-incomplete")
	}
	if err := record.Validate(view.TerminalAt.Value); err != nil {
		return lifecyclestate.RecordView{}, classified(ClassificationRecoveryRequired, "terminal-lifecycle-corrupt")
	}
	if view.State != lifecyclestate.StateDestroyed || view.CleanupRequired ||
		view.LastReconciliation != lifecyclestate.ReconciliationAuthoritativelyAbsent {
		return lifecyclestate.RecordView{}, classified(ClassificationNotTerminal, "terminal-lifecycle-incomplete")
	}
	if view.FirstFailure != lifecyclestate.FailureNone ||
		view.Operation != lifecyclestate.OperationDestroy ||
		view.OperationSequence != 6 || view.EffectStatus != lifecyclestate.EffectConfirmed ||
		view.LastConfirmedCheckpoint != lifecyclestate.CheckpointDestroy || !view.Instance.Present() {
		return lifecyclestate.RecordView{}, classified(ClassificationBinding, "terminal-lifecycle-contradiction")
	}
	bindings := view.Bindings.View()
	attempt := created.Attempt
	approval := created.Approval
	roles := created.PlanRoleBindings
	backend := bindings.Backend.View()
	if bindings.AttemptID != attempt.AttemptID || bindings.ApprovalID != attempt.ApprovalID ||
		bindings.AttemptNonce != attempt.AttemptNonce ||
		bindings.RegistrationID != attempt.RegistrationID ||
		bindings.RegistrationSequence != attempt.RegistrationSequence ||
		bindings.PlanDigest != attempt.PlanDigest || bindings.InstallationID != attempt.InstallationID ||
		bindings.EpochSequence != attempt.EpochSequence || bindings.EpochDigest != attempt.EpochDigest ||
		bindings.SupervisorID != attempt.SupervisorID || bindings.ApprovalPurpose != attempt.Purpose ||
		bindings.ApprovalAudience != attempt.Audience ||
		bindings.ApprovalPayloadDigest != attempt.ApprovalPayloadDigest ||
		bindings.AuthorizationIdentity != attempt.AuthorizationIdentity ||
		bindings.AttemptCreatedAt != attempt.CreatedAt ||
		bindings.AttemptStorageVersion != attempt.StorageFormatVersion ||
		bindings.ApprovalStorageVersion != approval.StorageFormatVersion ||
		bindings.RegistrationStorageVersion != created.Registration.StorageFormatVersion ||
		bindings.SourceManifestDigest != roles.SourceManifestDigest ||
		bindings.InlineInputDigest != roles.InlineInputDigest ||
		bindings.RuntimeBundleManifestDigest != roles.RuntimeBundleManifestDigest ||
		!equalReviewDigests(bindings.ProfileReviewAttestationDigests, roles.ProfileReviewAttestationDigests) ||
		bindings.ProfileRegistryEntryDigest != roles.ProfileRegistryEntryDigest ||
		bindings.BackendValidationRecordDigest != roles.BackendValidationRecordDigest ||
		bindings.BackendConfigurationDigest != roles.BackendConfigurationDigest ||
		bindings.TrustSnapshotDigest != roles.TrustSnapshotDigest ||
		bindings.PolicyDecisionDigest != roles.PolicyDecisionDigest ||
		backend.Kind != lifecyclestate.BackendFakeNoGuest || backend.CreatesGuest {
		return lifecyclestate.RecordView{}, classified(ClassificationBinding, "terminal-lifecycle-cross-link")
	}
	return view, nil
}

func equalReviewDigests(
	left, right []v0candidate.ProfileReviewAttestationDigest,
) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func buildProjection(
	created registrationstate.CreatedAttempt,
	lifecycle lifecyclestate.RecordView,
	completion CompletionRecord,
) TerminalProjectionView {
	bindings := lifecycle.Bindings.View()
	completionView := completion.View()
	return TerminalProjectionView{
		AttemptID: created.Attempt.AttemptID, ApprovalID: created.Attempt.ApprovalID,
		RegistrationID: created.Attempt.RegistrationID, PlanDigest: created.Attempt.PlanDigest,
		InstallationID: created.Attempt.InstallationID, EpochDigest: created.Attempt.EpochDigest,
		SupervisorID:                  created.Attempt.SupervisorID,
		ApprovalPayloadDigest:         created.Attempt.ApprovalPayloadDigest,
		AuthorizationIdentity:         created.Attempt.AuthorizationIdentity,
		SourceManifestDigest:          bindings.SourceManifestDigest,
		RuntimeBundleManifestDigest:   bindings.RuntimeBundleManifestDigest,
		ProfileRegistryEntryDigest:    bindings.ProfileRegistryEntryDigest,
		BackendValidationRecordDigest: bindings.BackendValidationRecordDigest,
		BackendConfigurationDigest:    bindings.BackendConfigurationDigest,
		BackendImplementationDigest:   bindings.Backend.View().ImplementationIdentityDigest,
		BackendInstanceDigest:         lifecycle.Instance.Digest(),
		ImmutableBindingDigest:        lifecycle.ImmutableBindingDigest,
		CompletionRecordDigest:        completion.Digest(), ResultJSON: completionView.ResultJSON,
		ResultJSONDigest: completionView.ResultJSONDigest,
		LifecycleState:   lifecycle.State, LifecycleOperationSequence: lifecycle.OperationSequence,
		LifecycleLastReconciliation: lifecycle.LastReconciliation,
		CleanupRequired:             lifecycle.CleanupRequired,
	}
}

type transcriptJSON struct {
	ObjectType                    string `json:"objectType"`
	ObjectVersion                 uint64 `json:"objectVersion"`
	Scope                         string `json:"scope"`
	AttemptID                     string `json:"attemptId"`
	ApprovalID                    string `json:"approvalId"`
	RegistrationID                string `json:"registrationId"`
	PlanDigest                    string `json:"planDigest"`
	InstallationID                string `json:"installationId"`
	EpochDigest                   string `json:"epochDigest"`
	SupervisorID                  string `json:"supervisorId"`
	ApprovalPayloadDigest         string `json:"approvalPayloadDigest"`
	AuthorizationIdentity         string `json:"authorizationIdentity"`
	SourceManifestDigest          string `json:"sourceManifestDigest"`
	RuntimeBundleManifestDigest   string `json:"runtimeBundleManifestDigest"`
	ProfileRegistryEntryDigest    string `json:"profileRegistryEntryDigest"`
	BackendValidationRecordDigest string `json:"backendValidationRecordDigest"`
	BackendConfigurationDigest    string `json:"backendConfigurationDigest"`
	BackendImplementationDigest   string `json:"backendImplementationDigest"`
	BackendInstanceDigest         string `json:"backendInstanceDigest"`
	ImmutableBindingDigest        string `json:"immutableBindingDigest"`
	CompletionRecordDigest        string `json:"completionRecordDigest"`
	ResultJSONBytes               uint64 `json:"resultJsonBytes"`
	ResultJSONDigest              string `json:"resultJsonDigest"`
	ResultIntegrity               string `json:"resultIntegrity"`
	Completion                    string `json:"completion"`
	TerminalClassification        string `json:"terminalClassification"`
	Lifecycle                     string `json:"lifecycle"`
	LifecycleOperationSequence    uint64 `json:"lifecycleOperationSequence"`
	Teardown                      string `json:"teardown"`
	Absence                       string `json:"absence"`
	RunnerIdentity                string `json:"runnerIdentity"`
	CleanupRequired               bool   `json:"cleanupRequired"`
	Limitation                    string `json:"limitation"`
}

func encodeTranscript(projection TerminalProjectionView) ([]byte, TranscriptDigest, error) {
	value := transcriptJSON{
		ObjectType: "capsule.unwired-fake-terminal-transcript", ObjectVersion: 0,
		Scope:                         "fake-no-guest-local-mechanic",
		AttemptID:                     hex.EncodeToString(projection.AttemptID[:]),
		ApprovalID:                    hex.EncodeToString(projection.ApprovalID[:]),
		RegistrationID:                hex.EncodeToString(projection.RegistrationID[:]),
		PlanDigest:                    hex.EncodeToString(projection.PlanDigest[:]),
		InstallationID:                hex.EncodeToString(projection.InstallationID[:]),
		EpochDigest:                   hex.EncodeToString(projection.EpochDigest[:]),
		SupervisorID:                  hex.EncodeToString(projection.SupervisorID[:]),
		ApprovalPayloadDigest:         hex.EncodeToString(projection.ApprovalPayloadDigest[:]),
		AuthorizationIdentity:         hex.EncodeToString(projection.AuthorizationIdentity[:]),
		SourceManifestDigest:          hex.EncodeToString(projection.SourceManifestDigest[:]),
		RuntimeBundleManifestDigest:   hex.EncodeToString(projection.RuntimeBundleManifestDigest[:]),
		ProfileRegistryEntryDigest:    hex.EncodeToString(projection.ProfileRegistryEntryDigest[:]),
		BackendValidationRecordDigest: hex.EncodeToString(projection.BackendValidationRecordDigest[:]),
		BackendConfigurationDigest:    hex.EncodeToString(projection.BackendConfigurationDigest[:]),
		BackendImplementationDigest:   hex.EncodeToString(projection.BackendImplementationDigest[:]),
		BackendInstanceDigest:         hex.EncodeToString(projection.BackendInstanceDigest[:]),
		ImmutableBindingDigest:        hex.EncodeToString(projection.ImmutableBindingDigest[:]),
		CompletionRecordDigest:        hex.EncodeToString(projection.CompletionRecordDigest[:]),
		ResultJSONBytes:               uint64(len(projection.ResultJSON)),
		ResultJSONDigest:              hex.EncodeToString(projection.ResultJSONDigest[:]),
		ResultIntegrity:               string(ResultIntegrityTypedJSONSHA256),
		Completion:                    string(CompletionCommittedLast),
		TerminalClassification:        string(TerminalCompletedFakeLocal),
		Lifecycle:                     string(projection.LifecycleState),
		LifecycleOperationSequence:    uint64(projection.LifecycleOperationSequence),
		Teardown:                      "destroy-confirmed", Absence: string(projection.LifecycleLastReconciliation),
		RunnerIdentity:  string(RunnerIdentityUnresolvedFake),
		CleanupRequired: projection.CleanupRequired,
		Limitation:      "no-guest-local-mechanic-not-signed-attestation-or-product-success",
	}
	exact, err := json.Marshal(value)
	if err != nil || len(exact) > MaximumTranscriptBytes {
		return nil, TranscriptDigest{}, classified(ClassificationLocalFailure, "terminal-transcript-encode")
	}
	digest := TranscriptDigest(sha256.Sum256(exact))
	return exact, digest, nil
}

type publicSummaryJSON struct {
	State        string `json:"state"`
	AttemptID    string `json:"attemptId"`
	TranscriptID string `json:"transcriptId"`
}

func encodePublicSummary(
	attemptID approvalattempt.AttemptID,
	transcriptDigest TranscriptDigest,
) ([]byte, error) {
	value := publicSummaryJSON{
		State:        "completed",
		AttemptID:    "attempt_" + hex.EncodeToString(attemptID[:]),
		TranscriptID: "transcript_sha256_" + hex.EncodeToString(transcriptDigest[:]),
	}
	exact, err := json.Marshal(value)
	if err != nil || len(exact) > MaximumPublicSummaryBytes {
		return nil, classified(ClassificationLocalFailure, "terminal-summary-encode")
	}
	return exact, nil
}
