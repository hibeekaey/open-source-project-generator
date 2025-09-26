package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new version command
func NewVersionCommand(cli interfaces.CLIInterface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  "Display the current version of the open source project generator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVersionCommand(cmd, args, cli)
		},
	}
	cmd.Flags().Bool("json", false, "Output version information in JSON format")
	return cmd
}

// RunVersionCommand displays the current version of the generator
func RunVersionCommand(cmd *cobra.Command, args []string, cli interfaces.CLIInterface) error {
	// Get version manager from CLI
	versionManager := cli.GetVersionManager()
	if versionManager == nil {
		fmt.Println("dev")
		return nil
	}

	currentVersion := versionManager.GetCurrentVersion()

	// Check for various output format flags
	jsonOutput, _ := cmd.Flags().GetBool("json")
	format, _ := cmd.Flags().GetString("format")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	short, _ := cmd.Flags().GetBool("short")

	// Simple version output (just the version string) - this is the default
	if short || (!jsonOutput && format != "json" && outputFormat != "json") {
		fmt.Println(currentVersion)
		return nil
	}

	// JSON output for automation
	if jsonOutput || format == "json" || outputFormat == "json" {
		versionInfo := map[string]string{
			"version": currentVersion,
		}
		jsonBytes, err := json.Marshal(versionInfo)
		if err != nil {
			return fmt.Errorf("🚫 Couldn't format version info: %w", err)
		}
		fmt.Println(string(jsonBytes))
	}

	return nil
}
