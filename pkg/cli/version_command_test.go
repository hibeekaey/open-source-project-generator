package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCLIInterface is a mock implementation of CLIInterface for testing
type MockCLIInterface struct {
	mock.Mock
}

func (m *MockCLIInterface) GetVersionManager() interfaces.VersionManager {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(interfaces.VersionManager)
}

func (m *MockCLIInterface) GetBuildInfo() (version, gitCommit, buildTime string) {
	args := m.Called()
	return args.String(0), args.String(1), args.String(2)
}

// Add missing interface methods to satisfy CLIInterface
func (m *MockCLIInterface) Run(args []string) error {
	mockArgs := m.Called(args)
	return mockArgs.Error(0)
}

func (m *MockCLIInterface) PromptProjectDetails() (*models.ProjectConfig, error) {
	args := m.Called()
	return args.Get(0).(*models.ProjectConfig), args.Error(1)
}

func (m *MockCLIInterface) ConfirmGeneration(config *models.ProjectConfig) bool {
	args := m.Called(config)
	return args.Bool(0)
}

func (m *MockCLIInterface) AuditProject(path string, options interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	args := m.Called(path, options)
	return args.Get(0).(*interfaces.AuditResult), args.Error(1)
}

func (m *MockCLIInterface) ValidateProject(path string, options interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	args := m.Called(path, options)
	return args.Get(0).(*interfaces.ValidationResult), args.Error(1)
}

func (m *MockCLIInterface) GenerateFromConfig(path string, options interfaces.GenerateOptions) error {
	args := m.Called(path, options)
	return args.Error(0)
}

func (m *MockCLIInterface) ListTemplates(filter interfaces.TemplateFilter) ([]interfaces.TemplateInfo, error) {
	args := m.Called(filter)
	return args.Get(0).([]interfaces.TemplateInfo), args.Error(1)
}

func (m *MockCLIInterface) ShowConfig() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) SetConfig(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockCLIInterface) ShowVersion(options interfaces.VersionOptions) error {
	args := m.Called(options)
	return args.Error(0)
}

func (m *MockCLIInterface) CheckUpdates() (*interfaces.UpdateInfo, error) {
	args := m.Called()
	return args.Get(0).(*interfaces.UpdateInfo), args.Error(1)
}

func (m *MockCLIInterface) ShowCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) ClearCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) ShowLogs() error {
	args := m.Called()
	return args.Error(0)
}

// Add more missing interface methods
func (m *MockCLIInterface) PromptAdvancedOptions() (*interfaces.AdvancedOptions, error) {
	args := m.Called()
	return args.Get(0).(*interfaces.AdvancedOptions), args.Error(1)
}

func (m *MockCLIInterface) ConfirmAdvancedGeneration(config *models.ProjectConfig, options *interfaces.AdvancedOptions) bool {
	args := m.Called(config, options)
	return args.Bool(0)
}

func (m *MockCLIInterface) SelectTemplateInteractively(filter interfaces.TemplateFilter) (*interfaces.TemplateInfo, error) {
	args := m.Called(filter)
	return args.Get(0).(*interfaces.TemplateInfo), args.Error(1)
}

func (m *MockCLIInterface) GenerateWithAdvancedOptions(config *models.ProjectConfig, options *interfaces.AdvancedOptions) error {
	args := m.Called(config, options)
	return args.Error(0)
}

func (m *MockCLIInterface) ValidateProjectAdvanced(path string, options *interfaces.ValidationOptions) (*interfaces.ValidationResult, error) {
	args := m.Called(path, options)
	return args.Get(0).(*interfaces.ValidationResult), args.Error(1)
}

func (m *MockCLIInterface) AuditProjectAdvanced(path string, options *interfaces.AuditOptions) (*interfaces.AuditResult, error) {
	args := m.Called(path, options)
	return args.Get(0).(*interfaces.AuditResult), args.Error(1)
}

func (m *MockCLIInterface) GetTemplateInfo(name string) (*interfaces.TemplateInfo, error) {
	args := m.Called(name)
	return args.Get(0).(*interfaces.TemplateInfo), args.Error(1)
}

func (m *MockCLIInterface) ValidateTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	args := m.Called(path)
	return args.Get(0).(*interfaces.TemplateValidationResult), args.Error(1)
}

func (m *MockCLIInterface) SearchTemplates(query string) ([]interfaces.TemplateInfo, error) {
	args := m.Called(query)
	return args.Get(0).([]interfaces.TemplateInfo), args.Error(1)
}

func (m *MockCLIInterface) GetTemplateMetadata(name string) (*interfaces.TemplateMetadata, error) {
	args := m.Called(name)
	return args.Get(0).(*interfaces.TemplateMetadata), args.Error(1)
}

func (m *MockCLIInterface) GetTemplateDependencies(name string) ([]string, error) {
	args := m.Called(name)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCLIInterface) ValidateCustomTemplate(path string) (*interfaces.TemplateValidationResult, error) {
	args := m.Called(path)
	return args.Get(0).(*interfaces.TemplateValidationResult), args.Error(1)
}

func (m *MockCLIInterface) EditConfig() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) ValidateConfig() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) ExportConfig(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockCLIInterface) LoadConfiguration(sources []string) (*models.ProjectConfig, error) {
	args := m.Called(sources)
	return args.Get(0).(*models.ProjectConfig), args.Error(1)
}

func (m *MockCLIInterface) MergeConfigurations(configs []*models.ProjectConfig) (*models.ProjectConfig, error) {
	args := m.Called(configs)
	return args.Get(0).(*models.ProjectConfig), args.Error(1)
}

func (m *MockCLIInterface) ValidateConfigurationSchema(config *models.ProjectConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockCLIInterface) GetConfigurationSources() ([]interfaces.ConfigSource, error) {
	args := m.Called()
	return args.Get(0).([]interfaces.ConfigSource), args.Error(1)
}

func (m *MockCLIInterface) InstallUpdates() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) GetPackageVersions() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockCLIInterface) GetLatestPackageVersions() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockCLIInterface) CheckCompatibility(path string) (*interfaces.CompatibilityResult, error) {
	args := m.Called(path)
	return args.Get(0).(*interfaces.CompatibilityResult), args.Error(1)
}

func (m *MockCLIInterface) CleanCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) GetCacheStats() (*interfaces.CacheStats, error) {
	args := m.Called()
	return args.Get(0).(*interfaces.CacheStats), args.Error(1)
}

func (m *MockCLIInterface) ValidateCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) RepairCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) EnableOfflineMode() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) DisableOfflineMode() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockCLIInterface) SetLogLevel(level string) error {
	args := m.Called(level)
	return args.Error(0)
}

func (m *MockCLIInterface) GetLogLevel() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockCLIInterface) ShowRecentLogs(lines int, level string) error {
	args := m.Called(lines, level)
	return args.Error(0)
}

func (m *MockCLIInterface) GetLogFileLocations() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCLIInterface) RunNonInteractive(config *models.ProjectConfig, options *interfaces.AdvancedOptions) error {
	args := m.Called(config, options)
	return args.Error(0)
}

func (m *MockCLIInterface) GenerateReport(reportType string, format string, outputFile string) error {
	args := m.Called(reportType, format, outputFile)
	return args.Error(0)
}

func (m *MockCLIInterface) GetExitCode() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockCLIInterface) SetExitCode(code int) {
	m.Called(code)
}

// MockVersionManager is a mock implementation of VersionManager for testing
type MockVersionManager struct {
	mock.Mock
}

func (m *MockVersionManager) GetCurrentVersion() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockVersionManager) GetLatestVersion() (*interfaces.VersionInfo, error) {
	args := m.Called()
	return args.Get(0).(*interfaces.VersionInfo), args.Error(1)
}

func (m *MockVersionManager) CheckForUpdates() (*interfaces.UpdateInfo, error) {
	args := m.Called()
	return args.Get(0).(*interfaces.UpdateInfo), args.Error(1)
}

func (m *MockVersionManager) DownloadUpdate(version string) error {
	args := m.Called(version)
	return args.Error(0)
}

func (m *MockVersionManager) InstallUpdate(version string) error {
	args := m.Called(version)
	return args.Error(0)
}

func (m *MockVersionManager) ValidateVersion(version string) error {
	args := m.Called(version)
	return args.Error(0)
}

func (m *MockVersionManager) CompareVersions(v1, v2 string) int {
	args := m.Called(v1, v2)
	return args.Int(0)
}

func (m *MockVersionManager) IsNewerVersion(current, latest string) bool {
	args := m.Called(current, latest)
	return args.Bool(0)
}

// Add missing interface methods
func (m *MockVersionManager) GetLatestNodeVersion() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockVersionManager) GetLatestGoVersion() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockVersionManager) GetLatestNPMPackage(packageName string) (string, error) {
	args := m.Called(packageName)
	return args.String(0), args.Error(1)
}

func (m *MockVersionManager) GetLatestGoModule(moduleName string) (string, error) {
	args := m.Called(moduleName)
	return args.String(0), args.Error(1)
}

func (m *MockVersionManager) UpdateVersionsConfig() (*models.VersionConfig, error) {
	args := m.Called()
	return args.Get(0).(*models.VersionConfig), args.Error(1)
}

func (m *MockVersionManager) GetLatestGitHubRelease(owner, repo string) (string, error) {
	args := m.Called(owner, repo)
	return args.String(0), args.Error(1)
}

func (m *MockVersionManager) GetAllPackageVersions() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockVersionManager) GetLatestPackageVersions() (map[string]string, error) {
	args := m.Called()
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockVersionManager) GetDetailedVersionHistory() ([]interfaces.VersionInfo, error) {
	args := m.Called()
	return args.Get(0).([]interfaces.VersionInfo), args.Error(1)
}

func (m *MockVersionManager) GetUpdateChannel() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockVersionManager) SetUpdateChannel(channel string) error {
	args := m.Called(channel)
	return args.Error(0)
}

func (m *MockVersionManager) CheckCompatibility(projectPath string) (*interfaces.CompatibilityResult, error) {
	args := m.Called(projectPath)
	return args.Get(0).(*interfaces.CompatibilityResult), args.Error(1)
}

func (m *MockVersionManager) GetSupportedVersions() (map[string][]string, error) {
	args := m.Called()
	return args.Get(0).(map[string][]string), args.Error(1)
}

func (m *MockVersionManager) ValidateVersionRequirements(requirements map[string]string) (*interfaces.VersionValidationResult, error) {
	args := m.Called(requirements)
	return args.Get(0).(*interfaces.VersionValidationResult), args.Error(1)
}

func (m *MockVersionManager) CacheVersionInfo(info *interfaces.VersionInfo) error {
	args := m.Called(info)
	return args.Error(0)
}

func (m *MockVersionManager) GetCachedVersionInfo() (*interfaces.VersionInfo, error) {
	args := m.Called()
	return args.Get(0).(*interfaces.VersionInfo), args.Error(1)
}

func (m *MockVersionManager) RefreshVersionCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockVersionManager) ClearVersionCache() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockVersionManager) GetReleaseNotes(version string) (*interfaces.ReleaseNotes, error) {
	args := m.Called(version)
	return args.Get(0).(*interfaces.ReleaseNotes), args.Error(1)
}

func (m *MockVersionManager) GetChangeLog(fromVersion, toVersion string) (*interfaces.ChangeLog, error) {
	args := m.Called(fromVersion, toVersion)
	return args.Get(0).(*interfaces.ChangeLog), args.Error(1)
}

func (m *MockVersionManager) GetSecurityAdvisories(version string) ([]interfaces.SecurityAdvisory, error) {
	args := m.Called(version)
	return args.Get(0).([]interfaces.SecurityAdvisory), args.Error(1)
}

func (m *MockVersionManager) GetPackageInfo(packageName string) (*interfaces.PackageInfo, error) {
	args := m.Called(packageName)
	return args.Get(0).(*interfaces.PackageInfo), args.Error(1)
}

func (m *MockVersionManager) GetPackageVersions(packageName string) ([]string, error) {
	args := m.Called(packageName)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockVersionManager) CheckPackageUpdates(packages map[string]string) (map[string]interfaces.PackageUpdate, error) {
	args := m.Called(packages)
	return args.Get(0).(map[string]interfaces.PackageUpdate), args.Error(1)
}

func (m *MockVersionManager) SetVersionConfig(config *interfaces.VersionConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockVersionManager) GetVersionConfig() (*interfaces.VersionConfig, error) {
	args := m.Called()
	return args.Get(0).(*interfaces.VersionConfig), args.Error(1)
}

func (m *MockVersionManager) SetAutoUpdate(enabled bool) error {
	args := m.Called(enabled)
	return args.Error(0)
}

func (m *MockVersionManager) SetUpdateNotifications(enabled bool) error {
	args := m.Called(enabled)
	return args.Error(0)
}

func (m *MockVersionManager) GetVersionHistory(packageName string) ([]string, error) {
	args := m.Called(packageName)
	return args.Get(0).([]string), args.Error(1)
}

// TestNewVersionCommand tests version command creation
func TestNewVersionCommand(t *testing.T) {
	mockCLI := &MockCLIInterface{}
	cmd := NewVersionCommand(mockCLI)

	assert.NotNil(t, cmd)
	assert.Equal(t, "version", cmd.Use)
	assert.Equal(t, "Display version information", cmd.Short)

	// Check that flags are properly set up
	jsonFlag := cmd.Flags().Lookup("json")
	assert.NotNil(t, jsonFlag)
	assert.Equal(t, "false", jsonFlag.DefValue)

	formatFlag := cmd.Flags().Lookup("format")
	assert.NotNil(t, formatFlag)

	outputFormatFlag := cmd.Flags().Lookup("output-format")
	assert.NotNil(t, outputFormatFlag)

	shortFlag := cmd.Flags().Lookup("short")
	assert.NotNil(t, shortFlag)
	assert.Equal(t, "false", shortFlag.DefValue)
}

// TestRunVersionCommand tests version command execution
func TestRunVersionCommand(t *testing.T) {
	tests := []struct {
		name           string
		setupFlags     func(*cobra.Command)
		setupMocks     func(*MockCLIInterface, *MockVersionManager)
		expectedOutput string
		expectedError  bool
		errorContains  string
		validateJSON   bool
	}{
		{
			name: "text output default",
			setupFlags: func(cmd *cobra.Command) {
				// No flags set - should default to text output
			},
			setupMocks: func(mockCLI *MockCLIInterface, mockVM *MockVersionManager) {
				mockCLI.On("GetVersionManager").Return(mockVM)
				mockVM.On("GetCurrentVersion").Return("1.0.0")
			},
			expectedOutput: "1.0.0\n",
			expectedError:  false,
		},
		{
			name: "short flag output",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("short", "true")
			},
			setupMocks: func(mockCLI *MockCLIInterface, mockVM *MockVersionManager) {
				mockCLI.On("GetVersionManager").Return(mockVM)
				mockVM.On("GetCurrentVersion").Return("1.0.0")
			},
			expectedOutput: "1.0.0\n",
			expectedError:  false,
		},
		{
			name: "JSON output via json flag",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("json", "true")
			},
			setupMocks: func(mockCLI *MockCLIInterface, mockVM *MockVersionManager) {
				mockCLI.On("GetVersionManager").Return(mockVM)
				mockVM.On("GetCurrentVersion").Return("1.0.0")
				mockCLI.On("GetBuildInfo").Return("1.0.0", "abc123", "2023-01-01T00:00:00Z")
			},
			expectedError: false,
			validateJSON:  true,
		},
		{
			name: "JSON output via format flag",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("format", "json")
			},
			setupMocks: func(mockCLI *MockCLIInterface, mockVM *MockVersionManager) {
				mockCLI.On("GetVersionManager").Return(mockVM)
				mockVM.On("GetCurrentVersion").Return("1.0.0")
				mockCLI.On("GetBuildInfo").Return("1.0.0", "abc123", "2023-01-01T00:00:00Z")
			},
			expectedError: false,
			validateJSON:  true,
		},
		{
			name: "JSON output via output-format flag",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("output-format", "json")
			},
			setupMocks: func(mockCLI *MockCLIInterface, mockVM *MockVersionManager) {
				mockCLI.On("GetVersionManager").Return(mockVM)
				mockVM.On("GetCurrentVersion").Return("1.0.0")
				mockCLI.On("GetBuildInfo").Return("1.0.0", "abc123", "2023-01-01T00:00:00Z")
			},
			expectedError: false,
			validateJSON:  true,
		},
		{
			name: "nil version manager",
			setupFlags: func(cmd *cobra.Command) {
				// Default text output
			},
			setupMocks: func(mockCLI *MockCLIInterface, mockVM *MockVersionManager) {
				mockCLI.On("GetVersionManager").Return(nil)
			},
			expectedOutput: "dev\n",
			expectedError:  false,
		},
		{
			name: "JSON output with nil version manager",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("json", "true")
			},
			setupMocks: func(mockCLI *MockCLIInterface, mockVM *MockVersionManager) {
				mockCLI.On("GetVersionManager").Return(nil)
				mockCLI.On("GetBuildInfo").Return("dev", "unknown", "unknown")
			},
			expectedError: false,
			validateJSON:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			mockCLI := &MockCLIInterface{}
			mockVM := &MockVersionManager{}
			tt.setupMocks(mockCLI, mockVM)

			cmd := &cobra.Command{}
			cmd.Flags().Bool("json", false, "")
			cmd.Flags().String("format", "", "")
			cmd.Flags().String("output-format", "", "")
			cmd.Flags().Bool("short", false, "")
			tt.setupFlags(cmd)

			err := RunVersionCommand(cmd, []string{}, mockCLI)

			// Restore stdout and capture output
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)

				if tt.validateJSON {
					// Validate that output is valid JSON
					var versionInfo VersionInfo
					err := json.Unmarshal([]byte(output), &versionInfo)
					assert.NoError(t, err, "Output should be valid JSON")

					// Validate JSON structure
					assert.NotEmpty(t, versionInfo.Version)
					assert.NotEmpty(t, versionInfo.GitCommit)
					assert.NotEmpty(t, versionInfo.BuildTime)
					assert.NotEmpty(t, versionInfo.GoVersion)
					assert.NotEmpty(t, versionInfo.Platform)
					assert.NotEmpty(t, versionInfo.Architecture)
				} else if tt.expectedOutput != "" {
					assert.Equal(t, tt.expectedOutput, output)
				}
			}

			mockCLI.AssertExpectations(t)
			mockVM.AssertExpectations(t)
		})
	}
}

// TestBuildVersionInfo tests version info structure building
func TestBuildVersionInfo(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		setupMocks  func(*MockCLIInterface)
		expectedErr bool
	}{
		{
			name:    "complete build info",
			version: "1.0.0",
			setupMocks: func(mockCLI *MockCLIInterface) {
				mockCLI.On("GetBuildInfo").Return("1.0.0", "abc123", "2023-01-01T00:00:00Z")
			},
			expectedErr: false,
		},
		{
			name:    "missing build info",
			version: "1.0.0",
			setupMocks: func(mockCLI *MockCLIInterface) {
				mockCLI.On("GetBuildInfo").Return("1.0.0", "", "")
			},
			expectedErr: false,
		},
		{
			name:        "nil CLI",
			version:     "1.0.0",
			setupMocks:  nil, // No mock setup, will pass nil CLI
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mockCLI *MockCLIInterface
			if tt.setupMocks != nil {
				mockCLI = &MockCLIInterface{}
				tt.setupMocks(mockCLI)
			}

			versionInfo, err := buildVersionInfo(tt.version, mockCLI)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, versionInfo)
				assert.Equal(t, tt.version, versionInfo.Version)
				assert.NotEmpty(t, versionInfo.GoVersion)
				assert.NotEmpty(t, versionInfo.Platform)
				assert.NotEmpty(t, versionInfo.Architecture)

				// Check default values for missing build info
				if mockCLI == nil {
					assert.Equal(t, "unknown", versionInfo.GitCommit)
					assert.Equal(t, "unknown", versionInfo.BuildTime)
				}
			}

			if mockCLI != nil {
				mockCLI.AssertExpectations(t)
			}
		})
	}
}

// TestVersionInfoJSONMarshaling tests JSON marshaling of VersionInfo
func TestVersionInfoJSONMarshaling(t *testing.T) {
	versionInfo := &VersionInfo{
		Version:      "1.0.0",
		GitCommit:    "abc123",
		BuildTime:    "2023-01-01T00:00:00Z",
		GoVersion:    "go1.21.0",
		Platform:     "linux",
		Architecture: "amd64",
	}

	jsonBytes, err := json.MarshalIndent(versionInfo, "", "  ")
	assert.NoError(t, err)

	// Validate that the JSON is valid
	assert.True(t, json.Valid(jsonBytes))

	// Unmarshal and verify structure
	var unmarshaled VersionInfo
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	assert.NoError(t, err)

	assert.Equal(t, versionInfo.Version, unmarshaled.Version)
	assert.Equal(t, versionInfo.GitCommit, unmarshaled.GitCommit)
	assert.Equal(t, versionInfo.BuildTime, unmarshaled.BuildTime)
	assert.Equal(t, versionInfo.GoVersion, unmarshaled.GoVersion)
	assert.Equal(t, versionInfo.Platform, unmarshaled.Platform)
	assert.Equal(t, versionInfo.Architecture, unmarshaled.Architecture)
}

// TestVersionInfoJSONTags tests that JSON tags are properly set
func TestVersionInfoJSONTags(t *testing.T) {
	versionInfo := &VersionInfo{
		Version:      "1.0.0",
		GitCommit:    "abc123",
		BuildTime:    "2023-01-01T00:00:00Z",
		GoVersion:    "go1.21.0",
		Platform:     "linux",
		Architecture: "amd64",
	}

	jsonBytes, err := json.Marshal(versionInfo)
	assert.NoError(t, err)

	// Parse as generic map to check field names
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	assert.NoError(t, err)

	// Verify JSON field names match tags
	expectedFields := []string{
		"version", "git_commit", "build_time",
		"go_version", "platform", "architecture",
	}

	for _, field := range expectedFields {
		assert.Contains(t, jsonMap, field, "JSON should contain field %s", field)
	}
}

// TestVersionCommandEnhancedErrorHandling tests enhanced error handling scenarios
func TestVersionCommandEnhancedErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*MockCLIInterface)
		setupFlags    func(*cobra.Command)
		expectedError bool
		errorContains string
	}{
		{
			name: "build info error simulation",
			setupMocks: func(mockCLI *MockCLIInterface) {
				mockCLI.On("GetVersionManager").Return(nil)
				// Simulate a scenario where GetBuildInfo might cause issues
				mockCLI.On("GetBuildInfo").Return("", "", "")
			},
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("json", "true")
			},
			expectedError: false, // Should handle gracefully with defaults
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCLI := &MockCLIInterface{}
			tt.setupMocks(mockCLI)

			cmd := &cobra.Command{}
			cmd.Flags().Bool("json", false, "")
			cmd.Flags().String("format", "", "")
			cmd.Flags().String("output-format", "", "")
			cmd.Flags().Bool("short", false, "")
			tt.setupFlags(cmd)

			err := RunVersionCommand(cmd, []string{}, mockCLI)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			mockCLI.AssertExpectations(t)
		})
	}
}

// TestVersionCommandFlagPrecedence tests flag precedence for output format
func TestVersionCommandFlagPrecedence(t *testing.T) {
	tests := []struct {
		name        string
		setupFlags  func(*cobra.Command)
		expectJSON  bool
		description string
	}{
		{
			name: "json flag takes precedence",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("json", "true")
				cmd.Flags().Set("format", "text")
			},
			expectJSON:  true,
			description: "json flag should override format flag",
		},
		{
			name: "format flag works when json not set",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("format", "json")
			},
			expectJSON:  true,
			description: "format=json should enable JSON output",
		},
		{
			name: "output-format flag works when others not set",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("output-format", "json")
			},
			expectJSON:  true,
			description: "output-format=json should enable JSON output",
		},
		{
			name: "no JSON flags defaults to text",
			setupFlags: func(cmd *cobra.Command) {
				// No JSON-related flags set
			},
			expectJSON:  false,
			description: "should default to text output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCLI := &MockCLIInterface{}
			mockVM := &MockVersionManager{}

			mockCLI.On("GetVersionManager").Return(mockVM)
			mockVM.On("GetCurrentVersion").Return("1.0.0")

			if tt.expectJSON {
				mockCLI.On("GetBuildInfo").Return("1.0.0", "abc123", "2023-01-01T00:00:00Z")
			}

			cmd := &cobra.Command{}
			cmd.Flags().Bool("json", false, "")
			cmd.Flags().String("format", "", "")
			cmd.Flags().String("output-format", "", "")
			cmd.Flags().Bool("short", false, "")
			tt.setupFlags(cmd)

			// Capture output to verify format
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := RunVersionCommand(cmd, []string{}, mockCLI)
			assert.NoError(t, err, tt.description)

			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if tt.expectJSON {
				// Should be valid JSON
				var versionInfo VersionInfo
				err := json.Unmarshal([]byte(output), &versionInfo)
				assert.NoError(t, err, "Output should be valid JSON for test: %s", tt.description)
			} else {
				// Should be simple text
				assert.Equal(t, "1.0.0\n", output, "Output should be simple text for test: %s", tt.description)
			}

			mockCLI.AssertExpectations(t)
			mockVM.AssertExpectations(t)
		})
	}
}
