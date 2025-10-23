package models

import (
	"time"
)

// GenerationResult represents the complete result of project generation
type GenerationResult struct {
	Success            bool               // Overall success status
	ProjectRoot        string             // Root directory of generated project
	Components         []*ComponentResult // Results for each component
	Errors             []error            // Any errors encountered
	Warnings           []string           // Warning messages
	Duration           time.Duration      // Total generation time
	LogFile            string             // Path to detailed log file
	DryRun             bool               // Whether this was a dry run
	SecurityScanResult interface{}        // Security scan result (if performed)
}

// ComponentResult represents the result of generating a single component
type ComponentResult struct {
	Type        string        // Component type (e.g., "nextjs", "go-backend")
	Name        string        // Component name
	Success     bool          // Whether generation succeeded
	Method      string        // "bootstrap" or "fallback"
	ToolUsed    string        // Tool used for generation (if bootstrap)
	OutputPath  string        // Path where component was generated
	Error       error         // Error if generation failed
	ManualSteps []string      // Manual steps required after generation
	Duration    time.Duration // Time taken to generate
	Warnings    []string      // Component-specific warnings
}

// ExecutionResult represents the result of executing a bootstrap tool
type ExecutionResult struct {
	Success   bool          // Whether execution succeeded
	OutputDir string        // Directory where output was generated
	Stdout    string        // Standard output from tool
	Stderr    string        // Standard error from tool
	Duration  time.Duration // Execution time
	ToolUsed  string        // Tool that was executed
	ExitCode  int           // Exit code from tool
}

// PreviewResult represents the result of a dry-run preview
type PreviewResult struct {
	ProjectRoot string              // Where project would be generated
	Components  []*ComponentPreview // Preview for each component
	Structure   []string            // Directory structure that would be created
	Files       []string            // Files that would be created
	Warnings    []string            // Potential issues
}

// ComponentPreview represents a preview of what would be generated
type ComponentPreview struct {
	Type       string   // Component type
	Name       string   // Component name
	TargetPath string   // Where it would be generated
	Method     string   // "bootstrap" or "fallback"
	ToolUsed   string   // Tool that would be used
	Files      []string // Files that would be created
	Warnings   []string // Potential issues
}

// ValidationResult represents the result of validating a configuration or structure
type ValidationResult struct {
	Valid    bool     // Whether validation passed
	Errors   []string // Validation errors
	Warnings []string // Validation warnings
}

// MigrationResult represents the result of migrating a legacy project
type MigrationResult struct {
	Success           bool          // Whether migration succeeded
	ProjectRoot       string        // Root directory of migrated project
	Changes           []string      // List of changes made
	BackupPath        string        // Path to backup of original project
	ManualSteps       []string      // Manual steps required after migration
	ComponentsUpdated []string      // Components that were updated
	Duration          time.Duration // Migration time
	Errors            []error       // Any errors encountered
}
