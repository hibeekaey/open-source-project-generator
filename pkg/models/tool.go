package models

import (
	"time"
)

// Tool represents metadata about a bootstrap tool
type Tool struct {
	Name                string            // Tool name (e.g., "npx", "go")
	Command             string            // Command to execute
	VersionCommand      string            // Command to get version
	MinVersion          string            // Minimum required version
	Available           bool              // Whether tool is available
	InstalledVersion    string            // Currently installed version
	InstallInstructions map[string]string // OS -> installation instructions
	LastChecked         time.Time         // Last availability check
}

// ToolMetadata defines static metadata for a tool
type ToolMetadata struct {
	Name              string            // Tool identifier
	Command           string            // Base command
	VersionFlag       string            // Flag to get version (e.g., "--version")
	MinVersion        string            // Minimum required version
	InstallDocs       map[string]string // OS -> documentation URL
	FallbackAvailable bool              // Whether fallback generation exists
	ComponentTypes    []string          // Component types this tool supports
}

// ToolCheckResult represents the result of checking tool availability
type ToolCheckResult struct {
	AllAvailable bool             // Whether all required tools are available
	Tools        map[string]*Tool // Tool name -> tool details
	Missing      []string         // List of missing tools
	Outdated     []string         // List of outdated tools
	CheckedAt    time.Time        // When the check was performed
}

// ToolRegistry manages registered bootstrap tools
type ToolRegistry struct {
	Tools map[string]*ToolMetadata // Tool name -> metadata
}

// CachedTool represents a cached tool availability check
type CachedTool struct {
	Available bool          // Whether tool is available
	Version   string        // Tool version
	CachedAt  time.Time     // Cache timestamp
	TTL       time.Duration // Time to live
}
