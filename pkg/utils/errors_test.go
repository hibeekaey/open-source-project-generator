package utils

import (
	"errors"
	"strings"
	"testing"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		context string
		want    string
	}{
		{"nil error", nil, "test context", ""},
		{"with error", errors.New("original error"), "test context", "test context: original error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HandleError(tt.err, tt.context)
			if tt.want == "" {
				if result != nil {
					t.Errorf("HandleError() = %v, want nil", result)
				}
			} else {
				if result == nil || result.Error() != tt.want {
					t.Errorf("HandleError() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")

	tests := []struct {
		name   string
		err    error
		format string
		args   []interface{}
		want   string
	}{
		{"nil error", nil, "context %s", []interface{}{"test"}, ""},
		{"with error", originalErr, "context %s", []interface{}{"test"}, "context test: original error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.format, tt.args...)
			if tt.want == "" {
				if result != nil {
					t.Errorf("WrapError() = %v, want nil", result)
				}
			} else {
				if result == nil || result.Error() != tt.want {
					t.Errorf("WrapError() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestIsNilError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		operation string
		wantError bool
		wantMsg   string
	}{
		{"nil error", nil, "test operation", false, ""},
		{"with error", errors.New("test error"), "test operation", true, "failed to test operation: test error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNilError(tt.err, tt.operation)
			if (result != nil) != tt.wantError {
				t.Errorf("IsNilError() error = %v, wantError %v", result, tt.wantError)
			}
			if result != nil && !strings.Contains(result.Error(), tt.wantMsg) {
				t.Errorf("IsNilError() = %v, want to contain %v", result.Error(), tt.wantMsg)
			}
		})
	}
}
