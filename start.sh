#!/bin/bash
# Claw Watch - 快速启动脚本

set -e

cd "$(dirname "$0")"

echo "🚀 Starting Claw Watch..."
echo ""

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker not found. Please install Docker first."
    exit 1
fi

# 检查 docker-compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose not found. Please install docker-compose first."
    exit 1
fi

# 启动服务
echo "📦 Starting containers..."
docker-compose up -d

echo ""
echo "✅ Claw Watch is running!"
echo ""
echo "🌐 Open in browser: http://localhost:3939"
echo ""
echo "Useful commands:"
echo "  docker-compose logs -f          # View logs"
echo "  docker-compose down             # Stop service"
echo "  docker-compose restart          # Restart service"
echo ""
