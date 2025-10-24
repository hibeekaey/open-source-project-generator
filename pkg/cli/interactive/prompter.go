package interactive

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// PrompterInterface defines the interface for user input prompts
type PrompterInterface interface {
	// Input prompts for a single line of text
	Input(message string, defaultValue string) (string, error)

	// Select prompts for selection from a list
	Select(message string, options []string) (string, error)

	// MultiSelect prompts for multiple selections
	MultiSelect(message string, options []string) ([]string, error)

	// Confirm prompts for yes/no confirmation
	Confirm(message string, defaultValue bool) (bool, error)

	// Password prompts for password input (hidden)
	Password(message string) (string, error)
}

// CLIPrompter implements Prompter using standard input/output
type CLIPrompter struct {
	reader *bufio.Reader
	writer io.Writer
}

// NewCLIPrompter creates a new CLI prompter
func NewCLIPrompter() *CLIPrompter {
	return &CLIPrompter{
		reader: bufio.NewReader(os.Stdin),
		writer: os.Stdout,
	}
}

// Input prompts for a single line of text
func (p *CLIPrompter) Input(message string, defaultValue string) (string, error) {
	if defaultValue != "" {
		_, _ = fmt.Fprintf(p.writer, "%s [%s]: ", message, defaultValue) // Error intentionally ignored for interactive output
	} else {
		_, _ = fmt.Fprintf(p.writer, "%s: ", message) // Error intentionally ignored for interactive output
	}

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)
	if input == "" && defaultValue != "" {
		return defaultValue, nil
	}

	return input, nil
}

// Select prompts for selection from a list
func (p *CLIPrompter) Select(message string, options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options provided")
	}

	_, _ = fmt.Fprintln(p.writer, message) // Error intentionally ignored for interactive output
	for i, option := range options {
		_, _ = fmt.Fprintf(p.writer, "  %d) %s\n", i+1, option) // Error intentionally ignored for interactive output
	}
	_, _ = fmt.Fprintf(p.writer, "Select (1-%d): ", len(options)) // Error intentionally ignored for interactive output

	for {
		input, err := p.reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			_, _ = fmt.Fprintf(p.writer, "Please select an option (1-%d): ", len(options)) // Error intentionally ignored for interactive output
			continue
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(options) {
			_, _ = fmt.Fprintf(p.writer, "Invalid selection. Please enter a number between 1 and %d: ", len(options)) // Error intentionally ignored for interactive output
			continue
		}

		return options[choice-1], nil
	}
}

// MultiSelect prompts for multiple selections
func (p *CLIPrompter) MultiSelect(message string, options []string) ([]string, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("no options provided")
	}

	_, _ = fmt.Fprintln(p.writer, message) // Error intentionally ignored for interactive output
	for i, option := range options {
		_, _ = fmt.Fprintf(p.writer, "  %d) %s\n", i+1, option) // Error intentionally ignored for interactive output
	}
	_, _ = fmt.Fprintf(p.writer, "Select multiple (comma-separated, e.g., 1,3,4) or 'all': ") // Error intentionally ignored for interactive output

	for {
		input, err := p.reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			_, _ = fmt.Fprintf(p.writer, "Please select at least one option: ") // Error intentionally ignored for interactive output
			continue
		}

		// Handle "all" selection
		if strings.ToLower(input) == "all" {
			return options, nil
		}

		// Parse comma-separated selections
		parts := strings.Split(input, ",")
		selected := make([]string, 0, len(parts))
		valid := true

		for _, part := range parts {
			part = strings.TrimSpace(part)
			choice, err := strconv.Atoi(part)
			if err != nil || choice < 1 || choice > len(options) {
				_, _ = fmt.Fprintf(p.writer, "Invalid selection '%s'. Please enter numbers between 1 and %d (comma-separated): ", part, len(options)) // Error intentionally ignored for interactive output
				valid = false
				break
			}
			selected = append(selected, options[choice-1])
		}

		if valid && len(selected) > 0 {
			return selected, nil
		}
	}
}

// Confirm prompts for yes/no confirmation
func (p *CLIPrompter) Confirm(message string, defaultValue bool) (bool, error) {
	defaultStr := "y/N"
	if defaultValue {
		defaultStr = "Y/n"
	}

	_, _ = fmt.Fprintf(p.writer, "%s [%s]: ", message, defaultStr) // Error intentionally ignored for interactive output

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return defaultValue, nil
	}

	switch input {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		_, _ = fmt.Fprintf(p.writer, "Please answer 'y' or 'n': ") // Error intentionally ignored for interactive output
		return p.Confirm(message, defaultValue)
	}
}

// Password prompts for password input (hidden)
func (p *CLIPrompter) Password(message string) (string, error) {
	_, _ = fmt.Fprintf(p.writer, "%s: ", message) // Error intentionally ignored for interactive output

	// Check if stdin is a terminal
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		// Not a terminal, read normally
		input, err := p.reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		return strings.TrimSpace(input), nil
	}

	// Read password with hidden input
	password, err := term.ReadPassword(fd)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	_, _ = fmt.Fprintln(p.writer) // Error intentionally ignored for interactive output (newline after password)
	return string(password), nil
}
