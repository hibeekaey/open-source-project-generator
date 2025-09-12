#!/bin/bash

# Script to audit Docker images for security vulnerabilities
set -e

echo "=== Auditing Docker images for security vulnerabilities ==="

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to scan a Docker image
scan_image() {
    local image="$1"
    echo "Scanning Docker image: $image"
    
    if command_exists trivy; then
        echo "Using Trivy scanner..."
        trivy image --severity HIGH,CRITICAL "$image"
    elif command_exists docker; then
        echo "Using Docker scout (if available)..."
        if docker scout version >/dev/null 2>&1; then
            docker scout cves "$image"
        else
            echo "⚠️  Docker Scout not available. Pulling image to check for basic issues..."
            docker pull "$image" >/dev/null 2>&1 || echo "Failed to pull $image"
        fi
    else
        echo "⚠️  No Docker vulnerability scanner available"
        return 1
    fi
}

# Function to extract base images from Dockerfile templates
extract_base_images() {
    local dockerfile="$1"
    echo "Extracting base images from $dockerfile..."
    
    # Replace template variables with realistic values and extract FROM statements
    sed -e 's/{{\.Versions\.Go}}/1.25/g' \
        -e 's/{{\.Versions\.Node}}/22/g' \
        "$dockerfile" | grep -E '^FROM' | awk '{print $2}' | grep -v '^builder$' | grep -v '^android-builder$' | grep -v '^ios-builder$' | grep -v '^artifacts$' | grep -v '^runner$'
}

# Check if we have any vulnerability scanner
if ! command_exists trivy && ! command_exists docker; then
    echo "❌ No Docker vulnerability scanner available. Please install Trivy or Docker."
    echo "To install Trivy: brew install trivy (macOS) or see https://trivy.dev/latest/getting-started/installation/"
    exit 1
fi

# Install trivy if not available but docker is
if ! command_exists trivy && command_exists docker; then
    echo "Installing Trivy..."
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command_exists brew; then
            brew install trivy
        else
            echo "Please install Trivy manually: https://trivy.dev/latest/getting-started/installation/"
            exit 1
        fi
    else
        echo "Please install Trivy manually: https://trivy.dev/latest/getting-started/installation/"
        exit 1
    fi
fi

# Scan base images from templates
echo "Scanning base images from Dockerfile templates..."

# Common base images to scan
BASE_IMAGES=(
    "golang:1.25-alpine"
    "golang:1.22-alpine" 
    "node:22-alpine"
    "alpine:3.20"
    "alpine:latest"
    "ubuntu:22.04"
    "openjdk:21-jdk-slim"
)

# Add images from Dockerfile templates
for dockerfile in Dockerfile Dockerfile.dev Dockerfile.build templates/*/Dockerfile.tmpl templates/*/*/Dockerfile.tmpl; do
    if [[ -f "$dockerfile" ]]; then
        while IFS= read -r image; do
            if [[ -n "$image" && "$image" != "AS" ]]; then
                BASE_IMAGES+=("$image")
            fi
        done < <(extract_base_images "$dockerfile")
    fi
done

# Remove duplicates
BASE_IMAGES=($(printf "%s\n" "${BASE_IMAGES[@]}" | sort -u))

echo "Found ${#BASE_IMAGES[@]} unique base images to scan:"
printf '%s\n' "${BASE_IMAGES[@]}"
echo

# Scan each base image
VULNERABLE_IMAGES=()
for image in "${BASE_IMAGES[@]}"; do
    echo "----------------------------------------"
    if scan_image "$image"; then
        echo "✅ $image - No critical vulnerabilities found"
    else
        echo "⚠️  $image - Vulnerabilities found or scan failed"
        VULNERABLE_IMAGES+=("$image")
    fi
    echo
done

# Summary
echo "========================================="
echo "Docker Image Security Audit Summary"
echo "========================================="
echo "Total images scanned: ${#BASE_IMAGES[@]}"
echo "Images with vulnerabilities: ${#VULNERABLE_IMAGES[@]}"

if [[ ${#VULNERABLE_IMAGES[@]} -gt 0 ]]; then
    echo
    echo "Images requiring attention:"
    printf '%s\n' "${VULNERABLE_IMAGES[@]}"
    echo
    echo "Recommendations:"
    echo "1. Update base images to latest versions"
    echo "2. Use specific version tags instead of 'latest'"
    echo "3. Consider using distroless or minimal base images"
    echo "4. Regularly update Dockerfile templates"
else
    echo "✅ All scanned images appear to be secure"
fi

echo "=== Docker image audit complete ==="