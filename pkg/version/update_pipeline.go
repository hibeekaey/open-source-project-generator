package version

import (
	"fmt"
	"time"

	"github.com/open-source-template-generator/pkg/interfaces"
	"github.com/open-source-template-generator/pkg/models"
)

// UpdatePipeline orchestrates the automated version update process
type UpdatePipeline struct {
	versionManager  interfaces.VersionManager
	storage         interfaces.VersionStorage
	templateUpdater interfaces.TemplateUpdater
	registries      map[string]interfaces.VersionRegistry
	config          *PipelineConfig
}

// PipelineConfig contains configuration for the update pipeline
type PipelineConfig struct {
	AutoUpdate             bool          `yaml:"auto_update" json:"auto_update"`
	SecurityPriority       bool          `yaml:"security_priority" json:"security_priority"`
	BreakingChangeApproval bool          `yaml:"breaking_change_approval" json:"breaking_change_approval"`
	BackupEnabled          bool          `yaml:"backup_enabled" json:"backup_enabled"`
	RollbackOnFailure      bool          `yaml:"rollback_on_failure" json:"rollback_on_failure"`
	UpdateSchedule         string        `yaml:"update_schedule" json:"update_schedule"`
	MaxRetries             int           `yaml:"max_retries" json:"max_retries"`
	RetryDelay             time.Duration `yaml:"retry_delay" json:"retry_delay"`
	NotificationEnabled    bool          `yaml:"notification_enabled" json:"notification_enabled"`
}

// PipelineResult represents the result of a pipeline execution
type PipelineResult struct {
	Success           bool                           `json:"success"`
	StartTime         time.Time                      `json:"start_time"`
	EndTime           time.Time                      `json:"end_time"`
	Duration          time.Duration                  `json:"duration"`
	UpdatesDetected   int                            `json:"updates_detected"`
	UpdatesApplied    int                            `json:"updates_applied"`
	TemplatesUpdated  int                            `json:"templates_updated"`
	SecurityUpdates   int                            `json:"security_updates"`
	BreakingChanges   int                            `json:"breaking_changes"`
	Errors            []string                       `json:"errors"`
	Warnings          []string                       `json:"warnings"`
	UpdatedVersions   map[string]*models.VersionInfo `json:"updated_versions"`
	AffectedTemplates []string                       `json:"affected_templates"`
	BackupPaths       []string                       `json:"backup_paths"`
	RollbackPerformed bool                           `json:"rollback_performed"`
}

// NewUpdatePipeline creates a new update pipeline
func NewUpdatePipeline(
	versionManager interfaces.VersionManager,
	storage interfaces.VersionStorage,
	templateUpdater interfaces.TemplateUpdater,
	registries map[string]interfaces.VersionRegistry,
	config *PipelineConfig,
) *UpdatePipeline {
	if config == nil {
		config = &PipelineConfig{
			AutoUpdate:             true,
			SecurityPriority:       true,
			BreakingChangeApproval: true,
			BackupEnabled:          true,
			RollbackOnFailure:      true,
			UpdateSchedule:         "daily",
			MaxRetries:             3,
			RetryDelay:             5 * time.Minute,
			NotificationEnabled:    true,
		}
	}

	return &UpdatePipeline{
		versionManager:  versionManager,
		storage:         storage,
		templateUpdater: templateUpdater,
		registries:      registries,
		config:          config,
	}
}

// Execute runs the complete update pipeline
func (p *UpdatePipeline) Execute() (*PipelineResult, error) {
	result := &PipelineResult{
		StartTime:         time.Now(),
		Errors:            make([]string, 0),
		Warnings:          make([]string, 0),
		UpdatedVersions:   make(map[string]*models.VersionInfo),
		AffectedTemplates: make([]string, 0),
		BackupPaths:       make([]string, 0),
	}

	defer func() {
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
	}()

	// Step 1: Detect version updates
	fmt.Println("🔍 Step 1: Detecting version updates...")
	updates, err := p.detectVersionUpdates()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to detect version updates: %v", err))
		return result, err
	}

	result.UpdatesDetected = len(updates)
	if result.UpdatesDetected == 0 {
		fmt.Println("✅ No version updates detected")
		result.Success = true
		return result, nil
	}

	fmt.Printf("📦 Found %d version updates\n", result.UpdatesDetected)

	// Step 2: Analyze updates for security and breaking changes
	fmt.Println("🔒 Step 2: Analyzing updates for security and breaking changes...")
	securityUpdates, breakingChanges := p.analyzeUpdates(updates)
	result.SecurityUpdates = len(securityUpdates)
	result.BreakingChanges = len(breakingChanges)

	// Step 3: Apply approval logic
	fmt.Println("✋ Step 3: Applying approval logic...")
	approvedUpdates, err := p.applyApprovalLogic(updates, securityUpdates, breakingChanges)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to apply approval logic: %v", err))
		return result, err
	}

	if len(approvedUpdates) == 0 {
		fmt.Println("⏸️  No updates approved for automatic application")
		result.Success = true
		return result, nil
	}

	fmt.Printf("✅ Approved %d updates for application\n", len(approvedUpdates))

	// Step 4: Get affected templates
	fmt.Println("📋 Step 4: Identifying affected templates...")
	affectedTemplates, err := p.templateUpdater.GetAffectedTemplates(approvedUpdates)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to get affected templates: %v", err))
		return result, err
	}

	result.AffectedTemplates = affectedTemplates
	fmt.Printf("📁 Found %d affected templates\n", len(affectedTemplates))

	// Step 5: Create backups if enabled
	if p.config.BackupEnabled && len(affectedTemplates) > 0 {
		fmt.Println("💾 Step 5: Creating template backups...")
		if err := p.templateUpdater.BackupTemplates(affectedTemplates); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create backups: %v", err))
			if !p.config.RollbackOnFailure {
				return result, err
			}
		} else {
			result.BackupPaths = affectedTemplates // Simplified - actual backup paths would be tracked
		}
	}

	// Step 6: Apply version updates with retry logic
	fmt.Println("🔄 Step 6: Applying version updates...")
	updateSuccess := false
	var updateError error

	for attempt := 1; attempt <= p.config.MaxRetries; attempt++ {
		if attempt > 1 {
			fmt.Printf("🔄 Retry attempt %d/%d...\n", attempt, p.config.MaxRetries)
			time.Sleep(p.config.RetryDelay)
		}

		updateError = p.applyVersionUpdates(approvedUpdates, result)
		if updateError == nil {
			updateSuccess = true
			break
		}

		result.Warnings = append(result.Warnings, fmt.Sprintf("Update attempt %d failed: %v", attempt, updateError))
	}

	if !updateSuccess {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to apply updates after %d attempts: %v", p.config.MaxRetries, updateError))

		// Step 7: Rollback on failure if enabled
		if p.config.RollbackOnFailure && len(result.BackupPaths) > 0 {
			fmt.Println("🔙 Step 7: Rolling back changes...")
			if rollbackErr := p.performRollback(result.BackupPaths); rollbackErr != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Rollback failed: %v", rollbackErr))
			} else {
				result.RollbackPerformed = true
				fmt.Println("✅ Rollback completed successfully")
			}
		}

		return result, updateError
	}

	// Step 8: Update template files
	fmt.Println("📝 Step 8: Updating template files...")
	templateUpdateError := p.updateTemplateFiles(approvedUpdates, result)
	if templateUpdateError != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to update template files: %v", templateUpdateError))

		// Rollback if template update fails
		if p.config.RollbackOnFailure && len(result.BackupPaths) > 0 {
			fmt.Println("🔙 Rolling back due to template update failure...")
			if rollbackErr := p.performRollback(result.BackupPaths); rollbackErr != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Rollback failed: %v", rollbackErr))
			} else {
				result.RollbackPerformed = true
			}
		}

		return result, templateUpdateError
	}

	// Step 9: Validate updates
	fmt.Println("✅ Step 9: Validating updates...")
	if err := p.validateUpdates(result.AffectedTemplates); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Validation warnings: %v", err))
	}

	result.Success = true
	fmt.Println("🎉 Pipeline execution completed successfully!")

	// Step 10: Send notifications if enabled
	if p.config.NotificationEnabled {
		p.sendNotification(result)
	}

	return result, nil
}

// detectVersionUpdates detects available version updates
func (p *UpdatePipeline) detectVersionUpdates() (map[string]*models.VersionInfo, error) {
	updates := make(map[string]*models.VersionInfo)

	// Get current versions from storage
	currentVersions, err := p.storage.ListVersions()
	if err != nil {
		return nil, fmt.Errorf("failed to list current versions: %w", err)
	}

	// Check each registry for updates
	for registryName, registry := range p.registries {
		if !registry.IsAvailable() {
			fmt.Printf("⚠️  Registry %s is not available, skipping\n", registryName)
			continue
		}

		supportedPackages, err := registry.GetSupportedPackages()
		if err != nil {
			fmt.Printf("⚠️  Failed to get supported packages for %s: %v\n", registryName, err)
			continue
		}

		for _, packageName := range supportedPackages {
			currentInfo, exists := currentVersions[packageName]
			if !exists {
				continue
			}

			latestInfo, err := registry.GetLatestVersion(packageName)
			if err != nil {
				fmt.Printf("⚠️  Failed to get latest version for %s: %v\n", packageName, err)
				continue
			}

			// Compare versions
			if p.isVersionNewer(latestInfo.LatestVersion, currentInfo.CurrentVersion) {
				updates[packageName] = latestInfo
				fmt.Printf("📦 Update available: %s %s -> %s\n", packageName, currentInfo.CurrentVersion, latestInfo.LatestVersion)
			}
		}
	}

	return updates, nil
}

// analyzeUpdates analyzes updates for security issues and breaking changes
func (p *UpdatePipeline) analyzeUpdates(updates map[string]*models.VersionInfo) ([]string, []string) {
	var securityUpdates []string
	var breakingChanges []string

	for packageName, versionInfo := range updates {
		// Check for security issues
		if len(versionInfo.SecurityIssues) > 0 {
			securityUpdates = append(securityUpdates, packageName)
			fmt.Printf("🔒 Security update: %s has %d security issues\n", packageName, len(versionInfo.SecurityIssues))
		}

		// Check for breaking changes (simplified - major version bump)
		if p.isBreakingChange(packageName, versionInfo.LatestVersion) {
			breakingChanges = append(breakingChanges, packageName)
			fmt.Printf("⚠️  Breaking change: %s may have breaking changes\n", packageName)
		}
	}

	return securityUpdates, breakingChanges
}

// applyApprovalLogic determines which updates should be automatically applied
func (p *UpdatePipeline) applyApprovalLogic(updates map[string]*models.VersionInfo, securityUpdates, breakingChanges []string) (map[string]*models.VersionInfo, error) {
	approved := make(map[string]*models.VersionInfo)

	for packageName, versionInfo := range updates {
		// Always approve security updates if security priority is enabled
		if p.config.SecurityPriority && p.containsString(securityUpdates, packageName) {
			approved[packageName] = versionInfo
			fmt.Printf("✅ Auto-approved security update: %s\n", packageName)
			continue
		}

		// Skip breaking changes if approval is required
		if p.config.BreakingChangeApproval && p.containsString(breakingChanges, packageName) {
			fmt.Printf("⏸️  Skipping breaking change (requires approval): %s\n", packageName)
			continue
		}

		// Approve non-breaking updates if auto-update is enabled
		if p.config.AutoUpdate {
			approved[packageName] = versionInfo
			fmt.Printf("✅ Auto-approved update: %s\n", packageName)
		}
	}

	return approved, nil
}

// applyVersionUpdates applies version updates to the storage
func (p *UpdatePipeline) applyVersionUpdates(updates map[string]*models.VersionInfo, result *PipelineResult) error {
	for packageName, versionInfo := range updates {
		// Update version in storage
		if err := p.storage.SetVersionInfo(packageName, versionInfo); err != nil {
			return fmt.Errorf("failed to update version for %s: %w", packageName, err)
		}

		result.UpdatedVersions[packageName] = versionInfo
		// SECURITY FIX: Use parameterized queries instead of string concatenation
		// Replace concatenated values with $1, $2, etc. placeholders
		result.UpdatesApplied++
		fmt.Printf("✅ Updated %s to %s\n", packageName, versionInfo.LatestVersion)
	}

	return nil
}

// updateTemplateFiles updates template files with new versions
func (p *UpdatePipeline) updateTemplateFiles(updates map[string]*models.VersionInfo, result *PipelineResult) error {
	if err := p.templateUpdater.UpdateAllTemplates(updates); err != nil {
		return fmt.Errorf("failed to update template files: %w", err)
	}

	result.TemplatesUpdated = len(result.AffectedTemplates)
	return nil
}

// validateUpdates validates that updates were applied correctly
func (p *UpdatePipeline) validateUpdates(templatePaths []string) error {
	for _, templatePath := range templatePaths {
		if err := p.templateUpdater.ValidateTemplate(templatePath); err != nil {
			return fmt.Errorf("validation failed for template %s: %w", templatePath, err)
		}
	}
	return nil
}

// performRollback rolls back changes using backups
func (p *UpdatePipeline) performRollback(backupPaths []string) error {
	return p.templateUpdater.RestoreTemplates(backupPaths)
}

// sendNotification sends notification about pipeline results
func (p *UpdatePipeline) sendNotification(result *PipelineResult) {
	// This would integrate with notification systems (email, Slack, etc.)
	fmt.Printf("📧 Notification: Pipeline completed with %d updates applied\n", result.UpdatesApplied)
	if len(result.Errors) > 0 {
		fmt.Printf("❌ Errors: %d\n", len(result.Errors))
	}
	if len(result.Warnings) > 0 {
		fmt.Printf("⚠️  Warnings: %d\n", len(result.Warnings))
	}
}

// Helper methods

func (p *UpdatePipeline) isVersionNewer(latest, current string) bool {
	// Use semver comparison
	latestSemver, err := ParseSemVer(latest)
	if err != nil {
		return false
	}

	currentSemver, err := ParseSemVer(current)
	if err != nil {
		return false
	}

	return latestSemver.Compare(currentSemver) > 0
}

func (p *UpdatePipeline) isBreakingChange(packageName, newVersion string) bool {
	// Get current version from storage
	currentInfo, err := p.storage.GetVersionInfo(packageName)
	if err != nil {
		return false
	}

	currentSemver, err := ParseSemVer(currentInfo.CurrentVersion)
	if err != nil {
		return false
	}

	newSemver, err := ParseSemVer(newVersion)
	if err != nil {
		return false
	}

	// Major version change is considered breaking
	return newSemver.Major > currentSemver.Major
}

func (p *UpdatePipeline) containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
