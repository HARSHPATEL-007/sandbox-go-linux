package security

import (
	"testing"
)

func TestIsValidFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"Valid simple filename", "solution.cpp", false},
		{"Valid with numbers", "prog123.py", false},
		{"Empty filename", "", true},
		{"Too long filename", string(make([]byte, 256)), true},
		{"Path traversal relative", "../../etc/passwd", true},
		{"Path traversal nested", "app/../../secret", true},
		{"Absolute path Unix", "/etc/passwd", true},
		{"Absolute path Windows", "C:\\Windows\\system32", true},
		{"Hidden file blocked", ".env", true},
		{"Single dot blocked", ".", true},
		{"Double dot blocked", "..", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidFilename(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidFilename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAreFlagsAllowed(t *testing.T) {
	allowlist := []string{"-O0", "-O1", "-O2", "-O3", "-Wall", "-std=c++17"}

	tests := []struct {
		name      string
		requested []string
		wantErr   bool
	}{
		{"Empty flags allowed", []string{}, false},
		{"Valid optimization flags", []string{"-O2", "-Wall"}, false},
		{"Exact standard flag matching", []string{"-std=c++17"}, false},
		{"Malicious plugin injection", []string{"-fplugin=/tmp/malicious.so"}, true},
		{"Malicious linker flag", []string{"-Wl,-Bstatic"}, true},
		{"Code execution flag mix", []string{"-O3", "-x", "c"}, true},
		{"Partial match exploit attempt", []string{"-Wall-extra"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AreFlagsAllowed(tt.requested, allowlist)
			if (err != nil) != tt.wantErr {
				t.Errorf("AreFlagsAllowed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}