#!/bin/bash
# ClawWatch - 构建前端并部署

set -e

echo "🔨 Building frontend..."
cd frontend
npm run build
cd ..

echo "✅ Frontend built successfully!"
echo ""
echo "🔄 Restarting ClawWatch..."
./restart-daemon.sh

echo ""
echo "✅ Deployment complete!"
echo "   Access: http://localhost:3939"
