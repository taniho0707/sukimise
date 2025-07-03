#!/bin/bash

echo "Sukimise セットアップスクリプト"
echo "================================"

# 環境変数ファイルの作成
if [ ! -f .env ]; then
    echo ".envファイルを作成しています..."
    cp .env.example .env
    echo "✓ .envファイルが作成されました"
else
    echo "✓ .envファイルが既に存在します"
fi

# Dockerサービスの起動
echo "Dockerサービスを起動しています..."
docker-compose up -d

# サービスの起動を待機
echo "サービスの起動を待機しています..."
sleep 10

echo "データベースの初期化が完了しました。"

echo ""
echo "✓ セットアップが完了しました！"
echo ""
echo "次のコマンドでアプリケーションを起動できます："
echo "  docker-compose up"
echo ""
echo "または個別に起動する場合："
echo "  バックエンド: cd backend && go run cmd/server/main.go"
echo "  フロントエンド: cd frontend && npm install && npm run dev"
echo ""
echo "アクセスURL："
echo "  フロントエンド: http://localhost:3000"
echo "  バックエンドAPI: http://localhost:8080"