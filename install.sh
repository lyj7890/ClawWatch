#!/bin/sh
# ClawWatch Agent 一键安装脚本
# 用法: curl -fsSL https://raw.githubusercontent.com/lyj7890/ClawWatch/main/install.sh | sh

set -e

REPO="lyj7890/ClawWatch"
BINARY="clawwatch-agent"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# 检测系统和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "不支持的架构: $ARCH"; exit 1 ;;
esac

case "$OS" in
  linux|darwin) ;;
  *) echo "不支持的系统: $OS"; exit 1 ;;
esac

SUFFIX="${OS}-${ARCH}"
echo "检测到系统: ${OS}/${ARCH}"

# 获取最新版本
echo "正在获取最新版本..."
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "错误: 无法获取最新版本"
  exit 1
fi

echo "最新版本: ${LATEST}"

# 下载
URL="https://github.com/${REPO}/releases/download/${LATEST}/${BINARY}-${SUFFIX}"
echo "正在下载 ${URL} ..."
curl -fsSL "$URL" -o "/tmp/${BINARY}"
chmod +x "/tmp/${BINARY}"

# 安装
if [ -w "$INSTALL_DIR" ]; then
  mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
  echo "需要 sudo 权限安装到 ${INSTALL_DIR}"
  sudo mv "/tmp/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi

echo ""
echo "✅ ClawWatch Agent ${LATEST} 安装成功!"
echo ""
echo "启动命令:"
echo "  ${BINARY} --hub wss://clawwatch.intra.mlamp.cn"
echo ""
echo "查看帮助:"
echo "  ${BINARY} --help"
