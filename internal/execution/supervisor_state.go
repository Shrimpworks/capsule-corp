package execution

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
)

// This package is an OS-neutral executable state-machine contract. It creates no
// guest and does not select the final macOS Supervisor implementation language.

var (
	ErrApprovalBinding    = errors.New("approval binding mismatch")
	ErrApprovalConsumed   = errors.New("approval already consumed")
	ErrApprovalExpired    = errors.New("approval expired or not yet valid")
	ErrCleanupRequired    = errors.New("attempt still has a cleanup obligation")
	ErrCleanupUnresolved  = errors.New("backend cleanup outcome is unresolved")
	ErrDuplicateApproval  = errors.New("approval payload already recorded")
	ErrInvalidState       = errors.New("invalid supervisor state transition")
	ErrNotFound           = errors.New("supervisor record not found")
	ErrRealBackendBlocked = errors.New("state-machine scaffold refuses guest-creating backends")
	ErrTransitionFenced   = errors.New("attempts are fenced by trust transition state")
)

const (
	ApprovalPurpose  = "capsule.plan.approve"
	ApprovalAudience = "capsule.execution-supervisor"
)

type OpaqueID [16]byte

func (id OpaqueID) IsZero() bool {
	return id == OpaqueID{}
}

type Digest [32]byte

func (digest Digest) IsZero() bool {
	return digest == Digest{}
}

func DigestBytes(value []byte) Digest {
	return sha256.Sum256(value)
}

type TrustPhase string

const (
	TrustStable                      TrustPhase = "stable"
	TrustTransitionFenced            TrustPhase = "transition-fenced"
	TrustAwaitingComponentAcceptance TrustPhase = "awaiting-component-acceptance"
	TrustRepairRequired              TrustPhase = "repair-required"
)

type ApprovalState string

const (
	ApprovalUnused      ApprovalState = "unused"
	ApprovalConsumed    ApprovalState = "consumed"
	ApprovalInvalidated ApprovalState = "invalidated"
)

type AttemptState string

const (
	AttemptCreated             AttemptState = "created"
	AttemptPreparing           AttemptState = "preparing"
	AttemptBackendCreateIntent AttemptState = "backend-create-intent"
	AttemptBackendCreated      AttemptState = "backend-created"
	AttemptStaging             AttemptState = "staging"
	AttemptStarting            AttemptState = "starting"
	AttemptRunning             AttemptState = "running"
	AttemptCollecting          AttemptState = "collecting"
	AttemptDestroying          AttemptState = "destroying"
	AttemptUnresolved          AttemptState = "unresolved"
	AttemptTerminal            AttemptState = "terminal"
)

type TerminalClass string

const (
	TerminalNone              TerminalClass = ""
	TerminalSimulationPassed  TerminalClass = "simulation-passed"
	TerminalSimulationFailed  TerminalClass = "simulation-failed"
	TerminalBackendFailed     TerminalClass = "backend-failed"
	TerminalCleanupUnresolved TerminalClass = "cleanup-unresolved"
)

type ComponentRole string

const (
	ComponentDaemon     ComponentRole = "daemon"
	ComponentBroker     ComponentRole = "broker"
	ComponentSupervisor ComponentRole = "supervisor"
	ComponentUpdater    ComponentRole = "updater"
)

var requiredComponentRoles = [...]ComponentRole{
	ComponentDaemon,
	ComponentBroker,
	ComponentSupervisor,
	ComponentUpdater,
}

type InstallationState struct {
	InstallationID  OpaqueID
	SupervisorID    OpaqueID
	EpochDigest     Digest
	TrustPhase      TrustPhase
	AttemptsEnabled bool
	Transition      *TransitionRecord
	Registrations   map[OpaqueID]PlanRegistration
	Approvals       map[OpaqueID]ApprovalRecord
	Attempts        map[OpaqueID]AttemptRecord
}

type TransitionRecord struct {
	ID                   OpaqueID
	FromEpoch            Digest
	TargetEpoch          Digest
	AcceptedIncarnations map[ComponentRole]OpaqueID
}

type PlanRegistration struct {
	ID             OpaqueID
	ExactPlanBytes []byte
	PlanDigest     Digest
	InstallationID OpaqueID
	EpochDigest    Digest
	SupervisorID   OpaqueID
}

type ApprovalRecord struct {
	ID               OpaqueID
	RegistrationID   OpaqueID
	CanonicalPayload []byte
	SignedEnvelope   []byte
	PayloadDigest    Digest
	EnvelopeDigest   Digest
	AttemptNonce     OpaqueID
	IssuedAt         time.Time
	ExpiresAt        time.Time
	State            ApprovalState
	AttemptID        OpaqueID
	InvalidatedBy    string
}

type AttemptRecord struct {
	ID                OpaqueID
	RegistrationID    OpaqueID
	ApprovalID        OpaqueID
	EpochDigest       Digest
	State             AttemptState
	CleanupRequired   bool
	ReconciliationKey OpaqueID
	BackendHandle     string
	TerminalClass     TerminalClass
	Failure           string
}

type VerifiedApproval struct {
	CanonicalPayload []byte
	InstallationID   OpaqueID
	EpochDigest      Digest
	RegistrationID   OpaqueID
	PlanDigest       Digest
	SupervisorID     OpaqueID
	AttemptNonce     OpaqueID
	Purpose          string
	Audience         string
	IssuedAt         time.Time
	ExpiresAt        time.Time
}

type StateStore interface {
	Snapshot(context.Context) (InstallationState, error)
	Update(context.Context, func(*InstallationState) error) error
}

type PlanValidator interface {
	ValidatePlan(context.Context, []byte) error
}

type ApprovalVerifier interface {
	VerifyApproval(context.Context, []byte) (VerifiedApproval, error)
}

type IntegrityAssessor interface {
	AssessRegistration(context.Context, PlanRegistration) error
}

type IdentifierSource interface {
	NewIdentifier() (OpaqueID, error)
}

type Clock interface {
	Now() time.Time
}

type Checkpoint string

const (
	CheckpointBackendCreateIntent Checkpoint = "backend-create-intent-durable"
	CheckpointBackendHandle       Checkpoint = "backend-handle-durable"
)

type CheckpointHook func(Checkpoint) error

type SupervisorOptions struct {
	Store            StateStore
	PlanValidator    PlanValidator
	ApprovalVerifier ApprovalVerifier
	Integrity        IntegrityAssessor
	Identifiers      IdentifierSource
	Clock            Clock
	Checkpoint       CheckpointHook
	MaxPlanBytes     int
	MaxApprovalBytes int
}

type SupervisorCore struct {
	store            StateStore
	planValidator    PlanValidator
	approvalVerifier ApprovalVerifier
	integrity        IntegrityAssessor
	identifiers      IdentifierSource
	clock            Clock
	checkpoint       CheckpointHook
	maxPlanBytes     int
	maxApprovalBytes int
}

func NewSupervisorCore(options SupervisorOptions) (*SupervisorCore, error) {
	if options.Store == nil || options.PlanValidator == nil || options.ApprovalVerifier == nil ||
		options.Integrity == nil || options.Identifiers == nil || options.Clock == nil {
		return nil, errors.New("all supervisor authority ports are required")
	}
	if options.MaxPlanBytes <= 0 {
		options.MaxPlanBytes = 1 << 20
	}
	if options.MaxApprovalBytes <= 0 {
		options.MaxApprovalBytes = 16 << 10
	}
	if options.Checkpoint == nil {
		options.Checkpoint = func(Checkpoint) error { return nil }
	}
	return &SupervisorCore{
		store:            options.Store,
		planValidator:    options.PlanValidator,
		approvalVerifier: options.ApprovalVerifier,
		integrity:        options.Integrity,
		identifiers:      options.Identifiers,
		clock:            options.Clock,
		checkpoint:       options.Checkpoint,
		maxPlanBytes:     options.MaxPlanBytes,
		maxApprovalBytes: options.MaxApprovalBytes,
	}, nil
}

func (core *SupervisorCore) RegisterPlan(
	ctx context.Context,
	exactCanonicalPlanBytes []byte,
) (PlanRegistration, error) {
	if len(exactCanonicalPlanBytes) == 0 || len(exactCanonicalPlanBytes) > core.maxPlanBytes {
		return PlanRegistration{}, errors.New("canonical plan bytes are empty or exceed the bound")
	}
	planCopy := append([]byte(nil), exactCanonicalPlanBytes...)
	if err := core.planValidator.ValidatePlan(ctx, append([]byte(nil), planCopy...)); err != nil {
		return PlanRegistration{}, fmt.Errorf("independent plan validation: %w", err)
	}
	registrationID, err := core.nextIdentifier()
	if err != nil {
		return PlanRegistration{}, err
	}

	var registration PlanRegistration
	err = core.store.Update(ctx, func(state *InstallationState) error {
		if err := requireExecutionEnabled(state); err != nil {
			return err
		}
		if _, exists := state.Registrations[registrationID]; exists {
			return errors.New("identifier source returned a duplicate registration ID")
		}
		registration = PlanRegistration{
			ID:             registrationID,
			ExactPlanBytes: append([]byte(nil), planCopy...),
			PlanDigest:     DigestBytes(planCopy),
			InstallationID: state.InstallationID,
			EpochDigest:    state.EpochDigest,
			SupervisorID:   state.SupervisorID,
		}
		state.Registrations[registrationID] = cloneRegistration(registration)
		return nil
	})
	return cloneRegistration(registration), err
}

func (core *SupervisorCore) GetRegisteredPlan(
	ctx context.Context,
	registrationID OpaqueID,
) ([]byte, error) {
	state, err := core.store.Snapshot(ctx)
	if err != nil {
		return nil, err
	}
	registration, exists := state.Registrations[registrationID]
	if !exists {
		return nil, ErrNotFound
	}
	return append([]byte(nil), registration.ExactPlanBytes...), nil
}

func (core *SupervisorCore) SubmitApproval(
	ctx context.Context,
	registrationID OpaqueID,
	signedEnvelope []byte,
) (ApprovalRecord, error) {
	if registrationID.IsZero() || len(signedEnvelope) == 0 ||
		len(signedEnvelope) > core.maxApprovalBytes {
		return ApprovalRecord{}, errors.New("approval reference or envelope is invalid")
	}
	verified, err := core.approvalVerifier.VerifyApproval(
		ctx,
		append([]byte(nil), signedEnvelope...),
	)
	if err != nil {
		return ApprovalRecord{}, fmt.Errorf("approval verification: %w", err)
	}
	if len(verified.CanonicalPayload) == 0 ||
		len(verified.CanonicalPayload) > core.maxApprovalBytes {
		return ApprovalRecord{}, errors.New("verified approval payload is empty or exceeds the bound")
	}
	approvalID, err := core.nextIdentifier()
	if err != nil {
		return ApprovalRecord{}, err
	}
	now := core.clock.Now()

	var approval ApprovalRecord
	err = core.store.Update(ctx, func(state *InstallationState) error {
		if err := requireExecutionEnabled(state); err != nil {
			return err
		}
		registration, exists := state.Registrations[registrationID]
		if !exists {
			return ErrNotFound
		}
		if err := validateApprovalBindings(state, registration, verified, now); err != nil {
			return err
		}
		payloadDigest := DigestBytes(verified.CanonicalPayload)
		for _, recorded := range state.Approvals {
			if recorded.PayloadDigest == payloadDigest {
				return ErrDuplicateApproval
			}
		}
		if _, exists := state.Approvals[approvalID]; exists {
			return errors.New("identifier source returned a duplicate approval ID")
		}
		approval = ApprovalRecord{
			ID:               approvalID,
			RegistrationID:   registrationID,
			CanonicalPayload: append([]byte(nil), verified.CanonicalPayload...),
			SignedEnvelope:   append([]byte(nil), signedEnvelope...),
			PayloadDigest:    payloadDigest,
			EnvelopeDigest:   DigestBytes(signedEnvelope),
			AttemptNonce:     verified.AttemptNonce,
			IssuedAt:         verified.IssuedAt,
			ExpiresAt:        verified.ExpiresAt,
			State:            ApprovalUnused,
		}
		state.Approvals[approvalID] = cloneApproval(approval)
		return nil
	})
	return cloneApproval(approval), err
}

func (core *SupervisorCore) RequestAttempt(
	ctx context.Context,
	registrationID OpaqueID,
	approvalID OpaqueID,
) (AttemptRecord, error) {
	state, err := core.store.Snapshot(ctx)
	if err != nil {
		return AttemptRecord{}, err
	}
	registration, exists := state.Registrations[registrationID]
	if !exists {
		return AttemptRecord{}, ErrNotFound
	}
	if err := core.integrity.AssessRegistration(ctx, cloneRegistration(registration)); err != nil {
		return AttemptRecord{}, fmt.Errorf("runtime-integrity preflight: %w", err)
	}
	attemptID, err := core.nextIdentifier()
	if err != nil {
		return AttemptRecord{}, err
	}
	now := core.clock.Now()

	var attempt AttemptRecord
	err = core.store.Update(ctx, func(current *InstallationState) error {
		if err := requireExecutionEnabled(current); err != nil {
			return err
		}
		freshRegistration, exists := current.Registrations[registrationID]
		if !exists {
			return ErrNotFound
		}
		if freshRegistration.EpochDigest != current.EpochDigest {
			return ErrApprovalBinding
		}
		approval, exists := current.Approvals[approvalID]
		if !exists {
			return ErrNotFound
		}
		if approval.RegistrationID != registrationID {
			return ErrApprovalBinding
		}
		if approval.State != ApprovalUnused {
			return ErrApprovalConsumed
		}
		if now.Before(approval.IssuedAt) || !now.Before(approval.ExpiresAt) {
			return ErrApprovalExpired
		}
		if _, exists := current.Attempts[attemptID]; exists {
			return errors.New("identifier source returned a duplicate attempt ID")
		}
		attempt = AttemptRecord{
			ID:             attemptID,
			RegistrationID: registrationID,
			ApprovalID:     approvalID,
			EpochDigest:    current.EpochDigest,
			State:          AttemptCreated,
		}
		approval.State = ApprovalConsumed
		approval.AttemptID = attemptID
		current.Approvals[approvalID] = approval
		current.Attempts[attemptID] = attempt
		return nil
	})
	return attempt, err
}

func (core *SupervisorCore) FenceTrustTransition(
	ctx context.Context,
	transitionID OpaqueID,
	targetEpoch Digest,
) error {
	if transitionID.IsZero() || targetEpoch.IsZero() {
		return errors.New("transition ID and target epoch are required")
	}
	return core.store.Update(ctx, func(state *InstallationState) error {
		if err := requireExecutionEnabled(state); err != nil {
			return err
		}
		if targetEpoch == state.EpochDigest {
			return errors.New("target epoch must differ from the active epoch")
		}
		for _, attempt := range state.Attempts {
			if attempt.State != AttemptTerminal || attempt.CleanupRequired {
				return ErrCleanupRequired
			}
		}
		state.AttemptsEnabled = false
		state.TrustPhase = TrustTransitionFenced
		state.Transition = &TransitionRecord{
			ID:                   transitionID,
			FromEpoch:            state.EpochDigest,
			TargetEpoch:          targetEpoch,
			AcceptedIncarnations: make(map[ComponentRole]OpaqueID),
		}
		for id, approval := range state.Approvals {
			if approval.State == ApprovalUnused {
				approval.State = ApprovalInvalidated
				approval.InvalidatedBy = "trust-transition-fence"
				state.Approvals[id] = approval
			}
		}
		return nil
	})
}

func (core *SupervisorCore) CommitEpochForComponentAcceptance(
	ctx context.Context,
	transitionID OpaqueID,
) error {
	return core.store.Update(ctx, func(state *InstallationState) error {
		if state.TrustPhase != TrustTransitionFenced || state.AttemptsEnabled ||
			state.Transition == nil || state.Transition.ID != transitionID {
			return ErrInvalidState
		}
		state.EpochDigest = state.Transition.TargetEpoch
		state.TrustPhase = TrustAwaitingComponentAcceptance
		return nil
	})
}

func (core *SupervisorCore) AcceptCurrentComponent(
	ctx context.Context,
	transitionID OpaqueID,
	role ComponentRole,
	processIncarnation OpaqueID,
) error {
	if !isRequiredRole(role) || processIncarnation.IsZero() {
		return errors.New("required component role and process incarnation are required")
	}
	return core.store.Update(ctx, func(state *InstallationState) error {
		if state.TrustPhase != TrustAwaitingComponentAcceptance || state.AttemptsEnabled ||
			state.Transition == nil || state.Transition.ID != transitionID ||
			state.EpochDigest != state.Transition.TargetEpoch {
			return ErrInvalidState
		}
		state.Transition.AcceptedIncarnations[role] = processIncarnation
		for _, required := range requiredComponentRoles {
			if state.Transition.AcceptedIncarnations[required].IsZero() {
				return nil
			}
		}
		state.TrustPhase = TrustStable
		state.AttemptsEnabled = true
		return nil
	})
}

func (core *SupervisorCore) RecoverTrustState(ctx context.Context) error {
	return core.store.Update(ctx, func(state *InstallationState) error {
		if state.TrustPhase == TrustStable && state.AttemptsEnabled {
			return nil
		}
		state.TrustPhase = TrustRepairRequired
		state.AttemptsEnabled = false
		return nil
	})
}

func (core *SupervisorCore) nextIdentifier() (OpaqueID, error) {
	id, err := core.identifiers.NewIdentifier()
	if err != nil {
		return OpaqueID{}, fmt.Errorf("generate opaque identifier: %w", err)
	}
	if id.IsZero() {
		return OpaqueID{}, errors.New("identifier source returned the zero identifier")
	}
	return id, nil
}

func validateApprovalBindings(
	state *InstallationState,
	registration PlanRegistration,
	approval VerifiedApproval,
	now time.Time,
) error {
	if registration.InstallationID != state.InstallationID ||
		registration.EpochDigest != state.EpochDigest ||
		registration.SupervisorID != state.SupervisorID ||
		approval.InstallationID != state.InstallationID ||
		approval.EpochDigest != state.EpochDigest ||
		approval.RegistrationID != registration.ID ||
		approval.PlanDigest != registration.PlanDigest ||
		approval.SupervisorID != state.SupervisorID ||
		approval.AttemptNonce.IsZero() ||
		approval.Purpose != ApprovalPurpose ||
		approval.Audience != ApprovalAudience {
		return ErrApprovalBinding
	}
	if approval.IssuedAt.IsZero() || approval.ExpiresAt.IsZero() ||
		!approval.ExpiresAt.After(approval.IssuedAt) || now.Before(approval.IssuedAt) ||
		!now.Before(approval.ExpiresAt) {
		return ErrApprovalExpired
	}
	return nil
}

func requireExecutionEnabled(state *InstallationState) error {
	if !state.AttemptsEnabled || state.TrustPhase != TrustStable {
		return ErrTransitionFenced
	}
	return nil
}

func isRequiredRole(role ComponentRole) bool {
	for _, required := range requiredComponentRoles {
		if role == required {
			return true
		}
	}
	return false
}

func cloneRegistration(registration PlanRegistration) PlanRegistration {
	registration.ExactPlanBytes = append([]byte(nil), registration.ExactPlanBytes...)
	return registration
}

func cloneApproval(approval ApprovalRecord) ApprovalRecord {
	approval.CanonicalPayload = append([]byte(nil), approval.CanonicalPayload...)
	approval.SignedEnvelope = append([]byte(nil), approval.SignedEnvelope...)
	return approval
}

type CryptoIdentifierSource struct{}

func (CryptoIdentifierSource) NewIdentifier() (OpaqueID, error) {
	var id OpaqueID
	_, err := rand.Read(id[:])
	return id, err
}

type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}

type PlanValidatorFunc func(context.Context, []byte) error

func (function PlanValidatorFunc) ValidatePlan(ctx context.Context, plan []byte) error {
	return function(ctx, plan)
}

type ApprovalVerifierFunc func(context.Context, []byte) (VerifiedApproval, error)

func (function ApprovalVerifierFunc) VerifyApproval(
	ctx context.Context,
	envelope []byte,
) (VerifiedApproval, error) {
	return function(ctx, envelope)
}

type IntegrityAssessorFunc func(context.Context, PlanRegistration) error

func (function IntegrityAssessorFunc) AssessRegistration(
	ctx context.Context,
	registration PlanRegistration,
) error {
	return function(ctx, registration)
}
