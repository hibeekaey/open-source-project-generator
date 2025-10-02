// Package cli provides automation and integration support for the CLI
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EnvironmentConfig holds configuration loaded from environment variables
type EnvironmentConfig struct {
	// Project configuration
	ProjectName         string `env:"GENERATOR_PROJECT_NAME"`
	ProjectOrganization string `env:"GENERATOR_PROJECT_ORGANIZATION"`
	ProjectDescription  string `env:"GENERATOR_PROJECT_DESCRIPTION"`
	ProjectLicense      string `env:"GENERATOR_PROJECT_LICENSE"`
	OutputPath          string `env:"GENERATOR_OUTPUT_PATH"`

	// Generation options
	Force           bool   `env:"GENERATOR_FORCE"`
	Minimal         bool   `env:"GENERATOR_MINIMAL"`
	Offline         bool   `env:"GENERATOR_OFFLINE"`
	UpdateVersions  bool   `env:"GENERATOR_UPDATE_VERSIONS"`
	SkipValidation  bool   `env:"GENERATOR_SKIP_VALIDATION"`
	BackupExisting  bool   `env:"GENERATOR_BACKUP_EXISTING"`
	IncludeExamples bool   `env:"GENERATOR_INCLUDE_EXAMPLES"`
	Template        string `env:"GENERATOR_TEMPLATE"`

	// Component selection
	Frontend       bool `env:"GENERATOR_FRONTEND"`
	Backend        bool `env:"GENERATOR_BACKEND"`
	Mobile         bool `env:"GENERATOR_MOBILE"`
	Infrastructure bool `env:"GENERATOR_INFRASTRUCTURE"`

	// Technology selection
	FrontendTech       string `env:"GENERATOR_FRONTEND_TECH"`
	BackendTech        string `env:"GENERATOR_BACKEND_TECH"`
	MobileTech         string `env:"GENERATOR_MOBILE_TECH"`
	InfrastructureTech string `env:"GENERATOR_INFRASTRUCTURE_TECH"`

	// CLI behavior
	NonInteractive bool   `env:"GENERATOR_NON_INTERACTIVE"`
	OutputFormat   string `env:"GENERATOR_OUTPUT_FORMAT"`
	LogLevel       string `env:"GENERATOR_LOG_LEVEL"`
	Verbose        bool   `env:"GENERATOR_VERBOSE"`
	Quiet          bool   `env:"GENERATOR_QUIET"`
}

// CIEnvironment represents detected CI/CD environment information
type CIEnvironment struct {
	IsCI        bool   `json:"is_ci"`
	Provider    string `json:"provider"`
	BuildID     string `json:"build_id,omitempty"`
	BuildNumber string `json:"build_number,omitempty"`
	Branch      string `json:"branch,omitempty"`
	Commit      string `json:"commit,omitempty"`
	Repository  string `json:"repository,omitempty"`
	PullRequest string `json:"pull_request,omitempty"`
	JobID       string `json:"job_id,omitempty"`
	WorkflowID  string `json:"workflow_id,omitempty"`
	Environment string `json:"environment,omitempty"`
	Actor       string `json:"actor,omitempty"`
}

// CIDetector represents a CI environment detector with priority
type CIDetector struct {
	Name     string
	Priority int
	Detect   func() bool
	Extract  func() *CIEnvironment
}

// detectCIEnvironment detects if running in a CI/CD environment and returns details
func (c *CLI) detectCIEnvironment() *CIEnvironment {
	ci := &CIEnvironment{}

	// Define CI detectors with priority (higher priority = more specific)
	detectors := []CIDetector{
		{
			Name:     "github-actions",
			Priority: 100,
			Detect:   func() bool { return os.Getenv("GITHUB_ACTIONS") == "true" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:        true,
					Provider:    "github-actions",
					BuildID:     os.Getenv("GITHUB_RUN_ID"),
					BuildNumber: os.Getenv("GITHUB_RUN_NUMBER"),
					Branch:      os.Getenv("GITHUB_REF_NAME"),
					Commit:      os.Getenv("GITHUB_SHA"),
					Repository:  os.Getenv("GITHUB_REPOSITORY"),
					PullRequest: os.Getenv("GITHUB_EVENT_NUMBER"),
					JobID:       os.Getenv("GITHUB_JOB"),
					WorkflowID:  os.Getenv("GITHUB_WORKFLOW"),
					Environment: os.Getenv("GITHUB_ENVIRONMENT"),
					Actor:       os.Getenv("GITHUB_ACTOR"),
				}
			},
		},
		{
			Name:     "gitlab-ci",
			Priority: 100,
			Detect:   func() bool { return os.Getenv("GITLAB_CI") == "true" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:        true,
					Provider:    "gitlab-ci",
					BuildID:     os.Getenv("CI_PIPELINE_ID"),
					BuildNumber: os.Getenv("CI_PIPELINE_IID"),
					Branch:      os.Getenv("CI_COMMIT_REF_NAME"),
					Commit:      os.Getenv("CI_COMMIT_SHA"),
					Repository:  os.Getenv("CI_PROJECT_PATH"),
					PullRequest: os.Getenv("CI_MERGE_REQUEST_IID"),
					JobID:       os.Getenv("CI_JOB_ID"),
					Environment: os.Getenv("CI_ENVIRONMENT_NAME"),
					Actor:       os.Getenv("GITLAB_USER_LOGIN"),
				}
			},
		},
		{
			Name:     "jenkins",
			Priority: 90,
			Detect:   func() bool { return os.Getenv("JENKINS_URL") != "" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:        true,
					Provider:    "jenkins",
					BuildID:     os.Getenv("BUILD_ID"),
					BuildNumber: os.Getenv("BUILD_NUMBER"),
					Branch:      os.Getenv("GIT_BRANCH"),
					Commit:      os.Getenv("GIT_COMMIT"),
					Repository:  os.Getenv("GIT_URL"),
					JobID:       os.Getenv("JOB_NAME"),
				}
			},
		},
		{
			Name:     "travis-ci",
			Priority: 90,
			Detect:   func() bool { return os.Getenv("TRAVIS") == "true" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:        true,
					Provider:    "travis-ci",
					BuildID:     os.Getenv("TRAVIS_BUILD_ID"),
					BuildNumber: os.Getenv("TRAVIS_BUILD_NUMBER"),
					Branch:      os.Getenv("TRAVIS_BRANCH"),
					Commit:      os.Getenv("TRAVIS_COMMIT"),
					Repository:  os.Getenv("TRAVIS_REPO_SLUG"),
					PullRequest: os.Getenv("TRAVIS_PULL_REQUEST"),
					JobID:       os.Getenv("TRAVIS_JOB_ID"),
				}
			},
		},
		{
			Name:     "circleci",
			Priority: 90,
			Detect:   func() bool { return os.Getenv("CIRCLECI") == "true" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:        true,
					Provider:    "circleci",
					BuildID:     os.Getenv("CIRCLE_BUILD_NUM"),
					BuildNumber: os.Getenv("CIRCLE_BUILD_NUM"),
					Branch:      os.Getenv("CIRCLE_BRANCH"),
					Commit:      os.Getenv("CIRCLE_SHA1"),
					Repository:  os.Getenv("CIRCLE_PROJECT_REPONAME"),
					PullRequest: os.Getenv("CIRCLE_PR_NUMBER"),
					JobID:       os.Getenv("CIRCLE_JOB"),
					WorkflowID:  os.Getenv("CIRCLE_WORKFLOW_ID"),
				}
			},
		},
		{
			Name:     "azure-devops",
			Priority: 90,
			Detect:   func() bool { return os.Getenv("TF_BUILD") == "True" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:        true,
					Provider:    "azure-devops",
					BuildID:     os.Getenv("BUILD_BUILDID"),
					BuildNumber: os.Getenv("BUILD_BUILDNUMBER"),
					Branch:      os.Getenv("BUILD_SOURCEBRANCH"),
					Commit:      os.Getenv("BUILD_SOURCEVERSION"),
					Repository:  os.Getenv("BUILD_REPOSITORY_NAME"),
					PullRequest: os.Getenv("SYSTEM_PULLREQUEST_PULLREQUESTID"),
					JobID:       os.Getenv("AGENT_JOBNAME"),
				}
			},
		},
		{
			Name:     "bitbucket-pipelines",
			Priority: 90,
			Detect:   func() bool { return os.Getenv("BITBUCKET_BUILD_NUMBER") != "" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:        true,
					Provider:    "bitbucket-pipelines",
					BuildNumber: os.Getenv("BITBUCKET_BUILD_NUMBER"),
					Branch:      os.Getenv("BITBUCKET_BRANCH"),
					Commit:      os.Getenv("BITBUCKET_COMMIT"),
					Repository:  os.Getenv("BITBUCKET_REPO_FULL_NAME"),
					PullRequest: os.Getenv("BITBUCKET_PR_ID"),
				}
			},
		},
		{
			Name:     "aws-codebuild",
			Priority: 90,
			Detect:   func() bool { return os.Getenv("CODEBUILD_BUILD_ID") != "" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:     true,
					Provider: "aws-codebuild",
					BuildID:  os.Getenv("CODEBUILD_BUILD_ID"),
					Branch:   os.Getenv("CODEBUILD_WEBHOOK_HEAD_REF"),
					Commit:   os.Getenv("CODEBUILD_RESOLVED_SOURCE_VERSION"),
				}
			},
		},
		{
			Name:     "buildkite",
			Priority: 90,
			Detect:   func() bool { return os.Getenv("BUILDKITE") == "true" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:        true,
					Provider:    "buildkite",
					BuildID:     os.Getenv("BUILDKITE_BUILD_ID"),
					BuildNumber: os.Getenv("BUILDKITE_BUILD_NUMBER"),
					Branch:      os.Getenv("BUILDKITE_BRANCH"),
					Commit:      os.Getenv("BUILDKITE_COMMIT"),
					Repository:  os.Getenv("BUILDKITE_REPO"),
					PullRequest: os.Getenv("BUILDKITE_PULL_REQUEST"),
					JobID:       os.Getenv("BUILDKITE_JOB_ID"),
				}
			},
		},
		{
			Name:     "drone",
			Priority: 90,
			Detect:   func() bool { return os.Getenv("DRONE") == "true" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:        true,
					Provider:    "drone",
					BuildID:     os.Getenv("DRONE_BUILD_NUMBER"),
					BuildNumber: os.Getenv("DRONE_BUILD_NUMBER"),
					Branch:      os.Getenv("DRONE_BRANCH"),
					Commit:      os.Getenv("DRONE_COMMIT"),
					Repository:  os.Getenv("DRONE_REPO"),
					PullRequest: os.Getenv("DRONE_PULL_REQUEST"),
				}
			},
		},
		{
			Name:     "teamcity",
			Priority: 90,
			Detect:   func() bool { return os.Getenv("TEAMCITY_VERSION") != "" },
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:        true,
					Provider:    "teamcity",
					BuildID:     os.Getenv("BUILD_NUMBER"),
					BuildNumber: os.Getenv("BUILD_NUMBER"),
					Branch:      os.Getenv("BUILD_VCS_BRANCH"),
					Commit:      os.Getenv("BUILD_VCS_NUMBER"),
				}
			},
		},
		{
			Name:     "generic",
			Priority: 10, // Lowest priority - fallback for generic CI detection
			Detect: func() bool {
				return os.Getenv("CI") == "true" || os.Getenv("CONTINUOUS_INTEGRATION") == "true"
			},
			Extract: func() *CIEnvironment {
				return &CIEnvironment{
					IsCI:     true,
					Provider: "generic",
				}
			},
		},
	}

	// Find the highest priority detector that matches
	var selectedDetector *CIDetector
	for _, detector := range detectors {
		if detector.Detect() {
			if selectedDetector == nil || detector.Priority > selectedDetector.Priority {
				selectedDetector = &detector
			}
		}
	}

	// Extract CI information using the selected detector
	if selectedDetector != nil {
		ci = selectedDetector.Extract()
	}

	return ci
}

// isNonInteractiveMode checks if the CLI should run in non-interactive mode
func (c *CLI) isNonInteractiveMode() bool {
	// Check explicit flag
	if c.rootCmd != nil {
		if nonInteractive, err := c.rootCmd.PersistentFlags().GetBool("non-interactive"); err == nil && nonInteractive {
			return true
		}
	}

	// Check environment variable
	if parseBoolEnv("GENERATOR_NON_INTERACTIVE", false) {
		return true
	}

	// Check CI environment
	ci := c.detectCIEnvironment()
	if ci.IsCI {
		return true
	}

	// Check if stdin is not a terminal (piped input)
	if !isTerminal() {
		return true
	}

	return false
}

// outputMachineReadable outputs data in machine-readable format
func (c *CLI) outputMachineReadable(data interface{}, format string) error {
	switch strings.ToLower(format) {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	case "yaml":
		// For now, output as JSON since we don't have yaml package imported
		// This can be enhanced later with proper YAML support
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// Helper functions

// parseBoolEnv parses a boolean environment variable with a default value
func parseBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

// isTerminal checks if stdin is a terminal (not piped)
func isTerminal() bool {
	// Simple check for terminal - this can be enhanced with proper terminal detection
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// CLIError represents a CLI error with exit code information
type CLIError struct {
	Message string
	Code    int
}

func (e *CLIError) Error() string {
	return e.Message
}

// NewCLIError creates a new CLI error with the specified message and exit code
func NewCLIError(message string, code int) *CLIError {
	return &CLIError{
		Message: message,
		Code:    code,
	}
}
