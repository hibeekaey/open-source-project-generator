// Package quality provides code quality analysis functionality for the audit engine.
//
// This package contains analyzers for measuring code complexity, test coverage,
// and other quality metrics to help identify areas for improvement.
package quality

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Pre-compiled regular expressions for performance
var (
	controlFlowRegex = regexp.MustCompile(`\b(if|for|while|switch)\b`)
	logicalOpRegex   = regexp.MustCompile(`\b(&&|\|\|)\b`)
	funcDefRegex     = regexp.MustCompile(`\bfunc\b.*\{`)
)

// ComplexityAnalyzer analyzes code complexity metrics including cyclomatic
// complexity, cognitive complexity, and maintainability indices.
type ComplexityAnalyzer struct {
	maxComplexity int
	rules         []ComplexityRule
}

// ComplexityRule defines a rule for complexity analysis
type ComplexityRule struct {
	ID          string
	Name        string
	Pattern     string
	Weight      int
	Description string
}

// NewComplexityAnalyzer creates a new complexity analyzer instance
func NewComplexityAnalyzer() *ComplexityAnalyzer {
	return &ComplexityAnalyzer{
		maxComplexity: 10, // Default threshold
		rules:         getDefaultComplexityRules(),
	}
}

// MeasureComplexity measures code complexity for a project
func (ca *ComplexityAnalyzer) MeasureComplexity(path string) (*interfaces.ComplexityAnalysisResult, error) {
	if err := ca.projectExists(path); err != nil {
		return nil, err
	}

	result := &interfaces.ComplexityAnalysisResult{
		Files:     []interfaces.FileComplexity{},
		Functions: []interfaces.FunctionComplexity{},
		Summary: interfaces.ComplexityAnalysisSummary{
			TotalFiles:          0,
			TotalFunctions:      0,
			AverageComplexity:   0.0,
			HighComplexityFiles: 0,
			TechnicalDebtHours:  0.0,
		},
	}

	// Analyze complexity for different file types
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || ca.shouldSkipFile(filePath) {
			return nil
		}

		// Only analyze source code files
		if !ca.isSourceCodeFile(filePath) {
			return nil
		}

		fileComplexity, err := ca.analyzeFileComplexity(filePath)
		if err != nil {
			return nil // Continue on error
		}

		result.Files = append(result.Files, *fileComplexity)
		result.Summary.TotalFiles++

		if fileComplexity.CyclomaticComplexity > ca.maxComplexity {
			result.Summary.HighComplexityFiles++
		}

		// Add function complexities (simplified - would need proper parsing)
		functions := ca.extractFunctions(filePath)
		result.Functions = append(result.Functions, functions...)
		result.Summary.TotalFunctions += len(functions)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("complexity analysis failed: %w", err)
	}

	// Calculate summary metrics
	if result.Summary.TotalFiles > 0 {
		var totalComplexity float64
		var totalDebt float64

		for _, file := range result.Files {
			totalComplexity += float64(file.CyclomaticComplexity)

			// Estimate technical debt (simplified calculation)
			if file.CyclomaticComplexity > ca.maxComplexity {
				totalDebt += float64(file.CyclomaticComplexity-ca.maxComplexity) * 0.5 // 0.5 hours per excess complexity point
			}
		}

		result.Summary.AverageComplexity = totalComplexity / float64(result.Summary.TotalFiles)
		result.Summary.TechnicalDebtHours = totalDebt
	}

	return result, nil
}

// analyzeFileComplexity analyzes complexity for a single file
func (ca *ComplexityAnalyzer) analyzeFileComplexity(filePath string) (*interfaces.FileComplexity, error) {
	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	complexity := &interfaces.FileComplexity{
		Path:                 filePath,
		Lines:                len(lines),
		CyclomaticComplexity: ca.calculateCyclomaticComplexity(string(content)),
		CognitiveComplexity:  ca.calculateCognitiveComplexity(string(content)),
		Maintainability:      ca.calculateMaintainabilityIndex(string(content), len(lines)),
		TechnicalDebt:        ca.calculateTechnicalDebt(string(content)),
	}

	return complexity, nil
}

// calculateCyclomaticComplexity calculates cyclomatic complexity
func (ca *ComplexityAnalyzer) calculateCyclomaticComplexity(content string) int {
	// Base complexity is 1
	complexity := 1

	// Count decision points that increase complexity
	patterns := []string{
		`\bif\b`,        // if statements
		`\belse\s+if\b`, // else if statements
		`\bfor\b`,       // for loops
		`\bwhile\b`,     // while loops
		`\bswitch\b`,    // switch statements
		`\bcase\b`,      // case statements
		`\bcatch\b`,     // catch blocks
		`\b&&\b`,        // logical AND
		`\b\|\|\b`,      // logical OR
		`\?.*:`,         // ternary operators
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(content, -1)
		complexity += len(matches)
	}

	return complexity
}

// calculateCognitiveComplexity calculates cognitive complexity
func (ca *ComplexityAnalyzer) calculateCognitiveComplexity(content string) int {
	// Simplified cognitive complexity calculation
	// In practice, this would need proper AST parsing
	complexity := 0
	nestingLevel := 0

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Increase nesting for blocks
		if strings.Contains(line, "{") {
			nestingLevel++
		}
		if strings.Contains(line, "}") && nestingLevel > 0 {
			nestingLevel--
		}

		// Add complexity for control structures
		if controlFlowRegex.MatchString(line) {
			complexity += 1 + nestingLevel
		}

		// Add complexity for logical operators
		if logicalOpRegex.MatchString(line) {
			complexity += 1
		}

		// Add complexity for nested functions/methods
		if funcDefRegex.MatchString(line) && nestingLevel > 0 {
			complexity += nestingLevel
		}
	}

	return complexity
}

// calculateMaintainabilityIndex calculates maintainability index
func (ca *ComplexityAnalyzer) calculateMaintainabilityIndex(content string, lines int) float64 {
	// Simplified maintainability index calculation
	// Real implementation would use Halstead metrics

	cyclomaticComplexity := float64(ca.calculateCyclomaticComplexity(content))
	linesOfCode := float64(lines)

	// Simplified formula (not the actual MI formula)
	if linesOfCode == 0 {
		return 100.0
	}

	// Higher complexity and more lines reduce maintainability
	maintainability := 100.0 - (cyclomaticComplexity * 2) - (linesOfCode / 100)

	if maintainability < 0 {
		maintainability = 0
	}
	if maintainability > 100 {
		maintainability = 100
	}

	return maintainability
}

// calculateTechnicalDebt estimates technical debt
func (ca *ComplexityAnalyzer) calculateTechnicalDebt(content string) string {
	cyclomaticComplexity := ca.calculateCyclomaticComplexity(content)

	if cyclomaticComplexity <= 5 {
		return "low"
	} else if cyclomaticComplexity <= 10 {
		return "medium"
	} else if cyclomaticComplexity <= 20 {
		return "high"
	} else {
		return "very high"
	}
}

// extractFunctions extracts function complexity information
func (ca *ComplexityAnalyzer) extractFunctions(filePath string) []interfaces.FunctionComplexity {
	var functions []interfaces.FunctionComplexity

	// #nosec G304 - Audit tool legitimately reads files for analysis
	content, err := os.ReadFile(filePath)
	if err != nil {
		return functions
	}

	lines := strings.Split(string(content), "\n")

	// Simple function detection (would need proper parsing for accuracy)
	funcPattern := regexp.MustCompile(`^\s*(func|function|def|public|private|protected)?\s*(\w+)\s*\(([^)]*)\)`)

	for i, line := range lines {
		matches := funcPattern.FindStringSubmatch(line)
		if len(matches) >= 3 {
			funcName := matches[2]
			params := strings.Split(matches[3], ",")
			paramCount := len(params)
			if len(params) == 1 && strings.TrimSpace(params[0]) == "" {
				paramCount = 0
			}

			// Extract function body to calculate complexity
			funcBody := ca.extractFunctionBody(lines, i)

			functions = append(functions, interfaces.FunctionComplexity{
				Name:                 funcName,
				File:                 filePath,
				Line:                 i + 1,
				CyclomaticComplexity: ca.calculateCyclomaticComplexity(funcBody),
				CognitiveComplexity:  ca.calculateCognitiveComplexity(funcBody),
				Parameters:           paramCount,
				Lines:                strings.Count(funcBody, "\n") + 1,
			})
		}
	}

	return functions
}

// extractFunctionBody extracts the body of a function starting from a given line
func (ca *ComplexityAnalyzer) extractFunctionBody(lines []string, startLine int) string {
	var body strings.Builder
	braceCount := 0
	inFunction := false

	for i := startLine; i < len(lines); i++ {
		line := lines[i]
		body.WriteString(line + "\n")

		// Count braces to determine function boundaries
		for _, char := range line {
			switch char {
			case '{':
				braceCount++
				inFunction = true
			case '}':
				braceCount--
				if inFunction && braceCount == 0 {
					return body.String()
				}
			}
		}

		// Prevent infinite loops for malformed code
		if i-startLine > 1000 {
			break
		}
	}

	return body.String()
}

// projectExists checks if the project path exists and is a directory
func (ca *ComplexityAnalyzer) projectExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("project path does not exist: %s", path)
		}
		return fmt.Errorf("unable to access project path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("project path is not a directory: %s", path)
	}

	return nil
}

// shouldSkipFile determines if a file should be skipped during analysis
func (ca *ComplexityAnalyzer) shouldSkipFile(filePath string) bool {
	// Skip binary files, images, and other non-text files
	skipExtensions := []string{
		".exe", ".bin", ".dll", ".so", ".dylib",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg",
		".pdf", ".doc", ".docx", ".xls", ".xlsx",
		".zip", ".tar", ".gz", ".rar", ".7z",
		".mp3", ".mp4", ".avi", ".mov", ".wmv",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, skipExt := range skipExtensions {
		if ext == skipExt {
			return true
		}
	}

	// Skip hidden files and directories
	if strings.HasPrefix(filepath.Base(filePath), ".") {
		return true
	}

	// Skip common directories
	skipDirs := []string{
		"node_modules", "vendor", ".git", ".svn", ".hg",
		"build", "dist", "target", "bin", "obj",
	}

	for _, skipDir := range skipDirs {
		if strings.Contains(filePath, skipDir) {
			return true
		}
	}

	return false
}

// isSourceCodeFile determines if a file is a source code file
func (ca *ComplexityAnalyzer) isSourceCodeFile(filePath string) bool {
	sourceExtensions := []string{
		".go", ".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".c", ".cpp", ".h", ".hpp",
		".cs", ".php", ".rb", ".swift", ".kt", ".scala", ".rs", ".dart", ".vue",
		".html", ".css", ".scss", ".sass", ".less", ".sql", ".sh", ".bash", ".ps1",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, sourceExt := range sourceExtensions {
		if ext == sourceExt {
			return true
		}
	}

	return false
}

// getDefaultComplexityRules returns default complexity analysis rules
func getDefaultComplexityRules() []ComplexityRule {
	return []ComplexityRule{
		{
			ID:          "COMPLEX-001",
			Name:        "High Cyclomatic Complexity",
			Pattern:     `cyclomatic_complexity > 10`,
			Weight:      3,
			Description: "Function has high cyclomatic complexity",
		},
		{
			ID:          "COMPLEX-002",
			Name:        "High Cognitive Complexity",
			Pattern:     `cognitive_complexity > 15`,
			Weight:      3,
			Description: "Function has high cognitive complexity",
		},
		{
			ID:          "COMPLEX-003",
			Name:        "Long Function",
			Pattern:     `lines > 50`,
			Weight:      2,
			Description: "Function is too long",
		},
		{
			ID:          "COMPLEX-004",
			Name:        "Too Many Parameters",
			Pattern:     `parameters > 5`,
			Weight:      2,
			Description: "Function has too many parameters",
		},
		{
			ID:          "COMPLEX-005",
			Name:        "Low Maintainability",
			Pattern:     `maintainability < 20`,
			Weight:      3,
			Description: "File has low maintainability index",
		},
	}
}

// SetMaxComplexity sets the maximum allowed complexity threshold
func (ca *ComplexityAnalyzer) SetMaxComplexity(max int) {
	ca.maxComplexity = max
}

// GetMaxComplexity returns the current maximum complexity threshold
func (ca *ComplexityAnalyzer) GetMaxComplexity() int {
	return ca.maxComplexity
}

// SetRules sets the complexity analysis rules
func (ca *ComplexityAnalyzer) SetRules(rules []ComplexityRule) {
	ca.rules = rules
}

// GetRules returns the current complexity analysis rules
func (ca *ComplexityAnalyzer) GetRules() []ComplexityRule {
	return ca.rules
}
