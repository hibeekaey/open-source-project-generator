# Template Validation Script

This tool validates Go template files by generating sample Go code and checking for compilation errors.

## Features

- **Template Generation**: Generates Go files from templates using test data
- **Syntax Validation**: Uses Go's AST parser to validate generated code syntax
- **Import Validation**: Checks for proper import statements and organization
- **Comprehensive Reporting**: Provides detailed reports of validation results
- **Error Detection**: Identifies missing imports and compilation issues

## Usage

### Basic Usage

```bash
# Validate all templates in the templates directory
./validator ../../templates

# Validate a specific template directory
./validator ../../templates/backend/go-gin

# Specify custom output directory
./validator ../../templates validation-output
```

### Using the Shell Script

```bash
# Run the validation script (builds and runs automatically)
./run-validation.sh
```

## How It Works

1. **Template Discovery**: Scans the specified directory for `.tmpl` files
2. **Go File Filtering**: Identifies Go template files (`.go.tmpl` and `.mod.tmpl`)
3. **Template Generation**: Uses test data to generate actual Go files from templates
4. **Syntax Validation**: Parses generated Go files using Go's AST parser
5. **Import Analysis**: Validates import statements and checks for missing imports
6. **Report Generation**: Provides comprehensive validation results

## Test Data

The validator uses predefined test data to populate template variables:

- **Project Name**: `testproject`
- **Organization**: `testorg`
- **Description**: `A test project for template validation`
- **License**: `MIT`
- **Go Version**: `1.22`
- **Node Version**: `20.0.0`

## Validation Features

### Go Syntax Validation

- Parses Go source files using `go/parser`
- Detects syntax errors and malformed code
- Validates package declarations

### Import Validation

- Checks for duplicate imports
- Validates import path formats
- Detects missing standard library imports
- Ensures proper import organization

### Template-Specific Validation

- Handles Go module files (`go.mod`)
- Processes template variables (`{{.Name}}`, `{{.Organization}}`, etc.)
- Maintains template syntax during validation

## Output

The validator generates:

1. **Console Output**: Real-time validation progress and results
2. **Generated Files**: Sample Go files in the output directory (for debugging)
3. **Validation Report**: Summary of total, valid, and invalid templates
4. **Error Details**: Specific compilation errors with file locations

## Example Output

```
🔍 Starting template validation...
📝 Validating: ../../templates/backend/go-gin/main.go.tmpl
✅ Valid: ../../templates/backend/go-gin/main.go.tmpl
📝 Validating: ../../templates/backend/go-gin/go.mod.tmpl
✅ Valid: ../../templates/backend/go-gin/go.mod.tmpl

============================================================
📊 TEMPLATE VALIDATION REPORT
============================================================
📁 Total Templates: 75
✅ Valid Templates: 75
❌ Invalid Templates: 0
============================================================
🎉 All templates are valid!
```

## Building

```bash
# Build the validator
go build -o validator .

# Run tests (if any)
go test ./...
```

## Requirements

- Go 1.22 or later
- Access to the main project's `pkg/models` package

## Integration

This validation script can be integrated into:

- **CI/CD Pipelines**: Validate templates on every commit
- **Pre-commit Hooks**: Ensure template quality before commits
- **Development Workflow**: Regular validation during template development
- **Release Process**: Final validation before releases

## Troubleshooting

### Common Issues

1. **Module Resolution**: Ensure the `go.mod` file correctly references the main project
2. **Template Variables**: Verify all template variables are included in test data
3. **Import Paths**: Check that generated import paths are valid
4. **File Permissions**: Ensure the validator has read/write access to directories

### Debug Mode

To debug template generation issues:

1. Check the generated files in the output directory
2. Manually inspect the generated Go code
3. Run `go build` on individual generated files
4. Review the validation report for specific error details
