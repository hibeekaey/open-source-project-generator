# Template Development Guide

This guide covers creating, maintaining, and extending templates for the Open Source Project Generator.

## Table of Contents

- [Template Structure](#template-structure)
- [Template Syntax](#template-syntax)
- [Template Variables](#template-variables)
- [Template Functions](#template-functions)
- [Template Maintenance](#template-maintenance)
- [Custom Templates](#custom-templates)
- [Template Validation](#template-validation)
- [Best Practices](#best-practices)

## Template Structure

### Directory Organization

Templates are organized in a hierarchical structure:

```text
pkg/template/templates/
├── base/                  # Base project files
│   ├── README.md.tmpl
│   ├── LICENSE.tmpl
│   ├── .gitignore.tmpl
│   └── Makefile.tmpl
├── frontend/              # Frontend templates
│   ├── nextjs-app/
│   │   ├── package.json.tmpl
│   │   ├── next.config.js.tmpl
│   │   └── src/
│   ├── react-app/
│   └── vue-app/
├── backend/               # Backend templates
│   ├── go-gin/
│   │   ├── go.mod.tmpl
│   │   ├── main.go.tmpl
│   │   └── internal/
│   ├── node-express/
│   └── python-fastapi/
├── mobile/                # Mobile templates
│   ├── android-kotlin/
│   ├── ios-swift/
│   └── react-native/
└── infrastructure/        # Infrastructure templates
    ├── docker/
    ├── kubernetes/
    └── terraform/
```

### Template File Naming

Template files use the `.tmpl` extension:

- `package.json.tmpl` → `package.json`
- `main.go.tmpl` → `main.go`
- `Dockerfile.tmpl` → `Dockerfile`

### Template Metadata

Each template directory should include a `metadata.yaml` file:

```yaml
# metadata.yaml
name: "go-gin"
title: "Go Gin API Server"
description: "Production-ready Go API server with Gin framework"
category: "backend"
technology: "go"
version: "1.0.0"
author: "Generator Team"
license: "MIT"

# Template requirements
requirements:
  go: ">=1.21"
  node: ">=20"

# Template variables
variables:
  - name: "project_name"
    type: "string"
    required: true
    description: "Name of the project"
  - name: "organization"
    type: "string"
    required: true
    description: "Organization name"
  - name: "database_type"
    type: "string"
    required: false
    default: "postgresql"
    description: "Database type to use"

# Dependencies
dependencies:
  - name: "gin"
    version: "latest"
    type: "go_module"
  - name: "gorm"
    version: "latest"
    type: "go_module"

# Template features
features:
  - "RESTful API"
  - "JWT Authentication"
  - "Database Integration"
  - "Middleware Support"
  - "Testing Framework"
  - "Docker Support"

# Tags for categorization
tags:
  - "api"
  - "backend"
  - "go"
  - "gin"
  - "rest"
  - "jwt"
  - "database"
```

## Template Syntax

### Go Template Syntax

Templates use Go's `text/template` package syntax:

```go
// Basic variable substitution
{{.ProjectName}}

// Conditional rendering
{{if .EnableAuth}}
// Authentication code
{{end}}

// Loops
{{range .Dependencies}}
import "{{.}}"
{{end}}

// Template functions
{{.ProjectName | toUpper}}
{{.Version | formatVersion}}
```

### Template Variables

Access project configuration variables:

```go
// Project metadata
{{.Name}}                    // Project name
{{.Organization}}            // Organization name
{{.Description}}             // Project description
{{.License}}                 // License type
{{.Author}}                  // Author name
{{.Email}}                   // Author email
{{.Repository}}              // Repository URL

// Component selection
{{.Components.Frontend.MainApp}}     // Frontend main app
{{.Components.Backend.API}}          // Backend API
{{.Components.Mobile.Android}}       // Android mobile
{{.Components.Infrastructure.Docker}} // Docker infrastructure

// Generation options
{{.GenerateOptions.Force}}          // Force overwrite
{{.GenerateOptions.Minimal}}        // Minimal generation
{{.GenerateOptions.Offline}}        // Offline mode
{{.GenerateOptions.UpdateVersions}} // Update versions

// Custom variables
{{.CustomVars.DatabaseName}}        // Custom database name
{{.CustomVars.APIPort}}             // Custom API port
{{.CustomVars.FrontendPort}}        // Custom frontend port
```

### Conditional Rendering

Use conditional statements to include/exclude code based on component selection:

```go
{{if .Components.Frontend.MainApp}}
// Frontend-specific configuration
{
  "name": "{{.Name}}-frontend",
  "version": "1.0.0",
  "dependencies": {
    "next": "{{.Versions.NextJS}}",
    "react": "{{.Versions.React}}"
  }
}
{{end}}

{{if .Components.Backend.API}}
// Backend-specific configuration
package main

import (
    "github.com/gin-gonic/gin"
    "{{.ModuleName}}/internal/handlers"
)
{{end}}
```

### Loops and Iteration

Iterate over collections:

```go
{{range .Dependencies}}
import "{{.}}"
{{end}}

{{range $key, $value := .EnvironmentVariables}}
export {{$key}}="{{$value}}"
{{end}}

{{range .Components.Frontend}}
// {{.Name}} component
{{end}}
```

## Template Variables

### Built-in Variables

The generator provides these built-in variables to all templates:

```go
// Project information
.Name                      // Project name
.Organization             // Organization name
.Description              // Project description
.License                  // License type
.Author                   // Author name
.Email                    // Author email
.Repository               // Repository URL
.OutputPath               // Output directory path

// Component selection
.Components.Frontend.MainApp
.Components.Frontend.Home
.Components.Frontend.Admin
.Components.Frontend.SharedComponents
.Components.Backend.API
.Components.Backend.Auth
.Components.Backend.Database
.Components.Mobile.Android
.Components.Mobile.iOS
.Components.Mobile.Shared
.Components.Infrastructure.Docker
.Components.Infrastructure.Kubernetes
.Components.Infrastructure.Terraform
.Components.Infrastructure.Monitoring

// Generation options
.GenerateOptions.Force
.GenerateOptions.Minimal
.GenerateOptions.Offline
.GenerateOptions.UpdateVersions
.GenerateOptions.SkipValidation
.GenerateOptions.BackupExisting
.GenerateOptions.IncludeExamples

// Version information
.Versions.Node
.Versions.Go
.Versions.React
.Versions.TypeScript
.Versions.NextJS
.Versions.TailwindCSS
.Versions.Kotlin
.Versions.Swift
.Versions.Docker
.Versions.Kubernetes
.Versions.Terraform

// Custom variables
.CustomVars               // Map of custom variables
```

### Custom Variables

Define custom variables in your configuration:

```yaml
# project-config.yaml
custom_vars:
  database_name: "myapp_db"
  api_port: "8080"
  frontend_port: "3000"
  cache_type: "redis"
  monitoring: "prometheus"
```

Access in templates:

```go
{{.CustomVars.DatabaseName}}
{{.CustomVars.APIPort}}
{{.CustomVars.FrontendPort}}
{{.CustomVars.CacheType}}
{{.CustomVars.Monitoring}}
```

## Template Functions

### Built-in Functions

The generator provides these built-in template functions:

```go
// String manipulation
{{.Name | toUpper}}           // Convert to uppercase
{{.Name | toLower}}           // Convert to lowercase
{{.Name | toTitle}}           // Convert to title case
{{.Name | toCamel}}           // Convert to camelCase
{{.Name | toSnake}}           // Convert to snake_case
{{.Name | toKebab}}           // Convert to kebab-case
{{.Name | toPascal}}          // Convert to PascalCase

// String operations
{{.Name | replace " " "-"}}   // Replace characters
{{.Name | trim}}              // Trim whitespace
{{.Name | truncate 20}}       // Truncate to length
{{.Name | pad 10}}            // Pad to length

// Version formatting
{{.Version | formatVersion}}  // Format version (^1.0.0)
{{.Version | semver}}         // Semantic version validation
{{.Version | major}}          // Get major version
{{.Version | minor}}          // Get minor version
{{.Version | patch}}          // Get patch version

// Date and time
{{.Date | formatDate "2006-01-02"}}  // Format date
{{.Time | formatTime "15:04:05"}}     // Format time
{{.Now | formatDateTime}}             // Current date and time

// Security functions
{{.Secret | generateSecret}}          // Generate secure random string
{{.Password | hashPassword}}          // Hash password
{{.Token | generateToken}}            // Generate JWT token

// File operations
{{.Path | cleanPath}}                 // Clean file path
{{.Path | baseName}}                  // Get base name
{{.Path | dirName}}                   // Get directory name
{{.Path | extName}}                   // Get file extension

// Mathematical operations
{{.Number | add 10}}                  // Add numbers
{{.Number | multiply 2}}              // Multiply numbers
{{.Number | divide 3}}                // Divide numbers
{{.Number | round}}                   // Round numbers

// Conditional functions
{{.Value | ifEmpty "default"}}        // Default value if empty
{{.Value | ifTrue "yes" "no"}}        // Conditional value
{{.Value | ifExists "exists"}}        // Check if value exists
```

### Custom Functions

Register custom functions in your templates:

```go
// Register custom functions
engine.RegisterFunctions(template.FuncMap{
    "formatVersion": func(version string) string {
        return "^" + version
    },
    "generateSecret": func() string {
        return generateRandomString(32)
    },
    "formatDate": func(format string) string {
        return time.Now().Format(format)
    },
    "companyName": func() string {
        return "MyCompany"
    },
})
```

Use in templates:

```go
{{.Version | formatVersion}}
{{generateSecret}}
{{formatDate "2006-01-02"}}
{{companyName}}
```

## Template Maintenance

### Import Management

All template files must organize imports according to Go conventions:

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

### Common Import Requirements

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

### Template Validation

Before adding any function to a template:

1. **Identify all packages** required by the function
2. **Add necessary import statements**
3. **Verify import organization** follows conventions
4. **Test template compilation** with sample data
5. **Update documentation** if introducing new patterns

### Error Handling Patterns

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

## Custom Templates

### Creating Custom Templates

1. **Create template directory**:

```bash
mkdir -p ./custom-templates/my-template
cd ./custom-templates/my-template
```

1. **Create metadata file**:

```yaml
# metadata.yaml
name: "my-template"
title: "My Custom Template"
description: "A custom template for my specific needs"
category: "custom"
technology: "mixed"
version: "1.0.0"
author: "Your Name"
license: "MIT"

variables:
  - name: "custom_setting"
    type: "string"
    required: true
    description: "Custom setting for the template"

features:
  - "Custom Feature 1"
  - "Custom Feature 2"

tags:
  - "custom"
  - "template"
```

2. **Create template files**:

```go
// main.go.tmpl
package main

import (
    "fmt"
    "log"
)

func main() {
    fmt.Println("Hello from {{.Name}}!")
    log.Printf("Custom setting: %s", "{{.CustomVars.CustomSetting}}")
}
```

3. **Validate template**:

```bash
generator template validate ./custom-templates/my-template
```

### Custom Template Structure

Organize your custom template:

```text
my-template/
├── metadata.yaml          # Template metadata
├── templates/             # Template files
│   ├── main.go.tmpl
│   ├── package.json.tmpl
│   └── README.md.tmpl
├── examples/              # Example configurations
│   ├── basic.yaml
│   └── advanced.yaml
└── docs/                  # Template documentation
    ├── README.md
    └── USAGE.md
```

### Custom Template Variables

Define variables your template needs:

```yaml
# metadata.yaml
variables:
  - name: "api_port"
    type: "string"
    required: true
    default: "8080"
    description: "Port for the API server"
  - name: "database_url"
    type: "string"
    required: false
    default: "postgresql://localhost:5432/mydb"
    description: "Database connection URL"
  - name: "enable_auth"
    type: "boolean"
    required: false
    default: true
    description: "Enable authentication"
```

Use in templates:

```go
{{if .CustomVars.EnableAuth}}
// Authentication middleware
func authMiddleware() gin.HandlerFunc {
    return gin.BasicAuth(gin.Accounts{
        "admin": "password",
    })
}
{{end}}
```

## Template Validation

### Validation Commands

```bash
# Validate template structure
generator template validate ./my-template

# Validate with detailed output
generator template validate ./my-template --detailed

# Validate and auto-fix issues
generator template validate ./my-template --fix

# Validate specific template file
generator template validate ./my-template/templates/main.go.tmpl
```

### Validation Rules

The validator checks for:

1. **Template structure** - Required files and directories
2. **Metadata validation** - Required fields and format
3. **Template syntax** - Go template syntax errors
4. **Import organization** - Proper import grouping
5. **Variable usage** - All variables are defined
6. **Function usage** - All functions are available
7. **File permissions** - Proper file permissions
8. **Naming conventions** - Consistent naming

### Validation Output

```bash
$ generator template validate ./my-template --detailed

✅ Template validation completed successfully

Validation Results:
  ✅ Structure: Valid
  ✅ Metadata: Valid
  ✅ Syntax: Valid
  ✅ Imports: Valid
  ✅ Variables: Valid
  ✅ Functions: Valid
  ✅ Permissions: Valid
  ✅ Naming: Valid

Issues Found: 0
Warnings: 0
Errors: 0
```

## Best Practices

### Template Design

1. **Keep templates simple** - Avoid complex logic in templates
2. **Use clear variable names** - Make variables self-documenting
3. **Provide sensible defaults** - Set reasonable default values
4. **Document variables** - Include descriptions for all variables
5. **Test thoroughly** - Test with various configurations

### Code Quality

1. **Follow Go conventions** - Use standard Go formatting
2. **Include proper imports** - Add all required imports
3. **Handle errors properly** - Include error handling
4. **Use meaningful names** - Choose descriptive variable names
5. **Add comments** - Document complex logic

### Maintenance

1. **Version your templates** - Use semantic versioning
2. **Update dependencies** - Keep dependencies current
3. **Test regularly** - Validate templates frequently
4. **Document changes** - Keep changelog updated
5. **Backup templates** - Keep backups of working templates

### Performance

1. **Minimize complexity** - Keep templates simple
2. **Use efficient functions** - Choose optimal template functions
3. **Cache when possible** - Enable template caching
4. **Avoid deep nesting** - Keep template structure flat
5. **Profile performance** - Monitor template performance

### Security

1. **Validate inputs** - Sanitize all user inputs
2. **Use secure defaults** - Set secure default values
3. **Avoid dangerous functions** - Don't use unsafe functions
4. **Limit permissions** - Use minimal required permissions
5. **Audit regularly** - Regular security audits

### Team Collaboration

1. **Use version control** - Track template changes
2. **Code review** - Review template changes
3. **Documentation** - Keep documentation updated
4. **Testing** - Comprehensive testing procedures
5. **Standards** - Establish team standards

## Quick Reference

### Import Organization Template

```go
package {{.PackageName}}

import (
    // Standard library (alphabetical)
    "context"
    "fmt"
    "time"

    // Third-party (alphabetical)
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    
    // Local (alphabetical)
    "{{.ModuleName}}/internal/models"
    "{{.ModuleName}}/pkg/utils"
)
```

### Common Function → Import Mapping

| Function | Import | Function | Import |
|----------|--------|----------|--------|
| `time.Now()` | `"time"` | `fmt.Sprintf()` | `"fmt"` |
| `strings.Contains()` | `"strings"` | `strconv.Atoi()` | `"strconv"` |
| `context.Background()` | `"context"` | `log.Printf()` | `"log"` |
| `http.StatusOK` | `"net/http"` | `json.Marshal()` | `"encoding/json"` |
| `errors.New()` | `"errors"` | `os.Getenv()` | `"os"` |

### Quick Validation Commands

```bash
# Quick import check
go run scripts/validate-templates/main.go --check-imports

# Full validation
make validate-templates

# Test specific template
go run scripts/validate-templates/main.go --file path/to/template.tmpl

# Generate and test project
go run cmd/generator/main.go --config config/test-configs/test-config.yaml --output test-output
cd test-output && go mod tidy && go build ./...
```

### Pre-Commit Checklist

- [ ] All used functions have corresponding imports
- [ ] Imports are properly grouped and ordered
- [ ] No unused imports remain
- [ ] Template compiles with test data
- [ ] Integration tests pass

### Common Fixes

#### Missing Time Import

```go
// Add to imports
"time"

// Usage
if time.Now().After(expiry) { ... }
```

#### Missing Context Import

```go
// Add to imports
"context"

// Usage
func Handler(ctx context.Context) { ... }
```

#### Missing HTTP Import

```go
// Add to imports
"net/http"

// Usage
c.JSON(http.StatusOK, response)
```

This comprehensive template development guide provides all the information needed to create, maintain, and extend templates for the Open Source Project Generator.
