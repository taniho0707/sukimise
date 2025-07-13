# Sukimise Discord Bot

Sukimiseシステム用のDiscord botです。GoogleMapのURLから店舗情報を自動的に抽出し、Sukimiseデータベースに登録できます。

## 機能

- **アカウント連携**: DiscordアカウントとSukimiseアカウントの連携
- **店舗登録**: GoogleMap URLから店舗情報を自動抽出・登録
- **セキュア認証**: JWTトークンによる安全な認証システム
- **自動情報取得**: 店舗名、住所、座標、WebサイトURLの自動抽出

## スラッシュコマンド

### `/connect <username> <password>`
DiscordアカウントとSukimiseアカウントを連携します。

**例:**
```
/connect myusername mypassword
```

**注意事項:**
- 1つのSukimiseアカウントには1つのDiscordアカウントのみ連携可能
- 1つのDiscordアカウントには1つのSukimiseアカウントのみ連携可能

### `/add <google_maps_url>`
GoogleMap URLから店舗情報を取得してSukimiseに登録します。

**例:**
```
/add https://www.google.com/maps/place/Restaurant+Name/@35.6812,139.7671,17z/data=!3m1!4b1!4m5!3m4!1s0x...
```

**対応URL形式:**
- `https://www.google.com/maps/place/...`
- `https://maps.google.com/maps/place/...`
- `https://goo.gl/maps/...`

### `/disconnect`
DiscordアカウントとSukimiseアカウントの連携を切断します。

### `/help`
使用方法とコマンド一覧を表示します。

## セットアップ

### 1. Discord Developer Portal でボットを作成

1. [Discord Developer Portal](https://discord.com/developers/applications) にアクセス
2. "New Application" をクリック
3. アプリケーション名を入力（例: "Sukimise Bot"）
4. 左メニューから "Bot" を選択
5. "Add Bot" をクリック
6. TOKEN をコピーして保存（後で.envファイルで使用）

### 2. ボット権限の設定

Bot設定ページで以下の権限を有効化:
- `applications.commands` (スラッシュコマンド使用)
- `bot` (基本的なボット機能)

Botの権限として以下を設定:
- `Send Messages`
- `Use Slash Commands`

### 3. サーバーへの招待

1. 左メニューから "OAuth2" → "URL Generator" を選択
2. "Scopes" で `bot` と `applications.commands` を選択
3. "Bot Permissions" で必要な権限を選択
4. 生成されたURLでサーバーに招待

### 4. 環境変数の設定

`.env` ファイルに以下を追加:

```env
# Discord Bot Configuration
DISCORD_TOKEN=your_discord_bot_token_here
SUKIMISE_API_URL=http://backend:8080
BOT_PORT=8081
```

### 5. データベースマイグレーション

Discord-Sukimise連携用のテーブルを作成:

```sql
-- 006_create_discord_links_table.up.sql は自動的に実行されます
```

### 6. Docker Composeで起動

```bash
# Sukimiseシステム全体を起動
docker-compose up -d

# ログの確認
docker-compose logs discord-bot
```

## 使用フロー

1. **ボットをサーバーに招待**
2. **アカウント連携**: `/connect username password`
3. **店舗登録**: `/add <google_maps_url>`
4. **結果確認**: チャンネルに登録結果が表示されます

## 技術仕様

### アーキテクチャ
- **言語**: Go 1.21
- **フレームワーク**: discordgo
- **データベース**: PostgreSQL（Sukimiseシステムと共有）
- **認証**: JWT（Sukimise API連携）

### ディレクトリ構造
```
discord-bot/
├── cmd/
│   └── main.go              # エントリーポイント
├── internal/
│   ├── config/              # 設定管理
│   ├── handlers/            # Discord コマンドハンドラー
│   ├── models/              # データモデル
│   ├── services/            # ビジネスロジック
│   └── utils/               # ユーティリティ（GoogleMap解析等）
├── go.mod
└── README.md
```

### 自動抽出される情報
- **店舗名**: URLパスから抽出
- **座標**: URL内の緯度・経度情報
- **住所**: HTMLページから抽出（可能な場合）
- **WebサイトURL**: HTMLページから抽出（可能な場合）
- **GoogleMap URL**: 入力されたURL

### セキュリティ機能
- パスワードは認証にのみ使用（保存されません）
- JWTアクセストークンとリフレッシュトークンの安全な管理
- トークン自動リフレッシュ機能
- ロールベースアクセス制御

## トラブルシューティング

### ボットが応答しない
1. Discord Tokenが正しく設定されているか確認
2. ボットがサーバーに正しく招待されているか確認
3. 必要な権限が設定されているか確認

### 認証エラー
1. Sukimiseアカウントの認証情報が正しいか確認
2. SukimiseAPIサーバーが起動しているか確認
3. ネットワーク接続を確認

### GoogleMap URL解析エラー
1. 対応しているURL形式か確認
2. URLに座標情報が含まれているか確認
3. ページが公開されているか確認

## ログ例

### 成功時のログ
```
✅ Store Successfully Registered!

Store Details:
• Name: Restaurant Name
• Address: Tokyo, Japan
• Store ID: 123e4567-e89b-12d3-a456-426614174000

Registration Info:
• Registered by: @username
• Google Maps URL: https://www.google.com/maps/place/...
```

### エラー時のログ
```
❌ Store Registration Failed
Invalid Google Maps URL. Must start with https://www.google.com/maps/place/
```

## 開発者向け情報

### 新機能の追加
1. `internal/handlers/command_handler.go` で新しいコマンドを定義
2. `internal/services/discord_service.go` でビジネスロジックを実装
3. `commands` スライスに新しいコマンドを追加

### Google Map解析の改善
`internal/utils/google_maps.go` で情報抽出ロジックを改善できます。

### API連携の拡張
`internal/services/discord_service.go` でSukimise API呼び出しを追加できます。

## サポート

質問や問題報告は、Sukimiseプロジェクトのメインリポジトリまでお願いします。