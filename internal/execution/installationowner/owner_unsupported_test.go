//go:build !darwin

package installationowner

import (
	"context"
	"errors"
	"testing"
)

func TestUnsupportedPlatformHasNoFallback(t *testing.T) {
	root, err := OpenDarwinTrustedStateRoot(context.Background(), "ignored", OwnerLockEnrollment{})
	if root != nil || !errors.Is(err, ErrUnsupportedPlatform) {
		t.Fatalf("unsupported root result = %T, %v", root, err)
	}
	owner, err := AcquireDarwinInstallationOwner(context.Background(), nil, OwnerLockEnrollment{})
	if owner != nil || !errors.Is(err, ErrUnsupportedPlatform) {
		t.Fatalf("unsupported owner result = %T, %v", owner, err)
	}
}
