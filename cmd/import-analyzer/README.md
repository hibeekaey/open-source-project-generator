# Import Analyzer Tool

The Import Analyzer is a command-line tool that analyzes Go template files to detect missing import statements. It helps ensure that generated Go code will compile successfully without import-related errors.

## Features

- **Automated Detection**: Scans template files to identify functions and types that require specific imports
- **Comprehensive Mapping**: Includes mappings for common standard library packages (time, fmt, strings, etc.)
- **Template-Aware**: Handles Go template syntax including variables, conditionals, and loops
- **Structured Reporting**: Generates detailed reports in text or JSON format
- **Extensible**: Supports custom function-to-package mappings

## Usage

### Basic Usage

```bash
# Analyze templates in the default directory
./import-analyzer

# Analyze templates in a specific directory
./import-analyzer -dir /path/to/templates

# Generate JSON output
./import-analyzer -json

# Save report to file
./import-analyzer -output report.txt

# Enable verbose output
./import-analyzer -verbose
```

### Command Line Options

- `-dir string`: Directory containing template files to analyze (default: "templates")
- `-output string`: Output file for the report (default: stdout)
- `-json`: Output report in JSON format instead of text
- `-verbose`: Enable verbose output with additional information

### Exit Codes

- `0`: No issues found
- `1`: Missing imports or analysis errors detected

## Example Output

### Text Format

```
Template Import Analysis Report
===============================

Generated: 2023-12-07T10:30:00Z
Total Files Analyzed: 15
Files with Issues: 3
Total Missing Imports: 7

Most Common Missing Imports:
  time: 2 files
  strings: 2 files
  encoding/json: 1 files

Detailed Results:
-----------------

File: templates/backend/go-gin/internal/middleware/auth.go.tmpl
  Missing Imports:
    - time
  Function Usage:
    - time.Now (line 12) requires time
    - time.Since (line 18) requires time

File: templates/backend/go-gin/internal/handlers/user.go.tmpl
  Missing Imports:
    - encoding/json
    - strings
  Function Usage:
    - json.Marshal (line 25) requires encoding/json
    - strings.ToLower (line 30) requires strings
```

### JSON Format

```json
{
  "reports": [
    {
      "file_path": "templates/backend/go-gin/internal/middleware/auth.go.tmpl",
      "missing_imports": ["time"],
      "used_functions": [
        {
          "function": "time.Now",
          "line": 12,
          "column": 10,
          "required_package": "time"
        }
      ],
      "current_imports": ["net/http", "github.com/gin-gonic/gin"],
      "has_errors": false,
      "errors": []
    }
  ],
  "summary": {
    "total_files": 15,
    "files_with_issues": 3,
    "total_missing_imports": 7,
    "most_common_missing": {
      "time": 2,
      "strings": 2,
      "encoding/json": 1
    }
  },
  "generated_at": "2023-12-07T10:30:00Z"
}
```

## Supported Packages

The analyzer includes built-in mappings for common Go standard library packages:

### Standard Library Packages

- **time**: `time.Now`, `time.Since`, `time.Parse`, `time.Sleep`, etc.
- **fmt**: `fmt.Printf`, `fmt.Sprintf`, `fmt.Errorf`, etc.
- **strings**: `strings.Contains`, `strings.Split`, `strings.ToLower`, etc.
- **strconv**: `strconv.Atoi`, `strconv.Itoa`, `strconv.ParseInt`, etc.
- **encoding/json**: `json.Marshal`, `json.Unmarshal`, etc.
- **net/http**: `http.Get`, `http.Post`, `http.ListenAndServe`, etc.
- **context**: `context.Background`, `context.WithTimeout`, etc.
- **errors**: `errors.New`, `errors.Is`, `errors.As`, etc.
- **os**: `os.Open`, `os.Getenv`, `os.Exit`, etc.
- **io**: `io.Copy`, `io.ReadAll`, etc.
- **log**: `log.Printf`, `log.Fatal`, etc.

### Type Mappings

- **time.Time**, **time.Duration**: requires `time`
- **http.Request**, **http.Response**: requires `net/http`
- **context.Context**: requires `context`
- **sql.DB**: requires `database/sql`
- **url.URL**: requires `net/url`

## Template Syntax Handling

The analyzer handles Go template syntax by:

1. **Variable Replacement**: Replaces `{{.Variable}}` with placeholder values
2. **Conditional Processing**: Converts `{{if .Condition}}` to valid Go syntax
3. **Loop Processing**: Converts `{{range .Items}}` to valid Go syntax
4. **Comment Removal**: Removes template comments `{{/* comment */}}`

## Integration

### CI/CD Integration

```bash
# In your CI pipeline
./import-analyzer -dir templates -json -output import-analysis.json
if [ $? -ne 0 ]; then
  echo "Import analysis failed - missing imports detected"
  exit 1
fi
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit
./import-analyzer -dir templates
if [ $? -ne 0 ]; then
  echo "Commit rejected: template files have missing imports"
  echo "Run './import-analyzer -dir templates' to see details"
  exit 1
fi
```

## Building

```bash
# Build the tool
go build -o import-analyzer ./cmd/import-analyzer

# Run tests
go test ./pkg/template/...

# Run with coverage
go test -cover ./pkg/template/...
```

## Limitations

1. **Static Analysis Only**: Cannot detect dynamic function calls or reflection-based usage
2. **Standard Library Focus**: Built-in mappings focus on standard library packages
3. **Template Syntax**: Complex template logic may not be perfectly parsed
4. **Third-party Packages**: Requires manual mapping for third-party package functions

## Extending the Analyzer

The analyzer can be extended programmatically:

```go
analyzer := template.NewImportAnalyzer()

// Add custom function mappings
analyzer.AddFunctionMapping("custom.Function", "github.com/example/custom")

// Add custom type mappings  
analyzer.AddTypeMapping("custom.Type", "github.com/example/custom")
```
