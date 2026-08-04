package registeredlifecycle

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/installationowner"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/execution/registrationstate"
)

// ownedStartupCheckpoint is a fixed test-only call-order oracle. It is not an
// IPC, readiness, execution-result, or evidence surface.
type ownedStartupCheckpoint string

const (
	ownedStartupRootOpened         ownedStartupCheckpoint = "state-root-opened"
	ownedStartupOwnerAcquired      ownedStartupCheckpoint = "owner-acquired"
	ownedStartupOwnerChecked       ownedStartupCheckpoint = "owner-checked-before-store"
	ownedStartupStoreOpened        ownedStartupCheckpoint = "store-opened"
	ownedStartupCoordinatorCreated ownedStartupCheckpoint = "coordinator-created"
	ownedStartupLifecycleCreated   ownedStartupCheckpoint = "lifecycle-created"
	ownedStartupRecoveryOwnerHeld  ownedStartupCheckpoint = "owner-checked-before-recovery"
	ownedStartupRecoveryEnumerated ownedStartupCheckpoint = "recovery-enumerated"
	ownedStartupBeforeRecover      ownedStartupCheckpoint = "before-recover"
	ownedStartupRecoveryComplete   ownedStartupCheckpoint = "recovery-complete"
	ownedStartupLifecycleStopped   ownedStartupCheckpoint = "lifecycle-stopped"
	ownedStartupStoreClosed        ownedStartupCheckpoint = "store-closed"
	ownedStartupOwnerClosed        ownedStartupCheckpoint = "owner-closed"
	ownedStartupStateRootClosed    ownedStartupCheckpoint = "state-root-closed"
)

// ownedStartupCheckpointHook exists only for deterministic local ordering and
// fault tests. Production composition leaves it nil.
type ownedStartupCheckpointHook func(context.Context, ownedStartupCheckpoint, approvalattempt.AttemptID) error

// OwnedStartupOptions contains only trusted installed bootstrap inputs and
// the current no-guest v1 lifecycle dependencies. StorePath is the trusted
// installation-configured sibling path used by the existing v1 conformance
// opener; no caller or IPC body can supply it through this internal API.
type OwnedStartupOptions struct {
	StateRootPath string
	Enrollment    installationowner.OwnerLockEnrollment
	StorePath     string
	Attempts      AttemptResolver
	Backend       *FakeBackend
	Clock         registrationstate.TrustedClock
	EffectIDs     registrationstate.LifecycleEffectIDSource
}

type ownedStartupTestHooks struct {
	lifecycle CheckpointHook
	startup   ownedStartupCheckpointHook
}

// OwnedStartup is the bounded G2 composition of one live installation owner,
// one current FixedFileStoreV1, one same-session coordinator, and the no-guest
// registered lifecycle. It becomes available only after the initial sorted
// recovery pass succeeds.
type OwnedStartup struct {
	mu              sync.Mutex
	root            *installationowner.TrustedStateRoot
	owner           installationowner.InstallationOwner
	store           *registrationstate.FixedFileStoreV1
	component       *Component
	checkpoint      ownedStartupCheckpointHook
	initialRecovery []Snapshot
	closed          bool
}

// OpenOwnedStartup acquires G1 ownership before the current v1 store is
// inspected, binds the owner session to both store and coordinator, then runs
// sorted AttemptID-only recovery before returning a live session. It performs
// no migration, archive, IPC, adapter selection, runtime, backend, or guest
// wiring beyond the supplied FakeBackend.
func OpenOwnedStartup(ctx context.Context, options OwnedStartupOptions) (*OwnedStartup, error) {
	return openOwnedStartup(ctx, options, ownedStartupTestHooks{})
}

func openOwnedStartup(
	ctx context.Context,
	options OwnedStartupOptions,
	testHooks ownedStartupTestHooks,
) (*OwnedStartup, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if options.StateRootPath == "" || options.StorePath == "" || options.Attempts == nil ||
		options.Backend == nil || options.Clock == nil {
		return nil, errors.New("owned startup bootstrap, store, resolver, fake backend, and clock are required")
	}
	if options.Backend.CreatesGuest() || options.Backend.Binding().View().CreatesGuest {
		return nil, errors.New("owned startup requires the no-guest fake backend")
	}
	startup := &OwnedStartup{checkpoint: testHooks.startup}
	fail := func(primary error) (*OwnedStartup, error) {
		return nil, errors.Join(primary, startup.closeOwnedResources(context.Background(), false))
	}

	root, err := installationowner.OpenDarwinTrustedStateRoot(ctx, options.StateRootPath, options.Enrollment)
	if err != nil {
		return nil, err
	}
	startup.root = root
	if err := startup.observe(ctx, ownedStartupRootOpened, approvalattempt.AttemptID{}); err != nil {
		return fail(err)
	}
	owner, err := installationowner.AcquireDarwinInstallationOwner(ctx, root, options.Enrollment)
	if err != nil {
		return fail(err)
	}
	startup.owner = owner
	if err := startup.observe(ctx, ownedStartupOwnerAcquired, approvalattempt.AttemptID{}); err != nil {
		return fail(err)
	}
	if err := owner.CheckHeld(ctx); err != nil {
		return fail(fmt.Errorf("owned startup pre-store owner check: %w", err))
	}
	if err := startup.observe(ctx, ownedStartupOwnerChecked, approvalattempt.AttemptID{}); err != nil {
		return fail(err)
	}

	store, err := registrationstate.OpenFixedFileStoreV1WithOwner(
		ctx,
		options.StorePath,
		owner,
		registrationstate.FixedFileStoreV1Options{
			EffectIDs: options.EffectIDs, OwnerSessionID: owner.OwnerSessionID(),
		},
	)
	if err != nil {
		return fail(err)
	}
	startup.store = store
	if err := startup.observe(ctx, ownedStartupStoreOpened, approvalattempt.AttemptID{}); err != nil {
		return fail(err)
	}
	coordinator, err := NewCoordinator(owner.OwnerSessionID())
	if err != nil {
		return fail(err)
	}
	if err := startup.observe(ctx, ownedStartupCoordinatorCreated, approvalattempt.AttemptID{}); err != nil {
		return fail(err)
	}
	component, err := New(Options{
		Attempts: options.Attempts, Store: store, Backend: options.Backend,
		Coordinator: coordinator, Clock: options.Clock, Checkpoint: testHooks.lifecycle,
	})
	if err != nil {
		return fail(err)
	}
	startup.component = component
	if err := startup.observe(ctx, ownedStartupLifecycleCreated, approvalattempt.AttemptID{}); err != nil {
		return fail(err)
	}
	recovered, err := startup.recoverCreatedAttemptsLocked(ctx)
	if err != nil {
		return fail(err)
	}
	startup.initialRecovery = cloneSnapshots(recovered)
	return startup, nil
}

// OwnerSessionID exposes only the already-nominal in-process session value so
// tests can prove the owner, store, and coordinator used one identity. It does
// not expose or grant access to the owner descriptor.
func (startup *OwnedStartup) OwnerSessionID() lifecyclestate.OwnerSessionID {
	if startup == nil {
		return lifecyclestate.OwnerSessionID{}
	}
	startup.mu.Lock()
	defer startup.mu.Unlock()
	if startup.closed || startup.owner == nil || startup.store == nil {
		return lifecyclestate.OwnerSessionID{}
	}
	if startup.owner.OwnerSessionID() != startup.store.OwnerSessionID() {
		return lifecyclestate.OwnerSessionID{}
	}
	return startup.owner.OwnerSessionID()
}

// InitialRecovery returns a defensive copy of the clean initial recovery pass.
func (startup *OwnedStartup) InitialRecovery() []Snapshot {
	if startup == nil {
		return nil
	}
	startup.mu.Lock()
	defer startup.mu.Unlock()
	return cloneSnapshots(startup.initialRecovery)
}

// RecoverCreatedAttempts repeats the sorted recovery pass only while this
// owned session remains live. Any failed owner check permanently fences the
// underlying store and requires shutdown plus full reopen.
func (startup *OwnedStartup) RecoverCreatedAttempts(ctx context.Context) ([]Snapshot, error) {
	if startup == nil {
		return nil, registrationstate.ErrStoreClosed
	}
	startup.mu.Lock()
	defer startup.mu.Unlock()
	if startup.closed || startup.component == nil || startup.store == nil || startup.owner == nil {
		return nil, registrationstate.ErrStoreClosed
	}
	return startup.recoverCreatedAttemptsLocked(ctx)
}

func (startup *OwnedStartup) recoverCreatedAttemptsLocked(ctx context.Context) ([]Snapshot, error) {
	if err := startup.owner.CheckHeld(ctx); err != nil {
		// RecoveryAttemptIDs performs the store-owned recheck and installs the
		// permanent fence. Invoke it even after this composition check fails.
		_, storeErr := startup.store.RecoveryAttemptIDs(ctx)
		return nil, errors.Join(fmt.Errorf("owned startup recovery owner check: %w", err), storeErr)
	}
	if err := startup.observe(ctx, ownedStartupRecoveryOwnerHeld, approvalattempt.AttemptID{}); err != nil {
		return nil, err
	}
	identifiers, err := startup.store.RecoveryAttemptIDs(ctx)
	if err != nil {
		return nil, err
	}
	if !sortedUniqueAttemptIDs(identifiers) {
		return nil, fmt.Errorf("%w: startup recovery identifiers are not sorted and unique", registrationstate.ErrStoreRepairRequired)
	}
	if err := startup.observe(ctx, ownedStartupRecoveryEnumerated, approvalattempt.AttemptID{}); err != nil {
		return nil, err
	}
	results := make([]Snapshot, 0, len(identifiers))
	var firstErr error
	for _, attemptID := range identifiers {
		if err := startup.owner.CheckHeld(ctx); err != nil {
			_, storeErr := startup.store.RecoveryAttemptIDs(ctx)
			return results, errors.Join(fmt.Errorf("owned startup per-attempt owner check: %w", err), storeErr)
		}
		if err := startup.observe(ctx, ownedStartupBeforeRecover, attemptID); err != nil {
			return results, err
		}
		snapshot, recoverErr := startup.component.Recover(ctx, attemptID)
		results = append(results, snapshot)
		if firstErr == nil && recoverErr != nil {
			firstErr = recoverErr
		}
	}
	if firstErr != nil {
		return results, firstErr
	}
	if err := startup.observe(ctx, ownedStartupRecoveryComplete, approvalattempt.AttemptID{}); err != nil {
		return results, err
	}
	return results, nil
}

// CloseAfterShutdown disables lifecycle work, closes the logical store handle,
// then releases the owner descriptor and finally the retained root descriptor.
func (startup *OwnedStartup) CloseAfterShutdown(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if startup == nil {
		return registrationstate.ErrStoreClosed
	}
	startup.mu.Lock()
	defer startup.mu.Unlock()
	if startup.closed {
		return registrationstate.ErrStoreClosed
	}
	return startup.closeOwnedResources(ctx, true)
}

func (startup *OwnedStartup) closeOwnedResources(ctx context.Context, markClosed bool) error {
	var result error
	if startup.component != nil {
		startup.component = nil
		result = errors.Join(result, startup.observe(ctx, ownedStartupLifecycleStopped, approvalattempt.AttemptID{}))
	}
	if startup.store != nil {
		if err := startup.store.CloseAfterLifecycleShutdown(ctx); err != nil {
			return errors.Join(result, err)
		}
		startup.store = nil
		result = errors.Join(result, startup.observe(ctx, ownedStartupStoreClosed, approvalattempt.AttemptID{}))
	}
	if startup.owner != nil {
		if err := startup.owner.CloseAfterShutdown(ctx); err != nil {
			return errors.Join(result, err)
		}
		startup.owner = nil
		result = errors.Join(result, startup.observe(ctx, ownedStartupOwnerClosed, approvalattempt.AttemptID{}))
	}
	if startup.root != nil {
		if err := startup.root.Close(); err != nil {
			return errors.Join(result, err)
		}
		startup.root = nil
		result = errors.Join(result, startup.observe(ctx, ownedStartupStateRootClosed, approvalattempt.AttemptID{}))
	}
	if markClosed {
		startup.closed = true
	}
	return result
}

func (startup *OwnedStartup) observe(
	ctx context.Context,
	checkpoint ownedStartupCheckpoint,
	attemptID approvalattempt.AttemptID,
) error {
	if startup.checkpoint == nil {
		return nil
	}
	return startup.checkpoint(ctx, checkpoint, attemptID)
}

func sortedUniqueAttemptIDs(identifiers []approvalattempt.AttemptID) bool {
	for index := 1; index < len(identifiers); index++ {
		if bytes.Compare(identifiers[index-1][:], identifiers[index][:]) >= 0 {
			return false
		}
	}
	return true
}

func cloneSnapshots(snapshots []Snapshot) []Snapshot {
	return append([]Snapshot(nil), snapshots...)
}
