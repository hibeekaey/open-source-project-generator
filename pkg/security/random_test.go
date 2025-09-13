package security

import (
	"encoding/hex"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestNewSecureRandom(t *testing.T) {
	sr := NewSecureRandom()

	if sr.DefaultSuffixLength != 16 {
		t.Errorf("Expected default suffix length 16, got %d", sr.DefaultSuffixLength)
	}

	if sr.IDFormat != "hex" {
		t.Errorf("Expected default ID format 'hex', got %s", sr.IDFormat)
	}
}

func TestNewSecureRandomWithConfig(t *testing.T) {
	sr := NewSecureRandomWithConfig(32, "base64")

	if sr.DefaultSuffixLength != 32 {
		t.Errorf("Expected suffix length 32, got %d", sr.DefaultSuffixLength)
	}

	if sr.IDFormat != "base64" {
		t.Errorf("Expected ID format 'base64', got %s", sr.IDFormat)
	}
}

func TestGenerateBytes(t *testing.T) {
	sr := NewSecureRandom()

	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{"Valid length", 16, false},
		{"Small length", 1, false},
		{"Large length", 1024, false},
		{"Zero length", 0, true},
		{"Negative length", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := sr.GenerateBytes(tt.length)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for length %d, got nil", tt.length)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(bytes) != tt.length {
				t.Errorf("Expected %d bytes, got %d", tt.length, len(bytes))
			}
		})
	}
}

func TestGenerateHexString(t *testing.T) {
	sr := NewSecureRandom()

	tests := []struct {
		name   string
		length int
	}{
		{"Short hex", 8},
		{"Medium hex", 16},
		{"Long hex", 32},
		{"Odd length", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hexStr, err := sr.GenerateHexString(tt.length)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(hexStr) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(hexStr))
			}

			// Verify it contains only hex characters (0-9, a-f, A-F)
			if !isValidHexChars(hexStr) {
				t.Errorf("Generated string contains non-hex characters: %s", hexStr)
			}
		})
	}
}

func TestGenerateBase64String(t *testing.T) {
	sr := NewSecureRandom()

	tests := []struct {
		name   string
		length int
	}{
		{"Short base64", 8},
		{"Medium base64", 16},
		{"Long base64", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b64Str, err := sr.GenerateBase64String(tt.length)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(b64Str) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(b64Str))
			}

			// Verify it's valid base64 URL-safe characters
			if !isValidBase64URLSafe(b64Str) {
				t.Errorf("Generated string is not valid base64 URL-safe: %s", b64Str)
			}
		})
	}
}

func TestGenerateAlphanumeric(t *testing.T) {
	sr := NewSecureRandom()

	tests := []struct {
		name   string
		length int
	}{
		{"Short alphanumeric", 8},
		{"Medium alphanumeric", 16},
		{"Long alphanumeric", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alphaStr, err := sr.GenerateAlphanumeric(tt.length)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(alphaStr) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(alphaStr))
			}

			// Verify it's valid alphanumeric
			if !isValidAlphanumeric(alphaStr) {
				t.Errorf("Generated string is not valid alphanumeric: %s", alphaStr)
			}
		})
	}
}

func TestGenerateRandomSuffix(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		length   int
		expected string
	}{
		{"Hex format", "hex", 16, "hex"},
		{"Base64 format", "base64", 16, "base64"},
		{"Alphanumeric format", "alphanumeric", 16, "alphanumeric"},
		{"Default format", "", 16, "hex"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := NewSecureRandomWithConfig(16, tt.format)
			suffix, err := sr.GenerateRandomSuffix(tt.length)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(suffix) != tt.length {
				t.Errorf("Expected length %d, got %d", tt.length, len(suffix))
			}

			// Verify format
			switch tt.expected {
			case "hex":
				if !isValidHex(suffix) {
					t.Errorf("Expected hex format, got: %s", suffix)
				}
			case "base64":
				if !isValidBase64URLSafe(suffix) {
					t.Errorf("Expected base64 format, got: %s", suffix)
				}
			case "alphanumeric":
				if !isValidAlphanumeric(suffix) {
					t.Errorf("Expected alphanumeric format, got: %s", suffix)
				}
			}
		})
	}
}

func TestGenerateSecureID(t *testing.T) {
	sr := NewSecureRandom()

	tests := []struct {
		name   string
		prefix string
	}{
		{"With prefix", "audit"},
		{"Empty prefix", ""},
		{"Long prefix", "very_long_prefix_name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := sr.GenerateSecureID(tt.prefix)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.prefix == "" {
				// Should just be the suffix
				if len(id) != sr.DefaultSuffixLength {
					t.Errorf("Expected length %d, got %d", sr.DefaultSuffixLength, len(id))
				}
			} else {
				// Should be prefix + "_" + suffix
				expectedMinLength := len(tt.prefix) + 1 + sr.DefaultSuffixLength
				if len(id) != expectedMinLength {
					t.Errorf("Expected length %d, got %d", expectedMinLength, len(id))
				}

				if !strings.HasPrefix(id, tt.prefix+"_") {
					t.Errorf("Expected ID to start with '%s_', got: %s", tt.prefix, id)
				}
			}
		})
	}
}

// Test randomness quality
func TestRandomnessQuality(t *testing.T) {
	sr := NewSecureRandom()

	// Generate multiple random strings and ensure they're different
	const numSamples = 100
	const stringLength = 16

	samples := make(map[string]bool)

	for i := 0; i < numSamples; i++ {
		randomStr, err := sr.GenerateHexString(stringLength)
		if err != nil {
			t.Errorf("Unexpected error generating random string: %v", err)
			return
		}

		if samples[randomStr] {
			t.Errorf("Duplicate random string generated: %s", randomStr)
			return
		}

		samples[randomStr] = true
	}

	if len(samples) != numSamples {
		t.Errorf("Expected %d unique samples, got %d", numSamples, len(samples))
	}
}

// Test global convenience functions
func TestGlobalFunctions(t *testing.T) {
	t.Run("GenerateRandomSuffix", func(t *testing.T) {
		suffix, err := GenerateRandomSuffix(16)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(suffix) != 16 {
			t.Errorf("Expected length 16, got %d", len(suffix))
		}
	})

	t.Run("GenerateSecureID", func(t *testing.T) {
		id, err := GenerateSecureID("test")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !strings.HasPrefix(id, "test_") {
			t.Errorf("Expected ID to start with 'test_', got: %s", id)
		}
	})

	t.Run("GenerateBytes", func(t *testing.T) {
		bytes, err := GenerateBytes(16)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(bytes) != 16 {
			t.Errorf("Expected 16 bytes, got %d", len(bytes))
		}
	})

	t.Run("GenerateHexString", func(t *testing.T) {
		hexStr, err := GenerateHexString(16)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(hexStr) != 16 || !isValidHexChars(hexStr) {
			t.Errorf("Invalid hex string: %s", hexStr)
		}
	})

	t.Run("GenerateBase64String", func(t *testing.T) {
		b64Str, err := GenerateBase64String(16)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(b64Str) != 16 || !isValidBase64URLSafe(b64Str) {
			t.Errorf("Invalid base64 string: %s", b64Str)
		}
	})

	t.Run("GenerateAlphanumeric", func(t *testing.T) {
		alphaStr, err := GenerateAlphanumeric(16)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(alphaStr) != 16 || !isValidAlphanumeric(alphaStr) {
			t.Errorf("Invalid alphanumeric string: %s", alphaStr)
		}
	})
}

// Benchmark tests for performance
func BenchmarkGenerateBytes(b *testing.B) {
	sr := NewSecureRandom()

	for i := 0; i < b.N; i++ {
		_, err := sr.GenerateBytes(16)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGenerateHexString(b *testing.B) {
	sr := NewSecureRandom()

	for i := 0; i < b.N; i++ {
		_, err := sr.GenerateHexString(16)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGenerateBase64String(b *testing.B) {
	sr := NewSecureRandom()

	for i := 0; i < b.N; i++ {
		_, err := sr.GenerateBase64String(16)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGenerateAlphanumeric(b *testing.B) {
	sr := NewSecureRandom()

	for i := 0; i < b.N; i++ {
		_, err := sr.GenerateAlphanumeric(16)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGenerateSecureID(b *testing.B) {
	sr := NewSecureRandom()

	for i := 0; i < b.N; i++ {
		_, err := sr.GenerateSecureID("audit")
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

// Test concurrent access
func TestConcurrentAccess(t *testing.T) {
	sr := NewSecureRandom()
	const numGoroutines = 100
	const numOperations = 10

	results := make(chan string, numGoroutines*numOperations)

	// Start multiple goroutines generating random strings
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numOperations; j++ {
				str, err := sr.GenerateHexString(16)
				if err != nil {
					t.Errorf("Concurrent generation error: %v", err)
					return
				}
				results <- str
			}
		}()
	}

	// Collect all results
	seen := make(map[string]bool)
	for i := 0; i < numGoroutines*numOperations; i++ {
		select {
		case result := <-results:
			if seen[result] {
				t.Errorf("Duplicate result in concurrent test: %s", result)
			}
			seen[result] = true
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent results")
		}
	}
}

// Helper functions for validation
func isValidHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

func isValidHexChars(s string) bool {
	// Check if string contains only hex characters (0-9, a-f, A-F)
	matched, _ := regexp.MatchString("^[0-9a-fA-F]*$", s)
	return matched
}

func isValidBase64URLSafe(s string) bool {
	// Base64 URL-safe characters: A-Z, a-z, 0-9, -, _
	matched, _ := regexp.MatchString("^[A-Za-z0-9_-]*$", s)
	return matched
}

func isValidAlphanumeric(s string) bool {
	// Alphanumeric characters: A-Z, a-z, 0-9
	matched, _ := regexp.MatchString("^[A-Za-z0-9]*$", s)
	return matched
}
