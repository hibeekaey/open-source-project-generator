package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
)

// setupCacheCommand sets up the cache command and its subcommands
func (c *CLI) setupCacheCommand() {
	cacheCmd := &cobra.Command{
		Use:   "cache <command> [flags]",
		Short: "Manage local cache for offline mode and performance",
		Long: `Manage local cache for templates, package versions, and other data.

Enables offline mode operation and improves performance through intelligent caching.`,
	}

	// cache show
	cacheShowCmd := &cobra.Command{
		Use:   "show",
		Short: "Show cache status and statistics",
		Long:  "Display cache location, size, statistics, and health information",
		RunE:  c.runCacheShow,
	}
	cacheCmd.AddCommand(cacheShowCmd)

	// cache clear
	cacheClearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all cached data",
		Long:  "Remove all cached templates, versions, and other data",
		RunE:  c.runCacheClear,
	}
	cacheClearCmd.Flags().Bool("force", false, "Clear cache without confirmation")
	cacheCmd.AddCommand(cacheClearCmd)

	// cache clean
	cacheCleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Remove expired and invalid cache entries",
		Long:  "Clean up expired cache entries and repair corrupted cache data",
		RunE:  c.runCacheClean,
	}
	cacheCmd.AddCommand(cacheCleanCmd)

	// cache validate
	cacheValidateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate cache integrity",
		Long:  "Check cache integrity and report any issues",
		RunE:  c.runCacheValidate,
	}
	cacheCmd.AddCommand(cacheValidateCmd)

	// cache repair
	cacheRepairCmd := &cobra.Command{
		Use:   "repair",
		Short: "Repair corrupted cache data",
		Long:  "Attempt to repair corrupted cache entries and restore functionality",
		RunE:  c.runCacheRepair,
	}
	cacheCmd.AddCommand(cacheRepairCmd)

	// cache offline enable
	cacheOfflineEnableCmd := &cobra.Command{
		Use:   "offline-enable",
		Short: "Enable offline mode",
		Long:  "Enable offline mode to work without internet connectivity",
		RunE:  c.runCacheOfflineEnable,
	}
	cacheCmd.AddCommand(cacheOfflineEnableCmd)

	// cache offline disable
	cacheOfflineDisableCmd := &cobra.Command{
		Use:   "offline-disable",
		Short: "Disable offline mode",
		Long:  "Disable offline mode and restore internet connectivity",
		RunE:  c.runCacheOfflineDisable,
	}
	cacheCmd.AddCommand(cacheOfflineDisableCmd)

	// cache offline status
	cacheOfflineStatusCmd := &cobra.Command{
		Use:   "offline-status",
		Short: "Show offline mode status",
		Long:  "Display current offline mode status and configuration",
		RunE:  c.runCacheOfflineStatus,
	}
	cacheCmd.AddCommand(cacheOfflineStatusCmd)

	c.rootCmd.AddCommand(cacheCmd)
}

// runCacheShow shows cache status and statistics
func (c *CLI) runCacheShow(cmd *cobra.Command, args []string) error {
	c.SuccessOutput("📊 Cache Status and Statistics")
	c.QuietOutput("")

	// Get cache statistics
	stats, err := c.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get cache statistics: %w", err)
	}

	// Display cache information
	c.QuietOutput("Cache Location: %s", c.highlight(stats.CacheLocation))
	c.QuietOutput("Cache Size: %s", c.highlight(formatBytes(stats.TotalSize)))
	c.QuietOutput("Total Entries: %s", c.highlight(strconv.Itoa(stats.TotalEntries)))
	c.QuietOutput("Expired Entries: %s", c.highlight(strconv.Itoa(stats.ExpiredEntries)))
	c.QuietOutput("Last Cleanup: %s", c.highlight(stats.LastCleanup.Format(time.RFC3339)))
	c.QuietOutput("Hit Rate: %s", c.highlight(fmt.Sprintf("%.1f%%", stats.HitRate)))
	c.QuietOutput("Cache Health: %s", c.highlight(stats.CacheHealth))

	return nil
}

// runCacheClear clears all cached data
func (c *CLI) runCacheClear(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	if !force {
		c.WarningOutput("⚠️  This will clear all cached data. This action cannot be undone.")
		c.QuietOutput("")
		
		// In non-interactive mode, require force flag
		if c.isNonInteractiveMode() {
			return fmt.Errorf("use --force flag to clear cache in non-interactive mode")
		}

		// Interactive confirmation
		ctx := cmd.Context()
		config := interfaces.ConfirmConfig{
			Prompt:       "Are you sure you want to clear all cached data?",
			DefaultValue: false,
		}

		result, err := c.interactiveUI.PromptConfirm(ctx, config)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}

		if !result.Confirmed {
			c.QuietOutput("Cache clear cancelled.")
			return nil
		}
	}

	c.VerboseOutput("Clearing cache...")
	err := c.cacheManager.Clear()
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	c.SuccessOutput("✅ Cache cleared successfully")
	return nil
}

// runCacheClean removes expired and invalid cache entries
func (c *CLI) runCacheClean(cmd *cobra.Command, args []string) error {
	c.VerboseOutput("Cleaning cache...")
	
	err := c.cacheManager.Clean()
	if err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}

	c.SuccessOutput("✅ Cache cleaned successfully")
	return nil
}

// runCacheValidate validates cache integrity
func (c *CLI) runCacheValidate(cmd *cobra.Command, args []string) error {
	c.VerboseOutput("Validating cache integrity...")
	
	err := c.cacheManager.ValidateCache()
	if err != nil {
		return fmt.Errorf("cache validation failed: %w", err)
	}

	c.SuccessOutput("✅ Cache validation passed")
	return nil
}

// runCacheRepair repairs corrupted cache data
func (c *CLI) runCacheRepair(cmd *cobra.Command, args []string) error {
	c.VerboseOutput("Repairing cache...")
	
	err := c.cacheManager.RepairCache()
	if err != nil {
		return fmt.Errorf("failed to repair cache: %w", err)
	}

	c.SuccessOutput("✅ Cache repair completed")
	return nil
}

// runCacheOfflineEnable enables offline mode
func (c *CLI) runCacheOfflineEnable(cmd *cobra.Command, args []string) error {
	return c.cacheManager.EnableOfflineMode()
}

// runCacheOfflineDisable disables offline mode
func (c *CLI) runCacheOfflineDisable(cmd *cobra.Command, args []string) error {
	return c.cacheManager.DisableOfflineMode()
}

// runCacheOfflineStatus shows offline mode status
func (c *CLI) runCacheOfflineStatus(cmd *cobra.Command, args []string) error {
	c.SuccessOutput("📡 Offline Mode Status")
	c.QuietOutput("")

	// Get cache statistics to check offline status
	stats, err := c.cacheManager.GetStats()
	if err != nil {
		return fmt.Errorf("failed to get cache statistics: %w", err)
	}

	// Check if offline mode is enabled
	offlineEnabled := stats.OfflineMode
	status := "Disabled"
	if offlineEnabled {
		status = "Enabled"
	}

	c.QuietOutput("Offline Mode: %s", c.highlight(status))
	
	if offlineEnabled {
		c.QuietOutput("Sync Status: %s", c.highlight("Offline mode active"))
	} else {
		c.QuietOutput("Sync Status: %s", c.highlight("Real-time"))
	}

	return nil
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
