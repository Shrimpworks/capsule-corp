package teardownpassive

// LatchStop retains the earliest accepted teardown trigger without consulting
// or mutating the pending storage operation.
func (model *Model) LatchStop(trigger Trigger, tick uint64) error {
	if model.state.RecoveryRequired {
		return ErrRecovery
	}
	if model.state.TerminalConfirmed {
		return ErrState
	}
	if !validTrigger(trigger) || !validTick(tick) {
		return ErrClock
	}
	switch trigger {
	case TriggerCancellation:
		if !model.state.PreparationStarted {
			return ErrState
		}
	case TriggerWallDeadline:
		wallTick, ok := addTick(model.state.Start.Tick, model.clock.WallAfterStartMS)
		if !model.state.Start.Attempted || !ok || tick != wallTick {
			return ErrClock
		}
	case TriggerFatalFault:
		if model.state.Custody != CustodyExact {
			return ErrState
		}
	}
	if !model.state.Stop.Latched || tick < model.state.Stop.ActionTick {
		model.state.Stop = StopSnapshot{Latched: true, Trigger: trigger, ActionTick: tick}
		model.updateTiming()
	}
	return nil
}

// AttemptStart records the one start-token attempt and its original monotonic
// anchor. The passive model writes no token and starts no process.
func (model *Model) AttemptStart(tick uint64) error {
	if model.state.RecoveryRequired {
		return ErrRecovery
	}
	if !validTick(tick) {
		return ErrClock
	}
	if !model.Decision().MayStart {
		return ErrState
	}
	if _, ok := addTick(tick, model.clock.WallAfterStartMS); !ok {
		return ErrClock
	}
	model.state.Start = StartSnapshot{Attempted: true, Tick: tick}
	return nil
}

// RequestSignal records the sole current-lifetime destructive request. It does
// not signal a process and never waits for a storage result.
func (model *Model) RequestSignal(identity ProcessIdentity, tick uint64) error {
	if model.state.RecoveryRequired {
		return ErrRecovery
	}
	if !validTick(tick) {
		return ErrClock
	}
	if model.state.Custody != CustodyExact {
		return ErrState
	}
	if zero32(identity) || identity != model.state.ProcessIdentity {
		return ErrBinding
	}
	if !model.Decision().MaySignal || tick < model.state.Stop.ActionTick {
		return ErrState
	}
	if _, ok := addTick(tick, model.clock.ForceAbsenceMS); !ok {
		return ErrClock
	}
	model.state.Signal = SignalSnapshot{Requested: true, Attempts: 1, ForceTick: tick}
	model.updateTiming()
	return nil
}

// MarkSignalUncertain records loss of the sole request response. It does not
// restore permission to issue a second request.
func (model *Model) MarkSignalUncertain() error {
	if model.state.RecoveryRequired {
		return ErrRecovery
	}
	if !model.state.Signal.Requested || model.state.Signal.Uncertain || model.state.Absence.Observed {
		return ErrState
	}
	model.state.Signal.Uncertain = true
	return nil
}

// ObserveAbsence records one authoritative same-lifetime absence observation
// for the exact custody identity. It treats neither signal return nor EOF as absence.
func (model *Model) ObserveAbsence(identity ProcessIdentity, tick uint64) error {
	if model.state.RecoveryRequired {
		return ErrRecovery
	}
	if !validTick(tick) {
		return ErrClock
	}
	if model.state.Custody != CustodyExact || model.state.Absence.Observed {
		return ErrState
	}
	if zero32(identity) || identity != model.state.ProcessIdentity {
		return ErrBinding
	}
	if model.state.Stop.Latched && tick < model.state.Stop.ActionTick {
		return ErrClock
	}
	if model.state.Signal.Requested && tick < model.state.Signal.ForceTick {
		return ErrClock
	}
	model.state.Absence = AbsenceSnapshot{Observed: true, Tick: tick}
	model.state.Custody = CustodyNone
	model.updateTiming()
	return nil
}

// Restart retains durable facts and storage uncertainty but clears all
// current-lifetime custody, signal, absence, start, and clock observations.
func (model *Model) Restart() {
	if model.state.Pending.Active {
		model.state.LastSettledGeneration = model.state.Pending.Generation
		model.state.LastWriteOutcome = WriteIndeterminate
		model.state.Pending = PendingWrite{}
	}
	model.state.Custody = CustodyNone
	model.state.ProcessIdentity = ProcessIdentity{}
	model.state.Stop = StopSnapshot{}
	model.state.Signal = SignalSnapshot{}
	model.state.Absence = AbsenceSnapshot{}
	model.state.Start = StartSnapshot{}
	model.state.Timing = TimingUnknown
	model.state.RecoveryRequired = model.state.TerminalDisposition != TerminalCompleted
}

func (model *Model) updateTiming() {
	if !model.state.Stop.Latched || !model.state.Absence.Observed {
		model.state.Timing = TimingUnknown
		return
	}
	actionLimit, actionOK := addTick(
		model.state.Stop.ActionTick,
		model.clock.ActionToAbsenceMS,
	)
	satisfied := actionOK && model.state.Absence.Tick <= actionLimit
	if model.state.Signal.Requested {
		forceLimit, forceOK := addTick(
			model.state.Signal.ForceTick,
			model.clock.ForceAbsenceMS,
		)
		satisfied = satisfied && forceOK && model.state.Absence.Tick <= forceLimit
	}
	if satisfied {
		model.state.Timing = TimingSatisfied
	} else {
		model.state.Timing = TimingViolated
	}
}

func validTick(tick uint64) bool { return tick > 0 && tick <= maxSafeTick }

func addTick(left, right uint64) (uint64, bool) {
	if !validTick(left) || right > maxSafeTick-left {
		return 0, false
	}
	return left + right, true
}
