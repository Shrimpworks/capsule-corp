package approvalattempt

import (
	"bytes"
	"testing"
)

func TestDomainIdentifierCopiesInputAndReferencesRejectZero(t *testing.T) {
	source := bytes.Repeat([]byte{0x41}, 16)
	value, err := NewDomainIdentifier(DomainApprovalID, source)
	if err != nil {
		t.Fatalf("new domain identifier: %v", err)
	}
	source[0] ^= 0xff
	identifier, err := NewApprovalID(value)
	if err != nil {
		t.Fatalf("new approval ID: %v", err)
	}
	if identifier[0] != 0x41 {
		t.Fatal("approval ID retained caller-owned storage")
	}
	if _, err := NewApprovalReference(ApprovalID{}); err == nil {
		t.Fatal("zero approval reference was accepted")
	}
	if _, err := NewAttemptReference(AttemptID{}); err == nil {
		t.Fatal("zero attempt reference was accepted")
	}
}

func TestClassificationProjectionIsDefensive(t *testing.T) {
	first := Classifications()
	first[0] = "MUTATED"
	second := Classifications()
	if second[0] != ClassificationAuthentication {
		t.Fatal("classification projection exposed mutable package state")
	}
}
