#!/bin/bash
# Script to update versions in configs/versions.yaml
# Calls check-latest-versions.sh to get version data and updates outdated packages
#
# Usage: ./scripts/update-versions.sh [OPTIONS]
#
# Options:
#   --dry-run      Show what would be updated without making changes
#   --auto-update  Automatically update versions.yaml without prompting
#   --help         Show this help message
#
# Environment Variables:
#   GITHUB_TOKEN  Optional GitHub token for higher rate limits (5000/hour vs 60/hour)
#
# Requirements:
#   - yq (YAML processor) - Install: brew install yq (macOS) or snap install yq (Linux)
#   - jq (JSON processor)
#   - All requirements from check-latest-versions.sh
#
# Example:
#   ./scripts/update-versions.sh
#   ./scripts/update-versions.sh --dry-run
#   ./scripts/update-versions.sh --auto-update

set -e

# Parse arguments
DRY_RUN=false
AUTO_UPDATE=false
for arg in "$@"; do
    case $arg in
        --dry-run)
            DRY_RUN=true
            ;;
        --auto-update)
            AUTO_UPDATE=true
            ;;
        --help)
            echo "Usage: $0 [--dry-run] [--auto-update]"
            echo ""
            echo "Options:"
            echo "  --dry-run      Show what would be updated without making changes"
            echo "  --auto-update  Automatically update versions.yaml without prompting"
            echo "  --help         Show this help message"
            exit 0
            ;;
    esac
done

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

VERSIONS_FILE="configs/versions.yaml"
TEMP_FILE=$(mktemp)

# Cleanup on exit
cleanup() {
    if [ -f "$TEMP_FILE" ]; then
        rm -f "$TEMP_FILE"
    fi
}
trap cleanup EXIT

echo -e "${BLUE}=== Updating Versions from Latest ===${NC}\n"

# Check if versions.yaml exists
if [ ! -f "$VERSIONS_FILE" ]; then
    echo -e "${RED}Error: $VERSIONS_FILE not found${NC}"
    echo "Please run this script from the project root directory."
    exit 1
fi

# Check required tools
MISSING_TOOLS=()

if ! command -v yq >/dev/null 2>&1; then
    MISSING_TOOLS+=("yq")
fi

if ! command -v jq >/dev/null 2>&1; then
    MISSING_TOOLS+=("jq")
fi

if [ ${#MISSING_TOOLS[@]} -gt 0 ]; then
    echo -e "${RED}Error: Missing required tools: ${MISSING_TOOLS[*]}${NC}"
    echo ""
    echo "Install missing tools:"
    for tool in "${MISSING_TOOLS[@]}"; do
        case "$tool" in
            yq)
                echo "  yq:  brew install yq (macOS) or snap install yq (Linux)"
                echo "       Or download from: https://github.com/mikefarah/yq/releases"
                ;;
            jq)
                echo "  jq:  brew install jq (macOS) or apt install jq (Linux)"
                ;;
        esac
    done
    echo ""
    echo "Alternatively, manually update configs/versions.yaml with versions from:"
    echo "  make check-versions"
    exit 1
fi

# Get version data from check-latest-versions.sh in JSON format
echo "Fetching latest versions..."
VERSION_DATA=$(./scripts/check-latest-versions.sh --json 2>/dev/null)

if [ -z "$VERSION_DATA" ] || [ "$VERSION_DATA" = "[]" ]; then
    echo -e "${RED}Error: Failed to fetch version data${NC}"
    exit 1
fi

# Copy current versions file to temp
cp "$VERSIONS_FILE" "$TEMP_FILE"

# Parse JSON and update YAML
UPDATES_MADE=false

# Check if there are any outdated packages
OUTDATED_COUNT=$(echo "$VERSION_DATA" | jq '[.[] | select(.status == "outdated")] | length')

if [ "$OUTDATED_COUNT" -eq 0 ]; then
    echo -e "${GREEN}✓ All versions are up-to-date!${NC}"
    rm "$TEMP_FILE"
    exit 0
fi

echo -e "${YELLOW}Found $OUTDATED_COUNT outdated package(s)${NC}\n"

# Update each version from the JSON data
while read -r item; do
    package=$(echo "$item" | jq -r '.package')
    latest=$(echo "$item" | jq -r '.latest')
    
    # Map package names to YAML paths
    case "$package" in
            "create-next-app")
                yq -i ".frontend.nextjs.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated create-next-app → $latest${NC}"
                ;;
            "next")
                # Next.js version is same as create-next-app, skip
                ;;
            "react")
                yq -i ".frontend.react.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated react → $latest${NC}"
                ;;
            "react-dom")
                yq -i ".frontend.react_dom.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated react-dom → $latest${NC}"
                ;;
            "typescript")
                yq -i ".frontend.typescript.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated typescript → $latest${NC}"
                ;;
            "go")
                yq -i ".backend.go.version = \"$latest\"" "$TEMP_FILE"
                # Update docker_tag to match (e.g., "1.25.0" -> "1.25-alpine")
                docker_tag=$(echo "$latest" | sed -E 's/\.0$//')-alpine
                yq -i ".backend.go.docker_tag = \"$docker_tag\"" "$TEMP_FILE"
                # Also update docker.golang.version
                yq -i ".docker.golang.version = \"$docker_tag\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Go → $latest (docker: $docker_tag)${NC}"
                ;;
            "github.com/gin-gonic/gin")
                yq -i ".backend.frameworks.gin.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated gin → $latest${NC}"
                ;;
            "github.com/gin-contrib/cors")
                yq -i ".backend.frameworks.gin_cors.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated gin-cors → $latest${NC}"
                ;;
            "github.com/labstack/echo/v4")
                yq -i ".backend.frameworks.echo.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated echo → $latest${NC}"
                ;;
            "github.com/gofiber/fiber/v2")
                yq -i ".backend.frameworks.fiber.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated fiber → $latest${NC}"
                ;;
            "core-ktx")
                yq -i ".android.androidx.core_ktx.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated core-ktx → $latest${NC}"
                ;;
            "appcompat")
                yq -i ".android.androidx.appcompat.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated appcompat → $latest${NC}"
                ;;
            "material")
                yq -i ".android.androidx.material.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated material → $latest${NC}"
                ;;
            "constraintlayout")
                yq -i ".android.androidx.constraintlayout.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated constraintlayout → $latest${NC}"
                ;;
            "Kotlin")
                yq -i ".android.kotlin.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Kotlin → $latest${NC}"
                ;;
            "gradle")
                yq -i ".android.gradle.version = \"$latest\"" "$TEMP_FILE"
                yq -i ".android.gradle.distribution_url = \"https://services.gradle.org/distributions/gradle-$latest-bin.zip\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Gradle → $latest${NC}"
                ;;
            "Android Gradle Plugin")
                yq -i ".android.gradle_plugin.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Android Gradle Plugin → $latest${NC}"
                ;;
            "JUnit")
                yq -i ".android.testing.junit.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated JUnit → $latest${NC}"
                ;;
            "AndroidX JUnit")
                yq -i ".android.testing.androidx_junit.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated AndroidX JUnit → $latest${NC}"
                ;;
            "Espresso")
                yq -i ".android.testing.espresso.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Espresso → $latest${NC}"
                ;;
            "android-sdk")
                yq -i ".android.compile_sdk = \"$latest\"" "$TEMP_FILE"
                yq -i ".android.target_sdk = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Android SDK → $latest${NC}"
                ;;
            "Swift")
                yq -i ".ios.swift.version = \"$latest\"" "$TEMP_FILE"
                # Extract short version (e.g., "swift-6.2-RELEASE" -> "6.2")
                short_version=$(echo "$latest" | sed -E 's/swift-([0-9]+\.[0-9]+).*/\1/')
                yq -i ".ios.swift.short_version = \"$short_version\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Swift → $latest (short: $short_version)${NC}"
                ;;
            "xcode")
                yq -i ".ios.xcode.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Xcode → $latest${NC}"
                ;;
            "ios-deployment-target")
                yq -i ".ios.deployment_target = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated iOS deployment target → $latest${NC}"
                ;;
            "alpine")
                yq -i ".docker.alpine.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Alpine → $latest${NC}"
                ;;
            "golang")
                yq -i ".docker.golang.version = \"$latest\"" "$TEMP_FILE"
                # Extract Go version from docker tag (e.g., "1.25-alpine" -> "1.25.0")
                go_version=$(echo "$latest" | sed -E 's/([0-9]+\.[0-9]+).*/\1.0/')
                yq -i ".backend.go.version = \"$go_version\"" "$TEMP_FILE"
                yq -i ".backend.go.docker_tag = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Golang Docker image → $latest (Go: $go_version)${NC}"
                ;;
            "ubuntu")
                yq -i ".docker.ubuntu.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Ubuntu → $latest${NC}"
                ;;
            "Terraform")
                yq -i ".infrastructure.terraform.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Terraform → $latest${NC}"
                ;;
            "Kubernetes")
                yq -i ".infrastructure.kubernetes.version = \"$latest\"" "$TEMP_FILE"
                echo -e "${YELLOW}📦 Updated Kubernetes → $latest${NC}"
                ;;
            *)
                echo -e "${YELLOW}⚠ Unknown package: $package (skipping)${NC}"
                ;;
        esac
done < <(echo "$VERSION_DATA" | jq -r '.[] | select(.status == "outdated") | @json')

# Update last_updated timestamp
yq -i ".metadata.last_updated = \"$(date +%Y-%m-%d)\"" "$TEMP_FILE"

# Apply the updates
if ! diff -q "$VERSIONS_FILE" "$TEMP_FILE" >/dev/null 2>&1; then

    echo -e "\n${BLUE}=== Summary ===${NC}\n"
    
    if [ "$DRY_RUN" = true ]; then
        echo -e "${YELLOW}Dry run mode - no changes made${NC}"
        echo -e "\nProposed changes:"
        diff "$VERSIONS_FILE" "$TEMP_FILE" || true
        rm "$TEMP_FILE"
    elif [ "$AUTO_UPDATE" = true ]; then
        mv "$TEMP_FILE" "$VERSIONS_FILE"
        echo -e "${GREEN}✓ Updated $VERSIONS_FILE${NC}"
    else
        echo -e "${YELLOW}Updates available. Apply changes? (y/n)${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            mv "$TEMP_FILE" "$VERSIONS_FILE"
            echo -e "${GREEN}✓ Updated $VERSIONS_FILE${NC}"
        else
            echo -e "${YELLOW}Changes not applied${NC}"
            rm "$TEMP_FILE"
        fi
    fi
else
    # No changes detected (shouldn't happen since we checked OUTDATED_COUNT)
    echo -e "${YELLOW}No changes to apply${NC}"
    rm "$TEMP_FILE"
fi
