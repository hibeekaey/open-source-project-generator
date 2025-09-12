package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type ImportAnalysis struct {
	CircularDeps   []string
	UnusedImports  []string
	BadlyOrganized []string
	Issues         []string
	Warnings       []string
}

type ImportInfo struct {
	Path     string
	Name     string
	Position token.Pos
	Used     bool
}

func main() {
	analysis := &ImportAnalysis{}

	fmt.Println("=== Import Organization Analysis ===")

	// Analyze all Go files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Skip vendor and .git directories
		if strings.Contains(path, "vendor/") || strings.Contains(path, ".git/") {
			return nil
		}

		return analyzeFile(path, analysis)
	})

	if err != nil {
		fmt.Printf("Error analyzing files: %v\n", err)
		os.Exit(1)
	}

	// Check for circular dependencies using go list
	err = checkCircularDependencies(analysis)
	if err != nil {
		fmt.Printf("Error checking circular dependencies: %v\n", err)
	}

	// Print results
	printResults(analysis)

	if len(analysis.Issues) > 0 {
		os.Exit(1)
	}
}

func analyzeFile(filename string, analysis *ImportAnalysis) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Analyze imports
	imports := extractImports(node)

	// Check import organization
	if !areImportsWellOrganized(imports) {
		analysis.BadlyOrganized = append(analysis.BadlyOrganized, filename)
	}

	// Check for unused imports (basic check)
	unusedImports := findUnusedImports(node, imports)
	for _, imp := range unusedImports {
		analysis.UnusedImports = append(analysis.UnusedImports,
			fmt.Sprintf("%s: %s", filename, imp))
	}

	return nil
}

func extractImports(node *ast.File) []ImportInfo {
	var imports []ImportInfo

	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		importName := ""
		if imp.Name != nil {
			importName = imp.Name.Name
		}

		imports = append(imports, ImportInfo{
			Path:     importPath,
			Name:     importName,
			Position: imp.Pos(),
		})
	}

	return imports
}

func areImportsWellOrganized(imports []ImportInfo) bool {
	if len(imports) <= 1 {
		return true
	}

	// Check if imports are grouped properly:
	// 1. Standard library
	// 2. Third-party
	// 3. Local/internal

	var stdImports, thirdPartyImports, localImports []string

	for _, imp := range imports {
		if isStandardLibrary(imp.Path) {
			stdImports = append(stdImports, imp.Path)
		} else if isLocalImport(imp.Path) {
			localImports = append(localImports, imp.Path)
		} else {
			thirdPartyImports = append(thirdPartyImports, imp.Path)
		}
	}

	// Check if each group is sorted
	return isSorted(stdImports) && isSorted(thirdPartyImports) && isSorted(localImports)
}

func isStandardLibrary(importPath string) bool {
	// Basic check for standard library packages
	stdPkgs := []string{
		"bufio", "bytes", "context", "crypto", "database", "encoding", "errors",
		"fmt", "go", "hash", "html", "image", "io", "log", "math", "mime", "net",
		"os", "path", "reflect", "regexp", "runtime", "sort", "strconv", "strings",
		"sync", "syscall", "testing", "text", "time", "unicode", "unsafe",
	}

	for _, pkg := range stdPkgs {
		if strings.HasPrefix(importPath, pkg) {
			return true
		}
	}

	// Also check if it doesn't contain dots (simple heuristic)
	return !strings.Contains(importPath, ".")
}

func isLocalImport(importPath string) bool {
	// Check if it's a local import (starts with current module path or relative)
	return strings.HasPrefix(importPath, "./") ||
		strings.HasPrefix(importPath, "../") ||
		strings.Contains(importPath, "internal/") ||
		strings.Contains(importPath, "/pkg/")
}

func isSorted(slice []string) bool {
	for i := 1; i < len(slice); i++ {
		if slice[i-1] > slice[i] {
			return false
		}
	}
	return true
}

func findUnusedImports(node *ast.File, imports []ImportInfo) []string {
	var unused []string

	// This is a simplified check - in practice, you'd need more sophisticated analysis
	for _, imp := range imports {
		if imp.Name == "_" {
			// Blank imports are intentionally unused
			continue
		}

		// Extract package name from path
		pkgName := imp.Name
		if pkgName == "" {
			parts := strings.Split(imp.Path, "/")
			pkgName = parts[len(parts)-1]
		}

		// Check if package name is used in the file
		used := false
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectorExpr:
				if ident, ok := x.X.(*ast.Ident); ok {
					if ident.Name == pkgName {
						used = true
						return false
					}
				}
			case *ast.Ident:
				if x.Name == pkgName && x.Obj == nil {
					// This might be a reference to the imported package
					used = true
					return false
				}
			}
			return true
		})

		if !used {
			unused = append(unused, imp.Path)
		}
	}

	return unused
}

func checkCircularDependencies(analysis *ImportAnalysis) error {
	// This would require more sophisticated analysis
	// For now, we'll just note that this check should be done
	analysis.Warnings = append(analysis.Warnings,
		"Circular dependency check requires 'go list -deps' analysis")
	return nil
}

func printResults(analysis *ImportAnalysis) {
	fmt.Printf("\nCircular Dependencies: %d\n", len(analysis.CircularDeps))
	for _, dep := range analysis.CircularDeps {
		fmt.Printf("  ✗ %s\n", dep)
	}

	fmt.Printf("\nUnused Imports: %d\n", len(analysis.UnusedImports))
	for _, imp := range analysis.UnusedImports {
		fmt.Printf("  ⚠ %s\n", imp)
	}

	fmt.Printf("\nBadly Organized Imports: %d\n", len(analysis.BadlyOrganized))
	for _, file := range analysis.BadlyOrganized {
		fmt.Printf("  ⚠ %s\n", file)
	}

	fmt.Printf("\nWarnings: %d\n", len(analysis.Warnings))
	for _, warning := range analysis.Warnings {
		fmt.Printf("  ⚠ %s\n", warning)
	}

	fmt.Printf("\nIssues: %d\n", len(analysis.Issues))
	for _, issue := range analysis.Issues {
		fmt.Printf("  ✗ %s\n", issue)
	}

	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Circular Dependencies: %d\n", len(analysis.CircularDeps))
	fmt.Printf("- Unused Imports: %d\n", len(analysis.UnusedImports))
	fmt.Printf("- Badly Organized Files: %d\n", len(analysis.BadlyOrganized))
	fmt.Printf("- Warnings: %d\n", len(analysis.Warnings))
	fmt.Printf("- Issues: %d\n", len(analysis.Issues))
}
