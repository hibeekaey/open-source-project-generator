package fallback

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

//go:embed templates/android/*
var androidTemplates embed.FS //nolint:unused // Reserved for future template-based generation

// AndroidGenerator implements fallback generation for Android projects
type AndroidGenerator struct{}

// NewAndroidGenerator creates a new Android fallback generator
func NewAndroidGenerator() *AndroidGenerator {
	return &AndroidGenerator{}
}

// Generate creates a minimal Android project structure
func (g *AndroidGenerator) Generate(ctx context.Context, spec *models.FallbackSpec) (*models.ComponentResult, error) {
	startTime := time.Now()

	result := &models.ComponentResult{
		Type:     "android",
		Name:     getComponentName(spec),
		Method:   "fallback",
		ToolUsed: "embedded-templates",
	}

	// Create target directory
	if err := os.MkdirAll(spec.TargetDir, 0755); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to create target directory: %w", err)
		return result, result.Error
	}

	// Get package name from config or use default
	packageName := getStringConfig(spec.Config, "package", "com.example.app")
	appName := getStringConfig(spec.Config, "app_name", "MyApp")

	// Create directory structure
	if err := g.createDirectoryStructure(spec.TargetDir, packageName); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to create directory structure: %w", err)
		return result, result.Error
	}

	// Generate files
	if err := g.generateFiles(spec.TargetDir, packageName, appName); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to generate files: %w", err)
		return result, result.Error
	}

	result.Success = true
	result.OutputPath = spec.TargetDir
	result.Duration = time.Since(startTime)
	result.ManualSteps = g.GetRequiredManualSteps("android")
	result.Warnings = []string{
		"This is a minimal Android project structure",
		"Android Studio and Gradle setup required",
		"Dependencies need to be synced manually",
	}

	return result, nil
}

// SupportsComponent checks if this generator supports the component type
func (g *AndroidGenerator) SupportsComponent(componentType string) bool {
	return componentType == "android"
}

// GetRequiredManualSteps returns manual steps needed after generation
func (g *AndroidGenerator) GetRequiredManualSteps(componentType string) []string {
	return []string{
		"Install Android Studio from https://developer.android.com/studio",
		"Open the project in Android Studio",
		"Wait for Gradle sync to complete",
		"Configure Android SDK if not already installed",
		"Update dependencies in build.gradle files as needed",
		"Run the app on an emulator or physical device",
	}
}

// createDirectoryStructure creates the Android project directory structure
func (g *AndroidGenerator) createDirectoryStructure(targetDir, packageName string) error {
	packagePath := strings.ReplaceAll(packageName, ".", "/")

	dirs := []string{
		"app/src/main/java/" + packagePath,
		"app/src/main/res/layout",
		"app/src/main/res/values",
		"app/src/main/res/drawable",
		"app/src/main/res/mipmap-hdpi",
		"app/src/main/res/mipmap-mdpi",
		"app/src/main/res/mipmap-xhdpi",
		"app/src/main/res/mipmap-xxhdpi",
		"app/src/main/res/mipmap-xxxhdpi",
		"app/src/androidTest/java/" + packagePath,
		"app/src/test/java/" + packagePath,
		"gradle/wrapper",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(targetDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateFiles generates all necessary Android project files
func (g *AndroidGenerator) generateFiles(targetDir, packageName, appName string) error {
	packagePath := strings.ReplaceAll(packageName, ".", "/")

	files := map[string]string{
		"settings.gradle":                  g.generateSettingsGradle(appName),
		"build.gradle":                     g.generateRootBuildGradle(),
		"gradle.properties":                g.generateGradleProperties(),
		"app/build.gradle":                 g.generateAppBuildGradle(packageName),
		"app/proguard-rules.pro":           g.generateProguardRules(),
		"app/src/main/AndroidManifest.xml": g.generateAndroidManifest(packageName, appName),
		"app/src/main/java/" + packagePath + "/MainActivity.kt": g.generateMainActivity(packageName),
		"app/src/main/res/layout/activity_main.xml":             g.generateMainLayout(),
		"app/src/main/res/values/strings.xml":                   g.generateStrings(appName),
		"app/src/main/res/values/colors.xml":                    g.generateColors(),
		"app/src/main/res/values/themes.xml":                    g.generateThemes(),
		"gradle/wrapper/gradle-wrapper.properties":              g.generateGradleWrapperProperties(),
		".gitignore": g.generateGitignore(),
		"README.md":  g.generateReadme(appName, packageName),
	}

	for filePath, content := range files {
		fullPath := filepath.Join(targetDir, filePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}

func (g *AndroidGenerator) generateSettingsGradle(appName string) string {
	return fmt.Sprintf(`pluginManagement {
    repositories {
        google()
        mavenCentral()
        gradlePluginPortal()
    }
}

dependencyResolutionManagement {
    repositoriesMode.set(RepositoriesMode.FAIL_ON_PROJECT_REPOS)
    repositories {
        google()
        mavenCentral()
    }
}

rootProject.name = "%s"
include ':app'
`, appName)
}

func (g *AndroidGenerator) generateRootBuildGradle() string {
	return `// Top-level build file
plugins {
    id 'com.android.application' version '8.7.3' apply false
    id 'com.android.library' version '8.7.3' apply false
    id 'org.jetbrains.kotlin.android' version '2.2.21' apply false
}

task clean(type: Delete) {
    delete rootProject.buildDir
}
`
}

func (g *AndroidGenerator) generateGradleProperties() string {
	return `# Project-wide Gradle settings
org.gradle.jvmargs=-Xmx2048m -Dfile.encoding=UTF-8
android.useAndroidX=true
android.enableJetifier=true
kotlin.code.style=official
`
}

func (g *AndroidGenerator) generateAppBuildGradle(packageName string) string {
	return fmt.Sprintf(`plugins {
    id 'com.android.application'
    id 'org.jetbrains.kotlin.android'
}

android {
    namespace '%s'
    compileSdk 36

    defaultConfig {
        applicationId "%s"
        minSdk 24
        targetSdk 36
        versionCode 1
        versionName "1.0"

        testInstrumentationRunner "androidx.test.runner.AndroidJUnitRunner"
    }

    buildTypes {
        release {
            minifyEnabled false
            proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
        }
    }

    compileOptions {
        sourceCompatibility JavaVersion.VERSION_17
        targetCompatibility JavaVersion.VERSION_17
    }

    kotlinOptions {
        jvmTarget = '17'
    }

    buildFeatures {
        viewBinding true
    }
}

dependencies {
    implementation 'androidx.core:core-ktx:1.17.0'
    implementation 'androidx.appcompat:appcompat:1.7.1'
    implementation 'com.google.android.material:material:1.13.0'
    implementation 'androidx.constraintlayout:constraintlayout:2.2.1'
    
    testImplementation 'junit:junit:4.13.2'
    androidTestImplementation 'androidx.test.ext:junit:1.2.1'
    androidTestImplementation 'androidx.test.espresso:espresso-core:3.6.1'
}
`, packageName, packageName)
}

func (g *AndroidGenerator) generateProguardRules() string {
	return `# Add project specific ProGuard rules here.
# You can control the set of applied configuration files using the
# proguardFiles setting in build.gradle.
`
}

func (g *AndroidGenerator) generateAndroidManifest(packageName, appName string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android">

    <application
        android:allowBackup="true"
        android:icon="@mipmap/ic_launcher"
        android:label="@string/app_name"
        android:roundIcon="@mipmap/ic_launcher_round"
        android:supportsRtl="true"
        android:theme="@style/Theme.%s">
        <activity
            android:name=".MainActivity"
            android:exported="true">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
    </application>

</manifest>
`, strings.ReplaceAll(appName, " ", ""))
}

func (g *AndroidGenerator) generateMainActivity(packageName string) string {
	return fmt.Sprintf(`package %s

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity

class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
    }
}
`, packageName)
}

func (g *AndroidGenerator) generateMainLayout() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<androidx.constraintlayout.widget.ConstraintLayout 
    xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    xmlns:tools="http://schemas.android.com/tools"
    android:layout_width="match_parent"
    android:layout_height="match_parent"
    tools:context=".MainActivity">

    <TextView
        android:layout_width="wrap_content"
        android:layout_height="wrap_content"
        android:text="@string/hello_world"
        app:layout_constraintBottom_toBottomOf="parent"
        app:layout_constraintEnd_toEndOf="parent"
        app:layout_constraintStart_toStartOf="parent"
        app:layout_constraintTop_toTopOf="parent" />

</androidx.constraintlayout.widget.ConstraintLayout>
`
}

func (g *AndroidGenerator) generateStrings(appName string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="app_name">%s</string>
    <string name="hello_world">Hello World!</string>
</resources>
`, appName)
}

func (g *AndroidGenerator) generateColors() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<resources>
    <color name="purple_200">#FFBB86FC</color>
    <color name="purple_500">#FF6200EE</color>
    <color name="purple_700">#FF3700B3</color>
    <color name="teal_200">#FF03DAC5</color>
    <color name="teal_700">#FF018786</color>
    <color name="black">#FF000000</color>
    <color name="white">#FFFFFFFF</color>
</resources>
`
}

func (g *AndroidGenerator) generateThemes() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<resources>
    <style name="Theme.MyApp" parent="Theme.MaterialComponents.DayNight.DarkActionBar">
        <item name="colorPrimary">@color/purple_500</item>
        <item name="colorPrimaryVariant">@color/purple_700</item>
        <item name="colorOnPrimary">@color/white</item>
        <item name="colorSecondary">@color/teal_200</item>
        <item name="colorSecondaryVariant">@color/teal_700</item>
        <item name="colorOnSecondary">@color/black</item>
    </style>
</resources>
`
}

func (g *AndroidGenerator) generateGradleWrapperProperties() string {
	return `distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
distributionUrl=https\://services.gradle.org/distributions/gradle-9.1.0-bin.zip
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
`
}

func (g *AndroidGenerator) generateGitignore() string {
	return `# Built application files
*.apk
*.aar
*.ap_
*.aab

# Files for the ART/Dalvik VM
*.dex

# Java class files
*.class

# Generated files
bin/
gen/
out/
build/

# Gradle files
.gradle/
.gradle

# Local configuration file
local.properties

# Android Studio
.idea/
*.iml
.DS_Store
captures/
.externalNativeBuild
.cxx

# Keystore files
*.jks
*.keystore
`
}

func (g *AndroidGenerator) generateReadme(appName, packageName string) string {
	return fmt.Sprintf(`# %s

This is a minimal Android project generated using fallback templates.

## Package Name
%s

## Setup Instructions

### Prerequisites
- Android Studio (latest version recommended)
- Android SDK (API 24 or higher)
- Java Development Kit (JDK) 8 or higher

### Steps to Run

1. **Install Android Studio**
   - Download from: https://developer.android.com/studio
   - Follow the installation wizard

2. **Open the Project**
   - Launch Android Studio
   - Select "Open an Existing Project"
   - Navigate to this directory and select it

3. **Sync Dependencies**
   - Android Studio will automatically start syncing Gradle
   - Wait for the sync to complete (this may take a few minutes)
   - If prompted, install any missing SDK components

4. **Run the App**
   - Connect an Android device via USB (with USB debugging enabled)
   - Or create an Android Virtual Device (AVD) in Android Studio
   - Click the "Run" button (green play icon) in the toolbar
   - Select your device/emulator

## Project Structure

- `+"`app/src/main/java/`"+` - Kotlin source files
- `+"`app/src/main/res/`"+` - Resources (layouts, strings, etc.)
- `+"`app/src/main/AndroidManifest.xml`"+` - App manifest
- `+"`app/build.gradle`"+` - App-level build configuration
- `+"`build.gradle`"+` - Project-level build configuration

## Next Steps

- Customize the app name in `+"`app/src/main/res/values/strings.xml`"+`
- Modify the UI in `+"`app/src/main/res/layout/activity_main.xml`"+`
- Add more activities and features as needed
- Update dependencies in `+"`app/build.gradle`"+`

## Troubleshooting

- If Gradle sync fails, check your internet connection
- Ensure Android SDK is properly installed
- Try "File > Invalidate Caches / Restart" in Android Studio
- Check the Gradle console for specific error messages

## Resources

- [Android Developer Guide](https://developer.android.com/guide)
- [Kotlin Documentation](https://kotlinlang.org/docs/home.html)
- [Material Design](https://material.io/develop/android)
`, appName, packageName)
}

// Helper functions

func getComponentName(spec *models.FallbackSpec) string {
	if name, ok := spec.Config["name"].(string); ok && name != "" {
		return name
	}
	return "android-app"
}

func getStringConfig(config map[string]interface{}, key, defaultValue string) string {
	if val, ok := config[key].(string); ok && val != "" {
		return val
	}
	return defaultValue
}
