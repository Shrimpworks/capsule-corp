package lifecyclestate

import (
	"errors"
	"fmt"
)

// Classification is the closed, content-free validation vocabulary for this
// passive contract. It is not a public protocol error or execution result.
type Classification string

const (
	ClassificationMalformed         Classification = "MALFORMED"
	ClassificationUnsupported       Classification = "UNSUPPORTED"
	ClassificationSchema            Classification = "SCHEMA"
	ClassificationDomain            Classification = "DOMAIN"
	ClassificationBinding           Classification = "BINDING"
	ClassificationStale             Classification = "STALE"
	ClassificationCapacity          Classification = "CAPACITY"
	ClassificationLocalFailure      Classification = "LOCAL_FAILURE"
	ClassificationLifecycleFailure  Classification = "LIFECYCLE_FAILURE"
	ClassificationCleanupUnresolved Classification = "CLEANUP_UNRESOLVED"
	ClassificationTrustState        Classification = "TRUST_STATE"
	ClassificationRecoveryRequired  Classification = "RECOVERY_REQUIRED"
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

// ErrorClassification extracts a fixed validation classification without
// exposing caller-controlled bytes or text.
func ErrorClassification(err error) (Classification, bool) {
	var candidate *contractError
	if !errors.As(err, &candidate) {
		return "", false
	}
	return candidate.classification, true
}
