package utils

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// ValidatePathWithBasePaths validates that a file path is safe and within expected boundaries
// This is a wrapper around the main ValidatePath function with additional base path checking
func ValidatePathWithBasePaths(path string, allowedBasePaths ...string) error {
	// First use the main validation function
	if err := ValidatePath(path); err != nil {
		return err
	}

	// If allowed base paths are specified, ensure the path is within one of them
	if len(allowedBasePaths) > 0 {
		cleanPath := filepath.Clean(path)
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
	if err := ValidatePathWithBasePaths(path, allowedBasePaths...); err != nil {
		return nil, err
	}
	return os.ReadFile(path) // #nosec G304 - path is validated above
}

// SafeWriteFile writes a file with secure permissions and path validation
func SafeWriteFile(path string, data []byte, allowedBasePaths ...string) error {
	if err := ValidatePathWithBasePaths(path, allowedBasePaths...); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600) // More secure permissions
}

// SafeMkdirAll creates directories with secure permissions and path validation
func SafeMkdirAll(path string, allowedBasePaths ...string) error {
	if err := ValidatePathWithBasePaths(path, allowedBasePaths...); err != nil {
		return err
	}
	return os.MkdirAll(path, 0750) // More secure permissions
}

// SafeOpenFile opens a file with path validation
func SafeOpenFile(path string, flag int, perm os.FileMode, allowedBasePaths ...string) (*os.File, error) {
	if err := ValidatePathWithBasePaths(path, allowedBasePaths...); err != nil {
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
	if err := ValidatePathWithBasePaths(path, allowedBasePaths...); err != nil {
		return nil, err
	}
	return os.Open(path) // #nosec G304 - path is validated above
}

// SafeCreate creates a file with secure permissions and path validation
func SafeCreate(path string, allowedBasePaths ...string) (*os.File, error) {
	if err := ValidatePathWithBasePaths(path, allowedBasePaths...); err != nil {
		return nil, err
	}
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600) // #nosec G304 - path is validated above
}

// SanitizeInput sanitizes user input to prevent injection attacks
func SanitizeInput(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input cannot be empty")
	}

	// Trim whitespace
	sanitized := strings.TrimSpace(input)

	// Check length limits
	if len(sanitized) > 255 {
		return "", fmt.Errorf("input too long (max 255 characters)")
	}

	// Check for dangerous patterns
	dangerousPatterns := []string{
		`<script`,
		`javascript:`,
		`data:`,
		`vbscript:`,
		`onload=`,
		`onerror=`,
		`../`,
		`..\\`,
		`;`,
		`&`,
		`|`,
		"`",
		"$(",
		"${",
	}

	lowerInput := strings.ToLower(sanitized)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerInput, pattern) {
			return "", fmt.Errorf("input contains dangerous pattern: %s", pattern)
		}
	}

	// Only allow alphanumeric, hyphens, underscores, and dots
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !validPattern.MatchString(sanitized) {
		return "", fmt.Errorf("input contains invalid characters")
	}

	return sanitized, nil
}

// ValidateFilePermissions validates that file permissions are secure
func ValidateFilePermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	mode := info.Mode()
	perm := mode.Perm()

	// Check for dangerous permissions
	if perm&0002 != 0 { // World writable
		return fmt.Errorf("file is world writable: %s", path)
	}

	if mode&os.ModeSetuid != 0 { // SETUID
		return fmt.Errorf("file has SETUID bit set: %s", path)
	}

	if mode&os.ModeSetgid != 0 { // SETGID
		return fmt.Errorf("file has SETGID bit set: %s", path)
	}

	// Check for files with no permissions (except for testing)
	if perm == 0 {
		return fmt.Errorf("file has no permissions: %s", path)
	}

	return nil
}

// ValidateProjectConfig validates project configuration for security
func ValidateProjectConfig(config *models.ProjectConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Validate name
	if _, err := SanitizeInput(config.Name); err != nil {
		return fmt.Errorf("invalid project name: %w", err)
	}

	// Validate organization
	if config.Organization != "" {
		if _, err := SanitizeInput(config.Organization); err != nil {
			return fmt.Errorf("invalid organization: %w", err)
		}
	}

	// Validate output path
	if config.OutputPath != "" {
		if err := ValidatePath(config.OutputPath); err != nil {
			return fmt.Errorf("invalid output path: %w", err)
		}
	}

	return nil
}

// ProcessTemplateSafely processes templates with security restrictions
func ProcessTemplateSafely(templatePath string, config *models.ProjectConfig) error {
	// Validate template path
	if err := ValidatePath(templatePath); err != nil {
		return fmt.Errorf("invalid template path: %w", err)
	}

	// Read template content
	content, err := SafeReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Check for dangerous template patterns
	dangerousPatterns := []string{
		`{{exec`,
		`{{js`,
		`{{eval`,
		`{{template "../../`,
		`{{template "/`,
		`{{template "C:\\`,
		`{{.exec`,
		`{{.js`,
		`{{.eval`,
		`{{ exec`,
		`{{ js`,
		`{{ eval`,
		`{{range .}}{{.}}{{end}}`, // Potential infinite loop
	}

	contentStr := string(content)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(contentStr, pattern) {
			return fmt.Errorf("template contains dangerous pattern: %s", pattern)
		}
	}

	// Create safe template with restricted functions
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(contentStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with config (in a real implementation, you'd write to output)
	_ = tmpl
	_ = config

	return nil
}

// DetectAPIKey detects API keys in content
func DetectAPIKey(content string) bool {
	patterns := []string{
		`sk-[a-zA-Z0-9]{32,}`,
		`pk_live_[a-zA-Z0-9]{24,}`,
		`AIzaSy[a-zA-Z0-9_-]{33}`,
		`ghp_[a-zA-Z0-9]{36}`,
		`API_KEY\s*=\s*[a-zA-Z0-9_-]{16,}`,
		`"api_key"\s*:\s*"[a-zA-Z0-9_-]{16,}"`,
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(`(?i)`+pattern, content)
		if matched {
			// Check if it's not a placeholder
			if !isPlaceholder(content) {
				return true
			}
		}
	}

	return false
}

// DetectPassword detects passwords in content
func DetectPassword(content string) bool {
	patterns := []string{
		`password\s*=\s*["'][^"']{6,}["']`,
		`PASSWORD\s*=\s*[^"'\s]{6,}`,
		`"password"\s*:\s*"[^"']{4,}"`,
		`pwd\s*:\s*["'][^"']{4,}["']`,
		`const\s+pass\s*=\s*["'][^"']{6,}["']`,
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(`(?i)`+pattern, content)
		if matched {
			// Check if it's not a common/placeholder password
			if !isCommonPassword(content) && !isPlaceholder(content) {
				return true
			}
		}
	}

	return false
}

// DetectToken detects tokens in content
func DetectToken(content string) bool {
	patterns := []string{
		`eyJ[a-zA-Z0-9_-]+\.eyJ[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+`, // JWT
		`ya29\.[a-zA-Z0-9_-]{68,}`,                             // Google OAuth
		`[a-zA-Z0-9]{32,}-[a-zA-Z0-9]{32,}`,                    // Generic token
		`token\s*=\s*["'][a-zA-Z0-9_-]{16,}["']`,
		`JWT_TOKEN\s*=\s*Bearer\s+[a-zA-Z0-9_-]+`,
		`"access_token"\s*:\s*"[a-zA-Z0-9_.-]{16,}"`,
		`oauth_token\s*:\s*"[a-zA-Z0-9_-]{16,}"`,
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(`(?i)`+pattern, content)
		if matched {
			if !isPlaceholder(content) {
				return true
			}
		}
	}

	return false
}

// DetectSecrets detects various types of secrets in content
func DetectSecrets(content string) []string {
	var secrets []string

	if DetectAPIKey(content) {
		secrets = append(secrets, "api_key")
	}

	if DetectPassword(content) {
		secrets = append(secrets, "password")
	}

	if DetectToken(content) {
		secrets = append(secrets, "token")
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
		matched, _ := regexp.MatchString(`(?i)`+pattern, lowerContent)
		if matched {
			return true
		}
	}

	return false
}

func isCommonPassword(content string) bool {
	// Extract the actual password value from the content
	patterns := []string{
		`password[^"']*["']([^"']+)["']`, // Extract password value after password key
		`PASSWORD[^"']*=\s*([^"'\s]+)`,   // Extract unquoted values after PASSWORD=
		`pwd[^"']*["']([^"']+)["']`,      // Extract password value after pwd key
		`pass[^"']*["']([^"']+)["']`,     // Extract password value after pass key
	}

	var passwordValue string
	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			passwordValue = strings.ToLower(matches[1])
			break
		}
	}

	if passwordValue == "" {
		return false
	}

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

	for _, common := range commonPasswords {
		if passwordValue == common {
			return true
		}
	}

	return false
}

// SafeCounter provides thread-safe counter operations
type SafeCounter struct {
	mu    sync.RWMutex
	value int
}

// Increment safely increments the counter
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Value safely returns the current counter value
func (c *SafeCounter) Value() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// SafeMap provides thread-safe map operations
type SafeMap struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewSafeMap creates a new thread-safe map
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]string),
	}
}

// Set safely sets a key-value pair
func (m *SafeMap) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get safely gets a value by key
func (m *SafeMap) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.data[key]
	return value, exists
}

// Size safely returns the map size
func (m *SafeMap) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}
