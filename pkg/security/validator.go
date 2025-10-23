// Package security provides input sanitization and validation for the CLI tool.
package security

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Validator provides security validation functionality.
type Validator struct {
	sanitizer *Sanitizer
}

// NewValidator creates a new validator with a sanitizer.
func NewValidator() *Validator {
	return &Validator{
		sanitizer: NewSanitizer(),
	}
}

// ValidatePathSecurity validates that a path is safe to use.
// Checks for path traversal, null bytes, and other security issues.
func (v *Validator) ValidatePathSecurity(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Use sanitizer to validate
	_, err := v.sanitizer.SanitizePath(path)
	return err
}

// ValidatePathWithinBase validates that a path is within an allowed base directory.
// This prevents path traversal attacks.
func (v *Validator) ValidatePathWithinBase(path string, basePath string) error {
	// First validate the path itself
	if err := v.ValidatePathSecurity(path); err != nil {
		return err
	}

	// Validate base path
	if err := v.ValidatePathSecurity(basePath); err != nil {
		return fmt.Errorf("invalid base path: %w", err)
	}

	// Clean both paths
	cleanPath := filepath.Clean(path)
	cleanBase := filepath.Clean(basePath)

	// Convert to absolute paths for comparison
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	absBase, err := filepath.Abs(cleanBase)
	if err != nil {
		return fmt.Errorf("failed to resolve base path: %w", err)
	}

	// Check if path is within base
	relPath, err := filepath.Rel(absBase, absPath)
	if err != nil {
		return fmt.Errorf("path is not within base directory: %w", err)
	}

	// Check for path traversal
	if strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("path is outside base directory")
	}

	return nil
}

// ValidateFilePermissions validates that file permissions are secure.
// Checks for world-writable files and dangerous permission bits.
func (v *Validator) ValidateFilePermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	mode := info.Mode()
	perm := mode.Perm()

	// Check for world writable
	if perm&0002 != 0 {
		return fmt.Errorf("file is world writable: %s", path)
	}

	// Check for SETUID bit
	if mode&os.ModeSetuid != 0 {
		return fmt.Errorf("file has SETUID bit set: %s", path)
	}

	// Check for SETGID bit
	if mode&os.ModeSetgid != 0 {
		return fmt.Errorf("file has SETGID bit set: %s", path)
	}

	// Check for files with no permissions
	if perm == 0 {
		return fmt.Errorf("file has no permissions: %s", path)
	}

	return nil
}

// ValidateToolCommand validates that a tool command is safe to execute.
// Checks for command injection attempts and validates against whitelist.
func (v *Validator) ValidateToolCommand(command string, allowedCommands []string) error {
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	// Check if command is in whitelist
	allowed := false
	for _, allowedCmd := range allowedCommands {
		if command == allowedCmd {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("command '%s' is not in whitelist", command)
	}

	// Check for dangerous characters that could be used for injection
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\n", "\r"}
	for _, char := range dangerousChars {
		if strings.Contains(command, char) {
			return fmt.Errorf("command contains dangerous character '%s'", char)
		}
	}

	return nil
}

// ValidateToolFlags validates that tool flags are safe to use.
// Checks for injection attempts and validates flag format.
func (v *Validator) ValidateToolFlags(flags []string) error {
	for _, flag := range flags {
		// Check for empty flags
		if strings.TrimSpace(flag) == "" {
			return fmt.Errorf("flag cannot be empty")
		}

		// Check for dangerous characters
		dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\n", "\r"}
		for _, char := range dangerousChars {
			if strings.Contains(flag, char) {
				return fmt.Errorf("flag contains dangerous character '%s'", char)
			}
		}

		// Check for path traversal in flags
		if strings.Contains(flag, "..") {
			return fmt.Errorf("flag contains path traversal attempt")
		}
	}

	return nil
}

// ValidateProjectName validates a project name for security and format.
func (v *Validator) ValidateProjectName(name string) error {
	_, err := v.sanitizer.SanitizeProjectName(name)
	return err
}

// ValidateConfigValue validates a configuration value for security.
func (v *Validator) ValidateConfigValue(value string) error {
	_, err := v.sanitizer.SanitizeConfigValue(value)
	return err
}

// DetectSecrets detects potential secrets in content.
// Returns a list of secret types found.
func (v *Validator) DetectSecrets(content string) []string {
	var secrets []string

	// API key patterns
	apiKeyPatterns := []string{
		`sk-[a-zA-Z0-9]{32,}`,
		`pk_live_[a-zA-Z0-9]{24,}`,
		`AIzaSy[a-zA-Z0-9_-]{33}`,
		`ghp_[a-zA-Z0-9]{36}`,
		`API_KEY\s*=\s*[a-zA-Z0-9_-]{16,}`,
		`"api_key"\s*:\s*"[a-zA-Z0-9_-]{16,}"`,
	}

	for _, pattern := range apiKeyPatterns {
		if matched, _ := regexp.MatchString(`(?i)`+pattern, content); matched {
			if !isPlaceholder(content) {
				secrets = append(secrets, "api_key")
				break
			}
		}
	}

	// Password patterns
	passwordPatterns := []string{
		`password\s*=\s*["'][^"']{6,}["']`,
		`PASSWORD\s*=\s*[^"'\s]{6,}`,
		`"password"\s*:\s*"[^"']{4,}"`,
		`pwd\s*:\s*["'][^"']{4,}["']`,
	}

	for _, pattern := range passwordPatterns {
		if matched, _ := regexp.MatchString(`(?i)`+pattern, content); matched {
			if !isPlaceholder(content) && !isCommonPassword(content) {
				secrets = append(secrets, "password")
				break
			}
		}
	}

	// Token patterns (JWT, OAuth, etc.)
	tokenPatterns := []string{
		`eyJ[a-zA-Z0-9_-]+\.eyJ[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+`, // JWT
		`ya29\.[a-zA-Z0-9_-]{68,}`,                             // Google OAuth
		`token\s*=\s*["'][a-zA-Z0-9_-]{16,}["']`,
		`"access_token"\s*:\s*"[a-zA-Z0-9_.-]{16,}"`,
	}

	for _, pattern := range tokenPatterns {
		if matched, _ := regexp.MatchString(`(?i)`+pattern, content); matched {
			if !isPlaceholder(content) {
				secrets = append(secrets, "token")
				break
			}
		}
	}

	return secrets
}

// Helper functions

func isPlaceholder(content string) bool {
	placeholderPatterns := []string{
		`your_.*_here`,
		`YOUR_.*_HERE`,
		`<.*>`,
		`\{\{.*\}\}`,
		`example`,
		`placeholder`,
		`PLACEHOLDER`,
		`template`,
	}

	lowerContent := strings.ToLower(content)
	for _, pattern := range placeholderPatterns {
		if matched, _ := regexp.MatchString(`(?i)`+pattern, lowerContent); matched {
			return true
		}
	}

	return false
}

func isCommonPassword(content string) bool {
	commonPasswords := []string{
		"password",
		"123456",
		"admin",
		"root",
		"guest",
		"test",
		"example",
		"sample",
	}

	lowerContent := strings.ToLower(content)
	for _, common := range commonPasswords {
		if strings.Contains(lowerContent, common) {
			return true
		}
	}

	return false
}
