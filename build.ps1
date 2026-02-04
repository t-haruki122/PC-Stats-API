# Windows用ビルドスクリプト

Write-Host "Building Worker Monitoring Agent for Windows..." -ForegroundColor Cyan

# webディレクトリをcmd/agent/にコピー（embed用）
Write-Host "Copying web files for embedding..." -ForegroundColor Yellow
Copy-Item -Path "web" -Destination "cmd\agent\web" -Recurse -Force

# CGO無効化、サイズ削減
$env:CGO_ENABLED = "0"

# ビルド（-ldflags で デバッグ情報とシンボルテーブルを削除してサイズ削減）
go build -ldflags="-s -w" -o agent.exe ./cmd/agent

# webディレクトリを削除（クリーンアップ）
Remove-Item -Path "cmd\agent\web" -Recurse -Force -ErrorAction SilentlyContinue

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Build successful!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Binary size:"
    $size = (Get-Item agent.exe).Length
    $sizeMB = [math]::Round($size / 1MB, 2)
    Write-Host "  $sizeMB MB - agent.exe"
    Write-Host ""
    Write-Host "To run the agent:"
    Write-Host "  .\agent.exe"
    Write-Host ""
    Write-Host "Web UI will be available at:"
    Write-Host "  http://localhost:8080/ui/"
} else {
    Write-Host "✗ Build failed" -ForegroundColor Red
    # クリーンアップ
    Remove-Item -Path "cmd\agent\web" -Recurse -Force -ErrorAction SilentlyContinue
    exit 1
}
