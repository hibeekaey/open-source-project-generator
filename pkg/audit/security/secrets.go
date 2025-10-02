// Package security provides security auditing functionality for the
// Open Source Project Generator.
//
// Security Note: This package contains security audit functionality that legitimately needs
// to read files for security analysis. The G304 warnings from gosec are false positives
// in this context as file reading is the core functionality of a security audit tool.
package security

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// SecretDetector implements secret detection functionality for security auditing.
type SecretDetector struct {
	// rules contains the secret detection rules
	rules []SecretRule
}

// SecretRule defines a rule for detecting secrets with pattern-based matching.
type SecretRule struct {
	Name       string  // Human-readable name of the secret type
	Pattern    string  // Regular expression pattern to match
	Confidence float64 // Confidence score (0.0 to 1.0)
	Category   string  // Category of the secret (api_key, password, token, etc.)
}

// NewSecretDetector creates a new secret detector instance with default rules.
func NewSecretDetector() *SecretDetector {
	return &SecretDetector{
		rules: getDefaultSecretRules(),
	}
}

// NewSecretDetectorWithRules creates a new secret detector with custom rules.
func NewSecretDetectorWithRules(rules []SecretRule) *SecretDetector {
	return &SecretDetector{
		rules: rules,
	}
}

// DetectSecrets detects secrets in the project using pattern-based detection.
func (sd *SecretDetector) DetectSecrets(path string) (*interfaces.SecretScanResult, error) {
	result := &interfaces.SecretScanResult{
		ScanTime: time.Now(),
		Secrets:  []interfaces.SecretDetection{},
		Summary: interfaces.SecretScanSummary{
			TotalSecrets:     0,
			HighConfidence:   0,
			MediumConfidence: 0,
			LowConfidence:    0,
			FilesScanned:     0,
		},
	}

	// Walk through project files
	walkErr := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-text files
		if info.IsDir() || sd.shouldSkipFile(filePath) {
			return nil
		}

		result.Summary.FilesScanned++

		// Scan file for secrets
		secrets, err := sd.scanFileForSecrets(filePath)
		if err != nil {
			// Log error but continue scanning
			return nil
		}

		result.Secrets = append(result.Secrets, secrets...)
		return nil
	})

	if walkErr != nil {
		return nil, fmt.Errorf("failed to scan for secrets: %w", walkErr)
	}

	// Update summary with confidence scoring
	for _, secret := range result.Secrets {
		result.Summary.TotalSecrets++
		if secret.Confidence >= 0.8 {
			result.Summary.HighConfidence++
		} else if secret.Confidence >= 0.5 {
			result.Summary.MediumConfidence++
		} else {
			result.Summary.LowConfidence++
		}
	}

	return result, nil
}

// scanFileForSecrets scans a single file for potential secrets using all configured rules.
func (sd *SecretDetector) scanFileForSecrets(filePath string) ([]interfaces.SecretDetection, error) {
	var secrets []interfaces.SecretDetection

	content, err := os.ReadFile(filePath) // #nosec G304 - This is an audit tool that needs to read files
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	// Apply each rule to the file content
	for _, rule := range sd.rules {
		regex, err := regexp.Compile(rule.Pattern)
		if err != nil {
			// Skip invalid regex patterns
			continue
		}

		// Scan each line for matches
		for lineNum, line := range lines {
			matches := regex.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) > 0 {
					secret := match[0]

					// Calculate confidence score with context analysis
					confidence := sd.calculateConfidence(rule, secret, line, filePath)

					// Only include secrets above minimum confidence threshold
					if confidence >= 0.3 {
						secrets = append(secrets, interfaces.SecretDetection{
							Type:       rule.Name,
							File:       filePath,
							Line:       lineNum + 1,
							Column:     strings.Index(line, secret) + 1,
							Secret:     secret,
							Confidence: confidence,
							Rule:       rule.Name,
							Pattern:    rule.Pattern,
							Masked:     sd.maskSecret(secret),
						})
					}
				}
			}
		}
	}

	return secrets, nil
}

// calculateConfidence calculates the confidence score for a detected secret
// based on the rule, context, and file characteristics.
func (sd *SecretDetector) calculateConfidence(rule SecretRule, secret, line, filePath string) float64 {
	confidence := rule.Confidence

	// Reduce confidence for test files
	if sd.isTestFile(filePath) {
		confidence *= 0.7
	}

	// Reduce confidence for example/demo files
	if sd.isExampleFile(filePath) {
		confidence *= 0.5
	}

	// Reduce confidence for obvious placeholders
	if sd.isPlaceholder(secret) {
		confidence *= 0.3
	}

	// Increase confidence for high-entropy strings
	if sd.hasHighEntropy(secret) {
		confidence *= 1.2
	}

	// Reduce confidence for common false positives
	if sd.isFalsePositive(secret, line) {
		confidence *= 0.4
	}

	// Ensure confidence stays within bounds
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

// isTestFile checks if the file is a test file
func (sd *SecretDetector) isTestFile(filePath string) bool {
	testPatterns := []string{
		"test", "spec", "mock", "fixture", "example",
		"_test.", ".test.", "_spec.", ".spec.",
	}

	lowerPath := strings.ToLower(filePath)
	for _, pattern := range testPatterns {
		if strings.Contains(lowerPath, pattern) {
			return true
		}
	}
	return false
}

// isExampleFile checks if the file is an example or documentation file
func (sd *SecretDetector) isExampleFile(filePath string) bool {
	examplePatterns := []string{
		"example", "demo", "sample", "template", "readme",
		"doc", "docs", ".md", ".txt", ".example",
	}

	lowerPath := strings.ToLower(filePath)
	for _, pattern := range examplePatterns {
		if strings.Contains(lowerPath, pattern) {
			return true
		}
	}
	return false
}

// isPlaceholder checks if the secret appears to be a placeholder
func (sd *SecretDetector) isPlaceholder(secret string) bool {
	// Check for obvious placeholder patterns
	lowerSecret := strings.ToLower(secret)

	// Full word placeholders
	fullWordPlaceholders := []string{
		"your_api_key", "your_secret", "your_password",
		"placeholder", "dummy", "fake", "sample",
		"xxx", "yyy", "zzz", "test_key", "test_secret",
		"password", "secret", "key", "token",
	}

	for _, placeholder := range fullWordPlaceholders {
		if lowerSecret == placeholder {
			return true
		}
	}

	// Pattern-based placeholders
	if strings.HasPrefix(lowerSecret, "your_") ||
		strings.HasPrefix(lowerSecret, "example_") ||
		strings.HasPrefix(lowerSecret, "placeholder_") ||
		strings.HasPrefix(lowerSecret, "dummy_") ||
		strings.HasPrefix(lowerSecret, "fake_") ||
		strings.HasPrefix(lowerSecret, "test_") ||
		strings.HasPrefix(lowerSecret, "sample_") {
		return true
	}

	// Simple repeated patterns
	if lowerSecret == "1234" || lowerSecret == "abcd" ||
		lowerSecret == "12345678" || lowerSecret == "abcdefgh" {
		return true
	}

	return false
}

// hasHighEntropy checks if the string has high entropy (randomness)
func (sd *SecretDetector) hasHighEntropy(secret string) bool {
	if len(secret) < 8 {
		return false
	}

	// Count unique characters
	charMap := make(map[rune]bool)
	for _, char := range secret {
		charMap[char] = true
	}

	// Calculate entropy ratio
	entropyRatio := float64(len(charMap)) / float64(len(secret))

	// Check for sequential patterns (low entropy)
	if sd.hasSequentialPattern(secret) {
		return false
	}

	// Check for repeated patterns (low entropy)
	if sd.hasRepeatedPattern(secret) {
		return false
	}

	return entropyRatio > 0.6
}

// hasSequentialPattern checks for sequential character patterns
func (sd *SecretDetector) hasSequentialPattern(secret string) bool {
	lowerSecret := strings.ToLower(secret)

	// Check for numeric sequences
	sequences := []string{
		"0123456789", "1234567890", "9876543210",
		"abcdefghijklmnopqrstuvwxyz", "zyxwvutsrqponmlkjihgfedcba",
	}

	for _, seq := range sequences {
		for i := 0; i <= len(seq)-4; i++ {
			if strings.Contains(lowerSecret, seq[i:i+4]) {
				return true
			}
		}
	}

	return false
}

// hasRepeatedPattern checks for repeated character patterns
func (sd *SecretDetector) hasRepeatedPattern(secret string) bool {
	// Check for repeated characters (more than 3 in a row)
	for i := 0; i < len(secret)-3; i++ {
		if secret[i] == secret[i+1] && secret[i+1] == secret[i+2] && secret[i+2] == secret[i+3] {
			return true
		}
	}

	// Check for repeated short patterns
	if len(secret) >= 8 {
		for patternLen := 2; patternLen <= 4; patternLen++ {
			for i := 0; i <= len(secret)-patternLen*2; i++ {
				pattern := secret[i : i+patternLen]
				if strings.Count(secret, pattern) >= 3 {
					return true
				}
			}
		}
	}

	return false
}

// isFalsePositive checks for common false positive patterns
func (sd *SecretDetector) isFalsePositive(secret, line string) bool {
	// Check for common false positive patterns
	falsePositives := []string{
		"console.log", "print", "echo", "debug", "log",
		"//", "/*", "#", "<!--", "TODO", "FIXME",
		"http://", "https://", "ftp://", "file://",
	}

	lowerLine := strings.ToLower(line)
	for _, fp := range falsePositives {
		if strings.Contains(lowerLine, fp) {
			return true
		}
	}

	// Check if it's in a comment
	trimmedLine := strings.TrimSpace(line)
	if strings.HasPrefix(trimmedLine, "//") ||
		strings.HasPrefix(trimmedLine, "#") ||
		strings.HasPrefix(trimmedLine, "/*") ||
		strings.Contains(trimmedLine, "<!--") {
		return true
	}

	return false
}

// shouldSkipFile determines if a file should be skipped during scanning
func (sd *SecretDetector) shouldSkipFile(filePath string) bool {
	// Skip binary files, images, and other non-text files
	skipExtensions := []string{
		".exe", ".bin", ".dll", ".so", ".dylib",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico",
		".mp3", ".mp4", ".avi", ".mov", ".wmv", ".flv",
		".zip", ".tar", ".gz", ".rar", ".7z", ".bz2",
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".class", ".jar", ".war", ".ear", ".dex", ".apk",
		".o", ".obj", ".lib", ".a", ".pyc", ".pyo",
		".woff", ".woff2", ".ttf", ".eot", ".otf",
	}

	// Skip directories that typically don't contain secrets
	skipDirs := []string{
		"node_modules", "vendor", "target", "build", "dist", "out",
		".git", ".svn", ".hg", ".bzr", ".vscode", ".idea",
		"__pycache__", ".pytest_cache", ".coverage", "coverage",
		"logs", "log", "tmp", "temp", ".tmp", ".temp",
		"bin", "obj", "Debug", "Release", "packages",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, skipExt := range skipExtensions {
		if ext == skipExt {
			return true
		}
	}

	// Check if any part of the path contains skip directories
	pathParts := strings.Split(filePath, string(filepath.Separator))
	for _, part := range pathParts {
		for _, skipDir := range skipDirs {
			if part == skipDir {
				return true
			}
		}
	}

	// Skip hidden files (but allow .env files)
	baseName := filepath.Base(filePath)
	if strings.HasPrefix(baseName, ".") && !strings.HasPrefix(baseName, ".env") {
		return true
	}

	return false
}

// maskSecret masks a secret for safe display while preserving some context
func (sd *SecretDetector) maskSecret(secret string) string {
	if len(secret) <= 4 {
		return strings.Repeat("*", len(secret))
	}

	// Show first 2 and last 2 characters for context
	return secret[:2] + strings.Repeat("*", len(secret)-4) + secret[len(secret)-2:]
}

// AddRule adds a new secret detection rule
func (sd *SecretDetector) AddRule(rule SecretRule) error {
	// Validate the rule
	if rule.Name == "" {
		return fmt.Errorf("rule name cannot be empty")
	}
	if rule.Pattern == "" {
		return fmt.Errorf("rule pattern cannot be empty")
	}
	if rule.Confidence < 0.0 || rule.Confidence > 1.0 {
		return fmt.Errorf("rule confidence must be between 0.0 and 1.0")
	}

	// Test the regex pattern
	_, err := regexp.Compile(rule.Pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	sd.rules = append(sd.rules, rule)
	return nil
}

// GetRules returns a copy of the current secret detection rules
func (sd *SecretDetector) GetRules() []SecretRule {
	rules := make([]SecretRule, len(sd.rules))
	copy(rules, sd.rules)
	return rules
}

// SetRules replaces all current rules with the provided rules
func (sd *SecretDetector) SetRules(rules []SecretRule) error {
	// Validate all rules first
	for _, rule := range rules {
		if rule.Name == "" {
			return fmt.Errorf("rule name cannot be empty")
		}
		if rule.Pattern == "" {
			return fmt.Errorf("rule pattern cannot be empty")
		}
		if rule.Confidence < 0.0 || rule.Confidence > 1.0 {
			return fmt.Errorf("rule confidence must be between 0.0 and 1.0")
		}

		// Test the regex pattern
		_, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern in rule %s: %w", rule.Name, err)
		}
	}

	sd.rules = make([]SecretRule, len(rules))
	copy(sd.rules, rules)
	return nil
}

// getDefaultSecretRules returns the default set of secret detection rules
func getDefaultSecretRules() []SecretRule {
	return []SecretRule{
		{
			Name:       "AWS Access Key ID",
			Pattern:    `AKIA[0-9A-Z]{16}`,
			Confidence: 0.9,
			Category:   "aws_credential",
		},
		{
			Name:       "AWS Secret Access Key",
			Pattern:    `[0-9a-zA-Z/+]{40}`,
			Confidence: 0.6,
			Category:   "aws_credential",
		},
		{
			Name:       "AWS Session Token",
			Pattern:    `AQoEXAMPLEH4aoAH0gNCAPyJxz4BlCFFxWNE1OPTgk5TthT\+FvwqnKwRcOIfrRh3c/LTo6UDdyJwOOvEVPvLXCrrrUtdnniCEXAMPLE/IvU1BN1bdbxHiQjHlUt0L/ZEXAMPLE`,
			Confidence: 0.8,
			Category:   "aws_credential",
		},
		{
			Name:       "GitHub Personal Access Token",
			Pattern:    `ghp_[0-9a-zA-Z]{36}`,
			Confidence: 0.9,
			Category:   "github_token",
		},
		{
			Name:       "GitHub OAuth Token",
			Pattern:    `gho_[0-9a-zA-Z]{36}`,
			Confidence: 0.9,
			Category:   "github_token",
		},
		{
			Name:       "GitHub App Token",
			Pattern:    `ghs_[0-9a-zA-Z]{36}`,
			Confidence: 0.9,
			Category:   "github_token",
		},
		{
			Name:       "GitHub Refresh Token",
			Pattern:    `ghr_[0-9a-zA-Z]{76}`,
			Confidence: 0.9,
			Category:   "github_token",
		},
		{
			Name:       "Generic API Key",
			Pattern:    `(?i)(api[_-]?key|apikey)\s*[:=]\s*['""]?[0-9a-zA-Z]{20,}['""]?`,
			Confidence: 0.7,
			Category:   "api_key",
		},
		{
			Name:       "Generic Secret",
			Pattern:    `(?i)(secret|password|passwd|pwd)\s*[:=]\s*['""]?[0-9a-zA-Z!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]{8,}['""]?`,
			Confidence: 0.5,
			Category:   "password",
		},
		{
			Name:       "Private Key (RSA)",
			Pattern:    `-----BEGIN\s+(RSA\s+)?PRIVATE\s+KEY-----`,
			Confidence: 0.95,
			Category:   "private_key",
		},
		{
			Name:       "Private Key (OpenSSH)",
			Pattern:    `-----BEGIN\s+OPENSSH\s+PRIVATE\s+KEY-----`,
			Confidence: 0.95,
			Category:   "private_key",
		},
		{
			Name:       "Private Key (EC)",
			Pattern:    `-----BEGIN\s+EC\s+PRIVATE\s+KEY-----`,
			Confidence: 0.95,
			Category:   "private_key",
		},
		{
			Name:       "JWT Token",
			Pattern:    `eyJ[0-9a-zA-Z_-]*\.eyJ[0-9a-zA-Z_-]*\.[0-9a-zA-Z_-]*`,
			Confidence: 0.8,
			Category:   "jwt_token",
		},
		{
			Name:       "Google API Key",
			Pattern:    `AIza[0-9A-Za-z\\-_]{35}`,
			Confidence: 0.9,
			Category:   "google_api_key",
		},
		{
			Name:       "Slack Token",
			Pattern:    `xox[baprs]-([0-9a-zA-Z]{10,48})?`,
			Confidence: 0.8,
			Category:   "slack_token",
		},
		{
			Name:       "Stripe API Key",
			Pattern:    `sk_live_[0-9a-zA-Z]{24}`,
			Confidence: 0.9,
			Category:   "stripe_key",
		},
		{
			Name:       "Stripe Publishable Key",
			Pattern:    `pk_live_[0-9a-zA-Z]{24}`,
			Confidence: 0.8,
			Category:   "stripe_key",
		},
		{
			Name:       "Twilio API Key",
			Pattern:    `SK[0-9a-fA-F]{32}`,
			Confidence: 0.8,
			Category:   "twilio_key",
		},
		{
			Name:       "Database Connection String",
			Pattern:    `(?i)(mongodb|mysql|postgresql|postgres)://[^\s'"]+`,
			Confidence: 0.7,
			Category:   "database_url",
		},
		{
			Name:       "Generic Token",
			Pattern:    `(?i)(token|bearer)\s*[:=]\s*['""]?[0-9a-zA-Z]{16,}['""]?`,
			Confidence: 0.6,
			Category:   "token",
		},
	}
}
