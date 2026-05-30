package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"goboxd/internal/config"
	"goboxd/internal/engine"
)

// setupMockEnv sets up a minimal registry and engine for isolated router testing.
func setupMockEnv() (map[string]config.Language, *engine.Engine) {
	registry := map[string]config.Language{
		"cpp": {
			ID:             "cpp",
			Name:           "C++",
			SourceFilename: "solution.cpp",
			Run: config.Command{
				Cmd: "./solution",
			},
		},
	}
	// Start an engine with a single worker to handle queued test requests
	eng := engine.NewEngine(1)
	return registry, eng
}

func TestGetHealthz(t *testing.T) {
	registry, eng := setupMockEnv()
	router := NewRouter(registry, eng)

	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	expectedBody := `{"status":"ok"}`
	if strings.TrimSpace(rr.Body.String()) != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rr.Body.String())
	}
}

func TestPostRun_UnknownLanguage(t *testing.T) {
	registry, eng := setupMockEnv()
	router := NewRouter(registry, eng)

	// Payload attempting to request an unconfigured language registry key
	badPayload := map[string]interface{}{
		"language": "rust",
		"source":   "fn main() {}",
	}
	bodyBytes, _ := json.Marshal(badPayload)

	req, err := http.NewRequest("POST", "/run", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Spec demands HTTP 400 for structural or missing configuration targets
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for unknown language, got %d", rr.Code)
	}
}

func TestPostRun_OversizedPayloadBoundary(t *testing.T) {
	registry, eng := setupMockEnv()
	router := NewRouter(registry, eng)

	// Security Hole #4 Check: Generate a massive code string exceeding the 256 KiB (262,144 bytes) limit
	hugeSource := strings.Repeat("A", 300*1024) 
	oversizedPayload := map[string]interface{}{
		"language": "cpp",
		"source":   hugeSource,
	}
	bodyBytes, _ := json.Marshal(oversizedPayload)

	req, err := http.NewRequest("POST", "/run", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// The router's http.MaxBytesReader must choke on the read loop and bubble up an error
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for oversized payload, got %d", rr.Code)
	}
}

func TestPostRun_ValidFlow(t *testing.T) {
	registry, eng := setupMockEnv()
	router := NewRouter(registry, eng)

	validPayload := map[string]interface{}{
		"language": "cpp",
		"source":   "#include <iostream>\nint main() { return 0; }",
	}
	bodyBytes, _ := json.Marshal(validPayload)

	req, err := http.NewRequest("POST", "/run", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	
	// Process request through router and engine channel pipelines
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 for functional run, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp engine.RunResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal JSON response: %v", err)
	}

	// Verify top-level payload structure returned matches internal engine state resolution
	if resp.Status == "" {
		t.Error("expected response status key to be initialized, got empty string")
	}
}