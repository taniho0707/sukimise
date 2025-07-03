#!/bin/bash

echo "フロントエンドテストスクリプト"
echo "==============================="

cd frontend

echo "1. パッケージの確認..."
if [ -f "package.json" ]; then
    echo "✓ package.json が存在します"
else
    echo "✗ package.json が見つかりません"
    exit 1
fi

echo "2. 依存関係のインストール..."
npm install

echo "3. TypeScript設定の確認..."
if [ -f "tsconfig.json" ]; then
    echo "✓ tsconfig.json が存在します"
else
    echo "✗ tsconfig.json が見つかりません"
    exit 1
fi

echo "4. Vite設定の確認..."
if [ -f "vite.config.ts" ]; then
    echo "✓ vite.config.ts が存在します"
else
    echo "✗ vite.config.ts が見つかりません"
    exit 1
fi

echo "5. ビルドテスト..."
npm run build

echo "6. 開発サーバーを起動します..."
echo "ブラウザで http://localhost:3000 にアクセスしてください"
echo "停止するには Ctrl+C を押してください"
npm run dev