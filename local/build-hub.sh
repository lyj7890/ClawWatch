#!/usr/bin/env bash
# build-hub.sh - 构建 clawwatch-hub amd64 镜像并推送到 Harbor
# 用法: ./build-hub.sh [harbor地址] [harbor项目] [版本tag]
# 示例: ./build-hub.sh hub.intra.mlamp.cn public v1.0.0

set -e

HARBOR_HOST="${1:-hub.intra.mlamp.cn}"
HARBOR_PROJECT="${2:-public}"
VERSION="${3:-latest}"
IMAGE_NAME="clawwatch-hub"
FULL_TAG="${HARBOR_HOST}/${HARBOR_PROJECT}/${IMAGE_NAME}:${VERSION}"

echo "🏗️  构建 amd64 镜像: ${FULL_TAG}"
echo ""

# 检查 buildx
if ! docker buildx version &>/dev/null; then
  echo "❌ 需要 docker buildx，请升级 Docker Desktop"
  exit 1
fi

# 确保有 amd64 builder
docker buildx inspect multiarch-builder &>/dev/null || \
  docker buildx create --name multiarch-builder --driver docker-container --use

docker buildx use multiarch-builder

# 构建并推送
echo "📦 开始构建（跨平台 linux/amd64）..."
docker buildx build \
  --platform linux/amd64 \
  --file hub/Dockerfile \
  --tag "${FULL_TAG}" \
  --push \
  hub/

echo ""
echo "✅ 构建完成！"
echo "   镜像: ${FULL_TAG}"
echo ""
echo "💡 拉取验证:"
echo "   docker pull ${FULL_TAG}"
