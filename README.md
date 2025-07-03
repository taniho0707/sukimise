# Sukimise - お気に入りの店記録サービス

複数の編集者が共同でお気に入りの店舗情報を記録・管理し、閲覧者に共有できるWebサービスです。

## 技術スタック

### バックエンド
- Go 1.21
- Gin (Web Framework)
- PostgreSQL
- JWT認証

### フロントエンド
- React 18
- TypeScript
- Vite
- React Router
- React Query
- Leaflet (地図表示)

### インフラ
- Docker & Docker Compose

## セットアップ

### 前提条件
- Docker
- Docker Compose

### 起動手順

1. リポジトリをクローン
```bash
git clone <repository-url>
cd sukimise
```

2. 環境変数ファイルを作成
```bash
cp .env.example .env
```

3. アプリケーションを起動

以下の3つの方法から選択できます：

### 方法1: Docker Compose（推奨）
```bash
# 全サービスをDockerで起動
docker-compose up -d

# ログを確認
docker-compose logs -f
```

### 方法2: ローカル実行（開発向け）
```bash
# PostgreSQLのみDockerで起動し、アプリはローカルで実行
./scripts/run-local.sh

# 別ターミナルでバックエンド起動
cd backend && go run cmd/server/main.go

# 別ターミナルでフロントエンド起動
cd frontend && npm install && npm run dev
```

### 方法3: 開発用Docker Compose
```bash
# PostgreSQLポートも公開する開発用設定
docker-compose -f docker-compose.dev.yml up -d
```

4. アプリケーションにアクセス
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080

### デフォルトユーザー
- **管理者**: username: `admin`, password: `admin123`
- **編集者**: username: `editor`, password: `editor123`

### トラブルシューティング

#### Dockerビルドエラーが発生する場合
```bash
# 既存のコンテナとイメージを削除
docker-compose down --volumes --remove-orphans
docker system prune -f

# 再起動
docker-compose up -d
```

#### ローカル実行でPostgreSQLに接続できない場合
```bash
# PostgreSQLコンテナの状態確認
docker ps | grep postgres

# PostgreSQLログ確認
docker logs sukimise-postgres
```

#### フロントエンド
```bash
cd frontend
npm install
npm run dev
```

## データベース

PostgreSQLを使用しています。マイグレーションファイルは `backend/migrations/` に配置されています。

### マイグレーション実行
```bash
# マイグレーション実行用のツールをインストール
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# マイグレーション実行
migrate -path backend/migrations -database "postgres://sukimise_user:sukimise_password@localhost:5432/sukimise?sslmode=disable" up
```

## 主要機能

1. **店舗情報管理**
   - 店舗の登録・編集・削除
   - 店舗情報の詳細表示

2. **レビュー機能**
   - 個人別の評価・コメント
   - 訪問記録管理

3. **地図表示**
   - OpenStreetMapを使用
   - 店舗位置の表示

4. **認証・認可**
   - JWT認証
   - ロールベースアクセス制御

## API仕様

APIエンドポイントは `/api/v1` 配下に配置されています。

### 認証
- `POST /api/v1/auth/login` - ログイン
- `POST /api/v1/auth/refresh` - トークンリフレッシュ

### 店舗
- `GET /api/v1/stores` - 店舗一覧取得
- `GET /api/v1/stores/:id` - 店舗詳細取得
- `POST /api/v1/stores` - 店舗作成（要認証）
- `PUT /api/v1/stores/:id` - 店舗更新（要認証）
- `DELETE /api/v1/stores/:id` - 店舗削除（要認証）

### レビュー
- `POST /api/v1/reviews` - レビュー作成（要認証）
- `PUT /api/v1/reviews/:id` - レビュー更新（要認証）
- `DELETE /api/v1/reviews/:id` - レビュー削除（要認証）

## ライセンス

MIT License