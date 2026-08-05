// Package bootstrapauthpassive freezes only the passive I2B1 request and
// record byte boundary. It contains no production signer, key access, IPC,
// filesystem, service, process, runtime, backend, or guest operation.
package bootstrapauthpassive

import "bytes"

const (
	RequestObjectType  = "capsule.supervisor-bootstrap-request"
	RequestPurpose     = "capsule.installation.bootstrap.request"
	RequestAudience    = "capsule.execution-supervisor.bootstrap"
	RequestMediaType   = "application/capsule.supervisor-bootstrap-request+cbor;v=0"
	RecordObjectType   = "capsule.supervisor-bootstrap-record"
	RecordPurpose      = "capsule.installation.bootstrap.record"
	RecordAudience     = "capsule.execution-supervisor"
	RecordMediaType    = "application/capsule.supervisor-bootstrap-record+cbor;v=0"
	ObjectVersion      = uint64(0)
	TrustEpochSequence = uint64(1)

	StateRootClass              = "supervisor-private-app-sandbox"
	StateRootName               = "supervisor.state"
	RequestEntryName            = "supervisor.bootstrap-request"
	RecordEntryName             = "supervisor.bootstrap-record"
	OwnerLockEntryName          = "supervisor.owner"
	OwnerLockMode               = uint64(0o600)
	OwnerLockLinkCount          = uint64(1)
	OwnerLockProfile            = "darwin-openat-flock-v0"
	StoreEntryName              = "supervisor.store"
	StoreFormat                 = "capsule.supervisor-store/fixed-v1"
	StoreProfile                = "conformance-non-product-no-guest"
	RequestReplayDisposition    = "one-shot-exact-payload-replay-only"
	RequestRetentionPolicy      = "retain-exact-envelope-mode-0400-single-link"
	RecordRetentionPolicy       = "retain-exact-envelope-file-and-epoch-anchor"
	InstallationTransition      = "installation-candidate-to-protected-root-disabled"
	ProtectedRootTransition     = "protected-root-validated-disabled"
	RecordOneUseDisposition     = "commit-once"
	RecordReplayDisposition     = "exact-retained-envelope-response-replay-only"
	AuthorizationIdentityDomain = "capsule.installation-root-key-authorization/v0\x00"

	RequestEnvelopeRawMaxBytes  = 4096
	RequestPayloadRawMaxBytes   = 2048
	RequestProtectedRawMaxBytes = 256
	RecordEnvelopeRawMaxBytes   = 6144
	RecordPayloadRawMaxBytes    = 4096
	RecordProtectedRawMaxBytes  = 256

	// These are patched only when the complete closed maximum shapes change.
	// Tests independently encode those shapes and require exact equality.
	RequestPayloadCalculatedMaxBytes   = 860
	RequestProtectedCalculatedMaxBytes = 98
	RequestEnvelopeCalculatedMaxBytes  = 1032
	RecordPayloadCalculatedMaxBytes    = 1526
	RecordProtectedCalculatedMaxBytes  = 97
	RecordEnvelopeCalculatedMaxBytes   = 1697
)

type Digest [32]byte
type InstallationID [16]byte
type SupervisorID [16]byte
type Nonce [32]byte
type KeyID [32]byte
type AuthorizationIdentity [32]byte

// field-authority-object: capsule.installation-root-public-key v0
type InstallationRootPublicKey struct {
	KeyType   int64
	Algorithm int64
	Curve     int64
	X         [32]byte
	Y         [32]byte
}

// field-authority-object: capsule.supervisor-bootstrap-request v0
type Request struct {
	ObjectType                          string
	ObjectVersion                       uint64
	Purpose                             string
	Audience                            string
	InstallationID                      InstallationID
	InstallationRootPublicKey           InstallationRootPublicKey
	SignerKeyID                         KeyID
	SignerAuthorizationIdentity         AuthorizationIdentity
	SupervisorID                        SupervisorID
	ExpectedEffectiveUID                uint64
	ContainingReleaseDigest             Digest
	I1ProfileCandidateDigest            Digest
	ComponentProfileDigest              Digest
	InstallationManifestCandidateDigest Digest
	TrustEpochSequence                  uint64
	TrustEpochCandidateDigest           Digest
	StateRootClass                      string
	StateRootName                       string
	RequestEntryName                    string
	RecordEntryName                     string
	OwnerLockEntryName                  string
	OwnerLockMode                       uint64
	OwnerLockLinkCount                  uint64
	OwnerLockProfile                    string
	StoreEntryName                      string
	StoreFormat                         string
	StoreProfile                        string
	AttemptsEnabled                     bool
	FakeBackendCreatesGuest             bool
	Nonce                               Nonce
	IssuedAt                            uint64
	NotBefore                           uint64
	ExpiresAt                           uint64
	ReplayDisposition                   string
}

// field-authority-object: capsule.supervisor-bootstrap-record v0
type Record struct {
	ObjectType                          string
	ObjectVersion                       uint64
	Purpose                             string
	Audience                            string
	RequestPayloadDigest                Digest
	RequestEnvelopeDigest               Digest
	RequestEnvelopeLength               uint64
	RequestNonce                        Nonce
	InstallationID                      InstallationID
	InstallationRootPublicKey           InstallationRootPublicKey
	SignerKeyID                         KeyID
	SignerAuthorizationIdentity         AuthorizationIdentity
	SupervisorID                        SupervisorID
	ExpectedEffectiveUID                uint64
	ContainingReleaseDigest             Digest
	I1ProfileCandidateDigest            Digest
	ComponentProfileDigest              Digest
	InstallationManifestCandidateDigest Digest
	TrustEpochSequence                  uint64
	TrustEpochCandidateDigest           Digest
	StateRootClass                      string
	StateRootName                       string
	StateRootDevice                     uint64
	StateRootInode                      uint64
	StateRootType                       string
	StateRootUID                        uint64
	StateRootMode                       uint64
	OwnerLockEntryName                  string
	OwnerLockDevice                     uint64
	OwnerLockInode                      uint64
	OwnerLockType                       string
	OwnerLockUID                        uint64
	OwnerLockMode                       uint64
	OwnerLockLinkCount                  uint64
	OwnerLockProfile                    string
	StoreEntryName                      string
	StoreDevice                         uint64
	StoreType                           string
	StoreUID                            uint64
	StoreMode                           uint64
	StoreLinkCount                      uint64
	StoreFormat                         string
	StoreProfile                        string
	InitialStoreByteLength              uint64
	InitialStoreDigest                  Digest
	RequestEntryName                    string
	RecordEntryName                     string
	RequestRetentionPolicy              string
	RecordRetentionPolicy               string
	InstallationTransition              string
	Transition                          string
	FinalizedAt                         uint64
	AttemptsEnabled                     bool
	RuntimeAdmitted                     bool
	BackendAdmitted                     bool
	GuestCreated                        bool
	OneUseDisposition                   string
	ReplayDisposition                   string
	RequestIssuedAt                     uint64
	RequestNotBefore                    uint64
	RequestExpiresAt                    uint64
	RequestReplayDisposition            string
	RequestObjectType                   string
	RequestObjectVersion                uint64
	RequestPurpose                      string
	RequestAudience                     string
	IssuedAt                            uint64
	NotBefore                           uint64
	ExpiresAt                           uint64
}

type ReplayDisposition string

const (
	ReplayFresh     ReplayDisposition = "fresh"
	ReplayPending   ReplayDisposition = "pending"
	ReplayCompleted ReplayDisposition = "completed"
)

type ReplayState struct {
	Disposition    ReplayDisposition
	PayloadDigest  Digest
	EnvelopeDigest Digest
	Nonce          Nonce
}

type Decision string

const (
	DecisionAdmitOnce      Decision = "admit-once"
	DecisionResumeExact    Decision = "resume-exact"
	DecisionCommitOnce     Decision = "commit-once"
	DecisionReturnRetained Decision = "return-retained-envelope"
)

type VerifiedRequest struct {
	exactEnvelope  []byte
	exactProtected []byte
	exactPayload   []byte
	envelopeDigest Digest
	payloadDigest  Digest
	view           Request
	decision       Decision
}

func (v *VerifiedRequest) ExactEnvelope() []byte  { return bytes.Clone(v.exactEnvelope) }
func (v *VerifiedRequest) ExactProtected() []byte { return bytes.Clone(v.exactProtected) }
func (v *VerifiedRequest) ExactPayload() []byte   { return bytes.Clone(v.exactPayload) }
func (v *VerifiedRequest) EnvelopeDigest() Digest { return v.envelopeDigest }
func (v *VerifiedRequest) PayloadDigest() Digest  { return v.payloadDigest }
func (v *VerifiedRequest) View() Request          { return v.view }
func (v *VerifiedRequest) Decision() Decision     { return v.decision }

type VerifiedRecord struct {
	exactEnvelope  []byte
	exactProtected []byte
	exactPayload   []byte
	envelopeDigest Digest
	payloadDigest  Digest
	view           Record
	decision       Decision
}

func (v *VerifiedRecord) ExactEnvelope() []byte  { return bytes.Clone(v.exactEnvelope) }
func (v *VerifiedRecord) ExactProtected() []byte { return bytes.Clone(v.exactProtected) }
func (v *VerifiedRecord) ExactPayload() []byte   { return bytes.Clone(v.exactPayload) }
func (v *VerifiedRecord) EnvelopeDigest() Digest { return v.envelopeDigest }
func (v *VerifiedRecord) PayloadDigest() Digest  { return v.payloadDigest }
func (v *VerifiedRecord) View() Record           { return v.view }
func (v *VerifiedRecord) Decision() Decision     { return v.decision }
