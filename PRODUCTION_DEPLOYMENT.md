# 🚀 Sukimise 本番環境デプロイガイド

## 📋 前提条件

- Docker および Docker Compose がインストール済み
- `.env` ファイルが設定済み
- Discord Bot 用のトークンとGoogle Maps API キーが取得済み

## 🎯 **ワンコマンドデプロイ**

本番環境は以下のコマンド一つで起動できます：

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## 🔧 事前準備

### 1. 環境設定ファイルの作成

```bash
# .env ファイルを作成（例：.env.production.example をコピー）
cp .env.production.example .env

# 設定を編集
nano .env
```

### 2. 必須設定項目

```bash
# ポート設定
PORT=80                    # http://HOST_DOMAIN_NAME/ でアクセス
BACKEND_PORT=8080         # 内部のみ（/api/ プロキシ経由）
BOT_PORT=8082             # Discord Bot用

# セキュリティ設定（必須変更）
JWT_SECRET=CHANGE_THIS_TO_A_STRONG_RANDOM_JWT_SECRET_AT_LEAST_32_CHARS
POSTGRES_PASSWORD=CHANGE_THIS_STRONG_DATABASE_PASSWORD

# 外部サービス（必須設定）
DISCORD_TOKEN=YOUR_PRODUCTION_DISCORD_BOT_TOKEN
GOOGLE_MAPS_API_KEY=YOUR_PRODUCTION_GOOGLE_MAPS_API_KEY

# CORS設定（必須変更）
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com

# ユーザー設定（必須変更）
ADMIN_USERS=admin:$2a$10$GENERATE_PROPER_BCRYPT_HASH_FOR_ADMIN
EDITOR_USERS=editor:$2a$10$GENERATE_PROPER_BCRYPT_HASH_FOR_EDITOR
```

## 🌐 アクセス方法

### フロントエンド
```
http://HOST_DOMAIN_NAME:PORT/
```
- 静的ファイル配信（React SPA）
- `/api/` への自動プロキシ

### バックエンド API
```
http://HOST_DOMAIN_NAME:PORT/api/
```
- nginx 経由でバックエンドにプロキシ
- 直接アクセス不可（セキュリティ）

### Discord Bot
```
http://HOST_DOMAIN_NAME:BOT_PORT/
```
- 独立したサービスとして動作

## 🏗️ アーキテクチャ

```
Internet
    ↓
PORT (nginx frontend)
    ├── / → 静的ファイル (React)
    ├── /api/ → プロキシ → BACKEND_PORT (Go API)
    └── /uploads/ → プロキシ → BACKEND_PORT (ファイル)

BOT_PORT (Discord Bot) ← 独立サービス
    ↓
BACKEND_PORT ← 内部API呼び出し

PostgreSQL ← 内部のみ
```

## 🚀 デプロイ手順

### 1. 基本デプロイ
```bash
# 設定確認
docker-compose -f docker-compose.prod.yml config

# サービス起動
docker-compose -f docker-compose.prod.yml up -d

# 起動確認
docker-compose -f docker-compose.prod.yml ps
```

### 2. カスタムポートでのデプロイ
```bash
# ポート8080でフロントエンドを公開
PORT=8080 docker-compose -f docker-compose.prod.yml up -d
```

### 3. 初回起動の確認
```bash
# ヘルスチェック
curl http://localhost:${PORT:-80}/health

# サービス状態確認
docker-compose -f docker-compose.prod.yml logs frontend
docker-compose -f docker-compose.prod.yml logs backend
docker-compose -f docker-compose.prod.yml logs discord-bot
```

## 🔍 トラブルシューティング

### 設定確認
```bash
# 設定値確認
docker-compose -f docker-compose.prod.yml config

# 環境変数確認
docker-compose -f docker-compose.prod.yml exec frontend env
```

### ログ確認
```bash
# 全サービスのログ
docker-compose -f docker-compose.prod.yml logs

# 特定サービスのログ
docker-compose -f docker-compose.prod.yml logs frontend
docker-compose -f docker-compose.prod.yml logs backend
```

### 再起動
```bash
# サービス停止
docker-compose -f docker-compose.prod.yml down

# 再ビルドして起動
docker-compose -f docker-compose.prod.yml up -d --build
```

### ポート競合の解決
```bash
# ポート使用状況確認
netstat -tulpn | grep :80

# カスタムポートで起動
PORT=8080 BOT_PORT=8083 docker-compose -f docker-compose.prod.yml up -d
```

## 🔒 セキュリティチェックリスト

### デプロイ前チェック
- [ ] JWT_SECRET を強力なランダム値に変更
- [ ] データベースパスワード変更
- [ ] Discord Token と Google Maps API Key 設定
- [ ] CORS_ALLOWED_ORIGINS に本番ドメイン設定
- [ ] 管理者・編集者パスワード設定（bcrypt）

### デプロイ後チェック
- [ ] ファイアウォール設定（PORT, BOT_PORTのみ公開）
- [ ] HTTPS証明書設定（本番環境）
- [ ] バックアップ設定
- [ ] モニタリング設定

## 📊 ヘルスチェック

### 自動ヘルスチェック
各サービスには自動ヘルスチェック機能が組み込まれています：

```bash
# Docker の健康状態確認
docker-compose -f docker-compose.prod.yml ps
```

### 手動ヘルスチェック
```bash
# フロントエンド
curl -f http://localhost:${PORT:-80}/health

# バックエンド（内部）
docker-compose -f docker-compose.prod.yml exec backend wget -q --spider http://localhost:8080/health

# Discord Bot
curl -f http://localhost:${BOT_PORT:-8082}/health
```

## 📈 本番運用

### スケーリング
```bash
# バックエンドサービスをスケール
docker-compose -f docker-compose.prod.yml up -d --scale backend=3
```

### バックアップ
```bash
# データベースバックアップ
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U sukimise_user sukimise > backup.sql

# アップロードファイルバックアップ
docker-compose -f docker-compose.prod.yml exec backend tar -czf /tmp/uploads.tar.gz /app/uploads
```

### アップデート
```bash
# 新しいイメージで再デプロイ
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d --build
```

## 🎉 デプロイ完了

デプロイが成功すると：

1. **フロントエンド**: `http://HOST_DOMAIN_NAME:PORT/` でアクセス可能
2. **API**: `http://HOST_DOMAIN_NAME:PORT/api/` で自動プロキシ
3. **Discord Bot**: `http://HOST_DOMAIN_NAME:BOT_PORT/` で独立動作

これで Sukimise の本番環境が完全に稼働します！