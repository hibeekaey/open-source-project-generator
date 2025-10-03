package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// DeadCodeAnalyzer analyzes Go code for unused functions, methods, imports, etc.
type DeadCodeAnalyzer struct {
	fileSet     *token.FileSet
	packages    map[string]*ast.Package
	allFiles    []string
	results     *AnalysisResults
	projectRoot string
}

// AnalysisResults holds the results of dead code analysis
type AnalysisResults struct {
	UnusedFunctions    []UnusedFunction
	UnusedMethods      []UnusedMethod
	UnusedImports      []UnusedImport
	UnusedConstants    []UnusedConstant
	UnusedStructs      []UnusedStruct
	UnusedVariables    []UnusedVariable
	CommentedCode      []CommentedCodeBlock
	DuplicateFunctions []DuplicateFunction
}

// UnusedFunction represents an unused function
type UnusedFunction struct {
	Name       string
	File       string
	LineNumber int
	Exported   bool
	Receiver   string // For methods
}

// UnusedMethod represents an unused method
type UnusedMethod struct {
	Name       string
	File       string
	LineNumber int
	Receiver   string
	Exported   bool
}

// UnusedImport represents an unused import
type UnusedImport struct {
	Path       string
	Alias      string
	File       string
	LineNumber int
}

// UnusedConstant represents an unused constant
type UnusedConstant struct {
	Name       string
	File       string
	LineNumber int
	Exported   bool
}

// UnusedStruct represents an unused struct
type UnusedStruct struct {
	Name       string
	File       string
	LineNumber int
	Exported   bool
}

// UnusedVariable represents an unused variable
type UnusedVariable struct {
	Name       string
	File       string
	LineNumber int
	Exported   bool
}

// CommentedCodeBlock represents a block of commented-out code
type CommentedCodeBlock struct {
	File       string
	StartLine  int
	EndLine    int
	Content    string
	Confidence float64 // 0.0 to 1.0
}

// DuplicateFunction represents duplicate function implementations
type DuplicateFunction struct {
	Name      string
	Files     []string
	Signature string
}

// NewDeadCodeAnalyzer creates a new dead code analyzer
func NewDeadCodeAnalyzer(projectRoot string) *DeadCodeAnalyzer {
	return &DeadCodeAnalyzer{
		fileSet:     token.NewFileSet(),
		packages:    make(map[string]*ast.Package),
		projectRoot: projectRoot,
		results:     &AnalysisResults{},
	}
}

// Analyze performs comprehensive dead code analysis
func (dca *DeadCodeAnalyzer) Analyze() error {
	fmt.Println("Starting dead code analysis...")

	// Step 1: Discover all Go files
	if err := dca.discoverGoFiles(); err != nil {
		return fmt.Errorf("failed to discover Go files: %w", err)
	}

	// Step 2: Parse all Go files
	if err := dca.parseGoFiles(); err != nil {
		return fmt.Errorf("failed to parse Go files: %w", err)
	}

	// Step 3: Analyze for unused code
	dca.analyzeUnusedFunctions()
	dca.analyzeUnusedImports()
	dca.analyzeUnusedConstants()
	dca.analyzeUnusedStructs()
	dca.analyzeUnusedVariables()
	dca.analyzeCommentedCode()
	dca.analyzeDuplicateFunctions()

	fmt.Println("Dead code analysis completed.")
	return nil
}

// discoverGoFiles finds all Go source files in the project
func (dca *DeadCodeAnalyzer) discoverGoFiles() error {
	return filepath.WalkDir(dca.projectRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor, .git, and other non-source directories
		if d.IsDir() {
			name := d.Name()
			if name == "vendor" || name == ".git" || name == "node_modules" ||
				strings.HasPrefix(name, ".") && name != "." {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .go files, excluding test files for now
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			dca.allFiles = append(dca.allFiles, path)
		}

		return nil
	})
}

// parseGoFiles parses all discovered Go files
func (dca *DeadCodeAnalyzer) parseGoFiles() error {
	for _, file := range dca.allFiles {
		// #nosec G304 - This is a development tool that processes known Go files
		src, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Warning: failed to read file %s: %v\n", file, err)
			continue
		}

		// Parse the file
		astFile, err := parser.ParseFile(dca.fileSet, file, src, parser.ParseComments)
		if err != nil {
			fmt.Printf("Warning: failed to parse file %s: %v\n", file, err)
			continue
		}

		// Group by package
		pkgName := astFile.Name.Name
		if dca.packages[pkgName] == nil {
			dca.packages[pkgName] = &ast.Package{
				Name:  pkgName,
				Files: make(map[string]*ast.File),
			}
		}
		dca.packages[pkgName].Files[file] = astFile
	}

	return nil
}

// analyzeUnusedFunctions identifies unused functions and methods
func (dca *DeadCodeAnalyzer) analyzeUnusedFunctions() {
	fmt.Println("Analyzing unused functions...")

	// Collect all function definitions
	allFunctions := make(map[string]*FunctionInfo)

	for _, pkg := range dca.packages {
		for fileName, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.FuncDecl:
					if node.Name != nil {
						funcInfo := &FunctionInfo{
							Name:       node.Name.Name,
							File:       fileName,
							LineNumber: dca.fileSet.Position(node.Pos()).Line,
							Exported:   ast.IsExported(node.Name.Name),
							Node:       node,
						}

						// Check if it's a method
						if node.Recv != nil && len(node.Recv.List) > 0 {
							if recv := node.Recv.List[0]; recv.Type != nil {
								funcInfo.Receiver = dca.getTypeString(recv.Type)
							}
						}

						key := dca.getFunctionKey(funcInfo)
						allFunctions[key] = funcInfo
					}
				}
				return true
			})
		}
	}

	// Find usage of each function
	for _, funcInfo := range allFunctions {
		if !dca.isFunctionUsed(funcInfo, allFunctions) {
			if funcInfo.Receiver != "" {
				dca.results.UnusedMethods = append(dca.results.UnusedMethods, UnusedMethod{
					Name:       funcInfo.Name,
					File:       funcInfo.File,
					LineNumber: funcInfo.LineNumber,
					Receiver:   funcInfo.Receiver,
					Exported:   funcInfo.Exported,
				})
			} else {
				dca.results.UnusedFunctions = append(dca.results.UnusedFunctions, UnusedFunction{
					Name:       funcInfo.Name,
					File:       funcInfo.File,
					LineNumber: funcInfo.LineNumber,
					Exported:   funcInfo.Exported,
					Receiver:   funcInfo.Receiver,
				})
			}
		}
	}
}

// FunctionInfo holds information about a function
type FunctionInfo struct {
	Name       string
	File       string
	LineNumber int
	Exported   bool
	Receiver   string
	Node       *ast.FuncDecl
}

// getFunctionKey creates a unique key for a function
func (dca *DeadCodeAnalyzer) getFunctionKey(funcInfo *FunctionInfo) string {
	if funcInfo.Receiver != "" {
		return fmt.Sprintf("%s.%s", funcInfo.Receiver, funcInfo.Name)
	}
	return funcInfo.Name
}

// getTypeString extracts type string from AST node
func (dca *DeadCodeAnalyzer) getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + dca.getTypeString(t.X)
	case *ast.SelectorExpr:
		return dca.getTypeString(t.X) + "." + t.Sel.Name
	default:
		return "unknown"
	}
}

// isFunctionUsed checks if a function is used anywhere in the codebase
func (dca *DeadCodeAnalyzer) isFunctionUsed(funcInfo *FunctionInfo, allFunctions map[string]*FunctionInfo) bool {
	// Skip main functions and init functions
	if funcInfo.Name == "main" || funcInfo.Name == "init" {
		return true
	}

	// Skip exported functions for now (they might be used externally)
	if funcInfo.Exported {
		return true
	}

	// Skip interface implementations (check if it implements an interface)
	if dca.implementsInterface(funcInfo) {
		return true
	}

	// Search for usage in all files
	for _, pkg := range dca.packages {
		for _, file := range pkg.Files {
			if dca.findFunctionUsage(file, funcInfo) {
				return true
			}
		}
	}

	return false
}

// implementsInterface checks if a method implements an interface
func (dca *DeadCodeAnalyzer) implementsInterface(funcInfo *FunctionInfo) bool {
	if funcInfo.Receiver == "" {
		return false
	}

	// This is a simplified check - in a real implementation, you'd want to
	// check against actual interface definitions
	// For now, we'll be conservative and assume methods might implement interfaces
	return true
}

// findFunctionUsage searches for function usage in an AST file
func (dca *DeadCodeAnalyzer) findFunctionUsage(file *ast.File, funcInfo *FunctionInfo) bool {
	found := false

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if ident, ok := node.Fun.(*ast.Ident); ok {
				if ident.Name == funcInfo.Name {
					found = true
					return false
				}
			}
			if sel, ok := node.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == funcInfo.Name {
					found = true
					return false
				}
			}
		case *ast.Ident:
			if node.Name == funcInfo.Name {
				found = true
				return false
			}
		}
		return !found
	})

	return found
}

// analyzeUnusedImports identifies unused imports
func (dca *DeadCodeAnalyzer) analyzeUnusedImports() {
	fmt.Println("Analyzing unused imports...")

	for _, pkg := range dca.packages {
		for fileName, file := range pkg.Files {
			for _, imp := range file.Imports {
				importPath := strings.Trim(imp.Path.Value, `"`)
				alias := ""
				if imp.Name != nil {
					alias = imp.Name.Name
				}

				if !dca.isImportUsed(file, importPath, alias) {
					dca.results.UnusedImports = append(dca.results.UnusedImports, UnusedImport{
						Path:       importPath,
						Alias:      alias,
						File:       fileName,
						LineNumber: dca.fileSet.Position(imp.Pos()).Line,
					})
				}
			}
		}
	}
}

// isImportUsed checks if an import is used in the file
func (dca *DeadCodeAnalyzer) isImportUsed(file *ast.File, importPath, alias string) bool {
	// Get the package name from the import path
	pkgName := alias
	if pkgName == "" {
		parts := strings.Split(importPath, "/")
		pkgName = parts[len(parts)-1]
	}

	// Special case for dot imports
	if pkgName == "." {
		return true // Assume dot imports are used
	}

	// Search for usage of the package
	used := false
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := node.X.(*ast.Ident); ok {
				if ident.Name == pkgName {
					used = true
					return false
				}
			}
		}
		return !used
	})

	return used
}

// analyzeUnusedConstants identifies unused constants
func (dca *DeadCodeAnalyzer) analyzeUnusedConstants() {
	fmt.Println("Analyzing unused constants...")

	// Collect all constants
	allConstants := make(map[string]*ConstantInfo)

	for _, pkg := range dca.packages {
		for fileName, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.GenDecl:
					if node.Tok == token.CONST {
						for _, spec := range node.Specs {
							if valueSpec, ok := spec.(*ast.ValueSpec); ok {
								for _, name := range valueSpec.Names {
									constInfo := &ConstantInfo{
										Name:       name.Name,
										File:       fileName,
										LineNumber: dca.fileSet.Position(name.Pos()).Line,
										Exported:   ast.IsExported(name.Name),
									}
									allConstants[name.Name] = constInfo
								}
							}
						}
					}
				}
				return true
			})
		}
	}

	// Check usage
	for _, constInfo := range allConstants {
		if !dca.isConstantUsed(constInfo) {
			dca.results.UnusedConstants = append(dca.results.UnusedConstants, UnusedConstant{
				Name:       constInfo.Name,
				File:       constInfo.File,
				LineNumber: constInfo.LineNumber,
				Exported:   constInfo.Exported,
			})
		}
	}
}

// ConstantInfo holds information about a constant
type ConstantInfo struct {
	Name       string
	File       string
	LineNumber int
	Exported   bool
}

// isConstantUsed checks if a constant is used
func (dca *DeadCodeAnalyzer) isConstantUsed(constInfo *ConstantInfo) bool {
	// Skip exported constants
	if constInfo.Exported {
		return true
	}

	// Search for usage
	for _, pkg := range dca.packages {
		for _, file := range pkg.Files {
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				if ident, ok := n.(*ast.Ident); ok {
					if ident.Name == constInfo.Name {
						found = true
						return false
					}
				}
				return !found
			})
			if found {
				return true
			}
		}
	}

	return false
}

// analyzeUnusedStructs identifies unused struct types
func (dca *DeadCodeAnalyzer) analyzeUnusedStructs() {
	fmt.Println("Analyzing unused structs...")

	// Collect all struct types
	allStructs := make(map[string]*StructInfo)

	for _, pkg := range dca.packages {
		for fileName, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.GenDecl:
					if node.Tok == token.TYPE {
						for _, spec := range node.Specs {
							if typeSpec, ok := spec.(*ast.TypeSpec); ok {
								if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
									structInfo := &StructInfo{
										Name:       typeSpec.Name.Name,
										File:       fileName,
										LineNumber: dca.fileSet.Position(typeSpec.Pos()).Line,
										Exported:   ast.IsExported(typeSpec.Name.Name),
									}
									allStructs[typeSpec.Name.Name] = structInfo
								}
							}
						}
					}
				}
				return true
			})
		}
	}

	// Check usage
	for _, structInfo := range allStructs {
		if !dca.isStructUsed(structInfo) {
			dca.results.UnusedStructs = append(dca.results.UnusedStructs, UnusedStruct{
				Name:       structInfo.Name,
				File:       structInfo.File,
				LineNumber: structInfo.LineNumber,
				Exported:   structInfo.Exported,
			})
		}
	}
}

// StructInfo holds information about a struct
type StructInfo struct {
	Name       string
	File       string
	LineNumber int
	Exported   bool
}

// isStructUsed checks if a struct is used
func (dca *DeadCodeAnalyzer) isStructUsed(structInfo *StructInfo) bool {
	// Skip exported structs
	if structInfo.Exported {
		return true
	}

	// Search for usage
	for _, pkg := range dca.packages {
		for _, file := range pkg.Files {
			found := false
			ast.Inspect(file, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.Ident:
					if node.Name == structInfo.Name {
						found = true
						return false
					}
				case *ast.CompositeLit:
					if ident, ok := node.Type.(*ast.Ident); ok {
						if ident.Name == structInfo.Name {
							found = true
							return false
						}
					}
				}
				return !found
			})
			if found {
				return true
			}
		}
	}

	return false
}

// analyzeUnusedVariables identifies unused variables
func (dca *DeadCodeAnalyzer) analyzeUnusedVariables() {
	fmt.Println("Analyzing unused variables...")

	// This is a simplified implementation - a full implementation would need
	// to track variable scope and usage within functions
	for _, pkg := range dca.packages {
		for fileName, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				switch node := n.(type) {
				case *ast.GenDecl:
					if node.Tok == token.VAR {
						for _, spec := range node.Specs {
							if valueSpec, ok := spec.(*ast.ValueSpec); ok {
								for _, name := range valueSpec.Names {
									if name.Name != "_" && !ast.IsExported(name.Name) {
										// This is a simplified check
										dca.results.UnusedVariables = append(dca.results.UnusedVariables, UnusedVariable{
											Name:       name.Name,
											File:       fileName,
											LineNumber: dca.fileSet.Position(name.Pos()).Line,
											Exported:   ast.IsExported(name.Name),
										})
									}
								}
							}
						}
					}
				}
				return true
			})
		}
	}
}

// analyzeCommentedCode identifies blocks of commented-out code
func (dca *DeadCodeAnalyzer) analyzeCommentedCode() {
	fmt.Println("Analyzing commented-out code...")

	for _, file := range dca.allFiles {
		dca.findCommentedCodeInFile(file)
	}
}

// findCommentedCodeInFile finds commented-out code in a specific file
func (dca *DeadCodeAnalyzer) findCommentedCodeInFile(filename string) {
	// #nosec G304 - This is a development tool that processes known Go files
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	inCommentBlock := false
	blockStart := 0
	blockContent := []string{}

	// Patterns that suggest commented-out code
	codePatterns := []*regexp.Regexp{
		regexp.MustCompile(`^\s*//\s*(func\s+\w+|if\s+|for\s+|switch\s+|type\s+\w+|var\s+\w+|const\s+\w+)`),
		regexp.MustCompile(`^\s*//\s*[{}]\s*$`),
		regexp.MustCompile(`^\s*//\s*\w+\s*:=`),
		regexp.MustCompile(`^\s*//\s*return\s+`),
		regexp.MustCompile(`^\s*//\s*fmt\.`),
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check if this line looks like commented-out code
		isCodeComment := false
		for _, pattern := range codePatterns {
			if pattern.MatchString(line) {
				isCodeComment = true
				break
			}
		}

		if isCodeComment {
			if !inCommentBlock {
				inCommentBlock = true
				blockStart = lineNum
				blockContent = []string{}
			}
			blockContent = append(blockContent, line)
		} else {
			if inCommentBlock {
				// End of comment block
				if len(blockContent) >= 3 { // Only report blocks of 3+ lines
					confidence := dca.calculateCodeCommentConfidence(blockContent)
					if confidence > 0.6 { // Only report high-confidence matches
						dca.results.CommentedCode = append(dca.results.CommentedCode, CommentedCodeBlock{
							File:       filename,
							StartLine:  blockStart,
							EndLine:    lineNum - 1,
							Content:    strings.Join(blockContent, "\n"),
							Confidence: confidence,
						})
					}
				}
				inCommentBlock = false
			}
		}
	}

	// Handle case where file ends with commented code
	if inCommentBlock && len(blockContent) >= 3 {
		confidence := dca.calculateCodeCommentConfidence(blockContent)
		if confidence > 0.6 {
			dca.results.CommentedCode = append(dca.results.CommentedCode, CommentedCodeBlock{
				File:       filename,
				StartLine:  blockStart,
				EndLine:    lineNum,
				Content:    strings.Join(blockContent, "\n"),
				Confidence: confidence,
			})
		}
	}
}

// calculateCodeCommentConfidence calculates confidence that comments are actually code
func (dca *DeadCodeAnalyzer) calculateCodeCommentConfidence(lines []string) float64 {
	codeIndicators := 0
	totalLines := len(lines)

	for _, line := range lines {
		// Remove comment markers
		cleaned := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "//"))

		// Check for code-like patterns
		if strings.Contains(cleaned, "{") || strings.Contains(cleaned, "}") {
			codeIndicators++
		}
		if strings.Contains(cleaned, ":=") || strings.Contains(cleaned, "=") {
			codeIndicators++
		}
		if strings.HasPrefix(cleaned, "func ") || strings.HasPrefix(cleaned, "type ") {
			codeIndicators += 2
		}
		if strings.Contains(cleaned, "return ") || strings.Contains(cleaned, "if ") {
			codeIndicators++
		}
	}

	return float64(codeIndicators) / float64(totalLines*2) // Normalize to 0-1 range
}

// analyzeDuplicateFunctions identifies duplicate function implementations
func (dca *DeadCodeAnalyzer) analyzeDuplicateFunctions() {
	fmt.Println("Analyzing duplicate functions...")

	functionSignatures := make(map[string][]string)

	for _, pkg := range dca.packages {
		for fileName, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					if funcDecl.Name != nil {
						signature := dca.getFunctionSignature(funcDecl)
						key := funcDecl.Name.Name + signature
						functionSignatures[key] = append(functionSignatures[key], fileName)
					}
				}
				return true
			})
		}
	}

	// Find duplicates
	for signature, files := range functionSignatures {
		if len(files) > 1 {
			dca.results.DuplicateFunctions = append(dca.results.DuplicateFunctions, DuplicateFunction{
				Name:      strings.Split(signature, "(")[0],
				Files:     files,
				Signature: signature,
			})
		}
	}
}

// getFunctionSignature creates a signature string for a function
func (dca *DeadCodeAnalyzer) getFunctionSignature(funcDecl *ast.FuncDecl) string {
	var params []string
	if funcDecl.Type.Params != nil {
		for _, param := range funcDecl.Type.Params.List {
			paramType := dca.getTypeString(param.Type)
			params = append(params, paramType)
		}
	}

	var results []string
	if funcDecl.Type.Results != nil {
		for _, result := range funcDecl.Type.Results.List {
			resultType := dca.getTypeString(result.Type)
			results = append(results, resultType)
		}
	}

	signature := fmt.Sprintf("(%s)", strings.Join(params, ", "))
	if len(results) > 0 {
		signature += fmt.Sprintf(" (%s)", strings.Join(results, ", "))
	}

	return signature
}

// GenerateReport generates a comprehensive report of the analysis
func (dca *DeadCodeAnalyzer) GenerateReport() string {
	var report strings.Builder

	report.WriteString("# Dead Code Analysis Report\n\n")
	report.WriteString(fmt.Sprintf("Generated on: %s\n\n", fmt.Sprintf("%v", "analysis_time")))

	// Summary
	report.WriteString("## Summary\n\n")
	report.WriteString(fmt.Sprintf("- Unused Functions: %d\n", len(dca.results.UnusedFunctions)))
	report.WriteString(fmt.Sprintf("- Unused Methods: %d\n", len(dca.results.UnusedMethods)))
	report.WriteString(fmt.Sprintf("- Unused Imports: %d\n", len(dca.results.UnusedImports)))
	report.WriteString(fmt.Sprintf("- Unused Constants: %d\n", len(dca.results.UnusedConstants)))
	report.WriteString(fmt.Sprintf("- Unused Structs: %d\n", len(dca.results.UnusedStructs)))
	report.WriteString(fmt.Sprintf("- Unused Variables: %d\n", len(dca.results.UnusedVariables)))
	report.WriteString(fmt.Sprintf("- Commented Code Blocks: %d\n", len(dca.results.CommentedCode)))
	report.WriteString(fmt.Sprintf("- Duplicate Functions: %d\n", len(dca.results.DuplicateFunctions)))
	report.WriteString("\n")

	// Detailed sections
	if len(dca.results.UnusedFunctions) > 0 {
		report.WriteString("## Unused Functions\n\n")
		sort.Slice(dca.results.UnusedFunctions, func(i, j int) bool {
			return dca.results.UnusedFunctions[i].File < dca.results.UnusedFunctions[j].File
		})
		for _, fn := range dca.results.UnusedFunctions {
			report.WriteString(fmt.Sprintf("- `%s` in %s:%d (exported: %t)\n",
				fn.Name, fn.File, fn.LineNumber, fn.Exported))
		}
		report.WriteString("\n")
	}

	if len(dca.results.UnusedMethods) > 0 {
		report.WriteString("## Unused Methods\n\n")
		sort.Slice(dca.results.UnusedMethods, func(i, j int) bool {
			return dca.results.UnusedMethods[i].File < dca.results.UnusedMethods[j].File
		})
		for _, method := range dca.results.UnusedMethods {
			report.WriteString(fmt.Sprintf("- `%s.%s` in %s:%d (exported: %t)\n",
				method.Receiver, method.Name, method.File, method.LineNumber, method.Exported))
		}
		report.WriteString("\n")
	}

	if len(dca.results.UnusedImports) > 0 {
		report.WriteString("## Unused Imports\n\n")
		sort.Slice(dca.results.UnusedImports, func(i, j int) bool {
			return dca.results.UnusedImports[i].File < dca.results.UnusedImports[j].File
		})
		for _, imp := range dca.results.UnusedImports {
			alias := imp.Alias
			if alias == "" {
				alias = "default"
			}
			report.WriteString(fmt.Sprintf("- `%s` (as %s) in %s:%d\n",
				imp.Path, alias, imp.File, imp.LineNumber))
		}
		report.WriteString("\n")
	}

	if len(dca.results.CommentedCode) > 0 {
		report.WriteString("## Commented-Out Code Blocks\n\n")
		sort.Slice(dca.results.CommentedCode, func(i, j int) bool {
			return dca.results.CommentedCode[i].File < dca.results.CommentedCode[j].File
		})
		for _, block := range dca.results.CommentedCode {
			report.WriteString(fmt.Sprintf("- %s:%d-%d (confidence: %.2f)\n",
				block.File, block.StartLine, block.EndLine, block.Confidence))
		}
		report.WriteString("\n")
	}

	if len(dca.results.DuplicateFunctions) > 0 {
		report.WriteString("## Duplicate Functions\n\n")
		for _, dup := range dca.results.DuplicateFunctions {
			report.WriteString(fmt.Sprintf("- `%s` found in: %s\n",
				dup.Name, strings.Join(dup.Files, ", ")))
		}
		report.WriteString("\n")
	}

	return report.String()
}

// SaveReport saves the analysis report to a file
func (dca *DeadCodeAnalyzer) SaveReport(filename string) error {
	report := dca.GenerateReport()
	return os.WriteFile(filename, []byte(report), 0600)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run dead_code_analysis.go <project_root>")
		os.Exit(1)
	}

	projectRoot := os.Args[1]

	analyzer := NewDeadCodeAnalyzer(projectRoot)

	if err := analyzer.Analyze(); err != nil {
		fmt.Printf("Error during analysis: %v\n", err)
		os.Exit(1)
	}

	// Generate and save report
	reportFile := "dead_code_analysis_report.md"
	if err := analyzer.SaveReport(reportFile); err != nil {
		fmt.Printf("Error saving report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Analysis complete. Report saved to %s\n", reportFile)

	// Print summary to console
	fmt.Println("\n" + analyzer.GenerateReport())
}
