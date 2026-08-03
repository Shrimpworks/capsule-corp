// Package tsapprovedbytecandidate verifies only the retained passive
// TypeScript approved-byte conformance corpus. It is not a product decoder,
// transformation service, registration path, or runtime admission mechanism.
package tsapprovedbytecandidate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
)

const (
	ManifestMaxBytes = 32 * 1024
	ObjectMaxBytes   = 4 * 1024
	SourceMaxBytes   = 262_144
)

type OriginalManifestDigest [32]byte
type ExecutableManifestDigest [32]byte
type TransformationRecordSetDigest [32]byte

type PlanSourceBindings struct {
	OriginalManifest        OriginalManifestDigest
	ExecutableManifest      ExecutableManifestDigest
	TransformationRecordSet TransformationRecordSetDigest
}

type VerifiedFixtureSet struct {
	bytes    map[string][]byte
	bindings PlanSourceBindings
}

func (set *VerifiedFixtureSet) Bytes(path string) []byte {
	return bytes.Clone(set.bytes[path])
}

func (set *VerifiedFixtureSet) Bindings() PlanSourceBindings {
	return set.bindings
}

type manifest struct {
	ManifestVersion string        `json:"manifestVersion"`
	Status          string        `json:"status"`
	FixtureCount    int           `json:"fixtureCount"`
	MutationCount   int           `json:"mutationCount"`
	Caps            caps          `json:"caps"`
	KnownAnswers    []knownAnswer `json:"knownAnswers"`
	Mutations       []mutation    `json:"mutations"`
}

type caps struct {
	SourceFiles            int `json:"sourceFiles"`
	OriginalFileBytes      int `json:"originalFileBytes"`
	OriginalAggregateBytes int `json:"originalAggregateBytes"`
	EmittedFileBytes       int `json:"emittedFileBytes"`
	EmittedAggregateBytes  int `json:"emittedAggregateBytes"`
}

type knownAnswer struct {
	Path       string `json:"path"`
	ByteLength int    `json:"byteLength"`
	SHA256     string `json:"sha256"`
	MediaType  string `json:"mediaType,omitempty"`
}

type mutation struct {
	ID             string `json:"id"`
	Decision       string `json:"decision"`
	Classification string `json:"classification"`
}

type expectedAnswer struct {
	length    int
	digest    string
	mediaType string
}

var expectedAnswers = map[string]expectedAnswer{
	"files/ordinary.ts":                      {391, "697a9bed1307422ecae49d20f60dff7b367535f450fee514b63ffe4d5a22c6db", ""},
	"files/ordinary.js":                      {391, "f91911dd606409fed94c214381533f5ece3e2ae23ea861a3a55192cefad884cd", ""},
	"objects/transformer-profile.cbor":       {215, "ab82cd553d52490fb0e2cc2e6cfc8cc106440a601c99c70cebad43537efc7f97", "application/capsule.typescript-transformer-profile+cbor;v=0"},
	"objects/normalized-options.cbor":        {188, "8a43f134b568e983a7e4e24a763209a61405dc535a710e4799611033e4983b2e", "application/capsule.typescript-normalized-options+cbor;v=0"},
	"objects/original-manifest.cbor":         {171, "1010ae00c786a6266348173c7760e0190be4cc280be1f71c8549f09727e4b183", "application/capsule.original-authoring-source-manifest+cbor;v=0"},
	"objects/executable-manifest.cbor":       {174, "295138062d0785785373b8c468fee75f77a28131d0974f30f69c4050425e9814", "application/capsule.executable-javascript-source-manifest+cbor;v=0"},
	"objects/transformation-record.cbor":     {659, "3ddf8b242afbddaa1856f79a24a46f7e8f9c674699ef6b1bd63215d49e512c39", "application/capsule.typescript-transformation-record+cbor;v=0"},
	"objects/transformation-record-set.cbor": {714, "5738283a5accdbd8b736af81982dc46068172ec502f5c43e4113fe7de10c76eb", "application/capsule.typescript-transformation-record-set+cbor;v=0"},
	"objects/plan-source-bindings.cbor":      {150, "deb9342fc8c2c0ea18ff280aa3c409246d7c96fab9dc46dddb54862640e4cc28", "application/capsule.typescript-plan-source-bindings+cbor;v=0"},
}

var expectedMutations = map[string]string{
	"original-byte":                       "BYTE_BINDING",
	"emitted-byte":                        "BYTE_BINDING",
	"path-order":                          "ORDERING",
	"profile-digest":                      "TRANSFORMER_BINDING",
	"profile-node":                        "TRANSFORMER_BINDING",
	"options-bytes":                       "OPTIONS_BINDING",
	"options-digest":                      "OPTIONS_BINDING",
	"input-media-type":                    "MEDIA_TYPE",
	"output-media-type":                   "MEDIA_TYPE",
	"source-map-present":                  "DISPOSITION",
	"source-url-present":                  "DISPOSITION",
	"diagnostic-count-one":                "DIAGNOSTIC",
	"cross-domain-original-as-executable": "DOMAIN",
	"cross-domain-record-set-as-original": "DOMAIN",
}

func VerifyFixtureFS(fsys fs.FS) (*VerifiedFixtureSet, error) {
	manifestBytes, err := readBounded(fsys, "manifest.json", ManifestMaxBytes)
	if err != nil {
		return nil, fmt.Errorf("read manifest: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(manifestBytes))
	decoder.DisallowUnknownFields()
	var value manifest
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, errors.New("manifest has trailing JSON")
	}
	if value.ManifestVersion != "capsule.typescript-approved-byte-conformance/v0" ||
		value.Status != "passive-unwired-fixture-only" ||
		value.FixtureCount != len(expectedAnswers) || value.MutationCount != len(expectedMutations) {
		return nil, errors.New("manifest identity or count mismatch")
	}
	if value.Caps != (caps{32, 262_144, 1_048_576, 262_144, 1_048_576}) {
		return nil, errors.New("manifest cap mismatch")
	}
	if len(value.KnownAnswers) != len(expectedAnswers) || len(value.Mutations) != len(expectedMutations) {
		return nil, errors.New("manifest retained entry count mismatch")
	}

	retained := make(map[string][]byte, len(expectedAnswers))
	seen := make(map[string]bool, len(expectedAnswers))
	for _, answer := range value.KnownAnswers {
		expected, ok := expectedAnswers[answer.Path]
		if !ok || seen[answer.Path] || answer.ByteLength != expected.length ||
			answer.SHA256 != expected.digest || answer.MediaType != expected.mediaType {
			return nil, fmt.Errorf("known answer mismatch for %q", answer.Path)
		}
		seen[answer.Path] = true
		limit := ObjectMaxBytes
		if len(answer.Path) >= 6 && answer.Path[:6] == "files/" {
			limit = SourceMaxBytes
		}
		fixture, err := readBounded(fsys, answer.Path, limit)
		if err != nil {
			return nil, fmt.Errorf("read fixture %q: %w", answer.Path, err)
		}
		digest := sha256.Sum256(fixture)
		if len(fixture) != expected.length || hex.EncodeToString(digest[:]) != expected.digest {
			return nil, fmt.Errorf("fixture byte mismatch for %q", answer.Path)
		}
		retained[answer.Path] = bytes.Clone(fixture)
	}
	mutationSeen := make(map[string]bool, len(expectedMutations))
	for _, item := range value.Mutations {
		expected, ok := expectedMutations[item.ID]
		if !ok || mutationSeen[item.ID] || item.Decision != "reject" || item.Classification != expected {
			return nil, fmt.Errorf("mutation mismatch for %q", item.ID)
		}
		mutationSeen[item.ID] = true
	}

	return &VerifiedFixtureSet{
		bytes: retained,
		bindings: PlanSourceBindings{
			OriginalManifest:        OriginalManifestDigest(mustDigest(expectedAnswers["objects/original-manifest.cbor"].digest)),
			ExecutableManifest:      ExecutableManifestDigest(mustDigest(expectedAnswers["objects/executable-manifest.cbor"].digest)),
			TransformationRecordSet: TransformationRecordSetDigest(mustDigest(expectedAnswers["objects/transformation-record-set.cbor"].digest)),
		},
	}, nil
}

func readBounded(fsys fs.FS, path string, maximum int) ([]byte, error) {
	file, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	value, readErr := io.ReadAll(io.LimitReader(file, int64(maximum)+1))
	closeErr := file.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if len(value) > maximum {
		return nil, errors.New("fixture exceeds bounded read")
	}
	return value, nil
}

func mustDigest(value string) [32]byte {
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != 32 {
		panic("invalid retained approved-byte digest")
	}
	var result [32]byte
	copy(result[:], decoded)
	return result
}
