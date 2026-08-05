package authorityplane

import (
	"bytes"
	"context"
	"errors"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

type CallerRole string

const (
	Daemon CallerRole = "daemon"
	Broker CallerRole = "broker"
)
const (
	RegisterPlanV0Purpose      = "capsule.ipc.register-plan.v0"
	GetRegisteredPlanV0Purpose = "capsule.ipc.get-registered-plan.v0"
)

type CallContext struct {
	Authenticated bool
	Role          CallerRole
	Purpose       string
}

type RoleResolver interface {
	ResolvePlanRoles(context.Context, v0candidate.ExecutionPlanRoleBindings) (v0candidate.ExecutionPlanRoleBindings, error)
}
type IdentifierSource interface {
	NewRegistrationID(context.Context) (v0candidate.RegistrationID, error)
}

type SupervisorContext struct {
	InstallationID v0candidate.InstallationID
	EpochSequence  v0candidate.UInt53
	EpochDigest    v0candidate.TrustEpochDigest
	SupervisorID   v0candidate.SupervisorID
}

type Facade struct {
	store       *FixedStore
	roles       RoleResolver
	identifiers IdentifierSource
	supervisor  SupervisorContext
}

func NewFacade(store *FixedStore, roles RoleResolver, identifiers IdentifierSource, supervisor SupervisorContext) (*Facade, error) {
	if store == nil || roles == nil || identifiers == nil {
		return nil, errors.New("fixed store, role resolver, and identifier source are required")
	}
	return &Facade{store: store, roles: roles, identifiers: identifiers, supervisor: supervisor}, nil
}

// RegisterPlanV0 validates the complete request before one all-or-none store update.
func (f *Facade) RegisterPlanV0(ctx context.Context, call CallContext, request RegisterPlanV0Request) ([]byte, error) {
	if !call.Authenticated || call.Role != Daemon || call.Purpose != RegisterPlanV0Purpose {
		return nil, refused(Authentication, "register-plan-caller")
	}
	nominal, err := DecodeRoleBindingsV0(request.bindings)
	if err != nil {
		return nil, err
	}
	resolved, err := f.roles.ResolvePlanRoles(ctx, nominal)
	if err != nil {
		return nil, refused(Binding, "role-binding-resolution")
	}
	resolvedBytes, err := EncodeRoleBindingsV0(resolved)
	if err != nil {
		return nil, err
	}
	plan, err := v0candidate.DecodeExecutionPlan(request.plan, resolved)
	if err != nil {
		return nil, err
	}
	manifest, err := v0candidate.DecodeSourceManifest(request.manifest, v0candidate.SourceManifestMediaType, request.source)
	if err != nil {
		return nil, err
	}
	view := plan.View()
	if view.InstallationID != f.supervisor.InstallationID || view.EpochSequence != f.supervisor.EpochSequence || view.EpochDigest != f.supervisor.EpochDigest {
		return nil, refused(Binding, "plan-supervisor-state")
	}
	if view.SourceManifestDigest != manifest.Digest() || view.SourceEntrypoint != v0candidate.MJSMainPath || view.SourceByteLength != manifest.View().AggregateByteLength {
		return nil, refused(Binding, "plan-source-custody")
	}
	var registration []byte
	err = f.store.update(ctx, func(records *[]retainedRegistration, sequence *v0candidate.UInt53) error {
		if len(*records) >= MaxRetainedRegistrations {
			return refused(Capacity, "registration-capacity")
		}
		if uint64(*sequence) == v0candidate.MaxSafeInteger {
			return refused(Capacity, "registration-sequence")
		}
		id, err := f.identifiers.NewRegistrationID(ctx)
		if err != nil || id == (v0candidate.RegistrationID{}) {
			return refused(LocalFailure, "registration-identifier")
		}
		for _, record := range *records {
			if record.id == id {
				return refused(LocalFailure, "duplicate-registration-identifier")
			}
		}
		*sequence = v0candidate.UInt53(uint64(*sequence) + 1)
		wire := v0candidate.PlanRegistration{ObjectType: v0candidate.PlanRegistrationObjectType, ObjectVersion: v0candidate.CandidateObjectVersion, RegistrationID: id, RegistrationSequence: v0candidate.PositiveUInt53(*sequence), PlanDigest: plan.Digest(), InstallationID: f.supervisor.InstallationID, EpochSequence: f.supervisor.EpochSequence, EpochDigest: f.supervisor.EpochDigest, SupervisorID: f.supervisor.SupervisorID, ExpiresAt: view.ExpiresAt}
		registration = encodeRegistration(wire)
		if _, err := v0candidate.DecodePlanRegistration(registration, v0candidate.PlanRegistrationRoleBindings{RegistrationID: id, PlanDigest: plan.Digest(), InstallationID: f.supervisor.InstallationID, EpochDigest: f.supervisor.EpochDigest, SupervisorID: f.supervisor.SupervisorID}); err != nil {
			return refused(LocalFailure, "registration-self-check")
		}
		*records = append(*records, retainedRegistration{id: id, plan: plan.AuthoritativeBytes(), bindings: resolvedBytes, registration: bytes.Clone(registration), manifest: manifest.AuthoritativeBytes(), source: manifest.SourceBytes()})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bytes.Clone(registration), nil
}

// GetRegisteredPlanV0 is a Broker-only repeatable read. It never mutates the store.
func (f *Facade) GetRegisteredPlanV0(_ context.Context, call CallContext, request GetRegisteredPlanV0Request) (GetRegisteredPlanV0Reply, error) {
	if !call.Authenticated || call.Role != Broker || call.Purpose != GetRegisteredPlanV0Purpose {
		return GetRegisteredPlanV0Reply{}, refused(Authentication, "get-registered-plan-caller")
	}
	if request.RegistrationID == (v0candidate.RegistrationID{}) {
		return GetRegisteredPlanV0Reply{}, refused(Schema, "registration-id")
	}
	record, ok := f.store.get(request.RegistrationID)
	if !ok {
		return GetRegisteredPlanV0Reply{}, refused(Binding, "registration-not-found")
	}
	if len(record.plan)+len(record.bindings)+len(record.registration)+len(record.manifest)+len(record.source) > GetRegisteredPlanV0MaxBytes {
		return GetRegisteredPlanV0Reply{}, refused(LocalFailure, "stored-reply-cap")
	}
	return GetRegisteredPlanV0Reply{plan: record.plan, bindings: record.bindings, registration: record.registration, manifest: record.manifest, source: record.source}, nil
}

func encodeRegistration(view v0candidate.PlanRegistration) []byte {
	result := appendCBORArgument(nil, 5, 10)
	result = appendUnsigned(result, 1)
	result = appendText(result, view.ObjectType)
	result = appendUnsigned(result, 2)
	result = appendUnsigned(result, uint64(view.ObjectVersion))
	result = appendUnsigned(result, 3)
	result = appendBytes(result, view.RegistrationID[:])
	result = appendUnsigned(result, 4)
	result = appendUnsigned(result, uint64(view.RegistrationSequence))
	result = appendUnsigned(result, 5)
	result = appendBytes(result, view.PlanDigest[:])
	result = appendUnsigned(result, 6)
	result = appendBytes(result, view.InstallationID[:])
	result = appendUnsigned(result, 7)
	result = appendUnsigned(result, uint64(view.EpochSequence))
	result = appendUnsigned(result, 8)
	result = appendBytes(result, view.EpochDigest[:])
	result = appendUnsigned(result, 9)
	result = appendBytes(result, view.SupervisorID[:])
	result = appendUnsigned(result, 10)
	result = appendUnsigned(result, uint64(view.ExpiresAt))
	return result
}

func EqualRoleBindings(left, right v0candidate.ExecutionPlanRoleBindings) bool {
	leftBytes, leftErr := EncodeRoleBindingsV0(left)
	rightBytes, rightErr := EncodeRoleBindingsV0(right)
	return leftErr == nil && rightErr == nil && bytes.Equal(leftBytes, rightBytes)
}
