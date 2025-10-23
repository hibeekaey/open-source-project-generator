package interactive

import (
	"context"
	"errors"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/logger"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockPrompter implements Prompter interface for testing
type MockPrompter struct {
	inputs        []string
	inputIndex    int
	selections    []string
	selectIndex   int
	multiSelects  [][]string
	multiIndex    int
	confirms      []bool
	confirmIndex  int
	passwords     []string
	passwordIndex int
	shouldError   bool
}

func NewMockPrompter() *MockPrompter {
	return &MockPrompter{
		inputs:       []string{},
		selections:   []string{},
		multiSelects: [][]string{},
		confirms:     []bool{},
		passwords:    []string{},
	}
}

func (m *MockPrompter) WithInputs(inputs ...string) *MockPrompter {
	m.inputs = inputs
	return m
}

func (m *MockPrompter) WithSelections(selections ...string) *MockPrompter {
	m.selections = selections
	return m
}

func (m *MockPrompter) WithMultiSelects(multiSelects ...[]string) *MockPrompter {
	m.multiSelects = multiSelects
	return m
}

func (m *MockPrompter) WithConfirms(confirms ...bool) *MockPrompter {
	m.confirms = confirms
	return m
}

func (m *MockPrompter) WithPasswords(passwords ...string) *MockPrompter {
	m.passwords = passwords
	return m
}

func (m *MockPrompter) WithError() *MockPrompter {
	m.shouldError = true
	return m
}

func (m *MockPrompter) Input(message string, defaultValue string) (string, error) {
	if m.shouldError {
		return "", errors.New("mock error")
	}
	if m.inputIndex >= len(m.inputs) {
		return defaultValue, nil
	}
	result := m.inputs[m.inputIndex]
	m.inputIndex++
	return result, nil
}

func (m *MockPrompter) Select(message string, options []string) (string, error) {
	if m.shouldError {
		return "", errors.New("mock error")
	}
	if m.selectIndex >= len(m.selections) {
		return options[0], nil
	}
	result := m.selections[m.selectIndex]
	m.selectIndex++
	return result, nil
}

func (m *MockPrompter) MultiSelect(message string, options []string) ([]string, error) {
	if m.shouldError {
		return nil, errors.New("mock error")
	}
	if m.multiIndex >= len(m.multiSelects) {
		return options[:1], nil
	}
	result := m.multiSelects[m.multiIndex]
	m.multiIndex++
	return result, nil
}

func (m *MockPrompter) Confirm(message string, defaultValue bool) (bool, error) {
	if m.shouldError {
		return false, errors.New("mock error")
	}
	if m.confirmIndex >= len(m.confirms) {
		return defaultValue, nil
	}
	result := m.confirms[m.confirmIndex]
	m.confirmIndex++
	return result, nil
}

func (m *MockPrompter) Password(message string) (string, error) {
	if m.shouldError {
		return "", errors.New("mock error")
	}
	if m.passwordIndex >= len(m.passwords) {
		return "", nil
	}
	result := m.passwords[m.passwordIndex]
	m.passwordIndex++
	return result, nil
}

// TestInteractiveWizard_CompleteFlow tests the complete wizard flow
func TestInteractiveWizard_CompleteFlow(t *testing.T) {
	log := logger.NewLogger()
	wizard := NewInteractiveWizard(log)

	// Setup mock prompter with complete flow
	mockPrompter := NewMockPrompter().
		WithInputs(
			"test-project",     // project name
			"Test Description", // project description
			"./test-output",    // output directory
			"web-app",          // nextjs component name
		).
		WithMultiSelects(
			[]string{"nextjs - Next.js frontend with TypeScript and Tailwind"}, // component selection
		).
		WithConfirms(
			true, // typescript
			true, // tailwind
			true, // app_router
			true, // eslint
			true, // generate docker compose
			true, // generate scripts
			true, // final confirmation
		)

	wizard.prompter = mockPrompter

	// Run wizard
	config, err := wizard.Run(context.Background())

	// Verify results
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "test-project", config.Name)
	assert.Equal(t, "Test Description", config.Description)
	assert.Contains(t, config.OutputDir, "test-output")
	assert.Len(t, config.Components, 1)
	assert.Equal(t, "nextjs", config.Components[0].Type)
	assert.True(t, config.Integration.GenerateDockerCompose)
	assert.True(t, config.Integration.GenerateScripts)
}

// TestInteractiveWizard_UserCancellation tests cancellation at various steps
func TestInteractiveWizard_UserCancellation(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockPrompter) *MockPrompter
		expectError bool
	}{
		{
			name: "cancel at final confirmation",
			setupMock: func(m *MockPrompter) *MockPrompter {
				return m.
					WithInputs("test-project", "Test Description", "./test-output", "web-app").
					WithMultiSelects([]string{"nextjs - Next.js frontend with TypeScript and Tailwind"}).
					WithConfirms(true, true, true, true, true, true, false) // false at final confirmation
			},
			expectError: true,
		},
		{
			name: "error during project info collection",
			setupMock: func(m *MockPrompter) *MockPrompter {
				return m.WithError()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.NewLogger()
			wizard := NewInteractiveWizard(log)
			wizard.prompter = tt.setupMock(NewMockPrompter())

			config, err := wizard.Run(context.Background())

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, config)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
			}
		})
	}
}

// TestInteractiveWizard_InvalidInput tests invalid input handling
func TestInteractiveWizard_InvalidInput(t *testing.T) {
	log := logger.NewLogger()
	wizard := NewInteractiveWizard(log)

	tests := []struct {
		name        string
		projectName string
		expectError bool
	}{
		{
			name:        "valid project name",
			projectName: "my-project",
			expectError: false,
		},
		{
			name:        "project name with spaces gets sanitized",
			projectName: "my project",
			expectError: false,
		},
		{
			name:        "empty project name uses default",
			projectName: "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPrompter := NewMockPrompter().
				WithInputs(tt.projectName, "Test Description", "./test-output")

			wizard.prompter = mockPrompter

			projectInfo, err := wizard.CollectProjectInfo()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, projectInfo)
				assert.NotEmpty(t, projectInfo.Name)
			}
		})
	}
}

// TestInteractiveWizard_MultipleComponents tests configuration with multiple components
func TestInteractiveWizard_MultipleComponents(t *testing.T) {
	log := logger.NewLogger()
	wizard := NewInteractiveWizard(log)

	mockPrompter := NewMockPrompter().
		WithInputs(
			"fullstack-app",
			"Fullstack Application",
			"./fullstack-output",
			"web-app",         // nextjs name
			"api-server",      // go-backend name
			"example.com/api", // go module
		).
		WithMultiSelects(
			[]string{
				"nextjs - Next.js frontend with TypeScript and Tailwind",
				"go-backend - Go backend with Gin framework",
			},
		).
		WithSelections(
			"gin", // framework selection
		).
		WithConfirms(
			true, // typescript
			true, // tailwind
			true, // app_router
			true, // eslint
			true, // generate docker compose
			true, // generate scripts
			true, // final confirmation
		)

	wizard.prompter = mockPrompter

	config, err := wizard.Run(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Len(t, config.Components, 2)
	assert.Equal(t, "nextjs", config.Components[0].Type)
	assert.Equal(t, "go-backend", config.Components[1].Type)
}

// TestInteractiveWizard_ComponentConfiguration tests component-specific configuration
func TestInteractiveWizard_ComponentConfiguration(t *testing.T) {
	log := logger.NewLogger()
	wizard := NewInteractiveWizard(log)

	tests := []struct {
		name          string
		componentType string
		setupMock     func(*MockPrompter) *MockPrompter
		validate      func(*testing.T, map[string]interface{})
	}{
		{
			name:          "nextjs configuration",
			componentType: "nextjs",
			setupMock: func(m *MockPrompter) *MockPrompter {
				return m.
					WithInputs("web-app").
					WithConfirms(true, true, true, true)
			},
			validate: func(t *testing.T, config map[string]interface{}) {
				assert.Equal(t, true, config["typescript"])
				assert.Equal(t, true, config["tailwind"])
				assert.Equal(t, true, config["app_router"])
				assert.Equal(t, true, config["eslint"])
			},
		},
		{
			name:          "go-backend configuration",
			componentType: "go-backend",
			setupMock: func(m *MockPrompter) *MockPrompter {
				return m.
					WithInputs("api-server", "example.com/api", "8080").
					WithSelections("gin")
			},
			validate: func(t *testing.T, config map[string]interface{}) {
				assert.Equal(t, "example.com/api", config["module"])
				assert.Equal(t, "gin", config["framework"])
				assert.Equal(t, 8080, config["port"])
			},
		},
		{
			name:          "android configuration",
			componentType: "android",
			setupMock: func(m *MockPrompter) *MockPrompter {
				return m.
					WithInputs("mobile-app", "com.example.app", "21", "33").
					WithSelections("kotlin")
			},
			validate: func(t *testing.T, config map[string]interface{}) {
				assert.Equal(t, "com.example.app", config["package"])
				assert.Equal(t, 21, config["min_sdk"])
				assert.Equal(t, 33, config["target_sdk"])
				assert.Equal(t, "kotlin", config["language"])
			},
		},
		{
			name:          "ios configuration",
			componentType: "ios",
			setupMock: func(m *MockPrompter) *MockPrompter {
				return m.
					WithInputs("mobile-app", "com.example.app", "14.0").
					WithSelections("swift")
			},
			validate: func(t *testing.T, config map[string]interface{}) {
				assert.Equal(t, "com.example.app", config["bundle_id"])
				assert.Equal(t, "14.0", config["deployment_target"])
				assert.Equal(t, "swift", config["language"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wizard.prompter = tt.setupMock(NewMockPrompter())

			config, err := wizard.ConfigureComponent(tt.componentType)

			require.NoError(t, err)
			assert.Equal(t, tt.componentType, config.Type)
			assert.True(t, config.Enabled)
			if tt.validate != nil {
				tt.validate(t, config.Config)
			}
		})
	}
}

// TestInteractiveWizard_CollectProjectInfo tests project info collection
func TestInteractiveWizard_CollectProjectInfo(t *testing.T) {
	log := logger.NewLogger()
	wizard := NewInteractiveWizard(log)

	mockPrompter := NewMockPrompter().
		WithInputs("my-awesome-project", "An awesome project", "./output")

	wizard.prompter = mockPrompter

	projectInfo, err := wizard.CollectProjectInfo()

	require.NoError(t, err)
	assert.NotNil(t, projectInfo)
	assert.Equal(t, "my-awesome-project", projectInfo.Name)
	assert.Equal(t, "An awesome project", projectInfo.Description)
	assert.Contains(t, projectInfo.OutputDir, "output")
}

// TestInteractiveWizard_SelectComponents tests component selection
func TestInteractiveWizard_SelectComponents(t *testing.T) {
	log := logger.NewLogger()
	wizard := NewInteractiveWizard(log)

	tests := []struct {
		name          string
		selections    []string
		expectedTypes []string
		expectedCount int
	}{
		{
			name: "single component",
			selections: []string{
				"nextjs - Next.js frontend with TypeScript and Tailwind",
			},
			expectedTypes: []string{"nextjs"},
			expectedCount: 1,
		},
		{
			name: "multiple components",
			selections: []string{
				"nextjs - Next.js frontend with TypeScript and Tailwind",
				"go-backend - Go backend with Gin framework",
			},
			expectedTypes: []string{"nextjs", "go-backend"},
			expectedCount: 2,
		},
		{
			name: "all components",
			selections: []string{
				"nextjs - Next.js frontend with TypeScript and Tailwind",
				"go-backend - Go backend with Gin framework",
				"android - Android mobile app with Kotlin",
				"ios - iOS mobile app with Swift",
			},
			expectedTypes: []string{"nextjs", "go-backend", "android", "ios"},
			expectedCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPrompter := NewMockPrompter().
				WithMultiSelects(tt.selections)

			wizard.prompter = mockPrompter

			componentTypes, err := wizard.SelectComponents()

			require.NoError(t, err)
			assert.Len(t, componentTypes, tt.expectedCount)
			for i, expectedType := range tt.expectedTypes {
				assert.Equal(t, expectedType, componentTypes[i])
			}
		})
	}
}

// TestInteractiveWizard_ConfirmConfiguration tests configuration confirmation
func TestInteractiveWizard_ConfirmConfiguration(t *testing.T) {
	log := logger.NewLogger()
	wizard := NewInteractiveWizard(log)

	tests := []struct {
		name       string
		confirmed  bool
		expectTrue bool
	}{
		{
			name:       "user confirms",
			confirmed:  true,
			expectTrue: true,
		},
		{
			name:       "user declines",
			confirmed:  false,
			expectTrue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPrompter := NewMockPrompter().
				WithConfirms(tt.confirmed)

			wizard.prompter = mockPrompter

			// Create a minimal config for testing
			config := &models.ProjectConfig{
				Name:        "test-project",
				Description: "Test",
				OutputDir:   "./test",
				Components:  []models.ComponentConfig{},
			}

			confirmed, err := wizard.ConfirmConfiguration(config)

			require.NoError(t, err)
			assert.Equal(t, tt.expectTrue, confirmed)
		})
	}
}
