# Templates Directory

This directory contains all the template files used by the Open Source Template Generator.

## Structure

```
templates/
├── base/                    # Core project files
│   ├── root/               # Root-level files (Makefile, docker-compose.yml, etc.)
│   ├── docs/               # Documentation templates
│   └── github/             # GitHub-specific files
├── frontend/               # Frontend application templates
│   ├── nextjs-app/         # Main application template
│   ├── nextjs-home/        # Landing page template
│   └── nextjs-admin/       # Admin dashboard template
├── backend/                # Backend service templates
│   └── go-gin/             # Go + Gin API server template
├── mobile/                 # Mobile application templates
│   ├── android-kotlin/     # Android Kotlin template
│   └── ios-swift/          # iOS Swift template
├── infrastructure/         # Infrastructure templates
│   ├── terraform/          # Terraform configurations
│   ├── kubernetes/         # K8s manifests
│   └── docker/             # Docker configurations
└── config/                 # Configuration templates
    ├── versions.yaml       # Latest package versions
    └── defaults.yaml       # Default configurations
```

## Template Processing

Templates use Go's text/template syntax with custom functions for:
- Variable substitution
- Conditional rendering based on selected components
- Version management
- String manipulation

## Adding New Templates

1. Create the template files in the appropriate directory
2. Add template metadata (if required)
3. Update the template processing logic
4. Test the template generation