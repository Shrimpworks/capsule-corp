package archivestate

import "capsule.local/capsule/internal/protocol/v0candidate"

const (
	SupervisorStoreFormatV2        = uint64(2)
	SupervisorArchiveFormatV0      = uint64(0)
	MaxSupervisorStateV2Bytes      = int64(64 << 20)
	MaxSupervisorArchiveBytes      = int64(64 << 20)
	MaxArchiveCohortsPerSegment    = 256
	MaxArchiveApprovalsPerSegment  = 4_096
	MaxArchiveAttemptsPerSegment   = 4_096
	MaxArchiveLifecyclesPerSegment = 4_096
	MaxReferencedArchiveSegments   = 64
	MaxArchiveIndexEntries         = 262_144

	MigrationSourceFormatV1 = uint64(1)
)

const EffectTombstoneCoverageVisibleV1AndV2 = "visible-v1-seed-plus-all-v2-issued"

// ArchiveGeneration is a positive UInt53 generation in the archive domain.
type ArchiveGeneration v0candidate.PositiveUInt53

// ArchiveOrdinal is a positive UInt53 immutable segment or cohort ordinal.
type ArchiveOrdinal v0candidate.PositiveUInt53

// ArchiveLimits bounds one pure selection. Values above the fixed oracle caps
// are refused rather than silently narrowed.
type ArchiveLimits struct {
	MaxCohorts uint16
	MaxBytes   uint64
}

func (limits ArchiveLimits) validate() error {
	if limits.MaxCohorts == 0 || limits.MaxBytes == 0 {
		return classified(ClassificationSchema, "archive-limits-zero")
	}
	if int(limits.MaxCohorts) > MaxArchiveCohortsPerSegment ||
		limits.MaxBytes > uint64(MaxSupervisorArchiveBytes) {
		return classified(ClassificationCapacity, "archive-limits")
	}
	return nil
}

// NewArchiveLimits validates exact caller-selected limits without clamping.
func NewArchiveLimits(maxCohorts uint16, maxBytes uint64) (ArchiveLimits, error) {
	limits := ArchiveLimits{MaxCohorts: maxCohorts, MaxBytes: maxBytes}
	if err := limits.validate(); err != nil {
		return ArchiveLimits{}, err
	}
	return limits, nil
}

func validPositive(value uint64) bool {
	return value > 0 && value <= v0candidate.MaxSafeInteger
}
