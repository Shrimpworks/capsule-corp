package runtimec2bpassivev2

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const (
	BindingBytes  = 7115
	BindingSHA256 = "c59f7fdd27834dd7be5a05a3c44a973d6ffa99869b9b99e2531045926827190a"
	C1Bytes       = 9289
	C1SHA256      = "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e"
	C2ABytes      = 26850
	C2ASHA256     = "d4ce88888186266f5d251e6246c889b1fd46d7746bb0ba56bcc4b3ce4675992f"
	C2BV1Bytes    = 8221
	C2BV1SHA256   = "3540d5224bdc81edbceafa1f0f17ac119904a70feab604957ab349dd116961a6"
)

// field-authority-object: capsule.governed-deno-core-c2b-passive-binding v2
type Binding struct {
	ObjectType           string               `json:"objectType"`
	SchemaVersion        int                  `json:"schemaVersion"`
	Domain               string               `json:"domain"`
	Identity             string               `json:"identity"`
	Status               string               `json:"status"`
	Predecessors         Predecessors         `json:"predecessors"`
	ArchiveEvidence      ArchiveEvidence      `json:"archiveEvidence"`
	ConstructedArtifacts []Artifact           `json:"constructedArtifacts"`
	RuntimeManifest      RuntimeManifest      `json:"runtimeManifestCandidate"`
	HostPreflightHarness HostPreflightHarness `json:"hostPreflightHarness"`
	Unresolved           Unresolved           `json:"unresolved"`
	WorkStatus           WorkStatus           `json:"workStatus"`
	NextConsumerGate     NextConsumerGate     `json:"nextConsumerGate"`
	Effects              Effects              `json:"effects"`
}

type Predecessors struct {
	C1    Predecessor `json:"c1"`
	C2A   Predecessor `json:"c2a"`
	C2BV1 Predecessor `json:"c2bV1"`
}

type Predecessor struct {
	Identity    string `json:"identity"`
	Fixture     string `json:"fixture"`
	Bytes       int    `json:"bytes"`
	SHA256      string `json:"sha256"`
	Consumption string `json:"consumption"`
	Mutation    string `json:"mutation"`
}

type ArchiveEvidence struct {
	Repository                     string           `json:"repository"`
	PullRequest                    int              `json:"pullRequest"`
	PullRequestState               string           `json:"pullRequestState"`
	MergeCommit                    string           `json:"mergeCommit"`
	ReviewedHead                   string           `json:"reviewedHead"`
	Tree                           string           `json:"tree"`
	BaselineParent                 string           `json:"baselineParent"`
	RetainedRoot                   string           `json:"retainedRoot"`
	RuntimeBundleCandidate         RetainedManifest `json:"runtimeBundleCandidate"`
	ArtifactClosureReport          RetainedManifest `json:"artifactClosureReport"`
	EvidenceDigests                EvidenceDigests  `json:"evidenceDigests"`
	IndependentPassiveVerification string           `json:"independentPassiveVerification"`
	AuthorityEffect                string           `json:"authorityEffect"`
}

type RetainedManifest struct {
	RetainedPath               string `json:"retainedPath"`
	Bytes                      int    `json:"bytes"`
	SHA256                     string `json:"sha256"`
	SelfDigestOfNullFormSHA256 string `json:"selfDigestOfNullFormSha256"`
}

type EvidenceDigests struct {
	SourceRefs         string `json:"sourceRefs"`
	Verification       string `json:"verification"`
	Mutations          string `json:"mutations"`
	SBOM               string `json:"sbom"`
	Licenses           string `json:"licenses"`
	UnsignedProvenance string `json:"unsignedProvenance"`
}

type Artifact struct {
	Role         string `json:"role"`
	SelectedPath string `json:"selectedPath"`
	Mode         string `json:"mode"`
	Bytes        int    `json:"bytes"`
	SHA256       string `json:"sha256"`
	Admitted     bool   `json:"admitted"`
}

type RuntimeManifest struct {
	Role               string  `json:"role"`
	SelectedPath       string  `json:"selectedPath"`
	Bytes              int     `json:"bytes"`
	SHA256             string  `json:"sha256"`
	CanonicalAuthority *string `json:"canonicalAuthority"`
	Admitted           bool    `json:"admitted"`
}

type HostPreflightHarness struct {
	Role                    string  `json:"role"`
	SelectedPath            string  `json:"selectedPath"`
	Mode                    string  `json:"mode"`
	Bytes                   int     `json:"bytes"`
	SHA256                  string  `json:"sha256"`
	Usage                   string  `json:"usage"`
	FinalHostRunnerIdentity *string `json:"finalHostRunnerIdentity"`
	Final                   bool    `json:"final"`
	CallsLibkrun            bool    `json:"callsLibkrun"`
	Admitted                bool    `json:"admitted"`
}

type Unresolved struct {
	Refusal                  string  `json:"refusal"`
	FinalHostRunnerIdentity  *string `json:"finalHostRunnerIdentity"`
	SeparateFirmwareIdentity *string `json:"separateFirmwareIdentity"`
	ComposedProfileIdentity  *string `json:"composedProfileIdentity"`
	ComposedProfileDigest    *string `json:"composedProfileDigest"`
	CPUTimeLimitMS           *int64  `json:"cpuTimeLimitMs"`
	HostVMMMemoryLimitBytes  *int64  `json:"hostVmmMemoryLimitBytes"`
	ScratchMaximumBytes      *int64  `json:"scratchMaximumBytes"`
	GuestEvidence            *string `json:"guestEvidence"`
	RuntimeManifestAdmission *string `json:"runtimeManifestAdmission"`
	RuntimeProfileAdmission  *string `json:"runtimeProfileAdmission"`
	BackendAdmission         *string `json:"backendAdmission"`
}

type WorkStatus struct {
	NoGuestBuildArtifactClosure      string `json:"noGuestBuildArtifactClosure"`
	ParentGovernedRuntime            string `json:"parentGovernedRuntime"`
	C2BComposedProfileGuestExecution string `json:"c2bComposedProfileGuestExecution"`
	RuntimeProfileAdmission          string `json:"runtimeProfileAdmission"`
	Runtime001                       string `json:"runtime001"`
	VMM001                           string `json:"vmm001"`
}

type NextConsumerGate struct {
	State           string   `json:"state"`
	Consumer        string   `json:"consumer"`
	Requirements    []string `json:"requirements"`
	Implemented     bool     `json:"implemented"`
	AdmissionEffect bool     `json:"admissionEffect"`
}

type Effects struct {
	Process     bool `json:"process"`
	Runtime     bool `json:"runtime"`
	Profile     bool `json:"profile"`
	Consumer    bool `json:"consumer"`
	Adapter     bool `json:"adapter"`
	Backend     bool `json:"backend"`
	VM          bool `json:"vm"`
	Guest       bool `json:"guest"`
	Credential  bool `json:"credential"`
	Signing     bool `json:"signing"`
	Publication bool `json:"publication"`
	Release     bool `json:"release"`
	Admission   bool `json:"admission"`
}

func Decode(exact []byte) (*Binding, error) {
	if len(exact) > BindingBytes {
		return nil, fmt.Errorf("C2B_V2_BINDING_CAP: got %d want at most %d", len(exact), BindingBytes)
	}
	if err := checkExact(exact); err != nil {
		return nil, err
	}
	if err := rejectDuplicateJSON(exact); err != nil {
		return nil, fmt.Errorf("C2B_V2_BINDING_DUPLICATE: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.DisallowUnknownFields()
	var binding Binding
	if err := decoder.Decode(&binding); err != nil {
		return nil, fmt.Errorf("C2B_V2_BINDING_DECODE: %w", err)
	}
	if err := requireEOF(decoder); err != nil {
		return nil, err
	}
	if err := Validate(&binding); err != nil {
		return nil, err
	}
	return clone(&binding), nil
}

func Validate(binding *Binding) error {
	if binding == nil {
		return errors.New("C2B_V2_BINDING_MISMATCH: nil")
	}
	canonical, err := json.MarshalIndent(binding, "", "  ")
	if err != nil {
		return fmt.Errorf("C2B_V2_BINDING_MISMATCH: encode: %w", err)
	}
	canonical = append(canonical, '\n')
	if len(canonical) != BindingBytes || digest(canonical) != BindingSHA256 {
		return errors.New("C2B_V2_BINDING_MISMATCH: semantic known answer")
	}
	return nil
}

func checkExact(exact []byte) error {
	if len(exact) != BindingBytes {
		return fmt.Errorf("C2B_V2_BINDING_LENGTH: got %d want %d", len(exact), BindingBytes)
	}
	if digest(exact) != BindingSHA256 {
		return errors.New("C2B_V2_BINDING_DIGEST")
	}
	return nil
}

func digest(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

func rejectDuplicateJSON(exact []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.UseNumber()
	if err := readJSONValue(decoder); err != nil {
		return err
	}
	return requireEOF(decoder)
}

func readJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delimiter {
	case '{':
		seen := make(map[string]struct{})
		for decoder.More() {
			keyToken, keyErr := decoder.Token()
			if keyErr != nil {
				return keyErr
			}
			key, keyOK := keyToken.(string)
			if !keyOK {
				return errors.New("object key is not a string")
			}
			if _, exists := seen[key]; exists {
				return fmt.Errorf("duplicate key %q", key)
			}
			seen[key] = struct{}{}
			if err := readJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err = decoder.Token()
		return err
	case '[':
		for decoder.More() {
			if err := readJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err = decoder.Token()
		return err
	default:
		return errors.New("unexpected JSON delimiter")
	}
}

func requireEOF(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); err == io.EOF {
		return nil
	} else if err != nil {
		return fmt.Errorf("C2B_V2_BINDING_TRAILING: %w", err)
	}
	return errors.New("C2B_V2_BINDING_TRAILING: extra JSON value")
}

func clone(binding *Binding) *Binding {
	copy := *binding
	copy.ConstructedArtifacts = append([]Artifact(nil), binding.ConstructedArtifacts...)
	copy.NextConsumerGate.Requirements = append([]string(nil), binding.NextConsumerGate.Requirements...)
	return &copy
}
