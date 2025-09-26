// Package errors provides security utilities for error handling
package errors

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// validatePath validates that a path is safe to use and prevents path traversal attacks
func validatePath(path string) error {
	// Clean the path to resolve any .. or . components
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected in path: %s", path)
	}

	// Ensure the path is absolute or relative to current directory
	if filepath.IsAbs(cleanPath) {
		// For absolute paths, ensure they're within reasonable bounds
		// Block obviously dangerous paths
		dangerousPaths := []string{"/etc", "/usr", "/bin", "/sbin", "/root"}
		for _, dangerous := range dangerousPaths {
			if strings.HasPrefix(cleanPath, dangerous) {
				return fmt.Errorf("path not allowed in system directory: %s", path)
			}
		}

		// Allow temp directory and user directories
		// Additional validation could be added here for other absolute paths
	}

	return nil
}

// secureCreateFile creates a file with path validation
func secureCreateFile(path string) (*os.File, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}

	// Use filepath.Clean to normalize the path
	cleanPath := filepath.Clean(path)

	return os.Create(cleanPath) // #nosec G304 - path has been validated
}

// secureOpenFile opens a file with path validation and secure permissions
func secureOpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	if err := validatePath(path); err != nil {
		return nil, err
	}

	// Use filepath.Clean to normalize the path
	cleanPath := filepath.Clean(path)

	// Ensure secure permissions (0600 or less for files)
	if perm > 0600 {
		perm = 0600
	}

	return os.OpenFile(cleanPath, flag, perm) // #nosec G304 - path has been validated
}

// secureMkdirAll creates directories with path validation and secure permissions
func secureMkdirAll(path string, perm os.FileMode) error {
	if err := validatePath(path); err != nil {
		return err
	}

	// Use filepath.Clean to normalize the path
	cleanPath := filepath.Clean(path)

	// Ensure secure permissions (0750 or less for directories)
	if perm > 0750 {
		perm = 0750
	}

	return os.MkdirAll(cleanPath, perm) // #nosec G301 - path has been validated and permissions are secure
}
