# Design Document

## Overview

The solution involves systematically addressing each category of golangci-lint violations through targeted fixes that preserve functionality while improving code quality. The approach prioritizes critical issues first, then maintainability improvements, and finally cosmetic fixes.

## Architecture

### Current State Analysis

The codebase has accumulated technical debt across several dimensions:
- **Duplicate code**: Functions with identical logic in validation package
- **Error handling**: Many unchecked error returns, especially in test files
- **Code formatting**: Some files don't follow gofmt standards
- **Complexity**: Several functions exceed recommended complexity thresholds
- **Variable shadowing**: Common variable names reused in inner scopes
- **Constants**: Repeated string literals throughout the codebase

### Target State

A clean codebase that:
- Passes all golangci-lint checks without errors or warnings
- Follows Go best practices for error handling and code organization
- Uses appropriate abstraction levels and constant extraction
- Maintains all existing functionality and test coverage

### Migration Strategy

1. **Phase 1**: Fix critical issues (duplicates, error checking, formatting)
2. **Phase 2**: Reduce complexity and improve maintainability
3. **Phase 3**: Extract constants and fix minor issues
4. **Phase 4**: Address deprecation warnings and final cleanup

## Components and Interfaces

### 1. Error Handling Strategy

**Current Issues**: Many functions ignore error returns, especially in test code and async operations.

**Solution Approach**:
```go
// Instead of ignoring errors
cache.Set("key", "value")

// Check and handle appropriately
if err := cache.Set("key", "value"); err != nil {
    // In tests: use t.Errorf or require
    // In production: log and/or propagate
    t.Errorf("Failed to set cache value: %v", err)
}
```

**Implementation Strategy**:
- Test files: Add proper error checking with t.Errorf or require
- Production code: Add error logging or propagation as appropriate
- Async operations: Use error channels or logging for background goroutines

### 2. Duplicate Code Elimination

**Current Issue**: `pkg/validation/setup.go` has duplicate function logic.

**Solution Approach**:
```go
// Extract common verification logic
func (s *SetupEngine) verifyComponents(projectPath string, config *models.ProjectConfig, 
    verifyFunc func(string, *models.ProjectConfig, *models.ValidationResult) error) (*models.ValidationResult, error) {
    
    result := &models.ValidationResult{
        Valid:    true,
        Errors:   []models.ValidationError{},
        Warnings: []models.ValidationWarning{},
    }
    
    if err := verifyFunc(projectPath, config, result); err != nil {
        return nil, err
    }
    
    return result, nil
}
```

### 3. Complexity Reduction Strategy

**High Complexity Functions**:
- `(*App).displayDashboardTable` (complexity 29)
- `(*CLI).displayAnalysisReport` (complexity 21)
- And 12 others with complexity 16-19

**Reduction Techniques**:
1. **Extract Methods**: Break large functions into smaller, focused methods
2. **Early Returns**: Use guard clauses to reduce nesting
3. **Strategy Pattern**: Replace complex conditional logic with interfaces
4. **Helper Functions**: Extract repeated logic patterns

**Example Refactoring**:
```go
// Before: Complex function with many conditions
func (a *App) displayDashboardTable(data *models.DashboardData, showDetails bool) error {
    // 50+ lines of complex logic
}

// After: Decomposed into focused methods
func (a *App) displayDashboardTable(data *models.DashboardData, showDetails bool) error {
    if err := a.displayHeader(data); err != nil {
        return err
    }
    if err := a.displaySummary(data); err != nil {
        return err
    }
    if showDetails {
        return a.displayDetails(data)
    }
    return nil
}
```

### 4. Constants Extraction

**String Constants to Extract**:
- File type identifiers: "javascript", "typescript", "nodejs"
- Status values: "failed", "consistent", "critical", "high", "medium", "low"
- Package manager names: "npm", "yarn"
- Format types: "yaml", "json"

**Organization Strategy**:
```go
// pkg/constants/types.go
package constants

const (
    // Package managers
    PackageManagerNPM  = "npm"
    PackageManagerYarn = "yarn"
    
    // Status levels
    StatusFailed     = "failed"
    StatusConsistent = "consistent"
    
    // Severity levels
    SeverityCritical = "critical"
    SeverityHigh     = "high"
    SeverityMedium   = "medium"
    SeverityLow      = "low"
)
```

### 5. Modern API Usage

**Deprecated APIs to Replace**:
- `strings.Title` → `golang.org/x/text/cases.Title`
- Built-in types as context keys → custom types

**Replacement Strategy**:
```go
// Replace strings.Title
import "golang.org/x/text/cases"

// Instead of: strings.Title(str)
caser := cases.Title(language.English)
result := caser.String(str)

// Replace context keys
type contextKey string
const resourceManagerKey contextKey = "resourceManager"

// Instead of: context.WithValue(ctx, "resourceManager", rm)
ctx = context.WithValue(ctx, resourceManagerKey, rm)
```

## Data Models

### Error Context Enhancement

```go
// Enhanced error context handling
type ErrorContext struct {
    Component string
    Operation string
    File      string
    Line      int
}

func (e *GeneratorError) WithContext(key, value string) error {
    // Return error instead of ignoring
    if e.context == nil {
        e.context = make(map[string]string)
    }
    e.context[key] = value
    return e
}
```

### Constants Package Structure

```go
package constants

// File types
const (
    FileTypeJavaScript = "javascript"
    FileTypeTypeScript = "typescript"
    FileTypeNodeJS     = "nodejs"
)

// Package managers
const (
    PackageManagerNPM  = "npm"
    PackageManagerYarn = "yarn"
)

// Status indicators
const (
    StatusSuccess    = "✅"
    StatusFailure    = "❌"
    StatusConsistent = "consistent"
    StatusFailed     = "failed"
)
```

## Error Handling

### Test Error Handling Pattern

```go
// In test functions
func TestSomething(t *testing.T) {
    // Use require for critical setup
    require.NoError(t, setupOperation())
    
    // Use assert for non-critical checks
    assert.NoError(t, cache.Set("key", "value"))
    
    // Use t.Cleanup for resource cleanup
    t.Cleanup(func() {
        if err := os.RemoveAll(tempDir); err != nil {
            t.Logf("Cleanup warning: %v", err)
        }
    })
}
```

### Production Error Handling Pattern

```go
// In production code
func (s *Service) operation() error {
    if err := dependency.Call(); err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    return nil
}

// For background operations
func (c *Cache) backgroundSave() {
    if err := c.save(); err != nil {
        c.logger.Error("Background save failed", "error", err)
    }
}
```

## Testing Strategy

### Validation Approach

1. **Unit Tests**: Ensure all existing tests continue to pass
2. **Integration Tests**: Verify end-to-end functionality
3. **Linting Tests**: Add golangci-lint to CI pipeline
4. **Regression Tests**: Test specific areas where complexity was reduced

### Test Coverage Maintenance

- Maintain or improve existing test coverage
- Add tests for newly extracted functions
- Ensure error handling paths are covered

## Implementation Phases

### Phase 1: Critical Fixes (High Priority)
- Fix duplicate code
- Add error checking
- Fix formatting issues
- Address critical staticcheck issues

### Phase 2: Maintainability (Medium Priority)
- Reduce function complexity
- Fix variable shadowing
- Remove ineffectual assignments

### Phase 3: Constants and Standards (Low Priority)
- Extract string constants
- Fix misspellings
- Replace deprecated APIs

### Phase 4: Validation and Cleanup
- Run comprehensive tests
- Verify golangci-lint passes
- Document changes

## Security Considerations

- Ensure error handling doesn't expose sensitive information
- Validate that complexity reduction doesn't introduce security vulnerabilities
- Maintain proper input validation in refactored functions

## Performance Impact

### Expected Improvements
- Reduced memory allocations from better pointer usage
- Improved maintainability leading to easier optimization
- Cleaner code enabling better compiler optimizations

### Potential Concerns
- Constants extraction may have minimal memory overhead
- Function decomposition may have slight call overhead (negligible)

### Mitigation Strategies
- Use benchmarks to validate performance is maintained
- Profile critical paths before and after changes
- Keep hot paths optimized during refactoring
