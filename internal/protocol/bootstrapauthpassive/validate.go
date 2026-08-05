package bootstrapauthpassive

import (
	"crypto/sha256"
)

func requestFromWire(w requestWire) (Request, error) {
	installationID, err := installationIDFrom(w.InstallationID)
	if err != nil {
		return Request{}, err
	}
	publicKey, err := decodePublicKey(w.InstallationRootPublicKey)
	if err != nil {
		return Request{}, err
	}
	keyID, err := keyIDFrom(w.SignerKeyID)
	if err != nil {
		return Request{}, err
	}
	auth, err := authFrom(w.SignerAuthorizationIdentity)
	if err != nil {
		return Request{}, err
	}
	supervisorID, err := supervisorIDFrom(w.SupervisorID)
	if err != nil {
		return Request{}, err
	}
	release, err := digestFrom(w.ContainingReleaseDigest, "containing-release-digest")
	if err != nil {
		return Request{}, err
	}
	i1, err := digestFrom(w.I1ProfileCandidateDigest, "i1-profile-candidate-digest")
	if err != nil {
		return Request{}, err
	}
	component, err := digestFrom(w.ComponentProfileDigest, "component-profile-digest")
	if err != nil {
		return Request{}, err
	}
	manifest, err := digestFrom(w.InstallationManifestCandidateDigest, "installation-manifest-candidate-digest")
	if err != nil {
		return Request{}, err
	}
	epoch, err := digestFrom(w.TrustEpochCandidateDigest, "trust-epoch-candidate-digest")
	if err != nil {
		return Request{}, err
	}
	nonce, err := nonceFrom(w.Nonce)
	if err != nil {
		return Request{}, err
	}
	return Request{
		ObjectType: w.ObjectType, ObjectVersion: w.ObjectVersion, Purpose: w.Purpose, Audience: w.Audience,
		InstallationID: installationID, InstallationRootPublicKey: publicKey, SignerKeyID: keyID,
		SignerAuthorizationIdentity: auth, SupervisorID: supervisorID, ExpectedEffectiveUID: w.ExpectedEffectiveUID,
		ContainingReleaseDigest: release, I1ProfileCandidateDigest: i1, ComponentProfileDigest: component,
		InstallationManifestCandidateDigest: manifest, TrustEpochSequence: w.TrustEpochSequence,
		TrustEpochCandidateDigest: epoch, StateRootClass: w.StateRootClass, StateRootName: w.StateRootName,
		RequestEntryName: w.RequestEntryName, RecordEntryName: w.RecordEntryName, OwnerLockEntryName: w.OwnerLockEntryName,
		OwnerLockMode: w.OwnerLockMode, OwnerLockLinkCount: w.OwnerLockLinkCount, OwnerLockProfile: w.OwnerLockProfile,
		StoreEntryName: w.StoreEntryName, StoreFormat: w.StoreFormat, StoreProfile: w.StoreProfile,
		AttemptsEnabled: w.AttemptsEnabled, FakeBackendCreatesGuest: w.FakeBackendCreatesGuest, Nonce: nonce,
		IssuedAt: w.IssuedAt, NotBefore: w.NotBefore, ExpiresAt: w.ExpiresAt, ReplayDisposition: w.ReplayDisposition,
	}, nil
}

func recordFromWire(w recordWire) (Record, error) {
	requestPayload, err := digestFrom(w.RequestPayloadDigest, "request-payload-digest")
	if err != nil {
		return Record{}, err
	}
	requestEnvelope, err := digestFrom(w.RequestEnvelopeDigest, "request-envelope-digest")
	if err != nil {
		return Record{}, err
	}
	requestNonce, err := nonceFrom(w.RequestNonce)
	if err != nil {
		return Record{}, err
	}
	installationID, err := installationIDFrom(w.InstallationID)
	if err != nil {
		return Record{}, err
	}
	publicKey, err := decodePublicKey(w.InstallationRootPublicKey)
	if err != nil {
		return Record{}, err
	}
	keyID, err := keyIDFrom(w.SignerKeyID)
	if err != nil {
		return Record{}, err
	}
	auth, err := authFrom(w.SignerAuthorizationIdentity)
	if err != nil {
		return Record{}, err
	}
	supervisorID, err := supervisorIDFrom(w.SupervisorID)
	if err != nil {
		return Record{}, err
	}
	release, err := digestFrom(w.ContainingReleaseDigest, "containing-release-digest")
	if err != nil {
		return Record{}, err
	}
	i1, err := digestFrom(w.I1ProfileCandidateDigest, "i1-profile-candidate-digest")
	if err != nil {
		return Record{}, err
	}
	component, err := digestFrom(w.ComponentProfileDigest, "component-profile-digest")
	if err != nil {
		return Record{}, err
	}
	manifest, err := digestFrom(w.InstallationManifestCandidateDigest, "installation-manifest-candidate-digest")
	if err != nil {
		return Record{}, err
	}
	epoch, err := digestFrom(w.TrustEpochCandidateDigest, "trust-epoch-candidate-digest")
	if err != nil {
		return Record{}, err
	}
	initialStore, err := digestFrom(w.InitialStoreDigest, "initial-store-digest")
	if err != nil {
		return Record{}, err
	}
	return Record{
		ObjectType: w.ObjectType, ObjectVersion: w.ObjectVersion, Purpose: w.Purpose, Audience: w.Audience,
		RequestPayloadDigest: requestPayload, RequestEnvelopeDigest: requestEnvelope, RequestEnvelopeLength: w.RequestEnvelopeLength,
		RequestNonce: requestNonce, InstallationID: installationID, InstallationRootPublicKey: publicKey,
		SignerKeyID: keyID, SignerAuthorizationIdentity: auth, SupervisorID: supervisorID, ExpectedEffectiveUID: w.ExpectedEffectiveUID,
		ContainingReleaseDigest: release, I1ProfileCandidateDigest: i1, ComponentProfileDigest: component,
		InstallationManifestCandidateDigest: manifest, TrustEpochSequence: w.TrustEpochSequence, TrustEpochCandidateDigest: epoch,
		StateRootClass: w.StateRootClass, StateRootName: w.StateRootName, StateRootDevice: w.StateRootDevice,
		StateRootInode: w.StateRootInode, StateRootType: w.StateRootType, StateRootUID: w.StateRootUID, StateRootMode: w.StateRootMode,
		OwnerLockEntryName: w.OwnerLockEntryName, OwnerLockDevice: w.OwnerLockDevice, OwnerLockInode: w.OwnerLockInode,
		OwnerLockType: w.OwnerLockType, OwnerLockUID: w.OwnerLockUID, OwnerLockMode: w.OwnerLockMode,
		OwnerLockLinkCount: w.OwnerLockLinkCount, OwnerLockProfile: w.OwnerLockProfile, StoreEntryName: w.StoreEntryName,
		StoreDevice: w.StoreDevice, StoreType: w.StoreType, StoreUID: w.StoreUID, StoreMode: w.StoreMode,
		StoreLinkCount: w.StoreLinkCount, StoreFormat: w.StoreFormat, StoreProfile: w.StoreProfile,
		InitialStoreByteLength: w.InitialStoreByteLength, InitialStoreDigest: initialStore, RequestEntryName: w.RequestEntryName,
		RecordEntryName: w.RecordEntryName, RequestRetentionPolicy: w.RequestRetentionPolicy,
		RecordRetentionPolicy: w.RecordRetentionPolicy, InstallationTransition: w.InstallationTransition, Transition: w.Transition,
		FinalizedAt: w.FinalizedAt, AttemptsEnabled: w.AttemptsEnabled, RuntimeAdmitted: w.RuntimeAdmitted,
		BackendAdmitted: w.BackendAdmitted, GuestCreated: w.GuestCreated, OneUseDisposition: w.OneUseDisposition,
		ReplayDisposition: w.ReplayDisposition,
		RequestIssuedAt:   w.RequestIssuedAt, RequestNotBefore: w.RequestNotBefore,
		RequestExpiresAt: w.RequestExpiresAt, RequestReplayDisposition: w.RequestReplayDisposition,
		RequestObjectType: w.RequestObjectType, RequestObjectVersion: w.RequestObjectVersion,
		RequestPurpose: w.RequestPurpose, RequestAudience: w.RequestAudience,
		IssuedAt: w.IssuedAt, NotBefore: w.NotBefore, ExpiresAt: w.ExpiresAt,
	}, nil
}

func validateRequestShape(v Request) error {
	if v.ObjectType != RequestObjectType || v.ObjectVersion != ObjectVersion {
		return rejected(ClassificationUnsupported, "request-domain")
	}
	if v.Purpose != RequestPurpose || v.Audience != RequestAudience {
		return rejected(ClassificationBinding, "request-purpose-audience")
	}
	if v.InstallationID == (InstallationID{}) || v.SupervisorID == (SupervisorID{}) || v.Nonce == (Nonce{}) {
		return rejected(ClassificationSchema, "request-zero-identifier-or-nonce")
	}
	if v.ExpectedEffectiveUID == 0 || v.ExpectedEffectiveUID > 4294967295 {
		return rejected(ClassificationSchema, "request-effective-uid")
	}
	if v.TrustEpochSequence != TrustEpochSequence {
		return rejected(ClassificationBinding, "request-epoch-sequence")
	}
	if v.StateRootClass != StateRootClass || v.StateRootName != StateRootName || v.RequestEntryName != RequestEntryName || v.RecordEntryName != RecordEntryName ||
		v.OwnerLockEntryName != OwnerLockEntryName || v.OwnerLockMode != OwnerLockMode || v.OwnerLockLinkCount != OwnerLockLinkCount || v.OwnerLockProfile != OwnerLockProfile ||
		v.StoreEntryName != StoreEntryName || v.StoreFormat != StoreFormat || v.StoreProfile != StoreProfile {
		return rejected(ClassificationBinding, "request-root-lock-store-profile")
	}
	if v.AttemptsEnabled || v.FakeBackendCreatesGuest {
		return rejected(ClassificationBinding, "request-passive-no-guest")
	}
	if v.ReplayDisposition != RequestReplayDisposition {
		return rejected(ClassificationBinding, "request-replay-disposition")
	}
	if v.IssuedAt > v.NotBefore || v.NotBefore >= v.ExpiresAt || v.ExpiresAt-v.IssuedAt > 300 {
		return rejected(ClassificationTime, "request-time-window")
	}
	keyBytes := encodePublicKey(v.InstallationRootPublicKey)
	if sha256.Sum256(keyBytes) != v.SignerKeyID {
		return rejected(ClassificationBinding, "request-key-id")
	}
	if authorizationIdentity(v.SignerKeyID, v.InstallationID, v.TrustEpochSequence, v.TrustEpochCandidateDigest) != v.SignerAuthorizationIdentity {
		return rejected(ClassificationBinding, "request-authorization-identity")
	}
	return nil
}

func validateRecordShape(v Record) error {
	if v.ObjectType != RecordObjectType || v.ObjectVersion != ObjectVersion {
		return rejected(ClassificationUnsupported, "record-domain")
	}
	if v.Purpose != RecordPurpose || v.Audience != RecordAudience {
		return rejected(ClassificationBinding, "record-purpose-audience")
	}
	if v.InstallationID == (InstallationID{}) || v.SupervisorID == (SupervisorID{}) || v.RequestNonce == (Nonce{}) {
		return rejected(ClassificationSchema, "record-zero-identifier-or-nonce")
	}
	if v.ExpectedEffectiveUID == 0 || v.ExpectedEffectiveUID > 4294967295 || v.TrustEpochSequence != TrustEpochSequence {
		return rejected(ClassificationBinding, "record-identity-or-epoch")
	}
	if v.RequestEnvelopeLength == 0 || v.RequestEnvelopeLength > RequestEnvelopeRawMaxBytes {
		return rejected(ClassificationSchema, "record-request-envelope-length")
	}
	if v.StateRootClass != StateRootClass || v.StateRootName != StateRootName || v.StateRootType != "directory" || v.StateRootMode != 0o700 || v.StateRootDevice == 0 || v.StateRootInode == 0 ||
		v.OwnerLockEntryName != OwnerLockEntryName || v.OwnerLockType != "regular-file" || v.OwnerLockMode != OwnerLockMode || v.OwnerLockLinkCount != 1 || v.OwnerLockProfile != OwnerLockProfile || v.OwnerLockDevice == 0 || v.OwnerLockInode == 0 ||
		v.StoreEntryName != StoreEntryName || v.StoreType != "regular-file" || v.StoreMode != 0o600 || v.StoreLinkCount != 1 || v.StoreFormat != StoreFormat || v.StoreProfile != StoreProfile || v.StoreDevice == 0 {
		return rejected(ClassificationBinding, "record-root-lock-store-profile")
	}
	if v.StateRootUID != v.ExpectedEffectiveUID || v.OwnerLockUID != v.ExpectedEffectiveUID || v.StoreUID != v.ExpectedEffectiveUID || v.OwnerLockDevice != v.StateRootDevice || v.StoreDevice != v.StateRootDevice {
		return rejected(ClassificationBinding, "record-filesystem-observation")
	}
	if v.RequestEntryName != RequestEntryName || v.RecordEntryName != RecordEntryName || v.RequestRetentionPolicy != RequestRetentionPolicy || v.RecordRetentionPolicy != RecordRetentionPolicy ||
		v.InstallationTransition != InstallationTransition || v.Transition != ProtectedRootTransition || v.OneUseDisposition != RecordOneUseDisposition || v.ReplayDisposition != RecordReplayDisposition {
		return rejected(ClassificationBinding, "record-transition-retention-replay")
	}
	if v.AttemptsEnabled || v.RuntimeAdmitted || v.BackendAdmitted || v.GuestCreated {
		return rejected(ClassificationBinding, "record-passive-no-guest")
	}
	if v.RequestIssuedAt > v.RequestNotBefore || v.RequestNotBefore >= v.RequestExpiresAt || v.RequestExpiresAt-v.RequestIssuedAt > 300 ||
		v.FinalizedAt < v.RequestNotBefore || v.FinalizedAt > v.IssuedAt || v.IssuedAt > v.NotBefore || v.NotBefore >= v.ExpiresAt ||
		v.ExpiresAt-v.IssuedAt > 300 || v.ExpiresAt > v.RequestExpiresAt || v.RequestReplayDisposition != RequestReplayDisposition ||
		v.RequestObjectType != RequestObjectType || v.RequestObjectVersion != ObjectVersion || v.RequestPurpose != RequestPurpose || v.RequestAudience != RequestAudience {
		return rejected(ClassificationTime, "record-request-time-replay")
	}
	keyBytes := encodePublicKey(v.InstallationRootPublicKey)
	if sha256.Sum256(keyBytes) != v.SignerKeyID || authorizationIdentity(v.SignerKeyID, v.InstallationID, v.TrustEpochSequence, v.TrustEpochCandidateDigest) != v.SignerAuthorizationIdentity {
		return rejected(ClassificationBinding, "record-key-authorization")
	}
	return nil
}

func bindRecordToRequest(record Record, request Request) error {
	if record.RequestNonce != request.Nonce ||
		record.InstallationID != request.InstallationID ||
		record.InstallationRootPublicKey != request.InstallationRootPublicKey ||
		record.SignerKeyID != request.SignerKeyID ||
		record.SignerAuthorizationIdentity != request.SignerAuthorizationIdentity ||
		record.SupervisorID != request.SupervisorID ||
		record.ExpectedEffectiveUID != request.ExpectedEffectiveUID ||
		record.ContainingReleaseDigest != request.ContainingReleaseDigest ||
		record.I1ProfileCandidateDigest != request.I1ProfileCandidateDigest ||
		record.ComponentProfileDigest != request.ComponentProfileDigest ||
		record.InstallationManifestCandidateDigest != request.InstallationManifestCandidateDigest ||
		record.TrustEpochSequence != request.TrustEpochSequence ||
		record.TrustEpochCandidateDigest != request.TrustEpochCandidateDigest ||
		record.StateRootClass != request.StateRootClass ||
		record.StateRootName != request.StateRootName ||
		record.OwnerLockEntryName != request.OwnerLockEntryName ||
		record.OwnerLockMode != request.OwnerLockMode ||
		record.OwnerLockLinkCount != request.OwnerLockLinkCount ||
		record.OwnerLockProfile != request.OwnerLockProfile ||
		record.StoreEntryName != request.StoreEntryName ||
		record.StoreFormat != request.StoreFormat ||
		record.StoreProfile != request.StoreProfile ||
		record.RequestEntryName != request.RequestEntryName ||
		record.RecordEntryName != request.RecordEntryName ||
		record.RequestIssuedAt != request.IssuedAt ||
		record.RequestNotBefore != request.NotBefore ||
		record.RequestExpiresAt != request.ExpiresAt ||
		record.RequestReplayDisposition != request.ReplayDisposition ||
		record.RequestObjectType != request.ObjectType ||
		record.RequestObjectVersion != request.ObjectVersion ||
		record.RequestPurpose != request.Purpose ||
		record.RequestAudience != request.Audience {
		return rejected(ClassificationBinding, "request-stable-field-binding")
	}
	return nil
}

func evaluateRequestReplay(payloadDigest Digest, nonce Nonce, state ReplayState) (Decision, error) {
	switch state.Disposition {
	case ReplayFresh:
		return DecisionAdmitOnce, nil
	case ReplayPending:
		if state.PayloadDigest == payloadDigest && state.Nonce == nonce {
			return DecisionResumeExact, nil
		}
	case ReplayCompleted:
		return "", rejected(ClassificationReplay, "request-already-enrolled")
	default:
		return "", rejected(ClassificationReplay, "request-replay-state")
	}
	return "", rejected(ClassificationReplay, "request-substitution-or-reused-nonce")
}

func evaluateRecordReplay(payloadDigest Digest, nonce Nonce, state ReplayState) (Decision, error) {
	switch state.Disposition {
	case ReplayFresh, ReplayPending:
		return DecisionCommitOnce, nil
	case ReplayCompleted:
		if state.PayloadDigest == payloadDigest && state.Nonce == nonce {
			return DecisionReturnRetained, nil
		}
	default:
		return "", rejected(ClassificationReplay, "record-replay-state")
	}
	return "", rejected(ClassificationReplay, "record-substitution-or-reused-nonce")
}

func digestFrom(value []byte, code string) (Digest, error) {
	if len(value) != 32 {
		return Digest{}, rejected(ClassificationSchema, code)
	}
	var result Digest
	copy(result[:], value)
	return result, nil
}
func keyIDFrom(value []byte) (KeyID, error) {
	d, err := digestFrom(value, "signer-key-id")
	return KeyID(d), err
}
func authFrom(value []byte) (AuthorizationIdentity, error) {
	d, err := digestFrom(value, "signer-authorization-identity")
	return AuthorizationIdentity(d), err
}
func nonceFrom(value []byte) (Nonce, error) {
	d, err := digestFrom(value, "nonce")
	return Nonce(d), err
}
func installationIDFrom(value []byte) (InstallationID, error) {
	if len(value) != 16 {
		return InstallationID{}, rejected(ClassificationSchema, "installation-id")
	}
	var result InstallationID
	copy(result[:], value)
	return result, nil
}
func supervisorIDFrom(value []byte) (SupervisorID, error) {
	if len(value) != 16 {
		return SupervisorID{}, rejected(ClassificationSchema, "supervisor-id")
	}
	var result SupervisorID
	copy(result[:], value)
	return result, nil
}
