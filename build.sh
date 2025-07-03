#!/bin/bash

# Build script for Janus CLI
# Builds cross-platform binaries

echo "Building Janus CLI..."

# Create build directory
mkdir -p build

# Build for different platforms
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o build/janus-linux-amd64 main.go

echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o build/janus-macos-amd64 main.go

echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o build/janus-macos-arm64 main.go

echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o build/janus-windows-amd64.exe main.go

echo "Building for current platform..."
go build -o build/janus main.go

echo "✅ Build complete! Binaries available in ./build/"
ls -la build/

echo ""
echo "To install locally:"
echo "  chmod +x build/janus"
echo "  sudo mv build/janus /usr/local/bin/"
echo ""
echo "To create example templates:"
echo "  mkdir -p ~/.janus/templates/nextjs"
echo "  # Add your template files to ~/.janus/templates/nextjs/"
