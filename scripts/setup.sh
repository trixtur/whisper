#!/bin/bash

set -e

echo "Speech-to-Clipboard Setup Script"
echo "================================="
echo ""

# Check for Go installation
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

echo "✓ Go is installed: $(go version)"

# Check for PortAudio
echo ""
echo "Checking for PortAudio..."
if command -v pkg-config &> /dev/null && pkg-config --exists portaudio-2.0; then
    echo "✓ PortAudio is installed"
else
    echo "⚠ PortAudio is not installed or not found via pkg-config"
    echo ""
    echo "Please install PortAudio:"
    echo "  macOS:   brew install portaudio"
    echo "  Ubuntu:  sudo apt-get install portaudio19-dev"
    echo "  Fedora:  sudo dnf install portaudio-devel"
    echo ""
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Install Go dependencies
echo ""
echo "Installing Go dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies installed"

# Check for OpenAI API key
echo ""
if [ -z "$OPENAI_API_KEY" ]; then
    echo "⚠ OPENAI_API_KEY environment variable is not set"
    echo ""
    echo "To use this application, you need an OpenAI API key."
    echo "Get one at: https://platform.openai.com/api-keys"
    echo ""
    echo "Then set it with:"
    echo "  export OPENAI_API_KEY='your-api-key-here'"
    echo ""
    echo "You can also copy .env.example to .env and add your key there."
else
    echo "✓ OPENAI_API_KEY is set"
fi

# Run tests
echo ""
echo "Running tests..."
if go test ./pkg/clipboard ./pkg/stt ./internal/config > /dev/null 2>&1; then
    echo "✓ Business logic tests passed"
else
    echo "⚠ Some tests failed (this is expected if PortAudio isn't installed)"
fi

# Build application
echo ""
echo "Building application..."
if go build -o speech-to-clipboard ./cmd/speech-to-clipboard; then
    echo "✓ Application built successfully"
    echo ""
    echo "Setup complete! Run the application with:"
    echo "  ./speech-to-clipboard"
else
    echo "✗ Build failed"
    exit 1
fi
