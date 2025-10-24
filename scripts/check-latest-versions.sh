#!/bin/bash
# Script to check latest versions of dependencies
# Reads current versions from configs/versions.yaml and compares with latest available
#
# Usage: ./scripts/check-latest-versions.sh [OPTIONS]
#
# Options:
#   --json    Output results in JSON format
#   --quiet   Only show outdated packages
#   --help    Show this help message
#
# Environment Variables:
#   GITHUB_TOKEN  Optional GitHub token for higher rate limits (5000/hour vs 60/hour)
#
# Requirements:
#   - npm (for Node.js packages)
#   - go (for Go modules)
#   - curl (for API checks)
#   - jq (optional, for JSON processing)
#   - yq (optional, for reading versions.yaml - falls back to hardcoded values)
#
# Example:
#   ./scripts/check-latest-versions.sh
#   ./scripts/check-latest-versions.sh --json
#   GITHUB_TOKEN="your_token" ./scripts/check-latest-versions.sh

set -e

# Parse arguments
JSON_OUTPUT=false
QUIET_MODE=false
for arg in "$@"; do
    case $arg in
        --json)
            JSON_OUTPUT=true
            ;;
        --quiet)
            QUIET_MODE=true
            ;;
        --help)
            echo "Usage: $0 [--json] [--quiet]"
            echo ""
            echo "Options:"
            echo "  --json    Output results in JSON format"
            echo "  --quiet   Only show outdated packages"
            echo "  --help    Show this help message"
            exit 0
            ;;
    esac
done

# Colors for output (disabled in JSON mode)
if [ "$JSON_OUTPUT" = false ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    NC=''
fi

# JSON output array
JSON_RESULTS='[]'
VERSIONS_FILE="pkg/config/versions.yaml"

if [ "$JSON_OUTPUT" = false ]; then
    echo -e "${BLUE}=== Checking Latest Versions ===${NC}\n"
fi

# Check if yq is available for reading versions.yaml
if ! command -v yq >/dev/null 2>&1; then
    if [ "$JSON_OUTPUT" = false ]; then
        echo -e "${YELLOW}Warning: yq is not installed. Using fallback hardcoded versions.${NC}"
        echo -e "${YELLOW}Install yq for automatic version reading: brew install yq${NC}\n"
    fi
    USE_YAML=false
else
    USE_YAML=true
fi

# Function to get version from YAML or use fallback
get_version() {
    local yaml_path=$1
    local fallback=$2
    
    if [ "$USE_YAML" = true ] && [ -f "$VERSIONS_FILE" ]; then
        local version=$(yq "$yaml_path" "$VERSIONS_FILE" 2>/dev/null)
        if [ -n "$version" ] && [ "$version" != "null" ]; then
            echo "$version"
            return
        fi
    fi
    
    echo "$fallback"
}

# Function to print section header
print_section() {
    if [ "$JSON_OUTPUT" = false ]; then
        echo -e "\n${BLUE}### $1 ###${NC}"
    fi
}

# Function to compare versions (returns 0 if outdated, 1 if up-to-date)
version_compare() {
    local current=$1
    local latest=$2
    
    # Remove 'v' prefix and '^' prefix if present
    current=$(echo "$current" | sed 's/^[v^]//')
    latest=$(echo "$latest" | sed 's/^[v^]//')
    
    if [ "$current" = "$latest" ]; then
        return 1  # Up-to-date
    fi
    
    # Simple version comparison (works for most cases)
    if [ "$(printf '%s\n' "$current" "$latest" | sort -V | head -n1)" = "$current" ] && [ "$current" != "$latest" ]; then
        return 0  # Outdated
    fi
    
    return 1  # Up-to-date or unable to determine
}

# Function to add JSON result
add_json_result() {
    local category=$1
    local package=$2
    local current=$3
    local latest=$4
    local status=$5
    
    if [ "$JSON_OUTPUT" = true ]; then
        JSON_RESULTS=$(echo "$JSON_RESULTS" | jq -c ". += [{\"category\": \"$category\", \"package\": \"$package\", \"current\": \"$current\", \"latest\": \"$latest\", \"status\": \"$status\"}]")
    fi
}

# Function to check npm package version
check_npm() {
    local package=$1
    local current=$2
    local category=${3:-"npm"}
    
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -n "Checking $package... "
    fi
    
    if command -v npm >/dev/null 2>&1; then
        latest=$(npm view "$package" version 2>/dev/null || echo "N/A")
        if [ "$latest" != "N/A" ]; then
            if version_compare "$current" "$latest"; then
                status="outdated"
                color=$RED
                symbol="⚠️"
            else
                status="up-to-date"
                color=$GREEN
                symbol="✓"
            fi
            
            add_json_result "$category" "$package" "$current" "$latest" "$status"
            
            if [ "$JSON_OUTPUT" = false ]; then
                if [ "$QUIET_MODE" = true ] && [ "$status" = "up-to-date" ]; then
                    return
                fi
                echo -e "${color}${symbol} Latest: $latest${NC} (Current: $current)"
            fi
        else
            add_json_result "$category" "$package" "$current" "N/A" "error"
            if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
                echo -e "${RED}Failed to fetch${NC}"
            fi
        fi
    else
        add_json_result "$category" "$package" "$current" "N/A" "tool-missing"
        if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
            echo -e "${YELLOW}npm not installed${NC}"
        fi
    fi
}

# Function to check Go module version
check_go_module() {
    local module=$1
    local current=$2
    local category=${3:-"go"}
    
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -n "Checking $module... "
    fi
    
    if command -v go >/dev/null 2>&1; then
        latest=$(go list -m -versions "$module" 2>/dev/null | awk '{print $NF}' || echo "N/A")
        if [ "$latest" != "N/A" ] && [ -n "$latest" ]; then
            if version_compare "$current" "$latest"; then
                status="outdated"
                color=$RED
                symbol="⚠️"
            else
                status="up-to-date"
                color=$GREEN
                symbol="✓"
            fi
            
            add_json_result "$category" "$module" "$current" "$latest" "$status"
            
            if [ "$JSON_OUTPUT" = false ]; then
                if [ "$QUIET_MODE" = true ] && [ "$status" = "up-to-date" ]; then
                    return
                fi
                echo -e "${color}${symbol} Latest: $latest${NC} (Current: $current)"
            fi
        else
            add_json_result "$category" "$module" "$current" "N/A" "error"
            if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
                echo -e "${YELLOW}Unable to fetch${NC}"
            fi
        fi
    else
        add_json_result "$category" "$module" "$current" "N/A" "tool-missing"
        if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
            echo -e "${YELLOW}go not installed${NC}"
        fi
    fi
}

# Function to check Docker image tag
check_docker_image() {
    local image=$1
    local current=$2
    local category=${3:-"docker"}
    
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -n "Checking $image... "
    fi
    
    if command -v curl >/dev/null 2>&1; then
        # Try Docker Hub API (works if logged in or within rate limits)
        # Get authentication token if available
        local token=""
        if [ -f ~/.docker/config.json ]; then
            # Try to get token from Docker Hub
            token=$(curl -s "https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/$image:pull" 2>/dev/null | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
        fi
        
        # Query Docker Hub API for tags
        if [ -n "$token" ]; then
            latest=$(curl -s -H "Authorization: Bearer $token" "https://registry.hub.docker.com/v2/library/$image/tags/list" 2>/dev/null | grep -o '"latest"' | head -1)
            if [ -n "$latest" ]; then
                # Get the actual version tag (not just "latest")
                # For Alpine, Ubuntu, etc., we need to check their specific versioning
                case $image in
                    alpine)
                        # Alpine uses version numbers like 3.19, 3.20
                        latest=$(curl -s -H "Authorization: Bearer $token" "https://registry.hub.docker.com/v2/library/$image/tags/list" 2>/dev/null | grep -o '"3\.[0-9]*"' | sort -V | tail -1 | tr -d '"')
                        ;;
                    ubuntu)
                        # Ubuntu uses version numbers like 24.04, 22.04
                        latest=$(curl -s -H "Authorization: Bearer $token" "https://registry.hub.docker.com/v2/library/$image/tags/list" 2>/dev/null | grep -o '"[0-9]*\.[0-9]*"' | sort -V | tail -1 | tr -d '"')
                        ;;
                    golang)
                        # Golang uses version like 1.25-alpine
                        latest=$(curl -s -H "Authorization: Bearer $token" "https://registry.hub.docker.com/v2/library/$image/tags/list" 2>/dev/null | grep -o '"1\.[0-9]*-alpine"' | sort -V | tail -1 | tr -d '"')
                        ;;
                esac
            fi
        else
            # Fallback: Try without authentication (may hit rate limits)
            case $image in
                alpine)
                    latest=$(curl -s "https://registry.hub.docker.com/v2/repositories/library/$image/tags?page_size=100" 2>/dev/null | grep -o '"name":"3\.[0-9]*"' | grep -v 'rc\|beta\|alpha' | sort -V | tail -1 | cut -d'"' -f4)
                    ;;
                ubuntu)
                    latest=$(curl -s "https://registry.hub.docker.com/v2/repositories/library/$image/tags?page_size=100" 2>/dev/null | grep -o '"name":"[0-9]*\.[0-9]*"' | grep -v 'rc\|beta\|alpha' | sort -V | tail -1 | cut -d'"' -f4)
                    ;;
                golang)
                    latest=$(curl -s "https://registry.hub.docker.com/v2/repositories/library/$image/tags?page_size=100" 2>/dev/null | grep -o '"name":"1\.[0-9]*-alpine"' | sort -V | tail -1 | cut -d'"' -f4)
                    ;;
            esac
        fi
        
        if [ -n "$latest" ] && [ "$latest" != "N/A" ]; then
            if version_compare "$current" "$latest"; then
                status="outdated"
                color=$RED
                symbol="⚠️"
            else
                status="up-to-date"
                color=$GREEN
                symbol="✓"
            fi
            
            add_json_result "$category" "$image" "$current" "$latest" "$status"
            
            if [ "$JSON_OUTPUT" = false ]; then
                if [ "$QUIET_MODE" = true ] && [ "$status" = "up-to-date" ]; then
                    return
                fi
                echo -e "${color}${symbol} Latest: $latest${NC} (Current: $current)"
            fi
        else
            add_json_result "$category" "$image" "$current" "N/A" "error"
            if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
                echo -e "${YELLOW}Unable to fetch (check manually: https://hub.docker.com/_/$image)${NC}"
            fi
        fi
    else
        add_json_result "$category" "$image" "$current" "N/A" "tool-missing"
        if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
            echo -e "${YELLOW}curl not installed${NC}"
        fi
    fi
}

# Node.js & Frontend
print_section "Node.js & Frontend"
check_npm "create-next-app" "$(get_version '.frontend.nextjs.version' '16.0.0')" "frontend"
check_npm "next" "$(get_version '.frontend.nextjs.version' '16.0.0')" "frontend"
check_npm "react" "$(get_version '.frontend.react.version' '19.2.0')" "frontend"
check_npm "react-dom" "$(get_version '.frontend.react_dom.version' '19.2.0')" "frontend"
check_npm "typescript" "$(get_version '.frontend.typescript.version' '5.9.3')" "frontend"

# Go Backend
print_section "Go Backend"

# Check Go version from Docker Hub golang image
current_go_version="$(get_version '.backend.go.version' '1.25.0')"
if command -v curl >/dev/null 2>&1; then
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -n "Checking Go... "
    fi
    
    # Get latest golang docker image tag
    latest_docker_tag=$(curl -s "https://registry.hub.docker.com/v2/repositories/library/golang/tags?page_size=100" 2>/dev/null | grep -o '"name":"[0-9]*\.[0-9]*-alpine"' | head -1 | sed 's/"name":"//;s/-alpine"//')
    
    if [ -n "$latest_docker_tag" ]; then
        latest_go_version="${latest_docker_tag}.0"
        
        if version_compare "$current_go_version" "$latest_go_version"; then
            status="outdated"
            add_json_result "backend" "go" "$current_go_version" "$latest_go_version" "$status"
            if [ "$JSON_OUTPUT" = false ]; then
                if [ "$QUIET_MODE" = true ] && [ "$status" = "up-to-date" ]; then
                    return
                fi
                echo -e "${RED}⚠️  Latest: $latest_go_version (Current: $current_go_version)${NC}"
            fi
        else
            status="up-to-date"
            add_json_result "backend" "go" "$current_go_version" "$latest_go_version" "$status"
            if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
                echo -e "${GREEN}✓ Latest: $latest_go_version (Current: $current_go_version)${NC}"
            fi
        fi
    else
        add_json_result "backend" "go" "$current_go_version" "N/A" "error"
        if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
            echo -e "${YELLOW}Unable to fetch${NC}"
        fi
    fi
else
    add_json_result "backend" "go" "$current_go_version" "N/A" "tool-missing"
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -e "${YELLOW}curl not installed${NC}"
    fi
fi

# Go Backend Frameworks
print_section "Go Backend Frameworks"
check_go_module "github.com/gin-gonic/gin" "$(get_version '.backend.frameworks.gin.version' 'v1.11.0')" "backend"
check_go_module "github.com/gin-contrib/cors" "$(get_version '.backend.frameworks.gin_cors.version' 'v1.7.6')" "backend"
check_go_module "github.com/labstack/echo/v4" "$(get_version '.backend.frameworks.echo.version' 'v4.13.4')" "backend"
check_go_module "github.com/gofiber/fiber/v2" "$(get_version '.backend.frameworks.fiber.version' 'v2.52.9')" "backend"

# Function to check GitHub releases
check_github_release() {
    local repo=$1
    local current=$2
    local category=${3:-"github"}
    local package_name=${4:-$repo}
    
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -n "Checking $package_name... "
    fi
    
    if command -v curl >/dev/null 2>&1; then
        # Build auth header if GITHUB_TOKEN is set
        local auth_header=""
        if [ -n "$GITHUB_TOKEN" ]; then
            auth_header="-H \"Authorization: Bearer $GITHUB_TOKEN\""
        fi
        
        # Try to get latest release from GitHub API (with timeout and follow redirects)
        latest=$(eval curl -s -L -m 5 $auth_header "https://api.github.com/repos/$repo/releases/latest" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
        
        if [ -z "$latest" ] || [ "$latest" = "N/A" ]; then
            # Fallback: try tags endpoint
            latest=$(eval curl -s -L -m 5 $auth_header "https://api.github.com/repos/$repo/tags" 2>/dev/null | grep '"name":' | head -1 | sed -E 's/.*"name": "([^"]+)".*/\1/')
        fi
        
        if [ "$latest" != "N/A" ] && [ -n "$latest" ]; then
            if version_compare "$current" "$latest"; then
                status="outdated"
                color=$RED
                symbol="⚠️"
            else
                status="up-to-date"
                color=$GREEN
                symbol="✓"
            fi
            
            add_json_result "$category" "$package_name" "$current" "$latest" "$status"
            
            if [ "$JSON_OUTPUT" = false ]; then
                if [ "$QUIET_MODE" = true ] && [ "$status" = "up-to-date" ]; then
                    return
                fi
                echo -e "${color}${symbol} Latest: $latest${NC} (Current: $current)"
            fi
        else
            add_json_result "$category" "$package_name" "$current" "N/A" "error"
            if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
                echo -e "${YELLOW}Unable to fetch${NC}"
            fi
        fi
    else
        add_json_result "$category" "$package_name" "$current" "N/A" "tool-missing"
        if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
            echo -e "${YELLOW}curl not installed${NC}"
        fi
    fi
}

# Function to check Maven Central (for Android libraries)
check_maven() {
    local group=$1
    local artifact=$2
    local current=$3
    local category=${4:-"android"}
    local display_name=${5:-"$artifact"}
    local package_name="$group:$artifact"
    
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -n "Checking $display_name... "
    fi
    
    if command -v curl >/dev/null 2>&1; then
        group_path=$(echo "$group" | tr '.' '/')
        
        # Get metadata and find latest stable version
        metadata=$(curl -s -m 5 "https://dl.google.com/android/maven2/$group_path/$artifact/maven-metadata.xml" 2>/dev/null)
        
        # Method 1: Try release tag (but verify it's stable)
        latest=$(echo "$metadata" | grep -o '<release>[^<]*</release>' | sed 's/<[^>]*>//g')
        
        # If release contains alpha/beta/rc, get latest stable from version list
        if [ -n "$latest" ] && echo "$latest" | grep -qE "(alpha|beta|rc)"; then
            latest=$(echo "$metadata" | grep '<version>' | grep -v -E "(alpha|beta|rc)" | tail -1 | sed 's/.*<version>\([^<]*\)<\/version>.*/\1/')
        fi
        
        # Method 2: If no release tag, try latest tag
        if [ -z "$latest" ]; then
            latest=$(echo "$metadata" | grep -o '<latest>[^<]*</latest>' | sed 's/<[^>]*>//g')
            # Filter out pre-release versions
            if echo "$latest" | grep -qE "(alpha|beta|rc)"; then
                latest=$(echo "$metadata" | grep '<version>' | grep -v -E "(alpha|beta|rc)" | tail -1 | sed 's/.*<version>\([^<]*\)<\/version>.*/\1/')
            fi
        fi
        
        # Method 3: Fallback to Maven Central
        if [ -z "$latest" ] || [ "$latest" = "N/A" ]; then
            latest=$(curl -s -m 5 "https://repo1.maven.org/maven2/$group_path/$artifact/maven-metadata.xml" 2>/dev/null | grep -o '<latest>[^<]*</latest>' | sed 's/<[^>]*>//g')
        fi
        
        if [ -n "$latest" ] && [ "$latest" != "N/A" ]; then
            if version_compare "$current" "$latest"; then
                status="outdated"
                color=$RED
                symbol="⚠️"
            else
                status="up-to-date"
                color=$GREEN
                symbol="✓"
            fi
            
            add_json_result "$category" "$display_name" "$current" "$latest" "$status"
            
            if [ "$JSON_OUTPUT" = false ]; then
                if [ "$QUIET_MODE" = true ] && [ "$status" = "up-to-date" ]; then
                    return
                fi
                echo -e "${color}${symbol} Latest: $latest${NC} (Current: $current)"
            fi
        else
            add_json_result "$category" "$display_name" "$current" "N/A" "error"
            if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
                echo -e "${YELLOW}Unable to fetch (check manually)${NC}"
            fi
        fi
    else
        add_json_result "$category" "$display_name" "$current" "N/A" "tool-missing"
        if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
            echo -e "${YELLOW}curl not installed${NC}"
        fi
    fi
}

# Function to check Gradle version
check_gradle() {
    local current=$1
    local category=${2:-"android"}
    
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -n "Checking Gradle... "
    fi
    
    if command -v curl >/dev/null 2>&1; then
        # Get latest Gradle version from services API
        latest=$(curl -s -m 10 "https://services.gradle.org/versions/current" 2>/dev/null | grep -o '"version" : "[^"]*"' | head -1 | cut -d'"' -f4)
        
        if [ -n "$latest" ] && [ "$latest" != "N/A" ]; then
            if version_compare "$current" "$latest"; then
                status="outdated"
                color=$RED
                symbol="⚠️"
            else
                status="up-to-date"
                color=$GREEN
                symbol="✓"
            fi
            
            add_json_result "$category" "gradle" "$current" "$latest" "$status"
            
            if [ "$JSON_OUTPUT" = false ]; then
                if [ "$QUIET_MODE" = true ] && [ "$status" = "up-to-date" ]; then
                    return
                fi
                echo -e "${color}${symbol} Latest: $latest${NC} (Current: $current)"
            fi
        else
            add_json_result "$category" "gradle" "$current" "N/A" "error"
            if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
                echo -e "${YELLOW}Unable to fetch (check manually)${NC}"
            fi
        fi
    else
        add_json_result "$category" "gradle" "$current" "N/A" "tool-missing"
        if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
            echo -e "${YELLOW}curl not installed${NC}"
        fi
    fi
}

# Android Dependencies (Maven Central)
print_section "Android Dependencies"
check_maven "androidx.core" "core-ktx" "$(get_version '.android.androidx.core_ktx.version' '1.17.0')" "android"
check_maven "androidx.appcompat" "appcompat" "$(get_version '.android.androidx.appcompat.version' '1.7.1')" "android"
check_maven "com.google.android.material" "material" "$(get_version '.android.androidx.material.version' '1.13.0')" "android"
check_maven "androidx.constraintlayout" "constraintlayout" "$(get_version '.android.androidx.constraintlayout.version' '2.2.1')" "android"

# Android Testing Dependencies
print_section "Android Testing Dependencies"
check_maven "junit" "junit" "$(get_version '.android.testing.junit.version' '4.13.2')" "android" "JUnit"
check_maven "androidx.test.ext" "junit" "$(get_version '.android.testing.androidx_junit.version' '1.2.1')" "android" "AndroidX JUnit"
check_maven "androidx.test.espresso" "espresso-core" "$(get_version '.android.testing.espresso.version' '3.6.1')" "android" "Espresso"

# Function to check Android SDK API level
check_android_sdk() {
    local current=$1
    local category=${2:-"android"}
    
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -n "Checking Android SDK... "
    fi
    
    if command -v curl >/dev/null 2>&1; then
        # Use endoflife.date API for Android versions (get apiVersion, not cycle)
        latest=$(curl -s -m 5 "https://endoflife.date/api/android.json" 2>/dev/null | grep -o '"apiVersion":"[0-9]*"' | head -1 | grep -o "[0-9]*")
        
        if [ -n "$latest" ] && [ "$latest" != "N/A" ]; then
            if version_compare "$current" "$latest"; then
                status="outdated"
                color=$RED
                symbol="⚠️"
            else
                status="up-to-date"
                color=$GREEN
                symbol="✓"
            fi
            
            add_json_result "$category" "android-sdk" "$current" "$latest" "$status"
            
            if [ "$JSON_OUTPUT" = false ]; then
                if [ "$QUIET_MODE" = true ] && [ "$status" = "up-to-date" ]; then
                    return
                fi
                echo -e "${color}${symbol} Latest API: $latest${NC} (Current: $current)"
            fi
        else
            add_json_result "$category" "android-sdk" "$current" "N/A" "error"
            if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
                echo -e "${YELLOW}Unable to fetch${NC}"
            fi
        fi
    else
        add_json_result "$category" "android-sdk" "$current" "N/A" "tool-missing"
        if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
            echo -e "${YELLOW}curl not installed${NC}"
        fi
    fi
}

# Function to check Xcode version
check_xcode() {
    local current=$1
    local category=${2:-"ios"}
    
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -n "Checking Xcode... "
    fi
    
    if command -v curl >/dev/null 2>&1; then
        # Scrape Apple Developer releases page for latest stable Xcode
        latest=$(curl -s -m 5 "https://developer.apple.com/news/releases/" 2>/dev/null | grep -o 'Xcode [0-9]*\.[0-9]*\.[0-9]* ([0-9]*[A-Z][0-9]*)' | grep -v beta | head -1 | grep -o '[0-9]*\.[0-9]*\.[0-9]*' | head -1)
        
        if [ -n "$latest" ] && [ "$latest" != "N/A" ]; then
            if version_compare "$current" "$latest"; then
                status="outdated"
                color=$RED
                symbol="⚠️"
            else
                status="up-to-date"
                color=$GREEN
                symbol="✓"
            fi
            
            add_json_result "$category" "xcode" "$current" "$latest" "$status"
            
            if [ "$JSON_OUTPUT" = false ]; then
                if [ "$QUIET_MODE" = true ] && [ "$status" = "up-to-date" ]; then
                    return
                fi
                echo -e "${color}${symbol} Latest: $latest${NC} (Current: $current)"
            fi
        else
            add_json_result "$category" "xcode" "$current" "N/A" "error"
            if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
                echo -e "${YELLOW}Unable to fetch${NC}"
            fi
        fi
    else
        add_json_result "$category" "xcode" "$current" "N/A" "tool-missing"
        if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
            echo -e "${YELLOW}curl not installed${NC}"
        fi
    fi
}

# Android Build Tools
print_section "Android Build Tools"
check_gradle "$(get_version '.android.gradle.version' '9.1.0')" "android"
check_maven "com.android.tools.build" "gradle" "$(get_version '.android.gradle_plugin.version' '8.7.3')" "android" "Android Gradle Plugin"
check_github_release "JetBrains/kotlin" "$(get_version '.android.kotlin.version' '2.2.21')" "android" "Kotlin"
check_android_sdk "$(get_version '.android.compile_sdk' '36')" "android"

# Function to check iOS version
check_ios_version() {
    local current=$1
    local category=${2:-"ios"}
    
    if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
        echo -n "Checking iOS Deployment Target... "
    fi
    
    if command -v curl >/dev/null 2>&1; then
        # Use endoflife.date API for iOS versions
        latest=$(curl -s -m 5 "https://endoflife.date/api/ios.json" 2>/dev/null | grep -o '"cycle":"[0-9]*"' | head -1 | grep -o "[0-9]*")
        
        if [ -n "$latest" ] && [ "$latest" != "N/A" ]; then
            # Add .0 to make it a proper version
            latest="${latest}.0"
            
            if version_compare "$current" "$latest"; then
                status="outdated"
                color=$RED
                symbol="⚠️"
            else
                status="up-to-date"
                color=$GREEN
                symbol="✓"
            fi
            
            add_json_result "$category" "ios-deployment-target" "$current" "$latest" "$status"
            
            if [ "$JSON_OUTPUT" = false ]; then
                if [ "$QUIET_MODE" = true ] && [ "$status" = "up-to-date" ]; then
                    return
                fi
                echo -e "${color}${symbol} Latest: $latest${NC} (Current: $current)"
            fi
        else
            add_json_result "$category" "ios-deployment-target" "$current" "N/A" "error"
            if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
                echo -e "${YELLOW}Unable to fetch${NC}"
            fi
        fi
    else
        add_json_result "$category" "ios-deployment-target" "$current" "N/A" "tool-missing"
        if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
            echo -e "${YELLOW}curl not installed${NC}"
        fi
    fi
}

# iOS
print_section "iOS"
check_github_release "apple/swift" "$(get_version '.ios.swift.version' 'swift-6.2-RELEASE')" "ios" "Swift"
check_xcode "$(get_version '.ios.xcode.version' '26.0.1')" "ios"
check_ios_version "$(get_version '.ios.deployment_target' '26.0')" "ios"

# Docker Base Images
print_section "Docker Base Images"
check_docker_image "alpine" "$(get_version '.docker.alpine.version' '3.22')" "docker"
check_docker_image "golang" "$(get_version '.docker.golang.version' '1.25-alpine')" "docker"
check_docker_image "ubuntu" "$(get_version '.docker.ubuntu.version' '25.10')" "docker"

# Go Toolchain
print_section "Go Toolchain"
if [ "$JSON_OUTPUT" = false ]; then
    if command -v go >/dev/null 2>&1; then
        installed=$(go version | awk '{print $3}' | sed 's/go//')
        echo -e "Installed: ${GREEN}$installed${NC} (Required: 1.25.0)"
        echo "Latest: Check https://go.dev/dl/"
    else
        echo -e "${YELLOW}go not installed${NC}"
    fi
fi

# Infrastructure Tools
print_section "Infrastructure Tools"
check_github_release "hashicorp/terraform" "$(get_version '.infrastructure.terraform.version' 'v1.13.4')" "infrastructure" "Terraform"
check_github_release "kubernetes/kubernetes" "$(get_version '.infrastructure.kubernetes.version' 'v1.34.1')" "infrastructure" "Kubernetes"

# Development Tools
print_section "Go Development Tools"
if [ "$JSON_OUTPUT" = false ] && [ "$QUIET_MODE" = false ]; then
    echo "Check latest versions with:"
    echo "  go list -m -versions github.com/air-verse/air"
    echo "  go list -m -versions gotest.tools/gotestsum"
    echo "  go list -m -versions github.com/securego/gosec/v2"
    echo "  go list -m -versions golang.org/x/vuln"
    echo "  go list -m -versions honnef.co/go/tools"
fi

if [ "$JSON_OUTPUT" = true ]; then
    # Output JSON results
    echo "$JSON_RESULTS" | jq '.'
else
    echo -e "\n${BLUE}=== Summary ===${NC}"
    
    # Count outdated packages
    if command -v jq >/dev/null 2>&1; then
        outdated_count=$(echo "$JSON_RESULTS" | jq '[.[] | select(.status == "outdated")] | length' 2>/dev/null || echo "0")
        uptodate_count=$(echo "$JSON_RESULTS" | jq '[.[] | select(.status == "up-to-date")] | length' 2>/dev/null || echo "0")
        
        if [ "$outdated_count" -gt 0 ]; then
            echo -e "${RED}⚠️  $outdated_count package(s) outdated${NC}"
        fi
        if [ "$uptodate_count" -gt 0 ]; then
            echo -e "${GREEN}✓ $uptodate_count package(s) up-to-date${NC}"
        fi
    fi
    
    echo ""
    echo "This script checks versions where possible via CLI tools."
    echo "For packages requiring manual checks, URLs are provided above."
    echo ""
    echo "To update versions:"
    echo "  1. Review the output above"
    echo "  2. Update hardcoded values in source files"
    echo "  3. Run 'make test' to verify changes"
    echo "  4. Update CHANGELOG.md"
    echo ""
    echo "Options:"
    echo "  --json    Output results in JSON format"
    echo "  --quiet   Only show outdated packages"
    
    echo -e "\n${YELLOW}Note: Some checks require internet connectivity${NC}"
fi
