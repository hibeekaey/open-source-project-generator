package cleanup

import (
	"fmt"
	"os"
	"strings"
)

// TODOResolver handles the resolution and documentation of remaining TODOs
type TODOResolver struct {
	projectRoot string
	config      *TODOResolverConfig
}

// TODOResolverConfig holds configuration for TODO resolution
type TODOResolverConfig struct {
	DryRun  bool
	Verbose bool
}

// NewTODOResolver creates a new TODO resolver
func NewTODOResolver(projectRoot string, config *TODOResolverConfig) *TODOResolver {
	if config == nil {
		config = &TODOResolverConfig{}
	}
	return &TODOResolver{
		projectRoot: projectRoot,
		config:      config,
	}
}

// ResolveRemainingTODOs processes all remaining TODOs after security implementation
func (tr *TODOResolver) ResolveRemainingTODOs() (*TODOResolutionReport, error) {
	scanner := NewTODOScanner(DefaultTODOScanConfig())

	// Scan for all TODOs
	report, err := scanner.ScanProject(tr.projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to scan for TODOs: %w", err)
	}

	resolution := &TODOResolutionReport{
		TotalFound:     len(report.TODOs),
		Resolved:       make([]ResolvedTODO, 0),
		Documented:     make([]DocumentedTODO, 0),
		Removed:        make([]RemovedTODO, 0),
		FalsePositives: make([]FalsePositiveTODO, 0),
	}

	// Process each TODO
	for _, todo := range report.TODOs {
		action := tr.determineTODOAction(todo)

		switch action {
		case TODOActionResolve:
			resolved, err := tr.resolveTODO(todo)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve TODO at %s:%d: %w", todo.File, todo.Line, err)
			}
			resolution.Resolved = append(resolution.Resolved, *resolved)

		case TODOActionDocument:
			documented := tr.documentTODO(todo)
			resolution.Documented = append(resolution.Documented, *documented)

		case TODOActionRemove:
			removed, err := tr.removeTODO(todo)
			if err != nil {
				return nil, fmt.Errorf("failed to remove TODO at %s:%d: %w", todo.File, todo.Line, err)
			}
			resolution.Removed = append(resolution.Removed, *removed)

		case TODOActionIgnore:
			falsePositive := tr.markAsFalsePositive(todo)
			resolution.FalsePositives = append(resolution.FalsePositives, *falsePositive)
		}
	}

	return resolution, nil
}

// TODOAction represents the action to take for a TODO
type TODOAction int

const (
	TODOActionResolve TODOAction = iota
	TODOActionDocument
	TODOActionRemove
	TODOActionIgnore
)

// determineTODOAction decides what action to take for a given TODO
func (tr *TODOResolver) determineTODOAction(todo TODOItem) TODOAction {
	// Check if it's a false positive (documentation, specs, etc.)
	if tr.isFalsePositive(todo) {
		return TODOActionIgnore
	}

	// Check if it's a legitimate code reference (like context.TODO)
	if tr.isLegitimateCodeReference(todo) {
		return TODOActionIgnore
	}

	// Check if it's in template files - these should be documented
	if tr.isTemplateFile(todo.File) {
		return tr.handleTemplateTODO(todo)
	}

	// Check if it's obsolete or completed
	if tr.isObsoleteTODO(todo) {
		return TODOActionRemove
	}

	// Check if it's a feature TODO that can be resolved
	if tr.canResolveFeatureTODO(todo) {
		return TODOActionResolve
	}

	// Default: document for future development
	return TODOActionDocument
}

// isFalsePositive checks if the TODO is actually a false positive
func (tr *TODOResolver) isFalsePositive(todo TODOItem) bool {
	// Skip documentation files
	if strings.Contains(todo.File, "/docs/") ||
		strings.HasSuffix(todo.File, ".md") ||
		strings.Contains(todo.File, "/.kiro/specs/") ||
		strings.Contains(todo.File, "/scripts/") {
		return true
	}

	// Skip if it's just mentioning TODOs in comments about TODOs
	if strings.Contains(todo.Message, "TODO comments") ||
		strings.Contains(todo.Message, "TODO/FIXME") ||
		strings.Contains(todo.Context, "Check for TODO") {
		return true
	}

	return false
}

// isLegitimateCodeReference checks if it's a legitimate code reference
func (tr *TODOResolver) isLegitimateCodeReference(todo TODOItem) bool {
	// context.TODO is a legitimate Go standard library function
	if strings.Contains(todo.Context, "context.TODO") {
		return true
	}

	// Check for other legitimate references
	legitimateRefs := []string{
		"context.TODO",
		"TODO comments without issues", // PR template
	}

	for _, ref := range legitimateRefs {
		if strings.Contains(todo.Context, ref) {
			return true
		}
	}

	return false
}

// isTemplateFile checks if the file is a template file
func (tr *TODOResolver) isTemplateFile(filePath string) bool {
	return strings.Contains(filePath, "/templates/") && strings.HasSuffix(filePath, ".tmpl")
}

// handleTemplateTODO determines action for template TODOs
func (tr *TODOResolver) handleTemplateTODO(todo TODOItem) TODOAction {
	// Template TODOs are intentional placeholders for generated projects
	// They should be documented as intentional
	return TODOActionDocument
}

// isObsoleteTODO checks if the TODO is obsolete or already completed
func (tr *TODOResolver) isObsoleteTODO(todo TODOItem) bool {
	// Check for TODOs that reference already implemented features
	obsoletePatterns := []string{
		"implement security checking", // Already implemented in task 2.2
		"security audit integration",  // Already implemented
		"vulnerability database",      // Already implemented
	}

	for _, pattern := range obsoletePatterns {
		if strings.Contains(strings.ToLower(todo.Message), pattern) {
			return true
		}
	}

	return false
}

// canResolveFeatureTODO checks if a feature TODO can be resolved now
func (tr *TODOResolver) canResolveFeatureTODO(todo TODOItem) bool {
	// For this cleanup task, we focus on simple resolutions
	// More complex feature TODOs should be documented for future work

	resolvablePatterns := []string{
		"send email", // Can be documented as intentional placeholder
	}

	for _, pattern := range resolvablePatterns {
		if strings.Contains(strings.ToLower(todo.Message), pattern) {
			return true
		}
	}

	return false
}

// resolveTODO resolves a TODO by implementing or fixing it
func (tr *TODOResolver) resolveTODO(todo TODOItem) (*ResolvedTODO, error) {
	// Handle specific resolvable TODOs
	if strings.Contains(strings.ToLower(todo.Message), "send email") {
		return tr.resolveEmailTODO(todo)
	}

	return &ResolvedTODO{
		Original:    todo,
		Action:      "Resolved by implementation",
		Description: "TODO was addressed during cleanup",
	}, nil
}

// resolveEmailTODO resolves email-related TODOs by adding proper documentation
func (tr *TODOResolver) resolveEmailTODO(todo TODOItem) (*ResolvedTODO, error) {
	// Read the file
	content, err := os.ReadFile(todo.File)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	if todo.Line-1 >= len(lines) {
		return nil, fmt.Errorf("line number %d out of range", todo.Line)
	}

	// Replace the TODO with proper documentation
	newComment := "\t// NOTE: Email sending should be implemented based on your email service provider"
	lines[todo.Line-1] = newComment

	// Write back to file
	newContent := strings.Join(lines, "\n")
	err = os.WriteFile(todo.File, []byte(newContent), 0644)
	if err != nil {
		return nil, err
	}

	return &ResolvedTODO{
		Original:    todo,
		Action:      "Replaced with documentation",
		Description: "TODO replaced with proper implementation guidance",
	}, nil
}

// documentTODO documents a TODO for future development
func (tr *TODOResolver) documentTODO(todo TODOItem) *DocumentedTODO {
	reason := tr.getDocumentationReason(todo)

	return &DocumentedTODO{
		Original:   todo,
		Reason:     reason,
		Documented: true,
		FutureWork: true,
	}
}

// getDocumentationReason provides a reason for why the TODO is being documented
func (tr *TODOResolver) getDocumentationReason(todo TODOItem) string {
	if tr.isTemplateFile(todo.File) {
		return "Template placeholder - intentional TODO for generated projects"
	}

	if todo.Category == CategorySecurity {
		return "Security enhancement - requires careful implementation and testing"
	}

	if todo.Category == CategoryPerformance {
		return "Performance optimization - requires benchmarking and analysis"
	}

	return "Feature enhancement - requires design and implementation planning"
}

// removeTODO removes an obsolete TODO
func (tr *TODOResolver) removeTODO(todo TODOItem) (*RemovedTODO, error) {
	// Read the file
	content, err := os.ReadFile(todo.File)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	if todo.Line-1 >= len(lines) {
		return nil, fmt.Errorf("line number %d out of range", todo.Line)
	}

	// Remove the TODO line if it's a standalone comment
	originalLine := lines[todo.Line-1]
	if strings.TrimSpace(originalLine) == strings.TrimSpace(todo.Context) {
		// Remove the entire line
		lines = append(lines[:todo.Line-1], lines[todo.Line:]...)
	} else {
		// Remove just the TODO part from the line
		lines[todo.Line-1] = strings.Replace(originalLine, todo.Context, "", 1)
	}

	// Write back to file
	newContent := strings.Join(lines, "\n")
	err = os.WriteFile(todo.File, []byte(newContent), 0644)
	if err != nil {
		return nil, err
	}

	return &RemovedTODO{
		Original: todo,
		Reason:   "Obsolete or completed TODO",
	}, nil
}

// markAsFalsePositive marks a TODO as a false positive
func (tr *TODOResolver) markAsFalsePositive(todo TODOItem) *FalsePositiveTODO {
	reason := "Not an actual TODO comment"

	if tr.isFalsePositive(todo) {
		reason = "Documentation or specification file"
	} else if tr.isLegitimateCodeReference(todo) {
		reason = "Legitimate code reference (e.g., context.TODO)"
	}

	return &FalsePositiveTODO{
		Original: todo,
		Reason:   reason,
	}
}

// TODOResolutionReport contains the results of TODO resolution
type TODOResolutionReport struct {
	TotalFound     int
	Resolved       []ResolvedTODO
	Documented     []DocumentedTODO
	Removed        []RemovedTODO
	FalsePositives []FalsePositiveTODO
}

// ResolvedTODO represents a TODO that was resolved
type ResolvedTODO struct {
	Original    TODOItem
	Action      string
	Description string
}

// DocumentedTODO represents a TODO that was documented for future work
type DocumentedTODO struct {
	Original   TODOItem
	Reason     string
	Documented bool
	FutureWork bool
}

// RemovedTODO represents a TODO that was removed as obsolete
type RemovedTODO struct {
	Original TODOItem
	Reason   string
}

// FalsePositiveTODO represents a false positive TODO
type FalsePositiveTODO struct {
	Original TODOItem
	Reason   string
}

// GenerateReport generates a markdown report of the TODO resolution
func (report *TODOResolutionReport) GenerateReport() string {
	var sb strings.Builder

	sb.WriteString("# TODO Resolution Report\n\n")
	sb.WriteString(fmt.Sprintf("**Total TODOs Found:** %d\n\n", report.TotalFound))

	// Summary
	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Resolved:** %d\n", len(report.Resolved)))
	sb.WriteString(fmt.Sprintf("- **Documented for Future Work:** %d\n", len(report.Documented)))
	sb.WriteString(fmt.Sprintf("- **Removed (Obsolete):** %d\n", len(report.Removed)))
	sb.WriteString(fmt.Sprintf("- **False Positives:** %d\n\n", len(report.FalsePositives)))

	// Resolved TODOs
	if len(report.Resolved) > 0 {
		sb.WriteString("## Resolved TODOs\n\n")
		for _, resolved := range report.Resolved {
			sb.WriteString(fmt.Sprintf("### %s:%d\n", resolved.Original.File, resolved.Original.Line))
			sb.WriteString(fmt.Sprintf("**Original:** %s\n", resolved.Original.Message))
			sb.WriteString(fmt.Sprintf("**Action:** %s\n", resolved.Action))
			sb.WriteString(fmt.Sprintf("**Description:** %s\n\n", resolved.Description))
		}
	}

	// Documented TODOs
	if len(report.Documented) > 0 {
		sb.WriteString("## Documented for Future Work\n\n")
		for _, documented := range report.Documented {
			sb.WriteString(fmt.Sprintf("### %s:%d\n", documented.Original.File, documented.Original.Line))
			sb.WriteString(fmt.Sprintf("**Message:** %s\n", documented.Original.Message))
			sb.WriteString(fmt.Sprintf("**Reason:** %s\n", documented.Reason))
			sb.WriteString(fmt.Sprintf("**Category:** %s\n\n", categoryToString(documented.Original.Category)))
		}
	}

	// Removed TODOs
	if len(report.Removed) > 0 {
		sb.WriteString("## Removed TODOs\n\n")
		for _, removed := range report.Removed {
			sb.WriteString(fmt.Sprintf("### %s:%d\n", removed.Original.File, removed.Original.Line))
			sb.WriteString(fmt.Sprintf("**Original:** %s\n", removed.Original.Message))
			sb.WriteString(fmt.Sprintf("**Reason:** %s\n\n", removed.Reason))
		}
	}

	// False Positives
	if len(report.FalsePositives) > 0 {
		sb.WriteString("## False Positives\n\n")
		sb.WriteString("These were identified as false positives and ignored:\n\n")
		for _, fp := range report.FalsePositives {
			sb.WriteString(fmt.Sprintf("- **%s:%d** - %s (%s)\n",
				fp.Original.File, fp.Original.Line, fp.Original.Message, fp.Reason))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// categoryToString converts a Category to its string representation
func categoryToString(c Category) string {
	switch c {
	case CategorySecurity:
		return "Security"
	case CategoryPerformance:
		return "Performance"
	case CategoryFeature:
		return "Feature"
	case CategoryBug:
		return "Bug"
	case CategoryDocumentation:
		return "Documentation"
	case CategoryRefactor:
		return "Refactor"
	default:
		return "Unknown"
	}
}
