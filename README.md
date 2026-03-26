# ClawWatch 🦞

OpenClaw 实时会话监控器 - 现代化的 Web 界面，用于实时查看和监控 OpenClaw 会话日志。

## ✨ 功能特性

**📊 实时监控**
- 自动轮询更新（3秒刷新）
- 多 Agent 支持
- Session Tabs（最多5个活跃会话）
- 🔴 未读消息提示

**🎯 消息展示**
- Markdown 渲染（AI 回复）
- Tool 执行结果可视化（✅成功/❌失败）
- Thinking 过程展示（可折叠）
- User/Assistant 视觉区分

**🔧 便捷操作**
- 批量展开/折叠 Thinking 和 Args
- 消息排序（正序/倒序）
- 角色过滤和搜索
- Tool Result 整合显示

## 🚀 快速开始

```bash
cd ~/.openclaw/ClawWatch

# 启动服务（自动构建前端）
./clawwatch.sh

# 或者明确指定命令
./clawwatch.sh start
```

访问：**http://localhost:3939**

### 管理命令

```bash
./clawwatch.sh start     # 启动服务
./clawwatch.sh stop      # 停止服务
./clawwatch.sh restart   # 重启服务
./clawwatch.sh status    # 查看状态
./clawwatch.sh logs      # 查看日志
./clawwatch.sh help      # 帮助信息
```

## 🏗️ 技术栈

**前端：**
- Vue 3 (Composition API)
- Vite (构建工具)
- Tailwind CSS v4 (样式)
- Marked (Markdown 渲染)

**后端：**
- Node.js
- HTTP API (RESTful)
- 文件系统监控

## 📁 项目结构

```
ClawWatch/
├── frontend/              # Vue 3 前端
│   ├── src/
│   │   ├── App.vue        # 主应用
│   │   ├── components/
│   │   │   └── MessageCard.vue
│   │   └── style.css
│   ├── package.json
│   └── vite.config.js
├── web-viewer-server.js   # Node.js API 服务器
├── start-daemon.sh        # 启动后端脚本
├── stop-daemon.sh         # 停止后端脚本
├── restart-daemon.sh      # 重启后端脚本
├── status.sh              # 查看状态
└── package.json
```

## 🔌 架构说明

**单端口架构：**
- 端口 3939 同时提供前端和 API
- 前端：静态资源（`frontend/dist/`）
- API：`/api/*` 路径
- Node.js 后端同时处理静态文件和 API 请求

### API 接口

- `GET /api/agents` - 获取所有 agent 列表
- `GET /api/sessions?agent=main` - 获取指定 agent 的所有 sessions
- `GET /api/latest-session?agent=main` - 获取最新 session
- `GET /api/session?path=<path>` - 读取指定 session 的内容

## 📖 使用指南

### Session Tabs

顶部显示最近 5 个活跃 sessions：
- 🔴 红点标记最新 session
- 点击切换到不同 session
- 有新消息的 session 会显示红点提示

### 消息过滤

侧边栏提供过滤选项：
- 👤 User - 用户消息
- 🤖 Assistant - AI 回复
- 🔧 Tools - 工具调用
- 💭 Thinking - 思考过程

### 批量操作

顶部工具栏：
- 🔽 Newest First / 🔼 Oldest First - 切换排序
- 💭 Expand/Collapse All Thinking - 批量展开/折叠思考
- 🔧 Show/Hide All Args - 批量展开/折叠工具参数

### Tool 执行状态

工具调用卡片会根据执行结果改变颜色：
- 🟡 黄色 - 等待执行
- 🟢 绿色 - 执行成功（✅）
- 🔴 红色 - 执行失败（❌）

## 🛠️ 开发

### 修改代码

```bash
# 1. 修改前端或后端代码
# 2. 重启服务（自动重新构建）
./clawwatch.sh restart

# 查看日志
./clawwatch.sh logs

# 或实时查看
tail -f clawwatch.log
```

### 首次安装

```bash
cd frontend
npm install
```

## 📝 开机自启（可选）

### macOS (launchd)

```bash
# 创建 plist 文件（替换 YOUR_USERNAME）
cat > ~/Library/LaunchAgents/com.openclaw.clawwatch.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.openclaw.clawwatch</string>
    <key>ProgramArguments</key>
    <array>
        <string>/Users/YOUR_USERNAME/.openclaw/ClawWatch/clawwatch.sh</string>
        <string>start</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/Users/YOUR_USERNAME/.openclaw/ClawWatch/clawwatch.log</string>
    <key>StandardErrorPath</key>
    <string>/Users/YOUR_USERNAME/.openclaw/ClawWatch/clawwatch.log</string>
</dict>
</plist>
EOF

# 加载服务
launchctl load ~/Library/LaunchAgents/com.openclaw.clawwatch.plist
```

### Linux (systemd)

```bash
# 创建服务文件（替换 YOUR_USERNAME）
sudo tee /etc/systemd/system/clawwatch.service << 'EOF'
[Unit]
Description=ClawWatch - OpenClaw Session Monitor
After=network.target

[Service]
Type=simple
User=YOUR_USERNAME
WorkingDirectory=/home/YOUR_USERNAME/.openclaw/ClawWatch
ExecStart=/home/YOUR_USERNAME/.openclaw/ClawWatch/clawwatch.sh start
Restart=always

[Install]
WantedBy=multi-user.target
EOF

# 启用并启动
sudo systemctl daemon-reload
sudo systemctl enable clawwatch
sudo systemctl start clawwatch
```

## 🐛 故障排查

### 服务无法启动

```bash
# 检查状态
./clawwatch.sh status

# 查看日志
./clawwatch.sh logs

# 检查端口占用
lsof -i :3939

# 停止旧进程
./clawwatch.sh stop
```

### 页面无法访问

```bash
# 1. 确认服务已启动
./clawwatch.sh status

# 2. 测试 API
curl http://localhost:3939/api/agents

# 3. 重启服务
./clawwatch.sh restart
```

### Sessions 不更新

1. 检查 OpenClaw 会话路径：`~/.openclaw/agents/main/sessions/`
2. 确认有写入权限
3. 重启服务：`./clawwatch.sh restart`

## 📄 许可证

MIT

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 支持

- 仓库：git@codex.mlamp.cn:paas/ClawWatch.git
- 问题反馈：创建 Issue
