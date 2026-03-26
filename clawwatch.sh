#!/bin/bash
# ClawWatch - 统一管理脚本

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

PID_FILE="$SCRIPT_DIR/clawwatch.pid"
LOG_FILE="$SCRIPT_DIR/clawwatch.log"

# 帮助信息
show_help() {
    cat << EOF
🦞 ClawWatch - OpenClaw Session Monitor

Usage: ./clawwatch.sh [command]

Commands:
  start       构建前端并启动服务（默认）
  stop        停止服务
  restart     重启服务
  status      查看服务状态
  logs        查看日志
  help        显示此帮助信息

Examples:
  ./clawwatch.sh          # 启动服务
  ./clawwatch.sh start    # 同上
  ./clawwatch.sh stop     # 停止服务
  ./clawwatch.sh status   # 查看状态
  ./clawwatch.sh logs     # 查看日志

访问: http://localhost:3939
EOF
}

# 检查服务状态
check_status() {
    if [ -f "$PID_FILE" ]; then
        PID=$(cat "$PID_FILE")
        if ps -p "$PID" > /dev/null 2>&1; then
            echo "running:$PID"
            return 0
        fi
    fi
    echo "stopped"
    return 1
}

# 启动服务
start_service() {
    # 检查是否已运行
    if check_status | grep -q "running"; then
        PID=$(cat "$PID_FILE")
        echo "⚠️  ClawWatch is already running (PID: $PID)"
        echo "   Use './clawwatch.sh restart' to restart"
        return 1
    fi

    echo "🔨 Building frontend..."
    cd frontend
    npm run build
    cd ..

    echo "🚀 Starting ClawWatch..."
    nohup node web-viewer-server.js > "$LOG_FILE" 2>&1 &
    PID=$!
    echo $PID > "$PID_FILE"

    sleep 2

    if ps -p "$PID" > /dev/null 2>&1; then
        echo "✅ ClawWatch started successfully!"
        echo ""
        echo "   PID: $PID"
        echo "   Port: 3939"
        echo "   URL: http://localhost:3939"
        echo "   Logs: tail -f $LOG_FILE"
        echo ""
        echo "Commands:"
        echo "   ./clawwatch.sh stop     - Stop the service"
        echo "   ./clawwatch.sh restart  - Restart the service"
        echo "   ./clawwatch.sh status   - Check status"
        echo "   ./clawwatch.sh logs     - View logs"
    else
        echo "❌ Failed to start ClawWatch"
        rm -f "$PID_FILE"
        echo "   Check logs: cat $LOG_FILE"
        return 1
    fi
}

# 停止服务
stop_service() {
    if ! check_status | grep -q "running"; then
        echo "⚠️  ClawWatch is not running"
        rm -f "$PID_FILE"
        return 0
    fi

    PID=$(cat "$PID_FILE")
    echo "🛑 Stopping ClawWatch (PID: $PID)..."

    kill "$PID" 2>/dev/null || true

    # 等待进程结束
    for i in {1..10}; do
        if ! ps -p "$PID" > /dev/null 2>&1; then
            break
        fi
        sleep 0.5
    done

    # 强制杀死
    if ps -p "$PID" > /dev/null 2>&1; then
        kill -9 "$PID" 2>/dev/null || true
    fi

    rm -f "$PID_FILE"
    echo "✅ ClawWatch stopped"
}

# 重启服务
restart_service() {
    echo "🔄 Restarting ClawWatch..."
    echo ""
    stop_service
    echo ""
    start_service
}

# 查看状态
show_status() {
    STATUS=$(check_status)

    if echo "$STATUS" | grep -q "running"; then
        PID=$(echo "$STATUS" | cut -d: -f2)
        echo "✅ ClawWatch is running"
        echo ""
        echo "   PID: $PID"
        echo "   Port: 3939"
        echo "   URL: http://localhost:3939"
        echo "   Log: $LOG_FILE"
        echo ""
        echo "   Memory: $(ps -o rss= -p $PID | awk '{printf "%.1f MB", $1/1024}')"
        echo "   Uptime: $(ps -o etime= -p $PID | xargs)"
    else
        echo "⚠️  ClawWatch is not running"
        echo ""
        echo "   Start with: ./clawwatch.sh start"
    fi
}

# 查看日志
show_logs() {
    if [ ! -f "$LOG_FILE" ]; then
        echo "⚠️  Log file not found: $LOG_FILE"
        return 1
    fi

    echo "📋 ClawWatch Logs (last 50 lines):"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    tail -50 "$LOG_FILE"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Follow logs: tail -f $LOG_FILE"
}

# 主逻辑
COMMAND=${1:-start}

case "$COMMAND" in
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        restart_service
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo "❌ Unknown command: $COMMAND"
        echo ""
        show_help
        exit 1
        ;;
esac
