package cleanup

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CodeAnalyzer provides utilities for analyzing Go source code
type CodeAnalyzer struct {
	fileSet *token.FileSet
}

// TODOItem represents a TODO/FIXME comment found in code
type TODOItem struct {
	File     string
	Line     int
	Type     string
	Message  string
	Context  string
	Priority Priority
	Category Category
}

// Priority represents the priority level of an issue
type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

// Category represents the category of an issue
type Category int

const (
	CategorySecurity Category = iota
	CategoryPerformance
	CategoryFeature
	CategoryBug
	CategoryDocumentation
	CategoryRefactor
)

// DuplicateCodeBlock represents a block of duplicate code
type DuplicateCodeBlock struct {
	Files      []string
	StartLines []int
	EndLines   []int
	Similarity float64
	Content    string
	Suggestion string
}

// UnusedCodeItem represents unused code elements
type UnusedCodeItem struct {
	File   string
	Line   int
	Type   string
	Name   string
	Reason string
}

// FunctionInfo holds detailed information about a function
type FunctionInfo struct {
	Name      string
	File      string
	StartLine int
	EndLine   int
	Node      *ast.FuncDecl
	Signature string
	Body      string
}

// Declaration represents a code declaration
type Declaration struct {
	Name     string
	Type     string // function, variable, type, const
	File     string
	Line     int
	Exported bool
	Package  string
}

// Usage represents a usage of a declaration
type Usage struct {
	File    string
	Line    int
	Context string
}

// ImportIssue represents import organization issues
type ImportIssue struct {
	File       string
	Line       int
	Type       string
	Import     string
	Suggestion string
}

// NewCodeAnalyzer creates a new code analyzer
func NewCodeAnalyzer() *CodeAnalyzer {
	return &CodeAnalyzer{
		fileSet: token.NewFileSet(),
	}
}

// AnalyzeTODOComments scans for TODO/FIXME comments in Go files
func (ca *CodeAnalyzer) AnalyzeTODOComments(rootDir string) ([]TODOItem, error) {
	var todos []TODOItem

	todoRegex := regexp.MustCompile(`(?i)(TODO|FIXME|HACK|XXX|BUG|NOTE)[\s:]*(.*)`)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		lines := strings.Split(string(content), "\n")
		for lineNum, line := range lines {
			if matches := todoRegex.FindStringSubmatch(line); matches != nil {
				todo := TODOItem{
					File:     path,
					Line:     lineNum + 1,
					Type:     strings.ToUpper(matches[1]),
					Message:  strings.TrimSpace(matches[2]),
					Context:  strings.TrimSpace(line),
					Priority: ca.determinePriority(matches[1], matches[2]),
					Category: ca.determineCategory(matches[2]),
				}
				todos = append(todos, todo)
			}
		}

		return nil
	})

	return todos, err
}

// FindDuplicateCode identifies duplicate code blocks using AST-based analysis
func (ca *CodeAnalyzer) FindDuplicateCode(rootDir string) ([]DuplicateCodeBlock, error) {
	var duplicates []DuplicateCodeBlock

	// Collect all functions and their AST representations
	functions := make(map[string]*FunctionInfo)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") || strings.Contains(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(ca.fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}

		ast.Inspect(file, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok && fn.Body != nil {
				pos := ca.fileSet.Position(fn.Pos())
				endPos := ca.fileSet.Position(fn.End())

				funcInfo := &FunctionInfo{
					Name:      fn.Name.Name,
					File:      path,
					StartLine: pos.Line,
					EndLine:   endPos.Line,
					Node:      fn,
					Signature: ca.getDetailedFunctionSignature(fn),
					Body:      ca.getFunctionBodyHash(fn),
				}

				key := fmt.Sprintf("%s:%d", path, pos.Line)
				functions[key] = funcInfo
			}
			return true
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Find duplicate functions by comparing AST structures
	duplicates = append(duplicates, ca.findDuplicateFunctions(functions)...)

	// Find duplicate validation logic
	duplicates = append(duplicates, ca.findDuplicateValidationLogic(rootDir)...)

	// Find duplicate helper functions
	duplicates = append(duplicates, ca.findDuplicateHelperFunctions(functions)...)

	return duplicates, nil
}

// IdentifyUnusedCode finds unused functions, variables, and imports
func (ca *CodeAnalyzer) IdentifyUnusedCode(rootDir string) ([]UnusedCodeItem, error) {
	var unused []UnusedCodeItem

	// Track all declarations and their usage across the project
	declarations := make(map[string]*Declaration)
	usages := make(map[string][]Usage)

	// First pass: collect all declarations
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") || strings.Contains(path, "_test.go") {
			return nil
		}

		file, err := parser.ParseFile(ca.fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}

		// Collect declarations
		ca.collectDeclarations(file, path, declarations)

		// Collect usages
		ca.collectUsages(file, path, usages)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Second pass: identify unused items
	for name, decl := range declarations {
		if ca.isUnused(name, decl, usages) {
			unused = append(unused, UnusedCodeItem{
				File:   decl.File,
				Line:   decl.Line,
				Type:   decl.Type,
				Name:   decl.Name,
				Reason: ca.getUnusedReason(decl, usages[name]),
			})
		}
	}

	// Also check for unused imports in each file
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") {
			return nil
		}

		file, err := parser.ParseFile(ca.fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return nil // Skip files that can't be parsed
		}

		// Check for unused imports
		for _, imp := range file.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if !ca.isImportUsed(file, importPath) {
				pos := ca.fileSet.Position(imp.Pos())
				unused = append(unused, UnusedCodeItem{
					File:   path,
					Line:   pos.Line,
					Type:   "import",
					Name:   importPath,
					Reason: "Import not used in file",
				})
			}
		}

		return nil
	})

	return unused, err
}

// ValidateImportOrganization checks import organization
func (ca *CodeAnalyzer) ValidateImportOrganization(rootDir string) ([]ImportIssue, error) {
	var issues []ImportIssue

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") {
			return nil
		}

		file, err := parser.ParseFile(ca.fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}

		// Check import organization
		if len(file.Imports) > 1 {
			for i, imp := range file.Imports {
				importPath := strings.Trim(imp.Path.Value, `"`)
				pos := ca.fileSet.Position(imp.Pos())

				// Check if imports are properly grouped
				if i > 0 {
					prevImport := strings.Trim(file.Imports[i-1].Path.Value, `"`)
					if ca.shouldBeGrouped(prevImport, importPath) {
						issues = append(issues, ImportIssue{
							File:       path,
							Line:       pos.Line,
							Type:       "grouping",
							Import:     importPath,
							Suggestion: "Group standard library, third-party, and local imports separately",
						})
					}
				}
			}
		}

		return nil
	})

	return issues, err
}

// Helper methods

func (ca *CodeAnalyzer) determinePriority(todoType, message string) Priority {
	message = strings.ToLower(message)
	todoType = strings.ToLower(todoType)

	if todoType == "fixme" || todoType == "bug" || strings.Contains(message, "security") {
		return PriorityHigh
	}
	if todoType == "hack" || strings.Contains(message, "performance") {
		return PriorityMedium
	}
	return PriorityLow
}

func (ca *CodeAnalyzer) determineCategory(message string) Category {
	message = strings.ToLower(message)

	if strings.Contains(message, "security") || strings.Contains(message, "vulnerability") {
		return CategorySecurity
	}
	if strings.Contains(message, "performance") || strings.Contains(message, "optimize") {
		return CategoryPerformance
	}
	if strings.Contains(message, "doc") || strings.Contains(message, "comment") {
		return CategoryDocumentation
	}
	if strings.Contains(message, "refactor") || strings.Contains(message, "cleanup") {
		return CategoryRefactor
	}
	if strings.Contains(message, "bug") || strings.Contains(message, "fix") {
		return CategoryBug
	}

	return CategoryFeature
}

func (ca *CodeAnalyzer) getFunctionSignature(fn *ast.FuncDecl) string {
	var params []string
	if fn.Type.Params != nil {
		for _, param := range fn.Type.Params.List {
			for range param.Names {
				params = append(params, "param")
			}
		}
	}

	return fmt.Sprintf("%s(%s)", fn.Name.Name, strings.Join(params, ","))
}

func (ca *CodeAnalyzer) isImportUsed(file *ast.File, importPath string) bool {
	// Simplified check - would need more sophisticated analysis
	packageName := filepath.Base(importPath)

	used := false
	ast.Inspect(file, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			if strings.HasPrefix(ident.Name, packageName) {
				used = true
				return false
			}
		}
		return true
	})

	return used
}

func (ca *CodeAnalyzer) shouldBeGrouped(prev, current string) bool {
	// Standard library imports should come first
	prevIsStd := !strings.Contains(prev, ".")
	currentIsStd := !strings.Contains(current, ".")

	// If previous is standard and current is not, they should be grouped separately
	return prevIsStd && !currentIsStd
}

// getDetailedFunctionSignature creates a detailed signature including parameter types
func (ca *CodeAnalyzer) getDetailedFunctionSignature(fn *ast.FuncDecl) string {
	var params []string
	if fn.Type.Params != nil {
		for _, param := range fn.Type.Params.List {
			paramType := ca.nodeToString(param.Type)
			for _, name := range param.Names {
				params = append(params, fmt.Sprintf("%s %s", name.Name, paramType))
			}
		}
	}

	var results []string
	if fn.Type.Results != nil {
		for _, result := range fn.Type.Results.List {
			resultType := ca.nodeToString(result.Type)
			if len(result.Names) > 0 {
				for _, name := range result.Names {
					results = append(results, fmt.Sprintf("%s %s", name.Name, resultType))
				}
			} else {
				results = append(results, resultType)
			}
		}
	}

	signature := fmt.Sprintf("func %s(%s)", fn.Name.Name, strings.Join(params, ", "))
	if len(results) > 0 {
		if len(results) == 1 {
			signature += " " + results[0]
		} else {
			signature += " (" + strings.Join(results, ", ") + ")"
		}
	}

	return signature
}

// getFunctionBodyHash creates a normalized hash of the function body
func (ca *CodeAnalyzer) getFunctionBodyHash(fn *ast.FuncDecl) string {
	if fn.Body == nil {
		return ""
	}

	// Normalize the function body by removing variable names and focusing on structure
	normalized := ca.normalizeNode(fn.Body)
	return normalized
}

// normalizeNode converts an AST node to a normalized string representation
func (ca *CodeAnalyzer) normalizeNode(node ast.Node) string {
	if node == nil {
		return ""
	}

	switch n := node.(type) {
	case *ast.BlockStmt:
		var stmts []string
		for _, stmt := range n.List {
			stmts = append(stmts, ca.normalizeNode(stmt))
		}
		return "{" + strings.Join(stmts, ";") + "}"

	case *ast.IfStmt:
		cond := ca.normalizeNode(n.Cond)
		body := ca.normalizeNode(n.Body)
		result := "if(" + cond + ")" + body
		if n.Else != nil {
			result += "else" + ca.normalizeNode(n.Else)
		}
		return result

	case *ast.ForStmt:
		init := ca.normalizeNode(n.Init)
		cond := ca.normalizeNode(n.Cond)
		post := ca.normalizeNode(n.Post)
		body := ca.normalizeNode(n.Body)
		return fmt.Sprintf("for(%s;%s;%s)%s", init, cond, post, body)

	case *ast.ReturnStmt:
		var results []string
		for _, result := range n.Results {
			results = append(results, ca.normalizeNode(result))
		}
		return "return(" + strings.Join(results, ",") + ")"

	case *ast.AssignStmt:
		return "assign"

	case *ast.CallExpr:
		return "call"

	case *ast.BinaryExpr:
		return "binary"

	default:
		return "stmt"
	}
}

// nodeToString converts an AST node to its string representation
func (ca *CodeAnalyzer) nodeToString(node ast.Node) string {
	if node == nil {
		return ""
	}

	switch n := node.(type) {
	case *ast.Ident:
		return n.Name
	case *ast.SelectorExpr:
		return ca.nodeToString(n.X) + "." + n.Sel.Name
	case *ast.StarExpr:
		return "*" + ca.nodeToString(n.X)
	case *ast.ArrayType:
		return "[]" + ca.nodeToString(n.Elt)
	case *ast.MapType:
		return "map[" + ca.nodeToString(n.Key) + "]" + ca.nodeToString(n.Value)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return "unknown"
	}
}

// findDuplicateFunctions identifies functions with similar implementations
func (ca *CodeAnalyzer) findDuplicateFunctions(functions map[string]*FunctionInfo) []DuplicateCodeBlock {
	var duplicates []DuplicateCodeBlock

	// Group functions by normalized body
	bodyGroups := make(map[string][]*FunctionInfo)
	for _, funcInfo := range functions {
		if funcInfo.Body != "" && len(funcInfo.Body) > 20 { // Only consider substantial functions
			bodyGroups[funcInfo.Body] = append(bodyGroups[funcInfo.Body], funcInfo)
		}
	}

	// Find groups with multiple functions
	for _, funcs := range bodyGroups {
		if len(funcs) > 1 {
			var files []string
			var startLines []int
			var endLines []int

			for _, f := range funcs {
				files = append(files, f.File)
				startLines = append(startLines, f.StartLine)
				endLines = append(endLines, f.EndLine)
			}

			duplicate := DuplicateCodeBlock{
				Files:      files,
				StartLines: startLines,
				EndLines:   endLines,
				Similarity: 1.0,
				Content:    fmt.Sprintf("Duplicate function implementations: %s", funcs[0].Name),
				Suggestion: "Extract common functionality into a shared utility function",
			}
			duplicates = append(duplicates, duplicate)
		}
	}

	return duplicates
}

// findDuplicateValidationLogic identifies redundant validation patterns
func (ca *CodeAnalyzer) findDuplicateValidationLogic(rootDir string) []DuplicateCodeBlock {
	var duplicates []DuplicateCodeBlock

	validationPatterns := make(map[string][]string) // pattern -> files

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Look for common validation patterns
		patterns := []string{
			`strings\.TrimSpace\([^)]+\)\s*==\s*""`,
			`len\([^)]+\)\s*==\s*0`,
			`err\s*!=\s*nil`,
			`==\s*nil`,
		}

		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			if re.Match(content) {
				validationPatterns[pattern] = append(validationPatterns[pattern], path)
			}
		}

		return nil
	})

	if err != nil {
		return duplicates
	}

	// Find patterns that appear in multiple files
	for pattern, files := range validationPatterns {
		if len(files) > 2 { // Only report if pattern appears in more than 2 files
			duplicate := DuplicateCodeBlock{
				Files:      files,
				Similarity: 0.8,
				Content:    fmt.Sprintf("Repeated validation pattern: %s", pattern),
				Suggestion: "Create a shared validation utility function",
			}
			duplicates = append(duplicates, duplicate)
		}
	}

	return duplicates
}

// findDuplicateHelperFunctions identifies similar helper functions
func (ca *CodeAnalyzer) findDuplicateHelperFunctions(functions map[string]*FunctionInfo) []DuplicateCodeBlock {
	var duplicates []DuplicateCodeBlock

	// Group helper functions by name pattern
	helperGroups := make(map[string][]*FunctionInfo)

	for _, funcInfo := range functions {
		// Identify helper functions by common naming patterns
		name := strings.ToLower(funcInfo.Name)
		if strings.Contains(name, "helper") ||
			strings.Contains(name, "util") ||
			strings.Contains(name, "setup") ||
			strings.Contains(name, "cleanup") ||
			strings.HasPrefix(name, "create") ||
			strings.HasPrefix(name, "build") ||
			strings.HasPrefix(name, "make") {

			// Group by function name (ignoring case)
			key := strings.ToLower(funcInfo.Name)
			helperGroups[key] = append(helperGroups[key], funcInfo)
		}
	}

	// Find groups with similar functions in different packages
	for name, funcs := range helperGroups {
		if len(funcs) > 1 {
			// Check if they're in different packages
			packages := make(map[string]bool)
			for _, f := range funcs {
				pkg := filepath.Dir(f.File)
				packages[pkg] = true
			}

			if len(packages) > 1 {
				var files []string
				var startLines []int
				var endLines []int

				for _, f := range funcs {
					files = append(files, f.File)
					startLines = append(startLines, f.StartLine)
					endLines = append(endLines, f.EndLine)
				}

				duplicate := DuplicateCodeBlock{
					Files:      files,
					StartLines: startLines,
					EndLines:   endLines,
					Similarity: 0.9,
					Content:    fmt.Sprintf("Similar helper function: %s", name),
					Suggestion: "Consider consolidating into a shared utility package",
				}
				duplicates = append(duplicates, duplicate)
			}
		}
	}

	return duplicates
}

// collectDeclarations collects all declarations in a file
func (ca *CodeAnalyzer) collectDeclarations(file *ast.File, filePath string, declarations map[string]*Declaration) {
	packageName := file.Name.Name

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			if node.Name != nil {
				pos := ca.fileSet.Position(node.Pos())
				key := fmt.Sprintf("%s.%s", packageName, node.Name.Name)
				declarations[key] = &Declaration{
					Name:     node.Name.Name,
					Type:     "function",
					File:     filePath,
					Line:     pos.Line,
					Exported: ast.IsExported(node.Name.Name),
					Package:  packageName,
				}
			}
		case *ast.GenDecl:
			for _, spec := range node.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if name.Name != "_" {
							pos := ca.fileSet.Position(name.Pos())
							key := fmt.Sprintf("%s.%s", packageName, name.Name)
							declType := "variable"
							if node.Tok == token.CONST {
								declType = "constant"
							}
							declarations[key] = &Declaration{
								Name:     name.Name,
								Type:     declType,
								File:     filePath,
								Line:     pos.Line,
								Exported: ast.IsExported(name.Name),
								Package:  packageName,
							}
						}
					}
				case *ast.TypeSpec:
					if s.Name.Name != "_" {
						pos := ca.fileSet.Position(s.Pos())
						key := fmt.Sprintf("%s.%s", packageName, s.Name.Name)
						declarations[key] = &Declaration{
							Name:     s.Name.Name,
							Type:     "type",
							File:     filePath,
							Line:     pos.Line,
							Exported: ast.IsExported(s.Name.Name),
							Package:  packageName,
						}
					}
				}
			}
		}
		return true
	})
}

// collectUsages collects all usages of identifiers in a file
func (ca *CodeAnalyzer) collectUsages(file *ast.File, filePath string, usages map[string][]Usage) {
	packageName := file.Name.Name

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.Ident:
			if node.Name != "_" && node.Obj == nil { // External reference
				pos := ca.fileSet.Position(node.Pos())
				key := fmt.Sprintf("%s.%s", packageName, node.Name)
				usages[key] = append(usages[key], Usage{
					File:    filePath,
					Line:    pos.Line,
					Context: "identifier",
				})
			}
		case *ast.SelectorExpr:
			if ident, ok := node.X.(*ast.Ident); ok {
				pos := ca.fileSet.Position(node.Pos())
				key := fmt.Sprintf("%s.%s", ident.Name, node.Sel.Name)
				usages[key] = append(usages[key], Usage{
					File:    filePath,
					Line:    pos.Line,
					Context: "selector",
				})
			}
		}
		return true
	})
}

// isUnused determines if a declaration is unused
func (ca *CodeAnalyzer) isUnused(name string, decl *Declaration, usages map[string][]Usage) bool {
	// Don't remove exported functions/types (they might be used by external packages)
	if decl.Exported {
		return false
	}

	// Don't remove main functions or init functions
	if decl.Name == "main" || decl.Name == "init" {
		return false
	}

	// Don't remove test functions
	if strings.HasPrefix(decl.Name, "Test") || strings.HasPrefix(decl.Name, "Benchmark") || strings.HasPrefix(decl.Name, "Example") {
		return false
	}

	// Check if there are any usages
	usageList := usages[name]

	// If no usages found, it's unused
	if len(usageList) == 0 {
		return true
	}

	// If only used in the same file where it's declared, it might be unused
	// (this is a simplified heuristic)
	usedElsewhere := false
	for _, usage := range usageList {
		if usage.File != decl.File {
			usedElsewhere = true
			break
		}
	}

	return !usedElsewhere
}

// getUnusedReason provides a reason why something is considered unused
func (ca *CodeAnalyzer) getUnusedReason(decl *Declaration, usages []Usage) string {
	if len(usages) == 0 {
		return fmt.Sprintf("%s '%s' is never used", decl.Type, decl.Name)
	}

	return fmt.Sprintf("%s '%s' is only used in the same file where it's declared", decl.Type, decl.Name)
}
