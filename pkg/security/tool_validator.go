// Package security provides tool execution validation for the CLI tool.
package security

import (
	"fmt"
	"regexp"
	"strings"
)

// ToolValidator validates tool commands and flags before execution
type ToolValidator struct {
	whitelistedTools map[string]*ToolWhitelist
	validator        *Validator
}

// ToolWhitelist defines allowed commands and flags for a tool
type ToolWhitelist struct {
	Command       string
	AllowedFlags  []string
	FlagPatterns  []*regexp.Regexp
	MaxFlagCount  int
	RequiresFlags bool
}

// NewToolValidator creates a new tool validator with default whitelists
func NewToolValidator() *ToolValidator {
	tv := &ToolValidator{
		whitelistedTools: make(map[string]*ToolWhitelist),
		validator:        NewValidator(),
	}

	// Register default tool whitelists
	tv.registerDefaultWhitelists()

	return tv
}

// registerDefaultWhitelists registers the default tool whitelists
func (tv *ToolValidator) registerDefaultWhitelists() {
	// NPX (for Next.js and other npm tools)
	tv.whitelistedTools["npx"] = &ToolWhitelist{
		Command: "npx",
		AllowedFlags: []string{
			"create-next-app",
			"create-react-app",
			"create-vite",
			"--typescript",
			"--javascript",
			"--tailwind",
			"--no-tailwind",
			"--app",
			"--no-app",
			"--no-git",
			"--use-npm",
			"--use-yarn",
			"--use-pnpm",
			"--eslint",
			"--no-eslint",
			"--src-dir",
			"--no-src-dir",
			"--import-alias",
		},
		FlagPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^--[a-z][a-z0-9-]*$`),
			regexp.MustCompile(`^-[a-z]$`),
			regexp.MustCompile(`^[a-z][a-z0-9-]*@[0-9.]+$`), // package@version
			regexp.MustCompile(`^[a-z][a-z0-9-]*@latest$`),  // package@latest
			regexp.MustCompile(`^[a-z][a-z0-9_-]+$`),        // project names (alphanumeric, dash, underscore)
			regexp.MustCompile(`^@[a-z0-9-]+/[a-z0-9-]+$`),  // scoped packages like @types/node
		},
		MaxFlagCount: 20,
	}

	// Go
	tv.whitelistedTools["go"] = &ToolWhitelist{
		Command: "go",
		AllowedFlags: []string{
			"mod",
			"init",
			"get",
			"install",
			"build",
			"run",
			"test",
			"tidy",
			"vendor",
			"-v",
			"-u",
			"-d",
		},
		FlagPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^-[a-z]+$`),
			regexp.MustCompile(`^[a-z][a-z0-9/._-]+$`),                  // module path (github.com/user/repo)
			regexp.MustCompile(`^[a-z][a-z0-9/._-]+@[a-z0-9.]+$`),       // module@version
			regexp.MustCompile(`^[a-z][a-z0-9/._-]+@latest$`),           // module@latest
			regexp.MustCompile(`^github\.com/[a-z0-9_-]+/[a-z0-9_-]+$`), // github module paths
		},
		MaxFlagCount: 15,
	}

	// Gradle
	tv.whitelistedTools["gradle"] = &ToolWhitelist{
		Command: "gradle",
		AllowedFlags: []string{
			"init",
			"build",
			"clean",
			"assemble",
			"test",
			"--type",
			"--dsl",
			"--test-framework",
			"--project-name",
			"--package",
			"-q",
			"-i",
			"-d",
		},
		FlagPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^--[a-z][a-z0-9-]*$`),
			regexp.MustCompile(`^-[a-z]$`),
			regexp.MustCompile(`^[a-z][a-z0-9._-]+$`),                        // project names and package names
			regexp.MustCompile(`^com\.[a-z][a-z0-9._-]*\.[a-z][a-z0-9_-]*$`), // Java package names
		},
		MaxFlagCount: 15,
	}

	// Xcodebuild
	tv.whitelistedTools["xcodebuild"] = &ToolWhitelist{
		Command: "xcodebuild",
		AllowedFlags: []string{
			"-project",
			"-scheme",
			"-configuration",
			"-sdk",
			"-destination",
			"build",
			"clean",
			"test",
		},
		FlagPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^-[a-z]+$`),
		},
		MaxFlagCount: 10,
	}

	// Swift
	tv.whitelistedTools["swift"] = &ToolWhitelist{
		Command: "swift",
		AllowedFlags: []string{
			"package",
			"init",
			"build",
			"test",
			"run",
			"--type",
			"--name",
		},
		FlagPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^--[a-z][a-z0-9-]*$`),
			regexp.MustCompile(`^-[a-z]$`),
			regexp.MustCompile(`^[A-Z][a-zA-Z0-9]+$`), // Swift project names (PascalCase)
			regexp.MustCompile(`^[a-z][a-z0-9_-]+$`),  // lowercase project names
		},
		MaxFlagCount: 10,
	}

	// Docker
	tv.whitelistedTools["docker"] = &ToolWhitelist{
		Command: "docker",
		AllowedFlags: []string{
			"build",
			"run",
			"compose",
			"up",
			"down",
			"ps",
			"logs",
			"-d",
			"-f",
			"-t",
			"--rm",
			"--name",
			"--build-arg",
		},
		FlagPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^--[a-z][a-z0-9-]*$`),
			regexp.MustCompile(`^-[a-z]+$`),
		},
		MaxFlagCount: 20,
	}

	// Terraform
	tv.whitelistedTools["terraform"] = &ToolWhitelist{
		Command: "terraform",
		AllowedFlags: []string{
			"init",
			"plan",
			"apply",
			"destroy",
			"validate",
			"fmt",
			"-auto-approve",
			"-var",
			"-var-file",
		},
		FlagPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^-[a-z][a-z0-9-]*$`),
		},
		MaxFlagCount: 15,
	}
}

// ValidateToolCommand validates a tool command before execution
func (tv *ToolValidator) ValidateToolCommand(toolName string, args []string) error {
	// Check if tool is whitelisted
	whitelist, exists := tv.whitelistedTools[toolName]
	if !exists {
		return fmt.Errorf("tool '%s' is not whitelisted for execution", toolName)
	}

	// Validate command matches whitelist
	if whitelist.Command != toolName {
		return fmt.Errorf("tool command mismatch: expected '%s', got '%s'", whitelist.Command, toolName)
	}

	// Check flag count
	if len(args) > whitelist.MaxFlagCount {
		return fmt.Errorf("too many arguments: maximum %d allowed, got %d", whitelist.MaxFlagCount, len(args))
	}

	// Validate each argument/flag
	for _, arg := range args {
		if err := tv.validateArgument(arg, whitelist); err != nil {
			return fmt.Errorf("invalid argument '%s': %w", arg, err)
		}
	}

	// Use validator to check for command injection
	if err := tv.validator.ValidateToolFlags(args); err != nil {
		return err
	}

	return nil
}

// validateArgument validates a single argument against the whitelist
func (tv *ToolValidator) validateArgument(arg string, whitelist *ToolWhitelist) error {
	// Check if argument is in allowed list
	for _, allowed := range whitelist.AllowedFlags {
		if arg == allowed {
			return nil
		}
	}

	// Check if argument matches any allowed pattern
	for _, pattern := range whitelist.FlagPatterns {
		if pattern.MatchString(arg) {
			return nil
		}
	}

	return fmt.Errorf("argument not in whitelist")
}

// ValidateToolFlags validates tool flags for security issues
func (tv *ToolValidator) ValidateToolFlags(flags []string) error {
	return tv.validator.ValidateToolFlags(flags)
}

// IsToolWhitelisted checks if a tool is whitelisted
func (tv *ToolValidator) IsToolWhitelisted(toolName string) bool {
	_, exists := tv.whitelistedTools[toolName]
	return exists
}

// GetWhitelistedTools returns a list of all whitelisted tools
func (tv *ToolValidator) GetWhitelistedTools() []string {
	tools := make([]string, 0, len(tv.whitelistedTools))
	for tool := range tv.whitelistedTools {
		tools = append(tools, tool)
	}
	return tools
}

// AddToolWhitelist adds a custom tool whitelist
func (tv *ToolValidator) AddToolWhitelist(toolName string, whitelist *ToolWhitelist) error {
	if toolName == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	if whitelist == nil {
		return fmt.Errorf("whitelist cannot be nil")
	}

	if whitelist.Command == "" {
		return fmt.Errorf("whitelist command cannot be empty")
	}

	tv.whitelistedTools[toolName] = whitelist
	return nil
}

// ValidateCommandString validates a complete command string
func (tv *ToolValidator) ValidateCommandString(commandStr string) error {
	// Check for dangerous characters
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\n", "\r", "\\"}
	for _, char := range dangerousChars {
		if strings.Contains(commandStr, char) {
			return fmt.Errorf("command contains dangerous character '%s'", char)
		}
	}

	// Check for command chaining attempts
	if strings.Contains(commandStr, "&&") || strings.Contains(commandStr, "||") {
		return fmt.Errorf("command chaining is not allowed")
	}

	// Check for redirection attempts
	if strings.Contains(commandStr, ">") || strings.Contains(commandStr, "<") {
		return fmt.Errorf("command redirection is not allowed")
	}

	// Check for pipe attempts
	if strings.Contains(commandStr, "|") {
		return fmt.Errorf("command piping is not allowed")
	}

	return nil
}

// SanitizeCommandOutput sanitizes command output before displaying to user
func (tv *ToolValidator) SanitizeCommandOutput(output string) string {
	// Remove potential sensitive information
	sanitized := output

	// Remove absolute paths that might contain usernames
	homePattern := regexp.MustCompile(`/home/[^/\s]+`)
	sanitized = homePattern.ReplaceAllString(sanitized, "/home/[user]")

	usersPattern := regexp.MustCompile(`/Users/[^/\s]+`)
	sanitized = usersPattern.ReplaceAllString(sanitized, "/Users/[user]")

	// Remove potential API keys or tokens in output
	apiKeyPattern := regexp.MustCompile(`[a-zA-Z0-9_-]{32,}`)
	sanitized = apiKeyPattern.ReplaceAllStringFunc(sanitized, func(match string) string {
		// Only redact if it looks like a key (all caps or mixed case with numbers)
		if regexp.MustCompile(`[A-Z0-9_-]{32,}`).MatchString(match) {
			return "[REDACTED]"
		}
		return match
	})

	return sanitized
}

// ValidateToolPath validates that a tool path is safe
func (tv *ToolValidator) ValidateToolPath(toolPath string) error {
	// Use validator to check path security
	if err := tv.validator.ValidatePathSecurity(toolPath); err != nil {
		return err
	}

	// Additional checks for tool paths
	if strings.Contains(toolPath, "..") {
		return fmt.Errorf("tool path contains path traversal")
	}

	// Check for suspicious paths
	suspiciousPaths := []string{
		"/tmp/",
		"/var/tmp/",
		"/dev/",
		"/proc/",
		"/sys/",
	}

	for _, suspicious := range suspiciousPaths {
		if strings.HasPrefix(toolPath, suspicious) {
			return fmt.Errorf("tool path in suspicious location: %s", suspicious)
		}
	}

	return nil
}
