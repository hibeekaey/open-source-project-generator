# Import Detection Utility

A Go tool for analyzing template files to detect missing import statements that would cause compilation errors in generated code.

## Features

- **AST-based Analysis**: Uses Go's AST parser to accurately detect function calls and required imports
- **Template Preprocessing**: Handles Go template syntax by preprocessing files before analysis
- **Comprehensive Mapping**: Includes mappings for common standard library functions
- **Multiple Output Formats**: Supports JSON, text, and summary output formats
- **Structured Reports**: Provides detailed reports with file locations and line numbers

## Usage

### Basic Usage

```bash
# Analyze templates directory with text output
./import-detector -dir templates

# Generate JSON report
./import-detector -dir templates -format json -output report.json

# Show summary only
./import-detector -dir templates -format summary

# Verbose output
./import-detector -dir templates -verbose
```

### Command Line Options

- `-dir`: Directory containing template files (default: "templates")
- `-output`: Output file path (default: stdout)
- `-format`: Output format - json, text, or summary (default: "text")
- `-verbose`: Enable verbose output for debugging

### Output Formats

#### Text Format

Provides detailed information about each file with missing imports:

```
Import Detection Report
======================

File: templates/backend/go-gin/internal/middleware/auth.go.tmpl
  Missing imports:
    - time
  Function usages requiring imports:
    - time.Now (line 25) -> requires time

Summary: 1/45 files have missing imports
```

#### JSON Format

Machine-readable format suitable for integration with other tools:

```json
{
  "total_files": 45,
  "files_with_issues": 1,
  "reports": [
    {
      "file_path": "templates/backend/go-gin/internal/middleware/auth.go.tmpl",
      "missing_imports": ["time"],
      "used_functions": [
        {
          "function": "time.Now",
          "line": 25,
          "column": 10,
          "required_package": "time"
        }
      ],
      "current_imports": [
        {
          "package": "fmt",
          "alias": "",
          "is_stdlib": true
        }
      ]
    }
  ],
  "summary": {
    "time": 1
  }
}
```

#### Summary Format

High-level overview of analysis results:

```
Import Detection Summary
=======================

Total files analyzed: 45
Files with issues: 1
Success rate: 97.8%

Most commonly missing packages:
  time                          1 files
```

## How It Works

### 1. Template Preprocessing

The tool preprocesses Go template files to handle template syntax that would break Go parsing:

- Replaces template variables (e.g., `{{.Name}}` → `TemplateName`)
- Removes template control structures (`{{if}}`, `{{range}}`, etc.)
- Converts remaining template expressions to valid Go literals

### 2. AST Analysis

Uses Go's `go/parser` and `go/ast` packages to:

- Parse preprocessed Go code into an Abstract Syntax Tree
- Extract existing import statements
- Find function calls and selector expressions
- Map function calls to required packages

### 3. Import Detection

Compares used functions against a comprehensive mapping of standard library functions to determine missing imports.

### 4. Report Generation

Generates structured reports with:

- File paths and line numbers
- Missing import packages
- Function usage details
- Current import statements
- Error information for unparseable files

## Function Mapping

The tool includes mappings for common standard library packages:

- **time**: `time.Now`, `time.Since`, `time.Parse`, etc.
- **fmt**: `fmt.Printf`, `fmt.Sprintf`, `fmt.Errorf`, etc.
- **strings**: `strings.Contains`, `strings.Split`, `strings.Join`, etc.
- **encoding/json**: `json.Marshal`, `json.Unmarshal`, etc.
- **net/http**: `http.Get`, `http.Post`, `http.StatusOK`, etc.
- **os**: `os.Getenv`, `os.Open`, `os.Create`, etc.
- And many more...

## Integration

### CI/CD Integration

The tool can be integrated into CI/CD pipelines to automatically detect import issues:

```bash
# In your CI script
./import-detector -dir templates -format json -output import-report.json
if [ $? -ne 0 ]; then
  echo "Import detection failed"
  exit 1
fi

# Check if any issues were found
issues=$(jq '.files_with_issues' import-report.json)
if [ "$issues" -gt 0 ]; then
  echo "Found $issues files with missing imports"
  exit 1
fi
```

### Pre-commit Hooks

Add to your pre-commit configuration to catch issues early:

```yaml
- repo: local
  hooks:
    - id: import-detector
      name: Template Import Detection
      entry: ./import-detector
      args: [-dir, templates, -format, summary]
      language: system
      pass_filenames: false
```

## Building

```bash
# Build the import detector
go build -o import-detector ./cmd/import-detector

# Run tests
go test ./pkg/template/...

# Install globally
go install ./cmd/import-detector
```

## Limitations

- Currently focuses on standard library packages
- Template preprocessing may not handle all edge cases
- Requires valid Go syntax after preprocessing
- Does not detect unused imports (use `goimports` for that)

## Contributing

When adding new function mappings:

1. Add the function to the `buildFunctionPackageMap()` function
2. Include appropriate test cases
3. Update documentation if adding new package categories
