package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// FunctionToPackage maps common Go functions to their required packages
var FunctionToPackage = map[string]string{
	// time package
	"time.Now":           "time",
	"time.Since":         "time",
	"time.Until":         "time",
	"time.Sleep":         "time",
	"time.Parse":         "time",
	"time.ParseDuration": "time",
	"time.Date":          "time",
	"time.Unix":          "time",
	"time.After":         "time",
	"time.Tick":          "time",
	"time.NewTimer":      "time",
	"time.NewTicker":     "time",

	// fmt package
	"fmt.Printf":  "fmt",
	"fmt.Sprintf": "fmt",
	"fmt.Fprintf": "fmt",
	"fmt.Print":   "fmt",
	"fmt.Println": "fmt",
	"fmt.Errorf":  "fmt",
	"fmt.Scanf":   "fmt",
	"fmt.Sscanf":  "fmt",
	"fmt.Fscanf":  "fmt",

	// strings package
	"strings.Contains":   "strings",
	"strings.HasPrefix":  "strings",
	"strings.HasSuffix":  "strings",
	"strings.Index":      "strings",
	"strings.LastIndex":  "strings",
	"strings.Replace":    "strings",
	"strings.ReplaceAll": "strings",
	"strings.Split":      "strings",
	"strings.Join":       "strings",
	"strings.ToLower":    "strings",
	"strings.ToUpper":    "strings",
	"strings.TrimSpace":  "strings",
	"strings.Trim":       "strings",
	"strings.TrimPrefix": "strings",
	"strings.TrimSuffix": "strings",
	"strings.Fields":     "strings",
	"strings.Compare":    "strings",
	"strings.EqualFold":  "strings",
	"strings.NewReader":  "strings",
	"strings.Builder":    "strings",

	// strconv package
	"strconv.Atoi":        "strconv",
	"strconv.Itoa":        "strconv",
	"strconv.ParseInt":    "strconv",
	"strconv.ParseFloat":  "strconv",
	"strconv.ParseBool":   "strconv",
	"strconv.FormatInt":   "strconv",
	"strconv.FormatFloat": "strconv",
	"strconv.FormatBool":  "strconv",
	"strconv.Quote":       "strconv",
	"strconv.Unquote":     "strconv",

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
	"os.Getwd":      "os",
	"os.Chdir":      "os",

	// io package
	"io.Copy":        "io",
	"io.CopyN":       "io",
	"io.ReadAll":     "io",
	"io.WriteString": "io",
	"io.EOF":         "io",
	"io.Reader":      "io",
	"io.Writer":      "io",
	"io.Closer":      "io",

	// ioutil package (deprecated but still used)
	"ioutil.ReadFile":  "io/ioutil",
	"ioutil.WriteFile": "io/ioutil",
	"ioutil.ReadAll":   "io/ioutil",
	"ioutil.TempDir":   "io/ioutil",
	"ioutil.TempFile":  "io/ioutil",

	// json package
	"json.Marshal":    "encoding/json",
	"json.Unmarshal":  "encoding/json",
	"json.NewEncoder": "encoding/json",
	"json.NewDecoder": "encoding/json",

	// http package
	"http.StatusOK":                  "net/http",
	"http.StatusCreated":             "net/http",
	"http.StatusBadRequest":          "net/http",
	"http.StatusUnauthorized":        "net/http",
	"http.StatusForbidden":           "net/http",
	"http.StatusNotFound":            "net/http",
	"http.StatusInternalServerError": "net/http",
	"http.Get":                       "net/http",
	"http.Post":                      "net/http",
	"http.NewRequest":                "net/http",
	"http.DefaultClient":             "net/http",
	"http.ListenAndServe":            "net/http",
	"http.HandleFunc":                "net/http",
	"http.Handler":                   "net/http",
	"http.HandlerFunc":               "net/http",

	// url package
	"url.Parse":       "net/url",
	"url.QueryEscape": "net/url",
	"url.Values":      "net/url",

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

	// sync package
	"sync.Mutex":     "sync",
	"sync.RWMutex":   "sync",
	"sync.WaitGroup": "sync",
	"sync.Once":      "sync",

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

	// math package
	"math.Abs":   "math",
	"math.Max":   "math",
	"math.Min":   "math",
	"math.Ceil":  "math",
	"math.Floor": "math",
	"math.Round": "math",
	"math.Sqrt":  "math",
	"math.Pow":   "math",

	// crypto packages
	"md5.Sum":       "crypto/md5",
	"sha1.Sum":      "crypto/sha1",
	"sha256.Sum256": "crypto/sha256",
	"rand.Read":     "crypto/rand",

	// base64 package
	"base64.StdEncoding": "encoding/base64",
	"base64.URLEncoding": "encoding/base64",
}

// MissingImport represents a missing import in a template file
type MissingImport struct {
	FilePath       string
	MissingPackage string
	UsedFunction   string
	LineNumber     int
	CurrentImports []string
}

// TemplateAnalysis represents the analysis result for a template file
type TemplateAnalysis struct {
	FilePath       string
	MissingImports []MissingImport
	CurrentImports []string
	UsedFunctions  []string
}

// ScanResult represents the overall scan results
type ScanResult struct {
	TotalFiles      int
	GoTemplateFiles int
	FilesWithIssues int
	Analyses        []TemplateAnalysis
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <templates-directory> [--json output.json]")
	}

	templatesDir := os.Args[1]
	var jsonOutput string

	// Check for JSON output flag
	if len(os.Args) >= 4 && os.Args[2] == "--json" {
		jsonOutput = os.Args[3]
	}

	fmt.Printf("Scanning template files in: %s\n", templatesDir)
	fmt.Println("=" + strings.Repeat("=", 50))

	result, err := scanTemplateDirectory(templatesDir)
	if err != nil {
		log.Fatal("Error scanning templates:", err)
	}

	printScanResults(result)

	// Generate JSON report if requested
	if jsonOutput != "" {
		report, err := generateJSONReport(result, templatesDir)
		if err != nil {
			log.Printf("Error generating JSON report: %v", err)
		} else {
			err = saveJSONReport(report, jsonOutput)
			if err != nil {
				log.Printf("Error saving JSON report: %v", err)
			} else {
				fmt.Printf("\n📄 JSON report saved to: %s\n", jsonOutput)
			}
		}
	}
}

// scanTemplateDirectory scans all .tmpl files in the given directory
func scanTemplateDirectory(dir string) (*ScanResult, error) {
	result := &ScanResult{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		result.TotalFiles++

		// Only process .tmpl files that appear to be Go files
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Check if it's a Go template by looking at the file extension before .tmpl
		if !strings.HasSuffix(strings.TrimSuffix(path, ".tmpl"), ".go") {
			return nil
		}

		result.GoTemplateFiles++

		analysis, err := analyzeTemplateFile(path)
		if err != nil {
			fmt.Printf("Warning: Could not analyze %s: %v\n", path, err)
			return nil
		}

		if len(analysis.MissingImports) > 0 {
			result.FilesWithIssues++
		}

		result.Analyses = append(result.Analyses, *analysis)
		return nil
	})

	return result, err
}

// analyzeTemplateFile analyzes a single template file for missing imports
func analyzeTemplateFile(filePath string) (*TemplateAnalysis, error) {
	analysis := &TemplateAnalysis{
		FilePath: filePath,
	}

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Extract current imports
	currentImports := extractImports(string(content))
	analysis.CurrentImports = currentImports

	// Find used functions and check for missing imports
	usedFunctions := extractUsedFunctions(string(content))
	analysis.UsedFunctions = usedFunctions

	// Check for missing imports
	for _, function := range usedFunctions {
		if requiredPackage, exists := FunctionToPackage[function]; exists {
			if !containsImport(currentImports, requiredPackage) {
				lineNumber := findFunctionLineNumber(string(content), function)
				missing := MissingImport{
					FilePath:       filePath,
					MissingPackage: requiredPackage,
					UsedFunction:   function,
					LineNumber:     lineNumber,
					CurrentImports: currentImports,
				}
				analysis.MissingImports = append(analysis.MissingImports, missing)
			}
		}
	}

	return analysis, nil
}

// extractImports extracts import statements from Go template content
func extractImports(content string) []string {
	var imports []string

	// Try to parse as Go code first (removing template syntax temporarily)
	cleanContent := removeTemplateSyntax(content)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", cleanContent, parser.ImportsOnly)
	if err == nil && node != nil {
		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			imports = append(imports, importPath)
		}
		return imports
	}

	// Fallback to regex-based extraction
	importRegex := regexp.MustCompile(`import\s*\(\s*([^)]+)\s*\)|import\s+"([^"]+)"`)
	matches := importRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if match[1] != "" {
			// Multi-line import block
			lines := strings.Split(match[1], "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "//") {
					// Extract import path from quotes
					quotedRegex := regexp.MustCompile(`"([^"]+)"`)
					if quotedMatch := quotedRegex.FindStringSubmatch(line); len(quotedMatch) > 1 {
						imports = append(imports, quotedMatch[1])
					}
				}
			}
		} else if match[2] != "" {
			// Single import
			imports = append(imports, match[2])
		}
	}

	return imports
}

// extractUsedFunctions extracts function calls from the content
func extractUsedFunctions(content string) []string {
	var functions []string
	functionMap := make(map[string]bool)

	// Look for function calls in the format package.Function()
	functionRegex := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*\.[A-Z][a-zA-Z0-9_]*)\s*\(`)
	matches := functionRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			function := match[1]
			if !functionMap[function] {
				functions = append(functions, function)
				functionMap[function] = true
			}
		}
	}

	// Also look for constants and types (like http.StatusOK)
	constantRegex := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*\.[A-Z][a-zA-Z0-9_]*)\b`)
	constantMatches := constantRegex.FindAllStringSubmatch(content, -1)

	for _, match := range constantMatches {
		if len(match) > 1 {
			constant := match[1]
			if !functionMap[constant] && !strings.Contains(constant, "(") {
				functions = append(functions, constant)
				functionMap[constant] = true
			}
		}
	}

	return functions
}

// removeTemplateSyntax removes Go template syntax to make content parseable
func removeTemplateSyntax(content string) string {
	// Remove template actions like {{.Name}}, {{range}}, etc.
	templateRegex := regexp.MustCompile(`\{\{[^}]*\}\}`)
	cleaned := templateRegex.ReplaceAllString(content, "placeholder")

	// Replace template variables in strings with placeholders
	stringTemplateRegex := regexp.MustCompile(`"[^"]*\{\{[^}]*\}\}[^"]*"`)
	cleaned = stringTemplateRegex.ReplaceAllString(cleaned, `"placeholder"`)

	return cleaned
}

// containsImport checks if an import is already present
func containsImport(imports []string, targetPackage string) bool {
	for _, imp := range imports {
		if imp == targetPackage {
			return true
		}
		// Handle aliased imports
		if strings.Contains(imp, " ") {
			parts := strings.Fields(imp)
			if len(parts) > 1 && parts[len(parts)-1] == targetPackage {
				return true
			}
		}
	}
	return false
}

// findFunctionLineNumber finds the line number where a function is used
func findFunctionLineNumber(content, function string) int {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, function) {
			return i + 1
		}
	}
	return 0
}

// printScanResults prints the scan results in a formatted way
func printScanResults(result *ScanResult) {
	fmt.Printf("\nScan Results:\n")
	fmt.Printf("Total files scanned: %d\n", result.TotalFiles)
	fmt.Printf("Go template files: %d\n", result.GoTemplateFiles)
	fmt.Printf("Files with missing imports: %d\n", result.FilesWithIssues)
	fmt.Println()

	if result.FilesWithIssues == 0 {
		fmt.Println("✅ No missing imports found!")
		return
	}

	fmt.Println("Missing Imports Report:")
	fmt.Println(strings.Repeat("-", 80))

	// Group by file
	for _, analysis := range result.Analyses {
		if len(analysis.MissingImports) == 0 {
			continue
		}

		fmt.Printf("\n📁 File: %s\n", analysis.FilePath)
		fmt.Printf("Current imports: %v\n", analysis.CurrentImports)

		// Group missing imports by package
		packageMap := make(map[string][]MissingImport)
		for _, missing := range analysis.MissingImports {
			packageMap[missing.MissingPackage] = append(packageMap[missing.MissingPackage], missing)
		}

		// Sort packages for consistent output
		var packages []string
		for pkg := range packageMap {
			packages = append(packages, pkg)
		}
		sort.Strings(packages)

		for _, pkg := range packages {
			fmt.Printf("  ❌ Missing import: \"%s\"\n", pkg)
			for _, missing := range packageMap[pkg] {
				fmt.Printf("     Used: %s (line %d)\n", missing.UsedFunction, missing.LineNumber)
			}
		}
	}

	// Summary by package
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("Summary by Package:")

	packageCount := make(map[string]int)
	for _, analysis := range result.Analyses {
		for _, missing := range analysis.MissingImports {
			packageCount[missing.MissingPackage]++
		}
	}

	var sortedPackages []string
	for pkg := range packageCount {
		sortedPackages = append(sortedPackages, pkg)
	}
	sort.Strings(sortedPackages)

	for _, pkg := range sortedPackages {
		fmt.Printf("  %s: %d occurrences\n", pkg, packageCount[pkg])
	}
}
