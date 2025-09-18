package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cuesoftinc/open-source-project-generator/pkg/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new version command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  "Display the current version of the open source project generator",
		RunE:  RunVersionCommand,
	}
	cmd.Flags().Bool("json", false, "Output version information in JSON format")
	return cmd
}

// RunVersionCommand displays the current version of the generator
func RunVersionCommand(cmd *cobra.Command, args []string) error {
	manager := version.NewManager()
	jsonOutput, _ := cmd.Flags().GetBool("json")
	format, _ := cmd.Flags().GetString("format")
	outputFormat, _ := cmd.Flags().GetString("output-format")

	currentVersion := manager.GetCurrentVersion()

	// Check if JSON output is requested via any flag
	if jsonOutput || format == "json" || outputFormat == "json" {
		versionInfo := map[string]string{
			"version": currentVersion,
			"build":   "development", // placeholder for build info
		}
		jsonBytes, err := json.Marshal(versionInfo)
		if err != nil {
			return fmt.Errorf("failed to marshal version info: %w", err)
		}
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("Generator version: %s\n", currentVersion)
	}
	return nil
}
