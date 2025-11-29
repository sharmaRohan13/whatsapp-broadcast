#!/bin/bash

# WhatsApp Broadcast - Build Script
# Builds executables for macOS and Windows and moves them to bin folder

echo "🔨 Building WhatsApp Broadcast executables..."
echo ""

# Create bin folder if it doesn't exist
mkdir -p bin

# Build for macOS
echo "📦 Building for macOS..."
go build -o bin/whatsapp-mac main.go
if [ $? -eq 0 ]; then
    echo "✅ macOS build complete: bin/whatsapp-mac"
else
    echo "❌ macOS build failed"
    exit 1
fi

echo ""

# Build for Windows
echo "📦 Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o bin/whatsapp-win.exe main.go
if [ $? -eq 0 ]; then
    echo "✅ Windows build complete: bin/whatsapp-win.exe"
else
    echo "❌ Windows build failed"
    exit 1
fi

echo ""
echo "🎉 All builds completed successfully!"
echo ""
echo "Files created in bin/ folder:"
ls -lh bin/whatsapp-mac bin/whatsapp-win.exe 2>/dev/null | awk '{print "  - " $9 " (" $5 ")"}'
