package security

import (
	"bufio"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfirmationManager_SetDefaultAnswer(t *testing.T) {
	cm := NewConfirmationManager()

	// Test setting default answer
	cm.SetDefaultAnswer(true)
	assert.Equal(t, true, cm.defaultAnswer)

	cm.SetDefaultAnswer(false)
	assert.Equal(t, false, cm.defaultAnswer)
}

func TestConfirmationManager_SetTimeout(t *testing.T) {
	cm := NewConfirmationManager()

	// Test setting timeout
	timeout := 10 * time.Second
	cm.SetTimeout(timeout)
	assert.Equal(t, timeout, cm.timeout)
}

func TestConfirmationManager_IsNonInteractive(t *testing.T) {
	cm := NewConfirmationManager()

	// Test default value
	assert.False(t, cm.IsNonInteractive())

	// Test after setting
	cm.SetNonInteractive(true)
	assert.True(t, cm.IsNonInteractive())
}

func TestConfirmationManager_Confirm_NonInteractive(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(true)
	cm.SetDefaultAnswer(true)

	request := &ConfirmationRequest{
		Message:       "Test confirmation",
		Impact:        "safe",
		DefaultAnswer: false, // This should be overridden by manager's default
	}

	result, err := cm.Confirm(request)
	require.NoError(t, err)
	assert.True(t, result.Confirmed)
	assert.True(t, result.NonInteractive)
	assert.True(t, result.DefaultUsed)
}

func TestConfirmationManager_ConfirmDirectoryDelete(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(true)
	cm.SetDefaultAnswer(false) // Default to false for destructive operations

	result, err := cm.ConfirmDirectoryDelete("/test/dir", 10, 1024)
	require.NoError(t, err)
	assert.False(t, result.Confirmed) // Should use default (false)
	assert.True(t, result.NonInteractive)
	assert.True(t, result.DefaultUsed)
}

func TestConfirmationManager_ConfirmBulkOperation(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(true)

	tests := []struct {
		name           string
		operationType  string
		itemCount      int
		defaultAnswer  bool
		expectedImpact string
	}{
		{
			name:           "small safe operation",
			operationType:  "update",
			itemCount:      5,
			defaultAnswer:  true,
			expectedImpact: "warning",
		},
		{
			name:           "large operation",
			operationType:  "process",
			itemCount:      150,
			defaultAnswer:  false,
			expectedImpact: "destructive",
		},
		{
			name:           "delete operation",
			operationType:  "delete files",
			itemCount:      10,
			defaultAnswer:  false,
			expectedImpact: "destructive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm.SetDefaultAnswer(tt.defaultAnswer)

			details := []string{"Detail 1", "Detail 2"}
			result, err := cm.ConfirmBulkOperation(tt.operationType, tt.itemCount, details)

			require.NoError(t, err)
			assert.Equal(t, tt.defaultAnswer, result.Confirmed)
			assert.True(t, result.NonInteractive)
			assert.True(t, result.DefaultUsed)
		})
	}
}

func TestConfirmationManager_ConfirmSecurityRisk(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(true)
	cm.SetDefaultAnswer(false) // Always default to false for security risks

	riskDescription := "Potential security vulnerability detected"
	riskLevel := "high"
	details := []string{"CVE-2023-1234", "Affects authentication"}

	result, err := cm.ConfirmSecurityRisk(riskDescription, riskLevel, details)
	require.NoError(t, err)
	assert.False(t, result.Confirmed) // Should default to false for security risks
	assert.True(t, result.NonInteractive)
	assert.True(t, result.DefaultUsed)
}

func TestConfirmationManager_ConfirmWithDryRun(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(true)

	tests := []struct {
		name           string
		dryRunSummary  map[string]interface{}
		defaultAnswer  bool
		expectedResult bool
	}{
		{
			name: "safe operations only",
			dryRunSummary: map[string]interface{}{
				"total_operations":       10,
				"safe_operations":        10,
				"warning_operations":     0,
				"destructive_operations": 0,
				"total_files_affected":   5,
			},
			defaultAnswer:  true,
			expectedResult: true,
		},
		{
			name: "with warning operations",
			dryRunSummary: map[string]interface{}{
				"total_operations":       10,
				"safe_operations":        8,
				"warning_operations":     2,
				"destructive_operations": 0,
				"total_files_affected":   5,
			},
			defaultAnswer:  false,
			expectedResult: false,
		},
		{
			name: "with destructive operations",
			dryRunSummary: map[string]interface{}{
				"total_operations":       10,
				"safe_operations":        7,
				"warning_operations":     2,
				"destructive_operations": 1,
				"total_files_affected":   5,
			},
			defaultAnswer:  false,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm.SetDefaultAnswer(tt.defaultAnswer)

			result, err := cm.ConfirmWithDryRun(tt.dryRunSummary)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result.Confirmed)
			assert.True(t, result.NonInteractive)
		})
	}
}

func TestConfirmationManager_Confirm_Interactive_RequireExplicit(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(false)

	// Mock reader with "yes" input
	cm.reader = bufio.NewReader(strings.NewReader("yes\n"))

	request := &ConfirmationRequest{
		Message:         "Destructive operation",
		Impact:          "destructive",
		DefaultAnswer:   false,
		RequireExplicit: true,
	}

	result, err := cm.Confirm(request)
	require.NoError(t, err)
	assert.True(t, result.Confirmed)
	assert.False(t, result.NonInteractive)
	assert.False(t, result.DefaultUsed)
	assert.Equal(t, "yes", result.UserInput)
}

func TestConfirmationManager_Confirm_Interactive_RequireExplicit_No(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(false)

	// Mock reader with "no" input
	cm.reader = bufio.NewReader(strings.NewReader("no\n"))

	request := &ConfirmationRequest{
		Message:         "Destructive operation",
		Impact:          "destructive",
		DefaultAnswer:   false,
		RequireExplicit: true,
	}

	result, err := cm.Confirm(request)
	require.NoError(t, err)
	assert.False(t, result.Confirmed)
	assert.False(t, result.NonInteractive)
	assert.False(t, result.DefaultUsed)
	assert.Equal(t, "no", result.UserInput)
}

func TestConfirmationManager_Confirm_Interactive_SimpleYes(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(false)

	// Mock reader with "y" input
	cm.reader = bufio.NewReader(strings.NewReader("y\n"))

	request := &ConfirmationRequest{
		Message:         "Simple confirmation",
		Impact:          "safe",
		DefaultAnswer:   false,
		RequireExplicit: false,
	}

	result, err := cm.Confirm(request)
	require.NoError(t, err)
	assert.True(t, result.Confirmed)
	assert.False(t, result.NonInteractive)
	assert.False(t, result.DefaultUsed)
	assert.Equal(t, "y", result.UserInput)
}

func TestConfirmationManager_Confirm_Interactive_EmptyInput(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(false)

	// Mock reader with empty input (just Enter)
	cm.reader = bufio.NewReader(strings.NewReader("\n"))

	request := &ConfirmationRequest{
		Message:         "Simple confirmation",
		Impact:          "safe",
		DefaultAnswer:   true,
		RequireExplicit: false,
	}

	result, err := cm.Confirm(request)
	require.NoError(t, err)
	assert.True(t, result.Confirmed) // Should use default answer
	assert.False(t, result.NonInteractive)
	assert.True(t, result.DefaultUsed)
	assert.Equal(t, "", result.UserInput)
}

func TestConfirmationManager_Confirm_Interactive_InvalidInput(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(false)

	// Mock reader with invalid input
	cm.reader = bufio.NewReader(strings.NewReader("invalid\n"))

	request := &ConfirmationRequest{
		Message:         "Simple confirmation",
		Impact:          "safe",
		DefaultAnswer:   false,
		RequireExplicit: false,
	}

	result, err := cm.Confirm(request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input")
	assert.Equal(t, "invalid", result.UserInput)
}

func TestConfirmationManager_Confirm_WithTimeout(t *testing.T) {
	cm := NewConfirmationManager()
	cm.SetNonInteractive(false)
	cm.SetTimeout(10 * time.Millisecond) // Short timeout

	// Create a pipe that will block on read
	r, w := io.Pipe()
	cm.reader = bufio.NewReader(r)

	// Close the writer after the test to clean up
	defer func() { _ = w.Close() }()

	request := &ConfirmationRequest{
		Message:       "Timeout test",
		Impact:        "safe",
		DefaultAnswer: true,
	}

	result, err := cm.Confirm(request)
	require.NoError(t, err)
	assert.True(t, result.Confirmed) // Should use default on timeout
	assert.True(t, result.TimedOut)
	assert.True(t, result.DefaultUsed)
}

func TestConfirmationHistory_Record(t *testing.T) {
	ch := NewConfirmationHistory()

	request := &ConfirmationRequest{
		Message: "Test request",
		Impact:  "safe",
	}

	result := &ConfirmationResult{
		Confirmed: true,
		Timestamp: time.Now(),
	}

	ch.Record(request, result)

	history := ch.GetHistory()
	assert.Len(t, history, 1)
	assert.Equal(t, request.Message, history[0].Request.Message)
	assert.Equal(t, result.Confirmed, history[0].Result.Confirmed)
}

func TestConfirmationHistory_GetRecentHistory(t *testing.T) {
	ch := NewConfirmationHistory()

	// Add multiple entries
	for i := 0; i < 5; i++ {
		request := &ConfirmationRequest{
			Message: "Test request",
			Impact:  "safe",
		}

		result := &ConfirmationResult{
			Confirmed: true,
			Timestamp: time.Now(),
		}

		ch.Record(request, result)
	}

	// Test getting recent history
	recent := ch.GetRecentHistory(3)
	assert.Len(t, recent, 3)

	// Test getting more than available
	recent = ch.GetRecentHistory(10)
	assert.Len(t, recent, 5)

	// Test getting zero or negative
	recent = ch.GetRecentHistory(0)
	assert.Len(t, recent, 5)

	recent = ch.GetRecentHistory(-1)
	assert.Len(t, recent, 5)
}

func TestConfirmationHistory_Clear(t *testing.T) {
	ch := NewConfirmationHistory()

	// Add entries
	request := &ConfirmationRequest{Message: "Test"}
	result := &ConfirmationResult{Confirmed: true}
	ch.Record(request, result)

	assert.Len(t, ch.GetHistory(), 1)

	// Clear history
	ch.Clear()
	assert.Len(t, ch.GetHistory(), 0)
}

func TestConfirmationManager_Integration(t *testing.T) {
	cm := NewConfirmationManager()
	ch := NewConfirmationHistory()

	// Test complete workflow in non-interactive mode
	cm.SetNonInteractive(true)
	cm.SetDefaultAnswer(true)

	// Test file overwrite confirmation
	result, err := cm.ConfirmFileOverwrite("/test/file.txt", 1024)
	require.NoError(t, err)
	assert.True(t, result.Confirmed) // Should be true due to default answer
	ch.Record(&ConfirmationRequest{Message: "File overwrite"}, result)

	// Test directory delete confirmation
	result, err = cm.ConfirmDirectoryDelete("/test/dir", 5, 2048)
	require.NoError(t, err)
	assert.True(t, result.Confirmed) // Should be true due to default answer
	ch.Record(&ConfirmationRequest{Message: "Directory delete"}, result)

	// Test bulk operation confirmation
	result, err = cm.ConfirmBulkOperation("process", 50, []string{"detail1", "detail2"})
	require.NoError(t, err)
	assert.True(t, result.Confirmed)
	ch.Record(&ConfirmationRequest{Message: "Bulk operation"}, result)

	// Verify history
	history := ch.GetHistory()
	assert.Len(t, history, 3)

	// Test switching to interactive mode
	cm.SetNonInteractive(false)
	cm.reader = bufio.NewReader(strings.NewReader("yes\n"))

	result, err = cm.ConfirmSecurityRisk("Test risk", "medium", []string{"detail"})
	require.NoError(t, err)
	assert.True(t, result.Confirmed)
	assert.False(t, result.NonInteractive)
}

func TestConfirmationManager_EdgeCases(t *testing.T) {
	cm := NewConfirmationManager()

	// Test with nil request (should not panic)
	request := &ConfirmationRequest{}
	cm.SetNonInteractive(true)

	result, err := cm.Confirm(request)
	require.NoError(t, err)
	assert.False(t, result.Confirmed) // Default should be false

	// Test with empty details and metadata
	request = &ConfirmationRequest{
		Message:  "Test",
		Impact:   "safe",
		Details:  []string{},
		Metadata: map[string]interface{}{},
	}

	result, err = cm.Confirm(request)
	require.NoError(t, err)
	assert.False(t, result.Confirmed)
}
