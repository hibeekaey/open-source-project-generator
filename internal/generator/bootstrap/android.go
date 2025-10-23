package bootstrap

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// AndroidExecutor handles Android project generation
type AndroidExecutor struct {
	*BaseExecutor
	toolDiscovery ToolDiscovery
}

// ToolDiscovery interface for checking tool availability
type ToolDiscovery interface {
	IsAvailable(toolName string) (bool, error)
	GetVersion(toolName string) (string, error)
}

// NewAndroidExecutor creates a new Android executor
func NewAndroidExecutor(toolDiscovery ToolDiscovery) *AndroidExecutor {
	return &AndroidExecutor{
		BaseExecutor:  NewBaseExecutor("gradle"),
		toolDiscovery: toolDiscovery,
	}
}

// Execute generates an Android project if tools are available
func (ae *AndroidExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error) {
	// Check for available Android tools
	tool, err := ae.detectAndroidTool()
	if err != nil {
		return nil, fmt.Errorf("no Android tools available: %w (fallback required)", err)
	}

	// Attempt to generate project based on available tool
	switch tool {
	case "gradle":
		return ae.generateWithGradle(ctx, spec)
	case "android":
		return ae.generateWithAndroidCLI(ctx, spec)
	default:
		return nil, fmt.Errorf("unsupported Android tool: %s (fallback required)", tool)
	}
}

// SupportsComponent checks if this executor supports the given component type
func (ae *AndroidExecutor) SupportsComponent(componentType string) bool {
	return componentType == "android" || componentType == "mobile-android"
}

// GetDefaultFlags returns default flags for Android generation
func (ae *AndroidExecutor) GetDefaultFlags(componentType string) []string {
	if !ae.SupportsComponent(componentType) {
		return []string{}
	}

	return []string{"init", "--type", "basic"}
}

// ValidateConfig validates component-specific configuration
func (ae *AndroidExecutor) ValidateConfig(config map[string]interface{}) error {
	// Validate name
	if name, ok := config["name"].(string); !ok || name == "" {
		return fmt.Errorf("name is required and must be a string")
	}

	// Validate package name if provided
	if pkg, exists := config["package"]; exists {
		if pkgStr, ok := pkg.(string); ok {
			// Basic validation for Java package name format
			if len(pkgStr) == 0 {
				return fmt.Errorf("package name cannot be empty")
			}
			// Package should contain at least one dot
			hasDot := false
			for _, c := range pkgStr {
				if c == '.' {
					hasDot = true
					break
				}
			}
			if !hasDot {
				return fmt.Errorf("package name must be a valid Java package (e.g., com.example.app)")
			}
		} else {
			return fmt.Errorf("package must be a string")
		}
	}

	// Validate SDK levels if provided
	if minSDK, exists := config["min_sdk"]; exists {
		switch v := minSDK.(type) {
		case int:
			if v < 21 || v > 36 {
				return fmt.Errorf("min_sdk must be between 21 and 36")
			}
		case float64:
			if v < 21 || v > 34 {
				return fmt.Errorf("min_sdk must be between 21 and 34")
			}
		default:
			return fmt.Errorf("min_sdk must be a number")
		}
	}

	if targetSDK, exists := config["target_sdk"]; exists {
		switch v := targetSDK.(type) {
		case int:
			if v < 21 || v > 36 {
				return fmt.Errorf("target_sdk must be between 21 and 36")
			}
		case float64:
			if v < 21 || v > 34 {
				return fmt.Errorf("target_sdk must be between 21 and 34")
			}
		default:
			return fmt.Errorf("target_sdk must be a number")
		}
	}

	// Validate language if provided
	if lang, exists := config["language"]; exists {
		if langStr, ok := lang.(string); ok {
			validLangs := map[string]bool{"kotlin": true, "java": true}
			if !validLangs[langStr] {
				return fmt.Errorf("language must be either 'kotlin' or 'java'")
			}
		} else {
			return fmt.Errorf("language must be a string")
		}
	}

	return nil
}

// detectAndroidTool detects which Android tool is available
func (ae *AndroidExecutor) detectAndroidTool() (string, error) {
	// Check for gradle
	if ae.toolDiscovery != nil {
		if available, _ := ae.toolDiscovery.IsAvailable("gradle"); available {
			return "gradle", nil
		}

		// Check for android command (from Android SDK)
		if available, _ := ae.toolDiscovery.IsAvailable("android"); available {
			return "android", nil
		}
	} else {
		// Fallback to direct check if no tool discovery provided
		if _, err := exec.LookPath("gradle"); err == nil {
			return "gradle", nil
		}

		if _, err := exec.LookPath("android"); err == nil {
			return "android", nil
		}
	}

	return "", fmt.Errorf("neither gradle nor android CLI tools found")
}

// generateWithGradle generates an Android project using Gradle
func (ae *AndroidExecutor) generateWithGradle(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error) {
	// Get project configuration
	projectName, ok := spec.Config["name"].(string)
	if !ok || projectName == "" {
		return nil, fmt.Errorf("project name is required in config")
	}

	packageName, ok := spec.Config["package"].(string)
	if !ok || packageName == "" {
		packageName = fmt.Sprintf("com.example.%s", projectName)
	}

	// Note: Gradle doesn't have a built-in project generator like create-next-app
	// This would typically require Android Studio or a custom Gradle init script
	// For now, we return an error to trigger fallback generation
	return nil, fmt.Errorf("gradle-based generation requires Android Studio or custom init script (fallback required)")
}

// generateWithAndroidCLI generates an Android project using Android CLI
func (ae *AndroidExecutor) generateWithAndroidCLI(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error) {
	// Get project configuration
	projectName, ok := spec.Config["name"].(string)
	if !ok || projectName == "" {
		return nil, fmt.Errorf("project name is required in config")
	}

	packageName, ok := spec.Config["package"].(string)
	if !ok || packageName == "" {
		packageName = fmt.Sprintf("com.example.%s", projectName)
	}

	// Build android create project command
	args := []string{
		"create", "project",
		"--name", projectName,
		"--package", packageName,
		"--path", spec.TargetDir,
	}

	// Add activity name if specified
	if activity, ok := spec.Config["activity"].(string); ok && activity != "" {
		args = append(args, "--activity", activity)
	} else {
		args = append(args, "--activity", "MainActivity")
	}

	// Add target SDK if specified
	if target, ok := spec.Config["target"].(string); ok && target != "" {
		args = append(args, "--target", target)
	}

	execSpec := &BootstrapSpec{
		ComponentType: spec.ComponentType,
		TargetDir:     spec.TargetDir,
		Config:        spec.Config,
		Flags:         args,
		Timeout:       spec.Timeout,
	}

	// Update tool name for execution
	ae.toolName = "android"

	result, err := ae.BaseExecutor.Execute(ctx, execSpec)
	if err != nil {
		return result, fmt.Errorf("Android CLI generation failed: %w (fallback may be required)", err)
	}

	return result, nil
}

// GetManualSteps returns manual steps required after Android generation
func (ae *AndroidExecutor) GetManualSteps(spec *BootstrapSpec) []string {
	return []string{
		"Install Android Studio from https://developer.android.com/studio",
		"Open the project in Android Studio",
		"Wait for Gradle sync to complete",
		"Configure Android SDK if prompted",
		"Connect an Android device or start an emulator",
		"Click 'Run' to build and deploy the app",
	}
}

// GetInstallInstructions returns installation instructions for Android tools
func (ae *AndroidExecutor) GetInstallInstructions(os string) string {
	instructions := map[string]string{
		"darwin": "Install Android Studio from https://developer.android.com/studio\n" +
			"Or install via Homebrew: brew install --cask android-studio",
		"linux": "Install Android Studio from https://developer.android.com/studio\n" +
			"Or use your package manager:\n" +
			"  Ubuntu/Debian: sudo snap install android-studio --classic\n" +
			"  Arch: yay -S android-studio",
		"windows": "Download and install Android Studio from https://developer.android.com/studio\n" +
			"Or use Chocolatey: choco install androidstudio",
	}

	if instruction, ok := instructions[os]; ok {
		return instruction
	}

	return "Install Android Studio from https://developer.android.com/studio"
}

// RequiresFallback checks if fallback generation is required
func (ae *AndroidExecutor) RequiresFallback() bool {
	tool, err := ae.detectAndroidTool()
	return err != nil || tool == ""
}
