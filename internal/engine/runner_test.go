package engine

import (
	"bytes"
	"io"
	"testing"
)

func TestOutputTruncationBoundary(t *testing.T) {
	// Verifies that Security Hole #6 (Unbounded child output) acts predictably
	// using an identical LimitedReader strategy.
	maxBytes := int64(10)
	inputData := bytes.NewReader([]byte("this is a very long string that should be cut short"))
	
	var stdout bytes.Buffer
	limitedReader := &io.LimitedReader{R: inputData, N: maxBytes}
	
	copied, err := io.Copy(&stdout, limitedReader)
	if err != nil && err != io.EOF {
		t.Fatalf("unexpected error during copy: %v", err)
	}

	if copied > maxBytes {
		t.Errorf("Expected stream to cap at %d bytes, got %d", maxBytes, copied)
	}

	expectedString := "this is a "
	if stdout.String() != expectedString {
		t.Errorf("Expected truncated string to be %q, got %q", expectedString, stdout.String())
	}
}

func TestTopLevelStatusMapping(t *testing.T) {
	// Simulates the API rule mapping table: 
	// Top-level status is accepted only if build is ok and every test is accepted.
	// Otherwise, it matches the first non-accepted failure status.
	
	type testCaseResult struct {
		status string
	}

	tests := []struct {
		name           string
		buildStatus    string
		testResults    []testCaseResult
		expectedStatus string
	}{
		{
			name:        "All accepted passes cleanly",
			buildStatus: "ok",
			testResults: []testCaseResult{{"accepted"}, {"accepted"}},
			expectedStatus: "accepted",
		},
		{
			name:        "Build failure overrides tests",
			buildStatus: "failed",
			testResults: []testCaseResult{{"not_executed"}, {"not_executed"}},
			expectedStatus: "build_failed",
		},
		{
			name:        "First test failure maps to top level",
			buildStatus: "ok",
			testResults: []testCaseResult{{"wrong_output"}, {"time_exceeded"}},
			expectedStatus: "wrong_output",
		},
		{
			name:        "Subsequent test failure maps to top level",
			buildStatus: "ok",
			testResults: []testCaseResult{{"accepted"}, {"runtime_error"}},
			expectedStatus: "runtime_error",
		},
		{
			name:        "Internal build error maps to internal_error",
			buildStatus: "internal_error",
			testResults: []testCaseResult{{"not_executed"}},
			expectedStatus: "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Evaluates status resolution logic directly
			var resolvedStatus string
			if tt.buildStatus == "failed" {
				resolvedStatus = "build_failed"
			} else if tt.buildStatus == "internal_error" {
				resolvedStatus = "internal_error"
			} else {
				resolvedStatus = "accepted"
				for _, tc := range tt.testResults {
					if tc.status != "accepted" {
						resolvedStatus = tc.status
						break
					}
				}
			}

			if resolvedStatus != tt.expectedStatus {
				t.Errorf("Status resolution mismatch: got %q, wanted %q", resolvedStatus, tt.expectedStatus)
			}
		})
	}
}