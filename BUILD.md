# Sukimise ビルド・デプロイガイド

## 📋 概要

Sukimiseの本番環境デプロイに関する包括的なガイドです。

## 🏗️ ビルド済みスクリプト

### 1. フロントエンドビルド
```bash
# 個別ビルド
./scripts/build-frontend.sh

# 手動ビルド
cd frontend
npm ci
npm run build
```

### 2. 本番環境デプロイ
```bash
# 自動デプロイ（推奨）
./scripts/deploy-production.sh

# 手動デプロイ
cp .env.production.example .env.production
# .env.production を編集
nano .env.production
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d --build
```

## 🔧 本番環境設定

### 1. 環境設定ファイル
```bash
# 本番環境設定をコピー
cp .env.production.example .env.production

# 必須設定項目を編集
nano .env.production
```

### 2. 必須設定項目
```bash
# セキュリティ（必須変更）
JWT_SECRET=CHANGE_THIS_TO_A_STRONG_RANDOM_JWT_SECRET_AT_LEAST_32_CHARS
POSTGRES_PASSWORD=CHANGE_THIS_STRONG_DATABASE_PASSWORD

# 外部サービス（必須設定）
DISCORD_TOKEN=YOUR_PRODUCTION_DISCORD_BOT_TOKEN
GOOGLE_MAPS_API_KEY=YOUR_PRODUCTION_GOOGLE_MAPS_API_KEY

# ドメイン設定（必須変更）
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# ユーザー設定（必須変更）
ADMIN_USERS=admin:$2a$10$GENERATE_PROPER_BCRYPT_HASH_FOR_ADMIN
EDITOR_USERS=editor:$2a$10$GENERATE_PROPER_BCRYPT_HASH_FOR_EDITOR
```

### 3. ポート設定
```bash
# 本番環境のポート設定
FRONTEND_PORT=80            # nginx（メインアクセスポイント）
BOT_PORT=8082              # Discord Bot
PORT=8080                  # Backend（内部のみ）
```

## 🐳 Docker構成

### 開発環境
```bash
# 開発環境起動
docker-compose up -d

# アクセス
Frontend: http://localhost:3000 (Vite + Proxy)
Discord Bot: http://localhost:8082
```

### 本番環境
```bash
# 本番環境起動
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d

# アクセス
Frontend: http://localhost:80 (nginx + Static)
Discord Bot: http://localhost:8082
```

## 📁 ビルド出力

### フロントエンド
```
frontend/dist/
├── index.html                    # SPAエントリーポイント
├── assets/
│   ├── index-[hash].js          # バンドルされたJavaScript
│   ├── index-[hash].css         # バンドルされたCSS
│   └── [other-assets]           # 画像・フォントなど
└── favicon.ico
```

### バックエンド
```bash
# Goバイナリ（本番環境ではDockerビルド）
go build -o server cmd/server/main.go
```

## 🌐 アーキテクチャ

### 開発環境
```
Browser → localhost:3000 (Vite)
                ↓ /api proxy
            localhost:8080 (Go Backend) ← localhost:8082 (Discord Bot)
                ↓
            PostgreSQL (internal)
```

### 本番環境
```
Internet → Port 80 (nginx)
              ├── Static Files (React build)
              └── /api proxy → Port 8080 (Go Backend)
                                    ↑
Internet → Port 8082 (Discord Bot) ┘
                ↓
            PostgreSQL (internal)
```

## ✅ ヘルスチェック

### 自動ヘルスチェック
```bash
# デプロイスクリプトが自動実行
curl -f http://localhost/health          # Frontend
curl -f http://localhost:8080/health     # Backend（内部）
```

### 手動ヘルスチェック
```bash
# サービス状態確認
docker-compose -f docker-compose.prod.yml ps

# ログ確認
docker-compose -f docker-compose.prod.yml logs nginx
docker-compose -f docker-compose.prod.yml logs backend
docker-compose -f docker-compose.prod.yml logs discord-bot
```

## 🛡️ セキュリティチェックリスト

### デプロイ前
- [ ] JWT_SECRET を強力なランダム値に変更
- [ ] データベースパスワード変更
- [ ] Discord Token と Google Maps API Key 設定
- [ ] CORS_ALLOWED_ORIGINS に本番ドメイン設定
- [ ] 管理者・編集者パスワード設定（bcrypt）

### デプロイ後
- [ ] ファイアウォール設定（80, 8082のみ公開）
- [ ] HTTPS証明書設定
- [ ] バックアップ設定
- [ ] モニタリング設定

## 🚨 トラブルシューティング

### ビルドエラー
```bash
# 依存関係の問題
cd frontend
rm -rf node_modules package-lock.json
npm install

# TypeScriptエラー
npm run lint:fix
```

### デプロイエラー
```bash
# 設定確認
docker-compose -f docker-compose.prod.yml config

# ログ確認
docker-compose -f docker-compose.prod.yml logs [service-name]

# コンテナ再構築
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --build
```

### ポート競合
```bash
# ポート使用状況確認
netstat -tulpn | grep :80
lsof -i :80

# カスタムポート使用
FRONTEND_PORT=8080 ./scripts/deploy-production.sh
```

## 📚 関連ドキュメント

- [DEPLOYMENT.md](docs/DEPLOYMENT.md) - 詳細なデプロイメントガイド
- [SECURITY_AUDIT_REPORT_2025-07-13.md](SECURITY_AUDIT_REPORT_2025-07-13.md) - セキュリティ監査結果
- [CLAUDE.md](CLAUDE.md) - プロジェクト全体の仕様