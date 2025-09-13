package common

import (
	"errors"
	"testing"
)

// Mock vulnerability database for testing
type mockVulnerabilityDB struct {
	issues []SecurityIssue
	err    error
}

func (m *mockVulnerabilityDB) CheckVulnerabilities(packageName, version string) ([]SecurityIssue, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.issues, nil
}

func TestCheckSecurity(t *testing.T) {
	tests := []struct {
		name    string
		pkg     string
		version string
		db      VulnerabilityDB
		wantErr bool
	}{
		{
			name:    "successful check",
			pkg:     "test-package",
			version: "1.0.0",
			db:      &mockVulnerabilityDB{issues: []SecurityIssue{{ID: "CVE-2023-1234", Severity: "high"}}},
			wantErr: false,
		},
		{
			name:    "nil database",
			pkg:     "test-package",
			version: "1.0.0",
			db:      nil,
			wantErr: true,
		},
		{
			name:    "database error",
			pkg:     "test-package",
			version: "1.0.0",
			db:      &mockVulnerabilityDB{err: errors.New("db error")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CheckSecurity(tt.pkg, tt.version, tt.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSecurity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilterSecurityIssues(t *testing.T) {
	issues := []SecurityIssue{
		{ID: "1", Severity: "low"},
		{ID: "2", Severity: "moderate"},
		{ID: "3", Severity: "high"},
		{ID: "4", Severity: "critical"},
	}

	tests := []struct {
		name        string
		minSeverity string
		wantCount   int
	}{
		{"filter low", "low", 4},
		{"filter moderate", "moderate", 3},
		{"filter high", "high", 2},
		{"filter critical", "critical", 1},
		{"invalid severity", "invalid", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := FilterSecurityIssues(issues, tt.minSeverity)
			if len(filtered) != tt.wantCount {
				t.Errorf("FilterSecurityIssues() count = %v, want %v", len(filtered), tt.wantCount)
			}
		})
	}
}

func TestHasCriticalIssues(t *testing.T) {
	tests := []struct {
		name   string
		issues []SecurityIssue
		want   bool
	}{
		{
			name:   "has critical",
			issues: []SecurityIssue{{Severity: "critical"}, {Severity: "high"}},
			want:   true,
		},
		{
			name:   "no critical",
			issues: []SecurityIssue{{Severity: "high"}, {Severity: "moderate"}},
			want:   false,
		},
		{
			name:   "empty issues",
			issues: []SecurityIssue{},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasCriticalIssues(tt.issues)
			if got != tt.want {
				t.Errorf("HasCriticalIssues() = %v, want %v", got, tt.want)
			}
		})
	}
}
