// Package classification is the shared error-classification vocabulary
// backing v0candidate, bootstrapauthpassive, and sourcevalidatorpassive's
// refusal codes. It exists so the same underlying concept — a syntactically
// malformed payload, a value outside a closed vocabulary, a value that is
// well-formed but semantically out of bounds, a cross-object identity
// mismatch, a value that disagrees with a caller-supplied binding, and so
// on — is spelled identically everywhere those packages classify a refusal,
// instead of each package independently reinventing its own words for the
// same concepts (and drifting, as v0candidate/bootstrapauthpassive/
// sourcevalidatorpassive previously had).
//
// No caller package is expected to produce every value here: each keeps
// its own closed subset (documented at its own Classification declaration)
// and its own error type carrying one. This package supplies only the
// shared, closed string vocabulary those subsets are drawn from.
package classification

// Classification is one word from the shared refusal-classification
// vocabulary.
type Classification string

const (
	// Malformed marks bytes that do not even parse as their declared wire
	// shape (truncated, non-canonical, or otherwise structurally invalid).
	Malformed Classification = "MALFORMED"
	// Unsupported marks a value outside a closed, enumerated vocabulary
	// (an unknown version, method, or field) that this boundary has no
	// defined behavior for.
	Unsupported Classification = "UNSUPPORTED"
	// Schema marks a value that parses but violates a structural
	// invariant of its own object (a required field's zero value, a
	// reserved byte that must be zero, and so on).
	Schema Classification = "SCHEMA"
	// Semantic marks a value that is well-formed and schema-valid but
	// semantically out of bounds (for example, content that exceeds a
	// declared maximum).
	Semantic Classification = "SEMANTIC"
	// Domain marks a value that disagrees with this boundary's own closed
	// domain/identity constants (a digest domain tag, an object-type
	// cross-link, a fixed frame magic, and so on).
	Domain Classification = "DOMAIN"
	// Binding marks a value that disagrees with a separately supplied,
	// caller-owned binding (a declared length, a digest, a correlation
	// binding between two independently decoded objects).
	Binding Classification = "BINDING"
	// Signature marks a cryptographic signature verification failure.
	Signature Classification = "SIGNATURE"
	// Time marks a time-bound (freshness/expiry) violation.
	Time Classification = "TIME"
	// Replay marks a replay-protection violation.
	Replay Classification = "REPLAY"
)
