package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/internal/config"
	"github.com/cuesoftinc/open-source-project-generator/internal/generator/bootstrap"
	"github.com/cuesoftinc/open-source-project-generator/internal/generator/fallback"
	"github.com/cuesoftinc/open-source-project-generator/internal/generator/mapper"
	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/security"
)

// ProjectCoordinator orchestrates the complete project generation workflow
type ProjectCoordinator struct {
	validator          *config.Validator
	toolDiscovery      *ToolDiscovery
	offlineDetector    *OfflineDetector
	structureMapper    interfaces.StructureMapperInterface
	componentMapper    interfaces.ComponentMapperInterface
	integrationManager *IntegrationManager //nolint:unused // Reserved for future use
	fallbackRegistry   *fallback.Registry
	logger             *logger.Logger
	backupDir          string
	sanitizer          *security.Sanitizer
	rollbackManager    *RollbackManager
}

// NewProjectCoordinator creates a new project coordinator
func NewProjectCoordinator(log *logger.Logger) *ProjectCoordinator {
	componentMapper := mapper.NewComponentMapper()
	structureMapper := mapper.NewStructureMapper(componentMapper)
	offlineDetector := NewOfflineDetector(DefaultOfflineDetectorConfig(), log)

	return &ProjectCoordinator{
		validator:        config.NewValidator(),
		toolDiscovery:    NewToolDiscovery(log),
		offlineDetector:  offlineDetector,
		structureMapper:  structureMapper,
		componentMapper:  componentMapper,
		fallbackRegistry: fallback.DefaultRegistry(),
		logger:           log,
		backupDir:        ".backups",
		sanitizer:        security.NewSanitizer(),
		rollbackManager:  NewRollbackManager(log),
	}
}

// SetOfflineMode sets the offline mode for the coordinator
func (pc *ProjectCoordinator) SetOfflineMode(offline bool) {
	if pc.offlineDetector != nil {
		pc.offlineDetector.ForceOffline(offline)
	}
	if pc.toolDiscovery != nil {
		pc.toolDiscovery.SetOfflineMode(offline)
	}
}

// IsOffline returns whether the coordinator is in offline mode
func (pc *ProjectCoordinator) IsOffline() bool {
	if pc.offlineDetector != nil {
		return pc.offlineDetector.IsOffline()
	}
	return false
}

// GetOfflineMessage returns a user-friendly offline status message
func (pc *ProjectCoordinator) GetOfflineMessage() string {
	if pc.offlineDetector != nil {
		return pc.offlineDetector.GetOfflineMessage()
	}
	return ""
}

// Ensure ProjectCoordinator implements the interface
var _ interfaces.ProjectCoordinatorInterface = (*ProjectCoordinator)(nil)

// enableFileLogging enables logging to a file
func (pc *ProjectCoordinator) enableFileLogging(logFilePath string) error {
	// Create log directory
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0750); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Enable file logging
	if err := pc.logger.EnableFileLogging(logFilePath); err != nil {
		return fmt.Errorf("failed to enable file logging: %w", err)
	}

	return nil
}

// Generate orchestrates the complete project generation workflow
func (pc *ProjectCoordinator) Generate(ctx context.Context, configInterface interface{}) (interface{}, error) {
	// Type assert the configuration
	config, ok := configInterface.(*models.ProjectConfig)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type: expected *models.ProjectConfig")
	}

	startTime := time.Now()

	// Initialize result
	result := &models.GenerationResult{
		Success:     false,
		ProjectRoot: config.OutputDir,
		Components:  make([]*models.ComponentResult, 0),
		Errors:      make([]error, 0),
		Warnings:    make([]string, 0),
		DryRun:      config.Options.DryRun,
	}

	// Enable file logging if not in dry-run mode
	var logFilePath string
	if !config.Options.DryRun {
		logFilePath = filepath.Join(config.OutputDir, ".logs", fmt.Sprintf("generation-%s.log", time.Now().Format("20060102-150405")))
		if err := pc.enableFileLogging(logFilePath); err != nil {
			pc.logger.Warn(fmt.Sprintf("Failed to enable file logging: %v", err))
		} else {
			result.LogFile = logFilePath
			pc.logger.Info(fmt.Sprintf("Logging to file: %s", logFilePath))
		}
	}

	pc.logger.Info(fmt.Sprintf("Starting project generation: %s", config.Name))
	pc.logger.Info(fmt.Sprintf("Output directory: %s", config.OutputDir))
	pc.logger.Info(fmt.Sprintf("Dry-run mode: %v", config.Options.DryRun))
	pc.logger.Info(fmt.Sprintf("Components to generate: %d", len(config.Components)))

	// Step 1: Validate configuration
	pc.logger.Info("Step 1/11: Validating configuration...")
	pc.logger.Debug(fmt.Sprintf("Project name: %s", config.Name))
	pc.logger.Debug(fmt.Sprintf("Description: %s", config.Description))
	if err := pc.Validate(config); err != nil {
		pc.logger.Error(fmt.Sprintf("Configuration validation failed: %v", err))
		result.Errors = append(result.Errors, fmt.Errorf("configuration validation failed: %w", err))
		result.Duration = time.Since(startTime)
		return result, err
	}
	pc.logger.Info("Configuration validation passed")

	// Step 2: Apply defaults
	pc.logger.Info("Step 2/11: Applying default configuration values...")
	if err := pc.validator.ApplyDefaults(config); err != nil {
		pc.logger.Error(fmt.Sprintf("Failed to apply defaults: %v", err))
		result.Errors = append(result.Errors, fmt.Errorf("failed to apply defaults: %w", err))
		result.Duration = time.Since(startTime)
		return result, err
	}
	pc.logger.Info("Defaults applied successfully")

	// Step 3: Sanitize inputs
	pc.logger.Info("Step 3/11: Sanitizing user inputs...")
	pc.logger.Debug(fmt.Sprintf("Sanitizing project name: %s", config.Name))
	pc.logger.Debug(fmt.Sprintf("Sanitizing output directory: %s", config.OutputDir))
	if err := pc.sanitizeConfiguration(config); err != nil {
		pc.logger.Error(fmt.Sprintf("Input sanitization failed: %v", err))
		result.Errors = append(result.Errors, fmt.Errorf("input sanitization failed: %w", err))
		result.Duration = time.Since(startTime)
		return result, err
	}
	pc.logger.Info("Input sanitization completed")

	// Step 4: Discover available tools
	pc.logger.Info("Step 4/11: Discovering available bootstrap tools...")
	toolCheckResult, err := pc.discoverTools(config)
	if err != nil {
		pc.logger.Error(fmt.Sprintf("Tool discovery failed: %v", err))
		result.Errors = append(result.Errors, fmt.Errorf("tool discovery failed: %w", err))
		result.Duration = time.Since(startTime)
		return result, err
	}

	// Log tool availability
	pc.logToolAvailability(toolCheckResult)
	pc.logger.Info("Tool discovery completed")

	// Step 5: Create backup if needed
	if config.Options.CreateBackup && !config.Options.DryRun {
		pc.logger.Info("Step 5/11: Creating backup of existing directory...")
		pc.logger.Debug(fmt.Sprintf("Backup directory: %s", pc.backupDir))
		if err := pc.createBackup(config.OutputDir); err != nil {
			pc.logger.Warn(fmt.Sprintf("Failed to create backup: %v", err))
			result.Warnings = append(result.Warnings, fmt.Sprintf("Failed to create backup: %v", err))
		} else {
			pc.logger.Info("Backup created successfully")
		}
	} else {
		pc.logger.Info("Step 5/11: Skipping backup (not requested or dry-run mode)")
	}

	// Step 6: Prepare output directory
	if !config.Options.DryRun {
		pc.logger.Info("Step 6/11: Preparing output directory...")
		pc.logger.Debug(fmt.Sprintf("Target directory: %s", config.OutputDir))
		if err := pc.prepareOutputDirectory(config); err != nil {
			pc.logger.Error(fmt.Sprintf("Failed to prepare output directory: %v", err))
			result.Errors = append(result.Errors, fmt.Errorf("failed to prepare output directory: %w", err))
			result.Duration = time.Since(startTime)
			return result, err
		}
		pc.logger.Info("Output directory prepared successfully")
	} else {
		pc.logger.Info("Step 6/11: Skipping directory preparation (dry-run mode)")
	}

	// Step 7: Generate each component
	pc.logger.Info(fmt.Sprintf("Step 7/11: Generating %d components...", len(config.Components)))
	for i, comp := range config.Components {
		if comp.Enabled {
			pc.logger.Debug(fmt.Sprintf("  Component %d: %s (%s)", i+1, comp.Name, comp.Type))
		}
	}
	componentResults, err := pc.generateComponents(ctx, config, toolCheckResult)
	if err != nil {
		result.Errors = append(result.Errors, err)
		result.Components = componentResults
		result.Duration = time.Since(startTime)

		// Print formatted error message
		pc.logger.PrintHeader("PROJECT GENERATION FAILED")
		pc.logger.Error(fmt.Sprintf("Generation failed: %v", err))

		// Attempt rollback on error
		if !config.Options.DryRun {
			pc.logger.Warn("Attempting to rollback changes...")
			if rollbackErr := pc.rollbackManager.Rollback(ctx); rollbackErr != nil {
				pc.logger.Error(fmt.Sprintf("Rollback failed: %v", rollbackErr))
				result.Warnings = append(result.Warnings, fmt.Sprintf("Rollback failed: %v", rollbackErr))
			} else {
				pc.logger.Success("Rollback completed successfully")
			}
		}

		// Show what was attempted
		if len(componentResults) > 0 {
			pc.logger.PrintSection("Component Status")
			for _, comp := range componentResults {
				if comp.Success {
					pc.logger.Success(fmt.Sprintf("%s (%s) - Generated", comp.Name, comp.Type))
				} else {
					pc.logger.Error(fmt.Sprintf("%s (%s) - Failed", comp.Name, comp.Type))
				}
			}
		}

		return result, err
	}

	result.Components = componentResults

	// Step 8: Map components to target structure
	if !config.Options.DryRun {
		pc.logger.Info("Step 8/11: Mapping components to target structure...")
		pc.logger.Debug("Target structure: App/, CommonServer/, Mobile/, Deploy/")
		if err := pc.mapComponentsToStructure(ctx, config, componentResults); err != nil {
			pc.logger.Error(fmt.Sprintf("Structure mapping failed: %v", err))
			result.Errors = append(result.Errors, fmt.Errorf("structure mapping failed: %w", err))
			result.Duration = time.Since(startTime)

			// Attempt rollback
			pc.logger.Warn("Mapping failed, attempting rollback...")
			if rollbackErr := pc.rollbackManager.Rollback(ctx); rollbackErr != nil {
				pc.logger.Error(fmt.Sprintf("Rollback failed: %v", rollbackErr))
				result.Warnings = append(result.Warnings, fmt.Sprintf("Rollback failed: %v", rollbackErr))
			} else {
				pc.logger.Info("Rollback completed successfully")
			}

			return result, err
		}
		pc.logger.Info("Structure mapping completed successfully")
	} else {
		pc.logger.Info("Step 8/11: Skipping structure mapping (dry-run mode)")
	}

	// Step 9: Integrate components
	if !config.Options.DryRun {
		pc.logger.Info("Step 9/11: Integrating components...")
		pc.logger.Debug("Generating Docker Compose, environment files, and scripts")
		if err := pc.integrateComponents(ctx, config, componentResults); err != nil {
			pc.logger.Error(fmt.Sprintf("Component integration failed: %v", err))
			result.Errors = append(result.Errors, fmt.Errorf("component integration failed: %w", err))
			result.Warnings = append(result.Warnings, "Components generated but integration incomplete")
			// Don't rollback here - components are already generated
		} else {
			pc.logger.Info("Component integration completed successfully")
		}
	} else {
		pc.logger.Info("Step 9/11: Skipping component integration (dry-run mode)")
	}

	// Step 10: Validate final structure
	if !config.Options.DryRun {
		pc.logger.Info("Step 10/11: Validating final project structure...")
		if err := pc.validateFinalStructure(config, componentResults); err != nil {
			pc.logger.Warn(fmt.Sprintf("Structure validation warnings: %v", err))
			result.Warnings = append(result.Warnings, fmt.Sprintf("Structure validation warnings: %v", err))
		} else {
			pc.logger.Info("Structure validation passed")
		}
	} else {
		pc.logger.Info("Step 10/11: Skipping structure validation (dry-run mode)")
	}

	// Step 11: Run security scan
	if !config.Options.DryRun {
		pc.logger.Info("Step 11/11: Running security scan...")
		scanResult, err := pc.runSecurityScan(ctx, config.OutputDir)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Security scan failed: %v", err))
		} else {
			// Add security scan results to generation result
			result.SecurityScanResult = scanResult

			// Log security issues
			if len(scanResult.Issues) > 0 {
				pc.logger.Warn(fmt.Sprintf("Security scan found %d issues", len(scanResult.Issues)))
				for _, issue := range scanResult.Issues {
					if issue.Severity == "critical" || issue.Severity == "high" {
						pc.logger.Warn(fmt.Sprintf("  [%s] %s: %s (line %d)", issue.Severity, issue.File, issue.Description, issue.Line))
					}
				}
			}

			// Log security warnings
			if len(scanResult.Warnings) > 0 {
				pc.logger.Info(fmt.Sprintf("Security scan found %d warnings", len(scanResult.Warnings)))
			}

			if !scanResult.Passed {
				result.Warnings = append(result.Warnings, "Security scan found critical or high severity issues - please review and fix")
			}
		}
	}

	// Mark as successful
	result.Success = true
	result.Duration = time.Since(startTime)

	// Print formatted success message
	pc.logger.PrintHeader("PROJECT GENERATION COMPLETED")
	pc.logger.Success(fmt.Sprintf("Project generated successfully in %v", result.Duration))
	pc.logger.PrintKeyValue("Project Name", config.Name)
	pc.logger.PrintKeyValue("Location", config.OutputDir)
	pc.logger.PrintKeyValue("Components", fmt.Sprintf("%d", len(componentResults)))

	// List generated components
	if len(componentResults) > 0 {
		pc.logger.PrintSection("Generated Components")
		for _, comp := range componentResults {
			if comp.Success {
				pc.logger.Success(fmt.Sprintf("%s (%s) - %s", comp.Name, comp.Type, comp.Method))
			}
		}
	}

	// Show warnings if any
	if len(result.Warnings) > 0 {
		pc.logger.PrintSection("Warnings")
		for _, warning := range result.Warnings {
			pc.logger.Warn(warning)
		}
	}

	// Show manual steps if any
	hasManualSteps := false
	for _, comp := range componentResults {
		if len(comp.ManualSteps) > 0 {
			hasManualSteps = true
			break
		}
	}

	if hasManualSteps {
		pc.logger.PrintSection("Manual Setup Required")
		for _, comp := range componentResults {
			if len(comp.ManualSteps) > 0 {
				pc.logger.Info(fmt.Sprintf("%s (%s):", comp.Name, comp.Type))
				for i, step := range comp.ManualSteps {
					pc.logger.PrintBullet(fmt.Sprintf("%d. %s", i+1, step))
				}
			}
		}
	}

	// Show next steps
	pc.logger.PrintSection("Next Steps")
	pc.logger.PrintBullet(fmt.Sprintf("Navigate to: cd %s", config.OutputDir))
	pc.logger.PrintBullet("Review the README.md for detailed instructions")
	if config.Integration.GenerateDockerCompose {
		pc.logger.PrintBullet("Run with Docker: docker-compose up")
	}
	if config.Integration.GenerateScripts {
		pc.logger.PrintBullet("Run development mode: ./dev.sh")
	}

	if result.LogFile != "" {
		pc.logger.PrintKeyValue("Log File", result.LogFile)
	}

	return result, nil
}

// DryRun previews what would be generated without creating files
func (pc *ProjectCoordinator) DryRun(ctx context.Context, configInterface interface{}) (interface{}, error) {
	// Type assert the configuration
	config, ok := configInterface.(*models.ProjectConfig)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type: expected *models.ProjectConfig")
	}

	pc.logger.Info(fmt.Sprintf("Starting dry-run for project: %s", config.Name))

	// Validate configuration
	if err := pc.Validate(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Apply defaults
	if err := pc.validator.ApplyDefaults(config); err != nil {
		return nil, fmt.Errorf("failed to apply defaults: %w", err)
	}

	// Discover tools
	toolCheckResult, err := pc.discoverTools(config)
	if err != nil {
		return nil, fmt.Errorf("tool discovery failed: %w", err)
	}

	// Build preview result
	preview := &models.PreviewResult{
		ProjectRoot: config.OutputDir,
		Components:  make([]*models.ComponentPreview, 0),
		Structure:   make([]string, 0),
		Files:       make([]string, 0),
		Warnings:    make([]string, 0),
	}

	// Preview each component
	for _, comp := range config.Components {
		if !comp.Enabled {
			continue
		}

		compPreview := &models.ComponentPreview{
			Type:       comp.Type,
			Name:       comp.Name,
			TargetPath: filepath.Join(config.OutputDir, pc.structureMapper.GetTargetPath(comp.Type)),
			Files:      make([]string, 0),
			Warnings:   make([]string, 0),
		}

		// Determine generation method
		toolsForComponent := pc.toolDiscovery.GetToolsForComponent(comp.Type)
		hasTools := true
		for _, tool := range toolsForComponent {
			if toolInfo, exists := toolCheckResult.Tools[tool]; !exists || !toolInfo.Available {
				hasTools = false
				break
			}
		}

		if hasTools && config.Options.UseExternalTools {
			compPreview.Method = "bootstrap"
			compPreview.ToolUsed = toolsForComponent[0]
		} else if pc.fallbackRegistry.Supports(comp.Type) {
			compPreview.Method = "fallback"
			compPreview.ToolUsed = "fallback-generator"
			compPreview.Warnings = append(compPreview.Warnings, "Using fallback generation - manual setup may be required")
		} else {
			compPreview.Method = "none"
			compPreview.Warnings = append(compPreview.Warnings, "No generation method available for this component type")
		}

		// Add expected files (simplified)
		compPreview.Files = pc.getExpectedFiles(comp.Type)

		preview.Components = append(preview.Components, compPreview)

		// Add to structure list
		preview.Structure = append(preview.Structure, compPreview.TargetPath)
	}

	// Add integration files
	if config.Integration.GenerateDockerCompose {
		preview.Files = append(preview.Files, filepath.Join(config.OutputDir, "docker-compose.yml"))
	}
	if config.Integration.GenerateScripts {
		preview.Files = append(preview.Files,
			filepath.Join(config.OutputDir, "build.sh"),
			filepath.Join(config.OutputDir, "dev.sh"),
			filepath.Join(config.OutputDir, "prod.sh"),
			filepath.Join(config.OutputDir, "docker.sh"),
		)
	}
	preview.Files = append(preview.Files,
		filepath.Join(config.OutputDir, ".env"),
		filepath.Join(config.OutputDir, "README.md"),
		filepath.Join(config.OutputDir, "TROUBLESHOOTING.md"),
	)

	pc.logger.Info("Dry-run completed successfully")

	return preview, nil
}

// Validate validates the project configuration
func (pc *ProjectCoordinator) Validate(configInterface interface{}) error {
	// Type assert the configuration
	config, ok := configInterface.(*models.ProjectConfig)
	if !ok {
		return fmt.Errorf("invalid configuration type: expected *models.ProjectConfig")
	}

	return pc.validator.Validate(config)
}

// sanitizeConfiguration sanitizes all user inputs in the configuration
func (pc *ProjectCoordinator) sanitizeConfiguration(config *models.ProjectConfig) error {
	// Sanitize project name
	sanitizedName, err := pc.sanitizer.SanitizeProjectName(config.Name)
	if err != nil {
		return fmt.Errorf("invalid project name: %w", err)
	}
	config.Name = sanitizedName

	// Sanitize output directory
	sanitizedPath, err := pc.sanitizer.SanitizePath(config.OutputDir)
	if err != nil {
		return fmt.Errorf("invalid output directory: %w", err)
	}
	config.OutputDir = sanitizedPath

	// Sanitize component names
	for i := range config.Components {
		sanitizedCompName, err := pc.sanitizer.SanitizeProjectName(config.Components[i].Name)
		if err != nil {
			return fmt.Errorf("invalid component name '%s': %w", config.Components[i].Name, err)
		}
		config.Components[i].Name = sanitizedCompName
	}

	return nil
}

// discoverTools discovers available bootstrap tools
func (pc *ProjectCoordinator) discoverTools(config *models.ProjectConfig) (*models.ToolCheckResult, error) {
	// Check offline mode first
	if pc.offlineDetector != nil {
		isOffline := pc.offlineDetector.IsOffline()
		if isOffline {
			pc.logger.Info(pc.offlineDetector.GetOfflineMessage())
			pc.toolDiscovery.SetOfflineMode(true)
		}
	}

	// Collect all required tools
	requiredTools := make(map[string]bool)

	for _, comp := range config.Components {
		if !comp.Enabled {
			continue
		}

		tools := pc.toolDiscovery.GetToolsForComponent(comp.Type)
		for _, tool := range tools {
			requiredTools[tool] = true
		}
	}

	// Convert to slice
	toolList := make([]string, 0, len(requiredTools))
	for tool := range requiredTools {
		toolList = append(toolList, tool)
	}

	// Check tool availability
	result, err := pc.toolDiscovery.CheckRequirements(toolList)
	if err != nil {
		return nil, err
	}

	toolCheckResult, ok := result.(*models.ToolCheckResult)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from tool discovery")
	}

	return toolCheckResult, nil
}

// logToolAvailability logs the availability of tools
func (pc *ProjectCoordinator) logToolAvailability(result *models.ToolCheckResult) {
	if result.AllAvailable {
		pc.logger.Info("All required tools are available")
	} else {
		pc.logger.Warn(fmt.Sprintf("%d of %d required tools are missing", len(result.Missing), len(result.Tools)))
		for _, missing := range result.Missing {
			pc.logger.Warn(fmt.Sprintf("  - %s is not available", missing))
		}
	}
}

// createBackup creates a backup of the output directory
func (pc *ProjectCoordinator) createBackup(outputDir string) error {
	// Use rollback manager to create and register backup
	backupPath, err := pc.rollbackManager.CreateBackup(outputDir)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	if backupPath != "" {
		pc.logger.Info(fmt.Sprintf("Backup created at: %s", backupPath))
	}

	return nil
}

// prepareOutputDirectory prepares the output directory for generation
func (pc *ProjectCoordinator) prepareOutputDirectory(config *models.ProjectConfig) error {
	// Check if directory exists
	if _, err := os.Stat(config.OutputDir); err == nil {
		// Directory exists
		if !config.Options.ForceOverwrite {
			return fmt.Errorf("output directory already exists: %s (use force_overwrite option to overwrite)", config.OutputDir)
		}

		// Remove existing directory
		pc.logger.Warn(fmt.Sprintf("Removing existing directory: %s", config.OutputDir))
		if err := os.RemoveAll(config.OutputDir); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	pc.logger.Info(fmt.Sprintf("Output directory prepared: %s", config.OutputDir))

	return nil
}

// generateComponents generates all enabled components
func (pc *ProjectCoordinator) generateComponents(ctx context.Context, config *models.ProjectConfig, toolCheckResult *models.ToolCheckResult) ([]*models.ComponentResult, error) {
	// Check if parallel generation is enabled (default: true)
	enableParallel := true
	if config.Options.DisableParallel {
		enableParallel = false
	}

	if enableParallel {
		return pc.generateComponentsParallel(ctx, config, toolCheckResult)
	}

	return pc.generateComponentsSequential(ctx, config, toolCheckResult)
}

// generateComponentsSequential generates components one at a time
func (pc *ProjectCoordinator) generateComponentsSequential(ctx context.Context, config *models.ProjectConfig, toolCheckResult *models.ToolCheckResult) ([]*models.ComponentResult, error) {
	results := make([]*models.ComponentResult, 0)

	for _, comp := range config.Components {
		if !comp.Enabled {
			continue
		}

		pc.logger.Info(fmt.Sprintf("Generating component: %s (%s)", comp.Name, comp.Type))

		result, err := pc.generateComponent(ctx, &comp, config, toolCheckResult)
		if err != nil {
			pc.logger.Error(fmt.Sprintf("Failed to generate component %s: %v", comp.Name, err))
			result.Error = err
			result.Success = false
		}

		results = append(results, result)

		// Stop on first error if not in continue-on-error mode
		if !result.Success {
			return results, fmt.Errorf("component generation failed: %s", comp.Name)
		}
	}

	return results, nil
}

// generateComponentsParallel generates independent components in parallel
func (pc *ProjectCoordinator) generateComponentsParallel(ctx context.Context, config *models.ProjectConfig, toolCheckResult *models.ToolCheckResult) ([]*models.ComponentResult, error) {
	// Collect enabled components
	enabledComponents := make([]models.ComponentConfig, 0)
	for _, comp := range config.Components {
		if comp.Enabled {
			enabledComponents = append(enabledComponents, comp)
		}
	}

	if len(enabledComponents) == 0 {
		return []*models.ComponentResult{}, nil
	}

	pc.logger.Info(fmt.Sprintf("Generating %d components in parallel...", len(enabledComponents)))

	// Create progress tracker
	progressTracker := NewProgressTracker(pc.logger.GetWriter(), len(enabledComponents), !config.Options.Verbose)

	// Create channels for results and errors
	type componentJob struct {
		index  int
		config models.ComponentConfig
	}

	type componentJobResult struct {
		index  int
		result *models.ComponentResult
		err    error
	}

	jobs := make(chan componentJob, len(enabledComponents))
	jobResults := make(chan componentJobResult, len(enabledComponents))

	// Determine number of workers (max 4 to avoid overwhelming the system)
	numWorkers := len(enabledComponents)
	if numWorkers > 4 {
		numWorkers = 4
	}

	// Start worker goroutines
	for w := 0; w < numWorkers; w++ {
		go func(workerID int) {
			for job := range jobs {
				if config.Options.Verbose {
					pc.logger.Info(fmt.Sprintf("[Worker %d] Generating component: %s (%s)", workerID, job.config.Name, job.config.Type))
				}

				result, err := pc.generateComponent(ctx, &job.config, config, toolCheckResult)
				if err != nil {
					pc.logger.Error(fmt.Sprintf("[Worker %d] Failed to generate component %s: %v", workerID, job.config.Name, err))
					if result != nil {
						result.Error = err
						result.Success = false
					}
				}

				jobResults <- componentJobResult{
					index:  job.index,
					result: result,
					err:    err,
				}
			}
		}(w)
	}

	// Send jobs to workers
	for i, comp := range enabledComponents {
		jobs <- componentJob{
			index:  i,
			config: comp,
		}
	}
	close(jobs)

	// Collect results with progress tracking
	results := make([]*models.ComponentResult, len(enabledComponents))
	var firstError error

	for i := 0; i < len(enabledComponents); i++ {
		jobResult := <-jobResults
		results[jobResult.index] = jobResult.result

		// Update progress
		componentName := enabledComponents[jobResult.index].Name
		if jobResult.err == nil {
			progressTracker.Increment(fmt.Sprintf("✓ %s", componentName))
		} else {
			progressTracker.Increment(fmt.Sprintf("✗ %s", componentName))
		}

		// Track first error
		if jobResult.err != nil && firstError == nil {
			firstError = jobResult.err
		}
	}
	close(jobResults)

	// Complete progress tracking
	progressTracker.Complete()

	// Return error if any component failed
	if firstError != nil {
		return results, firstError
	}

	pc.logger.Info("All components generated successfully")
	return results, nil
}

// generateComponent generates a single component with retry and fallback logic
func (pc *ProjectCoordinator) generateComponent(ctx context.Context, comp *models.ComponentConfig, config *models.ProjectConfig, toolCheckResult *models.ToolCheckResult) (*models.ComponentResult, error) {
	startTime := time.Now()

	// Create logger with component context
	componentLogger := pc.logger.WithContext("component", comp.Name).WithContext("type", comp.Type)
	componentLogger.Info("Starting component generation")

	result := &models.ComponentResult{
		Type:        comp.Type,
		Name:        comp.Name,
		Success:     false,
		ManualSteps: make([]string, 0),
		Warnings:    make([]string, 0),
	}

	// Create error context for recovery decisions
	errCtx := &ErrorContext{
		Operation:     "component_generation",
		Component:     comp.Name,
		Phase:         "generation",
		AttemptNumber: 0,
		CanRetry:      true,
		CanFallback:   pc.fallbackRegistry.Supports(comp.Type),
	}

	// Determine generation strategy
	useBootstrap := false
	toolsForComponent := pc.toolDiscovery.GetToolsForComponent(comp.Type)

	// Check if we should use fallback (offline mode or missing tools)
	shouldUseFallback, reason := pc.toolDiscovery.ShouldUseFallback(comp.Type)

	if shouldUseFallback {
		pc.logger.Info(fmt.Sprintf("Using fallback generation for %s: %s", comp.Type, reason))
		useBootstrap = false
	} else if config.Options.UseExternalTools && len(toolsForComponent) > 0 {
		// Check if all required tools are available
		allToolsAvailable := true
		for _, tool := range toolsForComponent {
			if toolInfo, exists := toolCheckResult.Tools[tool]; !exists || !toolInfo.Available {
				allToolsAvailable = false
				break
			}
		}
		useBootstrap = allToolsAvailable
	}

	var lastErr error

	// Try bootstrap generation with retry logic
	if useBootstrap {
		componentLogger.Info(fmt.Sprintf("Using bootstrap tool: %s", toolsForComponent[0]))
		result.Method = "bootstrap"
		result.ToolUsed = toolsForComponent[0]

		// Attempt bootstrap with retry
		for errCtx.AttemptNumber = 1; errCtx.AttemptNumber <= 2; errCtx.AttemptNumber++ {
			if errCtx.AttemptNumber > 1 {
				componentLogger.Info(fmt.Sprintf("Retrying component generation (attempt %d/2)", errCtx.AttemptNumber))
			} else {
				componentLogger.Debug(fmt.Sprintf("Executing bootstrap tool (attempt %d/2)", errCtx.AttemptNumber))
			}

			err := pc.generateWithBootstrap(ctx, comp, config, result)
			if err == nil {
				// Success
				result.Success = true
				result.Duration = time.Since(startTime)
				componentLogger.Info(fmt.Sprintf("Component generated successfully in %v using bootstrap tool", result.Duration))
				return result, nil
			}

			lastErr = err

			// Convert to GenerationError if not already
			var genErr *GenerationError
			if !errors.As(err, &genErr) {
				genErr = NewToolExecutionError(toolsForComponent[0], comp.Name, err)
			}

			// Check if we should retry
			if !ShouldRetry(genErr, errCtx) {
				componentLogger.Warn(fmt.Sprintf("Bootstrap generation failed, not retrying: %v", err))
				break
			}

			componentLogger.Warn(fmt.Sprintf("Bootstrap generation failed (attempt %d): %v", errCtx.AttemptNumber, err))
		}

		// Check if we should fallback after bootstrap failure
		var genErr *GenerationError
		if errors.As(lastErr, &genErr) && ShouldFallback(genErr, errCtx) {
			pc.logger.Info(fmt.Sprintf("Falling back to custom generation for %s", comp.Type))
			useBootstrap = false
		} else {
			// No fallback, return error
			return result, lastErr
		}
	}

	// Try fallback generation if bootstrap failed or wasn't available
	if !useBootstrap && pc.fallbackRegistry.Supports(comp.Type) {
		componentLogger.Info("Using fallback generator (bootstrap tool unavailable or failed)")
		result.Method = "fallback"
		result.ToolUsed = "fallback-generator"

		if err := pc.generateWithFallback(ctx, comp, config, result); err != nil {
			componentLogger.Error(fmt.Sprintf("Fallback generation failed: %v", err))
			return result, err
		}

		result.Success = true
		result.Duration = time.Since(startTime)
		componentLogger.Info(fmt.Sprintf("Component generated successfully in %v using fallback generator", result.Duration))
		if len(result.ManualSteps) > 0 {
			componentLogger.Warn(fmt.Sprintf("Manual setup required: %d steps", len(result.ManualSteps)))
		}
		return result, nil
	}

	// No generation method available
	if lastErr != nil {
		return result, lastErr
	}

	return result, NewGenerationError(
		ErrCategoryToolNotFound,
		fmt.Sprintf("no generation method available for component type: %s", comp.Type),
		nil,
	).WithSuggestions(
		"Install required tools for this component type",
		"Check if fallback generation is available",
	)
}

// generateWithBootstrap generates a component using a bootstrap tool
func (pc *ProjectCoordinator) generateWithBootstrap(ctx context.Context, comp *models.ComponentConfig, config *models.ProjectConfig, result *models.ComponentResult) error {
	// Create temporary directory for generation
	tempDir := filepath.Join(config.OutputDir, ".temp", comp.Name)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Register for cleanup
	pc.rollbackManager.RegisterTempDir(tempDir)

	// Build bootstrap spec
	spec := &bootstrap.BootstrapSpec{
		ComponentType: comp.Type,
		TargetDir:     tempDir,
		Config:        comp.Config,
		Flags:         []string{},
		Timeout:       5 * time.Minute,
	}

	// Get appropriate executor
	executor, err := pc.getBootstrapExecutor(comp.Type)
	if err != nil {
		return err
	}

	var execResult *models.ExecutionResult

	// Check if streaming output is enabled
	if config.Options.StreamOutput && config.Options.Verbose {
		// Use streaming execution with progress indicator
		pc.logger.Info(fmt.Sprintf("Executing bootstrap tool for %s (streaming output)...", comp.Name))

		// Create streaming writer
		streamWriter := NewStreamingWriter(pc.logger.GetWriter(), comp.Name, true)

		// Execute with streaming
		execResult, err = executor.ExecuteWithStreaming(ctx, spec, streamWriter)

		// Flush any remaining output
		streamWriter.Flush()

		if err != nil {
			return fmt.Errorf("bootstrap execution failed: %w", err)
		}
	} else {
		// Use standard execution with progress indicator
		progress := NewProgressIndicator(pc.logger.GetWriter(), fmt.Sprintf("Generating %s...", comp.Name), !config.Options.Verbose)
		progress.Start()

		// Execute bootstrap tool
		execResult, err = executor.Execute(ctx, spec)

		progress.Stop()

		if err != nil {
			return fmt.Errorf("bootstrap execution failed: %w", err)
		}

		// Log output if verbose
		if config.Options.Verbose && execResult.Stdout != "" {
			pc.logger.Info(fmt.Sprintf("Bootstrap output:\n%s", execResult.Stdout))
		}
	}

	// Check execution result
	if !execResult.Success {
		return fmt.Errorf("bootstrap tool failed with exit code %d", execResult.ExitCode)
	}

	// Set output path
	result.OutputPath = tempDir

	return nil
}

// generateWithFallback generates a component using fallback generator
func (pc *ProjectCoordinator) generateWithFallback(ctx context.Context, comp *models.ComponentConfig, config *models.ProjectConfig, result *models.ComponentResult) error {
	// Get fallback generator
	generator, err := pc.fallbackRegistry.Get(comp.Type)
	if err != nil {
		return err
	}

	// Create temporary directory for generation
	tempDir := filepath.Join(config.OutputDir, ".temp", comp.Name)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Register for cleanup
	pc.rollbackManager.RegisterTempDir(tempDir)

	// Build fallback spec
	spec := &models.FallbackSpec{
		ComponentType: comp.Type,
		TargetDir:     tempDir,
		Config:        comp.Config,
	}

	// Generate using fallback
	fallbackResult, err := generator.Generate(ctx, spec)
	if err != nil {
		return fmt.Errorf("fallback generation failed: %w", err)
	}

	// Get manual steps
	result.ManualSteps = generator.GetRequiredManualSteps(comp.Type)
	if len(result.ManualSteps) > 0 {
		result.Warnings = append(result.Warnings, "Manual setup steps required after generation")
	}

	// Set output path
	result.OutputPath = fallbackResult.OutputPath

	return nil
}

// getBootstrapExecutor returns the appropriate bootstrap executor for a component type
func (pc *ProjectCoordinator) getBootstrapExecutor(componentType string) (*bootstrap.BaseExecutor, error) {
	// For now, return a base executor
	// In a full implementation, this would return type-specific executors
	switch componentType {
	case "nextjs":
		return bootstrap.NewBaseExecutor("npx"), nil
	case "go-backend":
		return bootstrap.NewBaseExecutor("go"), nil
	case "android":
		return bootstrap.NewBaseExecutor("gradle"), nil
	case "ios":
		return bootstrap.NewBaseExecutor("xcodebuild"), nil
	default:
		return nil, fmt.Errorf("no bootstrap executor for component type: %s", componentType)
	}
}

// mapComponentsToStructure maps generated components to the target structure
func (pc *ProjectCoordinator) mapComponentsToStructure(ctx context.Context, config *models.ProjectConfig, results []*models.ComponentResult) error {
	for _, result := range results {
		if !result.Success {
			continue
		}

		pc.logger.Info(fmt.Sprintf("Mapping component %s to target structure", result.Name))

		// Map to target structure
		if err := pc.structureMapper.Map(ctx, result.OutputPath, config.OutputDir, result.Type); err != nil {
			return fmt.Errorf("failed to map component %s: %w", result.Name, err)
		}

		// Update output path to final location
		result.OutputPath = filepath.Join(config.OutputDir, pc.structureMapper.GetTargetPath(result.Type))
	}

	return nil
}

// integrateComponents integrates generated components
func (pc *ProjectCoordinator) integrateComponents(ctx context.Context, config *models.ProjectConfig, results []*models.ComponentResult) error {
	// Build component list for integration
	components := make([]*models.Component, 0)
	for _, result := range results {
		if !result.Success {
			continue
		}

		component := &models.Component{
			Type:        result.Type,
			Name:        result.Name,
			Path:        result.OutputPath,
			Config:      make(map[string]interface{}),
			GeneratedAt: time.Now(),
		}
		components = append(components, component)
	}

	// Create integration manager
	integrationMgr := NewIntegrationManager(config.OutputDir, config.Options.Verbose)

	// Integrate components
	if err := integrationMgr.Integrate(ctx, components, &config.Integration); err != nil {
		return err
	}

	return nil
}

// validateFinalStructure validates the final project structure
func (pc *ProjectCoordinator) validateFinalStructure(config *models.ProjectConfig, results []*models.ComponentResult) error {
	// Collect component types that were successfully generated
	componentTypes := make([]string, 0)
	for _, result := range results {
		if result.Success {
			componentTypes = append(componentTypes, result.Type)
		}
	}

	// Validate structure
	if err := pc.structureMapper.ValidateStructureWithComponents(config.OutputDir, componentTypes); err != nil {
		return err
	}

	pc.logger.Info("Final structure validation passed")

	return nil
}

// getExpectedFiles returns expected files for a component type (for dry-run)
func (pc *ProjectCoordinator) getExpectedFiles(componentType string) []string {
	switch componentType {
	case "nextjs":
		return []string{
			"package.json",
			"next.config.js",
			"tsconfig.json",
			"app/page.tsx",
			"app/layout.tsx",
		}
	case "go-backend":
		return []string{
			"go.mod",
			"go.sum",
			"main.go",
		}
	case "android":
		return []string{
			"build.gradle",
			"settings.gradle",
			"app/build.gradle",
			"app/src/main/AndroidManifest.xml",
		}
	case "ios":
		return []string{
			"Package.swift",
			"*.xcodeproj",
		}
	default:
		return []string{}
	}
}

// runSecurityScan runs a security scan on the generated project
func (pc *ProjectCoordinator) runSecurityScan(ctx context.Context, projectRoot string) (*security.ScanResult, error) {
	scanner := security.NewSecurityScanner()

	pc.logger.Info("Scanning for security issues...")
	scanResult, err := scanner.Scan(ctx, projectRoot)
	if err != nil {
		return nil, fmt.Errorf("security scan failed: %w", err)
	}

	return scanResult, nil
}
