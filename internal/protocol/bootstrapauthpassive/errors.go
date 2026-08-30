package bootstrapauthpassive

import (
	"fmt"

	"capsule.local/capsule/internal/protocol/classification"
)

// Classification is the shared internal/protocol/classification vocabulary
// — see that package for the full cross-package classification set. This
// package only ever produces the seven values below.
type Classification = classification.Classification

// Classification values this package can produce, aliased from the shared
// internal/protocol/classification vocabulary.
const (
	ClassificationMalformed   = classification.Malformed
	ClassificationUnsupported = classification.Unsupported
	ClassificationSchema      = classification.Schema
	ClassificationBinding     = classification.Binding
	ClassificationSignature   = classification.Signature
	ClassificationTime        = classification.Time
	ClassificationReplay      = classification.Replay
)

type VerificationError struct {
	classification Classification
	code           string
}

func (e *VerificationError) Error() string                  { return fmt.Sprintf("%s: %s", e.classification, e.code) }
func (e *VerificationError) Classification() Classification { return e.classification }
func (e *VerificationError) Code() string                   { return e.code }

func rejected(classification Classification, code string) error {
	return &VerificationError{classification: classification, code: code}
}
