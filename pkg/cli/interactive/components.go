package interactive

import (
	"fmt"
	"strconv"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// configureNextJS prompts for Next.js-specific configuration
func (iw *InteractiveWizard) configureNextJS(config *models.ComponentConfig) error {
	fmt.Println("Configuring Next.js options...")

	// TypeScript
	typescript, err := iw.prompter.Confirm("Use TypeScript?", true)
	if err != nil {
		return err
	}
	config.Config["typescript"] = typescript

	// Tailwind CSS
	tailwind, err := iw.prompter.Confirm("Use Tailwind CSS?", true)
	if err != nil {
		return err
	}
	config.Config["tailwind"] = tailwind

	// App Router
	appRouter, err := iw.prompter.Confirm("Use App Router (recommended)?", true)
	if err != nil {
		return err
	}
	config.Config["app_router"] = appRouter

	// ESLint
	eslint, err := iw.prompter.Confirm("Use ESLint?", true)
	if err != nil {
		return err
	}
	config.Config["eslint"] = eslint

	return nil
}

// configureGoBackend prompts for Go backend-specific configuration
func (iw *InteractiveWizard) configureGoBackend(config *models.ComponentConfig) error {
	fmt.Println("Configuring Go backend options...")

	// Go module path
	defaultModule := "github.com/user/" + config.Name
	module, err := InputWithValidation(
		iw.prompter,
		"Go module path",
		defaultModule,
		ValidateGoModule,
	)
	if err != nil {
		return err
	}
	config.Config["module"] = module

	// Framework selection
	frameworkOptions := []string{
		"gin",
		"echo",
		"fiber",
	}
	framework, err := iw.prompter.Select("Select web framework:", frameworkOptions)
	if err != nil {
		return err
	}
	config.Config["framework"] = framework

	// Port
	port, err := InputWithValidation(
		iw.prompter,
		"Server port",
		"8080",
		ValidatePort,
	)
	if err != nil {
		return err
	}

	// Convert port to integer
	portInt, _ := strconv.Atoi(port)
	config.Config["port"] = portInt

	return nil
}

// configureAndroid prompts for Android-specific configuration
func (iw *InteractiveWizard) configureAndroid(config *models.ComponentConfig) error {
	fmt.Println("Configuring Android options...")

	// Package name
	defaultPackage := "com.example." + config.Name
	packageName, err := InputWithValidation(
		iw.prompter,
		"Android package name",
		defaultPackage,
		ValidatePackageName,
	)
	if err != nil {
		return err
	}
	config.Config["package"] = packageName

	// Minimum SDK
	minSDK, err := InputWithValidation(
		iw.prompter,
		"Minimum SDK version",
		"24",
		ValidateAPILevel,
	)
	if err != nil {
		return err
	}

	// Convert to integer
	minSDKInt, _ := strconv.Atoi(minSDK)
	config.Config["min_sdk"] = minSDKInt

	// Target SDK
	targetSDK, err := InputWithValidation(
		iw.prompter,
		"Target SDK version",
		"34",
		ValidateAPILevel,
	)
	if err != nil {
		return err
	}

	// Convert to integer
	targetSDKInt, _ := strconv.Atoi(targetSDK)
	config.Config["target_sdk"] = targetSDKInt

	// Language
	languageOptions := []string{
		"kotlin",
		"java",
	}
	language, err := iw.prompter.Select("Select programming language:", languageOptions)
	if err != nil {
		return err
	}
	config.Config["language"] = language

	return nil
}

// configureIOS prompts for iOS-specific configuration
func (iw *InteractiveWizard) configureIOS(config *models.ComponentConfig) error {
	fmt.Println("Configuring iOS options...")

	// Bundle ID
	defaultBundleID := "com.example." + config.Name
	bundleID, err := InputWithValidation(
		iw.prompter,
		"Bundle identifier",
		defaultBundleID,
		ValidateBundleID,
	)
	if err != nil {
		return err
	}
	config.Config["bundle_id"] = bundleID

	// Deployment target
	deploymentTarget, err := InputWithValidation(
		iw.prompter,
		"Minimum iOS deployment target",
		"15.0",
		ValidateIOSVersion,
	)
	if err != nil {
		return err
	}
	config.Config["deployment_target"] = deploymentTarget

	// Language
	languageOptions := []string{
		"swift",
		"objective-c",
	}
	language, err := iw.prompter.Select("Select programming language:", languageOptions)
	if err != nil {
		return err
	}
	config.Config["language"] = language

	return nil
}
