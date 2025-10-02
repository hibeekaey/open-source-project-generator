package security

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDependencyChecker(t *testing.T) {
	checker := NewDependencyChecker()

	assert.NotNil(t, checker)
	assert.NotNil(t, checker.vulnerabilityDatabase)
	assert.NotNil(t, checker.ecosystemParsers)

	// Check that vulnerability database is initialized
	assert.Greater(t, len(checker.vulnerabilityDatabase), 0)

	// Check for some known vulnerabilities
	lodashVulns, exists := checker.vulnerabilityDatabase["lodash"]
	assert.True(t, exists)
	assert.Greater(t, len(lodashVulns), 0)
}

func TestScanDependencyVulnerabilities(t *testing.T) {
	checker := NewDependencyChecker()

	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "dependency_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a package.json with vulnerable dependencies
	packageJSON := `{
		"name": "test-project",
		"dependencies": {
			"lodash": "4.17.10",
			"axios": "0.21.1",
			"express": "4.17.0"
		}
	}`

	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	require.NoError(t, err)

	// Scan for vulnerabilities
	vulnerabilities, err := checker.ScanDependencyVulnerabilities(tempDir)
	require.NoError(t, err)

	// Should find vulnerabilities
	assert.Greater(t, len(vulnerabilities), 0)

	// Check that we found lodash vulnerabilities
	foundLodash := false
	for _, vuln := range vulnerabilities {
		if vuln.Package == "lodash" {
			foundLodash = true
			assert.Equal(t, "high", vuln.Severity)
			assert.Contains(t, vuln.Title, "lodash")
		}
	}
	assert.True(t, foundLodash, "Should find lodash vulnerability")
}

func TestScanDependencyVulnerabilitiesGoMod(t *testing.T) {
	checker := NewDependencyChecker()

	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "dependency_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a go.mod with vulnerable dependencies
	goMod := `module test-project

go 1.19

require (
	github.com/gin-gonic/gin v1.8.0
	github.com/stretchr/testify v1.8.0
)`

	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0644)
	require.NoError(t, err)

	// Scan for vulnerabilities
	vulnerabilities, err := checker.ScanDependencyVulnerabilities(tempDir)
	require.NoError(t, err)

	// Should find vulnerabilities
	assert.Greater(t, len(vulnerabilities), 0)

	// Check that we found gin vulnerabilities
	foundGin := false
	for _, vuln := range vulnerabilities {
		if vuln.Package == "github.com/gin-gonic/gin" {
			foundGin = true
			assert.Equal(t, "medium", vuln.Severity)
		}
	}
	assert.True(t, foundGin, "Should find gin vulnerability")
}

func TestAnalyzeDependencies(t *testing.T) {
	checker := NewDependencyChecker()

	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "dependency_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a package.json
	packageJSON := `{
		"name": "test-project",
		"dependencies": {
			"lodash": "4.17.21",
			"axios": "0.21.2"
		},
		"devDependencies": {
			"jest": "27.0.0"
		}
	}`

	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	require.NoError(t, err)

	// Analyze dependencies
	result, err := checker.AnalyzeDependencies(tempDir)
	require.NoError(t, err)

	assert.NotNil(t, result)
	assert.Greater(t, len(result.Dependencies), 0)
	assert.Equal(t, len(result.Dependencies), result.Summary.TotalDependencies)

	// Check that we have both direct and dev dependencies
	foundDirect := false
	foundDev := false
	for _, dep := range result.Dependencies {
		if dep.Type == "direct" {
			foundDirect = true
		}
		if dep.Type == "dev" {
			foundDev = true
		}
	}
	assert.True(t, foundDirect, "Should find direct dependencies")
	assert.True(t, foundDev, "Should find dev dependencies")
}

func TestCheckDependencyVulnerabilities(t *testing.T) {
	checker := NewDependencyChecker()

	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "dependency_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a package.json with vulnerable dependencies
	packageJSON := `{
		"dependencies": {
			"lodash": "4.17.10"
		}
	}`

	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	require.NoError(t, err)

	// Check for policy violations
	violations, err := checker.CheckDependencyVulnerabilities(tempDir)
	require.NoError(t, err)

	assert.Greater(t, len(violations), 0)

	// Check violation details
	found := false
	for _, violation := range violations {
		if violation.Policy == "SEC-002" {
			found = true
			assert.Contains(t, violation.Description, "lodash")
			assert.Equal(t, "high", violation.Severity)
		}
	}
	assert.True(t, found, "Should find SEC-002 policy violation")
}

func TestGetEcosystemFromFile(t *testing.T) {
	checker := NewDependencyChecker()

	tests := []struct {
		filename string
		expected string
	}{
		{"package.json", "npm"},
		{"package-lock.json", "npm"},
		{"yarn.lock", "npm"},
		{"go.mod", "go"},
		{"go.sum", "go"},
		{"requirements.txt", "pip"},
		{"Pipfile", "pip"},
		{"pom.xml", "maven"},
		{"build.gradle", "gradle"},
		{"Cargo.toml", "cargo"},
		{"composer.json", "composer"},
		{"Gemfile", "bundler"},
		{"unknown.txt", ""},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			result := checker.getEcosystemFromFile(test.filename)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestParseNPMDependencies(t *testing.T) {
	checker := NewDependencyChecker()

	packageJSON := `{
		"name": "test-project",
		"dependencies": {
			"lodash": "^4.17.21",
			"axios": "~0.21.2"
		},
		"devDependencies": {
			"jest": "^27.0.0"
		}
	}`

	dependencies := checker.parseNPMDependencies(packageJSON)

	assert.Len(t, dependencies, 3)

	// Check dependencies
	foundLodash := false
	foundAxios := false
	foundJest := false

	for _, dep := range dependencies {
		switch dep.Name {
		case "lodash":
			foundLodash = true
			assert.Equal(t, "^4.17.21", dep.Version)
			assert.Equal(t, "direct", dep.Type)
		case "axios":
			foundAxios = true
			assert.Equal(t, "~0.21.2", dep.Version)
			assert.Equal(t, "direct", dep.Type)
		case "jest":
			foundJest = true
			assert.Equal(t, "^27.0.0", dep.Version)
			assert.Equal(t, "dev", dep.Type)
		}
	}

	assert.True(t, foundLodash, "Should find lodash dependency")
	assert.True(t, foundAxios, "Should find axios dependency")
	assert.True(t, foundJest, "Should find jest dev dependency")
}

func TestParseGoDependencies(t *testing.T) {
	checker := NewDependencyChecker()

	goMod := `module test-project

go 1.19

require (
	github.com/gin-gonic/gin v1.9.0
	github.com/stretchr/testify v1.8.0
)

require github.com/gorilla/mux v1.8.0
`

	dependencies := checker.parseGoDependencies(goMod)

	assert.Greater(t, len(dependencies), 0)

	// Check that we found gin dependency
	foundGin := false
	for _, dep := range dependencies {
		if dep.Name == "github.com/gin-gonic/gin" {
			foundGin = true
			assert.Equal(t, "v1.9.0", dep.Version)
			assert.Equal(t, "direct", dep.Type)
		}
	}
	assert.True(t, foundGin, "Should find gin dependency")
}

func TestParsePipDependencies(t *testing.T) {
	checker := NewDependencyChecker()

	requirements := `# This is a comment
Django==3.2.0
requests>=2.25.0
pytest==6.2.0
# Another comment
flask~=2.0.0`

	dependencies := checker.parsePipDependencies(requirements)

	assert.Len(t, dependencies, 4)

	// Check dependencies
	foundDjango := false
	for _, dep := range dependencies {
		if dep.Name == "Django" {
			foundDjango = true
			assert.Equal(t, "3.2.0", dep.Version)
			assert.Equal(t, "direct", dep.Type)
		}
	}
	assert.True(t, foundDjango, "Should find Django dependency")
}

func TestParseCargoDependencies(t *testing.T) {
	checker := NewDependencyChecker()

	cargoToml := `[package]
name = "test-project"
version = "0.1.0"

[dependencies]
serde = "1.0"
tokio = { version = "1.0", features = ["full"] }
reqwest = "0.11"

[dev-dependencies]
tokio-test = "0.4"`

	dependencies := checker.parseCargoDependencies(cargoToml)

	assert.Greater(t, len(dependencies), 0)

	// Check that we found serde dependency
	foundSerde := false
	for _, dep := range dependencies {
		if dep.Name == "serde" {
			foundSerde = true
			assert.Equal(t, "1.0", dep.Version)
			assert.Equal(t, "direct", dep.Type)
		}
	}
	assert.True(t, foundSerde, "Should find serde dependency")
}

func TestIsVersionVulnerable(t *testing.T) {
	checker := NewDependencyChecker()

	tests := []struct {
		version    string
		constraint string
		expected   bool
	}{
		{"4.17.10", "<4.17.21", true},
		{"4.17.21", "<4.17.21", false},
		{"4.17.22", "<4.17.21", false},
		{"1.0.0", "1.0.0", true},
		{"1.0.1", "1.0.0", false},
		{"", "<4.17.21", false},
		{"4.17.10", "", false},
	}

	for _, test := range tests {
		t.Run(test.version+"_"+test.constraint, func(t *testing.T) {
			result := checker.isVersionVulnerable(test.version, test.constraint)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestCheckDependencyForVulnerabilities(t *testing.T) {
	checker := NewDependencyChecker()

	// Test with known vulnerable package
	vulns := checker.checkDependencyForVulnerabilities("lodash", "4.17.10", "npm")
	assert.Greater(t, len(vulns), 0)

	// Test with safe version
	vulns = checker.checkDependencyForVulnerabilities("lodash", "4.17.21", "npm")
	// Note: This might still return vulnerabilities due to simplified version checking
	// In a real implementation with proper semver, this would return 0

	// Test with unknown package
	unknownVulns := checker.checkDependencyForVulnerabilities("unknown-package", "1.0.0", "npm")
	assert.Equal(t, 0, len(unknownVulns))
	assert.Len(t, vulns, 0)
}

func TestScanForVulnerabilityPatterns(t *testing.T) {
	checker := NewDependencyChecker()

	// Test NPM patterns
	npmContent := `{
		"dependencies": {
			"lodash": "4.17.10",
			"express": "4.17.0"
		}
	}`

	vulns := checker.scanForVulnerabilityPatterns(npmContent, "npm")
	assert.Greater(t, len(vulns), 0)

	// Test Go patterns
	goContent := `require (
		github.com/gin-gonic/gin v1.8.0
	)`

	vulns = checker.scanForVulnerabilityPatterns(goContent, "go")
	assert.Greater(t, len(vulns), 0)
}

func TestScanDependencyFileErrors(t *testing.T) {
	checker := NewDependencyChecker()

	// Test with non-existent file
	vulns, err := checker.scanDependencyFile("/non/existent/file.json")
	assert.Error(t, err)
	assert.Nil(t, vulns)
}

func TestAnalyzeDependencyFileErrors(t *testing.T) {
	checker := NewDependencyChecker()

	// Test with non-existent file
	deps, err := checker.analyzeDependencyFile("/non/existent/file.json", "npm")
	assert.Error(t, err)
	assert.Nil(t, deps)
}

func TestParseComposerDependencies(t *testing.T) {
	checker := NewDependencyChecker()

	composerJSON := `{
		"name": "test/project",
		"require": {
			"php": "^8.0",
			"symfony/console": "^5.0"
		},
		"require-dev": {
			"phpunit/phpunit": "^9.0"
		}
	}`

	dependencies := checker.parseComposerDependencies(composerJSON)

	assert.Greater(t, len(dependencies), 0)

	// Check that we found symfony dependency
	foundSymfony := false
	foundPHPUnit := false
	for _, dep := range dependencies {
		if dep.Name == "symfony/console" {
			foundSymfony = true
			assert.Equal(t, "^5.0", dep.Version)
			assert.Equal(t, "direct", dep.Type)
		}
		if dep.Name == "phpunit/phpunit" {
			foundPHPUnit = true
			assert.Equal(t, "^9.0", dep.Version)
			assert.Equal(t, "dev", dep.Type)
		}
	}
	assert.True(t, foundSymfony, "Should find symfony dependency")
	assert.True(t, foundPHPUnit, "Should find phpunit dev dependency")
}

func TestParseBundlerDependencies(t *testing.T) {
	checker := NewDependencyChecker()

	gemfile := `source 'https://rubygems.org'

gem 'rails', '~> 7.0.0'
gem 'pg', '~> 1.1'
gem 'puma', '~> 5.0'

group :development, :test do
  gem 'rspec-rails'
end`

	dependencies := checker.parseBundlerDependencies(gemfile)

	assert.Greater(t, len(dependencies), 0)

	// Check that we found rails dependency
	foundRails := false
	for _, dep := range dependencies {
		if dep.Name == "rails" {
			foundRails = true
			assert.Equal(t, "~> 7.0.0", dep.Version)
			assert.Equal(t, "direct", dep.Type)
		}
	}
	assert.True(t, foundRails, "Should find rails dependency")
}

func TestParseMavenDependencies(t *testing.T) {
	checker := NewDependencyChecker()

	pomXML := `<?xml version="1.0" encoding="UTF-8"?>
<project>
	<dependencies>
		<dependency>
			<groupId>org.springframework</groupId>
			<artifactId>spring-core</artifactId>
			<version>5.3.0</version>
		</dependency>
		<dependency>
			<groupId>junit</groupId>
			<artifactId>junit</artifactId>
			<version>4.13.2</version>
		</dependency>
	</dependencies>
</project>`

	dependencies := checker.parseMavenDependencies(pomXML)

	assert.Greater(t, len(dependencies), 0)

	// Check that we found spring dependency
	foundSpring := false
	for _, dep := range dependencies {
		if dep.Name == "org.springframework:spring-core" {
			foundSpring = true
			assert.Equal(t, "5.3.0", dep.Version)
			assert.Equal(t, "direct", dep.Type)
		}
	}
	assert.True(t, foundSpring, "Should find spring dependency")
}

func TestParseGradleDependencies(t *testing.T) {
	checker := NewDependencyChecker()

	buildGradle := `dependencies {
    implementation 'org.springframework:spring-core:5.3.0'
    testImplementation 'junit:junit:4.13.2'
    api 'com.google.guava:guava:30.0-jre'
}`

	dependencies := checker.parseGradleDependencies(buildGradle)

	assert.Greater(t, len(dependencies), 0)

	// Check that we found spring dependency
	foundSpring := false
	for _, dep := range dependencies {
		if dep.Name == "org.springframework:spring-core" {
			foundSpring = true
			assert.Equal(t, "5.3.0", dep.Version)
			assert.Equal(t, "direct", dep.Type)
		}
	}
	assert.True(t, foundSpring, "Should find spring dependency")
}

func TestCompareVersions(t *testing.T) {
	checker := NewDependencyChecker()

	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "1.0.1", -1},
		{"1.0.1", "1.0.0", 1},
		{"2.0.0", "1.9.9", 1},
	}

	for _, test := range tests {
		t.Run(test.v1+"_vs_"+test.v2, func(t *testing.T) {
			result := checker.compareVersions(test.v1, test.v2)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestVulnerabilityDatabaseInitialization(t *testing.T) {
	checker := NewDependencyChecker()

	// Check that common vulnerabilities are present
	expectedPackages := []string{"lodash", "axios", "express", "github.com/gin-gonic/gin"}

	for _, pkg := range expectedPackages {
		vulns, exists := checker.vulnerabilityDatabase[pkg]
		assert.True(t, exists, "Should have vulnerabilities for %s", pkg)
		assert.Greater(t, len(vulns), 0, "Should have at least one vulnerability for %s", pkg)

		// Check vulnerability structure
		for _, vuln := range vulns {
			assert.NotEmpty(t, vuln.ID, "Vulnerability should have an ID")
			assert.NotEmpty(t, vuln.Severity, "Vulnerability should have a severity")
			assert.NotEmpty(t, vuln.Title, "Vulnerability should have a title")
			assert.NotEmpty(t, vuln.Package, "Vulnerability should have a package")
		}
	}
}

func TestIntegrationWithMultipleEcosystems(t *testing.T) {
	checker := NewDependencyChecker()

	// Create a temporary directory with multiple ecosystem files
	tempDir, err := os.MkdirTemp("", "multi_ecosystem_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create package.json
	packageJSON := `{
		"dependencies": {
			"lodash": "4.17.10"
		}
	}`
	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	require.NoError(t, err)

	// Create go.mod
	goMod := `module test
require github.com/gin-gonic/gin v1.8.0`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0644)
	require.NoError(t, err)

	// Create requirements.txt
	requirements := `Django==3.0.0`
	err = os.WriteFile(filepath.Join(tempDir, "requirements.txt"), []byte(requirements), 0644)
	require.NoError(t, err)

	// Scan for vulnerabilities
	vulnerabilities, err := checker.ScanDependencyVulnerabilities(tempDir)
	require.NoError(t, err)

	// Should find vulnerabilities from multiple ecosystems
	assert.Greater(t, len(vulnerabilities), 0)

	// Analyze dependencies
	result, err := checker.AnalyzeDependencies(tempDir)
	require.NoError(t, err)

	assert.Greater(t, len(result.Dependencies), 0)
	assert.Greater(t, result.Summary.TotalDependencies, 0)
}

// Benchmark tests
func BenchmarkScanDependencyVulnerabilities(b *testing.B) {
	checker := NewDependencyChecker()

	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "benchmark_test")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	// Create a large package.json
	packageJSON := `{
		"dependencies": {
			"lodash": "4.17.10",
			"axios": "0.21.1",
			"express": "4.17.0",
			"react": "17.0.0",
			"vue": "3.0.0"
		}
	}`

	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := checker.ScanDependencyVulnerabilities(tempDir)
		require.NoError(b, err)
	}
}

func BenchmarkAnalyzeDependencies(b *testing.B) {
	checker := NewDependencyChecker()

	// Create a temporary directory with test files
	tempDir, err := os.MkdirTemp("", "benchmark_test")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	// Create a large package.json
	packageJSON := `{
		"dependencies": {
			"lodash": "4.17.21",
			"axios": "0.21.2",
			"express": "4.18.0",
			"react": "17.0.0",
			"vue": "3.0.0"
		}
	}`

	err = os.WriteFile(filepath.Join(tempDir, "package.json"), []byte(packageJSON), 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := checker.AnalyzeDependencies(tempDir)
		require.NoError(b, err)
	}
}
