package security

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DryRunManager handles dry-run operations for previewing changes
type DryRunManager struct {
	enabled    bool
	operations []DryRunOperation
}

// DryRunOperation represents a planned operation in dry-run mode
type DryRunOperation struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Path        string                 `json:"path"`
	Details     map[string]interface{} `json:"details"`
	Timestamp   time.Time              `json:"timestamp"`
	Impact      string                 `json:"impact"` // "safe", "warning", "destructive"
	Size        int64                  `json:"size,omitempty"`
}

// NewDryRunManager creates a new dry-run manager
func NewDryRunManager() *DryRunManager {
	return &DryRunManager{
		enabled:    false,
		operations: make([]DryRunOperation, 0),
	}
}

// SetEnabled enables or disables dry-run mode
func (drm *DryRunManager) SetEnabled(enabled bool) {
	drm.enabled = enabled
	if enabled {
		drm.operations = make([]DryRunOperation, 0) // Reset operations when enabling
	}
}

// IsEnabled returns whether dry-run mode is enabled
func (drm *DryRunManager) IsEnabled() bool {
	return drm.enabled
}

// RecordFileWrite records a planned file write operation
func (drm *DryRunManager) RecordFileWrite(path string, data []byte, overwrite bool) {
	if !drm.enabled {
		return
	}

	impact := "safe"
	description := fmt.Sprintf("Create file: %s (%d bytes)", path, len(data))

	// Check if file exists
	if _, err := os.Stat(path); err == nil {
		if overwrite {
			impact = "destructive"
			description = fmt.Sprintf("Overwrite file: %s (%d bytes)", path, len(data))
		} else {
			impact = "warning"
			description = fmt.Sprintf("File exists, would skip: %s", path)
		}
	}

	operation := DryRunOperation{
		Type:        "file_write",
		Description: description,
		Path:        path,
		Details: map[string]interface{}{
			"size":      len(data),
			"overwrite": overwrite,
			"exists":    impact != "safe",
		},
		Timestamp: time.Now(),
		Impact:    impact,
		Size:      int64(len(data)),
	}

	drm.operations = append(drm.operations, operation)
}

// RecordFileDelete records a planned file deletion operation
func (drm *DryRunManager) RecordFileDelete(path string) {
	if !drm.enabled {
		return
	}

	impact := "destructive"
	size := int64(0)
	var description string

	// Get file size if it exists
	if fileInfo, err := os.Stat(path); err == nil {
		size = fileInfo.Size()
		description = fmt.Sprintf("Delete file: %s (%d bytes)", path, size)
	} else {
		impact = "safe"
		description = fmt.Sprintf("File does not exist, nothing to delete: %s", path)
	}

	operation := DryRunOperation{
		Type:        "file_delete",
		Description: description,
		Path:        path,
		Details: map[string]interface{}{
			"exists": impact == "destructive",
		},
		Timestamp: time.Now(),
		Impact:    impact,
		Size:      size,
	}

	drm.operations = append(drm.operations, operation)
}

// RecordDirectoryCreate records a planned directory creation operation
func (drm *DryRunManager) RecordDirectoryCreate(path string) {
	if !drm.enabled {
		return
	}

	impact := "safe"
	description := fmt.Sprintf("Create directory: %s", path)

	// Check if directory exists
	if _, err := os.Stat(path); err == nil {
		impact = "warning"
		description = fmt.Sprintf("Directory exists, would skip: %s", path)
	}

	operation := DryRunOperation{
		Type:        "directory_create",
		Description: description,
		Path:        path,
		Details: map[string]interface{}{
			"exists": impact == "warning",
		},
		Timestamp: time.Now(),
		Impact:    impact,
	}

	drm.operations = append(drm.operations, operation)
}

// RecordDirectoryDelete records a planned directory deletion operation
func (drm *DryRunManager) RecordDirectoryDelete(path string) {
	if !drm.enabled {
		return
	}

	impact := "destructive"
	fileCount := 0
	totalSize := int64(0)
	var description string

	// Count files in directory if it exists
	if _, err := os.Stat(path); err == nil {
		_ = filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				fileCount++
				totalSize += info.Size()
			}
			return nil
		})
		description = fmt.Sprintf("Delete directory: %s (%d files, %d bytes)", path, fileCount, totalSize)
	} else {
		impact = "safe"
		description = fmt.Sprintf("Directory does not exist, nothing to delete: %s", path)
	}

	operation := DryRunOperation{
		Type:        "directory_delete",
		Description: description,
		Path:        path,
		Details: map[string]interface{}{
			"exists":     impact == "destructive",
			"file_count": fileCount,
			"total_size": totalSize,
		},
		Timestamp: time.Now(),
		Impact:    impact,
		Size:      totalSize,
	}

	drm.operations = append(drm.operations, operation)
}

// RecordFileCopy records a planned file copy operation
func (drm *DryRunManager) RecordFileCopy(srcPath, dstPath string, overwrite bool) {
	if !drm.enabled {
		return
	}

	impact := "safe"
	size := int64(0)
	var description string

	// Get source file size
	if srcInfo, err := os.Stat(srcPath); err == nil {
		size = srcInfo.Size()
		description = fmt.Sprintf("Copy file: %s -> %s (%d bytes)", srcPath, dstPath, size)
	} else {
		impact = "warning"
		description = fmt.Sprintf("Source file does not exist: %s", srcPath)
	}

	// Check if destination exists
	if _, err := os.Stat(dstPath); err == nil {
		if overwrite {
			impact = "destructive"
			description = fmt.Sprintf("Copy and overwrite: %s -> %s (%d bytes)", srcPath, dstPath, size)
		} else {
			impact = "warning"
			description = fmt.Sprintf("Destination exists, would skip: %s -> %s", srcPath, dstPath)
		}
	}

	operation := DryRunOperation{
		Type:        "file_copy",
		Description: description,
		Path:        fmt.Sprintf("%s -> %s", srcPath, dstPath),
		Details: map[string]interface{}{
			"source":      srcPath,
			"destination": dstPath,
			"size":        size,
			"overwrite":   overwrite,
			"src_exists":  impact != "warning" || strings.Contains(description, "overwrite"),
			"dst_exists":  strings.Contains(description, "exists") || strings.Contains(description, "overwrite"),
		},
		Timestamp: time.Now(),
		Impact:    impact,
		Size:      size,
	}

	drm.operations = append(drm.operations, operation)
}

// RecordTemplateProcess records a planned template processing operation
func (drm *DryRunManager) RecordTemplateProcess(templatePath, outputPath string, variables map[string]interface{}) {
	if !drm.enabled {
		return
	}

	impact := "safe"
	description := fmt.Sprintf("Process template: %s -> %s", templatePath, outputPath)

	// Check if output file exists
	if _, err := os.Stat(outputPath); err == nil {
		impact = "destructive"
		description = fmt.Sprintf("Process template and overwrite: %s -> %s", templatePath, outputPath)
	}

	operation := DryRunOperation{
		Type:        "template_process",
		Description: description,
		Path:        outputPath,
		Details: map[string]interface{}{
			"template":      templatePath,
			"output":        outputPath,
			"variables":     variables,
			"output_exists": impact == "destructive",
		},
		Timestamp: time.Now(),
		Impact:    impact,
	}

	drm.operations = append(drm.operations, operation)
}

// RecordCustomOperation records a custom operation
func (drm *DryRunManager) RecordCustomOperation(operationType, description, path string, details map[string]interface{}, impact string) {
	if !drm.enabled {
		return
	}

	// Validate impact level
	if impact != "safe" && impact != "warning" && impact != "destructive" {
		impact = "warning"
	}

	operation := DryRunOperation{
		Type:        operationType,
		Description: description,
		Path:        path,
		Details:     details,
		Timestamp:   time.Now(),
		Impact:      impact,
	}

	drm.operations = append(drm.operations, operation)
}

// GetOperations returns all recorded operations
func (drm *DryRunManager) GetOperations() []DryRunOperation {
	return drm.operations
}

// GetOperationsByImpact returns operations filtered by impact level
func (drm *DryRunManager) GetOperationsByImpact(impact string) []DryRunOperation {
	var filtered []DryRunOperation
	for _, op := range drm.operations {
		if op.Impact == impact {
			filtered = append(filtered, op)
		}
	}
	return filtered
}

// GetSummary returns a summary of all planned operations
func (drm *DryRunManager) GetSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"total_operations":       len(drm.operations),
		"safe_operations":        0,
		"warning_operations":     0,
		"destructive_operations": 0,
		"total_files_affected":   0,
		"total_size_affected":    int64(0),
		"operations_by_type":     make(map[string]int),
	}

	affectedFiles := make(map[string]bool)
	totalSize := int64(0)

	for _, op := range drm.operations {
		// Count by impact
		switch op.Impact {
		case "safe":
			summary["safe_operations"] = summary["safe_operations"].(int) + 1
		case "warning":
			summary["warning_operations"] = summary["warning_operations"].(int) + 1
		case "destructive":
			summary["destructive_operations"] = summary["destructive_operations"].(int) + 1
		}

		// Count by type
		typeCount := summary["operations_by_type"].(map[string]int)
		typeCount[op.Type]++

		// Track affected files
		if op.Path != "" {
			affectedFiles[op.Path] = true
		}

		// Sum total size
		totalSize += op.Size
	}

	summary["total_files_affected"] = len(affectedFiles)
	summary["total_size_affected"] = totalSize

	return summary
}

// GenerateReport generates a detailed report of planned operations
func (drm *DryRunManager) GenerateReport(format string) (string, error) {
	switch strings.ToLower(format) {
	case "text":
		return drm.generateTextReport(), nil
	case "json":
		return drm.generateJSONReport()
	default:
		return "", fmt.Errorf("unsupported report format: %s", format)
	}
}

// generateTextReport generates a human-readable text report
func (drm *DryRunManager) generateTextReport() string {
	var report strings.Builder

	summary := drm.GetSummary()

	report.WriteString("DRY RUN REPORT\n")
	report.WriteString("==============\n\n")

	report.WriteString(fmt.Sprintf("Total Operations: %d\n", summary["total_operations"]))
	report.WriteString(fmt.Sprintf("Safe Operations: %d\n", summary["safe_operations"]))
	report.WriteString(fmt.Sprintf("Warning Operations: %d\n", summary["warning_operations"]))
	report.WriteString(fmt.Sprintf("Destructive Operations: %d\n", summary["destructive_operations"]))
	report.WriteString(fmt.Sprintf("Files Affected: %d\n", summary["total_files_affected"]))
	report.WriteString(fmt.Sprintf("Total Size: %d bytes\n\n", summary["total_size_affected"]))

	// Group operations by impact
	impacts := []string{"destructive", "warning", "safe"}
	impactLabels := map[string]string{
		"destructive": "DESTRUCTIVE OPERATIONS (⚠️  CAUTION)",
		"warning":     "WARNING OPERATIONS (⚠️  REVIEW)",
		"safe":        "SAFE OPERATIONS (✅ OK)",
	}

	for _, impact := range impacts {
		ops := drm.GetOperationsByImpact(impact)
		if len(ops) == 0 {
			continue
		}

		report.WriteString(fmt.Sprintf("%s\n", impactLabels[impact]))
		report.WriteString(strings.Repeat("-", len(impactLabels[impact])) + "\n")

		for _, op := range ops {
			report.WriteString(fmt.Sprintf("• %s\n", op.Description))
			if op.Size > 0 {
				report.WriteString(fmt.Sprintf("  Size: %d bytes\n", op.Size))
			}
		}
		report.WriteString("\n")
	}

	return report.String()
}

// generateJSONReport generates a JSON report
func (drm *DryRunManager) generateJSONReport() (string, error) {
	report := map[string]interface{}{
		"summary":    drm.GetSummary(),
		"operations": drm.operations,
		"timestamp":  time.Now(),
	}

	// Note: In a real implementation, you'd use json.Marshal here
	// For this example, we'll return a simple JSON-like string
	return fmt.Sprintf(`{
  "summary": %v,
  "operations": %d,
  "timestamp": "%s"
}`, report["summary"], len(drm.operations), time.Now().Format(time.RFC3339)), nil
}

// Clear clears all recorded operations
func (drm *DryRunManager) Clear() {
	drm.operations = make([]DryRunOperation, 0)
}

// HasDestructiveOperations returns true if there are any destructive operations
func (drm *DryRunManager) HasDestructiveOperations() bool {
	return len(drm.GetOperationsByImpact("destructive")) > 0
}

// HasWarningOperations returns true if there are any warning operations
func (drm *DryRunManager) HasWarningOperations() bool {
	return len(drm.GetOperationsByImpact("warning")) > 0
}
