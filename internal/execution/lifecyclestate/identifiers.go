package lifecyclestate

import "bytes"

const identifierBytes = 16

// IdentifierDomain is the trusted semantic role attached before converting
// opaque bytes into one of this package's nominal identifier types.
type IdentifierDomain string

const (
	// DomainEffectID labels bytes intended for the lifecycle effect namespace.
	DomainEffectID       IdentifierDomain = "lifecycle-effect"
	DomainOwnerSessionID IdentifierDomain = "lifecycle-owner-session"
)

// DomainIdentifier prevents correct-width bytes with the wrong trusted role
// from silently entering the effect or owner-session namespace.
type DomainIdentifier struct {
	domain IdentifierDomain
	bytes  [identifierBytes]byte
}

// NewDomainIdentifier validates a trusted identifier role and copies one
// exact nonzero 16-byte value into that role's namespace.
func NewDomainIdentifier(domain IdentifierDomain, value []byte) (DomainIdentifier, error) {
	if domain != DomainEffectID && domain != DomainOwnerSessionID {
		return DomainIdentifier{}, classified(ClassificationDomain, "unknown-identifier-domain")
	}
	if len(value) != identifierBytes {
		return DomainIdentifier{}, classified(ClassificationSchema, "identifier-length")
	}
	var identifier DomainIdentifier
	identifier.domain = domain
	copy(identifier.bytes[:], value)
	if identifier.bytes == ([identifierBytes]byte{}) {
		return DomainIdentifier{}, classified(ClassificationSchema, "identifier-zero")
	}
	return identifier, nil
}

// Domain returns the trusted semantic namespace attached at construction.
func (identifier DomainIdentifier) Domain() IdentifierDomain { return identifier.domain }

// Bytes returns a defensive copy of the exact identifier bytes.
func (identifier DomainIdentifier) Bytes() []byte { return bytes.Clone(identifier.bytes[:]) }

// EffectID is a distinct nonzero 16-byte Supervisor-issued namespace. It is
// not interchangeable with an approval, attempt, nonce, registration, owner
// session, or backend instance identity.
type EffectID [identifierBytes]byte

// NewEffectID converts only a domain-checked lifecycle-effect identifier.
func NewEffectID(identifier DomainIdentifier) (EffectID, error) {
	if identifier.domain != DomainEffectID {
		return EffectID{}, classified(ClassificationDomain, "effect-id-domain")
	}
	return EffectID(identifier.bytes), nil
}

// IsZero reports whether the effect identifier is absent.
func (identifier EffectID) IsZero() bool { return identifier == (EffectID{}) }

// OwnerSessionID seals a permit to one future in-process lifecycle owner. E1
// defines only the nominal value; it establishes no ownership mechanism.
type OwnerSessionID [identifierBytes]byte

// NewOwnerSessionID converts only a domain-checked lifecycle-owner identifier.
func NewOwnerSessionID(identifier DomainIdentifier) (OwnerSessionID, error) {
	if identifier.domain != DomainOwnerSessionID {
		return OwnerSessionID{}, classified(ClassificationDomain, "owner-session-id-domain")
	}
	return OwnerSessionID(identifier.bytes), nil
}

// IsZero reports whether the owner-session identifier is absent.
func (identifier OwnerSessionID) IsZero() bool { return identifier == (OwnerSessionID{}) }
