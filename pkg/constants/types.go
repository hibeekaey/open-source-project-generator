// Package constants provides centralized constants for the Open Source Project Generator.
//
// This package defines commonly used constants across the application to ensure
// consistency, reduce duplication, and improve maintainability. Constants are
// organized by category for easy navigation and use.
package constants

// Package Manager Constants
const (
	// PackageManagerNPM represents the npm package manager
	PackageManagerNPM = "npm"

	// PackageManagerYarn represents the yarn package manager
	PackageManagerYarn = "yarn"

	// PackageManagerPNPM represents the pnpm package manager
	PackageManagerPNPM = "pnpm"

	// LanguageJavaScript represents the JavaScript language
	LanguageJavaScript = "javascript"

	// LanguageTypeScript represents the TypeScript language
	LanguageTypeScript = "typescript"

	// LanguageNodeJS represents the Node.js runtime
	LanguageNodeJS = "nodejs"

	// LanguageGo represents the Go programming language
	LanguageGo = "go"

	// LanguagePython represents the Python programming language
	LanguagePython = "python"
)

// Version and Update Constants
const (
	// VersionLatest represents the latest version identifier
	VersionLatest = "latest"

	// VersionUnknown represents an unknown version
	VersionUnknown = "unknown"

	// VersionDev represents a development version
	VersionDev = "dev"
)

// Severity Level Constants
const (
	// SeverityCritical represents critical severity level
	SeverityCritical = "critical"

	// SeverityHigh represents high severity level
	SeverityHigh = "high"

	// SeverityMedium represents medium severity level
	SeverityMedium = "medium"

	// SeverityLow represents low severity level
	SeverityLow = "low"

	// SeverityInfo represents informational severity level
	SeverityInfo = "info"
)

// Status Constants
const (
	// StatusSuccess represents a successful operation
	StatusSuccess = "success"

	// StatusFailed represents a failed operation
	StatusFailed = "failed"

	// StatusPending represents a pending operation
	StatusPending = "pending"

	// StatusSkipped represents a skipped operation
	StatusSkipped = "skipped"

	// StatusInProgress represents an operation in progress
	StatusInProgress = "in_progress"
)

// File Format Constants
const (
	// FormatJSON represents JSON file format
	FormatJSON = "json"

	// FormatYAML represents YAML file format
	FormatYAML = "yaml"

	// FormatYML represents YML file format (alternative YAML extension)
	FormatYML = "yml"

	// FormatTOML represents TOML file format
	FormatTOML = "toml"

	// FormatXML represents XML file format
	FormatXML = "xml"
)

// File Type Constants
const (
	// FileTypePackage represents package.json files
	FileTypePackage = "package"

	// FileTypeConfig represents configuration files
	FileTypeConfig = "config"

	// FileTypeTemplate represents template files
	FileTypeTemplate = "template"

	// FileTypeDocumentation represents documentation files
	FileTypeDocumentation = "documentation"

	// FileTypeScript represents script files
	FileTypeScript = "script"
)

// Template Type Constants
const (
	// TemplateFrontend represents frontend templates
	TemplateFrontend = "frontend"

	// TemplateBackend represents backend templates
	TemplateBackend = "backend"

	// TemplateMobile represents mobile templates
	TemplateMobile = "mobile"

	// TemplateInfrastructure represents infrastructure templates
	TemplateInfrastructure = "infrastructure"

	// TemplateBase represents base templates
	TemplateBase = "base"
)

// UI Symbol Constants
const (
	// SymbolSuccess represents a success symbol (✅)
	SymbolSuccess = "✅"

	// SymbolFailure represents a failure symbol (❌)
	SymbolFailure = "❌"

	// SymbolWarning represents a warning symbol (⚠️)
	SymbolWarning = "⚠️"

	// SymbolInfo represents an info symbol (ℹ️)
	SymbolInfo = "ℹ️"

	// SymbolCheck represents a check mark (✓)
	SymbolCheck = "✓"

	// SymbolCross represents a cross mark (✗)
	SymbolCross = "✗"

	// SymbolBullet represents a bullet point (•)
	SymbolBullet = "•"

	// SymbolArrow represents an arrow (→)
	SymbolArrow = "→"
)

// Template Path Constants
const (
	// TemplateBaseDir represents the base templates directory in embedded filesystem
	TemplateBaseDir = "templates"

	// EmbeddedTemplateDir represents the embedded templates directory name
	EmbeddedTemplateDir = "templates"

	// EmbeddedTemplateBasePath represents the full path to embedded templates from project root
	EmbeddedTemplateBasePath = "pkg/template/templates"

	// TemplateRelativeBasePath represents the relative path to templates from pkg/template directory
	// This is used for tests and operations that need filesystem access relative to pkg/template/
	TemplateRelativeBasePath = "templates"

	// TemplateConfigDir represents the template config directory
	TemplateConfigDir = "config"

	// TemplateDefaultsFile represents the defaults configuration file
	TemplateDefaultsFile = "defaults.yaml"
)

// Validation Constants
const (
	// ValidationConsistent represents consistent validation status
	ValidationConsistent = "consistent"

	// ValidationInconsistent represents inconsistent validation status
	ValidationInconsistent = "inconsistent"

	// ValidationRequired represents a required validation
	ValidationRequired = "required"

	// ValidationOptional represents an optional validation
	ValidationOptional = "optional"
)

// String Type Constants
const (
	// StringType represents the string data type
	StringType = "string"

	// NumberType represents the number data type
	NumberType = "number"

	// BooleanType represents the boolean data type
	BooleanType = "boolean"

	// ObjectType represents the object data type
	ObjectType = "object"

	// ArrayType represents the array data type
	ArrayType = "array"
)

// Additional Constants
const (
	// Priority levels
	PriorityHigh     = "high"
	PriorityCritical = "critical"

	// Status values
	StatusPresent = "present"
	StatusPassed  = "passed"

	// File names
	FilePackageJSON = "package.json"
	FileGoMod       = "go.mod"
	FileDockerfile  = "Dockerfile"

	// Language type
	TypeLanguage = "language"
)

// Error and Message Constants
const (
	// MessageFailed represents a generic failure message
	MessageFailed = "failed"

	// MessageSuccess represents a generic success message
	MessageSuccess = "success"

	// MessageProcessing represents a processing message
	MessageProcessing = "processing"

	// MessageCompleted represents a completion message
	MessageCompleted = "completed"
)

// Build and Development Constants
const (
	// BuildModeDevelopment represents development build mode
	BuildModeDevelopment = "development"

	// BuildModeProduction represents production build mode
	BuildModeProduction = "production"

	// BuildModeTest represents test build mode
	BuildModeTest = "test"

	// BuildModeStaging represents staging build mode
	BuildModeStaging = "staging"
)
