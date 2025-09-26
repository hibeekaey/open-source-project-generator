package security

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSecurityManagerIntegration(t *testing.T) {
	// Create a temporary workspace for testing
	tempDir, err := os.MkdirTemp("", "security_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create security manager
	manager := NewSecurityManager(tempDir)

	t.Run("Input Sanitization", func(t *testing.T) {
		// Test dangerous input sanitization
		result := manager.SanitizeString("<script>alert('xss')</script>", "test_input")
		if result.IsValid {
			t.Error("Expected dangerous input to be invalid")
		}
		if len(result.Errors) == 0 {
			t.Error("Expected errors for dangerous input")
		}

		// Test project name sanitization
		projectResult := manager.SanitizeProjectName("My Project Name!")
		if !projectResult.IsValid {
			t.Error("Expected valid project name after sanitization")
		}
		if projectResult.Sanitized != "my-project-name" {
			t.Errorf("Expected 'my-project-name', got '%s'", projectResult.Sanitized)
		}
	})

	t.Run("File Operations with Backup", func(t *testing.T) {
		// Enable backup
		manager.SetBackupEnabled(true)

		// Create a test file
		testFile := filepath.Join(tempDir, "test.txt")
		originalContent := []byte("original content")

		writeResult, err := manager.SecureWriteFile(testFile, originalContent, 0600)
		if err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
		if !writeResult.Success {
			t.Error("Expected successful file write")
		}

		// Modify the file (should create backup)
		newContent := []byte("modified content")
		writeResult2, err := manager.SecureWriteFile(testFile, newContent, 0600)
		if err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
		if !writeResult2.Success {
			t.Error("Expected successful file write")
		}

		// Check that backup was created
		backups, err := manager.ListBackups(testFile)
		if err != nil {
			t.Fatalf("Failed to list backups: %v", err)
		}
		if len(backups) == 0 {
			t.Error("Expected backup to be created")
		}
	})

	t.Run("Dry Run Mode", func(t *testing.T) {
		// Enable dry run mode
		manager.SetDryRunMode(true)

		// Attempt file operations
		testFile := filepath.Join(tempDir, "dryrun_test.txt")
		testContent := []byte("test content")

		result, err := manager.SecureWriteFile(testFile, testContent, 0600)
		if err != nil {
			t.Fatalf("Dry run should not fail: %v", err)
		}
		if !result.Success {
			t.Error("Expected dry run to succeed")
		}

		// File should not actually exist
		if _, err := os.Stat(testFile); !os.IsNotExist(err) {
			t.Error("File should not exist in dry run mode")
		}

		// Check dry run operations were recorded
		operations := manager.GetDryRunOperations()
		if len(operations) == 0 {
			t.Error("Expected dry run operations to be recorded")
		}

		// Generate dry run report
		report, err := manager.GenerateDryRunReport("text")
		if err != nil {
			t.Fatalf("Failed to generate dry run report: %v", err)
		}
		if report == "" {
			t.Error("Expected non-empty dry run report")
		}

		// Disable dry run mode
		manager.SetDryRunMode(false)
	})

	t.Run("Template Security", func(t *testing.T) {
		// Test dangerous template content
		dangerousTemplate := `{{exec "rm -rf /"}}`
		result := manager.ValidateTemplateContent(dangerousTemplate, "test.tmpl")
		if result.IsSecure {
			t.Error("Expected dangerous template to be flagged as insecure")
		}
		if len(result.SecurityIssues) == 0 {
			t.Error("Expected security issues to be found")
		}

		// Test safe template content
		safeTemplate := `Hello {{.Name}}!`
		safeResult := manager.ValidateTemplateContent(safeTemplate, "safe.tmpl")
		if !safeResult.IsSecure {
			t.Error("Expected safe template to be secure")
		}
	})

	t.Run("Configuration Management", func(t *testing.T) {
		// Test security configuration
		config := map[string]interface{}{
			"max_file_size":   int64(50 * 1024 * 1024), // 50MB
			"backup_enabled":  true,
			"dry_run_enabled": false,
			"non_interactive": true,
		}

		err := manager.SetSecurityConfig(config)
		if err != nil {
			t.Fatalf("Failed to set security config: %v", err)
		}

		// Verify configuration was applied
		retrievedConfig := manager.GetSecurityConfig()
		if retrievedConfig["backup_enabled"] != true {
			t.Error("Expected backup_enabled to be true")
		}
		if retrievedConfig["non_interactive"] != true {
			t.Error("Expected non_interactive to be true")
		}
	})

	t.Run("Path Validation", func(t *testing.T) {
		// Test valid path
		validPath := filepath.Join(tempDir, "valid", "file.txt")
		err := manager.ValidateFilePath(validPath, "write")
		if err != nil {
			t.Errorf("Expected valid path to pass validation: %v", err)
		}

		// Test path traversal attempt
		maliciousPath := filepath.Join(tempDir, "..", "..", "etc", "passwd")
		err = manager.ValidateFilePath(maliciousPath, "write")
		if err == nil {
			t.Error("Expected path traversal to be blocked")
		}
	})
}

func TestSecurityManagerNonInteractiveMode(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "security_test_ni")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	manager := NewSecurityManager(tempDir)
	manager.SetNonInteractive(true)

	// Test confirmation in non-interactive mode
	result, err := manager.ConfirmFileOverwrite("test.txt", 1024)
	if err != nil {
		t.Fatalf("Non-interactive confirmation should not fail: %v", err)
	}
	if result.NonInteractive != true {
		t.Error("Expected non-interactive flag to be set")
	}
	if result.DefaultUsed != true {
		t.Error("Expected default answer to be used")
	}
}

func TestSecurityManagerPerformance(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "security_perf_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	manager := NewSecurityManager(tempDir)

	// Test performance of input sanitization
	start := time.Now()
	for i := 0; i < 1000; i++ {
		manager.SanitizeString("test input string", "test")
	}
	duration := time.Since(start)

	if duration > time.Second {
		t.Errorf("Input sanitization too slow: %v for 1000 operations", duration)
	}

	// Test performance of path validation
	start = time.Now()
	for i := 0; i < 1000; i++ {
		_ = manager.ValidateFilePath(filepath.Join(tempDir, "test.txt"), "read")
	}
	duration = time.Since(start)

	if duration > time.Second {
		t.Errorf("Path validation too slow: %v for 1000 operations", duration)
	}
}
