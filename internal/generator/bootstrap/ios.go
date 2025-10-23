package bootstrap

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// iOSExecutor handles iOS project generation
type iOSExecutor struct {
	*BaseExecutor
	toolDiscovery ToolDiscovery
}

// NewiOSExecutor creates a new iOS executor
func NewiOSExecutor(toolDiscovery ToolDiscovery) *iOSExecutor {
	return &iOSExecutor{
		BaseExecutor:  NewBaseExecutor("swift"),
		toolDiscovery: toolDiscovery,
	}
}

// Execute generates an iOS project if tools are available
func (ie *iOSExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error) {
	// Check if running on macOS
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("iOS development requires macOS (fallback required)")
	}

	// Check for available iOS tools
	tool, err := ie.detectiOSTool()
	if err != nil {
		return nil, fmt.Errorf("no iOS tools available: %w (fallback required)", err)
	}

	// Attempt to generate project based on available tool
	switch tool {
	case "swift":
		return ie.generateWithSwiftPackage(ctx, spec)
	case "xcodebuild":
		return ie.generateWithXcodebuild(ctx, spec)
	default:
		return nil, fmt.Errorf("unsupported iOS tool: %s (fallback required)", tool)
	}
}

// SupportsComponent checks if this executor supports the given component type
func (ie *iOSExecutor) SupportsComponent(componentType string) bool {
	return componentType == "ios" || componentType == "mobile-ios"
}

// GetDefaultFlags returns default flags for iOS generation
func (ie *iOSExecutor) GetDefaultFlags(componentType string) []string {
	if !ie.SupportsComponent(componentType) {
		return []string{}
	}

	return []string{"package", "init", "--type", "executable"}
}

// ValidateConfig validates component-specific configuration
func (ie *iOSExecutor) ValidateConfig(config map[string]interface{}) error {
	// Validate name
	if name, ok := config["name"].(string); !ok || name == "" {
		return fmt.Errorf("name is required and must be a string")
	}

	// Validate bundle_id if provided
	if bundleID, exists := config["bundle_id"]; exists {
		if bundleStr, ok := bundleID.(string); ok {
			// Basic validation for bundle identifier format
			if len(bundleStr) == 0 {
				return fmt.Errorf("bundle_id cannot be empty")
			}
			// Bundle ID should contain at least one dot
			hasDot := false
			for _, c := range bundleStr {
				if c == '.' {
					hasDot = true
					break
				}
			}
			if !hasDot {
				return fmt.Errorf("bundle_id must be a valid bundle identifier (e.g., com.example.app)")
			}
		} else {
			return fmt.Errorf("bundle_id must be a string")
		}
	}

	// Validate deployment_target if provided
	if target, exists := config["deployment_target"]; exists {
		if targetStr, ok := target.(string); ok {
			// Basic validation for iOS version format (e.g., "14.0", "15.0")
			if len(targetStr) == 0 {
				return fmt.Errorf("deployment_target cannot be empty")
			}
			// Should contain at least one dot
			hasDot := false
			for _, c := range targetStr {
				if c == '.' {
					hasDot = true
					break
				}
			}
			if !hasDot {
				return fmt.Errorf("deployment_target must be a valid iOS version (e.g., 14.0)")
			}
		} else {
			return fmt.Errorf("deployment_target must be a string")
		}
	}

	// Validate language if provided
	if lang, exists := config["language"]; exists {
		if langStr, ok := lang.(string); ok {
			validLangs := map[string]bool{"swift": true, "objective-c": true}
			if !validLangs[langStr] {
				return fmt.Errorf("language must be either 'swift' or 'objective-c'")
			}
		} else {
			return fmt.Errorf("language must be a string")
		}
	}

	return nil
}

// detectiOSTool detects which iOS tool is available
func (ie *iOSExecutor) detectiOSTool() (string, error) {
	// iOS development only works on macOS
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("iOS development requires macOS")
	}

	// Check for swift package manager
	if ie.toolDiscovery != nil {
		if available, _ := ie.toolDiscovery.IsAvailable("swift"); available {
			return "swift", nil
		}

		// Check for xcodebuild
		if available, _ := ie.toolDiscovery.IsAvailable("xcodebuild"); available {
			return "xcodebuild", nil
		}
	} else {
		// Fallback to direct check if no tool discovery provided
		if _, err := exec.LookPath("swift"); err == nil {
			return "swift", nil
		}

		if _, err := exec.LookPath("xcodebuild"); err == nil {
			return "xcodebuild", nil
		}
	}

	return "", fmt.Errorf("neither swift nor xcodebuild found")
}

// generateWithSwiftPackage generates an iOS project using Swift Package Manager
func (ie *iOSExecutor) generateWithSwiftPackage(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error) {
	// Get project configuration
	_, ok := spec.Config["name"].(string)
	if !ok {
		return nil, fmt.Errorf("project name is required in config")
	}

	// Swift Package Manager is primarily for libraries/frameworks
	// For iOS apps, we need Xcode project structure
	// This will trigger fallback generation
	return nil, fmt.Errorf("Swift Package Manager cannot create iOS app projects (fallback required)")
}

// generateWithXcodebuild generates an iOS project using xcodebuild
func (ie *iOSExecutor) generateWithXcodebuild(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error) {
	// Get project configuration
	projectName, ok := spec.Config["name"].(string)
	if !ok || projectName == "" {
		return nil, fmt.Errorf("project name is required in config")
	}

	// xcodebuild doesn't have a built-in project generator
	// It's used for building existing projects, not creating new ones
	// Project creation is typically done through Xcode GUI or templates
	// This will trigger fallback generation
	return nil, fmt.Errorf("xcodebuild cannot create new projects (fallback required)")
}

// GetManualSteps returns manual steps required after iOS generation
func (ie *iOSExecutor) GetManualSteps(spec *BootstrapSpec) []string {
	steps := []string{
		"Install Xcode from the Mac App Store",
		"Open Xcode and accept the license agreement",
		"Install Xcode Command Line Tools: xcode-select --install",
		"Open the .xcodeproj file in Xcode",
		"Select a development team in the project settings",
		"Choose a simulator or connect an iOS device",
		"Click 'Run' to build and deploy the app",
	}

	// Add SwiftUI-specific steps if configured
	if useSwiftUI, ok := spec.Config["swiftui"].(bool); ok && useSwiftUI {
		steps = append(steps, "The project uses SwiftUI for the user interface")
	}

	return steps
}

// GetInstallInstructions returns installation instructions for iOS tools
func (ie *iOSExecutor) GetInstallInstructions(os string) string {
	if os != "darwin" {
		return "iOS development requires macOS. Install Xcode from the Mac App Store."
	}

	return `Install Xcode from the Mac App Store:
1. Open the App Store application
2. Search for "Xcode"
3. Click "Get" or "Install"
4. After installation, open Xcode and accept the license agreement
5. Install Command Line Tools: xcode-select --install

Alternative: Install Xcode Command Line Tools only (for Swift Package Manager):
  xcode-select --install
`
}

// RequiresFallback checks if fallback generation is required
func (ie *iOSExecutor) RequiresFallback() bool {
	// iOS always requires fallback since there's no CLI tool for project creation
	// or if not running on macOS
	if runtime.GOOS != "darwin" {
		return true
	}

	tool, err := ie.detectiOSTool()
	// Even if tools are available, we need fallback for project creation
	return err != nil || tool == "" || true
}

// CanGenerateWithTools checks if tools are available (even if fallback is needed)
func (ie *iOSExecutor) CanGenerateWithTools() bool {
	if runtime.GOOS != "darwin" {
		return false
	}

	tool, err := ie.detectiOSTool()
	return err == nil && tool != ""
}
