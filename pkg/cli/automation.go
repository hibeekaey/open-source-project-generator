// Package cli provides automation and integration support for the CLI
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
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

// loadEnvironmentConfig loads configuration from environment variables
func (c *CLI) loadEnvironmentConfig() (*EnvironmentConfig, error) {
	config := &EnvironmentConfig{}

	// Project configuration
	config.ProjectName = os.Getenv("GENERATOR_PROJECT_NAME")
	config.ProjectOrganization = os.Getenv("GENERATOR_PROJECT_ORGANIZATION")
	config.ProjectDescription = os.Getenv("GENERATOR_PROJECT_DESCRIPTION")
	config.ProjectLicense = os.Getenv("GENERATOR_PROJECT_LICENSE")
	config.OutputPath = os.Getenv("GENERATOR_OUTPUT_PATH")

	// Generation options
	config.Force = parseBoolEnv("GENERATOR_FORCE", false)
	config.Minimal = parseBoolEnv("GENERATOR_MINIMAL", false)
	config.Offline = parseBoolEnv("GENERATOR_OFFLINE", false)
	config.UpdateVersions = parseBoolEnv("GENERATOR_UPDATE_VERSIONS", false)
	config.SkipValidation = parseBoolEnv("GENERATOR_SKIP_VALIDATION", false)
	config.BackupExisting = parseBoolEnv("GENERATOR_BACKUP_EXISTING", true)
	config.IncludeExamples = parseBoolEnv("GENERATOR_INCLUDE_EXAMPLES", true)
	config.Template = os.Getenv("GENERATOR_TEMPLATE")

	// Component selection
	config.Frontend = parseBoolEnv("GENERATOR_FRONTEND", false)
	config.Backend = parseBoolEnv("GENERATOR_BACKEND", false)
	config.Mobile = parseBoolEnv("GENERATOR_MOBILE", false)
	config.Infrastructure = parseBoolEnv("GENERATOR_INFRASTRUCTURE", false)

	// Technology selection
	config.FrontendTech = os.Getenv("GENERATOR_FRONTEND_TECH")
	config.BackendTech = os.Getenv("GENERATOR_BACKEND_TECH")
	config.MobileTech = os.Getenv("GENERATOR_MOBILE_TECH")
	config.InfrastructureTech = os.Getenv("GENERATOR_INFRASTRUCTURE_TECH")

	// CLI behavior
	config.NonInteractive = parseBoolEnv("GENERATOR_NON_INTERACTIVE", false)
	config.OutputFormat = getEnvWithDefault("GENERATOR_OUTPUT_FORMAT", "text")
	config.LogLevel = getEnvWithDefault("GENERATOR_LOG_LEVEL", "info")
	config.Verbose = parseBoolEnv("GENERATOR_VERBOSE", false)
	config.Quiet = parseBoolEnv("GENERATOR_QUIET", false)

	return config, nil
}

// detectCIEnvironment detects if running in a CI/CD environment and returns details
func (c *CLI) detectCIEnvironment() *CIEnvironment {
	ci := &CIEnvironment{}

	// Check for common CI environment variables
	if os.Getenv("CI") == "true" || os.Getenv("CONTINUOUS_INTEGRATION") == "true" {
		ci.IsCI = true
	}

	// GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		ci.IsCI = true
		ci.Provider = "github-actions"
		ci.BuildID = os.Getenv("GITHUB_RUN_ID")
		ci.BuildNumber = os.Getenv("GITHUB_RUN_NUMBER")
		ci.Branch = os.Getenv("GITHUB_REF_NAME")
		ci.Commit = os.Getenv("GITHUB_SHA")
		ci.Repository = os.Getenv("GITHUB_REPOSITORY")
		ci.PullRequest = os.Getenv("GITHUB_EVENT_NUMBER")
		ci.JobID = os.Getenv("GITHUB_JOB")
		ci.WorkflowID = os.Getenv("GITHUB_WORKFLOW")
		ci.Environment = os.Getenv("GITHUB_ENVIRONMENT")
		ci.Actor = os.Getenv("GITHUB_ACTOR")
	}

	// GitLab CI
	if os.Getenv("GITLAB_CI") == "true" {
		ci.IsCI = true
		ci.Provider = "gitlab-ci"
		ci.BuildID = os.Getenv("CI_PIPELINE_ID")
		ci.BuildNumber = os.Getenv("CI_PIPELINE_IID")
		ci.Branch = os.Getenv("CI_COMMIT_REF_NAME")
		ci.Commit = os.Getenv("CI_COMMIT_SHA")
		ci.Repository = os.Getenv("CI_PROJECT_PATH")
		ci.PullRequest = os.Getenv("CI_MERGE_REQUEST_IID")
		ci.JobID = os.Getenv("CI_JOB_ID")
		ci.Environment = os.Getenv("CI_ENVIRONMENT_NAME")
		ci.Actor = os.Getenv("GITLAB_USER_LOGIN")
	}

	// Jenkins
	if os.Getenv("JENKINS_URL") != "" {
		ci.IsCI = true
		ci.Provider = "jenkins"
		ci.BuildID = os.Getenv("BUILD_ID")
		ci.BuildNumber = os.Getenv("BUILD_NUMBER")
		ci.Branch = os.Getenv("GIT_BRANCH")
		ci.Commit = os.Getenv("GIT_COMMIT")
		ci.Repository = os.Getenv("GIT_URL")
		ci.JobID = os.Getenv("JOB_NAME")
	}

	// Travis CI
	if os.Getenv("TRAVIS") == "true" {
		ci.IsCI = true
		ci.Provider = "travis-ci"
		ci.BuildID = os.Getenv("TRAVIS_BUILD_ID")
		ci.BuildNumber = os.Getenv("TRAVIS_BUILD_NUMBER")
		ci.Branch = os.Getenv("TRAVIS_BRANCH")
		ci.Commit = os.Getenv("TRAVIS_COMMIT")
		ci.Repository = os.Getenv("TRAVIS_REPO_SLUG")
		ci.PullRequest = os.Getenv("TRAVIS_PULL_REQUEST")
		ci.JobID = os.Getenv("TRAVIS_JOB_ID")
	}

	// CircleCI
	if os.Getenv("CIRCLECI") == "true" {
		ci.IsCI = true
		ci.Provider = "circleci"
		ci.BuildID = os.Getenv("CIRCLE_BUILD_NUM")
		ci.BuildNumber = os.Getenv("CIRCLE_BUILD_NUM")
		ci.Branch = os.Getenv("CIRCLE_BRANCH")
		ci.Commit = os.Getenv("CIRCLE_SHA1")
		ci.Repository = os.Getenv("CIRCLE_PROJECT_REPONAME")
		ci.PullRequest = os.Getenv("CIRCLE_PR_NUMBER")
		ci.JobID = os.Getenv("CIRCLE_JOB")
		ci.WorkflowID = os.Getenv("CIRCLE_WORKFLOW_ID")
	}

	// Azure DevOps
	if os.Getenv("TF_BUILD") == "True" {
		ci.IsCI = true
		ci.Provider = "azure-devops"
		ci.BuildID = os.Getenv("BUILD_BUILDID")
		ci.BuildNumber = os.Getenv("BUILD_BUILDNUMBER")
		ci.Branch = os.Getenv("BUILD_SOURCEBRANCH")
		ci.Commit = os.Getenv("BUILD_SOURCEVERSION")
		ci.Repository = os.Getenv("BUILD_REPOSITORY_NAME")
		ci.PullRequest = os.Getenv("SYSTEM_PULLREQUEST_PULLREQUESTID")
		ci.JobID = os.Getenv("AGENT_JOBNAME")
	}

	// Bitbucket Pipelines
	if os.Getenv("BITBUCKET_BUILD_NUMBER") != "" {
		ci.IsCI = true
		ci.Provider = "bitbucket-pipelines"
		ci.BuildNumber = os.Getenv("BITBUCKET_BUILD_NUMBER")
		ci.Branch = os.Getenv("BITBUCKET_BRANCH")
		ci.Commit = os.Getenv("BITBUCKET_COMMIT")
		ci.Repository = os.Getenv("BITBUCKET_REPO_FULL_NAME")
		ci.PullRequest = os.Getenv("BITBUCKET_PR_ID")
	}

	// AWS CodeBuild
	if os.Getenv("CODEBUILD_BUILD_ID") != "" {
		ci.IsCI = true
		ci.Provider = "aws-codebuild"
		ci.BuildID = os.Getenv("CODEBUILD_BUILD_ID")
		ci.Branch = os.Getenv("CODEBUILD_WEBHOOK_HEAD_REF")
		ci.Commit = os.Getenv("CODEBUILD_RESOLVED_SOURCE_VERSION")
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

// convertEnvironmentConfigToProjectConfig converts environment config to project config
func (c *CLI) convertEnvironmentConfigToProjectConfig(envConfig *EnvironmentConfig) (*models.ProjectConfig, error) {
	config := &models.ProjectConfig{
		Name:         envConfig.ProjectName,
		Organization: envConfig.ProjectOrganization,
		Description:  envConfig.ProjectDescription,
		License:      envConfig.ProjectLicense,
		OutputPath:   envConfig.OutputPath,
	}

	// Set components based on environment variables
	components := models.Components{}

	if envConfig.Frontend {
		// Set frontend components based on technology
		if envConfig.FrontendTech == "nextjs-app" || envConfig.FrontendTech == "" {
			components.Frontend.NextJS.App = true
			components.Frontend.NextJS.Home = true
			components.Frontend.NextJS.Shared = true
		}
	}

	if envConfig.Backend {
		// Set backend components based on technology
		if envConfig.BackendTech == "go-gin" || envConfig.BackendTech == "" {
			components.Backend.GoGin = true
		}
	}

	if envConfig.Mobile {
		// Set mobile components based on technology
		if envConfig.MobileTech == "android-kotlin" || envConfig.MobileTech == "android" {
			components.Mobile.Android = true
		}
		if envConfig.MobileTech == "ios-swift" || envConfig.MobileTech == "ios" {
			components.Mobile.IOS = true
		}
	}

	if envConfig.Infrastructure {
		// Set infrastructure components based on technology
		components.Infrastructure.Docker = true
		if envConfig.InfrastructureTech == "kubernetes" {
			components.Infrastructure.Kubernetes = true
		}
		if envConfig.InfrastructureTech == "terraform" {
			components.Infrastructure.Terraform = true
		}
	}

	config.Components = components

	return config, nil
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
		// This can be improved later with proper YAML support
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

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// isTerminal checks if stdin is a terminal (not piped)
func isTerminal() bool {
	// Basic check for terminal - this can be improved with proper terminal detection
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
