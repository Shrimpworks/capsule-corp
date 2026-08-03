package archivestate

import (
	"errors"
	"fmt"
)

// Classification is the closed content-free F1 validation vocabulary. It is
// internal conformance output, not a public protocol error.
type Classification string

const (
	ClassificationMalformed    Classification = "MALFORMED"
	ClassificationUnsupported  Classification = "UNSUPPORTED"
	ClassificationSchema       Classification = "SCHEMA"
	ClassificationDomain       Classification = "DOMAIN"
	ClassificationBinding      Classification = "BINDING"
	ClassificationCapacity     Classification = "CAPACITY"
	ClassificationTrustState   Classification = "TRUST_STATE"
	ClassificationRepairNeeded Classification = "REPAIR_REQUIRED"
)

type contractError struct {
	classification Classification
	code           string
}

func (e *contractError) Error() string {
	return fmt.Sprintf("%s: %s", e.classification, e.code)
}

func (e *contractError) Classification() Classification { return e.classification }

func classified(classification Classification, code string) error {
	return &contractError{classification: classification, code: code}
}

// ErrorClassification extracts a fixed classification without returning
// caller-controlled record data.
func ErrorClassification(err error) (Classification, bool) {
	var candidate *contractError
	if !errors.As(err, &candidate) {
		return "", false
	}
	return candidate.classification, true
}
