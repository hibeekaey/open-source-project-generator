package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-source-template-generator/pkg/template"
)

// Config holds configuration for the import detector
type Config struct {
	TemplateDir string
	OutputFile  string
	Format      string // json, text, or summary
	Verbose     bool
}

// AnalysisResult holds the complete analysis results
type AnalysisResult struct {
	TotalFiles      int                             `json:"total_files"`
	FilesWithIssues int                             `json:"files_with_issues"`
	Reports         []*template.MissingImportReport `json:"reports"`
	Summary         map[string]int                  `json:"summary"` // package -> count of missing occurrences
}

func main() {
	config := parseFlags()

	if err := runAnalysis(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.TemplateDir, "dir", "templates", "Directory containing template files")
	flag.StringVar(&config.OutputFile, "output", "", "Output file (default: stdout)")
	flag.StringVar(&config.Format, "format", "text", "Output format: json, text, or summary")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nAnalyze Go template files for missing imports\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -dir templates -format json -output report.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -dir templates -format summary\n", os.Args[0])
	}

	flag.Parse()

	return config
}

func runAnalysis(config *Config) error {
	detector := template.NewImportDetector()

	// Find all template files
	templateFiles, err := findTemplateFiles(config.TemplateDir)
	if err != nil {
		return fmt.Errorf("failed to find template files: %w", err)
	}

	if config.Verbose {
		fmt.Fprintf(os.Stderr, "Found %d template files\n", len(templateFiles))
	}

	// Analyze each file
	result := &AnalysisResult{
		TotalFiles: len(templateFiles),
		Reports:    make([]*template.MissingImportReport, 0),
		Summary:    make(map[string]int),
	}

	for _, filePath := range templateFiles {
		if config.Verbose {
			fmt.Fprintf(os.Stderr, "Analyzing: %s\n", filePath)
		}

		report, err := detector.AnalyzeTemplateFile(filePath)
		if err != nil {
			if config.Verbose {
				fmt.Fprintf(os.Stderr, "Warning: Failed to analyze %s: %v\n", filePath, err)
			}
			continue
		}

		result.Reports = append(result.Reports, report)

		if len(report.MissingImports) > 0 || len(report.Errors) > 0 {
			result.FilesWithIssues++
		}

		// Update summary
		for _, pkg := range report.MissingImports {
			result.Summary[pkg]++
		}
	}

	// Output results
	return outputResults(result, config)
}

func findTemplateFiles(dir string) ([]string, error) {
	var templateFiles []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Check if it's a Go template file
		if strings.HasSuffix(path, ".go.tmpl") {
			templateFiles = append(templateFiles, path)
		}

		return nil
	})

	return templateFiles, err
}

func outputResults(result *AnalysisResult, config *Config) error {
	var output string
	var err error

	switch config.Format {
	case "json":
		output, err = formatJSON(result)
	case "summary":
		output = formatSummary(result)
	default:
		output = formatText(result)
	}

	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Write output
	if config.OutputFile != "" {
		return os.WriteFile(config.OutputFile, []byte(output), 0644)
	}

	fmt.Print(output)
	return nil
}

func formatJSON(result *AnalysisResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func formatSummary(result *AnalysisResult) string {
	var sb strings.Builder

	sb.WriteString("Import Detection Summary\n")
	sb.WriteString("=======================\n\n")
	sb.WriteString(fmt.Sprintf("Total files analyzed: %d\n", result.TotalFiles))
	sb.WriteString(fmt.Sprintf("Files with issues: %d\n", result.FilesWithIssues))
	sb.WriteString(fmt.Sprintf("Success rate: %.1f%%\n\n",
		float64(result.TotalFiles-result.FilesWithIssues)/float64(result.TotalFiles)*100))

	if len(result.Summary) > 0 {
		sb.WriteString("Most commonly missing packages:\n")
		for pkg, count := range result.Summary {
			sb.WriteString(fmt.Sprintf("  %-30s %d files\n", pkg, count))
		}
	} else {
		sb.WriteString("No missing imports detected!\n")
	}

	return sb.String()
}

func formatText(result *AnalysisResult) string {
	var sb strings.Builder

	sb.WriteString("Import Detection Report\n")
	sb.WriteString("======================\n\n")

	for _, report := range result.Reports {
		if len(report.MissingImports) == 0 && len(report.Errors) == 0 {
			continue // Skip files with no issues
		}

		sb.WriteString(fmt.Sprintf("File: %s\n", report.FilePath))

		if len(report.Errors) > 0 {
			sb.WriteString("  Errors:\n")
			for _, err := range report.Errors {
				sb.WriteString(fmt.Sprintf("    - %s\n", err))
			}
		}

		if len(report.MissingImports) > 0 {
			sb.WriteString("  Missing imports:\n")
			for _, pkg := range report.MissingImports {
				sb.WriteString(fmt.Sprintf("    - %s\n", pkg))
			}

			sb.WriteString("  Function usages requiring imports:\n")
			for _, usage := range report.UsedFunctions {
				found := false
				for _, missing := range report.MissingImports {
					if missing == usage.RequiredPackage {
						found = true
						break
					}
				}
				if found {
					sb.WriteString(fmt.Sprintf("    - %s (line %d) -> requires %s\n",
						usage.Function, usage.Line, usage.RequiredPackage))
				}
			}
		}

		sb.WriteString("\n")
	}

	// Summary
	sb.WriteString(fmt.Sprintf("Summary: %d/%d files have missing imports\n",
		result.FilesWithIssues, result.TotalFiles))

	return sb.String()
}
