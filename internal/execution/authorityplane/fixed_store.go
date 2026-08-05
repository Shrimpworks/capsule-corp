package authorityplane

import (
	"bytes"
	"context"
	"sync"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

const MaxRetainedRegistrations = 4096

type retainedRegistration struct {
	id                                             v0candidate.RegistrationID
	plan, bindings, registration, manifest, source []byte
}

func cloneRecord(value retainedRegistration) retainedRegistration {
	value.plan = bytes.Clone(value.plan)
	value.bindings = bytes.Clone(value.bindings)
	value.registration = bytes.Clone(value.registration)
	value.manifest = bytes.Clone(value.manifest)
	value.source = bytes.Clone(value.source)
	return value
}

// FixedStore is the fixture/facade transaction oracle. update publishes one
// cloned state only after the mutation succeeds, so partial custody is unrepresentable.
type FixedStore struct {
	mu         sync.Mutex
	records    []retainedRegistration
	sequence   v0candidate.UInt53
	failCommit bool
}

func (s *FixedStore) update(ctx context.Context, mutation func(*[]retainedRegistration, *v0candidate.UInt53) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	working := make([]retainedRegistration, len(s.records))
	for index := range s.records {
		working[index] = cloneRecord(s.records[index])
	}
	sequence := s.sequence
	if err := mutation(&working, &sequence); err != nil {
		return err
	}
	if s.failCommit {
		return refused(LocalFailure, "fixed-store-commit-failed")
	}
	s.records, s.sequence = working, sequence
	return nil
}

func (s *FixedStore) get(id v0candidate.RegistrationID) (retainedRegistration, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, record := range s.records {
		if record.id == id {
			return cloneRecord(record), true
		}
	}
	return retainedRegistration{}, false
}

func (s *FixedStore) count() int { s.mu.Lock(); defer s.mu.Unlock(); return len(s.records) }
func (s *FixedStore) setCommitFailureForTest(value bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failCommit = value
}
