package models

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityOperationError(t *testing.T) {
	t.Run("create security error with all fields", func(t *testing.T) {
		err := NewSecurityError(
			SecurityErrorType,
			SecuritySeverityHigh,
			"test-component",
			"test-operation",
			"test security error",
			nil,
		)

		assert.NotNil(t, err)
		assert.Equal(t, SecuritySeverityHigh, err.Severity)
		assert.Equal(t, "test-component", err.Component)
		assert.Equal(t, "test-operation", err.Operation)
		assert.Equal(t, "test security error", err.Error())
		assert.Equal(t, SecurityErrorType, err.Type)
	})

	t.Run("create security error with cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := NewSecurityError(
			CryptographicErrorType,
			SecuritySeverityCritical,
			"crypto",
			"random_gen",
			"crypto operation failed",
			cause,
		)

		assert.NotNil(t, err)
		assert.Equal(t, cause, err.Unwrap())
		assert.Contains(t, err.Error(), "crypto operation failed")
		assert.Contains(t, err.Error(), "underlying error")
	})

	t.Run("add remediation to security error", func(t *testing.T) {
		err := NewSecurityError(
			FileSecurityErrorType,
			SecuritySeverityMedium,
			"filesystem",
			"write",
			"file operation failed",
			nil,
		).WithRemediation("Check file permissions")

		assert.Equal(t, "Check file permissions", err.Remediation)
	})

	t.Run("check if error is critical", func(t *testing.T) {
		criticalErr := NewSecurityError(
			SecurityErrorType,
			SecuritySeverityCritical,
			"test",
			"test",
			"critical error",
			nil,
		)

		nonCriticalErr := NewSecurityError(
			SecurityErrorType,
			SecuritySeverityLow,
			"test",
			"test",
			"low error",
			nil,
		)

		assert.True(t, criticalErr.IsCritical())
		assert.False(t, nonCriticalErr.IsCritical())
	})
}

func TestSecuritySeverity(t *testing.T) {
	tests := []struct {
		severity SecuritySeverity
		expected string
	}{
		{SecuritySeverityLow, "low"},
		{SecuritySeverityMedium, "medium"},
		{SecuritySeverityHigh, "high"},
		{SecuritySeverityCritical, "critical"},
		{SecuritySeverity(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("severity_%s", tt.expected), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.severity.String())
		})
	}
}

func TestPredefinedSecurityErrors(t *testing.T) {
	t.Run("all predefined errors are security errors", func(t *testing.T) {
		predefinedErrors := []*SecurityOperationError{
			ErrInsufficientEntropy,
			ErrCryptographicFailure,
			ErrInvalidPath,
			ErrTempFileCreation,
			ErrAtomicWrite,
			ErrInsecurePermissions,
			ErrDangerousDirectory,
			ErrSecurityViolation,
			ErrUnauthorizedOperation,
		}

		for _, err := range predefinedErrors {
			assert.NotNil(t, err, "Predefined error should not be nil")
			assert.NotEmpty(t, err.Error(), "Error message should not be empty")
			assert.NotEmpty(t, err.Component, "Component should be specified")
			assert.NotEmpty(t, err.Operation, "Operation should be specified")
			assert.NotEmpty(t, err.Remediation, "Remediation should be provided")
			assert.True(t, IsSecurityError(err), "Should be identified as security error")
		}
	})

	t.Run("critical errors are properly marked", func(t *testing.T) {
		assert.True(t, ErrInsufficientEntropy.IsCritical())
		assert.False(t, ErrTempFileCreation.IsCritical())
	})

	t.Run("error types are correctly assigned", func(t *testing.T) {
		assert.Equal(t, CryptographicErrorType, ErrInsufficientEntropy.Type)
		assert.Equal(t, PathValidationErrorType, ErrInvalidPath.Type)
		assert.Equal(t, FileSecurityErrorType, ErrTempFileCreation.Type)
	})
}

func TestIsSecurityError(t *testing.T) {
	t.Run("identifies security errors correctly", func(t *testing.T) {
		secErr := NewSecurityError(
			SecurityErrorType,
			SecuritySeverityHigh,
			"test",
			"test",
			"security error",
			nil,
		)

		regularErr := NewGeneratorError(ValidationErrorType, "validation error", nil)
		stdErr := errors.New("standard error")

		assert.True(t, IsSecurityError(secErr))
		assert.False(t, IsSecurityError(regularErr))
		assert.False(t, IsSecurityError(stdErr))
		assert.False(t, IsSecurityError(nil))
	})

	t.Run("identifies wrapped security errors", func(t *testing.T) {
		baseErr := errors.New("base error")
		wrappedErr := WrapWithSecurityContext(
			baseErr,
			SecuritySeverityMedium,
			"wrapper",
			"wrap",
		)

		assert.True(t, IsSecurityError(wrappedErr))
		assert.Equal(t, baseErr, wrappedErr.Unwrap())
	})
}

func TestGetSecuritySeverity(t *testing.T) {
	t.Run("extracts severity from security errors", func(t *testing.T) {
		highErr := NewSecurityError(
			SecurityErrorType,
			SecuritySeverityHigh,
			"test",
			"test",
			"high severity error",
			nil,
		)

		assert.Equal(t, SecuritySeverityHigh, GetSecuritySeverity(highErr))
	})

	t.Run("returns low severity for non-security errors", func(t *testing.T) {
		regularErr := errors.New("regular error")
		assert.Equal(t, SecuritySeverityLow, GetSecuritySeverity(regularErr))
		assert.Equal(t, SecuritySeverityLow, GetSecuritySeverity(nil))
	})
}

func TestWrapWithSecurityContext(t *testing.T) {
	t.Run("wraps error with security context", func(t *testing.T) {
		baseErr := errors.New("base error")
		wrappedErr := WrapWithSecurityContext(
			baseErr,
			SecuritySeverityHigh,
			"test-component",
			"test-operation",
		)

		require.NotNil(t, wrappedErr)
		assert.Equal(t, SecuritySeverityHigh, wrappedErr.Severity)
		assert.Equal(t, "test-component", wrappedErr.Component)
		assert.Equal(t, "test-operation", wrappedErr.Operation)
		assert.Equal(t, baseErr, wrappedErr.Unwrap())
		assert.Contains(t, wrappedErr.Error(), "base error")
	})

	t.Run("returns nil for nil error", func(t *testing.T) {
		wrappedErr := WrapWithSecurityContext(
			nil,
			SecuritySeverityHigh,
			"test",
			"test",
		)

		assert.Nil(t, wrappedErr)
	})
}

func TestMessageSanitization(t *testing.T) {
	t.Run("sanitizes long messages", func(t *testing.T) {
		longMessage := strings.Repeat("a", 250)
		err := NewSecurityError(
			SecurityErrorType,
			SecuritySeverityLow,
			"test",
			"test",
			longMessage,
			nil,
		)

		errorMsg := err.Error()
		assert.True(t, len(errorMsg) <= 250, "Message should be truncated")
		assert.Contains(t, errorMsg, "[message truncated for security]")
	})

	t.Run("preserves short messages", func(t *testing.T) {
		shortMessage := "short message"
		err := NewSecurityError(
			SecurityErrorType,
			SecuritySeverityLow,
			"test",
			"test",
			shortMessage,
			nil,
		)

		assert.Equal(t, shortMessage, err.Error())
	})
}

func TestSecurityErrorContext(t *testing.T) {
	t.Run("adds context to security error", func(t *testing.T) {
		err := NewSecurityError(
			SecurityErrorType,
			SecuritySeverityMedium,
			"test",
			"test",
			"test error",
			nil,
		).WithContext("file_path", "/safe/path").WithContext("user_id", "user123")

		assert.Equal(t, "/safe/path", err.Context["file_path"])
		assert.Equal(t, "user123", err.Context["user_id"])
	})
}

func TestEntropyFailureScenarios(t *testing.T) {
	t.Run("entropy failure error properties", func(t *testing.T) {
		err := ErrInsufficientEntropy

		assert.Equal(t, CryptographicErrorType, err.Type)
		assert.Equal(t, SecuritySeverityCritical, err.Severity)
		assert.Equal(t, "crypto", err.Component)
		assert.Equal(t, "random_generation", err.Operation)
		assert.Contains(t, err.Remediation, "crypto/rand")
		assert.True(t, err.IsCritical())
	})

	t.Run("wrapping entropy failure with context", func(t *testing.T) {
		baseErr := errors.New("entropy source unavailable")
		wrappedErr := WrapWithSecurityContext(
			baseErr,
			SecuritySeverityCritical,
			"random_generator",
			"generate_secure_bytes",
		)

		assert.Equal(t, SecuritySeverityCritical, wrappedErr.Severity)
		assert.Contains(t, wrappedErr.Error(), "entropy source unavailable")
	})
}

func TestPathValidationErrorScenarios(t *testing.T) {
	t.Run("path validation error properties", func(t *testing.T) {
		err := ErrInvalidPath

		assert.Equal(t, PathValidationErrorType, err.Type)
		assert.Equal(t, SecuritySeverityHigh, err.Severity)
		assert.Equal(t, "filesystem", err.Component)
		assert.Equal(t, "path_validation", err.Operation)
		assert.Contains(t, err.Remediation, "validated")
	})

	t.Run("dangerous directory error", func(t *testing.T) {
		err := ErrDangerousDirectory

		assert.Equal(t, PathValidationErrorType, err.Type)
		assert.Equal(t, SecuritySeverityHigh, err.Severity)
		assert.Contains(t, err.Remediation, "trusted directory")
	})
}

func TestFileSecurityErrorScenarios(t *testing.T) {
	t.Run("atomic write failure", func(t *testing.T) {
		err := ErrAtomicWrite

		assert.Equal(t, FileSecurityErrorType, err.Type)
		assert.Equal(t, SecuritySeverityMedium, err.Severity)
		assert.Equal(t, "filesystem", err.Component)
		assert.Equal(t, "atomic_write", err.Operation)
	})

	t.Run("insecure permissions error", func(t *testing.T) {
		err := ErrInsecurePermissions

		assert.Equal(t, FileSecurityErrorType, err.Type)
		assert.Equal(t, SecuritySeverityHigh, err.Severity)
		assert.Contains(t, err.Remediation, "0600")
	})
}

func TestSecurityViolationScenarios(t *testing.T) {
	t.Run("general security violation", func(t *testing.T) {
		err := ErrSecurityViolation

		assert.Equal(t, SecurityErrorType, err.Type)
		assert.Equal(t, SecuritySeverityHigh, err.Severity)
		assert.Equal(t, "security", err.Component)
		assert.Equal(t, "validation", err.Operation)
	})

	t.Run("unauthorized operation", func(t *testing.T) {
		err := ErrUnauthorizedOperation

		assert.Equal(t, SecurityErrorType, err.Type)
		assert.Equal(t, SecuritySeverityHigh, err.Severity)
		assert.Equal(t, "security", err.Component)
		assert.Equal(t, "authorization", err.Operation)
	})
}
