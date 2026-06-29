# ClawWatch 使用文档

## 概述

ClawWatch 是 OpenClaw 实时会话监控器，通过 Web 界面实时查看 AI Agent 的运行日志、工具调用和会话状态。

## 访问地址

- **Hub Console**: https://clawhub.intra.mlamp.cn

## 界面说明

### 顶部状态栏

- **连接状态**: 绿色 "Connected" 表示 WebSocket 连接正常
- **统计信息**: 显示当前 hosts 数量、agents 数量、sessions 数量

### 控制栏

| 字段 | 说明 |
|------|------|
| Token | 输入主机 token 查看对应主机数据，输入管理员 token 查看全部 |
| Subscribe | 订阅过滤，`*` 表示全部，可输入特定 agentId |
| Reconnect | 重新建立 WebSocket 连接 |
| Clear | 清空当前日志显示 |

### 消息展示

- **User 消息**: 用户输入
- **Assistant 消息**: AI 回复，支持 Markdown 渲染
- **Tool 调用**: 工具执行结果（🟢成功 / 🔴失败）
- **Thinking**: AI 推理过程（可折叠）

### 功能按钮

- **Expand thinking**: 展开所有 thinking 块
- **Show tool details**: 展开工具调用参数和返回值
- **Settings**: 显示设置

## Token 使用

### 获取 Token

1. Agent 首次连接 Hub 时自动生成 per-agent token
2. Agent 进程日志会打印: `[agent] Console token for this host: <token>`
3. 或者请管理员从 Hub 服务端查询: `cat /app/data/tokens.json`

### Token 类型

| 类型 | 权限 | 用途 |
|------|------|------|
| Per-Agent Token | 只能看到对应主机的数据 | 分发给各主机使用者 |
| Admin Token (HUB_CONSOLE_TOKEN) | 看到所有主机数据 | 管理员使用 |
| 留空 | 被拒绝 (401) | — |

### 使用流程

1. 在 Token 输入框填入你的 token
2. 点击 **Reconnect**
3. 页面刷新后只显示你有权限的主机数据

## REST API

所有 REST API 均通过 `?token=<your-token>` 参数鉴权。

| 端点 | 说明 |
|------|------|
| `GET /api/health` | 健康检查（无需 token） |
| `GET /api/agents?token=xxx` | 在线 Agent 列表 |
| `GET /api/history?token=xxx&agentId=*&limit=500` | 历史日志消息 |
| `GET /api/trajectory?token=xxx&agentId=*&limit=500` | Session 运行轨迹事件 |
| `GET /api/session-history?token=xxx&agentId=X&openclawAgentId=Y&sessionId=Z` | 获取指定 session 的对话历史 |

## 常见问题

**Q: 不输入 token 为什么看不到数据？**
A: Hub 配置了全局 ConsoleToken 后，未认证请求会被拒绝（401）。

**Q: 如何获取我的主机 token？**
A: 重启 agent 进程查看日志输出，或联系管理员查看 Hub 上的 `/app/data/tokens.json`。

**Q: 连接显示 Connected 但看不到日志？**
A: 检查 token 是否正确，以及你的 agent 是否在线运行。
