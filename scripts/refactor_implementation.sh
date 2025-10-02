#!/bin/bash

# Code Splitting Implementation Script
# This script helps implement the refactoring plan systematically

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
check_project_root() {
    if [[ ! -f "go.mod" ]] || [[ ! -d "pkg" ]]; then
        log_error "This script must be run from the project root directory"
        exit 1
    fi
}

# Create backup of current state
create_backup() {
    local backup_dir="backups/refactor_$(date +%Y%m%d_%H%M%S)"
    log_info "Creating backup in $backup_dir"
    
    mkdir -p "$backup_dir"
    cp -r pkg/ "$backup_dir/"
    cp -r internal/ "$backup_dir/"
    cp -r cmd/ "$backup_dir/"
    
    log_success "Backup created in $backup_dir"
}

# Run tests to ensure current functionality
run_tests() {
    log_info "Running tests to verify current functionality"
    
    if ! go test ./... -v; then
        log_error "Tests are failing. Please fix tests before refactoring."
        exit 1
    fi
    
    log_success "All tests passing"
}

# Phase 1: CLI Package Split
phase1_cli_split() {
    log_info "Phase 1: Splitting CLI package"
    
    # Create new directories
    mkdir -p pkg/cli/commands
    mkdir -p pkg/cli/utils
    
    # Create output utilities file
    cat > pkg/cli/output.go << 'EOF'
package cli

import (
    "fmt"
    "os"
    "time"
)

// Color constants for beautiful CLI output
const (
    ColorReset  = "\033[0m"
    ColorRed    = "\033[31m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorBlue   = "\033[34m"
    ColorPurple = "\033[35m"
    ColorCyan   = "\033[36m"
    ColorWhite  = "\033[37m"
    ColorBold   = "\033[1m"
    ColorDim    = "\033[2m"
)

// Color helper functions for beautiful output
func (c *CLI) colorize(color, text string) string {
    if c.quietMode {
        return text // No colors in quiet mode
    }
    return color + text + ColorReset
}

func (c *CLI) success(text string) string {
    return c.colorize(ColorGreen+ColorBold, text)
}

func (c *CLI) info(text string) string {
    return c.colorize(ColorBlue, text)
}

func (c *CLI) warning(text string) string {
    return c.colorize(ColorYellow, text)
}

func (c *CLI) error(text string) string {
    return c.colorize(ColorRed+ColorBold, text)
}

func (c *CLI) highlight(text string) string {
    return c.colorize(ColorCyan+ColorBold, text)
}

func (c *CLI) dim(text string) string {
    return c.colorize(ColorDim, text)
}

// Output methods
func (c *CLI) VerboseOutput(format string, args ...interface{}) {
    if c.verboseMode && !c.quietMode {
        fmt.Printf(format+"\n", args...)
    }
}

func (c *CLI) DebugOutput(format string, args ...interface{}) {
    if c.debugMode && !c.quietMode {
        fmt.Printf("🐛 "+format+"\n", args...)
    }
}

func (c *CLI) QuietOutput(format string, args ...interface{}) {
    if !c.quietMode {
        fmt.Printf(format+"\n", args...)
    }
}

func (c *CLI) ErrorOutput(format string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
}

func (c *CLI) WarningOutput(format string, args ...interface{}) {
    if !c.quietMode {
        fmt.Printf("⚠️  "+format+"\n", args...)
    }
}

func (c *CLI) SuccessOutput(format string, args ...interface{}) {
    if !c.quietMode {
        fmt.Printf(format+"\n", args...)
    }
}

func (c *CLI) PerformanceOutput(operation string, duration time.Duration, metrics map[string]interface{}) {
    if c.debugMode && !c.quietMode {
        fmt.Printf("⚡ %s completed in %v\n", operation, duration)
        if len(metrics) > 0 {
            fmt.Printf("   Metrics: %+v\n", metrics)
        }
    }
}
EOF

    log_success "Created pkg/cli/output.go"
    
    # Create flags utilities
    cat > pkg/cli/flags.go << 'EOF'
package cli

import (
    "github.com/spf13/cobra"
)

// setupGlobalFlags adds global flags that apply to all commands
func (c *CLI) setupGlobalFlags() {
    c.rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging with detailed operation information")
    c.rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-error output (quiet mode)")
    c.rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode with detailed diagnostic information")
    c.rootCmd.PersistentFlags().String("log-level", "info", "Set logging level (debug, info, warn, error)")
    c.rootCmd.PersistentFlags().String("log-format", "text", "Set log format (text, json)")
}

// handleGlobalFlags processes global flags that apply to all commands
func (c *CLI) handleGlobalFlags(cmd *cobra.Command) error {
    // Get global flags
    verbose, _ := cmd.Flags().GetBool("verbose")
    quiet, _ := cmd.Flags().GetBool("quiet")
    debug, _ := cmd.Flags().GetBool("debug")
    
    c.verboseMode = verbose
    c.quietMode = quiet
    c.debugMode = debug
    
    return nil
}
EOF

    log_success "Created pkg/cli/flags.go"
    
    log_success "Phase 1 completed: CLI package split initiated"
}

# Phase 2: Audit Engine Split
phase2_audit_split() {
    log_info "Phase 2: Splitting Audit Engine"
    
    # Create audit subdirectories
    mkdir -p pkg/audit/security
    mkdir -p pkg/audit/quality
    mkdir -p pkg/audit/license
    mkdir -p pkg/audit/performance
    mkdir -p pkg/audit/utils
    
    # Create security scanner
    cat > pkg/audit/security/scanner.go << 'EOF'
package security

import (
    "github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Scanner handles security-related auditing
type Scanner struct {
    rules []interfaces.AuditRule
}

// NewScanner creates a new security scanner
func NewScanner() *Scanner {
    return &Scanner{
        rules: getSecurityRules(),
    }
}

// ScanForVulnerabilities scans for security vulnerabilities
func (s *Scanner) ScanForVulnerabilities(projectPath string) ([]interfaces.Vulnerability, error) {
    // Implementation will be moved from main engine
    return []interfaces.Vulnerability{}, nil
}

// getSecurityRules returns security-specific audit rules
func getSecurityRules() []interfaces.AuditRule {
    return []interfaces.AuditRule{
        {
            ID:          "security-001",
            Name:        "No hardcoded secrets",
            Description: "Check for hardcoded secrets in source code",
            Category:    interfaces.AuditCategorySecurity,
            Type:        interfaces.AuditCategorySecurity,
            Severity:    interfaces.AuditSeverityCritical,
            Enabled:     true,
            Pattern:     `(?i)(password|secret|key|token)\s*[:=]\s*["'][^"']+["']`,
            FileTypes:   []string{".go", ".js", ".ts", ".py", ".java", ".cs"},
        },
    }
}
EOF

    log_success "Created pkg/audit/security/scanner.go"
    
    log_success "Phase 2 completed: Audit engine split initiated"
}

# Phase 3: Template Manager Split
phase3_template_split() {
    log_info "Phase 3: Splitting Template Manager"
    
    # Create template subdirectories
    mkdir -p pkg/template/processor
    mkdir -p pkg/template/metadata
    
    # Create discovery module
    cat > pkg/template/discovery.go << 'EOF'
package template

import (
    "github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
    "github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TemplateDiscovery handles template discovery and scanning
type TemplateDiscovery struct {
    cache map[string]*models.TemplateInfo
}

// NewTemplateDiscovery creates a new template discovery instance
func NewTemplateDiscovery() *TemplateDiscovery {
    return &TemplateDiscovery{
        cache: make(map[string]*models.TemplateInfo),
    }
}

// DiscoverTemplates discovers available templates
func (td *TemplateDiscovery) DiscoverTemplates() ([]*models.TemplateInfo, error) {
    // Implementation will be moved from manager
    return []*models.TemplateInfo{}, nil
}

// ApplyFilters applies filters to template list
func (td *TemplateDiscovery) ApplyFilters(templates []*models.TemplateInfo, filter interfaces.TemplateFilter) []*models.TemplateInfo {
    // Implementation will be moved from manager
    return templates
}
EOF

    log_success "Created pkg/template/discovery.go"
    
    log_success "Phase 3 completed: Template manager split initiated"
}

# Phase 4: Validation and Cache Split
phase4_validation_cache_split() {
    log_info "Phase 4: Splitting Validation and Cache components"
    
    # Create validation format directories
    mkdir -p pkg/validation/formats
    
    # Create JSON validator
    cat > pkg/validation/formats/json.go << 'EOF'
package formats

import (
    "encoding/json"
    "github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
    "github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// JSONValidator handles JSON file validation
type JSONValidator struct{}

// NewJSONValidator creates a new JSON validator
func NewJSONValidator() *JSONValidator {
    return &JSONValidator{}
}

// ValidateFile validates a JSON file
func (jv *JSONValidator) ValidateFile(filePath string) (*interfaces.ConfigValidationResult, error) {
    content, err := utils.SafeReadFile(filePath)
    if err != nil {
        return nil, err
    }
    
    var jsonData interface{}
    if err := json.Unmarshal(content, &jsonData); err != nil {
        return &interfaces.ConfigValidationResult{
            Valid: false,
            Errors: []interfaces.ConfigValidationError{
                {
                    Field:   "syntax",
                    Message: "Invalid JSON syntax: " + err.Error(),
                    Type:    "syntax_error",
                },
            },
        }, nil
    }
    
    return &interfaces.ConfigValidationResult{
        Valid: true,
    }, nil
}
EOF

    log_success "Created pkg/validation/formats/json.go"
    
    # Create cache storage module
    cat > pkg/cache/storage.go << 'EOF'
package cache

import (
    "github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// Storage handles cache storage operations
type Storage struct {
    basePath string
}

// NewStorage creates a new cache storage instance
func NewStorage(basePath string) *Storage {
    return &Storage{
        basePath: basePath,
    }
}

// Store stores data in cache
func (s *Storage) Store(key string, data []byte) error {
    // Implementation will be moved from manager
    return nil
}

// Retrieve retrieves data from cache
func (s *Storage) Retrieve(key string) ([]byte, error) {
    // Implementation will be moved from manager
    return nil, nil
}

// Delete removes data from cache
func (s *Storage) Delete(key string) error {
    // Implementation will be moved from manager
    return nil
}
EOF

    log_success "Created pkg/cache/storage.go"
    
    log_success "Phase 4 completed: Validation and cache split initiated"
}

# Phase 5: Dead Code Analysis
phase5_dead_code_analysis() {
    log_info "Phase 5: Analyzing for dead code"
    
    # Create a script to find potentially unused functions
    cat > scripts/find_unused_functions.sh << 'EOF'
#!/bin/bash

# Find potentially unused exported functions
echo "=== Potentially Unused Exported Functions ==="

# Find all exported functions
grep -r "^func [A-Z]" pkg/ --include="*.go" | while read -r line; do
    file=$(echo "$line" | cut -d: -f1)
    func_line=$(echo "$line" | cut -d: -f2-)
    func_name=$(echo "$func_line" | sed 's/^func \([A-Z][a-zA-Z0-9_]*\).*/\1/')
    
    # Search for usage of this function (excluding the definition)
    usage_count=$(grep -r "\b$func_name\b" pkg/ cmd/ internal/ --include="*.go" | grep -v "$file:" | wc -l)
    
    if [ "$usage_count" -eq 0 ]; then
        echo "Potentially unused: $func_name in $file"
    fi
done

echo ""
echo "=== Potentially Unused Imports ==="

# Find files with potentially unused imports
find pkg/ -name "*.go" -exec go list -f '{{.ImportPath}}: {{.Imports}}' {} \; 2>/dev/null | grep -v "command-line-arguments"
EOF

    chmod +x scripts/find_unused_functions.sh
    
    log_info "Running dead code analysis..."
    ./scripts/find_unused_functions.sh > dead_code_analysis.txt
    
    log_success "Dead code analysis completed. Results in dead_code_analysis.txt"
}

# Run all phases
run_all_phases() {
    log_info "Starting complete refactoring process"
    
    check_project_root
    create_backup
    run_tests
    
    phase1_cli_split
    phase2_audit_split
    phase3_template_split
    phase4_validation_cache_split
    phase5_dead_code_analysis
    
    log_success "All phases completed successfully!"
    log_info "Next steps:"
    log_info "1. Review the created files and move implementation from original files"
    log_info "2. Update imports and references"
    log_info "3. Run tests after each component migration"
    log_info "4. Review dead_code_analysis.txt for cleanup opportunities"
}

# Main execution
case "${1:-all}" in
    "1"|"cli")
        check_project_root
        phase1_cli_split
        ;;
    "2"|"audit")
        check_project_root
        phase2_audit_split
        ;;
    "3"|"template")
        check_project_root
        phase3_template_split
        ;;
    "4"|"validation")
        check_project_root
        phase4_validation_cache_split
        ;;
    "5"|"deadcode")
        check_project_root
        phase5_dead_code_analysis
        ;;
    "all")
        run_all_phases
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [phase]"
        echo "Phases:"
        echo "  1, cli        - Split CLI package"
        echo "  2, audit      - Split audit engine"
        echo "  3, template   - Split template manager"
        echo "  4, validation - Split validation and cache"
        echo "  5, deadcode   - Analyze dead code"
        echo "  all           - Run all phases (default)"
        echo "  help          - Show this help"
        ;;
    *)
        log_error "Unknown phase: $1"
        log_info "Use '$0 help' for usage information"
        exit 1
        ;;
esac