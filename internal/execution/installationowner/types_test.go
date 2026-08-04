package installationowner

import (
	"strings"
	"testing"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

func TestOwnerLockEnrollmentIsClosedAndNonzero(t *testing.T) {
	ordinary := OwnerLockEnrollmentView{
		InstallationID: repeated16[v0candidate.InstallationID](0x11),
		SupervisorID:   repeated16[v0candidate.SupervisorID](0x55),
		ExpectedUID:    501, StateRootDevice: 1, StateRootInode: 2,
		LockEntryName: "supervisor.owner", LockDevice: 1, LockInode: 3,
	}
	if _, err := NewOwnerLockEnrollment(ordinary); err != nil {
		t.Fatalf("ordinary enrollment: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*OwnerLockEnrollmentView)
	}{
		{name: "installation", mutate: func(view *OwnerLockEnrollmentView) { view.InstallationID = v0candidate.InstallationID{} }},
		{name: "supervisor", mutate: func(view *OwnerLockEnrollmentView) { view.SupervisorID = v0candidate.SupervisorID{} }},
		{name: "root uid", mutate: func(view *OwnerLockEnrollmentView) { view.ExpectedUID = 0 }},
		{name: "root device", mutate: func(view *OwnerLockEnrollmentView) { view.StateRootDevice = 0 }},
		{name: "root inode", mutate: func(view *OwnerLockEnrollmentView) { view.StateRootInode = 0 }},
		{name: "lock device", mutate: func(view *OwnerLockEnrollmentView) { view.LockDevice = 0 }},
		{name: "lock inode", mutate: func(view *OwnerLockEnrollmentView) { view.LockInode = 0 }},
		{name: "empty name", mutate: func(view *OwnerLockEnrollmentView) { view.LockEntryName = "" }},
		{name: "dot", mutate: func(view *OwnerLockEnrollmentView) { view.LockEntryName = "." }},
		{name: "dot dot", mutate: func(view *OwnerLockEnrollmentView) { view.LockEntryName = ".." }},
		{name: "slash", mutate: func(view *OwnerLockEnrollmentView) { view.LockEntryName = "other/owner" }},
		{name: "backslash", mutate: func(view *OwnerLockEnrollmentView) { view.LockEntryName = `other\owner` }},
		{name: "non ascii", mutate: func(view *OwnerLockEnrollmentView) { view.LockEntryName = "ownér" }},
		{name: "control", mutate: func(view *OwnerLockEnrollmentView) { view.LockEntryName = "owner\n" }},
		{name: "too long", mutate: func(view *OwnerLockEnrollmentView) { view.LockEntryName = strings.Repeat("x", 256) }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := ordinary
			test.mutate(&candidate)
			if _, err := NewOwnerLockEnrollment(candidate); err != ErrInvalidEnrollment {
				t.Fatalf("invalid enrollment error = %v", err)
			}
		})
	}
}

func repeated16[T ~[16]byte](value byte) (result T) {
	for index := range result {
		result[index] = value
	}
	return result
}
