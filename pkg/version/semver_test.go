package version

import (
	"testing"
)

func TestParseSemVer(t *testing.T) {
	tests := []struct {
		input       string
		expected    *SemVer
		expectError bool
	}{
		{
			input: "1.2.3",
			expected: &SemVer{
				Major: 1, Minor: 2, Patch: 3,
				Original: "1.2.3",
			},
			expectError: false,
		},
		{
			input: "v1.2.3",
			expected: &SemVer{
				Major: 1, Minor: 2, Patch: 3,
				Original: "v1.2.3",
			},
			expectError: false,
		},
		{
			input: "1.2.3-alpha.1",
			expected: &SemVer{
				Major: 1, Minor: 2, Patch: 3,
				Prerelease: "alpha.1",
				Original:   "1.2.3-alpha.1",
			},
			expectError: false,
		},
		{
			input: "1.2.3+build.1",
			expected: &SemVer{
				Major: 1, Minor: 2, Patch: 3,
				Build:    "build.1",
				Original: "1.2.3+build.1",
			},
			expectError: false,
		},
		{
			input: "1.2.3-alpha.1+build.1",
			expected: &SemVer{
				Major: 1, Minor: 2, Patch: 3,
				Prerelease: "alpha.1",
				Build:      "build.1",
				Original:   "1.2.3-alpha.1+build.1",
			},
			expectError: false,
		},
		{
			input:       "invalid",
			expectError: true,
		},
		{
			input:       "1.2",
			expectError: true,
		},
		{
			input:       "1.2.3.4",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseSemVer(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", tt.input, err)
				return
			}

			if result.Major != tt.expected.Major {
				t.Errorf("Expected major %d, got %d", tt.expected.Major, result.Major)
			}
			if result.Minor != tt.expected.Minor {
				t.Errorf("Expected minor %d, got %d", tt.expected.Minor, result.Minor)
			}
			if result.Patch != tt.expected.Patch {
				t.Errorf("Expected patch %d, got %d", tt.expected.Patch, result.Patch)
			}
			if result.Prerelease != tt.expected.Prerelease {
				t.Errorf("Expected prerelease %s, got %s", tt.expected.Prerelease, result.Prerelease)
			}
			if result.Build != tt.expected.Build {
				t.Errorf("Expected build %s, got %s", tt.expected.Build, result.Build)
			}
			if result.Original != tt.expected.Original {
				t.Errorf("Expected original %s, got %s", tt.expected.Original, result.Original)
			}
		})
	}
}

func TestSemVerString(t *testing.T) {
	tests := []struct {
		semver   *SemVer
		expected string
	}{
		{
			semver:   &SemVer{Major: 1, Minor: 2, Patch: 3},
			expected: "1.2.3",
		},
		{
			semver:   &SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "alpha.1"},
			expected: "1.2.3-alpha.1",
		},
		{
			semver:   &SemVer{Major: 1, Minor: 2, Patch: 3, Build: "build.1"},
			expected: "1.2.3+build.1",
		},
		{
			semver:   &SemVer{Major: 1, Minor: 2, Patch: 3, Prerelease: "alpha.1", Build: "build.1"},
			expected: "1.2.3-alpha.1+build.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.semver.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestSemVerCompare(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.1.0", "1.0.0", 1},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0-alpha", "1.0.0", -1},
		{"1.0.0", "1.0.0-alpha", 1},
		{"1.0.0-alpha", "1.0.0-beta", -1},
		{"1.0.0-alpha.1", "1.0.0-alpha.2", -1},
		{"1.0.0-alpha.10", "1.0.0-alpha.2", 1},
		{"1.0.0-alpha", "1.0.0-alpha.1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.v1+" vs "+tt.v2, func(t *testing.T) {
			v1, err := ParseSemVer(tt.v1)
			if err != nil {
				t.Fatalf("Failed to parse v1 %s: %v", tt.v1, err)
			}

			v2, err := ParseSemVer(tt.v2)
			if err != nil {
				t.Fatalf("Failed to parse v2 %s: %v", tt.v2, err)
			}

			result := v1.Compare(v2)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}

			// Test convenience methods
			if tt.expected > 0 && !v1.IsGreaterThan(v2) {
				t.Errorf("Expected %s > %s", tt.v1, tt.v2)
			}
			if tt.expected < 0 && !v1.IsLessThan(v2) {
				t.Errorf("Expected %s < %s", tt.v1, tt.v2)
			}
			if tt.expected == 0 && !v1.IsEqual(v2) {
				t.Errorf("Expected %s == %s", tt.v1, tt.v2)
			}
		})
	}
}

func TestSemVerIsCompatible(t *testing.T) {
	tests := []struct {
		version    string
		constraint string
		expected   bool
		expectErr  bool
	}{
		// Exact match
		{"1.2.3", "1.2.3", true, false},
		{"1.2.3", "1.2.4", false, false},

		// Caret constraints
		{"1.2.3", "^1.2.0", true, false},
		{"1.3.0", "^1.2.0", true, false},
		{"2.0.0", "^1.2.0", false, false},
		{"0.9.0", "^1.2.0", false, false},

		// Tilde constraints
		{"1.2.3", "~1.2.0", true, false},
		{"1.2.9", "~1.2.0", true, false},
		{"1.3.0", "~1.2.0", false, false},

		// Comparison operators
		{"1.2.3", ">=1.2.0", true, false},
		{"1.1.9", ">=1.2.0", false, false},
		{"1.2.3", "<=1.2.5", true, false},
		{"1.2.6", "<=1.2.5", false, false},
		{"1.2.3", ">1.2.0", true, false},
		{"1.2.0", ">1.2.0", false, false},
		{"1.2.3", "<1.2.5", true, false},
		{"1.2.5", "<1.2.5", false, false},

		// Invalid constraints
		{"1.2.3", "invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.version+" "+tt.constraint, func(t *testing.T) {
			v, err := ParseSemVer(tt.version)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt.version, err)
			}

			result, err := v.IsCompatible(tt.constraint)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for constraint %s", tt.constraint)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for constraint %s: %v", tt.constraint, err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %t for %s %s, got %t", tt.expected, tt.version, tt.constraint, result)
			}
		})
	}
}

func TestSortVersions(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"1.0.0", "2.0.0", "1.1.0"},
			expected: []string{"1.0.0", "1.1.0", "2.0.0"},
		},
		{
			input:    []string{"1.0.0-beta", "1.0.0-alpha", "1.0.0"},
			expected: []string{"1.0.0-alpha", "1.0.0-beta", "1.0.0"},
		},
		{
			input:    []string{"v2.1.0", "v1.0.0", "v2.0.0"},
			expected: []string{"v1.0.0", "v2.0.0", "v2.1.0"},
		},
	}

	for _, tt := range tests {
		t.Run("sort versions", func(t *testing.T) {
			result, err := SortVersions(tt.input)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d versions, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected version %s at index %d, got %s", expected, i, result[i])
				}
			}
		})
	}
}

func TestGetLatestVersion(t *testing.T) {
	tests := []struct {
		input       []string
		expected    string
		expectError bool
	}{
		{
			input:    []string{"1.0.0", "2.0.0", "1.1.0"},
			expected: "2.0.0",
		},
		{
			input:    []string{"1.0.0-beta", "1.0.0-alpha", "1.0.0"},
			expected: "1.0.0",
		},
		{
			input:       []string{},
			expectError: true,
		},
		{
			input:       []string{"invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run("get latest", func(t *testing.T) {
			result, err := GetLatestVersion(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestComparePrereleases(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		expected int
	}{
		{"", "", 0},
		{"", "alpha", 1},
		{"alpha", "", -1},
		{"alpha", "beta", -1},
		{"alpha.1", "alpha.2", -1},
		{"alpha.10", "alpha.2", 1},
		{"alpha", "alpha.1", -1},
		{"1", "2", -1},
		{"10", "2", 1},
	}

	for _, tt := range tests {
		t.Run(tt.a+" vs "+tt.b, func(t *testing.T) {
			result := comparePrereleases(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}
