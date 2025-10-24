// Package security provides security scanning functionality for generated projects.
// #nosec G304 - This package scans files for security issues, file operations use validated paths
package security

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SecurityScanner scans generated projects for security issues
type SecurityScanner struct {
	validator *Validator
}

// NewSecurityScanner creates a new security scanner
func NewSecurityScanner() *SecurityScanner {
	return &SecurityScanner{
		validator: NewValidator(),
	}
}

// ScanResult represents the result of a security scan
type ScanResult struct {
	Path     string
	Issues   []*SecurityIssue
	Warnings []*SecurityWarning
	Passed   bool
}

// SecurityIssue represents a security issue found during scanning
type SecurityIssue struct {
	Type        string
	Severity    string // "critical", "high", "medium", "low"
	File        string
	Line        int
	Description string
	Suggestion  string
}

// SecurityWarning represents a security warning
type SecurityWarning struct {
	Type        string
	File        string
	Description string
	Suggestion  string
}

// Scan performs a comprehensive security scan on a directory
func (s *SecurityScanner) Scan(ctx context.Context, rootPath string) (*ScanResult, error) {
	result := &ScanResult{
		Path:     rootPath,
		Issues:   make([]*SecurityIssue, 0),
		Warnings: make([]*SecurityWarning, 0),
		Passed:   true,
	}

	// Scan for exposed secrets
	secretIssues, err := s.scanForSecrets(ctx, rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan for secrets: %w", err)
	}
	result.Issues = append(result.Issues, secretIssues...)

	// Scan for insecure configurations
	configIssues, err := s.scanForInsecureConfigs(ctx, rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan for insecure configs: %w", err)
	}
	result.Issues = append(result.Issues, configIssues...)

	// Scan for insecure file permissions
	permWarnings, err := s.scanFilePermissions(ctx, rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan file permissions: %w", err)
	}
	result.Warnings = append(result.Warnings, permWarnings...)

	// Scan for hardcoded credentials
	credIssues, err := s.scanForHardcodedCredentials(ctx, rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to scan for hardcoded credentials: %w", err)
	}
	result.Issues = append(result.Issues, credIssues...)

	// Determine if scan passed
	for _, issue := range result.Issues {
		if issue.Severity == "critical" || issue.Severity == "high" {
			result.Passed = false
			break
		}
	}

	return result, nil
}

// scanForSecrets scans files for exposed secrets
func (s *SecurityScanner) scanForSecrets(ctx context.Context, rootPath string) ([]*SecurityIssue, error) {
	issues := make([]*SecurityIssue, 0)

	// Define file patterns to scan
	scanPatterns := []string{
		"*.js", "*.ts", "*.jsx", "*.tsx",
		"*.go", "*.py", "*.java", "*.kt",
		"*.yaml", "*.yml", "*.json",
		"*.env", ".env.*",
		"*.config.js", "*.config.ts",
	}

	// Walk through directory
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Skip directories and non-matching files
		if info.IsDir() {
			// Skip common directories that shouldn't be scanned
			if shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file matches scan patterns
		if !matchesAnyPattern(info.Name(), scanPatterns) {
			return nil
		}

		// Scan file for secrets
		fileIssues, err := s.scanFileForSecrets(path)
		if err != nil {
			// Log error but continue scanning - don't fail entire scan for one file
			_ = err // Explicitly ignore error to continue scanning
			return nil
		}

		issues = append(issues, fileIssues...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return issues, nil
}

// scanFileForSecrets scans a single file for secrets
func (s *SecurityScanner) scanFileForSecrets(filePath string) ([]*SecurityIssue, error) {
	issues := make([]*SecurityIssue, 0)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Log error but don't override return value
		}
	}()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Use validator to detect secrets
		secrets := s.validator.DetectSecrets(line)

		for _, secretType := range secrets {
			issue := &SecurityIssue{
				Type:        "exposed_secret",
				Severity:    "critical",
				File:        filePath,
				Line:        lineNum,
				Description: fmt.Sprintf("Potential %s detected in file", secretType),
				Suggestion:  "Remove hardcoded secrets and use environment variables or secret management systems",
			}
			issues = append(issues, issue)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return issues, nil
}

// scanForInsecureConfigs scans for insecure default configurations
func (s *SecurityScanner) scanForInsecureConfigs(ctx context.Context, rootPath string) ([]*SecurityIssue, error) {
	issues := make([]*SecurityIssue, 0)

	// Scan for insecure Docker configurations
	dockerIssues, err := s.scanDockerConfigs(ctx, rootPath)
	if err != nil {
		return nil, err
	}
	issues = append(issues, dockerIssues...)

	// Scan for insecure environment files
	envIssues, err := s.scanEnvFiles(ctx, rootPath)
	if err != nil {
		return nil, err
	}
	issues = append(issues, envIssues...)

	// Scan for insecure CORS configurations
	corsIssues, err := s.scanCORSConfigs(ctx, rootPath)
	if err != nil {
		return nil, err
	}
	issues = append(issues, corsIssues...)

	return issues, nil
}

// scanDockerConfigs scans Docker configurations for security issues
func (s *SecurityScanner) scanDockerConfigs(ctx context.Context, rootPath string) ([]*SecurityIssue, error) {
	issues := make([]*SecurityIssue, 0)

	// Find Dockerfile and docker-compose.yml files
	dockerFiles := []string{"Dockerfile", "docker-compose.yml", "docker-compose.yaml"}

	for _, dockerFile := range dockerFiles {
		filePath := filepath.Join(rootPath, dockerFile)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}

		// Scan file
		fileIssues, err := s.scanDockerFile(filePath)
		if err != nil {
			continue
		}

		issues = append(issues, fileIssues...)
	}

	return issues, nil
}

// scanDockerFile scans a Docker file for security issues
func (s *SecurityScanner) scanDockerFile(filePath string) ([]*SecurityIssue, error) {
	issues := make([]*SecurityIssue, 0)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Log error but don't override return value
		}
	}()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Check for running as root
		if strings.HasPrefix(line, "USER root") {
			issue := &SecurityIssue{
				Type:        "insecure_docker_config",
				Severity:    "high",
				File:        filePath,
				Line:        lineNum,
				Description: "Container running as root user",
				Suggestion:  "Create and use a non-root user in the Dockerfile",
			}
			issues = append(issues, issue)
		}

		// Check for privileged mode
		if strings.Contains(line, "privileged: true") {
			issue := &SecurityIssue{
				Type:        "insecure_docker_config",
				Severity:    "critical",
				File:        filePath,
				Line:        lineNum,
				Description: "Container running in privileged mode",
				Suggestion:  "Remove privileged mode unless absolutely necessary",
			}
			issues = append(issues, issue)
		}

		// Check for exposed ports without restrictions
		if strings.Contains(line, "0.0.0.0:") {
			issue := &SecurityIssue{
				Type:        "insecure_docker_config",
				Severity:    "medium",
				File:        filePath,
				Line:        lineNum,
				Description: "Port exposed on all interfaces (0.0.0.0)",
				Suggestion:  "Bind to localhost (127.0.0.1) if external access is not required",
			}
			issues = append(issues, issue)
		}
	}

	return issues, nil
}

// scanEnvFiles scans environment files for security issues
func (s *SecurityScanner) scanEnvFiles(ctx context.Context, rootPath string) ([]*SecurityIssue, error) {
	issues := make([]*SecurityIssue, 0)

	// Find .env files
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if it's an env file
		if strings.HasPrefix(info.Name(), ".env") && !strings.HasSuffix(info.Name(), ".example") {
			// Check if .env is in .gitignore
			gitignorePath := filepath.Join(filepath.Dir(path), ".gitignore")
			if !isInGitignore(gitignorePath, info.Name()) {
				issue := &SecurityIssue{
					Type:        "exposed_env_file",
					Severity:    "high",
					File:        path,
					Line:        0,
					Description: "Environment file not in .gitignore",
					Suggestion:  "Add .env files to .gitignore to prevent committing secrets",
				}
				issues = append(issues, issue)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return issues, nil
}

// scanCORSConfigs scans for insecure CORS configurations
func (s *SecurityScanner) scanCORSConfigs(ctx context.Context, rootPath string) ([]*SecurityIssue, error) {
	issues := make([]*SecurityIssue, 0)

	// Patterns to scan
	scanPatterns := []string{"*.js", "*.ts", "*.go", "*.py"}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if !matchesAnyPattern(info.Name(), scanPatterns) {
			return nil
		}

		// Scan file for CORS issues
		fileIssues, err := s.scanFileForCORS(path)
		if err != nil {
			// Log error but continue scanning - don't fail entire scan for one file
			_ = err // Explicitly ignore error to continue scanning
			return nil
		}

		issues = append(issues, fileIssues...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return issues, nil
}

// scanFileForCORS scans a file for CORS configuration issues
func (s *SecurityScanner) scanFileForCORS(filePath string) ([]*SecurityIssue, error) {
	issues := make([]*SecurityIssue, 0)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Log error but don't override return value
		}
	}()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Patterns for insecure CORS
	insecureCORSPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)Access-Control-Allow-Origin.*\*`),
		regexp.MustCompile(`(?i)AllowOrigins.*\*`),
		regexp.MustCompile(`(?i)cors.*origin.*\*`),
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		for _, pattern := range insecureCORSPatterns {
			if pattern.MatchString(line) {
				issue := &SecurityIssue{
					Type:        "insecure_cors",
					Severity:    "medium",
					File:        filePath,
					Line:        lineNum,
					Description: "CORS configured to allow all origins (*)",
					Suggestion:  "Restrict CORS to specific trusted origins",
				}
				issues = append(issues, issue)
				break
			}
		}
	}

	return issues, nil
}

// scanFilePermissions scans file permissions for security issues
func (s *SecurityScanner) scanFilePermissions(ctx context.Context, rootPath string) ([]*SecurityWarning, error) {
	warnings := make([]*SecurityWarning, 0)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check file permissions
		if err := s.validator.ValidateFilePermissions(path); err != nil {
			warning := &SecurityWarning{
				Type:        "insecure_permissions",
				File:        path,
				Description: err.Error(),
				Suggestion:  "Review and fix file permissions",
			}
			warnings = append(warnings, warning)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return warnings, nil
}

// scanForHardcodedCredentials scans for hardcoded credentials
func (s *SecurityScanner) scanForHardcodedCredentials(ctx context.Context, rootPath string) ([]*SecurityIssue, error) {
	issues := make([]*SecurityIssue, 0)

	scanPatterns := []string{"*.js", "*.ts", "*.go", "*.py", "*.java", "*.kt"}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if !matchesAnyPattern(info.Name(), scanPatterns) {
			return nil
		}

		// Scan file for hardcoded credentials
		fileIssues, err := s.scanFileForCredentials(path)
		if err != nil {
			// Log error but continue scanning - don't fail entire scan for one file
			_ = err // Explicitly ignore error to continue scanning
			return nil
		}

		issues = append(issues, fileIssues...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return issues, nil
}

// scanFileForCredentials scans a file for hardcoded credentials
func (s *SecurityScanner) scanFileForCredentials(filePath string) ([]*SecurityIssue, error) {
	issues := make([]*SecurityIssue, 0)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Log error but don't override return value
		}
	}()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Patterns for hardcoded credentials
	credentialPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(username|user)\s*[:=]\s*["'][^"']+["']`),
		regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*["'][^"']+["']`),
		regexp.MustCompile(`(?i)(secret|token)\s*[:=]\s*["'][^"']+["']`),
		regexp.MustCompile(`(?i)(api_key|apikey)\s*[:=]\s*["'][^"']+["']`),
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "//") || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		for _, pattern := range credentialPatterns {
			if pattern.MatchString(line) {
				// Check if it's a placeholder
				if !isPlaceholder(line) && !isCommonPassword(line) {
					issue := &SecurityIssue{
						Type:        "hardcoded_credential",
						Severity:    "high",
						File:        filePath,
						Line:        lineNum,
						Description: "Potential hardcoded credential detected",
						Suggestion:  "Use environment variables or secret management for credentials",
					}
					issues = append(issues, issue)
					break
				}
			}
		}
	}

	return issues, nil
}

// Helper functions

func shouldSkipDirectory(dirName string) bool {
	skipDirs := []string{
		"node_modules",
		".git",
		".next",
		"dist",
		"build",
		"vendor",
		".temp",
		".backups",
		"__pycache__",
		".pytest_cache",
		".venv",
		"venv",
	}

	for _, skip := range skipDirs {
		if dirName == skip {
			return true
		}
	}

	return false
}

func matchesAnyPattern(filename string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, _ := filepath.Match(pattern, filename)
		if matched {
			return true
		}
	}
	return false
}

func isInGitignore(gitignorePath, filename string) bool {
	file, err := os.Open(gitignorePath)
	if err != nil {
		return false
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Log error but don't override return value
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == filename || line == "/"+filename {
			return true
		}
	}

	return false
}
