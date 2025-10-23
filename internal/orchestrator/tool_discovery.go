package orchestrator

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/versions"
)

// ToolDiscovery implements tool discovery and validation functionality
type ToolDiscovery struct {
	registry  *models.ToolRegistry
	cache     *ToolCache
	logger    *logger.Logger
	cacheTTL  time.Duration
	isOffline bool
}

// NewToolDiscovery creates a new tool discovery instance with registered tools
func NewToolDiscovery(log *logger.Logger) *ToolDiscovery {
	// Create persistent cache
	cache, err := NewToolCache(DefaultToolCacheConfig(), log)
	if err != nil && log != nil {
		log.Warn(fmt.Sprintf("Failed to create persistent cache: %v", err))
		// Continue without persistent cache
	}

	td := &ToolDiscovery{
		registry: &models.ToolRegistry{
			Tools: make(map[string]*models.ToolMetadata),
		},
		cache:     cache,
		logger:    log,
		cacheTTL:  5 * time.Minute, // Default cache TTL
		isOffline: false,
	}

	// Register known tools
	td.registerKnownTools()

	return td
}

// NewToolDiscoveryWithCache creates a new tool discovery instance with a custom cache
func NewToolDiscoveryWithCache(log *logger.Logger, cache *ToolCache) *ToolDiscovery {
	td := &ToolDiscovery{
		registry: &models.ToolRegistry{
			Tools: make(map[string]*models.ToolMetadata),
		},
		cache:     cache,
		logger:    log,
		cacheTTL:  5 * time.Minute,
		isOffline: false,
	}

	// Register known tools
	td.registerKnownTools()

	return td
}

// registerKnownTools registers all known bootstrap tools with their metadata
func (td *ToolDiscovery) registerKnownTools() {
	// Get versions from centralized config
	versionConfig, err := versions.Get()
	if err != nil {
		panic(fmt.Sprintf("failed to load version config: %v", err))
	}

	// Register npx for Next.js projects
	td.registry.Tools["npx"] = &models.ToolMetadata{
		Name:              "npx",
		Command:           "npx",
		VersionFlag:       "--version",
		MinVersion:        "", // No minimum version - npx is just a package runner
		FallbackAvailable: false,
		ComponentTypes:    []string{"nextjs"},
		InstallDocs: map[string]string{
			"linux":   "https://nodejs.org/en/download/package-manager",
			"darwin":  "https://nodejs.org/en/download/package-manager",
			"windows": "https://nodejs.org/en/download",
		},
	}

	// Register go for Go backend projects
	td.registry.Tools["go"] = &models.ToolMetadata{
		Name:              "go",
		Command:           "go",
		VersionFlag:       "version",
		MinVersion:        versionConfig.Backend.Go.Version,
		FallbackAvailable: false,
		ComponentTypes:    []string{"go-backend"},
		InstallDocs: map[string]string{
			"linux":   "https://go.dev/doc/install",
			"darwin":  "https://go.dev/doc/install",
			"windows": "https://go.dev/doc/install",
		},
	}

	// Register gradle for Android projects
	td.registry.Tools["gradle"] = &models.ToolMetadata{
		Name:              "gradle",
		Command:           "gradle",
		VersionFlag:       "--version",
		MinVersion:        versionConfig.Android.Gradle.Version,
		FallbackAvailable: true,
		ComponentTypes:    []string{"android"},
		InstallDocs: map[string]string{
			"linux":   "https://gradle.org/install/",
			"darwin":  "https://gradle.org/install/",
			"windows": "https://gradle.org/install/",
		},
	}

	// Register xcodebuild for iOS projects
	td.registry.Tools["xcodebuild"] = &models.ToolMetadata{
		Name:              "xcodebuild",
		Command:           "xcodebuild",
		VersionFlag:       "-version",
		MinVersion:        versionConfig.IOS.Xcode.Version,
		FallbackAvailable: true,
		ComponentTypes:    []string{"ios"},
		InstallDocs: map[string]string{
			"darwin": "https://developer.apple.com/xcode/",
		},
	}

	// Register docker for containerization
	td.registry.Tools["docker"] = &models.ToolMetadata{
		Name:              "docker",
		Command:           "docker",
		VersionFlag:       "--version",
		MinVersion:        "", // No minimum version requirement
		FallbackAvailable: false,
		ComponentTypes:    []string{"docker"},
		InstallDocs: map[string]string{
			"linux":   "https://docs.docker.com/engine/install/",
			"darwin":  "https://docs.docker.com/desktop/install/mac-install/",
			"windows": "https://docs.docker.com/desktop/install/windows-install/",
		},
	}

	// Register terraform for infrastructure
	td.registry.Tools["terraform"] = &models.ToolMetadata{
		Name:              "terraform",
		Command:           "terraform",
		VersionFlag:       "version",
		MinVersion:        versionConfig.Infrastructure.Terraform.Version,
		FallbackAvailable: false,
		ComponentTypes:    []string{"terraform"},
		InstallDocs: map[string]string{
			"linux":   "https://developer.hashicorp.com/terraform/install",
			"darwin":  "https://developer.hashicorp.com/terraform/install",
			"windows": "https://developer.hashicorp.com/terraform/install",
		},
	}
}

// Ensure ToolDiscovery implements the interface
var _ interfaces.ToolDiscoveryInterface = (*ToolDiscovery)(nil)

// IsAvailable checks if a tool is available in the system PATH
func (td *ToolDiscovery) IsAvailable(toolName string) (bool, error) {
	// Check cache first
	if td.cache != nil {
		// Use offline-aware cache retrieval
		var cached *models.CachedTool
		var found bool

		if td.isOffline {
			cached, found = td.cache.GetWithOfflineSupport(toolName)
		} else {
			cached, found = td.cache.Get(toolName)
		}

		if found {
			if td.logger != nil {
				td.logger.Debug(fmt.Sprintf("Tool '%s' availability from cache: %v", toolName, cached.Available))
			}
			return cached.Available, nil
		}
	}

	// In offline mode, if not in cache, assume unavailable
	if td.isOffline {
		if td.logger != nil {
			td.logger.Debug(fmt.Sprintf("Tool '%s' not in cache and offline mode active", toolName))
		}
		return false, fmt.Errorf("tool not in cache and offline mode active")
	}

	// Check if tool exists in PATH
	_, err := exec.LookPath(toolName)
	available := err == nil

	// Cache the result
	if td.cache != nil {
		td.cache.Set(toolName, available, "")
	}

	if td.logger != nil {
		if available {
			td.logger.Debug(fmt.Sprintf("Tool '%s' is available", toolName))
		} else {
			td.logger.Debug(fmt.Sprintf("Tool '%s' is not available: %v", toolName, err))
		}
	}

	return available, nil
}

// GetVersion retrieves the installed version of a tool
func (td *ToolDiscovery) GetVersion(toolName string) (string, error) {
	// Check cache first
	if td.cache != nil {
		var cached *models.CachedTool
		var found bool

		if td.isOffline {
			cached, found = td.cache.GetWithOfflineSupport(toolName)
		} else {
			cached, found = td.cache.Get(toolName)
		}

		if found && cached.Version != "" {
			if td.logger != nil {
				td.logger.Debug(fmt.Sprintf("Tool '%s' version from cache: %s", toolName, cached.Version))
			}
			return cached.Version, nil
		}
	}

	// In offline mode, if not in cache, return error
	if td.isOffline {
		return "", fmt.Errorf("tool version not in cache and offline mode active")
	}

	// Get tool metadata
	metadata, exists := td.registry.Tools[toolName]
	if !exists {
		return "", fmt.Errorf("tool '%s' not registered", toolName)
	}

	// Check if tool is available
	available, err := td.IsAvailable(toolName)
	if err != nil || !available {
		return "", fmt.Errorf("tool '%s' is not available", toolName)
	}

	// Execute version command
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if metadata.VersionFlag != "" {
		cmd = exec.CommandContext(ctx, metadata.Command, metadata.VersionFlag)
	} else {
		cmd = exec.CommandContext(ctx, metadata.Command, "version")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get version for '%s': %w", toolName, err)
	}

	version := strings.TrimSpace(string(output))

	// Cache the result
	if td.cache != nil {
		td.cache.Set(toolName, true, version)
	}

	if td.logger != nil {
		td.logger.Debug(fmt.Sprintf("Tool '%s' version: %s", toolName, version))
	}

	return version, nil
}

// CheckRequirements validates that all required tools are available
func (td *ToolDiscovery) CheckRequirements(tools []string) (interface{}, error) {
	result := &models.ToolCheckResult{
		AllAvailable: true,
		Tools:        make(map[string]*models.Tool),
		Missing:      []string{},
		Outdated:     []string{},
		CheckedAt:    time.Now(),
	}

	for _, toolName := range tools {
		metadata, exists := td.registry.Tools[toolName]
		if !exists {
			result.AllAvailable = false
			result.Missing = append(result.Missing, toolName)
			continue
		}

		tool := &models.Tool{
			Name:                toolName,
			Command:             metadata.Command,
			VersionCommand:      metadata.VersionFlag,
			MinVersion:          metadata.MinVersion,
			InstallInstructions: metadata.InstallDocs,
		}

		// Check availability
		available, _ := td.IsAvailable(toolName)
		tool.Available = available
		tool.LastChecked = time.Now()

		if !available {
			result.AllAvailable = false
			result.Missing = append(result.Missing, toolName)
		} else {
			// Get version if available
			version, err := td.GetVersion(toolName)
			if err == nil {
				tool.InstalledVersion = version
				// TODO: Add version comparison logic if needed
			}
		}

		result.Tools[toolName] = tool
	}

	if td.logger != nil {
		td.logger.Info(fmt.Sprintf("Tool check complete: %d/%d available",
			len(tools)-len(result.Missing), len(tools)))
	}

	return result, nil
}

// ClearCache clears the tool availability cache
func (td *ToolDiscovery) ClearCache() {
	if td.cache != nil {
		td.cache.Clear()
	}

	if td.logger != nil {
		td.logger.Debug("Tool cache cleared")
	}
}

// SetCacheTTL sets the cache time-to-live duration
func (td *ToolDiscovery) SetCacheTTL(ttl time.Duration) {
	td.cacheTTL = ttl
	if td.cache != nil {
		td.cache.SetTTL(ttl)
	}
}

// SaveCache persists the cache to disk
func (td *ToolDiscovery) SaveCache() error {
	if td.cache == nil {
		return fmt.Errorf("cache not initialized")
	}
	return td.cache.Save()
}

// GetCacheStats returns statistics about the cache
func (td *ToolDiscovery) GetCacheStats() map[string]interface{} {
	if td.cache == nil {
		return map[string]interface{}{
			"enabled": false,
		}
	}
	stats := td.cache.GetStats()
	stats["enabled"] = true
	return stats
}

// GetInstallInstructions returns OS-specific installation instructions for a tool
func (td *ToolDiscovery) GetInstallInstructions(toolName string, os string) string {
	// Normalize OS name
	normalizedOS := normalizeOS(os)

	// Get tool metadata
	metadata, exists := td.registry.Tools[toolName]
	if !exists {
		return fmt.Sprintf("Tool '%s' is not registered. Please check the tool name.", toolName)
	}

	// Get OS-specific instructions
	instructions, exists := metadata.InstallDocs[normalizedOS]
	if !exists {
		// Try to provide generic instructions
		if normalizedOS == "darwin" && toolName == "xcodebuild" {
			return "xcodebuild is only available on macOS. Install Xcode from the App Store."
		}
		return fmt.Sprintf("Installation instructions for '%s' on '%s' are not available. Please visit the official documentation.", toolName, os)
	}

	// Build formatted instructions
	return formatInstallInstructions(toolName, normalizedOS, instructions, metadata.FallbackAvailable)
}

// normalizeOS normalizes OS names to standard values
func normalizeOS(os string) string {
	os = strings.ToLower(strings.TrimSpace(os))

	switch os {
	case "darwin", "macos", "osx", "mac":
		return "darwin"
	case "linux", "unix":
		return "linux"
	case "windows", "win":
		return "windows"
	default:
		// Use runtime.GOOS if os is empty or unknown
		if os == "" {
			return runtime.GOOS
		}
		return os
	}
}

// formatInstallInstructions formats installation instructions for display
func formatInstallInstructions(toolName, os, docURL string, hasFallback bool) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Installation instructions for '%s' on %s:\n", toolName, os))
	builder.WriteString(fmt.Sprintf("  Documentation: %s\n", docURL))

	// Add OS-specific quick install commands
	switch os {
	case "darwin":
		switch toolName {
		case "npx":
			builder.WriteString("  Quick install: brew install node\n")
		case "go":
			builder.WriteString("  Quick install: brew install go\n")
		case "gradle":
			builder.WriteString("  Quick install: brew install gradle\n")
		case "docker":
			builder.WriteString("  Quick install: brew install --cask docker\n")
		case "terraform":
			builder.WriteString("  Quick install: brew install terraform\n")
		}
	case "linux":
		switch toolName {
		case "npx":
			builder.WriteString("  Quick install: sudo apt-get install nodejs npm (Debian/Ubuntu)\n")
			builder.WriteString("                 sudo yum install nodejs npm (RHEL/CentOS)\n")
		case "go":
			builder.WriteString("  Quick install: sudo apt-get install golang (Debian/Ubuntu)\n")
			builder.WriteString("                 sudo yum install golang (RHEL/CentOS)\n")
		case "gradle":
			builder.WriteString("  Quick install: sudo apt-get install gradle (Debian/Ubuntu)\n")
		case "docker":
			builder.WriteString("  Quick install: curl -fsSL https://get.docker.com | sh\n")
		case "terraform":
			builder.WriteString("  Quick install: sudo apt-get install terraform (with HashiCorp repo)\n")
		}
	case "windows":
		switch toolName {
		case "npx":
			builder.WriteString("  Quick install: Download from nodejs.org or use 'choco install nodejs'\n")
		case "go":
			builder.WriteString("  Quick install: Download from go.dev or use 'choco install golang'\n")
		case "gradle":
			builder.WriteString("  Quick install: choco install gradle\n")
		case "docker":
			builder.WriteString("  Quick install: Download Docker Desktop from docker.com\n")
		case "terraform":
			builder.WriteString("  Quick install: choco install terraform\n")
		}
	}

	if hasFallback {
		builder.WriteString("\n  Note: Fallback generation is available if this tool cannot be installed.\n")
	}

	return builder.String()
}

// GetToolMetadata returns metadata for a specific tool
func (td *ToolDiscovery) GetToolMetadata(toolName string) (*models.ToolMetadata, error) {
	metadata, exists := td.registry.Tools[toolName]
	if !exists {
		return nil, fmt.Errorf("tool '%s' not registered", toolName)
	}
	return metadata, nil
}

// ListRegisteredTools returns a list of all registered tool names
func (td *ToolDiscovery) ListRegisteredTools() []string {
	tools := make([]string, 0, len(td.registry.Tools))
	for name := range td.registry.Tools {
		tools = append(tools, name)
	}
	return tools
}

// GetToolsForComponent returns tools required for a specific component type
func (td *ToolDiscovery) GetToolsForComponent(componentType string) []string {
	var tools []string
	for name, metadata := range td.registry.Tools {
		for _, ct := range metadata.ComponentTypes {
			if ct == componentType {
				tools = append(tools, name)
				break
			}
		}
	}
	return tools
}

// HasFallback checks if a component type has fallback generation available
func (td *ToolDiscovery) HasFallback(componentType string) bool {
	for _, metadata := range td.registry.Tools {
		for _, ct := range metadata.ComponentTypes {
			if ct == componentType {
				return metadata.FallbackAvailable
			}
		}
	}
	return false
}
