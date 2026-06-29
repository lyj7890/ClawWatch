# ClawWatch 架构与原理

## 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        ClawWatch Hub                             │
│                     (Go, 单进程, 端口 4848)                       │
│                                                                 │
│  ┌──────────┐   ┌──────────┐   ┌───────────┐   ┌────────────┐  │
│  │ AgentMgr │   │ ConsoleMgr│   │ TokenStore│   │ MessageStore│  │
│  └────┬─────┘   └─────┬────┘   └─────┬─────┘   └─────┬──────┘  │
│       │                │              │                │         │
│       │    ┌───────────┴──────┐       │                │         │
│       │    │    Hub (fanout)  │───────┘                │         │
│       │    └──────────────────┘                        │         │
└───────┼────────────────┼───────────────────────────────┼─────────┘
        │                │                               │
   WebSocket          WebSocket                     JSONL Files
   /ws/agent          /ws/console                   (持久化)
        │                │
        │                │
┌───────┴───────┐  ┌─────┴──────────┐
│ ClawWatch Agent│  │  浏览器 Console │
│  (Go, 每台主机) │  │  (内置前端 UI)  │
│               │  └────────────────┘
│ - 扫描 session│
│ - 增量推送日志 │
│ - 推送 trajectory│
└───────────────┘
        │
        │ (读取本地文件)
        │
┌───────┴───────┐
│   OpenClaw    │
│ ~/.openclaw/  │
│  agents/      │
│   sessions/   │
└───────────────┘
```

## 核心组件

### Hub（中继服务器）

- **语言**: Go
- **职责**: 接收 Agent 上报，按权限分发给 Console
- **核心结构**: `Hub` struct，包含 agents map、consoles map、broadcast channel
- **并发模型**: 单 goroutine 事件循环 (`Hub.Run()`) + 每连接两个 goroutine (read/write pump)

### Agent（采集端）

- **语言**: Go
- **职责**: 扫描本地 OpenClaw session 文件，增量推送日志
- **数据源**: 
  - Session JSONL 文件（对话日志）
  - `.trajectory.jsonl` 文件（运行轨迹：模型调用、工具执行等）
- **上报间隔**: 1 秒轮询文件变更

### Console（查看端）

- **形式**: Hub 内置的 Web UI（`index_html.go`）
- **连接**: WebSocket 到 Hub，接收实时推送
- **功能**: 实时日志、Agent 状态、Session 选择、消息过滤

## 通信协议

### Agent → Hub (WebSocket)

**连接**: `ws://hub:4848/ws/agent?token=<agent-token>&agentId=<host-id>`

**握手消息 (agent_hello)**:
```json
{
  "type": "agent_hello",
  "protocolVersion": 3,
  "agentVersion": "0.3.1",
  "host": {
    "id": "prod-mac-mini-a1b2c3",
    "hostname": "prod-mac-mini",
    "ips": ["10.0.0.12"],
    "os": "darwin",
    "arch": "arm64"
  },
  "agents": [
    {"id": "main", "sessionCount": 5, "lastActivity": 1715608800000}
  ]
}
```

**日志消息 (log)**:
```json
{
  "type": "log",
  "openclawAgentId": "main",
  "sessionId": "c391fc75-...",
  "lines": [
    {"role": "user", "content": "hello", "timestamp": "..."}
  ]
}
```

**轨迹消息 (trajectory)**:
```json
{
  "type": "trajectory",
  "eventType": "tool.call",
  "openclawAgentId": "main",
  "sessionId": "c391fc75-...",
  "data": { "name": "exec", "status": "success" }
}
```

**Session 列表 (session_list)**:
```json
{
  "type": "session_list",
  "sessions": [
    {"agentId": "main", "sessionId": "xxx", "mtime": 1715608800, "size": 4096}
  ]
}
```

### Hub → Agent (WebSocket)

**Token 分配 (agent_hello_ack)**:
```json
{
  "type": "agent_hello_ack",
  "agentId": "prod-mac-mini-a1b2c3",
  "consoleToken": "d24f34d27d9859dc1fb6798360598f38",
  "timestamp": 1715608800000
}
```

**Session 历史请求 (fetch_session)**:
```json
{
  "type": "fetch_session",
  "requestId": "req-1715608800000",
  "openclawAgentId": "main",
  "sessionId": "c391fc75-..."
}
```

### Console → Hub (WebSocket)

**连接**: `ws://hub:4848/ws/console?token=<console-token>&subscribe=*`

Console 目前只接收数据，不主动发送消息（read pump 仅用于保活）。

### Hub → Console (WebSocket)

Hub 将 Agent 上报的消息注入 `agentId`、`hostname` 等元数据后 fanout 给有权限的 Console：

```json
{
  "type": "log",
  "agentId": "prod-mac-mini-a1b2c3",
  "hostname": "prod-mac-mini",
  "hostIPs": ["10.0.0.12"],
  "openclawAgentId": "main",
  "sessionId": "c391fc75-...",
  "lines": [...],
  "timestamp": 1715608800000
}
```

## 数据流

```
[OpenClaw Session 文件]
        │ (Agent 增量读取)
        ▼
[Agent 解析 JSONL → 构造消息]
        │ (WebSocket)
        ▼
[Hub 接收 → 注入元数据]
        │
        ├──→ [MessageStore 持久化到 JSONL]
        │
        └──→ [fanout → 有权限的 Console]
                    │ (WebSocket)
                    ▼
              [浏览器渲染展示]
```

## Token 鉴权模型

详见 [Token 鉴权](./auth.md)。

## 持久化

| 文件 | 说明 | 默认限制 |
|------|------|---------|
| `data/messages.jsonl` | 日志消息 | 2000 条 |
| `data/metrics.jsonl` | 轨迹事件 | 10000 条 |
| `data/tokens.json` | AgentID → Token 映射 | 无限制 |

采用 append + truncate 策略：追加写入，超过限制时截断旧条目。

## 心跳与保活

- **Agent ↔ Hub**: WebSocket Ping/Pong，间隔 45s，超时 180s
- **Console ↔ Hub**: 同上
- **Hub 统计**: 每 30s 打印一次连接数

## 安全边界

- Agent → Hub: 可选 `HUB_AGENT_TOKEN` 校验
- Console → Hub: `HUB_CONSOLE_TOKEN` + per-agent token 双层校验
- 数据隔离: 非 admin Console 只能看到自己 token 对应的 Agent 数据
- Trajectory 脱敏: 默认不上报敏感内容（prompt、工具参数等）
