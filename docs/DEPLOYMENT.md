# Sukimise デプロイメントガイド

## ポート構成

### 開発環境
```
PORT=8080          # Backend API (外部公開)
FRONTEND_PORT=3000 # Frontend Dev Server (外部公開)  
BOT_PORT=8082      # Discord Bot (外部公開)
```

### 本番環境
```
PORT=8080          # Backend API (リバースプロキシ経由)
BOT_PORT=8082      # Discord Bot (内部またはメトリクス用)
```

## ネットワーク設計

### 開発環境 (docker-compose.yml)
```
Host:3000 → Frontend Container:3000 (React Dev Server)
Host:8080 → Backend Container:8080 (Go API)
Host:8082 → Discord Bot Container:8082 (Go Bot)
```

### 本番環境 (docker-compose.prod.yml)
```
Internet:443 → Nginx:443 (HTTPS終端)
           → Backend:8080 (API)
           → Static Files (Frontend)

Internal:8082 → Discord Bot:8082 (ヘルスチェック)
```

## デプロイメント設定

### 1. 環境変数設定
```bash
# 本番環境用設定をコピー
cp .env.production.example .env.production

# 必須項目を編集
nano .env.production
```

### 2. 本番環境でのポート公開
```bash
# 本番環境起動
docker-compose -f docker-compose.prod.yml up -d

# 公開ポート確認
docker-compose -f docker-compose.prod.yml ps
```

### 3. リバースプロキシ設定 (推奨)
```nginx
# nginx.conf 例
server {
    listen 80;
    server_name yourdomain.com;
    
    # Frontend静的ファイル
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }
    
    # Backend API
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## セキュリティ考慮事項

### ファイアウォール設定
```bash
# 必要なポートのみ公開
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw deny 8080/tcp   # Backend直接アクセス拒否
ufw deny 8082/tcp   # Bot直接アクセス拒否
```

### 本番環境でのポート管理
- **外部公開**: 80, 443のみ (リバースプロキシ経由)
- **内部のみ**: 8080 (Backend), 8082 (Bot), 5432 (PostgreSQL)

## 設定変更手順

### ポート変更
```bash
# .envファイルでポート指定
echo "PORT=9000" >> .env
echo "BOT_PORT=9002" >> .env

# 再起動
docker-compose down
docker-compose up -d
```

### 本番環境への移行
```bash
# 開発環境停止
docker-compose down

# 本番環境起動
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d
```

## ヘルスチェック

### エンドポイント
- Backend: `http://localhost:8080/health`
- Discord Bot: `http://localhost:8082/health`

### 監視設定
```bash
# ヘルスチェック確認
curl http://localhost:8080/health
curl http://localhost:8082/health
```

## トラブルシューティング

### ポート競合の解決
```bash
# ポート使用状況確認
netstat -tulpn | grep :8080

# 使用中のプロセス確認
lsof -i :8080
```

### ログ確認
```bash
# サービス別ログ
docker-compose logs backend
docker-compose logs discord-bot
docker-compose logs frontend
```