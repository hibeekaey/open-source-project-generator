package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidatePath validates that a file path is safe and within expected boundaries
func ValidatePath(path string, allowedBasePaths ...string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Clean the path to resolve any .. or . components
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected: %s", path)
	}

	// If allowed base paths are specified, ensure the path is within one of them
	if len(allowedBasePaths) > 0 {
		allowed := false
		for _, basePath := range allowedBasePaths {
			cleanBasePath := filepath.Clean(basePath)
			if strings.HasPrefix(cleanPath, cleanBasePath) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("path not within allowed directories: %s", path)
		}
	}

	return nil
}

// SafeReadFile reads a file with path validation
func SafeReadFile(path string, allowedBasePaths ...string) ([]byte, error) {
	if err := ValidatePath(path, allowedBasePaths...); err != nil {
		return nil, err
	}
	return os.ReadFile(path) // #nosec G304 - path is validated above
}

// SafeWriteFile writes a file with secure permissions and path validation
func SafeWriteFile(path string, data []byte, allowedBasePaths ...string) error {
	if err := ValidatePath(path, allowedBasePaths...); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600) // More secure permissions
}

// SafeMkdirAll creates directories with secure permissions and path validation
func SafeMkdirAll(path string, allowedBasePaths ...string) error {
	if err := ValidatePath(path, allowedBasePaths...); err != nil {
		return err
	}
	return os.MkdirAll(path, 0750) // More secure permissions
}

// SafeOpenFile opens a file with path validation
func SafeOpenFile(path string, flag int, perm os.FileMode, allowedBasePaths ...string) (*os.File, error) {
	if err := ValidatePath(path, allowedBasePaths...); err != nil {
		return nil, err
	}

	// Ensure secure permissions for new files
	if flag&os.O_CREATE != 0 && perm > 0600 {
		perm = 0600
	}

	return os.OpenFile(path, flag, perm) // #nosec G304 - path is validated above
}

// SafeOpen opens a file for reading with path validation
func SafeOpen(path string, allowedBasePaths ...string) (*os.File, error) {
	if err := ValidatePath(path, allowedBasePaths...); err != nil {
		return nil, err
	}
	return os.Open(path) // #nosec G304 - path is validated above
}

// SafeCreate creates a file with secure permissions and path validation
func SafeCreate(path string, allowedBasePaths ...string) (*os.File, error) {
	if err := ValidatePath(path, allowedBasePaths...); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600) // #nosec G304 - path is validated above
}
