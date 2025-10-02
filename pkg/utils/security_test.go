package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		allowedBasePaths []string
		wantError        bool
		errorContains    string
	}{
		{"empty path", "", nil, true, "path cannot be empty"},
		{"valid relative path", "test/path", nil, false, ""},
		{"path traversal", "../../../etc/passwd", nil, true, "path traversal detected"},
		{"absolute dangerous path", "/etc/passwd", nil, true, "access to system path denied"},
		{"windows dangerous path", "C:\\Windows\\System32", nil, true, "access to system path denied"},
		{"URI scheme", "http://example.com", nil, true, "URI schemes not allowed"},
		{"valid with allowed base", "project/src/main.go", []string{"project"}, false, ""},
		{"invalid with allowed base", "other/src/main.go", []string{"project"}, true, "path not within allowed directories"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.path, tt.allowedBasePaths...)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidatePath() error = %v, wantError %v", err, tt.wantError)
			}
			if err != nil && tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("ValidatePath() error = %v, want to contain %v", err.Error(), tt.errorContains)
			}
		})
	}
}

func TestSafeReadFile(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("test content")

	err := os.WriteFile(testFile, testContent, 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name             string
		path             string
		allowedBasePaths []string
		wantError        bool
		wantContent      []byte
	}{
		{"valid file", testFile, []string{tmpDir}, false, testContent},
		{"invalid path", "../../../etc/passwd", nil, true, nil},
		{"non-existent file", filepath.Join(tmpDir, "nonexistent.txt"), []string{tmpDir}, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := SafeReadFile(tt.path, tt.allowedBasePaths...)
			if (err != nil) != tt.wantError {
				t.Errorf("SafeReadFile() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && string(content) != string(tt.wantContent) {
				t.Errorf("SafeReadFile() content = %v, want %v", string(content), string(tt.wantContent))
			}
		})
	}
}

func TestSafeWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("test content")

	tests := []struct {
		name             string
		path             string
		data             []byte
		allowedBasePaths []string
		wantError        bool
	}{
		{"valid write", testFile, testContent, []string{tmpDir}, false},
		{"invalid path", "../../../tmp/malicious.txt", testContent, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SafeWriteFile(tt.path, tt.data, tt.allowedBasePaths...)
			if (err != nil) != tt.wantError {
				t.Errorf("SafeWriteFile() error = %v, wantError %v", err, tt.wantError)
			}

			// Verify file was written correctly if no error expected
			if !tt.wantError {
				content, readErr := os.ReadFile(tt.path)
				if readErr != nil {
					t.Errorf("Failed to read written file: %v", readErr)
				}
				if string(content) != string(tt.data) {
					t.Errorf("Written content = %v, want %v", string(content), string(tt.data))
				}
			}
		})
	}
}

func TestSafeMkdirAll(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "test", "nested", "dir")

	tests := []struct {
		name             string
		path             string
		allowedBasePaths []string
		wantError        bool
	}{
		{"valid directory", testDir, []string{tmpDir}, false},
		{"invalid path", "../../../tmp/malicious", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SafeMkdirAll(tt.path, tt.allowedBasePaths...)
			if (err != nil) != tt.wantError {
				t.Errorf("SafeMkdirAll() error = %v, wantError %v", err, tt.wantError)
			}

			// Verify directory was created if no error expected
			if !tt.wantError {
				if _, statErr := os.Stat(tt.path); os.IsNotExist(statErr) {
					t.Errorf("Directory was not created: %v", tt.path)
				}
			}
		})
	}
}

func TestSafeOpenFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	tests := []struct {
		name             string
		path             string
		flag             int
		perm             os.FileMode
		allowedBasePaths []string
		wantError        bool
	}{
		{"create new file", testFile, os.O_CREATE | os.O_WRONLY, 0644, []string{tmpDir}, false},
		{"invalid path", "../../../tmp/malicious.txt", os.O_CREATE | os.O_WRONLY, 0644, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := SafeOpenFile(tt.path, tt.flag, tt.perm, tt.allowedBasePaths...)
			if (err != nil) != tt.wantError {
				t.Errorf("SafeOpenFile() error = %v, wantError %v", err, tt.wantError)
			}
			if file != nil {
				file.Close()
			}
		})
	}
}

func TestSafeOpen(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// Create test file
	err := os.WriteFile(testFile, []byte("test"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name             string
		path             string
		allowedBasePaths []string
		wantError        bool
	}{
		{"valid file", testFile, []string{tmpDir}, false},
		{"invalid path", "../../../etc/passwd", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := SafeOpen(tt.path, tt.allowedBasePaths...)
			if (err != nil) != tt.wantError {
				t.Errorf("SafeOpen() error = %v, wantError %v", err, tt.wantError)
			}
			if file != nil {
				file.Close()
			}
		})
	}
}

func TestSafeCreate(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	tests := []struct {
		name             string
		path             string
		allowedBasePaths []string
		wantError        bool
	}{
		{"valid create", testFile, []string{tmpDir}, false},
		{"invalid path", "../../../tmp/malicious.txt", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := SafeCreate(tt.path, tt.allowedBasePaths...)
			if (err != nil) != tt.wantError {
				t.Errorf("SafeCreate() error = %v, wantError %v", err, tt.wantError)
			}
			if file != nil {
				file.Close()
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantError     bool
		wantResult    string
		errorContains string
	}{
		{"empty input", "", true, "", "input cannot be empty"},
		{"valid input", "test-project_1.0", false, "test-project_1.0", ""},
		{"input with whitespace", "  test-project  ", false, "test-project", ""},
		{"too long input", strings.Repeat("a", 256), true, "", "input too long"},
		{"dangerous script", "<script>alert('xss')</script>", true, "", "dangerous pattern"},
		{"path traversal", "../../../etc", true, "", "dangerous pattern"},
		{"invalid characters", "test@project!", true, "", "invalid characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizeInput(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("SanitizeInput() error = %v, wantError %v", err, tt.wantError)
			}
			if !tt.wantError && result != tt.wantResult {
				t.Errorf("SanitizeInput() result = %v, want %v", result, tt.wantResult)
			}
			if err != nil && tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("SanitizeInput() error = %v, want to contain %v", err.Error(), tt.errorContains)
			}
		})
	}
}

func TestValidateFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files with different permissions
	normalFile := filepath.Join(tmpDir, "normal.txt")
	err := os.WriteFile(normalFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create normal file: %v", err)
	}

	worldWritableFile := filepath.Join(tmpDir, "world_writable.txt")
	err = os.WriteFile(worldWritableFile, []byte("test"), 0666)
	if err != nil {
		t.Fatalf("Failed to create world writable file: %v", err)
	}
	// Make it actually world writable
	err = os.Chmod(worldWritableFile, 0666)
	if err != nil {
		t.Fatalf("Failed to chmod world writable file: %v", err)
	}

	noPermFile := filepath.Join(tmpDir, "no_perm.txt")
	err = os.WriteFile(noPermFile, []byte("test"), 0000)
	if err != nil {
		t.Fatalf("Failed to create no permission file: %v", err)
	}

	tests := []struct {
		name          string
		path          string
		wantError     bool
		errorContains string
	}{
		{"normal file", normalFile, false, ""},
		{"world writable", worldWritableFile, true, "world writable"},
		{"no permissions", noPermFile, true, "no permissions"},
		{"non-existent file", filepath.Join(tmpDir, "nonexistent.txt"), true, "failed to get file info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilePermissions(tt.path)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateFilePermissions() error = %v, wantError %v", err, tt.wantError)
			}
			if err != nil && tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("ValidateFilePermissions() error = %v, want to contain %v", err.Error(), tt.errorContains)
			}
		})
	}
}

func TestProcessTemplateSafely(t *testing.T) {
	tmpDir := t.TempDir()

	// Create safe template
	safeTemplate := filepath.Join(tmpDir, "safe.tmpl")
	err := os.WriteFile(safeTemplate, []byte("Hello {{.Name}}!"), 0644)
	if err != nil {
		t.Fatalf("Failed to create safe template: %v", err)
	}

	// Create dangerous template
	dangerousTemplate := filepath.Join(tmpDir, "dangerous.tmpl")
	err = os.WriteFile(dangerousTemplate, []byte("{{exec \"rm -rf /\"}}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create dangerous template: %v", err)
	}

	config := &models.ProjectConfig{Name: "test"}

	tests := []struct {
		name          string
		templatePath  string
		config        *models.ProjectConfig
		wantError     bool
		errorContains string
	}{
		{"safe template", safeTemplate, config, false, ""},
		{"dangerous template", dangerousTemplate, config, true, "dangerous pattern"},
		{"invalid path", "../../../etc/passwd", config, true, "invalid template path"},
		{"non-existent template", filepath.Join(tmpDir, "nonexistent.tmpl"), config, true, "failed to read template"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ProcessTemplateSafely(tt.templatePath, tt.config)
			if (err != nil) != tt.wantError {
				t.Errorf("ProcessTemplateSafely() error = %v, wantError %v", err, tt.wantError)
			}
			if err != nil && tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("ProcessTemplateSafely() error = %v, want to contain %v", err.Error(), tt.errorContains)
			}
		})
	}
}

func TestDetectAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"no API key", "This is normal content", false},
		{"OpenAI API key", "sk-1234567890abcdef1234567890abcdef", true},
		{"Stripe API key", "pk_live_1234567890abcdef1234567890abcdef", true},
		{"Google API key", "AIzaSyAbCdEfGhIjKlMnOpQrStUvWxYz123456789", true},
		{"GitHub token", "ghp_1234567890abcdef1234567890abcdef123456", true},
		{"API key in config", "API_KEY = abc123def456ghi789jkl012mno345pqr678", true},
		{"API key in JSON", `"api_key": "abc123def456ghi789jkl012mno345pqr678"`, true},
		{"placeholder API key", "API_KEY = your_api_key_here", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectAPIKey(tt.content); got != tt.want {
				t.Errorf("DetectAPIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectPassword(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"no password", "This is normal content", false},
		{"password in config", `password = "secretpassword123"`, true},
		{"PASSWORD env var", "PASSWORD = secretpass123", true},
		{"password in JSON", `"password": "mysecretpass"`, true},
		{"pwd field", `pwd: "mypassword"`, true},
		{"const password", `const pass = "secretpass123"`, true},
		{"common password", `password = "password"`, false},
		{"placeholder password", `password = "your_password_here"`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectPassword(tt.content); got != tt.want {
				t.Errorf("DetectPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectToken(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"no token", "This is normal content", false},
		{"JWT token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", true},
		{"Google OAuth token", "ya29.1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234", true},
		{"generic token", "abc123def456ghi789jkl012mno345pqr678-stu901vwx234yz567abc123def456ghi789", true},
		{"token in config", `token = "abc123def456ghi789jkl012mno345pqr678"`, true},
		{"JWT token env", "JWT_TOKEN = Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", true},
		{"access token JSON", `"access_token": "abc123def456ghi789jkl012mno345pqr678"`, true},
		{"oauth token", `oauth_token: "abc123def456ghi789jkl012mno345pqr678"`, true},
		{"placeholder token", `token = "your_token_here"`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectToken(tt.content); got != tt.want {
				t.Errorf("DetectToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectSecrets(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{"no secrets", "This is normal content", []string{}},
		{"API key only", "sk-1234567890abcdef1234567890abcdef", []string{"api_key"}},
		{"password only", `password = "secretpassword123"`, []string{"password"}},
		{"token only", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", []string{"token"}},
		{"multiple secrets", `sk-1234567890abcdef1234567890abcdef password = "secretpass123"`, []string{"api_key", "password"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectSecrets(tt.content)
			if len(got) != len(tt.want) {
				t.Errorf("DetectSecrets() = %v, want %v", got, tt.want)
				return
			}
			for i, secret := range got {
				if secret != tt.want[i] {
					t.Errorf("DetectSecrets() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestSafeCounter(t *testing.T) {
	counter := &SafeCounter{}

	// Test initial value
	if counter.Value() != 0 {
		t.Errorf("Initial counter value = %v, want 0", counter.Value())
	}

	// Test increment
	counter.Increment()
	if counter.Value() != 1 {
		t.Errorf("Counter value after increment = %v, want 1", counter.Value())
	}

	// Test multiple increments
	for i := 0; i < 10; i++ {
		counter.Increment()
	}
	if counter.Value() != 11 {
		t.Errorf("Counter value after 11 increments = %v, want 11", counter.Value())
	}
}

func TestSafeMap(t *testing.T) {
	safeMap := NewSafeMap()

	// Test initial size
	if safeMap.Size() != 0 {
		t.Errorf("Initial map size = %v, want 0", safeMap.Size())
	}

	// Test set and get
	safeMap.Set("key1", "value1")
	value, exists := safeMap.Get("key1")
	if !exists || value != "value1" {
		t.Errorf("Get() = %v, %v, want value1, true", value, exists)
	}

	// Test size after set
	if safeMap.Size() != 1 {
		t.Errorf("Map size after set = %v, want 1", safeMap.Size())
	}

	// Test get non-existent key
	value, exists = safeMap.Get("nonexistent")
	if exists || value != "" {
		t.Errorf("Get() for non-existent key = %v, %v, want '', false", value, exists)
	}

	// Test multiple sets
	safeMap.Set("key2", "value2")
	safeMap.Set("key3", "value3")
	if safeMap.Size() != 3 {
		t.Errorf("Map size after multiple sets = %v, want 3", safeMap.Size())
	}
}
