#!/bin/bash

# CiliKube Docker Build Script
# Usage: ./build.sh [version] [environment]

set -e

# Default parameters
VERSION=${1:-"dev"}
ENVIRONMENT=${2:-"production"}
IMAGE_NAME="cilikube"
REGISTRY=${REGISTRY:-"cilliantech"}

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed or not in PATH"
        exit 1
    fi
    log_info "Docker version: $(docker --version)"
}

# Get build information
get_build_info() {
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    
    log_info "Build information:"
    log_info "  Version: $VERSION"
    log_info "  Environment: $ENVIRONMENT"
    log_info "  Build time: $BUILD_TIME"
    log_info "  Git commit: $GIT_COMMIT"
}

# Build Docker image
build_image() {
    local dockerfile_path="../../Dockerfile"
    local context_path="../.."
    
    log_info "Starting Docker image build..."
    
    # Build arguments
    local build_args=(
        "--build-arg" "VERSION=$VERSION"
        "--build-arg" "BUILD_TIME=$BUILD_TIME"
        "--build-arg" "GIT_COMMIT=$GIT_COMMIT"
        "--tag" "$REGISTRY/$IMAGE_NAME:$VERSION"
        "--tag" "$REGISTRY/$IMAGE_NAME:latest"
        "--file" "$dockerfile_path"
        "$context_path"
    )
    
    # Add extra optimizations for production environment
    if [[ "$ENVIRONMENT" == "production" ]]; then
        build_args+=("--no-cache")
    fi
    
    # Execute build
    if docker build "${build_args[@]}"; then
        log_info "Image build successful!"
        log_info "Image tags:"
        log_info "  $REGISTRY/$IMAGE_NAME:$VERSION"
        log_info "  $REGISTRY/$IMAGE_NAME:latest"
    else
        log_error "Image build failed!"
        exit 1
    fi
}

# Show image information
show_image_info() {
    log_info "Image information:"
    docker images "$REGISTRY/$IMAGE_NAME" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"
}

# Push image (optional)
push_image() {
    if [[ "$PUSH" == "true" ]]; then
        log_info "Pushing image to registry..."
        docker push "$REGISTRY/$IMAGE_NAME:$VERSION"
        docker push "$REGISTRY/$IMAGE_NAME:latest"
        log_info "Image push completed!"
    fi
}

# Cleanup old images (optional)
cleanup_old_images() {
    if [[ "$CLEANUP" == "true" ]]; then
        log_warn "Cleaning up dangling images..."
        docker image prune -f
        log_info "Cleanup completed!"
    fi
}

# Main function
main() {
    log_info "CiliKube Docker Build Script"
    log_info "============================"
    
    # Change to script directory
    cd "$(dirname "$0")"
    
    check_docker
    get_build_info
    build_image
    show_image_info
    push_image
    cleanup_old_images
    
    log_info "Build completed! 🎉"
    log_info ""
    log_info "Run container:"
    log_info "  docker run -d --name cilikube -p 8080:8080 $REGISTRY/$IMAGE_NAME:$VERSION"
    log_info ""
    log_info "Use Docker Compose:"
    log_info "  cd ../../ && docker-compose up -d"
}

# Help information
show_help() {
    echo "CiliKube Docker Build Script"
    echo ""
    echo "Usage:"
    echo "  $0 [version] [environment]"
    echo ""
    echo "Arguments:"
    echo "  version      Image version tag (default: dev)"
    echo "  environment  Build environment (default: production)"
    echo ""
    echo "Environment variables:"
    echo "  REGISTRY     Image registry address (default: cilliantech)"
    echo "  PUSH         Whether to push image (true/false)"
    echo "  CLEANUP      Whether to cleanup old images (true/false)"
    echo ""
    echo "Examples:"
    echo "  $0 v1.0.0 production"
    echo "  PUSH=true $0 v1.0.0"
    echo "  CLEANUP=true $0 latest development"
}

# Handle command line arguments
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    *)
        main "$@"
        ;;
esac