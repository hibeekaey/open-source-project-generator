package standards

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewTemplateUpdater(t *testing.T) {
	updater := NewTemplateUpdater()
	if updater == nil {
		t.Fatal("NewTemplateUpdater() returned nil")
	}
}

func TestIsValidTemplateType(t *testing.T) {
	testCases := []struct {
		templateType string
		expected     bool
	}{
		{"nextjs-app", true},
		{"nextjs-home", true},
		{"nextjs-admin", true},
		{"invalid-template", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.templateType, func(t *testing.T) {
			result := isValidTemplateType(tc.templateType)
			if result != tc.expected {
				t.Errorf("isValidTemplateType(%s) = %v, expected %v", tc.templateType, result, tc.expected)
			}
		})
	}
}

func TestFindTemplateFiles(t *testing.T) {
	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "template-files-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{
		"package.json.tmpl",
		"tsconfig.json.tmpl",
		"README.md",
		"src/component.tsx.tmpl",
		"src/utils.ts",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Find template files
	files, err := findTemplateFiles(tempDir)
	if err != nil {
		t.Fatalf("findTemplateFiles failed: %v", err)
	}

	// Check results
	expectedTemplateFiles := []string{
		"package.json.tmpl",
		"tsconfig.json.tmpl",
		"src/component.tsx.tmpl",
	}

	if len(files) != len(expectedTemplateFiles) {
		t.Errorf("Expected %d template files, got %d", len(expectedTemplateFiles), len(files))
	}

	for _, expectedFile := range expectedTemplateFiles {
		found := false
		for _, file := range files {
			if strings.HasSuffix(file, expectedFile) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected template file %s not found", expectedFile)
		}
	}
}

func TestCopyDir(t *testing.T) {
	// Create source directory with test files
	srcDir, err := os.MkdirTemp("", "copy-src")
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	defer os.RemoveAll(srcDir)

	// Create test files in source
	testFiles := map[string]string{
		"file1.txt":        "content1",
		"subdir/file2.txt": "content2",
		"subdir/file3.txt": "content3",
	}

	for file, content := range testFiles {
		filePath := filepath.Join(srcDir, file)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create destination directory
	dstDir, err := os.MkdirTemp("", "copy-dst")
	if err != nil {
		t.Fatalf("Failed to create destination directory: %v", err)
	}
	defer os.RemoveAll(dstDir)

	// Copy directory
	if err := copyDir(srcDir, dstDir); err != nil {
		t.Fatalf("copyDir failed: %v", err)
	}

	// Verify copied files
	for file, expectedContent := range testFiles {
		dstFilePath := filepath.Join(dstDir, file)

		if !fileExists(dstFilePath) {
			t.Errorf("Copied file %s does not exist", file)
			continue
		}

		content, err := os.ReadFile(dstFilePath)
		if err != nil {
			t.Errorf("Failed to read copied file %s: %v", file, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("Copied file %s has incorrect content: got %s, expected %s", file, string(content), expectedContent)
		}
	}
}

func TestApplyVersionUpdates(t *testing.T) {
	// Create a temporary directory with test template files
	tempDir, err := os.MkdirTemp("", "version-update-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test template file with version placeholders
	testContent := `{
  "name": "test-app",
  "dependencies": {
    "next": "{{.Versions.NextJS}}",
    "react": "{{.Versions.React}}"
  }
}`

	testFile := filepath.Join(tempDir, "package.json.tmpl")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create updater and apply version updates
	updater := NewTemplateUpdater()
	versions := map[string]string{
		"NextJS": "15.5.3",
		"React":  "19.1.0",
	}

	if err := updater.ApplyVersionUpdates(tempDir, versions); err != nil {
		t.Fatalf("ApplyVersionUpdates failed: %v", err)
	}

	// Verify updates
	updatedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	updatedStr := string(updatedContent)
	if !strings.Contains(updatedStr, `"next": "15.5.3"`) {
		t.Error("NextJS version was not updated correctly")
	}

	if !strings.Contains(updatedStr, `"react": "19.1.0"`) {
		t.Error("React version was not updated correctly")
	}

	// Verify placeholders are replaced
	if strings.Contains(updatedStr, "{{.Versions.NextJS}}") {
		t.Error("NextJS placeholder was not replaced")
	}

	if strings.Contains(updatedStr, "{{.Versions.React}}") {
		t.Error("React placeholder was not replaced")
	}
}

func TestBackupAndRestoreTemplate(t *testing.T) {
	// Create a temporary template directory
	templateDir, err := os.MkdirTemp("", "template-backup-test")
	if err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}
	defer os.RemoveAll(templateDir)

	// Create test files in template
	testFiles := map[string]string{
		"package.json.tmpl": `{"name": "test"}`,
		"src/index.ts.tmpl": `console.log("test");`,
	}

	for file, content := range testFiles {
		filePath := filepath.Join(templateDir, file)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	updater := NewTemplateUpdater()

	// Create backup
	backupPath, err := updater.BackupTemplate(templateDir)
	if err != nil {
		t.Fatalf("BackupTemplate failed: %v", err)
	}
	defer os.RemoveAll(backupPath)

	// Verify backup exists
	if !fileExists(backupPath) {
		t.Fatal("Backup directory was not created")
	}

	// Verify backup contains all files
	for file := range testFiles {
		backupFilePath := filepath.Join(backupPath, file)
		if !fileExists(backupFilePath) {
			t.Errorf("Backup file %s does not exist", file)
		}
	}

	// Modify original template
	modifiedContent := `{"name": "modified"}`
	modifiedFile := filepath.Join(templateDir, "package.json.tmpl")
	if err := os.WriteFile(modifiedFile, []byte(modifiedContent), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Restore from backup
	if err := updater.RestoreTemplate(templateDir, backupPath); err != nil {
		t.Fatalf("RestoreTemplate failed: %v", err)
	}

	// Verify restoration
	restoredContent, err := os.ReadFile(modifiedFile)
	if err != nil {
		t.Fatalf("Failed to read restored file: %v", err)
	}

	if string(restoredContent) != testFiles["package.json.tmpl"] {
		t.Errorf("File was not restored correctly: got %s, expected %s", string(restoredContent), testFiles["package.json.tmpl"])
	}
}

func TestGetStandardizationRules(t *testing.T) {
	updater := &StandardTemplateUpdater{
		generator: NewTemplateGenerator(),
		validator: NewTemplateValidator(),
		standards: GetFrontendStandards(),
	}

	rules := updater.GetStandardizationRules()

	if len(rules) == 0 {
		t.Error("GetStandardizationRules should return non-empty slice")
	}

	// Check that all rules have required fields
	for _, rule := range rules {
		if rule.Name == "" {
			t.Error("Rule should have a name")
		}

		if rule.Description == "" {
			t.Error("Rule should have a description")
		}

		if rule.FilePattern == "" {
			t.Error("Rule should have a file pattern")
		}

		if rule.Validator == nil {
			t.Error("Rule should have a validator function")
		}

		if rule.Fixer == nil {
			t.Error("Rule should have a fixer function")
		}
	}

	// Check for expected rules
	expectedRules := []string{
		"PackageJSONScripts",
		"TypeScriptConfig",
		"ESLintConfig",
		"PrettierConfig",
		"VercelConfig",
	}

	for _, expectedRule := range expectedRules {
		found := false
		for _, rule := range rules {
			if rule.Name == expectedRule {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected rule %s not found", expectedRule)
		}
	}
}

func TestUpdaterValidatePackageJSONScripts(t *testing.T) {
	updater := &StandardTemplateUpdater{
		generator: NewTemplateGenerator(),
		validator: NewTemplateValidator(),
		standards: GetFrontendStandards(),
	}

	testCases := []struct {
		name        string
		content     string
		expectValid bool
	}{
		{
			name: "valid package.json with all scripts",
			content: `{
				"scripts": {
					"dev": "next dev",
					"build": "next build",
					"start": "next start",
					"lint": "next lint",
					"type-check": "tsc --noEmit",
					"test": "jest"
				}
			}`,
			expectValid: true,
		},
		{
			name: "missing required scripts",
			content: `{
				"scripts": {
					"dev": "next dev",
					"build": "next build"
				}
			}`,
			expectValid: false,
		},
		{
			name: "no scripts section",
			content: `{
				"name": "test-app"
			}`,
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid, issues := updater.validatePackageJSONScripts(tc.content)

			if isValid != tc.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tc.expectValid, isValid)
			}

			if !tc.expectValid && len(issues) == 0 {
				t.Error("Expected validation issues but got none")
			}

			if tc.expectValid && len(issues) > 0 {
				t.Errorf("Expected no validation issues but got: %v", issues)
			}
		})
	}
}

func TestUpdaterValidateTypeScriptConfig(t *testing.T) {
	updater := &StandardTemplateUpdater{
		generator: NewTemplateGenerator(),
		validator: NewTemplateValidator(),
		standards: GetFrontendStandards(),
	}

	testCases := []struct {
		name        string
		content     string
		expectValid bool
	}{
		{
			name: "valid tsconfig.json",
			content: `{
				"compilerOptions": {
					"strict": true,
					"noEmit": true,
					"jsx": "preserve",
					"moduleResolution": "bundler"
				}
			}`,
			expectValid: true,
		},
		{
			name: "missing required options",
			content: `{
				"compilerOptions": {
					"strict": true
				}
			}`,
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid, issues := updater.validateTypeScriptConfig(tc.content)

			if isValid != tc.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tc.expectValid, isValid)
			}

			if !tc.expectValid && len(issues) == 0 {
				t.Error("Expected validation issues but got none")
			}
		})
	}
}

func TestUpdaterValidateESLintConfig(t *testing.T) {
	updater := &StandardTemplateUpdater{
		generator: NewTemplateGenerator(),
		validator: NewTemplateValidator(),
		standards: GetFrontendStandards(),
	}

	testCases := []struct {
		name        string
		content     string
		expectValid bool
	}{
		{
			name: "valid eslint config",
			content: `{
				"extends": [
					"next/core-web-vitals",
					"@typescript-eslint/recommended"
				]
			}`,
			expectValid: true,
		},
		{
			name: "missing required extends",
			content: `{
				"extends": ["next/core-web-vitals"]
			}`,
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid, issues := updater.validateESLintConfig(tc.content)

			if isValid != tc.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tc.expectValid, isValid)
			}

			if !tc.expectValid && len(issues) == 0 {
				t.Error("Expected validation issues but got none")
			}
		})
	}
}

func TestUpdaterValidatePrettierConfig(t *testing.T) {
	updater := &StandardTemplateUpdater{
		generator: NewTemplateGenerator(),
		validator: NewTemplateValidator(),
		standards: GetFrontendStandards(),
	}

	testCases := []struct {
		name        string
		content     string
		expectValid bool
	}{
		{
			name: "valid prettier config",
			content: `{
				"semi": true,
				"singleQuote": true,
				"printWidth": 80,
				"tabWidth": 2
			}`,
			expectValid: true,
		},
		{
			name: "missing required options",
			content: `{
				"semi": true
			}`,
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid, issues := updater.validatePrettierConfig(tc.content)

			if isValid != tc.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tc.expectValid, isValid)
			}

			if !tc.expectValid && len(issues) == 0 {
				t.Error("Expected validation issues but got none")
			}
		})
	}
}

func TestUpdaterValidateVercelConfig(t *testing.T) {
	updater := &StandardTemplateUpdater{
		generator: NewTemplateGenerator(),
		validator: NewTemplateValidator(),
		standards: GetFrontendStandards(),
	}

	testCases := []struct {
		name        string
		content     string
		expectValid bool
	}{
		{
			name: "valid vercel config",
			content: `{
				"framework": "nextjs",
				"buildCommand": "npm run build"
			}`,
			expectValid: true,
		},
		{
			name: "missing required settings",
			content: `{
				"framework": "nextjs"
			}`,
			expectValid: false,
		},
		{
			name: "incorrect framework",
			content: `{
				"framework": "react",
				"buildCommand": "npm run build"
			}`,
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid, issues := updater.validateVercelConfig(tc.content)

			if isValid != tc.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tc.expectValid, isValid)
			}

			if !tc.expectValid && len(issues) == 0 {
				t.Error("Expected validation issues but got none")
			}
		})
	}
}
