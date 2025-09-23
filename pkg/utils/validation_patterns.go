package utils

import "regexp"

// Common validation patterns used across the application
var (
	// ProjectNamePattern validates project names (alphanumeric, hyphens, underscores)
	ProjectNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_]*$`)

	// EmailPattern validates email addresses
	EmailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// URLPattern validates HTTP/HTTPS URLs
	URLPattern = regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)

	// PackageNamePattern validates package names (lowercase, hyphens)
	PackageNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`)

	// VersionPattern validates semantic version numbers
	VersionPattern = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*))?(?:\+([a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*))?$`)

	// GitHubRepoPattern validates GitHub repository names
	GitHubRepoPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$`)

	// SemverPattern validates semantic version strings
	SemverPattern = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

	// GoVersionPattern validates Go version strings
	GoVersionPattern = regexp.MustCompile(`^1\.\d+(\.\d+)?$`)

	// NPMNamePattern validates NPM package names
	NPMNamePattern = regexp.MustCompile(`^(@[a-z0-9-~][a-z0-9-._~]*/)?[a-z0-9-~][a-z0-9-._~]*$`)

	// PythonNamePattern validates Python package names
	PythonNamePattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9._-]*[a-zA-Z0-9])?$`)

	// CamelCasePattern validates camelCase strings
	CamelCasePattern = regexp.MustCompile(`^[a-z]+([A-Z][a-z]*)+$`)

	// EnvKeyPattern validates environment variable keys
	EnvKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
)
