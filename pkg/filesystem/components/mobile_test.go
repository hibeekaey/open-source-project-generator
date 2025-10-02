package components

import (
	"strings"
	"testing"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
)

func TestMobileGenerator_GenerateFiles(t *testing.T) {
	tests := []struct {
		name          string
		config        *models.ProjectConfig
		expectedFiles []string
		expectedError bool
	}{
		{
			name: "Generate all mobile components",
			config: &models.ProjectConfig{
				Name:         "testapp",
				Organization: "testorg",
				Components: models.Components{
					Mobile: models.MobileComponents{
						Android: true,
						IOS:     true,
						Shared:  true,
					},
				},
			},
			expectedFiles: []string{
				"testproject/Mobile/Android/app/build.gradle",
				"testproject/Mobile/Android/build.gradle",
				"testproject/Mobile/Android/app/src/main/AndroidManifest.xml",
				"testproject/Mobile/Android/app/src/main/java/com/testorg/testapp/MainActivity.kt",
				"testproject/Mobile/Android/app/src/main/res/values/strings.xml",
				"testproject/Mobile/Android/app/src/main/res/values/colors.xml",
				"testproject/Mobile/iOS/Package.swift",
				"testproject/Mobile/iOS/Sources/testapp/ContentView.swift",
				"testproject/Mobile/iOS/Sources/testapp/testappApp.swift",
				"testproject/Mobile/iOS/Resources/Info.plist",
				"testproject/Mobile/Shared/api/client.ts",
				"testproject/Mobile/Shared/types/index.ts",
				"testproject/Mobile/Shared/constants/index.ts",
			},
			expectedError: false,
		},
		{
			name: "Generate only Android",
			config: &models.ProjectConfig{
				Name:         "testapp",
				Organization: "testorg",
				Components: models.Components{
					Mobile: models.MobileComponents{
						Android: true,
					},
				},
			},
			expectedFiles: []string{
				"testproject/Mobile/Android/app/build.gradle",
				"testproject/Mobile/Android/build.gradle",
				"testproject/Mobile/Android/app/src/main/AndroidManifest.xml",
			},
			expectedError: false,
		},
		{
			name: "Generate only iOS",
			config: &models.ProjectConfig{
				Name:         "testapp",
				Organization: "testorg",
				Components: models.Components{
					Mobile: models.MobileComponents{
						IOS: true,
					},
				},
			},
			expectedFiles: []string{
				"testproject/Mobile/iOS/Package.swift",
				"testproject/Mobile/iOS/Sources/testapp/ContentView.swift",
				"testproject/Mobile/iOS/Sources/testapp/testappApp.swift",
				"testproject/Mobile/iOS/Resources/Info.plist",
			},
			expectedError: false,
		},
		{
			name:          "Nil config should return error",
			config:        nil,
			expectedFiles: []string{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := NewMockFileSystemOperations()
			mg := NewMobileGenerator(mockFS)

			err := mg.GenerateFiles("testproject", tt.config)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check that expected files were created
			for _, expectedFile := range tt.expectedFiles {
				if !mockFS.FileExists(expectedFile) {
					t.Errorf("Expected file %s was not created", expectedFile)
				}
			}
		})
	}
}

func TestMobileGenerator_generateAndroidBuildGradle(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name:         "testapp",
		Organization: "testorg",
	}

	content := mg.generateAndroidBuildGradle(config)

	expectedElements := []string{
		"plugins {",
		"com.android.application",
		"org.jetbrains.kotlin.android",
		"namespace 'com.testorg.testapp'",
		"applicationId \"com.testorg.testapp\"",
		"compileSdk 34",
		"minSdk 24",
		"targetSdk 34",
		"androidx.core:core-ktx",
		"androidx.compose.ui:ui",
		"androidx.navigation:navigation-compose",
		"com.squareup.retrofit2:retrofit",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Android build.gradle should contain %s", element)
		}
	}
}

func TestMobileGenerator_generateAndroidManifest(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := mg.generateAndroidManifest(config)

	expectedElements := []string{
		"<?xml version=\"1.0\" encoding=\"utf-8\"?>",
		"<manifest",
		"<uses-permission android:name=\"android.permission.INTERNET\" />",
		"<application",
		"android:label=\"@string/app_name\"",
		"@style/Theme.testapp",
		"<activity",
		"android:name=\".MainActivity\"",
		"android.intent.action.MAIN",
		"android.intent.category.LAUNCHER",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("AndroidManifest.xml should contain %s", element)
		}
	}
}

func TestMobileGenerator_generateMainActivity(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name:         "testapp",
		Organization: "testorg",
	}

	content := mg.generateMainActivity(config)

	expectedElements := []string{
		"package com.testorg.testapp",
		"class MainActivity : ComponentActivity()",
		"override fun onCreate(savedInstanceState: Bundle?)",
		"testappTheme {",
		"fun Greeting(name: String",
		"Hello $name from testapp!",
		"@Preview",
		"fun GreetingPreview()",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("MainActivity.kt should contain %s", element)
		}
	}
}

func TestMobileGenerator_generateAndroidStrings(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := mg.generateAndroidStrings(config)

	expectedElements := []string{
		"<?xml version=\"1.0\" encoding=\"utf-8\"?>",
		"<resources>",
		"<string name=\"app_name\">testapp</string>",
		"<string name=\"welcome_message\">Welcome to testapp</string>",
		"<string name=\"loading\">Loading...</string>",
		"<string name=\"error_network\">Network error occurred</string>",
		"<string name=\"retry\">Retry</string>",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("strings.xml should contain %s", element)
		}
	}
}

func TestMobileGenerator_generateIOSPackageSwift(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := mg.generateIOSPackageSwift(config)

	expectedElements := []string{
		"// swift-tools-version: 5.9",
		"// testapp iOS App",
		"import PackageDescription",
		"let package = Package(",
		"name: \"testapp\"",
		"platforms: [",
		".iOS(.v15)",
		".library(",
		".target(",
		".testTarget(",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Package.swift should contain %s", element)
		}
	}
}

func TestMobileGenerator_generateIOSContentView(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := mg.generateIOSContentView(config)

	expectedElements := []string{
		"import SwiftUI",
		"struct ContentView: View {",
		"var body: some View {",
		"VStack {",
		"Image(systemName: \"globe\")",
		"Text(\"Hello, testapp!\")",
		"#Preview {",
		"ContentView()",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("ContentView.swift should contain %s", element)
		}
	}
}

func TestMobileGenerator_generateIOSApp(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := mg.generateIOSApp(config)

	expectedElements := []string{
		"import SwiftUI",
		"@main",
		"struct testappApp: App {",
		"var body: some Scene {",
		"WindowGroup {",
		"ContentView()",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("testappApp.swift should contain %s", element)
		}
	}
}

func TestMobileGenerator_generateSharedAPIClient(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := mg.generateSharedAPIClient(config)

	expectedElements := []string{
		"// testapp Shared API Client",
		"export interface APIResponse<T>",
		"export interface APIError",
		"export class APIClient",
		"private baseURL: string",
		"async get<T>(endpoint: string)",
		"async post<T>(endpoint: string, data?: any)",
		"async put<T>(endpoint: string, data?: any)",
		"async delete<T>(endpoint: string)",
		"private async request<T>(",
		"export const apiClient = new APIClient()",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("API client should contain %s", element)
		}
	}
}

func TestMobileGenerator_generateSharedTypes(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := mg.generateSharedTypes(config)

	expectedElements := []string{
		"// testapp Shared Types",
		"export interface User",
		"export interface AuthResponse",
		"export interface LoginRequest",
		"export interface RegisterRequest",
		"export interface HealthCheck",
		"export interface StatusResponse",
		"export type LoadingState",
		"export interface AsyncState<T>",
		"export interface NavigationItem",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Shared types should contain %s", element)
		}
	}
}

func TestMobileGenerator_generateSharedConstants(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name: "testapp",
	}

	content := mg.generateSharedConstants(config)

	expectedElements := []string{
		"// testapp Shared Constants",
		"export const API_ENDPOINTS",
		"AUTH: {",
		"LOGIN: '/auth/login'",
		"USER: {",
		"PROFILE: '/user/profile'",
		"export const APP_CONFIG",
		"NAME: 'testapp'",
		"export const STORAGE_KEYS",
		"AUTH_TOKEN: 'auth_token'",
		"export const COLORS",
		"PRIMARY: '#6366f1'",
		"export const BREAKPOINTS",
	}

	for _, element := range expectedElements {
		if !strings.Contains(content, element) {
			t.Errorf("Shared constants should contain %s", element)
		}
	}
}

func TestMobileGenerator_getAndroidPackageName(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	tests := []struct {
		name        string
		config      *models.ProjectConfig
		expectedPkg string
	}{
		{
			name: "Simple names",
			config: &models.ProjectConfig{
				Name:         "testapp",
				Organization: "testorg",
			},
			expectedPkg: "com.testorg.testapp",
		},
		{
			name: "Names with spaces",
			config: &models.ProjectConfig{
				Name:         "Test App",
				Organization: "Test Org",
			},
			expectedPkg: "com.testorg.testapp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mg.getAndroidPackageName(tt.config)
			if result != tt.expectedPkg {
				t.Errorf("Expected package name %s, got %s", tt.expectedPkg, result)
			}
		})
	}
}

func TestMobileGenerator_getAndroidPackagePath(t *testing.T) {
	mockFS := NewMockFileSystemOperations()
	mg := NewMobileGenerator(mockFS)

	config := &models.ProjectConfig{
		Name:         "testapp",
		Organization: "testorg",
	}

	result := mg.getAndroidPackagePath(config)
	expected := "com/testorg/testapp"

	if result != expected {
		t.Errorf("Expected package path %s, got %s", expected, result)
	}
}
