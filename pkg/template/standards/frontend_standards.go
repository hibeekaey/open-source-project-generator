package standards

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FrontendStandards defines the standardized configurations for frontend templates
type FrontendStandards struct {
	PackageJSON PackageJSONStandard `json:"package_json"`
	TypeScript  TypeScriptStandard  `json:"typescript"`
	ESLint      ESLintStandard      `json:"eslint"`
	Prettier    PrettierStandard    `json:"prettier"`
	Vercel      VercelStandard      `json:"vercel"`
	TailwindCSS TailwindCSSStandard `json:"tailwindcss"`
	NextConfig  NextConfigStandard  `json:"next_config"`
}

// PackageJSONStandard defines standardized package.json configuration
type PackageJSONStandard struct {
	Scripts      map[string]string `json:"scripts"`
	Dependencies map[string]string `json:"dependencies"`
	DevDeps      map[string]string `json:"dev_dependencies"`
	Engines      map[string]string `json:"engines"`
	Metadata     PackageMetadata   `json:"metadata"`
	Ports        map[string]int    `json:"ports"`
}

// PackageMetadata defines common package.json metadata
type PackageMetadata struct {
	Version string `json:"version"`
	Private bool   `json:"private"`
	License string `json:"license"`
}

// TypeScriptStandard defines standardized TypeScript configuration
type TypeScriptStandard struct {
	CompilerOptions map[string]interface{} `json:"compiler_options"`
	Include         []string               `json:"include"`
	Exclude         []string               `json:"exclude"`
	Paths           map[string][]string    `json:"paths"`
}

// ESLintStandard defines standardized ESLint configuration
type ESLintStandard struct {
	Extends        []string          `json:"extends"`
	Parser         string            `json:"parser"`
	Plugins        []string          `json:"plugins"`
	Rules          map[string]string `json:"rules"`
	IgnorePatterns []string          `json:"ignore_patterns"`
}

// PrettierStandard defines standardized Prettier configuration
type PrettierStandard struct {
	Semi          bool     `json:"semi"`
	TrailingComma string   `json:"trailing_comma"`
	SingleQuote   bool     `json:"single_quote"`
	PrintWidth    int      `json:"print_width"`
	TabWidth      int      `json:"tab_width"`
	UseTabs       bool     `json:"use_tabs"`
	Plugins       []string `json:"plugins"`
}

// VercelStandard defines standardized Vercel deployment configuration
type VercelStandard struct {
	Framework      string                 `json:"framework"`
	BuildCommand   string                 `json:"build_command"`
	DevCommand     string                 `json:"dev_command"`
	InstallCommand string                 `json:"install_command"`
	Regions        []string               `json:"regions"`
	Headers        []VercelHeader         `json:"headers"`
	Functions      map[string]interface{} `json:"functions"`
	Build          map[string]interface{} `json:"build"`
}

// VercelHeader defines security headers for Vercel deployment
type VercelHeader struct {
	Source  string              `json:"source"`
	Headers []map[string]string `json:"headers"`
}

// TailwindCSSStandard defines standardized Tailwind CSS configuration
type TailwindCSSStandard struct {
	Content []string               `json:"content"`
	Theme   map[string]interface{} `json:"theme"`
	Plugins []string               `json:"plugins"`
}

// NextConfigStandard defines standardized Next.js configuration
type NextConfigStandard struct {
	TypeScript   map[string]interface{} `json:"typescript"`
	ESLint       map[string]interface{} `json:"eslint"`
	Experimental map[string]interface{} `json:"experimental"`
	Images       map[string]interface{} `json:"images"`
}

// GetFrontendStandards returns the standardized frontend configuration
func GetFrontendStandards() *FrontendStandards {
	return &FrontendStandards{
		PackageJSON: PackageJSONStandard{
			Scripts: map[string]string{
				"dev":           "next dev",
				"build":         "next build",
				"start":         "next start",
				"lint":          "next lint",
				"lint:fix":      "next lint --fix",
				"type-check":    "tsc --noEmit",
				"test":          "jest",
				"test:watch":    "jest --watch",
				"test:coverage": "jest --coverage",
				"format":        "prettier --write .",
				"format:check":  "prettier --check .",
				"clean":         "rm -rf .next out dist",
			},
			Dependencies: map[string]string{
				"next":                     "{{.Versions.NextJS}}",
				"react":                    "{{.Versions.React}}",
				"react-dom":                "{{.Versions.React}}",
				"typescript":               "^5.3.0",
				"@types/node":              "^20.10.0",
				"@types/react":             "^18.2.0",
				"@types/react-dom":         "^18.2.0",
				"tailwindcss":              "^3.4.0",
				"autoprefixer":             "^10.4.0",
				"postcss":                  "^8.4.0",
				"clsx":                     "^2.0.0",
				"class-variance-authority": "^0.7.0",
				"tailwind-merge":           "^2.2.0",
				"lucide-react":             "^0.300.0",
				"@radix-ui/react-slot":     "^1.0.0",
				"tailwindcss-animate":      "^1.0.7",
			},
			DevDeps: map[string]string{
				"eslint":                           "^8.55.0",
				"eslint-config-next":               "{{.Versions.NextJS}}",
				"@typescript-eslint/eslint-plugin": "^6.15.0",
				"@typescript-eslint/parser":        "^6.15.0",
				"prettier":                         "^3.1.0",
				"prettier-plugin-tailwindcss":      "^0.5.0",
				"jest":                             "^29.7.0",
				"jest-environment-jsdom":           "^29.7.0",
				"@testing-library/react":           "^14.1.0",
				"@testing-library/jest-dom":        "^6.1.0",
				"@testing-library/user-event":      "^14.5.0",
				"@types/jest":                      "^29.5.0",
			},
			Engines: map[string]string{
				"node": ">=22.0.0",
				"npm":  ">=10.0.0",
			},
			Metadata: PackageMetadata{
				Version: "0.1.0",
				Private: true,
				License: "{{.License}}",
			},
			Ports: map[string]int{
				"nextjs-app":   3000,
				"nextjs-home":  3001,
				"nextjs-admin": 3002,
			},
		},
		TypeScript: TypeScriptStandard{
			CompilerOptions: map[string]interface{}{
				"target":            "es5",
				"lib":               []string{"dom", "dom.iterable", "es6"},
				"allowJs":           true,
				"skipLibCheck":      true,
				"strict":            true,
				"noEmit":            true,
				"esModuleInterop":   true,
				"module":            "esnext",
				"moduleResolution":  "bundler",
				"resolveJsonModule": true,
				"isolatedModules":   true,
				"jsx":               "preserve",
				"incremental":       true,
				"plugins": []map[string]string{
					{"name": "next"},
				},
				"baseUrl": ".",
			},
			Paths: map[string][]string{
				"@/*":            {"./src/*"},
				"@/components/*": {"./src/components/*"},
				"@/lib/*":        {"./src/lib/*"},
				"@/hooks/*":      {"./src/hooks/*"},
				"@/context/*":    {"./src/context/*"},
				"@/types/*":      {"./src/types/*"},
				"@/api/*":        {"./src/api/*"},
			},
			Include: []string{"next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"},
			Exclude: []string{"node_modules"},
		},
		ESLint: ESLintStandard{
			Extends: []string{
				"next/core-web-vitals",
				"@typescript-eslint/recommended",
			},
			Parser:  "@typescript-eslint/parser",
			Plugins: []string{"@typescript-eslint"},
			Rules: map[string]string{
				"@typescript-eslint/no-unused-vars":                 "error",
				"@typescript-eslint/no-explicit-any":                "warn",
				"@typescript-eslint/explicit-function-return-type":  "off",
				"@typescript-eslint/explicit-module-boundary-types": "off",
				"@typescript-eslint/no-empty-function":              "warn",
				"prefer-const":                                      "error",
				"no-var":                                            "error",
			},
			IgnorePatterns: []string{"node_modules/", ".next/", "out/", "dist/"},
		},
		Prettier: PrettierStandard{
			Semi:          true,
			TrailingComma: "es5",
			SingleQuote:   true,
			PrintWidth:    80,
			TabWidth:      2,
			UseTabs:       false,
			Plugins:       []string{"prettier-plugin-tailwindcss"},
		},
		Vercel: VercelStandard{
			Framework:      "nextjs",
			BuildCommand:   "npm run build",
			DevCommand:     "npm run dev",
			InstallCommand: "npm install",
			Regions:        []string{"iad1"},
			Headers: []VercelHeader{
				{
					Source: "/(.*)",
					Headers: []map[string]string{
						{"key": "X-Frame-Options", "value": "DENY"},
						{"key": "X-Content-Type-Options", "value": "nosniff"},
						{"key": "Referrer-Policy", "value": "strict-origin-when-cross-origin"},
						{"key": "Permissions-Policy", "value": "camera=(), microphone=(), geolocation=()"},
					},
				},
			},
			Functions: map[string]interface{}{
				"app/api/**/*.ts": map[string]int{"maxDuration": 30},
			},
			Build: map[string]interface{}{
				"env": map[string]string{"NEXT_TELEMETRY_DISABLED": "1"},
			},
		},
		TailwindCSS: TailwindCSSStandard{
			Content: []string{
				"./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
				"./src/components/**/*.{js,ts,jsx,tsx,mdx}",
				"./src/app/**/*.{js,ts,jsx,tsx,mdx}",
			},
			Theme: map[string]interface{}{
				"extend": map[string]interface{}{
					"colors": map[string]interface{}{
						"border":     "hsl(var(--border))",
						"input":      "hsl(var(--input))",
						"ring":       "hsl(var(--ring))",
						"background": "hsl(var(--background))",
						"foreground": "hsl(var(--foreground))",
						"primary": map[string]string{
							"DEFAULT":    "hsl(var(--primary))",
							"foreground": "hsl(var(--primary-foreground))",
						},
						"secondary": map[string]string{
							"DEFAULT":    "hsl(var(--secondary))",
							"foreground": "hsl(var(--secondary-foreground))",
						},
						"destructive": map[string]string{
							"DEFAULT":    "hsl(var(--destructive))",
							"foreground": "hsl(var(--destructive-foreground))",
						},
						"muted": map[string]string{
							"DEFAULT":    "hsl(var(--muted))",
							"foreground": "hsl(var(--muted-foreground))",
						},
						"accent": map[string]string{
							"DEFAULT":    "hsl(var(--accent))",
							"foreground": "hsl(var(--accent-foreground))",
						},
						"popover": map[string]string{
							"DEFAULT":    "hsl(var(--popover))",
							"foreground": "hsl(var(--popover-foreground))",
						},
						"card": map[string]string{
							"DEFAULT":    "hsl(var(--card))",
							"foreground": "hsl(var(--card-foreground))",
						},
					},
					"borderRadius": map[string]string{
						"lg": "var(--radius)",
						"md": "calc(var(--radius) - 2px)",
						"sm": "calc(var(--radius) - 4px)",
					},
				},
			},
			Plugins: []string{"tailwindcss-animate"},
		},
		NextConfig: NextConfigStandard{
			TypeScript: map[string]interface{}{
				"ignoreBuildErrors": false,
			},
			ESLint: map[string]interface{}{
				"ignoreDuringBuilds": false,
			},
			Experimental: map[string]interface{}{
				"typedRoutes": true,
			},
			Images: map[string]interface{}{
				"domains": []string{},
			},
		},
	}
}

// GetTemplateSpecificDependencies returns additional dependencies for specific template types
func GetTemplateSpecificDependencies(templateType string) map[string]string {
	switch templateType {
	case "nextjs-home":
		return map[string]string{
			"@radix-ui/react-accordion":       "^1.1.0",
			"@radix-ui/react-navigation-menu": "^1.1.0",
			"framer-motion":                   "^10.16.0",
			"react-intersection-observer":     "^9.5.0",
		}
	case "nextjs-admin":
		return map[string]string{
			"@radix-ui/react-dropdown-menu": "^2.0.0",
			"@radix-ui/react-select":        "^2.0.0",
			"@radix-ui/react-checkbox":      "^1.0.0",
			"@radix-ui/react-switch":        "^1.0.0",
			"@radix-ui/react-tabs":          "^1.0.0",
			"@radix-ui/react-toast":         "^1.1.0",
			"@radix-ui/react-tooltip":       "^1.0.0",
			"@tanstack/react-table":         "^8.11.0",
			"react-hook-form":               "^7.48.0",
			"@hookform/resolvers":           "^3.3.0",
			"zod":                           "^3.22.0",
			"date-fns":                      "^3.0.0",
			"recharts":                      "^2.8.0",
		}
	case "nextjs-app":
		return map[string]string{
			"@radix-ui/react-dialog":        "^1.0.0",
			"@radix-ui/react-dropdown-menu": "^2.0.0",
			"@radix-ui/react-toast":         "^1.1.0",
		}
	default:
		return map[string]string{}
	}
}

// GeneratePackageJSON generates a standardized package.json for a specific template type
func (fs *FrontendStandards) GeneratePackageJSON(templateType, name, description string) ([]byte, error) {
	pkg := map[string]interface{}{
		"name":        fmt.Sprintf("%s-%s", name, strings.TrimPrefix(templateType, "nextjs-")),
		"version":     fs.PackageJSON.Metadata.Version,
		"private":     fs.PackageJSON.Metadata.Private,
		"description": description,
		"scripts":     fs.PackageJSON.Scripts,
		"engines":     fs.PackageJSON.Engines,
	}

	// Add template-specific port to dev and start scripts
	if port, exists := fs.PackageJSON.Ports[templateType]; exists && templateType != "nextjs-app" {
		scripts := make(map[string]string)
		for k, v := range fs.PackageJSON.Scripts {
			scripts[k] = v
		}
		scripts["dev"] = fmt.Sprintf("next dev -p %d", port)
		scripts["start"] = fmt.Sprintf("next start -p %d", port)
		pkg["scripts"] = scripts
	}

	// Merge base dependencies with template-specific ones
	dependencies := make(map[string]string)
	for k, v := range fs.PackageJSON.Dependencies {
		dependencies[k] = v
	}
	for k, v := range GetTemplateSpecificDependencies(templateType) {
		dependencies[k] = v
	}
	pkg["dependencies"] = dependencies
	pkg["devDependencies"] = fs.PackageJSON.DevDeps

	return json.MarshalIndent(pkg, "", "  ")
}

// GenerateTSConfig generates a standardized tsconfig.json
func (fs *FrontendStandards) GenerateTSConfig() ([]byte, error) {
	config := map[string]interface{}{
		"compilerOptions": fs.TypeScript.CompilerOptions,
		"include":         fs.TypeScript.Include,
		"exclude":         fs.TypeScript.Exclude,
	}

	// Add paths to compiler options
	compilerOptions := config["compilerOptions"].(map[string]interface{})
	compilerOptions["paths"] = fs.TypeScript.Paths

	return json.MarshalIndent(config, "", "  ")
}

// GenerateESLintConfig generates a standardized .eslintrc.json
func (fs *FrontendStandards) GenerateESLintConfig() ([]byte, error) {
	config := map[string]interface{}{
		"extends":        fs.ESLint.Extends,
		"parser":         fs.ESLint.Parser,
		"plugins":        fs.ESLint.Plugins,
		"rules":          fs.ESLint.Rules,
		"ignorePatterns": fs.ESLint.IgnorePatterns,
	}

	return json.MarshalIndent(config, "", "  ")
}

// GeneratePrettierConfig generates a standardized .prettierrc
func (fs *FrontendStandards) GeneratePrettierConfig() ([]byte, error) {
	config := map[string]interface{}{
		"semi":          fs.Prettier.Semi,
		"trailingComma": fs.Prettier.TrailingComma,
		"singleQuote":   fs.Prettier.SingleQuote,
		"printWidth":    fs.Prettier.PrintWidth,
		"tabWidth":      fs.Prettier.TabWidth,
		"useTabs":       fs.Prettier.UseTabs,
		"plugins":       fs.Prettier.Plugins,
	}

	return json.MarshalIndent(config, "", "  ")
}

// GenerateVercelConfig generates a standardized vercel.json
func (fs *FrontendStandards) GenerateVercelConfig(templateType string) ([]byte, error) {
	config := map[string]interface{}{
		"buildCommand":   fs.Vercel.BuildCommand,
		"devCommand":     fs.Vercel.DevCommand,
		"installCommand": fs.Vercel.InstallCommand,
		"framework":      fs.Vercel.Framework,
		"regions":        fs.Vercel.Regions,
		"headers":        fs.Vercel.Headers,
		"functions":      fs.Vercel.Functions,
		"build":          fs.Vercel.Build,
	}

	// Add template-specific environment variables
	env := map[string]string{
		"NEXT_PUBLIC_APP_NAME":    "{{.Name}}",
		"NEXT_PUBLIC_APP_VERSION": "0.1.0",
	}
	config["env"] = env

	// Add rewrites for API routes
	rewrites := []map[string]string{
		{
			"source":      "/api/:path*",
			"destination": "/api/:path*",
		},
	}
	config["rewrites"] = rewrites

	return json.MarshalIndent(config, "", "  ")
}

// GenerateTailwindConfig generates a standardized tailwind.config.js
func (fs *FrontendStandards) GenerateTailwindConfig() string {
	return `/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ["class"],
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
}
`
}

// GenerateNextConfig generates a standardized next.config.js
func (fs *FrontendStandards) GenerateNextConfig() string {
	return `/** @type {import('next').NextConfig} */
const nextConfig = {
  typescript: {
    ignoreBuildErrors: false,
  },
  eslint: {
    ignoreDuringBuilds: false,
  },
  experimental: {
    typedRoutes: true,
  },
  images: {
    domains: [],
  },
}

module.exports = nextConfig
`
}
