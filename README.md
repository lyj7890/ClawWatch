# Claw Watch 🦞

OpenClaw 实时会话监控器 - 一个强大的 Web 界面，用于实时查看和监控 OpenClaw 会话日志，支持 WebSocket 实时推送、思考过程展示、过滤搜索等功能。

## ✨ 功能特性

- 🔴 **实时监控**：WebSocket 实时推送，无需刷新
- 💭 **思考过程**：显示 AI 的思考过程（可折叠）
- 🎯 **智能过滤**：按角色（用户/助手）过滤消息
- 🔍 **全文搜索**：快速搜索会话内容
- 🎨 **美观界面**：暗色主题，响应式设计
- 🔧 **工具调用**：高亮显示工具调用和结果
- 📊 **Token 统计**：显示每条消息的 Token 使用情况
- 🔄 **多 Agent**：支持查看不同 Agent 的会话

## 🚀 快速开始

### 使用 Docker Compose（推荐）

```bash
# 进入项目目录
cd ~/.openclaw/ClawWatch

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

访问：**http://localhost:3939**

### 使用 Docker 命令

```bash
cd ~/.openclaw/ClawWatch

# 构建镜像
docker build -t openclaw-viewer .

# 运行容器
docker run -d \
  --name openclaw-viewer \
  -p 3939:3939 \
  -v ~/.openclaw:/root/.openclaw:ro \
  openclaw-viewer

# 查看日志
docker logs -f openclaw-viewer

# 停止并删除容器
docker stop openclaw-viewer && docker rm openclaw-viewer
```

## 📖 使用指南

### 界面说明

打开浏览器访问 http://localhost:3939，你会看到：

1. **顶部控制栏**
   - Agent 选择器：切换不同的 agent
   - 刷新按钮：手动刷新会话
   - 过滤按钮：All / 👤 User / 🤖 Assistant
   - Show Thinking：显示/隐藏思考过程
   - Auto Scroll：自动滚动到最新消息
   - 搜索框：全文搜索

2. **消息列表**
   - 用户消息：蓝色标记 👤
   - 助手消息：绿色标记 🤖
   - 思考块：点击展开/折叠 💭
   - 工具调用：黄色标记 🔧

3. **状态指示器**
   - 🟢 绿点：WebSocket 已连接，实时监控中
   - 🔴 红点：连接断开，自动重连中

### 实时监控

- 服务会自动监控最新的会话文件
- 有新消息时会实时推送到浏览器
- 无需手动刷新页面

### 思考过程

- 默认显示思考块
- 点击 "💭 Thinking" 可展开/折叠详细内容
- 取消勾选 "Show Thinking" 可隐藏所有思考块

## 🛠️ 配置

### 修改端口

编辑 `docker-compose.yml`：

```yaml
ports:
  - "8080:3939"  # 将宿主机端口改为 8080
```

然后重启：

```bash
docker-compose down
docker-compose up -d
```

### 自定义样式

前端页面使用 CSS 变量，可以轻松自定义主题。

编辑 `web-viewer.html`，找到 `:root` 部分：

```css
:root {
  --bg-primary: #1a1a1a;      /* 主背景色 */
  --bg-secondary: #2a2a2a;    /* 次背景色 */
  --accent: #1890ff;          /* 强调色 */
  --user-color: #4a9eff;      /* 用户消息颜色 */
  --assistant-color: #52c41a; /* 助手消息颜色 */
  /* ... 更多配置 */
}
```

## 📝 快捷命令（可选）

添加到 `~/.bashrc` 或 `~/.zshrc`：

```bash
# 启动 Web 查看器
alias vsw-web="cd ~/.openclaw/ClawWatch && docker-compose up -d && echo '🌐 http://localhost:3939'"

# 停止 Web 查看器
alias vsw-stop="cd ~/.openclaw/ClawWatch && docker-compose down"

# 查看日志
alias vsw-logs="docker logs -f openclaw-session-viewer"

# 重启服务
alias vsw-restart="cd ~/.openclaw/ClawWatch && docker-compose restart"

# 重新构建并启动
alias vsw-rebuild="cd ~/.openclaw/ClawWatch && docker-compose up -d --build"
```

使配置生效：

```bash
source ~/.bashrc  # 或 source ~/.zshrc
```

之后就可以使用：

```bash
vsw-web      # 启动服务
vsw-logs     # 查看日志
vsw-stop     # 停止服务
vsw-restart  # 重启服务
```

## 🔧 故障排除

### 1. 端口已被占用

```bash
# 查看端口占用
lsof -i :3939

# 解决方法：修改 docker-compose.yml 中的端口，或停止占用端口的程序
```

### 2. 镜像拉取失败

```bash
# 测试能否访问内部镜像仓库
docker pull hub.intra.mlamp.cn/public/node:latest

# 如果失败，检查网络或 VPN 连接
```

### 3. 无法读取会话文件

检查 Docker 是否有权限访问 `~/.openclaw` 目录：

- **macOS**: Docker Desktop → Settings → Resources → File Sharing
- **Linux**: 检查目录权限

### 4. WebSocket 连接失败

- 检查防火墙设置
- 确认容器正常运行：`docker ps`
- 查看容器日志：`docker logs openclaw-session-viewer`

### 5. 页面显示空白

- 检查浏览器控制台是否有错误
- 确认 agent 目录下有会话文件
- 尝试刷新页面或切换 agent

## 📂 项目结构

```
~/.openclaw/ClawWatch/
├── web-viewer-server.js   # 后端服务器（Node.js + WebSocket）
├── web-viewer.html         # 前端页面（HTML + CSS + JS）
├── package.json            # Node.js 依赖配置
├── Dockerfile              # Docker 镜像定义
├── docker-compose.yml      # Docker Compose 配置
├── .dockerignore           # Docker 构建忽略文件
└── README.md               # 本文档
```

## 🌟 技术栈

- **后端**: Node.js + HTTP + WebSocket (ws)
- **前端**: 原生 JavaScript + HTML5 + CSS3
- **容器**: Docker + Docker Compose
- **基础镜像**: hub.intra.mlamp.cn/public/node:latest

## 🔄 开发说明

### 本地开发（不使用 Docker）

```bash
cd ~/.openclaw/ClawWatch

# 安装依赖
npm install

# 启动服务
node web-viewer-server.js

# 访问 http://localhost:3939
```

### 修改代码后重启

```bash
# 方式 1: 重启容器
docker-compose restart

# 方式 2: 重新构建（代码有重大改动时）
docker-compose up -d --build
```

## 📜 License

MIT

---

**问题反馈**: 如有问题或建议，请联系项目维护者。

---

## 🖥️ 本地部署（非容器化）

如果你不想使用容器，可以直接在本地运行 Claw Watch。

### 前置要求
- Node.js 18+ 
- npm 或 pnpm

### 快速启动

**1. 前台运行（适合开发/测试）**
```bash
./start-local.sh
```
按 `Ctrl+C` 停止。

**2. 后台运行（适合生产）**
```bash
# 启动
./start-daemon.sh

# 查看状态
./status.sh

# 查看日志
tail -f clawwatch.log

# 停止
./stop-daemon.sh

# 重启
./restart-daemon.sh
```

### 配置端口
```bash
# 默认端口 3939
PORT=8080 ./start-local.sh

# 或
export PORT=8080
./start-daemon.sh
```

### 对比：容器 vs 本地

| 特性 | 容器化 (podman-compose) | 本地部署 |
|------|----------------------|---------|
| 启动速度 | 慢（需要构建镜像） | 快 |
| 资源占用 | 高 | 低 |
| 隔离性 | 强 | 弱 |
| 依赖管理 | 自动 | 手动 |
| 适用场景 | 生产环境、多实例 | 开发、单实例 |

### 常见问题

**Q: 端口被占用？**
```bash
# 查看占用端口的进程
lsof -i:3939

# 或换个端口
PORT=8080 ./start-daemon.sh
```

**Q: 启动失败？**
```bash
# 查看日志
cat clawwatch.log

# 检查依赖
npm install
```

**Q: 如何开机自启？**

**macOS (launchd):**
```bash
# 创建 plist 文件
cat > ~/Library/LaunchAgents/com.openclaw.clawwatch.plist << PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.openclaw.clawwatch</string>
    <key>ProgramArguments</key>
    <array>
        <string>/bin/bash</string>
        <string>$HOME/.openclaw/ClawWatch/start-daemon.sh</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
PLIST

# 加载
launchctl load ~/Library/LaunchAgents/com.openclaw.clawwatch.plist
```

**Linux (systemd):**
```bash
# 创建服务文件
sudo tee /etc/systemd/system/clawwatch.service << SERVICE
[Unit]
Description=Claw Watch
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$HOME/.openclaw/ClawWatch
ExecStart=/usr/bin/node web-viewer-server.js
Restart=always

[Install]
WantedBy=multi-user.target
SERVICE

# 启用
sudo systemctl enable clawwatch
sudo systemctl start clawwatch
```

