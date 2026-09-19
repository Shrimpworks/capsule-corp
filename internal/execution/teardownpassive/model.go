package teardownpassive

// Model is the deterministic C5b18 no-effect state machine. All fields are
// in-memory specification state; none represents a product store or live process.
type Model struct {
	bindings Bindings
	clock    ClockPolicy
	state    Snapshot
	used     map[OperationID]struct{}
}

// NewModel validates the immutable cleanup-only bindings and returns an empty
// passive successor state.
func NewModel(bindings Bindings) (*Model, error) {
	if !bindings.valid() {
		return nil, ErrBinding
	}
	return &Model{
		bindings: bindings,
		clock:    fixedClockPolicy,
		state: Snapshot{
			Contract:            ContractIdentity,
			RecordVersion:       RecordVersionV1,
			LastSettled:         SettledWrite{Outcome: WriteNone},
			Custody:             CustodyNone,
			Timing:              TimingUnknown,
			TerminalDisposition: TerminalNone,
		},
		used: make(map[OperationID]struct{}, 3),
	}, nil
}

// Bindings returns the immutable bound identities by value.
func (model *Model) Bindings() Bindings { return model.bindings }

// ClockPolicy returns the exact fixed v1 policy by value.
func (model *Model) ClockPolicy() ClockPolicy { return model.clock }

// Snapshot returns a defensive value projection of all passive state.
func (model *Model) Snapshot() Snapshot { return model.state }

// Decision derives the current closed permission projection. Pending storage
// blocks create/start/publication but never an already prepared stop path.
func (model *Model) Decision() Decision {
	snapshot := model.state
	blocked := snapshot.RecoveryRequired || snapshot.TerminalConfirmed
	return Decision{
		MayCreate: snapshot.Prepared && !snapshot.CreationConsumed && !blocked &&
			snapshot.Custody == CustodyNone &&
			!snapshot.Stop.Latched && !snapshot.Pending.Active,
		MayStart: snapshot.Prepared && !blocked && snapshot.Custody == CustodyExact &&
			snapshot.RunnerIdentityConfirmed && !snapshot.Stop.Latched &&
			!snapshot.Start.Attempted && !snapshot.Pending.Active,
		MaySignal: !blocked && snapshot.Custody == CustodyExact && snapshot.Stop.Latched &&
			!snapshot.Signal.Requested && !snapshot.Absence.Observed,
		MayPublishTerminal: !blocked && snapshot.Absence.Observed &&
			!snapshot.TerminalWriteStarted && !snapshot.Pending.Active,
		MayRelease: snapshot.TerminalConfirmed &&
			snapshot.TerminalDisposition == TerminalCompleted,
		StopIndependentOfStorage: !blocked && snapshot.Custody == CustodyExact && snapshot.Stop.Latched,
		RecoveryRequired:         snapshot.RecoveryRequired,
	}
}

// BeginPreparation starts the sole versioned preparation operation. Its ID must
// equal the immutable Supervisor-generated preparation identity.
func (model *Model) BeginPreparation(operationID OperationID) (PendingWrite, error) {
	if model.state.RecoveryRequired {
		return PendingWrite{}, ErrRecovery
	}
	if model.state.PreparationStarted || operationID != model.bindings.PreparationOperationID {
		return PendingWrite{}, ErrState
	}
	candidate := model.durableProjection()
	candidate.Prepared = true
	pending, err := model.beginWrite(operationID, WritePreparation, candidate)
	if err != nil {
		return PendingWrite{}, err
	}
	model.state.PreparationStarted = true
	return pending, nil
}

// ObserveCreatedCustody records one exact process identity observed after
// confirmed preparation in the current Supervisor lifetime. It creates nothing.
func (model *Model) ObserveCreatedCustody(identity ProcessIdentity) error {
	if model.state.RecoveryRequired {
		return ErrRecovery
	}
	if zero32(identity) {
		return ErrBinding
	}
	if !model.Decision().MayCreate {
		return ErrState
	}
	model.state.Custody = CustodyExact
	model.state.ProcessIdentity = identity
	model.state.CreationConsumed = true
	return nil
}

// BeginRunnerIdentityWrite freezes the one runner-publication operation while
// start remains closed.
func (model *Model) BeginRunnerIdentityWrite(operationID OperationID) (PendingWrite, error) {
	if model.state.RecoveryRequired {
		return PendingWrite{}, ErrRecovery
	}
	if model.state.Pending.Active {
		return PendingWrite{}, ErrPendingWrite
	}
	if !model.state.Prepared || model.state.Custody != CustodyExact ||
		model.state.RunnerWriteStarted || model.state.Stop.Latched ||
		model.state.TerminalConfirmed {
		return PendingWrite{}, ErrState
	}
	candidate := model.durableProjection()
	candidate.RunnerIdentityConfirmed = true
	pending, err := model.beginWrite(operationID, WriteRunnerIdentity, candidate)
	if err != nil {
		return PendingWrite{}, err
	}
	model.state.RunnerWriteStarted = true
	return pending, nil
}

// BeginTerminalWrite freezes the only terminal-record candidate. Authoritative
// absence must already be present; timing failure freezes an unresolved disposition.
func (model *Model) BeginTerminalWrite(operationID OperationID) (PendingWrite, error) {
	if model.state.RecoveryRequired {
		return PendingWrite{}, ErrRecovery
	}
	if model.state.Pending.Active {
		return PendingWrite{}, ErrPendingWrite
	}
	if !model.Decision().MayPublishTerminal {
		return PendingWrite{}, ErrState
	}
	candidate := model.durableProjection()
	candidate.TerminalConfirmed = true
	candidate.TerminalDisposition = TerminalCompleted
	if model.state.Timing == TimingViolated {
		candidate.TerminalDisposition = TerminalTimingViolated
	}
	pending, err := model.beginWrite(operationID, WriteTerminalJoin, candidate)
	if err != nil {
		return PendingWrite{}, err
	}
	model.state.TerminalWriteStarted = true
	return pending, nil
}

func (model *Model) beginWrite(
	operationID OperationID,
	kind WriteKind,
	candidate DurableProjection,
) (PendingWrite, error) {
	if model.state.Pending.Active {
		return PendingWrite{}, ErrPendingWrite
	}
	if zero16(operationID) {
		return PendingWrite{}, ErrBinding
	}
	if _, exists := model.used[operationID]; exists {
		return PendingWrite{}, ErrBinding
	}
	generation := model.state.LastSettled.Generation + 1
	if generation == 0 || generation > maxSafeTick {
		return PendingWrite{}, ErrState
	}
	pending := PendingWrite{
		Active: true, OperationID: operationID, Generation: generation, Kind: kind,
		Candidate: candidate,
	}
	model.used[operationID] = struct{}{}
	model.state.Pending = pending
	return pending, nil
}

// SettleWrite accepts only the exact outstanding token and one closed outcome.
// Late confirmation updates its durable fact without changing any control latch.
func (model *Model) SettleWrite(pending PendingWrite, outcome WriteOutcome) error {
	if !validWriteOutcome(outcome) {
		return ErrState
	}
	if !samePending(model.state.Pending, pending) {
		return ErrPendingWrite
	}
	model.state.Pending = PendingWrite{}
	model.state.LastSettled = settledWrite(pending, outcome)
	if outcome != WriteConfirmed {
		return nil
	}
	switch pending.Kind {
	case WritePreparation:
	case WriteRunnerIdentity:
	case WriteTerminalJoin:
	default:
		return ErrState
	}
	model.applyDurableProjection(pending.Candidate)
	if pending.Kind == WriteTerminalJoin {
		if pending.Candidate.TerminalDisposition == TerminalCompleted {
			model.state.OutputReleased = true
			model.state.CapacityReleased = true
		} else {
			model.state.RecoveryRequired = true
		}
	}
	return nil
}

func (model *Model) durableProjection() DurableProjection {
	return DurableProjection{
		Prepared:                model.state.Prepared,
		CreationConsumed:        model.state.CreationConsumed,
		RunnerIdentityConfirmed: model.state.RunnerIdentityConfirmed,
		TerminalConfirmed:       model.state.TerminalConfirmed,
		TerminalDisposition:     model.state.TerminalDisposition,
	}
}

func (model *Model) applyDurableProjection(candidate DurableProjection) {
	model.state.Prepared = candidate.Prepared
	model.state.CreationConsumed = candidate.CreationConsumed
	model.state.RunnerIdentityConfirmed = candidate.RunnerIdentityConfirmed
	model.state.TerminalConfirmed = candidate.TerminalConfirmed
	model.state.TerminalDisposition = candidate.TerminalDisposition
}

func samePending(left, right PendingWrite) bool {
	return left.Active && right.Active && left == right
}

func settledWrite(pending PendingWrite, outcome WriteOutcome) SettledWrite {
	return SettledWrite{
		Present:     true,
		OperationID: pending.OperationID,
		Generation:  pending.Generation,
		Kind:        pending.Kind,
		Candidate:   pending.Candidate,
		Outcome:     outcome,
	}
}
