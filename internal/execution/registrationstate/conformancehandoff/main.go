// Command conformancehandoff is a local-only cross-language test fixture. It
// is not an IPC protocol, daemon endpoint, or product transport.
package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const maximumRequestBytes = 32 << 10

type request struct {
	ExactPlanHex string       `json:"exactPlanHex"`
	RoleBindings roleBindings `json:"roleBindings"`
}

type roleBindings struct {
	InstallationIDHex                   string   `json:"installationIdHex"`
	EpochDigestHex                      string   `json:"epochDigestHex"`
	SourceManifestDigestHex             string   `json:"sourceManifestDigestHex"`
	InlineInputDigestHex                string   `json:"inlineInputDigestHex"`
	RuntimeBundleManifestDigestHex      string   `json:"runtimeBundleManifestDigestHex"`
	ProfileReviewAttestationDigestHexes []string `json:"profileReviewAttestationDigestHexes"`
	ProfileRegistryEntryDigestHex       string   `json:"profileRegistryEntryDigestHex"`
	BackendValidationRecordDigestHex    string   `json:"backendValidationRecordDigestHex"`
	BackendConfigurationDigestHex       string   `json:"backendConfigurationDigestHex"`
	TrustSnapshotDigestHex              string   `json:"trustSnapshotDigestHex"`
	PolicyDecisionDigestHex             string   `json:"policyDecisionDigestHex"`
}

type response struct {
	Accepted                      bool   `json:"accepted"`
	Classification                string `json:"classification,omitempty"`
	StateUnchanged                bool   `json:"stateUnchanged"`
	RetainedExactPlanHex          string `json:"retainedExactPlanHex,omitempty"`
	RecomputedPlanDigestHex       string `json:"recomputedPlanDigestHex,omitempty"`
	WireRegistrationPlanDigestHex string `json:"wireRegistrationPlanDigestHex,omitempty"`
}

func main() {
	result, err := run(context.Background(), os.Stdin)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(result); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "encode conformance response")
		os.Exit(2)
	}
}

func run(ctx context.Context, input io.Reader) (response, error) {
	decoder := json.NewDecoder(io.LimitReader(input, maximumRequestBytes+1))
	decoder.DisallowUnknownFields()
	var message request
	if err := decoder.Decode(&message); err != nil {
		return response{}, errors.New("decode closed conformance handoff")
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return response{}, errors.New("conformance handoff has trailing data")
	}
	exactPlan, err := decodeRequiredHex(message.ExactPlanHex, v0candidate.ExecutionPlanMaxCBORBytes)
	if err != nil {
		return response{}, fmt.Errorf("decode exact plan: %w", err)
	}
	bindings, err := decodeBindings(message.RoleBindings)
	if err != nil {
		return response{}, err
	}

	directory, err := os.MkdirTemp("", "capsule-registration-handoff-*")
	if err != nil {
		return response{}, errors.New("create conformance store directory")
	}
	defer func() { _ = os.RemoveAll(directory) }()
	storePath := filepath.Join(directory, "registration-state.json")
	store, err := registrationstate.NewFixedFileStore(storePath, ordinaryInitialState())
	if err != nil {
		return response{}, fmt.Errorf("create conformance store: %w", err)
	}
	before, err := os.ReadFile(storePath) //nolint:gosec // G304: storePath is a tempdir path this function created above, not caller-supplied
	if err != nil {
		return response{}, errors.New("read initial conformance store")
	}
	component, err := registrationstate.New(registrationstate.Options{
		Store:       store,
		Clock:       fixedClock(1_785_456_000),
		Identifiers: fixedIdentifierSource{identifier: repeated16[v0candidate.RegistrationID](0x77)},
	})
	if err != nil {
		return response{}, fmt.Errorf("create registration component: %w", err)
	}
	issued, registrationErr := component.RegisterPlan(
		ctx,
		registrationstate.AuthenticatedCallContext{
			Authenticated: true,
			Role:          registrationstate.CallerDaemon,
			Purpose:       registrationstate.RegisterPlanPurpose,
		},
		exactPlan,
		bindings,
	)
	after, err := os.ReadFile(storePath) //nolint:gosec // G304: storePath is a tempdir path this function created above, not caller-supplied
	if err != nil {
		return response{}, errors.New("read final conformance store")
	}
	if registrationErr != nil {
		classification, ok := registrationstate.ErrorClassification(registrationErr)
		if !ok {
			return response{}, fmt.Errorf("unclassified registration refusal: %w", registrationErr)
		}
		return response{
			Accepted:       false,
			Classification: string(classification),
			StateUnchanged: bytes.Equal(before, after),
		}, nil
	}
	record, err := component.ResolveUsable(ctx, issued.View().RegistrationID)
	if err != nil {
		return response{}, fmt.Errorf("resolve retained registration: %w", err)
	}
	wirePlanDigest := issued.View().PlanDigest
	return response{
		Accepted:                      true,
		StateUnchanged:                bytes.Equal(before, after),
		RetainedExactPlanHex:          hex.EncodeToString(record.ExactPlanBytes),
		RecomputedPlanDigestHex:       hex.EncodeToString(record.RecomputedPlanDigest[:]),
		WireRegistrationPlanDigestHex: hex.EncodeToString(wirePlanDigest[:]),
	}, nil
}

func decodeBindings(input roleBindings) (v0candidate.ExecutionPlanRoleBindings, error) {
	installationID, err := decodeOptionalHex16[v0candidate.InstallationID](input.InstallationIDHex)
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode installation role: %w", err)
	}
	epochDigest, err := decodeOptionalHex32[v0candidate.TrustEpochDigest](input.EpochDigestHex)
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode epoch role: %w", err)
	}
	sourceManifestDigest, err := decodeOptionalHex32[v0candidate.SourceManifestDigest](input.SourceManifestDigestHex)
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode source role: %w", err)
	}
	inlineInputDigest, err := decodeOptionalHex32[v0candidate.InlineInputDigest](input.InlineInputDigestHex)
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode input role: %w", err)
	}
	runtimeBundleDigest, err := decodeOptionalHex32[v0candidate.RuntimeBundleManifestDigest](input.RuntimeBundleManifestDigestHex)
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode runtime role: %w", err)
	}
	reviews := make([]v0candidate.ProfileReviewAttestationDigest, len(input.ProfileReviewAttestationDigestHexes))
	for index, value := range input.ProfileReviewAttestationDigestHexes {
		reviews[index], err = decodeOptionalHex32[v0candidate.ProfileReviewAttestationDigest](value)
		if err != nil {
			return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode profile review role: %w", err)
		}
	}
	profileRegistryDigest, err := decodeOptionalHex32[v0candidate.ProfileRegistryEntryDigest](input.ProfileRegistryEntryDigestHex)
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode profile registry role: %w", err)
	}
	backendValidationDigest, err := decodeOptionalHex32[v0candidate.BackendValidationRecordDigest](input.BackendValidationRecordDigestHex)
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode backend validation role: %w", err)
	}
	backendConfigurationDigest, err := decodeOptionalHex32[v0candidate.BackendConfigurationDigest](input.BackendConfigurationDigestHex)
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode backend configuration role: %w", err)
	}
	trustSnapshotDigest, err := decodeOptionalHex32[v0candidate.TrustSnapshotDigest](input.TrustSnapshotDigestHex)
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode trust snapshot role: %w", err)
	}
	policyDecisionDigest, err := decodeOptionalHex32[v0candidate.PolicyDecisionDigest](input.PolicyDecisionDigestHex)
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, fmt.Errorf("decode policy role: %w", err)
	}
	return v0candidate.ExecutionPlanRoleBindings{
		InstallationID:                  installationID,
		EpochDigest:                     epochDigest,
		SourceManifestDigest:            sourceManifestDigest,
		InlineInputDigest:               inlineInputDigest,
		RuntimeBundleManifestDigest:     runtimeBundleDigest,
		ProfileReviewAttestationDigests: reviews,
		ProfileRegistryEntryDigest:      profileRegistryDigest,
		BackendValidationRecordDigest:   backendValidationDigest,
		BackendConfigurationDigest:      backendConfigurationDigest,
		TrustSnapshotDigest:             trustSnapshotDigest,
		PolicyDecisionDigest:            policyDecisionDigest,
	}, nil
}

func decodeRequiredHex(value string, maximumBytes int) ([]byte, error) {
	if value == "" {
		return nil, errors.New("required bytes are missing")
	}
	if len(value) > maximumBytes*2 {
		return nil, errors.New("hex bytes exceed the conformance bound")
	}
	decoded, err := hex.DecodeString(value)
	if err != nil {
		return nil, errors.New("invalid hexadecimal bytes")
	}
	return decoded, nil
}

func decodeOptionalHex16[T ~[16]byte](value string) (T, error) {
	var result T
	if value == "" {
		return result, nil
	}
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != len(result) {
		return result, errors.New("role must contain exactly 16 bytes")
	}
	copy(result[:], decoded)
	return result, nil
}

func decodeOptionalHex32[T ~[32]byte](value string) (T, error) {
	var result T
	if value == "" {
		return result, nil
	}
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != len(result) {
		return result, errors.New("role must contain exactly 32 bytes")
	}
	copy(result[:], decoded)
	return result, nil
}

func ordinaryInitialState() registrationstate.InitialState {
	return registrationstate.InitialState{
		InstallationID:           repeated16[v0candidate.InstallationID](0x11),
		SupervisorID:             repeated16[v0candidate.SupervisorID](0x55),
		EpochSequence:            7,
		EpochDigest:              repeated32[v0candidate.TrustEpochDigest](0x22),
		TrustPhase:               registrationstate.TrustStable,
		TimeHighWaterUnixSeconds: 1_785_456_000,
	}
}

type fixedClock uint64

func (clock fixedClock) ObserveUnixSeconds(context.Context) (uint64, error) {
	return uint64(clock), nil
}

type fixedIdentifierSource struct {
	identifier v0candidate.RegistrationID
}

func (source fixedIdentifierSource) NewRegistrationID(context.Context) (v0candidate.RegistrationID, error) {
	return source.identifier, nil
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
