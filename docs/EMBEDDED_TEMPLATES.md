# Embedded Templates

The Open Source Template Generator uses Go's `embed` package to bundle all templates directly into the binary, making it completely self-contained and portable.

## How It Works

Templates are embedded at compile time using the `//go:embed` directive:

```go
//go:embed templates/*
var embeddedTemplates embed.FS
```

This means:

- **No external dependencies**: The binary contains all templates
- **Portable**: Works anywhere without needing template files
- **Secure**: Templates can't be modified after compilation
- **Fast**: No file system lookups needed

## Benefits

### 1. **Complete Portability**

```bash
# Works anywhere - no setup needed
./generator generate
```

### 2. **No Installation Complexity**

- Single binary distribution
- No need to manage template directories
- No path resolution issues

### 3. **Consistent Behavior**

- Same templates on all systems
- No version mismatches between binary and templates
- Reproducible builds

## Development

### Running from Source

```bash
go run cmd/generator/main.go generate
```

Templates are embedded from `pkg/template/templates/`

### Building

```bash
go build -o generator cmd/generator/main.go
```

All templates are automatically included in the binary.

### Template Updates

When templates are updated:

1. Copy changes to `pkg/template/templates/`
2. Rebuild the binary
3. Templates are automatically embedded

## Template Structure

The embedded templates are stored in the package directory:

```text
pkg/template/templates/
├── base/                  # Base project files
├── frontend/              # Frontend templates
├── backend/               # Backend templates
├── mobile/                # Mobile templates
└── infrastructure/        # Infrastructure templates
```

**Note**: The original `templates/` directory in the project root has been removed since all templates are now embedded in the binary.

## Implementation Details

### Embedded Engine

The `EmbeddedEngine` handles template processing from the embedded filesystem:

```go
type EmbeddedEngine struct {
    funcMap template.FuncMap
    fs      fs.FS
}
```

### Template Processing

1. Templates are loaded from `embed.FS`
2. Processed using Go's `text/template` package
3. Output written to the target directory

### Error Handling

- Missing templates are gracefully skipped
- Clear error messages for template processing issues
- Fallback behavior for optional components

## Migration from File-based Templates

The embedded system maintains full compatibility:

- Same template syntax and functions
- Same directory structure in output
- Same configuration options
- Same CLI interface

## Distribution

### Single Binary

```bash
# Just distribute the binary - no other files needed
./generator generate
```

### Package Managers

- DEB/RPM packages only need the binary
- No template directories to manage
- Simplified installation scripts

### Docker

```dockerfile
FROM scratch
COPY generator /generator
ENTRYPOINT ["/generator"]
```

## Troubleshooting

### Template Not Found

If you see template-related errors:

1. Ensure you're using the embedded engine
2. Check that templates exist in `pkg/template/templates/`
3. Rebuild the binary to include template changes

### Development vs Production

- Development: Templates loaded from `pkg/template/templates/`
- Production: Templates embedded in binary
- Both use the same `EmbeddedEngine` for consistency
