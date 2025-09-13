package utils

import (
	"strings"
	"testing"
)

func TestValidateNonEmptyString(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		wantError bool
	}{
		{"valid string", "hello", "field", false},
		{"empty string", "", "field", true},
		{"whitespace only", "   ", "field", true},
		{"string with content", "  hello  ", "field", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonEmptyString(tt.value, tt.fieldName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateNonEmptyString() error = %v, wantError %v", err, tt.wantError)
			}
			if err != nil && !strings.Contains(err.Error(), tt.fieldName) {
				t.Errorf("Error should contain field name %s, got %s", tt.fieldName, err.Error())
			}
		})
	}
}

func TestValidateNonEmptySlice(t *testing.T) {
	tests := []struct {
		name      string
		slice     interface{}
		fieldName string
		wantError bool
	}{
		{"valid slice", []string{"a", "b"}, "items", false},
		{"empty slice", []string{}, "items", true},
		{"nil slice", []string(nil), "items", true},
		{"not a slice", "not a slice", "items", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonEmptySlice(tt.slice, tt.fieldName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateNonEmptySlice() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateNotNil(t *testing.T) {
	var nilPtr *string
	validPtr := new(string)

	tests := []struct {
		name      string
		value     interface{}
		fieldName string
		wantError bool
	}{
		{"valid pointer", validPtr, "ptr", false},
		{"nil pointer", nilPtr, "ptr", true},
		{"nil interface", nil, "interface", true},
		{"valid string", "hello", "str", false},
		{"valid int", 42, "num", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotNil(tt.value, tt.fieldName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateNotNil() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
