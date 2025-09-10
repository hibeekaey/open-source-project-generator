package version

import (
	"testing"
)

func TestCompatibilityMatrix(t *testing.T) {
	matrix := NewCompatibilityMatrix()

	// Add a test rule
	rule := CompatibilityRule{
		Package: "react",
		Version: "18.0.0",
		Compatible: map[string]string{
			"react-dom": "^18.0.0",
		},
		Incompatible: map[string]string{
			"react-dom": "<17.0.0",
		},
		Notes: "React 18 requires React DOM 18+",
	}
	matrix.AddRule(rule)

	t.Run("compatible packages", func(t *testing.T) {
		packages := map[string]string{
			"react":     "18.0.0",
			"react-dom": "18.2.0",
		}

		result, err := matrix.CheckCompatibility(packages)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !result.Compatible {
			t.Errorf("Expected packages to be compatible")
		}

		if len(result.Issues) != 0 {
			t.Errorf("Expected no issues, got %d", len(result.Issues))
		}

		if len(result.Warnings) != 1 {
			t.Errorf("Expected 1 warning, got %d", len(result.Warnings))
		}
	})

	t.Run("incompatible packages", func(t *testing.T) {
		packages := map[string]string{
			"react":     "18.0.0",
			"react-dom": "17.0.0",
		}

		result, err := matrix.CheckCompatibility(packages)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result.Compatible {
			t.Errorf("Expected packages to be incompatible")
		}

		if len(result.Issues) == 0 {
			t.Errorf("Expected compatibility issues")
		}
	})

	t.Run("known incompatible packages", func(t *testing.T) {
		packages := map[string]string{
			"react":     "18.0.0",
			"react-dom": "16.14.0",
		}

		result, err := matrix.CheckCompatibility(packages)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result.Compatible {
			t.Errorf("Expected packages to be incompatible")
		}

		// Should have both incompatible version and known incompatible issues
		if len(result.Issues) == 0 {
			t.Errorf("Expected compatibility issues")
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		packages := map[string]string{
			"react":     "invalid-version",
			"react-dom": "18.0.0",
		}

		result, err := matrix.CheckCompatibility(packages)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result.Compatible {
			t.Errorf("Expected packages to be incompatible due to invalid version")
		}

		found := false
		for _, issue := range result.Issues {
			if issue.Type == "invalid_version" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected invalid_version issue")
		}
	})

	t.Run("no rules for package", func(t *testing.T) {
		packages := map[string]string{
			"unknown-package": "1.0.0",
		}

		result, err := matrix.CheckCompatibility(packages)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !result.Compatible {
			t.Errorf("Expected packages to be compatible when no rules exist")
		}
	})
}

func TestGetDefaultCompatibilityMatrix(t *testing.T) {
	matrix := GetDefaultCompatibilityMatrix()

	t.Run("has Next.js rules", func(t *testing.T) {
		rules, exists := matrix.rules["next"]
		if !exists {
			t.Errorf("Expected Next.js rules to exist")
		}
		if len(rules) == 0 {
			t.Errorf("Expected at least one Next.js rule")
		}
	})

	t.Run("Next.js 15 compatibility", func(t *testing.T) {
		packages := map[string]string{
			"next":      "15.0.0",
			"react":     "18.2.0",
			"react-dom": "18.2.0",
		}

		result, err := matrix.CheckCompatibility(packages)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !result.Compatible {
			t.Errorf("Expected Next.js 15 with React 18 to be compatible")
		}
	})

	t.Run("Next.js with old React", func(t *testing.T) {
		packages := map[string]string{
			"next":      "15.0.0",
			"react":     "17.0.0",
			"react-dom": "17.0.0",
		}

		result, err := matrix.CheckCompatibility(packages)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result.Compatible {
			t.Errorf("Expected Next.js 15 with React 17 to be incompatible")
		}
	})
}

func TestValidatePackageSet(t *testing.T) {
	t.Run("valid package set", func(t *testing.T) {
		packages := map[string]string{
			"next":      "15.0.0",
			"react":     "18.2.0",
			"react-dom": "18.2.0",
		}

		result, err := ValidatePackageSet(packages)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !result.Compatible {
			t.Errorf("Expected package set to be valid")
		}
	})

	t.Run("invalid package set", func(t *testing.T) {
		packages := map[string]string{
			"next":  "15.0.0",
			"react": "17.0.0",
		}

		result, err := ValidatePackageSet(packages)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result.Compatible {
			t.Errorf("Expected package set to be invalid")
		}
	})
}

func TestSuggestCompatibleVersions(t *testing.T) {
	t.Run("already compatible", func(t *testing.T) {
		packages := map[string]string{
			"next":      "15.0.0",
			"react":     "18.2.0",
			"react-dom": "18.2.0",
		}

		availableVersions := map[string][]string{
			"react":     {"17.0.0", "18.0.0", "18.1.0", "18.2.0"},
			"react-dom": {"17.0.0", "18.0.0", "18.1.0", "18.2.0"},
		}

		suggestions, err := SuggestCompatibleVersions(packages, availableVersions)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Should return the same versions since they're already compatible
		for pkg, version := range packages {
			if suggestions[pkg] != version {
				t.Errorf("Expected suggestion for %s to be %s, got %s", pkg, version, suggestions[pkg])
			}
		}
	})

	t.Run("suggest compatible version", func(t *testing.T) {
		packages := map[string]string{
			"next":  "15.0.0",
			"react": "17.0.0", // Incompatible with Next.js 15
		}

		availableVersions := map[string][]string{
			"react": {"17.0.0", "18.0.0", "18.1.0", "18.2.0"},
		}

		suggestions, err := SuggestCompatibleVersions(packages, availableVersions)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Should suggest a React 18 version
		suggestedReact := suggestions["react"]
		if suggestedReact == "17.0.0" {
			t.Errorf("Expected React version to be upgraded from 17.0.0")
		}

		// Verify the suggested version is compatible
		reactVersion, err := ParseSemVer(suggestedReact)
		if err != nil {
			t.Errorf("Suggested React version is invalid: %v", err)
		}

		compatible, err := reactVersion.IsCompatible("^18.0.0")
		if err != nil {
			t.Errorf("Error checking compatibility: %v", err)
		}
		if !compatible {
			t.Errorf("Suggested React version %s is not compatible with ^18.0.0", suggestedReact)
		}
	})
}

func TestCompatibilityRule(t *testing.T) {
	matrix := NewCompatibilityMatrix()

	t.Run("multiple rules for same package", func(t *testing.T) {
		rule1 := CompatibilityRule{
			Package: "test-package",
			Version: "1.0.0",
			Compatible: map[string]string{
				"dep1": "^1.0.0",
			},
		}

		rule2 := CompatibilityRule{
			Package: "test-package",
			Version: "2.0.0",
			Compatible: map[string]string{
				"dep1": "^2.0.0",
			},
		}

		matrix.AddRule(rule1)
		matrix.AddRule(rule2)

		rules := matrix.rules["test-package"]
		if len(rules) != 2 {
			t.Errorf("Expected 2 rules for test-package, got %d", len(rules))
		}
	})

	t.Run("constraint error handling", func(t *testing.T) {
		rule := CompatibilityRule{
			Package: "test-package",
			Version: "1.0.0",
			Compatible: map[string]string{
				"dep1": "invalid-constraint",
			},
		}
		matrix.AddRule(rule)

		packages := map[string]string{
			"test-package": "1.0.0",
			"dep1":         "1.0.0",
		}

		result, err := matrix.CheckCompatibility(packages)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result.Compatible {
			t.Errorf("Expected incompatible result due to constraint error")
		}

		found := false
		for _, issue := range result.Issues {
			if issue.Type == "constraint_error" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected constraint_error issue")
		}
	})
}
