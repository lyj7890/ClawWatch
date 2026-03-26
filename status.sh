#!/bin/bash
# Claw Watch - 状态检查脚本

cd "$(dirname "$0")"

PID_FILE="clawwatch.pid"
LOG_FILE="clawwatch.log"

echo "📊 Claw Watch Status"
echo "===================="
echo ""

# 检查 PID 文件
if [ ! -f "$PID_FILE" ]; then
    echo "Status: ❌ Not running (no PID file)"
    exit 1
fi

PID=$(cat "$PID_FILE")

# 检查进程
if ps -p "$PID" > /dev/null 2>&1; then
    echo "Status: ✅ Running"
    echo "PID: $PID"
    
    # 获取端口
    PORT=$(lsof -nP -iTCP -sTCP:LISTEN | grep "$PID" | awk '{print $9}' | cut -d: -f2 || echo "3939")
    echo "Port: $PORT"
    echo "URL: http://localhost:$PORT"
    echo ""
    
    # 检查服务是否响应
    if curl -s http://localhost:${PORT}/api/agents > /dev/null 2>&1; then
        echo "Health: ✅ Responding"
    else
        echo "Health: ⚠️  Not responding"
    fi
    
    echo ""
    echo "Recent logs (last 10 lines):"
    echo "----------------------------"
    tail -10 "$LOG_FILE" 2>/dev/null || echo "(no logs)"
else
    echo "Status: ❌ Not running (process $PID not found)"
    echo ""
    echo "Removing stale PID file..."
    rm -f "$PID_FILE"
    exit 1
fi

