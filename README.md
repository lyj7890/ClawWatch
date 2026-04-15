# ClawWatch 🦞

OpenClaw 实时会话监控器 - 现代化的 Web 界面，用于实时查看和监控 OpenClaw 会话日志。

## 🚀 快速开始

```bash
npx @lyj7890/claw-watch
```

访问：**http://localhost:3939**

或全局安装后使用：

```bash
npm i -g @lyj7890/claw-watch

clawwatch start    # 启动
clawwatch stop     # 停止
clawwatch restart  # 重启
clawwatch status   # 查看状态
clawwatch logs     # 查看日志
```

## ✨ 功能特性

**📊 实时监控**
- WebSocket 实时推送，无需手动刷新
- 多 Agent 支持
- Session Tabs（最多 5 个活跃会话）
- 🔴 未读消息提示

**🎯 消息展示**
- Markdown 渲染（AI 回复）
- Tool 执行结果可视化（✅成功 / ❌失败）
- Thinking 过程展示（可折叠）
- User / Assistant 视觉区分

**🔧 便捷操作**
- 批量展开/折叠 Thinking 和 Args
- 消息排序（正序/倒序）
- 角色过滤
- Tool Result 整合显示

## 🏗️ 技术栈

- **前端**：Vue 3 + Vite + Tailwind CSS v4
- **后端**：Node.js + WebSocket
- **架构**：单端口（3939），前后端合一

## 🔌 API 接口

| 接口 | 说明 |
|------|------|
| `GET /api/agents` | 获取所有 agent 列表 |
| `GET /api/sessions?agent=main` | 获取指定 agent 的所有 sessions |
| `GET /api/latest-session?agent=main` | 获取最新 session |
| `GET /api/session?path=<path>` | 读取指定 session 内容 |

## 📖 使用说明

### Session Tabs
顶部显示最近 5 个活跃 sessions，🔴 红点标记最新，点击可切换。

### 消息过滤
侧边栏支持按角色过滤：👤 User / 🤖 Assistant / 🔧 Tools / 💭 Thinking

### 批量操作
- 🔽 / 🔼 切换消息排序
- 💭 批量展开/折叠 Thinking
- 🔧 批量展开/折叠 Args

### Tool 执行状态
- 🟡 黄色 — 等待执行
- 🟢 绿色 — 执行成功
- 🔴 红色 — 执行失败

## 🐛 故障排查

```bash
# 检查状态
clawwatch status

# 检查端口占用
lsof -i :3939

# 查看日志
clawwatch logs
```

Sessions 不更新？检查 OpenClaw 会话路径：`~/.openclaw/agents/main/sessions/`

## 📄 许可证

MIT
