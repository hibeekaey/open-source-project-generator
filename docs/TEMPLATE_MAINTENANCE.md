# Template Maintenance Guidelines

## Overview

This document provides comprehensive guidelines for maintaining Go template files in the project generator. Following these guidelines ensures that generated code compiles successfully and follows Go best practices.

## Import Management

### Standard Library Import Organization

All template files must organize imports according to Go conventions:

1. **Standard library imports** - grouped first, alphabetically ordered
2. **Third-party imports** - grouped second, alphabetically ordered
3. **Local imports** - grouped last, alphabetically ordered

### Example of Correct Import Organization

```go
package {{.PackageName}}

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    
    "{{.ModuleName}}/internal/models"
    "{{.ModuleName}}/pkg/utils"
)
```

### Import Detection and Validation

Before adding any function calls to templates, ensure the corresponding package is imported:

#### Common Standard Library Functions and Their Required Imports

| Function/Type | Required Import | Example Usage |
|---------------|----------------|---------------|
| `time.Now()`, `time.Duration` | `"time"` | `time.Now().Unix()` |
| `fmt.Sprintf()`, `fmt.Errorf()` | `"fmt"` | `fmt.Sprintf("user_%d", id)` |
| `strings.Contains()`, `strings.Split()` | `"strings"` | `strings.Contains(email, "@")` |
| `strconv.Atoi()`, `strconv.Itoa()` | `"strconv"` | `strconv.Atoi(userID)` |
| `context.Context`, `context.Background()` | `"context"` | `ctx context.Context` |
| `log.Printf()`, `log.Fatal()` | `"log"` | `log.Printf("Error: %v", err)` |
| `http.StatusOK`, `http.Request` | `"net/http"` | `c.JSON(http.StatusOK, data)` |
| `json.Marshal()`, `json.Unmarshal()` | `"encoding/json"` | `json.Marshal(user)` |
| `errors.New()`, `errors.Is()` | `"errors"` | `errors.New("invalid input")` |
| `os.Getenv()`, `os.Exit()` | `"os"` | `os.Getenv("PORT")` |

## Guidelines for Adding New Functions

### 1. Function Addition Checklist

Before adding any new function to a template:

- [ ] Identify all packages required by the function
- [ ] Add necessary import statements
- [ ] Verify import organization follows conventions
- [ ] Test template compilation with sample data
- [ ] Update this documentation if introducing new patterns

### 2. Template Variable Handling

When adding functions that interact with template variables:

```go
// ✅ Correct - Handle template variables properly
func (s *{{.Name}}Service) CreateUser(ctx context.Context, user *models.User) error {
    if user.Email == "" {
        return fmt.Errorf("email is required")
    }
    // ... rest of function
}

// ❌ Incorrect - Missing fmt import for fmt.Errorf
func (s *{{.Name}}Service) CreateUser(ctx context.Context, user *models.User) error {
    if user.Email == "" {
        return fmt.Errorf("email is required") // This will cause compilation error
    }
}
```

### 3. Error Handling Patterns

Always include proper error handling imports:

```go
import (
    "errors"
    "fmt"
)

// Standard error creation patterns
var ErrUserNotFound = errors.New("user not found")

func validateUser(user *User) error {
    if user.Email == "" {
        return fmt.Errorf("email is required")
    }
    return nil
}
```

## Validation Process for Template Changes

### 1. Pre-Commit Validation

Before committing template changes:

1. **Run Import Detection**: Use the import analyzer tool to scan for missing imports

   ```bash
   go run scripts/validate-templates/main.go --scan-imports
   ```

2. **Generate Test Project**: Create a sample project using modified templates

   ```bash
   go run cmd/generator/main.go --config config/test-configs/test-config.yaml --output test-output
   ```

3. **Compile Generated Code**: Verify all generated Go files compile

   ```bash
   cd test-output && go mod tidy && go build ./...
   ```

### 2. Template Compilation Testing

Use the validation script to test template compilation:

```bash
# Run comprehensive template validation
go run scripts/validate-templates/main.go --validate-all

# Test specific template
go run scripts/validate-templates/main.go --template templates/backend/go-gin/internal/middleware/auth.go.tmpl
```

### 3. Integration Testing

After template modifications:

1. Generate multiple project types to test template variations
2. Run generated project tests to ensure functionality
3. Verify Docker builds work with generated code
4. Test with different configuration parameters

## Common Import Issues and Solutions

### Issue 1: Missing Time Import

**Problem**: Using `time.Now()` without importing `"time"`

```go
// ❌ Missing import
func (m *AuthMiddleware) ValidateToken(token string) error {
    if time.Now().After(expiry) { // Compilation error
        return errors.New("token expired")
    }
}
```

**Solution**: Add time import

```go
import (
    "errors"
    "time"
)

func (m *AuthMiddleware) ValidateToken(token string) error {
    if time.Now().After(expiry) { // ✅ Works correctly
        return errors.New("token expired")
    }
}
```

### Issue 2: Incorrect Import Grouping

**Problem**: Mixed import organization

```go
// ❌ Incorrect grouping
import (
    "github.com/gin-gonic/gin"
    "fmt"
    "{{.ModuleName}}/internal/models"
    "time"
)
```

**Solution**: Proper grouping with blank lines

```go
// ✅ Correct grouping
import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"

    "{{.ModuleName}}/internal/models"
)
```

### Issue 3: Unused Imports

**Problem**: Imports that aren't used in the template

```go
// ❌ Unused imports
import (
    "fmt"
    "log"    // Not used anywhere
    "time"
)
```

**Solution**: Remove unused imports

```go
// ✅ Only necessary imports
import (
    "fmt"
    "time"
)
```

## Template-Specific Considerations

### 1. Conditional Imports

For templates with conditional logic, ensure all code paths have required imports:

```go
{{if .EnableAuth}}
import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)
{{else}}
import (
    "fmt"
)
{{end}}
```

### 2. Template Variable Imports

When template variables affect import requirements:

```go
import (
    "{{.ModuleName}}/internal/config"
    "{{.ModuleName}}/internal/models"
    {{if .EnableDatabase}}
    "{{.ModuleName}}/internal/database"
    {{end}}
)
```

### 3. Cross-Template Dependencies

Ensure consistent imports across related templates:

- Service templates should import their corresponding model packages
- Controller templates should import service and model packages
- Test templates should import testing packages and the code under test

## Maintenance Workflow

### 1. Regular Validation

Run template validation weekly or before releases:

```bash
# Full validation suite
make validate-templates

# Quick import check
go run scripts/validate-templates/main.go --check-imports
```

### 2. Documentation Updates

When adding new template patterns:

1. Update this documentation with new function mappings
2. Add examples of correct usage
3. Document any new validation requirements
4. Update the import detection tool if needed

### 3. Version Control

- Always commit template changes with corresponding documentation updates
- Include validation results in commit messages
- Tag releases that include template changes

## Tools and Scripts

### Available Validation Tools

1. **Import Analyzer**: `scripts/validate-templates/main.go`
   - Scans templates for missing imports
   - Generates comprehensive reports
   - Validates import organization

2. **Template Compiler**: `scripts/validate-templates/compiler.go`
   - Generates sample Go files from templates
   - Tests compilation with various configurations
   - Reports compilation errors

3. **Integration Tester**: `pkg/template/template_compilation_integration_test.go`
   - End-to-end template testing
   - Multi-project generation testing
   - Docker build validation

### Usage Examples

```bash
# Scan all templates for import issues
go run scripts/validate-templates/main.go --scan

# Fix missing imports automatically
go run scripts/validate-templates/main.go --fix

# Validate specific template file
go run scripts/validate-templates/main.go --file templates/backend/go-gin/internal/services/auth_service.go.tmpl

# Generate test project and validate
go run cmd/generator/main.go --config config/test-configs/minimal-test-config.yaml --output validation-test
cd validation-test && go mod tidy && go build ./...
```

## Best Practices Summary

1. **Always validate imports** before committing template changes
2. **Follow Go import conventions** with proper grouping and ordering
3. **Test template compilation** with sample data regularly
4. **Document new patterns** and update guidelines accordingly
5. **Use provided tools** for automated validation and fixing
6. **Maintain consistency** across related templates
7. **Include error handling** imports when adding error-prone code
8. **Remove unused imports** to keep templates clean
9. **Test with multiple configurations** to ensure template flexibility
10. **Update documentation** when introducing new template patterns

Following these guidelines ensures that all generated Go code compiles successfully and maintains high code quality standards.
