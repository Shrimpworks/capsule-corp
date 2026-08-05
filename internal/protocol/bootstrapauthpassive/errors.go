package bootstrapauthpassive

import "fmt"

type Classification string

const (
	ClassificationMalformed   Classification = "MALFORMED"
	ClassificationUnsupported Classification = "UNSUPPORTED"
	ClassificationSchema      Classification = "SCHEMA"
	ClassificationBinding     Classification = "BINDING"
	ClassificationSignature   Classification = "SIGNATURE"
	ClassificationTime        Classification = "TIME"
	ClassificationReplay      Classification = "REPLAY"
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
