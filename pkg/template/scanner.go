package template

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/utils"
)

// TemplateScanner analyzes template configurations and identifies inconsistencies
type TemplateScanner struct {
	templateDir string
}

// NewTemplateScanner creates a new template scanner
func NewTemplateScanner(templateDir string) *TemplateScanner {
	return &TemplateScanner{
		templateDir: templateDir,
	}
}

// ConfigurationAnalysis represents the analysis results for template configurations
type ConfigurationAnalysis struct {
	Templates          []TemplateInfo            `json:"templates"`
	Inconsistencies    []Inconsistency           `json:"inconsistencies"`
	MissingFiles       []MissingFile             `json:"missing_files"`
	VersionReferences  map[string][]string       `json:"version_references"`
	DependencyPatterns map[string]DependencyInfo `json:"dependency_patterns"`
}

// TemplateInfo contains information about a single template
type TemplateInfo struct {
	Name            string            `json:"name"`
	Path            string            `json:"path"`
	Type            string            `json:"type"`
	ConfigFiles     []string          `json:"config_files"`
	PackageJSON     *PackageJSONInfo  `json:"package_json,omitempty"`
	Scripts         map[string]string `json:"scripts,omitempty"`
	Dependencies    []string          `json:"dependencies,omitempty"`
	DevDependencies []string          `json:"dev_dependencies,omitempty"`
	Port            string            `json:"port,omitempty"`
}

// PackageJSONInfo contains parsed package.json information
type PackageJSONInfo struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Engines         map[string]string `json:"engines"`
}

// Inconsistency represents a configuration inconsistency between templates
type Inconsistency struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Templates   []string `json:"templates"`
	Details     string   `json:"details"`
}

// MissingFile represents a configuration file missing from a template
type MissingFile struct {
	Template string `json:"template"`
	File     string `json:"file"`
	Reason   string `json:"reason"`
}

// DependencyInfo contains information about dependency usage patterns
type DependencyInfo struct {
	Package   string   `json:"package"`
	Versions  []string `json:"versions"`
	Templates []string `json:"templates"`
}

// ScanFrontendTemplates analyzes all frontend templates and returns configuration analysis
func (s *TemplateScanner) ScanFrontendTemplates() (*ConfigurationAnalysis, error) {
	analysis := &ConfigurationAnalysis{
		Templates:          []TemplateInfo{},
		Inconsistencies:    []Inconsistency{},
		MissingFiles:       []MissingFile{},
		VersionReferences:  make(map[string][]string),
		DependencyPatterns: make(map[string]DependencyInfo),
	}

	frontendDir := filepath.Join(s.templateDir, constants.TemplateFrontend)

	// Scan each frontend template
	templates := []string{"nextjs-app", "nextjs-home", "nextjs-admin"}

	for _, templateName := range templates {
		templatePath := filepath.Join(frontendDir, templateName)
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			continue
		}

		templateInfo, err := s.analyzeTemplate(templateName, templatePath)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze template %s: %w", templateName, err)
		}

		analysis.Templates = append(analysis.Templates, *templateInfo)
	}

	// Analyze inconsistencies
	s.findInconsistencies(analysis)
	s.findMissingFiles(analysis)
	s.analyzeVersionReferences(analysis)
	s.analyzeDependencyPatterns(analysis)

	return analysis, nil
}

// analyzeTemplate analyzes a single template directory
func (s *TemplateScanner) analyzeTemplate(name, path string) (*TemplateInfo, error) {
	info := &TemplateInfo{
		Name:        name,
		Path:        path,
		Type:        s.determineTemplateType(name),
		ConfigFiles: []string{},
		Scripts:     make(map[string]string),
	}

	// Walk through template directory to find configuration files
	err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileName := d.Name()

		// Check for configuration files
		if s.isConfigFile(fileName) {
			relPath, _ := filepath.Rel(path, filePath)
			info.ConfigFiles = append(info.ConfigFiles, relPath)

			// Parse package.json specifically
			if fileName == "package.json.tmpl" {
				packageInfo, err := s.parsePackageJSON(filePath)
				if err == nil {
					info.PackageJSON = packageInfo
					info.Scripts = packageInfo.Scripts

					// Extract port from dev script
					if devScript, exists := packageInfo.Scripts["dev"]; exists {
						info.Port = s.extractPortFromScript(devScript)
					}

					// Convert dependencies to slices
					for dep := range packageInfo.Dependencies {
						info.Dependencies = append(info.Dependencies, dep)
					}
					for dep := range packageInfo.DevDependencies {
						info.DevDependencies = append(info.DevDependencies, dep)
					}
				}
			}
		}

		return nil
	})

	return info, err
}

// isConfigFile checks if a file is a configuration file
func (s *TemplateScanner) isConfigFile(fileName string) bool {
	configFiles := []string{
		"package.json.tmpl",
		"tsconfig.json.tmpl",
		".eslintrc.json.tmpl",
		".prettierrc.tmpl",
		"next.config.js.tmpl",
		"tailwind.config.js.tmpl",
		"vercel.json.tmpl",
		"jest.config.js.tmpl",
		"postcss.config.js.tmpl",
		".gitignore.tmpl",
		".env.local.example.tmpl",
	}

	for _, configFile := range configFiles {
		if fileName == configFile {
			return true
		}
	}
	return false
}

// parsePackageJSON parses a package.json template file
func (s *TemplateScanner) parsePackageJSON(filePath string) (*PackageJSONInfo, error) {
	content, err := utils.SafeReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Remove template syntax for basic parsing (this is a simplified approach)
	contentStr := string(content)

	// For now, we'll extract the basic structure without full template parsing
	// In a real implementation, we'd need to handle Go template syntax properly

	var packageInfo PackageJSONInfo

	// Try to parse as JSON (will work for static parts)
	// This is a simplified approach - in production we'd need proper template parsing
	if err := json.Unmarshal(content, &packageInfo); err != nil {
		// If JSON parsing fails, extract information manually
		packageInfo = PackageJSONInfo{
			Scripts:         make(map[string]string),
			Dependencies:    make(map[string]string),
			DevDependencies: make(map[string]string),
			Engines:         make(map[string]string),
		}

		// Extract scripts section manually
		s.extractScriptsFromTemplate(contentStr, &packageInfo)
	}

	return &packageInfo, nil
}

// extractScriptsFromTemplate extracts scripts from template content
func (s *TemplateScanner) extractScriptsFromTemplate(content string, packageInfo *PackageJSONInfo) {
	// First extract dependencies
	s.extractDependenciesFromTemplate(content, packageInfo)

	// Then extract scripts
	lines := strings.Split(content, "\n")
	inScripts := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, `"scripts"`) {
			inScripts = true
			continue
		}

		if inScripts {
			if strings.Contains(line, "}") && !strings.Contains(line, ":") {
				inScripts = false
				continue
			}

			// Extract script name and command
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					scriptName := strings.Trim(strings.TrimSpace(parts[0]), `"`)
					scriptCmd := strings.Trim(strings.TrimSpace(parts[1]), `",`)
					scriptCmd = strings.Trim(scriptCmd, `"`)
					packageInfo.Scripts[scriptName] = scriptCmd
				}
			}
		}
	}
}

// extractDependenciesFromTemplate extracts dependency names from template content
func (s *TemplateScanner) extractDependenciesFromTemplate(content string, packageInfo *PackageJSONInfo) {
	lines := strings.Split(content, "\n")
	inDependencies := false
	inDevDependencies := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for dependencies sections
		if strings.Contains(line, `"dependencies"`) && strings.Contains(line, "{") {
			inDependencies = true
			inDevDependencies = false
			continue
		}
		if strings.Contains(line, `"devDependencies"`) && strings.Contains(line, "{") {
			inDevDependencies = true
			inDependencies = false
			continue
		}

		// End of section
		if (inDependencies || inDevDependencies) && strings.Contains(line, "}") {
			inDependencies = false
			inDevDependencies = false
			continue
		}

		// Extract dependency names
		if inDependencies || inDevDependencies {
			if strings.Contains(line, ":") && strings.Contains(line, `"`) {
				// Extract package name from line like: "react": "{{.Versions.React}}"
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					packageName := strings.Trim(strings.TrimSpace(parts[0]), `"`)
					if packageName != "" {
						if inDependencies {
							packageInfo.Dependencies[packageName] = constants.FileTypeTemplate
						} else {
							packageInfo.DevDependencies[packageName] = constants.FileTypeTemplate
						}
					}
				}
			}
		}
	}
}

// extractPortFromScript extracts port number from dev script
func (s *TemplateScanner) extractPortFromScript(script string) string {
	if strings.Contains(script, "-p ") {
		parts := strings.Split(script, "-p ")
		if len(parts) > 1 {
			portPart := strings.Fields(parts[1])[0]
			return portPart
		}
	}
	return "3000" // default
}

// determineTemplateType determines the type of template based on name
func (s *TemplateScanner) determineTemplateType(name string) string {
	switch name {
	case "nextjs-app":
		return "application"
	case "nextjs-home":
		return "landing"
	case "nextjs-admin":
		return "dashboard"
	default:
		return "unknown"
	}
}

// findInconsistencies identifies configuration inconsistencies between templates
func (s *TemplateScanner) findInconsistencies(analysis *ConfigurationAnalysis) {
	// Check for script inconsistencies
	s.checkScriptInconsistencies(analysis)

	// Check for missing configuration files
	s.checkConfigFileInconsistencies(analysis)

	// Check for port conflicts
	s.checkPortInconsistencies(analysis)

	// Check for dependency version inconsistencies
	s.checkDependencyInconsistencies(analysis)
}

// checkScriptInconsistencies checks for inconsistent npm scripts
func (s *TemplateScanner) checkScriptInconsistencies(analysis *ConfigurationAnalysis) {
	scriptMap := make(map[string]map[string]string) // script name -> template -> command

	for _, template := range analysis.Templates {
		for scriptName, scriptCmd := range template.Scripts {
			if scriptMap[scriptName] == nil {
				scriptMap[scriptName] = make(map[string]string)
			}
			scriptMap[scriptName][template.Name] = scriptCmd
		}
	}

	for scriptName, templateCommands := range scriptMap {
		// Skip port-specific scripts as they're expected to differ
		if scriptName == "dev" || scriptName == "start" {
			continue
		}

		commands := make(map[string][]string)
		for template, command := range templateCommands {
			commands[command] = append(commands[command], template)
		}

		if len(commands) > 1 {
			var details []string
			for command, templates := range commands {
				details = append(details, fmt.Sprintf("%s: %s", strings.Join(templates, ", "), command))
			}

			analysis.Inconsistencies = append(analysis.Inconsistencies, Inconsistency{
				Type:        "script_inconsistency",
				Description: fmt.Sprintf("Script '%s' has different implementations", scriptName),
				Templates:   getAllTemplateNames(templateCommands),
				Details:     strings.Join(details, "; "),
			})
		}
	}
}

// checkConfigFileInconsistencies checks for missing configuration files
func (s *TemplateScanner) checkConfigFileInconsistencies(analysis *ConfigurationAnalysis) {
	allConfigFiles := make(map[string]bool)
	templateFiles := make(map[string][]string)

	// Collect all config files across templates
	for _, template := range analysis.Templates {
		templateFiles[template.Name] = template.ConfigFiles
		for _, file := range template.ConfigFiles {
			allConfigFiles[file] = true
		}
	}

	// Check which files are missing from each template
	for fileName := range allConfigFiles {
		var templatesWithFile []string
		var templatesWithoutFile []string

		for templateName, files := range templateFiles {
			hasFile := false
			for _, file := range files {
				if file == fileName {
					hasFile = true
					break
				}
			}

			if hasFile {
				templatesWithFile = append(templatesWithFile, templateName)
			} else {
				templatesWithoutFile = append(templatesWithoutFile, templateName)
			}
		}

		if len(templatesWithoutFile) > 0 && len(templatesWithFile) > 0 {
			analysis.Inconsistencies = append(analysis.Inconsistencies, Inconsistency{
				Type:        "missing_config_file",
				Description: fmt.Sprintf("Configuration file '%s' is missing from some templates", fileName),
				Templates:   templatesWithoutFile,
				Details:     fmt.Sprintf("Present in: %s; Missing from: %s", strings.Join(templatesWithFile, ", "), strings.Join(templatesWithoutFile, ", ")),
			})
		}
	}
}

// checkPortInconsistencies checks for port conflicts
func (s *TemplateScanner) checkPortInconsistencies(analysis *ConfigurationAnalysis) {
	portMap := make(map[string][]string)

	for _, template := range analysis.Templates {
		if template.Port != "" {
			portMap[template.Port] = append(portMap[template.Port], template.Name)
		}
	}

	for port, templates := range portMap {
		if len(templates) > 1 {
			analysis.Inconsistencies = append(analysis.Inconsistencies, Inconsistency{
				Type:        "port_conflict",
				Description: fmt.Sprintf("Multiple templates use the same port %s", port),
				Templates:   templates,
				Details:     fmt.Sprintf("Port %s is used by: %s", port, strings.Join(templates, ", ")),
			})
		}
	}
}

// checkDependencyInconsistencies checks for dependency version inconsistencies
func (s *TemplateScanner) checkDependencyInconsistencies(analysis *ConfigurationAnalysis) {
	// This would require parsing the actual dependency versions from package.json
	// For now, we'll add a placeholder for this functionality
	// In a real implementation, we'd parse the template variables and compare versions
}

// findMissingFiles identifies files that should be present in all templates
func (s *TemplateScanner) findMissingFiles(analysis *ConfigurationAnalysis) {
	requiredFiles := []string{
		"package.json.tmpl",
		"next.config.js.tmpl",
		"tailwind.config.js.tmpl",
	}

	for _, template := range analysis.Templates {
		for _, requiredFile := range requiredFiles {
			hasFile := false
			for _, file := range template.ConfigFiles {
				if file == requiredFile {
					hasFile = true
					break
				}
			}

			if !hasFile {
				analysis.MissingFiles = append(analysis.MissingFiles, MissingFile{
					Template: template.Name,
					File:     requiredFile,
					Reason:   "Required for consistent frontend template configuration",
				})
			}
		}
	}
}

// analyzeVersionReferences analyzes version references in templates
func (s *TemplateScanner) analyzeVersionReferences(analysis *ConfigurationAnalysis) {
	// This would parse template files to find version references like {{.Versions.NextJS}}
	// For now, we'll add common version references
	analysis.VersionReferences["NextJS"] = []string{"nextjs-app", "nextjs-home", "nextjs-admin"}
	analysis.VersionReferences["React"] = []string{"nextjs-app", "nextjs-home", "nextjs-admin"}
	analysis.VersionReferences["Node"] = []string{"nextjs-app", "nextjs-home", "nextjs-admin"}
}

// analyzeDependencyPatterns analyzes dependency usage patterns
func (s *TemplateScanner) analyzeDependencyPatterns(analysis *ConfigurationAnalysis) {
	dependencyMap := make(map[string]map[string]bool) // package -> template -> exists

	for _, template := range analysis.Templates {
		for _, dep := range template.Dependencies {
			if dependencyMap[dep] == nil {
				dependencyMap[dep] = make(map[string]bool)
			}
			dependencyMap[dep][template.Name] = true
		}

		for _, dep := range template.DevDependencies {
			if dependencyMap[dep] == nil {
				dependencyMap[dep] = make(map[string]bool)
			}
			dependencyMap[dep][template.Name] = true
		}
	}

	for packageName, templateMap := range dependencyMap {
		var templates []string
		for template := range templateMap {
			templates = append(templates, template)
		}

		analysis.DependencyPatterns[packageName] = DependencyInfo{
			Package:   packageName,
			Templates: templates,
			Versions:  []string{}, // Would be populated with actual version parsing
		}
	}
}

// getAllTemplateNames extracts all template names from a map
func getAllTemplateNames(templateMap map[string]string) []string {
	var names []string
	for name := range templateMap {
		names = append(names, name)
	}
	return names
}
