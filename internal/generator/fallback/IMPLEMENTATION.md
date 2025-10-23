# Fallback Generator Implementation

This document describes the implementation of the fallback generator system for the CLI Architecture Modernization project.

## Overview

The fallback generator system provides a way to generate minimal project structures when external bootstrap tools (like gradle, xcodebuild) are not available. This ensures that users can still create projects even in environments where these tools cannot be installed.

## Components Implemented

### 1. Generator Interface and Registry (`generator.go`)

**Generator Interface**
- `Generate(ctx, spec)` - Creates a project component using custom templates
- `SupportsComponent(componentType)` - Checks if generator handles the component type
- `GetRequiredManualSteps(componentType)` - Returns post-generation manual steps

**Registry**
- Manages fallback generators for different component types
- Provides registration and lookup functionality
- Includes a `DefaultRegistry()` with Android and iOS generators pre-registered

### 2. Android Fallback Generator (`android.go`)

Generates a minimal Android project with:

**Project Structure**
- Standard Android app directory layout
- Kotlin source files with proper package structure
- Resource directories (layout, values, drawable, mipmap)
- Test directories (androidTest, test)
- Gradle wrapper configuration

**Generated Files**
- `settings.gradle` - Project settings
- `build.gradle` (root and app) - Build configuration
- `gradle.properties` - Gradle properties
- `AndroidManifest.xml` - App manifest
- `MainActivity.kt` - Main activity in Kotlin
- `activity_main.xml` - Main layout
- `strings.xml`, `colors.xml`, `themes.xml` - Resources
- `gradle-wrapper.properties` - Gradle wrapper config
- `.gitignore` - Git ignore rules
- `README.md` - Comprehensive setup instructions

**Configuration Options**
- `package` - Java package name (default: com.example.app)
- `app_name` - Application name (default: MyApp)

**Manual Steps Required**
1. Install Android Studio
2. Open project in Android Studio
3. Wait for Gradle sync
4. Configure Android SDK
5. Update dependencies as needed
6. Run on emulator or device

### 3. iOS Fallback Generator (`ios.go`)

Generates a minimal iOS project with:

**Project Structure**
- Xcode project structure (.xcodeproj)
- SwiftUI app files
- Asset catalogs (AppIcon, AccentColor)
- Preview content
- Test directories (Tests, UITests)

**Generated Files**
- `project.pbxproj` - Xcode project file
- `{AppName}App.swift` - App entry point
- `ContentView.swift` - Main SwiftUI view
- Asset catalog JSON files
- `Info.plist` - App information
- Test files (unit and UI tests)
- `.gitignore` - Git ignore rules
- `README.md` - Comprehensive setup instructions

**Configuration Options**
- `app_name` - Application name (default: MyApp)
- `bundle_id` - Bundle identifier (default: com.example.app)
- `organization` - Organization name (default: MyOrganization)

**Manual Steps Required**
1. Install Xcode from Mac App Store
2. Open .xcodeproj file in Xcode
3. Configure code signing
4. Select development team
5. Choose simulator or device
6. Build and run (Cmd+R)

### 4. Embedded Templates (`templates/`)

Created template directories with placeholder files to satisfy Go's embed directive:
- `templates/android/` - Android templates
- `templates/ios/` - iOS templates
- `templates/README.md` - Documentation

Note: Current implementation generates files programmatically rather than using template files. The embed directive is prepared for future template-based generation if needed.

## Testing

Comprehensive test suite (`generator_test.go`) covering:

**Registry Tests**
- Registry creation and initialization
- Generator registration and retrieval
- Default registry configuration
- Error handling for non-existent generators

**Android Generator Tests**
- Component type support
- Manual steps generation
- Full project generation with file verification

**iOS Generator Tests**
- Component type support
- Manual steps generation
- Full project generation with file verification

All tests pass successfully.

## Usage Example

```go
// Create a registry with default generators
registry := fallback.DefaultRegistry()

// Get the Android generator
androidGen, err := registry.Get("android")
if err != nil {
    log.Fatal(err)
}

// Create a fallback spec
spec := &models.FallbackSpec{
    ComponentType: "android",
    TargetDir:     "/path/to/output",
    Config: map[string]interface{}{
        "package":  "com.mycompany.app",
        "app_name": "MyAwesomeApp",
    },
}

// Generate the project
result, err := androidGen.Generate(context.Background(), spec)
if err != nil {
    log.Fatal(err)
}

// Check result
if result.Success {
    fmt.Printf("Generated %s project at %s\n", result.Type, result.OutputPath)
    fmt.Println("Manual steps required:")
    for _, step := range result.ManualSteps {
        fmt.Printf("  - %s\n", step)
    }
}
```

## Integration Points

The fallback generators integrate with:

1. **Tool Discovery** - Used when bootstrap tools are not available
2. **Project Coordinator** - Called as part of the generation workflow
3. **Structure Mapper** - Output is mapped to target directory structure
4. **Error Handling** - Provides clear error messages and recovery options

## Future Enhancements

Potential improvements:
1. Add more platform support (Flutter, React Native)
2. Template-based generation using embedded files
3. Customizable templates via configuration
4. Incremental generation (add components to existing projects)
5. Version-specific templates for different SDK versions

## Requirements Satisfied

This implementation satisfies the following requirements from the specification:

- **Requirement 6.1**: Custom generation logic when no Bootstrap Tool is available
- **Requirement 6.2**: Registry mapping component types to fallback generators
- **Requirement 6.3**: Minimal project structures with manual setup instructions

## Files Created

- `internal/generator/fallback/generator.go` - Interface and registry
- `internal/generator/fallback/android.go` - Android generator
- `internal/generator/fallback/ios.go` - iOS generator
- `internal/generator/fallback/generator_test.go` - Test suite
- `internal/generator/fallback/templates/` - Template directories
- `internal/generator/fallback/IMPLEMENTATION.md` - This document
