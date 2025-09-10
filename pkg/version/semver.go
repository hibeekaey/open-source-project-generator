package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// SemVer represents a semantic version
type SemVer struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
	Original   string
}

// semverRegex matches semantic version strings
var semverRegex = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z\-\.]+))?(?:\+([0-9A-Za-z\-\.]+))?$`)

// ParseSemVer parses a semantic version string
func ParseSemVer(version string) (*SemVer, error) {
	matches := semverRegex.FindStringSubmatch(strings.TrimSpace(version))
	if matches == nil {
		return nil, fmt.Errorf("invalid semantic version: %s", version)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", matches[3])
	}

	return &SemVer{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: matches[4],
		Build:      matches[5],
		Original:   version,
	}, nil
}

// String returns the string representation of the semantic version
func (v *SemVer) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)

	if v.Prerelease != "" {
		version += "-" + v.Prerelease
	}

	if v.Build != "" {
		version += "+" + v.Build
	}

	return version
}

// Compare compares two semantic versions
// Returns -1 if v < other, 0 if v == other, 1 if v > other
func (v *SemVer) Compare(other *SemVer) int {
	// Compare major version
	if v.Major < other.Major {
		return -1
	}
	if v.Major > other.Major {
		return 1
	}

	// Compare minor version
	if v.Minor < other.Minor {
		return -1
	}
	if v.Minor > other.Minor {
		return 1
	}

	// Compare patch version
	if v.Patch < other.Patch {
		return -1
	}
	if v.Patch > other.Patch {
		return 1
	}

	// Compare prerelease versions
	return comparePrereleases(v.Prerelease, other.Prerelease)
}

// IsGreaterThan returns true if v > other
func (v *SemVer) IsGreaterThan(other *SemVer) bool {
	return v.Compare(other) > 0
}

// IsLessThan returns true if v < other
func (v *SemVer) IsLessThan(other *SemVer) bool {
	return v.Compare(other) < 0
}

// IsEqual returns true if v == other
func (v *SemVer) IsEqual(other *SemVer) bool {
	return v.Compare(other) == 0
}

// IsCompatible checks if the version is compatible with a constraint
// Supports basic constraints like "^1.2.3", "~1.2.3", ">=1.2.3", etc.
func (v *SemVer) IsCompatible(constraint string) (bool, error) {
	constraint = strings.TrimSpace(constraint)

	// Handle caret constraint (^1.2.3)
	if strings.HasPrefix(constraint, "^") {
		return v.isCaretCompatible(constraint[1:])
	}

	// Handle tilde constraint (~1.2.3)
	if strings.HasPrefix(constraint, "~") {
		return v.isTildeCompatible(constraint[1:])
	}

	// Handle comparison operators
	if strings.HasPrefix(constraint, ">=") {
		target, err := ParseSemVer(constraint[2:])
		if err != nil {
			return false, err
		}
		return v.Compare(target) >= 0, nil
	}

	if strings.HasPrefix(constraint, "<=") {
		target, err := ParseSemVer(constraint[2:])
		if err != nil {
			return false, err
		}
		return v.Compare(target) <= 0, nil
	}

	if strings.HasPrefix(constraint, ">") {
		target, err := ParseSemVer(constraint[1:])
		if err != nil {
			return false, err
		}
		return v.Compare(target) > 0, nil
	}

	if strings.HasPrefix(constraint, "<") {
		target, err := ParseSemVer(constraint[1:])
		if err != nil {
			return false, err
		}
		return v.Compare(target) < 0, nil
	}

	// Handle exact match
	target, err := ParseSemVer(constraint)
	if err != nil {
		return false, err
	}

	return v.IsEqual(target), nil
}

// isCaretCompatible checks caret compatibility (^1.2.3)
// Compatible if major version matches and version is >= constraint
func (v *SemVer) isCaretCompatible(constraint string) (bool, error) {
	target, err := ParseSemVer(constraint)
	if err != nil {
		return false, err
	}

	// Major version must match
	if v.Major != target.Major {
		return false, nil
	}

	// Version must be >= target
	return v.Compare(target) >= 0, nil
}

// isTildeCompatible checks tilde compatibility (~1.2.3)
// Compatible if major and minor versions match and version is >= constraint
func (v *SemVer) isTildeCompatible(constraint string) (bool, error) {
	target, err := ParseSemVer(constraint)
	if err != nil {
		return false, err
	}

	// Major and minor versions must match
	if v.Major != target.Major || v.Minor != target.Minor {
		return false, nil
	}

	// Version must be >= target
	return v.Compare(target) >= 0, nil
}

// comparePrereleases compares prerelease versions according to SemVer spec
func comparePrereleases(a, b string) int {
	// No prerelease is greater than any prerelease
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return 1
	}
	if b == "" {
		return -1
	}

	// Split prerelease identifiers
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	// Compare each identifier
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		var aPart, bPart string

		if i < len(aParts) {
			aPart = aParts[i]
		}
		if i < len(bParts) {
			bPart = bParts[i]
		}

		// Missing part is less than any part
		if aPart == "" {
			return -1
		}
		if bPart == "" {
			return 1
		}

		// Try to parse as numbers
		aNum, aIsNum := parseNumber(aPart)
		bNum, bIsNum := parseNumber(bPart)

		// Numeric identifiers are compared numerically
		if aIsNum && bIsNum {
			if aNum < bNum {
				return -1
			}
			if aNum > bNum {
				return 1
			}
			continue
		}

		// Numeric identifiers are always less than non-numeric
		if aIsNum && !bIsNum {
			return -1
		}
		if !aIsNum && bIsNum {
			return 1
		}

		// Both are non-numeric, compare lexically
		if aPart < bPart {
			return -1
		}
		if aPart > bPart {
			return 1
		}
	}

	return 0
}

// parseNumber tries to parse a string as a number
func parseNumber(s string) (int, bool) {
	num, err := strconv.Atoi(s)
	return num, err == nil
}

// SortVersions sorts a slice of version strings in ascending order
func SortVersions(versions []string) ([]string, error) {
	semvers := make([]*SemVer, 0, len(versions))

	// Parse all versions
	for _, v := range versions {
		semver, err := ParseSemVer(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse version %s: %w", v, err)
		}
		semvers = append(semvers, semver)
	}

	// Sort using bubble sort (simple implementation)
	for i := 0; i < len(semvers); i++ {
		for j := i + 1; j < len(semvers); j++ {
			if semvers[i].Compare(semvers[j]) > 0 {
				semvers[i], semvers[j] = semvers[j], semvers[i]
			}
		}
	}

	// Convert back to strings
	result := make([]string, len(semvers))
	for i, v := range semvers {
		result[i] = v.Original
	}

	return result, nil
}

// GetLatestVersion returns the latest version from a slice of version strings
func GetLatestVersion(versions []string) (string, error) {
	if len(versions) == 0 {
		return "", fmt.Errorf("no versions provided")
	}

	sorted, err := SortVersions(versions)
	if err != nil {
		return "", err
	}

	return sorted[len(sorted)-1], nil
}
