#!/bin/bash
# Claw Watch - 后台启动脚本

set -e

cd "$(dirname "$0")"

PID_FILE="clawwatch.pid"
LOG_FILE="clawwatch.log"

# 检查是否已经在运行
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        echo "⚠️  Claw Watch is already running (PID: $PID)"
        echo "   Use ./stop-daemon.sh to stop it first"
        exit 1
    else
        echo "🧹 Removing stale PID file..."
        rm -f "$PID_FILE"
    fi
fi

# 检查依赖
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# 启动服务
export PORT=${PORT:-3939}
export NODE_ENV=${NODE_ENV:-production}

echo "🚀 Starting Claw Watch in background..."
nohup node web-viewer-server.js > "$LOG_FILE" 2>&1 &

# 保存 PID
echo $! > "$PID_FILE"

sleep 2

# 检查是否成功启动
if ps -p $(cat "$PID_FILE") > /dev/null 2>&1; then
    echo "✅ Claw Watch started successfully!"
    echo ""
    echo "   PID: $(cat $PID_FILE)"
    echo "   Port: $PORT"
    echo "   URL: http://localhost:$PORT"
    echo "   Logs: tail -f $LOG_FILE"
    echo ""
    echo "Commands:"
    echo "   ./stop-daemon.sh    - Stop the service"
    echo "   ./restart-daemon.sh - Restart the service"
else
    echo "❌ Failed to start. Check $LOG_FILE for errors."
    rm -f "$PID_FILE"
    exit 1
fi

