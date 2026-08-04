//go:build !darwin

package installationowner

import "context"

// TrustedStateRoot is unavailable off Darwin. The type remains present so
// callers receive a fixed unsupported-platform result instead of a build-time
// fallback to weaker locking semantics.
type TrustedStateRoot struct{}

// OpenDarwinTrustedStateRoot refuses on non-Darwin platforms.
func OpenDarwinTrustedStateRoot(
	context.Context,
	string,
	OwnerLockEnrollment,
) (*TrustedStateRoot, error) {
	return nil, ErrUnsupportedPlatform
}

// Close refuses on non-Darwin platforms.
func (*TrustedStateRoot) Close() error { return ErrUnsupportedPlatform }

// AcquireDarwinInstallationOwner refuses on non-Darwin platforms. No weaker
// lock or process-local fallback is selected.
func AcquireDarwinInstallationOwner(
	context.Context,
	*TrustedStateRoot,
	OwnerLockEnrollment,
) (InstallationOwner, error) {
	return nil, ErrUnsupportedPlatform
}
