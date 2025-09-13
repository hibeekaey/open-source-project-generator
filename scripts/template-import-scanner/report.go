package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// DetailedReport represents a comprehensive scan report
type DetailedReport struct {
	Timestamp       time.Time        `json:"timestamp"`
	ScanDirectory   string           `json:"scan_directory"`
	TotalFiles      int              `json:"total_files"`
	GoTemplateFiles int              `json:"go_template_files"`
	FilesWithIssues int              `json:"files_with_issues"`
	Issues          []FileIssue      `json:"issues"`
	PackageSummary  []PackageSummary `json:"package_summary"`
}

// FileIssue represents issues in a specific file
type FileIssue struct {
	FilePath       string          `json:"file_path"`
	CurrentImports []string        `json:"current_imports"`
	MissingImports []MissingImport `json:"missing_imports"`
}

// PackageSummary represents summary statistics for a package
type PackageSummary struct {
	Package     string `json:"package"`
	Occurrences int    `json:"occurrences"`
	FileCount   int    `json:"file_count"`
}

// generateJSONReport generates a JSON report from scan results
func generateJSONReport(result *ScanResult, directory string) (*DetailedReport, error) {
	report := &DetailedReport{
		Timestamp:       time.Now(),
		ScanDirectory:   directory,
		TotalFiles:      result.TotalFiles,
		GoTemplateFiles: result.GoTemplateFiles,
		FilesWithIssues: result.FilesWithIssues,
	}

	// Process issues
	for _, analysis := range result.Analyses {
		if len(analysis.MissingImports) > 0 {
			issue := FileIssue{
				FilePath:       analysis.FilePath,
				CurrentImports: analysis.CurrentImports,
				MissingImports: analysis.MissingImports,
			}
			report.Issues = append(report.Issues, issue)
		}
	}

	// Generate package summary
	packageCount := make(map[string]int)
	fileCount := make(map[string]int)

	for _, analysis := range result.Analyses {
		packageSet := make(map[string]bool)
		for _, missing := range analysis.MissingImports {
			packageCount[missing.MissingPackage]++
			if !packageSet[missing.MissingPackage] {
				fileCount[missing.MissingPackage]++
				packageSet[missing.MissingPackage] = true
			}
		}
	}

	for pkg, count := range packageCount {
		summary := PackageSummary{
			Package:     pkg,
			Occurrences: count,
			FileCount:   fileCount[pkg],
		}
		report.PackageSummary = append(report.PackageSummary, summary)
	}

	return report, nil
}

// saveJSONReport saves the report to a JSON file
func saveJSONReport(report *DetailedReport, filename string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
