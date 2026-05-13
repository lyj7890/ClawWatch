# ClawWatch Hub

ClawWatch 远端中继服务器（Go 实现），支持多 Agent 接入和浏览器实时查看。

## 快速启动

```bash
# 编译
make build

# 运行（开发模式，无 token）
./clawwatch-hub --port 4848

# 生产模式（带 token）
./clawwatch-hub --port 4848 \
  --agent-token your-agent-secret \
  --console-token your-console-secret

# 环境变量方式
HUB_PORT=4848 \
HUB_AGENT_TOKEN=xxx \
HUB_CONSOLE_TOKEN=xxx \
./clawwatch-hub
```

访问 http://localhost:4848 打开 Console UI。

## 连接协议

### Agent 端（推送日志）

```
ws://hub:4848/ws/agent?token=<agent-token>&agentId=<your-agent-id>
```

发送 JSON 消息，Hub 自动注入 `agentId` 和 `timestamp` 后广播给订阅的 Console。

### Console 端（查看日志）

```
ws://hub:4848/ws/console?token=<console-token>&subscribe=*
ws://hub:4848/ws/console?token=<console-token>&subscribe=my-mac-pro
```

`subscribe=*` 接收所有 Agent 的消息，也可指定 agentId 只看某台。

## 消息格式

Agent 推送：
```json
{ "type": "log", "session": "xxx.jsonl", "lines": [...] }
```
Hub 转发给 Console 时自动添加：
```json
{ "type": "log", "agentId": "my-mac", "session": "xxx.jsonl", "lines": [...], "timestamp": 1715608800000 }
```

Agent 状态通知（Hub 自动生成）：
```json
{ "type": "agent_status", "agentId": "my-mac", "online": true, "timestamp": 1715608800000 }
```

## API

| 端点 | 说明 |
|------|------|
| `GET /api/health` | 健康检查 |
| `GET /api/agents` | 当前在线 Agent 列表 |
| `GET /` | Console Web UI |

## 部署

```bash
# 编译 Linux 二进制（在 Mac 上交叉编译）
make build-linux

# Docker
docker build -t clawwatch-hub .
docker run -p 4848:4848 \
  -e HUB_AGENT_TOKEN=xxx \
  -e HUB_CONSOLE_TOKEN=xxx \
  clawwatch-hub
```
