# Claw Watch - OpenClaw 实时会话监控器
FROM hub.intra.mlamp.cn/public/node:20-alpine

# 设置工作目录
WORKDIR /app

# 复制 package.json 并安装依赖
COPY package.json .
RUN npm install --production

# 复制应用文件
COPY web-viewer-server.js .
COPY web-viewer.html .

# 暴露端口
EXPOSE 3939

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:3939/api/agents', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"

# 启动服务
CMD ["node", "web-viewer-server.js"]
