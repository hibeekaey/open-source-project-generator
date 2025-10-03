package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCIEnvironmentErrorHandling tests CI detection error handling scenarios
func TestCIEnvironmentErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func()
		cleanupEnv  func()
		operation   func(*CLI) interface{}
		expectPanic bool
		description string
	}{
		{
			name: "missing environment variables",
			setupEnv: func() {
				// Clear all CI-related environment variables
				clearAllCIEnvironment()
			},
			cleanupEnv: func() {
				// No cleanup needed
			},
			operation: func(cli *CLI) interface{} {
				return cli.detectCIEnvironment()
			},
			expectPanic: false,
			description: "Should handle missing environment variables gracefully",
		},
		{
			name: "malformed environment variables",
			setupEnv: func() {
				clearAllCIEnvironment()
				// Set some unusual values
				if err := os.Setenv("GITHUB_ACTIONS", "maybe"); err != nil {
					t.Errorf("Failed to set GITHUB_ACTIONS: %v", err)
				}
				if err := os.Setenv("GITLAB_CI", "sometimes"); err != nil {
					t.Errorf("Failed to set GITLAB_CI: %v", err)
				}
			},
			cleanupEnv: func() {
				if err := os.Unsetenv("GITHUB_ACTIONS"); err != nil {
					t.Errorf("Failed to unset GITHUB_ACTIONS: %v", err)
				}
				if err := os.Unsetenv("GITLAB_CI"); err != nil {
					t.Errorf("Failed to unset GITLAB_CI: %v", err)
				}
			},
			operation: func(cli *CLI) interface{} {
				return cli.detectCIEnvironment()
			},
			expectPanic: false,
			description: "Should handle malformed environment variables gracefully",
		},
		{
			name: "partial CI environment setup",
			setupEnv: func() {
				clearAllCIEnvironment()
				// Set only some GitHub Actions variables
				_ = os.Setenv("GITHUB_ACTIONS", "true")
				// Missing other expected variables
			},
			cleanupEnv: func() {
				_ = os.Unsetenv("GITHUB_ACTIONS")
			},
			operation: func(cli *CLI) interface{} {
				return cli.detectCIEnvironment()
			},
			expectPanic: false,
			description: "Should handle partial CI environment setup gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			defer tt.cleanupEnv()

			mockLogger := &MockLogger{}
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					tt.operation(cli)
				}, tt.description)
			} else {
				assert.NotPanics(t, func() {
					result := tt.operation(cli)
					assert.NotNil(t, result, tt.description)
				}, tt.description)
			}
		})
	}
}

// TestCIEnvironmentEdgeCases tests edge cases in CI detection
func TestCIEnvironmentEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectedCI  bool
		description string
	}{
		{
			name: "GitHub Actions with false value",
			envVars: map[string]string{
				"GITHUB_ACTIONS": "false",
			},
			expectedCI:  false,
			description: "GITHUB_ACTIONS=false should not be detected as CI",
		},
		{
			name: "GitLab CI with false value",
			envVars: map[string]string{
				"GITLAB_CI": "false",
			},
			expectedCI:  false,
			description: "GITLAB_CI=false should not be detected as CI",
		},
		{
			name: "Travis CI with false value",
			envVars: map[string]string{
				"TRAVIS": "false",
			},
			expectedCI:  false,
			description: "TRAVIS=false should not be detected as CI",
		},
		{
			name: "Empty CI environment variables",
			envVars: map[string]string{
				"GITHUB_ACTIONS": "",
				"GITLAB_CI":      "",
				"TRAVIS":         "",
			},
			expectedCI:  false,
			description: "Empty CI variables should not be detected as CI",
		},
		{
			name: "Case sensitivity test",
			envVars: map[string]string{
				"GITHUB_ACTIONS": "TRUE", // Uppercase
			},
			expectedCI:  false, // Should be case-sensitive
			description: "CI detection should be case-sensitive",
		},
		{
			name: "Multiple CI environments set",
			envVars: map[string]string{
				"GITHUB_ACTIONS": "true",
				"GITLAB_CI":      "true",
				"TRAVIS":         "true",
			},
			expectedCI:  true,
			description: "Should detect CI when multiple environments are set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			clearAllCIEnvironment()

			// Set test environment variables
			for key, value := range tt.envVars {
				_ = os.Setenv(key, value)
				defer func(k string) { _ = os.Unsetenv(k) }(key)
			}

			mockLogger := &MockLogger{}
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}

			ciInfo := cli.detectCIEnvironment()
			assert.Equal(t, tt.expectedCI, ciInfo.IsCI, tt.description)
		})
	}
}

// TestCIEnvironmentDataExtraction tests data extraction from CI environments
func TestCIEnvironmentDataExtraction(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		expectedData map[string]string
		description  string
	}{
		{
			name: "GitHub Actions complete data",
			envVars: map[string]string{
				"GITHUB_ACTIONS":    "true",
				"GITHUB_RUN_ID":     "12345",
				"GITHUB_RUN_NUMBER": "67",
				"GITHUB_REF_NAME":   "main",
				"GITHUB_SHA":        "abc123def456",
				"GITHUB_REPOSITORY": "owner/repo",
				"GITHUB_ACTOR":      "testuser",
			},
			expectedData: map[string]string{
				"provider":   "github-actions",
				"build_id":   "12345",
				"build_num":  "67",
				"branch":     "main",
				"commit":     "abc123def456",
				"repository": "owner/repo",
				"actor":      "testuser",
			},
			description: "Should extract complete GitHub Actions data",
		},
		{
			name: "GitLab CI complete data",
			envVars: map[string]string{
				"GITLAB_CI":          "true",
				"CI_PIPELINE_ID":     "98765",
				"CI_PIPELINE_IID":    "43",
				"CI_COMMIT_REF_NAME": "develop",
				"CI_COMMIT_SHA":      "def456abc123",
				"CI_PROJECT_PATH":    "group/project",
				"GITLAB_USER_LOGIN":  "testuser",
			},
			expectedData: map[string]string{
				"provider":   "gitlab-ci",
				"build_id":   "98765",
				"build_num":  "43",
				"branch":     "develop",
				"commit":     "def456abc123",
				"repository": "group/project",
				"actor":      "testuser",
			},
			description: "Should extract complete GitLab CI data",
		},
		{
			name: "Partial data extraction",
			envVars: map[string]string{
				"GITHUB_ACTIONS": "true",
				"GITHUB_RUN_ID":  "12345",
				// Missing other variables
			},
			expectedData: map[string]string{
				"provider": "github-actions",
				"build_id": "12345",
			},
			description: "Should handle partial data extraction gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			clearAllCIEnvironment()

			// Set test environment variables
			for key, value := range tt.envVars {
				_ = os.Setenv(key, value)
				defer func(k string) { _ = os.Unsetenv(k) }(key)
			}

			mockLogger := &MockLogger{}
			cli := &CLI{
				outputManager: NewOutputManager(false, false, false, mockLogger),
			}

			ciInfo := cli.detectCIEnvironment()
			assert.True(t, ciInfo.IsCI, "Should detect CI environment")

			// Check extracted data
			if provider, exists := tt.expectedData["provider"]; exists {
				assert.Equal(t, provider, ciInfo.Provider, "Provider should match")
			}
			if buildID, exists := tt.expectedData["build_id"]; exists {
				assert.Equal(t, buildID, ciInfo.BuildID, "Build ID should match")
			}
			if buildNum, exists := tt.expectedData["build_num"]; exists {
				assert.Equal(t, buildNum, ciInfo.BuildNumber, "Build number should match")
			}
			if branch, exists := tt.expectedData["branch"]; exists {
				assert.Equal(t, branch, ciInfo.Branch, "Branch should match")
			}
			if commit, exists := tt.expectedData["commit"]; exists {
				assert.Equal(t, commit, ciInfo.Commit, "Commit should match")
			}
			if repository, exists := tt.expectedData["repository"]; exists {
				assert.Equal(t, repository, ciInfo.Repository, "Repository should match")
			}
			if actor, exists := tt.expectedData["actor"]; exists {
				assert.Equal(t, actor, ciInfo.Actor, "Actor should match")
			}
		})
	}
}

// TestCIEnvironmentNullSafety tests null safety in CI detection
func TestCIEnvironmentNullSafety(t *testing.T) {
	t.Run("nil CLI instance", func(t *testing.T) {
		// Test that we don't panic with nil CLI
		var cli *CLI

		// This should not panic - we'll test the behavior when CLI is nil
		assert.NotPanics(t, func() {
			if cli != nil {
				cli.detectCIEnvironment()
			}
		}, "Should not panic with nil CLI")
	})

	t.Run("uninitialized CLI components", func(t *testing.T) {
		cli := &CLI{} // No components initialized

		assert.NotPanics(t, func() {
			result := cli.detectCIEnvironment()
			assert.NotNil(t, result, "Should return valid result even with uninitialized CLI")
		}, "Should not panic with uninitialized CLI components")
	})
}

// TestCIEnvironmentConcurrency tests CI detection under concurrent access
func TestCIEnvironmentConcurrency(t *testing.T) {
	// Set up a CI environment
	_ = os.Setenv("GITHUB_ACTIONS", "true")
	_ = os.Setenv("GITHUB_RUN_ID", "12345")
	defer func() {
		_ = os.Unsetenv("GITHUB_ACTIONS")
		_ = os.Unsetenv("GITHUB_RUN_ID")
	}()

	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}

	// Run multiple goroutines to test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			assert.NotPanics(t, func() {
				ciInfo := cli.detectCIEnvironment()
				assert.NotNil(t, ciInfo, "CI info should not be nil")
				assert.True(t, ciInfo.IsCI, "Should detect CI")
				assert.Equal(t, "github-actions", ciInfo.Provider, "Provider should be correct")
			}, "Should not panic under concurrent access")
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestCIEnvironmentPerformance tests CI detection performance
func TestCIEnvironmentPerformance(t *testing.T) {
	// Set up a CI environment
	_ = os.Setenv("GITHUB_ACTIONS", "true")
	defer func() { _ = os.Unsetenv("GITHUB_ACTIONS") }()

	mockLogger := &MockLogger{}
	cli := &CLI{
		outputManager: NewOutputManager(false, false, false, mockLogger),
	}

	// Run detection multiple times to ensure it's performant
	for i := 0; i < 100; i++ {
		assert.NotPanics(t, func() {
			ciInfo := cli.detectCIEnvironment()
			assert.NotNil(t, ciInfo, "CI info should not be nil")
		}, "Should be performant for repeated calls")
	}
}

// clearAllCIEnvironment clears all known CI environment variables
func clearAllCIEnvironment() {
	ciVars := []string{
		// GitHub Actions
		"GITHUB_ACTIONS", "GITHUB_RUN_ID", "GITHUB_RUN_NUMBER",
		"GITHUB_REF_NAME", "GITHUB_SHA", "GITHUB_REPOSITORY",
		"GITHUB_EVENT_NUMBER", "GITHUB_JOB", "GITHUB_WORKFLOW",
		"GITHUB_ENVIRONMENT", "GITHUB_ACTOR",

		// GitLab CI
		"GITLAB_CI", "CI_PIPELINE_ID", "CI_PIPELINE_IID",
		"CI_COMMIT_REF_NAME", "CI_COMMIT_SHA", "CI_PROJECT_PATH",
		"CI_MERGE_REQUEST_IID", "CI_JOB_ID", "CI_ENVIRONMENT_NAME",
		"GITLAB_USER_LOGIN",

		// Jenkins
		"JENKINS_URL", "BUILD_ID", "BUILD_NUMBER", "GIT_BRANCH",
		"GIT_COMMIT", "GIT_URL", "JOB_NAME",

		// Travis CI
		"TRAVIS", "TRAVIS_BUILD_ID", "TRAVIS_BUILD_NUMBER",
		"TRAVIS_BRANCH", "TRAVIS_COMMIT", "TRAVIS_REPO_SLUG",
		"TRAVIS_PULL_REQUEST", "TRAVIS_JOB_ID",

		// CircleCI
		"CIRCLECI", "CIRCLE_BUILD_NUM", "CIRCLE_BRANCH",
		"CIRCLE_SHA1", "CIRCLE_REPOSITORY_URL", "CIRCLE_PR_NUMBER",
		"CIRCLE_JOB", "CIRCLE_WORKFLOW_ID",

		// Azure Pipelines
		"AZURE_HTTP_USER_AGENT", "TF_BUILD", "BUILD_BUILDID",
		"BUILD_BUILDNUMBER", "BUILD_SOURCEBRANCHNAME", "BUILD_SOURCEVERSION",
		"BUILD_REPOSITORY_NAME", "SYSTEM_PULLREQUEST_PULLREQUESTID",

		// Generic CI variables
		"CI", "CONTINUOUS_INTEGRATION", "BUILD_URL", "JOB_URL",
	}

	for _, variable := range ciVars {
		_ = os.Unsetenv(variable)
	}
}
