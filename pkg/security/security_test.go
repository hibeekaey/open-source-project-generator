package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// TestInputValidationAndSanitization tests input validation and sanitization
func TestInputValidationAndSanitization(t *testing.T) {
	t.Run("path_traversal_prevention", func(t *testing.T) {
		testPathTraversalPrevention(t)
	})

	t.Run("input_sanitization", func(t *testing.T) {
		testInputSanitization(t)
	})

	t.Run("file_permission_validation", func(t *testing.T) {
		testFilePermissionValidation(t)
	})

	t.Run("configuration_validation", func(t *testing.T) {
		testConfigurationValidation(t)
	})
}

func testPathTraversalPrevention(t *testing.T) {
	// Test cases for path traversal attacks
	maliciousPaths := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32\\config\\sam",
		"/etc/passwd",
		"C:\\Windows\\System32\\config\\SAM",
		"file:///etc/passwd",
		"../config/../../../etc/passwd",
		"config/../../etc/passwd",
		"./../../etc/passwd",
		"config\\..\\..\\..\\etc\\passwd",
	}

	for _, maliciousPath := range maliciousPaths {
		t.Run("path_"+maliciousPath, func(t *testing.T) {
			err := utils.ValidatePath(maliciousPath)
			if err == nil {
				t.Errorf("Expected path validation to fail for malicious path: %s", maliciousPath)
			}
		})
	}

	// Test valid paths
	validPaths := []string{
		"config.yaml",
		"src/main.go",
		"./config/app.yaml",
		"templates/backend/go-gin",
		"output/project",
	}

	for _, validPath := range validPaths {
		t.Run("valid_path_"+validPath, func(t *testing.T) {
			err := utils.ValidatePath(validPath)
			if err != nil {
				t.Errorf("Expected path validation to pass for valid path: %s, got error: %v", validPath, err)
			}
		})
	}
}

func testInputSanitization(t *testing.T) {
	// Test project name sanitization
	testCases := []struct {
		input    string
		expected string
		valid    bool
	}{
		{"valid-project-name", "valid-project-name", true},
		{"ValidProjectName", "ValidProjectName", true},
		{"project_name_123", "project_name_123", true},
		{"<script>alert('xss')</script>", "", false},
		{"project; rm -rf /", "", false},
		{"project`whoami`", "", false},
		{"project$(whoami)", "", false},
		{"project & echo 'injected'", "", false},
		{"project | cat /etc/passwd", "", false},
		{"../../../etc/passwd", "", false},
		{"", "", false},
		{"a", "a", true},                      // Single character should be valid
		{strings.Repeat("a", 256), "", false}, // Too long
	}

	for _, tc := range testCases {
		t.Run("sanitize_"+tc.input, func(t *testing.T) {
			sanitized, err := utils.SanitizeInput(tc.input)

			if tc.valid {
				if err != nil {
					t.Errorf("Expected input '%s' to be valid, got error: %v", tc.input, err)
				}
				if sanitized != tc.expected {
					t.Errorf("Expected sanitized input '%s', got '%s'", tc.expected, sanitized)
				}
			} else {
				if err == nil {
					t.Errorf("Expected input '%s' to be invalid", tc.input)
				}
			}
		})
	}
}

func testFilePermissionValidation(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files with different permissions
	testFiles := []struct {
		name        string
		permissions os.FileMode
		shouldPass  bool
	}{
		{"normal_file.txt", 0644, true},
		{"executable_file.sh", 0755, true},
		{"readonly_file.txt", 0444, true},
		{"world_writable.txt", 0666, false}, // World writable is dangerous
		{"no_permissions.txt", 0000, false},
		{"setuid_file", 0755 | os.ModeSetuid, false}, // SETUID is dangerous
		{"setgid_file", 0755 | os.ModeSetgid, false}, // SETGID is dangerous
		{"sticky_bit", 0755 | os.ModeSticky, true},   // Sticky bit is OK
	}

	for _, tf := range testFiles {
		t.Run("permission_"+tf.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, tf.name)

			// Create file
			err := os.WriteFile(filePath, []byte("test content"), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Set the actual permissions after file creation
			err = os.Chmod(filePath, tf.permissions)
			if err != nil {
				t.Fatalf("Failed to set file permissions: %v", err)
			}

			// Check if the permissions were actually set (some systems don't allow setuid/setgid)
			info, err := os.Stat(filePath)
			if err != nil {
				t.Fatalf("Failed to stat file: %v", err)
			}

			// Skip test if setuid/setgid bits couldn't be set (system limitation)
			if (tf.permissions&os.ModeSetuid != 0 && info.Mode()&os.ModeSetuid == 0) ||
				(tf.permissions&os.ModeSetgid != 0 && info.Mode()&os.ModeSetgid == 0) {
				t.Skipf("System doesn't allow setting setuid/setgid bits for %s", tf.name)
			}

			// Validate permissions
			err = utils.ValidateFilePermissions(filePath)

			if tf.shouldPass && err != nil {
				t.Errorf("Expected file permissions to be valid for %s, got error: %v", tf.name, err)
			}

			if !tf.shouldPass && err == nil {
				t.Errorf("Expected file permissions to be invalid for %s", tf.name)
			}
		})
	}
}

func testConfigurationValidation(t *testing.T) {
	// Test secure configuration validation
	testConfigs := []struct {
		name   string
		config *models.ProjectConfig
		valid  bool
	}{
		{
			name: "valid_config",
			config: &models.ProjectConfig{
				Name:         "valid-project",
				Organization: "valid-org",
				License:      "MIT",
				OutputPath:   "./output",
			},
			valid: true,
		},
		{
			name: "malicious_name",
			config: &models.ProjectConfig{
				Name:         "../../../etc/passwd",
				Organization: "org",
				License:      "MIT",
			},
			valid: false,
		},
		{
			name: "script_injection_name",
			config: &models.ProjectConfig{
				Name:         "<script>alert('xss')</script>",
				Organization: "org",
				License:      "MIT",
			},
			valid: false,
		},
		{
			name: "command_injection_org",
			config: &models.ProjectConfig{
				Name:         "project",
				Organization: "org; rm -rf /",
				License:      "MIT",
			},
			valid: false,
		},
		{
			name: "malicious_output_path",
			config: &models.ProjectConfig{
				Name:         "project",
				Organization: "org",
				License:      "MIT",
				OutputPath:   "../../../etc",
			},
			valid: false,
		},
	}

	for _, tc := range testConfigs {
		t.Run(tc.name, func(t *testing.T) {
			err := utils.ValidateProjectConfig(tc.config)

			if tc.valid && err != nil {
				t.Errorf("Expected config to be valid, got error: %v", err)
			}

			if !tc.valid && err == nil {
				t.Error("Expected config to be invalid")
			}
		})
	}
}

// TestTemplateProcessingSecurity tests security in template processing
func TestTemplateProcessingSecurity(t *testing.T) {
	t.Run("template_injection_prevention", func(t *testing.T) {
		testTemplateInjectionPrevention(t)
	})

	t.Run("file_inclusion_prevention", func(t *testing.T) {
		testFileInclusionPrevention(t)
	})

	t.Run("safe_template_execution", func(t *testing.T) {
		testSafeTemplateExecution(t)
	})
}

func testTemplateInjectionPrevention(t *testing.T) {
	tempDir := t.TempDir()

	// Create malicious template content
	maliciousTemplates := []string{
		`{{.Name}}{{range .}}{{.}}{{end}}`,         // Potential infinite loop
		`{{.Name}}{{template "../../etc/passwd"}}`, // File inclusion attempt
		`{{.Name}}{{exec "rm -rf /"}}`,             // Command execution attempt
		`{{.Name}}{{js "alert('xss')"}}`,           // JavaScript execution attempt
		`{{.Name}}{{eval "process.exit(1)"}}`,      // Code evaluation attempt
	}

	for i, maliciousTemplate := range maliciousTemplates {
		t.Run(fmt.Sprintf("malicious_template_%d", i), func(t *testing.T) {
			templatePath := filepath.Join(tempDir, fmt.Sprintf("malicious_%d.tmpl", i))

			err := os.WriteFile(templatePath, []byte(maliciousTemplate), 0644)
			if err != nil {
				t.Fatalf("Failed to create malicious template: %v", err)
			}

			// Attempt to process template
			config := &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
			}

			err = utils.ProcessTemplateSafely(templatePath, config)
			if err == nil {
				t.Error("Expected malicious template processing to fail")
			}
		})
	}
}

func testFileInclusionPrevention(t *testing.T) {
	tempDir := t.TempDir()

	// Create template with file inclusion attempts
	inclusionAttempts := []string{
		`{{template "../../etc/passwd"}}`,
		`{{template "/etc/passwd"}}`,
		`{{template "C:\\Windows\\System32\\config\\SAM"}}`,
		`{{template "../../../sensitive.txt"}}`,
	}

	for i, attempt := range inclusionAttempts {
		t.Run(fmt.Sprintf("inclusion_attempt_%d", i), func(t *testing.T) {
			templatePath := filepath.Join(tempDir, fmt.Sprintf("inclusion_%d.tmpl", i))

			templateContent := fmt.Sprintf(`# {{.Name}}

%s

Organization: {{.Organization}}`, attempt)

			err := os.WriteFile(templatePath, []byte(templateContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create inclusion template: %v", err)
			}

			config := &models.ProjectConfig{
				Name:         "test-project",
				Organization: "test-org",
			}

			err = utils.ProcessTemplateSafely(templatePath, config)
			if err == nil {
				t.Error("Expected file inclusion attempt to fail")
			}
		})
	}
}

func testSafeTemplateExecution(t *testing.T) {
	tempDir := t.TempDir()

	// Create safe template
	safeTemplate := `# {{.Name}}

This is a safe template for {{.Organization}}.

License: {{.License}}

## Features

{{range .Features}}
- {{.}}
{{end}}
`

	templatePath := filepath.Join(tempDir, "safe.tmpl")
	err := os.WriteFile(templatePath, []byte(safeTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to create safe template: %v", err)
	}

	config := &models.ProjectConfig{
		Name:         "safe-project",
		Organization: "safe-org",
		License:      "MIT",
		Features:     []string{"feature1", "feature2", "feature3"},
	}

	err = utils.ProcessTemplateSafely(templatePath, config)
	if err != nil {
		t.Errorf("Expected safe template processing to succeed, got error: %v", err)
	}
}

// TestSecretDetection tests secret detection capabilities
func TestSecretDetection(t *testing.T) {
	t.Run("api_key_detection", func(t *testing.T) {
		testAPIKeyDetection(t)
	})

	t.Run("password_detection", func(t *testing.T) {
		testPasswordDetection(t)
	})

	t.Run("token_detection", func(t *testing.T) {
		testTokenDetection(t)
	})

	t.Run("false_positive_prevention", func(t *testing.T) {
		testFalsePositivePrevention(t)
	})
}

func testAPIKeyDetection(t *testing.T) {
	// Test various API key patterns
	apiKeyPatterns := []struct {
		content      string
		shouldDetect bool
	}{
		{`const apiKey = "sk-1234567890abcdef1234567890abcdef";`, true},
		{`API_KEY=pk_live_1234567890abcdef`, true},
		{`"api_key": "AIzaSyD1234567890abcdef"`, true},
		{`apikey: "ghp_1234567890abcdef1234567890abcdef123456"`, true},
		{`const key = "not-a-real-api-key";`, false},
		{`// Example: API_KEY=your_key_here`, false},
		{`const placeholder = "YOUR_API_KEY_HERE";`, false},
	}

	for i, pattern := range apiKeyPatterns {
		t.Run(fmt.Sprintf("api_key_pattern_%d", i), func(t *testing.T) {
			detected := utils.DetectAPIKey(pattern.content)

			if pattern.shouldDetect && !detected {
				t.Errorf("Expected to detect API key in: %s", pattern.content)
			}

			if !pattern.shouldDetect && detected {
				t.Errorf("False positive API key detection in: %s", pattern.content)
			}
		})
	}
}

func testPasswordDetection(t *testing.T) {
	// Test password detection patterns
	passwordPatterns := []struct {
		content      string
		shouldDetect bool
	}{
		{`password = "mySecretPassword123"`, true},
		{`PASSWORD=supersecret456`, true},
		{`"password": "p@ssw0rd!"`, true},
		{`pwd: "admin123"`, true},
		{`const pass = "password123";`, true},
		{`password = "password"`, false},             // Too common
		{`password = ""`, false},                     // Empty
		{`// password: your_password_here`, false},   // Comment
		{`password = "PASSWORD_PLACEHOLDER"`, false}, // Placeholder
	}

	for i, pattern := range passwordPatterns {
		t.Run(fmt.Sprintf("password_pattern_%d", i), func(t *testing.T) {
			detected := utils.DetectPassword(pattern.content)

			if pattern.shouldDetect && !detected {
				t.Errorf("Expected to detect password in: %s", pattern.content)
			}

			if !pattern.shouldDetect && detected {
				t.Errorf("False positive password detection in: %s", pattern.content)
			}
		})
	}
}

func testTokenDetection(t *testing.T) {
	// Test token detection patterns
	tokenPatterns := []struct {
		content      string
		shouldDetect bool
	}{
		{`token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`, true},
		{`JWT_TOKEN=Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9`, true},
		{`"access_token": "ya29.1234567890abcdef"`, true},
		{`oauth_token: "1234567890abcdef-1234567890abcdef"`, true},
		{`token = "token"`, false},           // Too generic
		{`token = ""`, false},                // Empty
		{`// token: your_token_here`, false}, // Comment
	}

	for i, pattern := range tokenPatterns {
		t.Run(fmt.Sprintf("token_pattern_%d", i), func(t *testing.T) {
			detected := utils.DetectToken(pattern.content)

			if pattern.shouldDetect && !detected {
				t.Errorf("Expected to detect token in: %s", pattern.content)
			}

			if !pattern.shouldDetect && detected {
				t.Errorf("False positive token detection in: %s", pattern.content)
			}
		})
	}
}

func testFalsePositivePrevention(t *testing.T) {
	// Test content that should NOT trigger secret detection
	falsePositives := []string{
		`// Example: API_KEY=your_api_key_here`,
		`const placeholder = "YOUR_SECRET_HERE";`,
		`password = "password"; // Default password`,
		`token = "TOKEN_PLACEHOLDER";`,
		`apiKey = ""; // Set your API key here`,
		`const example = "sk-example1234567890abcdef";`,
		`# Configuration template`,
		`export API_KEY=your_key_here`,
		`<API_KEY>your_api_key</API_KEY>`,
		`{{.ApiKey}}`, // Template variable
	}

	for i, content := range falsePositives {
		t.Run(fmt.Sprintf("false_positive_%d", i), func(t *testing.T) {
			detected := utils.DetectSecrets(content)

			if len(detected) > 0 {
				t.Errorf("False positive secret detection in: %s, detected: %v", content, detected)
			}
		})
	}
}

// TestConcurrentSecurity tests thread safety and concurrent access security
func TestConcurrentSecurity(t *testing.T) {
	t.Run("concurrent_file_access", func(t *testing.T) {
		testConcurrentFileAccess(t)
	})

	t.Run("concurrent_validation", func(t *testing.T) {
		testConcurrentValidationSecurity(t)
	})

	t.Run("race_condition_prevention", func(t *testing.T) {
		testRaceConditionPrevention(t)
	})
}

func testConcurrentFileAccess(t *testing.T) {
	tempDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tempDir, "concurrent_test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test concurrent file access
	numGoroutines := 10
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			// Attempt to read file safely
			content, err := utils.SafeReadFile(testFile)
			if err != nil {
				errors <- fmt.Errorf("goroutine %d failed to read file: %v", id, err)
				return
			}

			if string(content) != "test content" {
				errors <- fmt.Errorf("goroutine %d got unexpected content: %s", id, string(content))
				return
			}

			errors <- nil
		}(i)
	}

	// Check for errors
	for i := 0; i < numGoroutines; i++ {
		if err := <-errors; err != nil {
			t.Error(err)
		}
	}
}

func testConcurrentValidationSecurity(t *testing.T) {
	tempDir := t.TempDir()

	// Create test project
	createSecureTestProject(t, tempDir)

	// Test concurrent validation
	numGoroutines := 5
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			config := &models.ProjectConfig{
				Name:         fmt.Sprintf("concurrent-project-%d", id),
				Organization: "test-org",
				License:      "MIT",
				OutputPath:   filepath.Join(tempDir, fmt.Sprintf("output-%d", id)),
			}

			err := utils.ValidateProjectConfig(config)
			if err != nil {
				errors <- fmt.Errorf("goroutine %d validation failed: %v", id, err)
				return
			}

			errors <- nil
		}(i)
	}

	// Check for errors
	for i := 0; i < numGoroutines; i++ {
		if err := <-errors; err != nil {
			t.Error(err)
		}
	}
}

func testRaceConditionPrevention(t *testing.T) {
	// Test that shared resources are properly protected
	sharedCounter := &utils.SafeCounter{}

	numGoroutines := 100
	incrementsPerGoroutine := 100

	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < incrementsPerGoroutine; j++ {
				sharedCounter.Increment()
			}
			errors <- nil
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		if err := <-errors; err != nil {
			t.Error(err)
		}
	}

	// Verify final count
	expectedCount := numGoroutines * incrementsPerGoroutine
	actualCount := sharedCounter.Value()

	if actualCount != expectedCount {
		t.Errorf("Race condition detected: expected count %d, got %d", expectedCount, actualCount)
	}
}

// Helper functions

func createSecureTestProject(t *testing.T, projectDir string) {
	err := os.MkdirAll(projectDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create secure package.json
	packageJSON := `{
		"name": "secure-test-project",
		"version": "1.0.0",
		"description": "A secure test project",
		"license": "MIT",
		"dependencies": {
			"express": "^4.18.0"
		}
	}`

	err = os.WriteFile(filepath.Join(projectDir, "package.json"), []byte(packageJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Create secure main file
	mainJS := `const express = require('express');
const app = express();

// Security middleware
app.use(express.json({ limit: '10mb' }));

app.get('/', (req, res) => {
	res.json({ message: 'Secure Hello World' });
});

const port = process.env.PORT || 3000;
app.listen(port, () => {
	console.log('Secure server running on port ' + port);
});
`

	err = os.WriteFile(filepath.Join(projectDir, "main.js"), []byte(mainJS), 0644)
	if err != nil {
		t.Fatalf("Failed to create main.js: %v", err)
	}
}
