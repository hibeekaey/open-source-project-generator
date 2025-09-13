package cleanup

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// CodeConsolidator handles consolidation of duplicate code
type CodeConsolidator struct {
	fileSet     *token.FileSet
	analyzer    *CodeAnalyzer
	projectRoot string
}

// ConsolidationPlan represents a plan for consolidating duplicate code
type ConsolidationPlan struct {
	Duplicates      []DuplicateCodeBlock
	SharedUtilities []SharedUtility
	Refactorings    []Refactoring
	ValidationFixes []ValidationFix
}

// SharedUtility represents a utility function to be created
type SharedUtility struct {
	Name        string
	Package     string
	FilePath    string
	Function    string
	UsedBy      []string
	Description string
}

// Refactoring represents a code refactoring operation
type Refactoring struct {
	Type        string
	SourceFiles []string
	TargetFile  string
	Description string
}

// ValidationFix represents a validation pattern consolidation
type ValidationFix struct {
	Pattern     string
	Files       []string
	UtilityName string
	Description string
}

// NewCodeConsolidator creates a new code consolidator
func NewCodeConsolidator(projectRoot string) *CodeConsolidator {
	return &CodeConsolidator{
		fileSet:     token.NewFileSet(),
		analyzer:    NewCodeAnalyzer(),
		projectRoot: projectRoot,
	}
}

// CreateConsolidationPlan analyzes duplicates and creates a consolidation plan
func (cc *CodeConsolidator) CreateConsolidationPlan() (*ConsolidationPlan, error) {
	// Get duplicate code analysis
	duplicates, err := cc.analyzer.FindDuplicateCode(cc.projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze duplicates: %w", err)
	}

	plan := &ConsolidationPlan{
		Duplicates: duplicates,
	}

	// Create shared utilities for common patterns
	plan.SharedUtilities = cc.createSharedUtilities(duplicates)

	// Create refactoring plans
	plan.Refactorings = cc.createRefactoringPlans(duplicates)

	// Create validation fixes
	plan.ValidationFixes = cc.createValidationFixes(duplicates)

	return plan, nil
}

// createSharedUtilities identifies utilities that should be created
func (cc *CodeConsolidator) createSharedUtilities(duplicates []DuplicateCodeBlock) []SharedUtility {
	var utilities []SharedUtility

	for _, dup := range duplicates {
		if strings.Contains(dup.Content, "validation pattern") || strings.Contains(dup.Content, "TrimSpace") || strings.Contains(dup.Content, "len(") {
			// Create validation utilities
			if strings.Contains(dup.Content, "TrimSpace") {
				utilities = append(utilities, SharedUtility{
					Name:        "ValidateNonEmptyString",
					Package:     "pkg/utils",
					FilePath:    "pkg/utils/validation.go",
					Function:    cc.generateValidationFunction("ValidateNonEmptyString"),
					UsedBy:      dup.Files,
					Description: "Validates that a string is not empty after trimming whitespace",
				})
			}
			if strings.Contains(dup.Content, "len") && strings.Contains(dup.Content, "== 0") {
				utilities = append(utilities, SharedUtility{
					Name:        "ValidateNonEmptySlice",
					Package:     "pkg/utils",
					FilePath:    "pkg/utils/validation.go",
					Function:    cc.generateValidationFunction("ValidateNonEmptySlice"),
					UsedBy:      dup.Files,
					Description: "Validates that a slice is not empty",
				})
			}
			if strings.Contains(dup.Content, "err != nil") {
				utilities = append(utilities, SharedUtility{
					Name:        "HandleError",
					Package:     "pkg/utils",
					FilePath:    "pkg/utils/errors.go",
					Function:    cc.generateErrorHandlingFunction(),
					UsedBy:      dup.Files,
					Description: "Standard error handling utility",
				})
			}
		}

		// Create utilities for duplicate functions
		if strings.Contains(dup.Content, "GetLatestVersion") {
			utilities = append(utilities, SharedUtility{
				Name:        "GetLatestVersion",
				Package:     "pkg/version/common",
				FilePath:    "pkg/version/common/registry.go",
				Function:    cc.generateVersionFunction("GetLatestVersion"),
				UsedBy:      dup.Files,
				Description: "Common version retrieval logic for all registries",
			})
		}

		if strings.Contains(dup.Content, "CheckSecurity") {
			utilities = append(utilities, SharedUtility{
				Name:        "CheckSecurity",
				Package:     "pkg/version/common",
				FilePath:    "pkg/version/common/security.go",
				Function:    cc.generateSecurityCheckFunction(),
				UsedBy:      dup.Files,
				Description: "Common security checking logic for all registries",
			})
		}
	}

	return utilities
}

// createRefactoringPlans creates plans for refactoring duplicate code
func (cc *CodeConsolidator) createRefactoringPlans(duplicates []DuplicateCodeBlock) []Refactoring {
	var refactorings []Refactoring

	for _, dup := range duplicates {
		if len(dup.Files) > 1 && dup.Similarity > 0.9 {
			// High similarity functions should be consolidated
			refactorings = append(refactorings, Refactoring{
				Type:        "extract_function",
				SourceFiles: dup.Files,
				TargetFile:  cc.determineTargetFile(dup),
				Description: fmt.Sprintf("Extract common implementation: %s", dup.Content),
			})
		}
	}

	return refactorings
}

// createValidationFixes creates fixes for validation patterns
func (cc *CodeConsolidator) createValidationFixes(duplicates []DuplicateCodeBlock) []ValidationFix {
	var fixes []ValidationFix

	for _, dup := range duplicates {
		if strings.Contains(dup.Content, "validation pattern") {
			fixes = append(fixes, ValidationFix{
				Pattern:     dup.Content,
				Files:       dup.Files,
				UtilityName: cc.getValidationUtilityName(dup.Content),
				Description: dup.Suggestion,
			})
		}
	}

	return fixes
}

// ExecuteConsolidationPlan executes the consolidation plan
func (cc *CodeConsolidator) ExecuteConsolidationPlan(plan *ConsolidationPlan, dryRun bool) error {
	if dryRun {
		return cc.previewConsolidation(plan)
	}

	// Create shared utility files
	for _, utility := range plan.SharedUtilities {
		if err := cc.createSharedUtility(utility); err != nil {
			return fmt.Errorf("failed to create utility %s: %w", utility.Name, err)
		}
	}

	// Execute refactorings
	for _, refactoring := range plan.Refactorings {
		if err := cc.executeRefactoring(refactoring); err != nil {
			return fmt.Errorf("failed to execute refactoring: %w", err)
		}
	}

	return nil
}

// previewConsolidation shows what would be done without making changes
func (cc *CodeConsolidator) previewConsolidation(plan *ConsolidationPlan) error {
	fmt.Printf("=== Consolidation Plan Preview ===\n\n")

	fmt.Printf("Shared Utilities to Create (%d):\n", len(plan.SharedUtilities))
	for i, utility := range plan.SharedUtilities {
		fmt.Printf("%d. %s in %s\n", i+1, utility.Name, utility.FilePath)
		fmt.Printf("   Description: %s\n", utility.Description)
		fmt.Printf("   Used by %d files\n\n", len(utility.UsedBy))
	}

	fmt.Printf("Refactorings to Execute (%d):\n", len(plan.Refactorings))
	for i, refactoring := range plan.Refactorings {
		fmt.Printf("%d. %s\n", i+1, refactoring.Type)
		fmt.Printf("   Description: %s\n", refactoring.Description)
		fmt.Printf("   Affects %d files\n\n", len(refactoring.SourceFiles))
	}

	fmt.Printf("Validation Fixes (%d):\n", len(plan.ValidationFixes))
	for i, fix := range plan.ValidationFixes {
		fmt.Printf("%d. %s\n", i+1, fix.UtilityName)
		fmt.Printf("   Pattern: %s\n", fix.Pattern)
		fmt.Printf("   Affects %d files\n\n", len(fix.Files))
	}

	return nil
}

// createSharedUtility creates a shared utility file
func (cc *CodeConsolidator) createSharedUtility(utility SharedUtility) error {
	// Ensure directory exists
	dir := filepath.Dir(utility.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Check if file already exists
	if _, err := os.Stat(utility.FilePath); err == nil {
		// File exists, append to it
		return cc.appendToUtilityFile(utility)
	}

	// Create new file
	content := fmt.Sprintf(`package %s

// %s
%s
`, filepath.Base(utility.Package), utility.Description, utility.Function)

	return os.WriteFile(utility.FilePath, []byte(content), 0644)
}

// appendToUtilityFile adds a function to an existing utility file
func (cc *CodeConsolidator) appendToUtilityFile(utility SharedUtility) error {
	content, err := os.ReadFile(utility.FilePath)
	if err != nil {
		return err
	}

	// Parse existing file
	file, err := parser.ParseFile(cc.fileSet, utility.FilePath, content, parser.ParseComments)
	if err != nil {
		return err
	}

	// Check if function already exists
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == utility.Name {
				// Function already exists
				return nil
			}
		}
	}

	// Append new function
	newContent := string(content) + "\n" + utility.Function + "\n"
	return os.WriteFile(utility.FilePath, []byte(newContent), 0644)
}

// executeRefactoring executes a refactoring operation
func (cc *CodeConsolidator) executeRefactoring(refactoring Refactoring) error {
	// This is a simplified implementation
	// In a real scenario, this would involve complex AST manipulation
	fmt.Printf("Executing refactoring: %s\n", refactoring.Description)
	return nil
}

// Helper methods for generating utility functions

func (cc *CodeConsolidator) generateValidationFunction(name string) string {
	switch name {
	case "ValidateNonEmptyString":
		return `
// ValidateNonEmptyString validates that a string is not empty after trimming whitespace
func ValidateNonEmptyString(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}`
	case "ValidateNonEmptySlice":
		return `
// ValidateNonEmptySlice validates that a slice is not empty
func ValidateNonEmptySlice(slice interface{}, fieldName string) error {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("%s must be a slice", fieldName)
	}
	if v.Len() == 0 {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}`
	default:
		return ""
	}
}

func (cc *CodeConsolidator) generateErrorHandlingFunction() string {
	return `
// HandleError provides standard error handling with context
func HandleError(err error, context string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", context, err)
	}
	return nil
}`
}

func (cc *CodeConsolidator) generateVersionFunction(name string) string {
	return `
// GetLatestVersion retrieves the latest version from a registry
func GetLatestVersion(client RegistryClient, packageName string) (string, error) {
	versions, err := client.GetVersions(packageName)
	if err != nil {
		return "", fmt.Errorf("failed to get versions for %s: %w", packageName, err)
	}
	
	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found for %s", packageName)
	}
	
	// Return the latest version (assuming versions are sorted)
	return versions[0], nil
}`
}

func (cc *CodeConsolidator) generateSecurityCheckFunction() string {
	return `
// CheckSecurity performs security checks on a package version
func CheckSecurity(packageName, version string, vulnerabilityDB VulnerabilityDB) ([]SecurityIssue, error) {
	issues, err := vulnerabilityDB.CheckVulnerabilities(packageName, version)
	if err != nil {
		return nil, fmt.Errorf("failed to check vulnerabilities for %s@%s: %w", packageName, version, err)
	}
	
	return issues, nil
}`
}

func (cc *CodeConsolidator) determineTargetFile(dup DuplicateCodeBlock) string {
	// Simple heuristic: use the first file in pkg/ directory if available
	for _, file := range dup.Files {
		if strings.HasPrefix(file, "pkg/") {
			return filepath.Dir(file) + "/common.go"
		}
	}
	return "pkg/utils/common.go"
}

func (cc *CodeConsolidator) getValidationUtilityName(pattern string) string {
	if strings.Contains(pattern, "TrimSpace") {
		return "ValidateNonEmptyString"
	}
	if strings.Contains(pattern, "len") {
		return "ValidateNonEmptySlice"
	}
	if strings.Contains(pattern, "err") {
		return "HandleError"
	}
	return "ValidationUtility"
}
