# Template Maintenance Quick Reference

## Import Organization Template

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

## Common Function → Import Mapping

| Function | Import | Function | Import |
|----------|--------|----------|--------|
| `time.Now()` | `"time"` | `fmt.Sprintf()` | `"fmt"` |
| `strings.Contains()` | `"strings"` | `strconv.Atoi()` | `"strconv"` |
| `context.Background()` | `"context"` | `log.Printf()` | `"log"` |
| `http.StatusOK` | `"net/http"` | `json.Marshal()` | `"encoding/json"` |
| `errors.New()` | `"errors"` | `os.Getenv()` | `"os"` |

## Validation Commands

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

## Pre-Commit Checklist

- [ ] All used functions have corresponding imports
- [ ] Imports are properly grouped and ordered
- [ ] No unused imports remain
- [ ] Template compiles with test data
- [ ] Integration tests pass

## Common Fixes

### Missing Time Import

```go
// Add to imports
"time"

// Usage
if time.Now().After(expiry) { ... }
```

### Missing Context Import

```go
// Add to imports
"context"

// Usage
func Handler(ctx context.Context) { ... }
```

### Missing HTTP Import

```go
// Add to imports
"net/http"

// Usage
c.JSON(http.StatusOK, response)
```

For detailed guidelines, see [TEMPLATE_MAINTENANCE.md](./TEMPLATE_MAINTENANCE.md)
