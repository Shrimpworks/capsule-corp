package v0candidate

import (
	"errors"
	"fmt"
)

// Classification is the internal conformance classification owned by the
// first strict boundary that rejects candidate bytes. It is not a public error
// code.
type Classification string

const (
	ClassificationMalformed   Classification = "MALFORMED"
	ClassificationUnsupported Classification = "UNSUPPORTED"
	ClassificationSchema      Classification = "SCHEMA"
	ClassificationDomain      Classification = "DOMAIN"
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

func domain(detail string) error {
	return classified(ClassificationDomain, detail)
}
