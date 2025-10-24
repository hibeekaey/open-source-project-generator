package fallback

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

//go:embed templates/ios/*
var iosTemplates embed.FS //nolint:unused // Reserved for future template-based generation

// IOSGenerator implements fallback generation for iOS projects
type IOSGenerator struct{}

// NewIOSGenerator creates a new iOS fallback generator
func NewIOSGenerator() *IOSGenerator {
	return &IOSGenerator{}
}

// Generate creates a minimal iOS project structure
func (g *IOSGenerator) Generate(ctx context.Context, spec *models.FallbackSpec) (*models.ComponentResult, error) {
	startTime := time.Now()

	result := &models.ComponentResult{
		Type:     "ios",
		Name:     getIOSComponentName(spec),
		Method:   "fallback",
		ToolUsed: "embedded-templates",
	}

	// Create target directory
	if err := os.MkdirAll(spec.TargetDir, 0750); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to create target directory: %w", err)
		return result, result.Error
	}

	// Get configuration
	appName := getStringConfig(spec.Config, "app_name", "MyApp")
	bundleID := getStringConfig(spec.Config, "bundle_id", "com.example.app")
	organizationName := getStringConfig(spec.Config, "organization", "MyOrganization")

	// Create directory structure
	if err := g.createDirectoryStructure(spec.TargetDir, appName); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to create directory structure: %w", err)
		return result, result.Error
	}

	// Generate files
	if err := g.generateFiles(spec.TargetDir, appName, bundleID, organizationName); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to generate files: %w", err)
		return result, result.Error
	}

	result.Success = true
	result.OutputPath = spec.TargetDir
	result.Duration = time.Since(startTime)
	result.ManualSteps = g.GetRequiredManualSteps("ios")
	result.Warnings = []string{
		"This is a minimal iOS project structure",
		"Xcode is required to build and run the project",
		"Code signing configuration needed",
	}

	return result, nil
}

// SupportsComponent checks if this generator supports the component type
func (g *IOSGenerator) SupportsComponent(componentType string) bool {
	return componentType == "ios"
}

// GetRequiredManualSteps returns manual steps needed after generation
func (g *IOSGenerator) GetRequiredManualSteps(componentType string) []string {
	return []string{
		"Install Xcode from the Mac App Store",
		"Open the .xcodeproj file in Xcode",
		"Configure code signing in Xcode project settings",
		"Select a development team in the Signing & Capabilities tab",
		"Choose a simulator or connect a physical iOS device",
		"Build and run the project using Cmd+R",
	}
}

// createDirectoryStructure creates the iOS project directory structure
func (g *IOSGenerator) createDirectoryStructure(targetDir, appName string) error {
	dirs := []string{
		appName + ".xcodeproj",
		appName,
		appName + "/Assets.xcassets",
		appName + "/Assets.xcassets/AppIcon.appiconset",
		appName + "/Assets.xcassets/AccentColor.colorset",
		appName + "/Preview Content",
		appName + "/Preview Content/Preview Assets.xcassets",
		appName + "Tests",
		appName + "UITests",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(targetDir, dir)
		if err := os.MkdirAll(fullPath, 0750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateFiles generates all necessary iOS project files
func (g *IOSGenerator) generateFiles(targetDir, appName, bundleID, organizationName string) error {
	files := map[string]string{
		appName + ".xcodeproj/project.pbxproj":                             g.generateProjectPbxproj(appName, bundleID, organizationName),
		appName + "/" + appName + "App.swift":                              g.generateAppFile(appName),
		appName + "/ContentView.swift":                                     g.generateContentView(),
		appName + "/Assets.xcassets/Contents.json":                         g.generateAssetsContents(),
		appName + "/Assets.xcassets/AppIcon.appiconset/Contents.json":      g.generateAppIconContents(),
		appName + "/Assets.xcassets/AccentColor.colorset/Contents.json":    g.generateAccentColorContents(),
		appName + "/Preview Content/Preview Assets.xcassets/Contents.json": g.generateAssetsContents(),
		appName + "/Info.plist":                                            g.generateInfoPlist(),
		appName + "Tests/" + appName + "Tests.swift":                       g.generateTests(appName),
		appName + "UITests/" + appName + "UITests.swift":                   g.generateUITests(appName),
		appName + "UITests/" + appName + "UITestsLaunchTests.swift":        g.generateUITestsLaunch(appName),
		".gitignore": g.generateIOSGitignore(),
		"README.md":  g.generateIOSReadme(appName, bundleID),
	}

	for filePath, content := range files {
		fullPath := filepath.Join(targetDir, filePath)
		// Use 0644 for public files (README, .gitignore), 0600 for source code
		perm := os.FileMode(0600)
		if filePath == "README.md" || filePath == ".gitignore" {
			perm = 0644 // #nosec G306 - Public documentation files
		}
		if err := os.WriteFile(fullPath, []byte(content), perm); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}

func (g *IOSGenerator) generateProjectPbxproj(appName, bundleID, organizationName string) string {
	return `// !$*UTF8*$!
{
	archiveVersion = 1;
	classes = {
	};
	objectVersion = 56;
	objects = {
		/* Begin PBXBuildFile section */
		/* End PBXBuildFile section */
		
		/* Begin PBXFileReference section */
		/* End PBXFileReference section */
		
		/* Begin PBXFrameworksBuildPhase section */
		/* End PBXFrameworksBuildPhase section */
		
		/* Begin PBXGroup section */
		/* End PBXGroup section */
		
		/* Begin PBXNativeTarget section */
		/* End PBXNativeTarget section */
		
		/* Begin PBXProject section */
		/* End PBXProject section */
		
		/* Begin PBXResourcesBuildPhase section */
		/* End PBXResourcesBuildPhase section */
		
		/* Begin PBXSourcesBuildPhase section */
		/* End PBXSourcesBuildPhase section */
		
		/* Begin XCBuildConfiguration section */
		/* End XCBuildConfiguration section */
		
		/* Begin XCConfigurationList section */
		/* End XCConfigurationList section */
	};
	rootObject = /* Project object */;
}
`
}

func (g *IOSGenerator) generateAppFile(appName string) string {
	return fmt.Sprintf(`import SwiftUI

@main
struct %sApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}
`, appName)
}

func (g *IOSGenerator) generateContentView() string {
	return `import SwiftUI

struct ContentView: View {
    var body: some View {
        VStack {
            Image(systemName: "globe")
                .imageScale(.large)
                .foregroundStyle(.tint)
            Text("Hello, World!")
        }
        .padding()
    }
}

#Preview {
    ContentView()
}
`
}

func (g *IOSGenerator) generateAssetsContents() string {
	return `{
  "info" : {
    "author" : "xcode",
    "version" : 1
  }
}
`
}

func (g *IOSGenerator) generateAppIconContents() string {
	return `{
  "images" : [
    {
      "idiom" : "universal",
      "platform" : "ios",
      "size" : "1024x1024"
    }
  ],
  "info" : {
    "author" : "xcode",
    "version" : 1
  }
}
`
}

func (g *IOSGenerator) generateAccentColorContents() string {
	return `{
  "colors" : [
    {
      "idiom" : "universal"
    }
  ],
  "info" : {
    "author" : "xcode",
    "version" : 1
  }
}
`
}

func (g *IOSGenerator) generateInfoPlist() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>UIApplicationSceneManifest</key>
	<dict>
		<key>UIApplicationSupportsMultipleScenes</key>
		<true/>
	</dict>
</dict>
</plist>
`
}

func (g *IOSGenerator) generateTests(appName string) string {
	return fmt.Sprintf(`import XCTest
@testable import %s

final class %sTests: XCTestCase {
    override func setUpWithError() throws {
        // Put setup code here. This method is called before the invocation of each test method in the class.
    }

    override func tearDownWithError() throws {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
    }

    func testExample() throws {
        // This is an example of a functional test case.
        // Use XCTAssert and related functions to verify your tests produce the correct results.
        XCTAssertTrue(true)
    }

    func testPerformanceExample() throws {
        // This is an example of a performance test case.
        self.measure {
            // Put the code you want to measure the time of here.
        }
    }
}
`, appName, appName)
}

func (g *IOSGenerator) generateUITests(appName string) string {
	return fmt.Sprintf(`import XCTest

final class %sUITests: XCTestCase {
    override func setUpWithError() throws {
        // Put setup code here. This method is called before the invocation of each test method in the class.
        continueAfterFailure = false
    }

    override func tearDownWithError() throws {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
    }

    func testExample() throws {
        // UI tests must launch the application that they test.
        let app = XCUIApplication()
        app.launch()

        // Use XCTAssert and related functions to verify your tests produce the correct results.
    }

    func testLaunchPerformance() throws {
        if #available(macOS 10.15, iOS 13.0, tvOS 13.0, watchOS 7.0, *) {
            // This measures how long it takes to launch your application.
            measure(metrics: [XCTApplicationLaunchMetric()]) {
                XCUIApplication().launch()
            }
        }
    }
}
`, appName)
}

func (g *IOSGenerator) generateUITestsLaunch(appName string) string {
	return fmt.Sprintf(`import XCTest

final class %sUITestsLaunchTests: XCTestCase {
    override class var runsForEachTargetApplicationUIConfiguration: Bool {
        true
    }

    override func setUpWithError() throws {
        continueAfterFailure = false
    }

    func testLaunch() throws {
        let app = XCUIApplication()
        app.launch()

        // Insert steps here to perform after app launch but before taking a screenshot,
        // such as logging into a test account or navigating somewhere in the app

        let attachment = XCTAttachment(screenshot: app.screenshot())
        attachment.name = "Launch Screen"
        attachment.lifetime = .keepAlways
        add(attachment)
    }
}
`, appName)
}

func (g *IOSGenerator) generateIOSGitignore() string {
	return `# Xcode
#
# gitignore contributors: remember to update Global/Xcode.gitignore, Objective-C.gitignore & Swift.gitignore

## User settings
xcuserdata/

## compatibility with Xcode 8 and earlier (ignoring not required starting Xcode 9)
*.xcscmblueprint
*.xccheckout

## compatibility with Xcode 3 and earlier (ignoring not required starting Xcode 4)
build/
DerivedData/
*.moved-aside
*.pbxuser
!default.pbxuser
*.mode1v3
!default.mode1v3
*.mode2v3
!default.mode2v3
*.perspectivev3
!default.perspectivev3

## Obj-C/Swift specific
*.hmap

## App packaging
*.ipa
*.dSYM.zip
*.dSYM

## Playgrounds
timeline.xctimeline
playground.xcworkspace

# Swift Package Manager
.build/
.swiftpm/

# CocoaPods
Pods/

# Carthage
Carthage/Build/

# Accio dependency management
Dependencies/
.accio/

# fastlane
fastlane/report.xml
fastlane/Preview.html
fastlane/screenshots/**/*.png
fastlane/test_output

# Code Injection
iOSInjectionProject/

# macOS
.DS_Store
`
}

func (g *IOSGenerator) generateIOSReadme(appName, bundleID string) string {
	return fmt.Sprintf(`# %s

This is a minimal iOS project generated using fallback templates.

## Bundle Identifier
%s

## Setup Instructions

### Prerequisites
- macOS (required for iOS development)
- Xcode 26.0 or later
- Apple Developer account (for device deployment)

### Steps to Run

1. **Install Xcode**
   - Download from the Mac App Store
   - Launch Xcode and accept the license agreement
   - Install additional components if prompted

2. **Open the Project**
   - Double-click the `+"`%s.xcodeproj`"+` file
   - Or open Xcode and select "Open a project or file"

3. **Configure Code Signing**
   - Select the project in the navigator
   - Select the target under "Targets"
   - Go to "Signing & Capabilities" tab
   - Select your development team
   - Xcode will automatically manage provisioning profiles

4. **Run the App**
   - Select a simulator from the device menu (e.g., iPhone 15)
   - Or connect a physical iOS device
   - Press Cmd+R or click the Run button

## Project Structure

- `+"`%s/`"+` - Main app source code
  - `+"`%sApp.swift`"+` - App entry point
  - `+"`ContentView.swift`"+` - Main view
  - `+"`Assets.xcassets/`"+` - Images and colors
- `+"`%sTests/`"+` - Unit tests
- `+"`%sUITests/`"+` - UI tests

## Development

### SwiftUI
This project uses SwiftUI for the user interface. SwiftUI is Apple's modern declarative framework for building UIs.

### Adding New Views
1. Create a new Swift file (File > New > File > Swift File)
2. Import SwiftUI
3. Define your view struct conforming to the View protocol

### Adding Assets
- Drag images into Assets.xcassets
- Use Image("imageName") in your SwiftUI views

## Building for Device

1. Connect your iOS device via USB
2. Select your device from the device menu
3. Ensure your device is registered in your Apple Developer account
4. Build and run (Cmd+R)

## Troubleshooting

### Code Signing Issues
- Ensure you're logged into Xcode with your Apple ID (Xcode > Settings > Accounts)
- Select a valid development team in project settings
- Try "Automatically manage signing"

### Simulator Not Working
- Quit and restart Xcode
- Reset simulator: Device > Erase All Content and Settings
- Check that Xcode Command Line Tools are installed

### Build Errors
- Clean build folder: Product > Clean Build Folder (Cmd+Shift+K)
- Delete derived data: ~/Library/Developer/Xcode/DerivedData
- Restart Xcode

## Resources

- [SwiftUI Documentation](https://developer.apple.com/documentation/swiftui/)
- [Swift Programming Language](https://docs.swift.org/swift-book/)
- [Human Interface Guidelines](https://developer.apple.com/design/human-interface-guidelines/)
- [Xcode Help](https://developer.apple.com/documentation/xcode)

## Next Steps

- Customize the UI in ContentView.swift
- Add new views and navigation
- Integrate with backend APIs
- Add app icons and launch screens
- Configure app capabilities (push notifications, etc.)
`, appName, bundleID, appName, appName, appName, appName, appName)
}

func getIOSComponentName(spec *models.FallbackSpec) string {
	if name, ok := spec.Config["name"].(string); ok && name != "" {
		return name
	}
	return "ios-app"
}
