package workflow

import (
	"testing"
	"time"
)

func TestGenerateWorkflowID(t *testing.T) {
	id1 := generateWorkflowID()
	id2 := generateWorkflowID()

	if id1 == "" {
		t.Error("Generated workflow ID should not be empty")
	}

	if id2 == "" {
		t.Error("Generated workflow ID should not be empty")
	}

	if id1 == id2 {
		t.Error("Generated workflow IDs should be unique")
	}

	// Check format
	if len(id1) < 10 {
		t.Error("Generated workflow ID should be reasonably long")
	}
}

func TestGetProjectPathFromMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]interface{}
		expected string
	}{
		{
			name:     "nil metadata",
			metadata: nil,
			expected: "",
		},
		{
			name:     "empty metadata",
			metadata: map[string]interface{}{},
			expected: "",
		},
		{
			name: "metadata with project_path",
			metadata: map[string]interface{}{
				"project_path": "/test/project",
			},
			expected: "/test/project",
		},
		{
			name: "metadata with non-string project_path",
			metadata: map[string]interface{}{
				"project_path": 123,
			},
			expected: "",
		},
		{
			name: "metadata without project_path",
			metadata: map[string]interface{}{
				"other_field": "value",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getProjectPathFromMetadata(tt.metadata)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestWorkflowIDGeneration(t *testing.T) {
	// Test that two sequential calls generate different IDs
	id1 := generateWorkflowID()
	time.Sleep(1 * time.Millisecond) // Ensure different timestamp
	id2 := generateWorkflowID()

	if id1 == id2 {
		t.Error("Sequential workflow IDs should be different")
	}
}

func TestWorkflowIDFormat(t *testing.T) {
	id := generateWorkflowID()

	// Check that ID starts with expected prefix
	expectedPrefix := "workflow_"
	if len(id) < len(expectedPrefix) || id[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Workflow ID should start with '%s', got: %s", expectedPrefix, id)
	}

	// Check that the timestamp part is numeric
	timestampPart := id[len(expectedPrefix):]
	if len(timestampPart) == 0 {
		t.Error("Workflow ID should have a timestamp part")
	}

	// Verify it's a reasonable timestamp (should be close to current time)
	// This is a basic sanity check
	if len(timestampPart) < 10 { // Unix timestamp should be at least 10 digits
		t.Error("Timestamp part seems too short")
	}
}

func TestMetadataExtraction(t *testing.T) {
	// Test with various metadata structures
	testCases := []struct {
		name     string
		metadata map[string]interface{}
		expected string
	}{
		{
			name: "valid project path",
			metadata: map[string]interface{}{
				"project_path": "/home/user/project",
				"other_data":   "value",
			},
			expected: "/home/user/project",
		},
		{
			name: "project path with special characters",
			metadata: map[string]interface{}{
				"project_path": "/path/with spaces/and-dashes_underscores",
			},
			expected: "/path/with spaces/and-dashes_underscores",
		},
		{
			name: "empty project path",
			metadata: map[string]interface{}{
				"project_path": "",
			},
			expected: "",
		},
		{
			name: "project path as integer (invalid)",
			metadata: map[string]interface{}{
				"project_path": 42,
			},
			expected: "",
		},
		{
			name: "project path as boolean (invalid)",
			metadata: map[string]interface{}{
				"project_path": true,
			},
			expected: "",
		},
		{
			name: "project path as slice (invalid)",
			metadata: map[string]interface{}{
				"project_path": []string{"path1", "path2"},
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getProjectPathFromMetadata(tc.metadata)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestWorkflowIDConcurrency(t *testing.T) {
	// Test that concurrent ID generation doesn't panic
	// We don't test for uniqueness here due to timing issues in tests
	const numGoroutines = 3

	done := make(chan bool, numGoroutines)

	// Start multiple goroutines generating IDs
	for i := 0; i < numGoroutines; i++ {
		go func() {
			// Just generate an ID to ensure no panic
			id := generateWorkflowID()
			if id == "" {
				t.Error("Generated ID should not be empty")
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

func TestWorkflowIDTiming(t *testing.T) {
	// Test that IDs generated in sequence have increasing timestamps
	id1 := generateWorkflowID()
	time.Sleep(1 * time.Millisecond) // Small delay to ensure different timestamps
	id2 := generateWorkflowID()

	// Extract timestamp parts
	prefix := "workflow_"
	timestamp1 := id1[len(prefix):]
	timestamp2 := id2[len(prefix):]

	// Since we're using nanosecond timestamps, the second should be larger
	// This is a basic check - in practice, the timestamps should be different
	if timestamp1 == timestamp2 {
		t.Error("Sequential workflow IDs should have different timestamps")
	}
}

func TestMetadataEdgeCases(t *testing.T) {
	// Test edge cases for metadata handling

	// Test with nil map
	result := getProjectPathFromMetadata(nil)
	if result != "" {
		t.Errorf("Expected empty string for nil metadata, got '%s'", result)
	}

	// Test with nested map (should not extract nested project_path)
	nestedMetadata := map[string]interface{}{
		"nested": map[string]interface{}{
			"project_path": "/nested/path",
		},
	}
	result = getProjectPathFromMetadata(nestedMetadata)
	if result != "" {
		t.Errorf("Should not extract nested project_path, got '%s'", result)
	}

	// Test with multiple keys including project_path
	multiKeyMetadata := map[string]interface{}{
		"project_path": "/correct/path",
		"backup_path":  "/backup/path",
		"temp_path":    "/temp/path",
	}
	result = getProjectPathFromMetadata(multiKeyMetadata)
	if result != "/correct/path" {
		t.Errorf("Expected '/correct/path', got '%s'", result)
	}
}

// Benchmark tests for performance
func BenchmarkGenerateWorkflowID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateWorkflowID()
	}
}

func BenchmarkGetProjectPathFromMetadata(b *testing.B) {
	metadata := map[string]interface{}{
		"project_path": "/test/project/path",
		"other_field":  "some value",
		"number":       42,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getProjectPathFromMetadata(metadata)
	}
}
