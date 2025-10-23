#!/bin/bash
# Script to check latest versions of dependencies
# This helps identify which hardcoded versions are outdated

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Checking Latest Versions ===${NC}\n"

# Function to print section header
print_section() {
    echo -e "\n${BLUE}### $1 ###${NC}"
}

# Function to check npm package version
check_npm() {
    local package=$1
    local current=$2
    echo -n "Checking $package... "
    if command -v npm >/dev/null 2>&1; then
        latest=$(npm view "$package" version 2>/dev/null || echo "N/A")
        if [ "$latest" != "N/A" ]; then
            echo -e "${GREEN}Latest: $latest${NC} (Current: $current)"
        else
            echo -e "${RED}Failed to fetch${NC}"
        fi
    else
        echo -e "${YELLOW}npm not installed${NC}"
    fi
}

# Function to check Go module version
check_go_module() {
    local module=$1
    local current=$2
    echo -n "Checking $module... "
    if command -v go >/dev/null 2>&1; then
        latest=$(go list -m -versions "$module" 2>/dev/null | awk '{print $NF}' || echo "N/A")
        if [ "$latest" != "N/A" ] && [ -n "$latest" ]; then
            echo -e "${GREEN}Latest: $latest${NC} (Current: $current)"
        else
            echo -e "${YELLOW}Unable to fetch${NC}"
        fi
    else
        echo -e "${YELLOW}go not installed${NC}"
    fi
}

# Function to check Docker image tag
check_docker_image() {
    local image=$1
    local current=$2
    echo -n "Checking $image... "
    if command -v docker >/dev/null 2>&1; then
        # Note: This requires internet and docker hub access
        echo -e "${YELLOW}Current: $current${NC} (Check manually: https://hub.docker.com/_/$image)"
    else
        echo -e "${YELLOW}docker not installed${NC}"
    fi
}

# Node.js & Frontend
print_section "Node.js & Frontend"
check_npm "create-next-app" "@latest"
check_npm "next" "14.2.0"
check_npm "react" "18.3.1"
check_npm "react-dom" "18.3.1"
check_npm "typescript" "5.3.3"

# Go Backend Frameworks
print_section "Go Backend Frameworks"
check_go_module "github.com/gin-gonic/gin" "@latest"
check_go_module "github.com/gin-contrib/cors" "@latest"
check_go_module "github.com/labstack/echo/v4" "@latest"
check_go_module "github.com/gofiber/fiber/v2" "@latest"

# Android (Maven Central - requires manual check)
print_section "Android Dependencies"
echo "AndroidX Core KTX (current: 1.12.0)"
echo "  Check: https://developer.android.com/jetpack/androidx/releases/core"
echo "AndroidX AppCompat (current: 1.6.1)"
echo "  Check: https://developer.android.com/jetpack/androidx/releases/appcompat"
echo "Material Components (current: 1.10.0)"
echo "  Check: https://github.com/material-components/material-components-android/releases"
echo "ConstraintLayout (current: 2.1.4)"
echo "  Check: https://developer.android.com/jetpack/androidx/releases/constraintlayout"

# Android Build Tools
print_section "Android Build Tools"
echo "Gradle (current: 8.1.0 / wrapper: 8.0)"
echo "  Check: https://gradle.org/releases/"
echo "Kotlin (current: 1.9.0)"
echo "  Check: https://kotlinlang.org/docs/releases.html"
echo "Android SDK (compile: 34, min: 24, target: 34)"
echo "  Check: https://developer.android.com/studio/releases/platforms"

# iOS
print_section "iOS"
echo "Swift (current: 5.9)"
echo "  Check: https://www.swift.org/download/"
echo "Xcode (current: 15.0)"
echo "  Check: https://developer.apple.com/xcode/resources/"

# Docker Base Images
print_section "Docker Base Images"
check_docker_image "alpine" "3.19"
check_docker_image "golang" "1.25-alpine"
check_docker_image "ubuntu" "24.04"

# Go Toolchain
print_section "Go Toolchain"
if command -v go >/dev/null 2>&1; then
    installed=$(go version | awk '{print $3}' | sed 's/go//')
    echo -e "Installed: ${GREEN}$installed${NC} (Required: 1.25.0)"
    echo "Latest: Check https://go.dev/dl/"
else
    echo -e "${YELLOW}go not installed${NC}"
fi

# Development Tools
print_section "Go Development Tools"
echo "Check latest versions with:"
echo "  go list -m -versions github.com/air-verse/air"
echo "  go list -m -versions gotest.tools/gotestsum"
echo "  go list -m -versions github.com/securego/gosec/v2"
echo "  go list -m -versions golang.org/x/vuln"
echo "  go list -m -versions honnef.co/go/tools"

echo -e "\n${BLUE}=== Summary ===${NC}"
echo "This script checks versions where possible via CLI tools."
echo "For packages requiring manual checks, URLs are provided above."
echo ""
echo "To update versions:"
echo "  1. Review the output above"
echo "  2. Update configs/versions.yaml (once implemented)"
echo "  3. Update hardcoded values in source files"
echo "  4. Test thoroughly before committing"

echo -e "\n${YELLOW}Note: Some checks require internet connectivity${NC}"
