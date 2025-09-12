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

// APIAnalysis represents the analysis of API separation
type APIAnalysis struct {
	PublicAPIs  []string
	PrivateAPIs []string
	Issues      []string
	Warnings    []string
}

func main() {
	analysis := &APIAnalysis{}

	fmt.Println("=== API Separation Analysis ===")

	// Analyze pkg directory (public APIs)
	err := analyzeDirectory("pkg", analysis, true)
	if err != nil {
		fmt.Printf("Error analyzing pkg directory: %v\n", err)
		os.Exit(1)
	}

	// Analyze internal directory (private APIs)
	err = analyzeDirectory("internal", analysis, false)
	if err != nil {
		fmt.Printf("Error analyzing internal directory: %v\n", err)
		os.Exit(1)
	}

	// Print results
	printResults(analysis)

	if len(analysis.Issues) > 0 {
		os.Exit(1)
	}
}

func analyzeDirectory(dir string, analysis *APIAnalysis, isPublic bool) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		return analyzeFile(path, analysis, isPublic)
	})
}

func analyzeFile(filename string, analysis *APIAnalysis, isPublic bool) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Check package naming
	packageName := node.Name.Name
	if !isValidPackageName(packageName) {
		analysis.Issues = append(analysis.Issues,
			fmt.Sprintf("Invalid package name '%s' in %s", packageName, filename))
	}

	// Analyze exported vs unexported identifiers
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.IsExported() {
				if isPublic {
					analysis.PublicAPIs = append(analysis.PublicAPIs,
						fmt.Sprintf("%s.%s (func) in %s", packageName, x.Name.Name, filename))
				} else {
					analysis.Warnings = append(analysis.Warnings,
						fmt.Sprintf("Exported function %s.%s in internal package %s",
							packageName, x.Name.Name, filename))
				}
			}
		case *ast.TypeSpec:
			if x.Name.IsExported() {
				if isPublic {
					analysis.PublicAPIs = append(analysis.PublicAPIs,
						fmt.Sprintf("%s.%s (type) in %s", packageName, x.Name.Name, filename))
				} else {
					analysis.Warnings = append(analysis.Warnings,
						fmt.Sprintf("Exported type %s.%s in internal package %s",
							packageName, x.Name.Name, filename))
				}
			}
		case *ast.GenDecl:
			if x.Tok == token.VAR || x.Tok == token.CONST {
				for _, spec := range x.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range valueSpec.Names {
							if name.IsExported() {
								if isPublic {
									analysis.PublicAPIs = append(analysis.PublicAPIs,
										fmt.Sprintf("%s.%s (%s) in %s", packageName, name.Name,
											x.Tok.String(), filename))
								} else {
									analysis.Warnings = append(analysis.Warnings,
										fmt.Sprintf("Exported %s %s.%s in internal package %s",
											x.Tok.String(), packageName, name.Name, filename))
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	return nil
}

func isValidPackageName(name string) bool {
	if name == "" {
		return false
	}

	// Package names should be lowercase
	if strings.ToLower(name) != name {
		return false
	}

	// Should not contain underscores or hyphens
	if strings.Contains(name, "_") || strings.Contains(name, "-") {
		return false
	}

	// Should start with a letter
	if name[0] < 'a' || name[0] > 'z' {
		return false
	}

	return true
}

func printResults(analysis *APIAnalysis) {
	fmt.Printf("\nPublic APIs found: %d\n", len(analysis.PublicAPIs))
	for _, api := range analysis.PublicAPIs {
		fmt.Printf("  ✓ %s\n", api)
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
	fmt.Printf("- Public APIs: %d\n", len(analysis.PublicAPIs))
	fmt.Printf("- Warnings: %d\n", len(analysis.Warnings))
	fmt.Printf("- Issues: %d\n", len(analysis.Issues))
}
