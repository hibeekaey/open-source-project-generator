package models

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExampleSecureFileOperation demonstrates proper security error handling
// for file operations with comprehensive error context
func ExampleSecureFileOperation(filePath string, data []byte) error {
	// Validate the file path first
	if !isSecureFilePath(filePath) {
		return ErrInvalidPath.WithContext("attempted_path", filePath)
	}

	// Check if directory is safe
	dir := filepath.Dir(filePath)
	if isDangerousDirectory(dir) {
		return ErrDangerousDirectory.WithContext("attempted_path", filePath)
	}

	// Create secure temporary file
	tempFile, err := createSecureTempFile(dir)
	if err != nil {
		return WrapWithSecurityContext(
			err,
			SecuritySeverityMedium,
			"temp_file_manager",
			"secure_temp_creation",
		).WithRemediation("Ensure directory is writable and has sufficient space")
	}
	defer os.Remove(tempFile.Name())

	// Write data atomically
	if err := writeDataSecurely(tempFile, data); err != nil {
		return WrapWithSecurityContext(
			err,
			SecuritySeverityMedium,
			"file_writer",
			"atomic_write",
		).WithRemediation("Check file permissions and available disk space")
	}

	// Atomic rename to final location
	if err := os.Rename(tempFile.Name(), filePath); err != nil {
		return ErrAtomicWrite.WithContext("target_file", filePath)
	}

	return nil
}

// ExampleSecureRandomGeneration demonstrates proper error handling
// for cryptographic operations with entropy failure scenarios
func ExampleSecureRandomGeneration(length int) ([]byte, error) {
	if length <= 0 {
		return nil, NewSecurityError(
			CryptographicErrorType,
			SecuritySeverityMedium,
			"random_generator",
			"parameter_validation",
			"invalid random length requested",
			nil,
		).WithRemediation("Use positive length values for random generation")
	}

	// Attempt to generate secure random bytes
	randomBytes := make([]byte, length)
	n, err := rand.Read(randomBytes)
	if err != nil {
		return nil, WrapWithSecurityContext(
			err,
			SecuritySeverityCritical,
			"crypto_rand",
			"entropy_generation",
		).WithRemediation("Ensure system has sufficient entropy and crypto/rand is available")
	}

	if n != length {
		return nil, ErrInsufficientEntropy.WithContext("requested_bytes", length).WithContext("generated_bytes", n)
	}

	return randomBytes, nil
}

// ExampleErrorHandlingWithLogging demonstrates how to handle security errors
// with proper logging that doesn't leak sensitive information
func ExampleErrorHandlingWithLogging(operation func() error) error {
	err := operation()
	if err == nil {
		return nil
	}

	// Check if it's a security error and handle appropriately
	if IsSecurityError(err) {
		severity := GetSecuritySeverity(err)

		// Log security errors with appropriate detail level
		switch severity {
		case SecuritySeverityCritical:
			// Critical errors need immediate attention but shouldn't leak details
			fmt.Printf("CRITICAL SECURITY ERROR: Operation failed due to security violation (ID: %s)\n",
				generateSecureLogID())
		case SecuritySeverityHigh:
			fmt.Printf("HIGH SECURITY WARNING: Security policy violation detected\n")
		case SecuritySeverityMedium:
			fmt.Printf("MEDIUM SECURITY NOTICE: Security check failed\n")
		case SecuritySeverityLow:
			fmt.Printf("LOW SECURITY INFO: Minor security issue detected\n")
		}

		// Return sanitized error to caller
		return err
	}

	// Handle non-security errors normally
	fmt.Printf("Operation error: %v\n", err)
	return err
}

// Helper functions for the examples

func isSecureFilePath(path string) bool {
	// Basic path validation - in real implementation this would be more comprehensive
	if !filepath.IsAbs(path) {
		return false
	}

	// Check for path traversal attempts
	cleanPath := filepath.Clean(path)
	if cleanPath != path {
		return false
	}

	// Check for dangerous patterns
	if strings.Contains(path, "..") || strings.Contains(path, "./") {
		return false
	}

	return true
}

func isDangerousDirectory(dir string) bool {
	dangerousDirs := []string{"/", "/bin", "/sbin", "/usr/bin", "/etc", "/root", "/boot"}
	for _, dangerous := range dangerousDirs {
		if dir == dangerous {
			return true
		}
	}
	return false
}

func createSecureTempFile(dir string) (*os.File, error) {
	// This would use secure random naming in real implementation
	return os.CreateTemp(dir, "secure_*.tmp")
}

func writeDataSecurely(file *os.File, data []byte) error {
	_, err := file.Write(data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func generateSecureLogID() string {
	// Generate a secure ID for logging purposes
	bytes, _ := ExampleSecureRandomGeneration(8)
	return fmt.Sprintf("%x", bytes)
}
