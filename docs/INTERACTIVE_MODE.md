# Interactive Mode User Guide

Complete guide to using the interactive configuration wizard in the Open Source Project Generator.

## What is Interactive Mode?

Interactive mode is a guided wizard that helps you configure and generate projects without manually creating configuration files. It walks you through each step with prompts, validation, and helpful defaults.

## When to Use Interactive Mode

**Use Interactive Mode When:**

- You're new to the generator and want guidance
- You want to quickly prototype a project structure
- You prefer visual prompts over editing YAML files
- You're not sure what configuration options are available
- You want real-time validation of your inputs

**Use Configuration Files When:**

- You need to version control your project configuration
- You're automating project generation in CI/CD
- You want to reuse the same configuration multiple times
- You need to share configuration with team members
- You're generating multiple similar projects

## Starting Interactive Mode

```bash
# Start the interactive wizard
generator generate --interactive

# Or use the short flag
generator generate -i
```

## Interactive Mode Workflow

The wizard guides you through 5 main steps:

### Step 1: Project Information

The wizard first collects basic project information.

**Prompts:**

```text
? Project name: _
? Project description (optional): _
? Output directory: _
? Author (optional): _
? License (optional): _
```

**Example:**

```text
? Project name: my-awesome-app
? Project description (optional): A full-stack web application
? Output directory: ./my-awesome-app
? Author (optional): John Doe
? License (optional): MIT
```

**Tips:**

- **Project name**: Use lowercase with hyphens (e.g., `my-app`, `api-server`)
- **Output directory**: Can be relative (`./my-app`) or absolute (`/home/user/projects/my-app`)
- **Press Enter** to skip optional fields
- **Press Ctrl+C** at any time to cancel

### Step 2: Component Selection

Select which components to include in your project.

**Prompt:**

```text
? Select components to include (use space to select, enter to confirm):
  ◯ Next.js Frontend - Modern React framework with TypeScript
  ◯ Go Backend - RESTful API server with Gin framework
  ◯ Android App - Native Android app with Kotlin
  ◯ iOS App - Native iOS app with Swift
```

**Navigation:**

- **↑/↓ arrows**: Move between options
- **Space**: Select/deselect option
- **Enter**: Confirm selection
- **Ctrl+C**: Cancel

**Example:**

```text
? Select components to include:
  ◉ Next.js Frontend - Modern React framework with TypeScript
  ◉ Go Backend - RESTful API server with Gin framework
  ◯ Android App - Native Android app with Kotlin
  ◯ iOS App - Native iOS app with Swift

Selected: Next.js Frontend, Go Backend
```

**Tips:**

- Select at least one component
- You can select multiple components
- Components will be integrated automatically

### Step 3: Component Configuration

For each selected component, configure specific options.

#### Next.js Configuration

```text
Configuring: Next.js Frontend

? Enable TypeScript? (Y/n): Y
? Include Tailwind CSS? (Y/n): Y
? Use App Router? (Y/n): Y
? Include ESLint? (Y/n): Y
? Use src/ directory? (y/N): N
```

**Options:**

- **TypeScript**: Type-safe JavaScript (recommended: Yes)
- **Tailwind CSS**: Utility-first CSS framework (recommended: Yes)
- **App Router**: Next.js 13+ App Router vs Pages Router (recommended: Yes)
- **ESLint**: Code linting (recommended: Yes)
- **src/ directory**: Organize code in src/ folder (default: No)

**Defaults:**

- All options default to recommended values
- Press Enter to accept defaults
- Type `y` or `n` to change

#### Go Backend Configuration

```text
Configuring: Go Backend

? Go module path: github.com/myorg/my-awesome-app
? Web framework (gin/echo/fiber): gin
? Server port: 8080
? Enable CORS? (y/N): Y
? Include authentication? (y/N): N
```

**Options:**

- **Module path**: Go module name (required)
  - Format: `github.com/org/project` or `example.com/api`
  - Must be valid Go module path
- **Framework**: Web framework choice
  - `gin`: Fast, popular (recommended)
  - `echo`: Minimalist, high performance
  - `fiber`: Express-inspired, very fast
- **Port**: HTTP server port (default: 8080)
  - Must be 1-65535
  - Common: 8080, 8000, 3000, 9000
- **CORS**: Enable Cross-Origin Resource Sharing
- **Authentication**: Include auth middleware

**Tips:**

- Use your GitHub username/org in module path
- Port 8080 is standard for development
- Enable CORS if frontend is on different port

#### Android Configuration

```text
Configuring: Android App

? Package name: com.example.myawesomeapp
? Minimum SDK level: 24
? Target SDK level: 34
? Programming language (kotlin/java): kotlin
? Use Jetpack Compose? (Y/n): Y
```

**Options:**

- **Package name**: Java package name (required)
  - Format: `com.company.app` (lowercase, dots)
  - Must be valid Java package name
  - At least 2 segments
- **Minimum SDK**: Minimum Android version
  - 21 = Android 5.0 (Lollipop)
  - 24 = Android 7.0 (Nougat) - recommended
  - Higher = fewer devices, newer features
- **Target SDK**: Target Android version
  - Should be latest (currently 34)
  - Must be >= Minimum SDK
- **Language**: Programming language
  - `kotlin`: Modern, recommended
  - `java`: Legacy support
- **Jetpack Compose**: Modern UI toolkit (recommended: Yes)

**Tips:**

- Use reverse domain for package name
- Min SDK 24 covers ~95% of devices
- Always target latest SDK

#### iOS Configuration

```text
Configuring: iOS App

? Bundle identifier: com.example.myawesomeapp
? Deployment target: 15.0
? Programming language (swift/objective-c): swift
? Use SwiftUI? (Y/n): Y
```

**Options:**

- **Bundle ID**: iOS bundle identifier (required)
  - Format: `com.company.app` (reverse domain)
  - Must be valid bundle identifier
  - At least 2 segments
- **Deployment target**: Minimum iOS version
  - "15.0" = iOS 15 (recommended)
  - "16.0" = iOS 16
  - "17.0" = iOS 17
  - Higher = fewer devices, newer features
- **Language**: Programming language
  - `swift`: Modern, recommended
  - `objective-c`: Legacy support
- **SwiftUI**: Modern UI framework (recommended: Yes)

**Tips:**

- Use same domain as Android package
- iOS 15.0 covers ~95% of devices
- SwiftUI is the future of iOS development

### Step 4: Integration Options

Configure how components integrate together.

```text
Integration Options

? Generate Docker Compose file? (Y/n): Y
? Generate build scripts? (Y/n): Y
? Backend API URL: http://localhost:8080
? Frontend URL: http://localhost:3000
```

**Options:**

- **Docker Compose**: Generate docker-compose.yml
  - Includes all components as services
  - Configures networking
  - Sets up volumes
- **Build scripts**: Generate Makefile and scripts
  - Build, run, test commands
  - Development helpers
- **API URL**: Backend API endpoint
  - Used by frontend and mobile apps
  - Default: http://localhost:8080
- **Frontend URL**: Frontend application URL
  - Used for CORS configuration
  - Default: http://localhost:3000

**Tips:**

- Enable Docker Compose for easy development
- Build scripts save time
- Use localhost URLs for development

### Step 5: Confirmation

Review your configuration and confirm generation.

```text
Configuration Summary
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Project Information:
  Name: my-awesome-app
  Description: A full-stack web application
  Output Directory: ./my-awesome-app
  Author: John Doe
  License: MIT

Components:
  ✓ Next.js Frontend (web-app)
    - TypeScript: Yes
    - Tailwind CSS: Yes
    - App Router: Yes
    - ESLint: Yes

  ✓ Go Backend (api-server)
    - Module: github.com/myorg/my-awesome-app
    - Framework: gin
    - Port: 8080
    - CORS: Enabled

Integration:
  ✓ Docker Compose: Enabled
  ✓ Build Scripts: Enabled
  ✓ API URL: http://localhost:8080
  ✓ Frontend URL: http://localhost:3000

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

? Proceed with generation? (Y/n): _
```

**Actions:**

- **Y (Yes)**: Proceed with generation
- **N (No)**: Cancel and exit
- **E (Edit)**: Go back and edit configuration (if supported)

**What Happens Next:**

1. Configuration is validated
2. Required tools are checked
3. Components are generated
4. Files are integrated
5. Success message is displayed

## Input Validation

Interactive mode validates all inputs in real-time.

### Validation Examples

**Invalid Project Name:**

```text
? Project name: My Project!
✗ Invalid project name. Use lowercase letters, numbers, and hyphens only.
? Project name: my-project
✓ Valid
```

**Invalid Go Module Path:**

```text
? Go module path: my project
✗ Invalid Go module path. Use format: github.com/org/project
? Go module path: github.com/myorg/my-project
✓ Valid
```

**Invalid Port Number:**

```text
? Server port: 70000
✗ Invalid port number. Must be between 1 and 65535.
? Server port: 8080
✓ Valid
```

**Invalid Package Name:**

```text
? Package name: Com.Example.App
✗ Invalid package name. Must be lowercase with dots (e.g., com.example.app)
? Package name: com.example.app
✓ Valid
```

## Error Handling

### Cancellation

Press **Ctrl+C** at any prompt to cancel:

```text
? Project name: my-app
^C
Operation cancelled by user.
Exit code: 5
```

### Invalid Input

If you enter invalid input, you'll be prompted again:

```text
? Server port: abc
✗ Invalid input. Please enter a number.
? Server port: 8080
✓ Valid
```

### Tool Not Found

If required tools are missing, you'll see installation instructions:

```text
✗ Required tool 'npx' not found

Install instructions:
  macOS: brew install node
  Ubuntu: sudo apt install nodejs npm
  Windows: Download from https://nodejs.org/

? Continue with fallback generation? (y/N): _
```

## Common Workflows

### Full-Stack Web Application

```bash
generator generate --interactive

# Follow prompts:
# 1. Project name: fullstack-app
# 2. Select: Next.js Frontend, Go Backend
# 3. Configure Next.js: All defaults (TypeScript, Tailwind, etc.)
# 4. Configure Go: Module path, port 8080
# 5. Integration: Enable Docker Compose and scripts
# 6. Confirm and generate
```

### Frontend-Only Application

```bash
generator generate --interactive

# Follow prompts:
# 1. Project name: frontend-app
# 2. Select: Next.js Frontend only
# 3. Configure Next.js: TypeScript, Tailwind, App Router
# 4. Integration: Disable Docker Compose
# 5. Confirm and generate
```

### Mobile Application with Backend

```bash
generator generate --interactive

# Follow prompts:
# 1. Project name: mobile-app
# 2. Select: Android, iOS, Go Backend
# 3. Configure Android: Package name, SDK levels
# 4. Configure iOS: Bundle ID, deployment target
# 5. Configure Go: Module path, port 8080
# 6. Integration: Enable Docker Compose, set API URL
# 7. Confirm and generate
```

### Microservice

```bash
generator generate --interactive

# Follow prompts:
# 1. Project name: user-service
# 2. Select: Go Backend only
# 3. Configure Go: Module path, framework (gin/echo/fiber)
# 4. Integration: Enable Docker Compose
# 5. Confirm and generate
```

## Tips and Best Practices

### 1. Use Descriptive Names

```text
✓ Good: user-authentication-service
✓ Good: ecommerce-frontend
✗ Bad: proj1
✗ Bad: test
```

### 2. Accept Defaults When Unsure

Most defaults are sensible and recommended:

- TypeScript: Yes
- Tailwind CSS: Yes
- App Router: Yes
- Kotlin: Yes
- Swift: Yes

### 3. Use Consistent Naming

Use the same domain for Android package and iOS bundle ID:

```text
Android: com.mycompany.myapp
iOS: com.mycompany.myapp
```

### 4. Enable Docker Compose

Docker Compose makes development easier:

- All services start with one command
- Networking configured automatically
- Easy to share with team

### 5. Use Standard Ports

Stick to standard development ports:

- Frontend: 3000
- Backend: 8080
- Database: 5432 (PostgreSQL), 3306 (MySQL)

### 6. Test with Dry Run First

If unsure, use dry run mode:

```bash
generator generate --interactive --dry-run
```

### 7. Save Configuration

After interactive mode, save the configuration:

```bash
# Interactive mode generates config
generator generate --interactive

# Save for reuse
cp .generator/generated-config.yaml my-project-config.yaml
git add my-project-config.yaml
```

## Combining with Other Flags

Interactive mode works with other flags:

### Verbose Output

```bash
generator generate --interactive --verbose
```

See detailed output during generation.

### Streaming Output

```bash
generator generate --interactive --stream-output
```

See real-time output from bootstrap tools.

### Dry Run

```bash
generator generate --interactive --dry-run
```

Preview without creating files.

### Offline Mode

```bash
generator generate --interactive --offline
```

Use cached tool information.

## Troubleshooting

### Interactive Mode Won't Start

**Problem:** Terminal doesn't support interactive prompts

**Solution:**

```bash
# Check terminal type
echo $TERM

# Use non-interactive mode instead
generator init-config project.yaml
generator generate --config project.yaml
```

### Input Not Recognized

**Problem:** Pressing keys doesn't work

**Solution:**

- Ensure terminal supports interactive input
- Try different terminal (iTerm2, Windows Terminal, etc.)
- Use arrow keys for navigation
- Press Enter to confirm

### Validation Keeps Failing

**Problem:** Can't get past validation

**Solution:**

```bash
# See validation rules
generator generate --interactive --verbose

# Or use configuration file
generator init-config --example fullstack
# Edit the file
generator generate --config project.yaml
```

### Want to Change Previous Answer

**Problem:** Made a mistake in earlier step

**Solution:**

- Press Ctrl+C to cancel
- Start over
- Or use configuration file for more control

## Comparison: Interactive vs Configuration File

| Feature | Interactive Mode | Configuration File |
|---------|------------------|-------------------|
| **Ease of Use** | ✓ Guided prompts | Manual editing |
| **Speed** | Slower (step-by-step) | ✓ Faster (if you know what you want) |
| **Validation** | ✓ Real-time | On generation |
| **Reusability** | Generate config file | ✓ Reuse directly |
| **Version Control** | Generate then save | ✓ Direct |
| **CI/CD** | Not suitable | ✓ Ideal |
| **Learning** | ✓ Discover options | Need documentation |
| **Flexibility** | Limited to prompts | ✓ Full control |

## Next Steps

After generating your project with interactive mode:

1. **Explore the generated structure**
   ```bash
   cd my-awesome-app
   tree -L 2
   ```

2. **Read the generated README**
   ```bash
   cat README.md
   ```

3. **Install dependencies**
   ```bash
   # Frontend
   cd App && npm install

   # Backend
   cd CommonServer && go mod download
   ```

4. **Start development**
   ```bash
   # Using Docker Compose
   docker-compose up

   # Or manually
   make dev
   ```

5. **Save configuration for reuse**
   ```bash
   # Configuration is saved in .generator/
   cp .generator/generated-config.yaml my-config.yaml
   git add my-config.yaml
   ```

## See Also

- [Getting Started](GETTING_STARTED.md) - Installation and quick start
- [CLI Commands](CLI_COMMANDS.md) - Command reference
- [Configuration Guide](CONFIGURATION.md) - Configuration file format
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues
- [Examples](EXAMPLES.md) - Example configurations
