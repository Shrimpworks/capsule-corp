package registrationstate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

type retainedLookupKeys struct {
	registrationID v0candidate.RegistrationID
	approvalID     approvalattempt.ApprovalID
	attemptID      approvalattempt.AttemptID
	nonce          approvalattempt.AttemptNonce
	effectID       lifecyclestate.EffectID
	instanceDigest lifecyclestate.BackendInstanceDigest
	payloadDigest  approvalattempt.ApprovalPayloadDigest
	payload        []byte
	authorization  approvalattempt.ApprovalKeyAuthorizationIdentity
}

type retainedLookupResults struct {
	registration  RetainedRegistration
	approval      approvalattempt.ApprovalRecord
	attempt       RetainedAttempt
	nonce         RetainedNonce
	effect        RetainedEffect
	instance      RetainedInstance
	approvalID    approvalattempt.ApprovalID
	approvalState approvalattempt.ApprovalState
	attemptID     approvalattempt.AttemptID
	attemptState  approvalattempt.AttemptState
}

func lookupKeysForStore(t *testing.T, store *FixedFileStoreV2) retainedLookupKeys {
	t.Helper()
	if len(store.state.Registrations) != 1 || len(store.state.Approvals) != 1 ||
		len(store.state.Attempts) != 1 || len(store.lifecycles) != 1 {
		t.Fatal("lookup fixture is not the expected one-cohort world")
	}
	approval := store.state.Approvals[0]
	attempt := store.state.Attempts[0]
	lifecycle := store.lifecycles[0].View()
	return retainedLookupKeys{
		registrationID: store.state.Registrations[0].Index.RegistrationID,
		approvalID:     approval.ApprovalID, attemptID: attempt.AttemptID, nonce: approval.AttemptNonce,
		effectID: lifecycle.EffectID, instanceDigest: lifecycle.Instance.Digest(),
		payloadDigest: approval.PayloadDigest, payload: bytes.Clone(approval.ExactPayloadBytes),
		authorization: approval.AuthorizationIdentity,
	}
}

func resolveLookupResults(t *testing.T, store *FixedFileStoreV2, keys retainedLookupKeys) retainedLookupResults {
	t.Helper()
	ctx := context.Background()
	registration, err := store.ResolveRegistration(ctx, keys.registrationID)
	if err != nil {
		t.Fatal(err)
	}
	approval, err := store.ResolveApproval(ctx, keys.approvalID)
	if err != nil {
		t.Fatal(err)
	}
	attempt, err := store.ResolveAttempt(ctx, keys.attemptID)
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := store.ResolveNonce(ctx, keys.nonce)
	if err != nil {
		t.Fatal(err)
	}
	effect, err := store.ResolveEffect(ctx, keys.effectID)
	if err != nil {
		t.Fatal(err)
	}
	instance, err := store.ResolveInstance(ctx, keys.instanceDigest)
	if err != nil {
		t.Fatal(err)
	}
	approvalReference, approvalState, err := store.ResolveApprovalReplay(ctx, keys.payloadDigest, keys.payload, keys.authorization)
	if err != nil {
		t.Fatal(err)
	}
	attemptReference, attemptState, err := store.ResolveAttemptReplay(ctx, keys.registrationID, keys.approvalID)
	if err != nil {
		t.Fatal(err)
	}
	return retainedLookupResults{
		registration: registration, approval: approval, attempt: attempt, nonce: nonce,
		effect: effect, instance: instance, approvalID: approvalReference.ApprovalID(),
		approvalState: approvalState, attemptID: attemptReference.AttemptID(), attemptState: attemptState,
	}
}

func activateLookupStore(t *testing.T, store *FixedFileStoreV2, owner ArchiveOwner) *FixedFileStoreV2 {
	t.Helper()
	verified := mustPreparedArchive(t, store, owner)
	reopened, err := store.ActivateArchive(context.Background(), owner, verified, nil)
	if err != nil {
		t.Fatal(err)
	}
	return reopened
}

func TestFixedStoreV2RetainedLookupsHotAndArchiveAreSemanticallyIdentical(t *testing.T) {
	_, hotStore, _ := newEligibleFixedStoreV2(t)
	hotKeys := lookupKeysForStore(t, hotStore)
	hot := resolveLookupResults(t, hotStore, hotKeys)

	archivePath, archiveStore, owner := newEligibleFixedStoreV2(t)
	archiveKeys := lookupKeysForStore(t, archiveStore)
	if !reflect.DeepEqual(hotKeys, archiveKeys) {
		t.Fatal("deterministic hot and archive fixtures disagree")
	}
	archivedStore := activateLookupStore(t, archiveStore, owner)
	activeBefore := mustReadFile(t, archivePath)
	segmentPath := archiveSegmentPath(archivePath, archivedStore.segments[0].Segment.Digest())
	segmentBefore := mustReadFile(t, segmentPath)
	archived := resolveLookupResults(t, archivedStore, archiveKeys)
	if !reflect.DeepEqual(hot, archived) {
		t.Fatalf("hot and archived closed projections differ\nhot=%#v\narchived=%#v", hot, archived)
	}
	if archived.approvalID != archiveKeys.approvalID || archived.approvalState != approvalattempt.ApprovalConsumed ||
		archived.attemptID != archiveKeys.attemptID || archived.attemptState != approvalattempt.AttemptCreated {
		t.Fatal("archived replay did not return original retained identities and state")
	}
	recovery, err := archivedStore.RecoveryAttemptIDs(context.Background())
	if err != nil || len(recovery) != 0 {
		t.Fatalf("archived terminal attempt entered recovery: %#v, %v", recovery, err)
	}
	if !bytes.Equal(activeBefore, mustReadFile(t, archivePath)) || !bytes.Equal(segmentBefore, mustReadFile(t, segmentPath)) {
		t.Fatal("read-only lookup changed retained bytes")
	}

	reopened, err := OpenFixedFileStoreV2(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	if again := resolveLookupResults(t, reopened, archiveKeys); !reflect.DeepEqual(again, archived) {
		t.Fatal("reopen changed retained lookup projection")
	}

	archived.registration.Record.ExactPlanBytes[0] ^= 0xff
	archived.registration.PlanRoleBindings.ProfileReviewAttestationDigests[0][0] ^= 0xff
	archived.approval.ExactPayloadBytes[0] ^= 0xff
	archived.nonce.Approval.ExactEnvelopeBytes[0] ^= 0xff
	if defensive := resolveLookupResults(t, reopened, archiveKeys); !reflect.DeepEqual(defensive, hot) {
		t.Fatal("retained lookup result aliased caller-mutated bytes")
	}
}

func TestFixedStoreV2RetainedCollisionAndReplayQueries(t *testing.T) {
	_, store, owner := newEligibleFixedStoreV2(t)
	keys := lookupKeysForStore(t, store)
	candidates := RetainedIdentityCandidates{
		RegistrationID: keys.registrationID, ApprovalID: keys.approvalID, AttemptID: keys.attemptID,
		AttemptNonce: keys.nonce, EffectID: keys.effectID, InstanceDigest: keys.instanceDigest,
		ApprovalReplayPayloadDigest: keys.payloadDigest, ApprovalReplayAuthorizationIdentity: keys.authorization,
		AttemptReplayRegistrationID: keys.registrationID, AttemptReplayApprovalID: keys.approvalID,
	}
	hotCollisions, err := store.QueryRetainedIdentityCollisions(context.Background(), candidates)
	if err != nil {
		t.Fatal(err)
	}
	store = activateLookupStore(t, store, owner)
	collisions, err := store.QueryRetainedIdentityCollisions(context.Background(), candidates)
	if err != nil || collisions != (RetainedIdentityCollisions{
		Registration: true, Approval: true, Attempt: true, Nonce: true, Effect: true, Instance: true,
		ApprovalReplay: true, AttemptReplay: true,
	}) {
		t.Fatalf("archived collisions = %#v, %v", collisions, err)
	}
	if collisions != hotCollisions {
		t.Fatalf("hot/archive collision routing differs: %#v %#v", hotCollisions, collisions)
	}
	if err := store.CheckRetainedIdentityAvailability(context.Background(), candidates); err == nil {
		t.Fatal("archived identities were reported available")
	} else if classification, ok := approvalattempt.ErrorClassification(err); !ok || classification != approvalattempt.ClassificationReplay {
		t.Fatalf("collision classification = %q, %v", classification, err)
	}

	novel := candidates
	novel.RegistrationID[0] ^= 0xff
	novel.ApprovalID[0] ^= 0xff
	novel.AttemptID[0] ^= 0xff
	novel.AttemptNonce[0] ^= 0xff
	novel.EffectID[0] ^= 0xff
	novel.InstanceDigest[0] ^= 0xff
	novel.ApprovalReplayPayloadDigest[0] ^= 0xff
	novel.AttemptReplayRegistrationID[0] ^= 0xff
	if got, queryErr := store.QueryRetainedIdentityCollisions(context.Background(), novel); queryErr != nil || got.Any() {
		t.Fatalf("novel collision query = %#v, %v", got, queryErr)
	}
	if err := store.CheckRetainedIdentityAvailability(context.Background(), novel); err != nil {
		t.Fatalf("novel identities unavailable: %v", err)
	}
	if _, err := store.QueryRetainedIdentityCollisions(context.Background(), RetainedIdentityCandidates{
		ApprovalReplayPayloadDigest: keys.payloadDigest,
	}); err == nil {
		t.Fatal("partial approval replay identity accepted")
	}

	reference, state, err := store.ResolveApprovalReplay(context.Background(), keys.payloadDigest, keys.payload, keys.authorization)
	if err != nil || reference.ApprovalID() != keys.approvalID || state != approvalattempt.ApprovalConsumed {
		t.Fatalf("exact archived approval replay = %#v %q %v", reference, state, err)
	}
	collisionPayload := bytes.Clone(keys.payload)
	collisionPayload[0] ^= 0xff
	if _, _, err := store.ResolveApprovalReplay(context.Background(), keys.payloadDigest, collisionPayload, keys.authorization); err == nil {
		t.Fatal("approval payload digest collision accepted")
	} else if classification, ok := approvalattempt.ErrorClassification(err); !ok || classification != approvalattempt.ClassificationReplay {
		t.Fatalf("payload collision classification = %q, %v", classification, err)
	}
	wrongAuthorization := keys.authorization
	wrongAuthorization[0] ^= 0xff
	if _, _, err := store.ResolveApprovalReplay(context.Background(), keys.payloadDigest, keys.payload, wrongAuthorization); err == nil {
		t.Fatal("approval replay authorization substitution accepted")
	} else if classification, ok := approvalattempt.ErrorClassification(err); !ok || classification != approvalattempt.ClassificationBinding {
		t.Fatalf("authorization substitution classification = %q, %v", classification, err)
	}
}

func TestFixedStoreV2RetainedLookupMissingAndCorruptLocationsFailClosed(t *testing.T) {
	_, store, owner := newEligibleFixedStoreV2(t)
	keys := lookupKeysForStore(t, store)
	store = activateLookupStore(t, store, owner)
	missingRegistration, missingApproval, missingAttempt := keys.registrationID, keys.approvalID, keys.attemptID
	missingNonce, missingEffect, missingInstance := keys.nonce, keys.effectID, keys.instanceDigest
	missingRegistration[0] ^= 0xff
	missingApproval[0] ^= 0xff
	missingAttempt[0] ^= 0xff
	missingNonce[0] ^= 0xff
	missingEffect[0] ^= 0xff
	missingInstance[0] ^= 0xff
	missingCases := []struct {
		name string
		call func() error
	}{
		{name: "registration", call: func() error {
			_, err := store.ResolveRegistration(context.Background(), missingRegistration)
			return err
		}},
		{name: "approval", call: func() error { _, err := store.ResolveApproval(context.Background(), missingApproval); return err }},
		{name: "attempt", call: func() error { _, err := store.ResolveAttempt(context.Background(), missingAttempt); return err }},
		{name: "nonce", call: func() error { _, err := store.ResolveNonce(context.Background(), missingNonce); return err }},
		{name: "effect", call: func() error { _, err := store.ResolveEffect(context.Background(), missingEffect); return err }},
		{name: "instance", call: func() error { _, err := store.ResolveInstance(context.Background(), missingInstance); return err }},
		{name: "approval replay", call: func() error {
			_, _, err := store.ResolveApprovalReplay(context.Background(), approvalattempt.ApprovalPayloadDigest(missingInstance), keys.payload, keys.authorization)
			return err
		}},
		{name: "attempt replay", call: func() error {
			_, _, err := store.ResolveAttemptReplay(context.Background(), missingRegistration, keys.approvalID)
			return err
		}},
	}
	for _, test := range missingCases {
		if err := test.call(); err == nil {
			t.Fatalf("missing %s lookup succeeded", test.name)
		} else if classification, ok := approvalattempt.ErrorClassification(err); !ok || classification != approvalattempt.ClassificationBinding {
			t.Fatalf("missing %s classification = %q, %v", test.name, classification, err)
		}
	}

	tests := []struct {
		name   string
		mutate func(*diskEnvelopeV2)
	}{
		{name: "stale ordinal", mutate: func(envelope *diskEnvelopeV2) {
			envelope.Indexes.Registrations[0].Location.Archive.RecordOrdinal++
		}},
		{name: "wrong record kind", mutate: func(envelope *diskEnvelopeV2) {
			envelope.Indexes.Registrations[0].Location.RecordKind = archivestate.RecordApproval
		}},
		{name: "wrong index domain", mutate: func(envelope *diskEnvelopeV2) {
			envelope.Indexes.Scope = archivestate.ArchiveIndexScopeSegmentDerived
		}},
		{name: "missing replay tombstone", mutate: func(envelope *diskEnvelopeV2) {
			envelope.Indexes.ApprovalReplay = envelope.Indexes.ApprovalReplay[:0]
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, candidate, candidateOwner := newEligibleFixedStoreV2(t)
			candidateKeys := lookupKeysForStore(t, candidate)
			candidate = activateLookupStore(t, candidate, candidateOwner)
			var envelope diskEnvelopeV2
			if err := decodeOneClosedJSON(mustReadFile(t, path), &envelope); err != nil {
				t.Fatal(err)
			}
			test.mutate(&envelope)
			writeActiveLookupFixture(t, path, envelope)
			before := mustReadFile(t, path)
			if _, err := candidate.ResolveRegistration(context.Background(), candidateKeys.registrationID); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
				t.Fatalf("corrupt retained lookup error = %v", err)
			}
			if !bytes.Equal(before, mustReadFile(t, path)) {
				t.Fatal("corrupt lookup rewrote active evidence")
			}
		})
	}
}

func TestFixedStoreV2RecoveryAttemptIDsPreserveHotAbsenceAndExcludeArchive(t *testing.T) {
	state, _ := stateAndLifecycleRecord(t)
	path := newV1PathFromState(t, state, nil)
	hot := mustMigrateV2(t, path)
	recovery, err := hot.RecoveryAttemptIDs(context.Background())
	if err != nil || len(recovery) != 1 || recovery[0] != state.Attempts[0].AttemptID {
		t.Fatalf("hot absent-lifecycle recovery = %#v, %v", recovery, err)
	}

	_, archive, owner := newEligibleFixedStoreV2(t)
	archivedAttempt := archive.state.Attempts[0].AttemptID
	archive = activateLookupStore(t, archive, owner)
	recovery, err = archive.RecoveryAttemptIDs(context.Background())
	if err != nil || len(recovery) != 0 {
		t.Fatalf("archived terminal recovery = %#v, %v", recovery, err)
	}
	resolved, err := archive.ResolveAttempt(context.Background(), archivedAttempt)
	if err != nil || !resolved.LifecyclePresent || resolved.Lifecycle.View().State != lifecyclestate.StateDestroyed {
		t.Fatalf("archived terminal attempt lookup = %#v, %v", resolved, err)
	}
}

func TestFixedStoreV2RetainedLookupReadFaultCapSubstitutionAndRestoration(t *testing.T) {
	path, store, owner := newEligibleFixedStoreV2(t)
	keys := lookupKeysForStore(t, store)
	store = activateLookupStore(t, store, owner)
	segmentPath := archiveSegmentPath(path, store.segments[0].Segment.Digest())
	original := mustReadFile(t, segmentPath)
	heldPath := filepath.Join(t.TempDir(), "held-segment.json")
	if err := os.Rename(segmentPath, heldPath); err != nil {
		t.Fatal(err)
	}
	if _, err := store.ResolveRegistration(context.Background(), keys.registrationID); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
		t.Fatalf("missing segment lookup error = %v", err)
	}
	if err := os.Rename(heldPath, segmentPath); err != nil {
		t.Fatal(err)
	}
	if _, err := store.ResolveRegistration(context.Background(), keys.registrationID); err != nil {
		t.Fatalf("lookup after deliberate missing-segment restoration: %v", err)
	}

	if err := os.Chmod(segmentPath, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Truncate(segmentPath, archivestate.MaxSupervisorArchiveBytes+1); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(segmentPath, 0o400); err != nil {
		t.Fatal(err)
	}
	if _, err := store.ResolveRegistration(context.Background(), keys.registrationID); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
		t.Fatalf("cap-plus-one segment lookup error = %v", err)
	}
	replaceArchiveFixtureFile(t, segmentPath, original)
	if _, err := store.ResolveRegistration(context.Background(), keys.registrationID); err != nil {
		t.Fatalf("lookup after deliberate cap restoration: %v", err)
	}

	_, otherStore, otherOwner := newEligibleFixedStoreV2At(t, 1_785_456_301)
	other := mustPreparedArchive(t, otherStore, otherOwner)
	if bytes.Equal(other.SegmentBytes(), original) {
		t.Fatal("substitution fixture did not change")
	}
	replaceArchiveFixtureFile(t, segmentPath, other.SegmentBytes())
	if _, err := store.ResolveRegistration(context.Background(), keys.registrationID); err == nil || !errors.Is(err, ErrStoreRepairRequired) {
		t.Fatalf("substituted segment lookup error = %v", err)
	}
	replaceArchiveFixtureFile(t, segmentPath, original)
	if restored := resolveLookupResults(t, store, keys); restored.registration.RegistrationID != keys.registrationID {
		t.Fatal("deliberate segment restoration did not restore exact lookup")
	}
}

func TestFixedStoreV2RetainedLookupConcurrencyCancellationAndNoRewrite(t *testing.T) {
	path, store, owner := newEligibleFixedStoreV2(t)
	keys := lookupKeysForStore(t, store)
	store = activateLookupStore(t, store, owner)
	segmentPath := archiveSegmentPath(path, store.segments[0].Segment.Digest())
	activeBefore, segmentBefore := mustReadFile(t, path), mustReadFile(t, segmentPath)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := store.ResolveRegistration(ctx, keys.registrationID); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled lookup error = %v", err)
	}
	if _, err := store.ResolveRegistration(context.Background(), keys.registrationID); err != nil {
		t.Fatalf("retry after cancelled read: %v", err)
	}

	var wait sync.WaitGroup
	errorsSeen := make(chan error, 32)
	for range 32 {
		wait.Add(1)
		go func() {
			defer wait.Done()
			if _, err := store.ResolveEffect(context.Background(), keys.effectID); err != nil {
				errorsSeen <- err
				return
			}
			if _, _, err := store.ResolveAttemptReplay(context.Background(), keys.registrationID, keys.approvalID); err != nil {
				errorsSeen <- err
				return
			}
			if recovery, err := store.RecoveryAttemptIDs(context.Background()); err != nil || len(recovery) != 0 {
				if err == nil {
					err = errors.New("archived attempt entered recovery")
				}
				errorsSeen <- err
			}
		}()
	}
	wait.Wait()
	close(errorsSeen)
	for err := range errorsSeen {
		t.Fatal(err)
	}
	if !bytes.Equal(activeBefore, mustReadFile(t, path)) || !bytes.Equal(segmentBefore, mustReadFile(t, segmentPath)) {
		t.Fatal("concurrent read-only lookups changed retained files")
	}
}

func writeActiveLookupFixture(t *testing.T, path string, envelope diskEnvelopeV2) {
	t.Helper()
	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}
