package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/security"
)

// BootstrapSpec defines the specification for executing a bootstrap tool
type BootstrapSpec struct {
	ComponentType string                 // Component type (e.g., "nextjs", "go-backend")
	TargetDir     string                 // Directory where project should be generated
	Config        map[string]interface{} // Component-specific configuration
	Flags         []string               // Additional CLI flags
	Timeout       time.Duration          // Execution timeout (0 for default)
}

// BaseExecutor provides common functionality for all bootstrap executors
type BaseExecutor struct {
	toolName       string
	defaultTimeout time.Duration
	toolValidator  *security.ToolValidator
}

// NewBaseExecutor creates a new base executor
func NewBaseExecutor(toolName string) *BaseExecutor {
	return &BaseExecutor{
		toolName:       toolName,
		defaultTimeout: 5 * time.Minute,
		toolValidator:  security.NewToolValidator(),
	}
}

// Execute runs a command and captures output
func (be *BaseExecutor) Execute(ctx context.Context, spec *BootstrapSpec) (*models.ExecutionResult, error) {
	// Validate tool is whitelisted
	if !be.toolValidator.IsToolWhitelisted(be.toolName) {
		return nil, fmt.Errorf("tool not whitelisted: %s", be.toolName)
	}

	// Build command
	cmdArgs, err := be.buildCommand(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to build command: %w", err)
	}

	// Validate command and flags
	if err := be.toolValidator.ValidateToolCommand(be.toolName, cmdArgs); err != nil {
		return nil, fmt.Errorf("tool command validation failed: %w", err)
	}

	// Set timeout
	timeout := spec.Timeout
	if timeout == 0 {
		timeout = be.defaultTimeout
	}

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Execute command
	startTime := time.Now()
	result, err := be.executeCommand(execCtx, be.toolName, cmdArgs, spec.TargetDir)
	result.Duration = time.Since(startTime)
	result.ToolUsed = be.toolName

	if err != nil {
		result.Success = false
		return result, err
	}

	// Sanitize output before returning
	result.Stdout = be.toolValidator.SanitizeCommandOutput(result.Stdout)
	result.Stderr = be.toolValidator.SanitizeCommandOutput(result.Stderr)

	result.Success = true
	return result, nil
}

// ExecuteWithStreaming runs a command and streams output to the provided writer
func (be *BaseExecutor) ExecuteWithStreaming(ctx context.Context, spec *BootstrapSpec, output io.Writer) (*models.ExecutionResult, error) {
	// Validate tool is whitelisted
	if !be.toolValidator.IsToolWhitelisted(be.toolName) {
		return nil, fmt.Errorf("tool not whitelisted: %s", be.toolName)
	}

	// Build command
	cmdArgs, err := be.buildCommand(spec)
	if err != nil {
		return nil, fmt.Errorf("failed to build command: %w", err)
	}

	// Validate command and flags
	if err := be.toolValidator.ValidateToolCommand(be.toolName, cmdArgs); err != nil {
		return nil, fmt.Errorf("tool command validation failed: %w", err)
	}

	// Set timeout
	timeout := spec.Timeout
	if timeout == 0 {
		timeout = be.defaultTimeout
	}

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Execute with streaming
	startTime := time.Now()
	result, err := be.executeWithStreaming(execCtx, be.toolName, cmdArgs, spec.TargetDir, output)
	result.Duration = time.Since(startTime)
	result.ToolUsed = be.toolName

	if err != nil {
		result.Success = false
		return result, err
	}

	result.Success = true
	return result, nil
}

// executeCommand runs a command and captures output
func (be *BaseExecutor) executeCommand(ctx context.Context, command string, args []string, workDir string) (*models.ExecutionResult, error) {
	result := &models.ExecutionResult{
		OutputDir: workDir,
	}

	cmd := exec.CommandContext(ctx, command, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}

	// Capture output
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// Split stdout and stderr (combined in this case)
	result.Stdout = outputStr
	result.Stderr = ""

	if err != nil {
		// Check if it's a context error
		if ctx.Err() == context.DeadlineExceeded {
			return result, fmt.Errorf("command timed out after %v", be.defaultTimeout)
		}
		if ctx.Err() == context.Canceled {
			return result, fmt.Errorf("command was canceled")
		}

		// Get exit code
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			result.ExitCode = exitErr.ExitCode()
		}

		return result, fmt.Errorf("command failed: %w\nOutput: %s", err, outputStr)
	}

	result.ExitCode = 0
	return result, nil
}

// executeWithStreaming runs a command and streams output
func (be *BaseExecutor) executeWithStreaming(ctx context.Context, command string, args []string, workDir string, output io.Writer) (*models.ExecutionResult, error) {
	result := &models.ExecutionResult{
		OutputDir: workDir,
	}

	cmd := exec.CommandContext(ctx, command, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}

	// Stream output
	cmd.Stdout = output
	cmd.Stderr = output

	err := cmd.Run()

	if err != nil {
		// Check if it's a context error
		if ctx.Err() == context.DeadlineExceeded {
			return result, fmt.Errorf("command timed out after %v", be.defaultTimeout)
		}
		if ctx.Err() == context.Canceled {
			return result, fmt.Errorf("command was canceled")
		}

		// Get exit code
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			result.ExitCode = exitErr.ExitCode()
		}

		return result, fmt.Errorf("command failed: %w", err)
	}

	result.ExitCode = 0
	return result, nil
}

// buildCommand builds the command arguments (to be overridden by specific executors)
func (be *BaseExecutor) buildCommand(spec *BootstrapSpec) ([]string, error) {
	// Validate flags before building command
	if err := be.toolValidator.ValidateToolFlags(spec.Flags); err != nil {
		return nil, fmt.Errorf("flag validation failed: %w", err)
	}

	// Default implementation - specific executors should override
	return spec.Flags, nil
}

// ValidateFlags validates command flags for security
func (be *BaseExecutor) ValidateFlags(flags []string) error {
	return be.toolValidator.ValidateToolFlags(flags)
}

// SanitizeOutput sanitizes command output before displaying
func (be *BaseExecutor) SanitizeOutput(output string) string {
	return be.toolValidator.SanitizeCommandOutput(output)
}

// SupportsComponent is a default implementation that returns false
// Specific executors should override this method
func (be *BaseExecutor) SupportsComponent(componentType string) bool {
	return false
}

// GetDefaultFlags is a default implementation that returns empty flags
// Specific executors should override this method
func (be *BaseExecutor) GetDefaultFlags(componentType string) []string {
	return []string{}
}

// ValidateConfig is a default implementation that does no validation
// Specific executors should override this method
func (be *BaseExecutor) ValidateConfig(config map[string]interface{}) error {
	return nil
}
