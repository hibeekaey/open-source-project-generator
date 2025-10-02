package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSecurityManager(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	assert.NotNil(t, manager)

	// Verify default configuration
	config := manager.GetSecurityConfig()
	assert.Equal(t, tempDir, config["workspace_dir"])
	assert.Equal(t, filepath.Join(tempDir, ".generator", "backups"), config["backup_dir"])
	assert.Equal(t, true, config["backup_enabled"])
	assert.Equal(t, false, config["dry_run_enabled"])
	assert.Equal(t, false, config["non_interactive"])
	assert.Equal(t, int64(100*1024*1024), config["max_file_size"])
	assert.Equal(t, 10, config["max_backups"])
}

func TestSecurityManager_InputSanitization(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	// Test string sanitization
	result := manager.SanitizeString("valid string", "test")
	assert.True(t, result.IsValid)
	assert.Equal(t, "valid string", result.Sanitized)

	result = manager.SanitizeString("<script>alert('xss')</script>", "test")
	assert.False(t, result.IsValid)
	assert.Greater(t, len(result.Errors), 0)

	// Test project name sanitization
	result = manager.SanitizeProjectName("My Project Name!")
	assert.True(t, result.IsValid)
	assert.True(t, result.WasModified)
	assert.Equal(t, "my-project-name", result.Sanitized)

	// Test file path sanitization
	result = manager.SanitizeFilePath("valid/path.txt", "file_path")
	assert.True(t, result.IsValid)

	result = manager.SanitizeFilePath("../../../etc/passwd", "file_path")
	assert.False(t, result.IsValid)

	// Test URL sanitization
	result = manager.SanitizeURL("https://example.com", "url")
	assert.True(t, result.IsValid)

	result = manager.SanitizeURL("javascript:alert('xss')", "url")
	assert.False(t, result.IsValid)

	// Test email sanitization
	result = manager.SanitizeEmail("user@example.com", "email")
	assert.True(t, result.IsValid)

	result = manager.SanitizeEmail("invalid-email", "email")
	assert.False(t, result.IsValid)
}

func TestSecurityManager_ValidateAndSanitizeMap(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	testData := map[string]interface{}{
		"name":        "valid-name",
		"email":       "user@example.com",
		"description": "A valid description",
		"dangerous":   "<script>alert('xss')</script>",
		"nested": map[string]interface{}{
			"value": "nested-value",
			"bad":   "javascript:alert('xss')",
		},
	}

	results := manager.ValidateAndSanitizeMap(testData, "")

	// Check valid fields
	assert.True(t, results["name"].IsValid)
	assert.True(t, results["email"].IsValid)
	assert.True(t, results["description"].IsValid)

	// Check dangerous field
	assert.False(t, results["dangerous"].IsValid)

	// Check nested dangerous field
	assert.False(t, results["nested.bad"].IsValid)
}

func TestSecurityManager_TemplateSecurityIntegration(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	// Test safe template content
	safeContent := "Hello {{.Name}}!"
	result := manager.ValidateTemplateContent(safeContent, "safe.tmpl")
	assert.True(t, result.IsSecure)
	assert.Empty(t, result.SecurityIssues)

	// Test dangerous template content
	dangerousContent := "{{exec \"rm -rf /\"}}"
	result = manager.ValidateTemplateContent(dangerousContent, "dangerous.tmpl")
	assert.False(t, result.IsSecure)
	assert.Greater(t, len(result.SecurityIssues), 0)

	// Test template file validation
	templateFile := filepath.Join(tempDir, "test.tmpl")
	err := os.WriteFile(templateFile, []byte(safeContent), 0644)
	require.NoError(t, err)

	result, err = manager.ValidateTemplateFile(templateFile)
	assert.NoError(t, err)
	assert.True(t, result.IsSecure)

	// Test template directory scanning
	templateDir := filepath.Join(tempDir, "templates")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)

	template1 := filepath.Join(templateDir, "template1.tmpl")
	template2 := filepath.Join(templateDir, "template2.tmpl")

	err = os.WriteFile(template1, []byte(safeContent), 0644)
	require.NoError(t, err)
	err = os.WriteFile(template2, []byte(dangerousContent), 0644)
	require.NoError(t, err)

	results, err := manager.ScanTemplateDirectory(templateDir)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.True(t, results[template1].IsSecure)
	assert.False(t, results[template2].IsSecure)
}

func TestSecurityManager_FileOperationsIntegration(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	testContent := []byte("test content")

	// Test secure write
	result, err := manager.SecureWriteFile(testFile, testContent, 0644)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "write", result.Operation)

	// Test secure read
	result, data, err := manager.SecureReadFile(testFile)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, testContent, data)

	// Test secure copy
	copyFile := filepath.Join(tempDir, "copy.txt")
	result, err = manager.SecureCopyFile(testFile, copyFile)
	assert.NoError(t, err)
	assert.True(t, result.Success)

	// Test secure directory creation
	testDir := filepath.Join(tempDir, "testdir")
	result, err = manager.SecureCreateDirectory(testDir, 0755)
	assert.NoError(t, err)
	assert.True(t, result.Success)

	// Test path validation
	err = manager.ValidateFilePath(testFile, "read")
	assert.NoError(t, err)

	err = manager.ValidateFilePath("../../../etc/passwd", "read")
	assert.Error(t, err)

	// Test file permissions
	perms, err := manager.GetFilePermissions(testFile)
	assert.NoError(t, err)
	assert.NotNil(t, perms)
	assert.Contains(t, perms, "mode")
}

func TestSecurityManager_BackupIntegration(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	testFile := filepath.Join(tempDir, "backup_test.txt")
	originalContent := []byte("original content")

	// Create original file
	err := os.WriteFile(testFile, originalContent, 0644)
	require.NoError(t, err)

	// Test backup creation
	result, err := manager.BackupFile(testFile)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotEmpty(t, result.BackupPath)

	// Test backup listing
	backups, err := manager.ListBackups(testFile)
	assert.NoError(t, err)
	assert.Len(t, backups, 1)

	// Test backup restoration
	modifiedContent := []byte("modified content")
	err = os.WriteFile(testFile, modifiedContent, 0644)
	require.NoError(t, err)

	result, err = manager.RestoreFile(testFile, backups[0].Timestamp)
	assert.NoError(t, err)
	assert.True(t, result.Success)

	// Verify restoration
	restoredContent, err := os.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, originalContent, restoredContent)

	// Test backup directory
	sourceDir := filepath.Join(tempDir, "source")
	err = os.MkdirAll(sourceDir, 0755)
	require.NoError(t, err)

	file1 := filepath.Join(sourceDir, "file1.txt")
	file2 := filepath.Join(sourceDir, "file2.txt")
	err = os.WriteFile(file1, []byte("content1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file2, []byte("content2"), 0644)
	require.NoError(t, err)

	results, err := manager.BackupDirectory(sourceDir)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// Test backup enable/disable
	assert.True(t, manager.IsBackupEnabled())
	manager.SetBackupEnabled(false)
	assert.False(t, manager.IsBackupEnabled())
}

func TestSecurityManager_DryRunIntegration(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	// Enable dry run mode
	manager.SetDryRunMode(true)
	assert.True(t, manager.IsDryRunMode())

	testFile := filepath.Join(tempDir, "dryrun_test.txt")
	testContent := []byte("test content")

	// Test dry run file operations
	result, err := manager.SecureWriteFile(testFile, testContent, 0644)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Contains(t, result.Operation, "dry-run")

	// File should not actually exist
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))

	// Test dry run operations recording
	manager.RecordFileWrite(testFile, testContent, false)
	manager.RecordDirectoryCreate(filepath.Join(tempDir, "testdir"))
	manager.RecordFileCopy("src", "dst", false)
	manager.RecordTemplateProcess("template.tmpl", "output.txt", map[string]interface{}{"key": "value"})

	operations := manager.GetDryRunOperations()
	assert.Greater(t, len(operations), 0)

	summary := manager.GetDryRunSummary()
	assert.Greater(t, summary["total_operations"].(int), 0)

	// Test report generation
	report, err := manager.GenerateDryRunReport("text")
	assert.NoError(t, err)
	assert.Contains(t, report, "DRY RUN REPORT")

	// Disable dry run mode
	manager.SetDryRunMode(false)
	assert.False(t, manager.IsDryRunMode())
}

func TestSecurityManager_ConfirmationIntegration(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	// Test in non-interactive mode
	manager.SetNonInteractive(true)
	assert.True(t, manager.IsNonInteractive())

	// Test file overwrite confirmation
	result, err := manager.ConfirmFileOverwrite("/test/file.txt", 1024)
	assert.NoError(t, err)
	assert.True(t, result.NonInteractive)

	// Test directory delete confirmation
	result, err = manager.ConfirmDirectoryDelete("/test/dir", 5, 2048)
	assert.NoError(t, err)
	assert.True(t, result.NonInteractive)

	// Test bulk operation confirmation
	result, err = manager.ConfirmBulkOperation("process", 50, []string{"detail1", "detail2"})
	assert.NoError(t, err)
	assert.True(t, result.NonInteractive)

	// Test security risk confirmation
	result, err = manager.ConfirmSecurityRisk("Test risk", "medium", []string{"detail"})
	assert.NoError(t, err)
	assert.True(t, result.NonInteractive)

	// Test dry run confirmation
	dryRunSummary := map[string]interface{}{
		"total_operations":       10,
		"safe_operations":        8,
		"warning_operations":     2,
		"destructive_operations": 0,
		"total_files_affected":   5,
	}

	result, err = manager.ConfirmWithDryRun(dryRunSummary)
	assert.NoError(t, err)
	assert.True(t, result.NonInteractive)
}

func TestSecurityManager_ConfigurationManagement(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	// Test setting security configuration
	config := map[string]interface{}{
		"max_file_size":   int64(50 * 1024 * 1024), // 50MB
		"backup_enabled":  false,
		"dry_run_enabled": true,
		"non_interactive": true,
		"max_backups":     5,
		"default_answer":  true,
	}

	err := manager.SetSecurityConfig(config)
	assert.NoError(t, err)

	// Verify configuration was applied
	retrievedConfig := manager.GetSecurityConfig()
	assert.Equal(t, int64(50*1024*1024), retrievedConfig["max_file_size"])
	assert.Equal(t, false, retrievedConfig["backup_enabled"])
	assert.Equal(t, true, retrievedConfig["dry_run_enabled"])
	assert.Equal(t, true, retrievedConfig["non_interactive"])
	assert.Equal(t, 5, retrievedConfig["max_backups"])

	// Verify components were configured
	assert.False(t, manager.IsBackupEnabled())
	assert.True(t, manager.IsDryRunMode())
	assert.True(t, manager.IsNonInteractive())
}

func TestSecurityManager_FileOperationsWithBackup(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	testFile := filepath.Join(tempDir, "backup_integration.txt")
	originalContent := []byte("original content")

	// Create original file
	err := os.WriteFile(testFile, originalContent, 0644)
	require.NoError(t, err)

	// Enable backup
	manager.SetBackupEnabled(true)

	// Write new content (should create backup)
	newContent := []byte("new content")
	result, err := manager.SecureWriteFile(testFile, newContent, 0644)
	assert.NoError(t, err)
	assert.True(t, result.Success)

	// Verify backup was created
	backups, err := manager.ListBackups(testFile)
	assert.NoError(t, err)
	assert.Greater(t, len(backups), 0)

	// Copy file (should create backup of destination if it exists)
	copyFile := filepath.Join(tempDir, "copy_with_backup.txt")
	err = os.WriteFile(copyFile, []byte("existing content"), 0644)
	require.NoError(t, err)

	result, err = manager.SecureCopyFile(testFile, copyFile)
	assert.NoError(t, err)
	assert.True(t, result.Success)

	// Verify backup was created for copy destination
	copyBackups, err := manager.ListBackups(copyFile)
	assert.NoError(t, err)
	assert.Greater(t, len(copyBackups), 0)
}

func TestSecurityManager_DryRunWithBackup(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	testFile := filepath.Join(tempDir, "dryrun_backup_test.txt")
	testContent := []byte("test content")

	// Enable both dry run and backup
	manager.SetDryRunMode(true)
	manager.SetBackupEnabled(true)

	// File operations should be recorded but not executed
	result, err := manager.SecureWriteFile(testFile, testContent, 0644)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Contains(t, result.Operation, "dry-run")

	// File should not exist
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))

	// No backups should be created
	backups, err := manager.ListBackups(testFile)
	assert.NoError(t, err)
	assert.Empty(t, backups)

	// Operations should be recorded
	operations := manager.GetDryRunOperations()
	assert.Greater(t, len(operations), 0)
}

func TestSecurityManager_ErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	// Test invalid file operations
	invalidFile := "/etc/passwd"

	result, err := manager.SecureWriteFile(invalidFile, []byte("content"), 0644)
	assert.Error(t, err)
	assert.False(t, result.Success)

	result, _, err = manager.SecureReadFile(invalidFile)
	assert.Error(t, err)
	assert.False(t, result.Success)

	result, err = manager.SecureCopyFile(invalidFile, filepath.Join(tempDir, "copy.txt"))
	assert.Error(t, err)
	assert.False(t, result.Success)

	// Test backup operations with invalid files
	backupResult, err := manager.BackupFile(invalidFile)
	assert.Error(t, err)
	assert.False(t, backupResult.Success)

	// Test template validation with invalid content
	templateResult := manager.ValidateTemplateContent("{{exec \"rm -rf /\"}}", "dangerous.tmpl")
	assert.False(t, templateResult.IsSecure)
	assert.Greater(t, len(templateResult.SecurityIssues), 0)
}

func TestSecurityManager_PerformanceAndLimits(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	// Test file size limits
	config := map[string]interface{}{
		"max_file_size": int64(1024), // 1KB limit
	}

	err := manager.SetSecurityConfig(config)
	assert.NoError(t, err)

	// Try to write file larger than limit
	largeContent := make([]byte, 2048) // 2KB
	testFile := filepath.Join(tempDir, "large.txt")

	result, err := manager.SecureWriteFile(testFile, largeContent, 0644)
	assert.Error(t, err)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "exceeds maximum allowed size")

	// Test template size limits
	largeTemplate := strings.Repeat("a", 2*1024*1024) // 2MB
	templateResult := manager.ValidateTemplateContent(largeTemplate, "large.tmpl")
	assert.False(t, templateResult.IsSecure)
	assert.Contains(t, templateResult.SecurityIssues[0], "exceeds maximum allowed size")
}

func TestSecurityManager_CompleteWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewSecurityManager(tempDir)

	// Configure security manager
	config := map[string]interface{}{
		"backup_enabled":  true,
		"dry_run_enabled": false,
		"non_interactive": true,
		"max_backups":     5,
	}

	err := manager.SetSecurityConfig(config)
	require.NoError(t, err)

	// 1. Sanitize project configuration
	projectData := map[string]interface{}{
		"name":        "My Awesome Project!",
		"description": "A great project",
		"author":      "john@example.com",
		"url":         "https://github.com/user/project",
	}

	sanitizationResults := manager.ValidateAndSanitizeMap(projectData, "project")

	// Verify sanitization
	assert.True(t, sanitizationResults["name"].IsValid)
	assert.True(t, sanitizationResults["name"].WasModified)
	assert.Equal(t, "my-awesome-project", sanitizationResults["name"].Sanitized)

	// 2. Create project directory
	projectDir := filepath.Join(tempDir, "my-awesome-project")
	result, err := manager.SecureCreateDirectory(projectDir, 0755)
	require.NoError(t, err)
	assert.True(t, result.Success)

	// 3. Validate and process templates
	templateContent := `# {{.Name | upper}}

Welcome to {{.Name}}!

Author: {{.Author}}
Description: {{.Description}}
`

	templateResult := manager.ValidateTemplateContent(templateContent, "README.tmpl")
	assert.True(t, templateResult.IsSecure)

	// 4. Write project files with backup
	readmeFile := filepath.Join(projectDir, "README.md")
	readmeContent := []byte("# MY-AWESOME-PROJECT\n\nWelcome to my-awesome-project!\n")

	result, err = manager.SecureWriteFile(readmeFile, readmeContent, 0644)
	require.NoError(t, err)
	assert.True(t, result.Success)

	// 5. Update file (should create backup)
	updatedContent := []byte("# MY-AWESOME-PROJECT\n\nUpdated welcome message!\n")
	result, err = manager.SecureWriteFile(readmeFile, updatedContent, 0644)
	require.NoError(t, err)
	assert.True(t, result.Success)

	// 6. Verify backup was created
	backups, err := manager.ListBackups(readmeFile)
	assert.NoError(t, err)
	assert.Greater(t, len(backups), 0)

	// 7. Copy files
	licenseFile := filepath.Join(projectDir, "LICENSE")
	licenseContent := []byte("MIT License\n\nCopyright (c) 2023\n")

	err = os.WriteFile(licenseFile, licenseContent, 0644)
	require.NoError(t, err)

	copyFile := filepath.Join(projectDir, "LICENSE.backup")
	result, err = manager.SecureCopyFile(licenseFile, copyFile)
	require.NoError(t, err)
	assert.True(t, result.Success)

	// 8. Validate file permissions
	perms, err := manager.GetFilePermissions(readmeFile)
	assert.NoError(t, err)
	assert.False(t, perms["is_dir"].(bool))
	assert.True(t, perms["is_readable"].(bool))

	// 9. Test confirmation for potentially dangerous operation
	confirmResult, err := manager.ConfirmDirectoryDelete(projectDir, 3, 1024)
	assert.NoError(t, err)
	assert.True(t, confirmResult.NonInteractive)

	// 10. Scan for template security issues
	templateDir := filepath.Join(projectDir, "templates")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)

	safeTemplate := filepath.Join(templateDir, "safe.tmpl")
	err = os.WriteFile(safeTemplate, []byte("Hello {{.Name}}!"), 0644)
	require.NoError(t, err)

	scanResults, err := manager.ScanTemplateDirectory(templateDir)
	assert.NoError(t, err)
	assert.Len(t, scanResults, 1)
	assert.True(t, scanResults[safeTemplate].IsSecure)
}
