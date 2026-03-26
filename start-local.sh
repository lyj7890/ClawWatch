#!/bin/bash
# Claw Watch - 本地启动脚本（非容器化）

set -e

cd "$(dirname "$0")"

echo "🚀 Starting Claw Watch (Local Mode)..."
echo ""

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js not found. Please install Node.js first."
    exit 1
fi

# 检查依赖
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# 设置环境变量
export PORT=${PORT:-3939}
export NODE_ENV=${NODE_ENV:-production}

echo "✅ Starting server on port $PORT..."
echo ""

# 启动服务
node web-viewer-server.js

