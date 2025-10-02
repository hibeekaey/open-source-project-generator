package components

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

// MobileGenerator handles mobile component file generation
type MobileGenerator struct {
	fsOps FileSystemOperations
}

// NewMobileGenerator creates a new mobile generator
func NewMobileGenerator(fsOps FileSystemOperations) *MobileGenerator {
	return &MobileGenerator{
		fsOps: fsOps,
	}
}

// GenerateFiles creates mobile component files based on configuration
func (mg *MobileGenerator) GenerateFiles(projectPath string, config *models.ProjectConfig) error {
	if config == nil {
		return fmt.Errorf("project config cannot be nil")
	}

	// Generate Android files if selected
	if config.Components.Mobile.Android {
		if err := mg.generateAndroidFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate Android files: %w", err)
		}
	}

	// Generate iOS files if selected
	if config.Components.Mobile.IOS {
		if err := mg.generateIOSFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate iOS files: %w", err)
		}
	}

	// Generate shared mobile files if selected
	if config.Components.Mobile.Shared {
		if err := mg.generateSharedFiles(projectPath, config); err != nil {
			return fmt.Errorf("failed to generate shared mobile files: %w", err)
		}
	}

	return nil
}

// generateAndroidFiles creates Android application files
func (mg *MobileGenerator) generateAndroidFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate build.gradle (app level)
	buildGradleContent := mg.generateAndroidBuildGradle(config)
	buildGradlePath := filepath.Join(projectPath, "Mobile/Android/app/build.gradle")
	if err := mg.fsOps.WriteFile(buildGradlePath, []byte(buildGradleContent), 0644); err != nil {
		return fmt.Errorf("failed to create Android app/build.gradle: %w", err)
	}

	// Generate build.gradle (project level)
	projectBuildGradleContent := mg.generateAndroidProjectBuildGradle(config)
	projectBuildGradlePath := filepath.Join(projectPath, "Mobile/Android/build.gradle")
	if err := mg.fsOps.WriteFile(projectBuildGradlePath, []byte(projectBuildGradleContent), 0644); err != nil {
		return fmt.Errorf("failed to create Android build.gradle: %w", err)
	}

	// Generate AndroidManifest.xml
	manifestContent := mg.generateAndroidManifest(config)
	manifestPath := filepath.Join(projectPath, "Mobile/Android/app/src/main/AndroidManifest.xml")
	if err := mg.fsOps.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		return fmt.Errorf("failed to create AndroidManifest.xml: %w", err)
	}

	// Generate MainActivity.kt
	mainActivityContent := mg.generateMainActivity(config)
	mainActivityPath := filepath.Join(projectPath, fmt.Sprintf("Mobile/Android/app/src/main/java/%s/MainActivity.kt", mg.getAndroidPackagePath(config)))
	if err := mg.fsOps.WriteFile(mainActivityPath, []byte(mainActivityContent), 0644); err != nil {
		return fmt.Errorf("failed to create MainActivity.kt: %w", err)
	}

	// Generate strings.xml
	stringsContent := mg.generateAndroidStrings(config)
	stringsPath := filepath.Join(projectPath, "Mobile/Android/app/src/main/res/values/strings.xml")
	if err := mg.fsOps.WriteFile(stringsPath, []byte(stringsContent), 0644); err != nil {
		return fmt.Errorf("failed to create strings.xml: %w", err)
	}

	// Generate colors.xml
	colorsContent := mg.generateAndroidColors(config)
	colorsPath := filepath.Join(projectPath, "Mobile/Android/app/src/main/res/values/colors.xml")
	if err := mg.fsOps.WriteFile(colorsPath, []byte(colorsContent), 0644); err != nil {
		return fmt.Errorf("failed to create colors.xml: %w", err)
	}

	return nil
}

// generateIOSFiles creates iOS application files
func (mg *MobileGenerator) generateIOSFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate Package.swift
	packageSwiftContent := mg.generateIOSPackageSwift(config)
	packageSwiftPath := filepath.Join(projectPath, "Mobile/iOS/Package.swift")
	if err := mg.fsOps.WriteFile(packageSwiftPath, []byte(packageSwiftContent), 0644); err != nil {
		return fmt.Errorf("failed to create iOS Package.swift: %w", err)
	}

	// Generate ContentView.swift
	contentViewContent := mg.generateIOSContentView(config)
	contentViewPath := filepath.Join(projectPath, fmt.Sprintf("Mobile/iOS/Sources/%s/ContentView.swift", config.Name))
	if err := mg.fsOps.WriteFile(contentViewPath, []byte(contentViewContent), 0644); err != nil {
		return fmt.Errorf("failed to create ContentView.swift: %w", err)
	}

	// Generate App.swift
	appSwiftContent := mg.generateIOSApp(config)
	appSwiftPath := filepath.Join(projectPath, fmt.Sprintf("Mobile/iOS/Sources/%s/%sApp.swift", config.Name, config.Name))
	if err := mg.fsOps.WriteFile(appSwiftPath, []byte(appSwiftContent), 0644); err != nil {
		return fmt.Errorf("failed to create %sApp.swift: %w", config.Name, err)
	}

	// Generate Info.plist
	infoPlistContent := mg.generateIOSInfoPlist(config)
	infoPlistPath := filepath.Join(projectPath, "Mobile/iOS/Resources/Info.plist")
	if err := mg.fsOps.WriteFile(infoPlistPath, []byte(infoPlistContent), 0644); err != nil {
		return fmt.Errorf("failed to create Info.plist: %w", err)
	}

	return nil
}

// generateSharedFiles creates shared mobile files
func (mg *MobileGenerator) generateSharedFiles(projectPath string, config *models.ProjectConfig) error {
	// Generate API client interface
	apiClientContent := mg.generateSharedAPIClient(config)
	apiClientPath := filepath.Join(projectPath, "Mobile/Shared/api/client.ts")
	if err := mg.fsOps.WriteFile(apiClientPath, []byte(apiClientContent), 0644); err != nil {
		return fmt.Errorf("failed to create shared API client: %w", err)
	}

	// Generate shared types
	typesContent := mg.generateSharedTypes(config)
	typesPath := filepath.Join(projectPath, "Mobile/Shared/types/index.ts")
	if err := mg.fsOps.WriteFile(typesPath, []byte(typesContent), 0644); err != nil {
		return fmt.Errorf("failed to create shared types: %w", err)
	}

	// Generate shared constants
	constantsContent := mg.generateSharedConstants(config)
	constantsPath := filepath.Join(projectPath, "Mobile/Shared/constants/index.ts")
	if err := mg.fsOps.WriteFile(constantsPath, []byte(constantsContent), 0644); err != nil {
		return fmt.Errorf("failed to create shared constants: %w", err)
	}

	return nil
}

// generateAndroidBuildGradle generates Android app-level build.gradle content
func (mg *MobileGenerator) generateAndroidBuildGradle(config *models.ProjectConfig) string {
	packageName := mg.getAndroidPackageName(config)

	return fmt.Sprintf(`plugins {
    id 'com.android.application'
    id 'org.jetbrains.kotlin.android'
}

android {
    namespace '%s'
    compileSdk 34

    defaultConfig {
        applicationId "%s"
        minSdk 24
        targetSdk 34
        versionCode 1
        versionName "1.0"

        testInstrumentationRunner "androidx.test.runner.AndroidJUnitRunner"
        vectorDrawables {
            useSupportLibrary true
        }
    }

    buildTypes {
        release {
            minifyEnabled false
            proguardFiles getDefaultProguardFile('proguard-android-optimize.txt'), 'proguard-rules.pro'
        }
    }
    
    compileOptions {
        sourceCompatibility JavaVersion.VERSION_1_8
        targetCompatibility JavaVersion.VERSION_1_8
    }
    
    kotlinOptions {
        jvmTarget = '1.8'
    }
    
    buildFeatures {
        compose true
    }
    
    composeOptions {
        kotlinCompilerExtensionVersion '1.5.4'
    }
    
    packaging {
        resources {
            excludes += '/META-INF/{AL2.0,LGPL2.1}'
        }
    }
}

dependencies {
    implementation 'androidx.core:core-ktx:1.12.0'
    implementation 'androidx.lifecycle:lifecycle-runtime-ktx:2.7.0'
    implementation 'androidx.activity:activity-compose:1.8.2'
    implementation platform('androidx.compose:compose-bom:2023.10.01')
    implementation 'androidx.compose.ui:ui'
    implementation 'androidx.compose.ui:ui-graphics'
    implementation 'androidx.compose.ui:ui-tooling-preview'
    implementation 'androidx.compose.material3:material3'
    
    // Navigation
    implementation 'androidx.navigation:navigation-compose:2.7.5'
    
    // Networking
    implementation 'com.squareup.retrofit2:retrofit:2.9.0'
    implementation 'com.squareup.retrofit2:converter-gson:2.9.0'
    implementation 'com.squareup.okhttp3:logging-interceptor:4.12.0'
    
    // Testing
    testImplementation 'junit:junit:4.13.2'
    androidTestImplementation 'androidx.test.ext:junit:1.1.5'
    androidTestImplementation 'androidx.test.espresso:espresso-core:3.5.1'
    androidTestImplementation platform('androidx.compose:compose-bom:2023.10.01')
    androidTestImplementation 'androidx.compose.ui:ui-test-junit4'
    debugImplementation 'androidx.compose.ui:ui-tooling'
    debugImplementation 'androidx.compose.ui:ui-test-manifest'
}`, packageName, packageName)
}

// generateAndroidProjectBuildGradle generates Android project-level build.gradle content
func (mg *MobileGenerator) generateAndroidProjectBuildGradle(config *models.ProjectConfig) string {
	return `// Top-level build file where you can add configuration options common to all sub-projects/modules.
plugins {
    id 'com.android.application' version '8.1.4' apply false
    id 'org.jetbrains.kotlin.android' version '1.9.10' apply false
}`
}

// generateAndroidManifest generates AndroidManifest.xml content
func (mg *MobileGenerator) generateAndroidManifest(config *models.ProjectConfig) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:tools="http://schemas.android.com/tools">

    <uses-permission android:name="android.permission.INTERNET" />

    <application
        android:allowBackup="true"
        android:dataExtractionRules="@xml/data_extraction_rules"
        android:fullBackupContent="@xml/backup_rules"
        android:icon="@mipmap/ic_launcher"
        android:label="@string/app_name"
        android:roundIcon="@mipmap/ic_launcher_round"
        android:supportsRtl="true"
        android:theme="@style/Theme.%s"
        tools:targetApi="31">
        
        <activity
            android:name=".MainActivity"
            android:exported="true"
            android:label="@string/app_name"
            android:theme="@style/Theme.%s">
            <intent-filter>
                <action android:name="android.intent.action.MAIN" />
                <category android:name="android.intent.category.LAUNCHER" />
            </intent-filter>
        </activity>
    </application>

</manifest>`, config.Name, config.Name)
}

// generateMainActivity generates MainActivity.kt content
func (mg *MobileGenerator) generateMainActivity(config *models.ProjectConfig) string {
	packageName := mg.getAndroidPackageName(config)

	return fmt.Sprintf(`package %s

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import %s.ui.theme.%sTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            %sTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    Greeting("Android")
                }
            }
        }
    }
}

@Composable
fun Greeting(name: String, modifier: Modifier = Modifier) {
    Text(
        text = "Hello $name from %s!",
        modifier = modifier
    )
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    %sTheme {
        Greeting("Android")
    }
}`, packageName, packageName, config.Name, config.Name, config.Name, config.Name)
}

// generateAndroidStrings generates strings.xml content
func (mg *MobileGenerator) generateAndroidStrings(config *models.ProjectConfig) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<resources>
    <string name="app_name">%s</string>
    <string name="welcome_message">Welcome to %s</string>
    <string name="loading">Loading...</string>
    <string name="error_network">Network error occurred</string>
    <string name="retry">Retry</string>
</resources>`, config.Name, config.Name)
}

// generateAndroidColors generates colors.xml content
func (mg *MobileGenerator) generateAndroidColors(config *models.ProjectConfig) string {
	return `<?xml version="1.0" encoding="utf-8"?>
<resources>
    <color name="purple_200">#FFBB86FC</color>
    <color name="purple_500">#FF6200EE</color>
    <color name="purple_700">#FF3700B3</color>
    <color name="teal_200">#FF03DAC5</color>
    <color name="teal_700">#FF018786</color>
    <color name="black">#FF000000</color>
    <color name="white">#FFFFFFFF</color>
    <color name="primary">#FF6200EE</color>
    <color name="primary_dark">#FF3700B3</color>
    <color name="accent">#FF03DAC5</color>
</resources>`
}

// generateIOSPackageSwift generates Package.swift content for iOS
func (mg *MobileGenerator) generateIOSPackageSwift(config *models.ProjectConfig) string {
	return fmt.Sprintf(`// swift-tools-version: 5.9
// %s iOS App
// Generated by Open Source Project Generator

import PackageDescription

let package = Package(
    name: "%s",
    platforms: [
        .iOS(.v15)
    ],
    products: [
        .library(
            name: "%s",
            targets: ["%s"]
        ),
    ],
    dependencies: [
        // Add your dependencies here
    ],
    targets: [
        .target(
            name: "%s",
            dependencies: []
        ),
        .testTarget(
            name: "%sTests",
            dependencies: ["%s"]
        ),
    ]
)`, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name)
}

// generateIOSContentView generates ContentView.swift content
func (mg *MobileGenerator) generateIOSContentView(config *models.ProjectConfig) string {
	return fmt.Sprintf(`import SwiftUI

struct ContentView: View {
    var body: some View {
        VStack {
            Image(systemName: "globe")
                .imageScale(.large)
                .foregroundStyle(.tint)
            Text("Hello, %s!")
        }
        .padding()
    }
}

#Preview {
    ContentView()
}`, config.Name)
}

// generateIOSApp generates App.swift content
func (mg *MobileGenerator) generateIOSApp(config *models.ProjectConfig) string {
	return fmt.Sprintf(`import SwiftUI

@main
struct %sApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}`, config.Name)
}

// generateIOSInfoPlist generates Info.plist content
func (mg *MobileGenerator) generateIOSInfoPlist(config *models.ProjectConfig) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>$(DEVELOPMENT_LANGUAGE)</string>
    <key>CFBundleDisplayName</key>
    <string>%s</string>
    <key>CFBundleExecutable</key>
    <string>$(EXECUTABLE_NAME)</string>
    <key>CFBundleIdentifier</key>
    <string>$(PRODUCT_BUNDLE_IDENTIFIER)</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>$(PRODUCT_NAME)</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSRequiresIPhoneOS</key>
    <true/>
    <key>UIApplicationSceneManifest</key>
    <dict>
        <key>UIApplicationSupportsMultipleScenes</key>
        <false/>
        <key>UISceneConfigurations</key>
        <dict>
            <key>UIWindowSceneSessionRoleApplication</key>
            <array>
                <dict>
                    <key>UISceneConfigurationName</key>
                    <string>Default Configuration</string>
                    <key>UISceneDelegateClassName</key>
                    <string>$(PRODUCT_MODULE_NAME).SceneDelegate</string>
                    <key>UISceneStoryboardFile</key>
                    <string>Main</string>
                </dict>
            </array>
        </dict>
    </dict>
    <key>UIRequiredDeviceCapabilities</key>
    <array>
        <string>armv7</string>
    </array>
    <key>UISupportedInterfaceOrientations</key>
    <array>
        <string>UIInterfaceOrientationPortrait</string>
        <string>UIInterfaceOrientationLandscapeLeft</string>
        <string>UIInterfaceOrientationLandscapeRight</string>
    </array>
    <key>UISupportedInterfaceOrientations~ipad</key>
    <array>
        <string>UIInterfaceOrientationPortrait</string>
        <string>UIInterfaceOrientationPortraitUpsideDown</string>
        <string>UIInterfaceOrientationLandscapeLeft</string>
        <string>UIInterfaceOrientationLandscapeRight</string>
    </array>
</dict>
</plist>`, config.Name)
}

// generateSharedAPIClient generates shared API client content
func (mg *MobileGenerator) generateSharedAPIClient(config *models.ProjectConfig) string {
	return fmt.Sprintf(`// %s Shared API Client
// Generated by Open Source Project Generator

export interface APIResponse<T> {
  data: T;
  message?: string;
  status: number;
}

export interface APIError {
  message: string;
  status: number;
  code?: string;
}

export class APIClient {
  private baseURL: string;
  private timeout: number;

  constructor(baseURL: string = 'http://localhost:8080/api/v1', timeout: number = 10000) {
    this.baseURL = baseURL;
    this.timeout = timeout;
  }

  async get<T>(endpoint: string): Promise<APIResponse<T>> {
    return this.request<T>('GET', endpoint);
  }

  async post<T>(endpoint: string, data?: any): Promise<APIResponse<T>> {
    return this.request<T>('POST', endpoint, data);
  }

  async put<T>(endpoint: string, data?: any): Promise<APIResponse<T>> {
    return this.request<T>('PUT', endpoint, data);
  }

  async delete<T>(endpoint: string): Promise<APIResponse<T>> {
    return this.request<T>('DELETE', endpoint);
  }

  private async request<T>(
    method: string,
    endpoint: string,
    data?: any
  ): Promise<APIResponse<T>> {
    const url = `+"`${this.baseURL}${endpoint}`"+`;
    
    const config: RequestInit = {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
    };

    if (data) {
      config.body = JSON.stringify(data);
    }

    try {
      const response = await fetch(url, config);
      const responseData = await response.json();

      if (!response.ok) {
        throw {
          message: responseData.message || 'Request failed',
          status: response.status,
          code: responseData.code,
        } as APIError;
      }

      return {
        data: responseData,
        status: response.status,
      };
    } catch (error) {
      if (error instanceof Error) {
        throw {
          message: error.message,
          status: 0,
        } as APIError;
      }
      throw error;
    }
  }
}

export const apiClient = new APIClient();`, config.Name)
}

// generateSharedTypes generates shared types content
func (mg *MobileGenerator) generateSharedTypes(config *models.ProjectConfig) string {
	return fmt.Sprintf(`// %s Shared Types
// Generated by Open Source Project Generator

export interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export interface AuthResponse {
  user: User;
  token: string;
  refreshToken: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface HealthCheck {
  status: string;
  timestamp: string;
  service: string;
  version: string;
}

export interface StatusResponse {
  message: string;
  service: string;
}

// Common utility types
export type LoadingState = 'idle' | 'loading' | 'success' | 'error';

export interface AsyncState<T> {
  data: T | null;
  loading: LoadingState;
  error: string | null;
}

// Navigation types
export interface NavigationItem {
  id: string;
  title: string;
  icon?: string;
  route: string;
  children?: NavigationItem[];
}`, config.Name)
}

// generateSharedConstants generates shared constants content
func (mg *MobileGenerator) generateSharedConstants(config *models.ProjectConfig) string {
	return fmt.Sprintf(`// %s Shared Constants
// Generated by Open Source Project Generator

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/auth/login',
    REGISTER: '/auth/register',
    REFRESH: '/auth/refresh',
    LOGOUT: '/auth/logout',
  },
  USER: {
    PROFILE: '/user/profile',
    UPDATE: '/user/update',
  },
  HEALTH: '/health',
  STATUS: '/status',
} as const;

export const APP_CONFIG = {
  NAME: '%s',
  VERSION: '1.0.0',
  API_TIMEOUT: 10000,
  RETRY_ATTEMPTS: 3,
} as const;

export const STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  REFRESH_TOKEN: 'refresh_token',
  USER_DATA: 'user_data',
  THEME: 'theme',
  LANGUAGE: 'language',
} as const;

export const COLORS = {
  PRIMARY: '#6366f1',
  SECONDARY: '#8b5cf6',
  SUCCESS: '#10b981',
  WARNING: '#f59e0b',
  ERROR: '#ef4444',
  INFO: '#3b82f6',
  LIGHT: '#f8fafc',
  DARK: '#1e293b',
} as const;

export const BREAKPOINTS = {
  SM: 640,
  MD: 768,
  LG: 1024,
  XL: 1280,
  '2XL': 1536,
} as const;`, config.Name, config.Name)
}

// Helper functions
func (mg *MobileGenerator) getAndroidPackageName(config *models.ProjectConfig) string {
	org := strings.ToLower(strings.ReplaceAll(config.Organization, " ", ""))
	name := strings.ToLower(strings.ReplaceAll(config.Name, " ", ""))
	return fmt.Sprintf("com.%s.%s", org, name)
}

func (mg *MobileGenerator) getAndroidPackagePath(config *models.ProjectConfig) string {
	packageName := mg.getAndroidPackageName(config)
	return strings.ReplaceAll(packageName, ".", "/")
}
