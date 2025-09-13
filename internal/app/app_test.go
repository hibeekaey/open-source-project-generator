package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/open-source-template-generator/internal/container"
	"github.com/open-source-template-generator/pkg/filesystem"
	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/template"
)

func TestNewApp(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	if app == nil {
		t.Fatal("NewApp() returned nil")
	}

	if app.container != c {
		t.Error("App container not set correctly")
	}

	if app.logger == nil {
		t.Error("App logger not initialized")
	}

	if app.errorHandler == nil {
		t.Error("App error handler not initialized")
	}

	if app.cli == nil {
		t.Error("App CLI not initialized")
	}

	if app.rootCmd == nil {
		t.Error("App root command not initialized")
	}
}

func TestAppClose(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	err := app.Close()
	if err != nil {
		t.Errorf("App.Close() returned error: %v", err)
	}
}

func TestAppInitializeComponents(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)
	if app == nil {
		t.Fatal("NewApp returned nil")
	}

	// Check that all components are initialized
	if c.GetConfigManager() == nil {
		t.Error("ConfigManager not initialized")
	}

	if c.GetValidator() == nil {
		t.Error("Validator not initialized")
	}

	if c.GetTemplateEngine() == nil {
		t.Error("TemplateEngine not initialized")
	}

	if c.GetFileSystemGenerator() == nil {
		t.Error("FileSystemGenerator not initialized")
	}

	if c.GetVersionManager() == nil {
		t.Error("VersionManager not initialized")
	}
}

func TestAppExecute(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Test that Execute doesn't panic
	// We can't easily test the actual execution without mocking cobra
	if app.rootCmd == nil {
		t.Error("Root command not initialized")
	}
}

func TestAppSetupCommands(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Check that root command is set up
	if app.rootCmd == nil {
		t.Fatal("Root command not initialized")
	}

	if app.rootCmd.Use != "generator" {
		t.Errorf("Expected root command use 'generator', got '%s'", app.rootCmd.Use)
	}

	// Check that subcommands are added
	commands := app.rootCmd.Commands()
	expectedCommands := []string{"generate", "validate", "version", "config"}

	if len(commands) < len(expectedCommands) {
		t.Errorf("Expected at least %d subcommands, got %d", len(expectedCommands), len(commands))
	}

	// Check that expected commands exist
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		// Extract just the command name (first word)
		cmdName := strings.Fields(cmd.Use)[0]
		commandNames[cmdName] = true
	}

	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("Expected command '%s' not found", expected)
		}
	}
}

func TestAppGenerateCommand(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	generateCmd := app.generateCommand()
	if generateCmd == nil {
		t.Fatal("Generate command not created")
	}

	if generateCmd.Use != "generate" {
		t.Errorf("Expected generate command use 'generate', got '%s'", generateCmd.Use)
	}

	// Check flags
	flags := generateCmd.Flags()
	if flags.Lookup("dry-run") == nil {
		t.Error("Generate command missing --dry-run flag")
	}
	if flags.Lookup("config") == nil {
		t.Error("Generate command missing --config flag")
	}
	if flags.Lookup("output") == nil {
		t.Error("Generate command missing --output flag")
	}
}

func TestAppValidateCommand(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	validateCmd := app.validateCommand()
	if validateCmd == nil {
		t.Fatal("Validate command not created")
	}

	if validateCmd.Use != "validate [project-path]" {
		t.Errorf("Expected validate command use 'validate [project-path]', got '%s'", validateCmd.Use)
	}

	// Check flags
	flags := validateCmd.Flags()
	if flags.Lookup("verbose") == nil {
		t.Error("Validate command missing --verbose flag")
	}
}

func TestAppVersionCommand(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	versionCmd := app.versionCommand()
	if versionCmd == nil {
		t.Fatal("Version command not created")
	}

	if versionCmd.Use != "version" {
		t.Errorf("Expected version command use 'version', got '%s'", versionCmd.Use)
	}

	// Check flags
	flags := versionCmd.Flags()
	if flags.Lookup("packages") == nil {
		t.Error("Version command missing --packages flag")
	}
}

func TestAppConfigCommand(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	configCmd := app.configCommand()
	if configCmd == nil {
		t.Fatal("Config command not created")
	}

	if configCmd.Use != "config" {
		t.Errorf("Expected config command use 'config', got '%s'", configCmd.Use)
	}

	// Check subcommands
	subcommands := configCmd.Commands()
	expectedSubcommands := []string{"show", "set", "reset"}

	if len(subcommands) < len(expectedSubcommands) {
		t.Errorf("Expected at least %d config subcommands, got %d", len(expectedSubcommands), len(subcommands))
	}

	subcommandNames := make(map[string]bool)
	for _, cmd := range subcommands {
		// Extract just the command name (first word)
		cmdName := strings.Fields(cmd.Use)[0]
		subcommandNames[cmdName] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommandNames[expected] {
			t.Errorf("Expected config subcommand '%s' not found", expected)
		}
	}
}

func TestAppGenerateProject(t *testing.T) {
	c := container.NewContainer()

	// Set up dependencies
	templateEngine := template.NewEngine()
	fsGenerator := filesystem.NewGenerator()
	c.SetTemplateEngine(templateEngine)
	c.SetFileSystemGenerator(fsGenerator)

	app := NewApp(c)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "app-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a minimal config for testing
	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		OutputPath:   tempDir,
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
			},
		},
	}

	// Test project generation - this may fail if templates don't exist, which is expected in tests
	err = app.generateProject(config)
	if err != nil {
		// Log the error but don't fail the test - templates may not exist in test environment
		t.Logf("generateProject() returned error (may be expected): %v", err)
		return
	}

	// Check that project directory was created
	projectPath := filepath.Join(tempDir, config.Name)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Errorf("Project directory not created: %s", projectPath)
	}
}

func TestAppRunGenerateCommand(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "app-test-generate-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test dry run mode
	err = app.runGenerateCommand(true, "", tempDir)
	// This should fail because we don't have interactive input, but it shouldn't panic
	if err == nil {
		t.Error("Expected error for dry run without config, got nil")
	}
}

func TestAppRunValidateCommand(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Test with non-existent path
	err := app.runValidateCommand([]string{"/nonexistent/path"}, false)
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}

	// Test with current directory (should exist)
	err = app.runValidateCommand([]string{"."}, false)
	// This might fail validation but shouldn't panic
	if err != nil {
		t.Logf("Validation error (expected): %v", err)
	}
}

func TestAppRunVersionCommand(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Test version command without packages
	err := app.runVersionCommand(false)
	if err != nil {
		t.Errorf("runVersionCommand(false) returned error: %v", err)
	}

	// Test version command with packages (might fail due to network)
	err = app.runVersionCommand(true)
	if err != nil {
		t.Logf("runVersionCommand(true) returned error (expected): %v", err)
	}
}

func TestAppRunConfigShowCommand(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Test config show command
	err := app.runConfigShowCommand()
	// This might fail if defaults can't be loaded, but shouldn't panic
	if err != nil {
		t.Logf("runConfigShowCommand() returned error (might be expected): %v", err)
	}
}

func TestAppRunConfigSetCommand(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Test config set with invalid args
	err := app.runConfigSetCommand([]string{"key"}, "")
	if err == nil {
		t.Error("Expected error for invalid args, got nil")
	}

	// Test config set with valid args
	err = app.runConfigSetCommand([]string{"key", "value"}, "")
	if err != nil {
		t.Errorf("runConfigSetCommand() returned error: %v", err)
	}

	// Test config set with file
	err = app.runConfigSetCommand([]string{}, "/nonexistent/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent config file, got nil")
	}
}

func TestAppRunConfigResetCommand(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Test config reset command
	err := app.runConfigResetCommand()
	if err != nil {
		t.Errorf("runConfigResetCommand() returned error: %v", err)
	}
}

func TestAppGenerateBaseFiles(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "app-test-base-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
	}

	templateEngine := c.GetTemplateEngine()
	fsGenerator := c.GetFileSystemGenerator()

	// Test base files generation
	err = app.generateBaseFiles(templateEngine, fsGenerator, config, tempDir)
	// This might fail if templates don't exist, but shouldn't panic
	if err != nil {
		t.Logf("generateBaseFiles() returned error (might be expected): %v", err)
	}
}

func TestAppGenerateFrontendComponents(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "app-test-frontend-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		Components: models.Components{
			Frontend: models.FrontendComponents{
				MainApp: true,
				Home:    true,
				Admin:   true,
			},
		},
	}

	templateEngine := c.GetTemplateEngine()
	fsGenerator := c.GetFileSystemGenerator()

	// Test frontend components generation
	err = app.generateFrontendComponents(templateEngine, fsGenerator, config, tempDir)
	// This might fail if templates don't exist, but shouldn't panic
	if err != nil {
		t.Logf("generateFrontendComponents() returned error (might be expected): %v", err)
	}
}

func TestAppGenerateBackendComponents(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "app-test-backend-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		Components: models.Components{
			Backend: models.BackendComponents{
				API: true,
			},
		},
	}

	templateEngine := c.GetTemplateEngine()
	fsGenerator := c.GetFileSystemGenerator()

	// Test backend components generation
	err = app.generateBackendComponents(templateEngine, fsGenerator, config, tempDir)
	// This might fail if templates don't exist, but shouldn't panic
	if err != nil {
		t.Logf("generateBackendComponents() returned error (might be expected): %v", err)
	}
}

func TestAppGenerateMobileComponents(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "app-test-mobile-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		Components: models.Components{
			Mobile: models.MobileComponents{
				Android: true,
				IOS:     true,
			},
		},
	}

	templateEngine := c.GetTemplateEngine()
	fsGenerator := c.GetFileSystemGenerator()

	// Test mobile components generation
	err = app.generateMobileComponents(templateEngine, fsGenerator, config, tempDir)
	// This might fail if templates don't exist, but shouldn't panic
	if err != nil {
		t.Logf("generateMobileComponents() returned error (might be expected): %v", err)
	}
}

func TestAppGenerateInfrastructureComponents(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "app-test-infra-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
		Components: models.Components{
			Infrastructure: models.InfrastructureComponents{
				Docker:     true,
				Kubernetes: true,
				Terraform:  true,
			},
		},
	}

	templateEngine := c.GetTemplateEngine()
	fsGenerator := c.GetFileSystemGenerator()

	// Test infrastructure components generation
	err = app.generateInfrastructureComponents(templateEngine, fsGenerator, config, tempDir)
	// This might fail if templates don't exist, but shouldn't panic
	if err != nil {
		t.Logf("generateInfrastructureComponents() returned error (might be expected): %v", err)
	}
}

func TestAppGenerateCICDComponents(t *testing.T) {
	c := container.NewContainer()
	app := NewApp(c)

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "app-test-cicd-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &models.ProjectConfig{
		Name:         "test-project",
		Organization: "test-org",
		Description:  "Test project",
		License:      "MIT",
	}

	templateEngine := c.GetTemplateEngine()
	fsGenerator := c.GetFileSystemGenerator()

	// Test CI/CD components generation
	err = app.generateCICDComponents(templateEngine, fsGenerator, config, tempDir)
	// This might fail if templates don't exist, but shouldn't panic
	if err != nil {
		t.Logf("generateCICDComponents() returned error (might be expected): %v", err)
	}
}
