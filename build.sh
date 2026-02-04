#!/bin/bash
# Linux用ビルドスクリプト

echo "Building Worker Monitoring Agent for Linux..."

# ビルド
go build -o agent ./cmd/agent

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo ""
    echo "To run the agent:"
    echo "  ./agent"
    echo ""
    echo "Web UI will be available at:"
    echo "  http://localhost:8080/ui/"
else
    echo "✗ Build failed"
    exit 1
fi
