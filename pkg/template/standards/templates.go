package standards

import (
	"fmt"
	"strings"
)

// TemplateGenerator generates standardized template files
type TemplateGenerator struct {
	standards *FrontendStandards
}

// NewTemplateGenerator creates a new template generator
func NewTemplateGenerator() *TemplateGenerator {
	return &TemplateGenerator{
		standards: GetFrontendStandards(),
	}
}

// GenerateStandardizedPackageJSON generates a standardized package.json template
func (tg *TemplateGenerator) GenerateStandardizedPackageJSON(templateType string) string {
	templateName := strings.TrimPrefix(templateType, "nextjs-")

	// Base template
	template := `{
  "name": "{{.Name}}-` + templateName + `",
  "version": "` + tg.standards.PackageJSON.Metadata.Version + `",
  "private": ` + fmt.Sprintf("%t", tg.standards.PackageJSON.Metadata.Private) + `,
  "description": "{{.Description}} - ` + getTemplateDescription(templateType) + `",
  "author": "{{if .Author}}{{.Author}}{{if .Email}} <{{.Email}}>{{end}}{{end}}",
  "license": "{{.License}}",
  {{if .Repository}}"repository": {
    "type": "git",
    "url": "{{.Repository}}"
  },{{end}}
  "scripts": {`

	// Add scripts
	scripts := tg.getScriptsForTemplate(templateType)
	scriptLines := make([]string, 0, len(scripts))
	for name, command := range scripts {
		scriptLines = append(scriptLines, fmt.Sprintf(`    "%s": "%s"`, name, command))
	}
	template += "\n" + strings.Join(scriptLines, ",\n") + "\n  },"

	// Add dependencies
	template += `
  "dependencies": {`

	deps := tg.getDependenciesForTemplate(templateType)
	depLines := make([]string, 0, len(deps))
	for name, version := range deps {
		depLines = append(depLines, fmt.Sprintf(`    "%s": "%s"`, name, version))
	}
	template += "\n" + strings.Join(depLines, ",\n") + "\n  },"

	// Add dev dependencies
	template += `
  "devDependencies": {`

	devDepLines := make([]string, 0, len(tg.standards.PackageJSON.DevDeps))
	for name, version := range tg.standards.PackageJSON.DevDeps {
		devDepLines = append(devDepLines, fmt.Sprintf(`    "%s": "%s"`, name, version))
	}
	template += "\n" + strings.Join(devDepLines, ",\n") + "\n  },"

	// Add engines
	template += `
  "engines": {`

	engineLines := make([]string, 0, len(tg.standards.PackageJSON.Engines))
	for name, version := range tg.standards.PackageJSON.Engines {
		engineLines = append(engineLines, fmt.Sprintf(`    "%s": "%s"`, name, version))
	}
	template += "\n" + strings.Join(engineLines, ",\n") + "\n  }"

	template += "\n}"

	return template
}

// GenerateStandardizedTSConfig generates a standardized tsconfig.json template
func (tg *TemplateGenerator) GenerateStandardizedTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "es6"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"],
      "@/components/*": ["./src/components/*"],
      "@/lib/*": ["./src/lib/*"],
      "@/hooks/*": ["./src/hooks/*"],
      "@/context/*": ["./src/context/*"],
      "@/types/*": ["./src/types/*"],
      "@/api/*": ["./src/api/*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}`
}

// GenerateStandardizedESLintConfig generates a standardized .eslintrc.json template
func (tg *TemplateGenerator) GenerateStandardizedESLintConfig() string {
	return `{
  "extends": [
    "next/core-web-vitals",
    "@typescript-eslint/recommended"
  ],
  "parser": "@typescript-eslint/parser",
  "plugins": ["@typescript-eslint"],
  "rules": {
    "@typescript-eslint/no-unused-vars": "error",
    "@typescript-eslint/no-explicit-any": "warn",
    "@typescript-eslint/explicit-function-return-type": "off",
    "@typescript-eslint/explicit-module-boundary-types": "off",
    "@typescript-eslint/no-empty-function": "warn",
    "prefer-const": "error",
    "no-var": "error"
  },
  "ignorePatterns": ["node_modules/", ".next/", "out/", "dist/"]
}`
}

// GenerateStandardizedPrettierConfig generates a standardized .prettierrc template
func (tg *TemplateGenerator) GenerateStandardizedPrettierConfig() string {
	return `{
  "semi": true,
  "trailingComma": "es5",
  "singleQuote": true,
  "printWidth": 80,
  "tabWidth": 2,
  "useTabs": false,
  "plugins": ["prettier-plugin-tailwindcss"]
}`
}

// GenerateStandardizedVercelConfig generates a standardized vercel.json template
func (tg *TemplateGenerator) GenerateStandardizedVercelConfig() string {
	return `{
  "buildCommand": "npm run build",
  "devCommand": "npm run dev",
  "installCommand": "npm install",
  "framework": "nextjs",
  "regions": ["iad1"],
  "env": {
    "NEXT_PUBLIC_APP_NAME": "{{.Name}}",
    "NEXT_PUBLIC_APP_VERSION": "0.1.0"
  },
  "build": {
    "env": {
      "NEXT_TELEMETRY_DISABLED": "1"
    }
  },
  "functions": {
    "app/api/**/*.ts": {
      "maxDuration": 30
    }
  },
  "headers": [
    {
      "source": "/(.*)",
      "headers": [
        {
          "key": "X-Frame-Options",
          "value": "DENY"
        },
        {
          "key": "X-Content-Type-Options",
          "value": "nosniff"
        },
        {
          "key": "Referrer-Policy",
          "value": "strict-origin-when-cross-origin"
        },
        {
          "key": "Permissions-Policy",
          "value": "camera=(), microphone=(), geolocation=()"
        }
      ]
    }
  ],
  "rewrites": [
    {
      "source": "/api/:path*",
      "destination": "/api/:path*"
    }
  ]
}`
}

// GenerateStandardizedTailwindConfig generates a standardized tailwind.config.js template
func (tg *TemplateGenerator) GenerateStandardizedTailwindConfig() string {
	return tg.standards.GenerateTailwindConfig()
}

// GenerateStandardizedNextConfig generates a standardized next.config.js template
func (tg *TemplateGenerator) GenerateStandardizedNextConfig() string {
	return tg.standards.GenerateNextConfig()
}

// GenerateStandardizedPostCSSConfig generates a standardized postcss.config.js template
func (tg *TemplateGenerator) GenerateStandardizedPostCSSConfig() string {
	return `module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
`
}

// GenerateStandardizedJestConfig generates a standardized jest.config.js template
func (tg *TemplateGenerator) GenerateStandardizedJestConfig() string {
	return `const nextJest = require('next/jest')

const createJestConfig = nextJest({
  // Provide the path to your Next.js app to load next.config.js and .env files
  dir: './',
})

// Add any custom config to be passed to Jest
const customJestConfig = {
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
  moduleNameMapping: {
    // Handle module aliases (this will be automatically configured for you based on your tsconfig.json paths)
    '^@/components/(.*)$': '<rootDir>/src/components/$1',
    '^@/pages/(.*)$': '<rootDir>/src/pages/$1',
    '^@/lib/(.*)$': '<rootDir>/src/lib/$1',
    '^@/hooks/(.*)$': '<rootDir>/src/hooks/$1',
    '^@/context/(.*)$': '<rootDir>/src/context/$1',
    '^@/types/(.*)$': '<rootDir>/src/types/$1',
    '^@/api/(.*)$': '<rootDir>/src/api/$1',
  },
  testEnvironment: 'jest-environment-jsdom',
}

// createJestConfig is exported this way to ensure that next/jest can load the Next.js config which is async
module.exports = createJestConfig(customJestConfig)
`
}

// GenerateStandardizedJestSetup generates a standardized jest.setup.js template
func (tg *TemplateGenerator) GenerateStandardizedJestSetup() string {
	return `import '@testing-library/jest-dom'
`
}

// Helper functions

// getScriptsForTemplate returns scripts for a specific template type
func (tg *TemplateGenerator) getScriptsForTemplate(templateType string) map[string]string {
	scripts := make(map[string]string)
	for k, v := range tg.standards.PackageJSON.Scripts {
		scripts[k] = v
	}

	// Adjust for template-specific ports
	if port, exists := tg.standards.PackageJSON.Ports[templateType]; exists && templateType != "nextjs-app" {
		scripts["dev"] = fmt.Sprintf("next dev -p %d", port)
		scripts["start"] = fmt.Sprintf("next start -p %d", port)
	}

	return scripts
}

// getDependenciesForTemplate returns dependencies for a specific template type
func (tg *TemplateGenerator) getDependenciesForTemplate(templateType string) map[string]string {
	deps := make(map[string]string)

	// Add base dependencies
	for k, v := range tg.standards.PackageJSON.Dependencies {
		deps[k] = v
	}

	// Add template-specific dependencies
	for k, v := range GetTemplateSpecificDependencies(templateType) {
		deps[k] = v
	}

	return deps
}

// getTemplateDescription returns a description for a template type
func getTemplateDescription(templateType string) string {
	switch templateType {
	case "nextjs-app":
		return "Main Application"
	case "nextjs-home":
		return "Landing Page"
	case "nextjs-admin":
		return "Admin Dashboard"
	default:
		return "Frontend Application"
	}
}

// StandardizedTemplateFiles represents all standardized template files
type StandardizedTemplateFiles struct {
	PackageJSON    string
	TSConfig       string
	ESLintConfig   string
	PrettierConfig string
	VercelConfig   string
	TailwindConfig string
	NextConfig     string
	PostCSSConfig  string
	JestConfig     string
	JestSetup      string
}

// GenerateAllStandardizedFiles generates all standardized template files for a template type
func (tg *TemplateGenerator) GenerateAllStandardizedFiles(templateType string) *StandardizedTemplateFiles {
	return &StandardizedTemplateFiles{
		PackageJSON:    tg.GenerateStandardizedPackageJSON(templateType),
		TSConfig:       tg.GenerateStandardizedTSConfig(),
		ESLintConfig:   tg.GenerateStandardizedESLintConfig(),
		PrettierConfig: tg.GenerateStandardizedPrettierConfig(),
		VercelConfig:   tg.GenerateStandardizedVercelConfig(),
		TailwindConfig: tg.GenerateStandardizedTailwindConfig(),
		NextConfig:     tg.GenerateStandardizedNextConfig(),
		PostCSSConfig:  tg.GenerateStandardizedPostCSSConfig(),
		JestConfig:     tg.GenerateStandardizedJestConfig(),
		JestSetup:      tg.GenerateStandardizedJestSetup(),
	}
}

// GetStandardizedFilePaths returns the file paths for standardized template files
func GetStandardizedFilePaths() map[string]string {
	return map[string]string{
		"PackageJSON":    "package.json.tmpl",
		"TSConfig":       "tsconfig.json.tmpl",
		"ESLintConfig":   ".eslintrc.json.tmpl",
		"PrettierConfig": ".prettierrc.tmpl",
		"VercelConfig":   "vercel.json.tmpl",
		"TailwindConfig": "tailwind.config.js.tmpl",
		"NextConfig":     "next.config.js.tmpl",
		"PostCSSConfig":  "postcss.config.js.tmpl",
		"JestConfig":     "jest.config.js.tmpl",
		"JestSetup":      "jest.setup.js.tmpl",
	}
}
