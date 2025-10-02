package formats

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/interfaces"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// DockerValidator provides specialized Dockerfile validation
type DockerValidator struct {
	securityRules []DockerSecurityRule
}

// DockerSecurityRule defines security rules for Dockerfile validation
type DockerSecurityRule struct {
	Name       string
	Pattern    *regexp.Regexp
	Severity   string
	Message    string
	Suggestion string
	Rule       string
}

// DockerInstruction represents a parsed Dockerfile instruction
type DockerInstruction struct {
	Command  string
	Args     string
	LineNum  int
	FullLine string
}

// NewDockerValidator creates a new Dockerfile validator
func NewDockerValidator() *DockerValidator {
	validator := &DockerValidator{}
	validator.initializeSecurityRules()
	return validator
}

// ValidateDockerfile validates Dockerfile configuration
func (dv *DockerValidator) ValidateDockerfile(filePath string) (*interfaces.ConfigValidationResult, error) {
	result := &interfaces.ConfigValidationResult{
		Valid:    true,
		Errors:   []interfaces.ConfigValidationError{},
		Warnings: []interfaces.ConfigValidationError{},
		Summary: interfaces.ConfigValidationSummary{
			TotalProperties: 0,
			ValidProperties: 0,
			ErrorCount:      0,
			WarningCount:    0,
			MissingRequired: 0,
		},
	}

	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "file",
			Value:    filePath,
			Type:     "read_error",
			Message:  fmt.Sprintf("Failed to read file: %v", err),
			Severity: interfaces.ValidationSeverityError,
			Rule:     "config.file_access",
		})
		result.Summary.ErrorCount++
		return result, nil
	}

	instructions := dv.parseDockerfile(string(content))

	// Validate basic structure
	dv.validateDockerfileStructure(instructions, result)

	// Validate individual instructions
	for _, instruction := range instructions {
		dv.validateInstruction(instruction, result)
	}

	// Validate security best practices
	dv.validateSecurityPractices(instructions, result)

	// Validate performance best practices
	dv.validatePerformancePractices(instructions, result)

	return result, nil
}

// parseDockerfile parses Dockerfile content into instructions
func (dv *DockerValidator) parseDockerfile(content string) []DockerInstruction {
	var instructions []DockerInstruction
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle line continuations
		fullLine := line
		for strings.HasSuffix(line, "\\") && lineNum+1 < len(lines) {
			lineNum++
			nextLine := strings.TrimSpace(lines[lineNum])
			line = strings.TrimSuffix(line, "\\") + " " + nextLine
			fullLine += "\n" + lines[lineNum]
		}

		// Parse instruction
		parts := strings.SplitN(line, " ", 2)
		if len(parts) >= 1 {
			command := strings.ToUpper(parts[0])
			args := ""
			if len(parts) > 1 {
				args = parts[1]
			}

			instructions = append(instructions, DockerInstruction{
				Command:  command,
				Args:     args,
				LineNum:  lineNum + 1,
				FullLine: fullLine,
			})
		}
	}

	return instructions
}

// validateDockerfileStructure validates the overall structure of the Dockerfile
func (dv *DockerValidator) validateDockerfileStructure(instructions []DockerInstruction, result *interfaces.ConfigValidationResult) {
	hasFrom := false
	fromIndex := -1

	for i, instruction := range instructions {
		result.Summary.TotalProperties++

		switch instruction.Command {
		case "FROM":
			hasFrom = true
			fromIndex = i
			result.Summary.ValidProperties++
		default:
			result.Summary.ValidProperties++
		}
	}

	// Check if FROM instruction exists
	if !hasFrom {
		result.Valid = false
		result.Errors = append(result.Errors, interfaces.ConfigValidationError{
			Field:    "dockerfile",
			Value:    "missing FROM",
			Type:     "missing_instruction",
			Message:  "Dockerfile must contain a FROM instruction",
			Severity: interfaces.ValidationSeverityError,
			Rule:     "docker.from_required",
		})
		result.Summary.ErrorCount++
		result.Summary.MissingRequired++
	}

	// Check if FROM is the first instruction (ignoring ARG for build args)
	if hasFrom && fromIndex > 0 {
		// Check if there are non-ARG instructions before FROM
		for i := 0; i < fromIndex; i++ {
			if instructions[i].Command != "ARG" {
				result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
					Field:      fmt.Sprintf("line_%d", instructions[fromIndex].LineNum),
					Value:      instructions[fromIndex].FullLine,
					Type:       "structure_warning",
					Message:    "FROM should be the first instruction (except for ARG)",
					Suggestion: "Move FROM instruction to the beginning of the Dockerfile",
					Severity:   interfaces.ValidationSeverityWarning,
					Rule:       "docker.from_first",
				})
				result.Summary.WarningCount++
				break
			}
		}
	}
}

// validateInstruction validates individual Dockerfile instructions
func (dv *DockerValidator) validateInstruction(instruction DockerInstruction, result *interfaces.ConfigValidationResult) {
	switch instruction.Command {
	case "FROM":
		dv.validateFromInstruction(instruction, result)
	case "RUN":
		dv.validateRunInstruction(instruction, result)
	case "COPY", "ADD":
		dv.validateCopyAddInstruction(instruction, result)
	case "USER":
		dv.validateUserInstruction(instruction, result)
	case "WORKDIR":
		dv.validateWorkdirInstruction(instruction, result)
	case "EXPOSE":
		dv.validateExposeInstruction(instruction, result)
	case "ENV":
		dv.validateEnvInstruction(instruction, result)
	case "LABEL":
		dv.validateLabelInstruction(instruction, result)
	case "HEALTHCHECK":
		dv.validateHealthcheckInstruction(instruction, result)
	}
}

// validateFromInstruction validates FROM instructions
func (dv *DockerValidator) validateFromInstruction(instruction DockerInstruction, result *interfaces.ConfigValidationResult) {
	// Check for latest tag usage
	if strings.Contains(instruction.Args, ":latest") || !strings.Contains(instruction.Args, ":") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      instruction.Args,
			Type:       "best_practice",
			Message:    "Avoid using 'latest' tag or no tag in production",
			Suggestion: "Use specific version tags for better reproducibility",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker.latest_tag",
		})
		result.Summary.WarningCount++
	}

	// Check for official images without explicit registry
	if !strings.Contains(instruction.Args, "/") && !strings.Contains(instruction.Args, ".") {
		// This is likely an official image, which is good
		return
	}

	// Check for insecure registries (HTTP)
	if strings.HasPrefix(instruction.Args, "http://") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      instruction.Args,
			Type:       "security",
			Message:    "Insecure registry detected (HTTP)",
			Suggestion: "Use HTTPS registries for security",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker.insecure_registry",
		})
		result.Summary.WarningCount++
	}
}

// validateRunInstruction validates RUN instructions
func (dv *DockerValidator) validateRunInstruction(instruction DockerInstruction, result *interfaces.ConfigValidationResult) {
	args := instruction.Args

	// Check for package manager cache cleanup
	if strings.Contains(args, "apt-get install") && !strings.Contains(args, "rm -rf /var/lib/apt/lists/*") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      args,
			Type:       "best_practice",
			Message:    "apt-get install should clean up package cache",
			Suggestion: "Add '&& rm -rf /var/lib/apt/lists/*' to reduce image size",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker.apt_cleanup",
		})
		result.Summary.WarningCount++
	}

	if strings.Contains(args, "yum install") && !strings.Contains(args, "yum clean all") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      args,
			Type:       "best_practice",
			Message:    "yum install should clean up package cache",
			Suggestion: "Add '&& yum clean all' to reduce image size",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker.yum_cleanup",
		})
		result.Summary.WarningCount++
	}

	// Check for dangerous commands
	dangerousCommands := []string{"rm -rf /", "chmod 777", "chown -R root"}
	for _, dangerous := range dangerousCommands {
		if strings.Contains(args, dangerous) {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", instruction.LineNum),
				Value:      args,
				Type:       "security",
				Message:    fmt.Sprintf("Potentially dangerous command detected: %s", dangerous),
				Suggestion: "Review command for security implications",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "docker.dangerous_command",
			})
			result.Summary.WarningCount++
		}
	}

	// Check for curl/wget without verification
	if (strings.Contains(args, "curl") || strings.Contains(args, "wget")) &&
		!strings.Contains(args, "--verify") && !strings.Contains(args, "gpg") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      args,
			Type:       "security",
			Message:    "Downloading files without verification",
			Suggestion: "Verify downloaded files with checksums or GPG signatures",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker.unverified_download",
		})
		result.Summary.WarningCount++
	}
}

// validateCopyAddInstruction validates COPY and ADD instructions
func (dv *DockerValidator) validateCopyAddInstruction(instruction DockerInstruction, result *interfaces.ConfigValidationResult) {
	// Prefer COPY over ADD
	if instruction.Command == "ADD" && !strings.Contains(instruction.Args, "http") && !strings.Contains(instruction.Args, ".tar") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      instruction.Args,
			Type:       "best_practice",
			Message:    "Prefer COPY over ADD for simple file copying",
			Suggestion: "Use COPY instead of ADD unless you need URL fetching or archive extraction",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "docker.prefer_copy",
		})
		result.Summary.WarningCount++
	}

	// Check for copying entire context
	if strings.Contains(instruction.Args, ". /") || strings.Contains(instruction.Args, "./ /") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      instruction.Args,
			Type:       "best_practice",
			Message:    "Copying entire build context",
			Suggestion: "Copy only necessary files to reduce build context and improve caching",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker.copy_context",
		})
		result.Summary.WarningCount++
	}

	// Check for absolute paths in source
	parts := strings.Fields(instruction.Args)
	if len(parts) >= 2 && strings.HasPrefix(parts[0], "/") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      instruction.Args,
			Type:       "best_practice",
			Message:    "Using absolute path in COPY/ADD source",
			Suggestion: "Use relative paths for better portability",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "docker.absolute_path",
		})
		result.Summary.WarningCount++
	}
}

// validateUserInstruction validates USER instructions
func (dv *DockerValidator) validateUserInstruction(instruction DockerInstruction, result *interfaces.ConfigValidationResult) {
	user := strings.TrimSpace(instruction.Args)

	// Check for root user
	if user == "root" || user == "0" {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      user,
			Type:       "security",
			Message:    "Running as root user is a security risk",
			Suggestion: "Create and use a non-root user",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker.root_user",
		})
		result.Summary.WarningCount++
	}

	// Check for numeric UID without name
	if regexp.MustCompile(`^\d+$`).MatchString(user) && user != "0" {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      user,
			Type:       "best_practice",
			Message:    "Using numeric UID without username",
			Suggestion: "Consider using username for better readability",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "docker.numeric_user",
		})
		result.Summary.WarningCount++
	}
}

// validateWorkdirInstruction validates WORKDIR instructions
func (dv *DockerValidator) validateWorkdirInstruction(instruction DockerInstruction, result *interfaces.ConfigValidationResult) {
	workdir := strings.TrimSpace(instruction.Args)

	// Check for relative paths
	if !strings.HasPrefix(workdir, "/") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      workdir,
			Type:       "best_practice",
			Message:    "WORKDIR should use absolute paths",
			Suggestion: "Use absolute paths for WORKDIR to avoid confusion",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker.workdir_absolute",
		})
		result.Summary.WarningCount++
	}
}

// validateExposeInstruction validates EXPOSE instructions
func (dv *DockerValidator) validateExposeInstruction(instruction DockerInstruction, result *interfaces.ConfigValidationResult) {
	ports := strings.Fields(instruction.Args)

	for _, port := range ports {
		// Remove protocol if present
		portNum := strings.Split(port, "/")[0]

		// Check for common insecure ports
		insecurePorts := map[string]string{
			"21":   "FTP (insecure)",
			"23":   "Telnet (insecure)",
			"80":   "HTTP (consider HTTPS on 443)",
			"1433": "SQL Server (consider security)",
			"3306": "MySQL (consider security)",
			"5432": "PostgreSQL (consider security)",
		}

		if description, isInsecure := insecurePorts[portNum]; isInsecure {
			result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
				Field:      fmt.Sprintf("line_%d", instruction.LineNum),
				Value:      port,
				Type:       "security",
				Message:    fmt.Sprintf("Exposing potentially insecure port: %s", description),
				Suggestion: "Consider security implications and use secure alternatives",
				Severity:   interfaces.ValidationSeverityWarning,
				Rule:       "docker.insecure_port",
			})
			result.Summary.WarningCount++
		}
	}
}

// validateEnvInstruction validates ENV instructions
func (dv *DockerValidator) validateEnvInstruction(instruction DockerInstruction, result *interfaces.ConfigValidationResult) {
	// Check for hardcoded secrets
	if dv.containsPotentialSecret(instruction.Args) {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      instruction.Args,
			Type:       "security",
			Message:    "Potential secret in ENV instruction",
			Suggestion: "Use build args or runtime environment variables for secrets",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker.env_secret",
		})
		result.Summary.WarningCount++
	}
}

// validateLabelInstruction validates LABEL instructions
func (dv *DockerValidator) validateLabelInstruction(instruction DockerInstruction, result *interfaces.ConfigValidationResult) {
	// Check for recommended labels
	recommendedLabels := []string{"maintainer", "version", "description"}

	// This is a simple check - in practice, you'd want to track labels across the entire Dockerfile
	for _, label := range recommendedLabels {
		if strings.Contains(strings.ToLower(instruction.Args), label) {
			// Found a recommended label, which is good
			return
		}
	}
}

// validateHealthcheckInstruction validates HEALTHCHECK instructions
func (dv *DockerValidator) validateHealthcheckInstruction(instruction DockerInstruction, result *interfaces.ConfigValidationResult) {
	// Check for NONE healthcheck
	if strings.Contains(instruction.Args, "NONE") {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      fmt.Sprintf("line_%d", instruction.LineNum),
			Value:      instruction.Args,
			Type:       "best_practice",
			Message:    "Healthcheck disabled with NONE",
			Suggestion: "Consider implementing a proper healthcheck for better container management",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "docker.healthcheck_none",
		})
		result.Summary.WarningCount++
	}
}

// validateSecurityPractices validates overall security practices
func (dv *DockerValidator) validateSecurityPractices(instructions []DockerInstruction, result *interfaces.ConfigValidationResult) {
	hasUser := false
	hasHealthcheck := false

	for _, instruction := range instructions {
		switch instruction.Command {
		case "USER":
			hasUser = true
		case "HEALTHCHECK":
			hasHealthcheck = true
		}
	}

	// Check if USER instruction is present
	if !hasUser {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "dockerfile",
			Value:      "missing USER",
			Type:       "security",
			Message:    "No USER instruction found",
			Suggestion: "Add USER instruction to run container as non-root user",
			Severity:   interfaces.ValidationSeverityWarning,
			Rule:       "docker.missing_user",
		})
		result.Summary.WarningCount++
	}

	// Check if HEALTHCHECK instruction is present
	if !hasHealthcheck {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "dockerfile",
			Value:      "missing HEALTHCHECK",
			Type:       "best_practice",
			Message:    "No HEALTHCHECK instruction found",
			Suggestion: "Add HEALTHCHECK instruction for better container monitoring",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "docker.missing_healthcheck",
		})
		result.Summary.WarningCount++
	}
}

// validatePerformancePractices validates performance best practices
func (dv *DockerValidator) validatePerformancePractices(instructions []DockerInstruction, result *interfaces.ConfigValidationResult) {
	runCount := 0
	copyCount := 0

	for _, instruction := range instructions {
		switch instruction.Command {
		case "RUN":
			runCount++
		case "COPY", "ADD":
			copyCount++
		}
	}

	// Check for too many RUN instructions
	if runCount > 10 {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "dockerfile",
			Value:      fmt.Sprintf("%d RUN instructions", runCount),
			Type:       "performance",
			Message:    "Too many RUN instructions can create unnecessary layers",
			Suggestion: "Combine related RUN instructions with && to reduce layers",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "docker.too_many_runs",
		})
		result.Summary.WarningCount++
	}

	// Check for too many COPY instructions
	if copyCount > 5 {
		result.Warnings = append(result.Warnings, interfaces.ConfigValidationError{
			Field:      "dockerfile",
			Value:      fmt.Sprintf("%d COPY instructions", copyCount),
			Type:       "performance",
			Message:    "Multiple COPY instructions can impact build performance",
			Suggestion: "Consider combining COPY operations where possible",
			Severity:   interfaces.ValidationSeverityInfo,
			Rule:       "docker.too_many_copies",
		})
		result.Summary.WarningCount++
	}
}

// containsPotentialSecret checks if instruction arguments contain potential secrets
func (dv *DockerValidator) containsPotentialSecret(args string) bool {
	secretKeywords := []string{"password", "secret", "key", "token", "api"}
	argsLower := strings.ToLower(args)

	for _, keyword := range secretKeywords {
		if strings.Contains(argsLower, keyword) {
			// Check if it looks like a real secret (not just a variable name)
			if regexp.MustCompile(keyword + `\s*=\s*[^\s]+`).MatchString(argsLower) {
				return true
			}
		}
	}

	return false
}

// initializeSecurityRules initializes Docker security rules
func (dv *DockerValidator) initializeSecurityRules() {
	rules := []struct {
		name       string
		pattern    string
		severity   string
		message    string
		suggestion string
		rule       string
	}{
		{
			name:       "root_user",
			pattern:    `USER\s+(root|0)`,
			severity:   interfaces.ValidationSeverityWarning,
			message:    "Running as root user",
			suggestion: "Use non-root user for security",
			rule:       "docker.root_user",
		},
		{
			name:       "latest_tag",
			pattern:    `FROM\s+[^:]+(:latest)?$`,
			severity:   interfaces.ValidationSeverityWarning,
			message:    "Using latest tag",
			suggestion: "Use specific version tags",
			rule:       "docker.latest_tag",
		},
	}

	for _, r := range rules {
		compiled := regexp.MustCompile(r.pattern)
		dv.securityRules = append(dv.securityRules, DockerSecurityRule{
			Name:       r.name,
			Pattern:    compiled,
			Severity:   r.severity,
			Message:    r.message,
			Suggestion: r.suggestion,
			Rule:       r.rule,
		})
	}
}
