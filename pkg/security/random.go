// Package security provides cryptographically secure utilities for random generation
// and file operations. This package addresses common security vulnerabilities by:
//
// 1. Using crypto/rand instead of predictable sources like math/rand or timestamps
// 2. Implementing atomic file operations to prevent race conditions
// 3. Providing path validation to prevent directory traversal attacks
// 4. Ensuring secure temporary file creation with unpredictable names
//
// SECURITY RATIONALE:
// The primary motivation for this package is to replace insecure patterns found
// throughout the codebase, particularly:
// - Predictable temporary file names using time.Now().UnixNano()
// - Timestamp-based ID generation for security-sensitive operations
// - Non-atomic file operations vulnerable to race conditions
//
// All random generation in this package uses crypto/rand, which provides
// cryptographically secure pseudorandom numbers suitable for security-sensitive
// applications. Never use math/rand or timestamp-based generation for security.
package security

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// SecureRandom interface defines methods for cryptographically secure random generation
type SecureRandom interface {
	// GenerateRandomSuffix creates a cryptographically secure random suffix
	GenerateRandomSuffix(length int) (string, error)

	// GenerateSecureID creates a secure identifier for audit trails
	GenerateSecureID(prefix string) (string, error)

	// GenerateBytes creates random bytes for cryptographic operations
	GenerateBytes(length int) ([]byte, error)

	// GenerateHexString creates a random hex string of specified length
	GenerateHexString(length int) (string, error)

	// GenerateBase64String creates a random base64 string of specified length
	GenerateBase64String(length int) (string, error)

	// GenerateAlphanumeric creates a random alphanumeric string of specified length
	GenerateAlphanumeric(length int) (string, error)
}

// DefaultSecureRandom is the default implementation of SecureRandom
type DefaultSecureRandom struct {
	// DefaultSuffixLength is the default length for random suffixes
	DefaultSuffixLength int
	// IDFormat specifies the format for secure IDs (hex, base64, alphanumeric)
	IDFormat string
}

// NewSecureRandom creates a new instance of DefaultSecureRandom with default settings
func NewSecureRandom() *DefaultSecureRandom {
	return &DefaultSecureRandom{
		DefaultSuffixLength: 16,
		IDFormat:            "hex",
	}
}

// NewSecureRandomWithConfig creates a new instance with custom configuration
func NewSecureRandomWithConfig(suffixLength int, idFormat string) *DefaultSecureRandom {
	return &DefaultSecureRandom{
		DefaultSuffixLength: suffixLength,
		IDFormat:            idFormat,
	}
}

// GenerateBytes creates random bytes for cryptographic operations
//
// SECURITY RATIONALE:
// This function uses crypto/rand.Read() which provides cryptographically secure
// random bytes suitable for:
// - Cryptographic keys and initialization vectors
// - Session tokens and CSRF tokens
// - Secure temporary file suffixes
// - Any security-sensitive random data
//
// crypto/rand.Read() uses the operating system's cryptographically secure
// random number generator (/dev/urandom on Unix, CryptGenRandom on Windows).
// This is fundamentally different from math/rand which is deterministic and
// predictable given the seed.
//
// ERROR HANDLING:
// Entropy failures can occur in virtualized environments or systems with
// insufficient entropy. Always handle errors from this function - a failure
// to generate secure random data should cause the operation to fail rather
// than fall back to predictable alternatives.
func (sr *DefaultSecureRandom) GenerateBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, fmt.Errorf("length must be positive, got %d", length)
	}

	bytes := make([]byte, length)
	// crypto/rand.Read() is the ONLY acceptable source for security-sensitive randomness
	_, err := rand.Read(bytes)
	if err != nil {
		// Entropy failure is a serious security issue - never ignore this error
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return bytes, nil
}

// GenerateRandomSuffix creates a cryptographically secure random suffix
func (sr *DefaultSecureRandom) GenerateRandomSuffix(length int) (string, error) {
	if length <= 0 {
		length = sr.DefaultSuffixLength
	}

	switch sr.IDFormat {
	case "base64":
		return sr.GenerateBase64String(length)
	case "alphanumeric":
		return sr.GenerateAlphanumeric(length)
	default: // hex
		return sr.GenerateHexString(length)
	}
}

// GenerateSecureID creates a secure identifier for audit trails
func (sr *DefaultSecureRandom) GenerateSecureID(prefix string) (string, error) {
	// Generate a secure random suffix with default length
	suffix, err := sr.GenerateRandomSuffix(sr.DefaultSuffixLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate secure ID: %w", err)
	}

	if prefix == "" {
		return suffix, nil
	}

	return fmt.Sprintf("%s_%s", prefix, suffix), nil
}

// GenerateHexString creates a random hex string of specified length
func (sr *DefaultSecureRandom) GenerateHexString(length int) (string, error) {
	// For hex string, we need length/2 bytes (each byte produces 2 hex chars)
	byteLength := (length + 1) / 2
	bytes, err := sr.GenerateBytes(byteLength)
	if err != nil {
		return "", err
	}

	hexString := hex.EncodeToString(bytes)

	// Truncate to exact length if needed
	if len(hexString) > length {
		hexString = hexString[:length]
	}

	return hexString, nil
}

// GenerateBase64String creates a random base64 string of specified length
func (sr *DefaultSecureRandom) GenerateBase64String(length int) (string, error) {
	// For base64, we need approximately length*3/4 bytes
	byteLength := (length*3 + 3) / 4
	bytes, err := sr.GenerateBytes(byteLength)
	if err != nil {
		return "", err
	}

	base64String := base64.RawURLEncoding.EncodeToString(bytes)

	// Remove padding and truncate to exact length
	base64String = strings.TrimRight(base64String, "=")
	if len(base64String) > length {
		base64String = base64String[:length]
	}

	return base64String, nil
}

// GenerateAlphanumeric creates a random alphanumeric string of specified length
//
// SECURITY RATIONALE:
// This function generates cryptographically secure alphanumeric strings using
// crypto/rand.Int() with uniform distribution. Each character is independently
// selected from the charset using secure random selection.
//
// ENTROPY CONSIDERATIONS:
// Alphanumeric charset (62 characters) provides ~5.95 bits of entropy per character.
// For comparison:
// - Hex (16 chars): 4 bits per character
// - Base64 (64 chars): 6 bits per character
// - Alphanumeric (62 chars): ~5.95 bits per character
//
// For high-security applications, consider using hex or base64 formats which
// provide better entropy density and are less prone to character confusion.
//
// UNIFORM DISTRIBUTION:
// Using crypto/rand.Int() ensures uniform distribution across the charset,
// preventing bias that could reduce effective entropy. Simple modulo operations
// can introduce bias and should be avoided for security applications.
func (sr *DefaultSecureRandom) GenerateAlphanumeric(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLen := big.NewInt(int64(len(charset)))

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		// Use crypto/rand.Int for uniform distribution - never use modulo operations
		// which can introduce bias and reduce effective entropy
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random alphanumeric character: %w", err)
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}

// Global instance for convenience
var defaultSecureRandom = NewSecureRandom()

// Convenience functions using the global instance

// GenerateRandomSuffix creates a cryptographically secure random suffix using the global instance
func GenerateRandomSuffix(length int) (string, error) {
	return defaultSecureRandom.GenerateRandomSuffix(length)
}

// GenerateSecureID creates a secure identifier using the global instance
func GenerateSecureID(prefix string) (string, error) {
	return defaultSecureRandom.GenerateSecureID(prefix)
}

// GenerateBytes creates random bytes using the global instance
func GenerateBytes(length int) ([]byte, error) {
	return defaultSecureRandom.GenerateBytes(length)
}

// GenerateHexString creates a random hex string using the global instance
func GenerateHexString(length int) (string, error) {
	return defaultSecureRandom.GenerateHexString(length)
}

// GenerateBase64String creates a random base64 string using the global instance
func GenerateBase64String(length int) (string, error) {
	return defaultSecureRandom.GenerateBase64String(length)
}

// GenerateAlphanumeric creates a random alphanumeric string using the global instance
func GenerateAlphanumeric(length int) (string, error) {
	return defaultSecureRandom.GenerateAlphanumeric(length)
}
