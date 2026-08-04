package archivestate

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"capsule.local/capsule/internal/execution/lifecyclestate"
)

var updateArchiveFormatKnownAnswers = flag.Bool(
	"update-archive-format-known-answers",
	false,
	"regenerate the passive F2 format-correction known-answer fixture",
)

type archiveFormatKnownAnswers struct {
	EmptyRetainedRegistrationIndex string        `json:"emptyRetainedRegistrationIndex"`
	EmptyRetainedCombinedIndex     string        `json:"emptyRetainedCombinedIndex"`
	EmptySegmentRegistrationIndex  string        `json:"emptySegmentRegistrationIndex"`
	EmptySegmentCombinedIndex      string        `json:"emptySegmentCombinedIndex"`
	EmptyDescriptorSet             string        `json:"emptyDescriptorSet"`
	EmptyVisibleV1EffectSeed       string        `json:"emptyVisibleV1EffectSeed"`
	NonemptyVisibleV1EffectSeed    string        `json:"nonemptyVisibleV1EffectSeed"`
	ArchivedRetainedCombinedIndex  string        `json:"archivedRetainedCombinedIndex"`
	MigrationHotCombinedIndex      string        `json:"migrationHotCombinedIndex"`
	MigrationGenesisCheckpoint     string        `json:"migrationGenesisCheckpoint"`
	MissingLifecycleCombinedIndex  string        `json:"missingLifecycleCombinedIndex"`
	MissingLifecycleGenesis        string        `json:"missingLifecycleGenesis"`
	ActivationCheckpoint           string        `json:"activationCheckpoint"`
	OneCohortSegment               string        `json:"oneCohortSegment"`
	MultiCohortSegment             string        `json:"multiCohortSegment"`
	MigrationVisibleV1SeedCount    uint64        `json:"migrationVisibleV1SeedCount"`
	MigrationHotCounts             ArchiveCounts `json:"migrationHotCounts"`
	MigrationArchivedCounts        ArchiveCounts `json:"migrationArchivedCounts"`
	MigrationTotalCounts           ArchiveCounts `json:"migrationTotalCounts"`
	MissingLifecycleHotCounts      ArchiveCounts `json:"missingLifecycleHotCounts"`
}

func TestArchiveFormatCorrectionKnownAnswers(t *testing.T) {
	fixturePath := filepath.Join("testdata", "format-correction-known-answers.json")
	got := calculateArchiveFormatKnownAnswers(t)
	if *updateArchiveFormatKnownAnswers {
		encoded, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		encoded = append(encoded, '\n')
		if err := os.WriteFile(fixturePath, encoded, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	expectedBytes, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}
	var expected archiveFormatKnownAnswers
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("archive format known answers changed\n got: %#v\nwant: %#v\nregenerate with go test ./internal/execution/archivestate -run TestArchiveFormatCorrectionKnownAnswers -update-archive-format-known-answers", got, expected)
	}
}

func calculateArchiveFormatKnownAnswers(t *testing.T) archiveFormatKnownAnswers {
	t.Helper()
	emptyRetained := EmptyRetainedIndexes()
	emptySegment := EmptySegmentDerivedIndexes()
	descriptorDigest, err := DigestArchiveDescriptorSet([]ArchiveDescriptor{})
	if err != nil {
		t.Fatal(err)
	}
	emptySeed, err := DigestVisibleV1EffectSeed([]lifecyclestate.EffectID{})
	if err != nil {
		t.Fatal(err)
	}
	hotIndexes := completeHotRetainedIndexes(t, true)
	seed := []lifecyclestate.EffectID{effectID(1)}
	genesis, err := NewMigrationGenesisCheckpoint(MigrationGenesisCheckpointView{
		StoreFormatVersion: SupervisorStoreFormatV2, MigrationSourceVersion: MigrationSourceFormatV1,
		ResultSnapshotGeneration: 1, ArchiveGeneration: 1, DescriptorSetDigest: descriptorDigest,
		Indexes: hotIndexes, HotSetDigests: hotSetDigests(1), VisibleV1EffectSeed: seed,
		HotCounts: hotIndexes.counts(), InstallationID: installationID(1), SupervisorID: supervisorID(2),
		EpochSequence: 3, EpochDigest: epochDigest(4), DurableTimeHighWater: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	missingLifecycleIndexes := completeHotRetainedIndexes(t, false).View()
	missingLifecycleIndexes.Attempts[0].Lifecycle = NoAttemptLifecycle()
	missingLifecycleIndexes.Effects = []EffectIndexEntry{}
	missingLifecycleIndexes.Instances = []InstanceIndexEntry{}
	missingLifecycle, err := NewArchiveIndexes(missingLifecycleIndexes)
	if err != nil {
		t.Fatal(err)
	}
	missingLifecycleGenesis, err := NewMigrationGenesisCheckpoint(MigrationGenesisCheckpointView{
		StoreFormatVersion: SupervisorStoreFormatV2, MigrationSourceVersion: MigrationSourceFormatV1,
		ResultSnapshotGeneration: 1, ArchiveGeneration: 1, DescriptorSetDigest: descriptorDigest,
		Indexes: missingLifecycle, HotSetDigests: hotSetDigests(1), VisibleV1EffectSeed: []lifecyclestate.EffectID{},
		HotCounts: missingLifecycle.counts(), InstallationID: installationID(1), SupervisorID: supervisorID(2),
		EpochSequence: 3, EpochDigest: epochDigest(4), DurableTimeHighWater: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	activation, err := NewArchiveCheckpoint(ArchiveCheckpointView{
		PreviousCheckpoint: genesis.Reference(), NewSegmentDigest: segmentDigest(1),
		DescriptorSetDigest: descriptorSetDigest(1), ArchiveIndexDigest: combinedIndexDigest(1),
		HotSetDigests: hotSetDigests(1), SourceSnapshotGeneration: 1, ResultSnapshotGeneration: 2,
		ArchiveGeneration: 2, InstallationID: installationID(1), SupervisorID: supervisorID(2),
		EpochSequence: 3, EpochDigest: epochDigest(4), DurableTimeHighWater: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if activation.Digest() == genesis.Digest() || activation.Reference().Kind == genesis.Reference().Kind {
		t.Fatal("migration genesis and activation checkpoint domains are not distinct")
	}
	zeroCounts := ArchiveCounts{}
	one := mustCohortProjection(t, closedCohort(t, 1, 90, 1))
	multi := mustCohortProjection(t, closedCohort(t, 2, 90, 2))
	return archiveFormatKnownAnswers{
		EmptyRetainedRegistrationIndex: hex32(emptyRetained.Digests().Registrations),
		EmptyRetainedCombinedIndex:     hex32(emptyRetained.CombinedDigest()),
		EmptySegmentRegistrationIndex:  hex32(emptySegment.Digests().Registrations),
		EmptySegmentCombinedIndex:      hex32(emptySegment.CombinedDigest()),
		EmptyDescriptorSet:             hex32(descriptorDigest),
		EmptyVisibleV1EffectSeed:       hex32(emptySeed),
		NonemptyVisibleV1EffectSeed:    hex32(genesis.SeedDigest()),
		ArchivedRetainedCombinedIndex:  hex32(completeArchiveIndexes(t).CombinedDigest()),
		MigrationHotCombinedIndex:      hex32(hotIndexes.CombinedDigest()),
		MigrationGenesisCheckpoint:     hex32(genesis.Digest()),
		MissingLifecycleCombinedIndex:  hex32(missingLifecycle.CombinedDigest()),
		MissingLifecycleGenesis:        hex32(missingLifecycleGenesis.Digest()),
		ActivationCheckpoint:           hex32(activation.Digest()),
		OneCohortSegment:               hex32(mustSegment(t, []CohortProjection{one}, 1_024).Digest()),
		MultiCohortSegment:             hex32(mustSegment(t, []CohortProjection{one, multi}, 2_048).Digest()),
		MigrationVisibleV1SeedCount:    1,
		MigrationHotCounts:             hotIndexes.counts(),
		MigrationArchivedCounts:        zeroCounts,
		MigrationTotalCounts:           hotIndexes.counts(),
		MissingLifecycleHotCounts:      missingLifecycle.counts(),
	}
}

func hex32[T ~[32]byte](value T) string { return hex.EncodeToString(value[:]) }
