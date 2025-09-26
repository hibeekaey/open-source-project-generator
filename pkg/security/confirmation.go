package security

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// ConfirmationManager handles user confirmation prompts for dangerous operations
type ConfirmationManager struct {
	nonInteractive bool
	defaultAnswer  bool
	timeout        time.Duration
	reader         *bufio.Reader
}

// NewConfirmationManager creates a new confirmation manager
func NewConfirmationManager() *ConfirmationManager {
	return &ConfirmationManager{
		nonInteractive: false,
		defaultAnswer:  false,
		timeout:        30 * time.Second,
		reader:         bufio.NewReader(os.Stdin),
	}
}

// ConfirmationRequest represents a confirmation request
type ConfirmationRequest struct {
	Message         string                 `json:"message"`
	Details         []string               `json:"details,omitempty"`
	Impact          string                 `json:"impact"` // "safe", "warning", "destructive"
	DefaultAnswer   bool                   `json:"default_answer"`
	RequireExplicit bool                   `json:"require_explicit"` // Require explicit "yes" for dangerous operations
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// ConfirmationResult represents the result of a confirmation request
type ConfirmationResult struct {
	Confirmed      bool      `json:"confirmed"`
	UserInput      string    `json:"user_input"`
	Timestamp      time.Time `json:"timestamp"`
	NonInteractive bool      `json:"non_interactive"`
	TimedOut       bool      `json:"timed_out"`
	DefaultUsed    bool      `json:"default_used"`
}

// SetNonInteractive sets whether to run in non-interactive mode
func (cm *ConfirmationManager) SetNonInteractive(nonInteractive bool) {
	cm.nonInteractive = nonInteractive
}

// SetDefaultAnswer sets the default answer for non-interactive mode
func (cm *ConfirmationManager) SetDefaultAnswer(defaultAnswer bool) {
	cm.defaultAnswer = defaultAnswer
}

// SetTimeout sets the timeout for confirmation prompts
func (cm *ConfirmationManager) SetTimeout(timeout time.Duration) {
	cm.timeout = timeout
}

// IsNonInteractive returns whether running in non-interactive mode
func (cm *ConfirmationManager) IsNonInteractive() bool {
	return cm.nonInteractive
}

// Confirm prompts the user for confirmation
func (cm *ConfirmationManager) Confirm(request *ConfirmationRequest) (*ConfirmationResult, error) {
	result := &ConfirmationResult{
		Timestamp:      time.Now(),
		NonInteractive: cm.nonInteractive,
	}

	// In non-interactive mode, use default answer
	if cm.nonInteractive {
		result.Confirmed = cm.defaultAnswer
		result.DefaultUsed = true
		return result, nil
	}

	// Display confirmation prompt
	fmt.Println()
	fmt.Printf("⚠️  CONFIRMATION REQUIRED ⚠️\n")
	fmt.Printf("Impact Level: %s\n", strings.ToUpper(request.Impact))
	fmt.Printf("Message: %s\n", request.Message)

	// Display details if provided
	if len(request.Details) > 0 {
		fmt.Println("\nDetails:")
		for _, detail := range request.Details {
			fmt.Printf("  • %s\n", detail)
		}
	}

	// Display metadata if provided
	if len(request.Metadata) > 0 {
		fmt.Println("\nAdditional Information:")
		for key, value := range request.Metadata {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	fmt.Println()

	// Determine prompt text based on impact and requirements
	var promptText string
	var validAnswers []string

	if request.RequireExplicit || request.Impact == "destructive" {
		promptText = "Type 'yes' to confirm, 'no' to cancel"
		validAnswers = []string{"yes", "no"}
	} else {
		defaultText := "n"
		if request.DefaultAnswer {
			defaultText = "y"
		}
		promptText = fmt.Sprintf("Continue? [y/N] (default: %s)", defaultText)
		validAnswers = []string{"y", "yes", "n", "no", ""}
	}

	// Get user input with timeout
	userInput, timedOut, err := cm.getUserInputWithTimeout(promptText)
	if err != nil {
		return result, fmt.Errorf("failed to get user input: %w", err)
	}

	result.UserInput = userInput
	result.TimedOut = timedOut

	// Handle timeout
	if timedOut {
		fmt.Printf("\nTimeout reached. Using default answer: %t\n", request.DefaultAnswer)
		result.Confirmed = request.DefaultAnswer
		result.DefaultUsed = true
		return result, nil
	}

	// Parse user input
	userInput = strings.ToLower(strings.TrimSpace(userInput))

	// Validate input
	validInput := false
	for _, valid := range validAnswers {
		if userInput == valid {
			validInput = true
			break
		}
	}

	if !validInput {
		return result, fmt.Errorf("invalid input: %s (expected: %s)", userInput, strings.Join(validAnswers, ", "))
	}

	// Determine confirmation result
	if request.RequireExplicit || request.Impact == "destructive" {
		result.Confirmed = userInput == "yes"
	} else {
		switch userInput {
		case "y", "yes":
			result.Confirmed = true
		case "n", "no":
			result.Confirmed = false
		case "":
			result.Confirmed = request.DefaultAnswer
			result.DefaultUsed = true
		}
	}

	return result, nil
}

// ConfirmFileOverwrite prompts for confirmation before overwriting a file
func (cm *ConfirmationManager) ConfirmFileOverwrite(filePath string, fileSize int64) (*ConfirmationResult, error) {
	request := &ConfirmationRequest{
		Message:         fmt.Sprintf("File '%s' already exists and will be overwritten", filePath),
		Impact:          "destructive",
		DefaultAnswer:   false,
		RequireExplicit: true,
		Metadata: map[string]interface{}{
			"file_path": filePath,
			"file_size": fileSize,
		},
	}

	if fileSize > 0 {
		request.Details = []string{
			fmt.Sprintf("File size: %d bytes", fileSize),
			"This operation cannot be undone without a backup",
		}
	}

	return cm.Confirm(request)
}

// ConfirmDirectoryDelete prompts for confirmation before deleting a directory
func (cm *ConfirmationManager) ConfirmDirectoryDelete(dirPath string, fileCount int, totalSize int64) (*ConfirmationResult, error) {
	request := &ConfirmationRequest{
		Message:         fmt.Sprintf("Directory '%s' and all its contents will be permanently deleted", dirPath),
		Impact:          "destructive",
		DefaultAnswer:   false,
		RequireExplicit: true,
		Details: []string{
			fmt.Sprintf("Files to be deleted: %d", fileCount),
			fmt.Sprintf("Total size: %d bytes", totalSize),
			"This operation cannot be undone",
		},
		Metadata: map[string]interface{}{
			"directory_path": dirPath,
			"file_count":     fileCount,
			"total_size":     totalSize,
		},
	}

	return cm.Confirm(request)
}

// ConfirmBulkOperation prompts for confirmation before bulk operations
func (cm *ConfirmationManager) ConfirmBulkOperation(operationType string, itemCount int, details []string) (*ConfirmationResult, error) {
	impact := "warning"
	if itemCount > 100 || strings.Contains(strings.ToLower(operationType), "delete") {
		impact = "destructive"
	}

	request := &ConfirmationRequest{
		Message:         fmt.Sprintf("About to perform %s on %d items", operationType, itemCount),
		Impact:          impact,
		DefaultAnswer:   false,
		RequireExplicit: impact == "destructive",
		Details:         details,
		Metadata: map[string]interface{}{
			"operation_type": operationType,
			"item_count":     itemCount,
		},
	}

	return cm.Confirm(request)
}

// ConfirmSecurityRisk prompts for confirmation when security risks are detected
func (cm *ConfirmationManager) ConfirmSecurityRisk(riskDescription string, riskLevel string, details []string) (*ConfirmationResult, error) {
	request := &ConfirmationRequest{
		Message:         fmt.Sprintf("Security risk detected: %s", riskDescription),
		Impact:          "destructive",
		DefaultAnswer:   false,
		RequireExplicit: true,
		Details:         append([]string{fmt.Sprintf("Risk Level: %s", riskLevel)}, details...),
		Metadata: map[string]interface{}{
			"risk_description": riskDescription,
			"risk_level":       riskLevel,
		},
	}

	return cm.Confirm(request)
}

// ConfirmWithDryRun prompts for confirmation after showing dry-run results
func (cm *ConfirmationManager) ConfirmWithDryRun(dryRunSummary map[string]interface{}) (*ConfirmationResult, error) {
	destructiveOps := 0
	warningOps := 0
	if val, ok := dryRunSummary["destructive_operations"].(int); ok {
		destructiveOps = val
	}
	if val, ok := dryRunSummary["warning_operations"].(int); ok {
		warningOps = val
	}

	impact := "safe"
	if destructiveOps > 0 {
		impact = "destructive"
	} else if warningOps > 0 {
		impact = "warning"
	}

	details := []string{
		fmt.Sprintf("Total operations: %v", dryRunSummary["total_operations"]),
		fmt.Sprintf("Safe operations: %v", dryRunSummary["safe_operations"]),
		fmt.Sprintf("Warning operations: %d", warningOps),
		fmt.Sprintf("Destructive operations: %d", destructiveOps),
		fmt.Sprintf("Files affected: %v", dryRunSummary["total_files_affected"]),
	}

	request := &ConfirmationRequest{
		Message:         "Proceed with the operations shown above?",
		Impact:          impact,
		DefaultAnswer:   impact == "safe",
		RequireExplicit: destructiveOps > 0,
		Details:         details,
		Metadata:        dryRunSummary,
	}

	return cm.Confirm(request)
}

// getUserInputWithTimeout gets user input with a timeout
func (cm *ConfirmationManager) getUserInputWithTimeout(prompt string) (string, bool, error) {
	fmt.Printf("%s: ", prompt)

	// Create a channel to receive the input
	inputChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	// Start a goroutine to read input
	go func() {
		input, err := cm.reader.ReadString('\n')
		if err != nil {
			errorChan <- err
			return
		}
		inputChan <- strings.TrimSpace(input)
	}()

	// Wait for input or timeout
	select {
	case input := <-inputChan:
		return input, false, nil
	case err := <-errorChan:
		return "", false, err
	case <-time.After(cm.timeout):
		return "", true, nil
	}
}

// ConfirmationHistory tracks confirmation history for auditing
type ConfirmationHistory struct {
	entries []ConfirmationHistoryEntry
}

// ConfirmationHistoryEntry represents a single confirmation in history
type ConfirmationHistoryEntry struct {
	Request   *ConfirmationRequest `json:"request"`
	Result    *ConfirmationResult  `json:"result"`
	Timestamp time.Time            `json:"timestamp"`
}

// NewConfirmationHistory creates a new confirmation history tracker
func NewConfirmationHistory() *ConfirmationHistory {
	return &ConfirmationHistory{
		entries: make([]ConfirmationHistoryEntry, 0),
	}
}

// Record records a confirmation in the history
func (ch *ConfirmationHistory) Record(request *ConfirmationRequest, result *ConfirmationResult) {
	entry := ConfirmationHistoryEntry{
		Request:   request,
		Result:    result,
		Timestamp: time.Now(),
	}
	ch.entries = append(ch.entries, entry)
}

// GetHistory returns all confirmation history entries
func (ch *ConfirmationHistory) GetHistory() []ConfirmationHistoryEntry {
	return ch.entries
}

// GetRecentHistory returns recent confirmation history entries
func (ch *ConfirmationHistory) GetRecentHistory(limit int) []ConfirmationHistoryEntry {
	if limit <= 0 || limit >= len(ch.entries) {
		return ch.entries
	}
	return ch.entries[len(ch.entries)-limit:]
}

// Clear clears the confirmation history
func (ch *ConfirmationHistory) Clear() {
	ch.entries = make([]ConfirmationHistoryEntry, 0)
}
