package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/cuesoftinc/open-source-project-generator/internal/orchestrator"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

var (
	// Version information (set during build)
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"

	// Global flags
	configFile      string
	outputDir       string
	verbose         bool
	dryRun          bool
	noExternalTools bool
	createBackup    bool
	forceOverwrite  bool
	interactive     bool
	offlineMode     bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "generator",
	Short: "Open Source Project Generator - Modern tool-orchestration architecture",
	Long: `A CLI tool that generates production-ready project scaffolding across multiple
technology stacks by orchestrating industry-standard bootstrap tools.

Supports: Next.js, Go/Gin, Android/Kotlin, iOS/Swift, Docker, Kubernetes, Terraform`,
	Version: fmt.Sprintf("%s (built: %s, commit: %s)", version, buildTime, gitCommit),
}

func init() {
	// Add commands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(checkToolsCmd)
	rootCmd.AddCommand(initConfigCmd)
	rootCmd.AddCommand(cacheToolsCmd)
	rootCmd.AddCommand(versionCmd)
}

var checkToolsCmd = &cobra.Command{
	Use:   "check-tools",
	Short: "Check availability of required bootstrap tools",
	Long: `Check which bootstrap tools are available on your system and display
installation instructions for any missing tools.

This command validates your environment before project generation and helps
you install any missing dependencies.

Examples:
  # Check all registered tools
  generator check-tools

  # Check with verbose output
  generator check-tools --verbose

  # Check specific tools
  generator check-tools npx go gradle`,
	RunE: runCheckTools,
}

func init() {
	checkToolsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func runCheckTools(cmd *cobra.Command, args []string) error {
	// Initialize logger
	log := logger.NewLogger()
	if verbose {
		log.SetLevel(logger.DebugLevel)
	}

	// Create tool discovery
	toolDiscovery := orchestrator.NewToolDiscovery(log)

	// Determine which tools to check
	var toolsToCheck []string
	if len(args) > 0 {
		// Check specific tools provided as arguments
		toolsToCheck = args
	} else {
		// Check all registered tools
		toolsToCheck = toolDiscovery.ListRegisteredTools()
	}

	log.Info(fmt.Sprintf("Checking %d tools...", len(toolsToCheck)))

	// Check tool requirements
	result, err := toolDiscovery.CheckRequirements(toolsToCheck)
	if err != nil {
		return fmt.Errorf("tool check failed: %w", err)
	}

	toolCheckResult := result.(*models.ToolCheckResult)

	// Display results
	displayToolCheckResults(toolCheckResult, toolDiscovery, log)

	// Return error if not all tools are available (non-zero exit code)
	if !toolCheckResult.AllAvailable {
		return fmt.Errorf("some required tools are not available")
	}

	return nil
}

func displayToolCheckResults(result *models.ToolCheckResult, discovery *orchestrator.ToolDiscovery, log *logger.Logger) {
	fmt.Println("\n" + separator("="))
	fmt.Println("TOOL AVAILABILITY CHECK")
	fmt.Println(separator("="))

	availableCount := len(result.Tools) - len(result.Missing)
	fmt.Printf("\nStatus: %d of %d tools available\n", availableCount, len(result.Tools))
	fmt.Printf("Checked at: %s\n", result.CheckedAt.Format(time.RFC3339))

	// Display available tools
	if availableCount > 0 {
		fmt.Println("\n✓ Available Tools:")
		fmt.Println(separator("-"))
		for name, tool := range result.Tools {
			if tool.Available {
				versionInfo := ""
				if tool.InstalledVersion != "" {
					versionInfo = fmt.Sprintf(" (version: %s)", tool.InstalledVersion)
				}
				fmt.Printf("  ✓ %s%s\n", name, versionInfo)

				// Show component types this tool supports
				if metadata, err := discovery.GetToolMetadata(name); err == nil && len(metadata.ComponentTypes) > 0 {
					fmt.Printf("    Supports: %v\n", metadata.ComponentTypes)
				}
			}
		}
	}

	// Display missing tools
	if len(result.Missing) > 0 {
		fmt.Println("\n✗ Missing Tools:")
		fmt.Println(separator("-"))
		for _, name := range result.Missing {
			fmt.Printf("  ✗ %s\n", name)

			// Get metadata to show what it's used for
			if metadata, err := discovery.GetToolMetadata(name); err == nil {
				if len(metadata.ComponentTypes) > 0 {
					fmt.Printf("    Required for: %v\n", metadata.ComponentTypes)
				}
				if metadata.FallbackAvailable {
					fmt.Println("    Note: Fallback generation available")
				}
			}

			// Show installation instructions
			instructions := discovery.GetInstallInstructions(name, "")
			fmt.Printf("\n%s\n", indentText(instructions, "    "))
		}
	}

	// Display outdated tools
	if len(result.Outdated) > 0 {
		fmt.Println("\n⚠ Outdated Tools:")
		fmt.Println(separator("-"))
		for _, name := range result.Outdated {
			if tool, exists := result.Tools[name]; exists {
				fmt.Printf("  ⚠ %s (installed: %s, minimum: %s)\n",
					name, tool.InstalledVersion, tool.MinVersion)
			}
		}
	}

	fmt.Println("\n" + separator("="))
	if result.AllAvailable {
		fmt.Println("✓ All required tools are available!")
		fmt.Println("You can proceed with project generation.")
	} else {
		fmt.Println("✗ Some tools are missing.")
		fmt.Println("Install missing tools or use --no-external-tools flag for fallback generation.")
	}
	fmt.Println(separator("=") + "\n")
}

// indentText indents each line of text with the given prefix
func indentText(text, indent string) string {
	lines := splitLines(text)
	result := ""
	for _, line := range lines {
		result += indent + line + "\n"
	}
	return result
}

// splitLines splits text into lines
func splitLines(text string) []string {
	var lines []string
	current := ""
	for _, char := range text {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

var initConfigCmd = &cobra.Command{
	Use:   "init-config [output-file]",
	Short: "Generate a configuration template file",
	Long: `Generate a configuration template file with all available options documented.

This command creates a YAML configuration file that you can customize for your
project. The template includes examples for all supported component types and
integration options.

Examples:
  # Generate config in current directory
  generator init-config

  # Generate config with custom name
  generator init-config my-project.yaml

  # Generate minimal config (only required fields)
  generator init-config --minimal

  # Generate full-stack example
  generator init-config --example fullstack`,
	RunE: runInitConfig,
}

var (
	minimalConfig bool
	exampleType   string
)

func init() {
	initConfigCmd.Flags().BoolVar(&minimalConfig, "minimal", false, "Generate minimal configuration with only required fields")
	initConfigCmd.Flags().StringVar(&exampleType, "example", "", "Generate example configuration (fullstack, frontend, backend, mobile)")
}

func runInitConfig(cmd *cobra.Command, args []string) error {
	// Determine output file
	outputFile := "project-config.yaml"
	if len(args) > 0 {
		outputFile = args[0]
	}

	// Check if file already exists
	if _, err := os.Stat(outputFile); err == nil {
		return fmt.Errorf("file already exists: %s (remove it first or choose a different name)", outputFile)
	}

	// Generate configuration based on flags
	var config *models.ProjectConfig
	if exampleType != "" {
		config = generateExampleConfig(exampleType)
	} else if minimalConfig {
		config = generateMinimalConfig()
	} else {
		config = generateFullConfig()
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Add header comments
	header := generateConfigHeader()
	fullContent := header + "\n" + string(data)

	// Write to file
	if err := os.WriteFile(outputFile, []byte(fullContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✓ Configuration template created: %s\n", outputFile)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Edit the configuration file to match your project requirements")
	fmt.Println("2. Enable/disable components as needed")
	fmt.Println("3. Run: generator generate --config " + outputFile)

	return nil
}

func generateConfigHeader() string {
	return `# Open Source Project Generator Configuration
# 
# This file defines the project structure and components to generate.
# Edit this file to customize your project.
#
# For more information, see: https://github.com/cuesoftinc/open-source-project-generator

`
}

func generateMinimalConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:        "my-project",
		Description: "A new project",
		OutputDir:   "./my-project",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name": "web-app",
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
			APIEndpoints:          map[string]string{},
			SharedEnvironment:     map[string]string{},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
			CreateBackup:     true,
			ForceOverwrite:   false,
		},
	}
}

func generateFullConfig() *models.ProjectConfig {
	return &models.ProjectConfig{
		Name:        "my-fullstack-project",
		Description: "A full-stack project with frontend, backend, and mobile apps",
		OutputDir:   "./my-fullstack-project",
		Components: []models.ComponentConfig{
			{
				Type:    "nextjs",
				Name:    "web-app",
				Enabled: true,
				Config: map[string]interface{}{
					"name":       "web-app",
					"typescript": true,
					"tailwind":   true,
					"app_router": true,
					"eslint":     true,
				},
			},
			{
				Type:    "go-backend",
				Name:    "api-server",
				Enabled: true,
				Config: map[string]interface{}{
					"name":      "api-server",
					"module":    "github.com/user/my-fullstack-project",
					"framework": "gin",
					"port":      8080,
				},
			},
			{
				Type:    "android",
				Name:    "mobile-android",
				Enabled: false,
				Config: map[string]interface{}{
					"name":       "mobile-android",
					"package":    "com.example.myapp",
					"min_sdk":    24,
					"target_sdk": 34,
					"language":   "kotlin",
				},
			},
			{
				Type:    "ios",
				Name:    "mobile-ios",
				Enabled: false,
				Config: map[string]interface{}{
					"name":              "mobile-ios",
					"bundle_id":         "com.example.myapp",
					"deployment_target": "15.0",
					"language":          "swift",
				},
			},
		},
		Integration: models.IntegrationConfig{
			GenerateDockerCompose: true,
			GenerateScripts:       true,
			APIEndpoints: map[string]string{
				"backend": "http://localhost:8080",
			},
			SharedEnvironment: map[string]string{
				"NODE_ENV":    "development",
				"API_URL":     "http://localhost:8080",
				"API_TIMEOUT": "30000",
			},
		},
		Options: models.ProjectOptions{
			UseExternalTools: true,
			DryRun:           false,
			Verbose:          false,
			CreateBackup:     true,
			ForceOverwrite:   false,
		},
	}
}

func generateExampleConfig(exampleType string) *models.ProjectConfig {
	switch exampleType {
	case "fullstack":
		return generateFullConfig()
	case "frontend":
		return &models.ProjectConfig{
			Name:        "my-frontend-project",
			Description: "A frontend-only project",
			OutputDir:   "./my-frontend-project",
			Components: []models.ComponentConfig{
				{
					Type:    "nextjs",
					Name:    "web-app",
					Enabled: true,
					Config: map[string]interface{}{
						"name":       "web-app",
						"typescript": true,
						"tailwind":   true,
						"app_router": true,
					},
				},
			},
			Integration: models.IntegrationConfig{
				GenerateDockerCompose: true,
				GenerateScripts:       true,
				APIEndpoints:          map[string]string{},
			},
			Options: models.ProjectOptions{
				UseExternalTools: true,
				CreateBackup:     true,
			},
		}
	case "backend":
		return &models.ProjectConfig{
			Name:        "my-backend-project",
			Description: "A backend-only project",
			OutputDir:   "./my-backend-project",
			Components: []models.ComponentConfig{
				{
					Type:    "go-backend",
					Name:    "api-server",
					Enabled: true,
					Config: map[string]interface{}{
						"name":      "api-server",
						"module":    "github.com/user/my-backend-project",
						"framework": "gin",
						"port":      8080,
					},
				},
			},
			Integration: models.IntegrationConfig{
				GenerateDockerCompose: true,
				GenerateScripts:       true,
			},
			Options: models.ProjectOptions{
				UseExternalTools: true,
				CreateBackup:     true,
			},
		}
	case "mobile":
		return &models.ProjectConfig{
			Name:        "my-mobile-project",
			Description: "A mobile-only project",
			OutputDir:   "./my-mobile-project",
			Components: []models.ComponentConfig{
				{
					Type:    "android",
					Name:    "mobile-android",
					Enabled: true,
					Config: map[string]interface{}{
						"name":       "mobile-android",
						"package":    "com.example.myapp",
						"min_sdk":    24,
						"target_sdk": 34,
					},
				},
				{
					Type:    "ios",
					Name:    "mobile-ios",
					Enabled: true,
					Config: map[string]interface{}{
						"name":              "mobile-ios",
						"bundle_id":         "com.example.myapp",
						"deployment_target": "15.0",
					},
				},
			},
			Integration: models.IntegrationConfig{
				GenerateDockerCompose: false,
				GenerateScripts:       true,
			},
			Options: models.ProjectOptions{
				UseExternalTools: true,
				CreateBackup:     true,
			},
		}
	default:
		// Default to fullstack
		return generateFullConfig()
	}
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new project from configuration",
	Long: `Generate a new project using the specified configuration file.

The generate command orchestrates bootstrap tools (like create-next-app, go mod init)
to create project components, then maps them to a standardized directory structure
and integrates them together.

Examples:
  # Interactive mode (prompts for configuration)
  generator generate --interactive

  # Generate from config file
  generator generate --config project.yaml

  # Generate with custom output directory
  generator generate --config project.yaml --output ./my-project

  # Dry run (preview without creating files)
  generator generate --config project.yaml --dry-run

  # Force fallback generation (don't use external tools)
  generator generate --config project.yaml --no-external-tools

  # Verbose output for debugging
  generator generate --config project.yaml --verbose`,
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to project configuration file (YAML/JSON)")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for generated project (overrides config)")
	generateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview what would be generated without creating files")
	generateCmd.Flags().BoolVar(&noExternalTools, "no-external-tools", false, "Force fallback generation (don't use external bootstrap tools)")
	generateCmd.Flags().BoolVar(&offlineMode, "offline", false, "Force offline mode (use cached tools and fallback generators)")
	generateCmd.Flags().BoolVar(&createBackup, "backup", true, "Create backup before overwriting existing directory")
	generateCmd.Flags().BoolVar(&forceOverwrite, "force", false, "Force overwrite existing directory")
	generateCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode (prompts for configuration)")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\nReceived interrupt signal, cancelling generation...")
		cancel()
	}()

	// Initialize logger
	log := logger.NewLogger()
	if verbose {
		log.SetLevel(logger.DebugLevel)
	}

	// Load or create configuration
	var config *models.ProjectConfig
	var err error

	if interactive {
		// Interactive mode
		config, err = runInteractiveMode(log)
		if err != nil {
			return fmt.Errorf("interactive mode failed: %w", err)
		}
	} else {
		// Non-interactive mode - require config file
		if configFile == "" {
			return fmt.Errorf("config file is required in non-interactive mode (use --config or --interactive)")
		}

		config, err = loadConfigFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
	}

	// Apply command-line overrides
	if outputDir != "" {
		config.OutputDir = outputDir
	}
	if dryRun {
		config.Options.DryRun = true
	}
	if noExternalTools {
		config.Options.UseExternalTools = false
	} else {
		config.Options.UseExternalTools = true
	}
	if createBackup {
		config.Options.CreateBackup = true
	}
	if forceOverwrite {
		config.Options.ForceOverwrite = true
	}
	if verbose {
		config.Options.Verbose = true
	}

	// Create project coordinator
	coordinator := orchestrator.NewProjectCoordinator(log)

	// Set offline mode if requested
	if offlineMode {
		log.Info("Forcing offline mode - will use cached tools and fallback generators")
		coordinator.SetOfflineMode(true)
	}

	// Check and display offline status
	if coordinator.IsOffline() {
		fmt.Println(coordinator.GetOfflineMessage())
	}

	// Execute generation
	startTime := time.Now()

	if dryRun {
		log.Info("Running in dry-run mode (no files will be created)")
		result, err := coordinator.DryRun(ctx, config)
		if err != nil {
			return fmt.Errorf("dry-run failed: %w", err)
		}

		// Display preview
		displayPreview(result.(*models.PreviewResult), log)
	} else {
		log.Info(fmt.Sprintf("Starting project generation: %s", config.Name))
		result, err := coordinator.Generate(ctx, config)
		if err != nil {
			return fmt.Errorf("generation failed: %w", err)
		}

		// Display results
		displayResults(result.(*models.GenerationResult), log)
	}

	duration := time.Since(startTime)
	log.Info(fmt.Sprintf("Completed in %v", duration))

	return nil
}

// loadConfigFile loads a project configuration from a YAML or JSON file
func loadConfigFile(path string) (*models.ProjectConfig, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML/JSON
	var config models.ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set default output directory if not specified
	if config.OutputDir == "" {
		config.OutputDir = filepath.Join(".", config.Name)
	}

	return &config, nil
}

// runInteractiveMode runs the interactive configuration wizard
func runInteractiveMode(log *logger.Logger) (*models.ProjectConfig, error) {
	// TODO: Implement interactive mode in a future task
	// For now, return an error indicating it's not yet implemented
	return nil, fmt.Errorf("interactive mode is not yet implemented - please use --config flag")
}

// displayPreview displays the dry-run preview results
func displayPreview(preview *models.PreviewResult, log *logger.Logger) {
	fmt.Println("\n" + separator("="))
	fmt.Println("DRY RUN PREVIEW")
	fmt.Println(separator("="))

	fmt.Printf("\nProject Root: %s\n", preview.ProjectRoot)

	fmt.Println("\nComponents:")
	fmt.Println(separator("-"))
	for i, comp := range preview.Components {
		fmt.Printf("\n%d. %s (%s)\n", i+1, comp.Name, comp.Type)
		fmt.Printf("   Target Path: %s\n", comp.TargetPath)
		fmt.Printf("   Method: %s\n", comp.Method)
		fmt.Printf("   Tool: %s\n", comp.ToolUsed)

		if len(comp.Files) > 0 {
			fmt.Printf("   Expected Files: %d files\n", len(comp.Files))
		}

		if len(comp.Warnings) > 0 {
			fmt.Println("   Warnings:")
			for _, warning := range comp.Warnings {
				fmt.Printf("     ⚠ %s\n", warning)
			}
		}
	}

	if len(preview.Structure) > 0 {
		fmt.Println("\nDirectory Structure:")
		fmt.Println(separator("-"))
		for _, dir := range preview.Structure {
			fmt.Printf("  📁 %s\n", dir)
		}
	}

	if len(preview.Files) > 0 {
		fmt.Println("\nIntegration Files:")
		fmt.Println(separator("-"))
		for _, file := range preview.Files {
			fmt.Printf("  📄 %s\n", file)
		}
	}

	if len(preview.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		fmt.Println(separator("-"))
		for _, warning := range preview.Warnings {
			fmt.Printf("  ⚠ %s\n", warning)
		}
	}

	fmt.Println("\n" + separator("="))
	fmt.Println("This was a dry run. No files were created.")
	fmt.Println("Run without --dry-run to generate the project.")
	fmt.Println(separator("=") + "\n")
}

// displayResults displays the generation results
func displayResults(result *models.GenerationResult, log *logger.Logger) {
	fmt.Println("\n" + separator("="))
	if result.Success {
		fmt.Println("✓ PROJECT GENERATION SUCCESSFUL")
	} else {
		fmt.Println("✗ PROJECT GENERATION FAILED")
	}
	fmt.Println(separator("="))

	fmt.Printf("\nProject Root: %s\n", result.ProjectRoot)
	fmt.Printf("Duration: %v\n", result.Duration)

	fmt.Println("\nComponents:")
	fmt.Println(separator("-"))
	for i, comp := range result.Components {
		status := "✓"
		if !comp.Success {
			status = "✗"
		}

		fmt.Printf("\n%s %d. %s (%s)\n", status, i+1, comp.Name, comp.Type)
		fmt.Printf("   Method: %s\n", comp.Method)
		fmt.Printf("   Tool: %s\n", comp.ToolUsed)
		fmt.Printf("   Output: %s\n", comp.OutputPath)
		fmt.Printf("   Duration: %v\n", comp.Duration)

		if comp.Error != nil {
			fmt.Printf("   Error: %v\n", comp.Error)
		}

		if len(comp.ManualSteps) > 0 {
			fmt.Println("   Manual Steps Required:")
			for j, step := range comp.ManualSteps {
				fmt.Printf("     %d. %s\n", j+1, step)
			}
		}

		if len(comp.Warnings) > 0 {
			fmt.Println("   Warnings:")
			for _, warning := range comp.Warnings {
				fmt.Printf("     ⚠ %s\n", warning)
			}
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		fmt.Println(separator("-"))
		for _, warning := range result.Warnings {
			fmt.Printf("  ⚠ %s\n", warning)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		fmt.Println(separator("-"))
		for i, err := range result.Errors {
			fmt.Printf("  %d. %v\n", i+1, err)
		}
	}

	if result.Success {
		fmt.Println("\n" + separator("="))
		fmt.Println("Next Steps:")
		fmt.Println(separator("-"))
		fmt.Printf("1. Navigate to project: cd %s\n", result.ProjectRoot)
		fmt.Println("2. Read the README.md for setup instructions")
		fmt.Println("3. Review any manual steps listed above")
		fmt.Println("4. Start development!")
		fmt.Println(separator("=") + "\n")
	}
}

// separator creates a separator line
func separator(char string) string {
	return repeatString(char, 70)
}

// repeatString repeats a string n times
func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

var cacheToolsCmd = &cobra.Command{
	Use:   "cache-tools",
	Short: "Manage tool cache for offline use",
	Long: `Manage the tool availability cache for offline project generation.

This command allows you to:
- View cache statistics and status
- Clear the cache to force fresh tool checks
- Save current tool availability for offline use
- Display information about which tools can be cached

The cache stores tool availability and version information to speed up
subsequent tool checks and enable offline operation.

Examples:
  # View cache statistics
  generator cache-tools --stats

  # Clear the cache
  generator cache-tools --clear

  # Save current tool availability
  generator cache-tools --save

  # Show cache file location
  generator cache-tools --info`,
	RunE: runCacheTools,
}

var (
	showCacheStats bool
	clearCache     bool
	saveCache      bool
	showCacheInfo  bool
)

func init() {
	cacheToolsCmd.Flags().BoolVar(&showCacheStats, "stats", false, "Show cache statistics")
	cacheToolsCmd.Flags().BoolVar(&clearCache, "clear", false, "Clear the tool cache")
	cacheToolsCmd.Flags().BoolVar(&saveCache, "save", false, "Save current tool availability to cache")
	cacheToolsCmd.Flags().BoolVar(&showCacheInfo, "info", false, "Show cache information and location")
	cacheToolsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func runCacheTools(cmd *cobra.Command, args []string) error {
	// Initialize logger
	log := logger.NewLogger()
	if verbose {
		log.SetLevel(logger.DebugLevel)
	}

	// Create tool discovery with cache
	toolDiscovery := orchestrator.NewToolDiscovery(log)

	// If no flags specified, show stats by default
	if !showCacheStats && !clearCache && !saveCache && !showCacheInfo {
		showCacheStats = true
	}

	// Handle clear cache
	if clearCache {
		log.Info("Clearing tool cache...")
		toolDiscovery.ClearCache()
		fmt.Println("✓ Tool cache cleared successfully")
		return nil
	}

	// Handle save cache
	if saveCache {
		log.Info("Checking and caching all registered tools...")

		// Get all registered tools
		allTools := toolDiscovery.ListRegisteredTools()

		// Check each tool to populate cache
		for _, toolName := range allTools {
			log.Debug(fmt.Sprintf("Checking tool: %s", toolName))
			available, _ := toolDiscovery.IsAvailable(toolName)
			if available {
				// Also get version to cache it
				toolDiscovery.GetVersion(toolName)
			}
		}

		// Save cache to disk
		if err := toolDiscovery.SaveCache(); err != nil {
			return fmt.Errorf("failed to save cache: %w", err)
		}

		fmt.Printf("✓ Tool availability cached for %d tools\n", len(allTools))
		fmt.Println("\nCached tools can be used for offline project generation.")
		fmt.Println("Note: External bootstrap tools themselves are not cached,")
		fmt.Println("only their availability status. For true offline operation,")
		fmt.Println("ensure tools are installed before going offline.")
		return nil
	}

	// Handle show info
	if showCacheInfo {
		stats := toolDiscovery.GetCacheStats()

		fmt.Println("\n" + separator("="))
		fmt.Println("TOOL CACHE INFORMATION")
		fmt.Println(separator("="))

		if enabled, ok := stats["enabled"].(bool); ok && enabled {
			if cacheFile, ok := stats["cache_file"].(string); ok {
				fmt.Printf("\nCache File: %s\n", cacheFile)
			}
			if ttl, ok := stats["ttl"].(string); ok {
				fmt.Printf("Cache TTL: %s\n", ttl)
			}
			if lastSaved, ok := stats["last_saved"].(time.Time); ok {
				fmt.Printf("Last Saved: %s\n", lastSaved.Format(time.RFC3339))
			}

			fmt.Println("\nCacheable Information:")
			fmt.Println("  - Tool availability (whether tool is in PATH)")
			fmt.Println("  - Tool versions")
			fmt.Println("  - Last check timestamp")

			fmt.Println("\nNot Cached:")
			fmt.Println("  - Tool executables themselves")
			fmt.Println("  - Tool dependencies")
			fmt.Println("  - Network-based resources")

			fmt.Println("\nOffline Operation:")
			fmt.Println("  - Cache stores tool availability for quick checks")
			fmt.Println("  - Fallback generators work fully offline")
			fmt.Println("  - Bootstrap tools require prior installation")
		} else {
			fmt.Println("\n✗ Cache is not enabled")
		}

		fmt.Println("\n" + separator("=") + "\n")
		return nil
	}

	// Handle show stats (default)
	if showCacheStats {
		stats := toolDiscovery.GetCacheStats()

		fmt.Println("\n" + separator("="))
		fmt.Println("TOOL CACHE STATISTICS")
		fmt.Println(separator("="))

		if enabled, ok := stats["enabled"].(bool); ok && !enabled {
			fmt.Println("\n✗ Cache is not enabled")
			fmt.Println("\n" + separator("=") + "\n")
			return nil
		}

		if total, ok := stats["total"].(int); ok {
			fmt.Printf("\nTotal Entries: %d\n", total)
		}
		if available, ok := stats["available"].(int); ok {
			fmt.Printf("Available Tools: %d\n", available)
		}
		if unavailable, ok := stats["unavailable"].(int); ok {
			fmt.Printf("Unavailable Tools: %d\n", unavailable)
		}
		if expired, ok := stats["expired"].(int); ok {
			fmt.Printf("Expired Entries: %d\n", expired)
		}

		if cacheFile, ok := stats["cache_file"].(string); ok {
			fmt.Printf("\nCache File: %s\n", cacheFile)
		}
		if ttl, ok := stats["ttl"].(string); ok {
			fmt.Printf("Cache TTL: %s\n", ttl)
		}
		if lastSaved, ok := stats["last_saved"].(time.Time); ok {
			fmt.Printf("Last Saved: %s\n", lastSaved.Format(time.RFC3339))
		}
		if timeSinceCheck, ok := stats["time_since_check"].(string); ok {
			fmt.Printf("Time Since Check: %s\n", timeSinceCheck)
		}

		fmt.Println("\nCommands:")
		fmt.Println("  --clear    Clear the cache")
		fmt.Println("  --save     Cache all registered tools")
		fmt.Println("  --info     Show cache information")

		fmt.Println("\n" + separator("=") + "\n")
	}

	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Open Source Project Generator\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Built: %s\n", buildTime)
		fmt.Printf("Commit: %s\n", gitCommit)
	},
}
