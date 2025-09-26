package template

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// FunctionUsage represents a function call found in template code
type FunctionUsage struct {
	Function        string `json:"function"`
	Line            int    `json:"line"`
	Column          int    `json:"column"`
	RequiredPackage string `json:"required_package"`
}

// ImportStatement represents an import in a Go file
type ImportStatement struct {
	Package  string `json:"package"`
	Alias    string `json:"alias"`
	IsStdLib bool   `json:"is_stdlib"`
}

// MissingImportReport contains analysis results for a template file
type MissingImportReport struct {
	FilePath       string            `json:"file_path"`
	MissingImports []string          `json:"missing_imports"`
	UsedFunctions  []FunctionUsage   `json:"used_functions"`
	CurrentImports []ImportStatement `json:"current_imports"`
	Errors         []string          `json:"errors,omitempty"`
}

// ImportDetector analyzes Go template files for missing imports
type ImportDetector struct {
	functionPackageMap map[string]string
	fileSet            *token.FileSet
}

// NewImportDetector creates a new import detector with predefined function mappings
func NewImportDetector() *ImportDetector {
	return &ImportDetector{
		functionPackageMap: buildFunctionPackageMap(),
		fileSet:            token.NewFileSet(),
	}
}

// buildFunctionPackageMap creates a mapping of common functions to their required packages
func buildFunctionPackageMap() map[string]string {
	return map[string]string{
		// time package
		"time.Now":           "time",
		"time.Since":         "time",
		"time.Until":         "time",
		"time.Parse":         "time",
		"time.ParseDuration": "time",
		"time.Sleep":         "time",
		"time.Tick":          "time",
		"time.After":         "time",
		"time.NewTimer":      "time",
		"time.NewTicker":     "time",

		// fmt package
		"fmt.Printf":  "fmt",
		"fmt.Sprintf": "fmt",
		"fmt.Fprintf": "fmt",
		"fmt.Print":   "fmt",
		"fmt.Println": "fmt",
		"fmt.Errorf":  "fmt",
		"fmt.Scan":    "fmt",
		"fmt.Scanf":   "fmt",

		// strings package
		"strings.Contains":   "strings",
		"strings.HasPrefix":  "strings",
		"strings.HasSuffix":  "strings",
		"strings.Split":      "strings",
		"strings.Join":       "strings",
		"strings.Replace":    "strings",
		"strings.ReplaceAll": "strings",
		"strings.ToLower":    "strings",
		"strings.ToUpper":    "strings",
		"strings.TrimSpace":  "strings",
		"strings.Trim":       "strings",

		// strconv package
		"strconv.Atoi":        "strconv",
		"strconv.Itoa":        "strconv",
		"strconv.ParseInt":    "strconv",
		"strconv.ParseFloat":  "strconv",
		"strconv.ParseBool":   "strconv",
		"strconv.FormatInt":   "strconv",
		"strconv.FormatFloat": "strconv",
		"strconv.FormatBool":  "strconv",

		// os package
		"os.Getenv":     "os",
		"os.Setenv":     "os",
		"os.Exit":       "os",
		"os.Open":       "os",
		"os.Create":     "os",
		"os.Remove":     "os",
		"os.Mkdir":      "os",
		"os.MkdirAll":   "os",
		"os.Stat":       "os",
		"os.IsNotExist": "os",

		// log package
		"log.Printf":  "log",
		"log.Print":   "log",
		"log.Println": "log",
		"log.Fatal":   "log",
		"log.Fatalf":  "log",
		"log.Panic":   "log",
		"log.Panicf":  "log",

		// errors package
		"errors.New":    "errors",
		"errors.Is":     "errors",
		"errors.As":     "errors",
		"errors.Unwrap": "errors",

		// context package
		"context.Background":   "context",
		"context.TODO":         "context",
		"context.WithCancel":   "context",
		"context.WithTimeout":  "context",
		"context.WithDeadline": "context",
		"context.WithValue":    "context",

		// json package
		"json.Marshal":    "encoding/json",
		"json.Unmarshal":  "encoding/json",
		"json.NewEncoder": "encoding/json",
		"json.NewDecoder": "encoding/json",

		// http package
		"http.Get":                       "net/http",
		"http.Post":                      "net/http",
		"http.NewRequest":                "net/http",
		"http.ListenAndServe":            "net/http",
		"http.HandleFunc":                "net/http",
		"http.StatusOK":                  "net/http",
		"http.StatusNotFound":            "net/http",
		"http.StatusInternalServerError": "net/http",
		"http.StatusBadRequest":          "net/http",

		// url package
		"url.Parse":         "net/url",
		"url.QueryEscape":   "net/url",
		"url.QueryUnescape": "net/url",

		// filepath package
		"filepath.Join":  "path/filepath",
		"filepath.Dir":   "path/filepath",
		"filepath.Base":  "path/filepath",
		"filepath.Ext":   "path/filepath",
		"filepath.Clean": "path/filepath",
		"filepath.Abs":   "path/filepath",

		// io package
		"io.Copy":        "io",
		"io.ReadAll":     "io",
		"io.WriteString": "io",
		"io.EOF":         "io",

		// ioutil package (deprecated but still used)
		"ioutil.ReadFile":  "io/ioutil",
		"ioutil.WriteFile": "io/ioutil",
		"ioutil.ReadAll":   "io/ioutil",

		// regexp package
		"regexp.Compile":     "regexp",
		"regexp.MustCompile": "regexp",
		"regexp.Match":       "regexp",
		"regexp.MatchString": "regexp",

		// sort package
		"sort.Strings": "sort",
		"sort.Ints":    "sort",
		"sort.Sort":    "sort",
		"sort.Slice":   "sort",

		// sync package
		"sync.Mutex":     "sync",
		"sync.RWMutex":   "sync",
		"sync.WaitGroup": "sync",
		"sync.Once":      "sync",

		// crypto packages
		"md5.New":    "crypto/md5",
		"sha1.New":   "crypto/sha1",
		"sha256.New": "crypto/sha256",
		"rand.Read":  "crypto/rand",

		// base64 package
		"base64.StdEncoding": "encoding/base64",
		"base64.URLEncoding": "encoding/base64",

		// math package
		"math.Abs":  "math",
		"math.Max":  "math",
		"math.Min":  "math",
		"math.Sqrt": "math",
		"math.Pow":  "math",

		// reflect package
		"reflect.TypeOf":  "reflect",
		"reflect.ValueOf": "reflect",
	}
}

// AnalyzeTemplateFile analyzes a single template file for missing imports
func (id *ImportDetector) AnalyzeTemplateFile(filePath string) (*MissingImportReport, error) {
	report := &MissingImportReport{
		FilePath:       filePath,
		MissingImports: []string{},
		UsedFunctions:  []FunctionUsage{},
		CurrentImports: []ImportStatement{},
		Errors:         []string{},
	}

	// Read and preprocess template content
	content, readErr := id.readAndPreprocessTemplate(filePath)
	if readErr != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Failed to read file: %v", readErr))
		return report, readErr
	}

	// Parse the Go code
	file, err := parser.ParseFile(id.fileSet, filePath, content, parser.ParseComments)
	if err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Failed to parse Go code: %v", err))
		return report, err
	}

	// Extract current imports
	report.CurrentImports = id.extractImports(file)

	// Find function usages
	report.UsedFunctions = id.findFunctionUsages(file)

	// Determine missing imports
	report.MissingImports = id.findMissingImports(report.UsedFunctions, report.CurrentImports)

	return report, nil
}

// readAndPreprocessTemplate reads template file and preprocesses it for Go parsing
func (id *ImportDetector) readAndPreprocessTemplate(filePath string) (string, error) {
	content, err := readFileContent(filePath)
	if err != nil {
		return "", err
	}

	// Remove template directives that would break Go parsing
	// Replace template variables with valid Go identifiers
	processed := id.preprocessTemplateContent(content)

	return processed, nil
}

// preprocessTemplateContent removes/replaces template syntax that breaks Go parsing
func (id *ImportDetector) preprocessTemplateContent(content string) string {
	// Replace common template variables with valid Go identifiers
	replacements := map[string]string{
		`{{.Name}}`:        "TemplateName",
		`{{.ProjectName}}`: "TemplateProjectName",
		`{{.Package}}`:     "TemplatePackage",
		`{{.Version}}`:     "TemplateVersion",
		`{{.Author}}`:      "TemplateAuthor",
		`{{.Email}}`:       "TemplateEmail",
		`{{.Description}}`: "TemplateDescription",
	}

	processed := content
	for template, replacement := range replacements {
		processed = strings.ReplaceAll(processed, template, replacement)
	}

	// Remove template control structures that break parsing
	templateDirectives := []string{
		`{{- if .* -}}`,
		`{{- else -}}`,
		`{{- end -}}`,
		`{{- range .* -}}`,
		`{{ if .* }}`,
		`{{ else }}`,
		`{{ end }}`,
		`{{ range .* }}`,
	}

	for _, directive := range templateDirectives {
		re := regexp.MustCompile(directive)
		processed = re.ReplaceAllString(processed, "")
	}

	// Remove any remaining template expressions that might break parsing
	// We need to handle template expressions differently based on context
	processed = id.replaceTemplateExpressions(processed)

	return processed
}

// replaceTemplateExpressions replaces template expressions based on context
func (id *ImportDetector) replaceTemplateExpressions(content string) string {
	result := ""
	inString := false
	escapeNext := false
	i := 0

	for i < len(content) {
		if escapeNext {
			result += string(content[i])
			escapeNext = false
			i++
			continue
		}

		if content[i] == '\\' {
			result += string(content[i])
			escapeNext = true
			i++
			continue
		}

		if content[i] == '"' {
			inString = !inString
			result += string(content[i])
			i++
			continue
		}

		// Check for template expression start
		if i < len(content)-1 && content[i:i+2] == "{{" {
			// Find the end of the template expression
			end := i + 2
			for end < len(content)-1 && content[end:end+2] != "}}" {
				end++
			}
			if end < len(content)-1 {
				end += 2 // Include the closing }}

				// Replace based on context
				if inString {
					result += "template_value"
				} else {
					result += `"template_placeholder"`
				}
				i = end
				continue
			}
		}

		result += string(content[i])
		i++
	}

	return result
}

// extractImports extracts import statements from the parsed AST
func (id *ImportDetector) extractImports(file *ast.File) []ImportStatement {
	var imports []ImportStatement

	for _, imp := range file.Imports {
		importStmt := ImportStatement{
			Package:  strings.Trim(imp.Path.Value, `"`),
			IsStdLib: id.isStandardLibrary(strings.Trim(imp.Path.Value, `"`)),
		}

		if imp.Name != nil {
			importStmt.Alias = imp.Name.Name
		}

		imports = append(imports, importStmt)
	}

	return imports
}

// findFunctionUsages finds all function calls in the AST
func (id *ImportDetector) findFunctionUsages(file *ast.File) []FunctionUsage {
	var usages []FunctionUsage

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if funcName := id.extractFunctionName(node.Fun); funcName != "" {
				if requiredPackage, exists := id.functionPackageMap[funcName]; exists {
					pos := id.fileSet.Position(node.Pos())
					usage := FunctionUsage{
						Function:        funcName,
						Line:            pos.Line,
						Column:          pos.Column,
						RequiredPackage: requiredPackage,
					}
					usages = append(usages, usage)
				}
			}
		case *ast.SelectorExpr:
			// Handle cases like http.StatusOK (constants/variables)
			if funcName := id.extractSelectorName(node); funcName != "" {
				if requiredPackage, exists := id.functionPackageMap[funcName]; exists {
					pos := id.fileSet.Position(node.Pos())
					usage := FunctionUsage{
						Function:        funcName,
						Line:            pos.Line,
						Column:          pos.Column,
						RequiredPackage: requiredPackage,
					}
					usages = append(usages, usage)
				}
			}
		}
		return true
	})

	return usages
}

// extractFunctionName extracts function name from call expression
func (id *ImportDetector) extractFunctionName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		if ident, ok := e.X.(*ast.Ident); ok {
			return ident.Name + "." + e.Sel.Name
		}
	case *ast.Ident:
		return e.Name
	}
	return ""
}

// extractSelectorName extracts selector name for constants/variables
func (id *ImportDetector) extractSelectorName(expr *ast.SelectorExpr) string {
	if ident, ok := expr.X.(*ast.Ident); ok {
		return ident.Name + "." + expr.Sel.Name
	}
	return ""
}

// findMissingImports determines which imports are missing
func (id *ImportDetector) findMissingImports(usages []FunctionUsage, currentImports []ImportStatement) []string {
	requiredPackages := make(map[string]bool)
	currentPackages := make(map[string]bool)

	// Collect required packages
	for _, usage := range usages {
		requiredPackages[usage.RequiredPackage] = true
	}

	// Collect current packages
	for _, imp := range currentImports {
		currentPackages[imp.Package] = true
	}

	// Find missing packages
	var missing []string
	for pkg := range requiredPackages {
		if !currentPackages[pkg] {
			missing = append(missing, pkg)
		}
	}

	return missing
}

// isStandardLibrary checks if a package is part of Go standard library
func (id *ImportDetector) isStandardLibrary(pkg string) bool {
	stdLibPackages := map[string]bool{
		"bufio": true, "bytes": true, "context": true, "crypto": true,
		"database": true, "encoding": true, "errors": true, "fmt": true,
		"go": true, "hash": true, "html": true, "image": true, "io": true,
		"log": true, "math": true, "mime": true, "net": true, "os": true,
		"path": true, "reflect": true, "regexp": true, "runtime": true,
		"sort": true, "strconv": true, "strings": true, "sync": true,
		"syscall": true, "testing": true, "text": true, "time": true,
		"unicode": true, "unsafe": true,
	}

	// Check if it's a direct standard library package
	if stdLibPackages[pkg] {
		return true
	}

	// Check if it's a subpackage of standard library
	parts := strings.Split(pkg, "/")
	if len(parts) > 0 && stdLibPackages[parts[0]] {
		return true
	}

	return false
}

// readFileContent reads the content of a file
func readFileContent(filePath string) (string, error) {
	data, err := utils.SafeReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// AnalysisReport contains the complete analysis results for multiple files
type AnalysisReport struct {
	Reports     []MissingImportReport `json:"reports"`
	Summary     AnalysisSummary       `json:"summary"`
	GeneratedAt string                `json:"generated_at"`
}

// AnalysisSummary provides overview statistics
type AnalysisSummary struct {
	TotalFiles          int            `json:"total_files"`
	FilesWithIssues     int            `json:"files_with_issues"`
	TotalMissingImports int            `json:"total_missing_imports"`
	MostCommonMissing   map[string]int `json:"most_common_missing"`
}

// AnalyzeDirectory analyzes all template files in a directory recursively
func (id *ImportDetector) AnalyzeDirectory(dirPath string) (*AnalysisReport, error) {
	var reports []MissingImportReport

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only analyze .tmpl files
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		report, analyzeErr := id.AnalyzeTemplateFile(path)
		if analyzeErr != nil {
			// Create error report for files that couldn't be analyzed
			errorReport := MissingImportReport{
				FilePath: path,
				Errors:   []string{fmt.Sprintf("Analysis failed: %v", analyzeErr)},
			}
			reports = append(reports, errorReport)
		} else if report != nil {
			// Only include reports with issues or errors
			if len(report.MissingImports) > 0 || len(report.Errors) > 0 {
				reports = append(reports, *report)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	summary := id.generateSummary(reports)

	return &AnalysisReport{
		Reports:     reports,
		Summary:     summary,
		GeneratedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// generateSummary creates analysis summary statistics
func (id *ImportDetector) generateSummary(reports []MissingImportReport) AnalysisSummary {
	summary := AnalysisSummary{
		TotalFiles:          len(reports),
		FilesWithIssues:     0,
		TotalMissingImports: 0,
		MostCommonMissing:   make(map[string]int),
	}

	for _, report := range reports {
		if len(report.MissingImports) > 0 || len(report.Errors) > 0 {
			summary.FilesWithIssues++
		}

		summary.TotalMissingImports += len(report.MissingImports)

		for _, missing := range report.MissingImports {
			summary.MostCommonMissing[missing]++
		}
	}

	return summary
}

// GenerateTextReport creates a formatted text report from analysis results
func (id *ImportDetector) GenerateTextReport(analysis *AnalysisReport) string {
	var report strings.Builder

	report.WriteString("Template Import Analysis Report\n")
	report.WriteString("===============================\n\n")

	report.WriteString(fmt.Sprintf("Generated: %s\n", analysis.GeneratedAt))
	report.WriteString(fmt.Sprintf("Total Files Analyzed: %d\n", analysis.Summary.TotalFiles))
	report.WriteString(fmt.Sprintf("Files with Issues: %d\n", analysis.Summary.FilesWithIssues))
	report.WriteString(fmt.Sprintf("Total Missing Imports: %d\n\n", analysis.Summary.TotalMissingImports))

	if len(analysis.Summary.MostCommonMissing) > 0 {
		report.WriteString("Most Common Missing Imports:\n")
		for pkg, count := range analysis.Summary.MostCommonMissing {
			report.WriteString(fmt.Sprintf("  %s: %d files\n", pkg, count))
		}
		report.WriteString("\n")
	}

	if len(analysis.Reports) > 0 {
		report.WriteString("Detailed Results:\n")
		report.WriteString("-----------------\n\n")

		for _, fileReport := range analysis.Reports {
			report.WriteString(fmt.Sprintf("File: %s\n", fileReport.FilePath))

			if len(fileReport.Errors) > 0 {
				report.WriteString("  Errors:\n")
				for _, err := range fileReport.Errors {
					report.WriteString(fmt.Sprintf("    - %s\n", err))
				}
			}

			if len(fileReport.MissingImports) > 0 {
				report.WriteString("  Missing Imports:\n")
				for _, missing := range fileReport.MissingImports {
					report.WriteString(fmt.Sprintf("    - %s\n", missing))
				}
			}

			if len(fileReport.UsedFunctions) > 0 {
				report.WriteString("  Function Usage:\n")
				for _, usage := range fileReport.UsedFunctions {
					report.WriteString(fmt.Sprintf("    - %s (line %d) requires %s\n",
						usage.Function, usage.Line, usage.RequiredPackage))
				}
			}

			report.WriteString("\n")
		}
	} else {
		report.WriteString("No issues found in analyzed template files.\n")
	}

	return report.String()
}

// AddFunctionMapping adds a custom function to package mapping
func (id *ImportDetector) AddFunctionMapping(function, packageName string) {
	id.functionPackageMap[function] = packageName
}

// GetFunctionMappings returns a copy of the current function mappings
func (id *ImportDetector) GetFunctionMappings() map[string]string {
	mappings := make(map[string]string)
	for k, v := range id.functionPackageMap {
		mappings[k] = v
	}
	return mappings
}
