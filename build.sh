#!/bin/bash
# Linux用ビルドスクリプト

echo "Building Worker Monitoring Agent for Linux..."

# webディレクトリをcmd/agent/にコピー（embed用）
echo "Copying web files for embedding..."
cp -r web cmd/agent/web

# CGO無効化、サイズ削減
export CGO_ENABLED=0

# ビルド（-ldflags で デバッグ情報とシンボルテーブルを削除してサイズ削減）
go build -ldflags="-s -w" -o agent ./cmd/agent

# webディレクトリを削除（クリーンアップ）
rm -rf cmd/agent/web

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo ""
    echo "Binary size:"
    ls -lh agent | awk '{print "  " $5 " - agent"}'
    echo ""
    echo "To run the agent:"
    echo "  ./agent"
    echo ""
    echo "Web UI will be available at:"
    echo "  http://localhost:8080/ui/"
else
    echo "✗ Build failed"
    # クリーンアップ
    rm -rf cmd/agent/web
    exit 1
fi
