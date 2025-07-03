#!/bin/bash

echo "Sukimise ローカル実行スクリプト"
echo "================================"

# PostgreSQLのみをDockerで起動
echo "PostgreSQLコンテナを起動しています..."
docker run -d \
  --name sukimise-postgres \
  -e POSTGRES_DB=sukimise \
  -e POSTGRES_USER=sukimise_user \
  -e POSTGRES_PASSWORD=sukimise_password \
  -p 5432:5432 \
  -v sukimise_postgres_data:/var/lib/postgresql/data \
  -v "$(pwd)/backend/init-db.sql:/docker-entrypoint-initdb.d/init-db.sql" \
  postgres:15

echo "PostgreSQLの起動を待機しています..."
sleep 10

# 環境変数の設定
export DATABASE_URL="postgres://sukimise_user:sukimise_password@localhost:5432/sukimise?sslmode=disable"
export JWT_SECRET="your-jwt-secret-key"
export PORT="8080"

echo ""
echo "PostgreSQLが起動しました。"
echo ""
echo "バックエンドを起動するには："
echo "  cd backend && go run cmd/server/main.go"
echo ""
echo "フロントエンドを起動するには："
echo "  cd frontend && npm install && npm run dev"
echo ""
echo "PostgreSQLを停止するには："
echo "  docker stop sukimise-postgres && docker rm sukimise-postgres"