package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// TransformForExport transforms a configuration for export
func (t *ConfigTransformer) TransformForExport(config *SavedConfiguration, options *ConfigExportOptions) (interface{}, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	// Create export structure
	exportData := make(map[string]interface{})

	// Include metadata if requested
	if options.IncludeMeta {
		exportData["metadata"] = map[string]interface{}{
			"name":           config.Name,
			"description":    config.Description,
			"version":        config.Version,
			"created_at":     config.CreatedAt,
			"updated_at":     config.UpdatedAt,
			"tags":           config.Tags,
			"exported_at":    time.Now(),
			"export_version": "1.0.0",
		}
	}

	// Add project configuration
	if config.ProjectConfig != nil {
		projectData, err := t.transformProjectConfig(config.ProjectConfig, options)
		if err != nil {
			return nil, fmt.Errorf("failed to transform project config: %w", err)
		}
		exportData["project"] = projectData
	}

	// Add selected templates
	if len(config.SelectedTemplates) > 0 {
		templatesData, err := t.transformTemplateSelections(config.SelectedTemplates, options)
		if err != nil {
			return nil, fmt.Errorf("failed to transform template selections: %w", err)
		}
		exportData["templates"] = templatesData
	}

	// Add generation settings
	if config.GenerationSettings != nil {
		settingsData, err := t.transformGenerationSettings(config.GenerationSettings, options)
		if err != nil {
			return nil, fmt.Errorf("failed to transform generation settings: %w", err)
		}
		exportData["generation"] = settingsData
	}

	// Add user preferences
	if config.UserPreferences != nil {
		prefsData, err := t.transformUserPreferences(config.UserPreferences, options)
		if err != nil {
			return nil, fmt.Errorf("failed to transform user preferences: %w", err)
		}
		exportData["preferences"] = prefsData
	}

	// Apply field exclusions
	if len(options.ExcludeFields) > 0 {
		exportData = t.excludeFields(exportData, options.ExcludeFields)
	}

	return exportData, nil
}

// TransformForImport transforms imported data into a configuration
func (t *ConfigTransformer) TransformForImport(data interface{}, options *ConfigImportOptions) (*SavedConfiguration, error) {
	if data == nil {
		return nil, fmt.Errorf("import data cannot be nil")
	}

	// Convert to map for easier processing
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("import data must be an object/map")
	}

	config := &SavedConfiguration{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   "1.0.0",
	}

	// Extract metadata if present
	if metaData, exists := dataMap["metadata"]; exists {
		if err := t.extractMetadata(config, metaData, options); err != nil {
			return nil, fmt.Errorf("failed to extract metadata: %w", err)
		}
	}

	// Extract project configuration
	if projectData, exists := dataMap["project"]; exists {
		projectConfig, err := t.extractProjectConfig(projectData, options)
		if err != nil {
			return nil, fmt.Errorf("failed to extract project config: %w", err)
		}
		config.ProjectConfig = projectConfig
	}

	// Extract template selections
	if templatesData, exists := dataMap["templates"]; exists {
		templates, err := t.extractTemplateSelections(templatesData, options)
		if err != nil {
			return nil, fmt.Errorf("failed to extract template selections: %w", err)
		}
		config.SelectedTemplates = templates
	}

	// Extract generation settings
	if settingsData, exists := dataMap["generation"]; exists {
		settings, err := t.extractGenerationSettings(settingsData, options)
		if err != nil {
			return nil, fmt.Errorf("failed to extract generation settings: %w", err)
		}
		config.GenerationSettings = settings
	}

	// Extract user preferences
	if prefsData, exists := dataMap["preferences"]; exists {
		prefs, err := t.extractUserPreferences(prefsData, options)
		if err != nil {
			return nil, fmt.Errorf("failed to extract user preferences: %w", err)
		}
		config.UserPreferences = prefs
	}

	// Apply field mappings if specified
	if len(options.FieldMappings) > 0 {
		if err := t.applyFieldMappings(config, options.FieldMappings); err != nil {
			return nil, fmt.Errorf("failed to apply field mappings: %w", err)
		}
	}

	// Apply transformations if requested
	if options.Transform {
		if err := t.applyImportTransformations(config); err != nil {
			return nil, fmt.Errorf("failed to apply transformations: %w", err)
		}
	}

	return config, nil
}

// transformProjectConfig transforms project configuration for export
func (t *ConfigTransformer) transformProjectConfig(config *models.ProjectConfig, options *ConfigExportOptions) (interface{}, error) {
	data := map[string]interface{}{
		"name":         config.Name,
		"organization": config.Organization,
		"description":  config.Description,
		"license":      config.License,
		"author":       config.Author,
		"email":        config.Email,
		"repository":   config.Repository,
		"output_path":  config.OutputPath,
		"features":     config.Features,
	}

	// Transform components
	if !t.isEmptyComponents(&config.Components) {
		data["components"] = t.transformComponents(&config.Components)
	}

	// Transform versions
	if config.Versions != nil {
		data["versions"] = t.transformVersions(config.Versions)
	}

	// Add generation metadata if including meta
	if options.IncludeMeta {
		data["generated_at"] = config.GeneratedAt
		data["generator_version"] = config.GeneratorVersion
	}

	return data, nil
}

// transformTemplateSelections transforms template selections for export
func (t *ConfigTransformer) transformTemplateSelections(templates []TemplateSelection, options *ConfigExportOptions) (interface{}, error) {
	var result []map[string]interface{}

	for _, template := range templates {
		templateData := map[string]interface{}{
			"name":       template.TemplateName,
			"category":   template.Category,
			"technology": template.Technology,
			"version":    template.Version,
			"selected":   template.Selected,
		}

		if len(template.Options) > 0 {
			templateData["options"] = template.Options
		}

		result = append(result, templateData)
	}

	return result, nil
}

// transformGenerationSettings transforms generation settings for export
func (t *ConfigTransformer) transformGenerationSettings(settings *GenerationSettings, options *ConfigExportOptions) (interface{}, error) {
	return map[string]interface{}{
		"include_examples":   settings.IncludeExamples,
		"include_tests":      settings.IncludeTests,
		"include_docs":       settings.IncludeDocs,
		"update_versions":    settings.UpdateVersions,
		"minimal_mode":       settings.MinimalMode,
		"exclude_patterns":   settings.ExcludePatterns,
		"include_only_paths": settings.IncludeOnlyPaths,
		"backup_existing":    settings.BackupExisting,
		"overwrite_existing": settings.OverwriteExisting,
	}, nil
}

// transformUserPreferences transforms user preferences for export
func (t *ConfigTransformer) transformUserPreferences(prefs *UserPreferences, options *ConfigExportOptions) (interface{}, error) {
	data := map[string]interface{}{
		"default_license":      prefs.DefaultLicense,
		"default_author":       prefs.DefaultAuthor,
		"default_email":        prefs.DefaultEmail,
		"default_organization": prefs.DefaultOrganization,
		"preferred_format":     prefs.PreferredFormat,
	}

	if len(prefs.CustomDefaults) > 0 {
		data["custom_defaults"] = prefs.CustomDefaults
	}

	return data, nil
}

// transformComponents transforms component configuration
func (t *ConfigTransformer) transformComponents(components *models.Components) interface{} {
	return map[string]interface{}{
		"frontend": map[string]interface{}{
			"nextjs": map[string]interface{}{
				"app":    components.Frontend.NextJS.App,
				"home":   components.Frontend.NextJS.Home,
				"admin":  components.Frontend.NextJS.Admin,
				"shared": components.Frontend.NextJS.Shared,
			},
		},
		"backend": map[string]interface{}{
			"go_gin": components.Backend.GoGin,
		},
		"mobile": map[string]interface{}{
			"android": components.Mobile.Android,
			"ios":     components.Mobile.IOS,
			"shared":  components.Mobile.Shared,
		},
		"infrastructure": map[string]interface{}{
			"docker":     components.Infrastructure.Docker,
			"kubernetes": components.Infrastructure.Kubernetes,
			"terraform":  components.Infrastructure.Terraform,
		},
		"database": map[string]interface{}{
			"postgresql": components.Database.PostgreSQL,
			"mysql":      components.Database.MySQL,
			"mongodb":    components.Database.MongoDB,
			"sqlite":     components.Database.SQLite,
		},
		"cache": map[string]interface{}{
			"redis":     components.Cache.Redis,
			"memcached": components.Cache.Memcached,
		},
		"devops": map[string]interface{}{
			"cicd":           components.DevOps.CICD,
			"github_actions": components.DevOps.GitHubActions,
			"gitlab_ci":      components.DevOps.GitLabCI,
			"jenkins":        components.DevOps.Jenkins,
		},
		"monitoring": map[string]interface{}{
			"prometheus": components.Monitoring.Prometheus,
			"grafana":    components.Monitoring.Grafana,
			"jaeger":     components.Monitoring.Jaeger,
			"elk":        components.Monitoring.ELK,
		},
	}
}

// transformVersions transforms version configuration
func (t *ConfigTransformer) transformVersions(versions *models.VersionConfig) interface{} {
	data := map[string]interface{}{
		"node": versions.Node,
		"go":   versions.Go,
	}

	if len(versions.Packages) > 0 {
		data["packages"] = versions.Packages
	}

	return data
}

// extractMetadata extracts metadata from import data
func (t *ConfigTransformer) extractMetadata(config *SavedConfiguration, data interface{}, options *ConfigImportOptions) error {
	metaMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("metadata must be an object")
	}

	if name, exists := metaMap["name"]; exists {
		if nameStr, ok := name.(string); ok {
			config.Name = nameStr
		}
	}

	if desc, exists := metaMap["description"]; exists {
		if descStr, ok := desc.(string); ok {
			config.Description = descStr
		}
	}

	if version, exists := metaMap["version"]; exists {
		if versionStr, ok := version.(string); ok {
			config.Version = versionStr
		}
	}

	if tags, exists := metaMap["tags"]; exists {
		if tagsList, ok := tags.([]interface{}); ok {
			var tagStrings []string
			for _, tag := range tagsList {
				if tagStr, ok := tag.(string); ok {
					tagStrings = append(tagStrings, tagStr)
				}
			}
			config.Tags = tagStrings
		}
	}

	return nil
}

// extractProjectConfig extracts project configuration from import data
func (t *ConfigTransformer) extractProjectConfig(data interface{}, options *ConfigImportOptions) (*models.ProjectConfig, error) {
	projectMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("project config must be an object")
	}

	config := &models.ProjectConfig{}

	// Extract basic fields
	if name, exists := projectMap["name"]; exists {
		if nameStr, ok := name.(string); ok {
			config.Name = nameStr
		}
	}

	if org, exists := projectMap["organization"]; exists {
		if orgStr, ok := org.(string); ok {
			config.Organization = orgStr
		}
	}

	if desc, exists := projectMap["description"]; exists {
		if descStr, ok := desc.(string); ok {
			config.Description = descStr
		}
	}

	if license, exists := projectMap["license"]; exists {
		if licenseStr, ok := license.(string); ok {
			config.License = licenseStr
		}
	}

	if author, exists := projectMap["author"]; exists {
		if authorStr, ok := author.(string); ok {
			config.Author = authorStr
		}
	}

	if email, exists := projectMap["email"]; exists {
		if emailStr, ok := email.(string); ok {
			config.Email = emailStr
		}
	}

	if repo, exists := projectMap["repository"]; exists {
		if repoStr, ok := repo.(string); ok {
			config.Repository = repoStr
		}
	}

	if output, exists := projectMap["output_path"]; exists {
		if outputStr, ok := output.(string); ok {
			config.OutputPath = outputStr
		}
	}

	// Extract features
	if features, exists := projectMap["features"]; exists {
		if featuresList, ok := features.([]interface{}); ok {
			var featureStrings []string
			for _, feature := range featuresList {
				if featureStr, ok := feature.(string); ok {
					featureStrings = append(featureStrings, featureStr)
				}
			}
			config.Features = featureStrings
		}
	}

	// Extract components
	if components, exists := projectMap["components"]; exists {
		if err := t.extractComponents(&config.Components, components); err != nil {
			return nil, fmt.Errorf("failed to extract components: %w", err)
		}
	}

	// Extract versions
	if versions, exists := projectMap["versions"]; exists {
		versionConfig, err := t.extractVersionConfig(versions)
		if err != nil {
			return nil, fmt.Errorf("failed to extract versions: %w", err)
		}
		config.Versions = versionConfig
	}

	return config, nil
}

// extractTemplateSelections extracts template selections from import data
func (t *ConfigTransformer) extractTemplateSelections(data interface{}, options *ConfigImportOptions) ([]TemplateSelection, error) {
	templatesList, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("templates must be an array")
	}

	var templates []TemplateSelection

	for _, templateData := range templatesList {
		templateMap, ok := templateData.(map[string]interface{})
		if !ok {
			continue
		}

		template := TemplateSelection{}

		if name, exists := templateMap["name"]; exists {
			if nameStr, ok := name.(string); ok {
				template.TemplateName = nameStr
			}
		}

		if category, exists := templateMap["category"]; exists {
			if categoryStr, ok := category.(string); ok {
				template.Category = categoryStr
			}
		}

		if technology, exists := templateMap["technology"]; exists {
			if techStr, ok := technology.(string); ok {
				template.Technology = techStr
			}
		}

		if version, exists := templateMap["version"]; exists {
			if versionStr, ok := version.(string); ok {
				template.Version = versionStr
			}
		}

		if selected, exists := templateMap["selected"]; exists {
			if selectedBool, ok := selected.(bool); ok {
				template.Selected = selectedBool
			}
		}

		if options, exists := templateMap["options"]; exists {
			if optionsMap, ok := options.(map[string]interface{}); ok {
				template.Options = optionsMap
			}
		}

		templates = append(templates, template)
	}

	return templates, nil
}

// extractGenerationSettings extracts generation settings from import data
func (t *ConfigTransformer) extractGenerationSettings(data interface{}, options *ConfigImportOptions) (*GenerationSettings, error) {
	settingsMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("generation settings must be an object")
	}

	settings := &GenerationSettings{}

	if includeExamples, exists := settingsMap["include_examples"]; exists {
		if includeExamplesBool, ok := includeExamples.(bool); ok {
			settings.IncludeExamples = includeExamplesBool
		}
	}

	if includeTests, exists := settingsMap["include_tests"]; exists {
		if includeTestsBool, ok := includeTests.(bool); ok {
			settings.IncludeTests = includeTestsBool
		}
	}

	if includeDocs, exists := settingsMap["include_docs"]; exists {
		if includeDocsBool, ok := includeDocs.(bool); ok {
			settings.IncludeDocs = includeDocsBool
		}
	}

	if updateVersions, exists := settingsMap["update_versions"]; exists {
		if updateVersionsBool, ok := updateVersions.(bool); ok {
			settings.UpdateVersions = updateVersionsBool
		}
	}

	if minimalMode, exists := settingsMap["minimal_mode"]; exists {
		if minimalModeBool, ok := minimalMode.(bool); ok {
			settings.MinimalMode = minimalModeBool
		}
	}

	if backupExisting, exists := settingsMap["backup_existing"]; exists {
		if backupExistingBool, ok := backupExisting.(bool); ok {
			settings.BackupExisting = backupExistingBool
		}
	}

	if overwriteExisting, exists := settingsMap["overwrite_existing"]; exists {
		if overwriteExistingBool, ok := overwriteExisting.(bool); ok {
			settings.OverwriteExisting = overwriteExistingBool
		}
	}

	// Extract string arrays
	if excludePatterns, exists := settingsMap["exclude_patterns"]; exists {
		if patternsList, ok := excludePatterns.([]interface{}); ok {
			var patterns []string
			for _, pattern := range patternsList {
				if patternStr, ok := pattern.(string); ok {
					patterns = append(patterns, patternStr)
				}
			}
			settings.ExcludePatterns = patterns
		}
	}

	if includeOnlyPaths, exists := settingsMap["include_only_paths"]; exists {
		if pathsList, ok := includeOnlyPaths.([]interface{}); ok {
			var paths []string
			for _, path := range pathsList {
				if pathStr, ok := path.(string); ok {
					paths = append(paths, pathStr)
				}
			}
			settings.IncludeOnlyPaths = paths
		}
	}

	return settings, nil
}

// extractUserPreferences extracts user preferences from import data
func (t *ConfigTransformer) extractUserPreferences(data interface{}, options *ConfigImportOptions) (*UserPreferences, error) {
	prefsMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("user preferences must be an object")
	}

	prefs := &UserPreferences{}

	if defaultLicense, exists := prefsMap["default_license"]; exists {
		if licenseStr, ok := defaultLicense.(string); ok {
			prefs.DefaultLicense = licenseStr
		}
	}

	if defaultAuthor, exists := prefsMap["default_author"]; exists {
		if authorStr, ok := defaultAuthor.(string); ok {
			prefs.DefaultAuthor = authorStr
		}
	}

	if defaultEmail, exists := prefsMap["default_email"]; exists {
		if emailStr, ok := defaultEmail.(string); ok {
			prefs.DefaultEmail = emailStr
		}
	}

	if defaultOrg, exists := prefsMap["default_organization"]; exists {
		if orgStr, ok := defaultOrg.(string); ok {
			prefs.DefaultOrganization = orgStr
		}
	}

	if preferredFormat, exists := prefsMap["preferred_format"]; exists {
		if formatStr, ok := preferredFormat.(string); ok {
			prefs.PreferredFormat = formatStr
		}
	}

	if customDefaults, exists := prefsMap["custom_defaults"]; exists {
		if defaultsMap, ok := customDefaults.(map[string]interface{}); ok {
			defaults := make(map[string]string)
			for key, value := range defaultsMap {
				if valueStr, ok := value.(string); ok {
					defaults[key] = valueStr
				}
			}
			prefs.CustomDefaults = defaults
		}
	}

	return prefs, nil
}

// extractComponents extracts component configuration from import data
func (t *ConfigTransformer) extractComponents(components *models.Components, data interface{}) error {
	componentsMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("components must be an object")
	}

	// Extract frontend components
	if frontend, exists := componentsMap["frontend"]; exists {
		if frontendMap, ok := frontend.(map[string]interface{}); ok {
			if nextjs, exists := frontendMap["nextjs"]; exists {
				if nextjsMap, ok := nextjs.(map[string]interface{}); ok {
					if app, exists := nextjsMap["app"]; exists {
						if appBool, ok := app.(bool); ok {
							components.Frontend.NextJS.App = appBool
						}
					}
					if home, exists := nextjsMap["home"]; exists {
						if homeBool, ok := home.(bool); ok {
							components.Frontend.NextJS.Home = homeBool
						}
					}
					if admin, exists := nextjsMap["admin"]; exists {
						if adminBool, ok := admin.(bool); ok {
							components.Frontend.NextJS.Admin = adminBool
						}
					}
					if shared, exists := nextjsMap["shared"]; exists {
						if sharedBool, ok := shared.(bool); ok {
							components.Frontend.NextJS.Shared = sharedBool
						}
					}
				}
			}
		}
	}

	// Extract backend components
	if backend, exists := componentsMap["backend"]; exists {
		if backendMap, ok := backend.(map[string]interface{}); ok {
			if goGin, exists := backendMap["go_gin"]; exists {
				if goGinBool, ok := goGin.(bool); ok {
					components.Backend.GoGin = goGinBool
				}
			}
		}
	}

	// Extract mobile components
	if mobile, exists := componentsMap["mobile"]; exists {
		if mobileMap, ok := mobile.(map[string]interface{}); ok {
			if android, exists := mobileMap["android"]; exists {
				if androidBool, ok := android.(bool); ok {
					components.Mobile.Android = androidBool
				}
			}
			if ios, exists := mobileMap["ios"]; exists {
				if iosBool, ok := ios.(bool); ok {
					components.Mobile.IOS = iosBool
				}
			}
			if shared, exists := mobileMap["shared"]; exists {
				if sharedBool, ok := shared.(bool); ok {
					components.Mobile.Shared = sharedBool
				}
			}
		}
	}

	// Extract infrastructure components
	if infrastructure, exists := componentsMap["infrastructure"]; exists {
		if infraMap, ok := infrastructure.(map[string]interface{}); ok {
			if docker, exists := infraMap["docker"]; exists {
				if dockerBool, ok := docker.(bool); ok {
					components.Infrastructure.Docker = dockerBool
				}
			}
			if kubernetes, exists := infraMap["kubernetes"]; exists {
				if kubeBool, ok := kubernetes.(bool); ok {
					components.Infrastructure.Kubernetes = kubeBool
				}
			}
			if terraform, exists := infraMap["terraform"]; exists {
				if terraformBool, ok := terraform.(bool); ok {
					components.Infrastructure.Terraform = terraformBool
				}
			}
		}
	}

	return nil
}

// extractVersionConfig extracts version configuration from import data
func (t *ConfigTransformer) extractVersionConfig(data interface{}) (*models.VersionConfig, error) {
	versionsMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("versions must be an object")
	}

	versions := &models.VersionConfig{
		Packages: make(map[string]string),
	}

	if node, exists := versionsMap["node"]; exists {
		if nodeStr, ok := node.(string); ok {
			versions.Node = nodeStr
		}
	}

	if goVersion, exists := versionsMap["go"]; exists {
		if goStr, ok := goVersion.(string); ok {
			versions.Go = goStr
		}
	}

	if packages, exists := versionsMap["packages"]; exists {
		if packagesMap, ok := packages.(map[string]interface{}); ok {
			for pkg, version := range packagesMap {
				if versionStr, ok := version.(string); ok {
					versions.Packages[pkg] = versionStr
				}
			}
		}
	}

	return versions, nil
}

// excludeFields removes specified fields from export data
func (t *ConfigTransformer) excludeFields(data map[string]interface{}, excludeFields []string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		excluded := false
		for _, excludeField := range excludeFields {
			if strings.EqualFold(key, excludeField) {
				excluded = true
				break
			}
		}

		if !excluded {
			result[key] = value
		}
	}

	return result
}

// applyFieldMappings applies field name mappings during import
func (t *ConfigTransformer) applyFieldMappings(config *SavedConfiguration, mappings map[string]string) error {
	// This would implement field name mapping logic
	// For now, it's a placeholder
	return nil
}

// applyImportTransformations applies various transformations during import
func (t *ConfigTransformer) applyImportTransformations(config *SavedConfiguration) error {
	// Normalize project name
	if config.ProjectConfig != nil && config.ProjectConfig.Name != "" {
		config.ProjectConfig.Name = strings.ToLower(strings.ReplaceAll(config.ProjectConfig.Name, " ", "-"))
	}

	// Set default values
	if config.Version == "" {
		config.Version = "1.0.0"
	}

	// Validate and fix component dependencies
	if config.ProjectConfig != nil {
		t.fixComponentDependencies(&config.ProjectConfig.Components)
	}

	return nil
}

// fixComponentDependencies fixes component dependencies
func (t *ConfigTransformer) fixComponentDependencies(components *models.Components) {
	// If Kubernetes is enabled, enable Docker as well
	if components.Infrastructure.Kubernetes && !components.Infrastructure.Docker {
		components.Infrastructure.Docker = true
	}

	// If admin is enabled, enable shared components
	if components.Frontend.NextJS.Admin && !components.Frontend.NextJS.Shared {
		components.Frontend.NextJS.Shared = true
	}
}

// isEmptyComponents checks if components configuration is empty
func (t *ConfigTransformer) isEmptyComponents(components *models.Components) bool {
	v := reflect.ValueOf(components).Elem()
	return t.isZeroValue(v)
}

// isZeroValue checks if a reflect.Value is zero
func (t *ConfigTransformer) isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.String:
		return v.String() == ""
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !t.isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return v.IsZero()
	}
}
