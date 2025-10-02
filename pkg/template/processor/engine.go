// Package processor provides template processing capabilities for the Open Source Project Generator.
//
// This package implements the core template processing engine that handles:
//   - Template compilation and execution
//   - Template function registration and management
//   - Directory processing for complete project generation
//   - Template rendering with variable substitution
//
// The ProcessingEngine provides a focused interface for template processing operations,
// separating concerns from template discovery, validation, and caching.
package processor

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode"

	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Buffer pool for efficient memory reuse in template rendering
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// ProcessingEngine implements template processing functionality
type ProcessingEngine struct {
	funcMap template.FuncMap
}

// NewProcessingEngine creates a new template processing engine
func NewProcessingEngine() *ProcessingEngine {
	engine := &ProcessingEngine{
		funcMap: make(template.FuncMap),
	}
	engine.registerDefaultFunctions()
	return engine
}

// ProcessTemplate processes a single template file with the given configuration
func (e *ProcessingEngine) ProcessTemplate(templatePath string, config *models.ProjectConfig) ([]byte, error) {
	// Load the template
	tmpl, err := e.LoadTemplate(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	// Render the template
	content, err := e.RenderTemplate(tmpl, config)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	return content, nil
}

// ProcessDirectory processes an entire template directory recursively
func (e *ProcessingEngine) ProcessDirectory(templateDir string, outputDir string, config *models.ProjectConfig) error {
	// Create output directory if it doesn't exist with secure permissions
	if err := utils.SafeMkdirAll(outputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Walk through template directory
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Skip if it's the template directory itself
		if relPath == "." {
			return nil
		}

		// Skip disabled template files
		if strings.HasSuffix(path, ".tmpl.disabled") {
			return nil
		}

		// Calculate output path
		outputPath := filepath.Join(outputDir, relPath)

		// Remove .tmpl extension if present
		outputPath = strings.TrimSuffix(outputPath, ".tmpl")

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(outputPath, info.Mode())
		}

		// Skip disabled template files
		if strings.HasSuffix(path, ".tmpl.disabled") {
			return nil // Skip disabled templates
		}

		// Process file
		if strings.HasSuffix(path, ".tmpl") {
			// Process template file
			content, err := e.ProcessTemplate(path, config)
			if err != nil {
				return fmt.Errorf("failed to process template %s: %w", path, err)
			}

			// Write processed content
			return os.WriteFile(outputPath, content, info.Mode())
		} else {
			// Copy non-template file
			return e.copyFile(path, outputPath, info.Mode())
		}
	})
}

// LoadTemplate loads and parses a template from the given path
func (e *ProcessingEngine) LoadTemplate(templatePath string) (*template.Template, error) {
	// Validate path to prevent directory traversal
	if err := utils.ValidatePath(templatePath); err != nil {
		return nil, fmt.Errorf("invalid template path: %w", err)
	}

	content, err := utils.SafeReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Create template with custom functions
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(e.funcMap).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

// RenderTemplate renders a template with the provided data
func (e *ProcessingEngine) RenderTemplate(tmpl *template.Template, data interface{}) ([]byte, error) {
	// Get buffer from pool for efficient memory reuse
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()               // Clear any existing content
	defer bufferPool.Put(buf) // Return buffer to pool when done

	if err := tmpl.Execute(buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Make a copy of the buffer contents since we're returning the buffer to the pool
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// RegisterFunctions registers custom template functions
func (e *ProcessingEngine) RegisterFunctions(funcMap template.FuncMap) {
	for name, fn := range funcMap {
		e.funcMap[name] = fn
	}
}

// GetFunctionMap returns the current function map
func (e *ProcessingEngine) GetFunctionMap() template.FuncMap {
	return e.funcMap
}

// copyFile copies a file from src to dst with the specified mode
func (e *ProcessingEngine) copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := utils.SafeOpen(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := utils.SafeCreate(dst)
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return os.Chmod(dst, mode)
}

// registerDefaultFunctions registers the default template functions
func (e *ProcessingEngine) registerDefaultFunctions() {
	e.funcMap = template.FuncMap{
		// String manipulation functions
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"title":      cases.Title(language.English).String,
		"camelCase":  toCamelCase,
		"snakeCase":  toSnakeCase,
		"kebabCase":  toKebabCase,
		"pascalCase": toPascalCase,
		"trim":       strings.TrimSpace,
		"replace":    strings.ReplaceAll,
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"split":      strings.Split,
		"join":       strings.Join,

		// Version handling functions
		"semverMajor":    getSemverMajor,
		"semverMinor":    getSemverMinor,
		"semverPatch":    getSemverPatch,
		"semverCompare":  compareSemver,
		"latestVersion":  getLatestVersion,
		"versionPrefix":  addVersionPrefix,
		"stripVersion":   stripVersionPrefix,
		"nodeVersion":    getNodeVersion,
		"goVersion":      getGoVersion,
		"nextjsVersion":  getNextjsVersion,
		"reactVersion":   getReactVersion,
		"kotlinVersion":  getKotlinVersion,
		"swiftVersion":   getSwiftVersion,
		"packageVersion": getPackageVersion,
		"hasPackage":     hasPackage,

		// GitHub Actions template helpers
		"secrets": githubSecrets,
		"matrix":  githubMatrix,
		"github":  githubContext,

		// Enhanced Node.js version functions
		"nodeRuntime":      getNodeRuntime,
		"nodeTypesVersion": getNodeTypesVersion,
		"nodeNPMVersion":   getNodeNPMVersion,
		"nodeDockerImage":  getNodeDockerImage,
		"isNodeLTS":        isNodeLTS,

		// Conditional functions
		"if":       templateIf,
		"ifnot":    templateIfNot,
		"and":      templateAnd,
		"or":       templateOr,
		"not":      templateNot,
		"eq":       templateEq,
		"ne":       templateNe,
		"lt":       templateLt,
		"le":       templateLe,
		"gt":       templateGt,
		"ge":       templateGe,
		"empty":    isEmpty,
		"nonempty": isNonEmpty,

		// Component checking functions
		"hasFrontend":       hasFrontendComponent,
		"hasBackend":        hasBackendComponent,
		"hasMobile":         hasMobileComponent,
		"hasInfrastructure": hasInfrastructureComponent,
		"hasComponent":      hasComponent,

		// Utility functions
		"now":        time.Now,
		"formatTime": formatTime,
		"add":        add,
		"sub":        sub,
		"mul":        mul,
		"div":        div,
		"mod":        mod,
		"default":    defaultValue,
		"coalesce":   coalesce,
		"quote":      quote,
		"squote":     singleQuote,
		"indent":     indent,
		"nindent":    nindent,
		"env":        getEnvVar,
		"nonce":      generateNonce,
		"customVar":  getCustomVar,
		"slice":      sliceString,
	}
}

// String manipulation functions

func toCamelCase(s string) string {
	if s == "" {
		return s
	}

	// Handle already camelCase strings
	if unicode.IsLower(rune(s[0])) {
		// Check if it's already camelCase
		hasUpper := false
		for _, r := range s[1:] {
			if unicode.IsUpper(r) {
				hasUpper = true
				break
			}
		}
		if hasUpper && !strings.ContainsAny(s, "-_ ") {
			return s // Already camelCase
		}
	}

	// Handle PascalCase (convert to camelCase)
	if unicode.IsUpper(rune(s[0])) && !strings.ContainsAny(s, "-_ ") {
		return strings.ToLower(string(s[0])) + s[1:]
	}

	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	if len(words) == 0 {
		return strings.ToLower(s)
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		word := strings.ToLower(words[i])
		if len(word) > 0 {
			result += strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return result
}

func toSnakeCase(s string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

func toKebabCase(s string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	kebab := re.ReplaceAllString(s, "${1}-${2}")
	return strings.ToLower(kebab)
}

func toPascalCase(s string) string {
	if s == "" {
		return s
	}

	// Handle camelCase input (e.g., "helloWorld" -> "HelloWorld")
	if unicode.IsLower(rune(s[0])) && !strings.ContainsAny(s, "-_ ") {
		// Check if it's camelCase
		hasUpper := false
		for _, r := range s[1:] {
			if unicode.IsUpper(r) {
				hasUpper = true
				break
			}
		}
		if hasUpper {
			return strings.ToUpper(string(s[0])) + s[1:]
		}
	}

	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	if len(words) == 0 {
		// Handle single word case
		if len(s) > 0 {
			return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
		}
		return s
	}

	var result strings.Builder
	for _, word := range words {
		word = strings.ToLower(word)
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(string(word[0])) + word[1:])
		}
	}

	return result.String()
}

// Version handling functions

func getSemverMajor(version string) string {
	if version == "" {
		return "0"
	}
	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return "0"
}

func getSemverMinor(version string) string {
	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")
	if len(parts) > 1 {
		return parts[1]
	}
	return "0"
}

func getSemverPatch(version string) string {
	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")
	if len(parts) > 2 {
		// Remove any pre-release or build metadata
		patch := strings.Split(parts[2], "-")[0]
		patch = strings.Split(patch, "+")[0]
		return patch
	}
	return "0"
}

func compareSemver(v1, v2 string) int {
	// Simple semver comparison - returns -1, 0, or 1
	v1Clean := strings.TrimPrefix(v1, "v")
	v2Clean := strings.TrimPrefix(v2, "v")

	if v1Clean == v2Clean {
		return 0
	}
	if v1Clean < v2Clean {
		return -1
	}
	return 1
}

func getLatestVersion(config *models.ProjectConfig, packageName string) string {
	if config.Versions != nil && config.Versions.Packages != nil {
		if version, exists := config.Versions.Packages[packageName]; exists {
			return version
		}
	}

	// Check core language/framework versions
	if config.Versions != nil {
		switch packageName {
		case "node", "nodejs":
			return config.Versions.Node
		case "go", "golang":
			return config.Versions.Go
		case "next", "nextjs":
			return config.Versions.Packages["next"]
		case "react":
			return config.Versions.Packages["react"]
		case "kotlin":
			return config.Versions.Packages["kotlin"]
		case "swift":
			return config.Versions.Packages["swift"]
		}
	}

	return constants.VersionLatest
}

func addVersionPrefix(version string) string {
	if !strings.HasPrefix(version, "v") {
		return "v" + version
	}
	return version
}

func stripVersionPrefix(version string) string {
	return strings.TrimPrefix(version, "v")
}

// Specific version getter functions for common languages/frameworks
func getNodeVersion(config *models.ProjectConfig) string {
	if config.Versions != nil && config.Versions.Node != "" {
		return config.Versions.Node
	}
	return "20.11.0" // Default fallback
}

func getGoVersion(config *models.ProjectConfig) string {
	if config.Versions != nil && config.Versions.Go != "" {
		return config.Versions.Go
	}
	return "1.22.0" // Default fallback
}

func getNextjsVersion(config *models.ProjectConfig) string {
	if config.Versions != nil && config.Versions.Packages["next"] != "" {
		return config.Versions.Packages["next"]
	}
	return "14.2.0" // Default fallback
}

func getReactVersion(config *models.ProjectConfig) string {
	if config.Versions != nil && config.Versions.Packages["react"] != "" {
		return config.Versions.Packages["react"]
	}
	return "^18.3.1" // Default fallback
}

func getKotlinVersion(config *models.ProjectConfig) string {
	if config.Versions != nil && config.Versions.Packages["kotlin"] != "" {
		return config.Versions.Packages["kotlin"]
	}
	return "2.0.0" // Default fallback
}

func getSwiftVersion(config *models.ProjectConfig) string {
	if config.Versions != nil && config.Versions.Packages["swift"] != "" {
		return config.Versions.Packages["swift"]
	}
	return "5.9.0" // Default fallback
}

func getPackageVersion(config *models.ProjectConfig, packageName string) string {
	if config.Versions != nil && config.Versions.Packages != nil {
		if version, exists := config.Versions.Packages[packageName]; exists {
			return version
		}
	}
	return constants.VersionLatest
}

func hasPackage(config *models.ProjectConfig, packageName string) bool {
	if config.Versions != nil && config.Versions.Packages != nil {
		_, exists := config.Versions.Packages[packageName]
		return exists
	}
	return false
}

// Conditional functions

func templateIf(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func templateIfNot(condition bool, trueVal, falseVal interface{}) interface{} {
	return templateIf(!condition, trueVal, falseVal)
}

func templateAnd(a, b bool) bool {
	return a && b
}

func templateOr(a, b bool) bool {
	return a || b
}

func templateNot(a bool) bool {
	return !a
}

func templateEq(a, b interface{}) bool {
	return a == b
}

func templateNe(a, b interface{}) bool {
	return a != b
}

func templateLt(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
}

func templateLe(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) <= fmt.Sprintf("%v", b)
}

func templateGt(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) > fmt.Sprintf("%v", b)
}

func templateGe(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) >= fmt.Sprintf("%v", b)
}

func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return false
	}
}

func isNonEmpty(value interface{}) bool {
	return !isEmpty(value)
}

// Component checking functions

func hasFrontendComponent(config *models.ProjectConfig) bool {
	return config.Components.Frontend.NextJS.App ||
		config.Components.Frontend.NextJS.Home ||
		config.Components.Frontend.NextJS.Admin ||
		config.Components.Frontend.NextJS.Shared
}

func hasBackendComponent(config *models.ProjectConfig) bool {
	return config.Components.Backend.GoGin
}

func hasMobileComponent(config *models.ProjectConfig) bool {
	return config.Components.Mobile.Android || config.Components.Mobile.IOS || config.Components.Mobile.Shared
}

func hasInfrastructureComponent(config *models.ProjectConfig) bool {
	return config.Components.Infrastructure.Docker ||
		config.Components.Infrastructure.Kubernetes ||
		config.Components.Infrastructure.Terraform
}

func hasComponent(config *models.ProjectConfig, componentType, componentName string) bool {
	switch componentType {
	case constants.TemplateFrontend:
		switch componentName {
		case "main_app":
			return config.Components.Frontend.NextJS.App
		case "home":
			return config.Components.Frontend.NextJS.Home
		case "admin":
			return config.Components.Frontend.NextJS.Admin
		case "shared":
			return config.Components.Frontend.NextJS.Shared
		}
	case constants.TemplateBackend:
		switch componentName {
		case "api":
			return config.Components.Backend.GoGin
		}
	case "mobile":
		switch componentName {
		case "android":
			return config.Components.Mobile.Android
		case "ios":
			return config.Components.Mobile.IOS
		case "shared":
			return config.Components.Mobile.Shared
		}
	case "infrastructure":
		switch componentName {
		case "docker":
			return config.Components.Infrastructure.Docker
		case "kubernetes":
			return config.Components.Infrastructure.Kubernetes
		case "terraform":
			return config.Components.Infrastructure.Terraform
		}
	}
	return false
}

// Utility functions

func formatTime(t time.Time, format string) string {
	return t.Format(format)
}

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func mul(a, b int) int {
	return a * b
}

func div(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

func mod(a, b int) int {
	if b == 0 {
		return 0
	}
	return a % b
}

func defaultValue(value, defaultVal interface{}) interface{} {
	if isEmpty(value) {
		return defaultVal
	}
	return value
}

func coalesce(values ...interface{}) interface{} {
	for _, value := range values {
		if !isEmpty(value) {
			return value
		}
	}
	return nil
}

func quote(s string) string {
	return strconv.Quote(s)
}

func singleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "\\'") + "'"
}

func indent(spaces int, text string) string {
	indentation := strings.Repeat(" ", spaces)
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indentation + line
		}
	}
	return strings.Join(lines, "\n")
}

func nindent(spaces int, text string) string {
	return "\n" + indent(spaces, text)
}

// Enhanced Node.js version functions

func getNodeRuntime(config *models.ProjectConfig) string {
	if config.Versions != nil && config.Versions.Node != "" {
		return config.Versions.Node
	}
	return "20.0.0" // Default fallback
}

func getNodeTypesVersion(config *models.ProjectConfig) string {
	if config.Versions != nil && config.Versions.Packages["@types/node"] != "" {
		return config.Versions.Packages["@types/node"]
	}
	return "^20.17.0" // Default fallback
}

func getNodeNPMVersion(config *models.ProjectConfig) string {
	if config.Versions != nil && config.Versions.Packages["npm"] != "" {
		return config.Versions.Packages["npm"]
	}
	return "10.0.0" // Default fallback
}

func getNodeDockerImage(config *models.ProjectConfig) string {
	if config.Versions != nil && config.Versions.Node != "" {
		return "node:" + config.Versions.Node + "-alpine"
	}
	return "node:20-alpine" // Default fallback
}

func isNodeLTS(config *models.ProjectConfig) bool {
	// Assume Node 20+ is LTS
	if config.Versions != nil && config.Versions.Node != "" {
		return strings.HasPrefix(config.Versions.Node, "20")
	}
	return true // Default to LTS
}

// getEnvVar returns an environment variable value or a default
func getEnvVar(key string, defaultValue ...string) string {
	// For templates, we return a placeholder that will be replaced during CI/CD
	// This allows the template to generate valid GitHub Actions syntax
	if len(defaultValue) > 0 {
		return "${{ env." + key + " || '" + defaultValue[0] + "' }}"
	}
	return "${{ env." + key + " }}"
}

// githubSecrets returns a GitHub Actions secrets expression
func githubSecrets(secretName string) string {
	return "${{ secrets." + secretName + " }}"
}

// githubMatrix returns a GitHub Actions matrix expression
func githubMatrix(field string) string {
	return "${{ matrix." + field + " }}"
}

// githubContext returns a GitHub Actions github context expression
func githubContext(field string) string {
	return "${{ github." + field + " }}"
}

// generateNonce generates a random nonce for Content Security Policy
func generateNonce() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Return a fallback nonce if random generation fails
		return "fallback-nonce"
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// getCustomVar retrieves a custom variable from the template context
// This function is used to access custom variables in templates
func getCustomVar(config interface{}, key string) string {
	if _, ok := config.(*models.ProjectConfig); ok {
		if value, exists := map[string]string{}[key]; exists {
			return value
		}
	}
	return ""
}

// sliceString slices a string from start to end index
// This function works with Go template pipes: {{.Name | slice 0 1}}
func sliceString(args ...interface{}) string {
	if len(args) < 3 {
		return ""
	}

	s, ok := args[0].(string)
	if !ok {
		return ""
	}

	start, ok := args[1].(int)
	if !ok {
		return ""
	}

	end, ok := args[2].(int)
	if !ok {
		return ""
	}

	if start < 0 || end < 0 || start >= len(s) || end > len(s) || start >= end {
		return ""
	}
	return s[start:end]
}
