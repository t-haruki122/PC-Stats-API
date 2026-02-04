# Worker Monitoring Agent

軽量なシステムモニタリングエージェント。CPU、RAM、GPUのメトリクスを収集し、HTTP APIとWeb UIで可視化します。

## 特徴

- **マルチプラットフォーム**: Windows / Linux 対応
- **GPU対応**: NVIDIA / AMD GPU をサポート
- **リアルタイム可視化**: モダンなWeb UIでメトリクスを表示
- **軽量設計**: 単一バイナリ、DB不要
- **サービス化**: systemd / Windows Service として常駐可能

## クイックスタート

### 依存関係のインストール

```bash
go mod download
```

### ビルド

```bash
# Windows
go build -o agent.exe ./cmd/agent

# Linux
go build -o agent ./cmd/agent
```

### 実行

```bash
# Windows
.\agent.exe

# Linux
./agent
```

ブラウザで `http://localhost:8080/ui/` にアクセスしてWeb UIを開きます。

## 設定

環境変数で設定をカスタマイズできます：

| 環境変数 | デフォルト | 説明 |
|---------|-----------|------|
| `PORT` | 8080 | HTTPサーバーのポート |
| `INTERVAL_SEC` | 30 | メトリクス収集間隔（秒） |
| `HISTORY_SIZE` | 720 | 履歴保持数（720 = 6時間@30秒） |

### 例

```bash
# Windows
$env:PORT="9090"
$env:INTERVAL_SEC="10"
.\agent.exe

# Linux
PORT=9090 INTERVAL_SEC=10 ./agent
```

## API エンドポイント

### `GET /health`

ヘルスチェック

```json
{
  "status": "ok"
}
```

### `GET /metrics/latest`

最新のメトリクス

```json
{
  "timestamp": "2026-02-05T07:48:11.299145233+09:00",
  "cpu": {
    "model": "AMD Ryzen 5 3500 6-Core Processor",
    "cores": 6,
    "threads": 6,
    "usage": 0.112225006961855,
    "load_avg": [0.53, 0.78, 0.85],
    "frequency_mhz": 3600
  },
  "ram": {
    "total_mb": 15909,
    "used_mb": 8982,
    "free_mb": 6555,
    "usage": 0.564581093453417
  },
  "gpu": {
    "vendor": "amd",
    "model": "AMD Radeon RX 6600 XT",
    "util": 0,
    "temperature_c": 36,
    "vram_total_mb": 8176,
    "vram_used_mb": 87
  }
}
```

### `GET /metrics/history?seconds=300`

履歴データ（デフォルト: 直近5分）

```json
{
  "interval_sec": 5,
  "samples": [...]
}
```

## サービスとして実行

### Linux (systemd)

```bash
# ビルド
go build -o agent ./cmd/agent

# バイナリを配置
sudo mkdir -p /opt/worker-agent
sudo cp agent /opt/worker-agent/
sudo cp -r web /opt/worker-agent/

# サービスファイルをコピー
sudo cp scripts/agent.service /etc/systemd/system/

# サービスを有効化・起動
sudo systemctl daemon-reload
sudo systemctl enable agent.service
sudo systemctl start agent.service

# ステータス確認
sudo systemctl status agent.service

# ログ確認
journalctl -u agent.service -f
```

### Windows (Windows Service)

```powershell
# ビルド
go build -o agent.exe ./cmd/agent

# 管理者権限でPowerShellを開く
# サービスをインストール
.\scripts\install-service.ps1 install

# サービスを起動
.\scripts\install-service.ps1 start

# サービスのステータス確認
Get-Service WorkerMonitorAgent

# アンインストール
.\scripts\install-service.ps1 uninstall
```

## GPU サポート

### NVIDIA

`nvidia-smi` コマンドが必要です。NVIDIA ドライバーに含まれています。

### AMD

- **Linux**: `rocm-smi` が必要（ROCm インストール時に含まれる）
- **Windows**: WMI経由で基本情報のみ取得（使用率は制限あり）
