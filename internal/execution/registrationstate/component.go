package registrationstate

import (
	"bytes"
	"context"
	"errors"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

// Options supplies only Supervisor-owned registration dependencies. Store,
// trusted clock, and identifier source are mandatory; Checkpoint is an
// optional test hook and must not perform product work or grant authority.
type Options struct {
	Store       StateStore
	Clock       TrustedClock
	Identifiers IdentifierSource
	Checkpoint  CheckpointHook
}

// Component is the unwired Supervisor registration slice. It neither wraps
// nor activates any other Supervisor implementation.
type Component struct {
	store       StateStore
	clock       TrustedClock
	identifiers IdentifierSource
	checkpoint  CheckpointHook
}

// New constructs the unwired registration component after rejecting missing
// authority dependencies. Construction opens no IPC endpoint, persists no
// state, admits no runtime/backend, and creates no guest.
func New(options Options) (*Component, error) {
	if options.Store == nil || options.Clock == nil || options.Identifiers == nil {
		return nil, errors.New("registration store, trusted clock, and identifier source are required")
	}
	if options.Checkpoint == nil {
		options.Checkpoint = func(context.Context, Checkpoint, *v0candidate.DecodedExecutionPlan) {}
	}
	return &Component{
		store:       options.Store,
		clock:       options.Clock,
		identifiers: options.Identifiers,
		checkpoint:  options.Checkpoint,
	}, nil
}

// RegisterPlan accepts only exact candidate plan bytes, complete trusted role
// bindings, and authenticated daemon call context. The Task 4A wrapper owns all
// byte/predecoder/schema/role validation before any trusted-time persistence.
func (component *Component) RegisterPlan(
	ctx context.Context,
	call AuthenticatedCallContext,
	exactPlanBytes []byte,
	bindings v0candidate.ExecutionPlanRoleBindings,
) (PlanRegistration, error) {
	if err := authenticateRegistrationCall(call); err != nil {
		return PlanRegistration{}, err
	}
	decoded, err := v0candidate.DecodeExecutionPlan(exactPlanBytes, bindings)
	if err != nil {
		return PlanRegistration{}, err
	}
	component.checkpoint(ctx, CheckpointPlanDecoded, decoded)
	if err := ctx.Err(); err != nil {
		return PlanRegistration{}, classified(ClassificationLocalFailure, "registration-context-cancelled")
	}

	observation, err := component.clock.ObserveUnixSeconds(ctx)
	if err != nil || observation > v0candidate.MaxSafeInteger {
		return PlanRegistration{}, classified(ClassificationLocalFailure, "trusted-clock-observation-failed")
	}
	if err := component.store.persistTimeHighWater(ctx, v0candidate.UInt53(observation)); err != nil {
		if errors.Is(err, ErrCommitOutcomeIndeterminate) {
			return PlanRegistration{}, err
		}
		return PlanRegistration{}, classified(ClassificationLocalFailure, "time-high-water-write-failed")
	}
	component.checkpoint(ctx, CheckpointTimeHighWaterDone, decoded)

	planBytes := decoded.AuthoritativeBytes()
	planDigest := decoded.Digest()
	planView := decoded.View()
	var issued PlanRegistration
	var committedRejection error
	err = component.store.update(ctx, func(state *installationState) error {
		effectiveNow := uint64(state.TimeHighWaterUnixSeconds)
		if err := validatePlanStateBinding(state, planView); err != nil {
			return err
		}
		if err := validateTrustAdmission(state); err != nil {
			return err
		}
		expiresAt := uint64(planView.ExpiresAt)
		if expiresAt <= effectiveNow {
			return classified(ClassificationStale, "registration-expired")
		}
		if expiresAt-effectiveNow > RegistrationLifetimeSeconds {
			return classified(ClassificationPolicy, "registration-lifetime-exceeded")
		}
		if len(state.Registrations) >= MaxRetainedRegistrations {
			return classified(ClassificationCapacity, "retained-registration-capacity")
		}
		if countUnexpired(state.Registrations, state.TimeHighWaterUnixSeconds) >= MaxUnexpiredRegistrations {
			return classified(ClassificationCapacity, "unexpired-registration-capacity")
		}
		if uint64(state.LastRegistrationSequence) == v0candidate.MaxSafeInteger {
			state.TrustPhase = TrustRepairRequired
			state.TrustReason = sequenceExhaustedReason
			committedRejection = classified(ClassificationStale, "registration-sequence-exhausted")
			return nil
		}

		registrationID, err := component.identifiers.NewRegistrationID(ctx)
		if err != nil || registrationID == (v0candidate.RegistrationID{}) {
			return classified(ClassificationLocalFailure, "registration-identifier-failed")
		}
		if findRegistration(state.Registrations, registrationID) >= 0 {
			return classified(ClassificationLocalFailure, "duplicate-registration-identifier")
		}
		sequence := v0candidate.PositiveUInt53(uint64(state.LastRegistrationSequence) + 1)
		wireView := v0candidate.PlanRegistration{
			ObjectType:           v0candidate.PlanRegistrationObjectType,
			ObjectVersion:        v0candidate.CandidateObjectVersion,
			RegistrationID:       registrationID,
			RegistrationSequence: sequence,
			PlanDigest:           planDigest,
			InstallationID:       state.InstallationID,
			EpochSequence:        state.EpochSequence,
			EpochDigest:          state.EpochDigest,
			SupervisorID:         state.SupervisorID,
			ExpiresAt:            planView.ExpiresAt,
		}
		wireBytes := encodePlanRegistration(wireView)
		if _, err := v0candidate.DecodePlanRegistration(
			wireBytes,
			v0candidate.PlanRegistrationRoleBindings{
				RegistrationID: registrationID,
				PlanDigest:     planDigest,
				InstallationID: state.InstallationID,
				EpochDigest:    state.EpochDigest,
				SupervisorID:   state.SupervisorID,
			},
		); err != nil {
			return classified(ClassificationLocalFailure, "registration-wire-self-check-failed")
		}
		record := StoredRegistration{
			WireRegistrationBytes:   bytes.Clone(wireBytes),
			ExactPlanBytes:          bytes.Clone(planBytes),
			RecomputedPlanDigest:    planDigest,
			RegisteredAtUnixSeconds: state.TimeHighWaterUnixSeconds,
			StorageFormatVersion:    StorageFormatVersion,
			RetentionState:          RetentionStateRetained,
		}
		state.Registrations = append(state.Registrations, registrationEntry{
			Index: registrationIndex{
				RegistrationID:       registrationID,
				RegistrationSequence: sequence,
				InstallationID:       state.InstallationID,
				EpochSequence:        state.EpochSequence,
				EpochDigest:          state.EpochDigest,
				SupervisorID:         state.SupervisorID,
				ExpiresAt:            planView.ExpiresAt,
			},
			PlanBindings: clonePlanBindings(bindings),
			Record:       record,
		})
		state.LastRegistrationSequence = v0candidate.UInt53(sequence)
		registrationSetDigest, err := appendRegistrationSetDigest(state.RegistrationSetDigest, record)
		if err != nil {
			return classified(ClassificationLocalFailure, "registration-digest-encode-failed")
		}
		state.RegistrationSetDigest = registrationSetDigest
		issued = PlanRegistration{bytes: bytes.Clone(wireBytes), view: wireView}
		return nil
	})
	if err != nil {
		if _, classifiedError := ErrorClassification(err); classifiedError {
			return PlanRegistration{}, err
		}
		if errors.Is(err, ErrCommitOutcomeIndeterminate) {
			return PlanRegistration{}, err
		}
		return PlanRegistration{}, classified(ClassificationLocalFailure, "registration-commit-failed")
	}
	if committedRejection != nil {
		return PlanRegistration{}, committedRejection
	}
	return PlanRegistration{bytes: issued.Bytes(), view: issued.view}, nil
}

// ResolveUsable performs an ID-only final lookup for the later fake-backend
// slice. It persists trusted time first and never accepts replacement plan or
// registration bytes.
func (component *Component) ResolveUsable(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) (StoredRegistration, error) {
	observation, err := component.clock.ObserveUnixSeconds(ctx)
	if err != nil || observation > v0candidate.MaxSafeInteger {
		return StoredRegistration{}, classified(ClassificationLocalFailure, "trusted-clock-observation-failed")
	}
	if err := component.store.persistTimeHighWater(ctx, v0candidate.UInt53(observation)); err != nil {
		if errors.Is(err, ErrCommitOutcomeIndeterminate) {
			return StoredRegistration{}, err
		}
		return StoredRegistration{}, classified(ClassificationLocalFailure, "time-high-water-write-failed")
	}
	state, err := component.store.snapshot(ctx)
	if err != nil {
		return StoredRegistration{}, classified(ClassificationLocalFailure, "registration-read-failed")
	}
	index := findRegistration(state.Registrations, registrationID)
	if index < 0 {
		return StoredRegistration{}, classified(ClassificationBinding, "registration-binding-mismatch")
	}
	entry := state.Registrations[index]
	if entry.Index.InstallationID != state.InstallationID ||
		entry.Index.EpochSequence != state.EpochSequence ||
		entry.Index.EpochDigest != state.EpochDigest ||
		entry.Index.SupervisorID != state.SupervisorID {
		return StoredRegistration{}, classified(ClassificationBinding, "registration-binding-mismatch")
	}
	if err := validateTrustAdmission(&state); err != nil {
		return StoredRegistration{}, err
	}
	if uint64(entry.Index.ExpiresAt) <= uint64(state.TimeHighWaterUnixSeconds) {
		return StoredRegistration{}, classified(ClassificationStale, "registration-expired")
	}
	return cloneStoredRegistration(entry.Record), nil
}

func authenticateRegistrationCall(call AuthenticatedCallContext) error {
	if !call.Authenticated || call.Role != CallerDaemon || call.Purpose != RegisterPlanPurpose {
		return classified(ClassificationAuthentication, "registration-caller-not-authorized")
	}
	return nil
}

func validatePlanStateBinding(state *installationState, plan v0candidate.ExecutionPlan) error {
	if plan.InstallationID != state.InstallationID ||
		plan.EpochSequence != state.EpochSequence ||
		plan.EpochDigest != state.EpochDigest {
		return classified(ClassificationBinding, "plan-installation-epoch-binding-mismatch")
	}
	return nil
}

func validateTrustAdmission(state *installationState) error {
	if state.Quarantined || state.TrustPhase == TrustRepairRequired {
		return classified(ClassificationTrustState, "registration-trust-state-denied")
	}
	if state.TrustPhase != TrustStable {
		return classified(ClassificationStale, "registration-transition-fenced")
	}
	return nil
}

func findRegistration(entries []registrationEntry, id v0candidate.RegistrationID) int {
	for index := range entries {
		if entries[index].Index.RegistrationID == id {
			return index
		}
	}
	return -1
}

func countUnexpired(entries []registrationEntry, effectiveNow v0candidate.UInt53) int {
	count := 0
	for _, entry := range entries {
		if uint64(entry.Index.ExpiresAt) > uint64(effectiveNow) {
			count++
		}
	}
	return count
}
