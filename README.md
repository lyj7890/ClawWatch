# ClawWatch 🦞

OpenClaw / Hermes 实时会话监控器 - 现代化的 Web 界面，用于实时查看和监控会话日志。

支持 OpenClaw 和 Hermes Agent 两种模式。

## 🚀 快速开始

```bash
# OpenClaw 模式
npx @lyj7890/claw-watch

# Hermes 模式
npx @lyj7890/claw-watch start --hermes
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

### Session 列表
Sidebar 左侧按 **agent 主机名**分组展示所有 sessions，点击 agent 名可展开/折叠。Session 列表无数量限制，可无限滚动。🔴 标记最新 session，有新消息时显示红色未读 badge。

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

Sessions 不更新？检查会话路径：

- OpenClaw：`~/.openclaw/agents/main/sessions/`
- Hermes：`~/.hermes/sessions/`

## 🗺️ Roadmap

- [ ] **远端模式（Hub-and-Spoke）**
  - `clawwatch agent` — 本地监控，主动推送到 Hub
  - `clawwatch hub` — 公网中继服务器，支持多 Agent 接入
  - `clawwatch console` — 浏览器连接 Hub 远程查看
  - 适合手机随时查看 Agent 状态，或多人共享 Agent 运行情况

## 📄 许可证

MIT
# Test webhook trigger - Wed May 20 22:41:42 CST 2026
# Webhook delivery test v2 - Wed May 20 22:48:20 CST 2026
# Final webhook test - Wed May 20 22:51:11 CST 2026
