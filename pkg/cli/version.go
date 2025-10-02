package cli

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/spf13/cobra"
)

// VersionInfo represents the structured version information for JSON output
type VersionInfo struct {
	Version      string `json:"version"`
	GitCommit    string `json:"git_commit"`
	BuildTime    string `json:"build_time"`
	GoVersion    string `json:"go_version"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
}

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
	cmd.Flags().String("format", "", "Output format (json, text)")
	cmd.Flags().String("output-format", "", "Output format (json, text)")
	cmd.Flags().Bool("short", false, "Output only the version string")
	return cmd
}

// RunVersionCommand displays the current version of the generator
func RunVersionCommand(cmd *cobra.Command, args []string, cli interfaces.CLIInterface) error {
	// Check for various output format flags
	jsonOutput, _ := cmd.Flags().GetBool("json")
	format, _ := cmd.Flags().GetString("format")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	short, _ := cmd.Flags().GetBool("short")

	// Determine if JSON output is requested
	isJSONOutput := jsonOutput || format == "json" || outputFormat == "json"

	// Get version manager from CLI
	versionManager := cli.GetVersionManager()
	var currentVersion string
	if versionManager == nil {
		currentVersion = "dev"
	} else {
		currentVersion = versionManager.GetCurrentVersion()
	}

	// Simple version output (just the version string) - this is the default for text output
	if short || (!isJSONOutput) {
		fmt.Println(currentVersion)
		return nil
	}

	// JSON output for automation
	if isJSONOutput {
		versionInfo, err := buildVersionInfo(currentVersion, cli)
		if err != nil {
			return fmt.Errorf("failed to build version information: %w", err)
		}

		jsonBytes, err := json.MarshalIndent(versionInfo, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format version info as JSON: %w", err)
		}

		// Validate JSON output before printing
		if !json.Valid(jsonBytes) {
			return fmt.Errorf("generated invalid JSON output")
		}

		fmt.Println(string(jsonBytes))
	}

	return nil
}

// buildVersionInfo creates a structured VersionInfo with all available build information
func buildVersionInfo(version string, cli interfaces.CLIInterface) (*VersionInfo, error) {
	var gitCommit, buildTime string

	// Get build information from CLI if available
	if cli != nil && !reflect.ValueOf(cli).IsNil() {
		_, gitCommit, buildTime = cli.GetBuildInfo()
	}

	// Provide default values if build info is not available
	if gitCommit == "" {
		gitCommit = "unknown"
	}
	if buildTime == "" {
		buildTime = "unknown"
	}

	versionInfo := &VersionInfo{
		Version:      version,
		GitCommit:    gitCommit,
		BuildTime:    buildTime,
		GoVersion:    runtime.Version(),
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
	}

	return versionInfo, nil
}
