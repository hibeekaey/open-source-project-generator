package cleanup

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Manager coordinates the cleanup process
type Manager struct {
	projectRoot string
	backupMgr   *BackupManager
	analyzer    *CodeAnalyzer
	validator   *ValidationFramework
	config      *Config
	logger      *log.Logger
}

// Config holds cleanup configuration
type Config struct {
	BackupDir       string
	DryRun          bool
	Verbose         bool
	SkipPatterns    []string
	PreserveTODOs   []string
	MaxBackupAge    time.Duration
	ValidationLevel ValidationLevel
}

// ValidationLevel represents the level of validation to perform
type ValidationLevel int

const (
	ValidationBasic ValidationLevel = iota
	ValidationStandard
	ValidationStrict
)

// CleanupResult represents the overall result of cleanup operations
type CleanupResult struct {
	Success          bool
	FilesModified    []string
	IssuesFixed      []FixedIssue
	IssuesRemaining  []RemainingIssue
	ValidationResult *ValidationResult
	Summary          *CleanupSummary
	Duration         time.Duration
}

// FixedIssue represents an issue that was successfully fixed
type FixedIssue struct {
	Type        string
	File        string
	Line        int
	Description string
	Action      string
}

// RemainingIssue represents an issue that couldn't be automatically fixed
type RemainingIssue struct {
	Type        string
	File        string
	Line        int
	Description string
	Reason      string
	Suggestion  string
}

// CleanupSummary provides a summary of cleanup operations
type CleanupSummary struct {
	TotalFilesScanned int
	FilesModified     int
	TODOsAnalyzed     int
	TODOsResolved     int
	DuplicatesFound   int
	DuplicatesFixed   int
	UnusedCodeRemoved int
	ImportsOrganized  int
	TestsFixed        int
	ValidationsPassed bool
}

// NewManager creates a new cleanup manager
func NewManager(projectRoot string, config *Config) (*Manager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Ensure backup directory exists
	backupDir := config.BackupDir
	if backupDir == "" {
		backupDir = filepath.Join(projectRoot, ".cleanup-backups")
	}

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	logger := log.New(os.Stdout, "[CLEANUP] ", log.LstdFlags)
	if !config.Verbose {
		logger.SetOutput(os.Stderr)
	}

	return &Manager{
		projectRoot: projectRoot,
		backupMgr:   NewBackupManager(backupDir),
		analyzer:    NewCodeAnalyzer(),
		validator:   NewValidationFramework(projectRoot),
		config:      config,
		logger:      logger,
	}, nil
}

// DefaultConfig returns default cleanup configuration
func DefaultConfig() *Config {
	return &Config{
		BackupDir:       "",
		DryRun:          false,
		Verbose:         true,
		SkipPatterns:    []string{"vendor/", ".git/", "node_modules/"},
		PreserveTODOs:   []string{},
		MaxBackupAge:    7 * 24 * time.Hour, // 7 days
		ValidationLevel: ValidationStandard,
	}
}

// Initialize sets up the cleanup infrastructure
func (m *Manager) Initialize() error {
	m.logger.Println("Initializing cleanup infrastructure...")

	// 1. Validate project integrity
	if err := m.validator.EnsureProjectIntegrity(); err != nil {
		return fmt.Errorf("project integrity check failed: %w", err)
	}

	// 2. Create initial validation checkpoint
	initialValidation, err := m.validator.CreateValidationCheckpoint()
	if err != nil {
		return fmt.Errorf("failed to create initial validation checkpoint: %w", err)
	}

	if !initialValidation.Success {
		m.logger.Printf("Warning: Initial validation failed with %d errors", len(initialValidation.Errors))
		for _, err := range initialValidation.Errors {
			m.logger.Printf("  - %s: %s", err.Type, err.Message)
		}

		if m.config.ValidationLevel == ValidationStrict {
			return fmt.Errorf("project validation failed - cannot proceed with strict validation level")
		}
	}

	// 3. Clean up old backups
	if err := m.backupMgr.CleanupOldBackups(m.config.MaxBackupAge); err != nil {
		m.logger.Printf("Warning: Failed to cleanup old backups: %v", err)
	}

	m.logger.Println("Cleanup infrastructure initialized successfully")
	return nil
}

// CreateBackup creates a backup of specified files
func (m *Manager) CreateBackup(files []string) (*Backup, error) {
	if m.config.DryRun {
		m.logger.Printf("DRY RUN: Would create backup of %d files", len(files))
		return &Backup{
			ID:        fmt.Sprintf("dryrun_%d", time.Now().Unix()),
			Timestamp: time.Now(),
			Files:     []BackupFile{},
		}, nil
	}

	m.logger.Printf("Creating backup of %d files...", len(files))
	backup, err := m.backupMgr.CreateBackup(files)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	m.logger.Printf("Backup created successfully: %s", backup.ID)
	return backup, nil
}

// AnalyzeProject performs comprehensive project analysis
func (m *Manager) AnalyzeProject() (*ProjectAnalysis, error) {
	m.logger.Println("Starting project analysis...")

	analysis := &ProjectAnalysis{
		ProjectRoot: m.projectRoot,
		Timestamp:   time.Now(),
	}

	// 1. Analyze TODO comments
	m.logger.Println("Analyzing TODO/FIXME comments...")
	todos, err := m.analyzer.AnalyzeTODOComments(m.projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze TODO comments: %w", err)
	}
	analysis.TODOs = todos
	m.logger.Printf("Found %d TODO/FIXME comments", len(todos))

	// 2. Find duplicate code
	m.logger.Println("Scanning for duplicate code...")
	duplicates, err := m.analyzer.FindDuplicateCode(m.projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to find duplicate code: %w", err)
	}
	analysis.Duplicates = duplicates
	m.logger.Printf("Found %d potential duplicate code blocks", len(duplicates))

	// 3. Identify unused code
	m.logger.Println("Identifying unused code...")
	unused, err := m.analyzer.IdentifyUnusedCode(m.projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to identify unused code: %w", err)
	}
	analysis.UnusedCode = unused
	m.logger.Printf("Found %d unused code items", len(unused))

	// 4. Validate import organization
	m.logger.Println("Validating import organization...")
	importIssues, err := m.analyzer.ValidateImportOrganization(m.projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to validate imports: %w", err)
	}
	analysis.ImportIssues = importIssues
	m.logger.Printf("Found %d import organization issues", len(importIssues))

	m.logger.Println("Project analysis completed")
	return analysis, nil
}

// ValidateProject runs project validation
func (m *Manager) ValidateProject() (*ValidationResult, error) {
	m.logger.Println("Running project validation...")

	result, err := m.validator.ValidateProject()
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if result.Success {
		m.logger.Println("Project validation passed")
	} else {
		m.logger.Printf("Project validation failed with %d errors", len(result.Errors))
		for _, err := range result.Errors {
			m.logger.Printf("  - %s: %s", err.Type, err.Message)
		}
	}

	return result, nil
}

// Shutdown performs cleanup manager shutdown
func (m *Manager) Shutdown() error {
	m.logger.Println("Shutting down cleanup manager...")

	// Perform any necessary cleanup
	// For now, just log the shutdown

	m.logger.Println("Cleanup manager shutdown complete")
	return nil
}

// ProjectAnalysis holds the results of project analysis
type ProjectAnalysis struct {
	ProjectRoot  string
	Timestamp    time.Time
	TODOs        []TODOItem
	Duplicates   []DuplicateCodeBlock
	UnusedCode   []UnusedCodeItem
	ImportIssues []ImportIssue
}

// GetSummary returns a summary of the analysis
func (pa *ProjectAnalysis) GetSummary() string {
	return fmt.Sprintf(
		"Project Analysis Summary:\n"+
			"  TODO/FIXME comments: %d\n"+
			"  Duplicate code blocks: %d\n"+
			"  Unused code items: %d\n"+
			"  Import issues: %d\n",
		len(pa.TODOs),
		len(pa.Duplicates),
		len(pa.UnusedCode),
		len(pa.ImportIssues),
	)
}
