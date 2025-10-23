# Fallback Templates

This directory contains minimal embedded templates used for fallback project generation when external bootstrap tools are not available.

## Structure

- `android/` - Android project templates (Kotlin, Gradle)
- `ios/` - iOS project templates (Swift, Xcode)

## Usage

These templates are embedded into the binary using Go's `embed` directive. The fallback generators use these templates when:

1. Required external tools (gradle, xcodebuild) are not installed
2. User explicitly requests fallback generation with `--no-external-tools` flag
3. Operating in offline mode without cached tools

## Template Format

Templates in this directory are minimal and focus on:

- Basic project structure
- Essential configuration files
- Starter code with comments
- Comprehensive README with setup instructions

## Adding New Templates

To add templates for a new platform:

1. Create a new subdirectory (e.g., `flutter/`)
2. Add minimal project files
3. Update the corresponding generator in `../`
4. Add embed directive: `//go:embed templates/flutter/*`

## Notes

- Keep templates minimal to reduce binary size
- Include clear documentation in generated README files
- Focus on getting developers started quickly
- Provide links to official documentation
