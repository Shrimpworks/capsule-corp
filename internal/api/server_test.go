package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	NewServer(Options{}).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected ok health status, got %q", body["status"])
	}

	if response.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("expected responses to disable caching")
	}
}

func TestRuntimeProfileIsNotAvailable(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/v1/runtimes", nil)
	response := httptest.NewRecorder()

	NewServer(Options{}).ServeHTTP(response, request)

	var body struct {
		Profiles []struct {
			Name      string `json:"name"`
			Available bool   `json:"available"`
		} `json:"profiles"`
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(body.Profiles) != 1 {
		t.Fatalf("expected one profile, got %d", len(body.Profiles))
	}

	if body.Profiles[0].Name != "bun-data@1" || body.Profiles[0].Available {
		t.Fatalf("draft runtime must not be advertised as available: %+v", body.Profiles[0])
	}
}

func TestDiagnosticHTTPHasNoSubmissionOrMutationRoute(t *testing.T) {
	for _, path := range []string{"/v1/jobs", "/v1/proposals", "/v1/submit", "/v1/register", "/v1/attempts"} {
		request := httptest.NewRequest(http.MethodPost, path, nil)
		response := httptest.NewRecorder()

		NewServer(Options{}).ServeHTTP(response, request)

		if response.Code == http.StatusOK || response.Code == http.StatusAccepted || response.Code == http.StatusCreated {
			t.Fatalf("diagnostic HTTP exposed mutation route %s with status %d", path, response.Code)
		}
	}
}
