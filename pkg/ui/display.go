// Package ui provides display components for tables, trees, and progress tracking.
//
// This file contains implementations for visual display components that complement
// the interactive UI framework with formatted output capabilities.
package ui

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
)

// ShowTable displays a formatted table with optional pagination and sorting
func (ui *InteractiveUI) ShowTable(ctx context.Context, config interfaces.TableConfig) error {
	if len(config.Headers) == 0 {
		return fmt.Errorf("table must have at least one header")
	}

	ui.clearScreen()
	if config.Title != "" {
		ui.showHeader(config.Title, "")
	}

	// Calculate column widths
	colWidths := ui.calculateColumnWidths(config.Headers, config.Rows, config.MaxWidth)

	// Handle pagination
	pageSize := config.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	currentPage := 0
	totalPages := int(math.Ceil(float64(len(config.Rows)) / float64(pageSize)))

	// Handle sorting
	rows := make([][]string, len(config.Rows))
	copy(rows, config.Rows)
	sortColumn := -1
	sortAscending := true

	// Search functionality
	searchQuery := ""
	filteredRows := make([][]string, len(rows))
	copy(filteredRows, rows)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Apply search filter
		if searchQuery != "" {
			filteredRows = ui.filterTableRows(rows, searchQuery)
		} else {
			filteredRows = rows
		}

		// Apply sorting
		if config.Sortable && sortColumn >= 0 && sortColumn < len(config.Headers) {
			ui.sortTableRows(filteredRows, sortColumn, sortAscending)
		}

		// Calculate pagination for filtered results
		totalPages = int(math.Ceil(float64(len(filteredRows)) / float64(pageSize)))
		if currentPage >= totalPages && totalPages > 0 {
			currentPage = totalPages - 1
		}

		ui.displayTable(config, filteredRows, colWidths, currentPage, pageSize, totalPages, searchQuery, sortColumn, sortAscending)

		if !config.Pagination && !config.Searchable && !config.Sortable {
			break
		}

		ui.showTableHelp(config.Pagination, config.Searchable, config.Sortable)

		input, err := ui.readInput()
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		shouldExit := ui.handleTableInput(input, &currentPage, &totalPages, &searchQuery, &sortColumn, &sortAscending, config)
		if shouldExit {
			break
		}
	}

	return nil
}

// calculateColumnWidths calculates optimal column widths for table display
func (ui *InteractiveUI) calculateColumnWidths(headers []string, rows [][]string, maxWidth int) []int {
	if maxWidth <= 0 {
		maxWidth = 120 // Default terminal width
	}

	colCount := len(headers)
	colWidths := make([]int, colCount)

	// Initialize with header lengths
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	// Find maximum width needed for each column
	for _, row := range rows {
		for i, cell := range row {
			if i < colCount && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Adjust widths to fit within maxWidth
	totalWidth := 0
	for _, width := range colWidths {
		totalWidth += width + 3 // +3 for padding and separators
	}

	if totalWidth > maxWidth {
		// Proportionally reduce column widths
		ratio := float64(maxWidth-colCount*3) / float64(totalWidth-colCount*3)
		for i := range colWidths {
			colWidths[i] = int(float64(colWidths[i]) * ratio)
			if colWidths[i] < 5 { // Minimum column width
				colWidths[i] = 5
			}
		}
	}

	return colWidths
}

// displayTable renders the table with current settings
func (ui *InteractiveUI) displayTable(config interfaces.TableConfig, rows [][]string, colWidths []int, currentPage, pageSize, totalPages int, searchQuery string, sortColumn int, sortAscending bool) {
	ui.clearScreen()
	if config.Title != "" {
		ui.showHeader(config.Title, "")
	}

	// Show search query if active
	if searchQuery != "" {
		fmt.Printf("Search: %s (showing %d results)\n\n", ui.colorize(searchQuery, "cyan"), len(rows))
	}

	// Display headers
	ui.displayTableRow(config.Headers, colWidths, true, sortColumn, sortAscending)
	ui.displayTableSeparator(colWidths)

	// Calculate page bounds
	startIdx := currentPage * pageSize
	endIdx := startIdx + pageSize
	if endIdx > len(rows) {
		endIdx = len(rows)
	}

	// Display rows for current page
	for i := startIdx; i < endIdx; i++ {
		ui.displayTableRow(rows[i], colWidths, false, -1, true)
	}

	// Show pagination info
	if config.Pagination && totalPages > 1 {
		fmt.Printf("\nPage %d of %d (showing %d-%d of %d rows)\n",
			currentPage+1, totalPages, startIdx+1, endIdx, len(rows))
	}
}

// displayTableRow renders a single table row
func (ui *InteractiveUI) displayTableRow(row []string, colWidths []int, isHeader bool, sortColumn int, sortAscending bool) {
	fmt.Print("│")
	for i, cell := range row {
		if i >= len(colWidths) {
			break
		}

		// Truncate cell if too long
		if len(cell) > colWidths[i] {
			cell = cell[:colWidths[i]-3] + "..."
		}

		// Pad cell to column width
		padded := fmt.Sprintf(" %-*s ", colWidths[i], cell)

		if isHeader {
			padded = ui.colorize(padded, "bold")
			// Add sort indicator
			if i == sortColumn {
				indicator := "↑"
				if !sortAscending {
					indicator = "↓"
				}
				padded = strings.TrimRight(padded, " ") + ui.colorize(indicator, "cyan") + " "
			}
		}

		fmt.Print(padded)
		fmt.Print("│")
	}
	fmt.Println()
}

// displayTableSeparator renders the table separator line
func (ui *InteractiveUI) displayTableSeparator(colWidths []int) {
	fmt.Print("├")
	for i, width := range colWidths {
		fmt.Print(strings.Repeat("─", width+2))
		if i < len(colWidths)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤")
}

// filterTableRows filters table rows based on search query
func (ui *InteractiveUI) filterTableRows(rows [][]string, query string) [][]string {
	query = strings.ToLower(query)
	var filtered [][]string

	for _, row := range rows {
		for _, cell := range row {
			if strings.Contains(strings.ToLower(cell), query) {
				filtered = append(filtered, row)
				break
			}
		}
	}

	return filtered
}

// sortTableRows sorts table rows by the specified column
func (ui *InteractiveUI) sortTableRows(rows [][]string, column int, ascending bool) {
	sort.Slice(rows, func(i, j int) bool {
		if column >= len(rows[i]) || column >= len(rows[j]) {
			return false
		}

		a, b := rows[i][column], rows[j][column]

		// Try numeric comparison first
		if numA, errA := strconv.ParseFloat(a, 64); errA == nil {
			if numB, errB := strconv.ParseFloat(b, 64); errB == nil {
				if ascending {
					return numA < numB
				}
				return numA > numB
			}
		}

		// Fall back to string comparison
		if ascending {
			return strings.ToLower(a) < strings.ToLower(b)
		}
		return strings.ToLower(a) > strings.ToLower(b)
	})
}

// handleTableInput processes user input for table navigation
func (ui *InteractiveUI) handleTableInput(input string, currentPage *int, totalPages *int, searchQuery *string, sortColumn *int, sortAscending *bool, config interfaces.TableConfig) bool {
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "q", "quit", "exit", "":
		return true

	case "n", "next":
		if config.Pagination && *currentPage < *totalPages-1 {
			*currentPage++
		}

	case "p", "prev", "previous":
		if config.Pagination && *currentPage > 0 {
			*currentPage--
		}

	case "/":
		if config.Searchable {
			fmt.Print("Search: ")
			search, _ := ui.readInput()
			*searchQuery = strings.TrimSpace(search)
		}

	case "clear":
		*searchQuery = ""

	default:
		// Handle column sorting
		if config.Sortable && strings.HasPrefix(input, "sort ") {
			parts := strings.Fields(input)
			if len(parts) >= 2 {
				if col, err := strconv.Atoi(parts[1]); err == nil && col > 0 {
					newColumn := col - 1
					if newColumn == *sortColumn {
						*sortAscending = !*sortAscending
					} else {
						*sortColumn = newColumn
						*sortAscending = true
					}
				}
			}
		}

		// Handle page numbers
		if config.Pagination {
			if page, err := strconv.Atoi(input); err == nil && page > 0 && page <= *totalPages {
				*currentPage = page - 1
			}
		}
	}

	return false
}

// showTableHelp displays help for table navigation
func (ui *InteractiveUI) showTableHelp(pagination, searchable, sortable bool) {
	if !ui.config.ShowShortcuts {
		return
	}

	help := []string{"q: Quit"}

	if pagination {
		help = append(help, "n: Next page", "p: Previous page", "#: Go to page")
	}
	if searchable {
		help = append(help, "/: Search", "clear: Clear search")
	}
	if sortable {
		help = append(help, "sort #: Sort by column")
	}

	fmt.Printf("\n%s\n", ui.colorize(strings.Join(help, " | "), "gray"))
}

// ShowTree displays a tree structure with optional expansion
func (ui *InteractiveUI) ShowTree(ctx context.Context, config interfaces.TreeConfig) error {
	ui.clearScreen()
	if config.Title != "" {
		ui.showHeader(config.Title, "")
	}

	expandedNodes := make(map[string]bool)
	currentPath := ""

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ui.displayTree(config.Root, "", true, expandedNodes, config.ShowIcons, config.MaxDepth, 0)

		if !config.Expandable {
			break
		}

		ui.showTreeHelp()

		input, err := ui.readInput()
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		shouldExit := ui.handleTreeInput(input, &config.Root, expandedNodes, &currentPath)
		if shouldExit {
			break
		}
	}

	return nil
}

// displayTree renders the tree structure recursively
func (ui *InteractiveUI) displayTree(node interfaces.TreeNode, prefix string, isLast bool, expandedNodes map[string]bool, showIcons bool, maxDepth, currentDepth int) {
	if maxDepth > 0 && currentDepth >= maxDepth {
		return
	}

	// Determine tree characters
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	// Show icon if enabled
	icon := ""
	if showIcons && node.Icon != "" {
		icon = node.Icon + " "
	}

	// Show expansion indicator for nodes with children
	expansion := ""
	if len(node.Children) > 0 {
		if node.Expanded {
			expansion = ui.colorize("▼ ", "cyan")
		} else {
			expansion = ui.colorize("▶ ", "cyan")
		}
	}

	// Display the node
	label := node.Label
	if node.Selectable {
		label = ui.colorize(label, "white")
	}

	fmt.Printf("%s%s%s%s%s\n", prefix, connector, expansion, icon, label)

	// Display children if expanded
	if node.Expanded && len(node.Children) > 0 {
		childPrefix := prefix
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}

		for i, child := range node.Children {
			isChildLast := i == len(node.Children)-1
			ui.displayTree(child, childPrefix, isChildLast, expandedNodes, showIcons, maxDepth, currentDepth+1)
		}
	}
}

// handleTreeInput processes user input for tree navigation
func (ui *InteractiveUI) handleTreeInput(input string, root *interfaces.TreeNode, expandedNodes map[string]bool, currentPath *string) bool {
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "q", "quit", "exit", "":
		return true

	case "expand", "e":
		ui.expandAllNodes(root)

	case "collapse", "c":
		ui.collapseAllNodes(root)

	case "toggle", "t":
		// Toggle current node (would need path tracking for full implementation)
		root.Expanded = !root.Expanded
	}

	return false
}

// expandAllNodes expands all nodes in the tree
func (ui *InteractiveUI) expandAllNodes(node *interfaces.TreeNode) {
	node.Expanded = true
	for i := range node.Children {
		ui.expandAllNodes(&node.Children[i])
	}
}

// collapseAllNodes collapses all nodes in the tree
func (ui *InteractiveUI) collapseAllNodes(node *interfaces.TreeNode) {
	node.Expanded = false
	for i := range node.Children {
		ui.collapseAllNodes(&node.Children[i])
	}
}

// showTreeHelp displays help for tree navigation
func (ui *InteractiveUI) showTreeHelp() {
	if !ui.config.ShowShortcuts {
		return
	}

	help := []string{
		"q: Quit",
		"e: Expand all",
		"c: Collapse all",
		"t: Toggle current",
	}

	fmt.Printf("\n%s\n", ui.colorize(strings.Join(help, " | "), "gray"))
}

// ShowProgress creates and returns a progress tracker
func (ui *InteractiveUI) ShowProgress(ctx context.Context, config interfaces.ProgressConfig) (interfaces.ProgressTracker, error) {
	tracker := &ProgressTracker{
		ui:          ui,
		config:      config,
		ctx:         ctx,
		startTime:   time.Now(),
		progress:    0.0,
		currentStep: 0,
		logs:        []string{},
		mutex:       &sync.RWMutex{},
		done:        make(chan struct{}),
	}

	// Start the display goroutine
	go tracker.displayLoop()

	return tracker, nil
}

// ProgressTracker implements the ProgressTracker interface
type ProgressTracker struct {
	ui          *InteractiveUI
	config      interfaces.ProgressConfig
	ctx         context.Context
	startTime   time.Time
	progress    float64
	currentStep int
	stepDesc    string
	logs        []string
	completed   bool
	failed      bool
	error       error
	cancelled   bool
	mutex       *sync.RWMutex
	done        chan struct{}
}

// SetProgress updates the progress value (0.0 to 1.0)
func (pt *ProgressTracker) SetProgress(progress float64) error {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	if progress < 0.0 {
		progress = 0.0
	} else if progress > 1.0 {
		progress = 1.0
	}

	pt.progress = progress
	return nil
}

// SetCurrentStep sets the current step
func (pt *ProgressTracker) SetCurrentStep(step int, description string) error {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	pt.currentStep = step
	pt.stepDesc = description
	return nil
}

// AddLog adds a log message
func (pt *ProgressTracker) AddLog(message string) error {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	timestamp := time.Now().Format("15:04:05")
	pt.logs = append(pt.logs, fmt.Sprintf("[%s] %s", timestamp, message))

	// Keep only last 10 log messages
	if len(pt.logs) > 10 {
		pt.logs = pt.logs[len(pt.logs)-10:]
	}

	return nil
}

// Complete marks the progress as complete
func (pt *ProgressTracker) Complete() error {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	pt.completed = true
	pt.progress = 1.0
	close(pt.done)
	return nil
}

// Fail marks the progress as failed
func (pt *ProgressTracker) Fail(err error) error {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	pt.failed = true
	pt.error = err
	close(pt.done)
	return nil
}

// IsCancelled checks if the operation was cancelled
func (pt *ProgressTracker) IsCancelled() bool {
	select {
	case <-pt.ctx.Done():
		pt.mutex.Lock()
		pt.cancelled = true
		pt.mutex.Unlock()
		return true
	default:
		pt.mutex.RLock()
		cancelled := pt.cancelled
		pt.mutex.RUnlock()
		return cancelled
	}
}

// Close closes the progress tracker
func (pt *ProgressTracker) Close() error {
	pt.mutex.Lock()
	defer pt.mutex.Unlock()

	if !pt.completed && !pt.failed && !pt.cancelled {
		pt.cancelled = true
		close(pt.done)
	}

	return nil
}

// displayLoop runs the progress display loop
func (pt *ProgressTracker) displayLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-pt.done:
			pt.displayFinal()
			return
		case <-pt.ctx.Done():
			pt.mutex.Lock()
			pt.cancelled = true
			pt.mutex.Unlock()
			pt.displayFinal()
			return
		case <-ticker.C:
			pt.displayCurrent()
		}
	}
}

// displayCurrent displays the current progress state
func (pt *ProgressTracker) displayCurrent() {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()

	pt.ui.clearScreen()
	pt.ui.showHeader(pt.config.Title, pt.config.Description)

	// Show progress bar
	if pt.config.ShowPercent {
		percentage := int(pt.progress * 100)
		barWidth := 50
		filledWidth := int(float64(barWidth) * pt.progress)

		bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)
		fmt.Printf("Progress: [%s] %d%%\n", pt.ui.colorize(bar, "green"), percentage)
	}

	// Show current step
	if len(pt.config.Steps) > 0 && pt.currentStep < len(pt.config.Steps) {
		fmt.Printf("Step %d/%d: %s\n", pt.currentStep+1, len(pt.config.Steps), pt.config.Steps[pt.currentStep])
	} else if pt.stepDesc != "" {
		fmt.Printf("Current: %s\n", pt.stepDesc)
	}

	// Show ETA
	if pt.config.ShowETA && pt.progress > 0 {
		elapsed := time.Since(pt.startTime)
		estimated := time.Duration(float64(elapsed) / pt.progress)
		remaining := estimated - elapsed
		if remaining > 0 {
			fmt.Printf("ETA: %s\n", remaining.Round(time.Second))
		}
	}

	// Show recent logs
	if len(pt.logs) > 0 {
		fmt.Println()
		fmt.Println("Recent activity:")
		for _, log := range pt.logs {
			fmt.Printf("  %s\n", pt.ui.colorize(log, "gray"))
		}
	}

	// Show cancellation option
	if pt.config.Cancellable {
		fmt.Printf("\n%s\n", pt.ui.colorize("Press Ctrl+C to cancel", "gray"))
	}
}

// displayFinal displays the final progress state
func (pt *ProgressTracker) displayFinal() {
	pt.mutex.RLock()
	defer pt.mutex.RUnlock()

	pt.ui.clearScreen()
	pt.ui.showHeader(pt.config.Title, pt.config.Description)

	duration := time.Since(pt.startTime)

	if pt.completed {
		fmt.Printf("%s Completed successfully in %s\n", pt.ui.colorize("✓", "green"), duration.Round(time.Second))
	} else if pt.failed {
		fmt.Printf("%s Failed after %s\n", pt.ui.colorize("✗", "red"), duration.Round(time.Second))
		if pt.error != nil {
			fmt.Printf("Error: %s\n", pt.error.Error())
		}
	} else if pt.cancelled {
		fmt.Printf("%s Cancelled after %s\n", pt.ui.colorize("⚠", "yellow"), duration.Round(time.Second))
	}

	fmt.Print("\nPress Enter to continue...")
	if _, err := pt.ui.readInput(); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("\nError reading input: %v\n", err)
	}
}
