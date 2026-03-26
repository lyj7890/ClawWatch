#!/bin/bash
# Claw Watch - 重启脚本

set -e

cd "$(dirname "$0")"

echo "🔄 Restarting Claw Watch..."
echo ""

./stop-daemon.sh || true
sleep 1
./start-daemon.sh

