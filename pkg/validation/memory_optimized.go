package validation

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/open-source-template-generator/pkg/models"
	"github.com/open-source-template-generator/pkg/utils"
)

// MemoryOptimizedValidator provides memory-efficient validation
type MemoryOptimizedValidator struct {
	bufferPool    *utils.MemoryPool
	stringBuilder *utils.StringBuilder
	readerPool    sync.Pool
}

// NewMemoryOptimizedValidator creates a new memory-optimized validator
func NewMemoryOptimizedValidator() *MemoryOptimizedValidator {
	return &MemoryOptimizedValidator{
		bufferPool:    utils.NewMemoryPool(),
		stringBuilder: utils.NewStringBuilder(),
		readerPool: sync.Pool{
			New: func() interface{} {
				return bufio.NewReaderSize(nil, 32*1024) // 32KB buffer
			},
		},
	}
}

// ValidateFileStreaming validates a file using streaming to reduce memory usage
func (mov *MemoryOptimizedValidator) ValidateFileStreaming(path string, validator func(line string, lineNum int) error) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close() // Ignore close error
	}()

	// Get a reader from the pool
	reader := mov.readerPool.Get().(*bufio.Reader)
	defer mov.readerPool.Put(reader)
	reader.Reset(file)

	lineNum := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read line %d: %w", lineNum, err)
		}

		if line != "" {
			lineNum++
			// Trim newline for processing
			line = strings.TrimSuffix(line, "\n")
			line = strings.TrimSuffix(line, "\r")

			if validationErr := validator(line, lineNum); validationErr != nil {
				return fmt.Errorf("validation failed at line %d: %w", lineNum, validationErr)
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

// ValidatePackageJSONStreaming validates package.json using streaming
func (mov *MemoryOptimizedValidator) ValidatePackageJSONStreaming(path string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Use a smaller buffer for JSON parsing to reduce memory usage
	buffer := mov.bufferPool.Get(8192) // 8KB buffer
	defer mov.bufferPool.Put(buffer)

	file, err := os.Open(path)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "file",
			Tag:     "access",
			Value:   path,
			Message: fmt.Sprintf("Failed to open file: %s", err.Error()),
		})
		return result, nil
	}
	defer func() {
		_ = file.Close() // Ignore close error
	}()

	// Read file in chunks to avoid loading entire file into memory
	var jsonContent strings.Builder
	reader := bufio.NewReaderSize(file, len(buffer))

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			jsonContent.Write(buffer[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   "file",
				Tag:     "read",
				Value:   path,
				Message: fmt.Sprintf("Failed to read file: %s", err.Error()),
			})
			return result, nil
		}
	}

	// Validate JSON structure without fully parsing into memory
	content := jsonContent.String()
	if err := mov.validateJSONStructure(content, result); err != nil {
		return result, err
	}

	return result, nil
}

// validateJSONStructure validates JSON structure efficiently
func (mov *MemoryOptimizedValidator) validateJSONStructure(content string, result *models.ValidationResult) error {
	// Basic JSON validation without full parsing
	content = strings.TrimSpace(content)

	if !strings.HasPrefix(content, "{") || !strings.HasSuffix(content, "}") {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "json_structure",
			Tag:     "format",
			Value:   "invalid",
			Message: "JSON must be an object starting with { and ending with }",
		})
		return nil
	}

	// Check for required fields using string operations instead of full JSON parsing
	requiredFields := []string{"name", "version"}
	for _, field := range requiredFields {
		fieldPattern := fmt.Sprintf(`"%s"`, field)
		if !strings.Contains(content, fieldPattern) {
			result.Valid = false
			result.Errors = append(result.Errors, models.ValidationError{
				Field:   field,
				Tag:     "required",
				Value:   "",
				Message: fmt.Sprintf("Required field '%s' is missing", field),
			})
		}
	}

	return nil
}

// ValidateGoModStreaming validates go.mod using streaming
func (mov *MemoryOptimizedValidator) ValidateGoModStreaming(path string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	hasModule := false
	hasGoVersion := false

	err := mov.ValidateFileStreaming(path, func(line string, lineNum int) error {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "//") {
			return nil
		}

		if strings.HasPrefix(line, "module ") {
			hasModule = true
			parts := strings.Fields(line)
			if len(parts) < 2 {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   "module",
					Tag:     "format",
					Value:   line,
					Message: fmt.Sprintf("Invalid module declaration at line %d", lineNum),
				})
			}
		}

		if strings.HasPrefix(line, "go ") {
			hasGoVersion = true
			parts := strings.Fields(line)
			if len(parts) < 2 {
				result.Valid = false
				result.Errors = append(result.Errors, models.ValidationError{
					Field:   "go_version",
					Tag:     "format",
					Value:   line,
					Message: fmt.Sprintf("Invalid go version declaration at line %d", lineNum),
				})
			}
		}

		return nil
	})

	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "file",
			Tag:     "read",
			Value:   path,
			Message: err.Error(),
		})
		return result, nil
	}

	// Check for required declarations
	if !hasModule {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "module",
			Tag:     "required",
			Value:   "",
			Message: "Missing module declaration",
		})
	}

	if !hasGoVersion {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "go_version",
			Tag:     "required",
			Value:   "",
			Message: "Missing go version declaration",
		})
	}

	return result, nil
}

// ValidateProjectMemoryEfficient validates a project with memory efficiency
func (mov *MemoryOptimizedValidator) ValidateProjectMemoryEfficient(projectPath string) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationWarning{},
	}

	// Use streaming directory walk to avoid loading all file info into memory
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Process files based on type without loading full content
		fileName := filepath.Base(path)
		relativePath, _ := filepath.Rel(projectPath, path)

		switch fileName {
		case "package.json":
			fileResult, err := mov.ValidatePackageJSONStreaming(path)
			if err != nil {
				return err
			}
			mov.mergeValidationResults(result, fileResult, relativePath)

		case "go.mod":
			fileResult, err := mov.ValidateGoModStreaming(path)
			if err != nil {
				return err
			}
			mov.mergeValidationResults(result, fileResult, relativePath)

		case "Dockerfile":
			// Validate Dockerfile using streaming
			if err := mov.validateDockerfileStreaming(path, result, relativePath); err != nil {
				return err
			}
		}

		// Force garbage collection periodically to keep memory usage low
		utils.ForceGlobalGC()

		return nil
	})

	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   "project_walk",
			Tag:     "access",
			Value:   projectPath,
			Message: fmt.Sprintf("Failed to walk project directory: %s", err.Error()),
		})
	}

	return result, nil
}

// validateDockerfileStreaming validates Dockerfile using streaming
func (mov *MemoryOptimizedValidator) validateDockerfileStreaming(path string, result *models.ValidationResult, relativePath string) error {
	hasFrom := false
	hasWorkdir := false
	hasCopy := false

	err := mov.ValidateFileStreaming(path, func(line string, lineNum int) error {
		line = strings.TrimSpace(line)
		upperLine := strings.ToUpper(line)

		if strings.HasPrefix(upperLine, "FROM ") {
			hasFrom = true
		}
		if strings.HasPrefix(upperLine, "WORKDIR ") {
			hasWorkdir = true
		}
		if strings.HasPrefix(upperLine, "COPY ") || strings.HasPrefix(upperLine, "ADD ") {
			hasCopy = true
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Check for required instructions
	if !hasFrom {
		result.Valid = false
		result.Errors = append(result.Errors, models.ValidationError{
			Field:   relativePath,
			Tag:     "required",
			Value:   "FROM",
			Message: "Dockerfile missing FROM instruction",
		})
	}

	if !hasWorkdir {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   relativePath,
			Message: "Dockerfile missing WORKDIR instruction (recommended)",
		})
	}

	if !hasCopy {
		result.Warnings = append(result.Warnings, models.ValidationWarning{
			Field:   relativePath,
			Message: "Dockerfile missing COPY or ADD instruction",
		})
	}

	return nil
}

// mergeValidationResults merges validation results efficiently
func (mov *MemoryOptimizedValidator) mergeValidationResults(target, source *models.ValidationResult, prefix string) {
	if !source.Valid {
		target.Valid = false
	}

	// Add errors with path prefix
	for _, err := range source.Errors {
		err.Field = fmt.Sprintf("%s.%s", prefix, err.Field)
		target.Errors = append(target.Errors, err)
	}

	// Add warnings with path prefix
	for _, warning := range source.Warnings {
		warning.Field = fmt.Sprintf("%s.%s", prefix, warning.Field)
		target.Warnings = append(target.Warnings, warning)
	}
}

// Cleanup releases resources used by the validator
func (mov *MemoryOptimizedValidator) Cleanup() {
	if mov.stringBuilder != nil {
		mov.stringBuilder.Reset()
	}
	// Force garbage collection to clean up any remaining allocations
	utils.ForceGlobalGC()
}
