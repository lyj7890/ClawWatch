#!/bin/bash
# Claw Watch - 停止脚本

set -e

cd "$(dirname "$0")"

PID_FILE="clawwatch.pid"

if [ ! -f "$PID_FILE" ]; then
    echo "⚠️  Claw Watch is not running (no PID file found)"
    exit 1
fi

PID=$(cat "$PID_FILE")

if ps -p "$PID" > /dev/null 2>&1; then
    echo "🛑 Stopping Claw Watch (PID: $PID)..."
    kill "$PID"
    
    # 等待进程结束
    for i in {1..10}; do
        if ! ps -p "$PID" > /dev/null 2>&1; then
            break
        fi
        sleep 1
    done
    
    # 如果还在运行，强制杀死
    if ps -p "$PID" > /dev/null 2>&1; then
        echo "⚠️  Process still running, force killing..."
        kill -9 "$PID"
    fi
    
    rm -f "$PID_FILE"
    echo "✅ Claw Watch stopped"
else
    echo "⚠️  Process $PID not found, removing stale PID file"
    rm -f "$PID_FILE"
fi

