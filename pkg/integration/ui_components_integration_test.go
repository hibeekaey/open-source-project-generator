package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TestUIComponentsIntegration tests the integration of refactored UI components
func TestUIComponentsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := t.TempDir()

	t.Run("config_components_integration", func(t *testing.T) {
		testConfigComponentsIntegration(t, tempDir)
	})

	t.Run("preview_components_integration", func(t *testing.T) {
		testPreviewComponentsIntegration(t, tempDir)
	})

	t.Run("interactive_ui_workflow", func(t *testing.T) {
		testInteractiveUIWorkflow(t, tempDir)
	})

	t.Run("ui_component_coordination", func(t *testing.T) {
		testUIComponentCoordination(t, tempDir)
	})

	t.Run("ui_error_handling_integration", func(t *testing.T) {
		testUIErrorHandlingIntegration(t, tempDir)
	})
}

func testConfigComponentsIntegration(t *testing.T, tempDir string) {
	// Test config manager integration
	t.Run("config_manager", func(t *testing.T) {
		configManager := NewMockUIConfigManager()

		// Test configuration collection workflow
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"project_name":    "config-integration-test",
				"organization":    "config-org",
				"description":     "Test project for config integration",
				"license":         "MIT",
				"backend_enabled": true,
				"backend_type":    "go-gin",
			},
		}

		projectConfig, err := configManager.CollectConfiguration(ctx, mockInput)
		if err != nil {
			t.Errorf("Config collection failed: %v", err)
		}

		if projectConfig == nil {
			t.Fatal("Expected project config to be returned")
		}

		if projectConfig.Name != "config-integration-test" {
			t.Errorf("Expected project name 'config-integration-test', got '%s'", projectConfig.Name)
		}

		// Test configuration validation
		err = configManager.ValidateConfiguration(projectConfig)
		if err != nil {
			t.Errorf("Config validation failed: %v", err)
		}

		// Test configuration export
		configPath := filepath.Join(tempDir, "exported-config.yaml")
		err = configManager.ExportConfiguration(projectConfig, configPath)
		if err != nil {
			t.Errorf("Config export failed: %v", err)
		}

		// Verify config file was created
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Expected config file to be created")
		}

		// Test configuration import
		importedConfig, err := configManager.ImportConfiguration(configPath)
		if err != nil {
			t.Errorf("Config import failed: %v", err)
		}

		if importedConfig.Name != projectConfig.Name {
			t.Errorf("Expected imported config name '%s', got '%s'", projectConfig.Name, importedConfig.Name)
		}
	})

	// Test config collectors integration
	t.Run("config_collectors", func(t *testing.T) {
		// Test project info collector
		projectCollector := NewMockProjectInfoCollector()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"project_name": "collector-test-project",
				"organization": "collector-org",
				"description":  "Test project for collector integration",
				"license":      "Apache-2.0",
				"version":      "1.0.0",
			},
		}

		projectInfo, err := projectCollector.CollectProjectInfo(ctx, mockInput)
		if err != nil {
			t.Errorf("Project info collection failed: %v", err)
		}

		if projectInfo.Name != "collector-test-project" {
			t.Errorf("Expected project name 'collector-test-project', got '%s'", projectInfo.Name)
		}

		if projectInfo.License != "Apache-2.0" {
			t.Errorf("Expected license 'Apache-2.0', got '%s'", projectInfo.License)
		}

		// Test component collector
		componentCollector := NewMockComponentCollector()

		componentConfig, err := componentCollector.CollectComponents(ctx, mockInput)
		if err != nil {
			t.Errorf("Component collection failed: %v", err)
		}

		if componentConfig == nil {
			t.Error("Expected component config to be returned")
		}

		// Test advanced options collector
		advancedCollector := NewMockAdvancedOptionsCollector()

		advancedOptions, err := advancedCollector.CollectAdvancedOptions(ctx, mockInput)
		if err != nil {
			t.Errorf("Advanced options collection failed: %v", err)
		}

		if advancedOptions == nil {
			t.Error("Expected advanced options to be returned")
		}
	})

	// Test config validators integration
	t.Run("config_validators", func(t *testing.T) {
		// Test project validator
		projectValidator := NewMockProjectValidator()

		validProject := &models.ProjectConfig{
			Name:         "valid-project",
			Organization: "valid-org",
			License:      "MIT",
		}

		err := projectValidator.ValidateProject(validProject)
		if err != nil {
			t.Errorf("Valid project validation failed: %v", err)
		}

		invalidProject := &models.ProjectConfig{
			Name: "", // Invalid empty name
		}

		err = projectValidator.ValidateProject(invalidProject)
		if err == nil {
			t.Error("Expected invalid project to fail validation")
		}

		// Test component validator
		componentValidator := NewMockComponentValidator()

		validComponents := &models.Components{
			Backend: models.BackendComponents{
				GoGin: true,
			},
		}

		err = componentValidator.ValidateComponents(validComponents)
		if err != nil {
			t.Errorf("Valid components validation failed: %v", err)
		}

		// Test conflicting components
		conflictingComponents := &models.Components{
			Backend: models.BackendComponents{
				GoGin: true,
			},
		}

		err = componentValidator.ValidateComponents(conflictingComponents)
		if err == nil {
			t.Error("Expected conflicting components to fail validation")
		}
	})

	// Test config formatters integration
	t.Run("config_formatters", func(t *testing.T) {
		config := &models.ProjectConfig{
			Name:         "formatter-test-project",
			Organization: "formatter-org",
			Description:  "Test project for formatter integration",
			License:      "MIT",
			Components: models.Components{
				Backend: models.BackendComponents{
					GoGin: true,
				},
				Frontend: models.FrontendComponents{
					NextJS: models.NextJSComponents{
						App: true,
					},
				},
			},
		}

		// Test summary formatter
		summaryFormatter := NewMockSummaryFormatter()

		summary, err := summaryFormatter.FormatSummary(config)
		if err != nil {
			t.Errorf("Summary formatting failed: %v", err)
		}

		if summary == "" {
			t.Error("Expected summary to be generated")
		}

		// Test export formatter
		exportFormatter := NewMockExportFormatter()

		yamlData, err := exportFormatter.FormatYAML(config)
		if err != nil {
			t.Errorf("YAML formatting failed: %v", err)
		}

		if len(yamlData) == 0 {
			t.Error("Expected YAML data to be generated")
		}

		jsonData, err := exportFormatter.FormatJSON(config)
		if err != nil {
			t.Errorf("JSON formatting failed: %v", err)
		}

		if len(jsonData) == 0 {
			t.Error("Expected JSON data to be generated")
		}
	})
}

func testPreviewComponentsIntegration(t *testing.T, tempDir string) {
	// Test preview manager integration
	t.Run("preview_manager", func(t *testing.T) {
		previewManager := NewMockPreviewManager()

		config := &models.ProjectConfig{
			Name:         "preview-test-project",
			Organization: "preview-org",
			Description:  "Test project for preview integration",
			License:      "MIT",
			OutputPath:   filepath.Join(tempDir, "preview-output"),
			Components: models.Components{
				Backend: models.BackendComponents{
					GoGin: true,
				},
				Frontend: models.FrontendComponents{
					NextJS: models.NextJSComponents{
						App: true,
					},
				},
			},
		}

		// Test project structure preview generation
		structurePreview, err := previewManager.GenerateStructurePreview(config)
		if err != nil {
			t.Errorf("Structure preview generation failed: %v", err)
		}

		if structurePreview == nil {
			t.Error("Expected structure preview to be generated")
		}

		// Test component preview generation
		componentPreview, err := previewManager.GenerateComponentPreview(config)
		if err != nil {
			t.Errorf("Component preview generation failed: %v", err)
		}

		if componentPreview == nil {
			t.Error("Expected component preview to be generated")
		}

		// Test preview display
		err = previewManager.DisplayPreview(structurePreview)
		if err != nil {
			t.Errorf("Preview display failed: %v", err)
		}
	})

	// Test tree components integration
	t.Run("tree_components", func(t *testing.T) {
		// Test tree builder
		treeBuilder := NewMockTreeBuilder()

		config := &models.ProjectConfig{
			Name:       "tree-test-project",
			OutputPath: filepath.Join(tempDir, "tree-output"),
			Components: models.Components{
				Backend: models.BackendComponents{
					GoGin: true,
				},
			},
		}

		tree, err := treeBuilder.BuildProjectTree(config)
		if err != nil {
			t.Errorf("Tree building failed: %v", err)
		}

		if tree == nil {
			t.Error("Expected project tree to be built")
		}

		// Test tree renderer
		treeRenderer := NewMockTreeRenderer()

		renderedTree, err := treeRenderer.RenderTree(tree)
		if err != nil {
			t.Errorf("Tree rendering failed: %v", err)
		}

		if renderedTree == "" {
			t.Error("Expected rendered tree to be generated")
		}

		// Test tree formatter
		treeFormatter := NewMockTreeFormatter()

		formattedTree, err := treeFormatter.FormatTree(tree, "ascii")
		if err != nil {
			t.Errorf("Tree formatting failed: %v", err)
		}

		if formattedTree == "" {
			t.Error("Expected formatted tree to be generated")
		}

		// Test different format
		unicodeTree, err := treeFormatter.FormatTree(tree, "unicode")
		if err != nil {
			t.Errorf("Unicode tree formatting failed: %v", err)
		}

		if unicodeTree == "" {
			t.Error("Expected unicode tree to be generated")
		}
	})

	// Test component preview integration
	t.Run("component_previews", func(t *testing.T) {
		// Test frontend preview
		frontendPreview := NewMockFrontendPreview()

		frontendConfig := &models.FrontendComponents{
			NextJS: models.NextJSComponents{
				App:   true,
				Admin: true,
			},
		}

		frontendPreviewData, err := frontendPreview.GeneratePreview(frontendConfig)
		if err != nil {
			t.Errorf("Frontend preview generation failed: %v", err)
		}

		if frontendPreviewData == nil {
			t.Error("Expected frontend preview data to be generated")
		}

		// Test backend preview
		backendPreview := NewMockBackendPreview()

		backendConfig := &models.BackendComponents{
			GoGin: true,
		}

		backendPreviewData, err := backendPreview.GeneratePreview(backendConfig)
		if err != nil {
			t.Errorf("Backend preview generation failed: %v", err)
		}

		if backendPreviewData == nil {
			t.Error("Expected backend preview data to be generated")
		}

		// Test infrastructure preview
		infraPreview := NewMockInfrastructurePreview()

		infraConfig := &models.InfrastructureComponents{
			Docker:     true,
			Kubernetes: true,
		}

		infraPreviewData, err := infraPreview.GeneratePreview(infraConfig)
		if err != nil {
			t.Errorf("Infrastructure preview generation failed: %v", err)
		}

		if infraPreviewData == nil {
			t.Error("Expected infrastructure preview data to be generated")
		}
	})

	// Test display components integration
	t.Run("display_components", func(t *testing.T) {
		// Test console display
		consoleDisplay := NewMockConsoleDisplay()

		previewData := &interfaces.PreviewData{
			Content: "Test Preview Content",
			Type:    "project_structure",
			Metadata: map[string]interface{}{
				"title":       "Test Preview",
				"description": "Test preview for display integration",
			},
		}

		err := consoleDisplay.DisplayPreview(previewData)
		if err != nil {
			t.Errorf("Console display failed: %v", err)
		}

		// Test interactive display
		interactiveDisplay := NewMockInteractiveDisplay()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"navigate": "expand",
				"action":   "view",
			},
		}

		err = interactiveDisplay.DisplayInteractivePreview(ctx, previewData, mockInput)
		if err != nil {
			t.Errorf("Interactive display failed: %v", err)
		}
	})
}

func testInteractiveUIWorkflow(t *testing.T, tempDir string) {
	// Test complete interactive UI workflow
	t.Run("complete_workflow", func(t *testing.T) {
		// Create UI manager
		uiManager := NewMockUIManager()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"project_name":     "workflow-test-project",
				"organization":     "workflow-org",
				"description":      "Complete workflow test project",
				"license":          "MIT",
				"backend_enabled":  true,
				"backend_type":     "go-gin",
				"frontend_enabled": true,
				"frontend_type":    "nextjs-app",
				"mobile_enabled":   false,
				"infra_enabled":    true,
				"infra_type":       "docker",
				"confirm_config":   true,
				"generate_project": true,
			},
		}

		// Step 1: Collect project configuration
		projectConfig, err := uiManager.CollectProjectConfiguration(ctx, mockInput)
		if err != nil {
			t.Errorf("Project configuration collection failed: %v", err)
		}

		if projectConfig == nil {
			t.Fatal("Expected project config to be returned")
		}

		// Step 2: Generate and display preview
		preview, err := uiManager.GenerateProjectPreview(projectConfig)
		if err != nil {
			t.Errorf("Project preview generation failed: %v", err)
		}

		err = uiManager.DisplayPreview(preview)
		if err != nil {
			t.Errorf("Preview display failed: %v", err)
		}

		// Step 3: Confirm configuration
		confirmed, err := uiManager.ConfirmConfiguration(ctx, projectConfig, mockInput)
		if err != nil {
			t.Errorf("Configuration confirmation failed: %v", err)
		}

		if !confirmed {
			t.Error("Expected configuration to be confirmed")
		}

		// Step 4: Display generation progress
		progressChan := make(chan interfaces.ProgressUpdate, 10)
		go func() {
			// Simulate progress updates
			progressChan <- interfaces.ProgressUpdate{
				Step:     1,
				Total:    2,
				Progress: 10.0,
				Message:  "Setting up project structure",
			}
			progressChan <- interfaces.ProgressUpdate{
				Step:     1,
				Total:    2,
				Progress: 50.0,
				Message:  "Creating backend components",
			}
			progressChan <- interfaces.ProgressUpdate{
				Step:     2,
				Total:    2,
				Progress: 100.0,
				Message:  "Project generation complete",
			}
			close(progressChan)
		}()

		err = uiManager.DisplayProgress(ctx, progressChan)
		if err != nil {
			t.Errorf("Progress display failed: %v", err)
		}

		// Step 5: Display completion summary
		summary := &interfaces.GenerationSummary{
			FilesCreated: 25,
			Components:   []string{"backend", "frontend", "infrastructure"},
			Duration:     2 * time.Second,
			Metadata: map[string]interface{}{
				"project_name": projectConfig.Name,
				"output_path":  projectConfig.OutputPath,
			},
		}

		err = uiManager.DisplayCompletionSummary(summary)
		if err != nil {
			t.Errorf("Completion summary display failed: %v", err)
		}
	})

	// Test error handling workflow
	t.Run("error_handling_workflow", func(t *testing.T) {
		uiManager := NewMockUIManager()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Test validation error handling
		invalidConfig := &models.ProjectConfig{
			Name: "", // Invalid empty name
		}

		err := uiManager.ValidateConfiguration(invalidConfig)
		if err == nil {
			t.Error("Expected validation to fail for invalid config")
		}

		// Test error display
		validationError := &interfaces.ValidationError{
			Field:   "name",
			Message: "Project name cannot be empty",
			Code:    "REQUIRED_FIELD",
		}

		err = uiManager.DisplayValidationError(validationError)
		if err != nil {
			t.Errorf("Validation error display failed: %v", err)
		}

		// Test error recovery
		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"retry":        true,
				"project_name": "recovered-project",
			},
		}

		recovered, err := uiManager.HandleValidationError(ctx, validationError, mockInput)
		if err != nil {
			t.Errorf("Validation error handling failed: %v", err)
		}

		if !recovered {
			t.Error("Expected error recovery to succeed")
		}
	})

	// Test navigation workflow
	t.Run("navigation_workflow", func(t *testing.T) {
		uiManager := NewMockUIManager()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"menu_choice":   "project_setup",
				"navigate_back": false,
				"navigate_next": true,
				"navigate_exit": false,
			},
		}

		// Test main menu navigation
		choice, err := uiManager.DisplayMainMenu(ctx, mockInput)
		if err != nil {
			t.Errorf("Main menu display failed: %v", err)
		}

		if choice != "project_setup" {
			t.Errorf("Expected menu choice 'project_setup', got '%s'", choice)
		}

		// Test navigation controls
		navigation, err := uiManager.HandleNavigation(ctx, mockInput)
		if err != nil {
			t.Errorf("Navigation handling failed: %v", err)
		}

		if navigation.Action != "next" {
			t.Errorf("Expected navigation action 'next', got '%s'", navigation.Action)
		}

		// Test breadcrumb navigation
		breadcrumbs := []string{"Home", "Project Setup", "Components"}
		err = uiManager.DisplayBreadcrumbs(breadcrumbs)
		if err != nil {
			t.Errorf("Breadcrumb display failed: %v", err)
		}
	})
}

func testUIComponentCoordination(t *testing.T, tempDir string) {
	// Test coordination between config and preview components
	t.Run("config_preview_coordination", func(t *testing.T) {
		configManager := NewMockUIConfigManager()
		previewManager := NewMockPreviewManager()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"project_name":    "coordination-test",
				"organization":    "coordination-org",
				"backend_enabled": true,
				"backend_type":    "go-gin",
			},
		}

		// Collect configuration
		projectConfig, err := configManager.CollectConfiguration(ctx, mockInput)
		if err != nil {
			t.Errorf("Configuration collection failed: %v", err)
		}

		// Generate preview from configuration
		structurePreview, err := previewManager.GenerateStructurePreview(projectConfig)
		if err != nil {
			t.Errorf("Preview generation from config failed: %v", err)
		}

		if structurePreview == nil {
			t.Error("Expected structure preview to be generated from config")
		}

		// Update configuration based on preview feedback
		mockInput.responses["modify_structure"] = true
		mockInput.responses["add_component"] = "frontend"

		updatedConfig, err := configManager.UpdateConfigurationFromPreview(projectConfig, structurePreview, mockInput)
		if err != nil {
			t.Errorf("Configuration update from preview failed: %v", err)
		}

		if updatedConfig == nil {
			t.Error("Expected updated configuration to be returned")
		}
	})

	// Test coordination between different UI managers
	t.Run("ui_managers_coordination", func(t *testing.T) {
		uiManager := NewMockUIManager()
		configManager := NewMockUIConfigManager()
		previewManager := NewMockPreviewManager()

		// Register managers with UI coordinator
		err := uiManager.RegisterConfigManager(configManager)
		if err != nil {
			t.Errorf("Config manager registration failed: %v", err)
		}

		err = uiManager.RegisterPreviewManager(previewManager)
		if err != nil {
			t.Errorf("Preview manager registration failed: %v", err)
		}

		// Test coordinated workflow
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"workflow_type": "guided",
				"project_name":  "coordinated-project",
			},
		}

		result, err := uiManager.ExecuteCoordinatedWorkflow(ctx, "project_creation", mockInput)
		if err != nil {
			t.Errorf("Coordinated workflow execution failed: %v", err)
		}

		if result == nil {
			t.Error("Expected workflow result to be returned")
		}
	})

	// Test component state synchronization
	t.Run("component_state_sync", func(t *testing.T) {
		// Create state manager
		stateManager := NewMockStateManager()

		// Create components with shared state
		configComponent := NewMockManagerWithState(stateManager)
		previewComponent := NewMockManagerWithState(stateManager)

		// Update state in config component
		state := &interfaces.UIState{
			CurrentStep:   "component_selection",
			ProjectConfig: &models.ProjectConfig{Name: "sync-test"},
			Data: map[string]interface{}{
				"validation_state": map[string]interface{}{
					"valid": true,
				},
			},
		}

		err := configComponent.UpdateState(state)
		if err != nil {
			t.Errorf("Config component state update failed: %v", err)
		}

		// Verify state is synchronized in preview component
		syncedState, err := previewComponent.GetState()
		if err != nil {
			t.Errorf("Preview component state retrieval failed: %v", err)
		}

		if syncedState.CurrentStep != state.CurrentStep {
			t.Errorf("Expected synced step '%s', got '%s'", state.CurrentStep, syncedState.CurrentStep)
		}

		stateConfig := state.ProjectConfig.(*models.ProjectConfig)
		syncedConfig := syncedState.ProjectConfig.(*models.ProjectConfig)
		if syncedConfig.Name != stateConfig.Name {
			t.Errorf("Expected synced project name '%s', got '%s'", stateConfig.Name, syncedConfig.Name)
		}
	})
}

func testUIErrorHandlingIntegration(t *testing.T, tempDir string) {
	// Test comprehensive error handling across UI components
	t.Run("validation_error_handling", func(t *testing.T) {
		uiManager := NewMockUIManager()

		// Test validation error display and recovery
		validationErrors := []interfaces.ValidationError{
			{
				Field:   "name",
				Message: "Project name must be at least 3 characters",
				Code:    "MIN_LENGTH",
			},
			{
				Field:   "organization",
				Message: "Organization name contains invalid characters",
				Code:    "INVALID_CHARS",
			},
		}

		err := uiManager.DisplayValidationErrors(validationErrors)
		if err != nil {
			t.Errorf("Validation errors display failed: %v", err)
		}

		// Test error correction workflow
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"correct_name":         "corrected-project-name",
				"correct_organization": "corrected-org",
				"retry_validation":     true,
			},
		}

		corrected, err := uiManager.HandleValidationErrors(ctx, validationErrors, mockInput)
		if err != nil {
			t.Errorf("Validation error handling failed: %v", err)
		}

		if !corrected {
			t.Error("Expected validation errors to be corrected")
		}
	})

	// Test generation error handling
	t.Run("generation_error_handling", func(t *testing.T) {
		uiManager := NewMockUIManager()

		// Test generation error display
		generationError := &interfaces.GenerationError{
			Type:      "PERMISSION_ERROR",
			Message:   "Failed to create Go module: permission denied",
			Component: "backend_generation",
			Metadata: map[string]interface{}{
				"cause": "insufficient file system permissions",
			},
		}

		err := uiManager.DisplayGenerationError(generationError)
		if err != nil {
			t.Errorf("Generation error display failed: %v", err)
		}

		// Test error recovery options
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"recovery_action":  "change_output_path",
				"new_output_path":  filepath.Join(tempDir, "recovery-output"),
				"retry_generation": true,
			},
		}

		recovered, err := uiManager.HandleGenerationError(ctx, generationError, mockInput)
		if err != nil {
			t.Errorf("Generation error handling failed: %v", err)
		}

		if !recovered {
			t.Error("Expected generation error recovery to succeed")
		}
	})

	// Test network error handling
	t.Run("network_error_handling", func(t *testing.T) {
		uiManager := NewMockUIManager()

		// Test network error display
		networkError := &interfaces.NetworkError{
			Type:    "TIMEOUT",
			URL:     "https://example.com/templates/go-gin.zip",
			Message: "Connection timeout",
			Timeout: true,
			Metadata: map[string]interface{}{
				"operation": "template_download",
			},
		}

		err := uiManager.DisplayNetworkError(networkError)
		if err != nil {
			t.Errorf("Network error display failed: %v", err)
		}

		// Test offline mode fallback
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		mockInput := &MockUIInput{
			responses: map[string]interface{}{
				"use_offline_mode":     true,
				"use_cached_templates": true,
			},
		}

		fallback, err := uiManager.HandleNetworkError(ctx, networkError, mockInput)
		if err != nil {
			t.Errorf("Network error handling failed: %v", err)
		}

		if !fallback {
			t.Error("Expected network error fallback to succeed")
		}
	})

	// Test error aggregation and reporting
	t.Run("error_aggregation", func(t *testing.T) {
		uiManager := NewMockUIManager()

		// Collect multiple errors
		errors := []interfaces.UIError{
			interfaces.UIError{
				Type:    "VALIDATION_ERROR",
				Message: "Invalid project name",
				Metadata: map[string]interface{}{
					"field": "name",
					"code":  "INVALID_NAME",
				},
			},
			interfaces.UIError{
				Type:    "GENERATION_ERROR",
				Message: "Template not found",
				Metadata: map[string]interface{}{
					"stage": "frontend_generation",
					"code":  "TEMPLATE_NOT_FOUND",
				},
			},
			interfaces.UIError{
				Type:    "NETWORK_ERROR",
				Message: "Network unreachable",
				Metadata: map[string]interface{}{
					"operation": "dependency_check",
					"code":      "NETWORK_ERROR",
				},
			},
		}

		// Test error aggregation
		aggregatedReport, err := uiManager.AggregateErrors(errors)
		if err != nil {
			t.Errorf("Error aggregation failed: %v", err)
		}

		if aggregatedReport == nil {
			t.Error("Expected aggregated error report to be generated")
		}

		errorCount := aggregatedReport.Metadata["error_count"].(int)
		if errorCount != len(errors) {
			t.Errorf("Expected %d errors in report, got %d", len(errors), errorCount)
		}

		// Test error report display
		err = uiManager.DisplayErrorReport(aggregatedReport)
		if err != nil {
			t.Errorf("Error report display failed: %v", err)
		}

		// Test error report export
		reportPath := filepath.Join(tempDir, "error-report.json")
		err = uiManager.ExportErrorReport(aggregatedReport, reportPath)
		if err != nil {
			t.Errorf("Error report export failed: %v", err)
		}

		// Verify report file was created
		if _, err := os.Stat(reportPath); os.IsNotExist(err) {
			t.Error("Expected error report file to be created")
		}
	})
}

// Mock implementations for UI testing

type MockUIInput struct {
	responses map[string]interface{}
}

func (m *MockUIInput) GetStringInput(prompt string, key string) (string, error) {
	if response, exists := m.responses[key]; exists {
		if str, ok := response.(string); ok {
			return str, nil
		}
	}
	return "", nil
}

func (m *MockUIInput) GetBoolInput(prompt string, key string) (bool, error) {
	if response, exists := m.responses[key]; exists {
		if b, ok := response.(bool); ok {
			return b, nil
		}
		if str, ok := response.(string); ok {
			return str == "true" || str == "yes" || str == "y", nil
		}
	}
	return false, nil
}

func (m *MockUIInput) GetSelectInput(prompt string, key string, options []string) (string, error) {
	if response, exists := m.responses[key]; exists {
		if str, ok := response.(string); ok {
			// Verify response is in options
			for _, option := range options {
				if option == str {
					return str, nil
				}
			}
		}
	}
	// Return first option as default
	if len(options) > 0 {
		return options[0], nil
	}
	return "", nil
}

func (m *MockUIInput) GetMultiSelectInput(prompt string, key string, options []string) ([]string, error) {
	if response, exists := m.responses[key]; exists {
		if str, ok := response.(string); ok {
			return []string{str}, nil
		}
		if slice, ok := response.([]string); ok {
			return slice, nil
		}
	}
	return []string{}, nil
}

func (m *MockUIInput) GetIntInput(prompt string, key string) (int, error) {
	if response, exists := m.responses[key]; exists {
		if i, ok := response.(int); ok {
			return i, nil
		}
	}
	return 0, nil
}

func (m *MockUIInput) ConfirmInput(prompt string, key string) (bool, error) {
	return m.GetBoolInput(prompt, key)
}

// Mock interface definitions (since they don't exist in the actual interfaces package yet)
type MockPreviewData struct {
	Title       string
	Description string
	Structure   []MockPreviewItem
}

type MockPreviewItem struct {
	Name     string
	Type     string
	Children []MockPreviewItem
}

type MockProjectTree struct {
	Root *MockTreeNode
}

type MockTreeNode struct {
	Name     string
	Type     string
	Children []*MockTreeNode
}

type MockUIState struct {
	CurrentStep     string
	ProjectConfig   *models.ProjectConfig
	ValidationState *MockValidationState
}

type MockValidationState struct {
	Valid bool
}

type MockProgressUpdate struct {
	Stage   string
	Percent int
	Message string
}

type MockGenerationSummary struct {
	ProjectName    string
	OutputPath     string
	ComponentCount int
	FileCount      int
	Duration       time.Duration
}

type MockValidationError struct {
	Field   string
	Message string
	Code    string
}

func (e *MockValidationError) Error() string {
	return e.Message
}

type MockGenerationError struct {
	Stage   string
	Message string
	Code    string
	Cause   string
}

func (e *MockGenerationError) Error() string {
	return e.Message
}

type MockNetworkError struct {
	Operation string
	URL       string
	Message   string
	Code      string
}

func (e *MockNetworkError) Error() string {
	return e.Message
}

type MockNavigationResult struct {
	Action string
}

type MockWorkflowResult struct {
	Success bool
}

type MockErrorReport struct {
	Errors []MockUIErrorInterface
}

type MockUIErrorInterface interface {
	Error() string
}

type MockStorageBackend interface {
	Initialize() error
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Exists(key string) bool
	Close() error
}

// Mock implementations for UI testing

// Mock Config Components
type MockUIConfigManager struct{}

func NewMockUIConfigManager() *MockUIConfigManager {
	return &MockUIConfigManager{}
}

func (m *MockUIConfigManager) CollectConfiguration(ctx context.Context, input *MockUIInput) (*models.ProjectConfig, error) {
	return &models.ProjectConfig{
		Name:         input.responses["project_name"].(string),
		Organization: input.responses["organization"].(string),
		Description:  input.responses["description"].(string),
		License:      input.responses["license"].(string),
	}, nil
}

func (m *MockUIConfigManager) ValidateConfiguration(config *models.ProjectConfig) error {
	if config.Name == "" {
		return NewMockUIError("project name is required")
	}
	return nil
}

func (m *MockUIConfigManager) ExportConfiguration(config *models.ProjectConfig, path string) error {
	content := "name: " + config.Name + "\norganization: " + config.Organization
	return os.WriteFile(path, []byte(content), 0644)
}

func (m *MockUIConfigManager) ImportConfiguration(path string) (*models.ProjectConfig, error) {
	return &models.ProjectConfig{
		Name:         "imported-project",
		Organization: "imported-org",
	}, nil
}

func (m *MockUIConfigManager) UpdateConfigurationFromPreview(config *models.ProjectConfig, preview *MockPreviewData, input *MockUIInput) (*models.ProjectConfig, error) {
	// Mock implementation that adds a component based on input
	if addComponent, exists := input.responses["add_component"]; exists && addComponent == "frontend" {
		config.Components.Frontend.NextJS.App = true
	}
	return config, nil
}

type MockProjectInfoCollector struct{}

func NewMockProjectInfoCollector() *MockProjectInfoCollector {
	return &MockProjectInfoCollector{}
}

func (m *MockProjectInfoCollector) CollectProjectInfo(ctx context.Context, input *MockUIInput) (*models.ProjectConfig, error) {
	return &models.ProjectConfig{
		Name:         input.responses["project_name"].(string),
		Organization: input.responses["organization"].(string),
		Description:  input.responses["description"].(string),
		License:      input.responses["license"].(string),
	}, nil
}

type MockComponentCollector struct{}

func NewMockComponentCollector() *MockComponentCollector {
	return &MockComponentCollector{}
}

func (m *MockComponentCollector) CollectComponents(ctx context.Context, input *MockUIInput) (*models.Components, error) {
	return &models.Components{
		Backend: models.BackendComponents{
			GoGin: true,
		},
	}, nil
}

type MockAdvancedOptionsCollector struct{}

func NewMockAdvancedOptionsCollector() *MockAdvancedOptionsCollector {
	return &MockAdvancedOptionsCollector{}
}

func (m *MockAdvancedOptionsCollector) CollectAdvancedOptions(ctx context.Context, input *MockUIInput) (*interfaces.AdvancedOptions, error) {
	return &interfaces.AdvancedOptions{}, nil
}

type MockProjectValidator struct{}

func NewMockProjectValidator() *MockProjectValidator {
	return &MockProjectValidator{}
}

func (m *MockProjectValidator) ValidateProject(config *models.ProjectConfig) error {
	if config.Name == "" {
		return NewMockUIError("project name is required")
	}
	return nil
}

type MockComponentValidator struct{}

func NewMockComponentValidator() *MockComponentValidator {
	return &MockComponentValidator{}
}

func (m *MockComponentValidator) ValidateComponents(components *models.Components) error {
	// Since we only have GoGin now, no conflicts to check
	return nil
	return nil
}

type MockSummaryFormatter struct{}

func NewMockSummaryFormatter() *MockSummaryFormatter {
	return &MockSummaryFormatter{}
}

func (m *MockSummaryFormatter) FormatSummary(config *models.ProjectConfig) (string, error) {
	return "Project: " + config.Name + "\nOrganization: " + config.Organization, nil
}

type MockExportFormatter struct{}

func NewMockExportFormatter() *MockExportFormatter {
	return &MockExportFormatter{}
}

func (m *MockExportFormatter) FormatYAML(config *models.ProjectConfig) ([]byte, error) {
	content := "name: " + config.Name + "\norganization: " + config.Organization
	return []byte(content), nil
}

func (m *MockExportFormatter) FormatJSON(config *models.ProjectConfig) ([]byte, error) {
	content := `{"name": "` + config.Name + `", "organization": "` + config.Organization + `"}`
	return []byte(content), nil
}

// Mock Preview Components
type MockPreviewManager struct{}

func NewMockPreviewManager() *MockPreviewManager {
	return &MockPreviewManager{}
}

func (m *MockPreviewManager) GenerateStructurePreview(config *models.ProjectConfig) (*MockPreviewData, error) {
	return &MockPreviewData{
		Title:       "Project Structure Preview",
		Description: "Preview for " + config.Name,
		Structure: []MockPreviewItem{
			{
				Name: "src/",
				Type: "directory",
			},
		},
	}, nil
}

func (m *MockPreviewManager) GenerateComponentPreview(config *models.ProjectConfig) (*MockPreviewData, error) {
	return &MockPreviewData{
		Title:       "Component Preview",
		Description: "Components for " + config.Name,
	}, nil
}

func (m *MockPreviewManager) DisplayPreview(preview *MockPreviewData) error {
	return nil
}

type MockTreeBuilder struct{}

func NewMockTreeBuilder() *MockTreeBuilder {
	return &MockTreeBuilder{}
}

func (m *MockTreeBuilder) BuildProjectTree(config *models.ProjectConfig) (*interfaces.ProjectTree, error) {
	return &interfaces.ProjectTree{
		Root: &interfaces.TreeNode{
			Label: config.Name,
			Metadata: map[string]interface{}{
				"type": "directory",
			},
		},
	}, nil
}

type MockTreeRenderer struct{}

func NewMockTreeRenderer() *MockTreeRenderer {
	return &MockTreeRenderer{}
}

func (m *MockTreeRenderer) RenderTree(tree *interfaces.ProjectTree) (string, error) {
	return "├── " + tree.Root.Label + "/", nil
}

type MockTreeFormatter struct{}

func NewMockTreeFormatter() *MockTreeFormatter {
	return &MockTreeFormatter{}
}

func (m *MockTreeFormatter) FormatTree(tree *interfaces.ProjectTree, format string) (string, error) {
	if format == "unicode" {
		return "└── " + tree.Root.Label + "/", nil
	}
	return "├── " + tree.Root.Label + "/", nil
}

type MockFrontendPreview struct{}

func NewMockFrontendPreview() *MockFrontendPreview {
	return &MockFrontendPreview{}
}

func (m *MockFrontendPreview) GeneratePreview(config *models.FrontendComponents) (*interfaces.PreviewData, error) {
	return &interfaces.PreviewData{
		Content: "Frontend Preview Content",
		Type:    "frontend_preview",
		Metadata: map[string]interface{}{
			"title": "Frontend Preview",
		},
	}, nil
}

type MockBackendPreview struct{}

func NewMockBackendPreview() *MockBackendPreview {
	return &MockBackendPreview{}
}

func (m *MockBackendPreview) GeneratePreview(config *models.BackendComponents) (*interfaces.PreviewData, error) {
	return &interfaces.PreviewData{
		Content: "Backend Preview Content",
		Type:    "backend_preview",
		Metadata: map[string]interface{}{
			"title": "Backend Preview",
		},
	}, nil
}

type MockInfrastructurePreview struct{}

func NewMockInfrastructurePreview() *MockInfrastructurePreview {
	return &MockInfrastructurePreview{}
}

func (m *MockInfrastructurePreview) GeneratePreview(config *models.InfrastructureComponents) (*interfaces.PreviewData, error) {
	return &interfaces.PreviewData{
		Content: "Infrastructure Preview Content",
		Type:    "infrastructure_preview",
		Metadata: map[string]interface{}{
			"title": "Infrastructure Preview",
		},
	}, nil
}

type MockConsoleDisplay struct{}

func NewMockConsoleDisplay() *MockConsoleDisplay {
	return &MockConsoleDisplay{}
}

func (m *MockConsoleDisplay) DisplayPreview(preview *interfaces.PreviewData) error {
	return nil
}

type MockInteractiveDisplay struct{}

func NewMockInteractiveDisplay() *MockInteractiveDisplay {
	return &MockInteractiveDisplay{}
}

func (m *MockInteractiveDisplay) DisplayInteractivePreview(ctx context.Context, preview *interfaces.PreviewData, input *MockUIInput) error {
	return nil
}

// Mock UI Manager
type MockUIManager struct{}

func NewMockUIManager() *MockUIManager {
	return &MockUIManager{}
}

func (m *MockUIManager) CollectProjectConfiguration(ctx context.Context, input *MockUIInput) (*models.ProjectConfig, error) {
	return &models.ProjectConfig{
		Name:         input.responses["project_name"].(string),
		Organization: input.responses["organization"].(string),
		Description:  input.responses["description"].(string),
		License:      input.responses["license"].(string),
	}, nil
}

func (m *MockUIManager) GenerateProjectPreview(config *models.ProjectConfig) (*interfaces.PreviewData, error) {
	return &interfaces.PreviewData{
		Content: "Project Preview Content",
		Type:    "project_preview",
		Metadata: map[string]interface{}{
			"title":       "Project Preview",
			"description": "Preview for " + config.Name,
		},
	}, nil
}

func (m *MockUIManager) DisplayPreview(preview *interfaces.PreviewData) error {
	return nil
}

func (m *MockUIManager) ConfirmConfiguration(ctx context.Context, config *models.ProjectConfig, input *MockUIInput) (bool, error) {
	if confirmed, exists := input.responses["confirm_config"]; exists {
		return confirmed.(bool), nil
	}
	return true, nil
}

func (m *MockUIManager) DisplayProgress(ctx context.Context, progressChan <-chan interfaces.ProgressUpdate) error {
	for range progressChan {
		// Mock progress display
	}
	return nil
}

func (m *MockUIManager) DisplayCompletionSummary(summary *interfaces.GenerationSummary) error {
	return nil
}

func (m *MockUIManager) ValidateConfiguration(config *models.ProjectConfig) error {
	if config.Name == "" {
		return NewMockUIError("project name is required")
	}
	return nil
}

func (m *MockUIManager) DisplayValidationError(err *interfaces.ValidationError) error {
	return nil
}

func (m *MockUIManager) HandleValidationError(ctx context.Context, err *interfaces.ValidationError, input *MockUIInput) (bool, error) {
	if retry, exists := input.responses["retry"]; exists {
		return retry.(bool), nil
	}
	return true, nil
}

func (m *MockUIManager) DisplayMainMenu(ctx context.Context, input *MockUIInput) (string, error) {
	if choice, exists := input.responses["menu_choice"]; exists {
		return choice.(string), nil
	}
	return "project_setup", nil
}

func (m *MockUIManager) HandleNavigation(ctx context.Context, input *MockUIInput) (*interfaces.NavigationResult, error) {
	return &interfaces.NavigationResult{
		Action: "next",
	}, nil
}

func (m *MockUIManager) DisplayBreadcrumbs(breadcrumbs []string) error {
	return nil
}

func (m *MockUIManager) RegisterConfigManager(manager interface{}) error {
	return nil
}

func (m *MockUIManager) RegisterPreviewManager(manager interface{}) error {
	return nil
}

func (m *MockUIManager) ExecuteCoordinatedWorkflow(ctx context.Context, workflowType string, input *MockUIInput) (interface{}, error) {
	return &MockWorkflowResult{
		Success: true,
	}, nil
}

func (m *MockUIManager) DisplayValidationErrors(errors []interfaces.ValidationError) error {
	return nil
}

func (m *MockUIManager) HandleValidationErrors(ctx context.Context, errors []interfaces.ValidationError, input *MockUIInput) (bool, error) {
	return true, nil
}

func (m *MockUIManager) DisplayGenerationError(err *interfaces.GenerationError) error {
	return nil
}

func (m *MockUIManager) HandleGenerationError(ctx context.Context, err *interfaces.GenerationError, input *MockUIInput) (bool, error) {
	return true, nil
}

func (m *MockUIManager) DisplayNetworkError(err *interfaces.NetworkError) error {
	return nil
}

func (m *MockUIManager) HandleNetworkError(ctx context.Context, err *interfaces.NetworkError, input *MockUIInput) (bool, error) {
	return true, nil
}

func (m *MockUIManager) AggregateErrors(errors []interfaces.UIError) (*interfaces.ErrorReport, error) {
	return &interfaces.ErrorReport{
		ID:      "test-error-report",
		Type:    "AGGREGATED_ERRORS",
		Message: fmt.Sprintf("Aggregated %d errors", len(errors)),
		Metadata: map[string]interface{}{
			"error_count": len(errors),
			"errors":      errors,
		},
	}, nil
}

func (m *MockUIManager) DisplayErrorReport(report *interfaces.ErrorReport) error {
	return nil
}

func (m *MockUIManager) ExportErrorReport(report *interfaces.ErrorReport, path string) error {
	content := "Error Report\n"
	return os.WriteFile(path, []byte(content), 0644)
}

type MockStateManager struct{}

func NewMockStateManager() *MockStateManager {
	return &MockStateManager{}
}

type MockManagerWithState struct {
	stateManager *MockStateManager
	state        *interfaces.UIState
}

func NewMockManagerWithState(stateManager *MockStateManager) *MockManagerWithState {
	return &MockManagerWithState{
		stateManager: stateManager,
		state:        &interfaces.UIState{},
	}
}

func (m *MockManagerWithState) UpdateState(state *interfaces.UIState) error {
	m.state = state
	return nil
}

func (m *MockManagerWithState) GetState() (*interfaces.UIState, error) {
	return m.state, nil
}

// Mock Error Types
type MockUIError struct {
	message string
}

func NewMockUIError(message string) *MockUIError {
	return &MockUIError{message: message}
}

func (e *MockUIError) Error() string {
	return e.message
}
