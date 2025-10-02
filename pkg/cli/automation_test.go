package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectCIEnvironment(t *testing.T) {
	// Store original environment to restore later
	originalEnv := make(map[string]string)
	ciEnvVars := []string{
		"CI", "CONTINUOUS_INTEGRATION", "GITHUB_ACTIONS", "GITLAB_CI", "JENKINS_URL",
		"TRAVIS", "CIRCLECI", "TF_BUILD", "BITBUCKET_BUILD_NUMBER", "CODEBUILD_BUILD_ID",
		"BUILDKITE", "DRONE", "TEAMCITY_VERSION",
	}

	// Store and clean environment
	for _, envVar := range ciEnvVars {
		originalEnv[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}

	// Restore environment after test
	defer func() {
		for envVar, value := range originalEnv {
			if value != "" {
				os.Setenv(envVar, value)
			} else {
				os.Unsetenv(envVar)
			}
		}
	}()

	tests := []struct {
		name             string
		envVar           string
		envValue         string
		expectedCI       bool
		expectedProvider string
	}{
		{
			name:             "GitHub Actions",
			envVar:           "GITHUB_ACTIONS",
			envValue:         "true",
			expectedCI:       true,
			expectedProvider: "github-actions",
		},
		{
			name:             "GitLab CI",
			envVar:           "GITLAB_CI",
			envValue:         "true",
			expectedCI:       true,
			expectedProvider: "gitlab-ci",
		},
		{
			name:             "Jenkins",
			envVar:           "JENKINS_URL",
			envValue:         "http://jenkins.example.com",
			expectedCI:       true,
			expectedProvider: "jenkins",
		},
		{
			name:             "Travis CI",
			envVar:           "TRAVIS",
			envValue:         "true",
			expectedCI:       true,
			expectedProvider: "travis-ci",
		},
		{
			name:             "CircleCI",
			envVar:           "CIRCLECI",
			envValue:         "true",
			expectedCI:       true,
			expectedProvider: "circleci",
		},
		{
			name:             "Azure DevOps",
			envVar:           "TF_BUILD",
			envValue:         "True",
			expectedCI:       true,
			expectedProvider: "azure-devops",
		},
		{
			name:             "Bitbucket Pipelines",
			envVar:           "BITBUCKET_BUILD_NUMBER",
			envValue:         "123",
			expectedCI:       true,
			expectedProvider: "bitbucket-pipelines",
		},
		{
			name:             "AWS CodeBuild",
			envVar:           "CODEBUILD_BUILD_ID",
			envValue:         "build-123",
			expectedCI:       true,
			expectedProvider: "aws-codebuild",
		},
		{
			name:             "Buildkite",
			envVar:           "BUILDKITE",
			envValue:         "true",
			expectedCI:       true,
			expectedProvider: "buildkite",
		},
		{
			name:             "Drone",
			envVar:           "DRONE",
			envValue:         "true",
			expectedCI:       true,
			expectedProvider: "drone",
		},
		{
			name:             "TeamCity",
			envVar:           "TEAMCITY_VERSION",
			envValue:         "2021.1",
			expectedCI:       true,
			expectedProvider: "teamcity",
		},
		{
			name:             "Generic CI",
			envVar:           "CI",
			envValue:         "true",
			expectedCI:       true,
			expectedProvider: "generic",
		},
		{
			name:             "Generic CONTINUOUS_INTEGRATION",
			envVar:           "CONTINUOUS_INTEGRATION",
			envValue:         "true",
			expectedCI:       true,
			expectedProvider: "generic",
		},
		{
			name:             "No CI environment",
			envVar:           "",
			envValue:         "",
			expectedCI:       false,
			expectedProvider: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment for each test
			for _, envVar := range ciEnvVars {
				os.Unsetenv(envVar)
			}

			// Set test environment variable
			if tt.envVar != "" {
				os.Setenv(tt.envVar, tt.envValue)
				defer os.Unsetenv(tt.envVar)
			}

			// Create CLI instance and test detection
			cli := &CLI{}
			result := cli.detectCIEnvironment()

			assert.Equal(t, tt.expectedCI, result.IsCI, "IsCI should match expected value")
			assert.Equal(t, tt.expectedProvider, result.Provider, "Provider should match expected value")
		})
	}
}

func TestDetectCIEnvironmentPriority(t *testing.T) {
	// Store original environment to restore later
	originalEnv := make(map[string]string)
	ciEnvVars := []string{"CI", "GITHUB_ACTIONS", "GITLAB_CI"}

	// Store and clean environment
	for _, envVar := range ciEnvVars {
		originalEnv[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}

	// Restore environment after test
	defer func() {
		for envVar, value := range originalEnv {
			if value != "" {
				os.Setenv(envVar, value)
			} else {
				os.Unsetenv(envVar)
			}
		}
	}()

	tests := []struct {
		name             string
		envVars          map[string]string
		expectedProvider string
		description      string
	}{
		{
			name: "GitHub Actions takes priority over generic CI",
			envVars: map[string]string{
				"CI":             "true",
				"GITHUB_ACTIONS": "true",
			},
			expectedProvider: "github-actions",
			description:      "Specific CI platforms should take priority over generic CI detection",
		},
		{
			name: "GitLab CI takes priority over generic CI",
			envVars: map[string]string{
				"CI":        "true",
				"GITLAB_CI": "true",
			},
			expectedProvider: "gitlab-ci",
			description:      "Specific CI platforms should take priority over generic CI detection",
		},
		{
			name: "GitHub Actions and GitLab CI both present - GitHub Actions wins",
			envVars: map[string]string{
				"GITHUB_ACTIONS": "true",
				"GITLAB_CI":      "true",
			},
			expectedProvider: "github-actions",
			description:      "When multiple specific CI platforms are detected, priority should be consistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment for each test
			for _, envVar := range ciEnvVars {
				os.Unsetenv(envVar)
			}

			// Set test environment variables
			for envVar, envValue := range tt.envVars {
				os.Setenv(envVar, envValue)
				defer os.Unsetenv(envVar)
			}

			// Create CLI instance and test detection
			cli := &CLI{}
			result := cli.detectCIEnvironment()

			assert.True(t, result.IsCI, "Should detect CI environment")
			assert.Equal(t, tt.expectedProvider, result.Provider, tt.description)
		})
	}
}

func TestDetectCIEnvironmentExtractedData(t *testing.T) {
	// Store original environment to restore later
	originalEnv := make(map[string]string)
	githubEnvVars := []string{
		"GITHUB_ACTIONS", "GITHUB_RUN_ID", "GITHUB_RUN_NUMBER", "GITHUB_REF_NAME",
		"GITHUB_SHA", "GITHUB_REPOSITORY", "GITHUB_JOB", "GITHUB_WORKFLOW",
	}

	// Store and clean environment
	for _, envVar := range githubEnvVars {
		originalEnv[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}

	// Restore environment after test
	defer func() {
		for envVar, value := range originalEnv {
			if value != "" {
				os.Setenv(envVar, value)
			} else {
				os.Unsetenv(envVar)
			}
		}
	}()

	t.Run("GitHub Actions extracts all data", func(t *testing.T) {
		// Set GitHub Actions environment
		testData := map[string]string{
			"GITHUB_ACTIONS":    "true",
			"GITHUB_RUN_ID":     "123456789",
			"GITHUB_RUN_NUMBER": "42",
			"GITHUB_REF_NAME":   "main",
			"GITHUB_SHA":        "abc123def456",
			"GITHUB_REPOSITORY": "owner/repo",
			"GITHUB_JOB":        "test-job",
			"GITHUB_WORKFLOW":   "CI",
		}

		for envVar, envValue := range testData {
			os.Setenv(envVar, envValue)
			defer os.Unsetenv(envVar)
		}

		// Create CLI instance and test detection
		cli := &CLI{}
		result := cli.detectCIEnvironment()

		assert.True(t, result.IsCI)
		assert.Equal(t, "github-actions", result.Provider)
		assert.Equal(t, "123456789", result.BuildID)
		assert.Equal(t, "42", result.BuildNumber)
		assert.Equal(t, "main", result.Branch)
		assert.Equal(t, "abc123def456", result.Commit)
		assert.Equal(t, "owner/repo", result.Repository)
		assert.Equal(t, "test-job", result.JobID)
		assert.Equal(t, "CI", result.WorkflowID)
	})
}
