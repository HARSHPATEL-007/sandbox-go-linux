package security

import (
	"fmt"
	"path/filepath"
	"strings"
)

// IsValidFilename prevents path traversal (Security Hole #1)
// Ensures the filename is a single component, no slashes, no leading dots.
func IsValidFilename(name string) error {
	if name == "" || len(name) > 255 {
		return fmt.Errorf("invalid filename length")
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("path separators not allowed")
	}
	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("hidden files not allowed")
	}
	if filepath.Base(name) != name {
		return fmt.Errorf("must be a single path component")
	}
	return nil
}

// AreFlagsAllowed prevents compiler-flag injection (Security Hole #3)
func AreFlagsAllowed(requested []string, allowlist []string) error {
	allowedMap := make(map[string]bool)
	for _, f := range allowlist {
		allowedMap[f] = true
	}

	for _, flag := range requested {
		if !allowedMap[flag] {
			return fmt.Errorf("flag not in allowlist: %s", flag)
		}
	}
	return nil
}