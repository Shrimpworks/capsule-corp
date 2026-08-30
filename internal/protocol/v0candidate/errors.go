package v0candidate

import (
	"errors"
	"fmt"

	"capsule.local/capsule/internal/protocol/classification"
)

// Classification is the internal conformance classification owned by the
// first strict boundary that rejects candidate bytes. It is not a public error
// code. It is the shared internal/protocol/classification vocabulary — see
// that package for the full cross-package classification set. This package
// only ever produces the five values below.
type Classification = classification.Classification

const (
	ClassificationMalformed   = classification.Malformed
	ClassificationUnsupported = classification.Unsupported
	ClassificationSchema      = classification.Schema
	ClassificationSemantic    = classification.Semantic
	ClassificationDomain      = classification.Domain
)

// DecodeError reports a bounded internal classification without including
// received bytes or decoded caller-controlled text.
type DecodeError struct {
	classification Classification
	detail         string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("%s: %s", e.classification, e.detail)
}

// Classification returns the first owning conformance classification.
func (e *DecodeError) Classification() Classification {
	return e.classification
}

// ErrorClassification extracts a strict-decoder classification.
func ErrorClassification(err error) (Classification, bool) {
	var decodeError *DecodeError
	if !errors.As(err, &decodeError) {
		return "", false
	}
	return decodeError.classification, true
}

func classified(classification Classification, detail string) error {
	return &DecodeError{classification: classification, detail: detail}
}

func malformed(detail string) error {
	return classified(ClassificationMalformed, detail)
}

func unsupported(detail string) error {
	return classified(ClassificationUnsupported, detail)
}

func schema(detail string) error {
	return classified(ClassificationSchema, detail)
}

func semantic(detail string) error {
	return classified(ClassificationSemantic, detail)
}

func domain(detail string) error {
	return classified(ClassificationDomain, detail)
}
