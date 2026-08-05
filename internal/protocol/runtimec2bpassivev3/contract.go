package runtimec2bpassivev3

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
	BindingBytes       = 18357
	BindingSHA256      = "d72327bba369484a56db7d543a32e8bbd4eac403230ac65d63709ac3ba3bbdfb"
	ContractSHA256     = "8b1ec936a7b56370716d28557125e46866dea8f21a149704a01f251a0dddbcc1"
	SemanticBytes      = 18523
	SemanticSHA256     = "e550212a9e10f2ecc30cb862fe4c228e0ddbcbb0dc807b14b154e3eead7a89d0"
	C1Bytes            = 9289
	C1SHA256           = "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e"
	C2ABytes           = 26850
	C2ASHA256          = "d4ce88888186266f5d251e6246c889b1fd46d7746bb0ba56bcc4b3ce4675992f"
	C2BV1Bytes         = 8221
	C2BV1SHA256        = "3540d5224bdc81edbceafa1f0f17ac119904a70feab604957ab349dd116961a6"
	C2BV2Bytes         = 7115
	C2BV2SHA256        = "c59f7fdd27834dd7be5a05a3c44a973d6ffa99869b9b99e2531045926827190a"
	zeroContractSHA256 = "0000000000000000000000000000000000000000000000000000000000000000"
)

// field-authority-object: capsule.governed-deno-core-c2b-passive-binding v3
type Binding struct {
	ObjectType                  string          `json:"objectType"`
	SchemaVersion               int             `json:"schemaVersion"`
	Domain                      string          `json:"domain"`
	Identity                    string          `json:"identity"`
	Status                      string          `json:"status"`
	Predecessors                json.RawMessage `json:"predecessors"`
	GovernedInputs              json.RawMessage `json:"governedInputs"`
	ArtifactRoles               json.RawMessage `json:"artifactRoles"`
	HostRunnerContract          json.RawMessage `json:"hostRunnerContract"`
	DescriptorAndDeviceContract json.RawMessage `json:"descriptorAndDeviceContract"`
	RuntimeContract             json.RawMessage `json:"runtimeContract"`
	ResourceContract            json.RawMessage `json:"resourceContract"`
	TransportContract           json.RawMessage `json:"transportContract"`
	TeardownContract            json.RawMessage `json:"teardownContract"`
	ComposedProfile             json.RawMessage `json:"composedProfile"`
	WorkStatus                  json.RawMessage `json:"workStatus"`
	NextGuestGate               json.RawMessage `json:"nextGuestGate"`
	Effects                     json.RawMessage `json:"effects"`
}

func Decode(exact []byte) (*Binding, error) {
	if len(exact) > BindingBytes {
		return nil, fmt.Errorf("C2B_V3_BINDING_CAP: got %d want at most %d", len(exact), BindingBytes)
	}
	if err := checkExact(exact); err != nil {
		return nil, err
	}
	if err := rejectDuplicateJSON(exact); err != nil {
		return nil, fmt.Errorf("C2B_V3_BINDING_DUPLICATE: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.DisallowUnknownFields()
	var binding Binding
	if err := decoder.Decode(&binding); err != nil {
		return nil, fmt.Errorf("C2B_V3_BINDING_DECODE: %w", err)
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
		return errors.New("C2B_V3_BINDING_MISMATCH: nil")
	}
	canonical, err := json.MarshalIndent(binding, "", "  ")
	if err != nil {
		return fmt.Errorf("C2B_V3_BINDING_MISMATCH: encode: %w", err)
	}
	canonical = append(canonical, '\n')
	var generic any
	if err := json.Unmarshal(canonical, &generic); err != nil {
		return fmt.Errorf("C2B_V3_BINDING_MISMATCH: null scan: %w", err)
	}
	if containsNull(generic) {
		return errors.New("C2B_V3_AUTHORITY_NULL_FORBIDDEN")
	}
	if len(canonical) != SemanticBytes || digest(canonical) != SemanticSHA256 {
		return fmt.Errorf(
			"C2B_V3_BINDING_MISMATCH: semantic known answer: got bytes=%d sha256=%s",
			len(canonical),
			digest(canonical),
		)
	}
	var composed struct {
		ContractDigestSHA256 string `json:"contractDigestSha256"`
	}
	if err := json.Unmarshal(binding.ComposedProfile, &composed); err != nil || composed.ContractDigestSHA256 != ContractSHA256 {
		return errors.New("C2B_V3_COMPOSED_PROFILE_DIGEST")
	}
	return nil
}

func containsNull(value any) bool {
	if value == nil {
		return true
	}
	switch typed := value.(type) {
	case []any:
		for _, child := range typed {
			if containsNull(child) {
				return true
			}
		}
	case map[string]any:
		for _, child := range typed {
			if containsNull(child) {
				return true
			}
		}
	}
	return false
}

func checkExact(exact []byte) error {
	if len(exact) != BindingBytes {
		return fmt.Errorf("C2B_V3_BINDING_LENGTH: got %d want %d", len(exact), BindingBytes)
	}
	if digest(exact) != BindingSHA256 {
		return errors.New("C2B_V3_BINDING_DIGEST")
	}
	zeroForm := bytes.Replace(exact, []byte(ContractSHA256), []byte(zeroContractSHA256), 1)
	if bytes.Equal(zeroForm, exact) || digest(zeroForm) != ContractSHA256 {
		return errors.New("C2B_V3_COMPOSED_PROFILE_DIGEST")
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
		return fmt.Errorf("C2B_V3_BINDING_TRAILING: %w", err)
	}
	return errors.New("C2B_V3_BINDING_TRAILING: extra JSON value")
}

func clone(binding *Binding) *Binding {
	copy := *binding
	copy.Predecessors = append(json.RawMessage(nil), binding.Predecessors...)
	copy.GovernedInputs = append(json.RawMessage(nil), binding.GovernedInputs...)
	copy.ArtifactRoles = append(json.RawMessage(nil), binding.ArtifactRoles...)
	copy.HostRunnerContract = append(json.RawMessage(nil), binding.HostRunnerContract...)
	copy.DescriptorAndDeviceContract = append(json.RawMessage(nil), binding.DescriptorAndDeviceContract...)
	copy.RuntimeContract = append(json.RawMessage(nil), binding.RuntimeContract...)
	copy.ResourceContract = append(json.RawMessage(nil), binding.ResourceContract...)
	copy.TransportContract = append(json.RawMessage(nil), binding.TransportContract...)
	copy.TeardownContract = append(json.RawMessage(nil), binding.TeardownContract...)
	copy.ComposedProfile = append(json.RawMessage(nil), binding.ComposedProfile...)
	copy.WorkStatus = append(json.RawMessage(nil), binding.WorkStatus...)
	copy.NextGuestGate = append(json.RawMessage(nil), binding.NextGuestGate...)
	copy.Effects = append(json.RawMessage(nil), binding.Effects...)
	return &copy
}
