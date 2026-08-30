package runtimec2bpassive

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"capsule.local/capsule/internal/protocol/strictjson"
)

const (
	BindingBytes       = 8221
	BindingSHA256      = "3540d5224bdc81edbceafa1f0f17ac119904a70feab604957ab349dd116961a6"
	C1Bytes            = 9289
	C1SHA256           = "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e"
	C2ABytes           = 26850
	C2ASHA256          = "d4ce88888186266f5d251e6246c889b1fd46d7746bb0ba56bcc4b3ce4675992f"
	ForkBytes          = 1895
	ForkSHA256         = "41350bcfc854338ded5e62f77475daf86486351356104dbbf647a8f8b5f11946"
	EvidenceBytes      = 6026
	EvidenceSHA256     = "7de51e88ccc7cf35ee04b822d028f6c2184b72eaa833352fe9fc656b4a7baa13"
	EvidenceSelfDigest = "732301bf8553b0c59b3fe0e4f2b9e070dcc3a1b478e742dc13bd438873b7e488"
)

type Inputs struct {
	C1             []byte
	C2A            []byte
	ForkSupplement []byte
	BuildEvidence  []byte
}

// field-authority-object: capsule.governed-deno-core-c2b-passive-binding v1
type Binding struct {
	ObjectType            string                `json:"objectType"`
	SchemaVersion         int                   `json:"schemaVersion"`
	Domain                string                `json:"domain"`
	Identity              string                `json:"identity"`
	Status                string                `json:"status"`
	Predecessors          Predecessors          `json:"predecessors"`
	ForkSupplement        ForkSupplement        `json:"forkSupplement"`
	BuildEvidence         BuildEvidence         `json:"buildEvidence"`
	FixedFixture          FixedFixture          `json:"fixedFixture"`
	Artifacts             []Artifact            `json:"artifacts"`
	DependencyMergePolicy DependencyMergePolicy `json:"dependencyMergePolicy"`
	WorkStatus            WorkStatus            `json:"workStatus"`
	Limitations           Limitations           `json:"limitations"`
	NextConsumerGate      NextConsumerGate      `json:"nextConsumerGate"`
	Effects               Effects               `json:"effects"`
}

type Predecessors struct {
	C1  Predecessor `json:"c1"`
	C2A Predecessor `json:"c2a"`
}

type Predecessor struct {
	ObjectType  string `json:"objectType"`
	Version     int    `json:"version"`
	Identity    string `json:"identity"`
	Fixture     string `json:"fixture"`
	Bytes       int    `json:"bytes"`
	SHA256      string `json:"sha256"`
	Consumption string `json:"consumption"`
	Mutation    string `json:"mutation"`
}

type ForkSupplement struct {
	ObjectType      string `json:"objectType"`
	SchemaVersion   int    `json:"schemaVersion"`
	Identity        string `json:"identity"`
	Fixture         string `json:"fixture"`
	Bytes           int    `json:"bytes"`
	SHA256          string `json:"sha256"`
	Repository      string `json:"repository"`
	PullRequest     string `json:"pullRequest"`
	Branch          string `json:"branch"`
	Base            string `json:"base"`
	Head            string `json:"head"`
	Tree            string `json:"tree"`
	RetainedPath    string `json:"retainedPath"`
	AuthorityEffect string `json:"authorityEffect"`
}

type BuildEvidence struct {
	ObjectType              string `json:"objectType"`
	SchemaVersion           int    `json:"schemaVersion"`
	Identity                string `json:"identity"`
	Fixture                 string `json:"fixture"`
	Bytes                   int    `json:"bytes"`
	RawSHA256               string `json:"rawSha256"`
	SelfDigest              string `json:"selfDigest"`
	PredecessorIdentity     string `json:"predecessorIdentity"`
	PredecessorSelfDigest   string `json:"predecessorSelfDigest"`
	PredecessorState        string `json:"predecessorState"`
	Repository              string `json:"repository"`
	PullRequest             string `json:"pullRequest"`
	Branch                  string `json:"branch"`
	Head                    string `json:"head"`
	Tree                    string `json:"tree"`
	RetainedPath            string `json:"retainedPath"`
	ForkSupplementIdentity  string `json:"forkSupplementIdentity"`
	ForkSupplementSHA256    string `json:"forkSupplementSha256"`
	SameHostReproducibility string `json:"sameHostReproducibility"`
	IndependentBuilder      bool   `json:"independentBuilder"`
}

type FixedFixture struct {
	Source         Source         `json:"source"`
	SourceManifest SourceManifest `json:"sourceManifest"`
	Input          Input          `json:"input"`
	Completion     Completion     `json:"completion"`
	Authority      Authority      `json:"authority"`
}

type Source struct {
	Path         string `json:"path"`
	MediaType    string `json:"mediaType"`
	UTF8         string `json:"utf8"`
	Bytes        int    `json:"bytes"`
	SHA256       string `json:"sha256"`
	MaximumBytes int    `json:"maximumBytes"`
}

type SourceManifest struct {
	ObjectType string `json:"objectType"`
	Version    int    `json:"version"`
	MediaType  string `json:"mediaType"`
	Bytes      int    `json:"bytes"`
	SHA256     string `json:"sha256"`
}

type Input struct {
	MediaType    string `json:"mediaType"`
	Encoding     string `json:"encoding"`
	UTF8         string `json:"utf8"`
	Bytes        int    `json:"bytes"`
	SHA256       string `json:"sha256"`
	MaximumBytes int    `json:"maximumBytes"`
}

type Completion struct {
	MediaType    string `json:"mediaType"`
	UTF8         string `json:"utf8"`
	Bytes        int    `json:"bytes"`
	SHA256       string `json:"sha256"`
	MaximumBytes int    `json:"maximumBytes"`
	Commit       string `json:"commit"`
}

type Authority struct {
	WorkloadBytes     string `json:"workloadBytes"`
	ArbitrarySource   bool   `json:"arbitrarySource"`
	CallerPaths       bool   `json:"callerPaths"`
	CallerArguments   bool   `json:"callerArguments"`
	CallerEnvironment bool   `json:"callerEnvironment"`
}

type Artifact struct {
	Role         string `json:"role"`
	Bytes        int    `json:"bytes"`
	SHA256       string `json:"sha256"`
	PackageCount *int   `json:"packageCount"`
}

type DependencyMergePolicy struct {
	State           string       `json:"state"`
	Dependencies    []Dependency `json:"dependencies"`
	ChangeRule      string       `json:"changeRule"`
	InPlaceMutation string       `json:"inPlaceMutation"`
}

type Dependency struct {
	Role          string `json:"role"`
	Repository    string `json:"repository"`
	PullRequest   int    `json:"pullRequest"`
	Head          string `json:"head"`
	Tree          string `json:"tree"`
	RequiredState string `json:"requiredState"`
}

type WorkStatus struct {
	PassiveReconciliation            string `json:"passiveReconciliation"`
	ParentGovernedRuntime            string `json:"parentGovernedRuntime"`
	C2BComposedProfileGuestExecution string `json:"c2bComposedProfileGuestExecution"`
	RuntimeProfileAdmission          string `json:"runtimeProfileAdmission"`
	Runtime001                       string `json:"runtime001"`
	VMM001                           string `json:"vmm001"`
}

type Limitations struct {
	SameHostBuildEvidence    string `json:"sameHostBuildEvidence"`
	RuntimeRoot              string `json:"runtimeRoot"`
	ComposedProfile          string `json:"composedProfile"`
	GuestExecution           string `json:"guestExecution"`
	SecurityOrReadinessClaim string `json:"securityOrReadinessClaim"`
	ProductConsumer          string `json:"productConsumer"`
}

type NextConsumerGate struct {
	State           string   `json:"state"`
	Consumer        string   `json:"consumer"`
	Requirements    []string `json:"requirements"`
	Implemented     bool     `json:"implemented"`
	AdmissionEffect bool     `json:"admissionEffect"`
}

type Effects struct {
	Process    bool `json:"process"`
	Runtime    bool `json:"runtime"`
	Profile    bool `json:"profile"`
	Consumer   bool `json:"consumer"`
	Adapter    bool `json:"adapter"`
	Backend    bool `json:"backend"`
	VM         bool `json:"vm"`
	Guest      bool `json:"guest"`
	Credential bool `json:"credential"`
	Signing    bool `json:"signing"`
	Release    bool `json:"release"`
	Admission  bool `json:"admission"`
}

func Decode(exact []byte, inputs Inputs) (*Binding, error) {
	if err := checkExact("C2B_BINDING", exact, BindingBytes, BindingSHA256); err != nil {
		return nil, err
	}
	binding, err := decodeClosed(exact)
	if err != nil {
		return nil, err
	}
	if err := Validate(binding); err != nil {
		return nil, err
	}
	if err := validateInputs(binding, inputs); err != nil {
		return nil, err
	}
	return strictjson.Clone(binding), nil
}

func decodeClosed(exact []byte) (*Binding, error) {
	if err := rejectDuplicateJSON(exact); err != nil {
		return nil, fmt.Errorf("C2B_BINDING_DUPLICATE: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.DisallowUnknownFields()
	var binding Binding
	if err := decoder.Decode(&binding); err != nil {
		return nil, fmt.Errorf("C2B_BINDING_DECODE: %w", err)
	}
	if err := strictjson.RequireEOF(decoder, "C2B_BINDING_TRAILING"); err != nil {
		return nil, err
	}
	return &binding, nil
}

func Validate(binding *Binding) error {
	if binding == nil {
		return errors.New("C2B_BINDING_MISMATCH: nil")
	}
	canonical, err := json.MarshalIndent(binding, "", "  ")
	if err != nil {
		return fmt.Errorf("C2B_BINDING_MISMATCH: encode: %w", err)
	}
	canonical = append(canonical, '\n')
	if len(canonical) != BindingBytes || digest(canonical) != BindingSHA256 {
		return errors.New("C2B_BINDING_MISMATCH: semantic known answer")
	}
	return nil
}

func validateInputs(binding *Binding, inputs Inputs) error {
	checks := []struct {
		name   string
		value  []byte
		length int
		digest string
	}{
		{"C1", inputs.C1, C1Bytes, C1SHA256},
		{"C2A", inputs.C2A, C2ABytes, C2ASHA256},
		{"FORK", inputs.ForkSupplement, ForkBytes, ForkSHA256},
		{"EVIDENCE", inputs.BuildEvidence, EvidenceBytes, EvidenceSHA256},
	}
	for _, check := range checks {
		if err := checkExact(check.name, check.value, check.length, check.digest); err != nil {
			return err
		}
		if err := rejectDuplicateJSON(check.value); err != nil {
			return fmt.Errorf("C2B_%s_DUPLICATE: %w", check.name, err)
		}
	}

	if err := validateEvidenceSelfDigest(inputs.BuildEvidence); err != nil {
		return err
	}

	var c2a struct {
		KnownAnswer struct {
			Source struct {
				UTF8                 string `json:"utf8"`
				Bytes                int    `json:"bytes"`
				SHA256               string `json:"sha256"`
				SourceManifestBytes  int    `json:"sourceManifestBytes"`
				SourceManifestSHA256 string `json:"sourceManifestSha256"`
			} `json:"source"`
			CanonicalInput struct {
				UTF8   string `json:"utf8"`
				Bytes  int    `json:"bytes"`
				SHA256 string `json:"sha256"`
			} `json:"canonicalInput"`
			ExpectedCompletion struct {
				JSONUTF8   string `json:"jsonUtf8"`
				JSONBytes  int    `json:"jsonBytes"`
				JSONSHA256 string `json:"jsonSha256"`
				Commit     string `json:"commit"`
			} `json:"expectedCompletion"`
		} `json:"knownAnswer"`
	}
	if err := json.Unmarshal(inputs.C2A, &c2a); err != nil {
		return fmt.Errorf("C2B_C2A_CROSS_LINK: %w", err)
	}
	if c2a.KnownAnswer.Source.UTF8 != binding.FixedFixture.Source.UTF8 ||
		c2a.KnownAnswer.Source.Bytes != binding.FixedFixture.Source.Bytes ||
		c2a.KnownAnswer.Source.SHA256 != binding.FixedFixture.Source.SHA256 ||
		c2a.KnownAnswer.Source.SourceManifestBytes != binding.FixedFixture.SourceManifest.Bytes ||
		c2a.KnownAnswer.Source.SourceManifestSHA256 != binding.FixedFixture.SourceManifest.SHA256 ||
		c2a.KnownAnswer.CanonicalInput.UTF8 != binding.FixedFixture.Input.UTF8 ||
		c2a.KnownAnswer.CanonicalInput.Bytes != binding.FixedFixture.Input.Bytes ||
		c2a.KnownAnswer.CanonicalInput.SHA256 != binding.FixedFixture.Input.SHA256 ||
		c2a.KnownAnswer.ExpectedCompletion.JSONUTF8 != binding.FixedFixture.Completion.UTF8 ||
		c2a.KnownAnswer.ExpectedCompletion.JSONBytes != binding.FixedFixture.Completion.Bytes ||
		c2a.KnownAnswer.ExpectedCompletion.JSONSHA256 != binding.FixedFixture.Completion.SHA256 ||
		c2a.KnownAnswer.ExpectedCompletion.Commit != binding.FixedFixture.Completion.Commit {
		return errors.New("C2B_C2A_CROSS_LINK: known answer")
	}

	var fork struct {
		Identity     string `json:"identity"`
		Predecessors struct {
			C1 struct {
				Bytes  int    `json:"bytes"`
				SHA256 string `json:"sha256"`
			} `json:"c1"`
			C2A struct {
				Bytes  int    `json:"bytes"`
				SHA256 string `json:"sha256"`
			} `json:"c2a"`
		} `json:"predecessors"`
		Source struct {
			Bytes  int    `json:"bytes"`
			SHA256 string `json:"sha256"`
		} `json:"source"`
		Input struct {
			Bytes  int    `json:"bytes"`
			SHA256 string `json:"sha256"`
		} `json:"input"`
		Completion struct {
			Bytes  int    `json:"bytes"`
			SHA256 string `json:"sha256"`
		} `json:"completion"`
	}
	if err := json.Unmarshal(inputs.ForkSupplement, &fork); err != nil {
		return fmt.Errorf("C2B_FORK_CROSS_LINK: %w", err)
	}
	if fork.Identity != binding.ForkSupplement.Identity ||
		fork.Predecessors.C1.Bytes != C1Bytes || fork.Predecessors.C1.SHA256 != C1SHA256 ||
		fork.Predecessors.C2A.Bytes != C2ABytes || fork.Predecessors.C2A.SHA256 != C2ASHA256 ||
		fork.Source.Bytes != binding.FixedFixture.Source.Bytes || fork.Source.SHA256 != binding.FixedFixture.Source.SHA256 ||
		fork.Input.Bytes != binding.FixedFixture.Input.Bytes || fork.Input.SHA256 != binding.FixedFixture.Input.SHA256 ||
		fork.Completion.Bytes != binding.FixedFixture.Completion.Bytes || fork.Completion.SHA256 != binding.FixedFixture.Completion.SHA256 {
		return errors.New("C2B_FORK_CROSS_LINK: identity or known answer")
	}

	var evidence struct {
		Identity   string `json:"identity"`
		SelfDigest struct {
			SHA256 string `json:"sha256"`
		} `json:"selfDigest"`
		ImmutableBuildOnlySupplement struct {
			Identity string `json:"identity"`
			SHA256   string `json:"sha256"`
		} `json:"immutableBuildOnlySupplement"`
		PredecessorBuildEvidence struct {
			Identity   string `json:"identity"`
			SelfDigest string `json:"selfDigest"`
		} `json:"predecessorBuildEvidence"`
		Sources struct {
			Deno struct {
				Base   string `json:"base"`
				Commit string `json:"commit"`
				Tree   string `json:"tree"`
			} `json:"deno"`
		} `json:"sources"`
		Artifacts struct {
			Binary struct {
				Bytes  int    `json:"bytes"`
				SHA256 string `json:"sha256"`
			} `json:"binary"`
			Snapshot struct {
				Bytes  int    `json:"bytes"`
				SHA256 string `json:"sha256"`
			} `json:"snapshot"`
			TwoFileBundle struct {
				Bytes  int    `json:"bytes"`
				SHA256 string `json:"sha256"`
			} `json:"twoFileBundle"`
			BundleManifest struct {
				Bytes  int    `json:"bytes"`
				SHA256 string `json:"sha256"`
			} `json:"bundleManifest"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal(inputs.BuildEvidence, &evidence); err != nil {
		return fmt.Errorf("C2B_EVIDENCE_CROSS_LINK: %w", err)
	}
	if evidence.Identity != binding.BuildEvidence.Identity || evidence.SelfDigest.SHA256 != binding.BuildEvidence.SelfDigest ||
		evidence.ImmutableBuildOnlySupplement.Identity != binding.ForkSupplement.Identity || evidence.ImmutableBuildOnlySupplement.SHA256 != binding.ForkSupplement.SHA256 ||
		evidence.PredecessorBuildEvidence.Identity != binding.BuildEvidence.PredecessorIdentity || evidence.PredecessorBuildEvidence.SelfDigest != binding.BuildEvidence.PredecessorSelfDigest ||
		evidence.Sources.Deno.Base != binding.ForkSupplement.Base || evidence.Sources.Deno.Commit != binding.ForkSupplement.Head || evidence.Sources.Deno.Tree != binding.ForkSupplement.Tree {
		return errors.New("C2B_EVIDENCE_CROSS_LINK: identity, ancestry, or supplement")
	}
	expectedArtifacts := binding.Artifacts[1:5]
	observedArtifacts := []struct {
		Bytes  int
		SHA256 string
	}{
		{evidence.Artifacts.Binary.Bytes, evidence.Artifacts.Binary.SHA256},
		{evidence.Artifacts.Snapshot.Bytes, evidence.Artifacts.Snapshot.SHA256},
		{evidence.Artifacts.TwoFileBundle.Bytes, evidence.Artifacts.TwoFileBundle.SHA256},
		{evidence.Artifacts.BundleManifest.Bytes, evidence.Artifacts.BundleManifest.SHA256},
	}
	for index, observed := range observedArtifacts {
		if observed.Bytes != expectedArtifacts[index].Bytes || observed.SHA256 != expectedArtifacts[index].SHA256 {
			return fmt.Errorf("C2B_EVIDENCE_CROSS_LINK: artifact %d", index)
		}
	}
	return nil
}

func validateEvidenceSelfDigest(exact []byte) error {
	needle := []byte(`"sha256": "` + EvidenceSelfDigest + `"`)
	replacement := []byte(`"sha256": null`)
	if bytes.Count(exact, needle) != 1 {
		return errors.New("C2B_EVIDENCE_SELF_DIGEST: field")
	}
	normalized := bytes.Replace(exact, needle, replacement, 1)
	if digest(normalized) != EvidenceSelfDigest {
		return errors.New("C2B_EVIDENCE_SELF_DIGEST: mismatch")
	}
	return nil
}

func checkExact(name string, exact []byte, length int, expectedDigest string) error {
	if len(exact) != length {
		return fmt.Errorf("C2B_%s_LENGTH: got %d want %d", name, len(exact), length)
	}
	if digest(exact) != expectedDigest {
		return fmt.Errorf("C2B_%s_DIGEST", name)
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
	return strictjson.RequireEOF(decoder, "C2B_BINDING_TRAILING")
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
