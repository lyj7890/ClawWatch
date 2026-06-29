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
HUB_STORAGE_PATH=data/messages.jsonl \
HUB_HISTORY_LIMIT=2000 \
HUB_METRICS_PATH=data/metrics.jsonl \
HUB_METRICS_LIMIT=10000 \
./clawwatch-hub
```

访问 http://localhost:4848 打开 Console UI。

## 连接协议

### Go Agent（推荐）

Agent 使用独立 Go 二进制，默认监听 `~/.openclaw/agents`：

```bash
cd ../agent
go build -o clawwatch-agent .

./clawwatch-agent --hub wss://clawatch.intra.mlamp.cn
```

首次启动会在 `~/.clawwatch/host-id` 自动生成稳定且唯一的主机 ID。
hostname、IP、OS 和架构由 Agent 自动采集。

```bash
./clawwatch-agent --hub wss://clawatch.intra.mlamp.cn
```

开发阶段仅保留两个参数：

| 参数 | 说明 |
|------|------|
| `--hub` | Hub 地址，必需 |
| `--dir` | OpenClaw agents 目录，默认 `~/.openclaw/agents` |

### Agent 端协议

```
ws://hub:4848/ws/agent?token=<agent-token>&agentId=<your-agent-id>
```

协议 v2 上报三层身份：

- 主机：`host.id`、`host.hostname`、`host.ips`、OS、架构
- OpenClaw Agent：`openclawAgentId`
- Session：`sessionId`

```json
{
  "type": "log",
  "protocolVersion": 2,
  "host": {
    "id": "prod-mac-mini-a1b2c3",
    "hostname": "prod-mac-mini",
    "ips": ["10.0.0.12"],
    "os": "darwin",
    "arch": "arm64"
  },
  "openclawAgentId": "bajie",
  "sessionId": "c391fc75-e048-44c1-9c3e-d51b4ffcdded",
  "lines": []
}
```

Hub 会继续兼容旧 Node Agent 消息。

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
| `GET /api/history?agentId=*&limit=500` | 最近持久化的日志消息 |
| `GET /api/trajectory?agentId=*&limit=500` | 最近持久化的 session trajectory 运行事件 |
| `GET /` | Console Web UI |

Hub 默认将最近 2000 个 `log` 消息信封写入 `data/messages.jsonl`。消息文件权限为
`0600`，目录权限为 `0700`；Kubernetes 部署使用 PVC 挂载到 `/data`。
Trajectory 运行事件独立写入 `data/metrics.jsonl`，默认保留最近 10000 条，不会挤占对话历史。

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
