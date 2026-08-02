package registrationstate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

const conformanceRoot = "../../../schemas/conformance/v0"

type conformanceManifest struct {
	Cases []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Fixture struct {
			Path string `json:"path"`
		} `json:"fixture"`
		Context struct {
			Kind      string `json:"kind"`
			Operation struct {
				Path string `json:"path"`
			} `json:"operation"`
			Before struct {
				Path string `json:"path"`
			} `json:"before"`
		} `json:"context"`
		Expected struct {
			Decision              string  `json:"decision"`
			Classification        *string `json:"classification"`
			AuthorityStateChanged bool    `json:"authorityStateChanged"`
			StateDelta            struct {
				After struct {
					Path string `json:"path"`
				} `json:"after"`
			} `json:"stateDelta"`
		} `json:"expected"`
		Implementations map[string]string `json:"implementations"`
	} `json:"cases"`
}

type operationFixture struct {
	Method string `json:"method"`
	Caller struct {
		Authenticated bool   `json:"authenticated"`
		Role          string `json:"role"`
		Purpose       string `json:"purpose"`
	} `json:"caller"`
	TrustedClockObservationUnixSeconds uint64  `json:"trustedClockObservationUnixSeconds"`
	ConcurrentHighWaterUnixSeconds     *uint64 `json:"concurrentHighWaterUnixSeconds"`
	Identifier                         *struct {
		Kind     string `json:"kind"`
		ValueHex string `json:"valueHex"`
		Failure  string `json:"failure"`
	} `json:"identifier"`
	Mutation string  `json:"mutation"`
	Fault    *string `json:"fault"`
}

type stateFixture struct {
	InstallationIDHex        string  `json:"installationIdHex"`
	SupervisorIDHex          string  `json:"supervisorIdHex"`
	EpochSequence            uint64  `json:"epochSequence"`
	EpochDigestHex           string  `json:"epochDigestHex"`
	TrustPhase               string  `json:"trustPhase"`
	TrustReason              *string `json:"trustReason"`
	Quarantined              bool    `json:"quarantined"`
	TimeHighWaterUnixSeconds uint64  `json:"timeHighWaterUnixSeconds"`
	LastRegistrationSequence uint64  `json:"lastRegistrationSequence"`
	RegistrationPopulation   struct {
		StoredCount    int    `json:"storedCount"`
		UnexpiredCount int    `json:"unexpiredCount"`
		SetDigest      string `json:"setDigest"`
	} `json:"registrationPopulation"`
	MaterializedRecords []storedRecordFixture `json:"materializedRecords"`
}

type storedRecordFixture struct {
	WireRegistrationHex     string `json:"wireRegistrationHex"`
	ExactPlanHex            string `json:"exactPlanHex"`
	RecomputedPlanDigestHex string `json:"recomputedPlanDigestHex"`
	RegisteredAtUnixSeconds uint64 `json:"registeredAtUnixSeconds"`
	StorageFormatVersion    uint64 `json:"storageFormatVersion"`
	RetentionState          string `json:"retentionState"`
}

func TestRegistrationStateConformanceMatrix(t *testing.T) {
	manifest := readJSON[conformanceManifest](t, filepath.Join(conformanceRoot, "manifest.json"))
	stateCases := 0
	accepts := 0
	rejects := 0
	for _, testCase := range manifest.Cases {
		if testCase.Context.Kind != "registration-state" {
			continue
		}
		stateCases++
		if testCase.Implementations["go"] != "verified" {
			t.Fatalf("registration-state case %s is not marked as verified for Go", testCase.ID)
		}
		if testCase.Expected.Decision == "accept" {
			accepts++
		} else {
			rejects++
		}
		testCase := testCase
		t.Run(testCase.ID, func(t *testing.T) {
			operation := readJSON[operationFixture](
				t, filepath.Join(conformanceRoot, testCase.Context.Operation.Path),
			)
			beforeFixture := readJSON[stateFixture](
				t, filepath.Join(conformanceRoot, testCase.Context.Before.Path),
			)
			afterFixture := readJSON[stateFixture](
				t, filepath.Join(conformanceRoot, testCase.Expected.StateDelta.After.Path),
			)
			planOrRegistration := readBytes(
				t, filepath.Join(conformanceRoot, testCase.Fixture.Path),
			)
			state := buildStateFromFixture(t, beforeFixture)
			storePath := filepath.Join(t.TempDir(), "registration-state.json")
			if err := persistState(storePath, state); err != nil {
				t.Fatalf("persist fixture pre-state: %v", err)
			}
			store, err := NewFixedFileStore(storePath, InitialState{})
			if err != nil {
				t.Fatalf("open fixture store: %v", err)
			}
			if operation.Fault != nil {
				switch *operation.Fault {
				case string(FaultTimeHighWaterWrite):
					store.InjectFailure(FaultTimeHighWaterWrite, errors.New("injected high-water failure"))
				case string(FaultRegistrationCommitAbort):
					store.InjectFailure(FaultRegistrationCommitAbort, errors.New("injected confirmed abort"))
				default:
					t.Fatalf("unimplemented fixture fault %q", *operation.Fault)
				}
			}

			clock := fixedTrustedClock{observation: operation.TrustedClockObservationUnixSeconds}
			identifier := identifierSourceFromFixture(t, operation.Identifier)
			var requestBytes []byte
			component, err := New(Options{
				Store:       store,
				Clock:       clock,
				Identifiers: identifier,
				Checkpoint: func(
					ctx context.Context,
					checkpoint Checkpoint,
					decoded *v0candidate.DecodedExecutionPlan,
				) {
					if checkpoint == CheckpointPlanDecoded {
						switch operation.Mutation {
						case "caller-buffer-after-submission":
							requestBytes[0] ^= 0xff
						case "validator-private-copy":
							privateCopy := decoded.AuthoritativeBytes()
							privateCopy[0] ^= 0xff
						}
					}
					if checkpoint == CheckpointTimeHighWaterDone &&
						operation.ConcurrentHighWaterUnixSeconds != nil {
						if err := store.persistTimeHighWater(
							ctx,
							v0candidate.UInt53(*operation.ConcurrentHighWaterUnixSeconds),
						); err != nil {
							t.Fatalf("persist concurrent high water: %v", err)
						}
					}
				},
			})
			if err != nil {
				t.Fatalf("new registration component: %v", err)
			}

			beforeSnapshot, err := snapshotComponent(context.Background(), component)
			if err != nil {
				t.Fatalf("snapshot before: %v", err)
			}
			var operationErr error
			var issued PlanRegistration
			switch operation.Method {
			case "register-plan":
				requestBytes = bytes.Clone(planOrRegistration)
				issued, operationErr = component.RegisterPlan(
					context.Background(),
					AuthenticatedCallContext{
						Authenticated: operation.Caller.Authenticated,
						Role:          CallerRole(operation.Caller.Role),
						Purpose:       operation.Caller.Purpose,
					},
					requestBytes,
					bindingsForPlanFixture(testCase.Fixture.Path),
				)
			case "use-registration":
				registrationID := decodeUseFixtureID(t, testCase.Fixture.Path, planOrRegistration)
				_, operationErr = component.ResolveUsable(context.Background(), registrationID)
			default:
				t.Fatalf("unimplemented fixture method %q", operation.Method)
			}

			assertDecision(t, testCase.Expected.Decision, testCase.Expected.Classification, operationErr)
			afterSnapshot, err := snapshotComponent(context.Background(), component)
			if err != nil {
				t.Fatalf("snapshot after: %v", err)
			}
			assertSnapshotMatches(t, afterSnapshot, afterFixture)
			authorityChanged := beforeSnapshot.StoredRegistrationCount != afterSnapshot.StoredRegistrationCount ||
				beforeSnapshot.LastRegistrationSequence != afterSnapshot.LastRegistrationSequence
			if authorityChanged != testCase.Expected.AuthorityStateChanged {
				t.Fatalf(
					"authorityStateChanged = %t, want %t",
					authorityChanged,
					testCase.Expected.AuthorityStateChanged,
				)
			}
			if operationErr == nil && operation.Method == "register-plan" {
				assertIssuedMatchesMaterializedRecord(t, issued, afterFixture.MaterializedRecords)
			}
		})
	}
	if stateCases != 40 || accepts != 18 || rejects != 22 {
		t.Fatalf("registration matrix = %d cases (%d accept, %d reject), want 40 (18, 22)", stateCases, accepts, rejects)
	}
}

func TestFixedStoreRestartAndDefensiveCopies(t *testing.T) {
	plan := readBytes(t, filepath.Join(conformanceRoot, "execution-plan/ordinary.cbor"))
	path := filepath.Join(t.TempDir(), "registration-state.json")
	store, err := NewFixedFileStore(path, ordinaryInitialState())
	if err != nil {
		t.Fatalf("new fixed store: %v", err)
	}
	component := mustComponent(t, store, fixedTrustedClock{observation: 1_785_456_000}, &sequenceIdentifierSource{next: 1})
	issued, err := component.RegisterPlan(
		context.Background(),
		AuthenticatedCallContext{Authenticated: true, Role: CallerDaemon, Purpose: RegisterPlanPurpose},
		plan,
		ordinaryPlanBindings(),
	)
	if err != nil {
		t.Fatalf("register plan: %v", err)
	}
	wireCopy := issued.Bytes()
	wireCopy[0] ^= 0xff
	snapshot, err := snapshotComponent(context.Background(), component)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	snapshot.Records[0].ExactPlanBytes[0] ^= 0xff
	snapshot.Records[0].WireRegistrationBytes[0] ^= 0xff

	reopened, err := NewFixedFileStore(path, InitialState{})
	if err != nil {
		t.Fatalf("reopen fixed store: %v", err)
	}
	restarted := mustComponent(t, reopened, fixedTrustedClock{observation: 1_785_456_001}, &sequenceIdentifierSource{next: 2})
	record, err := restarted.ResolveUsable(context.Background(), issued.View().RegistrationID)
	if err != nil {
		t.Fatalf("resolve after restart: %v", err)
	}
	if !bytes.Equal(record.ExactPlanBytes, plan) || !bytes.Equal(record.WireRegistrationBytes, issued.Bytes()) {
		t.Fatal("caller mutation changed durable registration bytes")
	}
	record.ExactPlanBytes[0] ^= 0xff
	again, err := restarted.ResolveUsable(context.Background(), issued.View().RegistrationID)
	if err != nil {
		t.Fatalf("resolve defensive copy again: %v", err)
	}
	if !bytes.Equal(again.ExactPlanBytes, plan) {
		t.Fatal("resolved record did not return a defensive copy")
	}
}

func TestFixedStoreRestartRefusesStoredPlanDigestMismatch(t *testing.T) {
	plan := readBytes(t, filepath.Join(conformanceRoot, "execution-plan/ordinary.cbor"))
	path := filepath.Join(t.TempDir(), "registration-state.json")
	store, err := NewFixedFileStore(path, ordinaryInitialState())
	if err != nil {
		t.Fatalf("new fixed store: %v", err)
	}
	component := mustComponent(t, store, fixedTrustedClock{observation: 1_785_456_000}, &sequenceIdentifierSource{next: 1})
	if _, err := component.RegisterPlan(
		context.Background(),
		AuthenticatedCallContext{Authenticated: true, Role: CallerDaemon, Purpose: RegisterPlanPurpose},
		plan,
		ordinaryPlanBindings(),
	); err != nil {
		t.Fatalf("register plan: %v", err)
	}
	state, err := store.snapshot(context.Background())
	if err != nil {
		t.Fatalf("snapshot registered state: %v", err)
	}
	state.Registrations[0].Record.RecomputedPlanDigest[0] ^= 0xff
	corrupt, err := json.Marshal(diskEnvelope{StoreFormatVersion: 0, State: state})
	if err != nil {
		t.Fatalf("marshal digest-mismatched store: %v", err)
	}
	corrupt = append(corrupt, '\n')
	if err := os.WriteFile(path, corrupt, 0o600); err != nil {
		t.Fatalf("write digest-mismatched store: %v", err)
	}
	if _, err := NewFixedFileStore(path, InitialState{}); err == nil ||
		!strings.Contains(err.Error(), "stored plan digest mismatch") {
		t.Fatalf("restart error = %v, want stored plan digest mismatch", err)
	}
}

func TestTrustedClockObservationFailureChangesNoState(t *testing.T) {
	plan := readBytes(t, filepath.Join(conformanceRoot, "execution-plan/ordinary.cbor"))
	for _, test := range []struct {
		name  string
		clock fixedTrustedClock
	}{
		{name: "source error", clock: fixedTrustedClock{err: errors.New("clock unavailable")}},
		{name: "outside UInt53", clock: fixedTrustedClock{observation: v0candidate.MaxSafeInteger + 1}},
	} {
		t.Run(test.name, func(t *testing.T) {
			store, err := NewFixedFileStore(
				filepath.Join(t.TempDir(), "state.json"),
				ordinaryInitialState(),
			)
			if err != nil {
				t.Fatalf("new fixed store: %v", err)
			}
			component := mustComponent(t, store, test.clock, &sequenceIdentifierSource{next: 1})
			before, err := snapshotComponent(context.Background(), component)
			if err != nil {
				t.Fatalf("snapshot before: %v", err)
			}
			_, err = component.RegisterPlan(
				context.Background(),
				AuthenticatedCallContext{Authenticated: true, Role: CallerDaemon, Purpose: RegisterPlanPurpose},
				plan,
				ordinaryPlanBindings(),
			)
			classification, ok := ErrorClassification(err)
			if !ok || classification != ClassificationLocalFailure {
				t.Fatalf("clock failure classification = %q (%t), want LOCAL_FAILURE: %v", classification, ok, err)
			}
			after, snapshotErr := snapshotComponent(context.Background(), component)
			if snapshotErr != nil {
				t.Fatalf("snapshot after: %v", snapshotErr)
			}
			if !reflect.DeepEqual(before, after) {
				t.Fatal("trusted-clock observation failure changed registration state")
			}
		})
	}
}

func TestConcurrentRegistrationCommitsCompleteUniqueSequences(t *testing.T) {
	const workers = 24
	plan := readBytes(t, filepath.Join(conformanceRoot, "execution-plan/ordinary.cbor"))
	store, err := NewFixedFileStore(filepath.Join(t.TempDir(), "state.json"), ordinaryInitialState())
	if err != nil {
		t.Fatalf("new fixed store: %v", err)
	}
	component := mustComponent(
		t,
		store,
		fixedTrustedClock{observation: 1_785_456_000},
		&sequenceIdentifierSource{next: 1},
	)
	results := make(chan PlanRegistration, workers)
	errorsSeen := make(chan error, workers)
	var group sync.WaitGroup
	for range workers {
		group.Add(1)
		go func() {
			defer group.Done()
			registration, err := component.RegisterPlan(
				context.Background(),
				AuthenticatedCallContext{Authenticated: true, Role: CallerDaemon, Purpose: RegisterPlanPurpose},
				plan,
				ordinaryPlanBindings(),
			)
			if err != nil {
				errorsSeen <- err
				return
			}
			results <- registration
		}()
	}
	group.Wait()
	close(results)
	close(errorsSeen)
	for err := range errorsSeen {
		t.Fatalf("concurrent registration: %v", err)
	}
	sequences := make([]int, 0, workers)
	identifiers := make(map[v0candidate.RegistrationID]struct{}, workers)
	for registration := range results {
		view := registration.View()
		sequences = append(sequences, int(view.RegistrationSequence))
		identifiers[view.RegistrationID] = struct{}{}
	}
	sort.Ints(sequences)
	if len(sequences) != workers || len(identifiers) != workers {
		t.Fatalf("got %d sequences and %d IDs, want %d", len(sequences), len(identifiers), workers)
	}
	for index, sequence := range sequences {
		if sequence != index+1 {
			t.Fatalf("sequence[%d] = %d, want %d", index, sequence, index+1)
		}
	}
	snapshot, err := snapshotComponent(context.Background(), component)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if snapshot.StoredRegistrationCount != workers || len(snapshot.Records) != workers {
		t.Fatalf("store exposed partial commits: count=%d records=%d", snapshot.StoredRegistrationCount, len(snapshot.Records))
	}
}

func buildStateFromFixture(t *testing.T, fixture stateFixture) installationState {
	t.Helper()
	state := installationState{
		InitialState: InitialState{
			InstallationID:           mustHex16[v0candidate.InstallationID](t, fixture.InstallationIDHex),
			SupervisorID:             mustHex16[v0candidate.SupervisorID](t, fixture.SupervisorIDHex),
			EpochSequence:            v0candidate.UInt53(fixture.EpochSequence),
			EpochDigest:              mustHex32[v0candidate.TrustEpochDigest](t, fixture.EpochDigestHex),
			TrustPhase:               TrustPhase(fixture.TrustPhase),
			TrustReason:              optionalString(fixture.TrustReason),
			Quarantined:              fixture.Quarantined,
			TimeHighWaterUnixSeconds: v0candidate.UInt53(fixture.TimeHighWaterUnixSeconds),
			LastRegistrationSequence: v0candidate.UInt53(fixture.LastRegistrationSequence),
		},
		RegistrationSetDigest: mustHex32[[32]byte](t, fixture.RegistrationPopulation.SetDigest),
		Registrations:         make([]registrationEntry, 0, fixture.RegistrationPopulation.StoredCount),
	}
	for _, materialized := range fixture.MaterializedRecords {
		state.Registrations = append(state.Registrations, materializedEntry(t, materialized))
	}
	materializedUnexpired := countUnexpired(state.Registrations, state.TimeHighWaterUnixSeconds)
	missingCount := fixture.RegistrationPopulation.StoredCount - len(state.Registrations)
	missingUnexpired := fixture.RegistrationPopulation.UnexpiredCount - materializedUnexpired
	if missingCount < 0 || missingUnexpired < 0 || missingUnexpired > missingCount {
		t.Fatal("invalid compact registration population fixture")
	}
	for index := 0; index < missingCount; index++ {
		sequence := uint64(index + 1)
		if fixture.LastRegistrationSequence == v0candidate.MaxSafeInteger && missingCount == 1 {
			sequence = v0candidate.MaxSafeInteger
		}
		state.Registrations = append(
			state.Registrations,
			syntheticEntry(t, state.InitialState, sequence, index < missingUnexpired),
		)
	}
	if got := countUnexpired(state.Registrations, state.TimeHighWaterUnixSeconds); got != fixture.RegistrationPopulation.UnexpiredCount {
		t.Fatalf("built unexpired count = %d, want %d", got, fixture.RegistrationPopulation.UnexpiredCount)
	}
	if err := validateState(state); err != nil {
		t.Fatalf("built fixture state is invalid: %v", err)
	}
	return state
}

func materializedEntry(t *testing.T, fixture storedRecordFixture) registrationEntry {
	t.Helper()
	record := storedRecordFromFixture(t, fixture)
	decoded, err := v0candidate.DecodePlanRegistration(
		record.WireRegistrationBytes,
		ordinaryRegistrationBindings(0x77),
	)
	if err != nil {
		t.Fatalf("decode materialized registration: %v", err)
	}
	view := decoded.View()
	return registrationEntry{
		Index: registrationIndex{
			RegistrationID:       view.RegistrationID,
			RegistrationSequence: view.RegistrationSequence,
			InstallationID:       view.InstallationID,
			EpochSequence:        view.EpochSequence,
			EpochDigest:          view.EpochDigest,
			SupervisorID:         view.SupervisorID,
			ExpiresAt:            view.ExpiresAt,
		},
		PlanBindings: ordinaryPlanBindings(),
		Record:       record,
	}
}

func syntheticEntry(
	t *testing.T,
	state InitialState,
	sequence uint64,
	unexpired bool,
) registrationEntry {
	t.Helper()
	expiresAt := uint64(state.TimeHighWaterUnixSeconds)
	registeredAt := uint64(state.TimeHighWaterUnixSeconds)
	if unexpired {
		expiresAt++
	} else if registeredAt > 0 {
		registeredAt--
	}
	plan := bytes.Clone(readBytes(t, filepath.Join(conformanceRoot, "execution-plan/ordinary.cbor")))
	if len(plan) < 4 {
		t.Fatal("ordinary plan fixture is truncated")
	}
	binary.BigEndian.PutUint32(plan[len(plan)-4:], uint32(expiresAt))
	digest := v0candidate.ExecutionPlanDigest(sha256.Sum256(plan))
	identifier := syntheticID(sequence)
	bindings := ordinaryPlanBindings()
	view := v0candidate.PlanRegistration{
		ObjectType:           v0candidate.PlanRegistrationObjectType,
		ObjectVersion:        v0candidate.CandidateObjectVersion,
		RegistrationID:       identifier,
		RegistrationSequence: v0candidate.PositiveUInt53(sequence),
		PlanDigest:           digest,
		InstallationID:       bindings.InstallationID,
		EpochSequence:        7,
		EpochDigest:          bindings.EpochDigest,
		SupervisorID:         state.SupervisorID,
		ExpiresAt:            v0candidate.UInt53(expiresAt),
	}
	return registrationEntry{
		Index: registrationIndex{
			RegistrationID:       identifier,
			RegistrationSequence: v0candidate.PositiveUInt53(sequence),
			InstallationID:       bindings.InstallationID,
			EpochSequence:        7,
			EpochDigest:          bindings.EpochDigest,
			SupervisorID:         state.SupervisorID,
			ExpiresAt:            v0candidate.UInt53(expiresAt),
		},
		PlanBindings: bindings,
		Record: StoredRegistration{
			WireRegistrationBytes:   encodePlanRegistration(view),
			ExactPlanBytes:          plan,
			RecomputedPlanDigest:    digest,
			RegisteredAtUnixSeconds: v0candidate.UInt53(registeredAt),
			StorageFormatVersion:    StorageFormatVersion,
			RetentionState:          RetentionStateRetained,
		},
	}
}

func syntheticID(sequence uint64) v0candidate.RegistrationID {
	var identifier v0candidate.RegistrationID
	identifier[0] = 0xd0
	binary.BigEndian.PutUint64(identifier[8:], sequence)
	return identifier
}

func assertSnapshotMatches(t *testing.T, actual stateSnapshot, expected stateFixture) {
	t.Helper()
	if actual.InstallationID != mustHex16[v0candidate.InstallationID](t, expected.InstallationIDHex) ||
		actual.SupervisorID != mustHex16[v0candidate.SupervisorID](t, expected.SupervisorIDHex) ||
		uint64(actual.EpochSequence) != expected.EpochSequence ||
		actual.EpochDigest != mustHex32[v0candidate.TrustEpochDigest](t, expected.EpochDigestHex) ||
		string(actual.TrustPhase) != expected.TrustPhase ||
		actual.TrustReason != optionalString(expected.TrustReason) ||
		actual.Quarantined != expected.Quarantined ||
		uint64(actual.TimeHighWaterUnixSeconds) != expected.TimeHighWaterUnixSeconds ||
		uint64(actual.LastRegistrationSequence) != expected.LastRegistrationSequence ||
		actual.StoredRegistrationCount != expected.RegistrationPopulation.StoredCount ||
		actual.UnexpiredRegistrationCount != expected.RegistrationPopulation.UnexpiredCount ||
		actual.RegistrationSetDigest != mustHex32[[32]byte](t, expected.RegistrationPopulation.SetDigest) {
		t.Fatalf("post-state mismatch:\nactual: %#v\nexpected: %#v", actual, expected)
	}
	for _, fixtureRecord := range expected.MaterializedRecords {
		expectedRecord := storedRecordFromFixture(t, fixtureRecord)
		found := false
		for _, actualRecord := range actual.Records {
			if recordsEqual(actualRecord, expectedRecord) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected materialized record was not retained: %#v", fixtureRecord)
		}
	}
}

func assertIssuedMatchesMaterializedRecord(
	t *testing.T,
	issued PlanRegistration,
	records []storedRecordFixture,
) {
	t.Helper()
	for _, record := range records {
		if hex.EncodeToString(issued.Bytes()) == record.WireRegistrationHex {
			if bytes.Contains(issued.Bytes(), mustDecodeHex(t, record.ExactPlanHex)) {
				t.Fatal("wire PlanRegistration unexpectedly embeds exact plan bytes")
			}
			return
		}
	}
	t.Fatal("issued wire registration did not match an exact post-state record")
}

func assertDecision(t *testing.T, decision string, expected *string, err error) {
	t.Helper()
	if decision == "accept" {
		if err != nil {
			t.Fatalf("accepted fixture failed: %v", err)
		}
		return
	}
	if err == nil {
		t.Fatal("rejected fixture unexpectedly succeeded")
	}
	classification, ok := ErrorClassification(err)
	if !ok || expected == nil || string(classification) != *expected {
		t.Fatalf("classification = %q (%t), want %v; error: %v", classification, ok, expected, err)
	}
}

func bindingsForPlanFixture(path string) v0candidate.ExecutionPlanRoleBindings {
	bindings := ordinaryPlanBindings()
	switch path {
	case "execution-plan/epoch-eight.cbor":
		bindings.EpochDigest = repeated32[v0candidate.TrustEpochDigest](0x23)
	case "execution-plan/wrong-epoch-binding.cbor":
		bindings.EpochDigest = repeated32[v0candidate.TrustEpochDigest](0x21)
	}
	return bindings
}

func decodeUseFixtureID(t *testing.T, path string, wire []byte) v0candidate.RegistrationID {
	t.Helper()
	idByte := byte(0x77)
	sequence := v0candidate.PositiveUInt53(1)
	if strings.Contains(path, "plan-a-registration-b") {
		idByte = 0x78
		sequence = 2
	}
	bindings := ordinaryRegistrationBindings(idByte)
	decoded, err := v0candidate.DecodePlanRegistration(wire, bindings)
	if err != nil {
		t.Fatalf("decode use-registration fixture: %v", err)
	}
	if decoded.View().RegistrationSequence != sequence {
		t.Fatalf("use-registration sequence = %d, want %d", decoded.View().RegistrationSequence, sequence)
	}
	return decoded.View().RegistrationID
}

func identifierSourceFromFixture(
	t *testing.T,
	identifier *struct {
		Kind     string `json:"kind"`
		ValueHex string `json:"valueHex"`
		Failure  string `json:"failure"`
	},
) IdentifierSource {
	t.Helper()
	if identifier == nil {
		return identifierSourceFunc(func(context.Context) (v0candidate.RegistrationID, error) {
			return v0candidate.RegistrationID{}, errors.New("identifier source must not be called")
		})
	}
	switch identifier.Kind {
	case "failure":
		return identifierSourceFunc(func(context.Context) (v0candidate.RegistrationID, error) {
			return v0candidate.RegistrationID{}, errors.New("injected identifier source error")
		})
	case "value":
		value := mustHex16[v0candidate.RegistrationID](t, identifier.ValueHex)
		return identifierSourceFunc(func(context.Context) (v0candidate.RegistrationID, error) {
			return value, nil
		})
	default:
		t.Fatalf("unknown identifier fixture kind %q", identifier.Kind)
		return nil
	}
}

func storedRecordFromFixture(t *testing.T, fixture storedRecordFixture) StoredRegistration {
	t.Helper()
	return StoredRegistration{
		WireRegistrationBytes:   mustDecodeHex(t, fixture.WireRegistrationHex),
		ExactPlanBytes:          mustDecodeHex(t, fixture.ExactPlanHex),
		RecomputedPlanDigest:    mustHex32[v0candidate.ExecutionPlanDigest](t, fixture.RecomputedPlanDigestHex),
		RegisteredAtUnixSeconds: v0candidate.UInt53(fixture.RegisteredAtUnixSeconds),
		StorageFormatVersion:    fixture.StorageFormatVersion,
		RetentionState:          fixture.RetentionState,
	}
}

func ordinaryInitialState() InitialState {
	return InitialState{
		InstallationID:           repeated16[v0candidate.InstallationID](0x11),
		SupervisorID:             repeated16[v0candidate.SupervisorID](0x55),
		EpochSequence:            7,
		EpochDigest:              repeated32[v0candidate.TrustEpochDigest](0x22),
		TrustPhase:               TrustStable,
		TimeHighWaterUnixSeconds: 1_785_456_000,
	}
}

func ordinaryPlanBindings() v0candidate.ExecutionPlanRoleBindings {
	return v0candidate.ExecutionPlanRoleBindings{
		InstallationID:                  repeated16[v0candidate.InstallationID](0x11),
		EpochDigest:                     repeated32[v0candidate.TrustEpochDigest](0x22),
		SourceManifestDigest:            mustHex32WithoutTest[v0candidate.SourceManifestDigest]("e5e09b2435baedf897526a89c698c0b0531437a69472372ae426f62d801fc171"),
		InlineInputDigest:               mustHex32WithoutTest[v0candidate.InlineInputDigest]("bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72"),
		RuntimeBundleManifestDigest:     repeated32[v0candidate.RuntimeBundleManifestDigest](0x55),
		ProfileReviewAttestationDigests: []v0candidate.ProfileReviewAttestationDigest{repeated32[v0candidate.ProfileReviewAttestationDigest](0x66), repeated32[v0candidate.ProfileReviewAttestationDigest](0x67)},
		ProfileRegistryEntryDigest:      repeated32[v0candidate.ProfileRegistryEntryDigest](0x77),
		BackendValidationRecordDigest:   repeated32[v0candidate.BackendValidationRecordDigest](0x88),
		BackendConfigurationDigest:      repeated32[v0candidate.BackendConfigurationDigest](0x99),
		TrustSnapshotDigest:             repeated32[v0candidate.TrustSnapshotDigest](0xaa),
		PolicyDecisionDigest:            repeated32[v0candidate.PolicyDecisionDigest](0xbb),
	}
}

func ordinaryRegistrationBindings(id byte) v0candidate.PlanRegistrationRoleBindings {
	return v0candidate.PlanRegistrationRoleBindings{
		RegistrationID: repeated16[v0candidate.RegistrationID](id),
		PlanDigest:     mustHex32WithoutTest[v0candidate.ExecutionPlanDigest]("627f9524479000dab6f3cee1d70c0428c63285bcadbc2cb3c6e8018b2dea008c"),
		InstallationID: repeated16[v0candidate.InstallationID](0x11),
		EpochDigest:    repeated32[v0candidate.TrustEpochDigest](0x22),
		SupervisorID:   repeated16[v0candidate.SupervisorID](0x55),
	}
}

func mustComponent(
	t *testing.T,
	store StateStore,
	clock TrustedClock,
	identifiers IdentifierSource,
) *Component {
	t.Helper()
	component, err := New(Options{Store: store, Clock: clock, Identifiers: identifiers})
	if err != nil {
		t.Fatalf("new component: %v", err)
	}
	return component
}

func snapshotComponent(ctx context.Context, component *Component) (stateSnapshot, error) {
	state, err := component.store.snapshot(ctx)
	if err != nil {
		return stateSnapshot{}, err
	}
	return snapshotProjection(state), nil
}

type fixedTrustedClock struct {
	observation uint64
	err         error
}

func (clock fixedTrustedClock) ObserveUnixSeconds(context.Context) (uint64, error) {
	return clock.observation, clock.err
}

type identifierSourceFunc func(context.Context) (v0candidate.RegistrationID, error)

func (function identifierSourceFunc) NewRegistrationID(
	ctx context.Context,
) (v0candidate.RegistrationID, error) {
	return function(ctx)
}

type sequenceIdentifierSource struct {
	mu   sync.Mutex
	next uint64
}

func (source *sequenceIdentifierSource) NewRegistrationID(
	context.Context,
) (v0candidate.RegistrationID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	id := syntheticID(source.next)
	source.next++
	return id, nil
}

func readJSON[T any](t *testing.T, path string) T {
	t.Helper()
	data := readBytes(t, path)
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return value
}

func readBytes(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func mustDecodeHex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	if err != nil {
		t.Fatalf("decode hex: %v", err)
	}
	return decoded
}

func mustHex16[T ~[16]byte](t *testing.T, value string) T {
	t.Helper()
	decoded := mustDecodeHex(t, value)
	if len(decoded) != 16 {
		t.Fatalf("16-byte value has length %d", len(decoded))
	}
	var result T
	copy(result[:], decoded)
	return result
}

func mustHex32[T ~[32]byte](t *testing.T, value string) T {
	t.Helper()
	decoded := mustDecodeHex(t, value)
	if len(decoded) != 32 {
		t.Fatalf("32-byte value has length %d", len(decoded))
	}
	var result T
	copy(result[:], decoded)
	return result
}

func mustHex32WithoutTest[T ~[32]byte](value string) T {
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != 32 {
		panic(fmt.Sprintf("invalid fixed digest %q", value))
	}
	var result T
	copy(result[:], decoded)
	return result
}

func repeated16[T ~[16]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func repeated32[T ~[32]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func optionalString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
