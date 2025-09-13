package models

// SecurityTestFixtures provides standardized test data for security tests
type SecurityTestFixtures struct{}

// NewSecurityTestFixtures creates a new instance of security test fixtures
func NewSecurityTestFixtures() *SecurityTestFixtures {
	return &SecurityTestFixtures{}
}

// StandardSecurityConfig returns a standard security configuration for tests
func (stf *SecurityTestFixtures) StandardSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		TempFileRandomLength: 16,
		AllowedTempDirs:      []string{"/tmp", "/var/tmp"},
		FilePermissions:      0o600,
		EnablePathValidation: true,
		MaxFileSize:          10 * 1024 * 1024, // 10MB
		SecureCleanup:        true,
	}
}

// MinimalSecurityConfig returns a minimal valid security configuration for tests
func (stf *SecurityTestFixtures) MinimalSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		TempFileRandomLength: 8, // Minimum valid
		AllowedTempDirs:      []string{"/tmp"},
		FilePermissions:      0o600,
		MaxFileSize:          1024, // 1KB minimum
	}
}

// StandardRandomConfig returns a standard random configuration for tests
func (stf *SecurityTestFixtures) StandardRandomConfig() *RandomConfig {
	return &RandomConfig{
		DefaultSuffixLength: 16,
		IDFormat:            "hex",
		MinEntropyBytes:     32,
		IDPrefixLength:      4,
		EnableEntropyCheck:  true,
	}
}

// MinimalRandomConfig returns a minimal valid random configuration for tests
func (stf *SecurityTestFixtures) MinimalRandomConfig() *RandomConfig {
	return &RandomConfig{
		DefaultSuffixLength: 8, // Minimum valid
		IDFormat:            "hex",
		MinEntropyBytes:     16, // Minimum valid
		IDPrefixLength:      0,
		EnableEntropyCheck:  false,
	}
}

// DangerousDirectories returns a list of directories that should be flagged as dangerous
func (stf *SecurityTestFixtures) DangerousDirectories() []string {
	return []string{
		"/",
		"/bin",
		"/sbin",
		"/usr/bin",
		"/etc",
		"/root",
		"/boot",
	}
}

// SafeDirectories returns a list of directories that should be considered safe
func (stf *SecurityTestFixtures) SafeDirectories() []string {
	return []string{
		"/tmp",
		"/var/tmp",
		"/home/user/temp",
		"/opt/app/tmp",
	}
}

// MockFileSystemFactory is a placeholder for creating mock file systems in tests
// The actual implementation is in the test files where MockFileSystem is defined
type MockFileSystemFactory struct{}

// NewMockFileSystemFactory creates a new mock file system factory
func NewMockFileSystemFactory() *MockFileSystemFactory {
	return &MockFileSystemFactory{}
}

// SecurityTestCase represents a standardized security test case
type SecurityTestCase struct {
	Name        string
	Config      interface{} // Can be *SecurityConfig or *RandomConfig
	ExpectValid bool
	ExpectError string
	Description string
}

// SecurityValidationTestSuite provides a complete test suite for security validation
type SecurityValidationTestSuite struct {
	Fixtures *SecurityTestFixtures
}

// NewSecurityValidationTestSuite creates a new security validation test suite
func NewSecurityValidationTestSuite() *SecurityValidationTestSuite {
	return &SecurityValidationTestSuite{
		Fixtures: NewSecurityTestFixtures(),
	}
}

// StandardSecurityConfigTestCases returns standard test cases for security configuration validation
func (svts *SecurityValidationTestSuite) StandardSecurityConfigTestCases() []SecurityTestCase {
	return []SecurityTestCase{
		{
			Name:        "valid standard config",
			Config:      svts.Fixtures.StandardSecurityConfig(),
			ExpectValid: true,
			Description: "Standard configuration should pass validation",
		},
		{
			Name:        "valid minimal config",
			Config:      svts.Fixtures.MinimalSecurityConfig(),
			ExpectValid: true,
			Description: "Minimal valid configuration should pass validation",
		},
		{
			Name: "invalid temp file length",
			Config: &SecurityConfig{
				TempFileRandomLength: 4, // Too small
				AllowedTempDirs:      []string{"/tmp"},
				FilePermissions:      0o600,
				MaxFileSize:          1024,
			},
			ExpectValid: false,
			ExpectError: "min",
			Description: "Configuration with insufficient random length should fail",
		},
		{
			Name: "empty temp dirs",
			Config: &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{}, // Empty
				FilePermissions:      0o600,
				MaxFileSize:          1024,
			},
			ExpectValid: false,
			ExpectError: "min",
			Description: "Configuration with no temp directories should fail",
		},
		{
			Name: "dangerous directory",
			Config: &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{"/etc"}, // Dangerous
				FilePermissions:      0o600,
				MaxFileSize:          1024,
			},
			ExpectValid: false,
			ExpectError: "dangerous_path",
			Description: "Configuration with dangerous directory should fail",
		},
	}
}

// StandardRandomConfigTestCases returns standard test cases for random configuration validation
func (svts *SecurityValidationTestSuite) StandardRandomConfigTestCases() []SecurityTestCase {
	return []SecurityTestCase{
		{
			Name:        "valid standard config",
			Config:      svts.Fixtures.StandardRandomConfig(),
			ExpectValid: true,
			Description: "Standard random configuration should pass validation",
		},
		{
			Name:        "valid minimal config",
			Config:      svts.Fixtures.MinimalRandomConfig(),
			ExpectValid: true,
			Description: "Minimal valid random configuration should pass validation",
		},
		{
			Name: "invalid suffix length",
			Config: &RandomConfig{
				DefaultSuffixLength: 4, // Too small
				IDFormat:            "hex",
				MinEntropyBytes:     32,
			},
			ExpectValid: false,
			ExpectError: "min",
			Description: "Configuration with insufficient suffix length should fail",
		},
		{
			Name: "invalid ID format",
			Config: &RandomConfig{
				DefaultSuffixLength: 16,
				IDFormat:            "invalid", // Invalid
				MinEntropyBytes:     32,
			},
			ExpectValid: false,
			ExpectError: "oneof",
			Description: "Configuration with invalid ID format should fail",
		},
		{
			Name: "insufficient entropy",
			Config: &RandomConfig{
				DefaultSuffixLength: 16,
				IDFormat:            "hex",
				MinEntropyBytes:     8, // Too small
			},
			ExpectValid: false,
			ExpectError: "min",
			Description: "Configuration with insufficient entropy should fail",
		},
	}
}

// PermissionTestCases returns standard test cases for file permission validation
func (svts *SecurityValidationTestSuite) PermissionTestCases() []SecurityTestCase {
	return []SecurityTestCase{
		{
			Name: "secure permissions (600)",
			Config: &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{"/tmp"},
				FilePermissions:      0o600,
				MaxFileSize:          1024,
			},
			ExpectValid: true,
			Description: "Secure file permissions should pass validation",
		},
		{
			Name: "readable permissions (644)",
			Config: &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{"/tmp"},
				FilePermissions:      0o644,
				MaxFileSize:          1024,
			},
			ExpectValid: true,
			Description: "Read-only permissions should pass with warnings",
		},
		{
			Name: "no owner write (400)",
			Config: &SecurityConfig{
				TempFileRandomLength: 16,
				AllowedTempDirs:      []string{"/tmp"},
				FilePermissions:      0o400,
				MaxFileSize:          1024,
			},
			ExpectValid: false,
			ExpectError: "no_write_permission",
			Description: "Missing owner write permission should fail",
		},
	}
}
