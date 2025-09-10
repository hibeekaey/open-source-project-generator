package version

import (
	"fmt"
)

// CompatibilityRule represents a version compatibility rule
type CompatibilityRule struct {
	Package      string            `json:"package"`
	Version      string            `json:"version"`
	Compatible   map[string]string `json:"compatible"`   // package -> version constraint
	Incompatible map[string]string `json:"incompatible"` // package -> version constraint
	Notes        string            `json:"notes,omitempty"`
}

// CompatibilityMatrix manages version compatibility rules
type CompatibilityMatrix struct {
	rules map[string][]CompatibilityRule
}

// NewCompatibilityMatrix creates a new compatibility matrix
func NewCompatibilityMatrix() *CompatibilityMatrix {
	return &CompatibilityMatrix{
		rules: make(map[string][]CompatibilityRule),
	}
}

// AddRule adds a compatibility rule to the matrix
func (cm *CompatibilityMatrix) AddRule(rule CompatibilityRule) {
	cm.rules[rule.Package] = append(cm.rules[rule.Package], rule)
}

// CheckCompatibility checks if a set of package versions are compatible
func (cm *CompatibilityMatrix) CheckCompatibility(packages map[string]string) (*CompatibilityResult, error) {
	result := &CompatibilityResult{
		Compatible: true,
		Issues:     []CompatibilityIssue{},
		Warnings:   []CompatibilityWarning{},
	}

	// Check each package against the rules
	for packageName, version := range packages {
		rules, exists := cm.rules[packageName]
		if !exists {
			continue // No rules for this package
		}

		packageVersion, err := ParseSemVer(version)
		if err != nil {
			result.Issues = append(result.Issues, CompatibilityIssue{
				Type:        "invalid_version",
				Package:     packageName,
				Version:     version,
				Description: fmt.Sprintf("Invalid version format: %v", err),
			})
			result.Compatible = false
			continue
		}

		// Find applicable rules for this version
		for _, rule := range rules {
			ruleVersion, err := ParseSemVer(rule.Version)
			if err != nil {
				continue // Skip invalid rule versions
			}

			// Check if this rule applies to the current version
			if !packageVersion.IsEqual(ruleVersion) && !isVersionInRange(packageVersion, rule.Version) {
				continue
			}

			// Check compatible packages
			for compatPackage, constraint := range rule.Compatible {
				if otherVersion, exists := packages[compatPackage]; exists {
					compatible, err := cm.checkConstraint(otherVersion, constraint)
					if err != nil {
						result.Issues = append(result.Issues, CompatibilityIssue{
							Type:        "constraint_error",
							Package:     compatPackage,
							Version:     otherVersion,
							Description: fmt.Sprintf("Error checking constraint %s: %v", constraint, err),
						})
						result.Compatible = false
					} else if !compatible {
						result.Issues = append(result.Issues, CompatibilityIssue{
							Type:        "incompatible_version",
							Package:     compatPackage,
							Version:     otherVersion,
							Constraint:  constraint,
							Description: fmt.Sprintf("%s@%s is not compatible with %s@%s (requires %s)", packageName, version, compatPackage, otherVersion, constraint),
						})
						result.Compatible = false
					}
				}
			}

			// Check incompatible packages
			for incompatPackage, constraint := range rule.Incompatible {
				if otherVersion, exists := packages[incompatPackage]; exists {
					incompatible, err := cm.checkConstraint(otherVersion, constraint)
					if err != nil {
						result.Issues = append(result.Issues, CompatibilityIssue{
							Type:        "constraint_error",
							Package:     incompatPackage,
							Version:     otherVersion,
							Description: fmt.Sprintf("Error checking incompatibility constraint %s: %v", constraint, err),
						})
						result.Compatible = false
					} else if incompatible {
						result.Issues = append(result.Issues, CompatibilityIssue{
							Type:        "known_incompatible",
							Package:     incompatPackage,
							Version:     otherVersion,
							Constraint:  constraint,
							Description: fmt.Sprintf("%s@%s is known to be incompatible with %s@%s", packageName, version, incompatPackage, otherVersion),
						})
						result.Compatible = false
					}
				}
			}

			// Add notes as warnings
			if rule.Notes != "" {
				result.Warnings = append(result.Warnings, CompatibilityWarning{
					Package:     packageName,
					Version:     version,
					Description: rule.Notes,
				})
			}
		}
	}

	return result, nil
}

// checkConstraint checks if a version satisfies a constraint
func (cm *CompatibilityMatrix) checkConstraint(version, constraint string) (bool, error) {
	v, err := ParseSemVer(version)
	if err != nil {
		return false, err
	}

	return v.IsCompatible(constraint)
}

// isVersionInRange checks if a version is in a range (simplified implementation)
func isVersionInRange(version *SemVer, rangeStr string) bool {
	// For now, just check exact match
	// This could be extended to support version ranges
	rangeVersion, err := ParseSemVer(rangeStr)
	if err != nil {
		return false
	}

	return version.IsEqual(rangeVersion)
}

// CompatibilityResult represents the result of a compatibility check
type CompatibilityResult struct {
	Compatible bool                   `json:"compatible"`
	Issues     []CompatibilityIssue   `json:"issues"`
	Warnings   []CompatibilityWarning `json:"warnings"`
}

// CompatibilityIssue represents a compatibility issue
type CompatibilityIssue struct {
	Type        string `json:"type"`
	Package     string `json:"package"`
	Version     string `json:"version"`
	Constraint  string `json:"constraint,omitempty"`
	Description string `json:"description"`
}

// CompatibilityWarning represents a compatibility warning
type CompatibilityWarning struct {
	Package     string `json:"package"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// GetDefaultCompatibilityMatrix returns a matrix with common compatibility rules
func GetDefaultCompatibilityMatrix() *CompatibilityMatrix {
	matrix := NewCompatibilityMatrix()

	// React and Next.js compatibility rules
	matrix.AddRule(CompatibilityRule{
		Package: "next",
		Version: "15.0.0",
		Compatible: map[string]string{
			"react":     "^18.0.0",
			"react-dom": "^18.0.0",
		},
		Notes: "Next.js 15 requires React 18 or later",
	})

	matrix.AddRule(CompatibilityRule{
		Package: "next",
		Version: "14.0.0",
		Compatible: map[string]string{
			"react":     "^18.0.0",
			"react-dom": "^18.0.0",
		},
		Notes: "Next.js 14 requires React 18 or later",
	})

	// TypeScript compatibility
	matrix.AddRule(CompatibilityRule{
		Package: "typescript",
		Version: "5.0.0",
		Compatible: map[string]string{
			"@types/node":  "^20.0.0",
			"@types/react": "^18.0.0",
		},
	})

	// Tailwind CSS compatibility
	matrix.AddRule(CompatibilityRule{
		Package: "tailwindcss",
		Version: "3.4.0",
		Compatible: map[string]string{
			"autoprefixer": "^10.0.0",
			"postcss":      "^8.0.0",
		},
	})

	// Go framework compatibility
	matrix.AddRule(CompatibilityRule{
		Package: "github.com/gin-gonic/gin",
		Version: "1.9.0",
		Compatible: map[string]string{
			"go": ">=1.19",
		},
		Notes: "Gin v1.9+ requires Go 1.19 or later",
	})

	// GORM compatibility
	matrix.AddRule(CompatibilityRule{
		Package: "gorm.io/gorm",
		Version: "1.25.0",
		Compatible: map[string]string{
			"go": ">=1.18",
		},
	})

	// Android/Kotlin compatibility
	matrix.AddRule(CompatibilityRule{
		Package: "kotlin",
		Version: "2.0.0",
		Compatible: map[string]string{
			"android-gradle-plugin": "^8.0.0",
		},
		Notes: "Kotlin 2.0 requires Android Gradle Plugin 8.0+",
	})

	return matrix
}

// ValidatePackageSet validates a set of packages for compatibility
func ValidatePackageSet(packages map[string]string) (*CompatibilityResult, error) {
	matrix := GetDefaultCompatibilityMatrix()
	return matrix.CheckCompatibility(packages)
}

// SuggestCompatibleVersions suggests compatible versions for a package set
func SuggestCompatibleVersions(packages map[string]string, availableVersions map[string][]string) (map[string]string, error) {
	suggestions := make(map[string]string)

	// Start with the provided versions
	for pkg, version := range packages {
		suggestions[pkg] = version
	}

	// Check compatibility and suggest alternatives if needed
	result, err := ValidatePackageSet(suggestions)
	if err != nil {
		return nil, err
	}

	if result.Compatible {
		return suggestions, nil
	}

	// Try to resolve compatibility issues
	for _, issue := range result.Issues {
		if issue.Type == "incompatible_version" && issue.Constraint != "" {
			// Try to find a compatible version
			if versions, exists := availableVersions[issue.Package]; exists {
				for _, version := range versions {
					v, err := ParseSemVer(version)
					if err != nil {
						continue
					}

					compatible, err := v.IsCompatible(issue.Constraint)
					if err != nil {
						continue
					}

					if compatible {
						suggestions[issue.Package] = version
						break
					}
				}
			}
		}
	}

	return suggestions, nil
}
